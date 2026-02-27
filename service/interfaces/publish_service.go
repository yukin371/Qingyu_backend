package interfaces

import (
	"Qingyu_backend/models/dto"
	"context"
)

// PublishService 发布管理服务接口
type PublishService interface {
	// 项目发布
	PublishProject(ctx context.Context, projectID, userID string, req *PublishProjectRequest) (*PublicationRecord, error)
	UnpublishProject(ctx context.Context, projectID, userID string) error
	GetProjectPublicationStatus(ctx context.Context, projectID string) (*PublicationStatus, error)

	// 文档（章节）发布
	PublishDocument(ctx context.Context, documentID, projectID, userID string, req *PublishDocumentRequest) (*PublicationRecord, error)
	UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *UpdateDocumentPublishStatusRequest) error
	BatchPublishDocuments(ctx context.Context, projectID, userID string, req *BatchPublishDocumentsRequest) (*BatchPublishResult, error)

	// 发布记录管理
	GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*PublicationRecord, int64, error)
	GetPublicationRecord(ctx context.Context, recordID string) (*PublicationRecord, error)
}

// PublishProjectRequest 发布项目请求
// Deprecated: 使用 dto.PublishProjectRequest 替代
type PublishProjectRequest = dto.PublishProjectRequest

// PublishDocumentRequest 发布文档（章节）请求
// Deprecated: 使用 dto.PublishDocumentRequest 替代
type PublishDocumentRequest = dto.PublishDocumentRequest

// UpdateDocumentPublishStatusRequest 更新文档发布状态请求
// Deprecated: 使用 dto.UpdateDocumentPublishStatusRequest 替代
type UpdateDocumentPublishStatusRequest = dto.UpdateDocumentPublishStatusRequest

// BatchPublishDocumentsRequest 批量发布文档请求
// Deprecated: 使用 dto.BatchPublishDocumentsRequest 替代
type BatchPublishDocumentsRequest = dto.BatchPublishDocumentsRequest

// BatchPublishResult 批量发布结果
// Deprecated: 使用 dto.BatchPublishResult 替代
type BatchPublishResult = dto.BatchPublishResult

// BatchPublishItem 批量发布项
// Deprecated: 使用 dto.BatchPublishItem 替代
type BatchPublishItem = dto.BatchPublishItem

// PublicationRecord 发布记录
// Deprecated: 使用 dto.PublicationRecord 替代
type PublicationRecord = dto.PublicationRecord

// PublicationStatus 发布状态
// Deprecated: 使用 dto.PublicationStatus 替代
type PublicationStatus = dto.PublicationStatus

// PublicationStatus 发布状态常量
const (
	PublicationStatusPending     = dto.PublicationStatusPending     // 等待发布
	PublicationStatusPublished   = dto.PublicationStatusPublished   // 已发布
	PublicationStatusUnpublished = dto.PublicationStatusUnpublished // 已取消发布
	PublicationStatusFailed      = dto.PublicationStatusFailed      // 发布失败
)

// PublishType 发布类型常量
const (
	PublishTypeSerial   = dto.PublishTypeSerial   // 连载
	PublishTypeComplete = dto.PublishTypeComplete // 完本
)
