package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/auth"
)

// MockPermissionServiceForMiddleware 模拟权限服务（用于中间件测试）
type MockPermissionServiceForMiddleware struct {
	permissions map[string]*auth.Permission
	roles       map[string]*auth.Role
	userRoles   map[string][]string // userID -> roleNames
}

func NewMockPermissionServiceForMiddleware() *MockPermissionServiceForMiddleware {
	return &MockPermissionServiceForMiddleware{
		permissions: make(map[string]*auth.Permission),
		roles:       make(map[string]*auth.Role),
		userRoles:   make(map[string][]string),
	}
}

// ==================== 权限管理 ====================

func (m *MockPermissionServiceForMiddleware) GetAllPermissions(ctx context.Context) ([]*auth.Permission, error) {
	result := make([]*auth.Permission, 0, len(m.permissions))
	for _, p := range m.permissions {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockPermissionServiceForMiddleware) GetPermissionByCode(ctx context.Context, code string) (*auth.Permission, error) {
	if p, ok := m.permissions[code]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockPermissionServiceForMiddleware) CreatePermission(ctx context.Context, permission *auth.Permission) error {
	m.permissions[permission.Code] = permission
	return nil
}

func (m *MockPermissionServiceForMiddleware) UpdatePermission(ctx context.Context, permission *auth.Permission) error {
	if _, ok := m.permissions[permission.Code]; ok {
		m.permissions[permission.Code] = permission
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForMiddleware) DeletePermission(ctx context.Context, code string) error {
	delete(m.permissions, code)
	return nil
}

// ==================== 角色管理 ====================

func (m *MockPermissionServiceForMiddleware) GetAllRoles(ctx context.Context) ([]*auth.Role, error) {
	result := make([]*auth.Role, 0, len(m.roles))
	for _, r := range m.roles {
		result = append(result, r)
	}
	return result, nil
}

func (m *MockPermissionServiceForMiddleware) GetRoleByID(ctx context.Context, roleID string) (*auth.Role, error) {
	if r, ok := m.roles[roleID]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *MockPermissionServiceForMiddleware) GetRoleByName(ctx context.Context, name string) (*auth.Role, error) {
	for _, r := range m.roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, nil
}

func (m *MockPermissionServiceForMiddleware) CreateRole(ctx context.Context, role *auth.Role) error {
	role.ID = primitive.NewObjectID().Hex()
	m.roles[role.ID] = role
	return nil
}

func (m *MockPermissionServiceForMiddleware) UpdateRole(ctx context.Context, role *auth.Role) error {
	if _, ok := m.roles[role.ID]; ok {
		m.roles[role.ID] = role
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForMiddleware) DeleteRole(ctx context.Context, roleID string) error {
	delete(m.roles, roleID)
	return nil
}

func (m *MockPermissionServiceForMiddleware) AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error {
	if r, ok := m.roles[roleID]; ok {
		r.Permissions = append(r.Permissions, permissionCode)
		return nil
	}
	return nil
}

func (m *MockPermissionServiceForMiddleware) RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error {
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

func (m *MockPermissionServiceForMiddleware) GetRolePermissions(ctx context.Context, roleID string) ([]*auth.Permission, error) {
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

func (m *MockPermissionServiceForMiddleware) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	if roles, ok := m.userRoles[userID]; ok {
		return roles, nil
	}
	return []string{}, nil
}

func (m *MockPermissionServiceForMiddleware) AssignRoleToUser(ctx context.Context, userID, roleName string) error {
	m.userRoles[userID] = append(m.userRoles[userID], roleName)
	return nil
}

func (m *MockPermissionServiceForMiddleware) RemoveRoleFromUser(ctx context.Context, userID, roleName string) error {
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

func (m *MockPermissionServiceForMiddleware) ClearUserRoles(ctx context.Context, userID string) error {
	delete(m.userRoles, userID)
	return nil
}

// ==================== 权限检查 ====================

func (m *MockPermissionServiceForMiddleware) UserHasPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
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

func (m *MockPermissionServiceForMiddleware) UserHasAnyPermission(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	for _, code := range permissionCodes {
		has, _ := m.UserHasPermission(ctx, userID, code)
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPermissionServiceForMiddleware) UserHasAllPermissions(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	for _, code := range permissionCodes {
		has, _ := m.UserHasPermission(ctx, userID, code)
		if !has {
			return false, nil
		}
	}
	return true, nil
}

func (m *MockPermissionServiceForMiddleware) GetUserPermissions(ctx context.Context, userID string) ([]*auth.Permission, error) {
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

func (m *MockPermissionServiceForMiddleware) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	permissionCode := resource + "." + action
	return m.UserHasPermission(ctx, userID, permissionCode)
}

// ==================== 测试用例 ====================

// setupTestContext 设置测试上下文
func setupTestContext() (*gin.Engine, *RBACMiddleware, *MockPermissionServiceForMiddleware) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := NewMockPermissionServiceForMiddleware()
	middleware := NewRBACMiddleware(mockService)

	return router, middleware, mockService
}

// TestRequirePermission_Unauthenticated 测试未认证用户
func TestRequirePermission_Unauthenticated(t *testing.T) {
	router, middleware, _ := setupTestContext()

	router.GET("/test", middleware.RequirePermission("users.read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "未认证")
}

// TestRequirePermission_NoPermission 测试无权限用户
func TestRequirePermission_NoPermission(t *testing.T) {
	router, middleware, _ := setupTestContext()

	// 设置测试用户（但没有权限）
	userID := primitive.NewObjectID().Hex()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequirePermission("users.read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "权限不足")
}

// TestRequirePermission_HasPermission 测试有权限用户
func TestRequirePermission_HasPermission(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	// 创建权限
	ctx := context.Background()
	perm := &auth.Permission{
		Code:     "users.read",
		Name:     "读取用户",
		Resource: "users",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm)

	// 创建角色并分配权限
	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"users.read"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户角色
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "admin")

	// 设置路由
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequirePermission("users.read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRequireAnyPermission_NoPermission 测试RequireAnyPermission无权限
func TestRequireAnyPermission_NoPermission(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限1
	perm1 := &auth.Permission{
		Code:     "users.read",
		Name:     "读取用户",
		Resource: "users",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm1)

	// 创建权限2（不给用户）
	perm2 := &auth.Permission{
		Code:     "users.write",
		Name:     "写入用户",
		Resource: "users",
		Action:   "write",
	}
	mockService.CreatePermission(ctx, perm2)

	// 创建角色只给权限1
	role := &auth.Role{
		Name:        "editor",
		Description: "编辑",
		Permissions: []string{"users.read"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "editor")

	// 需要任意一个权限：users.write 或 users.delete（用户都没有）
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireAnyPermission("users.write", "users.delete"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

// TestRequireAnyPermission_HasPermission 测试RequireAnyPermission有权限
func TestRequireAnyPermission_HasPermission(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限
	perm1 := &auth.Permission{
		Code:     "users.read",
		Name:     "读取用户",
		Resource: "users",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm1)

	perm2 := &auth.Permission{
		Code:     "users.write",
		Name:     "写入用户",
		Resource: "users",
		Action:   "write",
	}
	mockService.CreatePermission(ctx, perm2)

	// 创建角色给权限1
	role := &auth.Role{
		Name:        "editor",
		Description: "编辑",
		Permissions: []string{"users.read"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "editor")

	// 需要任意一个权限：users.write 或 users.read（用户有read）
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireAnyPermission("users.write", "users.read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRequireAllPermissions_MissingPermission 测试RequireAllPermissions缺少权限
func TestRequireAllPermissions_MissingPermission(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限
	perm1 := &auth.Permission{
		Code:     "users.read",
		Name:     "读取用户",
		Resource: "users",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm1)

	perm2 := &auth.Permission{
		Code:     "users.write",
		Name:     "写入用户",
		Resource: "users",
		Action:   "write",
	}
	mockService.CreatePermission(ctx, perm2)

	// 创建角色只给权限1
	role := &auth.Role{
		Name:        "editor",
		Description: "编辑",
		Permissions: []string{"users.read"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "editor")

	// 需要所有权限：users.read 和 users.write（用户只有read）
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireAllPermissions("users.read", "users.write"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

// TestRequireAllPermissions_HasAllPermissions 测试RequireAllPermissions拥有所有权限
func TestRequireAllPermissions_HasAllPermissions(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限
	perm1 := &auth.Permission{
		Code:     "users.read",
		Name:     "读取用户",
		Resource: "users",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm1)

	perm2 := &auth.Permission{
		Code:     "users.write",
		Name:     "写入用户",
		Resource: "users",
		Action:   "write",
	}
	mockService.CreatePermission(ctx, perm2)

	// 创建角色给两个权限
	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"users.read", "users.write"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "admin")

	// 需要所有权限：users.read 和 users.write（用户都有）
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireAllPermissions("users.read", "users.write"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRequireResourcePermission 测试资源权限检查
func TestRequireResourcePermission(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限
	perm := &auth.Permission{
		Code:     "books.delete",
		Name:     "删除书籍",
		Resource: "books",
		Action:   "delete",
	}
	mockService.CreatePermission(ctx, perm)

	// 创建角色
	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{"books.delete"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "author")

	// 使用资源权限中间件
	router.DELETE("/books/:id", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireResourcePermission("books", "delete"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("DELETE", "/books/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRequireRole_NoRole 测试无角色用户
func TestRequireRole_NoRole(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建角色
	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户（但没分配角色）
	userID := primitive.NewObjectID().Hex()

	// 需要admin角色
	router.GET("/admin", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "权限不足")
}

// TestRequireRole_HasRole 测试有角色用户
func TestRequireRole_HasRole(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建角色
	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户并分配角色
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "admin")

	// 需要admin角色
	router.GET("/admin", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRequireRole_MultipleRoles 测试多角色检查
func TestRequireRole_MultipleRoles(t *testing.T) {
	router, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建角色
	role1 := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{},
		IsSystem:    false,
	}
	_ = mockService.CreateRole(ctx, role1)

	role2 := &auth.Role{
		Name:        "editor",
		Description: "编辑",
		Permissions: []string{},
		IsSystem:    false,
	}
	_ = mockService.CreateRole(ctx, role2)

	// 设置用户为editor角色
	userID := primitive.NewObjectID().Hex()
	_ = mockService.AssignRoleToUser(ctx, userID, "editor")

	// 需要admin或editor角色之一
	router.GET("/content", func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}, middleware.RequireRole("admin", "editor"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/content", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRBAC_GetUserPermissions 测试获取用户权限辅助函数
func TestRBAC_GetUserPermissions(t *testing.T) {
	_, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限
	perm1 := &auth.Permission{
		Code:     "users.read",
		Name:     "读取用户",
		Resource: "users",
		Action:   "read",
	}
	mockService.CreatePermission(ctx, perm1)

	perm2 := &auth.Permission{
		Code:     "users.write",
		Name:     "写入用户",
		Resource: "users",
		Action:   "write",
	}
	mockService.CreatePermission(ctx, perm2)

	// 创建角色
	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"users.read", "users.write"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "admin")

	// 创建gin上下文（带请求）
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建请求并设置到上下文
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req
	c.Set("user_id", userID)

	// 获取用户权限
	permissions, err := middleware.GetUserPermissions(c)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(permissions))
}

// TestHasPermission 测试权限检查辅助函数
func TestHasPermission(t *testing.T) {
	_, middleware, mockService := setupTestContext()

	ctx := context.Background()
	// 创建权限
	perm := &auth.Permission{
		Code:     "books.publish",
		Name:     "发布书籍",
		Resource: "books",
		Action:   "publish",
	}
	mockService.CreatePermission(ctx, perm)

	// 创建角色
	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{"books.publish"},
		IsSystem:    false,
	}
	mockService.CreateRole(ctx, role)

	// 设置用户
	userID := primitive.NewObjectID().Hex()
	mockService.AssignRoleToUser(ctx, userID, "author")

	// 创建gin上下文（带请求）
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建请求并设置到上下文
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req
	c.Set("user_id", userID)

	// 检查权限
	has, err := middleware.HasPermission(c, "books.publish")
	assert.NoError(t, err)
	assert.True(t, has)

	// 检查不存在的权限
	has, err = middleware.HasPermission(c, "books.delete")
	assert.NoError(t, err)
	assert.False(t, has)
}
