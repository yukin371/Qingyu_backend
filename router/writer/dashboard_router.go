package writer

import (
	writerApi "Qingyu_backend/api/v1/writer"

	"github.com/gin-gonic/gin"
)

// InitDashboardRouter 初始化仪表板统计路由
func InitDashboardRouter(r *gin.RouterGroup, dashboardApi *writerApi.DashboardApi) {
	dashboardGroup := r.Group("/dashboard")
	{
		dashboardGroup.GET("/stats", dashboardApi.GetDashboardStats)
	}
}
