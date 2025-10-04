package recommendation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/recommendation/similar/:booksID 相似推荐
func GetSimilarRecommendations(c *gin.Context) {
	// TODO: 相似推荐
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}
