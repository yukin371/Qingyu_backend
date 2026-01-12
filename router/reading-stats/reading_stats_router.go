package readingstats

import (
	"github.com/gin-gonic/gin"

	statsAPI "Qingyu_backend/api/v1/stats"
	"Qingyu_backend/middleware"
	readingStatsService "Qingyu_backend/service/reader/stats"
)

// RegisterReadingStatsRoutes 注册阅读统计路由
func RegisterReadingStatsRoutes(
	r *gin.RouterGroup,
	statsService *readingStatsService.ReadingStatsService,
) {
	if statsService == nil {
		return
	}

	// 创建API处理器
	readingStatsAPI := statsAPI.NewReadingStatsAPI(statsService)

	// 应用中间件
	statsGroup := r.Group("/stats")
	statsGroup.Use(middleware.ResponseFormatterMiddleware())
	statsGroup.Use(middleware.ResponseTimingMiddleware())
	statsGroup.Use(middleware.CORSMiddleware())
	statsGroup.Use(middleware.Recovery())

	// ============ 我的统计路由（需要认证） ============
	myStats := statsGroup.Group("/my")
	myStats.Use(middleware.JWTAuth())
	{
		// 获取我的统计
		myStats.GET("/stats", readingStatsAPI.GetMyStats)

		// 获取每日统计
		myStats.GET("/daily", readingStatsAPI.GetMyDailyStats)

		// 获取排名
		myStats.GET("/ranking", readingStatsAPI.GetMyRanking)

		// 获取阅读时长
		myStats.GET("/reading-time", readingStatsAPI.GetMyReadingTime)

		// 获取阅读历史
		myStats.GET("/history", readingStatsAPI.GetMyHistory)
	}

	// ============ 推荐路由（需要认证） ============
	recommendations := statsGroup.Group("/recommendations")
	recommendations.Use(middleware.JWTAuth())
	{
		recommendations.GET("", readingStatsAPI.GetRecommendations)
	}
}
