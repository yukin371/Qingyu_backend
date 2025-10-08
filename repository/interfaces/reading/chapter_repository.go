package reading

import (
	"context"

	"Qingyu_backend/models/reading/reader"
)

// ChapterRepository 章节仓储接口
type ChapterRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, chapter *reader.Chapter) error
	GetByID(ctx context.Context, id string) (*reader.Chapter, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 章节查询
	GetByBookID(ctx context.Context, bookID string) ([]*reader.Chapter, error)
	GetByBookIDWithPagination(ctx context.Context, bookID string, limit, offset int64) ([]*reader.Chapter, error)
	GetByChapterNum(ctx context.Context, bookID string, chapterNum int) (*reader.Chapter, error)

	// 章节导航
	GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error)
	GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error)
	GetFirstChapter(ctx context.Context, bookID string) (*reader.Chapter, error)
	GetLastChapter(ctx context.Context, bookID string) (*reader.Chapter, error)

	// 章节状态查询
	GetPublishedChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error)
	GetVIPChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error)
	GetFreeChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error)

	// 统计查询
	CountByBookID(ctx context.Context, bookID string) (int64, error)
	CountByStatus(ctx context.Context, bookID string, status int) (int64, error)
	CountVIPChapters(ctx context.Context, bookID string) (int64, error)

	// 批量操作
	BatchCreate(ctx context.Context, chapters []*reader.Chapter) error
	BatchUpdateStatus(ctx context.Context, chapterIDs []string, status int) error
	BatchDelete(ctx context.Context, chapterIDs []string) error

	// VIP权限检查
	CheckVIPAccess(ctx context.Context, chapterID string) (bool, error)
	GetChapterPrice(ctx context.Context, chapterID string) (int64, error)

	// 章节内容管理
	GetChapterContent(ctx context.Context, chapterID string) (string, error)
	UpdateChapterContent(ctx context.Context, chapterID string, content string) error

	// 健康检查
	Health(ctx context.Context) error
}
