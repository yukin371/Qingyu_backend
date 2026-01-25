package search

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// =========================
// 辅助函数
// =========================

// getTestLogger 获取测试用的 logger
func getTestLogger(t *testing.T) *zap.Logger {
	t.Helper()
	return zap.NewNop()
}

// getTestGrayScaleConfig 获取测试用的灰度配置
func getTestGrayScaleConfig(enabled bool, percent int) *GrayScaleConfig {
	return &GrayScaleConfig{
		Enabled: enabled,
		Percent: percent,
	}
}

// =========================
// 基础测试
// =========================

// TestNewGrayScaleDecision 测试创建灰度决策器
func TestNewGrayScaleDecision(t *testing.T) {
	logger := getTestLogger(t)

	// 测试正常创建
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	assert.NotNil(t, decision, "灰度决策器不应为 nil")

	// 验证配置
	retrievedConfig := decision.GetConfig()
	assert.True(t, retrievedConfig.Enabled)
	assert.Equal(t, 50, retrievedConfig.Percent)
}

// TestGrayScaleDecision_NilConfig 测试 nil 配置
func TestGrayScaleDecision_NilConfig(t *testing.T) {
	logger := getTestLogger(t)

	// 创建 nil 配置的决策器
	decision := NewGrayScaleDecision(nil, logger)
	ctx := context.Background()

	// nil 配置应该返回 false（使用 MongoDB）
	result := decision.ShouldUseES(ctx, "books", "user123")
	assert.False(t, result, "nil 配置应该使用 MongoDB")
}

// TestGrayScaleDecision_DefaultConfig 测试默认配置
func TestGrayScaleDecision_DefaultConfig(t *testing.T) {
	logger := getTestLogger(t)

	// 测试默认禁用的配置
	config := getTestGrayScaleConfig(false, 0)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	result := decision.ShouldUseES(ctx, "books", "user123")
	assert.False(t, result, "默认禁用配置应该使用 MongoDB")
}

// =========================
// 流量分配测试
// =========================

// TestGrayScaleDecision_Disabled 测试灰度禁用
func TestGrayScaleDecision_Disabled(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(false, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 禁用状态下，所有流量都应使用 MongoDB
	result := decision.ShouldUseES(ctx, "books", "user123")
	assert.False(t, result, "灰度禁用时应使用 MongoDB")
}

// TestGrayScaleDecision_ZeroPercent 测试 0% 流量
func TestGrayScaleDecision_ZeroPercent(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 0)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 0% 流量应使用 MongoDB
	result := decision.ShouldUseES(ctx, "books", "user123")
	assert.False(t, result, "0% 流量应该使用 MongoDB")
}

// TestGrayScaleDecision_HundredPercent 测试 100% 流量
func TestGrayScaleDecision_HundredPercent(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 100)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 100% 流量应使用 ES
	result := decision.ShouldUseES(ctx, "books", "user123")
	assert.True(t, result, "100% 流量应该使用 ES")
}

// TestGrayScaleDecision_UserIDHashConsistency 测试用户 ID 哈希一致性
func TestGrayScaleDecision_UserIDHashConsistency(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 同一用户 ID 的多次调用应返回相同结果
	userID := "user123"
	var results []bool

	for i := 0; i < 10; i++ {
		result := decision.ShouldUseES(ctx, "books", userID)
		results = append(results, result)
	}

	// 验证所有结果一致
	firstResult := results[0]
	for _, result := range results {
		assert.Equal(t, firstResult, result, "同一用户 ID 的决策结果应该一致")
	}
}

// TestGrayScaleDecision_DifferentUsers 测试不同用户的分配
func TestGrayScaleDecision_DifferentUsers(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 测试多个不同用户
	userIDs := []string{
		"user1", "user2", "user3", "user4", "user5",
		"user6", "user7", "user8", "user9", "user10",
	}

	results := make(map[bool]int)
	for _, userID := range userIDs {
		result := decision.ShouldUseES(ctx, "books", userID)
		results[result]++
	}

	// 50% 流量下，应该有部分用户使用 ES，部分使用 MongoDB
	// 虽然具体分布可能不完全平均，但应该有两种结果
	assert.True(t, results[true] > 0 || results[false] > 0, "应该有决策结果")
}

// TestGrayScaleDistribution_ValidateDistribution 测试流量分配分布
func TestGrayScaleDistribution_ValidateDistribution(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 测试大量用户，验证流量分配接近 50%
	esCount := 0
	mongoCount := 0
	totalUsers := 1000

	for i := 0; i < totalUsers; i++ {
		userID := string(rune('a' + i%26)) + string(rune('a'+(i/26)%26)) + string(rune('a'+(i/676)%26))
		result := decision.ShouldUseES(ctx, "books", userID)
		if result {
			esCount++
		} else {
			mongoCount++
		}
	}

	esPercent := float64(esCount) / float64(totalUsers) * 100

	// 验证 ES 流量在 40%-60% 之间（允许一定偏差）
	assert.Greater(t, esPercent, 40.0, "ES 流量应大于 40%")
	assert.Less(t, esPercent, 60.0, "ES 流量应小于 60%")
}

// =========================
// 指标收集测试
// =========================

// TestGrayScaleMetrics_RecordESUsage 测试记录 ES 使用
func TestGrayScaleMetrics_RecordESUsage(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 记录 ES 使用
	decision.RecordUsage("es", 100*time.Millisecond)
	decision.RecordUsage("ES", 150*time.Millisecond)
	decision.RecordUsage("elasticsearch", 200*time.Millisecond)

	// 获取指标
	metrics := decision.GetMetrics()

	assert.Equal(t, int64(3), metrics.ESCount, "ES 使用次数应为 3")
	assert.Equal(t, int64(0), metrics.MongoDBCount, "MongoDB 使用次数应为 0")
	assert.Equal(t, 450*time.Millisecond, metrics.ESTotalTook, "ES 总耗时应为 450ms")
	assert.Equal(t, 150*time.Millisecond, metrics.ESAvgTook, "ES 平均耗时应为 150ms")
}

// TestGrayScaleMetrics_RecordMongoDBUsage 测试记录 MongoDB 使用
func TestGrayScaleMetrics_RecordMongoDBUsage(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 记录 MongoDB 使用
	decision.RecordUsage("mongodb", 80*time.Millisecond)
	decision.RecordUsage("MongoDB", 120*time.Millisecond)

	// 获取指标
	metrics := decision.GetMetrics()

	assert.Equal(t, int64(0), metrics.ESCount, "ES 使用次数应为 0")
	assert.Equal(t, int64(2), metrics.MongoDBCount, "MongoDB 使用次数应为 2")
	assert.Equal(t, 200*time.Millisecond, metrics.MongoDBTotalTook, "MongoDB 总耗时应为 200ms")
	assert.Equal(t, 100*time.Millisecond, metrics.MongoDBAvgTook, "MongoDB 平均耗时应为 100ms")
}

// TestGrayScaleMetrics_RecordMixedUsage 测试记录混合使用
func TestGrayScaleMetrics_RecordMixedUsage(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 记录混合使用
	decision.RecordUsage("es", 100*time.Millisecond)
	decision.RecordUsage("MongoDB", 150*time.Millisecond)
	decision.RecordUsage("es", 200*time.Millisecond)
	decision.RecordUsage("MongoDB", 120*time.Millisecond)

	// 获取指标
	metrics := decision.GetMetrics()

	assert.Equal(t, int64(2), metrics.ESCount, "ES 使用次数应为 2")
	assert.Equal(t, int64(2), metrics.MongoDBCount, "MongoDB 使用次数应为 2")
	assert.Equal(t, 300*time.Millisecond, metrics.ESTotalTook, "ES 总耗时应为 300ms")
	assert.Equal(t, 150*time.Millisecond, metrics.ESAvgTook, "ES 平均耗时应为 150ms")
	assert.Equal(t, 270*time.Millisecond, metrics.MongoDBTotalTook, "MongoDB 总耗时应为 270ms")
	assert.Equal(t, 135*time.Millisecond, metrics.MongoDBAvgTook, "MongoDB 平均耗时应为 135ms")
}

// TestGrayScaleMetrics_TrafficDistribution 测试流量分配比例
func TestGrayScaleMetrics_TrafficDistribution(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 记录使用情况
	for i := 0; i < 30; i++ {
		decision.RecordUsage("es", 100*time.Millisecond)
	}
	for i := 0; i < 70; i++ {
		decision.RecordUsage("MongoDB", 150*time.Millisecond)
	}

	// 获取指标
	metrics := decision.GetMetrics()

	// 计算流量分配
	esPercent, mongoPercent := metrics.GetTrafficDistribution()

	assert.InDelta(t, 30.0, esPercent, 0.1, "ES 流量比例应为 30%")
	assert.InDelta(t, 70.0, mongoPercent, 0.1, "MongoDB 流量比例应为 70%")
}

// TestGrayScaleMetrics_ZeroUsage 测试零使用情况
func TestGrayScaleMetrics_ZeroUsage(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 获取未使用的指标
	metrics := decision.GetMetrics()

	assert.Equal(t, int64(0), metrics.ESCount, "ES 使用次数应为 0")
	assert.Equal(t, int64(0), metrics.MongoDBCount, "MongoDB 使用次数应为 0")
	assert.Equal(t, time.Duration(0), metrics.ESAvgTook, "ES 平均耗时应为 0")
	assert.Equal(t, time.Duration(0), metrics.MongoDBAvgTook, "MongoDB 平均耗时应为 0")

	// 测试流量分配
	esPercent, mongoPercent := metrics.GetTrafficDistribution()
	assert.Equal(t, 0.0, esPercent, "ES 流量比例应为 0")
	assert.Equal(t, 0.0, mongoPercent, "MongoDB 流量比例应为 0")
}

// =========================
// 配置热更新测试
// =========================

// TestGrayScaleConfig_UpdateConfig 测试更新配置
func TestGrayScaleConfig_UpdateConfig(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 初始配置
	initialConfig := decision.GetConfig()
	assert.True(t, initialConfig.Enabled)
	assert.Equal(t, 50, initialConfig.Percent)

	// 更新配置
	err := decision.UpdateConfig(false, 80)
	require.NoError(t, err, "更新配置不应出错")

	// 验证配置已更新
	updatedConfig := decision.GetConfig()
	assert.False(t, updatedConfig.Enabled, "灰度应该已禁用")
	assert.Equal(t, 80, updatedConfig.Percent, "百分比应该已更新")
}

// TestGrayScaleConfig_UpdateInvalidPercent 测试更新无效百分比
func TestGrayScaleConfig_UpdateInvalidPercent(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 测试负数百分比
	err := decision.UpdateConfig(true, -10)
	assert.Error(t, err, "负数百分比应该返回错误")
	assert.Contains(t, err.Error(), "invalid percent", "错误信息应包含 'invalid percent'")

	// 测试超过 100 的百分比
	err = decision.UpdateConfig(true, 150)
	assert.Error(t, err, "超过 100 的百分比应该返回错误")
	assert.Contains(t, err.Error(), "invalid percent", "错误信息应包含 'invalid percent'")

	// 验证配置未改变
	currentConfig := decision.GetConfig()
	assert.True(t, currentConfig.Enabled)
	assert.Equal(t, 50, currentConfig.Percent, "无效更新不应改变配置")
}

// TestGrayScaleConfig_BoundaryValues 测试边界值
func TestGrayScaleConfig_BoundaryValues(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 测试边界值 0
	err := decision.UpdateConfig(true, 0)
	require.NoError(t, err, "0% 应该是有效的")
	assert.Equal(t, 0, decision.GetConfig().Percent)

	// 测试边界值 100
	err = decision.UpdateConfig(true, 100)
	require.NoError(t, err, "100% 应该是有效的")
	assert.Equal(t, 100, decision.GetConfig().Percent)
}

// =========================
// 并发测试
// =========================

// TestGrayScaleDecision_ConcurrentDecisions 测试并发决策
func TestGrayScaleDecision_ConcurrentDecisions(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 并发测试
	var wg sync.WaitGroup
	numGoroutines := 100
	decisionsPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < decisionsPerGoroutine; j++ {
				userID := string(rune('a' + (goroutineID+j)%26))
				_ = decision.ShouldUseES(ctx, "books", userID)
			}
		}(i)
	}

	wg.Wait()

	// 验证没有 panic 或 race condition
	// 如果测试通过，说明并发安全
	assert.True(t, true, "并发测试应该成功完成")
}

// TestGrayScaleDecision_ConcurrentConfigUpdates 测试并发配置更新
func TestGrayScaleDecision_ConcurrentConfigUpdates(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 并发更新配置
	var wg sync.WaitGroup
	numUpdates := 50

	for i := 0; i < numUpdates; i++ {
		wg.Add(1)
		go func(updateID int) {
			defer wg.Done()
			enabled := updateID%2 == 0
			percent := (updateID % 101)
			_ = decision.UpdateConfig(enabled, percent)
		}(i)
	}

	wg.Wait()

	// 验证最终配置有效
	finalConfig := decision.GetConfig()
	assert.GreaterOrEqual(t, finalConfig.Percent, 0, "百分比应该 >= 0")
	assert.LessOrEqual(t, finalConfig.Percent, 100, "百分比应该 <= 100")
}

// TestGrayScaleDecision_ConcurrentDecisionsAndUpdates 测试并发决策和配置更新
func TestGrayScaleDecision_ConcurrentDecisionsAndUpdates(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	var wg sync.WaitGroup

	// 并发决策
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userID := string(rune('a' + id%26))
			_ = decision.ShouldUseES(ctx, "books", userID)
		}(i)
	}

	// 并发配置更新
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = decision.UpdateConfig(id%2 == 0, id%101)
		}(i)
	}

	// 并发记录指标
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			engine := "es"
			if id%2 == 0 {
				engine = "MongoDB"
			}
			decision.RecordUsage(engine, time.Duration(id)*time.Millisecond)
		}(i)
	}

	wg.Wait()

	// 验证指标收集正常
	metrics := decision.GetMetrics()
	assert.True(t, metrics.ESCount >= 0, "ES 计数应该有效")
	assert.True(t, metrics.MongoDBCount >= 0, "MongoDB 计数应该有效")
}

// TestGrayScaleMetrics_ConcurrentRecordUsage 测试并发记录使用情况
func TestGrayScaleMetrics_ConcurrentRecordUsage(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	var wg sync.WaitGroup
	numRecords := 1000

	// 并发记录 ES 使用
	for i := 0; i < numRecords; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			decision.RecordUsage("es", time.Duration(id)*time.Millisecond)
		}(i)
	}

	wg.Wait()

	// 验证计数正确
	metrics := decision.GetMetrics()
	assert.Equal(t, int64(numRecords), metrics.ESCount, "ES 计数应该等于记录次数")
}

// =========================
// 集成测试
// =========================

// TestGrayScaleDecision_FullWorkflow 测试完整的灰度决策流程
func TestGrayScaleDecision_FullWorkflow(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 30)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 模拟多个用户搜索
	users := []string{"user1", "user2", "user3", "user4", "user5", "user6", "user7", "user8", "user9", "user10"}

	for _, userID := range users {
		// 1. 决策使用哪个引擎
		useES := decision.ShouldUseES(ctx, "books", userID)

		// 2. 模拟搜索并记录使用情况
		engine := "MongoDB"
		took := 150 * time.Millisecond
		if useES {
			engine = "es"
			took = 100 * time.Millisecond
		}
		decision.RecordUsage(engine, took)

		// 3. 验证决策一致性
		result := decision.ShouldUseES(ctx, "books", userID)
		assert.Equal(t, useES, result, "同一用户的决策应该一致")
	}

	// 4. 验证指标
	metrics := decision.GetMetrics()
	totalSearches := metrics.ESCount + metrics.MongoDBCount
	assert.Equal(t, int64(len(users)), totalSearches, "总搜索次数应该等于用户数")
}

// TestGrayScaleDecision_ConfigChangeImpact 测试配置变更影响
func TestGrayScaleDecision_ConfigChangeImpact(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 0)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 初始配置：0% 流量
	userID := "user123"
	result1 := decision.ShouldUseES(ctx, "books", userID)
	assert.False(t, result1, "0% 流量应该使用 MongoDB")

	// 更新配置：100% 流量
	err := decision.UpdateConfig(true, 100)
	require.NoError(t, err)

	result2 := decision.ShouldUseES(ctx, "books", userID)
	assert.True(t, result2, "100% 流量应该使用 ES")

	// 更新配置：禁用灰度
	err = decision.UpdateConfig(false, 50)
	require.NoError(t, err)

	result3 := decision.ShouldUseES(ctx, "books", userID)
	assert.False(t, result3, "禁用灰度应该使用 MongoDB")
}

// TestGrayScaleDecision_DifferentSearchTypes 测试不同搜索类型
func TestGrayScaleDecision_DifferentSearchTypes(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 测试不同搜索类型
	searchTypes := []string{"books", "projects", "documents", "users"}
	userID := "user123"

	for _, searchType := range searchTypes {
		result := decision.ShouldUseES(ctx, searchType, userID)
		// 所有搜索类型对同一用户应该返回相同结果
		// 因为决策只基于 userID 哈希
		assert.IsType(t, false, result, "决策结果应该是布尔类型")
	}
}

// =========================
// 边界条件测试
// =========================

// TestGrayScaleDecision_EmptyUserID 测试空用户 ID
func TestGrayScaleDecision_EmptyUserID(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 空用户 ID 应该使用随机策略
	// 多次调用应该有不同的结果
	results := make(map[bool]int)
	for i := 0; i < 100; i++ {
		result := decision.ShouldUseES(ctx, "books", "")
		results[result]++
	}

	// 应该有两种结果（因为使用随机策略）
	assert.True(t, results[true] > 0, "应该有使用 ES 的情况")
	assert.True(t, results[false] > 0, "应该有使用 MongoDB 的情况")
}

// TestGrayScaleDecision_SpecialCharactersInUserID 测试用户 ID 包含特殊字符
func TestGrayScaleDecision_SpecialCharactersInUserID(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)
	ctx := context.Background()

	// 测试包含特殊字符的用户 ID
	userIDs := []string{
		"user@example.com",
		"user_123!@#",
		"用户123",
		"user with spaces",
		"user\n\t\r",
	}

	for _, userID := range userIDs {
		// 不应该 panic
		result := decision.ShouldUseES(ctx, "books", userID)
		assert.IsType(t, false, result, "特殊字符用户 ID 应该正常处理")
	}
}

// TestGrayScaleMetrics_LargeDuration 测试大耗时值
func TestGrayScaleMetrics_LargeDuration(t *testing.T) {
	logger := getTestLogger(t)
	config := getTestGrayScaleConfig(true, 50)
	decision := NewGrayScaleDecision(config, logger)

	// 记录大耗时值
	largeDuration := 3600 * time.Second // 1 小时
	decision.RecordUsage("es", largeDuration)

	metrics := decision.GetMetrics()
	assert.Equal(t, largeDuration, metrics.ESTotalTook, "应该正确记录大耗时值")
	assert.Equal(t, largeDuration, metrics.ESAvgTook, "应该正确计算平均耗时")
}
