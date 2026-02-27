package writer_test

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
	"go.mongodb.org/mongo-driver/bson/primitive"

	writerAPI "Qingyu_backend/api/v1/writer"
	"Qingyu_backend/models/writer"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// MockOutlineService Mock大纲服务
type MockOutlineService struct {
	mock.Mock
}

func (m *MockOutlineService) Create(ctx context.Context, projectID, userID string, req *serviceInterfaces.CreateOutlineRequest) (*writer.OutlineNode, error) {
	args := m.Called(ctx, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.OutlineNode), args.Error(1)
}

func (m *MockOutlineService) GetByID(ctx context.Context, outlineID, projectID string) (*writer.OutlineNode, error) {
	args := m.Called(ctx, outlineID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.OutlineNode), args.Error(1)
}

func (m *MockOutlineService) List(ctx context.Context, projectID string) ([]*writer.OutlineNode, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.OutlineNode), args.Error(1)
}

func (m *MockOutlineService) Update(ctx context.Context, outlineID, projectID string, req *serviceInterfaces.UpdateOutlineRequest) (*writer.OutlineNode, error) {
	args := m.Called(ctx, outlineID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.OutlineNode), args.Error(1)
}

func (m *MockOutlineService) Delete(ctx context.Context, outlineID, projectID string) error {
	args := m.Called(ctx, outlineID, projectID)
	return args.Error(0)
}

func (m *MockOutlineService) GetTree(ctx context.Context, projectID string) ([]*serviceInterfaces.OutlineTreeNode, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*serviceInterfaces.OutlineTreeNode), args.Error(1)
}

func (m *MockOutlineService) GetChildren(ctx context.Context, projectID, parentID string) ([]*writer.OutlineNode, error) {
	args := m.Called(ctx, projectID, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.OutlineNode), args.Error(1)
}

// setupOutlineTestRouter 设置测试路由
func setupOutlineTestRouter(outlineService *MockOutlineService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	// 添加错误处理中间件
	r.Use(func(c *gin.Context) {
		c.Next()
		// 检查是否有错误写入
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(500, gin.H{
				"code":    5000,
				"message": "内部服务器错误",
				"details": err.Error(),
			})
		}
	})

	api := writerAPI.NewOutlineApi(outlineService)
	r.POST("/api/v1/writer/projects/:projectId/outlines", api.CreateOutline)
	r.GET("/api/v1/writer/projects/:projectId/outlines", api.ListOutlines)
	r.GET("/api/v1/writer/projects/:projectId/outlines/tree", api.GetOutlineTree)
	r.GET("/api/v1/writer/projects/:projectId/outlines/children", api.GetOutlineChildren)
	r.GET("/api/v1/writer/outlines/:outlineId", api.GetOutline)
	r.PUT("/api/v1/writer/outlines/:outlineId", api.UpdateOutline)
	r.DELETE("/api/v1/writer/outlines/:outlineId", api.DeleteOutline)

	return r
}

// TestOutlineApi_CreateOutline_Success 测试成功创建大纲
func TestOutlineApi_CreateOutline_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"title":    "第一卷",
		"summary":  "这是第一卷的内容",
		"type":     "volume",
		"tension":  7,
		"order":    0,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/outlines", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	expectedOutline := &writer.OutlineNode{}
	expectedOutline.ID = primitive.NewObjectID()
	expectedOutline.ProjectID = projectID
	expectedOutline.Title = "第一卷"
	expectedOutline.Summary = "这是第一卷的内容"
	expectedOutline.Type = "volume"
	expectedOutline.Tension = 7
	expectedOutline.Order = 0
	expectedOutline.ParentID = ""

	mockService.On("Create", mock.Anything, projectID, userID, mock.Anything).Return(expectedOutline, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 0 = Success
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "第一卷", data["title"])
	assert.Equal(t, "volume", data["type"])

	mockService.AssertExpectations(t)
}

// TestOutlineApi_CreateOutline_MissingProjectID 测试缺少项目ID
func TestOutlineApi_CreateOutline_MissingProjectID(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	router := setupOutlineTestRouter(mockService, "")

	reqBody := map[string]interface{}{
		"title": "第一卷",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects//outlines", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 1001 = InvalidParams
}

// TestOutlineApi_CreateOutline_InvalidJSON 测试无效的JSON
func TestOutlineApi_CreateOutline_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/outlines", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestOutlineApi_CreateOutline_ServiceError 测试服务错误
func TestOutlineApi_CreateOutline_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"title": "第一卷",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/outlines", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	mockService.On("Create", mock.Anything, projectID, userID, mock.Anything).Return(nil, errors.New("service error"))

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(5000), response["code"]) // 5000 = InternalError

	mockService.AssertExpectations(t)
}

// TestOutlineApi_ListOutlines_Success 测试成功获取大纲列表
func TestOutlineApi_ListOutlines_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/outlines", nil)

	expectedOutlines := []*writer.OutlineNode{
		func() *writer.OutlineNode {
			o := &writer.OutlineNode{}
			o.ID = primitive.NewObjectID()
			o.Title = "第一章"
			o.ProjectID = projectID
			return o
		}(),
		func() *writer.OutlineNode {
			o := &writer.OutlineNode{}
			o.ID = primitive.NewObjectID()
			o.Title = "第二章"
			o.ProjectID = projectID
			return o
		}(),
	}

	mockService.On("List", mock.Anything, projectID).Return(expectedOutlines, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	mockService.AssertExpectations(t)
}

// TestOutlineApi_GetOutlineTree_Success 测试成功获取大纲树
func TestOutlineApi_GetOutlineTree_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/outlines/tree", nil)

	volume := &writer.OutlineNode{}
	volume.ID = primitive.NewObjectID()
	volume.Title = "第一卷"
	volume.ProjectID = projectID

	chapter1 := &writer.OutlineNode{}
	chapter1.ID = primitive.NewObjectID()
	chapter1.Title = "第一章"
	chapter1.ProjectID = projectID

	expectedTree := []*serviceInterfaces.OutlineTreeNode{
		{
			OutlineNode: volume,
			Children: []*serviceInterfaces.OutlineTreeNode{
				{OutlineNode: chapter1},
			},
		},
	}

	mockService.On("GetTree", mock.Anything, projectID).Return(expectedTree, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	mockService.AssertExpectations(t)
}

// TestOutlineApi_GetOutlineChildren_Success 测试成功获取子节点
func TestOutlineApi_GetOutlineChildren_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/outlines/children?parentId="+parentID, nil)

	expectedChildren := []*writer.OutlineNode{
		func() *writer.OutlineNode {
			o := &writer.OutlineNode{}
			o.ID = primitive.NewObjectID()
			o.Title = "第一章"
			o.ProjectID = projectID
			o.ParentID = parentID
			return o
		}(),
	}

	mockService.On("GetChildren", mock.Anything, projectID, parentID).Return(expectedChildren, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	mockService.AssertExpectations(t)
}

// TestOutlineApi_GetOutlineChildren_RootNodes 测试获取根节点
func TestOutlineApi_GetOutlineChildren_RootNodes(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/outlines/children", nil)

	expectedRoots := []*writer.OutlineNode{
		func() *writer.OutlineNode {
			o := &writer.OutlineNode{}
			o.ID = primitive.NewObjectID()
			o.Title = "第一卷"
			o.ProjectID = projectID
			o.ParentID = ""
			return o
		}(),
	}

	mockService.On("GetChildren", mock.Anything, projectID, "").Return(expectedRoots, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	mockService.AssertExpectations(t)
}

// TestOutlineApi_GetOutline_Success 测试成功获取大纲详情
func TestOutlineApi_GetOutline_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/outlines/"+outlineID+"?projectId="+projectID, nil)

	objID, _ := primitive.ObjectIDFromHex(outlineID)
	expectedOutline := &writer.OutlineNode{}
	expectedOutline.ID = objID
	expectedOutline.Title = "第一章"
	expectedOutline.ProjectID = projectID
	expectedOutline.Summary = "这是第一章的内容"
	expectedOutline.Type = "chapter"
	expectedOutline.Tension = 5

	mockService.On("GetByID", mock.Anything, outlineID, projectID).Return(expectedOutline, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "第一章", data["title"])
	assert.Equal(t, "这是第一章的内容", data["summary"])

	mockService.AssertExpectations(t)
}

// TestOutlineApi_GetOutline_MissingProjectID 测试缺少项目ID
func TestOutlineApi_GetOutline_MissingProjectID(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/outlines/"+outlineID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 1001 = InvalidParams
}

// TestOutlineApi_GetOutline_NotFound 测试大纲不存在
func TestOutlineApi_GetOutline_NotFound(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/outlines/"+outlineID+"?projectId="+projectID, nil)

	mockService.On("GetByID", mock.Anything, outlineID, projectID).Return(nil, errors.New("not found"))

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"].(string), "大纲不存在")

	mockService.AssertExpectations(t)
}

// TestOutlineApi_UpdateOutline_Success 测试成功更新大纲
func TestOutlineApi_UpdateOutline_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	newTitle := "新标题"
	newSummary := "新摘要"

	reqBody := map[string]interface{}{
		"title":   newTitle,
		"summary": newSummary,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/writer/outlines/"+outlineID+"?projectId="+projectID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	objID, _ := primitive.ObjectIDFromHex(outlineID)
	expectedOutline := &writer.OutlineNode{}
	expectedOutline.ID = objID
	expectedOutline.Title = newTitle
	expectedOutline.ProjectID = projectID
	expectedOutline.Summary = newSummary

	mockService.On("Update", mock.Anything, outlineID, projectID, mock.Anything).Return(expectedOutline, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, newTitle, data["title"])
	assert.Equal(t, newSummary, data["summary"])

	mockService.AssertExpectations(t)
}

// TestOutlineApi_UpdateOutline_MissingParams 测试缺少必需参数
func TestOutlineApi_UpdateOutline_MissingParams(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	reqBody := map[string]interface{}{
		"title": "新标题",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/writer/outlines/"+outlineID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 1001 = InvalidParams
}

// TestOutlineApi_DeleteOutline_Success 测试成功删除大纲
func TestOutlineApi_DeleteOutline_Success(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("DELETE", "/api/v1/writer/outlines/"+outlineID+"?projectId="+projectID, nil)

	mockService.On("Delete", mock.Anything, outlineID, projectID).Return(nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	mockService.AssertExpectations(t)
}

// TestOutlineApi_DeleteOutline_MissingParams 测试缺少必需参数
func TestOutlineApi_DeleteOutline_MissingParams(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("DELETE", "/api/v1/writer/outlines/"+outlineID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 1001 = InvalidParams
}

// TestOutlineApi_DeleteOutline_ServiceError 测试删除服务错误
func TestOutlineApi_DeleteOutline_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockOutlineService)
	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	router := setupOutlineTestRouter(mockService, "")

	req, _ := http.NewRequest("DELETE", "/api/v1/writer/outlines/"+outlineID+"?projectId="+projectID, nil)

	mockService.On("Delete", mock.Anything, outlineID, projectID).Return(errors.New("service error"))

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(5000), response["code"]) // 5000 = InternalError

	mockService.AssertExpectations(t)
}
