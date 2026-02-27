package reader

import (
	"context"
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

// MockReaderServiceForBooks 模拟ReaderService (用于BooksAPI测试)
type MockReaderServiceForBooks struct {
	mock.Mock
}

func (m *MockReaderServiceForBooks) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderServiceForBooks) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForBooks) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*readerModels.ReadingProgress, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderServiceForBooks) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModels.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.ReadingProgress), args.Error(1)
}

func (m *MockReaderServiceForBooks) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) CreateAnnotation(ctx context.Context, annotation *readerModels.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	args := m.Called(ctx, annotationID, updates)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) DeleteAnnotation(ctx context.Context, annotationID string) error {
	args := m.Called(ctx, annotationID)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetNotes(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModels.Annotation, error) {
	args := m.Called(ctx, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModels.Annotation), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetReadingSettings(ctx context.Context, userID string) (*readerModels.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.ReadingSettings), args.Error(1)
}

func (m *MockReaderServiceForBooks) SaveReadingSettings(ctx context.Context, settings *readerModels.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) BatchCreateAnnotations(ctx context.Context, annotations []*readerModels.Annotation) error {
	args := m.Called(ctx, annotations)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	args := m.Called(ctx, annotationIDs)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReaderServiceForBooks) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReaderServiceForBooks) BatchUpdateBookStatus(ctx context.Context, userID string, bookIDs []string, status string) error {
	args := m.Called(ctx, userID, bookIDs, status)
	return args.Error(0)
}

func (m *MockReaderServiceForBooks) UpdateBookStatus(ctx context.Context, userID, bookID, status string) error {
	args := m.Called(ctx, userID, bookID, status)
	return args.Error(0)
}

func setupBooksTestRouter(readerService interfaces.ReaderService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 错误处理中间件
	r.Use(func(c *gin.Context) {
		c.Next()
		// 检查是否有错误
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": c.Errors.String(),
			})
		}
	})

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewBooksAPI(readerService)

	v1 := r.Group("/api/v1/reader/books")
	{
		v1.GET("", api.GetBookshelf)
		v1.POST("/:bookId", api.AddToBookshelf)
		v1.DELETE("/:bookId", api.RemoveFromBookshelf)
		v1.GET("/recent", api.GetRecentReading)
		v1.GET("/unfinished", api.GetUnfinishedBooks)
		v1.GET("/finished", api.GetFinishedBooks)
		v1.PUT("/:bookId/status", api.UpdateBookStatus)
	}

	return r
}

func TestBooksAPI_GetBookshelf_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	userIDObj, _ := primitive.ObjectIDFromHex(userID)
	router := setupBooksTestRouter(mockService, userID)

	expectedProgresses := []*readerModels.ReadingProgress{
		{UserID: userIDObj, BookID: primitive.NewObjectID()},
	}

	mockService.On("GetReadingHistory", mock.Anything, userID, 1, 20).
		Return(expectedProgresses, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooksAPI_GetBookshelf_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	router := setupBooksTestRouter(mockService, "") // No userID

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBooksAPI_AddToBookshelf_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupBooksTestRouter(mockService, userID)

	mockService.On("SaveReadingProgress", mock.Anything, userID, bookID, "", float64(0)).
		Return(nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooksAPI_AddToBookshelf_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	router := setupBooksTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+primitive.NewObjectID().Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBooksAPI_RemoveFromBookshelf_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupBooksTestRouter(mockService, userID)

	mockService.On("DeleteReadingProgress", mock.Anything, userID, bookID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/books/"+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooksAPI_RemoveFromBookshelf_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	router := setupBooksTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/books/"+primitive.NewObjectID().Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBooksAPI_GetRecentReading_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	router := setupBooksTestRouter(mockService, userID)

	expectedProgresses := []*readerModels.ReadingProgress{}

	mockService.On("GetRecentReading", mock.Anything, userID, 10).
		Return(expectedProgresses, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/recent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooksAPI_GetRecentReading_CustomLimit(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	router := setupBooksTestRouter(mockService, userID)

	expectedProgresses := []*readerModels.ReadingProgress{}

	mockService.On("GetRecentReading", mock.Anything, userID, 50).
		Return(expectedProgresses, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/recent?limit=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooksAPI_GetUnfinishedBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	router := setupBooksTestRouter(mockService, userID)

	expectedProgresses := []*readerModels.ReadingProgress{}

	mockService.On("GetUnfinishedBooks", mock.Anything, userID).
		Return(expectedProgresses, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/unfinished", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBooksAPI_GetFinishedBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockReaderServiceForBooks)
	userID := primitive.NewObjectID().Hex()
	router := setupBooksTestRouter(mockService, userID)

	expectedProgresses := []*readerModels.ReadingProgress{}

	mockService.On("GetFinishedBooks", mock.Anything, userID).
		Return(expectedProgresses, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/finished", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}
