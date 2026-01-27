package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	migrationpkg "Qingyu_backend/migration"
	mongodbpkg "Qingyu_backend/migration/mongodb"
)

func main() {
	// 命令行参数
	command := flag.String("command", "", "迁移命令: up 或 down")
	name := flag.String("name", "", "迁移名称")
	env := flag.String("env", "dev", "环境: dev, staging, production")
	flag.Parse()

	// 验证命令
	if *command == "" {
		log.Fatal("❌ 请指定 -command 参数 (up 或 down)")
	}
	if *command != "up" && *command != "down" {
		log.Fatalf("❌ 无效的命令: %s (只支持 up 或 down)", *command)
	}
	if *name == "" {
		log.Fatal("❌ 请指定 -name 参数 (迁移名称)")
	}

	// 验证环境
	if *env != "dev" && *env != "staging" && *env != "production" {
		log.Fatalf("❌ 无效的环境: %s (只支持 dev, staging, production)", *env)
	}

	// 获取MongoDB连接字符串
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// 确定数据库名称
	var dbName string
	switch *env {
	case "production":
		dbName = "qingyu"
	default:
		dbName = "qingyu_staging"
	}

	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ 连接MongoDB失败: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("⚠️  断开MongoDB连接失败: %v", err)
		}
	}()

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB连接测试失败: %v", err)
	}

	db := client.Database(dbName)
	log.Printf("✅ 已连接到数据库: %s", dbName)

	// 创建迁移器
	migrator := migrationpkg.NewMigrator(db)

	// 注册所有索引迁移
	migrator.Register("002_create_users_indexes", &mongodbpkg.CreateUsersIndexes{})
	migrator.Register("003_create_books_indexes_p0", &mongodbpkg.CreateBooksIndexesP0{})
	migrator.Register("004_create_chapters_indexes", &mongodbpkg.CreateChaptersIndexes{})
	migrator.Register("005_create_reading_progress_indexes", &mongodbpkg.CreateReadingProgressIndexes{})

	// 执行命令
	switch *command {
	case "up":
		log.Printf("🚀 开始执行迁移: %s", *name)
		if err := migrator.Up(ctx, *name); err != nil {
			log.Fatalf("❌ 迁移执行失败: %v", err)
		}
		log.Printf("✅ 迁移执行成功: %s", *name)

	case "down":
		log.Printf("🔄 开始回滚迁移: %s", *name)
		if err := migrator.Down(ctx, *name); err != nil {
			log.Fatalf("❌ 迁移回滚失败: %v", err)
		}
		log.Printf("✅ 迁移回滚成功: %s", *name)
	}

	fmt.Println("\n✨ 操作完成！")
}
