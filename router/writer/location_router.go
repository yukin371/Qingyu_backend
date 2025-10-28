package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitLocationRoutes 初始化地点路由
func InitLocationRoutes(router *gin.RouterGroup, locationService interfaces.LocationService) {
	api := writer.NewLocationApi(locationService)

	// 项目级别的地点路由
	projectGroup := router.Group("/projects/:projectId")
	{
		projectGroup.POST("/locations", api.CreateLocation)
		projectGroup.GET("/locations", api.ListLocations)
		projectGroup.GET("/locations/tree", api.GetLocationTree)
		projectGroup.GET("/locations/relations", api.ListLocationRelations)
	}

	// 地点级别的路由
	locationGroup := router.Group("/locations")
	{
		locationGroup.GET("/:locationId", api.GetLocation)
		locationGroup.PUT("/:locationId", api.UpdateLocation)
		locationGroup.DELETE("/:locationId", api.DeleteLocation)

		locationGroup.POST("/relations", api.CreateLocationRelation)
		locationGroup.DELETE("/relations/:relationId", api.DeleteLocationRelation)
	}
}
