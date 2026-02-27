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

	publishAPI := writer.NewPublishApi(publishService)
	r.POST("/api/v1/writer/projects/:id/publish", publishAPI.PublishProject)

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
