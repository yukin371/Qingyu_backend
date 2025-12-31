package messaging

import (
	"github.com/gin-gonic/gin"

	messagingAPI "Qingyu_backend/api/v1/messaging"
	messagingService "Qingyu_backend/service/messaging"
	"Qingyu_backend/middleware"
)

// RegisterRoutes 注册messaging模块路由
func RegisterRoutes(
	r *gin.RouterGroup,
	announcementService messagingService.AnnouncementService,
) {
	// 应用中间件
	r.Use(middleware.ResponseFormatterMiddleware())
	r.Use(middleware.ResponseTimingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.Recovery())

	// 创建API处理器
	announcementPublicAPI := messagingAPI.NewAnnouncementPublicAPI(announcementService)

	// ============ 公告公开路由 ============
	announcementGroup := r.Group("/announcements")
	{
		// 公开路由（无需认证，可选速率限制）
		publicAnnouncement := announcementGroup.Group("")
		publicAnnouncement.Use(middleware.RateLimitMiddleware(30, 60)) // 30次/分钟
		{
			publicAnnouncement.GET("/effective", announcementPublicAPI.GetEffectiveAnnouncements)
			publicAnnouncement.GET("/:id", announcementPublicAPI.GetAnnouncementByID)
			publicAnnouncement.POST("/:id/view", announcementPublicAPI.IncrementViewCount)
		}
	}
}
