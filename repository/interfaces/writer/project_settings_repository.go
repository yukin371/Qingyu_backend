package writer

import (
	"context"

	"Qingyu_backend/models/writer"
)

// ProjectSettingsRepository 项目设置仓库接口
type ProjectSettingsRepository interface {
	// Create 创建项目设置
	Create(ctx context.Context, settings *writer.ProjectSettings) error

	// FindByProjectID 根据项目ID查找设置
	FindByProjectID(ctx context.Context, projectID string) (*writer.ProjectSettings, error)

	// Update 更新项目设置
	Update(ctx context.Context, projectID string, settings *writer.ProjectSettings) error

	// AddCharacterRole 添加自定义角色类型
	AddCharacterRole(ctx context.Context, projectID string, role *writer.CharacterRole) error

	// UpdateCharacterRole 更新角色类型
	UpdateCharacterRole(ctx context.Context, projectID, roleID string, role *writer.CharacterRole) error

	// DeleteCharacterRole 删除角色类型
	DeleteCharacterRole(ctx context.Context, projectID, roleID string) error

	// GetDefaultRoles 获取项目的角色类型列表（包含默认和自定义）
	GetDefaultRoles(ctx context.Context, projectID string) ([]writer.CharacterRole, error)
}
