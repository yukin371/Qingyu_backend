package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPermissionServiceForMiddleware 模拟权限服务（用于中间件测试）
type MockPermissionServiceForMiddleware struct {
	mock.Mock
}

func (m *MockPermissionServiceForMiddleware) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	args := m.Called(ctx, userID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionServiceForMiddleware) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionServiceForMiddleware) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionServiceForMiddleware) HasRole(ctx context.Context, userID, role string) (bool, error) {
	args := m.Called(ctx, userID, role)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionServiceForMiddleware) InvalidateUserPermissionsCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPermissionServiceForMiddleware) SetChecker(checker interface{}) {
	// Not used in tests
}

func (m *MockPermissionServiceForMiddleware) LoadPermissionsToChecker(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPermissionServiceForMiddleware) LoadUserRolesToChecker(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPermissionServiceForMiddleware) ReloadAllFromDatabase(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// 设置测试模式
func init() {
	gin.SetMode(gin.TestMode)
}
// TestRequireAdminPermission_Success 测试权限检查成功
func TestRequireAdminPermission_Success(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	mockService.On("HasRole", mock.Anything, "user123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.read").Return(true, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// 创建测试路由
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})
	router.GET("/test", RequireAdminPermission("user.read", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestRequireAdminPermission_NoUserID 测试没有用户ID
func TestRequireAdminPermission_NoUserID(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router := gin.New()
	router.GET("/test", RequireAdminPermission("user.read", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", resp["code"])
}

// TestRequireAdminPermission_NoPermission 测试权限不足
func TestRequireAdminPermission_NoPermission(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	mockService.On("HasRole", mock.Anything, "user123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.read").Return(false, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})
	router.GET("/test", RequireAdminPermission("user.read", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", resp["code"])

	mockService.AssertExpectations(t)
}

// TestRequireAdminPermission_SuperAdminBypass 测试超级管理员绕过
func TestRequireAdminPermission_SuperAdminBypass(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	// 超级管理员
	mockService.On("HasRole", mock.Anything, "admin123", "super_admin").Return(true, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "admin123")
		c.Next()
	})
	router.GET("/test", RequireAdminPermission("user.read", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestRequireAnyAdminPermission_Success 测试任意权限检查成功
func TestRequireAnyAdminPermission_Success(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	// 有 user.write 权限，可以访问需要 user.read 或 user.write 的路由
	mockService.On("HasRole", mock.Anything, "user123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.read").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.write").Return(true, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-any", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})
	router.GET("/test-any", RequireAnyAdminPermission(WithPermissionService(mockService))("user.read", "user.write"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestRequireAnyAdminPermission_NoPermission 测试没有任何权限
func TestRequireAnyAdminPermission_NoPermission(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	mockService.On("HasRole", mock.Anything, "user123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.read").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.write").Return(false, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-any", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})
	router.GET("/test-any", RequireAnyAdminPermission(WithPermissionService(mockService))("user.read", "user.write"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", resp["code"])

	mockService.AssertExpectations(t)
}

// TestRequireResourceOwnerPermission_OwnerAccess 测试资源所有者访问
func TestRequireResourceOwnerPermission_OwnerAccess(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	// 资源所有者应该被允许访问
	mockService.On("HasRole", mock.Anything, "user123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "resource.manage").Return(false, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/resource/user123", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})
	router.GET("/resource/:id", RequireResourceOwnerPermission("resource", "id", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestRequireResourceOwnerPermission_AdminAccess 测试管理员访问
func TestRequireResourceOwnerPermission_AdminAccess(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	mockService.On("HasRole", mock.Anything, "admin123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "admin123", "resource.manage").Return(true, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/resource/user123", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "admin123")
		c.Next()
	})
	router.GET("/resource/:id", RequireResourceOwnerPermission("resource", "id", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestRequireResourceOwnerPermission_Forbidden 测试无权限访问
func TestRequireResourceOwnerPermission_Forbidden(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	mockService.On("HasRole", mock.Anything, "user456", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user456", "resource.manage").Return(false, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/resource/user123", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user456")
		c.Next()
	})
	router.GET("/resource/:id", RequireResourceOwnerPermission("resource", "id", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", resp["code"])

	mockService.AssertExpectations(t)
}

// TestMatchWildcardPermission 测试通配符权限匹配
func TestMatchWildcardPermission(t *testing.T) {
	tests := []struct {
		name              string
		userPermissions   []string
		requiredPermission string
		expected          bool
	}{
		{
			name:              "精确匹配",
			userPermissions:   []string{"user.read"},
			requiredPermission: "user.read",
			expected:          true,
		},
		{
			name:              "通配符匹配成功",
			userPermissions:   []string{"user.*"},
			requiredPermission: "user.read",
			expected:          true,
		},
		{
			name:              "通配符匹配失败",
			userPermissions:   []string{"book.*"},
			requiredPermission: "user.read",
			expected:          false,
		},
		{
			name:              "全通配符",
			userPermissions:   []string{"*"},
			requiredPermission: "anything",
			expected:          true,
		},
		{
			name:              "多层通配符",
			userPermissions:   []string{"admin.user.*"},
			requiredPermission: "admin.user.read",
			expected:          true,
		},
		{
			name:              "无权限",
			userPermissions:   []string{},
			requiredPermission: "user.read",
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchWildcardPermission(tt.userPermissions, tt.requiredPermission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRequireAdminPermission_ServiceError 测试权限服务错误
func TestRequireAdminPermission_ServiceError(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	mockService.On("HasRole", mock.Anything, "user123", "super_admin").Return(false, nil)
	mockService.On("CheckPermission", mock.Anything, "user123", "user.read").Return(false, errors.New("service error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})
	router.GET("/test", RequireAdminPermission("user.read", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "PERMISSION_CHECK_FAILED", resp["code"])

	mockService.AssertExpectations(t)
}

// TestRequireAdminPermission_InvalidUserID 测试无效的用户ID
func TestRequireAdminPermission_InvalidUserID(t *testing.T) {
	mockService := new(MockPermissionServiceForMiddleware)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 12345) // 数字而不是字符串
		c.Next()
	})
	router.GET("/test", RequireAdminPermission("user.read", WithPermissionService(mockService)), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "INVALID_USER_ID", resp["code"])
}
