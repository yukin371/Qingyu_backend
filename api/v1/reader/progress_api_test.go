package reader_test

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

	progressAPI "Qingyu_backend/api/v1/reader"
	"Qingyu_backend/models/reader"
	"Qingyu_backend/service/interfaces"
)

// MockReaderService 模拟阅读器服务接口
type MockReaderService struct {
	mock.Mock
}

// GetReadingProgress 模拟获取阅读进度
func (m *MockReaderService) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingProgress), args.Error(1)
}

// SaveReadingProgress 模拟保存阅读进度
func (m *MockReaderService) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

// UpdateReadingTime 模拟更新阅读时长
func (m *MockReaderService) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

// GetRecentReading 模拟获取最近阅读记录
func (m *MockReaderService) GetRecentReading(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

// GetReadingHistory 模拟获取阅读历史
func (m *MockReaderService) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*reader.ReadingProgress, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Get(1).(int64), args.Error(2)
}

// GetTotalReadingTime 模拟获取总阅读时长
func (m *MockReaderService) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// GetReadingTimeByPeriod 模拟获取某时段阅读时长
func (m *MockReaderService) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

// GetUnfinishedBooks 模拟获取未读完的书
func (m *MockReaderService) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

// GetFinishedBooks 模拟获取已读完的书
func (m *MockReaderService) GetFinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

// DeleteReadingProgress 模拟删除阅读进度
func (m *MockReaderService) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

// 实现其他接口方法（简化版本，仅用于测试）
func (m *MockReaderService) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockReaderService) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	args := m.Called(ctx, chapterID)
	return args.Get(0), args.Error(1)
}

func (m *MockReaderService) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderService) CreateAnnotation(ctx context.Context, annotation *reader.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockReaderService) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	args := m.Called(ctx, annotationID, updates)
	return args.Error(0)
}

func (m *MockReaderService) DeleteAnnotation(ctx context.Context, annotationID string) error {
	args := m.Called(ctx, annotationID)
	return args.Error(0)
}

func (m *MockReaderService) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) SearchNotes(ctx context.Context, userID, keyword string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderService) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReaderService) GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReaderService) SaveReadingSettings(ctx context.Context, settings *reader.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockReaderService) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockReaderService) BatchCreateAnnotations(ctx context.Context, annotations []*reader.Annotation) error {
	args := m.Called(ctx, annotations)
	return args.Error(0)
}

func (m *MockReaderService) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	args := m.Called(ctx, annotationIDs)
	return args.Error(0)
}

func (m *MockReaderService) SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// setupProgressTestRouter 设置测试路由
func setupProgressTestRouter(readerService interfaces.ReaderService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := progressAPI.NewProgressAPI(readerService)

	v1 := r.Group("/api/v1/reader/progress")
	{
		v1.GET("/:bookId", api.GetReadingProgress)
		v1.POST("", api.SaveReadingProgress)
		v1.PUT("/reading-time", api.UpdateReadingTime)
		v1.GET("/recent", api.GetRecentReading)
		v1.GET("/history", api.GetReadingHistory)
		v1.GET("/stats", api.GetReadingStats)
		v1.GET("/unfinished", api.GetUnfinishedBooks)
		v1.GET("/finished", api.GetFinishedBooks)
	}

	return r
}

// =========================
// 阅读进度测试
// =========================

// TestProgressAPI_GetReadingProgress_Success 测试获取阅读进度成功
func TestProgressAPI_GetReadingProgress_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	expectedProgress := &reader.ReadingProgress{
		ID:         primitive.NewObjectID().Hex(),
		UserID:     userID,
		BookID:     bookID,
		ChapterID:  primitive.NewObjectID().Hex(),
		Progress:   0.5,
		LastReadAt: time.Now(),
	}

	mockService.On("GetReadingProgress", mock.Anything, userID, bookID).
		Return(expectedProgress, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/"+bookID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_SaveReadingProgress_Success 测试保存阅读进度成功
func TestProgressAPI_SaveReadingProgress_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	mockService.On("SaveReadingProgress", mock.Anything, userID, bookID, chapterID, 0.75).
		Return(nil)

	reqBody := map[string]interface{}{
		"bookId":    bookID,
		"chapterId": chapterID,
		"progress":  0.75,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_SaveReadingProgress_MissingBookID 测试保存阅读进度-缺少bookId
func TestProgressAPI_SaveReadingProgress_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"chapterId": primitive.NewObjectID().Hex(),
		"progress":  0.75,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestProgressAPI_SaveReadingProgress_InvalidProgress 测试保存阅读进度-进度超出范围
func TestProgressAPI_SaveReadingProgress_InvalidProgress(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"bookId":    primitive.NewObjectID().Hex(),
		"chapterId": primitive.NewObjectID().Hex(),
		"progress":  1.5, // 超出范围
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/progress", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestProgressAPI_UpdateReadingTime_Success 测试更新阅读时长成功
func TestProgressAPI_UpdateReadingTime_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	mockService.On("UpdateReadingTime", mock.Anything, userID, bookID, int64(3600)).
		Return(nil)

	reqBody := map[string]interface{}{
		"bookId":   bookID,
		"duration": 3600,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/progress/reading-time", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_UpdateReadingTime_MissingBookID 测试更新阅读时长-缺少bookId
func TestProgressAPI_UpdateReadingTime_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"duration": 3600,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/progress/reading-time", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestProgressAPI_GetRecentReading_Success 测试获取最近阅读记录成功
func TestProgressAPI_GetRecentReading_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	expectedProgresses := []*reader.ReadingProgress{
		{
			UserID: userID,
			BookID: primitive.NewObjectID().Hex(),
		},
	}

	mockService.On("GetRecentReading", mock.Anything, userID, 20).
		Return(expectedProgresses, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/recent", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_GetRecentReading_CustomLimit 测试获取最近阅读记录-自定义限制
func TestProgressAPI_GetRecentReading_CustomLimit(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	mockService.On("GetRecentReading", mock.Anything, userID, 10).
		Return([]*reader.ReadingProgress{}, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/recent?limit=10", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// TestProgressAPI_GetReadingHistory_Success 测试获取阅读历史成功
func TestProgressAPI_GetReadingHistory_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	expectedProgresses := []*reader.ReadingProgress{
		{
			UserID: userID,
			BookID: primitive.NewObjectID().Hex(),
		},
	}

	mockService.On("GetReadingHistory", mock.Anything, userID, 1, 20).
		Return(expectedProgresses, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/history?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_GetReadingStats_Success 测试获取阅读统计成功
func TestProgressAPI_GetReadingStats_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	mockService.On("GetTotalReadingTime", mock.Anything, userID).
		Return(int64(7200), nil)
	mockService.On("GetUnfinishedBooks", mock.Anything, userID).
		Return([]*reader.ReadingProgress{}, nil)
	mockService.On("GetFinishedBooks", mock.Anything, userID).
		Return([]*reader.ReadingProgress{}, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_GetUnfinishedBooks_Success 测试获取未读完的书成功
func TestProgressAPI_GetUnfinishedBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	expectedProgresses := []*reader.ReadingProgress{
		{
			UserID:   userID,
			BookID:   primitive.NewObjectID().Hex(),
			Progress: 0.3,
		},
	}

	mockService.On("GetUnfinishedBooks", mock.Anything, userID).
		Return(expectedProgresses, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/unfinished", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_GetFinishedBooks_Success 测试获取已读完的书成功
func TestProgressAPI_GetFinishedBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	userID := primitive.NewObjectID().Hex()
	router := setupProgressTestRouter(mockService, userID)

	expectedProgresses := []*reader.ReadingProgress{
		{
			UserID:   userID,
			BookID:   primitive.NewObjectID().Hex(),
			Progress: 1.0,
		},
	}

	mockService.On("GetFinishedBooks", mock.Anything, userID).
		Return(expectedProgresses, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/finished", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProgressAPI_Unauthorized 测试未授权访问
func TestProgressAPI_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderService)
	router := setupProgressTestRouter(mockService, "") // 不设置用户ID

	bookID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("GET", "/api/v1/reader/progress/"+bookID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
