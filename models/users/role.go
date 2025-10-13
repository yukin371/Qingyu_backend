package users

import "time"

// 角色常量
const (
	RoleUser   = "user"   // 普通用户
	RoleAuthor = "author" // 作者
	RoleAdmin  = "admin"  // 管理员
)

// 权限常量
const (
	// 用户权限
	PermissionUserRead   = "user:read"
	PermissionUserWrite  = "user:write"
	PermissionUserDelete = "user:delete"

	// 文档权限
	PermissionDocumentRead    = "document:read"
	PermissionDocumentWrite   = "document:write"
	PermissionDocumentDelete  = "document:delete"
	PermissionDocumentPublish = "document:publish"

	// 书籍权限
	PermissionBookRead    = "book:read"
	PermissionBookWrite   = "book:write"
	PermissionBookDelete  = "book:delete"
	PermissionBookPublish = "book:publish"

	// 评论权限
	PermissionCommentRead   = "comment:read"
	PermissionCommentWrite  = "comment:write"
	PermissionCommentDelete = "comment:delete"

	// 管理权限
	PermissionAdminAccess = "admin:access"
	PermissionAdminUsers  = "admin:users"
	PermissionAdminBooks  = "admin:books"
	PermissionAdminReview = "admin:review"
)

// Role 角色模型
type Role struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name" validate:"required,min=2,max=50"`
	Description string    `bson:"description,omitempty" json:"description,omitempty" validate:"max=200"`
	IsDefault   bool      `bson:"is_default" json:"isDefault"`
	Permissions []string  `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

// HasPermission 检查角色是否拥有指定权限
func (r *Role) HasPermission(permission string) bool {
	for _, p := range r.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// AddPermission 添加权限
func (r *Role) AddPermission(permission string) {
	if !r.HasPermission(permission) {
		r.Permissions = append(r.Permissions, permission)
	}
}

// RemovePermission 移除权限
func (r *Role) RemovePermission(permission string) {
	for i, p := range r.Permissions {
		if p == permission {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			break
		}
	}
}

// GetDefaultPermissions 获取角色的默认权限集
func GetDefaultPermissions(roleName string) []string {
	switch roleName {
	case RoleAdmin:
		return []string{
			PermissionUserRead, PermissionUserWrite, PermissionUserDelete,
			PermissionDocumentRead, PermissionDocumentWrite, PermissionDocumentDelete, PermissionDocumentPublish,
			PermissionBookRead, PermissionBookWrite, PermissionBookDelete, PermissionBookPublish,
			PermissionCommentRead, PermissionCommentWrite, PermissionCommentDelete,
			PermissionAdminAccess, PermissionAdminUsers, PermissionAdminBooks, PermissionAdminReview,
		}
	case RoleAuthor:
		return []string{
			PermissionUserRead, PermissionUserWrite,
			PermissionDocumentRead, PermissionDocumentWrite, PermissionDocumentDelete, PermissionDocumentPublish,
			PermissionBookRead, PermissionBookWrite, PermissionBookPublish,
			PermissionCommentRead, PermissionCommentWrite,
		}
	case RoleUser:
		return []string{
			PermissionUserRead, PermissionUserWrite,
			PermissionDocumentRead,
			PermissionBookRead,
			PermissionCommentRead, PermissionCommentWrite,
		}
	default:
		return []string{}
	}
}
