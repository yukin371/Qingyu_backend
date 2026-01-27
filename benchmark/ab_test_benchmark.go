package benchmark

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

// ABTestBenchmark A/B测试基准测试工具
type ABTestBenchmark struct {
	client   *http.Client
	baseURL  string
	resultMu sync.Mutex // 互斥锁保护result字段
}

// NewABTestBenchmark 创建A/B测试基准测试工具
func NewABTestBenchmark(baseURL string) *ABTestBenchmark {
	return &ABTestBenchmark{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// makeRequest 执行HTTP请求
func (b *ABTestBenchmark) makeRequest(ctx context.Context, endpoint string) error {
	url := b.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// RunABTest 执行A/B测试
func (b *ABTestBenchmark) RunABTest(
	ctx context.Context,
	scenario TestScenario,
	withCache bool,
) (*TestResult, error) {
	result := &TestResult{
		Scenario:      scenario.Name,
		WithCache:     withCache,
		TotalRequests: scenario.Requests,
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, scenario.Concurrent)
	latencies := make([]time.Duration, scenario.Requests)

	startTime := time.Now()

	for i := 0; i < scenario.Requests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			reqStart := time.Now()
			err := b.makeRequest(ctx, scenario.Endpoints[idx%len(scenario.Endpoints)])
			latency := time.Since(reqStart)

			// 使用互斥锁保护并发写入
			b.resultMu.Lock()
			if err != nil {
				result.ErrorCount++
			} else {
				result.SuccessCount++
			}
			b.resultMu.Unlock()

			latencies[idx] = latency
		}(i)
	}

	wg.Wait()
	result.Duration = time.Since(startTime)

	// 计算统计数据
	result.calculateStatistics(latencies)

	return result, nil
}

// calculateStatistics 计算统计数据
func (r *TestResult) calculateStatistics(latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}

	// 计算平均延迟
	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	r.AvgLatency = total / time.Duration(len(latencies))

	// 使用标准库排序 (O(n log n))
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	p95Index := int(float64(len(sorted)) * 0.95)
	p99Index := int(float64(len(sorted)) * 0.99)

	if p95Index < len(sorted) {
		r.P95Latency = sorted[p95Index]
	}
	if p99Index < len(sorted) {
		r.P99Latency = sorted[p99Index]
	}

	// 计算吞吐量
	r.Throughput = float64(r.TotalRequests) / r.Duration.Seconds()
}
