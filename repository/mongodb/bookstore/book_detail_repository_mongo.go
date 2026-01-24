package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
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

// MongoBookDetailRepository MongoDB书籍详情仓储实现
type MongoBookDetailRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBookDetailRepository 创建MongoDB书籍详情仓储实例
func NewMongoBookDetailRepository(client *mongo.Client, database string) BookstoreInterface.BookDetailRepository {
	return &MongoBookDetailRepository{
		collection: client.Database(database).Collection("books"),
		client:     client,
	}
}

// Create 创建书籍详情
func (r *MongoBookDetailRepository) Create(ctx context.Context, bookDetail *bookstore2.BookDetail) error {
	if bookDetail == nil {
		return errors.New("book detail cannot be nil")
	}

	bookDetail.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, bookDetail)
	if err != nil {
		return err
	}

	// ID 是 primitive.ObjectID 类型
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		bookDetail.ID = oid
	}
	return nil
}

// GetByID 根据ID获取书籍详情
func (r *MongoBookDetailRepository) GetByID(ctx context.Context, id string) (*bookstore2.BookDetail, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid book detail ID: %w", err)
	}
	var bookDetail bookstore2.BookDetail
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&bookDetail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bookDetail, nil
}

// Update 更新书籍详情
func (r *MongoBookDetailRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	if len(updates) == 0 {
		return errors.New("updates cannot be nil or empty")
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID}
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
func (r *MongoBookDetailRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("book detail not found")
	}

	return nil
}

// GetAll 获取所有书籍详情
func (r *MongoBookDetailRepository) GetAll(ctx context.Context, limit, offset int) ([]*bookstore2.BookDetail, error) {
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

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// Count 统计书籍详情总数
func (r *MongoBookDetailRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.collection.CountDocuments(ctx, query)
}

// GetByTitle 根据标题获取书籍详情
func (r *MongoBookDetailRepository) GetByTitle(ctx context.Context, title string) (*bookstore2.BookDetail, error) {
	var bookDetail bookstore2.BookDetail
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
func (r *MongoBookDetailRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore2.BookDetail, error) {
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

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByAuthorID 根据作者ID获取书籍详情列表
func (r *MongoBookDetailRepository) GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstore2.BookDetail, error) {
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, fmt.Errorf("invalid author ID: %w", err)
	}
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"author_id": objectID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByCategory 根据分类获取书籍详情列表
func (r *MongoBookDetailRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*bookstore2.BookDetail, error) {
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

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByStatus 根据状态获取书籍详情列表
func (r *MongoBookDetailRepository) GetByStatus(ctx context.Context, status bookstore2.BookStatus, limit, offset int) ([]*bookstore2.BookDetail, error) {
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

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByISBN 根据ISBN获取书籍详情
func (r *MongoBookDetailRepository) GetByISBN(ctx context.Context, isbn string) (*bookstore2.BookDetail, error) {
	var bookDetail bookstore2.BookDetail
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

// GetByPublisher 根据出版社获取书籍详情列表
func (r *MongoBookDetailRepository) GetByPublisher(ctx context.Context, publisher string, limit, offset int) ([]*bookstore2.BookDetail, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"publisher": publisher}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*bookstore2.BookDetail
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetByTags 根据标签获取书籍详情列表
func (r *MongoBookDetailRepository) GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore2.BookDetail, error) {
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

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// Search 搜索书籍详情
func (r *MongoBookDetailRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.BookDetail, error) {
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

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// SearchByFilter 根据过滤器搜索书籍详情
func (r *MongoBookDetailRepository) SearchByFilter(ctx context.Context, filter *BookstoreInterface.BookDetailFilter) ([]*bookstore2.BookDetail, error) {
	opts := options.Find()

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
	if len(filter.CategoryIDs) > 0 {
		query["category_ids"] = bson.M{"$in": filter.CategoryIDs}
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
	if filter.MinRating != nil || filter.MaxRating != nil {
		ratingQuery := bson.M{}
		if filter.MinRating != nil {
			ratingQuery["$gte"] = *filter.MinRating
		}
		if filter.MaxRating != nil {
			ratingQuery["$lte"] = *filter.MaxRating
		}
		query["rating"] = ratingQuery
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
	if filter.SerializedFrom != nil || filter.SerializedTo != nil {
		serializedQuery := bson.M{}
		if filter.SerializedFrom != nil {
			serializedQuery["$gte"] = *filter.SerializedFrom
		}
		if filter.SerializedTo != nil {
			serializedQuery["$lte"] = *filter.SerializedTo
		}
		query["serialized_at"] = serializedQuery
	}
	if filter.CompletedFrom != nil || filter.CompletedTo != nil {
		completedQuery := bson.M{}
		if filter.CompletedFrom != nil {
			completedQuery["$gte"] = *filter.CompletedFrom
		}
		if filter.CompletedTo != nil {
			completedQuery["$lte"] = *filter.CompletedTo
		}
		query["completed_at"] = completedQuery
	}
	if filter.CreatedAtFrom != nil || filter.CreatedAtTo != nil {
		createdQuery := bson.M{}
		if filter.CreatedAtFrom != nil {
			createdQuery["$gte"] = *filter.CreatedAtFrom
		}
		if filter.CreatedAtTo != nil {
			createdQuery["$lte"] = *filter.CreatedAtTo
		}
		query["created_at"] = createdQuery
	}
	if filter.UpdatedAtFrom != nil || filter.UpdatedAtTo != nil {
		updatedQuery := bson.M{}
		if filter.UpdatedAtFrom != nil {
			updatedQuery["$gte"] = *filter.UpdatedAtFrom
		}
		if filter.UpdatedAtTo != nil {
			updatedQuery["$lte"] = *filter.UpdatedAtTo
		}
		query["updated_at"] = updatedQuery
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookDetails []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var bookDetail bookstore2.BookDetail
		if err := cursor.Decode(&bookDetail); err != nil {
			return nil, err
		}
		bookDetails = append(bookDetails, &bookDetail)
	}

	return bookDetails, cursor.Err()
}

// GetByBookID 根据书籍基础ID获取详情
func (r *MongoBookDetailRepository) GetByBookID(ctx context.Context, bookID string) (*bookstore2.BookDetail, error) {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}
	var detail bookstore2.BookDetail
	err = r.collection.FindOne(ctx, bson.M{"book_id": objectID}).Decode(&detail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

// GetByBookIDs 批量根据书籍基础ID获取详情
func (r *MongoBookDetailRepository) GetByBookIDs(ctx context.Context, bookIDs []string) ([]*bookstore2.BookDetail, error) {
	// 转换 string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid book ID: %w", err)
		}
		objectIDs = append(objectIDs, objectID)
	}
	cursor, err := r.collection.Find(ctx, bson.M{"book_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*bookstore2.BookDetail
	for cursor.Next(ctx) {
		var d bookstore2.BookDetail
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &d)
	}
	return results, cursor.Err()
}

// UpdateAuthor 更新作者信息
func (r *MongoBookDetailRepository) UpdateAuthor(ctx context.Context, bookID string, authorID string, authorName string) error {
	bookObjectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	authorObjectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return fmt.Errorf("invalid author ID: %w", err)
	}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bookObjectID},
		bson.M{"$set": bson.M{"author_id": authorObjectID, "author": authorName, "updated_at": time.Now()}},
	)
	return err
}

// GetSimilarBooks 获取相似书籍（基于标签和分类）
func (r *MongoBookDetailRepository) GetSimilarBooks(ctx context.Context, bookID string, limit int) ([]*bookstore2.BookDetail, error) {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID: %w", err)
	}
	// 先获取目标书的标签
	var current bookstore2.BookDetail
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&current); err != nil {
		if err == mongo.ErrNoDocuments {
			return []*bookstore2.BookDetail{}, nil
		}
		return nil, err
	}
	query := bson.M{
		"_id": bson.M{"$ne": bookID},
		"$or": []bson.M{
			{"tags": bson.M{"$in": current.Tags}},
			{"category_ids": bson.M{"$in": current.CategoryIDs}},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*bookstore2.BookDetail
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
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
func (r *MongoBookDetailRepository) CountByStatus(ctx context.Context, status bookstore2.BookStatus) (int64, error) {
	filter := bson.M{"status": status}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByTags 根据标签统计书籍数量
func (r *MongoBookDetailRepository) CountByTags(ctx context.Context, tags []string) (int64, error) {
	filter := bson.M{"tags": bson.M{"$in": tags}}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByPublisher 根据出版社统计数量
func (r *MongoBookDetailRepository) CountByPublisher(ctx context.Context, publisher string) (int64, error) {
	filter := bson.M{"publisher": publisher}
	return r.collection.CountDocuments(ctx, filter)
}

// BatchUpdateStatus 批量更新书籍状态
func (r *MongoBookDetailRepository) BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstore2.BookStatus) error {
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
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// DecrementCommentCount递减评论数
func (r *MongoBookDetailRepository) DecrementCommentCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"comment_count": -1}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// BatchUpdateCategories 批量更新书籍分类
func (r *MongoBookDetailRepository) BatchUpdateCategories(ctx context.Context, bookIDs []string, categoryIDs []string) error {
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
			"category_ids": categoryIDs,
			"updated_at":   time.Now(),
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchUpdatePublisher 批量更新出版社
func (r *MongoBookDetailRepository) BatchUpdatePublisher(ctx context.Context, bookIDs []string, publisher string) error {
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
			"publisher":  publisher,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchUpdateTags 批量更新书籍标签
func (r *MongoBookDetailRepository) BatchUpdateTags(ctx context.Context, bookIDs []string, tags []string) error {
	// 转换 string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid book ID: %w", err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{
		"$set": bson.M{
			"tags":       tags,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// IncrementViewCount 递增浏览量
func (r *MongoBookDetailRepository) IncrementViewCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"view_count": 1}, "$set": bson.M{"updated_at": time.Now()}})
	return err
}

// IncrementLikeCount 递增点赞数
func (r *MongoBookDetailRepository) IncrementLikeCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"like_count": 1}, "$set": bson.M{"updated_at": time.Now()}})
	return err
}

// DecrementLikeCount 递减点赞数
func (r *MongoBookDetailRepository) DecrementLikeCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"like_count": -1}, "$set": bson.M{"updated_at": time.Now()}})
	return err
}

// IncrementCommentCount 递增评论数
func (r *MongoBookDetailRepository) IncrementCommentCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"comment_count": 1}, "$set": bson.M{"updated_at": time.Now()}})
	return err
}

// IncrementShareCount 递增分享数
func (r *MongoBookDetailRepository) IncrementShareCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"share_count": 1}, "$set": bson.M{"updated_at": time.Now()}})
	return err
}

// UpdateRating 更新评分统计
func (r *MongoBookDetailRepository) UpdateRating(ctx context.Context, bookID string, rating float64, ratingCount int64) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"rating": rating, "rating_count": ratingCount, "updated_at": time.Now()}},
	)
	return err
}

// UpdateLastChapter 更新最后章节标题
func (r *MongoBookDetailRepository) UpdateLastChapter(ctx context.Context, bookID string, chapterTitle string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("invalid book ID: %w", err)
	}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"last_chapter": chapterTitle, "updated_at": time.Now()}},
	)
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

// List 根据过滤条件列出书籍详情
func (r *MongoBookDetailRepository) List(ctx context.Context, filter infra.Filter) ([]*bookstore2.BookDetail, error) {
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
	var results []*bookstore2.BookDetail
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Exists 判断记录是否存在
func (r *MongoBookDetailRepository) Exists(ctx context.Context, id string) (bool, error) {
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

// CountByAuthorID 根据作者ID统计书籍数量
func (r *MongoBookDetailRepository) CountByAuthorID(ctx context.Context, authorID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return 0, fmt.Errorf("invalid author ID: %w", err)
	}
	filter := bson.M{"author_id": objectID}
	return r.collection.CountDocuments(ctx, filter)
}

// CountByFilter 根据过滤条件统计书籍数量
// 注意：复用现有 SearchByFilter 的过滤逻辑
func (r *MongoBookDetailRepository) CountByFilter(ctx context.Context, filterObj *BookstoreInterface.BookDetailFilter) (int64, error) {
	// 构建查询条件（与 SearchByFilter 中的逻辑保持一致）
	query := bson.M{}

	if filterObj.Title != "" {
		query["title"] = bson.M{"$regex": filterObj.Title, "$options": "i"}
	}
	if filterObj.Author != "" {
		query["author"] = bson.M{"$regex": filterObj.Author, "$options": "i"}
	}
	if filterObj.AuthorID != nil {
		query["author_id"] = *filterObj.AuthorID
	}
	if len(filterObj.CategoryIDs) > 0 {
		query["category_ids"] = bson.M{"$in": filterObj.CategoryIDs}
	}
	if len(filterObj.Tags) > 0 {
		query["tags"] = bson.M{"$in": filterObj.Tags}
	}
	if filterObj.Status != nil {
		query["status"] = *filterObj.Status
	}

	return r.collection.CountDocuments(ctx, query)
}
