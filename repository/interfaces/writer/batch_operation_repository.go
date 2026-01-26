package writer

import (
	"Qingyu_backend/models/writer"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
)

// BatchOperationRepository 批量操作仓储接口
type BatchOperationRepository interface {
	// 继承CRUDRepository接口
	base.CRUDRepository[*writer.BatchOperation, string]

	// 继承 HealthRepository 接口
	base.HealthRepository

	// BatchOperation特定的查询方法

	// GetByProjectID 获取项目的批量操作列表
	GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.BatchOperation, error)

	// GetByProjectAndType 按项目和操作类型查询
	GetByProjectAndType(ctx context.Context, projectID string, opType writer.BatchOperationType, limit, offset int64) ([]*writer.BatchOperation, error)

	// GetByProjectAndStatus 按项目和状态查询
	GetByProjectAndStatus(ctx context.Context, projectID string, status writer.BatchOperationStatus, limit, offset int64) ([]*writer.BatchOperation, error)

	// GetByClientRequestID 根据客户端请求ID查询（幂等性支持）
	GetByClientRequestID(ctx context.Context, clientRequestID string) (*writer.BatchOperation, error)

	// UpdateStatus 更新操作状态
	UpdateStatus(ctx context.Context, operationID string, status writer.BatchOperationStatus) error

	// UpdateItemStatus 更新操作项状态
	UpdateItemStatus(ctx context.Context, operationID, targetID string, itemStatus writer.BatchItemStatus, errCode, errMsg string) error

	// AddOperationLog 添加操作日志（如果有独立的日志表）
	AddOperationLog(ctx context.Context, operationID string, log map[string]interface{}) error

	// GetPendingOperations 获取待处理的操作列表
	GetPendingOperations(ctx context.Context, limit int64) ([]*writer.BatchOperation, error)

	// GetProcessingOperations 获取正在处理的操作列表
	GetProcessingOperations(ctx context.Context, limit int64) ([]*writer.BatchOperation, error)

	// CountByProjectAndStatus 按项目和状态统计
	CountByProjectAndStatus(ctx context.Context, projectID string, status writer.BatchOperationStatus) (int64, error)

	// DeleteByProject 按项目删除操作
	DeleteByProject(ctx context.Context, operationID, projectID string) error

	// SoftDeleteByProject 按项目软删除
	SoftDeleteByProject(ctx context.Context, operationID, projectID string) error

	// IsProjectMember 检查操作是否属于项目
	IsProjectMember(ctx context.Context, operationID, projectID string) (bool, error)
}
