package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitExportRoutes 初始化导出路由
func InitExportRoutes(router *gin.RouterGroup, exportService interfaces.ExportService) {
	api := writer.NewExportApi(exportService)

	// 文档导出路由
	documentGroup := router.Group("/documents/:id")
	{
		documentGroup.POST("/export", api.ExportDocument)
	}

	// 项目导出路由
	projectGroup := router.Group("/projects/:projectId")
	{
		projectGroup.POST("/export", api.ExportProject)
		projectGroup.GET("/exports", api.ListExportTasks)
	}

	// 导出任务管理路由
	exportsGroup := router.Group("/exports")
	{
		exportsGroup.GET("/:id", api.GetExportTask)
		exportsGroup.GET("/:id/download", api.DownloadExportFile)
		exportsGroup.DELETE("/:id", api.DeleteExportTask)
		exportsGroup.POST("/:id/cancel", api.CancelExportTask)
	}
}
