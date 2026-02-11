package impl

import (
	"context"

	"Qingyu_backend/models/reader"
	serviceReader "Qingyu_backend/service/interfaces/reader"
	readerService "Qingyu_backend/service/reader"
)

// ReadingProgressImpl 阅读进度管理端口实现
type ReadingProgressImpl struct {
	readerService *readerService.ReaderService
	serviceName   string
	version       string
}

// NewReadingProgressImpl 创建阅读进度管理端口实现
func NewReadingProgressImpl(readerService *readerService.ReaderService) serviceReader.ReadingProgressPort {
	return &ReadingProgressImpl{
		readerService: readerService,
		serviceName:   "ReadingProgressPort",
		version:       "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (r *ReadingProgressImpl) Initialize(ctx context.Context) error {
	return r.readerService.Initialize(ctx)
}

func (r *ReadingProgressImpl) Health(ctx context.Context) error {
	return r.readerService.Health(ctx)
}

func (r *ReadingProgressImpl) Close(ctx context.Context) error {
	return r.readerService.Close(ctx)
}

func (r *ReadingProgressImpl) GetServiceName() string {
	return r.serviceName
}

func (r *ReadingProgressImpl) GetVersion() string {
	return r.version
}

// ============================================================================
// ReadingProgressPort 方法实现
// ============================================================================

// GetReadingProgress 获取阅读进度
func (r *ReadingProgressImpl) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	return r.readerService.GetReadingProgress(ctx, userID, bookID)
}

// SaveReadingProgress 保存阅读进度
func (r *ReadingProgressImpl) SaveReadingProgress(ctx context.Context, req *serviceReader.SaveReadingProgressRequest) error {
	// 委托给现有 Service
	return r.readerService.SaveReadingProgress(ctx, req.UserID, req.BookID, req.ChapterID, req.Progress)
}

// UpdateReadingTime 更新阅读时长
func (r *ReadingProgressImpl) UpdateReadingTime(ctx context.Context, req *serviceReader.UpdateReadingTimeRequest) error {
	// 委托给现有 Service
	return r.readerService.UpdateReadingTime(ctx, req.UserID, req.BookID, req.Duration)
}

// GetRecentReading 获取最近阅读记录
func (r *ReadingProgressImpl) GetRecentReading(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error) {
	return r.readerService.GetRecentReading(ctx, userID, limit)
}

// GetReadingHistory 获取阅读历史
func (r *ReadingProgressImpl) GetReadingHistory(ctx context.Context, req *serviceReader.GetReadingHistoryRequest) (*serviceReader.GetReadingHistoryResponse, error) {
	// 委托给现有 Service，并进行类型转换
	progresses, total, err := r.readerService.GetReadingHistory(ctx, req.UserID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / req.Size
	if int(total)%req.Size > 0 {
		totalPages++
	}

	return &serviceReader.GetReadingHistoryResponse{
		Progresses: progresses,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
	}, nil
}

// GetTotalReadingTime 获取总阅读时长
func (r *ReadingProgressImpl) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	return r.readerService.GetTotalReadingTime(ctx, userID)
}

// GetReadingTimeByPeriod 获取时间段内的阅读时长
func (r *ReadingProgressImpl) GetReadingTimeByPeriod(ctx context.Context, req *serviceReader.GetReadingTimeByPeriodRequest) (int64, error) {
	// 委托给现有 Service
	return r.readerService.GetReadingTimeByPeriod(ctx, req.UserID, req.StartTime, req.EndTime)
}

// GetUnfinishedBooks 获取未读完的书籍
func (r *ReadingProgressImpl) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	return r.readerService.GetUnfinishedBooks(ctx, userID)
}

// GetFinishedBooks 获取已读完的书籍
func (r *ReadingProgressImpl) GetFinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	return r.readerService.GetFinishedBooks(ctx, userID)
}

// DeleteReadingProgress 删除阅读进度
func (r *ReadingProgressImpl) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	return r.readerService.DeleteReadingProgress(ctx, userID, bookID)
}

// UpdateBookStatus 更新书籍状态（在读/想读/读完）
func (r *ReadingProgressImpl) UpdateBookStatus(ctx context.Context, req *serviceReader.UpdateBookStatusRequest) error {
	// 委托给现有 Service
	return r.readerService.UpdateBookStatus(ctx, req.UserID, req.BookID, req.Status)
}

// BatchUpdateBookStatus 批量更新书籍状态
func (r *ReadingProgressImpl) BatchUpdateBookStatus(ctx context.Context, req *serviceReader.BatchUpdateBookStatusRequest) error {
	// 委托给现有 Service
	return r.readerService.BatchUpdateBookStatus(ctx, req.UserID, req.BookIDs, req.Status)
}
