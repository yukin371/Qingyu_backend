package social_test

import (
	"bytes"
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

// BatchLikeBooks 模拟批量点赞书籍
func (m *MockLikeService) BatchLikeBooks(ctx context.Context, userID string, bookIDs []string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, bookIDs)
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
		v1.POST("/books/batch-like", api.BatchLikeBooks)
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
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.Equal(t, "操作成功", response["message"]) // 成功响应message为"操作成功"

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
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.Equal(t, "操作成功", response["message"]) // 成功响应message为"操作成功"

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
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

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
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.Equal(t, "操作成功", response["message"]) // 成功响应message为"操作成功"

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
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.Equal(t, "操作成功", response["message"]) // 成功响应message为"操作成功"

	mockService.AssertExpectations(t)
}

// ========================= 批量点赞测试 =========================

// BatchLikeBookResult 批量点赞书籍结果
type BatchLikeBookResult struct {
	BookID    string `json:"book_id"`
	Success   bool   `json:"success"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// TestLikeAPI_BatchLikeBooks_Success 测试成功批量点赞书籍
func TestLikeAPI_BatchLikeBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	// 准备测试数据
	bookID1 := primitive.NewObjectID().Hex()
	bookID2 := primitive.NewObjectID().Hex()
	bookID3 := primitive.NewObjectID().Hex()

	// Mock服务层调用
	mockService.On("BatchLikeBooks", mock.Anything, userID, []string{bookID1, bookID2, bookID3}).Return(map[string]interface{}{
		"success_count": 3,
		"failed_count":  0,
		"results": []BatchLikeBookResult{
			{BookID: bookID1, Success: true},
			{BookID: bookID2, Success: true},
			{BookID: bookID3, Success: true},
		},
	}, nil)

	// 准备请求体
	requestBody := map[string][]string{
		"book_ids": {bookID1, bookID2, bookID3},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/batch-like", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["success_count"])
	assert.Equal(t, float64(0), data["failed_count"])

	mockService.AssertExpectations(t)
}

// TestLikeAPI_BatchLikeBooks_EmptyBookIDs 测试空书籍ID列表
func TestLikeAPI_BatchLikeBooks_EmptyBookIDs(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	// 准备空的请求体
	requestBody := map[string][]string{
		"book_ids": {},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/batch-like", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestLikeAPI_BatchLikeBooks_ExceedLimit 测试超过数量限制
func TestLikeAPI_BatchLikeBooks_ExceedLimit(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	// 准备超过50个书籍ID的请求体
	bookIDs := make([]string, 51)
	for i := 0; i < 51; i++ {
		bookIDs[i] = primitive.NewObjectID().Hex()
	}

	requestBody := map[string][]string{
		"book_ids": bookIDs,
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/batch-like", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestLikeAPI_BatchLikeBooks_PartialSuccess 测试部分成功
func TestLikeAPI_BatchLikeBooks_PartialSuccess(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	userID := primitive.NewObjectID().Hex()
	router := setupLikeTestRouter(mockService, userID)

	// 准备测试数据
	bookID1 := primitive.NewObjectID().Hex()
	bookID2 := primitive.NewObjectID().Hex()
	bookID3 := primitive.NewObjectID().Hex()

	// Mock服务层调用 - 部分成功
	mockService.On("BatchLikeBooks", mock.Anything, userID, []string{bookID1, bookID2, bookID3}).Return(map[string]interface{}{
		"success_count": 2,
		"failed_count":  1,
		"results": []BatchLikeBookResult{
			{BookID: bookID1, Success: true},
			{BookID: bookID2, Success: false, ErrorMsg: "书籍不存在"},
			{BookID: bookID3, Success: true},
		},
	}, nil)

	requestBody := map[string][]string{
		"book_ids": {bookID1, bookID2, bookID3},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/batch-like", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["success_count"])
	assert.Equal(t, float64(1), data["failed_count"])

	mockService.AssertExpectations(t)
}

// TestLikeAPI_BatchLikeBooks_Unauthorized 测试未授权访问
func TestLikeAPI_BatchLikeBooks_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockLikeService)
	router := setupLikeTestRouter(mockService, "") // 不设置userID

	bookID := primitive.NewObjectID().Hex()
	requestBody := map[string][]string{
		"book_ids": {bookID},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/batch-like", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
