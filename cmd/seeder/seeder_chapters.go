// Package main 提供章节数据填充功能
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/generators"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterSeeder 章节数据填充器
type ChapterSeeder struct {
	db       *utils.Database
	config   *config.Config
	baseGen  *generators.BaseGenerator
	inserter *utils.BulkInserter
}

// NewChapterSeeder 创建章节填充器
func NewChapterSeeder(db *utils.Database, cfg *config.Config) *ChapterSeeder {
	return &ChapterSeeder{
		db:       db,
		config:   cfg,
		baseGen:  generators.NewBaseGenerator(),
		inserter: utils.NewBulkInserter(db.Collection("chapters"), cfg.BatchSize),
	}
}

// SeedChapters 为现有书籍生成章节数据
func (s *ChapterSeeder) SeedChapters() error {
	ctx := context.Background()

	// 获取所有书籍
	booksCollection := s.db.Collection("books")
	cursor, err := booksCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err = cursor.All(ctx, &books); err != nil {
		return fmt.Errorf("解析书籍列表失败: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("  没有找到书籍，请先运行 bookstore 命令创建书籍")
		return nil
	}

	fmt.Printf("  找到 %d 本书，开始生成章节...\n", len(books))

	var allChapters []models.Chapter
	totalChapters := 0

	// 为每本书生成章节
	for _, book := range books {
		// 根据书籍状态决定章节数量
		chapterCount := s.determineChapterCount(book)

		// 生成章节
		chapters := s.generateChaptersForBook(book, chapterCount)
		allChapters = append(allChapters, chapters...)
		totalChapters += len(chapters)

		fmt.Printf("  为《%s》生成 %d 章\n", book.Title, len(chapters))
	}

	// 批量插入章节
	if err := s.inserter.InsertMany(ctx, allChapters); err != nil {
		return fmt.Errorf("插入章节失败: %w", err)
	}

	fmt.Printf("  成功插入 %d 章\n", totalChapters)
	return nil
}

// SeedChapterContents 为章节生成内容
func (s *ChapterSeeder) SeedChapterContents() error {
	ctx := context.Background()

	// 获取所有章节
	chapterCollection := s.db.Collection("chapters")
	cursor, err := chapterCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("获取章节列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var chapters []models.Chapter
	if err = cursor.All(ctx, &chapters); err != nil {
		return fmt.Errorf("解析章节列表失败: %w", err)
	}

	if len(chapters) == 0 {
		fmt.Println("  没有找到章节")
		return nil
	}

	fmt.Printf("  找到 %d 章，开始生成内容...\n", len(chapters))

	// 生成章节内容
	contentCollection := s.db.Collection("chapter_contents")
	var allContents []interface{}

	for _, chapter := range chapters {
		content := models.ChapterContent{
			ChapterID: chapter.ID,
			Content:   s.baseGen.ChapterContent(800, 1500),
			WordCount: 1000 + rand.Intn(500),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		allContents = append(allContents, content)
	}

	// 批量插入（每次100条）
	batchSize := 100
	for i := 0; i < len(allContents); i += batchSize {
		end := i + batchSize
		if end > len(allContents) {
			end = len(allContents)
		}

		_, err := contentCollection.InsertMany(ctx, allContents[i:end])
		if err != nil {
			return fmt.Errorf("插入章节内容失败（批次 %d）: %w", i/batchSize, err)
		}
	}

	fmt.Printf("  成功生成 %d 章内容\n", len(allContents))
	return nil
}

// determineChapterCount 根据配置的数据规模决定章节数量
func (s *ChapterSeeder) determineChapterCount(book models.Book) int {
	// 如果书籍已有章节数，使用它
	if book.ChapterCount > 0 {
		return book.ChapterCount
	}

	// 根据配置的数据规模决定章节数量
	scale := config.GetScaleConfig(s.config.Scale)

	// 生成 MinChapters 到 MaxChapters 之间的随机章节数量
	// rand.Intn(n) 返回 [0, n) 范围内的随机数，所以需要 +1 确保包含 MaxChapters
	return scale.MinChapters + rand.Intn(scale.MaxChapters-scale.MinChapters+1)
}

// generateChaptersForBook 为单本书生成章节
func (s *ChapterSeeder) generateChaptersForBook(book models.Book, count int) []models.Chapter {
	now := time.Now()
	publishedAt := book.PublishedAt

	chapters := make([]models.Chapter, count)
	for i := 0; i < count; i++ {
		chapterNum := i + 1
		isFree := chapterNum <= 10 // 前10章免费

		chapters[i] = models.Chapter{
			ID:          primitive.NewObjectID(),
			BookID:      book.ID,
			ChapterNum:  chapterNum,
			Title:       s.generateChapterTitle(chapterNum),
			WordCount:   1000 + rand.Intn(1000),
			Price:       s.calculateChapterPrice(chapterNum, isFree),
			IsFree:      isFree,
			Status:      "published",
			PublishedAt: publishedAt.Add(time.Duration(chapterNum) * 24 * time.Hour),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	return chapters
}

// generateChapterTitle 生成章节标题
func (s *ChapterSeeder) generateChapterTitle(chapterNum int) string {
	titles := []string{
		"第一章", "第二章", "第三章", "第四章", "第五章",
		"第六章", "第七章", "第八章", "第九章", "第十章",
	}

	if chapterNum <= 10 {
		return fmt.Sprintf("%s %s", titles[chapterNum-1], s.getRandomTitleSuffix())
	}

	return fmt.Sprintf("第%d章 %s", chapterNum, s.getRandomTitleSuffix())
}

// getRandomTitleSuffix 获取随机标题后缀
func (s *ChapterSeeder) getRandomTitleSuffix() string {
	suffixes := []string{
		"初入江湖", "机缘巧合", "实力大增", "遭遇强敌", "突破境界",
		"险象环生", "绝地反击", "获得传承", "扬名立万", "再攀高峰",
		"风云变幻", "暗流涌动", "一触即发", "生死时刻", "峰回路转",
	}
	return suffixes[rand.Intn(len(suffixes))]
}

// calculateChapterPrice 计算章节价格
func (s *ChapterSeeder) calculateChapterPrice(chapterNum int, isFree bool) float64 {
	if isFree {
		return 0
	}

	// 章节价格逐渐增加
	basePrice := 0.1
	if chapterNum > 50 {
		basePrice = 0.15
	}
	if chapterNum > 100 {
		basePrice = 0.2
	}

	return basePrice
}

// Clean 清空章节数据
func (s *ChapterSeeder) Clean() error {
	ctx := context.Background()

	// 清空章节
	_, err := s.db.Collection("chapters").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 chapters 集合失败: %w", err)
	}

	// 清空章节内容
	_, err = s.db.Collection("chapter_contents").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 chapter_contents 集合失败: %w", err)
	}

	fmt.Println("  已清空 chapters 和 chapter_contents 集合")
	return nil
}

// Count 统计章节数量
func (s *ChapterSeeder) Count() (int64, error) {
	ctx := context.Background()
	count, err := s.db.Collection("chapters").CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("统计章节数量失败: %w", err)
	}
	return count, nil
}
