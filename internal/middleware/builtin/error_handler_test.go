package builtin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestErrorHandlerMiddleware_Name 测试中间件名称
func TestErrorHandlerMiddleware_Name(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewErrorHandlerMiddleware(logger)
	assert.Equal(t, "error_handler", middleware.Name())
}

// TestErrorHandlerMiddleware_Priority 测试中间件优先级
func TestErrorHandlerMiddleware_Priority(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewErrorHandlerMiddleware(logger)
	assert.Equal(t, 3, middleware.Priority())
}

// TestErrorHandlerMiddleware_HandleGinErrors 测试捕获Gin错误
func TestErrorHandlerMiddleware_HandleGinErrors(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(observedZapCore)

	middleware := NewErrorHandlerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/error", func(c *gin.Context) {
		// 添加Gin错误
		_ = c.Error(errors.New("test error"))
	})

	// 发送请求
	w := performRequest(router, "GET", "/error", nil)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证响应体
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(http.StatusInternalServerError), response["code"])
	assert.Equal(t, "Internal Server Error", response["message"])

	// 验证错误日志被记录
	assert.Equal(t, 1, logs.FilterField(zap.String("path", "/error")).Len())
	assert.Equal(t, 1, logs.FilterField(zap.String("method", "GET")).Len())
}

// TestErrorHandlerMiddleware_HandleMiddlewareError 测试捕获中间件错误
func TestErrorHandlerMiddleware_HandleMiddlewareError(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(observedZapCore)

	middleware := NewErrorHandlerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/error", func(c *gin.Context) {
		// 设置中间件错误
		SetMiddlewareError(c, errors.New("middleware error"))
		c.Abort()
	})

	// 发送请求
	w := performRequest(router, "GET", "/error", nil)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证响应体
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(http.StatusInternalServerError), response["code"])
	assert.Equal(t, "Internal Server Error", response["message"])
	assert.Equal(t, "Middleware error occurred", response["error"])

	// 验证错误日志被记录
	assert.Equal(t, 1, logs.FilterField(zap.String("path", "/error")).Len())
}

// TestErrorHandlerMiddleware_NormalRequest 测试正常请求不受影响
func TestErrorHandlerMiddleware_NormalRequest(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewErrorHandlerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/normal", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/normal", nil)

	// 验证响应正常
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "ok", response["message"])
}

// TestErrorHandlerMiddleware_MultipleErrors 测试处理多个错误
func TestErrorHandlerMiddleware_MultipleErrors(t *testing.T) {
	// 创建测试观察器
	observedZapCore, _ := observer.New(zap.ErrorLevel)
	logger := zap.New(observedZapCore)

	middleware := NewErrorHandlerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/errors", func(c *gin.Context) {
		// 添加多个错误
		_ = c.Error(errors.New("first error"))
		_ = c.Error(errors.New("second error"))
		_ = c.Error(errors.New("third error"))
	})

	// 发送请求
	w := performRequest(router, "GET", "/errors", nil)

	// 验证响应（应该处理第一个错误）
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(http.StatusInternalServerError), response["code"])
	// Gin的错误消息包含错误类型前缀
	assert.Contains(t, response["error"], "first error")
}

// TestErrorHandlerMiddleware_ErrorAlreadyWritten 测试响应已写入时的行为
func TestErrorHandlerMiddleware_ErrorAlreadyWritten(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewErrorHandlerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/written", func(c *gin.Context) {
		// 先写入响应
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		// 然后添加错误
		_ = c.Error(errors.New("error after response"))
	})

	// 发送请求
	w := performRequest(router, "GET", "/written", nil)

	// 验证响应保持原样（不会被错误处理器覆盖）
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(400), response["code"])
	assert.Equal(t, "bad request", response["message"])
}

// TestErrorHandlerMiddleware_WithRequestID 测试带请求ID的错误处理
func TestErrorHandlerMiddleware_WithRequestID(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// 创建测试路由
	router := gin.New()
	router.Use(NewRequestIDMiddleware().Handler())
	router.Use(NewErrorHandlerMiddleware(logger).Handler())
	router.GET("/error", func(c *gin.Context) {
		_ = c.Error(errors.New("error with request ID"))
	})

	// 发送请求
	w := performRequest(router, "GET", "/error", nil)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证请求ID仍然存在于响应头中
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "请求ID应该存在于响应头中")
}

// TestErrorHandlerMiddleware_DifferentErrorTypes 测试不同类型的错误
func TestErrorHandlerMiddleware_DifferentErrorTypes(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		errorFunc   func(*gin.Context)
		description string
	}{
		{
			name: "标准错误",
			errorFunc: func(c *gin.Context) {
				_ = c.Error(errors.New("standard error"))
			},
			description: "测试标准错误类型的处理",
		},
		{
			name: "中间件错误",
			errorFunc: func(c *gin.Context) {
				SetMiddlewareError(c, "middleware error string")
			},
			description: "测试中间件错误的处理",
		},
		{
			name: "nil错误",
			errorFunc: func(c *gin.Context) {
				SetMiddlewareError(c, nil)
			},
			description: "测试nil错误的处理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := NewErrorHandlerMiddleware(logger)

			// 创建测试路由
			router := gin.New()
			router.Use(middleware.Handler())
			router.GET("/error", tt.errorFunc)

			// 发送请求
			w := performRequest(router, "GET", "/error", nil)

			// 验证响应
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}

// TestErrorHandlerMiddleware_MiddlewareChain 测试中间件链中的错误处理
func TestErrorHandlerMiddleware_MiddlewareChain(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// 创建一个会设置错误的中间件
	errorMiddleware := func(c *gin.Context) {
		if c.Request.URL.Path == "/error" {
			SetMiddlewareError(c, errors.New("chain error"))
			c.Abort()
			return
		}
		c.Next()
	}

	// 创建测试路由
	router := gin.New()
	router.Use(NewRequestIDMiddleware().Handler())
	router.Use(NewErrorHandlerMiddleware(logger).Handler())
	router.Use(errorMiddleware)
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})
	router.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试错误的请求
	w1 := performRequest(router, "GET", "/error", nil)
	assert.Equal(t, http.StatusInternalServerError, w1.Code)

	// 测试正常的请求
	w2 := performRequest(router, "GET", "/ok", nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 验证请求ID仍然存在
	requestID := w2.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
}

// TestErrorHandlerMiddleware_PanicError 测试panic后的错误处理
func TestErrorHandlerMiddleware_PanicError(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// 创建测试路由
	// Recovery应该先捕获panic，然后ErrorHandler处理后续错误
	router := gin.New()
	router.Use(NewRecoveryMiddleware(logger).Handler())
	router.Use(NewErrorHandlerMiddleware(logger).Handler())
	router.GET("/panic", func(c *gin.Context) {
		panic("intentional panic")
	})

	// 发送请求
	w := performRequest(router, "GET", "/panic", nil)

	// 验证panic被Recovery处理，返回500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// BenchmarkErrorHandlerMiddleware 性能测试
func BenchmarkErrorHandlerMiddleware(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	middleware := NewErrorHandlerMiddleware(logger)
	handler := middleware.Handler()

	// 创建测试路由
	router := gin.New()
	router.Use(handler)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
