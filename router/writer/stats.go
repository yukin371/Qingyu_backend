package writer

import (
	"Qingyu_backend/api/v1/writer"
	readingStats "Qingyu_backend/service/reading/stats"

	"github.com/gin-gonic/gin"
)

// InitStatsRouter 初始化统计路由（阅读/书店统计）
func InitStatsRouter(r *gin.RouterGroup, service *readingStats.ReadingStatsService) {
	// 创建API实例
	statsApi := writer.NewStatsApi(service)

	// 作品统计路由组
	bookStats := r.Group("/books/:book_id")
	{
		bookStats.GET("/stats", statsApi.GetBookStats)               // 获取作品统计
		bookStats.GET("/heatmap", statsApi.GetBookHeatmap)           // 获取阅读热力图
		bookStats.GET("/revenue", statsApi.GetBookRevenue)           // 获取收入统计
		bookStats.GET("/top-chapters", statsApi.GetTopChapters)      // 获取热门章节
		bookStats.GET("/daily-stats", statsApi.GetDailyStats)        // 获取每日统计
		bookStats.GET("/drop-off-points", statsApi.GetDropOffPoints) // 获取跳出点
		bookStats.GET("/retention", statsApi.GetRetentionRate)       // 获取留存率
	}

	// 章节统计路由组
	chapterStats := r.Group("/chapters/:chapter_id")
	{
		chapterStats.GET("/stats", statsApi.GetChapterStats) // 获取章节统计
	}

	// 读者行为记录（通常在读者端调用，这里也提供）
	r.POST("/reader/behavior", statsApi.RecordBehavior)
}
