package bookstore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"Qingyu_backend/models/bookstore"
)

// ChapterContentRepository 章节内容仓储接口
type ChapterContentRepository interface {
	// 基础 CRUD 操作
	Create(ctx context.Context, content *bookstore.ChapterContent) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.ChapterContent, error)
	GetByChapterID(ctx context.Context, chapterID primitive.ObjectID) (*bookstore.ChapterContent, error)
	Update(ctx context.Context, content *bookstore.ChapterContent) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	// 内容查询
	GetContentByChapterIDs(ctx context.Context, chapterIDs []primitive.ObjectID) ([]*bookstore.ChapterContent, error)

	// 批量操作
	BatchCreate(ctx context.Context, contents []*bookstore.ChapterContent) error
	BatchDelete(ctx context.Context, ids []primitive.ObjectID) error

	// 版本管理
	GetLatestVersion(ctx context.Context, chapterID primitive.ObjectID) (int, error)
	UpdateContent(ctx context.Context, chapterID primitive.ObjectID, content string) (*bookstore.ChapterContent, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
