package auth

import "time"

// Role 角色模型
type Role struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`               // 角色名称：reader, author, admin
	Description string    `json:"description" bson:"description"` // 角色描述
	Permissions []string  `json:"permissions" bson:"permissions"` // 权限列表
	IsSystem    bool      `json:"is_system" bson:"is_system"`     // 是否系统角色（不可删除）
	IsDefault   bool      `json:"is_default" bson:"is_default"`   // 是否默认角色（新用户默认分配）
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Permission 权限定义（可选，可以只用字符串）
type Permission struct {
	Code        string    `json:"code" bson:"code"`               // 权限代码：user.read, book.write
	Name        string    `json:"name" bson:"name"`               // 权限名称
	Description string    `json:"description" bson:"description"` // 权限描述
	Resource    string    `json:"resource" bson:"resource"`       // 资源类型：user, book, wallet
	Action      string    `json:"action" bson:"action"`           // 操作类型：read, write, delete
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// UserRole 用户角色关联（存储在users集合中的roles字段，或单独的集合）
type UserRole struct {
	UserID     string    `json:"user_id" bson:"user_id"`
	RoleID     string    `json:"role_id" bson:"role_id"`
	AssignedAt time.Time `json:"assigned_at" bson:"assigned_at"`
	AssignedBy string    `json:"assigned_by,omitempty" bson:"assigned_by,omitempty"` // 分配者ID
}

// 预定义角色
const (
	RoleReader = "reader" // 读者
	RoleAuthor = "author" // 作者
	RoleAdmin  = "admin"  // 管理员
)

// 预定义权限
const (
	// 用户权限
	PermUserRead   = "user.read"
	PermUserWrite  = "user.write"
	PermUserDelete = "user.delete"

	// 书籍权限
	PermBookRead   = "book.read"
	PermBookWrite  = "book.write"
	PermBookDelete = "book.delete"
	PermBookReview = "book.review"

	// 钱包权限
	PermWalletRead     = "wallet.read"
	PermWalletRecharge = "wallet.recharge"
	PermWalletWithdraw = "wallet.withdraw"
	PermWalletReview   = "wallet.review"

	// 管理权限
	PermAdminAccess = "admin.access"
	PermAdminReview = "admin.review"
	PermAdminManage = "admin.manage"

	// 文档权限
	PermDocumentRead    = "document.read"
	PermDocumentWrite   = "document.write"
	PermDocumentDelete  = "document.delete"
	PermDocumentPublish = "document.publish"

	// 评论权限
	PermCommentRead   = "comment.read"
	PermCommentWrite  = "comment.write"
	PermCommentDelete = "comment.delete"
)
