package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
)

// InitEncyclopediaRouter 初始化设定百科路由
func InitEncyclopediaRouter(r *gin.RouterGroup) {
	encyclopediaApi := writer.NewEncyclopediaApi(nil) // conceptRepo 暂为nil，返回空列表

	projectGroup := r.Group("/projects/:id")
	{
		projectGroup.GET("/concepts", encyclopediaApi.ListConcepts)
		projectGroup.GET("/concepts/search", encyclopediaApi.SearchConcepts)
		projectGroup.GET("/concepts/:conceptId", encyclopediaApi.GetConcept)
	}
}
