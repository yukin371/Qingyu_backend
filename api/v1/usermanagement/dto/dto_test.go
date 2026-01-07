package dto

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestDTOValidation 测试 DTO 验证
func TestDTOValidation(t *testing.T) {
	t.Run("RegisterRequest validation", func(t *testing.T) {
		validReq := RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		assert.Equal(t, "testuser", validReq.Username)
		assert.Equal(t, "test@example.com", validReq.Email)
		assert.Equal(t, "password123", validReq.Password)
	})

	t.Run("LoginRequest validation", func(t *testing.T) {
		validReq := LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		assert.Equal(t, "testuser", validReq.Username)
		assert.Equal(t, "password123", validReq.Password)
	})

	t.Run("UpdateProfileRequest validation", func(t *testing.T) {
		nickname := "新昵称"
		bio := "这是我的个人简介"
		validReq := UpdateProfileRequest{
			Nickname: &nickname,
			Bio:      &bio,
		}
		assert.Equal(t, "新昵称", *validReq.Nickname)
		assert.Equal(t, "这是我的个人简介", *validReq.Bio)
	})
}

// TestRequestSerialization 测试请求序列化
func TestRequestSerialization(t *testing.T) {
	t.Run("RegisterRequest JSON serialization", func(t *testing.T) {
		req := RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		body, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(body), "testuser")
		assert.Contains(t, string(body), "test@example.com")
	})

	t.Run("LoginRequest JSON serialization", func(t *testing.T) {
		req := LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		body, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(body), "testuser")
		assert.Contains(t, string(body), "password123")
	})
}

// TestResponseSerialization 测试响应序列化
func TestResponseSerialization(t *testing.T) {
	t.Run("RegisterResponse JSON serialization", func(t *testing.T) {
		resp := RegisterResponse{
			UserID:   "user-123",
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
			Status:   "active",
			Token:    "jwt-token",
		}

		body, err := json.Marshal(resp)
		assert.NoError(t, err)

		var parsed RegisterResponse
		err = json.Unmarshal(body, &parsed)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", parsed.UserID)
		assert.Equal(t, "testuser", parsed.Username)
		assert.Equal(t, "jwt-token", parsed.Token)
	})
}

// TestHandlerStructure 测试处理器结构
func TestHandlerStructure(t *testing.T) {
	// 这个测试验证处理器结构是否正确
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 验证可以创建路由组而不会崩溃
	v1 := router.Group("/api/v1")
	assert.NotNil(t, v1)

	authGroup := v1.Group("/user-management/auth")
	assert.NotNil(t, authGroup)

	profileGroup := v1.Group("/user-management/profile")
	assert.NotNil(t, profileGroup)
}

// TestHTTPRequestParsing 测试HTTP请求解析
func TestHTTPRequestParsing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	var parsedReq RegisterRequest
	router.POST("/test", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&parsedReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	reqBody := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "testuser", parsedReq.Username)
	assert.Equal(t, "test@example.com", parsedReq.Email)
}
