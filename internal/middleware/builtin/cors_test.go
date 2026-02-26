package builtin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestCORSMiddleware_Name 测试中间件名称
func TestCORSMiddleware_Name(t *testing.T) {
	middleware := NewCORSMiddleware()
	assert.Equal(t, "cors", middleware.Name())
}

// TestCORSMiddleware_Priority 测试中间件优先级
func TestCORSMiddleware_Priority(t *testing.T) {
	middleware := NewCORSMiddleware()
	assert.Equal(t, 5, middleware.Priority())
}

// TestCORSMiddleware_PreflightRequest 测试预检请求处理
func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.OPTIONS("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建预检请求
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证预检请求返回204
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 验证CORS响应头
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
}

// TestCORSMiddleware_WildcardOrigin 测试通配符origin配置
func TestCORSMiddleware_WildcardOrigin(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试多个不同的origin
	origins := []string{
		"https://example.com",
		"https://api.example.com",
		"http://localhost:8080",
	}

	for _, origin := range origins {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", origin)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 验证响应正常
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证CORS头（默认配置使用通配符，但要求凭证时返回具体origin）
		assert.Equal(t, origin, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	}
}

// TestCORSMiddleware_AllowCredentials 测试凭证配置
func TestCORSMiddleware_AllowCredentials(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送带Origin的请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证凭证头被设置
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))

	// 验证返回具体的origin而不是通配符
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSMiddleware_CustomAllowedOrigins 测试自定义允许的源列表
func TestCORSMiddleware_CustomAllowedOrigins(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 加载自定义配置
	config := map[string]interface{}{
		"allowed_origins":   []interface{}{"https://allowed.com", "https://also-allowed.com"},
		"allow_credentials": false,
	}
	err := middleware.LoadConfig(config)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	tests := []struct {
		name           string
		origin         string
		expectedStatus int
		expectedAllow  string
	}{
		{
			name:           "允许的源",
			origin:         "https://allowed.com",
			expectedStatus: http.StatusOK,
			expectedAllow:  "https://allowed.com",
		},
		{
			name:           "另一个允许的源",
			origin:         "https://also-allowed.com",
			expectedStatus: http.StatusOK,
			expectedAllow:  "https://also-allowed.com",
		},
		{
			name:           "不允许的源",
			origin:         "https://not-allowed.com",
			expectedStatus: http.StatusForbidden,
			expectedAllow:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tt.origin)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应状态码
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 验证CORS头
			if tt.expectedAllow != "" {
				assert.Equal(t, tt.expectedAllow, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

// TestCORSMiddleware_ResponseHeaders 测试CORS响应头设置
func TestCORSMiddleware_ResponseHeaders(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证所有CORS响应头
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Expose-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Max-Age"))

	// 验证暴露的响应头包含预期的头
	exposeHeaders := w.Header().Get("Access-Control-Expose-Headers")
	assert.Contains(t, exposeHeaders, "X-Request-ID")
	assert.Contains(t, exposeHeaders, "X-Response-Time")
}

// TestCORSMiddleware_NoOrigin 测试没有Origin头的请求
func TestCORSMiddleware_NoOrigin(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送不带Origin的请求
	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证请求正常通过（没有Origin不是同源请求，不需要CORS）
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestCORSMiddleware_LoadConfig 测试配置加载
func TestCORSMiddleware_LoadConfig(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 测试加载配置
	config := map[string]interface{}{
		"allowed_origins":   []interface{}{"https://example.com", "https://api.example.com"},
		"allowed_methods":   []interface{}{"GET", "POST", "PUT"},
		"allowed_headers":   []interface{}{"Content-Type", "Authorization"},
		"exposed_headers":   []interface{}{"X-Custom-Header"},
		"allow_credentials": true,
		"max_age":           3600,
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.Equal(t, []string{"https://example.com", "https://api.example.com"}, middleware.config.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST", "PUT"}, middleware.config.AllowedMethods)
	assert.Equal(t, []string{"Content-Type", "Authorization"}, middleware.config.AllowedHeaders)
	assert.Equal(t, []string{"X-Custom-Header"}, middleware.config.ExposedHeaders)
	assert.True(t, middleware.config.AllowCredentials)
	assert.Equal(t, 3600, middleware.config.MaxAge)
}

// TestCORSMiddleware_ValidateConfig 测试配置验证
func TestCORSMiddleware_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *CORSConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &CORSConfig{
				AllowedOrigins:   []string{"https://example.com"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: false,
				MaxAge:           86400,
			},
			wantErr: false,
		},
		{
			name: "空AllowedOrigins",
			config: &CORSConfig{
				AllowedOrigins:   []string{},
				AllowedMethods:   []string{"GET"},
				AllowCredentials: false,
			},
			wantErr: true,
		},
		{
			name: "AllowCredentials为true时使用通配符",
			config: &CORSConfig{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET"},
				AllowCredentials: true,
			},
			wantErr: true,
		},
		{
			name: "空AllowedMethods",
			config: &CORSConfig{
				AllowedOrigins:   []string{"https://example.com"},
				AllowedMethods:   []string{},
				AllowCredentials: false,
			},
			wantErr: true,
		},
		{
			name: "负数MaxAge",
			config: &CORSConfig{
				AllowedOrigins:   []string{"https://example.com"},
				AllowedMethods:   []string{"GET"},
				AllowCredentials: false,
				MaxAge:           -1,
			},
			wantErr: true,
		},
		{
			name: "自定义配置（AllowCredentials为true且不使用通配符）",
			config: &CORSConfig{
				AllowedOrigins:   []string{"https://example.com"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
				MaxAge:           86400,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &CORSMiddleware{
				config: tt.config,
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

// TestCORSMiddleware_AllowCredentialsWithWildcard 测试AllowCredentials与通配符的兼容性验证
func TestCORSMiddleware_AllowCredentialsWithWildcard(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 尝试设置不兼容的配置
	middleware.config.AllowCredentials = true
	middleware.config.AllowedOrigins = []string{"*"}

	// 验证应该失败
	err := middleware.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "通配符")
}

// TestCORSMiddleware_NonPreflightWithCustomHeaders 测试非预检请求的自定义头
func TestCORSMiddleware_NonPreflightWithCustomHeaders(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送带自定义请求头的请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("Authorization", "Bearer token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应正常
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证CORS头存在
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSMiddleware_MethodsNotAllowed 测试不允许的HTTP方法
func TestCORSMiddleware_MethodsNotAllowed(t *testing.T) {
	middleware := NewCORSMiddleware()

	// 加载只允许GET的配置
	config := map[string]interface{}{
		"allowed_origins": []interface{}{"*"},
		"allowed_methods": []interface{}{"GET"},
	}
	err := middleware.LoadConfig(config)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送POST请求（不在允许的方法列表中）
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 注意：CORS中间件本身不阻止请求，只是设置响应头
	// 实际的方法限制由路由处理
	// 这里我们只验证CORS头被正确设置
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
}

// BenchmarkCORSMiddleware 性能测试
func BenchmarkCORSMiddleware(b *testing.B) {
	middleware := NewCORSMiddleware()
	handler := middleware.Handler()

	// 创建测试路由
	router := gin.New()
	router.Use(handler)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
