package ai

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// 测试用的gRPC服务器地址
	testGRPCAddress = "localhost:50051"
)

func TestNewPhase3Client(t *testing.T) {
	client, err := NewPhase3Client(testGRPCAddress)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()
}

func TestPhase3Client_HealthCheck(t *testing.T) {
	client, err := NewPhase3Client(testGRPCAddress)
	if err != nil {
		t.Skipf("无法连接到gRPC服务器: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()
	response, err := client.HealthCheck(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	t.Logf("健康状态: %s", response.Status)
	t.Logf("检查结果: %v", response.Checks)
}

func TestPhase3Client_GenerateOutline(t *testing.T) {
	client, err := NewPhase3Client(testGRPCAddress)
	if err != nil {
		t.Skipf("无法连接到gRPC服务器: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	response, err := client.GenerateOutline(
		ctx,
		"创作一个修仙小说大纲，主角是天才少年，包含5章内容",
		"test_user",
		"test_project",
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.Outline)

	t.Logf("✅ 大纲生成成功")
	t.Logf("标题: %s", response.Outline.Title)
	t.Logf("类型: %s", response.Outline.Genre)
	t.Logf("章节数: %d", len(response.Outline.Chapters))
	t.Logf("执行时间: %.2f秒", response.ExecutionTime)

	// 验证章节
	assert.Greater(t, len(response.Outline.Chapters), 0)
	if len(response.Outline.Chapters) > 0 {
		chapter := response.Outline.Chapters[0]
		t.Logf("第一章: %s", chapter.Title)
		t.Logf("概要: %s", chapter.Summary)
	}
}

func TestPhase3Client_ExecuteCreativeWorkflow(t *testing.T) {
	// 这个测试会比较慢（30-60秒），可以标记为长测试
	if testing.Short() {
		t.Skip("跳过长测试")
	}

	client, err := NewPhase3Client(testGRPCAddress)
	if err != nil {
		t.Skipf("无法连接到gRPC服务器: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	startTime := time.Now()
	response, err := client.ExecuteCreativeWorkflow(
		ctx,
		"创作一个都市爱情小说的完整设定，包含3章内容",
		"test_user",
		"test_project",
		3,     // 最大反思次数
		false, // 不启用人工审核
		nil,
	)
	duration := time.Since(startTime)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	t.Logf("\n%s", strings.Repeat("=", 60))
	t.Logf("✅ 完整工作流执行成功")
	t.Logf("%s", strings.Repeat("=", 60))
	t.Logf("执行ID: %s", response.ExecutionId)
	t.Logf("审核通过: %v", response.ReviewPassed)
	t.Logf("反思次数: %d", response.ReflectionCount)
	t.Logf("总耗时: %.2f秒", duration.Seconds())

	// 大纲信息
	if response.Outline != nil {
		t.Logf("\n📖 大纲:")
		t.Logf("  标题: %s", response.Outline.Title)
		t.Logf("  章节数: %d", len(response.Outline.Chapters))
	}

	// 角色信息
	if response.Characters != nil {
		t.Logf("\n👥 角色:")
		t.Logf("  角色数: %d", len(response.Characters.Characters))
		if len(response.Characters.Characters) > 0 {
			char := response.Characters.Characters[0]
			t.Logf("  主角: %s (%s)", char.Name, char.RoleType)
		}
	}

	// 情节信息
	if response.Plot != nil {
		t.Logf("\n📊 情节:")
		t.Logf("  事件数: %d", len(response.Plot.TimelineEvents))
		t.Logf("  情节线: %d", len(response.Plot.PlotThreads))
	}

	// 执行时间分析
	if len(response.ExecutionTimes) > 0 {
		t.Logf("\n⏱️  执行时间:")
		totalTime := float32(0)
		for stage, execTime := range response.ExecutionTimes {
			t.Logf("  %s: %.2f秒", stage, execTime)
			totalTime += execTime
		}
		t.Logf("  总计: %.2f秒", totalTime)
	}
}

// Benchmark测试
func BenchmarkPhase3Client_GenerateOutline(b *testing.B) {
	client, err := NewPhase3Client(testGRPCAddress)
	if err != nil {
		b.Skipf("无法连接到gRPC服务器: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GenerateOutline(
			ctx,
			"创作一个科幻小说大纲",
			"bench_user",
			"bench_project",
			nil,
		)
		if err != nil {
			b.Fatalf("生成失败: %v", err)
		}
	}
}
