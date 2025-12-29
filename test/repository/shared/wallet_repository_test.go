package shared_test

import (
	walletModel "Qingyu_backend/models/wallet"
	sharedInterfaces "Qingyu_backend/repository/interfaces/shared"
	"Qingyu_backend/repository/mongodb/shared"
	"Qingyu_backend/test/testutil"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 测试辅助函数
func setupWalletRepo(t *testing.T) (sharedInterfaces.WalletRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := shared.NewWalletRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

func createTestWallet(userID string) *walletModel.Wallet {
	return &walletModel.Wallet{
		UserID:  userID,
		Balance: 100.0,
		Frozen:  false,
	}
}

func createTestTransaction(userID string, amount float64) *walletModel.Transaction {
	return &walletModel.Transaction{
		UserID:          userID,
		Type:            walletModel.TransactionTypeRecharge,
		Amount:          amount,
		Balance:         100 + amount,
		Status:          walletModel.TransactionStatusSuccess,
		TransactionTime: time.Now(),
	}
}

// 1. 测试创建钱包
func TestWalletRepository_CreateWallet(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	wallet := createTestWallet("user123")
	err := repo.CreateWallet(ctx, wallet)

	require.NoError(t, err)
	assert.NotEmpty(t, wallet.ID)
	assert.NotZero(t, wallet.CreatedAt)
	assert.NotZero(t, wallet.UpdatedAt)
}

// 2. 测试获取钱包（根据UserID）
func TestWalletRepository_GetWallet(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	// 创建钱包
	userID := "user123"
	wallet := createTestWallet(userID)
	err := repo.CreateWallet(ctx, wallet)
	require.NoError(t, err)

	// 获取钱包
	retrieved, err := repo.GetWallet(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, userID, retrieved.UserID)
	assert.Equal(t, 100.0, retrieved.Balance)
}

// 3. 测试获取不存在的钱包
func TestWalletRepository_GetWallet_NotFound(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	_, err := repo.GetWallet(ctx, "non_existent_user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "钱包不存在")
}

// 4. 测试更新钱包
func TestWalletRepository_UpdateWallet(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	// 创建钱包
	wallet := createTestWallet("user123")
	err := repo.CreateWallet(ctx, wallet)
	require.NoError(t, err)

	// 更新钱包
	updates := map[string]interface{}{
		"frozen":    true,
		"frozen_at": time.Now(),
	}
	err = repo.UpdateWallet(ctx, wallet.ID, updates)
	require.NoError(t, err)
}

// 5. 测试更新余额（原子操作）
func TestWalletRepository_UpdateBalance(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	// 创建钱包
	wallet := createTestWallet("user123")
	err := repo.CreateWallet(ctx, wallet)
	require.NoError(t, err)

	// 增加余额
	err = repo.UpdateBalance(ctx, wallet.ID, 50.0)
	require.NoError(t, err)

	// 减少余额
	err = repo.UpdateBalance(ctx, wallet.ID, -20.0)
	require.NoError(t, err)
}

// 6. 测试创建交易记录
func TestWalletRepository_CreateTransaction(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	transaction := createTestTransaction("user123", 50.0)
	err := repo.CreateTransaction(ctx, transaction)

	require.NoError(t, err)
	assert.NotEmpty(t, transaction.ID)
	assert.NotZero(t, transaction.CreatedAt)
}

// 7. 测试获取交易记录
func TestWalletRepository_GetTransaction(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	// 创建交易
	transaction := createTestTransaction("user123", 50.0)
	err := repo.CreateTransaction(ctx, transaction)
	require.NoError(t, err)

	// 获取交易
	retrieved, err := repo.GetTransaction(ctx, transaction.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, transaction.UserID, retrieved.UserID)
	assert.Equal(t, 50.0, retrieved.Amount)
}

// 8. 测试列出交易记录
func TestWalletRepository_ListTransactions(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建多个交易
	for i := 0; i < 3; i++ {
		transaction := createTestTransaction(userID, float64(i+1)*10)
		err := repo.CreateTransaction(ctx, transaction)
		require.NoError(t, err)
	}

	// 查询交易列表
	filter := &sharedInterfaces.TransactionFilter{
		UserID: userID,
		Limit:  10,
	}
	transactions, err := repo.ListTransactions(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, transactions, 3)
}

// 9. 测试统计交易数量
func TestWalletRepository_CountTransactions(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建多个交易
	for i := 0; i < 5; i++ {
		transaction := createTestTransaction(userID, float64(i+1)*10)
		err := repo.CreateTransaction(ctx, transaction)
		require.NoError(t, err)
	}

	// 统计交易数量
	filter := &sharedInterfaces.TransactionFilter{
		UserID: userID,
	}
	count, err := repo.CountTransactions(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

// 10. 测试创建提现请求
func TestWalletRepository_CreateWithdrawRequest(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	request := &walletModel.WithdrawRequest{
		UserID:       "user123",
		Amount:       100.0,
		Fee:          1.0,
		ActualAmount: 99.0,
		Account:      "alipay@test.com",
		AccountType:  "alipay",
		Status:       walletModel.WithdrawStatusPending,
		OrderNo:      "WD" + primitive.NewObjectID().Hex(),
	}

	err := repo.CreateWithdrawRequest(ctx, request)
	require.NoError(t, err)
	assert.NotEmpty(t, request.ID)
	assert.NotZero(t, request.CreatedAt)
}

// 11. 测试获取提现请求
func TestWalletRepository_GetWithdrawRequest(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	// 创建提现请求
	request := &walletModel.WithdrawRequest{
		UserID:       "user123",
		Amount:       100.0,
		Fee:          1.0,
		ActualAmount: 99.0,
		Account:      "alipay@test.com",
		AccountType:  "alipay",
		Status:       walletModel.WithdrawStatusPending,
		OrderNo:      "WD" + primitive.NewObjectID().Hex(),
	}
	err := repo.CreateWithdrawRequest(ctx, request)
	require.NoError(t, err)

	// 获取提现请求
	retrieved, err := repo.GetWithdrawRequest(ctx, request.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, request.UserID, retrieved.UserID)
	assert.Equal(t, 100.0, retrieved.Amount)
}

// 12. 测试更新提现请求
func TestWalletRepository_UpdateWithdrawRequest(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	// 创建提现请求
	request := &walletModel.WithdrawRequest{
		UserID:       "user123",
		Amount:       100.0,
		Fee:          1.0,
		ActualAmount: 99.0,
		Account:      "alipay@test.com",
		AccountType:  "alipay",
		Status:       walletModel.WithdrawStatusPending,
		OrderNo:      "WD" + primitive.NewObjectID().Hex(),
	}
	err := repo.CreateWithdrawRequest(ctx, request)
	require.NoError(t, err)

	// 更新提现请求
	updates := map[string]interface{}{
		"status":      walletModel.WithdrawStatusApproved,
		"reviewed_by": "admin123",
		"reviewed_at": time.Now(),
	}
	err = repo.UpdateWithdrawRequest(ctx, request.ID, updates)
	require.NoError(t, err)
}

// 13. 测试列出提现请求
func TestWalletRepository_ListWithdrawRequests(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建多个提现请求
	for i := 0; i < 3; i++ {
		request := &walletModel.WithdrawRequest{
			UserID:       userID,
			Amount:       float64(i+1) * 100,
			Fee:          1.0,
			ActualAmount: float64(i+1)*100 - 1,
			Account:      "alipay@test.com",
			AccountType:  "alipay",
			Status:       walletModel.WithdrawStatusPending,
			OrderNo:      "WD" + primitive.NewObjectID().Hex(),
		}
		err := repo.CreateWithdrawRequest(ctx, request)
		require.NoError(t, err)
	}

	// 查询提现列表
	filter := &sharedInterfaces.WithdrawFilter{
		UserID: userID,
		Limit:  10,
	}
	requests, err := repo.ListWithdrawRequests(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, requests, 3)
}

// 14. 测试统计提现请求数量
func TestWalletRepository_CountWithdrawRequests(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建多个提现请求
	for i := 0; i < 4; i++ {
		request := &walletModel.WithdrawRequest{
			UserID:       userID,
			Amount:       float64(i+1) * 100,
			Fee:          1.0,
			ActualAmount: float64(i+1)*100 - 1,
			Account:      "alipay@test.com",
			AccountType:  "alipay",
			Status:       walletModel.WithdrawStatusPending,
			OrderNo:      "WD" + primitive.NewObjectID().Hex(),
		}
		err := repo.CreateWithdrawRequest(ctx, request)
		require.NoError(t, err)
	}

	// 统计提现数量
	filter := &sharedInterfaces.WithdrawFilter{
		UserID: userID,
	}
	count, err := repo.CountWithdrawRequests(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(4), count)
}

// 15. 测试健康检查
func TestWalletRepository_Health(t *testing.T) {
	repo, ctx, cleanup := setupWalletRepo(t)
	defer cleanup()

	err := repo.Health(ctx)
	assert.NoError(t, err)
}
