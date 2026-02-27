package content

import (
	"github.com/gin-gonic/gin"
)

// ContentAPI 内容管理API路由聚合
// 提供统一的内容管理相关路由，包括文档、章节、进度和项目管理
type ContentAPI struct {
	documentAPI *DocumentAPI
	chapterAPI  *ChapterAPI
	progressAPI *ProgressAPI
	projectAPI  *ProjectAPI
}

// NewContentAPI 创建内容管理API实例
func NewContentAPI(
	documentAPI *DocumentAPI,
	chapterAPI *ChapterAPI,
	progressAPI *ProgressAPI,
	projectAPI *ProjectAPI,
) *ContentAPI {
	return &ContentAPI{
		documentAPI: documentAPI,
		chapterAPI:  chapterAPI,
		progressAPI: progressAPI,
		projectAPI:  projectAPI,
	}
}

// RegisterRoutes 注册内容管理路由
// 路由前缀: /api/v1/content
//
// 路由结构:
// - /documents      文档管理
// - /books/:bookId/chapters  章节内容
// - /progress       阅读进度
// - /projects       项目管理
func (api *ContentAPI) RegisterRoutes(router *gin.RouterGroup) {
	contentGroup := router.Group("/content")
	{
		// 文档管理路由
		api.registerDocumentRoutes(contentGroup)

		// 章节内容路由
		api.registerChapterRoutes(contentGroup)

		// 阅读进度路由
		api.registerProgressRoutes(contentGroup)

		// 项目管理路由
		api.registerProjectRoutes(contentGroup)
	}
}

// registerDocumentRoutes 注册文档管理路由
func (api *ContentAPI) registerDocumentRoutes(group *gin.RouterGroup) {
	documents := group.Group("/documents")
	{
		// 基础CRUD操作
		documents.POST("", api.documentAPI.CreateDocument)
		documents.GET("/:id", api.documentAPI.GetDocument)
		documents.PUT("/:id", api.documentAPI.UpdateDocument)
		documents.DELETE("/:id", api.documentAPI.DeleteDocument)
		documents.GET("", api.documentAPI.ListDocuments)

		// 文档操作
		documents.POST("/:id/duplicate", api.documentAPI.DuplicateDocument)
		documents.PUT("/:id/move", api.documentAPI.MoveDocument)
		documents.GET("/tree/:projectId", api.documentAPI.GetDocumentTree)

		// 文档内容管理
		documents.GET("/:id/content", api.documentAPI.GetDocumentContent)
		documents.PUT("/:id/content", api.documentAPI.UpdateDocumentContent)
		documents.POST("/autosave", api.documentAPI.AutoSaveDocument)

		// 版本控制
		documents.GET("/:id/versions", api.documentAPI.GetVersionHistory)
		documents.POST("/:id/versions/:versionId/restore", api.documentAPI.RestoreVersion)
	}
}

// registerChapterRoutes 注册章节内容路由
func (api *ContentAPI) registerChapterRoutes(group *gin.RouterGroup) {
	chapters := group.Group("/books/:bookId/chapters")
	{
		// 章节内容获取
		chapters.GET("/:chapterId", api.chapterAPI.GetChapter)
		chapters.GET("", api.chapterAPI.ListChapters)

		// 章节导航
		chapters.GET("/:chapterId/next", api.chapterAPI.GetNextChapter)
		chapters.GET("/:chapterId/previous", api.chapterAPI.GetPreviousChapter)
		chapters.GET("/by-number/:chapterNum", api.chapterAPI.GetChapterByNumber)

		// 章节信息（不含内容）
		chapters.GET("/:chapterId/info", api.chapterAPI.GetChapterInfo)
	}
}

// registerProgressRoutes 注册阅读进度路由
func (api *ContentAPI) registerProgressRoutes(group *gin.RouterGroup) {
	progress := group.Group("/progress")
	{
		// 进度管理
		progress.GET("/:bookId", api.progressAPI.GetProgress)
		progress.POST("", api.progressAPI.SaveProgress)
		progress.PUT("/reading-time", api.progressAPI.UpdateReadingTime)

		// 最近阅读
		progress.GET("/recent", api.progressAPI.GetRecentBooks)

		// 阅读统计
		progress.GET("/stats", api.progressAPI.GetReadingStats)

		// 阅读历史
		progress.GET("/history", api.progressAPI.GetReadingHistory)

		// 书籍状态
		progress.GET("/unfinished", api.progressAPI.GetUnfinishedBooks)
		progress.GET("/finished", api.progressAPI.GetFinishedBooks)
	}
}

// registerProjectRoutes 注册项目管理路由
func (api *ContentAPI) registerProjectRoutes(group *gin.RouterGroup) {
	projects := group.Group("/projects")
	{
		// 基础CRUD操作
		projects.POST("", api.projectAPI.CreateProject)
		projects.GET("/:id", api.projectAPI.GetProject)
		projects.PUT("/:id", api.projectAPI.UpdateProject)
		projects.DELETE("/:id", api.projectAPI.DeleteProject)
		projects.GET("", api.projectAPI.ListProjects)

		// 项目统计
		projects.GET("/:id/statistics", api.projectAPI.GetProjectStatistics)
		projects.PUT("/:id/statistics", api.projectAPI.UpdateProjectStatistics)

		// 项目操作
		projects.POST("/:id/duplicate", api.projectAPI.DuplicateProject)
		projects.PUT("/:id/archive", api.projectAPI.ArchiveProject)
		projects.DELETE("/:id/archive", api.projectAPI.UnarchiveProject)

		// 协作管理
		projects.GET("/:id/collaborators", api.projectAPI.ListCollaborators)
		projects.POST("/:id/collaborators", api.projectAPI.AddCollaborator)
		projects.DELETE("/:id/collaborators/:userId", api.projectAPI.RemoveCollaborator)
		projects.PUT("/:id/collaborators/:userId", api.projectAPI.UpdateCollaboratorRole)

		// 导出
		projects.GET("/:id/export", api.projectAPI.GetExportTasks)
		projects.POST("/:id/export", api.projectAPI.ExportProject)
	}
}
