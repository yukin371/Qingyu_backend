package adapter

import (
	"context"
)

// DeepSeekAdapter DeepSeek适配器实现（基于OpenAI兼容API）
type DeepSeekAdapter struct {
	*OpenAIAdapter // 继承OpenAI适配器的所有功能
}

// NewDeepSeekAdapter 创建DeepSeek适配器实例
func NewDeepSeekAdapter(apiKey, baseURL string) *DeepSeekAdapter {
	if baseURL == "" {
		// DeepSeek的文本补全API需要使用/beta endpoint
		baseURL = "https://api.deepseek.com/beta"
	}

	// 创建OpenAI适配器作为基础
	openaiAdapter := NewOpenAIAdapter(apiKey, baseURL)

	return &DeepSeekAdapter{
		OpenAIAdapter: openaiAdapter,
	}
}

// GetName 获取适配器名称
func (a *DeepSeekAdapter) GetName() string {
	return "deepseek"
}

// GetSupportedModels 获取支持的模型列表
func (a *DeepSeekAdapter) GetSupportedModels() []string {
	return []string{
		"deepseek-chat",     // DeepSeek-V3 标准对话模型
		"deepseek-reasoner", // DeepSeek-V3 推理模型 (思考链)
	}
}

// HealthCheck 健康检查
func (a *DeepSeekAdapter) HealthCheck(ctx context.Context) error {
	// 复用OpenAI适配器的健康检查
	return a.OpenAIAdapter.HealthCheck(ctx)
}

// TextGeneration 文本生成（复用OpenAI适配器）
func (a *DeepSeekAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 如果没有指定模型，使用默认模型
	if req.Model == "" {
		req.Model = "deepseek-chat"
	}
	return a.OpenAIAdapter.TextGeneration(ctx, req)
}

// ChatCompletion 对话完成（复用OpenAI适配器）
func (a *DeepSeekAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 如果没有指定模型，使用默认模型
	if req.Model == "" {
		req.Model = "deepseek-chat"
	}
	return a.OpenAIAdapter.ChatCompletion(ctx, req)
}

// TextGenerationStream 流式文本生成（复用OpenAI适配器）
func (a *DeepSeekAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	// 如果没有指定模型，使用默认模型
	if req.Model == "" {
		req.Model = "deepseek-chat"
	}
	return a.OpenAIAdapter.TextGenerationStream(ctx, req)
}

// ImageGeneration 图像生成（DeepSeek不支持图像生成）
func (a *DeepSeekAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	return nil, &AdapterError{
		Code:    ErrorTypeNotImplemented,
		Message: "DeepSeek不支持图像生成功能",
		Type:    ErrorTypeNotImplemented,
	}
}
