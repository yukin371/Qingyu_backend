package bookstore

import (
	"context"

	"Qingyu_backend/models/bookstore"
)

// ChapterContentRepository 章节内容仓储接口
type ChapterContentRepository interface {
	// 基础 CRUD 操作
	Create(ctx context.Context, content *bookstore.ChapterContent) error
	GetByID(ctx context.Context, id string) (*bookstore.ChapterContent, error)
	GetByChapterID(ctx context.Context, chapterID string) (*bookstore.ChapterContent, error)
	Update(ctx context.Context, content *bookstore.ChapterContent) error
	Delete(ctx context.Context, id string) error

	// 内容查询
	GetContentByChapterIDs(ctx context.Context, chapterIDs []string) ([]*bookstore.ChapterContent, error)

	// 批量操作
	BatchCreate(ctx context.Context, contents []*bookstore.ChapterContent) error
	BatchDelete(ctx context.Context, ids []string) error

	// 版本管理
	GetLatestVersion(ctx context.Context, chapterID string) (int, error)
	UpdateContent(ctx context.Context, chapterID string, content string) (*bookstore.ChapterContent, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
