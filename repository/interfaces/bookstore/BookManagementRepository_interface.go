package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
)

// BookManagementRepository 书籍管理接口 - 专注于书籍的批量管理和元数据操作
// 用于后台管理系统、批量操作场景
//
// 职责：
// - 批量更新书籍状态、分类、推荐等
// - 获取元数据（年份、标签等）
//
// 方法总数：6个
type BookManagementRepository interface {
	// BatchUpdateStatus 批量更新书籍状态
	BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstore2.BookStatus) error

	// BatchUpdateCategory 批量更新书籍分类
	BatchUpdateCategory(ctx context.Context, bookIDs []string, categoryIDs []string) error

	// BatchUpdateRecommended 批量更新推荐状态
	BatchUpdateRecommended(ctx context.Context, bookIDs []string, isRecommended bool) error

	// BatchUpdateFeatured 批量更新精选状态
	BatchUpdateFeatured(ctx context.Context, bookIDs []string, isFeatured bool) error

	// GetYears 获取所有书籍的发布年份列表（去重，倒序）
	// 用于筛选项
	GetYears(ctx context.Context) ([]int, error)

	// GetTags 获取所有标签列表（去重，排序）
	// 如果提供了 categoryID，则只返回该分类下的书籍标签
	GetTags(ctx context.Context, categoryID *string) ([]string, error)
}
