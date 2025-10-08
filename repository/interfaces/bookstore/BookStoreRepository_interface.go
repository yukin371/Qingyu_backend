package bookstore

import (
	"Qingyu_backend/models/reading/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookRepository 书籍仓储接口 - 专注于书籍列表展示和基础管理
// 用于书城首页、分类页面、搜索结果等列表场景
type BookRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.Book, primitive.ObjectID]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 列表查询方法 - 用于书城展示
	GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error)
	GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.Book, error)
	GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error)
	GetByStatus(ctx context.Context, status bookstore.BookStatus, limit, offset int) ([]*bookstore.Book, error)
	GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore.Book, error)

	// 搜索方法 - 用于书城搜索
	Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Book, error)
	SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error)

	// 统计方法 - 用于列表计数
	CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error)
	CountByAuthor(ctx context.Context, author string) (int64, error)
	CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error)
	CountByFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error)

	// 批量操作 - 用于管理
	BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error
	BatchUpdateCategory(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []primitive.ObjectID) error
	BatchUpdateRecommended(ctx context.Context, bookIDs []primitive.ObjectID, isRecommended bool) error
	BatchUpdateFeatured(ctx context.Context, bookIDs []primitive.ObjectID, isFeatured bool) error

	// 统计和计数操作
	GetStats(ctx context.Context) (*bookstore.BookStats, error)
	IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.Category, primitive.ObjectID]

	// 分类特定查询方法
	GetByName(ctx context.Context, name string) (*bookstore.Category, error)
	GetByParent(ctx context.Context, parentID primitive.ObjectID, limit, offset int) ([]*bookstore.Category, error)
	GetByLevel(ctx context.Context, level int, limit, offset int) ([]*bookstore.Category, error)
	GetRootCategories(ctx context.Context) ([]*bookstore.Category, error)
	GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error)

	// 统计方法
	CountByParent(ctx context.Context, parentID primitive.ObjectID) (int64, error)
	UpdateBookCount(ctx context.Context, categoryID primitive.ObjectID, count int64) error

	// 层级操作
	GetChildren(ctx context.Context, parentID primitive.ObjectID) ([]*bookstore.Category, error)
	GetAncestors(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstore.Category, error)
	GetDescendants(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstore.Category, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, categoryIDs []primitive.ObjectID, isActive bool) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BannerRepository Banner仓储接口
type BannerRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.Banner, primitive.ObjectID]

	// Banner特定查询方法
	GetActive(ctx context.Context, limit, offset int) ([]*bookstore.Banner, error)
	GetByTargetType(ctx context.Context, targetType string, limit, offset int) ([]*bookstore.Banner, error)
	GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*bookstore.Banner, error)

	// 统计方法
	IncrementClickCount(ctx context.Context, bannerID primitive.ObjectID) error
	GetClickStats(ctx context.Context, bannerID primitive.ObjectID) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, bannerIDs []primitive.ObjectID, isActive bool) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
