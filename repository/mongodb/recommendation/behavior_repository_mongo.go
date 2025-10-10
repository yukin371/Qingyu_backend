package recommendation

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	reco "Qingyu_backend/models/recommendation/reco"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
)

// BehaviorRepositoryMongo MongoDB用户行为仓储实现
type BehaviorRepositoryMongo struct {
	collection *mongo.Collection
}

// NewBehaviorRepository 创建行为仓储实例
func NewBehaviorRepository(db *mongo.Database) recoRepo.BehaviorRepository {
	collection := db.Collection("user_behaviors")

	// 创建索引
	ctx := context.Background()
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "occurred_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "item_id", Value: 1},
				{Key: "behavior_type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "occurred_at", Value: -1},
			},
		},
	}

	_, _ = collection.Indexes().CreateMany(ctx, indexes)

	return &BehaviorRepositoryMongo{
		collection: collection,
	}
}

// Create 创建单个行为记录
func (r *BehaviorRepositoryMongo) Create(ctx context.Context, b *reco.Behavior) error {
	if b.ID == "" {
		b.ID = primitive.NewObjectID().Hex()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	if b.OccurredAt.IsZero() {
		b.OccurredAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, b)
	if err != nil {
		return fmt.Errorf("failed to create behavior: %w", err)
	}

	return nil
}

// BatchCreate 批量创建行为记录
func (r *BehaviorRepositoryMongo) BatchCreate(ctx context.Context, bs []*reco.Behavior) error {
	if len(bs) == 0 {
		return nil
	}

	docs := make([]interface{}, len(bs))
	now := time.Now()

	for i, b := range bs {
		if b.ID == "" {
			b.ID = primitive.NewObjectID().Hex()
		}
		if b.CreatedAt.IsZero() {
			b.CreatedAt = now
		}
		if b.OccurredAt.IsZero() {
			b.OccurredAt = now
		}
		docs[i] = b
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to batch create behaviors: %w", err)
	}

	return nil
}

// GetByUser 获取用户行为记录（最近N条）
func (r *BehaviorRepositoryMongo) GetByUser(ctx context.Context, userID string, limit int) ([]*reco.Behavior, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "occurred_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find behaviors: %w", err)
	}
	defer cursor.Close(ctx)

	var behaviors []*reco.Behavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, fmt.Errorf("failed to decode behaviors: %w", err)
	}

	return behaviors, nil
}
