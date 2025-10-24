package ai

import (
	"context"
	"time"
)

// ModelAdapter 模型适配器接口
type ModelAdapter interface {
	// Initialize 初始化适配器
	Initialize(ctx context.Context) error

	// Health 健康检查
	Health(ctx context.Context) error

	// Close 关闭适配器
	Close(ctx context.Context) error

	// GetModelID 获取模型ID
	GetModelID() string

	// GetProvider 获取提供商
	GetProvider() string

	// GetConfig 获取配置
	GetConfig() *ModelConfig

	// GenerateContent 生成内容
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)

	// GenerateContentStream 流式生成内容
	GenerateContentStream(ctx context.Context, req *GenerateContentRequest) (<-chan *StreamResponse, error)
}

// ModelConfig 模型配置
type ModelConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Type        string            `json:"type"`
	MaxTokens   int               `json:"max_tokens"`
	Temperature float64           `json:"temperature"`
	TopP        float64           `json:"top_p"`
	TopK        int               `json:"top_k"`
	InputPrice  float64           `json:"input_price"`
	OutputPrice float64           `json:"output_price"`
	Features    []string          `json:"features"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// GetModelConfigRequest 获取模型配置请求
type GetModelConfigRequest struct {
	ModelID string `json:"model_id" validate:"required"`
}

// GetModelConfigResponse 获取模型配置响应
type GetModelConfigResponse struct {
	Config *ModelConfig `json:"config"`
}

// UpdateModelConfigRequest 更新模型配置请求
type UpdateModelConfigRequest struct {
	ModelID     string            `json:"model_id" validate:"required"`
	Name        string            `json:"name,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	TopP        float64           `json:"top_p,omitempty"`
	TopK        int               `json:"top_k,omitempty"`
	Features    []string          `json:"features,omitempty"`
	Status      string            `json:"status,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateModelConfigResponse 更新模型配置响应
type UpdateModelConfigResponse struct {
	Updated   bool      `json:"updated"`
	UpdatedAt time.Time `json:"updated_at"`
}
