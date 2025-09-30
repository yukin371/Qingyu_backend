package shared

import (
	"context"
	"time"

	adminModel "Qingyu_backend/models/shared/admin"
	authModel "Qingyu_backend/models/shared/auth"
	recommendationModel "Qingyu_backend/models/shared/recommendation"
	storageModel "Qingyu_backend/models/shared/storage"
	walletModel "Qingyu_backend/models/shared/wallet"
)

// ============ Auth Repository ============

// AuthRepository 认证相关Repository
type AuthRepository interface {
	// 角色管理
	CreateRole(ctx context.Context, role *authModel.Role) error
	GetRole(ctx context.Context, roleID string) (*authModel.Role, error)
	GetRoleByName(ctx context.Context, name string) (*authModel.Role, error)
	UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error
	DeleteRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context) ([]*authModel.Role, error)

	// 用户角色关联
	AssignUserRole(ctx context.Context, userID, roleID string) error
	RemoveUserRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error)
	HasUserRole(ctx context.Context, userID, roleID string) (bool, error)

	// 权限查询
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)

	// 健康检查
	Health(ctx context.Context) error
}

// ============ Wallet Repository ============

// WalletRepository 钱包相关Repository
type WalletRepository interface {
	// 钱包管理
	CreateWallet(ctx context.Context, wallet *walletModel.Wallet) error
	GetWallet(ctx context.Context, userID string) (*walletModel.Wallet, error)
	UpdateWallet(ctx context.Context, userID string, updates map[string]interface{}) error
	UpdateBalance(ctx context.Context, userID string, amount float64) error

	// 交易记录
	CreateTransaction(ctx context.Context, transaction *walletModel.Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (*walletModel.Transaction, error)
	ListTransactions(ctx context.Context, filter *TransactionFilter) ([]*walletModel.Transaction, error)
	CountTransactions(ctx context.Context, filter *TransactionFilter) (int64, error)

	// 提现申请
	CreateWithdrawRequest(ctx context.Context, request *walletModel.WithdrawRequest) error
	GetWithdrawRequest(ctx context.Context, withdrawID string) (*walletModel.WithdrawRequest, error)
	UpdateWithdrawRequest(ctx context.Context, withdrawID string, updates map[string]interface{}) error
	ListWithdrawRequests(ctx context.Context, filter *WithdrawFilter) ([]*walletModel.WithdrawRequest, error)
	CountWithdrawRequests(ctx context.Context, filter *WithdrawFilter) (int64, error)

	// 健康检查
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

// ============ Recommendation Repository ============

// RecommendationRepository 推荐相关Repository
type RecommendationRepository interface {
	// 用户行为记录
	RecordBehavior(ctx context.Context, behavior *recommendationModel.UserBehavior) error
	GetUserBehaviors(ctx context.Context, userID string, limit int) ([]*recommendationModel.UserBehavior, error)
	GetItemBehaviors(ctx context.Context, itemID string, limit int) ([]*recommendationModel.UserBehavior, error)

	// 推荐结果存储（可选，可以完全用缓存）
	// SaveRecommendations(ctx context.Context, userID string, items []*recommendationModel.RecommendedItem) error
	// GetRecommendations(ctx context.Context, userID string) ([]*recommendationModel.RecommendedItem, error)

	// 健康检查
	Health(ctx context.Context) error
}

// ============ Storage Repository ============

// StorageRepository 文件存储相关Repository
type StorageRepository interface {
	// 文件元数据管理
	CreateFile(ctx context.Context, file *storageModel.FileInfo) error
	GetFile(ctx context.Context, fileID string) (*storageModel.FileInfo, error)
	UpdateFile(ctx context.Context, fileID string, updates map[string]interface{}) error
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, filter *FileFilter) ([]*storageModel.FileInfo, error)
	CountFiles(ctx context.Context, filter *FileFilter) (int64, error)

	// 权限管理（可选，也可以通过FileInfo.IsPublic处理）
	GrantAccess(ctx context.Context, fileID, userID string) error
	RevokeAccess(ctx context.Context, fileID, userID string) error
	CheckAccess(ctx context.Context, fileID, userID string) (bool, error)

	// 健康检查
	Health(ctx context.Context) error
}

// FileFilter 文件过滤器
type FileFilter struct {
	UserID   string
	Category string
	IsPublic *bool
	Limit    int64
	Offset   int64
}

// ============ Admin Repository ============

// AdminRepository 管理后台相关Repository
type AdminRepository interface {
	// 审核记录
	CreateAuditRecord(ctx context.Context, record *adminModel.AuditRecord) error
	GetAuditRecord(ctx context.Context, contentID, contentType string) (*adminModel.AuditRecord, error)
	UpdateAuditRecord(ctx context.Context, recordID string, updates map[string]interface{}) error
	ListAuditRecords(ctx context.Context, filter *AuditFilter) ([]*adminModel.AuditRecord, error)

	// 操作日志
	CreateAdminLog(ctx context.Context, log *adminModel.AdminLog) error
	GetAdminLog(ctx context.Context, logID string) (*adminModel.AdminLog, error)
	ListAdminLogs(ctx context.Context, filter *AdminLogFilter) ([]*adminModel.AdminLog, error)
	CountAdminLogs(ctx context.Context, filter *AdminLogFilter) (int64, error)

	// 健康检查
	Health(ctx context.Context) error
}

// AuditFilter 审核记录过滤器
type AuditFilter struct {
	ContentType string
	Status      string
	ReviewerID  string
	StartDate   time.Time
	EndDate     time.Time
	Limit       int64
	Offset      int64
}

// AdminLogFilter 管理员日志过滤器
type AdminLogFilter struct {
	AdminID   string
	Operation string
	StartDate time.Time
	EndDate   time.Time
	Limit     int64
	Offset    int64
}
