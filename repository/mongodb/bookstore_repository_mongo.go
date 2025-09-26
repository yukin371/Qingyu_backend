package mongodb

import (
	"context"
	"time"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"Qingyu_backend/models/reading/bookstore"
	"Qingyu_backend/repository/interfaces"
)

// MongoBookRepository MongoDB书籍仓储实现
type MongoBookRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBookRepository 创建MongoDB书籍仓储实例
func NewMongoBookRepository(client *mongo.Client, database string) interfaces.BookRepository {
	return &MongoBookRepository{
		collection: client.Database(database).Collection("books"),
		client:     client,
	}
}

// Create 创建书籍
func (r *MongoBookRepository) Create(ctx context.Context, book *bookstore.Book) error {
	if book == nil {
		return errors.New("book cannot be nil")
	}
	
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		return err
	}
	
	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取书籍
func (r *MongoBookRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Book, error) {
	var book bookstore.Book
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

// Update 更新书籍
func (r *MongoBookRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return errors.New("book not found")
	}
	
	return nil
}

// Delete 删除书籍
func (r *MongoBookRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return errors.New("book not found")
	}
	
	return nil
}

// Health 健康检查
func (r *MongoBookRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// GetByTitle 根据标题获取书籍
func (r *MongoBookRepository) GetByTitle(ctx context.Context, title string) (*bookstore.Book, error) {
	var book bookstore.Book
	err := r.collection.FindOne(ctx, bson.M{"title": title}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

// GetByAuthor 根据作者获取书籍列表
func (r *MongoBookRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{"author": author}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// GetByCategory 根据分类获取书籍列表
func (r *MongoBookRepository) GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{"category_ids": categoryID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// GetByStatus 根据状态获取书籍列表
func (r *MongoBookRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// GetRecommended 获取推荐书籍
func (r *MongoBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"rating", -1}, {"view_count", -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{"is_recommended": true, "status": "published"}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// GetFeatured 获取精选书籍
func (r *MongoBookRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"rating", -1}, {"view_count", -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{"is_featured": true, "status": "published"}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// Search 搜索书籍
func (r *MongoBookRepository) Search(ctx context.Context, keyword string, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	query := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"author": bson.M{"$regex": keyword, "$options": "i"}},
			{"introduction": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	
	// 应用过滤器
	if filter != nil {
		if filter.Status != nil {
			query["status"] = *filter.Status
		}
		if filter.CategoryID != nil {
			query["category_ids"] = *filter.CategoryID
		}
		if filter.IsRecommended != nil {
			query["is_recommended"] = *filter.IsRecommended
		}
		if filter.IsFeatured != nil {
			query["is_featured"] = *filter.IsFeatured
		}
		if filter.MinRating != nil {
			query["rating"] = bson.M{"$gte": *filter.MinRating}
		}
		if len(filter.Tags) > 0 {
			query["tags"] = bson.M{"$in": filter.Tags}
		}
	}
	
	opts := options.Find()
	if filter != nil {
		if filter.Limit > 0 {
			opts.SetLimit(int64(filter.Limit))
		}
		if filter.Offset > 0 {
			opts.SetSkip(int64(filter.Offset))
		}
		
		// 排序
		sortBy := "created_at"
		sortOrder := -1
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder == "asc" {
			sortOrder = 1
		}
		opts.SetSort(bson.D{{sortBy, sortOrder}})
	}
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// SearchWithFilter 使用过滤器搜索书籍
func (r *MongoBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	query := bson.M{}
	
	if filter.Status != nil {
		query["status"] = *filter.Status
	}
	if filter.CategoryID != nil {
		query["category_ids"] = *filter.CategoryID
	}
	if filter.Author != nil {
		query["author"] = bson.M{"$regex": *filter.Author, "$options": "i"}
	}
	if filter.IsRecommended != nil {
		query["is_recommended"] = *filter.IsRecommended
	}
	if filter.IsFeatured != nil {
		query["is_featured"] = *filter.IsFeatured
	}
	if filter.MinRating != nil {
		query["rating"] = bson.M{"$gte": *filter.MinRating}
	}
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}
	if filter.Keyword != nil {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": *filter.Keyword, "$options": "i"}},
			{"author": bson.M{"$regex": *filter.Keyword, "$options": "i"}},
			{"introduction": bson.M{"$regex": *filter.Keyword, "$options": "i"}},
		}
	}
	
	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}
	
	// 排序
	sortBy := "created_at"
	sortOrder := -1
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	opts.SetSort(bson.D{{sortBy, sortOrder}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// CountByCategory 统计分类下的书籍数量
func (r *MongoBookRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"category_ids": categoryID})
}

// CountByAuthor 统计作者的书籍数量
func (r *MongoBookRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"author": author})
}

// CountByStatus 统计指定状态的书籍数量
func (r *MongoBookRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// GetStats 获取书籍统计信息
func (r *MongoBookRepository) GetStats(ctx context.Context) (*bookstore.BookStats, error) {
	stats := &bookstore.BookStats{}
	
	// 总书籍数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalBooks = total
	
	// 已发布书籍数
	published, err := r.collection.CountDocuments(ctx, bson.M{"status": "published"})
	if err != nil {
		return nil, err
	}
	stats.PublishedBooks = published
	
	// 草稿书籍数
	draft, err := r.collection.CountDocuments(ctx, bson.M{"status": "draft"})
	if err != nil {
		return nil, err
	}
	stats.DraftBooks = draft
	
	// 推荐书籍数
	recommended, err := r.collection.CountDocuments(ctx, bson.M{"is_recommended": true})
	if err != nil {
		return nil, err
	}
	stats.RecommendedBooks = recommended
	
	// 精选书籍数
	featured, err := r.collection.CountDocuments(ctx, bson.M{"is_featured": true})
	if err != nil {
		return nil, err
	}
	stats.FeaturedBooks = featured
	
	return stats, nil
}

// IncrementViewCount 增加浏览量
func (r *MongoBookRepository) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bookID},
		bson.M{
			"$inc": bson.M{"view_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// IncrementLikeCount 增加点赞数
func (r *MongoBookRepository) IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bookID},
		bson.M{
			"$inc": bson.M{"like_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// IncrementCommentCount 增加评论数
func (r *MongoBookRepository) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bookID},
		bson.M{
			"$inc": bson.M{"comment_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// UpdateRating 更新评分
func (r *MongoBookRepository) UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	// 这里简化处理，实际应该计算平均评分
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bookID},
		bson.M{
			"$set": bson.M{
				"rating": rating,
				"updated_at": time.Now(),
			},
			"$inc": bson.M{"rating_count": 1},
		},
	)
	return err
}

// BatchUpdateStatus 批量更新状态
func (r *MongoBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": bookIDs}},
		bson.M{
			"$set": bson.M{
				"status": status,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// BatchUpdateCategory 批量更新分类
func (r *MongoBookRepository) BatchUpdateCategory(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []primitive.ObjectID) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": bookIDs}},
		bson.M{
			"$set": bson.M{
				"category_ids": categoryIDs,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// Transaction 执行事务
func (r *MongoBookRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	
	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		return session.WithTransaction(sc, func(sc mongo.SessionContext) (interface{}, error) {
			return nil, fn(sc)
		})
	})
}