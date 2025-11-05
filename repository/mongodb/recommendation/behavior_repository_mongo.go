package recommendation

import (
	reco "Qingyu_backend/models/recommendation"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
)

// MongoBehaviorRepository 用户行为Repository的MongoDB实现
type MongoBehaviorRepository struct {
	collection *mongo.Collection
}

// NewMongoBehaviorRepository 创建MongoBehaviorRepository实例
func NewMongoBehaviorRepository(db *mongo.Database) recoRepo.BehaviorRepository {
	return &MongoBehaviorRepository{
		collection: db.Collection("user_behaviors"),
	}
}

// Create 创建用户行为记录
func (r *MongoBehaviorRepository) Create(ctx context.Context, b *reco.Behavior) error {
	if b == nil {
		return fmt.Errorf("behavior cannot be nil")
	}

	// 生成ID（如果为空）
	if b.ID == "" {
		b.ID = generateID()
	}

	// 设置时间戳
	now := time.Now()
	if b.OccurredAt.IsZero() {
		b.OccurredAt = now
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}

	_, err := r.collection.InsertOne(ctx, b)
	if err != nil {
		return fmt.Errorf("failed to create behavior: %w", err)
	}

	return nil
}

// BatchCreate 批量创建用户行为记录
func (r *MongoBehaviorRepository) BatchCreate(ctx context.Context, bs []*reco.Behavior) error {
	if len(bs) == 0 {
		return nil
	}

	// 转换为interface{}切片
	docs := make([]interface{}, len(bs))
	now := time.Now()
	for i, b := range bs {
		if b.OccurredAt.IsZero() {
			b.OccurredAt = now
		}
		if b.CreatedAt.IsZero() {
			b.CreatedAt = now
		}
		docs[i] = b
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to batch create behaviors: %w", err)
	}

	return nil
}

// GetByUser 获取用户的行为记录
func (r *MongoBehaviorRepository) GetByUser(ctx context.Context, userID string, limit int) ([]*reco.Behavior, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{
			{Key: "occurred_at", Value: -1}, // 按时间倒序
			{Key: "_id", Value: -1},         // 时间相同时按ID倒序，确保排序稳定
		}).
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

// generateID 生成唯一ID（使用雪花算法生成的时间戳+随机数）
func generateID() string {
	return fmt.Sprintf("beh_%d", time.Now().UnixNano())
}
