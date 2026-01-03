package finance

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthorEarning 作者收入记录
type AuthorEarning struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID        string             `bson:"author_id" json:"author_id"`           // 作者ID
	BookID          primitive.ObjectID `bson:"book_id" json:"book_id"`               // 书籍ID
	BookTitle       string             `bson:"book_title" json:"book_title"`         // 书名
	ChapterID       primitive.ObjectID `bson:"chapter_id,omitempty" json:"chapter_id,omitempty"` // 章节ID（如果是章节购买）
	ChapterTitle    string             `bson:"chapter_title,omitempty" json:"chapter_title,omitempty"` // 章节标题
	Type            string             `bson:"type" json:"type"`                     // 收入类型：chapter_purchase, reward, vip_reading
	Amount          float64            `bson:"amount" json:"amount"`                 // 收入金额
	ReaderID        string             `bson:"reader_id,omitempty" json:"reader_id,omitempty"` // 读者ID
	ReaderNickname  string             `bson:"reader_nickname,omitempty" json:"reader_nickname,omitempty"` // 读者昵称
	PlatformFee     float64            `bson:"platform_fee" json:"platform_fee"`     // 平台抽成
	AuthorIncome    float64            `bson:"author_income" json:"author_income"`   // 作者收入
	WordCount       int                `bson:"word_count,omitempty" json:"word_count,omitempty"` // 字数（VIP阅读）
	SettlementID    primitive.ObjectID `bson:"settlement_id,omitempty" json:"settlement_id,omitempty"` // 结算ID
	IsSettled       bool               `bson:"is_settled" json:"is_settled"`         // 是否已结算
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// WithdrawalRequest 提现申请
type WithdrawalRequest struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           string             `bson:"user_id" json:"user_id"`                     // 用户ID
	Amount           float64            `bson:"amount" json:"amount"`                       // 提现金额
	Fee              float64            `bson:"fee" json:"fee"`                             // 手续费
	ActualAmount     float64            `bson:"actual_amount" json:"actual_amount"`         // 实际到账金额
	Method           string             `bson:"method" json:"method"`                       // 提现方式：alipay, wechat, bank
	AccountInfo      WithdrawAccount    `bson:"account_info" json:"account_info"`          // 账户信息
	Status           string             `bson:"status" json:"status"`                       // 状态：pending, approved, rejected, completed, failed
	RejectReason     string             `bson:"reject_reason,omitempty" json:"reject_reason,omitempty"` // 拒绝原因
	ApprovedBy       string             `bson:"approved_by,omitempty" json:"approved_by,omitempty"`     // 审批人
	ApprovedAt       *time.Time         `bson:"approved_at,omitempty" json:"approved_at,omitempty"`     // 审批时间
	CompletedAt      *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`   // 完成时间
	TransactionID    string             `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"` // 交易流水号
	Note             string             `bson:"note,omitempty" json:"note,omitempty"`       // 备注
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// WithdrawAccount 提现账户信息
type WithdrawAccount struct {
	AccountType string `bson:"account_type" json:"account_type"` // 账户类型
	AccountName string `bson:"account_name" json:"account_name"` // 账户名
	AccountNo   string `bson:"account_no" json:"account_no"`     // 账号
	BankName    string `bson:"bank_name,omitempty" json:"bank_name,omitempty"` // 银行名称
	BranchName  string `bson:"branch_name,omitempty" json:"branch_name,omitempty"` // 开户支行
}

// Settlement 结算记录
type Settlement struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID          string             `bson:"author_id" json:"author_id"`                 // 作者ID
	AuthorNickname    string             `bson:"author_nickname" json:"author_nickname"`     // 作者昵称
	PeriodStart       time.Time          `bson:"period_start" json:"period_start"`           // 结算周期开始
	PeriodEnd         time.Time          `bson:"period_end" json:"period_end"`               // 结算周期结束
	TotalRevenue      float64            `bson:"total_revenue" json:"total_revenue"`         // 总收入
	PlatformFee       float64            `bson:"platform_fee" json:"platform_fee"`           // 平台费用
	ActualIncome      float64            `bson:"actual_income" json:"actual_income"`         // 实际收入
	TaxFee            float64            `bson:"tax_fee" json:"tax_fee"`                     // 税费
	FinalIncome       float64            `bson:"final_income" json:"final_income"`           // 最终收入
	EarningCount      int                `bson:"earning_count" json:"earning_count"`         // 收入记录数
	Status            string             `bson:"status" json:"status"`                       // 状态：pending, processing, completed, failed
	ProcessedAt       *time.Time         `bson:"processed_at,omitempty" json:"processed_at,omitempty"` // 处理时间
	TransactionID     string             `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"` // 交易流水号
	Note              string             `bson:"note,omitempty" json:"note,omitempty"`       // 备注
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

// RevenueDetail 收入明细
type RevenueDetail struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID         string             `bson:"author_id" json:"author_id"`               // 作者ID
	BookID           primitive.ObjectID `bson:"book_id" json:"book_id"`                   // 书籍ID
	BookTitle        string             `bson:"book_title" json:"book_title"`             // 书名
	Type             string             `bson:"type" json:"type"`                         // 收入类型
	TotalAmount      float64            `bson:"total_amount" json:"total_amount"`         // 总金额
	TotalIncome      float64            `bson:"total_income" json:"total_income"`         // 作者收入
	TransactionCount int                `bson:"transaction_count" json:"transaction_count"` // 交易笔数
	FirstEarningAt   *time.Time         `bson:"first_earning_at,omitempty" json:"first_earning_at,omitempty"` // 首次收入时间
	LastEarningAt    *time.Time         `bson:"last_earning_at,omitempty" json:"last_earning_at,omitempty"`   // 最后收入时间
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// RevenueStatistics 收入统计
type RevenueStatistics struct {
	AuthorID         string             `bson:"author_id" json:"author_id"`             // 作者ID
	Period           string             `bson:"period" json:"period"`                   // 统计周期：daily, monthly, yearly
	PeriodStart      time.Time          `bson:"period_start" json:"period_start"`       // 周期开始
	PeriodEnd        time.Time          `bson:"period_end" json:"period_end"`           // 周期结束
	TotalRevenue     float64            `bson:"total_revenue" json:"total_revenue"`     // 总收入
	ChapterIncome    float64            `bson:"chapter_income" json:"chapter_income"`   // 章节购买收入
	RewardIncome     float64            `bson:"reward_income" json:"reward_income"`     // 打赏收入
	VIPReadingIncome float64            `bson:"vip_reading_income" json:"vip_reading_income"` // VIP阅读收入
	TransactionCount int                `bson:"transaction_count" json:"transaction_count"` // 交易笔数
	ReaderCount      int                `bson:"reader_count" json:"reader_count"`       // 读者数
	BookCount        int                `bson:"book_count" json:"book_count"`           // 书籍数
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// TaxInfo 税务信息
type TaxInfo struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           string             `bson:"user_id" json:"user_id"`                   // 用户ID
	IDType           string             `bson:"id_type" json:"id_type"`                   // 证件类型：id_card, passport, other
	IDNumber         string             `bson:"id_number" json:"id_number"`               // 证件号
	Name             string             `bson:"name" json:"name"`                         // 真实姓名
	TaxType          string             `bson:"tax_type" json:"tax_type"`                 // 纳税人类型：individual, company
	TaxRate          float64            `bson:"tax_rate" json:"tax_rate"`                 // 税率
	IsVerified       bool               `bson:"is_verified" json:"is_verified"`           // 是否已认证
	VerifiedAt       *time.Time         `bson:"verified_at,omitempty" json:"verified_at,omitempty"` // 认证时间
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// EarningType 收入类型枚举
const (
	EarningTypeChapterPurchase = "chapter_purchase" // 章节购买
	EarningTypeReward          = "reward"          // 打赏
	EarningTypeVIPReading      = "vip_reading"     // VIP阅读
)

// WithdrawStatus 提现状态枚举
const (
	WithdrawStatusPending   = "pending"   // 待审核
	WithdrawStatusApproved  = "approved"  // 已批准
	WithdrawStatusRejected  = "rejected"  // 已拒绝
	WithdrawStatusCompleted = "completed" // 已完成
	WithdrawStatusFailed    = "failed"    // 失败
)

// SettlementStatus 结算状态枚举
const (
	SettlementStatusPending    = "pending"    // 待结算
	SettlementStatusProcessing = "processing" // 处理中
	SettlementStatusCompleted  = "completed"  // 已完成
	SettlementStatusFailed     = "failed"     // 失败
)

// 收入分成规则常量
const (
	// 章节购买分成
	ChapterPurchaseAuthorRate = 0.70 // 作者70%
	ChapterPurchasePlatformRate = 0.30 // 平台30%

	// 打赏分成
	RewardAuthorRate = 0.90  // 作者90%
	RewardPlatformRate = 0.10 // 平台10%

	// VIP阅读分成
	VIPReadingAuthorRate = 0.60  // 作者60%
	VIPReadingPlatformRate = 0.40 // 平台40%

	// 默认税率
	DefaultTaxRate = 0.00 // 默认0%（根据实际情况调整）
)

// WithdrawMethod 提现方式枚举
const (
	WithdrawMethodAlipay = "alipay" // 支付宝
	WithdrawMethodWechat = "wechat" // 微信
	WithdrawMethodBank   = "bank"   // 银行卡
)
