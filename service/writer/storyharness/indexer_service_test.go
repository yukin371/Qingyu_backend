package storyharness

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIndexerServiceTriggerChapterIndexPrefersAIResults(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	aliceID := primitive.NewObjectID()
	bobID := primitive.NewObjectID()

	documentRepo := &stubIndexerDocumentRepository{
		document: &writer.Document{
			IdentifiedEntity: writerBase.IdentifiedEntity{ID: chapterID},
			ProjectID:        projectID,
			Type:             writer.TypeChapter,
			Title:            "第一章",
		},
	}
	contentRepo := &stubIndexerDocumentContentRepository{
		content: &writer.DocumentContent{
			DocumentID:  chapterID,
			Content:     "林昭拔剑护住沈砚，青铜钥匙第一次发出微光。",
			ContentType: "markdown",
		},
	}
	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, aliceID, "林昭", "平静"),
			newTestCharacter(projectID, bobID, "沈砚", "警惕"),
		},
	}
	changeRepo := &stubIndexerChangeRequestRepository{}
	aiClient := &stubChapterAnalysisClient{
		response: &ChapterAnalysisResponse{
			StateChanges: []AIStateChange{
				{
					EntityName: "林昭",
					FieldKey:   "情绪",
					OldValue:   "平静",
					NewValue:   "紧张",
					Evidence:   "林昭拔剑护住沈砚",
				},
			},
			RelationChanges: []AIRelationChange{
				{
					FromEntity:   "林昭",
					ToEntity:     "沈砚",
					RelationType: "保护",
					ChangeType:   "new",
					Evidence:     "林昭拔剑护住沈砚",
				},
			},
			NewEntities: []AINewEntity{
				{
					Name:         "青铜钥匙",
					EntityType:   "item",
					FirstMention: "青铜钥匙第一次发出微光",
					Description:  "一把能发光的旧钥匙",
				},
			},
		},
	}

	service := NewIndexerService(documentRepo, contentRepo, characterRepo, changeRepo, aiClient)

	result, err := service.TriggerChapterIndex(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("TriggerChapterIndex returned error: %v", err)
	}

	if result.Source != "ai" {
		t.Fatalf("expected source ai, got %q", result.Source)
	}
	if result.Generated != 3 {
		t.Fatalf("expected 3 generated requests, got %d", result.Generated)
	}
	if len(changeRepo.createdRequests) != 3 {
		t.Fatalf("expected 3 stored requests, got %d", len(changeRepo.createdRequests))
	}

	stateRequest := changeRepo.createdRequests[0]
	if stateRequest.Category != writer.CRCategoryCharacterState {
		t.Fatalf("expected first request category %q, got %q", writer.CRCategoryCharacterState, stateRequest.Category)
	}
	if stateRequest.Source != "ai" {
		t.Fatalf("expected AI request source, got %q", stateRequest.Source)
	}
	if got := stateRequest.SuggestedChange["characterId"]; got != aliceID.Hex() {
		t.Fatalf("expected mapped characterId %s, got %v", aliceID.Hex(), got)
	}

	relationRequest := changeRepo.createdRequests[1]
	if relationRequest.Category != writer.CRCategoryRelationChange {
		t.Fatalf("expected relation change category, got %q", relationRequest.Category)
	}
	if got := relationRequest.SuggestedChange["fromId"]; got != aliceID.Hex() {
		t.Fatalf("expected fromId %s, got %v", aliceID.Hex(), got)
	}
	if got := relationRequest.SuggestedChange["toId"]; got != bobID.Hex() {
		t.Fatalf("expected toId %s, got %v", bobID.Hex(), got)
	}

	newEntityRequest := changeRepo.createdRequests[2]
	if newEntityRequest.Category != writer.CRCategoryScopeDrift {
		t.Fatalf("expected new entity request category %q, got %q", writer.CRCategoryScopeDrift, newEntityRequest.Category)
	}
	if got := newEntityRequest.SuggestedChange["entityType"]; got != "item" {
		t.Fatalf("expected entityType item, got %v", got)
	}
}

func TestIndexerServiceTriggerChapterIndexFallsBackToRuleEngine(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()

	documentRepo := &stubIndexerDocumentRepository{
		document: &writer.Document{
			IdentifiedEntity: writerBase.IdentifiedEntity{ID: chapterID},
			ProjectID:        projectID,
			Type:             writer.TypeChapter,
			Title:            "第一章",
		},
	}
	contentRepo := &stubIndexerDocumentContentRepository{
		content: &writer.DocumentContent{
			DocumentID:  chapterID,
			Content:     "林昭受伤后几乎站不稳。",
			ContentType: "markdown",
		},
	}
	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, characterID, "林昭", "平静"),
		},
	}
	changeRepo := &stubIndexerChangeRequestRepository{}
	aiClient := &stubChapterAnalysisClient{err: errors.New("ai unavailable")}

	service := NewIndexerService(documentRepo, contentRepo, characterRepo, changeRepo, aiClient)

	result, err := service.TriggerChapterIndex(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("TriggerChapterIndex returned error: %v", err)
	}

	if result.Source != "rule" {
		t.Fatalf("expected source rule, got %q", result.Source)
	}
	if result.Generated != 1 {
		t.Fatalf("expected 1 generated request, got %d", result.Generated)
	}
	if len(changeRepo.createdRequests) != 1 {
		t.Fatalf("expected 1 stored request, got %d", len(changeRepo.createdRequests))
	}
	if changeRepo.createdRequests[0].Source != "rule" {
		t.Fatalf("expected fallback request source rule, got %q", changeRepo.createdRequests[0].Source)
	}
}

func TestIndexerServiceTriggerChapterIndexDeduplicatesPendingRequests(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()

	documentRepo := &stubIndexerDocumentRepository{
		document: &writer.Document{
			IdentifiedEntity: writerBase.IdentifiedEntity{ID: chapterID},
			ProjectID:        projectID,
			Type:             writer.TypeChapter,
			Title:            "第一章",
		},
	}
	contentRepo := &stubIndexerDocumentContentRepository{
		content: &writer.DocumentContent{
			DocumentID:  chapterID,
			Content:     "林昭受伤后几乎站不稳。",
			ContentType: "markdown",
		},
	}
	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, characterID, "林昭", "平静"),
		},
	}
	changeRepo := &stubIndexerChangeRequestRepository{
		pendingRequests: []*writer.ChangeRequest{
			{
				Category: writer.CRCategoryCharacterState,
				Title:    "建议更新角色状态：林昭",
			},
		},
	}
	aiClient := &stubChapterAnalysisClient{
		response: &ChapterAnalysisResponse{
			StateChanges: []AIStateChange{
				{
					EntityName: "林昭",
					FieldKey:   "状态",
					NewValue:   "受伤",
					Evidence:   "林昭受伤后几乎站不稳。",
				},
			},
		},
	}

	service := NewIndexerService(documentRepo, contentRepo, characterRepo, changeRepo, aiClient)

	result, err := service.TriggerChapterIndex(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("TriggerChapterIndex returned error: %v", err)
	}

	if result.Generated != 0 {
		t.Fatalf("expected 0 generated requests after dedupe, got %d", result.Generated)
	}
	if result.Deduplicated != 1 {
		t.Fatalf("expected 1 deduplicated request, got %d", result.Deduplicated)
	}
	if len(changeRepo.createdRequests) != 0 {
		t.Fatalf("expected no newly stored requests, got %d", len(changeRepo.createdRequests))
	}
}

type stubChapterAnalysisClient struct {
	response *ChapterAnalysisResponse
	err      error
}

func (s *stubChapterAnalysisClient) AnalyzeChapter(ctx context.Context, request *ChapterAnalysisRequest) (*ChapterAnalysisResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}

type stubIndexerChangeRequestRepository struct {
	createdRequests []*writer.ChangeRequest
	pendingRequests []*writer.ChangeRequest
}

func (s *stubIndexerChangeRequestRepository) CreateRequest(ctx context.Context, cr *writer.ChangeRequest) error {
	s.createdRequests = append(s.createdRequests, cr)
	return nil
}

func (s *stubIndexerChangeRequestRepository) FindRequestByID(ctx context.Context, id string) (*writer.ChangeRequest, error) {
	return nil, nil
}

func (s *stubIndexerChangeRequestRepository) FindRequestsByBatchID(ctx context.Context, batchID string) ([]*writer.ChangeRequest, error) {
	return nil, nil
}

func (s *stubIndexerChangeRequestRepository) FindPendingByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequest, error) {
	return s.pendingRequests, nil
}

func (s *stubIndexerChangeRequestRepository) FindByChapterAndStatus(ctx context.Context, projectID, chapterID string, status writer.ChangeRequestStatus) ([]*writer.ChangeRequest, error) {
	return nil, nil
}

func (s *stubIndexerChangeRequestRepository) CountPendingByChapter(ctx context.Context, projectID, chapterID string) (int64, error) {
	return int64(len(s.pendingRequests)), nil
}

func (s *stubIndexerChangeRequestRepository) UpdateRequestStatus(ctx context.Context, id string, status writer.ChangeRequestStatus, processedBy string) error {
	return nil
}

func (s *stubIndexerChangeRequestRepository) DeleteRequest(ctx context.Context, id string) error {
	return nil
}

func (s *stubIndexerChangeRequestRepository) CreateBatch(ctx context.Context, batch *writer.ChangeRequestBatch) error {
	return nil
}

func (s *stubIndexerChangeRequestRepository) FindBatchByID(ctx context.Context, id string) (*writer.ChangeRequestBatch, error) {
	return nil, nil
}

func (s *stubIndexerChangeRequestRepository) FindBatchesByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequestBatch, error) {
	return nil, nil
}

func (s *stubIndexerChangeRequestRepository) UpdateBatchCounts(ctx context.Context, id string, total, pending int) error {
	return nil
}

type stubIndexerDocumentRepository struct {
	document *writer.Document
}

func (s *stubIndexerDocumentRepository) Create(ctx context.Context, entity *writer.Document) error {
	return nil
}

func (s *stubIndexerDocumentRepository) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	return s.document, nil
}

func (s *stubIndexerDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}

func (s *stubIndexerDocumentRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *stubIndexerDocumentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Document, error) {
	return nil, nil
}

func (s *stubIndexerDocumentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}

func (s *stubIndexerDocumentRepository) Exists(ctx context.Context, id string) (bool, error) {
	return s.document != nil, nil
}

func (s *stubIndexerDocumentRepository) Health(ctx context.Context) error {
	return nil
}

func (s *stubIndexerDocumentRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Document, error) {
	return nil, nil
}

func (s *stubIndexerDocumentRepository) GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*writer.Document, error) {
	return nil, nil
}

func (s *stubIndexerDocumentRepository) GetByIDs(ctx context.Context, ids []string) ([]*writer.Document, error) {
	return nil, nil
}

func (s *stubIndexerDocumentRepository) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	return nil
}

func (s *stubIndexerDocumentRepository) DeleteByProject(ctx context.Context, documentID, projectID string) error {
	return nil
}

func (s *stubIndexerDocumentRepository) RestoreByProject(ctx context.Context, documentID, projectID string) error {
	return nil
}

func (s *stubIndexerDocumentRepository) IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error) {
	return true, nil
}

func (s *stubIndexerDocumentRepository) SoftDelete(ctx context.Context, documentID, projectID string) error {
	return nil
}

func (s *stubIndexerDocumentRepository) HardDelete(ctx context.Context, documentID string) error {
	return nil
}

func (s *stubIndexerDocumentRepository) GetByIDUnscoped(ctx context.Context, id string) (*writer.Document, error) {
	return s.document, nil
}

func (s *stubIndexerDocumentRepository) CountByProject(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}

func (s *stubIndexerDocumentRepository) CreateWithTransaction(ctx context.Context, document *writer.Document, callback func(ctx context.Context) error) error {
	if callback != nil {
		return callback(ctx)
	}
	return nil
}

type stubIndexerDocumentContentRepository struct {
	content *writer.DocumentContent
}

func (s *stubIndexerDocumentContentRepository) Create(ctx context.Context, entity *writer.DocumentContent) error {
	return nil
}

func (s *stubIndexerDocumentContentRepository) GetByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	return s.content, nil
}

func (s *stubIndexerDocumentContentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}

func (s *stubIndexerDocumentContentRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *stubIndexerDocumentContentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.DocumentContent, error) {
	return nil, nil
}

func (s *stubIndexerDocumentContentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}

func (s *stubIndexerDocumentContentRepository) Exists(ctx context.Context, id string) (bool, error) {
	return s.content != nil, nil
}

func (s *stubIndexerDocumentContentRepository) Health(ctx context.Context) error {
	return nil
}

func (s *stubIndexerDocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (*writer.DocumentContent, error) {
	return s.content, nil
}

func (s *stubIndexerDocumentContentRepository) UpdateWithVersion(ctx context.Context, documentID string, updates map[string]interface{}, expectedVersion int) error {
	return nil
}

func (s *stubIndexerDocumentContentRepository) BatchUpdateContent(ctx context.Context, updates map[string]string) error {
	return nil
}

func (s *stubIndexerDocumentContentRepository) GetContentStats(ctx context.Context, documentID string) (int, int, error) {
	return 0, 0, nil
}

func (s *stubIndexerDocumentContentRepository) StoreToGridFS(ctx context.Context, documentID string, content []byte) (string, error) {
	return "", nil
}

func (s *stubIndexerDocumentContentRepository) LoadFromGridFS(ctx context.Context, gridFSID string) ([]byte, error) {
	return nil, nil
}

func (s *stubIndexerDocumentContentRepository) CreateWithTransaction(ctx context.Context, content *writer.DocumentContent, callback func(ctx context.Context) error) error {
	if callback != nil {
		return callback(ctx)
	}
	return nil
}
