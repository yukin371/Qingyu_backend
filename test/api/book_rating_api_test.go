package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readingAPI "Qingyu_backend/api/v1/reading"
	"Qingyu_backend/models/reading/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookRatingService 模拟评分服务
type MockBookRatingService struct {
	mock.Mock
}

func (m *MockBookRatingService) CreateRating(ctx context.Context, rating *bookstore.BookRating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockBookRatingService) GetRatingByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookRating, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookRating), args.Error(1)
}

func (m *MockBookRatingService) GetRatingsByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, bookID, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRatingService) GetRatingsByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRatingService) GetAverageRating(ctx context.Context, bookID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockBookRatingService) GetRatingDistribution(ctx context.Context, bookID primitive.ObjectID) (map[int]int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(map[int]int64), args.Error(1)
}

func (m *MockBookRatingService) UpdateRating(ctx context.Context, rating *bookstore.BookRating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockBookRatingService) DeleteRating(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRatingService) LikeRating(ctx context.Context, ratingID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, ratingID, userID)
	return args.Error(0)
}

func (m *MockBookRatingService) UnlikeRating(ctx context.Context, ratingID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, ratingID, userID)
	return args.Error(0)
}

func (m *MockBookRatingService) GetRatingByBookIDAndUserID(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error) {
	args := m.Called(ctx, bookID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookRating), args.Error(1)
}

func (m *MockBookRatingService) GetRatingsByRating(ctx context.Context, rating float64, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, rating, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRatingService) GetRatingsByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, tags, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRatingService) GetRatingCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRatingService) GetRatingStats(ctx context.Context, bookID primitive.ObjectID) (map[string]interface{}, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBookRatingService) GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstore.BookRating, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.BookRating), args.Error(1)
}

func (m *MockBookRatingService) GetRatingLikes(ctx context.Context, ratingID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, ratingID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRatingService) HasUserRated(ctx context.Context, bookID, userID primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, bookID, userID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockBookRatingService) GetUserRatingForBook(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error) {
	args := m.Called(ctx, bookID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookRating), args.Error(1)
}

func (m *MockBookRatingService) UpdateUserRating(ctx context.Context, bookID, userID primitive.ObjectID, rating float64, comment string, tags []string) error {
	args := m.Called(ctx, bookID, userID, rating, comment, tags)
	return args.Error(0)
}

func (m *MockBookRatingService) DeleteUserRating(ctx context.Context, bookID, userID primitive.ObjectID) error {
	args := m.Called(ctx, bookID, userID)
	return args.Error(0)
}

func (m *MockBookRatingService) BatchUpdateRatingTags(ctx context.Context, ratingIDs []primitive.ObjectID, tags []string) error {
	args := m.Called(ctx, ratingIDs, tags)
	return args.Error(0)
}

func (m *MockBookRatingService) BatchDeleteRatings(ctx context.Context, ratingIDs []primitive.ObjectID) error {
	args := m.Called(ctx, ratingIDs)
	return args.Error(0)
}

func (m *MockBookRatingService) BatchDeleteRatingsByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRatingService) BatchDeleteRatingsByUserID(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockBookRatingService) SearchRatings(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRatingService) GetRatingsWithComments(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, bookID, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRatingService) GetHighRatedComments(ctx context.Context, bookID primitive.ObjectID, minRating float64, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	args := m.Called(ctx, bookID, minRating, page, pageSize)
	return args.Get(0).([]*bookstore.BookRating), args.Get(1).(int64), args.Error(2)
}

// setupRatingTestRouter 设置评分测试路由
func setupRatingTestRouter(service bookstoreService.BookRatingService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := readingAPI.NewBookRatingAPI(service)

	v1 := router.Group("/api/v1/reading")
	{
		// 评分相关路由
		v1.GET("/ratings/:id", api.GetBookRating)
		v1.POST("/ratings", api.CreateRating)
		v1.PUT("/ratings/:id", api.UpdateRating)
		v1.DELETE("/ratings/:id", api.DeleteRating)
		v1.POST("/ratings/:id/like", api.LikeRating)
		v1.POST("/ratings/:id/unlike", api.UnlikeRating)
		v1.GET("/ratings/search", api.SearchRatings)

		// 图书评分路由
		v1.GET("/books/:book_id/ratings", api.GetRatingsByBookID)
		v1.GET("/books/:book_id/average-rating", api.GetAverageRating)
		v1.GET("/books/:book_id/rating-distribution", api.GetRatingDistribution)

		// 用户评分路由
		v1.GET("/users/:user_id/ratings", api.GetRatingsByUserID)
	}

	return router
}

// TestGetBookRating 测试获取评分详情
func TestGetBookRating(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	ratingID := primitive.NewObjectID()
	rating := &bookstore.BookRating{
		ID:      ratingID,
		BookID:  primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		Rating:  5,
		Comment: "很棒的书籍",
	}

	mockService.On("GetRatingByID", mock.Anything, ratingID).Return(rating, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/ratings/"+ratingID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetBookRating_InvalidID 测试无效的评分ID
func TestGetBookRating_InvalidID(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/ratings/invalid-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "无效的评分ID格式", response.Message)

	mockService.AssertNotCalled(t, "GetRatingByID")
}

// TestGetBookRating_NotFound 测试评分不存在
func TestGetBookRating_NotFound(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	ratingID := primitive.NewObjectID()
	mockService.On("GetRatingByID", mock.Anything, ratingID).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/reading/ratings/"+ratingID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// TestGetRatingsByBookID 测试获取图书的所有评分
func TestGetRatingsByBookID(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	bookID := primitive.NewObjectID()
	ratings := []*bookstore.BookRating{
		{
			ID:      primitive.NewObjectID(),
			BookID:  bookID,
			UserID:  primitive.NewObjectID(),
			Rating:  5,
			Comment: "很好",
		},
		{
			ID:      primitive.NewObjectID(),
			BookID:  bookID,
			UserID:  primitive.NewObjectID(),
			Rating:  4,
			Comment: "不错",
		},
	}

	mockService.On("GetRatingsByBookID", mock.Anything, bookID, 1, 10).Return(ratings, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/books/"+bookID.Hex()+"/ratings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取成功", response.Message)
	assert.Equal(t, int64(2), response.Total)

	mockService.AssertExpectations(t)
}

// TestGetRatingsByUserID 测试获取用户的所有评分
func TestGetRatingsByUserID(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	userID := primitive.NewObjectID()
	ratings := []*bookstore.BookRating{
		{
			ID:      primitive.NewObjectID(),
			BookID:  primitive.NewObjectID(),
			UserID:  userID,
			Rating:  5,
			Comment: "很好",
		},
	}

	mockService.On("GetRatingsByUserID", mock.Anything, userID, 1, 10).Return(ratings, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/users/"+userID.Hex()+"/ratings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetAverageRating 测试获取图书平均评分
func TestGetAverageRating(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	bookID := primitive.NewObjectID()
	avgRating := 4.5

	mockService.On("GetAverageRating", mock.Anything, bookID).Return(avgRating, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/books/"+bookID.Hex()+"/average-rating", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)
	assert.Equal(t, avgRating, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetRatingDistribution 测试获取评分分布
func TestGetRatingDistribution(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	bookID := primitive.NewObjectID()
	distribution := map[int]int64{
		5: 10,
		4: 5,
		3: 2,
		2: 1,
		1: 0,
	}

	mockService.On("GetRatingDistribution", mock.Anything, bookID).Return(distribution, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/books/"+bookID.Hex()+"/rating-distribution", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestCreateRating 测试创建评分
func TestCreateRating(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	mockService.On("CreateRating", mock.Anything, mock.AnythingOfType("*bookstore.BookRating")).Return(nil)

	requestBody := `{
		"book_id": "` + primitive.NewObjectID().Hex() + `",
		"user_id": "` + primitive.NewObjectID().Hex() + `",
		"rating": 5,
		"comment": "很棒的书籍"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/reading/ratings", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, "创建成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestCreateRating_InvalidJSON 测试创建评分时JSON格式错误
func TestCreateRating_InvalidJSON(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	requestBody := `{invalid json}`

	req, _ := http.NewRequest("POST", "/api/v1/reading/ratings", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "请求参数格式错误", response.Message)

	mockService.AssertNotCalled(t, "CreateRating")
}

// TestUpdateRating 测试更新评分
func TestUpdateRating(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	ratingID := primitive.NewObjectID()
	mockService.On("UpdateRating", mock.Anything, mock.AnythingOfType("*bookstore.BookRating")).Return(nil)

	requestBody := `{
		"rating": 4,
		"comment": "修改后的评价"
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/reading/ratings/"+ratingID.Hex(), strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "更新成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestDeleteRating 测试删除评分
func TestDeleteRating(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	ratingID := primitive.NewObjectID()
	mockService.On("DeleteRating", mock.Anything, ratingID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reading/ratings/"+ratingID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "删除成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestLikeRating 测试点赞评分
func TestLikeRating(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	// 设置测试路由，添加用户ID到上下文
	router = gin.New()
	gin.SetMode(gin.TestMode)

	api := readingAPI.NewBookRatingAPI(mockService)

	v1 := router.Group("/api/v1/reading")
	v1.Use(func(c *gin.Context) {
		// 模拟中间件设置用户ID
		c.Set("userID", primitive.NewObjectID().Hex())
		c.Next()
	})
	v1.POST("/ratings/:id/like", api.LikeRating)

	ratingID := primitive.NewObjectID()
	mockService.On("LikeRating", mock.Anything, ratingID, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reading/ratings/"+ratingID.Hex()+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "点赞成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestLikeRating_Unauthorized 测试未登录用户点赞
func TestLikeRating_Unauthorized(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	ratingID := primitive.NewObjectID()

	req, _ := http.NewRequest("POST", "/api/v1/reading/ratings/"+ratingID.Hex()+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "用户未登录", response.Message)

	mockService.AssertNotCalled(t, "LikeRating")
}

// TestUnlikeRating 测试取消点赞
func TestUnlikeRating(t *testing.T) {
	mockService := new(MockBookRatingService)

	// 设置测试路由，添加用户ID到上下文
	router := gin.New()
	gin.SetMode(gin.TestMode)

	api := readingAPI.NewBookRatingAPI(mockService)

	v1 := router.Group("/api/v1/reading")
	v1.Use(func(c *gin.Context) {
		// 模拟中间件设置用户ID
		c.Set("userID", primitive.NewObjectID().Hex())
		c.Next()
	})
	v1.POST("/ratings/:id/unlike", api.UnlikeRating)

	ratingID := primitive.NewObjectID()
	mockService.On("UnlikeRating", mock.Anything, ratingID, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reading/ratings/"+ratingID.Hex()+"/unlike", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "取消点赞成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestSearchRatings 测试搜索评分
func TestSearchRatings(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/ratings/search?keyword=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "搜索功能开发中", response.Message)
}

// TestGetRatingsByBookID_Pagination 测试分页功能
func TestGetRatingsByBookID_Pagination(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	bookID := primitive.NewObjectID()
	ratings := []*bookstore.BookRating{
		{ID: primitive.NewObjectID(), Rating: 5},
	}

	mockService.On("GetRatingsByBookID", mock.Anything, bookID, 2, 20).Return(ratings, int64(25), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/books/"+bookID.Hex()+"/ratings?page=2&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 20, response.Limit)

	mockService.AssertExpectations(t)
}

// TestCreateRating_ServiceError 测试服务层错误
func TestCreateRating_ServiceError(t *testing.T) {
	mockService := new(MockBookRatingService)
	router := setupRatingTestRouter(mockService)

	mockService.On("CreateRating", mock.Anything, mock.AnythingOfType("*bookstore.BookRating")).Return(assert.AnError)

	requestBody := `{
		"book_id": "` + primitive.NewObjectID().Hex() + `",
		"user_id": "` + primitive.NewObjectID().Hex() + `",
		"rating": 5,
		"comment": "测试评价"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/reading/ratings", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Contains(t, response.Message, "创建评分失败")

	mockService.AssertExpectations(t)
}
