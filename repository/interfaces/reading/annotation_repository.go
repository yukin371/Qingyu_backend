package reading

import (
	"context"

	"Qingyu_backend/models/reading/reader"
)

// AnnotationRepository 标注（笔记、书签）仓储接口
type AnnotationRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, annotation *reader.Annotation) error
	GetByID(ctx context.Context, id string) (*reader.Annotation, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 查询操作
	GetByUserAndBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
	GetByUserAndChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error)
	GetByType(ctx context.Context, userID, bookID string, annotationType int) ([]*reader.Annotation, error)

	// 笔记操作
	GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
	GetNotesByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error)
	SearchNotes(ctx context.Context, userID string, keyword string) ([]*reader.Annotation, error)

	// 书签操作
	GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
	GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, startOffset int) (*reader.Annotation, error)
	GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error)

	// 高亮操作
	GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
	GetHighlightsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error)

	// 统计操作
	CountByUser(ctx context.Context, userID string) (int64, error)
	CountByBook(ctx context.Context, userID, bookID string) (int64, error)
	CountByType(ctx context.Context, userID string, annotationType int) (int64, error)

	// 批量操作
	BatchCreate(ctx context.Context, annotations []*reader.Annotation) error
	BatchDelete(ctx context.Context, annotationIDs []string) error
	DeleteByBook(ctx context.Context, userID, bookID string) error
	DeleteByChapter(ctx context.Context, userID, bookID, chapterID string) error

	// 数据同步
	SyncAnnotations(ctx context.Context, userID string, annotations []*reader.Annotation) error
	GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error)

	// 分享功能
	GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error)
	GetSharedAnnotations(ctx context.Context, userID string) ([]*reader.Annotation, error)

	// 健康检查
	Health(ctx context.Context) error
}
