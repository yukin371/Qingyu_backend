package shared

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/shared/container"
)

// RegisterRoutes 注册共享服务路由
func RegisterRoutes(r *gin.RouterGroup, serviceContainer *container.SharedServiceContainer) {
	// 创建API处理器
	authAPI := shared.NewAuthAPI(serviceContainer.AuthService())
	walletAPI := shared.NewWalletAPI(serviceContainer.WalletService())
	storageAPI := shared.NewStorageAPI(serviceContainer.StorageService())
	adminAPI := shared.NewAdminAPI(serviceContainer.AdminService())

	// ============ 认证服务路由 ============
	authGroup := r.Group("/auth")
	{
		// 公开路由
		authGroup.POST("/register", authAPI.Register)
		authGroup.POST("/login", authAPI.Login)

		// 需要认证的路由
		authProtected := authGroup.Group("")
		authProtected.Use(middleware.JWTAuth())
		{
			authProtected.POST("/logout", authAPI.Logout)
			authProtected.POST("/refresh", authAPI.RefreshToken)
			authProtected.GET("/permissions", authAPI.GetUserPermissions)
			authProtected.GET("/roles", authAPI.GetUserRoles)
		}
	}

	// ============ 钱包服务路由 ============
	walletGroup := r.Group("/wallet")
	walletGroup.Use(middleware.JWTAuth()) // 所有钱包接口都需要认证
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
	storageGroup.Use(middleware.JWTAuth()) // 所有存储接口都需要认证
	{
		// 文件操作
		storageGroup.POST("/upload", storageAPI.UploadFile)
		storageGroup.GET("/download/:file_id", storageAPI.DownloadFile)
		storageGroup.DELETE("/files/:file_id", storageAPI.DeleteFile)
		storageGroup.GET("/files/:file_id", storageAPI.GetFileInfo)
		storageGroup.GET("/files", storageAPI.ListFiles)
		storageGroup.GET("/files/:file_id/url", storageAPI.GetFileURL)
	}

	// ============ 管理服务路由 ============
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JWTAuth()) // 所有管理接口都需要认证
	// TODO: 添加管理员权限验证中间件
	{
		// 内容审核
		adminGroup.GET("/reviews/pending", adminAPI.GetPendingReviews)
		adminGroup.POST("/reviews", adminAPI.ReviewContent)

		// 提现审核
		adminGroup.POST("/withdraw/review", adminAPI.ReviewWithdraw)

		// 用户管理
		adminGroup.GET("/users/:user_id/statistics", adminAPI.GetUserStatistics)

		// 操作日志
		adminGroup.GET("/operation-logs", adminAPI.GetOperationLogs)
	}
}
