package admin

import (
	"github.com/gin-gonic/gin"

	adminApi "Qingyu_backend/api/v1/admin"
	"Qingyu_backend/middleware"
	aiService "Qingyu_backend/service/ai"
	auditService "Qingyu_backend/service/interfaces/audit"
	userService "Qingyu_backend/service/interfaces/user"
	sharedService "Qingyu_backend/service/shared"
	adminService "Qingyu_backend/service/shared/admin"
)

// RegisterAdminRoutes 注册管理员路由
func RegisterAdminRoutes(
	r *gin.RouterGroup,
	userSvc userService.UserService,
	quotaSvc *aiService.QuotaService,
	auditSvc auditService.ContentAuditService,
	adminSvc adminService.AdminService,
	configSvc *sharedService.ConfigService,
) {
	// 创建admin API实例
	userAdminAPI := adminApi.NewUserAdminAPI(userSvc)
	quotaAdminAPI := adminApi.NewQuotaAdminAPI(quotaSvc)
	auditAdminAPI := adminApi.NewAuditAdminAPI(auditSvc)
	systemAdminAPI := adminApi.NewSystemAdminAPI(adminSvc)
	configAdminAPI := adminApi.NewConfigAPI(configSvc)

	// 管理员路由组 - 需要JWT认证 + 管理员权限
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JWTAuth())            // JWT认证
	adminGroup.Use(middleware.RequireRole("admin")) // 管理员权限验证
	{
		// ===========================
		// 用户管理
		// ===========================
		usersGroup := adminGroup.Group("/users")
		{
			usersGroup.GET("", userAdminAPI.ListUsers)            // 获取用户列表
			usersGroup.GET("/:id", userAdminAPI.GetUser)          // 获取用户详情
			usersGroup.PUT("/:id", userAdminAPI.UpdateUser)       // 更新用户信息
			usersGroup.DELETE("/:id", userAdminAPI.DeleteUser)    // 删除用户
			usersGroup.POST("/:id/ban", userAdminAPI.BanUser)     // 封禁用户
			usersGroup.POST("/:id/unban", userAdminAPI.UnbanUser) // 解除封禁

			// 用户统计
			usersGroup.GET("/:id/statistics", systemAdminAPI.GetUserStatistics)
		}

		// ===========================
		// AI配额管理
		// ===========================
		quotaGroup := adminGroup.Group("/quota")
		{
			quotaGroup.GET("/:userId", quotaAdminAPI.GetUserQuotaDetails)         // 获取用户配额详情
			quotaGroup.PUT("/:userId", quotaAdminAPI.UpdateUserQuota)             // 更新用户配额
			quotaGroup.POST("/:userId/suspend", quotaAdminAPI.SuspendUserQuota)   // 暂停用户配额
			quotaGroup.POST("/:userId/activate", quotaAdminAPI.ActivateUserQuota) // 激活用户配额
		}

		// ===========================
		// 审核管理
		// ===========================
		auditGroup := adminGroup.Group("/audit")
		{
			auditGroup.GET("/pending", auditAdminAPI.GetPendingAudits)        // 获取待审核内容
			auditGroup.GET("/high-risk", auditAdminAPI.GetHighRiskAudits)     // 获取高风险审核记录
			auditGroup.GET("/statistics", auditAdminAPI.GetAuditStatistics)   // 获取审核统计
			auditGroup.POST("/:id/review", auditAdminAPI.ReviewAudit)         // 审核内容
			auditGroup.POST("/:id/appeal/review", auditAdminAPI.ReviewAppeal) // 审核申诉
		}

		// ===========================
		// 系统管理
		// ===========================

		// 系统统计
		adminGroup.GET("/stats", systemAdminAPI.GetSystemStats)

		// 操作日志
		adminGroup.GET("/operation-logs", systemAdminAPI.GetOperationLogs)

		// ===========================
		// 配置管理（新版）
		// ===========================
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

		// 提现管理
		adminGroup.POST("/withdraw/review", systemAdminAPI.ReviewWithdraw)

		// 公告管理
		adminGroup.GET("/announcements", systemAdminAPI.GetAnnouncements)
		adminGroup.POST("/announcements", systemAdminAPI.CreateAnnouncement)
	}
}
