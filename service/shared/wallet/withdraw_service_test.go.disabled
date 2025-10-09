package wallet

import (
	"context"
	"testing"
)

// TestCreateWithdrawRequest 测试创建提现请求
func TestCreateWithdrawRequest(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 创建提现请求
	request, err := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 50.0, "alipay", "user@example.com")

	if err != nil {
		t.Fatalf("创建提现请求失败: %v", err)
	}

	if request.Amount != 50.0 {
		t.Errorf("提现金额错误: %f", request.Amount)
	}

	if request.Status != "pending" {
		t.Errorf("提现状态应为pending，实际: %s", request.Status)
	}

	// 验证余额被冻结（扣除）
	balance, _ := walletService.GetBalance(ctx, wallet.ID)
	if balance != 50.0 {
		t.Errorf("余额应为50（100-50冻结），实际: %f", balance)
	}

	t.Logf("创建提现请求测试通过: amount=%.2f, status=%s", request.Amount, request.Status)
}

// TestCreateWithdrawRequest_InvalidAmount 测试无效提现金额
func TestCreateWithdrawRequest_InvalidAmount(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 提现负数
	_, err := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, -50.0, "alipay", "user@example.com")
	if err == nil {
		t.Error("提现负数应该失败")
	}

	// 提现小于最小金额
	_, err = withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 5.0, "alipay", "user@example.com")
	if err == nil {
		t.Error("提现金额小于10应该失败")
	}

	t.Logf("无效提现金额测试通过")
}

// TestCreateWithdrawRequest_InsufficientBalance 测试余额不足提现
func TestCreateWithdrawRequest_InsufficientBalance(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包（余额为0）
	wallet, _ := walletService.CreateWallet(ctx, "user123")

	// 尝试提现
	_, err := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 50.0, "alipay", "user@example.com")

	if err == nil {
		t.Error("余额不足提现应该失败")
	}

	t.Logf("余额不足提现测试通过: %v", err)
}

// TestApproveWithdraw 测试审核通过提现
func TestApproveWithdraw(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 创建提现请求
	request, _ := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 50.0, "alipay", "user@example.com")

	// 审核通过
	err := withdrawService.ApproveWithdraw(ctx, request.ID, "admin_001", "审核通过")
	if err != nil {
		t.Fatalf("审核通过失败: %v", err)
	}

	// 验证状态更新
	approved, _ := withdrawService.GetWithdrawRequest(ctx, request.ID)
	if approved.Status != "approved" {
		t.Errorf("状态应为approved，实际: %s", approved.Status)
	}

	if approved.ReviewerID != "admin_001" {
		t.Errorf("审核人错误: %s", approved.ReviewerID)
	}

	// 验证余额不变（已在申请时冻结）
	balance, _ := walletService.GetBalance(ctx, wallet.ID)
	if balance != 50.0 {
		t.Errorf("余额应保持50，实际: %f", balance)
	}

	t.Logf("审核通过测试通过: status=%s, reviewer=%s", approved.Status, approved.ReviewerID)
}

// TestRejectWithdraw 测试审核拒绝提现
func TestRejectWithdraw(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 创建提现请求
	request, _ := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 50.0, "alipay", "user@example.com")

	// 审核拒绝
	err := withdrawService.RejectWithdraw(ctx, request.ID, "admin_001", "账户信息错误")
	if err != nil {
		t.Fatalf("审核拒绝失败: %v", err)
	}

	// 验证状态更新
	rejected, _ := withdrawService.GetWithdrawRequest(ctx, request.ID)
	if rejected.Status != "rejected" {
		t.Errorf("状态应为rejected，实际: %s", rejected.Status)
	}

	// 验证余额退还
	balance, _ := walletService.GetBalance(ctx, wallet.ID)
	if balance != 100.0 {
		t.Errorf("余额应退还至100，实际: %f", balance)
	}

	t.Logf("审核拒绝测试通过: status=%s, balance=%.2f", rejected.Status, balance)
}

// TestGetWithdrawRequest 测试获取提现请求
func TestGetWithdrawRequest(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 100.0, "alipay", "order_001")

	// 创建提现请求
	created, _ := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 50.0, "alipay", "user@example.com")

	// 获取提现请求
	request, err := withdrawService.GetWithdrawRequest(ctx, created.ID)
	if err != nil {
		t.Fatalf("获取提现请求失败: %v", err)
	}

	if request.ID != created.ID {
		t.Errorf("提现请求ID不匹配")
	}

	if request.Amount != 50.0 {
		t.Errorf("提现金额错误: %f", request.Amount)
	}

	t.Logf("获取提现请求测试通过: %+v", request)
}

// TestListWithdrawRequests 测试列出提现请求
func TestListWithdrawRequests(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 500.0, "alipay", "order_001")

	// 创建多个提现请求
	withdrawService.CreateWithdrawRequest(ctx, "user123", wallet.ID, 50.0, "alipay", "user@example.com")
	withdrawService.CreateWithdrawRequest(ctx, "user123", wallet.ID, 100.0, "wechat", "user@example.com")

	// 列出所有提现请求
	requests, err := withdrawService.ListWithdrawRequests(ctx, "user123", "", 10, 0)
	if err != nil {
		t.Fatalf("列出提现请求失败: %v", err)
	}

	if len(requests) != 2 {
		t.Errorf("提现请求数量错误: %d", len(requests))
	}

	t.Logf("列出提现请求测试通过，共%d条记录", len(requests))
	for i, req := range requests {
		t.Logf("  [%d] amount=%.2f, status=%s", i+1, req.Amount, req.Status)
	}
}

// TestListWithdrawRequests_FilterByStatus 测试按状态筛选提现请求
func TestListWithdrawRequests_FilterByStatus(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 创建钱包并充值
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	txService.Recharge(ctx, wallet.ID, 500.0, "alipay", "order_001")

	// 创建多个提现请求
	req1, _ := withdrawService.CreateWithdrawRequest(ctx, "user123", wallet.ID, 50.0, "alipay", "user@example.com")
	req2, _ := withdrawService.CreateWithdrawRequest(ctx, "user123", wallet.ID, 100.0, "wechat", "user@example.com")

	// 审核第一个请求
	withdrawService.ApproveWithdraw(ctx, req1.ID, "admin_001", "审核通过")

	// 拒绝第二个请求
	withdrawService.RejectWithdraw(ctx, req2.ID, "admin_001", "账户错误")

	// 筛选pending状态（应该没有）
	pendingRequests, _ := withdrawService.ListWithdrawRequests(ctx, "user123", "pending", 10, 0)
	if len(pendingRequests) != 0 {
		t.Errorf("pending状态的请求应为0，实际: %d", len(pendingRequests))
	}

	// 筛选approved状态（应该有1个）
	approvedRequests, _ := withdrawService.ListWithdrawRequests(ctx, "user123", "approved", 10, 0)
	if len(approvedRequests) != 1 {
		t.Errorf("approved状态的请求应为1，实际: %d", len(approvedRequests))
	}

	// 筛选rejected状态（应该有1个）
	rejectedRequests, _ := withdrawService.ListWithdrawRequests(ctx, "user123", "rejected", 10, 0)
	if len(rejectedRequests) != 1 {
		t.Errorf("rejected状态的请求应为1，实际: %d", len(rejectedRequests))
	}

	t.Logf("按状态筛选测试通过: pending=%d, approved=%d, rejected=%d",
		len(pendingRequests), len(approvedRequests), len(rejectedRequests))
}

// TestWithdrawWorkflow 测试完整提现流程
func TestWithdrawWorkflow(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	walletService := NewWalletService(mockRepo)
	txService := NewTransactionService(mockRepo)
	withdrawService := NewWithdrawService(mockRepo)

	ctx := context.Background()

	// 1. 创建钱包
	wallet, _ := walletService.CreateWallet(ctx, "user123")
	t.Logf("1. 创建钱包: ID=%s", wallet.ID)

	// 2. 充值
	txService.Recharge(ctx, wallet.ID, 200.0, "alipay", "order_001")
	balance1, _ := walletService.GetBalance(ctx, wallet.ID)
	t.Logf("2. 充值200元, 余额: %.2f", balance1)

	// 3. 申请提现
	request, _ := withdrawService.CreateWithdrawRequest(ctx,
		"user123", wallet.ID, 100.0, "alipay", "user@example.com")
	balance2, _ := walletService.GetBalance(ctx, wallet.ID)
	t.Logf("3. 申请提现100元（冻结）, 余额: %.2f, 状态: %s", balance2, request.Status)

	// 4. 审核通过
	withdrawService.ApproveWithdraw(ctx, request.ID, "admin_001", "审核通过")
	approved, _ := withdrawService.GetWithdrawRequest(ctx, request.ID)
	balance3, _ := walletService.GetBalance(ctx, wallet.ID)
	t.Logf("4. 审核通过, 余额: %.2f, 状态: %s", balance3, approved.Status)

	// 验证最终余额
	if balance3 != 100.0 {
		t.Errorf("最终余额应为100，实际: %f", balance3)
	}

	// 验证状态
	if approved.Status != "approved" {
		t.Errorf("状态应为approved，实际: %s", approved.Status)
	}

	t.Logf("完整提现流程测试通过 ✅")
}
