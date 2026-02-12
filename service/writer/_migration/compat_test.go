package migration

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
	writeriface "Qingyu_backend/service/interfaces/writer"
)

// 创建 Mock 实现
type mockProjectPort struct{}

func (m *mockProjectPort) CreateProject(ctx context.Context, req *writeriface.CreateProjectRequest) (*writeriface.CreateProjectResponse, error) {
	return &writeriface.CreateProjectResponse{ProjectID: "test-id"}, nil
}

func (m *mockProjectPort) GetProject(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}

func (m *mockProjectPort) GetProjectByID(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}

func (m *mockProjectPort) GetByIDWithoutAuth(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}

func (m *mockProjectPort) ListMyProjects(ctx context.Context, req *writeriface.ListProjectsRequest) (*writeriface.ListProjectsResponse, error) {
	return &writeriface.ListProjectsResponse{}, nil
}

func (m *mockProjectPort) GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*writer.Project, error) {
	return []*writer.Project{}, nil
}

func (m *mockProjectPort) UpdateProject(ctx context.Context, projectID string, req *writeriface.UpdateProjectRequest) error {
	return nil
}

func (m *mockProjectPort) UpdateProjectByID(ctx context.Context, projectID, userID string, req *writer.Project) error {
	return nil
}

func (m *mockProjectPort) DeleteProject(ctx context.Context, projectID string) error {
	return nil
}

func (m *mockProjectPort) DeleteProjectByID(ctx context.Context, projectID, userID string) error {
	return nil
}

func (m *mockProjectPort) RestoreProjectByID(ctx context.Context, projectID, userID string) error {
	return nil
}

func (m *mockProjectPort) DeleteHard(ctx context.Context, projectID string) error {
	return nil
}

func (m *mockProjectPort) UpdateProjectStatistics(ctx context.Context, projectID string, stats *writer.ProjectStats) error {
	return nil
}

func (m *mockProjectPort) RecalculateProjectStatistics(ctx context.Context, projectID string) error {
	return nil
}

func (m *mockProjectPort) Initialize(ctx context.Context) error { return nil }
func (m *mockProjectPort) Health(ctx context.Context) error     { return nil }
func (m *mockProjectPort) Close(ctx context.Context) error       { return nil }
func (m *mockProjectPort) GetServiceName() string               { return "mock-project" }
func (m *mockProjectPort) GetVersion() string                  { return "1.0.0" }

// 其他 Port 的简化 Mock
type mockDocumentPort struct{}

func (m *mockDocumentPort) Initialize(ctx context.Context) error { return nil }
func (m *mockDocumentPort) Health(ctx context.Context) error     { return nil }
func (m *mockDocumentPort) Close(ctx context.Context) error       { return nil }
func (m *mockDocumentPort) GetServiceName() string               { return "mock-document" }
func (m *mockDocumentPort) GetVersion() string                  { return "1.0.0" }

func (m *mockDocumentPort) CreateDocument(ctx context.Context, req *writeriface.CreateDocumentRequest) (*writeriface.CreateDocumentResponse, error) {
	return &writeriface.CreateDocumentResponse{}, nil
}

func (m *mockDocumentPort) GetDocument(ctx context.Context, documentID string) (*writer.Document, error) {
	return &writer.Document{}, nil
}

func (m *mockDocumentPort) GetDocumentTree(ctx context.Context, projectID string) (*writeriface.DocumentTreeResponse, error) {
	return &writeriface.DocumentTreeResponse{}, nil
}

func (m *mockDocumentPort) ListDocuments(ctx context.Context, req *writeriface.ListDocumentsRequest) (*writeriface.ListDocumentsResponse, error) {
	return &writeriface.ListDocumentsResponse{}, nil
}

func (m *mockDocumentPort) UpdateDocument(ctx context.Context, documentID string, req *writeriface.UpdateDocumentRequest) error {
	return nil
}

func (m *mockDocumentPort) DeleteDocument(ctx context.Context, documentID string) error {
	return nil
}

func (m *mockDocumentPort) MoveDocument(ctx context.Context, req *writeriface.MoveDocumentRequest) error {
	return nil
}

func (m *mockDocumentPort) ReorderDocuments(ctx context.Context, req *writeriface.ReorderDocumentsRequest) error {
	return nil
}

func (m *mockDocumentPort) AutoSaveDocument(ctx context.Context, req *writeriface.AutoSaveRequest) (*writeriface.AutoSaveResponse, error) {
	return &writeriface.AutoSaveResponse{}, nil
}

func (m *mockDocumentPort) GetSaveStatus(ctx context.Context, documentID string) (*writeriface.SaveStatusResponse, error) {
	return &writeriface.SaveStatusResponse{}, nil
}

func (m *mockDocumentPort) GetDocumentContent(ctx context.Context, documentID string) (*writeriface.DocumentContentResponse, error) {
	return &writeriface.DocumentContentResponse{}, nil
}

func (m *mockDocumentPort) UpdateDocumentContent(ctx context.Context, req *writeriface.UpdateContentRequest) error {
	return nil
}

func (m *mockDocumentPort) DuplicateDocument(ctx context.Context, documentID string, req *writeriface.DuplicateRequest) (*writeriface.DuplicateResponse, error) {
	return &writeriface.DuplicateResponse{}, nil
}

func (m *mockDocumentPort) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentPort) UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentPort) RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentPort) ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*writer.FileRevision, error) {
	return []*writer.FileRevision{}, nil
}

func (m *mockDocumentPort) GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*writeriface.VersionHistoryResponse, error) {
	return &writeriface.VersionHistoryResponse{}, nil
}

func (m *mockDocumentPort) GetVersionDetail(ctx context.Context, documentID, versionID string) (*writeriface.VersionDetail, error) {
	return &writeriface.VersionDetail{}, nil
}

func (m *mockDocumentPort) CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*writeriface.VersionDiff, error) {
	return &writeriface.VersionDiff{}, nil
}

func (m *mockDocumentPort) RestoreVersion(ctx context.Context, documentID, versionID string) error {
	return nil
}

func (m *mockDocumentPort) CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*writer.FilePatch, error) {
	return &writer.FilePatch{}, nil
}

func (m *mockDocumentPort) ApplyPatch(projectID, patchID, applierID string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}

func (m *mockDocumentPort) DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*writer.ConflictInfo, error) {
	return &writer.ConflictInfo{}, nil
}

func (m *mockDocumentPort) CreateCommit(ctx context.Context, projectID, authorID, message string, files []writer.CommitFile) (*writer.Commit, error) {
	return &writer.Commit{}, nil
}

func (m *mockDocumentPort) ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*writer.Commit, error) {
	return []*writer.Commit{}, nil
}

func (m *mockDocumentPort) GetCommitDetails(ctx context.Context, projectID, commitID string) (*writer.Commit, []*writer.FileRevision, error) {
	return &writer.Commit{}, []*writer.FileRevision{}, nil
}

func (m *mockDocumentPort) StartAutoSave(documentID, projectID, nodeID, userID string) error {
	return nil
}

func (m *mockDocumentPort) StopAutoSave(documentID string) error {
	return nil
}

func (m *mockDocumentPort) GetAutoSaveStatus(documentID string) (bool, *interface{}) {
	return true, nil
}

func (m *mockDocumentPort) SaveImmediately(documentID string) error {
	return nil
}

type mockCollaborationPort struct{}

func (m *mockCollaborationPort) Initialize(ctx context.Context) error { return nil }
func (m *mockCollaborationPort) Health(ctx context.Context) error     { return nil }
func (m *mockCollaborationPort) Close(ctx context.Context) error       { return nil }
func (m *mockCollaborationPort) GetServiceName() string               { return "mock-collab" }
func (m *mockCollaborationPort) GetVersion() string                  { return "1.0.0" }

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

type mockContentPort struct{}

func (m *mockContentPort) Initialize(ctx context.Context) error { return nil }
func (m *mockContentPort) Health(ctx context.Context) error     { return nil }
func (m *mockContentPort) Close(ctx context.Context) error       { return nil }
func (m *mockContentPort) GetServiceName() string               { return "mock-content" }
func (m *mockContentPort) GetVersion() string                  { return "1.0.0" }

func (m *mockContentPort) CreateCharacter(ctx context.Context, projectID, userID string, req *writeriface.CreateCharacterRequest) (*writer.Character, error) {
	return &writer.Character{}, nil
}

func (m *mockContentPort) GetCharacterByID(ctx context.Context, characterID, projectID string) (*writer.Character, error) {
	return &writer.Character{}, nil
}

func (m *mockContentPort) ListCharacters(ctx context.Context, projectID string) ([]*writer.Character, error) {
	return []*writer.Character{}, nil
}

func (m *mockContentPort) UpdateCharacter(ctx context.Context, characterID, projectID string, req *writeriface.UpdateCharacterRequest) (*writer.Character, error) {
	return &writer.Character{}, nil
}

func (m *mockContentPort) DeleteCharacter(ctx context.Context, characterID, projectID string) error {
	return nil
}

func (m *mockContentPort) CreateCharacterRelation(ctx context.Context, projectID string, req *writeriface.CreateRelationRequest) (*writer.CharacterRelation, error) {
	return &writer.CharacterRelation{}, nil
}

func (m *mockContentPort) ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	return []*writer.CharacterRelation{}, nil
}

func (m *mockContentPort) DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error {
	return nil
}

func (m *mockContentPort) GetCharacterGraph(ctx context.Context, projectID string) (*writeriface.CharacterGraph, error) {
	return &writeriface.CharacterGraph{}, nil
}

func (m *mockContentPort) CreateLocation(ctx context.Context, projectID, userID string, req *writeriface.CreateLocationRequest) (*writer.Location, error) {
	return &writer.Location{}, nil
}

func (m *mockContentPort) GetLocationByID(ctx context.Context, locationID, projectID string) (*writer.Location, error) {
	return &writer.Location{}, nil
}

func (m *mockContentPort) ListLocations(ctx context.Context, projectID string) ([]*writer.Location, error) {
	return []*writer.Location{}, nil
}

func (m *mockContentPort) UpdateLocation(ctx context.Context, locationID, projectID string, req *writeriface.UpdateLocationRequest) (*writer.Location, error) {
	return &writer.Location{}, nil
}

func (m *mockContentPort) DeleteLocation(ctx context.Context, locationID, projectID string) error {
	return nil
}

func (m *mockContentPort) GetLocationTree(ctx context.Context, projectID string) ([]*writeriface.LocationNode, error) {
	return []*writeriface.LocationNode{}, nil
}

func (m *mockContentPort) GetLocationPath(ctx context.Context, locationID string) ([]string, error) {
	return []string{}, nil
}

func (m *mockContentPort) CreateLocationRelation(ctx context.Context, projectID string, req *writeriface.CreateLocationRelationRequest) (*writer.LocationRelation, error) {
	return &writer.LocationRelation{}, nil
}

func (m *mockContentPort) ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error) {
	return []*writer.LocationRelation{}, nil
}

func (m *mockContentPort) DeleteLocationRelation(ctx context.Context, relationID, projectID string) error {
	return nil
}

func (m *mockContentPort) CreateTimeline(ctx context.Context, projectID string, req *writeriface.CreateTimelineRequest) (*writer.Timeline, error) {
	return &writer.Timeline{}, nil
}

func (m *mockContentPort) GetTimeline(ctx context.Context, timelineID, projectID string) (*writer.Timeline, error) {
	return &writer.Timeline{}, nil
}

func (m *mockContentPort) ListTimelines(ctx context.Context, projectID string) ([]*writer.Timeline, error) {
	return []*writer.Timeline{}, nil
}

func (m *mockContentPort) DeleteTimeline(ctx context.Context, timelineID, projectID string) error {
	return nil
}

func (m *mockContentPort) CreateTimelineEvent(ctx context.Context, projectID string, req *writeriface.CreateTimelineEventRequest) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}

func (m *mockContentPort) GetTimelineEvent(ctx context.Context, eventID, projectID string) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}

func (m *mockContentPort) ListTimelineEvents(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error) {
	return []*writer.TimelineEvent{}, nil
}

func (m *mockContentPort) UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *writeriface.UpdateTimelineEventRequest) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}

func (m *mockContentPort) DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error {
	return nil
}

func (m *mockContentPort) GetTimelineVisualization(ctx context.Context, timelineID string) (*writeriface.TimelineVisualization, error) {
	return &writeriface.TimelineVisualization{}, nil
}

type mockPublishExportPort struct{}

func (m *mockPublishExportPort) Initialize(ctx context.Context) error { return nil }
func (m *mockPublishExportPort) Health(ctx context.Context) error     { return nil }
func (m *mockPublishExportPort) Close(ctx context.Context) error       { return nil }
func (m *mockPublishExportPort) GetServiceName() string               { return "mock-publish" }
func (m *mockPublishExportPort) GetVersion() string                  { return "1.0.0" }

func (m *mockPublishExportPort) PublishProject(ctx context.Context, projectID, userID string, req *writeriface.PublishProjectRequest) (*writeriface.PublicationRecord, error) {
	return &writeriface.PublicationRecord{}, nil
}

func (m *mockPublishExportPort) UnpublishProject(ctx context.Context, projectID, userID string) error {
	return nil
}

func (m *mockPublishExportPort) GetProjectPublicationStatus(ctx context.Context, projectID string) (*writeriface.PublicationStatus, error) {
	return &writeriface.PublicationStatus{}, nil
}

func (m *mockPublishExportPort) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.PublishDocumentRequest) (*writeriface.PublicationRecord, error) {
	return &writeriface.PublicationRecord{}, nil
}

func (m *mockPublishExportPort) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *writeriface.UpdateDocumentPublishStatusRequest) error {
	return nil
}

func (m *mockPublishExportPort) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *writeriface.BatchPublishDocumentsRequest) (*writeriface.BatchPublishResult, error) {
	return &writeriface.BatchPublishResult{}, nil
}

func (m *mockPublishExportPort) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.PublicationRecord, int64, error) {
	return []*writeriface.PublicationRecord{}, 0, nil
}

func (m *mockPublishExportPort) GetPublicationRecord(ctx context.Context, recordID string) (*writeriface.PublicationRecord, error) {
	return &writeriface.PublicationRecord{}, nil
}

func (m *mockPublishExportPort) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.ExportDocumentRequest) (*writeriface.ExportTask, error) {
	return &writeriface.ExportTask{}, nil
}

func (m *mockPublishExportPort) ExportProject(ctx context.Context, projectID, userID string, req *writeriface.ExportProjectRequest) (*writeriface.ExportTask, error) {
	return &writeriface.ExportTask{}, nil
}

func (m *mockPublishExportPort) GetExportTask(ctx context.Context, taskID string) (*writeriface.ExportTask, error) {
	return &writeriface.ExportTask{}, nil
}

func (m *mockPublishExportPort) DownloadExportFile(ctx context.Context, taskID string) (*writeriface.ExportFile, error) {
	return &writeriface.ExportFile{}, nil
}

func (m *mockPublishExportPort) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.ExportTask, int64, error) {
	return []*writeriface.ExportTask{}, 0, nil
}

func (m *mockPublishExportPort) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	return nil
}

func (m *mockPublishExportPort) CancelExportTask(ctx context.Context, taskID, userID string) error {
	return nil
}

// TestNewWriterServiceAdapter 测试适配器创建
func TestNewWriterServiceAdapter(t *testing.T) {
	projectPort := &mockProjectPort{}
	documentPort := &mockDocumentPort{}
	collaborationPort := &mockCollaborationPort{}
	contentPort := &mockContentPort{}
	publishExportPort := &mockPublishExportPort{}

	adapter := NewWriterServiceAdapter(
		projectPort,
		documentPort,
		collaborationPort,
		contentPort,
		publishExportPort,
	)

	if adapter == nil {
		t.Fatal("NewWriterServiceAdapter returned nil")
	}

	if adapter.projectPort != projectPort {
		t.Error("projectPort not set correctly")
	}

	if adapter.documentPort != documentPort {
		t.Error("documentPort not set correctly")
	}
}

// TestWriterServiceAdapter_BaseService 测试 BaseService 接口实现
func TestWriterServiceAdapter_BaseService(t *testing.T) {
	adapter := NewWriterServiceAdapter(
		&mockProjectPort{},
		&mockDocumentPort{},
		&mockCollaborationPort{},
		&mockContentPort{},
		&mockPublishExportPort{},
	)

	ctx := context.Background()

	// 测试 Initialize
	if err := adapter.Initialize(ctx); err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	// 测试 Health
	if err := adapter.Health(ctx); err != nil {
		t.Errorf("Health failed: %v", err)
	}

	// 测试 GetServiceName
	if name := adapter.GetServiceName(); name != "mock-project" {
		t.Errorf("GetServiceName = %v, want %v", name, "mock-project")
	}

	// 测试 GetVersion
	if version := adapter.GetVersion(); version != "1.0.0" {
		t.Errorf("GetVersion = %v, want %v", version, "1.0.0")
	}

	// 测试 Close
	if err := adapter.Close(ctx); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

// TestWriterServiceAdapter_ProjectManagement 测试项目管理方法委托
func TestWriterServiceAdapter_ProjectManagement(t *testing.T) {
	adapter := NewWriterServiceAdapter(
		&mockProjectPort{},
		&mockDocumentPort{},
		&mockCollaborationPort{},
		&mockContentPort{},
		&mockPublishExportPort{},
	)

	ctx := context.Background()

	// 测试 CreateProject
	resp, err := adapter.CreateProject(ctx, &writeriface.CreateProjectRequest{})
	if err != nil {
		t.Errorf("CreateProject failed: %v", err)
	}

	if resp.ProjectID != "test-id" {
		t.Errorf("CreateProject returned wrong ID: %v", resp.ProjectID)
	}

	// 测试 GetProject
	_, err = adapter.GetProject(ctx, "test-id")
	if err != nil {
		t.Errorf("GetProject failed: %v", err)
	}

	// 测试 ListMyProjects
	_, err = adapter.ListMyProjects(ctx, &writeriface.ListProjectsRequest{})
	if err != nil {
		t.Errorf("ListMyProjects failed: %v", err)
	}
}

// TestWriterServiceAdapter_DocumentManagement 测试文档管理方法委托
func TestWriterServiceAdapter_DocumentManagement(t *testing.T) {
	adapter := NewWriterServiceAdapter(
		&mockProjectPort{},
		&mockDocumentPort{},
		&mockCollaborationPort{},
		&mockContentPort{},
		&mockPublishExportPort{},
	)

	ctx := context.Background()

	// 测试 CreateDocument
	_, err := adapter.CreateDocument(ctx, &writeriface.CreateDocumentRequest{})
	if err != nil {
		t.Errorf("CreateDocument failed: %v", err)
	}

	// 测试 GetDocument
	_, err = adapter.GetDocument(ctx, "test-id")
	if err != nil {
		t.Errorf("GetDocument failed: %v", err)
	}

	// 测试 AutoSaveDocument
	_, err = adapter.AutoSaveDocument(ctx, &writeriface.AutoSaveRequest{})
	if err != nil {
		t.Errorf("AutoSaveDocument failed: %v", err)
	}
}

// TestWriterServiceAdapter_Collaboration 测试协作方法委托
func TestWriterServiceAdapter_Collaboration(t *testing.T) {
	adapter := NewWriterServiceAdapter(
		&mockProjectPort{},
		&mockDocumentPort{},
		&mockCollaborationPort{},
		&mockContentPort{},
		&mockPublishExportPort{},
	)

	ctx := context.Background()

	// 测试 CreateComment
	_, err := adapter.CreateComment(ctx, &writer.DocumentComment{})
	if err != nil {
		t.Errorf("CreateComment failed: %v", err)
	}

	// 测试 ResolveComment
	err = adapter.ResolveComment(ctx, "comment-id", "user-id")
	if err != nil {
		t.Errorf("ResolveComment failed: %v", err)
	}
}

// TestWriterServiceAdapter_ContentManagement 测试内容管理方法委托
func TestWriterServiceAdapter_ContentManagement(t *testing.T) {
	adapter := NewWriterServiceAdapter(
		&mockProjectPort{},
		&mockDocumentPort{},
		&mockCollaborationPort{},
		&mockContentPort{},
		&mockPublishExportPort{},
	)

	ctx := context.Background()

	// 测试 CreateCharacter
	_, err := adapter.CreateCharacter(ctx, "project-id", "user-id", &writeriface.CreateCharacterRequest{})
	if err != nil {
		t.Errorf("CreateCharacter failed: %v", err)
	}

	// 测试 CreateLocation
	_, err = adapter.CreateLocation(ctx, "project-id", "user-id", &writeriface.CreateLocationRequest{})
	if err != nil {
		t.Errorf("CreateLocation failed: %v", err)
	}

	// 测试 CreateTimeline
	_, err = adapter.CreateTimeline(ctx, "project-id", &writeriface.CreateTimelineRequest{})
	if err != nil {
		t.Errorf("CreateTimeline failed: %v", err)
	}
}

// TestWriterServiceAdapter_PublishExport 测试发布导出方法委托
func TestWriterServiceAdapter_PublishExport(t *testing.T) {
	adapter := NewWriterServiceAdapter(
		&mockProjectPort{},
		&mockDocumentPort{},
		&mockCollaborationPort{},
		&mockContentPort{},
		&mockPublishExportPort{},
	)

	ctx := context.Background()

	// 测试 PublishProject
	_, err := adapter.PublishProject(ctx, "project-id", "user-id", &writeriface.PublishProjectRequest{})
	if err != nil {
		t.Errorf("PublishProject failed: %v", err)
	}

	// 测试 ExportDocument
	_, err = adapter.ExportDocument(ctx, "doc-id", "project-id", "user-id", &writeriface.ExportDocumentRequest{})
	if err != nil {
		t.Errorf("ExportDocument failed: %v", err)
	}
}
