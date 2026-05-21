// Package main 提供AI配额激活功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"
	aiModels "Qingyu_backend/models/ai"
	aiMongoRepo "Qingyu_backend/repository/mongodb/ai"
	aiService "Qingyu_backend/service/ai"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type aiQuotaUsageSample struct {
	Username     string
	WorkflowType string
	Model        string
	Tokens       int
	RequestID    string
}

var defaultAIQuotaUsageSamples = []aiQuotaUsageSample{
	{Username: "testadmin001", WorkflowType: "admin_audit", Model: "local-seed-admin", Tokens: 320, RequestID: "local-quota-seed-admin-audit"},
	{Username: "testauthor001", WorkflowType: "story_write", Model: "local-seed-writer", Tokens: 860, RequestID: "local-quota-seed-story-write"},
	{Username: "testauthor001", WorkflowType: "chapter_expand", Model: "local-seed-writer", Tokens: 540, RequestID: "local-quota-seed-chapter-expand"},
	{Username: "testuser001", WorkflowType: "reader_chat", Model: "local-seed-reader", Tokens: 210, RequestID: "local-quota-seed-reader-chat"},
	{Username: "testuser002", WorkflowType: "reader_rewrite", Model: "local-seed-reader", Tokens: 430, RequestID: "local-quota-seed-reader-rewrite"},
}

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

type quotaSeedUser struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Roles    []string           `bson:"roles"`
}

// SeedQuotaUsageSamples 为本地 quota 对账补 backend 侧交易样本。
func (s *AIQuotaSeeder) SeedQuotaUsageSamples() error {
	ctx := context.Background()

	quotaRepo := aiMongoRepo.NewMongoQuotaRepository(s.db.Database)
	quotaService := aiService.NewQuotaService(quotaRepo)
	transactionCollection := s.db.Collection(aiModels.QuotaTransaction{}.CollectionName())

	targetUsernames := make([]string, 0, len(defaultAIQuotaUsageSamples))
	for _, sample := range defaultAIQuotaUsageSamples {
		targetUsernames = append(targetUsernames, sample.Username)
	}

	cursor, err := s.db.Collection("users").Find(ctx, bson.M{
		"username": bson.M{"$in": targetUsernames},
	})
	if err != nil {
		return fmt.Errorf("查询 quota usage 样本用户失败: %w", err)
	}
	defer cursor.Close(ctx)

	var users []quotaSeedUser
	if err := cursor.All(ctx, &users); err != nil {
		return fmt.Errorf("解析 quota usage 样本用户失败: %w", err)
	}

	userByUsername := make(map[string]quotaSeedUser, len(users))
	for _, user := range users {
		userByUsername[user.Username] = user
	}

	seededCount := 0
	skippedCount := 0

	for _, sample := range defaultAIQuotaUsageSamples {
		user, ok := userByUsername[sample.Username]
		if !ok {
			return fmt.Errorf("未找到 quota usage 样本用户: %s，请先执行 users/baseline seeder", sample.Username)
		}

		existingCount, err := transactionCollection.CountDocuments(ctx, bson.M{
			"request_id": sample.RequestID,
		})
		if err != nil {
			return fmt.Errorf("检查 quota usage 样本是否已存在失败: %w", err)
		}
		if existingCount > 0 {
			skippedCount++
			continue
		}

		userID := user.ID.Hex()
		role := seedPrimaryRole(user.Roles)
		if role == "" {
			role = "reader"
		}

		if err := quotaService.InitializeUserQuota(ctx, userID, role, "normal"); err != nil {
			return fmt.Errorf("初始化 quota 失败 user=%s: %w", sample.Username, err)
		}
		if err := ensureQuotaCapacity(ctx, quotaService, userID, sample.Tokens); err != nil {
			return fmt.Errorf("补足 quota 容量失败 user=%s: %w", sample.Username, err)
		}

		if err := quotaService.ConsumeQuota(
			ctx,
			userID,
			sample.Tokens,
			sample.WorkflowType,
			sample.Model,
			sample.RequestID,
		); err != nil {
			return fmt.Errorf("写入 quota usage 样本失败 user=%s workflow=%s: %w", sample.Username, sample.WorkflowType, err)
		}

		fmt.Printf(
			"  已写入 backend quota usage: username=%s user_id=%s workflow=%s tokens=%d request_id=%s\n",
			sample.Username,
			userID,
			sample.WorkflowType,
			sample.Tokens,
			sample.RequestID,
		)
		seededCount++
	}

	fmt.Printf("  新增 backend quota usage 样本: %d 条\n", seededCount)
	fmt.Printf("  已存在跳过: %d 条\n", skippedCount)

	return nil
}

func seedPrimaryRole(roles []string) string {
	if len(roles) == 0 {
		return ""
	}

	for _, preferred := range []string{"admin", "author", "reader"} {
		for _, role := range roles {
			if role == preferred {
				return role
			}
		}
	}

	return roles[0]
}

func ensureQuotaCapacity(ctx context.Context, quotaService *aiService.QuotaService, userID string, requiredAmount int) error {
	quotaInfo, err := quotaService.GetQuotaInfo(ctx, userID)
	if err != nil {
		if err == aiModels.ErrQuotaNotFound {
			return quotaService.UpdateUserQuota(ctx, userID, aiModels.QuotaTypeDaily, requiredAmount+1000)
		}
		return err
	}

	if quotaInfo.RemainingQuota >= requiredAmount {
		return nil
	}

	targetTotal := quotaInfo.UsedQuota + requiredAmount + 1000
	if targetTotal <= quotaInfo.TotalQuota {
		targetTotal = quotaInfo.TotalQuota + requiredAmount + 1000
	}

	return quotaService.UpdateUserQuota(ctx, userID, aiModels.QuotaTypeDaily, targetTotal)
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
