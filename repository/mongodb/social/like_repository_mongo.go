package social

import (
	"Qingyu_backend/models/social"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoLikeRepository MongoDB点赞仓储实现
type MongoLikeRepository struct {
	collection *mongo.Collection
	// redisClient cache.RedisClient // TODO: 添加Redis支持
}

// NewMongoLikeRepository 创建MongoDB点赞仓储实例
func NewMongoLikeRepository(db *mongo.Database) *MongoLikeRepository {
	collection := db.Collection("likes")

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			// 唯一索引：防止重复点赞
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "target_type", Value: 1},
				{Key: "target_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			// 用户点赞列表查询
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			// 目标点赞数统计
			Keys: bson.D{
				{Key: "target_type", Value: 1},
				{Key: "target_id", Value: 1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create like indexes: %v\n", err)
	}

	return &MongoLikeRepository{
		collection: collection,
	}
}

// AddLike 添加点赞
func (r *MongoLikeRepository) AddLike(ctx context.Context, like *social.Like) error {
	if like.ID.IsZero() {
		like.ID = primitive.NewObjectID()
	}

	if like.CreatedAt.IsZero() {
		like.CreatedAt = time.Now()
	}

	// 参数验证
	if like.UserID == "" || like.TargetType == "" || like.TargetID == "" {
		return fmt.Errorf("点赞参数不完整")
	}

	// 插入点赞记录
	_, err := r.collection.InsertOne(ctx, like)
	if err != nil {
		// 检查是否为重复点赞（唯一索引冲突）
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("已经点赞过了")
		}
		return fmt.Errorf("failed to add like: %w", err)
	}

	// TODO: 更新Redis缓存

	return nil
}

// RemoveLike 取消点赞
func (r *MongoLikeRepository) RemoveLike(ctx context.Context, userID, targetType, targetID string) error {
	if userID == "" || targetType == "" || targetID == "" {
		return fmt.Errorf("取消点赞参数不完整")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{
		"user_id":     userID,
		"target_type": targetType,
		"target_id":   targetID,
	})

	if err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("点赞记录不存在")
	}

	// TODO: 更新Redis缓存

	return nil
}

// IsLiked 检查是否已点赞
func (r *MongoLikeRepository) IsLiked(ctx context.Context, userID, targetType, targetID string) (bool, error) {
	if userID == "" || targetType == "" || targetID == "" {
		return false, fmt.Errorf("查询参数不完整")
	}

	// TODO: 优先从Redis缓存查询

	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id":     userID,
		"target_type": targetType,
		"target_id":   targetID,
	})

	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}

	return count > 0, nil
}

// GetByID 根据ID获取点赞记录
func (r *MongoLikeRepository) GetByID(ctx context.Context, id string) (*social.Like, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid like ID: %w", err)
	}

	var like social.Like
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&like)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("like not found")
		}
		return nil, fmt.Errorf("failed to get like: %w", err)
	}

	return &like, nil
}

// GetUserLikes 获取用户点赞列表
func (r *MongoLikeRepository) GetUserLikes(ctx context.Context, userID, targetType string, page, size int) ([]*social.Like, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	filter := bson.M{"user_id": userID}
	if targetType != "" {
		filter["target_type"] = targetType
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user likes: %w", err)
	}

	// 计算跳过数
	skip := int64((page - 1) * size)

	// 查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(size))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find user likes: %w", err)
	}
	defer cursor.Close(ctx)

	var likes []*social.Like
	if err := cursor.All(ctx, &likes); err != nil {
		return nil, 0, fmt.Errorf("failed to decode user likes: %w", err)
	}

	return likes, total, nil
}

// GetLikeCount 获取目标点赞数
func (r *MongoLikeRepository) GetLikeCount(ctx context.Context, targetType, targetID string) (int64, error) {
	if targetType == "" || targetID == "" {
		return 0, fmt.Errorf("目标参数不完整")
	}

	// TODO: 优先从Redis缓存查询

	count, err := r.collection.CountDocuments(ctx, bson.M{
		"target_type": targetType,
		"target_id":   targetID,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to count likes: %w", err)
	}

	return count, nil
}

// GetLikesCountBatch 批量获取点赞数
func (r *MongoLikeRepository) GetLikesCountBatch(ctx context.Context, targetType string, targetIDs []string) (map[string]int64, error) {
	if targetType == "" || len(targetIDs) == 0 {
		return make(map[string]int64), nil
	}

	// 使用聚合管道批量统计
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"target_type": targetType,
			"target_id":   bson.M{"$in": targetIDs},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$target_id",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate like counts: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]int64)
	for cursor.Next(ctx) {
		var doc struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		result[doc.ID] = doc.Count
	}

	// 为未点赞的目标补0
	for _, targetID := range targetIDs {
		if _, exists := result[targetID]; !exists {
			result[targetID] = 0
		}
	}

	return result, nil
}

// GetUserLikeStatusBatch 批量检查用户点赞状态
func (r *MongoLikeRepository) GetUserLikeStatusBatch(ctx context.Context, userID, targetType string, targetIDs []string) (map[string]bool, error) {
	if userID == "" || targetType == "" || len(targetIDs) == 0 {
		return make(map[string]bool), nil
	}

	// 查询用户对这些目标的点赞记录
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":     userID,
		"target_type": targetType,
		"target_id":   bson.M{"$in": targetIDs},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find user like status: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]bool)
	// 初始化所有为false
	for _, targetID := range targetIDs {
		result[targetID] = false
	}

	// 标记已点赞的
	for cursor.Next(ctx) {
		var like social.Like
		if err := cursor.Decode(&like); err != nil {
			continue
		}
		result[like.TargetID] = true
	}

	return result, nil
}

// CountUserLikes 统计用户点赞总数
func (r *MongoLikeRepository) CountUserLikes(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return 0, fmt.Errorf("failed to count user likes: %w", err)
	}

	return count, nil
}

// CountTargetLikes 统计目标点赞数
func (r *MongoLikeRepository) CountTargetLikes(ctx context.Context, targetType, targetID string) (int64, error) {
	return r.GetLikeCount(ctx, targetType, targetID)
}

// Health 健康检查
func (r *MongoLikeRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
