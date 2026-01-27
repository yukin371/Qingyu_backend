package benchmark

import "time"

// TestScenario 测试场景定义
type TestScenario struct {
	Name       string
	Requests   int
	Concurrent int
	Endpoints  []string
}

// TestResult 测试结果
type TestResult struct {
	Scenario      string
	WithCache     bool
	TotalRequests int
	SuccessCount  int
	ErrorCount    int
	AvgLatency    time.Duration
	P95Latency    time.Duration
	P99Latency    time.Duration
	Throughput    float64 // req/s
	Duration      time.Duration
}
