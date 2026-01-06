package writer

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// PublishService 发布管理服务实现
type PublishService struct {
	projectRepo     ProjectRepository
	documentRepo    DocumentRepository
	publicationRepo PublicationRepository
	bookstoreClient BookstoreClient // 书城客户端接口
	eventBus        EventBus
}

// PublicationRepository 发布记录仓储接口
type PublicationRepository interface {
	Create(ctx context.Context, record *serviceInterfaces.PublicationRecord) error
	FindByID(ctx context.Context, id string) (*serviceInterfaces.PublicationRecord, error)
	FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error)
	FindByResourceID(ctx context.Context, resourceID string) (*serviceInterfaces.PublicationRecord, error)
	Update(ctx context.Context, record *serviceInterfaces.PublicationRecord) error
	Delete(ctx context.Context, id string) error
	FindPublishedByProjectID(ctx context.Context, projectID string) (*serviceInterfaces.PublicationRecord, error)
}

// BookstoreClient 书城客户端接口
type BookstoreClient interface {
	PublishProject(ctx context.Context, req *BookstorePublishProjectRequest) (*BookstorePublishResponse, error)
	UnpublishProject(ctx context.Context, projectID, bookstoreID string) error
	PublishChapter(ctx context.Context, req *BookstorePublishChapterRequest) (*BookstorePublishResponse, error)
	UnpublishChapter(ctx context.Context, chapterID, bookstoreID string) error
	UpdateChapter(ctx context.Context, req *BookstoreUpdateChapterRequest) error
	GetStatistics(ctx context.Context, projectID, bookstoreID string) (*serviceInterfaces.PublicationStatistics, error)
}

// EventBus 事件总线接口
type EventBus interface {
	PublishAsync(ctx context.Context, event interface{}) error
}

// BookstorePublishProjectRequest 书城发布项目请求
type BookstorePublishProjectRequest struct {
	ProjectID      string
	ProjectTitle   string
	AuthorID       string
	CategoryID     string
	Tags           []string
	Description    string
	CoverImage     string
	PublishType    string
	Price          float64
	FreeChapters   int
	AuthorNote     string
	EnableComment  bool
	EnableShare    bool
}

// BookstorePublishChapterRequest 书城发布章节请求
type BookstorePublishChapterRequest struct {
	ProjectID     string
	DocumentID    string
	ChapterTitle  string
	ChapterNumber int
	Content       string
	IsFree        bool
	AuthorNote    string
}

// BookstoreUpdateChapterRequest 书城更新章节请求
type BookstoreUpdateChapterRequest struct {
	ChapterID     string
	ChapterTitle  string
	ChapterNumber int
	Content       string
	IsFree        bool
}

// BookstorePublishResponse 书城发布响应
type BookstorePublishResponse struct {
	Success       bool
	BookstoreID   string
	BookstoreName string
	ExternalID    string // 书城中的外部ID
	Message       string
}

// NewPublishService 创建PublishService实例
func NewPublishService(
	projectRepo ProjectRepository,
	documentRepo DocumentRepository,
	publicationRepo PublicationRepository,
	bookstoreClient BookstoreClient,
	eventBus EventBus,
) serviceInterfaces.PublishService {
	return &PublishService{
		projectRepo:     projectRepo,
		documentRepo:    documentRepo,
		publicationRepo: publicationRepo,
		bookstoreClient: bookstoreClient,
		eventBus:        eventBus,
	}
}

// PublishProject 发布项目到书城
func (s *PublishService) PublishProject(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.PublishProjectRequest,
) (*serviceInterfaces.PublicationRecord, error) {
	// 验证项目是否存在
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "项目不存在", "", err)
	}

	// 验证权限
	if !project.IsOwner(userID) {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权发布此项目", "", nil)
	}

	// 检查是否已发布
	existingRecord, _ := s.publicationRepo.FindPublishedByProjectID(ctx, projectID)
	if existingRecord != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorValidation, "项目已发布，请先取消发布", "", nil)
	}

	// 创建发布记录
	record := &serviceInterfaces.PublicationRecord{
		ID:            primitive.NewObjectID().Hex(),
		Type:          "project",
		ResourceID:    projectID,
		ResourceTitle: project.Title,
		BookstoreID:   req.BookstoreID,
		Status:        serviceInterfaces.PublicationStatusPending,
		ScheduledTime: req.PublishTime,
		Metadata: serviceInterfaces.PublicationMetadata{
			CategoryID:    req.CategoryID,
			Tags:          req.Tags,
			Description:   req.Description,
			CoverImage:    req.CoverImage,
			PublishType:   req.PublishType,
			FreeChapters:  req.FreeChapters,
			AuthorNote:    req.AuthorNote,
			EnableComment: req.EnableComment,
			EnableShare:   req.EnableShare,
		},
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Price != nil {
		record.Metadata.Price = *req.Price
	}

	// 保存发布记录
	if err := s.publicationRepo.Create(ctx, record); err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorInternal, "创建发布记录失败", "", err)
	}

	// 如果是立即发布，调用书城API
	if req.PublishTime == nil || req.PublishTime.Before(time.Now()) {
		go s.executeProjectPublish(context.Background(), record, project, req)
	} else {
		// 定时发布任务应该由定时任务调度器处理
		// 这里暂时返回记录
	}

	return record, nil
}

// executeProjectPublish 执行项目发布
func (s *PublishService) executeProjectPublish(
	ctx context.Context,
	record *serviceInterfaces.PublicationRecord,
	project *writer.Project,
	req *serviceInterfaces.PublishProjectRequest,
) {
	// 调用书城API发布
	bookstoreReq := &BookstorePublishProjectRequest{
		ProjectID:     project.ID,
		ProjectTitle:  project.Title,
		AuthorID:      project.AuthorID,
		CategoryID:    req.CategoryID,
		Tags:          req.Tags,
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		PublishType:   req.PublishType,
		FreeChapters:  req.FreeChapters,
		AuthorNote:    req.AuthorNote,
		EnableComment: req.EnableComment,
		EnableShare:   req.EnableShare,
	}

	if req.Price != nil {
		bookstoreReq.Price = *req.Price
	}

	resp, err := s.bookstoreClient.PublishProject(ctx, bookstoreReq)
	if err != nil {
		s.failPublication(ctx, record, "发布到书城失败: "+err.Error())
		return
	}

	if !resp.Success {
		s.failPublication(ctx, record, resp.Message)
		return
	}

	// 更新发布记录为成功
	record.Status = serviceInterfaces.PublicationStatusPublished
	record.BookstoreName = resp.BookstoreName
	record.PublishTime = &time.Time{}
	*record.PublishTime = time.Now()
	record.UpdatedAt = time.Now()

	_ = s.publicationRepo.Update(ctx, record)

	// 发布事件
	if s.eventBus != nil {
		event := map[string]interface{}{
			"eventType": "project.published",
			"projectId": project.ID,
			"bookstoreId": req.BookstoreID,
			"publishedAt": time.Now(),
		}
		_ = s.eventBus.PublishAsync(ctx, event)
	}
}

// failPublication 标记发布失败
func (s *PublishService) failPublication(ctx context.Context, record *serviceInterfaces.PublicationRecord, errorMsg string) {
	record.Status = serviceInterfaces.PublicationStatusFailed
	record.UpdatedAt = time.Now()
	_ = s.publicationRepo.Update(ctx, record)
}

// UnpublishProject 取消发布项目
func (s *PublishService) UnpublishProject(ctx context.Context, projectID, userID string) error {
	// 查找发布记录
	record, err := s.publicationRepo.FindPublishedByProjectID(ctx, projectID)
	if err != nil {
		return errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "项目未发布或发布记录不存在", "", err)
	}

	// 验证权限
	if record.CreatedBy != userID {
		return errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权取消此项目的发布", "", nil)
	}

	// 调用书城API取消发布
	if err := s.bookstoreClient.UnpublishProject(ctx, projectID, record.BookstoreID); err != nil {
		return errors.NewServiceError("PublishService", errors.ServiceErrorInternal, "取消发布失败", "", err)
	}

	// 更新发布记录
	now := time.Now()
	record.Status = serviceInterfaces.PublicationStatusUnpublished
	record.UnpublishTime = &now
	record.UpdatedAt = now
	if err := s.publicationRepo.Update(ctx, record); err != nil {
		return err
	}

	// 发布事件
	if s.eventBus != nil {
		event := map[string]interface{}{
			"eventType": "project.unpublished",
			"projectId": projectID,
			"unpublishedAt": now,
		}
		_ = s.eventBus.PublishAsync(ctx, event)
	}

	return nil
}

// GetProjectPublicationStatus 获取项目发布状态
func (s *PublishService) GetProjectPublicationStatus(ctx context.Context, projectID string) (*serviceInterfaces.PublicationStatus, error) {
	// 获取项目信息
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "项目不存在", "", err)
	}

	// 查找发布记录
	record, _ := s.publicationRepo.FindPublishedByProjectID(ctx, projectID)

	status := &serviceInterfaces.PublicationStatus{
		ProjectID:   projectID,
		ProjectTitle: project.Title,
		IsPublished: record != nil && record.Status == serviceInterfaces.PublicationStatusPublished,
	}

	if record != nil {
		status.BookstoreID = record.BookstoreID
		status.BookstoreName = record.BookstoreName
		status.PublishedAt = record.PublishTime
	}

	// 获取项目统计
	if record != nil {
		stats, _ := s.bookstoreClient.GetStatistics(ctx, projectID, record.BookstoreID)
		if stats != nil {
			status.Statistics = *stats
		}
	}

	// 获取文档统计
	documents, _ := s.documentRepo.FindByProjectID(ctx, projectID)
	status.TotalChapters = len(documents)
	status.PublishedChapters = len(documents) // 这里应该从实际发布记录统计

	return status, nil
}

// PublishDocument 发布文档（章节）
func (s *PublishService) PublishDocument(
	ctx context.Context,
	documentID, projectID, userID string,
	req *serviceInterfaces.PublishDocumentRequest,
) (*serviceInterfaces.PublicationRecord, error) {
	// 验证文档是否存在
	document, err := s.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "文档不存在", "", err)
	}

	// 验证项目权限
	if document.ProjectID != projectID {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权访问此文档", "", nil)
	}

	// 获取项目信息
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !project.IsOwner(userID) {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权发布此文档", "", nil)
	}

	// 检查是否已发布
	existingRecord, _ := s.publicationRepo.FindByResourceID(ctx, documentID)
	if existingRecord != nil && existingRecord.Status == serviceInterfaces.PublicationStatusPublished {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorValidation, "文档已发布", "", nil)
	}

	// 获取发布记录
	projectRecord, _ := s.publicationRepo.FindPublishedByProjectID(ctx, projectID)
	bookstoreID := ""
	bookstoreName := ""
	if projectRecord != nil {
		bookstoreID = projectRecord.BookstoreID
		bookstoreName = projectRecord.BookstoreName
	}

	// 创建发布记录
	record := &serviceInterfaces.PublicationRecord{
		ID:            primitive.NewObjectID().Hex(),
		Type:          "document",
		ResourceID:    documentID,
		ResourceTitle: req.ChapterTitle,
		BookstoreID:   bookstoreID,
		Status:        serviceInterfaces.PublicationStatusPending,
		ScheduledTime: req.PublishTime,
		Metadata: serviceInterfaces.PublicationMetadata{
			ChapterTitle:  req.ChapterTitle,
			ChapterNumber: req.ChapterNumber,
			IsFree:        req.IsFree,
			AuthorNote:    req.AuthorNote,
		},
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if bookstoreName != "" {
		record.BookstoreName = bookstoreName
	}

	// 保存发布记录
	if err := s.publicationRepo.Create(ctx, record); err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorInternal, "创建发布记录失败", "", err)
	}

	// 如果是立即发布
	if req.PublishTime == nil || req.PublishTime.Before(time.Now()) {
		go s.executeDocumentPublish(context.Background(), record, document, req)
	}

	return record, nil
}

// executeDocumentPublish 执行文档发布
func (s *PublishService) executeDocumentPublish(
	ctx context.Context,
	record *serviceInterfaces.PublicationRecord,
	document *writer.Document,
	req *serviceInterfaces.PublishDocumentRequest,
) {
	// 这里需要获取文档内容
	// content, err := s.documentContentRepo.FindByID(ctx, document.ID)
	// if err != nil {
	//     s.failPublication(ctx, record, "获取文档内容失败: "+err.Error())
	//     return
	// }

	// 调用书城API
	bookstoreReq := &BookstorePublishChapterRequest{
		ProjectID:     document.ProjectID,
		DocumentID:    document.ID,
		ChapterTitle:  req.ChapterTitle,
		ChapterNumber: req.ChapterNumber,
		IsFree:        req.IsFree,
		AuthorNote:    req.AuthorNote,
	}

	resp, err := s.bookstoreClient.PublishChapter(ctx, bookstoreReq)
	if err != nil {
		s.failPublication(ctx, record, "发布章节失败: "+err.Error())
		return
	}

	if !resp.Success {
		s.failPublication(ctx, record, resp.Message)
		return
	}

	// 更新发布记录
	record.Status = serviceInterfaces.PublicationStatusPublished
	if record.BookstoreName == "" {
		record.BookstoreName = resp.BookstoreName
	}
	record.PublishTime = &time.Time{}
	*record.PublishTime = time.Now()
	record.UpdatedAt = time.Now()

	_ = s.publicationRepo.Update(ctx, record)
}

// UpdateDocumentPublishStatus 更新文档发布状态
func (s *PublishService) UpdateDocumentPublishStatus(
	ctx context.Context,
	documentID, projectID, userID string,
	req *serviceInterfaces.UpdateDocumentPublishStatusRequest,
) error {
	// 验证权限
	document, err := s.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "文档不存在", "", err)
	}

	if document.ProjectID != projectID {
		return errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权访问此文档", "", nil)
	}

	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}

	if !project.IsOwner(userID) {
		return errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权更新此文档的发布状态", "", nil)
	}

	// 查找发布记录
	record, err := s.publicationRepo.FindByResourceID(ctx, documentID)
	if err != nil {
		return errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "发布记录不存在", "", err)
	}

	// 取消发布
	if !req.IsPublished {
		return s.unpublishDocument(ctx, record, req.UnpublishReason)
	}

	// 更新发布信息
	if req.IsFree {
		record.Metadata.IsFree = req.IsFree
	}
	if req.ChapterNumber != nil {
		record.Metadata.ChapterNumber = *req.ChapterNumber
	}
	record.UpdatedAt = time.Now()

	return s.publicationRepo.Update(ctx, record)
}

// unpublishDocument 取消文档发布
func (s *PublishService) unpublishDocument(ctx context.Context, record *serviceInterfaces.PublicationRecord, reason string) error {
	// 调用书城API
	if err := s.bookstoreClient.UnpublishChapter(ctx, record.ResourceID, record.BookstoreID); err != nil {
		return errors.NewServiceError("PublishService", errors.ServiceErrorInternal, "取消发布失败", "", err)
	}

	// 更新记录
	now := time.Now()
	record.Status = serviceInterfaces.PublicationStatusUnpublished
	record.UnpublishTime = &now
	record.UnpublishReason = reason
	record.UpdatedAt = now

	return s.publicationRepo.Update(ctx, record)
}

// BatchPublishDocuments 批量发布文档
func (s *PublishService) BatchPublishDocuments(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.BatchPublishDocumentsRequest,
) (*serviceInterfaces.BatchPublishResult, error) {
	// 验证权限
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "项目不存在", "", err)
	}

	if !project.IsOwner(userID) {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorForbidden, "无权批量发布此项目的文档", "", nil)
	}

	result := &serviceInterfaces.BatchPublishResult{
		Results: make([]serviceInterfaces.BatchPublishItem, 0, len(req.DocumentIDs)),
	}

	// 批量发布
	for i, docID := range req.DocumentIDs {
		chapterNumber := req.StartNumber + i
		if !req.AutoNumbering {
			chapterNumber = 0
		}

		publishReq := &serviceInterfaces.PublishDocumentRequest{
			ChapterTitle:  "Chapter " + string(rune(i+1)), // 实际应该从文档标题获取
			ChapterNumber: chapterNumber,
			IsFree:        req.IsFree,
			PublishTime:   req.PublishTime,
		}

		record, err := s.PublishDocument(ctx, docID, projectID, userID, publishReq)
		if err != nil {
			result.FailCount++
			result.Results = append(result.Results, serviceInterfaces.BatchPublishItem{
				DocumentID: docID,
				Success:    false,
				Error:      err.Error(),
			})
		} else {
			result.SuccessCount++
			result.Results = append(result.Results, serviceInterfaces.BatchPublishItem{
				DocumentID: docID,
				Success:    true,
				RecordID:   record.ID,
			})
		}
	}

	return result, nil
}

// GetPublicationRecords 获取发布记录列表
func (s *PublishService) GetPublicationRecords(
	ctx context.Context,
	projectID string,
	page, pageSize int,
) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return s.publicationRepo.FindByProjectID(ctx, projectID, page, pageSize)
}

// GetPublicationRecord 获取发布记录详情
func (s *PublishService) GetPublicationRecord(ctx context.Context, recordID string) (*serviceInterfaces.PublicationRecord, error) {
	record, err := s.publicationRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, errors.NewServiceError("PublishService", errors.ServiceErrorNotFound, "发布记录不存在", "", err)
	}
	return record, nil
}
