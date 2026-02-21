package writer

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
	writeriface "Qingyu_backend/service/interfaces/writer"
	writermigration "Qingyu_backend/service/writer/_migration"
)

// 类型别名，方便测试代码使用
type WriterServiceAdapter = writermigration.WriterServiceAdapter

// Mock 实现
type mockFactoryProjectPort struct{}

func (m *mockFactoryProjectPort) Initialize(ctx context.Context) error { return nil }
func (m *mockFactoryProjectPort) Health(ctx context.Context) error     { return nil }
func (m *mockFactoryProjectPort) Close(ctx context.Context) error      { return nil }
func (m *mockFactoryProjectPort) GetServiceName() string               { return "mock" }
func (m *mockFactoryProjectPort) GetVersion() string                   { return "1.0.0" }
func (m *mockFactoryProjectPort) CreateProject(ctx context.Context, req *writeriface.CreateProjectRequest) (*writeriface.CreateProjectResponse, error) {
	return &writeriface.CreateProjectResponse{}, nil
}
func (m *mockFactoryProjectPort) GetProject(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}
func (m *mockFactoryProjectPort) GetProjectByID(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}
func (m *mockFactoryProjectPort) GetByIDWithoutAuth(ctx context.Context, projectID string) (*writer.Project, error) {
	return &writer.Project{}, nil
}
func (m *mockFactoryProjectPort) ListMyProjects(ctx context.Context, req *writeriface.ListProjectsRequest) (*writeriface.ListProjectsResponse, error) {
	return &writeriface.ListProjectsResponse{}, nil
}
func (m *mockFactoryProjectPort) GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*writer.Project, error) {
	return []*writer.Project{}, nil
}
func (m *mockFactoryProjectPort) UpdateProject(ctx context.Context, projectID string, req *writeriface.UpdateProjectRequest) error {
	return nil
}
func (m *mockFactoryProjectPort) UpdateProjectByID(ctx context.Context, projectID, userID string, req *writer.Project) error {
	return nil
}
func (m *mockFactoryProjectPort) DeleteProject(ctx context.Context, projectID string) error {
	return nil
}
func (m *mockFactoryProjectPort) DeleteProjectByID(ctx context.Context, projectID, userID string) error {
	return nil
}
func (m *mockFactoryProjectPort) RestoreProjectByID(ctx context.Context, projectID, userID string) error {
	return nil
}
func (m *mockFactoryProjectPort) DeleteHard(ctx context.Context, projectID string) error {
	return nil
}
func (m *mockFactoryProjectPort) UpdateProjectStatistics(ctx context.Context, projectID string, stats *writer.ProjectStats) error {
	return nil
}
func (m *mockFactoryProjectPort) RecalculateProjectStatistics(ctx context.Context, projectID string) error {
	return nil
}

type mockFactoryDocumentPort struct{}

func (m *mockFactoryDocumentPort) Initialize(ctx context.Context) error { return nil }
func (m *mockFactoryDocumentPort) Health(ctx context.Context) error     { return nil }
func (m *mockFactoryDocumentPort) Close(ctx context.Context) error      { return nil }
func (m *mockFactoryDocumentPort) GetServiceName() string               { return "mock" }
func (m *mockFactoryDocumentPort) GetVersion() string                   { return "1.0.0" }
func (m *mockFactoryDocumentPort) CreateDocument(ctx context.Context, req *writeriface.CreateDocumentRequest) (*writeriface.CreateDocumentResponse, error) {
	return &writeriface.CreateDocumentResponse{}, nil
}
func (m *mockFactoryDocumentPort) GetDocument(ctx context.Context, documentID string) (*writer.Document, error) {
	return &writer.Document{}, nil
}
func (m *mockFactoryDocumentPort) GetDocumentTree(ctx context.Context, projectID string) (*writeriface.DocumentTreeResponse, error) {
	return &writeriface.DocumentTreeResponse{}, nil
}
func (m *mockFactoryDocumentPort) ListDocuments(ctx context.Context, req *writeriface.ListDocumentsRequest) (*writeriface.ListDocumentsResponse, error) {
	return &writeriface.ListDocumentsResponse{}, nil
}
func (m *mockFactoryDocumentPort) UpdateDocument(ctx context.Context, documentID string, req *writeriface.UpdateDocumentRequest) error {
	return nil
}
func (m *mockFactoryDocumentPort) DeleteDocument(ctx context.Context, documentID string) error {
	return nil
}
func (m *mockFactoryDocumentPort) MoveDocument(ctx context.Context, req *writeriface.MoveDocumentRequest) error {
	return nil
}
func (m *mockFactoryDocumentPort) ReorderDocuments(ctx context.Context, req *writeriface.ReorderDocumentsRequest) error {
	return nil
}
func (m *mockFactoryDocumentPort) AutoSaveDocument(ctx context.Context, req *writeriface.AutoSaveRequest) (*writeriface.AutoSaveResponse, error) {
	return &writeriface.AutoSaveResponse{}, nil
}
func (m *mockFactoryDocumentPort) GetSaveStatus(ctx context.Context, documentID string) (*writeriface.SaveStatusResponse, error) {
	return &writeriface.SaveStatusResponse{}, nil
}
func (m *mockFactoryDocumentPort) GetDocumentContent(ctx context.Context, documentID string) (*writeriface.DocumentContentResponse, error) {
	return &writeriface.DocumentContentResponse{}, nil
}
func (m *mockFactoryDocumentPort) UpdateDocumentContent(ctx context.Context, req *writeriface.UpdateContentRequest) error {
	return nil
}
func (m *mockFactoryDocumentPort) DuplicateDocument(ctx context.Context, documentID string, req *writeriface.DuplicateRequest) (*writeriface.DuplicateResponse, error) {
	return &writeriface.DuplicateResponse{}, nil
}
func (m *mockFactoryDocumentPort) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}
func (m *mockFactoryDocumentPort) UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}
func (m *mockFactoryDocumentPort) RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}
func (m *mockFactoryDocumentPort) ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*writer.FileRevision, error) {
	return make([]*writer.FileRevision, 0), nil
}
func (m *mockFactoryDocumentPort) GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*writeriface.VersionHistoryResponse, error) {
	return &writeriface.VersionHistoryResponse{}, nil
}
func (m *mockFactoryDocumentPort) GetVersionDetail(ctx context.Context, documentID, versionID string) (*writeriface.VersionDetail, error) {
	return &writeriface.VersionDetail{}, nil
}
func (m *mockFactoryDocumentPort) CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*writeriface.VersionDiff, error) {
	return &writeriface.VersionDiff{}, nil
}
func (m *mockFactoryDocumentPort) RestoreVersion(ctx context.Context, documentID, versionID string) error {
	return nil
}
func (m *mockFactoryDocumentPort) CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*writer.FilePatch, error) {
	return &writer.FilePatch{}, nil
}
func (m *mockFactoryDocumentPort) ApplyPatch(projectID, patchID, applierID string) (*writer.FileRevision, error) {
	return &writer.FileRevision{}, nil
}
func (m *mockFactoryDocumentPort) DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*writer.ConflictInfo, error) {
	return &writer.ConflictInfo{}, nil
}
func (m *mockFactoryDocumentPort) CreateCommit(ctx context.Context, projectID, authorID, message string, files []writer.CommitFile) (*writer.Commit, error) {
	return &writer.Commit{}, nil
}
func (m *mockFactoryDocumentPort) ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*writer.Commit, error) {
	return []*writer.Commit{}, nil
}
func (m *mockFactoryDocumentPort) GetCommitDetails(ctx context.Context, projectID, commitID string) (*writer.Commit, []*writer.FileRevision, error) {
	return &writer.Commit{}, make([]*writer.FileRevision, 0), nil
}
func (m *mockFactoryDocumentPort) StartAutoSave(documentID, projectID, nodeID, userID string) error {
	return nil
}
func (m *mockFactoryDocumentPort) StopAutoSave(documentID string) error {
	return nil
}
func (m *mockFactoryDocumentPort) GetAutoSaveStatus(documentID string) (bool, *interface{}) {
	return true, nil
}
func (m *mockFactoryDocumentPort) SaveImmediately(documentID string) error {
	return nil
}

type mockFactoryCollaborationPort struct{}

func (m *mockFactoryCollaborationPort) Initialize(ctx context.Context) error { return nil }
func (m *mockFactoryCollaborationPort) Health(ctx context.Context) error     { return nil }
func (m *mockFactoryCollaborationPort) Close(ctx context.Context) error      { return nil }
func (m *mockFactoryCollaborationPort) GetServiceName() string               { return "mock" }
func (m *mockFactoryCollaborationPort) GetVersion() string                   { return "1.0.0" }
func (m *mockFactoryCollaborationPort) CreateComment(ctx context.Context, comment *writer.DocumentComment) (*writer.DocumentComment, error) {
	return &writer.DocumentComment{}, nil
}
func (m *mockFactoryCollaborationPort) GetComment(ctx context.Context, commentID string) (*writer.DocumentComment, error) {
	return &writer.DocumentComment{}, nil
}
func (m *mockFactoryCollaborationPort) UpdateComment(ctx context.Context, commentID string, comment *writer.DocumentComment) error {
	return nil
}
func (m *mockFactoryCollaborationPort) DeleteComment(ctx context.Context, commentID string) error {
	return nil
}
func (m *mockFactoryCollaborationPort) ListComments(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	return []*writer.DocumentComment{}, 0, nil
}
func (m *mockFactoryCollaborationPort) GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	return []*writer.DocumentComment{}, nil
}
func (m *mockFactoryCollaborationPort) GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	return []*writer.DocumentComment{}, nil
}
func (m *mockFactoryCollaborationPort) ResolveComment(ctx context.Context, commentID, userID string) error {
	return nil
}
func (m *mockFactoryCollaborationPort) UnresolveComment(ctx context.Context, commentID string) error {
	return nil
}
func (m *mockFactoryCollaborationPort) ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*writer.DocumentComment, error) {
	return &writer.DocumentComment{}, nil
}
func (m *mockFactoryCollaborationPort) GetCommentThread(ctx context.Context, threadID string) (*writer.CommentThread, error) {
	return &writer.CommentThread{}, nil
}
func (m *mockFactoryCollaborationPort) GetCommentReplies(ctx context.Context, parentID string) ([]*writer.DocumentComment, error) {
	return []*writer.DocumentComment{}, nil
}
func (m *mockFactoryCollaborationPort) GetCommentStats(ctx context.Context, documentID string) (*writer.CommentStats, error) {
	return &writer.CommentStats{}, nil
}
func (m *mockFactoryCollaborationPort) BatchDeleteComments(ctx context.Context, commentIDs []string) error {
	return nil
}
func (m *mockFactoryCollaborationPort) SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	return []*writer.DocumentComment{}, 0, nil
}

type mockFactoryContentPort struct{}

func (m *mockFactoryContentPort) Initialize(ctx context.Context) error { return nil }
func (m *mockFactoryContentPort) Health(ctx context.Context) error     { return nil }
func (m *mockFactoryContentPort) Close(ctx context.Context) error      { return nil }
func (m *mockFactoryContentPort) GetServiceName() string               { return "mock" }
func (m *mockFactoryContentPort) GetVersion() string                   { return "1.0.0" }
func (m *mockFactoryContentPort) CreateCharacter(ctx context.Context, projectID, userID string, req *writeriface.CreateCharacterRequest) (*writer.Character, error) {
	return &writer.Character{}, nil
}
func (m *mockFactoryContentPort) GetCharacterByID(ctx context.Context, characterID, projectID string) (*writer.Character, error) {
	return &writer.Character{}, nil
}
func (m *mockFactoryContentPort) ListCharacters(ctx context.Context, projectID string) ([]*writer.Character, error) {
	return []*writer.Character{}, nil
}
func (m *mockFactoryContentPort) UpdateCharacter(ctx context.Context, characterID, projectID string, req *writeriface.UpdateCharacterRequest) (*writer.Character, error) {
	return &writer.Character{}, nil
}
func (m *mockFactoryContentPort) DeleteCharacter(ctx context.Context, characterID, projectID string) error {
	return nil
}
func (m *mockFactoryContentPort) CreateCharacterRelation(ctx context.Context, projectID string, req *writeriface.CreateRelationRequest) (*writer.CharacterRelation, error) {
	return &writer.CharacterRelation{}, nil
}
func (m *mockFactoryContentPort) ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	return []*writer.CharacterRelation{}, nil
}
func (m *mockFactoryContentPort) DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error {
	return nil
}
func (m *mockFactoryContentPort) GetCharacterGraph(ctx context.Context, projectID string) (*writeriface.CharacterGraph, error) {
	return &writeriface.CharacterGraph{}, nil
}
func (m *mockFactoryContentPort) CreateLocation(ctx context.Context, projectID, userID string, req *writeriface.CreateLocationRequest) (*writer.Location, error) {
	return &writer.Location{}, nil
}
func (m *mockFactoryContentPort) GetLocationByID(ctx context.Context, locationID, projectID string) (*writer.Location, error) {
	return &writer.Location{}, nil
}
func (m *mockFactoryContentPort) ListLocations(ctx context.Context, projectID string) ([]*writer.Location, error) {
	return []*writer.Location{}, nil
}
func (m *mockFactoryContentPort) UpdateLocation(ctx context.Context, locationID, projectID string, req *writeriface.UpdateLocationRequest) (*writer.Location, error) {
	return &writer.Location{}, nil
}
func (m *mockFactoryContentPort) DeleteLocation(ctx context.Context, locationID, projectID string) error {
	return nil
}
func (m *mockFactoryContentPort) GetLocationTree(ctx context.Context, projectID string) ([]*writeriface.LocationNode, error) {
	return []*writeriface.LocationNode{}, nil
}
func (m *mockFactoryContentPort) GetLocationPath(ctx context.Context, locationID string) ([]string, error) {
	return []string{}, nil
}
func (m *mockFactoryContentPort) CreateLocationRelation(ctx context.Context, projectID string, req *writeriface.CreateLocationRelationRequest) (*writer.LocationRelation, error) {
	return &writer.LocationRelation{}, nil
}
func (m *mockFactoryContentPort) ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error) {
	return []*writer.LocationRelation{}, nil
}
func (m *mockFactoryContentPort) DeleteLocationRelation(ctx context.Context, relationID, projectID string) error {
	return nil
}
func (m *mockFactoryContentPort) CreateTimeline(ctx context.Context, projectID string, req *writeriface.CreateTimelineRequest) (*writer.Timeline, error) {
	return &writer.Timeline{}, nil
}
func (m *mockFactoryContentPort) GetTimeline(ctx context.Context, timelineID, projectID string) (*writer.Timeline, error) {
	return &writer.Timeline{}, nil
}
func (m *mockFactoryContentPort) ListTimelines(ctx context.Context, projectID string) ([]*writer.Timeline, error) {
	return []*writer.Timeline{}, nil
}
func (m *mockFactoryContentPort) DeleteTimeline(ctx context.Context, timelineID, projectID string) error {
	return nil
}
func (m *mockFactoryContentPort) CreateTimelineEvent(ctx context.Context, projectID string, req *writeriface.CreateTimelineEventRequest) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}
func (m *mockFactoryContentPort) GetTimelineEvent(ctx context.Context, eventID, projectID string) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}
func (m *mockFactoryContentPort) ListTimelineEvents(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error) {
	return []*writer.TimelineEvent{}, nil
}
func (m *mockFactoryContentPort) UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *writeriface.UpdateTimelineEventRequest) (*writer.TimelineEvent, error) {
	return &writer.TimelineEvent{}, nil
}
func (m *mockFactoryContentPort) DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error {
	return nil
}
func (m *mockFactoryContentPort) GetTimelineVisualization(ctx context.Context, timelineID string) (*writeriface.TimelineVisualization, error) {
	return &writeriface.TimelineVisualization{}, nil
}

type mockFactoryPublishExportPort struct{}

func (m *mockFactoryPublishExportPort) Initialize(ctx context.Context) error { return nil }
func (m *mockFactoryPublishExportPort) Health(ctx context.Context) error     { return nil }
func (m *mockFactoryPublishExportPort) Close(ctx context.Context) error      { return nil }
func (m *mockFactoryPublishExportPort) GetServiceName() string               { return "mock" }
func (m *mockFactoryPublishExportPort) GetVersion() string                   { return "1.0.0" }
func (m *mockFactoryPublishExportPort) PublishProject(ctx context.Context, projectID, userID string, req *writeriface.PublishProjectRequest) (*writeriface.PublicationRecord, error) {
	return &writeriface.PublicationRecord{}, nil
}
func (m *mockFactoryPublishExportPort) UnpublishProject(ctx context.Context, projectID, userID string) error {
	return nil
}
func (m *mockFactoryPublishExportPort) GetProjectPublicationStatus(ctx context.Context, projectID string) (*writeriface.PublicationStatus, error) {
	return &writeriface.PublicationStatus{}, nil
}
func (m *mockFactoryPublishExportPort) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.PublishDocumentRequest) (*writeriface.PublicationRecord, error) {
	return &writeriface.PublicationRecord{}, nil
}
func (m *mockFactoryPublishExportPort) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *writeriface.UpdateDocumentPublishStatusRequest) error {
	return nil
}
func (m *mockFactoryPublishExportPort) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *writeriface.BatchPublishDocumentsRequest) (*writeriface.BatchPublishResult, error) {
	return &writeriface.BatchPublishResult{}, nil
}
func (m *mockFactoryPublishExportPort) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.PublicationRecord, int64, error) {
	return []*writeriface.PublicationRecord{}, 0, nil
}
func (m *mockFactoryPublishExportPort) GetPublicationRecord(ctx context.Context, recordID string) (*writeriface.PublicationRecord, error) {
	return &writeriface.PublicationRecord{}, nil
}
func (m *mockFactoryPublishExportPort) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.ExportDocumentRequest) (*writeriface.ExportTask, error) {
	return &writeriface.ExportTask{}, nil
}
func (m *mockFactoryPublishExportPort) ExportProject(ctx context.Context, projectID, userID string, req *writeriface.ExportProjectRequest) (*writeriface.ExportTask, error) {
	return &writeriface.ExportTask{}, nil
}
func (m *mockFactoryPublishExportPort) GetExportTask(ctx context.Context, taskID string) (*writeriface.ExportTask, error) {
	return &writeriface.ExportTask{}, nil
}
func (m *mockFactoryPublishExportPort) DownloadExportFile(ctx context.Context, taskID string) (*writeriface.ExportFile, error) {
	return &writeriface.ExportFile{}, nil
}
func (m *mockFactoryPublishExportPort) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.ExportTask, int64, error) {
	return []*writeriface.ExportTask{}, 0, nil
}
func (m *mockFactoryPublishExportPort) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	return nil
}
func (m *mockFactoryPublishExportPort) CancelExportTask(ctx context.Context, taskID, userID string) error {
	return nil
}

// TestNewWriterServiceFactory 测试工厂创建
func TestNewWriterServiceFactory(t *testing.T) {
	factory := NewWriterServiceFactory()

	if factory == nil {
		t.Fatal("NewWriterServiceFactory returned nil")
	}
}

// TestWriterServiceFactory_CreateWithPorts 测试使用 Port 接口创建服务
func TestWriterServiceFactory_CreateWithPorts(t *testing.T) {
	factory := NewWriterServiceFactory()

	projectPort := &mockFactoryProjectPort{}
	documentPort := &mockFactoryDocumentPort{}
	collaborationPort := &mockFactoryCollaborationPort{}
	contentPort := &mockFactoryContentPort{}
	publishExportPort := &mockFactoryPublishExportPort{}

	service := factory.CreateWithPorts(
		projectPort,
		documentPort,
		collaborationPort,
		contentPort,
		publishExportPort,
	)

	if service == nil {
		t.Fatal("CreateWithPorts returned nil")
	}

	// 验证返回的是适配器类型
	_, ok := service.(*WriterServiceAdapter)
	if !ok {
		t.Error("CreateWithPorts did not return *WriterServiceAdapter")
	}
}

// TestWriterServiceFactory_CreateFromImplementations 测试从结构体创建服务
func TestWriterServiceFactory_CreateFromImplementations(t *testing.T) {
	factory := NewWriterServiceFactory()

	ports := PortImplementations{
		ProjectPort:       &mockFactoryProjectPort{},
		DocumentPort:      &mockFactoryDocumentPort{},
		CollaborationPort: &mockFactoryCollaborationPort{},
		ContentPort:       &mockFactoryContentPort{},
		PublishExportPort: &mockFactoryPublishExportPort{},
	}

	service := factory.CreateFromImplementations(ports)

	if service == nil {
		t.Fatal("CreateFromImplementations returned nil")
	}

	// 验证返回的是适配器类型
	_, ok := service.(*WriterServiceAdapter)
	if !ok {
		t.Error("CreateFromImplementations did not return *WriterServiceAdapter")
	}
}
