package admin

import (
	"context"
	"fmt"
	"time"

	adminModel "Qingyu_backend/models/users"
	adminRepo "Qingyu_backend/repository/interfaces/admin"
)

// AuditLogService 审计日志服务接口
type AuditLogService interface {
	// LogOperationWithAudit 记录带审计信息的操作日志
	LogOperationWithAudit(ctx context.Context, req *LogOperationWithAuditRequest) error

	// QueryAuditLogs 查询审计日志
	QueryAuditLogs(ctx context.Context, req *QueryAuditLogsRequest) ([]*adminModel.AdminLog, int64, error)

	// GetLogsByResource 获取资源相关日志
	GetLogsByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error)

	// GetLogsByDateRange 获取日期范围内的日志
	GetLogsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error)

	// CleanOldLogs 清理旧日志
	CleanOldLogs(ctx context.Context, beforeDate time.Time) error
}

// LogOperationWithAuditRequest 带审计信息的操作日志请求
type LogOperationWithAuditRequest struct {
	AdminID      string                 `json:"admin_id" binding:"required"`
	AdminName    string                 `json:"admin_name"`
	Operation    string                 `json:"operation" binding:"required"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	OldValues    map[string]interface{} `json:"old_values,omitempty"`
	NewValues    map[string]interface{} `json:"new_values,omitempty"`
	Changes      map[string]adminModel.ChangeRecord
	IP           string                 `json:"ip"`
	UserAgent    string                 `json:"user_agent"`
	// 敏感操作相关字段
	IsSensitive bool   `json:"is_sensitive,omitempty"` // 是否为敏感操作
	BatchID     string `json:"batch_id,omitempty"`     // 批次ID，用于批量操作去重
}

// QueryAuditLogsRequest 查询审计日志请求
type QueryAuditLogsRequest struct {
	AdminID      string     `json:"admin_id"`
	Operation    string     `json:"operation"`
	ResourceType string     `json:"resource_type"`
	ResourceID   string     `json:"resource_id"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Keyword      string     `json:"keyword"`
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
}

// auditLogServiceImpl 审计日志服务实现
type auditLogServiceImpl struct {
	logRepo adminRepo.AdminLogRepository
}

// NewAuditLogService 创建审计日志服务
func NewAuditLogService(logRepo adminRepo.AdminLogRepository) AuditLogService {
	return &auditLogServiceImpl{
		logRepo: logRepo,
	}
}

// LogOperationWithAudit 记录带审计信息的操作日志
func (s *auditLogServiceImpl) LogOperationWithAudit(ctx context.Context, req *LogOperationWithAuditRequest) error {
	if req == nil {
		return fmt.Errorf("请求参数不能为空")
	}

	// 如果提供了旧值和新值，自动生成变更记录
	if req.Changes == nil && len(req.OldValues) > 0 && len(req.NewValues) > 0 {
		req.Changes = s.generateChanges(req.OldValues, req.NewValues)
	}

	// 构建日志实体
	log := &adminModel.AdminLog{
		AdminID:      req.AdminID,
		AdminName:    req.AdminName,
		Operation:    req.Operation,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		OldValues:    req.OldValues,
		NewValues:    req.NewValues,
		Changes:      req.Changes,
		IP:           req.IP,
		UserAgent:    req.UserAgent,
		CreatedAt:    time.Now(),
	}

	// 将变更记录也转换为 Details 字段存储（保持向后兼容）
	if len(req.Changes) > 0 {
		if log.Details == nil {
			log.Details = make(map[string]interface{})
		}
		for field, change := range req.Changes {
			log.Details[field] = map[string]interface{}{
				"old": change.OldValue,
				"new": change.NewValue,
			}
		}
	}

	return s.logRepo.CreateAdminLog(ctx, log)
}

// QueryAuditLogs 查询审计日志
func (s *auditLogServiceImpl) QueryAuditLogs(ctx context.Context, req *QueryAuditLogsRequest) ([]*adminModel.AdminLog, int64, error) {
	if req == nil {
		return nil, 0, fmt.Errorf("请求参数不能为空")
	}

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建过滤条件
	filter := &adminRepo.AdminLogFilter{
		AdminID:      req.AdminID,
		Operation:    req.Operation,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Limit:        int64(pageSize),
		Offset:       int64((page - 1) * pageSize),
	}

	if req.StartDate != nil {
		filter.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		filter.EndDate = *req.EndDate
	}

	// 查询日志
	logs, err := s.logRepo.ListAdminLogs(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("查询日志失败: %w", err)
	}

	// 统计总数
	total, err := s.logRepo.CountAdminLogs(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计日志数量失败: %w", err)
	}

	return logs, total, nil
}

// GetLogsByResource 获取资源相关日志
func (s *auditLogServiceImpl) GetLogsByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error) {
	if resourceType == "" {
		return nil, fmt.Errorf("资源类型不能为空")
	}
	if resourceID == "" {
		return nil, fmt.Errorf("资源ID不能为空")
	}

	return s.logRepo.GetByResource(ctx, resourceType, resourceID)
}

// GetLogsByDateRange 获取日期范围内的日志
func (s *auditLogServiceImpl) GetLogsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, fmt.Errorf("日期范围不能为空")
	}
	if startDate.After(endDate) {
		return nil, fmt.Errorf("开始日期不能晚于结束日期")
	}

	return s.logRepo.GetByDateRange(ctx, startDate, endDate)
}

// CleanOldLogs 清理旧日志
func (s *auditLogServiceImpl) CleanOldLogs(ctx context.Context, beforeDate time.Time) error {
	if beforeDate.IsZero() {
		return fmt.Errorf("清理日期不能为空")
	}

	return s.logRepo.CleanOldLogs(ctx, beforeDate)
}

// generateChanges 根据旧值和新值生成变更记录
func (s *auditLogServiceImpl) generateChanges(oldValues, newValues map[string]interface{}) map[string]adminModel.ChangeRecord {
	changes := make(map[string]adminModel.ChangeRecord)

	// 收集所有字段
	allFields := make(map[string]bool)
	for field := range oldValues {
		allFields[field] = true
	}
	for field := range newValues {
		allFields[field] = true
	}

	// 生成变更记录
	for field := range allFields {
		oldValue, oldExists := oldValues[field]
		newValue, newExists := newValues[field]

		// 如果值发生变化，记录变更
		if !oldExists || !newExists || oldValue != newValue {
			changes[field] = adminModel.ChangeRecord{
				Field:    field,
				OldValue: oldValue,
				NewValue: newValue,
			}
		}
	}

	return changes
}
