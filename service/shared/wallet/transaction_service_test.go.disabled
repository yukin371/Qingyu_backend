package wallet

import (
	"context"
	"testing"
)

// TestRecharge 测试充值
func TestRecharge(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包
	wallet, _ := walletService.CreateWallet(ctx, "user123")

	// 充值
	tx, err := txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")
	if err != nil {
		t.Fatalf("充值失败: %v", err)
	}

	if tx.Amount != 100.0 {
		t.Errorf("充值金额错误: %f", tx.Amount)
	}

	if tx.Type != "recharge" {
		t.Errorf("交易类型错误: %s", tx.Type)
	}

	// 验证余额
	balance, _ := walletService.GetBalance(ctx, wallet.ID)
	if balance != 100.0 {
		t.Errorf("余额应为100，实际: %f", balance)
	}

	t.Logf("充值测试通过: amount=%.2f, balance=%.2f", tx.Amount, balance)
}

// TestRecharge_InvalidAmount 测试充值无效金额
func TestRecharge_InvalidAmount(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包
	wallet, _ := walletService.CreateWallet(ctx, "user123")

	// 充值负数
	_, err := txService.Recharge(ctx, wallet.ID, -100.0, "alipay", "order_001")
	if err == nil {
		t.Error("充值负数应该失败")
	}

	// 充值0
	_, err = txService.Recharge(ctx, wallet.ID, 0, "alipay", "order_002")
	if err == nil {
		t.Error("充值0应该失败")
	}

	t.Logf("充值无效金额测试通过")
}

// TestConsume 测试消费
func TestConsume(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 消费
	tx, err := txService.Consume(ctx, wallet.ID, 30.0, "购买VIP")
	if err != nil {
		t.Fatalf("消费失败: %v", err)
	}

	if tx.Amount != -30.0 {
		t.Errorf("消费金额错误: %f（应为负数）", tx.Amount)
	}

	if tx.Type != "consume" {
		t.Errorf("交易类型错误: %s", tx.Type)
	}

	// 验证余额
	balance, _ := walletService.GetBalance(ctx, wallet.ID)
	if balance != 70.0 {
		t.Errorf("余额应为70，实际: %f", balance)
	}

	t.Logf("消费测试通过: amount=%.2f, balance=%.2f", tx.Amount, balance)
}

// TestConsume_InsufficientBalance 测试余额不足
func TestConsume_InsufficientBalance(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包（余额为0）
	wallet, _ := walletService.CreateWallet(ctx, "user123")

	// 尝试消费
	_, err := txService.Consume(ctx, wallet.ID, 100.0, "购买VIP")
	if err == nil {
		t.Error("余额不足应该失败")
	}

	t.Logf("余额不足测试通过: %v", err)
}

// TestTransfer 测试转账
func TestTransfer(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建两个钱包
	walletA, _ := walletService.CreateWallet(ctx, "userA")
	walletB, _ := walletService.CreateWallet(ctx, "userB")

	// A充值100
	txService.Recharge(ctx, walletA.ID, 100.0, "alipay", "order_001")

	// A转账50给B
	err := txService.Transfer(ctx, walletA.ID, walletB.ID, 50.0, "转账")
	if err != nil {
		t.Fatalf("转账失败: %v", err)
	}

	// 验证A的余额
	balanceA, _ := walletService.GetBalance(ctx, walletA.ID)
	if balanceA != 50.0 {
		t.Errorf("A余额应为50，实际: %f", balanceA)
	}

	// 验证B的余额
	balanceB, _ := walletService.GetBalance(ctx, walletB.ID)
	if balanceB != 50.0 {
		t.Errorf("B余额应为50，实际: %f", balanceB)
	}

	t.Logf("转账测试通过: A余额=%.2f, B余额=%.2f", balanceA, balanceB)
}

// TestTransfer_InsufficientBalance 测试转账余额不足
func TestTransfer_InsufficientBalance(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建两个钱包
	walletA, _ := walletService.CreateWallet(ctx, "userA")
	walletB, _ := walletService.CreateWallet(ctx, "userB")

	// A余额为0，尝试转账
	err := txService.Transfer(ctx, walletA.ID, walletB.ID, 100.0, "转账")
	if err == nil {
		t.Error("余额不足转账应该失败")
	}

	t.Logf("转账余额不足测试通过: %v", err)
}

// TestGetTransaction 测试获取交易记录
func TestGetTransaction(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	tx, _ := txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 获取交易记录
	retrieved, err := txService.GetTransaction(ctx, tx.ID)
	if err != nil {
		t.Fatalf("获取交易记录失败: %v", err)
	}

	if retrieved.ID != tx.ID {
		t.Errorf("交易ID不匹配")
	}

	if retrieved.Amount != 100.0 {
		t.Errorf("交易金额错误: %f", retrieved.Amount)
	}

	t.Logf("获取交易记录测试通过: %+v", retrieved)
}

// TestListTransactions 测试列出交易记录
func TestListTransactions(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包
	wallet, _ := walletService.CreateWallet(ctx, "user123")

	// 执行多次交易
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")
	txService.Consume(ctx, wallet.ID, 30.0, "购买VIP")
	txService.Consume(ctx, wallet.ID, 20.0, "购买道具")

	// 列出交易记录
	transactions, err := txService.ListTransactions(ctx, wallet.ID, 10, 0)
	if err != nil {
		t.Fatalf("列出交易记录失败: %v", err)
	}

	if len(transactions) != 3 {
		t.Errorf("交易记录数量错误: %d", len(transactions))
	}

	t.Logf("列出交易记录测试通过，共%d条记录", len(transactions))
	for i, tx := range transactions {
		t.Logf("  [%d] %s: %.2f", i+1, tx.Type, tx.Amount)
	}
}

// TestMultipleTransactions 测试多次交易余额正确性
func TestMultipleTransactions(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)

	ctx := context.Background()

	// 创建钱包
	wallet, _ := walletService.CreateWallet(ctx, "user123")

	// 充值200
	txService.Recharge(ctx, wallet.ID, 200.0, "alipay", "order_001")

	// 消费50
	txService.Consume(ctx, wallet.ID, 50.0, "购买VIP")

	// 充值100
	txService.Recharge(ctx, wallet.ID, 100.0, "wechat", "order_002")

	// 消费80
	txService.Consume(ctx, wallet.ID, 80.0, "购买道具")

	// 最终余额应为: 200 - 50 + 100 - 80 = 170
	balance, _ := walletService.GetBalance(ctx, wallet.ID)
	expectedBalance := 170.0

	if balance != expectedBalance {
		t.Errorf("最终余额错误: 期望%.2f，实际%.2f", expectedBalance, balance)
	}

	t.Logf("多次交易测试通过，最终余额: %.2f", balance)
}
