// Package main 提供钱包数据填充功能
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

// WalletSeeder 钱包数据填充器
type WalletSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewWalletSeeder 创建钱包数据填充器
func NewWalletSeeder(db *utils.Database, cfg *config.Config) *WalletSeeder {
	return &WalletSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedWallets 填充钱包数据
func (s *WalletSeeder) SeedWallets() error {
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

	// 创建钱包和交易
	walletCollection := s.db.Collection("wallets")
	transactionCollection := s.db.Collection("transactions")

	successCount := 0
	skipCount := 0

	for _, user := range users {
		// 检查钱包是否已存在
		var existingWallet map[string]interface{}
		err := walletCollection.FindOne(ctx, bson.M{"user_id": user.ID}).Decode(&existingWallet)
		if err == nil {
			skipCount++
			continue
		}

		// 根据角色设置初始余额
		var initialBalance float64
		var transactions []interface{}

		switch user.Role {
		case "admin":
			initialBalance = 10000.0
			transactions = s.generateAdminTransactions(user.ID)
		case "vip":
			initialBalance = 5000.0 + rand.Float64()*5000
			transactions = s.generateAuthorTransactions(user.ID, initialBalance)
		case "author":
			initialBalance = 3000.0 + rand.Float64()*3000
			transactions = s.generateAuthorTransactions(user.ID, initialBalance)
		default: // user/reader
			initialBalance = 100.0 + rand.Float64()*900
			transactions = s.generateReaderTransactions(user.ID, initialBalance)
		}

		// 创建钱包
		now := time.Now()
		wallet := bson.M{
			"_id":                primitive.NewObjectID(),
			"user_id":            user.ID,
			"balance":            initialBalance,
			"frozen":             false,
			"frozen_amount":      0,
			"total_recharge":     initialBalance * 0.5,
			"total_consume":      initialBalance * 0.1,
			"created_at":         now,
			"updated_at":         now,
		}

		_, err = walletCollection.InsertOne(ctx, wallet)
		if err != nil {
			return fmt.Errorf("创建钱包失败: %w", err)
		}

		// 创建交易记录
		if len(transactions) > 0 {
			_, err = transactionCollection.InsertMany(ctx, transactions)
			if err != nil {
				return fmt.Errorf("创建交易记录失败: %w", err)
			}
		}

		successCount++
	}

	fmt.Printf("  创建成功: %d 个钱包\n", successCount)
	fmt.Printf("  已存在跳过: %d 个钱包\n", skipCount)

	return nil
}

// generateAdminTransactions 生成管理员交易记录
func (s *WalletSeeder) generateAdminTransactions(userID string) []interface{} {
	var transactions []interface{}

	// 充值记录
	for i := 0; i < 5; i++ {
		amount := 1000.0 + rand.Float64()*1000
		transactions = append(transactions, bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      userID,
			"type":         "recharge",
			"amount":       amount,
			"balance":      10000.0 - float64(5-i-1)*2000,
			"description":  "系统充值",
			"status":       "completed",
			"created_at":   time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
		})
	}

	return transactions
}

// generateAuthorTransactions 生成作者交易记录
func (s *WalletSeeder) generateAuthorTransactions(userID string, balance float64) []interface{} {
	var transactions []interface{}

	// 充值记录
	for i := 0; i < 3; i++ {
		amount := 500.0 + rand.Float64()*500
		transactions = append(transactions, bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      userID,
			"type":         "recharge",
			"amount":       amount,
			"balance":      balance - float64(3-i-1)*1000,
			"description":  "充值",
			"status":       "completed",
			"created_at":   time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
		})
	}

	// 收入记录（稿费）
	for i := 0; i < 2; i++ {
		amount := 100.0 + rand.Float64()*200
		transactions = append(transactions, bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      userID,
			"type":         "income",
			"amount":       amount,
			"balance":      balance + float64(i)*200,
			"description":  "稿费收入",
			"status":       "completed",
			"created_at":   time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
		})
	}

	return transactions
}

// generateReaderTransactions 生成读者交易记录
func (s *WalletSeeder) generateReaderTransactions(userID string, balance float64) []interface{} {
	var transactions []interface{}

	// 充值记录
	rechargeCount := 1 + rand.Intn(3)
	for i := 0; i < rechargeCount; i++ {
		amount := 50.0 + rand.Float64()*100
		transactions = append(transactions, bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      userID,
			"type":         "recharge",
			"amount":       amount,
			"balance":      balance - float64(rechargeCount-i-1)*50,
			"description":  "充值",
			"status":       "completed",
			"created_at":   time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
		})
	}

	// 消费记录
	consumeCount := rand.Intn(5)
	for i := 0; i < consumeCount; i++ {
		amount := 5.0 + rand.Float64()*20
		transactions = append(transactions, bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      userID,
			"type":         "consume",
			"amount":       amount,
			"balance":      balance - float64(i)*10,
			"description":  "购买章节",
			"status":       "completed",
			"created_at":   time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
		})
	}

	return transactions
}

// Clean 清空钱包数据
func (s *WalletSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("wallets").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 wallets 集合失败: %w", err)
	}

	_, err = s.db.Collection("transactions").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 transactions 集合失败: %w", err)
	}

	fmt.Println("  已清空 wallets 和 transactions 集合")
	return nil
}
