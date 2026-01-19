package seeds

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Transaction 交易记录
type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	Type        string             `bson:"type"` // recharge, consume, transfer, withdraw, income
	Amount      float64            `bson:"amount"`
	Balance     float64            `bson:"balance"` // 交易后余额
	Description string             `bson:"description"`
	Status      string             `bson:"status"` // completed, pending, failed, cancelled
	CreatedAt   time.Time          `bson:"created_at"`
}

// WalletData 钱包数据
type WalletData struct {
	UserID        string    `bson:"user_id"`
	Balance       float64   `bson:"balance"`
	Frozen        bool      `bson:"frozen"`
	FrozenAmount  float64   `bson:"frozen_amount"`
	TotalRecharge float64   `bson:"total_recharge"`
	TotalConsume  float64   `bson:"total_consume"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

// SeedWallets 钱包种子数据
func SeedWallets(ctx context.Context, db *mongo.Database) error {
	walletCollection := db.Collection("wallets")
	transactionCollection := db.Collection("transactions")

	fmt.Println("========================================")
	fmt.Println("开始创建钱包和交易测试数据...")
	fmt.Println("========================================")

	// 获取用户ID列表
	userCollection := db.Collection("users")
	cursor, err := userCollection.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID    string `bson:"_id"`
		Role  string `bson:"role"`
		Email string `bson:"email"`
	}
	if err = cursor.All(ctx, &users); err != nil {
		return fmt.Errorf("解析用户列表失败: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("警告：没有找到活跃用户，跳过钱包创建")
		return nil
	}

	successCount := 0
	skipCount := 0

	for _, user := range users {
		fmt.Printf("处理用户: %s (%s)...", user.Email, user.Role)

		// 检查钱包是否已存在
		var existingWallet WalletData
		err := walletCollection.FindOne(ctx, bson.M{"user_id": user.ID}).Decode(&existingWallet)
		if err == nil {
			fmt.Println(" [钱包已存在，跳过]")
			skipCount++
			continue
		}

		// 根据角色设置初始余额
		var initialBalance float64
		var transactions []Transaction

		switch user.Role {
		case "admin":
			initialBalance = 10000.0
			transactions = generateAdminTransactions(user.ID)
		case "author":
			initialBalance = 5000.0 + rand.Float64()*5000
			transactions = generateAuthorTransactions(user.ID, initialBalance)
		case "reader":
			initialBalance = 100.0 + rand.Float64()*900
			transactions = generateReaderTransactions(user.ID, initialBalance)
		default:
			initialBalance = 50.0
			transactions = generateReaderTransactions(user.ID, initialBalance)
		}

		// 创建钱包
		now := time.Now()
		wallet := WalletData{
			UserID:        user.ID,
			Balance:       initialBalance,
			Frozen:        false,
			FrozenAmount:  0,
			TotalRecharge: calculateTotalRecharge(transactions),
			TotalConsume:  calculateTotalConsume(transactions),
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		_, err = walletCollection.InsertOne(ctx, wallet)
		if err != nil {
			fmt.Printf(" [创建钱包失败: %v]\n", err)
			continue
		}

		// 创建交易记录
		if len(transactions) > 0 {
			transDocs := make([]interface{}, len(transactions))
			for i, trans := range transactions {
				transDocs[i] = trans
			}
			_, err = transactionCollection.InsertMany(ctx, transDocs)
			if err != nil {
				fmt.Printf(" [创建交易记录失败: %v]\n", err)
				continue
			}
		}

		fmt.Printf(" [成功，余额: %.2f，交易: %d条]\n", initialBalance, len(transactions))
		successCount++
	}

	fmt.Println("========================================")
	fmt.Println("钱包创建完成")
	fmt.Println("========================================")
	fmt.Printf("成功创建: %d 个\n", successCount)
	fmt.Printf("已存在跳过: %d 个\n", skipCount)
	fmt.Printf("总计: %d 个\n", len(users))
	fmt.Println()

	return nil
}

func generateAdminTransactions(userID string) []Transaction {
	now := time.Now()
	transactions := []Transaction{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "recharge",
			Amount:      10000.0,
			Balance:     10000.0,
			Description: "系统充值",
			Status:      "completed",
			CreatedAt:   now.Add(-30 * 24 * time.Hour),
		},
	}
	return transactions
}

func generateAuthorTransactions(userID string, balance float64) []Transaction {
	now := time.Now()
	transactions := []Transaction{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "recharge",
			Amount:      5000.0 + rand.Float64()*5000,
			Balance:     balance,
			Description: "初始充值",
			Status:      "completed",
			CreatedAt:   now.Add(-60 * 24 * time.Hour),
		},
	}

	// 收入记录（来自作品销售）
	incomeCount := 3 + rand.Intn(10)
	for i := 0; i < incomeCount; i++ {
		amount := 10.0 + rand.Float64()*100
		transactions = append(transactions, Transaction{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "income",
			Amount:      amount,
			Balance:     balance,
			Description: "作品销售收入",
			Status:      "completed",
			CreatedAt:   now.Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
		})
	}

	// 消费记录
	consumeCount := 2 + rand.Intn(5)
	for i := 0; i < consumeCount; i++ {
		amount := 5.0 + rand.Float64()*50
		transactions = append(transactions, Transaction{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "consume",
			Amount:      -amount,
			Balance:     balance,
			Description: "购买章节",
			Status:      "completed",
			CreatedAt:   now.Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
		})
	}

	return transactions
}

func generateReaderTransactions(userID string, balance float64) []Transaction {
	now := time.Now()
	transactions := []Transaction{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "recharge",
			Amount:      100.0 + rand.Float64()*900,
			Balance:     balance,
			Description: "账户充值",
			Status:      "completed",
			CreatedAt:   now.Add(-time.Duration(10+rand.Intn(60)) * 24 * time.Hour),
		},
	}

	// 购买记录
	purchaseCount := 5 + rand.Intn(20)
	for i := 0; i < purchaseCount; i++ {
		amount := 0.05 + rand.Float64()*5
		transactions = append(transactions, Transaction{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Type:        "consume",
			Amount:      -amount,
			Balance:     balance,
			Description: "购买章节",
			Status:      "completed",
			CreatedAt:   now.Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
		})
	}

	return transactions
}

func calculateTotalRecharge(transactions []Transaction) float64 {
	total := 0.0
	for _, t := range transactions {
		if t.Type == "recharge" || t.Type == "income" {
			total += t.Amount
		}
	}
	return total
}

func calculateTotalConsume(transactions []Transaction) float64 {
	total := 0.0
	for _, t := range transactions {
		if t.Type == "consume" {
			total += t.Amount
		}
	}
	return total
}
