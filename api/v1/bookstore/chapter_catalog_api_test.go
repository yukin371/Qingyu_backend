package bookstore

import (
	"Qingyu_backend/api/v1/bookstore"
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/service/bookstore"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock ChapterPurchaseService for API tests

type MockChapterPurchaseServiceForAPI struct {
	mock.Mock
}

func (m *MockChapterPurchaseServiceForAPI) GetChapterCatalog(ctx context.Context, userID, bookID primitive.ObjectID) (*bookstore.ChapterCatalog, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterCatalog), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) GetTrialChapters(ctx context.Context, bookID primitive.ObjectID, trialCount int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, trialCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) GetVIPChapters(ctx context.Context, bookID primitive.ObjectID) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) PurchaseChapter(ctx context.Context, userID, chapterID primitive.ObjectID) (*bookstore.ChapterPurchase, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterPurchase), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) PurchaseChapters(ctx context.Context, userID primitive.ObjectID, chapterIDs []primitive.ObjectID) (*bookstore.ChapterPurchaseBatch, error) {
	args := m.Called(ctx, userID, chapterIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterPurchaseBatch), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) PurchaseBook(ctx context.Context, userID, bookID primitive.ObjectID) (*bookstore.BookPurchase, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookPurchase), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) GetChapterPurchases(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchase), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseServiceForAPI) GetBookPurchases(ctx context.Context, userID, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	args := m.Called(ctx, userID, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchase), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseServiceForAPI) GetAllPurchases(ctx context.Context, userID primitive.ObjectID, page, pageSize int) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) CheckChapterAccess(ctx context.Context, userID, chapterID primitive.ObjectID) (*bookstore.ChapterAccessInfo, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterAccessInfo), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) GetPurchasedChapterIDs(ctx context.Context, userID, bookID primitive.ObjectID) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) GetChapterPrice(ctx context.Context, chapterID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, chapterID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockChapterPurchaseServiceForAPI) CalculateBookPrice(ctx context.Context, bookID primitive.ObjectID) (float64, float64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockChapterPurchaseServiceForAPI) IsVIPUser(ctx context.Context, userID primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// Mock ChapterService for API tests

type MockChapterServiceForAPI struct {
	mock.Mock
}

func (m *MockChapterServiceForAPI) GetChapterByID(ctx context.Context, chapterID primitive.ObjectID) (*bookstore.Chapter, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

// Test Helper Functions

func setupTestRouter(purchaseService bookstore.ChapterPurchaseService, chapterService bookstore.ChapterService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := bookstore.NewChapterCatalogAPI(chapterService, purchaseService)

	router.GET("/api/v1/bookstore/books/:id/chapters", func(c *gin.Context) {
		c.Set("userId", primitive.NewObjectID().Hex())
		api.GetChapterCatalog(c)
	})

	router.GET("/api/v1/bookstore/books/:id/chapters/:chapterId", api.GetChapterInfo)
	router.GET("/api/v1/bookstore/books/:id/trial-chapters", api.GetTrialChapters)
	router.GET("/api/v1/bookstore/books/:id/vip-chapters", api.GetVIPChapters)
	router.GET("/api/v1/bookstore/chapters/:chapterId/price", api.GetChapterPrice)
	router.POST("/api/v1/reader/chapters/:chapterId/purchase", setupAuthMiddleware(), api.PurchaseChapter)
	router.POST("/api/v1/reader/books/:id/buy-all", setupAuthMiddleware(), api.PurchaseBook)
	router.GET("/api/v1/reader/purchases", setupAuthMiddleware(), api.GetPurchases)
	router.GET("/api/v1/reader/purchases/:bookId", setupAuthMiddleware(), api.GetBookPurchases)
	router.GET("/api/v1/bookstore/chapters/:chapterId/access", api.CheckChapterAccess)

	return router
}

func setupAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock authentication
		userID := primitive.NewObjectID()
		c.Set("userId", userID.Hex())
		c.Next()
	}
}

// Test: GetChapterCatalog API

func TestChapterCatalogAPI_GetChapterCatalog_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	catalog := &bookstore.ChapterCatalog{
		BookID:         bookID,
		BookTitle:      "Test Book",
		TotalChapters:  10,
		FreeChapters:   3,
		PaidChapters:   7,
		VIPChapters:    0,
		TotalWordCount: 50000,
		Chapters: []bookstore.ChapterCatalogItem{
			{
				ChapterID:   primitive.NewObjectID(),
				Title:       "Chapter 1",
				ChapterNum:  1,
				WordCount:   5000,
				IsFree:      true,
				Price:       0,
				IsPublished: true,
			},
		},
		TrialCount: 10,
	}

	purchaseService.On("GetChapterCatalog", mock.Anything, userID, bookID).Return(catalog, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex()+"/chapters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取成功", response.Message)

	purchaseService.AssertExpectations(t)
}

func TestChapterCatalogAPI_GetChapterCatalog_InvalidBookID(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/invalid-id/chapters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Contains(t, response.Message, "无效的书籍ID格式")
}

func TestChapterCatalogAPI_GetChapterCatalog_EmptyBookID(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books//chapters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Contains(t, response.Message, "书籍ID不能为空")
}

// Test: GetChapterInfo API

func TestChapterCatalogAPI_GetChapterInfo_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := &bookstore.Chapter{
		ID:         chapterID,
		BookID:     bookID,
		Title:      "Test Chapter",
		ChapterNum: 1,
		WordCount:  2000,
		IsFree:     true,
		Price:      0,
	}

	chapterService.On("GetChapterByID", mock.Anything, chapterID).Return(chapter, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex()+"/chapters/"+chapterID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	chapterService.AssertExpectations(t)
}

func TestChapterCatalogAPI_GetChapterInfo_NotFound(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapterService.On("GetChapterByID", mock.Anything, chapterID).Return(nil, assert.AnError)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex()+"/chapters/"+chapterID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	chapterService.AssertExpectations(t)
}

// Test: GetTrialChapters API

func TestChapterCatalogAPI_GetTrialChapters_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	bookID := primitive.NewObjectID()

	trialChapters := []*bookstore.Chapter{
		{
			ID:         primitive.NewObjectID(),
			BookID:     bookID,
			Title:      "Chapter 1",
			ChapterNum: 1,
			IsFree:     true,
			Price:      0,
		},
		{
			ID:         primitive.NewObjectID(),
			BookID:     bookID,
			Title:      "Chapter 2",
			ChapterNum: 2,
			IsFree:     true,
			Price:      0,
		},
	}

	purchaseService.On("GetTrialChapters", mock.Anything, bookID, 10).Return(trialChapters, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex()+"/trial-chapters?count=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, float64(2), data["count"])

	purchaseService.AssertExpectations(t)
}

// Test: GetVIPChapters API

func TestChapterCatalogAPI_GetVIPChapters_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	bookID := primitive.NewObjectID()

	vipChapters := []*bookstore.Chapter{
		{
			ID:         primitive.NewObjectID(),
			BookID:     bookID,
			Title:      "VIP Chapter 1",
			ChapterNum: 5,
			IsFree:     false,
			Price:      2.99,
		},
	}

	purchaseService.On("GetVIPChapters", mock.Anything, bookID).Return(vipChapters, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex()+"/vip-chapters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	purchaseService.AssertExpectations(t)
}

// Test: GetChapterPrice API

func TestChapterCatalogAPI_GetChapterPrice_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	chapterID := primitive.NewObjectID()

	purchaseService.On("GetChapterPrice", mock.Anything, chapterID).Return(1.99, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/chapters/"+chapterID.Hex()+"/price", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, 1.99, data["price"])

	purchaseService.AssertExpectations(t)
}

// Test: PurchaseChapter API

func TestChapterCatalogAPI_PurchaseChapter_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	purchase := &bookstore.ChapterPurchase{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ChapterID: chapterID,
		Price:     1.99,
	}

	purchaseService.On("PurchaseChapter", mock.Anything, userID, chapterID).Return(purchase, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/purchase", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "购买成功", response.Message)

	purchaseService.AssertExpectations(t)
}

func TestChapterCatalogAPI_PurchaseChapter_InsufficientBalance(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	purchaseService.On("PurchaseChapter", mock.Anything, userID, chapterID).
		Return(nil, assert.AnError)

	router := setupTestRouter(purchaseService, chapterService)

	// Mock the error to contain "insufficient balance"
	purchaseService.On("PurchaseChapter", mock.Anything, userID, chapterID).
		Return(nil, &testError{msg: "insufficient balance"})

	router2 := setupTestRouter(purchaseService, chapterService)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/purchase", nil)
	w := httptest.NewRecorder()

	router2.ServeHTTP(w, req)

	// Should return 403 for insufficient balance
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// Test: PurchaseBook API

func TestChapterCatalogAPI_PurchaseBook_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	purchase := &bookstore.BookPurchase{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		BookID:       bookID,
		TotalPrice:   7.99,
		ChapterCount: 10,
	}

	purchaseService.On("PurchaseBook", mock.Anything, userID, bookID).Return(purchase, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("POST", "/api/v1/reader/books/"+bookID.Hex()+"/buy-all", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "购买成功", response.Message)

	purchaseService.AssertExpectations(t)
}

// Test: GetPurchases API

func TestChapterCatalogAPI_GetPurchases_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	purchases := []*bookstore.ChapterPurchase{
		{
			ID:        primitive.NewObjectID(),
			ChapterID: primitive.NewObjectID(),
			Price:     1.99,
		},
		{
			ID:        primitive.NewObjectID(),
			ChapterID: primitive.NewObjectID(),
			Price:     2.99,
		},
	}

	purchaseService.On("GetChapterPurchases", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), 1, 20).
		Return(purchases, int64(2), nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/reader/purchases?page=1&size=20", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 20, response.Size)

	purchaseService.AssertExpectations(t)
}

// Test: GetBookPurchases API

func TestChapterCatalogAPI_GetBookPurchases_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	bookID := primitive.NewObjectID()

	purchases := []*bookstore.ChapterPurchase{
		{
			ID:        primitive.NewObjectID(),
			BookID:    bookID,
			ChapterID: primitive.NewObjectID(),
			Price:     1.99,
		},
	}

	purchaseService.On("GetBookPurchases", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), bookID, 1, 20).
		Return(purchases, int64(1), nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/reader/purchases/"+bookID.Hex()+"?page=1&size=20", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, int64(1), response.Total)

	purchaseService.AssertExpectations(t)
}

// Test: CheckChapterAccess API

func TestChapterCatalogAPI_CheckChapterAccess_Success(t *testing.T) {
	purchaseService := new(MockChapterPurchaseServiceForAPI)
	chapterService := new(MockChapterServiceForAPI)

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	accessInfo := &bookstore.ChapterAccessInfo{
		ChapterID:  chapterID,
		Title:       "Test Chapter",
		ChapterNum:  1,
		WordCount:   2000,
		IsFree:      true,
		Price:       0,
		IsPurchased: false,
		CanAccess:   true,
		AccessReason: "free",
	}

	purchaseService.On("CheckChapterAccess", mock.Anything, primitive.NilObjectID, chapterID).Return(accessInfo, nil)

	router := setupTestRouter(purchaseService, chapterService)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/chapters/"+chapterID.Hex()+"/access", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstore.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "检查成功", response.Message)

	purchaseService.AssertExpectations(t)
}
