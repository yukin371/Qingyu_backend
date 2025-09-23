package ai

import (
	aiApi "Qingyu_backend/api/v1/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// AIRouter AI路由组
type AIRouter struct {
	api             *aiApi.AIApi
	chatApi         *aiApi.ChatApi
	novelContextApi *aiApi.NovelContextApi
}

// NewAIRouter 创建AI路由
func NewAIRouter(service *aiService.Service) *AIRouter {
	// 创建聊天仓库
	chatRepository := aiService.NewChatRepository()
	
	// 创建聊天服务
	chatService := aiService.NewChatService(service, chatRepository)

	// 创建小说上下文服务 - 需要传入必要的依赖
	// 这里暂时传入nil，实际使用时需要注入正确的依赖
	novelContextService := aiService.NewNovelContextService(nil, nil, nil, nil, nil, nil, nil, nil)

	return &AIRouter{
		api:             aiApi.NewAIApi(service),
		chatApi:         aiApi.NewChatApi(chatService),
		novelContextApi: aiApi.NewNovelContextApi(novelContextService),
	}
}

// InitAIRouter 初始化AI路由
func (r *AIRouter) InitAIRouter(Router *gin.RouterGroup) {
	aiRouter := Router.Group("ai")
	{
		// 内容生成相关
		aiRouter.POST("generate", r.api.GenerateContent) // 生成内容
		aiRouter.POST("continue", r.api.ContinueWriting) // 续写内容
		aiRouter.POST("optimize", r.api.OptimizeText)    // 优化文本
		aiRouter.POST("outline", r.api.GenerateOutline)  // 生成大纲

		// 内容分析相关
		aiRouter.POST("analyze", r.api.AnalyzeContent) // 分析内容

		// 上下文管理相关
		aiRouter.GET("context/:projectId", r.api.GetContextInfo)            // 获取项目上下文
		aiRouter.GET("context/:projectId/:chapterId", r.api.GetContextInfo) // 获取章节上下文
		aiRouter.POST("context/feedback", r.api.UpdateContextFeedback)      // 更新上下文反馈

		// AI聊天相关
		chatRouter := aiRouter.Group("chat")
		{
			chatRouter.POST("", r.chatApi.StartChat)                         // 开始聊天
			chatRouter.POST("stream", r.chatApi.ContinueChat)                // 流式聊天
			chatRouter.GET("sessions", r.chatApi.ListChatSessions)           // 获取会话列表
			chatRouter.GET(":sessionId/history", r.chatApi.GetChatHistory)   // 获取聊天历史
			chatRouter.PUT(":sessionId", r.chatApi.UpdateChatSession)        // 更新会话信息
			chatRouter.DELETE(":sessionId", r.chatApi.DeleteChatSession)     // 删除会话
			chatRouter.GET("statistics", r.chatApi.GetChatStatistics)        // 获取统计信息
			chatRouter.GET(":sessionId/export", r.chatApi.ExportChatHistory) // 导出聊天历史
		}

		// 小说上下文相关
		novelRouter := aiRouter.Group("novel")
		{
			novelRouter.GET("context/:projectId", r.novelContextApi.GetNovelContext)    // 获取小说上下文
			novelRouter.POST("memory", r.novelContextApi.CreateNovelMemory)             // 创建记忆
			novelRouter.PUT("memory/:memoryId", r.novelContextApi.UpdateNovelMemory)    // 更新记忆
			novelRouter.DELETE("memory/:memoryId", r.novelContextApi.DeleteNovelMemory) // 删除记忆
			novelRouter.POST("search", r.novelContextApi.SearchNovelContext)            // 搜索上下文
		}
	}
}