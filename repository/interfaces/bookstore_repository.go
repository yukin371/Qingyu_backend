package interfaces

import (
	"context"
	"time"
	"Qingyu_backend/models/reading/bookstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookRepository 书籍仓储接口
type BookRepository interface {
	// 继承基础CRUD接口
	CRUDRepository[*bookstore.Book, primitive.ObjectID]

	// 书籍特定查询方法
	GetByTitle(ctx context.Context, title string) (*bookstore.Book, error)
	GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.Book, error)
	GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*bookstore.Book, error)
	GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore.Book, error)
	
	// 搜索方法
	Search(ctx context.Context, keyword string, filter *bookstore.BookFilter) ([]*bookstore.Book, error)
	SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error)
	
	// 统计方法
	CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error)
	CountByAuthor(ctx context.Context, author string) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	GetStats(ctx context.Context) (*bookstore.BookStats, error)
	
	// 更新统计数据
	IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error
	
	// 批量操作
	BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status string) error
	BatchUpdateCategory(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []primitive.ObjectID) error
	
	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	// 继承基础CRUD接口
	CRUDRepository[*bookstore.Category, primitive.ObjectID]

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
	CRUDRepository[*bookstore.Banner, primitive.ObjectID]

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

// BookstoreRepositoryFactory 书城仓储工厂接口
type BookstoreRepositoryFactory interface {
	GetBookRepository() BookRepository
	GetCategoryRepository() CategoryRepository
	GetBannerRepository() BannerRepository
	Health(ctx context.Context) error
	Close() error
}