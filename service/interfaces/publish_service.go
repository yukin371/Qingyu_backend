package interfaces

import (
	"context"
	"time"
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
type PublishProjectRequest struct {
	BookstoreID   string     `json:"bookstoreId" validate:"required"`                       // 书城ID
	CategoryID    string     `json:"categoryId" validate:"required"`                        // 分类ID
	Tags          []string   `json:"tags" validate:"max=10"`                                // 标签
	Description   string     `json:"description" validate:"max=1000"`                       // 作品简介
	CoverImage    string     `json:"coverImage" validate:"omitempty,url"`                   // 封面图
	PublishType   string     `json:"publishType" validate:"required,oneof=serial complete"` // 发布方式：serial(连载) complete(完本)
	PublishTime   *time.Time `json:"publishTime,omitempty"`                                 // 定时发布时间（空则立即发布）
	Price         *float64   `json:"price,omitempty"`                                       // 定价（分）
	FreeChapters  int        `json:"freeChapters"`                                          // 免费章节数量
	AuthorNote    string     `json:"authorNote" validate:"max=500"`                         // 作者的话
	EnableComment bool       `json:"enableComment"`                                         // 是否开启评论
	EnableShare   bool       `json:"enableShare"`                                           // 是否允许分享
}

// PublishDocumentRequest 发布文档（章节）请求
type PublishDocumentRequest struct {
	ChapterTitle  string     `json:"chapterTitle" validate:"required"`        // 章节标题
	ChapterNumber int        `json:"chapterNumber" validate:"required,min=1"` // 章节序号
	IsFree        bool       `json:"isFree"`                                  // 是否免费
	PublishTime   *time.Time `json:"publishTime,omitempty"`                   // 定时发布时间（空则立即发布）
	AuthorNote    string     `json:"authorNote" validate:"max=500"`           // 章节作者的话
}

// UpdateDocumentPublishStatusRequest 更新文档发布状态请求
type UpdateDocumentPublishStatusRequest struct {
	IsPublished     bool       `json:"isPublished"`               // 是否已发布
	PublishTime     *time.Time `json:"publishTime,omitempty"`     // 发布时间
	UnpublishReason string     `json:"unpublishReason,omitempty"` // 取消发布原因
	IsFree          bool       `json:"isFree"`                    // 是否免费
	ChapterNumber   *int       `json:"chapterNumber,omitempty"`   // 章节序号
}

// BatchPublishDocumentsRequest 批量发布文档请求
type BatchPublishDocumentsRequest struct {
	DocumentIDs   []string   `json:"documentIds" validate:"required,min=1,max=50"` // 文档ID列表
	AutoNumbering bool       `json:"autoNumbering"`                                // 是否自动编号（从起始章节号开始）
	StartNumber   int        `json:"startNumber" validate:"min=1"`                 // 起始章节号
	IsFree        bool       `json:"isFree"`                                       // 是否全部免费
	PublishTime   *time.Time `json:"publishTime,omitempty"`                        // 发布时间
}

// BatchPublishResult 批量发布结果
type BatchPublishResult struct {
	SuccessCount int                `json:"successCount"`
	FailCount    int                `json:"failCount"`
	Results      []BatchPublishItem `json:"results"`
}

// BatchPublishItem 批量发布项
type BatchPublishItem struct {
	DocumentID string `json:"documentId"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	RecordID   string `json:"recordId,omitempty"`
}

// PublicationRecord 发布记录
type PublicationRecord struct {
	ID              string              `json:"id"`
	Type            string              `json:"type"`       // project, document
	ResourceID      string              `json:"resourceId"` // 项目ID或文档ID
	ResourceTitle   string              `json:"resourceTitle"`
	BookstoreID     string              `json:"bookstoreId"`
	BookstoreName   string              `json:"bookstoreName"`
	Status          string              `json:"status"`                    // pending, published, unpublished, failed
	PublishTime     *time.Time          `json:"publishTime"`               // 实际发布时间
	ScheduledTime   *time.Time          `json:"scheduledTime"`             // 计划发布时间
	UnpublishTime   *time.Time          `json:"unpublishTime,omitempty"`   // 取消发布时间
	UnpublishReason string              `json:"unpublishReason,omitempty"` // 取消发布原因
	Metadata        PublicationMetadata `json:"metadata"`                  // 发布元数据
	CreatedBy       string              `json:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

// PublicationMetadata 发布元数据
type PublicationMetadata struct {
	// 项目发布元数据
	CategoryID    string   `json:"categoryId,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Description   string   `json:"description,omitempty"`
	CoverImage    string   `json:"coverImage,omitempty"`
	PublishType   string   `json:"publishType,omitempty"`
	Price         float64  `json:"price,omitempty"`
	FreeChapters  int      `json:"freeChapters,omitempty"`
	AuthorNote    string   `json:"authorNote,omitempty"`
	EnableComment bool     `json:"enableComment,omitempty"`
	EnableShare   bool     `json:"enableShare,omitempty"`

	// 文档发布元数据
	ChapterTitle  string `json:"chapterTitle,omitempty"`
	ChapterNumber int    `json:"chapterNumber,omitempty"`
	IsFree        bool   `json:"isFree,omitempty"`
}

// PublicationStatus 发布状态
type PublicationStatus struct {
	ProjectID           string                `json:"projectId"`
	ProjectTitle        string                `json:"projectTitle"`
	IsPublished         bool                  `json:"isPublished"`
	BookstoreID         string                `json:"bookstoreId,omitempty"`
	BookstoreName       string                `json:"bookstoreName,omitempty"`
	PublishedAt         *time.Time            `json:"publishedAt,omitempty"`
	TotalChapters       int                   `json:"totalChapters"`
	PublishedChapters   int                   `json:"publishedChapters"`
	UnpublishedChapters int                   `json:"unpublishedChapters"`
	PendingChapters     int                   `json:"pendingChapters"` // 待发布章节
	Statistics          PublicationStatistics `json:"statistics"`
}

// PublicationStatistics 发布统计
type PublicationStatistics struct {
	TotalViews     int64     `json:"totalViews"`     // 总阅读量
	TotalLikes     int64     `json:"totalLikes"`     // 总点赞数
	TotalComments  int64     `json:"totalComments"`  // 总评论数
	TotalShares    int64     `json:"totalShares"`    // 总分享数
	TotalFavorites int64     `json:"totalFavorites"` // 总收藏数
	LastSyncTime   time.Time `json:"lastSyncTime"`   // 最后同步时间
}

// PublicationStatus 发布状态常量
const (
	PublicationStatusPending     = "pending"     // 等待发布
	PublicationStatusPublished   = "published"   // 已发布
	PublicationStatusUnpublished = "unpublished" // 已取消发布
	PublicationStatusFailed      = "failed"      // 发布失败
)

// PublishType 发布类型常量
const (
	PublishTypeSerial   = "serial"   // 连载
	PublishTypeComplete = "complete" // 完本
)
