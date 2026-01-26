package writer

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
	writingInterface "Qingyu_backend/repository/interfaces/writer"
)

// MongoDocumentContentRepository MongoDB文档内容仓储实现
type MongoDocumentContentRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoDocumentContentRepository 创建MongoDB文档内容仓储
func NewMongoDocumentContentRepository(db *mongo.Database) writingInterface.DocumentContentRepository {
	return &MongoDocumentContentRepository{
		db:         db,
		collection: db.Collection("document_contents"),
	}
}

// Create 创建文档内容
func (r *MongoDocumentContentRepository) Create(ctx context.Context, content *writer.DocumentContent) error {
	if content == nil {
		return fmt.Errorf("文档内容对象不能为空")
	}

	if content.ID.IsZero() {
		content.ID = primitive.NewObjectID()
	}

	content.TouchForCreate()

	if err := content.Validate(); err != nil {
		return fmt.Errorf("文档内容数据验证失败: %w", err)
	}

	_, err := r.collection.InsertOne(ctx, content)
	if err != nil {
		return fmt.Errorf("创建文档内容失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取文档内容
func (r *MongoDocumentContentRepository) GetByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID: %w", err)
	}

	var content writer.DocumentContent

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文档内容失败: %w", err)
	}

	return &content, nil
}

// GetByDocumentID 根据DocumentID获取文档内容
func (r *MongoDocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (*writer.DocumentContent, error) {
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, fmt.Errorf("无效的DocumentID: %w", err)
	}

	var content writer.DocumentContent
	err = r.collection.FindOne(ctx, bson.M{"document_id": objectID}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文档内容失败: %w", err)
	}

	return &content, nil
}

// Update 更新文档内容
func (r *MongoDocumentContentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新文档内容失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("文档内容不存在")
	}

	return nil
}

// UpdateWithVersion 带版本号的更新（乐观锁）
func (r *MongoDocumentContentRepository) UpdateWithVersion(ctx context.Context, documentID string, content string, expectedVersion int) error {
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return fmt.Errorf("无效的DocumentID: %w", err)
	}

	filter := bson.M{
		"document_id": objectID,
		"version":     expectedVersion,
	}

	// 计算字数统计
	wordCount := len([]rune(content))
	charCount := len(content)

	update := bson.M{
		"$set": bson.M{
			"content":    content,
			"word_count": wordCount,
			"char_count": charCount,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新文档内容失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("版本冲突，请重新获取最新内容")
	}

	return nil
}

// Delete 删除文档内容
func (r *MongoDocumentContentRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("删除文档内容失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("文档内容不存在")
	}

	return nil
}

// List 列出文档内容
func (r *MongoDocumentContentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.DocumentContent, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("查询文档内容列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var contents []*writer.DocumentContent
	if err = cursor.All(ctx, &contents); err != nil {
		return nil, fmt.Errorf("解析文档内容列表失败: %w", err)
	}

	return contents, nil
}

// Count 统计文档内容数量
func (r *MongoDocumentContentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("统计文档内容数量失败: %w", err)
	}
	return count, nil
}

// Exists 检查文档内容是否存在
func (r *MongoDocumentContentRepository) Exists(ctx context.Context, id string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, fmt.Errorf("检查文档内容存在性失败: %w", err)
	}

	return count > 0, nil
}

// BatchUpdateContent 批量更新内容（简化实现）
func (r *MongoDocumentContentRepository) BatchUpdateContent(ctx context.Context, updates map[string]string) error {
	// TODO: 实现批量更新
	return fmt.Errorf("批量更新暂未实现")
}

// GetContentStats 获取内容统计
func (r *MongoDocumentContentRepository) GetContentStats(ctx context.Context, documentID string) (wordCount, charCount int, err error) {
	var content writer.DocumentContent
	err = r.collection.FindOne(ctx, bson.M{"document_id": documentID}, options.FindOne().SetProjection(bson.M{"word_count": 1, "char_count": 1})).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("查询内容统计失败: %w", err)
	}

	return content.WordCount, content.CharCount, nil
}

// StoreToGridFS 存储大文档到GridFS（简化实现）
func (r *MongoDocumentContentRepository) StoreToGridFS(ctx context.Context, documentID string, content []byte) (string, error) {
	// TODO: 实现GridFS存储
	return "", fmt.Errorf("GridFS存储暂未实现")
}

// LoadFromGridFS 从GridFS加载大文档（简化实现）
func (r *MongoDocumentContentRepository) LoadFromGridFS(ctx context.Context, gridFSID string) ([]byte, error) {
	// TODO: 实现GridFS加载
	return nil, fmt.Errorf("GridFS加载暂未实现")
}

// CreateWithTransaction 在事务中创建文档内容
func (r *MongoDocumentContentRepository) CreateWithTransaction(ctx context.Context, content *writer.DocumentContent, callback func(ctx context.Context) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err := r.Create(sessCtx, content); err != nil {
			return nil, err
		}

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
func (r *MongoDocumentContentRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
