//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config 清理配置
type Config struct {
	MongoURI string
	Database string
	DryRun   bool // 是否为试运行
	BatchSize int  // 批量处理大小
}

func main() {
	// 配置
	config := Config{
		MongoURI:   "mongodb://localhost:27017",
		Database:   "qingyu",
		DryRun:     false, // 实际清理模式
		BatchSize:  1000,
	}

	fmt.Println("========================================")
	fmt.Println("章节内容清理工具")
	fmt.Println("========================================")
	fmt.Printf("MongoDB URI: %s\n", config.MongoURI)
	fmt.Printf("Database: %s\n", config.Database)
	fmt.Printf("Dry Run: %v\n", config.DryRun)
	fmt.Printf("Batch Size: %d\n", config.BatchSize)
	fmt.Println("========================================\n")
	fmt.Println("此工具将删除 chapters 集合中的 content 字段")
	fmt.Println("（前提：内容已迁移到 chapter_contents 集合）\n")

	// 确认
	if !config.DryRun {
		fmt.Print("警告：此操作将删除 chapters.content 字段！确认继续？(yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("已取消清理")
			return
		}
	}

	// 执行清理
	if err := cleanup(context.Background(), config); err != nil {
		log.Fatalf("清理失败: %v", err)
	}

	fmt.Println("\n清理完成！")
}

func cleanup(ctx context.Context, config Config) error {
	// 连接 MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(config.Database)
	chaptersCol := db.Collection("chapters")
	contentsCol := db.Collection("chapter_contents")

	// 1. 统计
	fmt.Println("1. 统计需要清理的章节...")
	totalWithContent, err := chaptersCol.CountDocuments(ctx, bson.M{
		"content": bson.M{"$exists": true, "$ne": ""},
	})
	if err != nil {
		return fmt.Errorf("统计章节失败: %w", err)
	}
	fmt.Printf("   找到 %d 个包含 content 字段的章节\n\n", totalWithContent)

	if totalWithContent == 0 {
		fmt.Println("没有需要清理的章节")
		return nil
	}

	// 2. 验证内容是否已迁移
	fmt.Println("2. 验证内容是否已迁移...")
	totalContentDocs, err := contentsCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计 chapter_contents 失败: %w", err)
	}
	fmt.Printf("   chapter_contents 集合文档数: %d\n", totalContentDocs)

	if totalContentDocs < totalWithContent {
		fmt.Printf("   警告：chapter_contents 文档数 (%d) 小于需要清理的章节数 (%d)\n", totalContentDocs, totalWithContent)
		fmt.Println("   建议先运行迁移脚本，确保所有内容都已迁移")
		return fmt.Errorf("内容未完全迁移")
	}
	fmt.Println("   ✓ 所有章节内容均已迁移\n")

	// 3. 清理 content 字段
	fmt.Println("3. 开始清理 content 字段...")

	var cleanedCount int64

	if !config.DryRun {
		// 使用 $unset 操作符删除 content 字段
		update := bson.M{"$unset": bson.M{"content": ""}}

		// 一次性删除所有 content 字段
		result, err := chaptersCol.UpdateMany(ctx,
			bson.M{"content": bson.M{"$exists": true}},
			update,
		)
		if err != nil {
			return fmt.Errorf("清理 content 字段失败: %w", err)
		}

		cleanedCount = result.MatchedCount
		fmt.Printf("   已清理: %d 个章节\n", cleanedCount)
	} else {
		// 试运行模式
		cleanedCount = totalWithContent
		fmt.Printf("   [DRY RUN] 将清理: %d 个章节\n", cleanedCount)
	}

	// 4. 验证清理结果
	fmt.Println("\n4. 验证清理结果...")
	remainingWithContent, _ := chaptersCol.CountDocuments(ctx, bson.M{
		"content": bson.M{"$exists": true, "$ne": ""},
	})
	fmt.Printf("   chapters 集合仍有 content 字段的文档数: %d\n", remainingWithContent)

	// 5. 打印统计
	fmt.Println("\n========================================")
	fmt.Println("清理统计")
	fmt.Println("========================================")
	fmt.Printf("总数: %d\n", totalWithContent)
	fmt.Printf("已清理: %d\n", cleanedCount)
	fmt.Printf("剩余: %d\n", remainingWithContent)
	if totalWithContent > 0 {
		fmt.Printf("成功率: %.2f%%\n", float64(cleanedCount)/float64(totalWithContent)*100)
	}
	fmt.Println("========================================")

	if remainingWithContent > 0 {
		fmt.Println("\n警告：仍有章节包含 content 字段，可能需要重新运行清理")
	} else {
		fmt.Println("\n✓ 所有章节的 content 字段已成功清理！")
	}

	return nil
}
