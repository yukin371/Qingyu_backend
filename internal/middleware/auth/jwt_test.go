package auth

import (
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
