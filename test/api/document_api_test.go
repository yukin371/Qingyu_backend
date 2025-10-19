package api

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
	"github.com/stretchr/testify/require"

	writerAPI "Qingyu_backend/api/v1/writer"
	documentModel "Qingyu_backend/models/document"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/base"
	"Qingyu_backend/service/document"
)

// === Mock Repositories ===

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *documentModel.Document) error {
	args := m.Called(ctx, doc)
	if args.Error(0) == nil {
		if doc.ID == "" {
			doc.ID = "mock_document_id"
		}
		doc.CreatedAt = time.Now()
		doc.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*documentModel.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*documentModel.Document), args.Error(1)
}

func (m *MockDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*documentModel.Document, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*documentModel.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*documentModel.Document, error) {
	args := m.Called(ctx, projectID, documentType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*documentModel.Document), args.Error(1)
}

func (m *MockDocumentRepository) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	args := m.Called(ctx, documentID, projectID, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) DeleteByProject(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) RestoreByProject(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error) {
	args := m.Called(ctx, documentID, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentRepository) SoftDelete(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) HardDelete(ctx context.Context, documentID string) error {
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

func (m *MockDocumentRepository) CountByProject(ctx context.Context, projectID string) (int64, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) CreateWithTransaction(ctx context.Context, doc *documentModel.Document, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, doc, callback)
	return args.Error(0)
}

func (m *MockDocumentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*documentModel.Document, error) {
	return nil, nil
}

func (m *MockDocumentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}

func (m *MockDocumentRepository) Exists(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (m *MockDocumentRepository) Health(ctx context.Context) error {
	return nil
}

// === Mock EventBus ===

type MockDocumentEventBus struct {
	mock.Mock
}

func (m *MockDocumentEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return nil
}

func (m *MockDocumentEventBus) Unsubscribe(eventType string, handlerName string) error {
	return nil
}

func (m *MockDocumentEventBus) Publish(ctx context.Context, event base.Event) error {
	return nil
}

func (m *MockDocumentEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	return nil
}

// === 测试辅助函数 ===

func setupDocumentTestRouter(documentService *document.DocumentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := writerAPI.NewDocumentApi(documentService)

	v1 := r.Group("/api/v1")
	{
		// 文档路由
		documents := v1.Group("/documents")
		{
			documents.GET("/:id", api.GetDocument)
			documents.PUT("/:id", api.UpdateDocument)
			documents.DELETE("/:id", api.DeleteDocument)
			documents.PUT("/:id/move", api.MoveDocument)
		}

		// 项目下的文档路由
		projects := v1.Group("/projects/:projectId")
		{
			projects.POST("/documents", api.CreateDocument)
			projects.GET("/documents", api.ListDocuments)
			projects.GET("/documents/tree", api.GetDocumentTree)
			projects.PUT("/documents/reorder", api.ReorderDocuments)
		}
	}

	return r
}

func createTestDocument(projectID string) *documentModel.Document {
	now := time.Now()
	return &documentModel.Document{
		ID:        "test_document_id",
		ProjectID: projectID,
		Title:     "测试文档",
		Content:   "这是测试内容",
		Type:      "chapter",
		ParentID:  "",
		Order:     1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// === 测试用例 ===

func TestDocumentApi_CreateDocument(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		requestBody    document.CreateDocumentRequest
		setupMock      func(*MockDocumentRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "成功创建文档",
			projectID: "project123",
			userID:    "user123",
			requestBody: document.CreateDocumentRequest{
				Title:   "新文档",
				Content: "文档内容",
				Type:    "chapter",
			},
			setupMock: func(repo *MockDocumentRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*document.Document")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(201), resp["code"])
				assert.Equal(t, "创建成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["documentId"])
			},
		},
		{
			name:      "缺少必填字段",
			projectID: "project123",
			userID:    "user123",
			requestBody: document.CreateDocumentRequest{
				Content: "文档内容",
			},
			setupMock:      func(repo *MockDocumentRepository) {},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockRepo)

			documentService := document.NewDocumentService(mockRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/projects/"+tt.projectID+"/documents", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_GetDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		userID         string
		setupMock      func(*MockDocumentRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功获取文档",
			documentID: "doc123",
			userID:     "user123",
			setupMock: func(repo *MockDocumentRepository) {
				testDoc := createTestDocument("project123")
				testDoc.ID = "doc123"
				repo.On("GetByID", mock.Anything, "doc123").Return(testDoc, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "doc123", data["id"])
			},
		},
		{
			name:       "文档不存在",
			documentID: "nonexistent",
			userID:     "user123",
			setupMock: func(repo *MockDocumentRepository) {
				repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockRepo)

			documentService := document.NewDocumentService(mockRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService)

			req := httptest.NewRequest("GET", "/api/v1/documents/"+tt.documentID, nil)

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_UpdateDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		userID         string
		requestBody    document.UpdateDocumentRequest
		setupMock      func(*MockDocumentRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功更新文档",
			documentID: "doc123",
			userID:     "user123",
			requestBody: document.UpdateDocumentRequest{
				Title:   "更新后的标题",
				Content: "更新后的内容",
			},
			setupMock: func(repo *MockDocumentRepository) {
				testDoc := createTestDocument("project123")
				testDoc.ID = "doc123"
				repo.On("GetByID", mock.Anything, "doc123").Return(testDoc, nil)
				repo.On("Update", mock.Anything, "doc123", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "更新成功", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockRepo)

			documentService := document.NewDocumentService(mockRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/v1/documents/"+tt.documentID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_DeleteDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		userID         string
		setupMock      func(*MockDocumentRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功删除文档",
			documentID: "doc123",
			userID:     "user123",
			setupMock: func(repo *MockDocumentRepository) {
				testDoc := createTestDocument("project123")
				testDoc.ID = "doc123"
				repo.On("GetByID", mock.Anything, "doc123").Return(testDoc, nil)
				repo.On("SoftDelete", mock.Anything, "doc123", "project123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "删除成功", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockRepo)

			documentService := document.NewDocumentService(mockRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService)

			req := httptest.NewRequest("DELETE", "/api/v1/documents/"+tt.documentID, nil)

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_ListDocuments(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		queryParams    string
		setupMock      func(*MockDocumentRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "成功获取文档列表",
			projectID:   "project123",
			userID:      "user123",
			queryParams: "?page=1&pageSize=20",
			setupMock: func(repo *MockDocumentRepository) {
				docs := []*documentModel.Document{
					createTestDocument("project123"),
					createTestDocument("project123"),
				}
				repo.On("GetByProjectID", mock.Anything, "project123", int64(20), int64(0)).Return(docs, nil)
				repo.On("CountByProject", mock.Anything, "project123").Return(int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].(map[string]interface{})
				docs := data["documents"].([]interface{})
				assert.Len(t, docs, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDocumentRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockRepo)

			documentService := document.NewDocumentService(mockRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService)

			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.projectID+"/documents"+tt.queryParams, nil)

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "userID", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockRepo.AssertExpectations(t)
		})
	}
}
