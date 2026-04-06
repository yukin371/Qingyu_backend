package writer

import (
	"context"

	modelwriter "Qingyu_backend/models/writer"
)

// ProjectionRepository 章节上下文投影仓储接口。
type ProjectionRepository interface {
	GetByChapter(ctx context.Context, projectID, chapterID string) (*modelwriter.ChapterProjection, error)
	UpsertByChapter(ctx context.Context, projection *modelwriter.ChapterProjection) error
}
