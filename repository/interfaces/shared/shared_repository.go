package shared

import (
	authModel "Qingyu_backend/models/auth"
	recommendationModel "Qingyu_backend/models/recommendation"
	storageModel "Qingyu_backend/models/storage"
	adminModel "Qingyu_backend/models/users"
	walletModel "Qingyu_backend/models/wallet"
	"context"
	"time"
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

	// Health 健康检查
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

	// Health 健康检查
	Health(ctx context.Context) error
}

// ============ Storage Repository ============

// StorageRepository 文件存储相关Repository
type StorageRepository interface {
	// 文件元数据管理
	CreateFile(ctx context.Context, file *storageModel.FileInfo) error
	GetFile(ctx context.Context, fileID string) (*storageModel.FileInfo, error)
	GetFileByMD5(ctx context.Context, md5Hash string) (*storageModel.FileInfo, error)
	UpdateFile(ctx context.Context, fileID string, updates map[string]interface{}) error
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, filter *FileFilter) ([]*storageModel.FileInfo, int64, error)
	CountFiles(ctx context.Context, filter *FileFilter) (int64, error)

	// 权限管理
	GrantAccess(ctx context.Context, fileID, userID string) error
	RevokeAccess(ctx context.Context, fileID, userID string) error
	CheckAccess(ctx context.Context, fileID, userID string) (bool, error)

	// 分片上传管理
	CreateMultipartUpload(ctx context.Context, upload *storageModel.MultipartUpload) error
	GetMultipartUpload(ctx context.Context, uploadID string) (*storageModel.MultipartUpload, error)
	UpdateMultipartUpload(ctx context.Context, uploadID string, updates map[string]interface{}) error
	CompleteMultipartUpload(ctx context.Context, uploadID string) error
	AbortMultipartUpload(ctx context.Context, uploadID string) error
	ListMultipartUploads(ctx context.Context, userID string, status string) ([]*storageModel.MultipartUpload, error)

	// 统计功能
	IncrementDownloadCount(ctx context.Context, fileID string) error

	// Health 健康检查
	Health(ctx context.Context) error
}

// FileFilter 文件过滤器
type FileFilter struct {
	UserID    string
	Category  string
	FileType  string
	Status    string
	IsPublic  *bool
	Tags      []string
	Keyword   string
	StartDate *time.Time
	EndDate   *time.Time
	MinSize   *int64
	MaxSize   *int64
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Limit     int64
	Offset    int64
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

	// Health 健康检查
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
