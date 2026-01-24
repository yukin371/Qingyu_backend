// Qingyu_backend/test/integration/ai_migration_test.go

package integration

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	"Qingyu_backend/pkg/circuitbreaker"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContext 测试上下文
type TestContext struct {
	aiService      *ai.AIService
	quotaService   *ai.QuotaService
	grpcClient     *ai.GRPCClient
	initialQuota   int64
	circuitBreaker *circuitbreaker.CircuitBreaker
}

func setupTestEnvironment(t *testing.T) *TestContext {
	ctx := &TestContext{}

	// TODO: 初始化 AI 服务
	// ctx.aiService = ai.NewAIService(...)
	// ctx.quotaService = ai.NewQuotaService(...)
	// ctx.grpcClient = ai.NewGRPCClient(...)

	// TODO: 设置初始配额
	// ctx.initialQuota = 100000

	t.Skip("Integration tests require full environment setup - skipping for now")

	return ctx
}

func teardownTestEnvironment(t *testing.T, ctx *TestContext) {
	// TODO: 清理测试环境
}

func TestAIMigrationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 设置测试环境
	ctx := setupTestEnvironment(t)
	defer teardownTestEnvironment(t, ctx)

	t.Run("完整AI调用流程", func(t *testing.T) {
		// 1. 发起 AI 请求
		req := &ai.AgentRequest{
			UserID:       "test-user-123",
			WorkflowType: "chat",
			Parameters: map[string]interface{}{
				"message": "Hello, AI!",
			},
		}

		resp, err := ctx.aiService.ExecuteAgent(context.Background(), req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Content)
		assert.Greater(t, resp.TokensUsed, 0)

		// 2. 验证配额扣除
		quota, err := ctx.quotaService.CheckQuota(context.Background(), "test-user-123", 1000)
		require.NoError(t, err)
		assert.Less(t, quota.Remaining, ctx.initialQuota)

		// 3. 验证 AI 服务记录（通过 gRPC 查询）
		consumption, err := ctx.grpcClient.GetQuotaConsumption(
			context.Background(),
			&pb.QuotaConsumptionQuery{
				UserId:     "test-user-123",
				TimeRange: "day",
			},
		)
		require.NoError(t, err)
		assert.Greater(t, consumption.TotalTokens, 0)
	})

	t.Run("熔断器测试", func(t *testing.T) {
		// TODO: 停止 AI 服务模拟故障
		// ctx.stopAIService()
		// defer ctx.startAIService()

		// 连续失败触发熔断
		for i := 0; i < 6; i++ {
			_, err := ctx.aiService.ExecuteAgent(
				context.Background(),
				&ai.AgentRequest{UserID: "test-user", WorkflowType: "chat"},
			)
			assert.Error(t, err)
		}

		// 验证熔断器打开
		assert.Equal(t, circuitbreaker.StateOpen, ctx.aiService.GetCircuitBreakerState())

		// 验证降级生效
		if ctx.aiService.HasFallback() {
			resp, err := ctx.aiService.ExecuteAgent(
				context.Background(),
				&ai.AgentRequest{UserID: "test-user", WorkflowType: "chat"},
			)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		}
	})

	t.Run("配额一致性测试", func(t *testing.T) {
		userID := "test-user-quota"

		// 记录初始配额
		initialQuota, _ := ctx.quotaService.CheckQuota(context.Background(), userID, 0)

		// 执行 AI 调用
		_, err := ctx.aiService.ExecuteAgent(
			context.Background(),
			&ai.AgentRequest{UserID: userID, WorkflowType: "chat"},
		)
		require.NoError(t, err)

		// 等待同步
		time.Sleep(2 * time.Second)

		// 验证后端配额已扣除
		backendQuota, _ := ctx.quotaService.CheckQuota(context.Background(), userID, 0)
		assert.Less(t, backendQuota.Remaining, initialQuota.Remaining)

		// 验证 AI 服务记录
		aiConsumption, _ := ctx.grpcClient.GetQuotaConsumption(
			context.Background(),
			&pb.QuotaConsumptionQuery{UserId: userID, TimeRange: "day"},
		)
		assert.Greater(t, aiConsumption.TotalTokens, 0)

		// 验证一致性（误差 < 1%）
		diff := abs(initialQuota.Remaining - backendQuota.Remaining - int64(aiConsumption.TotalTokens))
		assert.Less(t, float64(diff)/float64(aiConsumption.TotalTokens), 0.01)
	})
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
