package mongodb

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	"Qingyu_backend/repository/interfaces/writing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TimelineRepositoryMongo Timeline Repository的MongoDB实现
type TimelineRepositoryMongo struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewTimelineRepository 创建TimelineRepository实例
func NewTimelineRepository(db *mongo.Database) writing.TimelineRepository {
	return &TimelineRepositoryMongo{
		db:         db,
		collection: db.Collection("timelines"),
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

	_, err := r.collection.InsertOne(ctx, timeline)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create timeline failed", err)
	}

	return nil
}

// FindByID 根据ID查询时间线
func (r *TimelineRepositoryMongo) FindByID(ctx context.Context, timelineID string) (*writer.Timeline, error) {
	var timeline writer.Timeline
	filter := bson.M{"_id": timelineID}

	err := r.collection.FindOne(ctx, filter).Decode(&timeline)
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

	cursor, err := r.collection.Find(ctx, filter)
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

	result, err := r.collection.UpdateOne(ctx, filter, update)
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

	result, err := r.collection.DeleteOne(ctx, filter)
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
	db         *mongo.Database
	collection *mongo.Collection
}

// NewTimelineEventRepository 创建TimelineEventRepository实例
func NewTimelineEventRepository(db *mongo.Database) writing.TimelineEventRepository {
	return &TimelineEventRepositoryMongo{
		db:         db,
		collection: db.Collection("timeline_events"),
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

	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create timeline event failed", err)
	}

	return nil
}

// FindByID 根据ID查询事件
func (r *TimelineEventRepositoryMongo) FindByID(ctx context.Context, eventID string) (*writer.TimelineEvent, error) {
	var event writer.TimelineEvent
	filter := bson.M{"_id": eventID}

	err := r.collection.FindOne(ctx, filter).Decode(&event)
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

	cursor, err := r.collection.Find(ctx, filter)
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

	cursor, err := r.collection.Find(ctx, filter)
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

	result, err := r.collection.UpdateOne(ctx, filter, update)
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

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete timeline event failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", nil)
	}

	return nil
}
