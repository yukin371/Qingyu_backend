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

// MongoQuotaPolicyRepository MongoDB配额策略Repository实现
type MongoQuotaPolicyRepository struct {
	collection *mongo.Collection
}

// NewMongoQuotaPolicyRepository 创建MongoDB配额策略Repository
func NewMongoQuotaPolicyRepository(db *mongo.Database) aiInterfaces.QuotaPolicyRepository {
	return &MongoQuotaPolicyRepository{
		collection: db.Collection(aiModels.QuotaPolicy{}.CollectionName()),
	}
}

// ensureIndexes 创建必要的索引
func (r *MongoQuotaPolicyRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_role", Value: 1},
				{Key: "membership_level", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "is_default", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create 创建配额策略
func (r *MongoQuotaPolicyRepository) Create(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	policy.BeforeCreate()

	// 确保索引存在
	_ = r.ensureIndexes(ctx)

	_, err := r.collection.InsertOne(ctx, policy)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("配额策略已存在")
		}
		return fmt.Errorf("创建配额策略失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取策略
func (r *MongoQuotaPolicyRepository) GetByID(ctx context.Context, id string) (*aiModels.QuotaPolicy, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的策略ID")
	}

	filter := bson.M{"_id": objectID}

	var policy aiModels.QuotaPolicy
	err = r.collection.FindOne(ctx, filter).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("配额策略不存在")
		}
		return nil, fmt.Errorf("查询配额策略失败: %w", err)
	}

	return &policy, nil
}

// GetByRoleAndLevel 根据角色和等级获取策略
func (r *MongoQuotaPolicyRepository) GetByRoleAndLevel(ctx context.Context, role aiModels.UserRole, level aiModels.MembershipLevel) (*aiModels.QuotaPolicy, error) {
	// 优先精确匹配
	filter := bson.M{
		"user_role":         role,
		"membership_level":  level,
		"status":            aiModels.QuotaPolicyStatusActive,
	}

	var policy aiModels.QuotaPolicy
	err := r.collection.FindOne(ctx, filter).Decode(&policy)
	if err == nil {
		return &policy, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("查询配额策略失败: %w", err)
	}

	// 回退到该角色的默认策略
	fallbackFilter := bson.M{
		"user_role":   role,
		"is_default": true,
		"status":      aiModels.QuotaPolicyStatusActive,
	}

	err = r.collection.FindOne(ctx, fallbackFilter).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("未找到适用的配额策略")
		}
		return nil, fmt.Errorf("查询默认配额策略失败: %w", err)
	}

	return &policy, nil
}

// List 获取策略列表
func (r *MongoQuotaPolicyRepository) List(ctx context.Context, role string, status string, page, limit int) ([]*aiModels.QuotaPolicy, int64, error) {
	filter := bson.M{}

	if role != "" {
		filter["user_role"] = role
	}

	if status != "" {
		filter["status"] = status
	}

	// 统计总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计配额策略数量失败: %w", err)
	}

	// 分页查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询配额策略列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var policies []*aiModels.QuotaPolicy
	if err = cursor.All(ctx, &policies); err != nil {
		return nil, 0, fmt.Errorf("解析配额策略列表失败: %w", err)
	}

	return policies, total, nil
}

// Update 更新策略
func (r *MongoQuotaPolicyRepository) Update(ctx context.Context, policy *aiModels.QuotaPolicy) error {
	policy.BeforeUpdate()

	filter := bson.M{"_id": policy.ID}

	update := bson.M{
		"$set": bson.M{
			"name":             policy.Name,
			"user_role":        policy.UserRole,
			"membership_level": policy.MembershipLevel,
			"daily_quota":      policy.DailyQuota,
			"monthly_quota":    policy.MonthlyQuota,
			"total_quota":      policy.TotalQuota,
			"is_default":       policy.IsDefault,
			"status":           policy.Status,
			"description":      policy.Description,
			"updated_at":       policy.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新配额策略失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("配额策略不存在")
	}

	return nil
}

// Delete 删除策略（软删除）
func (r *MongoQuotaPolicyRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的策略ID")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     aiModels.QuotaPolicyStatusDisabled,
			"updated_at": primitive.NewObjectID().Timestamp(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("删除配额策略失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("配额策略不存在")
	}

	return nil
}

// Health 健康检查
func (r *MongoQuotaPolicyRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
