package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Qingyu_backend/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionalJWTAuth_AllowsAnonymousRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(OptionalJWTAuth())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"exists":  exists,
			"user_id": userID,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"exists\":false")
}

func TestOptionalJWTAuth_PopulatesClaimsWhenTokenValid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret-key-123456",
			ExpirationHours: 24,
		},
	}
	defer func() {
		config.GlobalConfig = previousConfig
	}()

	token, err := GenerateToken("user-123", "tester", []string{"reader"})
	require.NoError(t, err)

	router := gin.New()
	router.Use(OptionalJWTAuth())
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user-123")
	assert.Contains(t, w.Body.String(), "tester")
}

func TestOptionalJWTAuth_IgnoresInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(OptionalJWTAuth())
	router.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"exists": exists,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"exists\":false")
}

func TestOptionalJWTAuth_IgnoresBlacklistedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret-key-123456",
			ExpirationHours: 24,
		},
	}
	defer func() {
		config.GlobalConfig = previousConfig
		SetSharedBlacklist(nil)
	}()

	token, err := GenerateToken("user-123", "tester", []string{"reader"})
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	require.NoError(t, blacklist.Add(context.Background(), token, time.Hour))
	SetSharedBlacklist(blacklist)

	router := gin.New()
	router.Use(OptionalJWTAuth())
	router.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"exists": exists,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"exists\":false")
}

func TestJWTAuth_RejectsBlacklistedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret-key-123456",
			ExpirationHours: 24,
		},
	}
	defer func() {
		config.GlobalConfig = previousConfig
		SetSharedBlacklist(nil)
	}()

	token, err := GenerateToken("user-123", "tester", []string{"reader"})
	require.NoError(t, err)

	blacklist := NewMockBlacklist()
	require.NoError(t, blacklist.Add(context.Background(), token, time.Hour))
	SetSharedBlacklist(blacklist)

	router := gin.New()
	router.Use(JWTAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "\"code\":\"2016\"")
}
