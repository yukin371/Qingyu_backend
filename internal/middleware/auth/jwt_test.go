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

	// 等待一小段时间，确保新Token的iat时间戳不同
	time.Sleep(1 * time.Second)

	// 使用RefreshToken获取新的AccessToken
	newAccessToken, err := jwtManager.RefreshToken(refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, newAccessToken)

	// 验证新Token
	claims, err := jwtManager.ValidateToken(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "access", claims.TokenType)

	// 验证新Token的iat时间比旧Token晚
	oldClaims, err := jwtManager.ValidateToken(accessToken)
	require.NoError(t, err)
	newClaims, err := jwtManager.ValidateToken(newAccessToken)
	require.NoError(t, err)
	assert.Greater(t, newClaims.IssuedAt.Time, oldClaims.IssuedAt.Time)
}

// TestJWTAuthMiddleware_NameAndPriority 测试中间件名称和优先级
func TestJWTAuthMiddleware_NameAndPriority(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	assert.Equal(t, "jwt_auth", middleware.Name())
	assert.Equal(t, 9, middleware.Priority())
}

// TestJWTAuthMiddleware_Getters 测试Getter方法
func TestJWTAuthMiddleware_Getters(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	config := middleware.GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "qingyu", config.Issuer)
	assert.Equal(t, "Authorization", config.TokenHeader)
	assert.Equal(t, "Bearer", config.TokenPrefix)
	assert.True(t, config.Enabled)

	assert.Same(t, jwtManager, middleware.GetJWTManager())
	assert.Same(t, blacklist, middleware.GetBlacklist())
}

// TestJWTAuthMiddleware_LoadConfig 测试配置加载
func TestJWTAuthMiddleware_LoadConfig(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 加载配置
	config := map[string]interface{}{
		"enabled": false,
		"secret":  "new-secret",
		"issuer":  "test-issuer",
		"token_header": "X-Auth-Token",
		"token_prefix": "Token",
		"skip_paths": []interface{}{"/public", "/static"},
	}

	err = middleware.LoadConfig(config)
	require.NoError(t, err)

	// 验证配置已加载
	middlewareConfig := middleware.GetConfig()
	assert.False(t, middlewareConfig.Enabled)
	assert.Equal(t, "test-issuer", middlewareConfig.Issuer)
	assert.Equal(t, "X-Auth-Token", middlewareConfig.TokenHeader)
	assert.Equal(t, "Token", middlewareConfig.TokenPrefix)
	assert.Equal(t, []string{"/public", "/static"}, middlewareConfig.SkipPaths)
}

// TestJWTAuthMiddleware_ValidateConfig 测试配置验证
func TestJWTAuthMiddleware_ValidateConfig(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		secret := "test-secret-key-123456"
		jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
		require.NoError(t, err)

		blacklist := NewMockBlacklist()
		middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
		require.NotNil(t, middleware)

		// 设置有效配置
		config := middleware.GetConfig()
		config.Secret = secret

		err = middleware.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("EmptySecret", func(t *testing.T) {
		secret := "test-secret-key-123456"
		jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
		require.NoError(t, err)

		blacklist := NewMockBlacklist()
		middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
		require.NotNil(t, middleware)

		// 设置空Secret
		config := middleware.GetConfig()
		config.Secret = ""

		err = middleware.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "secret")
	})

	t.Run("InvalidAccessExpiration", func(t *testing.T) {
		secret := "test-secret-key-123456"
		jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
		require.NoError(t, err)

		blacklist := NewMockBlacklist()
		middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
		require.NotNil(t, middleware)

		// 设置无效的过期时间
		config := middleware.GetConfig()
		config.Secret = secret
		config.AccessExpiration = -1 * time.Hour

		err = middleware.ValidateConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "access_expiration")
	})
}

// TestJWTManager_ParseToken 测试解析Token（不验证签名）
func TestJWTManager_ParseToken(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	// 生成Token
	userID := "user123"
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, nil)
	require.NoError(t, err)

	// 解析Token
	claims, err := jwtManager.ParseToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "access", claims.TokenType)
}

// TestBlacklist_Remove 测试从黑名单移除Token
func TestBlacklist_Remove(t *testing.T) {
	blacklist := NewMockBlacklist()
	ctx := context.Background()

	token := "test-token"

	// 添加到黑名单
	err := blacklist.Add(ctx, token, 1*time.Hour)
	require.NoError(t, err)

	// 验证在黑名单中
	isBlacklisted, err := blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.True(t, isBlacklisted)

	// 从黑名单移除
	err = blacklist.Remove(ctx, token)
	require.NoError(t, err)

	// 验证不在黑名单中
	isBlacklisted, err = blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.False(t, isBlacklisted)
}

// TestAuthMiddleware_Disabled 测试禁用认证时跳过验证
func TestAuthMiddleware_Disabled(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 禁用认证
	config := middleware.GetConfig()
	config.Enabled = false

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 执行 - 不提供Token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证 - 应该通过（认证被禁用）
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthMiddleware_CustomHeader 测试自定义Token Header
func TestAuthMiddleware_CustomHeader(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 使用自定义Header
	config := middleware.GetConfig()
	config.TokenHeader = "X-Auth-Token"
	config.TokenPrefix = "Token"

	// 生成Token
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

	// 执行 - 使用自定义Header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Auth-Token", "Token "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)
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

// TestJWTManager_NewJWTManager_WithError 测试创建JWT管理器时的错误
func TestJWTManager_NewJWTManager_WithError(t *testing.T) {
	// 空Secret应该返回错误
	_, err := NewJWTManager("", DefaultAccessExpiration, DefaultRefreshExpiration)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secret")
}

// TestJWTManager_ParseToken_InvalidToken 测试解析无效Token
func TestJWTManager_ParseToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	// 解析无效Token
	_, err = jwtManager.ParseToken("invalid.token")
	assert.Error(t, err)
}

// TestAuthMiddleware_ValidateToken_UnexpectedSigningMethod 测试非预期签名算法
func TestAuthMiddleware_ValidateToken_UnexpectedSigningMethod(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 创建使用不同算法的Token（使用None算法）
	now := time.Now()
	claims := &Claims{
		UserID:    "user123",
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "qingyu",
			Subject:   "user123",
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	token.Header["alg"] = "none"
	accessToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
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

	// 验证 - 应该返回Token无效错误
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "2008", response["code"])
}

// TestAuthMiddleware_WithUserInfoInContext 测试用户信息注入到上下文
func TestAuthMiddleware_WithUserInfoInContext(t *testing.T) {
	secret := "test-secret-key-123456"
	jwtManager, err := NewJWTManager(secret, DefaultAccessExpiration, DefaultRefreshExpiration)
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	middleware := NewJWTAuthMiddleware(jwtManager, blacklist, nil)
	require.NotNil(t, middleware)

	// 生成包含额外信息的Token
	userID := "user123"
	extraClaims := map[string]interface{}{
		"username": "testuser",
		"roles":    []string{"admin", "user"},
	}
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, extraClaims)
	require.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, "user123", userID)

		username, exists := c.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "testuser", username)

		roles, exists := c.Get("roles")
		assert.True(t, exists)
		assert.Equal(t, []string{"admin", "user"}, roles)

		tokenType, exists := c.Get("token_type")
		assert.True(t, exists)
		assert.Equal(t, "access", tokenType)

		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
		})
	})

	// 执行
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)
}
