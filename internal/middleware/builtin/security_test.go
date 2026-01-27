package builtin

import (
	"crypto/tls"
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

// TestSecurityMiddleware_Name 测试中间件名称
func TestSecurityMiddleware_Name(t *testing.T) {
	middleware := NewSecurityMiddleware()
	assert.Equal(t, "security", middleware.Name())
}

// TestSecurityMiddleware_Priority 测试中间件优先级
func TestSecurityMiddleware_Priority(t *testing.T) {
	middleware := NewSecurityMiddleware()
	assert.Equal(t, 3, middleware.Priority())
}

// TestSecurityMiddleware_XFrameOptions 测试X-Frame-Options头
func TestSecurityMiddleware_XFrameOptions(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
}

// TestSecurityMiddleware_XContentTypeOptions 测试X-Content-Type-Options头
func TestSecurityMiddleware_XContentTypeOptions(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}

// TestSecurityMiddleware_XXSSProtection 测试X-XSS-Protection头
func TestSecurityMiddleware_XXSSProtection(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}

// TestSecurityMiddleware_HSTS 测试Strict-Transport-Security头
func TestSecurityMiddleware_HSTS(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建带TLS的请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.TLS = &tls.ConnectionState{} // 模拟HTTPS请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应头
	hsts := w.Header().Get("Strict-Transport-Security")
	assert.Contains(t, hsts, "max-age=31536000")
	assert.Contains(t, hsts, "includeSubDomains")
}

// TestSecurityMiddleware_HSTS_HTTP 测试HTTP请求不设置HSTS
func TestSecurityMiddleware_HSTS_HTTP(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建HTTP请求（没有TLS）
	w := performRequest(router, "GET", "/test", nil)

	// 验证没有HSTS头
	hsts := w.Header().Get("Strict-Transport-Security")
	assert.Empty(t, hsts, "HTTP请求不应该设置HSTS头")
}

// TestSecurityMiddleware_ReferrerPolicy 测试Referrer-Policy头
func TestSecurityMiddleware_ReferrerPolicy(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
}

// TestSecurityMiddleware_PermissionsPolicy 测试Permissions-Policy头
func TestSecurityMiddleware_PermissionsPolicy(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头
	permissionsPolicy := w.Header().Get("Permissions-Policy")
	assert.NotEmpty(t, permissionsPolicy)
	assert.Contains(t, permissionsPolicy, "geolocation=()")
}

// TestSecurityMiddleware_DisableXFrameOptions 测试禁用X-Frame-Options
func TestSecurityMiddleware_DisableXFrameOptions(t *testing.T) {
	middleware := NewSecurityMiddleware()
	middleware.config.EnableXFrameOptions = false

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头不存在
	xFrameOptions := w.Header().Get("X-Frame-Options")
	assert.Empty(t, xFrameOptions)
}

// TestSecurityMiddleware_CSP 测试Content-Security-Policy头
func TestSecurityMiddleware_CSP(t *testing.T) {
	middleware := NewSecurityMiddleware()
	middleware.config.EnableCSP = true
	middleware.config.CSPContent = "default-src 'self'"

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证响应头
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
}

// TestSecurityMiddleware_CSPReportOnly 测试Content-Security-Policy-Report-Only头
func TestSecurityMiddleware_CSPReportOnly(t *testing.T) {
	middleware := NewSecurityMiddleware()
	middleware.config.EnableCSP = true
	middleware.config.CSPContent = "default-src 'self'"
	middleware.config.EnableContentSecurityPolicyReportOnly = true

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证Report-Only头存在
	cspReportOnly := w.Header().Get("Content-Security-Policy-Report-Only")
	assert.NotEmpty(t, cspReportOnly)

	// 验证普通的CSP头不存在
	csp := w.Header().Get("Content-Security-Policy")
	assert.Empty(t, csp)
}

// TestSecurityMiddleware_CustomHeaders 测试自定义安全头
func TestSecurityMiddleware_CustomHeaders(t *testing.T) {
	middleware := NewSecurityMiddleware()
	middleware.config.CustomHeaders = map[string]string{
		"X-Custom-Header": "custom-value",
		"X-Another-Header": "another-value",
	}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test", nil)

	// 验证自定义头
	assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))
	assert.Equal(t, "another-value", w.Header().Get("X-Another-Header"))
}

// TestSecurityMiddleware_LoadConfig 测试配置加载
func TestSecurityMiddleware_LoadConfig(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 测试加载配置
	config := map[string]interface{}{
		"enable_x_frame_options": false,
		"x_frame_options":        "SAMEORIGIN",
		"enable_hsts":            true,
		"hsts_max_age":           63072000, // 2年
		"enable_csp":             true,
		"csp_content":            "default-src 'self'",
		"referrer_policy":        "no-referrer",
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.False(t, middleware.config.EnableXFrameOptions)
	assert.Equal(t, "SAMEORIGIN", middleware.config.XFrameOptions)
	assert.True(t, middleware.config.EnableHSTS)
	assert.Equal(t, 63072000, middleware.config.HSTSMaxAge)
	assert.True(t, middleware.config.EnableCSP)
	assert.Equal(t, "default-src 'self'", middleware.config.CSPContent)
	assert.Equal(t, "no-referrer", middleware.config.ReferrerPolicy)
}

// TestSecurityMiddleware_ValidateConfig 测试配置验证
func TestSecurityMiddleware_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *SecurityConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &SecurityConfig{
				EnableXFrameOptions: true,
				XFrameOptions:       "DENY",
				EnableHSTS:          true,
				HSTSMaxAge:          31536000,
			},
			wantErr: false,
		},
		{
			name: "无效的XFrameOptions",
			config: &SecurityConfig{
				EnableXFrameOptions: true,
				XFrameOptions:       "INVALID",
			},
			wantErr: true,
		},
		{
			name: "负数HSTSMaxAge",
			config: &SecurityConfig{
				EnableHSTS:   true,
				HSTSMaxAge:   -1,
			},
			wantErr: true,
		},
		{
			name: "启用CSP但无内容",
			config: &SecurityConfig{
				EnableCSP:    true,
				CSPContent:   "",
			},
			wantErr: true,
		},
		{
			name: "无效的ReferrerPolicy",
			config: &SecurityConfig{
				EnableReferrerPolicy: true,
				ReferrerPolicy:       "invalid-policy",
			},
			wantErr: true,
		},
		{
			name:   "默认配置",
			config: DefaultSecurityConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &SecurityMiddleware{config: tt.config}
			err := middleware.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSecurityMiddleware_AllSecurityHeaders 测试所有安全头是否设置
func TestSecurityMiddleware_AllSecurityHeaders(t *testing.T) {
	middleware := NewSecurityMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建带TLS的请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.TLS = &tls.ConnectionState{}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证所有默认启用的安全头
	assert.NotEmpty(t, w.Header().Get("X-Frame-Options"))
	assert.NotEmpty(t, w.Header().Get("X-Content-Type-Options"))
	assert.NotEmpty(t, w.Header().Get("X-XSS-Protection"))
	assert.NotEmpty(t, w.Header().Get("Strict-Transport-Security"))
	assert.NotEmpty(t, w.Header().Get("Referrer-Policy"))
	assert.NotEmpty(t, w.Header().Get("Permissions-Policy"))
}

// BenchmarkSecurityMiddleware 性能测试
func BenchmarkSecurityMiddleware(b *testing.B) {
	middleware := NewSecurityMiddleware()
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
