package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitPublishRoutes 初始化发布管理路由
func InitPublishRoutes(router *gin.RouterGroup, publishService interfaces.PublishService) {
	api := writer.NewPublishApi(publishService)

	// 项目发布路由（使用 :id 以避免与现有项目路由冲突）
	projectGroup := router.Group("/projects/:id")
	{
		// 项目发布
		projectGroup.POST("/publish", api.PublishProject)
		projectGroup.POST("/unpublish", api.UnpublishProject)
		projectGroup.GET("/publication-status", api.GetProjectPublicationStatus)

		// 发布记录
		projectGroup.GET("/publications", api.GetPublicationRecords)

		// 批量发布文档
		projectGroup.POST("/documents/batch-publish", api.BatchPublishDocuments)
	}

	// 文档发布路由
	documentGroup := router.Group("/documents/:id")
	{
		documentGroup.POST("/publish", api.PublishDocument)
		documentGroup.PUT("/publish-status", api.UpdateDocumentPublishStatus)
	}

	// 发布记录详情路由
	publicationsGroup := router.Group("/publications")
	{
		publicationsGroup.GET("/:id", api.GetPublicationRecord)
	}
}
