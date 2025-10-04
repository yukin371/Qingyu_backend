package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutConfig 超时中间件配置
type TimeoutConfig struct {
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
	Message     string        `json:"message" yaml:"message"`
	StatusCode  int           `json:"status_code" yaml:"status_code"`
	SkipPaths   []string      `json:"skip_paths" yaml:"skip_paths"`
	ErrorHandler func(*gin.Context) `json:"-" yaml:"-"`
}

// DefaultTimeoutConfig 默认超时配置
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout:    30 * time.Second,
		Message:    "请求超时",
		StatusCode: http.StatusRequestTimeout,
		SkipPaths:  []string{"/health", "/metrics"},
		ErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":      40801,
				"message":   "请求超时",
				"timestamp": time.Now().Unix(),
				"data":      nil,
			})
		},
	}
}

// Timeout 默认超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	config := DefaultTimeoutConfig()
	config.Timeout = timeout
	return TimeoutWithConfig(config)
}

// TimeoutWithConfig 带配置的超时中间件
func TimeoutWithConfig(config TimeoutConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过超时检查
		if shouldSkipTimeout(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()

		// 替换请求上下文
		c.Request = c.Request.WithContext(ctx)

		// 创建完成通道
		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			// 请求正常完成
			return
		case p := <-panicChan:
			// 处理panic
			panic(p)
		case <-ctx.Done():
			// 请求超时
			c.Header("Connection", "close")
			if config.ErrorHandler != nil {
				config.ErrorHandler(c)
			} else {
				c.JSON(config.StatusCode, gin.H{
					"code":      40801,
					"message":   config.Message,
					"timestamp": time.Now().Unix(),
					"data":      nil,
				})
			}
			c.Abort()
			return
		}
	}
}

// shouldSkipTimeout 检查是否应该跳过超时检查
func shouldSkipTimeout(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// TimeoutWriter 超时响应写入器
type TimeoutWriter struct {
	gin.ResponseWriter
	h    http.Header
	wbuf []byte
	code int
	mu   chan struct{}
}

// NewTimeoutWriter 创建超时写入器
func NewTimeoutWriter(w gin.ResponseWriter) *TimeoutWriter {
	return &TimeoutWriter{
		ResponseWriter: w,
		h:              make(http.Header),
		mu:             make(chan struct{}, 1),
	}
}

// Header 获取响应头
func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}

// Write 写入响应体
func (tw *TimeoutWriter) Write(data []byte) (int, error) {
	select {
	case <-tw.mu:
		return 0, http.ErrHandlerTimeout
	default:
		tw.wbuf = append(tw.wbuf, data...)
		return len(data), nil
	}
}

// WriteHeader 写入响应状态码
func (tw *TimeoutWriter) WriteHeader(code int) {
	select {
	case <-tw.mu:
		return
	default:
		tw.code = code
	}
}

// Flush 刷新响应
func (tw *TimeoutWriter) Flush() {
	select {
	case <-tw.mu:
		return
	default:
		// 复制头部
		for k, v := range tw.h {
			tw.ResponseWriter.Header()[k] = v
		}
		
		if tw.code != 0 {
			tw.ResponseWriter.WriteHeader(tw.code)
		}
		
		if len(tw.wbuf) > 0 {
			tw.ResponseWriter.Write(tw.wbuf)
		}
		
		if flusher, ok := tw.ResponseWriter.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

// Timeout 标记超时
func (tw *TimeoutWriter) Timeout() {
	select {
	case tw.mu <- struct{}{}:
	default:
	}
}

// CreateTimeoutMiddleware 创建超时中间件（用于中间件工厂）
func CreateTimeoutMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	timeoutConfig := DefaultTimeoutConfig()
	
	// 解析配置
	if timeout, ok := config["timeout"].(string); ok {
		if duration, err := time.ParseDuration(timeout); err == nil {
			timeoutConfig.Timeout = duration
		}
	}
	if timeoutInt, ok := config["timeout"].(int); ok {
		timeoutConfig.Timeout = time.Duration(timeoutInt) * time.Second
	}
	if message, ok := config["message"].(string); ok {
		timeoutConfig.Message = message
	}
	if statusCode, ok := config["status_code"].(int); ok {
		timeoutConfig.StatusCode = statusCode
	}
	if skipPaths, ok := config["skip_paths"].([]string); ok {
		timeoutConfig.SkipPaths = skipPaths
	}
	
	return TimeoutWithConfig(timeoutConfig), nil
}