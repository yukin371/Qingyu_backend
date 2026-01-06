package social

import (
	"context"

	"Qingyu_backend/models/social"
)

// BookListRepository 书单仓储接口
type BookListRepository interface {
	// ========== 书单管理 ==========

	// CreateBookList 创建书单
	CreateBookList(ctx context.Context, bookList *social.BookList) error

	// GetBookListByID 根据ID获取书单
	GetBookListByID(ctx context.Context, bookListID string) (*social.BookList, error)

	// GetBookListsByUser 获取用户的书单列表
	GetBookListsByUser(ctx context.Context, userID string, page, size int) ([]*social.BookList, int64, error)

	// GetPublicBookLists 获取公开书单列表
	GetPublicBookLists(ctx context.Context, page, size int) ([]*social.BookList, int64, error)

	// GetBookListsByCategory 根据分类获取书单
	GetBookListsByCategory(ctx context.Context, category string, page, size int) ([]*social.BookList, int64, error)

	// GetBookListsByTag 根据标签获取书单
	GetBookListsByTag(ctx context.Context, tag string, page, size int) ([]*social.BookList, int64, error)

	// SearchBookLists 搜索书单
	SearchBookLists(ctx context.Context, keyword string, page, size int) ([]*social.BookList, int64, error)

	// UpdateBookList 更新书单
	UpdateBookList(ctx context.Context, bookListID string, updates map[string]interface{}) error

	// DeleteBookList 删除书单
	DeleteBookList(ctx context.Context, bookListID string) error

	// ========== 书单书籍管理 ==========

	// AddBookToList 添加书籍到书单
	AddBookToList(ctx context.Context, bookListID string, bookItem *social.BookListItem) error

	// RemoveBookFromList 从书单中移除书籍
	RemoveBookFromList(ctx context.Context, bookListID, bookID string) error

	// UpdateBookInList 更新书单中的书籍
	UpdateBookInList(ctx context.Context, bookListID, bookID string, updates map[string]interface{}) error

	// ReorderBooks 重新排序书籍
	ReorderBooks(ctx context.Context, bookListID string, bookOrders map[string]int) error

	// GetBooksInList 获取书单中的书籍
	GetBooksInList(ctx context.Context, bookListID string) ([]*social.BookListItem, error)

	// ========== 书单点赞 ==========

	// CreateBookListLike 创建书单点赞
	CreateBookListLike(ctx context.Context, bookListLike *social.BookListLike) error

	// DeleteBookListLike 删除书单点赞
	DeleteBookListLike(ctx context.Context, bookListID, userID string) error

	// GetBookListLike 获取书单点赞记录
	GetBookListLike(ctx context.Context, bookListID, userID string) (*social.BookListLike, error)

	// IsBookListLiked 检查是否已点赞
	IsBookListLiked(ctx context.Context, bookListID, userID string) (bool, error)

	// GetBookListLikes 获取书单点赞列表
	GetBookListLikes(ctx context.Context, bookListID string, page, size int) ([]*social.BookListLike, int64, error)

	// IncrementBookListLikeCount 增加书单点赞数
	IncrementBookListLikeCount(ctx context.Context, bookListID string) error

	// DecrementBookListLikeCount 减少书单点赞数
	DecrementBookListLikeCount(ctx context.Context, bookListID string) error

	// ========== 书单复制 ==========

	// ForkBookList 复制书单
	ForkBookList(ctx context.Context, originalID, userID string) (*social.BookList, error)

	// IncrementForkCount 增加被复制次数
	IncrementForkCount(ctx context.Context, bookListID string) error

	// GetForkedBookLists 获取复制的书单列表
	GetForkedBookLists(ctx context.Context, originalID string, page, size int) ([]*social.BookList, int64, error)

	// ========== 统计 ==========

	// IncrementViewCount 增加浏览次数
	IncrementViewCount(ctx context.Context, bookListID string) error

	// CountUserBookLists 统计用户书单数
	CountUserBookLists(ctx context.Context, userID string) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
