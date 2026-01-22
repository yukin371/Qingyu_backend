package wallet

import (
	"context"
	"fmt"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"
	mongoFinance "Qingyu_backend/repository/mongodb/finance"

	"go.mongodb.org/mongo-driver/mongo"
)

// ExampleWalletServiceUsage 展示如何使用Wallet服务
// 这是一个示例文件，演示如何创建和使用统一的Wallet服务
// 注意：所有金额单位都是"分"（int64），100分 = 1元
func ExampleWalletServiceUsage(db *mongo.Database) {
	ctx := context.Background()

	// 1. 创建Repository
	walletRepo := mongoFinance.NewWalletRepository(db)

	// 2. 创建Wallet服务（使用统一的服务）
	walletService := NewUnifiedWalletService(walletRepo)

	// 3. 使用示例
	userID := "user123"

	// 3.1 创建钱包
	wallet, err := walletService.CreateWallet(ctx, userID)
	if err != nil {
		fmt.Printf("创建钱包失败: %v\n", err)
		return
	}
	fmt.Printf("创建钱包成功: %+v\n", wallet)

	// 3.2 充值（10000分 = 100元）
	transaction, err := walletService.Recharge(ctx, userID, 10000, "alipay")
	if err != nil {
		fmt.Printf("充值失败: %v\n", err)
		return
	}
	fmt.Printf("充值成功: %+v\n", transaction)

	// 3.3 查询余额（返回分）
	balance, err := walletService.GetBalance(ctx, userID)
	if err != nil {
		fmt.Printf("查询余额失败: %v\n", err)
		return
	}
	fmt.Printf("当前余额: %.2f元\n", float64(balance)/100.0)

	// 3.4 消费（3000分 = 30元）
	consumeTx, err := walletService.Consume(ctx, userID, 3000, "购买VIP会员")
	if err != nil {
		fmt.Printf("消费失败: %v\n", err)
		return
	}
	fmt.Printf("消费成功: %+v\n", consumeTx)

	// 3.5 查询交易记录
	txList, err := walletService.ListTransactions(ctx, userID, &ListTransactionsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("查询交易记录失败: %v\n", err)
		return
	}
	fmt.Printf("交易记录数量: %d\n", len(txList))

	// 3.6 申请提现（5000分 = 50元）
	withdrawReq, err := walletService.RequestWithdraw(ctx, userID, 5000, "alipay@example.com")
	if err != nil {
		fmt.Printf("申请提现失败: %v\n", err)
		return
	}
	fmt.Printf("提现申请成功: %+v\n", withdrawReq)

	fmt.Println("Wallet服务使用示例完成")
}

// ExampleDirectServiceUsage 展示如何直接使用分离的服务（不推荐）
// 推荐使用UnifiedWalletService
func ExampleDirectServiceUsage(walletRepo sharedRepo.WalletRepository) {
	// 这是旧的使用方式，不推荐
	// 这里展示是为了向后兼容
	_ = NewWalletService(walletRepo)      // 钱包管理
	_ = NewTransactionService(walletRepo) // 交易服务
	_ = NewWithdrawService(walletRepo)    // 提现服务

	// 新代码应该使用UnifiedWalletService
	// walletService := NewUnifiedWalletService(walletRepo)
}
