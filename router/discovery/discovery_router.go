package discovery

import (
	discoveryAPI "Qingyu_backend/api/v1/discovery"
	"Qingyu_backend/internal/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterDiscoveryRoutes(router *gin.RouterGroup, api *discoveryAPI.DiscoveryAPI) {
	discovery := router.Group("/discovery")
	discovery.Use(auth.OptionalJWTAuth())
	{
		discovery.GET("/recommendations", api.GetRecommendations)
		discovery.GET("/personalized", api.GetPersonalized)
		discovery.GET("/new-releases", api.GetNewReleases)
		discovery.GET("/editors-pick", api.GetEditorsPick)
		discovery.GET("/trending", api.GetTrending)
		discovery.GET("/topics", api.GetTopics)
		discovery.PUT("/preferences", api.UpdatePreferences)
		discovery.POST("/track", api.TrackAction)
	}
}
