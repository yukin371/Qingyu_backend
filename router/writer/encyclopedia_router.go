package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitEncyclopediaRouter 初始化设定百科路由
func InitEncyclopediaRouter(r *gin.RouterGroup, conceptRepo interfaces.ConceptCatalogPort) {
	encyclopediaApi := writer.NewEncyclopediaApi(conceptRepo)

	projectGroup := r.Group("/projects/:id")
	{
		projectGroup.GET("/concepts", encyclopediaApi.ListConcepts)
		projectGroup.GET("/concepts/search", encyclopediaApi.SearchConcepts)
		projectGroup.GET("/concepts/:conceptId", encyclopediaApi.GetConcept)
	}
}
