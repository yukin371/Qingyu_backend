package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestAuthMiddleware_ValidToken 测试有效Token通过验证
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 生成有效Token
	userID := "user123"
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, nil)
	require.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, "user123", userID)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// 执行
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthMiddleware_MissingToken 测试缺少Token返回2010错误
func TestAuthMiddleware_MissingToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行 - 不设置Authorization Header
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// 验证错误码为2010（Token缺失）
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "2010", response["code"])
}

// TestAuthMiddleware_InvalidFormat 测试Token格式错误返回2009错误
func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行 - 设置无效格式的Token（缺少Bearer前缀）
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "invalid-token-format")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 验证错误码为2009（Token格式错误）
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "2009", response["code"])
}

// TestAuthMiddleware_InvalidToken 测试Token无效返回2008错误
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行 - 设置无效的Token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// 验证错误码为2008（Token无效）
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "2008", response["code"])
}

// TestAuthMiddleware_ExpiredToken 测试Token过期返回2007错误
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 创建一个已过期的Token（手动构造JWT）
	now := time.Now()
	claims := &Claims{
		UserID:    "user123",
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "qingyu",
			Subject:   "user123",
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)), // 1小时前过期
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// 验证错误码为2007（Token过期）
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "2007", response["code"])
}

// TestJWTManager_RefreshToken 测试Token刷新功能
func TestJWTManager_RefreshToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	// 生成Token对
	userID := "user123"
	accessToken, refreshToken, err := jwtManager.GenerateTokenPair(userID, nil)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	// 使用RefreshToken获取新的AccessToken
	newAccessToken, err := jwtManager.RefreshToken(refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, newAccessToken)
	require.NotEqual(t, accessToken, newAccessToken)

	// 验证新Token
	claims, err := jwtManager.ValidateToken(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "access", claims.TokenType)
}

// TestJWTManager_RefreshToken_WithInvalidToken 测试使用无效的RefreshToken
func TestJWTManager_RefreshToken_WithInvalidToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	// 使用无效的RefreshToken
	_, err = jwtManager.RefreshToken("invalid.refresh.token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

// TestJWTManager_RefreshToken_WithAccessToken 测试使用AccessToken而不是RefreshToken
func TestJWTManager_RefreshToken_WithAccessToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	// 生成Token对
	userID := "user123"
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, nil)
	require.NoError(t, err)

	// 尝试使用AccessToken刷新（应该失败）
	_, err = jwtManager.RefreshToken(accessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a refresh token")
}

// TestAuthMiddleware_BlacklistedToken 测试黑名单Token被拒绝
func TestAuthMiddleware_BlacklistedToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 生成有效Token
	userID := "user123"
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, nil)
	require.NoError(t, err)

	// 将Token添加到黑名单
	ctx := context.Background()
	err = blacklist.Add(ctx, accessToken, 1*time.Hour)
	require.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证 - 黑名单Token应该被拒绝
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "2016", response["code"])
}

// TestAuthMiddleware_NotBlacklistedToken 测试未黑名单Token通过验证
func TestAuthMiddleware_NotBlacklistedToken(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 生成有效Token
	userID := "user123"
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, nil)
	require.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行 - Token不在黑名单中
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证 - 未黑名单Token应该通过
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthMiddleware_SkipPaths 测试跳过指定路径
func TestAuthMiddleware_SkipPaths(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 配置跳过路径
	config := middleware.GetConfig()
	config.SkipPaths = []string{"/health", "/metrics"}

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"metrics": "ok"})
	})
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"api": "ok"})
	})

	// 测试/health路径 - 不需要Token
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 测试/metrics路径 - 不需要Token
	req = httptest.NewRequest("GET", "/metrics", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 测试/api/test路径 - 需要Token
	req = httptest.NewRequest("GET", "/api/test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestJWTManager_GenerateTokenPair_WithExtraClaims 测试生成Token时包含额外信息
func TestJWTManager_GenerateTokenPair_WithExtraClaims(t *testing.T) {
	// 准备
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	// 生成Token对，包含额外信息
	userID := "user123"
	extraClaims := map[string]interface{}{
		"username": "testuser",
		"roles":    []string{"admin", "user"},
	}
	accessToken, refreshToken, err := jwtManager.GenerateTokenPair(userID, extraClaims)
	require.NoError(t, err)

	// 验证AccessToken包含额外信息
	claims, err := jwtManager.ValidateToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, []string{"admin", "user"}, claims.Roles)

	// 验证RefreshToken不包含额外信息（只包含基本信息）
	refreshClaims, err := jwtManager.ValidateToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
	assert.Equal(t, "refresh", refreshClaims.TokenType)
}
