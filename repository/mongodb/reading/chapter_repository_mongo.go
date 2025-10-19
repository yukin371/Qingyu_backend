package reading

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/reading/reader"
)

// MongoChapterRepository 章节仓储MongoDB实现
type MongoChapterRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

// NewMongoChapterRepository 创建章节仓储实例
func NewMongoChapterRepository(db *mongo.Database) *MongoChapterRepository {
	return &MongoChapterRepository{
		collection: db.Collection("chapters"),
		db:         db,
	}
}

// Create 创建章节
func (r *MongoChapterRepository) Create(ctx context.Context, chapter *reader.Chapter) error {
	if chapter.ID == "" {
		chapter.ID = generateID()
	}
	chapter.CreatedAt = time.Now()
	chapter.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, chapter)
	if err != nil {
		return fmt.Errorf("创建章节失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取章节
func (r *MongoChapterRepository) GetByID(ctx context.Context, id string) (*reader.Chapter, error) {
	var chapter reader.Chapter
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("章节不存在")
		}
		return nil, fmt.Errorf("查询章节失败: %w", err)
	}

	return &chapter, nil
}

// Update 更新章节
func (r *MongoChapterRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新章节失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("章节不存在")
	}

	return nil
}

// Delete 删除章节
func (r *MongoChapterRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("删除章节失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("章节不存在")
	}

	return nil
}

// GetByBookID 获取书籍的所有章节
func (r *MongoChapterRepository) GetByBookID(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	filter := bson.M{
		"book_id": bookID,
		"status":  1, // 正常状态
	}

	opts := options.Find().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询章节列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var chapters []*reader.Chapter
	if err = cursor.All(ctx, &chapters); err != nil {
		return nil, fmt.Errorf("解析章节数据失败: %w", err)
	}

	return chapters, nil
}

// GetByBookIDWithPagination 分页获取书籍章节
func (r *MongoChapterRepository) GetByBookIDWithPagination(ctx context.Context, bookID string, limit, offset int64) ([]*reader.Chapter, error) {
	filter := bson.M{
		"book_id": bookID,
		"status":  1,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询章节列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var chapters []*reader.Chapter
	if err = cursor.All(ctx, &chapters); err != nil {
		return nil, fmt.Errorf("解析章节数据失败: %w", err)
	}

	return chapters, nil
}

// GetByChapterNum 根据章节号获取章节
func (r *MongoChapterRepository) GetByChapterNum(ctx context.Context, bookID string, chapterNum int) (*reader.Chapter, error) {
	var chapter reader.Chapter
	err := r.collection.FindOne(ctx, bson.M{
		"book_id":     bookID,
		"chapter_num": chapterNum,
		"status":      1,
	}).Decode(&chapter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("章节不存在")
		}
		return nil, fmt.Errorf("查询章节失败: %w", err)
	}

	return &chapter, nil
}

// GetPrevChapter 获取上一章
func (r *MongoChapterRepository) GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
	var chapter reader.Chapter
	filter := bson.M{
		"book_id":     bookID,
		"chapter_num": bson.M{"$lt": currentChapterNum},
		"status":      1,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: -1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 已经是第一章
		}
		return nil, fmt.Errorf("查询上一章失败: %w", err)
	}

	return &chapter, nil
}

// GetNextChapter 获取下一章
func (r *MongoChapterRepository) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
	var chapter reader.Chapter
	filter := bson.M{
		"book_id":     bookID,
		"chapter_num": bson.M{"$gt": currentChapterNum},
		"status":      1,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 已经是最后一章
		}
		return nil, fmt.Errorf("查询下一章失败: %w", err)
	}

	return &chapter, nil
}

// GetFirstChapter 获取第一章
func (r *MongoChapterRepository) GetFirstChapter(ctx context.Context, bookID string) (*reader.Chapter, error) {
	var chapter reader.Chapter
	filter := bson.M{
		"book_id": bookID,
		"status":  1,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("该书没有章节")
		}
		return nil, fmt.Errorf("查询第一章失败: %w", err)
	}

	return &chapter, nil
}

// GetLastChapter 获取最后一章
func (r *MongoChapterRepository) GetLastChapter(ctx context.Context, bookID string) (*reader.Chapter, error) {
	var chapter reader.Chapter
	filter := bson.M{
		"book_id": bookID,
		"status":  1,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: -1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("该书没有章节")
		}
		return nil, fmt.Errorf("查询最后一章失败: %w", err)
	}

	return &chapter, nil
}

// GetPublishedChapters 获取已发布的章节
func (r *MongoChapterRepository) GetPublishedChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	filter := bson.M{
		"book_id":      bookID,
		"status":       1,
		"publish_time": bson.M{"$lte": time.Now()},
	}

	opts := options.Find().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询已发布章节失败: %w", err)
	}
	defer cursor.Close(ctx)

	var chapters []*reader.Chapter
	if err = cursor.All(ctx, &chapters); err != nil {
		return nil, fmt.Errorf("解析章节数据失败: %w", err)
	}

	return chapters, nil
}

// GetVIPChapters 获取VIP章节
func (r *MongoChapterRepository) GetVIPChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	filter := bson.M{
		"book_id": bookID,
		"status":  1,
		"is_vip":  true,
	}

	opts := options.Find().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询VIP章节失败: %w", err)
	}
	defer cursor.Close(ctx)

	var chapters []*reader.Chapter
	if err = cursor.All(ctx, &chapters); err != nil {
		return nil, fmt.Errorf("解析章节数据失败: %w", err)
	}

	return chapters, nil
}

// GetFreeChapters 获取免费章节
func (r *MongoChapterRepository) GetFreeChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	filter := bson.M{
		"book_id": bookID,
		"status":  1,
		"is_vip":  false,
	}

	opts := options.Find().SetSort(bson.D{{Key: "chapter_num", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询免费章节失败: %w", err)
	}
	defer cursor.Close(ctx)

	var chapters []*reader.Chapter
	if err = cursor.All(ctx, &chapters); err != nil {
		return nil, fmt.Errorf("解析章节数据失败: %w", err)
	}

	return chapters, nil
}

// CountByBookID 统计书籍章节数
func (r *MongoChapterRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"book_id": bookID,
		"status":  1,
	})
	if err != nil {
		return 0, fmt.Errorf("统计章节数失败: %w", err)
	}

	return count, nil
}

// CountByStatus 按状态统计章节数
func (r *MongoChapterRepository) CountByStatus(ctx context.Context, bookID string, status int) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"book_id": bookID,
		"status":  status,
	})
	if err != nil {
		return 0, fmt.Errorf("统计章节数失败: %w", err)
	}

	return count, nil
}

// CountVIPChapters 统计VIP章节数
func (r *MongoChapterRepository) CountVIPChapters(ctx context.Context, bookID string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"book_id": bookID,
		"status":  1,
		"is_vip":  true,
	})
	if err != nil {
		return 0, fmt.Errorf("统计VIP章节数失败: %w", err)
	}

	return count, nil
}

// BatchCreate 批量创建章节
func (r *MongoChapterRepository) BatchCreate(ctx context.Context, chapters []*reader.Chapter) error {
	if len(chapters) == 0 {
		return nil
	}

	docs := make([]interface{}, len(chapters))
	now := time.Now()
	for i, chapter := range chapters {
		if chapter.ID == "" {
			chapter.ID = generateID()
		}
		chapter.CreatedAt = now
		chapter.UpdatedAt = now
		docs[i] = chapter
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("批量创建章节失败: %w", err)
	}

	return nil
}

// BatchUpdateStatus 批量更新章节状态
func (r *MongoChapterRepository) BatchUpdateStatus(ctx context.Context, chapterIDs []string, status int) error {
	if len(chapterIDs) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"$in": chapterIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("批量更新章节状态失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除章节
func (r *MongoChapterRepository) BatchDelete(ctx context.Context, chapterIDs []string) error {
	if len(chapterIDs) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"$in": chapterIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("批量删除章节失败: %w", err)
	}

	return nil
}

// CheckVIPAccess 检查章节是否需要VIP权限
func (r *MongoChapterRepository) CheckVIPAccess(ctx context.Context, chapterID string) (bool, error) {
	var chapter reader.Chapter
	err := r.collection.FindOne(ctx, bson.M{"_id": chapterID}).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, fmt.Errorf("章节不存在")
		}
		return false, fmt.Errorf("查询章节失败: %w", err)
	}

	return chapter.IsVIP, nil
}

// GetChapterPrice 获取章节价格
func (r *MongoChapterRepository) GetChapterPrice(ctx context.Context, chapterID string) (int64, error) {
	var chapter reader.Chapter
	err := r.collection.FindOne(ctx, bson.M{"_id": chapterID}).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, fmt.Errorf("章节不存在")
		}
		return 0, fmt.Errorf("查询章节失败: %w", err)
	}

	return chapter.Price, nil
}

// GetChapterContent 获取章节内容
func (r *MongoChapterRepository) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	var chapter reader.Chapter
	projection := bson.M{"content": 1}
	opts := options.FindOne().SetProjection(projection)

	err := r.collection.FindOne(ctx, bson.M{"_id": chapterID}, opts).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("章节不存在")
		}
		return "", fmt.Errorf("查询章节内容失败: %w", err)
	}

	return chapter.Content, nil
}

// UpdateChapterContent 更新章节内容
func (r *MongoChapterRepository) UpdateChapterContent(ctx context.Context, chapterID string, content string) error {
	update := bson.M{
		"$set": bson.M{
			"content":    content,
			"updated_at": time.Now(),
			"word_count": len([]rune(content)), // 计算字数
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": chapterID}, update)
	if err != nil {
		return fmt.Errorf("更新章节内容失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("章节不存在")
	}

	return nil
}

// Health 健康检查
func (r *MongoChapterRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// 全局ID计数器，用于避免并发时的ID冲突
var idCounter int64

// generateID 生成唯一ID (简化实现，实际应使用更可靠的ID生成方案)
func generateID() string {
	idCounter++
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), idCounter)
}
