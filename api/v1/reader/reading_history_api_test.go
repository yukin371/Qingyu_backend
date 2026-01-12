package reader

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

	readerModels "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
	"Qingyu_backend/service/interfaces"
	readerservice "Qingyu_backend/service/reader"
)

// MockReadingHistoryService 模拟ReadingHistoryService
type MockReadingHistoryService struct {
	mock.Mock
}

func (m *MockReadingHistoryService) RecordReading(ctx context.Context, userID, bookID, chapterID string, startTime, endTime time.Time, progress float64, deviceType, deviceID string) error {
	args := m.Called(ctx, userID, bookID, chapterID, startTime, endTime, progress, deviceType, deviceID)
	return args.Error(0)
}

func (m *MockReadingHistoryService) GetUserHistories(ctx context.Context, userID string, page, pageSize int) ([]*readerModels.ReadingHistory, *readerservice.PaginationInfo, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*readerservice.PaginationInfo), args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*readerModels.ReadingHistory), nil, args.Error(2)
	}
	return args.Get(0).([]*readerModels.ReadingHistory), args.Get(1).(*readerservice.PaginationInfo), args.Error(2)
}

func (m *MockReadingHistoryService) GetUserHistoriesByBook(ctx context.Context, userID, bookID string, page, pageSize int) ([]*readerModels.ReadingHistory, *readerservice.PaginationInfo, error) {
	args := m.Called(ctx, userID, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*readerservice.PaginationInfo), args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*readerModels.ReadingHistory), nil, args.Error(2)
	}
	return args.Get(0).([]*readerModels.ReadingHistory), args.Get(1).(*readerservice.PaginationInfo), args.Error(2)
}

func (m *MockReadingHistoryService) GetUserReadingStats(ctx context.Context, userID string) (*readerRepo.ReadingStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerRepo.ReadingStats), args.Error(1)
}

func (m *MockReadingHistoryService) GetUserDailyReadingStats(ctx context.Context, userID string, days int) ([]readerRepo.DailyReadingStats, error) {
	args := m.Called(ctx, userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]readerRepo.DailyReadingStats), args.Error(1)
}

func (m *MockReadingHistoryService) DeleteHistory(ctx context.Context, userID, historyID string) error {
	args := m.Called(ctx, userID, historyID)
	return args.Error(0)
}

func (m *MockReadingHistoryService) ClearUserHistories(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func setupReadingHistoryTestRouter(historyService interfaces.ReadingHistoryService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewReadingHistoryAPI(historyService)

	v1 := r.Group("/api/v1/reader/reading-history")
	{
		v1.POST("", api.RecordReading)
		v1.GET("", api.GetReadingHistories)
		v1.GET("/stats", api.GetReadingStats)
		v1.DELETE("/:id", api.DeleteHistory)
		v1.DELETE("", api.ClearHistories)
	}

	return r
}

func TestReadingHistoryAPI_RecordReading_Success(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupReadingHistoryTestRouter(mockService, userID)

	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()

	reqBody := map[string]interface{}{
		"book_id":     bookID,
		"chapter_id":  chapterID,
		"start_time":  startTime,
		"end_time":    endTime,
		"progress":    50.0,
		"device_type": "web",
		"device_id":   "test-device",
	}

	mockService.On("RecordReading", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything, mock.Anything, 50.0, "web", "test-device").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/reading-history", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 201, w.Code)
	mockService.AssertExpectations(t)
}

func TestReadingHistoryAPI_RecordReading_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	router := setupReadingHistoryTestRouter(mockService, "") // No userID

	reqBody := map[string]interface{}{
		"book_id":    primitive.NewObjectID().Hex(),
		"chapter_id": primitive.NewObjectID().Hex(),
		"start_time": time.Now(),
		"end_time":   time.Now(),
		"progress":   50.0,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/reading-history", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 401, w.Code)
}

func TestReadingHistoryAPI_GetReadingHistories_Success(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	userID := primitive.NewObjectID().Hex()
	router := setupReadingHistoryTestRouter(mockService, userID)

	expectedHistories := []*readerModels.ReadingHistory{}
	expectedPagination := &readerservice.PaginationInfo{
		Page:       1,
		PageSize:   20,
		Total:      100,
		TotalPages: 5,
	}

	mockService.On("GetUserHistories", mock.Anything, userID, 1, 20).
		Return(expectedHistories, expectedPagination, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/reading-history?page=1&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestReadingHistoryAPI_GetReadingHistories_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	router := setupReadingHistoryTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/reading-history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 401, w.Code)
}

func TestReadingHistoryAPI_GetReadingHistoriesByBook_Success(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupReadingHistoryTestRouter(mockService, userID)

	expectedHistories := []*readerModels.ReadingHistory{}
	expectedPagination := &readerservice.PaginationInfo{
		Page:       1,
		PageSize:   20,
		Total:      10,
		TotalPages: 1,
	}

	mockService.On("GetUserHistoriesByBook", mock.Anything, userID, bookID, 1, 20).
		Return(expectedHistories, expectedPagination, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/reading-history?book_id="+bookID+"&page=1&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestReadingHistoryAPI_GetReadingStats_Success(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	userID := primitive.NewObjectID().Hex()
	router := setupReadingHistoryTestRouter(mockService, userID)

	expectedStats := &readerRepo.ReadingStats{
		TotalDuration:    3600,
		TotalBooks:       10,
		TotalChapters:    100,
		LastReadTime:     time.Now(),
		AvgDailyDuration: 120,
	}

	expectedDailyStats := []readerRepo.DailyReadingStats{
		{Date: "2026-01-01", Duration: 3600, Books: 1},
		{Date: "2026-01-02", Duration: 7200, Books: 2},
	}

	mockService.On("GetUserReadingStats", mock.Anything, userID).Return(expectedStats, nil)
	mockService.On("GetUserDailyReadingStats", mock.Anything, userID, 30).Return(expectedDailyStats, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/reading-history/stats?days=30", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestReadingHistoryAPI_GetReadingStats_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	router := setupReadingHistoryTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/reading-history/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 401, w.Code)
}

func TestReadingHistoryAPI_DeleteHistory_Success(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	userID := primitive.NewObjectID().Hex()
	historyID := primitive.NewObjectID().Hex()
	router := setupReadingHistoryTestRouter(mockService, userID)

	mockService.On("DeleteHistory", mock.Anything, userID, historyID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/reading-history/"+historyID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestReadingHistoryAPI_DeleteHistory_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	router := setupReadingHistoryTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/reading-history/"+primitive.NewObjectID().Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 401, w.Code)
}

func TestReadingHistoryAPI_ClearHistories_Success(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	userID := primitive.NewObjectID().Hex()
	router := setupReadingHistoryTestRouter(mockService, userID)

	mockService.On("ClearUserHistories", mock.Anything, userID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/reading-history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestReadingHistoryAPI_ClearHistories_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReadingHistoryService)
	router := setupReadingHistoryTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/reading-history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 401, w.Code)
}
