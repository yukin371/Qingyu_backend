package events

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============ Replay 指标测试 ============

// TestMetrics_RecordEventReplayed_Success 测试记录成功的事件回放
func TestMetrics_RecordEventReplayed_Success(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.1"
	replayedCount := int64(10)
	failedCount := int64(0)
	duration := 100 * time.Millisecond

	// Act
	metrics.RecordEventReplayed(eventType, replayedCount, failedCount, duration)

	// Assert
	// 验证 EventsReplayTotal 指标
	count := getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": eventType, "status": "success"})
	assert.Equal(t, float64(replayedCount), count, "EventsReplayTotal should match replayed count")

	// 验证 EventsReplayDuration 指标被记录
	histogram := getHistogramValue(metrics.EventsReplayDuration, map[string]string{"event_type": eventType})
	assert.NotNil(t, histogram, "EventsReplayDuration histogram should exist")
	assert.Greater(t, histogram.sampleCount, uint64(0), "Histogram should have samples")
	assert.Equal(t, duration.Seconds(), histogram.sampleSum, "Duration sum should match")
}

// newMetricsWithRegistry 使用指定的registry创建指标
func newMetricsWithRegistry(registry prometheus.Registerer) *Metrics {
	factory := promauto.With(registry)

	return &Metrics{
		// Replay 指标
		EventsReplayTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_events_replay_total",
				Help: "Total number of events replayed.",
			},
			[]string{"event_type", "status"},
		),
		EventsReplayFailed: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_events_replay_failed_total",
				Help: "Total number of failed event replay attempts.",
			},
			[]string{"event_type", "error_type"},
		),
		EventsReplayDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "test_events_replay_duration_seconds",
				Help:    "Duration of event replay in seconds.",
				Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60, 120, 300},
			},
			[]string{"event_type"},
		),
	}
}

// TestMetrics_RecordEventReplayed_PartialFailure 测试记录部分失败的事件回放
func TestMetrics_RecordEventReplayed_PartialFailure(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.2"
	replayedCount := int64(8)
	failedCount := int64(2)
	duration := 200 * time.Millisecond

	// Act
	metrics.RecordEventReplayed(eventType, replayedCount, failedCount, duration)

	// Assert
	count := getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": eventType, "status": "partial_failure"})
	assert.Equal(t, float64(replayedCount), count, "EventsReplayTotal should match replayed count")
}

// TestMetrics_RecordEventReplayFailed 测试记录回放失败
func TestMetrics_RecordEventReplayFailed(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.3"
	errorType := "handler_error"

	// Act
	metrics.RecordEventReplayFailed(eventType, errorType)

	// Assert
	count := getCounterValue(metrics.EventsReplayFailed, map[string]string{"event_type": eventType, "error_type": errorType})
	assert.Equal(t, float64(1), count, "EventsReplayFailed should be incremented")

	// 多次调用
	metrics.RecordEventReplayFailed(eventType, errorType)
	count = getCounterValue(metrics.EventsReplayFailed, map[string]string{"event_type": eventType, "error_type": errorType})
	assert.Equal(t, float64(2), count, "EventsReplayFailed should be incremented again")
}

// TestMetrics_RecordEventReplay_MultipleEventTypes 测试多种事件类型
func TestMetrics_RecordEventReplay_MultipleEventTypes(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventTypes := []string{"test.event.multi.a", "test.event.multi.b", "test.event.multi.c"}
	duration := 150 * time.Millisecond

	// Act
	for _, eventType := range eventTypes {
		metrics.RecordEventReplayed(eventType, 5, 0, duration)
	}

	// Assert
	for _, eventType := range eventTypes {
		count := getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": eventType, "status": "success"})
		assert.Equal(t, float64(5), count, "EventsReplayTotal for %s should match", eventType)
	}
}

// TestMetrics_RecordEventReplay_DurationBuckets 测试耗时分桶
func TestMetrics_RecordEventReplay_DurationBuckets(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.4"
	durations := []time.Duration{
		50 * time.Millisecond,  // 0.1s 桶
		300 * time.Millisecond, // 0.5s 桶
		2 * time.Second,        // 2.5s 桶
		30 * time.Second,       // 30s 桶
		150 * time.Second,      // 300s 桶
	}

	// Act
	for _, duration := range durations {
		metrics.RecordEventReplayed(eventType, 1, 0, duration)
	}

	// Assert
	histogram := getHistogramValue(metrics.EventsReplayDuration, map[string]string{"event_type": eventType})
	require.NotNil(t, histogram, "Histogram should exist")
	assert.Equal(t, uint64(5), histogram.sampleCount, "Should have 5 samples")
}

// TestMetrics_GlobalMetricsInstance 测试全局指标实例
func TestMetrics_GlobalMetricsInstance(t *testing.T) {
	// Arrange & Act
	globalMetrics := GetGlobalMetrics()

	// Assert
	assert.NotNil(t, globalMetrics, "Global metrics instance should not be nil")
	assert.NotNil(t, globalMetrics.EventsReplayTotal, "EventsReplayTotal should be initialized")
	assert.NotNil(t, globalMetrics.EventsReplayFailed, "EventsReplayFailed should be initialized")
	assert.NotNil(t, globalMetrics.EventsReplayDuration, "EventsReplayDuration should be initialized")
}

// TestMetrics_NewMetrics_CreatesAllMetrics 测试创建所有指标
func TestMetrics_NewMetrics_CreatesAllMetrics(t *testing.T) {
	// Act
	metrics := GetGlobalMetrics()

	// Assert
	require.NotNil(t, metrics, "Metrics should not be nil")

	// Replay 指标
	assert.NotNil(t, metrics.EventsReplayTotal, "EventsReplayTotal should not be nil")
	assert.NotNil(t, metrics.EventsReplayFailed, "EventsReplayFailed should not be nil")
	assert.NotNil(t, metrics.EventsReplayDuration, "EventsReplayDuration should not be nil")

	// 验证指标类型（通过检查它们是否为nil）
	assert.NotNil(t, metrics.EventsReplayTotal, "EventsReplayTotal should be initialized")
	assert.NotNil(t, metrics.EventsReplayFailed, "EventsReplayFailed should be initialized")
	assert.NotNil(t, metrics.EventsReplayDuration, "EventsReplayDuration should be initialized")
}

// TestMetrics_RecordEventReplay_ZeroCount 测试零计数情况
func TestMetrics_RecordEventReplay_ZeroCount(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.5"

	// Act
	metrics.RecordEventReplayed(eventType, 0, 0, 50*time.Millisecond)

	// Assert
	count := getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": eventType, "status": "success"})
	assert.Equal(t, float64(0), count, "EventsReplayTotal should be 0")
}

// TestMetrics_RecordEventReplay_AllFailed 测试全部失败情况
func TestMetrics_RecordEventReplay_AllFailed(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.6"
	replayedCount := int64(0)
	failedCount := int64(10)
	duration := 500 * time.Millisecond

	// Act
	metrics.RecordEventReplayed(eventType, replayedCount, failedCount, duration)

	// Assert
	count := getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": eventType, "status": "partial_failure"})
	assert.Equal(t, float64(0), count, "EventsReplayTotal should be 0 for all failures")
}

// TestMetrics_RecordEventReplayDifferentErrorTypes 测试不同错误类型
func TestMetrics_RecordEventReplayDifferentErrorTypes(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()
	eventType := "test.event.unique.7"
	errorTypes := []string{
		"query_error",
		"decode_error",
		"handler_error",
		"storage_error",
		"timeout_error",
	}

	// Act
	for _, errorType := range errorTypes {
		metrics.RecordEventReplayFailed(eventType, errorType)
	}

	// Assert
	for _, errorType := range errorTypes {
		count := getCounterValue(metrics.EventsReplayFailed, map[string]string{"event_type": eventType, "error_type": errorType})
		assert.Equal(t, float64(1), count, "EventsReplayFailed for %s should be 1", errorType)
	}
}

// ============ 辅助函数 ============

// getCounterValue 获取Counter的当前值
func getCounterValue(counterVec *prometheus.CounterVec, labels map[string]string) float64 {
	labelValues := make([]string, 0, len(labels))
	// 注意：这里假设label的顺序与定义时一致，实际使用时需要注意顺序
	// 为简化测试，这里使用简单的字符串匹配
	for _, v := range labels {
		labelValues = append(labelValues, v)
	}

	counter, err := counterVec.GetMetricWith(labelValuesToLabels(labels))
	if err != nil {
		return 0
	}

	return getCounterMetricValue(counter)
}

// labelValuesToLabels 将map转换为prometheus.Labels
func labelValuesToLabels(labels map[string]string) prometheus.Labels {
	result := make(prometheus.Labels)
	for k, v := range labels {
		result[k] = v
	}
	return result
}

// getCounterMetricValue 获取Counter指标的值
func getCounterMetricValue(counter prometheus.Counter) float64 {
	ch := make(chan prometheus.Metric, 1)
	counter.Collect(ch)
	close(ch)

	metric := <-ch
	if metric == nil {
		return 0
	}

	var pbMetric io_prometheus_client.Metric
	if err := metric.Write(&pbMetric); err != nil {
		return 0
	}

	return pbMetric.Counter.GetValue()
}

// getHistogramValue 获取Histogram的样本数和总和
func getHistogramValue(histogramVec *prometheus.HistogramVec, labels map[string]string) *histmetricData {
	labelValues := make([]string, 0, len(labels))
	for _, v := range labels {
		labelValues = append(labelValues, v)
	}

	histogram, err := histogramVec.GetMetricWith(labelValuesToLabels(labels))
	if err != nil {
		return nil
	}

	// 使用类型断言来获取prometheus.Histogram接口
	if h, ok := histogram.(prometheus.Histogram); ok {
		return getHistogramData(h)
	}
	return nil
}

type histmetricData struct {
	sampleCount uint64
	sampleSum   float64
}

func getHistogramData(histogram prometheus.Histogram) *histmetricData {
	ch := make(chan prometheus.Metric, 1)
	histogram.Collect(ch)
	close(ch)

	metric := <-ch
	if metric == nil {
		return nil
	}

	var pbMetric io_prometheus_client.Metric
	if err := metric.Write(&pbMetric); err != nil {
		return nil
	}

	h := pbMetric.Histogram
	if h == nil {
		return nil
	}

	return &histmetricData{
		sampleCount: h.GetSampleCount(),
		sampleSum:   h.GetSampleSum(),
	}
}

// TestMetrics_ReplayMetrics_Labels 测试指标标签
func TestMetrics_ReplayMetrics_Labels(t *testing.T) {
	// Arrange
	metrics := GetGlobalMetrics()

	// Act & Assert - EventsReplayTotal 标签
	t.Run("EventsReplayTotal labels", func(t *testing.T) {
		metrics.RecordEventReplayed("test.event.label.a", 1, 0, 100*time.Millisecond)

		// success status
		count := getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": "test.event.label.a", "status": "success"})
		assert.Equal(t, float64(1), count)

		// partial_failure status
		metrics.RecordEventReplayed("test.event.label.b", 5, 2, 100*time.Millisecond)
		count = getCounterValue(metrics.EventsReplayTotal, map[string]string{"event_type": "test.event.label.b", "status": "partial_failure"})
		assert.Equal(t, float64(5), count)
	})

	// Act & Assert - EventsReplayFailed 标签
	t.Run("EventsReplayFailed labels", func(t *testing.T) {
		metrics.RecordEventReplayFailed("test.event.label.c", "custom_error")
		count := getCounterValue(metrics.EventsReplayFailed, map[string]string{"event_type": "test.event.label.c", "error_type": "custom_error"})
		assert.Equal(t, float64(1), count)
	})
}
