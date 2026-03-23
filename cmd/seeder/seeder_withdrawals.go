// Package main 提供提现申请数据填充功能
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WithdrawalSeeder 提现申请数据填充器
type WithdrawalSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewWithdrawalSeeder 创建提现申请数据填充器
func NewWithdrawalSeeder(db *utils.Database, cfg *config.Config) *WithdrawalSeeder {
	collection := db.Collection("withdraw_requests")
	return &WithdrawalSeeder{
		db:       db,
		config:   cfg,
		inserter: utils.NewBulkInserter(collection, cfg.BatchSize),
	}
}

// SeedWithdrawals 填充提现申请数据
func (s *WithdrawalSeeder) SeedWithdrawals() error {
	ctx := context.Background()

	// 获取作者用户
	authors, err := s.getAuthorIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取作者列表失败: %w", err)
	}

	if len(authors) == 0 {
		fmt.Println("  没有找到作者用户，请先运行 users 命令创建作者")
		return nil
	}

	// 创建提现申请
	if err := s.seedWithdrawalsData(ctx, authors); err != nil {
		return err
	}

	return nil
}

// getAuthorIDs 获取作者ID列表
func (s *WithdrawalSeeder) getAuthorIDs(ctx context.Context) ([]string, error) {
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{"roles": bson.M{"$in": []string{"author"}}})
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

	authorIDs := make([]string, len(users))
	for i, u := range users {
		authorIDs[i] = u.ID
	}
	return authorIDs, nil
}

// seedWithdrawalsData 创建提现申请数据
func (s *WithdrawalSeeder) seedWithdrawalsData(ctx context.Context, authorIDs []string) error {
	collection := s.db.Collection("withdraw_requests")

	var withdrawals []interface{}
	now := time.Now()

	// 提现状态
	statuses := []string{"pending", "approved", "rejected", "processed"}
	accountTypes := []string{"alipay", "wechat", "bank"}

	// 为 50% 的作者创建提现申请
	authorCount := len(authorIDs) / 2
	if authorCount < 1 && len(authorIDs) > 0 {
		authorCount = 1
	}

	for i := 0; i < authorCount; i++ {
		authorID := authorIDs[i]

		// 每个作者 1-5 条提现申请
		withdrawalCount := 1 + rand.Intn(5)

		for j := 0; j < withdrawalCount; j++ {
			amount := 100.0 + rand.Float64()*900 // 100-1000元
			fee := amount * 0.01                  // 1% 手续费
			actualAmount := amount - fee

			status := statuses[rand.Intn(len(statuses))]
			accountType := accountTypes[rand.Intn(len(accountTypes))]

			var account string
			var accountName string
			switch accountType {
			case "alipay":
				account = fmt.Sprintf("alipay_%d@qq.com", rand.Intn(100000))
				accountName = "支付宝账户"
			case "wechat":
				account = fmt.Sprintf("wx_%d", rand.Intn(100000))
				accountName = "微信账户"
			case "bank":
				account = fmt.Sprintf("6222%012d", rand.Int63n(1000000000000))
				accountName = "银行卡账户"
			}

			withdrawal := bson.M{
				"_id":               primitive.NewObjectID(),
				"user_id":           authorID,
				"amount_cents":      int64(amount * 100),
				"fee_cents":         int64(fee * 100),
				"actual_amount_cents": int64(actualAmount * 100),
				"account":           account,
				"account_type":      accountType,
				"account_name":      accountName,
				"status":            status,
				"order_no":          fmt.Sprintf("WD%s%04d", now.Format("20060102150405"), rand.Intn(10000)),
				"created_at":        now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
				"updated_at":        now,
			}

			// 根据状态设置额外字段
			if status == "approved" || status == "processed" {
				withdrawal["reviewed_at"] = now.Add(-time.Duration(rand.Intn(48)) * time.Hour)
				withdrawal["reviewed_by"] = "admin"
			}
			if status == "rejected" {
				withdrawal["reviewed_at"] = now.Add(-time.Duration(rand.Intn(48)) * time.Hour)
				withdrawal["reviewed_by"] = "admin"
				withdrawal["reject_reason"] = "账户信息有误，请核实后重新提交"
			}
			if status == "processed" {
				withdrawal["processed_at"] = now.Add(-time.Duration(rand.Intn(24)) * time.Hour)
				withdrawal["transaction_id"] = primitive.NewObjectID().Hex()
			}

			withdrawals = append(withdrawals, withdrawal)
		}
	}

	// 批量插入
	if len(withdrawals) > 0 {
		if _, err := collection.InsertMany(ctx, withdrawals); err != nil {
			return fmt.Errorf("插入提现申请失败: %w", err)
		}
		fmt.Printf("  创建了 %d 条提现申请\n", len(withdrawals))
	}

	return nil
}

// Clean 清空提现申请数据
func (s *WithdrawalSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("withdraw_requests").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 withdraw_requests 集合失败: %w", err)
	}

	fmt.Println("  已清空 withdraw_requests 集合")
	return nil
}
