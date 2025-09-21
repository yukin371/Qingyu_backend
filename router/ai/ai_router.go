package ai

import (
	aiApi "Qingyu_backend/api/v1/ai"
	aiService "Qingyu_backend/service/ai"
	"github.com/gin-gonic/gin"
)

// AIRouter AI路由组
type AIRouter struct {
	api *aiApi.AIApi
}

// NewAIRouter 创建AI路由
func NewAIRouter(service *aiService.Service) *AIRouter {
	return &AIRouter{
		api: aiApi.NewAIApi(service),
	}
}

// InitAIRouter 初始化AI路由
func (r *AIRouter) InitAIRouter(Router *gin.RouterGroup) {
	aiRouter := Router.Group("ai")
	{
		// 内容生成相关
		aiRouter.POST("generate", r.api.GenerateContent)      // 生成内容
		aiRouter.POST("continue", r.api.ContinueWriting)      // 续写内容
		aiRouter.POST("optimize", r.api.OptimizeText)         // 优化文本
		aiRouter.POST("outline", r.api.GenerateOutline)       // 生成大纲
		
		// 内容分析相关
		aiRouter.POST("analyze", r.api.AnalyzeContent)        // 分析内容
		
		// 上下文管理相关
		aiRouter.GET("context/:projectId", r.api.GetContextInfo)                    // 获取项目上下文
		aiRouter.GET("context/:projectId/:chapterId", r.api.GetContextInfo)         // 获取章节上下文
		aiRouter.POST("context/feedback", r.api.UpdateContextFeedback)              // 更新上下文反馈
	}
}