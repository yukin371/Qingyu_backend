package interfaces

import (
	"context"

	readerModels "Qingyu_backend/models/reader"
	readerservice "Qingyu_backend/service/reader"
)

// BookmarkService 书签服务接口
type BookmarkService interface {
	// CreateBookmark 创建书签
	CreateBookmark(ctx context.Context, bookmark *readerModels.Bookmark) error

	// GetBookmark 获取书签详情
	GetBookmark(ctx context.Context, bookmarkID string) (*readerModels.Bookmark, error)

	// GetUserBookmarks 获取用户书签列表
	GetUserBookmarks(ctx context.Context, userID string, filter *readerModels.BookmarkFilter, page, size int) (*readerservice.BookmarkListResponse, error)

	// GetBookBookmarks 获取书籍书签
	GetBookBookmarks(ctx context.Context, userID, bookID string, page, size int) (*readerservice.BookmarkListResponse, error)

	// UpdateBookmark 更新书签
	UpdateBookmark(ctx context.Context, bookmarkID string, bookmark *readerModels.Bookmark) error

	// DeleteBookmark 删除书签
	DeleteBookmark(ctx context.Context, bookmarkID string) error

	// ExportBookmarks 导出书签
	ExportBookmarks(ctx context.Context, userID, format string) ([]byte, string, error)

	// GetBookmarkStats 获取书签统计
	GetBookmarkStats(ctx context.Context, userID string) (*readerModels.BookmarkStats, error)

	// SearchBookmarks 搜索书签
	SearchBookmarks(ctx context.Context, userID, keyword string, page, size int) (*readerservice.BookmarkListResponse, error)
}
