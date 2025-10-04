package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/reader/chapters/:id 获取章节内容
func GetChapterContent(c *gin.Context) {
	chapterId := c.Param("id")
	// TODO: 根据chapterId获取章节内容
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{"chapterId": chapterId},
	})
}
