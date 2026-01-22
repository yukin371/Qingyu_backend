package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"context"
	"time"
)

// WalletRepository 钱包相关Repository
type WalletRepository interface {
	// 钱包管理
	CreateWallet(ctx context.Context, wallet *financeModel.Wallet) error
	GetWallet(ctx context.Context, userID string) (*financeModel.Wallet, error)
	UpdateWallet(ctx context.Context, userID string, updates map[string]interface{}) error
	UpdateBalance(ctx context.Context, userID string, amount int64) error

	// 交易记录
	CreateTransaction(ctx context.Context, transaction *financeModel.Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (*financeModel.Transaction, error)
	ListTransactions(ctx context.Context, filter *TransactionFilter) ([]*financeModel.Transaction, error)
	CountTransactions(ctx context.Context, filter *TransactionFilter) (int64, error)

	// 提现申请
	CreateWithdrawRequest(ctx context.Context, request *financeModel.WithdrawRequest) error
	GetWithdrawRequest(ctx context.Context, withdrawID string) (*financeModel.WithdrawRequest, error)
	UpdateWithdrawRequest(ctx context.Context, withdrawID string, updates map[string]interface{}) error
	ListWithdrawRequests(ctx context.Context, filter *WithdrawFilter) ([]*financeModel.WithdrawRequest, error)
	CountWithdrawRequests(ctx context.Context, filter *WithdrawFilter) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// TransactionFilter 交易过滤器
type TransactionFilter struct {
	UserID    string
	Type      string
	Status    string
	StartDate time.Time
	EndDate   time.Time
	Limit     int64
	Offset    int64
}

// WithdrawFilter 提现过滤器
type WithdrawFilter struct {
	UserID    string
	Status    string
	StartDate time.Time
	EndDate   time.Time
	Limit     int64
	Offset    int64
}
