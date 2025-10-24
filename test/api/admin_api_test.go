//go:build integration
// +build integration

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	adminAPI "Qingyu_backend/api/v1/admin"
	usersModel "Qingyu_backend/models/users"
	userService "Qingyu_backend/service/user"
)

// TestAdminUserAPI_ListUsers 测试管理员获取用户列表
func TestAdminUserAPI_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建Mock Repository
	mockRepo := new(MockUserRepository)

	// 准备测试数据
	users := []*usersModel.User{
		{
			ID:       "user1",
			Username: "testuser1",
			Email:    "test1@example.com",
			Role:     "user",
			Status:   usersModel.UserStatusActive,
		},
		{
			ID:       "user2",
			Username: "testuser2",
			Email:    "test2@example.com",
			Role:     "author",
			Status:   usersModel.UserStatusActive,
		},
	}

	// 设置Mock返回 - List方法接受Filter参数
	mockRepo.On("List", mock.Anything, mock.AnythingOfType("infrastructure.Filter")).Return(users, nil)
	mockRepo.On("Count", mock.Anything, mock.AnythingOfType("infrastructure.Filter")).Return(int64(2), nil)

	// 创建Service和API
	userSvc := userService.NewUserService(mockRepo)
	api := adminAPI.NewUserAdminAPI(userSvc)

	// 设置路由
	router := gin.New()
	router.Use(mockAdminAuthMiddleware())
	router.GET("/api/v1/admin/users", api.ListUsers)

	// 创建请求
	req, _ := http.NewRequest("GET", "/api/v1/admin/users?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	mockRepo.AssertExpectations(t)
}

// TestAdminUserAPI_GetUser 测试管理员获取单个用户
func TestAdminUserAPI_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)

	user := &usersModel.User{
		ID:       "user1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
		Status:   usersModel.UserStatusActive,
	}

	mockRepo.On("GetByID", mock.Anything, "user1").Return(user, nil)

	userSvc := userService.NewUserService(mockRepo)
	api := adminAPI.NewUserAdminAPI(userSvc)

	router := gin.New()
	router.Use(mockAdminAuthMiddleware())
	router.GET("/api/v1/admin/users/:id", api.GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/admin/users/user1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	mockRepo.AssertExpectations(t)
}

// TestAdminUserAPI_UpdateUser 测试管理员更新用户
func TestAdminUserAPI_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)

	user := &usersModel.User{
		ID:       "user1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
		Status:   usersModel.UserStatusActive,
	}

	mockRepo.On("GetByID", mock.Anything, "user1").Return(user, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)

	userSvc := userService.NewUserService(mockRepo)
	api := adminAPI.NewUserAdminAPI(userSvc)

	router := gin.New()
	router.Use(mockAdminAuthMiddleware())
	router.PUT("/api/v1/admin/users/:id", api.UpdateUser)

	updateData := map[string]interface{}{
		"nickname": "Updated Name",
		"role":     "author",
	}
	body, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/user1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

// TestAdminUserAPI_DeleteUser 测试管理员删除用户
func TestAdminUserAPI_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)

	mockRepo.On("Delete", mock.Anything, "user1").Return(nil)

	userSvc := userService.NewUserService(mockRepo)
	api := adminAPI.NewUserAdminAPI(userSvc)

	router := gin.New()
	router.Use(mockAdminAuthMiddleware())
	router.DELETE("/api/v1/admin/users/:id", api.DeleteUser)

	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/user1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

// TestAdminUserAPI_BanUser 测试管理员封禁用户
func TestAdminUserAPI_BanUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)

	user := &usersModel.User{
		ID:       "user1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
		Status:   usersModel.UserStatusActive,
	}

	mockRepo.On("GetByID", mock.Anything, "user1").Return(user, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)

	userSvc := userService.NewUserService(mockRepo)
	api := adminAPI.NewUserAdminAPI(userSvc)

	router := gin.New()
	router.Use(mockAdminAuthMiddleware())
	router.POST("/api/v1/admin/users/:id/ban", api.BanUser)

	banData := map[string]interface{}{
		"reason":       "违反规则",
		"duration":     7,
		"durationUnit": "days",
	}
	body, _ := json.Marshal(banData)

	req, _ := http.NewRequest("POST", "/api/v1/admin/users/user1/ban", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

// mockAdminAuthMiddleware 模拟管理员认证中间件
func mockAdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟管理员用户
		c.Set("user_id", "admin1")
		c.Set("username", "admin")
		c.Set("role", "admin")
		c.Next()
	}
}
