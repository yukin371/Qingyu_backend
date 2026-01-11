//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Chapter 旧章节结构（包含 Content 字段）
type OldChapter struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BookID    primitive.ObjectID `bson:"book_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"` // 需要迁移的字段
	WordCount int                `bson:"word_count"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ChapterContent 新章节内容结构
type ChapterContent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChapterID primitive.ObjectID `bson:"chapter_id"`
	Content   string             `bson:"content"`
	Format    string             `bson:"format"`
	Version   int                `bson:"version"`
	WordCount int                `bson:"word_count"`
	Hash      string             `bson:"hash,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// Config 迁移配置
type Config struct {
	MongoURI   string
	Database   string
	DryRun     bool // 是否为试运行（不实际写入数据）
	BatchSize  int  // 批量处理大小
	SkipErrors bool // 遇到错误是否继续
}

func main() {
	// 配置
	config := Config{
		MongoURI:   "mongodb://localhost:27017",
		Database:   "qingyu",
		DryRun:     false, // 实际迁移模式
		BatchSize:  100,
		SkipErrors: true,
	}

	fmt.Println("========================================")
	fmt.Println("章节内容数据迁移工具")
	fmt.Println("========================================")
	fmt.Printf("MongoDB URI: %s\n", config.MongoURI)
	fmt.Printf("Database: %s\n", config.Database)
	fmt.Printf("Dry Run: %v\n", config.DryRun)
	fmt.Printf("Batch Size: %d\n", config.BatchSize)
	fmt.Println("========================================\n")

	// 确认
	if !config.DryRun {
		fmt.Print("警告：此操作将修改数据库！确认继续？(yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("已取消迁移")
			return
		}
	}

	// 执行迁移
	if err := migrate(context.Background(), config); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	fmt.Println("\n迁移完成！")
}

func migrate(ctx context.Context, config Config) error {
	// 连接 MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(config.Database)
	chaptersCol := db.Collection("chapters")
	contentsCol := db.Collection("chapter_contents")

	// 1. 统计需要迁移的章节数
	fmt.Println("1. 统计需要迁移的章节...")
	totalCount, err := chaptersCol.CountDocuments(ctx, bson.M{
		"content": bson.M{"$exists": true, "$ne": ""},
	})
	if err != nil {
		return fmt.Errorf("统计章节失败: %w", err)
	}
	fmt.Printf("   找到 %d 个包含内容的章节\n\n", totalCount)

	if totalCount == 0 {
		fmt.Println("没有需要迁移的章节")
		return nil
	}

	// 2. 批量迁移
	fmt.Printf("2. 开始批量迁移（每批 %d 条）...\n\n", config.BatchSize)

	var migratedCount int
	var failedCount int
	var skippedCount int

	batchNum := 0
	skip := 0

	for {
		batchNum++

		// 查询一批数据
		cursor, err := chaptersCol.Find(ctx, bson.M{
			"content": bson.M{"$exists": true, "$ne": ""},
		}, options.Find().SetSkip(int64(skip)).SetLimit(int64(config.BatchSize)))
		if err != nil {
			return fmt.Errorf("查询章节失败: %w", err)
		}

		var chapters []OldChapter
		if err := cursor.All(ctx, &chapters); err != nil {
			cursor.Close(ctx)
			return fmt.Errorf("解析章节数据失败: %w", err)
		}
		cursor.Close(ctx)

		if len(chapters) == 0 {
			break // 没有更多数据
		}

		fmt.Printf("批次 #%d: 处理 %d 个章节...\n", batchNum, len(chapters))

		// 处理这一批数据
		for i, chapter := range chapters {
			// 检查是否已经迁移过
			existing, _ := contentsCol.CountDocuments(ctx, bson.M{"chapter_id": chapter.ID})
			if existing > 0 {
				skippedCount++
				if (i+1)%20 == 0 {
					fmt.Printf("   进度: %d/%d (跳过已迁移)\n", i+1, len(chapters))
				}
				continue
			}

			// 创建 ChapterContent 文档
			now := time.Now()
			content := ChapterContent{
				ID:        primitive.NewObjectID(),
				ChapterID: chapter.ID,
				Content:   chapter.Content,
				Format:    "markdown",
				Version:   1,
				WordCount: len([]rune(chapter.Content)),
				Hash:      fmt.Sprintf("%s:%d", chapter.ID.Hex(), 1),
				CreatedAt: now,
				UpdatedAt: now,
			}

			// 更新 Chapter 文档（添加内容引用信息）
			chapterUpdates := bson.M{
				"content_url":      fmt.Sprintf("/api/v1/bookstore/chapters/%s/content", chapter.ID.Hex()),
				"content_size":     int64(len([]rune(chapter.Content))),
				"content_hash":     content.Hash,
				"content_version":  1,
				"updated_at":       now,
			}

			if !config.DryRun {
				// 写入 ChapterContent
				_, err := contentsCol.InsertOne(ctx, content)
				if err != nil {
					failedCount++
					if config.SkipErrors {
						log.Printf("   警告: 章节 %s 内容插入失败: %v", chapter.ID.Hex(), err)
						continue
					}
					return fmt.Errorf("插入章节 %s 内容失败: %w", chapter.ID.Hex(), err)
				}

				// 更新 Chapter
				_, err = chaptersCol.UpdateOne(ctx, bson.M{"_id": chapter.ID}, bson.M{"$set": chapterUpdates})
				if err != nil {
					failedCount++
					if config.SkipErrors {
						log.Printf("   警告: 章节 %s 更新失败: %v", chapter.ID.Hex(), err)
						continue
					}
					return fmt.Errorf("更新章节 %s 失败: %w", chapter.ID.Hex(), err)
				}
			} else {
				// 试运行 - 只打印信息
				fmt.Printf("   [DRY RUN] 章节 %s (%s) - %d 字符\n",
					chapter.ID.Hex(), chapter.Title, len([]rune(chapter.Content)))
			}

			migratedCount++
			if (i+1)%20 == 0 || !config.DryRun {
				fmt.Printf("   进度: %d/%d\n", i+1, len(chapters))
			}
		}

		skip += len(chapters)

		// 检查是否已经处理完
		if int64(migratedCount+skippedCount) >= totalCount {
			break
		}
	}

	// 3. 验证迁移结果
	fmt.Println("\n3. 验证迁移结果...")
	contentCount, _ := contentsCol.CountDocuments(ctx, bson.M{})
	remainingCount, _ := chaptersCol.CountDocuments(ctx, bson.M{"content": bson.M{"$exists": true, "$ne": ""}})

	fmt.Printf("   chapter_contents 集合文档数: %d\n", contentCount)
	fmt.Printf("   chapters 集合仍有 content 字段的文档数: %d\n", remainingCount)

	// 4. 打印统计
	fmt.Println("\n========================================")
	fmt.Println("迁移统计")
	fmt.Println("========================================")
	fmt.Printf("总数: %d\n", totalCount)
	fmt.Printf("已迁移: %d\n", migratedCount)
	fmt.Printf("跳过（已存在）: %d\n", skippedCount)
	fmt.Printf("失败: %d\n", failedCount)
	fmt.Printf("成功率: %.2f%%\n", float64(migratedCount)/float64(totalCount)*100)
	fmt.Println("========================================")

	if remainingCount > 0 {
		fmt.Println("\n警告：仍有章节包含 content 字段，可能需要重新运行迁移")
	}

	// 5. 生成迁移报告
	report := map[string]interface{}{
		"timestamp":       time.Now().Format(time.RFC3339),
		"total":           totalCount,
		"migrated":        migratedCount,
		"skipped":         skippedCount,
		"failed":          failedCount,
		"success_rate":    float64(migratedCount) / float64(totalCount) * 100,
		"content_count":   contentCount,
		"remaining_count": remainingCount,
	}

	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	fmt.Printf("\n迁移报告（JSON）:\n%s\n", string(reportJSON))

	return nil
}
