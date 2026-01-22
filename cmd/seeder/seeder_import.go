// Package main 提供小说导入功能
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ImportSeeder 小说导入器
type ImportSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewImportSeeder 创建小说导入器
func NewImportSeeder(db *utils.Database, cfg *config.Config) *ImportSeeder {
	return &ImportSeeder{
		db:     db,
		config: cfg,
	}
}

// NovelData 从外部导入的小说数据结构
type NovelData struct {
	Metadata NovelsMetadata `json:"metadata"`
	Novels   []NovelItem    `json:"novels"`
}

// NovelsMetadata 元数据
type NovelsMetadata struct {
	Source        string    `json:"source"`
	TotalNovels   int       `json:"total_novels"`
	TotalChapters int       `json:"total_chapters"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// NovelItem 单本小说数据
type NovelItem struct {
	Title        string        `json:"title"`
	Author       string        `json:"author"`
	Introduction string        `json:"introduction"`
	Category     string        `json:"category"`
	WordCount    int64         `json:"word_count"`
	ChapterCount int           `json:"chapter_count"`
	Rating       float64       `json:"rating"`
	Status       string        `json:"status"`
	IsFree       bool          `json:"is_free"`
	Chapters     []ChapterItem `json:"chapters"`
}

// ChapterItem 章节数据
type ChapterItem struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	WordCount int    `json:"word_count"`
}

// SeedFromJSON 从JSON文件导入小说
func (s *ImportSeeder) SeedFromJSON(filepath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("数据文件不存在: %s", filepath)
	}

	fmt.Printf("  从文件导入小说: %s\n", filepath)

	// 读取 JSON 文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	var novelData NovelData
	if err := json.Unmarshal(data, &novelData); err != nil {
		return fmt.Errorf("解析 JSON 失败: %w", err)
	}

	fmt.Printf("  元数据: 来源=%s, 小说数=%d, 章节数=%d\n",
		novelData.Metadata.Source,
		novelData.Metadata.TotalNovels,
		novelData.Metadata.TotalChapters)

	// 导入书籍
	ctx := context.Background()
	if err := s.importNovels(ctx, novelData.Novels); err != nil {
		return err
	}

	fmt.Println("  小说导入成功")

	// 创建索引
	fmt.Println("  创建索引...")
	if err := s.createIndexes(ctx); err != nil {
		fmt.Printf("  ⚠️  创建索引失败: %v\n", err)
	} else {
		fmt.Println("  索引创建成功")
	}

	return nil
}

// importNovels 导入小说数据
func (s *ImportSeeder) importNovels(ctx context.Context, novels []NovelItem) error {
	bookCollection := s.db.Collection("books")
	chapterCollection := s.db.Collection("chapters")
	contentCollection := s.db.Collection("chapter_contents")

	now := time.Now()
	successCount := 0

	for _, novel := range novels {
		// 创建书籍
		bookID := primitive.NewObjectID()
		book := bson.M{
			"_id":           bookID,
			"title":         novel.Title,
			"author":        novel.Author,
			"introduction":  novel.Introduction,
			"cover":         "/images/covers/default.jpg",
			"categories":    []string{novel.Category},
			"tags":          []string{},
			"status":        novel.Status,
			"rating":        novel.Rating,
			"rating_count":  0,
			"view_count":    0,
			"word_count":    novel.WordCount,
			"chapter_count": len(novel.Chapters),
			"price":         0,
			"is_free":       novel.IsFree,
			"is_recommended": false,
			"is_featured":   false,
			"is_hot":        false,
			"published_at":  now.Add(-time.Duration(len(novel.Chapters)) * 24 * time.Hour),
			"last_update_at": now,
			"created_at":    now,
			"updated_at":    now,
		}

		// 插入书籍
		_, err := bookCollection.InsertOne(ctx, book)
		if err != nil {
			fmt.Printf("  ⚠️  插入书籍失败: %s - %v\n", novel.Title, err)
			continue
		}

		// 插入章节
		var chapters []interface{}
		var contents []interface{}

		for i, chapterItem := range novel.Chapters {
			chapterID := primitive.NewObjectID()

			// 章节
			chapter := bson.M{
				"_id":         chapterID,
				"book_id":     bookID.Hex(),
				"chapter_num": i + 1,
				"title":       chapterItem.Title,
				"word_count":  chapterItem.WordCount,
				"price":       0,
				"is_free":     true,
				"status":      "published",
				"published_at": now.Add(time.Duration(i) * 24 * time.Hour),
				"created_at":  now,
				"updated_at":  now,
			}
			chapters = append(chapters, chapter)

			// 章节内容
			content := bson.M{
				"chapter_id": chapterID.Hex(),
				"content":    chapterItem.Content,
				"word_count": chapterItem.WordCount,
				"created_at": now,
				"updated_at": now,
			}
			contents = append(contents, content)
		}

		// 批量插入章节
		if len(chapters) > 0 {
			_, err = chapterCollection.InsertMany(ctx, chapters)
			if err != nil {
				fmt.Printf("  ⚠️  插入章节失败: %s - %v\n", novel.Title, err)
				continue
			}
		}

		// 批量插入章节内容
		if len(contents) > 0 {
			_, err = contentCollection.InsertMany(ctx, contents)
			if err != nil {
				fmt.Printf("  ⚠️  插入章节内容失败: %s - %v\n", novel.Title, err)
				continue
			}
		}

		successCount++
		fmt.Printf("  ✓ %s (%d章)\n", novel.Title, len(novel.Chapters))
	}

	fmt.Printf("  成功导入 %d 本小说\n", successCount)
	return nil
}

// createIndexes 创建索引
func (s *ImportSeeder) createIndexes(ctx context.Context) error {
	// 为书籍创建索引
	bookCollection := s.db.Collection("books")

	// 创建 title 文本索引
	_, err := bookCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: "text"}},
		Options: options.Index().SetName("title_text"),
	})
	if err != nil {
		// 索引可能已存在，忽略错误
	}

	_, err = bookCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "author", Value: 1}},
		Options: options.Index().SetName("author_1"),
	})
	if err != nil {
		// 索引可能已存在，忽略错误
	}

	_, err = bookCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "categories", Value: 1}},
		Options: options.Index().SetName("categories_1"),
	})
	if err != nil {
		// 索引可能已存在，忽略错误
	}

	// 为章节创建索引
	chapterCollection := s.db.Collection("chapters")

	_, err = chapterCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "book_id", Value: 1}},
		Options: options.Index().SetName("book_id_1"),
	})
	if err != nil {
		// 索引可能已存在，忽略错误
	}

	_, err = chapterCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "book_id", Value: 1}, {Key: "chapter_num", Value: 1}},
		Options: options.Index().SetName("book_id_chapter_num_1"),
	})
	if err != nil {
		// 索引可能已存在，忽略错误
	}

	return nil
}

// Clean 清空导入的小说数据
func (s *ImportSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("books").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 books 集合失败: %w", err)
	}

	_, err = s.db.Collection("chapters").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 chapters 集合失败: %w", err)
	}

	_, err = s.db.Collection("chapter_contents").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 chapter_contents 集合失败: %w", err)
	}

	fmt.Println("  已清空 books, chapters 和 chapter_contents 集合")
	return nil
}

// GetStats 获取导入统计信息
func (s *ImportSeeder) GetStats() error {
	ctx := context.Background()

	booksCount, _ := s.db.Collection("books").CountDocuments(ctx, bson.M{})
	chaptersCount, _ := s.db.Collection("chapters").CountDocuments(ctx, bson.M{})
	contentsCount, _ := s.db.Collection("chapter_contents").CountDocuments(ctx, bson.M{})

	fmt.Printf("  书籍: %d 本\n", booksCount)
	fmt.Printf("  章节: %d 章\n", chaptersCount)
	fmt.Printf("  内容: %d 条\n", contentsCount)

	return nil
}
