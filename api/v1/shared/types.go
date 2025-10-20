package shared

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WalletBalanceResponse 钱包余额响应
type WalletBalanceResponse struct {
	UserID           primitive.ObjectID `json:"userId"`
	Balance          float64            `json:"balance"`
	FrozenBalance    float64            `json:"frozenBalance"`
	AvailableBalance float64            `json:"availableBalance"`
	TotalIncome      float64            `json:"totalIncome"`
	TotalExpense     float64            `json:"totalExpense"`
	Currency         string             `json:"currency"`
	UpdatedAt        time.Time          `json:"updatedAt"`
}

// TransactionResponse 交易响应
type TransactionResponse struct {
	ID          primitive.ObjectID `json:"id"`
	UserID      primitive.ObjectID `json:"userId"`
	Type        string             `json:"type"` // deposit, withdraw, purchase, reward
	Amount      float64            `json:"amount"`
	Balance     float64            `json:"balance"` // 交易后余额
	Status      string             `json:"status"`  // pending, completed, failed
	Description string             `json:"description"`
	ReferenceID string             `json:"referenceId,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
}

// WithdrawResponse 提现响应
type WithdrawResponse struct {
	ID           primitive.ObjectID `json:"id"`
	UserID       primitive.ObjectID `json:"userId"`
	Amount       float64            `json:"amount"`
	Fee          float64            `json:"fee"`
	ActualAmount float64            `json:"actualAmount"` // 实际到账金额
	Status       string             `json:"status"`       // pending, approved, rejected, completed
	BankAccount  string             `json:"bankAccount"`
	AppliedAt    time.Time          `json:"appliedAt"`
	ProcessedAt  *time.Time         `json:"processedAt,omitempty"`
	ProcessedBy  string             `json:"processedBy,omitempty"`
	RejectReason string             `json:"rejectReason,omitempty"`
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	ID            primitive.ObjectID `json:"id"`
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Avatar        string             `json:"avatar"`
	Role          string             `json:"role"`
	VIPLevel      int                `json:"vipLevel"`
	VIPExpireAt   *time.Time         `json:"vipExpireAt,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
	EmailVerified bool               `json:"emailVerified"`
	PhoneVerified bool               `json:"phoneVerified"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string              `json:"token"`
	TokenType string              `json:"tokenType"`
	ExpiresIn int64               `json:"expiresIn"` // 秒
	User      UserProfileResponse `json:"user"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID     string    `json:"fileId"`
	Filename   string    `json:"filename"`
	URL        string    `json:"url"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mimeType"`
	UploadedAt time.Time `json:"uploadedAt"`
}

// StorageQuotaResponse 存储配额响应
type StorageQuotaResponse struct {
	UserID         primitive.ObjectID `json:"userId"`
	TotalQuota     int64              `json:"totalQuota"` // 字节
	UsedQuota      int64              `json:"usedQuota"`
	AvailableQuota int64              `json:"availableQuota"`
	FileCount      int                `json:"fileCount"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}

// AdminStatsResponse 管理员统计响应
type AdminStatsResponse struct {
	TotalUsers       int64   `json:"totalUsers"`
	ActiveUsers      int64   `json:"activeUsers"`
	TotalProjects    int64   `json:"totalProjects"`
	TotalBooks       int64   `json:"totalBooks"`
	TotalRevenue     float64 `json:"totalRevenue"`
	PendingAudits    int64   `json:"pendingAudits"`
	PendingWithdraws int64   `json:"pendingWithdraws"`
}
