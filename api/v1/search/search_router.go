package search

import (
	"github.com/gin-gonic/gin"
	searchService "Qingyu_backend/service/search"
)

// RegisterSearchRoutes 注册搜索路由
func RegisterSearchRoutes(router *gin.RouterGroup, searchSvc *searchService.SearchService) {
	searchAPI := NewSearchAPI(searchSvc)

	// 统一搜索路由组
	searchGroup := router.Group("/search")
	{
		// 统一搜索入口
		searchGroup.POST("/search", searchAPI.Search)

		// 批量搜索
		searchGroup.POST("/batch", searchAPI.SearchBatch)

		// 健康检查
		searchGroup.GET("/health", searchAPI.Health)
	}
}
