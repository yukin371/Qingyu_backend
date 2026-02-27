package admin

import (
	"context"
	"testing"
	"time"

	adminModel "Qingyu_backend/models/users"
)

// MockAuditLogServiceForSensitive Mock审计日志服务（用于敏感操作服务测试）
type MockAuditLogServiceForSensitive struct {
	LogFunc func(ctx context.Context, req *LogOperationWithAuditRequest) error
}

func (m *MockAuditLogServiceForSensitive) LogOperationWithAudit(ctx context.Context, req *LogOperationWithAuditRequest) error {
	if m.LogFunc != nil {
		return m.LogFunc(ctx, req)
	}
	return nil
}

func (m *MockAuditLogServiceForSensitive) QueryAuditLogs(ctx context.Context, req *QueryAuditLogsRequest) ([]*adminModel.AdminLog, int64, error) {
	return []*adminModel.AdminLog{}, 0, nil
}

func (m *MockAuditLogServiceForSensitive) GetLogsByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error) {
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogServiceForSensitive) GetLogsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error) {
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogServiceForSensitive) CleanOldLogs(ctx context.Context, beforeDate time.Time) error {
	return nil
}

// TestSensitiveOperationService_DetectAndNotify_Success 测试检测并通知
func TestSensitiveOperationService_DetectAndNotify_Success(t *testing.T) {
	ctx := context.Background()
	notified := false

	mockService := &MockAuditLogServiceForSensitive{
		LogFunc: func(ctx context.Context, req *LogOperationWithAuditRequest) error {
			// 验证敏感操作被标记
			if !req.IsSensitive {
				t.Error("期望 IsSensitive 为 true")
			}
			notified = true
			return nil
		},
	}

	service := NewSensitiveOperationService(mockService)

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		Operation:    "delete",
		ResourceType: "user",
		ResourceID:   "user456",
		IP:           "192.168.1.100",
	}

	err := service.LogSensitiveOperation(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}

	if !notified {
		t.Error("期望触发敏感操作通知")
	}
}

// TestSensitiveOperationService_LogAndAlert_Success 测试记录并警告
func TestSensitiveOperationService_LogAndAlert_Success(t *testing.T) {
	ctx := context.Background()
	var loggedReq *LogOperationWithAuditRequest

	mockService := &MockAuditLogServiceForSensitive{
		LogFunc: func(ctx context.Context, req *LogOperationWithAuditRequest) error {
			loggedReq = req
			return nil
		},
	}

	service := NewSensitiveOperationService(mockService)

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		Operation:    "update",
		ResourceType: "role",
		ResourceID:   "role1",
		IP:           "192.168.1.100",
	}

	err := service.LogSensitiveOperation(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}

	if !loggedReq.IsSensitive {
		t.Error("期望 IsSensitive 为 true")
	}

	// 验证 OldValues 中包含敏感操作标记
	if loggedReq.OldValues == nil {
		t.Error("期望 OldValues 不为空")
	}
}

// TestSensitiveOperationService_WhitelistedOperation_NoAlert 测试白名单操作不警告
func TestSensitiveOperationService_WhitelistedOperation_NoAlert(t *testing.T) {
	ctx := context.Background()
	notified := false

	mockService := &MockAuditLogServiceForSensitive{
		LogFunc: func(ctx context.Context, req *LogOperationWithAuditRequest) error {
			if req.IsSensitive {
				notified = true
			}
			return nil
		},
	}

	service := NewSensitiveOperationService(mockService)

	// 将 "update:user" 添加到白名单
	service.AddToWhitelist("update", "user")

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		Operation:    "update",
		ResourceType: "user",
		ResourceID:   "user456",
		IP:           "192.168.1.100",
	}

	err := service.LogSensitiveOperation(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}

	if notified {
		t.Error("白名单操作不应该触发敏感操作警告")
	}
}

// TestSensitiveOperationService_BatchOperation_AlertOnce 测试批量操作只警告一次
func TestSensitiveOperationService_BatchOperation_AlertOnce(t *testing.T) {
	ctx := context.Background()
	alertCount := 0

	mockService := &MockAuditLogServiceForSensitive{
		LogFunc: func(ctx context.Context, req *LogOperationWithAuditRequest) error {
			if req.IsSensitive {
				alertCount++
			}
			return nil
		},
	}

	service := NewSensitiveOperationService(mockService)

	// 模拟批量删除用户操作
	operations := []string{"user1", "user2", "user3"}

	for _, userID := range operations {
		req := &LogOperationWithAuditRequest{
			AdminID:      "admin123",
			Operation:    "delete",
			ResourceType: "user",
			ResourceID:   userID,
			IP:           "192.168.1.100",
			BatchID:      "batch123", // 相同的批次ID
		}

		err := service.LogSensitiveOperation(ctx, req)
		if err != nil {
			t.Fatalf("期望记录成功, 但得到错误: %v", err)
		}
	}

	// 批量操作应该只警告一次
	if alertCount != 1 {
		t.Errorf("期望警告 1 次, 实际警告 %d 次", alertCount)
	}
}

// TestSensitiveOperationService_IsSensitiveOperation 测试敏感操作检测
func TestSensitiveOperationService_IsSensitiveOperation(t *testing.T) {
	service := NewSensitiveOperationService(nil)

	testCases := []struct {
		name         string
		action       string
		resourceType string
		expect       bool
	}{
		{"删除用户", "delete", "user", true},
		{"修改角色", "update", "role", true},
		{"删除内容", "delete", "content", true},
		{"修改系统配置", "update", "system", true},
		{"修改权限", "update", "permission", true},
		{"普通查询", "query", "user", false},
		{"更新用户信息", "update", "user", false},
		{"创建角色", "create", "role", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.IsSensitiveOperation(tc.action, tc.resourceType)
			if result != tc.expect {
				t.Errorf("操作 %s:%s 期望敏感=%v, 实际=%v", tc.action, tc.resourceType, tc.expect, result)
			}
		})
	}
}

// TestSensitiveOperationService_RemoveFromWhitelist 测试从白名单移除
func TestSensitiveOperationService_RemoveFromWhitelist(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{
		LogFunc: func(ctx context.Context, req *LogOperationWithAuditRequest) error {
			return nil
		},
	}

	service := NewSensitiveOperationService(mockService)

	// 添加到白名单
	service.AddToWhitelist("delete", "user")

	// 验证不在敏感操作列表中
	if service.IsSensitiveOperation("delete", "user") {
		t.Error("添加到白名单后, 不应该被识别为敏感操作")
	}

	// 从白名单移除
	service.RemoveFromWhitelist("delete", "user")

	// 验证重新变为敏感操作
	if !service.IsSensitiveOperation("delete", "user") {
		t.Error("从白名单移除后, 应该被识别为敏感操作")
	}
}

// TestSensitiveOperationService_AddToWhitelist_EmptyAction 测试空操作类型
func TestSensitiveOperationService_AddToWhitelist_EmptyAction(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	err := service.AddToWhitelist("", "user")
	if err == nil {
		t.Error("期望返回错误, 但得到 nil")
	}
	if err.Error() != "操作类型和资源类型不能为空" {
		t.Errorf("期望错误信息为 '操作类型和资源类型不能为空', 实际为 '%v'", err)
	}
}

// TestSensitiveOperationService_AddToWhitelist_EmptyResourceType 测试空资源类型
func TestSensitiveOperationService_AddToWhitelist_EmptyResourceType(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	err := service.AddToWhitelist("delete", "")
	if err == nil {
		t.Error("期望返回错误, 但得到 nil")
	}
	if err.Error() != "操作类型和资源类型不能为空" {
		t.Errorf("期望错误信息为 '操作类型和资源类型不能为空', 实际为 '%v'", err)
	}
}

// TestSensitiveOperationService_RemoveFromWhitelist_EmptyAction 测试移除空操作类型
func TestSensitiveOperationService_RemoveFromWhitelist_EmptyAction(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	err := service.RemoveFromWhitelist("", "user")
	if err == nil {
		t.Error("期望返回错误, 但得到 nil")
	}
	if err.Error() != "操作类型和资源类型不能为空" {
		t.Errorf("期望错误信息为 '操作类型和资源类型不能为空', 实际为 '%v'", err)
	}
}

// TestSensitiveOperationService_RemoveFromWhitelist_EmptyResourceType 测试移除空资源类型
func TestSensitiveOperationService_RemoveFromWhitelist_EmptyResourceType(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	err := service.RemoveFromWhitelist("delete", "")
	if err == nil {
		t.Error("期望返回错误, 但得到 nil")
	}
	if err.Error() != "操作类型和资源类型不能为空" {
		t.Errorf("期望错误信息为 '操作类型和资源类型不能为空', 实际为 '%v'", err)
	}
}

// TestSensitiveOperationService_LogSensitiveOperation_NilRequest 测试空请求
func TestSensitiveOperationService_LogSensitiveOperation_NilRequest(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	err := service.LogSensitiveOperation(context.Background(), nil)
	if err == nil {
		t.Error("期望返回错误, 但得到 nil")
	}
	if err.Error() != "请求参数不能为空" {
		t.Errorf("期望错误信息为 '请求参数不能为空', 实际为 '%v'", err)
	}
}

// TestSensitiveOperationService_LogSensitiveOperation_NonSensitive 测试非敏感操作
func TestSensitiveOperationService_LogSensitiveOperation_NonSensitive(t *testing.T) {
	ctx := context.Background()
	var markedSensitive bool

	mockService := &MockAuditLogServiceForSensitive{
		LogFunc: func(ctx context.Context, req *LogOperationWithAuditRequest) error {
			markedSensitive = req.IsSensitive
			return nil
		},
	}

	service := NewSensitiveOperationService(mockService)

	req := &LogOperationWithAuditRequest{
		AdminID:      "admin123",
		Operation:    "query",
		ResourceType: "user",
		ResourceID:   "user123",
	}

	err := service.LogSensitiveOperation(ctx, req)
	if err != nil {
		t.Fatalf("期望记录成功, 但得到错误: %v", err)
	}
	if markedSensitive {
		t.Error("非敏感操作不应该被标记为敏感")
	}
}

// TestSensitiveOperationService_IsSensitiveOperation_CaseInsensitive 测试大小写不敏感
func TestSensitiveOperationService_IsSensitiveOperation_CaseInsensitive(t *testing.T) {
	service := NewSensitiveOperationService(nil)

	// 测试大小写变化
	if !service.IsSensitiveOperation("DELETE", "USER") {
		t.Error("大写的 DELETE:USER 应该被识别为敏感操作")
	}
	if !service.IsSensitiveOperation("Delete", "User") {
		t.Error("混合大小写的 Delete:User 应该被识别为敏感操作")
	}
}

// TestSensitiveOperationService_Whitelist_CaseInsensitive 测试白名单大小写不敏感
func TestSensitiveOperationService_Whitelist_CaseInsensitive(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	// 添加到白名单
	service.AddToWhitelist("DELETE", "USER")

	// 验证小写也不敏感
	if service.IsSensitiveOperation("delete", "user") {
		t.Error("添加到白名单后, 大小写变化都应该不敏感")
	}
}

// TestSensitiveOperationService_GetSensitiveOperations 测试获取敏感操作列表
func TestSensitiveOperationService_GetSensitiveOperations(t *testing.T) {
	operations := GetSensitiveOperations()

	if len(operations) == 0 {
		t.Error("期望敏感操作列表不为空")
	}

	// 验证包含已知敏感操作
	foundDeleteUser := false
	for _, op := range operations {
		if op == "delete:user" {
			foundDeleteUser = true
			break
		}
	}
	if !foundDeleteUser {
		t.Error("期望敏感操作列表包含 'delete:user'")
	}
}

// TestSensitiveOperationService_AddSensitiveOperation 测试动态添加敏感操作
func TestSensitiveOperationService_AddSensitiveOperation(t *testing.T) {
	service := NewSensitiveOperationService(nil)

	// 添加新的敏感操作
	AddSensitiveOperation("archive", "book")

	// 验证被识别为敏感操作
	if !service.IsSensitiveOperation("archive", "book") {
		t.Error("动态添加的敏感操作应该被识别")
	}

	// 清理
	RemoveSensitiveOperation("archive", "book")
}

// TestSensitiveOperationService_RemoveSensitiveOperation 测试动态移除敏感操作
func TestSensitiveOperationService_RemoveSensitiveOperation(t *testing.T) {
	service := NewSensitiveOperationService(nil)

	// 先添加
	AddSensitiveOperation("temp", "action")

	// 验证存在
	if !service.IsSensitiveOperation("temp", "action") {
		t.Error("添加后应该被识别为敏感操作")
	}

	// 移除
	RemoveSensitiveOperation("temp", "action")

	// 验证不再敏感
	if service.IsSensitiveOperation("temp", "action") {
		t.Error("移除后不应该被识别为敏感操作")
	}
}

// TestSensitiveOperationService_ConcurrentWhitelistOperations 测试并发白名单操作
func TestSensitiveOperationService_ConcurrentWhitelistOperations(t *testing.T) {
	mockService := &MockAuditLogServiceForSensitive{}
	service := NewSensitiveOperationService(mockService)

	// 并发添加和移除
	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			service.AddToWhitelist("action", "resource")
		}
		done <- true
	}()
	go func() {
		for i := 0; i < 100; i++ {
			service.RemoveFromWhitelist("action", "resource")
		}
		done <- true
	}()

	<-done
	<-done

	// 验证服务仍然可用
	if service.IsSensitiveOperation("delete", "user") {
		// 这是预期的, delete:user 是敏感操作
	}
}
