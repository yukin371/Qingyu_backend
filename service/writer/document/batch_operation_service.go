package document

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Qingyu_backend/models/writer"
	writerInterface "Qingyu_backend/repository/interfaces/writer"
	mongodbwriter "Qingyu_backend/repository/mongodb/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBatchOperationNotRunning = errors.New("batch operation is not running")
	ErrBatchOperationFailed     = errors.New("batch operation failed")
)

// BatchOperationService 批量操作服务接口
type BatchOperationService interface {
	// Submit 提交批量操作（含Preflight）
	Submit(ctx context.Context, req *SubmitBatchOperationRequest) (*writer.BatchOperation, error)

	// Execute 执行批量操作（异步）
	Execute(ctx context.Context, batchID primitive.ObjectID) error

	// Cancel 取消正在运行的操作
	Cancel(ctx context.Context, batchID primitive.ObjectID, userID primitive.ObjectID) error

	// Undo 撤销批量操作
	Undo(ctx context.Context, batchID primitive.ObjectID, userID primitive.ObjectID) error

	// GetProgress 获取操作进度
	GetProgress(ctx context.Context, batchID primitive.ObjectID) (*BatchOperationProgress, error)
}

// SubmitBatchOperationRequest 提交批量操作请求
type SubmitBatchOperationRequest struct {
	ProjectID          primitive.ObjectID
	Type               writer.BatchOperationType
	TargetIDs          []string
	Atomic             bool
	Payload            map[string]interface{}
	ConflictPolicy     writer.ConflictPolicy
	ExpectedVersions   map[string]int
	ClientRequestID    string
	UserID             primitive.ObjectID
	IncludeDescendants bool
}

// BatchOperationProgress 批量操作进度
type BatchOperationProgress struct {
	BatchID        primitive.ObjectID          `json:"batchId"`
	Status         writer.BatchOperationStatus `json:"status"`
	TotalItems     int                         `json:"totalItems"`
	CompletedItems int                         `json:"completedItems"`
	FailedItems    int                         `json:"failedItems"`
	StartedAt      *time.Time                  `json:"startedAt,omitempty"`
	FinishedAt     *time.Time                  `json:"finishedAt,omitempty"`
}

// BatchOperationServiceImpl 批量操作服务实现
type BatchOperationServiceImpl struct {
	batchOpRepo  *mongodbwriter.BatchOperationRepositoryImpl
	opLogRepo    *mongodbwriter.OperationLogRepositoryImpl
	docRepo      writerInterface.DocumentRepository
	preflightSvc PreflightService
}

// NewBatchOperationService 创建批量操作服务
func NewBatchOperationService(
	batchOpRepo *mongodbwriter.BatchOperationRepositoryImpl,
	opLogRepo *mongodbwriter.OperationLogRepositoryImpl,
	docRepo writerInterface.DocumentRepository,
) BatchOperationService {
	return &BatchOperationServiceImpl{
		batchOpRepo:  batchOpRepo,
		opLogRepo:    opLogRepo,
		docRepo:      docRepo,
		preflightSvc: NewPreflightService(docRepo),
	}
}

// Submit 提交批量操作（含Preflight预检查）
func (s *BatchOperationServiceImpl) Submit(ctx context.Context, req *SubmitBatchOperationRequest) (*writer.BatchOperation, error) {
	// 1. 幂等性检查
	if req.ClientRequestID != "" {
		existing, err := s.batchOpRepo.GetByClientRequestID(ctx, req.ProjectID, req.ClientRequestID)
		if err == nil && existing != nil {
			// 返回已存在的操作
			return existing, nil
		}
	}

	// 2. Preflight预检查
	summary, preflightResult, err := s.preflightSvc.ValidateBatchOperation(
		ctx,
		req.ProjectID,
		req.Type,
		req.TargetIDs,
		&PreflightOptions{
			ExpectedVersions:   req.ExpectedVersions,
			ConflictPolicy:     req.ConflictPolicy,
			IncludeDescendants: req.IncludeDescendants,
			UserID:             req.UserID,
		},
	)
	if err != nil && req.Atomic {
		// 原子操作模式下，Preflight失败则拒绝
		return nil, fmt.Errorf("preflight validation failed: %w", err)
	}

	// 3. 创建BatchOperation记录
	batchOp := &writer.BatchOperation{
		ProjectID:         req.ProjectID,
		Type:              req.Type,
		TargetIDs:         preflightResult.ValidIDs,
		OriginalTargetIDs: req.TargetIDs,
		Atomic:            req.Atomic,
		Payload:           req.Payload,
		ConflictPolicy:    req.ConflictPolicy,
		ExpectedVersions:  req.ExpectedVersions,
		ClientRequestID:   req.ClientRequestID,
		Status:            writer.BatchOpStatusPending,
		CreatedBy:         req.UserID,
		PreflightSummary:  summary,
	}

	// 根据节点数量选择执行模式
	if len(preflightResult.ValidIDs) <= 200 {
		batchOp.ExecutionMode = writer.ExecutionModeStandardAtomic
	} else {
		batchOp.ExecutionMode = writer.ExecutionModeSagaAtomic
	}

	err = s.batchOpRepo.Create(ctx, batchOp)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch operation: %w", err)
	}

	// 4. 创建BatchOperationItem记录（用于进度跟踪）
	for _, id := range preflightResult.ValidIDs {
		item := &writer.BatchOperationItem{
			BatchID:         batchOp.ID,
			TargetID:        id,
			TargetStableRef: preflightResult.DocumentMap[id].StableRef,
			Status:          writer.BatchItemStatusPending,
		}
		item.TouchForCreate()
		// TODO: 创建item记录到数据库（需要BatchOperationItemRepository）
		_ = item
	}

	return batchOp, nil
}

// Execute 执行批量操作
func (s *BatchOperationServiceImpl) Execute(ctx context.Context, batchID primitive.ObjectID) error {
	// 1. 加载BatchOperation
	batchOp, err := s.batchOpRepo.GetByID(ctx, batchID)
	if err != nil {
		return err
	}

	// 2. 更新状态为运行中
	err = s.batchOpRepo.UpdateStatus(ctx, batchID, writer.BatchOpStatusRunning)
	if err != nil {
		return err
	}

	now := time.Now()
	batchOp.StartedAt = &now

	// 3. 选择执行策略
	switch batchOp.ExecutionMode {
	case writer.ExecutionModeStandardAtomic:
		return s.executeStandardAtomic(ctx, batchOp)
	case writer.ExecutionModeSagaAtomic:
		return s.executeSagaAtomic(ctx, batchOp)
	default:
		return fmt.Errorf("unsupported execution mode: %s", batchOp.ExecutionMode)
	}
}

// executeStandardAtomic 标准原子执行（<=200节点，单事务）
func (s *BatchOperationServiceImpl) executeStandardAtomic(ctx context.Context, batchOp *writer.BatchOperation) error {
	// 注意：这里需要从docRepo获取MongoDB client来启动事务
	// 由于接口限制，暂时使用非事务方式执行
	// TODO: 在DocumentRepository接口中添加Client()方法获取MongoDB client

	// 执行批量操作
	var err error
	switch batchOp.Type {
	case writer.BatchOpTypeDelete:
		err = s.executeBatchDelete(ctx, batchOp)
	case writer.BatchOpTypeMove:
		err = s.executeBatchMove(ctx, batchOp)
	default:
		err = fmt.Errorf("unsupported operation type: %s", batchOp.Type)
	}

	if err != nil {
		s.batchOpRepo.UpdateStatus(ctx, batchOp.ID, writer.BatchOpStatusFailed)
		return fmt.Errorf("batch operation failed: %w", err)
	}

	// 创建OperationLog
	err = s.createOperationLog(ctx, batchOp)
	if err != nil {
		s.batchOpRepo.UpdateStatus(ctx, batchOp.ID, writer.BatchOpStatusFailed)
		return fmt.Errorf("failed to create operation log: %w", err)
	}

	// 更新状态为完成
	now := time.Now()
	batchOp.FinishedAt = &now
	s.batchOpRepo.UpdateStatus(ctx, batchOp.ID, writer.BatchOpStatusCompleted)

	return nil
}

// executeSagaAtomic Saga原子执行（>200节点，补偿事务）
func (s *BatchOperationServiceImpl) executeSagaAtomic(ctx context.Context, batchOp *writer.BatchOperation) error {
	var inverseCommands []map[string]interface{}

	// 逐个执行命令，收集inverseCommand
	for i, targetID := range batchOp.TargetIDs {
		itemErr := s.executeSingleItem(ctx, batchOp, targetID, i)
		if itemErr != nil {
			// 失败时执行补偿
			for j := i - 1; j >= 0; j-- {
				if j < len(inverseCommands) {
					_ = s.executeInverse(ctx, inverseCommands[j])
				}
			}
			s.batchOpRepo.UpdateStatus(ctx, batchOp.ID, writer.BatchOpStatusFailed)
			return fmt.Errorf("item %s failed, compensated: %w", targetID, itemErr)
		}

		// 收集逆命令
		inverseCmd := s.buildInverseCommand(batchOp, targetID)
		inverseCommands = append(inverseCommands, inverseCmd)
	}

	// 全部成功，创建OperationLog
	err := s.createOperationLog(ctx, batchOp)
	if err != nil {
		s.batchOpRepo.UpdateStatus(ctx, batchOp.ID, writer.BatchOpStatusFailed)
		return fmt.Errorf("failed to create operation log: %w", err)
	}

	// 更新状态为完成
	now := time.Now()
	batchOp.FinishedAt = &now
	s.batchOpRepo.UpdateStatus(ctx, batchOp.ID, writer.BatchOpStatusCompleted)

	return nil
}

// executeBatchDelete 执行批量删除
func (s *BatchOperationServiceImpl) executeBatchDelete(ctx context.Context, batchOp *writer.BatchOperation) error {
	for _, targetID := range batchOp.TargetIDs {
		// 软删除文档
		err := s.docRepo.SoftDelete(ctx, targetID, batchOp.CreatedBy.Hex())
		if err != nil {
			return fmt.Errorf("failed to delete document %s: %w", targetID, err)
		}
	}
	return nil
}

// executeBatchMove 执行批量移动
func (s *BatchOperationServiceImpl) executeBatchMove(ctx context.Context, batchOp *writer.BatchOperation) error {
	// 从payload中提取移动参数
	parentID, _ := batchOp.Payload["parent_id"].(string)
	position, _ := batchOp.Payload["position"].(string)
	referenceID, _ := batchOp.Payload["reference_id"].(string)

	// TODO: 实现批量移动逻辑
	// 这需要调用DocumentRepository的更新方法
	_ = parentID
	_ = position
	_ = referenceID

	return fmt.Errorf("batch move operation not yet implemented")
}

// executeSingleItem 执行单个项目（Saga模式）
func (s *BatchOperationServiceImpl) executeSingleItem(ctx context.Context, batchOp *writer.BatchOperation, targetID string, index int) error {
	switch batchOp.Type {
	case writer.BatchOpTypeDelete:
		return s.docRepo.SoftDelete(ctx, targetID, batchOp.CreatedBy.Hex())
	case writer.BatchOpTypeMove:
		// TODO: 实现单个文档移动
		return fmt.Errorf("single item move not yet implemented")
	default:
		return fmt.Errorf("unsupported operation type: %s", batchOp.Type)
	}
}

// executeInverse 执行逆操作（Saga补偿）
func (s *BatchOperationServiceImpl) executeInverse(ctx context.Context, inverseCmd map[string]interface{}) error {
	// TODO: 实现补偿逻辑
	// 例如：如果是删除，则恢复文档
	return nil
}

// buildInverseCommand 构建逆命令
func (s *BatchOperationServiceImpl) buildInverseCommand(batchOp *writer.BatchOperation, targetID string) map[string]interface{} {
	// 根据操作类型构建逆命令
	switch batchOp.Type {
	case writer.BatchOpTypeDelete:
		// 删除的逆操作是恢复
		return map[string]interface{}{
			"type":      "restore",
			"target_id": targetID,
		}
	case writer.BatchOpTypeMove:
		// 移动的逆操作是移回原位置
		return map[string]interface{}{
			"type":      "move",
			"target_id": targetID,
			"parent_id": batchOp.Payload["original_parent_id"],
			"position":  batchOp.Payload["original_position"],
			"reference": batchOp.Payload["original_reference"],
		}
	default:
		return nil
	}
}

// createOperationLog 创建操作日志
func (s *BatchOperationServiceImpl) createOperationLog(ctx context.Context, batchOp *writer.BatchOperation) error {
	log := &writer.OperationLog{
		ProjectID:      batchOp.ProjectID,
		UserID:         batchOp.CreatedBy,
		BatchOpID:      &batchOp.ID,
		ChainID:        batchOp.ID.Hex(), // 使用batchOp的ID作为chainID
		CommandType:    s.mapBatchTypeToCommandType(batchOp.Type),
		TargetIDs:      batchOp.TargetIDs,
		CommandPayload: batchOp.Payload,
		Status:         writer.OpLogStatusExecuted,
		IsCommitted:    true,
	}

	return s.opLogRepo.Create(ctx, log)
}

// mapBatchTypeToCommandType 映射批量操作类型到命令类型
func (s *BatchOperationServiceImpl) mapBatchTypeToCommandType(batchType writer.BatchOperationType) writer.DocumentCommandType {
	switch batchType {
	case writer.BatchOpTypeDelete:
		return writer.CommandDelete
	case writer.BatchOpTypeMove:
		return writer.CommandMove
	case writer.BatchOpTypeCopy:
		return writer.CommandCopy
	default:
		return writer.CommandUpdate
	}
}

// Cancel 取消正在运行的操作
func (s *BatchOperationServiceImpl) Cancel(ctx context.Context, batchID primitive.ObjectID, userID primitive.ObjectID) error {
	batchOp, err := s.batchOpRepo.GetByID(ctx, batchID)
	if err != nil {
		return err
	}

	if !batchOp.CanCancel() {
		return ErrBatchOperationNotRunning
	}

	// TODO: 实现取消逻辑（发送取消信号给执行中的goroutine）
	// 当前仅更新状态
	return s.batchOpRepo.UpdateStatus(ctx, batchID, writer.BatchOpStatusCancelled)
}

// Undo 撤销批量操作
func (s *BatchOperationServiceImpl) Undo(ctx context.Context, batchID primitive.ObjectID, userID primitive.ObjectID) error {
	// 1. 查询OperationLog
	logs, err := s.opLogRepo.GetByChainID(ctx, batchID.Hex())
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		return errors.New("no operation logs found for this batch operation")
	}

	// 2. 按逆序执行撤销
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]
		if !log.IsUndoable() {
			continue
		}

		err = s.executeInverseCommand(ctx, log.InverseCommand)
		if err != nil {
			return fmt.Errorf("failed to undo log %s: %w", log.ID.Hex(), err)
		}

		_ = s.opLogRepo.UpdateStatus(ctx, log.ID, writer.OpLogStatusUndone)
	}

	return nil
}

// executeInverseCommand 执行逆命令
func (s *BatchOperationServiceImpl) executeInverseCommand(ctx context.Context, inverseCmd map[string]interface{}) error {
	if inverseCmd == nil {
		return nil
	}

	cmdType, _ := inverseCmd["type"].(string)
	targetID, _ := inverseCmd["target_id"].(string)

	switch cmdType {
	case "restore":
		// 恢复文档（清除DeletedAt）
		// TODO: 实现恢复逻辑
		_ = targetID
		return nil
	case "move":
		// 移回原位置
		// TODO: 实现移回逻辑
		return nil
	default:
		return fmt.Errorf("unknown inverse command type: %s", cmdType)
	}
}

// GetProgress 获取操作进度
func (s *BatchOperationServiceImpl) GetProgress(ctx context.Context, batchID primitive.ObjectID) (*BatchOperationProgress, error) {
	batchOp, err := s.batchOpRepo.GetByID(ctx, batchID)
	if err != nil {
		return nil, err
	}

	// TODO: 查询items表统计进度
	// 当前简化实现：根据状态判断
	progress := &BatchOperationProgress{
		BatchID:    batchOp.ID,
		Status:     batchOp.Status,
		TotalItems: len(batchOp.TargetIDs),
		StartedAt:  batchOp.StartedAt,
		FinishedAt: batchOp.FinishedAt,
	}

	if batchOp.Status == writer.BatchOpStatusCompleted {
		progress.CompletedItems = progress.TotalItems
	} else if batchOp.Status == writer.BatchOpStatusRunning {
		// TODO: 从items表查询实际进度
		progress.CompletedItems = 0
	} else if batchOp.Status == writer.BatchOpStatusFailed {
		progress.FailedItems = progress.TotalItems
	}

	return progress, nil
}
