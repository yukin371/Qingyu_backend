package reader

import (
	"Qingyu_backend/models/reader"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/repository/mongodb/base"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoReadingProgressRepository 阅读进度仓储MongoDB实现
type MongoReadingProgressRepository struct {
	*base.BaseMongoRepository  // 嵌入基类，继承ID转换和通用CRUD方法喵~
	db *mongo.Database
}

// NewMongoReadingProgressRepository 创建阅读进度仓储实例
func NewMongoReadingProgressRepository(db *mongo.Database) *MongoReadingProgressRepository {
	return &MongoReadingProgressRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "reading_progress"),
		db:                 db,
	}
}

// Create 创建阅读进度
func (r *MongoReadingProgressRepository) Create(ctx context.Context, progress *reader.ReadingProgress) error {
	if progress.ID.IsZero() {
		progress.ID = primitive.NewObjectID()
	}
	progress.CreatedAt = time.Now()
	progress.UpdatedAt = time.Now()
	progress.LastReadAt = time.Now()

	_, err := r.GetCollection().InsertOne(ctx, progress)
	if err != nil {
		return fmt.Errorf("创建阅读进度失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取阅读进度
func (r *MongoReadingProgressRepository) GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error) {
	objectID, err := r.ParseID(id)  // 使用基类的ParseID方法喵~
	if err != nil {
		return nil, err
	}

	var progress reader.ReadingProgress
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&progress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("阅读进度不存在")
		}
		return nil, fmt.Errorf("查询阅读进度失败: %w", err)
	}

	return &progress, nil
}

// Update 更新阅读进度
func (r *MongoReadingProgressRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
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
		return fmt.Errorf("更新阅读进度失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("阅读进度不存在")
	}

	return nil
}

// Delete 删除阅读进度
func (r *MongoReadingProgressRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return err
	}

	result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除阅读进度失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("阅读进度不存在")
	}

	return nil
}

// GetByUserAndBook 获取用户对特定书籍的阅读进度
func (r *MongoReadingProgressRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}
	bookOID, err := r.ParseID(bookID)
	if err != nil {
		return nil, err
	}

	var progress reader.ReadingProgress
	err = r.GetCollection().FindOne(ctx, bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	}).Decode(&progress)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没有阅读记录，返回nil而不是错误
		}
		return nil, fmt.Errorf("查询阅读进度失败: %w", err)
	}

	return &progress, nil
}

// GetByUser 获取用户的所有阅读进度
func (r *MongoReadingProgressRepository) GetByUser(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_id": userOID}
	opts := options.Find().SetSort(bson.D{{Key: "last_read_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询阅读进度失败: %w", err)
	}
	defer cursor.Close(ctx)

	var progresses []*reader.ReadingProgress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, fmt.Errorf("解析阅读进度数据失败: %w", err)
	}

	return progresses, nil
}

// GetRecentReadingByUser 获取用户最近的阅读记录
func (r *MongoReadingProgressRepository) GetRecentReadingByUser(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_id": userOID}
	opts := options.Find().
		SetSort(bson.D{{Key: "last_read_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询最近阅读记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var progresses []*reader.ReadingProgress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, fmt.Errorf("解析阅读进度数据失败: %w", err)
	}

	return progresses, nil
}

// SaveProgress 保存或更新阅读进度
func (r *MongoReadingProgressRepository) SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("无效的用户ID: %w", err)
	}
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("无效的书籍ID: %w", err)
	}
	chapterOID, err := r.ParseID(chapterID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	}

	update := bson.M{
		"$set": bson.M{
			"chapter_id":   chapterOID,
			"progress":     progress,
			"last_read_at": time.Now(),
			"updated_at":   time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":          primitive.NewObjectID(),
			"reading_time": int64(0),
			"created_at":   time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.GetCollection().UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("保存阅读进度失败: %w", err)
	}

	return nil
}

// UpdateReadingTime 更新阅读时长
func (r *MongoReadingProgressRepository) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("无效的用户ID: %w", err)
	}
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("无效的书籍ID: %w", err)
	}

	filter := bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	}

	update := bson.M{
		"$inc": bson.M{
			"reading_time": duration,
		},
		"$set": bson.M{
			"last_read_at": time.Now(),
			"updated_at":   time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新阅读时长失败: %w", err)
	}

	if result.MatchedCount == 0 {
		// 如果没有记录，创建新记录
		progress := &reader.ReadingProgress{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			UserID:           userOID,
			BookID:           bookOID,
			ReadingTime:      duration,
			Progress:         0,
			LastReadAt:       time.Now(),
		}
		return r.Create(ctx, progress)
	}

	return nil
}

// UpdateLastReadAt 更新最后阅读时间
func (r *MongoReadingProgressRepository) UpdateLastReadAt(ctx context.Context, userID, bookID string) error {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("无效的用户ID: %w", err)
	}
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("无效的书籍ID: %w", err)
	}

	filter := bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	}

	update := bson.M{
		"$set": bson.M{
			"last_read_at": time.Now(),
			"updated_at":   time.Now(),
		},
	}

	_, err = r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新最后阅读时间失败: %w", err)
	}

	return nil
}

// BatchUpdateProgress 批量更新阅读进度
func (r *MongoReadingProgressRepository) BatchUpdateProgress(ctx context.Context, progresses []*reader.ReadingProgress) error {
	if len(progresses) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, len(progresses))
	for i, progress := range progresses {
		filter := bson.M{
			"user_id": progress.UserID,
			"book_id": progress.BookID,
		}

		update := bson.M{
			"$set": bson.M{
				"chapter_id":   progress.ChapterID,
				"progress":     progress.Progress,
				"reading_time": progress.ReadingTime,
				"last_read_at": progress.LastReadAt,
				"updated_at":   time.Now(),
			},
			"$setOnInsert": bson.M{
				"_id":        progress.ID,
				"created_at": time.Now(),
			},
		}

		models[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.GetCollection().BulkWrite(ctx, models, opts)
	if err != nil {
		return fmt.Errorf("批量更新阅读进度失败: %w", err)
	}

	return nil
}

// GetTotalReadingTime 获取用户总阅读时长
func (r *MongoReadingProgressRepository) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userOID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$reading_time"},
		}}},
	}

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("统计总阅读时长失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("解析统计结果失败: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Total, nil
}

// GetReadingTimeByBook 获取用户阅读某本书的时长
func (r *MongoReadingProgressRepository) GetReadingTimeByBook(ctx context.Context, userID, bookID string) (int64, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return 0, err
	}

	var progress reader.ReadingProgress
	err = r.GetCollection().FindOne(ctx, bson.M{
		"user_id": userOID,
		"book_id": bookOID,
	}).Decode(&progress)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, fmt.Errorf("查询阅读时长失败: %w", err)
	}

	return progress.ReadingTime, nil
}

// GetReadingTimeByPeriod 获取用户在指定时间段的阅读时长
func (r *MongoReadingProgressRepository) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	filter := bson.M{
		"user_id": userOID,
		"last_read_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$reading_time"},
		}}},
	}

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("统计时间段阅读时长失败: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("解析统计结果失败: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Total, nil
}

// CountReadingBooks 统计用户正在阅读的书籍数量
func (r *MongoReadingProgressRepository) CountReadingBooks(ctx context.Context, userID string) (int64, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.GetCollection().CountDocuments(ctx, bson.M{"user_id": userOID})
	if err != nil {
		return 0, fmt.Errorf("统计阅读书籍数失败: %w", err)
	}

	return count, nil
}

// GetReadingHistory 获取阅读历史
func (r *MongoReadingProgressRepository) GetReadingHistory(ctx context.Context, userID string, limit, offset int) ([]*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_id": userOID}
	opts := options.Find().
		SetSort(bson.D{{Key: "last_read_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询阅读历史失败: %w", err)
	}
	defer cursor.Close(ctx)

	var progresses []*reader.ReadingProgress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, fmt.Errorf("解析阅读历史数据失败: %w", err)
	}

	return progresses, nil
}

// GetUnfinishedBooks 获取未读完的书籍
func (r *MongoReadingProgressRepository) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id":  userOID,
		"progress": bson.M{"$lt": 1.0}, // 进度小于100%
	}
	opts := options.Find().SetSort(bson.D{{Key: "last_read_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询未读完书籍失败: %w", err)
	}
	defer cursor.Close(ctx)

	var progresses []*reader.ReadingProgress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, fmt.Errorf("解析数据失败: %w", err)
	}

	return progresses, nil
}

// GetFinishedBooks 获取已读完的书籍
func (r *MongoReadingProgressRepository) GetFinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id":  userOID,
		"progress": bson.M{"$gte": 1.0}, // 进度达到100%
	}
	opts := options.Find().SetSort(bson.D{{Key: "last_read_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询已读完书籍失败: %w", err)
	}
	defer cursor.Close(ctx)

	var progresses []*reader.ReadingProgress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, fmt.Errorf("解析数据失败: %w", err)
	}

	return progresses, nil
}

// SyncProgress 同步进度数据
func (r *MongoReadingProgressRepository) SyncProgress(ctx context.Context, userID string, progresses []*reader.ReadingProgress) error {
	return r.BatchUpdateProgress(ctx, progresses)
}

// GetProgressesByUser 获取用户在指定时间后更新的进度
func (r *MongoReadingProgressRepository) GetProgressesByUser(ctx context.Context, userID string, updatedAfter time.Time) ([]*reader.ReadingProgress, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id":    userOID,
		"updated_at": bson.M{"$gt": updatedAfter},
	}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询进度数据失败: %w", err)
	}
	defer cursor.Close(ctx)

	var progresses []*reader.ReadingProgress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, fmt.Errorf("解析进度数据失败: %w", err)
	}

	return progresses, nil
}

// DeleteOldProgress 删除旧的阅读进度
func (r *MongoReadingProgressRepository) DeleteOldProgress(ctx context.Context, beforeTime time.Time) error {
	filter := bson.M{
		"last_read_at": bson.M{"$lt": beforeTime},
	}

	_, err := r.GetCollection().DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除旧阅读进度失败: %w", err)
	}

	return nil
}

// DeleteByBook 删除某本书的所有阅读进度
func (r *MongoReadingProgressRepository) DeleteByBook(ctx context.Context, bookID string) error {
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return fmt.Errorf("无效的书籍ID: %w", err)
	}

	filter := bson.M{"book_id": bookOID}
	_, err = r.GetCollection().DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除书籍阅读进度失败: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoReadingProgressRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
