package content

import (
	"context"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/service/interfaces/base"
)

// ============================================================================
// Content 模块 Port 接口定义
// ============================================================================

// DocumentServicePort 文档管理端口
// 负责文档的 CRUD 操作、树形结构管理和版本控制
type DocumentServicePort interface {
	base.BaseService

	// CreateDocument 创建新文档
	CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error)

	// UpdateDocument 更新文档信息
	UpdateDocument(ctx context.Context, id string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error)

	// GetDocument 获取文档详情
	GetDocument(ctx context.Context, id string) (*dto.DocumentResponse, error)

	// DeleteDocument 删除文档（软删除）
	DeleteDocument(ctx context.Context, id string) error

	// ListDocuments 获取文档列表
	ListDocuments(ctx context.Context, req *dto.ListDocumentsRequest) (*dto.ListDocumentsResponse, error)

	// DuplicateDocument 复制文档
	DuplicateDocument(ctx context.Context, id string) (*dto.DocumentResponse, error)

	// MoveDocument 移动文档到新的父节点
	MoveDocument(ctx context.Context, id string, newParentID string, order int) error

	// GetDocumentTree 获取文档树结构
	GetDocumentTree(ctx context.Context, projectID string) (*dto.DocumentTreeResponse, error)

	// AutoSaveDocument 自动保存文档内容
	AutoSaveDocument(ctx context.Context, req *dto.AutoSaveRequest) (*dto.AutoSaveResponse, error)

	// GetDocumentContent 获取文档内容
	GetDocumentContent(ctx context.Context, documentID string) (*dto.DocumentContentResponse, error)

	// UpdateDocumentContent 更新文档内容
	UpdateDocumentContent(ctx context.Context, req *dto.UpdateContentRequest) error

	// GetVersionHistory 获取版本历史
	GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*dto.VersionHistoryResponse, error)

	// RestoreVersion 恢复到指定版本
	RestoreVersion(ctx context.Context, documentID, versionID string) error
}

// ChapterServicePort 章节内容端口
// 负责章节内容的获取和导航功能
type ChapterServicePort interface {
	base.BaseService

	// GetChapter 获取章节内容
	GetChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error)

	// GetChapterByNumber 根据章节号获取章节
	GetChapterByNumber(ctx context.Context, bookID string, chapterNum int) (*dto.ChapterResponse, error)

	// GetNextChapter 获取下一章
	GetNextChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error)

	// GetPreviousChapter 获取上一章
	GetPreviousChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error)

	// ListChapters 获取章节列表
	ListChapters(ctx context.Context, bookID string) (*dto.ChapterListResponse, error)

	// GetChapterInfo 获取章节信息（不含内容）
	GetChapterInfo(ctx context.Context, chapterID string) (*dto.ChapterInfo, error)
}

// ReadingProgressServicePort 阅读进度端口
// 负责阅读进度的保存、查询和统计功能
type ReadingProgressServicePort interface {
	base.BaseService

	// GetProgress 获取阅读进度
	GetProgress(ctx context.Context, userID, bookID string) (*dto.ReadingProgressResponse, error)

	// SaveProgress 保存阅读进度
	SaveProgress(ctx context.Context, req *dto.SaveProgressRequest) error

	// UpdateReadingTime 更新阅读时长
	UpdateReadingTime(ctx context.Context, userID, bookID string, duration int) error

	// GetRecentBooks 获取最近阅读的书籍
	GetRecentBooks(ctx context.Context, userID string, limit int) (*dto.RecentBooksResponse, error)

	// GetReadingStats 获取阅读统计信息
	GetReadingStats(ctx context.Context, userID string) (*dto.ReadingStatsResponse, error)

	// GetReadingHistory 获取阅读历史
	GetReadingHistory(ctx context.Context, userID string, page, pageSize int) (*dto.ReadingHistoryResponse, error)

	// GetUnfinishedBooks 获取未读完的书籍
	GetUnfinishedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error)

	// GetFinishedBooks 获取已读完的书籍
	GetFinishedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error)
}

// ProjectServicePort 项目管理端口
// 负责项目的 CRUD 操作和统计管理
type ProjectServicePort interface {
	base.BaseService

	// CreateProject 创建新项目
	CreateProject(ctx context.Context, req *dto.CreateProjectRequest) (*dto.ProjectResponse, error)

	// UpdateProject 更新项目信息
	UpdateProject(ctx context.Context, id string, req *dto.UpdateProjectRequest) (*dto.ProjectResponse, error)

	// GetProject 获取项目详情
	GetProject(ctx context.Context, id string) (*dto.ProjectResponse, error)

	// DeleteProject 删除项目（软删除）
	DeleteProject(ctx context.Context, id string) error

	// ListProjects 获取项目列表
	ListProjects(ctx context.Context, req *dto.ListProjectsRequest) (*dto.ListProjectsResponse, error)

	// GetProjectStatistics 获取项目统计信息
	GetProjectStatistics(ctx context.Context, projectID string) (*dto.ProjectStatistics, error)

	// UpdateProjectStatistics 更新项目统计信息
	UpdateProjectStatistics(ctx context.Context, projectID string, stats *dto.ProjectStatistics) error
}

// ============================================================================
// 内容管理接口聚合
// ============================================================================

// ContentManagementService 内容管理聚合服务
// 提供内容管理相关的所有核心功能
type ContentManagementService interface {
	// 子服务接口
	DocumentService() DocumentServicePort
	ChapterService() ChapterServicePort
	ReadingProgressService() ReadingProgressServicePort
	ProjectService() ProjectServicePort
}
