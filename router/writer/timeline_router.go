package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitTimelineRoutes 初始化时间线路由
func InitTimelineRoutes(router *gin.RouterGroup, timelineService interfaces.TimelineService) {
	api := writer.NewTimelineApi(timelineService)

	// 项目级别的时间线路由
	projectGroup := router.Group("/projects/:projectId")
	{
		projectGroup.POST("/timelines", api.CreateTimeline)
		projectGroup.GET("/timelines", api.ListTimelines)
	}

	// 时间线级别的路由
	timelineGroup := router.Group("/timelines")
	{
		timelineGroup.GET("/:timelineId", api.GetTimeline)
		timelineGroup.DELETE("/:timelineId", api.DeleteTimeline)

		// 事件管理
		timelineGroup.POST("/:timelineId/events", api.CreateTimelineEvent)
		timelineGroup.GET("/:timelineId/events", api.ListTimelineEvents)
		timelineGroup.GET("/:timelineId/visualization", api.GetTimelineVisualization)
	}

	// 事件级别的路由
	eventGroup := router.Group("/timeline-events")
	{
		eventGroup.GET("/:eventId", api.GetTimelineEvent)
		eventGroup.PUT("/:eventId", api.UpdateTimelineEvent)
		eventGroup.DELETE("/:eventId", api.DeleteTimelineEvent)
	}
}
