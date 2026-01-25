package metrics_test

import (
	"fmt"
	"testing"
	"time"

	"Qingyu_backend/pkg/metrics"
)

// DemoSearchGrayscaleMetrics 演示灰度指标的使用
func DemoSearchGrayscaleMetrics() {
	// 1. 记录搜索请求
	metrics.RecordSearch("elasticsearch", 150*time.Millisecond)
	metrics.RecordSearch("mongodb", 80*time.Millisecond)
	metrics.RecordSearch("elasticsearch", 200*time.Millisecond)
	metrics.RecordSearch("mongodb", 95*time.Millisecond)

	// 2. 更新灰度百分比
	metrics.UpdateGrayscalePercent("books", 30)
	metrics.UpdateGrayscalePercent("projects", 50)
	metrics.UpdateGrayscalePercent("documents", 20)

	// 3. 记录灰度决策
	metrics.RecordGrayscaleDecision("books", true)
	metrics.RecordGrayscaleDecision("books", false)
	metrics.RecordGrayscaleDecision("projects", true)

	// 4. 获取指标快照
	snapshot := metrics.GetGrayScaleMetrics()
	fmt.Printf("ES Count: %d\n", snapshot.ESCount)
	fmt.Printf("MongoDB Count: %d\n", snapshot.MongoDBCount)

	// 5. 获取流量分配比例
	esPercent, mongoPercent := metrics.GetTrafficDistribution()
	fmt.Printf("ES Traffic: %.2f%%\n", esPercent)
	fmt.Printf("MongoDB Traffic: %.2f%%\n", mongoPercent)

	// 6. 获取平均耗时
	esAvg := metrics.GetAverageDuration("elasticsearch")
	mongoAvg := metrics.GetAverageDuration("mongodb")
	fmt.Printf("ES Avg Duration: %v\n", esAvg)
	fmt.Printf("MongoDB Avg Duration: %v\n", mongoAvg)
}

// TestRecordSearch 测试记录搜索请求
func TestRecordSearch(t *testing.T) {
	// 重置指标
	metrics.ResetMetrics()

	// 记录一些搜索
	metrics.RecordSearch("elasticsearch", 100*time.Millisecond)
	metrics.RecordSearch("mongodb", 50*time.Millisecond)
	metrics.RecordSearch("ES", 150*time.Millisecond) // 测试名称规范化
	metrics.RecordSearch("mongo", 75*time.Millisecond) // 测试名称规范化

	// 获取指标
	snapshot := metrics.GetGrayScaleMetrics()

	if snapshot.ESCount != 2 {
		t.Errorf("Expected ES count 2, got %d", snapshot.ESCount)
	}

	if snapshot.MongoDBCount != 2 {
		t.Errorf("Expected MongoDB count 2, got %d", snapshot.MongoDBCount)
	}

	// 检查总耗时
	expectedESTook := 250 * time.Millisecond
	if snapshot.ESTotalTook != expectedESTook {
		t.Errorf("Expected ES total took %v, got %v", expectedESTook, snapshot.ESTotalTook)
	}
}

// TestGetTrafficDistribution 测试流量分配比例
func TestGetTrafficDistribution(t *testing.T) {
	metrics.ResetMetrics()

	// 记录搜索
	for i := 0; i < 7; i++ {
		metrics.RecordSearch("elasticsearch", 100*time.Millisecond)
	}
	for i := 0; i < 3; i++ {
		metrics.RecordSearch("mongodb", 50*time.Millisecond)
	}

	// 获取流量分配
	esPercent, mongoPercent := metrics.GetTrafficDistribution()

	if esPercent != 70.0 {
		t.Errorf("Expected ES percent 70.0, got %.2f", esPercent)
	}

	if mongoPercent != 30.0 {
		t.Errorf("Expected MongoDB percent 30.0, got %.2f", mongoPercent)
	}
}

// TestGetAverageDuration 测试平均耗时
func TestGetAverageDuration(t *testing.T) {
	metrics.ResetMetrics()

	// 记录不同耗时的搜索
	durations := []time.Duration{
		100 * time.Millisecond,
		150 * time.Millisecond,
		200 * time.Millisecond,
	}

	for _, d := range durations {
		metrics.RecordSearch("elasticsearch", d)
	}

	// 计算预期平均值
	expectedAvg := (100*time.Millisecond + 150*time.Millisecond + 200*time.Millisecond) / 3

	// 获取平均耗时
	esAvg := metrics.GetAverageDuration("elasticsearch")

	if esAvg != expectedAvg {
		t.Errorf("Expected ES avg duration %v, got %v", expectedAvg, esAvg)
	}
}

// TestUpdateGrayscalePercent 测试更新灰度百分比
func TestUpdateGrayscalePercent(t *testing.T) {
	// 这个测试主要验证函数不会 panic
	// 实际的百分比存储在 Prometheus 中，可以通过 Prometheus 端点查询

	metrics.UpdateGrayscalePercent("books", 30)
	metrics.UpdateGrayscalePercent("projects", 50)
	metrics.UpdateGrayscalePercent("documents", 100)
	metrics.UpdateGrayscalePercent("users", 0)

	// 如果执行到这里没有 panic，测试通过
}

// TestRecordGrayscaleDecision 测试记录灰度决策
func TestRecordGrayscaleDecision(t *testing.T) {
	// 这个测试主要验证函数不会 panic
	// 实际的决策计数存储在 Prometheus 中

	metrics.RecordGrayscaleDecision("books", true)
	metrics.RecordGrayscaleDecision("books", false)
	metrics.RecordGrayscaleDecision("projects", true)
	metrics.RecordGrayscaleDecision("documents", false)

	// 如果执行到这里没有 panic，测试通过
}

// BenchmarkRecordSearch 性能基准测试
func BenchmarkRecordSearch(b *testing.B) {
	metrics.ResetMetrics()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.RecordSearch("elasticsearch", 100*time.Millisecond)
		}
	})
}
