package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/reader/annotations/:id 获取注释列表
func GetAnnotations(c *gin.Context) {
	// TODO: 获取注释列表
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

// POST /api/v1/reader/annotations 创建注释
func CreateAnnotation(c *gin.Context) {
	// TODO: 创建注释
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

// PUT /api/v1/reader/annotations/:id 更新注释
func UpdateAnnotation(c *gin.Context) {
	// TODO: 更新注释
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}
