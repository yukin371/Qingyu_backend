package validation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestValidationMiddleware_Name 测试中间件名称
func TestValidationMiddleware_Name(t *testing.T) {
	middleware := NewValidationMiddleware()
	assert.Equal(t, "validation", middleware.Name())
}

// TestValidationMiddleware_Priority 测试中间件优先级
func TestValidationMiddleware_Priority(t *testing.T) {
	middleware := NewValidationMiddleware()
	assert.Equal(t, 11, middleware.Priority(), "Validation中间件应该在业务层，优先级为11")
}

// TestValidationMiddleware_HandlerExists 测试Handler函数存在
func TestValidationMiddleware_HandlerExists(t *testing.T) {
	middleware := NewValidationMiddleware()
	handler := middleware.Handler()
	assert.NotNil(t, handler, "Handler函数不应该为nil")
}

// TestValidationMiddleware_DefaultConfig 测试默认配置
func TestValidationMiddleware_DefaultConfig(t *testing.T) {
	middleware := NewValidationMiddleware()

	assert.NotNil(t, middleware.config, "配置不应该为nil")
	assert.True(t, middleware.config.Enabled, "默认应该启用validation")
	assert.NotEmpty(t, middleware.config.AllowedContentTypes, "默认应该有允许的Content-Type列表")
}

// TestValidationMiddleware_ValidContentType 测试有效的Content-Type
func TestValidationMiddleware_ValidContentType(t *testing.T) {
	middleware := NewValidationMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送有效的Content-Type请求
	body := `{"name": "test"}`
	w := performRequestWithBody(router, "POST", "/test", body, "application/json")
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestValidationMiddleware_InvalidContentType 测试无效的Content-Type
func TestValidationMiddleware_InvalidContentType(t *testing.T) {
	middleware := NewValidationMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送无效的Content-Type请求
	body := `{"name": "test"}`
	w := performRequestWithBody(router, "POST", "/test", body, "text/plain")

	// 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证错误响应格式
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1001", response["code"], "应该返回InvalidParams错误码")
}

// TestValidationMiddleware_RequestSize 测试请求体大小限制
func TestValidationMiddleware_RequestSize(t *testing.T) {
	maxSize := int64(1024) // 1KB
	middleware := NewValidationMiddleware()
	middleware.config.MaxBodySize = maxSize

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送超过大小限制的请求
	largeBody := strings.Repeat("a", int(maxSize)+1)
	w := performRequestWithBody(router, "POST", "/test", largeBody, "application/json")

	// 应该返回400错误（请求体过大属于参数验证错误）
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证错误响应格式
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1001", response["code"], "应该返回InvalidParams错误码")
}

// TestValidationMiddleware_RequestSizeWithinLimit 测试请求体在限制内
func TestValidationMiddleware_RequestSizeWithinLimit(t *testing.T) {
	maxSize := int64(1024) // 1KB
	middleware := NewValidationMiddleware()
	middleware.config.MaxBodySize = maxSize

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送在限制内的有效JSON请求
	smallBody := `{"data": "` + strings.Repeat("a", 50) + `"}`
	w := performRequestWithBody(router, "POST", "/test", smallBody, "application/json")

	// 应该成功
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestValidationMiddleware_RequiredQueryParams 测试必填查询参数
func TestValidationMiddleware_RequiredQueryParams(t *testing.T) {
	middleware := NewValidationMiddleware()
	middleware.config.RequiredQueryParams = []string{"id", "name"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试缺少必填参数
	w := performRequest(router, "/test", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证错误响应
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1001", response["code"])

	// 测试提供所有必填参数
	w2 := performRequest(router, "/test?id=123&name=test", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
}

// TestValidationMiddleware_RequiredFields 测试必填字段（JSON body）
func TestValidationMiddleware_RequiredFields(t *testing.T) {
	middleware := NewValidationMiddleware()
	middleware.config.RequiredFields = []string{"name", "email"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试缺少必填字段
	body := `{"name": "test"}`
	w := performRequestWithBody(router, "POST", "/test", body, "application/json")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证错误响应
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1001", response["code"])

	// 测试提供所有必填字段
	body2 := `{"name": "test", "email": "test@example.com"}`
	w2 := performRequestWithBody(router, "POST", "/test", body2, "application/json")
	assert.Equal(t, http.StatusOK, w2.Code)
}

// TestValidationMiddleware_InvalidJSON 测试无效的JSON
func TestValidationMiddleware_InvalidJSON(t *testing.T) {
	middleware := NewValidationMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送无效的JSON
	w := performRequestWithBody(router, "POST", "/test", "invalid json", "application/json")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证错误响应
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1001", response["code"])
}

// TestValidationMiddleware_LoadConfig 测试配置加载
func TestValidationMiddleware_LoadConfig(t *testing.T) {
	middleware := NewValidationMiddleware()

	config := map[string]interface{}{
		"enabled": false,
		"max_body_size": int64(2048),
		"allowed_content_types": []interface{}{"application/json", "application/xml"},
		"required_query_params": []interface{}{"id", "token"},
		"required_fields": []interface{}{"name", "email"},
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, int64(2048), middleware.config.MaxBodySize)
	assert.Equal(t, 2, len(middleware.config.AllowedContentTypes))
	assert.Equal(t, 2, len(middleware.config.RequiredQueryParams))
	assert.Equal(t, 2, len(middleware.config.RequiredFields))
}

// TestValidationMiddleware_ValidateConfig 测试配置验证
func TestValidationMiddleware_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *ValidationConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &ValidationConfig{
				Enabled:              true,
				MaxBodySize:          1024 * 1024,
				AllowedContentTypes:  []string{"application/json"},
				RequiredQueryParams:  []string{},
				RequiredFields:       []string{},
			},
			wantErr: false,
		},
		{
			name: "负数请求体大小",
			config: &ValidationConfig{
				Enabled:              true,
				MaxBodySize:          -1,
				AllowedContentTypes:  []string{"application/json"},
				RequiredQueryParams:  []string{},
				RequiredFields:       []string{},
			},
			wantErr: true,
		},
		{
			name: "空Content-Type列表",
			config: &ValidationConfig{
				Enabled:              true,
				MaxBodySize:          1024 * 1024,
				AllowedContentTypes:  []string{},
				RequiredQueryParams:  []string{},
				RequiredFields:       []string{},
			},
			wantErr: true,
		},
		{
			name:    "默认配置",
			config:  DefaultValidationConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &ValidationMiddleware{config: tt.config}
			err := middleware.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidationMiddleware_DisabledValidation 测试禁用验证
func TestValidationMiddleware_DisabledValidation(t *testing.T) {
	middleware := NewValidationMiddleware()
	middleware.config.Enabled = false
	middleware.config.RequiredQueryParams = []string{"id"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 即使没有必填参数，也应该成功（验证已禁用）
	w := performRequest(router, "/test", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestValidationMiddleware_AllowAllContentTypes 测试允许所有Content-Type
func TestValidationMiddleware_AllowAllContentTypes(t *testing.T) {
	middleware := NewValidationMiddleware()
	middleware.config.AllowedContentTypes = []string{"*"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 任何Content-Type都应该被接受
	w := performRequestWithBody(router, "POST", "/test", "test", "text/plain")
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestValidationMiddleware_MultipleRequiredFields 测试多个必填字段
func TestValidationMiddleware_MultipleRequiredFields(t *testing.T) {
	middleware := NewValidationMiddleware()
	middleware.config.RequiredFields = []string{"name", "email", "age"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试缺少部分字段
	body := `{"name": "test"}`
	w := performRequestWithBody(router, "POST", "/test", body, "application/json")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证错误消息包含缺失的字段名
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "email")
	assert.Contains(t, response["message"], "age")

	// 测试提供所有字段
	body2 := `{"name": "test", "email": "test@example.com", "age": 25}`
	w2 := performRequestWithBody(router, "POST", "/test", body2, "application/json")
	assert.Equal(t, http.StatusOK, w2.Code)
}

// TestValidationMiddleware_ErrorResponseFormat 测试错误响应格式
func TestValidationMiddleware_ErrorResponseFormat(t *testing.T) {
	middleware := NewValidationMiddleware()
	middleware.config.AllowedContentTypes = []string{"application/json"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送无效的Content-Type
	w := performRequestWithBody(router, "POST", "/test", "{}", "text/plain")

	// 验证错误响应格式
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证错误码存在且是字符串（因为UnifiedError.Code是string类型）
	code, ok := response["code"].(string)
	assert.True(t, ok, "code字段应该存在且是字符串")
	assert.NotEmpty(t, code, "错误码不应该为空")

	// 验证错误码是4位数字字符串
	assert.GreaterOrEqual(t, len(code), 4, "错误码应该是至少4位")
	assert.LessOrEqual(t, len(code), 5, "错误码应该最多5位")

	// 验证message字段存在
	_, ok = response["message"].(string)
	assert.True(t, ok, "message字段应该存在且是字符串")
}

// TestDefaultValidationConfig 测试默认配置
func TestDefaultValidationConfig(t *testing.T) {
	config := DefaultValidationConfig()

	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Greater(t, config.MaxBodySize, int64(0))
	assert.NotEmpty(t, config.AllowedContentTypes)
}

// TestValidationMiddleware_EmptyBody 测试空请求体
func TestValidationMiddleware_EmptyBody(t *testing.T) {
	middleware := NewValidationMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送空请求体
	w := performRequestWithBody(router, "POST", "/test", "", "application/json")

	// 应该成功（空请求体是允许的）
	assert.Equal(t, http.StatusOK, w.Code)
}

// performRequest 辅助函数：执行HTTP请求
func performRequest(router *gin.Engine, url string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, nil)

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// performRequestWithMethod 辅助函数：执行指定方法的HTTP请求
func performRequestWithMethod(router *gin.Engine, method, url string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// performRequestWithBody 辅助函数：执行带body的HTTP请求
func performRequestWithBody(router *gin.Engine, method, path, body, contentType string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// BenchmarkValidationMiddleware 性能测试
func BenchmarkValidationMiddleware(b *testing.B) {
	middleware := NewValidationMiddleware()
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
