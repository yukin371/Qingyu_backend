package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/writer"
)

// DashboardApi 仪表板统计API
type DashboardApi struct {
	dashboardService *writer.DashboardService
}

// NewDashboardApi 创建仪表板统计API
func NewDashboardApi(dashboardService *writer.DashboardService) *DashboardApi {
	return &DashboardApi{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats 获取作者仪表板统计
// @Summary 获取作者仪表板统计数据
// @Description 获取作者工作台的汇总统计，包括总字数、项目数、今日字数、待审核数、连续写作天数
// @Tags Writer-Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=writer.DashboardStats}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/dashboard/stats [get]
func (api *DashboardApi) GetDashboardStats(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	stats, err := api.dashboardService.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, stats)
}
