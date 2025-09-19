package adapter

import (
	"context"
)

// AIAdapter 适配器接口，定义所有AI服务提供商需要实现的方法
type AIAdapter interface {
	// 文本生成相关方法
	TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error)

	// 图像生成相关方法
	ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error)
}

// TextGenerationRequest 文本生成请求结构
type TextGenerationRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

// TextGenerationResponse 文本生成响应结构
type TextGenerationResponse struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Usage Usage  `json:"usage"`
	Model string `json:"model"`
}

// Usage 使用情况统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ImageGenerationRequest 图像生成请求结构
type ImageGenerationRequest struct {
	// TODO: 定义图像生成请求结构
}

// ImageGenerationResponse 图像生成响应结构
type ImageGenerationResponse struct {
	// TODO: 定义图像生成响应结构
}
