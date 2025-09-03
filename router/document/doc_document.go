package document

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文档相关路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 文档路由组
	docRouter := r.Group("/document")
	{
		// 创建文档
		docRouter.POST("/", nil) // TODO: 添加创建文档处理函数

		// 获取文档列表
		docRouter.GET("/", nil) // TODO: 添加获取文档列表处理函数

		// 获取单个文档
		docRouter.GET("/:id", nil) // TODO: 添加获取单个文档处理函数

		// 更新文档
		docRouter.PUT("/:id", nil) // TODO: 添加更新文档处理函数

		// 删除文档
		docRouter.DELETE("/:id", nil) // TODO: 添加删除文档处理函数
	}
}