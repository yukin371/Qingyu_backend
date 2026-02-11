package impl

import (
	"context"
	"time"

	serviceInterfaces "Qingyu_backend/service/interfaces"
	serviceWriter "Qingyu_backend/service/interfaces/writer"
	writerservice "Qingyu_backend/service/writer"
)

// PublishExportImpl 发布导出端口实现
type PublishExportImpl struct {
	publishService writerservice.PublishService
	exportService  writerservice.ExportService
	serviceName    string
	version        string
}

// NewPublishExportImpl 创建发布导出端口实现
func NewPublishExportImpl(
	publishService writerservice.PublishService,
	exportService writerservice.ExportService,
) serviceWriter.PublishExportPort {
	return &PublishExportImpl{
		publishService: publishService,
		exportService:  exportService,
		serviceName:    "PublishExportPort",
		version:        "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (p *PublishExportImpl) Initialize(ctx context.Context) error {
	return nil
}

func (p *PublishExportImpl) Health(ctx context.Context) error {
	return nil
}

func (p *PublishExportImpl) Close(ctx context.Context) error {
	return nil
}

func (p *PublishExportImpl) GetServiceName() string {
	return p.serviceName
}

func (p *PublishExportImpl) GetVersion() string {
	return p.version
}

// ============================================================================
// PublishExportPort 发布管理方法
// ============================================================================

// PublishProject 发布项目到书城
func (p *PublishExportImpl) PublishProject(ctx context.Context, projectID, userID string, req *serviceWriter.PublishProjectRequest) (*serviceWriter.PublicationRecord, error) {
	// 转换请求类型
	publishReq := &serviceInterfaces.PublishProjectRequest{
		BookstoreID:   req.BookstoreID,
		CategoryID:    req.CategoryID,
		Tags:          req.Tags,
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		PublishType:   req.PublishType,
		PublishTime:   nil, // Port DTO 使用 *string，Service DTO 使用 *time.Time
		Price:         req.Price,
		FreeChapters:  req.FreeChapters,
		AuthorNote:    req.AuthorNote,
		EnableComment: req.EnableComment,
		EnableShare:   req.EnableShare,
	}
	if req.PublishTime != nil {
		t, err := time.Parse(time.RFC3339, *req.PublishTime)
		if err == nil {
			publishReq.PublishTime = &t
		}
	}
	record, err := p.publishService.PublishProject(ctx, projectID, userID, publishReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return p.convertPublicationRecord(record), nil
}

// UnpublishProject 取消发布项目
func (p *PublishExportImpl) UnpublishProject(ctx context.Context, projectID, userID string) error {
	return p.publishService.UnpublishProject(ctx, projectID, userID)
}

// GetProjectPublicationStatus 获取项目发布状态
func (p *PublishExportImpl) GetProjectPublicationStatus(ctx context.Context, projectID string) (*serviceWriter.PublicationStatus, error) {
	status, err := p.publishService.GetProjectPublicationStatus(ctx, projectID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return p.convertPublicationStatus(status), nil
}

// PublishDocument 发布文档（章节）
func (p *PublishExportImpl) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *serviceWriter.PublishDocumentRequest) (*serviceWriter.PublicationRecord, error) {
	// 转换请求类型
	publishReq := &serviceInterfaces.PublishDocumentRequest{
		ChapterTitle: req.ChapterTitle,
		ChapterNumber: req.ChapterNumber,
		IsFree:       req.IsFree,
		PublishTime:  nil, // Port DTO 使用 *string，Service DTO 使用 *time.Time
		AuthorNote:   req.AuthorNote,
	}
	if req.PublishTime != nil {
		t, err := time.Parse(time.RFC3339, *req.PublishTime)
		if err == nil {
			publishReq.PublishTime = &t
		}
	}
	record, err := p.publishService.PublishDocument(ctx, documentID, projectID, userID, publishReq)
	if err != nil {
		return nil, err
	}
	return p.convertPublicationRecord(record), nil
}

// UpdateDocumentPublishStatus 更新文档发布状态
func (p *PublishExportImpl) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *serviceWriter.UpdateDocumentPublishStatusRequest) error {
	// 转换请求类型
	updateReq := &serviceInterfaces.UpdateDocumentPublishStatusRequest{
		IsPublished:     req.IsPublished,
		PublishTime:     nil,
		UnpublishReason: req.UnpublishReason,
		IsFree:          req.IsFree,
		ChapterNumber:   req.ChapterNumber,
	}
	return p.publishService.UpdateDocumentPublishStatus(ctx, documentID, projectID, userID, updateReq)
}

// BatchPublishDocuments 批量发布文档
func (p *PublishExportImpl) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *serviceWriter.BatchPublishDocumentsRequest) (*serviceWriter.BatchPublishResult, error) {
	// 转换请求类型
	batchReq := &serviceInterfaces.BatchPublishDocumentsRequest{
		DocumentIDs:   req.DocumentIDs,
		AutoNumbering: req.AutoNumbering,
		StartNumber:   req.StartNumber,
		IsFree:        req.IsFree,
		PublishTime:   nil,
	}
	result, err := p.publishService.BatchPublishDocuments(ctx, projectID, userID, batchReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	items := make([]serviceWriter.BatchPublishItem, 0, len(result.Results))
	for _, item := range result.Results {
		items = append(items, serviceWriter.BatchPublishItem{
			DocumentID: item.DocumentID,
			Success:    item.Success,
			RecordID:   item.RecordID,
			Error:      item.Error,
		})
	}
	return &serviceWriter.BatchPublishResult{
		SuccessCount: result.SuccessCount,
		FailCount:    result.FailCount,
		Results:      items,
	}, nil
}

// GetPublicationRecords 获取发布记录列表
func (p *PublishExportImpl) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*serviceWriter.PublicationRecord, int64, error) {
	records, total, err := p.publishService.GetPublicationRecords(ctx, projectID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	// 转换响应类型
	result := make([]*serviceWriter.PublicationRecord, 0, len(records))
	for _, record := range records {
		result = append(result, p.convertPublicationRecord(record))
	}
	return result, total, nil
}

// GetPublicationRecord 获取发布记录详情
func (p *PublishExportImpl) GetPublicationRecord(ctx context.Context, recordID string) (*serviceWriter.PublicationRecord, error) {
	record, err := p.publishService.GetPublicationRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}
	return p.convertPublicationRecord(record), nil
}

// convertPublicationRecord 转换发布记录
func (p *PublishExportImpl) convertPublicationRecord(record *serviceInterfaces.PublicationRecord) *serviceWriter.PublicationRecord {
	r := &serviceWriter.PublicationRecord{
		ID:            record.ID,
		Type:          record.Type,
		ResourceID:    record.ResourceID,
		ResourceTitle: record.ResourceTitle,
		BookstoreID:   record.BookstoreID,
		BookstoreName: record.BookstoreName,
		Status:        record.Status,
		CreatedBy:     record.CreatedBy,
		CreatedAt:     record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if record.PublishTime != nil {
		formatted := record.PublishTime.Format("2006-01-02T15:04:05Z07:00")
		r.PublishTime = &formatted
	}
	return r
}

// convertPublicationStatus 转换发布状态
func (p *PublishExportImpl) convertPublicationStatus(status *serviceInterfaces.PublicationStatus) *serviceWriter.PublicationStatus {
	s := &serviceWriter.PublicationStatus{
		ProjectID:         status.ProjectID,
		ProjectTitle:      status.ProjectTitle,
		IsPublished:       status.IsPublished,
		BookstoreID:       status.BookstoreID,
		BookstoreName:     status.BookstoreName,
		TotalChapters:     status.TotalChapters,
		PublishedChapters: status.PublishedChapters,
	}
	if status.PublishedAt != nil {
		formatted := status.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		s.PublishedAt = &formatted
	}
	return s
}

// ============================================================================
// PublishExportPort 导出管理方法
// ============================================================================

// ExportDocument 导出文档
func (p *PublishExportImpl) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *serviceWriter.ExportDocumentRequest) (*serviceWriter.ExportTask, error) {
	// TODO: 实现类型转换
	return nil, nil
}

// ExportProject 导出项目
func (p *PublishExportImpl) ExportProject(ctx context.Context, projectID, userID string, req *serviceWriter.ExportProjectRequest) (*serviceWriter.ExportTask, error) {
	// TODO: 实现类型转换
	return nil, nil
}

// GetExportTask 获取导出任务
func (p *PublishExportImpl) GetExportTask(ctx context.Context, taskID string) (*serviceWriter.ExportTask, error) {
	// TODO: 实现类型转换
	return nil, nil
}

// DownloadExportFile 下载导出文件
func (p *PublishExportImpl) DownloadExportFile(ctx context.Context, taskID string) (*serviceWriter.ExportFile, error) {
	// TODO: 实现类型转换
	return nil, nil
}

// ListExportTasks 列出导出任务
func (p *PublishExportImpl) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*serviceWriter.ExportTask, int64, error) {
	// TODO: 实现类型转换
	return nil, 0, nil
}

// DeleteExportTask 删除导出任务
func (p *PublishExportImpl) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	// TODO: 实现类型转换
	return nil
}

// CancelExportTask 取消导出任务
func (p *PublishExportImpl) CancelExportTask(ctx context.Context, taskID, userID string) error {
	// TODO: 实现类型转换
	return nil
}
