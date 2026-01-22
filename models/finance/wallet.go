package finance

import "time"

// Wallet 钱包模型
type Wallet struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"` // 用户ID（唯一）
	Balance   int64     `json:"balance" bson:"balance"` // 余额 (分)
	Frozen    bool      `json:"frozen" bson:"frozen"`   // 是否冻结
	FrozenAt  time.Time `json:"frozen_at,omitempty" bson:"frozen_at,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Transaction 交易记录
type Transaction struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	UserID          string    `json:"user_id" bson:"user_id"`                                     // 用户ID
	Type            string    `json:"type" bson:"type"`                                           // 交易类型
	Amount          int64     `json:"amount" bson:"amount"`                                       // 交易金额 (分)
	Balance         int64     `json:"balance" bson:"balance"`                                     // 交易后余额 (分)
	RelatedUserID   string    `json:"related_user_id,omitempty" bson:"related_user_id,omitempty"` // 关联用户（转账）
	Method          string    `json:"method,omitempty" bson:"method,omitempty"`                   // 支付方式
	Reason          string    `json:"reason,omitempty" bson:"reason,omitempty"`                   // 交易原因
	Status          string    `json:"status" bson:"status"`                                       // 交易状态
	OrderNo         string    `json:"order_no,omitempty" bson:"order_no,omitempty"`               // 订单号（雪花算法）
	ThirdPartyNo    string    `json:"third_party_no,omitempty" bson:"third_party_no,omitempty"`   // 第三方流水号
	TransactionTime time.Time `json:"transaction_time" bson:"transaction_time"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// WithdrawRequest 提现申请
type WithdrawRequest struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	UserID        string    `json:"user_id" bson:"user_id"`                               // 用户ID
	Amount        int64     `json:"amount" bson:"amount"`                                 // 提现金额 (分)
	Fee           int64     `json:"fee" bson:"fee"`                                       // 手续费 (分)
	ActualAmount  int64     `json:"actual_amount" bson:"actual_amount"`                   // 实际到账金额 (分)
	Account       string    `json:"account" bson:"account"`                               // 提现账号
	AccountType   string    `json:"account_type" bson:"account_type"`                     // 账号类型：alipay, wechat, bank
	AccountName   string    `json:"account_name,omitempty" bson:"account_name,omitempty"` // 账户名
	Status        string    `json:"status" bson:"status"`                                 // 状态
	ReviewedBy    string    `json:"reviewed_by,omitempty" bson:"reviewed_by,omitempty"`
	ReviewedAt    time.Time `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
	RejectReason  string    `json:"reject_reason,omitempty" bson:"reject_reason,omitempty"`
	ProcessedAt   time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"` // 关联交易ID
	OrderNo       string    `json:"order_no" bson:"order_no"`                                 // 订单号
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
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
