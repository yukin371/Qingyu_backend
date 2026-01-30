package notifications

import (
	"github.com/gin-gonic/gin"

	notificationsAPI "Qingyu_backend/api/v1/notifications"
	"Qingyu_backend/internal/middleware/auth"
)

// RegisterRoutes 注册通知路由
func RegisterRoutes(r *gin.RouterGroup, notificationAPI *notificationsAPI.NotificationAPI) {
	// 通知路由组（需要认证）
	notificationGroup := r.Group("/notifications")
	notificationGroup.Use(auth.JWTAuth())

	// 通知基础操作
	notificationGroup.GET("", notificationAPI.GetNotifications)
	notificationGroup.GET("/unread-count", notificationAPI.GetUnreadCount)
	notificationGroup.GET("/stats", notificationAPI.GetNotificationStats)
	notificationGroup.GET("/:id", notificationAPI.GetNotification)

	// 标记已读
	notificationGroup.POST("/:id/read", notificationAPI.MarkAsRead)
	notificationGroup.POST("/batch-read", notificationAPI.MarkMultipleAsRead)
	notificationGroup.POST("/read-all", notificationAPI.MarkAllAsRead)

	// 清除已读通知
	notificationGroup.POST("/clear-read", notificationAPI.ClearReadNotifications)

	// 重新发送通知
	notificationGroup.POST("/:id/resend", notificationAPI.ResendNotification)

	// WebSocket端点
	notificationGroup.GET("/ws-endpoint", notificationAPI.GetWSEndpoint)

	// 删除通知
	notificationGroup.DELETE("/:id", notificationAPI.DeleteNotification)
	notificationGroup.POST("/batch-delete", notificationAPI.BatchDeleteNotifications)
	notificationGroup.DELETE("/delete-all", notificationAPI.DeleteAllNotifications)

	// 通知偏好设置
	notificationGroup.GET("/preferences", notificationAPI.GetNotificationPreference)
	notificationGroup.PUT("/preferences", notificationAPI.UpdateNotificationPreference)
	notificationGroup.POST("/preferences/reset", notificationAPI.ResetNotificationPreference)

	// 推送设备管理
	notificationGroup.POST("/push/register", notificationAPI.RegisterPushDevice)
	notificationGroup.DELETE("/push/unregister/:deviceId", notificationAPI.UnregisterPushDevice)
	notificationGroup.GET("/push/devices", notificationAPI.GetPushDevices)

	// 邮件通知设置（在user-management路径下）
	notificationGroup.GET("/email-settings", notificationAPI.GetEmailNotificationSettings)
	notificationGroup.PUT("/email-settings", notificationAPI.UpdateEmailNotificationSettings)

	// 短信通知设置（在user-management路径下）
	notificationGroup.GET("/sms-settings", notificationAPI.GetSMSNotificationSettings)
	notificationGroup.PUT("/sms-settings", notificationAPI.UpdateSMSNotificationSettings)
}

// RegisterUserManagementRoutes 注册用户管理路由（包含邮件和短信通知设置）
func RegisterUserManagementRoutes(r *gin.RouterGroup, notificationAPI *notificationsAPI.NotificationAPI) {
	// 用户管理路由组（需要认证）
	userManagementGroup := r.Group("/user-management")
	userManagementGroup.Use(auth.JWTAuth())

	// 邮件通知设置
	userManagementGroup.GET("/email-notifications", notificationAPI.GetEmailNotificationSettings)
	userManagementGroup.PUT("/email-notifications", notificationAPI.UpdateEmailNotificationSettings)

	// 短信通知设置
	userManagementGroup.GET("/sms-notifications", notificationAPI.GetSMSNotificationSettings)
	userManagementGroup.PUT("/sms-notifications", notificationAPI.UpdateSMSNotificationSettings)
}
