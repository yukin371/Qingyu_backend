package reader

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reader"
)

// BookmarkRepository 书签仓储接口
type BookmarkRepository interface {
	// Create 创建书签
	Create(ctx context.Context, bookmark *reader.Bookmark) error

	// GetByID 根据ID获取书签
	GetByID(ctx context.Context, id primitive.ObjectID) (*reader.Bookmark, error)

	// GetByUser 获取用户的所有书签
	GetByUser(ctx context.Context, userID primitive.ObjectID, filter *reader.BookmarkFilter, page, size int) ([]*reader.Bookmark, int64, error)

	// GetByBook 获取某本书的书签
	GetByBook(ctx context.Context, userID, bookID primitive.ObjectID, page, size int) ([]*reader.Bookmark, int64, error)

	// GetByChapter 获取某章节的书签
	GetByChapter(ctx context.Context, userID, chapterID primitive.ObjectID, page, size int) ([]*reader.Bookmark, int64, error)

	// Update 更新书签
	Update(ctx context.Context, bookmark *reader.Bookmark) error

	// Delete 删除书签
	Delete(ctx context.Context, id primitive.ObjectID) error

	// DeleteByBook 删除某本书的所有书签
	DeleteByBook(ctx context.Context, userID, bookID primitive.ObjectID) error

	// Exists 检查书签是否存在
	Exists(ctx context.Context, userID, chapterID primitive.ObjectID, position int) (bool, error)

	// GetStats 获取书签统计
	GetStats(ctx context.Context, userID primitive.ObjectID) (*reader.BookmarkStats, error)

	// GetPublicBookmarks 获取公开书签
	GetPublicBookmarks(ctx context.Context, page, size int) ([]*reader.Bookmark, int64, error)

	// Search 搜索书签
	Search(ctx context.Context, userID primitive.ObjectID, keyword string, page, size int) ([]*reader.Bookmark, int64, error)
}
