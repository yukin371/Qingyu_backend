package ai

import (
	"context"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/repository/interfaces/ai"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoQuotaRepository MongoDB配额Repository实现
type MongoQuotaRepository struct {
	quotaCollection       *mongo.Collection
	transactionCollection *mongo.Collection
}

// NewMongoQuotaRepository 创建MongoDB配额Repository
func NewMongoQuotaRepository(db *mongo.Database) ai.QuotaRepository {
	return &MongoQuotaRepository{
		quotaCollection:       db.Collection(aiModels.UserQuota{}.CollectionName()),
		transactionCollection: db.Collection(aiModels.QuotaTransaction{}.CollectionName()),
	}
}

// CreateQuota 创建配额
func (r *MongoQuotaRepository) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	quota.BeforeCreate()

	_, err := r.quotaCollection.InsertOne(ctx, quota)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("配额记录已存在")
		}
		return fmt.Errorf("创建配额失败: %w", err)
	}

	return nil
}

// GetQuotaByUserID 根据用户ID和配额类型获取配额
func (r *MongoQuotaRepository) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	filter := bson.M{
		"user_id":    userID,
		"quota_type": quotaType,
	}

	var quota aiModels.UserQuota
	err := r.quotaCollection.FindOne(ctx, filter).Decode(&quota)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, aiModels.ErrQuotaNotFound
		}
		return nil, fmt.Errorf("查询配额失败: %w", err)
	}

	// 检查是否需要重置
	if quota.ShouldReset() {
		quota.Reset()
		if err := r.UpdateQuota(ctx, &quota); err != nil {
			return nil, fmt.Errorf("重置配额失败: %w", err)
		}
	}

	return &quota, nil
}

// UpdateQuota 更新配额
func (r *MongoQuotaRepository) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	quota.BeforeUpdate()

	filter := bson.M{
		"user_id":    quota.UserID,
		"quota_type": quota.QuotaType,
	}

	update := bson.M{
		"$set": bson.M{
			"used_quota":      quota.UsedQuota,
			"remaining_quota": quota.RemainingQuota,
			"status":          quota.Status,
			"reset_at":        quota.ResetAt,
			"metadata":        quota.Metadata,
			"updated_at":      quota.UpdatedAt,
		},
	}

	result, err := r.quotaCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新配额失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return aiModels.ErrQuotaNotFound
	}

	return nil
}

// DeleteQuota 删除配额
func (r *MongoQuotaRepository) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	filter := bson.M{
		"user_id":    userID,
		"quota_type": quotaType,
	}

	result, err := r.quotaCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除配额失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return aiModels.ErrQuotaNotFound
	}

	return nil
}

// GetAllQuotasByUserID 获取用户所有配额
func (r *MongoQuotaRepository) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	filter := bson.M{"user_id": userID}

	cursor, err := r.quotaCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询配额列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var quotas []*aiModels.UserQuota
	if err = cursor.All(ctx, &quotas); err != nil {
		return nil, fmt.Errorf("解析配额列表失败: %w", err)
	}

	// 检查并重置过期的配额
	for _, quota := range quotas {
		if quota.ShouldReset() {
			quota.Reset()
			_ = r.UpdateQuota(ctx, quota) // 忽略错误，继续处理其他配额
		}
	}

	return quotas, nil
}

// BatchResetQuotas 批量重置配额
func (r *MongoQuotaRepository) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	filter := bson.M{
		"quota_type": quotaType,
		"reset_at":   bson.M{"$lte": time.Now()},
	}

	// 根据配额类型设置下次重置时间
	var resetAt time.Time
	now := time.Now()
	switch quotaType {
	case aiModels.QuotaTypeDaily:
		resetAt = now.AddDate(0, 0, 1)
	case aiModels.QuotaTypeMonthly:
		resetAt = now.AddDate(0, 1, 0)
	default:
		return fmt.Errorf("不支持的配额类型: %s", quotaType)
	}

	update := bson.M{
		"$set": bson.M{
			"used_quota": 0,
			"status":     aiModels.QuotaStatusActive,
			"reset_at":   resetAt,
			"updated_at": now,
		},
	}

	_, err := r.quotaCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("批量重置配额失败: %w", err)
	}

	return nil
}

// CreateTransaction 创建配额事务记录
func (r *MongoQuotaRepository) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}
	if transaction.Timestamp.IsZero() {
		transaction.Timestamp = time.Now()
	}

	_, err := r.transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return fmt.Errorf("创建配额事务失败: %w", err)
	}

	return nil
}

// GetTransactionsByUserID 获取用户配额事务记录
func (r *MongoQuotaRepository) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	filter := bson.M{"user_id": userID}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.transactionCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询配额事务失败: %w", err)
	}
	defer cursor.Close(ctx)

	var transactions []*aiModels.QuotaTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("解析配额事务失败: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByTimeRange 获取时间范围内的配额事务
func (r *MongoQuotaRepository) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	filter := bson.M{
		"user_id": userID,
		"timestamp": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := r.transactionCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询配额事务失败: %w", err)
	}
	defer cursor.Close(ctx)

	var transactions []*aiModels.QuotaTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("解析配额事务失败: %w", err)
	}

	return transactions, nil
}

// GetQuotaStatistics 获取配额统计信息
func (r *MongoQuotaRepository) GetQuotaStatistics(ctx context.Context, userID string) (*ai.QuotaStatistics, error) {
	// 获取所有配额
	quotas, err := r.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &ai.QuotaStatistics{
		UserID:         userID,
		QuotaByType:    make(map[string]int),
		QuotaByService: make(map[string]int),
	}

	// 统计配额
	for _, quota := range quotas {
		stats.TotalQuota += quota.TotalQuota
		stats.UsedQuota += quota.UsedQuota
		stats.RemainingQuota += quota.RemainingQuota
		stats.QuotaByType[string(quota.QuotaType)] = quota.UsedQuota
	}

	// 计算使用百分比
	if stats.TotalQuota > 0 {
		stats.UsagePercentage = float64(stats.UsedQuota) / float64(stats.TotalQuota) * 100
	}

	// 统计事务
	transactions, err := r.GetTransactionsByUserID(ctx, userID, 1000, 0)
	if err == nil {
		stats.TotalTransactions = len(transactions)

		// 按服务统计
		for _, tx := range transactions {
			if tx.Type == "consume" {
				stats.QuotaByService[tx.Service] += tx.Amount
			}
		}

		// 计算日均消费（最近30天）
		if len(transactions) > 0 {
			thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
			recentTx, _ := r.GetTransactionsByTimeRange(ctx, userID, thirtyDaysAgo, time.Now())
			if len(recentTx) > 0 {
				totalConsumption := 0
				for _, tx := range recentTx {
					if tx.Type == "consume" {
						totalConsumption += tx.Amount
					}
				}
				stats.DailyAverage = float64(totalConsumption) / 30.0
			}
		}
	}

	return stats, nil
}

// GetTotalConsumption 获取总消费量
func (r *MongoQuotaRepository) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	filter := bson.M{
		"user_id":    userID,
		"quota_type": quotaType,
		"type":       "consume",
		"timestamp": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("统计消费量失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("解析统计结果失败: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Total, nil
}

// Health 健康检查
func (r *MongoQuotaRepository) Health(ctx context.Context) error {
	return r.quotaCollection.Database().Client().Ping(ctx, nil)
}
