package auth

import "time"

// PermissionTemplate 权限模板模型
// 权限模板用于预定义一组权限，可以快速应用到角色
type PermissionTemplate struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`                   // 模板名称
	Code        string    `bson:"code" json:"code"`                   // 模板代码（唯一标识）
	Description string    `bson:"description" json:"description"`     // 模板描述
	Permissions []string  `bson:"permissions" json:"permissions"`     // 权限列表
	IsSystem    bool      `bson:"is_system" json:"is_system"`         // 是否系统模板（不可删除）
	Category    string    `bson:"category" json:"category"`           // 分类：reader, author, admin, custom
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	CreatedBy   string    `bson:"created_by,omitempty" json:"created_by,omitempty"` // 创建者ID
}

// 预定义模板代码
const (
	TemplateReader = "template_reader" // 读者模板
	TemplateAuthor = "template_author" // 作者模板
	TemplateAdmin  = "template_admin"  // 管理员模板
)

// 预定义模板分类
const (
	CategoryReader = "reader"
	CategoryAuthor = "author"
	CategoryAdmin  = "admin"
	CategoryCustom = "custom"
)

// Validate 验证模板数据
func (pt *PermissionTemplate) Validate() error {
	if pt.Name == "" {
		return ErrTemplateNameEmpty
	}
	if pt.Code == "" {
		return ErrTemplateCodeEmpty
	}
	if len(pt.Permissions) == 0 {
		return ErrTemplatePermissionsEmpty
	}
	return nil
}

// 模板验证错误
var (
	ErrTemplateNameEmpty        = ErrPermissionBase("template name cannot be empty")
	ErrTemplateCodeEmpty        = ErrPermissionBase("template code cannot be empty")
	ErrTemplatePermissionsEmpty = ErrPermissionBase("template permissions cannot be empty")
	ErrTemplateIsSystem         = ErrPermissionBase("cannot delete system template")
	ErrTemplateNotFound         = ErrPermissionBase("template not found")
	ErrTemplateCodeExists       = ErrPermissionBase("template code already exists")
)

// ErrPermissionBase 权限基础错误类型
type ErrPermissionBase string

func (e ErrPermissionBase) Error() string {
	return string(e)
}
