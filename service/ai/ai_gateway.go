package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/service/ai/adapter"
)

// AIGateway AI 网关 - 统一请求入口
type AIGateway struct {
	adapterManager *adapter.AdapterManager
	quotaService   *QuotaService
	rateLimiter    RateLimiter
	authService    AuthService
}

// UnifiedRequest 统一 AI 请求结构
type UnifiedRequest struct {
	RequestID   string                 `json:"request_id"`
	UserID      string                 `json:"user_id"`
	Model       string                 `json:"model,omitempty"`
	TaskType    string                 `json:"task_type"`
	Prompt      string                 `json:"prompt,omitempty"`
	Messages    []adapter.Message      `json:"messages,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
}

// UnifiedResponse 统一 AI 响应结构
type UnifiedResponse struct {
	RequestID string         `json:"request_id"`
	Success   bool           `json:"success"`
	Data      interface{}    `json:"data,omitempty"`
	Error     *GatewayError  `json:"error,omitempty"`
	Usage     *adapter.Usage `json:"usage,omitempty"`
	Model     string         `json:"model"`
	Provider  string         `json:"provider"`
	Latency   time.Duration  `json:"latency"`
}

// GatewayError 网关错误
type GatewayError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
}

// AuthService 鉴权服务接口
type AuthService interface {
	Authenticate(ctx context.Context, token string) (string, error)
	Authorize(ctx context.Context, userID, action string) bool
}

// NewAIGateway 创建 AI 网关
func NewAIGateway(
	adapterManager *adapter.AdapterManager,
	quotaService *QuotaService,
	rateLimiter RateLimiter,
	authService AuthService,
) *AIGateway {
	return &AIGateway{
		adapterManager: adapterManager,
		quotaService:   quotaService,
		rateLimiter:    rateLimiter,
		authService:    authService,
	}
}

// ProcessRequest 处理统一 AI 请求
// 执行流程：鉴权 → 限流 → 配额校验 → 模型选择 → 调用编排
func (g *AIGateway) ProcessRequest(ctx context.Context, req *UnifiedRequest, token string) (*UnifiedResponse, error) {
	startTime := time.Now()
	requestID := req.RequestID

	// 鉴权
	userID, err := g.authenticate(ctx, token)
	if err != nil {
		return g.buildErrorResponse(requestID, "AUTH_FAILED", err.Error(), time.Since(startTime)), nil
	}
	// 限流
	if !g.checkRateLimit(ctx, userID) {
		return g.buildErrorResponse(requestID, "RATE_LIMIT_EXCEEDED", "请求过于频繁", time.Since(startTime)), nil
	}
	// 配额校验
	estimatedTokens := g.estimateTokens(req)
	if err := g.checkQuota(ctx, userID, estimatedTokens); err != nil {
		return g.buildErrorResponse(requestID, "QUOTA_EXCEEDED", err.Error(), time.Since(startTime)), nil
	}
	// 模型选择
	selectedModel, selectedProvider, err := g.selectModel(ctx, req)
	if err != nil {
		return g.buildErrorResponse(requestID, "MODEL_SELECTION_FAILED", err.Error(), time.Since(startTime)), nil
	}
	// 调用编排
	result, usage, err := g.orchestrateCall(ctx, req, selectedProvider, selectedModel)
	if err != nil {
		return g.buildErrorResponse(requestID, "EXECUTION_FAILED", err.Error(), time.Since(startTime)), nil
	}
	// 消费配额
	if err := g.consumeQuota(ctx, userID, usage.TotalTokens); err != nil {
		fmt.Printf("配额消费失败: %v\n", err)
	}

	latency := time.Since(startTime)
	return &UnifiedResponse{
		RequestID: requestID,
		Success:   true,
		Data:      result,
		Usage:     usage,
		Model:     selectedModel,
		Provider:  selectedProvider,
		Latency:   latency,
	}, nil
}

// authenticate 鉴权
func (g *AIGateway) authenticate(ctx context.Context, token string) (string, error) {
	if g.authService == nil {
		return "anonymous", nil
	}
	return g.authService.Authenticate(ctx, token)
}

// checkRateLimit 限流检查
func (g *AIGateway) checkRateLimit(ctx context.Context, userID string) bool {
	if g.rateLimiter == nil {
		return true
	}
	return g.rateLimiter.Allow(ctx, "user:"+userID)
}

// checkQuota 配额校验
func (g *AIGateway) checkQuota(ctx context.Context, userID string, tokens int) error {
	if g.quotaService == nil {
		return nil
	}
	return g.quotaService.Check(ctx, userID, tokens)
}

// consumeQuota 消费配额
func (g *AIGateway) consumeQuota(ctx context.Context, userID string, tokens int) error {
	if g.quotaService == nil {
		return nil
	}
	err := g.quotaService.ConsumeQuota(ctx, userID, tokens, "ai-gateway", "default", "consume")
	return err
}

// selectModel 模型选择
func (g *AIGateway) selectModel(ctx context.Context, req *UnifiedRequest) (string, string, error) {
	if req.Model != "" {
		adapter, err := g.adapterManager.GetAdapterByModel(req.Model)
		if err != nil {
			return "", "", fmt.Errorf("不支持的模型: %s", req.Model)
		}
		return req.Model, adapter.GetName(), nil
	}

	adapter, err := g.adapterManager.GetDefaultAdapter()
	if err != nil {
		return "", "", fmt.Errorf("未配置默认模型")
	}

	models := adapter.GetSupportedModels()
	if len(models) == 0 {
		return "", "", fmt.Errorf("适配器 %s 无可用模型", adapter.GetName())
	}

	return models[0], adapter.GetName(), nil
}

// orchestrateCall 调用编排
func (g *AIGateway) orchestrateCall(
	ctx context.Context,
	req *UnifiedRequest,
	provider, model string,
) (interface{}, *adapter.Usage, error) {
	adapter, err := g.adapterManager.GetAdapter(provider)
	if err != nil {
		return nil, nil, err
	}

	switch req.TaskType {
	case "chat":
		return g.executeChat(ctx, adapter, req, model)
	case "text_generation":
		return g.executeTextGeneration(ctx, adapter, req, model)
	default:
		return nil, nil, fmt.Errorf("不支持的任务类型: %s", req.TaskType)
	}
}

// executeChat 执行聊天任务
func (g *AIGateway) executeChat(
	ctx context.Context,
	aiAdapter adapter.AIAdapter,
	req *UnifiedRequest,
	model string,
) (interface{}, *adapter.Usage, error) {
	chatReq := &adapter.ChatCompletionRequest{
		Model:       model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	resp, err := aiAdapter.ChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, nil, err
	}

	return resp.Message, &resp.Usage, nil
}

// executeTextGeneration 执行文本生成任务
func (g *AIGateway) executeTextGeneration(
	ctx context.Context,
	aiAdapter adapter.AIAdapter,
	req *UnifiedRequest,
	model string,
) (interface{}, *adapter.Usage, error) {
	textReq := &adapter.TextGenerationRequest{
		Model:       model,
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	resp, err := aiAdapter.TextGeneration(ctx, textReq)
	if err != nil {
		return nil, nil, err
	}

	return resp.Text, &resp.Usage, nil
}

// estimateTokens 估算 token 数量
func (g *AIGateway) estimateTokens(req *UnifiedRequest) int {
	total := 0

	if req.Prompt != "" {
		total += len(req.Prompt) / 3
	}

	for _, msg := range req.Messages {
		total += len(msg.Content) / 3
	}

	if req.MaxTokens > 0 {
		total += req.MaxTokens
	}

	if total < 100 {
		return 100
	}
	return total
}

// buildErrorResponse 构建错误响应
func (g *AIGateway) buildErrorResponse(requestID, code, message string, latency time.Duration) *UnifiedResponse {
	return &UnifiedResponse{
		RequestID: requestID,
		Success:   false,
		Error: &GatewayError{
			Code:    code,
			Message: message,
		},
		Latency: latency,
	}
}
