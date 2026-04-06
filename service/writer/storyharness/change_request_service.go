package storyharness

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangeRequestService 变更建议服务
// 负责建议的 CRUD 和状态流转
type ChangeRequestService struct {
	crRepo         writerRepo.ChangeRequestRepository
	projectionRepo writerRepo.ProjectionRepository
	characterRepo  writerRepo.CharacterRepository
}

// RebuildProjectionResult 章节投影重建结果。
type RebuildProjectionResult struct {
	ProjectID     string `json:"projectId"`
	ChapterID     string `json:"chapterId"`
	ReplayedCount int    `json:"replayedCount"`
	LastRequestID string `json:"lastRequestId,omitempty"`
}

// NewChangeRequestService 创建 ChangeRequestService 实例
func NewChangeRequestService(
	crRepo writerRepo.ChangeRequestRepository,
	projectionRepo writerRepo.ProjectionRepository,
	characterRepo writerRepo.CharacterRepository,
) *ChangeRequestService {
	return &ChangeRequestService{
		crRepo:         crRepo,
		projectionRepo: projectionRepo,
		characterRepo:  characterRepo,
	}
}

// ListByChapter 获取章节下的建议列表
func (s *ChangeRequestService) ListByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequest, error) {
	return s.crRepo.FindPendingByChapter(ctx, projectID, chapterID)
}

// GetByID 获取单条建议
func (s *ChangeRequestService) GetByID(ctx context.Context, id string) (*writer.ChangeRequest, error) {
	return s.crRepo.FindRequestByID(ctx, id)
}

// Process 处理建议（接受/忽略/延后）
func (s *ChangeRequestService) Process(ctx context.Context, requestID string, newStatus writer.ChangeRequestStatus, processedBy string) error {
	// 验证状态合法性
	switch newStatus {
	case writer.CRStatusAccepted, writer.CRStatusIgnored, writer.CRStatusDeferred:
		// valid
	default:
		return errors.NewServiceError("ChangeRequestService", errors.ServiceErrorValidation,
			"无效的处理状态", string(newStatus), nil)
	}

	request, err := s.crRepo.FindRequestByID(ctx, requestID)
	if err != nil {
		return errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"获取建议详情失败", requestID, err)
	}
	if request == nil {
		return errors.NewServiceError("ChangeRequestService", errors.ServiceErrorNotFound,
			"建议不存在", requestID, nil)
	}

	if err := s.crRepo.UpdateRequestStatus(ctx, requestID, newStatus, processedBy); err != nil {
		return errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"更新建议状态失败", requestID, err)
	}

	if newStatus != writer.CRStatusAccepted || s.projectionRepo == nil {
		return nil
	}

	projection, err := s.ensureProjection(ctx, request.ProjectID.Hex(), request.ChapterID.Hex())
	if err != nil {
		return err
	}

	s.applyAcceptedChange(projection, request)
	projection.Checkpoint = writer.ProjectionCheckpoint{
		LastRequestID: request.ID.Hex(),
		LastCategory:  request.Category,
		RefreshedAt:   time.Now(),
	}

	if err := s.projectionRepo.UpsertByChapter(ctx, projection); err != nil {
		return errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"刷新章节投影失败", requestID, err)
	}

	return nil
}

// CountPending 获取章节待处理建议数
func (s *ChangeRequestService) CountPending(ctx context.Context, projectID, chapterID string) (int64, error) {
	return s.crRepo.CountPendingByChapter(ctx, projectID, chapterID)
}

// RebuildProjection 重新基于 accepted 建议构建章节投影。
func (s *ChangeRequestService) RebuildProjection(ctx context.Context, projectID, chapterID string) (*RebuildProjectionResult, error) {
	if s.projectionRepo == nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"投影仓储不可用", chapterID, nil)
	}

	projection, err := s.buildBaseProjection(ctx, projectID, chapterID)
	if err != nil {
		return nil, err
	}

	accepted, err := s.crRepo.FindByChapterAndStatus(ctx, projectID, chapterID, writer.CRStatusAccepted)
	if err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"加载已接受建议失败", chapterID, err)
	}

	for _, request := range accepted {
		s.applyAcceptedChange(projection, request)
	}

	result := &RebuildProjectionResult{
		ProjectID:     projectID,
		ChapterID:     chapterID,
		ReplayedCount: len(accepted),
	}
	if len(accepted) > 0 {
		last := accepted[len(accepted)-1]
		projection.Checkpoint = writer.ProjectionCheckpoint{
			LastRequestID: last.ID.Hex(),
			LastCategory:  last.Category,
			RefreshedAt:   time.Now(),
		}
		result.LastRequestID = last.ID.Hex()
	} else {
		projection.Checkpoint = writer.ProjectionCheckpoint{
			RefreshedAt: time.Now(),
		}
	}

	if err := s.projectionRepo.UpsertByChapter(ctx, projection); err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"写入重建后的章节投影失败", chapterID, err)
	}

	return result, nil
}

func (s *ChangeRequestService) ensureProjection(ctx context.Context, projectID, chapterID string) (*writer.ChapterProjection, error) {
	projection, err := s.projectionRepo.GetByChapter(ctx, projectID, chapterID)
	if err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"获取章节投影失败", chapterID, err)
	}
	if projection != nil {
		return projection, nil
	}

	return s.buildBaseProjection(ctx, projectID, chapterID)
}

func (s *ChangeRequestService) buildBaseProjection(ctx context.Context, projectID, chapterID string) (*writer.ChapterProjection, error) {
	projectOID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorValidation,
			"无效的项目ID", projectID, err)
	}
	chapterOID, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorValidation,
			"无效的章节ID", chapterID, err)
	}

	projection := &writer.ChapterProjection{
		ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectOID},
		ChapterID:           chapterOID,
	}
	projection.TouchForCreate()

	if s.characterRepo == nil {
		return projection, nil
	}

	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"加载角色基线失败", projectID, err)
	}
	projection.Characters = make([]writer.CharacterSnapshot, 0, len(characters))
	for _, item := range characters {
		if item == nil {
			continue
		}
		projection.Characters = append(projection.Characters, writer.CharacterSnapshot{
			CharacterID:   item.ID.Hex(),
			CharacterName: item.Name,
			Summary:       item.Summary,
			CurrentState:  item.CurrentState,
		})
	}

	relations, err := s.characterRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		return nil, errors.NewServiceError("ChangeRequestService", errors.ServiceErrorInternal,
			"加载关系基线失败", projectID, err)
	}
	projection.Relations = make([]writer.RelationSnapshot, 0, len(relations))
	for _, item := range relations {
		if item == nil {
			continue
		}
		projection.Relations = append(projection.Relations, writer.RelationSnapshot{
			FromID:   item.FromID,
			ToID:     item.ToID,
			Relation: string(item.Type),
			Strength: item.Strength,
		})
	}

	return projection, nil
}

func (s *ChangeRequestService) applyAcceptedChange(projection *writer.ChapterProjection, request *writer.ChangeRequest) {
	if projection == nil || request == nil {
		return
	}

	switch request.Category {
	case writer.CRCategoryCharacterState:
		s.applyCharacterStateChange(projection, request)
	case writer.CRCategoryRelationChange:
		s.applyRelationChange(projection, request)
	}
}

func (s *ChangeRequestService) applyCharacterStateChange(projection *writer.ChapterProjection, request *writer.ChangeRequest) {
	characterID := readStringField(request.SuggestedChange, "characterId")
	if characterID == "" {
		return
	}

	snapshot := writer.CharacterSnapshot{
		CharacterID:   characterID,
		CharacterName: readStringField(request.SuggestedChange, "characterName"),
		CurrentState:  readStringField(request.SuggestedChange, "stateSummary"),
	}

	for idx := range projection.Characters {
		if projection.Characters[idx].CharacterID != characterID {
			continue
		}
		if snapshot.CharacterName == "" {
			snapshot.CharacterName = projection.Characters[idx].CharacterName
		}
		if snapshot.CurrentState == "" {
			snapshot.CurrentState = projection.Characters[idx].CurrentState
		}
		if projection.Characters[idx].Summary != "" {
			snapshot.Summary = projection.Characters[idx].Summary
		}
		projection.Characters[idx] = snapshot
		return
	}

	projection.Characters = append(projection.Characters, snapshot)
}

func (s *ChangeRequestService) applyRelationChange(projection *writer.ChapterProjection, request *writer.ChangeRequest) {
	fromID := readStringField(request.SuggestedChange, "fromId")
	toID := readStringField(request.SuggestedChange, "toId")
	if fromID == "" || toID == "" {
		return
	}

	snapshot := writer.RelationSnapshot{
		FromID:   fromID,
		ToID:     toID,
		FromName: readStringField(request.SuggestedChange, "fromName"),
		ToName:   readStringField(request.SuggestedChange, "toName"),
		Relation: readStringField(request.SuggestedChange, "relation"),
		Strength: readIntField(request.SuggestedChange, "strength"),
	}

	for idx := range projection.Relations {
		if projection.Relations[idx].FromID != fromID || projection.Relations[idx].ToID != toID {
			continue
		}
		if snapshot.FromName == "" {
			snapshot.FromName = projection.Relations[idx].FromName
		}
		if snapshot.ToName == "" {
			snapshot.ToName = projection.Relations[idx].ToName
		}
		if snapshot.Relation == "" {
			snapshot.Relation = projection.Relations[idx].Relation
		}
		if snapshot.Strength == 0 {
			snapshot.Strength = projection.Relations[idx].Strength
		}
		projection.Relations[idx] = snapshot
		return
	}

	projection.Relations = append(projection.Relations, snapshot)
}

func readStringField(payload map[string]interface{}, key string) string {
	if payload == nil {
		return ""
	}
	raw, ok := payload[key]
	if !ok || raw == nil {
		return ""
	}

	switch value := raw.(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}

func readIntField(payload map[string]interface{}, key string) int {
	if payload == nil {
		return 0
	}
	raw, ok := payload[key]
	if !ok || raw == nil {
		return 0
	}

	switch value := raw.(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	default:
		return 0
	}
}
