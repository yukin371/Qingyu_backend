package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/shared/types"
	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimelineRepositoryMongo struct {
	*base.BaseMongoRepository
}

func NewTimelineRepository(db *mongo.Database) writerRepo.TimelineRepository {
	return &TimelineRepositoryMongo{BaseMongoRepository: base.NewBaseMongoRepository(db, "timelines")}
}

func timelineObjectID(id string) (primitive.ObjectID, error) {
	oid, err := types.ParseObjectID(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}

func (r *TimelineRepositoryMongo) Create(ctx context.Context, timeline *writer.Timeline) error {
	if timeline.ID.IsZero() {
		timeline.ID = primitive.NewObjectID()
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

func (r *TimelineRepositoryMongo) FindByID(ctx context.Context, timelineID string) (*writer.Timeline, error) {
	objectID, err := timelineObjectID(timelineID)
	if err != nil {
		return nil, err
	}

	var timeline writer.Timeline
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&timeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline failed", err)
	}

	return &timeline, nil
}

func (r *TimelineRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Timeline, error) {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	cursor, err := r.GetCollection().Find(ctx, bson.M{"project_id": projectObjectID})
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

func (r *TimelineRepositoryMongo) Update(ctx context.Context, timeline *writer.Timeline) error {
	timeline.UpdatedAt = time.Now()

	result, err := r.GetCollection().UpdateOne(ctx, bson.M{"_id": timeline.ID}, bson.M{"$set": timeline})
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update timeline failed", err)
	}
	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline not found", nil)
	}

	return nil
}

func (r *TimelineRepositoryMongo) Delete(ctx context.Context, timelineID string) error {
	objectID, err := timelineObjectID(timelineID)
	if err != nil {
		return err
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete timeline failed", err)
	}
	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline not found", nil)
	}

	return nil
}

type TimelineEventRepositoryMongo struct {
	*base.BaseMongoRepository
}

func NewTimelineEventRepository(db *mongo.Database) writerRepo.TimelineEventRepository {
	return &TimelineEventRepositoryMongo{BaseMongoRepository: base.NewBaseMongoRepository(db, "timeline_events")}
}

func (r *TimelineEventRepositoryMongo) Create(ctx context.Context, event *writer.TimelineEvent) error {
	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
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

func (r *TimelineEventRepositoryMongo) FindByID(ctx context.Context, eventID string) (*writer.TimelineEvent, error) {
	objectID, err := timelineObjectID(eventID)
	if err != nil {
		return nil, err
	}

	var event writer.TimelineEvent
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline event failed", err)
	}

	return &event, nil
}

func (r *TimelineEventRepositoryMongo) FindByTimelineID(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error) {
	cursor, err := r.GetCollection().Find(ctx, bson.M{"timeline_id": timelineID})
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

func (r *TimelineEventRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.TimelineEvent, error) {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	cursor, err := r.GetCollection().Find(ctx, bson.M{"project_id": projectObjectID})
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

func (r *TimelineEventRepositoryMongo) Update(ctx context.Context, event *writer.TimelineEvent) error {
	event.UpdatedAt = time.Now()

	result, err := r.GetCollection().UpdateOne(ctx, bson.M{"_id": event.ID}, bson.M{"$set": event})
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update timeline event failed", err)
	}
	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", nil)
	}

	return nil
}

func (r *TimelineEventRepositoryMongo) Delete(ctx context.Context, eventID string) error {
	objectID, err := timelineObjectID(eventID)
	if err != nil {
		return err
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete timeline event failed", err)
	}
	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "timeline event not found", nil)
	}

	return nil
}
