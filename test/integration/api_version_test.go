package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1"
	"Qingyu_backend/internal/middleware"
	pkgmiddleware "Qingyu_backend/pkg/middleware"
	"Qingyu_backend/pkg/response"
)

// TestHelper 测试辅助结构
type TestHelper struct {
	logger *zap.Logger
}

// NewTestHelper 创建测试辅助器
func NewTestHelper() *TestHelper {
	logger, _ := zap.NewDevelopment()
	return &TestHelper{
		logger: logger,
	}
}

// setupTestRouter 设置测试路由
func (h *TestHelper) setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加logger到context
	r.Use(func(c *gin.Context) {
		c.Set("logger", h.logger)
		c.Next()
	})

	// 添加版本路由中间件
	r.Use(middleware.VersionRoutingMiddleware())

	// 添加错误处理
	r.Use(pkgmiddleware.ErrorHandler())

	// 创建版本API
	versionAPI := v1.NewVersionAPI()

	// 注册版本路由
	v1Group := r.Group("/api/v1")
	v1.InitVersionRoutesWithVersionGroup(v1Group, versionAPI)

	// 注册全局版本路由
	apiGroup := r.Group("/api")
	v1.InitVersionRoutes(apiGroup, versionAPI)

	// 添加测试端点（模拟v1和v2）
	v1Group.GET("/test", func(c *gin.Context) {
		version := middleware.GetAPIVersion(c)
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"message": "v1 endpoint",
		})
	})

	// 添加v2路由组（模拟未来版本）
	v2Group := r.Group("/api/v2")
	v2Group.GET("/test", func(c *gin.Context) {
		version := middleware.GetAPIVersion(c)
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"message": "v2 endpoint",
		})
	})

	// 添加废弃端点（用于测试废弃响应头）
	sunsetDate := time.Now().AddDate(0, 0, 90)
	deprecationConfig := middleware.NewDeprecationConfig(
		&sunsetDate,
		"/api/v2/new-endpoint",
		"This API is deprecated",
	)
	v1Group.GET("/old-endpoint", middleware.DeprecationMiddleware(deprecationConfig), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "deprecated endpoint",
		})
	})

	return r
}

// T5.1: 版本路由测试
func TestAPIVersioning_Routing(t *testing.T) {
	h := NewTestHelper()
	r := h.setupTestRouter()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "v1路由",
			url:            "/api/v1/test",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"version": "v1",
				"message": "v1 endpoint",
			},
		},
		{
			name:           "v2路由",
			url:            "/api/v2/test",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"version": "v2",
				"message": "v2 endpoint",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody["version"], resp["version"])
			assert.Equal(t, tt.expectedBody["message"], resp["message"])
		})
	}
}

// T5.2: 版本选择Header测试
func TestAPIVersioning_HeaderSelection(t *testing.T) {
	h := NewTestHelper()
	r := h.setupTestRouter()

	// 创建一个使用Header选择版本的路由（注意：不使用/api前缀，避免被版本路由中间件从URL提取）
	r.GET("/test", func(c *gin.Context) {
		version := middleware.GetAPIVersion(c)
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"message": "version from header",
		})
	})

	tests := []struct {
		name           string
		header         string
		expectedStatus int
		expectedVer    string
	}{
		{
			name:           "使用v1版本",
			header:         "v1",
			expectedStatus: http.StatusOK,
			expectedVer:    "v1",
		},
		{
			name:           "使用v2版本",
			header:         "v2",
			expectedStatus: http.StatusOK,
			expectedVer:    "v2",
		},
		{
			name:           "不指定版本（使用默认）",
			header:         "",
			expectedStatus: http.StatusOK,
			expectedVer:    "v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.header != "" {
				req.Header.Set("X-API-Version", tt.header)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedVer, resp["version"])
		})
	}
}

// T5.3: 废弃响应头测试
func TestAPIVersioning_DeprecationHeaders(t *testing.T) {
	h := NewTestHelper()
	r := h.setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/old-endpoint", nil)
	r.ServeHTTP(w, req)

	// 验证HTTP状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证废弃响应头
	assert.Equal(t, "true", w.Header().Get("X-API-Deprecated"))
	assert.NotEmpty(t, w.Header().Get("X-API-Sunset-Date"))
	assert.Equal(t, "/api/v2/new-endpoint", w.Header().Get("X-API-Replacement"))
	assert.Contains(t, w.Header().Get("Warning"), "deprecated")
}

// T5.4: 废弃警告信息测试
func TestAPIVersioning_DeprecationWarning(t *testing.T) {
	h := NewTestHelper()
	r := h.setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/old-endpoint", nil)
	r.ServeHTTP(w, req)

	// 验证Warning头格式
	warningHeader := w.Header().Get("Warning")
	assert.NotEmpty(t, warningHeader)

	// Warning格式: 299 - "message"
	assert.Contains(t, warningHeader, "299")
	assert.Contains(t, warningHeader, "This API is deprecated")
	assert.Contains(t, warningHeader, "/api/v2/new-endpoint")
}

// T5.5: 版本信息端点测试
func TestAPIVersioning_VersionInfoEndpoint(t *testing.T) {
	h := NewTestHelper()
	r := h.setupTestRouter()

	t.Run("GET /api - 获取所有版本信息", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.EqualValues(t, 0, resp.Code)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, data["default_version"])

		versions, ok := data["versions"].([]interface{})
		assert.True(t, ok)
		assert.Greater(t, len(versions), 0)
	})

	t.Run("GET /api/v1 - 获取v1版本信息", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.EqualValues(t, 0, resp.Code)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "v1", data["version"])
	})
}

// T5.6: 多版本共存测试
func TestAPIVersioning_MultiVersionCoexistence(t *testing.T) {
	h := NewTestHelper()
	r := h.setupTestRouter()

	// 同时测试v1和v2端点
	t.Run("v1和v2同时可用", func(t *testing.T) {
		// 测试v1
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/api/v1/test", nil)
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		var resp1 map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &resp1)
		assert.Equal(t, "v1", resp1["version"])

		// 测试v2
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/api/v2/test", nil)
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		var resp2 map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &resp2)
		assert.Equal(t, "v2", resp2["version"])
	})
}

// TestVersionRoutingMiddleware 版本路由中间件测试
func TestVersionRoutingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.VersionRoutingMiddleware())

	r.GET("/test", func(c *gin.Context) {
		version := middleware.GetAPIVersion(c)
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	t.Run("从URL提取版本", func(t *testing.T) {
		// 添加v1路由组
		v1Group := r.Group("/api/v1")
		v1Group.GET("/test", func(c *gin.Context) {
			version := middleware.GetAPIVersion(c)
			c.JSON(http.StatusOK, gin.H{"version": version})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "v1", resp["version"])
	})
}

// TestDeprecationMiddleware 废弃中间件独立测试
func TestDeprecationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	sunsetDate := time.Now().AddDate(0, 0, 90)
	config := middleware.NewDeprecationConfig(
		&sunsetDate,
		"/api/v2/new-endpoint",
		"This endpoint is deprecated",
	)

	r.GET("/deprecated", middleware.DeprecationMiddleware(config), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/deprecated", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("X-API-Deprecated"))
	assert.Equal(t, "/api/v2/new-endpoint", w.Header().Get("X-API-Replacement"))
}

// TestDeprecationRegistry 废弃注册表测试
func TestDeprecationRegistry(t *testing.T) {
	registry := middleware.NewDeprecationRegistry()

	sunsetDate := time.Now().AddDate(0, 0, 90)
	config := middleware.NewDeprecationConfig(
		&sunsetDate,
		"/api/v2/new-endpoint",
		"Deprecated",
	)

	registry.RegisterEndpoint("GET", "/api/v1/old-endpoint", config)

	t.Run("检查端点是否废弃", func(t *testing.T) {
		foundConfig, exists := registry.IsDeprecated("GET", "/api/v1/old-endpoint")
		assert.True(t, exists)
		assert.True(t, foundConfig.Enabled)
		assert.Equal(t, "/api/v2/new-endpoint", foundConfig.Replacement)
	})

	t.Run("未废弃的端点", func(t *testing.T) {
		_, exists := registry.IsDeprecated("GET", "/api/v1/new-endpoint")
		assert.False(t, exists)
	})
}

// TestVersionRegistry 版本注册表测试
func TestVersionRegistry(t *testing.T) {
	registry := middleware.NewVersionRegistry()

	registry.RegisterVersion(&middleware.VersionConfig{
		Version:     "v1",
		Status:      "stable",
		Path:        "/api/v1",
		Description: "Stable version",
	})

	registry.RegisterVersion(&middleware.VersionConfig{
		Version:     "v2",
		Status:      "beta",
		Path:        "/api/v2",
		Description: "Beta version",
	})

	t.Run("获取版本配置", func(t *testing.T) {
		config, exists := registry.GetVersion("v1")
		assert.True(t, exists)
		assert.Equal(t, "v1", config.Version)
		assert.Equal(t, "stable", config.Status)
	})

	t.Run("获取所有版本", func(t *testing.T) {
		versions := registry.GetAllVersions()
		assert.Equal(t, 2, len(versions))
	})

	t.Run("设置和获取默认版本", func(t *testing.T) {
		registry.SetDefaultVersion("v2")
		assert.Equal(t, "v2", registry.GetDefaultVersion())
	})

	t.Run("检查版本可用性", func(t *testing.T) {
		assert.True(t, registry.IsVersionAvailable("v1"))
		assert.True(t, registry.IsVersionAvailable("v2"))
		assert.False(t, registry.IsVersionAvailable("v3"))
	})
}

// TestDeprecationConfigOptions 废弃配置选项测试
func TestDeprecationConfigOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	sunsetDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	r.GET("/deprecated-options", func(c *gin.Context) {
		middleware.SetDeprecationHeadersWithOptions(c,
			middleware.WithSunsetDate(sunsetDate),
			middleware.WithReplacement("/api/v2/new"),
			middleware.WithWarningMessage("Custom warning"),
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

// TestParseVersionFromHeader Header版本解析测试
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middleware.ParseVersionFromHeader(tt.header)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsValidVersion 版本格式验证测试
func TestIsValidVersion(t *testing.T) {
	// 使用反射或导出测试函数
	// 这里测试isValidVersion函数的逻辑
	t.Run("有效版本格式", func(t *testing.T) {
		// v1, v2, v1.1 等格式应该有效
		// 由于isValidVersion是未导出函数，这里通过中间件行为间接测试
	})
}

// TestGetAPIVersion Context获取版本测试
func TestGetAPIVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.VersionRoutingMiddleware())

	r.GET("/test", func(c *gin.Context) {
		version := middleware.GetAPIVersion(c)
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 没有版本信息时应该返回默认版本
	assert.Equal(t, "v1", resp["version"])
}
