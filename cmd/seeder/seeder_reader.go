// Package main 提供阅读数据填充功能
package main

import (
	"context"
	"fmt"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/generators"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReaderSeeder 阅读数据填充器
type ReaderSeeder struct {
	db       *utils.Database
	config   *config.Config
	gen      *generators.ReaderGenerator
	inserter *utils.BulkInserter
}

// NewReaderSeeder 创建阅读数据填充器
func NewReaderSeeder(db *utils.Database, cfg *config.Config) *ReaderSeeder {
	return &ReaderSeeder{
		db:     db,
		config: cfg,
		gen:    generators.NewReaderGenerator(),
	}
}

// SeedReadingData 填充所有阅读数据
func (s *ReaderSeeder) SeedReadingData() error {
	ctx := context.Background()

	// 获取用户和书籍
	users, err := s.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	books, err := s.getBookIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}

	chapters, err := s.getChapterIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取章节列表失败: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("  没有找到用户，请先运行 users 命令创建用户")
		return nil
	}

	if len(books) == 0 {
		fmt.Println("  没有找到书籍，请先运行 bookstore 命令创建书籍")
		return nil
	}

	if len(chapters) == 0 {
		fmt.Println("  没有找到章节，请先运行 chapters 命令创建章节")
		return nil
	}

	// 填充阅读进度
	if err := s.seedReadingProgresses(ctx, users, books); err != nil {
		return err
	}

	// 填充阅读历史
	if err := s.seedReadingHistories(ctx, users, books, chapters); err != nil {
		return err
	}

	// 填充书签
	if err := s.seedBookmarks(ctx, users, books, chapters); err != nil {
		return err
	}

	// 填充批注
	if err := s.seedAnnotations(ctx, users, books, chapters); err != nil {
		return err
	}

	return nil
}

// getUserIDs 获取用户ID列表
func (s *ReaderSeeder) getUserIDs(ctx context.Context) ([]string, error) {
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}
	return userIDs, nil
}

// getBookIDs 获取书籍ID列表
func (s *ReaderSeeder) getBookIDs(ctx context.Context) ([]string, error) {
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	bookIDs := make([]string, len(books))
	for i, b := range books {
		bookIDs[i] = b.ID
	}
	return bookIDs, nil
}

// getChapterIDs 获取章节ID列表
func (s *ReaderSeeder) getChapterIDs(ctx context.Context) ([]string, error) {
	cursor, err := s.db.Collection("chapters").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &chapters); err != nil {
		return nil, err
	}

	chapterIDs := make([]string, len(chapters))
	for i, c := range chapters {
		chapterIDs[i] = c.ID
	}
	return chapterIDs, nil
}

// seedReadingProgresses 创建阅读进度
func (s *ReaderSeeder) seedReadingProgresses(ctx context.Context, users, books []string) error {
	collection := s.db.Collection("reading_progress")

	// 70%的用户有阅读进度
	userCount := int(float64(len(users)) * 0.7)
	activeUsers := users[:userCount]

	// 用于跟踪已生成的用户-书籍组合，避免重复
	progressMap := make(map[string]bool) // key: "userID-bookID"

	for _, user := range activeUsers {
		// 每个用户平均阅读 3-10 本书
		bookCount := 3 + len(books)%8
		userBooks := s.getRandomItems(books, bookCount)

		for _, book := range userBooks {
			// 检查是否已存在该用户-书籍组合
			key := user + "-" + book
			if progressMap[key] {
				continue // 跳过已存在的组合
			}

			// 生成单个阅读进度
			progress := s.gen.GenerateReadingProgress(user, book)
			progressMap[key] = true

			// 使用 Upsert 插入或更新，确保唯一性
			filter := bson.M{"user_id": user, "book_id": book}
			update := bson.M{"$set": progress}

			_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
			if err != nil {
				return fmt.Errorf("插入阅读进度失败: %w", err)
			}
		}
	}

	fmt.Printf("  创建了 %d 条阅读进度数据\n", len(progressMap))
	return nil
}

// seedReadingHistories 创建阅读历史
func (s *ReaderSeeder) seedReadingHistories(ctx context.Context, users, books, chapters []string) error {
	collection := s.db.Collection("reading_history")

	// 为70%的用户生成阅读历史
	userCount := len(users) * 7 / 10
	historyCount := userCount * 3 // 每个用户平均3条历史

	// 使用生成器生成阅读历史
	histories := s.gen.GenerateReadingHistories(users, books, chapters, historyCount)

	if len(histories) > 0 {
		if _, err := collection.InsertMany(ctx, s.toInterfaceSlice(histories)); err != nil {
			return fmt.Errorf("插入阅读历史失败: %w", err)
		}
		fmt.Printf("  创建了 %d 条阅读历史\n", len(histories))
	}

	return nil
}

// seedBookmarks 创建书签
func (s *ReaderSeeder) seedBookmarks(ctx context.Context, users, books, chapters []string) error {
	collection := s.db.Collection("bookmarks")

	// 20%的活跃用户添加书签
	userCount := len(users) / 5
	bookmarkCount := userCount * 2 // 每个用户平均2个书签

	// 使用生成器生成书签
	bookmarks := s.gen.GenerateBookmarks(users, books, chapters, bookmarkCount)

	if len(bookmarks) > 0 {
		if _, err := collection.InsertMany(ctx, s.toInterfaceSlice(bookmarks)); err != nil {
			return fmt.Errorf("插入书签失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个书签\n", len(bookmarks))
	}

	return nil
}

// seedAnnotations 创建批注
func (s *ReaderSeeder) seedAnnotations(ctx context.Context, users, books, chapters []string) error {
	collection := s.db.Collection("annotations")

	// 10%的活跃用户添加批注
	userCount := len(users) / 10
	annotationCount := userCount // 每个用户平均1个批注

	// 使用生成器生成批注
	annotations := s.gen.GenerateAnnotations(users, books, chapters, annotationCount)

	if len(annotations) > 0 {
		if _, err := collection.InsertMany(ctx, s.toInterfaceSlice(annotations)); err != nil {
			return fmt.Errorf("插入批注失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个批注\n", len(annotations))
	}

	return nil
}

// getRandomItems 从数组中随机获取指定数量的元素
func (s *ReaderSeeder) getRandomItems(items []string, count int) []string {
	if count > len(items) {
		count = len(items)
	}
	return items[:count]
}

// toInterfaceSlice 转换为interface{}切片
func (s *ReaderSeeder) toInterfaceSlice(items interface{}) []interface{} {
	v := make([]interface{}, 0)
	switch items := items.(type) {
	case []models.ReadingProgress:
		for _, item := range items {
			v = append(v, item)
		}
	case []models.ReadingHistory:
		for _, item := range items {
			v = append(v, item)
		}
	case []models.Bookmark:
		for _, item := range items {
			v = append(v, item)
		}
	case []models.Annotation:
		for _, item := range items {
			v = append(v, item)
		}
	}
	return v
}

// Clean 清空阅读数据
func (s *ReaderSeeder) Clean() error {
	ctx := context.Background()

	collections := []string{"reading_progress", "reading_history", "bookmarks", "annotations"}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空阅读数据集合")
	return nil
}
