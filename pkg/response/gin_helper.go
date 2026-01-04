package response

import (
	"github.com/gin-gonic/gin"
)

// JSON 是gin.Context.JSON的替代函数，不转义非ASCII字符
// 使用示例: response.JSON(c, 200, data)
func JSON(c *gin.Context, code int, obj any) {
	JsonWithNoEscape(c, code, obj)
}

// SuccessJSON 返回成功响应（不转义中文）
func SuccessJSON(c *gin.Context, message string, data any) {
	JSON(c, 200, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}

// PaginatedJSON 返回分页响应（不转义中文）
func PaginatedJSON(c *gin.Context, message string, data any, total int64, page, pageSize int) {
	JSON(c, 200, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
		"total":   total,
		"page":    page,
		"size":    pageSize,
	})
}

// ErrorJSON 返回错误响应（不转义中文）
func ErrorJSON(c *gin.Context, code int, message string) {
	JSON(c, code, gin.H{
		"code":    code,
		"message": message,
	})
}
