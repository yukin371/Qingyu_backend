package admin

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/admin"
	"Qingyu_backend/internal/middleware/auth"
	adminservice "Qingyu_backend/service/admin"
	aiService "Qingyu_backend/service/ai"
	bookstoreService "Qingyu_backend/service/bookstore"
	eventservice "Qingyu_backend/service/events"
	publishService "Qingyu_backend/service/interfaces"
	auditService "Qingyu_backend/service/interfaces/audit"
	userService "Qingyu_backend/service/interfaces/user"
	messagingService "Qingyu_backend/service/messaging"
	sharedService "Qingyu_backend/service/shared"
)

// RegisterAdminRoutes 注册管理员路由
func RegisterAdminRoutes(
	r *gin.RouterGroup,
	userSvc userService.UserService,
	quotaSvc *aiService.QuotaService,
	auditSvc auditService.ContentAuditService,
	adminSvc adminservice.AdminService,
	configSvc *sharedService.ConfigService,
	announcementSvc messagingService.AnnouncementService,
	userAdminSvc adminservice.UserAdminService,
	permissionSvc sharedService.PermissionService,
	eventBus eventservice.PersistedEventBusInterface,
	categorySvc adminservice.CategoryAdminService,
	publicationSvc publishService.PublishService,
	bannerSvc bookstoreService.BannerService,
	quotaAdminSvc *aiService.QuotaAdminService,
	quotaDashboardSvc *aiService.QuotaDashboardService,
	quotaPolicySvc *aiService.QuotaPolicyService,
	quotaAlertSvc *aiService.QuotaAlertService,
) {
	// 创建admin API实例
	auditAdminAPI := admin.NewAuditAdminAPI(auditSvc)
	systemAdminAPI := admin.NewSystemAdminAPI(adminSvc)
	configAdminAPI := admin.NewConfigAPI(configSvc)
	announcementAdminAPI := admin.NewAnnouncementAPI(announcementSvc)
	publicationAdminAPI := admin.NewPublicationAdminAPI(publicationSvc)

	var quotaAdminAPI *admin.QuotaAdminAPI
	var quotaDashboardAPI *admin.QuotaDashboardAPI
	var quotaPolicyAPI *admin.QuotaPolicyAPI
	var quotaAlertAPI *admin.QuotaAlertAPI
	if quotaSvc != nil && quotaAdminSvc != nil && quotaDashboardSvc != nil && quotaPolicySvc != nil && quotaAlertSvc != nil {
		quotaAdminAPI = admin.NewQuotaAdminAPI(quotaSvc, quotaAdminSvc)
		quotaDashboardAPI = admin.NewQuotaDashboardAPI(quotaDashboardSvc)
		quotaPolicyAPI = admin.NewQuotaPolicyAPI(quotaPolicySvc)
		quotaAlertAPI = admin.NewQuotaAlertAPI(quotaAlertSvc)
	}

	// 分类管理API
	var categoryAdminAPI *admin.CategoryAdminAPI
	if categorySvc != nil {
		categoryAdminAPI = admin.NewCategoryAdminAPI(categorySvc)
	}

	// 权限管理API
	var permissionAPI *admin.PermissionAPI
	if permissionSvc != nil {
		permissionAPI = admin.NewPermissionAPI(permissionSvc)
	}

	// 事件管理API
	var eventsAPI *admin.EventsAdminAPI
	if eventBus != nil {
		eventsAPI = admin.NewEventsAdminAPI(eventBus)
	}

	// 管理员路由组 - 需要JWT认证 + 管理员权限
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuth())            // JWT认证
	adminGroup.Use(auth.RequireRole("admin")) // 管理员权限验证
	{
		// ===========================
		// 用户管理（管理员专用）
		// ===========================
		if userAdminSvc != nil {
			userAdminAPI := admin.NewUserAdminAPI(userAdminSvc)
			usersGroup := adminGroup.Group("/users")
			{
				// 用户列表和搜索
				usersGroup.GET("", userAdminAPI.ListUsers)                     // 获取用户列表
				usersGroup.POST("", userAdminAPI.CreateUser)                   // 创建用户
				usersGroup.GET("/search", userAdminAPI.SearchUsers)            // 搜索用户
				usersGroup.GET("/count-by-status", userAdminAPI.CountByStatus) // 按状态统计

				// 批量操作
				usersGroup.POST("/batch-create", userAdminAPI.BatchCreateUsers)         // 批量创建
				usersGroup.POST("/batch-update-status", userAdminAPI.BatchUpdateStatus) // 批量更新状态
				usersGroup.POST("/batch-delete", userAdminAPI.BatchDeleteUsers)         // 批量删除

				// 单个用户操作
				usersGroup.GET("/:id", userAdminAPI.GetUserDetail)                     // 获取用户详情
				usersGroup.DELETE("/:id", userAdminAPI.DeleteUser)                     // 删除用户
				usersGroup.PUT("/:id/status", userAdminAPI.UpdateUserStatus)           // 更新用户状态
				usersGroup.PUT("/:id/role", userAdminAPI.UpdateUserRole)               // 更新用户角色
				usersGroup.POST("/:id/reset-password", userAdminAPI.ResetUserPassword) // 重置密码
				usersGroup.GET("/:id/activities", userAdminAPI.GetUserActivities)      // 获取活动记录
				usersGroup.GET("/:id/statistics", userAdminAPI.GetUserStatistics)      // 获取统计信息
			}
		}

		// ===========================
		// AI配额管理
		// ===========================
		if quotaSvc != nil && quotaAdminAPI != nil && quotaDashboardAPI != nil && quotaPolicyAPI != nil && quotaAlertAPI != nil {
			quotaGroup := adminGroup.Group("/quota")
			{
				// 仪表盘
				quotaGroup.GET("/dashboard", quotaDashboardAPI.GetDashboard)
				quotaGroup.GET("/statistics/global", quotaDashboardAPI.GetStatistics)
				quotaGroup.GET("/statistics/trend", quotaDashboardAPI.GetTrend)
				quotaGroup.POST("/dashboard/refresh", quotaDashboardAPI.RefreshCache)

				// 用户配额管理
				quotaGroup.GET("/users", quotaAdminAPI.ListUserQuotas)
				quotaGroup.GET("/users/:userId", quotaAdminAPI.GetUserQuotaDetails)
				quotaGroup.PUT("/users/:userId", quotaAdminAPI.UpdateUserQuota)
				quotaGroup.POST("/users/:userId/recharge", quotaAdminAPI.RechargeUserQuota)
				quotaGroup.POST("/users/:userId/suspend", quotaAdminAPI.SuspendUserQuota)
				quotaGroup.POST("/users/:userId/activate", quotaAdminAPI.ActivateUserQuota)

				// 批量操作
				quotaGroup.POST("/batch-recharge", quotaAdminAPI.BatchRecharge)
				quotaGroup.POST("/batch-update", quotaAdminAPI.BatchUpdateQuota)
				quotaGroup.POST("/batch-suspend", quotaAdminAPI.BatchSuspend)
				quotaGroup.POST("/batch-activate", quotaAdminAPI.BatchActivate)

				// 策略管理
				quotaGroup.GET("/policies", quotaPolicyAPI.ListPolicies)
				quotaGroup.GET("/policies/:id", quotaPolicyAPI.GetPolicy)
				quotaGroup.POST("/policies", quotaPolicyAPI.CreatePolicy)
				quotaGroup.PUT("/policies/:id", quotaPolicyAPI.UpdatePolicy)
				quotaGroup.DELETE("/policies/:id", quotaPolicyAPI.DeletePolicy)
				quotaGroup.POST("/policies/initialize", quotaPolicyAPI.InitializeDefaultPolicies)

				// 告警管理
				quotaGroup.GET("/alerts", quotaAlertAPI.ListAlerts)
				quotaGroup.GET("/alerts/:id", quotaAlertAPI.GetAlert)
				quotaGroup.PUT("/alerts/:id/acknowledge", quotaAlertAPI.AcknowledgeAlert)
				quotaGroup.PUT("/alerts/:id/resolve", quotaAlertAPI.ResolveAlert)
				quotaGroup.PUT("/alerts/:id/ignore", quotaAlertAPI.IgnoreAlert)
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

			adminGroup.GET("/stats/reviews", auditAdminAPI.GetAuditStatistics)
		}

		if publicationSvc != nil {
			publicationGroup := adminGroup.Group("/publications")
			{
				publicationGroup.GET("/pending", publicationAdminAPI.GetPendingPublications)
				publicationGroup.POST("/:id/review", publicationAdminAPI.ReviewPublication)
			}
		}

		// ===========================
		// 系统管理
		// ===========================

		// 系统统计
		if adminSvc != nil {
			adminGroup.GET("/stats", systemAdminAPI.GetSystemStats)
			adminGroup.GET("/stats/overview", systemAdminAPI.GetSystemStats)
			adminGroup.GET("/stats/users", systemAdminAPI.GetSystemStats)
			adminGroup.GET("/stats/content", systemAdminAPI.GetSystemStats)
			adminGroup.GET("/stats/revenue", systemAdminAPI.GetWithdrawalStatistics)
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
			adminGroup.GET("/withdrawals", systemAdminAPI.ListWithdrawals)
			adminGroup.GET("/withdrawals/stats", systemAdminAPI.GetWithdrawalStatistics)
			adminGroup.POST("/withdraw/review", systemAdminAPI.ReviewWithdraw)
		}

		// ===========================
		// 公告管理
		// ===========================
		if announcementSvc != nil {
			announcementsGroup := adminGroup.Group("/announcements")
			{
				announcementsGroup.GET("", announcementAdminAPI.GetAnnouncements)                // 获取公告列表
				announcementsGroup.GET("/:id", announcementAdminAPI.GetAnnouncementByID)         // 获取公告详情
				announcementsGroup.POST("", announcementAdminAPI.CreateAnnouncement)             // 创建公告
				announcementsGroup.PUT("/:id", announcementAdminAPI.UpdateAnnouncement)          // 更新公告
				announcementsGroup.DELETE("/:id", announcementAdminAPI.DeleteAnnouncement)       // 删除公告
				announcementsGroup.PUT("/batch-status", announcementAdminAPI.BatchUpdateStatus)  // 批量更新状态（OpenAPI 契约）
				announcementsGroup.POST("/batch-status", announcementAdminAPI.BatchUpdateStatus) // 批量更新状态
				announcementsGroup.DELETE("/batch-delete", announcementAdminAPI.BatchDelete)     // 批量删除
			}
		}

		// ===========================
		// 权限和角色管理
		// ===========================
		if permissionAPI != nil {
			// 权限管理
			permissionsGroup := adminGroup.Group("/permissions")
			{
				permissionsGroup.GET("", permissionAPI.GetAllPermissions)         // 获取所有权限
				permissionsGroup.GET("/:code", permissionAPI.GetPermission)       // 获取权限详情
				permissionsGroup.POST("", permissionAPI.CreatePermission)         // 创建权限
				permissionsGroup.PUT("/:code", permissionAPI.UpdatePermission)    // 更新权限
				permissionsGroup.DELETE("/:code", permissionAPI.DeletePermission) // 删除权限
			}

			// 角色管理
			rolesGroup := adminGroup.Group("/roles")
			{
				rolesGroup.GET("", permissionAPI.GetAllRoles)                                                 // 获取所有角色
				rolesGroup.GET("/:id", permissionAPI.GetRole)                                                 // 获取角色详情
				rolesGroup.POST("", permissionAPI.CreateRole)                                                 // 创建角色
				rolesGroup.PUT("/:id", permissionAPI.UpdateRole)                                              // 更新角色
				rolesGroup.DELETE("/:id", permissionAPI.DeleteRole)                                           // 删除角色
				rolesGroup.GET("/:id/permissions", permissionAPI.GetRolePermissions)                          // 获取角色权限
				rolesGroup.POST("/:id/permissions/:permissionCode", permissionAPI.AssignPermissionToRole)     // 为角色分配权限
				rolesGroup.DELETE("/:id/permissions/:permissionCode", permissionAPI.RemovePermissionFromRole) // 移除角色权限
			}
		}

		// ===========================
		// 事件管理
		// ===========================
		if eventsAPI != nil {
			eventsGroup := adminGroup.Group("/events")
			{
				eventsGroup.POST("/replay", eventsAPI.ReplayEvents) // 事件回放
			}
		}

		// ===========================
		// 分类管理
		// ===========================
		if categorySvc != nil {
			categoriesGroup := adminGroup.Group("/categories")
			{
				categoriesGroup.GET("", categoryAdminAPI.GetCategories)
				categoriesGroup.GET("/tree", categoryAdminAPI.GetCategoryTree)
				categoriesGroup.GET("/:id", categoryAdminAPI.GetCategoryByID)
				categoriesGroup.POST("", categoryAdminAPI.CreateCategory)
				categoriesGroup.PUT("/:id", categoryAdminAPI.UpdateCategory)
				categoriesGroup.DELETE("/:id", categoryAdminAPI.DeleteCategory)
				categoriesGroup.PUT("/:id/move", categoryAdminAPI.MoveCategory)
				categoriesGroup.PUT("/:id/sort", categoryAdminAPI.SortCategory)
			}
		}

		// ===========================
		// Banner管理
		// ===========================
		if bannerSvc != nil {
			bannerAdminAPI := admin.NewBannerAPI(bannerSvc)
			bannersGroup := adminGroup.Group("/banners")
			{
				bannersGroup.GET("", bannerAdminAPI.GetBanners)                     // 获取Banner列表
				bannersGroup.GET("/:id", bannerAdminAPI.GetBannerByID)              // 获取Banner详情
				bannersGroup.POST("", bannerAdminAPI.CreateBanner)                  // 创建Banner
				bannersGroup.PUT("/:id", bannerAdminAPI.UpdateBanner)               // 更新Banner
				bannersGroup.DELETE("/:id", bannerAdminAPI.DeleteBanner)            // 删除Banner
				bannersGroup.PUT("/batch-status", bannerAdminAPI.BatchUpdateStatus) // 批量更新状态
				bannersGroup.PUT("/batch-sort", bannerAdminAPI.BatchUpdateSort)     // 批量更新排序
			}
		}
	}
}
