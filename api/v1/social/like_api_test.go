package social_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	socialAPI "Qingyu_backend/api/v1/social"
	"Qingyu_backend/models/social"
	"Qingyu_backend/service/interfaces"
)

// MockLikeService 模拟点赞服务接口
type MockLikeService struct {
	mock.Mock
}

// LikeBook 模拟点赞书籍
func (m *MockLikeService) LikeBook(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

// UnlikeBook 模拟取消点赞书籍
func (m *MockLikeService) UnlikeBook(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

// GetBookLikeCount 模拟获取书籍点赞数
func (m *MockLikeService) GetBookLikeCount(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

// IsBookLiked 模拟检查是否点赞书籍
func (m *MockLikeService) IsBookLiked(ctx context.Context, userID, bookID string) (bool, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Bool(0), args.Error(1)
}

// LikeComment 模拟点赞评论
func (m *MockLikeService) LikeComment(ctx context.Context, userID, commentID string) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}

// UnlikeComment 模拟取消点赞评论
func (m *MockLikeService) UnlikeComment(ctx context.Context, userID, commentID string) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}

// GetUserLikedBooks 模拟获取用户点赞的书籍列表
func (m *MockLikeService) GetUserLikedBooks(ctx context.Context, userID string, page, size int) ([]*social.Like, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Like), args.Get(1).(int64), args.Error(2)
}

// GetUserLikedComments 模拟获取用户点赞的评论列表
func (m *MockLikeService) GetUserLikedComments(ctx context.Context, userID string, page, size int) ([]*social.Like, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Like), args.Get(1).(int64), args.Error(2)
}

// GetBooksLikeCount 模拟批量获取书籍点赞数
func (m *MockLikeService) GetBooksLikeCount(ctx context.Context, bookIDs []string) (map[string]int64, error) {
	args := m.Called(ctx, bookIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

// GetUserLikeStatus 模拟获取用户点赞状态
func (m *MockLikeService) GetUserLikeStatus(ctx context.Context, userID string, bookIDs []string) (map[string]bool, error) {
	args := m.Called(ctx, userID, bookIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}

// GetUserLikeStats 模拟获取用户点赞统计
func (m *MockLikeService) GetUserLikeStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// setupLikeTestRouter 设置测试路由
func setupLikeTestRouter(likeService interfaces.LikeService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId（用于需要认证的端点）
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := socialAPI.NewLikeAPI(likeService)

	v1 := r.Group("/api/v1/reader")
	{
		v1.POST("/books/:bookId/like", api.LikeBook)
		v1.DELETE("/books/:bookId/like", api.UnlikeBook)
		v1.GET("/books/:bookId/like/info", api.GetBookLikeInfo)
		v1.GET("/likes/books", api.GetUserLikedBooks)
		v1.GET("/likes/stats", api.GetUserLikeStats)
		v1.POST("/comments/:id/like", api.LikeComment)
		v1.DELETE("/comments/:id/like", api.UnlikeComment)
	}

	return r
}

// TestLikeAPI_LikeBook_Success 测试成功点赞书籍
func TestLikeAPI_LikeBook_Success(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	mockService.On("LikeBook", mock.Anything, userID, bookID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "点赞成功")

	mockService.AssertExpectations(t)
}

// TestLikeAPI_LikeBook_MissingBookID 测试缺少书籍ID
func TestLikeAPI_LikeBook_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books//like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestLikeAPI_LikeBook_Unauthorized 测试未授权访问
func TestLikeAPI_LikeBook_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	router := setupLikeTestRouter(mockService, "") // 不设置userID

	bookID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLikeAPI_UnlikeBook_Success 测试成功取消点赞书籍
func TestLikeAPI_UnlikeBook_Success(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	mockService.On("UnlikeBook", mock.Anything, userID, bookID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/books/"+bookID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "取消点赞成功")

	mockService.AssertExpectations(t)
}

// TestLikeAPI_GetBookLikeInfo_Success 测试成功获取书籍点赞信息
func TestLikeAPI_GetBookLikeInfo_Success(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	mockService.On("GetBookLikeCount", mock.Anything, bookID).Return(int64(100), nil)
	mockService.On("IsBookLiked", mock.Anything, userID, bookID).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/like/info", nil)

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
	assert.Equal(t, float64(100), data["like_count"])
	assert.Equal(t, true, data["is_liked"])

	mockService.AssertExpectations(t)
}

// TestLikeAPI_GetBookLikeInfo_NotLoggedIn 测试未登录用户获取书籍点赞信息
func TestLikeAPI_GetBookLikeInfo_NotLoggedIn(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	router := setupLikeTestRouter(mockService, "") // 不设置userID

	bookID := primitive.NewObjectID().Hex()
	mockService.On("GetBookLikeCount", mock.Anything, bookID).Return(int64(50), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/like/info", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(50), data["like_count"])
	assert.Equal(t, false, data["is_liked"])

	mockService.AssertExpectations(t)
}

// TestLikeAPI_LikeComment_Success 测试成功点赞评论
func TestLikeAPI_LikeComment_Success(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	commentID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	mockService.On("LikeComment", mock.Anything, userID, commentID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+commentID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "点赞成功")

	mockService.AssertExpectations(t)
}

// TestLikeAPI_UnlikeComment_Success 测试成功取消点赞评论
func TestLikeAPI_UnlikeComment_Success(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	commentID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	mockService.On("UnlikeComment", mock.Anything, userID, commentID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "取消点赞成功")

	mockService.AssertExpectations(t)
}
