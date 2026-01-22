// Package main 提供AI配额激活功能
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

// AIQuotaSeeder AI配额激活器
type AIQuotaSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewAIQuotaSeeder 创建AI配额激活器
func NewAIQuotaSeeder(db *utils.Database, cfg *config.Config) *AIQuotaSeeder {
	return &AIQuotaSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedAIQuota 激活用户AI配额
func (s *AIQuotaSeeder) SeedAIQuota() error {
	ctx := context.Background()

	// 获取用户列表
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID   string `bson:"_id"`
		Role string `bson:"role"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return fmt.Errorf("解析用户列表失败: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("  没有找到用户，请先运行 users 命令创建用户")
		return nil
	}

	// 清空所有旧配额
	_, _ = s.db.Collection("ai_quotas").DeleteMany(ctx, bson.M{})

	collection := s.db.Collection("ai_quotas")
	var quotas []interface{}

	activatedCount := 0

	for _, user := range users {
		now := time.Now()
		resetAt := now.AddDate(0, 0, 1) // 明天重置

		// 根据角色设置配额
		var totalQuota int
		switch user.Role {
		case "admin":
			totalQuota = 999999 // 管理员：超大配额
		case "vip":
			totalQuota = 100000 // VIP：十万配额
		default:
			totalQuota = 10000 // 普通用户：一万配额
		}

		quota := bson.M{
			"_id":             primitive.NewObjectID(),
			"user_id":         user.ID,
			"quota_type":      "daily",
			"total_quota":     totalQuota,
			"used_quota":      0,
			"remaining_quota": totalQuota,
			"status":          "active",
			"reset_at":        resetAt,
			"last_reset_date": now,
			"created_at":      now,
			"updated_at":      now,
			"monthly_limit":   totalQuota * 30,
			"used_this_month": 0,
		}

		quotas = append(quotas, quota)
		activatedCount++
	}

	// 批量插入
	if len(quotas) > 0 {
		batchSize := 100
		for i := 0; i < len(quotas); i += batchSize {
			end := i + batchSize
			if end > len(quotas) {
				end = len(quotas)
			}
			_, err := collection.InsertMany(ctx, quotas[i:end])
			if err != nil {
				return fmt.Errorf("插入AI配额失败: %w", err)
			}
		}
	}

	fmt.Printf("  激活了 %d 个用户的AI配额\n", activatedCount)

	// 显示配额分布
	s.showQuotaDistribution(users)

	return nil
}

// showQuotaDistribution 显示配额分布
func (s *AIQuotaSeeder) showQuotaDistribution(users []struct {
	ID   string `bson:"_id"`
	Role string `bson:"role"`
}) {
	adminCount := 0
	vipCount := 0
	normalCount := 0

	for _, user := range users {
		switch user.Role {
		case "admin":
			adminCount++
		case "vip":
			vipCount++
		default:
			normalCount++
		}
	}

	fmt.Println("  配额分布:")
	fmt.Printf("    管理员: %d 人 × 999999 = 无限配额\n", adminCount)
	fmt.Printf("    VIP用户: %d 人 × 100000 = %d 配额\n", vipCount, vipCount*100000)
	fmt.Printf("    普通用户: %d 人 × 10000 = %d 配额\n", normalCount, normalCount*10000)
}

// Clean 清空AI配额数据
func (s *AIQuotaSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("ai_quotas").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 ai_quotas 集合失败: %w", err)
	}

	fmt.Println("  已清空 ai_quotas 集合")
	return nil
}
