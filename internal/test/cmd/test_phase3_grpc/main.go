package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"Qingyu_backend/service/ai"
)

var (
	grpcAddr = flag.String("addr", "localhost:50051", "gRPC服务器地址")
	task     = flag.String("task", "创作一个修仙小说大纲，主角是天才少年", "创作任务")
	workflow = flag.Bool("workflow", false, "是否执行完整工作流")
)

func main() {
	flag.Parse()

	fmt.Println("========================================")
	fmt.Println("Phase3 gRPC客户端测试")
	fmt.Println("========================================")
	fmt.Println()

	// 创建客户端
	fmt.Printf("连接到gRPC服务器: %s\n", *grpcAddr)
	client, err := ai.NewPhase3Client(*grpcAddr)
	if err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ 连接成功")
	fmt.Println()

	// 1. 健康检查
	fmt.Println("1️⃣  健康检查...")
	ctx := context.Background()
	healthResp, err := client.HealthCheck(ctx)
	if err != nil {
		log.Fatalf("❌ 健康检查失败: %v", err)
	}
	fmt.Printf("   状态: %s\n", healthResp.Status)
	fmt.Println("   组件状态:")
	for name, status := range healthResp.Checks {
		fmt.Printf("     - %s: %s\n", name, status)
	}
	fmt.Println()

	if *workflow {
		// 执行完整工作流
		testCompleteWorkflow(client, *task)
	} else {
		// 测试单个Agent
		testIndividualAgents(client, *task)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✅ 测试完成")
	fmt.Println("========================================")
}

func testIndividualAgents(client *ai.Phase3Client, task string) {
	ctx := context.Background()

	// 2. 生成大纲
	fmt.Println("2️⃣  生成大纲...")
	fmt.Printf("   任务: %s\n", task)

	startTime := time.Now()
	outlineResp, err := client.GenerateOutline(ctx, task, "test_user", "test_project", nil)
	if err != nil {
		log.Fatalf("❌ 大纲生成失败: %v", err)
	}
	duration := time.Since(startTime)

	fmt.Printf("   ✅ 成功! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   📖 标题: %s\n", outlineResp.Outline.Title)
	fmt.Printf("   🎭 类型: %s\n", outlineResp.Outline.Genre)
	fmt.Printf("   📚 章节数: %d\n", len(outlineResp.Outline.Chapters))

	if len(outlineResp.Outline.Chapters) > 0 {
		fmt.Println("   章节列表:")
		for i, chapter := range outlineResp.Outline.Chapters {
			if i < 3 { // 只显示前3章
				fmt.Printf("     %d. %s\n", i+1, chapter.Title)
				if chapter.Summary != "" {
					summary := chapter.Summary
					if len(summary) > 50 {
						summary = summary[:50] + "..."
					}
					fmt.Printf("        %s\n", summary)
				}
			}
		}
		if len(outlineResp.Outline.Chapters) > 3 {
			fmt.Printf("     ... 还有 %d 章\n", len(outlineResp.Outline.Chapters)-3)
		}
	}
	fmt.Println()

	// 3. 生成角色
	fmt.Println("3️⃣  生成角色...")

	startTime = time.Now()
	charResp, err := client.GenerateCharacters(
		ctx,
		"根据大纲创建主要角色",
		"test_user",
		"test_project",
		outlineResp.Outline,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ 角色生成失败: %v", err)
	}
	duration = time.Since(startTime)

	fmt.Printf("   ✅ 成功! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   👥 角色数: %d\n", len(charResp.Characters.Characters))

	if len(charResp.Characters.Characters) > 0 {
		fmt.Println("   角色列表:")
		for i, char := range charResp.Characters.Characters {
			if i < 3 { // 只显示前3个角色
				fmt.Printf("     %d. %s (%s)\n", i+1, char.Name, char.RoleType)
				if char.Personality != nil && len(char.Personality.Traits) > 0 {
					fmt.Printf("        性格: %v\n", char.Personality.Traits[:min(3, len(char.Personality.Traits))])
				}
			}
		}
		if len(charResp.Characters.Characters) > 3 {
			fmt.Printf("     ... 还有 %d 个角色\n", len(charResp.Characters.Characters)-3)
		}
	}
	fmt.Println()

	// 4. 生成情节
	fmt.Println("4️⃣  生成情节...")

	startTime = time.Now()
	plotResp, err := client.GeneratePlot(
		ctx,
		"根据大纲和角色设计情节",
		"test_user",
		"test_project",
		outlineResp.Outline,
		charResp.Characters,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ 情节生成失败: %v", err)
	}
	duration = time.Since(startTime)

	fmt.Printf("   ✅ 成功! 耗时: %.2f秒\n", duration.Seconds())
	fmt.Printf("   📅 事件数: %d\n", len(plotResp.Plot.TimelineEvents))
	fmt.Printf("   🧵 情节线: %d\n", len(plotResp.Plot.PlotThreads))

	if len(plotResp.Plot.TimelineEvents) > 0 {
		fmt.Println("   主要事件:")
		for i, event := range plotResp.Plot.TimelineEvents {
			if i < 3 { // 只显示前3个事件
				fmt.Printf("     %d. %s (%s)\n", i+1, event.Title, event.Timestamp)
				fmt.Printf("        类型: %s\n", event.EventType)
			}
		}
		if len(plotResp.Plot.TimelineEvents) > 3 {
			fmt.Printf("     ... 还有 %d 个事件\n", len(plotResp.Plot.TimelineEvents)-3)
		}
	}
}

func testCompleteWorkflow(client *ai.Phase3Client, task string) {
	ctx := context.Background()

	fmt.Println("🎨 执行完整创作工作流...")
	fmt.Printf("   任务: %s\n", task)
	fmt.Println("   ⏳ 这可能需要30-60秒...")
	fmt.Println()

	startTime := time.Now()
	resp, err := client.ExecuteCreativeWorkflow(
		ctx,
		task,
		"test_user",
		"test_project",
		3,     // 最大反思次数
		false, // 不启用人工审核
		nil,
	)
	if err != nil {
		log.Fatalf("❌ 工作流执行失败: %v", err)
	}
	duration := time.Since(startTime)

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("✅ 工作流执行成功!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("🆔 执行ID: %s\n", resp.ExecutionId)
	fmt.Printf("✓  审核状态: %v\n", resp.ReviewPassed)
	fmt.Printf("🔄 反思次数: %d\n", resp.ReflectionCount)
	fmt.Printf("⏱️  总耗时: %.2f秒\n", duration.Seconds())
	fmt.Println()

	// 大纲信息
	if resp.Outline != nil {
		fmt.Println("📖 大纲:")
		fmt.Printf("   标题: %s\n", resp.Outline.Title)
		fmt.Printf("   章节数: %d\n", len(resp.Outline.Chapters))
	}
	fmt.Println()

	// 角色信息
	if resp.Characters != nil {
		fmt.Println("👥 角色:")
		fmt.Printf("   角色数: %d\n", len(resp.Characters.Characters))
		if len(resp.Characters.Characters) > 0 {
			for i, char := range resp.Characters.Characters {
				if i < 2 {
					fmt.Printf("   - %s (%s)\n", char.Name, char.RoleType)
				}
			}
		}
	}
	fmt.Println()

	// 情节信息
	if resp.Plot != nil {
		fmt.Println("📊 情节:")
		fmt.Printf("   事件数: %d\n", len(resp.Plot.TimelineEvents))
		fmt.Printf("   情节线: %d\n", len(resp.Plot.PlotThreads))
	}
	fmt.Println()

	// 执行时间分析
	if len(resp.ExecutionTimes) > 0 {
		fmt.Println("⏱️  执行时间分析:")
		totalTime := float32(0)
		for stage, execTime := range resp.ExecutionTimes {
			fmt.Printf("   %s: %.2f秒\n", stage, execTime)
			totalTime += execTime
		}
		fmt.Printf("   总计: %.2f秒\n", totalTime)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
