package admin

import (
	"context"
	"testing"
	"time"

	adminModel "Qingyu_backend/models/users"
	adminRepo "Qingyu_backend/repository/interfaces/admin"
)

// MockAuditLogRepository Mock审计日志仓储
type MockAuditLogRepository struct {
	CreateFunc         func(ctx context.Context, log *adminModel.AdminLog) error
	ListFunc           func(ctx context.Context, filter *adminRepo.AdminLogFilter) ([]*adminModel.AdminLog, error)
	CountFunc          func(ctx context.Context, filter *adminRepo.AdminLogFilter) (int64, error)
	GetByResourceFunc  func(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error)
	GetByDateRangeFunc func(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error)
	CleanOldLogsFunc   func(ctx context.Context, beforeDate time.Time) error
}

func (m *MockAuditLogRepository) CreateAdminLog(ctx context.Context, log *adminModel.AdminLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, log)
	}
	return nil
}

func (m *MockAuditLogRepository) ListAdminLogs(ctx context.Context, filter *adminRepo.AdminLogFilter) ([]*adminModel.AdminLog, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter)
	}
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogRepository) CountAdminLogs(ctx context.Context, filter *adminRepo.AdminLogFilter) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, filter)
	}
	return 0, nil
}

func (m *MockAuditLogRepository) GetAdminLog(ctx context.Context, logID string) (*adminModel.AdminLog, error) {
	return &adminModel.AdminLog{}, nil
}

func (m *MockAuditLogRepository) GetByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error) {
	if m.GetByResourceFunc != nil {
		return m.GetByResourceFunc(ctx, resourceType, resourceID)
	}
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error) {
	if m.GetByDateRangeFunc != nil {
		return m.GetByDateRangeFunc(ctx, startDate, endDate)
	}
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogRepository) CleanOldLogs(ctx context.Context, beforeDate time.Time) error {
	if m.CleanOldLogsFunc != nil {
		return m.CleanOldLogsFunc(ctx, beforeDate)
	}
	return nil
}

func (m *MockAuditLogRepository) Health(ctx context.Context) error {
	return nil
}

// TestAuditLogService_LogUserOperation_Success 测试记录用户操作
func TestAuditLogService_LogUserOperation_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAuditLogRepository{
		CreateFunc: func(ctx context.Context, log *adminModel.AdminLog) error {
			if log.AdminID != "admin123" {
				t.Errorf("期望 AdminID 为 admin123, 实际为 %s", log.AdminID)
			}
			if log.Operation != "ban_user" {
				t.Errorf("期望 Operation 为 ban_user, 实际为 %s", log.Operation)
			}
			if log.ResourceType != "user" {
				t.Errorf("期望 ResourceType 为 user, 实际为 %s", log.ResourceType)
			}
			if log.ResourceID != "user456" {
				t.Errorf("期望 ResourceID 为 user456, 实际为 %s", log.ResourceID)
			}
			return nil
		},
	}

	service := NewAuditLogService(mockRepo)

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		AdminName:    "管理员A",
		Operation:    "ban_user",
		ResourceType: "user",
		ResourceID:   "user456",
		IP:           "192.168.1.100",
		UserAgent:    "Mozilla/5.0",
	}

	err := service.LogOperationWithAudit(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}
}

// TestAuditLogService_LogContentOperation_Success 测试记录内容操作
func TestAuditLogService_LogContentOperation_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAuditLogRepository{
		CreateFunc: func(ctx context.Context, log *adminModel.AdminLog) error {
			if log.ResourceType != "content" {
				t.Errorf("期望 ResourceType 为 content, 实际为 %s", log.ResourceType)
			}
			if log.Operation != "delete_content" {
				t.Errorf("期望 Operation 为 delete_content, 实际为 %s", log.Operation)
			}
			return nil
		},
	}

	service := NewAuditLogService(mockRepo)

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		AdminName:    "管理员A",
		Operation:    "delete_content",
		ResourceType: "content",
		ResourceID:   "book789",
		IP:           "192.168.1.100",
	}

	err := service.LogOperationWithAudit(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}
}

// TestAuditLogService_LogWithChanges_Success 测试记录变更前后值
func TestAuditLogService_LogWithChanges_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAuditLogRepository{
		CreateFunc: func(ctx context.Context, log *adminModel.AdminLog) error {
			// 检查变更记录
			if log.OldValues == nil {
				t.Error("期望 OldValues 不为空")
			}
			if log.NewValues == nil {
				t.Error("期望 NewValues 不为空")
			}
			if log.Changes == nil {
				t.Error("期望 Changes 不为空")
			}
			return nil
		},
	}

	service := NewAuditLogService(mockRepo)

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		AdminName:    "管理员A",
		Operation:    "update_role",
		ResourceType: "role",
		ResourceID:   "role1",
		OldValues: map[string]interface{}{
			"name":  "普通用户",
			"level": 1,
		},
		NewValues: map[string]interface{}{
			"name":  "高级用户",
			"level": 2,
		},
		IP: "192.168.1.100",
	}

	err := service.LogOperationWithAudit(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}
}

// TestAuditLogService_QueryLogsByUser_Success 测试按用户查询
func TestAuditLogService_QueryLogsByUser_Success(t *testing.T) {
	ctx := context.Background()
	expectedLogs := []*adminModel.AdminLog{
		{
			ID:        "log1",
			AdminID:   "admin123",
			Operation: "ban_user",
			CreatedAt: time.Now(),
		},
		{
			ID:        "log2",
			AdminID:   "admin123",
			Operation: "update_role",
			CreatedAt: time.Now(),
		},
	}

	mockRepo := &MockAuditLogRepository{
		ListFunc: func(ctx context.Context, filter *adminRepo.AdminLogFilter) ([]*adminModel.AdminLog, error) {
			return expectedLogs, nil
		},
		CountFunc: func(ctx context.Context, filter *adminRepo.AdminLogFilter) (int64, error) {
			return 2, nil
		},
	}

	service := NewAuditLogService(mockRepo)

	req := &QueryAuditLogsRequest{
		AdminID:  "admin123",
		Page:     1,
		PageSize: 10,
	}

	logs, total, err := service.QueryAuditLogs(ctx, req)
	if err != nil {
		t.Fatalf("期望查询成功, 但得到错误: %v", err)
	}
	if total != 2 {
		t.Errorf("期望总数为 2, 实际为 %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("期望返回 2 条日志, 实际为 %d", len(logs))
	}
	if logs[0].AdminID != "admin123" {
		t.Errorf("期望 AdminID 为 admin123, 实际为 %s", logs[0].AdminID)
	}
}

// TestAuditLogService_QueryLogsByDateRange_Success 测试按日期查询
func TestAuditLogService_QueryLogsByDateRange_Success(t *testing.T) {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	expectedLogs := []*adminModel.AdminLog{
		{
			ID:        "log1",
			AdminID:   "admin123",
			Operation: "ban_user",
			CreatedAt: time.Now().AddDate(0, 0, -15),
		},
	}

	mockRepo := &MockAuditLogRepository{
		ListFunc: func(ctx context.Context, filter *adminRepo.AdminLogFilter) ([]*adminModel.AdminLog, error) {
			return expectedLogs, nil
		},
		CountFunc: func(ctx context.Context, filter *adminRepo.AdminLogFilter) (int64, error) {
			return 1, nil
		},
	}

	service := NewAuditLogService(mockRepo)

	req := &QueryAuditLogsRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
		Page:      1,
		PageSize:  10,
	}

	logs, total, err := service.QueryAuditLogs(ctx, req)
	if err != nil {
		t.Fatalf("期望查询成功, 但得到错误: %v", err)
	}
	if total != 1 {
		t.Errorf("期望总数为 1, 实际为 %d", total)
	}
	if len(logs) != 1 {
		t.Errorf("期望返回 1 条日志, 实际为 %d", len(logs))
	}
}

// TestAuditLogService_GetLogsByResource_Success 测试按资源查询
func TestAuditLogService_GetLogsByResource_Success(t *testing.T) {
	ctx := context.Background()
	expectedLogs := []*adminModel.AdminLog{
		{
			ID:           "log1",
			AdminID:      "admin123",
			Operation:    "update_role",
			ResourceType: "role",
			ResourceID:   "role1",
			CreatedAt:    time.Now(),
		},
	}

	mockRepo := &MockAuditLogRepository{
		GetByResourceFunc: func(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error) {
			if resourceType != "role" {
				t.Errorf("期望 ResourceType 为 role, 实际为 %s", resourceType)
			}
			if resourceID != "role1" {
				t.Errorf("期望 ResourceID 为 role1, 实际为 %s", resourceID)
			}
			return expectedLogs, nil
		},
	}

	service := NewAuditLogService(mockRepo)

	logs, err := service.GetLogsByResource(ctx, "role", "role1")
	if err != nil {
		t.Fatalf("期望查询成功, 但得到错误: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("期望返回 1 条日志, 实际为 %d", len(logs))
	}
}

// TestAuditLogService_CleanOldLogs_Success 测试清理旧日志
func TestAuditLogService_CleanOldLogs_Success(t *testing.T) {
	ctx := context.Background()
	cleanDate := time.Now().AddDate(-1, 0, 0) // 一年前

	mockRepo := &MockAuditLogRepository{
		CleanOldLogsFunc: func(ctx context.Context, beforeDate time.Time) error {
			if !beforeDate.Before(cleanDate.Add(time.Minute)) && !beforeDate.After(cleanDate.Add(-time.Minute)) {
				t.Errorf("期望清理日期接近 %v, 实际为 %v", cleanDate, beforeDate)
			}
			return nil
		},
	}

	service := NewAuditLogService(mockRepo)

	err := service.CleanOldLogs(ctx, cleanDate)
	if err != nil {
		t.Fatalf("期望清理成功, 但得到错误: %v", err)
	}
}
