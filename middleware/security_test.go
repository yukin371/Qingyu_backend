package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurity(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		config         SecurityConfig
		checkHeaders   map[string]string
		skipCSRF       bool
		method         string
		csrfToken      string
		expectStatus   int
	}{
		{
			name: "default security headers",
			config: DefaultSecurityConfig(),
			checkHeaders: map[string]string{
				"X-XSS-Protection":         "1; mode=block",
				"X-Content-Type-Options":   "nosniff",
				"X-Frame-Options":          "DENY",
				"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
				"Content-Security-Policy":  "default-src 'self'",
				"Referrer-Policy":          "strict-origin-when-cross-origin",
			},
			skipCSRF:     true,
			method:       "GET",
			expectStatus: 200,
		},
		{
			name: "custom security headers",
			config: SecurityConfig{
				XSSProtection:      true,
				ContentTypeNoSniff: true,
				FrameOptions:       "SAMEORIGIN",
				HSTSMaxAge:         31536000,
				CSRFProtection:     false,
			},
			checkHeaders: map[string]string{
				"X-XSS-Protection":       "1; mode=block",
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options":        "SAMEORIGIN",
			},
			skipCSRF:     true,
			method:       "GET",
			expectStatus: 200,
		},
		{
			name: "CSRF protection enabled - no token",
			config: SecurityConfig{
				CSRFProtection: true,
				CSRFTokenName:  "X-CSRF-Token",
				CSRFCookieName: "_csrf",
			},
			checkHeaders: map[string]string{},
			skipCSRF:     false,
			method:       "POST",
			csrfToken:    "",
			expectStatus: 403,
		},
		{
			name: "CSRF protection - skip GET requests",
			config: SecurityConfig{
				CSRFProtection: true,
				CSRFTokenName:  "X-CSRF-Token",
				CSRFCookieName: "_csrf",
			},
			checkHeaders: map[string]string{},
			skipCSRF:     true,
			method:       "GET",
			expectStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(SecurityWithConfig(tt.config))
			router.GET("/test", func(c *gin.Context) {
				c.Status(200)
			})
			router.POST("/test", func(c *gin.Context) {
				c.Status(200)
			})

			var req *http.Request
			if tt.method == "POST" {
				req = httptest.NewRequest("POST", "/test", nil)
				if tt.csrfToken != "" {
					req.Header.Set(tt.config.CSRFTokenName, tt.csrfToken)
				}
			} else {
				req = httptest.NewRequest("GET", "/test", nil)
			}

			// 设置CSRF cookie（如果需要）
			if tt.config.CSRFProtection && !tt.skipCSRF && tt.csrfToken != "" {
				req.AddCookie(&http.Cookie{
					Name:  tt.config.CSRFCookieName,
					Value: tt.csrfToken,
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			// 检查安全头
			for header, expectedValue := range tt.checkHeaders {
				actualValue := w.Header().Get(header)
				if actualValue != expectedValue {
					t.Errorf("Header %s: expected %q, got %q", header, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestGenerateCSRFToken(t *testing.T) {
	token1 := GenerateCSRFToken()
	token2 := GenerateCSRFToken()

	// Token不应该为空
	if token1 == "" || token2 == "" {
		t.Error("Generated token should not be empty")
	}

	// 每次生成的token应该不同
	if token1 == token2 {
		t.Error("Tokens should be unique")
	}

	// Token长度应该是64（32字节的hex编码）
	if len(token1) != 64 {
		t.Errorf("Expected token length 64, got %d", len(token1))
	}
}

func TestSetCSRFToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := DefaultSecurityConfig()
	router := gin.New()
	router.GET("/csrf", func(c *gin.Context) {
		SetCSRFToken(c, config)
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/csrf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应头中的token
	token := w.Header().Get(config.CSRFTokenName)
	if token == "" {
		t.Error("CSRF token should be set in response header")
	}

	// 检查cookie
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("CSRF cookie should be set")
		return
	}

	var csrfCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == config.CSRFCookieName {
			csrfCookie = cookie
			break
		}
	}

	if csrfCookie == nil {
		t.Error("CSRF cookie not found")
		return
	}

	if csrfCookie.Value != token {
		t.Errorf("Cookie token %q should match header token %q", csrfCookie.Value, token)
	}

	if !csrfCookie.HttpOnly {
		t.Error("CSRF cookie should be HttpOnly")
	}
}

func TestCSRFTokenMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CSRFToken())
	router.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查token是否在响应头中
	token := w.Header().Get("X-CSRF-Token")
	if token == "" {
		t.Error("CSRF token should be set in response header")
	}

	// 检查cookie
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("CSRF cookie should be set")
		return
	}

	var hasCSRFCookie bool
	for _, cookie := range cookies {
		if cookie.Name == "_csrf" {
			hasCSRFCookie = true
			break
		}
	}

	if !hasCSRFCookie {
		t.Error("CSRF cookie not found")
	}
}

func TestShouldSkipCSRF(t *testing.T) {
	tests := []struct {
		path     string
		skipPaths []string
		expected bool
	}{
		{
			path:      "/api/auth/login",
			skipPaths: []string{"/health", "/metrics", "/api/auth/login"},
			expected:  true,
		},
		{
			path:      "/api/auth/register",
			skipPaths: []string{"/health", "/metrics", "/api/auth/login"},
			expected:  false,
		},
		{
			path:      "/health",
			skipPaths: []string{"/health", "/metrics"},
			expected:  true,
		},
		{
			path:      "/api/users",
			skipPaths: []string{},
			expected:  false,
		},
		{
			path:      "/api/auth/login/callback",
			skipPaths: []string{"/api/auth/login"},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldSkipCSRF(tt.path, tt.skipPaths)
			if result != tt.expected {
				t.Errorf("shouldSkipCSRF(%q, %v) = %v, want %v", tt.path, tt.skipPaths, result, tt.expected)
			}
		})
	}
}
