package migrateapp

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	migrationpkg "Qingyu_backend/migration"
	mongodbpkg "Qingyu_backend/migration/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run executes one named MongoDB migration command.
func Run(args []string) error {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	command := fs.String("command", "", "迁移命令: up 或 down")
	name := fs.String("name", "", "迁移名称")
	env := fs.String("env", "dev", "环境: dev, staging, production")
	force := fs.Bool("force", false, "强制执行生产环境迁移")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := validateOptions(*command, *name, *env); err != nil {
		return err
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := databaseName(*env)

	if *env == "production" && !*force {
		fmt.Println("⚠️  警告：即将在生产环境执行迁移！")
		fmt.Print("请输入 'yes' 确认: ")
		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if confirm != "yes" {
			return errors.New("❌ 操作已取消")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("❌ 连接MongoDB失败: %w", err)
	}
	defer func() {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("⚠️  断开MongoDB连接失败: %v", disconnectErr)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("❌ MongoDB连接测试失败: %w", err)
	}

	db := client.Database(dbName)
	log.Printf("✅ 已连接到数据库: %s", dbName)
	migrator := migrationpkg.NewMigrator(db)
	if err := registerMigrations(migrator); err != nil {
		return err
	}

	switch *command {
	case "up":
		log.Printf("🚀 开始执行迁移: %s", *name)
		if err := migrator.Up(ctx, *name); err != nil {
			return fmt.Errorf("❌ 迁移执行失败: %w", err)
		}
		log.Printf("✅ 迁移执行成功: %s", *name)
	case "down":
		log.Printf("🔄 开始回滚迁移: %s", *name)
		if err := migrator.Down(ctx, *name); err != nil {
			return fmt.Errorf("❌ 迁移回滚失败: %w", err)
		}
		log.Printf("✅ 迁移回滚成功: %s", *name)
	}

	fmt.Println("\n✨ 操作完成！")
	return nil
}

func validateOptions(command, name, env string) error {
	if command == "" {
		return errors.New("❌ 请指定 -command 参数 (up 或 down)")
	}
	if command != "up" && command != "down" {
		return fmt.Errorf("❌ 无效的命令: %s (只支持 up 或 down)", command)
	}
	if name == "" {
		return errors.New("❌ 请指定 -name 参数 (迁移名称)")
	}
	if env != "dev" && env != "staging" && env != "production" {
		return fmt.Errorf("❌ 无效的环境: %s (只支持 dev, staging, production)", env)
	}
	return nil
}

func databaseName(env string) string {
	switch env {
	case "production":
		return "qingyu"
	case "staging":
		return "qingyu_staging"
	default:
		return "qingyu_dev"
	}
}

func registerMigrations(migrator *migrationpkg.Migrator) error {
	migrations := []struct {
		name      string
		migration migrationpkg.SimpleMigration
	}{
		{"002_create_users_indexes", &mongodbpkg.CreateUsersIndexes{}},
		{"003_create_books_indexes_p0", &mongodbpkg.CreateBooksIndexesP0{}},
		{"004_create_chapters_indexes", &mongodbpkg.CreateChaptersIndexes{}},
		{"005_create_reading_progress_indexes", &mongodbpkg.CreateReadingProgressIndexes{}},
		{"006_create_core_query_indexes", &mongodbpkg.CreateCoreQueryIndexes{}},
	}
	for _, migration := range migrations {
		if err := migrator.Register(migration.name, migration.migration); err != nil {
			return fmt.Errorf("❌ 注册迁移失败 %s: %w", migration.name, err)
		}
	}
	return nil
}
