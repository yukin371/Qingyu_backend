package reader

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
	"Qingyu_backend/service/base"
)

// ReadingHistoryService 阅读历史Service
type ReadingHistoryService struct {
	historyRepo readerRepo.ReadingHistoryRepository
	eventBus    base.EventBus
}

// NewReadingHistoryService 创建阅读历史服务
func NewReadingHistoryService(
	historyRepo readerRepo.ReadingHistoryRepository,
	eventBus base.EventBus,
) *ReadingHistoryService {
	return &ReadingHistoryService{
		historyRepo: historyRepo,
		eventBus:    eventBus,
	}
}

// =======================
// BaseService接口实现
// =======================

// Initialize 初始化服务
func (s *ReadingHistoryService) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *ReadingHistoryService) Health(ctx context.Context) error {
	return s.historyRepo.Health(ctx)
}

// Close 关闭服务
func (s *ReadingHistoryService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *ReadingHistoryService) GetServiceName() string {
	return "ReadingHistoryService"
}

// GetVersion 获取服务版本
func (s *ReadingHistoryService) GetVersion() string {
	return "1.0.0"
}

// =======================
// 记录阅读行为
// =======================

// RecordReading 记录一次阅读行为
func (s *ReadingHistoryService) RecordReading(
	ctx context.Context,
	userID, bookID, chapterID string,
	startTime, endTime time.Time,
	progress float64,
	deviceType, deviceID string,
) error {
	// 参数验证
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return fmt.Errorf("书籍ID不能为空")
	}
	if chapterID == "" {
		return fmt.Errorf("章节ID不能为空")
	}

	// 验证时间合理性
	if endTime.Before(startTime) {
		return fmt.Errorf("结束时间不能早于开始时间")
	}

	// 验证时长不能超过24小时
	duration := int(endTime.Sub(startTime).Seconds())
	if duration > 86400 {
		return fmt.Errorf("阅读时长不合理")
	}

	// 验证进度范围
	if progress < 0 || progress > 100 {
		return fmt.Errorf("阅读进度必须在0-100之间")
	}

	// 创建历史记录
	history := &reader.ReadingHistory{
		UserID:       userID,
		BookID:       bookID,
		ChapterID:    chapterID,
		ReadDuration: duration,
		Progress:     progress,
		DeviceType:   deviceType,
		DeviceID:     deviceID,
		StartTime:    startTime,
		EndTime:      endTime,
		CreatedAt:    time.Now(),
	}

	// 保存到数据库
	if err := s.historyRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("记录阅读历史失败: %w", err)
	}

	// 发布阅读记录事件
	if s.eventBus != nil {
		event := &ReadingRecordedEvent{
			UserID:     userID,
			BookID:     bookID,
			ChapterID:  chapterID,
			Duration:   duration,
			Progress:   progress,
			RecordedAt: history.CreatedAt,
		}
		_ = s.eventBus.PublishAsync(ctx, event)
	}

	return nil
}

// =======================
// 查询历史记录
// =======================

// GetUserHistories 获取用户阅读历史（分页）
func (s *ReadingHistoryService) GetUserHistories(
	ctx context.Context,
	userID string,
	page, pageSize int,
) ([]*reader.ReadingHistory, *PaginationInfo, error) {
	// 参数验证
	if userID == "" {
		return nil, nil, fmt.Errorf("用户ID不能为空")
	}

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 计算offset
	offset := (page - 1) * pageSize

	// 查询历史列表
	histories, err := s.historyRepo.GetUserHistories(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}

	// 查询总数
	total, err := s.historyRepo.CountUserHistories(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("统计历史记录失败: %w", err)
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	pagination := &PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return histories, pagination, nil
}

// GetUserHistoriesByBook 获取用户指定书籍的阅读历史
func (s *ReadingHistoryService) GetUserHistoriesByBook(
	ctx context.Context,
	userID, bookID string,
	page, pageSize int,
) ([]*reader.ReadingHistory, *PaginationInfo, error) {
	// 参数验证
	if userID == "" {
		return nil, nil, fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return nil, nil, fmt.Errorf("书籍ID不能为空")
	}

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 计算offset
	offset := (page - 1) * pageSize

	// 查询历史列表
	histories, err := s.historyRepo.GetUserHistoriesByBook(ctx, userID, bookID, pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}

	// 注意：这里为了简化，不单独统计该书籍的记录数
	// 实际可以通过Repository添加CountUserHistoriesByBook方法

	pagination := &PaginationInfo{
		Page:     page,
		PageSize: pageSize,
		Total:    len(histories), // 简化处理
	}

	return histories, pagination, nil
}

// GetUserHistoriesByTimeRange 获取指定时间范围内的阅读历史
func (s *ReadingHistoryService) GetUserHistoriesByTimeRange(
	ctx context.Context,
	userID string,
	startTime, endTime time.Time,
	page, pageSize int,
) ([]*reader.ReadingHistory, *PaginationInfo, error) {
	// 参数验证
	if userID == "" {
		return nil, nil, fmt.Errorf("用户ID不能为空")
	}
	if endTime.Before(startTime) {
		return nil, nil, fmt.Errorf("结束时间不能早于开始时间")
	}

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 计算offset
	offset := (page - 1) * pageSize

	// 查询历史列表
	histories, err := s.historyRepo.GetUserHistoriesByTimeRange(ctx, userID, startTime, endTime, pageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}

	pagination := &PaginationInfo{
		Page:     page,
		PageSize: pageSize,
		Total:    len(histories),
	}

	return histories, pagination, nil
}

// =======================
// 阅读统计
// =======================

// GetUserReadingStats 获取用户阅读统计
func (s *ReadingHistoryService) GetUserReadingStats(
	ctx context.Context,
	userID string,
) (*readerRepo.ReadingStats, error) {
	// 参数验证
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 查询统计
	stats, err := s.historyRepo.GetUserReadingStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询阅读统计失败: %w", err)
	}

	return stats, nil
}

// GetUserDailyReadingStats 获取用户每日阅读统计
func (s *ReadingHistoryService) GetUserDailyReadingStats(
	ctx context.Context,
	userID string,
	days int,
) ([]readerRepo.DailyReadingStats, error) {
	// 参数验证
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 验证天数范围
	if days < 1 || days > 90 {
		days = 7 // 默认7天
	}

	// 查询统计
	stats, err := s.historyRepo.GetUserDailyReadingStats(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("查询每日阅读统计失败: %w", err)
	}

	return stats, nil
}

// =======================
// 删除操作
// =======================

// DeleteHistory 删除单条历史记录
func (s *ReadingHistoryService) DeleteHistory(
	ctx context.Context,
	userID, historyID string,
) error {
	// 参数验证
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if historyID == "" {
		return fmt.Errorf("历史记录ID不能为空")
	}

	// 先查询记录，验证权限
	history, err := s.historyRepo.GetByID(ctx, historyID)
	if err != nil {
		return fmt.Errorf("查询历史记录失败: %w", err)
	}

	// 验证是否为该用户的记录
	if history.UserID != userID {
		return fmt.Errorf("无权删除该历史记录")
	}

	// 删除记录
	if err := s.historyRepo.Delete(ctx, historyID); err != nil {
		return fmt.Errorf("删除历史记录失败: %w", err)
	}

	// 发布删除事件
	if s.eventBus != nil {
		event := &HistoryDeletedEvent{
			UserID:    userID,
			HistoryID: historyID,
			DeletedAt: time.Now(),
		}
		_ = s.eventBus.PublishAsync(ctx, event)
	}

	return nil
}

// ClearUserHistories 清空用户所有历史记录
func (s *ReadingHistoryService) ClearUserHistories(
	ctx context.Context,
	userID string,
) error {
	// 参数验证
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	// 删除所有记录
	if err := s.historyRepo.DeleteUserHistories(ctx, userID); err != nil {
		return fmt.Errorf("清空历史记录失败: %w", err)
	}

	// 发布清空事件
	if s.eventBus != nil {
		event := &HistoryClearedEvent{
			UserID:    userID,
			ClearedAt: time.Now(),
		}
		_ = s.eventBus.PublishAsync(ctx, event)
	}

	return nil
}

// =======================
// 自动清理
// =======================

// CleanOldHistories 清理过期历史记录
func (s *ReadingHistoryService) CleanOldHistories(
	ctx context.Context,
	daysToKeep int,
) (int64, error) {
	// 验证天数
	if daysToKeep < 1 {
		daysToKeep = 90 // 默认保留90天
	}

	// 计算截止日期
	beforeDate := time.Now().AddDate(0, 0, -daysToKeep)

	// 批量删除过期记录
	deleted, err := s.historyRepo.DeleteOldHistories(ctx, beforeDate)
	if err != nil {
		return 0, fmt.Errorf("清理过期历史记录失败: %w", err)
	}

	return deleted, nil
}

// =======================
// 事件定义
// =======================

// ReadingRecordedEvent 阅读记录事件
type ReadingRecordedEvent struct {
	UserID     string
	BookID     string
	ChapterID  string
	Duration   int
	Progress   float64
	RecordedAt time.Time
}

// GetEventType 获取事件类型
func (e *ReadingRecordedEvent) GetEventType() string {
	return "readerRepo.recorded"
}

// GetEventData 获取事件数据
func (e *ReadingRecordedEvent) GetEventData() interface{} {
	return e
}

// GetTimestamp 获取事件时间戳
func (e *ReadingRecordedEvent) GetTimestamp() time.Time {
	return e.RecordedAt
}

// GetSource 获取事件来源
func (e *ReadingRecordedEvent) GetSource() string {
	return "ReadingHistoryService"
}

// HistoryDeletedEvent 历史记录删除事件
type HistoryDeletedEvent struct {
	UserID    string
	HistoryID string
	DeletedAt time.Time
}

// GetEventType 获取事件类型
func (e *HistoryDeletedEvent) GetEventType() string {
	return "readerRepo.history.deleted"
}

// GetEventData 获取事件数据
func (e *HistoryDeletedEvent) GetEventData() interface{} {
	return e
}

// GetTimestamp 获取事件时间戳
func (e *HistoryDeletedEvent) GetTimestamp() time.Time {
	return e.DeletedAt
}

// GetSource 获取事件来源
func (e *HistoryDeletedEvent) GetSource() string {
	return "ReadingHistoryService"
}

// HistoryClearedEvent 历史记录清空事件
type HistoryClearedEvent struct {
	UserID    string
	ClearedAt time.Time
}

// GetEventType 获取事件类型
func (e *HistoryClearedEvent) GetEventType() string {
	return "readerRepo.history.cleared"
}

// GetEventData 获取事件数据
func (e *HistoryClearedEvent) GetEventData() interface{} {
	return e
}

// GetTimestamp 获取事件时间戳
func (e *HistoryClearedEvent) GetTimestamp() time.Time {
	return e.ClearedAt
}

// GetSource 获取事件来源
func (e *HistoryClearedEvent) GetSource() string {
	return "ReadingHistoryService"
}

// =======================
// 辅助结构
// =======================

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages,omitempty"`
}
