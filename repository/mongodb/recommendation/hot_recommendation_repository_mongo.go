package recommendation

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
)

// MongoHotRecommendationRepository 热门推荐Repository的MongoDB实现
// 基于book_statistics集合的统计数据实现热门推荐
type MongoHotRecommendationRepository struct {
	statisticsCollection *mongo.Collection // book_statistics集合
	booksCollection      *mongo.Collection // books集合
}

// NewMongoHotRecommendationRepository 创建MongoHotRecommendationRepository实例
func NewMongoHotRecommendationRepository(db *mongo.Database) recoRepo.HotRecommendationRepository {
	return &MongoHotRecommendationRepository{
		statisticsCollection: db.Collection("book_statistics"),
		booksCollection:      db.Collection("books"),
	}
}

// GetHotBooks 获取热门书籍列表
// 基于浏览量、收藏量、评分综合计算热度分数
func (r *MongoHotRecommendationRepository) GetHotBooks(ctx context.Context, limit int, days int) ([]string, error) {
	// 计算时间阈值
	timeThreshold := time.Now().AddDate(0, 0, -days)

	// 聚合管道：计算热度分数
	pipeline := mongo.Pipeline{
		// 筛选最近N天的数据
		bson.D{{Key: "$match", Value: bson.M{
			"updated_at": bson.M{"$gte": timeThreshold},
		}}},
		// 计算热度分数
		bson.D{{Key: "$addFields", Value: bson.M{
			"hot_score": bson.M{
				"$add": []interface{}{
					bson.M{"$multiply": []interface{}{"$views", 0.2}},         // 浏览量权重 0.2
					bson.M{"$multiply": []interface{}{"$favorites", 0.7}},     // 收藏量权重 0.7（收藏行为更重要）
					bson.M{"$multiply": []interface{}{"$average_rating", 20}}, // 评分权重（转换到相同量级）
				},
			},
		}}},
		// 按热度分数排序
		bson.D{{Key: "$sort", Value: bson.M{"hot_score": -1}}},
		// 限制返回数量
		bson.D{{Key: "$limit", Value: limit}},
		// 只返回book_id
		bson.D{{Key: "$project", Value: bson.M{"book_id": 1, "_id": 0}}},
	}

	cursor, err := r.statisticsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get hot books: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		BookID string `bson:"book_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode hot books: %w", err)
	}

	// 提取book_id列表
	bookIDs := make([]string, len(results))
	for i, result := range results {
		bookIDs[i] = result.BookID
	}

	return bookIDs, nil
}

// GetHotBooksByCategory 获取分类下的热门书籍
func (r *MongoHotRecommendationRepository) GetHotBooksByCategory(ctx context.Context, category string, limit int, days int) ([]string, error) {
	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	// 计算时间阈值
	timeThreshold := time.Now().AddDate(0, 0, -days)

	// 聚合管道：联合查询books和book_statistics
	pipeline := mongo.Pipeline{
		// 从books集合开始，筛选分类
		bson.D{{Key: "$match", Value: bson.M{
			"category": category,
			"status":   "published", // 只推荐已发布的书籍
		}}},
		// 关联book_statistics集合
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "book_statistics",
			"localField":   "_id",
			"foreignField": "book_id",
			"as":           "statistics",
		}}},
		// 展开statistics数组
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$statistics",
			"preserveNullAndEmptyArrays": true,
		}}},
		// 筛选最近N天的数据
		bson.D{{Key: "$match", Value: bson.M{
			"statistics.updated_at": bson.M{"$gte": timeThreshold},
		}}},
		// 计算热度分数
		bson.D{{Key: "$addFields", Value: bson.M{
			"hot_score": bson.M{
				"$add": []interface{}{
					bson.M{"$multiply": []interface{}{bson.M{"$ifNull": []interface{}{"$statistics.views", 0}}, 0.3}},
					bson.M{"$multiply": []interface{}{bson.M{"$ifNull": []interface{}{"$statistics.favorites", 0}}, 0.5}},
					bson.M{"$multiply": []interface{}{bson.M{"$ifNull": []interface{}{"$statistics.average_rating", 0}}, 20}},
				},
			},
		}}},
		// 按热度分数排序
		bson.D{{Key: "$sort", Value: bson.M{"hot_score": -1}}},
		// 限制返回数量
		bson.D{{Key: "$limit", Value: limit}},
		// 只返回book_id
		bson.D{{Key: "$project", Value: bson.M{"_id": 1}}},
	}

	cursor, err := r.booksCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get hot books by category: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode hot books by category: %w", err)
	}

	// 提取book_id列表
	bookIDs := make([]string, len(results))
	for i, result := range results {
		bookIDs[i] = result.ID
	}

	return bookIDs, nil
}

// GetTrendingBooks 获取正在飙升的书籍（增长趋势）
// 基于最近7天和前7天的数据对比，计算增长率
func (r *MongoHotRecommendationRepository) GetTrendingBooks(ctx context.Context, limit int) ([]string, error) {
	// 简化实现：使用最近3天的数据作为"飙升"标准
	// 后续可以优化为真正的增长率计算

	recentThreshold := time.Now().AddDate(0, 0, -3)

	pipeline := mongo.Pipeline{
		// 筛选最近3天的数据
		bson.D{{Key: "$match", Value: bson.M{
			"updated_at": bson.M{"$gte": recentThreshold},
		}}},
		// 计算趋势分数（最近浏览量+收藏量）
		bson.D{{Key: "$addFields", Value: bson.M{
			"trend_score": bson.M{
				"$add": []interface{}{
					"$views",
					bson.M{"$multiply": []interface{}{"$favorites", 2}}, // 收藏权重更高
				},
			},
		}}},
		// 按趋势分数排序
		bson.D{{Key: "$sort", Value: bson.M{"trend_score": -1}}},
		// 限制返回数量
		bson.D{{Key: "$limit", Value: limit}},
		// 只返回book_id
		bson.D{{Key: "$project", Value: bson.M{"book_id": 1, "_id": 0}}},
	}

	cursor, err := r.statisticsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending books: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		BookID string `bson:"book_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode trending books: %w", err)
	}

	// 提取book_id列表
	bookIDs := make([]string, len(results))
	for i, result := range results {
		bookIDs[i] = result.BookID
	}

	return bookIDs, nil
}

// GetNewPopularBooks 获取新书中的热门书籍
func (r *MongoHotRecommendationRepository) GetNewPopularBooks(ctx context.Context, limit int, daysThreshold int) ([]string, error) {
	// 计算新书时间阈值
	timeThreshold := time.Now().AddDate(0, 0, -daysThreshold)

	// 聚合管道
	pipeline := mongo.Pipeline{
		// 筛选新书（最近N天上架）
		bson.D{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": timeThreshold},
			"status":     "published",
		}}},
		// 关联book_statistics集合
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "book_statistics",
			"localField":   "_id",
			"foreignField": "book_id",
			"as":           "statistics",
		}}},
		// 展开statistics数组
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$statistics",
			"preserveNullAndEmptyArrays": true,
		}}},
		// 计算新书热度分数
		bson.D{{Key: "$addFields", Value: bson.M{
			"new_book_score": bson.M{
				"$add": []interface{}{
					bson.M{"$multiply": []interface{}{bson.M{"$ifNull": []interface{}{"$statistics.views", 0}}, 0.4}},
					bson.M{"$multiply": []interface{}{bson.M{"$ifNull": []interface{}{"$statistics.favorites", 0}}, 0.6}},
				},
			},
		}}},
		// 按热度分数排序
		bson.D{{Key: "$sort", Value: bson.M{"new_book_score": -1}}},
		// 限制返回数量
		bson.D{{Key: "$limit", Value: limit}},
		// 只返回book_id
		bson.D{{Key: "$project", Value: bson.M{"_id": 1}}},
	}

	cursor, err := r.booksCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get new popular books: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode new popular books: %w", err)
	}

	// 提取book_id列表
	bookIDs := make([]string, len(results))
	for i, result := range results {
		bookIDs[i] = result.ID
	}

	return bookIDs, nil
}

// Health 健康检查
func (r *MongoHotRecommendationRepository) Health(ctx context.Context) error {
	return r.statisticsCollection.Database().Client().Ping(ctx, nil)
}
