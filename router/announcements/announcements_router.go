package announcements

import (
	"github.com/gin-gonic/gin"

	announcementsAPI "Qingyu_backend/api/v1/announcements"
	"Qingyu_backend/internal/middleware/builtin"
	"Qingyu_backend/internal/middleware/ratelimit"
	messagingService "Qingyu_backend/service/messaging"
)

// RegisterAnnouncementRoutes 注册公告模块路由
func RegisterAnnouncementRoutes(
	r *gin.RouterGroup,
	announcementService messagingService.AnnouncementService,
) {
	// 应用中间件
	// 使用新架构的中间件
	r.Use(builtin.NewRequestIDMiddleware().Handler())
	r.Use(builtin.NewRecoveryMiddleware(nil).Handler())
	r.Use(builtin.NewCORSMiddleware().Handler())

	// 创建API处理器
	announcementPublicAPI := announcementsAPI.NewAnnouncementPublicAPI(announcementService)

	// ============ 公告公开路由 ============
	announcementGroup := r.Group("/announcements")
	{
		// 公开路由（无需认证，可选速率限制）
		publicAnnouncement := announcementGroup.Group("")
		publicAnnouncement.Use(ratelimit.RateLimitMiddlewareSimple(30, 60)) // 30次/分钟
		{
			publicAnnouncement.GET("/effective", announcementPublicAPI.GetEffectiveAnnouncements)
			publicAnnouncement.GET("/:id", announcementPublicAPI.GetAnnouncementByID)
			publicAnnouncement.POST("/:id/view", announcementPublicAPI.IncrementViewCount)
		}
	}
}
