package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler 处理 /ping 请求
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
