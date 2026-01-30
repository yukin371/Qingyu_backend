package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TimelineRepositoryMongo Timeline Repository的MongoDB实现
type TimelineRepositoryMongo struct {
	*base.BaseMongoRepository
}

// NewTimelineRepository 创建TimelineRepository实例
func NewTimelineRepository(db *mongo.Database) writerRepo.TimelineRepository {
	return &TimelineRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "timelines"),
	}
}

// Create 创建时间线
func (r *TimelineRepositoryMongo) Create(ctx context.Context, timeline *writer.Timeline) error {
	if timeline.ID == "" {
		timeline.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	timeline.CreatedAt = now
	timeline.UpdatedAt = now

	_, err := r.GetCollection().InsertOne(ctx, timeline)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create timeline failed", err)
	}

	return nil
}

// FindByID 根据ID查询时间线
func (r *TimelineRepositoryMongo) FindByID(ctx context.Context, timelineID string) (*writer.Timeline, error) {
	var timeline writer.Timeline
	filter := bson.M{"_id": timelineID}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&timeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline failed", err)
	}

	return &timeline, nil
}

// FindByProjectID 查询项目下的所有时间线
func (r *TimelineRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Timeline, error) {
	filter := bson.M{"project_id": projectID}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timelines failed", err)
	}
	defer cursor.Close(ctx)

	var timelines []*writer.Timeline
	if err = cursor.All(ctx, &timelines); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode timelines failed", err)
	}

	return timelines, nil
}

// Update 更新时间线
func (r *TimelineRepositoryMongo) Update(ctx context.Context, timeline *writer.Timeline) error {
	timeline.UpdatedAt = time.Now()

	filter := bson.M{"_id": timeline.ID}
	update := bson.M{"$set": timeline}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update timeline failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline not found", nil)
	}

	return nil
}

// Delete 删除时间线
func (r *TimelineRepositoryMongo) Delete(ctx context.Context, timelineID string) error {
	filter := bson.M{"_id": timelineID}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete timeline failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline not found", nil)
	}

	return nil
}

// TimelineEventRepositoryMongo TimelineEvent Repository的MongoDB实现
type TimelineEventRepositoryMongo struct {
	*base.BaseMongoRepository
}

// NewTimelineEventRepository 创建TimelineEventRepository实例
func NewTimelineEventRepository(db *mongo.Database) writerRepo.TimelineEventRepository {
	return &TimelineEventRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "timeline_events"),
	}
}

// Create 创建时间线事件
func (r *TimelineEventRepositoryMongo) Create(ctx context.Context, event *writer.TimelineEvent) error {
	if event.ID == "" {
		event.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	_, err := r.GetCollection().InsertOne(ctx, event)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create timeline event failed", err)
	}

	return nil
}

// FindByID 根据ID查询事件
func (r *TimelineEventRepositoryMongo) FindByID(ctx context.Context, eventID string) (*writer.TimelineEvent, error) {
	var event writer.TimelineEvent
	filter := bson.M{"_id": eventID}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline event failed", err)
	}

	return &event, nil
}

// FindByTimelineID 查询时间线下的所有事件
func (r *TimelineEventRepositoryMongo) FindByTimelineID(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error) {
	filter := bson.M{"timeline_id": timelineID}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline events failed", err)
	}
	defer cursor.Close(ctx)

	var events []*writer.TimelineEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode timeline events failed", err)
	}

	return events, nil
}

// FindByProjectID 查询项目下的所有事件
func (r *TimelineEventRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.TimelineEvent, error) {
	filter := bson.M{"project_id": projectID}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline events failed", err)
	}
	defer cursor.Close(ctx)

	var events []*writer.TimelineEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode timeline events failed", err)
	}

	return events, nil
}

// Update 更新事件
func (r *TimelineEventRepositoryMongo) Update(ctx context.Context, event *writer.TimelineEvent) error {
	event.UpdatedAt = time.Now()

	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": event}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update timeline event failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", nil)
	}

	return nil
}

// Delete 删除事件
func (r *TimelineEventRepositoryMongo) Delete(ctx context.Context, eventID string) error {
	filter := bson.M{"_id": eventID}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete timeline event failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", nil)
	}

	return nil
}
