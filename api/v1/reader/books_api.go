package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/reader/books/:id 获取章节列表
func GetBookChapters(c *gin.Context) {
	bookId := c.Param("id")
	// TODO: 根据bookId获取章节列表
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{"bookId": bookId},
	})
}
