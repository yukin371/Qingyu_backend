package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/document"
	"Qingyu_backend/service/project"
)

// InitWriterRouter 初始化写作端路由
func InitWriterRouter(r *gin.RouterGroup, projectService *project.ProjectService, documentService *document.DocumentService) {
	// 创建API实例
	projectApi := writer.NewProjectApi(projectService)
	documentApi := writer.NewDocumentApi(documentService)

	// 写作端路由组
	writerGroup := r.Group("/writer")
	writerGroup.Use(middleware.JWTAuth())
	{
		// 项目管理路由
		InitProjectRouter(writerGroup, projectApi)

		// 文档管理路由
		InitDocumentRouter(writerGroup, documentApi)
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
	}
}

// InitDocumentRouter 初始化文档路由
func InitDocumentRouter(r *gin.RouterGroup, documentApi *writer.DocumentApi) {
	// 项目下的文档管理
	projectDocGroup := r.Group("/projects/:projectId/documents")
	{
		projectDocGroup.POST("", documentApi.CreateDocument)
		projectDocGroup.GET("/tree", documentApi.GetDocumentTree)
	}

	// 文档直接操作
	documentGroup := r.Group("/documents")
	{
		documentGroup.GET("/:id", documentApi.GetDocument)
		documentGroup.PUT("/:id", documentApi.UpdateDocument)
		documentGroup.DELETE("/:id", documentApi.DeleteDocument)
	}
}
