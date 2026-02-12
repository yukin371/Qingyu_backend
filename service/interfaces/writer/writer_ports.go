package writer

import (
	"context"
	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/interfaces/base"
)

// ============================================================================
// Writer 模块 Port 接口定义
// ============================================================================

// ProjectManagementPort 项目管理端口
// 负责项目的 CRUD 操作和统计管理
type ProjectManagementPort interface {
	base.BaseService

	// CreateProject 创建新项目
	CreateProject(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error)

	// GetProject 获取项目详情
	GetProject(ctx context.Context, projectID string) (*writer.Project, error)

	// GetProjectByID 根据ID获取项目（别名方法，兼容API层调用）
	GetProjectByID(ctx context.Context, projectID string) (*writer.Project, error)

	// GetByIDWithoutAuth 获取项目详情（无权限检查，用于内部服务调用如AI）
	GetByIDWithoutAuth(ctx context.Context, projectID string) (*writer.Project, error)

	// ListMyProjects 获取我的项目列表
	ListMyProjects(ctx context.Context, req *ListProjectsRequest) (*ListProjectsResponse, error)

	// GetProjectList 获取项目列表（兼容API层调用）
	GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*writer.Project, error)

	// UpdateProject 更新项目
	UpdateProject(ctx context.Context, projectID string, req *UpdateProjectRequest) error

	// UpdateProjectByID 更新项目（兼容API层调用）
	UpdateProjectByID(ctx context.Context, projectID, userID string, req *writer.Project) error

	// DeleteProject 删除项目
	DeleteProject(ctx context.Context, projectID string) error

	// DeleteProjectByID 删除项目（兼容API层调用）
	DeleteProjectByID(ctx context.Context, projectID, userID string) error

	// RestoreProjectByID 恢复已删除的项目
	RestoreProjectByID(ctx context.Context, projectID, userID string) error

	// DeleteHard 物理删除项目（永久删除）
	DeleteHard(ctx context.Context, projectID string) error

	// UpdateProjectStatistics 更新项目统计
	UpdateProjectStatistics(ctx context.Context, projectID string, stats *writer.ProjectStats) error

	// RecalculateProjectStatistics 重新计算项目统计信息
	RecalculateProjectStatistics(ctx context.Context, projectID string) error
}

// DocumentManagementPort 文档管理端口
// 负责文档的 CRUD 操作、版本控制和自动保存
type DocumentManagementPort interface {
	base.BaseService

	// CreateDocument 创建文档
	CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error)

	// GetDocument 获取文档详情
	GetDocument(ctx context.Context, documentID string) (*writer.Document, error)

	// GetDocumentTree 获取文档树
	GetDocumentTree(ctx context.Context, projectID string) (*DocumentTreeResponse, error)

	// ListDocuments 获取文档列表
	ListDocuments(ctx context.Context, req *ListDocumentsRequest) (*ListDocumentsResponse, error)

	// UpdateDocument 更新文档
	UpdateDocument(ctx context.Context, documentID string, req *UpdateDocumentRequest) error

	// DeleteDocument 删除文档
	DeleteDocument(ctx context.Context, documentID string) error

	// MoveDocument 移动文档
	MoveDocument(ctx context.Context, req *MoveDocumentRequest) error

	// ReorderDocuments 重新排序文档
	ReorderDocuments(ctx context.Context, req *ReorderDocumentsRequest) error

	// AutoSaveDocument 自动保存文档
	AutoSaveDocument(ctx context.Context, req *AutoSaveRequest) (*AutoSaveResponse, error)

	// GetSaveStatus 获取保存状态
	GetSaveStatus(ctx context.Context, documentID string) (*SaveStatusResponse, error)

	// GetDocumentContent 获取文档内容
	GetDocumentContent(ctx context.Context, documentID string) (*DocumentContentResponse, error)

	// UpdateDocumentContent 更新文档内容
	UpdateDocumentContent(ctx context.Context, req *UpdateContentRequest) error

	// DuplicateDocument 复制文档
	DuplicateDocument(ctx context.Context, documentID string, req *DuplicateRequest) (*DuplicateResponse, error)

	// 版本控制方法

	// BumpVersionAndCreateRevision 创建新版本并记录修订
	BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*writer.FileRevision, error)

	// UpdateContentWithVersion 使用乐观并发控制更新内容
	UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*writer.FileRevision, error)

	// RollbackToVersion 回滚到指定的历史版本
	RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*writer.FileRevision, error)

	// ListRevisions 列表修订
	ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*writer.FileRevision, error)

	// GetVersionHistory 获取版本历史
	GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*VersionHistoryResponse, error)

	// GetVersionDetail 获取特定版本
	GetVersionDetail(ctx context.Context, documentID, versionID string) (*VersionDetail, error)

	// CompareVersions 比较两个版本
	CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*VersionDiff, error)

	// RestoreVersion 恢复到特定版本
	RestoreVersion(ctx context.Context, documentID, versionID string) error

	// CreatePatch 提交一个候选补丁
	CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*writer.FilePatch, error)

	// ApplyPatch 审核并应用补丁
	ApplyPatch(projectID, patchID, applierID string) (*writer.FileRevision, error)

	// DetectConflicts 检测文件的版本冲突
	DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*writer.ConflictInfo, error)

	// CreateCommit 创建批量提交
	CreateCommit(ctx context.Context, projectID, authorID, message string, files []writer.CommitFile) (*writer.Commit, error)

	// ListCommits 查询提交历史
	ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*writer.Commit, error)

	// GetCommitDetails 获取提交详情
	GetCommitDetails(ctx context.Context, projectID, commitID string) (*writer.Commit, []*writer.FileRevision, error)

	// 自动保存方法

	// StartAutoSave 启动文档自动保存
	StartAutoSave(documentID, projectID, nodeID, userID string) error

	// StopAutoSave 停止文档自动保存
	StopAutoSave(documentID string) error

	// GetAutoSaveStatus 获取自动保存状态
	GetAutoSaveStatus(documentID string) (isRunning bool, lastSaved *interface{})

	// SaveImmediately 立即执行一次保存（手动触发）
	SaveImmediately(documentID string) error
}

// CollaborationPort 协作批注端口
// 负责文档批注和协作功能
type CollaborationPort interface {
	base.BaseService

	// CreateComment 创建批注
	CreateComment(ctx context.Context, comment *writer.DocumentComment) (*writer.DocumentComment, error)

	// GetComment 获取批注详情
	GetComment(ctx context.Context, commentID string) (*writer.DocumentComment, error)

	// UpdateComment 更新批注
	UpdateComment(ctx context.Context, commentID string, comment *writer.DocumentComment) error

	// DeleteComment 删除批注
	DeleteComment(ctx context.Context, commentID string) error

	// ListComments 查询批注列表
	ListComments(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error)

	// GetDocumentComments 获取文档的所有批注
	GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*writer.DocumentComment, error)

	// GetChapterComments 获取章节的所有批注
	GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*writer.DocumentComment, error)

	// ResolveComment 标记批注为已解决
	ResolveComment(ctx context.Context, commentID, userID string) error

	// UnresolveComment 标记批注为未解决
	UnresolveComment(ctx context.Context, commentID string) error

	// ReplyComment 回复批注
	ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*writer.DocumentComment, error)

	// GetCommentThread 获取批注线程
	GetCommentThread(ctx context.Context, threadID string) (*writer.CommentThread, error)

	// GetCommentReplies 获取批注的回复
	GetCommentReplies(ctx context.Context, parentID string) ([]*writer.DocumentComment, error)

	// GetCommentStats 获取批注统计
	GetCommentStats(ctx context.Context, documentID string) (*writer.CommentStats, error)

	// BatchDeleteComments 批量删除批注
	BatchDeleteComments(ctx context.Context, commentIDs []string) error

	// SearchComments 搜索批注
	SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*writer.DocumentComment, int64, error)
}

// ContentManagementPort 内容管理端口
// 负责角色、地点、时间线等创作内容管理
type ContentManagementPort interface {
	base.BaseService

	// === 角色管理 ===

	// CreateCharacter 创建角色
	CreateCharacter(ctx context.Context, projectID, userID string, req *CreateCharacterRequest) (*writer.Character, error)

	// GetCharacterByID 根据ID获取角色
	GetCharacterByID(ctx context.Context, characterID, projectID string) (*writer.Character, error)

	// ListCharacters 获取项目下的所有角色
	ListCharacters(ctx context.Context, projectID string) ([]*writer.Character, error)

	// UpdateCharacter 更新角色
	UpdateCharacter(ctx context.Context, characterID, projectID string, req *UpdateCharacterRequest) (*writer.Character, error)

	// DeleteCharacter 删除角色
	DeleteCharacter(ctx context.Context, characterID, projectID string) error

	// CreateCharacterRelation 创建角色关系
	CreateCharacterRelation(ctx context.Context, projectID string, req *CreateRelationRequest) (*writer.CharacterRelation, error)

	// ListCharacterRelations 获取角色关系列表
	ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error)

	// DeleteCharacterRelation 删除角色关系
	DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error

	// GetCharacterGraph 获取角色关系图
	GetCharacterGraph(ctx context.Context, projectID string) (*CharacterGraph, error)

	// === 地点管理 ===

	// CreateLocation 创建地点
	CreateLocation(ctx context.Context, projectID, userID string, req *CreateLocationRequest) (*writer.Location, error)

	// GetLocationByID 根据ID获取地点
	GetLocationByID(ctx context.Context, locationID, projectID string) (*writer.Location, error)

	// ListLocations 获取项目下的所有地点
	ListLocations(ctx context.Context, projectID string) ([]*writer.Location, error)

	// UpdateLocation 更新地点
	UpdateLocation(ctx context.Context, locationID, projectID string, req *UpdateLocationRequest) (*writer.Location, error)

	// DeleteLocation 删除地点
	DeleteLocation(ctx context.Context, locationID, projectID string) error

	// GetLocationTree 获取地点层级树
	GetLocationTree(ctx context.Context, projectID string) ([]*LocationNode, error)

	// GetLocationPath 获取地点的完整路径
	GetLocationPath(ctx context.Context, locationID string) ([]string, error)

	// CreateLocationRelation 创建地点关系
	CreateLocationRelation(ctx context.Context, projectID string, req *CreateLocationRelationRequest) (*writer.LocationRelation, error)

	// ListLocationRelations 获取地点关系列表
	ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error)

	// DeleteLocationRelation 删除地点关系
	DeleteLocationRelation(ctx context.Context, relationID, projectID string) error

	// === 时间线管理 ===

	// CreateTimeline 创建时间线
	CreateTimeline(ctx context.Context, projectID string, req *CreateTimelineRequest) (*writer.Timeline, error)

	// GetTimeline 根据ID获取时间线
	GetTimeline(ctx context.Context, timelineID, projectID string) (*writer.Timeline, error)

	// ListTimelines 获取项目下的所有时间线
	ListTimelines(ctx context.Context, projectID string) ([]*writer.Timeline, error)

	// DeleteTimeline 删除时间线
	DeleteTimeline(ctx context.Context, timelineID, projectID string) error

	// CreateTimelineEvent 创建时间线事件
	CreateTimelineEvent(ctx context.Context, projectID string, req *CreateTimelineEventRequest) (*writer.TimelineEvent, error)

	// GetTimelineEvent 根据ID获取事件
	GetTimelineEvent(ctx context.Context, eventID, projectID string) (*writer.TimelineEvent, error)

	// ListTimelineEvents 获取时间线下的所有事件
	ListTimelineEvents(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error)

	// UpdateTimelineEvent 更新时间线事件
	UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *UpdateTimelineEventRequest) (*writer.TimelineEvent, error)

	// DeleteTimelineEvent 删除时间线事件
	DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error

	// GetTimelineVisualization 获取时间线可视化数据
	GetTimelineVisualization(ctx context.Context, timelineID string) (*TimelineVisualization, error)
}

// PublishExportPort 发布导出端口
// 负责内容发布到书城和导出功能
type PublishExportPort interface {
	base.BaseService

	// === 发布管理 ===

	// PublishProject 发布项目到书城
	PublishProject(ctx context.Context, projectID, userID string, req *PublishProjectRequest) (*PublicationRecord, error)

	// UnpublishProject 取消发布项目
	UnpublishProject(ctx context.Context, projectID, userID string) error

	// GetProjectPublicationStatus 获取项目发布状态
	GetProjectPublicationStatus(ctx context.Context, projectID string) (*PublicationStatus, error)

	// PublishDocument 发布文档（章节）
	PublishDocument(ctx context.Context, documentID, projectID, userID string, req *PublishDocumentRequest) (*PublicationRecord, error)

	// UpdateDocumentPublishStatus 更新文档发布状态
	UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *UpdateDocumentPublishStatusRequest) error

	// BatchPublishDocuments 批量发布文档
	BatchPublishDocuments(ctx context.Context, projectID, userID string, req *BatchPublishDocumentsRequest) (*BatchPublishResult, error)

	// GetPublicationRecords 获取发布记录列表
	GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*PublicationRecord, int64, error)

	// GetPublicationRecord 获取发布记录详情
	GetPublicationRecord(ctx context.Context, recordID string) (*PublicationRecord, error)

	// === 导出管理 ===

	// ExportDocument 导出文档
	ExportDocument(ctx context.Context, documentID, projectID, userID string, req *ExportDocumentRequest) (*ExportTask, error)

	// ExportProject 导出项目
	ExportProject(ctx context.Context, projectID, userID string, req *ExportProjectRequest) (*ExportTask, error)

	// GetExportTask 获取导出任务
	GetExportTask(ctx context.Context, taskID string) (*ExportTask, error)

	// DownloadExportFile 下载导出文件
	DownloadExportFile(ctx context.Context, taskID string) (*ExportFile, error)

	// ListExportTasks 列出导出任务
	ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*ExportTask, int64, error)

	// DeleteExportTask 删除导出任务
	DeleteExportTask(ctx context.Context, taskID, userID string) error

	// CancelExportTask 取消导出任务
	CancelExportTask(ctx context.Context, taskID, userID string) error
}
