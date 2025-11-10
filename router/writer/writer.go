package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/document"
	"Qingyu_backend/service/project"
)

// InitWriterRouter 初始化写作端路由
func InitWriterRouter(r *gin.RouterGroup, projectService *project.ProjectService, documentService *document.DocumentService, versionService *project.VersionService) {
	// 创建API实例
	projectApi := writer.NewProjectApi(projectService)
	documentApi := writer.NewDocumentApi(documentService)
	versionApi := writer.NewVersionApi(versionService)
	editorApi := writer.NewEditorApi(documentService)

	// 写作端路由组
	writerGroup := r.Group("/writer")
	writerGroup.Use(middleware.JWTAuth())
	{
		// 项目管理路由
		InitProjectRouter(writerGroup, projectApi)

		// 文档管理路由
		InitDocumentRouter(writerGroup, documentApi)

		// 版本控制路由
		InitVersionRouter(writerGroup, versionApi)

		// 编辑器路由
		InitEditorRouter(writerGroup, editorApi)
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
func InitDocumentRouter(r *gin.RouterGroup, documentApi *writer.DocumentApi) {
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
