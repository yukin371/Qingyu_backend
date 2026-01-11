package social_test

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
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	socialAPI "Qingyu_backend/api/v1/social"
	"Qingyu_backend/models/social"
	"Qingyu_backend/service/interfaces"
)

// MockBookListService 模拟书单服务接口
type MockBookListService struct {
	mock.Mock
}

// CreateBookList 模拟创建书单
func (m *MockBookListService) CreateBookList(ctx context.Context, userID, userName, userAvatar, title, description, cover, category string, tags []string, isPublic bool) (*social.BookList, error) {
	args := m.Called(ctx, userID, userName, userAvatar, title, description, cover, category, tags, isPublic)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.BookList), args.Error(1)
}

// GetBookLists 模拟获取书单列表
func (m *MockBookListService) GetBookLists(ctx context.Context, page, size int) ([]*social.BookList, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.BookList), args.Get(1).(int64), args.Error(2)
}

// GetBookListByID 模拟获取书单详情
func (m *MockBookListService) GetBookListByID(ctx context.Context, bookListID string) (*social.BookList, error) {
	args := m.Called(ctx, bookListID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.BookList), args.Error(1)
}

// UpdateBookList 模拟更新书单
func (m *MockBookListService) UpdateBookList(ctx context.Context, userID, bookListID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, bookListID, updates)
	return args.Error(0)
}

// DeleteBookList 模拟删除书单
func (m *MockBookListService) DeleteBookList(ctx context.Context, userID, bookListID string) error {
	args := m.Called(ctx, userID, bookListID)
	return args.Error(0)
}

// LikeBookList 模拟点赞书单
func (m *MockBookListService) LikeBookList(ctx context.Context, userID, bookListID string) error {
	args := m.Called(ctx, userID, bookListID)
	return args.Error(0)
}

// ForkBookList 模拟复制书单
func (m *MockBookListService) ForkBookList(ctx context.Context, userID, bookListID string) (*social.BookList, error) {
	args := m.Called(ctx, userID, bookListID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.BookList), args.Error(1)
}

// GetBooksInList 模拟获取书单中的书籍
func (m *MockBookListService) GetBooksInList(ctx context.Context, bookListID string) ([]*social.BookListItem, error) {
	args := m.Called(ctx, bookListID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*social.BookListItem), args.Error(1)
}

// setupBookListTestRouter 设置测试路由
func setupBookListTestRouter(bookListService interfaces.BookListService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId（用于需要认证的端点）
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("username", "testuser")
			c.Set("user_avatar", "")
		}
		c.Next()
	})

	api := socialAPI.NewBookListAPI(bookListService)

	v1 := r.Group("/api/v1/social")
	{
		v1.POST("/booklists", api.CreateBookList)
		v1.GET("/booklists", api.GetBookLists)
		v1.GET("/booklists/:id", api.GetBookListDetail)
		v1.PUT("/booklists/:id", api.UpdateBookList)
		v1.DELETE("/booklists/:id", api.DeleteBookList)
		v1.POST("/booklists/:id/like", api.LikeBookList)
		v1.POST("/booklists/:id/fork", api.ForkBookList)
		v1.GET("/booklists/:id/books", api.GetBooksInList)
	}

	return r
}

// TestBookListAPI_CreateBookList_Success 测试创建书单成功
func TestBookListAPI_CreateBookList_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	expectedBookList := &social.BookList{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		UserName:    "testuser",
		Title:       "我的书单",
		Description: "这是一个测试书单",
		Category:    "玄幻",
		Tags:        []string{"推荐", "经典"},
		IsPublic:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateBookList", mock.Anything, userID, "testuser", "",
		"我的书单", "这是一个测试书单", "", "玄幻",
		[]string{"推荐", "经典"}, true).Return(expectedBookList, nil)

	reqBody := map[string]interface{}{
		"title":       "我的书单",
		"description": "这是一个测试书单",
		"category":    "玄幻",
		"tags":        []string{"推荐", "经典"},
		"is_public":   true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/booklists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusCreated), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_CreateBookList_MissingTitle 测试创建书单缺少标题
func TestBookListAPI_CreateBookList_MissingTitle(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"description": "这是一个测试书单",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/booklists", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
}

// TestBookListAPI_GetBookLists_Success 测试获取书单列表成功
func TestBookListAPI_GetBookLists_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	router := setupBookListTestRouter(mockService, "")

	expectedBookLists := []*social.BookList{
		{
			ID:        primitive.NewObjectID(),
			Title:     "书单1",
			BookCount: 5,
		},
		{
			ID:        primitive.NewObjectID(),
			Title:     "书单2",
			BookCount: 10,
		},
	}

	mockService.On("GetBookLists", mock.Anything, 1, 20).Return(expectedBookLists, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/booklists?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_GetBookListDetail_Success 测试获取书单详情成功
func TestBookListAPI_GetBookListDetail_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	router := setupBookListTestRouter(mockService, "")

	bookListID := primitive.NewObjectID().Hex()
	expectedBookList := &social.BookList{
		ID:        primitive.NewObjectID(),
		Title:     "我的书单",
		BookCount: 5,
	}

	mockService.On("GetBookListByID", mock.Anything, bookListID).Return(expectedBookList, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/booklists/"+bookListID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_UpdateBookList_Success 测试更新书单成功
func TestBookListAPI_UpdateBookList_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	bookListID := primitive.NewObjectID().Hex()
	newTitle := "更新后的书单标题"

	mockService.On("UpdateBookList", mock.Anything, userID, bookListID, mock.Anything).Return(nil)

	reqBody := map[string]interface{}{
		"title": &newTitle,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/social/booklists/"+bookListID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_DeleteBookList_Success 测试删除书单成功
func TestBookListAPI_DeleteBookList_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	bookListID := primitive.NewObjectID().Hex()

	mockService.On("DeleteBookList", mock.Anything, userID, bookListID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/booklists/"+bookListID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_LikeBookList_Success 测试点赞书单成功
func TestBookListAPI_LikeBookList_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	bookListID := primitive.NewObjectID().Hex()

	mockService.On("LikeBookList", mock.Anything, userID, bookListID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/social/booklists/"+bookListID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_ForkBookList_Success 测试复制书单成功
func TestBookListAPI_ForkBookList_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	bookListID := primitive.NewObjectID().Hex()
	expectedBookList := &social.BookList{
		ID:        primitive.NewObjectID(),
		Title:     "复制的书单",
		BookCount: 5,
	}

	mockService.On("ForkBookList", mock.Anything, userID, bookListID).Return(expectedBookList, nil)

	req, _ := http.NewRequest("POST", "/api/v1/social/booklists/"+bookListID+"/fork", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusCreated), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_GetBooksInList_Success 测试获取书单中的书籍成功
func TestBookListAPI_GetBooksInList_Success(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	router := setupBookListTestRouter(mockService, "")

	bookListID := primitive.NewObjectID().Hex()
	expectedBooks := []*social.BookListItem{
		{
			BookID:    "book123",
			BookTitle: "测试书籍1",
		},
		{
			BookID:    "book456",
			BookTitle: "测试书籍2",
		},
	}

	mockService.On("GetBooksInList", mock.Anything, bookListID).Return(expectedBooks, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/booklists/"+bookListID+"/books", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestBookListAPI_UpdateBookList_NoFields 测试更新书单没有字段
func TestBookListAPI_UpdateBookList_NoFields(t *testing.T) {
	// Given
	mockService := new(MockBookListService)
	userID := primitive.NewObjectID().Hex()
	router := setupBookListTestRouter(mockService, userID)

	bookListID := primitive.NewObjectID().Hex()

	reqBody := map[string]interface{}{}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/social/booklists/"+bookListID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
}
