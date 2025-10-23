package api

import (
	documentModel "Qingyu_backend/models/writer"
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

func setupDocumentTestRouter(documentService *document.DocumentService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userID到request context
	r.Use(func(c *gin.Context) {
		if userID != "" {
			ctx := context.WithValue(c.Request.Context(), "userID", userID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	})

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
		Type:      "chapter",
		ParentID:  "",
		Order:     1,
		WordCount: 1000,
		Status:    "writing",
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
		setupMock      func(*MockDocumentRepository, *MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "成功创建文档",
			projectID: "project123",
			userID:    "user123",
			requestBody: document.CreateDocumentRequest{
				Title: "新文档",
				Type:  "chapter",
			},
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				testProject := createTestProject("user123")
				projRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
				docRepo.On("Create", mock.Anything, mock.AnythingOfType("*document.Document")).Return(nil)
				// updateProjectStatistics会异步调用GetByProjectID来统计（使用Maybe因为是异步的）
				docRepo.On("GetByProjectID", mock.Anything, "project123", int64(10000), int64(0)).Return([]*documentModel.Document{}, nil).Maybe()
				projRepo.On("Update", mock.Anything, "project123", mock.AnythingOfType("map[string]interface {}")).Return(nil).Maybe()
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(201), resp["code"])
				assert.Equal(t, "创建成功", resp["message"])
				if data, ok := resp["data"].(map[string]interface{}); ok {
					assert.NotEmpty(t, data["documentId"])
				}
			},
		},
		{
			name:           "缺少必填字段",
			projectID:      "project123",
			userID:         "user123",
			requestBody:    document.CreateDocumentRequest{},
			setupMock:      func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
		{
			name:      "项目不存在",
			projectID: "nonexistent",
			userID:    "user123",
			requestBody: document.CreateDocumentRequest{
				Title: "新文档",
				Type:  "chapter",
			},
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				projRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
				assert.Contains(t, resp["error"], "项目不存在")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDocRepo := new(MockDocumentRepository)
			mockProjRepo := new(MockProjectRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockDocRepo, mockProjRepo)

			documentService := document.NewDocumentService(mockDocRepo, mockProjRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService, tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/projects/"+tt.projectID+"/documents", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockDocRepo.AssertExpectations(t)
			mockProjRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_GetDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		userID         string
		setupMock      func(*MockDocumentRepository, *MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功获取文档",
			documentID: "doc123",
			userID:     "user123",
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				testDoc := createTestDocument("project123")
				testDoc.ID = "doc123"
				testProject := createTestProject("user123")
				docRepo.On("GetByID", mock.Anything, "doc123").Return(testDoc, nil)
				projRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
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
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				docRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDocRepo := new(MockDocumentRepository)
			mockProjRepo := new(MockProjectRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockDocRepo, mockProjRepo)

			documentService := document.NewDocumentService(mockDocRepo, mockProjRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/documents/"+tt.documentID, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockDocRepo.AssertExpectations(t)
			mockProjRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_UpdateDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		userID         string
		requestBody    document.UpdateDocumentRequest
		setupMock      func(*MockDocumentRepository, *MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功更新文档",
			documentID: "doc123",
			userID:     "user123",
			requestBody: document.UpdateDocumentRequest{
				Title:  "更新后的标题",
				Status: "completed",
			},
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				testDoc := createTestDocument("project123")
				testDoc.ID = "doc123"
				testProject := createTestProject("user123")
				docRepo.On("GetByID", mock.Anything, "doc123").Return(testDoc, nil)
				projRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
				docRepo.On("UpdateByProject", mock.Anything, "doc123", "project123", mock.AnythingOfType("map[string]interface {}")).Return(nil)
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
			mockDocRepo := new(MockDocumentRepository)
			mockProjRepo := new(MockProjectRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockDocRepo, mockProjRepo)

			documentService := document.NewDocumentService(mockDocRepo, mockProjRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService, tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/v1/documents/"+tt.documentID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockDocRepo.AssertExpectations(t)
			mockProjRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_DeleteDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		userID         string
		setupMock      func(*MockDocumentRepository, *MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功删除文档",
			documentID: "doc123",
			userID:     "user123",
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				testDoc := createTestDocument("project123")
				testDoc.ID = "doc123"
				testProject := createTestProject("user123")
				docRepo.On("GetByID", mock.Anything, "doc123").Return(testDoc, nil)
				projRepo.On("GetByID", mock.Anything, "project123").Return(testProject, nil)
				docRepo.On("SoftDelete", mock.Anything, "doc123", "project123").Return(nil)
				// DeleteDocument也会异步调用updateProjectStatistics（使用Maybe因为是异步的）
				docRepo.On("GetByProjectID", mock.Anything, "project123", int64(10000), int64(0)).Return([]*documentModel.Document{}, nil).Maybe()
				projRepo.On("Update", mock.Anything, "project123", mock.AnythingOfType("map[string]interface {}")).Return(nil).Maybe()
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
			mockDocRepo := new(MockDocumentRepository)
			mockProjRepo := new(MockProjectRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockDocRepo, mockProjRepo)

			documentService := document.NewDocumentService(mockDocRepo, mockProjRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService, tt.userID)

			req := httptest.NewRequest("DELETE", "/api/v1/documents/"+tt.documentID, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockDocRepo.AssertExpectations(t)
			mockProjRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentApi_ListDocuments(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		userID         string
		queryParams    string
		setupMock      func(*MockDocumentRepository, *MockProjectRepository)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "成功获取文档列表",
			projectID:   "project123",
			userID:      "user123",
			queryParams: "?page=1&pageSize=20",
			setupMock: func(docRepo *MockDocumentRepository, projRepo *MockProjectRepository) {
				docs := []*documentModel.Document{
					createTestDocument("project123"),
					createTestDocument("project123"),
				}
				// ListDocuments方法会调用两次GetByProjectID
				docRepo.On("GetByProjectID", mock.Anything, "project123", int64(20), int64(0)).Return(docs, nil)
				docRepo.On("GetByProjectID", mock.Anything, "project123", int64(10000), int64(0)).Return(docs, nil)
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
			mockDocRepo := new(MockDocumentRepository)
			mockProjRepo := new(MockProjectRepository)
			mockEventBus := new(MockDocumentEventBus)
			tt.setupMock(mockDocRepo, mockProjRepo)

			documentService := document.NewDocumentService(mockDocRepo, mockProjRepo, mockEventBus)
			router := setupDocumentTestRouter(documentService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.projectID+"/documents"+tt.queryParams, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockDocRepo.AssertExpectations(t)
		})
	}
}
