package mongodb

import (
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared/types"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	infra "Qingyu_backend/repository/interfaces/infrastructure"
)

// MongoBookStatisticsRepository MongoDB书籍统计仓储实现
type MongoBookStatisticsRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBookStatisticsRepository 创建MongoDB书籍统计仓储实例
func NewMongoBookStatisticsRepository(client *mongo.Client, database string) BookstoreInterface.BookStatisticsRepository {
	return &MongoBookStatisticsRepository{
		collection: client.Database(database).Collection("book_statistics"),
		client:     client,
	}
}

// Create 创建书籍统计
func (r *MongoBookStatisticsRepository) Create(ctx context.Context, statistics *bookstore.BookStatistics) error {
	if statistics == nil {
		return errors.New("book statistics cannot be nil")
	}

	statistics.BeforeUpdate()

	result, err := r.collection.InsertOne(ctx, statistics)
	if err != nil {
		return err
	}

	// ID 是 string 类型，需要转换
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		statistics.ID = oid.Hex()
	}
	return nil
}

// GetByID 根据ID获取书籍统计
func (r *MongoBookStatisticsRepository) GetByID(ctx context.Context, id string) (*bookstore.BookStatistics, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid statistics ID: %w", err)
	}
	var statistics bookstore.BookStatistics
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&statistics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &statistics, nil
}

// Update 更新书籍统计（符合基础CRUD接口）
func (r *MongoBookStatisticsRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid statistics ID: %w", err)
	}
	if len(updates) == 0 {
		return errors.New("updates cannot be empty")
	}
	updates["updated_at"] = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("book statistics not found")
	}
	return nil
}

// Delete 删除书籍统计
func (r *MongoBookStatisticsRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("book statistics not found")
	}

	return nil
}

// GetAll 获取所有书籍统计
func (r *MongoBookStatisticsRepository) GetAll(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "hot_score", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// Count 统计书籍统计总数
func (r *MongoBookStatisticsRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.collection.CountDocuments(ctx, query)
}

// List 根据过滤条件列出统计
func (r *MongoBookStatisticsRepository) List(ctx context.Context, filter infra.Filter) ([]*bookstore.BookStatistics, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	opts := options.Find()
	if filter != nil {
		sort := filter.GetSort()
		if len(sort) > 0 {
			var sortDoc bson.D
			for k, v := range sort {
				sortDoc = append(sortDoc, bson.E{Key: k, Value: v})
			}
			opts.SetSort(sortDoc)
		}
	}
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*bookstore.BookStatistics
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Exists 判断统计是否存在
func (r *MongoBookStatisticsRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid book ID: %w", err)
	}
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByBookID 根据书籍ID获取统计信息
func (r *MongoBookStatisticsRepository) GetByBookID(ctx context.Context, bookID string) (*bookstore.BookStatistics, error) {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}
	var statistics bookstore.BookStatistics
	filter := bson.M{"book_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&statistics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &statistics, nil
}

// GetTopViewed 获取浏览量最高的书籍统计
func (r *MongoBookStatisticsRepository) GetTopViewed(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "view_count", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// GetTopFavorited 获取收藏量最高的书籍统计
func (r *MongoBookStatisticsRepository) GetTopFavorited(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "favorite_count", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// GetMostShared 获取分享最多的书籍统计
func (r *MongoBookStatisticsRepository) GetMostShared(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "share_count", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// GetTopRated 获取评分最高的书籍统计
func (r *MongoBookStatisticsRepository) GetTopRated(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "average_rating", Value: -1}})

	filter := bson.M{"rating_count": bson.M{"$gte": 5}} // 至少5个评分
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// GetMostCommented 获取评论最多的书籍统计
func (r *MongoBookStatisticsRepository) GetMostCommented(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "comment_count", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// GetHottest 获取热度最高的书籍统计
func (r *MongoBookStatisticsRepository) GetHottest(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "hot_score", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// Search 搜索书籍统计
func (r *MongoBookStatisticsRepository) Search(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookStatistics, int64, error) {
	// 计算总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// 设置分页选项
	opts := options.Find()
	if pageSize > 0 {
		opts.SetLimit(int64(pageSize))
		if page > 0 {
			opts.SetSkip(int64((page - 1) * pageSize))
		}
	}
	opts.SetSort(bson.D{{Key: "hot_score", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, 0, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, total, cursor.Err()
}

// IncrementViewCount 增加浏览量
func (r *MongoBookStatisticsRepository) IncrementViewCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$inc": bson.M{"view_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// IncrementFavoriteCount 增加收藏量
func (r *MongoBookStatisticsRepository) IncrementFavoriteCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$inc": bson.M{"favorite_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// DecrementFavoriteCount 减少收藏量
func (r *MongoBookStatisticsRepository) DecrementFavoriteCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$inc": bson.M{"favorite_count": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// IncrementCommentCount 增加评论量
func (r *MongoBookStatisticsRepository) IncrementCommentCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$inc": bson.M{"comment_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// DecrementCommentCount 减少评论量
func (r *MongoBookStatisticsRepository) DecrementCommentCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$inc": bson.M{"comment_count": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// IncrementShareCount 增加分享量
func (r *MongoBookStatisticsRepository) IncrementShareCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$inc": bson.M{"share_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateRating 更新评分统计
func (r *MongoBookStatisticsRepository) UpdateRatingSummary(ctx context.Context, bookID primitive.ObjectID, averageRating float64, ratingCount int64, ratingDistribution map[int]int64) error {
	filter := bson.M{"book_id": bookID}
	update := bson.M{
		"$set": bson.M{
			"average_rating":      averageRating,
			"rating_count":        ratingCount,
			"rating_distribution": ratingDistribution,
			"updated_at":          time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateRating 根据新增评分更新统计（按增量更新平均值与计数）
func (r *MongoBookStatisticsRepository) UpdateRating(ctx context.Context, bookID string, rating int) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	// 读取当前统计
	var stats bookstore.BookStatistics
	err = r.collection.FindOne(ctx, bson.M{"book_id": objectID}).Decode(&stats)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		// 初始创建
		stats = bookstore.BookStatistics{BookID: bookID, AverageRating: 0, RatingCount: 0}
	}

	newCount := stats.RatingCount + 1
	// Rating 类型需要转换为 float64 进行运算
	newAvg := ((float64(stats.AverageRating) * float64(stats.RatingCount)) + float64(rating)) / float64(newCount)

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"book_id": bookID},
		bson.M{"$set": bson.M{"average_rating": types.Rating(newAvg), "rating_count": newCount, "updated_at": time.Now()}},
		options.Update().SetUpsert(true),
	)
	return err
}

// RemoveRating 根据移除评分更新统计
func (r *MongoBookStatisticsRepository) RemoveRating(ctx context.Context, bookID string, rating int) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	var stats bookstore.BookStatistics
	err = r.collection.FindOne(ctx, bson.M{"book_id": objectID}).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}

	if stats.RatingCount <= 0 {
		return nil
	}

	newCount := stats.RatingCount - 1
	var newAvg float64
	if newCount > 0 {
		// Rating 类型需要转换为 float64 进行运算
		newAvg = ((float64(stats.AverageRating) * float64(stats.RatingCount)) - float64(rating)) / float64(newCount)
	} else {
		newAvg = 0
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"book_id": bookID},
		bson.M{"$set": bson.M{"average_rating": types.Rating(newAvg), "rating_count": newCount, "updated_at": time.Now()}},
	)
	return err
}

// UpdateHotScore 更新热度分数
func (r *MongoBookStatisticsRepository) UpdateHotScore(ctx context.Context, bookID string, hotScore float64) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"book_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"hot_score":  hotScore,
			"updated_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// BatchUpdateHotScore 批量更新热度分数
func (r *MongoBookStatisticsRepository) BatchUpdateHotScore(ctx context.Context, bookIDs []string) error {
	// 转换 string ID 为 ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("无效的ID: %s", id)
		}
		objectIDs = append(objectIDs, oid)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchIncrementViewCount 批量增加浏览量
func (r *MongoBookStatisticsRepository) BatchIncrementViewCount(ctx context.Context, increments map[string]int64) error {
	var operations []mongo.WriteModel

	for bookIDStr, count := range increments {
		bookID, err := primitive.ObjectIDFromHex(bookIDStr)
		if err != nil {
			return fmt.Errorf("无效的ID: %s", bookIDStr)
		}

		filter := bson.M{"book_id": bookID}
		update := bson.M{
			"$inc": bson.M{"view_count": count},
			"$set": bson.M{"updated_at": time.Now()},
		}

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(filter)
		operation.SetUpdate(update)
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	if len(operations) == 0 {
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, operations)
	return err
}

// GetAggregatedStatistics 获取聚合统计信息
func (r *MongoBookStatisticsRepository) GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":               nil,
			"total_books":       bson.M{"$sum": 1},
			"total_views":       bson.M{"$sum": "$view_count"},
			"total_favorites":   bson.M{"$sum": "$favorite_count"},
			"total_comments":    bson.M{"$sum": "$comment_count"},
			"total_shares":      bson.M{"$sum": "$share_count"},
			"average_rating":    bson.M{"$avg": "$average_rating"},
			"average_hot_score": bson.M{"$avg": "$hot_score"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result map[string]interface{}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		delete(result, "_id") // 移除_id字段
		return result, nil
	}

	return make(map[string]interface{}), nil
}

// GetTotalViews 获取总浏览量
func (r *MongoBookStatisticsRepository) GetTotalViews(ctx context.Context) (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":         nil,
			"total_views": bson.M{"$sum": "$view_count"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalViews int64 `bson:"total_views"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.TotalViews, nil
	}

	return 0, nil
}

// GetTotalFavorites 获取总收藏量
func (r *MongoBookStatisticsRepository) GetTotalFavorites(ctx context.Context) (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":             nil,
			"total_favorites": bson.M{"$sum": "$favorite_count"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalFavorites int64 `bson:"total_favorites"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.TotalFavorites, nil
	}

	return 0, nil
}

// GetTotalComments 获取总评论量
func (r *MongoBookStatisticsRepository) GetTotalComments(ctx context.Context) (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":            nil,
			"total_comments": bson.M{"$sum": "$comment_count"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalComments int64 `bson:"total_comments"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.TotalComments, nil
	}

	return 0, nil
}

// GetTotalShares 获取总分享量
func (r *MongoBookStatisticsRepository) GetTotalShares(ctx context.Context) (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":          nil,
			"total_shares": bson.M{"$sum": "$share_count"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalShares int64 `bson:"total_shares"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.TotalShares, nil
	}

	return 0, nil
}

// GetAverageRating 获取平均评分
func (r *MongoBookStatisticsRepository) GetAverageRating(ctx context.Context) (float64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"rating_count": bson.M{"$gt": 0},
		}},
		{"$group": bson.M{
			"_id":            nil,
			"average_rating": bson.M{"$avg": "$average_rating"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		AverageRating float64 `bson:"average_rating"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.AverageRating, nil
	}

	return 0, nil
}

// GetViewsInRange 获取指定时间范围内的浏览量
func (r *MongoBookStatisticsRepository) GetViewsInRange(ctx context.Context, bookID string, days int) (int64, error) {
	// 由于当前模型没有时间范围的浏览量记录，这里返回当前的浏览量
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return 0, fmt.Errorf("invalid book ID: %w", err)
	}
	stats, err := r.GetByBookID(ctx, objectID.Hex())
	if err != nil {
		return 0, err
	}
	if stats == nil {
		return 0, nil
	}
	return stats.ViewCount, nil
}

// GetFavoritesInRange 获取指定时间范围内的收藏量
func (r *MongoBookStatisticsRepository) GetFavoritesInRange(ctx context.Context, bookID string, days int) (int64, error) {
	// 由于当前模型没有时间范围的收藏量记录，这里返回当前的收藏量
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return 0, fmt.Errorf("invalid book ID: %w", err)
	}
	stats, err := r.GetByBookID(ctx, objectID.Hex())
	if err != nil {
		return 0, err
	}
	if stats == nil {
		return 0, nil
	}
	return stats.FavoriteCount, nil
}

// GetCommentsInRange 获取指定时间范围内的评论量
func (r *MongoBookStatisticsRepository) GetCommentsInRange(ctx context.Context, bookID string, days int) (int64, error) {
	// 由于当前模型没有时间范围的评论量记录，这里返回当前的评论量
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return 0, fmt.Errorf("invalid book ID: %w", err)
	}
	stats, err := r.GetByBookID(ctx, objectID.Hex())
	if err != nil {
		return 0, err
	}
	if stats == nil {
		return 0, nil
	}
	return stats.CommentCount, nil
}

// GetStatisticsByTimeRange 根据时间范围获取统计信息
func (r *MongoBookStatisticsRepository) GetStatisticsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*bookstore.BookStatistics, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	filter := bson.M{
		"updated_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// GetTrendingBooks 获取趋势书籍（基于最近的统计变化）
func (r *MongoBookStatisticsRepository) GetTrendingBooks(ctx context.Context, days int, limit, offset int) ([]*bookstore.BookStatistics, error) {
	// 计算时间范围
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "hot_score", Value: -1}})

	filter := bson.M{
		"updated_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, cursor.Err()
}

// BatchRecalculateStatistics 批量重新计算统计数据
func (r *MongoBookStatisticsRepository) BatchRecalculateStatistics(ctx context.Context, bookIDs []string) error {
	// 转换 string ID 为 ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("无效的ID: %s", id)
		}
		objectIDs = append(objectIDs, oid)
	}

	// 为每个书籍重新计算统计数据
	for _, bookID := range objectIDs {
		// 这里可以根据实际需求实现统计数据的重新计算逻辑
		// 例如：重新计算浏览量、收藏量、评分等
		filter := bson.M{"book_id": bookID}
		update := bson.M{
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		}

		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}

// Transaction 执行事务
func (r *MongoBookStatisticsRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := fn(sc); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})
}

// Health 健康检查
func (r *MongoBookStatisticsRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// BatchUpdateViewCount 批量更新浏览量
func (r *MongoBookStatisticsRepository) BatchUpdateViewCount(ctx context.Context, bookIDs []string, increment int64) error {
	var operations []mongo.WriteModel

	for _, bookID := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(bookID)
		if err != nil {
			return fmt.Errorf("invalid book ID: %w", err)
		}

		filter := bson.M{"book_id": objectID}
		update := bson.M{
			"$inc": bson.M{"view_count": increment},
			"$set": bson.M{"updated_at": time.Now()},
		}

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(filter)
		operation.SetUpdate(update)
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	if len(operations) == 0 {
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, operations)
	return err
}

// SearchByFilter 根据过滤器搜索（带分页）
func (r *MongoBookStatisticsRepository) SearchByFilter(ctx context.Context, filter *BookstoreInterface.BookStatisticsFilter, page, pageSize int) ([]*bookstore.BookStatistics, int64, error) {
	// 构建查询条件
	query := bson.M{}

	if filter.BookID != nil {
		query["book_id"] = *filter.BookID
	}
	if filter.MinViewCount != nil {
		if query["view_count"] == nil {
			query["view_count"] = bson.M{}
		}
		query["view_count"].(bson.M)["$gte"] = *filter.MinViewCount
	}
	if filter.MaxViewCount != nil {
		if query["view_count"] == nil {
			query["view_count"] = bson.M{}
		}
		query["view_count"].(bson.M)["$lte"] = *filter.MaxViewCount
	}
	// 可以继续添加其他过滤条件...

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 设置分页选项
	opts := options.Find()
	if pageSize > 0 {
		opts.SetLimit(int64(pageSize))
		if page > 0 {
			opts.SetSkip(int64((page - 1) * pageSize))
		}
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var statisticsList []*bookstore.BookStatistics
	for cursor.Next(ctx) {
		var statistics bookstore.BookStatistics
		if err := cursor.Decode(&statistics); err != nil {
			return nil, 0, err
		}
		statisticsList = append(statisticsList, &statistics)
	}

	return statisticsList, total, cursor.Err()
}
