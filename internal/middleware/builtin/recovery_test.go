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
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestRecoveryMiddleware_Name 测试中间件名称
func TestRecoveryMiddleware_Name(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)
	assert.Equal(t, "recovery", middleware.Name())
}

// TestRecoveryMiddleware_Priority 测试中间件优先级
func TestRecoveryMiddleware_Priority(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)
	assert.Equal(t, 2, middleware.Priority())
}

// TestRecoveryMiddleware_RecoverFromPanic 测试panic恢复
func TestRecoveryMiddleware_RecoverFromPanic(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/panic", func(c *gin.Context) {
		// 故意触发panic
		panic("intentional panic")
	})

	// 发送请求
	w := performRequest(router, "GET", "/panic", nil)

	// 验证响应状态码为500（服务已恢复，未崩溃）
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证响应体
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(500), response["code"])
	assert.Equal(t, "Internal Server Error", response["message"])
}

// TestRecoveryMiddleware_NormalRequest 测试正常请求不受影响
func TestRecoveryMiddleware_NormalRequest(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)

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

// TestRecoveryMiddleware_DifferentPanicTypes 测试不同类型的panic
func TestRecoveryMiddleware_DifferentPanicTypes(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		panicValue  interface{}
		description string
	}{
		{
			name:        "字符串panic",
			panicValue:  "panic string",
			description: "测试字符串类型的panic",
		},
		{
			name:        "错误panic",
			panicValue:  errors.New("panic error"),
			description: "测试错误类型的panic",
		},
		{
			name:        "整数panic",
			panicValue:  42,
			description: "测试整数类型的panic",
		},
		{
			name:        "nil panic",
			panicValue:  nil,
			description: "测试nil类型的panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := NewRecoveryMiddleware(logger)

			// 创建测试路由
			router := gin.New()
			router.Use(middleware.Handler())
			router.GET("/panic", func(c *gin.Context) {
				panic(tt.panicValue)
			})

			// 发送请求
			w := performRequest(router, "GET", "/panic", nil)

			// 验证响应状态码为500（服务已恢复）
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}

// TestRecoveryMiddleware_WithRequestID 测试带请求ID的panic恢复
func TestRecoveryMiddleware_WithRequestID(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// 创建测试路由
	router := gin.New()
	router.Use(NewRequestIDMiddleware().Handler())
	router.Use(NewRecoveryMiddleware(logger).Handler())
	router.GET("/panic", func(c *gin.Context) {
		panic("panic with request ID")
	})

	// 发送请求
	w := performRequest(router, "GET", "/panic", nil)

	// 验证响应状态码为500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证请求ID仍然存在于响应头中
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "请求ID应该存在于响应头中")
}

// TestRecoveryMiddleware_LoadConfig 测试配置加载
func TestRecoveryMiddleware_LoadConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)

	// 测试加载配置
	config := map[string]interface{}{
		"stack_size":    8192,
		"disable_print": false,
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.Equal(t, 8192, middleware.config.StackSize)
	assert.False(t, middleware.config.DisablePrint)
}

// TestRecoveryMiddleware_ValidateConfig 测试配置验证
func TestRecoveryMiddleware_ValidateConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name    string
		config  *RecoveryConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &RecoveryConfig{
				StackSize:    4096,
				DisablePrint: true,
			},
			wantErr: false,
		},
		{
			name: "负数StackSize",
			config: &RecoveryConfig{
				StackSize:    -1,
				DisablePrint: true,
			},
			wantErr: true,
		},
		{
			name: "零StackSize（有效）",
			config: &RecoveryConfig{
				StackSize:    0,
				DisablePrint: true,
			},
			wantErr: false,
		},
		{
			name:    "默认配置",
			config:  DefaultRecoveryConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &RecoveryMiddleware{
				config: tt.config,
				logger: logger,
			}
			err := middleware.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRecoveryMiddleware_CustomRecoveryFunc 测试自定义恢复函数
func TestRecoveryMiddleware_CustomRecoveryFunc(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)

	// 设置自定义恢复函数
	customRecovered := false
	middleware.SetCustomRecovery(func(c *gin.Context, recovered interface{}) {
		customRecovered = true
		c.JSON(500, gin.H{
			"code":    500,
			"message": "custom recovery",
			"panic":   recovered,
		})
	})

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// 发送请求
	w := performRequest(router, "GET", "/panic", nil)

	// 验证自定义恢复函数被调用
	assert.True(t, customRecovered, "自定义恢复函数应该被调用")

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "custom recovery", response["message"])
	assert.Equal(t, "test panic", response["panic"])
}

// TestRecoveryMiddleware_PanicInHandler 测试处理器中的panic恢复
func TestRecoveryMiddleware_PanicInHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/panic", func(c *gin.Context) {
		// 在处理器中触发panic
		panic("handler panic")
	})
	router.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试panic的请求
	w1 := performRequest(router, "GET", "/panic", nil)
	assert.Equal(t, http.StatusInternalServerError, w1.Code)

	// 测试正常的请求仍然正常工作（服务未崩溃）
	w2 := performRequest(router, "GET", "/ok", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
}

// TestRecoveryMiddleware_MiddlewareChain 测试中间件链中的panic恢复
func TestRecoveryMiddleware_MiddlewareChain(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// 创建一个会panic的中间件
	panicMiddleware := func(c *gin.Context) {
		if c.Request.URL.Path == "/panic" {
			panic("middleware panic")
		}
		c.Next()
	}

	// 创建测试路由
	// 注意：Recovery必须在panic之前注册才能捕获panic
	router := gin.New()
	router.Use(NewRequestIDMiddleware().Handler())
	router.Use(NewRecoveryMiddleware(logger).Handler())
	router.Use(panicMiddleware)
	router.GET("/panic", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/panic", nil)

	// 验证panic被恢复
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证请求ID仍然存在
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
}

// BenchmarkRecoveryMiddleware 性能测试
func BenchmarkRecoveryMiddleware(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	middleware := NewRecoveryMiddleware(logger)
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

// 辅助函数：执行HTTP请求
func performRequest(router *gin.Engine, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}
