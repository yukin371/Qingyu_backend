package auth

import (
	"errors"
	"strings"
	"time"
)

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
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Code        string    `json:"code" bson:"code"`                     // 权限代码：user.read, book.write
	Name        string    `json:"name" bson:"name"`                     // 权限名称
	Description string    `json:"description" bson:"description"`       // 权限描述
	Resource    string    `json:"resource" bson:"resource"`             // 资源类型：user, book, wallet
	Action      string    `json:"action" bson:"action"`                 // 操作类型：read, write, delete
	Effect      string    `json:"effect" bson:"effect"`                 // 权限效果：allow 或 deny
	Priority    int       `json:"priority" bson:"priority"`             // 优先级，数值越大优先级越高
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

// 权限效果常量
const (
	EffectAllow = "allow"
	EffectDeny  = "deny"
)

// 权限验证错误
var (
	ErrPermissionCodeEmpty     = errors.New("permission code cannot be empty")
	ErrPermissionEffectInvalid = errors.New("permission effect must be 'allow' or 'deny'")
	ErrPermissionPriorityNeg   = errors.New("permission priority cannot be negative")
)

// ResourcePermission 资源级权限
type ResourcePermission struct {
	PermissionID string   `json:"permissionId" bson:"permissionId"` // 权限ID
	ResourceIDs  []string `json:"resourceIds" bson:"resourceIds"`   // 资源ID列表，空表示所有资源
	Effect       string   `json:"effect" bson:"effect"`             // 权限效果：allow 或 deny
}

// HasPermission 检查权限是否匹配
func (p *Permission) HasPermission(requiredCode string) bool {
	// 精确匹配或通配符匹配
	if p.Code == requiredCode || p.Code == "*" {
		return p.Effect == EffectAllow
	}
	// 支持前缀匹配: user.* 匹配 user.read, user.write
	if strings.HasSuffix(p.Code, ".*") {
		prefix := strings.TrimSuffix(p.Code, ".*")
		return strings.HasPrefix(requiredCode, prefix+".") && p.Effect == EffectAllow
	}
	return false
}

// GetHigherPriority 获取优先级更高的权限
func (p *Permission) GetHigherPriority(other *Permission) *Permission {
	if other == nil {
		return p
	}
	// deny 优先级高于 allow
	if p.Effect == EffectDeny && other.Effect == EffectAllow {
		return p
	}
	if other.Effect == EffectDeny && p.Effect == EffectAllow {
		return other
	}
	// 同等效果，priority数值大的优先
	if p.Priority >= other.Priority {
		return p
	}
	return other
}

// Validate 验证权限数据
func (p *Permission) Validate() error {
	if p.Code == "" {
		return ErrPermissionCodeEmpty
	}
	if p.Effect != EffectAllow && p.Effect != EffectDeny {
		return ErrPermissionEffectInvalid
	}
	if p.Priority < 0 {
		return ErrPermissionPriorityNeg
	}
	return nil
}

// IsResourceAllowed 检查资源是否被允许访问
func (rp *ResourcePermission) IsResourceAllowed(resourceID string) bool {
	// 空或nil ResourceIDs 表示所有资源
	if len(rp.ResourceIDs) == 0 {
		return rp.Effect == EffectAllow
	}
	// 检查资源ID是否在列表中
	for _, rid := range rp.ResourceIDs {
		if rid == "*" || rid == resourceID {
			return rp.Effect == EffectAllow
		}
	}
	return false
}
