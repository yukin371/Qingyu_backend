package writer

import (
	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/writer"
)

// InitDashboardRouter 初始化仪表板统计路由
func InitDashboardRouter(r *gin.RouterGroup, dashboardApi *writer.DashboardApi) {
	dashboardGroup := r.Group("/dashboard")
	{
		dashboardGroup.GET("/stats", dashboardApi.GetDashboardStats)
	}
}
