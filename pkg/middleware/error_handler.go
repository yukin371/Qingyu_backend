package middleware

import (
	"Qingyu_backend/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last().Err

			// 转换为应用错误
			appErr := errors.GetAppError(err)

			// 添加请求ID
			requestID := c.GetString("requestId")
			if requestID != "" {
				appErr.WithRequestID(requestID)
			}

			// 格式化错误响应
			handler := errors.DefaultErrorHandler
			status, response := handler.Handle(appErr)

			// 记录错误日志（5xx错误）
			if status >= 500 {
				// TODO: 使用结构化日志记录
				// logger.Error("server error",
				//     zap.String("request_id", requestID),
				//     zap.Error(err),
				//     zap.String("path", c.Request.URL.Path),
				//     zap.String("method", c.Request.Method),
				// )
			}

			// 设置响应头
			c.Header("Content-Type", "application/json")
			if requestID != "" {
				c.Header("X-Request-ID", requestID)
			}

			// 如果还没有响应，则发送错误响应
			if !c.Writer.Written() {
				c.JSON(status, response)
			}

			// 阻止后续中间件执行
			c.Abort()
		}
	}
}

// ResponseHandler 统一响应处理中间件
// 用于处理业务逻辑返回的正常响应
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果已经有响应或错误，跳过
		if c.Writer.Written() || len(c.Errors) > 0 {
			return
		}

		// 处理成功响应
		// 业务逻辑可以通过 c.Set("response", data) 设置响应数据
		if resp, exists := c.Get("response"); exists {
			c.JSON(200, gin.H{
				"code":    errors.Success,
				"message": "success",
				"data":    resp,
			})
		}
	}
}

// SuccessResponse 发送成功响应的辅助函数
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    errors.Success,
		"message": "success",
		"data":    data,
	})
}

// SuccessWithMessage 发送带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"code":    errors.Success,
		"message": message,
		"data":    data,
	})
}

// PaginatedResponse 发送分页响应
func PaginatedResponse(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(200, gin.H{
		"code":    errors.Success,
		"message": "success",
		"data":    data,
		"total":   total,
		"page":    page,
		"size":    pageSize,
	})
}

// ErrorResponse 发送错误响应的辅助函数
func ErrorResponse(c *gin.Context, err error) {
	if err == nil {
		SuccessResponse(c, nil)
		return
	}

	// 转换为应用错误
	appErr := errors.GetAppError(err)

	// 添加请求ID
	requestID := c.GetString("requestId")
	if requestID != "" {
		appErr.WithRequestID(requestID)
	}

	// 格式化错误响应
	handler := errors.DefaultErrorHandler
	status, response := handler.Handle(appErr)

	// 记录错误日志（5xx错误）
	if status >= 500 {
		// TODO: 使用结构化日志记录
	}

	// 设置响应头
	if requestID != "" {
		c.Header("X-Request-ID", requestID)
	}

	c.JSON(status, response)
}

// BadRequest 发送400错误
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, errors.New(errors.InvalidParams, message))
}

// Unauthorized 发送401错误
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, errors.NewUnauthorized(message))
}

// Forbidden 发送403错误
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, errors.NewForbidden(message))
}

// NotFound 发送404错误
func NotFound(c *gin.Context, resource string) {
	ErrorResponse(c, errors.NewNotFound(resource))
}

// InternalError 发送500错误
func InternalError(c *gin.Context, err error) {
	ErrorResponse(c, errors.NewInternal(err, ""))
}
