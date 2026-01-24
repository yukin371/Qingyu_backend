package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/pkg/circuitbreaker"
	"Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/ai/adapter"

	pb "Qingyu_backend/pkg/grpc/pb"

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
		quota, err := s.quotaService.CheckQuota(ctx, req.UserID)
		if err != nil {
			return nil, fmt.Errorf("quota check failed: %w", err)
		}

		if quota.Remaining <= 0 {
			return nil, errors.NewAIError(
				errors.ErrQuotaExhausted,
				"Insufficient quota",
				nil,
			)
		}
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
				return nil, errors.NewAIError(
					errors.ErrAIUnavailable,
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
		s.quotaService.ConsumeQuota(ctx, req.UserID, resp.TokensUsed, "ai_generation")
	}

	return resp, nil
}

// executeFallback 降级执行
func (s *AIService) executeFallback(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	if s.fallbackAdapter == nil {
		return nil, errors.NewAIError(
			errors.ErrAIUnavailable,
			"AI service unavailable and no fallback configured",
			nil,
		)
	}

	logrus.Warn("Using deprecated fallback adapter")

	// 根据工作流类型选择适配器
	// 这里需要根据实际的适配器接口进行调用
	// 由于适配器接口可能不兼容，这里返回一个通用错误
	return nil, errors.NewAIError(
		errors.ErrAIUnavailable,
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

// QuotaService 配额服务接口（假设存在）
type QuotaService struct{}

// CheckQuota 检查配额
func (q *QuotaService) CheckQuota(ctx context.Context, userID string) (*QuotaInfo, error) {
	// TODO: 实现配额检查逻辑
	return &QuotaInfo{Remaining: 1000}, nil
}

// ConsumeQuota 消费配额
func (q *QuotaService) ConsumeQuota(ctx context.Context, userID string, amount int64, workflowType string) error {
	// TODO: 实现配额消费逻辑
	return nil
}

// QuotaInfo 配额信息
type QuotaInfo struct {
	Remaining int64
	Total     int64
	Used      int64
}
