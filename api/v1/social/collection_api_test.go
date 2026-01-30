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
	socialModels "Qingyu_backend/models/social"
	"Qingyu_backend/service/interfaces"
)

// MockCollectionService 模拟收藏服务接口
type MockCollectionService struct {
	mock.Mock
}

// AddToCollection 模拟添加收藏
func (m *MockCollectionService) AddToCollection(ctx context.Context, userID, bookID, folderID, note string, tags []string, isPublic bool) (*socialModels.Collection, error) {
	args := m.Called(ctx, userID, bookID, folderID, note, tags, isPublic)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*socialModels.Collection), args.Error(1)
}

// RemoveFromCollection 模拟移除收藏
func (m *MockCollectionService) RemoveFromCollection(ctx context.Context, userID, collectionID string) error {
	args := m.Called(ctx, userID, collectionID)
	return args.Error(0)
}

// UpdateCollection 模拟更新收藏
func (m *MockCollectionService) UpdateCollection(ctx context.Context, userID, collectionID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, collectionID, updates)
	return args.Error(0)
}

// GetUserCollections 模拟获取用户收藏列表
func (m *MockCollectionService) GetUserCollections(ctx context.Context, userID, folderID string, page, size int) ([]*socialModels.Collection, int64, error) {
	args := m.Called(ctx, userID, folderID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*socialModels.Collection), args.Get(1).(int64), args.Error(2)
}

// GetCollectionsByTag 模拟按标签获取收藏
func (m *MockCollectionService) GetCollectionsByTag(ctx context.Context, userID, tag string, page, size int) ([]*socialModels.Collection, int64, error) {
	args := m.Called(ctx, userID, tag, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*socialModels.Collection), args.Get(1).(int64), args.Error(2)
}

// IsCollected 模拟检查是否已收藏
func (m *MockCollectionService) IsCollected(ctx context.Context, userID, bookID string) (bool, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Bool(0), args.Error(1)
}

// CreateFolder 模拟创建收藏夹
func (m *MockCollectionService) CreateFolder(ctx context.Context, userID, name, description string, isPublic bool) (*socialModels.CollectionFolder, error) {
	args := m.Called(ctx, userID, name, description, isPublic)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*socialModels.CollectionFolder), args.Error(1)
}

// GetUserFolders 模拟获取用户收藏夹列表
func (m *MockCollectionService) GetUserFolders(ctx context.Context, userID string) ([]*socialModels.CollectionFolder, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*socialModels.CollectionFolder), args.Error(1)
}

// UpdateFolder 模拟更新收藏夹
func (m *MockCollectionService) UpdateFolder(ctx context.Context, userID, folderID string, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, folderID, updates)
	return args.Error(0)
}

// DeleteFolder 模拟删除收藏夹
func (m *MockCollectionService) DeleteFolder(ctx context.Context, userID, folderID string) error {
	args := m.Called(ctx, userID, folderID)
	return args.Error(0)
}

// ShareCollection 模拟分享收藏
func (m *MockCollectionService) ShareCollection(ctx context.Context, userID, collectionID string) error {
	args := m.Called(ctx, userID, collectionID)
	return args.Error(0)
}

// UnshareCollection 模拟取消分享收藏
func (m *MockCollectionService) UnshareCollection(ctx context.Context, userID, collectionID string) error {
	args := m.Called(ctx, userID, collectionID)
	return args.Error(0)
}

// GetPublicCollections 模拟获取公开收藏列表
func (m *MockCollectionService) GetPublicCollections(ctx context.Context, page, size int) ([]*socialModels.Collection, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*socialModels.Collection), args.Get(1).(int64), args.Error(2)
}

// GetUserCollectionStats 模拟获取用户收藏统计
func (m *MockCollectionService) GetUserCollectionStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// setupCollectionTestRouter 设置测试路由
func setupCollectionTestRouter(collectionService interfaces.CollectionService, userID string) *gin.Engine {
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

	api := socialAPI.NewCollectionAPI(collectionService)

	v1 := r.Group("/api/v1/reader")
	{
		v1.POST("/collections", api.AddCollection)
		v1.GET("/collections", api.GetCollections)
		v1.GET("/collections/check/:book_id", api.CheckCollected)
		v1.GET("/collections/tag/:tag", api.GetCollectionsByTag)
		v1.PUT("/collections/:id", api.UpdateCollection)
		v1.DELETE("/collections/:id", api.DeleteCollection)
		v1.POST("/collections/:id/share", api.ShareCollection)
		v1.DELETE("/collections/:id/share", api.UnshareCollection)

		v1.POST("/folders", api.CreateFolder)
		v1.GET("/folders", api.GetFolders)
		v1.PUT("/folders/:id", api.UpdateFolder)
		v1.DELETE("/folders/:id", api.DeleteFolder)

		v1.GET("/collections/stats", api.GetCollectionStats)
	}

	v1Public := r.Group("/api/v1")
	{
		v1Public.GET("/public/collections", api.GetPublicCollections)
	}

	return r
}

// TestCollectionAPI_AddCollection_Success 测试添加收藏成功
func TestCollectionAPI_AddCollection_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	expectedCollection := &socialModels.Collection{
		UserID: userID,
		BookID: bookID,
		Note:   "很好看的书",
	}
	expectedCollection.ID = primitive.NewObjectID()
	expectedCollection.CreatedAt = time.Now()

	mockService.On("AddToCollection", mock.Anything, userID, bookID, "", "很好看的书", []string{"推荐"}, true).
		Return(expectedCollection, nil)

	reqBody := map[string]interface{}{
		"book_id":   bookID,
		"note":      "很好看的书",
		"tags":      []string{"推荐"},
		"is_public": true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/collections", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_AddCollection_MissingBookID 测试添加收藏缺少book_id
func TestCollectionAPI_AddCollection_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"note": "很好看的书",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/collections", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 参数错误code为1001
}

// TestCollectionAPI_GetCollections_Success 测试获取收藏列表成功
func TestCollectionAPI_GetCollections_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	expectedCollections := []*socialModels.Collection{
		func() *socialModels.Collection {
			c := &socialModels.Collection{BookID: primitive.NewObjectID().Hex(), Note: "收藏1"}
			c.ID = primitive.NewObjectID()
			return c
		}(),
		func() *socialModels.Collection {
			c := &socialModels.Collection{BookID: primitive.NewObjectID().Hex(), Note: "收藏2"}
			c.ID = primitive.NewObjectID()
			return c
		}(),
	}

	mockService.On("GetUserCollections", mock.Anything, userID, "", 1, 20).
		Return(expectedCollections, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/collections?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_CheckCollected_True 测试检查已收藏-已收藏
func TestCollectionAPI_CheckCollected_True(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	mockService.On("IsCollected", mock.Anything, userID, bookID).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/collections/check/"+bookID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, true, data["is_collected"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_UpdateCollection_Success 测试更新收藏成功
func TestCollectionAPI_UpdateCollection_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	collectionID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	newNote := "更新后的笔记"
	mockService.On("UpdateCollection", mock.Anything, userID, collectionID, mock.Anything).Return(nil)

	reqBody := map[string]interface{}{
		"note": &newNote,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/collections/"+collectionID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_DeleteCollection_Success 测试删除收藏成功
func TestCollectionAPI_DeleteCollection_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	collectionID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	mockService.On("RemoveFromCollection", mock.Anything, userID, collectionID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/collections/"+collectionID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_CreateFolder_Success 测试创建收藏夹成功
func TestCollectionAPI_CreateFolder_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	expectedFolder := &socialModels.CollectionFolder{
		UserID:      userID,
		Name:        "我的收藏夹",
		Description: "收藏我喜欢的书",
	}
	expectedFolder.ID = primitive.NewObjectID()
	expectedFolder.CreatedAt = time.Now()

	mockService.On("CreateFolder", mock.Anything, userID, "我的收藏夹", "收藏我喜欢的书", true).
		Return(expectedFolder, nil)

	reqBody := map[string]interface{}{
		"name":        "我的收藏夹",
		"description": "收藏我喜欢的书",
		"is_public":   true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/folders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_CreateFolder_MissingName 测试创建收藏夹缺少名称
func TestCollectionAPI_CreateFolder_MissingName(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"description": "收藏夹描述",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/folders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 参数错误code为1001
}

// TestCollectionAPI_GetFolders_Success 测试获取收藏夹列表成功
func TestCollectionAPI_GetFolders_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	expectedFolders := []*socialModels.CollectionFolder{
		func() *socialModels.CollectionFolder {
			f := &socialModels.CollectionFolder{Name: "收藏夹1"}
			f.ID = primitive.NewObjectID()
			return f
		}(),
		func() *socialModels.CollectionFolder {
			f := &socialModels.CollectionFolder{Name: "收藏夹2"}
			f.ID = primitive.NewObjectID()
			return f
		}(),
	}

	mockService.On("GetUserFolders", mock.Anything, userID).Return(expectedFolders, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/folders", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_UpdateFolder_Success 测试更新收藏夹成功
func TestCollectionAPI_UpdateFolder_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	folderID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	newName := "更新后的收藏夹名称"
	mockService.On("UpdateFolder", mock.Anything, userID, folderID, mock.Anything).Return(nil)

	reqBody := map[string]interface{}{
		"name": &newName,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/folders/"+folderID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_DeleteFolder_Success 测试删除收藏夹成功
func TestCollectionAPI_DeleteFolder_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	folderID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	mockService.On("DeleteFolder", mock.Anything, userID, folderID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/folders/"+folderID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_ShareCollection_Success 测试分享收藏成功
func TestCollectionAPI_ShareCollection_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	collectionID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	mockService.On("ShareCollection", mock.Anything, userID, collectionID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reader/collections/"+collectionID+"/share", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_UnshareCollection_Success 测试取消分享收藏成功
func TestCollectionAPI_UnshareCollection_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	collectionID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	mockService.On("UnshareCollection", mock.Anything, userID, collectionID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/collections/"+collectionID+"/share", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_GetPublicCollections_Success 测试获取公开收藏列表成功
func TestCollectionAPI_GetPublicCollections_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	router := setupCollectionTestRouter(mockService, "")

	expectedCollections := []*socialModels.Collection{
		func() *socialModels.Collection {
			c := &socialModels.Collection{BookID: primitive.NewObjectID().Hex(), Note: "公开收藏1"}
			c.ID = primitive.NewObjectID()
			return c
		}(),
		func() *socialModels.Collection {
			c := &socialModels.Collection{BookID: primitive.NewObjectID().Hex(), Note: "公开收藏2"}
			c.ID = primitive.NewObjectID()
			return c
		}(),
	}

	mockService.On("GetPublicCollections", mock.Anything, 1, 20).
		Return(expectedCollections, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/public/collections?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_GetCollectionStats_Success 测试获取收藏统计成功
func TestCollectionAPI_GetCollectionStats_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	router := setupCollectionTestRouter(mockService, userID)

	expectedStats := map[string]interface{}{
		"total_collections": 50,
		"total_folders":     5,
		"public_count":      10,
	}

	mockService.On("GetUserCollectionStats", mock.Anything, userID).Return(expectedStats, nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/collections/stats", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestCollectionAPI_GetCollectionsByTag_Success 测试按标签获取收藏成功
func TestCollectionAPI_GetCollectionsByTag_Success(t *testing.T) {
	// Given
	mockService := new(MockCollectionService)
	userID := primitive.NewObjectID().Hex()
	tag := "推荐"
	router := setupCollectionTestRouter(mockService, userID)

	expectedCollections := []*socialModels.Collection{
		func() *socialModels.Collection {
			c := &socialModels.Collection{BookID: primitive.NewObjectID().Hex(), Note: "推荐书籍1", Tags: []string{"推荐"}}
			c.ID = primitive.NewObjectID()
			return c
		}(),
	}

	mockService.On("GetCollectionsByTag", mock.Anything, userID, tag, 1, 20).
		Return(expectedCollections, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/collections/tag/"+tag+"?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 成功响应code为0

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}
