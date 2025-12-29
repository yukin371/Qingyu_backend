package mocks

import (
	walletModel "Qingyu_backend/models/wallet"
	"context"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"

	"github.com/stretchr/testify/mock"
)

// MockWalletRepository WalletRepository 的 Mock 实现
type MockWalletRepository struct {
	mock.Mock
}

// CreateWallet 创建钱包
func (m *MockWalletRepository) CreateWallet(ctx context.Context, wallet *walletModel.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

// GetWallet 获取钱包
func (m *MockWalletRepository) GetWallet(ctx context.Context, userID string) (*walletModel.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletModel.Wallet), args.Error(1)
}

// UpdateWallet 更新钱包
func (m *MockWalletRepository) UpdateWallet(ctx context.Context, userID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

// UpdateBalance 更新余额
func (m *MockWalletRepository) UpdateBalance(ctx context.Context, userID string, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

// CreateTransaction 创建交易记录
func (m *MockWalletRepository) CreateTransaction(ctx context.Context, transaction *walletModel.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

// GetTransaction 获取交易记录
func (m *MockWalletRepository) GetTransaction(ctx context.Context, transactionID string) (*walletModel.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletModel.Transaction), args.Error(1)
}

// ListTransactions 列出交易记录
func (m *MockWalletRepository) ListTransactions(ctx context.Context, filter *sharedRepo.TransactionFilter) ([]*walletModel.Transaction, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*walletModel.Transaction), args.Error(1)
}

// CountTransactions 统计交易数量
func (m *MockWalletRepository) CountTransactions(ctx context.Context, filter *sharedRepo.TransactionFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// CreateWithdrawRequest 创建提现申请
func (m *MockWalletRepository) CreateWithdrawRequest(ctx context.Context, request *walletModel.WithdrawRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

// GetWithdrawRequest 获取提现申请
func (m *MockWalletRepository) GetWithdrawRequest(ctx context.Context, withdrawID string) (*walletModel.WithdrawRequest, error) {
	args := m.Called(ctx, withdrawID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*walletModel.WithdrawRequest), args.Error(1)
}

// UpdateWithdrawRequest 更新提现申请
func (m *MockWalletRepository) UpdateWithdrawRequest(ctx context.Context, withdrawID string, updates map[string]interface{}) error {
	args := m.Called(ctx, withdrawID, updates)
	return args.Error(0)
}

// ListWithdrawRequests 列出提现申请
func (m *MockWalletRepository) ListWithdrawRequests(ctx context.Context, filter *sharedRepo.WithdrawFilter) ([]*walletModel.WithdrawRequest, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*walletModel.WithdrawRequest), args.Error(1)
}

// CountWithdrawRequests 统计提现申请数量
func (m *MockWalletRepository) CountWithdrawRequests(ctx context.Context, filter *sharedRepo.WithdrawFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// Health 健康检查
func (m *MockWalletRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
