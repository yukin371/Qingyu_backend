package ai

import (
	"context"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	aiInterfaces "Qingyu_backend/repository/interfaces/ai"

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
func NewMongoQuotaRepository(db *mongo.Database) aiInterfaces.QuotaRepository {
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
func (r *MongoQuotaRepository) GetQuotaStatistics(ctx context.Context, userID string) (*aiInterfaces.QuotaStatistics, error) {
	// 获取所有配额
	quotas, err := r.GetAllQuotasByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &aiInterfaces.QuotaStatistics{
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

// GetDashboardSummary 获取仪表盘汇总数据
func (r *MongoQuotaRepository) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	// 统计总用户数
	totalUsers, err := r.quotaCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("统计总用户数失败: %w", err)
	}

	// 统计活跃用户数
	activeFilter := bson.M{"status": aiModels.QuotaStatusActive}
	activeUsers, err := r.quotaCollection.CountDocuments(ctx, activeFilter)
	if err != nil {
		return nil, fmt.Errorf("统计活跃用户数失败: %w", err)
	}

	// 统计配额耗尽用户数
	exhaustedFilter := bson.M{"status": aiModels.QuotaStatusExhausted}
	exhaustedUsers, err := r.quotaCollection.CountDocuments(ctx, exhaustedFilter)
	if err != nil {
		return nil, fmt.Errorf("统计耗尽用户数失败: %w", err)
	}

	// 统计暂停用户数
	suspendedFilter := bson.M{"status": aiModels.QuotaStatusSuspended}
	suspendedUsers, err := r.quotaCollection.CountDocuments(ctx, suspendedFilter)
	if err != nil {
		return nil, fmt.Errorf("统计暂停用户数失败: %w", err)
	}

	// 统计总消费量和计算平均消费
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$used_quota"},
		}}},
	}

	cursor, err := r.quotaCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("统计总消费量失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解析统计结果失败: %w", err)
	}

	totalConsumption := int64(0)
	if len(result) > 0 {
		totalConsumption = int64(result[0].Total)
	}

	avgConsumption := float64(0)
	if totalUsers > 0 {
		avgConsumption = float64(totalConsumption) / float64(totalUsers)
	}

	return &aiModels.DashboardSummary{
		TotalUsers:      totalUsers,
		ActiveUsers:     activeUsers,
		ExhaustedUsers:  exhaustedUsers,
		NearExhaustUsers: 0, // 需要额外查询接近耗尽的用户
		SuspendedUsers:  suspendedUsers,
		TotalConsumption: totalConsumption,
		AvgConsumption:  avgConsumption,
	}, nil
}

// GetQuotaDistribution 获取配额分布统计
func (r *MongoQuotaRepository) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	distribution := &aiModels.QuotaDistribution{
		ByRole:    make(map[string]int),
		ByLevel:   make(map[string]int),
		ByService: make(map[string]int),
		ByStatus:  make(map[string]int),
	}

	// 按角色统计
	rolePipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$metadata.user_role",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.quotaCollection.Aggregate(ctx, rolePipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var roleResults []struct {
			Role  string `bson:"_id"`
			Count int    `bson:"count"`
		}
		_ = cursor.All(ctx, &roleResults)
		for _, r := range roleResults {
			distribution.ByRole[r.Role] = r.Count
		}
	}

	// 按等级统计
	levelPipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$metadata.membership_level",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err = r.quotaCollection.Aggregate(ctx, levelPipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var levelResults []struct {
			Level string `bson:"_id"`
			Count int    `bson:"count"`
		}
		_ = cursor.All(ctx, &levelResults)
		for _, r := range levelResults {
			distribution.ByLevel[r.Level] = r.Count
		}
	}

	// 按状态统计
	statusPipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err = r.quotaCollection.Aggregate(ctx, statusPipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var statusResults []struct {
			Status string `bson:"_id"`
			Count  int    `bson:"count"`
		}
		_ = cursor.All(ctx, &statusResults)
		for _, r := range statusResults {
			distribution.ByStatus[r.Status] = r.Count
		}
	}

	return distribution, nil
}

// GetTopConsumers 获取消费排行
func (r *MongoQuotaRepository) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	// 聚合查询消费最高的用户
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":        "$user_id",
			"usedQuota":  bson.M{"$sum": "$used_quota"},
			"totalQuota": bson.M{"$first": "$total_quota"},
			"role":       bson.M{"$first": "$metadata.user_role"},
			"username":   bson.M{"$first": "$user_id"},
		}}},
		{{Key: "$sort", Value: bson.M{"usedQuota": -1}}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := r.quotaCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("查询消费排行失败: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		UserID     string `bson:"_id"`
		UsedQuota  int    `bson:"usedQuota"`
		TotalQuota int    `bson:"totalQuota"`
		Role       string `bson:"role"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析消费排行失败: %w", err)
	}

	rankings := make([]aiModels.UserQuotaRanking, 0, len(results))
	for _, r := range results {
		usagePercent := float64(0)
		if r.TotalQuota > 0 {
			usagePercent = float64(r.UsedQuota) / float64(r.TotalQuota) * 100
		}
		rankings = append(rankings, aiModels.UserQuotaRanking{
			UserID:       r.UserID,
			Username:     r.UserID, // 实际应关联用户表获取真实用户名
			Role:         r.Role,
			UsedQuota:    r.UsedQuota,
			TotalQuota:   r.TotalQuota,
			UsagePercent: usagePercent,
		})
	}

	return rankings, nil
}

// GetConsumptionTrend 获取消费趋势
func (r *MongoQuotaRepository) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	// 从事务表获取消费趋势
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"type":       "consume",
			"timestamp":  bson.M{"$gte": time.Now().AddDate(0, 0, -days)},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$timestamp"},
			},
			"consumption": bson.M{"$sum": "$amount"},
			"users":       bson.M{"$addToSet": "$user_id"},
		}}},
		{{Key: "$project", Value: bson.M{
			"date":        "$_id",
			"consumption": 1,
			"users":       bson.M{"$size": "$users"},
		}}},
		{{Key: "$sort", Value: bson.M{"date": 1}}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("查询消费趋势失败: %w", err)
	}
	defer cursor.Close(ctx)

	var results []aiModels.TrendPoint
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析消费趋势失败: %w", err)
	}

	return results, nil
}

// ListUserQuotas 获取用户配额列表
func (r *MongoQuotaRepository) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	filter := bson.M{}

	if role != "" {
		filter["metadata.user_role"] = role
	}

	if status != "" {
		filter["status"] = status
	}

	if search != "" {
		filter["user_id"] = bson.M{"$regex": search, "$options": "i"}
	}

	// 统计总数
	total, err := r.quotaCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计用户配额数量失败: %w", err)
	}

	// 分页查询
	opts := options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.quotaCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户配额列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var quotas []*aiModels.UserQuota
	if err = cursor.All(ctx, &quotas); err != nil {
		return nil, 0, fmt.Errorf("解析用户配额列表失败: %w", err)
	}

	items := make([]*aiModels.UserQuotaListItem, 0, len(quotas))
	for _, q := range quotas {
		usagePercent := float64(0)
		if q.TotalQuota > 0 {
			usagePercent = float64(q.UsedQuota) / float64(q.TotalQuota) * 100
		}

		memberLevel := ""
		if q.Metadata != nil {
			memberLevel = q.Metadata.MembershipLevel
		}

		items = append(items, &aiModels.UserQuotaListItem{
			UserID:       q.UserID,
			Username:     q.UserID, // 实际应关联用户表获取真实用户名
			Role:         string(q.QuotaType),
			MemberLevel:  memberLevel,
			DailyQuota:   q.TotalQuota,
			DailyUsed:    q.UsedQuota,
			UsagePercent: usagePercent,
			Status:       string(q.Status),
		})
	}

	return items, total, nil
}
