package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// MetricsRegistry 指标注册表
	MetricsRegistry = prometheus.NewRegistry()

	once sync.Once
)

// Metrics 指标收集器
type Metrics struct {
	// HTTP请求指标
	httpRequestTotal    *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSize     *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec

	// 数据库连接池指标
	dbPoolConnections      *prometheus.GaugeVec
	dbPoolIdleConnections  *prometheus.GaugeVec
	dbPoolWaitDuration     *prometheus.HistogramVec
	dbQueryDuration        *prometheus.HistogramVec
	dbQueryTotal           *prometheus.CounterVec
	dbQueryErrors          *prometheus.CounterVec

	// 缓存指标
	cacheHits      *prometheus.CounterVec
	cacheMisses    *prometheus.CounterVec
	cacheDuration  *prometheus.HistogramVec

	// 业务指标
	userRegistrations    *prometheus.CounterVec
	userLogins           *prometheus.CounterVec
	bookPublish          *prometheus.CounterVec
	bookUpdate           *prometheus.CounterVec
	chapterPublish       *prometheus.CounterVec
	contentRead          *prometheus.CounterVec
	commentCreate        *prometheus.CounterVec
	paymentSuccess       *prometheus.CounterVec
	paymentFailed        *prometheus.CounterVec
	vipSubscriptions     *prometheus.CounterVec

	// 系统资源指标
	systemMemoryUsage *prometheus.GaugeVec
	systemCPUUsage    *prometheus.GaugeVec
	goroutineCount    *prometheus.GaugeVec
}

var (
	// 全局指标实例
	globalMetrics *Metrics
)

// GetMetrics 获取全局指标实例（单例）
func GetMetrics() *Metrics {
	once.Do(func() {
		globalMetrics = NewMetrics()
		MetricsRegistry.MustRegister(globalMetrics.getAllCollectors()...)
	})
	return globalMetrics
}

// NewMetrics 创建新的指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP请求指标
		httpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latencies in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		httpRequestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request sizes in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 7),
			},
			[]string{"method", "path"},
		),
		httpResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response sizes in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 7),
			},
			[]string{"method", "path"},
		),

		// 数据库连接池指标
		dbPoolConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_pool_connections",
				Help: "Number of database connections in the pool",
			},
			[]string{"database"},
		),
		dbPoolIdleConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_pool_idle_connections",
				Help: "Number of idle database connections",
			},
			[]string{"database"},
		),
		dbPoolWaitDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_pool_wait_duration_seconds",
				Help:    "Time waiting for database connection",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"database"},
		),
		dbQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"database", "operation"},
		),
		dbQueryTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"database", "operation", "status"},
		),
		dbQueryErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_query_errors_total",
				Help: "Total number of database query errors",
			},
			[]string{"database", "operation"},
		),

		// 缓存指标
		cacheHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache"},
		),
		cacheMisses: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache"},
		),
		cacheDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cache_duration_seconds",
				Help:    "Cache operation duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"cache", "operation"},
		),

		// 业务指标
		userRegistrations: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_registrations_total",
				Help: "Total number of user registrations",
			},
			[]string{"source"},
		),
		userLogins: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_logins_total",
				Help: "Total number of user logins",
			},
			[]string{"method"},
		),
		bookPublish: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "book_publish_total",
				Help: "Total number of published books",
			},
			[]string{"category"},
		),
		bookUpdate: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "book_update_total",
				Help: "Total number of book updates",
			},
			[]string{"category"},
		),
		chapterPublish: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "chapter_publish_total",
				Help: "Total number of published chapters",
			},
			[]string{"status"},
		),
		contentRead: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "content_read_total",
				Help: "Total number of content reads",
			},
			[]string{"type"},
		),
		commentCreate: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "comment_create_total",
				Help: "Total number of created comments",
			},
			[]string{"type"},
		),
		paymentSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_success_total",
				Help: "Total number of successful payments",
			},
			[]string{"type"},
		),
		paymentFailed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_failed_total",
				Help: "Total number of failed payments",
			},
			[]string{"type", "reason"},
		),
		vipSubscriptions: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "vip_subscriptions_total",
				Help: "Total number of VIP subscriptions",
			},
			[]string{"plan", "duration"},
		),

		// 系统资源指标
		systemMemoryUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_memory_usage_bytes",
				Help: "System memory usage in bytes",
			},
			[]string{"type"},
		),
		systemCPUUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_cpu_usage_percent",
				Help: "System CPU usage percentage",
			},
			[]string{},
		),
		goroutineCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goroutine_count",
				Help: "Number of goroutines",
			},
			[]string{},
		),
	}
}

// getAllCollectors 获取所有收集器
func (m *Metrics) getAllCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.httpRequestTotal,
		m.httpRequestDuration,
		m.httpRequestSize,
		m.httpResponseSize,
		m.dbPoolConnections,
		m.dbPoolIdleConnections,
		m.dbPoolWaitDuration,
		m.dbQueryDuration,
		m.dbQueryTotal,
		m.dbQueryErrors,
		m.cacheHits,
		m.cacheMisses,
		m.cacheDuration,
		m.userRegistrations,
		m.userLogins,
		m.bookPublish,
		m.bookUpdate,
		m.chapterPublish,
		m.contentRead,
		m.commentCreate,
		m.paymentSuccess,
		m.paymentFailed,
		m.vipSubscriptions,
		m.systemMemoryUsage,
		m.systemCPUUsage,
		m.goroutineCount,
	}
}

// ==================== HTTP请求指标 ====================

// RecordHttpRequest 记录HTTP请求
func (m *Metrics) RecordHttpRequest(method, path string, status int, duration time.Duration, requestSize, responseSize int64) {
	m.httpRequestTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	if requestSize > 0 {
		m.httpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	}
	if responseSize > 0 {
		m.httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
	}
}

// ==================== 数据库连接池指标 ====================

// UpdateDbPoolMetrics 更新数据库连接池指标
func (m *Metrics) UpdateDbPoolMetrics(database string, total, idle int, waitDuration time.Duration) {
	m.dbPoolConnections.WithLabelValues(database).Set(float64(total))
	m.dbPoolIdleConnections.WithLabelValues(database).Set(float64(idle))
	if waitDuration > 0 {
		m.dbPoolWaitDuration.WithLabelValues(database).Observe(waitDuration.Seconds())
	}
}

// RecordDbQuery 记录数据库查询
func (m *Metrics) RecordDbQuery(database, operation string, duration time.Duration, success bool) {
	m.dbQueryDuration.WithLabelValues(database, operation).Observe(duration.Seconds())
	status := "success"
	if !success {
		status = "error"
		m.dbQueryErrors.WithLabelValues(database, operation).Inc()
	}
	m.dbQueryTotal.WithLabelValues(database, operation, status).Inc()
}

// ==================== 缓存指标 ====================

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit(cache string) {
	m.cacheHits.WithLabelValues(cache).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss(cache string) {
	m.cacheMisses.WithLabelValues(cache).Inc()
}

// RecordCacheOperation 记录缓存操作
func (m *Metrics) RecordCacheOperation(cache, operation string, duration time.Duration) {
	m.cacheDuration.WithLabelValues(cache, operation).Observe(duration.Seconds())
}

// ==================== 业务指标 ====================

// RecordUserRegistration 记录用户注册
func (m *Metrics) RecordUserRegistration(source string) {
	m.userRegistrations.WithLabelValues(source).Inc()
}

// RecordUserLogin 记录用户登录
func (m *Metrics) RecordUserLogin(method string) {
	m.userLogins.WithLabelValues(method).Inc()
}

// RecordBookPublish 记录书籍发布
func (m *Metrics) RecordBookPublish(category string) {
	m.bookPublish.WithLabelValues(category).Inc()
}

// RecordBookUpdate 记录书籍更新
func (m *Metrics) RecordBookUpdate(category string) {
	m.bookUpdate.WithLabelValues(category).Inc()
}

// RecordChapterPublish 记录章节发布
func (m *Metrics) RecordChapterPublish(status string) {
	m.chapterPublish.WithLabelValues(status).Inc()
}

// RecordContentRead 记录内容阅读
func (m *Metrics) RecordContentRead(contentType string) {
	m.contentRead.WithLabelValues(contentType).Inc()
}

// RecordCommentCreate 记录评论创建
func (m *Metrics) RecordCommentCreate(commentType string) {
	m.commentCreate.WithLabelValues(commentType).Inc()
}

// RecordPaymentSuccess 记录支付成功
func (m *Metrics) RecordPaymentSuccess(paymentType string) {
	m.paymentSuccess.WithLabelValues(paymentType).Inc()
}

// RecordPaymentFailed 记录支付失败
func (m *Metrics) RecordPaymentFailed(paymentType, reason string) {
	m.paymentFailed.WithLabelValues(paymentType, reason).Inc()
}

// RecordVipSubscription 记录VIP订阅
func (m *Metrics) RecordVipSubscription(plan, duration string) {
	m.vipSubscriptions.WithLabelValues(plan, duration).Inc()
}

// ==================== 系统资源指标 ====================

// UpdateSystemMemory 更新系统内存使用
func (m *Metrics) UpdateSystemMemory(used, total uint64) {
	m.systemMemoryUsage.WithLabelValues("used").Set(float64(used))
	m.systemMemoryUsage.WithLabelValues("total").Set(float64(total))
}

// UpdateSystemCPU 更新系统CPU使用
func (m *Metrics) UpdateSystemCPU(percent float64) {
	m.systemCPUUsage.WithLabelValues().Set(percent)
}

// UpdateGoroutineCount 更新goroutine数量
func (m *Metrics) UpdateGoroutineCount(count int) {
	m.goroutineCount.WithLabelValues().Set(float64(count))
}

// ==================== HTTP处理器 ====================

// MetricsHandler Prometheus指标暴露端点
func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(MetricsRegistry, promhttp.HandlerOpts{})
}

// GinMetricsHandler Gin框架的指标处理器
func GinMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		MetricsHandler().ServeHTTP(c.Writer, c.Request)
	}
}
