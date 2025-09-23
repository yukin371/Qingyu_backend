package adapter

import (
	"context"
	"fmt"
	"time"
)

// AIAdapter 适配器接口，定义所有AI服务提供商需要实现的方法
type AIAdapter interface {
	// 获取适配器名称
	GetName() string
	
	// 获取支持的模型列表
	GetSupportedModels() []string
	
	// 文本生成相关方法
	TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error)
	
	// 对话生成方法（支持多轮对话）
	ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	
	// 流式文本生成
	TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error)
	
	// 图像生成相关方法
	ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error)
	
	// 健康检查
	HealthCheck(ctx context.Context) error
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // 消息内容
}

// TextGenerationRequest 文本生成请求结构
type TextGenerationRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
	Stop        []string `json:"stop,omitempty"`
	User        string  `json:"user,omitempty"`
}

// ChatCompletionRequest 对话完成请求结构
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Stop        []string  `json:"stop,omitempty"`
	User        string    `json:"user,omitempty"`
}

// TextGenerationResponse 文本生成响应结构
type TextGenerationResponse struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	Usage        Usage     `json:"usage"`
	Model        string    `json:"model"`
	FinishReason string    `json:"finish_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ChatCompletionResponse 对话完成响应结构
type ChatCompletionResponse struct {
	ID           string    `json:"id"`
	Message      Message   `json:"message"`
	Usage        Usage     `json:"usage"`
	Model        string    `json:"model"`
	FinishReason string    `json:"finish_reason"`
	CreatedAt    time.Time `json:"created_at"`
}

// Choice 选择结构
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用情况统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ImageGenerationRequest 图像生成请求结构
type ImageGenerationRequest struct {
	Model  string `json:"model,omitempty"`
	Prompt string `json:"prompt"`
	Size   string `json:"size,omitempty"`   // 图像尺寸，如 "1024x1024"
	Style  string `json:"style,omitempty"`  // 图像风格
	N      int    `json:"n,omitempty"`      // 生成图像数量
}

// ImageGenerationResponse 图像生成响应结构
type ImageGenerationResponse struct {
	ID        string      `json:"id"`
	Images    []ImageData `json:"images"`
	Usage     Usage       `json:"usage"`
	Model     string      `json:"model"`
	CreatedAt time.Time   `json:"created_at"`
}

// ImageData 图像数据结构
type ImageData struct {
	URL           string `json:"url,omitempty"`            // 图像URL
	Base64        string `json:"base64,omitempty"`         // Base64编码的图像数据
	RevisedPrompt string `json:"revised_prompt,omitempty"` // 修订后的提示词
}

// AdapterError 适配器错误
type AdapterError struct {
	Type       string `json:"type"`       // 错误类型
	Message    string `json:"message"`    // 错误消息
	Code       string `json:"code"`       // 错误代码
	StatusCode int    `json:"statusCode"` // HTTP状态码
	Provider   string `json:"provider"`   // 提供商名称
	Retryable  bool   `json:"retryable"`  // 是否可重试
}

func (e *AdapterError) Error() string {
	return fmt.Sprintf("[%s:%s] %s (code: %s, status: %d, retryable: %t)", 
		e.Provider, e.Type, e.Message, e.Code, e.StatusCode, e.Retryable)
}

// NewAdapterError 创建适配器错误
func NewAdapterError(provider, errorType, message, code string, statusCode int, retryable bool) *AdapterError {
	return &AdapterError{
		Type:       errorType,
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
		Provider:   provider,
		Retryable:  retryable,
	}
}

// IsRetryable 检查错误是否可重试
func (e *AdapterError) IsRetryable() bool {
	return e.Retryable
}

// IsRateLimitError 检查是否为速率限制错误
func (e *AdapterError) IsRateLimitError() bool {
	return e.Type == ErrorTypeRateLimit
}

// IsAuthError 检查是否为认证错误
func (e *AdapterError) IsAuthError() bool {
	return e.Type == ErrorTypeAuthentication
}

// 常见错误类型
const (
	ErrorTypeInvalidRequest     = "invalid_request"
	ErrorTypeAuthentication     = "authentication_error"
	ErrorTypePermission         = "permission_error"
	ErrorTypeRateLimit          = "rate_limit_exceeded"
	ErrorTypeQuotaExceeded      = "quota_exceeded"
	ErrorTypeModelNotFound      = "model_not_found"
	ErrorTypeServiceUnavailable = "service_unavailable"
	ErrorTypeTimeout            = "timeout"
	ErrorTypeNetworkError       = "network_error"
	ErrorTypeInvalidResponse    = "invalid_response"
	ErrorTypeNotImplemented     = "not_implemented"
	ErrorTypeUnknown            = "unknown_error"
)
