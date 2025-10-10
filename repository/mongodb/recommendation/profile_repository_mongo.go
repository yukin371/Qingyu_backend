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

// ProfileRepositoryMongo MongoDB用户画像仓储实现
type ProfileRepositoryMongo struct {
	collection *mongo.Collection
}

// NewProfileRepository 创建画像仓储实例
func NewProfileRepository(db *mongo.Database) recoRepo.ProfileRepository {
	collection := db.Collection("user_profiles")

	// 创建索引
	ctx := context.Background()
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
	}

	_, _ = collection.Indexes().CreateMany(ctx, indexes)

	return &ProfileRepositoryMongo{
		collection: collection,
	}
}

// Upsert 创建或更新用户画像
func (r *ProfileRepositoryMongo) Upsert(ctx context.Context, p *reco.UserProfile) error {
	if p.ID == "" {
		p.ID = primitive.NewObjectID().Hex()
	}
	p.UpdatedAt = time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	filter := bson.M{"user_id": p.UserID}
	update := bson.M{
		"$set": bson.M{
			"tags":       p.Tags,
			"authors":    p.Authors,
			"categories": p.Categories,
			"updated_at": p.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        p.ID,
			"user_id":    p.UserID,
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

// GetByUserID 根据用户ID获取画像
func (r *ProfileRepositoryMongo) GetByUserID(ctx context.Context, userID string) (*reco.UserProfile, error) {
	filter := bson.M{"user_id": userID}

	var profile reco.UserProfile
	err := r.collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 返回nil表示不存在
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}
