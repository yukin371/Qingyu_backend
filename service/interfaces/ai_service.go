package interfaces

import (
	"context"
	"time"

	"Qingyu_backend/service/base"
)

// AIService AI服务接口
// 定义AI相关的所有服务方法
type AIService interface {
	base.BaseService

	// GenerateContent 生成内容
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)

	// GenerateContentStream 流式生成内容
	GenerateContentStream(ctx context.Context, req *GenerateContentRequest) (<-chan *StreamResponse, error)

	// AnalyzeContent 分析内容
	AnalyzeContent(ctx context.Context, req *AnalyzeContentRequest) (*AnalyzeContentResponse, error)

	// ContinueWriting 续写内容
	ContinueWriting(ctx context.Context, req *ContinueWritingRequest) (*ContinueWritingResponse, error)

	// OptimizeText 优化文本
	OptimizeText(ctx context.Context, req *OptimizeTextRequest) (*OptimizeTextResponse, error)

	// GenerateOutline 生成大纲
	GenerateOutline(ctx context.Context, req *GenerateOutlineRequest) (*GenerateOutlineResponse, error)

	// GetContextInfo 获取上下文信息
	GetContextInfo(ctx context.Context, req *GetContextInfoRequest) (*GetContextInfoResponse, error)

	// UpdateContextWithFeedback 根据反馈更新上下文
	UpdateContextWithFeedback(ctx context.Context, req *UpdateContextWithFeedbackRequest) (*UpdateContextWithFeedbackResponse, error)

	// GetSupportedModels 获取支持的AI模型列表
	GetSupportedModels(ctx context.Context) (*GetSupportedModelsResponse, error)

	// GetModelInfo 获取模型信息
	GetModelInfo(ctx context.Context, req *GetModelInfoRequest) (*GetModelInfoResponse, error)

	// ValidateAPIKey 验证API密钥
	ValidateAPIKey(ctx context.Context, req *ValidateAPIKeyRequest) (*ValidateAPIKeyResponse, error)
}

// ContextService 上下文服务接口
type ContextService interface {
	base.BaseService

	// CreateContext 创建上下文
	CreateContext(ctx context.Context, req *CreateContextRequest) (*CreateContextResponse, error)

	// GetContext 获取上下文
	GetContext(ctx context.Context, req *GetContextRequest) (*GetContextResponse, error)

	// UpdateContext 更新上下文
	UpdateContext(ctx context.Context, req *UpdateContextRequest) (*UpdateContextResponse, error)

	// DeleteContext 删除上下文
	DeleteContext(ctx context.Context, req *DeleteContextRequest) (*DeleteContextResponse, error)

	// ListContexts 列出上下文
	ListContexts(ctx context.Context, req *ListContextsRequest) (*ListContextsResponse, error)

	// AddMessage 添加消息到上下文
	AddMessage(ctx context.Context, req *AddMessageRequest) (*AddMessageResponse, error)

	// GetMessages 获取上下文消息
	GetMessages(ctx context.Context, req *GetMessagesRequest) (*GetMessagesResponse, error)

	// ClearMessages 清空上下文消息
	ClearMessages(ctx context.Context, req *ClearMessagesRequest) (*ClearMessagesResponse, error)
}

// ExternalAPIService 外部API服务接口
type ExternalAPIService interface {
	base.BaseService

	// CallAPI 调用外部API
	CallAPI(ctx context.Context, req *CallAPIRequest) (*CallAPIResponse, error)

	// GetAPIStatus 获取API状态
	GetAPIStatus(ctx context.Context, req *GetAPIStatusRequest) (*GetAPIStatusResponse, error)

	// GetAPIUsage 获取API使用情况
	GetAPIUsage(ctx context.Context, req *GetAPIUsageRequest) (*GetAPIUsageResponse, error)

	// RefreshAPIKey 刷新API密钥
	RefreshAPIKey(ctx context.Context, req *RefreshAPIKeyRequest) (*RefreshAPIKeyResponse, error)
}

// AdapterManager 适配器管理器接口
type AdapterManager interface {
	base.BaseService

	// GetAdapter 获取适配器
	GetAdapter(ctx context.Context, req *GetAdapterRequest) (*GetAdapterResponse, error)

	// ListAdapters 列出适配器
	ListAdapters(ctx context.Context, req *ListAdaptersRequest) (*ListAdaptersResponse, error)

	// RegisterAdapter 注册适配器
	RegisterAdapter(ctx context.Context, req *RegisterAdapterRequest) (*RegisterAdapterResponse, error)

	// UnregisterAdapter 注销适配器
	UnregisterAdapter(ctx context.Context, req *UnregisterAdapterRequest) (*UnregisterAdapterResponse, error)

	// UpdateAdapter 更新适配器
	UpdateAdapter(ctx context.Context, req *UpdateAdapterRequest) (*UpdateAdapterResponse, error)

	// GetModelConfig 获取模型配置
	GetModelConfig(ctx context.Context, req *GetModelConfigRequest) (*GetModelConfigResponse, error)

	// UpdateModelConfig 更新模型配置
	UpdateModelConfig(ctx context.Context, req *UpdateModelConfigRequest) (*UpdateModelConfigResponse, error)
}

// 请求和响应结构体定义

// GenerateContentRequest 生成内容请求
type GenerateContentRequest struct {
	Model       string            `json:"model" validate:"required"`
	Prompt      string            `json:"prompt" validate:"required"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	TopP        float64           `json:"top_p,omitempty"`
	TopK        int               `json:"top_k,omitempty"`
	Stop        []string          `json:"stop,omitempty"`
	Stream      bool              `json:"stream,omitempty"`
	Context     map[string]string `json:"context,omitempty"`
	UserID      string            `json:"user_id,omitempty"`
	SessionID   string            `json:"session_id,omitempty"`
}

// GenerateContentResponse 生成内容响应
type GenerateContentResponse struct {
	Content      string            `json:"content"`
	Model        string            `json:"model"`
	TokensUsed   int               `json:"tokens_used"`
	FinishReason string            `json:"finish_reason"`
	ResponseTime time.Duration     `json:"response_time"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	RequestID    string            `json:"request_id"`
}

// StreamResponse 流式响应
type StreamResponse struct {
	Content      string            `json:"content"`
	Delta        string            `json:"delta"`
	Done         bool              `json:"done"`
	TokensUsed   int               `json:"tokens_used"`
	FinishReason string            `json:"finish_reason,omitempty"`
	Error        string            `json:"error,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// AnalyzeContentRequest 分析内容请求
type AnalyzeContentRequest struct {
	Content      string   `json:"content" validate:"required"`
	AnalysisType string   `json:"analysis_type" validate:"required"`
	Options      []string `json:"options,omitempty"`
	UserID       string   `json:"user_id,omitempty"`
}

// AnalyzeContentResponse 分析内容响应
type AnalyzeContentResponse struct {
	Analysis     map[string]interface{} `json:"analysis"`
	Summary      string                 `json:"summary"`
	Keywords     []string               `json:"keywords"`
	Sentiment    string                 `json:"sentiment,omitempty"`
	Topics       []string               `json:"topics,omitempty"`
	Language     string                 `json:"language,omitempty"`
	Confidence   float64                `json:"confidence"`
	ResponseTime time.Duration          `json:"response_time"`
}

// ContinueWritingRequest 续写内容请求
type ContinueWritingRequest struct {
	Content string            `json:"content" validate:"required"`
	Style   string            `json:"style,omitempty"`
	Length  int               `json:"length,omitempty"`
	Context map[string]string `json:"context,omitempty"`
	UserID  string            `json:"user_id,omitempty"`
}

// ContinueWritingResponse 续写内容响应
type ContinueWritingResponse struct {
	ContinuedContent string        `json:"continued_content"`
	OriginalLength   int           `json:"original_length"`
	AddedLength      int           `json:"added_length"`
	ResponseTime     time.Duration `json:"response_time"`
}

// OptimizeTextRequest 优化文本请求
type OptimizeTextRequest struct {
	Text             string   `json:"text" validate:"required"`
	OptimizationType string   `json:"optimization_type" validate:"required"`
	TargetAudience   string   `json:"target_audience,omitempty"`
	Style            string   `json:"style,omitempty"`
	Options          []string `json:"options,omitempty"`
	UserID           string   `json:"user_id,omitempty"`
}

// OptimizeTextResponse 优化文本响应
type OptimizeTextResponse struct {
	OptimizedText   string            `json:"optimized_text"`
	Changes         []string          `json:"changes"`
	Improvements    []string          `json:"improvements"`
	OriginalLength  int               `json:"original_length"`
	OptimizedLength int               `json:"optimized_length"`
	Suggestions     []string          `json:"suggestions,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	ResponseTime    time.Duration     `json:"response_time"`
}

// GenerateOutlineRequest 生成大纲请求
type GenerateOutlineRequest struct {
	Topic    string   `json:"topic" validate:"required"`
	Type     string   `json:"type,omitempty"`
	Depth    int      `json:"depth,omitempty"`
	Sections []string `json:"sections,omitempty"`
	Style    string   `json:"style,omitempty"`
	UserID   string   `json:"user_id,omitempty"`
}

// GenerateOutlineResponse 生成大纲响应
type GenerateOutlineResponse struct {
	Outline         []OutlineItem `json:"outline"`
	Title           string        `json:"title"`
	Summary         string        `json:"summary"`
	EstimatedLength int           `json:"estimated_length"`
	ResponseTime    time.Duration `json:"response_time"`
}

// OutlineItem 大纲项目
type OutlineItem struct {
	Level       int           `json:"level"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Children    []OutlineItem `json:"children,omitempty"`
}

// GetContextInfoRequest 获取上下文信息请求
type GetContextInfoRequest struct {
	ContextID string `json:"context_id" validate:"required"`
	UserID    string `json:"user_id,omitempty"`
}

// GetContextInfoResponse 获取上下文信息响应
type GetContextInfoResponse struct {
	ContextID    string            `json:"context_id"`
	MessageCount int               `json:"message_count"`
	TokensUsed   int               `json:"tokens_used"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Status       string            `json:"status"`
}

// UpdateContextWithFeedbackRequest 根据反馈更新上下文请求
type UpdateContextWithFeedbackRequest struct {
	ContextID string            `json:"context_id" validate:"required"`
	Feedback  string            `json:"feedback" validate:"required"`
	Rating    int               `json:"rating,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	UserID    string            `json:"user_id,omitempty"`
}

// UpdateContextWithFeedbackResponse 根据反馈更新上下文响应
type UpdateContextWithFeedbackResponse struct {
	ContextID    string        `json:"context_id"`
	Updated      bool          `json:"updated"`
	Changes      []string      `json:"changes,omitempty"`
	ResponseTime time.Duration `json:"response_time"`
}

// GetSupportedModelsResponse 获取支持的AI模型列表响应
type GetSupportedModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Type        string            `json:"type"`
	MaxTokens   int               `json:"max_tokens"`
	InputPrice  float64           `json:"input_price"`
	OutputPrice float64           `json:"output_price"`
	Features    []string          `json:"features"`
	Status      string            `json:"status"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// GetModelInfoRequest 获取模型信息请求
type GetModelInfoRequest struct {
	ModelID string `json:"model_id" validate:"required"`
}

// GetModelInfoResponse 获取模型信息响应
type GetModelInfoResponse struct {
	Model ModelInfo `json:"model"`
}

// ValidateAPIKeyRequest 验证API密钥请求
type ValidateAPIKeyRequest struct {
	Provider string `json:"provider" validate:"required"`
	APIKey   string `json:"api_key" validate:"required"`
}

// ValidateAPIKeyResponse 验证API密钥响应
type ValidateAPIKeyResponse struct {
	Valid        bool          `json:"valid"`
	Provider     string        `json:"provider"`
	Message      string        `json:"message,omitempty"`
	Quota        *APIQuota     `json:"quota,omitempty"`
	ResponseTime time.Duration `json:"response_time"`
}

// APIQuota API配额信息
type APIQuota struct {
	Used      int64     `json:"used"`
	Limit     int64     `json:"limit"`
	Remaining int64     `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
}

// Context相关请求响应结构体

// CreateContextRequest 创建上下文请求
type CreateContextRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	UserID      string            `json:"user_id" validate:"required"`
}

// CreateContextResponse 创建上下文响应
type CreateContextResponse struct {
	ContextID string    `json:"context_id"`
	CreatedAt time.Time `json:"created_at"`
}

// GetContextRequest 获取上下文请求
type GetContextRequest struct {
	ContextID string `json:"context_id" validate:"required"`
	UserID    string `json:"user_id,omitempty"`
}

// GetContextResponse 获取上下文响应
type GetContextResponse struct {
	Context ContextInfo `json:"context"`
}

// ContextInfo 上下文信息
type ContextInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	UserID      string            `json:"user_id"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// UpdateContextRequest 更新上下文请求
type UpdateContextRequest struct {
	ContextID   string            `json:"context_id" validate:"required"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	UserID      string            `json:"user_id,omitempty"`
}

// UpdateContextResponse 更新上下文响应
type UpdateContextResponse struct {
	Updated   bool      `json:"updated"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteContextRequest 删除上下文请求
type DeleteContextRequest struct {
	ContextID string `json:"context_id" validate:"required"`
	UserID    string `json:"user_id,omitempty"`
}

// DeleteContextResponse 删除上下文响应
type DeleteContextResponse struct {
	Deleted   bool      `json:"deleted"`
	DeletedAt time.Time `json:"deleted_at"`
}

// ListContextsRequest 列出上下文请求
type ListContextsRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Type     string `json:"type,omitempty"`
	Status   string `json:"status,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

// ListContextsResponse 列出上下文响应
type ListContextsResponse struct {
	Contexts   []ContextInfo `json:"contexts"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// AddMessageRequest 添加消息到上下文请求
type AddMessageRequest struct {
	ContextID string            `json:"context_id" validate:"required"`
	Role      string            `json:"role" validate:"required"`
	Content   string            `json:"content" validate:"required"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	UserID    string            `json:"user_id,omitempty"`
}

// AddMessageResponse 添加消息到上下文响应
type AddMessageResponse struct {
	MessageID string    `json:"message_id"`
	AddedAt   time.Time `json:"added_at"`
}

// GetMessagesRequest 获取上下文消息请求
type GetMessagesRequest struct {
	ContextID string `json:"context_id" validate:"required"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}

// GetMessagesResponse 获取上下文消息响应
type GetMessagesResponse struct {
	Messages []MessageInfo `json:"messages"`
	Total    int           `json:"total"`
}

// MessageInfo 消息信息
type MessageInfo struct {
	ID        string            `json:"id"`
	ContextID string            `json:"context_id"`
	Role      string            `json:"role"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
}

// ClearMessagesRequest 清空上下文消息请求
type ClearMessagesRequest struct {
	ContextID string `json:"context_id" validate:"required"`
	UserID    string `json:"user_id,omitempty"`
}

// ClearMessagesResponse 清空上下文消息响应
type ClearMessagesResponse struct {
	Cleared   bool      `json:"cleared"`
	ClearedAt time.Time `json:"cleared_at"`
}

// ExternalAPI相关请求响应结构体

// CallAPIRequest 调用外部API请求
type CallAPIRequest struct {
	Provider string            `json:"provider" validate:"required"`
	Endpoint string            `json:"endpoint" validate:"required"`
	Method   string            `json:"method" validate:"required"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     interface{}       `json:"body,omitempty"`
	Timeout  time.Duration     `json:"timeout,omitempty"`
	UserID   string            `json:"user_id,omitempty"`
}

// CallAPIResponse 调用外部API响应
type CallAPIResponse struct {
	StatusCode   int               `json:"status_code"`
	Headers      map[string]string `json:"headers"`
	Body         interface{}       `json:"body"`
	ResponseTime time.Duration     `json:"response_time"`
	Success      bool              `json:"success"`
	Error        string            `json:"error,omitempty"`
}

// GetAPIStatusRequest 获取API状态请求
type GetAPIStatusRequest struct {
	Provider string `json:"provider" validate:"required"`
}

// GetAPIStatusResponse 获取API状态响应
type GetAPIStatusResponse struct {
	Provider     string        `json:"provider"`
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
	LastChecked  time.Time     `json:"last_checked"`
	Message      string        `json:"message,omitempty"`
}

// GetAPIUsageRequest 获取API使用情况请求
type GetAPIUsageRequest struct {
	Provider  string    `json:"provider" validate:"required"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
}

// GetAPIUsageResponse 获取API使用情况响应
type GetAPIUsageResponse struct {
	Provider  string      `json:"provider"`
	Usage     APIUsage    `json:"usage"`
	Period    TimePeriod  `json:"period"`
	Breakdown []UsageItem `json:"breakdown,omitempty"`
}

// APIUsage API使用情况
type APIUsage struct {
	TotalRequests int64   `json:"total_requests"`
	TotalTokens   int64   `json:"total_tokens"`
	TotalCost     float64 `json:"total_cost"`
	SuccessRate   float64 `json:"success_rate"`
}

// TimePeriod 时间段
type TimePeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// UsageItem 使用项目
type UsageItem struct {
	Date     time.Time `json:"date"`
	Requests int64     `json:"requests"`
	Tokens   int64     `json:"tokens"`
	Cost     float64   `json:"cost"`
}

// RefreshAPIKeyRequest 刷新API密钥请求
type RefreshAPIKeyRequest struct {
	Provider string `json:"provider" validate:"required"`
	UserID   string `json:"user_id,omitempty"`
}

// RefreshAPIKeyResponse 刷新API密钥响应
type RefreshAPIKeyResponse struct {
	Provider    string    `json:"provider"`
	Refreshed   bool      `json:"refreshed"`
	RefreshedAt time.Time `json:"refreshed_at"`
	Message     string    `json:"message,omitempty"`
}

// Adapter相关请求响应结构体

// GetAdapterRequest 获取适配器请求
type GetAdapterRequest struct {
	AdapterID string `json:"adapter_id" validate:"required"`
}

// GetAdapterResponse 获取适配器响应
type GetAdapterResponse struct {
	Adapter AdapterInfo `json:"adapter"`
}

// AdapterInfo 适配器信息
type AdapterInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Provider    string            `json:"provider"`
	Version     string            `json:"version"`
	Status      string            `json:"status"`
	Config      map[string]string `json:"config"`
	Features    []string          `json:"features"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ListAdaptersRequest 列出适配器请求
type ListAdaptersRequest struct {
	Type     string `json:"type,omitempty"`
	Provider string `json:"provider,omitempty"`
	Status   string `json:"status,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

// ListAdaptersResponse 列出适配器响应
type ListAdaptersResponse struct {
	Adapters   []AdapterInfo `json:"adapters"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// RegisterAdapterRequest 注册适配器请求
type RegisterAdapterRequest struct {
	Name        string            `json:"name" validate:"required"`
	Type        string            `json:"type" validate:"required"`
	Provider    string            `json:"provider" validate:"required"`
	Version     string            `json:"version" validate:"required"`
	Config      map[string]string `json:"config,omitempty"`
	Features    []string          `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
}

// RegisterAdapterResponse 注册适配器响应
type RegisterAdapterResponse struct {
	AdapterID    string    `json:"adapter_id"`
	Registered   bool      `json:"registered"`
	RegisteredAt time.Time `json:"registered_at"`
}

// UnregisterAdapterRequest 注销适配器请求
type UnregisterAdapterRequest struct {
	AdapterID string `json:"adapter_id" validate:"required"`
}

// UnregisterAdapterResponse 注销适配器响应
type UnregisterAdapterResponse struct {
	Unregistered   bool      `json:"unregistered"`
	UnregisteredAt time.Time `json:"unregistered_at"`
}

// UpdateAdapterRequest 更新适配器请求
type UpdateAdapterRequest struct {
	AdapterID   string            `json:"adapter_id" validate:"required"`
	Name        string            `json:"name,omitempty"`
	Version     string            `json:"version,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
	Features    []string          `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status,omitempty"`
}

// UpdateAdapterResponse 更新适配器响应
type UpdateAdapterResponse struct {
	Updated   bool      `json:"updated"`
	UpdatedAt time.Time `json:"updated_at"`
}
