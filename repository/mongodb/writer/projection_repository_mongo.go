package writer

import (
	"context"
	"time"

	modelwriter "Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	mongoBase "Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProjectionRepositoryMongo 章节投影 MongoDB 仓储实现。
type ProjectionRepositoryMongo struct {
	*mongoBase.BaseMongoRepository
}

// NewProjectionRepository 创建 ProjectionRepository 实例。
func NewProjectionRepository(db *mongo.Database) writerRepo.ProjectionRepository {
	return &ProjectionRepositoryMongo{
		BaseMongoRepository: mongoBase.NewBaseMongoRepository(db, "story_harness_projections"),
	}
}

// GetByChapter 根据项目和章节读取投影。
func (r *ProjectionRepositoryMongo) GetByChapter(ctx context.Context, projectID, chapterID string) (*modelwriter.ChapterProjection, error) {
	projectOID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	chapterOID, err := r.ParseID(chapterID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid chapter ID", err)
	}

	var projection modelwriter.ChapterProjection
	err = r.GetCollection().FindOne(ctx, bson.M{
		"project_id": projectOID,
		"chapter_id": chapterOID,
	}).Decode(&projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find projection failed", err)
	}

	return &projection, nil
}

// UpsertByChapter 按项目和章节写入投影。
func (r *ProjectionRepositoryMongo) UpsertByChapter(ctx context.Context, projection *modelwriter.ChapterProjection) error {
	if projection == nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "projection is nil", nil)
	}

	if projection.ID.IsZero() {
		projection.TouchForCreate()
	}
	projection.UpdatedAt = time.Now()
	if projection.CreatedAt.IsZero() {
		projection.CreatedAt = projection.UpdatedAt
	}

	filter := bson.M{
		"project_id": projection.ProjectID,
		"chapter_id": projection.ChapterID,
	}
	update := bson.M{
		"$set": bson.M{
			"project_id": projection.ProjectID,
			"chapter_id": projection.ChapterID,
			"characters": projection.Characters,
			"relations":  projection.Relations,
			"checkpoint": projection.Checkpoint,
			"updated_at": projection.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        projection.ID,
			"created_at": projection.CreatedAt,
		},
	}

	_, err := r.GetCollection().UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "upsert projection failed", err)
	}

	return nil
}
