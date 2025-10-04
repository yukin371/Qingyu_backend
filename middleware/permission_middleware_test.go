package middleware

import (
	"context"
	"net/http/httptest"
	"testing"

	"Qingyu_backend/service/shared/auth"

	"github.com/gin-gonic/gin"
)

// MockAuthServiceForPermission Mock Auth服务（权限）
type MockAuthServiceForPermission struct {
	checkPermFunc func(ctx context.Context, userID, permission string) (bool, error)
}

func (m *MockAuthServiceForPermission) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) Logout(ctx context.Context, token string) error {
	return nil
}

func (m *MockAuthServiceForPermission) RefreshToken(ctx context.Context, token string) (string, error) {
	return "", nil
}

func (m *MockAuthServiceForPermission) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	if m.checkPermFunc != nil {
		return m.checkPermFunc(ctx, userID, permission)
	}
	return false, nil
}

func (m *MockAuthServiceForPermission) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return []string{"book.read", "book.write"}, nil
}

func (m *MockAuthServiceForPermission) HasRole(ctx context.Context, userID, role string) (bool, error) {
	return false, nil
}

func (m *MockAuthServiceForPermission) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) UpdateRole(ctx context.Context, roleID string, req *auth.UpdateRoleRequest) error {
	return nil
}

func (m *MockAuthServiceForPermission) DeleteRole(ctx context.Context, roleID string) error {
	return nil
}

func (m *MockAuthServiceForPermission) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthServiceForPermission) RemoveRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthServiceForPermission) CreateSession(ctx context.Context, userID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthServiceForPermission) DestroySession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthServiceForPermission) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthServiceForPermission) Health(ctx context.Context) error {
	return nil
}

// ============ 测试用例 ============

// TestRequirePermission_Success 测试权限检查成功
func TestRequirePermission_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{
		checkPermFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			return userID == "user123" && permission == "book.read", nil
		},
	}

	middleware := NewPermissionMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequirePermission("book.read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("权限检查成功测试通过")
}

// TestRequirePermission_Fail 测试权限不足
func TestRequirePermission_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{
		checkPermFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			return false, nil
		},
	}

	middleware := NewPermissionMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequirePermission("book.delete"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("期望状态码403，实际%d", w.Code)
	}

	t.Logf("权限不足测试通过")
}

// TestRequireAnyPermission_Success 测试任一权限成功
func TestRequireAnyPermission_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{
		checkPermFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			return permission == "book.write", nil
		},
	}

	middleware := NewPermissionMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequireAnyPermission("book.read", "book.write"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("任一权限成功测试通过")
}

// TestRequireAllPermissions_Success 测试所有权限成功
func TestRequireAllPermissions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{
		checkPermFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			return true, nil
		},
	}

	middleware := NewPermissionMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequireAllPermissions("book.read", "book.write"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("所有权限成功测试通过")
}

// TestRequireAllPermissions_Fail 测试所有权限失败
func TestRequireAllPermissions_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{
		checkPermFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			// 只有book.read权限
			return permission == "book.read", nil
		},
	}

	middleware := NewPermissionMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequireAllPermissions("book.read", "book.write"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("期望状态码403，实际%d", w.Code)
	}

	t.Logf("所有权限失败测试通过")
}

// TestCheckResourcePermission 测试资源权限
func TestCheckResourcePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{
		checkPermFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			return permission == "book.read", nil
		},
	}

	middleware := NewPermissionMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.CheckResourcePermission("book", "read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("资源权限测试通过")
}

// TestGetUserPermissions 测试获取用户权限
func TestGetUserPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthServiceForPermission{}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "user123")

	perms, err := GetUserPermissions(c, mockAuth)
	if err != nil {
		t.Fatalf("获取权限失败: %v", err)
	}

	if len(perms) != 2 {
		t.Errorf("期望2个权限，实际%d个", len(perms))
	}

	t.Logf("获取用户权限测试通过: %v", perms)
}
