package adapter

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"Qingyu_backend/config"
	"github.com/sirupsen/logrus"
)

// Deprecated: Use Qingyu-Ai-Service gRPC API instead.
// This adapter manager is kept only for emergency fallback.
// Will be removed in v2.0.0
//
// Migration guide:
// 1. Use service/ai/grpc_client.go to call Qingyu-Ai-Service
// 2. Configure AI_SERVICE_ENDPOINT environment variable
// 3. Enable fallback with AI_ENABLE_FALLBACK=true if needed
type AdapterManager struct {
	adapters        map[string]AIAdapter
	config          *config.ExternalAPIConfig
	defaultProvider string
	mu              sync.RWMutex
}

// Deprecated: Use gRPC client instead
func NewAdapterManager(cfg *config.ExternalAPIConfig) *AdapterManager {
	// Log deprecation warning
	logrus.Warn("AdapterManager is deprecated. Use Qingyu-Ai-Service gRPC API instead. " +
		"This is kept only for emergency fallback and will be removed in v2.0.0. " +
		"See docs/plans/2026-01-24-ai-service-complete-migration-design.md for migration guide.")

	manager := &AdapterManager{
		adapters:        make(map[string]AIAdapter),
		config:          cfg,
		defaultProvider: cfg.DefaultProvider,
	}

	// 初始化所有启用的适配器
	manager.initializeAdapters()

	return manager
}

// initializeAdapters 初始化所有启用的适配器
func (m *AdapterManager) initializeAdapters() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, providerConfig := range m.config.Providers {
		if !providerConfig.Enabled {
			continue
		}

		var adapter AIAdapter
		switch name {
		case "openai":
			adapter = NewOpenAIAdapter(providerConfig.APIKey, providerConfig.BaseURL)
		case "deepseek":
			// DeepSeek API兼容OpenAI格式，使用专用DeepSeek适配器
			adapter = NewDeepSeekAdapter(providerConfig.APIKey, providerConfig.BaseURL)
		case "claude":
			adapter = NewClaudeAdapter(providerConfig.APIKey, providerConfig.BaseURL)
		case "gemini":
			adapter = NewGeminiAdapter(providerConfig.APIKey, providerConfig.BaseURL)
		case "wenxin":
			adapter = NewWenxinAdapter(providerConfig.APIKey, providerConfig.SecretKey, providerConfig.BaseURL)
		case "qwen":
			adapter = NewQwenAdapter(providerConfig.APIKey, providerConfig.BaseURL)
		case "glm", "zhipu":
			adapter = NewGLMAdapter(providerConfig.APIKey, providerConfig.BaseURL)
		default:
			continue
		}

		if adapter != nil {
			m.adapters[name] = adapter
		}
	}
}

// GetAdapter 获取指定的适配器
func (m *AdapterManager) GetAdapter(provider string) (AIAdapter, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if provider == "" {
		provider = m.defaultProvider
		// 检查默认提供商是否为空
		if provider == "" {
			return nil, &AdapterError{
				Code:    ErrorTypeInvalidRequest,
				Message: "未配置默认AI提供商，请检查配置文件",
				Type:    ErrorTypeInvalidRequest,
			}
		}
	}

	adapter, exists := m.adapters[provider]
	if !exists {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: fmt.Sprintf("适配器 '%s' 不存在或未启用", provider),
			Type:    ErrorTypeInvalidRequest,
		}
	}

	return adapter, nil
}

// GetDefaultAdapter 获取默认适配器
func (m *AdapterManager) GetDefaultAdapter() (AIAdapter, error) {
	if m.defaultProvider == "" {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: "未配置默认AI提供商，请检查配置文件",
			Type:    ErrorTypeInvalidRequest,
		}
	}
	return m.GetAdapter(m.defaultProvider)
}

// GetAvailableAdapters 获取所有可用的适配器
func (m *AdapterManager) GetAvailableAdapters() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var providers []string
	for name := range m.adapters {
		providers = append(providers, name)
	}

	// 按优先级排序
	sort.Slice(providers, func(i, j int) bool {
		configI := m.config.Providers[providers[i]]
		configJ := m.config.Providers[providers[j]]
		return configI.Priority < configJ.Priority
	})

	return providers
}

// GetAdapterByModel 根据模型名称获取适配器
func (m *AdapterManager) GetAdapterByModel(model string) (AIAdapter, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 按优先级排序的提供商列表
	var sortedProviders []string
	for name := range m.adapters {
		sortedProviders = append(sortedProviders, name)
	}

	sort.Slice(sortedProviders, func(i, j int) bool {
		configI := m.config.Providers[sortedProviders[i]]
		configJ := m.config.Providers[sortedProviders[j]]
		return configI.Priority < configJ.Priority
	})

	// 查找支持该模型的适配器
	for _, provider := range sortedProviders {
		config := m.config.Providers[provider]
		for _, supportedModel := range config.SupportedModels {
			if supportedModel == model {
				return m.adapters[provider], nil
			}
		}
	}

	return nil, &AdapterError{
		Code:    ErrorTypeInvalidRequest,
		Message: fmt.Sprintf("没有找到支持模型 '%s' 的适配器", model),
		Type:    ErrorTypeInvalidRequest,
	}
}

// HealthCheck 检查所有适配器的健康状态
func (m *AdapterManager) HealthCheck(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]error)

	for name, adapter := range m.adapters {
		err := adapter.HealthCheck(ctx)
		results[name] = err
	}

	return results
}

// GetSupportedModels 获取所有支持的模型
func (m *AdapterManager) GetSupportedModels() map[string][]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	models := make(map[string][]string)

	for name, adapter := range m.adapters {
		models[name] = adapter.GetSupportedModels()
	}

	return models
}

// ReloadConfig 重新加载配置
func (m *AdapterManager) ReloadConfig(cfg *config.ExternalAPIConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = cfg
	m.defaultProvider = cfg.DefaultProvider

	// 清空现有适配器
	m.adapters = make(map[string]AIAdapter)

	// 重新初始化适配器
	m.initializeAdapters()
}

// AddAdapter 动态添加适配器
func (m *AdapterManager) AddAdapter(name string, adapter AIAdapter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.adapters[name] = adapter
}

// RemoveAdapter 移除适配器
func (m *AdapterManager) RemoveAdapter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.adapters, name)
}

// SetDefaultProvider 设置默认提供商
func (m *AdapterManager) SetDefaultProvider(provider string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.adapters[provider]; !exists {
		return &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: fmt.Sprintf("适配器 '%s' 不存在", provider),
			Type:    ErrorTypeInvalidRequest,
		}
	}

	m.defaultProvider = provider
	return nil
}

// GetProviderConfig 获取提供商配置
func (m *AdapterManager) GetProviderConfig(provider string) (*config.ProviderConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, exists := m.config.Providers[provider]
	if !exists {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: fmt.Sprintf("提供商 '%s' 配置不存在", provider),
			Type:    ErrorTypeInvalidRequest,
		}
	}

	return config, nil
}

// TextGeneration 使用指定提供商进行文本生成
func (m *AdapterManager) TextGeneration(ctx context.Context, provider string, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	adapter, err := m.GetAdapter(provider)
	if err != nil {
		return nil, err
	}

	return adapter.TextGeneration(ctx, req)
}

// ChatCompletion 使用指定提供商进行对话完成
func (m *AdapterManager) ChatCompletion(ctx context.Context, provider string, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	adapter, err := m.GetAdapter(provider)
	if err != nil {
		return nil, err
	}

	return adapter.ChatCompletion(ctx, req)
}

// AutoTextGenerationStream 自动选择最佳适配器进行流式文本生成
func (m *AdapterManager) AutoTextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	// 如果指定了模型，尝试根据模型选择适配器
	if req.Model != "" {
		adapter, err := m.GetAdapterByModel(req.Model)
		if err == nil {
			return adapter.TextGenerationStream(ctx, req)
		}
	}

	// 否则使用默认适配器
	adapter, err := m.GetDefaultAdapter()
	if err != nil {
		return nil, err
	}

	return adapter.TextGenerationStream(ctx, req)
}

// TextGenerationStream 使用指定提供商进行流式文本生成
func (m *AdapterManager) TextGenerationStream(ctx context.Context, provider string, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	adapter, err := m.GetAdapter(provider)
	if err != nil {
		return nil, err
	}

	return adapter.TextGenerationStream(ctx, req)
}

// ImageGeneration 使用指定提供商进行图像生成
func (m *AdapterManager) ImageGeneration(ctx context.Context, provider string, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	adapter, err := m.GetAdapter(provider)
	if err != nil {
		return nil, err
	}

	return adapter.ImageGeneration(ctx, req)
}

// AutoTextGeneration 自动选择最佳适配器进行文本生成
func (m *AdapterManager) AutoTextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 如果指定了模型，尝试根据模型选择适配器
	if req.Model != "" {
		adapter, err := m.GetAdapterByModel(req.Model)
		if err == nil {
			return adapter.TextGeneration(ctx, req)
		}
	}

	// 否则使用默认适配器
	adapter, err := m.GetDefaultAdapter()
	if err != nil {
		return nil, err
	}

	return adapter.TextGeneration(ctx, req)
}

// AutoChatCompletion 自动选择最佳适配器进行对话完成
func (m *AdapterManager) AutoChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 如果指定了模型，尝试根据模型选择适配器
	if req.Model != "" {
		adapter, err := m.GetAdapterByModel(req.Model)
		if err == nil {
			return adapter.ChatCompletion(ctx, req)
		}
	}

	// 否则使用默认适配器
	adapter, err := m.GetDefaultAdapter()
	if err != nil {
		return nil, err
	}

	return adapter.ChatCompletion(ctx, req)
}
