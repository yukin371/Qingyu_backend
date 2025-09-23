package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"Qingyu_backend/api/v1/ai"
	aiModels "Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"
	"Qingyu_backend/service/ai/adapter"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CloseNotifierResponseRecorder 支持CloseNotifier的测试响应记录器
type CloseNotifierResponseRecorder struct {
	*httptest.ResponseRecorder
	closeCh chan bool
}

// NewCloseNotifierResponseRecorder 创建新的CloseNotifierResponseRecorder
func NewCloseNotifierResponseRecorder() *CloseNotifierResponseRecorder {
	return &CloseNotifierResponseRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closeCh:          make(chan bool, 1),
	}
}

// CloseNotify 实现http.CloseNotifier接口
func (c *CloseNotifierResponseRecorder) CloseNotify() <-chan bool {
	return c.closeCh
}

// Close 模拟连接关闭
func (c *CloseNotifierResponseRecorder) Close() {
	select {
	case c.closeCh <- true:
	default:
	}
}

// MockChatService 模拟聊天服务
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) StartChat(ctx context.Context, req *aiService.ChatRequest) (*aiService.ChatResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*aiService.ChatResponse), args.Error(1)
}

func (m *MockChatService) StartChatStream(ctx context.Context, req *aiService.ChatRequest) (<-chan *aiService.StreamChatResponse, error) {
	args := m.Called(ctx, req)
	// 直接返回通道，不进行类型转换
	if ch, ok := args.Get(0).(chan *aiService.StreamChatResponse); ok {
		return ch, args.Error(1)
	}
	return args.Get(0).(<-chan *aiService.StreamChatResponse), args.Error(1)
}

func (m *MockChatService) GetChatHistory(ctx context.Context, sessionID string) (*aiService.ChatSession, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*aiService.ChatSession), args.Error(1)
}

func (m *MockChatService) ListChatSessions(ctx context.Context, projectID string, limit, offset int) ([]*aiService.ChatSession, error) {
	args := m.Called(ctx, projectID, limit, offset)
	return args.Get(0).([]*aiService.ChatSession), args.Error(1)
}

func (m *MockChatService) DeleteChatSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockChatService) GetStatistics(projectID string) (*aiService.ChatStatistics, error) {
	args := m.Called(projectID)
	return args.Get(0).(*aiService.ChatStatistics), args.Error(1)
}

// TestStreamChatAPI 测试流式聊天API
func TestStreamChatAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockService := new(MockChatService)
	
	// 创建流式响应通道
	streamChan := make(chan *aiService.StreamChatResponse, 3)
	go func() {
		defer close(streamChan)
		streamChan <- &aiService.StreamChatResponse{
			SessionID: "test-session",
			MessageID: "test-message-1",
			Content:   "Hello",
			Delta:     "Hello",
			IsComplete: false,
		}
		streamChan <- &aiService.StreamChatResponse{
			SessionID: "test-session",
			MessageID: "test-message-1",
			Content:   "Hello, how can I help you?",
			Delta:     ", how can I help you?",
			IsComplete: false,
		}
		streamChan <- &aiService.StreamChatResponse{
			SessionID: "test-session",
			MessageID: "test-message-1",
			Content:   "Hello, how can I help you?",
			Delta:     "",
			IsComplete: true,
		}
	}()

	// 设置模拟期望
	mockService.On("StartChatStream", mock.Anything, mock.AnythingOfType("*ai.ChatRequest")).Return(streamChan, nil)

	// 创建API实例
	chatAPI := ai.NewChatApi(mockService)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/ai/chat/stream", chatAPI.ContinueChat)

	// 准备请求数据
	requestData := aiService.ChatRequest{
		ProjectID: "test-project",
		Message:   "Hello",
		Options: &aiModels.GenerateOptions{
			Stream: true,
		},
	}
	requestBody, _ := json.Marshal(requestData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat/stream", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := NewCloseNotifierResponseRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	assert.Equal(t, "keep-alive", w.Header().Get("Connection"))

	// 验证SSE响应内容
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "event: message")
	assert.Contains(t, responseBody, "data:")
	assert.Contains(t, responseBody, "Hello")
	assert.Contains(t, responseBody, "event: end")

	// 验证模拟调用
	mockService.AssertExpectations(t)
}

// TestOpenAIAdapterStream 测试OpenAI适配器流式响应
func TestOpenAIAdapterStream(t *testing.T) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
		assert.Equal(t, "text/event-stream", r.Header.Get("Accept"))

		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 发送模拟SSE数据
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\" world\"}}]}\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建OpenAI适配器
	openaiAdapter := adapter.NewOpenAIAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Hello",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   100,
		Temperature: 0.7,
		Stream:      true,
	}

	// 执行流式请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	responseChan, err := openaiAdapter.TextGenerationStream(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, responseChan)

	// 收集响应
	var responses []*adapter.TextGenerationResponse
	for response := range responseChan {
		responses = append(responses, response)
	}

	// 验证响应
	assert.Greater(t, len(responses), 0)
	
	// 验证增量响应
	var fullContent strings.Builder
	for _, resp := range responses {
		fullContent.WriteString(resp.Text)
	}
	
	assert.Contains(t, fullContent.String(), "Hello")
	assert.Contains(t, fullContent.String(), "world")
}

// TestClaudeAdapterStream 测试Claude适配器流式响应
func TestClaudeAdapterStream(t *testing.T) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.NotEmpty(t, r.Header.Get("x-api-key"))
		assert.Equal(t, "text/event-stream", r.Header.Get("Accept"))

		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 发送模拟SSE数据
		fmt.Fprint(w, "data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\"Hello\"}}\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\" Claude\"}}\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建Claude适配器
	claudeAdapter := adapter.NewClaudeAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Hello",
		Model:       "claude-3-haiku-20240307",
		MaxTokens:   100,
		Temperature: 0.7,
		Stream:      true,
	}

	// 执行流式请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	responseChan, err := claudeAdapter.TextGenerationStream(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, responseChan)

	// 收集响应
	var responses []*adapter.TextGenerationResponse
	for response := range responseChan {
		responses = append(responses, response)
	}

	// 验证响应
	assert.Greater(t, len(responses), 0)
	
	// 验证增量响应
	var fullContent strings.Builder
	for _, resp := range responses {
		fullContent.WriteString(resp.Text)
	}
	
	assert.Contains(t, fullContent.String(), "Hello")
	assert.Contains(t, fullContent.String(), "Claude")
}

// TestGeminiAdapterStream 测试Gemini适配器流式响应
func TestGeminiAdapterStream(t *testing.T) {
	// 创建模拟HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "text/event-stream", r.Header.Get("Accept"))
		assert.Contains(t, r.URL.RawQuery, "key=")
		assert.Contains(t, r.URL.RawQuery, "alt=sse")

		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 发送模拟SSE数据
		fmt.Fprint(w, "data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"Hello\"}]}}]}\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\" Gemini\"}]}}]}\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	// 创建Gemini适配器
	geminiAdapter := adapter.NewGeminiAdapter("test-key", server.URL)

	// 创建请求
	req := &adapter.TextGenerationRequest{
		Prompt:      "Hello",
		Model:       "gemini-pro",
		MaxTokens:   100,
		Temperature: 0.7,
		Stream:      true,
	}

	// 执行流式请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	responseChan, err := geminiAdapter.TextGenerationStream(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, responseChan)

	// 收集响应
	var responses []*adapter.TextGenerationResponse
	for response := range responseChan {
		responses = append(responses, response)
	}

	// 验证响应
	assert.Greater(t, len(responses), 0)
	
	// 验证增量响应
	var fullContent strings.Builder
	for _, resp := range responses {
		fullContent.WriteString(resp.Text)
	}
	
	assert.Contains(t, fullContent.String(), "Hello")
	assert.Contains(t, fullContent.String(), "Gemini")
}

// TestStreamErrorHandling 测试流式响应错误处理
func TestStreamErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建模拟服务，返回错误
	mockService := new(MockChatService)
	mockService.On("StartChatStream", mock.Anything, mock.AnythingOfType("*ai.ChatRequest")).Return(
		(<-chan *aiService.StreamChatResponse)(nil), 
		fmt.Errorf("服务不可用"),
	)

	// 创建API实例
	chatAPI := ai.NewChatApi(mockService)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/ai/chat/stream", chatAPI.ContinueChat)

	// 准备请求数据
	requestData := aiService.ChatRequest{
		ProjectID: "test-project",
		Message:   "Hello",
		Options: &aiModels.GenerateOptions{
			Stream: true,
		},
	}
	requestBody, _ := json.Marshal(requestData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat/stream", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := NewCloseNotifierResponseRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))

	// 验证错误响应内容
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "event: error")
	assert.Contains(t, responseBody, "服务不可用")

	// 验证模拟调用
	mockService.AssertExpectations(t)
}