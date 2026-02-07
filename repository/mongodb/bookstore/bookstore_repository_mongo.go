package mongodb

import (
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/mongodb/base"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	infra "Qingyu_backend/repository/interfaces/infrastructure"
)

// MongoBookRepository MongoDB书籍仓储实现
type MongoBookRepository struct {
	*base.BaseMongoRepository
	client *mongo.Client
}

// NewMongoBookRepository 创建MongoDB书籍仓储实例
func NewMongoBookRepository(client *mongo.Client, database string) BookstoreInterface.BookRepository {
	db := client.Database(database)
	return &MongoBookRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "books"),
		client:              client,
	}
}

// Create 创建书籍
func (r *MongoBookRepository) Create(ctx context.Context, book *bookstore.Book) error {
	if book == nil {
		return errors.New("book cannot be nil")
	}

	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	result, err := r.GetCollection().InsertOne(ctx, book)
	if err != nil {
		return err
	}

	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取书籍
func (r *MongoBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
	log.Printf("[DEBUG] GetByID(%s) database: %s, collection: %s\n",
		id, r.GetDB().Name(), r.GetCollection().Name())

	// 由于数据库中的 _id 字段是字符串类型，直接使用字符串查询
	// 而不是转换为 ObjectID
	var book bookstore.Book
	err := r.GetCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[DEBUG] GetByID(%s) no documents found for query: {_id: %s}\n", id, id)
			return nil, nil
		}
		log.Printf("[DEBUG] GetByID(%s) decode error: %v\n", id, err)
		return nil, err
	}

	log.Printf("[DEBUG] GetByID(%s) found book: %s, status: %s\n", id, book.Title, book.Status)
	return &book, nil
}

// getSampleBookIDs 辅助函数：从 BSON 文档中提取 ID
func getSampleBookIDs(docs []bson.M) []string {
	ids := make([]string, 0, len(docs))
	for _, doc := range docs {
		if id, ok := doc["_id"].(primitive.ObjectID); ok {
			ids = append(ids, id.Hex())
		}
	}
	return ids
}

// Update 更新书籍
func (r *MongoBookRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()

	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
func (r *MongoBookRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return err
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
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
	return r.GetCollection().CountDocuments(ctx, query)
}

// GetByTitle 根据标题获取书籍
func (r *MongoBookRepository) GetByTitle(ctx context.Context, title string) (*bookstore.Book, error) {
	var book bookstore.Book
	err := r.GetCollection().FindOne(ctx, bson.M{"title": title}).Decode(&book)
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
	cursor, err := r.GetCollection().Find(ctx, bson.M{"author": author}, opts)
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
func (r *MongoBookRepository) GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstore.Book, error) {
	objectID, err := r.ParseID(authorID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.GetCollection().Find(ctx, bson.M{"author_id": objectID}, opts)
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
func (r *MongoBookRepository) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstore.Book, error) {
	objectID, err := r.ParseID(categoryID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.GetCollection().Find(ctx, bson.M{"category_ids": objectID}, opts)
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
	cursor, err := r.GetCollection().Find(ctx, bson.M{"status": status}, opts)
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

// GetRecommended 获取所有书籍（书库列表）
// 修改：返回所有书籍，按创建时间倒序排序
func (r *MongoBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}, {Key: "title", Value: 1}})

	// 查询所有书籍，不再只返回推荐的
	cursor, err := r.GetCollection().Find(ctx, bson.M{}, opts)
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
		SetSort(bson.D{{Key: "rating", Value: -1}, {Key: "view_count", Value: -1}})

	// 查询精选书籍，包含所有已发布状态（published, ongoing, completed）
	cursor, err := r.GetCollection().Find(ctx, bson.M{
		"is_featured": true,
		"status":      bson.M{"$in": []string{"published", "ongoing", "completed"}},
	}, opts)
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
		SetSort(bson.D{{Key: "view_count", Value: -1}, {Key: "rating", Value: -1}})

	// 查询热门书籍，包含所有已发布状态（published, ongoing, completed）
	cursor, err := r.GetCollection().Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"published", "ongoing", "completed"}},
	}, opts)
	if err != nil {
		log.Printf("[DEBUG] GetHotBooks failed: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		log.Printf("[DEBUG] GetHotBooks decode failed: %v\n", err)
		return nil, err
	}

	log.Printf("[DEBUG] GetHotBooks found %d books, IDs: %v\n", len(books), getHotBookIDs(books))
	return books, nil
}

// getHotBookIDs 辅助函数：从 Book 指针切片中提取 ID
func getHotBookIDs(books []*bookstore.Book) []string {
	ids := make([]string, 0, len(books))
	for _, book := range books {
		ids = append(ids, book.ID.Hex())
	}
	return ids
}

// GetNewReleases 获取新上架书籍
func (r *MongoBookRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "published_at", Value: -1}})
	// 查询新上架书籍，包含所有已发布状态（published, ongoing, completed）
	cursor, err := r.GetCollection().Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"published", "ongoing", "completed"}},
	}, opts)
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
		SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.GetCollection().Find(ctx, bson.M{"is_free": true}, opts)
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
		SetSort(bson.D{{Key: "price", Value: 1}})
	cursor, err := r.GetCollection().Find(ctx, bson.M{"price": priceQuery}, opts)
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
	// 为了避免正则表达式 UTF-8 编码问题，我们获取所有书籍然后在 Go 代码中进行过滤
	// 这不是最优解，但可以确保中文搜索正常工作
	// TODO: 考虑使用 MongoDB Atlas Search 或文本索引
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 查询所有已发布的书籍
	cursor, err := r.GetCollection().Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"published", "ongoing", "completed"}},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allBooks []*bookstore.Book
	if err = cursor.All(ctx, &allBooks); err != nil {
		return nil, err
	}

	// 在 Go 代码中进行过滤
	keywordLower := strings.ToLower(keyword)
	var filteredBooks []*bookstore.Book
	for _, book := range allBooks {
		if strings.Contains(strings.ToLower(book.Title), keywordLower) ||
			strings.Contains(strings.ToLower(book.Author), keywordLower) ||
			strings.Contains(strings.ToLower(book.Introduction), keywordLower) {
			filteredBooks = append(filteredBooks, book)
		}
	}

	return filteredBooks, nil
}

// SearchWithPagination 搜索书籍（带分页和过滤，内部使用）
func (r *MongoBookRepository) SearchWithPagination(ctx context.Context, keyword string, filter *bookstore.BookFilter, page, pageSize int) ([]*bookstore.Book, int64, error) {
	// 使用 $indexOfCP 进行字符串搜索，避免正则表达式的 UTF-8 编码问题
	query := bson.M{
		"$or": []bson.M{
			{
				"$expr": bson.M{
					"$gt": bson.A{
						bson.M{"$indexOfCP": bson.A{"$title", keyword}},
						-1,
					},
				},
			},
			{
				"$expr": bson.M{
					"$gt": bson.A{
						bson.M{"$indexOfCP": bson.A{"$author", keyword}},
						-1,
					},
				},
			},
			{
				"$expr": bson.M{
					"$gt": bson.A{
						bson.M{"$indexOfCP": bson.A{"$introduction", keyword}},
						-1,
					},
				},
			},
		},
	}

	// 应用过滤器
	if filter != nil {
		if filter.Status != nil {
			query["status"] = *filter.Status
		}
		if filter.CategoryID != nil {
			objectID, err := primitive.ObjectIDFromHex(*filter.CategoryID)
			if err == nil {
				query["category_ids"] = objectID
			}
		}
		if filter.Author != nil {
			// 使用 $indexOfCP 避免正则表达式编码问题
			query["$expr"] = bson.M{
				"$gt": bson.A{
					bson.M{"$indexOfCP": bson.A{"$author", *filter.Author}},
					-1,
				},
			}
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
		// 多字段排序：先按创建时间倒序，时间相同时按标题拼音排序
		sortFields := bson.D{
			{Key: "created_at", Value: -1}, // 创建时间倒序
			{Key: "title", Value: 1},       // 标题拼音升序
		}

		// 如果指定了排序字段，使用指定的排序
		if filter.SortBy != "" {
			sortOrder := -1
			if filter.SortOrder == "asc" {
				sortOrder = 1
			}
			// 对于创建时间排序，仍然保持标题作为次要排序
			if filter.SortBy == "created_at" {
				sortFields = bson.D{
					{Key: "created_at", Value: sortOrder},
					{Key: "title", Value: 1},
				}
			} else {
				// 其他字段只按指定字段排序
				sortFields = bson.D{{Key: filter.SortBy, Value: sortOrder}}
			}
		}
		opts.SetSort(sortFields)
	}

	cursor, err := r.GetCollection().Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var books []*bookstore.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, 0, err
	}

	total, err := r.GetCollection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// SearchWithFilter 使用过滤器搜索书籍
func (r *MongoBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	// 先构建基础查询（不包括关键词）
	query := bson.M{}

	if filter.Status != nil {
		query["status"] = *filter.Status
	}
	if filter.CategoryID != nil {
		objectID, err := primitive.ObjectIDFromHex(*filter.CategoryID)
		if err == nil {
			query["category_ids"] = objectID
		}
	}
	if filter.Author != nil {
		query["author"] = *filter.Author
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

	// 如果有关键词，使用正则表达式进行搜索
	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := *filter.Keyword
		log.Printf("[DEBUG] 搜索关键词: %s (字节: %v)", keyword, []byte(keyword))

		// 直接在Go代码中构建原始的BSON正则表达式
		// 使用bson.RawValue来避免编码问题
		orConditions := []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"author": bson.M{"$regex": keyword, "$options": "i"}},
			{"introduction": bson.M{"$regex": keyword, "$options": "i"}},
		}
		query["$or"] = orConditions
		log.Printf("[DEBUG] MongoDB正则查询条件已添加: %s (长度: %d)", keyword, len(keyword))
	}

	// 没有关键词，直接查询
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
	opts.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := r.GetCollection().Find(ctx, query, opts)
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

// containsStringIgnoreCase 检查字符串是否包含子字符串（不区分大小写）
func containsStringIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// CountByCategory 统计分类下的书籍数量
func (r *MongoBookRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return 0, err
	}
	return r.GetCollection().CountDocuments(ctx, bson.M{"category_ids": objectID})
}

// CountByAuthor 统计作者的书籍数量
func (r *MongoBookRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	return r.GetCollection().CountDocuments(ctx, bson.M{"author": author})
}

// CountByStatus 统计指定状态的书籍数量
func (r *MongoBookRepository) CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	return r.GetCollection().CountDocuments(ctx, bson.M{"status": status})
}

// CountByFilter 根据过滤器统计
func (r *MongoBookRepository) CountByFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	query := bson.M{}
	if filter != nil {
		if filter.Status != nil {
			query["status"] = *filter.Status
		}
		if filter.CategoryID != nil {
			objectID, err := primitive.ObjectIDFromHex(*filter.CategoryID)
			if err == nil {
				query["category_ids"] = objectID
			}
		}
		if filter.Author != nil {
			// 使用 $indexOfCP 避免正则表达式编码问题
			query["$expr"] = bson.M{
				"$gt": bson.A{
					bson.M{"$indexOfCP": bson.A{"$author", *filter.Author}},
					-1,
				},
			}
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
			keyword := *filter.Keyword
			// 使用 $or 条件配合 $indexOfCP 进行关键词搜索
			orConditions := []bson.M{
				{
					"$expr": bson.M{
						"$gt": bson.A{
							bson.M{"$indexOfCP": bson.A{"$title", keyword}},
							-1,
						},
					},
				},
				{
					"$expr": bson.M{
						"$gt": bson.A{
							bson.M{"$indexOfCP": bson.A{"$author", keyword}},
							-1,
						},
					},
				},
				{
					"$expr": bson.M{
						"$gt": bson.A{
							bson.M{"$indexOfCP": bson.A{"$introduction", keyword}},
							-1,
						},
					},
				},
			}
			query["$or"] = orConditions
		}
	}
	return r.GetCollection().CountDocuments(ctx, query)
}

// GetStats 获取书籍统计信息
func (r *MongoBookRepository) GetStats(ctx context.Context) (*bookstore.BookStats, error) {
	stats := &bookstore.BookStats{}

	// 总书籍数
	total, err := r.GetCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalBooks = total

	// 已发布书籍数
	published, err := r.GetCollection().CountDocuments(ctx, bson.M{"status": "published"})
	if err != nil {
		return nil, err
	}
	stats.PublishedBooks = published

	// 草稿书籍数
	draft, err := r.GetCollection().CountDocuments(ctx, bson.M{"status": "draft"})
	if err != nil {
		return nil, err
	}
	stats.DraftBooks = draft

	// 推荐书籍数
	recommended, err := r.GetCollection().CountDocuments(ctx, bson.M{"is_recommended": true})
	if err != nil {
		return nil, err
	}
	stats.RecommendedBooks = recommended

	// 精选书籍数
	featured, err := r.GetCollection().CountDocuments(ctx, bson.M{"is_featured": true})
	if err != nil {
		return nil, err
	}
	stats.FeaturedBooks = featured

	return stats, nil
}

// IncrementLikeCount 增加点赞数
func (r *MongoBookRepository) IncrementLikeCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return err
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"like_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// IncrementCommentCount 增加评论数
func (r *MongoBookRepository) IncrementCommentCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return err
	}

	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{"comment_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// UpdateRating 更新评分
func (r *MongoBookRepository) UpdateRating(ctx context.Context, bookID string, rating float64) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return err
	}

	// 这里简化处理，实际应该计算平均评分
	_, err = r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
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
func (r *MongoBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstore.BookStatus) error {
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	_, err := r.GetCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": objectIDs}},
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
func (r *MongoBookRepository) BatchUpdateFeatured(ctx context.Context, bookIDs []string, isFeatured bool) error {
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	_, err := r.GetCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": objectIDs}},
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
func (r *MongoBookRepository) BatchUpdateRecommended(ctx context.Context, bookIDs []string, isRecommended bool) error {
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	_, err := r.GetCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": objectIDs}},
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
func (r *MongoBookRepository) BatchUpdateCategory(ctx context.Context, bookIDs []string, categoryIDs []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(bookIDs))
	for _, id := range bookIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	categoryObjectIDs := make([]primitive.ObjectID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		categoryObjectIDs = append(categoryObjectIDs, objectID)
	}

	_, err := r.GetCollection().UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": objectIDs}},
		bson.M{
			"$set": bson.M{
				"category_ids": categoryObjectIDs,
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
	cursor, err := r.GetCollection().Find(ctx, query, opts)
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
func (r *MongoBookRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	count, err := r.GetCollection().CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IncrementViewCount 增加浏览计数
func (r *MongoBookRepository) IncrementViewCount(ctx context.Context, bookID string) error {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$inc": bson.M{"view_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err = r.GetCollection().UpdateOne(ctx, filter, update)
	return err
}

// GetYears 获取所有书籍的发布年份列表（去重，倒序）
func (r *MongoBookRepository) GetYears(ctx context.Context) ([]int, error) {
	// 使用聚合管道提取年份并去重
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"published_at": bson.M{"$ne": nil}, // 只查询有发布时间的书籍
			},
		},
		{
			"$project": bson.M{
				"year": bson.M{"$year": "$published_at"}, // 提取年份
			},
		},
		{
			"$group": bson.M{
				"_id": "$year", // 按年份分组去重
			},
		},
		{
			"$sort": bson.M{
				"_id": -1, // 按年份倒序
			},
		},
	}

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type YearResult struct {
		Year int `bson:"_id"`
	}

	var results []YearResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	years := make([]int, 0, len(results))
	for _, result := range results {
		years = append(years, result.Year)
	}

	return years, nil
}

// GetTags 获取所有标签列表（去重，排序）
// 如果提供了 categoryID，则只返回该分类下的书籍标签
func (r *MongoBookRepository) GetTags(ctx context.Context, categoryID *string) ([]string, error) {
	// 构建聚合管道
	pipeline := []bson.M{}

	// 如果提供了分类ID，先过滤书籍
	if categoryID != nil && *categoryID != "" {
		objectID, err := primitive.ObjectIDFromHex(*categoryID)
		if err != nil {
			return nil, err
		}
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"category_ids": objectID,
			},
		})
	}

	// 展开标签数组
	pipeline = append(pipeline, bson.M{
		"$unwind": "$tags",
	})

	// 按标签分组去重
	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id": "$tags",
		},
	})

	// 按标签名排序
	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"_id": 1, // 按标签名升序
		},
	})

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type TagResult struct {
		Tag string `bson:"_id"`
	}

	var results []TagResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	tags := make([]string, 0, len(results))
	for _, result := range results {
		if result.Tag != "" { // 过滤空标签
			tags = append(tags, result.Tag)
		}
	}

	return tags, nil
}

// escapeRegex 转义正则表达式特殊字符，避免中文编码问题
func escapeRegex(s string) string {
	// 特殊字符转义
	specialChars := []string{`\`, `.`, `^`, `$`, `*`, `+`, `?`, `(`, `)`, `[`, `]`, `{`, `}`, `|`}
	for _, c := range specialChars {
		s = strings.ReplaceAll(s, c, `\`+c)
	}
	return s
}
