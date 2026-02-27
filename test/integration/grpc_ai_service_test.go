package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"Qingyu_backend/service/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	grpcServerAddr = "localhost:50051"
	testTimeout    = 120 * time.Second
)

func requireAIService(t *testing.T) *ai.Phase3Client {
	t.Helper()

	client, err := ai.NewPhase3Client(grpcServerAddr)
	if err != nil {
		t.Skipf("AI gRPC服务不可用，跳过测试: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := client.HealthCheck(ctx); err != nil {
		client.Close()
		t.Skipf("AI gRPC健康检查失败，跳过测试: %v", err)
	}

	return client
}

// TestGRPCConnection 测试gRPC连接
func TestGRPCConnection(t *testing.T) {
	client := requireAIService(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.HealthCheck(ctx)
	require.NoError(t, err, "健康检查失败")
	assert.Equal(t, "healthy", resp.Status, "服务状态不健康")

	fmt.Printf("✅ gRPC连接成功 - 状态: %s\n", resp.Status)
	for name, status := range resp.Checks {
		fmt.Printf("   - %s: %s\n", name, status)
	}
}

// TestGenerateOutline 测试大纲生成
func TestGenerateOutline(t *testing.T) {
	client := requireAIService(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	task := "创作一个修仙小说大纲，主角是天才少年，包含5章内容"

	fmt.Printf("\n📖 测试大纲生成\n")
	fmt.Printf("   任务: %s\n", task)

	startTime := time.Now()
	resp, err := client.GenerateOutline(ctx, task, "test_user", "test_project", nil)
	duration := time.Since(startTime)

	require.NoError(t, err, "大纲生成失败")
	require.NotNil(t, resp, "大纲响应为空")
	require.NotNil(t, resp.Outline, "大纲数据为空")

	// 验证大纲数据
	assert.NotEmpty(t, resp.Outline.Title, "大纲标题为空")
	assert.NotEmpty(t, resp.Outline.Genre, "大纲类型为空")
	assert.Greater(t, len(resp.Outline.Chapters), 0, "章节数量为0")

	fmt.Printf("\n✅ 大纲生成成功! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   📖 标题: %s\n", resp.Outline.Title)
	fmt.Printf("   🎭 类型: %s\n", resp.Outline.Genre)
	fmt.Printf("   📚 章节数: %d\n", len(resp.Outline.Chapters))

	if len(resp.Outline.Chapters) > 0 {
		fmt.Println("   章节列表:")
		for i, chapter := range resp.Outline.Chapters {
			if i < 3 {
				fmt.Printf("     %d. %s\n", i+1, chapter.Title)
			}
		}
		if len(resp.Outline.Chapters) > 3 {
			fmt.Printf("     ... 还有 %d 章\n", len(resp.Outline.Chapters)-3)
		}
	}
}

// TestGenerateCharacters 测试角色生成
func TestGenerateCharacters(t *testing.T) {
	client := requireAIService(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// 先生成大纲
	outlineResp, err := client.GenerateOutline(
		ctx,
		"创作一个修仙小说大纲",
		"test_user",
		"test_project",
		nil,
	)
	require.NoError(t, err, "大纲生成失败")

	fmt.Printf("\n👤 测试角色生成\n")

	// 生成角色
	startTime := time.Now()
	resp, err := client.GenerateCharacters(
		ctx,
		"根据大纲创建主要角色",
		"test_user",
		"test_project",
		outlineResp.Outline,
		nil,
	)
	duration := time.Since(startTime)

	require.NoError(t, err, "角色生成失败")
	require.NotNil(t, resp, "角色响应为空")
	require.NotNil(t, resp.Characters, "角色数据为空")
	assert.Greater(t, len(resp.Characters.Characters), 0, "角色数量为0")

	fmt.Printf("\n✅ 角色生成成功! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   👥 角色数量: %d\n", len(resp.Characters.Characters))

	if len(resp.Characters.Characters) > 0 {
		fmt.Println("   角色列表:")
		for i, char := range resp.Characters.Characters {
			if i < 3 {
				fmt.Printf("     %d. %s (%s)\n", i+1, char.Name, char.RoleType)
				if char.Personality != nil && len(char.Personality.Traits) > 0 {
					fmt.Printf("        性格: %v\n", char.Personality.Traits[:min(3, len(char.Personality.Traits))])
				}
			}
		}
	}
}

// TestGeneratePlot 测试情节生成
func TestGeneratePlot(t *testing.T) {
	client := requireAIService(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// 先生成大纲和角色
	outlineResp, err := client.GenerateOutline(
		ctx,
		"创作一个修仙小说大纲",
		"test_user",
		"test_project",
		nil,
	)
	require.NoError(t, err, "大纲生成失败")

	charResp, err := client.GenerateCharacters(
		ctx,
		"根据大纲创建主要角色",
		"test_user",
		"test_project",
		outlineResp.Outline,
		nil,
	)
	require.NoError(t, err, "角色生成失败")

	fmt.Printf("\n📊 测试情节生成\n")

	// 生成情节
	startTime := time.Now()
	resp, err := client.GeneratePlot(
		ctx,
		"根据大纲和角色设计情节",
		"test_user",
		"test_project",
		outlineResp.Outline,
		charResp.Characters,
		nil,
	)
	duration := time.Since(startTime)

	require.NoError(t, err, "情节生成失败")
	require.NotNil(t, resp, "情节响应为空")
	require.NotNil(t, resp.Plot, "情节数据为空")
	assert.Greater(t, len(resp.Plot.TimelineEvents), 0, "事件数量为0")

	fmt.Printf("\n✅ 情节生成成功! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   📅 事件数量: %d\n", len(resp.Plot.TimelineEvents))
	fmt.Printf("   🧵 情节线数: %d\n", len(resp.Plot.PlotThreads))

	if len(resp.Plot.TimelineEvents) > 0 {
		fmt.Println("   主要事件:")
		for i, event := range resp.Plot.TimelineEvents {
			if i < 3 {
				fmt.Printf("     %d. %s (%s)\n", i+1, event.Title, event.Timestamp)
				fmt.Printf("        类型: %s\n", event.EventType)
			}
		}
	}
}

// TestCompleteWorkflow 测试完整工作流
func TestCompleteWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过长时间测试")
	}

	client := requireAIService(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	fmt.Printf("\n🎨 测试完整创作工作流\n")
	fmt.Printf("   ⏳ 这可能需要30-60秒...\n")

	startTime := time.Now()
	resp, err := client.ExecuteCreativeWorkflow(
		ctx,
		"创作一个都市爱情小说的完整设定，包含3章内容",
		"test_user",
		"test_project",
		3,     // 最大反思次数
		false, // 不启用人工审核
		nil,
	)
	duration := time.Since(startTime)

	require.NoError(t, err, "工作流执行失败")
	require.NotNil(t, resp, "工作流响应为空")

	// 验证响应
	assert.NotEmpty(t, resp.ExecutionId, "执行ID为空")
	assert.NotNil(t, resp.Outline, "大纲数据为空")
	assert.NotNil(t, resp.Characters, "角色数据为空")
	assert.NotNil(t, resp.Plot, "情节数据为空")

	fmt.Printf("\n✅ 工作流执行成功! 总耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   🆔 执行ID: %s\n", resp.ExecutionId)
	fmt.Printf("   ✓  审核状态: %v\n", resp.ReviewPassed)
	fmt.Printf("   🔄 反思次数: %d\n", resp.ReflectionCount)

	if resp.Outline != nil {
		fmt.Printf("\n   📖 大纲: %s\n", resp.Outline.Title)
		fmt.Printf("      章节数: %d\n", len(resp.Outline.Chapters))
	}

	if resp.Characters != nil {
		fmt.Printf("\n   👥 角色数: %d\n", len(resp.Characters.Characters))
	}

	if resp.Plot != nil {
		fmt.Printf("\n   📊 事件数: %d\n", len(resp.Plot.TimelineEvents))
		fmt.Printf("      情节线: %d\n", len(resp.Plot.PlotThreads))
	}

	if len(resp.ExecutionTimes) > 0 {
		fmt.Println("\n   ⏱️  执行时间分析:")
		totalTime := float32(0)
		for stage, execTime := range resp.ExecutionTimes {
			fmt.Printf("      %s: %.2f秒\n", stage, execTime)
			totalTime += execTime
		}
		fmt.Printf("      总计: %.2f秒\n", totalTime)
	}
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发测试")
	}

	client := requireAIService(t)
	defer client.Close()

	fmt.Printf("\n🔀 测试并发请求\n")

	concurrency := 3
	done := make(chan bool, concurrency)
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
			defer cancel()

			task := fmt.Sprintf("创作第%d个故事大纲", id+1)
			resp, err := client.GenerateOutline(ctx, task, "test_user", "test_project", nil)

			if err != nil {
				errors <- err
			} else if resp == nil || resp.Outline == nil {
				errors <- fmt.Errorf("响应为空")
			} else {
				fmt.Printf("   [%d] ✓ 完成: %s\n", id+1, resp.Outline.Title)
				done <- true
			}
		}(i)
	}

	// 等待所有请求完成
	successCount := 0
	for i := 0; i < concurrency; i++ {
		select {
		case <-done:
			successCount++
		case err := <-errors:
			t.Logf("请求失败: %v", err)
		case <-time.After(testTimeout):
			t.Fatal("并发测试超时")
		}
	}

	assert.Equal(t, concurrency, successCount, "部分并发请求失败")
	fmt.Printf("\n✅ 并发测试通过! 成功: %d/%d\n", successCount, concurrency)
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	client := requireAIService(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("\n⚠️  测试错误处理\n")

	// 测试空任务
	_, err := client.GenerateOutline(ctx, "", "test_user", "test_project", nil)
	assert.Error(t, err, "空任务应该返回错误")
	fmt.Printf("   [1/2] ✓ 空任务错误处理正常\n")

	// 测试超时
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer shortCancel()

	_, err = client.GenerateOutline(shortCtx, "测试任务", "test_user", "test_project", nil)
	assert.Error(t, err, "超时应该返回错误")
	fmt.Printf("   [2/2] ✓ 超时错误处理正常\n")

	fmt.Printf("\n✅ 错误处理测试通过\n")
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
