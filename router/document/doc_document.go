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
		d := api.NewDocumentApi()
		v := api.NewVersionApi()
		// 创建文档
		docRouter.POST("/", d.Create)

		// 获取文档列表
		docRouter.GET("/", d.List)

		// 获取单个文档
		docRouter.GET("/:id", d.Get)

		// 更新文档
		docRouter.PUT("/:id", d.Update)

		// 删除文档
		docRouter.DELETE("/:id", d.Delete)

		// 版本相关
		// 创建新版本（提交内容）
		docRouter.POST(":projectId/:nodeId/version", v.CreateVersion)
		docRouter.POST(":projectId/:nodeId/rollback", v.Rollback)
		docRouter.POST(":projectId/:nodeId/patch", v.CreatePatch)
		docRouter.POST(":projectId/:nodeId/patch/:patchId/apply", v.ApplyPatch)
		docRouter.GET(":projectId/:nodeId/versions", v.ListVersions)
	}
}
