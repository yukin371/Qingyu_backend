package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitOutlineRoutes 初始化大纲路由
func InitOutlineRoutes(router *gin.RouterGroup, outlineService interfaces.OutlineService) {
	api := writer.NewOutlineApi(outlineService)

	// 项目级别的大纲路由（需要认证）
	projectGroup := router.Group("/projects/:projectId")
	{
		// 创建和列表查询
		projectGroup.POST("/outlines", api.CreateOutline)
		projectGroup.GET("/outlines", api.ListOutlines)

		// 树形结构
		projectGroup.GET("/outlines/tree", api.GetOutlineTree)

		// 子节点查询
		projectGroup.GET("/outlines/children", api.GetOutlineChildren)
	}

	// 大纲级别的路由（需要认证）
	outlineGroup := router.Group("/outlines")
	{
		// 单个大纲操作（需要传递projectId作为查询参数进行权限验证）
		outlineGroup.GET("/:outlineId", api.GetOutline)
		outlineGroup.PUT("/:outlineId", api.UpdateOutline)
		outlineGroup.DELETE("/:outlineId", api.DeleteOutline)
	}
}
