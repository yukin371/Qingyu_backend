package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestNewDeprecationConfig 测试创建废弃配置
func TestNewDeprecationConfig(t *testing.T) {
	sunsetDate := time.Now().AddDate(0, 0, 90)
	config := NewDeprecationConfig(&sunsetDate, "/api/v2/new", "Deprecated")

	assert.True(t, config.Enabled)
	assert.Equal(t, "/api/v2/new", config.Replacement)
	assert.Equal(t, "Deprecated", config.Message)
	assert.NotNil(t, config.SunsetDate)
}

// TestDeprecationRegistry 测试废弃注册表
func TestDeprecationRegistry(t *testing.T) {
	registry := NewDeprecationRegistry()

	sunsetDate := time.Now().AddDate(0, 0, 90)
	config := NewDeprecationConfig(&sunsetDate, "/api/v2/new", "Deprecated")

	registry.RegisterEndpoint("GET", "/api/v1/old", config)

	t.Run("检查端点是否废弃", func(t *testing.T) {
		foundConfig, exists := registry.IsDeprecated("GET", "/api/v1/old")
		assert.True(t, exists)
		assert.True(t, foundConfig.Enabled)
		assert.Equal(t, "/api/v2/new", foundConfig.Replacement)
	})

	t.Run("未废弃的端点", func(t *testing.T) {
		_, exists := registry.IsDeprecated("GET", "/api/v1/new")
		assert.False(t, exists)
	})
}

// TestDeprecationMiddleware 测试废弃中间件
func TestDeprecationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	sunsetDate := time.Now().AddDate(0, 0, 90)
	config := NewDeprecationConfig(&sunsetDate, "/api/v2/new", "This endpoint is deprecated")

	r.GET("/deprecated", DeprecationMiddleware(config), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/deprecated", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("X-API-Deprecated"))
	assert.Equal(t, "/api/v2/new", w.Header().Get("X-API-Replacement"))
	assert.NotEmpty(t, w.Header().Get("X-API-Sunset-Date"))
	assert.Contains(t, w.Header().Get("Warning"), "deprecated")
}

// TestSetDeprecationHeadersWithOptions 测试选项模式设置废弃头
func TestSetDeprecationHeadersWithOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	sunsetDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	r.GET("/deprecated-options", func(c *gin.Context) {
		SetDeprecationHeadersWithOptions(c,
			WithSunsetDate(sunsetDate),
			WithReplacement("/api/v2/new"),
			WithWarningMessage("Custom warning"),
		)
		c.JSON(http.StatusOK, gin.H{})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/deprecated-options", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "true", w.Header().Get("X-API-Deprecated"))
	assert.Contains(t, w.Header().Get("X-API-Sunset-Date"), "2026-06-01")
	assert.Equal(t, "/api/v2/new", w.Header().Get("X-API-Replacement"))
	assert.Contains(t, w.Header().Get("Warning"), "Custom warning")
}

// TestVersionRegistryIsVersionAvailable 测试版本可用性检查
func TestVersionRegistryIsVersionAvailable(t *testing.T) {
	registry := NewVersionRegistry()

	registry.RegisterVersion(&VersionConfig{
		Version:     "v1",
		Status:      "stable",
		Path:        "/api/v1",
		Description: "Stable version",
	})

	registry.RegisterVersion(&VersionConfig{
		Version:     "v2",
		Status:      "beta",
		Path:        "/api/v2",
		Description: "Beta version",
	})

	t.Run("检查已注册版本", func(t *testing.T) {
		assert.True(t, registry.IsVersionAvailable("v1"))
		assert.True(t, registry.IsVersionAvailable("v2"))
	})

	t.Run("检查未注册版本", func(t *testing.T) {
		assert.False(t, registry.IsVersionAvailable("v3"))
		assert.False(t, registry.IsVersionAvailable(""))
	})
}

// TestParseVersionFromHeader 测试Header版本解析
func TestParseVersionFromHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"小写v1", "v1", "v1"},
		{"大写V1", "V1", "v1"},
		{"latest", "latest", "v1"},
		{"纯数字1", "1", "v1"},
		{"空字符串", "", ""},
		{"无效格式", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseVersionFromHeader(tt.header)
			assert.Equal(t, tt.expected, result)
		})
	}
}
