package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"Qingyu_backend/config"
	"Qingyu_backend/migration"
	"Qingyu_backend/migration/examples"
	"Qingyu_backend/migration/seeds"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 定义命令行参数
	command := flag.String("command", "status", "Command to run: up, down, status, reset, seed, import-novels, clean-novels")
	steps := flag.String("steps", "0", "Number of steps for down command (0 means all)")
	configPath := flag.String("config", ".", "Path to config file")
	novelFile := flag.String("file", "data/novels.json", "Path to novels JSON file for import-novels command")
	dryRun := flag.Bool("dry-run", false, "Dry run mode (validate only, don't insert)")
	category := flag.String("category", "", "Category filter for clean-novels command")
	flag.Parse()

	// 加载配置
	_, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 创建迁移管理器
	manager := migration.NewManager(db)

	// 注册所有迁移
	registerMigrations(manager)

	// 执行命令
	ctx := context.Background()

	switch *command {
	case "up":
		err = manager.Up(ctx)
	case "down":
		stepsInt, parseErr := strconv.Atoi(*steps)
		if parseErr != nil {
			log.Fatalf("Invalid steps value: %v", parseErr)
		}
		err = manager.Down(ctx, stepsInt)
	case "status":
		err = manager.Status(ctx)
	case "reset":
		fmt.Println("⚠️  WARNING: This will rollback all migrations!")
		fmt.Print("Are you sure? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm == "yes" {
			err = manager.Reset(ctx)
		} else {
			fmt.Println("Reset cancelled")
			return
		}
	case "seed":
		err = runSeeds(ctx, db)
	case "import-novels":
		err = importNovels(ctx, db, *novelFile, *dryRun)
	case "clean-novels":
		err = cleanNovels(ctx, db, *category)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}

	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	fmt.Println("\n✓ Command completed successfully")
}

// connectDB 连接数据库
func connectDB() (*mongo.Database, error) {
	cfg := config.GlobalConfig.Database
	if cfg == nil || cfg.Primary.MongoDB == nil {
		return nil, fmt.Errorf("database configuration is missing")
	}

	mongoCfg := cfg.Primary.MongoDB

	clientOptions := options.Client().
		ApplyURI(mongoCfg.URI).
		SetConnectTimeout(mongoCfg.ConnectTimeout).
		SetMaxPoolSize(mongoCfg.MaxPoolSize)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(mongoCfg.Database), nil
}

// registerMigrations 注册所有迁移
func registerMigrations(manager *migration.Manager) {
	// 注册示例迁移
	manager.RegisterMultiple(
		&examples.AddUserIndexes{},
		&examples.AddBookFields{},
		// 在这里添加更多迁移...
	)
}

// runSeeds 运行种子数据
func runSeeds(ctx context.Context, db *mongo.Database) error {
	fmt.Println("\n=== Running Seeds ===")

	// 运行用户种子
	if err := seeds.SeedUsers(ctx, db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// 运行分类种子
	if err := seeds.SeedCategories(ctx, db); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	// 运行书籍种子
	if err := seeds.SeedBooks(ctx, db); err != nil {
		return fmt.Errorf("failed to seed books: %w", err)
	}

	fmt.Println("\n✓ All seeds completed")
	return nil
}

// importNovels 导入小说数据
func importNovels(ctx context.Context, db *mongo.Database, filepath string, dryRun bool) error {
	fmt.Println("\n=== Importing Novels from CNNovel125K ===")

	// 创建导入器
	importer := seeds.NewNovelImporter(db, dryRun)

	// 从 JSON 文件导入
	if err := importer.ImportFromJSON(ctx, filepath); err != nil {
		return fmt.Errorf("failed to import novels: %w", err)
	}

	// 如果不是试运行模式，创建索引
	if !dryRun {
		if err := importer.CreateIndexes(ctx); err != nil {
			return fmt.Errorf("failed to create indexes: %w", err)
		}

		// 显示统计信息
		if err := importer.GetStats(ctx); err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}
	}

	return nil
}

// cleanNovels 清理小说数据
func cleanNovels(ctx context.Context, db *mongo.Database, category string) error {
	fmt.Println("\n=== Cleaning Novel Data ===")

	// 创建清理器
	cleaner := seeds.NewNovelCleaner(db)

	// 显示清理前统计
	if err := cleaner.GetStats(ctx); err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	// 确认操作
	if category == "" {
		fmt.Println("\n⚠️  WARNING: This will delete ALL books and chapters!")
		fmt.Print("Are you sure? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Clean cancelled")
			return nil
		}

		// 清理所有数据
		if err := cleaner.Clean(ctx); err != nil {
			return fmt.Errorf("failed to clean novels: %w", err)
		}
	} else {
		fmt.Printf("\n⚠️  WARNING: This will delete all books in category [%s]!\n", category)
		fmt.Print("Are you sure? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Clean cancelled")
			return nil
		}

		// 按分类清理
		if err := cleaner.CleanByCategory(ctx, category); err != nil {
			return fmt.Errorf("failed to clean novels by category: %w", err)
		}
	}

	// 显示清理后统计
	fmt.Println()
	if err := cleaner.GetStats(ctx); err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	return nil
}
