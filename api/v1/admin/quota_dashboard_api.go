package admin

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	ai "Qingyu_backend/service/ai"
)

// QuotaDashboardAPI 配额仪表盘API处理器
type QuotaDashboardAPI struct {
	dashboardService *ai.QuotaDashboardService
}

// NewQuotaDashboardAPI 创建配额仪表盘API实例
func NewQuotaDashboardAPI(dashboardService *ai.QuotaDashboardService) *QuotaDashboardAPI {
	return &QuotaDashboardAPI{
		dashboardService: dashboardService,
	}
}

// GetDashboard 获取配额仪表盘数据
//
//	@Summary		获取配额仪表盘
//	@Description	管理员获取配额仪表盘完整数据
//	@Tags			管理员-配额仪表盘
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/dashboard [get]
func (api *QuotaDashboardAPI) GetDashboard(c *gin.Context) {
	dashboard, err := api.dashboardService.GetDashboard(c.Request.Context())
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取仪表盘数据失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", dashboard)
}

// GetStatistics 获取全局统计
//
//	@Summary		获取全局配额统计
//	@Description	管理员获取配额全局统计数据
//	@Tags			管理员-配额仪表盘
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/statistics/global [get]
func (api *QuotaDashboardAPI) GetStatistics(c *gin.Context) {
	summary, err := api.dashboardService.GetDashboardSummary(c.Request.Context())
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取统计数据失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", summary)
}

// GetTrend 获取消耗趋势
//
//	@Summary		获取消耗趋势
//	@Description	管理员获取配额消耗趋势数据
//	@Tags			管理员-配额仪表盘
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			days	query		int	false	"天数(默认7天)"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/statistics/trend [get]
func (api *QuotaDashboardAPI) GetTrend(c *gin.Context) {
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	trend, err := api.dashboardService.GetConsumptionTrend(c.Request.Context(), days)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取趋势数据失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", trend)
}

// GetReconciliationSummary 获取跨服务聚合对账摘要
//
//	@Summary		获取跨服务聚合对账摘要
//	@Description	管理员获取 backend 与 AI service 的聚合配额消费对账结果
//	@Tags			管理员-配额仪表盘
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			timeRange		query		string	false	"时间范围(day/week/month/all，默认day)"
//	@Param			workflowType	query		string	false	"工作流过滤"
//	@Param			groupBy			query		string	false	"聚合维度(user/workflow，默认user)"
//	@Param			page			query		int		false	"页码"
//	@Param			pageSize		query		int		false	"每页数量"
//	@Success		200				{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/statistics/reconciliation [get]
func (api *QuotaDashboardAPI) GetReconciliationSummary(c *gin.Context) {
	page := 1
	pageSize := 20
	if pageStr := c.Query("page"); pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if parsed, err := strconv.Atoi(pageSizeStr); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	summary, err := api.dashboardService.GetReconciliationSummary(
		c.Request.Context(),
		c.DefaultQuery("timeRange", "day"),
		c.Query("workflowType"),
		c.DefaultQuery("groupBy", "user"),
		page,
		pageSize,
	)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取对账摘要失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", summary)
}

// RefreshCache 刷新仪表盘缓存
//
//	@Summary		刷新仪表盘缓存
//	@Description	管理员强制刷新配额仪表盘缓存
//	@Tags			管理员-配额仪表盘
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/dashboard/refresh [post]
func (api *QuotaDashboardAPI) RefreshCache(c *gin.Context) {
	if err := api.dashboardService.RefreshDashboardCache(c.Request.Context()); err != nil {
		response.InternalError(c, fmt.Errorf("刷新缓存失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "刷新成功", nil)
}

// RunConsistencyCheck 手动执行跨服务一致性检查
//
//	@Summary		执行跨服务一致性检查
//	@Description	管理员手动触发一次 backend 与 AI service 的配额对账检查
//	@Tags			管理员-配额仪表盘
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/statistics/reconciliation/check [post]
func (api *QuotaDashboardAPI) RunConsistencyCheck(c *gin.Context) {
	if err := api.dashboardService.RunConsistencyCheck(c.Request.Context()); err != nil {
		response.InternalError(c, fmt.Errorf("执行一致性检查失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "一致性检查已执行", nil)
}
