package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/reading/bookstore"
	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
)

// MongoBookRatingRepository MongoDB书籍评分仓储实现
type MongoBookRatingRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBookRatingRepository 创建MongoDB书籍评分仓储实例
func NewMongoBookRatingRepository(client *mongo.Client, database string) BookstoreInterface.BookRatingRepository {
	return &MongoBookRatingRepository{
		collection: client.Database(database).Collection("book_ratings"),
		client:     client,
	}
}

// Create 创建书籍评分
func (r *MongoBookRatingRepository) Create(ctx context.Context, rating *bookstore.BookRating) error {
	if rating == nil {
		return errors.New("book rating cannot be nil")
	}

	rating.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, rating)
	if err != nil {
		return err
	}

	rating.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取书籍评分
func (r *MongoBookRatingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookRating, error) {
	var rating bookstore.BookRating
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&rating)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &rating, nil
}

// Update 更新书籍评分
func (r *MongoBookRatingRepository) Update(ctx context.Context, rating *bookstore.BookRating) error {
	if rating == nil {
		return errors.New("book rating cannot be nil")
	}

	rating.BeforeUpdate()

	filter := bson.M{"_id": rating.ID}
	update := bson.M{"$set": rating}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("book rating not found")
	}

	return nil
}

// Delete 删除书籍评分
func (r *MongoBookRatingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("book rating not found")
	}

	return nil
}

// GetAll 获取所有书籍评分
func (r *MongoBookRatingRepository) GetAll(ctx context.Context, limit, offset int) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var rating bookstore.BookRating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, cursor.Err()
}

// Count 统计书籍评分总数
func (r *MongoBookRatingRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// GetByBookID 根据书籍ID获取评分列表
func (r *MongoBookRatingRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"book_id": bookID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var rating bookstore.BookRating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, cursor.Err()
}

// GetByUserID 根据用户ID获取评分列表
func (r *MongoBookRatingRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var rating bookstore.BookRating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, cursor.Err()
}

// GetByBookIDAndUserID 根据书籍ID和用户ID获取评分
func (r *MongoBookRatingRepository) GetByBookIDAndUserID(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error) {
	var rating bookstore.BookRating
	filter := bson.M{"book_id": bookID, "user_id": userID}
	err := r.collection.FindOne(ctx, filter).Decode(&rating)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &rating, nil
}

// GetByRating 根据评分值获取评分列表
func (r *MongoBookRatingRepository) GetByRating(ctx context.Context, rating int, limit, offset int) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"rating": rating}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var ratingDoc bookstore.BookRating
		if err := cursor.Decode(&ratingDoc); err != nil {
			return nil, err
		}
		ratings = append(ratings, &ratingDoc)
	}

	return ratings, cursor.Err()
}

// GetByTags 根据标签获取评分列表
func (r *MongoBookRatingRepository) GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"tags": bson.M{"$in": tags}}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var rating bookstore.BookRating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, cursor.Err()
}

// Search 搜索评分
func (r *MongoBookRatingRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{
		"$or": []bson.M{
			{"comment": bson.M{"$regex": keyword, "$options": "i"}},
			{"tags": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var rating bookstore.BookRating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, cursor.Err()
}

// SearchByFilter 根据过滤器搜索评分
func (r *MongoBookRatingRepository) SearchByFilter(ctx context.Context, filter *BookstoreInterface.BookRatingFilter) ([]*bookstore.BookRating, error) {
	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	// 构建排序
	sortField := "created_at"
	sortOrder := -1
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// 构建查询条件
	query := bson.M{}
	
	if filter.BookID != nil {
		query["book_id"] = *filter.BookID
	}
	if filter.UserID != nil {
		query["user_id"] = *filter.UserID
	}
	if filter.MinRating != nil {
		if query["rating"] == nil {
			query["rating"] = bson.M{}
		}
		query["rating"].(bson.M)["$gte"] = *filter.MinRating
	}
	if filter.MaxRating != nil {
		if query["rating"] == nil {
			query["rating"] = bson.M{}
		}
		query["rating"].(bson.M)["$lte"] = *filter.MaxRating
	}
	if filter.HasComment != nil {
		if *filter.HasComment {
			query["comment"] = bson.M{"$ne": ""}
		} else {
			query["comment"] = ""
		}
	}
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	if filter.MinLikes != nil {
		query["likes"] = bson.M{"$gte": *filter.MinLikes}
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []*bookstore.BookRating
	for cursor.Next(ctx) {
		var rating bookstore.BookRating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, cursor.Err()
}

// CountByBookID 根据书籍ID统计评分数量
func (r *MongoBookRatingRepository) CountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	filter := bson.M{"book_id": bookID}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByUserID 根据用户ID统计评分数量
func (r *MongoBookRatingRepository) CountByUserID(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	filter := bson.M{"user_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByRating 根据评分值统计数量
func (r *MongoBookRatingRepository) CountByRating(ctx context.Context, bookID primitive.ObjectID, rating int) (int64, error) {
	filter := bson.M{
		"book_id": bookID,
		"rating":  rating,
	}
	return r.collection.CountDocuments(ctx, filter)
}

// GetAverageRating 获取书籍平均评分
func (r *MongoBookRatingRepository) GetAverageRating(ctx context.Context, bookID primitive.ObjectID) (float64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"book_id": bookID}},
		{"$group": bson.M{
			"_id":     nil,
			"average": bson.M{"$avg": "$rating"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Average float64 `bson:"average"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Average, nil
	}

	return 0, nil
}

// GetRatingDistribution 获取评分分布
func (r *MongoBookRatingRepository) GetRatingDistribution(ctx context.Context, bookID primitive.ObjectID) (map[int]int64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"book_id": bookID}},
		{"$group": bson.M{
			"_id":   bson.M{"$floor": "$rating"},
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	distribution := make(map[int]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    int   `bson:"_id"`
			Count int64 `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		distribution[result.ID] = result.Count
	}

	return distribution, cursor.Err()
}

// GetTopRatedBooks 获取评分最高的书籍
func (r *MongoBookRatingRepository) GetTopRatedBooks(ctx context.Context, limit int) ([]primitive.ObjectID, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":     "$book_id",
			"average": bson.M{"$avg": "$rating"},
			"count":   bson.M{"$sum": 1},
		}},
		{"$match": bson.M{"count": bson.M{"$gte": 5}}}, // 至少5个评分
		{"$sort": bson.M{"average": -1}},
		{"$limit": limit},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookIDs []primitive.ObjectID
	for cursor.Next(ctx) {
		var result struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		bookIDs = append(bookIDs, result.ID)
	}

	return bookIDs, cursor.Err()
}

// IncrementLikes 增加点赞数
func (r *MongoBookRatingRepository) IncrementLikes(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"likes": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("book rating not found")
	}

	return nil
}

// DecrementLikes 减少点赞数
func (r *MongoBookRatingRepository) DecrementLikes(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"likes": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("book rating not found")
	}

	return nil
}

// UpdateRating 更新评分
func (r *MongoBookRatingRepository) UpdateRating(ctx context.Context, id primitive.ObjectID, rating float64, comment string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"rating":     rating,
			"comment":    comment,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("book rating not found")
	}

	return nil
}

// BatchUpdateTags 批量更新标签
func (r *MongoBookRatingRepository) BatchUpdateTags(ctx context.Context, ratingIDs []primitive.ObjectID, tags []string) error {
	filter := bson.M{"_id": bson.M{"$in": ratingIDs}}
	update := bson.M{
		"$set": bson.M{
			"tags":       tags,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchDelete 批量删除评分
func (r *MongoBookRatingRepository) BatchDelete(ctx context.Context, ratingIDs []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ratingIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// BatchDeleteByBookID 根据书籍ID批量删除评分
func (r *MongoBookRatingRepository) BatchDeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	filter := bson.M{"book_id": bookID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// BatchDeleteByUserID 根据用户ID批量删除评分
func (r *MongoBookRatingRepository) BatchDeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// DeleteByBookID 根据书籍ID删除评分
func (r *MongoBookRatingRepository) DeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	filter := bson.M{"book_id": bookID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// DeleteByUserID 根据用户ID删除评分
func (r *MongoBookRatingRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// Transaction 执行事务
func (r *MongoBookRatingRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
func (r *MongoBookRatingRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}