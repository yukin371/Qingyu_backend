package social

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
)

// MongoUserRelationRepository MongoDB用户关系仓储实现
type MongoUserRelationRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRelationRepository 创建MongoDB用户关系仓储
func NewMongoUserRelationRepository(db *mongo.Database) socialRepo.UserRelationRepository {
	return &MongoUserRelationRepository{
		collection: db.Collection("user_relations"),
	}
}

// Create 创建用户关系
func (r *MongoUserRelationRepository) Create(ctx context.Context, relation *social.UserRelation) error {
	if relation.ID.IsZero() {
		relation.ID = primitive.NewObjectID()
	}
	now := time.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, relation)
	return err
}

// GetByID 根据ID获取用户关系
func (r *MongoUserRelationRepository) GetByID(ctx context.Context, id string) (*social.UserRelation, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var relation social.UserRelation
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&relation)
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

// Update 更新用户关系
func (r *MongoUserRelationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	return err
}

// Delete 删除用户关系
func (r *MongoUserRelationRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// GetRelation 获取用户关系
func (r *MongoUserRelationRepository) GetRelation(ctx context.Context, followerID, followeeID string) (*social.UserRelation, error) {
	var relation social.UserRelation
	err := r.collection.FindOne(ctx, bson.M{
		"follower_id": followerID,
		"followee_id": followeeID,
	}).Decode(&relation)
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

// IsFollowing 检查是否关注
func (r *MongoUserRelationRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"follower_id": followerID,
		"followee_id": followeeID,
		"status":      "active",
	})
	return count > 0, err
}

// GetFollowers 获取粉丝列表
func (r *MongoUserRelationRepository) GetFollowers(ctx context.Context, followeeID string, limit, offset int) ([]*social.UserRelation, int64, error) {
	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{
		"followee_id": followeeID,
		"status":      "active",
	})
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{
		"followee_id": followeeID,
		"status":      "active",
	}, opts)
	if err != nil {
		return nil, 0, err
	}

	var relations []*social.UserRelation
	if err = cursor.All(ctx, &relations); err != nil {
		return nil, 0, err
	}

	return relations, total, nil
}

// GetFollowing 获取关注列表
func (r *MongoUserRelationRepository) GetFollowing(ctx context.Context, followerID string, limit, offset int) ([]*social.UserRelation, int64, error) {
	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{
		"follower_id": followerID,
		"status":      "active",
	})
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{
		"follower_id": followerID,
		"status":      "active",
	}, opts)
	if err != nil {
		return nil, 0, err
	}

	var relations []*social.UserRelation
	if err = cursor.All(ctx, &relations); err != nil {
		return nil, 0, err
	}

	return relations, total, nil
}

// CountFollowers 统计粉丝数
func (r *MongoUserRelationRepository) CountFollowers(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{
		"followee_id": userID,
		"status":      "active",
	})
}

// CountFollowing 统计关注数
func (r *MongoUserRelationRepository) CountFollowing(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{
		"follower_id": userID,
		"status":      "active",
	})
}

// BatchCreate 批量创建用户关系
func (r *MongoUserRelationRepository) BatchCreate(ctx context.Context, relations []*social.UserRelation) error {
	if len(relations) == 0 {
		return nil
	}

	now := time.Now()
	documents := make([]interface{}, len(relations))
	for i, relation := range relations {
		if relation.ID.IsZero() {
			relation.ID = primitive.NewObjectID()
		}
		relation.CreatedAt = now
		relation.UpdatedAt = now
		documents[i] = relation
	}

	_, err := r.collection.InsertMany(ctx, documents)
	return err
}
