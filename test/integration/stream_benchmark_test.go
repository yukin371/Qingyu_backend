package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"Qingyu_backend/service/ai/adapter"
)

// BenchmarkOpenAIStream 基准测试OpenAI流式响应性能
func BenchmarkOpenAIStream(b *testing.B) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 模拟发送多个数据块
		for i := 0; i < 10; i++ {
			fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"content\":\"chunk %d \"}}]}\n\n", i)
			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Millisecond) // 模拟网络延迟
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建OpenAI适配器
	openaiAdapter := adapter.NewOpenAIAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Generate a long response",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   1000,
		Temperature: 0.7,
		Stream:      true,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			responseChan, err := openaiAdapter.TextGenerationStream(ctx, req)
			if err != nil {
				b.Fatalf("流式请求失败: %v", err)
			}

			// 消费所有响应
			var responseCount int
			for range responseChan {
				responseCount++
			}

			if responseCount == 0 {
				b.Fatal("未收到任何响应")
			}

			cancel()
		}
	})
}

// BenchmarkClaudeStream 基准测试Claude流式响应性能
func BenchmarkClaudeStream(b *testing.B) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 模拟发送多个数据块
		for i := 0; i < 10; i++ {
			fmt.Fprintf(w, "data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\"chunk %d \"}}\n\n", i)
			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Millisecond) // 模拟网络延迟
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建Claude适配器
	claudeAdapter := adapter.NewClaudeAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Generate a long response",
		Model:       "claude-3-haiku-20240307",
		MaxTokens:   1000,
		Temperature: 0.7,
		Stream:      true,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			responseChan, err := claudeAdapter.TextGenerationStream(ctx, req)
			if err != nil {
				b.Fatalf("流式请求失败: %v", err)
			}

			// 消费所有响应
			var responseCount int
			for range responseChan {
				responseCount++
			}

			if responseCount == 0 {
				b.Fatal("未收到任何响应")
			}

			cancel()
		}
	})
}

// BenchmarkGeminiStream 基准测试Gemini流式响应性能
func BenchmarkGeminiStream(b *testing.B) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 模拟发送多个数据块
		for i := 0; i < 10; i++ {
			fmt.Fprintf(w, "data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"chunk %d \"}]}}]}\n\n", i)
			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Millisecond) // 模拟网络延迟
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建Gemini适配器
	geminiAdapter := adapter.NewGeminiAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Generate a long response",
		Model:       "gemini-pro",
		MaxTokens:   1000,
		Temperature: 0.7,
		Stream:      true,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			responseChan, err := geminiAdapter.TextGenerationStream(ctx, req)
			if err != nil {
				b.Fatalf("流式请求失败: %v", err)
			}

			// 消费所有响应
			var responseCount int
			for range responseChan {
				responseCount++
			}

			if responseCount == 0 {
				b.Fatal("未收到任何响应")
			}

			cancel()
		}
	})
}

// BenchmarkStreamResponseProcessing 基准测试流式响应处理性能
func BenchmarkStreamResponseProcessing(b *testing.B) {
	// 模拟大量流式数据
	testData := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		testData[i] = fmt.Sprintf("data: {\"choices\":[{\"delta\":{\"content\":\"word%d \"}}]}\n\n", i)
	}
	testData = append(testData, "data: [DONE]\n\n")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 创建模拟HTTP服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				for _, data := range testData {
					fmt.Fprint(w, data)
					w.(http.Flusher).Flush()
				}
			}))

			// 创建OpenAI适配器
			openaiAdapter := adapter.NewOpenAIAdapter("test-key", server.URL)

			// 创建请求
			req := &adapter.TextGenerationRequest{
				Prompt:      "Generate a very long response",
				Model:       "gpt-3.5-turbo",
				MaxTokens:   2000,
				Temperature: 0.7,
				Stream:      true,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			responseChan, err := openaiAdapter.TextGenerationStream(ctx, req)
			if err != nil {
				b.Fatalf("流式请求失败: %v", err)
			}

			// 处理所有响应并构建完整内容
			var fullContent strings.Builder
			var responseCount int
			for response := range responseChan {
				fullContent.WriteString(response.Text)
				responseCount++
			}

			if responseCount == 0 {
				b.Fatal("未收到任何响应")
			}

			if fullContent.Len() == 0 {
				b.Fatal("未构建完整内容")
			}

			cancel()
			server.Close()
		}
	})
}

// BenchmarkConcurrentStreams 基准测试并发流式请求性能
func BenchmarkConcurrentStreams(b *testing.B) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 模拟发送5个数据块
		for i := 0; i < 5; i++ {
			fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"content\":\"chunk %d \"}}]}\n\n", i)
			w.(http.Flusher).Flush()
			time.Sleep(2 * time.Millisecond) // 模拟网络延迟
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建OpenAI适配器
	openaiAdapter := adapter.NewOpenAIAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Generate response",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   100,
		Temperature: 0.7,
		Stream:      true,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 并发执行多个流式请求
			const concurrency = 10
			done := make(chan bool, concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					defer func() { done <- true }()

					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					responseChan, err := openaiAdapter.TextGenerationStream(ctx, req)
					if err != nil {
						return
					}

					// 消费所有响应
					for range responseChan {
						// 处理响应
					}
				}()
			}

			// 等待所有协程完成
			for i := 0; i < concurrency; i++ {
				<-done
			}
		}
	})
}
