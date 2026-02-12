package writer

import (
	"context"

	writeriface "Qingyu_backend/service/interfaces/writer"
	writermigration "Qingyu_backend/service/writer/_migration"
)

// WriterService 统一的 Writer 服务接口
// 兼容层适配器实现的接口，供 API 层使用
type WriterService interface {
	// BaseService 方法
	Initialize(ctx context.Context) error
	Health(ctx context.Context) error
	Close(ctx context.Context) error
	GetServiceName() string
	GetVersion() string

	// 项目管理方法
	CreateProject(ctx context.Context, req *writeriface.CreateProjectRequest) (*writeriface.CreateProjectResponse, error)
	GetProject(ctx context.Context, projectID string) (*interface{}, error)
	GetProjectByID(ctx context.Context, projectID string) (*interface{}, error)
	GetByIDWithoutAuth(ctx context.Context, projectID string) (*interface{}, error)
	ListMyProjects(ctx context.Context, req *writeriface.ListProjectsRequest) (*writeriface.ListProjectsResponse, error)
	GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*interface{}, error)
	UpdateProject(ctx context.Context, projectID string, req *writeriface.UpdateProjectRequest) error
	UpdateProjectByID(ctx context.Context, projectID, userID string, req *interface{}) error
	DeleteProject(ctx context.Context, projectID string) error
	DeleteProjectByID(ctx context.Context, projectID, userID string) error
	RestoreProjectByID(ctx context.Context, projectID, userID string) error
	DeleteHard(ctx context.Context, projectID string) error
	UpdateProjectStatistics(ctx context.Context, projectID string, stats *interface{}) error
	RecalculateProjectStatistics(ctx context.Context, projectID string) error

	// 文档管理方法
	CreateDocument(ctx context.Context, req *writeriface.CreateDocumentRequest) (*writeriface.CreateDocumentResponse, error)
	GetDocument(ctx context.Context, documentID string) (*interface{}, error)
	GetDocumentTree(ctx context.Context, projectID string) (*writeriface.DocumentTreeResponse, error)
	ListDocuments(ctx context.Context, req *writeriface.ListDocumentsRequest) (*writeriface.ListDocumentsResponse, error)
	UpdateDocument(ctx context.Context, documentID string, req *writeriface.UpdateDocumentRequest) error
	DeleteDocument(ctx context.Context, documentID string) error
	MoveDocument(ctx context.Context, req *writeriface.MoveDocumentRequest) error
	ReorderDocuments(ctx context.Context, req *writeriface.ReorderDocumentsRequest) error
	AutoSaveDocument(ctx context.Context, req *writeriface.AutoSaveRequest) (*writeriface.AutoSaveResponse, error)
	GetSaveStatus(ctx context.Context, documentID string) (*writeriface.SaveStatusResponse, error)
	GetDocumentContent(ctx context.Context, documentID string) (*writeriface.DocumentContentResponse, error)
	UpdateDocumentContent(ctx context.Context, req *writeriface.UpdateContentRequest) error
	DuplicateDocument(ctx context.Context, documentID string, req *writeriface.DuplicateRequest) (*writeriface.DuplicateResponse, error)

	// 协作批注方法
	CreateComment(ctx context.Context, comment *interface{}) (*interface{}, error)
	GetComment(ctx context.Context, commentID string) (*interface{}, error)
	UpdateComment(ctx context.Context, commentID string, comment *interface{}) error
	DeleteComment(ctx context.Context, commentID string) error
	ListComments(ctx context.Context, filter *interface{}, page, pageSize int) ([]*interface{}, int64, error)
	GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*interface{}, error)
	GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*interface{}, error)
	ResolveComment(ctx context.Context, commentID, userID string) error
	UnresolveComment(ctx context.Context, commentID string) error
	ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*interface{}, error)
	GetCommentThread(ctx context.Context, threadID string) (*interface{}, error)
	GetCommentReplies(ctx context.Context, parentID string) ([]*interface{}, error)
	GetCommentStats(ctx context.Context, documentID string) (*interface{}, error)
	BatchDeleteComments(ctx context.Context, commentIDs []string) error
	SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*interface{}, int64, error)

	// 内容管理方法 - 角色
	CreateCharacter(ctx context.Context, projectID, userID string, req *writeriface.CreateCharacterRequest) (*interface{}, error)
	GetCharacterByID(ctx context.Context, characterID, projectID string) (*interface{}, error)
	ListCharacters(ctx context.Context, projectID string) ([]*interface{}, error)
	UpdateCharacter(ctx context.Context, characterID, projectID string, req *writeriface.UpdateCharacterRequest) (*interface{}, error)
	DeleteCharacter(ctx context.Context, characterID, projectID string) error
	CreateCharacterRelation(ctx context.Context, projectID string, req *writeriface.CreateRelationRequest) (*interface{}, error)
	ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*interface{}, error)
	DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error
	GetCharacterGraph(ctx context.Context, projectID string) (*writeriface.CharacterGraph, error)

	// 内容管理方法 - 地点
	CreateLocation(ctx context.Context, projectID, userID string, req *writeriface.CreateLocationRequest) (*interface{}, error)
	GetLocationByID(ctx context.Context, locationID, projectID string) (*interface{}, error)
	ListLocations(ctx context.Context, projectID string) ([]*interface{}, error)
	UpdateLocation(ctx context.Context, locationID, projectID string, req *writeriface.UpdateLocationRequest) (*interface{}, error)
	DeleteLocation(ctx context.Context, locationID, projectID string) error
	GetLocationTree(ctx context.Context, projectID string) ([]*writeriface.LocationNode, error)
	GetLocationPath(ctx context.Context, locationID string) ([]string, error)
	CreateLocationRelation(ctx context.Context, projectID string, req *writeriface.CreateLocationRelationRequest) (*interface{}, error)
	ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*interface{}, error)
	DeleteLocationRelation(ctx context.Context, relationID, projectID string) error

	// 内容管理方法 - 时间线
	CreateTimeline(ctx context.Context, projectID string, req *writeriface.CreateTimelineRequest) (*interface{}, error)
	GetTimeline(ctx context.Context, timelineID, projectID string) (*interface{}, error)
	ListTimelines(ctx context.Context, projectID string) ([]*interface{}, error)
	DeleteTimeline(ctx context.Context, timelineID, projectID string) error
	CreateTimelineEvent(ctx context.Context, projectID string, req *writeriface.CreateTimelineEventRequest) (*interface{}, error)
	GetTimelineEvent(ctx context.Context, eventID, projectID string) (*interface{}, error)
	ListTimelineEvents(ctx context.Context, timelineID string) ([]*interface{}, error)
	UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *writeriface.UpdateTimelineEventRequest) (*interface{}, error)
	DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error
	GetTimelineVisualization(ctx context.Context, timelineID string) (*writeriface.TimelineVisualization, error)

	// 发布导出方法
	PublishProject(ctx context.Context, projectID, userID string, req *writeriface.PublishProjectRequest) (*writeriface.PublicationRecord, error)
	UnpublishProject(ctx context.Context, projectID, userID string) error
	GetProjectPublicationStatus(ctx context.Context, projectID string) (*writeriface.PublicationStatus, error)
	PublishDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.PublishDocumentRequest) (*writeriface.PublicationRecord, error)
	UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *writeriface.UpdateDocumentPublishStatusRequest) error
	BatchPublishDocuments(ctx context.Context, projectID, userID string, req *writeriface.BatchPublishDocumentsRequest) (*writeriface.BatchPublishResult, error)
	GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.PublicationRecord, int64, error)
	GetPublicationRecord(ctx context.Context, recordID string) (*writeriface.PublicationRecord, error)
	ExportDocument(ctx context.Context, documentID, projectID, userID string, req *writeriface.ExportDocumentRequest) (*writeriface.ExportTask, error)
	ExportProject(ctx context.Context, projectID, userID string, req *writeriface.ExportProjectRequest) (*writeriface.ExportTask, error)
	GetExportTask(ctx context.Context, taskID string) (*writeriface.ExportTask, error)
	DownloadExportFile(ctx context.Context, taskID string) (*writeriface.ExportFile, error)
	ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*writeriface.ExportTask, int64, error)
	DeleteExportTask(ctx context.Context, taskID, userID string) error
	CancelExportTask(ctx context.Context, taskID, userID string) error
}

// WriterServiceFactory Writer 服务工厂
// 提供创建和组装 Writer 服务的方法
type WriterServiceFactory struct{}

// NewWriterServiceFactory 创建工厂实例
func NewWriterServiceFactory() *WriterServiceFactory {
	return &WriterServiceFactory{}
}

// CreateWithPorts 使用 Port 接口创建服务（推荐方式）
//
// 新架构推荐的使用方式：
// 1. 实现 5 个 Port 接口的具体实现
// 2. 使用 WriterServiceAdapter 组装它们
// 3. 返回 WriterService 接口供 API 层使用
func (f *WriterServiceFactory) CreateWithPorts(
	projectPort writeriface.ProjectManagementPort,
	documentPort writeriface.DocumentManagementPort,
	collaborationPort writeriface.CollaborationPort,
	contentPort writeriface.ContentManagementPort,
	publishExportPort writeriface.PublishExportPort,
) WriterService {
	return writermigration.NewWriterServiceAdapter(
		projectPort,
		documentPort,
		collaborationPort,
		contentPort,
		publishExportPort,
	)
}

// PortImplementations Port 接口实现集合
type PortImplementations struct {
	ProjectPort      writeriface.ProjectManagementPort
	DocumentPort     writeriface.DocumentManagementPort
	CollaborationPort writeriface.CollaborationPort
	ContentPort      writeriface.ContentManagementPort
	PublishExportPort writeriface.PublishExportPort
}

// CreateFromImplementations 从结构体创建服务
func (f *WriterServiceFactory) CreateFromImplementations(ports PortImplementations) WriterService {
	return writermigration.NewWriterServiceAdapter(
		ports.ProjectPort,
		ports.DocumentPort,
		ports.CollaborationPort,
		ports.ContentPort,
		ports.PublishExportPort,
	)
}
