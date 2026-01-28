package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.BaseURL, "base-url", "http://localhost:8080", "Base URL for testing")
	flag.StringVar(&config.Name, "name", "Performance Test", "Test scenario name")
	flag.IntVar(&config.Requests, "requests", 1000, "Total number of requests")
	flag.IntVar(&config.Concurrent, "concurrent", 50, "Number of concurrent requests")
	flag.BoolVar(&config.WithCache, "with-cache", true, "Enable cache")
	flag.StringVar(&config.Output, "output", "result.json", "Output JSON file path")
	flag.DurationVar(&config.Timeout, "timeout", 30*time.Minute, "Test timeout")

	flag.Parse()
	return config
}

func main() {
	config := parseFlags()

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
		log.Fatalf("测试失败: %v", err)
	}

	// 保存结果到JSON文件
	if err := saveResult(result, config.Output); err != nil {
		log.Fatalf("保存结果失败: %v", err)
	}

	// 输出摘要
	fmt.Printf("测试完成:\n")
	fmt.Printf("  总请求数: %d\n", result.TotalRequests)
	fmt.Printf("  成功: %d\n", result.SuccessCount)
	fmt.Printf("  失败: %d\n", result.ErrorCount)
	fmt.Printf("  平均延迟: %v\n", result.AvgLatency)
	fmt.Printf("  P95延迟: %v\n", result.P95Latency)
	fmt.Printf("  吞吐量: %.2f req/s\n", result.Throughput)
}

func saveResult(result *benchmark.TestResult, path string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
