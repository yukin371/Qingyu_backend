package admin

import (
	adminModel "Qingyu_backend/models/users"
	"context"
	"errors"
	"testing"
	"time"
)

// ============ Mock Repositories ============

type MockAuditRepository struct {
	createFunc        func(ctx context.Context, record *AuditRecord) error
	getFunc           func(ctx context.Context, recordID string) (*AuditRecord, error)
	updateFunc        func(ctx context.Context, recordID string, updates map[string]interface{}) error
	listByStatusFunc  func(ctx context.Context, contentType, status string) ([]*AuditRecord, error)
	listByContentFunc func(ctx context.Context, contentID string) ([]*AuditRecord, error)
}

func (m *MockAuditRepository) Create(ctx context.Context, record *AuditRecord) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, record)
	}
	return nil
}

func (m *MockAuditRepository) Get(ctx context.Context, recordID string) (*AuditRecord, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, recordID)
	}
	return nil, nil
}

func (m *MockAuditRepository) Update(ctx context.Context, recordID string, updates map[string]interface{}) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, recordID, updates)
	}
	return nil
}

func (m *MockAuditRepository) ListByStatus(ctx context.Context, contentType, status string) ([]*AuditRecord, error) {
	if m.listByStatusFunc != nil {
		return m.listByStatusFunc(ctx, contentType, status)
	}
	return []*AuditRecord{}, nil
}

func (m *MockAuditRepository) ListByContent(ctx context.Context, contentID string) ([]*AuditRecord, error) {
	if m.listByContentFunc != nil {
		return m.listByContentFunc(ctx, contentID)
	}
	return []*AuditRecord{}, nil
}

type MockLogRepository struct {
	createFunc func(ctx context.Context, log *AdminLog) error
	listFunc   func(ctx context.Context, filter *LogFilter) ([]*AdminLog, error)
}

func (m *MockLogRepository) Create(ctx context.Context, log *AdminLog) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, log)
	}
	return nil
}

func (m *MockLogRepository) List(ctx context.Context, filter *LogFilter) ([]*AdminLog, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return []*AdminLog{}, nil
}

type MockUserRepository struct {
	getStatisticsFunc func(ctx context.Context, userID string) (*UserStatistics, error)
	banUserFunc       func(ctx context.Context, userID, reason string, until time.Time) error
	unbanUserFunc     func(ctx context.Context, userID string) error
}

func (m *MockUserRepository) GetStatistics(ctx context.Context, userID string) (*UserStatistics, error) {
	if m.getStatisticsFunc != nil {
		return m.getStatisticsFunc(ctx, userID)
	}
	return &UserStatistics{}, nil
}

func (m *MockUserRepository) BanUser(ctx context.Context, userID, reason string, until time.Time) error {
	if m.banUserFunc != nil {
		return m.banUserFunc(ctx, userID, reason, until)
	}
	return nil
}

func (m *MockUserRepository) UnbanUser(ctx context.Context, userID string) error {
	if m.unbanUserFunc != nil {
		return m.unbanUserFunc(ctx, userID)
	}
	return nil
}

// ============ 测试用例 ============

func TestReviewContent_Approve(t *testing.T) {
	mockAudit := &MockAuditRepository{
		listByContentFunc: func(ctx context.Context, contentID string) ([]*AuditRecord, error) {
			return []*AuditRecord{
				{
					ID:        "audit1",
					ContentID: contentID,
					Status:    adminModel.AuditStatusPending,
				},
			}, nil
		},
	}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &ReviewContentRequest{
		ContentID:   "content123",
		ContentType: "book",
		Action:      "approve",
		ReviewerID:  "admin1",
	}

	err := service.ReviewContent(context.Background(), req)
	if err != nil {
		t.Errorf("审核内容失败: %v", err)
	}
}

func TestReviewContent_Reject(t *testing.T) {
	mockAudit := &MockAuditRepository{
		listByContentFunc: func(ctx context.Context, contentID string) ([]*AuditRecord, error) {
			return []*AuditRecord{
				{
					ID:        "audit1",
					ContentID: contentID,
					Status:    adminModel.AuditStatusPending,
				},
			}, nil
		},
	}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &ReviewContentRequest{
		ContentID:   "content123",
		ContentType: "book",
		Action:      "reject",
		Reason:      "违规内容",
		ReviewerID:  "admin1",
	}

	err := service.ReviewContent(context.Background(), req)
	if err != nil {
		t.Errorf("审核内容失败: %v", err)
	}
}

func TestReviewContent_InvalidAction(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &ReviewContentRequest{
		ContentID:   "content123",
		ContentType: "book",
		Action:      "invalid",
		ReviewerID:  "admin1",
	}

	err := service.ReviewContent(context.Background(), req)
	if err == nil {
		t.Error("期望无效操作错误，但成功了")
	}
}

func TestGetPendingReviews(t *testing.T) {
	mockAudit := &MockAuditRepository{
		listByStatusFunc: func(ctx context.Context, contentType, status string) ([]*AuditRecord, error) {
			if status != adminModel.AuditStatusPending {
				t.Errorf("期望查询pending状态，实际: %s", status)
			}
			return []*AuditRecord{
				{ID: "audit1", ContentType: contentType},
				{ID: "audit2", ContentType: contentType},
			}, nil
		},
	}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	records, err := service.GetPendingReviews(context.Background(), "book")
	if err != nil {
		t.Errorf("获取待审核内容失败: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("期望2条记录，实际 %d 条", len(records))
	}
}

func TestBanUser(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{
		banUserFunc: func(ctx context.Context, userID, reason string, until time.Time) error {
			if userID != "user123" {
				t.Errorf("用户ID错误: %s", userID)
			}
			return nil
		},
	}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	err := service.BanUser(context.Background(), "user123", "违规", 24*time.Hour)
	if err != nil {
		t.Errorf("封禁用户失败: %v", err)
	}
}

func TestBanUser_Permanent(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{
		banUserFunc: func(ctx context.Context, userID, reason string, until time.Time) error {
			// 永久封禁应该设置到2099年
			if until.Year() != 2099 {
				t.Errorf("永久封禁时间错误: %v", until)
			}
			return nil
		},
	}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	err := service.BanUser(context.Background(), "user123", "严重违规", 0)
	if err != nil {
		t.Errorf("永久封禁用户失败: %v", err)
	}
}

func TestUnbanUser(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	err := service.UnbanUser(context.Background(), "user123")
	if err != nil {
		t.Errorf("解封用户失败: %v", err)
	}
}

func TestManageUser_Ban(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &ManageUserRequest{
		UserID:   "user123",
		Action:   "ban",
		Reason:   "违规",
		Duration: 86400, // 1天
		AdminID:  "admin1",
	}

	err := service.ManageUser(context.Background(), req)
	if err != nil {
		t.Errorf("管理用户失败: %v", err)
	}
}

func TestManageUser_InvalidAction(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &ManageUserRequest{
		UserID:  "user123",
		Action:  "invalid",
		AdminID: "admin1",
	}

	err := service.ManageUser(context.Background(), req)
	if err == nil {
		t.Error("期望无效操作错误，但成功了")
	}
}

func TestGetUserStatistics(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{
		getStatisticsFunc: func(ctx context.Context, userID string) (*UserStatistics, error) {
			return &UserStatistics{
				UserID:     userID,
				TotalBooks: 10,
				TotalWords: 100000,
			}, nil
		},
	}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	stats, err := service.GetUserStatistics(context.Background(), "user123")
	if err != nil {
		t.Errorf("获取用户统计失败: %v", err)
	}

	if stats.TotalBooks != 10 {
		t.Errorf("统计数据错误: %v", stats)
	}
}

func TestReviewWithdraw_Approve(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	err := service.ReviewWithdraw(context.Background(), "withdraw123", "admin1", true, "")
	if err != nil {
		t.Errorf("审核提现失败: %v", err)
	}
}

func TestReviewWithdraw_Reject(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	err := service.ReviewWithdraw(context.Background(), "withdraw123", "admin1", false, "信息不完整")
	if err != nil {
		t.Errorf("审核提现失败: %v", err)
	}
}

func TestLogOperation(t *testing.T) {
	logCreated := false
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{
		createFunc: func(ctx context.Context, log *AdminLog) error {
			logCreated = true
			if log.AdminID != "admin1" {
				t.Errorf("管理员ID错误: %s", log.AdminID)
			}
			return nil
		},
	}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &LogOperationRequest{
		AdminID:   "admin1",
		Operation: "test_operation",
		Target:    "target123",
		IP:        "127.0.0.1",
	}

	err := service.LogOperation(context.Background(), req)
	if err != nil {
		t.Errorf("记录操作日志失败: %v", err)
	}

	if !logCreated {
		t.Error("日志未创建")
	}
}

func TestGetOperationLogs(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{
		listFunc: func(ctx context.Context, filter *LogFilter) ([]*AdminLog, error) {
			return []*AdminLog{
				{ID: "log1", AdminID: "admin1"},
				{ID: "log2", AdminID: "admin1"},
			}, nil
		},
	}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &GetLogsRequest{
		AdminID:  "admin1",
		Page:     1,
		PageSize: 50,
	}

	logs, err := service.GetOperationLogs(context.Background(), req)
	if err != nil {
		t.Errorf("获取操作日志失败: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("期望2条日志，实际 %d 条", len(logs))
	}
}

func TestGetOperationLogs_DefaultPagination(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{
		listFunc: func(ctx context.Context, filter *LogFilter) ([]*AdminLog, error) {
			if filter.Page != 1 || filter.PageSize != 50 {
				t.Errorf("默认分页参数错误: page=%d, pageSize=%d", filter.Page, filter.PageSize)
			}
			return []*AdminLog{}, nil
		},
	}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &GetLogsRequest{}

	_, err := service.GetOperationLogs(context.Background(), req)
	if err != nil {
		t.Errorf("获取操作日志失败: %v", err)
	}
}

func TestExportLogs(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{
		listFunc: func(ctx context.Context, filter *LogFilter) ([]*AdminLog, error) {
			return []*AdminLog{
				{
					ID:        "log1",
					AdminID:   "admin1",
					Operation: "test",
					Target:    "target1",
					IP:        "127.0.0.1",
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	csv, err := service.ExportLogs(context.Background(), time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		t.Errorf("导出日志失败: %v", err)
	}

	if len(csv) == 0 {
		t.Error("导出的CSV为空")
	}

	// 检查CSV包含表头
	csvStr := string(csv)
	if !contains(csvStr, "管理员ID") {
		t.Error("CSV缺少表头")
	}
}

func TestExportLogs_EmptyResult(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{
		listFunc: func(ctx context.Context, filter *LogFilter) ([]*AdminLog, error) {
			return []*AdminLog{}, nil
		},
	}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	csv, err := service.ExportLogs(context.Background(), time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		t.Errorf("导出日志失败: %v", err)
	}

	// 即使没有数据，也应该有表头
	if len(csv) == 0 {
		t.Error("CSV应该至少包含表头")
	}
}

func TestHealth(t *testing.T) {
	mockAudit := &MockAuditRepository{}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	// 初始化服务（使用类型断言访问实现类的Initialize方法）
	if serviceImpl, ok := service.(*AdminServiceImpl); ok {
		err := serviceImpl.Initialize(context.Background())
		if err != nil {
			t.Fatalf("初始化服务失败: %v", err)
		}
	}

	err := service.Health(context.Background())
	if err != nil {
		t.Errorf("健康检查失败: %v", err)
	}
}

func TestReviewContent_CreateRecordFailure(t *testing.T) {
	mockAudit := &MockAuditRepository{
		listByContentFunc: func(ctx context.Context, contentID string) ([]*AuditRecord, error) {
			return []*AuditRecord{}, nil
		},
		createFunc: func(ctx context.Context, record *AuditRecord) error {
			return errors.New("创建失败")
		},
	}
	mockLog := &MockLogRepository{}
	mockUser := &MockUserRepository{}

	service := NewAdminService(mockAudit, mockLog, mockUser)

	req := &ReviewContentRequest{
		ContentID:   "content123",
		ContentType: "book",
		Action:      "approve",
		ReviewerID:  "admin1",
	}

	err := service.ReviewContent(context.Background(), req)
	if err == nil {
		t.Error("期望创建记录失败，但成功了")
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
