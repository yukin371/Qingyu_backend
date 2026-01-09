package api

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreAPI "Qingyu_backend/api/v1/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookstoreService 模拟书城服务
type MockBookstoreService struct {
	mock.Mock
}

func (m *MockBookstoreService) GetHomepageData(ctx context.Context) (*bookstoreService.HomepageData, error) {
	args := m.Called(ctx)
	return args.Get(0).(*bookstoreService.HomepageData), args.Error(1)
}

func (m *MockBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore2.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore2.Book), args.Error(1)
}

func (m *MockBookstoreService) GetActiveBanners(ctx context.Context, limit int) ([]*bookstore2.Banner, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore2.Banner), args.Error(1)
}

func (m *MockBookstoreService) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookstoreService) GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookstoreService) SearchBooksWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore2.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetCategoryTree(ctx context.Context) ([]*bookstore2.CategoryTree, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore2.CategoryTree), args.Error(1)
}

func (m *MockBookstoreService) GetCategoryByID(ctx context.Context, id string) (*bookstore2.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore2.Category), args.Error(1)
}

func (m *MockBookstoreService) IncrementBookView(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookstoreService) IncrementBannerClick(ctx context.Context, bannerID string) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBookstoreService) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetRankingByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookstoreService) GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookstoreService) GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookstoreService) SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore2.Book, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	return args.Get(0).([]*bookstore2.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetRootCategories(ctx context.Context) ([]*bookstore2.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockBookstoreService) UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string) error {
	args := m.Called(ctx, rankingType, period)
	return args.Error(0)
}

func (m *MockBookstoreService) GetBookStats(ctx context.Context) (*bookstore2.BookStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*bookstore2.BookStats), args.Error(1)
}

// setupTestRouter 设置测试路由
func setupTestRouter(service bookstoreService.BookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := bookstoreAPI.NewBookstoreAPI(service)

	v1 := router.Group("/api/v1")
	bookstoreGroup := v1.Group("/bookstore")
	{
		// 首页数据
		bookstoreGroup.GET("/homepage", api.GetHomepage)

		// 书籍相关
		books := bookstoreGroup.Group("/books")
		{
			books.GET("/:id", api.GetBookByID)
			books.GET("/recommended", api.GetRecommendedBooks)
			books.GET("/featured", api.GetFeaturedBooks)
			books.GET("/search", api.SearchBooks)
			books.POST("/:id/view", api.IncrementBookView)
		}

		// 榜单相关
		rankings := bookstoreGroup.Group("/rankings")
		{
			rankings.GET("/realtime", api.GetRealtimeRanking)
			rankings.GET("/weekly", api.GetWeeklyRanking)
			rankings.GET("/monthly", api.GetMonthlyRanking)
			rankings.GET("/newbie", api.GetNewbieRanking)
			rankings.GET("/:type", api.GetRankingByType)
		}

		// Banner相关
		banners := bookstoreGroup.Group("/banners")
		{
			banners.GET("", api.GetActiveBanners)
			banners.POST("/:id/click", api.IncrementBannerClick)
		}

		// 分类相关
		categories := bookstoreGroup.Group("/categories")
		{
			categories.GET("/tree", api.GetCategoryTree)
			categories.GET("/:id", api.GetCategoryByID)
		}

		// 分类书籍路由（避免参数冲突）
		categoryBooks := bookstoreGroup.Group("/category-books")
		{
			categoryBooks.GET("/:categoryId", api.GetBooksByCategory)
		}
	}

	return router
}

// TestGetHomepage 测试获取首页数据API
func TestGetHomepage(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	// 准备测试数据
	homepageData := &bookstoreService.HomepageData{
		Banners: []*bookstore2.Banner{
			{
				ID:    primitive.NewObjectID(),
				Title: "Test Banner",
				Image: "https://example.com/banner.jpg",
			},
		},
		RecommendedBooks: []*bookstore2.Book{
			{
				ID:     primitive.NewObjectID(),
				Title:  "Test Book",
				Author: "Test Author",
			},
		},
		Rankings: map[string][]*bookstore2.RankingItem{
			"realtime": {
				{
					ID:   primitive.NewObjectID(),
					Rank: 1,
				},
			},
		},
	}

	mockService.On("GetHomepageData", mock.Anything).Return(homepageData, nil)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/homepage", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取首页数据成功", response.Message)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetRealtimeRanking 测试获取实时榜API
func TestGetRealtimeRankingApi(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	// 准备测试数据
	rankingItems := []*bookstore2.RankingItem{
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore2.RankingTypeRealtime,
			Rank:      1,
			Score:     100.0,
			ViewCount: 1000,
			LikeCount: 50,
			Period:    "2024-01-15",
		},
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore2.RankingTypeRealtime,
			Rank:      2,
			Score:     95.0,
			ViewCount: 800,
			LikeCount: 40,
			Period:    "2024-01-15",
		},
	}

	mockService.On("GetRealtimeRanking", mock.Anything, 20).Return(rankingItems, nil)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/realtime", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取实时榜成功", response.Message)

	// 验证返回的数据
	dataBytes, _ := json.Marshal(response.Data)
	var returnedItems []*bookstore2.RankingItem
	json.Unmarshal(dataBytes, &returnedItems)

	assert.Len(t, returnedItems, 2)
	assert.Equal(t, 1, returnedItems[0].Rank)
	assert.Equal(t, 2, returnedItems[1].Rank)

	mockService.AssertExpectations(t)
}

// TestGetRealtimeRankingWithLimit 测试带限制参数的实时榜API
func TestGetRealtimeRankingWithLimit(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	rankingItems := []*bookstore2.RankingItem{
		{
			ID:   primitive.NewObjectID(),
			Rank: 1,
		},
	}

	mockService.On("GetRealtimeRanking", mock.Anything, 10).Return(rankingItems, nil)

	// 发送带limit参数的请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/realtime?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetWeeklyRanking 测试获取周榜API
func TestGetWeeklyRankingApi(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	rankingItems := []*bookstore2.RankingItem{
		{
			ID:     primitive.NewObjectID(),
			Type:   bookstore2.RankingTypeWeekly,
			Rank:   1,
			Period: "2024-W03",
		},
	}

	mockService.On("GetWeeklyRanking", mock.Anything, "2024-W03", 20).Return(rankingItems, nil)

	// 发送带period参数的请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/weekly?period=2024-W03", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取周榜成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetRankingByType 测试根据类型获取榜单API
func TestGetRankingByType(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	rankingItems := []*bookstore2.RankingItem{
		{
			ID:   primitive.NewObjectID(),
			Type: bookstore2.RankingTypeMonthly,
			Rank: 1,
		},
	}

	mockService.On("GetMonthlyRanking", mock.Anything, "2024-01", 15).Return(rankingItems, nil)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/monthly?period=2024-01&limit=15", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetRankingByType_InvalidType 测试无效榜单类型
func TestGetRankingByType_InvalidType(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	// 发送无效类型的请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "无效的榜单类型", response.Message)

	// 确保没有调用服务
	mockService.AssertNotCalled(t, "GetRankingByType")
}

// TestGetBookByID 测试获取书籍详情API
func TestGetBookByID(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	book := &bookstore2.Book{
		ID:     primitive.NewObjectID(),
		Title:  "Test Book",
		Author: "Test Author",
		Status: "published",
	}

	bookID := book.ID.Hex()
	mockService.On("GetBookByID", mock.Anything, bookID).Return(book, nil)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取书籍详情成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetActiveBanners 测试获取Banner列表API
func TestGetActiveBanners(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	banners := []*bookstore2.Banner{
		{
			ID:       primitive.NewObjectID(),
			Title:    "Banner 1",
			IsActive: true,
		},
		{
			ID:       primitive.NewObjectID(),
			Title:    "Banner 2",
			IsActive: true,
		},
	}

	mockService.On("GetActiveBanners", mock.Anything, 5).Return(banners, nil)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取Banner列表成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestIncrementBannerClick 测试增加Banner点击次数API
func TestIncrementBannerClick(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	bannerID := primitive.NewObjectID().Hex()
	mockService.On("IncrementBannerClick", mock.Anything, bannerID).Return(nil)

	// 发送POST请求
	req, _ := http.NewRequest("POST", "/api/v1/bookstore/banners/"+bannerID+"/click", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "点击次数增加成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestSearchBooks 测试搜索书籍API
func TestSearchBooks(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	books := []*bookstore2.Book{
		{
			ID:     primitive.NewObjectID(),
			Title:  "Search Result Book",
			Author: "Search Author",
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.AnythingOfType("*bookstore.BookFilter")).Return(books, int64(1), nil)

	// 发送搜索请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "搜索书籍成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestAPIErrorHandling 测试API错误处理
func TestAPIErrorHandling(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	// 模拟服务返回错误
	mockService.On("GetRealtimeRanking", mock.Anything, 20).Return([]*bookstore2.RankingItem(nil), assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/realtime", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Contains(t, response.Message, "获取实时榜失败")

	mockService.AssertExpectations(t)
}

// TestGetBooksByCategory 测试根据分类获取书籍列表API
func TestGetBooksByCategory(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	categoryID := primitive.NewObjectID().Hex()
	books := []*bookstore2.Book{
		{
			ID:     primitive.NewObjectID(),
			Title:  "Category Book 1",
			Author: "Author 1",
		},
		{
			ID:     primitive.NewObjectID(),
			Title:  "Category Book 2",
			Author: "Author 2",
		},
	}

	mockService.On("GetBooksByCategory", mock.Anything, categoryID, 1, 20).Return(books, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/category-books/"+categoryID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取分类书籍成功", response.Message)
	assert.Equal(t, int64(2), response.Total)

	mockService.AssertExpectations(t)
}

// TestGetBooksByCategory_InvalidID 测试无效分类ID
func TestGetBooksByCategory_InvalidID(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	// 使用无效的ObjectID格式
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/category-books/invalid-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "分类ID格式无效", response.Message)

	mockService.AssertNotCalled(t, "GetBooksByCategory")
}

// TestGetRecommendedBooks 测试获取推荐书籍API
func TestGetRecommendedBooks(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	books := []*bookstore2.Book{
		{
			ID:     primitive.NewObjectID(),
			Title:  "Recommended Book 1",
			Author: "Author 1",
		},
	}

	mockService.On("GetRecommendedBooks", mock.Anything, 1, 20).Return(books, nil)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/recommended", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取推荐书籍成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetFeaturedBooks 测试获取精选书籍API
func TestGetFeaturedBooks(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	books := []*bookstore2.Book{
		{
			ID:     primitive.NewObjectID(),
			Title:  "Featured Book 1",
			Author: "Author 1",
		},
	}

	mockService.On("GetFeaturedBooks", mock.Anything, 1, 20).Return(books, nil)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/featured", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取精选书籍成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetCategoryTree 测试获取分类树API
func TestGetCategoryTree(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	tree := []*bookstore2.CategoryTree{
		{
			Category: bookstore2.Category{
				ID:   primitive.NewObjectID(),
				Name: "分类1",
			},
			Children: []*bookstore2.CategoryTree{},
		},
	}

	mockService.On("GetCategoryTree", mock.Anything).Return(tree, nil)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/tree", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取分类树成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetCategoryByID 测试获取分类详情API
func TestGetCategoryByID(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	categoryID := primitive.NewObjectID().Hex()
	category := &bookstore2.Category{
		ID:   primitive.NewObjectID(),
		Name: "测试分类",
	}

	mockService.On("GetCategoryByID", mock.Anything, categoryID).Return(category, nil)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/"+categoryID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取分类详情成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestIncrementBookView 测试增加书籍浏览量API
func TestIncrementBookView(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	bookID := primitive.NewObjectID().Hex()
	mockService.On("IncrementBookView", mock.Anything, bookID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/bookstore/books/"+bookID+"/view", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response bookstoreAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "浏览量增加成功", response.Message)

	mockService.AssertExpectations(t)
}
