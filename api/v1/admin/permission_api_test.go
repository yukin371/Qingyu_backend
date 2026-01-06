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

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/auth"
)

// MockPermissionServiceForAPI 模拟权限服务（用于API测试）
type MockPermissionServiceForAPI struct {
	permissions map[string]*auth.Permission
	roles       map[string]*auth.Role
	userRoles   map[string][]string // userID -> roleNames
}

func NewMockPermissionServiceForAPI() *MockPermissionServiceForAPI {
	return &MockPermissionServiceForAPI{
		permissions: make(map[string]*auth.Permission),
		roles:       make(map[string]*auth.Role),
		userRoles:   make(map[string][]string),
	}
}

// ==================== 权限管理 ====================

func (m *MockPermissionServiceForAPI) GetAllPermissions(ctx context.Context) ([]*auth.Permission, error) {
	result := make([]*auth.Permission, 0, len(m.permissions))
	for _, p := range m.permissions {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockPermissionServiceForAPI) GetPermissionByCode(ctx context.Context, code string) (*auth.Permission, error) {
	if p, ok := m.permissions[code]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockPermissionServiceForAPI) CreatePermission(ctx context.Context, permission *auth.Permission) error {
	m.permissions[permission.Code] = permission
	return nil
}

func (m *MockPermissionServiceForAPI) UpdatePermission(ctx context.Context, permission *auth.Permission) error {
	if _, ok := m.permissions[permission.Code]; ok {
		m.permissions[permission.Code] = permission
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForAPI) DeletePermission(ctx context.Context, code string) error {
	delete(m.permissions, code)
	return nil
}

// ==================== 角色管理 ====================

func (m *MockPermissionServiceForAPI) GetAllRoles(ctx context.Context) ([]*auth.Role, error) {
	result := make([]*auth.Role, 0, len(m.roles))
	for _, r := range m.roles {
		result = append(result, r)
	}
	return result, nil
}

func (m *MockPermissionServiceForAPI) GetRoleByID(ctx context.Context, roleID string) (*auth.Role, error) {
	if r, ok := m.roles[roleID]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *MockPermissionServiceForAPI) GetRoleByName(ctx context.Context, name string) (*auth.Role, error) {
	for _, r := range m.roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, nil
}

func (m *MockPermissionServiceForAPI) CreateRole(ctx context.Context, role *auth.Role) error {
	role.ID = primitive.NewObjectID().Hex()
	m.roles[role.ID] = role
	return nil
}

func (m *MockPermissionServiceForAPI) UpdateRole(ctx context.Context, role *auth.Role) error {
	if _, ok := m.roles[role.ID]; ok {
		m.roles[role.ID] = role
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForAPI) DeleteRole(ctx context.Context, roleID string) error {
	delete(m.roles, roleID)
	return nil
}

func (m *MockPermissionServiceForAPI) AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error {
	if r, ok := m.roles[roleID]; ok {
		r.Permissions = append(r.Permissions, permissionCode)
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForAPI) RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error {
	if r, ok := m.roles[roleID]; ok {
		for i, p := range r.Permissions {
			if p == permissionCode {
				r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
				break
			}
		}
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForAPI) GetRolePermissions(ctx context.Context, roleID string) ([]*auth.Permission, error) {
	r, err := m.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return []*auth.Permission{}, nil
	}

	result := make([]*auth.Permission, 0)
	for _, code := range r.Permissions {
		if p := m.permissions[code]; p != nil {
			result = append(result, p)
		}
	}
	return result, nil
}

// ==================== 用户角色管理 ====================

func (m *MockPermissionServiceForAPI) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	if roles, ok := m.userRoles[userID]; ok {
		return roles, nil
	}
	return []string{}, nil
}

func (m *MockPermissionServiceForAPI) AssignRoleToUser(ctx context.Context, userID, roleName string) error {
	m.userRoles[userID] = append(m.userRoles[userID], roleName)
	return nil
}

func (m *MockPermissionServiceForAPI) RemoveRoleFromUser(ctx context.Context, userID, roleName string) error {
	if roles, ok := m.userRoles[userID]; ok {
		for i, r := range roles {
			if r == roleName {
				m.userRoles[userID] = append(roles[:i], roles[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (m *MockPermissionServiceForAPI) ClearUserRoles(ctx context.Context, userID string) error {
	delete(m.userRoles, userID)
	return nil
}

// ==================== 权限检查 ====================

func (m *MockPermissionServiceForAPI) UserHasPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	roles, _ := m.GetUserRoles(ctx, userID)
	if len(roles) == 0 {
		return false, nil
	}

	for _, roleName := range roles {
		for _, r := range m.roles {
			if r.Name == roleName {
				for _, p := range r.Permissions {
					if p == permissionCode {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func (m *MockPermissionServiceForAPI) UserHasAnyPermission(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	for _, code := range permissionCodes {
		has, _ := m.UserHasPermission(ctx, userID, code)
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPermissionServiceForAPI) UserHasAllPermissions(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	for _, code := range permissionCodes {
		has, _ := m.UserHasPermission(ctx, userID, code)
		if !has {
			return false, nil
		}
	}
	return true, nil
}

func (m *MockPermissionServiceForAPI) GetUserPermissions(ctx context.Context, userID string) ([]*auth.Permission, error) {
	roles, _ := m.GetUserRoles(ctx, userID)
	if len(roles) == 0 {
		return []*auth.Permission{}, nil
	}

	permissionCodes := make(map[string]bool)
	for _, roleName := range roles {
		for _, r := range m.roles {
			if r.Name == roleName {
				for _, p := range r.Permissions {
					permissionCodes[p] = true
				}
			}
		}
	}

	result := make([]*auth.Permission, 0)
	for code := range permissionCodes {
		if p := m.permissions[code]; p != nil {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockPermissionServiceForAPI) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	permissionCode := resource + "." + action
	return m.UserHasPermission(ctx, userID, permissionCode)
}

// ==================== 测试辅助函数 ====================

// setupTestRouter 设置测试路由
func setupTestRouter() (*gin.Engine, *PermissionAPI, *MockPermissionServiceForAPI) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := NewMockPermissionServiceForAPI()
	api := NewPermissionAPI(mockService)

	return router, api, mockService
}

// ==================== 权限API测试 ====================

// TestPermissionAPI_GetAllPermissions 测试获取所有权限
func TestPermissionAPI_GetAllPermissions(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	// 创建测试权限
	perm1 := &auth.Permission{
		Code:        "users.read",
		Name:        "读取用户",
		Description: "查看用户信息",
		Resource:    "users",
		Action:      "read",
	}
	perm2 := &auth.Permission{
		Code:        "users.write",
		Name:        "写入用户",
		Description: "修改用户信息",
		Resource:    "users",
		Action:      "write",
	}
	mockService.CreatePermission(ctx, perm1)
	mockService.CreatePermission(ctx, perm2)

	// 设置路由
	router.GET("/admin/permissions", api.GetAllPermissions)

	// 发送请求
	req, _ := http.NewRequest("GET", "/admin/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_GetPermission 测试获取单个权限
func TestPermissionAPI_GetPermission(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm := &auth.Permission{
		Code:        "books.read",
		Name:        "阅读书籍",
		Description: "允许阅读书籍内容",
		Resource:    "books",
		Action:      "read",
	}
	mockService.CreatePermission(ctx, perm)

	router.GET("/admin/permissions/:code", api.GetPermission)

	req, _ := http.NewRequest("GET", "/admin/permissions/books.read", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_CreatePermission 测试创建权限
func TestPermissionAPI_CreatePermission(t *testing.T) {
	router, api, _ := setupTestRouter()

	router.POST("/admin/permissions", api.CreatePermission)

	reqBody := auth.Permission{
		Code:        "books.publish",
		Name:        "发布书籍",
		Description: "允许发布书籍",
		Resource:    "books",
		Action:      "publish",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/admin/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "创建成功", response["message"])
}

// TestPermissionAPI_UpdatePermission 测试更新权限
func TestPermissionAPI_UpdatePermission(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm := &auth.Permission{
		Code:        "books.edit",
		Name:        "编辑书籍",
		Description: "允许编辑书籍",
		Resource:    "books",
		Action:      "edit",
	}
	mockService.CreatePermission(ctx, perm)

	router.PUT("/admin/permissions/:code", api.UpdatePermission)

	reqBody := auth.Permission{
		Code:        "books.edit",
		Name:        "编辑书籍（已更新）",
		Description: "允许编辑书籍内容",
		Resource:    "books",
		Action:      "edit",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/admin/permissions/books.edit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "更新成功", response["message"])
}

// TestPermissionAPI_DeletePermission 测试删除权限
func TestPermissionAPI_DeletePermission(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm := &auth.Permission{
		Code:        "books.delete",
		Name:        "删除书籍",
		Description: "允许删除书籍",
		Resource:    "books",
		Action:      "delete",
	}
	mockService.CreatePermission(ctx, perm)

	router.DELETE("/admin/permissions/:code", api.DeletePermission)

	req, _ := http.NewRequest("DELETE", "/admin/permissions/books.delete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "删除成功", response["message"])
}

// ==================== 角色API测试 ====================

// TestPermissionAPI_GetAllRoles 测试获取所有角色
func TestPermissionAPI_GetAllRoles(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	role1 := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"users.read", "users.write"},
		IsSystem:    false,
	}
	role2 := &auth.Role{
		Name:        "editor",
		Description: "编辑",
		Permissions: []string{"books.write"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role1)
	mockService.CreateRole(ctx, role2)

	router.GET("/admin/roles", api.GetAllRoles)

	req, _ := http.NewRequest("GET", "/admin/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_GetRole 测试获取单个角色
func TestPermissionAPI_GetRole(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{"books.write", "books.publish"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.GET("/admin/roles/:id", api.GetRole)

	req, _ := http.NewRequest("GET", "/admin/roles/"+role.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_CreateRole 测试创建角色
func TestPermissionAPI_CreateRole(t *testing.T) {
	router, api, _ := setupTestRouter()

	router.POST("/admin/roles", api.CreateRole)

	reqBody := auth.Role{
		Name:        "moderator",
		Description: "审核员",
		Permissions: []string{"books.audit"},
		IsSystem:    false,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/admin/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "创建成功", response["message"])
}

// TestPermissionAPI_UpdateRole 测试更新角色
func TestPermissionAPI_UpdateRole(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	role := &auth.Role{
		Name:        "editor",
		Description: "编辑",
		Permissions: []string{"books.write"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.PUT("/admin/roles/:id", api.UpdateRole)

	reqBody := auth.Role{
		Name:        "editor",
		Description: "高级编辑",
		Permissions: []string{"books.write", "books.publish"},
		IsSystem:    false,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/admin/roles/"+role.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "更新成功", response["message"])
}

// TestPermissionAPI_DeleteRole 测试删除角色
func TestPermissionAPI_DeleteRole(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	role := &auth.Role{
		Name:        "temp_role",
		Description: "临时角色",
		Permissions: []string{},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.DELETE("/admin/roles/:id", api.DeleteRole)

	req, _ := http.NewRequest("DELETE", "/admin/roles/"+role.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "删除成功", response["message"])
}

// TestPermissionAPI_AssignPermissionToRole 测试为角色分配权限
func TestPermissionAPI_AssignPermissionToRole(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm := &auth.Permission{
		Code:     "books.manage",
		Name:     "管理书籍",
		Resource: "books",
		Action:   "manage",
	}
	mockService.CreatePermission(ctx, perm)

	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.POST("/admin/roles/:id/permissions/:permissionCode", api.AssignPermissionToRole)

	req, _ := http.NewRequest("POST", "/admin/roles/"+role.ID+"/permissions/books.manage", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "分配成功", response["message"])
}

// TestPermissionAPI_RemovePermissionFromRole 测试移除角色权限
func TestPermissionAPI_RemovePermissionFromRole(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm := &auth.Permission{
		Code:     "books.delete",
		Name:     "删除书籍",
		Resource: "books",
		Action:   "delete",
	}
	mockService.CreatePermission(ctx, perm)

	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"books.delete"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.DELETE("/admin/roles/:id/permissions/:permissionCode", api.RemovePermissionFromRole)

	req, _ := http.NewRequest("DELETE", "/admin/roles/"+role.ID+"/permissions/books.delete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "移除成功", response["message"])
}

// TestPermissionAPI_GetRolePermissions 测试获取角色权限
func TestPermissionAPI_GetRolePermissions(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm1 := &auth.Permission{
		Code:     "users.read",
		Name:     "读用户",
		Resource: "users",
		Action:   "read",
	}
	perm2 := &auth.Permission{
		Code:     "users.write",
		Name:     "写用户",
		Resource: "users",
		Action:   "write",
	}
	mockService.CreatePermission(ctx, perm1)
	mockService.CreatePermission(ctx, perm2)

	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"users.read", "users.write"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.GET("/admin/roles/:id/permissions", api.GetRolePermissions)

	req, _ := http.NewRequest("GET", "/admin/roles/"+role.ID+"/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_AssignRoleToUser 测试为用户分配角色
func TestPermissionAPI_AssignRoleToUser(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	router.POST("/admin/users/:userId/roles", api.AssignRoleToUser)

	reqBody := AssignRoleRequest{Role: "author"}
	body, _ := json.Marshal(reqBody)

	userID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/admin/users/"+userID+"/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "分配成功", response["message"])
}

// TestPermissionAPI_GetUserRoles 测试获取用户角色
func TestPermissionAPI_GetUserRoles(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	role := &auth.Role{
		Name:        "vip",
		Description: "VIP用户",
		Permissions: []string{},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "vip")

	router.GET("/admin/users/:userId/roles", api.GetUserRoles)

	req, _ := http.NewRequest("GET", "/admin/users/"+userID+"/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_GetUserPermissions 测试获取用户权限
func TestPermissionAPI_GetUserPermissions(t *testing.T) {
	router, api, mockService := setupTestRouter()

	ctx := context.Background()
	perm := &auth.Permission{
		Code:     "books.read",
		Name:     "读书",
		Resource: "books",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm)

	role := &auth.Role{
		Name:        "user",
		Description: "普通用户",
		Permissions: []string{"books.read"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "user")

	router.GET("/admin/users/:userId/permissions", api.GetUserPermissions)

	req, _ := http.NewRequest("GET", "/admin/users/"+userID+"/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response["message"])
}

// TestPermissionAPI_InvalidPermissionCode 测试无效权限代码
func TestPermissionAPI_InvalidPermissionCode(t *testing.T) {
	router, api, _ := setupTestRouter()

	router.GET("/admin/permissions/:code", api.GetPermission)

	req, _ := http.NewRequest("GET", "/admin/permissions/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin会返回404因为路径不匹配
	assert.Equal(t, 404, w.Code)
}

// TestPermissionAPI_InvalidRequestBody 测试无效请求体
func TestPermissionAPI_InvalidRequestBody(t *testing.T) {
	router, api, _ := setupTestRouter()

	router.POST("/admin/permissions", api.CreatePermission)

	// 发送无效JSON
	req, _ := http.NewRequest("POST", "/admin/permissions", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Gin的JSON绑定会返回解析错误的详细信息
	_, hasError := response["error"]
	assert.True(t, hasError, "响应应包含error字段")
}
