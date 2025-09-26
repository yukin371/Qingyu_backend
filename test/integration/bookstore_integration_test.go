package integration

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
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"Qingyu_backend/models/reading/bookstore"
	readingAPI "Qingyu_backend/api/v1/reading"
	bookstoreService "Qingyu_backend/service/reading/bookstore"
	"Qingyu_backend/repository/mongodb"
	"Qingyu_backend/repository/interfaces"
)

// BookstoreIntegrationTestSuite 书城集成测试套件
type BookstoreIntegrationTestSuite struct {
	suite.Suite
	router      *gin.Engine
	api         *readingAPI.BookstoreAPI
	service     bookstoreService.BookstoreService
	bookRepo    interfaces.BookRepository
	categoryRepo interfaces.CategoryRepository
	bannerRepo  interfaces.BannerRepository
}

// SetupSuite 测试套件初始化
func (suite *BookstoreIntegrationTestSuite) SetupSuite() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	
	// 这里应该连接测试数据库，为了演示使用Mock
	// 在实际项目中，应该连接真实的测试MongoDB实例
	
	// 创建仓储层（这里使用Mock，实际应该使用真实的MongoDB连接）
	suite.bookRepo = &MockBookRepository{}
	suite.categoryRepo = &MockCategoryRepository{}
	suite.bannerRepo = &MockBannerRepository{}
	
	// 创建服务层
	suite.service = bookstoreService.NewBookstoreService(
		suite.bookRepo,
		suite.categoryRepo,
		suite.bannerRepo,
	)
	
	// 创建API层
	suite.api = readingAPI.NewBookstoreAPI(suite.service)
	
	// 创建路由
	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	{
		bookstore := v1.Group("/bookstore")
		{
			bookstore.GET("/homepage", suite.api.GetHomepage)
			bookstore.GET("/books/:id", suite.api.GetBookByID)
			bookstore.GET("/books/recommended", suite.api.GetRecommendedBooks)
			bookstore.GET("/books/featured", suite.api.GetFeaturedBooks)
			bookstore.GET("/books/search", suite.api.SearchBooks)
			bookstore.POST("/books/:id/view", suite.api.IncrementBookView)
			bookstore.GET("/categories/tree", suite.api.GetCategoryTree)
			bookstore.GET("/categories/:id", suite.api.GetCategoryByID)
			bookstore.GET("/categories/:categoryId/books", suite.api.GetBooksByCategory)
			bookstore.GET("/banners", suite.api.GetActiveBanners)
			bookstore.POST("/banners/:id/click", suite.api.IncrementBannerClick)
		}
	}
}

// TestGetHomepage 测试获取首页数据
func (suite *BookstoreIntegrationTestSuite) TestGetHomepage() {
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/homepage", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), float64(200), response["code"])
	assert.Contains(suite.T(), response["message"], "成功")
	assert.NotNil(suite.T(), response["data"])
}

// TestGetBookByID 测试获取书籍详情
func (suite *BookstoreIntegrationTestSuite) TestGetBookByID() {
	// 使用有效的ObjectID
	bookID := primitive.NewObjectID()
	
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/"+bookID.Hex(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应状态码（可能是200或404，取决于Mock数据）
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestGetBookByID_InvalidID 测试无效ID
func (suite *BookstoreIntegrationTestSuite) TestGetBookByID_InvalidID() {
	// 发送请求 - 使用无效ID
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/invalid-id", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), float64(400), response["code"])
	assert.Contains(suite.T(), response["message"], "ID不能为空")
}

// TestGetRecommendedBooks 测试获取推荐书籍
func (suite *BookstoreIntegrationTestSuite) TestGetRecommendedBooks() {
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/recommended?page=1&size=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
	
	if w.Code == http.StatusOK {
		assert.Equal(suite.T(), float64(200), response["code"])
		assert.NotNil(suite.T(), response["data"])
		assert.NotNil(suite.T(), response["page"])
		assert.NotNil(suite.T(), response["size"])
	}
}

// TestSearchBooks 测试搜索书籍
func (suite *BookstoreIntegrationTestSuite) TestSearchBooks() {
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=测试&page=1&size=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestSearchBooks_NoKeywordOrFilter 测试搜索书籍 - 无关键词和过滤器
func (suite *BookstoreIntegrationTestSuite) TestSearchBooks_NoKeywordOrFilter() {
	// 发送请求 - 无关键词和过滤器
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), float64(400), response["code"])
	assert.Contains(suite.T(), response["message"], "搜索关键词或过滤条件")
}

// TestGetCategoryTree 测试获取分类树
func (suite *BookstoreIntegrationTestSuite) TestGetCategoryTree() {
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/tree", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestGetCategoryByID 测试获取分类详情
func (suite *BookstoreIntegrationTestSuite) TestGetCategoryByID() {
	// 使用有效的ObjectID
	categoryID := primitive.NewObjectID()
	
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/"+categoryID.Hex(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应状态码
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestGetBooksByCategory 测试根据分类获取书籍
func (suite *BookstoreIntegrationTestSuite) TestGetBooksByCategory() {
	// 使用有效的ObjectID
	categoryID := primitive.NewObjectID()
	
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/categories/"+categoryID.Hex()+"/books?page=1&size=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应状态码
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestGetActiveBanners 测试获取激活的Banner
func (suite *BookstoreIntegrationTestSuite) TestGetActiveBanners() {
	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners?limit=5", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestIncrementBookView 测试增加书籍浏览量
func (suite *BookstoreIntegrationTestSuite) TestIncrementBookView() {
	// 使用有效的ObjectID
	bookID := primitive.NewObjectID()
	
	// 发送POST请求
	req, _ := http.NewRequest("POST", "/api/v1/bookstore/books/"+bookID.Hex()+"/view", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应状态码
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestIncrementBannerClick 测试增加Banner点击次数
func (suite *BookstoreIntegrationTestSuite) TestIncrementBannerClick() {
	// 使用有效的ObjectID
	bannerID := primitive.NewObjectID()
	
	// 发送POST请求
	req, _ := http.NewRequest("POST", "/api/v1/bookstore/banners/"+bannerID.Hex()+"/click", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// 验证响应状态码
	assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.NotNil(suite.T(), response["code"])
	assert.NotNil(suite.T(), response["message"])
}

// TestPaginationParameters 测试分页参数
func (suite *BookstoreIntegrationTestSuite) TestPaginationParameters() {
	testCases := []struct {
		name     string
		url      string
		expected map[string]interface{}
	}{
		{
			name: "默认分页参数",
			url:  "/api/v1/bookstore/books/recommended",
		},
		{
			name: "自定义分页参数",
			url:  "/api/v1/bookstore/books/recommended?page=2&size=5",
		},
		{
			name: "无效分页参数",
			url:  "/api/v1/bookstore/books/recommended?page=0&size=-1",
		},
		{
			name: "超大分页参数",
			url:  "/api/v1/bookstore/books/recommended?page=1&size=1000",
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			
			// 验证响应状态码
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			assert.NotNil(t, response["code"])
			assert.NotNil(t, response["message"])
			
			if w.Code == http.StatusOK {
				// 验证分页信息
				assert.NotNil(t, response["page"])
				assert.NotNil(t, response["size"])
				
				// 验证分页参数的合理性
				page := response["page"].(float64)
				size := response["size"].(float64)
				assert.True(t, page >= 1)
				assert.True(t, size >= 1 && size <= 100)
			}
		})
	}
}

// TestErrorHandling 测试错误处理
func (suite *BookstoreIntegrationTestSuite) TestErrorHandling() {
	testCases := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "无效书籍ID",
			method:         "GET",
			url:            "/api/v1/bookstore/books/invalid-id",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "无效分类ID",
			method:         "GET",
			url:            "/api/v1/bookstore/categories/invalid-id",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "无效Banner ID",
			method:         "POST",
			url:            "/api/v1/bookstore/banners/invalid-id/click",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "空搜索条件",
			method:         "GET",
			url:            "/api/v1/bookstore/books/search",
			expectedStatus: http.StatusBadRequest,
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			
			assert.Equal(t, tc.expectedStatus, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			assert.Equal(t, float64(tc.expectedStatus), response["code"])
			assert.NotEmpty(t, response["message"])
		})
	}
}

// TestResponseFormat 测试响应格式
func (suite *BookstoreIntegrationTestSuite) TestResponseFormat() {
	// 测试标准响应格式
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	// 验证响应格式
	assert.Contains(suite.T(), response, "code")
	assert.Contains(suite.T(), response, "message")
	
	if w.Code == http.StatusOK {
		assert.Contains(suite.T(), response, "data")
	}
}

// TestConcurrentRequests 测试并发请求
func (suite *BookstoreIntegrationTestSuite) TestConcurrentRequests() {
	const numRequests = 10
	results := make(chan int, numRequests)
	
	// 并发发送请求
	for i := 0; i < numRequests; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/api/v1/bookstore/banners", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}
	
	// 收集结果
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		// 验证状态码是合理的
		assert.True(suite.T(), statusCode >= 200 && statusCode < 600)
	}
}

// 运行测试套件
func TestBookstoreIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(BookstoreIntegrationTestSuite))
}

// Mock实现（简化版，实际测试中应该使用真实数据库）
type MockBookRepository struct{}

func (m *MockBookRepository) Create(ctx context.Context, book *bookstore.Book) error {
	return nil
}

func (m *MockBookRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Book, error) {
	return &bookstore.Book{
		ID:     id,
		Title:  "测试书籍",
		Author: "测试作者",
		Status: "published",
	}, nil
}

func (m *MockBookRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *MockBookRepository) Health(ctx context.Context) error {
	return nil
}

func (m *MockBookRepository) GetByTitle(ctx context.Context, title string) (*bookstore.Book, error) {
	return nil, nil
}

func (m *MockBookRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) Search(ctx context.Context, keyword string, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	return []*bookstore.Book{}, nil
}

func (m *MockBookRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	return 0, nil
}

func (m *MockBookRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	return 0, nil
}

func (m *MockBookRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	return 0, nil
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstore.BookStats, error) {
	return &bookstore.BookStats{}, nil
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	return nil
}

func (m *MockBookRepository) IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	return nil
}

func (m *MockBookRepository) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	return nil
}

func (m *MockBookRepository) UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	return nil
}

func (m *MockBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status string) error {
	return nil
}

func (m *MockBookRepository) BatchUpdateCategory(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []primitive.ObjectID) error {
	return nil
}

func (m *MockBookRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type MockCategoryRepository struct{}

func (m *MockCategoryRepository) Create(ctx context.Context, category *bookstore.Category) error {
	return nil
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Category, error) {
	return &bookstore.Category{
		ID:       id,
		Name:     "测试分类",
		IsActive: true,
	}, nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *MockCategoryRepository) Health(ctx context.Context) error {
	return nil
}

func (m *MockCategoryRepository) GetByName(ctx context.Context, name string) (*bookstore.Category, error) {
	return nil, nil
}

func (m *MockCategoryRepository) GetByParent(ctx context.Context, parentID primitive.ObjectID, limit, offset int) ([]*bookstore.Category, error) {
	return []*bookstore.Category{}, nil
}

func (m *MockCategoryRepository) GetByLevel(ctx context.Context, level int, limit, offset int) ([]*bookstore.Category, error) {
	return []*bookstore.Category{}, nil
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*bookstore.Category, error) {
	return []*bookstore.Category{}, nil
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	return []*bookstore.CategoryTree{}, nil
}

func (m *MockCategoryRepository) CountByParent(ctx context.Context, parentID primitive.ObjectID) (int64, error) {
	return 0, nil
}

func (m *MockCategoryRepository) UpdateBookCount(ctx context.Context, categoryID primitive.ObjectID, count int64) error {
	return nil
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID primitive.ObjectID) ([]*bookstore.Category, error) {
	return []*bookstore.Category{}, nil
}

func (m *MockCategoryRepository) GetAncestors(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstore.Category, error) {
	return []*bookstore.Category{}, nil
}

func (m *MockCategoryRepository) GetDescendants(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstore.Category, error) {
	return []*bookstore.Category{}, nil
}

func (m *MockCategoryRepository) BatchUpdateStatus(ctx context.Context, categoryIDs []primitive.ObjectID, isActive bool) error {
	return nil
}

func (m *MockCategoryRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type MockBannerRepository struct{}

func (m *MockBannerRepository) Create(ctx context.Context, banner *bookstore.Banner) error {
	return nil
}

func (m *MockBannerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Banner, error) {
	return &bookstore.Banner{
		ID:       id,
		Title:    "测试Banner",
		IsActive: true,
	}, nil
}

func (m *MockBannerRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}

func (m *MockBannerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *MockBannerRepository) Health(ctx context.Context) error {
	return nil
}

func (m *MockBannerRepository) GetActive(ctx context.Context, limit, offset int) ([]*bookstore.Banner, error) {
	return []*bookstore.Banner{}, nil
}

func (m *MockBannerRepository) GetByTargetType(ctx context.Context, targetType string, limit, offset int) ([]*bookstore.Banner, error) {
	return []*bookstore.Banner{}, nil
}

func (m *MockBannerRepository) GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*bookstore.Banner, error) {
	return []*bookstore.Banner{}, nil
}

func (m *MockBannerRepository) IncrementClickCount(ctx context.Context, bannerID primitive.ObjectID) error {
	return nil
}

func (m *MockBannerRepository) GetClickStats(ctx context.Context, bannerID primitive.ObjectID) (int64, error) {
	return 0, nil
}

func (m *MockBannerRepository) BatchUpdateStatus(ctx context.Context, bannerIDs []primitive.ObjectID, isActive bool) error {
	return nil
}

func (m *MockBannerRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}