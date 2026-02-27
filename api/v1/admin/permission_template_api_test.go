package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authService "Qingyu_backend/service/auth"
)

// 设置测试模式
func init() {
	gin.SetMode(gin.TestMode)
}

// MockPermissionTemplateServiceForAPI 模拟权限模板服务（用于API测试）
type MockPermissionTemplateServiceForAPI struct {
	mock.Mock
}

func (m *MockPermissionTemplateServiceForAPI) CreateTemplate(ctx context.Context, req *authService.CreateTemplateRequest) (*authService.TemplateResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.TemplateResponse), args.Error(1)
}

func (m *MockPermissionTemplateServiceForAPI) GetTemplate(ctx context.Context, templateID string) (*authService.TemplateResponse, error) {
	args := m.Called(ctx, templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.TemplateResponse), args.Error(1)
}

func (m *MockPermissionTemplateServiceForAPI) GetTemplateByCode(ctx context.Context, code string) (*authService.TemplateResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.TemplateResponse), args.Error(1)
}

func (m *MockPermissionTemplateServiceForAPI) UpdateTemplate(ctx context.Context, templateID string, req *authService.UpdateTemplateRequest) error {
	args := m.Called(ctx, templateID, req)
	return args.Error(0)
}

func (m *MockPermissionTemplateServiceForAPI) DeleteTemplate(ctx context.Context, templateID string) error {
	args := m.Called(ctx, templateID)
	return args.Error(0)
}

func (m *MockPermissionTemplateServiceForAPI) ListTemplates(ctx context.Context, category string) ([]*authService.TemplateResponse, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authService.TemplateResponse), args.Error(1)
}

func (m *MockPermissionTemplateServiceForAPI) ApplyTemplate(ctx context.Context, templateID, roleID string) error {
	args := m.Called(ctx, templateID, roleID)
	return args.Error(0)
}

func (m *MockPermissionTemplateServiceForAPI) InitializeSystemTemplates(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// setupTestRouter 设置测试路由
func setupTestRouter(service *MockPermissionTemplateServiceForAPI) *gin.Engine {
	router := gin.New()
	api := NewPermissionTemplateAPI(service)
	v1 := router.Group("/api/v1/admin")
	api.RegisterRoutes(v1)
	return router
}

// TestPermissionTemplateAPI_CreateTemplate_Success 测试成功创建模板
func TestPermissionTemplateAPI_CreateTemplate_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	req := CreateTemplateRequest{
		Name:        "测试模板",
		Code:        "test_template",
		Description: "测试用模板",
		Permissions: []string{"user.read", "book.read"},
		Category:    "custom",
	}

	expectedResp := &authService.TemplateResponse{
		ID:          "template_id",
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Permissions: req.Permissions,
		IsSystem:    false,
		Category:    req.Category,
	}

	mockService.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(r *authService.CreateTemplateRequest) bool {
		return r.Name == req.Name && r.Code == req.Code
	})).Return(expectedResp, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/admin/permission-templates", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"].(bool))

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_CreateTemplate_ValidationError 测试参数验证错误
func TestPermissionTemplateAPI_CreateTemplate_ValidationError(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	req := map[string]interface{}{
		"name": "", // 空名称应该失败
		"code": "test",
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/admin/permission-templates", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPermissionTemplateAPI_GetTemplate_Success 测试获取模板
func TestPermissionTemplateAPI_GetTemplate_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	expectedResp := &authService.TemplateResponse{
		ID:          "template_id",
		Name:        "测试模板",
		Code:        "test_code",
		Description: "测试描述",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    "custom",
	}

	mockService.On("GetTemplate", mock.Anything, "template_id").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/admin/permission-templates/template_id", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"].(bool))

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_GetTemplate_MissingID 测试缺少ID
func TestPermissionTemplateAPI_GetTemplate_MissingID(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/admin/permission-templates/", nil)

	router.ServeHTTP(w, httpReq)

	// Gin路由重定向，应该返回301或者404
	// 因为路径不以/结尾，Gin会重定向到带/的路径
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMovedPermanently)
}

// TestPermissionTemplateAPI_UpdateTemplate_Success 测试更新模板
func TestPermissionTemplateAPI_UpdateTemplate_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	req := UpdateTemplateRequest{
		Name:        "新名称",
		Description: "新描述",
	}

	mockService.On("UpdateTemplate", mock.Anything, "template_id", mock.MatchedBy(func(r *authService.UpdateTemplateRequest) bool {
		return r.Name == req.Name && r.Description == req.Description
	})).Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/v1/admin/permission-templates/template_id", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_DeleteTemplate_Success 测试删除模板
func TestPermissionTemplateAPI_DeleteTemplate_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	mockService.On("DeleteTemplate", mock.Anything, "template_id").Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/v1/admin/permission-templates/template_id", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_ListTemplates_All 测试列出所有模板
func TestPermissionTemplateAPI_ListTemplates_All(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	expectedTemplates := []*authService.TemplateResponse{
		{
			ID:          "1",
			Name:        "模板1",
			Code:        "template_1",
			Description: "第一个模板",
			Permissions: []string{"user.read"},
			IsSystem:    false,
			Category:    "custom",
		},
		{
			ID:          "2",
			Name:        "模板2",
			Code:        "template_2",
			Description: "第二个模板",
			Permissions: []string{"book.read"},
			IsSystem:    false,
			Category:    "reader",
		},
	}

	mockService.On("ListTemplates", mock.Anything, "").Return(expectedTemplates, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/admin/permission-templates", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// 检查success字段
	if success, ok := resp["success"]; ok {
		assert.True(t, success.(bool))
	}

	// 检查data字段
	if data, ok := resp["data"].(map[string]interface{}); ok {
		assert.Equal(t, float64(2), data["total"])
	}

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_ListTemplates_WithCategory 测试按分类列出模板
func TestPermissionTemplateAPI_ListTemplates_WithCategory(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	expectedTemplates := []*authService.TemplateResponse{
		{
			ID:          "1",
			Name:        "读者模板",
			Code:        "reader_template",
			Description: "读者模板",
			Permissions: []string{"book.read"},
			IsSystem:    false,
			Category:    "reader",
		},
	}

	mockService.On("ListTemplates", mock.Anything, "reader").Return(expectedTemplates, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/admin/permission-templates?category=reader", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// 检查success字段
	if success, ok := resp["success"]; ok {
		assert.True(t, success.(bool))
	}

	// 检查data字段
	if data, ok := resp["data"].(map[string]interface{}); ok {
		assert.Equal(t, float64(1), data["total"])
	}

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_ApplyTemplate_Success 测试应用模板
func TestPermissionTemplateAPI_ApplyTemplate_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	req := struct {
		RoleID string `json:"roleId"`
	}{
		RoleID: "role_id",
	}

	mockService.On("ApplyTemplate", mock.Anything, "template_id", "role_id").Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/admin/permission-templates/template_id/apply", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_ApplyTemplate_MissingRoleID 测试缺少角色ID
func TestPermissionTemplateAPI_ApplyTemplate_MissingRoleID(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	req := struct {
		RoleID string `json:"roleId"`
	}{
		RoleID: "", // 空角色ID
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/admin/permission-templates/template_id/apply", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPermissionTemplateAPI_InitializeSystemTemplates_Success 测试初始化系统模板
func TestPermissionTemplateAPI_InitializeSystemTemplates_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	mockService.On("InitializeSystemTemplates", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/admin/permission-templates/initialize", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_GetTemplateByCode_Success 测试根据代码获取模板
func TestPermissionTemplateAPI_GetTemplateByCode_Success(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	expectedResp := &authService.TemplateResponse{
		ID:          "template_id",
		Name:        "测试模板",
		Code:        "test_code",
		Description: "测试描述",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    "custom",
	}

	mockService.On("GetTemplateByCode", mock.Anything, "test_code").Return(expectedResp, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/admin/permission-templates/code/test_code", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"].(bool))

	mockService.AssertExpectations(t)
}

// TestPermissionTemplateAPI_GetTemplateByCode_MissingCode 测试缺少代码
func TestPermissionTemplateAPI_GetTemplateByCode_MissingCode(t *testing.T) {
	mockService := new(MockPermissionTemplateServiceForAPI)
	router := setupTestRouter(mockService)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/admin/permission-templates/code/", nil)

	router.ServeHTTP(w, httpReq)

	// Gin路由重定向，应该返回301或者404
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMovedPermanently)
}
