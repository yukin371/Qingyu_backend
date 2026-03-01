package admin

import (
	"bytes"
	"context"
	"encoding/json"
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
