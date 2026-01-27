package monitor

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMetricsRegistration 测试指标是否成功注册
func TestMetricsRegistration(t *testing.T) {
	// 先调用一些函数确保包被初始化
	RecordSlowQuery("init_test")
	RecordQueryDuration("init_test", "find", 0.1)
	SetIndexUsage("init_test", 0.5)
	RecordQuery("init_test", "find", true)
	SetConnectionPool("active", 1.0)
	SetIndexEfficiency("init_test", "idx_test", 0.9)

	// 验证基础指标已注册（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err, "应该能够成功收集指标")

	var metricNames []string
	for _, m := range metrics {
		metricNames = append(metricNames, m.GetName())
	}

	// 验证基础指标存在
	assert.Contains(t, metricNames, "mongodb_slow_queries_total", "慢查询计数器应该已注册")
	assert.Contains(t, metricNames, "mongodb_query_duration_seconds", "查询延迟直方图应该已注册")
	assert.Contains(t, metricNames, "mongodb_index_usage_ratio", "索引使用率指标应该已注册")

	// 验证增强指标存在
	assert.Contains(t, metricNames, "mongodb_query_types_total", "查询类型分布指标应该已注册")
	assert.Contains(t, metricNames, "mongodb_connection_pool", "连接池状态指标应该已注册")
	assert.Contains(t, metricNames, "mongodb_index_efficiency", "索引效率指标应该已注册")
}

// TestRecordSlowQuery 测试慢查询计数功能
func TestRecordSlowQuery(t *testing.T) {
	collection := "test_collection"

	// 记录一次慢查询
	RecordSlowQuery(collection)

	// 验证指标值（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找慢查询指标
	var slowQueryMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_slow_queries_total" {
			slowQueryMetric = m
			break
		}
	}

	require.NotNil(t, slowQueryMetric, "慢查询指标应该存在")

	// 查找对应collection的metric
	var targetMetric *dto.Metric
	for _, m := range slowQueryMetric.Metric {
		for _, label := range m.Label {
			if label.GetName() == "collection" && label.GetValue() == collection {
				targetMetric = m
				break
			}
		}
	}

	require.NotNil(t, targetMetric, "应该找到test_collection的metric")

	// 验证label
	require.Len(t, targetMetric.Label, 1, "应该有一个label")
	assert.Equal(t, "collection", targetMetric.Label[0].GetName(), "label名称应该是collection")
	assert.Equal(t, collection, targetMetric.Label[0].GetValue(), "label值应该是test_collection")

	// 验证计数器值
	assert.NotNil(t, targetMetric.Counter, "应该是Counter类型")
	assert.Equal(t, float64(1), targetMetric.Counter.GetValue(), "计数器值应该是1")
}

// TestRecordQueryDuration 测试查询延迟记录功能
func TestRecordQueryDuration(t *testing.T) {
	collection := "books"
	operation := "find"
	duration := 0.5 // 500ms

	// 记录查询延迟
	RecordQueryDuration(collection, operation, duration)

	// 验证指标值（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找查询延迟指标
	var durationMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_query_duration_seconds" {
			durationMetric = m
			break
		}
	}

	require.NotNil(t, durationMetric, "查询延迟指标应该存在")
	assert.True(t, len(durationMetric.Metric) > 0, "应该有metric实例")

	metric := durationMetric.Metric[0]

	// 验证label
	assert.Len(t, metric.Label, 2, "应该有两个label")

	// 查找collection和operation label
	var collectionLabel, operationLabel *dto.LabelPair
	for _, label := range metric.Label {
		if label.GetName() == "collection" {
			collectionLabel = label
		} else if label.GetName() == "operation" {
			operationLabel = label
		}
	}

	require.NotNil(t, collectionLabel, "collection label应该存在")
	assert.Equal(t, collection, collectionLabel.GetValue(), "collection label值应该正确")

	require.NotNil(t, operationLabel, "operation label应该存在")
	assert.Equal(t, operation, operationLabel.GetValue(), "operation label值应该正确")

	// 验证直方图类型
	assert.NotNil(t, metric.Histogram, "应该是Histogram类型")
	assert.NotNil(t, metric.Histogram.SampleCount, "应该有样本计数")
	assert.Equal(t, uint64(1), *metric.Histogram.SampleCount, "样本计数应该是1")
}

// TestSetIndexUsage 测试索引使用率设置功能
func TestSetIndexUsage(t *testing.T) {
	collection := "users"
	ratio := 0.85 // 85%索引使用率

	// 设置索引使用率
	SetIndexUsage(collection, ratio)

	// 验证指标值（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找索引使用率指标
	var usageMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_index_usage_ratio" {
			usageMetric = m
			break
		}
	}

	require.NotNil(t, usageMetric, "索引使用率指标应该存在")

	// 查找对应collection的metric
	var targetMetric *dto.Metric
	for _, m := range usageMetric.Metric {
		for _, label := range m.Label {
			if label.GetName() == "collection" && label.GetValue() == collection {
				targetMetric = m
				break
			}
		}
	}

	require.NotNil(t, targetMetric, "应该找到users的metric")

	// 验证label
	assert.Len(t, targetMetric.Label, 1, "应该有一个label")
	assert.Equal(t, "collection", targetMetric.Label[0].GetName(), "label名称应该是collection")
	assert.Equal(t, collection, targetMetric.Label[0].GetValue(), "label值应该是users")

	// 验证Gauge值
	assert.NotNil(t, targetMetric.Gauge, "应该是Gauge类型")
	assert.Equal(t, ratio, targetMetric.Gauge.GetValue(), "索引使用率应该是0.85")
}

// TestRecordQuery 测试查询类型记录功能
func TestRecordQuery(t *testing.T) {
	collection := "orders"
	operation := "update"
	hasIndex := true

	// 记录查询
	RecordQuery(collection, operation, hasIndex)

	// 验证指标值（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找查询类型指标
	var queryTypeMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_query_types_total" {
			queryTypeMetric = m
			break
		}
	}

	require.NotNil(t, queryTypeMetric, "查询类型指标应该存在")

	// 查找对应collection和operation的metric
	var targetMetric *dto.Metric
	for _, m := range queryTypeMetric.Metric {
		labels := make(map[string]string)
		for _, label := range m.Label {
			labels[label.GetName()] = label.GetValue()
		}
		if labels["collection"] == collection && labels["operation"] == operation {
			targetMetric = m
			break
		}
	}

	require.NotNil(t, targetMetric, "应该找到orders的metric")

	// 验证label
	assert.Len(t, targetMetric.Label, 3, "应该有三个label")

	// 查找labels
	labels := make(map[string]string)
	for _, label := range targetMetric.Label {
		labels[label.GetName()] = label.GetValue()
	}

	assert.Equal(t, collection, labels["collection"], "collection label值应该正确")
	assert.Equal(t, operation, labels["operation"], "operation label值应该正确")
	assert.Equal(t, "true", labels["has_index"], "has_index label应该是true")

	// 验证Counter值
	assert.NotNil(t, targetMetric.Counter, "应该是Counter类型")
	assert.Equal(t, float64(1), targetMetric.Counter.GetValue(), "计数器值应该是1")
}

// TestSetConnectionPool 测试连接池状态设置功能
func TestSetConnectionPool(t *testing.T) {
	state := "active"
	count := 10.0

	// 设置连接池状态
	SetConnectionPool(state, count)

	// 验证指标值（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找连接池指标
	var poolMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_connection_pool" {
			poolMetric = m
			break
		}
	}

	require.NotNil(t, poolMetric, "连接池指标应该存在")
	assert.True(t, len(poolMetric.Metric) > 0, "应该有metric实例")

	metric := poolMetric.Metric[0]

	// 验证label
	assert.Len(t, metric.Label, 1, "应该有一个label")
	assert.Equal(t, "state", metric.Label[0].GetName(), "label名称应该是state")
	assert.Equal(t, state, metric.Label[0].GetValue(), "label值应该是active")

	// 验证Gauge值
	assert.NotNil(t, metric.Gauge, "应该是Gauge类型")
	assert.Equal(t, count, metric.Gauge.GetValue(), "连接数应该是10")
}

// TestSetIndexEfficiency 测试索引效率设置功能
func TestSetIndexEfficiency(t *testing.T) {
	collection := "products"
	indexName := "idx_product_id"
	ratio := 0.92 // 92%命中率

	// 设置索引效率
	SetIndexEfficiency(collection, indexName, ratio)

	// 验证指标值（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找索引效率指标
	var efficiencyMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_index_efficiency" {
			efficiencyMetric = m
			break
		}
	}

	require.NotNil(t, efficiencyMetric, "索引效率指标应该存在")

	// 查找对应collection和indexName的metric
	var targetMetric *dto.Metric
	for _, m := range efficiencyMetric.Metric {
		labels := make(map[string]string)
		for _, label := range m.Label {
			labels[label.GetName()] = label.GetValue()
		}
		if labels["collection"] == collection && labels["index_name"] == indexName {
			targetMetric = m
			break
		}
	}

	require.NotNil(t, targetMetric, "应该找到products的metric")

	// 验证label
	assert.Len(t, targetMetric.Label, 2, "应该有两个label")

	labels := make(map[string]string)
	for _, label := range targetMetric.Label {
		labels[label.GetName()] = label.GetValue()
	}

	assert.Equal(t, collection, labels["collection"], "collection label值应该正确")
	assert.Equal(t, indexName, labels["index_name"], "index_name label值应该正确")

	// 验证Gauge值
	assert.NotNil(t, targetMetric.Gauge, "应该是Gauge类型")
	assert.Equal(t, ratio, targetMetric.Gauge.GetValue(), "索引效率应该是0.92")
}

// TestConcurrentMetricRecording 测试并发记录指标的安全性
func TestConcurrentMetricRecording(t *testing.T) {
	// 使用唯一的collection名称避免与其他测试冲突
	collection := "concurrent_test_unique"

	done := make(chan bool)

	// 并发记录指标
	for i := 0; i < 100; i++ {
		go func() {
			RecordSlowQuery(collection)
			RecordQuery(collection, "find", true)
			SetIndexUsage(collection, 0.8)
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 100; i++ {
		<-done
	}

	// 验证指标值正确（使用默认registry）
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	// 查找慢查询指标
	var slowQueryMetric *dto.MetricFamily
	for _, m := range metrics {
		if m.GetName() == "mongodb_slow_queries_total" {
			slowQueryMetric = m
			break
		}
	}

	require.NotNil(t, slowQueryMetric, "慢查询指标应该存在")

	// 查找对应collection的metric
	var targetMetric *dto.Metric
	for _, m := range slowQueryMetric.Metric {
		for _, label := range m.Label {
			if label.GetName() == "collection" && label.GetValue() == collection {
				targetMetric = m
				break
			}
		}
	}

	require.NotNil(t, targetMetric, "应该找到concurrent_test_unique的metric")
	assert.Equal(t, float64(100), targetMetric.Counter.GetValue(), "应该记录了100次慢查询")
}

// BenchmarkRecordSlowQuery 性能测试：慢查询记录
func BenchmarkRecordSlowQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordSlowQuery("benchmark_collection")
	}
}

// BenchmarkRecordQueryDuration 性能测试：查询延迟记录
func BenchmarkRecordQueryDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordQueryDuration("benchmark_collection", "find", 0.1)
	}
}
