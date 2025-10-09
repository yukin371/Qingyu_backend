package document

import (
	api "Qingyu_backend/api/v1/document"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文档相关路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 文档路由组
	docRouter := r.Group("/document")
	{
		// 文档相关路由
		d := api.NewDocumentApi()
		v := api.NewVersionApi()

		// 基础文档操作
		docRouter.POST("/", d.Create)
		docRouter.GET("/", d.List)
		docRouter.GET("/doc/:id", d.Get)
		docRouter.PUT("/doc/:id", d.Update)
		docRouter.DELETE("/doc/:id", d.Delete)

		// 版本相关路由 - 使用不同的路径前缀避免冲突
		versionRouter := docRouter.Group("/version")
		{
			versionRouter.POST("/:projectId/:nodeId", v.CreateVersion)
			versionRouter.POST("/:projectId/:nodeId/rollback", v.Rollback)
			versionRouter.POST("/:projectId/:nodeId/patch", v.CreatePatch)
			versionRouter.POST("/:projectId/:nodeId/patch/:patchId/apply", v.ApplyPatch)
			versionRouter.GET("/:projectId/:nodeId/versions", v.ListVersions)
		}
	}
}
