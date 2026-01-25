package user

import (
	"github.com/gin-gonic/gin"

	managementApi "Qingyu_backend/api/v1/user"
	"Qingyu_backend/api/v1/user/handler"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	userService "Qingyu_backend/service/user"
	userConstants "Qingyu_backend/service/user"
	"Qingyu_backend/middleware"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"Qingyu_backend/service/shared/stats"
)

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(
	r *gin.RouterGroup,
	userSvc userServiceInterface.UserService,
	userRepo interface{},
	bookstoreService interface{},
	statsService stats.PlatformStatsService,
) {
	// 创建验证服务
	// 注意：这里使用类型断言来获取UserRepository接口
	// TODO: 未来应该直接传入具体的Repository接口类型，而不是使用interface{}
	verificationService := userService.NewVerificationService(
		userRepo.(repoInterfaces.UserRepository),
		nil, // authRepo
		nil, // emailService (暂时为nil，使用模拟发送)
	)

	// 创建密码服务
	passwordService := userService.NewPasswordService(
		verificationService,
		userRepo.(repoInterfaces.UserRepository),
	)

	// 创建处理器
	handlers := &Handlers{
		AuthHandler:           handler.NewAuthHandler(userSvc),
		ProfileHandler:        handler.NewProfileHandler(userSvc),
		PublicUserHandler:     handler.NewPublicUserHandler(userSvc),
		SecurityHandler:       managementApi.NewSecurityAPI(userSvc),
		VerificationAPI:       managementApi.NewVerificationAPI(verificationService, userSvc),
		PasswordAPI:           managementApi.NewPasswordAPI(passwordService),
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
		// 认证相关 - 完整路径（新架构）
		r.POST("/user/auth/register", handlers.AuthHandler.Register)
		r.POST("/user/auth/login", handlers.AuthHandler.Login)

		// 认证相关 - 简化路径（向后兼容）
		r.POST("/register", handlers.AuthHandler.Register)
		r.POST("/login", handlers.AuthHandler.Login)

		// 安全相关 - 邮箱验证（需要登录后发送验证码）
		// 密码重置（公开访问）
		r.POST("/user/password/reset-request", handlers.SecurityHandler.RequestPasswordReset)
		r.POST("/user/password/reset", handlers.SecurityHandler.ConfirmPasswordReset)

		// 验证和密码相关（新API） - 公开访问
		// 应用频率限制：每分钟最多3次请求
		verifyEmailRateLimit := middleware.RateLimitMiddleware(userConstants.VerificationRateLimitCount, userConstants.VerificationRateLimitWindow)
		verifyPhoneRateLimit := middleware.RateLimitMiddleware(userConstants.VerificationRateLimitCount, userConstants.VerificationRateLimitWindow)
		passwordResetRateLimit := middleware.RateLimitMiddleware(userConstants.VerificationRateLimitCount, userConstants.VerificationRateLimitWindow)

		r.POST("/users/verify/email/send", verifyEmailRateLimit, handlers.VerificationAPI.SendEmailVerifyCode)
		r.POST("/users/verify/phone/send", verifyPhoneRateLimit, handlers.VerificationAPI.SendPhoneVerifyCode)
		r.POST("/users/email/verify", handlers.VerificationAPI.VerifyEmail)
		r.POST("/users/password/reset/send", passwordResetRateLimit, handlers.PasswordAPI.SendPasswordResetCode)
		r.POST("/users/password/reset/verify", handlers.PasswordAPI.ResetPassword)

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

		// 验证和密码相关（新API - 需要认证）
		authenticated.DELETE("/users/email/unbind", handlers.VerificationAPI.UnbindEmail)
		authenticated.DELETE("/users/phone/unbind", handlers.VerificationAPI.UnbindPhone)
		authenticated.DELETE("/users/devices/:deviceId", handlers.VerificationAPI.DeleteDevice)
		authenticated.PUT("/users/password", handlers.PasswordAPI.UpdatePassword)

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
	VerificationAPI   *managementApi.VerificationAPI
	PasswordAPI       *managementApi.PasswordAPI
}
