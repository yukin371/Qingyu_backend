package repository

import (
	"context"

	"github.com/qingteng-qa/internal/domain"
)

// BookRepository 书籍Repository接口
// 定义书籍数据访问的抽象接口,遵循单一职责原则
// 接口方法数: 10个 (符合≤15个的要求)
type BookRepository interface {
	// Create 创建书籍
	// 返回: error - 创建失败时返回错误
	Create(ctx context.Context, book *domain.Book) error

	// GetByID 根据ID获取书籍
	// 返回: (*domain.Book, error) - 书籍不存在时返回ErrCodeBookNotFound
	GetByID(ctx context.Context, id string) (*domain.Book, error)

	// Update 更新书籍
	// 返回: error - 更新失败时返回错误,书籍不存在时返回ErrCodeBookNotFound
	Update(ctx context.Context, book *domain.Book) error

	// Delete 删除书籍 (软删除,更新状态为deleted)
	// 返回: error - 删除失败时返回错误,书籍不存在时返回ErrCodeBookNotFound
	Delete(ctx context.Context, id string) error

	// List 查询书籍列表
	// 参数:
	//   - filter: 查询过滤条件
	//   - page: 页码 (从1开始)
	//   - size: 每页数量
	// 返回: ([]*domain.Book, int64, error) - 书籍列表、总数、错误
	List(ctx context.Context, filter *BookFilter, page, size int) ([]*domain.Book, int64, error)

	// GetByISBN 根据ISBN获取书籍
	// 返回: (*domain.Book, error) - 书籍不存在时返回nil, nil
	GetByISBN(ctx context.Context, isbn string) (*domain.Book, error)

	// UpdateStock 更新库存
	// 返回: error - 更新失败时返回错误,书籍不存在或库存不足时返回对应错误码
	UpdateStock(ctx context.Context, bookID string, delta int) error

	// IncrementView 增加浏览次数
	// 返回: error - 更新失败时返回错误
	IncrementView(ctx context.Context, bookID string) error

	// IncrementLike 增加点赞数
	// 返回: error - 更新失败时返回错误
	IncrementLike(ctx context.Context, bookID string) error

	// DecrementLike 减少点赞数
	// 返回: error - 更新失败时返回错误
	DecrementLike(ctx context.Context, bookID string) error

	// TransactionAware 事务感知接口
	// 返回在事务上下文中执行的Repository
	TransactionAware
}

// BookFilter 书籍查询过滤条件
type BookFilter struct {
	// Keyword 关键词搜索 (匹配标题、作者、简介)
	Keyword *string

	// CategoryID 分类ID过滤
	CategoryID *string

	// Tags 标签过滤 (ANY语义: 只要包含任一标签即可)
	Tags []string

	// Status 状态过滤
	Status *domain.BookStatus

	// Author 作者过滤
	Author *string

	// MinPrice 最低价格
	MinPrice *int64

	// MaxPrice 最高价格
	MaxPrice *int64

	// InStock 是否只显示有库存书籍
	InStock *bool

	// SortBy 排序字段 (view_count, like_count, price, created_at, updated_at)
	SortBy *string

	// SortOrder 排序方向 (asc, desc)
	SortOrder *string
}
