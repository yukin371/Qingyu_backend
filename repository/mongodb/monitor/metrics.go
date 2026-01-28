package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// slowQueryCount 慢查询计数器
	// Label: collection - 集合名称
	slowQueryCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_slow_queries_total",
			Help: "Total number of slow queries executed",
		},
		[]string{"collection"},
	)

	// queryDuration 查询持续时间直方图
	// Label: collection - 集合名称, operation - 操作类型(find, findOne, update, delete等)
	// Buckets: 使用默认buckets [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]秒
	queryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mongodb_query_duration_seconds",
			Help:    "Query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"collection", "operation"},
	)

	// indexUsage 索引使用率指标
	// Label: collection - 集合名称
	// Value: 0.0-1.0 之间的浮点数，表示索引使用率
	indexUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_index_usage_ratio",
			Help: "Ratio of queries using indexes",
		},
		[]string{"collection"},
	)
)

// RecordSlowQuery 记录慢查询
// 参数:
//   - collection: 集合名称
//
// 用法示例:
//
//	RecordSlowQuery("users")
//	RecordSlowQuery("books")
func RecordSlowQuery(collection string) {
	slowQueryCount.WithLabelValues(collection).Inc()
}

// RecordQueryDuration 记录查询持续时间
// 参数:
//   - collection: 集合名称
//   - operation: 操作类型 (find, findOne, update, delete, aggregate等)
//   - duration: 查询持续时间（秒）
//
// 用法示例:
//
//	start := time.Now()
//	err := collection.Find(ctx, filter)
//	duration := time.Since(start).Seconds()
//	RecordQueryDuration("books", "find", duration)
func RecordQueryDuration(collection, operation string, duration float64) {
	queryDuration.WithLabelValues(collection, operation).Observe(duration)
}

// SetIndexUsage 设置索引使用率
// 参数:
//   - collection: 集合名称
//   - ratio: 索引使用率 (0.0 - 1.0)
//
// 用法示例:
//
//	// 计算80%的查询使用了索引
//	SetIndexUsage("users", 0.8)
func SetIndexUsage(collection string, ratio float64) {
	indexUsage.WithLabelValues(collection).Set(ratio)
}
