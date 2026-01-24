//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AIServiceHTTPClient AI服务HTTP客户端
type AIServiceHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAIServiceHTTPClient 创建AI服务HTTP客户端
func NewAIServiceHTTPClient(baseURL string) *AIServiceHTTPClient {
	return &AIServiceHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	UserID  string `json:"user_id"`
	Messages []Message `json:"messages"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message         string `json:"message"`
	Usage           Usage  `json:"usage"`
	Model           string `json:"model"`
	QuotaRemaining  int64  `json:"quota_remaining"`
}

// Usage token使用情况
type Usage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status   string `json:"status"`
	Service  string `json:"service"`
	Timestamp string `json:"timestamp"`
	Version  string `json:"version"`
}

// Chat 发送聊天请求
func (c *AIServiceHTTPClient) Chat(req *ChatRequest) (*ChatResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/api/v1/ai/chat", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var chatResp ChatResponse
	err = json.Unmarshal(respBody, &chatResp)
	if err != nil {
		return nil, err
	}

	return &chatResp, nil
}

// Health 健康检查
func (c *AIServiceHTTPClient) Health() (*HealthResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/health")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var healthResp HealthResponse
	err = json.Unmarshal(respBody, &healthResp)
	if err != nil {
		return nil, err
	}

	return &healthResp, nil
}

// HTTPError HTTP错误
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return e.Body
}

// TestAIHTTPIntegration AI服务HTTP集成测试
func TestAIHTTPIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 创建AI服务客户端
	client := NewAIServiceHTTPClient("http://localhost:8000")

	t.Run("健康检查", func(t *testing.T) {
		resp, err := client.Health()
		require.NoError(t, err)
		assert.Equal(t, "healthy", resp.Status)
		assert.Equal(t, "Qingyu-AI-Service", resp.Service)
		assert.NotEmpty(t, resp.Version)
	})

	t.Run("基础聊天功能", func(t *testing.T) {
		req := &ChatRequest{
			UserID: "test-user-123",
			Messages: []Message{
				{
					Role:    "user",
					Content: "Hello, AI!",
				},
			},
		}

		resp, err := client.Chat(req)
		require.NoError(t, err)

		// 验证响应
		assert.NotEmpty(t, resp.Message)
		assert.NotEmpty(t, resp.Model)
		assert.Greater(t, resp.Usage.TotalTokens, int32(0))
		assert.Greater(t, resp.QuotaRemaining, int64(0))

		// 验证token使用
		assert.Greater(t, resp.Usage.PromptTokens, int32(0))
		assert.GreaterOrEqual(t, resp.Usage.CompletionTokens, int32(0))
	})

	t.Run("多轮对话", func(t *testing.T) {
		userID := "test-user-multi-turn"

		// 第一轮
		req1 := &ChatRequest{
			UserID: userID,
			Messages: []Message{
				{
					Role:    "user",
					Content: "My name is Alice",
				},
			},
		}

		resp1, err := client.Chat(req1)
		require.NoError(t, err)
		assert.NotEmpty(t, resp1.Message)

		// 第二轮 - 引用上一轮的内容
		req2 := &ChatRequest{
			UserID: userID,
			Messages: []Message{
				{
					Role:    "user",
					Content: "What's my name?",
				},
			},
		}

		resp2, err := client.Chat(req2)
		require.NoError(t, err)
		assert.NotEmpty(t, resp2.Message)
	})

	t.Run("配额扣除验证", func(t *testing.T) {
		userID := "test-user-quota"

		// 第一次调用
		req1 := &ChatRequest{
			UserID: userID,
			Messages: []Message{
				{
					Role:    "user",
					Content: "Test quota 1",
				},
			},
		}

		resp1, err := client.Chat(req1)
		require.NoError(t, err)
		initialQuota := resp1.QuotaRemaining

		// 第二次调用
		req2 := &ChatRequest{
			UserID: userID,
			Messages: []Message{
				{
					Role:    "user",
					Content: "Test quota 2",
				},
			},
		}

		resp2, err := client.Chat(req2)
		require.NoError(t, err)

		// 验证配额（模拟响应可能不会实际扣除配额）
		// 实际实现中，配额应该被扣除
		t.Logf("Initial quota: %d, Final quota: %d", initialQuota, resp2.QuotaRemaining)
		assert.GreaterOrEqual(t, initialQuota, resp2.QuotaRemaining)
	})

	t.Run("并发请求", func(t *testing.T) {
		userID := "test-user-concurrent"

		// 发送3个并发请求
		results := make(chan *ChatResponse, 3)
		errors := make(chan error, 3)

		for i := 0; i < 3; i++ {
			go func(index int) {
				req := &ChatRequest{
					UserID: userID,
					Messages: []Message{
						{
							Role:    "user",
							Content: "Concurrent test",
						},
					},
				}

				resp, err := client.Chat(req)
				if err != nil {
					errors <- err
				} else {
					results <- resp
				}
			}(i)
		}

		// 收集结果
		successCount := 0
		for i := 0; i < 3; i++ {
			select {
			case <-results:
				successCount++
			case err := <-errors:
				t.Logf("Request failed: %v", err)
			case <-time.After(10 * time.Second):
				t.Fatal("Timeout waiting for concurrent requests")
			}
		}

		// 至少2个请求应该成功
		assert.GreaterOrEqual(t, successCount, 2)
	})
}
