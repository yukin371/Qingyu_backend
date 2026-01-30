package social_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	socialAPI "Qingyu_backend/api/v1/social"
	"Qingyu_backend/models/social"
	"Qingyu_backend/service/interfaces"
)

// MockFollowService 模拟关注服务接口
type MockFollowService struct {
	mock.Mock
}

// FollowUser 模拟关注用户
func (m *MockFollowService) FollowUser(ctx context.Context, followerID, followingID string) error {
	args := m.Called(ctx, followerID, followingID)
	return args.Error(0)
}

// UnfollowUser 模拟取消关注用户
func (m *MockFollowService) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	args := m.Called(ctx, followerID, followingID)
	return args.Error(0)
}

// CheckFollowStatus 模拟检查关注状态
func (m *MockFollowService) CheckFollowStatus(ctx context.Context, followerID, followingID string) (bool, error) {
	args := m.Called(ctx, followerID, followingID)
	return args.Bool(0), args.Error(1)
}

// GetFollowers 模拟获取粉丝列表
func (m *MockFollowService) GetFollowers(ctx context.Context, userID string, page, size int) ([]*social.FollowInfo, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.FollowInfo), args.Get(1).(int64), args.Error(2)
}

// GetFollowing 模拟获取关注列表
func (m *MockFollowService) GetFollowing(ctx context.Context, userID string, page, size int) ([]*social.FollowingInfo, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.FollowingInfo), args.Get(1).(int64), args.Error(2)
}

// GetFollowStats 模拟获取关注统计
func (m *MockFollowService) GetFollowStats(ctx context.Context, userID string) (*social.FollowStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.FollowStats), args.Error(1)
}

// FollowAuthor 模拟关注作者
func (m *MockFollowService) FollowAuthor(ctx context.Context, userID, authorID, authorName, authorAvatar string, notifyNewBook bool) error {
	args := m.Called(ctx, userID, authorID, authorName, authorAvatar, notifyNewBook)
	return args.Error(0)
}

// UnfollowAuthor 模拟取消关注作者
func (m *MockFollowService) UnfollowAuthor(ctx context.Context, userID, authorID string) error {
	args := m.Called(ctx, userID, authorID)
	return args.Error(0)
}

// GetFollowingAuthors 模拟获取关注的作者列表
func (m *MockFollowService) GetFollowingAuthors(ctx context.Context, userID string, page, size int) ([]*social.AuthorFollow, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.AuthorFollow), args.Get(1).(int64), args.Error(2)
}

// setupFollowTestRouter 设置测试路由
func setupFollowTestRouter(followService interfaces.FollowService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId（用于需要认证的端点）
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("username", "testuser")
			c.Set("user_avatar", "")
		}
		c.Next()
	})

	api := socialAPI.NewFollowAPI(followService)

	v1 := r.Group("/api/v1/social")
	{
		// 用户关注
		v1.POST("/users/:userId/follow", api.FollowUser)
		v1.DELETE("/users/:userId/unfollow", api.UnfollowUser)
		v1.GET("/users/:userId/followers", api.GetFollowers)
		v1.GET("/users/:userId/following", api.GetFollowing)
		v1.GET("/users/:userId/follow-status", api.CheckFollowStatus)

		// 作者关注
		v1.POST("/authors/:authorId/follow", api.FollowAuthor)
		v1.DELETE("/authors/:authorId/unfollow", api.UnfollowAuthor)
		v1.GET("/following/authors", api.GetFollowingAuthors)
	}

	return r
}

// =========================
// 用户关注测试
// =========================

// TestFollowAPI_FollowUser_Success 测试关注用户成功
func TestFollowAPI_FollowUser_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	mockService.On("FollowUser", mock.Anything, userID, targetUserID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/social/users/"+targetUserID+"/follow", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestFollowAPI_FollowUser_EmptyUserID 测试关注用户-用户ID为空
func TestFollowAPI_FollowUser_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	req, _ := http.NewRequest("POST", "/api/v1/social/users//follow", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 参数错误code为1001
}

// TestFollowAPI_UnfollowUser_Success 测试取消关注用户成功
func TestFollowAPI_UnfollowUser_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	mockService.On("UnfollowUser", mock.Anything, userID, targetUserID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/users/"+targetUserID+"/unfollow", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestFollowAPI_GetFollowers_Success 测试获取粉丝列表成功
func TestFollowAPI_GetFollowers_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, "")

	expectedFollowers := []*social.FollowInfo{
		{
			FollowerID:   "user1",
			FollowerName: "粉丝1",
			IsMutual:     false,
			CreatedAt:    time.Now(),
		},
		{
			FollowerID:   "user2",
			FollowerName: "粉丝2",
			IsMutual:     true,
			CreatedAt:    time.Now(),
		},
	}

	mockService.On("GetFollowers", mock.Anything, userID, 1, 20).
		Return(expectedFollowers, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+userID+"/followers?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestFollowAPI_GetFollowing_Success 测试获取关注列表成功
func TestFollowAPI_GetFollowing_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, "")

	expectedFollowing := []*social.FollowingInfo{
		{
			FollowingID:   "user1",
			FollowingName: "关注用户1",
			IsMutual:      true,
			CreatedAt:     time.Now(),
		},
		{
			FollowingID:   "user2",
			FollowingName: "关注用户2",
			IsMutual:      false,
			CreatedAt:     time.Now(),
		},
	}

	mockService.On("GetFollowing", mock.Anything, userID, 1, 20).
		Return(expectedFollowing, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+userID+"/following?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestFollowAPI_CheckFollowStatus_True 测试检查关注状态-已关注
func TestFollowAPI_CheckFollowStatus_True(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	mockService.On("CheckFollowStatus", mock.Anything, userID, targetUserID).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+targetUserID+"/follow-status", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, true, data["is_following"])

	mockService.AssertExpectations(t)
}

// TestFollowAPI_CheckFollowStatus_False 测试检查关注状态-未关注
func TestFollowAPI_CheckFollowStatus_False(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	mockService.On("CheckFollowStatus", mock.Anything, userID, targetUserID).Return(false, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+targetUserID+"/follow-status", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, false, data["is_following"])

	mockService.AssertExpectations(t)
}

// =========================
// 作者关注测试
// =========================

// TestFollowAPI_FollowAuthor_Success 测试关注作者成功
func TestFollowAPI_FollowAuthor_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	authorID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	mockService.On("FollowAuthor", mock.Anything, userID, authorID, "作者名", "", true).
		Return(nil)

	reqBody := map[string]interface{}{
		"author_name":     "作者名",
		"notify_new_book": true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/authors/"+authorID+"/follow", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestFollowAPI_FollowAuthor_MissingAuthorName 测试关注作者缺少作者名
func TestFollowAPI_FollowAuthor_MissingAuthorName(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	authorID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"notify_new_book": true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/authors/"+authorID+"/follow", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 参数错误code为1001
}

// TestFollowAPI_UnfollowAuthor_Success 测试取消关注作者成功
func TestFollowAPI_UnfollowAuthor_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	authorID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	mockService.On("UnfollowAuthor", mock.Anything, userID, authorID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/authors/"+authorID+"/unfollow", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestFollowAPI_GetFollowingAuthors_Success 测试获取关注的作者列表成功
func TestFollowAPI_GetFollowingAuthors_Success(t *testing.T) {
	// Given
	mockService := new(MockFollowService)
	userID := primitive.NewObjectID().Hex()
	router := setupFollowTestRouter(mockService, userID)

	expectedAuthors := []*social.AuthorFollow{
		{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			AuthorID:      "author1",
			AuthorName:    "作者1",
			NotifyNewBook: true,
			CreatedAt:     time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			AuthorID:      "author2",
			AuthorName:    "作者2",
			NotifyNewBook: false,
			CreatedAt:     time.Now(),
		},
	}

	mockService.On("GetFollowingAuthors", mock.Anything, userID, 1, 20).
		Return(expectedAuthors, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/following/authors?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}
