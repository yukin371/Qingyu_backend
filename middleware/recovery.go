package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryConfig 恢复中间件配置
type RecoveryConfig struct {
	EnableStackTrace bool   `json:"enable_stack_trace" yaml:"enable_stack_trace"`
	EnableLogging    bool   `json:"enable_logging" yaml:"enable_logging"`
	StackSize        int    `json:"stack_size" yaml:"stack_size"`
	Message          string `json:"message" yaml:"message"`
	StatusCode       int    `json:"status_code" yaml:"status_code"`
	LogLevel         string `json:"log_level" yaml:"log_level"`
}

// DefaultRecoveryConfig 默认恢复配置
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		EnableStackTrace: true,
		EnableLogging:    true,
		StackSize:        4096,
		Message:          "服务器内部错误",
		StatusCode:       http.StatusInternalServerError,
		LogLevel:         "error",
	}
}

// Recovery 默认恢复中间件
func Recovery() gin.HandlerFunc {
	return RecoveryWithConfig(DefaultRecoveryConfig())
}

// RecoveryWithConfig 带配置的恢复中间件
func RecoveryWithConfig(config RecoveryConfig) gin.HandlerFunc {
	logger := initZapLogger()
	
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := make([]byte, config.StackSize)
				length := runtime.Stack(stack, false)
				stackTrace := string(stack[:length])
				
				// 记录错误日志
				if config.EnableLogging {
					logPanicError(logger, c, err, stackTrace, config)
				}
				
				// 构建错误响应
				response := buildErrorResponse(err, stackTrace, config)
				
				// 返回错误响应
				c.JSON(config.StatusCode, response)
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// logPanicError 记录panic错误
func logPanicError(logger *zap.Logger, c *gin.Context, err interface{}, stackTrace string, config RecoveryConfig) {
	// 获取请求信息
	requestID := ""
	if id, exists := c.Get("request_id"); exists {
		requestID = fmt.Sprintf("%v", id)
	}
	
	userID := ""
	if user, exists := c.Get("user"); exists {
		if userCtx, ok := user.(*UserContext); ok {
			userID = userCtx.UserID
		}
	}
	
	// 构建日志字段
	fields := []zap.Field{
		zap.String("request_id", requestID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("user_id", userID),
		zap.Any("panic", err),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}
	
	// 添加堆栈信息
	if config.EnableStackTrace {
		fields = append(fields, zap.String("stack_trace", stackTrace))
	}
	
	// 记录日志
	switch strings.ToLower(config.LogLevel) {
	case "fatal":
		logger.Fatal("Panic recovered", fields...)
	case "error":
		logger.Error("Panic recovered", fields...)
	case "warn":
		logger.Warn("Panic recovered", fields...)
	default:
		logger.Error("Panic recovered", fields...)
	}
}

// buildErrorResponse 构建错误响应
func buildErrorResponse(err interface{}, stackTrace string, config RecoveryConfig) gin.H {
	response := gin.H{
		"code":      50000,
		"message":   config.Message,
		"timestamp": time.Now().Unix(),
		"data":      nil,
	}
	
	// 在开发环境下添加详细错误信息
	if gin.Mode() == gin.DebugMode {
		response["error"] = fmt.Sprintf("%v", err)
		if config.EnableStackTrace {
			response["stack_trace"] = stackTrace
		}
	}
	
	return response
}

// CustomRecoveryWithWriter 自定义恢复中间件（兼容gin.Recovery）
func CustomRecoveryWithWriter(out io.Writer, recovery ...gin.RecoveryFunc) gin.HandlerFunc {
	if len(recovery) > 0 {
		return gin.CustomRecoveryWithWriter(out, recovery[0])
	}
	return gin.RecoveryWithWriter(out)
}

// CreateRecoveryMiddleware 创建恢复中间件（用于中间件工厂）
func CreateRecoveryMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	recoveryConfig := DefaultRecoveryConfig()
	
	// 解析配置
	if enableStackTrace, ok := config["enable_stack_trace"].(bool); ok {
		recoveryConfig.EnableStackTrace = enableStackTrace
	}
	if enableLogging, ok := config["enable_logging"].(bool); ok {
		recoveryConfig.EnableLogging = enableLogging
	}
	if stackSize, ok := config["stack_size"].(int); ok {
		recoveryConfig.StackSize = stackSize
	}
	if message, ok := config["message"].(string); ok {
		recoveryConfig.Message = message
	}
	if statusCode, ok := config["status_code"].(int); ok {
		recoveryConfig.StatusCode = statusCode
	}
	if logLevel, ok := config["log_level"].(string); ok {
		recoveryConfig.LogLevel = logLevel
	}
	
	return RecoveryWithConfig(recoveryConfig), nil
}