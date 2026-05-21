package storyharness

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestChangeRequestServiceProcessAcceptedCharacterStateRefreshesProjection(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	requestID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()

	crRepo := &stubChangeRequestRepository{
		requestByID: map[string]*writer.ChangeRequest{
			requestID.Hex(): {
				IdentifiedEntity:    writerBase.IdentifiedEntity{ID: requestID},
				ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
				ChapterID:           chapterID,
				Category:            writer.CRCategoryCharacterState,
				SuggestedChange: map[string]interface{}{
					"characterId":   characterID.Hex(),
					"characterName": "林昭",
					"stateSummary":  "身体状态受损",
				},
			},
		},
	}
	projectionRepo := &stubProjectionRepository{}
	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, characterID, "林昭", "初始状态"),
		},
	}

	svc := NewChangeRequestService(crRepo, projectionRepo, characterRepo)

	if err := svc.Process(context.Background(), requestID.Hex(), writer.CRStatusAccepted, "user-1"); err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	if crRepo.updatedStatus != writer.CRStatusAccepted {
		t.Fatalf("expected updated status %q, got %q", writer.CRStatusAccepted, crRepo.updatedStatus)
	}
	if projectionRepo.saved == nil {
		t.Fatal("expected projection to be saved")
	}
	if got := projectionRepo.saved.Checkpoint.LastRequestID; got != requestID.Hex() {
		t.Fatalf("expected checkpoint request id %s, got %s", requestID.Hex(), got)
	}
	if len(projectionRepo.saved.Characters) != 1 {
		t.Fatalf("expected 1 projected character, got %d", len(projectionRepo.saved.Characters))
	}
	if got := projectionRepo.saved.Characters[0].CurrentState; got != "身体状态受损" {
		t.Fatalf("expected projected current state updated, got %q", got)
	}
}

func TestChangeRequestServiceProcessAcceptedRelationChangeRefreshesProjection(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	requestID := primitive.NewObjectID()

	crRepo := &stubChangeRequestRepository{
		requestByID: map[string]*writer.ChangeRequest{
			requestID.Hex(): {
				IdentifiedEntity:    writerBase.IdentifiedEntity{ID: requestID},
				ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
				ChapterID:           chapterID,
				Category:            writer.CRCategoryRelationChange,
				SuggestedChange: map[string]interface{}{
					"fromId":   "char-a",
					"toId":     "char-b",
					"fromName": "甲",
					"toName":   "乙",
					"relation": "关系改善",
					"strength": 70,
				},
			},
		},
	}
	projectionRepo := &stubProjectionRepository{
		projection: &writer.ChapterProjection{
			ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
			ChapterID:           chapterID,
			Relations: []writer.RelationSnapshot{
				{FromID: "char-a", ToID: "char-b", Relation: "关系恶化", Strength: 20},
			},
		},
	}

	svc := NewChangeRequestService(crRepo, projectionRepo, &stubCharacterRepository{})

	if err := svc.Process(context.Background(), requestID.Hex(), writer.CRStatusAccepted, "user-1"); err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	if projectionRepo.saved == nil {
		t.Fatal("expected projection to be saved")
	}
	if len(projectionRepo.saved.Relations) != 1 {
		t.Fatalf("expected 1 projected relation, got %d", len(projectionRepo.saved.Relations))
	}
	if got := projectionRepo.saved.Relations[0].Relation; got != "关系改善" {
		t.Fatalf("expected projected relation label updated, got %q", got)
	}
	if got := projectionRepo.saved.Relations[0].Strength; got != 70 {
		t.Fatalf("expected projected relation strength updated, got %d", got)
	}
}

func TestChangeRequestServiceProcessRejectsInvalidStatus(t *testing.T) {
	svc := NewChangeRequestService(&stubChangeRequestRepository{}, &stubProjectionRepository{}, &stubCharacterRepository{})

	err := svc.Process(context.Background(), primitive.NewObjectID().Hex(), writer.ChangeRequestStatus("unknown"), "user-1")
	if err == nil {
		t.Fatal("expected validation error for invalid status")
	}
}

func TestContextServiceGetChapterContextPrefersProjection(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()

	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, primitive.NewObjectID(), "旧基线", "旧状态"),
		},
		relations: []*writer.CharacterRelation{
			{
				ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
				FromID:              "baseline-a",
				ToID:                "baseline-b",
				Type:                writer.RelationEnemy,
				Strength:            10,
			},
		},
	}
	changeRepo := &stubChangeRequestRepository{pendingCount: 3}
	projectionRepo := &stubProjectionRepository{
		projection: &writer.ChapterProjection{
			ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
			ChapterID:           chapterID,
			Characters: []writer.CharacterSnapshot{
				{CharacterID: characterID.Hex(), CharacterName: "投影角色", CurrentState: "新状态"},
			},
			Relations: []writer.RelationSnapshot{
				{FromID: "char-a", ToID: "char-b", Relation: "关系改善", Strength: 80},
			},
		},
	}

	svc := NewContextService(characterRepo, changeRepo, projectionRepo)

	result, err := svc.GetChapterContext(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("GetChapterContext returned error: %v", err)
	}

	if got := len(result.Characters); got != 1 {
		t.Fatalf("expected 1 character from projection, got %d", got)
	}
	if got := result.Characters[0].Name; got != "投影角色" {
		t.Fatalf("expected projection character name, got %q", got)
	}
	if got := result.PendingCRs; got != 3 {
		t.Fatalf("expected pending count 3, got %d", got)
	}
	if got := result.Relations[0].Notes; got != "关系改善" {
		t.Fatalf("expected projection relation notes, got %q", got)
	}
}

func TestContextServiceGetChapterContextFallsBackToBaseline(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()

	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, characterID, "基线角色", "基线状态"),
		},
	}
	changeRepo := &stubChangeRequestRepository{pendingCount: 1}
	projectionRepo := &stubProjectionRepository{}

	svc := NewContextService(characterRepo, changeRepo, projectionRepo)

	result, err := svc.GetChapterContext(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("GetChapterContext returned error: %v", err)
	}

	if got := len(result.Characters); got != 1 {
		t.Fatalf("expected 1 baseline character, got %d", got)
	}
	if got := result.Characters[0].Name; got != "基线角色" {
		t.Fatalf("expected baseline character name, got %q", got)
	}
}

func TestChangeRequestServiceRebuildProjectionReplaysAcceptedRequests(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()
	requestID := primitive.NewObjectID()

	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, characterID, "林昭", "初始状态"),
		},
	}
	crRepo := &stubChangeRequestRepository{
		requestsByStatus: map[writer.ChangeRequestStatus][]*writer.ChangeRequest{
			writer.CRStatusAccepted: {
				{
					IdentifiedEntity:    writerBase.IdentifiedEntity{ID: requestID},
					ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
					ChapterID:           chapterID,
					Category:            writer.CRCategoryCharacterState,
					SuggestedChange: map[string]interface{}{
						"characterId":   characterID.Hex(),
						"characterName": "林昭",
						"stateSummary":  "体力明显下降",
					},
				},
			},
		},
	}
	projectionRepo := &stubProjectionRepository{}

	svc := NewChangeRequestService(crRepo, projectionRepo, characterRepo)

	result, err := svc.RebuildProjection(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("RebuildProjection returned error: %v", err)
	}

	if result.ReplayedCount != 1 {
		t.Fatalf("expected replayed count 1, got %d", result.ReplayedCount)
	}
	if projectionRepo.saved == nil {
		t.Fatal("expected projection to be saved")
	}
	if got := projectionRepo.saved.Characters[0].CurrentState; got != "体力明显下降" {
		t.Fatalf("expected rebuilt projection state updated, got %q", got)
	}
	if got := projectionRepo.saved.Checkpoint.LastRequestID; got != requestID.Hex() {
		t.Fatalf("expected checkpoint last request id %s, got %s", requestID.Hex(), got)
	}
}

func TestChangeRequestServiceRebuildProjectionAppliesAcceptedScopeDrift(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	requestID := primitive.NewObjectID()

	crRepo := &stubChangeRequestRepository{
		requestsByStatus: map[writer.ChangeRequestStatus][]*writer.ChangeRequest{
			writer.CRStatusAccepted: {
				{
					IdentifiedEntity:    writerBase.IdentifiedEntity{ID: requestID},
					ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
					ChapterID:           chapterID,
					Category:            writer.CRCategoryScopeDrift,
					SuggestedChange: map[string]interface{}{
						"action":      "create",
						"entityType":  string(writer.EntityTypeCharacter),
						"name":        "奶茶店老板",
						"description": "在街角经营奶茶店的新角色",
					},
				},
			},
		},
	}
	projectionRepo := &stubProjectionRepository{}

	svc := NewChangeRequestService(crRepo, projectionRepo, &stubCharacterRepository{})

	result, err := svc.RebuildProjection(context.Background(), projectID.Hex(), chapterID.Hex())
	if err != nil {
		t.Fatalf("RebuildProjection returned error: %v", err)
	}

	if result.ReplayedCount != 1 {
		t.Fatalf("expected replayed count 1, got %d", result.ReplayedCount)
	}
	if projectionRepo.saved == nil {
		t.Fatal("expected projection to be saved")
	}
	if got := len(projectionRepo.saved.Characters); got != 1 {
		t.Fatalf("expected 1 projected character, got %d", got)
	}
	if got := projectionRepo.saved.Characters[0].CharacterName; got != "奶茶店老板" {
		t.Fatalf("expected projected entity name %q, got %q", "奶茶店老板", got)
	}
	if got := projectionRepo.saved.Characters[0].CharacterID; got != requestID.Hex() {
		t.Fatalf("expected projected entity id %s, got %s", requestID.Hex(), got)
	}
	if got := projectionRepo.saved.Characters[0].EntityType; got != writer.EntityTypeCharacter {
		t.Fatalf("expected projected entity type %q, got %q", writer.EntityTypeCharacter, got)
	}
}

func TestChangeRequestServiceProcessAcceptedScopeDriftBackfillsAcceptedRelation(t *testing.T) {
	projectID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	existingCharacterID := primitive.NewObjectID()
	relationRequestID := primitive.NewObjectID()
	scopeRequestID := primitive.NewObjectID()

	relationRequest := &writer.ChangeRequest{
		IdentifiedEntity:    writerBase.IdentifiedEntity{ID: relationRequestID},
		ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
		ChapterID:           chapterID,
		Category:            writer.CRCategoryRelationChange,
		Status:              writer.CRStatusAccepted,
		SuggestedChange: map[string]interface{}{
			"fromId":   existingCharacterID.Hex(),
			"fromName": "诺艾尔",
			"toName":   "奶茶店老板",
			"relation": "债务往来",
			"strength": 65,
		},
	}
	scopeRequest := &writer.ChangeRequest{
		IdentifiedEntity:    writerBase.IdentifiedEntity{ID: scopeRequestID},
		ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
		ChapterID:           chapterID,
		Category:            writer.CRCategoryScopeDrift,
		Status:              writer.CRStatusPending,
		SuggestedChange: map[string]interface{}{
			"action":      "create",
			"entityType":  string(writer.EntityTypeCharacter),
			"name":        "奶茶店老板",
			"description": "新登场的店主",
		},
	}

	crRepo := &stubChangeRequestRepository{
		requestByID: map[string]*writer.ChangeRequest{
			relationRequestID.Hex(): relationRequest,
			scopeRequestID.Hex():    scopeRequest,
		},
		requestsByStatus: map[writer.ChangeRequestStatus][]*writer.ChangeRequest{
			writer.CRStatusAccepted: {relationRequest},
		},
	}
	characterRepo := &stubCharacterRepository{
		characters: []*writer.Character{
			newTestCharacter(projectID, existingCharacterID, "诺艾尔", "正在整理账本"),
		},
	}
	projectionRepo := &stubProjectionRepository{}

	svc := NewChangeRequestService(crRepo, projectionRepo, characterRepo)

	if err := svc.Process(context.Background(), scopeRequestID.Hex(), writer.CRStatusAccepted, "user-1"); err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	if projectionRepo.saved == nil {
		t.Fatal("expected projection to be saved")
	}
	if got := len(projectionRepo.saved.Characters); got != 2 {
		t.Fatalf("expected 2 projected characters, got %d", got)
	}
	if got := len(projectionRepo.saved.Relations); got != 1 {
		t.Fatalf("expected 1 projected relation, got %d", got)
	}
	if got := projectionRepo.saved.Relations[0].ToID; got != scopeRequestID.Hex() {
		t.Fatalf("expected relation target id %s, got %s", scopeRequestID.Hex(), got)
	}
	if got := projectionRepo.saved.Relations[0].ToName; got != "奶茶店老板" {
		t.Fatalf("expected relation target name %q, got %q", "奶茶店老板", got)
	}
}

func newTestCharacter(projectID, characterID primitive.ObjectID, name, currentState string) *writer.Character {
	character := &writer.Character{
		IdentifiedEntity:    writerBase.IdentifiedEntity{ID: characterID},
		ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
		CurrentState:        currentState,
	}
	character.Name = name
	return character
}

type stubChangeRequestRepository struct {
	requestByID      map[string]*writer.ChangeRequest
	requestsByStatus map[writer.ChangeRequestStatus][]*writer.ChangeRequest
	updatedStatus    writer.ChangeRequestStatus
	updatedBy        string
	pendingCount     int64
}

func (s *stubChangeRequestRepository) CreateRequest(ctx context.Context, cr *writer.ChangeRequest) error {
	if s.requestByID == nil {
		s.requestByID = make(map[string]*writer.ChangeRequest)
	}
	s.requestByID[cr.ID.Hex()] = cr
	return nil
}

func (s *stubChangeRequestRepository) FindRequestByID(ctx context.Context, id string) (*writer.ChangeRequest, error) {
	if s.requestByID == nil {
		return nil, nil
	}
	return s.requestByID[id], nil
}

func (s *stubChangeRequestRepository) FindRequestsByBatchID(ctx context.Context, batchID string) ([]*writer.ChangeRequest, error) {
	return nil, nil
}

func (s *stubChangeRequestRepository) FindPendingByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequest, error) {
	return nil, nil
}

func (s *stubChangeRequestRepository) FindByChapterAndStatus(ctx context.Context, projectID, chapterID string, status writer.ChangeRequestStatus) ([]*writer.ChangeRequest, error) {
	if s.requestsByStatus == nil {
		return nil, nil
	}
	return s.requestsByStatus[status], nil
}

func (s *stubChangeRequestRepository) CountPendingByChapter(ctx context.Context, projectID, chapterID string) (int64, error) {
	return s.pendingCount, nil
}

func (s *stubChangeRequestRepository) UpdateRequestStatus(ctx context.Context, id string, status writer.ChangeRequestStatus, processedBy string) error {
	s.updatedStatus = status
	s.updatedBy = processedBy
	if s.requestByID == nil {
		return nil
	}
	request := s.requestByID[id]
	if request == nil {
		return nil
	}
	request.Status = status
	if s.requestsByStatus == nil {
		s.requestsByStatus = make(map[writer.ChangeRequestStatus][]*writer.ChangeRequest)
	}
	for currentStatus, items := range s.requestsByStatus {
		filtered := items[:0]
		for _, item := range items {
			if item == nil || item.ID.Hex() == id {
				continue
			}
			filtered = append(filtered, item)
		}
		s.requestsByStatus[currentStatus] = filtered
	}
	s.requestsByStatus[status] = append(s.requestsByStatus[status], request)
	return nil
}

func (s *stubChangeRequestRepository) DeleteRequest(ctx context.Context, id string) error {
	return nil
}

func (s *stubChangeRequestRepository) CreateBatch(ctx context.Context, batch *writer.ChangeRequestBatch) error {
	return nil
}

func (s *stubChangeRequestRepository) FindBatchByID(ctx context.Context, id string) (*writer.ChangeRequestBatch, error) {
	return nil, nil
}

func (s *stubChangeRequestRepository) FindBatchesByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequestBatch, error) {
	return nil, nil
}

func (s *stubChangeRequestRepository) UpdateBatchCounts(ctx context.Context, id string, total, pending int) error {
	return nil
}

type stubProjectionRepository struct {
	projection *writer.ChapterProjection
	saved      *writer.ChapterProjection
}

func (s *stubProjectionRepository) GetByChapter(ctx context.Context, projectID, chapterID string) (*writer.ChapterProjection, error) {
	return s.projection, nil
}

func (s *stubProjectionRepository) UpsertByChapter(ctx context.Context, projection *writer.ChapterProjection) error {
	s.saved = projection
	s.projection = projection
	return nil
}

type stubCharacterRepository struct {
	characters []*writer.Character
	relations  []*writer.CharacterRelation
}

func (s *stubCharacterRepository) Create(ctx context.Context, character *writer.Character) error {
	return nil
}

func (s *stubCharacterRepository) FindByID(ctx context.Context, characterID string) (*writer.Character, error) {
	return nil, nil
}

func (s *stubCharacterRepository) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Character, error) {
	return s.characters, nil
}

func (s *stubCharacterRepository) Update(ctx context.Context, character *writer.Character) error {
	return nil
}

func (s *stubCharacterRepository) Delete(ctx context.Context, characterID string) error {
	return nil
}

func (s *stubCharacterRepository) CreateRelation(ctx context.Context, relation *writer.CharacterRelation) error {
	return nil
}

func (s *stubCharacterRepository) FindRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	return s.relations, nil
}

func (s *stubCharacterRepository) FindRelationByID(ctx context.Context, relationID string) (*writer.CharacterRelation, error) {
	return nil, nil
}

func (s *stubCharacterRepository) DeleteRelation(ctx context.Context, relationID string) error {
	return nil
}

func (s *stubCharacterRepository) CreateRelationTimelineEvent(ctx context.Context, relationID string, event *writer.RelationTimelineEvent) error {
	return nil
}

func (s *stubCharacterRepository) GetRelationTimeline(ctx context.Context, relationID string) ([]writer.RelationTimelineEvent, error) {
	return nil, nil
}

func (s *stubCharacterRepository) UpdateRelationTimelineEvent(ctx context.Context, relationID string, eventIndex int, event *writer.RelationTimelineEvent) error {
	return nil
}

func (s *stubCharacterRepository) DeleteRelationTimelineEvent(ctx context.Context, relationID string, eventIndex int) error {
	return nil
}

func (s *stubCharacterRepository) ExistsByID(ctx context.Context, characterID string) (bool, error) {
	return false, nil
}

func (s *stubCharacterRepository) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	return int64(len(s.characters)), nil
}
