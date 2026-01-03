package ai

import (
	aiApi "Qingyu_backend/api/v1/ai"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// InitAIRouter 初始化AI路由
func InitAIRouter(r *gin.RouterGroup, aiService *ai.Service, chatService *ai.ChatService, quotaService *ai.QuotaService, phase3Client *ai.Phase3Client) {
	// 创建API实例
	writingApiHandler := aiApi.NewWritingApi(aiService, quotaService)
	chatApiHandler := aiApi.NewChatApi(chatService, quotaService)
	systemApiHandler := aiApi.NewSystemApi(aiService)
	quotaApiHandler := aiApi.NewQuotaApi(quotaService)

	// Phase3创作API（如果客户端可用）
	var creativeApiHandler *aiApi.CreativeAPI
	if phase3Client != nil {
		creativeApiHandler = aiApi.NewCreativeAPI(phase3Client)
	}

	// 创建写作辅助服务实例
	summarizeService := ai.NewSummarizeService(aiService.GetAdapterManager())
	proofreadService := ai.NewProofreadService(aiService.GetAdapterManager())
	sensitiveWordsService := ai.NewSensitiveWordsService(aiService.GetAdapterManager())
	writingAssistantApiHandler := aiApi.NewWritingAssistantApi(summarizeService, proofreadService, sensitiveWordsService)

	// AI主路由组
	aiGroup := r.Group("/ai")
	aiGroup.Use(middleware.JWTAuth()) // 需要认证
	{
		// 系统功能（不需要配额）
		aiGroup.GET("/health", systemApiHandler.HealthCheck)
		aiGroup.GET("/providers", systemApiHandler.GetProviders)
		aiGroup.GET("/models", systemApiHandler.GetModels)

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
			writingGroup.POST("/continue", writingApiHandler.ContinueWriting)
			writingGroup.POST("/continue/stream", writingApiHandler.ContinueWritingStream)

			// 内容改写
			writingGroup.POST("/rewrite", writingApiHandler.RewriteText)
			writingGroup.POST("/rewrite/stream", writingApiHandler.RewriteTextStream)

			// ============ 新增：写作辅助功能 ============
			// 内容总结
			writingGroup.POST("/summarize", writingAssistantApiHandler.SummarizeContent)
			writingGroup.POST("/summarize-chapter", writingAssistantApiHandler.SummarizeChapter)

			// 文本校对
			writingGroup.POST("/proofread", writingAssistantApiHandler.ProofreadContent)
			writingGroup.GET("/suggestions/:id", writingAssistantApiHandler.GetProofreadSuggestion)
		}

		// AI内容审核功能
		auditGroup := aiGroup.Group("/audit")
		auditGroup.Use(middleware.QuotaCheckMiddleware(quotaService))
		{
			// 敏感词检测
			auditGroup.POST("/sensitive-words", writingAssistantApiHandler.CheckSensitiveWords)
			auditGroup.GET("/sensitive-words/:id", writingAssistantApiHandler.GetSensitiveWordsDetail)
		}

		// AI聊天功能（轻量级配额检查）
		chatGroup := aiGroup.Group("/chat")
		chatGroup.Use(middleware.LightQuotaCheckMiddleware(quotaService))
		{
			chatGroup.POST("", chatApiHandler.Chat)
			chatGroup.POST("/stream", chatApiHandler.ChatStream)
			chatGroup.GET("/sessions", chatApiHandler.GetChatSessions)
			chatGroup.GET("/sessions/:sessionId", chatApiHandler.GetChatHistory)
			chatGroup.DELETE("/sessions/:sessionId", chatApiHandler.DeleteChatSession)
		}

		// ============ 兼容性路由（支持旧的API路径） ============
		// 这些路由将重定向到新的API路径，保持向后兼容
		compatGroup := aiGroup.Group("")
		compatGroup.Use(middleware.QuotaCheckMiddleware(quotaService))
		{
			// /api/v1/ai/generate -> /api/v1/ai/writing/continue
			compatGroup.POST("/generate", writingApiHandler.ContinueWriting)

			// /api/v1/ai/rewrite, /api/v1/ai/expand, /api/v1/ai/polish -> /api/v1/ai/writing/rewrite
			compatGroup.POST("/rewrite", writingApiHandler.RewriteText)
			compatGroup.POST("/expand", writingApiHandler.RewriteText)
			compatGroup.POST("/polish", writingApiHandler.RewriteText)
		}
	}

	// ============ Phase3 创作路由 ============
	if creativeApiHandler != nil {
		InitCreativeRoutes(r, creativeApiHandler)
	}

	// 注意：管理员配额管理路由已迁移到 /api/v1/admin/quota
	// 参见: router/admin/admin_router.go
}
