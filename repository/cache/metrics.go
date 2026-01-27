package cache

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CacheHitsTotal 缓存命中总数
	CacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"prefix"},
	)

	// CacheMissesTotal 缓存未命中总数
	CacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"prefix"},
	)

	// CacheOperationDuration 缓存操作耗时
	CacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "Duration of cache operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"prefix", "operation"},
	)

	// CachePenetrationTotal 缓存穿透总数
	CachePenetrationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_penetration_total",
			Help: "Total number of cache penetration attacks prevented",
		},
		[]string{"prefix"},
	)

	// CacheBreakdownTotal 缓存击穿总数
	CacheBreakdownTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_breakdown_total",
			Help: "Total number of hot key breakdown events",
		},
		[]string{"prefix"},
	)
)

// RecordCacheHit 记录缓存命中
func RecordCacheHit(prefix string) {
	CacheHitsTotal.WithLabelValues(prefix).Inc()
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(prefix string) {
	CacheMissesTotal.WithLabelValues(prefix).Inc()
}

// RecordCacheOperation 记录缓存操作耗时
func RecordCacheOperation(prefix, operation string, duration float64) {
	CacheOperationDuration.WithLabelValues(prefix, operation).Observe(duration)
}

// RecordCachePenetration 记录缓存穿透
func RecordCachePenetration(prefix string) {
	CachePenetrationTotal.WithLabelValues(prefix).Inc()
}

// RecordCacheBreakdown 记录缓存击穿
func RecordCacheBreakdown(prefix string) {
	CacheBreakdownTotal.WithLabelValues(prefix).Inc()
}

// GetCacheHitRatioPromQL 返回计算缓存命中率的PromQL表达式
func GetCacheHitRatioPromQL(prefix string) string {
	return fmt.Sprintf(
		"rate(cache_hits_total{prefix=\"%s\"}[5m]) / (rate(cache_hits_total{prefix=\"%s\"}[5m]) + rate(cache_misses_total{prefix=\"%s\"}[5m]))",
		prefix, prefix, prefix,
	)
}
