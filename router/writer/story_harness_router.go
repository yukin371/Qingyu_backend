package writer

import (
	"github.com/gin-gonic/gin"

	apiWriter "Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/writer/storyharness"
)

// InitStoryHarnessRoutes 注册 Story Harness 相关路由
func InitStoryHarnessRoutes(
	router *gin.RouterGroup,
	contextSvc *storyharness.ContextService,
	indexerSvc *storyharness.IndexerService,
	crSvc *storyharness.ChangeRequestService,
) {
	// 创建 API 处理器
	storyHarnessAPI := apiWriter.NewStoryHarnessApi(contextSvc, indexerSvc, crSvc)
	crAPI := apiWriter.NewChangeRequestApi(crSvc)

	// 项目级路由组
	projectGroup := router.Group("/projects/:id")
	{
		// 章节上下文（Context Lens）
		projectGroup.GET("/chapters/:chapterId/context", storyHarnessAPI.GetChapterContext)

		// 手动触发章节索引
		projectGroup.POST("/chapters/:chapterId/trigger-index", storyHarnessAPI.TriggerChapterIndex)

		// 手动重建章节投影
		projectGroup.POST("/chapters/:chapterId/rebuild-projection", storyHarnessAPI.RebuildChapterProjection)

		// 章节建议列表
		projectGroup.GET("/chapters/:chapterId/change-requests", crAPI.ListChangeRequests)
	}

	// 资源级路由组
	resourceGroup := router.Group("/change-requests")
	{
		// 建议详情
		resourceGroup.GET("/:requestId", crAPI.GetChangeRequest)

		// 处理建议（接受/忽略/延后）
		resourceGroup.PUT("/:requestId/status", crAPI.ProcessChangeRequest)
	}
}
