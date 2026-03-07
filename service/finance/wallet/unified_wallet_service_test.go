package wallet

import (
	"context"
	"testing"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnifiedWalletService_CreateWallet(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	result, err := svc.CreateWallet(context.Background(), "user123")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user123", result.UserID)
	assert.Equal(t, int64(0), result.Balance)
}

func TestUnifiedWalletService_GetWallet(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
	}

	svc := NewUnifiedWalletService(repo)

	result, err := svc.GetWallet(context.Background(), "user123")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1000), result.Balance)
}

func TestUnifiedWalletService_GetBalance(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(5000),
	}

	svc := NewUnifiedWalletService(repo)

	balance, err := svc.GetBalance(context.Background(), "user123")
	require.NoError(t, err)
	assert.Equal(t, int64(5000), balance)
}

func TestUnifiedWalletService_FreezeWallet(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
		Frozen:  false,
	}

	svc := NewUnifiedWalletService(repo)

	err := svc.FreezeWallet(context.Background(), "user123")
	require.NoError(t, err)
	assert.True(t, repo.wallets["user123"].Frozen)
}

func TestUnifiedWalletService_UnfreezeWallet(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
		Frozen:  true,
	}

	svc := NewUnifiedWalletService(repo)

	err := svc.UnfreezeWallet(context.Background(), "user123")
	require.NoError(t, err)
	assert.False(t, repo.wallets["user123"].Frozen)
}

func TestUnifiedWalletService_Recharge(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Recharge(context.Background(), "user123", 500, "alipay")
	require.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, int64(1500), int64(repo.wallets["user123"].Balance))
}

func TestUnifiedWalletService_Recharge_WalletNotFound(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Recharge(context.Background(), "user_not_exist", 500, "alipay")
	require.Error(t, err)
	assert.Nil(t, tx)
	assert.Contains(t, err.Error(), "获取钱包失败")
}

func TestUnifiedWalletService_Consume(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Consume(context.Background(), "user123", 300, "购买书籍")
	require.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, int64(700), int64(repo.wallets["user123"].Balance))
}

func TestUnifiedWalletService_Consume_WalletNotFound(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Consume(context.Background(), "user_not_exist", 300, "购买书籍")
	require.Error(t, err)
	assert.Nil(t, tx)
	assert.Contains(t, err.Error(), "获取钱包失败")
}

func TestUnifiedWalletService_Transfer(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_from"] = &financeModel.Wallet{
		UserID:  "user_from",
		Balance: types.Money(1000),
	}
	repo.wallets["user_to"] = &financeModel.Wallet{
		UserID:  "user_to",
		Balance: types.Money(500),
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Transfer(context.Background(), "user_from", "user_to", 300, "转账")
	require.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, int64(700), int64(repo.wallets["user_from"].Balance))
	assert.Equal(t, int64(800), int64(repo.wallets["user_to"].Balance))
}

func TestUnifiedWalletService_Transfer_SourceWalletNotFound(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_to"] = &financeModel.Wallet{
		UserID:  "user_to",
		Balance: types.Money(500),
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Transfer(context.Background(), "user_from_not_exist", "user_to", 300, "转账")
	require.Error(t, err)
	assert.Nil(t, tx)
	assert.Contains(t, err.Error(), "获取源钱包失败")
}

func TestUnifiedWalletService_Transfer_TargetWalletNotFound(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_from"] = &financeModel.Wallet{
		UserID:  "user_from",
		Balance: types.Money(1000),
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Transfer(context.Background(), "user_from", "user_to_not_exist", 300, "转账")
	require.Error(t, err)
	assert.Nil(t, tx)
	assert.Contains(t, err.Error(), "获取目标钱包失败")
}

func TestUnifiedWalletService_GetTransaction(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
	}

	svc := NewUnifiedWalletService(repo)

	// 先创建一笔交易
	tx, err := svc.Recharge(context.Background(), "user123", 500, "alipay")
	require.NoError(t, err)

	// 获取交易
	result, err := svc.GetTransaction(context.Background(), tx.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tx.ID, result.ID)
}

func TestUnifiedWalletService_ListTransactions(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000),
	}

	svc := NewUnifiedWalletService(repo)

	// 创建几笔交易
	_, _ = svc.Recharge(context.Background(), "user123", 500, "alipay")
	_, _ = svc.Recharge(context.Background(), "user123", 300, "wechat")

	// 列出交易
	result, err := svc.ListTransactions(context.Background(), "user123", &ListTransactionsRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUnifiedWalletService_RequestWithdraw(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(10000),
	}

	svc := NewUnifiedWalletService(repo)

	req, err := svc.RequestWithdraw(context.Background(), "user123", 5000, "alipay_account")
	require.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, int64(5000), req.Amount)
}

func TestUnifiedWalletService_GetWithdrawRequest(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(10000),
	}

	svc := NewUnifiedWalletService(repo)

	// 先创建提现请求
	req, err := svc.RequestWithdraw(context.Background(), "user123", 5000, "alipay_account")
	require.NoError(t, err)

	// 获取提现请求
	result, err := svc.GetWithdrawRequest(context.Background(), req.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUnifiedWalletService_ListWithdrawRequests(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(10000),
	}

	svc := NewUnifiedWalletService(repo)

	// 创建几个提现请求
	_, _ = svc.RequestWithdraw(context.Background(), "user123", 1000, "alipay_account")
	_, _ = svc.RequestWithdraw(context.Background(), "user123", 2000, "wechat_account")

	// 列出提现请求
	result, err := svc.ListWithdrawRequests(context.Background(), &ListWithdrawRequestsRequest{
		UserID:   "user123",
		Status:   "pending",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUnifiedWalletService_ApproveWithdraw(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(10000),
	}

	svc := NewUnifiedWalletService(repo)

	// 创建提现请求
	req, err := svc.RequestWithdraw(context.Background(), "user123", 5000, "alipay_account")
	require.NoError(t, err)

	// 批准提现
	err = svc.ApproveWithdraw(context.Background(), req.ID, "admin")
	require.NoError(t, err)
}

func TestUnifiedWalletService_RejectWithdraw(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(10000),
	}

	svc := NewUnifiedWalletService(repo)

	// 创建提现请求
	req, err := svc.RequestWithdraw(context.Background(), "user123", 5000, "alipay_account")
	require.NoError(t, err)

	// 拒绝提现
	err = svc.RejectWithdraw(context.Background(), req.ID, "admin", "审核不通过")
	require.NoError(t, err)
}

func TestUnifiedWalletService_Health(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	// 先初始化服务
	err := svc.Initialize(context.Background())
	require.NoError(t, err)

	// 然后检查健康状态
	err = svc.Health(context.Background())
	require.NoError(t, err)
}

func TestNewUnifiedWalletServiceWithRunner(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	runner := NewRepositoryTransactionRunner(repo)

	svc := NewUnifiedWalletServiceWithRunner(repo, runner)
	assert.NotNil(t, svc)
}

func TestUnifiedWalletService_Initialize(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	err := svc.Initialize(context.Background())
	require.NoError(t, err)
}

func TestUnifiedWalletService_GetServiceName(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	name := svc.GetServiceName()
	assert.NotEmpty(t, name)
}

func TestUnifiedWalletService_GetVersion(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	version := svc.GetVersion()
	assert.NotEmpty(t, version)
}

func TestUnifiedWalletService_ProcessWithdraw(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(10000),
	}

	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	// 创建提现请求
	req, err := svc.RequestWithdraw(context.Background(), "user123", 5000, "alipay_account")
	require.NoError(t, err)

	// 先批准提现请求
	err = svc.ApproveWithdraw(context.Background(), req.ID, "admin")
	require.NoError(t, err)

	// 处理提现（批准后才能处理）
	err = svc.ProcessWithdraw(context.Background(), req.ID)
	require.NoError(t, err)
}

func TestUnifiedWalletService_Close(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	err := svc.Close(context.Background())
	require.NoError(t, err)
}

func TestUnifiedWalletService_Initialize_AlreadyInitialized(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	svc := NewUnifiedWalletService(repo).(*UnifiedWalletService)

	// 第一次初始化
	err := svc.Initialize(context.Background())
	require.NoError(t, err)

	// 再次初始化应该返回nil（幂等）
	err = svc.Initialize(context.Background())
	require.NoError(t, err)
}

func TestUnifiedWalletService_Transfer_InsufficientBalance(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user_from"] = &financeModel.Wallet{
		UserID:  "user_from",
		Balance: types.Money(100), // 余额不足
	}
	repo.wallets["user_to"] = &financeModel.Wallet{
		UserID:  "user_to",
		Balance: types.Money(500),
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Transfer(context.Background(), "user_from", "user_to", 300, "转账")
	require.Error(t, err)
	assert.Nil(t, tx)
}

func TestUnifiedWalletService_Consume_InsufficientBalance(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(100), // 余额不足
	}

	svc := NewUnifiedWalletService(repo)

	tx, err := svc.Consume(context.Background(), "user123", 300, "购买书籍")
	require.Error(t, err)
	assert.Nil(t, tx)
}

func TestUnifiedWalletService_RequestWithdraw_InsufficientBalance(t *testing.T) {
	repo := NewMockWalletRepositoryV2()
	repo.wallets["user123"] = &financeModel.Wallet{
		UserID:  "user123",
		Balance: types.Money(1000), // 余额不足
	}

	svc := NewUnifiedWalletService(repo)

	req, err := svc.RequestWithdraw(context.Background(), "user123", 5000, "alipay_account")
	require.Error(t, err)
	assert.Nil(t, req)
}
