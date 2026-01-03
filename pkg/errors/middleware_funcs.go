package errors

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
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
	// 简单的panic处理
	c.JSON(500, gin.H{
		"code":    500,
		"message": "内部服务器错误",
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
