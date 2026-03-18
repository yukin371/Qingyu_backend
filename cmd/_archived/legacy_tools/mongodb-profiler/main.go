package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProfilerConfig Profiler配置
type ProfilerConfig struct {
	Level        int   // 0=off, 1=slow only, 2=all
	SlowMS       int64 // 慢查询阈值（毫秒）
	ProfilerSize int64 // Profiler存储大小（字节）
}

// DefaultProfilerConfig 默认Profiler配置
var DefaultProfilerConfig = ProfilerConfig{
	Level:        1,
	SlowMS:       100,
	ProfilerSize: 100 * 1024 * 1024, // 100MB
}

func main() {
	// 从环境变量或命令行参数获取MongoDB URI
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "qingyu"
	}

	fmt.Printf("正在连接到MongoDB: %s\n", uri)
	fmt.Printf("目标数据库: %s\n\n", dbName)

	// 创建MongoDB客户端
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Ping MongoDB失败: %v", err)
	}
	fmt.Println("✓ MongoDB连接成功")

	// 应用profiling配置
	config := DefaultProfilerConfig
	if err := applyProfilerConfig(ctx, client, dbName, config); err != nil {
		log.Fatalf("应用profiling配置失败: %v", err)
	}

	// 验证配置
	if err := verifyProfilerConfig(ctx, client, dbName); err != nil {
		log.Fatalf("验证profiling配置失败: %v", err)
	}

	fmt.Println("\n✅ Profiler配置完成！")
}

// applyProfilerConfig 应用profiling配置
func applyProfilerConfig(ctx context.Context, client *mongo.Client, dbName string, config ProfilerConfig) error {
	db := client.Database(dbName)

	fmt.Println("=== 配置MongoDB Profiler ===")

	// 1. 设置profiling级别
	fmt.Printf("设置profiling级别: %d\n", config.Level)
	fmt.Printf("设置慢查询阈值: %dms\n", config.SlowMS)

	profileCmd := bson.D{
		{Key: "profile", Value: config.Level},
		{Key: "slowms", Value: config.SlowMS},
	}

	var profileResult bson.M
	if err := db.RunCommand(ctx, profileCmd).Decode(&profileResult); err != nil {
		return fmt.Errorf("设置profiling级别失败: %w", err)
	}
	fmt.Println("✓ Profiling级别配置成功")

	// 2. 设置system.profile collection为capped collection
	fmt.Printf("\n设置profiler存储上限: %dMB\n", config.ProfilerSize/(1024*1024))

	// 尝试创建capped collection（如果已存在且是capped的会报错，忽略即可）
	collections, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: "system.profile"}})
	if err == nil && len(collections) == 0 {
		// collection不存在，创建新的
		createCmd := bson.D{
			{Key: "create", Value: "system.profile"},
			{Key: "capped", Value: true},
			{Key: "size", Value: config.ProfilerSize},
		}
		var createResult bson.M
		if err := db.RunCommand(ctx, createCmd).Decode(&createResult); err != nil {
			// 可能已经存在，继续
			fmt.Printf("⚠ 创建capped collection时出现警告: %v\n", err)
		} else {
			fmt.Println("✓ 创建system.profile capped collection成功")
		}
	} else if err == nil && len(collections) > 0 {
		// collection已存在，尝试转换
		convertCmd := bson.D{
			{Key: "convertToCapped", Value: "system.profile"},
			{Key: "size", Value: config.ProfilerSize},
		}
		var convertResult bson.M
		if err := db.RunCommand(ctx, convertCmd).Decode(&convertResult); err != nil {
			// 可能已经是capped的，继续
			fmt.Printf("⚠ 转换capped collection时出现警告: %v\n", err)
		} else {
			fmt.Println("✓ 转换system.profile为capped collection成功")
		}
	}

	fmt.Println("\n=== 配置完成 ===")
	fmt.Printf("   级别: %d\n", config.Level)
	fmt.Printf("   阈值: %dms\n", config.SlowMS)
	fmt.Printf("   存储: %dMB\n", config.ProfilerSize/(1024*1024))

	return nil
}

// verifyProfilerConfig 验证profiling配置
func verifyProfilerConfig(ctx context.Context, client *mongo.Client, dbName string) error {
	db := client.Database(dbName)

	fmt.Println("\n=== 验证当前配置 ===")

	// 获取profiling状态
	statusCmd := bson.D{{Key: "profile", Value: -1}}
	var status bson.M
	if err := db.RunCommand(ctx, statusCmd).Decode(&status); err != nil {
		return fmt.Errorf("获取profiling状态失败: %w", err)
	}

	fmt.Println("当前配置:")
	if level, ok := status["was"].(int32); ok {
		fmt.Printf("   Profiling级别: %d\n", level)
	}
	if slowMs, ok := status["slowms"].(int32); ok {
		fmt.Printf("   慢查询阈值: %dms\n", slowMs)
	}

	// 获取system.profile统计信息
	coll := db.Collection("system.profile")
	count, err := coll.CountDocuments(ctx, bson.D{})
	if err == nil {
		fmt.Printf("\nsystem.profile集合统计:\n")
		fmt.Printf("   文档数量: %d\n", count)
	}

	// 显示示例查询
	fmt.Println("\n=== 示例查询 ===")
	fmt.Println("查看最慢的查询:")
	fmt.Println(`  db.system.profile.find().sort({millis: -1}).limit(5)`)
	fmt.Println("\n查看特定collection的慢查询:")
	fmt.Println(`  db.system.profile.find({ns: "qingyu.users"}).sort({millis: -1})`)

	return nil
}
