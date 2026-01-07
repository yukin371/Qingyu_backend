package usermanagement

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/usermanagement/handler"
	managementApi "Qingyu_backend/api/v1/usermanagement"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/shared/stats"
)

// RegisterUsermanagementRoutes 注册用户管理路由
func RegisterUsermanagementRoutes(
	r *gin.RouterGroup,
	userService userServiceInterface.UserService,
	bookstoreService interface{},
	statsService stats.PlatformStatsService,
) {
	// 创建处理器
	handlers := &Handlers{
		AuthHandler:       handler.NewAuthHandler(userService),
		ProfileHandler:    handler.NewProfileHandler(userService),
		PublicUserHandler: handler.NewPublicUserHandler(userService),
		AdminUserHandler:  handler.NewAdminUserHandler(userService),
		SecurityHandler:   managementApi.NewSecurityAPI(userService),
	}

	if statsService != nil {
		handlers.StatsHandler = handler.NewStatsHandler(statsService)
	}

	// 设置可选依赖
	if bookstoreService != nil {
		handlers.PublicUserHandler.SetBookstoreService(bookstoreService)
	}

	// ========================================
	// 公开路由（不需要认证）
	// ========================================
	{
		// 认证相关
		r.POST("/user-management/auth/register", handlers.AuthHandler.Register)
		r.POST("/user-management/auth/login", handlers.AuthHandler.Login)

		// 安全相关 - 邮箱验证（需要登录后发送验证码）
		// 密码重置（公开访问）
		r.POST("/user-management/password/reset-request", handlers.SecurityHandler.RequestPasswordReset)
		r.POST("/user-management/password/reset", handlers.SecurityHandler.ConfirmPasswordReset)

		// 公开用户信息
		r.GET("/user-management/users/:id", handlers.PublicUserHandler.GetUser)
		r.GET("/user-management/users/:id/profile", handlers.PublicUserHandler.GetUserProfile)
		r.GET("/user-management/users/:id/books", handlers.PublicUserHandler.GetUserBooks)
	}

	// ========================================
	// 需要认证的路由
	// ========================================
	authenticated := r.Group("/user-management")
	authenticated.Use(middleware.JWTAuth())
	{
		// 个人信息管理
		authenticated.GET("/profile", handlers.ProfileHandler.GetProfile)
		authenticated.PUT("/profile", handlers.ProfileHandler.UpdateProfile)
		authenticated.PUT("/password", handlers.ProfileHandler.ChangePassword)

		// 安全相关 - 邮箱验证（需要登录）
		authenticated.POST("/email/send-code", handlers.SecurityHandler.SendEmailVerification)
		authenticated.POST("/email/verify", handlers.SecurityHandler.VerifyEmail)

		// 统计信息（如果 StatsHandler 可用）
		if handlers.StatsHandler != nil {
			authenticated.GET("/stats/my", handlers.StatsHandler.GetMyStats)
			authenticated.GET("/stats/my/content", handlers.StatsHandler.GetMyContentStats)
			authenticated.GET("/stats/my/activity", handlers.StatsHandler.GetMyActivityStats)
			authenticated.GET("/stats/my/revenue", handlers.StatsHandler.GetMyRevenueStats)
		}
	}

	// ========================================
	// 需要管理员权限的路由
	// ========================================
	adminGroup := r.Group("/admin/users")
	adminGroup.Use(middleware.JWTAuth())
	adminGroup.Use(middleware.RequireRole("admin"))
	{
		// 用户管理
		adminGroup.GET("", handlers.AdminUserHandler.ListUsers)
		adminGroup.GET("/:id", handlers.AdminUserHandler.GetUser)
		adminGroup.PUT("/:id", handlers.AdminUserHandler.UpdateUser)
		adminGroup.DELETE("/:id", handlers.AdminUserHandler.DeleteUser)
		adminGroup.POST("/:id/ban", handlers.AdminUserHandler.BanUser)
		adminGroup.POST("/:id/unban", handlers.AdminUserHandler.UnbanUser)
	}
}

// Handlers 聚合所有处理器
type Handlers struct {
	AuthHandler       *handler.AuthHandler
	ProfileHandler    *handler.ProfileHandler
	PublicUserHandler *handler.PublicUserHandler
	AdminUserHandler  *handler.AdminUserHandler
	StatsHandler      *handler.StatsHandler
	SecurityHandler   *managementApi.SecurityAPI
}
