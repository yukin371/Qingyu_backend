package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMetricsInitialization 测试指标初始化
func TestMetricsInitialization(t *testing.T) {
	metrics := GetMetrics()

	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.httpRequestTotal)
	assert.NotNil(t, metrics.httpRequestDuration)
	assert.NotNil(t, metrics.dbPoolConnections)
	assert.NotNil(t, metrics.userRegistrations)
}

// TestRecordHttpRequest 测试HTTP请求记录
func TestRecordHttpRequest(t *testing.T) {
	metrics := GetMetrics()

	metrics.RecordHttpRequest(
		"GET",
		"/api/v1/books",
		200,
		100*time.Millisecond,
		1024,
		2048,
	)

	// 验证指标已记录（通过检查注册表）
	_, err := MetricsRegistry.Gather()
	assert.NoError(t, err)
}

// TestMiddleware 测试HTTP中间件
func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Middleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "test response")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "test response")
}

// TestDbMetrics 测试数据库指标
func TestDbMetrics(t *testing.T) {
	metrics := GetMetrics()

	// 测试连接池指标
	metrics.UpdateDbPoolMetrics("qingyu", 10, 5, 50*time.Millisecond)

	// 测试查询指标
	metrics.RecordDbQuery("qingyu", "find", 100*time.Millisecond, true)
	metrics.RecordDbQuery("qingyu", "insert", 50*time.Millisecond, false)

	// 验证指标不抛出异常
	assert.NotNil(t, metrics.dbPoolConnections)
	assert.NotNil(t, metrics.dbQueryTotal)
}

// TestCacheMetrics 测试缓存指标
func TestCacheMetrics(t *testing.T) {
	metrics := GetMetrics()

	// 测试缓存命中
	metrics.RecordCacheHit("redis")
	metrics.RecordCacheHit("redis")

	// 测试缓存未命中
	metrics.RecordCacheMiss("redis")

	// 测试缓存操作
	metrics.RecordCacheOperation("redis", "get", 10*time.Millisecond)
	metrics.RecordCacheOperation("redis", "set", 20*time.Millisecond)

	assert.NotNil(t, metrics.cacheHits)
	assert.NotNil(t, metrics.cacheMisses)
}

// TestBusinessMetrics 测试业务指标
func TestBusinessMetrics(t *testing.T) {
	metrics := GetMetrics()

	// 测试用户注册
	metrics.RecordUserRegistration("web")
	metrics.RecordUserRegistration("mobile")

	// 测试用户登录
	metrics.RecordUserLogin("password")
	metrics.RecordUserLogin("oauth")

	// 测试书籍发布
	metrics.RecordBookPublish("fantasy")
	metrics.RecordBookPublish("scifi")

	// 测试章节发布
	metrics.RecordChapterPublish("draft")
	metrics.RecordChapterPublish("published")

	// 测试内容阅读
	metrics.RecordContentRead("book")
	metrics.RecordContentRead("chapter")

	// 测试评论创建
	metrics.RecordCommentCreate("post")
	metrics.RecordCommentCreate("reply")

	// 测试支付
	metrics.RecordPaymentSuccess("vip")
	metrics.RecordPaymentFailed("vip", "insufficient_balance")

	// 测试VIP订阅
	metrics.RecordVipSubscription("monthly", "30")
	metrics.RecordVipSubscription("yearly", "365")

	// 验证指标不抛出异常
	assert.NotNil(t, metrics.userRegistrations)
	assert.NotNil(t, metrics.paymentSuccess)
}

// TestSystemMetrics 测试系统资源指标
func TestSystemMetrics(t *testing.T) {
	metrics := GetMetrics()

	// 测试内存指标
	metrics.UpdateSystemMemory(1024*1024*100, 1024*1024*1024) // 100MB used, 1GB total

	// 测试CPU指标
	metrics.UpdateSystemCPU(45.5)

	// 测试goroutine计数
	metrics.UpdateGoroutineCount(150)

	// 验证指标不抛出异常
	assert.NotNil(t, metrics.systemMemoryUsage)
	assert.NotNil(t, metrics.systemCPUUsage)
	assert.NotNil(t, metrics.goroutineCount)
}

// TestMetricsHandler 测试Prometheus端点
func TestMetricsHandler(t *testing.T) {
	handler := MetricsHandler()

	assert.NotNil(t, handler)

	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

// TestGinMetricsHandler 测试Gin指标处理器
func TestGinMetricsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/metrics", GinMetricsHandler())

	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

// TestMetricsSingleton 测试单例模式
func TestMetricsSingleton(t *testing.T) {
	metrics1 := GetMetrics()
	metrics2 := GetMetrics()

	// 验证返回的是同一个实例
	assert.Same(t, metrics1, metrics2)
}

// BenchmarkRecordHttpRequest 性能测试 - HTTP请求记录
func BenchmarkRecordHttpRequest(b *testing.B) {
	metrics := GetMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordHttpRequest(
			"GET",
			"/api/v1/books",
			200,
			100*time.Millisecond,
			1024,
			2048,
		)
	}
}

// BenchmarkRecordDbQuery 性能测试 - 数据库查询记录
func BenchmarkRecordDbQuery(b *testing.B) {
	metrics := GetMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordDbQuery("qingyu", "find", 50*time.Millisecond, true)
	}
}

// BenchmarkRecordCacheHit 性能测试 - 缓存命中记录
func BenchmarkRecordCacheHit(b *testing.B) {
	metrics := GetMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordCacheHit("redis")
	}
}
