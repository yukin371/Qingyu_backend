package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Qingyu_backend/pkg/circuitbreaker"
	pkgErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/ai/adapter"
	aiInterfaces "Qingyu_backend/service/interfaces/ai"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/sirupsen/logrus"
)

// AIServiceConfig AI服务配置
type AIServiceConfig struct {
	Endpoint       string        // AI 服务 gRPC 端点
	Timeout        time.Duration // 请求超时
	MaxRetries     int           // 最大重试次数
	RetryDelay     time.Duration // 重试延迟
	EnableFallback bool          // 启用降级
	EnableMonitor  bool          // 启用监控与追踪
}

// AIService AI服务（简化版，使用gRPC）
type AIService struct {
	grpcClient      *GRPCClient
	quotaService    *QuotaService // 假设存在配额服务
	fallbackAdapter *adapter.AdapterManager // 废弃的适配器（用于降级）
	circuitBreaker  *circuitbreaker.CircuitBreaker
	config          *AIServiceConfig
}

// NewAIService 创建 AI 服务
func NewAIService(
	conn *grpc.ClientConn,
	quotaService *QuotaService,
	config *AIServiceConfig,
) *AIService {
	if config == nil {
		config = &AIServiceConfig{
			Endpoint:   "localhost:50051",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			RetryDelay: time.Second,
		}
	}

	return &AIService{
		grpcClient:     NewGRPCClient(conn, config),
		quotaService:   quotaService,
		circuitBreaker: circuitbreaker.NewCircuitBreaker(5, 60*time.Second, 3),
		config:         config,
	}
}

// SetFallbackAdapter 设置降级适配器
func (s *AIService) SetFallbackAdapter(adapter *adapter.AdapterManager) {
	s.fallbackAdapter = adapter
	logrus.Warn("Fallback adapter configured. This is deprecated and will be removed in v2.0.0")
}

// ExecuteAgent 执行 AI Agent
func (s *AIService) ExecuteAgent(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	// 1. 检查熔断器
	if !s.circuitBreaker.AllowRequest() {
		logrus.Warn("Circuit breaker open, using fallback")
		return s.executeFallback(ctx, req)
	}

	// 2. 检查配额（如果配额服务存在）
	if s.quotaService != nil {
		// 配额检查集成在消费中，这里只检查是否存在错误
		// 实际的配额检查在 ConsumeQuota 中进行
	}

	// 3. 调用 AI 服务
	resp, err := s.grpcClient.ExecuteAgentWithRetry(ctx, req)
	if err != nil {
		s.circuitBreaker.RecordFailure()

		// 降级处理
		if s.config.EnableFallback {
			logrus.Warnf("AI service failed: %v, using fallback", err)
			return s.executeFallback(ctx, req)
		}

		// 判断错误类型
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.Unavailable {
				return nil, pkgErrors.NewAIError(
					pkgErrors.ErrAIUnavailable,
					"AI service unavailable",
					err,
				)
			}
		}
		return nil, err
	}

	s.circuitBreaker.RecordSuccess()

	// 4. 消费配额（如果配额服务存在）
	if s.quotaService != nil {
		err := s.quotaService.ConsumeQuota(ctx, req.UserID, int(resp.TokensUsed), "ai-service", "default", req.WorkflowType)
		if err != nil {
			// 配额消费失败记录日志，但不影响响应
			logrus.Errorf("Failed to consume quota: %v", err)
		}
	}

	return resp, nil
}

// executeFallback 降级执行
func (s *AIService) executeFallback(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	if s.fallbackAdapter == nil {
		return nil, pkgErrors.NewAIError(
			pkgErrors.ErrAIUnavailable,
			"AI service unavailable and no fallback configured",
			nil,
		)
	}

	logrus.Warn("Using deprecated fallback adapter")

	// 根据工作流类型选择适配器
	// 这里需要根据实际的适配器接口进行调用
	// 由于适配器接口可能不兼容，这里返回一个通用错误
	return nil, pkgErrors.NewAIError(
		pkgErrors.ErrAIUnavailable,
		"Fallback adapter not implemented yet",
		nil,
	)
}

// HealthCheck 健康检查
func (s *AIService) HealthCheck(ctx context.Context) error {
	return s.grpcClient.HealthCheck(ctx)
}

// GetCircuitBreakerState 获取熔断器状态
func (s *AIService) GetCircuitBreakerState() circuitbreaker.CircuitState {
	return s.circuitBreaker.GetState()
}

// GetCircuitBreakerStats 获取熔断器统计
func (s *AIService) GetCircuitBreakerStats() map[string]interface{} {
	return s.circuitBreaker.GetStats()
}

// HasFallback 是否有降级适配器
func (s *AIService) HasFallback() bool {
	return s.fallbackAdapter != nil
}

// Close 关闭服务
func (s *AIService) Close() error {
	return s.grpcClient.Close()
}

// Legacy Service compatibility
// 为了保持向后兼容，保留旧的 Service 结构，但标记为废弃

// Deprecated: Use AIService instead
type Service struct {
	legacyService *AIService
	adapterManager *adapter.AdapterManager
}

// Deprecated: Use NewAIService instead
func NewService() *Service {
	logrus.Warn("NewService is deprecated. Use NewAIService with gRPC client instead.")
	return &Service{}
}

// Deprecated: Use NewAIService instead
func NewServiceWithDependencies(projectService interface{}) *Service {
	logrus.Warn("NewServiceWithDependencies is deprecated. Use NewAIService with gRPC client instead.")
	return &Service{}
}

// GenerateContent 生成内容（废弃实现，仅用于保持接口兼容）
// Deprecated: Use AIService.ExecuteAgent instead
func (s *Service) GenerateContent(ctx context.Context, req *aiInterfaces.GenerateContentRequest) (*aiInterfaces.GenerateContentResponse, error) {
	logrus.Warn("Service.GenerateContent is deprecated. Use AIService with gRPC client instead.")
	return nil, errors.New("deprecated method, use AIService instead")
}

// GenerateContentStream 流式生成内容（废弃实现，仅用于保持接口兼容）
// Deprecated: Use AIService instead
func (s *Service) GenerateContentStream(ctx context.Context, req *GenerateContentRequest) (<-chan *aiInterfaces.GenerateContentResponse, error) {
	logrus.Warn("Service.GenerateContentStream is deprecated. Use AIService with gRPC client instead.")

	// 返回一个包含错误结果的 channel
	ch := make(chan *aiInterfaces.GenerateContentResponse, 1)
	ch <- &aiInterfaces.GenerateContentResponse{
		Content:      "",
		Model:        "deprecated",
		TokensUsed:   0,
		FinishReason: "error",
		Metadata:     map[string]string{"error": "streaming not supported in deprecated service"},
	}
	close(ch)
	return ch, errors.New("streaming not supported in deprecated service")
}

// GenerateContentStreamWithInterface 流式生成内容（接口实现，接受接口类型）
// Deprecated: Use AIService instead
func (s *Service) GenerateContentStreamWithInterface(ctx context.Context, req *aiInterfaces.GenerateContentRequest) (<-chan *aiInterfaces.GenerateContentResponse, error) {
	logrus.Warn("Service.GenerateContentStreamWithInterface is deprecated. Use AIService with gRPC client instead.")

	// 返回一个包含错误结果的 channel
	ch := make(chan *aiInterfaces.GenerateContentResponse, 1)
	ch <- &aiInterfaces.GenerateContentResponse{
		Content:      "",
		Model:        "deprecated",
		TokensUsed:   0,
		FinishReason: "error",
		Metadata:     map[string]string{"error": "streaming not supported in deprecated service"},
	}
	close(ch)
	return ch, errors.New("streaming not supported in deprecated service")
}

// QuotaInfo 配额信息
type QuotaInfo struct {
	Remaining int64
	Total     int64
	Used      int64
}

// ============ 向后兼容类型定义 ============
// 这些类型定义用于保持与现有 API 层的兼容性

// ContinueWritingRequest 续写请求（向后兼容）
type ContinueWritingRequest struct {
	ProjectID      string                `json:"projectId"`
	ChapterID      string                `json:"chapterId"`
	CurrentText    string                `json:"currentText"`
	ContinueLength int                   `json:"continueLength"`
	Options        interface{}           `json:"options"`
}

// ContinueWritingResponse 续写响应（向后兼容）
type ContinueWritingResponse struct {
	Content      string `json:"content"`
	TokensUsed   int    `json:"tokensUsed"`
	Model        string `json:"model"`
	FinishReason string `json:"finishReason"`
}

// OptimizeTextRequest 文本优化请求（向后兼容）
type OptimizeTextRequest struct {
	ProjectID    string      `json:"projectId"`
	ChapterID    string      `json:"chapterId"`
	OriginalText string      `json:"originalText"`
	OptimizeType string      `json:"optimizeType"`
	Instructions string      `json:"instructions"`
	Options      interface{} `json:"options"`
}

// OptimizeTextResponse 文本优化响应（向后兼容）
type OptimizeTextResponse struct {
	Content      string `json:"content"`
	TokensUsed   int    `json:"tokensUsed"`
	Model        string `json:"model"`
	Changes      []string `json:"changes,omitempty"`
}

// GenerateContentRequest 内容生成请求（向后兼容）
type GenerateContentRequest struct {
	ProjectID  string                 `json:"projectId"`
	ChapterID  string                 `json:"chapterId"`
	Prompt     string                 `json:"prompt"`
	Options    interface{}            `json:"options"`
	MaxTokens  int                    `json:"maxTokens,omitempty"`
	Temperature float64               `json:"temperature,omitempty"`
	Context    map[string]string      `json:"context,omitempty"`
	UserID     string                 `json:"userId,omitempty"`
}

// GenerateContentResponse 内容生成响应（向后兼容）
type GenerateContentResponse struct {
	Content      string `json:"content"`
	TokensUsed   int    `json:"tokensUsed"`
	Model        string `json:"model"`
	FinishReason string `json:"finishReason"`
	RequestID    string `json:"requestId,omitempty"`
}

// ============ 向后兼容方法实现 ============

// GetAdapterManager 获取适配器管理器（废弃方法，返回 nil）
// Deprecated: 适配器已废弃，此方法仅为向后兼容
func (s *Service) GetAdapterManager() *adapter.AdapterManager {
	logrus.Warn("GetAdapterManager is deprecated and will always return nil")
	return s.adapterManager
}

// ContinueWriting 智能续写（向后兼容实现）
func (s *Service) ContinueWriting(ctx context.Context, req *ContinueWritingRequest) (*ContinueWritingResponse, error) {
	logrus.Warn("Service.ContinueWriting is deprecated. Use new AIService with gRPC instead.")

	// 构建提示词
	prompt := fmt.Sprintf("请基于以下内容进行续写，保持风格和情节的连贯性：\n\n%s", req.CurrentText)
	if req.ContinueLength > 0 {
		prompt += fmt.Sprintf("\n\n请续写约%d字的内容。", req.ContinueLength)
	}

	// 调用 GenerateContent - 需要转换为接口类型
	genReq := &aiInterfaces.GenerateContentRequest{
		Model:    "gpt-4", // 默认模型
		Prompt:   prompt,
		MaxTokens: 2000,
	}

	genResp, err := s.GenerateContent(ctx, genReq)
	if err != nil {
		return nil, err
	}

	return &ContinueWritingResponse{
		Content:      genResp.Content,
		TokensUsed:   genResp.TokensUsed,
		Model:        genResp.Model,
		FinishReason: genResp.FinishReason,
	}, nil
}

// OptimizeText 文本优化（向后兼容实现）
func (s *Service) OptimizeText(ctx context.Context, req *OptimizeTextRequest) (*OptimizeTextResponse, error) {
	logrus.Warn("Service.OptimizeText is deprecated. Use new AIService with gRPC instead.")

	// 构建提示词
	var prompt string
	switch req.OptimizeType {
	case "expand":
		prompt = "请对以下文本进行扩写，增加细节描述和情节内容："
	case "shorten":
		prompt = "请对以下文本进行缩写，保留核心内容："
	case "style":
		prompt = "请对以下文本进行润色，优化表达方式："
	default:
		prompt = "请优化以下文本："
	}

	if req.Instructions != "" {
		prompt += fmt.Sprintf("\n\n具体要求：%s", req.Instructions)
	}

	prompt += fmt.Sprintf("\n\n原文：\n%s", req.OriginalText)

	// 调用 GenerateContent - 需要转换为接口类型
	genReq := &aiInterfaces.GenerateContentRequest{
		Model:    "gpt-4", // 默认模型
		Prompt:   prompt,
		MaxTokens: 2000,
	}

	genResp, err := s.GenerateContent(ctx, genReq)
	if err != nil {
		return nil, err
	}

	return &OptimizeTextResponse{
		Content:    genResp.Content,
		TokensUsed: genResp.TokensUsed,
		Model:      genResp.Model,
		Changes:    []string{"文本已优化"},
	}, nil
}
