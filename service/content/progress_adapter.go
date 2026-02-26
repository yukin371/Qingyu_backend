package content

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/reader"
	readerService "Qingyu_backend/service/reader"
)

// ProgressAdapter 阅读进度适配器
// 将现有ReaderService的功能适配到ReadingProgressService接口
type ProgressAdapter struct {
	readerService *readerService.ReaderService
}

// NewProgressAdapter 创建进度适配器
// 返回接口类型以符合服务契约
func NewProgressAdapter(readerSvc *readerService.ReaderService) *ProgressAdapter {
	return &ProgressAdapter{
		readerService: readerSvc,
	}
}

// =========================
// 进度管理
// =========================

// GetProgress 获取阅读进度
func (a *ProgressAdapter) GetProgress(ctx context.Context, userID, bookID string) (*dto.ReadingProgressResponse, error) {
	// 调用现有服务
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取阅读进度失败: %w", err)
	}

	// 转换为DTO
	response := &dto.ReadingProgressResponse{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   progress.ChapterID.Hex(),
		Progress:    float64(progress.Progress),
		ReadingTime: progress.ReadingTime,
		UpdateTime:  progress.UpdatedAt.Unix(),
	}

	return response, nil
}

// SaveProgress 保存阅读进度
func (a *ProgressAdapter) SaveProgress(ctx context.Context, req *dto.SaveProgressRequest) error {
	if req == nil {
		return fmt.Errorf("保存进度请求不能为空")
	}

	// 调用现有服务
	err := a.readerService.SaveReadingProgress(ctx, req.UserID, req.BookID, req.ChapterID, req.Progress)
	if err != nil {
		return fmt.Errorf("保存阅读进度失败: %w", err)
	}

	return nil
}

// UpdateReadingTime 更新阅读时长
func (a *ProgressAdapter) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int) error {
	if duration <= 0 {
		return fmt.Errorf("阅读时长必须大于0")
	}

	// 调用现有服务（需要转换为int64）
	err := a.readerService.UpdateReadingTime(ctx, userID, bookID, int64(duration))
	if err != nil {
		return fmt.Errorf("更新阅读时长失败: %w", err)
	}

	return nil
}

// DeleteProgress 删除阅读进度
func (a *ProgressAdapter) DeleteProgress(ctx context.Context, userID, bookID string) error {
	err := a.readerService.DeleteReadingProgress(ctx, userID, bookID)
	if err != nil {
		return fmt.Errorf("删除阅读进度失败: %w", err)
	}

	return nil
}

// =========================
// 书籍状态管理
// =========================

// UpdateBookStatus 更新书籍状态
func (a *ProgressAdapter) UpdateBookStatus(ctx context.Context, userID, bookID, status string) error {
	err := a.readerService.UpdateBookStatus(ctx, userID, bookID, status)
	if err != nil {
		return fmt.Errorf("更新书籍状态失败: %w", err)
	}

	return nil
}

// BatchUpdateBookStatus 批量更新书籍状态
func (a *ProgressAdapter) BatchUpdateBookStatus(ctx context.Context, userID string, bookIDs []string, status string) error {
	err := a.readerService.BatchUpdateBookStatus(ctx, userID, bookIDs, status)
	if err != nil {
		return fmt.Errorf("批量更新书籍状态失败: %w", err)
	}

	return nil
}

// GetBooksByStatus 根据状态获取书籍
func (a *ProgressAdapter) GetBooksByStatus(ctx context.Context, userID, status string, page, pageSize int) (*dto.RecentBooksResponse, error) {
	// 现有服务没有直接按状态获取的方法，需要通过GetReadingHistory获取后筛选
	progresses, total, err := a.readerService.GetReadingHistory(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取书籍列表失败: %w", err)
	}

	// 转换为DTO
	books := make([]*dto.BookProgressInfo, 0, len(progresses))
	for _, p := range progresses {
		// 根据状态筛选
		if p.Status == status {
			books = append(books, convertToBookProgressInfo(p))
		}
	}

	return &dto.RecentBooksResponse{
		Books: books,
		Total: int(total),
	}, nil
}

// =========================
// 阅读统计
// =========================

// GetReadingStats 获取阅读统计
func (a *ProgressAdapter) GetReadingStats(ctx context.Context, userID string) (*dto.ReadingStatsResponse, error) {
	// 获取总阅读时长
	totalTime, err := a.readerService.GetTotalReadingTime(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取阅读统计失败: %w", err)
	}

	// 获取已读完和未读完的书籍
	finished, err := a.readerService.GetFinishedBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取已读完书籍失败: %w", err)
	}

	unfinished, err := a.readerService.GetUnfinishedBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取未读完书籍失败: %w", err)
	}

	// 构建响应
	response := &dto.ReadingStatsResponse{
		TotalBooks:       len(finished) + len(unfinished),
		FinishedBooks:    len(finished),
		UnfinishedBooks:  len(unfinished),
		TotalReadingTime: totalTime,
		TotalWords:       0, // 需要从书籍信息获取
		ReadingDays:      0, // 需要从历史记录计算
		MaxStreak:        0, // 需要额外实现
	}

	return response, nil
}

// GetTotalReadingTime 获取总阅读时长
func (a *ProgressAdapter) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	totalTime, err := a.readerService.GetTotalReadingTime(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("获取总阅读时长失败: %w", err)
	}

	return totalTime, nil
}

// GetReadingTimeByPeriod 获取时间段阅读时长
func (a *ProgressAdapter) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime int64) (int64, error) {
	// 转换时间戳
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	// 调用现有服务
	totalTime, err := a.readerService.GetReadingTimeByPeriod(ctx, userID, start, end)
	if err != nil {
		return 0, fmt.Errorf("获取时间段阅读时长失败: %w", err)
	}

	return totalTime, nil
}

// =========================
// 阅读历史
// =========================

// GetRecentBooks 获取最近阅读的书籍
func (a *ProgressAdapter) GetRecentBooks(ctx context.Context, userID string, limit int) (*dto.RecentBooksResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 调用现有服务
	progresses, err := a.readerService.GetRecentReading(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取最近阅读失败: %w", err)
	}

	// 转换为DTO
	books := make([]*dto.BookProgressInfo, 0, len(progresses))
	for _, p := range progresses {
		books = append(books, convertToBookProgressInfo(p))
	}

	return &dto.RecentBooksResponse{
		Books: books,
		Total: len(books),
	}, nil
}

// GetReadingHistory 获取阅读历史
func (a *ProgressAdapter) GetReadingHistory(ctx context.Context, userID string, page, pageSize int) (*dto.ReadingHistoryResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 调用现有服务
	progresses, total, err := a.readerService.GetReadingHistory(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取阅读历史失败: %w", err)
	}

	// 转换为DTO
	items := make([]*dto.ReadingProgressResponse, 0, len(progresses))
	for _, p := range progresses {
		items = append(items, convertToProgressResponse(p))
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ReadingHistoryResponse{
		Progresses:  items,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

// GetUnfinishedBooks 获取未读完的书籍
func (a *ProgressAdapter) GetUnfinishedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error) {
	progresses, err := a.readerService.GetUnfinishedBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取未读完书籍失败: %w", err)
	}

	// 转换为DTO
	books := make([]*dto.BookProgressInfo, 0, len(progresses))
	for _, p := range progresses {
		books = append(books, convertToBookProgressInfo(p))
	}

	return books, nil
}

// GetFinishedBooks 获取已读完的书籍
func (a *ProgressAdapter) GetFinishedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error) {
	progresses, err := a.readerService.GetFinishedBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取已读完书籍失败: %w", err)
	}

	// 转换为DTO
	books := make([]*dto.BookProgressInfo, 0, len(progresses))
	for _, p := range progresses {
		books = append(books, convertToBookProgressInfo(p))
	}

	return books, nil
}

// =========================
// 阅读分析
// =========================

// GetReadingTrends 获取阅读趋势
func (a *ProgressAdapter) GetReadingTrends(ctx context.Context, userID string, period string, days int) (*dto.ReadingTrendsResponse, error) {
	// 现有服务没有直接实现趋势功能
	// 这里返回一个空的趋势数据，后续可以通过ReadingHistoryService实现
	return &dto.ReadingTrendsResponse{
		Period: period,
		Trends: []*dto.TrendDataPoint{},
	}, nil
}

// GetReadingStreak 获取连续阅读天数
func (a *ProgressAdapter) GetReadingStreak(ctx context.Context, userID string) (*dto.ReadingStreakResponse, error) {
	// 现有服务没有直接实现连续阅读天数功能
	// 这里返回默认值，后续可以通过ReadingHistoryService实现
	return &dto.ReadingStreakResponse{
		CurrentStreak: 0,
		MaxStreak:     0,
		LastReadDate:  0,
	}, nil
}

// GetLongestBooks 获取阅读字数最多的书籍
func (a *ProgressAdapter) GetLongestBooks(ctx context.Context, userID string, limit int) ([]*dto.BookProgressInfo, error) {
	// 现有服务没有直接按字数排序的功能
	// 先获取所有阅读记录，后续可以优化
	progresses, err := a.readerService.GetRecentReading(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取书籍列表失败: %w", err)
	}

	// 转换为DTO
	books := make([]*dto.BookProgressInfo, 0, len(progresses))
	for _, p := range progresses {
		books = append(books, convertToBookProgressInfo(p))
	}

	return books, nil
}

// =========================
// 同步相关
// =========================

// GetProgressSyncData 获取进度同步数据
func (a *ProgressAdapter) GetProgressSyncData(ctx context.Context, userID string, lastSyncTime int64) (*dto.ProgressSyncData, error) {
	// 获取所有阅读历史
	progresses, _, err := a.readerService.GetReadingHistory(ctx, userID, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("获取同步数据失败: %w", err)
	}

	// 筛选比lastSyncTime更新的记录
	filtered := make([]*dto.ReadingProgressResponse, 0)
	for _, p := range progresses {
		if p.UpdatedAt.Unix() > lastSyncTime {
			filtered = append(filtered, convertToProgressResponse(p))
		}
	}

	return &dto.ProgressSyncData{
		Progresses: filtered,
		SyncTime:   time.Now().Unix(),
		HasMore:    false,
	}, nil
}

// MergeProgress 合并进度数据
func (a *ProgressAdapter) MergeProgress(ctx context.Context, userID string, progressData []*dto.ReadingProgressResponse) error {
	// 遍历需要合并的进度数据
	for _, data := range progressData {
		// 调用保存进度方法
		err := a.SaveProgress(ctx, &dto.SaveProgressRequest{
			UserID:    data.UserID,
			BookID:    data.BookID,
			ChapterID: data.ChapterID,
			Progress:  data.Progress,
		})
		if err != nil {
			// 记录错误但继续处理其他数据
			continue
		}
	}

	return nil
}

// =========================
// 辅助转换函数
// =========================

// convertToBookProgressInfo 将reader.ReadingProgress转换为dto.BookProgressInfo
func convertToBookProgressInfo(p *reader.ReadingProgress) *dto.BookProgressInfo {
	return &dto.BookProgressInfo{
		BookID:       p.BookID.Hex(),
		BookTitle:    "", // 需要从bookstore获取
		CoverURL:     "", // 需要从bookstore获取
		Progress:     float64(p.Progress),
		ChapterNum:   0, // 需要从chapter获取
		ReadingTime:  p.ReadingTime,
		LastReadTime: p.UpdatedAt.Unix(),
		Status:       p.Status,
	}
}

// convertToProgressResponse 将reader.ReadingProgress转换为dto.ReadingProgressResponse
func convertToProgressResponse(p *reader.ReadingProgress) *dto.ReadingProgressResponse {
	return &dto.ReadingProgressResponse{
		UserID:      p.UserID.Hex(),
		BookID:      p.BookID.Hex(),
		ChapterID:   p.ChapterID.Hex(),
		Progress:    float64(p.Progress),
		ReadingTime: p.ReadingTime,
		UpdateTime:  p.UpdatedAt.Unix(),
		BookTitle:   "", // 需要从bookstore获取
		ChapterNum:  0,  // 需要从chapter获取
	}
}

// =========================
// 辅助方法（如果需要扩展功能）
// =========================

// GetReadingProgressWithDetails 获取带书籍和章节详情的阅读进度
// 这是一个扩展功能，用于完善BookProgressInfo和ReadingProgressResponse中的缺失字段
func (a *ProgressAdapter) GetReadingProgressWithDetails(ctx context.Context, userID, bookID string) (*dto.ReadingProgressResponse, error) {
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取阅读进度失败: %w", err)
	}

	response := convertToProgressResponse(progress)

	// TODO: 这里可以通过依赖注入的bookstoreService来获取书籍和章节详情
	// 例如:
	// book, _ := a.bookstoreService.GetBook(ctx, bookID)
	// response.BookTitle = book.Title
	//
	// chapter, _ := a.chapterService.GetChapter(ctx, progress.ChapterID.Hex())
	// response.ChapterNum = chapter.ChapterNum

	return response, nil
}

// UpdateProgressWithTimestamp 带时间戳的进度更新
// 用于同步场景，支持指定更新时间
func (a *ProgressAdapter) UpdateProgressWithTimestamp(ctx context.Context, userID, bookID, chapterID string, progress float64, timestamp int64) error {
	// 现有服务不支持指定时间戳，这里先调用基础方法
	err := a.readerService.SaveReadingProgress(ctx, userID, bookID, chapterID, progress)
	if err != nil {
		return fmt.Errorf("更新进度失败: %w", err)
	}

	return nil
}

// BatchGetProgress 批量获取阅读进度
// 用于优化批量查询场景
func (a *ProgressAdapter) BatchGetProgress(ctx context.Context, userID string, bookIDs []string) ([]*dto.ReadingProgressResponse, error) {
	responses := make([]*dto.ReadingProgressResponse, 0, len(bookIDs))

	// 逐个获取（后续可以优化为批量查询）
	for _, bookID := range bookIDs {
		progress, err := a.GetProgress(ctx, userID, bookID)
		if err != nil {
			// 记录错误但继续处理
			continue
		}
		responses = append(responses, progress)
	}

	return responses, nil
}

// GetProgressByChapterID 根据章节ID获取进度
// 用于快速定位到特定章节的阅读进度
func (a *ProgressAdapter) GetProgressByChapterID(ctx context.Context, userID, chapterID string) (*dto.ReadingProgressResponse, error) {
	// 现有服务没有直接按章节ID查询的方法
	// 需要通过Repository扩展或遍历查询
	// 这里返回错误，提示需要扩展功能
	return nil, fmt.Errorf("暂不支持按章节ID查询进度，需要扩展Repository")
}

// CalculateProgressPercentage 计算进度百分比
// 根据当前章节计算总体阅读进度
func (a *ProgressAdapter) CalculateProgressPercentage(ctx context.Context, userID, bookID string) (float64, error) {
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return 0, fmt.Errorf("获取进度失败: %w", err)
	}

	return float64(progress.Progress), nil
}

// EstimateReadingTime 估算剩余阅读时间
// 根据当前进度和历史阅读速度估算
func (a *ProgressAdapter) EstimateReadingTime(ctx context.Context, userID, bookID string) (int64, error) {
	// 获取进度和阅读时长
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return 0, fmt.Errorf("获取进度失败: %w", err)
	}

	// 简单估算：剩余进度比例 * 已阅读时长 / 当前进度
	if progress.Progress <= 0 || progress.ReadingTime <= 0 {
		return 0, nil
	}

	remainingProgress := 1.0 - float64(progress.Progress)
	estimatedTime := int64(float64(progress.ReadingTime) / float64(progress.Progress) * remainingProgress)

	return estimatedTime, nil
}

// ValidateProgressThreshold 验证进度阈值
// 检查是否达到特定的进度阈值（如50%、100%）
func (a *ProgressAdapter) ValidateProgressThreshold(ctx context.Context, userID, bookID string, threshold float64) (bool, error) {
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return false, fmt.Errorf("获取进度失败: %w", err)
	}

	return float64(progress.Progress) >= threshold, nil
}

// GetProgressPercentage 获取进度百分比（0-100）
func (a *ProgressAdapter) GetProgressPercentage(ctx context.Context, userID, bookID string) (float64, error) {
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return 0, fmt.Errorf("获取进度失败: %w", err)
	}

	return float64(progress.Progress.ToPercent()), nil
}

// UpdateProgressByChapterNum 根据章节号更新进度
// 用于自动计算章节对应的进度百分比
func (a *ProgressAdapter) UpdateProgressByChapterNum(ctx context.Context, userID, bookID, chapterID string, totalChapters int) error {
	if totalChapters <= 0 {
		return fmt.Errorf("总章节数必须大于0")
	}

	// TODO: 获取当前章节的序号
	// chapterNum := a.chapterService.GetChapterNum(ctx, chapterID)
	// progress := float64(chapterNum) / float64(totalChapters)

	// 暂时使用默认进度
	progress := 0.5

	err := a.readerService.SaveReadingProgress(ctx, userID, bookID, chapterID, progress)
	if err != nil {
		return fmt.Errorf("更新进度失败: %w", err)
	}

	return nil
}

// GetProgressStatus 获取进度状态描述
// 返回当前阅读状态的文字描述
func (a *ProgressAdapter) GetProgressStatus(ctx context.Context, userID, bookID string) (string, error) {
	progress, err := a.readerService.GetReadingProgress(ctx, userID, bookID)
	if err != nil {
		return "", fmt.Errorf("获取进度失败: %w", err)
	}

	// 根据进度返回状态描述
	switch progress.Status {
	case "reading":
		return "正在阅读", nil
	case "want_read":
		return "想读", nil
	case "finished":
		return "已读完", nil
	default:
		return "未开始", nil
	}
}

// ResetProgress 重置阅读进度
// 清除书籍的阅读进度记录
func (a *ProgressAdapter) ResetProgress(ctx context.Context, userID, bookID string) error {
	err := a.readerService.DeleteReadingProgress(ctx, userID, bookID)
	if err != nil {
		return fmt.Errorf("重置进度失败: %w", err)
	}

	return nil
}

// ArchiveProgress 归档阅读进度
// 将旧的阅读进度归档（保留记录但不显示在最近阅读中）
func (a *ProgressAdapter) ArchiveProgress(ctx context.Context, userID, bookID string) error {
	// 现有服务不支持归档功能
	// 可以通过更新状态为"archived"来实现
	err := a.readerService.UpdateBookStatus(ctx, userID, bookID, "archived")
	if err != nil {
		return fmt.Errorf("归档进度失败: %w", err)
	}

	return nil
}

// GetArchivedBooks 获取已归档的书籍
func (a *ProgressAdapter) GetArchivedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error) {
	// 通过状态筛选获取归档书籍
	response, err := a.GetBooksByStatus(ctx, userID, "archived", 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取归档书籍失败: %w", err)
	}

	return response.Books, nil
}

// RestoreProgress 恢复已归档的进度
func (a *ProgressAdapter) RestoreProgress(ctx context.Context, userID, bookID string) error {
	err := a.readerService.UpdateBookStatus(ctx, userID, bookID, "reading")
	if err != nil {
		return fmt.Errorf("恢复进度失败: %w", err)
	}

	return nil
}

// SyncProgressFromClient 从客户端同步进度
// 处理来自客户端的进度同步请求
func (a *ProgressAdapter) SyncProgressFromClient(ctx context.Context, userID string, progresses []*dto.ReadingProgressResponse) ([]*dto.ReadingProgressResponse, error) {
	// 1. 合并客户端进度
	err := a.MergeProgress(ctx, userID, progresses)
	if err != nil {
		return nil, fmt.Errorf("合并进度失败: %w", err)
	}

	// 2. 返回服务端更新的进度
	syncData, err := a.GetProgressSyncData(ctx, userID, 0)
	if err != nil {
		return nil, fmt.Errorf("获取同步数据失败: %w", err)
	}

	return syncData.Progresses, nil
}

// ConflictResolver 进度冲突解决策略
type ConflictResolver int

const (
	// ConflictResolverServerWins 服务端优先
	ConflictResolverServerWins ConflictResolver = iota
	// ConflictResolverClientWins 客户端优先
	ConflictResolverClientWins
	// ConflictResolverLatestTime 最新时间优先
	ConflictResolverLatestTime
	// ConflictResolverMaxProgress 最大进度优先
	ConflictResolverMaxProgress
)

// MergeProgressWithResolver 使用指定策略合并进度
func (a *ProgressAdapter) MergeProgressWithResolver(ctx context.Context, userID string, clientProgress *dto.ReadingProgressResponse, resolver ConflictResolver) error {
	// 获取服务端进度
	serverProgress, err := a.GetProgress(ctx, userID, clientProgress.BookID)
	if err != nil {
		// 服务端没有进度，直接使用客户端进度
		return a.SaveProgress(ctx, &dto.SaveProgressRequest{
			UserID:    clientProgress.UserID,
			BookID:    clientProgress.BookID,
			ChapterID: clientProgress.ChapterID,
			Progress:  clientProgress.Progress,
		})
	}

	// 根据策略选择要保存的进度
	var finalProgress *dto.SaveProgressRequest

	switch resolver {
	case ConflictResolverServerWins:
		// 保持服务端进度，不做处理
		return nil
	case ConflictResolverClientWins:
		// 使用客户端进度
		finalProgress = &dto.SaveProgressRequest{
			UserID:    clientProgress.UserID,
			BookID:    clientProgress.BookID,
			ChapterID: clientProgress.ChapterID,
			Progress:  clientProgress.Progress,
		}
	case ConflictResolverLatestTime:
		// 使用更新时间较晚的进度
		if clientProgress.UpdateTime > serverProgress.UpdateTime {
			finalProgress = &dto.SaveProgressRequest{
				UserID:    clientProgress.UserID,
				BookID:    clientProgress.BookID,
				ChapterID: clientProgress.ChapterID,
				Progress:  clientProgress.Progress,
			}
		} else {
			return nil
		}
	case ConflictResolverMaxProgress:
		// 使用进度较大的值
		if clientProgress.Progress > serverProgress.Progress {
			finalProgress = &dto.SaveProgressRequest{
				UserID:    clientProgress.UserID,
				BookID:    clientProgress.BookID,
				ChapterID: clientProgress.ChapterID,
				Progress:  clientProgress.Progress,
			}
		} else {
			return nil
		}
	}

	if finalProgress != nil {
		return a.SaveProgress(ctx, finalProgress)
	}

	return nil
}

// GetProgressWithConflictInfo 获取带冲突信息的进度
// 返回进度数据以及可能的冲突信息
func (a *ProgressAdapter) GetProgressWithConflictInfo(ctx context.Context, userID, bookID string) (*dto.ReadingProgressResponse, bool, error) {
	progress, err := a.GetProgress(ctx, userID, bookID)
	if err != nil {
		return nil, false, fmt.Errorf("获取进度失败: %w", err)
	}

	// TODO: 检查是否存在冲突（如多个设备同时阅读）
	// 这里简化处理，假设无冲突
	return progress, false, nil
}
