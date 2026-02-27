package admin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyticsService_GetUserGrowthTrend_Success 测试获取用户增长趋势成功
func TestAnalyticsService_GetUserGrowthTrend_Success(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

	req := &UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	}

	resp, err := service.GetUserGrowthTrend(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, startDate, resp.StartDate)
	assert.Equal(t, endDate, resp.EndDate)
	assert.Equal(t, "daily", resp.Interval)
	assert.Greater(t, resp.TotalNewUsers, int64(0))
	assert.NotEmpty(t, resp.Data)
	assert.Greater(t, len(resp.Data), 0)
}

// TestAnalyticsService_GetUserGrowthTrend_WeeklyInterval 测试周间隔
func TestAnalyticsService_GetUserGrowthTrend_WeeklyInterval(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	req := &UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "weekly",
	}

	resp, err := service.GetUserGrowthTrend(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "weekly", resp.Interval)
}

// TestAnalyticsService_GetUserGrowthTrend_MonthlyInterval 测试月间隔
func TestAnalyticsService_GetUserGrowthTrend_MonthlyInterval(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	req := &UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "monthly",
	}

	resp, err := service.GetUserGrowthTrend(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "monthly", resp.Interval)
}

// TestAnalyticsService_GetUserGrowthTrend_InvalidDateRange 测试无效日期范围
func TestAnalyticsService_GetUserGrowthTrend_InvalidDateRange(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	req := &UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	}

	resp, err := service.GetUserGrowthTrend(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "开始日期不能晚于结束日期")
}

// TestAnalyticsService_GetContentStatistics_Success 测试获取内容统计成功
func TestAnalyticsService_GetContentStatistics_Success(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	req := &ContentStatisticsRequest{}

	resp, err := service.GetContentStatistics(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Greater(t, resp.TotalBooks, int64(0))
	assert.Greater(t, resp.TotalChapters, int64(0))
	assert.Greater(t, resp.TotalComments, int64(0))
	assert.Greater(t, resp.TotalWords, int64(0))
	assert.GreaterOrEqual(t, resp.PendingReviews, int64(0))
	assert.NotEmpty(t, resp.CategoryStats)
	assert.NotEmpty(t, resp.TrendingBooks)
}

// TestAnalyticsService_GetContentStatistics_WithDateRange 测试带日期范围的内容统计
func TestAnalyticsService_GetContentStatistics_WithDateRange(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	req := &ContentStatisticsRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	resp, err := service.GetContentStatistics(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestAnalyticsService_GetRevenueReport_Success 测试获取收入报告成功
func TestAnalyticsService_GetRevenueReport_Success(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	req := &RevenueReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	}

	resp, err := service.GetRevenueReport(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, startDate, resp.StartDate)
	assert.Equal(t, endDate, resp.EndDate)
	assert.Equal(t, "daily", resp.Interval)
	assert.Greater(t, resp.TotalRevenue, 0.0)
	assert.NotEmpty(t, resp.Data)
}

// TestAnalyticsService_GetRevenueReport_WeeklyInterval 测试周间隔收入报告
func TestAnalyticsService_GetRevenueReport_WeeklyInterval(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)

	req := &RevenueReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "weekly",
	}

	resp, err := service.GetRevenueReport(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "weekly", resp.Interval)
	assert.NotEmpty(t, resp.Data)
	// 验证每个数据点都有订单数
	for _, dp := range resp.Data {
		assert.Greater(t, dp.Orders, int64(0))
	}
}

// TestAnalyticsService_GetRevenueReport_InvalidDateRange 测试无效日期范围
func TestAnalyticsService_GetRevenueReport_InvalidDateRange(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	req := &RevenueReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	}

	resp, err := service.GetRevenueReport(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "开始日期不能晚于结束日期")
}

// TestAnalyticsService_GetActiveUsersReport_DAU 测试获取DAU报告
func TestAnalyticsService_GetActiveUsersReport_DAU(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

	req := &ActiveUsersReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      "dau",
	}

	resp, err := service.GetActiveUsersReport(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, startDate, resp.StartDate)
	assert.Equal(t, endDate, resp.EndDate)
	assert.Equal(t, "dau", resp.Type)
	assert.Greater(t, resp.AverageActiveUsers, 0.0)
	assert.Greater(t, resp.PeakActiveUsers, int64(0))
	assert.NotEmpty(t, resp.PeakDate)
	assert.NotEmpty(t, resp.Data)
}

// TestAnalyticsService_GetActiveUsersReport_WAU 测试获取WAU报告
func TestAnalyticsService_GetActiveUsersReport_WAU(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	req := &ActiveUsersReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      "wau",
	}

	resp, err := service.GetActiveUsersReport(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "wau", resp.Type)
}

// TestAnalyticsService_GetActiveUsersReport_MAU 测试获取MAU报告
func TestAnalyticsService_GetActiveUsersReport_MAU(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)

	req := &ActiveUsersReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      "mau",
	}

	resp, err := service.GetActiveUsersReport(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "mau", resp.Type)
}

// TestAnalyticsService_GetActiveUsersReport_InvalidDateRange 测试无效日期范围
func TestAnalyticsService_GetActiveUsersReport_InvalidDateRange(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	req := &ActiveUsersReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      "dau",
	}

	resp, err := service.GetActiveUsersReport(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "开始日期不能晚于结束日期")
}

// TestAnalyticsService_GetSystemOverview_Success 测试获取系统概览成功
func TestAnalyticsService_GetSystemOverview_Success(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	resp, err := service.GetSystemOverview(ctx)

	require.NoError(t, err)
	assert.NotNil(t, resp)

	// 用户统计
	assert.Greater(t, resp.TotalUsers, int64(0))
	assert.GreaterOrEqual(t, resp.NewUsersToday, int64(0))
	assert.Greater(t, resp.ActiveUsers, int64(0))

	// 内容统计
	assert.Greater(t, resp.TotalBooks, int64(0))
	assert.Greater(t, resp.TotalChapters, int64(0))
	assert.Greater(t, resp.TotalComments, int64(0))
	assert.GreaterOrEqual(t, resp.PendingReviews, int64(0))

	// 收入统计
	assert.Greater(t, resp.TotalRevenue, 0.0)
	assert.GreaterOrEqual(t, resp.RevenueToday, 0.0)

	// 系统状态
	assert.NotEmpty(t, resp.SystemStatus)
	assert.False(t, resp.LastUpdated.IsZero())

	// 验证系统状态值
	assert.Contains(t, []string{"healthy", "degraded", "down"}, resp.SystemStatus)
}

// TestAnalyticsService_UserGrowthDataPoint_Structure 测试用户增长数据点结构
func TestAnalyticsService_UserGrowthDataPoint_Structure(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)

	req := &UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	}

	resp, err := service.GetUserGrowthTrend(ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, resp.Data)

	// 验证数据点结构
	for _, dp := range resp.Data {
		assert.NotEmpty(t, dp.Date, "日期不能为空")
		assert.GreaterOrEqual(t, dp.Count, int64(0), "计数不能为负数")
	}
}

// TestAnalyticsService_RevenueDataPoint_Structure 测试收入数据点结构
func TestAnalyticsService_RevenueDataPoint_Structure(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)

	req := &RevenueReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	}

	resp, err := service.GetRevenueReport(ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, resp.Data)

	// 验证数据点结构
	for _, dp := range resp.Data {
		assert.NotEmpty(t, dp.Date, "日期不能为空")
		assert.Greater(t, dp.Amount, 0.0, "金额必须大于0")
		assert.Greater(t, dp.Orders, int64(0), "订单数必须大于0")
	}
}

// TestAnalyticsService_ContentStatistics_CategoryStats 测试分类统计
func TestAnalyticsService_ContentStatistics_CategoryStats(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	req := &ContentStatisticsRequest{}

	resp, err := service.GetContentStatistics(ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, resp.CategoryStats)

	// 验证分类统计结构
	for _, cs := range resp.CategoryStats {
		assert.NotEmpty(t, cs.CategoryName, "分类名称不能为空")
		assert.GreaterOrEqual(t, cs.BookCount, int64(0), "书籍数不能为负数")
		assert.GreaterOrEqual(t, cs.ChapterCount, int64(0), "章节数不能为负数")
	}
}

// TestAnalyticsService_ContentStatistics_TrendingBooks 测试热门书籍
func TestAnalyticsService_ContentStatistics_TrendingBooks(t *testing.T) {
	service := NewAnalyticsService()

	ctx := context.Background()
	req := &ContentStatisticsRequest{}

	resp, err := service.GetContentStatistics(ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, resp.TrendingBooks)

	// 验证热门书籍结构
	for _, tb := range resp.TrendingBooks {
		assert.NotEmpty(t, tb.BookID, "书籍ID不能为空")
		assert.NotEmpty(t, tb.Title, "标题不能为空")
		assert.NotEmpty(t, tb.AuthorID, "作者ID不能为空")
		assert.GreaterOrEqual(t, tb.ViewCount, int64(0), "浏览数不能为负数")
		assert.GreaterOrEqual(t, tb.ChapterCount, int64(0), "章节数不能为负数")
	}
}
