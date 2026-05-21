package ai

import (
	"context"

	aiModels "Qingyu_backend/models/ai"
)

// QuotaPolicyRepository 配额策略Repository接口
type QuotaPolicyRepository interface {
	// Create 创建配额策略
	Create(ctx context.Context, policy *aiModels.QuotaPolicy) error

	// GetByID 根据ID获取策略
	GetByID(ctx context.Context, id string) (*aiModels.QuotaPolicy, error)

	// GetByRoleAndLevel 根据角色和等级获取策略
	GetByRoleAndLevel(ctx context.Context, role aiModels.UserRole, level aiModels.MembershipLevel) (*aiModels.QuotaPolicy, error)

	// List 获取策略列表
	List(ctx context.Context, role string, status string, page, limit int) ([]*aiModels.QuotaPolicy, int64, error)

	// Update 更新策略
	Update(ctx context.Context, policy *aiModels.QuotaPolicy) error

	// Delete 删除策略（软删除）
	Delete(ctx context.Context, id string) error

	// Health 健康检查
	Health(ctx context.Context) error
}
