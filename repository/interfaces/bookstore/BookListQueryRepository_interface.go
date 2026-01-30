package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
)

// BookListQueryRepository 书籍列表查询接口 - 专注于书籍的只读查询操作
// 用于书城首页、分类页面、搜索结果等列表场景
//
// 职责：
// - 根据ID、分类、作者等条件查询书籍
// - 获取推荐、精选、热门等列表
// - 健康检查
//
// 方法总数：14个
type BookListQueryRepository interface {
	// 继承基础CRUD接口的查询方法
	base.CRUDRepository[*bookstore2.Book, string]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 列表查询方法 - 用于书城展示
	GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstore2.Book, error)
	GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore2.Book, error)
	GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstore2.Book, error)
	GetByStatus(ctx context.Context, status bookstore2.BookStatus, limit, offset int) ([]*bookstore2.Book, error)

	// 推荐列表 - 用于首页展示
	GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error)
	GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error)
	GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error)
	GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error)
	GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error)
}
