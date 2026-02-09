package shared

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/internal/middleware/auth"
	"Qingyu_backend/internal/middleware/ratelimit"
	mw "Qingyu_backend/pkg/middleware"
	"Qingyu_backend/pkg/response"
	sharedAuth "Qingyu_backend/service/auth"
	"Qingyu_backend/service/shared/storage"
	walletService "Qingyu_backend/service/finance/wallet"
)

// RegisterRoutes 注册共享服务路由（向后兼容）
// 参数改为接收独立服务而不是整个容器，避免与全局 ServiceContainer 冲突
func RegisterRoutes(
	r *gin.RouterGroup,
	authService sharedAuth.AuthService,
	oauthService *sharedAuth.OAuthService,
	logger *zap.Logger,
	walletSvc walletService.WalletService,
	storageService *storage.StorageServiceImpl,
	multipartService *storage.MultipartUploadService,
	imageProcessor *storage.ImageProcessor,
) {
	// 应用全局中间件
	r.Use(mw.ResponseFormatterMiddleware()) // 响应格式化（RequestID生成）
	r.Use(mw.ResponseTimingMiddleware())    // 响应时间记录
	r.Use(mw.CORSMiddleware())              // 跨域处理
	r.Use(mw.Recovery())                    // Panic恢复
	r.Use(response.GzipMiddleware(5))       // Gzip压缩（压缩级别5）

	// 调用各个独立的注册函数
	if authService != nil && oauthService != nil {
		RegisterAuthRoutes(r, authService, oauthService, logger)
	}
	if walletSvc != nil {
		RegisterWalletRoutes(r, walletSvc)
	}
	if storageService != nil {
		RegisterStorageRoutes(r, storageService, multipartService, imageProcessor)
	}
	// TODO: 注册搜索建议路由（需要SearchService，将在writer路由初始化后创建）
}

// RegisterAuthRoutes 注册认证服务路由
func RegisterAuthRoutes(r *gin.RouterGroup, authService sharedAuth.AuthService, oauthService *sharedAuth.OAuthService, logger *zap.Logger) {
	// 创建API处理器
	authAPI := shared.NewAuthAPI(authService)
	oauthAPI := shared.NewOAuthAPI(oauthService, authService, logger)

	// ============ 认证服务路由 ============
	authGroup := r.Group("/auth")
	{
		// 公开路由（添加速率限制）
		publicAuth := authGroup.Group("")
		publicAuth.Use(ratelimit.RateLimitMiddlewareSimple(10, 60)) // 10次/分钟
		{
			publicAuth.POST("/register", authAPI.Register)
			publicAuth.POST("/login", authAPI.Login)
		}

		// 需要认证的路由
		authProtected := authGroup.Group("")
		authProtected.Use(auth.JWTAuth())
		authProtected.Use(ratelimit.RateLimitMiddlewareSimple(30, 60)) // 30次/分钟
		{
			authProtected.POST("/logout", authAPI.Logout)
			authProtected.POST("/refresh", authAPI.RefreshToken)
			authProtected.GET("/permissions", authAPI.GetUserPermissions)
			authProtected.GET("/roles", authAPI.GetUserRoles)
		}
	}

	// ============ OAuth认证路由 ============
	oauthGroup := r.Group("/oauth")
	{
		// 公开路由（添加速率限制）
		publicOAuth := oauthGroup.Group("")
		publicOAuth.Use(ratelimit.RateLimitMiddlewareSimple(10, 60)) // 10次/分钟
		{
			// 获取授权URL
			publicOAuth.POST("/:provider/authorize", oauthAPI.GetAuthorizeURL)
			// OAuth回调
			publicOAuth.POST("/:provider/callback", oauthAPI.HandleCallback)
		}

		// 需要认证的路由
		oauthProtected := oauthGroup.Group("")
		oauthProtected.Use(auth.JWTAuth())
		oauthProtected.Use(ratelimit.RateLimitMiddlewareSimple(20, 60)) // 20次/分钟
		{
			// 获取绑定的账号列表
			oauthProtected.GET("/accounts", oauthAPI.GetLinkedAccounts)
			// 解绑账号
			oauthProtected.DELETE("/accounts/:accountID", oauthAPI.UnlinkAccount)
			// 设置主账号
			oauthProtected.PUT("/accounts/:accountID/primary", oauthAPI.SetPrimaryAccount)
		}
	}
}

// RegisterWalletRoutes 注册钱包服务路由
func RegisterWalletRoutes(r *gin.RouterGroup, walletService walletService.WalletService) {
	// 创建API处理器
	walletAPI := shared.NewWalletAPI(walletService)

	// ============ 钱包服务路由 ============
	walletGroup := r.Group("/wallet")
	walletGroup.Use(auth.JWTAuth())                   // 所有钱包接口都需要认证
	walletGroup.Use(ratelimit.RateLimitMiddlewareSimple(50, 60)) // 50次/分钟
	{
		// 查询接口
		walletGroup.GET("/balance", walletAPI.GetBalance)
		walletGroup.GET("", walletAPI.GetWallet)
		walletGroup.GET("/transactions", walletAPI.GetTransactions)
		walletGroup.GET("/withdrawals", walletAPI.GetWithdrawRequests)

		// 操作接口
		walletGroup.POST("/recharge", walletAPI.Recharge)
		walletGroup.POST("/consume", walletAPI.Consume)
		walletGroup.POST("/transfer", walletAPI.Transfer)
		walletGroup.POST("/withdraw", walletAPI.RequestWithdraw)
	}
}

// RegisterStorageRoutes 注册存储服务路由
func RegisterStorageRoutes(
	r *gin.RouterGroup,
	storageService *storage.StorageServiceImpl,
	multipartService *storage.MultipartUploadService,
	imageProcessor *storage.ImageProcessor,
) {
	// 创建API处理器
	storageAPI := shared.NewStorageAPI(storageService, multipartService, imageProcessor)

	// ============ 存储服务路由 ============
	storageGroup := r.Group("/storage")
	storageGroup.Use(auth.JWTAuth())                   // 所有存储接口都需要认证
	storageGroup.Use(ratelimit.RateLimitMiddlewareSimple(20, 60)) // 20次/分钟（文件操作限制更严格）
	{
		// 文件操作
		storageGroup.POST("/upload", storageAPI.UploadFile)
		storageGroup.GET("/download/:file_id", storageAPI.DownloadFile)
		storageGroup.DELETE("/files/:file_id", storageAPI.DeleteFile)
		storageGroup.GET("/files/:file_id", storageAPI.GetFileInfo)
		storageGroup.GET("/files", storageAPI.ListFiles)
		storageGroup.GET("/files/:file_id/url", storageAPI.GetDownloadURL)
	}
}

// 注意：管理员路由已迁移到 /api/v1/admin
// 参见: router/admin/admin_router.go
