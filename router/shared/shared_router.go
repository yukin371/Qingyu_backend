package shared

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/middleware"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/shared/storage"
	"Qingyu_backend/service/shared/wallet"
)

// RegisterRoutes 注册共享服务路由
// 参数改为接收独立服务而不是整个容器，避免与全局 ServiceContainer 冲突
func RegisterRoutes(
	r *gin.RouterGroup,
	authService auth.AuthService,
	walletService wallet.WalletService,
	storageService *storage.StorageServiceImpl,
	multipartService *storage.MultipartUploadService,
	imageProcessor *storage.ImageProcessor,
) {
	// 应用全局中间件
	r.Use(middleware.ResponseFormatterMiddleware()) // 响应格式化（RequestID生成）
	r.Use(middleware.ResponseTimingMiddleware())    // 响应时间记录
	r.Use(middleware.CORSMiddleware())              // 跨域处理
	r.Use(middleware.Recovery())                    // Panic恢复
	r.Use(response.GzipMiddleware(5))               // Gzip压缩（压缩级别5）

	// 创建API处理器
	authAPI := shared.NewAuthAPI(authService)
	walletAPI := shared.NewWalletAPI(walletService)
	storageAPI := shared.NewStorageAPI(storageService, multipartService, imageProcessor)
	// 注意：AdminAPI已迁移到 admin 模块

	// ============ 认证服务路由 ============
	authGroup := r.Group("/auth")
	{
		// 公开路由（添加速率限制）
		publicAuth := authGroup.Group("")
		publicAuth.Use(middleware.RateLimitMiddleware(10, 60)) // 10次/分钟
		{
			publicAuth.POST("/register", authAPI.Register)
			publicAuth.POST("/login", authAPI.Login)
		}

		// 需要认证的路由
		authProtected := authGroup.Group("")
		authProtected.Use(middleware.JWTAuth())
		authProtected.Use(middleware.RateLimitMiddleware(30, 60)) // 30次/分钟
		{
			authProtected.POST("/logout", authAPI.Logout)
			authProtected.POST("/refresh", authAPI.RefreshToken)
			authProtected.GET("/permissions", authAPI.GetUserPermissions)
			authProtected.GET("/roles", authAPI.GetUserRoles)
		}
	}

	// ============ 钱包服务路由 ============
	walletGroup := r.Group("/wallet")
	walletGroup.Use(middleware.JWTAuth())                   // 所有钱包接口都需要认证
	walletGroup.Use(middleware.RateLimitMiddleware(50, 60)) // 50次/分钟
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

	// ============ 存储服务路由 ============
	storageGroup := r.Group("/storage")
	storageGroup.Use(middleware.JWTAuth())                   // 所有存储接口都需要认证
	storageGroup.Use(middleware.RateLimitMiddleware(20, 60)) // 20次/分钟（文件操作限制更严格）
	{
		// 文件操作
		storageGroup.POST("/upload", storageAPI.UploadFile)
		storageGroup.GET("/download/:file_id", storageAPI.DownloadFile)
		storageGroup.DELETE("/files/:file_id", storageAPI.DeleteFile)
		storageGroup.GET("/files/:file_id", storageAPI.GetFileInfo)
		storageGroup.GET("/files", storageAPI.ListFiles)
		storageGroup.GET("/files/:file_id/url", storageAPI.GetDownloadURL)
	}

	// 注意：管理员路由已迁移到 /api/v1/admin
	// 参见: router/admin/admin_router.go
}
