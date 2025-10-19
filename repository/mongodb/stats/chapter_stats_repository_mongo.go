package stats

import (
	"Qingyu_backend/models/stats"
	statsInterface "Qingyu_backend/repository/interfaces/stats"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoChapterStatsRepository MongoDB实现
type MongoChapterStatsRepository struct {
	collection *mongo.Collection
}

// NewMongoChapterStatsRepository 创建MongoDB Repository
func NewMongoChapterStatsRepository(db *mongo.Database) statsInterface.ChapterStatsRepository {
	return &MongoChapterStatsRepository{
		collection: db.Collection("chapter_stats"),
	}
}

// Create 创建章节统计
func (r *MongoChapterStatsRepository) Create(ctx context.Context, chapterStats *stats.ChapterStats) error {
	chapterStats.ID = primitive.NewObjectID().Hex()
	chapterStats.CreatedAt = time.Now()
	chapterStats.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, chapterStats)
	return err
}

// GetByID 根据ID获取
func (r *MongoChapterStatsRepository) GetByID(ctx context.Context, id string) (*stats.ChapterStats, error) {
	var chapterStats stats.ChapterStats
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&chapterStats)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &chapterStats, err
}

// Update 更新
func (r *MongoChapterStatsRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

// Delete 删除
func (r *MongoChapterStatsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetByChapterID 根据章节ID获取统计
func (r *MongoChapterStatsRepository) GetByChapterID(ctx context.Context, chapterID string) (*stats.ChapterStats, error) {
	var chapterStats stats.ChapterStats
	err := r.collection.FindOne(ctx, bson.M{"chapter_id": chapterID}).Decode(&chapterStats)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &chapterStats, err
}

// GetByBookID 根据作品ID获取章节统计列表
func (r *MongoChapterStatsRepository) GetByBookID(ctx context.Context, bookID string, limit, offset int64) ([]*stats.ChapterStats, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{"book_id": bookID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var chapterStats []*stats.ChapterStats
	if err = cursor.All(ctx, &chapterStats); err != nil {
		return nil, err
	}
	
	return chapterStats, nil
}

// GetByDateRange 根据日期范围获取统计
func (r *MongoChapterStatsRepository) GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.ChapterStats, error) {
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
	
	var chapterStats []*stats.ChapterStats
	if err = cursor.All(ctx, &chapterStats); err != nil {
		return nil, err
	}
	
	return chapterStats, nil
}

// GetChapterStatsAggregate 获取章节统计聚合
func (r *MongoChapterStatsRepository) GetChapterStatsAggregate(ctx context.Context, bookID string) ([]*stats.ChapterStatsAggregate, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$chapter_id",
			"title": bson.M{"$first": "$title"},
			"view_count": bson.M{"$sum": "$view_count"},
			"unique_viewers": bson.M{"$sum": "$unique_viewers"},
			"completion_rate": bson.M{"$avg": "$completion_rate"},
			"drop_off_rate": bson.M{"$avg": "$drop_off_rate"},
			"revenue": bson.M{"$sum": "$revenue"},
		}}},
		{{Key: "$project", Value: bson.M{
			"chapter_id": "$_id",
			"title": 1,
			"view_count": 1,
			"unique_viewers": 1,
			"completion_rate": 1,
			"drop_off_rate": 1,
			"revenue": 1,
			"_id": 0,
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*stats.ChapterStatsAggregate
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// GetTopViewedChapters 获取阅读量最高的章节
func (r *MongoChapterStatsRepository) GetTopViewedChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$sort", Value: bson.D{{Key: "view_count", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: bson.M{
			"chapter_id": 1,
			"title": 1,
			"view_count": 1,
			"unique_viewers": 1,
			"completion_rate": 1,
			"drop_off_rate": 1,
			"revenue": 1,
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*stats.ChapterStatsAggregate
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// GetTopRevenueChapters 获取收入最高的章节
func (r *MongoChapterStatsRepository) GetTopRevenueChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$sort", Value: bson.D{{Key: "revenue", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: bson.M{
			"chapter_id": 1,
			"title": 1,
			"view_count": 1,
			"unique_viewers": 1,
			"completion_rate": 1,
			"drop_off_rate": 1,
			"revenue": 1,
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*stats.ChapterStatsAggregate
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// GetLowestCompletionChapters 获取完读率最低的章节
func (r *MongoChapterStatsRepository) GetLowestCompletionChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$sort", Value: bson.D{{Key: "completion_rate", Value: 1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: bson.M{
			"chapter_id": 1,
			"title": 1,
			"view_count": 1,
			"unique_viewers": 1,
			"completion_rate": 1,
			"drop_off_rate": 1,
			"revenue": 1,
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*stats.ChapterStatsAggregate
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// GetHighestDropOffChapters 获取跳出率最高的章节
func (r *MongoChapterStatsRepository) GetHighestDropOffChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$sort", Value: bson.D{{Key: "drop_off_rate", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: bson.M{
			"chapter_id": 1,
			"title": 1,
			"view_count": 1,
			"unique_viewers": 1,
			"completion_rate": 1,
			"drop_off_rate": 1,
			"revenue": 1,
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*stats.ChapterStatsAggregate
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// GenerateHeatmap 生成热力图
func (r *MongoChapterStatsRepository) GenerateHeatmap(ctx context.Context, bookID string) ([]*stats.HeatmapPoint, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"book_id": bookID}}},
		{{Key: "$sort", Value: bson.D{{Key: "chapter_num", Value: 1}}}},
		{{Key: "$project", Value: bson.M{
			"chapter_num": 1,
			"chapter_id": 1,
			"view_count": 1,
			"completion_rate": 1,
			"drop_off_rate": 1,
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*stats.HeatmapPoint
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// GetTimeRangeStats 获取时间范围统计
func (r *MongoChapterStatsRepository) GetTimeRangeStats(ctx context.Context, bookID string, startDate, endDate time.Time) (*stats.TimeRangeStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"book_id": bookID,
			"stat_date": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"total_views": bson.M{"$sum": "$view_count"},
			"total_unique_viewers": bson.M{"$sum": "$unique_viewers"},
			"avg_completion_rate": bson.M{"$avg": "$completion_rate"},
			"avg_drop_off_rate": bson.M{"$avg": "$drop_off_rate"},
			"total_revenue": bson.M{"$sum": "$revenue"},
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []stats.TimeRangeStats
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return &stats.TimeRangeStats{
			StartDate: startDate,
			EndDate:   endDate,
		}, nil
	}
	
	result := &results[0]
	result.StartDate = startDate
	result.EndDate = endDate
	
	return result, nil
}

// BatchCreate 批量创建
func (r *MongoChapterStatsRepository) BatchCreate(ctx context.Context, chapterStats []*stats.ChapterStats) error {
	if len(chapterStats) == 0 {
		return nil
	}
	
	docs := make([]interface{}, len(chapterStats))
	for i, stat := range chapterStats {
		stat.ID = primitive.NewObjectID().Hex()
		stat.CreatedAt = time.Now()
		stat.UpdatedAt = time.Now()
		docs[i] = stat
	}
	
	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

// BatchUpdate 批量更新
func (r *MongoChapterStatsRepository) BatchUpdate(ctx context.Context, updates []map[string]interface{}) error {
	// TODO: 实现批量更新
	return fmt.Errorf("未实现")
}

// Count 统计数量
func (r *MongoChapterStatsRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	// 简化实现，使用空filter
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByBook 按作品统计
func (r *MongoChapterStatsRepository) CountByBook(ctx context.Context, bookID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"book_id": bookID})
}

// Health 健康检查
func (r *MongoChapterStatsRepository) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}

