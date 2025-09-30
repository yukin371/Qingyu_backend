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

// MongoChapterRepository MongoDB章节仓储实现
type MongoChapterRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoChapterRepository 创建MongoDB章节仓储实例
func NewMongoChapterRepository(client *mongo.Client, database string) BookstoreInterface.ChapterRepository {
	return &MongoChapterRepository{
		collection: client.Database(database).Collection("chapters"),
		client:     client,
	}
}

// Create 创建章节
func (r *MongoChapterRepository) Create(ctx context.Context, chapter *bookstore.Chapter) error {
	if chapter == nil {
		return errors.New("chapter cannot be nil")
	}

	chapter.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, chapter)
	if err != nil {
		return err
	}

	chapter.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取章节
func (r *MongoChapterRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Chapter, error) {
	var chapter bookstore.Chapter
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

// Update 更新章节
func (r *MongoChapterRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("updates cannot be empty")
	}

	// 添加更新时间戳
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("chapter not found")
	}

	return nil
}

// Delete 删除章节
func (r *MongoChapterRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("chapter not found")
	}

	return nil
}

// GetAll 获取所有章节
func (r *MongoChapterRepository) GetAll(ctx context.Context, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// Count 统计章节总数
func (r *MongoChapterRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.collection.CountDocuments(ctx, query)
}

// List 根据过滤条件列出章节
func (r *MongoChapterRepository) List(ctx context.Context, filter infra.Filter) ([]*bookstore.Chapter, error) {
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
	var results []*bookstore.Chapter
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Exists 判断章节是否存在
func (r *MongoChapterRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByBookID 根据书籍ID获取章节列表
func (r *MongoChapterRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	filter := bson.M{"book_id": bookID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// GetByBookIDAndChapterNum 根据书籍ID和章节号获取章节
func (r *MongoChapterRepository) GetByBookIDAndChapterNum(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	var chapter bookstore.Chapter
	filter := bson.M{"book_id": bookID, "chapter_num": chapterNum}
	err := r.collection.FindOne(ctx, filter).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

// GetByTitle 根据标题获取章节
func (r *MongoChapterRepository) GetByTitle(ctx context.Context, title string, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	filter := bson.M{"title": bson.M{"$regex": title, "$options": "i"}}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// GetFreeChapters 获取免费章节
func (r *MongoChapterRepository) GetFreeChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	filter := bson.M{"book_id": bookID, "is_free": true}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// GetPaidChapters 获取付费章节
func (r *MongoChapterRepository) GetPaidChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	filter := bson.M{"book_id": bookID, "is_free": false}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// GetPublishedChapters 获取已发布章节
func (r *MongoChapterRepository) GetPublishedChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	filter := bson.M{
		"book_id":      bookID,
		"publish_time": bson.M{"$lte": time.Now()},
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// Search 搜索章节
func (r *MongoChapterRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"content": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// SearchByFilter 根据过滤器搜索章节
func (r *MongoChapterRepository) SearchByFilter(ctx context.Context, filter *BookstoreInterface.ChapterFilter) ([]*bookstore.Chapter, error) {
	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	// 构建排序
	sortField := "chapter_num"
	sortOrder := 1
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// 构建查询条件
	query := bson.M{}

	if filter.BookID != nil {
		query["book_id"] = *filter.BookID
	}
	if filter.Title != "" {
		query["title"] = bson.M{"$regex": filter.Title, "$options": "i"}
	}
	if filter.IsFree != nil {
		query["is_free"] = *filter.IsFree
	}
	if filter.MinChapterNum != nil {
		if query["chapter_num"] == nil {
			query["chapter_num"] = bson.M{}
		}
		query["chapter_num"].(bson.M)["$gte"] = *filter.MinChapterNum
	}
	if filter.MaxChapterNum != nil {
		if query["chapter_num"] == nil {
			query["chapter_num"] = bson.M{}
		}
		query["chapter_num"].(bson.M)["$lte"] = *filter.MaxChapterNum
	}
	if filter.MinWordCount != nil {
		if query["word_count"] == nil {
			query["word_count"] = bson.M{}
		}
		query["word_count"].(bson.M)["$gte"] = *filter.MinWordCount
	}
	if filter.MaxWordCount != nil {
		if query["word_count"] == nil {
			query["word_count"] = bson.M{}
		}
		query["word_count"].(bson.M)["$lte"] = *filter.MaxWordCount
	}
	if filter.IsPublished != nil {
		if *filter.IsPublished {
			query["publish_time"] = bson.M{"$lte": time.Now()}
		} else {
			query["publish_time"] = bson.M{"$gt": time.Now()}
		}
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// CountByBookID 根据书籍ID统计章节数量
func (r *MongoChapterRepository) CountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	filter := bson.M{"book_id": bookID}
	return r.collection.CountDocuments(ctx, filter)
}

// CountFreeChapters 统计免费章节数量
func (r *MongoChapterRepository) CountFreeChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	filter := bson.M{"book_id": bookID, "is_free": true}
	return r.collection.CountDocuments(ctx, filter)
}

// CountPaidChapters 统计付费章节数量
func (r *MongoChapterRepository) CountPaidChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	filter := bson.M{"book_id": bookID, "is_free": false}
	return r.collection.CountDocuments(ctx, filter)
}

// GetChapterRange 获取章节范围
func (r *MongoChapterRepository) GetChapterRange(ctx context.Context, bookID primitive.ObjectID, startChapter, endChapter int) ([]*bookstore.Chapter, error) {
	filter := bson.M{
		"book_id": bookID,
		"chapter_num": bson.M{
			"$gte": startChapter,
			"$lte": endChapter,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []*bookstore.Chapter
	for cursor.Next(ctx) {
		var chapter bookstore.Chapter
		if err := cursor.Decode(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, cursor.Err()
}

// CountPublishedChapters 统计已发布章节数量
func (r *MongoChapterRepository) CountPublishedChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"book_id":      bookID,
		"is_published": true,
	}
	return r.collection.CountDocuments(ctx, filter)
}

// GetTotalWordCount 获取书籍总字数
func (r *MongoChapterRepository) GetTotalWordCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"book_id": bookID}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$word_count"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}

	return 0, nil
}

// GetPreviousChapter 获取上一章节
func (r *MongoChapterRepository) GetPreviousChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	var chapter bookstore.Chapter
	filter := bson.M{
		"book_id":     bookID,
		"chapter_num": bson.M{"$lt": chapterNum},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: -1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

// GetNextChapter 获取下一章节
func (r *MongoChapterRepository) GetNextChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	var chapter bookstore.Chapter
	filter := bson.M{
		"book_id":     bookID,
		"chapter_num": bson.M{"$gt": chapterNum},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

// GetFirstChapter 获取第一章节
func (r *MongoChapterRepository) GetFirstChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error) {
	var chapter bookstore.Chapter
	filter := bson.M{"book_id": bookID}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

// GetLastChapter 获取最后章节
func (r *MongoChapterRepository) GetLastChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error) {
	var chapter bookstore.Chapter
	filter := bson.M{"book_id": bookID}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: -1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

// BatchUpdatePrice 批量更新章节价格
func (r *MongoChapterRepository) BatchUpdatePrice(ctx context.Context, chapterIDs []primitive.ObjectID, price float64) error {
	filter := bson.M{"_id": bson.M{"$in": chapterIDs}}
	update := bson.M{
		"$set": bson.M{
			"price":      price,
			"is_free":    price == 0,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchUpdatePublishTime 批量更新发布时间
func (r *MongoChapterRepository) BatchUpdatePublishTime(ctx context.Context, chapterIDs []primitive.ObjectID, publishTime time.Time) error {
	filter := bson.M{"_id": bson.M{"$in": chapterIDs}}
	update := bson.M{
		"$set": bson.M{
			"publish_time": publishTime,
			"updated_at":   time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchUpdateFreeStatus 批量更新章节免费状态
func (r *MongoChapterRepository) BatchUpdateFreeStatus(ctx context.Context, chapterIDs []primitive.ObjectID, isFree bool) error {
	filter := bson.M{"_id": bson.M{"$in": chapterIDs}}
	update := bson.M{
		"$set": bson.M{
			"is_free":    isFree,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// BatchDelete 批量删除章节
func (r *MongoChapterRepository) BatchDelete(ctx context.Context, chapterIDs []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": chapterIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// DeleteByBookID 根据书籍ID删除所有章节
func (r *MongoChapterRepository) DeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"book_id": bookID})
	return err
}

// Transaction 执行事务
func (r *MongoChapterRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
func (r *MongoChapterRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}
