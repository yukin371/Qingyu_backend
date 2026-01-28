package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// queryTypes 查询类型分布计数器
	// Labels:
	//   - collection: 集合名称
	//   - operation: 操作类型 (find, findOne, update, delete, insert, aggregate等)
	//   - has_index: 是否使用了索引 (true/false)
	queryTypes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_query_types_total",
			Help: "Query types distribution with index usage information",
		},
		[]string{"collection", "operation", "has_index"},
	)

	// connectionPool 连接池状态指标
	// Labels:
	//   - state: 连接状态 (active, available, total)
	// Value: 当前连接数
	connectionPool = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_connection_pool",
			Help: "Connection pool statistics (active, available, total connections)",
		},
		[]string{"state"},
	)

	// indexEfficiency 索引效率指标
	// Labels:
	//   - collection: 集合名称
	//   - index_name: 索引名称
	// Value: 索引命中率 (0.0 - 1.0)
	indexEfficiency = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_index_efficiency",
			Help: "Index hit/miss ratio for specific indexes",
		},
		[]string{"collection", "index_name"},
	)
)

// RecordQuery 记录查询类型和索引使用情况
// 参数:
//   - collection: 集合名称
//   - operation: 操作类型 (find, findOne, update, delete, insert, aggregate等)
//   - hasIndex: 是否使用了索引
//
// 用法示例:
//
//	// 记录使用了索引的find查询
//	RecordQuery("users", "find", true)
//
//	// 记录没有使用索引的update查询
//	RecordQuery("orders", "update", false)
func RecordQuery(collection, operation string, hasIndex bool) {
	hasIndexStr := "false"
	if hasIndex {
		hasIndexStr = "true"
	}
	queryTypes.WithLabelValues(collection, operation, hasIndexStr).Inc()
}

// SetConnectionPool 设置连接池状态
// 参数:
//   - state: 连接状态，可选值: "active", "available", "total"
//   - count: 当前连接数
//
// 用法示例:
//
//	// 设置活跃连接数
//	SetConnectionPool("active", 10)
//
//	// 设置可用连接数
//	SetConnectionPool("available", 5)
//
//	// 设置总连接数
//	SetConnectionPool("total", 15)
func SetConnectionPool(state string, count float64) {
	connectionPool.WithLabelValues(state).Set(count)
}

// SetIndexEfficiency 设置索引效率
// 参数:
//   - collection: 集合名称
//   - indexName: 索引名称
//   - ratio: 索引命中率 (0.0 - 1.0)
//
// 用法示例:
//
//	// 设置users集合的email索引效率为95%
//	SetIndexEfficiency("users", "idx_email", 0.95)
//
//	// 设置books集合的title索引效率为88%
//	SetIndexEfficiency("books", "idx_title", 0.88)
func SetIndexEfficiency(collection, indexName string, ratio float64) {
	indexEfficiency.WithLabelValues(collection, indexName).Set(ratio)
}
