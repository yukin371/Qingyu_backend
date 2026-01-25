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

// MongoBatchOperationRepository MongoDB批量操作仓储实现
type MongoBatchOperationRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoBatchOperationRepository 创建MongoDB批量操作仓储
func NewMongoBatchOperationRepository(db *mongo.Database) writingInterface.BatchOperationRepository {
	return &MongoBatchOperationRepository{
		db:         db,
		collection: db.Collection("batch_operations"),
	}
}

// Create 创建批量操作
func (r *MongoBatchOperationRepository) Create(ctx context.Context, op *writer.BatchOperation) error {
	if op == nil {
		return fmt.Errorf("批量操作对象不能为空")
	}

	// 生成ID
	if op.IdentifiedEntity.ID.IsZero() {
		op.IdentifiedEntity.ID = primitive.NewObjectID()
	}

	// 设置时间戳和默认值
	op.TouchForCreate()

	// 插入数据库
	_, err := r.collection.InsertOne(ctx, op)
	if err != nil {
		return fmt.Errorf("创建批量操作失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取批量操作
func (r *MongoBatchOperationRepository) GetByID(ctx context.Context, id string) (*writer.BatchOperation, error) {
	var op writer.BatchOperation

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的批量操作ID: %w", err)
	}

	filter := bson.M{
		"_id":        objID,
		"deleted_at": nil,
	}

	err = r.collection.FindOne(ctx, filter).Decode(&op)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询批量操作失败: %w", err)
	}

	return &op, nil
}

// Update 更新批量操作
func (r *MongoBatchOperationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的批量操作ID: %w", err)
	}

	// 自动更新updated_at
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": objID, "deleted_at": nil}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新批量操作失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("批量操作不存在或已删除")
	}

	return nil
}

// Delete 物理删除批量操作
func (r *MongoBatchOperationRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的批量操作ID: %w", err)
	}

	filter := bson.M{"_id": objID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除批量操作失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("批量操作不存在")
	}

	return nil
}

// List 查询批量操作列表
func (r *MongoBatchOperationRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.BatchOperation, error) {
	mongoFilter := bson.M{"deleted_at": nil}

	// 如果有筛选条件，合并条件
	if filter != nil {
		conditions := filter.GetConditions()
		for key, value := range conditions {
			mongoFilter[key] = value
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 如果Filter提供了排序，使用Filter的排序
	if filter != nil && filter.GetSort() != nil {
		opts.SetSort(filter.GetSort())
	}

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询批量操作列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var operations []*writer.BatchOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, fmt.Errorf("解析批量操作数据失败: %w", err)
	}

	return operations, nil
}

// Exists 检查批量操作是否存在
func (r *MongoBatchOperationRepository) Exists(ctx context.Context, id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("无效的批量操作ID: %w", err)
	}

	filter := bson.M{"_id": objID, "deleted_at": nil}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查批量操作存在失败: %w", err)
	}

	return count > 0, nil
}

// Count 统计批量操作总数
func (r *MongoBatchOperationRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
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
		return 0, fmt.Errorf("统计批量操作数失败: %w", err)
	}

	return count, nil
}

// GetByProjectID 获取项目的批量操作列表
func (r *MongoBatchOperationRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.BatchOperation, error) {
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, fmt.Errorf("无效的项目ID: %w", err)
	}

	filter := bson.M{
		"project_id": objID,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询批量操作列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var operations []*writer.BatchOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, fmt.Errorf("解析批量操作数据失败: %w", err)
	}

	return operations, nil
}

// GetByProjectAndType 按项目和操作类型查询
func (r *MongoBatchOperationRepository) GetByProjectAndType(ctx context.Context, projectID string, opType writer.BatchOperationType, limit, offset int64) ([]*writer.BatchOperation, error) {
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, fmt.Errorf("无效的项目ID: %w", err)
	}

	filter := bson.M{
		"project_id": objID,
		"type":       opType,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询批量操作列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var operations []*writer.BatchOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, fmt.Errorf("解析批量操作数据失败: %w", err)
	}

	return operations, nil
}

// GetByProjectAndStatus 按项目和状态查询
func (r *MongoBatchOperationRepository) GetByProjectAndStatus(ctx context.Context, projectID string, status writer.BatchOperationStatus, limit, offset int64) ([]*writer.BatchOperation, error) {
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, fmt.Errorf("无效的项目ID: %w", err)
	}

	filter := bson.M{
		"project_id": objID,
		"status":     status,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询批量操作列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var operations []*writer.BatchOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, fmt.Errorf("解析批量操作数据失败: %w", err)
	}

	return operations, nil
}

// GetByClientRequestID 根据客户端请求ID查询（幂等性支持）
func (r *MongoBatchOperationRepository) GetByClientRequestID(ctx context.Context, clientRequestID string) (*writer.BatchOperation, error) {
	var op writer.BatchOperation

	filter := bson.M{
		"client_request_id": clientRequestID,
		"deleted_at":        nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&op)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询批量操作失败: %w", err)
	}

	return &op, nil
}

// UpdateStatus 更新操作状态
func (r *MongoBatchOperationRepository) UpdateStatus(ctx context.Context, operationID string, status writer.BatchOperationStatus) error {
	objID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return fmt.Errorf("无效的批量操作ID: %w", err)
	}

	filter := bson.M{"_id": objID, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新操作状态失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("批量操作不存在或已删除")
	}

	return nil
}

// UpdateItemStatus 更新操作项状态
func (r *MongoBatchOperationRepository) UpdateItemStatus(ctx context.Context, operationID, targetID string, itemStatus writer.BatchItemStatus, errCode, errMsg string) error {
	objID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return fmt.Errorf("无效的批量操作ID: %w", err)
	}

	now := primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{
		"_id":         objID,
		"deleted_at":  nil,
		"items.target_id": targetID,
	}

	update := bson.M{
		"$set": bson.M{
			"items.$.status":         itemStatus,
			"items.$.error_code":     errCode,
			"items.$.error_msg":      errMsg,
			"items.$.error_message":  errMsg,
			"items.$.updated_at":     now,
		},
	}

	// 如果是成功或失败状态，设置completed_at
	if itemStatus == writer.BatchItemStatusSucceeded || itemStatus == writer.BatchItemStatusFailed {
		update["$set"].(bson.M)["items.$.completed_at"] = &now
	}

	// 如果是处理中状态，设置started_at（如果还未设置）
	if itemStatus == writer.BatchItemStatusProcessing {
		update["$set"].(bson.M)["items.$.started_at"] = &now
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新操作项状态失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("批量操作或操作项不存在")
	}

	return nil
}

// AddOperationLog 添加操作日志
func (r *MongoBatchOperationRepository) AddOperationLog(ctx context.Context, operationID string, log map[string]interface{}) error {
	// TODO: 如果有独立的操作日志表，实现此方法
	// 当前简化实现：直接更新payload
	return nil
}

// GetPendingOperations 获取待处理的操作列表
func (r *MongoBatchOperationRepository) GetPendingOperations(ctx context.Context, limit int64) ([]*writer.BatchOperation, error) {
	filter := bson.M{
		"status":     writer.BatchOpStatusPending,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询待处理操作失败: %w", err)
	}
	defer cursor.Close(ctx)

	var operations []*writer.BatchOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, fmt.Errorf("解析待处理操作数据失败: %w", err)
	}

	return operations, nil
}

// GetProcessingOperations 获取正在处理的操作列表
func (r *MongoBatchOperationRepository) GetProcessingOperations(ctx context.Context, limit int64) ([]*writer.BatchOperation, error) {
	filter := bson.M{
		"status":     writer.BatchOpStatusProcessing,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "started_at", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询处理中操作失败: %w", err)
	}
	defer cursor.Close(ctx)

	var operations []*writer.BatchOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, fmt.Errorf("解析处理中操作数据失败: %w", err)
	}

	return operations, nil
}

// CountByProjectAndStatus 按项目和状态统计
func (r *MongoBatchOperationRepository) CountByProjectAndStatus(ctx context.Context, projectID string, status writer.BatchOperationStatus) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return 0, fmt.Errorf("无效的项目ID: %w", err)
	}

	filter := bson.M{
		"project_id": objID,
		"status":     status,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计批量操作数失败: %w", err)
	}

	return count, nil
}

// DeleteByProject 按项目删除操作
func (r *MongoBatchOperationRepository) DeleteByProject(ctx context.Context, operationID, projectID string) error {
	opObjID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return fmt.Errorf("无效的批量操作ID: %w", err)
	}

	projectObjID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return fmt.Errorf("无效的项目ID: %w", err)
	}

	filter := bson.M{
		"_id":        opObjID,
		"project_id": projectObjID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除批量操作失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("批量操作不存在")
	}

	return nil
}

// SoftDeleteByProject 按项目软删除
func (r *MongoBatchOperationRepository) SoftDeleteByProject(ctx context.Context, operationID, projectID string) error {
	opObjID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return fmt.Errorf("无效的批量操作ID: %w", err)
	}

	projectObjID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return fmt.Errorf("无效的项目ID: %w", err)
	}

	now := time.Now()

	filter := bson.M{
		"_id":        opObjID,
		"project_id": projectObjID,
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
		return fmt.Errorf("软删除批量操作失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("批量操作不存在或无权限")
	}

	return nil
}

// IsProjectMember 检查操作是否属于项目
func (r *MongoBatchOperationRepository) IsProjectMember(ctx context.Context, operationID, projectID string) (bool, error) {
	opObjID, err := primitive.ObjectIDFromHex(operationID)
	if err != nil {
		return false, fmt.Errorf("无效的批量操作ID: %w", err)
	}

	projectObjID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return false, fmt.Errorf("无效的项目ID: %w", err)
	}

	filter := bson.M{
		"_id":        opObjID,
		"project_id": projectObjID,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查批量操作失败: %w", err)
	}

	return count > 0, nil
}

// Health 健康检查
func (r *MongoBatchOperationRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// EnsureIndexes 创建索引
func (r *MongoBatchOperationRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "project_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "type", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "client_request_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}
