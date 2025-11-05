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

// MongoReaderBehaviorRepository MongoDB实现
type MongoReaderBehaviorRepository struct {
	collection          *mongo.Collection
	sessionCollection   *mongo.Collection
	retentionCollection *mongo.Collection
}

// NewMongoReaderBehaviorRepository 创建MongoDB Repository
func NewMongoReaderBehaviorRepository(db *mongo.Database) statsInterface.ReaderBehaviorRepository {
	return &MongoReaderBehaviorRepository{
		collection:          db.Collection("reader_behaviors"),
		sessionCollection:   db.Collection("reading_sessions"),
		retentionCollection: db.Collection("reader_retentions"),
	}
}

// Create 创建读者行为记录
func (r *MongoReaderBehaviorRepository) Create(ctx context.Context, behavior *stats.ReaderBehavior) error {
	behavior.ID = primitive.NewObjectID().Hex()
	behavior.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, behavior)
	return err
}

// GetByID 根据ID获取
func (r *MongoReaderBehaviorRepository) GetByID(ctx context.Context, id string) (*stats.ReaderBehavior, error) {
	var behavior stats.ReaderBehavior
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&behavior)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &behavior, err
}

// Delete 删除
func (r *MongoReaderBehaviorRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetByUserID 根据用户ID获取行为列表
func (r *MongoReaderBehaviorRepository) GetByUserID(ctx context.Context, userID string, limit, offset int64) ([]*stats.ReaderBehavior, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "read_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []*stats.ReaderBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// GetByBookID 根据作品ID获取行为列表
func (r *MongoReaderBehaviorRepository) GetByBookID(ctx context.Context, bookID string, limit, offset int64) ([]*stats.ReaderBehavior, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "read_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"book_id": bookID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []*stats.ReaderBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// GetByChapterID 根据章节ID获取行为列表
func (r *MongoReaderBehaviorRepository) GetByChapterID(ctx context.Context, chapterID string, limit, offset int64) ([]*stats.ReaderBehavior, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "read_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"chapter_id": chapterID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []*stats.ReaderBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// GetByBehaviorType 根据行为类型获取
func (r *MongoReaderBehaviorRepository) GetByBehaviorType(ctx context.Context, behaviorType string, limit, offset int64) ([]*stats.ReaderBehavior, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "read_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"behavior_type": behaviorType}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []*stats.ReaderBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// GetByDateRange 根据日期范围查询
func (r *MongoReaderBehaviorRepository) GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.ReaderBehavior, error) {
	filter := bson.M{
		"book_id": bookID,
		"read_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []*stats.ReaderBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// CountUniqueReaders 统计独立读者数
func (r *MongoReaderBehaviorRepository) CountUniqueReaders(ctx context.Context, bookID string) (int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$user_id",
		}}},
		{{Key: "$count", Value: "unique_readers"}},
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

	count, _ := result[0]["unique_readers"].(int32)
	return int64(count), nil
}

// CountUniqueReadersByChapter 统计章节独立读者数
func (r *MongoReaderBehaviorRepository) CountUniqueReadersByChapter(ctx context.Context, chapterID string) (int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"chapter_id": chapterID}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$user_id",
		}}},
		{{Key: "$count", Value: "unique_readers"}},
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

	count, _ := result[0]["unique_readers"].(int32)
	return int64(count), nil
}

// CalculateAvgReadTime 计算平均阅读时长
func (r *MongoReaderBehaviorRepository) CalculateAvgReadTime(ctx context.Context, chapterID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"chapter_id": chapterID}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"avg_read_time": bson.M{"$avg": "$read_duration"},
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

	avgTime, _ := result[0]["avg_read_time"].(float64)
	return avgTime, nil
}

// CalculateCompletionRate 计算完读率
func (r *MongoReaderBehaviorRepository) CalculateCompletionRate(ctx context.Context, chapterID string) (float64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{"chapter_id": chapterID})
	if err != nil {
		return 0, err
	}

	if totalCount == 0 {
		return 0, nil
	}

	completeCount, err := r.collection.CountDocuments(ctx, bson.M{
		"chapter_id":    chapterID,
		"behavior_type": stats.BehaviorTypeComplete,
	})
	if err != nil {
		return 0, err
	}

	return float64(completeCount) / float64(totalCount), nil
}

// CalculateDropOffRate 计算跳出率
func (r *MongoReaderBehaviorRepository) CalculateDropOffRate(ctx context.Context, chapterID string) (float64, error) {
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{"chapter_id": chapterID})
	if err != nil {
		return 0, err
	}

	if totalCount == 0 {
		return 0, nil
	}

	dropOffCount, err := r.collection.CountDocuments(ctx, bson.M{
		"chapter_id":    chapterID,
		"behavior_type": stats.BehaviorTypeDropOff,
	})
	if err != nil {
		return 0, err
	}

	return float64(dropOffCount) / float64(totalCount), nil
}

// CreateSession 创建阅读会话
func (r *MongoReaderBehaviorRepository) CreateSession(ctx context.Context, session *stats.ReadingSession) error {
	_, err := r.sessionCollection.InsertOne(ctx, session)
	return err
}

// GetSessionByID 根据ID获取会话
func (r *MongoReaderBehaviorRepository) GetSessionByID(ctx context.Context, sessionID string) (*stats.ReadingSession, error) {
	var session stats.ReadingSession
	err := r.sessionCollection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &session, err
}

// GetUserSessions 获取用户会话列表
func (r *MongoReaderBehaviorRepository) GetUserSessions(ctx context.Context, userID string, limit int) ([]*stats.ReadingSession, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "start_time", Value: -1}})

	cursor, err := r.sessionCollection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*stats.ReadingSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// CreateRetention 创建留存数据
func (r *MongoReaderBehaviorRepository) CreateRetention(ctx context.Context, retention *stats.ReaderRetention) error {
	retention.CreatedAt = time.Now()
	_, err := r.retentionCollection.InsertOne(ctx, retention)
	return err
}

// GetRetentionByBookID 获取作品留存数据
func (r *MongoReaderBehaviorRepository) GetRetentionByBookID(ctx context.Context, bookID string) (*stats.ReaderRetention, error) {
	var retention stats.ReaderRetention
	opts := options.FindOne().SetSort(bson.D{{Key: "stat_date", Value: -1}})
	err := r.retentionCollection.FindOne(ctx, bson.M{"book_id": bookID}, opts).Decode(&retention)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &retention, err
}

// CalculateRetention 计算留存率
func (r *MongoReaderBehaviorRepository) CalculateRetention(ctx context.Context, bookID string, days int) (float64, error) {
	// 简化实现：计算N天前的读者在今天还活跃的比例
	targetDate := time.Now().AddDate(0, 0, -days)

	// 统计N天前的读者数
	targetReaders, err := r.collection.Distinct(ctx, "user_id", bson.M{
		"book_id": bookID,
		"read_at": bson.M{
			"$gte": targetDate,
			"$lt":  targetDate.Add(24 * time.Hour),
		},
	})
	if err != nil {
		return 0, err
	}

	if len(targetReaders) == 0 {
		return 0, nil
	}

	// 统计这些读者今天还活跃的数量
	today := time.Now().Truncate(24 * time.Hour)
	activeCount, err := r.collection.CountDocuments(ctx, bson.M{
		"book_id": bookID,
		"user_id": bson.M{"$in": targetReaders},
		"read_at": bson.M{"$gte": today},
	})
	if err != nil {
		return 0, err
	}

	return float64(activeCount) / float64(len(targetReaders)), nil
}

// BatchCreate 批量创建
func (r *MongoReaderBehaviorRepository) BatchCreate(ctx context.Context, behaviors []*stats.ReaderBehavior) error {
	if len(behaviors) == 0 {
		return nil
	}

	docs := make([]interface{}, len(behaviors))
	for i, behavior := range behaviors {
		behavior.ID = primitive.NewObjectID().Hex()
		behavior.CreatedAt = time.Now()
		docs[i] = behavior
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

// Count 统计数量
func (r *MongoReaderBehaviorRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByBehaviorType 按行为类型统计
func (r *MongoReaderBehaviorRepository) CountByBehaviorType(ctx context.Context, behaviorType string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"behavior_type": behaviorType})
}

// Health 健康检查
func (r *MongoReaderBehaviorRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
