package adapter

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAIAdapter 模拟AI适配器
type MockAIAdapter struct {
	mock.Mock
}

func (m *MockAIAdapter) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAIAdapter) GetSupportedModels() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockAIAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*TextGenerationResponse), args.Error(1)
}

func (m *MockAIAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*ChatCompletionResponse), args.Error(1)
}

func (m *MockAIAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan *TextGenerationResponse), args.Error(1)
}

func (m *MockAIAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*ImageGenerationResponse), args.Error(1)
}

func (m *MockAIAdapter) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestAdapterManager_GetAdapter 测试获取适配器
func TestAdapterManager_GetAdapter(t *testing.T) {
	// 创建测试配置
	cfg := &config.ExternalAPIConfig{
		Providers: map[string]*config.ProviderConfig{
			"openai": {
				Name:     "openai",
				Priority: 1,
				Enabled:  true,
			},
		},
	}

	// 创建适配器管理器
	manager := &AdapterManager{
		adapters: make(map[string]AIAdapter),
		config:   cfg,
	}

	// 添加测试适配器
	mockAdapter := &MockAIAdapter{}
	manager.adapters["openai"] = mockAdapter

	// 测试获取存在的适配器
	adapter, err := manager.GetAdapter("openai")
	assert.NoError(t, err)
	assert.Equal(t, mockAdapter, adapter)

	// 测试获取不存在的适配器
	_, err = manager.GetAdapter("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "适配器 'non-existent' 不存在或未启用")
}

// TestAdapterManager_TextGeneration 测试文本生成
func TestAdapterManager_TextGeneration(t *testing.T) {
	// 创建模拟适配器
	mockAdapter := new(MockAIAdapter)

	// 创建测试请求和响应
	testReq := &TextGenerationRequest{
		Model:       "test-model",
		Prompt:      "测试提示词",
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	testResp := &TextGenerationResponse{
		ID:   "test-response-id",
		Text: "生成的测试文本",
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
		Model:     "test-model",
		CreatedAt: time.Now(),
	}

	// 设置模拟期望
	mockAdapter.On("TextGeneration", mock.Anything, testReq).Return(testResp, nil)

	// 创建适配器管理器
	manager := &AdapterManager{
		adapters: make(map[string]AIAdapter),
	}
	manager.AddAdapter("test-adapter", mockAdapter)

	// 执行测试
	ctx := context.Background()
	response, err := manager.TextGeneration(ctx, "test-adapter", testReq)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test-response-id", response.ID)
	assert.Equal(t, "生成的测试文本", response.Text)
	assert.Equal(t, "test-model", response.Model)
	assert.Equal(t, 30, response.Usage.TotalTokens)

	// 验证模拟调用
	mockAdapter.AssertExpectations(t)
}

// TestAdapterManager_ChatCompletion 测试对话完成
func TestAdapterManager_ChatCompletion(t *testing.T) {
	// 创建模拟适配器
	mockAdapter := new(MockAIAdapter)

	// 创建测试请求和响应
	testReq := &ChatCompletionRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "你好"},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	testResp := &ChatCompletionResponse{
		ID: "test-chat-response-id",
		Message: Message{
			Role:    "assistant",
			Content: "你好！有什么可以帮助你的吗？",
		},
		Usage: Usage{
			PromptTokens:     5,
			CompletionTokens: 15,
			TotalTokens:      20,
		},
		Model:     "test-model",
		CreatedAt: time.Now(),
	}

	// 设置模拟期望
	mockAdapter.On("ChatCompletion", mock.Anything, testReq).Return(testResp, nil)

	// 创建适配器管理器
	manager := &AdapterManager{
		adapters: make(map[string]AIAdapter),
	}
	manager.AddAdapter("test-adapter", mockAdapter)

	// 执行测试
	ctx := context.Background()
	response, err := manager.ChatCompletion(ctx, "test-adapter", testReq)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test-chat-response-id", response.ID)
	assert.Equal(t, "assistant", response.Message.Role)
	assert.Equal(t, "你好！有什么可以帮助你的吗？", response.Message.Content)
	assert.Equal(t, "test-model", response.Model)
	assert.Equal(t, 20, response.Usage.TotalTokens)

	// 验证模拟调用
	mockAdapter.AssertExpectations(t)
}

// TestAdapterManager_HealthCheck 测试健康检查
func TestAdapterManager_HealthCheck(t *testing.T) {
	// 创建模拟适配器
	mockAdapter1 := new(MockAIAdapter)
	mockAdapter2 := new(MockAIAdapter)

	// 设置模拟期望
	mockAdapter1.On("HealthCheck", mock.Anything).Return(nil)
	mockAdapter2.On("HealthCheck", mock.Anything).Return(assert.AnError)

	// 创建适配器管理器
	manager := &AdapterManager{
		adapters: make(map[string]AIAdapter),
	}
	manager.AddAdapter("healthy-adapter", mockAdapter1)
	manager.AddAdapter("unhealthy-adapter", mockAdapter2)

	// 执行测试
	ctx := context.Background()
	results := manager.HealthCheck(ctx)

	// 验证结果
	assert.Len(t, results, 2)
	assert.NoError(t, results["healthy-adapter"])
	assert.Error(t, results["unhealthy-adapter"])

	// 验证模拟调用
	mockAdapter1.AssertExpectations(t)
	mockAdapter2.AssertExpectations(t)
}

// TestAdapterError 测试适配器错误
func TestAdapterError(t *testing.T) {
	// 创建适配器错误
	err := NewAdapterError(
		"test-provider",
		ErrorTypeRateLimit,
		"Rate limit exceeded",
		"rate_limit_429",
		429,
		true,
	)

	// 验证错误属性
	assert.Equal(t, "test-provider", err.Provider)
	assert.Equal(t, ErrorTypeRateLimit, err.Type)
	assert.Equal(t, "Rate limit exceeded", err.Message)
	assert.Equal(t, "rate_limit_429", err.Code)
	assert.Equal(t, 429, err.StatusCode)
	assert.True(t, err.Retryable)

	// 验证错误方法
	assert.True(t, err.IsRetryable())
	assert.True(t, err.IsRateLimitError())
	assert.False(t, err.IsAuthError())

	// 验证错误字符串
	expectedErrorString := "[test-provider:rate_limit_exceeded] Rate limit exceeded (code: rate_limit_429, status: 429, retryable: true)"
	assert.Equal(t, expectedErrorString, err.Error())
}

// TestAdapterError_AuthError 测试认证错误
func TestAdapterError_AuthError(t *testing.T) {
	err := NewAdapterError(
		"test-provider",
		ErrorTypeAuthentication,
		"Invalid API key",
		"auth_401",
		401,
		false,
	)

	assert.True(t, err.IsAuthError())
	assert.False(t, err.IsRateLimitError())
	assert.False(t, err.IsRetryable())
}

// TestGetAvailableAdapters 测试获取可用适配器列表
func TestGetAvailableAdapters(t *testing.T) {
	// 创建测试配置
	cfg := &config.ExternalAPIConfig{
		Providers: map[string]*config.ProviderConfig{
			"openai": {
				Name:     "openai",
				Priority: 1,
				Enabled:  true,
			},
			"claude": {
				Name:     "claude",
				Priority: 2,
				Enabled:  true,
			},
		},
	}

	// 创建适配器管理器
	manager := &AdapterManager{
		adapters: make(map[string]AIAdapter),
		config:   cfg,
	}

	// 添加测试适配器
	mockAdapter1 := &MockAIAdapter{}
	mockAdapter2 := &MockAIAdapter{}

	manager.adapters["openai"] = mockAdapter1
	manager.adapters["claude"] = mockAdapter2

	// 获取可用适配器
	adapters := manager.GetAvailableAdapters()

	// 验证结果
	assert.Len(t, adapters, 2)
	assert.Contains(t, adapters, "openai")
	assert.Contains(t, adapters, "claude")

	// 验证排序（按优先级）
	assert.Equal(t, "openai", adapters[0])
	assert.Equal(t, "claude", adapters[1])
}

// TestGetSupportedModels 测试获取支持的模型列表
func TestGetSupportedModels(t *testing.T) {
	// 创建测试配置
	cfg := &config.ExternalAPIConfig{
		Providers: map[string]*config.ProviderConfig{
			"openai": {
				Name:     "openai",
				Priority: 1,
				Enabled:  true,
			},
			"claude": {
				Name:     "claude",
				Priority: 2,
				Enabled:  true,
			},
		},
	}

	// 创建适配器管理器
	manager := &AdapterManager{
		adapters: make(map[string]AIAdapter),
		config:   cfg,
	}

	// 创建模拟适配器
	mockAdapter1 := &MockAIAdapter{}
	mockAdapter2 := &MockAIAdapter{}

	// 设置模拟期望
	mockAdapter1.On("GetSupportedModels").Return([]string{"gpt-3.5-turbo", "gpt-4"})
	mockAdapter2.On("GetSupportedModels").Return([]string{"claude-3-sonnet", "claude-3-opus"})

	manager.adapters["openai"] = mockAdapter1
	manager.adapters["claude"] = mockAdapter2

	// 获取支持的模型
	models := manager.GetSupportedModels()

	// 验证结果
	assert.Len(t, models, 2)
	assert.Contains(t, models, "openai")
	assert.Contains(t, models, "claude")
	assert.Equal(t, []string{"gpt-3.5-turbo", "gpt-4"}, models["openai"])
	assert.Equal(t, []string{"claude-3-sonnet", "claude-3-opus"}, models["claude"])

	// 验证模拟调用
	mockAdapter1.AssertExpectations(t)
	mockAdapter2.AssertExpectations(t)
}
