package middleware

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/errors"
)

// ErrorMiddlewareConfig 错误中间件配置
type ErrorMiddlewareConfig struct {
	Service          string `json:"service" yaml:"service"`
	EnableLogging    bool   `json:"enable_logging" yaml:"enable_logging"`
	EnableStackTrace bool   `json:"enable_stack_trace" yaml:"enable_stack_trace"`
}

// DefaultErrorMiddlewareConfig 默认错误中间件配置
func DefaultErrorMiddlewareConfig(service string) ErrorMiddlewareConfig {
	return ErrorMiddlewareConfig{
		Service:          service,
		EnableLogging:    true,
		EnableStackTrace: false,
	}
}

// ErrorHandler 统一错误处理中间件
func ErrorHandler(service string) gin.HandlerFunc {
	return ErrorHandlerWithConfig(DefaultErrorMiddlewareConfig(service))
}

// ErrorHandlerWithConfig 带配置的错误处理中间件
func ErrorHandlerWithConfig(config ErrorMiddlewareConfig) gin.HandlerFunc {
	return errors.ErrorMiddleware(config.Service)
}

// BusinessErrorHandler 业务错误处理中间件
func BusinessErrorHandler(service string) gin.HandlerFunc {
	return BusinessErrorHandlerWithConfig(DefaultErrorMiddlewareConfig(service))
}

// BusinessErrorHandlerWithConfig 带配置的业务错误处理中间件
func BusinessErrorHandlerWithConfig(config ErrorMiddlewareConfig) gin.HandlerFunc {
	return errors.BusinessErrorMiddleware(config.Service)
}

// PanicRecovery 统一panic恢复中间件
func PanicRecovery(service string) gin.HandlerFunc {
	return PanicRecoveryWithConfig(DefaultErrorMiddlewareConfig(service))
}

// PanicRecoveryWithConfig 带配置的panic恢复中间件
func PanicRecoveryWithConfig(config ErrorMiddlewareConfig) gin.HandlerFunc {
	errorHandler := errors.NewErrorHandler(config.EnableLogging, config.EnableStackTrace)

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errorHandler.HandlePanic(c, r, config.Service, c.Request.URL.Path)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// CreateErrorMiddleware 创建错误中间件（工厂方法）
func CreateErrorMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	service, ok := config["service"].(string)
	if !ok {
		service = "unknown"
	}

	enableLogging := true
	if logging, exists := config["enable_logging"]; exists {
		if loggingBool, ok := logging.(bool); ok {
			enableLogging = loggingBool
		}
	}

	enableStackTrace := false
	if stackTrace, exists := config["enable_stack_trace"]; exists {
		if stackTraceBool, ok := stackTrace.(bool); ok {
			enableStackTrace = stackTraceBool
		}
	}

	middlewareConfig := ErrorMiddlewareConfig{
		Service:          service,
		EnableLogging:    enableLogging,
		EnableStackTrace: enableStackTrace,
	}

	return ErrorHandlerWithConfig(middlewareConfig), nil
}
