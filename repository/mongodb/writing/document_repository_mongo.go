package writing

import (
	"Qingyu_backend/models/writer"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/repository/interfaces/infrastructure"
	writingInterface "Qingyu_backend/repository/interfaces/writing"
)

// MongoDocumentRepository MongoDB文档仓储实现
type MongoDocumentRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoDocumentRepository 创建MongoDB文档仓储
func NewMongoDocumentRepository(db *mongo.Database) writingInterface.DocumentRepository {
	return &MongoDocumentRepository{
		db:         db,
		collection: db.Collection("documents"),
	}
}

// Create 创建文档
func (r *MongoDocumentRepository) Create(ctx context.Context, doc *writer.Document) error {
	if doc == nil {
		return fmt.Errorf("文档对象不能为空")
	}

	// 生成ID
	if doc.IdentifiedEntity.ID == "" {
		doc.IdentifiedEntity.ID = primitive.NewObjectID().Hex()
	}

	// 设置时间戳和默认值
	doc.TouchForCreate()

	// 基础验证（不验证文档类型，因为需要项目的writing_type）
	if err := doc.ValidateWithoutType(); err != nil {
		return fmt.Errorf("文档数据验证失败: %w", err)
	}

	// 插入数据库
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("创建文档失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取文档
func (r *MongoDocumentRepository) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	var doc writer.Document

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{
		"_id":        objID,
		"deleted_at": nil,
	}

	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文档失败: %w", err)
	}

	return &doc, nil
}

// Update 更新文档
func (r *MongoDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	// 自动更新updated_at
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": objID, "deleted_at": nil}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新文档失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("文档不存在或已删除")
	}

	return nil
}

// Delete 物理删除文档
func (r *MongoDocumentRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{"_id": objID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("文档不存在")
	}

	return nil
}

// List 查询文档列表
func (r *MongoDocumentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Document, error) {
	mongoFilter := bson.M{"deleted_at": nil}

	// 如果有筛选条件，合并条件
	if filter != nil {
		conditions := filter.GetConditions()
		for key, value := range conditions {
			mongoFilter[key] = value
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})

	// 如果Filter提供了排序，使用Filter的排序
	if filter != nil && filter.GetSort() != nil {
		opts.SetSort(filter.GetSort())
	}

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询文档列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []*writer.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("解析文档数据失败: %w", err)
	}

	return documents, nil
}

// Exists 检查文档是否存在
func (r *MongoDocumentRepository) Exists(ctx context.Context, id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{"_id": objID, "deleted_at": nil}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查文档存在失败: %w", err)
	}

	return count > 0, nil
}

// Count 统计文档总数
func (r *MongoDocumentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	mongoFilter := bson.M{"deleted_at": nil}

	// 如果有筛选条件，合并条件
	if filter != nil {
		conditions := filter.GetConditions()
		for key, value := range conditions {
			mongoFilter[key] = value
		}
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, fmt.Errorf("统计文档数失败: %w", err)
	}

	return count, nil
}

// GetByProjectID 获取项目的所有文档
func (r *MongoDocumentRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Document, error) {
	filter := bson.M{
		"project_id": projectID,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{
			{Key: "level", Value: 1},
			{Key: "order", Value: 1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询文档列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []*writer.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("解析文档数据失败: %w", err)
	}

	return documents, nil
}

// GetByProjectAndType 按项目和类型查询文档
func (r *MongoDocumentRepository) GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*writer.Document, error) {
	filter := bson.M{
		"project_id": projectID,
		"type":       documentType,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "order", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询文档列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []*writer.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("解析文档数据失败: %w", err)
	}

	return documents, nil
}

// UpdateByProject 按项目更新文档
func (r *MongoDocumentRepository) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	filter := bson.M{
		"_id":        objID,
		"project_id": projectID,
		"deleted_at": nil,
	}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新文档失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("文档不存在或无权限")
	}

	return nil
}

// DeleteByProject 按项目删除文档
func (r *MongoDocumentRepository) DeleteByProject(ctx context.Context, documentID, projectID string) error {
	objID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{
		"_id":        objID,
		"project_id": projectID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("文档不存在")
	}

	return nil
}

// RestoreByProject 按项目恢复文档
func (r *MongoDocumentRepository) RestoreByProject(ctx context.Context, documentID, projectID string) error {
	objID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{
		"_id":        objID,
		"project_id": projectID,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": nil,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("恢复文档失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("文档不存在")
	}

	return nil
}

// IsProjectMember 检查文档是否属于项目
func (r *MongoDocumentRepository) IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return false, fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{
		"_id":        objID,
		"project_id": projectID,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查文档失败: %w", err)
	}

	return count > 0, nil
}

// SoftDelete 软删除文档
func (r *MongoDocumentRepository) SoftDelete(ctx context.Context, documentID, projectID string) error {
	objID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	now := time.Now()

	filter := bson.M{
		"_id":        objID,
		"project_id": projectID,
		"deleted_at": nil,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("软删除文档失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("文档不存在或无权限")
	}

	return nil
}

// HardDelete 物理删除文档
func (r *MongoDocumentRepository) HardDelete(ctx context.Context, documentID string) error {
	objID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return fmt.Errorf("无效的文档ID: %w", err)
	}

	filter := bson.M{"_id": objID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("物理删除文档失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("文档不存在")
	}

	return nil
}

// CountByProject 统计项目的文档数
func (r *MongoDocumentRepository) CountByProject(ctx context.Context, projectID string) (int64, error) {
	filter := bson.M{
		"project_id": projectID,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计文档数失败: %w", err)
	}

	return count, nil
}

// CreateWithTransaction 在事务中创建文档
func (r *MongoDocumentRepository) CreateWithTransaction(ctx context.Context, doc *writer.Document, callback func(ctx context.Context) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 创建文档
		if err := r.Create(sessCtx, doc); err != nil {
			return nil, err
		}

		// 执行回调
		if callback != nil {
			if err := callback(sessCtx); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoDocumentRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// EnsureIndexes 创建索引
func (r *MongoDocumentRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "project_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "parent_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "level", Value: 1},
				{Key: "order", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}
