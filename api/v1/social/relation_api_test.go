package social_test

import (
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
)

// MockUserRelationService 模拟用户关系服务接口
type MockUserRelationService struct {
	mock.Mock
}

// FollowUser 模拟关注用户
func (m *MockUserRelationService) FollowUser(ctx interface{}, followerID, followeeID string) error {
	args := m.Called(ctx, followerID, followeeID)
	return args.Error(0)
}

// UnfollowUser 模拟取消关注用户
func (m *MockUserRelationService) UnfollowUser(ctx interface{}, followerID, followeeID string) error {
	args := m.Called(ctx, followerID, followeeID)
	return args.Error(0)
}

// IsFollowing 模拟检查是否关注
func (m *MockUserRelationService) IsFollowing(ctx interface{}, followerID, followeeID string) (bool, error) {
	args := m.Called(ctx, followerID, followeeID)
	return args.Bool(0), args.Error(1)
}

// GetFollowers 模拟获取粉丝列表
func (m *MockUserRelationService) GetFollowers(ctx interface{}, userID string, page, pageSize int) ([]*social.UserRelation, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.UserRelation), args.Get(1).(int64), args.Error(2)
}

// GetFollowing 模拟获取关注列表
func (m *MockUserRelationService) GetFollowing(ctx interface{}, userID string, page, pageSize int) ([]*social.UserRelation, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.UserRelation), args.Get(1).(int64), args.Error(2)
}

// GetFollowerCount 模拟获取粉丝数
func (m *MockUserRelationService) GetFollowerCount(ctx interface{}, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// GetFollowingCount 模拟获取关注数
func (m *MockUserRelationService) GetFollowingCount(ctx interface{}, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupRelationTestRouter 设置测试路由
func setupRelationTestRouter(relationService socialAPI.UserRelationServiceInterface, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId（用于需要认证的端点）
	// 注意：relation_api使用"userId"而不是"user_id"
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := socialAPI.NewUserRelationAPI(relationService)

	v1 := r.Group("/api/v1/social")
	{
		// 关注/取消关注
		v1.POST("/follow/:userId", api.FollowUser)
		v1.DELETE("/follow/:userId", api.UnfollowUser)
		v1.GET("/follow/:userId/status", api.CheckIsFollowing)

		// 粉丝/关注列表
		v1.GET("/users/:userId/followers", api.GetFollowers)
		v1.GET("/users/:userId/following", api.GetFollowing)
		v1.GET("/users/:userId/follow-stats", api.GetFollowStats)
	}

	return r
}

// TestRelationAPI_FollowUser_Success 测试关注用户成功
func TestRelationAPI_FollowUser_Success(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	mockService.On("FollowUser", mock.Anything, userID, targetUserID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/social/follow/"+targetUserID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	// 检查响应
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, true, data["followed"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_FollowUser_EmptyUserID 测试关注用户-用户ID为空
func TestRelationAPI_FollowUser_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	req, _ := http.NewRequest("POST", "/api/v1/social/follow/", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRelationAPI_FollowUser_AlreadyFollowing 测试关注用户-已经关注
func TestRelationAPI_FollowUser_AlreadyFollowing(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	// 返回"已经关注了该用户"错误
	mockService.On("FollowUser", mock.Anything, userID, targetUserID).
		Return(assert.AnError)

	req, _ := http.NewRequest("POST", "/api/v1/social/follow/"+targetUserID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// 会返回500而不是409，因为我们使用的是通用错误
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// TestRelationAPI_UnfollowUser_Success 测试取消关注用户成功
func TestRelationAPI_UnfollowUser_Success(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	mockService.On("UnfollowUser", mock.Anything, userID, targetUserID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/follow/"+targetUserID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, false, data["followed"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_UnfollowUser_NotFollowing 测试取消关注用户-未关注
func TestRelationAPI_UnfollowUser_NotFollowing(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	mockService.On("UnfollowUser", mock.Anything, userID, targetUserID).
		Return(assert.AnError)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/follow/"+targetUserID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// TestRelationAPI_CheckIsFollowing_True 测试检查关注状态-已关注
func TestRelationAPI_CheckIsFollowing_True(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	mockService.On("IsFollowing", mock.Anything, userID, targetUserID).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/follow/"+targetUserID+"/status", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, true, data["is_following"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_CheckIsFollowing_False 测试检查关注状态-未关注
func TestRelationAPI_CheckIsFollowing_False(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, userID)

	mockService.On("IsFollowing", mock.Anything, userID, targetUserID).Return(false, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/follow/"+targetUserID+"/status", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, false, data["is_following"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_GetFollowers_Success 测试获取粉丝列表成功
func TestRelationAPI_GetFollowers_Success(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, "")

	expectedRelations := []*social.UserRelation{
		func() *social.UserRelation {
			r := &social.UserRelation{
				FollowerID: "user1",
				FolloweeID: userID,
				Status:     social.RelationStatusActive,
				Timestamps: social.Timestamps{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			r.ID = primitive.NewObjectID()
			return r
		}(),
		func() *social.UserRelation {
			r := &social.UserRelation{
				FollowerID: "user2",
				FolloweeID: userID,
				Status:     social.RelationStatusActive,
				Timestamps: social.Timestamps{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			r.ID = primitive.NewObjectID()
			return r
		}(),
	}

	mockService.On("GetFollowers", mock.Anything, userID, 1, 20).
		Return(expectedRelations, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+userID+"/followers?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["followers"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_GetFollowers_EmptyUserID 测试获取粉丝列表-用户ID为空
func TestRelationAPI_GetFollowers_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	router := setupRelationTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/social/users//followers", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRelationAPI_GetFollowing_Success 测试获取关注列表成功
func TestRelationAPI_GetFollowing_Success(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, "")

	expectedRelations := []*social.UserRelation{
		{
			IdentifiedEntity: social.IdentifiedEntity{
				ID: primitive.NewObjectID().Hex(),
			},
			FollowerID: userID,
			FolloweeID: "user1",
			Status:     social.RelationStatusActive,
			Timestamps: social.Timestamps{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			IdentifiedEntity: social.IdentifiedEntity{
				ID: primitive.NewObjectID().Hex(),
			},
			FollowerID: userID,
			FolloweeID: "user2",
			Status:     social.RelationStatusActive,
			Timestamps: social.Timestamps{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	mockService.On("GetFollowing", mock.Anything, userID, 1, 20).
		Return(expectedRelations, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+userID+"/following?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["following"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_GetFollowing_DefaultPagination 测试获取关注列表-默认分页
func TestRelationAPI_GetFollowing_DefaultPagination(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, "")

	mockService.On("GetFollowing", mock.Anything, userID, 1, 20).
		Return([]*social.UserRelation{}, int64(0), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+userID+"/following", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestRelationAPI_GetFollowStats_Success 测试获取关注统计成功
func TestRelationAPI_GetFollowStats_Success(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	userID := primitive.NewObjectID().Hex()
	router := setupRelationTestRouter(mockService, "")

	mockService.On("GetFollowerCount", mock.Anything, userID).Return(int64(100), nil)
	mockService.On("GetFollowingCount", mock.Anything, userID).Return(int64(50), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/users/"+userID+"/follow-stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(100), data["follower_count"])
	assert.Equal(t, float64(50), data["following_count"])

	mockService.AssertExpectations(t)
}

// TestRelationAPI_GetFollowStats_EmptyUserID 测试获取关注统计-用户ID为空
func TestRelationAPI_GetFollowStats_EmptyUserID(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	router := setupRelationTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/social/users//follow-stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRelationAPI_FollowUser_Unauthorized 测试关注用户-未授权
func TestRelationAPI_FollowUser_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	router := setupRelationTestRouter(mockService, "") // 不设置用户ID

	targetUserID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/api/v1/social/follow/"+targetUserID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRelationAPI_CheckIsFollowing_Unauthorized 测试检查关注状态-未授权
func TestRelationAPI_CheckIsFollowing_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	router := setupRelationTestRouter(mockService, "") // 不设置用户ID

	targetUserID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("GET", "/api/v1/social/follow/"+targetUserID+"/status", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRelationAPI_UnfollowUser_Unauthorized 测试取消关注用户-未授权
func TestRelationAPI_UnfollowUser_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockUserRelationService)
	router := setupRelationTestRouter(mockService, "") // 不设置用户ID

	targetUserID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("DELETE", "/api/v1/social/follow/"+targetUserID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
