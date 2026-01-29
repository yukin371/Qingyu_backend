package recommendation

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
)

// GET /api/v1/recommendation/personal 个性化推荐
func GetPersonalRecommendations(c *gin.Context) {
	// TODO: 获取个人推荐
	response.Success(c, gin.H{})
}
