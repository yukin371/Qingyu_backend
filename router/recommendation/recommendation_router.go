package recommendation

import (
	"github.com/gin-gonic/gin"

	recoAPI "Qingyu_backend/api/v1/recommendation"
	"Qingyu_backend/internal/middleware/auth"
)

// RegisterRecommendationRoutes 注册推荐系统路由
func RegisterRecommendationRoutes(router *gin.RouterGroup, api *recoAPI.RecommendationAPI) {
	reco := router.Group("/recommendation")
	{
		// 需要认证的路由
		authenticated := reco.Group("")
		authenticated.Use(auth.JWTAuth())
		{
			// 个性化推荐
			authenticated.GET("/personalized", api.GetPersonalizedRecommendations)

			// 记录行为
			authenticated.POST("/behavior", api.RecordBehavior)
		}

		// 不需要认证的路由（公开推荐）
		// 首页推荐（混合推荐策略）
		reco.GET("/homepage", api.GetHomepageRecommendations)

		// 热门推荐
		reco.GET("/hot", api.GetHotRecommendations)

		// 分类推荐
		reco.GET("/category", api.GetCategoryRecommendations)

		// 相似物品推荐
		reco.GET("/similar", api.GetSimilarItems)

		// 推荐榜查询（公开）
		reco.GET("/tables", api.ListTables)
		reco.GET("/tables/:id", api.GetTable)

		// 推荐榜管理（管理员，当前使用JWT认证）
		admin := reco.Group("/admin")
		admin.Use(auth.JWTAuth())
		{
			admin.PUT("/tables/auto/:tableType/:period", api.UpsertAutoTable)
			admin.POST("/tables/manual", api.CreateManualTable)
			admin.PUT("/tables/manual/:id", api.UpdateManualTable)
			admin.DELETE("/tables/manual/:id", api.DeleteTable)
		}
	}
}
