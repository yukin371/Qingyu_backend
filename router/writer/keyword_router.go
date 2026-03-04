package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitKeywordRouter 初始化关键词检索路由
func InitKeywordRouter(r *gin.RouterGroup, characterService interfaces.CharacterService, locationService interfaces.LocationService) {
	keywordAPI := writer.NewKeywordApi(characterService, locationService)

	projectGroup := r.Group("/projects/:id")
	{
		projectGroup.GET("/keywords/search", keywordAPI.SearchKeywords)
	}
}
