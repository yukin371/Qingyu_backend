package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	apperrors "Qingyu_backend/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ErrorRecoveryConfig 错误恢复配置
type ErrorRecoveryConfig struct {
	// EnableStacktrace 是否启用堆栈跟踪
	EnableStacktrace bool
	// EnableLogging 是否启用日志记录
	EnableLogging bool
	// CustomErrorHandler 自定义错误处理函数
	CustomErrorHandler func(c *gin.Context, err error)
}

// DefaultErrorRecoveryConfig 默认错误恢复配置
var DefaultErrorRecoveryConfig = ErrorRecoveryConfig{
	EnableStacktrace: false, // 生产环境关闭堆栈跟踪
	EnableLogging:    true,
}

// ErrorRecoveryMiddleware 错误恢复中间件
// 自动捕获panic并返回统一的错误响应
func ErrorRecoveryMiddleware(config ...ErrorRecoveryConfig) gin.HandlerFunc {
	var cfg ErrorRecoveryConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultErrorRecoveryConfig
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := debug.Stack()
				if cfg.EnableLogging {
					// TODO: 集成日志系统记录panic信息
					fmt.Printf("PANIC recovered: %v\n%s\n", err, stack)
				}

				// 转换为UnifiedError
				unifiedErr := convertPanicToUnifiedError(err, stack, cfg.EnableStacktrace)

				// 添加请求上下文信息
				if requestID := c.GetString("request_id"); requestID != "" {
					unifiedErr.WithRequestID(requestID)
				}
				if userID := c.GetString("user_id"); userID != "" {
					unifiedErr.WithContext(userID, c.GetString("request_id"), c.GetString("trace_id"))
				}

				// 使用自定义错误处理器或默认处理器
				if cfg.CustomErrorHandler != nil {
					cfg.CustomErrorHandler(c, unifiedErr)
				} else {
					handleUnifiedError(c, unifiedErr)
				}

				// 阻止后续中间件执行
				c.Abort()
			}
		}()

		c.Next()
	}
}

// convertPanicToUnifiedError 将panic转换为UnifiedError
func convertPanicToUnifiedError(p interface{}, stack []byte, enableStacktrace bool) *apperrors.UnifiedError {
	var err error
	var message string

	switch v := p.(type) {
	case error:
		err = v
		message = v.Error()
	case string:
		err = errors.New(v)
		message = v
	default:
		err = fmt.Errorf("%v", v)
		message = "服务器内部错误"
	}

	// 获取错误工厂
	factory := apperrors.NewErrorFactory("api-gateway")

	// 构建UnifiedError
	unifiedErr := factory.InternalError("PANIC", message, err)

	// 在开发环境中添加堆栈信息
	if enableStacktrace {
		unifiedErr.Stack = string(stack)
	}

	return unifiedErr
}

// handleUnifiedError 处理UnifiedError并返回HTTP响应
func handleUnifiedError(c *gin.Context, err *apperrors.UnifiedError) {
	statusCode := err.GetHTTPStatus()

	// 构建响应
	response := gin.H{
		"code":      statusCode,
		"message":   err.Message,
		"error":     err.Code,
		"timestamp": err.Timestamp.Unix(),
	}

	// 添加请求ID
	if err.RequestID != "" {
		response["request_id"] = err.RequestID
	} else if requestID := c.GetString("request_id"); requestID != "" {
		response["request_id"] = requestID
	}

	// 在开发环境中添加详细信息
	// TODO: 从配置中读取环境变量
	if err.Details != "" {
		response["debug"] = err.Details
	}

	c.JSON(statusCode, response)
}

// ErrorHandlerMiddleware 统一错误处理中间件
// 捕获c.Error()中存储的错误并返回统一响应
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last().Err

			// 如果还没有响应，则处理错误
			if !c.Writer.Written() {
				handleErrorWithContext(c, err)
			}
		}
	}
}

// handleErrorWithContext 根据上下文处理错误
func handleErrorWithContext(c *gin.Context, err error) {
	// 检查是否为UnifiedError
	if unifiedErr, ok := err.(*apperrors.UnifiedError); ok {
		handleUnifiedError(c, unifiedErr)
		return
	}

	// 处理常见错误类型
	errMsg := strings.ToLower(err.Error())

	// 参数错误
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "parameter") ||
		strings.Contains(errMsg, "required") {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 未授权
	if strings.Contains(errMsg, "unauthorized") ||
		strings.Contains(errMsg, "authentication") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
			"error":   err.Error(),
		})
		return
	}

	// 禁止访问
	if strings.Contains(errMsg, "forbidden") ||
		strings.Contains(errMsg, "permission") {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "禁止访问",
			"error":   err.Error(),
		})
		return
	}

	// 未找到
	if strings.Contains(errMsg, "not found") {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "资源不存在",
			"error":   err.Error(),
		})
		return
	}

	// 默认返回内部错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": "服务器内部错误",
		"error":   err.Error(),
	})
}

// WrapError 将error包装为UnifiedError
// 用于在API处理器中快速创建标准错误
func WrapError(serviceName, operation string, err error) *apperrors.UnifiedError {
	if err == nil {
		return nil
	}

	if unifiedErr, ok := err.(*apperrors.UnifiedError); ok {
		return unifiedErr.WithService(serviceName, operation)
	}

	factory := apperrors.NewErrorFactory(serviceName)
	return factory.InternalError("INTERNAL_ERROR", "操作失败", err).
		WithService(serviceName, operation)
}

// NotFoundError 创建404错误
func NotFoundError(resource, id string) *apperrors.UnifiedError {
	factory := apperrors.NewErrorFactory("api-gateway")
	return factory.NotFoundError(resource, id)
}

// ValidationError 创建参数验证错误
func ValidationError(field, message string) *apperrors.UnifiedError {
	factory := apperrors.NewErrorFactory("api-gateway")
	return factory.ValidationError("VALIDATION_ERROR", message, fmt.Sprintf("Field: %s", field))
}

// BusinessError 创建业务逻辑错误
func BusinessError(code, message string) *apperrors.UnifiedError {
	factory := apperrors.NewErrorFactory("api-gateway")
	return factory.BusinessError(code, message)
}
