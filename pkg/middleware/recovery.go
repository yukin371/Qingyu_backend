package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery 默认恢复中间件
// 这是一个便捷函数，用于从panic中恢复
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    50000,
					"message": "服务器内部错误",
					"data":    nil,
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
