package search

import (
	"github.com/gin-gonic/gin"
)

// RegisterSearchRoutes 注册搜索路由
func RegisterSearchRoutes(router *gin.RouterGroup) {
	searchAPI := NewSearchAPI()

	// 统一搜索路由组
	searchGroup := router.Group("/search")
	{
		// 无需认证的搜索
		searchGroup.GET("/books", searchAPI.SearchBooks)
		searchGroup.GET("/users", searchAPI.SearchUsers)

		// 需要认证的搜索
		// TODO: 添加认证中间件
		searchGroup.GET("/projects", searchAPI.SearchProjects)
		searchGroup.GET("/documents", searchAPI.SearchDocuments)
	}
}
