package writer

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/interfaces/base"
)

// Mock 实现 Port 接口用于测试

type mockProjectManagementPort struct {
	base.BaseService
}

func (m *mockProjectManagementPort) CreateProject(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	return &CreateProjectResponse{}, nil
}

func (m *mockProjectManagementPort) GetProject(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}

func (m *mockProjectManagementPort) GetProjectByID(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}

func (m *mockProjectManagementPort) GetByIDWithoutAuth(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}

func (m *mockProjectManagementPort) ListMyProjects(ctx context.Context, req *ListProjectsRequest) (*ListProjectsResponse, error) {
	return &ListProjectsResponse{}, nil
}

func (m *mockProjectManagementPort) GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*writer.Project, error) {
	return []*writer.Project{}, nil
}

func (m *mockProjectManagementPort) UpdateProject(ctx context.Context, projectID string, req *UpdateProjectRequest) error {
	return nil
}

func (m *mockProjectManagementPort) UpdateProjectByID(ctx context.Context, projectID, userID string, req *writer.Project) error {
	return nil
}

func (m *mockProjectManagementPort) DeleteProject(ctx context.Context, projectID string) error {
	return nil
}

func (m *mockProjectManagementPort) DeleteProjectByID(ctx context.Context, projectID, userID string) error {
	return nil
}

func (m *mockProjectManagementPort) RestoreProjectByID(ctx context.Context, projectID, userID string) error {
	return nil
}

func (m *mockProjectManagementPort) DeleteHard(ctx context.Context, projectID string) error {
	return nil
}

func (m *mockProjectManagementPort) UpdateProjectStatistics(ctx context.Context, projectID string, stats *writer.ProjectStats) error {
	return nil
}

func (m *mockProjectManagementPort) RecalculateProjectStatistics(ctx context.Context, projectID string) error {
	return nil
}

func (m *mockProjectManagementPort) Initialize(ctx context.Context) error { return nil }
func (m *mockProjectManagementPort) Health(ctx context.Context) error     { return nil }
func (m *mockProjectManagementPort) Close(ctx context.Context) error       { return nil }
func (m *mockProjectManagementPort) GetServiceName() string               { return "mock" }
func (m *mockProjectManagementPort) GetVersion() string                  { return "1.0.0" }

type mockDocumentManagementPort struct {
	base.BaseService
}

// DocumentManagementPort 所需方法（简化版）
func (m *mockDocumentManagementPort) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	return &CreateDocumentResponse{}, nil
}

func (m *mockDocumentManagementPort) GetDocument(ctx context.Context, documentID string) (*writer.Document, error) {
	return &writer.Document{}, nil
}

func (m *mockDocumentManagementPort) GetDocumentTree(ctx context.Context, projectID string) (*DocumentTreeResponse, error) {
	return &DocumentTreeResponse{}, nil
}

func (m *mockDocumentManagementPort) ListDocuments(ctx context.Context, req *ListDocumentsRequest) (*ListDocumentsResponse, error) {
	return &ListDocumentsResponse{}, nil
}

func (m *mockDocumentManagementPort) UpdateDocument(ctx context.Context, documentID string, req *UpdateDocumentRequest) error {
	return nil
}

func (m *mockDocumentManagementPort) DeleteDocument(ctx context.Context, documentID string) error {
	return nil
}

func (m *mockDocumentManagementPort) MoveDocument(ctx context.Context, req *MoveDocumentRequest) error {
	return nil
}

func (m *mockDocumentManagementPort) ReorderDocuments(ctx context.Context, req *ReorderDocumentsRequest) error {
	return nil
}

func (m *mockDocumentManagementPort) AutoSaveDocument(ctx context.Context, req *AutoSaveRequest) (*AutoSaveResponse, error) {
	return &AutoSaveResponse{}, nil
}

func (m *mockDocumentManagementPort) GetSaveStatus(ctx context.Context, documentID string) (*SaveStatusResponse, error) {
	return &SaveStatusResponse{}, nil
}

func (m *mockDocumentManagementPort) GetDocumentContent(ctx context.Context, documentID string) (*DocumentContentResponse, error) {
	return &DocumentContentResponse{}, nil
}

func (m *mockDocumentManagementPort) UpdateDocumentContent(ctx context.Context, req *UpdateContentRequest) error {
	return nil
}

func (m *mockDocumentManagementPort) DuplicateDocument(ctx context.Context, documentID string, req *DuplicateRequest) (*DuplicateResponse, error) {
	return &DuplicateResponse{}, nil
}

// 版本控制方法（简化版）
func (m *mockDocumentManagementPort) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentManagementPort) UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentManagementPort) RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentManagementPort) ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*writer.FileRevision, error) {
	return []*writer.FileRevision{}, nil
}

func (m *mockDocumentManagementPort) GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*VersionHistoryResponse, error) {
	return &VersionHistoryResponse{}, nil
}

func (m *mockDocumentManagementPort) GetVersionDetail(ctx context.Context, documentID, versionID string) (*VersionDetail, error) {
	return &VersionDetail{}, nil
}

func (m *mockDocumentManagementPort) CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*VersionDiff, error) {
	return &VersionDiff{}, nil
}

func (m *mockDocumentManagementPort) RestoreVersion(ctx context.Context, documentID, versionID string) error {
	return nil
}

func (m *mockDocumentManagementPort) CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*writer.FilePatch, error) {
	return &writer.FilePatch{}, nil
}

func (m *mockDocumentManagementPort) ApplyPatch(projectID, patchID, applierID string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentManagementPort) DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*writer.ConflictInfo, error) {
	return &writer.ConflictInfo{}, nil
}

func (m *mockDocumentManagementPort) CreateCommit(ctx context.Context, projectID, authorID, message string, files []writer.CommitFile) (*writer.Commit, error) {
	return &writer.Commit{}, nil
}

func (m *mockDocumentManagementPort) ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*writer.Commit, error) {
	return []*writer.Commit{}, nil
}

func (m *mockDocumentManagementPort) GetCommitDetails(ctx context.Context, projectID, commitID string) (*writer.Commit, []*writer.FileRevision, error) {
	return &writer.Commit{}, []*writer.FileRevision{}, nil
}

// 自动保存方法（简化版）
func (m *mockDocumentManagementPort) StartAutoSave(documentID, projectID, nodeID, userID string) error {
	return nil
}

func (m *mockDocumentManagementPort) StopAutoSave(documentID string) error {
	return nil
}

func (m *mockDocumentManagementPort) GetAutoSaveStatus(documentID string) (bool, *interface{}) {
	return true, nil
}

func (m *mockDocumentManagementPort) SaveImmediately(documentID string) error {
	return nil
}

func (m *mockDocumentManagementPort) Initialize(ctx context.Context) error { return nil }
func (m *mockDocumentManagementPort) Health(ctx context.Context) error     { return nil }
func (m *mockDocumentManagementPort) Close(ctx context.Context) error       { return nil }
func (m *mockDocumentManagementPort) GetServiceName() string               { return "mock" }
func (m *mockDocumentManagementPort) GetVersion() string                  { return "1.0.0" }

type mockCollaborationPort struct {
	base.BaseService
}

func (m *mockCollaborationPort) CreateComment(ctx context.Context, comment *writer.DocumentComment) (*writer.DocumentComment, error) {
	return &writer.DocumentComment{}, nil
}

func (m *mockCollaborationPort) GetComment(ctx context.Context, commentID string) (*writer.DocumentComment, error) {
	return &writer.DocumentComment{}, nil
}

func (m *mockCollaborationPort) UpdateComment(ctx context.Context, commentID string, comment *writer.DocumentComment) error {
	return nil
}

func (m *mockCollaborationPort) DeleteComment(ctx context.Context, commentID string) error {
	return nil
}

func (m *mockCollaborationPort) ListComments(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	return []*writer.DocumentComment{}, 0, nil
}

func (m *mockCollaborationPort) GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	return []*writer.DocumentComment{}, nil
}

func (m *mockCollaborationPort) GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	return []*writer.DocumentComment{}, nil
}

func (m *mockCollaborationPort) ResolveComment(ctx context.Context, commentID, userID string) error {
	return nil
}

func (m *mockCollaborationPort) UnresolveComment(ctx context.Context, commentID string) error {
	return nil
}

func (m *mockCollaborationPort) ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*writer.DocumentComment, error) {
	return &writer.DocumentComment{}, nil
}

func (m *mockCollaborationPort) GetCommentThread(ctx context.Context, threadID string) (*writer.CommentThread, error) {
	return &writer.CommentThread{}, nil
}

func (m *mockCollaborationPort) GetCommentReplies(ctx context.Context, parentID string) ([]*writer.DocumentComment, error) {
	return []*writer.DocumentComment{}, nil
}

func (m *mockCollaborationPort) GetCommentStats(ctx context.Context, documentID string) (*writer.CommentStats, error) {
	return &writer.CommentStats{}, nil
}

func (m *mockCollaborationPort) BatchDeleteComments(ctx context.Context, commentIDs []string) error {
	return nil
}

func (m *mockCollaborationPort) SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	return []*writer.DocumentComment{}, 0, nil
}

func (m *mockCollaborationPort) Initialize(ctx context.Context) error { return nil }
func (m *mockCollaborationPort) Health(ctx context.Context) error     { return nil }
func (m *mockCollaborationPort) Close(ctx context.Context) error       { return nil }
func (m *mockCollaborationPort) GetServiceName() string               { return "mock" }
func (m *mockCollaborationPort) GetVersion() string                  { return "1.0.0" }

type mockContentManagementPort struct {
	base.BaseService
}

// ContentManagementPort 所需方法（简化版，只实现角色管理部分）
func (m *mockContentManagementPort) CreateCharacter(ctx context.Context, projectID, userID string, req *CreateCharacterRequest) (*writer.Character, error) {
	return &writer.Character{}, nil
}

func (m *mockContentManagementPort) GetCharacterByID(ctx context.Context, characterID, projectID string) (*writer.Character, error) {
	return &writer.Character{}, nil
}

func (m *mockContentManagementPort) ListCharacters(ctx context.Context, projectID string) ([]*writer.Character, error) {
	return []*writer.Character{}, nil
}

func (m *mockContentManagementPort) UpdateCharacter(ctx context.Context, characterID, projectID string, req *UpdateCharacterRequest) (*writer.Character, error) {
	return &writer.Character{}, nil
}

func (m *mockContentManagementPort) DeleteCharacter(ctx context.Context, characterID, projectID string) error {
	return nil
}

func (m *mockContentManagementPort) CreateCharacterRelation(ctx context.Context, projectID string, req *CreateRelationRequest) (*writer.CharacterRelation, error) {
	return &writer.CharacterRelation{}, nil
}

func (m *mockContentManagementPort) ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	return []*writer.CharacterRelation{}, nil
}

func (m *mockContentManagementPort) DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error {
	return nil
}

func (m *mockContentManagementPort) GetCharacterGraph(ctx context.Context, projectID string) (*CharacterGraph, error) {
	return &CharacterGraph{}, nil
}

// 地点管理方法（简化版）
func (m *mockContentManagementPort) CreateLocation(ctx context.Context, projectID, userID string, req *CreateLocationRequest) (*writer.Location, error) {
	return &writer.Location{}, nil
}

func (m *mockContentManagementPort) GetLocationByID(ctx context.Context, locationID, projectID string) (*writer.Location, error) {
	return &writer.Location{}, nil
}

func (m *mockContentManagementPort) ListLocations(ctx context.Context, projectID string) ([]*writer.Location, error) {
	return []*writer.Location{}, nil
}

func (m *mockContentManagementPort) UpdateLocation(ctx context.Context, locationID, projectID string, req *UpdateLocationRequest) (*writer.Location, error) {
	return &writer.Location{}, nil
}

func (m *mockContentManagementPort) DeleteLocation(ctx context.Context, locationID, projectID string) error {
	return nil
}

func (m *mockContentManagementPort) GetLocationTree(ctx context.Context, projectID string) ([]*LocationNode, error) {
	return []*LocationNode{}, nil
}

func (m *mockContentManagementPort) GetLocationPath(ctx context.Context, locationID string) ([]string, error) {
	return []string{}, nil
}

func (m *mockContentManagementPort) CreateLocationRelation(ctx context.Context, projectID string, req *CreateLocationRelationRequest) (*writer.LocationRelation, error) {
	return &writer.LocationRelation{}, nil
}

func (m *mockContentManagementPort) ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error) {
	return []*writer.LocationRelation{}, nil
}

func (m *mockContentManagementPort) DeleteLocationRelation(ctx context.Context, relationID, projectID string) error {
	return nil
}

// 时间线管理方法（简化版）
func (m *mockContentManagementPort) CreateTimeline(ctx context.Context, projectID string, req *CreateTimelineRequest) (*writer.Timeline, error) {
	return &writer.Timeline{}, nil
}

func (m *mockContentManagementPort) GetTimeline(ctx context.Context, timelineID, projectID string) (*writer.Timeline, error) {
	return &writer.Timeline{}, nil
}

func (m *mockContentManagementPort) ListTimelines(ctx context.Context, projectID string) ([]*writer.Timeline, error) {
	return []*writer.Timeline{}, nil
}

func (m *mockContentManagementPort) DeleteTimeline(ctx context.Context, timelineID, projectID string) error {
	return nil
}

func (m *mockContentManagementPort) CreateTimelineEvent(ctx context.Context, projectID string, req *CreateTimelineEventRequest) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}

func (m *mockContentManagementPort) GetTimelineEvent(ctx context.Context, eventID, projectID string) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}

func (m *mockContentManagementPort) ListTimelineEvents(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error) {
	return []*writer.TimelineEvent{}, nil
}

func (m *mockContentManagementPort) UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *UpdateTimelineEventRequest) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}

func (m *mockContentManagementPort) DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error {
	return nil
}

func (m *mockContentManagementPort) GetTimelineVisualization(ctx context.Context, timelineID string) (*TimelineVisualization, error) {
	return &TimelineVisualization{}, nil
}

func (m *mockContentManagementPort) Initialize(ctx context.Context) error { return nil }
func (m *mockContentManagementPort) Health(ctx context.Context) error     { return nil }
func (m *mockContentManagementPort) Close(ctx context.Context) error       { return nil }
func (m *mockContentManagementPort) GetServiceName() string               { return "mock" }
func (m *mockContentManagementPort) GetVersion() string                  { return "1.0.0" }

type mockPublishExportPort struct {
	base.BaseService
}

// PublishExportPort 所需方法（简化版）
func (m *mockPublishExportPort) PublishProject(ctx context.Context, projectID, userID string, req *PublishProjectRequest) (*PublicationRecord, error) {
	return &PublicationRecord{}, nil
}

func (m *mockPublishExportPort) UnpublishProject(ctx context.Context, projectID, userID string) error {
	return nil
}

func (m *mockPublishExportPort) GetProjectPublicationStatus(ctx context.Context, projectID string) (*PublicationStatus, error) {
	return &PublicationStatus{}, nil
}

func (m *mockPublishExportPort) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *PublishDocumentRequest) (*PublicationRecord, error) {
	return &PublicationRecord{}, nil
}

func (m *mockPublishExportPort) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *UpdateDocumentPublishStatusRequest) error {
	return nil
}

func (m *mockPublishExportPort) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *BatchPublishDocumentsRequest) (*BatchPublishResult, error) {
	return &BatchPublishResult{}, nil
}

func (m *mockPublishExportPort) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*PublicationRecord, int64, error) {
	return []*PublicationRecord{}, 0, nil
}

func (m *mockPublishExportPort) GetPublicationRecord(ctx context.Context, recordID string) (*PublicationRecord, error) {
	return &PublicationRecord{}, nil
}

func (m *mockPublishExportPort) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *ExportDocumentRequest) (*ExportTask, error) {
	return &ExportTask{}, nil
}

func (m *mockPublishExportPort) ExportProject(ctx context.Context, projectID, userID string, req *ExportProjectRequest) (*ExportTask, error) {
	return &ExportTask{}, nil
}

func (m *mockPublishExportPort) GetExportTask(ctx context.Context, taskID string) (*ExportTask, error) {
	return &ExportTask{}, nil
}

func (m *mockPublishExportPort) DownloadExportFile(ctx context.Context, taskID string) (*ExportFile, error) {
	return &ExportFile{}, nil
}

func (m *mockPublishExportPort) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*ExportTask, int64, error) {
	return []*ExportTask{}, 0, nil
}

func (m *mockPublishExportPort) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	return nil
}

func (m *mockPublishExportPort) CancelExportTask(ctx context.Context, taskID, userID string) error {
	return nil
}

func (m *mockPublishExportPort) Initialize(ctx context.Context) error { return nil }
func (m *mockPublishExportPort) Health(ctx context.Context) error     { return nil }
func (m *mockPublishExportPort) Close(ctx context.Context) error       { return nil }
func (m *mockPublishExportPort) GetServiceName() string               { return "mock" }
func (m *mockPublishExportPort) GetVersion() string                  { return "1.0.0" }

// TestProjectManagementPort 测试 ProjectManagementPort 接口
func TestProjectManagementPort(t *testing.T) {
	var port ProjectManagementPort = &mockProjectManagementPort{}

	ctx := context.Background()

	// 测试 CreateProject
	_, err := port.CreateProject(ctx, &CreateProjectRequest{})
	if err != nil {
		t.Errorf("CreateProject failed: %v", err)
	}

	// 测试 GetProject
	_, err = port.GetProject(ctx, "test-id")
	if err != nil {
		t.Errorf("GetProject failed: %v", err)
	}

	// 测试 BaseService 方法
	err = port.Initialize(ctx)
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	err = port.Health(ctx)
	if err != nil {
		t.Errorf("Health failed: %v", err)
	}

	err = port.Close(ctx)
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if port.GetServiceName() != "mock" {
		t.Errorf("GetServiceName returned wrong name: %s", port.GetServiceName())
	}

	if port.GetVersion() != "1.0.0" {
		t.Errorf("GetVersion returned wrong version: %s", port.GetVersion())
	}
}

// TestDocumentManagementPort 测试 DocumentManagementPort 接口
func TestDocumentManagementPort(t *testing.T) {
	var port DocumentManagementPort = &mockDocumentManagementPort{}

	ctx := context.Background()

	// 测试 CreateDocument
	_, err := port.CreateDocument(ctx, &CreateDocumentRequest{})
	if err != nil {
		t.Errorf("CreateDocument failed: %v", err)
	}

	// 测试 GetDocument
	_, err = port.GetDocument(ctx, "test-id")
	if err != nil {
		t.Errorf("GetDocument failed: %v", err)
	}

	// 测试版本控制方法
	_, err = port.BumpVersionAndCreateRevision("project-id", "node-id", "user-id", "test message")
	if err != nil {
		t.Errorf("BumpVersionAndCreateRevision failed: %v", err)
	}
}

// TestCollaborationPort 测试 CollaborationPort 接口
func TestCollaborationPort(t *testing.T) {
	var port CollaborationPort = &mockCollaborationPort{}

	ctx := context.Background()

	// 测试 CreateComment
	comment := &writer.DocumentComment{}
	_, err := port.CreateComment(ctx, comment)
	if err != nil {
		t.Errorf("CreateComment failed: %v", err)
	}

	// 测试 GetComment
	_, err = port.GetComment(ctx, "test-id")
	if err != nil {
		t.Errorf("GetComment failed: %v", err)
	}
}

// TestContentManagementPort 测试 ContentManagementPort 接口
func TestContentManagementPort(t *testing.T) {
	var port ContentManagementPort = &mockContentManagementPort{}

	ctx := context.Background()

	// 测试 CreateCharacter
	_, err := port.CreateCharacter(ctx, "project-id", "user-id", &CreateCharacterRequest{})
	if err != nil {
		t.Errorf("CreateCharacter failed: %v", err)
	}

	// 测试 CreateLocation
	_, err = port.CreateLocation(ctx, "project-id", "user-id", &CreateLocationRequest{})
	if err != nil {
		t.Errorf("CreateLocation failed: %v", err)
	}

	// 测试 CreateTimeline
	_, err = port.CreateTimeline(ctx, "project-id", &CreateTimelineRequest{})
	if err != nil {
		t.Errorf("CreateTimeline failed: %v", err)
	}
}

// TestPublishExportPort 测试 PublishExportPort 接口
func TestPublishExportPort(t *testing.T) {
	var port PublishExportPort = &mockPublishExportPort{}

	ctx := context.Background()

	// 测试 PublishProject
	_, err := port.PublishProject(ctx, "project-id", "user-id", &PublishProjectRequest{})
	if err != nil {
		t.Errorf("PublishProject failed: %v", err)
	}

	// 测试 ExportDocument
	_, err = port.ExportDocument(ctx, "doc-id", "project-id", "user-id", &ExportDocumentRequest{})
	if err != nil {
		t.Errorf("ExportDocument failed: %v", err)
	}
}
