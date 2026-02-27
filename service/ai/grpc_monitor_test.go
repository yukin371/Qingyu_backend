package ai

import (
	"context"
	"testing"
	"time"
)

// TestGRPCMetrics 测试gRPC指标收集
func TestGRPCMetrics(t *testing.T) {
	metrics := NewGRPCMetrics()

	// 测试记录调用
	metrics.RecordCall(ServiceExecuteAgent, true)
	metrics.RecordCall(ServiceExecuteAgent, true)
	metrics.RecordCall(ServiceExecuteAgent, false)

	// 获取统计信息
	stats, err := metrics.GetStatsByService(ServiceExecuteAgent)
	if err != nil {
		t.Fatalf("GetStatsByService failed: %v", err)
	}

	if stats.Total != 3 {
		t.Errorf("Expected 3 total calls, got %d", stats.Total)
	}
	if stats.Success != 2 {
		t.Errorf("Expected 2 success calls, got %d", stats.Success)
	}
	if stats.Failed != 1 {
		t.Errorf("Expected 1 failed call, got %d", stats.Failed)
	}

	successRate := stats.GetSuccessRate()
	if successRate != 66.66666666666666 {
		t.Errorf("Expected success rate 66.67, got %f", successRate)
	}

	// 测试延迟记录
	metrics.RecordLatency(ServiceExecuteAgent, 100*time.Millisecond)
	metrics.RecordLatency(ServiceExecuteAgent, 200*time.Millisecond)

	latencyStats, err := metrics.GetLatencyStats(ServiceExecuteAgent)
	if err != nil {
		t.Fatalf("GetLatencyStats failed: %v", err)
	}

	if latencyStats.Count != 2 {
		t.Errorf("Expected 2 latency records, got %d", latencyStats.Count)
	}
	if latencyStats.Min != 100 {
		t.Errorf("Expected min latency 100ms, got %d", latencyStats.Min)
	}
	if latencyStats.Max != 200 {
		t.Errorf("Expected max latency 200ms, got %d", latencyStats.Max)
	}

	// 测试超时和重试记录
	metrics.RecordTimeout(ServiceExecuteAgent)
	metrics.RecordRetry(ServiceExecuteAgent)

	timeoutCount := metrics.GetTimeoutCount(ServiceExecuteAgent)
	if timeoutCount != 1 {
		t.Errorf("Expected 1 timeout, got %d", timeoutCount)
	}

	retryCount := metrics.GetRetryCount(ServiceExecuteAgent)
	if retryCount != 1 {
		t.Errorf("Expected 1 retry, got %d", retryCount)
	}
}

// TestTracer 测试请求追踪
func TestTracer(t *testing.T) {
	tracer := NewTracer(100)

	// 测试开始追踪
	serviceName := ServiceExecuteAgent
	requestID := "test-request-1"
	tracer.StartTrace(serviceName, requestID)

	// 测试获取活动追踪
	activeTraces := tracer.GetActiveTraces()
	if len(activeTraces) != 1 {
		t.Errorf("Expected 1 active trace, got %d", len(activeTraces))
	}

	// 测试结束追踪
	tracer.EndTrace(requestID, TraceStatusSuccess, nil)

	// 测试获取追踪信息
	trace, err := tracer.GetTrace(requestID)
	if err != nil {
		t.Fatalf("GetTrace failed: %v", err)
	}

	if trace.ServiceName != serviceName {
		t.Errorf("Expected service name %s, got %s", serviceName, trace.ServiceName)
	}
	if trace.Status != TraceStatusSuccess {
		t.Errorf("Expected status %s, got %s", TraceStatusSuccess, trace.Status)
	}

	// 测试获取最近追踪
	recentTraces := tracer.GetRecentTraces(10)
	if len(recentTraces) != 1 {
		t.Errorf("Expected 1 recent trace, got %d", len(recentTraces))
	}

	// 测试按服务获取追踪
	serviceTraces := tracer.GetTracesByService(serviceName, 10)
	if len(serviceTraces) != 1 {
		t.Errorf("Expected 1 service trace, got %d", len(serviceTraces))
	}

	// 测试按状态获取追踪
	successTraces := tracer.GetTracesByStatus(TraceStatusSuccess, 10)
	if len(successTraces) != 1 {
		t.Errorf("Expected 1 success trace, got %d", len(successTraces))
	}

	// 测试获取统计信息
	stats := tracer.GetStats()
	if stats.TotalTraces != 1 {
		t.Errorf("Expected 1 total trace, got %d", stats.TotalTraces)
	}
	if stats.ActiveTraces != 0 {
		t.Errorf("Expected 0 active traces, got %d", stats.ActiveTraces)
	}
}

// TestTracerWithError 测试带错误的追踪
func TestTracerWithError(t *testing.T) {
	tracer := NewTracer(100)

	serviceName := ServiceExecuteAgent
	requestID := "test-request-error"
	tracer.StartTrace(serviceName, requestID)

	// 结束追踪并记录错误
	testErr := context.DeadlineExceeded
	tracer.EndTrace(requestID, TraceStatusTimeout, testErr)

	// 获取追踪信息
	trace, err := tracer.GetTrace(requestID)
	if err != nil {
		t.Fatalf("GetTrace failed: %v", err)
	}

	if trace.Status != TraceStatusTimeout {
		t.Errorf("Expected status %s, got %s", TraceStatusTimeout, trace.Status)
	}
	if trace.Error == "" {
		t.Error("Expected error message, got empty string")
	}
}

// TestMetricsFormatReport 测试格式化报告
func TestMetricsFormatReport(t *testing.T) {
	metrics := NewGRPCMetrics()

	// 记录一些数据
	metrics.RecordCall(ServiceExecuteAgent, true)
	metrics.RecordCall(ServiceGenerateOutline, true)
	metrics.RecordCall(ServiceGenerateOutline, false)
	metrics.RecordLatency(ServiceExecuteAgent, 150*time.Millisecond)
	metrics.RecordLatency(ServiceGenerateOutline, 200*time.Millisecond)

	// 生成报告
	report := metrics.FormatReport()

	if report == "" {
		t.Error("Expected non-empty report")
	}

	// 检查报告内容是否包含关键信息
	if len(report) < 50 {
		t.Errorf("Report seems too short: %s", report)
	}
}

// TestTraceStats 测试追踪统计
func TestTraceStats(t *testing.T) {
	tracer := NewTracer(100)

	// 添加多个追踪
	for i := 0; i < 5; i++ {
		requestID := "test-request-stats-" + string(rune('0'+i))
		tracer.StartTrace(ServiceExecuteAgent, requestID)
		tracer.EndTrace(requestID, TraceStatusSuccess, nil)
	}

	// 添加一些失败的追踪
	for i := 0; i < 2; i++ {
		requestID := "test-request-fail-" + string(rune('0'+i))
		tracer.StartTrace(ServiceGenerateOutline, requestID)
		tracer.EndTrace(requestID, TraceStatusFailed, context.Canceled)
	}

	// 获取统计信息
	stats := tracer.GetStats()

	if stats.TotalTraces != 7 {
		t.Errorf("Expected 7 total traces, got %d", stats.TotalTraces)
	}
	if stats.ServiceCounts[ServiceExecuteAgent] != 5 {
		t.Errorf("Expected 5 ExecuteAgent traces, got %d", stats.ServiceCounts[ServiceExecuteAgent])
	}
	if stats.ServiceCounts[ServiceGenerateOutline] != 2 {
		t.Errorf("Expected 2 GenerateOutline traces, got %d", stats.ServiceCounts[ServiceGenerateOutline])
	}
	if stats.StatusCounts[TraceStatusSuccess] != 5 {
		t.Errorf("Expected 5 success traces, got %d", stats.StatusCounts[TraceStatusSuccess])
	}
	if stats.StatusCounts[TraceStatusFailed] != 2 {
		t.Errorf("Expected 2 failed traces, got %d", stats.StatusCounts[TraceStatusFailed])
	}
}

// TestMetricsReset 测试重置统计信息
func TestMetricsReset(t *testing.T) {
	metrics := NewGRPCMetrics()

	// 记录一些数据
	metrics.RecordCall(ServiceExecuteAgent, true)
	metrics.RecordCall(ServiceGenerateOutline, false)

	// 重置指定服务的统计信息
	metrics.ResetService(ServiceExecuteAgent)

	// 检查重置是否生效
	_, err := metrics.GetStatsByService(ServiceExecuteAgent)
	if err == nil {
		t.Error("Expected error when getting stats for reset service")
	}

	// 检查其他服务是否保留
	stats, err := metrics.GetStatsByService(ServiceGenerateOutline)
	if err != nil {
		t.Fatalf("GetStatsByService failed: %v", err)
	}
	if stats.Total != 1 {
		t.Errorf("Expected 1 total call, got %d", stats.Total)
	}

	// 重置所有统计信息
	metrics.Reset()

	// 检查是否全部重置
	allStats := metrics.GetStats()
	if len(allStats) != 0 {
		t.Errorf("Expected 0 services after reset, got %d", len(allStats))
	}
}

// TestTracerClear 测试清空追踪记录
func TestTracerClear(t *testing.T) {
	tracer := NewTracer(100)

	// 添加一些追踪
	for i := 0; i < 3; i++ {
		requestID := "test-request-clear-" + string(rune('0'+i))
		tracer.StartTrace(ServiceExecuteAgent, requestID)
		tracer.EndTrace(requestID, TraceStatusSuccess, nil)
	}

	// 清空所有追踪
	tracer.Clear()

	// 检查是否清空成功
	stats := tracer.GetStats()
	if stats.TotalTraces != 0 {
		t.Errorf("Expected 0 total traces after clear, got %d", stats.TotalTraces)
	}
	if stats.ActiveTraces != 0 {
		t.Errorf("Expected 0 active traces after clear, got %d", stats.ActiveTraces)
	}
}

// TestUnifiedClientMonitoring 测试UnifiedClient监控功能
func TestUnifiedClientMonitoring(t *testing.T) {
	// 注意：这个测试需要实际的gRPC连接，这里只测试监控组件的初始化
	metrics := NewGRPCMetrics()
	tracer := NewTracer(100)

	if metrics == nil {
		t.Error("Failed to create metrics")
	}
	if tracer == nil {
		t.Error("Failed to create tracer")
	}

	// 测试启用/禁用监控
	client := &UnifiedClient{
		enableMonitor: false,
	}

	if client.IsMonitoringEnabled() {
		t.Error("Expected monitoring to be disabled")
	}

	client.EnableMonitoring()

	if !client.IsMonitoringEnabled() {
		t.Error("Expected monitoring to be enabled")
	}

	if client.GetMetrics() == nil {
		t.Error("Expected metrics to be initialized")
	}

	if client.GetTracer() == nil {
		t.Error("Expected tracer to be initialized")
	}

	client.DisableMonitoring()

	if client.IsMonitoringEnabled() {
		t.Error("Expected monitoring to be disabled")
	}
}
