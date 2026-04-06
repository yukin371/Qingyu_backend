package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	mongoBase "Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChangeRequestRepositoryMongo 变更建议 MongoDB 仓储实现
type ChangeRequestRepositoryMongo struct {
	*mongoBase.BaseMongoRepository
	db       *mongo.Database
	batchCol *mongo.Collection
}

// NewChangeRequestRepository 创建 ChangeRequestRepository 实例
func NewChangeRequestRepository(db *mongo.Database) writerRepo.ChangeRequestRepository {
	return &ChangeRequestRepositoryMongo{
		BaseMongoRepository: mongoBase.NewBaseMongoRepository(db, "change_requests"),
		db:                  db,
		batchCol:            db.Collection("change_request_batches"),
	}
}

// --- ChangeRequest ---

func (r *ChangeRequestRepositoryMongo) CreateRequest(ctx context.Context, cr *writer.ChangeRequest) error {
	if cr.ID.IsZero() {
		cr.ID = primitive.NewObjectID()
	}
	now := time.Now()
	cr.CreatedAt = now
	cr.UpdatedAt = now

	_, err := r.GetCollection().InsertOne(ctx, cr)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create change request failed", err)
	}
	return nil
}

func (r *ChangeRequestRepositoryMongo) FindRequestByID(ctx context.Context, id string) (*writer.ChangeRequest, error) {
	oid, err := r.ParseID(id)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid change request ID", err)
	}

	var cr writer.ChangeRequest
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": oid}).Decode(&cr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "change request not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find change request failed", err)
	}
	return &cr, nil
}

func (r *ChangeRequestRepositoryMongo) FindRequestsByBatchID(ctx context.Context, batchID string) ([]*writer.ChangeRequest, error) {
	batchOID, err := r.ParseID(batchID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid batch ID", err)
	}

	cursor, err := r.GetCollection().Find(ctx, bson.M{"batch_id": batchOID})
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find requests by batch failed", err)
	}
	defer cursor.Close(ctx)

	var results []*writer.ChangeRequest
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode requests failed", err)
	}
	return results, nil
}

func (r *ChangeRequestRepositoryMongo) FindPendingByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequest, error) {
	return r.FindByChapterAndStatus(ctx, projectID, chapterID, writer.CRStatusPending)
}

func (r *ChangeRequestRepositoryMongo) FindByChapterAndStatus(ctx context.Context, projectID, chapterID string, status writer.ChangeRequestStatus) ([]*writer.ChangeRequest, error) {
	projectOID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	chapterOID, err := r.ParseID(chapterID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid chapter ID", err)
	}

	filter := bson.M{
		"project_id": projectOID,
		"chapter_id": chapterOID,
		"status":     status,
	}

	cursor, err := r.GetCollection().Find(ctx, filter, options.Find().SetSort(bson.D{
		{Key: "processed_at", Value: 1},
		{Key: "created_at", Value: 1},
	}))
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find requests by chapter and status failed", err)
	}
	defer cursor.Close(ctx)

	var results []*writer.ChangeRequest
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode requests by chapter and status failed", err)
	}
	return results, nil
}

func (r *ChangeRequestRepositoryMongo) CountPendingByChapter(ctx context.Context, projectID, chapterID string) (int64, error) {
	projectOID, err := r.ParseID(projectID)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	chapterOID, err := r.ParseID(chapterID)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid chapter ID", err)
	}

	count, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"chapter_id": chapterOID,
		"status":     writer.CRStatusPending,
	})
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count pending requests failed", err)
	}
	return count, nil
}

func (r *ChangeRequestRepositoryMongo) UpdateRequestStatus(ctx context.Context, id string, status writer.ChangeRequestStatus, processedBy string) error {
	oid, err := r.ParseID(id)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid change request ID", err)
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"processed_at": now,
			"processed_by": processedBy,
			"updated_at":   time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update request status failed", err)
	}
	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "change request not found", nil)
	}
	return nil
}

func (r *ChangeRequestRepositoryMongo) DeleteRequest(ctx context.Context, id string) error {
	oid, err := r.ParseID(id)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid change request ID", err)
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete request failed", err)
	}
	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "change request not found", nil)
	}
	return nil
}

// --- Batch ---

func (r *ChangeRequestRepositoryMongo) CreateBatch(ctx context.Context, batch *writer.ChangeRequestBatch) error {
	if batch.ID.IsZero() {
		batch.ID = primitive.NewObjectID()
	}
	now := time.Now()
	batch.CreatedAt = now
	batch.UpdatedAt = now

	_, err := r.batchCol.InsertOne(ctx, batch)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create batch failed", err)
	}
	return nil
}

func (r *ChangeRequestRepositoryMongo) FindBatchByID(ctx context.Context, id string) (*writer.ChangeRequestBatch, error) {
	oid, err := r.ParseID(id)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid batch ID", err)
	}

	var batch writer.ChangeRequestBatch
	err = r.batchCol.FindOne(ctx, bson.M{"_id": oid}).Decode(&batch)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "batch not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find batch failed", err)
	}
	return &batch, nil
}

func (r *ChangeRequestRepositoryMongo) FindBatchesByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequestBatch, error) {
	projectOID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	chapterOID, err := r.ParseID(chapterID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid chapter ID", err)
	}

	cursor, err := r.batchCol.Find(ctx, bson.M{
		"project_id": projectOID,
		"chapter_id": chapterOID,
	})
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find batches failed", err)
	}
	defer cursor.Close(ctx)

	var results []*writer.ChangeRequestBatch
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode batches failed", err)
	}
	return results, nil
}

func (r *ChangeRequestRepositoryMongo) UpdateBatchCounts(ctx context.Context, id string, total, pending int) error {
	oid, err := r.ParseID(id)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid batch ID", err)
	}

	update := bson.M{
		"$set": bson.M{
			"total_count":   total,
			"pending_count": pending,
			"updated_at":    time.Now(),
		},
	}

	result, err := r.batchCol.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update batch counts failed", err)
	}
	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "batch not found", nil)
	}
	return nil
}
