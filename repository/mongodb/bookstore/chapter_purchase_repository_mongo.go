package mongodb

import (
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/mongodb/base"
	"context"
	"errors"
	"fmt"
	"time"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoChapterPurchaseRepository MongoDB章节购买仓储实现。
type MongoChapterPurchaseRepository struct {
	*base.BaseMongoRepository
	db                       *mongo.Database
	chapterPurchaseColl      *mongo.Collection
	chapterPurchaseBatchColl *mongo.Collection
	bookPurchaseColl         *mongo.Collection
}

// NewMongoChapterPurchaseRepository 创建章节购买仓储。
func NewMongoChapterPurchaseRepository(client *mongo.Client, database string) BookstoreInterface.ChapterPurchaseRepository {
	db := client.Database(database)
	return &MongoChapterPurchaseRepository{
		BaseMongoRepository:      base.NewBaseMongoRepository(db, "chapter_purchases"),
		db:                       db,
		chapterPurchaseColl:      db.Collection("chapter_purchases"),
		chapterPurchaseBatchColl: db.Collection("chapter_purchase_batches"),
		bookPurchaseColl:         db.Collection("book_purchases"),
	}
}

func (r *MongoChapterPurchaseRepository) parseID(id string) (primitive.ObjectID, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("无效ID: %w", err)
	}
	return objectID, nil
}

func (r *MongoChapterPurchaseRepository) paginate(page, pageSize int) (*options.FindOptions, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSort(bson.D{{Key: "purchase_time", Value: -1}, {Key: "_id", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))
	return opts, skip
}

func (r *MongoChapterPurchaseRepository) chapterFilterByUser(userID string) (bson.M, error) {
	userOID, err := r.parseID(userID)
	if err != nil {
		return nil, err
	}
	return bson.M{"user_id": userOID}, nil
}

func (r *MongoChapterPurchaseRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

func (r *MongoChapterPurchaseRepository) Create(ctx context.Context, purchase *bookstore.ChapterPurchase) error {
	if purchase == nil {
		return errors.New("purchase cannot be nil")
	}
	if purchase.ID.IsZero() {
		purchase.ID = primitive.NewObjectID()
	}
	purchase.BeforeCreate()

	_, err := r.chapterPurchaseColl.InsertOne(ctx, purchase)
	if err != nil {
		return fmt.Errorf("创建章节购买记录失败: %w", err)
	}
	return nil
}

func (r *MongoChapterPurchaseRepository) GetByID(ctx context.Context, id string) (*bookstore.ChapterPurchase, error) {
	objectID, err := r.parseID(id)
	if err != nil {
		return nil, err
	}

	var purchase bookstore.ChapterPurchase
	err = r.chapterPurchaseColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&purchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询章节购买记录失败: %w", err)
	}

	return &purchase, nil
}

func (r *MongoChapterPurchaseRepository) GetByUserAndChapter(ctx context.Context, userID, chapterID string) (*bookstore.ChapterPurchase, error) {
	userOID, err := r.parseID(userID)
	if err != nil {
		return nil, err
	}
	chapterOID, err := r.parseID(chapterID)
	if err != nil {
		return nil, err
	}

	var purchase bookstore.ChapterPurchase
	err = r.chapterPurchaseColl.FindOne(ctx, bson.M{
		"user_id":    userOID,
		"chapter_id": chapterOID,
	}).Decode(&purchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询章节购买记录失败: %w", err)
	}

	return &purchase, nil
}

func (r *MongoChapterPurchaseRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := r.parseID(id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()
	result, err := r.chapterPurchaseColl.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新章节购买记录失败: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("chapter purchase not found")
	}
	return nil
}

func (r *MongoChapterPurchaseRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.parseID(id)
	if err != nil {
		return err
	}

	result, err := r.chapterPurchaseColl.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除章节购买记录失败: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("chapter purchase not found")
	}
	return nil
}

func (r *MongoChapterPurchaseRepository) getChapterPurchases(ctx context.Context, filter bson.M, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	opts, _ := r.paginate(page, pageSize)
	cursor, err := r.chapterPurchaseColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询章节购买记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var purchases []*bookstore.ChapterPurchase
	if err := cursor.All(ctx, &purchases); err != nil {
		return nil, 0, fmt.Errorf("解析章节购买记录失败: %w", err)
	}

	total, err := r.chapterPurchaseColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计章节购买记录失败: %w", err)
	}
	return purchases, total, nil
}

func (r *MongoChapterPurchaseRepository) GetByUser(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	return r.getChapterPurchases(ctx, filter, page, pageSize)
}

func (r *MongoChapterPurchaseRepository) GetByUserAndBook(ctx context.Context, userID, bookID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	bookOID, err := r.parseID(bookID)
	if err != nil {
		return nil, 0, err
	}
	filter["book_id"] = bookOID
	return r.getChapterPurchases(ctx, filter, page, pageSize)
}

func (r *MongoChapterPurchaseRepository) CreateBatch(ctx context.Context, batch *bookstore.ChapterPurchaseBatch) error {
	if batch == nil {
		return errors.New("batch cannot be nil")
	}
	if batch.ID.IsZero() {
		batch.ID = primitive.NewObjectID()
	}
	batch.BeforeCreate()

	_, err := r.chapterPurchaseBatchColl.InsertOne(ctx, batch)
	if err != nil {
		return fmt.Errorf("创建批量购买记录失败: %w", err)
	}
	return nil
}

func (r *MongoChapterPurchaseRepository) GetBatchByID(ctx context.Context, id string) (*bookstore.ChapterPurchaseBatch, error) {
	objectID, err := r.parseID(id)
	if err != nil {
		return nil, err
	}

	var batch bookstore.ChapterPurchaseBatch
	err = r.chapterPurchaseBatchColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&batch)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询批量购买记录失败: %w", err)
	}
	return &batch, nil
}

func (r *MongoChapterPurchaseRepository) getBatches(ctx context.Context, filter bson.M, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error) {
	opts, _ := r.paginate(page, pageSize)
	cursor, err := r.chapterPurchaseBatchColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询批量购买记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var batches []*bookstore.ChapterPurchaseBatch
	if err := cursor.All(ctx, &batches); err != nil {
		return nil, 0, fmt.Errorf("解析批量购买记录失败: %w", err)
	}

	total, err := r.chapterPurchaseBatchColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计批量购买记录失败: %w", err)
	}
	return batches, total, nil
}

func (r *MongoChapterPurchaseRepository) GetBatchesByUser(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	return r.getBatches(ctx, filter, page, pageSize)
}

func (r *MongoChapterPurchaseRepository) GetBatchesByUserAndBook(ctx context.Context, userID, bookID string, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	bookOID, err := r.parseID(bookID)
	if err != nil {
		return nil, 0, err
	}
	filter["book_id"] = bookOID
	return r.getBatches(ctx, filter, page, pageSize)
}

func (r *MongoChapterPurchaseRepository) CreateBookPurchase(ctx context.Context, purchase *bookstore.BookPurchase) error {
	if purchase == nil {
		return errors.New("book purchase cannot be nil")
	}
	if purchase.ID.IsZero() {
		purchase.ID = primitive.NewObjectID()
	}
	purchase.BeforeCreate()

	_, err := r.bookPurchaseColl.InsertOne(ctx, purchase)
	if err != nil {
		return fmt.Errorf("创建全书购买记录失败: %w", err)
	}
	return nil
}

func (r *MongoChapterPurchaseRepository) GetBookPurchaseByID(ctx context.Context, id string) (*bookstore.BookPurchase, error) {
	objectID, err := r.parseID(id)
	if err != nil {
		return nil, err
	}

	var purchase bookstore.BookPurchase
	err = r.bookPurchaseColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&purchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询全书购买记录失败: %w", err)
	}
	return &purchase, nil
}

func (r *MongoChapterPurchaseRepository) GetBookPurchaseByUserAndBook(ctx context.Context, userID, bookID string) (*bookstore.BookPurchase, error) {
	userOID, err := r.parseID(userID)
	if err != nil {
		return nil, err
	}
	bookOID, err := r.parseID(bookID)
	if err != nil {
		return nil, err
	}

	var purchase bookstore.BookPurchase
	err = r.bookPurchaseColl.FindOne(ctx, bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	}).Decode(&purchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询全书购买记录失败: %w", err)
	}
	return &purchase, nil
}

func (r *MongoChapterPurchaseRepository) GetBookPurchasesByUser(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.BookPurchase, int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	opts, _ := r.paginate(page, pageSize)
	cursor, err := r.bookPurchaseColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询全书购买记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var purchases []*bookstore.BookPurchase
	if err := cursor.All(ctx, &purchases); err != nil {
		return nil, 0, fmt.Errorf("解析全书购买记录失败: %w", err)
	}

	total, err := r.bookPurchaseColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计全书购买记录失败: %w", err)
	}
	return purchases, total, nil
}

func (r *MongoChapterPurchaseRepository) CheckUserPurchasedChapter(ctx context.Context, userID, chapterID string) (bool, error) {
	userOID, err := r.parseID(userID)
	if err != nil {
		return false, err
	}
	chapterOID, err := r.parseID(chapterID)
	if err != nil {
		return false, err
	}

	count, err := r.chapterPurchaseColl.CountDocuments(ctx, bson.M{
		"user_id":    userOID,
		"chapter_id": chapterOID,
	})
	if err != nil {
		return false, fmt.Errorf("检查章节购买状态失败: %w", err)
	}
	return count > 0, nil
}

func (r *MongoChapterPurchaseRepository) CheckUserPurchasedBook(ctx context.Context, userID, bookID string) (bool, error) {
	userOID, err := r.parseID(userID)
	if err != nil {
		return false, err
	}
	bookOID, err := r.parseID(bookID)
	if err != nil {
		return false, err
	}

	count, err := r.bookPurchaseColl.CountDocuments(ctx, bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	})
	if err != nil {
		return false, fmt.Errorf("检查全书购买状态失败: %w", err)
	}
	return count > 0, nil
}

func (r *MongoChapterPurchaseRepository) GetPurchasedChapterIDs(ctx context.Context, userID, bookID string) ([]string, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, err
	}
	if bookID != "" {
		bookOID, err := r.parseID(bookID)
		if err != nil {
			return nil, err
		}
		filter["book_id"] = bookOID
	}

	cursor, err := r.chapterPurchaseColl.Find(ctx, filter, options.Find().SetProjection(bson.M{"chapter_id": 1}))
	if err != nil {
		return nil, fmt.Errorf("查询已购章节ID失败: %w", err)
	}
	defer cursor.Close(ctx)

	type item struct {
		ChapterID primitive.ObjectID `bson:"chapter_id"`
	}

	var rows []item
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, fmt.Errorf("解析已购章节ID失败: %w", err)
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		if !row.ChapterID.IsZero() {
			ids = append(ids, row.ChapterID.Hex())
		}
	}
	return ids, nil
}

func (r *MongoChapterPurchaseRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return 0, err
	}
	chapterCount, err := r.chapterPurchaseColl.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计章节购买记录失败: %w", err)
	}
	batchCount, err := r.chapterPurchaseBatchColl.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计批量购买记录失败: %w", err)
	}
	bookCount, err := r.bookPurchaseColl.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计全书购买记录失败: %w", err)
	}
	return chapterCount + batchCount + bookCount, nil
}

func (r *MongoChapterPurchaseRepository) CountByUserAndBook(ctx context.Context, userID, bookID string) (int64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return 0, err
	}
	bookOID, err := r.parseID(bookID)
	if err != nil {
		return 0, err
	}
	filter["book_id"] = bookOID
	chapterCount, err := r.chapterPurchaseColl.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计章节购买记录失败: %w", err)
	}
	batchCount, err := r.chapterPurchaseBatchColl.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计批量购买记录失败: %w", err)
	}
	bookCount, err := r.bookPurchaseColl.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计全书购买记录失败: %w", err)
	}
	return chapterCount + batchCount + bookCount, nil
}

func (r *MongoChapterPurchaseRepository) GetTotalSpentByUser(ctx context.Context, userID string) (float64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return 0, err
	}
	return r.getTotalSpent(ctx, filter)
}

func (r *MongoChapterPurchaseRepository) GetTotalSpentByUserAndBook(ctx context.Context, userID, bookID string) (float64, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return 0, err
	}
	bookOID, err := r.parseID(bookID)
	if err != nil {
		return 0, err
	}
	filter["book_id"] = bookOID
	return r.getTotalSpent(ctx, filter)
}

func (r *MongoChapterPurchaseRepository) getTotalSpent(ctx context.Context, filter bson.M) (float64, error) {
	sumField := func(coll *mongo.Collection, field string) (float64, error) {
		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$" + field}}}},
		}
		cursor, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return 0, err
		}
		defer cursor.Close(ctx)

		type result struct {
			Total float64 `bson:"total"`
		}

		var rows []result
		if err := cursor.All(ctx, &rows); err != nil {
			return 0, err
		}
		if len(rows) == 0 {
			return 0, nil
		}
		return rows[0].Total, nil
	}

	chapterTotal, err := sumField(r.chapterPurchaseColl, "price")
	if err != nil {
		return 0, fmt.Errorf("统计章节购买金额失败: %w", err)
	}
	batchTotal, err := sumField(r.chapterPurchaseBatchColl, "total_price")
	if err != nil {
		return 0, fmt.Errorf("统计批量购买金额失败: %w", err)
	}
	bookTotal, err := sumField(r.bookPurchaseColl, "total_price")
	if err != nil {
		return 0, fmt.Errorf("统计全书购买金额失败: %w", err)
	}
	return chapterTotal + batchTotal + bookTotal, nil
}

func (r *MongoChapterPurchaseRepository) GetPurchasesByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*bookstore.ChapterPurchase, error) {
	filter, err := r.chapterFilterByUser(userID)
	if err != nil {
		return nil, err
	}
	filter["purchase_time"] = bson.M{
		"$gte": startTime,
		"$lte": endTime,
	}
	cursor, err := r.chapterPurchaseColl.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "purchase_time", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("按时间范围查询章节购买记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var purchases []*bookstore.ChapterPurchase
	if err := cursor.All(ctx, &purchases); err != nil {
		return nil, fmt.Errorf("解析章节购买记录失败: %w", err)
	}
	return purchases, nil
}

func (r *MongoChapterPurchaseRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	})
	return err
}
