package document

import (
	"Qingyu_backend/models/writer"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	pkgErrors "Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	serviceBase "Qingyu_backend/service/base"
)

// BatchOperationService 批量操作服务
type BatchOperationService struct {
	batchOpRepo writerRepo.BatchOperationRepository
	docRepo     writerRepo.DocumentRepository
	projectRepo writerRepo.ProjectRepository
	retrySvc    *RetryService
	eventBus    serviceBase.EventBus
	serviceName string
	version     string
}

// NewBatchOperationService 创建批量操作服务
func NewBatchOperationService(
	batchOpRepo writerRepo.BatchOperationRepository,
	docRepo writerRepo.DocumentRepository,
	projectRepo writerRepo.ProjectRepository,
	eventBus serviceBase.EventBus,
) *BatchOperationService {
	return &BatchOperationService{
		batchOpRepo: batchOpRepo,
		docRepo:     docRepo,
		projectRepo: projectRepo,
		retrySvc:    NewRetryService(),
		eventBus:    eventBus,
		serviceName: "BatchOperationService",
		version:     "1.0.0",
	}
}

// Submit 提交批量操作
func (s *BatchOperationService) Submit(ctx context.Context, req *SubmitBatchOperationRequest) (*writer.BatchOperation, error) {
	// 1. 参数验证
	if err := s.validateSubmitRequest(req); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数验证失败", err.Error(), err)
	}

	// 2. 幂等性检查（如果提供了clientRequestID）
	if req.ClientRequestID != "" {
		existingOp, err := s.batchOpRepo.GetByClientRequestID(ctx, req.ClientRequestID)
		if err != nil {
			return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "检查幂等性失败", "", err)
		}
		if existingOp != nil {
			// 已存在相同请求，返回现有操作
			return existingOp, nil
		}
	}

	// 3. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanEdit(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限编辑该项目", "", nil)
	}

	// 4. 转换ProjectID
	projectID, err := primitive.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的项目ID", "", err)
	}

	// 5. 创建批量操作对象
	retryConfigData := s.convertRetryConfigToMap(req.RetryConfig)
	batchOp := &writer.BatchOperation{
		ProjectID:       projectID,
		Type:            req.Type,
		TargetIDs:       req.TargetIDs,
		Status:          writer.BatchOpStatusPending,
		Payload:         req.Payload,
		Atomic:          req.Atomic,
		ConflictPolicy:  req.ConflictPolicy,
		ClientRequestID: req.ClientRequestID,
		RetryConfig:     retryConfigData,
		CreatedBy:       userID,
	}

	// 初始化items
	batchOp.TouchForCreate()

	// 6. 预检查（Preflight）
	preflightSummary, err := s.runPreflight(ctx, batchOp)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "预检查失败", "", err)
	}
	batchOp.PreflightSummary = preflightSummary

	// 如果预检查发现所有项都无效，直接返回失败
	if preflightSummary.ValidCount == 0 {
		batchOp.Status = writer.BatchOpStatusFailed
		batchOp.ErrorCode = "PREFLIGHT_FAILED"
		batchOp.ErrorMessage = "没有有效的操作项"
	}

	// 7. 保存批量操作
	if err := s.batchOpRepo.Create(ctx, batchOp); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建批量操作失败", "", err)
	}

	// 8. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "batch_operation.created",
			EventData: map[string]interface{}{
				"operation_id": batchOp.ID.Hex(),
				"project_id":   batchOp.ProjectID.Hex(),
				"type":         string(batchOp.Type),
				"target_count": len(batchOp.TargetIDs),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return batchOp, nil
}

// Execute 执行批量操作
func (s *BatchOperationService) Execute(ctx context.Context, operationID string) error {
	// 1. 获取批量操作
	batchOp, err := s.batchOpRepo.GetByID(ctx, operationID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询批量操作失败", "", err)
	}

	if batchOp == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "批量操作不存在", "", nil)
	}

	// 2. 检查状态
	if batchOp.Status != writer.BatchOpStatusPending && batchOp.Status != writer.BatchOpStatusPreflight {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "批量操作状态不正确，无法执行", "", nil)
	}

	// 3. 更新状态为处理中
	now := primitive.NewDateTimeFromTime(time.Now())
	batchOp.Status = writer.BatchOpStatusProcessing
	batchOp.StartedAt = &now
	if err := s.batchOpRepo.UpdateStatus(ctx, operationID, batchOp.Status); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新操作状态失败", "", err)
	}

	// 4. 根据操作类型执行不同的逻辑
	var executeErr error
	switch batchOp.Type {
	case writer.BatchOpTypeDelete:
		executeErr = s.executeDelete(ctx, batchOp)
	case writer.BatchOpTypeMove:
		executeErr = s.executeMove(ctx, batchOp)
	case writer.BatchOpTypeExport:
		executeErr = s.executeExport(ctx, batchOp)
	case writer.BatchOpTypeCopy:
		executeErr = s.executeCopy(ctx, batchOp)
	case writer.BatchOpTypeApply:
		executeErr = s.executeApplyTemplate(ctx, batchOp)
	default:
		executeErr = fmt.Errorf("不支持的操作类型: %s", batchOp.Type)
	}

	// 5. 更新最终状态
	completedAt := primitive.NewDateTimeFromTime(time.Now())
	batchOp.CompletedAt = &completedAt

	if executeErr != nil {
		if batchOp.Atomic {
			// 原子操作：全部失败
			batchOp.Status = writer.BatchOpStatusFailed
			batchOp.ErrorCode = "EXECUTION_FAILED"
			batchOp.ErrorMessage = executeErr.Error()
		} else {
			// 非原子操作：部分成功
			successCount := 0
			failedCount := 0
			for _, item := range batchOp.Items {
				if item.Status == writer.BatchItemStatusSucceeded {
					successCount++
				} else if item.Status == writer.BatchItemStatusFailed {
					failedCount++
				}
			}

			if successCount > 0 && failedCount > 0 {
				batchOp.Status = writer.BatchOpStatusPartial
			} else if successCount == 0 {
				batchOp.Status = writer.BatchOpStatusFailed
			} else {
				batchOp.Status = writer.BatchOpStatusCompleted
			}
		}
	} else {
		// 全部成功
		batchOp.Status = writer.BatchOpStatusCompleted
	}

	// 更新数据库
	updates := map[string]interface{}{
		"status":       batchOp.Status,
		"completed_at": batchOp.CompletedAt,
		"items":        batchOp.Items,
	}
	if batchOp.ErrorCode != "" {
		updates["error_code"] = batchOp.ErrorCode
		updates["error_message"] = batchOp.ErrorMessage
	}
	if batchOp.PreflightSummary != nil {
		updates["preflight_summary"] = batchOp.PreflightSummary
	}

	if err := s.batchOpRepo.Update(ctx, operationID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新操作结果失败", "", err)
	}

	// 6. 发布完成事件
	if s.eventBus != nil {
		eventType := "batch_operation.completed"
		if batchOp.Status == writer.BatchOpStatusFailed {
			eventType = "batch_operation.failed"
		} else if batchOp.Status == writer.BatchOpStatusPartial {
			eventType = "batch_operation.partial"
		}

		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: eventType,
			EventData: map[string]interface{}{
				"operation_id": operationID,
				"project_id":   batchOp.ProjectID.Hex(),
				"type":         string(batchOp.Type),
				"status":       string(batchOp.Status),
				"summary":      batchOp.GetSummary(),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return executeErr
}

// GetOperation 获取批量操作详情
func (s *BatchOperationService) GetOperation(ctx context.Context, operationID string) (*writer.BatchOperation, error) {
	// 1. 查询批量操作
	batchOp, err := s.batchOpRepo.GetByID(ctx, operationID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询批量操作失败", "", err)
	}

	if batchOp == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "批量操作不存在", "", nil)
	}

	// 2. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, batchOp.ProjectID.Hex())
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanView(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限查看该操作", "", nil)
	}

	return batchOp, nil
}

// ListOperations 获取批量操作列表
func (s *BatchOperationService) ListOperations(ctx context.Context, req *ListBatchOperationsRequest) (*ListBatchOperationsResponse, error) {
	// 1. 验证项目权限
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanView(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限查看该项目", "", nil)
	}

	// 2. 查询批量操作列表
	var operations []*writer.BatchOperation
	if req.Status != "" {
		operations, err = s.batchOpRepo.GetByProjectAndStatus(ctx, req.ProjectID, req.Status, req.Limit, req.Offset)
	} else if req.Type != "" {
		operations, err = s.batchOpRepo.GetByProjectAndType(ctx, req.ProjectID, req.Type, req.Limit, req.Offset)
	} else {
		operations, err = s.batchOpRepo.GetByProjectID(ctx, req.ProjectID, req.Limit, req.Offset)
	}

	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询批量操作列表失败", "", err)
	}

	return &ListBatchOperationsResponse{
		Operations: operations,
		Total:      len(operations),
	}, nil
}

// Cancel 取消批量操作
func (s *BatchOperationService) Cancel(ctx context.Context, operationID string) error {
	// 1. 获取批量操作
	batchOp, err := s.batchOpRepo.GetByID(ctx, operationID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询批量操作失败", "", err)
	}

	if batchOp == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "批量操作不存在", "", nil)
	}

	// 2. 检查状态
	if batchOp.IsCompleted() {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorBusiness, "操作已完成，无法取消", "", nil)
	}

	// 3. 更新状态为已取消
	if err := s.batchOpRepo.UpdateStatus(ctx, operationID, writer.BatchOpStatusCancelled); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "取消操作失败", "", err)
	}

	// 4. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &serviceBase.BaseEvent{
			EventType: "batch_operation.cancelled",
			EventData: map[string]interface{}{
				"operation_id": operationID,
				"project_id":   batchOp.ProjectID.Hex(),
				"type":         string(batchOp.Type),
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// 私有方法

// convertRetryConfigToMap 将RetryConfig转换为map
func (s *BatchOperationService) convertRetryConfigToMap(config *RetryConfig) map[string]interface{} {
	if config == nil {
		return nil
	}
	return map[string]interface{}{
		"maxRetries":      config.MaxRetries,
		"retryDelay":      config.RetryDelay,
		"retryableErrors": config.RetryableErrors,
	}
}

// convertMapToRetryConfig 将map转换为RetryConfig
func (s *BatchOperationService) convertMapToRetryConfig(data map[string]interface{}) *RetryConfig {
	if data == nil {
		return nil
	}

	config := &RetryConfig{}
	if maxRetries, ok := data["maxRetries"].(int); ok {
		config.MaxRetries = maxRetries
	}
	if retryDelay, ok := data["retryDelay"].(int); ok {
		config.RetryDelay = retryDelay
	}
	if retryableErrors, ok := data["retryableErrors"].([]string); ok {
		config.RetryableErrors = retryableErrors
	}

	return config
}

// validateSubmitRequest 验证提交请求
func (s *BatchOperationService) validateSubmitRequest(req *SubmitBatchOperationRequest) error {
	if req.ProjectID == "" {
		return fmt.Errorf("项目ID不能为空")
	}
	if req.Type == "" {
		return fmt.Errorf("操作类型不能为空")
	}
	if len(req.TargetIDs) == 0 {
		return fmt.Errorf("目标ID列表不能为空")
	}
	if len(req.TargetIDs) > 1000 {
		return fmt.Errorf("目标ID列表不能超过1000个")
	}
	return nil
}

// runPreflight 运行预检查
func (s *BatchOperationService) runPreflight(ctx context.Context, batchOp *writer.BatchOperation) (*writer.PreflightSummary, error) {
	summary := &writer.PreflightSummary{
		TotalCount: len(batchOp.TargetIDs),
		ValidCount: 0,
		InvalidCount: 0,
		SkippedCount: 0,
	}

	// 验证每个目标ID
	for _, targetID := range batchOp.TargetIDs {
		doc, err := s.docRepo.GetByID(ctx, targetID)
		if err != nil {
			summary.InvalidCount++
			summary.Errors = append(summary.Errors, fmt.Sprintf("目标 %s 验证失败: %v", targetID, err))
			continue
		}

		if doc == nil {
			summary.InvalidCount++
			summary.Errors = append(summary.Errors, fmt.Sprintf("目标 %s 不存在", targetID))
			continue
		}

		// 检查是否属于同一项目
		if doc.ProjectID != batchOp.ProjectID {
			summary.InvalidCount++
			summary.Errors = append(summary.Errors, fmt.Sprintf("目标 %s 不属于该项目", targetID))
			continue
		}

		summary.ValidCount++
	}

	return summary, nil
}

// executeDelete 执行批量删除
func (s *BatchOperationService) executeDelete(ctx context.Context, batchOp *writer.BatchOperation) error {
	// 转换RetryConfig
	retryConfig := s.convertMapToRetryConfig(batchOp.RetryConfig)

	for i := range batchOp.Items {
		item := &batchOp.Items[i]

		// 更新状态为处理中
		item.Status = writer.BatchItemStatusProcessing
		s.batchOpRepo.UpdateItemStatus(ctx, batchOp.ID.Hex(), item.TargetID, item.Status, "", "")

		// 执行删除
		err := s.docRepo.SoftDelete(ctx, item.TargetID, batchOp.ProjectID.Hex())
		if err != nil {
			item.Status = writer.BatchItemStatusFailed
			item.ErrorCode = "DELETE_FAILED"
			item.ErrorMsg = err.Error()
			item.Retryable = s.retrySvc.ShouldRetry(err, retryConfig)

			s.batchOpRepo.UpdateItemStatus(ctx, batchOp.ID.Hex(), item.TargetID, item.Status, item.ErrorCode, item.ErrorMsg)

			if batchOp.Atomic {
				return fmt.Errorf("删除目标 %s 失败: %w", item.TargetID, err)
			}
			continue
		}

		// 成功
		item.Status = writer.BatchItemStatusSucceeded
		s.batchOpRepo.UpdateItemStatus(ctx, batchOp.ID.Hex(), item.TargetID, item.Status, "", "")
	}

	return nil
}

// executeMove 执行批量移动
func (s *BatchOperationService) executeMove(ctx context.Context, batchOp *writer.BatchOperation) error {
	// TODO: 实现批量移动逻辑
	// 这将在Task 2.1中实现
	return fmt.Errorf("批量移动功能尚未实现")
}

// executeExport 执行批量导出
func (s *BatchOperationService) executeExport(ctx context.Context, batchOp *writer.BatchOperation) error {
	// TODO: 实现批量导出逻辑
	// 这将在Task 3.x中实现
	return fmt.Errorf("批量导出功能尚未实现")
}

// executeCopy 执行批量复制
func (s *BatchOperationService) executeCopy(ctx context.Context, batchOp *writer.BatchOperation) error {
	// TODO: 实现批量复制逻辑
	// 这将在Task 2.3中实现
	return fmt.Errorf("批量复制功能尚未实现")
}

// executeApplyTemplate 执行批量应用模板
func (s *BatchOperationService) executeApplyTemplate(ctx context.Context, batchOp *writer.BatchOperation) error {
	// TODO: 实现批量应用模板逻辑
	// 这将在Task 3.x中实现
	return fmt.Errorf("批量应用模板功能尚未实现")
}

// BaseService接口实现
func (s *BatchOperationService) Initialize(ctx context.Context) error {
	return nil
}

func (s *BatchOperationService) Health(ctx context.Context) error {
	return s.batchOpRepo.Health(ctx)
}

func (s *BatchOperationService) Close(ctx context.Context) error {
	return nil
}

func (s *BatchOperationService) GetServiceName() string {
	return s.serviceName
}

func (s *BatchOperationService) GetVersion() string {
	return s.version
}

// 请求和响应DTO

// SubmitBatchOperationRequest 提交批量操作请求
type SubmitBatchOperationRequest struct {
	ProjectID       string                      `json:"projectId" validate:"required"`
	Type            writer.BatchOperationType   `json:"type" validate:"required"`
	TargetIDs       []string                    `json:"targetIds" validate:"required,min=1,max=1000"`
	Payload         map[string]interface{}      `json:"payload,omitempty"`
	Atomic          bool                        `json:"atomic"`
	ConflictPolicy  writer.ConflictPolicy       `json:"conflictPolicy,omitempty"`
	ClientRequestID string                      `json:"clientRequestId,omitempty"`
	RetryConfig     *RetryConfig                `json:"retryConfig,omitempty"`
}

// ListBatchOperationsRequest 获取批量操作列表请求
type ListBatchOperationsRequest struct {
	ProjectID string                        `json:"projectId" validate:"required"`
	Type      writer.BatchOperationType     `json:"type,omitempty"`
	Status    writer.BatchOperationStatus   `json:"status,omitempty"`
	Limit     int64                         `json:"limit,omitempty"`
	Offset    int64                         `json:"offset,omitempty"`
}

// ListBatchOperationsResponse 获取批量操作列表响应
type ListBatchOperationsResponse struct {
	Operations []*writer.BatchOperation `json:"operations"`
	Total      int                     `json:"total"`
}
