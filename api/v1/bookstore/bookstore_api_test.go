package bookstore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/pkg/logger"
	bookstoreService "Qingyu_backend/service/bookstore"
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
	args := m.Called(ctx, mock.AnythingOfType("*bookstoreModel.BookFilter"))
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

func setupBookstoreTestRouter(service *MockBookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	logger := logger.Get()
	api := NewBookstoreAPI(service, nil, logger)

	v1 := r.Group("/api/v1/bookstore")
	{
		v1.GET("/homepage", api.GetHomepage)
		v1.GET("/books", api.GetBooks)
		v1.GET("/books/:id", api.GetBookByID)
		v1.GET("/books/recommended", api.GetRecommendedBooks)
		v1.GET("/books/featured", api.GetFeaturedBooks)
		v1.GET("/books/search", api.SearchBooks)
		v1.GET("/books/tags", api.GetBooksByTags) // 新增：按标签筛选
		v1.GET("/books/status", api.GetBooksByStatus) // 新增：按状态筛选
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

func TestBookstoreAPI_GetBooksByTags_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?tags=玄幻,仙侠&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooksByTags_MissingTags(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When - Missing tags parameter
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_GetBooksByTags_InvalidPagination(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When - Invalid size (> 100)
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/tags?tags=玄幻&size=150", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should still return 200 (size will be normalized to 20)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooksByStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}
	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=ongoing&page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreAPI_GetBooksByStatus_MissingStatus(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When - Missing status parameter
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_GetBooksByStatus_InvalidStatus(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// When - Invalid status value
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookstoreAPI_GetBooksByStatus_ValidStatuses(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	books := []*bookstoreModel.Book{}

	// Test all valid status values
	validStatuses := []string{"ongoing", "completed", "paused"}

	for _, status := range validStatuses {
		// Reset mock for each iteration
		mockService.ExpectedCalls = nil
		mockService.Calls = nil
		mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

		// When
		req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/status?status="+status, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Then - All valid statuses should return 200
		assert.Equal(t, http.StatusOK, w.Code, "Status "+status+" should be valid")
	}
}
