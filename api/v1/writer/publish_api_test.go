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

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// MockPublishService Mock发布服务 - 完整实现所有接口方法
type MockPublishService struct {
	mock.Mock
}

func (m *MockPublishService) PublishProject(ctx context.Context, projectID, userID string, req *interfaces.PublishProjectRequest) (*interfaces.PublicationRecord, error) {
	args := m.Called(ctx, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.PublicationRecord), args.Error(1)
}

func (m *MockPublishService) UnpublishProject(ctx context.Context, projectID, userID string) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *MockPublishService) GetProjectPublicationStatus(ctx context.Context, projectID string) (*interfaces.PublicationStatus, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.PublicationStatus), args.Error(1)
}

func (m *MockPublishService) PublishDocument(ctx context.Context, documentID, projectID, userID string, req *interfaces.PublishDocumentRequest) (*interfaces.PublicationRecord, error) {
	args := m.Called(ctx, documentID, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.PublicationRecord), args.Error(1)
}

func (m *MockPublishService) UpdateDocumentPublishStatus(ctx context.Context, documentID, projectID, userID string, req *interfaces.UpdateDocumentPublishStatusRequest) error {
	args := m.Called(ctx, documentID, projectID, userID, req)
	return args.Error(0)
}

func (m *MockPublishService) BatchPublishDocuments(ctx context.Context, projectID, userID string, req *interfaces.BatchPublishDocumentsRequest) (*interfaces.BatchPublishResult, error) {
	args := m.Called(ctx, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.BatchPublishResult), args.Error(1)
}

func (m *MockPublishService) GetPublicationRecords(ctx context.Context, projectID string, page, pageSize int) ([]*interfaces.PublicationRecord, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	if args.Get(0) == nil {
		return nil, int64(0), args.Error(1)
	}
	return args.Get(0).([]*interfaces.PublicationRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockPublishService) GetPublicationRecord(ctx context.Context, recordID string) (*interfaces.PublicationRecord, error) {
	args := m.Called(ctx, recordID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.PublicationRecord), args.Error(1)
}

func (m *MockPublishService) GetPendingPublicationRecords(ctx context.Context, page, pageSize int) ([]*interfaces.PublicationRecord, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*interfaces.PublicationRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockPublishService) ReviewPublication(ctx context.Context, recordID, reviewerID string, approved bool, note string) (*interfaces.PublicationRecord, error) {
	args := m.Called(ctx, recordID, reviewerID, approved, note)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.PublicationRecord), args.Error(1)
}

// setupPublishTestRouter 设置测试路由
func setupPublishTestRouter(publishService interfaces.PublishService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置user_id（用于需要认证的端点）
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

	publishAPI := writer.NewPublishApi(publishService)
	r.POST("/api/v1/writer/projects/:id/publish", publishAPI.PublishProject)
	r.POST("/api/v1/writer/projects/:id/unpublish", publishAPI.UnpublishProject)
	r.GET("/api/v1/writer/projects/:id/publication-status", publishAPI.GetProjectPublicationStatus)
	r.POST("/api/v1/writer/documents/:id/publish", publishAPI.PublishDocument)
	r.PUT("/api/v1/writer/documents/:id/publish-status", publishAPI.UpdateDocumentPublishStatus)
	r.POST("/api/v1/writer/projects/:id/documents/batch-publish", publishAPI.BatchPublishDocuments)
	r.GET("/api/v1/writer/projects/:id/publications", publishAPI.GetPublicationRecords)
	r.GET("/api/v1/writer/publications/:id", publishAPI.GetPublicationRecord)

	return r
}

// TestPublishApi_PublishProject_Success 测试成功发布项目
func TestPublishApi_PublishProject_Success(t *testing.T) {
	// Given
	mockService := new(MockPublishService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"bookstoreId": "test_bookstore",
		"autoPublish": true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	expectedRecord := &interfaces.PublicationRecord{
		ID:          primitive.NewObjectID().Hex(),
		Type:        "project",
		ResourceID:  projectID,
		BookstoreID: "test_bookstore",
		Status:      "published",
	}

	mockService.On("PublishProject", mock.Anything, projectID, userID, mock.Anything).Return(expectedRecord, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 0 = Success
	assert.Equal(t, "操作成功", response["message"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestPublishApi_PublishProject_MissingProjectID 测试缺少项目ID
func TestPublishApi_PublishProject_MissingProjectID(t *testing.T) {
	// Given
	mockService := new(MockPublishService)
	router := setupPublishTestRouter(mockService, "")

	reqBody := map[string]interface{}{
		"bookstoreId": "test_bookstore",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects//publish", bytes.NewBuffer(jsonBody))
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

// TestPublishApi_PublishProject_ServiceError 测试服务错误
func TestPublishApi_PublishProject_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockPublishService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"bookstoreId": "test_bookstore",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	mockService.On("PublishProject", mock.Anything, projectID, userID, mock.Anything).Return(nil, errors.New("service error"))

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

// TestPublishApi_PublishProject_InvalidJSON 测试无效的JSON
func TestPublishApi_PublishProject_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockPublishService)
	router := setupPublishTestRouter(mockService, "")

	projectID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublishApi_GetProjectPublicationStatus_RejectsInvalidProjectID(t *testing.T) {
	mockService := new(MockPublishService)
	router := setupPublishTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/project!bad/publication-status", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID格式不正确")
	mockService.AssertNotCalled(t, "GetProjectPublicationStatus", mock.Anything, mock.Anything)
}

func TestPublishApi_PublishDocument_UsesOptionalProjectIDAndUserID(t *testing.T) {
	mockService := new(MockPublishService)
	documentID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"chapterTitle":  "第一章",
		"chapterNumber": 3,
		"isFree":        true,
		"authorNote":    "ready",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/documents/"+documentID+"/publish?projectId="+projectID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	expectedRecord := &interfaces.PublicationRecord{
		ID:         primitive.NewObjectID().Hex(),
		Type:       "document",
		ResourceID: documentID,
		Status:     "pending",
	}
	mockService.
		On("PublishDocument", mock.Anything, documentID, projectID, userID, mock.Anything).
		Return(expectedRecord, nil).
		Once()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), documentID)
	mockService.AssertExpectations(t)
}

func TestPublishApi_UpdateDocumentPublishStatus_RequiresProjectID(t *testing.T) {
	mockService := new(MockPublishService)
	documentID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, "")

	req, _ := http.NewRequest("PUT", "/api/v1/writer/documents/"+documentID+"/publish-status", bytes.NewBufferString(`{"isPublished":true}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "documentId和projectId不能为空")
	mockService.AssertNotCalled(t, "UpdateDocumentPublishStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPublishApi_BatchPublishDocuments_Success(t *testing.T) {
	mockService := new(MockPublishService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"documentIds":   []string{"doc-1", "doc-2"},
		"autoNumbering": true,
		"startNumber":   1,
		"isFree":        true,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/"+projectID+"/documents/batch-publish", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	expectedResult := &interfaces.BatchPublishResult{
		SuccessCount: 2,
		Results: []interfaces.BatchPublishItem{
			{DocumentID: "doc-1", Success: true, RecordID: "record-1"},
			{DocumentID: "doc-2", Success: true, RecordID: "record-2"},
		},
	}
	mockService.
		On("BatchPublishDocuments", mock.Anything, projectID, userID, mock.Anything).
		Return(expectedResult, nil).
		Once()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"successCount":2`)
	mockService.AssertExpectations(t)
}

func TestPublishApi_GetPublicationRecords_UsesSizeAliasPagination(t *testing.T) {
	mockService := new(MockPublishService)
	projectID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, "")

	records := []*interfaces.PublicationRecord{
		{ID: "record-1", Type: "project", ResourceID: projectID, Status: "published"},
	}
	mockService.
		On("GetPublicationRecords", mock.Anything, projectID, 2, 5).
		Return(records, int64(7), nil).
		Once()

	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/publications?page=2&size=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"total":7`)
	mockService.AssertExpectations(t)
}

func TestPublishApi_GetPublicationRecord_NotFound(t *testing.T) {
	mockService := new(MockPublishService)
	recordID := primitive.NewObjectID().Hex()
	router := setupPublishTestRouter(mockService, "")

	mockService.On("GetPublicationRecord", mock.Anything, recordID).Return(nil, errors.New("not found")).Once()

	req, _ := http.NewRequest("GET", "/api/v1/writer/publications/"+recordID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "发布记录不存在")
	mockService.AssertExpectations(t)
}
