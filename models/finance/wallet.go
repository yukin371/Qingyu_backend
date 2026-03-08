package finance

import (
	"Qingyu_backend/models/shared/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wallet 钱包模型
type Wallet struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Balance   types.Money        `json:"-" bson:"balance_cents"`
	Frozen    bool               `json:"frozen" bson:"frozen"`
	FrozenAt  time.Time          `json:"frozen_at,omitempty" bson:"frozen_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Transaction 交易记录
type Transaction struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	Type            string             `json:"type" bson:"type"`
	Amount          types.Money        `json:"-" bson:"amount_cents"`
	Balance         types.Money        `json:"-" bson:"balance_cents"`
	RelatedUserID   string             `json:"related_user_id,omitempty" bson:"related_user_id,omitempty"`
	Method          string             `json:"method,omitempty" bson:"method,omitempty"`
	Reason          string             `json:"reason,omitempty" bson:"reason,omitempty"`
	Status          string             `json:"status" bson:"status"`
	OrderNo         string             `json:"order_no,omitempty" bson:"order_no,omitempty"`
	ThirdPartyNo    string             `json:"third_party_no,omitempty" bson:"third_party_no,omitempty"`
	TransactionTime time.Time          `json:"transaction_time" bson:"transaction_time"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// WithdrawRequest 提现申请
type WithdrawRequest struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        string             `json:"user_id" bson:"user_id"`
	Amount        types.Money        `json:"-" bson:"amount_cents"`
	Fee           types.Money        `json:"-" bson:"fee_cents"`
	ActualAmount  types.Money        `json:"-" bson:"actual_amount_cents"`
	Account       string             `json:"account" bson:"account"`
	AccountType   string             `json:"account_type" bson:"account_type"`
	AccountName   string             `json:"account_name,omitempty" bson:"account_name,omitempty"`
	Status        string             `json:"status" bson:"status"`
	ReviewedBy    string             `json:"reviewed_by,omitempty" bson:"reviewed_by,omitempty"`
	ReviewedAt    time.Time          `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
	RejectReason  string             `json:"reject_reason,omitempty" bson:"reject_reason,omitempty"`
	ProcessedAt   time.Time          `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	TransactionID string             `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	OrderNo       string             `json:"order_no" bson:"order_no"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// 交易类型
const (
	TransactionTypeRecharge    = "recharge"     // 充值
	TransactionTypeConsume     = "consume"      // 消费
	TransactionTypeTransferIn  = "transfer_in"  // 转入
	TransactionTypeTransferOut = "transfer_out" // 转出
	TransactionTypeWithdraw    = "withdraw"     // 提现
	TransactionTypeRefund      = "refund"       // 退款
)

// 交易状态
const (
	TransactionStatusPending = "pending" // 处理中
	TransactionStatusSuccess = "success" // 成功
	TransactionStatusFailed  = "failed"  // 失败
)

// 支付方式
const (
	PaymentMethodAlipay = "alipay" // 支付宝
	PaymentMethodWechat = "wechat" // 微信支付
	PaymentMethodBank   = "bank"   // 银行卡
)
