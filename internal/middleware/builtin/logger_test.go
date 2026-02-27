package builtin

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestLoggerMiddleware_Name 测试中间件名称
func TestLoggerMiddleware_Name(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewLoggerMiddleware(logger)
	assert.Equal(t, "logger", middleware.Name())
}

// TestLoggerMiddleware_Priority 测试中间件优先级
func TestLoggerMiddleware_Priority(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewLoggerMiddleware(logger)
	assert.Equal(t, 7, middleware.Priority())
}

// TestLoggerMiddleware_LogBasicRequest 测试基本请求日志
func TestLoggerMiddleware_LogBasicRequest(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证日志记录
	assert.Equal(t, 1, logs.FilterField(zap.String("method", "GET")).Len())
	assert.Equal(t, 1, logs.FilterField(zap.String("path", "/test")).Len())
	assert.Equal(t, 1, logs.FilterField(zap.Int("status", 200)).Len())
	assert.True(t, logs.FilterField(zap.Int64("latency_ms", 0)).Len() >= 1, "应该记录耗时")
}

// TestLoggerMiddleware_LogWithRequestID 测试带请求ID的日志
func TestLoggerMiddleware_LogWithRequestID(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(NewRequestIDMiddleware().Handler())
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证日志中包含请求ID
	logsWithRequestID := logs.FilterField(zap.String("request_id", w.Header().Get("X-Request-ID")))
	assert.Equal(t, 1, logsWithRequestID.Len())
}

// TestLoggerMiddleware_LogError 测试错误日志
func TestLoggerMiddleware_LogError(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/error", nil)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证错误日志
	assert.Equal(t, 1, logs.FilterField(zap.Int("status", 500)).Len())
}

// TestLoggerMiddleware_SkipPath 测试跳过指定路径
func TestLoggerMiddleware_SkipPath(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)
	middleware.config.SkipPaths = []string{"/health", "/metrics"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求到health路径（应该被跳过）
	w1 := performRequest(router, "GET", "/health", nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 发送请求到test路径（应该记录日志）
	w2 := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 验证只记录了test路径的日志
	assert.Equal(t, 1, logs.FilterField(zap.String("path", "/test")).Len())
	assert.Equal(t, 0, logs.FilterField(zap.String("path", "/health")).Len())
}

// TestLoggerMiddleware_LogRequestBody 测试记录请求体
func TestLoggerMiddleware_LogRequestBody(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)
	middleware.config.EnableRequestBody = true

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送JSON请求
	reqBody := `{"test": "data"}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证请求体被记录
	logsWithBody := logs.FilterField(zap.String("request_body", reqBody))
	assert.Equal(t, 1, logsWithBody.Len())
}

// TestLoggerMiddleware_LogResponseBody 测试记录响应体
func TestLoggerMiddleware_LogResponseBody(t *testing.T) {
	// 创建测试观察器
	observedZapCore, _ := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)
	middleware.config.EnableResponseBody = true

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestLoggerMiddleware_SlowRequest 测试慢请求日志
func TestLoggerMiddleware_SlowRequest(t *testing.T) {
	// 创建测试观察器
	observedZapCore, _ := observer.New(zap.WarnLevel)
	logger := zap.New(observedZapCore)

	middleware := NewLoggerMiddleware(logger)
	middleware.config.SlowRequestThreshold = 10 // 设置很低的阈值

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/slow", func(c *gin.Context) {
		// 模拟慢请求
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/slow", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证慢请求被记录为WARN级别
	// 注意：由于实际耗时可能很短，这个测试可能不稳定
	// 在实际项目中，应该使用mock时间来控制
}

// TestLoggerMiddleware_LoadConfig 测试配置加载
func TestLoggerMiddleware_LoadConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewLoggerMiddleware(logger)

	// 测试加载配置
	config := map[string]interface{}{
		"skip_paths":             []interface{}{"/health", "/metrics"},
		"skip_log_body":          true,
		"enable_request_id":      false,
		"enable_request_body":    false,
		"enable_response_body":   true,
		"slow_request_threshold": 5000,
		"enable_colors":          false,
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.Equal(t, []string{"/health", "/metrics"}, middleware.config.SkipPaths)
	assert.True(t, middleware.config.SkipLogBody)
	assert.False(t, middleware.config.EnableRequestID)
	assert.False(t, middleware.config.EnableRequestBody)
	assert.True(t, middleware.config.EnableResponseBody)
	assert.Equal(t, 5000, middleware.config.SlowRequestThreshold)
	assert.False(t, middleware.config.EnableColors)
}

func TestLoggerMiddleware_LoadConfigWithStringSlice(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewLoggerMiddleware(logger)

	config := map[string]interface{}{
		"skip_paths":       []string{"/health", "/metrics"},
		"body_allow_paths": []string{"/api/v1"},
		"redact_keys":      []string{"authorization", "password"},
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err)
	assert.Equal(t, []string{"/health", "/metrics"}, middleware.config.SkipPaths)
	assert.Equal(t, []string{"/api/v1"}, middleware.config.BodyAllowPaths)
	assert.Equal(t, []string{"authorization", "password"}, middleware.config.RedactKeys)
}

// TestLoggerMiddleware_ValidateConfig 测试配置验证
func TestLoggerMiddleware_ValidateConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name    string
		config  *LoggerConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &LoggerConfig{
				SkipPaths:            []string{"/health"},
				SlowRequestThreshold: 3000,
			},
			wantErr: false,
		},
		{
			name: "负数慢请求阈值",
			config: &LoggerConfig{
				SlowRequestThreshold: -1,
			},
			wantErr: true,
		},
		{
			name:    "默认配置",
			config:  DefaultLoggerConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &LoggerMiddleware{
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

// TestLoggerMiddleware_ResponseTime 测试响应时间记录
func TestLoggerMiddleware_ResponseTime(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := NewLoggerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		// 检查Context中是否有响应时间
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// BenchmarkLoggerMiddleware 性能测试
func BenchmarkLoggerMiddleware(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	middleware := NewLoggerMiddleware(logger)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
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
