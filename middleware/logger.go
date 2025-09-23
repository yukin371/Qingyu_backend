package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig logger中间件配置
type LoggerConfig struct {
	// 是否启用彩色输出
	EnableColor bool
	// 是否记录请求体
	LogRequestBody bool
	// 是否记录响应体
	LogResponseBody bool
	// 最大请求体大小（字节）
	MaxRequestBodySize int64
	// 最大响应体大小（字节）
	MaxResponseBodySize int64
	// 跳过记录的路径
	SkipPaths []string
}

// DefaultLoggerConfig 默认配置
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		EnableColor:         true,
		LogRequestBody:      false,
		LogResponseBody:     false,
		MaxRequestBodySize:  1024 * 1024, // 1MB
		MaxResponseBodySize: 1024 * 1024, // 1MB
		SkipPaths:          []string{"/ping", "/health"},
	}
}

// responseBodyWriter 用于捕获响应体的writer
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger 创建logger中间件
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(DefaultLoggerConfig())
}

// LoggerWithConfig 使用自定义配置创建logger中间件
func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	// 创建标准输出logger
	logger := log.New(os.Stdout, "", 0)
	
	return func(c *gin.Context) {
		// 检查是否跳过记录
		path := c.Request.URL.Path
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// 记录开始时间
		start := time.Now()
		
		// 读取请求体（如果启用）
		var requestBody string
		if config.LogRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, config.MaxRequestBodySize))
			if err == nil {
				requestBody = string(bodyBytes)
				// 重新设置请求体，以便后续处理器可以读取
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 创建响应体捕获器（如果启用）
		var responseWriter *responseBodyWriter
		if config.LogResponseBody {
			responseWriter = &responseBodyWriter{
				ResponseWriter: c.Writer,
				body:          bytes.NewBufferString(""),
			}
			c.Writer = responseWriter
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		
		// 获取响应体（如果启用）
		var responseBody string
		if config.LogResponseBody && responseWriter != nil {
			responseBytes := responseWriter.body.Bytes()
			if len(responseBytes) > int(config.MaxResponseBodySize) {
				responseBody = string(responseBytes[:config.MaxResponseBodySize]) + "...[truncated]"
			} else {
				responseBody = string(responseBytes)
			}
		}

		// 构建日志信息
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// 根据状态码选择颜色
		var statusColor, methodColor, resetColor string
		if config.EnableColor {
			statusColor = getStatusColor(statusCode)
			methodColor = getMethodColor(method)
			resetColor = "\033[0m"
		}

		// 基础日志信息
		logMsg := fmt.Sprintf("[GIN] %s |%s %3d %s| %13v | %15s |%s %-7s %s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, resetColor,
			latency,
			clientIP,
			methodColor, method, resetColor,
			path,
		)

		// 添加查询参数
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			logMsg += fmt.Sprintf(" | Query: %s", rawQuery)
		}

		// 添加User-Agent
		if userAgent != "" {
			logMsg += fmt.Sprintf(" | UA: %s", userAgent)
		}

		// 添加请求体
		if config.LogRequestBody && requestBody != "" {
			logMsg += fmt.Sprintf(" | ReqBody: %s", requestBody)
		}

		// 添加响应体
		if config.LogResponseBody && responseBody != "" {
			logMsg += fmt.Sprintf(" | RespBody: %s", responseBody)
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			logMsg += fmt.Sprintf(" | Errors: %s", c.Errors.String())
		}

		// 输出日志
		logger.Println(logMsg)
	}
}

// getStatusColor 根据状态码获取颜色
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[97;42m" // 绿色背景
	case code >= 300 && code < 400:
		return "\033[90;47m" // 白色背景
	case code >= 400 && code < 500:
		return "\033[90;43m" // 黄色背景
	default:
		return "\033[97;41m" // 红色背景
	}
}

// getMethodColor 根据HTTP方法获取颜色
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[97;44m" // 蓝色背景
	case "POST":
		return "\033[97;42m" // 绿色背景
	case "PUT":
		return "\033[97;43m" // 黄色背景
	case "DELETE":
		return "\033[97;41m" // 红色背景
	case "PATCH":
		return "\033[97;42m" // 绿色背景
	case "HEAD":
		return "\033[97;45m" // 紫色背景
	case "OPTIONS":
		return "\033[90;47m" // 白色背景
	default:
		return "\033[0m" // 默认颜色
	}
}
