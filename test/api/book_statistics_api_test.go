package api

import (
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

	readingAPI "Qingyu_backend/api/v1/reading"
	"Qingyu_backend/models/reading/bookstore"
	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookStatisticsService 模拟统计服务
type MockBookStatisticsService struct {
	mock.Mock
}

func (m *MockBookStatisticsService) GetStatisticsByBookID(ctx context.Context, bookID primitive.ObjectID) (*bookstore.BookStatistics, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) GetTopViewedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) GetTopFavoritedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) GetHottestBooks(ctx context.Context, limit int) ([]*bookstore.BookStatistics, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) GetTrendingBooks(ctx context.Context, days int, limit int) ([]*bookstore.BookStatistics, error) {
	args := m.Called(ctx, days, limit)
	return args.Get(0).([]*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) IncrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBookStatisticsService) GetStatisticsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*bookstore.BookStatistics, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).([]*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) GenerateDailyReport(ctx context.Context, date time.Time) (map[string]interface{}, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBookStatisticsService) GenerateWeeklyReport(ctx context.Context, startDate time.Time) (map[string]interface{}, error) {
	args := m.Called(ctx, startDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBookStatisticsService) GenerateMonthlyReport(ctx context.Context, year int, month int) (map[string]interface{}, error) {
	args := m.Called(ctx, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBookStatisticsService) SearchStatistics(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookStatistics, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	return args.Get(0).([]*bookstore.BookStatistics), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookStatisticsService) CreateStatistics(ctx context.Context, stats *bookstore.BookStatistics) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockBookStatisticsService) GetStatisticsByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookStatistics, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStatistics), args.Error(1)
}

func (m *MockBookStatisticsService) UpdateStatistics(ctx context.Context, stats *bookstore.BookStatistics) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockBookStatisticsService) DeleteStatistics(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookStatisticsService) DecrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	args := m.Called(ctx, bookID, rating)
	return args.Error(0)
}

func (m *MockBookStatisticsService) RemoveRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	args := m.Called(ctx, bookID, rating)
	return args.Error(0)
}

func (m *MockBookStatisticsService) RecalculateRating(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) UpdateHotScore(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookStatisticsService) BatchUpdateHotScore(ctx context.Context, bookIDs []primitive.ObjectID) error {
	args := m.Called(ctx, bookIDs)
	return args.Error(0)
}

func (m *MockBookStatisticsService) RecalculateAllHotScores(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBookStatisticsService) GetBookPopularityLevel(ctx context.Context, bookID primitive.ObjectID) (string, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockBookStatisticsService) BatchUpdateViewCount(ctx context.Context, bookIDs []primitive.ObjectID, increment int64) error {
	args := m.Called(ctx, bookIDs, increment)
	return args.Error(0)
}

func (m *MockBookStatisticsService) BatchCreateStatistics(ctx context.Context, statsList []*bookstore.BookStatistics) error {
	args := m.Called(ctx, statsList)
	return args.Error(0)
}

func (m *MockBookStatisticsService) BatchDeleteStatistics(ctx context.Context, bookIDs []primitive.ObjectID) error {
	args := m.Called(ctx, bookIDs)
	return args.Error(0)
}

func (m *MockBookStatisticsService) GetStatisticsByFilter(ctx context.Context, filter *BookstoreRepo.BookStatisticsFilter, page, pageSize int) ([]*bookstore.BookStatistics, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	return args.Get(0).([]*bookstore.BookStatistics), args.Get(1).(int64), args.Error(2)
}

// setupStatisticsTestRouter 设置统计测试路由
func setupStatisticsTestRouter(service bookstoreService.BookStatisticsService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := readingAPI.NewBookStatisticsAPI(service)

	v1 := router.Group("/api/v1/reading")
	{
		// 统计相关路由
		v1.GET("/books/:book_id/statistics", api.GetBookStatistics)
		v1.GET("/statistics/top-viewed", api.GetTopViewedBooks)
		v1.GET("/statistics/top-favorited", api.GetTopFavoritedBooks)
		v1.GET("/statistics/top-rated", api.GetTopRatedBooks)
		v1.GET("/statistics/hottest", api.GetHottestBooks)
		v1.GET("/statistics/trending", api.GetTrendingBooks)
		v1.POST("/books/:book_id/view", api.IncrementViewCount)
		v1.POST("/books/:book_id/favorite", api.IncrementFavoriteCount)
		v1.GET("/statistics/aggregated", api.GetAggregatedStatistics)
		v1.GET("/statistics/time-range", api.GetStatisticsByTimeRange)
		v1.GET("/statistics/daily-report", api.GetDailyStatisticsReport)
		v1.GET("/statistics/weekly-report", api.GetWeeklyStatisticsReport)
		v1.GET("/statistics/monthly-report", api.GetMonthlyStatisticsReport)
		v1.GET("/statistics/search", api.SearchStatistics)
	}

	return router
}

// TestGetBookStatistics 测试获取图书统计信息
func TestGetBookStatistics(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	bookID := primitive.NewObjectID()
	stats := &bookstore.BookStatistics{
		ID:            primitive.NewObjectID(),
		BookID:        bookID,
		ViewCount:     1000,
		FavoriteCount: 50,
		CommentCount:  20,
		ShareCount:    10,
		AverageRating: 4.5,
	}

	mockService.On("GetStatisticsByBookID", mock.Anything, bookID).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/books/"+bookID.Hex()+"/statistics", nil)
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

// TestGetBookStatistics_InvalidID 测试无效的图书ID
func TestGetBookStatistics_InvalidID(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/books/invalid-id/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "无效的图书ID格式", response.Message)

	mockService.AssertNotCalled(t, "GetStatisticsByBookID")
}

// TestGetTopViewedBooks 测试获取最多浏览的图书
func TestGetTopViewedBooks(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			ViewCount: 10000,
		},
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			ViewCount: 8000,
		},
	}

	mockService.On("GetTopViewedBooks", mock.Anything, 10).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/top-viewed", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetTopViewedBooks_WithLimit 测试带限制的最多浏览图书
func TestGetTopViewedBooks_WithLimit(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{ID: primitive.NewObjectID(), ViewCount: 10000},
	}

	mockService.On("GetTopViewedBooks", mock.Anything, 20).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/top-viewed?limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetTopFavoritedBooks 测试获取最多收藏的图书
func TestGetTopFavoritedBooks(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{
			ID:            primitive.NewObjectID(),
			BookID:        primitive.NewObjectID(),
			FavoriteCount: 500,
		},
	}

	mockService.On("GetTopFavoritedBooks", mock.Anything, 10).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/top-favorited", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetTopRatedBooks 测试获取最高评分的图书
func TestGetTopRatedBooks(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{
			ID:            primitive.NewObjectID(),
			BookID:        primitive.NewObjectID(),
			AverageRating: 4.9,
		},
	}

	mockService.On("GetTopRatedBooks", mock.Anything, 10).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/top-rated", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetHottestBooks 测试获取最热门的图书
func TestGetHottestBooks(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{
			ID:       primitive.NewObjectID(),
			BookID:   primitive.NewObjectID(),
			HotScore: 95.5,
		},
	}

	mockService.On("GetHottestBooks", mock.Anything, 10).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/hottest", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetTrendingBooks 测试获取趋势图书
func TestGetTrendingBooks(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{ID: primitive.NewObjectID()},
	}

	mockService.On("GetTrendingBooks", mock.Anything, 0, 10).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/trending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestIncrementViewCount 测试增加浏览量
func TestIncrementViewCount(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	bookID := primitive.NewObjectID()
	mockService.On("IncrementViewCount", mock.Anything, bookID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reading/books/"+bookID.Hex()+"/view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "增加浏览量成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestIncrementViewCount_InvalidID 测试无效ID增加浏览量
func TestIncrementViewCount_InvalidID(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("POST", "/api/v1/reading/books/invalid-id/view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "IncrementViewCount")
}

// TestIncrementFavoriteCount 测试增加收藏量
func TestIncrementFavoriteCount(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	bookID := primitive.NewObjectID()
	mockService.On("IncrementFavoriteCount", mock.Anything, bookID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reading/books/"+bookID.Hex()+"/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "增加收藏量成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetAggregatedStatistics 测试获取聚合统计信息
func TestGetAggregatedStatistics(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	aggregatedStats := map[string]interface{}{
		"total_books":     1000,
		"total_views":     50000,
		"total_favorites": 10000,
		"total_comments":  5000,
	}

	mockService.On("GetAggregatedStatistics", mock.Anything).Return(aggregatedStats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/aggregated", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetStatisticsByTimeRange 测试获取时间范围内的统计信息
func TestGetStatisticsByTimeRange(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{ID: primitive.NewObjectID()},
	}

	mockService.On("GetStatisticsByTimeRange", mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/time-range?start_time=2024-01-01T00:00:00Z&end_time=2024-01-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetStatisticsByTimeRange_InvalidTime 测试无效时间格式
func TestGetStatisticsByTimeRange_InvalidTime(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/time-range?start_time=invalid&end_time=2024-01-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "开始时间格式错误", response.Message)

	mockService.AssertNotCalled(t, "GetStatisticsByTimeRange")
}

// TestGetDailyStatisticsReport 测试获取日统计报告
func TestGetDailyStatisticsReport(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	report := map[string]interface{}{
		"date":            "2024-01-15",
		"total_views":     5000,
		"total_favorites": 500,
	}

	mockService.On("GenerateDailyReport", mock.Anything, mock.AnythingOfType("time.Time")).Return(report, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/daily-report?date=2024-01-15", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetDailyStatisticsReport_InvalidDate 测试无效日期格式
func TestGetDailyStatisticsReport_InvalidDate(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/daily-report?date=invalid-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "日期格式错误", response.Message)

	mockService.AssertNotCalled(t, "GenerateDailyReport")
}

// TestGetWeeklyStatisticsReport 测试获取周统计报告
func TestGetWeeklyStatisticsReport(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	report := map[string]interface{}{
		"year": 2024,
		"week": 3,
	}

	mockService.On("GenerateWeeklyReport", mock.Anything, mock.AnythingOfType("time.Time")).Return(report, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/weekly-report?year=2024&week=3", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetWeeklyStatisticsReport_MissingParams 测试缺少参数
func TestGetWeeklyStatisticsReport_MissingParams(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/weekly-report?year=2024", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "年份和周数不能为空", response.Message)

	mockService.AssertNotCalled(t, "GenerateWeeklyReport")
}

// TestGetMonthlyStatisticsReport 测试获取月统计报告
func TestGetMonthlyStatisticsReport(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	report := map[string]interface{}{
		"year":  2024,
		"month": 1,
	}

	mockService.On("GenerateMonthlyReport", mock.Anything, 2024, 1).Return(report, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/monthly-report?year=2024&month=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestSearchStatistics 测试搜索统计信息
func TestSearchStatistics(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	stats := []*bookstore.BookStatistics{
		{ID: primitive.NewObjectID()},
	}

	mockService.On("SearchStatistics", mock.Anything, "test", 1, 10).Return(stats, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/search?keyword=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "搜索成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestSearchStatistics_EmptyKeyword 测试空关键词搜索
func TestSearchStatistics_EmptyKeyword(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/search?keyword=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "搜索关键词不能为空", response.Message)

	mockService.AssertNotCalled(t, "SearchStatistics")
}

// TestGetTopViewedBooks_ServiceError 测试服务层错误
func TestGetTopViewedBooks_ServiceError(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	mockService.On("GetTopViewedBooks", mock.Anything, 10).Return([]*bookstore.BookStatistics(nil), assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/top-viewed", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Contains(t, response.Message, "获取最多浏览图书失败")

	mockService.AssertExpectations(t)
}

// TestIncrementViewCount_ServiceError 测试增加浏览量服务错误
func TestIncrementViewCount_ServiceError(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	bookID := primitive.NewObjectID()
	mockService.On("IncrementViewCount", mock.Anything, bookID).Return(assert.AnError)

	req, _ := http.NewRequest("POST", "/api/v1/reading/books/"+bookID.Hex()+"/view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)

	mockService.AssertExpectations(t)
}

// TestGetStatisticsByTimeRange_MissingParams 测试缺少时间参数
func TestGetStatisticsByTimeRange_MissingParams(t *testing.T) {
	mockService := new(MockBookStatisticsService)
	router := setupStatisticsTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/reading/statistics/time-range?start_time=2024-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "开始时间和结束时间不能为空", response.Message)

	mockService.AssertNotCalled(t, "GetStatisticsByTimeRange")
}
