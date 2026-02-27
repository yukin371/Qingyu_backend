package errors

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 捕获panic并记录详细信息
				stack := debug.Stack()

				// 尝试记录日志
				if logger, exists := c.Get("logger"); exists {
					if zapLogger, ok := logger.(*zap.Logger); ok {
						zapLogger.Error("API panic recovered",
							zap.String("service", service),
							zap.String("path", c.Request.URL.Path),
							zap.String("method", c.Request.Method),
							zap.Any("error", err),
							zap.String("stack", string(stack)),
						)
					}
				}

				// 返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "内部服务器错误",
					"details": "服务器发生未预期的错误，请稍后重试",
				})
				c.Abort()
			}
		}()

		c.Next()

		// 检查是否有错误写入
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 使用MapToHTTPStatus映射错误到正确的HTTP状态码
			statusCode := MapToHTTPStatus(err.Err)
			errorMessage := GetErrorMessage(err.Err)

			c.JSON(statusCode, gin.H{
				"code":    statusCode,
				"message": errorMessage,
				"details": err.Error(),
			})
		}
	}
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(enableLogging, enableStackTrace bool) *ErrorHandler {
	return &ErrorHandler{
		enableLogging:    enableLogging,
		enableStackTrace: enableStackTrace,
	}
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	enableLogging    bool
	enableStackTrace bool
}

// Handle 处理错误并返回HTTP状态码和响应
func (h *ErrorHandler) Handle(err *UnifiedError) (int, gin.H) {
	status := err.HTTPStatus
	response := gin.H{
		"code":       err.Code,
		"message":    err.Message,
		"details":    err.Details,
		"request_id": err.ID,
	}

	// 在开发模式下添加堆栈跟踪
	if h.enableStackTrace && err.Cause != nil {
		response["cause"] = err.Cause.Error()
	}

	return status, response
}

// HandlePanic 处理panic
func (h *ErrorHandler) HandlePanic(c *gin.Context, r interface{}, service, path string) {
	// 记录详细的panic信息
	stack := debug.Stack()
	errorMsg := fmt.Sprintf("PANIC in %s at %s: %v\nStack:\n%s", service, path, r, string(stack))

	// 尝试记录日志
	if h.enableLogging {
		// 这里应该使用项目配置的logger
		// 暂时使用简单的日志记录
		fmt.Printf("[ERROR] %s\n", errorMsg)
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": "内部服务器错误",
		"details": "服务器发生未预期的错误",
	})
}

// BusinessErrorMiddleware 业务错误中间件
func BusinessErrorMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// DefaultErrorHandler 默认错误处理器
var DefaultErrorHandler = NewErrorHandler(false, false)

// GetAppError 将error转换为*UnifiedError
func GetAppError(err error) *UnifiedError {
	if err == nil {
		return nil
	}
	if ue, ok := err.(*UnifiedError); ok {
		return ue
	}
	// 如果不是UnifiedError，创建一个通用的内部错误
	return NewUnifiedError("Internal", "INTERNAL", "Unexpected error", err)
}

// New 创建一个新的错误
func New(code ErrorCode, message string) *UnifiedError {
	return &UnifiedError{
		Code:       strconv.Itoa(int(code)),
		Message:    message,
		HTTPStatus: int(GetHTTPStatus(code)),
	}
}

// NewUnauthorized 创建401未授权错误（函数形式）
func NewUnauthorized(message string) *UnifiedError {
	return UserServiceFactory.AuthError("UNAUTHORIZED", message)
}

// NewForbidden 创建403禁止访问错误（函数形式）
func NewForbidden(message string) *UnifiedError {
	return UserServiceFactory.ForbiddenError("FORBIDDEN", message)
}

// NewNotFound 创建404未找到错误（函数形式）
func NewNotFound(resource string) *UnifiedError {
	return UserServiceFactory.NotFoundError(resource, "")
}

// NewInternal 创建500内部错误（函数形式）
func NewInternal(err error, message string) *UnifiedError {
	if message == "" {
		message = "内部服务器错误"
	}
	return UserServiceFactory.InternalError("INTERNAL", message, err)
}

// NewRateLimit 创建429限流错误（函数形式）
func NewRateLimit() *UnifiedError {
	return &UnifiedError{
		Code:       strconv.Itoa(int(RateLimitExceeded)),
		Message:    GetDefaultMessage(RateLimitExceeded),
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// NewUnifiedError 创建统一错误（辅助函数）
func NewUnifiedError(service, code, message string, cause error) *UnifiedError {
	return &UnifiedError{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Cause:      cause,
		Service:    service,
	}
}
