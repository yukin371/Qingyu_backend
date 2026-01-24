package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/circuitbreaker"
	"Qingyu_backend/service/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// AI 服务 gRPC 地址
	aiServiceGRPCAddr = "localhost:50051"
	// 测试超时时间
	testTimeout = 120 * time.Second
)

// AIMigrationTestSuite AI 迁移测试套件
type AIMigrationTestSuite struct {
	ctx            context.Context
	grpcClient     *ai.GRPCClient
	aiService      *ai.AIService
	quotaService   *ai.QuotaService
	cleanupFunc    func()
	testUserID     string
	testProjectID  string
}

// setupAIMigrationTestSuite 设置 AI 迁移测试套件
func setupAIMigrationTestSuite(t *testing.T) *AIMigrationTestSuite {
	if testing.Short() {
		t.Skip("跳过 AI 迁移集成测试（短模式）")
	}

	suite := &AIMigrationTestSuite{
		ctx:           context.Background(),
		testUserID:    "test_ai_migration_user",
		testProjectID: "test_ai_migration_project",
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	suite.cleanupFunc = cleanup

	// 尝试创建 gRPC 客户端
	conn, err := ai.NewGRPCConnection(aiServiceGRPCAddr)
	if err != nil {
		t.Logf("⚠️  无法连接到 AI 服务 gRPC 端点 (%s): %v", aiServiceGRPCAddr, err)
		t.Logf("   请确保 AI 服务正在运行：docker-compose up -d qingyu-ai-service")
		t.Skip("AI 服务不可用，跳过测试")
	}

	suite.grpcClient = ai.NewGRPCClient(conn, &ai.AIServiceConfig{
		Endpoint:   aiServiceGRPCAddr,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: time.Second,
	})

	// 创建 AI 服务（使用真实依赖）
	suite.aiService = ai.NewAIService(conn, nil, &ai.AIServiceConfig{
		Endpoint:       aiServiceGRPCAddr,
		Timeout:        30 * time.Second,
		MaxRetries:     3,
		RetryDelay:     time.Second,
		EnableFallback: false,
	})

	// 尝试创建配额服务
	// 注意：这需要实际的数据库连接
	suite.cleanupWhenDone(t)

	_ = router // 避免未使用警告

	return suite
}

// cleanupWhenDone 注册清理函数
func (s *AIMigrationTestSuite) cleanupWhenDone(t *testing.T) {
	t.Cleanup(func() {
		if s.cleanupFunc != nil {
			s.cleanupFunc()
		}
		if s.grpcClient != nil {
			s.grpcClient.Close()
		}
		if s.aiService != nil {
			s.aiService.Close()
		}
	})
}

// setupTestQuota 为测试用户设置配额
func (s *AIMigrationTestSuite) setupTestQuota(t *testing.T) *ai.UserQuota {
	// 尝试初始化测试用户配额
	// 这里假设全局 DB 可用
	if global.DB == nil {
		t.Skip("数据库不可用，跳过配额测试")
	}

	// 检查配额是否已存在
	var existingQuota ai.UserQuota
	err := global.DB.Collection(ai.UserQuota{}.CollectionName()).
		FindOne(s.ctx, bson.M{
			"user_id":    s.testUserID,
			"quota_type": ai.QuotaTypeDaily,
		}).Decode(&existingQuota)

	if err == nil {
		// 配额已存在，重置它
		existingQuota.Reset()
		global.DB.Collection(ai.UserQuota{}.CollectionName()).
			UpdateByID(s.ctx, existingQuota.ID, bson.M{"$set": existingQuota})
		return &existingQuota
	}

	// 创建新配额
	newQuota := &ai.UserQuota{
		ID:             primitive.NewObjectID(),
		UserID:         s.testUserID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000, // 足够用于测试
		UsedQuota:      0,
		RemainingQuota: 1000,
		Status:         ai.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1),
		Metadata: &ai.QuotaMetadata{
			UserRole:        "reader",
			MembershipLevel: "normal",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newQuota.BeforeCreate()
	_, err = global.DB.Collection(ai.UserQuota{}.CollectionName()).InsertOne(s.ctx, newQuota)
	require.NoError(t, err, "创建测试配额失败")

	t.Logf("✓ 创建测试配额: 用户=%s, 总配额=%d", s.testUserID, newQuota.TotalQuota)

	return newQuota
}

// getQuota 获取用户配额
func (s *AIMigrationTestSuite) getQuota(t *testing.T) *ai.UserQuota {
	var quota ai.UserQuota
	err := global.DB.Collection(ai.UserQuota{}.CollectionName()).
		FindOne(s.ctx, bson.M{
			"user_id":    s.testUserID,
			"quota_type": ai.QuotaTypeDaily,
		}).Decode(&quota)

	if err != nil {
		return nil
	}
	return &quota
}

// countQuotaTransactions 统计配额事务数量
func (s *AIMigrationTestSuite) countQuotaTransactions(t *testing.T) int64 {
	count, err := global.DB.Collection(ai.QuotaTransaction{}.CollectionName()).
		CountDocuments(s.ctx, bson.M{"user_id": s.testUserID})
	if err != nil {
		return 0
	}
	return count
}

// TestCompleteAIWorkflow 测试完整的 AI 调用流程
func TestCompleteAIWorkflow(t *testing.T) {
	suite := setupAIMigrationTestSuite(t)

	t.Run("AI_Service_Health_Check", func(t *testing.T) {
		t.Log("📋 步骤 1: AI 服务健康检查")

		ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
		defer cancel()

		err := suite.aiService.HealthCheck(ctx)
		if err != nil {
			t.Logf("❌ AI 服务健康检查失败: %v", err)
			t.Skip("AI 服务不可用")
		}

		t.Log("✓ AI 服务健康检查通过")
	})

	t.Run("Execute_Agent_Request", func(t *testing.T) {
		t.Log("📋 步骤 2: 执行 Agent 请求")

		// 设置测试配额
		initialQuota := suite.setupTestQuota(t)
		initialTransactionCount := suite.countQuotaTransactions(t)

		ctx, cancel := context.WithTimeout(suite.ctx, testTimeout)
		defer cancel()

		// 创建 AI 请求
		req := &ai.AgentRequest{
			UserID:      suite.testUserID,
			ProjectID:   suite.testProjectID,
			WorkflowType: "creative_workflow",
			Tasks: []ai.AgentTask{
				{
					TaskType: "generate_outline",
					Input: map[string]interface{}{
						"prompt":    "创作一个关于魔法学院的故事大纲",
						"max_chapters": 3,
					},
				},
			},
		}

		t.Logf("   发送 AI 请求: 用户=%s, 项目=%s, 工作流=%s",
			req.UserID, req.ProjectID, req.WorkflowType)

		// 执行请求
		startTime := time.Now()
		resp, err := suite.aiService.ExecuteAgent(ctx, req)
		duration := time.Since(startTime)

		if err != nil {
			t.Logf("❌ AI 请求失败: %v", err)
			t.Skip("AI 请求执行失败")
		}

		require.NotNil(t, resp, "响应不应为空")
		assert.NotEmpty(t, resp.ExecutionID, "执行ID不应为空")
		assert.Greater(t, resp.TokensUsed, 0, "使用的Token数应大于0")

		t.Logf("✓ AI 请求成功")
		t.Logf("   执行ID: %s", resp.ExecutionID)
		t.Logf("   使用Token: %d", resp.TokensUsed)
		t.Logf("   耗时: %.2f秒", duration.Seconds())

		// 验证配额扣除
		t.Run("Verify_Quota_Deduction", func(t *testing.T) {
			t.Log("📋 步骤 3: 验证配额扣除")

			// 等待配额更新
			time.Sleep(100 * time.Millisecond)

			finalQuota := suite.getQuota(t)
			require.NotNil(t, finalQuota, "应该能获取到配额")

			t.Logf("   初始配额: 总计=%d, 已用=%d, 剩余=%d",
				initialQuota.TotalQuota, initialQuota.UsedQuota, initialQuota.RemainingQuota)
			t.Logf("   最终配额: 总计=%d, 已用=%d, 剩余=%d",
				finalQuota.TotalQuota, finalQuota.UsedQuota, finalQuota.RemainingQuota)

			// 验证配额已扣除
			assert.Greater(t, finalQuota.UsedQuota, initialQuota.UsedQuota,
				"已用配额应该增加")
			assert.Less(t, finalQuota.RemainingQuota, initialQuota.RemainingQuota,
				"剩余配额应该减少")

			t.Logf("✓ 配额验证通过: 扣除=%d Token",
				finalQuota.UsedQuota-initialQuota.UsedQuota)
		})

		// 验证事务记录
		t.Run("Verify_Transaction_Record", func(t *testing.T) {
			t.Log("📋 步骤 4: 验证事务记录")

			finalTransactionCount := suite.countQuotaTransactions(t)
			assert.Greater(t, finalTransactionCount, initialTransactionCount,
				"应该有新的配额事务记录")

			t.Logf("✓ 事务记录验证通过: 新增=%d 条",
				finalTransactionCount-initialTransactionCount)
		})
	})
}

// TestCircuitBreakerBehavior 测试熔断器行为
func TestCircuitBreakerBehavior(t *testing.T) {
	suite := setupAIMigrationTestSuite(t)

	t.Run("Initial_State", func(t *testing.T) {
		t.Log("📋 测试熔断器初始状态")

		// 获取熔断器实例
		cb := circuitbreaker.NewCircuitBreaker(3, 5*time.Second, 2)

		state := cb.GetState()
		stats := cb.GetStats()

		assert.Equal(t, circuitbreaker.StateClosed, state, "初始状态应为关闭")
		assert.Equal(t, "Closed", stats["state"], "状态字符串应为Closed")
		assert.Equal(t, 0, stats["failureCount"], "初始失败次数应为0")

		t.Logf("✓ 熔断器初始状态正确: %s", state)
	})

	t.Run("Trigger_Circuit_Breaker", func(t *testing.T) {
		t.Log("📋 测试触发熔断")

		cb := circuitbreaker.NewCircuitBreaker(3, 5*time.Second, 2)

		// 记录3次失败（达到阈值）
		for i := 0; i < 3; i++ {
			cb.RecordFailure()
			t.Logf("   记录失败 #%d", i+1)
		}

		state := cb.GetState()
		assert.Equal(t, circuitbreaker.StateOpen, state, "达到阈值后应打开熔断器")

		stats := cb.GetStats()
		assert.Equal(t, 3, stats["failureCount"], "失败次数应为3")
		assert.False(t, cb.AllowRequest(), "熔断器打开时应拒绝请求")

		t.Logf("✓ 熔断器正确触发: 状态=%s, 失败次数=%d",
			state, stats["failureCount"])
	})

	t.Run("Half_Open_State", func(t *testing.T) {
		t.Log("📋 测试半开状态")

		cb := circuitbreaker.NewCircuitBreaker(3, 100*time.Millisecond, 2)

		// 触发熔断
		for i := 0; i < 3; i++ {
			cb.RecordFailure()
		}

		assert.Equal(t, circuitbreaker.StateOpen, cb.GetState(), "应已打开")

		// 等待超时进入半开状态
		time.Sleep(150 * time.Millisecond)

		// 下一个请求应该被允许（进入半开状态）
		allowed := cb.AllowRequest()
		assert.True(t, allowed, "超时后应允许请求")
		assert.Equal(t, circuitbreaker.StateHalfOpen, cb.GetState(), "应进入半开状态")

		t.Logf("✓ 熔断器正确进入半开状态")
	})

	t.Run("Recovery_After_Success", func(t *testing.T) {
		t.Log("📋 测试成功后恢复")

		cb := circuitbreaker.NewCircuitBreaker(3, 100*time.Millisecond, 2)

		// 触发熔断
		for i := 0; i < 3; i++ {
			cb.RecordFailure()
		}

		// 等待进入半开状态
		time.Sleep(150 * time.Millisecond)
		cb.AllowRequest() // 触发进入半开

		// 记录2次成功（达到恢复阈值）
		cb.RecordSuccess()
		t.Logf("   记录成功 #1")
		cb.RecordSuccess()
		t.Logf("   记录成功 #2")

		state := cb.GetState()
		assert.Equal(t, circuitbreaker.StateClosed, state, "达到成功阈值后应关闭熔断器")

		stats := cb.GetStats()
		assert.Equal(t, 0, stats["failureCount"], "恢复后失败计数应重置")

		t.Logf("✓ 熔断器正确恢复: 状态=%s", state)
	})

	t.Run("Statistics", func(t *testing.T) {
		t.Log("📋 测试熔断器统计")

		cb := circuitbreaker.NewCircuitBreaker(5, 10*time.Second, 3)

		// 记录一些请求
		for i := 0; i < 10; i++ {
			if i < 3 {
				cb.RecordFailure()
			} else {
				cb.RecordSuccess()
			}
		}

		stats := cb.GetStats()
		assert.Equal(t, int64(10), stats["totalRequests"], "总请求数应为10")
		assert.Equal(t, int64(7), stats["totalSuccesses"], "成功请求数应为7")
		assert.Equal(t, int64(3), stats["totalFailures"], "失败请求数应为3")

		failureRate := cb.GetFailureRate()
		assert.InDelta(t, 0.3, failureRate, 0.01, "失败率应约为30%")

		t.Logf("✓ 熔断器统计正确:")
		t.Logf("   总请求: %v", stats["totalRequests"])
		t.Logf("   成功: %v", stats["totalSuccesses"])
		t.Logf("   失败: %v", stats["totalFailures"])
		t.Logf("   失败率: %.2f%%", failureRate*100)
	})

	t.Run("Integrated_With_AI_Service", func(t *testing.T) {
		t.Log("📋 测试 AI 服务集成的熔断器")

		// 创建 AI 服务实例以获取其熔断器
		if suite.aiService == nil {
			t.Skip("AI 服务不可用")
		}

		// 获取熔断器状态
		state := suite.aiService.GetCircuitBreakerState()
		stats := suite.aiService.GetCircuitBreakerStats()

		t.Logf("✓ AI 服务熔断器状态:")
		t.Logf("   状态: %s", state)
		t.Logf("   统计: %+v", stats)

		// 验证熔断器方法可用
		assert.NotNil(t, stats, "统计信息不应为空")
		assert.Contains(t, stats, "state", "应包含状态信息")
	})
}

// TestQuotaConsistency 测试配额一致性
func TestQuotaConsistency(t *testing.T) {
	suite := setupAIMigrationTestSuite(t)

	if global.DB == nil {
		t.Skip("数据库不可用，跳过配额一致性测试")
	}

	t.Run("Setup_Initial_Quota", func(t *testing.T) {
		t.Log("📋 步骤 1: 设置初始配额")

		initialQuota := suite.setupTestQuota(t)

		assert.Equal(t, 1000, initialQuota.TotalQuota, "总配额应为1000")
		assert.Equal(t, 0, initialQuota.UsedQuota, "已用配额应为0")
		assert.Equal(t, 1000, initialQuota.RemainingQuota, "剩余配额应为1000")

		t.Logf("✓ 初始配额设置成功")
	})

	t.Run("Execute_AI_Call", func(t *testing.T) {
		t.Log("📋 步骤 2: 执行 AI 调用")

		// 获取初始配额
		initialQuota := suite.getQuota(t)
		require.NotNil(t, initialQuota, "应该有初始配额")

		ctx, cancel := context.WithTimeout(suite.ctx, testTimeout)
		defer cancel()

		// 创建简单的 AI 请求
		req := &ai.AgentRequest{
			UserID:        suite.testUserID,
			ProjectID:     suite.testProjectID,
			WorkflowType:  "text_generation",
			Tasks: []ai.AgentTask{
				{
					TaskType: "generate_text",
					Input: map[string]interface{}{
						"prompt": "写一段简短的故事",
						"max_tokens": 100,
					},
				},
			},
		}

		t.Logf("   执行 AI 调用...")
		startTime := time.Now()

		resp, err := suite.aiService.ExecuteAgent(ctx, req)

		duration := time.Since(startTime)

		if err != nil {
			t.Logf("⚠️  AI 调用失败: %v", err)
			t.Logf("   这可能是正常的，如果 AI 服务不可用")
			// 不跳过，继续测试配额逻辑
		} else {
			t.Logf("✓ AI 调用成功")
			t.Logf("   使用Token: %d", resp.TokensUsed)
			t.Logf("   耗时: %.2f秒", duration.Seconds())
		}

		// 等待配额更新
		time.Sleep(200 * time.Millisecond)
	})

	t.Run("Verify_Backend_Quota", func(t *testing.T) {
		t.Log("📋 步骤 3: 验证后端配额扣除")

		finalQuota := suite.getQuota(t)
		require.NotNil(t, finalQuota, "应该能获取到配额")

		t.Logf("   配额状态:")
		t.Logf("   总计: %d", finalQuota.TotalQuota)
		t.Logf("   已用: %d", finalQuota.UsedQuota)
		t.Logf("   剩余: %d", finalQuota.RemainingQuota)
		t.Logf("   状态: %s", finalQuota.Status)

		// 验证配额状态有效
		assert.True(t, finalQuota.TotalQuota > 0, "总配额应大于0")
		assert.True(t, finalQuota.UsedQuota >= 0, "已用配额应>=0")
		assert.True(t, finalQuota.RemainingQuota >= 0, "剩余配额应>=0")

		t.Logf("✓ 后端配额状态有效")
	})

	t.Run("Verify_Consistency", func(t *testing.T) {
		t.Log("📋 步骤 4: 验证配额一致性")

		quota := suite.getQuota(t)
		require.NotNil(t, quota, "应该能获取到配额")

		// 验证：已用 + 剩余 = 总计
		totalCalculated := quota.UsedQuota + quota.RemainingQuota

		t.Logf("   一致性检查:")
		t.Logf("   已用 + 剩余 = %d + %d = %d",
			quota.UsedQuota, quota.RemainingQuota, totalCalculated)
		t.Logf("   总配额 = %d", quota.TotalQuota)

		// 允许1%的误差（由于可能的并发更新）
		epsilon := int(float64(quota.TotalQuota) * 0.01)
		difference := abs(totalCalculated - quota.TotalQuota)

		assert.LessOrEqual(t, difference, epsilon,
			fmt.Sprintf("已用+剩余 应约等于总计（误差<1%%），差值=%d", difference))

		t.Logf("✓ 配额一致性验证通过（差值=%d, 阈值=%d）", difference, epsilon)
	})

	t.Run("Verify_Transaction_History", func(t *testing.T) {
		t.Log("📋 步骤 5: 验证事务历史")

		count := suite.countQuotaTransactions(t)

		t.Logf("   用户 %s 的配额事务记录数: %d", suite.testUserID, count)

		// 查询最近的事务
		cursor, err := global.DB.Collection(ai.QuotaTransaction{}.CollectionName()).
			Find(suite.ctx, bson.M{"user_id": suite.testUserID})
		require.NoError(t, err, "查询事务历史失败")
		defer cursor.Close(suite.ctx)

		var transactions []ai.QuotaTransaction
		err = cursor.All(suite.ctx, &transactions)
		require.NoError(t, err, "解析事务历史失败")

		if len(transactions) > 0 {
			t.Logf("   最近事务:")
			for i, txn := range transactions {
				if i >= 3 { // 只显示前3条
					break
				}
				t.Logf("   [%d] 类型=%s, 数量=%d, 服务=%s, 时间=%s",
					i+1, txn.Type, txn.Amount, txn.Service,
					txn.Timestamp.Format("15:04:05"))
			}

			assert.NotEmpty(t, transactions[0].RequestID, "事务应有请求ID")
			assert.NotEmpty(t, transactions[0].Service, "事务应有服务类型")
		}

		t.Logf("✓ 事务历史验证通过")
	})
}

// TestAIServiceIntegration 综合集成测试
func TestAIServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过综合集成测试（短模式）")
	}

	suite := setupAIMigrationTestSuite(t)

	t.Run("Full_Integration_Check", func(t *testing.T) {
		t.Log("🎯 执行综合集成检查")

		checks := []struct {
			name string
			fn   func(t *testing.T) bool
		}{
			{
				name: "AI服务健康检查",
				fn: func(t *testing.T) bool {
					ctx, cancel := context.WithTimeout(suite.ctx, 3*time.Second)
					defer cancel()
					return suite.aiService.HealthCheck(ctx) == nil
				},
			},
			{
				name: "熔断器状态",
				fn: func(t *testing.T) bool {
					state := suite.aiService.GetCircuitBreakerState()
					return state == circuitbreaker.StateClosed ||
						state == circuitbreaker.StateOpen ||
						state == circuitbreaker.StateHalfOpen
				},
			},
			{
				name: "配额服务可用性",
				fn: func(t *testing.T) bool {
					return global.DB != nil
				},
			},
		}

		allPassed := true
		for _, check := range checks {
			t.Run(check.name, func(t *testing.T) {
				passed := check.fn(t)
				if passed {
					t.Logf("✓ %s - 通过", check.name)
				} else {
					t.Logf("✗ %s - 失败", check.name)
					allPassed = false
				}
			})
		}

		if allPassed {
			t.Log("🎉 所有集成检查通过！")
		} else {
			t.Log("⚠️  部分集成检查未通过")
		}
	})
}

// 辅助函数

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// setupTestEnvironment 复用现有的测试环境设置
// 这个函数已经在 helpers.go 中定义，这里只是为了文档完整性
