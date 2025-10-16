package recommendation

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	reco "Qingyu_backend/models/recommendation/reco"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
)

// MongoProfileRepository 用户画像Repository的MongoDB实现
type MongoProfileRepository struct {
	collection *mongo.Collection
}

// NewMongoProfileRepository 创建MongoProfileRepository实例
func NewMongoProfileRepository(db *mongo.Database) recoRepo.ProfileRepository {
	return &MongoProfileRepository{
		collection: db.Collection("user_profiles"),
	}
}

// Upsert 更新或创建用户画像
func (r *MongoProfileRepository) Upsert(ctx context.Context, p *reco.UserProfile) error {
	if p == nil {
		return fmt.Errorf("profile cannot be nil")
	}

	if p.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	// 设置时间戳
	now := time.Now()
	p.UpdatedAt = now
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}

	// 使用upsert操作
	filter := bson.M{"user_id": p.UserID}
	update := bson.M{
		"$set": bson.M{
			"tags":       p.Tags,
			"authors":    p.Authors,
			"categories": p.Categories,
			"updated_at": p.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": p.CreatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert profile: %w", err)
	}

	return nil
}

// GetByUserID 根据用户ID获取用户画像
func (r *MongoProfileRepository) GetByUserID(ctx context.Context, userID string) (*reco.UserProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	filter := bson.M{"user_id": userID}

	var profile reco.UserProfile
	err := r.collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 用户画像不存在
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}
