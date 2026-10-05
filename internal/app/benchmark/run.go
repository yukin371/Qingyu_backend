package benchmarkapp

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"Qingyu_backend/benchmark"
)

type Config struct {
	BaseURL    string
	Name       string
	Requests   int
	Concurrent int
	WithCache  bool
	Output     string
	Timeout    time.Duration
}

func parseFlags(args []string) (*Config, error) {
	config := &Config{}
	fs := flag.NewFlagSet("benchmark", flag.ContinueOnError)

	fs.StringVar(&config.BaseURL, "base-url", "http://localhost:9090", "Base URL for testing")
	fs.StringVar(&config.Name, "name", "Performance Test", "Test scenario name")
	fs.IntVar(&config.Requests, "requests", 1000, "Total number of requests")
	fs.IntVar(&config.Concurrent, "concurrent", 50, "Number of concurrent requests")
	fs.BoolVar(&config.WithCache, "with-cache", true, "Enable cache")
	fs.StringVar(&config.Output, "output", "result.json", "Output JSON file path")
	fs.DurationVar(&config.Timeout, "timeout", 30*time.Minute, "Test timeout")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return config, nil
}

// Run executes the benchmark command with the provided arguments.
func Run(args []string) error {
	config, err := parseFlags(args)
	if err != nil {
		return err
	}

	b := benchmark.NewABTestBenchmark(config.BaseURL)

	scenario := benchmark.TestScenario{
		Name:       config.Name,
		Requests:   config.Requests,
		Concurrent: config.Concurrent,
		Endpoints:  []string{"/api/v1/bookstore/homepage"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	result, err := b.RunABTest(ctx, scenario, config.WithCache)
	if err != nil {
		return fmt.Errorf("测试失败: %w", err)
	}

	// 保存结果到JSON文件
	if err := saveResult(result, config.Output); err != nil {
		return fmt.Errorf("保存结果失败: %w", err)
	}

	// 输出摘要
	fmt.Printf("测试完成:\n")
	fmt.Printf("  总请求数: %d\n", result.TotalRequests)
	fmt.Printf("  成功: %d\n", result.SuccessCount)
	fmt.Printf("  失败: %d\n", result.ErrorCount)
	fmt.Printf("  平均延迟: %v\n", result.AvgLatency)
	fmt.Printf("  P95延迟: %v\n", result.P95Latency)
	fmt.Printf("  吞吐量: %.2f req/s\n", result.Throughput)
	return nil
}

func saveResult(result *benchmark.TestResult, path string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
