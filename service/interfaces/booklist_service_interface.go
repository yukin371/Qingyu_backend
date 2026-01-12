package interfaces

import (
	socialModel "Qingyu_backend/models/social"
	"context"
)

// BookListService 书单服务接口
type BookListService interface {
	// =========================
	// 书单管理
	// =========================
	CreateBookList(ctx context.Context, userID, userName, userAvatar, title, description, cover, category string, tags []string, isPublic bool) (*socialModel.BookList, error)
	GetBookLists(ctx context.Context, page, size int) ([]*socialModel.BookList, int64, error)
	GetBookListByID(ctx context.Context, bookListID string) (*socialModel.BookList, error)
	UpdateBookList(ctx context.Context, userID, bookListID string, updates map[string]interface{}) error
	DeleteBookList(ctx context.Context, userID, bookListID string) error

	// =========================
	// 书单互动
	// =========================
	LikeBookList(ctx context.Context, userID, bookListID string) error
	ForkBookList(ctx context.Context, userID, bookListID string) (*socialModel.BookList, error)
	GetBooksInList(ctx context.Context, bookListID string) ([]*socialModel.BookListItem, error)
}
