package social_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	socialAPI "Qingyu_backend/api/v1/social"
	"Qingyu_backend/models/social"
)

// MockRatingService 模拟评分服务接口
type MockRatingService struct {
	mock.Mock
}

// GetRatingStats 模拟获取评分统计
func (m *MockRatingService) GetRatingStats(ctx context.Context, targetType, targetID string) (*social.RatingStats, error) {
	args := m.Called(ctx, targetType, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.RatingStats), args.Error(1)
}

// GetUserRating 模拟获取用户评分
func (m *MockRatingService) GetUserRating(ctx context.Context, userID, targetType, targetID string) (int, error) {
	args := m.Called(ctx, userID, targetType, targetID)
	return args.Int(0), args.Error(1)
}

// AggregateRatings 模拟聚合评分
func (m *MockRatingService) AggregateRatings(ctx context.Context, targetType, targetID string) (*social.RatingStats, error) {
	args := m.Called(ctx, targetType, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.RatingStats), args.Error(1)
}

// InvalidateCache 模拟使缓存失效
func (m *MockRatingService) InvalidateCache(ctx context.Context, targetType, targetID string) error {
	args := m.Called(ctx, targetType, targetID)
	return args.Error(0)
}

// setupRatingTestRouter 设置评分测试路由
func setupRatingTestRouter(ratingService *MockRatingService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置user_id（用于需要认证的端点）
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := socialAPI.NewRatingAPI(ratingService)

	v1 := r.Group("/api/v1/ratings")
	{
		v1.GET("/:targetType/:targetId/stats", api.GetRatingStats)
		v1.GET("/:targetType/:targetId/user-rating", api.GetUserRating)
	}

	return r
}

// TestGetRatingStats_Success 测试成功获取评分统计
func TestGetRatingStats_Success(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "")

	targetType := "comment"
	targetID := "507f1f77bcf86cd799439011"

	expectedStats := &social.RatingStats{
		TargetID:      targetID,
		TargetType:    targetType,
		AverageRating: 4.5,
		TotalRatings:  100,
		Distribution: map[int]int64{
			1: 5,
			2: 3,
			3: 10,
			4: 30,
			5: 52,
		},
		UpdatedAt: time.Now(),
	}

	mockService.On("GetRatingStats", mock.Anything, targetType, targetID).Return(expectedStats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, targetID, data["targetId"])
	assert.Equal(t, targetType, data["targetType"])
	assert.Equal(t, 4.5, data["averageRating"])
	assert.Equal(t, float64(100), data["totalRatings"])

	mockService.AssertExpectations(t)
}

// TestGetRatingStats_InvalidTargetType 测试无效的目标类型
func TestGetRatingStats_InvalidTargetType(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "")

	targetID := "507f1f77bcf86cd799439011"

	req, _ := http.NewRequest("GET", "/api/v1/ratings/invalid/"+targetID+"/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "参数错误", response["message"])
}

// TestGetRatingStats_EmptyTargetType 测试空的目标类型
func TestGetRatingStats_EmptyTargetType(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "")

	targetID := "507f1f77bcf86cd799439011"

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetID+"/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// 由于路由不匹配，会返回404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestGetRatingStats_ServiceError 测试服务层错误
func TestGetRatingStats_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "")

	targetType := "comment"
	targetID := "507f1f77bcf86cd799439011"

	mockService.On("GetRatingStats", mock.Anything, targetType, targetID).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "服务器内部错误")

	mockService.AssertExpectations(t)
}

// TestGetRatingStats_ForReviewType 测试获取书评评分统计
func TestGetRatingStats_ForReviewType(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "")

	targetType := "review"
	targetID := "507f1f77bcf86cd799439011"

	expectedStats := &social.RatingStats{
		TargetID:      targetID,
		TargetType:    targetType,
		AverageRating: 4.2,
		TotalRatings:  50,
		Distribution: map[int]int64{
			1: 2,
			2: 1,
			3: 8,
			4: 15,
			5: 24,
		},
		UpdatedAt: time.Now(),
	}

	mockService.On("GetRatingStats", mock.Anything, targetType, targetID).Return(expectedStats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestGetRatingStats_ForBookType 测试获取书籍评分统计
func TestGetRatingStats_ForBookType(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "")

	targetType := "book"
	targetID := "507f1f77bcf86cd799439011"

	expectedStats := &social.RatingStats{
		TargetID:      targetID,
		TargetType:    targetType,
		AverageRating: 4.8,
		TotalRatings:  200,
		Distribution: map[int]int64{
			1: 1,
			2: 0,
			3: 5,
			4: 20,
			5: 174,
		},
		UpdatedAt: time.Now(),
	}

	mockService.On("GetRatingStats", mock.Anything, targetType, targetID).Return(expectedStats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestGetUserRating_Success 测试成功获取用户评分
func TestGetUserRating_Success(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	userID := "507f1f77bcf86cd799439011"
	router := setupRatingTestRouter(mockService, userID)

	targetType := "book"
	targetID := "507f1f77bcf86cd799439012"
	expectedRating := 5

	mockService.On("GetUserRating", mock.Anything, userID, targetType, targetID).Return(expectedRating, nil)

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/user-rating", nil)

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
	assert.Equal(t, float64(expectedRating), data["rating"])

	mockService.AssertExpectations(t)
}

// TestGetUserRating_Unauthorized 测试未授权访问
func TestGetUserRating_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	router := setupRatingTestRouter(mockService, "") // 不设置userID

	targetType := "book"
	targetID := "507f1f77bcf86cd799439012"

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/user-rating", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "登录")
}

// TestGetUserRating_InvalidTargetType 测试无效的目标类型
func TestGetUserRating_InvalidTargetType(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	userID := "507f1f77bcf86cd799439011"
	router := setupRatingTestRouter(mockService, userID)

	targetType := "comment" // comment不是GetUserRating的有效类型
	targetID := "507f1f77bcf86cd799439012"

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/user-rating", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "参数错误", response["message"])
}

// TestGetUserRating_ServiceError 测试服务层错误
func TestGetUserRating_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	userID := "507f1f77bcf86cd799439011"
	router := setupRatingTestRouter(mockService, userID)

	targetType := "book"
	targetID := "507f1f77bcf86cd799439012"

	mockService.On("GetUserRating", mock.Anything, userID, targetType, targetID).Return(0, errors.New("database error"))

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/user-rating", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "服务器内部错误")

	mockService.AssertExpectations(t)
}

// TestGetUserRating_ForReviewType 测试获取用户对书评的评分
func TestGetUserRating_ForReviewType(t *testing.T) {
	// Given
	mockService := new(MockRatingService)
	userID := "507f1f77bcf86cd799439011"
	router := setupRatingTestRouter(mockService, userID)

	targetType := "review"
	targetID := "507f1f77bcf86cd799439012"
	expectedRating := 4

	mockService.On("GetUserRating", mock.Anything, userID, targetType, targetID).Return(expectedRating, nil)

	req, _ := http.NewRequest("GET", "/api/v1/ratings/"+targetType+"/"+targetID+"/user-rating", nil)

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
	assert.Equal(t, float64(expectedRating), data["rating"])

	mockService.AssertExpectations(t)
}
