package search

import (
	"github.com/gin-gonic/gin"
	"Qingyu_backend/middleware"
	searchService "Qingyu_backend/service/search"
)

// RegisterSearchRoutes 注册搜索路由
func RegisterSearchRoutes(router *gin.RouterGroup, searchSvc *searchService.SearchService) {
	searchAPI := NewSearchAPI(searchSvc)
	grayscaleAPI := NewGrayscaleAPI(searchSvc)

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

	// 灰度配置管理路由组（需要认证）
	grayscaleGroup := router.Group("/search/grayscale")
	grayscaleGroup.Use(middleware.JWTAuth())
	{
		// 获取灰度状态
		grayscaleGroup.GET("/status", grayscaleAPI.GetGrayscaleStatus)

		// 更新灰度配置
		grayscaleGroup.POST("/config", grayscaleAPI.UpdateGrayscaleConfig)
	}
}
