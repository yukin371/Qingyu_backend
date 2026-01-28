package search

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/models/bookstore"
	searchengine "Qingyu_backend/service/search/engine"
	searchprovider "Qingyu_backend/service/search/provider"
	searchService "Qingyu_backend/service/search"
)

// MockEngine 是一个简单的模拟搜索引擎实现
type MockEngine struct{}

func (m *MockEngine) Search(ctx context.Context, index string, query interface{}, opts *searchengine.SearchOptions) (*searchengine.SearchResult, error) {
	// 返回模拟数据
	mockBook := bson.M{
		"_id":          primitive.NewObjectID(),
		"title":        "修仙世界",
		"author":       "张三",
		"introduction": "这是一个修仙世界的故事",
		"status":       "ongoing",
		"word_count":   100000,
		"view_count":   1000,
		"rating":       8.5,
		"created_at":   "2024-01-01T00:00:00Z",
		"updated_at":   "2024-01-01T00:00:00Z",
	}

	return &searchengine.SearchResult{
		Total: 1,
		Hits: []searchengine.Hit{
			{
				ID:     mockBook["_id"].(primitive.ObjectID).Hex(),
				Score:  1.0,
				Source: mockBook,
			},
		},
	}, nil
}

func (m *MockEngine) Index(ctx context.Context, index string, documents []searchengine.Document) error {
	return nil
}

func (m *MockEngine) Update(ctx context.Context, index string, id string, document searchengine.Document) error {
	return nil
}

func (m *MockEngine) Delete(ctx context.Context, index string, id string) error {
	return nil
}

func (m *MockEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	return nil
}

func (m *MockEngine) Health(ctx context.Context) error {
	return nil
}

// setupTestAPI 设置测试环境
func setupTestAPI() (*SearchAPI, *gin.Engine) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建 Mock Engine
	mockEngine := &MockEngine{}

	// 创建 BookProvider
	bookProviderConfig := &searchprovider.BookProviderConfig{
		AllowedStatuses: []string{
			string(bookstore.BookStatusOngoing),
			string(bookstore.BookStatusCompleted),
		},
		AllowedPrivacy: []bool{false}, // 只允许公开书籍
	}
	bookProvider, err := searchprovider.NewBookProvider(mockEngine, bookProviderConfig)
	if err != nil {
		panic(err)
	}

	// 创建 SearchService
	searchConfig := &searchService.Config{
		EnableCache:           false, // 测试时禁用缓存
		DefaultCacheTTL:       300,
		MaxConcurrentSearches: 10,
	}

	// 创建标准日志记录器
	stdLogger := log.Default()
	zapLogger := zap.NewNop() // 测试时使用无输出日志

	// 创建灰度决策器（测试环境禁用灰度）
	grayScaleConfig := &searchService.GrayScaleConfig{
		Enabled: false, // 测试环境禁用灰度
		Percent: 0,
	}
	grayScaleDecision := searchService.NewGrayScaleDecision(grayScaleConfig, zapLogger)

	searchSvc := searchService.NewSearchService(stdLogger, searchConfig, grayScaleDecision)

	// 注册 BookProvider
	searchSvc.RegisterProvider(bookProvider)

	// 创建 SearchAPI
	searchAPI := NewSearchAPI(searchSvc)

	// 创建 Gin router
	router := gin.New()
	router.Use(gin.Recovery())

	return searchAPI, router
}

// TestSearch_Success 测试成功搜索
func TestSearch_Success(t *testing.T) {
	searchAPI, router := setupTestAPI()

	// 注册路由
	router.POST("/api/v1/search/search", searchAPI.Search)

	// 构造请求
	reqBody := SearchRequest{
		Type:     "books",
		Query:    "修仙",
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "搜索成功", response["message"])
	assert.NotNil(t, response["data"])

	// 验证数据结构
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "books", data["type"])
	assert.NotNil(t, data["total"])
	assert.NotNil(t, data["results"])
	assert.NotNil(t, data["took"])
	assert.NotNil(t, response["request_id"])
}

// TestSearch_EmptyType 测试空搜索类型
func TestSearch_EmptyType(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	// 发送空字符串（会被 Gin 的 required 验证拦截）
	reqBody := map[string]interface{}{
		"type":  "",
		"query": "修仙",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	// Gin 的 required 验证会返回 "Field validation for 'Type' failed on the 'required' tag"
	assert.Contains(t, response["error"], "required")
}

// TestSearch_EmptyQuery 测试空搜索关键词
func TestSearch_EmptyQuery(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	// 发送空字符串（会被 Gin 的 required 验证拦截）
	reqBody := map[string]interface{}{
		"type":  "books",
		"query": "",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	// Gin 的 required 验证会返回 "Field validation for 'Query' failed on the 'required' tag"
	assert.Contains(t, response["error"], "required")
}

// TestSearch_WithFilter 测试带过滤条件的搜索
func TestSearch_WithFilter(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:  "books",
		Query: "修仙",
		Filter: map[string]interface{}{
			"status": "ongoing",
		},
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])
}

// TestSearch_WithSort 测试带排序的搜索
func TestSearch_WithSort(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:  "books",
		Query: "修仙",
		Sort: []SortFieldRequest{
			{
				Field:     "created_at",
				Ascending: false,
			},
		},
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])
}

// TestSearch_Pagination 测试分页功能
func TestSearch_Pagination(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	testCases := []struct {
		name     string
		page     int
		pageSize int
	}{
		{"第一页", 1, 10},
		{"第二页", 2, 20},
		{"大页码", 5, 50},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := SearchRequest{
				Type:     "books",
				Query:    "修仙",
				Page:     tc.page,
				PageSize: tc.pageSize,
			}

			body, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			data := response["data"].(map[string]interface{})
			assert.Equal(t, float64(tc.page), data["page"])
			assert.Equal(t, float64(tc.pageSize), data["page_size"])
		})
	}
}

// TestSearchBatch_Success 测试批量搜索
func TestSearchBatch_Success(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/batch", searchAPI.SearchBatch)

	reqBody := BatchSearchRequest{
		Queries: []SearchRequest{
			{
				Type:     "books",
				Query:    "修仙",
				Page:     1,
				PageSize: 10,
			},
			{
				Type:     "books",
				Query:    "玄幻",
				Page:     1,
				PageSize: 10,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "批量搜索成功", response["message"])

	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// 验证每个搜索结果
	for _, result := range data {
		resultMap := result.(map[string]interface{})
		assert.NotNil(t, resultMap["success"])
		assert.NotNil(t, resultMap["data"])
	}
}

// TestSearchBatch_EmptyQueries 测试空查询列表
func TestSearchBatch_EmptyQueries(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/batch", searchAPI.SearchBatch)

	reqBody := BatchSearchRequest{
		Queries: []SearchRequest{},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	// Gin 的 min 验证会返回 "Field validation for 'Queries' failed on the 'min' tag"
	assert.Contains(t, response["error"], "min")
}

// TestSearchBatch_InvalidQuery 测试批量搜索中的无效查询
func TestSearchBatch_InvalidQuery(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/batch", searchAPI.SearchBatch)

	reqBody := BatchSearchRequest{
		Queries: []SearchRequest{
			{
				Type:     "books",
				Query:    "修仙",
				Page:     1,
				PageSize: 10,
			},
			{
				Type:  "", // 无效：空类型
				Query: "玄幻",
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
}

// TestHealth_Success 测试健康检查
func TestHealth_Success(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.GET("/api/v1/search/health", searchAPI.Health)

	req, _ := http.NewRequest("GET", "/api/v1/search/health", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "健康检查完成", response["message"])

	// 验证健康状态数据
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["status"])
	assert.NotNil(t, data["books"])
}

// TestSearch_UnsupportedType 测试不支持的搜索类型
func TestSearch_UnsupportedType(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:     "projects", // 未注册的类型
		Query:    "测试",
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 当 Provider 不存在时，API 返回 HTTP 200，但 code 字段为 400
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应格式 - HTTP 200 但业务 code 为 400
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Unsupported search type")
}

// TestSearch_DefaultPagination 测试默认分页参数
func TestSearch_DefaultPagination(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:  "books",
		Query: "修仙",
		// 不设置 page 和 page_size
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	// 验证默认分页参数
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(20), data["page_size"])
}

// TestSearch_MaxPageSize 测试最大每页数量限制
func TestSearch_MaxPageSize(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:     "books",
		Query:    "修仙",
		Page:     1,
		PageSize: 200, // 超过最大限制
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	// 验证被限制为最大值 100
	assert.Equal(t, float64(100), data["page_size"])
}

// TestSearch_ComplexFilter 测试复杂过滤条件
func TestSearch_ComplexFilter(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:  "books",
		Query: "修仙",
		Filter: map[string]interface{}{
			"status":         "ongoing",
			"word_count_min": 50000,
			"rating_min":     7.0,
			"tags":           []string{"热血", "冒险"},
		},
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])
}

// TestSearch_MultipleSortFields 测试多字段排序
func TestSearch_MultipleSortFields(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:  "books",
		Query: "修仙",
		Sort: []SortFieldRequest{
			{Field: "rating", Ascending: false},
			{Field: "created_at", Ascending: false},
		},
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])
}

// TestSearch_ResponseStructure 测试响应结构完整性
func TestSearch_ResponseStructure(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	reqBody := SearchRequest{
		Type:     "books",
		Query:    "修仙",
		Page:     1,
		PageSize: 10,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "test-request-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证顶层响应字段
	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "request_id")

	// 验证数据字段
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "type")
	assert.Contains(t, data, "total")
	assert.Contains(t, data, "page")
	assert.Contains(t, data, "page_size")
	assert.Contains(t, data, "results")
	assert.Contains(t, data, "took")

	// 验证 request_id 存在（SearchService 会生成新的 request_id）
	assert.NotEmpty(t, response["request_id"])
}

// TestSearch_InvalidJSON 测试无效 JSON
func TestSearch_InvalidJSON(t *testing.T) {
	searchAPI, router := setupTestAPI()

	router.POST("/api/v1/search/search", searchAPI.Search)

	req, _ := http.NewRequest("POST", "/api/v1/search/search", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "参数错误")
}
