package writer

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/lock"
	writerrepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/service"
	"Qingyu_backend/service/interfaces"
	searchservice "Qingyu_backend/service/search"
	writerservice "Qingyu_backend/service/writer"
	documentService "Qingyu_backend/service/writer/document"
	projectService "Qingyu_backend/service/writer/project"
)

// RegisterWriterRoutes 注册所有写作相关路由到 /api/v1/writer
func RegisterWriterRoutes(r *gin.RouterGroup, searchSvc *searchservice.SearchService) {
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

	// 创建ExportService（导出服务）
	// 注意：需要先实现ExportTaskRepository和FileStorage接口
	// exportTaskRepo := repositoryFactory.CreateExportTaskRepository()
	// fileStorage := serviceContainer.GetFileStorage()
	// exportSvc := writerService.NewExportService(documentRepo, documentContentRepo, projectRepo, exportTaskRepo, fileStorage)
	var exportSvc interfaces.ExportService

	// 创建PublishService（发布服务）
	// 注意：需要先实现PublicationRepository和BookstoreClient接口
	// publicationRepo := repositoryFactory.CreatePublicationRepository()
	// bookstoreClient := serviceContainer.GetBookstoreClient()
	// publishSvc := writerService.NewPublishService(projectRepo, documentRepo, publicationRepo, bookstoreClient, eventBus)
	var publishSvc interfaces.PublishService

	// 创建DocumentLockService（文档锁服务）
	// 从服务容器获取Redis客户端
	var lockSvc lock.DocumentLockService
	redisClient := serviceContainer.GetRedisClient()
	if redisClient != nil {
		lockSvc = lock.NewRedisDocumentLockService(redisClient, "doclock")
	}

	// 创建CommentService（批注服务）
	var commentSvc writerservice.CommentService
	// 直接创建 writer comment repository（因为 factory 的 CreateCommentRepository 返回的是 reading 模块的）
	commentRepo := writerrepo.NewMongoCommentRepository(mongoDB)
	commentSvc = writerservice.NewCommentService(commentRepo)

	// 创建TemplateService（模板服务）
	var templateSvc *documentService.TemplateService
	templateRepo := repositoryFactory.CreateTemplateRepository()
	if templateRepo != nil {
		// 创建一个简单的logger
		logger, _ := zap.NewDevelopment()
		templateSvc = documentService.NewTemplateService(templateRepo, logger)
	}

	// 调用InitWriterRouter初始化所有写作路由
	InitWriterRouter(r, projectSvc, documentSvc, versionSvc, searchSvc, exportSvc, publishSvc, lockSvc, commentSvc, templateSvc)
}
