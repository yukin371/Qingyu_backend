package recommendation

import (
	"github.com/gin-gonic/gin"

	recoAPI "Qingyu_backend/api/v1/recommendation"
	"Qingyu_backend/middleware"
)

// RegisterRecommendationRoutes 注册推荐系统路由
func RegisterRecommendationRoutes(router *gin.RouterGroup, api *recoAPI.RecommendationAPI, authMiddleware *middleware.AuthMiddleware) {
	reco := router.Group("/recommendation")
	{
		// 需要认证的路由
		authenticated := reco.Group("")
		authenticated.Use(authMiddleware.RequireAuth())
		{
			// 个性化推荐
			authenticated.GET("/personalized", api.GetPersonalizedRecommendations)

			// 记录行为
			authenticated.POST("/behavior", api.RecordBehavior)
		}

		// 不需要认证的路由（公开推荐）
		reco.GET("/similar", api.GetSimilarItems)
	}
}
