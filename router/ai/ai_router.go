package ai

import (
	aiApi "Qingyu_backend/api/v1/ai"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// InitAIRouter 初始化AI路由
func InitAIRouter(r *gin.RouterGroup, aiService *ai.Service, chatService *ai.ChatService, quotaService *ai.QuotaService) {
	// 创建API实例
	aiApiHandler := aiApi.NewAIApi(aiService, chatService, quotaService)
	quotaApiHandler := aiApi.NewQuotaApi(quotaService)

	// AI主路由组
	aiGroup := r.Group("/ai")
	aiGroup.Use(middleware.JWTAuth()) // 需要认证
	{
		// 健康检查（不需要配额）
		aiGroup.GET("/health", aiApiHandler.HealthCheck)

		// 配额管理（不需要配额检查）
		quotaGroup := aiGroup.Group("/quota")
		{
			quotaGroup.GET("", quotaApiHandler.GetQuotaInfo)
			quotaGroup.GET("/all", quotaApiHandler.GetAllQuotas)
			quotaGroup.GET("/statistics", quotaApiHandler.GetQuotaStatistics)
			quotaGroup.GET("/transactions", quotaApiHandler.GetTransactionHistory)
		}

		// AI写作功能（需要配额检查）
		writingGroup := aiGroup.Group("/writing")
		writingGroup.Use(middleware.QuotaCheckMiddleware(quotaService))
		{
			// 智能续写
			writingGroup.POST("/continue", aiApiHandler.ContinueWriting)
			writingGroup.POST("/continue/stream", aiApiHandler.ContinueWritingStream)

			// 内容改写
			writingGroup.POST("/rewrite", aiApiHandler.RewriteText)
			writingGroup.POST("/rewrite/stream", aiApiHandler.RewriteTextStream)
		}

		// AI聊天功能（轻量级配额检查）
		chatGroup := aiGroup.Group("/chat")
		chatGroup.Use(middleware.LightQuotaCheckMiddleware(quotaService))
		{
			chatGroup.POST("", aiApiHandler.Chat)
			chatGroup.POST("/stream", aiApiHandler.ChatStream)
			chatGroup.GET("/sessions", aiApiHandler.GetChatSessions)
			chatGroup.GET("/sessions/:sessionId", aiApiHandler.GetChatHistory)
			chatGroup.DELETE("/sessions/:sessionId", aiApiHandler.DeleteChatSession)
		}
	}

	// 管理员路由（配额管理）
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JWTAuth())
	adminGroup.Use(middleware.AdminPermissionMiddleware()) // 使用管理员权限检查
	{
		adminQuotaGroup := adminGroup.Group("/quota")
		{
			adminQuotaGroup.PUT("/:userId", quotaApiHandler.UpdateUserQuota)
			adminQuotaGroup.POST("/:userId/suspend", quotaApiHandler.SuspendUserQuota)
			adminQuotaGroup.POST("/:userId/activate", quotaApiHandler.ActivateUserQuota)
		}
	}
}
