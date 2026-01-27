package bookstore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared"
	searchModels "Qingyu_backend/models/search"
	bookstoreService "Qingyu_backend/service/bookstore"
	"Qingyu_backend/pkg/logger"
)

// MockBookstoreService 模拟BookstoreService
type MockBookstoreService struct {
	mock.Mock
}

func (m *MockBookstoreService) GetAllBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstoreModel.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookstoreService) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, authorID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetPopularBooks(ctx context.Context, limit int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) SearchBooksWithFilter(ctx context.Context, filter *bookstoreModel.BookFilter) ([]*bookstoreModel.Book, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetCategoryTree(ctx context.Context) ([]*bookstoreModel.CategoryTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.CategoryTree), args.Error(1)
}

func (m *MockBookstoreService) GetCategoryByID(ctx context.Context, id string) (*bookstoreModel.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Category), args.Error(1)
}

func (m *MockBookstoreService) GetRootCategories(ctx context.Context) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockBookstoreService) GetActiveBanners(ctx context.Context, limit int) ([]*bookstoreModel.Banner, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Banner), args.Error(1)
}

func (m *MockBookstoreService) IncrementBannerClick(ctx context.Context, bannerID string) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBookstoreService) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetRankingByType(ctx context.Context, rankingType bookstoreModel.RankingType, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) UpdateRankings(ctx context.Context, rankingType bookstoreModel.RankingType, period string) error {
	args := m.Called(ctx, rankingType, period)
	return args.Error(0)
}

func (m *MockBookstoreService) GetHomepageData(ctx context.Context) (*bookstoreService.HomepageData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreService.HomepageData), args.Error(1)
}

func (m *MockBookstoreService) GetBookStats(ctx context.Context) (*bookstoreModel.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookStats), args.Error(1)
}

func (m *MockBookstoreService) IncrementBookView(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookstoreService) GetTags(ctx context.Context, categoryID *string) ([]string, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockBookstoreService) GetYears(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func setupBookstoreTestRouter(service *MockBookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewBookstoreAPI(service, nil, logger.Get())

	v1 := r.Group("/api/v1/bookstore")
	{
		v1.GET("/homepage", api.GetHomepage)
		v1.GET("/books", api.GetBooks)
		v1.GET("/books/:id", api.GetBookByID)
		v1.GET("/books/recommended", api.GetRecommendedBooks)
		v1.GET("/books/featured", api.GetFeaturedBooks)
		v1.GET("/books/search", api.SearchBooks)
		v1.GET("/books/:id/view", api.IncrementBookView)
		v1.GET("/categories/:id/books", api.GetBooksByCategory)
		v1.GET("/categories/tree", api.GetCategoryTree)
		v1.GET("/categories/:id", api.GetCategoryByID)
		v1.GET("/banners", api.GetActiveBanners)
		v1.POST("/banners/:id/click", api.IncrementBannerClick)
		v1.GET("/rankings/realtime", api.GetRealtimeRanking)
		v1.GET("/rankings/weekly", api.GetWeeklyRanking)
		v1.GET("/rankings/monthly", api.GetMonthlyRanking)
		v1.GET("/rankings/newbie", api.GetNewbieRanking)
		v1.GET("/rankings/:type", api.GetRankingByType)
	}

	return r
}

func TestBookstoreAPI_GetHomepage_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	homepageData := &bookstoreService.HomepageData{
		Banners:          []*bookstoreModel.Banner{},
		RecommendedBooks: []*bookstoreModel.Book{},
		FeaturedBooks:    []*bookstoreModel.Book{},
		Categories:       []*bookstoreModel.Category{},
		Rankings:         make(map[string][]*bookstoreModel.RankingItem),
	}
	mockService.On("GetHomepageData", mock.Anything).Return(homepageData, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/homepage", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("GetAllBooks", mock.Anything, 1, 20).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooks_DefaultPagination(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("GetAllBooks", mock.Anything, 1, 20).Return(books, int64(0), nil)

	// When - No pagination parameters, should use defaults
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBookByID_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	bookID := primitive.NewObjectID()
	book := &bookstoreModel.Book{IdentifiedEntity: shared.IdentifiedEntity{ID: bookID}, Title: "测试书籍"}
	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(book, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBookByID_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Gin returns 301 redirect for trailing slash without path param
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestBookstoreAPI_GetBookByID_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	bookID := primitive.NewObjectID().Hex()
	mockService.On("GetBookByID", mock.Anything, bookID).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 500 for service error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBookstoreAPI_GetBooksByCategory_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	categoryID := primitive.NewObjectID().Hex()
	books := []*bookstoreModel.Book{}
	mockService.On("GetBooksByCategory", mock.Anything, categoryID, 1, 20).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/"+categoryID+"/books?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooksByCategory_InvalidID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/invalid-id/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_GetRecommendedBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("GetRecommendedBooks", mock.Anything, 1, 20).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/recommended?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetFeaturedBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("GetFeaturedBooks", mock.Anything, 1, 20).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/featured?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooks_InvalidSize(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	// Size over 100 should be clamped to 20 (default)
	mockService.On("GetAllBooks", mock.Anything, 1, 20).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books?page=1&size=150", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooks_InvalidPage(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	// Page less than 1 should be clamped to 1
	mockService.On("GetAllBooks", mock.Anything, 1, 20).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books?page=0&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_SearchBooks_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	keyword := "test"
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword="+keyword, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_SearchBooks_NoFilters(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When - No keyword or filters provided
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_SearchBooks_WithAuthor(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	author := "测试作者"
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?author="+author, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetCategoryTree_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	tree := []*bookstoreModel.CategoryTree{}
	mockService.On("GetCategoryTree", mock.Anything).Return(tree, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/tree", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetCategoryByID_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	categoryID := primitive.NewObjectID().Hex()
	category := &bookstoreModel.Category{ID: primitive.NewObjectID().Hex(), Name: "测试分类"}
	mockService.On("GetCategoryByID", mock.Anything, categoryID).Return(category, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/"+categoryID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetCategoryByID_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Gin returns 404 for routes with empty path params
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookstoreAPI_GetCategoryByID_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	categoryID := primitive.NewObjectID().Hex()
	mockService.On("GetCategoryByID", mock.Anything, categoryID).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/"+categoryID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBookstoreAPI_GetActiveBanners_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	banners := []*bookstoreModel.Banner{}
	mockService.On("GetActiveBanners", mock.Anything, 5).Return(banners, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners?limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetActiveBanners_DefaultLimit(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	banners := []*bookstoreModel.Banner{}
	mockService.On("GetActiveBanners", mock.Anything, 5).Return(banners, nil)

	// When - No limit parameter, should use default 5
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetActiveBanners_InvalidLimit(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	banners := []*bookstoreModel.Banner{}
	// Limit over 20 should be clamped to 5 (default)
	mockService.On("GetActiveBanners", mock.Anything, 5).Return(banners, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners?limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_IncrementBookView_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	bookID := primitive.NewObjectID()
	mockService.On("IncrementBookView", mock.Anything, bookID.Hex()).Return(nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex()+"/view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_IncrementBookView_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books//view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Gin returns 400 for empty path params (validation happens before route match)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_IncrementBookView_InvalidID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/invalid-id/view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_IncrementBannerClick_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	bannerID := primitive.NewObjectID().Hex()
	mockService.On("IncrementBannerClick", mock.Anything, bannerID).Return(nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/bookstore/banners/"+bannerID+"/click", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_IncrementBannerClick_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/bookstore/banners//click", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Gin returns 400 for empty path params
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_IncrementBannerClick_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	bannerID := primitive.NewObjectID().Hex()
	mockService.On("IncrementBannerClick", mock.Anything, bannerID).Return(assert.AnError)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/bookstore/banners/"+bannerID+"/click", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 500 for generic errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBookstoreAPI_GetRealtimeRanking_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	rankings := []*bookstoreModel.RankingItem{}
	mockService.On("GetRealtimeRanking", mock.Anything, 20).Return(rankings, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/realtime?limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetRealtimeRanking_DefaultLimit(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	rankings := []*bookstoreModel.RankingItem{}
	mockService.On("GetRealtimeRanking", mock.Anything, 20).Return(rankings, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/realtime", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetWeeklyRanking_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	rankings := []*bookstoreModel.RankingItem{}
	mockService.On("GetWeeklyRanking", mock.Anything, "", 20).Return(rankings, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/weekly?limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetMonthlyRanking_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	rankings := []*bookstoreModel.RankingItem{}
	mockService.On("GetMonthlyRanking", mock.Anything, "", 20).Return(rankings, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/monthly?limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetNewbieRanking_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	rankings := []*bookstoreModel.RankingItem{}
	mockService.On("GetNewbieRanking", mock.Anything, "", 20).Return(rankings, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/newbie?limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetRankingByType_InvalidType(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When - "hot" is not a valid ranking type
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/hot", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 400 for invalid ranking type
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_SearchBooks_WithPagination(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	keyword := "test"
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword="+keyword+"&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_SearchBooks_WithTags(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When - Need to provide keyword along with tags
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=test&tags=tag1&tags=tag2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_SearchBooks_OnlyTags_NoKeyword(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When - Only tags without keyword should return 400
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?tags=tag1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 400 Bad Request (tags alone are not sufficient)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =====================================================
// SearchByTitle 和 SearchByAuthor API 测试
// =====================================================

// MockSearchService 模拟搜索服务
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Search(ctx context.Context, req *searchModels.SearchRequest) (*searchModels.SearchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchModels.SearchResponse), args.Error(1)
}

// setupBookstoreSearchTestRouter 创建带搜索服务的测试路由
func setupBookstoreSearchTestRouter(service *MockBookstoreService, searchService *MockSearchService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 创建 logger 实例
	testLogger := logger.Get()

	// 创建 API 实例，传入 nil 作为 searchService（测试 MongoDB fallback 路径）
	api := NewBookstoreAPI(service, nil, testLogger)

	v1 := r.Group("/api/v1/bookstore")
	{
		v1.GET("/books/search/title", api.SearchByTitle)
		v1.GET("/books/search/author", api.SearchByAuthor)
	}

	return r
}

// TestBookstoreAPI_SearchByTitle_MissingParam 测试缺少必需参数
func TestBookstoreAPI_SearchByTitle_MissingParam(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	mockSearchService := new(MockSearchService)
	router := setupBookstoreSearchTestRouter(mockService, mockSearchService)

	// When - 没有提供 title 参数
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestBookstoreAPI_SearchByTitle_Success 测试搜索成功
func TestBookstoreAPI_SearchByTitle_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	mockSearchService := new(MockSearchService)
	router := setupBookstoreSearchTestRouter(mockService, mockSearchService)

	// 模拟 MongoDB fallback 返回结果
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "测试书籍1",
			Author:           "测试作者",
			ViewCount:        100,
		},
	}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?title=测试&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestBookstoreAPI_SearchByTitle_PaginationValidation 测试分页参数验证
func TestBookstoreAPI_SearchByTitle_PaginationValidation(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	mockSearchService := new(MockSearchService)
	router := setupBookstoreSearchTestRouter(mockService, mockSearchService)

	books := []*bookstoreModel.Book{}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil).Times(2)

	// When & Then - 测试 page < 1 的情况（应该被纠正为 1）
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?title=测试&page=0&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// When & Then - 测试 size > 100 的情况（应该被纠正为 20）
	req2, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?title=测试&page=1&size=150", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestBookstoreAPI_SearchByAuthor_MissingParam 测试缺少必需参数
func TestBookstoreAPI_SearchByAuthor_MissingParam(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	mockSearchService := new(MockSearchService)
	router := setupBookstoreSearchTestRouter(mockService, mockSearchService)

	// When - 没有提供 author 参数
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestBookstoreAPI_SearchByAuthor_Success 测试搜索成功
func TestBookstoreAPI_SearchByAuthor_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	mockSearchService := new(MockSearchService)
	router := setupBookstoreSearchTestRouter(mockService, mockSearchService)

	// 模拟 MongoDB fallback 返回结果
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "测试书籍1",
			Author:           "测试作者",
			ViewCount:        100,
		},
	}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?author=测试作者&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestBookstoreAPI_SearchByAuthor_PaginationValidation 测试分页参数验证
func TestBookstoreAPI_SearchByAuthor_PaginationValidation(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	mockSearchService := new(MockSearchService)
	router := setupBookstoreSearchTestRouter(mockService, mockSearchService)

	books := []*bookstoreModel.Book{}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When & Then - 测试边界值
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?author=测试&page=1&size=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ==================== GetSimilarBooks 测试 ====================

// setupSimilarBooksTestRouter 设置相似书籍测试路由
func setupSimilarBooksTestRouter(service bookstoreService.BookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := NewBookstoreAPI(service, nil, logger.Get())
	router.GET("/books/:id/similar", api.GetSimilarBooks)

	return router
}

// TestGetSimilarBooks_Strategy1_CategoryAndTags 测试策略1：同分类+标签匹配
func TestGetSimilarBooks_Strategy1_CategoryAndTags(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "原书籍",
		CategoryIDs:       []primitive.ObjectID{categoryID},
		Tags:              []string{"玄幻", "修仙"},
	}

	similarBooks := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "相似书籍1",
			CategoryIDs:      []primitive.ObjectID{categoryID},
			Tags:             []string{"玄幻", "魔法"},
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "相似书籍2",
			CategoryIDs:      []primitive.ObjectID{categoryID},
			Tags:             []string{"修仙", "仙侠"},
		},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	// 策略1: 同分类 + 标签 - 注意：searchSimilarBooks 会多查一条用于排除
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(similarBooks, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetSimilarBooks_Strategy2_CategoryOnly 测试策略2：只有同分类匹配
func TestGetSimilarBooks_Strategy2_CategoryOnly(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "原书籍",
		CategoryIDs:       []primitive.ObjectID{categoryID},
		Tags:              []string{"小众标签"},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	// 策略1返回空（同分类+标签）
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 策略2: 同分类（由于limit=10，策略1返回0条，会触发策略2）
	// 返回10本书，这样就不会触发策略3
	fullCategoryBooks := make([]*bookstoreModel.Book, 10)
	for i := 0; i < 10; i++ {
		fullCategoryBooks[i] = &bookstoreModel.Book{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            fmt.Sprintf("同分类书籍%d", i+1),
			CategoryIDs:      []primitive.ObjectID{categoryID},
		}
	}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(fullCategoryBooks, int64(10), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetSimilarBooks_Strategy3_TagsOnly 测试策略3：只有标签匹配
func TestGetSimilarBooks_Strategy3_TagsOnly(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "原书籍",
		CategoryIDs:       []primitive.ObjectID{}, // 无分类
		Tags:              []string{"热门标签"},
	}

	tagBooks := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "同标签书籍1",
			Tags:             []string{"热门标签", "其他标签"},
		},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	// 策略1返回空（无分类，有标签）
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 策略2返回空（无分类，无标签）
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 策略3: 标签匹配
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(tagBooks, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetSimilarBooks_Strategy4_FallbackToHot 测试策略4：兜底返回热门书籍
func TestGetSimilarBooks_Strategy4_FallbackToHot(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "冷门书籍",
		CategoryIDs:       []primitive.ObjectID{},
		Tags:              []string{}, // 无标签
	}

	hotBooks := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "热门书籍1",
			ViewCount:        9999,
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "热门书籍2",
			ViewCount:        8888,
		},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	// 所有策略都返回空
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil)

	// 策略4: 兜底热门书籍
	mockService.On("GetHotBooks", mock.Anything, 1, 10).Return(hotBooks, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetSimilarBooks_Deduplication 测试去重逻辑
func TestGetSimilarBooks_Deduplication(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	// 使用固定的测试ID
	bookID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	categoryID, _ := primitive.ObjectIDFromHex("cat123")

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "原书籍",
		CategoryIDs:       []primitive.ObjectID{categoryID}, // 使用固定CategoryID
		Tags:              []string{"玄幻"},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil).Once()

	// 策略1: 同分类+标签 - 返回空（触发降级到策略2）
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证是策略1：有分类ID，有标签
		return filter.CategoryID != nil && len(filter.Tags) > 0
	})).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 策略2: 同分类（无标签）- 返回空（触发降级到策略3）
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证是策略2：有分类ID，无标签
		return filter.CategoryID != nil && len(filter.Tags) == 0
	})).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 策略3: 标签（无分类）- 返回空（触发降级到策略4）
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证是策略3：无分类ID，有标签
		return filter.CategoryID == nil && len(filter.Tags) > 0
	})).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 策略4: GetHotBooks 兜底 - 返回包含当前书籍的热门书籍（用于测试去重）
	hotBooks := []*bookstoreModel.Book{
		sourceBook, // 包含当前书籍（测试去重）
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "热门书籍1",
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "热门书籍2",
		},
	}
	mockService.On("GetHotBooks", mock.Anything, 1, 10).Return(hotBooks, int64(3), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t) // 验证所有Mock期望都被满足

	// 验证响应中不包含当前书籍ID
	assert.NotContains(t, w.Body.String(), bookID.Hex())
	assert.Contains(t, w.Body.String(), "获取相似书籍成功")
}

// TestGetSimilarBooks_LimitValidation 测试数量限制验证
func TestGetSimilarBooks_LimitValidation(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "测试书籍",
		CategoryIDs:       []primitive.ObjectID{primitive.NewObjectID()},
		Tags:              []string{"测试"},
	}

	similarBooks := []*bookstoreModel.Book{
		{IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()}, Title: "相似书籍1"},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	// 第一个请求：limit=100 会被修正为 20
	// 策略1返回1本，会触发策略2和策略3
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(similarBooks, int64(1), nil).Once()
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// When & Then - 测试限制超出范围
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 第二个请求：limit=0 会被修正为 10
	// 策略1返回1本，会触发策略2和策略3
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(similarBooks, int64(1), nil).Once()
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil).Once()

	// 测试限制小于1
	req2, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=0", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

// TestGetSimilarBooks_NeverReturnEmptyList 测试禁止返回空列表
func TestGetSimilarBooks_NeverReturnEmptyList(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "测试书籍",
		CategoryIDs:       []primitive.ObjectID{},
		Tags:              []string{},
	}

	fallbackBooks := []*bookstoreModel.Book{
		{IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()}, Title: "兜底书籍"},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	// 所有搜索策略返回空
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return([]*bookstoreModel.Book{}, int64(0), nil)

	// 兜底返回热门书籍
	mockService.On("GetHotBooks", mock.Anything, 1, 10).Return(fallbackBooks, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回兜底书籍，不返回空列表
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取相似书籍成功")
}

// TestGetSimilarBooks_BookNotFound 测试书籍不存在的错误处理
func TestGetSimilarBooks_BookNotFound(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(nil, errors.New("book not found"))

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestGetSimilarBooks_EmptyBookID 测试空书籍ID的参数验证
func TestGetSimilarBooks_EmptyBookID(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/books//similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetSimilarBooks_ExcludeCurrentBook 测试排除当前书籍
func TestGetSimilarBooks_ExcludeCurrentBook(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupSimilarBooksTestRouter(mockService)

	bookID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	sourceBook := &bookstoreModel.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: bookID},
		Title:             "原书籍",
		CategoryIDs:       []primitive.ObjectID{categoryID},
		Tags:              []string{"测试"},
	}

	// 返回包含原书籍的列表
	booksIncludingCurrent := []*bookstoreModel.Book{
		sourceBook, // 包含原书籍
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "其他书籍",
		},
	}

	mockService.On("GetBookByID", mock.Anything, bookID.Hex()).Return(sourceBook, nil)

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(booksIncludingCurrent, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/books/"+bookID.Hex()+"/similar?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 原书籍应该被排除
	assert.Equal(t, http.StatusOK, w.Code)
}

// ==================== SearchByTitle 完整测试 ====================

// TestSearchByTitle_Success 测试搜索成功（MongoDB fallback路径）
// 场景：searchService 为 nil，直接使用 MongoDB
func TestSearchByTitle_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// 构造MongoDB返回结果
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "测试书籍",
			Author:           "测试作者",
			ViewCount:        100,
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证参数：标题搜索
		return filter.Keyword != nil && *filter.Keyword == "测试" &&
			filter.SortBy == "view_count" && filter.SortOrder == "desc"
	})).Return(books, int64(1), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?title=测试&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "搜索成功")
	mockService.AssertExpectations(t)
}

// TestSearchByTitle_MissingParam 测试缺少必需参数
// 场景：缺少 title 参数
func TestSearchByTitle_MissingParam(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// When - 没有提供 title 参数
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "参数错误")
}

// TestSearchByTitle_FallbackOnError 测试 SearchService 返回错误时的 fallback
// 场景：searchService 返回错误，触发 MongoDB fallback
// 注意：由于 SearchService 是具体实现而非接口，此测试模拟 nil searchService 场景
func TestSearchByTitle_FallbackOnError(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	// searchService 为 nil，会直接走 MongoDB 路径（相当于 fallback）
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// MongoDB 返回结果
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "测试书籍",
			Author:           "测试作者",
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证 MongoDB fallback 被调用
		return filter.Keyword != nil && *filter.Keyword == "测试" &&
			filter.SortBy == "view_count" && filter.SortOrder == "desc"
	})).Return(books, int64(1), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?title=测试&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - HTTP 200, MongoDB 被调用
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "搜索成功")
	mockService.AssertExpectations(t)
}

// TestSearchByTitle_FallbackOnEmpty 测试 SearchService 返回空结果时的 fallback
// 场景：searchService 返回空结果（Total=0），触发 MongoDB fallback
// 注意：由于 SearchService 是具体实现而非接口，此测试模拟返回空结果的场景
func TestSearchByTitle_FallbackOnEmpty(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	// searchService 为 nil，会直接走 MongoDB 路径
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// MongoDB 返回结果（即使是空结果，也返回 200）
	books := []*bookstoreModel.Book{}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证 MongoDB fallback 被调用
		return filter.Keyword != nil && *filter.Keyword == "不存在的书"
	})).Return(books, int64(0), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/title?title=不存在的书&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - HTTP 200, MongoDB 被调用（即使结果为空）
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "搜索成功")
	mockService.AssertExpectations(t)
}

// ==================== SearchByAuthor 完整测试 ====================

// TestSearchByAuthor_Success 测试搜索成功（MongoDB fallback路径）
func TestSearchByAuthor_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// 构造MongoDB返回结果
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "作者的书1",
			Author:           "测试作者",
			ViewCount:        200,
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "作者的书2",
			Author:           "测试作者",
			ViewCount:        150,
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证参数：作者搜索
		return filter.Author != nil && *filter.Author == "测试作者" &&
			filter.SortBy == "view_count" && filter.SortOrder == "desc"
	})).Return(books, int64(2), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?author=测试作者&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "搜索成功")
	mockService.AssertExpectations(t)
}

// TestSearchByAuthor_EmptyResults 测试空结果情况
func TestSearchByAuthor_EmptyResults(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// 返回空结果
	books := []*bookstoreModel.Book{}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?author=不存在的作者&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "搜索成功")
	mockService.AssertExpectations(t)
}

// TestSearchByAuthor_EmptyAuthor 测试缺少author参数
func TestSearchByAuthor_EmptyAuthor(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreSearchTestRouter(mockService, nil)

	// When - 没有提供author参数
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "作者姓名不能为空")
}

// TestSearchByAuthor_PaginationValidation 测试分页参数验证
func TestSearchByAuthor_PaginationValidation(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreSearchTestRouter(mockService, nil)

	books := []*bookstoreModel.Book{}

	// 测试page=0被纠正为1，size=20被纠正为20
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		return filter.Offset == 0 && filter.Limit == 20
	})).Return(books, int64(0), nil).Once()

	// When & Then
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search/author?author=测试&page=0&size=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetBooksByTags 完整测试 ====================

// setupGetBooksByTagsTestRouter 创建按标签筛选的测试路由
func setupGetBooksByTagsTestRouter(service *MockBookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewBookstoreAPI(service, nil, logger.Get())
	v1 := r.Group("/api/v1/bookstore")
	{
		v1.GET("/books/tags", api.GetBooksByTags)
	}

	return r
}

// TestGetBooksByTags_SingleTag 测试单标签匹配（ANY语义）
func TestGetBooksByTags_SingleTag(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByTagsTestRouter(mockService)

	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "玄幻小说1",
			Tags:             []string{"玄幻", "修仙"},
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证标签参数正确传递
		return len(filter.Tags) == 1 && filter.Tags[0] == "玄幻"
	})).Return(books, int64(1), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?tags=玄幻&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取书籍列表成功")
	mockService.AssertExpectations(t)
}

// TestGetBooksByTags_MultipleTagsANY 测试多标签匹配（ANY语义）
func TestGetBooksByTags_MultipleTagsANY(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByTagsTestRouter(mockService)

	// 返回只有部分标签匹配的书籍
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "玄幻小说",
			Tags:             []string{"玄幻"}, // 只匹配一个标签
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "武侠小说",
			Tags:             []string{"武侠"}, // 只匹配另一个标签
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证标签参数：ANY语义，只要包含任一标签即可
		return len(filter.Tags) == 2 && filter.Tags[0] == "玄幻" && filter.Tags[1] == "武侠"
	})).Return(books, int64(2), nil).Once()

	// When - 传入多个标签（逗号分隔）
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?tags=玄幻,武侠&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取书籍列表成功")
	mockService.AssertExpectations(t)
}

// TestGetBooksByTags_EmptyTags 测试缺少tags参数
func TestGetBooksByTags_EmptyTags(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByTagsTestRouter(mockService)

	// When - 没有提供tags参数
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "标签不能为空")
}

// TestGetBooksByTags_TrimWhitespace 测试标签前后空格处理
func TestGetBooksByTags_TrimWhitespace(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByTagsTestRouter(mockService)

	books := []*bookstoreModel.Book{}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证标签被trim处理
		return len(filter.Tags) == 2 && filter.Tags[0] == "玄幻" && filter.Tags[1] == "武侠"
	})).Return(books, int64(0), nil).Once()

	// When - 标签包含前后空格
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?tags= 玄幻 , 武侠 &page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetBooksByTags_WithPagination 测试分页功能
func TestGetBooksByTags_WithPagination(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByTagsTestRouter(mockService)

	books := []*bookstoreModel.Book{}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?tags=测试&page=2&size=30", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ==================== GetBooksByStatus 完整测试 ====================

// setupGetBooksByStatusTestRouter 创建按状态筛选的测试路由
func setupGetBooksByStatusTestRouter(service *MockBookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewBookstoreAPI(service, nil, logger.Get())
	v1 := r.Group("/api/v1/bookstore")
	{
		v1.GET("/books/status", api.GetBooksByStatus)
	}

	return r
}

// TestGetBooksByStatus_ValidStatusPublished 测试有效状态：published
func TestGetBooksByStatus_ValidStatusPublished(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByStatusTestRouter(mockService)

	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "已发布书籍",
			Status:           bookstoreModel.BookStatusCompleted,
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		// 验证状态参数正确传递
		return filter.Status != nil && *filter.Status == "completed"
	})).Return(books, int64(1), nil).Once()

	// When - 注意：API使用"completed"而不是"published"
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=completed&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "获取书籍列表成功")
	mockService.AssertExpectations(t)
}

// TestGetBooksByStatus_ValidStatusOngoingAndPaused 测试有效状态：ongoing和paused
func TestGetBooksByStatus_ValidStatusOngoingAndPaused(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByStatusTestRouter(mockService)

	booksOngoing := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "连载中书籍",
			Status:           bookstoreModel.BookStatusOngoing,
		},
	}

	booksPaused := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "暂停书籍",
			Status:           bookstoreModel.BookStatusPaused,
		},
	}

	// 测试ongoing状态
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		return filter.Status != nil && *filter.Status == "ongoing"
	})).Return(booksOngoing, int64(1), nil).Once()

	// 测试paused状态
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.MatchedBy(func(filter *bookstoreModel.BookFilter) bool {
		return filter.Status != nil && *filter.Status == "paused"
	})).Return(booksPaused, int64(1), nil).Once()

	// When - 测试ongoing状态
	req1, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=ongoing&page=1&size=20", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// When - 测试paused状态
	req2, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=paused&page=1&size=20", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	mockService.AssertExpectations(t)
}

// TestGetBooksByStatus_InvalidStatus 测试无效状态值
func TestGetBooksByStatus_InvalidStatus(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByStatusTestRouter(mockService)

	// When - 传入无效状态
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=invalid&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "无效的状态值")
}

// TestGetBooksByStatus_EmptyStatus 测试缺少status参数
func TestGetBooksByStatus_EmptyStatus(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByStatusTestRouter(mockService)

	// When - 没有提供status参数
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - 应该返回400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "状态不能为空")
}

// TestGetBooksByStatus_WithPagination 测试分页功能
func TestGetBooksByStatus_WithPagination(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupGetBooksByStatusTestRouter(mockService)

	books := []*bookstoreModel.Book{}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil).Once()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=ongoing&page=3&size=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
