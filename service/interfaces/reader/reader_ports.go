package reader

import (
	"context"
	readerModel "Qingyu_backend/models/reader"
	"time"
)

// ReadingProgressPort 阅读进度管理端口
// 负责阅读进度的保存、查询、统计和历史记录管理
type ReadingProgressPort interface {
	BaseService

	// GetReadingProgress 获取阅读进度
	GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error)

	// SaveReadingProgress 保存阅读进度
	SaveReadingProgress(ctx context.Context, req *SaveReadingProgressRequest) error

	// UpdateReadingTime 更新阅读时长
	UpdateReadingTime(ctx context.Context, req *UpdateReadingTimeRequest) error

	// GetRecentReading 获取最近阅读记录
	GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModel.ReadingProgress, error)

	// GetReadingHistory 获取阅读历史
	GetReadingHistory(ctx context.Context, req *GetReadingHistoryRequest) (*GetReadingHistoryResponse, error)

	// GetTotalReadingTime 获取总阅读时长
	GetTotalReadingTime(ctx context.Context, userID string) (int64, error)

	// GetReadingTimeByPeriod 获取时间段内的阅读时长
	GetReadingTimeByPeriod(ctx context.Context, req *GetReadingTimeByPeriodRequest) (int64, error)

	// GetUnfinishedBooks 获取未读完的书籍
	GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error)

	// GetFinishedBooks 获取已读完的书籍
	GetFinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error)

	// DeleteReadingProgress 删除阅读进度
	DeleteReadingProgress(ctx context.Context, userID, bookID string) error

	// UpdateBookStatus 更新书籍状态（在读/想读/读完）
	UpdateBookStatus(ctx context.Context, req *UpdateBookStatusRequest) error

	// BatchUpdateBookStatus 批量更新书籍状态
	BatchUpdateBookStatus(ctx context.Context, req *BatchUpdateBookStatusRequest) error
}

// AnnotationPort 标注与笔记管理端口
// 负责书签、高亮、笔记等标注功能的 CRUD 操作
type AnnotationPort interface {
	BaseService

	// CreateAnnotation 创建标注
	CreateAnnotation(ctx context.Context, annotation *readerModel.Annotation) error

	// UpdateAnnotation 更新标注
	UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error

	// DeleteAnnotation 删除标注
	DeleteAnnotation(ctx context.Context, annotationID string) error

	// GetAnnotationsByChapter 获取章节的标注
	GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModel.Annotation, error)

	// GetAnnotationsByBook 获取书籍的所有标注
	GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)

	// GetNotes 获取笔记
	GetNotes(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)

	// SearchNotes 搜索笔记
	SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModel.Annotation, error)

	// GetBookmarks 获取书签
	GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)

	// GetLatestBookmark 获取最新的书签
	GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModel.Annotation, error)

	// GetHighlights 获取高亮
	GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)

	// GetRecentAnnotations 获取最近的标注
	GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModel.Annotation, error)

	// GetPublicAnnotations 获取公开的标注
	GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModel.Annotation, error)

	// GetAnnotationStats 获取标注统计
	GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error)

	// BatchCreateAnnotations 批量创建标注
	BatchCreateAnnotations(ctx context.Context, annotations []*readerModel.Annotation) error

	// BatchDeleteAnnotations 批量删除标注
	BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error
}

// ChapterContentPort 章节内容获取端口
// 负责章节内容的获取和导航功能
type ChapterContentPort interface {
	BaseService

	// GetChapterContent 获取章节内容（基本）
	GetChapterContent(ctx context.Context, userID, chapterID string) (string, error)

	// GetChapterByID 获取章节信息（不含内容）
	GetChapterByID(ctx context.Context, chapterID string) (interface{}, error)

	// GetBookChapters 获取书籍的章节列表
	GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error)

	// GetChapterContentWithProgress 获取章节内容（带权限检查、进度保存）
	GetChapterContentWithProgress(ctx context.Context, req *GetChapterContentRequest) (*ChapterContentResponse, error)

	// GetChapterByNumber 根据章节号获取内容
	GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*ChapterContentResponse, error)

	// GetNextChapter 获取下一章
	GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error)

	// GetPreviousChapter 获取上一章
	GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error)

	// GetChapterList 获取章节目录
	GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*ChapterListResponse, error)

	// GetChapterInfo 获取章节信息（不含内容）
	GetChapterInfo(ctx context.Context, userID, chapterID string) (*ChapterInfo, error)
}

// ReaderSettingsPort 阅读设置管理端口
// 负责用户阅读偏好设置的 CRUD 操作
type ReaderSettingsPort interface {
	BaseService

	// GetReadingSettings 获取阅读设置
	GetReadingSettings(ctx context.Context, userID string) (*readerModel.ReadingSettings, error)

	// SaveReadingSettings 保存阅读设置
	SaveReadingSettings(ctx context.Context, settings *readerModel.ReadingSettings) error

	// UpdateReadingSettings 更新阅读设置
	UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error
}

// ReaderSyncPort 阅读数据同步端口
// 负责多端阅读数据的同步功能
type ReaderSyncPort interface {
	BaseService

	// SyncAnnotations 同步标注（多端同步）
	SyncAnnotations(ctx context.Context, req *SyncAnnotationsRequest) (*SyncAnnotationsResponse, error)
}

// ============================================================================
// 请求和响应结构体定义
// ============================================================================

// SaveReadingProgressRequest 保存阅读进度请求
type SaveReadingProgressRequest struct {
	UserID    string  `json:"user_id" validate:"required"`
	BookID    string  `json:"book_id" validate:"required"`
	ChapterID string  `json:"chapter_id" validate:"required"`
	Progress  float64 `json:"progress" validate:"min=0,max=1"`
}

// UpdateReadingTimeRequest 更新阅读时长请求
type UpdateReadingTimeRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	BookID   string `json:"book_id" validate:"required"`
	Duration int64  `json:"duration" validate:"required,gt=0"`
}

// GetReadingHistoryRequest 获取阅读历史请求
type GetReadingHistoryRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}

// GetReadingHistoryResponse 获取阅读历史响应
type GetReadingHistoryResponse struct {
	Progresses []*readerModel.ReadingProgress `json:"progresses"`
	Total      int64                          `json:"total"`
	Page       int                            `json:"page"`
	Size       int                            `json:"size"`
	TotalPages int                            `json:"total_pages"`
}

// GetReadingTimeByPeriodRequest 获取时间段阅读时长请求
type GetReadingTimeByPeriodRequest struct {
	UserID    string    `json:"user_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

// UpdateBookStatusRequest 更新书籍状态请求
type UpdateBookStatusRequest struct {
	UserID string `json:"user_id" validate:"required"`
	BookID string `json:"book_id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=reading want_read finished"`
}

// BatchUpdateBookStatusRequest 批量更新书籍状态请求
type BatchUpdateBookStatusRequest struct {
	UserID  string   `json:"user_id" validate:"required"`
	BookIDs []string `json:"book_ids" validate:"required,min=1,max=50"`
	Status  string   `json:"status" validate:"required,oneof=reading want_read finished"`
}

// GetChapterContentRequest 获取章节内容请求
type GetChapterContentRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	BookID    string `json:"book_id" validate:"required"`
	ChapterID string `json:"chapter_id" validate:"required"`
}

// SyncAnnotationsRequest 同步标注请求
type SyncAnnotationsRequest struct {
	UserID           string                 `json:"user_id" validate:"required"`
	BookID           string                 `json:"book_id" validate:"required"`
	LastSyncTime     int64                  `json:"last_sync_time"`
	LocalAnnotations []*readerModel.Annotation `json:"local_annotations"`
}

// SyncAnnotationsResponse 同步标注响应
type SyncAnnotationsResponse struct {
	NewAnnotations  []*readerModel.Annotation `json:"new_annotations"`
	SyncTime        int64                     `json:"sync_time"`
	UploadedCount   int                       `json:"uploaded_count"`
	DownloadedCount int                       `json:"downloaded_count"`
}

// ============================================================================
// BaseService 基础接口（复用）
// ============================================================================

// BaseService 基础Service接口
type BaseService interface {
	// 服务生命周期
	Initialize(ctx context.Context) error
	Health(ctx context.Context) error
	Close(ctx context.Context) error

	// 服务信息
	GetServiceName() string
	GetVersion() string
}

// ============================================================================
// 章节相关类型定义
// ============================================================================

// ChapterContentResponse 章节内容响应
type ChapterContentResponse struct {
	ChapterID   string `json:"chapter_id"`
	BookID      string `json:"book_id"`
	ChapterNum  int    `json:"chapter_num"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	WordCount   int    `json:"word_count"`
	IsVIP       bool   `json:"is_vip"`
	PublishedAt int64  `json:"published_at"`
}

// ChapterInfo 章节信息（不含内容）
type ChapterInfo struct {
	ChapterID   string `json:"chapter_id"`
	BookID      string `json:"book_id"`
	ChapterNum  int    `json:"chapter_num"`
	Title       string `json:"title"`
	WordCount   int    `json:"word_count"`
	IsVIP       bool   `json:"is_vip"`
	PublishedAt int64  `json:"published_at"`
}

// ChapterListResponse 章节列表响应
type ChapterListResponse struct {
	Chapters []*ChapterInfo `json:"chapters"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	Size     int            `json:"size"`
}
