package stats

import (
	"Qingyu_backend/models/stats"
	"Qingyu_backend/repository/interfaces/infrastructure"
	statsInterface "Qingyu_backend/repository/interfaces/stats"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoBookStatsRepository MongoDB实现
type MongoBookStatsRepository struct {
	collection      *mongo.Collection
	dailyCollection *mongo.Collection
}

// NewMongoBookStatsRepository 创建MongoDB Repository
func NewMongoBookStatsRepository(db *mongo.Database) statsInterface.BookStatsRepository {
	return &MongoBookStatsRepository{
		collection:      db.Collection("book_stats"),
		dailyCollection: db.Collection("book_stats_daily"),
	}
}

// Create 创建作品统计
func (r *MongoBookStatsRepository) Create(ctx context.Context, bookStats *stats.BookStats) error {
	bookStats.ID = primitive.NewObjectID().Hex()
	bookStats.CreatedAt = time.Now()
	bookStats.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, bookStats)
	return err
}

// GetByID 根据ID获取
func (r *MongoBookStatsRepository) GetByID(ctx context.Context, id string) (*stats.BookStats, error) {
	var bookStats stats.BookStats
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bookStats)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &bookStats, err
}

// Update 更新
func (r *MongoBookStatsRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

// Delete 删除
func (r *MongoBookStatsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetByBookID 根据作品ID获取统计
func (r *MongoBookStatsRepository) GetByBookID(ctx context.Context, bookID string) (*stats.BookStats, error) {
	var bookStats stats.BookStats
	opts := options.FindOne().SetSort(bson.D{{Key: "stat_date", Value: -1}})
	err := r.collection.FindOne(ctx, bson.M{"book_id": bookID}, opts).Decode(&bookStats)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &bookStats, err
}

// GetByAuthorID 根据作者ID获取作品列表
func (r *MongoBookStatsRepository) GetByAuthorID(ctx context.Context, authorID string, limit, offset int64) ([]*stats.BookStats, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "total_views", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"author_id": authorID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookStats []*stats.BookStats
	if err = cursor.All(ctx, &bookStats); err != nil {
		return nil, err
	}

	return bookStats, nil
}

// GetByDateRange 根据日期范围获取统计
func (r *MongoBookStatsRepository) GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.BookStats, error) {
	filter := bson.M{
		"book_id": bookID,
		"stat_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookStats []*stats.BookStats
	if err = cursor.All(ctx, &bookStats); err != nil {
		return nil, err
	}

	return bookStats, nil
}

// CreateDailyStats 创建每日统计
func (r *MongoBookStatsRepository) CreateDailyStats(ctx context.Context, dailyStats *stats.BookStatsDaily) error {
	dailyStats.ID = primitive.NewObjectID().Hex()
	dailyStats.CreatedAt = time.Now()
	dailyStats.UpdatedAt = time.Now()

	_, err := r.dailyCollection.InsertOne(ctx, dailyStats)
	return err
}

// GetDailyStats 获取指定日期的统计
func (r *MongoBookStatsRepository) GetDailyStats(ctx context.Context, bookID string, date time.Time) (*stats.BookStatsDaily, error) {
	var dailyStats stats.BookStatsDaily
	err := r.dailyCollection.FindOne(ctx, bson.M{
		"book_id": bookID,
		"date":    date.Truncate(24 * time.Hour),
	}).Decode(&dailyStats)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &dailyStats, err
}

// GetDailyStatsRange 获取日期范围的统计
func (r *MongoBookStatsRepository) GetDailyStatsRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.BookStatsDaily, error) {
	filter := bson.M{
		"book_id": bookID,
		"date": bson.M{
			"$gte": startDate.Truncate(24 * time.Hour),
			"$lte": endDate.Truncate(24 * time.Hour),
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}})
	cursor, err := r.dailyCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dailyStats []*stats.BookStatsDaily
	if err = cursor.All(ctx, &dailyStats); err != nil {
		return nil, err
	}

	return dailyStats, nil
}

// GetRevenueBreakdown 获取收入细分
func (r *MongoBookStatsRepository) GetRevenueBreakdown(ctx context.Context, bookID string, startDate, endDate time.Time) (*stats.RevenueBreakdown, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"book_id": bookID,
			"stat_date": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":               nil,
			"chapter_revenue":   bson.M{"$sum": "$chapter_revenue"},
			"subscribe_revenue": bson.M{"$sum": "$subscribe_revenue"},
			"reward_revenue":    bson.M{"$sum": "$reward_revenue"},
			"total_revenue":     bson.M{"$sum": "$total_revenue"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []stats.RevenueBreakdown
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &stats.RevenueBreakdown{
			BookID:    bookID,
			StartDate: startDate,
			EndDate:   endDate,
		}, nil
	}

	breakdown := &result[0]
	breakdown.BookID = bookID
	breakdown.StartDate = startDate
	breakdown.EndDate = endDate

	return breakdown, nil
}

// CalculateTotalRevenue 计算总收入
func (r *MongoBookStatsRepository) CalculateTotalRevenue(ctx context.Context, bookID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"total_revenue": bson.M{"$sum": "$total_revenue"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	totalRevenue, _ := result[0]["total_revenue"].(float64)
	return totalRevenue, nil
}

// CalculateRevenueByType 按类型计算收入
func (r *MongoBookStatsRepository) CalculateRevenueByType(ctx context.Context, bookID string, revenueType string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id":     nil,
			"revenue": bson.M{"$sum": "$" + revenueType},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	revenue, _ := result[0]["revenue"].(float64)
	return revenue, nil
}

// GetTopChapters 获取热门章节
func (r *MongoBookStatsRepository) GetTopChapters(ctx context.Context, bookID string) (*stats.TopChapters, error) {
	// 此方法需要与ChapterStatsRepository配合
	// 简化实现：返回空结果
	return &stats.TopChapters{
		BookID:           bookID,
		MostViewed:       []*stats.ChapterStatsAggregate{},
		HighestRevenue:   []*stats.ChapterStatsAggregate{},
		LowestCompletion: []*stats.ChapterStatsAggregate{},
		HighestDropOff:   []*stats.ChapterStatsAggregate{},
	}, nil
}

// AnalyzeViewTrend 分析阅读量趋势
func (r *MongoBookStatsRepository) AnalyzeViewTrend(ctx context.Context, bookID string, days int) (string, error) {
	// 获取最近N天的每日统计
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	dailyStats, err := r.GetDailyStatsRange(ctx, bookID, startDate, endDate)
	if err != nil {
		return stats.TrendStable, err
	}

	if len(dailyStats) < 2 {
		return stats.TrendStable, nil
	}

	// 简单趋势分析：比较前半段和后半段的平均值
	mid := len(dailyStats) / 2
	var firstHalfAvg, secondHalfAvg int64

	for i := 0; i < mid; i++ {
		firstHalfAvg += dailyStats[i].DailyViews
	}
	firstHalfAvg /= int64(mid)

	for i := mid; i < len(dailyStats); i++ {
		secondHalfAvg += dailyStats[i].DailyViews
	}
	secondHalfAvg /= int64(len(dailyStats) - mid)

	// 判断趋势
	if secondHalfAvg > firstHalfAvg*11/10 { // 增长超过10%
		return stats.TrendUp, nil
	} else if secondHalfAvg < firstHalfAvg*9/10 { // 下降超过10%
		return stats.TrendDown, nil
	}

	return stats.TrendStable, nil
}

// AnalyzeRevenueTrend 分析收入趋势
func (r *MongoBookStatsRepository) AnalyzeRevenueTrend(ctx context.Context, bookID string, days int) (string, error) {
	// 获取最近N天的每日统计
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	dailyStats, err := r.GetDailyStatsRange(ctx, bookID, startDate, endDate)
	if err != nil {
		return stats.TrendStable, err
	}

	if len(dailyStats) < 2 {
		return stats.TrendStable, nil
	}

	// 简单趋势分析：比较前半段和后半段的平均值
	mid := len(dailyStats) / 2
	var firstHalfAvg, secondHalfAvg float64

	for i := 0; i < mid; i++ {
		firstHalfAvg += dailyStats[i].DailyRevenue
	}
	firstHalfAvg /= float64(mid)

	for i := mid; i < len(dailyStats); i++ {
		secondHalfAvg += dailyStats[i].DailyRevenue
	}
	secondHalfAvg /= float64(len(dailyStats) - mid)

	// 判断趋势
	if secondHalfAvg > firstHalfAvg*1.1 { // 增长超过10%
		return stats.TrendUp, nil
	} else if secondHalfAvg < firstHalfAvg*0.9 { // 下降超过10%
		return stats.TrendDown, nil
	}

	return stats.TrendStable, nil
}

// CalculateAvgCompletionRate 计算平均完读率
func (r *MongoBookStatsRepository) CalculateAvgCompletionRate(ctx context.Context, bookID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id":                 nil,
			"avg_completion_rate": bson.M{"$avg": "$avg_completion_rate"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	avgRate, _ := result[0]["avg_completion_rate"].(float64)
	return avgRate, nil
}

// CalculateAvgDropOffRate 计算平均跳出率
func (r *MongoBookStatsRepository) CalculateAvgDropOffRate(ctx context.Context, bookID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id":               nil,
			"avg_drop_off_rate": bson.M{"$avg": "$avg_drop_off_rate"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	avgRate, _ := result[0]["avg_drop_off_rate"].(float64)
	return avgRate, nil
}

// CalculateAvgReadingDuration 计算平均阅读时长
func (r *MongoBookStatsRepository) CalculateAvgReadingDuration(ctx context.Context, bookID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id":                  nil,
			"avg_reading_duration": bson.M{"$avg": "$avg_reading_duration"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	avgDuration, _ := result[0]["avg_reading_duration"].(float64)
	return avgDuration, nil
}

// GetTopBooksByViews 获取阅读量最高的作品
func (r *MongoBookStatsRepository) GetTopBooksByViews(ctx context.Context, limit int) ([]*stats.BookStats, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "total_views", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookStats []*stats.BookStats
	if err = cursor.All(ctx, &bookStats); err != nil {
		return nil, err
	}

	return bookStats, nil
}

// GetTopBooksByRevenue 获取收入最高的作品
func (r *MongoBookStatsRepository) GetTopBooksByRevenue(ctx context.Context, limit int) ([]*stats.BookStats, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "total_revenue", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookStats []*stats.BookStats
	if err = cursor.All(ctx, &bookStats); err != nil {
		return nil, err
	}

	return bookStats, nil
}

// GetTopBooksByCompletion 获取完读率最高的作品
func (r *MongoBookStatsRepository) GetTopBooksByCompletion(ctx context.Context, limit int) ([]*stats.BookStats, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "avg_completion_rate", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookStats []*stats.BookStats
	if err = cursor.All(ctx, &bookStats); err != nil {
		return nil, err
	}

	return bookStats, nil
}

// BatchCreate 批量创建
func (r *MongoBookStatsRepository) BatchCreate(ctx context.Context, bookStats []*stats.BookStats) error {
	if len(bookStats) == 0 {
		return nil
	}

	docs := make([]interface{}, len(bookStats))
	for i, stat := range bookStats {
		stat.ID = primitive.NewObjectID().Hex()
		stat.CreatedAt = time.Now()
		stat.UpdatedAt = time.Now()
		docs[i] = stat
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

// BatchUpdate 批量更新
func (r *MongoBookStatsRepository) BatchUpdate(ctx context.Context, updates []map[string]interface{}) error {
	// TODO: 实现批量更新
	return nil
}

// Count 统计数量
func (r *MongoBookStatsRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByAuthor 按作者统计
func (r *MongoBookStatsRepository) CountByAuthor(ctx context.Context, authorID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"author_id": authorID})
}

// Health 健康检查
func (r *MongoBookStatsRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
