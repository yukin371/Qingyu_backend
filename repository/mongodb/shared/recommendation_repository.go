package shared

import (
	"context"
	"fmt"
	"time"

	recModel "Qingyu_backend/models/shared/recommendation"
	sharedInterfaces "Qingyu_backend/repository/interfaces/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RecommendationRepositoryImpl 推荐Repository实现
type RecommendationRepositoryImpl struct {
	db                 *mongo.Database
	behaviorCollection *mongo.Collection
	profileCollection  *mongo.Collection
}

// NewRecommendationRepository 创建推荐Repository
func NewRecommendationRepository(db *mongo.Database) sharedInterfaces.RecommendationRepository {
	return &RecommendationRepositoryImpl{
		db:                 db,
		behaviorCollection: db.Collection("user_behaviors"),
		profileCollection:  db.Collection("user_profiles"),
	}
}

// ============ 行为记录 ============

// RecordBehavior 记录用户行为
func (r *RecommendationRepositoryImpl) RecordBehavior(ctx context.Context, behavior *recModel.UserBehavior) error {
	behavior.CreatedAt = time.Now()

	result, err := r.behaviorCollection.InsertOne(ctx, behavior)
	if err != nil {
		return fmt.Errorf("记录用户行为失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		behavior.ID = oid.Hex()
	}

	return nil
}

// GetUserBehaviors 获取用户行为记录
func (r *RecommendationRepositoryImpl) GetUserBehaviors(ctx context.Context, userID string, limit int) ([]*recModel.UserBehavior, error) {
	query := bson.M{"user_id": userID}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.behaviorCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("查询用户行为失败: %w", err)
	}
	defer cursor.Close(ctx)

	var behaviors []*recModel.UserBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, fmt.Errorf("解析用户行为失败: %w", err)
	}

	return behaviors, nil
}

// GetItemBehaviors 获取物品的行为记录
func (r *RecommendationRepositoryImpl) GetItemBehaviors(ctx context.Context, itemID string, limit int) ([]*recModel.UserBehavior, error) {
	query := bson.M{"item_id": itemID}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.behaviorCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("查询物品行为失败: %w", err)
	}
	defer cursor.Close(ctx)

	var behaviors []*recModel.UserBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, fmt.Errorf("解析物品行为失败: %w", err)
	}

	return behaviors, nil
}

// GetItemStatistics 获取物品统计信息
func (r *RecommendationRepositoryImpl) GetItemStatistics(ctx context.Context, itemID string) (map[string]int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "item_id", Value: itemID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$action_type"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.behaviorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合统计失败: %w", err)
	}
	defer cursor.Close(ctx)

	stats := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats[result.ID] = result.Count
	}

	return stats, nil
}

// GetHotItems 获取热门物品
func (r *RecommendationRepositoryImpl) GetHotItems(ctx context.Context, itemType string, limit int) ([]string, error) {
	// 获取最近7天的热门物品
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "item_type", Value: itemType},
			{Key: "created_at", Value: bson.D{{Key: "$gte", Value: sevenDaysAgo}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$item_id"},
			{Key: "view_count", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{"$action_type", "view"}}},
					1,
					0,
				}},
			}}}},
			{Key: "favorite_count", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{"$action_type", "favorite"}}},
					1,
					0,
				}},
			}}}},
			{Key: "read_count", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{"$action_type", "read"}}},
					1,
					0,
				}},
			}}}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "score", Value: bson.D{{Key: "$add", Value: bson.A{
				"$view_count",
				bson.D{{Key: "$multiply", Value: bson.A{"$favorite_count", 3}}},
				bson.D{{Key: "$multiply", Value: bson.A{"$read_count", 2}}},
			}}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "score", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := r.behaviorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合热门物品失败: %w", err)
	}
	defer cursor.Close(ctx)

	var itemIDs []string
	for cursor.Next(ctx) {
		var result struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		itemIDs = append(itemIDs, result.ID)
	}

	return itemIDs, nil
}

// GetUserPreferences 获取用户偏好（基于行为分析）
func (r *RecommendationRepositoryImpl) GetUserPreferences(ctx context.Context, userID string) (map[string]float64, error) {
	// 获取用户最近30天的行为
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "created_at", Value: bson.D{{Key: "$gte", Value: thirtyDaysAgo}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$item_type"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "total_duration", Value: bson.D{{Key: "$sum", Value: "$duration"}}},
		}}},
	}

	cursor, err := r.behaviorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合用户偏好失败: %w", err)
	}
	defer cursor.Close(ctx)

	preferences := make(map[string]float64)
	for cursor.Next(ctx) {
		var result struct {
			ID            string  `bson:"_id"`
			Count         int64   `bson:"count"`
			TotalDuration float64 `bson:"total_duration"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		// 偏好分数 = 行为次数 + (总时长 / 3600)
		preferences[result.ID] = float64(result.Count) + (result.TotalDuration / 3600)
	}

	return preferences, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *RecommendationRepositoryImpl) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
