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

	"go.mongodb.org/mongo-driver/mongo"
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
	fmt.Printf("  数据库: %s\n", cfg.Database.Primary.MongoDB.URI)
	fmt.Printf("  数据库名: %s\n", cfg.Database.Primary.MongoDB.Database)
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

	// 显示菜单
	for {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("请选择操作：")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("1. 导入100本小说（完整数据，含章节内容）")
		fmt.Println("2. 试运行模式（仅验证数据，不写入）")
		fmt.Println("3. 创建索引")
		fmt.Println("4. 查看数据库统计")
		fmt.Println("5. 退出")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		var choice int
		fmt.Print("请输入选项 (1-5): ")
		fmt.Scanf("%d", &choice)

		fmt.Println()

		switch choice {
		case 1:
			importNovels(ctx, db, dataFile, false)
		case 2:
			importNovels(ctx, db, dataFile, true)
		case 3:
			createIndexes(ctx, db)
		case 4:
			showStatistics(ctx, db)
		case 5:
			fmt.Println("👋 再见！")
			os.Exit(0)
		default:
			fmt.Println("❌ 无效选项，请重新选择")
		}

		fmt.Println()
		fmt.Println("按 Enter 继续...")
		fmt.Scanln()
	}
}

func importNovels(ctx context.Context, db *mongo.Database, dataFile string, dryRun bool) {
	mode := "正式导入"
	if dryRun {
		mode = "试运行模式"
	}

	fmt.Printf("📚 开始导入小说 (%s)...\n", mode)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	importer := seeds.NewNovelImporter(db, dryRun)

	if err := importer.ImportFromJSON(ctx, dataFile); err != nil {
		fmt.Printf("❌ 导入失败: %v\n", err)
		return
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ 导入完成！")
	fmt.Println()
}

func createIndexes(ctx context.Context, db *mongo.Database) {
	fmt.Println("🔍 创建索引...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	importer := seeds.NewNovelImporter(db, false)

	if err := importer.CreateIndexes(ctx); err != nil {
		fmt.Printf("❌ 创建索引失败: %v\n", err)
		return
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ 索引创建成功！")
	fmt.Println()
}

func showStatistics(ctx context.Context, db *mongo.Database) {
	fmt.Println("📊 数据统计")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	importer := seeds.NewNovelImporter(db, false)

	if err := importer.GetStats(ctx); err != nil {
		fmt.Printf("❌ 获取统计失败: %v\n", err)
		return
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}
