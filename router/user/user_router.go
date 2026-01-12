package user

import (
	"github.com/gin-gonic/gin"

	managementApi "Qingyu_backend/api/v1/user"
	"Qingyu_backend/api/v1/user/handler"
	"Qingyu_backend/middleware"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"Qingyu_backend/service/shared/stats"
)

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(
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
		r.POST("/user/auth/register", handlers.AuthHandler.Register)
		r.POST("/user/auth/login", handlers.AuthHandler.Login)

		// 安全相关 - 邮箱验证（需要登录后发送验证码）
		// 密码重置（公开访问）
		r.POST("/user/password/reset-request", handlers.SecurityHandler.RequestPasswordReset)
		r.POST("/user/password/reset", handlers.SecurityHandler.ConfirmPasswordReset)

		// 公开用户信息
		r.GET("/user/users/:id", handlers.PublicUserHandler.GetUser)
		r.GET("/user/users/:id/profile", handlers.PublicUserHandler.GetUserProfile)
		r.GET("/user/users/:id/books", handlers.PublicUserHandler.GetUserBooks)
	}

	// ========================================
	// 需要认证的路由
	// ========================================
	authenticated := r.Group("/user")
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
}

// Handlers 聚合所有处理器
type Handlers struct {
	AuthHandler       *handler.AuthHandler
	ProfileHandler    *handler.ProfileHandler
	PublicUserHandler *handler.PublicUserHandler
	StatsHandler      *handler.StatsHandler
	SecurityHandler   *managementApi.SecurityAPI
}
