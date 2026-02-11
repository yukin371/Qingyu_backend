package migration

import (
	"context"

	writeriface "Qingyu_backend/service/interfaces/writer"
	"Qingyu_backend/models/writer"
)

// ============================================================================
// 兼容层 - 向后兼容支持
// ============================================================================

// WriterServiceAdapter 旧 WriterService 接口的适配器
// 将旧的 WriterService 方法调用委托给新的 Port 接口
type WriterServiceAdapter struct {
	projectPort      writeriface.ProjectManagementPort
	documentPort     writeriface.DocumentManagementPort
	collaborationPort writeriface.CollaborationPort
	contentPort      writeriface.ContentManagementPort
	publishExportPort writeriface.PublishExportPort
}

// ============================================================================
// BaseService 接口实现 - 委托给 projectPort
// ============================================================================

// Initialize 初始化服务
func (a *WriterServiceAdapter) Initialize(ctx context.Context) error {
	return a.projectPort.Initialize(ctx)
}

// Health 健康检查
func (a *WriterServiceAdapter) Health(ctx context.Context) error {
	return a.projectPort.Health(ctx)
}

// Close 关闭服务
func (a *WriterServiceAdapter) Close(ctx context.Context) error {
	// 关闭所有 Port
	if err := a.projectPort.Close(ctx); err != nil {
		return err
	}
	if err := a.documentPort.Close(ctx); err != nil {
		return err
	}
	if err := a.collaborationPort.Close(ctx); err != nil {
		return err
	}
	if err := a.contentPort.Close(ctx); err != nil {
		return err
	}
	if err := a.publishExportPort.Close(ctx); err != nil {
		return err
	}
	return nil
}

// GetServiceName 获取服务名称
func (a *WriterServiceAdapter) GetServiceName() string {
	return a.projectPort.GetServiceName()
}

// GetVersion 获取服务版本
func (a *WriterServiceAdapter) GetVersion() string {
	return a.projectPort.GetVersion()
}

// NewWriterServiceAdapter 创建新的适配器
func NewWriterServiceAdapter(
	projectPort writeriface.ProjectManagementPort,
	documentPort writeriface.DocumentManagementPort,
	collaborationPort writeriface.CollaborationPort,
	contentPort writeriface.ContentManagementPort,
	publishExportPort writeriface.PublishExportPort,
) *WriterServiceAdapter {
	return &WriterServiceAdapter{
		projectPort:      projectPort,
		documentPort:     documentPort,
		collaborationPort: collaborationPort,
		contentPort:      contentPort,
		publishExportPort: publishExportPort,
	}
}

// ============================================================================
// 辅助类型转换函数
// ============================================================================

// toInterfaceSlice 将切片转换为 interface{} 切片
func toInterfaceSlice(projects []*writer.Project) []*interface{} {
	result := make([]*interface{}, len(projects))
	for i, p := range projects {
		val := interface{}(p)
		result[i] = &val
	}
	return result
}

// fromInterfaceToProject 从 interface{} 转换为 Project
func fromInterfaceToProject(req *interface{}) *writer.Project {
	if req == nil {
		return nil
	}
	if p, ok := (*req).(*writer.Project); ok {
		return p
	}
	return nil
}

// fromInterfaceToProjectStats 从 interface{} 转换为 ProjectStats
func fromInterfaceToProjectStats(stats *interface{}) *writer.ProjectStats {
	if stats == nil {
		return nil
	}
	if s, ok := (*stats).(*writer.ProjectStats); ok {
		return s
	}
	return nil
}

// fromInterfaceToDocument 从 interface{} 转换为 Document
func fromInterfaceToDocument(doc *interface{}) *writer.Document {
	if doc == nil {
		return nil
	}
	if d, ok := (*doc).(*writer.Document); ok {
		return d
	}
	return nil
}

// fromInterfaceToDocumentComment 从 interface{} 转换为 DocumentComment
func fromInterfaceToDocumentComment(comment *interface{}) *writer.DocumentComment {
	if comment == nil {
		return nil
	}
	if c, ok := (*comment).(*writer.DocumentComment); ok {
		return c
	}
	return nil
}

// fromInterfaceToCommentFilter 从 interface{} 转换为 CommentFilter
func fromInterfaceToCommentFilter(filter *interface{}) *writer.CommentFilter {
	if filter == nil {
		return nil
	}
	if f, ok := (*filter).(*writer.CommentFilter); ok {
		return f
	}
	return nil
}

// toInterface 将具体类型转换为 *interface{}
func toInterface(v interface{}) *interface{} {
	if v == nil {
		return nil
	}
	result := v
	return &result
}

// toInterfaceSliceComment 将 DocumentComment 切片转换为 interface{} 切片
func toInterfaceSliceComment(comments []*writer.DocumentComment) []*interface{} {
	result := make([]*interface{}, len(comments))
	for i, c := range comments {
		result[i] = toInterface(c)
	}
	return result
}

// toInterfaceSliceCharacter 将 Character 切片转换为 interface{} 切片
func toInterfaceSliceCharacter(characters []*writer.Character) []*interface{} {
	result := make([]*interface{}, len(characters))
	for i, c := range characters {
		result[i] = toInterface(c)
	}
	return result
}

// toInterfaceSliceCharacterRelation 将 CharacterRelation 切片转换为 interface{} 切片
func toInterfaceSliceCharacterRelation(relations []*writer.CharacterRelation) []*interface{} {
	result := make([]*interface{}, len(relations))
	for i, r := range relations {
		result[i] = toInterface(r)
	}
	return result
}

// toInterfaceSliceLocation 将 Location 切片转换为 interface{} 切片
func toInterfaceSliceLocation(locations []*writer.Location) []*interface{} {
	result := make([]*interface{}, len(locations))
	for i, l := range locations {
		result[i] = toInterface(l)
	}
	return result
}

// toInterfaceSliceLocationRelation 将 LocationRelation 切片转换为 interface{} 切片
func toInterfaceSliceLocationRelation(relations []*writer.LocationRelation) []*interface{} {
	result := make([]*interface{}, len(relations))
	for i, r := range relations {
		result[i] = toInterface(r)
	}
	return result
}

// toInterfaceSliceTimeline 将 Timeline 切片转换为 interface{} 切片
func toInterfaceSliceTimeline(timelines []*writer.Timeline) []*interface{} {
	result := make([]*interface{}, len(timelines))
	for i, t := range timelines {
		result[i] = toInterface(t)
	}
	return result
}

// toInterfaceSliceTimelineEvent 将 TimelineEvent 切片转换为 interface{} 切片
func toInterfaceSliceTimelineEvent(events []*writer.TimelineEvent) []*interface{} {
	result := make([]*interface{}, len(events))
	for i, e := range events {
		result[i] = toInterface(e)
	}
	return result
}

// ============================================================================
// 项目管理方法 - 委托给 ProjectManagementPort
// ============================================================================

func (a *WriterServiceAdapter) CreateProject(ctx context.Context, req *writeriface.CreateProjectRequest) (*writeriface.CreateProjectResponse, error) {
	return a.projectPort.CreateProject(ctx, req)
}

func (a *WriterServiceAdapter) GetProject(ctx context.Context, projectID string) (*interface{}, error) {
	project, err := a.projectPort.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(project), nil
}

func (a *WriterServiceAdapter) GetProjectByID(ctx context.Context, projectID string) (*interface{}, error) {
	project, err := a.projectPort.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(project), nil
}

func (a *WriterServiceAdapter) GetByIDWithoutAuth(ctx context.Context, projectID string) (*interface{}, error) {
	project, err := a.projectPort.GetByIDWithoutAuth(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(project), nil
}

func (a *WriterServiceAdapter) ListMyProjects(ctx context.Context, req *writeriface.ListProjectsRequest) (*writeriface.ListProjectsResponse, error) {
	return a.projectPort.ListMyProjects(ctx, req)
}

func (a *WriterServiceAdapter) GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*interface{}, error) {
	projects, err := a.projectPort.GetProjectList(ctx, userID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	return toInterfaceSlice(projects), nil
}

func (a *WriterServiceAdapter) UpdateProject(ctx context.Context, projectID string, req *writeriface.UpdateProjectRequest) error {
	return a.projectPort.UpdateProject(ctx, projectID, req)
}

func (a *WriterServiceAdapter) UpdateProjectByID(ctx context.Context, projectID, userID string, req *interface{}) error {
	project := fromInterfaceToProject(req)
	return a.projectPort.UpdateProjectByID(ctx, projectID, userID, project)
}

func (a *WriterServiceAdapter) DeleteProject(ctx context.Context, projectID string) error {
	return a.projectPort.DeleteProject(ctx, projectID)
}

func (a *WriterServiceAdapter) DeleteProjectByID(ctx context.Context, projectID, userID string) error {
	return a.projectPort.DeleteProjectByID(ctx, projectID, userID)
}

func (a *WriterServiceAdapter) RestoreProjectByID(ctx context.Context, projectID, userID string) error {
	return a.projectPort.RestoreProjectByID(ctx, projectID, userID)
}

func (a *WriterServiceAdapter) DeleteHard(ctx context.Context, projectID string) error {
	return a.projectPort.DeleteHard(ctx, projectID)
}

func (a *WriterServiceAdapter) UpdateProjectStatistics(ctx context.Context, projectID string, stats *interface{}) error {
	projectStats := fromInterfaceToProjectStats(stats)
	return a.projectPort.UpdateProjectStatistics(ctx, projectID, projectStats)
}

func (a *WriterServiceAdapter) RecalculateProjectStatistics(ctx context.Context, projectID string) error {
	return a.projectPort.RecalculateProjectStatistics(ctx, projectID)
}

// ============================================================================
// 文档管理方法 - 委托给 DocumentManagementPort
// ============================================================================

func (a *WriterServiceAdapter) CreateDocument(ctx context.Context, req *writeriface.CreateDocumentRequest) (*writeriface.CreateDocumentResponse, error) {
	return a.documentPort.CreateDocument(ctx, req)
}

func (a *WriterServiceAdapter) GetDocument(ctx context.Context, documentID string) (*interface{}, error) {
	document, err := a.documentPort.GetDocument(ctx, documentID)
	if err != nil {
		return nil, err
	}
	return toInterface(document), nil
}

func (a *WriterServiceAdapter) GetDocumentTree(ctx context.Context, projectID string) (*writeriface.DocumentTreeResponse, error) {
	return a.documentPort.GetDocumentTree(ctx, projectID)
}

func (a *WriterServiceAdapter) ListDocuments(ctx context.Context, req *writeriface.ListDocumentsRequest) (*writeriface.ListDocumentsResponse, error) {
	return a.documentPort.ListDocuments(ctx, req)
}

func (a *WriterServiceAdapter) UpdateDocument(ctx context.Context, documentID string, req *writeriface.UpdateDocumentRequest) error {
	return a.documentPort.UpdateDocument(ctx, documentID, req)
}

func (a *WriterServiceAdapter) DeleteDocument(ctx context.Context, documentID string) error {
	return a.documentPort.DeleteDocument(ctx, documentID)
}

func (a *WriterServiceAdapter) MoveDocument(ctx context.Context, req *writeriface.MoveDocumentRequest) error {
	return a.documentPort.MoveDocument(ctx, req)
}

func (a *WriterServiceAdapter) ReorderDocuments(ctx context.Context, req *writeriface.ReorderDocumentsRequest) error {
	return a.documentPort.ReorderDocuments(ctx, req)
}

func (a *WriterServiceAdapter) AutoSaveDocument(ctx context.Context, req *writeriface.AutoSaveRequest) (*writeriface.AutoSaveResponse, error) {
	return a.documentPort.AutoSaveDocument(ctx, req)
}

func (a *WriterServiceAdapter) GetSaveStatus(ctx context.Context, documentID string) (*writeriface.SaveStatusResponse, error) {
	return a.documentPort.GetSaveStatus(ctx, documentID)
}

func (a *WriterServiceAdapter) GetDocumentContent(ctx context.Context, documentID string) (*writeriface.DocumentContentResponse, error) {
	return a.documentPort.GetDocumentContent(ctx, documentID)
}

func (a *WriterServiceAdapter) UpdateDocumentContent(ctx context.Context, req *writeriface.UpdateContentRequest) error {
	return a.documentPort.UpdateDocumentContent(ctx, req)
}

func (a *WriterServiceAdapter) DuplicateDocument(ctx context.Context, documentID string, req *writeriface.DuplicateRequest) (*writeriface.DuplicateResponse, error) {
	return a.documentPort.DuplicateDocument(ctx, documentID, req)
}

// ============================================================================
// 协作批注方法 - 委托给 CollaborationPort
// ============================================================================

func (a *WriterServiceAdapter) CreateComment(ctx context.Context, comment *interface{}) (*interface{}, error) {
	c := fromInterfaceToDocumentComment(comment)
	result, err := a.collaborationPort.CreateComment(ctx, c)
	if err != nil {
		return nil, err
	}
	return toInterface(result), nil
}

func (a *WriterServiceAdapter) GetComment(ctx context.Context, commentID string) (*interface{}, error) {
	comment, err := a.collaborationPort.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}
	return toInterface(comment), nil
}

func (a *WriterServiceAdapter) UpdateComment(ctx context.Context, commentID string, comment *interface{}) error {
	c := fromInterfaceToDocumentComment(comment)
	return a.collaborationPort.UpdateComment(ctx, commentID, c)
}

func (a *WriterServiceAdapter) DeleteComment(ctx context.Context, commentID string) error {
	return a.collaborationPort.DeleteComment(ctx, commentID)
}

func (a *WriterServiceAdapter) ListComments(ctx context.Context, filter *interface{}, page, pageSize int) ([]*interface{}, int64, error) {
	f := fromInterfaceToCommentFilter(filter)
	comments, total, err := a.collaborationPort.ListComments(ctx, f, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return toInterfaceSliceComment(comments), total, nil
}

func (a *WriterServiceAdapter) GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*interface{}, error) {
	comments, err := a.collaborationPort.GetDocumentComments(ctx, documentID, includeResolved)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceComment(comments), nil
}

func (a *WriterServiceAdapter) GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*interface{}, error) {
	comments, err := a.collaborationPort.GetChapterComments(ctx, chapterID, includeResolved)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceComment(comments), nil
}

func (a *WriterServiceAdapter) ResolveComment(ctx context.Context, commentID, userID string) error {
	return a.collaborationPort.ResolveComment(ctx, commentID, userID)
}

func (a *WriterServiceAdapter) UnresolveComment(ctx context.Context, commentID string) error {
	return a.collaborationPort.UnresolveComment(ctx, commentID)
}

func (a *WriterServiceAdapter) ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*interface{}, error) {
	comment, err := a.collaborationPort.ReplyComment(ctx, parentID, content, userID, userName)
	if err != nil {
		return nil, err
	}
	return toInterface(comment), nil
}

func (a *WriterServiceAdapter) GetCommentThread(ctx context.Context, threadID string) (*interface{}, error) {
	thread, err := a.collaborationPort.GetCommentThread(ctx, threadID)
	if err != nil {
		return nil, err
	}
	return toInterface(thread), nil
}

func (a *WriterServiceAdapter) GetCommentReplies(ctx context.Context, parentID string) ([]*interface{}, error) {
	replies, err := a.collaborationPort.GetCommentReplies(ctx, parentID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceComment(replies), nil
}

func (a *WriterServiceAdapter) GetCommentStats(ctx context.Context, documentID string) (*interface{}, error) {
	stats, err := a.collaborationPort.GetCommentStats(ctx, documentID)
	if err != nil {
		return nil, err
	}
	return toInterface(stats), nil
}

func (a *WriterServiceAdapter) BatchDeleteComments(ctx context.Context, commentIDs []string) error {
	return a.collaborationPort.BatchDeleteComments(ctx, commentIDs)
}

func (a *WriterServiceAdapter) SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*interface{}, int64, error) {
	comments, total, err := a.collaborationPort.SearchComments(ctx, keyword, documentID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return toInterfaceSliceComment(comments), total, nil
}

// ============================================================================
// 内容管理方法 - 委托给 ContentManagementPort
// ============================================================================

// 角色管理
func (a *WriterServiceAdapter) CreateCharacter(ctx context.Context, projectID, userID string, req *writeriface.CreateCharacterRequest) (*interface{}, error) {
	character, err := a.contentPort.CreateCharacter(ctx, projectID, userID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(character), nil
}

func (a *WriterServiceAdapter) GetCharacterByID(ctx context.Context, characterID, projectID string) (*interface{}, error) {
	character, err := a.contentPort.GetCharacterByID(ctx, characterID, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(character), nil
}

func (a *WriterServiceAdapter) ListCharacters(ctx context.Context, projectID string) ([]*interface{}, error) {
	characters, err := a.contentPort.ListCharacters(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceCharacter(characters), nil
}

func (a *WriterServiceAdapter) UpdateCharacter(ctx context.Context, characterID, projectID string, req *writeriface.UpdateCharacterRequest) (*interface{}, error) {
	character, err := a.contentPort.UpdateCharacter(ctx, characterID, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(character), nil
}

func (a *WriterServiceAdapter) DeleteCharacter(ctx context.Context, characterID, projectID string) error {
	return a.contentPort.DeleteCharacter(ctx, characterID, projectID)
}

func (a *WriterServiceAdapter) CreateCharacterRelation(ctx context.Context, projectID string, req *writeriface.CreateRelationRequest) (*interface{}, error) {
	relation, err := a.contentPort.CreateCharacterRelation(ctx, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(relation), nil
}

func (a *WriterServiceAdapter) ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*interface{}, error) {
	relations, err := a.contentPort.ListCharacterRelations(ctx, projectID, characterID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceCharacterRelation(relations), nil
}

func (a *WriterServiceAdapter) DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error {
	return a.contentPort.DeleteCharacterRelation(ctx, relationID, projectID)
}

func (a *WriterServiceAdapter) GetCharacterGraph(ctx context.Context, projectID string) (*writeriface.CharacterGraph, error) {
	return a.contentPort.GetCharacterGraph(ctx, projectID)
}

// 地点管理
func (a *WriterServiceAdapter) CreateLocation(ctx context.Context, projectID, userID string, req *writeriface.CreateLocationRequest) (*interface{}, error) {
	location, err := a.contentPort.CreateLocation(ctx, projectID, userID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(location), nil
}

func (a *WriterServiceAdapter) GetLocationByID(ctx context.Context, locationID, projectID string) (*interface{}, error) {
	location, err := a.contentPort.GetLocationByID(ctx, locationID, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(location), nil
}

func (a *WriterServiceAdapter) ListLocations(ctx context.Context, projectID string) ([]*interface{}, error) {
	locations, err := a.contentPort.ListLocations(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceLocation(locations), nil
}

func (a *WriterServiceAdapter) UpdateLocation(ctx context.Context, locationID, projectID string, req *writeriface.UpdateLocationRequest) (*interface{}, error) {
	location, err := a.contentPort.UpdateLocation(ctx, locationID, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(location), nil
}

func (a *WriterServiceAdapter) DeleteLocation(ctx context.Context, locationID, projectID string) error {
	return a.contentPort.DeleteLocation(ctx, locationID, projectID)
}

func (a *WriterServiceAdapter) GetLocationTree(ctx context.Context, projectID string) ([]*writeriface.LocationNode, error) {
	return a.contentPort.GetLocationTree(ctx, projectID)
}

func (a *WriterServiceAdapter) GetLocationPath(ctx context.Context, locationID string) ([]string, error) {
	return a.contentPort.GetLocationPath(ctx, locationID)
}

func (a *WriterServiceAdapter) CreateLocationRelation(ctx context.Context, projectID string, req *writeriface.CreateLocationRelationRequest) (*interface{}, error) {
	relation, err := a.contentPort.CreateLocationRelation(ctx, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(relation), nil
}

func (a *WriterServiceAdapter) ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*interface{}, error) {
	relations, err := a.contentPort.ListLocationRelations(ctx, projectID, locationID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceLocationRelation(relations), nil
}

func (a *WriterServiceAdapter) DeleteLocationRelation(ctx context.Context, relationID, projectID string) error {
	return a.contentPort.DeleteLocationRelation(ctx, relationID, projectID)
}

// 时间线管理
func (a *WriterServiceAdapter) CreateTimeline(ctx context.Context, projectID string, req *writeriface.CreateTimelineRequest) (*interface{}, error) {
	timeline, err := a.contentPort.CreateTimeline(ctx, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(timeline), nil
}

func (a *WriterServiceAdapter) GetTimeline(ctx context.Context, timelineID, projectID string) (*interface{}, error) {
	timeline, err := a.contentPort.GetTimeline(ctx, timelineID, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(timeline), nil
}

func (a *WriterServiceAdapter) ListTimelines(ctx context.Context, projectID string) ([]*interface{}, error) {
	timelines, err := a.contentPort.ListTimelines(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceTimeline(timelines), nil
}

func (a *WriterServiceAdapter) DeleteTimeline(ctx context.Context, timelineID, projectID string) error {
	return a.contentPort.DeleteTimeline(ctx, timelineID, projectID)
}

func (a *WriterServiceAdapter) CreateTimelineEvent(ctx context.Context, projectID string, req *writeriface.CreateTimelineEventRequest) (*interface{}, error) {
	event, err := a.contentPort.CreateTimelineEvent(ctx, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(event), nil
}

func (a *WriterServiceAdapter) GetTimelineEvent(ctx context.Context, eventID, projectID string) (*interface{}, error) {
	event, err := a.contentPort.GetTimelineEvent(ctx, eventID, projectID)
	if err != nil {
		return nil, err
	}
	return toInterface(event), nil
}

func (a *WriterServiceAdapter) ListTimelineEvents(ctx context.Context, timelineID string) ([]*interface{}, error) {
	events, err := a.contentPort.ListTimelineEvents(ctx, timelineID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSliceTimelineEvent(events), nil
}

func (a *WriterServiceAdapter) UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *writeriface.UpdateTimelineEventRequest) (*interface{}, error) {
	event, err := a.contentPort.UpdateTimelineEvent(ctx, eventID, projectID, req)
	if err != nil {
		return nil, err
	}
	return toInterface(event), nil
}

func (a *WriterServiceAdapter) DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error {
	return a.contentPort.DeleteTimelineEvent(ctx, eventID, projectID)
}

func (a *WriterServiceAdapter) GetTimelineVisualization(ctx context.Context, timelineID string) (*writeriface.TimelineVisualization, error) {
	return a.contentPort.GetTimelineVisualization(ctx, timelineID)
}

// ============================================================================
// 发布导出方法 - 委托给 PublishExportPort
// ============================================================================

// 发布管理
func (a *WriterServiceAdapter) PublishProject(ctx context.Context, projectID, userID string, req *writeriface.PublishProjectRequest) (*writeriface.PublicationRecord, error) {
	return a.publishExportPort.PublishProject(ctx, projectID, userID, req)
}

func (a *WriterServiceAdapter) UnpublishProject(ctx context.Context, projectID, userID string) error {
	return a.publishExportPort.UnpublishProject(ctx, projectID, userID)
}

func (a *WriterServiceAdapter) GetProjectPublicationStatus(ctx context.Context, projectID string) (*writeriface.PublicationStatus, error) {
	return a.publishExportPort.GetProjectPublicationStatus(ctx, projectID)
}

func (a *WriterServiceAdapter) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.PublishDocumentRequest) (*writeriface.PublicationRecord, error) {
	return a.publishExportPort.PublishDocument(ctx, documentID, projectID, userID, req)
}

func (a *WriterServiceAdapter) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *writeriface.UpdateDocumentPublishStatusRequest) error {
	return a.publishExportPort.UpdateDocumentPublishStatus(ctx, documentID, projectID, userID, req)
}

func (a *WriterServiceAdapter) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *writeriface.BatchPublishDocumentsRequest) (*writeriface.BatchPublishResult, error) {
	return a.publishExportPort.BatchPublishDocuments(ctx, projectID, userID, req)
}

func (a *WriterServiceAdapter) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.PublicationRecord, int64, error) {
	return a.publishExportPort.GetPublicationRecords(ctx, projectID, page, pageSize)
}

func (a *WriterServiceAdapter) GetPublicationRecord(ctx context.Context, recordID string) (*writeriface.PublicationRecord, error) {
	return a.publishExportPort.GetPublicationRecord(ctx, recordID)
}

// 导出管理
func (a *WriterServiceAdapter) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.ExportDocumentRequest) (*writeriface.ExportTask, error) {
	return a.publishExportPort.ExportDocument(ctx, documentID, projectID, userID, req)
}

func (a *WriterServiceAdapter) ExportProject(ctx context.Context, projectID, userID string, req *writeriface.ExportProjectRequest) (*writeriface.ExportTask, error) {
	return a.publishExportPort.ExportProject(ctx, projectID, userID, req)
}

func (a *WriterServiceAdapter) GetExportTask(ctx context.Context, taskID string) (*writeriface.ExportTask, error) {
	return a.publishExportPort.GetExportTask(ctx, taskID)
}

func (a *WriterServiceAdapter) DownloadExportFile(ctx context.Context, taskID string) (*writeriface.ExportFile, error) {
	return a.publishExportPort.DownloadExportFile(ctx, taskID)
}

func (a *WriterServiceAdapter) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.ExportTask, int64, error) {
	return a.publishExportPort.ListExportTasks(ctx, projectID, page, pageSize)
}

func (a *WriterServiceAdapter) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	return a.publishExportPort.DeleteExportTask(ctx, taskID, userID)
}

func (a *WriterServiceAdapter) CancelExportTask(ctx context.Context, taskID, userID string) error {
	return a.publishExportPort.CancelExportTask(ctx, taskID, userID)
}
