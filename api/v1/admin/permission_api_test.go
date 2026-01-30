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

	authModel "Qingyu_backend/models/auth"
	sharedService "Qingyu_backend/service/shared"
)

// MockPermissionService 模拟PermissionService
type MockPermissionService struct {
	mock.Mock
}

func (m *MockPermissionService) GetAllPermissions(ctx context.Context) ([]*authModel.Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Permission), args.Error(1)
}

func (m *MockPermissionService) GetPermissionByCode(ctx context.Context, code string) (*authModel.Permission, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Permission), args.Error(1)
}

func (m *MockPermissionService) CreatePermission(ctx context.Context, permission *authModel.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionService) UpdatePermission(ctx context.Context, permission *authModel.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionService) DeletePermission(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *MockPermissionService) GetAllRoles(ctx context.Context) ([]*authModel.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

func (m *MockPermissionService) GetRoleByID(ctx context.Context, roleID string) (*authModel.Role, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

func (m *MockPermissionService) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

func (m *MockPermissionService) CreateRole(ctx context.Context, role *authModel.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockPermissionService) UpdateRole(ctx context.Context, role *authModel.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockPermissionService) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockPermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error {
	args := m.Called(ctx, roleID, permissionCode)
	return args.Error(0)
}

func (m *MockPermissionService) RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error {
	args := m.Called(ctx, roleID, permissionCode)
	return args.Error(0)
}

func (m *MockPermissionService) GetRolePermissions(ctx context.Context, roleID string) ([]*authModel.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Permission), args.Error(1)
}

func (m *MockPermissionService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionService) AssignRoleToUser(ctx context.Context, userID, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

func (m *MockPermissionService) RemoveRoleFromUser(ctx context.Context, userID, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

func (m *MockPermissionService) GetUserPermissions(ctx context.Context, userID string) ([]*authModel.Permission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Permission), args.Error(1)
}

func (m *MockPermissionService) ClearUserRoles(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPermissionService) UserHasPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	args := m.Called(ctx, userID, permissionCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionService) UserHasAnyPermission(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	args := m.Called(ctx, userID, permissionCodes)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionService) UserHasAllPermissions(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	args := m.Called(ctx, userID, permissionCodes)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionService) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	args := m.Called(ctx, userID, resource, action)
	return args.Bool(0), args.Error(1)
}

// setupPermissionAPITestRouter 设置权限管理测试路由
func setupPermissionAPITestRouter(permissionService *MockPermissionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewPermissionAPI(permissionService)

	v1 := r.Group("/api/v1/admin")
	{
		// Permissions
		v1.GET("/permissions", api.GetAllPermissions)
		v1.GET("/permissions/:code", api.GetPermission)
		v1.POST("/permissions", api.CreatePermission)
		v1.PUT("/permissions/:code", api.UpdatePermission)
		v1.DELETE("/permissions/:code", api.DeletePermission)

		// Roles
		v1.GET("/roles", api.GetAllRoles)
		v1.GET("/roles/:id", api.GetRole)
		v1.POST("/roles", api.CreateRole)
		v1.PUT("/roles/:id", api.UpdateRole)
		v1.DELETE("/roles/:id", api.DeleteRole)
		v1.GET("/roles/:id/permissions", api.GetRolePermissions)
		v1.POST("/roles/:id/permissions/:permissionCode", api.AssignPermissionToRole)
		v1.DELETE("/roles/:id/permissions/:permissionCode", api.RemovePermissionFromRole)

		// User roles
		v1.GET("/users/:userId/roles", api.GetUserRoles)
		v1.POST("/users/:userId/roles", api.AssignRoleToUser)
		v1.DELETE("/users/:userId/roles", api.RemoveRoleFromUser)
		v1.GET("/users/:userId/permissions", api.GetUserPermissions)
	}

	return r
}

// ==================== GetAllPermissions Tests ====================

func TestPermissionAPI_GetAllPermissions_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	expectedPermissions := []*authModel.Permission{
		{Code: "read:books", Name: "Read Books", Description: "Can read books"},
		{Code: "write:books", Name: "Write Books", Description: "Can write books"},
	}

	mockService.On("GetAllPermissions", mock.Anything).Return(expectedPermissions, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetAllPermissions_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	mockService.On("GetAllPermissions", mock.Anything).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetPermission Tests ====================

func TestPermissionAPI_GetPermission_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	code := "read:books"
	expectedPermission := &authModel.Permission{Code: code, Name: "Read Books"}

	mockService.On("GetPermissionByCode", mock.Anything, code).Return(expectedPermission, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/permissions/"+code, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetPermission_EmptyCode(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/permissions/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Gin redirects trailing slash, so we expect 301
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestPermissionAPI_GetPermission_NotFound(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	code := "invalid:permission"
	mockService.On("GetPermissionByCode", mock.Anything, code).Return(nil, sharedService.ErrPermissionNotFound)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/permissions/"+code, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== CreatePermission Tests ====================

func TestPermissionAPI_CreatePermission_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	req := authModel.Permission{Code: "new:permission", Name: "New Permission", Description: "Test"}
	mockService.On("CreatePermission", mock.Anything, &req).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/permissions", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_CreatePermission_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/permissions", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== UpdatePermission Tests ====================

func TestPermissionAPI_UpdatePermission_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	code := "update:test"
	req := authModel.Permission{Code: code, Name: "Updated Name"}
	mockService.On("UpdatePermission", mock.Anything, &req).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/permissions/"+code, bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_UpdatePermission_EmptyCode(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	req := authModel.Permission{Name: "Test"}

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/permissions//", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	// Gin redirects trailing slash, so we expect 307 for PUT
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

// ==================== DeletePermission Tests ====================

func TestPermissionAPI_DeletePermission_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	code := "delete:test"
	mockService.On("DeletePermission", mock.Anything, code).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/permissions/"+code, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_DeletePermission_EmptyCode(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/permissions/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// DELETE with trailing slash returns 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ==================== GetAllRoles Tests ====================

func TestPermissionAPI_GetAllRoles_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	expectedRoles := []*authModel.Role{
		{ID: "role1", Name: "Admin", Permissions: []string{"read", "write"}},
		{ID: "role2", Name: "User", Permissions: []string{"read"}},
	}

	mockService.On("GetAllRoles", mock.Anything).Return(expectedRoles, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetAllRoles_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	mockService.On("GetAllRoles", mock.Anything).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetRole Tests ====================

func TestPermissionAPI_GetRole_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "role123"
	expectedRole := &authModel.Role{ID: roleID, Name: "Editor"}

	mockService.On("GetRoleByID", mock.Anything, roleID).Return(expectedRole, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles/"+roleID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetRole_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Gin redirects trailing slash, so we expect 301
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestPermissionAPI_GetRole_NotFound(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "invalid"
	mockService.On("GetRoleByID", mock.Anything, roleID).Return(nil, sharedService.ErrRoleNotFound)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles/"+roleID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== CreateRole Tests ====================

func TestPermissionAPI_CreateRole_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	req := authModel.Role{Name: "Moderator", Description: "Content moderator"}
	mockService.On("CreateRole", mock.Anything, &req).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== UpdateRole Tests ====================

func TestPermissionAPI_UpdateRole_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "role123"
	req := authModel.Role{ID: roleID, Name: "Updated Moderator"}
	mockService.On("UpdateRole", mock.Anything, &req).Return(nil)

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/roles/"+roleID, bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_UpdateRole_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	req := authModel.Role{Name: "Test"}

	// When
	jsonBody, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", "/api/v1/admin/roles//", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	// Gin redirects trailing slash, so we expect 307 for PUT
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

// ==================== DeleteRole Tests ====================

func TestPermissionAPI_DeleteRole_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "role123"
	mockService.On("DeleteRole", mock.Anything, roleID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/roles/"+roleID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== AssignPermissionToRole Tests ====================

func TestPermissionAPI_AssignPermissionToRole_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "role123"
	permissionCode := "write:books"
	mockService.On("AssignPermissionToRole", mock.Anything, roleID, permissionCode).Return(nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/roles/"+roleID+"/permissions/"+permissionCode, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_AssignPermissionToRole_EmptyParams(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/roles//permissions/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// POST with trailing slash returns 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ==================== RemovePermissionFromRole Tests ====================

func TestPermissionAPI_RemovePermissionFromRole_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "role123"
	permissionCode := "write:books"
	mockService.On("RemovePermissionFromRole", mock.Anything, roleID, permissionCode).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/roles/"+roleID+"/permissions/"+permissionCode, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetRolePermissions Tests ====================

func TestPermissionAPI_GetRolePermissions_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	roleID := "role123"
	expectedPermissions := []*authModel.Permission{
		{Code: "read:books", Name: "Read Books"},
		{Code: "write:books", Name: "Write Books"},
	}

	mockService.On("GetRolePermissions", mock.Anything, roleID).Return(expectedPermissions, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles/"+roleID+"/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetRolePermissions_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/roles//permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Empty path param returns 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== GetUserRoles Tests ====================

func TestPermissionAPI_GetUserRoles_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	userID := "user123"
	expectedRoles := []string{"reader", "author"}
	mockService.On("GetUserRoles", mock.Anything, userID).Return(expectedRoles, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID+"/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetUserRoles_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users//roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Empty path param returns 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== AssignRoleToUser Tests ====================

func TestPermissionAPI_AssignRoleToUser_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	userID := "user123"
	reqBody := AssignRoleRequest{Role: "admin"}
	mockService.On("AssignRoleToUser", mock.Anything, userID, "admin").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/users/"+userID+"/roles", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_AssignRoleToUser_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	reqBody := AssignRoleRequest{Role: "admin"}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/users//roles", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	// Empty path param returns 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPermissionAPI_AssignRoleToUser_MissingRequiredFields(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	userID := "user123"
	reqBody := map[string]string{} // Missing role

	// When
	jsonBody, _ := json.Marshal(reqBody)
	reqHTTP, _ := http.NewRequest("POST", "/api/v1/admin/users/"+userID+"/roles", bytes.NewBuffer(jsonBody))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== RemoveRoleFromUser Tests ====================

func TestPermissionAPI_RemoveRoleFromUser_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	userID := "user123"
	role := "admin"
	mockService.On("RemoveRoleFromUser", mock.Anything, userID, role).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/"+userID+"/roles?role="+role, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_RemoveRoleFromUser_EmptyParams(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users//roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Empty path param returns 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== GetUserPermissions Tests ====================

func TestPermissionAPI_GetUserPermissions_Success(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	userID := "user123"
	expectedPermissions := []*authModel.Permission{
		{Code: "read:books", Name: "Read Books"},
		{Code: "write:books", Name: "Write Books"},
	}
	mockService.On("GetUserPermissions", mock.Anything, userID).Return(expectedPermissions, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID+"/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestPermissionAPI_GetUserPermissions_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockPermissionService)
	router := setupPermissionAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users//permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Empty path param returns 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
