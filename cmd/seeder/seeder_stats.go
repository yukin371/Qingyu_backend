// Package main 提供统计数据填充功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StatsSeeder 统计数据填充器
type StatsSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewStatsSeeder 创建统计数据填充器
func NewStatsSeeder(db *utils.Database, cfg *config.Config) *StatsSeeder {
	return &StatsSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedStats 填充统计数据
func (s *StatsSeeder) SeedStats() error {
	ctx := context.Background()

	// 获取书籍和章节列表
	books, err := s.getBookIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}

	chapters, err := s.getChapterIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取章节列表失败: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("  没有找到书籍，请先运行 bookstore 命令创建书籍")
		return nil
	}

	// 填充书籍统计（最近30天）
	if err := s.seedBookStats(ctx, books, 30); err != nil {
		return err
	}

	// 填充章节统计（最近30天）
	if err := s.seedChapterStats(ctx, chapters, 30); err != nil {
		return err
	}

	return nil
}

// seedBookStats 填充书籍统计数据
func (s *StatsSeeder) seedBookStats(ctx context.Context, bookIDs []string, days int) error {
	collection := s.db.Collection("book_stats")

	var allStats []interface{}
	now := time.Now()

	for _, bookID := range bookIDs {
		// 为每本书生成最近30天的统计
		for day := 0; day < days; day++ {
			date := now.Add(-time.Duration(day) * 24 * time.Hour).Format("2006-01-02")

			// 模拟周期性波动（周末高，工作日低）
			weekday := now.Add(-time.Duration(day) * 24 * time.Hour).Weekday()
			weekendBoost := 1.0
			if weekday == time.Saturday || weekday == time.Sunday {
				weekendBoost = 1.5
			}

			baseViewCount := 50 + (day*7)%100
			viewCount := int(float64(baseViewCount) * weekendBoost)

			allStats = append(allStats, models.BookStats{
				ID:             primitive.NewObjectID().Hex(),
				BookID:         bookID,
				ViewCount:      viewCount,
				ReadCount:      viewCount * 6 / 10,
				FavoriteCount:  viewCount * 5 / 100,
				ShareCount:     viewCount / 20,
				AvgReadTime:    1800 + day*60,
				CompletionRate: 0.3 + float64(day%50)/100,
				Date:           date,
				CreatedAt:      now,
			})
		}
	}

	// 批量插入
	if len(allStats) > 0 {
		batchSize := 100
		for i := 0; i < len(allStats); i += batchSize {
			end := i + batchSize
			if end > len(allStats) {
				end = len(allStats)
			}

			_, err := collection.InsertMany(ctx, allStats[i:end])
			if err != nil {
				return fmt.Errorf("插入书籍统计失败（批次 %d）: %w", i/batchSize, err)
			}
		}

		fmt.Printf("  创建了 %d 条书籍统计记录\n", len(allStats))
	}

	return nil
}

// seedChapterStats 填充章节统计数据
func (s *StatsSeeder) seedChapterStats(ctx context.Context, chapterIDs []string, days int) error {
	collection := s.db.Collection("chapter_stats")

	var allStats []interface{}
	now := time.Now()

	// 为最近活跃的章节生成统计（避免数据量过大）
	activeChapterCount := len(chapterIDs)
	if activeChapterCount > 1000 {
		activeChapterCount = 1000
	}

	activeChapters := chapterIDs[:activeChapterCount]

	for _, chapterID := range activeChapters {
		// 为每个章节生成最近7天的统计
		for day := 0; day < 7; day++ {
			date := now.Add(-time.Duration(day) * 24 * time.Hour).Format("2006-01-02")

			baseViewCount := 20 + (day*5)%50

			allStats = append(allStats, models.ChapterStats{
				ID:        primitive.NewObjectID().Hex(),
				ChapterID: chapterID,
				BookID:    "", // 可以通过章节查询获取
				ViewCount: baseViewCount,
				ReadCount: baseViewCount * 8 / 10,
				StayTime:  180 + day*30,
				Date:      date,
				CreatedAt: now,
			})
		}
	}

	// 批量插入
	if len(allStats) > 0 {
		batchSize := 100
		for i := 0; i < len(allStats); i += batchSize {
			end := i + batchSize
			if end > len(allStats) {
				end = len(allStats)
			}

			_, err := collection.InsertMany(ctx, allStats[i:end])
			if err != nil {
				return fmt.Errorf("插入章节统计失败（批次 %d）: %w", i/batchSize, err)
			}
		}

		fmt.Printf("  创建了 %d 条章节统计记录\n", len(allStats))
	}

	return nil
}

// getBookIDs 获取书籍ID列表
func (s *StatsSeeder) getBookIDs(ctx context.Context) ([]string, error) {
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
func (s *StatsSeeder) getChapterIDs(ctx context.Context) ([]string, error) {
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

// Clean 清空统计数据
func (s *StatsSeeder) Clean() error {
	ctx := context.Background()

	collections := []string{"book_stats", "chapter_stats"}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空统计数据集合")
	return nil
}
