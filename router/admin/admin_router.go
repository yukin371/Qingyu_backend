package admin

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/middleware"
	aiService "Qingyu_backend/service/ai"
	auditService "Qingyu_backend/service/interfaces/audit"
	messagingService "Qingyu_backend/service/messaging"
	userService "Qingyu_backend/service/interfaces/user"
	sharedService "Qingyu_backend/service/shared"
	adminService "Qingyu_backend/service/shared/admin"
	"Qingyu_backend/api/v1/admin"
)

// RegisterAdminRoutes 注册管理员路由
func RegisterAdminRoutes(
	r *gin.RouterGroup,
	userSvc userService.UserService,
	quotaSvc *aiService.QuotaService,
	auditSvc auditService.ContentAuditService,
	adminSvc adminService.AdminService,
	configSvc *sharedService.ConfigService,
	announcementSvc messagingService.AnnouncementService,
) {
	// 创建admin API实例
	quotaAdminAPI := admin.NewQuotaAdminAPI(quotaSvc)
	auditAdminAPI := admin.NewAuditAdminAPI(auditSvc)
	systemAdminAPI := admin.NewSystemAdminAPI(adminSvc)
	configAdminAPI := admin.NewConfigAPI(configSvc)
	announcementAdminAPI := admin.NewAnnouncementAPI(announcementSvc)

	// 管理员路由组 - 需要JWT认证 + 管理员权限
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JWTAuth())            // JWT认证
	adminGroup.Use(middleware.RequireRole("admin")) // 管理员权限验证
	{
		// ===========================
		// 用户管理
		// ⚠️ 注意：用户管理功能已迁移到 user-management 模块
		// 新路由位置：/api/v1/user-management/users/* (需管理员权限)
		// 此处保留注释说明，实际路由已在 router/enter.go 中注册
		// ===========================

		// ===========================
		// AI配额管理
		// ===========================
		if quotaSvc != nil {
			quotaGroup := adminGroup.Group("/quota")
			{
				quotaGroup.GET("/:userId", quotaAdminAPI.GetUserQuotaDetails)         // 获取用户配额详情
				quotaGroup.PUT("/:userId", quotaAdminAPI.UpdateUserQuota)             // 更新用户配额
				quotaGroup.POST("/:userId/suspend", quotaAdminAPI.SuspendUserQuota)   // 暂停用户配额
				quotaGroup.POST("/:userId/activate", quotaAdminAPI.ActivateUserQuota) // 激活用户配额
			}
		}

		// ===========================
		// 审核管理
		// ===========================
		if auditSvc != nil {
			auditGroup := adminGroup.Group("/audit")
			{
				auditGroup.GET("/pending", auditAdminAPI.GetPendingAudits)        // 获取待审核内容
				auditGroup.GET("/high-risk", auditAdminAPI.GetHighRiskAudits)     // 获取高风险审核记录
				auditGroup.GET("/statistics", auditAdminAPI.GetAuditStatistics)   // 获取审核统计
				auditGroup.POST("/:id/review", auditAdminAPI.ReviewAudit)         // 审核内容
				auditGroup.POST("/:id/appeal/review", auditAdminAPI.ReviewAppeal) // 审核申诉
			}
		}

		// ===========================
		// 系统管理
		// ===========================

		// 系统统计
		if adminSvc != nil {
			adminGroup.GET("/stats", systemAdminAPI.GetSystemStats)
			adminGroup.GET("/operation-logs", systemAdminAPI.GetOperationLogs)
		}

		// ===========================
		// 配置管理（新版）
		// ===========================
		if configSvc != nil {
			configGroup := adminGroup.Group("/config")
			{
				// 读取配置
				configGroup.GET("", configAdminAPI.GetAllConfigs)       // 获取所有配置
				configGroup.GET("/:key", configAdminAPI.GetConfigByKey) // 获取单个配置

				// 更新配置
				configGroup.PUT("", configAdminAPI.UpdateConfig)            // 更新配置
				configGroup.PUT("/batch", configAdminAPI.BatchUpdateConfig) // 批量更新配置

				// 验证配置
				configGroup.POST("/validate", configAdminAPI.ValidateConfig) // 验证配置

				// 备份管理
				configGroup.GET("/backups", configAdminAPI.GetConfigBackups)     // 获取备份列表
				configGroup.POST("/restore", configAdminAPI.RestoreConfigBackup) // 恢复备份
			}

			// 系统配置（旧版，保持兼容）
			adminGroup.GET("/config-legacy", systemAdminAPI.GetSystemConfig)
			adminGroup.PUT("/config-legacy", systemAdminAPI.UpdateSystemConfig)
		}

		// 提现管理
		if adminSvc != nil {
			adminGroup.POST("/withdraw/review", systemAdminAPI.ReviewWithdraw)
		}

		// ===========================
		// 公告管理
		// ===========================
		if announcementSvc != nil {
			announcementsGroup := adminGroup.Group("/announcements")
			{
				announcementsGroup.GET("", announcementAdminAPI.GetAnnouncements)           // 获取公告列表
				announcementsGroup.GET("/:id", announcementAdminAPI.GetAnnouncementByID)    // 获取公告详情
				announcementsGroup.POST("", announcementAdminAPI.CreateAnnouncement)          // 创建公告
				announcementsGroup.PUT("/:id", announcementAdminAPI.UpdateAnnouncement)       // 更新公告
				announcementsGroup.DELETE("/:id", announcementAdminAPI.DeleteAnnouncement)    // 删除公告
				announcementsGroup.PUT("/batch-status", announcementAdminAPI.BatchUpdateStatus) // 批量更新状态
				announcementsGroup.DELETE("/batch-delete", announcementAdminAPI.BatchDelete)      // 批量删除
			}
		}
	}
}
