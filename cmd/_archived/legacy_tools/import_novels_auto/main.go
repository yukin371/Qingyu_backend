package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/migration/seeds"
	"Qingyu_backend/service"
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║    青羽写作平台 - 小说批量导入工具     ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 加载配置
	fmt.Println("📁 加载配置文件...")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v\n", err)
	}
	fmt.Printf("✓ 配置加载成功\n")

	// 使用GetMongoConfig获取MongoDB配置
	mongoConfig, err := cfg.Database.GetMongoConfig()
	if err != nil {
		log.Fatalf("❌ 获取MongoDB配置失败: %v\n", err)
	}
	fmt.Printf("  数据库: %s\n", mongoConfig.URI)
	fmt.Printf("  数据库名: %s\n", mongoConfig.Database)
	fmt.Println()

	// 初始化服务（包括数据库连接）
	fmt.Println("🔗 连接数据库...")
	if err := core.InitServices(); err != nil {
		log.Fatalf("❌ 初始化服务失败: %v\n", err)
	}
	fmt.Println("✓ 数据库连接成功")
	fmt.Println()

	ctx := context.Background()

	// 获取数据库连接
	db := service.GetServiceContainer().GetMongoDB()
	if db == nil {
		log.Fatal("❌ 数据库未初始化")
	}

	// 检查数据文件
	dataFile := "data/novels_100.json"
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		log.Fatalf("❌ 数据文件不存在: %s\n", dataFile)
	}
	fmt.Printf("📄 找到数据文件: %s\n\n", dataFile)

	// 直接执行导入
	fmt.Println("📚 开始导入100本小说...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	importer := seeds.NewNovelImporter(db, false)

	if err := importer.ImportFromJSON(ctx, dataFile); err != nil {
		fmt.Printf("❌ 导入失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ 导入完成！")
	fmt.Println()

	// 创建索引
	fmt.Println("🔍 创建索引...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := importer.CreateIndexes(ctx); err != nil {
		fmt.Printf("⚠️  创建索引失败: %v\n", err)
	} else {
		fmt.Println("✓ 索引创建成功！")
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 显示统计
	fmt.Println("📊 数据库统计")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := importer.GetStats(ctx); err != nil {
		fmt.Printf("❌ 获取统计失败: %v\n", err)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ 全部完成！")
}
