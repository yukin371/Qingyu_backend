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
	"Qingyu_backend/service/interfaces"
)

// MockReaderServiceForSettings 模拟ReaderService (用于SettingAPI测试)
type MockReaderServiceForSettings struct {
	mock.Mock
}

func (m *MockReaderServiceForSettings) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderServiceForSettings) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForSettings) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*readerModels.ReadingProgress, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderServiceForSettings) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForSettings) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) CreateAnnotation(ctx context.Context, annotation *readerModels.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	args := m.Called(ctx, annotationID, updates)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) DeleteAnnotation(ctx context.Context, annotationID string) error {
	args := m.Called(ctx, annotationID)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetNotes(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetReadingSettings(ctx context.Context, userID string) (*readerModels.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.ReadingSettings), args.Error(1)
}

func (m *MockReaderServiceForSettings) SaveReadingSettings(ctx context.Context, settings *readerModels.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) BatchCreateAnnotations(ctx context.Context, annotations []*readerModels.Annotation) error {
	args := m.Called(ctx, annotations)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	args := m.Called(ctx, annotationIDs)
	return args.Error(0)
}

func (m *MockReaderServiceForSettings) SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReaderServiceForSettings) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func setupSettingTestRouter(readerService interfaces.ReaderService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewSettingAPI(readerService)

	v1 := r.Group("/api/v1/reader/settings")
	{
		v1.GET("", api.GetReadingSettings)
		v1.POST("", api.SaveReadingSettings)
		v1.PUT("", api.UpdateReadingSettings)
	}

	return r
}

func TestSettingAPI_GetReadingSettings_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForSettings)
	userID := primitive.NewObjectID().Hex()
	router := setupSettingTestRouter(mockService, userID)

	expectedSettings := &readerModels.ReadingSettings{
		UserID:     userID,
		FontSize:   18,
		FontFamily: "Arial",
		LineHeight: 1.6,
	}

	mockService.On("GetReadingSettings", mock.Anything, mock.AnythingOfType("string")).Return(expectedSettings, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/settings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSettingAPI_GetReadingSettings_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForSettings)
	router := setupSettingTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/settings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSettingAPI_SaveReadingSettings_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForSettings)
	userID := primitive.NewObjectID().Hex()
	router := setupSettingTestRouter(mockService, userID)

	fontSize := 18
	fontFamily := "Arial"

	reqBody := map[string]interface{}{
		"fontSize":   &fontSize,
		"fontFamily": &fontFamily,
		"lineHeight": 1.6,
	}

	mockService.On("SaveReadingSettings", mock.Anything, mock.Anything).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/settings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSettingAPI_SaveReadingSettings_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForSettings)
	router := setupSettingTestRouter(mockService, "")

	reqBody := map[string]interface{}{
		"fontSize": 18,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/settings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSettingAPI_UpdateReadingSettings_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForSettings)
	userID := primitive.NewObjectID().Hex()
	router := setupSettingTestRouter(mockService, userID)

	fontSize := 20
	fontFamily := "Georgia"

	reqBody := map[string]interface{}{
		"fontSize":   &fontSize,
		"fontFamily": &fontFamily,
	}

	mockService.On("UpdateReadingSettings", mock.Anything, userID, mock.Anything).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/settings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSettingAPI_UpdateReadingSettings_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForSettings)
	router := setupSettingTestRouter(mockService, "")

	lineHeight := 1.8
	reqBody := map[string]interface{}{
		"lineHeight": &lineHeight,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/settings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
