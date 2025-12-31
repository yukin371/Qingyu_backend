package seeds

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NovelData 从 Python 脚本导出的小说数据结构
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
	ChapterSize   int       `json:"chapter_size"`
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

// NovelImporter 小说导入器
type NovelImporter struct {
	db                *mongo.Database
	bookCollection    *mongo.Collection
	chapterCollection *mongo.Collection
	dryRun            bool
}

// NewNovelImporter 创建小说导入器
func NewNovelImporter(db *mongo.Database, dryRun bool) *NovelImporter {
	return &NovelImporter{
		db:                db,
		bookCollection:    db.Collection("books"),
		chapterCollection: db.Collection("chapters"),
		dryRun:            dryRun,
	}
}

// ImportFromJSON 从 JSON 文件导入小说
func (ni *NovelImporter) ImportFromJSON(ctx context.Context, filepath string) error {
	log.Printf("开始从文件导入小说: %s", filepath)

	// 读取 JSON 文件
	data, err := ni.loadJSONFile(filepath)
	if err != nil {
		return fmt.Errorf("读取 JSON 文件失败: %w", err)
	}

	log.Printf("元数据:")
	log.Printf("  来源: %s", data.Metadata.Source)
	log.Printf("  小说数量: %d", data.Metadata.TotalNovels)
	log.Printf("  章节数量: %d", data.Metadata.TotalChapters)
	log.Printf("  生成时间: %s", data.Metadata.GeneratedAt.Format("2006-01-02 15:04:05"))

	if ni.dryRun {
		log.Println("\n[试运行模式] 不会实际写入数据库")
		return ni.validateData(data)
	}

	// 开始导入
	return ni.importNovels(ctx, data.Novels)
}

// loadJSONFile 加载 JSON 文件
func (ni *NovelImporter) loadJSONFile(filepath string) (*NovelData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var data NovelData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return &data, nil
}

// validateData 验证数据
func (ni *NovelImporter) validateData(data *NovelData) error {
	log.Println("\n开始验证数据...")

	invalidCount := 0
	for i, novel := range data.Novels {
		if err := ni.validateNovel(&novel); err != nil {
			log.Printf("  [%d] %s - 验证失败: %v", i+1, novel.Title, err)
			invalidCount++
		}
	}

	if invalidCount > 0 {
		return fmt.Errorf("发现 %d 本小说数据无效", invalidCount)
	}

	log.Printf("✓ 所有 %d 本小说数据验证通过", len(data.Novels))
	return nil
}

// validateNovel 验证单本小说
func (ni *NovelImporter) validateNovel(novel *NovelItem) error {
	if novel.Title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if len(novel.Title) > 100 {
		return fmt.Errorf("标题过长: %d 字符", len(novel.Title))
	}
	if novel.Author == "" {
		return fmt.Errorf("作者不能为空")
	}
	if len(novel.Chapters) == 0 {
		return fmt.Errorf("章节列表为空")
	}
	if novel.ChapterCount != len(novel.Chapters) {
		return fmt.Errorf("章节数量不匹配: 声明 %d，实际 %d", novel.ChapterCount, len(novel.Chapters))
	}
	return nil
}

// importNovels 导入小说列表
func (ni *NovelImporter) importNovels(ctx context.Context, novels []NovelItem) error {
	log.Printf("\n开始导入 %d 本小说...\n", len(novels))

	successCount := 0
	errorCount := 0

	for i, novel := range novels {
		if err := ni.importSingleNovel(ctx, &novel); err != nil {
			log.Printf("  [%d/%d] ✗ %s - 导入失败: %v", i+1, len(novels), novel.Title, err)
			errorCount++
		} else {
			successCount++
			if successCount%50 == 0 {
				log.Printf("  已成功导入 %d 本...", successCount)
			}
		}
	}

	log.Printf("\n导入完成:")
	log.Printf("  ✓ 成功: %d 本", successCount)
	log.Printf("  ✗ 失败: %d 本", errorCount)

	if errorCount > 0 {
		return fmt.Errorf("部分小说导入失败")
	}

	return nil
}

// importSingleNovel 导入单本小说
func (ni *NovelImporter) importSingleNovel(ctx context.Context, novel *NovelItem) error {
	// 验证数据
	if err := ni.validateNovel(novel); err != nil {
		return err
	}

	// 创建书籍记录
	book := ni.convertToBook(novel)

	// 插入书籍
	result, err := ni.bookCollection.InsertOne(ctx, book)
	if err != nil {
		return fmt.Errorf("插入书籍失败: %w", err)
	}

	bookID := result.InsertedID.(primitive.ObjectID)

	// 批量插入章节
	if err := ni.importChapters(ctx, bookID, novel.Chapters); err != nil {
		// 如果章节导入失败，删除已插入的书籍
		_, _ = ni.bookCollection.DeleteOne(ctx, primitive.M{"_id": bookID})
		return fmt.Errorf("导入章节失败: %w", err)
	}

	return nil
}

// convertToBook 转换为 Book 模型
func (ni *NovelImporter) convertToBook(novel *NovelItem) *bookstore2.Book {
	now := time.Now()
	publishedAt := now.Add(-30 * 24 * time.Hour) // 假设30天前发布

	// 解析状态
	var status bookstore2.BookStatus
	switch novel.Status {
	case "completed":
		status = bookstore2.BookStatusCompleted
	case "ongoing":
		status = bookstore2.BookStatusOngoing
	default:
		status = bookstore2.BookStatusPublished
	}

	// 默认封面
	defaultCover := "https://via.placeholder.com/300x400?text=" + novel.Title

	book := &bookstore2.Book{
		Title:         novel.Title,
		Author:        novel.Author,
		Introduction:  novel.Introduction,
		Cover:         defaultCover,
		Categories:    []string{novel.Category},
		Tags:          []string{novel.Category},
		Status:        status,
		WordCount:     novel.WordCount,
		ChapterCount:  novel.ChapterCount,
		Price:         0, // 免费
		IsFree:        novel.IsFree,
		IsRecommended: novel.Rating >= 4.5,      // 评分>=4.5的推荐
		IsFeatured:    novel.Rating >= 4.8,      // 评分>=4.8的精选
		IsHot:         novel.WordCount > 500000, // 字数>50万的标记为热门
		PublishedAt:   &publishedAt,
		LastUpdateAt:  &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return book
}

// importChapters 批量导入章节
func (ni *NovelImporter) importChapters(ctx context.Context, bookID primitive.ObjectID, chapters []ChapterItem) error {
	if len(chapters) == 0 {
		return nil
	}

	// 批量插入，每批500章
	batchSize := 500
	now := time.Now()

	for i := 0; i < len(chapters); i += batchSize {
		end := i + batchSize
		if end > len(chapters) {
			end = len(chapters)
		}

		batch := chapters[i:end]
		chapterDocs := make([]interface{}, 0, len(batch))

		for j, chapterData := range batch {
			chapterNum := i + j + 1
			// 注意：Content字段已从Chapter模型中移除，内容现在存储在ChapterContent中
			// 如果需要导入内容，应该使用ChapterContent模型
			chapter := &bookstore2.Chapter{
				BookID:      bookID,
				Title:       chapterData.Title,
				ChapterNum:  chapterNum,
				WordCount:   chapterData.WordCount,
				IsFree:      true,
				Price:       0,
				PublishTime: now,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			chapterDocs = append(chapterDocs, chapter)
		}

		if _, err := ni.chapterCollection.InsertMany(ctx, chapterDocs); err != nil {
			return fmt.Errorf("批量插入章节失败 (批次 %d-%d): %w", i, end, err)
		}
	}

	return nil
}

// CreateIndexes 创建索引
func (ni *NovelImporter) CreateIndexes(ctx context.Context) error {
	log.Println("创建索引...")

	// 书籍索引
	bookIndexes := []mongo.IndexModel{
		{
			Keys: primitive.D{
				{Key: "title", Value: "text"},
				{Key: "author", Value: "text"},
			},
		},
		{
			Keys: primitive.D{{Key: "categories", Value: 1}},
		},
		{
			Keys: primitive.D{{Key: "status", Value: 1}},
		},
		{
			Keys: primitive.D{{Key: "is_recommended", Value: 1}},
		},
		{
			Keys: primitive.D{{Key: "created_at", Value: -1}},
		},
	}

	if _, err := ni.bookCollection.Indexes().CreateMany(ctx, bookIndexes); err != nil {
		return fmt.Errorf("创建书籍索引失败: %w", err)
	}

	// 章节索引
	chapterIndexes := []mongo.IndexModel{
		{
			Keys: primitive.D{
				{Key: "book_id", Value: 1},
				{Key: "chapter_num", Value: 1},
			},
		},
		{
			Keys: primitive.D{{Key: "book_id", Value: 1}},
		},
	}

	if _, err := ni.chapterCollection.Indexes().CreateMany(ctx, chapterIndexes); err != nil {
		return fmt.Errorf("创建章节索引失败: %w", err)
	}

	log.Println("✓ 索引创建成功")
	return nil
}

// GetStats 获取导入统计
func (ni *NovelImporter) GetStats(ctx context.Context) error {
	bookCount, err := ni.bookCollection.CountDocuments(ctx, primitive.M{})
	if err != nil {
		return err
	}

	chapterCount, err := ni.chapterCollection.CountDocuments(ctx, primitive.M{})
	if err != nil {
		return err
	}

	log.Printf("\n数据库统计:")
	log.Printf("  书籍总数: %d", bookCount)
	log.Printf("  章节总数: %d", chapterCount)

	return nil
}
