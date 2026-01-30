package base

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/shared/types"
)

// BaseMongoRepository MongoDB仓储基类
//
// 提供统一的ID转换逻辑和通用CRUD方法，减少重复代码喵~
//
// 用法示例：
//
//	type MongoReadingProgressRepository struct {
//	    *mongodb.BaseMongoRepository  // 嵌入基类
//	}
//
//	func NewMongoReadingProgressRepository(db *mongo.Database) *MongoReadingProgressRepository {
//	    return &MongoReadingProgressRepository{
//	        BaseMongoRepository: NewBaseMongoRepository(db, "reading_progress"),
//	    }
//	}
//
//	// 使用基类的 ParseID 方法
//	func (r *MongoReadingProgressRepository) GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error) {
//	    objectID, err := r.ParseID(id)  // 使用基类的ID转换方法
//	    if err != nil {
//	        return nil, err
//	    }
//	    // ...
//	}
type BaseMongoRepository struct {
	collection *mongo.Collection
}

// NewBaseMongoRepository 创建基类实例
//
// 参数：
//   - db: MongoDB数据库实例
//   - collectionName: 集合名称
//
// 返回：
//   - *BaseMongoRepository: 基类实例
func NewBaseMongoRepository(db *mongo.Database, collectionName string) *BaseMongoRepository {
	return &BaseMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// GetCollection 获取集合实例（用于子类直接访问）
//
// 返回：
//   - *mongo.Collection: MongoDB集合实例
func (b *BaseMongoRepository) GetCollection() *mongo.Collection {
	return b.collection
}

// ===== ID转换辅助方法 =====
// 这些方法封装了 models/shared/types 中的ID转换逻辑，提供统一的错误处理喵~

// ParseID 解析ID字符串为ObjectID
//
// 使用 models/shared/types.ParseObjectID 进行转换，提供统一的错误处理喵~
//
// 参数：
//   - id: hex格式的ID字符串
//
// 返回：
//   - primitive.ObjectID: 解析后的ObjectID
//   - error: 转换错误（types.ErrInvalidIDFormat 或 types.ErrEmptyID）
//
// 示例：
//	objectID, err := r.ParseID("507f1f77bcf86cd799439011")
//	if err != nil {
//	    return nil, err
//	}
func (b *BaseMongoRepository) ParseID(id string) (primitive.ObjectID, error) {
	return types.ParseObjectID(id)
}

// ParseIDs 批量解析ID字符串
//
// 使用 models/shared/types.ParseObjectIDSlice 进行批量转换喵~
//
// 参数：
//   - ids: hex格式的ID字符串切片
//
// 返回：
//   - []primitive.ObjectID: 解析后的ObjectID切片
//   - error: 如果有任何ID转换失败，返回错误
//
// 示例：
//	oids, err := r.ParseIDs([]string{"id1", "id2", "id3"})
//	if err != nil {
//	    return nil, err
//	}
func (b *BaseMongoRepository) ParseIDs(ids []string) ([]primitive.ObjectID, error) {
	oids, errs := types.ParseObjectIDSlice(ids)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse %d IDs", len(errs))
	}
	return oids, nil
}

// IDToHex ObjectID转换为hex字符串
//
// 使用 models/shared/types.ToHex 进行转换喵~
//
// 参数：
//   - id: ObjectID
//
// 返回：
//   - string: hex格式的ID字符串
//
// 示例：
//	hexID := r.IDToHex(objectID)
func (b *BaseMongoRepository) IDToHex(id primitive.ObjectID) string {
	return types.ToHex(id)
}

// IDsToHex 批量转换ObjectID为hex字符串
//
// 使用 models/shared/types.ToHexSlice 进行批量转换喵~
//
// 参数：
//   - ids: ObjectID切片
//
// 返回：
//   - []string: hex格式的ID字符串切片
//
// 示例：
//	hexIDs := r.IDsToHex([]primitive.ObjectID{id1, id2})
func (b *BaseMongoRepository) IDsToHex(ids []primitive.ObjectID) []string {
	return types.ToHexSlice(ids)
}

// IsValidID 检查字符串是否为有效的ObjectID格式
//
// 使用 models/shared/types.IsValidObjectID 进行验证喵~
//
// 参数：
//   - id: 待验证的ID字符串
//
// 返回：
//   - bool: 是否为有效的ObjectID格式
//
// 示例：
//	if !r.IsValidID("invalid") {
//	    return nil, types.ErrInvalidIDFormat
//	}
func (b *BaseMongoRepository) IsValidID(id string) bool {
	return types.IsValidObjectID(id)
}

// GenerateID 生成新的ObjectID并返回hex字符串
//
// 使用 models/shared/types.GenerateNewObjectID 生成新ID喵~
//
// 返回：
//   - string: 新生成的hex格式ID字符串
//
// 示例：
//	newID := r.GenerateID()
func (b *BaseMongoRepository) GenerateID() string {
	return types.GenerateNewObjectID()
}

// ===== 通用CRUD方法 =====
// 提供常用的CRUD操作，减少重复代码喵~

// FindByID 根据ID查找文档
//
// 参数：
//   - ctx: 上下文
//   - id: 文档ID（hex字符串）
//   - result: 用于接收查询结果的对象指针
//
// 返回：
//   - error: 查询错误或mongo.ErrNoDocuments
//
// 示例：
//	var progress reader.ReadingProgress
//	err := r.FindByID(ctx, id, &progress)
//	if err != nil {
//	    return nil, err
//	}
func (b *BaseMongoRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	oid, err := b.ParseID(id)
	if err != nil {
		return err
	}

	return b.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(result)
}

// FindByIDWithOpts 根据ID查找文档（支持自定义选项）
//
// 参数：
//   - ctx: 上下文
//   - id: 文档ID（hex字符串）
//   - result: 用于接收查询结果的对象指针
//   - opts: 查询选项
//
// 返回：
//   - error: 查询错误或mongo.ErrNoDocuments
//
// 示例：
//	opts := options.FindOne().SetProjection(bson.M{"password": 0})
//	err := r.FindByIDWithOpts(ctx, id, &user, opts)
func (b *BaseMongoRepository) FindByIDWithOpts(ctx context.Context, id string, result interface{}, opts ...*options.FindOneOptions) error {
	oid, err := b.ParseID(id)
	if err != nil {
		return err
	}

	return b.collection.FindOne(ctx, bson.M{"_id": oid}, opts...).Decode(result)
}

// UpdateByID 根据ID更新文档
//
// 参数：
//   - ctx: 上下文
//   - id: 文档ID（hex字符串）
//   - update: 更新内容（bson.M格式，通常使用 bson.M{"$set": fields}）
//
// 返回：
//   - error: 更新错误
//
// 示例：
//	update := bson.M{"$set": bson.M{"status": "completed", "updated_at": time.Now()}}
//	err := r.UpdateByID(ctx, id, update)
func (b *BaseMongoRepository) UpdateByID(ctx context.Context, id string, update bson.M) error {
	oid, err := b.ParseID(id)
	if err != nil {
		return err
	}

	_, err = b.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

// DeleteByID 根据ID删除文档
//
// 参数：
//   - ctx: 上下文
//   - id: 文档ID（hex字符串）
//
// 返回：
//   - error: 删除错误
//
// 示例：
//	err := r.DeleteByID(ctx, id)
func (b *BaseMongoRepository) DeleteByID(ctx context.Context, id string) error {
	oid, err := b.ParseID(id)
	if err != nil {
		return err
	}

	_, err = b.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// ===== 通用查询方法 =====

// Find 查找多个文档
//
// 参数：
//   - ctx: 上下文
//   - filter: 查询条件
//   - results: 用于接收查询结果的切片指针
//   - opts: 可选的查询选项
//
// 返回：
//   - error: 查询错误
//
// 示例：
//	var progresses []*reader.ReadingProgress
//	filter := bson.M{"user_id": userOID}
//	err := r.Find(ctx, filter, &progresses)
func (b *BaseMongoRepository) Find(ctx context.Context, filter bson.M, results interface{}, opts ...*options.FindOptions) error {
	cursor, err := b.collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	return cursor.All(ctx, results)
}

// FindOne 查找单个文档
//
// 参数：
//   - ctx: 上下文
//   - filter: 查询条件
//   - result: 用于接收查询结果的对象指针
//   - opts: 可选的查询选项
//
// 返回：
//   - error: 查询错误或mongo.ErrNoDocuments
//
// 示例：
//	var progress reader.ReadingProgress
//	filter := bson.M{"user_id": userOID, "book_id": bookOID}
//	err := r.FindOne(ctx, filter, &progress)
func (b *BaseMongoRepository) FindOne(ctx context.Context, filter bson.M, result interface{}, opts ...*options.FindOneOptions) error {
	return b.collection.FindOne(ctx, filter, opts...).Decode(result)
}

// Count 统计文档数量
//
// 参数：
//   - ctx: 上下文
//   - filter: 查询条件
//
// 返回：
//   - int64: 文档数量
//   - error: 统计错误
//
// 示例：
//	count, err := r.Count(ctx, bson.M{"status": "active"})
func (b *BaseMongoRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return b.collection.CountDocuments(ctx, filter)
}

// Exists 检查指定ID的文档是否存在
//
// 参数：
//   - ctx: 上下文
//   - id: 文档ID（hex字符串）
//
// 返回：
//   - bool: 是否存在
//   - error: 查询错误
//
// 示例：
//	exists, err := r.Exists(ctx, id)
//	if err != nil {
//	    return false, err
//	}
//	if !exists {
//	    return fmt.Errorf("document not found")
//	}
func (b *BaseMongoRepository) Exists(ctx context.Context, id string) (bool, error) {
	oid, err := b.ParseID(id)
	if err != nil {
		return false, err
	}

	count, err := b.collection.CountDocuments(ctx, bson.M{"_id": oid})
	return count > 0, err
}

// ExistsByFilter 根据条件检查文档是否存在
//
// 参数：
//   - ctx: 上下文
//   - filter: 查询条件
//
// 返回：
//   - bool: 是否存在
//   - error: 查询错误
//
// 示例：
//	filter := bson.M{"user_id": userOID, "book_id": bookOID}
//	exists, err := r.ExistsByFilter(ctx, filter)
func (b *BaseMongoRepository) ExistsByFilter(ctx context.Context, filter bson.M) (bool, error) {
	count, err := b.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

// Create 创建文档
//
// 参数：
//   - ctx: 上下文
//   - document: 要创建的文档
//
// 返回：
//   - error: 创建错误
//
// 示例：
//	progress := &reader.ReadingProgress{...}
//	err := r.Create(ctx, progress)
func (b *BaseMongoRepository) Create(ctx context.Context, document interface{}) error {
	_, err := b.collection.InsertOne(ctx, document)
	return err
}

// CreateWithResult 创建文档并返回插入的ID
//
// 参数：
//   - ctx: 上下文
//   - document: 要创建的文档
//
// 返回：
//   - primitive.ObjectID: 插入的文档ID
//   - error: 创建错误
//
// 示例：
//	progress := &reader.ReadingProgress{...}
//	id, err := r.CreateWithResult(ctx, progress)
//	if err != nil {
//	    return err
//	}
//	progress.ID = id
func (b *BaseMongoRepository) CreateWithResult(ctx context.Context, document interface{}) (primitive.ObjectID, error) {
	result, err := b.collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
