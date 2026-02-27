package auth

import (
	authModel "Qingyu_backend/models/auth"
	"context"
)

// PermissionTemplateRepository 权限模板仓储接口
type PermissionTemplateRepository interface {
	// ==================== 模板管理 ====================

	// CreateTemplate 创建模板
	CreateTemplate(ctx context.Context, template *authModel.PermissionTemplate) error

	// GetTemplateByID 根据ID获取模板
	GetTemplateByID(ctx context.Context, templateID string) (*authModel.PermissionTemplate, error)

	// GetTemplateByCode 根据代码获取模板
	GetTemplateByCode(ctx context.Context, code string) (*authModel.PermissionTemplate, error)

	// UpdateTemplate 更新模板
	UpdateTemplate(ctx context.Context, templateID string, updates map[string]interface{}) error

	// DeleteTemplate 删除模板
	DeleteTemplate(ctx context.Context, templateID string) error

	// ListTemplates 列出所有模板
	ListTemplates(ctx context.Context) ([]*authModel.PermissionTemplate, error)

	// ListTemplatesByCategory 根据分类列出模板
	ListTemplatesByCategory(ctx context.Context, category string) ([]*authModel.PermissionTemplate, error)

	// ==================== 模板应用 ====================

	// ApplyTemplateToRole 将模板应用到角色
	ApplyTemplateToRole(ctx context.Context, templateID, roleID string) error

	// GetSystemTemplates 获取所有系统模板
	GetSystemTemplates(ctx context.Context) ([]*authModel.PermissionTemplate, error)

	// InitializeSystemTemplates 初始化系统预设模板
	InitializeSystemTemplates(ctx context.Context) error

	// Health 健康检查
	Health(ctx context.Context) error
}
