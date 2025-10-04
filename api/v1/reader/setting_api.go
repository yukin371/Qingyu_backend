package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/reader/settings 获取阅读设置
func GetReadingSettings(c *gin.Context) {
	// TODO: 获取阅读设置
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

// PUT /api/v1/reader/settings 更新阅读设置
func UpdateReadingSettings(c *gin.Context) {
	// TODO: 更新阅读设置
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}
