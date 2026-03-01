package admin

import (
	"context"
	"Qingyu_backend/models/bookstore"
)

// CategoryAdminRepository 分类管理仓储接口（管理员专用）
// 复用现有的 CategoryRepository，扩展管理员特有方法
type CategoryAdminRepository interface {
	// 嵌入基础分类仓储（复用现有方法）
	bookstore.CategoryRepository

	// UpdateBookCount 更新分类下的作品数量
	UpdateBookCount(ctx context.Context, categoryID string, count int64) error

	// BatchUpdateStatus 批量更新分类启用状态
	BatchUpdateStatus(ctx context.Context, categoryIDs []string, isActive bool) error

	// GetDescendantIDs 获取所有子孙分类ID（用于级联操作）
	GetDescendantIDs(ctx context.Context, categoryID string) ([]string, error)

	// HasChildren 检查是否有子分类（用于删除前校验）
	HasChildren(ctx context.Context, categoryID string) (bool, error)

	// NameExistsAtLevel 检查同级分类名称是否存在
	NameExistsAtLevel(ctx context.Context, parentID *string, name string, excludeID string) (bool, error)
}
