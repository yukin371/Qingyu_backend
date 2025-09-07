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
		a := api.NewDocumentApi()
		// 创建文档
		docRouter.POST("/", a.Create)

		// 获取文档列表
		docRouter.GET("/", a.List)

		// 获取单个文档
		docRouter.GET("/:id", a.Get)

		// 更新文档
		docRouter.PUT("/:id", a.Update)

		// 删除文档
		docRouter.DELETE("/:id", a.Delete)
	}
}