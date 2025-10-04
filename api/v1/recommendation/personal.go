package recommendation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/recommendation/personal 个性化推荐
func GetPersonalRecommendations(c *gin.Context) {
	// TODO: 获取个人推荐
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}
