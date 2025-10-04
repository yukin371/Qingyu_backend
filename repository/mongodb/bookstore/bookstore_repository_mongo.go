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
	infra "Qingyu_backend/repository/interfaces/infrastructure"
)

// MongoBookRepository MongoDB书籍仓储实现
type MongoBookRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoBookRepository 创建MongoDB书籍仓储实例
func NewMongoBookRepository(client *mongo.Client, database string) BookstoreInterface.BookRepository {
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

// Count 满足基础接口签名
func (r *MongoBookRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.collection.CountDocuments(ctx, query)
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

// GetByAuthorID 根据作者ID获取书籍
func (r *MongoBookRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{"author_id": authorID}, opts)
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
func (r *MongoBookRepository) GetByStatus(ctx context.Context, status bookstore.BookStatus, limit, offset int) ([]*bookstore.Book, error) {
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

// GetHotBooks 获取热门书籍（按浏览量/评分综合排序，这里按view_count降序）
func (r *MongoBookRepository) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"view_count", -1}, {"rating", -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"status": "published"}, opts)
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

// GetNewReleases 获取新上架书籍
func (r *MongoBookRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"published_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"status": "published"}, opts)
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

// GetFreeBooks 获取免费书籍
func (r *MongoBookRepository) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"updated_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"is_free": true}, opts)
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

// GetByPriceRange 按价格区间获取
func (r *MongoBookRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore.Book, error) {
	priceQuery := bson.M{}
	priceQuery["$gte"] = minPrice
	priceQuery["$lte"] = maxPrice
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"price", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"price": priceQuery}, opts)
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

// Search 搜索书籍（简化版本，符合接口定义）
func (r *MongoBookRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Book, error) {
	query := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"author": bson.M{"$regex": keyword, "$options": "i"}},
			{"introduction": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{"created_at", -1}})

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

// SearchWithPagination 搜索书籍（带分页和过滤，内部使用）
func (r *MongoBookRepository) SearchWithPagination(ctx context.Context, keyword string, filter *bookstore.BookFilter, page, pageSize int) ([]*bookstore.Book, int64, error) {
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
		if filter.Author != nil {
			query["author"] = bson.M{"$regex": *filter.Author, "$options": "i"}
		}
		if filter.IsRecommended != nil {
			query["is_recommended"] = *filter.IsRecommended
		}
		if filter.IsFeatured != nil {
			query["is_featured"] = *filter.IsFeatured
		}
		if len(filter.Tags) > 0 {
			query["tags"] = bson.M{"$in": filter.Tags}
		}
	}

	opts := options.Find()
	limit := pageSize
	offset := 0
	if page > 0 && pageSize > 0 {
		offset = (page - 1) * pageSize
	}
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}

	// 排序
	if filter != nil {
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
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
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
	// 忽略最小评分过滤以兼容当前模型
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
func (r *MongoBookRepository) CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// CountByFilter 根据过滤器统计
func (r *MongoBookRepository) CountByFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	query := bson.M{}
	if filter != nil {
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
	}
	return r.collection.CountDocuments(ctx, query)
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
				"rating":     rating,
				"updated_at": time.Now(),
			},
			"$inc": bson.M{"rating_count": 1},
		},
	)
	return err
}

// BatchUpdateStatus 批量更新状态
func (r *MongoBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": bookIDs}},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// BatchUpdateFeatured 批量更新精选状态
func (r *MongoBookRepository) BatchUpdateFeatured(ctx context.Context, bookIDs []primitive.ObjectID, isFeatured bool) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": bookIDs}},
		bson.M{
			"$set": bson.M{
				"is_featured": isFeatured,
				"updated_at":  time.Now(),
			},
		},
	)
	return err
}

// BatchUpdateRecommended 批量更新推荐状态
func (r *MongoBookRepository) BatchUpdateRecommended(ctx context.Context, bookIDs []primitive.ObjectID, isRecommended bool) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": bookIDs}},
		bson.M{
			"$set": bson.M{
				"is_recommended": isRecommended,
				"updated_at":     time.Now(),
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
				"updated_at":   time.Now(),
			},
		},
	)
	return err
}

// Transaction 执行事务
func (r *MongoBookRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.client.StartSession() // 启动会话
	if err != nil {
		return err
	}
	defer session.EndSession(ctx) // 确保会话结束

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	})
	return err
}

// List 实现基础接口的List
func (r *MongoBookRepository) List(ctx context.Context, filter infra.Filter) ([]*bookstore.Book, error) {
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
	var results []*bookstore.Book
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Exists 判断书籍是否存在
func (r *MongoBookRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
