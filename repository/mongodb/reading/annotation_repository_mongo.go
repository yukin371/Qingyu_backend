package reading

import (
	"Qingyu_backend/models/reader"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAnnotationRepository 标注仓储MongoDB实现
type MongoAnnotationRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

// NewMongoAnnotationRepository 创建标注仓储实例
func NewMongoAnnotationRepository(db *mongo.Database) *MongoAnnotationRepository {
	return &MongoAnnotationRepository{
		collection: db.Collection("annotations"),
		db:         db,
	}
}

// Create 创建标注
func (r *MongoAnnotationRepository) Create(ctx context.Context, annotation *reader.Annotation) error {
	if annotation.ID == "" {
		annotation.ID = generateAnnotationID()
	}
	annotation.CreatedAt = time.Now()
	annotation.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, annotation)
	if err != nil {
		return fmt.Errorf("创建标注失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取标注
func (r *MongoAnnotationRepository) GetByID(ctx context.Context, id string) (*reader.Annotation, error) {
	var annotation reader.Annotation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&annotation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("标注不存在")
		}
		return nil, fmt.Errorf("查询标注失败: %w", err)
	}

	return &annotation, nil
}

// Update 更新标注
func (r *MongoAnnotationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新标注失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("标注不存在")
	}

	return nil
}

// Delete 删除标注
func (r *MongoAnnotationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("删除标注失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("标注不存在")
	}

	return nil
}

// GetByUserAndBook 获取用户在某本书的所有标注
func (r *MongoAnnotationRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id": userID,
		"book_id": bookID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询标注列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析标注数据失败: %w", err)
	}

	return annotations, nil
}

// GetByUserAndChapter 获取用户在某章节的所有标注
func (r *MongoAnnotationRepository) GetByUserAndChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id":    userID,
		"book_id":    bookID,
		"chapter_id": chapterID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_offset", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询章节标注失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析标注数据失败: %w", err)
	}

	return annotations, nil
}

// GetByType 根据类型获取标注
func (r *MongoAnnotationRepository) GetByType(ctx context.Context, userID, bookID string, annotationType string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id": userID,
		"book_id": bookID,
		"type":    annotationType,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询标注失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析标注数据失败: %w", err)
	}

	return annotations, nil
}

// GetNotes 获取笔记
func (r *MongoAnnotationRepository) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return r.GetByType(ctx, userID, bookID, string(reader.AnnotationTypeNote))
}

// GetNotesByChapter 获取章节笔记
func (r *MongoAnnotationRepository) GetNotesByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id":    userID,
		"book_id":    bookID,
		"chapter_id": chapterID,
		"type":       string(reader.AnnotationTypeNote),
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_offset", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询章节笔记失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析笔记数据失败: %w", err)
	}

	return annotations, nil
}

// SearchNotes 搜索笔记
func (r *MongoAnnotationRepository) SearchNotes(ctx context.Context, userID string, keyword string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id": userID,
		"type":    string(reader.AnnotationTypeNote),
		"$or": []bson.M{
			{"content": bson.M{"$regex": keyword, "$options": "i"}},
			{"note": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("搜索笔记失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析笔记数据失败: %w", err)
	}

	return annotations, nil
}

// GetBookmarks 获取书签（type = 2）
func (r *MongoAnnotationRepository) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return r.GetByType(ctx, userID, bookID, string(reader.AnnotationTypeBookmark))
}

// GetBookmarkByPosition 根据位置获取书签
func (r *MongoAnnotationRepository) GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, startOffset int) (*reader.Annotation, error) {
	var annotation reader.Annotation
	filter := bson.M{
		"user_id":      userID,
		"book_id":      bookID,
		"chapter_id":   chapterID,
		"type":         string(reader.AnnotationTypeBookmark),
		"start_offset": startOffset,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&annotation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 书签不存在
		}
		return nil, fmt.Errorf("查询书签失败: %w", err)
	}

	return &annotation, nil
}

// GetLatestBookmark 获取最新的书签
func (r *MongoAnnotationRepository) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error) {
	var annotation reader.Annotation
	filter := bson.M{
		"user_id": userID,
		"book_id": bookID,
		"type":    string(reader.AnnotationTypeBookmark),
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&annotation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询最新书签失败: %w", err)
	}

	return &annotation, nil
}

// GetHighlights 获取高亮（type = 3）
func (r *MongoAnnotationRepository) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return r.GetByType(ctx, userID, bookID, string(reader.AnnotationTypeHighlight))
}

// GetHighlightsByChapter 获取章节高亮
func (r *MongoAnnotationRepository) GetHighlightsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id":    userID,
		"book_id":    bookID,
		"chapter_id": chapterID,
		"type":       string(reader.AnnotationTypeHighlight),
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_offset", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询章节高亮失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析高亮数据失败: %w", err)
	}

	return annotations, nil
}

// CountByUser 统计用户的标注数量
func (r *MongoAnnotationRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return 0, fmt.Errorf("统计标注数失败: %w", err)
	}

	return count, nil
}

// CountByBook 统计用户在某本书的标注数量
func (r *MongoAnnotationRepository) CountByBook(ctx context.Context, userID, bookID string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"book_id": bookID,
	})
	if err != nil {
		return 0, fmt.Errorf("统计书籍标注数失败: %w", err)
	}

	return count, nil
}

// CountByType 统计某类型的标注数量
func (r *MongoAnnotationRepository) CountByType(ctx context.Context, userID string, annotationType string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"type":    annotationType,
	})
	if err != nil {
		return 0, fmt.Errorf("统计标注数失败: %w", err)
	}

	return count, nil
}

// BatchCreate 批量创建标注
func (r *MongoAnnotationRepository) BatchCreate(ctx context.Context, annotations []*reader.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}

	docs := make([]interface{}, len(annotations))
	now := time.Now()
	for i, annotation := range annotations {
		if annotation.ID == "" {
			annotation.ID = generateAnnotationID()
		}
		annotation.CreatedAt = now
		annotation.UpdatedAt = now
		docs[i] = annotation
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("批量创建标注失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除标注
func (r *MongoAnnotationRepository) BatchDelete(ctx context.Context, annotationIDs []string) error {
	if len(annotationIDs) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"$in": annotationIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("批量删除标注失败: %w", err)
	}

	return nil
}

// DeleteByBook 删除某本书的所有标注
func (r *MongoAnnotationRepository) DeleteByBook(ctx context.Context, userID, bookID string) error {
	filter := bson.M{
		"user_id": userID,
		"book_id": bookID,
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除书籍标注失败: %w", err)
	}

	return nil
}

// DeleteByChapter 删除某章节的所有标注
func (r *MongoAnnotationRepository) DeleteByChapter(ctx context.Context, userID, bookID, chapterID string) error {
	filter := bson.M{
		"user_id":    userID,
		"book_id":    bookID,
		"chapter_id": chapterID,
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除章节标注失败: %w", err)
	}

	return nil
}

// SyncAnnotations 同步标注数据
func (r *MongoAnnotationRepository) SyncAnnotations(ctx context.Context, userID string, annotations []*reader.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, len(annotations))
	for i, annotation := range annotations {
		if annotation.ID == "" {
			annotation.ID = generateAnnotationID()
		}

		filter := bson.M{"_id": annotation.ID}
		update := bson.M{
			"$set": bson.M{
				"user_id":    annotation.UserID,
				"book_id":    annotation.BookID,
				"chapter_id": annotation.ChapterID,
				"type":       annotation.Type,
				"text":       annotation.Text,
				"note":       annotation.Note,
				"range":      annotation.Range,
				"updated_at": time.Now(),
			},
			"$setOnInsert": bson.M{
				"created_at": annotation.CreatedAt,
			},
		}

		models[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.collection.BulkWrite(ctx, models, opts)
	if err != nil {
		return fmt.Errorf("同步标注失败: %w", err)
	}

	return nil
}

// GetRecentAnnotations 获取最近的标注
func (r *MongoAnnotationRepository) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询最近标注失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析标注数据失败: %w", err)
	}

	return annotations, nil
}

// GetPublicAnnotations 获取公开的标注
func (r *MongoAnnotationRepository) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"book_id":    bookID,
		"chapter_id": chapterID,
		"is_public":  true,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询公开标注失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析标注数据失败: %w", err)
	}

	return annotations, nil
}

// GetSharedAnnotations 获取用户分享的标注
func (r *MongoAnnotationRepository) GetSharedAnnotations(ctx context.Context, userID string) ([]*reader.Annotation, error) {
	filter := bson.M{
		"user_id":   userID,
		"is_public": true,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询分享标注失败: %w", err)
	}
	defer cursor.Close(ctx)

	var annotations []*reader.Annotation
	if err = cursor.All(ctx, &annotations); err != nil {
		return nil, fmt.Errorf("解析标注数据失败: %w", err)
	}

	return annotations, nil
}

// Health 健康检查
func (r *MongoAnnotationRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// generateAnnotationID 生成标注ID
var annotationIDCounter int64

func generateAnnotationID() string {
	annotationIDCounter++
	return fmt.Sprintf("ann_%d_%d", time.Now().UnixNano(), annotationIDCounter)
}
