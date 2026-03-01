package admin

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

	"Qingyu_backend/models/bookstore"
	adminsvc "Qingyu_backend/service/admin"
)

// MockCategoryAdminService Mock服务
type MockCategoryAdminService struct {
	mock.Mock
}

func (m *MockCategoryAdminService) CreateCategory(ctx context.Context, req *adminsvc.CreateCategoryRequest) (*bookstore.Category, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminService) UpdateCategory(ctx context.Context, id string, req *adminsvc.UpdateCategoryRequest) (*bookstore.Category, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminService) DeleteCategory(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryAdminService) GetCategoryTree(ctx context.Context) ([]*adminsvc.CategoryTreeNode, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*adminsvc.CategoryTreeNode), args.Error(1)
}

func (m *MockCategoryAdminService) GetCategories(ctx context.Context, filter *adminsvc.CategoryFilter) ([]*bookstore.Category, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminService) GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminService) MoveCategory(ctx context.Context, id string, req *adminsvc.MoveCategoryRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockCategoryAdminService) SortCategory(ctx context.Context, id string, sortOrder int) error {
	args := m.Called(ctx, id, sortOrder)
	return args.Error(0)
}

// setupCategoryTestRouter 设置分类测试路由
func setupCategoryTestRouter(service adminsvc.CategoryAdminService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := NewCategoryAdminAPI(service)
	v1 := router.Group("/api/v1/admin/categories")
	{
		v1.POST("", api.CreateCategory)
		v1.GET("", api.GetCategories)
		v1.GET("/tree", api.GetCategoryTree)
		v1.GET("/:id", api.GetCategoryByID)
		v1.PUT("/:id", api.UpdateCategory)
		v1.DELETE("/:id", api.DeleteCategory)
		v1.PUT("/:id/move", api.MoveCategory)
		v1.PUT("/:id/sort", api.SortCategory)
	}

	return router
}

// 测试创建分类成功
func TestCategoryAdminAPI_CreateCategory_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.CreateCategoryRequest{
		Name:      "玄幻",
		SortOrder: 1,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	expectedCategory := &bookstore.Category{
		ID:   "123",
		Name: "玄幻",
	}

	mockService.On("CreateCategory", mock.Anything, &reqBody).Return(expectedCategory, nil)

	req, _ := http.NewRequest("POST", "/api/v1/admin/categories", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// 测试创建分类 - 参数错误
func TestCategoryAdminAPI_CreateCategory_InvalidJSON(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	req, _ := http.NewRequest("POST", "/api/v1/admin/categories", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试创建分类 - 服务错误
func TestCategoryAdminAPI_CreateCategory_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.CreateCategoryRequest{
		Name:      "玄幻",
		SortOrder: 1,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.On("CreateCategory", mock.Anything, &reqBody).Return(nil, errors.New("service error"))

	req, _ := http.NewRequest("POST", "/api/v1/admin/categories", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试获取分类列表成功
func TestCategoryAdminAPI_GetCategories_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	expectedCategories := []*bookstore.Category{
		{ID: "1", Name: "玄幻"},
		{ID: "2", Name: "武侠"},
	}

	mockService.On("GetCategories", mock.Anything, mock.MatchedBy(func(filter *adminsvc.CategoryFilter) bool {
		return filter != nil
	})).Return(expectedCategories, nil)

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// 测试获取分类列表 - 带筛选参数
func TestCategoryAdminAPI_GetCategories_WithFilter(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	expectedCategories := []*bookstore.Category{
		{ID: "1", Name: "玄幻"},
	}

	mockService.On("GetCategories", mock.Anything, mock.MatchedBy(func(filter *adminsvc.CategoryFilter) bool {
		return filter != nil && filter.ParentID != nil && *filter.ParentID == "123"
	})).Return(expectedCategories, nil)

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories?parent_id=123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试获取分类列表 - 带level筛选参数
func TestCategoryAdminAPI_GetCategories_WithLevelFilter(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	expectedCategories := []*bookstore.Category{
		{ID: "1", Name: "玄幻"},
	}

	mockService.On("GetCategories", mock.Anything, mock.MatchedBy(func(filter *adminsvc.CategoryFilter) bool {
		return filter != nil && filter.Level != nil && *filter.Level == 2
	})).Return(expectedCategories, nil)

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories?level=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// 测试获取分类列表 - level参数格式错误
func TestCategoryAdminAPI_GetCategories_InvalidLevel(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories?level=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试获取分类列表 - 服务错误
func TestCategoryAdminAPI_GetCategories_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	mockService.On("GetCategories", mock.Anything, mock.MatchedBy(func(filter *adminsvc.CategoryFilter) bool {
		return filter != nil
	})).Return(nil, errors.New("service error"))

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试获取分类树成功
func TestCategoryAdminAPI_GetCategoryTree_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	expectedTree := []*adminsvc.CategoryTreeNode{
		{
			Category: &bookstore.Category{ID: "1", Name: "小说"},
			Children: []*adminsvc.CategoryTreeNode{
				{Category: &bookstore.Category{ID: "2", Name: "玄幻"}},
			},
		},
	}

	mockService.On("GetCategoryTree", mock.Anything).Return(expectedTree, nil)

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories/tree", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// 测试获取分类树 - 服务错误
func TestCategoryAdminAPI_GetCategoryTree_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	mockService.On("GetCategoryTree", mock.Anything).Return(nil, errors.New("service error"))

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories/tree", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试获取分类详情成功
func TestCategoryAdminAPI_GetCategoryByID_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	expectedCategory := &bookstore.Category{
		ID:   "123",
		Name: "玄幻",
	}

	mockService.On("GetCategoryByID", mock.Anything, "123").Return(expectedCategory, nil)

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// 测试获取分类详情 - 不存在
func TestCategoryAdminAPI_GetCategoryByID_NotFound(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	mockService.On("GetCategoryByID", mock.Anything, "999").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/api/v1/admin/categories/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// 测试获取分类详情 - 空ID参数（通过特殊字符模拟）
func TestCategoryAdminAPI_GetCategoryByID_EmptyParam(t *testing.T) {
	mockService := new(MockCategoryAdminService)

	// 创建一个上下文并手动设置空参数
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	api := &CategoryAdminAPI{categoryService: mockService}
	api.GetCategoryByID(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

// 测试更新分类成功
func TestCategoryAdminAPI_UpdateCategory_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.UpdateCategoryRequest{
		Name:       strPtr("玄幻小说"),
		SortOrder:  intPtr(2),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	expectedCategory := &bookstore.Category{
		ID:        "123",
		Name:      "玄幻小说",
		SortOrder: 2,
	}

	mockService.On("UpdateCategory", mock.Anything, "123", &reqBody).Return(expectedCategory, nil)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// 测试更新分类 - 服务错误
func TestCategoryAdminAPI_UpdateCategory_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.UpdateCategoryRequest{
		Name: strPtr("玄幻"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.On("UpdateCategory", mock.Anything, "123", &reqBody).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试更新分类 - 路由不存在
func TestCategoryAdminAPI_UpdateCategory_NoRoute(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.UpdateCategoryRequest{
		Name: strPtr("玄幻"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试更新分类 - 空ID参数
func TestCategoryAdminAPI_UpdateCategory_EmptyParam(t *testing.T) {
	mockService := new(MockCategoryAdminService)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	api := &CategoryAdminAPI{categoryService: mockService}
	api.UpdateCategory(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

// 测试更新分类 - 参数错误
func TestCategoryAdminAPI_UpdateCategory_InvalidJSON(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试删除分类成功
func TestCategoryAdminAPI_DeleteCategory_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	mockService.On("DeleteCategory", mock.Anything, "123").Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/admin/categories/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试删除分类 - 路由不存在
func TestCategoryAdminAPI_DeleteCategory_NoRoute(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	req, _ := http.NewRequest("DELETE", "/api/v1/admin/categories/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试删除分类 - 空ID参数
func TestCategoryAdminAPI_DeleteCategory_EmptyParam(t *testing.T) {
	mockService := new(MockCategoryAdminService)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	api := &CategoryAdminAPI{categoryService: mockService}
	api.DeleteCategory(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

// 测试删除分类 - 服务错误
func TestCategoryAdminAPI_DeleteCategory_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	mockService.On("DeleteCategory", mock.Anything, "123").Return(errors.New("has children"))

	req, _ := http.NewRequest("DELETE", "/api/v1/admin/categories/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试移动分类成功
func TestCategoryAdminAPI_MoveCategory_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.MoveCategoryRequest{
		ParentID: strPtr("456"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.On("MoveCategory", mock.Anything, "123", &reqBody).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123/move", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试移动分类 - 服务错误
func TestCategoryAdminAPI_MoveCategory_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.MoveCategoryRequest{
		ParentID: strPtr("456"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.On("MoveCategory", mock.Anything, "123", &reqBody).Return(errors.New("circular reference"))

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123/move", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试移动分类 - 空ID
func TestCategoryAdminAPI_MoveCategory_EmptyID(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := adminsvc.MoveCategoryRequest{
		ParentID: strPtr("456"),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories//move", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试移动分类 - 参数错误
func TestCategoryAdminAPI_MoveCategory_InvalidJSON(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123/move", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试调整排序成功
func TestCategoryAdminAPI_SortCategory_Success(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := struct {
		SortOrder int `json:"sort_order"`
	}{
		SortOrder: 10,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.On("SortCategory", mock.Anything, "123", 10).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123/sort", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试调整排序 - 服务错误
func TestCategoryAdminAPI_SortCategory_ServiceError(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := struct {
		SortOrder int `json:"sort_order"`
	}{
		SortOrder: 10,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.On("SortCategory", mock.Anything, "123", 10).Return(errors.New("not found"))

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123/sort", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["code"])

	mockService.AssertExpectations(t)
}

// 测试调整排序 - 空ID
func TestCategoryAdminAPI_SortCategory_EmptyID(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := struct {
		SortOrder int `json:"sort_order"`
	}{
		SortOrder: 10,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories//sort", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试调整排序 - 参数错误
func TestCategoryAdminAPI_SortCategory_MissingSortOrder(t *testing.T) {
	mockService := new(MockCategoryAdminService)
	router := setupCategoryTestRouter(mockService)

	reqBody := struct {
		SortOrder int `json:"sort_order"`
	}{
		SortOrder: 0,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/123/sort", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 辅助函数
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
