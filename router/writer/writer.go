package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/middleware"
	"Qingyu_backend/pkg/lock"
	"Qingyu_backend/service/interfaces"
	searchservice "Qingyu_backend/service/search"
	writerservice "Qingyu_backend/service/writer"
	"Qingyu_backend/service/writer/document"
	projectService "Qingyu_backend/service/writer/project"
)

// InitWriterRouter 初始化写作端路由
func InitWriterRouter(
	r *gin.RouterGroup,
	projectService *projectService.ProjectService,
	documentService *document.DocumentService,
	versionService *projectService.VersionService,
	searchSvc *searchservice.SearchService,
	exportService interfaces.ExportService,
	publishService interfaces.PublishService,
	lockService lock.DocumentLockService,
	commentService writerservice.CommentService,
	batchOpService document.BatchOperationService,
	templateService *document.TemplateService,
) {
	// 创建API实例
	projectApi := writer.NewProjectApi(projectService)
	documentApi := writer.NewDocumentApi(documentService)
	versionApi := writer.NewVersionApi(versionService)
	editorApi := writer.NewEditorApi(documentService)

	// 锁定API（如果可用）
	var lockApi *writer.LockAPI
	if lockService != nil {
		lockApi = writer.NewLockAPI(lockService)
	}

	// 批注API（如果可用）
	var commentApi *writer.CommentAPI
	if commentService != nil {
		commentApi = writer.NewCommentAPI(commentService)
	}

	// 批量操作API（如果可用）
	var batchOpApi *writer.BatchOperationAPI
	if batchOpService != nil {
		batchOpApi = writer.NewBatchOperationAPI(batchOpService)
	}

	// 模板API
	var templateApi *writer.TemplateAPI
	if templateService != nil {
		templateApi = writer.NewTemplateAPI(templateService, nil) // logger由service内部处理
	}

	// 写作端路由组
	writerGroup := r.Group("/writer")
	writerGroup.Use(middleware.JWTAuth())
	{
		// 项目管理路由
		InitProjectRouter(writerGroup, projectApi)

		// 文档管理路由
		InitDocumentRouter(writerGroup, documentApi, lockApi)

		// 版本控制路由
		InitVersionRouter(writerGroup, versionApi)

		// 编辑器路由
		InitEditorRouter(writerGroup, editorApi)

		// 搜索路由
		if searchSvc != nil {
			InitSearchRouter(writerGroup, searchSvc)
		}

		// 导出功能路由
		if exportService != nil {
			InitExportRoutes(writerGroup, exportService)
		}

		// 发布管理路由
		if publishService != nil {
			InitPublishRoutes(writerGroup, publishService)
		}

		// 批注路由
		if commentApi != nil {
			InitCommentRouter(writerGroup, commentApi)
		}

		// 批量操作路由
		if batchOpApi != nil {
			batchOpApi.RegisterRoutes(writerGroup)
		}

		// 模板路由
		if templateApi != nil {
			InitTemplateRouter(writerGroup, templateApi)
		}
	}
}

// InitProjectRouter 初始化项目路由
func InitProjectRouter(r *gin.RouterGroup, projectApi *writer.ProjectApi) {
	projectGroup := r.Group("/projects")
	{
		// 项目CRUD
		projectGroup.POST("", projectApi.CreateProject)
		projectGroup.GET("", projectApi.ListProjects)
		projectGroup.GET("/:id", projectApi.GetProject)
		projectGroup.PUT("/:id", projectApi.UpdateProject)
		projectGroup.DELETE("/:id", projectApi.DeleteProject)

		// 项目统计
		projectGroup.PUT("/:id/statistics", projectApi.UpdateProjectStatistics)
	}
}

// InitDocumentRouter 初始化文档路由
func InitDocumentRouter(r *gin.RouterGroup, documentApi *writer.DocumentApi, lockApi *writer.LockAPI) {
	// 项目下的文档管理
	projectDocGroup := r.Group("/project/:projectId/documents")
	{
		projectDocGroup.POST("", documentApi.CreateDocument)
		projectDocGroup.GET("", documentApi.ListDocuments)
		projectDocGroup.GET("/tree", documentApi.GetDocumentTree)
		projectDocGroup.PUT("/reorder", documentApi.ReorderDocuments)
	}

	// 文档直接操作
	documentGroup := r.Group("/documents")
	{
		documentGroup.GET("/:id", documentApi.GetDocument)
		documentGroup.PUT("/:id", documentApi.UpdateDocument)
		documentGroup.DELETE("/:id", documentApi.DeleteDocument)
		documentGroup.PUT("/:id/move", documentApi.MoveDocument)
		documentGroup.POST("/:id/duplicate", documentApi.DuplicateDocument)

		// 文档锁定路由
		if lockApi != nil {
			documentGroup.POST("/:id/lock", lockApi.LockDocument)
			documentGroup.DELETE("/:id/lock", lockApi.UnlockDocument)
			documentGroup.PUT("/:id/lock/refresh", lockApi.RefreshLock)
			documentGroup.GET("/:id/lock/status", lockApi.GetLockStatus)
			documentGroup.POST("/:id/lock/force", lockApi.ForceUnlock)
			documentGroup.POST("/:id/lock/extend", lockApi.ExtendLock)
		}
	}
}

// InitVersionRouter 初始化版本控制路由
func InitVersionRouter(r *gin.RouterGroup, versionApi *writer.VersionApi) {
	documentGroup := r.Group("/document/:documentId")
	{
		versionGroup := documentGroup.Group("/versions")
		{
			// 版本历史
			versionGroup.GET("", versionApi.GetVersionHistory)
			versionGroup.GET("/:versionId", versionApi.GetVersion)

			// 版本比较
			versionGroup.GET("/compare", versionApi.CompareVersions)

			// 版本恢复
			versionGroup.POST("/:versionId/restore", versionApi.RestoreVersion)
		}
	}
}

// InitEditorRouter 初始化编辑器路由
func InitEditorRouter(r *gin.RouterGroup, editorApi *writer.EditorApi) {
	// 文档编辑相关
	documentGroup := r.Group("/documents/:id")
	{
		// 自动保存
		documentGroup.POST("/autosave", editorApi.AutoSaveDocument)

		// 保存状态
		documentGroup.GET("/save-status", editorApi.GetSaveStatus)

		// 文档内容
		documentGroup.GET("/content", editorApi.GetDocumentContent)
		documentGroup.PUT("/content", editorApi.UpdateDocumentContent)

		// 字数统计
		documentGroup.POST("/word-count", editorApi.CalculateWordCount)
	}

	// 用户快捷键配置
	userGroup := r.Group("/user")
	{
		shortcutGroup := userGroup.Group("/shortcuts")
		{
			shortcutGroup.GET("", editorApi.GetUserShortcuts)
			shortcutGroup.PUT("", editorApi.UpdateUserShortcuts)
			shortcutGroup.POST("/reset", editorApi.ResetUserShortcuts)
			shortcutGroup.GET("/help", editorApi.GetShortcutHelp)
		}
	}
}

// InitSearchRouter 初始化搜索路由
func InitSearchRouter(r *gin.RouterGroup, searchSvc *searchservice.SearchService) {
	// 创建搜索API实例
	searchAPI := writer.NewSearchAPI(searchSvc)

	// 搜索路由组
	searchGroup := r.Group("/search")
	{
		// 文档搜索
		searchGroup.GET("/documents", searchAPI.SearchDocuments)
	}
}

// InitCommentRouter 初始化批注路由
func InitCommentRouter(r *gin.RouterGroup, commentApi *writer.CommentAPI) {
	// 文档批注路由
	documentGroup := r.Group("/documents/:id")
	{
		// 批注管理
		documentGroup.POST("/comments", commentApi.CreateComment)
		documentGroup.GET("/comments", commentApi.GetComments)
		documentGroup.GET("/comments/stats", commentApi.GetCommentStats)
		documentGroup.GET("/comments/search", commentApi.SearchComments)
	}

	// 批注操作路由
	commentGroup := r.Group("/comments")
	{
		// 单个批注操作
		commentGroup.GET("/:id", commentApi.GetComment)
		commentGroup.PUT("/:id", commentApi.UpdateComment)
		commentGroup.DELETE("/:id", commentApi.DeleteComment)

		// 批注状态
		commentGroup.POST("/:id/resolve", commentApi.ResolveComment)
		commentGroup.POST("/:id/unresolve", commentApi.UnresolveComment)
		commentGroup.POST("/:id/reply", commentApi.ReplyComment)

		// 批注线程
		commentGroup.GET("/threads/:threadId", commentApi.GetCommentThread)

		// 批量操作
		commentGroup.POST("/batch-delete", commentApi.BatchDeleteComments)
	}
}

// InitTemplateRouter 初始化模板路由
func InitTemplateRouter(r *gin.RouterGroup, templateApi *writer.TemplateAPI) {
	// 模板管理路由
	templateGroup := r.Group("/templates")
	{
		// 模板CRUD
		templateGroup.POST("", templateApi.CreateTemplate)
		templateGroup.GET("", templateApi.ListTemplates)
		templateGroup.GET("/:id", templateApi.GetTemplate)
		templateGroup.PUT("/:id", templateApi.UpdateTemplate)
		templateGroup.DELETE("/:id", templateApi.DeleteTemplate)

		// 应用模板
		templateGroup.POST("/:id/apply", templateApi.ApplyTemplate)
	}
}
