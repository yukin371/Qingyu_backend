//go:build auto

package main

import (
	"context"
	"fmt"
	"log"

	"Qingyu_backend/config"
	"Qingyu_backend/migration/seeds"
	"Qingyu_backend/service"

	"go.mongodb.org/mongo-driver/mongo"
)

// 自动运行版本 - 非交互式
func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║    青羽写作平台 - 自动创建测试数据     ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 加载配置
	fmt.Println("📁 加载配置文件...")
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v\n", err)
	}
	fmt.Println("✓ 配置加载成功")
	fmt.Println()

	// 初始化服务容器
	fmt.Println("🔗 连接数据库...")
	if err := service.InitializeServices(); err != nil {
		log.Fatalf("❌ 初始化服务失败: %v\n", err)
	}
	serviceContainer := service.GetServiceContainer()
	if serviceContainer == nil {
		log.Fatal("❌ 服务容器未初始化")
	}
	fmt.Println("✓ 数据库连接成功")
	fmt.Println()

	ctx := context.Background()

	// 获取数据库连接
	db := serviceContainer.GetMongoDB()
	if db == nil {
		log.Fatal("❌ 无法获取数据库连接")
	}

	// 直接创建所有数据（不清理）
	createAllData(ctx, db)
}

func createAllData(ctx context.Context, db *mongo.Database) {
	fmt.Println("✨ 开始创建测试数据...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 1. 创建用户
	fmt.Println("【1/6】创建用户数据")
	if err := seeds.SeedEnhancedUsers(ctx, db); err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err)
	} else {
		fmt.Println("✓ 用户数据创建完成")
	}
	fmt.Println()

	// 2. 创建联调发布数据
	fmt.Println("【2/6】创建发布联调数据")
	publicationInfo, err := seeds.SeedPublicationFlowData(ctx, db)
	if err != nil {
		fmt.Printf("❌ 创建发布联调数据失败: %v\n", err)
	} else {
		fmt.Println("✓ 发布联调数据创建完成")
		fmt.Printf("  项目ID: %s\n", publicationInfo.ProjectID)
		fmt.Printf("  文档ID: %v\n", publicationInfo.DocumentIDs)
	}
	fmt.Println()

	// 3. 创建书籍
	fmt.Println("【3/6】创建书籍数据")
	if err := seeds.SeedBooks(ctx, db); err != nil {
		fmt.Printf("❌ 创建书籍失败: %v\n", err)
	} else {
		fmt.Println("✓ 书籍数据创建完成")
	}
	fmt.Println()

	// 4. 创建章节
	fmt.Println("【4/6】创建章节数据")
	if err := seeds.SeedChapters(ctx, db); err != nil {
		fmt.Printf("❌ 创建章节失败: %v\n", err)
	} else {
		fmt.Println("✓ 章节数据创建完成")
	}
	fmt.Println()

	// 5. 创建钱包
	fmt.Println("【5/6】创建钱包数据")
	if err := seeds.SeedWallets(ctx, db); err != nil {
		fmt.Printf("❌ 创建钱包失败: %v\n", err)
	} else {
		fmt.Println("✓ 钱包数据创建完成")
	}
	fmt.Println()

	// 6. 创建社交数据
	fmt.Println("【6/6】创建社交数据")
	if err := seeds.SeedSocialData(ctx, db); err != nil {
		fmt.Printf("❌ 创建社交数据失败: %v\n", err)
	} else {
		fmt.Println("✓ 社交数据创建完成")
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ 所有测试数据创建完成！")
	fmt.Println()
}
