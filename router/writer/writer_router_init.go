package writer

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/models/dto"
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
	// 注意：由于PublicationRepository和BookstoreClient尚未完全实现，暂时使用MockPublishService
	// 当书城服务完善后，可以替换为真实的PublishService
	var publishSvc interfaces.PublishService = NewMockPublishService()

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

	// 创建CharacterService（角色服务）
	characterRepo := repositoryFactory.CreateCharacterRepository()
	var characterSvc interfaces.CharacterService
	if characterRepo != nil {
		characterSvc = writerservice.NewCharacterService(characterRepo, eventBus)
	}

	// 创建LocationService（地点服务）
	locationRepo := repositoryFactory.CreateLocationRepository()
	var locationSvc interfaces.LocationService
	if locationRepo != nil {
		locationSvc = writerservice.NewLocationService(locationRepo, eventBus)
	}

	// 创建TimelineService（时间线服务）
	timelineRepo := repositoryFactory.CreateTimelineRepository()
	timelineEventRepo := repositoryFactory.CreateTimelineEventRepository()
	var timelineSvc interfaces.TimelineService
	if timelineRepo != nil && timelineEventRepo != nil {
		timelineSvc = writerservice.NewTimelineService(timelineRepo, timelineEventRepo, eventBus)
	}

	// 创建OutlineService（大纲服务）
	outlineRepo := repositoryFactory.CreateOutlineRepository()
	var outlineSvc interfaces.OutlineService
	if outlineRepo != nil {
		outlineSvc = writerservice.NewOutlineService(outlineRepo, eventBus)
	}

	// 调用InitWriterRouter初始化文档编辑相关路由
	InitWriterRouter(r, projectSvc, documentSvc, versionSvc, searchSvc, exportSvc, publishSvc, lockSvc, commentSvc, templateSvc)

	// 调用InitWriterRoutes初始化设定百科路由（角色、地点、时间线、大纲）
	InitWriterRoutes(r, characterSvc, locationSvc, timelineSvc, outlineSvc)
}

// MockPublishService 用于E2E测试的Mock发布服务
type MockPublishService struct{}

// NewMockPublishService 创建Mock发布服务实例
func NewMockPublishService() *MockPublishService {
	return &MockPublishService{}
}

// PublishProject 模拟发布项目
func (m *MockPublishService) PublishProject(ctx context.Context, projectID, userID string, req *interfaces.PublishProjectRequest) (*interfaces.PublicationRecord, error) {
	// 返回一个模拟的成功发布记录
	now := time.Now()
	return &interfaces.PublicationRecord{
		ID:            "mock_pub_" + projectID,
		Type:          "project",
		ResourceID:    projectID,
		ResourceTitle: "Mock Published Project",
		BookstoreID:   req.BookstoreID,
		BookstoreName: "Mock Bookstore",
		Status:        interfaces.PublicationStatusPublished,
		PublishTime:   &now,
		CreatedBy:     userID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// UnpublishProject 模拟取消发布项目
func (m *MockPublishService) UnpublishProject(ctx context.Context, projectID, userID string) error {
	return nil
}

// GetProjectPublicationStatus 获取项目发布状态
func (m *MockPublishService) GetProjectPublicationStatus(ctx context.Context, projectID string) (*interfaces.PublicationStatus, error) {
	return &interfaces.PublicationStatus{
		ProjectID:           projectID,
		ProjectTitle:        "Mock Project",
		IsPublished:         false,
		TotalChapters:       0,
		PublishedChapters:   0,
		UnpublishedChapters: 0,
		PendingChapters:     0,
	}, nil
}

// PublishDocument 模拟发布文档
func (m *MockPublishService) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *interfaces.PublishDocumentRequest) (*interfaces.PublicationRecord, error) {
	// 如果没有提供 chapterTitle，使用默认值
	chapterTitle := req.ChapterTitle
	if chapterTitle == "" {
		chapterTitle = "Mock Chapter"
	}

	now := time.Now()
	return &interfaces.PublicationRecord{
		ID:            "mock_pub_doc_" + documentID,
		Type:          "document",
		ResourceID:    documentID,
		ResourceTitle: chapterTitle,
		BookstoreID:   "mock_bookstore",
		BookstoreName: "Mock Bookstore",
		Status:        interfaces.PublicationStatusPublished,
		PublishTime:   &now,
		Metadata: dto.PublicationMetadata{
			ChapterTitle:  chapterTitle,
			ChapterNumber: req.ChapterNumber,
			IsFree:        req.IsFree,
		},
		CreatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateDocumentPublishStatus 更新文档发布状态
func (m *MockPublishService) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *interfaces.UpdateDocumentPublishStatusRequest) error {
	return nil
}

// BatchPublishDocuments 批量发布文档
func (m *MockPublishService) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *interfaces.BatchPublishDocumentsRequest) (*interfaces.BatchPublishResult, error) {
	result := &interfaces.BatchPublishResult{
		SuccessCount: len(req.DocumentIDs),
		FailCount:    0,
		Results:      make([]interfaces.BatchPublishItem, 0, len(req.DocumentIDs)),
	}

	for _, docID := range req.DocumentIDs {
		result.Results = append(result.Results, interfaces.BatchPublishItem{
			DocumentID: docID,
			Success:    true,
			RecordID:   "mock_pub_" + docID,
		})
	}

	return result, nil
}

// GetPublicationRecords 获取发布记录列表
func (m *MockPublishService) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*interfaces.PublicationRecord, int64, error) {
	return []*interfaces.PublicationRecord{}, 0, nil
}

// GetPublicationRecord 获取发布记录详情
func (m *MockPublishService) GetPublicationRecord(ctx context.Context, recordID string) (*interfaces.PublicationRecord, error) {
	return nil, nil
}
