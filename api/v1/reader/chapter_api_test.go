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

	"Qingyu_backend/service/interfaces"
	readerservice "Qingyu_backend/service/reader"
)

// MockChapterService 模拟ChapterService
type MockChapterService struct {
	mock.Mock
}

func (m *MockChapterService) GetChapterContent(ctx context.Context, userID, bookID, chapterID string) (*readerservice.ChapterContentResponse, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.ChapterContentResponse), args.Error(1)
}

func (m *MockChapterService) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*readerservice.ChapterContentResponse, error) {
	args := m.Called(ctx, userID, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.ChapterContentResponse), args.Error(1)
}

func (m *MockChapterService) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*readerservice.ChapterInfo, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.ChapterInfo), args.Error(1)
}

func (m *MockChapterService) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*readerservice.ChapterInfo, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.ChapterInfo), args.Error(1)
}

func (m *MockChapterService) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*readerservice.ChapterListResponse, error) {
	args := m.Called(ctx, userID, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.ChapterListResponse), args.Error(1)
}

func (m *MockChapterService) GetChapterInfo(ctx context.Context, userID, chapterID string) (*readerservice.ChapterInfo, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.ChapterInfo), args.Error(1)
}

func setupChapterTestRouter(chapterService interfaces.ReaderChapterService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewChapterAPI(chapterService)

	v1 := r.Group("/api/v1/reader/books/:bookId/chapters")
	{
		v1.GET("/:chapterId", api.GetChapterContent)
		v1.GET("/:chapterId/next", api.GetNextChapter)
		v1.GET("/:chapterId/previous", api.GetPreviousChapter)
		v1.GET("", api.GetChapterList)
		v1.GET("/by-number/:chapterNum", api.GetChapterByNumber)
	}

	r.GET("/api/v1/reader/chapters/:chapterId/info", api.GetChapterInfo)

	return r
}

func TestChapterAPI_GetChapterContent_Success(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	expectedContent := &readerservice.ChapterContentResponse{
		ChapterID:   chapterID,
		BookID:      bookID,
		Title:       "第一章",
		Content:     "章节内容",
		ChapterNum:  1,
		WordCount:   4,
		HasNext:     true,
		HasPrevious: false,
		Progress:    0.0,
		ReadingTime: 60,
	}

	mockService.On("GetChapterContent", mock.Anything, userID, bookID, chapterID).
		Return(expectedContent, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetChapterContent_NotFound(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	mockService.On("GetChapterContent", mock.Anything, userID, bookID, chapterID).
		Return(nil, readerservice.ErrChapterNotFound)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChapterAPI_GetChapterContent_NotPublished(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	mockService.On("GetChapterContent", mock.Anything, userID, bookID, chapterID).
		Return(nil, readerservice.ErrChapterNotPublished)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChapterAPI_GetChapterByNumber_Success(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	chapterID := primitive.NewObjectID().Hex()
	expectedContent := &readerservice.ChapterContentResponse{
		ChapterID:   chapterID,
		BookID:      bookID,
		Title:       "第一章",
		Content:     "章节内容",
		ChapterNum:  1,
		WordCount:   4,
		HasNext:     false,
		HasPrevious: false,
		Progress:    0.0,
		ReadingTime: 60,
	}

	mockService.On("GetChapterByNumber", mock.Anything, userID, bookID, 1).
		Return(expectedContent, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/by-number/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetChapterByNumber_InvalidNumber(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	bookID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/by-number/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterAPI_GetNextChapter_Success(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	nextChapterID := primitive.NewObjectID().Hex()
	nextChapter := &readerservice.ChapterInfo{
		ChapterID:   nextChapterID,
		BookID:      bookID,
		Title:       "第二章",
		ChapterNum:  2,
		WordCount:   500,
		IsFree:      true,
		Price:       0,
		PublishTime: time.Now(),
		Progress:    0,
		IsRead:      false,
	}

	mockService.On("GetNextChapter", mock.Anything, userID, bookID, chapterID).
		Return(nextChapter, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID+"/next", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetNextChapter_LastChapter(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	mockService.On("GetNextChapter", mock.Anything, userID, bookID, chapterID).
		Return(nil, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID+"/next", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetPreviousChapter_Success(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	prevChapterID := primitive.NewObjectID().Hex()
	prevChapter := &readerservice.ChapterInfo{
		ChapterID:   prevChapterID,
		BookID:      bookID,
		Title:       "第一章",
		ChapterNum:  1,
		WordCount:   500,
		IsFree:      true,
		Price:       0,
		PublishTime: time.Now(),
		Progress:    0,
		IsRead:      true,
	}

	mockService.On("GetPreviousChapter", mock.Anything, userID, bookID, chapterID).
		Return(prevChapter, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID+"/previous", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetPreviousChapter_FirstChapter(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	mockService.On("GetPreviousChapter", mock.Anything, userID, bookID, chapterID).
		Return(nil, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters/"+chapterID+"/previous", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetChapterList_Success(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	expectedList := &readerservice.ChapterListResponse{
		Chapters:   []*readerservice.ChapterInfo{},
		Total:      int64(100),
		Page:       1,
		Size:       50,
		BookID:     bookID,
		BookTitle:  "测试书籍",
		Author:     "测试作者",
		TotalWords: 50000,
	}

	mockService.On("GetChapterList", mock.Anything, userID, bookID, 1, 50).
		Return(expectedList, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/chapters?page=1&size=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetChapterList_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	router := setupChapterTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books//chapters", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterAPI_GetChapterInfo_Success(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	bookID := primitive.NewObjectID().Hex()
	expectedInfo := &readerservice.ChapterInfo{
		ChapterID:   chapterID,
		BookID:      bookID,
		Title:       "第一章",
		ChapterNum:  1,
		WordCount:   5000,
		IsFree:      true,
		Price:       0,
		PublishTime: time.Now(),
		Progress:    0,
		IsRead:      false,
	}

	mockService.On("GetChapterInfo", mock.Anything, userID, chapterID).
		Return(expectedInfo, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterAPI_GetChapterInfo_NotFound(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupChapterTestRouter(mockService, userID)

	mockService.On("GetChapterInfo", mock.Anything, userID, chapterID).
		Return(nil, readerservice.ErrChapterNotFound)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChapterAPI_GetChapterInfo_MissingChapterID(t *testing.T) {
	// Given
	mockService := new(MockChapterService)
	router := setupChapterTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters//info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
