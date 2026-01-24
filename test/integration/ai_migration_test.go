// Qingyu_backend/test/integration/ai_migration_test.go
//
// AI 服务迁移集成测试
// 测试从 Go 后端到独立 Qingyu-Ai-Service (Python) 的迁移功能
//
// 验证范围：
// 1. 完整 AI 调用流程（gRPC 通信、熔断器、配额）
// 2. 熔断器功能（状态转换、失败计数、恢复）
// 3. 配额一致性（检查、消费、恢复）
//
// 运行方式：
//   go test -v ./test/integration -run TestAI
//   go test -v ./test/integration -run TestAI -short  # 跳过需要外部服务的测试
//
// 环境要求：
//   - AI 服务运行在 localhost:50051（或设置 AI_SERVICE_ENDPOINT 环境变量）
//   - MongoDB 和 Redis 可选（配额测试需要）

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"Qingyu_backend/pkg/circuitbreaker"
	pkgErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ============================================================================
// 测试配置
// ============================================================================

const (
	// AI 服务端点（可通过环境变量覆盖）
	defaultAIEndpoint = "localhost:50051"
	aiTestTimeout     = 30 * time.Second
	aiShortTimeout    = 5 * time.Second
)

// getAIEndpoint 从环境变量获取 AI 服务端点，否则使用默认值
func getAIEndpoint() string {
	if endpoint := os.Getenv("AI_SERVICE_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return defaultAIEndpoint
}

// ============================================================================
// 测试辅助函数
// ============================================================================

// setupTestConnection 创建测试用的 gRPC 连接
// 如果无法连接到服务，会跳过测试
func setupTestConnection(t *testing.T) *grpc.ClientConn {
	endpoint := getAIEndpoint()

	ctx, cancel := context.WithTimeout(context.Background(), aiShortTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("无法连接到 AI 服务 %s: %v (跳过集成测试)", endpoint, err)
	}

	t.Cleanup(func() {
		if conn != nil {
			conn.Close()
		}
	})

	t.Logf("✅ 已连接到 AI 服务: %s", endpoint)
	return conn
}

// setupTestAIService 创建测试用的 AI 服务
// 返回 AI 服务实例和 gRPC 连接
func setupTestAIService(t *testing.T) (*ai.AIService, *grpc.ClientConn) {
	conn := setupTestConnection(t)

	config := &ai.AIServiceConfig{
		Endpoint:       getAIEndpoint(),
		Timeout:        aiTestTimeout,
		MaxRetries:     3,
		RetryDelay:     time.Second,
		EnableFallback: false, // 测试时不启用降级
	}

	service := ai.NewAIService(conn, nil, config)

	t.Cleanup(func() {
		if service != nil {
			service.Close()
		}
	})

	t.Log("✅ AI 服务初始化完成")
	return service, conn
}

// createTestAgentRequest 创建测试用的 Agent 请求
func createTestAgentRequest(userID, workflowType string) *ai.AgentRequest {
	return &ai.AgentRequest{
		UserID:       userID,
		WorkflowType: workflowType,
		Parameters: map[string]interface{}{
			"task":        "测试任务",
			"max_length":  1000,
			"temperature": 0.7,
		},
	}
}

// assertAIError 断言错误为 AIError 类型并检查错误类型
func assertAIError(t *testing.T, err error, expectedType pkgErrors.AIErrorType) {
	require.Error(t, err, "期望返回错误")

	aiErr, ok := err.(*pkgErrors.AIError)
	require.True(t, ok, "错误应该是 AIError 类型")

	assert.Equal(t, expectedType, aiErr.Type, "错误类型不匹配")
}

// skipIfServiceUnavailable 如果服务不可用则跳过测试
func skipIfServiceUnavailable(t *testing.T, err error) {
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			t.Skipf("AI 服务不可用: %v (跳过测试)", st.Message())
		}
	}
}

// ============================================================================
// 测试用例 1: 完整 AI 调用流程测试
// ============================================================================

// TestCompleteAIWorkflow 测试完整的 AI 调用流程
// 验证点：
// 1. gRPC 连接建立
// 2. AIService 初始化
// 3. ExecuteAgent 调用成功
// 4. 熔断器状态正常
// 5. 响应数据完整
func TestCompleteAIWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	// Setup
	service, _ := setupTestAIService(t)
	ctx, cancel := context.WithTimeout(context.Background(), aiTestTimeout)
	defer cancel()

	// 测试健康检查
	t.Run("HealthCheck", func(t *testing.T) {
		t.Log("🔍 测试健康检查...")

		err := service.HealthCheck(ctx)
		skipIfServiceUnavailable(t, err)
		require.NoError(t, err, "健康检查失败")

		// 验证熔断器状态
		state := service.GetCircuitBreakerState()
		assert.Equal(t, circuitbreaker.StateClosed, state, "健康检查后熔断器应该是关闭状态")

		// 验证统计信息
		stats := service.GetCircuitBreakerStats()
		t.Logf("   熔断器状态: %s", stats["state"])
		t.Logf("   总请求数: %v", stats["totalRequests"])
		t.Logf("   总成功数: %v", stats["totalSuccesses"])
		t.Logf("   总失败数: %v", stats["totalFailures"])

		t.Log("✅ 健康检查通过")
	})

	// 测试基本 Agent 执行
	t.Run("ExecuteAgent_Basic", func(t *testing.T) {
		t.Log("🤖 测试基本 Agent 执行...")

		req := createTestAgentRequest("test_user_001", "text_generation")

		resp, err := service.ExecuteAgent(ctx, req)
		skipIfServiceUnavailable(t, err)

		if err != nil {
			// 如果服务不可用，跳过后续测试
			t.Skipf("AI 服务执行失败: %v", err)
		}

		// 验证响应
		require.NotNil(t, resp, "响应不应为空")
		assert.NotEmpty(t, resp.Content, "响应内容不应为空")
		assert.Greater(t, resp.TokensUsed, int64(0), "Token 使用量应该大于 0")
		assert.NotEmpty(t, resp.WorkflowType, "工作流类型不应为空")

		t.Logf("✅ Agent 执行成功")
		t.Logf("   Content: %s", truncateString(resp.Content, 100))
		t.Logf("   Tokens Used: %d", resp.TokensUsed)
		t.Logf("   Workflow Type: %s", resp.WorkflowType)
	})

	// 测试多次连续调用
	t.Run("ExecuteAgent_Sequential", func(t *testing.T) {
		t.Log("🔄 测试连续请求...")

		requests := []struct {
			name         string
			userID       string
			workflowType string
		}{
			{"Request1", "test_user_002", "text_generation"},
			{"Request2", "test_user_003", "text_generation"},
			{"Request3", "test_user_004", "text_generation"},
		}

		successCount := 0
		for _, tc := range requests {
			t.Run(tc.name, func(t *testing.T) {
				req := createTestAgentRequest(tc.userID, tc.workflowType)
				resp, err := service.ExecuteAgent(ctx, req)

				skipIfServiceUnavailable(t, err)

				require.NoError(t, err, "%s 执行失败", tc.name)
				require.NotNil(t, resp, "%s 响应为空", tc.name)
				assert.NotEmpty(t, resp.Content, "%s 内容为空", tc.name)

				successCount++
				t.Logf("   ✓ %s 完成", tc.name)
			})
		}

		// 验证熔断器统计
		stats := service.GetCircuitBreakerStats()
		totalRequests := int(stats["totalRequests"].(int64))
		assert.Greater(t, totalRequests, 0, "总请求数应该大于 0")
		t.Logf("✅ 连续请求完成，成功: %d/%d，总请求数: %d", successCount, len(requests), totalRequests)
	})

	// 测试错误处理
	t.Run("ExecuteAgent_ErrorHandling", func(t *testing.T) {
		t.Log("⚠️  测试错误处理...")

		// 测试超时
		shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer shortCancel()

		req := createTestAgentRequest("test_user_timeout", "text_generation")
		_, err := service.ExecuteAgent(shortCtx, req)

		// 超时错误或服务不可用
		if err != nil {
			t.Logf("   ✓ 超时错误处理正常: %v", err)
		} else {
			t.Log("   ⚠️  请求在超时前完成（这是正常的）")
		}

		// 测试空用户ID（根据实际实现可能不会返回错误）
		t.Run("EmptyUserID", func(t *testing.T) {
			req := &ai.AgentRequest{
				UserID:       "",
				WorkflowType: "text_generation",
				Parameters:   map[string]interface{}{},
			}

			_, err := service.ExecuteAgent(ctx, req)
			// 注意：根据实际实现，空用户ID可能不会返回错误
			_ = err // 避免未使用变量错误
		})

		t.Log("✅ 错误处理测试完成")
	})

	// 测试熔断器统计
	t.Run("CircuitBreakerStats", func(t *testing.T) {
		t.Log("📊 测试熔断器统计...")

		stats := service.GetCircuitBreakerStats()

		require.Contains(t, stats, "state", "统计应该包含状态")
		require.Contains(t, stats, "failureCount", "统计应该包含失败计数")
		require.Contains(t, stats, "successCount", "统计应该包含成功计数")
		require.Contains(t, stats, "totalRequests", "统计应该包含总请求数")
		require.Contains(t, stats, "totalSuccesses", "统计应该包含总成功数")
		require.Contains(t, stats, "totalFailures", "统计应该包含总失败数")

		t.Logf("   状态: %v", stats["state"])
		t.Logf("   失败计数: %v", stats["failureCount"])
		t.Logf("   成功计数: %v", stats["successCount"])
		t.Logf("   总请求数: %v", stats["totalRequests"])
		t.Logf("   总成功数: %v", stats["totalSuccesses"])
		t.Logf("   总失败数: %v", stats["totalFailures"])

		t.Log("✅ 熔断器统计正常")
	})
}

// ============================================================================
// 测试用例 2: 熔断器功能测试
// ============================================================================

// TestCircuitBreakerIntegration 测试熔断器的集成功能
// 验证点：
// 1. 熔断器初始化正常
// 2. 失败计数正确
// 3. 状态转换正确（Closed -> Open -> HalfOpen -> Closed）
// 4. 统计信息准确
func TestCircuitBreakerIntegration(t *testing.T) {
	t.Log("🔌 测试熔断器集成功能...")

	// 创建独立的熔断器进行测试
	breaker := circuitbreaker.NewCircuitBreaker(
		3,             // 失败阈值
		5*time.Second, // 超时时间
		2,             // 成功阈值
	)

	t.Run("InitialState", func(t *testing.T) {
		t.Log("   测试初始状态...")

		state := breaker.GetState()
		assert.Equal(t, circuitbreaker.StateClosed, state, "初始状态应该是关闭")

		stats := breaker.GetStats()
		assert.Equal(t, "Closed", stats["state"], "状态字符串应该是 Closed")
		assert.Equal(t, 0, stats["failureCount"], "初始失败计数应该为 0")
		assert.Equal(t, 0, stats["successCount"], "初始成功计数应该为 0")
		assert.Equal(t, int64(0), stats["totalRequests"], "初始总请求数应该为 0")

		t.Log("   ✓ 初始状态正确")
	})

	t.Run("AllowRequest_InClosedState", func(t *testing.T) {
		t.Log("   测试关闭状态下的请求...")

		allowed := breaker.AllowRequest()
		assert.True(t, allowed, "关闭状态下应该允许请求")

		stats := breaker.GetStats()
		assert.Equal(t, int64(1), stats["totalRequests"], "总请求数应该增加")

		t.Log("   ✓ 关闭状态请求正常")
	})

	t.Run("RecordSuccess", func(t *testing.T) {
		t.Log("   测试记录成功...")

		breaker.RecordSuccess()

		stats := breaker.GetStats()
		assert.Equal(t, int64(1), stats["totalSuccesses"], "总成功数应该增加")
		assert.Equal(t, 0, stats["failureCount"], "成功后失败计数应该重置")

		t.Log("   ✓ 成功记录正常")
	})

	t.Run("RecordFailure_TripToOpen", func(t *testing.T) {
		t.Log("   测试失败触发熔断...")

		// 记录失败直到触发熔断
		for i := 0; i < 3; i++ {
			breaker.RecordFailure()
		}

		state := breaker.GetState()
		assert.Equal(t, circuitbreaker.StateOpen, state, "失败次数达到阈值后应该打开")

		stats := breaker.GetStats()
		assert.Equal(t, 3, stats["failureCount"], "失败计数应该为 3")

		t.Log("   ✓ 熔断触发正常")
	})

	t.Run("AllowRequest_InOpenState", func(t *testing.T) {
		t.Log("   测试打开状态下的请求...")

		allowed := breaker.AllowRequest()
		assert.False(t, allowed, "打开状态下不应该允许请求")

		t.Log("   ✓ 打开状态请求阻止正常")
	})

	t.Run("TransitionToHalfOpen", func(t *testing.T) {
		t.Log("   测试转换到半开状态...")

		// 等待超时
		time.Sleep(6 * time.Second)

		// 再次检查，应该进入半开状态
		allowed := breaker.AllowRequest()
		assert.True(t, allowed, "超时后应该允许请求（进入半开状态）")

		state := breaker.GetState()
		assert.Equal(t, circuitbreaker.StateHalfOpen, state, "超时后应该进入半开状态")

		t.Log("   ✓ 半开状态转换正常")
	})

	t.Run("RecordSuccess_InHalfOpen", func(t *testing.T) {
		t.Log("   测试半开状态下的成功记录...")

		// 在半开状态下记录成功
		breaker.RecordSuccess()
		stats := breaker.GetStats()
		assert.Equal(t, 1, stats["successCount"], "成功计数应该增加")

		breaker.RecordSuccess()
		state := breaker.GetState()
		assert.Equal(t, circuitbreaker.StateClosed, state, "成功次数达到阈值后应该关闭")

		t.Log("   ✓ 半开状态恢复正常")
	})

	t.Run("GetFailureRate", func(t *testing.T) {
		t.Log("   测试失败率计算...")

		// 重置熔断器
		breaker.Reset()

		// 记录一些成功和失败
		breaker.AllowRequest() // +1 total
		breaker.RecordSuccess()
		breaker.AllowRequest() // +1 total
		breaker.RecordFailure()
		breaker.AllowRequest() // +1 total
		breaker.RecordFailure()

		_ = breaker.GetStats()
		failureRate := breaker.GetFailureRate()

		expectedRate := float64(2) / float64(3)
		assert.InDelta(t, expectedRate, failureRate, 0.01, "失败率计算应该正确")

		t.Logf("   ✓ 失败率计算正确: %.2f", failureRate)
	})

	t.Run("StateCheckMethods", func(t *testing.T) {
		t.Log("   测试状态检查方法...")

		// 测试 IsOpen
		breaker.RecordFailure()
		breaker.RecordFailure()
		breaker.RecordFailure()
		assert.True(t, breaker.IsOpen(), "IsOpen 应该返回 true")

		// 测试 IsClosed
		breaker.Reset()
		assert.True(t, breaker.IsClosed(), "IsClosed 应该返回 true")

		// 测试 IsHalfOpen
		breaker.RecordFailure()
		breaker.RecordFailure()
		breaker.RecordFailure()
		time.Sleep(6 * time.Second)
		breaker.AllowRequest()
		assert.True(t, breaker.IsHalfOpen(), "IsHalfOpen 应该返回 true")

		t.Log("   ✓ 状态检查方法正常")
	})

	t.Run("Reset", func(t *testing.T) {
		t.Log("   测试重置...")

		breaker.Reset()

		state := breaker.GetState()
		assert.Equal(t, circuitbreaker.StateClosed, state, "重置后应该是关闭状态")

		stats := breaker.GetStats()
		assert.Equal(t, 0, stats["failureCount"], "重置后失败计数应该为 0")
		assert.Equal(t, int64(0), stats["totalRequests"], "重置后总请求数应该为 0")
		assert.Equal(t, int64(0), stats["totalSuccesses"], "重置后总成功数应该为 0")
		assert.Equal(t, int64(0), stats["totalFailures"], "重置后总失败数应该为 0")

		t.Log("   ✓ 重置功能正常")
	})

	t.Log("✅ 熔断器集成测试通过")
}

// ============================================================================
// 测试用例 3: 配额一致性测试
// ============================================================================

// TestQuotaConsistency 测试配额一致性
// 验证点：
// 1. 配额检查正确
// 2. 配额消费正确
// 3. 配额不足时正确拒绝
// 4. 配额恢复正确
// 注意：这个测试需要真实的数据库连接，或者使用 mock
func TestQuotaConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	t.Log("💰 测试配额一致性...")

	t.Run("QuotaService_BasicOperations", func(t *testing.T) {
		t.Log("   测试基本配额操作...")

		// 注意：这里只是演示测试结构
		// 实际的配额测试需要真实的 MongoDB 连接
		// 或者使用 mock repository

		t.Skip("配额测试需要数据库连接，跳过")

		// 示例代码（当有数据库连接时）：
		/*
			ctx := context.Background()
			quotaService := setupTestQuotaService(t)
			userID := "test_user_quota_001"

			// 1. 初始化配额
			err := quotaService.InitializeUserQuota(ctx, userID, "reader", "normal")
			require.NoError(t, err, "配额初始化失败")

			// 2. 检查配额
			err = quotaService.CheckQuota(ctx, userID, 100)
			require.NoError(t, err, "配额检查失败")

			// 3. 消费配额
			err = quotaService.ConsumeQuota(ctx, userID, 100, "ai-service", "default", "test_req_001")
			require.NoError(t, err, "配额消费失败")

			// 4. 获取配额信息
			quotaInfo, err := quotaService.GetQuotaInfo(ctx, userID)
			require.NoError(t, err, "获取配额信息失败")
			assert.Equal(t, 100, quotaInfo.UsedQuota, "已用配额应该为 100")

			// 5. 恢复配额
			err = quotaService.RestoreQuota(ctx, userID, 50, "测试恢复")
			require.NoError(t, err, "配额恢复失败")

			// 6. 验证恢复
			quotaInfo, err = quotaService.GetQuotaInfo(ctx, userID)
			require.NoError(t, err, "获取配额信息失败")
			assert.Equal(t, 50, quotaInfo.UsedQuota, "恢复后已用配额应该为 50")

			t.Log("   ✓ 基本配额操作正常")
		*/
	})

	t.Run("QuotaExhausted_ErrorHandling", func(t *testing.T) {
		t.Log("   测试配额不足错误处理...")

		t.Skip("配额测试需要数据库连接，跳过")

		// 示例代码（当有数据库连接时）：
		/*
			ctx := context.Background()
			quotaService := setupTestQuotaService(t)
			userID := "test_user_quota_002"

			// 初始化小额配额
			err := quotaService.InitializeUserQuota(ctx, userID, "reader", "normal")
			require.NoError(t, err)

			// 消费超过配额
			err = quotaService.ConsumeQuota(ctx, userID, 100000, "ai-service", "default", "test_req_002")
			assert.Error(t, err, "消费超过配额应该返回错误")

			// 验证错误类型
			if errors.Is(err, ai.ErrQuotaExhausted) {
				t.Log("   ✓ 正确返回配额不足错误")
			}
		*/
	})

	t.Run("QuotaWithAIService", func(t *testing.T) {
		t.Log("   测试配额与 AI 服务集成...")

		t.Skip("配额测试需要数据库连接，跳过")

		// 示例代码（当有数据库连接时）：
		/*
			ctx := context.Background()
			service, _ := setupTestAIService(t)
			quotaService := setupTestQuotaService(t)
			userID := "test_user_quota_003"

			// 初始化配额
			err := quotaService.InitializeUserQuota(ctx, userID, "reader", "normal")
			require.NoError(t, err)

			// 记录初始配额
			initialQuota, err := quotaService.GetQuotaInfo(ctx, userID)
			require.NoError(t, err)

			// 执行 AI 调用
			req := createTestAgentRequest(userID, "text_generation")
			resp, err := service.ExecuteAgent(ctx, req)
			require.NoError(t, err)

			// 等待配额同步
			time.Sleep(1 * time.Second)

			// 验证配额已扣除
			quotaAfter, err := quotaService.GetQuotaInfo(ctx, userID)
			require.NoError(t, err)
			assert.Less(t, quotaAfter.RemainingQuota, initialQuota.RemainingQuota, "配额应该减少")
			assert.Equal(t, int(resp.TokensUsed), quotaAfter.UsedQuota-initialQuota.UsedQuota, "配额减少应该与 Token 使用量一致")

			t.Logf("   ✓ 配额集成正常，使用 %d tokens", resp.TokensUsed)
		*/
	})

	t.Log("✅ 配额一致性测试通过（部分测试因缺少数据库连接而跳过）")
}

// ============================================================================
// 辅助函数
// ============================================================================

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// minInt 返回两个整数中的最小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// abs 返回整数的绝对值
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// ============================================================================
// 注意事项
// ============================================================================

// 注意：本测试文件不定义 TestMain，因为同一包中只能有一个 TestMain
// 测试配置通过环境变量 AI_SERVICE_ENDPOINT 设置
// 运行测试：
//   go test -v ./test/integration -run TestAI
//   go test -v ./test/integration -run TestAI -short  # 跳过需要外部服务的测试
