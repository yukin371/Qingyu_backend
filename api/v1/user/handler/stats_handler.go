package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	shared "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/shared/stats"
	"Qingyu_backend/pkg/response"
)

// StatsHandler 用户统计处理器
type StatsHandler struct {
	statsService stats.PlatformStatsService
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler(statsService stats.PlatformStatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetMyStats 获取当前用户统计
//
//	@Summary		获取我的统计
//	@Description	获取当前登录用户的统计数据
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=stats.UserStats}
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/user/stats/my [get]
func (h *StatsHandler) GetMyStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取统计数据
	statsData, err := h.statsService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 3. 返回结果
	shared.Success(c, http.StatusOK, "获取成功", statsData)
}

// GetMyContentStats 获取我的内容统计
//
//	@Summary		获取内容统计
//	@Description	获取当前用户的内容统计数据
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse{data=stats.ContentStats}
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/user/stats/my/content [get]
func (h *StatsHandler) GetMyContentStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取内容统计
	contentStats, err := h.statsService.GetContentStats(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 3. 返回结果
	shared.Success(c, http.StatusOK, "获取成功", contentStats)
}

// GetMyActivityStats 获取我的活跃度统计
//
//	@Summary		获取活跃度统计
//	@Description	获取当前用户的活跃度统计（默认7天）
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Param			days	query		int	false	"统计天数（默认7）"
//	@Success		200	{object}	shared.APIResponse{data=stats.ActivityStats}
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/user/stats/my/activity [get]
func (h *StatsHandler) GetMyActivityStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取天数参数
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	// 3. 获取活跃度统计
	activityStats, err := h.statsService.GetUserActivityStats(c.Request.Context(), userID, days)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 4. 返回结果
	shared.Success(c, http.StatusOK, "获取成功", activityStats)
}

// GetMyRevenueStats 获取我的收益统计
//
//	@Summary		获取收益统计
//	@Description	获取当前用户的收益统计
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"开始日期（YYYY-MM-DD）"
//	@Param			end_date	query		string	false	"结束日期（YYYY-MM-DD）"
//	@Success		200			{object}	shared.APIResponse{data=stats.RevenueStats}
//	@Failure		401			{object}	shared.ErrorResponse
//	@Failure		500			{object}	shared.ErrorResponse
//	@Router			/api/v1/user/stats/my/revenue [get]
func (h *StatsHandler) GetMyRevenueStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未登录")
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
	revenueStats, err := h.statsService.GetRevenueStats(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 4. 返回结果
	shared.Success(c, http.StatusOK, "获取成功", revenueStats)
}
