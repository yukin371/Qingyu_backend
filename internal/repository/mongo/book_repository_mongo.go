package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"qingteng-qa/internal/domain"
	"qingteng-qa/internal/repository"
	"qingteng-qa/pkg/errors"
)

// MongoBookRepository MongoDB书籍Repository实现
// 实现BookRepository接口,提供MongoDB的数据访问操作
type MongoBookRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	tx         *MongoTransaction
}

// NewMongoBookRepository 创建MongoBookRepository实例
func NewMongoBookRepository(db *mongo.Database) *MongoBookRepository {
	return &MongoBookRepository{
		db:         db,
		collection: db.Collection("books"),
	}
}

// WithTransaction 返回事务版本的Repository
// 实现TransactionAware接口
func (r *MongoBookRepository) WithTransaction(tx repository.Transaction) interface{} {
	mongoTx := tx.(*MongoTransaction)
	return &MongoBookRepository{
		db:         r.db,
		collection: r.collection,
		tx:         mongoTx,
	}
}

// Create 创建书籍
func (r *MongoBookRepository) Create(ctx context.Context, book *domain.Book) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	// 设置创建时间
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now

	// 插入文档
	result, err := coll.InsertOne(ctx, book)
	if err != nil {
		// 检查是否是重复键错误 (如ISBN重复)
		if mongo.IsDuplicateKeyError(err) {
			return errors.Wrap(err, errors.ErrCodeDuplicateKey, "书籍已存在")
		}
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "创建书籍失败")
	}

	// 设置生成的ID
	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取书籍
func (r *MongoBookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	// 转换ID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInvalidID, "无效的书籍ID")
	}

	// 查询文档
	var book domain.Book
	err = coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(err, errors.ErrCodeBookNotFound, "书籍不存在")
		}
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseError, "查询书籍失败")
	}

	return &book, nil
}

// Update 更新书籍
func (r *MongoBookRepository) Update(ctx context.Context, book *domain.Book) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	// 设置更新时间
	book.UpdatedAt = time.Now()

	// 构建更新文档
	updateDoc := bson.M{
		"$set": bson.M{
			"title":        book.Title,
			"author":       book.Author,
			"isbn":         book.ISBN,
			"description":  book.Description,
			"category_id":  book.CategoryID,
			"tags":         book.Tags,
			"cover_image":  book.CoverImage,
			"price":        book.Price,
			"stock":        book.Stock,
			"status":       book.Status,
			"view_count":   book.ViewCount,
			"like_count":   book.LikeCount,
			"collect_count": book.CollectCount,
			"word_count":   book.WordCount,
			"publish_date": book.PublishDate,
			"publisher":    book.Publisher,
			"updated_at":   book.UpdatedAt,
		},
	}

	// 执行更新
	result, err := coll.UpdateByID(ctx, book.ID, updateDoc)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "更新书籍失败")
	}

	// 检查是否找到文档
	if result.MatchedCount == 0 {
		return errors.New(errors.ErrCodeBookNotFound, "书籍不存在")
	}

	return nil
}

// Delete 删除书籍 (软删除)
func (r *MongoBookRepository) Delete(ctx context.Context, id string) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	// 转换ID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInvalidID, "无效的书籍ID")
	}

	// 软删除: 更新状态为deleted
	update := bson.M{
		"$set": bson.M{
			"status":     domain.BookStatusDeleted,
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateByID(ctx, objectID, update)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "删除书籍失败")
	}

	if result.MatchedCount == 0 {
		return errors.New(errors.ErrCodeBookNotFound, "书籍不存在")
	}

	return nil
}

// List 查询书籍列表
func (r *MongoBookRepository) List(ctx context.Context, filter *repository.BookFilter, page, size int) ([]*domain.Book, int64, error) {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	// 构建查询条件
	bsonFilter := r.buildFilter(filter)

	// 计算总数
	total, err := coll.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeDatabaseError, "查询书籍总数失败")
	}

	// 构建查询选项
	opts := options.Find()
	opts.SetSkip(int64((page - 1) * size))
	opts.SetLimit(int64(size))

	// 设置排序
	if filter != nil && filter.SortBy != nil {
		sortOrder := 1 // 默认升序
		if filter.SortOrder != nil && *filter.SortOrder == "desc" {
			sortOrder = -1
		}
		opts.SetSort(bson.D{{*filter.SortBy, sortOrder}})
	} else {
		// 默认按创建时间倒序
		opts.SetSort(bson.D{{"created_at", -1}})
	}

	// 执行查询
	cursor, err := coll.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeDatabaseError, "查询书籍列表失败")
	}
	defer cursor.Close(ctx)

	// 解码结果
	var books []*domain.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeDatabaseError, "解码书籍列表失败")
	}

	return books, total, nil
}

// GetByISBN 根据ISBN获取书籍
func (r *MongoBookRepository) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	var book domain.Book
	err := coll.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // ISBN不存在,返回nil而不是错误
		}
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseError, "查询书籍失败")
	}

	return &book, nil
}

// UpdateStock 更新库存
func (r *MongoBookRepository) UpdateStock(ctx context.Context, bookID string, delta int) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	// 转换ID
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInvalidID, "无效的书籍ID")
	}

	// 使用$inc原子操作更新库存
	update := bson.M{
		"$inc": bson.M{
			"stock": delta,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateByID(ctx, objectID, update)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "更新库存失败")
	}

	if result.MatchedCount == 0 {
		return errors.New(errors.ErrCodeBookNotFound, "书籍不存在")
	}

	// 检查库存是否变为负数 (MongoDB不会自动阻止)
	book, err := r.GetByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book.Stock < 0 {
		// 回滚库存
		r.UpdateStock(ctx, bookID, -delta)
		return errors.New(errors.ErrCodeBookOutOfStock, "库存不足")
	}

	return nil
}

// IncrementView 增加浏览次数
func (r *MongoBookRepository) IncrementView(ctx context.Context, bookID string) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInvalidID, "无效的书籍ID")
	}

	update := bson.M{
		"$inc": bson.M{
			"view_count": 1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateByID(ctx, objectID, update)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "更新浏览次数失败")
	}

	if result.MatchedCount == 0 {
		return errors.New(errors.ErrCodeBookNotFound, "书籍不存在")
	}

	return nil
}

// IncrementLike 增加点赞数
func (r *MongoBookRepository) IncrementLike(ctx context.Context, bookID string) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInvalidID, "无效的书籍ID")
	}

	update := bson.M{
		"$inc": bson.M{
			"like_count": 1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateByID(ctx, objectID, update)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "更新点赞数失败")
	}

	if result.MatchedCount == 0 {
		return errors.New(errors.ErrCodeBookNotFound, "书籍不存在")
	}

	return nil
}

// DecrementLike 减少点赞数
func (r *MongoBookRepository) DecrementLike(ctx context.Context, bookID string) error {
	coll := r.getCollection()
	ctx = r.getCTX(ctx)

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInvalidID, "无效的书籍ID")
	}

	update := bson.M{
		"$inc": bson.M{
			"like_count": -1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateByID(ctx, objectID, update)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "更新点赞数失败")
	}

	if result.MatchedCount == 0 {
		return errors.New(errors.ErrCodeBookNotFound, "书籍不存在")
	}

	return nil
}

// buildFilter 构建MongoDB查询过滤器
func (r *MongoBookRepository) buildFilter(filter *repository.BookFilter) bson.M {
	bsonFilter := bson.M{}

	if filter == nil {
		// 默认只查询未删除的书籍
		bsonFilter["status"] = bson.M{"$ne": domain.BookStatusDeleted}
		return bsonFilter
	}

	// 关键词搜索 (标题、作者、简介)
	if filter.Keyword != nil && *filter.Keyword != "" {
		regex := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", *filter.Keyword),
			Options: "i", // 不区分大小写
		}
		bsonFilter["$or"] = []bson.M{
			{"title": bson.M{"$regex": regex}},
			{"author": bson.M{"$regex": regex}},
			{"description": bson.M{"$regex": regex}},
		}
	}

	// 分类ID过滤
	if filter.CategoryID != nil {
		objectID, err := primitive.ObjectIDFromHex(*filter.CategoryID)
		if err == nil {
			bsonFilter["category_id"] = objectID
		}
	}

	// 标签过滤 (ANY语义)
	if len(filter.Tags) > 0 {
		bsonFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	// 状态过滤
	if filter.Status != nil {
		bsonFilter["status"] = *filter.Status
	} else {
		// 默认不包含已删除的书籍
		bsonFilter["status"] = bson.M{"$ne": domain.BookStatusDeleted}
	}

	// 作者过滤
	if filter.Author != nil && *filter.Author != "" {
		regex := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", *filter.Author),
			Options: "i",
		}
		bsonFilter["author"] = bson.M{"$regex": regex}
	}

	// 价格范围过滤
	priceFilter := bson.M{}
	if filter.MinPrice != nil {
		priceFilter["$gte"] = *filter.MinPrice
	}
	if filter.MaxPrice != nil {
		priceFilter["$lte"] = *filter.MaxPrice
	}
	if len(priceFilter) > 0 {
		bsonFilter["price"] = priceFilter
	}

	// 库存过滤
	if filter.InStock != nil && *filter.InStock {
		bsonFilter["stock"] = bson.M{"$gt": 0}
	}

	return bsonFilter
}

// getCollection 获取集合 (支持事务)
func (r *MongoBookRepository) getCollection() *mongo.Collection {
	if r.tx != nil {
		return r.collection.WithSession(r.tx.ctx)
	}
	return r.collection
}

// getCTX 获取上下文 (支持事务)
func (r *MongoBookRepository) getCTX(ctx context.Context) context.Context {
	if r.tx != nil {
		return r.tx.Context()
	}
	return ctx
}

// MongoTransaction MongoDB事务 (占位,完整实现在事务管理模块)
type MongoTransaction struct {
	session mongo.Session
	ctx     context.Context
}

// Context 返回事务上下文
func (t *MongoTransaction) Context() context.Context {
	if t.ctx == nil {
		return context.Background()
	}
	return mongo.NewSessionContext(t.ctx, t.session)
}

// Commit 提交事务
func (t *MongoTransaction) Commit(ctx context.Context) error {
	if t.session == nil {
		return errors.New(errors.ErrCodeTransactionError, "事务未初始化")
	}
	if err := t.session.CommitTransaction(ctx); err != nil {
		return errors.Wrap(err, errors.ErrCodeTransactionError, "提交事务失败")
	}
	return nil
}

// Rollback 回滚事务
func (t *MongoTransaction) Rollback(ctx context.Context) error {
	if t.session == nil {
		return errors.New(errors.ErrCodeTransactionError, "事务未初始化")
	}
	if err := t.session.AbortTransaction(ctx); err != nil {
		return errors.Wrap(err, errors.ErrCodeTransactionError, "回滚事务失败")
	}
	return nil
}
