package wallet

import (
	"context"
	"testing"
)

// ============ 测试用例 ============

// TestUnifiedWalletService_CreateWallet 测试创建钱包
func TestUnifiedWalletService_CreateWallet(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_create_test"
	wallet, err := service.CreateWallet(ctx, userID)

	if err != nil {
		t.Fatalf("创建钱包失败: %v", err)
	}

	if wallet.UserID != userID {
		t.Errorf("用户ID不匹配: got %s, want %s", wallet.UserID, userID)
	}

	if wallet.Balance != 0 {
		t.Errorf("初始余额应为0: got %.2f", wallet.Balance)
	}

	if wallet.Frozen {
		t.Error("钱包不应该被冻结")
	}
}

// TestUnifiedWalletService_Recharge 测试充值
func TestUnifiedWalletService_Recharge(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_recharge_test"

	// 先创建钱包
	_, err := service.CreateWallet(ctx, userID)
	if err != nil {
		t.Fatalf("创建钱包失败: %v", err)
	}

	// 充值
	amount := 100.0
	tx, err := service.Recharge(ctx, userID, amount, "alipay")
	if err != nil {
		t.Fatalf("充值失败: %v", err)
	}

	if tx.Type != "recharge" {
		t.Errorf("交易类型错误: got %s, want recharge", tx.Type)
	}

	if tx.Amount != amount {
		t.Errorf("充值金额错误: got %.2f, want %.2f", tx.Amount, amount)
	}

	// 检查余额
	balance, err := service.GetBalance(ctx, userID)
	if err != nil {
		t.Fatalf("获取余额失败: %v", err)
	}

	if balance != amount {
		t.Errorf("余额错误: got %.2f, want %.2f", balance, amount)
	}
}

// TestUnifiedWalletService_Consume 测试消费
func TestUnifiedWalletService_Consume(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_consume_test"

	// 创建钱包并充值
	_, _ = service.CreateWallet(ctx, userID)
	_, _ = service.Recharge(ctx, userID, 100.0, "alipay")

	// 消费
	consumeAmount := 30.0
	tx, err := service.Consume(ctx, userID, consumeAmount, "购买VIP")
	if err != nil {
		t.Fatalf("消费失败: %v", err)
	}

	if tx.Type != "consume" {
		t.Errorf("交易类型错误: got %s, want consume", tx.Type)
	}

	// 检查余额
	balance, _ := service.GetBalance(ctx, userID)
	expectedBalance := 100.0 - consumeAmount

	if balance != expectedBalance {
		t.Errorf("余额错误: got %.2f, want %.2f", balance, expectedBalance)
	}
}

// TestUnifiedWalletService_Transfer 测试转账
func TestUnifiedWalletService_Transfer(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	fromUserID := "user_transfer_from"
	toUserID := "user_transfer_to"

	// 创建两个钱包
	_, _ = service.CreateWallet(ctx, fromUserID)
	_, _ = service.CreateWallet(ctx, toUserID)

	// 给源用户充值
	_, _ = service.Recharge(ctx, fromUserID, 100.0, "alipay")

	// 转账
	transferAmount := 50.0
	tx, err := service.Transfer(ctx, fromUserID, toUserID, transferAmount, "转账测试")
	if err != nil {
		t.Fatalf("转账失败: %v", err)
	}

	if tx == nil {
		t.Fatal("转账交易记录为空")
	}

	// 检查两个用户的余额
	fromBalance, _ := service.GetBalance(ctx, fromUserID)
	toBalance, _ := service.GetBalance(ctx, toUserID)

	expectedFromBalance := 100.0 - transferAmount
	expectedToBalance := transferAmount

	if fromBalance != expectedFromBalance {
		t.Errorf("转出用户余额错误: got %.2f, want %.2f", fromBalance, expectedFromBalance)
	}

	if toBalance != expectedToBalance {
		t.Errorf("转入用户余额错误: got %.2f, want %.2f", toBalance, expectedToBalance)
	}
}

// TestUnifiedWalletService_FreezeAndUnfreeze 测试冻结和解冻
func TestUnifiedWalletService_FreezeAndUnfreeze(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_freeze_test"

	// 创建钱包
	_, _ = service.CreateWallet(ctx, userID)

	// 冻结钱包
	err := service.FreezeWallet(ctx, userID)
	if err != nil {
		t.Fatalf("冻结钱包失败: %v", err)
	}

	// 检查钱包是否已冻结
	wallet, _ := service.GetWallet(ctx, userID)
	if !wallet.Frozen {
		t.Error("钱包应该被冻结")
	}

	// 解冻钱包
	err = service.UnfreezeWallet(ctx, userID)
	if err != nil {
		t.Fatalf("解冻钱包失败: %v", err)
	}

	// 检查钱包是否已解冻
	wallet, _ = service.GetWallet(ctx, userID)
	if wallet.Frozen {
		t.Error("钱包应该被解冻")
	}
}

// TestUnifiedWalletService_RequestWithdraw 测试申请提现
func TestUnifiedWalletService_RequestWithdraw(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_withdraw_test"

	// 创建钱包并充值
	_, _ = service.CreateWallet(ctx, userID)
	_, _ = service.Recharge(ctx, userID, 100.0, "alipay")

	// 申请提现
	withdrawAmount := 50.0
	req, err := service.RequestWithdraw(ctx, userID, withdrawAmount, "alipay@example.com")
	if err != nil {
		t.Fatalf("申请提现失败: %v", err)
	}

	if req.UserID != userID {
		t.Errorf("用户ID错误: got %s, want %s", req.UserID, userID)
	}

	if req.Amount != withdrawAmount {
		t.Errorf("提现金额错误: got %.2f, want %.2f", req.Amount, withdrawAmount)
	}

	if req.Status != "pending" {
		t.Errorf("提现状态错误: got %s, want pending", req.Status)
	}

	// 检查余额（提现申请后余额应扣除）
	balance, _ := service.GetBalance(ctx, userID)
	expectedBalance := 100.0 - withdrawAmount

	if balance != expectedBalance {
		t.Errorf("余额错误: got %.2f, want %.2f", balance, expectedBalance)
	}
}

// TestUnifiedWalletService_WithdrawFlow 测试完整提现流程
func TestUnifiedWalletService_WithdrawFlow(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_withdrawflow_test"
	adminID := "admin001"

	// 创建钱包并充值
	_, _ = service.CreateWallet(ctx, userID)
	_, _ = service.Recharge(ctx, userID, 100.0, "alipay")

	// 1. 申请提现
	req, err := service.RequestWithdraw(ctx, userID, 50.0, "alipay@example.com")
	if err != nil {
		t.Fatalf("申请提现失败: %v", err)
	}

	// 2. 管理员批准
	err = service.ApproveWithdraw(ctx, req.ID, adminID)
	if err != nil {
		t.Fatalf("批准提现失败: %v", err)
	}

	// 检查状态
	approvedReq, _ := service.GetWithdrawRequest(ctx, req.ID)
	if approvedReq.Status != "approved" {
		t.Errorf("提现状态错误: got %s, want approved", approvedReq.Status)
	}

	// 3. 处理提现
	err = service.ProcessWithdraw(ctx, req.ID)
	if err != nil {
		t.Fatalf("处理提现失败: %v", err)
	}

	// 检查最终状态
	processedReq, _ := service.GetWithdrawRequest(ctx, req.ID)
	if processedReq.Status != "processed" {
		t.Errorf("提现状态错误: got %s, want processed", processedReq.Status)
	}
}

// TestUnifiedWalletService_ListTransactions 测试查询交易列表
func TestUnifiedWalletService_ListTransactions(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	userID := "user_listtx_test"

	// 创建钱包并进行多次交易
	_, _ = service.CreateWallet(ctx, userID)
	_, _ = service.Recharge(ctx, userID, 100.0, "alipay")
	_, _ = service.Consume(ctx, userID, 30.0, "购买VIP")
	_, _ = service.Consume(ctx, userID, 20.0, "购买道具")

	// 查询交易列表
	txList, err := service.ListTransactions(ctx, userID, &ListTransactionsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("查询交易列表失败: %v", err)
	}

	// 应该有3条交易记录（1充值 + 2消费）
	if len(txList) != 3 {
		t.Errorf("交易记录数量错误: got %d, want 3", len(txList))
	}
}

// TestUnifiedWalletService_Health 测试健康检查
func TestUnifiedWalletService_Health(t *testing.T) {
	ctx := context.Background()
	repo := NewMockWalletRepositoryV2()
	service := NewUnifiedWalletService(repo)

	err := service.Health(ctx)
	if err != nil {
		t.Fatalf("健康检查失败: %v", err)
	}
}
