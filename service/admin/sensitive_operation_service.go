package admin

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// SensitiveOperationService 敏感操作服务接口
type SensitiveOperationService interface {
	// IsSensitiveOperation 检查是否是敏感操作
	IsSensitiveOperation(action, resourceType string) bool

	// LogSensitiveOperation 记录敏感操作
	LogSensitiveOperation(ctx context.Context, req *LogOperationWithAuditRequest) error

	// AddToWhitelist 添加到白名单
	AddToWhitelist(action, resourceType string) error

	// RemoveFromWhitelist 从白名单移除
	RemoveFromWhitelist(action, resourceType string) error
}

// sensitiveOperationServiceImpl 敏感操作服务实现
type sensitiveOperationServiceImpl struct {
	auditLogService AuditLogService
	whitelist       map[string]bool
	whitelistMutex  sync.RWMutex
	alertedBatches  map[string]bool // 已警告的批次ID，用于批量操作去重
	batchMutex      sync.RWMutex
}

// 敏感操作列表
var sensitiveOperations = map[string]bool{
	"delete:user":       true,
	"update:role":       true,
	"delete:content":    true,
	"update:system":     true,
	"update:permission": true,
	"delete:role":       true,
	"delete:comment":    true,
	"ban:user":          true,
}

// NewSensitiveOperationService 创建敏感操作服务
func NewSensitiveOperationService(auditLogService AuditLogService) SensitiveOperationService {
	return &sensitiveOperationServiceImpl{
		auditLogService: auditLogService,
		whitelist:       make(map[string]bool),
		alertedBatches:   make(map[string]bool),
	}
}

// IsSensitiveOperation 检查是否是敏感操作
func (s *sensitiveOperationServiceImpl) IsSensitiveOperation(action, resourceType string) bool {
	// 检查白名单
	key := strings.ToLower(action + ":" + resourceType)

	s.whitelistMutex.RLock()
	_, inWhitelist := s.whitelist[key]
	s.whitelistMutex.RUnlock()

	if inWhitelist {
		return false
	}

	// 检查敏感操作列表
	return sensitiveOperations[key]
}

// LogSensitiveOperation 记录敏感操作
func (s *sensitiveOperationServiceImpl) LogSensitiveOperation(ctx context.Context, req *LogOperationWithAuditRequest) error {
	if req == nil {
		return fmt.Errorf("请求参数不能为空")
	}

	// 检查是否为敏感操作
	isSensitive := s.IsSensitiveOperation(req.Operation, req.ResourceType)

	// 批量操作去重处理
	if req.BatchID != "" && isSensitive {
		s.batchMutex.Lock()
		if s.alertedBatches[req.BatchID] {
			// 该批次已警告过，不再重复警告
			isSensitive = false
		} else {
			// 标记该批次已警告
			s.alertedBatches[req.BatchID] = true
		}
		s.batchMutex.Unlock()
	}

	// 标记敏感操作
	req.IsSensitive = isSensitive

	// 如果是敏感操作，添加额外信息到 Details
	if isSensitive {
		if req.OldValues == nil {
			req.OldValues = make(map[string]interface{})
		}
		req.OldValues["_sensitive_operation"] = map[string]interface{}{
			"action":        req.Operation,
			"resource_type": req.ResourceType,
			"detected_at":   ctx.Value("timestamp"),
		}
	}

	// 记录日志
	return s.auditLogService.LogOperationWithAudit(ctx, req)
}

// AddToWhitelist 添加到白名单
func (s *sensitiveOperationServiceImpl) AddToWhitelist(action, resourceType string) error {
	if action == "" || resourceType == "" {
		return fmt.Errorf("操作类型和资源类型不能为空")
	}

	key := strings.ToLower(action + ":" + resourceType)

	s.whitelistMutex.Lock()
	s.whitelist[key] = true
	s.whitelistMutex.Unlock()

	return nil
}

// RemoveFromWhitelist 从白名单移除
func (s *sensitiveOperationServiceImpl) RemoveFromWhitelist(action, resourceType string) error {
	if action == "" || resourceType == "" {
		return fmt.Errorf("操作类型和资源类型不能为空")
	}

	key := strings.ToLower(action + ":" + resourceType)

	s.whitelistMutex.Lock()
	delete(s.whitelist, key)
	s.whitelistMutex.Unlock()

	return nil
}

// GetSensitiveOperations 获取敏感操作列表（用于管理和展示）
func GetSensitiveOperations() []string {
	operations := make([]string, 0, len(sensitiveOperations))
	for op := range sensitiveOperations {
		operations = append(operations, op)
	}
	return operations
}

// AddSensitiveOperation 添加敏感操作到列表（用于动态配置）
func AddSensitiveOperation(action, resourceType string) {
	key := strings.ToLower(action + ":" + resourceType)
	sensitiveOperations[key] = true
}

// RemoveSensitiveOperation 从敏感操作列表移除（用于动态配置）
func RemoveSensitiveOperation(action, resourceType string) {
	key := strings.ToLower(action + ":" + resourceType)
	delete(sensitiveOperations, key)
}
