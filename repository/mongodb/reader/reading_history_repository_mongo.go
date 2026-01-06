package reader

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
)

// MongoReadingHistoryRepository MongoDB实现的阅读历史Repository
type MongoReadingHistoryRepository struct {
	collection *mongo.Collection
}

// NewMongoReadingHistoryRepository 创建MongoDB阅读历史Repository
func NewMongoReadingHistoryRepository(db *mongo.Database) readerRepo.ReadingHistoryRepository {
	collection := db.Collection("reading_histories")

	// 创建索引
	go createReadingHistoryIndexes(collection)

	return &MongoReadingHistoryRepository{
		collection: collection,
	}
}

// createReadingHistoryIndexes 创建索引
func createReadingHistoryIndexes(collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		// 用户ID + 创建时间（倒序）
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("idx_user_created"),
		},
		// 用户ID + 书籍ID + 创建时间（倒序）
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "book_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("idx_user_book_created"),
		},
		// TTL索引：90天自动清理
		{
			Keys: bson.D{
				{Key: "created_at", Value: 1},
			},
			Options: options.Index().
				SetName("idx_created_ttl").
				SetExpireAfterSeconds(7776000), // 90天 = 90 * 24 * 60 * 60
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// 日志记录错误，但不中断程序
		fmt.Printf("创建阅读历史索引失败: %v\n", err)
	}
}

// =======================
// 基础CRUD
// =======================

// Create 创建阅读历史记录
func (r *MongoReadingHistoryRepository) Create(ctx context.Context, history *reader.ReadingHistory) error {
	if history == nil {
		return fmt.Errorf("history不能为空")
	}

	// 设置创建时间
	if history.CreatedAt.IsZero() {
		history.CreatedAt = time.Now()
	}

	// 设置ID
	if history.ID.IsZero() {
		history.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, history)
	if err != nil {
		return fmt.Errorf("创建阅读历史失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取阅读历史
func (r *MongoReadingHistoryRepository) GetByID(ctx context.Context, id string) (*reader.ReadingHistory, error) {
	if id == "" {
		return nil, fmt.Errorf("ID不能为空")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	var history reader.ReadingHistory
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&history)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("阅读历史不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}

	return &history, nil
}

// Delete 删除阅读历史
func (r *MongoReadingHistoryRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("ID不能为空")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除阅读历史失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("阅读历史不存在")
	}

	return nil
}

// =======================
// 查询方法
// =======================

// GetUserHistories 获取用户阅读历史列表（按时间倒序）
func (r *MongoReadingHistoryRepository) GetUserHistories(ctx context.Context, userID string, limit, offset int) ([]*reader.ReadingHistory, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}
	defer cursor.Close(ctx)

	var histories []*reader.ReadingHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, fmt.Errorf("解析阅读历史失败: %w", err)
	}

	if histories == nil {
		histories = []*reader.ReadingHistory{}
	}

	return histories, nil
}

// GetUserHistoriesByBook 获取用户指定书籍的阅读历史
func (r *MongoReadingHistoryRepository) GetUserHistoriesByBook(ctx context.Context, userID, bookID string, limit, offset int) ([]*reader.ReadingHistory, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return nil, fmt.Errorf("书籍ID不能为空")
	}

	filter := bson.M{
		"user_id": userID,
		"book_id": bookID,
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}
	defer cursor.Close(ctx)

	var histories []*reader.ReadingHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, fmt.Errorf("解析阅读历史失败: %w", err)
	}

	if histories == nil {
		histories = []*reader.ReadingHistory{}
	}

	return histories, nil
}

// GetUserHistoriesByTimeRange 获取指定时间范围内的阅读历史
func (r *MongoReadingHistoryRepository) GetUserHistoriesByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time, limit, offset int) ([]*reader.ReadingHistory, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	filter := bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}
	defer cursor.Close(ctx)

	var histories []*reader.ReadingHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, fmt.Errorf("解析阅读历史失败: %w", err)
	}

	if histories == nil {
		histories = []*reader.ReadingHistory{}
	}

	return histories, nil
}

// =======================
// 统计方法
// =======================

// CountUserHistories 统计用户历史记录数量
func (r *MongoReadingHistoryRepository) CountUserHistories(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}

	filter := bson.M{"user_id": userID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计阅读历史失败: %w", err)
	}

	return count, nil
}

// GetUserReadingStats 获取用户阅读统计
func (r *MongoReadingHistoryRepository) GetUserReadingStats(ctx context.Context, userID string) (*readerRepo.ReadingStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 使用聚合管道统计
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$group", Value: bson.M{
			"_id":            nil,
			"total_duration": bson.M{"$sum": "$read_duration"},
			"total_books":    bson.M{"$addToSet": "$book_id"},
			"total_chapters": bson.M{"$sum": 1},
			"last_read_time": bson.M{"$max": "$end_time"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":            0,
			"total_duration": 1,
			"total_books":    bson.M{"$size": "$total_books"},
			"total_chapters": 1,
			"last_read_time": 1,
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合查询失败: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalDuration int       `bson:"total_duration"`
		TotalBooks    int       `bson:"total_books"`
		TotalChapters int       `bson:"total_chapters"`
		LastReadTime  time.Time `bson:"last_read_time"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析统计结果失败: %w", err)
	}

	// 如果没有数据，返回空统计
	if len(results) == 0 {
		return &readerRepo.ReadingStats{
			TotalDuration:    0,
			TotalBooks:       0,
			TotalChapters:    0,
			LastReadTime:     time.Time{},
			AvgDailyDuration: 0,
		}, nil
	}

	result := results[0]

	// 计算日均阅读时长（基于有阅读记录的天数）
	avgDailyDuration := 0
	if !result.LastReadTime.IsZero() {
		// 获取第一条记录的时间
		var firstHistory reader.ReadingHistory
		opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: 1}})
		err := r.collection.FindOne(ctx, bson.M{"user_id": userID}, opts).Decode(&firstHistory)
		if err == nil {
			days := int(time.Since(firstHistory.CreatedAt).Hours() / 24)
			if days == 0 {
				days = 1
			}
			avgDailyDuration = result.TotalDuration / days
		}
	}

	return &readerRepo.ReadingStats{
		TotalDuration:    result.TotalDuration,
		TotalBooks:       result.TotalBooks,
		TotalChapters:    result.TotalChapters,
		LastReadTime:     result.LastReadTime,
		AvgDailyDuration: avgDailyDuration,
	}, nil
}

// GetUserDailyReadingStats 获取用户每日阅读统计
func (r *MongoReadingHistoryRepository) GetUserDailyReadingStats(ctx context.Context, userID string, days int) ([]readerRepo.DailyReadingStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if days <= 0 {
		days = 7 // 默认7天
	}

	// 计算起始日期
	startDate := time.Now().AddDate(0, 0, -days+1).Truncate(24 * time.Hour)

	// 聚合管道
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id":    userID,
			"created_at": bson.M{"$gte": startDate},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   "$created_at",
				},
			},
			"duration": bson.M{"$sum": "$read_duration"},
			"books":    bson.M{"$addToSet": "$book_id"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":      0,
			"date":     "$_id",
			"duration": 1,
			"books":    bson.M{"$size": "$books"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "date", Value: 1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合查询失败: %w", err)
	}
	defer cursor.Close(ctx)

	var stats []readerRepo.DailyReadingStats
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, fmt.Errorf("解析统计结果失败: %w", err)
	}

	if stats == nil {
		stats = []readerRepo.DailyReadingStats{}
	}

	return stats, nil
}

// =======================
// 批量操作
// =======================

// DeleteUserHistories 删除用户所有历史记录
func (r *MongoReadingHistoryRepository) DeleteUserHistories(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	filter := bson.M{"user_id": userID}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除用户历史记录失败: %w", err)
	}

	return nil
}

// DeleteOldHistories 删除指定日期之前的历史记录
func (r *MongoReadingHistoryRepository) DeleteOldHistories(ctx context.Context, beforeDate time.Time) (int64, error) {
	filter := bson.M{
		"created_at": bson.M{"$lt": beforeDate},
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("删除过期历史记录失败: %w", err)
	}

	return result.DeletedCount, nil
}

// =======================
// 健康检查
// =======================

// Health 健康检查
func (r *MongoReadingHistoryRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
