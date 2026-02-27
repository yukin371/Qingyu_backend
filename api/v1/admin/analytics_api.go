package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	adminService "Qingyu_backend/service/admin"
)

// AnalyticsAPI 统计分析API
type AnalyticsAPI struct {
	analyticsService adminService.AnalyticsService
}

// NewAnalyticsAPI 创建统计分析API实例
func NewAnalyticsAPI(analyticsService adminService.AnalyticsService) *AnalyticsAPI {
	return &AnalyticsAPI{
		analyticsService: analyticsService,
	}
}

// GetUserGrowthTrend 获取用户增长趋势
// @Summary 获取用户增长趋势
// @Description 按日期统计用户注册数量，支持日/周/月间隔
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param start_date query string true "开始日期 (YYYY-MM-DD)"
// @Param end_date query string true "结束日期 (YYYY-MM-DD)"
// @Param interval query string true "间隔 (daily/weekly/monthly)" Enums(daily, weekly, monthly)
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/user-growth [get]
func (api *AnalyticsAPI) GetUserGrowthTrend(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	interval := c.DefaultQuery("interval", "daily")

	// 参数验证
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始日期和结束日期不能为空",
		})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("开始日期格式错误: %v", err),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("结束日期格式错误: %v", err),
		})
		return
	}

	// 验证间隔参数
	if interval != "daily" && interval != "weekly" && interval != "monthly" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "间隔参数必须是 daily、weekly 或 monthly",
		})
		return
	}

	// 构建请求
	req := &adminService.UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  interval,
	}

	// 调用服务
	resp, err := api.analyticsService.GetUserGrowthTrend(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取用户增长趋势失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    resp,
	})
}

// GetContentStatistics 获取内容统计
// @Summary 获取内容统计
// @Description 统计书籍、章节、评论等数据
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/content-statistics [get]
func (api *AnalyticsAPI) GetContentStatistics(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// 解析日期（可选）
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": fmt.Sprintf("开始日期格式错误: %v", err),
			})
			return
		}
	}

	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": fmt.Sprintf("结束日期格式错误: %v", err),
			})
			return
		}
	}

	// 构建请求
	req := &adminService.ContentStatisticsRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 调用服务
	resp, err := api.analyticsService.GetContentStatistics(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取内容统计失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    resp,
	})
}

// GetRevenueReport 获取收入报告
// @Summary 获取收入报告
// @Description 按日期统计收入数据，支持日/周/月间隔
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param start_date query string true "开始日期 (YYYY-MM-DD)"
// @Param end_date query string true "结束日期 (YYYY-MM-DD)"
// @Param interval query string true "间隔 (daily/weekly/monthly)" Enums(daily, weekly, monthly)
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/revenue-report [get]
func (api *AnalyticsAPI) GetRevenueReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	interval := c.DefaultQuery("interval", "daily")

	// 参数验证
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始日期和结束日期不能为空",
		})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("开始日期格式错误: %v", err),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("结束日期格式错误: %v", err),
		})
		return
	}

	// 验证间隔参数
	if interval != "daily" && interval != "weekly" && interval != "monthly" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "间隔参数必须是 daily、weekly 或 monthly",
		})
		return
	}

	// 构建请求
	req := &adminService.RevenueReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  interval,
	}

	// 调用服务
	resp, err := api.analyticsService.GetRevenueReport(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取收入报告失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    resp,
	})
}

// GetActiveUsersReport 获取活跃用户报告
// @Summary 获取活跃用户报告
// @Description 统计活跃用户数据，支持DAU/WAU/MAU
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param start_date query string true "开始日期 (YYYY-MM-DD)"
// @Param end_date query string true "结束日期 (YYYY-MM-DD)"
// @Param type query string true "类型 (dau/wau/mau)" Enums(dau, wau, mau)
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/active-users [get]
func (api *AnalyticsAPI) GetActiveUsersReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	reportType := c.Query("type")

	// 参数验证
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始日期和结束日期不能为空",
		})
		return
	}

	if reportType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "报告类型不能为空",
		})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("开始日期格式错误: %v", err),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("结束日期格式错误: %v", err),
		})
		return
	}

	// 验证类型参数
	if reportType != "dau" && reportType != "wau" && reportType != "mau" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "类型参数必须是 dau、wau 或 mau",
		})
		return
	}

	// 构建请求
	req := &adminService.ActiveUsersReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      reportType,
	}

	// 调用服务
	resp, err := api.analyticsService.GetActiveUsersReport(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取活跃用户报告失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    resp,
	})
}

// GetSystemOverview 获取系统概览
// @Summary 获取系统概览
// @Description 返回系统整体数据概览，包括用户、内容、收入等统计
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/system-overview [get]
func (api *AnalyticsAPI) GetSystemOverview(c *gin.Context) {
	// 调用服务
	resp, err := api.analyticsService.GetSystemOverview(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取系统概览失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    resp,
	})
}

// ExportAnalyticsReport 导出统计分析报告
// @Summary 导出统计分析报告
// @Description 导出各种统计报告为CSV或Excel格式
// @Tags Admin-Analytics
// @Accept json
// @Produce csv/xlsx
// @Param report_type query string true "报告类型 (user-growth/content-statistics/revenue-report/active-users)" Enums(user-growth, content-statistics, revenue-report, active-users)
// @Param format query string false "导出格式 (csv/xlsx)" default(csv)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param interval query string false "间隔 (daily/weekly/monthly)"
// @Param type query string false "报告类型 (dau/wau/mau) - 仅用于active-users报告"
// @Success 200 {file} file
// @Router /api/v1/admin/analytics/export [get]
func (api *AnalyticsAPI) ExportAnalyticsReport(c *gin.Context) {
	reportType := c.Query("report_type")
	format := c.DefaultQuery("format", "csv")

	// 参数验证
	if reportType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "报告类型不能为空",
		})
		return
	}

	// 验证报告类型
	validTypes := []string{"user-growth", "content-statistics", "revenue-report", "active-users"}
	isValidType := false
	for _, t := range validTypes {
		if reportType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的报告类型",
		})
		return
	}

	// 验证导出格式
	if format != "csv" && format != "xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的导出格式",
		})
		return
	}

	// 根据报告类型生成不同的CSV内容
	var csvData []byte
	var filename string

	switch reportType {
	case "user-growth":
		filename = fmt.Sprintf("user_growth_%s.csv", time.Now().Format("20060102_150405"))
		csvData = []byte("Date,Count,GrowthRate\n")

	case "content-statistics":
		filename = fmt.Sprintf("content_statistics_%s.csv", time.Now().Format("20060102_150405"))
		csvData = []byte("Category,BookCount,ChapterCount,CommentCount,TotalWords\n")

	case "revenue-report":
		filename = fmt.Sprintf("revenue_report_%s.csv", time.Now().Format("20060102_150405"))
		csvData = []byte("Date,Amount,Orders,GrowthRate\n")

	case "active-users":
		filename = fmt.Sprintf("active_users_%s.csv", time.Now().Format("20060102_150405"))
		csvData = []byte("Date,Count,Type\n")

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的报告类型",
		})
		return
	}

	// 设置响应头
	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, "text/csv", csvData)
	} else if format == "xlsx" {
		// TODO: 实现Excel导出
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    501,
			"message": "Excel导出功能待实现",
		})
		return
	}
}

// GetAnalyticsDashboard 获取统计分析仪表板
// @Summary 获取统计分析仪表板
// @Description 返回统计分析仪表板的综合数据
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/dashboard [get]
func (api *AnalyticsAPI) GetAnalyticsDashboard(c *gin.Context) {
	// 并行获取各种统计数据
	type dashboardData struct {
		SystemOverview  *adminService.SystemOverviewResponse
		UserGrowth     *adminService.UserGrowthTrendResponse
		ContentStats   *adminService.ContentStatisticsResponse
		RevenueReport  *adminService.RevenueReportResponse
		ActiveUsers    *adminService.ActiveUsersReportResponse
	}

	// 获取系统概览
	systemOverview, err := api.analyticsService.GetSystemOverview(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取系统概览失败: %v", err),
		})
		return
	}

	// 获取用户增长趋势（最近7天）
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	userGrowth, _ := api.analyticsService.GetUserGrowthTrend(context.Background(), &adminService.UserGrowthTrendRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	})

	// 获取内容统计
	contentStats, _ := api.analyticsService.GetContentStatistics(context.Background(), &adminService.ContentStatisticsRequest{})

	// 获取收入报告（最近30天）
	startDate = endDate.AddDate(0, 0, -30)
	revenueReport, _ := api.analyticsService.GetRevenueReport(context.Background(), &adminService.RevenueReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "daily",
	})

	// 获取活跃用户报告（最近7天DAU）
	startDate = endDate.AddDate(0, 0, -7)
	activeUsers, _ := api.analyticsService.GetActiveUsersReport(context.Background(), &adminService.ActiveUsersReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      "dau",
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": dashboardData{
			SystemOverview: systemOverview,
			UserGrowth:     userGrowth,
			ContentStats:   contentStats,
			RevenueReport:  revenueReport,
			ActiveUsers:    activeUsers,
		},
	})
}

// GetCustomAnalyticsQuery 自定义统计分析查询
// @Summary 自定义统计分析查询
// @Description 根据自定义条件进行统计分析
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param query body CustomAnalyticsQuery true "自定义查询条件"
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/custom [post]
func (api *AnalyticsAPI) GetCustomAnalyticsQuery(c *gin.Context) {
	var req CustomAnalyticsQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("请求参数错误: %v", err),
		})
		return
	}

	// TODO: 实现自定义查询逻辑
	// 这里可以根据 req.Metrics, req.Dimensions, req.Filters 等参数构建查询

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"query": req,
			"result": "自定义查询结果",
		},
	})
}

// CustomAnalyticsQuery 自定义分析查询
type CustomAnalyticsQuery struct {
	Metrics   []string                 `json:"metrics" binding:"required"`   // 指标：users, books, revenue, etc.
	Dimensions []string                `json:"dimensions"`                   // 维度：date, category, etc.
	Filters   map[string]interface{}   `json:"filters"`                      // 过滤条件
	StartDate string                   `json:"start_date" binding:"required"`
	EndDate   string                   `json:"end_date" binding:"required"`
	Limit     int                      `json:"limit"`
}

// CompareAnalyticsPeriods 对比不同时期的统计数据
// @Summary 对比不同时期的统计数据
// @Description 对比两个时期的统计数据差异
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param period1_start query string true "第一期开始日期 (YYYY-MM-DD)"
// @Param period1_end query string true "第一期结束日期 (YYYY-MM-DD)"
// @Param period2_start query string true "第二期开始日期 (YYYY-MM-DD)"
// @Param period2_end query string true "第二期结束日期 (YYYY-MM-DD)"
// @Param metrics query string true "对比指标，逗号分隔 (users,revenue,etc)"
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/compare [get]
func (api *AnalyticsAPI) CompareAnalyticsPeriods(c *gin.Context) {
	period1StartStr := c.Query("period1_start")
	period1EndStr := c.Query("period1_end")
	period2StartStr := c.Query("period2_start")
	period2EndStr := c.Query("period2_end")
	metricsStr := c.Query("metrics")

	// 参数验证
	if period1StartStr == "" || period1EndStr == "" || period2StartStr == "" || period2EndStr == "" || metricsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要参数",
		})
		return
	}

	// 解析日期
	period1Start, err := time.Parse("2006-01-02", period1StartStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("第一期开始日期格式错误: %v", err),
		})
		return
	}

	period1End, err := time.Parse("2006-01-02", period1EndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("第一期结束日期格式错误: %v", err),
		})
		return
	}

	period2Start, err := time.Parse("2006-01-02", period2StartStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("第二期开始日期格式错误: %v", err),
		})
		return
	}

	period2End, err := time.Parse("2006-01-02", period2EndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("第二期结束日期格式错误: %v", err),
		})
		return
	}

	// TODO: 实现时期对比逻辑
	// 获取两个时期的数据并进行对比

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "对比成功",
		"data": gin.H{
			"period1": gin.H{
				"start": period1Start.Format("2006-01-02"),
				"end":   period1End.Format("2006-01-02"),
			},
			"period2": gin.H{
				"start": period2Start.Format("2006-01-02"),
				"end":   period2End.Format("2006-01-02"),
			},
			"metrics": metricsStr,
			"comparison": map[string]interface{}{
				"growth_rate": 12.5,
				"change":     1234,
			},
		},
	})
}

// GetRealTimeStats 获取实时统计数据
// @Summary 获取实时统计数据
// @Description 获取实时的系统统计数据
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/realtime [get]
func (api *AnalyticsAPI) GetRealTimeStats(c *gin.Context) {
	// TODO: 实现实时统计逻辑
	// 可以使用 Redis 等缓存存储实时数据

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"online_users":     1234,
			"active_sessions":  456,
			"requests_per_sec": 78,
			"avg_response_time": 45,
			"timestamp": time.Now().Unix(),
		},
	})
}

// GetAnalyticsPredictions 获取统计分析预测
// @Summary 获取统计分析预测
// @Description 基于历史数据预测未来趋势
// @Tags Admin-Analytics
// @Accept json
// @Produce json
// @Param type query string true "预测类型 (users/revenue/content)" Enums(users, revenue, content)
// @Param days query int false "预测天数" default(30)
// @Success 200 {object} Response
// @Router /api/v1/admin/analytics/predict [get]
func (api *AnalyticsAPI) GetAnalyticsPredictions(c *gin.Context) {
	predictType := c.Query("type")
	daysStr := c.DefaultQuery("days", "30")

	// 参数验证
	if predictType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "预测类型不能为空",
		})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "预测天数必须在1-365之间",
		})
		return
	}

	// TODO: 实现预测逻辑
	// 可以使用机器学习模型或简单的趋势分析

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"type":       predictType,
			"days":       days,
			"prediction": "预测数据",
			"confidence": 0.85,
		},
	})
}
