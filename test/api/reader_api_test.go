//go:build apitest_draft

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	readerAPI "Qingyu_backend/api/v1/reader"
	readerModel "Qingyu_backend/models/reading/reader"
	"Qingyu_backend/service/reading"
)

// ===========================
// 测试工具函数
// ===========================

// setupTestRouter 设置测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// mockAuth 模拟认证中间件
func mockAuth(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", userID)
		c.Next()
	}
}

// makeRequest 执行HTTP请求并返回响应
func makeRequest(router *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// parseResponse 解析响应
func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

// ===========================
// Mock ReaderService
// ===========================

type MockReaderService struct {
	mock.Mock
}

func (m *MockReaderService) GetChapterByID(ctx context.Context, id string) (*readerModel.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.Chapter), args.Error(1)
}

func (m *MockReaderService) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockReaderService) GetChaptersByBookID(ctx context.Context, bookID string, page, size int) ([]*readerModel.Chapter, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*readerModel.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderService) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*readerModel.Chapter, error) {
	args := m.Called(ctx, bookID, currentChapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.Chapter), args.Error(1)
}

func (m *MockReaderService) GetPreviousChapter(ctx context.Context, bookID string, currentChapterNum int) (*readerModel.Chapter, error) {
	args := m.Called(ctx, bookID, currentChapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.Chapter), args.Error(1)
}

func (m *MockReaderService) GetFirstChapter(ctx context.Context, bookID string) (*readerModel.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.Chapter), args.Error(1)
}

func (m *MockReaderService) GetLastChapter(ctx context.Context, bookID string) (*readerModel.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.Chapter), args.Error(1)
}

func (m *MockReaderService) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

func (m *MockReaderService) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.ReadingProgress), args.Error(1)
}

func (m *MockReaderService) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

func (m *MockReaderService) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModel.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.ReadingProgress), args.Error(1)
}

func (m *MockReaderService) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*readerModel.ReadingProgress, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*readerModel.ReadingProgress), args.Get(1).(int64), args.Error(2)
}

func (m *MockReaderService) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReaderService) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.ReadingProgress), args.Error(1)
}

func (m *MockReaderService) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.ReadingProgress), args.Error(1)
}

func (m *MockReaderService) CreateAnnotation(ctx context.Context, annotation *readerModel.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockReaderService) UpdateAnnotation(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockReaderService) DeleteAnnotation(ctx context.Context, userID, annotationID string) error {
	args := m.Called(ctx, userID, annotationID)
	return args.Error(0)
}

func (m *MockReaderService) GetAnnotationsByChapter(ctx context.Context, userID, chapterID string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetNotesByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetBookmarksByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetHighlightsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetPublicAnnotations(ctx context.Context, chapterID string) ([]*readerModel.Annotation, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*readerModel.Annotation), args.Error(1)
}

func (m *MockReaderService) GetReadingSettings(ctx context.Context, userID string) (*readerModel.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*readerModel.ReadingSettings), args.Error(1)
}

func (m *MockReaderService) SaveReadingSettings(ctx context.Context, settings *readerModel.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockReaderService) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockReaderService) ServiceInfo() (string, string) {
	return "ReaderService", "1.0.0"
}

func (m *MockReaderService) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ===========================
// 测试用例 - 章节API
// ===========================

// TestChaptersAPI_GetChapterByID 测试获取章节信息
func TestChaptersAPI_GetChapterByID(t *testing.T) {
	// 设置
	mockService := new(MockReaderService)
	api := readerAPI.NewChaptersAPI(&reading.ReaderService{})
	router := setupTestRouter()
	router.GET("/api/v1/reader/chapters/:id", api.GetChapterByID)

	// 测试用例
	t.Run("成功获取章节信息", func(t *testing.T) {
		testChapter := &readerModel.Chapter{
			ID:         "chapter123",
			BookID:     "book456",
			Title:      "第一章",
			ChapterNum: 1,
			WordCount:  3000,
			IsVIP:      false,
		}

		mockService.On("GetChapterByID", mock.Anything, "chapter123").Return(testChapter, nil)

		w := makeRequest(router, "GET", "/api/v1/reader/chapters/chapter123", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])
		assert.Equal(t, "获取成功", response["message"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "chapter123", data["id"])
		assert.Equal(t, "第一章", data["title"])

		mockService.AssertExpectations(t)
	})

	t.Run("章节不存在", func(t *testing.T) {
		mockService.On("GetChapterByID", mock.Anything, "nonexistent").Return(nil, errors.New("章节不存在"))

		w := makeRequest(router, "GET", "/api/v1/reader/chapters/nonexistent", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(404), response["code"])

		mockService.AssertExpectations(t)
	})
}

// TestChaptersAPI_GetChapterContent 测试获取章节内容
func TestChaptersAPI_GetChapterContent(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewChaptersAPI(mockService)
	router := setupTestRouter()

	// 使用认证中间件
	router.GET("/api/v1/reader/chapters/:id/content", mockAuth("user123"), api.GetChapterContent)

	t.Run("成功获取章节内容", func(t *testing.T) {
		mockService.On("GetChapterContent", mock.Anything, "user123", "chapter123").Return("这是章节内容...", nil)

		w := makeRequest(router, "GET", "/api/v1/reader/chapters/chapter123/content", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "这是章节内容...", data["content"])

		mockService.AssertExpectations(t)
	})

	t.Run("VIP章节无权限", func(t *testing.T) {
		mockService.On("GetChapterContent", mock.Anything, "user123", "vip_chapter").Return("", errors.New("需要VIP权限"))

		w := makeRequest(router, "GET", "/api/v1/reader/chapters/vip_chapter/content", nil)

		assert.Equal(t, http.StatusForbidden, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(403), response["code"])

		mockService.AssertExpectations(t)
	})
}

// TestChaptersAPI_GetBookChapters 测试获取书籍章节列表
func TestChaptersAPI_GetBookChapters(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewChaptersAPI(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/reader/chapters", api.GetBookChapters)

	t.Run("成功获取章节列表", func(t *testing.T) {
		chapters := []*readerModel.Chapter{
			{ID: "c1", Title: "第一章", ChapterNum: 1},
			{ID: "c2", Title: "第二章", ChapterNum: 2},
		}

		mockService.On("GetChaptersByBookID", mock.Anything, "book123", 1, 20).Return(chapters, int64(2), nil)

		w := makeRequest(router, "GET", "/api/v1/reader/chapters?bookId=book123&page=1&size=20", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		data := response["data"].(map[string]interface{})
		chapterList := data["list"].([]interface{})
		assert.Equal(t, 2, len(chapterList))

		mockService.AssertExpectations(t)
	})

	t.Run("缺少bookId参数", func(t *testing.T) {
		w := makeRequest(router, "GET", "/api/v1/reader/chapters?page=1", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(400), response["code"])
	})
}

// ===========================
// 测试用例 - 进度API
// ===========================

// TestProgressAPI_SaveReadingProgress 测试保存阅读进度
func TestProgressAPI_SaveReadingProgress(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewProgressAPI(mockService)
	router := setupTestRouter()
	router.POST("/api/v1/reader/progress", mockAuth("user123"), api.SaveReadingProgress)

	t.Run("成功保存进度", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"bookId":    "book123",
			"chapterId": "chapter456",
			"progress":  0.75,
		}

		mockService.On("SaveReadingProgress", mock.Anything, "user123", "book123", "chapter456", 0.75).Return(nil)

		w := makeRequest(router, "POST", "/api/v1/reader/progress", requestBody)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])
		assert.Equal(t, "保存成功", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("无效的进度值", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"bookId":    "book123",
			"chapterId": "chapter456",
			"progress":  1.5, // 无效：超过1.0
		}

		w := makeRequest(router, "POST", "/api/v1/reader/progress", requestBody)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(400), response["code"])
	})
}

// TestProgressAPI_GetReadingProgress 测试获取阅读进度
func TestProgressAPI_GetReadingProgress(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewProgressAPI(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/reader/progress", mockAuth("user123"), api.GetReadingProgress)

	t.Run("成功获取进度", func(t *testing.T) {
		progress := &readerModel.ReadingProgress{
			UserID:      "user123",
			BookID:      "book123",
			ChapterID:   "chapter456",
			Progress:    0.75,
			ReadingTime: 3600,
		}

		mockService.On("GetReadingProgress", mock.Anything, "user123", "book123").Return(progress, nil)

		w := makeRequest(router, "GET", "/api/v1/reader/progress?bookId=book123", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "book123", data["bookId"])
		assert.Equal(t, 0.75, data["progress"])

		mockService.AssertExpectations(t)
	})

	t.Run("进度不存在", func(t *testing.T) {
		mockService.On("GetReadingProgress", mock.Anything, "user123", "new_book").Return(nil, nil)

		w := makeRequest(router, "GET", "/api/v1/reader/progress?bookId=new_book", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Nil(t, response["data"])

		mockService.AssertExpectations(t)
	})
}

// TestProgressAPI_GetRecentReading 测试获取最近阅读
func TestProgressAPI_GetRecentReading(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewProgressAPI(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/reader/progress/recent", mockAuth("user123"), api.GetRecentReading)

	t.Run("成功获取最近阅读", func(t *testing.T) {
		recentBooks := []*readerModel.ReadingProgress{
			{BookID: "book1", ChapterID: "chapter1", Progress: 0.5},
			{BookID: "book2", ChapterID: "chapter2", Progress: 0.8},
		}

		mockService.On("GetRecentReading", mock.Anything, "user123", 10).Return(recentBooks, nil)

		w := makeRequest(router, "GET", "/api/v1/reader/progress/recent?limit=10", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		data := response["data"].([]interface{})
		assert.Equal(t, 2, len(data))

		mockService.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 标注API
// ===========================

// TestAnnotationsAPI_CreateAnnotation 测试创建标注
func TestAnnotationsAPI_CreateAnnotation(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewAnnotationsAPI(mockService)
	router := setupTestRouter()
	router.POST("/api/v1/reader/annotations", mockAuth("user123"), api.CreateAnnotation)

	t.Run("成功创建书签", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"bookId":    "book123",
			"chapterId": "chapter456",
			"type":      "bookmark",
			"isPublic":  false,
		}

		mockService.On("CreateAnnotation", mock.Anything, mock.MatchedBy(func(a *readerModel.Annotation) bool {
			return a.UserID == "user123" && a.Type == "bookmark"
		})).Return(nil)

		w := makeRequest(router, "POST", "/api/v1/reader/annotations", requestBody)

		assert.Equal(t, http.StatusCreated, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(201), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("成功创建笔记", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"bookId":    "book123",
			"chapterId": "chapter456",
			"type":      "note",
			"content":   "这是一条笔记",
			"isPublic":  false,
		}

		mockService.On("CreateAnnotation", mock.Anything, mock.MatchedBy(func(a *readerModel.Annotation) bool {
			return a.Type == "note" && a.Note == "这是一条笔记"
		})).Return(nil)

		w := makeRequest(router, "POST", "/api/v1/reader/annotations", requestBody)

		assert.Equal(t, http.StatusCreated, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("无效的标注类型", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"bookId":    "book123",
			"chapterId": "chapter456",
			"type":      "invalid_type", // 无效类型
		}

		w := makeRequest(router, "POST", "/api/v1/reader/annotations", requestBody)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestAnnotationsAPI_DeleteAnnotation 测试删除标注
func TestAnnotationsAPI_DeleteAnnotation(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewAnnotationsAPI(mockService)
	router := setupTestRouter()
	router.DELETE("/api/v1/reader/annotations/:id", mockAuth("user123"), api.DeleteAnnotation)

	t.Run("成功删除标注", func(t *testing.T) {
		mockService.On("DeleteAnnotation", mock.Anything, "user123", "annotation123").Return(nil)

		w := makeRequest(router, "DELETE", "/api/v1/reader/annotations/annotation123", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("删除不存在的标注", func(t *testing.T) {
		mockService.On("DeleteAnnotation", mock.Anything, "user123", "nonexistent").Return(errors.New("标注不存在"))

		w := makeRequest(router, "DELETE", "/api/v1/reader/annotations/nonexistent", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})
}

// TestAnnotationsAPI_GetBookmarks 测试获取书签
func TestAnnotationsAPI_GetBookmarks(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewAnnotationsAPI(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/reader/bookmarks", mockAuth("user123"), api.GetBookmarks)

	t.Run("成功获取书签列表", func(t *testing.T) {
		bookmarks := []*readerModel.Annotation{
			{ID: "b1", Type: "bookmark", BookID: "book123", ChapterID: "c1"},
			{ID: "b2", Type: "bookmark", BookID: "book123", ChapterID: "c2"},
		}

		mockService.On("GetBookmarksByBook", mock.Anything, "user123", "book123").Return(bookmarks, nil)

		w := makeRequest(router, "GET", "/api/v1/reader/bookmarks?bookId=book123", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		data := response["data"].([]interface{})
		assert.Equal(t, 2, len(data))

		mockService.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 设置API
// ===========================

// TestSettingAPI_GetReadingSettings 测试获取阅读设置
func TestSettingAPI_GetReadingSettings(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewSettingAPI(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/reader/settings", mockAuth("user123"), api.GetReadingSettings)

	t.Run("成功获取设置", func(t *testing.T) {
		settings := &readerModel.ReadingSettings{
			UserID:     "user123",
			FontSize:   16,
			FontFamily: "serif",
			LineHeight: 1.8,
			Theme:      "light",
		}

		mockService.On("GetReadingSettings", mock.Anything, "user123").Return(settings, nil)

		w := makeRequest(router, "GET", "/api/v1/reader/settings", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(16), data["fontSize"])
		assert.Equal(t, "serif", data["fontFamily"])

		mockService.AssertExpectations(t)
	})
}

// TestSettingAPI_SaveReadingSettings 测试保存阅读设置
func TestSettingAPI_SaveReadingSettings(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewSettingAPI(mockService)
	router := setupTestRouter()
	router.POST("/api/v1/reader/settings", mockAuth("user123"), api.SaveReadingSettings)

	t.Run("成功保存设置", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"fontSize":   18,
			"fontFamily": "sans-serif",
			"lineHeight": 2.0,
			"theme":      "dark",
		}

		mockService.On("SaveReadingSettings", mock.Anything, mock.MatchedBy(func(s *readerModel.ReadingSettings) bool {
			return s.UserID == "user123" && s.FontSize == 18
		})).Return(nil)

		w := makeRequest(router, "POST", "/api/v1/reader/settings", requestBody)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("无效的字体大小", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"fontSize": 5, // 太小
		}

		w := makeRequest(router, "POST", "/api/v1/reader/settings", requestBody)

		// 应该返回400或由验证中间件处理
		assert.True(t, w.Code >= 400)
	})
}

// TestSettingAPI_UpdateReadingSettings 测试更新阅读设置
func TestSettingAPI_UpdateReadingSettings(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewSettingAPI(mockService)
	router := setupTestRouter()
	router.PATCH("/api/v1/reader/settings", mockAuth("user123"), api.UpdateReadingSettings)

	t.Run("成功更新单个设置项", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"fontSize": 20,
		}

		mockService.On("UpdateReadingSettings", mock.Anything, "user123", map[string]interface{}{
			"fontSize": 20,
		}).Return(nil)

		w := makeRequest(router, "PATCH", "/api/v1/reader/settings", requestBody)

		assert.Equal(t, http.StatusOK, w.Code)
		response := parseResponse(w)
		assert.Equal(t, float64(200), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("成功更新多个设置项", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"fontSize": 18,
			"theme":    "dark",
		}

		mockService.On("UpdateReadingSettings", mock.Anything, "user123", mock.Anything).Return(nil)

		w := makeRequest(router, "PATCH", "/api/v1/reader/settings", requestBody)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 认证和权限
// ===========================

// TestAuth_MissingToken 测试缺少认证Token
func TestAuth_MissingToken(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewChaptersAPI(mockService)
	router := setupTestRouter()

	// 不添加mockAuth中间件
	router.GET("/api/v1/reader/chapters/:id/content", api.GetChapterContent)

	w := makeRequest(router, "GET", "/api/v1/reader/chapters/chapter123/content", nil)

	// 应该返回401未授权
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ===========================
// 测试用例 - 并发安全
// ===========================

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	mockService := new(MockReaderService)
	api := readerAPI.NewChaptersAPI(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/reader/chapters/:id", api.GetChapterByID)

	testChapter := &readerModel.Chapter{
		ID:    "chapter123",
		Title: "测试章节",
	}

	// Mock可以处理多次调用
	mockService.On("GetChapterByID", mock.Anything, "chapter123").Return(testChapter, nil).Times(10)

	// 并发发送10个请求
	results := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go func() {
			w := makeRequest(router, "GET", "/api/v1/reader/chapters/chapter123", nil)
			results <- w.Code
		}()
	}

	// 验证所有请求都成功
	for i := 0; i < 10; i++ {
		code := <-results
		assert.Equal(t, http.StatusOK, code)
	}
}
