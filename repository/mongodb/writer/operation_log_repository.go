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
	ErrOperationLogNotFound  = errors.New("operation log not found")
	ErrOperationNotUndoable = errors.New("operation is not undoable")
	ErrOperationNotRedoable = errors.New("operation is not redoable")
)

// OperationLogRepository 操作日志仓储接口
type OperationLogRepository interface {
	// Create 创建操作日志
	Create(ctx context.Context, log *writer.OperationLog) error

	// GetByID 根据ID获取
	GetByID(ctx context.Context, id primitive.ObjectID) (*writer.OperationLog, error)

	// GetByChainID 根据链ID获取操作列表（批量操作的所有日志）
	GetByChainID(ctx context.Context, chainID string) ([]*writer.OperationLog, error)

	// GetLatestByProject 获取项目的最新操作日志（用于撤销栈）
	GetLatestByProject(ctx context.Context, projectID primitive.ObjectID, limit int) ([]*writer.OperationLog, error)

	// UpdateStatus 更新状态（executed/undone/redone）
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status writer.OperationLogStatus) error

	// MarkAsCommitted 标记为已提交
	MarkAsCommitted(ctx context.Context, id primitive.ObjectID) error
}

// OperationLogRepositoryImpl 操作日志仓储实现
type OperationLogRepositoryImpl struct {
	*base.BaseMongoRepository
}

// NewOperationLogRepository 创建操作日志仓储
func NewOperationLogRepository(db *mongo.Database) OperationLogRepository {
	return &OperationLogRepositoryImpl{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "operation_logs"),
	}
}

// Create 创建操作日志
func (r *OperationLogRepositoryImpl) Create(ctx context.Context, log *writer.OperationLog) error {
	log.TouchForCreate()
	_, err := r.GetCollection().InsertOne(ctx, log)
	return err
}

// GetByID 根据ID获取操作日志
func (r *OperationLogRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*writer.OperationLog, error) {
	var log writer.OperationLog
	err := r.GetCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&log)
	if err == mongo.ErrNoDocuments {
		return nil, ErrOperationLogNotFound
	}
	return &log, err
}

// GetByChainID 根据链ID获取操作列表（批量操作的所有日志）
func (r *OperationLogRepositoryImpl) GetByChainID(ctx context.Context, chainID string) ([]*writer.OperationLog, error) {
	filter := bson.M{"chain_id": chainID}
	opts := options.Find().SetSort(bson.M{"created_at": 1})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var logs []*writer.OperationLog
	err = cursor.All(ctx, &logs)
	return logs, err
}

// GetLatestByProject 获取项目的最新操作日志（用于撤销栈）
func (r *OperationLogRepositoryImpl) GetLatestByProject(ctx context.Context, projectID primitive.ObjectID, limit int) ([]*writer.OperationLog, error) {
	filter := bson.M{
		"project_id":   projectID,
		"is_committed": true,
	}
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(limit))

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var logs []*writer.OperationLog
	err = cursor.All(ctx, &logs)
	return logs, err
}

// UpdateStatus 更新状态（executed/undone/redone）
func (r *OperationLogRepositoryImpl) UpdateStatus(ctx context.Context, id primitive.ObjectID, status writer.OperationLogStatus) error {
	update := bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": primitive.NewDateTimeFromTime(time.Now()),
	}}

	// 根据状态添加时间戳
	switch status {
	case writer.OpLogStatusUndone:
		update["$set"].(bson.M)["undone_at"] = primitive.NewDateTimeFromTime(time.Now())
	case writer.OpLogStatusRedone:
		update["$set"].(bson.M)["redone_at"] = primitive.NewDateTimeFromTime(time.Now())
	}

	result, err := r.GetCollection().UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrOperationLogNotFound
	}
	return nil
}

// MarkAsCommitted 标记为已提交
func (r *OperationLogRepositoryImpl) MarkAsCommitted(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{"$set": bson.M{
		"is_committed": true,
		"updated_at":   primitive.NewDateTimeFromTime(time.Now()),
	}}
	result, err := r.GetCollection().UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrOperationLogNotFound
	}
	return nil
}
