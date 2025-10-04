package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readingAPI "Qingyu_backend/api/v1/reading"
	"Qingyu_backend/models/reading/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// setupTestRouter 设置测试路由
func setupTestRouter(service bookstoreService.BookstoreService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := readingAPI.NewBookstoreAPI(service)

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
			categories.GET("/:categoryId/books", api.GetBooksByCategory)
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
		Banners: []*bookstore.Banner{
			{
				ID:    primitive.NewObjectID(),
				Title: "Test Banner",
				Image: "https://example.com/banner.jpg",
			},
		},
		RecommendedBooks: []*bookstore.Book{
			{
				ID:     primitive.NewObjectID(),
				Title:  "Test Book",
				Author: "Test Author",
			},
		},
		Rankings: map[string][]*bookstore.RankingItem{
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

	var response readingAPI.APIResponse
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
	rankingItems := []*bookstore.RankingItem{
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeRealtime,
			Rank:      1,
			Score:     100.0,
			ViewCount: 1000,
			LikeCount: 50,
			Period:    "2024-01-15",
		},
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeRealtime,
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

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取实时榜成功", response.Message)

	// 验证返回的数据
	dataBytes, _ := json.Marshal(response.Data)
	var returnedItems []*bookstore.RankingItem
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

	rankingItems := []*bookstore.RankingItem{
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

	rankingItems := []*bookstore.RankingItem{
		{
			ID:     primitive.NewObjectID(),
			Type:   bookstore.RankingTypeWeekly,
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

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取周榜成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetRankingByType 测试根据类型获取榜单API
func TestGetRankingByType(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	rankingItems := []*bookstore.RankingItem{
		{
			ID:   primitive.NewObjectID(),
			Type: bookstore.RankingTypeMonthly,
			Rank: 1,
		},
	}

	mockService.On("GetRankingByType", mock.Anything, bookstore.RankingTypeMonthly, "2024-01", 15).Return(rankingItems, nil)

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

	var response readingAPI.APIResponse
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

	book := &bookstore.Book{
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

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取书籍详情成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestGetActiveBanners 测试获取Banner列表API
func TestGetActiveBanners(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	banners := []*bookstore.Banner{
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

	var response readingAPI.APIResponse
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

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "点击次数增加成功", response.Message)

	mockService.AssertExpectations(t)
}

// TestSearchBooks 测试搜索书籍API
func TestSearchBooks(t *testing.T) {
	mockService := new(MockBookstoreService)
	router := setupTestRouter(mockService)

	books := []*bookstore.Book{
		{
			ID:     primitive.NewObjectID(),
			Title:  "Search Result Book",
			Author: "Search Author",
		},
	}

	mockService.On("SearchBooks", mock.Anything, "test", mock.AnythingOfType("*bookstore.BookFilter")).Return(books, nil)

	// 发送搜索请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response readingAPI.PaginatedResponse
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
	mockService.On("GetRealtimeRanking", mock.Anything, 20).Return([]*bookstore.RankingItem(nil), assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/bookstore/rankings/realtime", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response readingAPI.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Contains(t, response.Message, "获取实时榜失败")

	mockService.AssertExpectations(t)
}
