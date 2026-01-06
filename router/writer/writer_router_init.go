package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/lock"
	"Qingyu_backend/service"
	documentService "Qingyu_backend/service/writer/document"
	projectService "Qingyu_backend/service/writer/project"
	searchService "Qingyu_backend/service/shared/search"
	writerservice "Qingyu_backend/service/writer"
	writerrepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/service/interfaces"
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

	// 创建SearchService（用于文档搜索）
	bookRepo := repositoryFactory.CreateBookRepository()
	searchSvc := searchService.NewSearchService(bookRepo, documentRepo)

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

	// 调用InitWriterRouter初始化所有写作路由
	InitWriterRouter(r, projectSvc, documentSvc, versionSvc, searchSvc, exportSvc, publishSvc, lockSvc, commentSvc)
}
