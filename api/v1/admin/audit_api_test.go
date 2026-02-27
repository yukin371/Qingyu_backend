package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	adminModel "Qingyu_backend/models/users"
	adminService "Qingyu_backend/service/admin"
)

// MockAuditLogServiceForAPI Mock审计日志服务（用于API测试）
type MockAuditLogServiceForAPI struct {
	QueryFunc        func(req interface{}) (interface{}, int64, error)
	GetByResourceFunc func(resourceType, resourceID string) (interface{}, error)
}

func (m *MockAuditLogServiceForAPI) LogOperationWithAudit(ctx context.Context, req *adminService.LogOperationWithAuditRequest) error {
	return nil
}

func (m *MockAuditLogServiceForAPI) QueryAuditLogs(ctx context.Context, req *adminService.QueryAuditLogsRequest) ([]*adminModel.AdminLog, int64, error) {
	if m.QueryFunc != nil {
		result, total, err := m.QueryFunc(req)
		return result.([]*adminModel.AdminLog), total, err
	}
	return []*adminModel.AdminLog{}, 0, nil
}

func (m *MockAuditLogServiceForAPI) GetLogsByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error) {
	if m.GetByResourceFunc != nil {
		result, _ := m.GetByResourceFunc(resourceType, resourceID)
		return result.([]*adminModel.AdminLog), nil
	}
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogServiceForAPI) GetLogsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error) {
	return []*adminModel.AdminLog{}, nil
}

func (m *MockAuditLogServiceForAPI) CleanOldLogs(ctx context.Context, beforeDate time.Time) error {
	return nil
}

// 设置 Gin 为测试模式
func init() {
	gin.SetMode(gin.TestMode)
}

// setupAuditTestRouter 设置审计API测试路由
func setupAuditTestRouter(auditAPI *AuditAPI) *gin.Engine {
	router := gin.New()
	apiGroup := router.Group("/api/v1/admin")
	{
		apiGroup.GET("/audit/trail", auditAPI.GetAuditTrail)
		apiGroup.GET("/audit/trail/resource/:type/:id", auditAPI.GetResourceAuditTrail)
		apiGroup.GET("/audit/trail/export", auditAPI.ExportAuditTrail)
		apiGroup.GET("/audit/statistics", auditAPI.GetAuditStatistics)
	}
	return router
}

// TestAuditAPI_GetAuditTrail_Success 测试获取审计追踪
func TestAuditAPI_GetAuditTrail_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		QueryFunc: func(req interface{}) (interface{}, int64, error) {
			return []*adminModel.AdminLog{
				{
					ID:           "log1",
					AdminID:      "admin123",
					AdminName:    "管理员A",
					Operation:    "ban_user",
					ResourceType: "user",
					ResourceID:   "user456",
					CreatedAt:    time.Now(),
				},
			}, 1, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/trail?page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 200 {
		t.Errorf("期望 code 为 200, 实际为 %v", response["code"])
	}
}

// TestAuditAPI_GetAuditTrailByResource_Success 测试按资源查询
func TestAuditAPI_GetAuditTrailByResource_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		GetByResourceFunc: func(resourceType, resourceID string) (interface{}, error) {
			return []*adminModel.AdminLog{
				{
					ID:           "log1",
					AdminID:      "admin123",
					Operation:    "update_role",
					ResourceType: resourceType,
					ResourceID:   resourceID,
				},
			}, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/trail/resource/role/role1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}
}

// TestAuditAPI_GetAuditTrailByAction_Success 测试按操作查询
func TestAuditAPI_GetAuditTrailByAction_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		QueryFunc: func(req interface{}) (interface{}, int64, error) {
			return []*adminModel.AdminLog{
				{
					ID:        "log1",
					AdminID:   "admin123",
					Operation: "delete_user",
				},
			}, 1, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/trail?operation=delete_user&page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}
}

// TestAuditAPI_GetAuditTrailWithPagination_Success 测试分页查询
func TestAuditAPI_GetAuditTrailWithPagination_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		QueryFunc: func(req interface{}) (interface{}, int64, error) {
			// 返回 25 条数据
			logs := make([]*adminModel.AdminLog, 25)
			for i := 0; i < 25; i++ {
				logs[i] = &adminModel.AdminLog{
					ID:         "log" + string(rune('0'+i)),
					AdminID:    "admin123",
					Operation:  "ban_user",
					CreatedAt:  time.Now(),
				}
			}
			return logs, 25, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	// 请求第 2 页，每页 10 条
	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/trail?page=2&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	total := data["total"].(float64)
	if total != 25 {
		t.Errorf("期望总数为 25, 实际为 %v", total)
	}
}

// TestAuditAPI_ExportAuditTrail_Success 测试导出审计日志
func TestAuditAPI_ExportAuditTrail_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		QueryFunc: func(req interface{}) (interface{}, int64, error) {
			return []*adminModel.AdminLog{
				{
					ID:           "log1",
					AdminID:      "admin123",
					AdminName:    "管理员A",
					Operation:    "ban_user",
					ResourceType: "user",
					ResourceID:   "user456",
					CreatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			}, 1, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/trail/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("期望 Content-Type 为 text/csv, 实际为 %s", contentType)
	}
}

// TestAuditAPI_GetAuditStatistics_Success 测试获取审计统计
func TestAuditAPI_GetAuditStatistics_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		QueryFunc: func(req interface{}) (interface{}, int64, error) {
			return []*adminModel.AdminLog{}, 100, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 200 {
		t.Errorf("期望 code 为 200, 实际为 %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["total_logs"].(float64) != 100 {
		t.Errorf("期望 total_logs 为 100, 实际为 %v", data["total_logs"])
	}
}

// TestAuditAPI_GetAuditTrailWithDateRange_Success 测试按日期范围查询
func TestAuditAPI_GetAuditTrailWithDateRange_Success(t *testing.T) {
	mockService := &MockAuditLogServiceForAPI{
		QueryFunc: func(req interface{}) (interface{}, int64, error) {
			return []*adminModel.AdminLog{
				{
					ID:         "log1",
					AdminID:    "admin123",
					Operation:  "ban_user",
					CreatedAt:  time.Now(),
				},
			}, 1, nil
		},
	}

	auditAPI := &AuditAPI{
		auditLogService: mockService,
	}
	router := setupAuditTestRouter(auditAPI)

	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/trail?start_date="+startDate+"&end_date="+endDate+"&page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际为 %d", w.Code)
	}
}
