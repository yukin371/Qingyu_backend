package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/repository/mongodb/base"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	infra "Qingyu_backend/repository/interfaces/infrastructure"
)

// MongoRankingRepository MongoDB榜单仓储实现
type MongoRankingRepository struct {
	*base.BaseMongoRepository
	bookCollection *mongo.Collection
	client         *mongo.Client
}

// NewMongoRankingRepository 创建MongoDB榜单仓储实例
func NewMongoRankingRepository(client *mongo.Client, database string) BookstoreInterface.RankingRepository {
	db := client.Database(database)
	return &MongoRankingRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "ranking_items"),
		bookCollection:      db.Collection("books"),
		client:              client,
	}
}

// ========== 基础CRUD方法 ==========

// Create 创建榜单项
func (r *MongoRankingRepository) Create(ctx context.Context, item *bookstore2.RankingItem) error {
	if item == nil {
		return errors.New("ranking item cannot be nil")
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	result, err := r.GetCollection().InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to create ranking item: %w", err)
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取榜单项
func (r *MongoRankingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.RankingItem, error) {
	var item bookstore2.RankingItem
	err := r.GetCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ranking item: %w", err)
	}
	return &item, nil
}

// Update 更新榜单项
func (r *MongoRankingRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.GetCollection().UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("failed to update ranking item: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("ranking item not found")
	}

	return nil
}

// Delete 删除榜单项
func (r *MongoRankingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete ranking item: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("ranking item not found")
	}

	return nil
}

// List 查询榜单项列表
func (r *MongoRankingRepository) List(ctx context.Context, filter infra.Filter) ([]*bookstore2.RankingItem, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}

	cursor, err := r.GetCollection().Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list ranking items: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*bookstore2.RankingItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode ranking items: %w", err)
	}

	return items, nil
}

// Count 统计榜单项数量
func (r *MongoRankingRepository) Count(ctx context.Context, filter infra.Filter) (int64, error) {
	var query bson.M
	if filter != nil {
		query = bson.M(filter.GetConditions())
	} else {
		query = bson.M{}
	}
	return r.GetCollection().CountDocuments(ctx, query)
}

// Exists 检查榜单项是否存在
func (r *MongoRankingRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.GetCollection().CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, fmt.Errorf("failed to check ranking item existence: %w", err)
	}
	return count > 0, nil
}

// Health 健康检查
func (r *MongoRankingRepository) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// ========== 榜单特定查询方法 ==========

// GetByType 根据榜单类型获取榜单项
func (r *MongoRankingRepository) GetByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	filter := bson.M{
		"type":   rankingType,
		"period": period,
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "rank", Value: 1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	// 调试日志
	collName := r.GetCollection().Name()
	fmt.Printf("[DEBUG] GetByType - collection: %s, type: %s (%T), period: %s, filter: %v\n", collName, rankingType, rankingType, period, filter)

	// 打印集合中的所有文档以调试
	allDocs, _ := r.GetCollection().Find(ctx, bson.M{})
	var allItems []*bookstore2.RankingItem
	if allDocs != nil {
		defer allDocs.Close(ctx)
		allDocs.All(ctx, &allItems)
		fmt.Printf("[DEBUG] Total documents in collection: %d\n", len(allItems))
		for i, item := range allItems {
			fmt.Printf("[DEBUG]   Item %d: type=%s, period=%s\n", i, item.Type, item.Period)
		}
	}

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get rankings by type: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*bookstore2.RankingItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode ranking items: %w", err)
	}

	fmt.Printf("[DEBUG] GetByType - found %d items\n", len(items))
	return items, nil
}

// GetByTypeWithBooks 根据榜单类型获取榜单项（包含书籍信息）
func (r *MongoRankingRepository) GetByTypeWithBooks(ctx context.Context, rankingType bookstore2.RankingType, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	// 先获取榜单项
	items, err := r.GetByType(ctx, rankingType, period, limit, offset)
	if err != nil {
		return nil, err
	}

	// 如果没有榜单项，直接返回
	if len(items) == 0 {
		return items, nil
	}

	// 提取所有书籍ID
	bookIDs := make([]primitive.ObjectID, len(items))
	for i, item := range items {
		bookIDs[i] = item.BookID
	}

	// 批量查询书籍信息
	cursor, err := r.bookCollection.Find(ctx, bson.M{"_id": bson.M{"$in": bookIDs}})
	if err != nil {
		return items, nil // 即使查询书籍失败，也返回榜单项（不包含书籍信息）
	}
	defer cursor.Close(ctx)

	var books []*bookstore2.Book
	if err = cursor.All(ctx, &books); err != nil {
		return items, nil
	}

	// 创建书籍ID到书籍的映射
	bookMap := make(map[string]*bookstore2.Book)
	for _, book := range books {
		bookMap[book.ID.Hex()] = book
	}

	// 填充书籍信息
	for _, item := range items {
		if book, exists := bookMap[item.BookID.Hex()]; exists {
			item.Book = book
		}
	}

	return items, nil
}

// GetByBookID 根据书籍ID获取榜单项
func (r *MongoRankingRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, rankingType bookstore2.RankingType, period string) (*bookstore2.RankingItem, error) {
	filter := bson.M{
		"book_id": bookID,
		"type":    rankingType,
		"period":  period,
	}

	var item bookstore2.RankingItem
	err := r.GetCollection().FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ranking by book ID: %w", err)
	}

	return &item, nil
}

// GetByPeriod 根据周期获取榜单项
func (r *MongoRankingRepository) GetByPeriod(ctx context.Context, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	filter := bson.M{"period": period}

	opts := options.Find().
		SetSort(bson.D{{Key: "type", Value: 1}, {Key: "rank", Value: 1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get rankings by period: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*bookstore2.RankingItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode ranking items: %w", err)
	}

	return items, nil
}

// ========== 榜单统计方法 ==========

// GetRankingStats 获取榜单统计信息
func (r *MongoRankingRepository) GetRankingStats(ctx context.Context, rankingType bookstore2.RankingType, period string) (*bookstore2.RankingStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"type":   rankingType,
			"period": period,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"total_books":   bson.M{"$sum": 1},
			"total_views":   bson.M{"$sum": "$view_count"},
			"total_likes":   bson.M{"$sum": "$like_count"},
			"average_score": bson.M{"$avg": "$score"},
			"last_updated":  bson.M{"$max": "$updated_at"},
		}}},
	}

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate ranking stats: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalBooks   int64     `bson:"total_books"`
		TotalViews   int64     `bson:"total_views"`
		TotalLikes   int64     `bson:"total_likes"`
		AverageScore float64   `bson:"average_score"`
		LastUpdated  time.Time `bson:"last_updated"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode ranking stats: %w", err)
	}

	if len(results) == 0 {
		return &bookstore2.RankingStats{
			Type:   rankingType,
			Period: period,
		}, nil
	}

	return &bookstore2.RankingStats{
		Type:          rankingType,
		Period:        period,
		TotalBooks:    results[0].TotalBooks,
		TotalViews:    results[0].TotalViews,
		TotalLikes:    results[0].TotalLikes,
		AverageScore:  results[0].AverageScore,
		LastUpdatedAt: results[0].LastUpdated,
	}, nil
}

// CountByType 统计某类型榜单的数量
func (r *MongoRankingRepository) CountByType(ctx context.Context, rankingType bookstore2.RankingType, period string) (int64, error) {
	filter := bson.M{
		"type":   rankingType,
		"period": period,
	}
	return r.GetCollection().CountDocuments(ctx, filter)
}

// GetTopBooks 获取榜单前N本书
func (r *MongoRankingRepository) GetTopBooks(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	return r.GetByType(ctx, rankingType, period, limit, 0)
}

// ========== 榜单更新方法 ==========

// UpsertRankingItem 插入或更新榜单项
func (r *MongoRankingRepository) UpsertRankingItem(ctx context.Context, item *bookstore2.RankingItem) error {
	if item == nil {
		return errors.New("ranking item cannot be nil")
	}

	filter := bson.M{
		"book_id": item.BookID,
		"type":    item.Type,
		"period":  item.Period,
	}

	item.UpdatedAt = time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.GetCollection().UpdateOne(
		ctx,
		filter,
		bson.M{"$set": item},
		opts,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert ranking item: %w", err)
	}

	return nil
}

// BatchUpsertRankingItems 批量插入或更新榜单项
func (r *MongoRankingRepository) BatchUpsertRankingItems(ctx context.Context, items []*bookstore2.RankingItem) error {
	if len(items) == 0 {
		return nil
	}

	// 使用BulkWrite进行批量操作
	var operations []mongo.WriteModel
	now := time.Now()

	for _, item := range items {
		if item == nil {
			continue
		}

		filter := bson.M{
			"book_id": item.BookID,
			"type":    item.Type,
			"period":  item.Period,
		}

		item.UpdatedAt = now
		if item.CreatedAt.IsZero() {
			item.CreatedAt = now
		}

		update := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(bson.M{"$set": item}).
			SetUpsert(true)

		operations = append(operations, update)
	}

	if len(operations) == 0 {
		return nil
	}

	_, err := r.GetCollection().BulkWrite(ctx, operations)
	if err != nil {
		return fmt.Errorf("failed to batch upsert ranking items: %w", err)
	}

	return nil
}

// GetBooksForRanking 获取用于榜单计算的书籍数据（纯数据查询）
func (r *MongoRankingRepository) GetBooksForRanking(ctx context.Context) ([]*bookstore2.Book, error) {
	cursor, err := r.bookCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query books for ranking: %w", err)
	}
	defer cursor.Close(ctx)

	var books []*bookstore2.Book
	if err := cursor.All(ctx, &books); err != nil {
		return nil, fmt.Errorf("failed to decode books for ranking: %w", err)
	}
	return books, nil
}

// UpdateRankings 在事务中替换指定榜单数据
func (r *MongoRankingRepository) UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string, items []*bookstore2.RankingItem) error {
	return r.Transaction(ctx, func(txCtx context.Context) error {
		if err := r.DeleteByTypeAndPeriod(txCtx, rankingType, period); err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return r.BatchUpsertRankingItems(txCtx, items)
	})
}

// ========== 榜单维护方法 ==========

// DeleteByPeriod 删除指定周期的榜单
func (r *MongoRankingRepository) DeleteByPeriod(ctx context.Context, period string) error {
	_, err := r.GetCollection().DeleteMany(ctx, bson.M{"period": period})
	if err != nil {
		return fmt.Errorf("failed to delete rankings by period: %w", err)
	}
	return nil
}

// DeleteByType 删除指定类型的榜单
func (r *MongoRankingRepository) DeleteByType(ctx context.Context, rankingType bookstore2.RankingType) error {
	_, err := r.GetCollection().DeleteMany(ctx, bson.M{"type": rankingType})
	if err != nil {
		return fmt.Errorf("failed to delete rankings by type: %w", err)
	}
	return nil
}

// DeleteByTypeAndPeriod 删除指定类型和周期的榜单
func (r *MongoRankingRepository) DeleteByTypeAndPeriod(ctx context.Context, rankingType bookstore2.RankingType, period string) error {
	_, err := r.GetCollection().DeleteMany(ctx, bson.M{
		"type":   rankingType,
		"period": period,
	})
	if err != nil {
		return fmt.Errorf("failed to delete rankings by type and period: %w", err)
	}
	return nil
}

// DeleteExpiredRankings 删除过期的榜单
func (r *MongoRankingRepository) DeleteExpiredRankings(ctx context.Context, beforeDate time.Time) error {
	// 转换为周期字符串格式
	periodStr := beforeDate.Format("2006-01-02")

	_, err := r.GetCollection().DeleteMany(ctx, bson.M{
		"period": bson.M{"$lt": periodStr},
	})
	if err != nil {
		return fmt.Errorf("failed to delete expired rankings: %w", err)
	}
	return nil
}

func (r *MongoRankingRepository) calculateRankingItems(ctx context.Context, rankingType bookstore2.RankingType, period string) ([]*bookstore2.RankingItem, error) {
	now := time.Now()
	filter := bson.M{}
	if rankingType != bookstore2.RankingTypeRealtime &&
		rankingType != bookstore2.RankingTypeWeekly &&
		rankingType != bookstore2.RankingTypeMonthly &&
		rankingType != bookstore2.RankingTypeNewbie {
		return nil, fmt.Errorf("unsupported ranking type: %s", rankingType)
	}

	cursor, err := r.bookCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query books for ranking: %w", err)
	}
	defer cursor.Close(ctx)

	var books []*bookstore2.Book
	if err := cursor.All(ctx, &books); err != nil {
		return nil, fmt.Errorf("failed to decode books for ranking: %w", err)
	}

	items := make([]*bookstore2.RankingItem, 0, len(books))
	for _, book := range books {
		if book == nil || !rankingEligible(book, rankingType, now) {
			continue
		}
		items = append(items, &bookstore2.RankingItem{
			BookID:    book.ID,
			Type:      rankingType,
			Score:     rankingScore(book, rankingType),
			ViewCount: book.ViewCount,
			LikeCount: book.RatingCount,
			Period:    period,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].BookID.Hex() < items[j].BookID.Hex()
		}
		return items[i].Score > items[j].Score
	})

	if len(items) > 100 {
		items = items[:100]
	}
	for i, item := range items {
		item.Rank = i + 1
	}

	return items, nil
}

func rankingScore(book *bookstore2.Book, rankingType bookstore2.RankingType) float64 {
	switch rankingType {
	case bookstore2.RankingTypeRealtime:
		return float64(book.ViewCount)*0.7 + float64(book.RatingCount)*0.3
	case bookstore2.RankingTypeWeekly:
		return float64(book.ViewCount)*0.6 + float64(book.ChapterCount)*10
	case bookstore2.RankingTypeMonthly:
		return float64(book.ViewCount)*0.5 + float64(book.RatingCount)*0.3 + float64(book.WordCount)*0.0001
	case bookstore2.RankingTypeNewbie:
		return float64(book.ViewCount)*0.6 + float64(book.RatingCount)*0.4
	default:
		return 0
	}
}

func rankingEligible(book *bookstore2.Book, rankingType bookstore2.RankingType, now time.Time) bool {
	if book.Status != bookstore2.BookStatusOngoing {
		return false
	}
	if rankingType != bookstore2.RankingTypeNewbie {
		return true
	}
	if book.PublishedAt == nil {
		return false
	}
	return now.Sub(*book.PublishedAt) <= 30*24*time.Hour
}

// ========== 事务支持 ==========

// Transaction 执行事务
func (r *MongoRankingRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
