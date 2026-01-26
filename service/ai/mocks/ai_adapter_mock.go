package mocks

import (
	"context"
	"time"

	"Qingyu_backend/service/ai/adapter"
)

// MockAIAdapter 模拟AI适配器，用于单元测试
type MockAIAdapter struct {
	// 配置字段
	Name             string
	SupportedModels  []string
	ShouldFail       bool
	FailureError     error
	ShouldTimeout    bool
	ResponseDelay    time.Duration
	TextResponse     *adapter.TextGenerationResponse
	ChatResponse     *adapter.ChatCompletionResponse
	ImageResponse    *adapter.ImageGenerationResponse
	HealthCheckError error
	CallCount        int
	LastRequest      interface{}
}

// NewMockAIAdapter 创建新的模拟适配器
func NewMockAIAdapter(name string) *MockAIAdapter {
	return &MockAIAdapter{
		Name:            name,
		SupportedModels: []string{"mock-model", "gpt-4", "claude-3"},
		ShouldFail:      false,
		ShouldTimeout:   false,
		ResponseDelay:   0,
	}
}

// GetName 获取适配器名称
func (m *MockAIAdapter) GetName() string {
	return m.Name
}

// GetSupportedModels 获取支持的模型列表
func (m *MockAIAdapter) GetSupportedModels() []string {
	return m.SupportedModels
}

// TextGeneration 文本生成
func (m *MockAIAdapter) TextGeneration(ctx context.Context, req *adapter.TextGenerationRequest) (*adapter.TextGenerationResponse, error) {
	m.CallCount++
	m.LastRequest = req

	// 模拟延迟
	if m.ResponseDelay > 0 {
		select {
		case <-time.After(m.ResponseDelay):
			// 延迟完成
		case <-ctx.Done():
			return nil, &adapter.AdapterError{
				Type:       adapter.ErrorTypeTimeout,
				Message:    "请求超时",
				Code:       "timeout",
				StatusCode: 408,
				Provider:   m.Name,
				Retryable:  true,
			}
		}
	}

	// 模拟失败
	if m.ShouldFail {
		if m.FailureError != nil {
			return nil, m.FailureError
		}
		return nil, &adapter.AdapterError{
			Type:       adapter.ErrorTypeServiceUnavailable,
			Message:    "模拟的服务错误",
			Code:       "service_error",
			StatusCode: 500,
			Provider:   m.Name,
			Retryable:  true,
		}
	}

	// 返回预设响应或默认响应
	if m.TextResponse != nil {
		return m.TextResponse, nil
	}

	// 默认响应
	return &adapter.TextGenerationResponse{
		ID:   "mock-response-id",
		Text: "这是模拟的AI生成文本",
		Usage: adapter.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
		Model:        "mock-model",
		FinishReason: "stop",
		CreatedAt:    time.Now(),
	}, nil
}

// ChatCompletion 对话完成
func (m *MockAIAdapter) ChatCompletion(ctx context.Context, req *adapter.ChatCompletionRequest) (*adapter.ChatCompletionResponse, error) {
	m.CallCount++
	m.LastRequest = req

	// 模拟延迟
	if m.ResponseDelay > 0 {
		select {
		case <-time.After(m.ResponseDelay):
			// 延迟完成
		case <-ctx.Done():
			return nil, &adapter.AdapterError{
				Type:       adapter.ErrorTypeTimeout,
				Message:    "请求超时",
				Code:       "timeout",
				StatusCode: 408,
				Provider:   m.Name,
				Retryable:  true,
			}
		}
	}

	// 模拟失败
	if m.ShouldFail {
		if m.FailureError != nil {
			return nil, m.FailureError
		}
		return nil, &adapter.AdapterError{
			Type:       adapter.ErrorTypeServiceUnavailable,
			Message:    "模拟的服务错误",
			Code:       "service_error",
			StatusCode: 500,
			Provider:   m.Name,
			Retryable:  true,
		}
	}

	// 返回预设响应或默认响应
	if m.ChatResponse != nil {
		return m.ChatResponse, nil
	}

	// 默认响应
	return &adapter.ChatCompletionResponse{
		ID: "mock-chat-id",
		Message: adapter.Message{
			Role:    "assistant",
			Content: "这是模拟的对话响应",
		},
		Usage: adapter.Usage{
			PromptTokens:     15,
			CompletionTokens: 25,
			TotalTokens:      40,
		},
		Model:        "mock-model",
		FinishReason: "stop",
		CreatedAt:    time.Now(),
	}, nil
}

// TextGenerationStream 流式文本生成
func (m *MockAIAdapter) TextGenerationStream(ctx context.Context, req *adapter.TextGenerationRequest) (<-chan *adapter.TextGenerationResponse, error) {
	m.CallCount++
	m.LastRequest = req

	// 模拟失败
	if m.ShouldFail {
		if m.FailureError != nil {
			return nil, m.FailureError
		}
		return nil, &adapter.AdapterError{
			Type:       adapter.ErrorTypeServiceUnavailable,
			Message:    "模拟的服务错误",
			Code:       "service_error",
			StatusCode: 500,
			Provider:   m.Name,
			Retryable:  true,
		}
	}

	// 创建响应通道
	respChan := make(chan *adapter.TextGenerationResponse, 3)

	// 发送多个流式响应
	go func() {
		defer close(respChan)
		for i := 0; i < 3; i++ {
			select {
			case <-ctx.Done():
				return
			case respChan <- &adapter.TextGenerationResponse{
				ID:   "mock-stream-id",
				Text: "流式文本片段",
				Usage: adapter.Usage{
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalTokens:      15,
				},
				Model:        "mock-model",
				FinishReason: "",
				CreatedAt:    time.Now(),
			}:
			}
		}
	}()

	return respChan, nil
}

// ImageGeneration 图像生成
func (m *MockAIAdapter) ImageGeneration(ctx context.Context, req *adapter.ImageGenerationRequest) (*adapter.ImageGenerationResponse, error) {
	m.CallCount++
	m.LastRequest = req

	// 模拟延迟
	if m.ResponseDelay > 0 {
		select {
		case <-time.After(m.ResponseDelay):
			// 延迟完成
		case <-ctx.Done():
			return nil, &adapter.AdapterError{
				Type:       adapter.ErrorTypeTimeout,
				Message:    "请求超时",
				Code:       "timeout",
				StatusCode: 408,
				Provider:   m.Name,
				Retryable:  true,
			}
		}
	}

	// 模拟失败
	if m.ShouldFail {
		if m.FailureError != nil {
			return nil, m.FailureError
		}
		return nil, &adapter.AdapterError{
			Type:       adapter.ErrorTypeServiceUnavailable,
			Message:    "模拟的服务错误",
			Code:       "service_error",
			StatusCode: 500,
			Provider:   m.Name,
			Retryable:  true,
		}
	}

	// 返回预设响应或默认响应
	if m.ImageResponse != nil {
		return m.ImageResponse, nil
	}

	// 默认响应
	return &adapter.ImageGenerationResponse{
		ID: "mock-image-id",
		Images: []adapter.ImageData{
			{
				URL:           "https://example.com/mock-image.png",
				Base64:        "",
				RevisedPrompt: "修订后的提示词",
			},
		},
		Usage: adapter.Usage{
			PromptTokens:     20,
			CompletionTokens: 0,
			TotalTokens:      20,
		},
		Model:     "mock-image-model",
		CreatedAt: time.Now(),
	}, nil
}

// HealthCheck 健康检查
func (m *MockAIAdapter) HealthCheck(ctx context.Context) error {
	if m.HealthCheckError != nil {
		return m.HealthCheckError
	}
	return nil
}

// Reset 重置模拟器状态
func (m *MockAIAdapter) Reset() {
	m.CallCount = 0
	m.LastRequest = nil
	m.ShouldFail = false
	m.FailureError = nil
	m.ShouldTimeout = false
	m.TextResponse = nil
	m.ChatResponse = nil
	m.ImageResponse = nil
	m.HealthCheckError = nil
}

// SetTextResponse 设置文本生成响应
func (m *MockAIAdapter) SetTextResponse(text string, tokens int) {
	m.TextResponse = &adapter.TextGenerationResponse{
		ID:   "mock-response-id",
		Text: text,
		Usage: adapter.Usage{
			PromptTokens:     tokens / 2,
			CompletionTokens: tokens - tokens/2,
			TotalTokens:      tokens,
		},
		Model:        "mock-model",
		FinishReason: "stop",
		CreatedAt:    time.Now(),
	}
}

// SetChatResponse 设置对话响应
func (m *MockAIAdapter) SetChatResponse(content string, tokens int) {
	m.ChatResponse = &adapter.ChatCompletionResponse{
		ID: "mock-chat-id",
		Message: adapter.Message{
			Role:    "assistant",
			Content: content,
		},
		Usage: adapter.Usage{
			PromptTokens:     tokens / 2,
			CompletionTokens: tokens - tokens/2,
			TotalTokens:      tokens,
		},
		Model:        "mock-model",
		FinishReason: "stop",
		CreatedAt:    time.Now(),
	}
}
