package wallet

import (
	"context"
	"time"
)

// WalletService 钱包服务接口（对外暴露）
type WalletService interface {
	// 钱包管理
	CreateWallet(ctx context.Context, userID string) (*Wallet, error)
	GetWallet(ctx context.Context, userID string) (*Wallet, error)
	GetBalance(ctx context.Context, userID string) (int64, error)
	FreezeWallet(ctx context.Context, userID string) error
	UnfreezeWallet(ctx context.Context, userID string) error

	// 交易操作
	Recharge(ctx context.Context, userID string, amount int64, method string) (*Transaction, error)
	Consume(ctx context.Context, userID string, amount int64, reason string) (*Transaction, error)
	Transfer(ctx context.Context, fromUserID, toUserID string, amount int64, reason string) (*Transaction, error)

	// 交易查询
	GetTransaction(ctx context.Context, transactionID string) (*Transaction, error)
	ListTransactions(ctx context.Context, userID string, req *ListTransactionsRequest) ([]*Transaction, error)

	// 提现管理
	RequestWithdraw(ctx context.Context, userID string, amount int64, account string) (*WithdrawRequest, error)
	GetWithdrawRequest(ctx context.Context, withdrawID string) (*WithdrawRequest, error)
	ListWithdrawRequests(ctx context.Context, req *ListWithdrawRequestsRequest) ([]*WithdrawRequest, error)
	ApproveWithdraw(ctx context.Context, withdrawID, adminID string) error
	RejectWithdraw(ctx context.Context, withdrawID, adminID, reason string) error
	ProcessWithdraw(ctx context.Context, withdrawID string) error

	// Health 健康检查
	Health(ctx context.Context) error
}

// ============ 请求/响应结构 ============

// ListTransactionsRequest 查询交易列表请求
type ListTransactionsRequest struct {
	TransactionType string    `json:"transaction_type"` // recharge, consume, transfer
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Page            int       `json:"page"`
	PageSize        int       `json:"page_size"`
}

// ListWithdrawRequestsRequest 查询提现列表请求
type ListWithdrawRequestsRequest struct {
	Status   string `json:"status"` // pending, approved, rejected, processed
	UserID   string `json:"user_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// ============ 数据结构 ============

// Wallet 钱包
type Wallet struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Balance   int64     `json:"balance" bson:"balance"` // 余额 (分)
	Frozen    bool      `json:"frozen" bson:"frozen"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Transaction 交易记录
type Transaction struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	UserID          string    `json:"user_id" bson:"user_id"`
	Type            string    `json:"type" bson:"type"` // recharge, consume, transfer_in, transfer_out
	Amount          int64     `json:"amount" bson:"amount"`   // 交易金额 (分)
	Balance         int64     `json:"balance" bson:"balance"` // 交易后余额 (分)
	RelatedUserID   string    `json:"related_user_id,omitempty" bson:"related_user_id,omitempty"`
	Method          string    `json:"method,omitempty" bson:"method,omitempty"` // 支付方式
	Reason          string    `json:"reason,omitempty" bson:"reason,omitempty"`
	Status          string    `json:"status" bson:"status"` // pending, success, failed
	TransactionTime time.Time `json:"transaction_time" bson:"transaction_time"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
}

// WithdrawRequest 提现申请
type WithdrawRequest struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	UserID        string    `json:"user_id" bson:"user_id"`
	Amount        int64     `json:"amount" bson:"amount"` // 提现金额 (分)
	Account       string    `json:"account" bson:"account"` // 提现账号
	Status        string    `json:"status" bson:"status"`   // pending, approved, rejected, processed
	ReviewedBy    string    `json:"reviewed_by,omitempty" bson:"reviewed_by,omitempty"`
	ReviewedAt    time.Time `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
	RejectReason  string    `json:"reject_reason,omitempty" bson:"reject_reason,omitempty"`
	ProcessedAt   time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}
