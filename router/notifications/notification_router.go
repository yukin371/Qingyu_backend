package notifications

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/notifications"
	"Qingyu_backend/middleware"
)

// RegisterRoutes 注册通知路由
func RegisterRoutes(r *gin.RouterGroup, notificationAPI *notifications.NotificationAPI) {
	// 通知路由组（需要认证）
	notificationGroup := r.Group("/notifications")
	notificationGroup.Use(middleware.AuthMiddleware())

	// 通知基础操作
	notificationGroup.GET("", notificationAPI.GetNotifications)
	notificationGroup.GET("/unread-count", notificationAPI.GetUnreadCount)
	notificationGroup.GET("/stats", notificationAPI.GetNotificationStats)
	notificationGroup.GET("/:id", notificationAPI.GetNotification)

	// 标记已读
	notificationGroup.PUT("/:id/read", notificationAPI.MarkAsRead)
	notificationGroup.PUT("/mark-read", notificationAPI.MarkMultipleAsRead)
	notificationGroup.PUT("/read-all", notificationAPI.MarkAllAsRead)

	// 删除通知
	notificationGroup.DELETE("/:id", notificationAPI.DeleteNotification)
	notificationGroup.DELETE("/batch-delete", notificationAPI.BatchDeleteNotifications)
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
func RegisterUserManagementRoutes(r *gin.RouterGroup, notificationAPI *notifications.NotificationAPI) {
	// 用户管理路由组（需要认证）
	userManagementGroup := r.Group("/user-management")
	userManagementGroup.Use(middleware.AuthMiddleware())

	// 邮件通知设置
	userManagementGroup.GET("/email-notifications", notificationAPI.GetEmailNotificationSettings)
	userManagementGroup.PUT("/email-notifications", notificationAPI.UpdateEmailNotificationSettings)

	// 短信通知设置
	userManagementGroup.GET("/sms-notifications", notificationAPI.GetSMSNotificationSettings)
	userManagementGroup.PUT("/sms-notifications", notificationAPI.UpdateSMSNotificationSettings)
}
