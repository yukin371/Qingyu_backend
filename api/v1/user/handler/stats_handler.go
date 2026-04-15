package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/user"
	"Qingyu_backend/service/writer"
)

// UserStatsService 切片 B 最小 seam：仅承接我的统计。
type UserStatsService interface {
	GetUserStats(ctx context.Context, userID string) (*user.UserStats, error)
}

// ContentStatsService 切片 C 最小 seam：仅承接我的内容统计。
type ContentStatsService interface {
	GetContentStats(ctx context.Context, userID string) (*writer.ContentStats, error)
}

// StatsHandler 用户统计处理器
type StatsHandler struct {
	userStatsService    UserStatsService
	contentStatsService ContentStatsService
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler(
	userStatsService UserStatsService,
	contentStatsService ContentStatsService,
) *StatsHandler {
	return &StatsHandler{
		userStatsService:    userStatsService,
		contentStatsService: contentStatsService,
	}
}

func respondStatsNotImplemented(c *gin.Context, message string) {
	response.JSON(c, http.StatusNotImplemented, response.APIResponse{
		Code:      response.CodeInternalError,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("requestId"),
	})
}

// GetMyStats 获取当前用户统计
//
//	@Summary		获取我的统计
//	@Description	获取当前登录用户的统计数据
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse{data=stats.UserStats}
//	@Failure		401	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/user/stats/my [get]
func (h *StatsHandler) GetMyStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取统计数据
	statsData, err := h.userStatsService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	// 3. 返回结果
	response.Success(c, statsData)
}

// GetMyContentStats 获取我的内容统计
//
//	@Summary		获取内容统计
//	@Description	获取当前用户的内容统计数据
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse{data=stats.ContentStats}
//	@Failure		401	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/user/stats/my/content [get]
func (h *StatsHandler) GetMyContentStats(c *gin.Context) {
	// 1. 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	userID := userIDInterface.(string)

	// 2. 获取内容统计
	contentStats, err := h.contentStatsService.GetContentStats(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	// 3. 返回结果
	response.Success(c, contentStats)
}

// GetMyActivityStats 获取我的活跃度统计
//
//	@Summary		获取活跃度统计
//	@Description	获取当前用户的活跃度统计（默认7天）
//	@Tags			用户管理-统计
//	@Accept			json
//	@Produce		json
//	@Param			days	query		int	false	"统计天数（默认7）"
//	@Success		200	{object}	response.APIResponse{data=stats.ActivityStats}
//	@Failure		401	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/user/stats/my/activity [get]
func (h *StatsHandler) GetMyActivityStats(c *gin.Context) {
	// 1. 获取当前用户ID
	_, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	respondStatsNotImplemented(c, "活跃度统计功能开发中")
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
//	@Success		200			{object}	response.APIResponse{data=stats.RevenueStats}
//	@Failure		401			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/user/stats/my/revenue [get]
func (h *StatsHandler) GetMyRevenueStats(c *gin.Context) {
	// 1. 获取当前用户ID
	_, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	respondStatsNotImplemented(c, "收益统计功能开发中")
}
