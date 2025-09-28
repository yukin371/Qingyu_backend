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

// MongoBookDetailRepository MongoDB书籍详情仓储实现
type MongoBookDetailRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBookDetailRepository 创建MongoDB书籍详情仓储实例
func NewMongoBookDetailRepository(client *mongo.Client, database string) BookstoreInterface.BookDetailRepository {
	return &MongoBookDetailRepository{
		collection: client.Database(database).Collection("book_details"),
		client:     client,
	}
}

// Create 创建书籍详情
func (r *MongoBookDetailRepository) Create(ctx context.Context, bookDetail *bookstore.BookDetail) error {
	if bookDetail == nil {
		return errors.New("book detail cannot be nil")
	}

	bookDetail.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, bookDetail)
	if err != nil {
		return err
	}

	bookDetail.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取书籍详情
func (r *MongoBookDetailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookDetail, error) {
	var bookDetail bookstore.BookDetail
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bookDetail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bookDetail, nil
}

// Update 更新书籍详情
func (r *MongoBookDetailRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	if updates == nil || len(updates) == 0 {
		return errors.New("updates cannot be nil or empty")
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("book detail not found")
	}

	return nil
}

// Delete 删除书籍详情
func (r *MongoBookDetailRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("book detail not found")
	}

	return nil
}

// GetAll 获取所有书籍详情
func (r *MongoBookDetailRepository) GetAll(ctx context.Context, limit, offset int) ([]*bookstore.BookDetail, error) {
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

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// Count 统计书籍详情总数
func (r *MongoBookDetailRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// GetByTitle 根据标题获取书籍详情
func (r *MongoBookDetailRepository) GetByTitle(ctx context.Context, title string) (*bookstore.BookDetail, error) {
	var bookDetail bookstore.BookDetail
	filter := bson.M{"title": bson.M{"$regex": title, "$options": "i"}}
	err := r.collection.FindOne(ctx, filter).Decode(&bookDetail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bookDetail, nil
}

// GetByAuthor 根据作者获取书籍详情列表
func (r *MongoBookDetailRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.BookDetail, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"author": bson.M{"$regex": author, "$options": "i"}}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByAuthorID 根据作者ID获取书籍详情列表
func (r *MongoBookDetailRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore.BookDetail, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"author_id": authorID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByCategory 根据分类获取书籍详情列表
func (r *MongoBookDetailRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*bookstore.BookDetail, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"categories": category}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByStatus 根据状态获取书籍详情列表
func (r *MongoBookDetailRepository) GetByStatus(ctx context.Context, status bookstore.BookStatus, limit, offset int) ([]*bookstore.BookDetail, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByISBN 根据ISBN获取书籍详情
func (r *MongoBookDetailRepository) GetByISBN(ctx context.Context, isbn string) (*bookstore.BookDetail, error) {
	var bookDetail bookstore.BookDetail
	filter := bson.M{"isbn": isbn}
	err := r.collection.FindOne(ctx, filter).Decode(&bookDetail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bookDetail, nil
}

// GetByTags 根据标签获取书籍详情列表
func (r *MongoBookDetailRepository) GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore.BookDetail, error) {
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

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// Search 搜索书籍详情
func (r *MongoBookDetailRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.BookDetail, error) {
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
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"author": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
			{"tags": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// SearchByFilter 根据过滤器搜索书籍详情
func (r *MongoBookDetailRepository) SearchByFilter(ctx context.Context, filter *BookstoreInterface.BookDetailFilter) ([]*bookstore.BookDetail, error) {
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
	
	if filter.Title != "" {
		query["title"] = bson.M{"$regex": filter.Title, "$options": "i"}
	}
	if filter.Author != "" {
		query["author"] = bson.M{"$regex": filter.Author, "$options": "i"}
	}
	if filter.AuthorID != nil {
		query["author_id"] = *filter.AuthorID
	}
	if len(filter.Categories) > 0 {
		query["categories"] = bson.M{"$in": filter.Categories}
	}
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	if filter.Status != nil {
		query["status"] = *filter.Status
	}
	if filter.IsFree != nil {
		query["is_free"] = *filter.IsFree
	}
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		priceQuery := bson.M{}
		if filter.MinPrice != nil {
			priceQuery["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceQuery["$lte"] = *filter.MaxPrice
		}
		query["price"] = priceQuery
	}
	if filter.MinWordCount != nil || filter.MaxWordCount != nil {
		wordCountQuery := bson.M{}
		if filter.MinWordCount != nil {
			wordCountQuery["$gte"] = *filter.MinWordCount
		}
		if filter.MaxWordCount != nil {
			wordCountQuery["$lte"] = *filter.MaxWordCount
		}
		query["word_count"] = wordCountQuery
	}
	if filter.Publisher != "" {
		query["publisher"] = bson.M{"$regex": filter.Publisher, "$options": "i"}
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// CountByCategory 根据分类统计书籍数量
func (r *MongoBookDetailRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	filter := bson.M{"categories": category}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByAuthor 根据作者统计书籍数量
func (r *MongoBookDetailRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	filter := bson.M{"author": bson.M{"$regex": author, "$options": "i"}}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByStatus 根据状态统计书籍数量
func (r *MongoBookDetailRepository) CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	filter := bson.M{"status": status}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByTags 根据标签统计书籍数量
func (r *MongoBookDetailRepository) CountByTags(ctx context.Context, tags []string) (int64, error) {
	filter := bson.M{"tags": bson.M{"$in": tags}}
	return r.collection.CountDocuments(ctx, filter)
}

// BatchUpdateStatus 批量更新书籍状态
func (r *MongoBookDetailRepository) BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error {
	filter := bson.M{"_id": bson.M{"$in": bookIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchUpdateCategories 批量更新书籍分类
func (r *MongoBookDetailRepository) BatchUpdateCategories(ctx context.Context, bookIDs []primitive.ObjectID, categories []string) error {
	filter := bson.M{"_id": bson.M{"$in": bookIDs}}
	update := bson.M{
		"$set": bson.M{
			"categories": categories,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchUpdateTags 批量更新书籍标签
func (r *MongoBookDetailRepository) BatchUpdateTags(ctx context.Context, bookIDs []primitive.ObjectID, tags []string) error {
	filter := bson.M{"_id": bson.M{"$in": bookIDs}}
	update := bson.M{
		"$set": bson.M{
			"tags":       tags,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// Transaction 执行事务
func (r *MongoBookDetailRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
func (r *MongoBookDetailRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}