package builtin

import (
	"net/http"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// MiddlewareErrorKey 中间件错误在Context中的key
	MiddlewareErrorKey = "middleware_error"
)

// ErrorHandlerMiddleware 轻量级错误处理中间件
//
// 优先级: 2（内层，在所有业务中间件之后执行）
// 用途: 捕获中间件链中的错误，返回统一的错误响应格式
type ErrorHandlerMiddleware struct {
	logger *zap.Logger
}

// NewErrorHandlerMiddleware 创建新的错误处理中间件
func NewErrorHandlerMiddleware(logger *zap.Logger) *ErrorHandlerMiddleware {
	if logger == nil {
		// 如果没有提供logger，创建一个开发环境的logger
		logger, _ = zap.NewDevelopment()
	}
	return &ErrorHandlerMiddleware{
		logger: logger,
	}
}

// Name 返回中间件名称
func (m *ErrorHandlerMiddleware) Name() string {
	return "error_handler"
}

// Priority 返回执行优先级
//
// 返回3，确保错误处理在recovery之后执行
// Recovery(2) -> ErrorHandler(3) -> Security(4) -> ...
func (m *ErrorHandlerMiddleware) Priority() int {
	return 3
}

// Handler 返回Gin处理函数
func (m *ErrorHandlerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行后续处理
		c.Next()

		// 检查是否有中间件错误
		if err, exists := c.Get(MiddlewareErrorKey); exists {
			m.handleMiddlewareError(c, err)
			return
		}

		// 检查是否有其他错误
		if len(c.Errors) > 0 {
			m.handleGinErrors(c)
			return
		}
	}
}

// handleMiddlewareError 处理中间件错误
func (m *ErrorHandlerMiddleware) handleMiddlewareError(c *gin.Context, err interface{}) {
	// 记录错误日志
	m.logger.Error("Middleware error",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.Any("error", err),
	)

	// 如果还没有响应，发送统一错误格式
	if !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   "Middleware error occurred",
		})
	}
}

// handleGinErrors 处理Gin错误
func (m *ErrorHandlerMiddleware) handleGinErrors(c *gin.Context) {
	// 获取第一个错误
	err := c.Errors[0]

	// 记录错误日志
	m.logger.Error("Gin error",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.Error(err.Err),
	)

	// 如果还没有响应，发送统一错误格式
	if !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}
}

// SetMiddlewareError 在Context中设置中间件错误
//
// 这是一个辅助函数，供其他中间件使用
// 当中间件遇到错误时，可以调用此函数设置错误，然后使用c.Abort()
//
// 示例:
//
//	if someError {
//	    builtin.SetMiddlewareError(c, someError)
//	    c.Abort()
//	    return
//	}
func SetMiddlewareError(c *gin.Context, err interface{}) {
	c.Set(MiddlewareErrorKey, err)
}

// 确保ErrorHandlerMiddleware实现了Middleware接口
var _ core.Middleware = (*ErrorHandlerMiddleware)(nil)
