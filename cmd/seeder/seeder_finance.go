// Package main 提供财务数据填充功能
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

// FinanceSeeder 财务数据填充器
type FinanceSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewFinanceSeeder 创建财务数据填充器
func NewFinanceSeeder(db *utils.Database, cfg *config.Config) *FinanceSeeder {
	return &FinanceSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedFinanceData 填充财务数据
func (s *FinanceSeeder) SeedFinanceData() error {
	ctx := context.Background()

	// 获取用户和书籍列表
	users, err := s.getUserIDs(ctx, "author")
	if err != nil {
		return fmt.Errorf("获取作者列表失败: %w", err)
	}

	allUsers, err := s.getUserIDs(ctx, "")
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	books, err := s.getBookIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}

	// 填充作者收入
	if err := s.seedAuthorRevenue(ctx, users, books); err != nil {
		return err
	}

	// 填充会员
	if err := s.seedMemberships(ctx, allUsers); err != nil {
		return err
	}

	return nil
}

// seedAuthorRevenue 填充作者收入
func (s *FinanceSeeder) seedAuthorRevenue(ctx context.Context, userIDs, bookIDs []string) error {
	collection := s.db.Collection("author_revenue")

	var revenues []interface{}
	now := time.Now()

	for _, userID := range userIDs {
		// 每个作者20-100条收入记录
		revenueCount := 20 + len(userID)%81

		for i := 0; i < revenueCount; i++ {
			bookID := bookIDs[i%len(bookIDs)]

			// 收入类型: 章节购买70%, 订阅20%, 打赏10%
			rand := i % 10
			var revType string
			var amount float64

			switch {
			case rand < 7:
				revType = "chapter"
				amount = 0.1 + float64(i%10)*0.1
			case rand < 9:
				revType = "subscription"
				amount = 1.0 + float64(i%5)
			default:
				revType = "reward"
				amount = 5.0 + float64(i%10)*5
			}

			// 80%已结算，20%待结算
			status := "settled"
			settledAt := now.Add(-time.Duration(i) * time.Hour)
			if rand > 7 {
				status = "pending"
				settledAt = time.Time{}
			}

			period := now.Add(-time.Duration(i/5) * 24 * time.Hour).Format("2006-01")

			revenues = append(revenues, models.AuthorRevenue{
				ID:        primitive.NewObjectID().Hex(),
				UserID:    userID,
				BookID:    bookID,
				Amount:    amount,
				Type:      revType,
				Status:    status,
				Period:    period,
				CreatedAt: now.Add(-time.Duration(i*24) * time.Hour),
				SettledAt: settledAt,
			})
		}
	}

	// 批量插入
	if len(revenues) > 0 {
		batchSize := 100
		for i := 0; i < len(revenues); i += batchSize {
			end := i + batchSize
			if end > len(revenues) {
				end = len(revenues)
			}

			_, err := collection.InsertMany(ctx, revenues[i:end])
			if err != nil {
				return fmt.Errorf("插入作者收入失败（批次 %d）: %w", i/batchSize, err)
			}
		}

		fmt.Printf("  创建了 %d 条作者收入记录\n", len(revenues))
	}

	return nil
}

// seedMemberships 填充会员
func (s *FinanceSeeder) seedMemberships(ctx context.Context, userIDs []string) error {
	collection := s.db.Collection("memberships")

	var memberships []interface{}
	now := time.Now()

	// 10%的用户是会员
	memberCount := len(userIDs) / 10
	if memberCount < 1 {
		memberCount = 1
	}

	for i := 0; i < memberCount; i++ {
		userID := userIDs[i]

		// 会员类型: 月会员70%, 年会员25%, 终身会员5%
		rand := i % 100
		var memType string
		var duration time.Duration

		switch {
		case rand < 70:
			memType = "monthly"
			duration = 30 * 24 * time.Hour
		case rand < 95:
			memType = "yearly"
			duration = 365 * 24 * time.Hour
		default:
			memType = "lifetime"
			duration = 10 * 365 * 24 * time.Hour
		}

		startAt := now.Add(-time.Duration(i) * 24 * time.Hour)
		expireAt := startAt.Add(duration)

		// 检查是否过期
		status := "active"
		if expireAt.Before(now) {
			status = "expired"
		}

		memberships = append(memberships, models.Membership{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    userID,
			Type:      memType,
			Status:    status,
			StartAt:   startAt,
			ExpireAt:  expireAt,
			CreatedAt: startAt,
		})
	}

	// 批量插入
	if len(memberships) > 0 {
		_, err := collection.InsertMany(ctx, memberships)
		if err != nil {
			return fmt.Errorf("插入会员失败: %w", err)
		}

		fmt.Printf("  创建了 %d 个会员\n", len(memberships))
	}

	return nil
}

// getUserIDs 获取用户ID列表
func (s *FinanceSeeder) getUserIDs(ctx context.Context, role string) ([]string, error) {
	filter := bson.M{}
	if role != "" {
		filter = bson.M{"role": role}
	}

	cursor, err := s.db.Collection("users").Find(ctx, filter)
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
func (s *FinanceSeeder) getBookIDs(ctx context.Context) ([]string, error) {
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

// Clean 清空财务数据
func (s *FinanceSeeder) Clean() error {
	ctx := context.Background()

	collections := []string{"author_revenue", "memberships"}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空财务数据集合")
	return nil
}
