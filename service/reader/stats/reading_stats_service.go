package stats

import (
	"Qingyu_backend/models/stats"
	statsRepo "Qingyu_backend/repository/interfaces/stats"
	"context"
	"fmt"
	"math"
	"time"
)

// ReadingStatsService 阅读统计服务
// 职责：书籍、章节、读者行为统计
// 领域：Reading/Bookstore业务域
type ReadingStatsService struct {
	chapterStatsRepo   statsRepo.ChapterStatsRepository
	readerBehaviorRepo statsRepo.ReaderBehaviorRepository
	bookStatsRepo      statsRepo.BookStatsRepository
	serviceName        string
	version            string
}

// NewReadingStatsService 创建阅读统计服务
func NewReadingStatsService(
	chapterStatsRepo statsRepo.ChapterStatsRepository,
	readerBehaviorRepo statsRepo.ReaderBehaviorRepository,
	bookStatsRepo statsRepo.BookStatsRepository,
) *ReadingStatsService {
	return &ReadingStatsService{
		chapterStatsRepo:   chapterStatsRepo,
		readerBehaviorRepo: readerBehaviorRepo,
		bookStatsRepo:      bookStatsRepo,
		serviceName:        "ReadingStatsService",
		version:            "1.0.0",
	}
}

// CalculateChapterStats 计算章节统计数据
func (s *ReadingStatsService) CalculateChapterStats(ctx context.Context, chapterID string) (*stats.ChapterStats, error) {
	// 1. 获取基础统计
	uniqueReaders, err := s.readerBehaviorRepo.CountUniqueReadersByChapter(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("统计独立读者数失败: %w", err)
	}

	// 2. 计算平均阅读时长
	avgReadTime, err := s.readerBehaviorRepo.CalculateAvgReadTime(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("计算平均阅读时长失败: %w", err)
	}

	// 3. 计算完读率
	completionRate, err := s.readerBehaviorRepo.CalculateCompletionRate(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("计算完读率失败: %w", err)
	}

	// 4. 计算跳出率
	dropOffRate, err := s.readerBehaviorRepo.CalculateDropOffRate(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("计算跳出率失败: %w", err)
	}

	// 5. 获取现有统计记录
	chapterStats, err := s.chapterStatsRepo.GetByChapterID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节统计失败: %w", err)
	}

	// 6. 更新统计数据
	if chapterStats != nil {
		updates := map[string]interface{}{
			"unique_viewers":  uniqueReaders,
			"avg_read_time":   avgReadTime,
			"completion_rate": completionRate,
			"drop_off_rate":   dropOffRate,
			"updated_at":      time.Now(),
		}

		err = s.chapterStatsRepo.Update(ctx, chapterStats.ID, updates)
		if err != nil {
			return nil, fmt.Errorf("更新章节统计失败: %w", err)
		}

		// 重新获取更新后的数据
		chapterStats, err = s.chapterStatsRepo.GetByChapterID(ctx, chapterID)
		if err != nil {
			return nil, fmt.Errorf("获取更新后的章节统计失败: %w", err)
		}
	}

	return chapterStats, nil
}

// CalculateBookStats 计算作品统计数据
func (s *ReadingStatsService) CalculateBookStats(ctx context.Context, bookID string) (*stats.BookStats, error) {
	// 1. 统计独立读者数
	uniqueReaders, err := s.readerBehaviorRepo.CountUniqueReaders(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("统计独立读者数失败: %w", err)
	}

	// 2. 计算平均完读率
	avgCompletionRate, err := s.bookStatsRepo.CalculateAvgCompletionRate(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("计算平均完读率失败: %w", err)
	}

	// 3. 计算平均跳出率
	avgDropOffRate, err := s.bookStatsRepo.CalculateAvgDropOffRate(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("计算平均跳出率失败: %w", err)
	}

	// 4. 计算平均阅读时长
	avgReadingDuration, err := s.bookStatsRepo.CalculateAvgReadingDuration(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("计算平均阅读时长失败: %w", err)
	}

	// 5. 计算总收入
	totalRevenue, err := s.bookStatsRepo.CalculateTotalRevenue(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("计算总收入失败: %w", err)
	}

	// 6. 分析阅读量趋势
	viewTrend, err := s.bookStatsRepo.AnalyzeViewTrend(ctx, bookID, 7)
	if err != nil {
		viewTrend = stats.TrendStable
	}

	// 7. 分析收入趋势
	revenueTrend, err := s.bookStatsRepo.AnalyzeRevenueTrend(ctx, bookID, 7)
	if err != nil {
		revenueTrend = stats.TrendStable
	}

	// 8. 获取现有统计记录
	bookStats, err := s.bookStatsRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取作品统计失败: %w", err)
	}

	// 9. 更新统计数据
	if bookStats != nil {
		updates := map[string]interface{}{
			"unique_readers":       uniqueReaders,
			"avg_completion_rate":  avgCompletionRate,
			"avg_drop_off_rate":    avgDropOffRate,
			"avg_reading_duration": avgReadingDuration,
			"total_revenue":        totalRevenue,
			"view_trend":           viewTrend,
			"revenue_trend":        revenueTrend,
			"updated_at":           time.Now(),
		}

		err = s.bookStatsRepo.Update(ctx, bookStats.ID, updates)
		if err != nil {
			return nil, fmt.Errorf("更新作品统计失败: %w", err)
		}

		// 重新获取更新后的数据
		bookStats, err = s.bookStatsRepo.GetByBookID(ctx, bookID)
		if err != nil {
			return nil, fmt.Errorf("获取更新后的作品统计失败: %w", err)
		}
	}

	return bookStats, nil
}

// GenerateHeatmap 生成阅读热力图
func (s *ReadingStatsService) GenerateHeatmap(ctx context.Context, bookID string) ([]*stats.HeatmapPoint, error) {
	// 调用Repository层生成热力图
	heatmap, err := s.chapterStatsRepo.GenerateHeatmap(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("生成热力图失败: %w", err)
	}

	// 计算热度分数 (0-100)
	if len(heatmap) > 0 {
		// 找到最大值用于归一化
		maxViews := int64(0)
		for _, point := range heatmap {
			if point.ViewCount > maxViews {
				maxViews = point.ViewCount
			}
		}

		// 计算每个点的热度分数
		for _, point := range heatmap {
			// 热度分数 = 阅读量权重(50%) + 完读率权重(30%) + (1-跳出率)权重(20%)
			viewScore := float64(point.ViewCount) / float64(maxViews) * 50
			completionScore := point.CompletionRate * 30
			dropOffScore := (1 - point.DropOffRate) * 20

			point.HeatScore = viewScore + completionScore + dropOffScore

			// 确保分数在0-100之间
			if point.HeatScore > 100 {
				point.HeatScore = 100
			}
			if point.HeatScore < 0 {
				point.HeatScore = 0
			}
		}
	}

	return heatmap, nil
}

// CalculateCompletionRate 计算完读率
func (s *ReadingStatsService) CalculateCompletionRate(ctx context.Context, chapterID string) (float64, error) {
	return s.readerBehaviorRepo.CalculateCompletionRate(ctx, chapterID)
}

// CalculateDropOffPoints 计算跳出点
func (s *ReadingStatsService) CalculateDropOffPoints(ctx context.Context, bookID string) ([]*stats.ChapterStatsAggregate, error) {
	// 获取跳出率最高的章节
	dropOffChapters, err := s.chapterStatsRepo.GetHighestDropOffChapters(ctx, bookID, 10)
	if err != nil {
		return nil, fmt.Errorf("获取跳出点失败: %w", err)
	}

	return dropOffChapters, nil
}

// GetTimeRangeStats 获取时间范围统计
func (s *ReadingStatsService) GetTimeRangeStats(ctx context.Context, bookID string, startDate, endDate time.Time) (*stats.TimeRangeStats, error) {
	return s.chapterStatsRepo.GetTimeRangeStats(ctx, bookID, startDate, endDate)
}

// GetRevenueBreakdown 获取收入细分
func (s *ReadingStatsService) GetRevenueBreakdown(ctx context.Context, bookID string, startDate, endDate time.Time) (*stats.RevenueBreakdown, error) {
	return s.bookStatsRepo.GetRevenueBreakdown(ctx, bookID, startDate, endDate)
}

// GetTopChapters 获取热门章节
func (s *ReadingStatsService) GetTopChapters(ctx context.Context, bookID string) (*stats.TopChapters, error) {
	return s.bookStatsRepo.GetTopChapters(ctx, bookID)
}

// RecordReaderBehavior 记录读者行为
func (s *ReadingStatsService) RecordReaderBehavior(ctx context.Context, behavior *stats.ReaderBehavior) error {
	// 设置创建时间
	behavior.CreatedAt = time.Now()

	// 保存行为记录
	err := s.readerBehaviorRepo.Create(ctx, behavior)
	if err != nil {
		return fmt.Errorf("记录读者行为失败: %w", err)
	}

	// 异步更新统计数据
	go func() {
		// 使用新的context避免请求超时
		bgCtx := context.Background()

		// 更新章节统计
		_, _ = s.CalculateChapterStats(bgCtx, behavior.ChapterID)

		// 更新作品统计
		_, _ = s.CalculateBookStats(bgCtx, behavior.BookID)
	}()

	return nil
}

// CalculateRetention 计算留存率
func (s *ReadingStatsService) CalculateRetention(ctx context.Context, bookID string, days int) (float64, error) {
	return s.readerBehaviorRepo.CalculateRetention(ctx, bookID, days)
}

// GetDailyStats 获取每日统计
func (s *ReadingStatsService) GetDailyStats(ctx context.Context, bookID string, days int) ([]*stats.BookStatsDaily, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	dailyStats, err := s.bookStatsRepo.GetDailyStatsRange(ctx, bookID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取每日统计失败: %w", err)
	}

	return dailyStats, nil
}

// AnalyzeTrend 分析趋势
func (s *ReadingStatsService) AnalyzeTrend(ctx context.Context, bookID string, metric string, days int) (string, error) {
	switch metric {
	case "view":
		return s.bookStatsRepo.AnalyzeViewTrend(ctx, bookID, days)
	case "revenue":
		return s.bookStatsRepo.AnalyzeRevenueTrend(ctx, bookID, days)
	default:
		return "", fmt.Errorf("不支持的指标类型: %s", metric)
	}
}

// CalculateAvgRevenuePerUser 计算用户平均贡献
func (s *ReadingStatsService) CalculateAvgRevenuePerUser(ctx context.Context, bookID string) (float64, error) {
	// 获取总收入
	totalRevenue, err := s.bookStatsRepo.CalculateTotalRevenue(ctx, bookID)
	if err != nil {
		return 0, err
	}

	// 获取独立读者数
	uniqueReaders, err := s.readerBehaviorRepo.CountUniqueReaders(ctx, bookID)
	if err != nil {
		return 0, err
	}

	if uniqueReaders == 0 {
		return 0, nil
	}

	// 计算平均值
	avgRevenue := totalRevenue / float64(uniqueReaders)

	// 保留两位小数
	return math.Round(avgRevenue*100) / 100, nil
}

// Health 健康检查
func (s *ReadingStatsService) Health(ctx context.Context) error {
	// 检查所有Repository的健康状态
	if err := s.chapterStatsRepo.Health(ctx); err != nil {
		return fmt.Errorf("章节统计Repository健康检查失败: %w", err)
	}

	if err := s.readerBehaviorRepo.Health(ctx); err != nil {
		return fmt.Errorf("读者行为Repository健康检查失败: %w", err)
	}

	if err := s.bookStatsRepo.Health(ctx); err != nil {
		return fmt.Errorf("作品统计Repository健康检查失败: %w", err)
	}

	return nil
}

// GetServiceName 获取服务名称
func (s *ReadingStatsService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *ReadingStatsService) GetVersion() string {
	return s.version
}

// Initialize 初始化服务
func (s *ReadingStatsService) Initialize(ctx context.Context) error {
	// ReadingStatsService 无需特殊初始化
	return nil
}

// Close 关闭服务
func (s *ReadingStatsService) Close(ctx context.Context) error {
	// ReadingStatsService 无需清理资源
	return nil
}
