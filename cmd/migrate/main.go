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
	command := flag.String("command", "status", "Command to run: up, down, status, reset, seed")
	steps := flag.String("steps", "0", "Number of steps for down command (0 means all)")
	configPath := flag.String("config", ".", "Path to config file")
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
	fmt.Println("\n=== Running Seeds ===\n")

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
