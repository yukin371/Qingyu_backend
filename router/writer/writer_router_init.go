package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/service"
	documentService "Qingyu_backend/service/document"
	projectService "Qingyu_backend/service/project"
)

// RegisterWriterRoutes 注册所有写作相关路由到 /api/v1/writer
func RegisterWriterRoutes(r *gin.RouterGroup) {
	// 从服务容器获取依赖
	serviceContainer := service.GetServiceContainer()
	if serviceContainer == nil {
		// 如果服务容器未初始化，跳过路由注册
		return
	}

	// 获取Repository工厂和EventBus
	repositoryFactory := serviceContainer.GetRepositoryFactory()
	if repositoryFactory == nil {
		return
	}

	eventBus := serviceContainer.GetEventBus()

	// 创建ProjectRepository
	projectRepo := repositoryFactory.CreateProjectRepository()

	// 创建ProjectService
	projectSvc := projectService.NewProjectService(projectRepo, eventBus)

	// 创建DocumentService所需的Repository
	documentRepo := repositoryFactory.CreateDocumentRepository()
	documentContentRepo := repositoryFactory.CreateDocumentContentRepository()
	projectRepoForDoc := repositoryFactory.CreateProjectRepository()

	// 创建DocumentService
	documentSvc := documentService.NewDocumentService(documentRepo, documentContentRepo, projectRepoForDoc, eventBus)

	// 获取MongoDB数据库连接用于VersionService
	mongoDB := serviceContainer.GetMongoDB()

	// 创建VersionService（直接使用MongoDB数据库）
	versionSvc := projectService.NewVersionService(mongoDB)

	// 调用InitWriterRouter初始化所有写作路由
	InitWriterRouter(r, projectSvc, documentSvc, versionSvc)
}
