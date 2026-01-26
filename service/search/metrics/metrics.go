package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 搜索指标收集器
type Metrics struct {
	// SearchRequestsTotal 搜索请求总数（按类型、状态、引擎分类）
	SearchRequestsTotal *prometheus.CounterVec

	// SearchDuration 搜索延迟直方图（按类型、引擎分类）
	SearchDuration *prometheus.HistogramVec

	// CacheHitsTotal 缓存命中总数
	CacheHitsTotal *prometheus.CounterVec

	// CacheMissesTotal 缓存未命中总数
	CacheMissesTotal *prometheus.CounterVec

	// ElasticsearchQueriesTotal ES 查询总数
	ElasticsearchQueriesTotal *prometheus.CounterVec

	// ElasticsearchQueryDuration ES 查询延迟
	ElasticsearchQueryDuration *prometheus.HistogramVec

	// RateLimitHitsTotal 限流触发次数
	RateLimitHitsTotal *prometheus.CounterVec

	// ActiveSearchesGauge 当前活跃搜索数
	ActiveSearchesGauge prometheus.Gauge

	// SearchResultSize 搜索结果大小
	SearchResultSize *prometheus.HistogramVec

	// DatabaseQueryDuration 数据库查询延迟
	DatabaseQueryDuration *prometheus.HistogramVec

	// SyncEventsTotal 同步事件总数
	SyncEventsTotal *prometheus.CounterVec

	// SyncQueueLength 同步队列长度
	SyncQueueLength prometheus.Gauge

	// DeadLetterQueueLength 死信队列长度
	DeadLetterQueueLength prometheus.Gauge
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		// 搜索请求总数
		SearchRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_requests_total",
				Help: "搜索请求总数",
			},
			[]string{"search_type", "status", "engine"},
		),

		// 搜索延迟直方图
		SearchDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "search_duration_seconds",
				Help:    "搜索延迟分布",
				Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"search_type", "engine"},
		),

		// 缓存命中总数
		CacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_cache_hits_total",
				Help: "缓存命中总数",
			},
			[]string{"search_type"},
		),

		// 缓存未命中总数
		CacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_cache_misses_total",
				Help: "缓存未命中总数",
			},
			[]string{"search_type"},
		),

		// ES 查询总数
		ElasticsearchQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_elasticsearch_queries_total",
				Help: "Elasticsearch 查询总数",
			},
			[]string{"index", "status"},
		),

		// ES 查询延迟
		ElasticsearchQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "search_elasticsearch_query_duration_seconds",
				Help:    "Elasticsearch 查询延迟分布",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
			},
			[]string{"index", "query_type"},
		),

		// 限流触发次数
		RateLimitHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_rate_limit_hits_total",
				Help: "限流触发次数",
			},
			[]string{"type"}, // user or ip
		),

		// 活跃搜索数
		ActiveSearchesGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "search_active_gauge",
				Help: "当前活跃搜索数",
			},
		),

		// 搜索结果大小
		SearchResultSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "search_result_size_bytes",
				Help:    "搜索结果大小分布",
				Buckets: []float64{1024, 10240, 102400, 1024000, 10485760},
			},
			[]string{"search_type"},
		),

		// 数据库查询延迟
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "search_database_query_duration_seconds",
				Help:    "数据库查询延迟分布",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5},
			},
			[]string{"collection"},
		),

		// 同步事件总数
		SyncEventsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_sync_events_total",
				Help: "同步事件总数",
			},
			[]string{"event_type", "status"},
		),

		// 同步队列长度
		SyncQueueLength: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "search_sync_queue_length",
				Help: "同步队列长度",
			},
		),

		// 死信队列长度
		DeadLetterQueueLength: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "search_dead_letter_queue_length",
				Help: "死信队列长度",
			},
		),
	}
}

// RecordSearch 记录搜索指标
func (m *Metrics) RecordSearch(searchType, status, engine string, duration time.Duration, resultSize int) {
	m.SearchRequestsTotal.WithLabelValues(searchType, status, engine).Inc()
	m.SearchDuration.WithLabelValues(searchType, engine).Observe(duration.Seconds())
	m.SearchResultSize.WithLabelValues(searchType).Observe(float64(resultSize))
}

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit(searchType string) {
	m.CacheHitsTotal.WithLabelValues(searchType).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss(searchType string) {
	m.CacheMissesTotal.WithLabelValues(searchType).Inc()
}

// RecordElasticsearchQuery 记录 ES 查询
func (m *Metrics) RecordElasticsearchQuery(index, status string, duration time.Duration) {
	m.ElasticsearchQueriesTotal.WithLabelValues(index, status).Inc()
	m.ElasticsearchQueryDuration.WithLabelValues(index, "search").Observe(duration.Seconds())
}

// RecordRateLimitHit 记录限流触发
func (m *Metrics) RecordRateLimitHit(limitType string) {
	m.RateLimitHitsTotal.WithLabelValues(limitType).Inc()
}

// IncActiveSearches 增加活跃搜索数
func (m *Metrics) IncActiveSearches() {
	m.ActiveSearchesGauge.Inc()
}

// DecActiveSearches 减少活跃搜索数
func (m *Metrics) DecActiveSearches() {
	m.ActiveSearchesGauge.Dec()
}

// RecordDatabaseQuery 记录数据库查询
func (m *Metrics) RecordDatabaseQuery(collection string, duration time.Duration) {
	m.DatabaseQueryDuration.WithLabelValues(collection).Observe(duration.Seconds())
}

// RecordSyncEvent 记录同步事件
func (m *Metrics) RecordSyncEvent(eventType, status string) {
	m.SyncEventsTotal.WithLabelValues(eventType, status).Inc()
}

// SetSyncQueueLength 设置同步队列长度
func (m *Metrics) SetSyncQueueLength(length float64) {
	m.SyncQueueLength.Set(length)
}

// SetDeadLetterQueueLength 设置死信队列长度
func (m *Metrics) SetDeadLetterQueueLength(length float64) {
	m.DeadLetterQueueLength.Set(length)
}

// GetCacheHitRate 获取缓存命中率
// 注意：此方法需要通过 Prometheus HTTP API 或 metric 接口获取实际值
// 这里返回占位符，实际使用时需要通过 /metrics 端点或 Prometheus 查询
func (m *Metrics) GetCacheHitRate(searchType string) (float64, error) {
	// TODO: 实现从 Prometheus metric 获取缓存命中率
	// 可以使用 prometheus.MetricVec 接口或 HTTP API
	return 0, nil
}

// 全局指标实例
var globalMetrics *Metrics

// InitMetrics 初始化全局指标
func InitMetrics() {
	if globalMetrics == nil {
		globalMetrics = NewMetrics()
	}
}

// GetMetrics 获取全局指标
func GetMetrics() *Metrics {
	if globalMetrics == nil {
		InitMetrics()
	}
	return globalMetrics
}
