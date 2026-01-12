package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics Prometheus指标集合
type PrometheusMetrics struct {
	// HTTP指标
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestSize     *prometheus.HistogramVec
	HTTPResponseSize    *prometheus.HistogramVec
	HTTPActiveRequests  prometheus.Gauge

	// 服务指标
	ServiceHealthStatus *prometheus.GaugeVec
	ServiceUptime       *prometheus.GaugeVec
	ServiceCallTotal    *prometheus.CounterVec
	ServiceCallDuration *prometheus.HistogramVec
	ServiceErrors       *prometheus.CounterVec

	// 业务指标 - AI配额
	AIQuotaUsage      *prometheus.GaugeVec
	AIQuotaTotal      *prometheus.GaugeVec
	AIQuotaRemaining  *prometheus.GaugeVec
	AIRequestsTotal   *prometheus.CounterVec
	AIRequestDuration *prometheus.HistogramVec
	AITokensConsumed  *prometheus.CounterVec

	// 业务指标 - 用户
	UserActiveTotal   prometheus.Gauge
	UserRegistrations prometheus.Counter
	UserLogins        prometheus.Counter

	// 业务指标 - 数据库
	DBConnectionsActive prometheus.Gauge
	DBConnectionsIdle   prometheus.Gauge
	DBQueryDuration     *prometheus.HistogramVec
	DBQueryTotal        *prometheus.CounterVec
	DBErrors            *prometheus.CounterVec

	// 业务指标 - Redis
	RedisHits            prometheus.Counter
	RedisMisses          prometheus.Counter
	RedisConnections     prometheus.Gauge
	RedisCommandDuration *prometheus.HistogramVec

	// 业务指标 - 书城
	BookViewsTotal     *prometheus.CounterVec
	BookPurchasesTotal *prometheus.CounterVec
	BookRatingsTotal   *prometheus.CounterVec
}

var (
	// DefaultMetrics 全局Prometheus指标实例
	DefaultMetrics *PrometheusMetrics
	metricsOnce    sync.Once
)

// InitPrometheusMetrics 初始化Prometheus指标（幂等，多次调用只注册一次）
func InitPrometheusMetrics() *PrometheusMetrics {
	// 如果已经初始化，直接返回
	if DefaultMetrics != nil {
		return DefaultMetrics
	}

	metricsOnce.Do(func() {
		metrics := &PrometheusMetrics{
			// ============ HTTP指标 ============
			HTTPRequestsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "HTTP请求总数",
				},
				[]string{"method", "path", "status"},
			),
			HTTPRequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_request_duration_seconds",
					Help:    "HTTP请求响应时间（秒）",
					Buckets: prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
				},
				[]string{"method", "path"},
			),
			HTTPRequestSize: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_request_size_bytes",
					Help:    "HTTP请求大小（字节）",
					Buckets: []float64{100, 1000, 10000, 100000, 1000000},
				},
				[]string{"method", "path"},
			),
			HTTPResponseSize: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_response_size_bytes",
					Help:    "HTTP响应大小（字节）",
					Buckets: []float64{100, 1000, 10000, 100000, 1000000},
				},
				[]string{"method", "path"},
			),
			HTTPActiveRequests: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "http_active_requests",
					Help: "当前活跃的HTTP请求数",
				},
			),

			// ============ 服务指标 ============
			ServiceHealthStatus: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "service_health_status",
					Help: "服务健康状态（1=健康，0=不健康）",
				},
				[]string{"service", "version"},
			),
			ServiceUptime: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "service_uptime_seconds",
					Help: "服务运行时间（秒）",
				},
				[]string{"service"},
			),
			ServiceCallTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "service_call_total",
					Help: "服务调用总数",
				},
				[]string{"service", "method", "status"},
			),
			ServiceCallDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "service_call_duration_seconds",
					Help:    "服务调用响应时间（秒）",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"service", "method"},
			),
			ServiceErrors: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "service_errors_total",
					Help: "服务错误总数",
				},
				[]string{"service", "error_type"},
			),

			// ============ AI配额指标 ============
			AIQuotaUsage: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "ai_quota_usage",
					Help: "AI配额使用量",
				},
				[]string{"user_id", "quota_type"},
			),
			AIQuotaTotal: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "ai_quota_total",
					Help: "AI配额总量",
				},
				[]string{"user_id", "quota_type"},
			),
			AIQuotaRemaining: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "ai_quota_remaining",
					Help: "AI配额剩余量",
				},
				[]string{"user_id", "quota_type"},
			),
			AIRequestsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "ai_requests_total",
					Help: "AI请求总数",
				},
				[]string{"service", "model", "status"},
			),
			AIRequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "ai_request_duration_seconds",
					Help:    "AI请求响应时间（秒）",
					Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
				},
				[]string{"service", "model"},
			),
			AITokensConsumed: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "ai_tokens_consumed_total",
					Help: "AI Token消耗总数",
				},
				[]string{"user_id", "model"},
			),

			// ============ 用户指标 ============
			UserActiveTotal: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "user_active_total",
					Help: "当前活跃用户数",
				},
			),
			UserRegistrations: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "user_registrations_total",
					Help: "用户注册总数",
				},
			),
			UserLogins: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "user_logins_total",
					Help: "用户登录总数",
				},
			),

			// ============ 数据库指标 ============
			DBConnectionsActive: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "db_connections_active",
					Help: "数据库活跃连接数",
				},
			),
			DBConnectionsIdle: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "db_connections_idle",
					Help: "数据库空闲连接数",
				},
			),
			DBQueryDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "db_query_duration_seconds",
					Help:    "数据库查询响应时间（秒）",
					Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
				},
				[]string{"operation", "collection"},
			),
			DBQueryTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "db_query_total",
					Help: "数据库查询总数",
				},
				[]string{"operation", "collection", "status"},
			),
			DBErrors: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "db_errors_total",
					Help: "数据库错误总数",
				},
				[]string{"operation", "error_type"},
			),

			// ============ Redis指标 ============
			RedisHits: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "redis_hits_total",
					Help: "Redis缓存命中总数",
				},
			),
			RedisMisses: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "redis_misses_total",
					Help: "Redis缓存未命中总数",
				},
			),
			RedisConnections: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "redis_connections",
					Help: "Redis连接数",
				},
			),
			RedisCommandDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "redis_command_duration_seconds",
					Help:    "Redis命令响应时间（秒）",
					Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1},
				},
				[]string{"command"},
			),

			// ============ 书城指标 ============
			BookViewsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "book_views_total",
					Help: "书籍查看总数",
				},
				[]string{"book_id"},
			),
			BookPurchasesTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "book_purchases_total",
					Help: "书籍购买总数",
				},
				[]string{"book_id"},
			),
			BookRatingsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "book_ratings_total",
					Help: "书籍评分总数",
				},
				[]string{"book_id", "rating"},
			),
		}

		DefaultMetrics = metrics
	})

	return DefaultMetrics
}

// GetDefaultMetrics 获取默认Prometheus指标实例
func GetDefaultMetrics() *PrometheusMetrics {
	if DefaultMetrics == nil {
		DefaultMetrics = InitPrometheusMetrics()
	}
	return DefaultMetrics
}
