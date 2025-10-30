package audit

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/audit"
	"Qingyu_backend/repository/interfaces/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SensitiveWordRepositoryMongo 敏感词Repository的MongoDB实现
type SensitiveWordRepositoryMongo struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewSensitiveWordRepository 创建敏感词Repository
func NewSensitiveWordRepository(db *mongo.Database) *SensitiveWordRepositoryMongo {
	return &SensitiveWordRepositoryMongo{
		db:         db,
		collection: db.Collection("sensitive_words"),
	}
}

// ============ 基础CRUD ============

// Create 创建敏感词
func (r *SensitiveWordRepositoryMongo) Create(ctx context.Context, word *audit.SensitiveWord) error {
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, word)
	if err != nil {
		return fmt.Errorf("插入敏感词失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		word.ID = oid.Hex()
	}

	return nil
}

// GetByID 根据ID获取敏感词
func (r *SensitiveWordRepositoryMongo) GetByID(ctx context.Context, id string) (*audit.SensitiveWord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的敏感词ID: %w", err)
	}

	var word audit.SensitiveWord
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&word)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("敏感词不存在: %s", id)
		}
		return nil, fmt.Errorf("查询敏感词失败: %w", err)
	}

	return &word, nil
}

// Update 更新敏感词
func (r *SensitiveWordRepositoryMongo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的敏感词ID: %w", err)
	}

	updates["updatedAt"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新敏感词失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("敏感词不存在: %s", id)
	}

	return nil
}

// Delete 删除敏感词
func (r *SensitiveWordRepositoryMongo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的敏感词ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除敏感词失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("敏感词不存在: %s", id)
	}

	return nil
}

// ============ 查询方法 ============

// GetByWord 根据敏感词内容获取
func (r *SensitiveWordRepositoryMongo) GetByWord(ctx context.Context, word string) (*audit.SensitiveWord, error) {
	var sensitiveWord audit.SensitiveWord
	err := r.collection.FindOne(ctx, bson.M{"word": word}).Decode(&sensitiveWord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("敏感词不存在: %s", word)
		}
		return nil, fmt.Errorf("查询敏感词失败: %w", err)
	}

	return &sensitiveWord, nil
}

// List 列出敏感词
func (r *SensitiveWordRepositoryMongo) List(ctx context.Context, filter infrastructure.Filter) ([]*audit.SensitiveWord, error) {
	mongoFilter := r.buildFilter(filter)

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询敏感词列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var words []*audit.SensitiveWord
	if err = cursor.All(ctx, &words); err != nil {
		return nil, fmt.Errorf("解析敏感词列表失败: %w", err)
	}

	return words, nil
}

// Count 统计敏感词数量
func (r *SensitiveWordRepositoryMongo) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	mongoFilter := r.buildFilter(filter)

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, fmt.Errorf("统计敏感词数量失败: %w", err)
	}

	return count, nil
}

// FindWithPagination 分页查询敏感词
func (r *SensitiveWordRepositoryMongo) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[audit.SensitiveWord], error) {
	mongoFilter := r.buildFilter(filter)

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("统计敏感词总数失败: %w", err)
	}

	// 计算分页参数
	pagination.CalculatePagination()

	// 查询数据
	opts := options.Find().
		SetSkip(int64(pagination.Skip)).
		SetLimit(int64(pagination.PageSize)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询敏感词失败: %w", err)
	}
	defer cursor.Close(ctx)

	var words []*audit.SensitiveWord
	if err = cursor.All(ctx, &words); err != nil {
		return nil, fmt.Errorf("解析敏感词列表失败: %w", err)
	}

	return infrastructure.NewPagedResult(words, total, pagination), nil
}

// ============ 业务方法 ============

// GetEnabledWords 获取所有启用的敏感词
func (r *SensitiveWordRepositoryMongo) GetEnabledWords(ctx context.Context) ([]*audit.SensitiveWord, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"isEnabled": true})
	if err != nil {
		return nil, fmt.Errorf("查询启用的敏感词失败: %w", err)
	}
	defer cursor.Close(ctx)

	var words []*audit.SensitiveWord
	if err = cursor.All(ctx, &words); err != nil {
		return nil, fmt.Errorf("解析敏感词列表失败: %w", err)
	}

	return words, nil
}

// GetByCategory 根据分类获取敏感词
func (r *SensitiveWordRepositoryMongo) GetByCategory(ctx context.Context, category string) ([]*audit.SensitiveWord, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, fmt.Errorf("查询分类敏感词失败: %w", err)
	}
	defer cursor.Close(ctx)

	var words []*audit.SensitiveWord
	if err = cursor.All(ctx, &words); err != nil {
		return nil, fmt.Errorf("解析敏感词列表失败: %w", err)
	}

	return words, nil
}

// GetByLevel 根据等级获取敏感词
func (r *SensitiveWordRepositoryMongo) GetByLevel(ctx context.Context, minLevel int) ([]*audit.SensitiveWord, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"level": bson.M{"$gte": minLevel}})
	if err != nil {
		return nil, fmt.Errorf("查询等级敏感词失败: %w", err)
	}
	defer cursor.Close(ctx)

	var words []*audit.SensitiveWord
	if err = cursor.All(ctx, &words); err != nil {
		return nil, fmt.Errorf("解析敏感词列表失败: %w", err)
	}

	return words, nil
}

// BatchCreate 批量创建敏感词
func (r *SensitiveWordRepositoryMongo) BatchCreate(ctx context.Context, words []*audit.SensitiveWord) error {
	if len(words) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(words))
	for i, word := range words {
		word.CreatedAt = now
		word.UpdatedAt = now
		docs[i] = word
	}

	result, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("批量插入敏感词失败: %w", err)
	}

	// 设置ID
	for i, id := range result.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok {
			words[i].ID = oid.Hex()
		}
	}

	return nil
}

// BatchUpdate 批量更新敏感词
func (r *SensitiveWordRepositoryMongo) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, oid)
		}
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的敏感词ID")
	}

	updates["updatedAt"] = time.Now()

	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": objectIDs}},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("批量更新敏感词失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除敏感词
func (r *SensitiveWordRepositoryMongo) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, oid)
		}
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的敏感词ID")
	}

	_, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return fmt.Errorf("批量删除敏感词失败: %w", err)
	}

	return nil
}

// ============ 统计方法 ============

// CountByCategory 按分类统计敏感词数量
func (r *SensitiveWordRepositoryMongo) CountByCategory(ctx context.Context) (map[string]int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("统计分类敏感词失败: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]int64)
	for cursor.Next(ctx) {
		var doc struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		result[doc.ID] = doc.Count
	}

	return result, nil
}

// CountByLevel 按等级统计敏感词数量
func (r *SensitiveWordRepositoryMongo) CountByLevel(ctx context.Context) (map[int]int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$level"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("统计等级敏感词失败: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[int]int64)
	for cursor.Next(ctx) {
		var doc struct {
			ID    int   `bson:"_id"`
			Count int64 `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		result[doc.ID] = doc.Count
	}

	return result, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *SensitiveWordRepositoryMongo) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// ============ 辅助方法 ============

// buildFilter 构建MongoDB过滤器
func (r *SensitiveWordRepositoryMongo) buildFilter(filter infrastructure.Filter) bson.M {
	mongoFilter := bson.M{}

	if filter == nil {
		return mongoFilter
	}

	// 这里可以根据 filter 的具体类型添加更多过滤条件
	// 目前返回空过滤器（查询所有）

	return mongoFilter
}
