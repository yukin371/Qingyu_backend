package ai

import (
	"context"
	"fmt"

	aiModels "Qingyu_backend/models/ai"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoQuotaAlertRepository MongoDB配额告警Repository实现
type MongoQuotaAlertRepository struct {
	collection *mongo.Collection
}

// NewMongoQuotaAlertRepository 创建MongoDB配额告警Repository
func NewMongoQuotaAlertRepository(db *mongo.Database) aiInterfaces.QuotaAlertRepository {
	return &MongoQuotaAlertRepository{
		collection: db.Collection(aiModels.QuotaAlert{}.CollectionName()),
	}
}

// ensureIndexes 创建必要的索引
func (r *MongoQuotaAlertRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "level", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create 创建告警
func (r *MongoQuotaAlertRepository) Create(ctx context.Context, alert *aiModels.QuotaAlert) error {
	alert.BeforeCreate()

	// 确保索引存在
	_ = r.ensureIndexes(ctx)

	_, err := r.collection.InsertOne(ctx, alert)
	if err != nil {
		return fmt.Errorf("创建配额告警失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取告警
func (r *MongoQuotaAlertRepository) GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的告警ID")
	}

	filter := bson.M{"_id": objectID}

	var alert aiModels.QuotaAlert
	err = r.collection.FindOne(ctx, filter).Decode(&alert)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("配额告警不存在")
		}
		return nil, fmt.Errorf("查询配额告警失败: %w", err)
	}

	return &alert, nil
}

// List 获取告警列表
func (r *MongoQuotaAlertRepository) List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	filter := bson.M{}

	if alertType != "" {
		filter["type"] = alertType
	}

	if level != "" {
		filter["level"] = level
	}

	if status != "" {
		filter["status"] = status
	}

	// 统计总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计配额告警数量失败: %w", err)
	}

	// 分页查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询配额告警列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var alerts []*aiModels.QuotaAlert
	if err = cursor.All(ctx, &alerts); err != nil {
		return nil, 0, fmt.Errorf("解析配额告警列表失败: %w", err)
	}

	return alerts, total, nil
}

// GetByUserID 获取用户的所有告警
func (r *MongoQuotaAlertRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	filter := bson.M{"user_id": userID}

	// 统计总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计用户配额告警数量失败: %w", err)
	}

	// 分页查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户配额告警列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var alerts []*aiModels.QuotaAlert
	if err = cursor.All(ctx, &alerts); err != nil {
		return nil, 0, fmt.Errorf("解析用户配额告警列表失败: %w", err)
	}

	return alerts, total, nil
}

// Update 更新告警
func (r *MongoQuotaAlertRepository) Update(ctx context.Context, alert *aiModels.QuotaAlert) error {
	filter := bson.M{"_id": alert.ID}

	update := bson.M{
		"$set": bson.M{
			"status":      alert.Status,
			"resolved_by": alert.ResolvedBy,
			"resolved_at": alert.ResolvedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新配额告警失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("配额告警不存在")
	}

	return nil
}

// GetRecentGlobal 获取最近的全剧告警
func (r *MongoQuotaAlertRepository) GetRecentGlobal(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error) {
	filter := bson.M{"user_id": ""}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询最近全局告警失败: %w", err)
	}
	defer cursor.Close(ctx)

	var alerts []*aiModels.QuotaAlert
	if err = cursor.All(ctx, &alerts); err != nil {
		return nil, fmt.Errorf("解析最近全局告警失败: %w", err)
	}

	return alerts, nil
}

// CountByStatus 按状态统计告警数量
func (r *MongoQuotaAlertRepository) CountByStatus(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("统计告警状态失败: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[aiModels.QuotaAlertStatus]int64)
	var items []struct {
		Status aiModels.QuotaAlertStatus `bson:"_id"`
		Count  int64                      `bson:"count"`
	}

	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("解析告警统计结果失败: %w", err)
	}

	for _, item := range items {
		result[item.Status] = item.Count
	}

	return result, nil
}

// Health 健康检查
func (r *MongoQuotaAlertRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
