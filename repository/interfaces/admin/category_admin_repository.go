package admin

import (
	"Qingyu_backend/models/bookstore"
	"context"
)

// CategoryAdminRepository 分类管理仓储接口（管理员专用）
// 扩展基础 CategoryRepository，添加管理员特有方法
type CategoryAdminRepository interface {
	// ============ 基础CRUD方法（复用CategoryRepository） ============
	// 这些方法由现有的 CategoryRepository 提供，这里不重复定义
	// 实现时需要嵌入 CategoryRepository 实例

	// Create 创建分类
	Create(ctx context.Context, category *bookstore.Category) error

	// GetByID 根据ID获取分类
	GetByID(ctx context.Context, id string) (*bookstore.Category, error)

	// Update 更新分类
	Update(ctx context.Context, category *bookstore.Category) error

	// Delete 删除分类
	Delete(ctx context.Context, id string) error

	// List 获取分类列表
	List(ctx context.Context, filter interface{}, opts ...interface{}) ([]*bookstore.Category, error)

	// GetTree 获取分类树
	GetTree(ctx context.Context) ([]*bookstore.Category, error)

	// ============ 管理员特有方法 ============

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
