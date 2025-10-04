package middleware

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"Qingyu_backend/service/shared/auth"

	"github.com/gin-gonic/gin"
)

// MockAuthService Mock Auth服务
type MockAuthService struct {
	validateFunc func(ctx context.Context, token string) (*auth.TokenClaims, error)
	hasRoleFunc  func(ctx context.Context, userID, role string) (bool, error)
}

func (m *MockAuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return nil, nil
}

func (m *MockAuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return nil, nil
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	return nil
}

func (m *MockAuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	return "", nil
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	if m.validateFunc != nil {
		return m.validateFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockAuthService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return false, nil
}

func (m *MockAuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (m *MockAuthService) HasRole(ctx context.Context, userID, role string) (bool, error) {
	if m.hasRoleFunc != nil {
		return m.hasRoleFunc(ctx, userID, role)
	}
	return false, nil
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (m *MockAuthService) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	return nil, nil
}

func (m *MockAuthService) UpdateRole(ctx context.Context, roleID string, req *auth.UpdateRoleRequest) error {
	return nil
}

func (m *MockAuthService) DeleteRole(ctx context.Context, roleID string) error {
	return nil
}

func (m *MockAuthService) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthService) RemoveRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *MockAuthService) CreateSession(ctx context.Context, userID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthService) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	return nil, nil
}

func (m *MockAuthService) DestroySession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthService) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockAuthService) Health(ctx context.Context) error {
	return nil
}

// ============ 测试用例 ============

// TestRequireAuth_Success 测试成功认证
func TestRequireAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{
		validateFunc: func(ctx context.Context, token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID: "user123",
				Roles:  []string{"reader"},
			}, nil
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(200, gin.H{"user_id": userID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("认证成功测试通过")
}

// TestRequireAuth_NoToken 测试缺少Token
func TestRequireAuth_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("期望状态码401，实际%d", w.Code)
	}

	t.Logf("缺少Token测试通过")
}

// TestRequireAuth_InvalidToken 测试无效Token
func TestRequireAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{
		validateFunc: func(ctx context.Context, token string) (*auth.TokenClaims, error) {
			return nil, fmt.Errorf("token无效")
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("期望状态码401，实际%d", w.Code)
	}

	t.Logf("无效Token测试通过")
}

// TestOptionalAuth_WithToken 测试可选认证（有Token）
func TestOptionalAuth_WithToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{
		validateFunc: func(ctx context.Context, token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID: "user123",
				Roles:  []string{"reader"},
			}, nil
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", middleware.OptionalAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		c.JSON(200, gin.H{"authenticated": exists, "user_id": userID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("可选认证（有Token）测试通过")
}

// TestOptionalAuth_NoToken 测试可选认证（无Token）
func TestOptionalAuth_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", middleware.OptionalAuth(), func(c *gin.Context) {
		_, exists := c.Get("user_id")
		c.JSON(200, gin.H{"authenticated": exists})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("可选认证（无Token）测试通过")
}

// TestRequireRole_Success 测试角色检查成功
func TestRequireRole_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{
		hasRoleFunc: func(ctx context.Context, userID, role string) (bool, error) {
			return userID == "user123" && role == "admin", nil
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	t.Logf("角色检查成功测试通过")
}

// TestRequireRole_Fail 测试角色检查失败
func TestRequireRole_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := &MockAuthService{
		hasRoleFunc: func(ctx context.Context, userID, role string) (bool, error) {
			return false, nil
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}, middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("期望状态码403，实际%d", w.Code)
	}

	t.Logf("角色检查失败测试通过")
}

// TestGetUserID 测试辅助函数
func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "user123")

	userID, exists := GetUserID(c)
	if !exists {
		t.Error("应该存在user_id")
	}
	if userID != "user123" {
		t.Errorf("user_id错误: %s", userID)
	}

	t.Logf("GetUserID测试通过")
}
