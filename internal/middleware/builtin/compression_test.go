package builtin

import (
	"compress/gzip"
	"fmt"
	"io"
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

// TestCompressionMiddleware_Name 测试中间件名称
func TestCompressionMiddleware_Name(t *testing.T) {
	middleware := NewCompressionMiddleware()
	assert.Equal(t, "compression", middleware.Name())
}

// TestCompressionMiddleware_Priority 测试中间件优先级
func TestCompressionMiddleware_Priority(t *testing.T) {
	middleware := NewCompressionMiddleware()
	assert.Equal(t, 12, middleware.Priority())
}

// TestCompressionMiddleware_CompressJSON 测试压缩JSON响应
func TestCompressionMiddleware_CompressJSON(t *testing.T) {
	middleware := NewCompressionMiddleware()
	middleware.config.MinLength = 100 // 设置较低的阈值

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		// 返回大量JSON数据
		largeData := make(map[string]string)
		for i := 0; i < 100; i++ {
			largeData[string(rune(i))] = strings.Repeat("data", 10)
		}
		c.JSON(http.StatusOK, largeData)
	})

	// 发送请求（支持gzip）
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
}

// TestCompressionMiddleware_NoCompressionWithoutGzipSupport 测试客户端不支持gzip时不压缩
func TestCompressionMiddleware_NoCompressionWithoutGzipSupport(t *testing.T) {
	middleware := NewCompressionMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求（不支持gzip）
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "identity")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Content-Encoding"), "不应该设置Content-Encoding头")
}

// TestCompressionMiddleware_Disabled 测试禁用压缩
func TestCompressionMiddleware_Disabled(t *testing.T) {
	middleware := NewCompressionMiddleware()
	middleware.config.Enabled = false

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		largeData := strings.Repeat("data", 1000)
		c.String(http.StatusOK, largeData)
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Content-Encoding"), "压缩禁用时不应该设置Content-Encoding头")
}

// TestCompressionMiddleware_SmallResponse 测试小响应不压缩
func TestCompressionMiddleware_SmallResponse(t *testing.T) {
	middleware := NewCompressionMiddleware()
	middleware.config.MinLength = 1024

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Content-Encoding"), "小响应不应该被压缩")
}

// TestCompressionMiddleware_ExcludedTypes 测试排除类型
func TestCompressionMiddleware_ExcludedTypes(t *testing.T) {
	middleware := NewCompressionMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/image", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/png", []byte(strings.Repeat("data", 2000)))
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/image", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Content-Encoding"), "图片不应该被压缩")
}

// TestCompressionMiddleware_CompressionLevels 测试不同压缩级别
func TestCompressionMiddleware_CompressionLevels(t *testing.T) {
	levels := []int{0, 1, 5, 9}

	for _, level := range levels {
		t.Run(fmt.Sprintf("Level_%d", level), func(t *testing.T) {
			middleware := NewCompressionMiddleware()
			middleware.config.Level = level
			middleware.config.MinLength = 100

			// 创建测试路由
			router := gin.New()
			router.Use(middleware.Handler())
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, strings.Repeat("data", 100))
			})

			// 发送请求
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestCompressionMiddleware_LoadConfig 测试配置加载
func TestCompressionMiddleware_LoadConfig(t *testing.T) {
	middleware := NewCompressionMiddleware()

	// 测试加载配置
	config := map[string]interface{}{
		"enabled":     false,
		"level":       7,
		"min_length":  512,
		"types":       []interface{}{"application/json", "text/html"},
		"excluded_types": []interface{}{"image/*", "video/*"},
	}

	err := middleware.LoadConfig(config)
	assert.NoError(t, err, "配置加载应该成功")

	// 验证配置已加载
	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, 7, middleware.config.Level)
	assert.Equal(t, 512, middleware.config.MinLength)
	assert.Equal(t, []string{"application/json", "text/html"}, middleware.config.Types)
	assert.Equal(t, []string{"image/*", "video/*"}, middleware.config.ExcludedTypes)
}

// TestCompressionMiddleware_ValidateConfig 测试配置验证
func TestCompressionMiddleware_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *CompressionConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &CompressionConfig{
				Enabled:   true,
				Level:     5,
				MinLength: 1024,
			},
			wantErr: false,
		},
		{
			name: "压缩级别过大",
			config: &CompressionConfig{
				Level: 10,
			},
			wantErr: true,
		},
		{
			name: "负数压缩级别",
			config: &CompressionConfig{
				Level: -1,
			},
			wantErr: true,
		},
		{
			name: "负数MinLength",
			config: &CompressionConfig{
				Level:     5,
				MinLength: -1,
			},
			wantErr: true,
		},
		{
			name:    "默认配置",
			config:  DefaultCompressionConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &CompressionMiddleware{config: tt.config}
			err := middleware.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCompressionMiddleware_DecompressResponse 测试解压缩响应
func TestCompressionMiddleware_DecompressResponse(t *testing.T) {
	middleware := NewCompressionMiddleware()
	middleware.config.MinLength = 100

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		largeData := strings.Repeat("hello world ", 100)
		c.String(http.StatusOK, largeData)
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

	// 解压缩响应
	reader, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Contains(t, string(decompressed), "hello world")
}

// BenchmarkCompressionMiddleware 性能测试
func BenchmarkCompressionMiddleware(b *testing.B) {
	middleware := NewCompressionMiddleware()

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
