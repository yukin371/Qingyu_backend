package writer

import (
	"context"
	"errors"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrBatchOperationNotFound     = errors.New("batch operation not found")
	ErrBatchOperationNotCancellable = errors.New("batch operation is not cancellable")
)

// BatchOperationRepository 批量操作仓储接口
type BatchOperationRepository interface {
	// Create 创建批量操作
	Create(ctx context.Context, op *writer.BatchOperation) error

	// GetByID 根据ID获取
	GetByID(ctx context.Context, id primitive.ObjectID) (*writer.BatchOperation, error)

	// GetByClientRequestID 根据客户端请求ID获取（幂等性检查）
	GetByClientRequestID(ctx context.Context, projectID primitive.ObjectID, clientRequestID string) (*writer.BatchOperation, error)

	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status writer.BatchOperationStatus) error

	// Update 更新整个操作
	Update(ctx context.Context, op *writer.BatchOperation) error

	// ListByProject 查询项目的批量操作列表
	ListByProject(ctx context.Context, projectID primitive.ObjectID, opts *ListOptions) ([]*writer.BatchOperation, error)

	// GetRunningCount 获取运行中的操作数量
	GetRunningCount(ctx context.Context, projectID primitive.ObjectID) (int64, error)
}

// BatchOperationRepositoryImpl 批量操作仓储实现
type BatchOperationRepositoryImpl struct {
	*base.BaseMongoRepository // 嵌入基类，继承ID转换和通用CRUD方法喵~
	itemCollection            *mongo.Collection // 辅助collection独立管理喵~
}

// NewBatchOperationRepository 创建批量操作仓储
func NewBatchOperationRepository(db *mongo.Database) BatchOperationRepository {
	return &BatchOperationRepositoryImpl{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "batch_operations"),
		itemCollection:      db.Collection("batch_operation_items"),
	}
}

// Create 创建批量操作
func (r *BatchOperationRepositoryImpl) Create(ctx context.Context, op *writer.BatchOperation) error {
	op.TouchForCreate()
	_, err := r.GetCollection().InsertOne(ctx, op)
	return err
}

// GetByID 根据ID获取批量操作
func (r *BatchOperationRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*writer.BatchOperation, error) {
	var op writer.BatchOperation
	err := r.GetCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&op)
	if err == mongo.ErrNoDocuments {
		return nil, ErrBatchOperationNotFound
	}
	return &op, err
}

// GetByClientRequestID 根据客户端请求ID获取（幂等性检查）
func (r *BatchOperationRepositoryImpl) GetByClientRequestID(ctx context.Context, projectID primitive.ObjectID, clientRequestID string) (*writer.BatchOperation, error) {
	var op writer.BatchOperation
	err := r.GetCollection().FindOne(ctx, bson.M{
		"project_id":        projectID,
		"client_request_id": clientRequestID,
	}).Decode(&op)
	if err == mongo.ErrNoDocuments {
		return nil, ErrBatchOperationNotFound
	}
	return &op, err
}

// UpdateStatus 更新状态
func (r *BatchOperationRepositoryImpl) UpdateStatus(ctx context.Context, id primitive.ObjectID, status writer.BatchOperationStatus) error {
	update := bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": primitive.NewDateTimeFromTime(time.Now()),
	}}
	result, err := r.GetCollection().UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrBatchOperationNotFound
	}
	return nil
}

// Update 更新整个操作
func (r *BatchOperationRepositoryImpl) Update(ctx context.Context, op *writer.BatchOperation) error {
	op.Timestamps.Touch()
	result, err := r.GetCollection().UpdateByID(ctx, op.ID, bson.M{"$set": op})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrBatchOperationNotFound
	}
	return nil
}

// ListOptions 列表查询选项
type ListOptions struct {
	Limit  int64
	Skip   int64
	Status writer.BatchOperationStatus
	SortBy string // 默认created_at
}

// ListByProject 查询项目的批量操作列表
func (r *BatchOperationRepositoryImpl) ListByProject(ctx context.Context, projectID primitive.ObjectID, opts *ListOptions) ([]*writer.BatchOperation, error) {
	filter := bson.M{"project_id": projectID}
	if opts != nil && opts.Status != "" {
		filter["status"] = opts.Status
	}

	findOpts := options.Find()
	if opts != nil {
		if opts.Limit > 0 {
			findOpts.SetLimit(opts.Limit)
		}
		if opts.Skip > 0 {
			findOpts.SetSkip(opts.Skip)
		}
		if opts.SortBy != "" {
			findOpts.SetSort(bson.M{opts.SortBy: -1})
		} else {
			findOpts.SetSort(bson.M{"created_at": -1})
		}
	}

	cursor, err := r.GetCollection().Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	var ops []*writer.BatchOperation
	err = cursor.All(ctx, &ops)
	return ops, err
}

// GetRunningCount 获取运行中的操作数量
func (r *BatchOperationRepositoryImpl) GetRunningCount(ctx context.Context, projectID primitive.ObjectID) (int64, error) {
	count, err := r.GetCollection().CountDocuments(ctx, bson.M{
		"project_id": projectID,
		"status":     bson.M{"$in": []writer.BatchOperationStatus{writer.BatchOpStatusRunning, writer.BatchOpStatusProcessing}},
	})
	return count, err
}
