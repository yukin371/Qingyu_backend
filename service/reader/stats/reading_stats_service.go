package stats

import (
	"Qingyu_backend/models/stats"
	statsRepo "Qingyu_backend/repository/interfaces/stats"
	"context"
	"fmt"
	"math"
	"sort"
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
	uniqueReaders, err := s.readerBehaviorRepo.CountUniqueReadersByChapter(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("统计独立读者数失败: %w", err)
	}

	avgReadTime, err := s.readerBehaviorRepo.CalculateAvgReadTime(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("计算平均阅读时长失败: %w", err)
	}

	completionRate, err := s.CalculateCompletionRate(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("计算完读率失败: %w", err)
	}

	dropOffRate, err := s.CalculateDropOffRate(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("计算跳出率失败: %w", err)
	}

	chapterStats, err := s.chapterStatsRepo.GetByChapterID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节统计失败: %w", err)
	}

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

// GenerateReadingTimeHeatmap 生成年/月阅读时段热力图（按星期*小时聚合）
func (s *ReadingStatsService) GenerateReadingTimeHeatmap(ctx context.Context, bookID string, days int) ([]map[string]int, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	behaviors, err := s.readerBehaviorRepo.GetByDateRange(ctx, bookID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取阅读行为失败: %w", err)
	}

	grid := make(map[[2]int]int)
	for _, behavior := range behaviors {
		readAt := behavior.ReadAt
		day := int(readAt.Weekday())
		if day == 0 {
			day = 6
		} else {
			day--
		}
		hour := readAt.Hour()
		grid[[2]int{hour, day}]++
	}

	result := make([]map[string]int, 0, 24*7)
	for day := 0; day < 7; day++ {
		for hour := 0; hour < 24; hour++ {
			result = append(result, map[string]int{
				"hour":  hour,
				"day":   day,
				"value": grid[[2]int{hour, day}],
			})
		}
	}

	return result, nil
}

// CalculateCompletionRate 计算完读率
func (s *ReadingStatsService) CalculateCompletionRate(ctx context.Context, chapterID string) (float64, error) {
	totalCount, err := s.readerBehaviorRepo.CountByChapter(ctx, chapterID)
	if err != nil {
		return 0, fmt.Errorf("获取总行为数失败: %w", err)
	}
	if totalCount == 0 {
		return 0, nil
	}

	completeCount, err := s.readerBehaviorRepo.CountByChapterAndType(ctx, chapterID, stats.BehaviorTypeComplete)
	if err != nil {
		return 0, fmt.Errorf("获取完读行为数失败: %w", err)
	}
	return float64(completeCount) / float64(totalCount), nil
}

// CalculateDropOffRate 计算跳出率
func (s *ReadingStatsService) CalculateDropOffRate(ctx context.Context, chapterID string) (float64, error) {
	totalCount, err := s.readerBehaviorRepo.CountByChapter(ctx, chapterID)
	if err != nil {
		return 0, fmt.Errorf("获取总行为数失败: %w", err)
	}
	if totalCount == 0 {
		return 0, nil
	}

	dropOffCount, err := s.readerBehaviorRepo.CountByChapterAndType(ctx, chapterID, stats.BehaviorTypeDropOff)
	if err != nil {
		return 0, fmt.Errorf("获取跳出行为数失败: %w", err)
	}
	return float64(dropOffCount) / float64(totalCount), nil
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
	mostViewed, err := s.chapterStatsRepo.GetTopViewedChapters(ctx, bookID, 10)
	if err != nil {
		return nil, fmt.Errorf("获取阅读量最高章节失败: %w", err)
	}
	highestRevenue, err := s.chapterStatsRepo.GetTopRevenueChapters(ctx, bookID, 10)
	if err != nil {
		return nil, fmt.Errorf("获取收入最高章节失败: %w", err)
	}
	lowestCompletion, err := s.chapterStatsRepo.GetLowestCompletionChapters(ctx, bookID, 10)
	if err != nil {
		return nil, fmt.Errorf("获取完读率最低章节失败: %w", err)
	}
	highestDropOff, err := s.chapterStatsRepo.GetHighestDropOffChapters(ctx, bookID, 10)
	if err != nil {
		return nil, fmt.Errorf("获取跳出率最高章节失败: %w", err)
	}

	return &stats.TopChapters{
		BookID:           bookID,
		MostViewed:       mostViewed,
		HighestRevenue:   highestRevenue,
		LowestCompletion: lowestCompletion,
		HighestDropOff:   highestDropOff,
	}, nil
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
	now := time.Now()
	targetDate := now.AddDate(0, 0, -days)
	targetDateEnd := targetDate.Add(24 * time.Hour)

	targetReaders, err := s.readerBehaviorRepo.GetDistinctUsersByBookAndDateRange(ctx, bookID, targetDate, targetDateEnd)
	if err != nil {
		return 0, fmt.Errorf("获取目标读者失败: %w", err)
	}
	if len(targetReaders) == 0 {
		return 0, nil
	}

	today := now.Truncate(24 * time.Hour)
	activeCount, err := s.readerBehaviorRepo.CountActiveUsersInList(ctx, bookID, targetReaders, today)
	if err != nil {
		return 0, fmt.Errorf("统计活跃读者失败: %w", err)
	}
	return float64(activeCount) / float64(len(targetReaders)), nil
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

// GetChapterRankings 获取章节排行列表
func (s *ReadingStatsService) GetChapterRankings(ctx context.Context, bookID string, page, size int) ([]*stats.ChapterStatsAggregate, int, error) {
	all, err := s.chapterStatsRepo.GetChapterStatsAggregate(ctx, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("获取章节统计聚合失败: %w", err)
	}

	sort.SliceStable(all, func(i, j int) bool {
		if all[i].ViewCount == all[j].ViewCount {
			return all[i].Revenue > all[j].Revenue
		}
		return all[i].ViewCount > all[j].ViewCount
	})

	total := len(all)
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * size
	if start >= total {
		return []*stats.ChapterStatsAggregate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// GetSubscribersTrend 获取订阅趋势
func (s *ReadingStatsService) GetSubscribersTrend(ctx context.Context, bookID string, days int) ([]*stats.BookStatsDaily, error) {
	return s.GetDailyStats(ctx, bookID, days)
}

// GetReaderActivity 获取读者活跃度分布
func (s *ReadingStatsService) GetReaderActivity(ctx context.Context, bookID string) ([]map[string]interface{}, error) {
	totalUniqueReaders, err := s.readerBehaviorRepo.CountUniqueReaders(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取独立读者数失败: %w", err)
	}

	now := time.Now()
	startOfDay := now.Truncate(24 * time.Hour)
	dailyUsers, err := s.readerBehaviorRepo.GetDistinctUsersByBookAndDateRange(ctx, bookID, startOfDay, now)
	if err != nil {
		return nil, fmt.Errorf("获取日活读者失败: %w", err)
	}
	weeklyUsers, err := s.readerBehaviorRepo.GetDistinctUsersByBookAndDateRange(ctx, bookID, now.AddDate(0, 0, -7), now)
	if err != nil {
		return nil, fmt.Errorf("获取周活读者失败: %w", err)
	}
	monthlyUsers, err := s.readerBehaviorRepo.GetDistinctUsersByBookAndDateRange(ctx, bookID, now.AddDate(0, 0, -30), now)
	if err != nil {
		return nil, fmt.Errorf("获取月活读者失败: %w", err)
	}

	monthlyCount := int64(len(monthlyUsers))
	inactiveCount := totalUniqueReaders - monthlyCount
	if inactiveCount < 0 {
		inactiveCount = 0
	}

	buildItem := func(activityType, label string, count int64) map[string]interface{} {
		percentage := 0.0
		if totalUniqueReaders > 0 {
			percentage = math.Round((float64(count)/float64(totalUniqueReaders))*10000) / 100
		}
		return map[string]interface{}{
			"type":       activityType,
			"label":      label,
			"count":      count,
			"percentage": percentage,
		}
	}

	return []map[string]interface{}{
		buildItem("daily", "日活跃", int64(len(dailyUsers))),
		buildItem("weekly", "周活跃", int64(len(weeklyUsers))),
		buildItem("monthly", "月活跃", monthlyCount),
		buildItem("inactive", "不活跃", inactiveCount),
	}, nil
}

// CompareBooks 对比多个作品的核心指标
func (s *ReadingStatsService) CompareBooks(ctx context.Context, bookIDs []string, metrics []string, startDate, endDate *time.Time) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0, len(bookIDs))
	for _, bookID := range bookIDs {
		bookStats, err := s.bookStatsRepo.GetByBookID(ctx, bookID)
		if err != nil {
			return nil, fmt.Errorf("获取作品 %s 统计失败: %w", bookID, err)
		}

		var rangeStats *stats.TimeRangeStats
		if startDate != nil && endDate != nil {
			rangeStats, err = s.chapterStatsRepo.GetTimeRangeStats(ctx, bookID, *startDate, *endDate)
			if err != nil {
				return nil, fmt.Errorf("获取作品 %s 时间范围统计失败: %w", bookID, err)
			}
		}

		item := map[string]interface{}{
			"bookId": bookID,
			"title":  "",
		}
		if bookStats != nil {
			item["title"] = bookStats.Title
		}

		for _, metric := range metrics {
			switch metric {
			case "views":
				if rangeStats != nil {
					item["views"] = rangeStats.TotalViews
				} else if bookStats != nil {
					item["views"] = bookStats.TotalViews
				} else {
					item["views"] = int64(0)
				}
			case "subscribers":
				if bookStats != nil {
					item["subscribers"] = bookStats.TotalSubscribers
				} else {
					item["subscribers"] = int64(0)
				}
			case "favorites":
				if bookStats != nil {
					item["favorites"] = bookStats.TotalBookmarks
				} else {
					item["favorites"] = int64(0)
				}
			case "comments":
				if bookStats != nil {
					item["comments"] = bookStats.TotalComments
				} else {
					item["comments"] = int64(0)
				}
			case "revenue":
				if rangeStats != nil {
					item["revenue"] = rangeStats.TotalRevenue
				} else if bookStats != nil {
					item["revenue"] = bookStats.TotalRevenue
				} else {
					item["revenue"] = float64(0)
				}
			case "retention":
				retention, err := s.CalculateRetention(ctx, bookID, 7)
				if err != nil {
					return nil, fmt.Errorf("计算作品 %s 留存率失败: %w", bookID, err)
				}
				item["retention"] = retention
			}
		}

		results = append(results, item)
	}

	return results, nil
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
