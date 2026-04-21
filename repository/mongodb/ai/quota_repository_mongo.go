package ai

import (
	"context"
	"fmt"
	"log"
	"sort"
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
	legacyQuotaCollection *mongo.Collection
	transactionCollection *mongo.Collection
	usersCollection       *mongo.Collection
}

const (
	quotaCollectionName             = "ai_quotas"
	legacyQuotaCollectionName       = "ai_user_quotas"
	quotaCollectionSelectionTimeout = 2 * time.Second
)

// NewMongoQuotaRepository 创建MongoDB配额Repository
func NewMongoQuotaRepository(db *mongo.Database) aiInterfaces.QuotaRepository {
	activeCollection, fallbackCollection := selectQuotaCollections(db)
	return &MongoQuotaRepository{
		quotaCollection:       activeCollection,
		legacyQuotaCollection: fallbackCollection,
		transactionCollection: db.Collection(aiModels.QuotaTransaction{}.CollectionName()),
		usersCollection:       db.Collection("users"),
	}
}

// selectQuotaCollections 在启动时选择主/备 quota 集合，主集合维持现有读写口径，
// 同时保留 fallback 集合供 legacy 兼容读写使用。
func selectQuotaCollections(db *mongo.Database) (*mongo.Collection, *mongo.Collection) {
	primary := db.Collection(quotaCollectionName)
	legacy := db.Collection(legacyQuotaCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), quotaCollectionSelectionTimeout)
	defer cancel()

	names, err := db.ListCollectionNames(ctx, bson.M{
		"name": bson.M{
			"$in": bson.A{quotaCollectionName, legacyQuotaCollectionName},
		},
	})
	if err != nil {
		return primary, legacy
	}

	hasPrimary := false
	hasLegacy := false
	for _, name := range names {
		switch name {
		case quotaCollectionName:
			hasPrimary = true
		case legacyQuotaCollectionName:
			hasLegacy = true
		}
	}

	switch {
	case hasPrimary && !hasLegacy:
		return primary, nil
	case !hasPrimary && hasLegacy:
		return legacy, primary
	case !hasPrimary && !hasLegacy:
		return primary, legacy
	}

	primaryCount, err := primary.CountDocuments(ctx, bson.M{})
	if err == nil && primaryCount > 0 {
		legacyCount, legacyErr := legacy.CountDocuments(ctx, bson.M{})
		if legacyErr == nil && legacyCount > 0 {
			log.Printf("[quota] detected mixed quota collections: %s=%d, %s=%d; reads will fallback to legacy collection",
				quotaCollectionName, primaryCount, legacyQuotaCollectionName, legacyCount)
		}
		return primary, legacy
	}

	legacyCount, err := legacy.CountDocuments(ctx, bson.M{})
	if err == nil && legacyCount > 0 {
		return legacy, primary
	}

	return primary, legacy
}

func (r *MongoQuotaRepository) quotaReadCollections() []*mongo.Collection {
	collections := []*mongo.Collection{r.quotaCollection}
	if r.legacyQuotaCollection != nil && r.legacyQuotaCollection.Name() != r.quotaCollection.Name() {
		collections = append(collections, r.legacyQuotaCollection)
	}
	return collections
}

func quotaDedupKey(userID string, quotaType aiModels.QuotaType) string {
	return userID + "::" + string(quotaType)
}

func selectLatestQuota(current, candidate *aiModels.UserQuota) *aiModels.UserQuota {
	if candidate == nil {
		return current
	}
	if current == nil || candidate.UpdatedAt.After(current.UpdatedAt) {
		return candidate
	}
	return current
}

func mergeLatestQuotasByKey(merged map[string]*aiModels.UserQuota, quotas []*aiModels.UserQuota) {
	for _, quota := range quotas {
		if quota == nil {
			continue
		}
		key := quotaDedupKey(quota.UserID, quota.QuotaType)
		merged[key] = selectLatestQuota(merged[key], quota)
	}
}

func (r *MongoQuotaRepository) findQuotaByFilter(ctx context.Context, filter bson.M) (*aiModels.UserQuota, error) {
	var selected *aiModels.UserQuota

	for _, collection := range r.quotaReadCollections() {
		var quota aiModels.UserQuota
		err := collection.FindOne(ctx, filter).Decode(&quota)
		if err == nil {
			quotaCopy := quota
			selected = selectLatestQuota(selected, &quotaCopy)
			continue
		}
		if err != mongo.ErrNoDocuments {
			return nil, fmt.Errorf("查询配额失败: %w", err)
		}
	}

	if selected != nil {
		return selected, nil
	}

	return nil, aiModels.ErrQuotaNotFound
}

func (r *MongoQuotaRepository) findQuotasByFilter(ctx context.Context, filter bson.M) ([]*aiModels.UserQuota, error) {
	merged := make(map[string]*aiModels.UserQuota)

	for _, collection := range r.quotaReadCollections() {
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("查询配额列表失败: %w", err)
		}

		var quotas []*aiModels.UserQuota
		if err = cursor.All(ctx, &quotas); err != nil {
			cursor.Close(ctx)
			return nil, fmt.Errorf("解析配额列表失败: %w", err)
		}
		cursor.Close(ctx)

		mergeLatestQuotasByKey(merged, quotas)
	}

	items := make([]*aiModels.UserQuota, 0, len(merged))
	for _, quota := range merged {
		items = append(items, quota)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}

type quotaUserProfile struct {
	Username string   `bson:"username"`
	Roles    []string `bson:"roles"`
	VIPLevel int      `bson:"vip_level"`
}

func primaryRoleFromRoles(roles []string) string {
	if len(roles) == 0 {
		return ""
	}

	for _, preferred := range []string{"admin", "author", "reader"} {
		for _, role := range roles {
			if role == preferred {
				return role
			}
		}
	}

	return roles[0]
}

func (r *MongoQuotaRepository) getUserProfile(ctx context.Context, userID string) (*quotaUserProfile, error) {
	if r.usersCollection == nil || userID == "" {
		return nil, nil
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, nil
	}

	var profile quotaUserProfile
	if err := r.usersCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&profile); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func mergeBSONMaps(filters ...bson.M) bson.M {
	merged := bson.M{}
	for _, filter := range filters {
		for key, value := range filter {
			merged[key] = value
		}
	}
	return merged
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

	quota, err := r.findQuotaByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 检查是否需要重置
	if quota.ShouldReset() {
		quota.Reset()
		if err := r.UpdateQuota(ctx, quota); err != nil {
			return nil, fmt.Errorf("重置配额失败: %w", err)
		}
	}

	return quota, nil
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

	matched := false
	for _, collection := range r.quotaReadCollections() {
		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("更新配额失败: %w", err)
		}
		if result.MatchedCount > 0 {
			matched = true
		}
	}

	if !matched {
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

	deleted := false
	for _, collection := range r.quotaReadCollections() {
		result, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			return fmt.Errorf("删除配额失败: %w", err)
		}
		if result.DeletedCount > 0 {
			deleted = true
		}
	}

	if !deleted {
		return aiModels.ErrQuotaNotFound
	}

	return nil
}

// GetAllQuotasByUserID 获取用户所有配额
func (r *MongoQuotaRepository) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	filter := bson.M{"user_id": userID}

	quotas, err := r.findQuotasByFilter(ctx, filter)
	if err != nil {
		return nil, err
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

	for _, collection := range r.quotaReadCollections() {
		if _, err := collection.UpdateMany(ctx, filter, update); err != nil {
			return fmt.Errorf("批量重置配额失败: %w", err)
		}
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
	dailyFilter := bson.M{"quota_type": aiModels.QuotaTypeDaily}
	quotas, err := r.findQuotasByFilter(ctx, dailyFilter)
	if err != nil {
		return nil, err
	}

	var activeUsers int64
	var exhaustedUsers int64
	var nearExhaustUsers int64
	var suspendedUsers int64
	for _, quota := range quotas {
		switch quota.Status {
		case aiModels.QuotaStatusActive:
			activeUsers++
			if quota.TotalQuota > 0 && float64(quota.UsedQuota)/float64(quota.TotalQuota) >= 0.8 {
				nearExhaustUsers++
			}
		case aiModels.QuotaStatusExhausted:
			exhaustedUsers++
		case aiModels.QuotaStatusSuspended:
			suspendedUsers++
		}
	}

	// 统计真实消费流水总量和计算平均消费
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"type":   "consume",
			"amount": bson.M{"$gt": 0},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("统计总消费量失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解析统计结果失败: %w", err)
	}

	totalConsumption := int64(0)
	if len(result) > 0 {
		totalConsumption = result[0].Total
	}

	avgConsumption := float64(0)
	totalUsers := int64(len(quotas))
	if totalUsers > 0 {
		avgConsumption = float64(totalConsumption) / float64(totalUsers)
	}

	return &aiModels.DashboardSummary{
		TotalUsers:       totalUsers,
		ActiveUsers:      activeUsers,
		ExhaustedUsers:   exhaustedUsers,
		NearExhaustUsers: nearExhaustUsers,
		SuspendedUsers:   suspendedUsers,
		TotalConsumption: totalConsumption,
		AvgConsumption:   avgConsumption,
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
	dailyMatch := bson.M{"quota_type": aiModels.QuotaTypeDaily}
	quotas, err := r.findQuotasByFilter(ctx, dailyMatch)
	if err != nil {
		return nil, err
	}

	for _, quota := range quotas {
		role := ""
		level := ""
		if quota.Metadata != nil {
			role = quota.Metadata.UserRole
			level = quota.Metadata.MembershipLevel
		}
		if role != "" {
			distribution.ByRole[role]++
		}
		if level != "" {
			distribution.ByLevel[level]++
		}
		if quota.Status != "" {
			distribution.ByStatus[string(quota.Status)]++
		}
	}

	if len(distribution.ByRole) == 0 && r.usersCollection != nil {
		userRolePipeline := mongo.Pipeline{
			{{Key: "$match", Value: bson.M{"roles": bson.M{"$exists": true, "$ne": bson.A{}}}}},
			{{Key: "$unwind", Value: "$roles"}},
			{{Key: "$group", Value: bson.M{
				"_id":   "$roles",
				"count": bson.M{"$sum": 1},
			}}},
		}

		cursor, err := r.usersCollection.Aggregate(ctx, userRolePipeline)
		if err == nil {
			defer cursor.Close(ctx)
			var userRoleResults []struct {
				Role  string `bson:"_id"`
				Count int    `bson:"count"`
			}
			_ = cursor.All(ctx, &userRoleResults)
			for _, result := range userRoleResults {
				distribution.ByRole[result.Role] = result.Count
			}
		}
	}

	servicePipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"type":    "consume",
			"amount":  bson.M{"$gt": 0},
			"service": bson.M{"$exists": true, "$ne": ""},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$service",
			"total": bson.M{"$sum": "$amount"},
		}}},
		{{Key: "$sort", Value: bson.M{"total": -1}}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, servicePipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var serviceResults []struct {
			Service string `bson:"_id"`
			Total   int64  `bson:"total"`
		}
		_ = cursor.All(ctx, &serviceResults)
		for _, result := range serviceResults {
			if result.Service == "" {
				continue
			}
			distribution.ByService[result.Service] = int(result.Total)
		}
	}

	return distribution, nil
}

// GetTopConsumers 获取消费排行
func (r *MongoQuotaRepository) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	quotas, err := r.findQuotasByFilter(ctx, bson.M{"quota_type": aiModels.QuotaTypeDaily})
	if err != nil {
		return nil, err
	}

	sort.Slice(quotas, func(i, j int) bool {
		if quotas[i].UsedQuota == quotas[j].UsedQuota {
			return quotas[i].TotalQuota > quotas[j].TotalQuota
		}
		return quotas[i].UsedQuota > quotas[j].UsedQuota
	})

	if limit > 0 && len(quotas) > limit {
		quotas = quotas[:limit]
	}

	rankings := make([]aiModels.UserQuotaRanking, 0, len(quotas))
	for _, quota := range quotas {
		usagePercent := float64(0)
		if quota.TotalQuota > 0 {
			usagePercent = float64(quota.UsedQuota) / float64(quota.TotalQuota) * 100
		}

		username := quota.UserID
		role := ""
		if quota.Metadata != nil {
			role = quota.Metadata.UserRole
		}
		if profile, err := r.getUserProfile(ctx, quota.UserID); err == nil && profile != nil {
			if profile.Username != "" {
				username = profile.Username
			}
			if role == "" {
				role = primaryRoleFromRoles(profile.Roles)
			}
		}

		rankings = append(rankings, aiModels.UserQuotaRanking{
			UserID:       quota.UserID,
			Username:     username,
			Role:         role,
			UsedQuota:    quota.UsedQuota,
			TotalQuota:   quota.TotalQuota,
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
			"type":      "consume",
			"timestamp": bson.M{"$gte": time.Now().AddDate(0, 0, -days)},
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

	quotas, err := r.findQuotasByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	total := int64(len(quotas))
	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}
	if start > len(quotas) {
		start = len(quotas)
	}
	end := start + limit
	if end > len(quotas) {
		end = len(quotas)
	}

	items := make([]*aiModels.UserQuotaListItem, 0, end-start)
	for _, q := range quotas[start:end] {
		usagePercent := float64(0)
		if q.TotalQuota > 0 {
			usagePercent = float64(q.UsedQuota) / float64(q.TotalQuota) * 100
		}

		username := q.UserID
		role := ""
		memberLevel := ""
		if q.Metadata != nil {
			role = q.Metadata.UserRole
			memberLevel = q.Metadata.MembershipLevel
		}
		if profile, err := r.getUserProfile(ctx, q.UserID); err == nil && profile != nil {
			if profile.Username != "" {
				username = profile.Username
			}
			if role == "" {
				role = primaryRoleFromRoles(profile.Roles)
			}
		}

		items = append(items, &aiModels.UserQuotaListItem{
			UserID:       q.UserID,
			Username:     username,
			Role:         role,
			MemberLevel:  memberLevel,
			DailyQuota:   q.TotalQuota,
			DailyUsed:    q.UsedQuota,
			UsagePercent: usagePercent,
			Status:       string(q.Status),
		})
	}

	return items, total, nil
}
