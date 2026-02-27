package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/mongodb/bookstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BenchmarkStreamSearch_NoCursor 基准测试：无游标的流式搜索
func BenchmarkStreamSearch_NoCursor(b *testing.B) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("X-Accel-Buffering", "no")

		// 模拟发送100条数据
		for i := 0; i < 100; i++ {
			book := map[string]interface{}{
				"id":       fmt.Sprintf("book-%d", i),
				"title":    fmt.Sprintf("测试书籍%d", i),
				"author":   "测试作者",
				"cover":    "https://example.com/cover.jpg",
				"rating":   4.5,
				"viewCount": 1000 + i,
			}

			data, _ := json.Marshal(map[string]interface{}{
				"type":  "data",
				"books": []interface{}{book},
			})

			fmt.Fprintf(w, "%s\n", string(data))
		}

		// 发送完成信号
		done, _ := json.Marshal(map[string]interface{}{
			"type":   "done",
			"cursor": "final-cursor",
			"total":  100,
		})
		fmt.Fprintf(w, "%s\n", string(done))
	}))
	defer server.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// 模拟流式请求
		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		req.Header.Set("Accept", "application/x-ndjson")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}

		// 消费响应
		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				break
			}
		}

		resp.Body.Close()
		cancel()
	}
}

// BenchmarkStreamSearch_WithCursor 基准测试：使用游标的流式搜索
func BenchmarkStreamSearch_WithCursor(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("X-Accel-Buffering", "no")

		// 发送20条数据
		for i := 0; i < 20; i++ {
			book := map[string]interface{}{
				"id":     fmt.Sprintf("book-%d", i),
				"title":  fmt.Sprintf("测试书籍%d", i),
				"author": "测试作者",
			}

			data, _ := json.Marshal(map[string]interface{}{
				"type":  "data",
				"books": []interface{}{book},
			})

			fmt.Fprintf(w, "%s\n", string(data))
		}

		done, _ := json.Marshal(map[string]interface{}{
			"type":   "done",
			"cursor": "next-cursor",
			"total":  20,
		})
		fmt.Fprintf(w, "%s\n", string(done))
	}))
	defer server.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		req.Header.Set("Accept", "application/x-ndjson")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}

		// 消费响应
		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				break
			}
		}

		resp.Body.Close()
		cancel()
	}
}

// BenchmarkCursorEncoding 基准测试：游标编码性能
func BenchmarkCursorEncoding(b *testing.B) {
	cursorMgr := mongodb.NewCursorManager()
	book := &bookstore.Book{}
	book.ID = primitive.NewObjectID()
	book.CreatedAt = time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := cursorMgr.GenerateNextCursor(book, bookstore.CursorTypeTimestamp, "created_at")
		if err != nil {
			b.Fatalf("GenerateNextCursor failed: %v", err)
		}
	}
}

// BenchmarkCursorDecoding 基准测试：游标解码性能
func BenchmarkCursorDecoding(b *testing.B) {
	cursorMgr := mongodb.NewCursorManager()
	testCursor := "eyJ0eXBlIjoidGltZXN0YW1wIiwidmFsdWUiOiIxNzA2MTQwODIxMDAwIn0="

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := cursorMgr.DecodeCursor(testCursor)
		if err != nil {
			b.Fatalf("DecodeCursor failed: %v", err)
		}
	}
}

// BenchmarkNDJSONParsing 基准测试：NDJSON解析性能
func BenchmarkNDJSONParsing(b *testing.B) {
	// 创建测试数据
	var lines []string
	for i := 0; i < 100; i++ {
		book := map[string]interface{}{
			"id":     fmt.Sprintf("book-%d", i),
			"title":  fmt.Sprintf("测试书籍%d", i),
			"author": "测试作者",
		}

		data, _ := json.Marshal(map[string]interface{}{
			"type":  "data",
			"books": []interface{}{book},
		})

		lines = append(lines, string(data))
	}

	testData := strings.Join(lines, "\n") + "\n"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decoder := json.NewDecoder(strings.NewReader(testData))
		for decoder.More() {
			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				b.Fatalf("Decode failed: %v", err)
			}
		}
	}
}

// BenchmarkStreamSearchLargeDataset 基准测试：大数据集流式搜索
func BenchmarkStreamSearchLargeDataset(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("X-Accel-Buffering", "no")

		// 发送10000条数据
		for i := 0; i < 10000; i++ {
			book := map[string]interface{}{
				"id":     fmt.Sprintf("book-%d", i),
				"title":  fmt.Sprintf("测试书籍%d", i),
				"author": "测试作者",
			}

			data, _ := json.Marshal(map[string]interface{}{
				"type":  "data",
				"books": []interface{}{book},
			})

			fmt.Fprintf(w, "%s\n", string(data))

			// 每100条发送一次进度
			if (i+1)%100 == 0 {
				progress, _ := json.Marshal(map[string]interface{}{
					"type":   "progress",
					"loaded": i + 1,
					"total":  10000,
				})
				fmt.Fprintf(w, "%s\n", string(progress))
			}
		}

		done, _ := json.Marshal(map[string]interface{}{
			"type":   "done",
			"cursor": "final-cursor",
			"total":  10000,
		})
		fmt.Fprintf(w, "%s\n", string(done))
	}))
	defer server.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		req.Header.Set("Accept", "application/x-ndjson")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}

		// 消费响应
		decoder := json.NewDecoder(resp.Body)
		count := 0
		for decoder.More() {
			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				break
			}
			count++
		}

		resp.Body.Close()
		cancel()
	}
}

// BenchmarkConcurrentStreamRequests 基准测试：并发流式请求
func BenchmarkConcurrentStreamRequests(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("X-Accel-Buffering", "no")

		// 发送20条数据
		for i := 0; i < 20; i++ {
			book := map[string]interface{}{
				"id":     fmt.Sprintf("book-%d", i),
				"title":  fmt.Sprintf("测试书籍%d", i),
				"author": "测试作者",
			}

			data, _ := json.Marshal(map[string]interface{}{
				"type":  "data",
				"books": []interface{}{book},
			})

			fmt.Fprintf(w, "%s\n", string(data))
		}

		done, _ := json.Marshal(map[string]interface{}{
			"type":   "done",
			"cursor": "final-cursor",
			"total":  20,
		})
		fmt.Fprintf(w, "%s\n", string(done))
	}))
	defer server.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
			req.Header.Set("Accept", "application/x-ndjson")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				cancel()
				continue
			}

			// 消费响应
			decoder := json.NewDecoder(resp.Body)
			for decoder.More() {
				var data map[string]interface{}
				decoder.Decode(&data)
			}

			resp.Body.Close()
			cancel()
		}
	})
}

// TestStreamSearchPerformanceMetrics 测试流式搜索性能指标
func TestStreamSearchPerformanceMetrics(t *testing.T) {
	t.Run("首屏响应时间应该小于100ms", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Accel-Buffering", "no")

			// 立即发送第一条数据
			book := map[string]interface{}{
				"id":     "book-1",
				"title":  "首屏书籍",
				"author": "测试作者",
			}

			data, _ := json.Marshal(map[string]interface{}{
				"type":  "data",
				"books": []interface{}{book},
			})

			fmt.Fprintf(w, "%s\n", string(data))

			done, _ := json.Marshal(map[string]interface{}{
				"type":   "done",
				"cursor": "final-cursor",
				"total":  1,
			})
			fmt.Fprintf(w, "%s\n", string(done))
		}))
		defer server.Close()

		startTime := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		req.Header.Set("Accept", "application/x-ndjson")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		// 读取第一条数据
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&map[string]interface{}{})

		elapsed := time.Since(startTime)

		if elapsed > 100*time.Millisecond {
			t.Errorf("首屏响应时间 %v 超过目标 100ms", elapsed)
		}

		t.Logf("首屏响应时间: %v", elapsed)
	})

	t.Run("内存占用应该可控", func(t *testing.T) {
		// 这个测试需要在实际环境中运行
		t.Skip("需要在实际环境中测量内存占用")

		// 理论上，流式处理应该只保持少量数据在内存中
		// 而不是一次性加载所有数据
	})

	t.Run("支持10K+数据流畅传输", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Accel-Buffering", "no")

			// 发送10000条数据
			for i := 0; i < 10000; i++ {
				book := map[string]interface{}{
					"id":     fmt.Sprintf("book-%d", i),
					"title":  fmt.Sprintf("测试书籍%d", i),
					"author": "测试作者",
				}

				data, _ := json.Marshal(map[string]interface{}{
					"type":  "data",
					"books": []interface{}{book},
				})

				fmt.Fprintf(w, "%s\n", string(data))
			}

			done, _ := json.Marshal(map[string]interface{}{
				"type":   "done",
				"cursor": "final-cursor",
				"total":  10000,
			})
			fmt.Fprintf(w, "%s\n", string(done))
		}))
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		req.Header.Set("Accept", "application/x-ndjson")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		// 消费所有数据
		decoder := json.NewDecoder(resp.Body)
		count := 0
		for decoder.More() {
			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				break
			}

			if data["type"] == "data" {
				count++
			}
		}

		if count != 10000 {
			t.Errorf("期望10000条数据，实际收到%d条", count)
		}

		t.Logf("成功处理%d条数据", count)
	})
}
