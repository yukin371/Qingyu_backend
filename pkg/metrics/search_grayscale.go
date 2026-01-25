package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// GrayScaleMetricsSnapshot 灰度指标快照
type GrayScaleMetricsSnapshot struct {
	ESCount         int64         // ES 使用次数
	MongoDBCount    int64         // MongoDB 使用次数
	ESTotalTook     time.Duration // ES 总耗时
	MongoDBTotalTook time.Duration // MongoDB 总耗时
	LastUpdated     time.Time     // 最后更新时间
}

// grayscaleMetricsCollector 灰度指标收集器
type grayscaleMetricsCollector struct {
	mu              sync.RWMutex
	esCount         int64
	mongoCount      int64
	esTotalTook     time.Duration
	mongoTotalTook  time.Duration
	lastUpdated     time.Time
}

var (
	// 搜索请求总数
	searchRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_requests_total",
			Help: "Total number of search requests by engine",
		},
		[]string{"engine"}, // elasticsearch, mongodb
	)

	// 搜索耗时
	searchDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "search_duration_seconds",
			Help:    "Search request duration in seconds by engine",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"engine"},
	)

	// 灰度流量分配
	grayscaleTrafficGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_grayscale_traffic_percent",
			Help: "Grayscale traffic percentage for Elasticsearch by search type",
		},
		[]string{"search_type"}, // books, projects, documents
	)

	// 搜索引擎切换次数
	searchEngineSwitchesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_engine_switches_total",
			Help: "Total number of search engine switches during grayscale",
		},
		[]string{"from_engine", "to_engine"},
	)

	// 灰度决策结果分布
	grayscaleDecisionResult = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_grayscale_decision_total",
			Help: "Total number of grayscale decisions by result",
		},
		[]string{"search_type", "decision"}, // decision: es, mongodb
	)

	// 全局指标收集器实例
	collector = &grayscaleMetricsCollector{
		lastUpdated: time.Now(),
	}
)

// RecordSearch 记录搜索请求
func RecordSearch(engine string, took time.Duration) {
	// 规范化引擎名称
	normalizedEngine := normalizeEngineName(engine)

	// 更新 Prometheus Counter
	searchRequestsTotal.WithLabelValues(normalizedEngine).Inc()

	// 更新 Prometheus Histogram
	searchDuration.WithLabelValues(normalizedEngine).Observe(took.Seconds())

	// 更新内部指标
	collector.mu.Lock()
	defer collector.mu.Unlock()

	if normalizedEngine == "elasticsearch" {
		collector.esCount++
		collector.esTotalTook += took
	} else if normalizedEngine == "mongodb" {
		collector.mongoCount++
		collector.mongoTotalTook += took
	}
	collector.lastUpdated = time.Now()
}

// UpdateGrayscalePercent 更新灰度百分比
func UpdateGrayscalePercent(searchType string, percent int) {
	// 更新 Prometheus Gauge
	grayscaleTrafficGauge.WithLabelValues(searchType).Set(float64(percent))
}

// RecordGrayscaleDecision 记录灰度决策结果
func RecordGrayscaleDecision(searchType string, useES bool) {
	decision := "mongodb"
	if useES {
		decision = "elasticsearch"
	}
	grayscaleDecisionResult.WithLabelValues(searchType, decision).Inc()
}

// RecordEngineSwitch 记录引擎切换
func RecordEngineSwitch(fromEngine, toEngine string) {
	from := normalizeEngineName(fromEngine)
	to := normalizeEngineName(toEngine)
	searchEngineSwitchesTotal.WithLabelValues(from, to).Inc()
}

// GetGrayScaleMetrics 获取当前灰度指标快照
func GetGrayScaleMetrics() *GrayScaleMetricsSnapshot {
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	return &GrayScaleMetricsSnapshot{
		ESCount:          collector.esCount,
		MongoDBCount:     collector.mongoCount,
		ESTotalTook:      collector.esTotalTook,
		MongoDBTotalTook: collector.mongoTotalTook,
		LastUpdated:      collector.lastUpdated,
	}
}

// GetTrafficDistribution 获取流量分配比例
func GetTrafficDistribution() (esPercent, mongoPercent float64) {
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	total := collector.esCount + collector.mongoCount
	if total == 0 {
		return 0, 0
	}

	esPercent = float64(collector.esCount) / float64(total) * 100
	mongoPercent = float64(collector.mongoCount) / float64(total) * 100

	return esPercent, mongoPercent
}

// GetAverageDuration 获取平均搜索耗时
func GetAverageDuration(engine string) time.Duration {
	normalizedEngine := normalizeEngineName(engine)

	collector.mu.RLock()
	defer collector.mu.RUnlock()

	if normalizedEngine == "elasticsearch" {
		if collector.esCount > 0 {
			return collector.esTotalTook / time.Duration(collector.esCount)
		}
	} else if normalizedEngine == "mongodb" {
		if collector.mongoCount > 0 {
			return collector.mongoTotalTook / time.Duration(collector.mongoCount)
		}
	}

	return 0
}

// ResetMetrics 重置所有指标（主要用于测试）
func ResetMetrics() {
	collector.mu.Lock()
	defer collector.mu.Unlock()

	collector.esCount = 0
	collector.mongoCount = 0
	collector.esTotalTook = 0
	collector.mongoTotalTook = 0
	collector.lastUpdated = time.Now()
}

// normalizeEngineName 规范化引擎名称
func normalizeEngineName(engine string) string {
	switch engine {
	case "es", "ES", "elasticsearch", "Elasticsearch", "ElasticSearch":
		return "elasticsearch"
	case "mongodb", "MongoDB", "mongo", "Mongo":
		return "mongodb"
	default:
		return engine
	}
}

// GetPrometheusCollectors 获取所有 Prometheus 收集器（用于注册）
func GetPrometheusCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		searchRequestsTotal,
		searchDuration,
		grayscaleTrafficGauge,
		searchEngineSwitchesTotal,
		grayscaleDecisionResult,
	}
}
