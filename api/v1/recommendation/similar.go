package recommendation

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
)

// GET /api/v1/recommendation/similar/:booksID 相似推荐
func GetSimilarRecommendations(c *gin.Context) {
	// TODO: 相似推荐
	response.Success(c, gin.H{})
}
