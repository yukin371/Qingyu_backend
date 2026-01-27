package builtin

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestRequestIDMiddleware_Name 测试中间件名称
func TestRequestIDMiddleware_Name(t *testing.T) {
	middleware := NewRequestIDMiddleware()
	assert.Equal(t, "request_id", middleware.Name())
}

// TestRequestIDMiddleware_Priority 测试中间件优先级
func TestRequestIDMiddleware_Priority(t *testing.T) {
	middleware := NewRequestIDMiddleware()
	assert.Equal(t, 1, middleware.Priority())
}

// TestRequestIDMiddleware_GenerateNewID 测试生成新请求ID
func TestRequestIDMiddleware_GenerateNewID(t *testing.T) {
	middleware := NewRequestIDMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	// 发送请求（不带请求ID）
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证请求ID已生成
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "请求ID应该被生成")

	// 验证响应体中的请求ID
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, requestID, response["request_id"])
}

// TestRequestIDMiddleware_UseExistingID 测试使用已有的请求ID
func TestRequestIDMiddleware_UseExistingID(t *testing.T) {
	middleware := NewRequestIDMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	// 发送请求（带请求ID）
	existingRequestID := "test-request-id-123"
	w := performRequest(router, "GET", "/test", map[string]string{
		"X-Request-ID": existingRequestID,
	})

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证使用的是已存在的请求ID
	requestID := w.Header().Get("X-Request-ID")
	assert.Equal(t, existingRequestID, requestID, "应该使用已存在的请求ID")
}

// TestRequestIDMiddleware_ForceGen 测试强制生成新ID
func TestRequestIDMiddleware_ForceGen(t *testing.T) {
	middleware := NewRequestIDMiddleware()
	middleware.config = &RequestIDConfig{
		HeaderName: "X-Request-ID",
		ForceGen:   true, // 强制生成
	}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	// 发送请求（带请求ID，但会被忽略）
	existingRequestID := "test-request-id-123"
	w := performRequest(router, "GET", "/test", map[string]string{
		"X-Request-ID": existingRequestID,
	})

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证生成的是新ID（不是已存在的ID）
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEqual(t, existingRequestID, requestID, "应该生成新的请求ID")
	assert.NotEmpty(t, requestID, "请求ID应该被生成")
}

// TestRequestIDMiddleware_CustomHeaderName 测试自定义头名称
func TestRequestIDMiddleware_CustomHeaderName(t *testing.T) {
	middleware := NewRequestIDMiddleware()
	middleware.config = &RequestIDConfig{
		HeaderName: "X-Trace-ID",
		ForceGen:   false,
	}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	// 发送请求（使用自定义头名称）
	existingTraceID := "trace-id-456"
	w := performRequest(router, "GET", "/test", map[string]string{
		"X-Trace-ID": existingTraceID,
	})

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证使用自定义头名称
	traceID := w.Header().Get("X-Trace-ID")
	assert.Equal(t, existingTraceID, traceID, "应该使用自定义头名称")

	// 验证默认头名称没有值
	defaultRequestID := w.Header().Get("X-Request-ID")
	assert.Empty(t, defaultRequestID, "默认头名称不应该有值")
}

// TestRequestIDMiddleware_LoadConfig 测试配置加载
func TestRequestIDMiddleware_LoadConfig(t *testing.T) {
	middleware := NewRequestIDMiddleware()

	// 测试加载配置
	config := map[string]interface{}{
		"header_name": "X-Custom-Request-ID",
		"force_gen":   true,
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.Equal(t, "X-Custom-Request-ID", middleware.config.HeaderName)
	assert.True(t, middleware.config.ForceGen)
}

// TestRequestIDMiddleware_ValidateConfig 测试配置验证
func TestRequestIDMiddleware_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *RequestIDConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &RequestIDConfig{
				HeaderName: "X-Request-ID",
				ForceGen:   false,
			},
			wantErr: false,
		},
		{
			name: "空头名称",
			config: &RequestIDConfig{
				HeaderName: "",
				ForceGen:   false,
			},
			wantErr: true,
		},
		{
			name:    "默认配置",
			config:  DefaultRequestIDConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &RequestIDMiddleware{config: tt.config}
			err := middleware.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetRequestID 测试GetRequestID辅助函数
func TestGetRequestID(t *testing.T) {
	router := gin.New()
	router.Use(NewRequestIDMiddleware().Handler())
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	// 测试正常情况
	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.NotEmpty(t, response["request_id"])

	// 测试没有设置请求ID的情况
	router2 := gin.New()
	router2.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	w2 := performRequest(router2, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	var response2 map[string]interface{}
	assert.NoError(t, json.Unmarshal(w2.Body.Bytes(), &response2))
	assert.Empty(t, response2["request_id"], "没有设置请求ID时应该返回空字符串")
}

// TestRequestIDMiddleware_MiddlewareChain 测试中间件链中的请求ID传递
func TestRequestIDMiddleware_MiddlewareChain(t *testing.T) {
	requestIDMiddleware := NewRequestIDMiddleware()

	// 创建一个自定义中间件来验证请求ID是否正确传递
	var capturedRequestID string
	customMiddleware := func(c *gin.Context) {
		capturedRequestID = GetRequestID(c)
		c.Next()
	}

	// 创建测试路由
	router := gin.New()
	router.Use(requestIDMiddleware.Handler())
	router.Use(customMiddleware)
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{
			"middleware_request_id": capturedRequestID,
			"handler_request_id":    requestID,
		})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证请求ID在整个中间件链中一致
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, response["middleware_request_id"], response["handler_request_id"])
	assert.NotEmpty(t, response["middleware_request_id"])
}

// BenchmarkRequestIDMiddleware 性能测试
func BenchmarkRequestIDMiddleware(b *testing.B) {
	middleware := NewRequestIDMiddleware()
	handler := middleware.Handler()

	// 创建测试Context
	c, _ := gin.CreateTestContext(nil)
	c.Request = nil
	c.Writer = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(c)
	}
}
