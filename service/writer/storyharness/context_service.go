package storyharness

import (
	"context"

	"Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContextService 章节上下文服务
// 负责组装 Context Lens 所需的角色快照、关系数据
type ContextService struct {
	characterRepo     writerRepo.CharacterRepository
	changeRequestRepo writerRepo.ChangeRequestRepository
	projectionRepo    writerRepo.ProjectionRepository
}

// NewContextService 创建 ContextService 实例
func NewContextService(
	characterRepo writerRepo.CharacterRepository,
	changeRequestRepo writerRepo.ChangeRequestRepository,
	projectionRepo writerRepo.ProjectionRepository,
) *ContextService {
	return &ContextService{
		characterRepo:     characterRepo,
		changeRequestRepo: changeRequestRepo,
		projectionRepo:    projectionRepo,
	}
}

// ChapterContextData 章节上下文数据
type ChapterContextData struct {
	Characters []*writer.Character
	Relations  []*writer.CharacterRelation
	PendingCRs int64
}

// GetChapterContext 获取章节上下文
// 第一版使用现有角色/关系数据作为基线
func (s *ContextService) GetChapterContext(ctx context.Context, projectID, chapterID string) (*ChapterContextData, error) {
	pendingCount, _ := s.changeRequestRepo.CountPendingByChapter(ctx, projectID, chapterID)

	if s.projectionRepo != nil {
		projection, err := s.projectionRepo.GetByChapter(ctx, projectID, chapterID)
		if err != nil {
			return nil, errors.NewServiceError("ContextService", errors.ServiceErrorInternal,
				"获取章节投影失败", chapterID, err)
		}
		if projection != nil {
			return &ChapterContextData{
				Characters: snapshotCharactersToModels(projection.Characters, projectID),
				Relations:  snapshotRelationsToModels(projection.Relations, projectID),
				PendingCRs: pendingCount,
			}, nil
		}
	}

	// 1. 获取项目角色列表
	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("ContextService", errors.ServiceErrorInternal,
			"获取角色列表失败", "", err)
	}

	// 2. 获取项目关系列表
	relations, err := s.characterRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		return nil, errors.NewServiceError("ContextService", errors.ServiceErrorInternal,
			"获取关系列表失败", "", err)
	}

	return &ChapterContextData{
		Characters: characters,
		Relations:  relations,
		PendingCRs: pendingCount,
	}, nil
}

func snapshotCharactersToModels(items []writer.CharacterSnapshot, projectID string) []*writer.Character {
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	results := make([]*writer.Character, 0, len(items))
	for _, item := range items {
		character := &writer.Character{
			ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectOID},
			Summary:             item.Summary,
			CurrentState:        item.CurrentState,
			ShortDescription:    item.Summary,
		}
		character.Name = item.CharacterName
		if oid, err := primitive.ObjectIDFromHex(item.CharacterID); err == nil {
			character.ID = oid
		}
		results = append(results, character)
	}
	return results
}

func snapshotRelationsToModels(items []writer.RelationSnapshot, projectID string) []*writer.CharacterRelation {
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	results := make([]*writer.CharacterRelation, 0, len(items))
	for _, item := range items {
		results = append(results, &writer.CharacterRelation{
			ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectOID},
			FromID:              item.FromID,
			ToID:                item.ToID,
			Type:                writer.RelationType(item.Relation),
			Strength:            item.Strength,
			Notes:               item.Relation,
		})
	}
	return results
}
