package shared

import (
	"Qingyu_backend/service/shared/stats"
	"time"

	"github.com/gin-gonic/gin"
)

// StatsAPI 统计API处理器
type StatsAPI struct {
	statsService stats.StatsService
}

// NewStatsAPI 创建统计API实例
func NewStatsAPI(statsService stats.StatsService) *StatsAPI {
	return &StatsAPI{
		statsService: statsService,
	}
}

// ============ 用户统计API ============

// GetMyStats 获取当前用户统计
//
//	@Summary		获取我的统计
//	@Description	获取当前登录用户的统计数据
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=stats.UserStats}
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/stats/my [get]
func (api *StatsAPI) GetMyStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取统计数据
	stats, err := api.statsService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		Error(c, 500, "获取统计失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, 200, "获取成功", stats)
}

// GetMyContentStats 获取我的内容统计
//
//	@Summary		获取内容统计
//	@Description	获取当前用户的内容统计数据
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=stats.ContentStats}
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/stats/my/content [get]
func (api *StatsAPI) GetMyContentStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取内容统计
	contentStats, err := api.statsService.GetContentStats(c.Request.Context(), userID)
	if err != nil {
		Error(c, 500, "获取内容统计失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, 200, "获取成功", contentStats)
}

// GetMyActivityStats 获取我的活跃度统计
//
//	@Summary		获取活跃度统计
//	@Description	获取当前用户的活跃度统计（默认7天）
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Param			days	query		int	false	"统计天数（默认7）"
//	@Success		200		{object}	APIResponse{data=stats.ActivityStats}
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/stats/my/activity [get]
func (api *StatsAPI) GetMyActivityStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取天数参数
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := time.ParseDuration(daysParam + "h"); err == nil {
			days = int(d.Hours() / 24)
		}
	}

	// 3. 获取活跃度统计
	activityStats, err := api.statsService.GetUserActivityStats(c.Request.Context(), userID, days)
	if err != nil {
		Error(c, 500, "获取活跃度统计失败", err.Error())
		return
	}

	// 4. 返回结果
	Success(c, 200, "获取成功", activityStats)
}

// GetMyRevenueStats 获取我的收益统计
//
//	@Summary		获取收益统计
//	@Description	获取当前用户的收益统计
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"开始日期（YYYY-MM-DD）"
//	@Param			end_date	query		string	false	"结束日期（YYYY-MM-DD）"
//	@Success		200			{object}	APIResponse{data=stats.RevenueStats}
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/stats/my/revenue [get]
func (api *StatsAPI) GetMyRevenueStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 解析日期参数
	var startDate, endDate time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = t
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // 默认最近一个月
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = t
		}
	} else {
		endDate = time.Now()
	}

	// 3. 获取收益统计
	revenueStats, err := api.statsService.GetRevenueStats(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		Error(c, 500, "获取收益统计失败", err.Error())
		return
	}

	// 4. 返回结果
	Success(c, 200, "获取成功", revenueStats)
}

// ============ 平台统计API（管理员）============

// GetPlatformUserStats 获取平台用户统计（管理员）
//
//	@Summary		获取平台用户统计
//	@Description	获取平台级别的用户统计数据（仅管理员）
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"开始日期"
//	@Param			end_date	query		string	false	"结束日期"
//	@Success		200			{object}	APIResponse{data=stats.PlatformUserStats}
//	@Failure		401			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/admin/stats/users [get]
func (api *StatsAPI) GetPlatformUserStats(c *gin.Context) {
	// 1. 解析日期参数
	startDate := time.Now().AddDate(0, -1, 0) // 默认最近一个月
	endDate := time.Now()

	// 2. 获取平台用户统计
	platformStats, err := api.statsService.GetPlatformUserStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		Error(c, 500, "获取平台用户统计失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, 200, "获取成功", platformStats)
}

// GetPlatformContentStats 获取平台内容统计（管理员）
//
//	@Summary		获取平台内容统计
//	@Description	获取平台级别的内容统计数据（仅管理员）
//	@Tags			统计
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"开始日期"
//	@Param			end_date	query		string	false	"结束日期"
//	@Success		200			{object}	APIResponse{data=stats.PlatformContentStats}
//	@Failure		401			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/admin/stats/content [get]
func (api *StatsAPI) GetPlatformContentStats(c *gin.Context) {
	// 1. 解析日期参数
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	// 2. 获取平台内容统计
	contentStats, err := api.statsService.GetPlatformContentStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		Error(c, 500, "获取平台内容统计失败", err.Error())
		return
	}

	// 3. 返回结果
	Success(c, 200, "获取成功", contentStats)
}

// TODO(Phase3): 高级统计API
// - [ ] 导出统计报表（Excel/PDF）
// - [ ] 自定义时间范围统计
// - [ ] 实时统计数据
// - [ ] 统计图表数据API
// - [ ] 对比分析API
