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
	readerservice "Qingyu_backend/service/reader"
)

// MockBookmarkService 模拟书签服务
type MockBookmarkService struct {
	mock.Mock
}

func (m *MockBookmarkService) CreateBookmark(ctx context.Context, bookmark *readerModels.Bookmark) error {
	args := m.Called(ctx, bookmark)
	return args.Error(0)
}

func (m *MockBookmarkService) GetBookmark(ctx context.Context, bookmarkID string) (*readerModels.Bookmark, error) {
	args := m.Called(ctx, bookmarkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.Bookmark), args.Error(1)
}

func (m *MockBookmarkService) GetUserBookmarks(ctx context.Context, userID string, filter *readerModels.BookmarkFilter, page, size int) (*readerservice.BookmarkListResponse, error) {
	args := m.Called(ctx, userID, filter, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.BookmarkListResponse), args.Error(1)
}

func (m *MockBookmarkService) GetBookBookmarks(ctx context.Context, userID, bookID string, page, size int) (*readerservice.BookmarkListResponse, error) {
	args := m.Called(ctx, userID, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.BookmarkListResponse), args.Error(1)
}

func (m *MockBookmarkService) UpdateBookmark(ctx context.Context, bookmarkID string, bookmark *readerModels.Bookmark) error {
	args := m.Called(ctx, bookmarkID, bookmark)
	return args.Error(0)
}

func (m *MockBookmarkService) DeleteBookmark(ctx context.Context, bookmarkID string) error {
	args := m.Called(ctx, bookmarkID)
	return args.Error(0)
}

func (m *MockBookmarkService) ExportBookmarks(ctx context.Context, userID, format string) ([]byte, string, error) {
	args := m.Called(ctx, userID, format)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]byte), args.Get(1).(string), args.Error(2)
}

func (m *MockBookmarkService) GetBookmarkStats(ctx context.Context, userID string) (*readerModels.BookmarkStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModels.BookmarkStats), args.Error(1)
}

func (m *MockBookmarkService) SearchBookmarks(ctx context.Context, userID, keyword string, page, size int) (*readerservice.BookmarkListResponse, error) {
	args := m.Called(ctx, userID, keyword, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerservice.BookmarkListResponse), args.Error(1)
}

func setupBookmarkTestRouter(bookmarkService interfaces.BookmarkService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewBookmarkAPI(bookmarkService)

	v1 := r.Group("/api/v1/reader")
	{
		v1.POST("/books/:bookId/bookmarks", api.CreateBookmark)
		v1.GET("/bookmarks", api.GetBookmarks)
		v1.GET("/books/:bookId/bookmarks", api.GetBookmarks)
		v1.GET("/bookmarks/:id", api.GetBookmark)
		v1.PUT("/bookmarks/:id", api.UpdateBookmark)
		v1.DELETE("/bookmarks/:id", api.DeleteBookmark)
		v1.GET("/bookmarks/export", api.ExportBookmarks)
		v1.GET("/bookmarks/stats", api.GetBookmarkStats)
		v1.GET("/bookmarks/search", api.SearchBookmarks)
	}

	return r
}

func TestBookmarkAPI_CreateBookmark_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"bookId":    bookID,
		"chapterId": chapterID,
		"position":  100,
		"note":      "重要内容",
		"color":     "yellow",
		"quote":     "这是一段引用",
		"isPublic":  false,
		"tags":      []string{"重点", "复习"},
	}

	mockService.On("CreateBookmark", mock.Anything, mock.Anything).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID+"/bookmarks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_CreateBookmark_AlreadyExists(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"bookId":    bookID,
		"chapterId": chapterID,
		"position":  100,
	}

	mockService.On("CreateBookmark", mock.Anything, mock.Anything).
		Return(readerservice.ErrBookmarkAlreadyExists)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID+"/bookmarks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestBookmarkAPI_CreateBookmark_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	bookID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, "") // No userID

	reqBody := map[string]interface{}{
		"bookId":    bookID,
		"chapterId": primitive.NewObjectID().Hex(),
		"position":  100,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID+"/bookmarks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBookmarkAPI_GetBookmarks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	expectedBookmarks := &readerservice.BookmarkListResponse{
		Bookmarks: []*readerModels.Bookmark{
			{
				ID:        primitive.NewObjectID(),
				UserID:    primitive.NewObjectID(),
				BookID:    primitive.NewObjectID(),
				ChapterID: primitive.NewObjectID(),
				Position:  100,
				Note:      "测试笔记",
				Color:     "yellow",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Total: int64(1),
		Page:  1,
		Size:  20,
	}

	mockService.On("GetUserBookmarks", mock.Anything, userID, mock.Anything, 1, 20).
		Return(expectedBookmarks, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_GetBookBookmarks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	expectedBookmarks := &readerservice.BookmarkListResponse{
		Bookmarks: []*readerModels.Bookmark{},
		Total:     int64(0),
		Page:      1,
		Size:      20,
	}

	mockService.On("GetBookBookmarks", mock.Anything, userID, bookID, 1, 20).
		Return(expectedBookmarks, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/books/"+bookID+"/bookmarks?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_GetBookmark_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookmarkID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	expectedBookmark := &readerModels.Bookmark{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		BookID:    primitive.NewObjectID(),
		ChapterID: primitive.NewObjectID(),
		Position:  100,
		Note:      "测试笔记",
		Color:     "yellow",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetBookmark", mock.Anything, bookmarkID).Return(expectedBookmark, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks/"+bookmarkID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_GetBookmark_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookmarkID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	mockService.On("GetBookmark", mock.Anything, bookmarkID).
		Return(nil, readerservice.ErrBookmarkNotFound)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks/"+bookmarkID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookmarkAPI_UpdateBookmark_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookmarkID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"note":     "更新后的笔记",
		"color":    "blue",
		"isPublic": true,
		"tags":     []string{"更新"},
	}

	mockService.On("UpdateBookmark", mock.Anything, bookmarkID, mock.Anything).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/bookmarks/"+bookmarkID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_UpdateBookmark_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookmarkID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"note": "更新后的笔记",
	}

	mockService.On("UpdateBookmark", mock.Anything, bookmarkID, mock.Anything).
		Return(readerservice.ErrBookmarkNotFound)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/bookmarks/"+bookmarkID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookmarkAPI_DeleteBookmark_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookmarkID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	mockService.On("DeleteBookmark", mock.Anything, bookmarkID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/bookmarks/"+bookmarkID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_DeleteBookmark_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	bookmarkID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	mockService.On("DeleteBookmark", mock.Anything, bookmarkID).
		Return(readerservice.ErrBookmarkNotFound)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/bookmarks/"+bookmarkID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookmarkAPI_ExportBookmarks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	expectedData := []byte(`[{"note": "测试"}]`)
	mockService.On("ExportBookmarks", mock.Anything, userID, "json").
		Return(expectedData, "application/json", nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks/export?format=json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "attachment; filename=bookmarks_json", w.Header().Get("Content-Disposition"))
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_GetBookmarkStats_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	expectedStats := &readerModels.BookmarkStats{
		TotalCount:     100,
		PublicCount:    50,
		PrivateCount:   50,
		ByColor:        map[string]int64{"yellow": 50, "blue": 30, "red": 20},
		ByBook:         map[string]int64{},
		ThisMonthCount: 10,
		ThisWeekCount:  5,
	}

	mockService.On("GetBookmarkStats", mock.Anything, userID).Return(expectedStats, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_SearchBookmarks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	expectedResult := &readerservice.BookmarkListResponse{
		Bookmarks: []*readerModels.Bookmark{},
		Total:     int64(0),
		Page:      1,
		Size:      20,
	}

	mockService.On("SearchBookmarks", mock.Anything, userID, "测试", 1, 20).
		Return(expectedResult, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks/search?keyword=测试&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkAPI_SearchBookmarks_EmptyKeyword(t *testing.T) {
	// Given
	mockService := new(MockBookmarkService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookmarkTestRouter(mockService, userID)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/bookmarks/search?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
