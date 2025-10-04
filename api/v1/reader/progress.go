package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/reader/progress/:bookId 获取阅读进度
func GetReadingProgress(c *gin.Context) {
	// TODO: 获取阅读进度
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

// PUT /api/v1/reader/progress 更新阅读进度
func UpdateReadingProgress(c *gin.Context) {
	// TODO: 更新阅读进度
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}
