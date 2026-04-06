package writer

import (
	"context"

	"Qingyu_backend/models/writer"
)

// ChangeRequestRepository 变更建议Repository接口
type ChangeRequestRepository interface {
	// ChangeRequest CRUD
	CreateRequest(ctx context.Context, cr *writer.ChangeRequest) error
	FindRequestByID(ctx context.Context, id string) (*writer.ChangeRequest, error)
	FindRequestsByBatchID(ctx context.Context, batchID string) ([]*writer.ChangeRequest, error)
	FindPendingByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequest, error)
	FindByChapterAndStatus(ctx context.Context, projectID, chapterID string, status writer.ChangeRequestStatus) ([]*writer.ChangeRequest, error)
	CountPendingByChapter(ctx context.Context, projectID, chapterID string) (int64, error)
	UpdateRequestStatus(ctx context.Context, id string, status writer.ChangeRequestStatus, processedBy string) error
	DeleteRequest(ctx context.Context, id string) error

	// Batch CRUD
	CreateBatch(ctx context.Context, batch *writer.ChangeRequestBatch) error
	FindBatchByID(ctx context.Context, id string) (*writer.ChangeRequestBatch, error)
	FindBatchesByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequestBatch, error)
	UpdateBatchCounts(ctx context.Context, id string, total, pending int) error
}
