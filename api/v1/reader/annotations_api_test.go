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

// MockReaderServiceForAnnotations 模拟ReaderService
type MockReaderServiceForAnnotations struct {
	mock.Mock
}

func (m *MockReaderServiceForAnnotations) CreateAnnotation(ctx context.Context, annotation *readerModels.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	args := m.Called(ctx, annotationID, updates)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) DeleteAnnotation(ctx context.Context, annotationID string) error {
	args := m.Called(ctx, annotationID)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetNotes(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

// Implement remaining required methods from ReaderService interface
func (m *MockReaderServiceForAnnotations) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderServiceForAnnotations) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*readerModels.ReadingProgress, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderServiceForAnnotations) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) GetReadingSettings(ctx context.Context, userID string) (*readerModels.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.ReadingSettings), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) SaveReadingSettings(ctx context.Context, settings *readerModels.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) BatchCreateAnnotations(ctx context.Context, annotations []*readerModels.Annotation) error {
	args := m.Called(ctx, annotations)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	args := m.Called(ctx, annotationIDs)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReaderServiceForAnnotations) BatchUpdateBookStatus(ctx context.Context, userID string, bookIDs []string, status string) error {
	args := m.Called(ctx, userID, bookIDs, status)
	return args.Error(0)
}

func (m *MockReaderServiceForAnnotations) UpdateBookStatus(ctx context.Context, userID, bookID, status string) error {
	args := m.Called(ctx, userID, bookID, status)
	return args.Error(0)
}

func setupAnnotationsTestRouter(readerService interfaces.ReaderService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewAnnotationsAPI(readerService)

	v1 := r.Group("/api/v1/reader/annotations")
	{
		v1.POST("", api.CreateAnnotation)
		v1.PUT("/:id", api.UpdateAnnotation)
		v1.DELETE("/:id", api.DeleteAnnotation)
		v1.GET("/chapter", api.GetAnnotationsByChapter)
		v1.GET("/book", api.GetAnnotationsByBook)
		v1.GET("/notes", api.GetNotes)
		v1.GET("/notes/search", api.SearchNotes)
		v1.GET("/bookmarks", api.GetBookmarks)
		v1.GET("/bookmarks/latest", api.GetLatestBookmark)
		v1.GET("/highlights", api.GetHighlights)
		v1.GET("/recent", api.GetRecentAnnotations)
		v1.GET("/public", api.GetPublicAnnotations)
	}

	return r
}

func TestAnnotationsAPI_CreateAnnotation_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"bookId":    bookID,
		"chapterId": chapterID,
		"type":      "note",
		"text":      "这是一段重要文本",
		"note":      "我的笔记",
		"range":     "0-100",
	}

	mockService.On("CreateAnnotation", mock.Anything, mock.Anything).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/annotations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAnnotationsAPI_CreateAnnotation_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	router := setupAnnotationsTestRouter(mockService, "") // No userID

	reqBody := map[string]interface{}{
		"bookId":    primitive.NewObjectID().Hex(),
		"chapterId": primitive.NewObjectID().Hex(),
		"type":      "note",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/annotations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAnnotationsAPI_UpdateAnnotation_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	annotationID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"text": strPtr("更新后的文本"),
		"note": strPtr("更新后的笔记"),
	}

	mockService.On("UpdateAnnotation", mock.Anything, annotationID, mock.Anything).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/annotations/"+annotationID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_DeleteAnnotation_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	annotationID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	mockService.On("DeleteAnnotation", mock.Anything, annotationID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/annotations/"+annotationID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetAnnotationsByChapter_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	userIDObj := primitive.NewObjectID()
	bookID := primitive.NewObjectID().Hex()
	bookIDObj := primitive.NewObjectID()
	chapterID := primitive.NewObjectID().Hex()
	chapterIDObj := primitive.NewObjectID()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedAnnotations := []*readerModels.Annotation{
		{ID: primitive.NewObjectID(), UserID: userIDObj, BookID: bookIDObj, ChapterID: chapterIDObj},
	}

	mockService.On("GetAnnotationsByChapter", mock.Anything, userID, bookID, chapterID).
		Return(expectedAnnotations, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/chapter?bookId="+bookID+"&chapterId="+chapterID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetAnnotationsByChapter_MissingParams(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/chapter", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnnotationsAPI_GetAnnotationsByBook_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedAnnotations := []*readerModels.Annotation{}

	mockService.On("GetAnnotationsByBook", mock.Anything, userID, bookID).
		Return(expectedAnnotations, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/book?bookId="+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetNotes_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedNotes := []*readerModels.Annotation{}

	mockService.On("GetNotes", mock.Anything, userID, bookID).Return(expectedNotes, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/notes?bookId="+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_SearchNotes_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedNotes := []*readerModels.Annotation{}

	mockService.On("SearchNotes", mock.Anything, userID, "测试").Return(expectedNotes, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/notes/search?keyword=测试", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_SearchNotes_EmptyKeyword(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/notes/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnnotationsAPI_GetBookmarks_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedBookmarks := []*readerModels.Annotation{}

	mockService.On("GetBookmarks", mock.Anything, userID, bookID).Return(expectedBookmarks, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/bookmarks?bookId="+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetLatestBookmark_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedBookmark := &readerModels.Annotation{ID: primitive.NewObjectID()}

	mockService.On("GetLatestBookmark", mock.Anything, userID, bookID).Return(expectedBookmark, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/bookmarks/latest?bookId="+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetHighlights_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedHighlights := []*readerModels.Annotation{}

	mockService.On("GetHighlights", mock.Anything, userID, bookID).Return(expectedHighlights, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/highlights?bookId="+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetRecentAnnotations_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedAnnotations := []*readerModels.Annotation{}

	mockService.On("GetRecentAnnotations", mock.Anything, userID, 20).Return(expectedAnnotations, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/recent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetRecentAnnotations_CustomLimit(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	userID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, userID)

	expectedAnnotations := []*readerModels.Annotation{}

	mockService.On("GetRecentAnnotations", mock.Anything, userID, 50).Return(expectedAnnotations, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/recent?limit=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetPublicAnnotations_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupAnnotationsTestRouter(mockService, "")

	expectedAnnotations := []*readerModels.Annotation{}

	mockService.On("GetPublicAnnotations", mock.Anything, bookID, chapterID).Return(expectedAnnotations, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/public?bookId="+bookID+"&chapterId="+chapterID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnnotationsAPI_GetPublicAnnotations_MissingParams(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForAnnotations)
	router := setupAnnotationsTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/annotations/public", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Helper function
func strPtr(s string) *string {
	return &s
}
