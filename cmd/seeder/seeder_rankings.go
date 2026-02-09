// Package main 提供榜单数据填充功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RankingSeeder 榜单数据填充器
type RankingSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewRankingSeeder 创建榜单数据填充器
func NewRankingSeeder(db *utils.Database, cfg *config.Config) *RankingSeeder {
	return &RankingSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedRankings 填充榜单数据
func (s *RankingSeeder) SeedRankings() error {
	ctx := context.Background()

	// 获取书籍列表
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID        string  `bson:"_id"`
		Title     string  `bson:"title"`
		Rating    float64 `bson:"rating"`
		ViewCount int64   `bson:"view_count"`
	}
	if err := cursor.All(ctx, &books); err != nil {
		return fmt.Errorf("解析书籍列表失败: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("  没有找到书籍，请先运行 bookstore 命令创建书籍")
		return nil
	}

	collection := s.db.Collection("ranking_items")

	// 清空现有榜单数据
	_, _ = collection.DeleteMany(ctx, bson.M{})

	now := time.Now()
	today := now.Format("2006-01-02")
	year, week := now.ISOWeek()
	weeklyPeriod := fmt.Sprintf("%d-W%02d", year, week)
	monthlyPeriod := now.Format("2006-01")

	var rankings []interface{}

	// 为每种榜单类型创建排名
	rankingTypes := []struct {
		Type        string
		Period      string
		Description string
		MaxRank     int
	}{
		{"realtime", today, "实时榜", 50},
		{"weekly", weeklyPeriod, "周榜", 50},
		{"monthly", monthlyPeriod, "月榜", 100},
		{"newbie", monthlyPeriod, "新人榜", 30},
	}

	for _, rankingType := range rankingTypes {
		// 根据榜单类型排序书籍
		sortedBooks := s.sortBooksForRanking(books, rankingType.Type)

		// 创建榜单项
		maxRank := rankingType.MaxRank
		if len(sortedBooks) < maxRank {
			maxRank = len(sortedBooks)
		}

		for i := 0; i < maxRank; i++ {
			book := sortedBooks[i]

			// 计算分数和统计数据
			score := s.calculateScore(book, i, rankingType.Type)
			viewCount := int64((maxRank - i) * 10000)
			likeCount := int64((maxRank - i) * 500)

			// 转换字符串ID为ObjectID
			bookID, err := primitive.ObjectIDFromHex(book.ID)
			if err != nil {
				return fmt.Errorf("转换书籍ID失败: %w", err)
			}

			rankings = append(rankings, bson.M{
				"_id":        primitive.NewObjectID(),
				"book_id":    bookID,
				"type":       rankingType.Type,
				"rank":       i + 1,
				"score":      score,
				"view_count": viewCount,
				"like_count": likeCount,
				"period":     rankingType.Period,
				"created_at": now,
				"updated_at": now,
			})
		}

		fmt.Printf("  创建 %s: %d 条\n", rankingType.Description, maxRank)
	}

	// 批量插入
	if len(rankings) > 0 {
		batchSize := 100
		for i := 0; i < len(rankings); i += batchSize {
			end := i + batchSize
			if end > len(rankings) {
				end = len(rankings)
			}
			_, err := collection.InsertMany(ctx, rankings[i:end])
			if err != nil {
				return fmt.Errorf("插入榜单失败: %w", err)
			}
		}
		fmt.Printf("  总计创建 %d 条榜单数据\n", len(rankings))
	}

	return nil
}

// sortBooksForRanking 根据榜单类型排序书籍
func (s *RankingSeeder) sortBooksForRanking(books []struct {
	ID        string  `bson:"_id"`
	Title     string  `bson:"title"`
	Rating    float64 `bson:"rating"`
	ViewCount int64   `bson:"view_count"`
}, rankingType string) []struct {
	ID        string  `bson:"_id"`
	Title     string  `bson:"title"`
	Rating    float64 `bson:"rating"`
	ViewCount int64   `bson:"view_count"`
} {
	// 简单排序：按评分和浏览量排序
	sorted := make([]struct {
		ID        string  `bson:"_id"`
		Title     string  `bson:"title"`
		Rating    float64 `bson:"rating"`
		ViewCount int64   `bson:"view_count"`
	}, len(books))
	copy(sorted, books)

	// 简单的冒泡排序（实际应用中应该使用更高效的排序）
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Rating < sorted[j+1].Rating {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// calculateScore 计算榜单分数
func (s *RankingSeeder) calculateScore(book struct {
	ID        string  `bson:"_id"`
	Title     string  `bson:"title"`
	Rating    float64 `bson:"rating"`
	ViewCount int64   `bson:"view_count"`
}, rank int, rankingType string) float64 {
	// 基础分数为评分
	score := book.Rating

	// 根据排名调整分数
	score += float64(100-rank) * 0.01

	// 根据榜单类型调整
	switch rankingType {
	case "realtime":
		score += float64(book.ViewCount/10000) * 0.1
	case "newbie":
		score += 0.5 // 新书奖励
	}

	// 确保分数在合理范围内
	if score > 10 {
		score = 10
	}
	if score < 0 {
		score = 0
	}

	return score
}

// Clean 清空榜单数据
func (s *RankingSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("ranking_items").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 ranking_items 集合失败: %w", err)
	}

	fmt.Println("  已清空 ranking_items 集合")
	return nil
}
