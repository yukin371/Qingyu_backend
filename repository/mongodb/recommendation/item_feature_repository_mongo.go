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

// MongoItemFeatureRepository 物品特征Repository的MongoDB实现
type MongoItemFeatureRepository struct {
	collection *mongo.Collection
}

// NewMongoItemFeatureRepository 创建MongoItemFeatureRepository实例
func NewMongoItemFeatureRepository(db *mongo.Database) recoRepo.ItemFeatureRepository {
	return &MongoItemFeatureRepository{
		collection: db.Collection("item_features"),
	}
}

// Create 创建物品特征
func (r *MongoItemFeatureRepository) Create(ctx context.Context, feature *reco.ItemFeature) error {
	if feature == nil {
		return fmt.Errorf("feature cannot be nil")
	}

	if feature.ItemID == "" {
		return fmt.Errorf("itemID cannot be empty")
	}

	// 设置时间戳
	now := time.Now()
	feature.CreatedAt = now
	feature.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, feature)
	if err != nil {
		return fmt.Errorf("failed to create item feature: %w", err)
	}

	return nil
}

// Upsert 更新或创建物品特征
func (r *MongoItemFeatureRepository) Upsert(ctx context.Context, feature *reco.ItemFeature) error {
	if feature == nil {
		return fmt.Errorf("feature cannot be nil")
	}

	if feature.ItemID == "" {
		return fmt.Errorf("itemID cannot be empty")
	}

	// 设置时间戳
	now := time.Now()
	feature.UpdatedAt = now
	if feature.CreatedAt.IsZero() {
		feature.CreatedAt = now
	}

	filter := bson.M{"item_id": feature.ItemID}
	update := bson.M{
		"$set": bson.M{
			"tags":       feature.Tags,
			"authors":    feature.Authors,
			"categories": feature.Categories,
			"vector":     feature.Vector,
			"updated_at": feature.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": feature.CreatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert item feature: %w", err)
	}

	return nil
}

// GetByItemID 根据物品ID获取特征
func (r *MongoItemFeatureRepository) GetByItemID(ctx context.Context, itemID string) (*reco.ItemFeature, error) {
	if itemID == "" {
		return nil, fmt.Errorf("itemID cannot be empty")
	}

	filter := bson.M{"item_id": itemID}

	var feature reco.ItemFeature
	err := r.collection.FindOne(ctx, filter).Decode(&feature)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 特征不存在
		}
		return nil, fmt.Errorf("failed to get item feature: %w", err)
	}

	return &feature, nil
}

// BatchGetByItemIDs 批量获取物品特征
func (r *MongoItemFeatureRepository) BatchGetByItemIDs(ctx context.Context, itemIDs []string) ([]*reco.ItemFeature, error) {
	if len(itemIDs) == 0 {
		return []*reco.ItemFeature{}, nil
	}

	filter := bson.M{"item_id": bson.M{"$in": itemIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get item features: %w", err)
	}
	defer cursor.Close(ctx)

	var features []*reco.ItemFeature
	if err = cursor.All(ctx, &features); err != nil {
		return nil, fmt.Errorf("failed to decode item features: %w", err)
	}

	return features, nil
}

// GetByCategory 根据分类获取物品特征列表
func (r *MongoItemFeatureRepository) GetByCategory(ctx context.Context, category string, limit int) ([]*reco.ItemFeature, error) {
	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	filter := bson.M{"categories": category}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get item features by category: %w", err)
	}
	defer cursor.Close(ctx)

	var features []*reco.ItemFeature
	if err = cursor.All(ctx, &features); err != nil {
		return nil, fmt.Errorf("failed to decode item features: %w", err)
	}

	return features, nil
}

// GetByTags 根据标签获取物品特征列表
// 使用简单的标签匹配策略（包含任一标签）
func (r *MongoItemFeatureRepository) GetByTags(ctx context.Context, tags map[string]float64, limit int) ([]*reco.ItemFeature, error) {
	if len(tags) == 0 {
		return []*reco.ItemFeature{}, nil
	}

	// 构建标签查询（匹配任一标签）
	tagKeys := make([]string, 0, len(tags))
	for key := range tags {
		tagKeys = append(tagKeys, key)
	}

	// 使用$exists查询，找到包含这些标签的物品
	filter := bson.M{}
	tagFilters := make([]bson.M, len(tagKeys))
	for i, key := range tagKeys {
		tagFilters[i] = bson.M{fmt.Sprintf("tags.%s", key): bson.M{"$exists": true}}
	}
	if len(tagFilters) > 0 {
		filter["$or"] = tagFilters
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get item features by tags: %w", err)
	}
	defer cursor.Close(ctx)

	var features []*reco.ItemFeature
	if err = cursor.All(ctx, &features); err != nil {
		return nil, fmt.Errorf("failed to decode item features: %w", err)
	}

	return features, nil
}

// Delete 删除物品特征
func (r *MongoItemFeatureRepository) Delete(ctx context.Context, itemID string) error {
	if itemID == "" {
		return fmt.Errorf("itemID cannot be empty")
	}

	filter := bson.M{"item_id": itemID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete item feature: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoItemFeatureRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
