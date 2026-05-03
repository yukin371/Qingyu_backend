package ai

import (
	"context"
	"time"

	"Qingyu_backend/pkg/circuitbreaker"
	pkgErrors "Qingyu_backend/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	grpcClient     *GRPCClient
	quotaService   *QuotaService // 假设存在配额服务
	circuitBreaker *circuitbreaker.CircuitBreaker
	config         *AIServiceConfig
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
	return nil, pkgErrors.NewAIError(
		pkgErrors.ErrAIUnavailable,
		"AI service unavailable and no fallback configured",
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
	return false
}

// Close 关闭服务
func (s *AIService) Close() error {
	return s.grpcClient.Close()
}

// Legacy Service compatibility
// 为了保持向后兼容，保留旧的 Service 结构，但标记为废弃

// Deprecated: Use AIService instead
type Service struct{}

// Deprecated: Use NewAIService instead
func NewServiceWithDependencies(projectService interface{}) *Service {
	logrus.Warn("NewServiceWithDependencies is deprecated. Use NewAIService with gRPC client instead.")
	return &Service{}
}
