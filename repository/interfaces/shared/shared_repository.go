package shared

import (
	// 重新导出新模块的接口，保持向后兼容
	"Qingyu_backend/repository/interfaces/admin"
	"Qingyu_backend/repository/interfaces/auth"
	"Qingyu_backend/repository/interfaces/finance"
	"Qingyu_backend/repository/interfaces/messaging"
	"Qingyu_backend/repository/interfaces/recommendation"
	"Qingyu_backend/repository/interfaces/storage"
)

// ============ 向后兼容的类型别名 ============
// 以下类型别名用于保持向后兼容性，新代码应直接使用对应模块的接口

// AuthRepository 认证相关Repository (别名，指向 auth.RoleRepository)
// Deprecated: 使用 auth.RoleRepository 替代
type AuthRepository = auth.RoleRepository

// WalletRepository 钱包相关Repository (别名，指向 finance.WalletRepository)
// Deprecated: 使用 finance.WalletRepository 替代
type WalletRepository = finance.WalletRepository

// TransactionFilter 交易过滤器 (别名，指向 finance.TransactionFilter)
// Deprecated: 使用 finance.TransactionFilter 替代
type TransactionFilter = finance.TransactionFilter

// WithdrawFilter 提现过滤器 (别名，指向 finance.WithdrawFilter)
// Deprecated: 使用 finance.WithdrawFilter 替代
type WithdrawFilter = finance.WithdrawFilter

// RecommendationRepository 推荐相关Repository (别名，指向 recommendation.RecommendationRepository)
// Deprecated: 使用 recommendation.RecommendationRepository 替代
type RecommendationRepository = recommendation.RecommendationRepository

// StorageRepository 文件存储相关Repository (别名，指向 storage.StorageRepository)
// Deprecated: 使用 storage.StorageRepository 替代
type StorageRepository = storage.StorageRepository

// FileFilter 文件过滤器 (别名，指向 storage.FileFilter)
// Deprecated: 使用 storage.FileFilter 替代
type FileFilter = storage.FileFilter

// AdminRepository 管理后台相关Repository (别名，指向 admin.AuditRepository)
// Deprecated: 使用 admin.AuditRepository 和 admin.AdminLogRepository 替代
type AdminRepository = admin.AuditRepository

// AuditFilter 审核记录过滤器 (别名，指向 admin.AuditFilter)
// Deprecated: 使用 admin.AuditFilter 替代
type AuditFilter = admin.AuditFilter

// AdminLogFilter 管理员日志过滤器 (别名，指向 admin.AdminLogFilter)
// Deprecated: 使用 admin.AdminLogFilter 替代
type AdminLogFilter = admin.AdminLogFilter

// MessageRepository 消息相关Repository (别名，指向 messaging.MessageRepository)
// Deprecated: 使用 messaging.MessageRepository 替代
type MessageRepository = messaging.MessageRepository

// MessageFilter 消息过滤器 (别名，指向 messaging.MessageFilter)
// Deprecated: 使用 messaging.MessageFilter 替代
type MessageFilter = messaging.MessageFilter

// NotificationFilter 通知过滤器 (别名，指向 messaging.NotificationFilter)
// Deprecated: 使用 messaging.NotificationFilter 替代
type NotificationFilter = messaging.NotificationFilter

