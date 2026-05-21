package writer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	"Qingyu_backend/repository/interfaces/infrastructure"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/writer/document"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockProjectRepository struct {
	mock.Mock
}

func (m *mockProjectRepository) Create(ctx context.Context, project *writer.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *mockProjectRepository) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

func (m *mockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *mockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Project, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *mockProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *mockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *mockProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	args := m.Called(ctx, projectID, ownerID, updates)
	return args.Error(0)
}

func (m *mockProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	args := m.Called(ctx, projectID, ownerID)
	return args.Bool(0), args.Error(1)
}

func (m *mockProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *mockProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *mockProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *mockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockProjectRepository) CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

type documentRepoStub struct {
	getByProjectIDErr error
}

func (s *documentRepoStub) Create(ctx context.Context, document *writer.Document) error { return nil }
func (s *documentRepoStub) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	return nil, nil
}
func (s *documentRepoStub) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}
func (s *documentRepoStub) Delete(ctx context.Context, id string) error { return nil }
func (s *documentRepoStub) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Document, error) {
	return nil, nil
}
func (s *documentRepoStub) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}
func (s *documentRepoStub) Exists(ctx context.Context, id string) (bool, error) { return false, nil }
func (s *documentRepoStub) Health(ctx context.Context) error                   { return nil }
func (s *documentRepoStub) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Document, error) {
	if s.getByProjectIDErr != nil {
		return nil, s.getByProjectIDErr
	}
	return nil, nil
}
func (s *documentRepoStub) GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*writer.Document, error) {
	return nil, nil
}
func (s *documentRepoStub) GetByIDs(ctx context.Context, ids []string) ([]*writer.Document, error) {
	return nil, nil
}
func (s *documentRepoStub) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	return nil
}
func (s *documentRepoStub) DeleteByProject(ctx context.Context, documentID, projectID string) error {
	return nil
}
func (s *documentRepoStub) RestoreByProject(ctx context.Context, documentID, projectID string) error {
	return nil
}
func (s *documentRepoStub) IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error) {
	return true, nil
}
func (s *documentRepoStub) SoftDelete(ctx context.Context, documentID, projectID string) error {
	return nil
}
func (s *documentRepoStub) HardDelete(ctx context.Context, documentID string) error { return nil }
func (s *documentRepoStub) GetByIDUnscoped(ctx context.Context, id string) (*writer.Document, error) {
	return nil, nil
}
func (s *documentRepoStub) CountByProject(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}
func (s *documentRepoStub) CreateWithTransaction(ctx context.Context, document *writer.Document, callback func(ctx context.Context) error) error {
	return nil
}

var _ writerRepo.DocumentRepository = (*documentRepoStub)(nil)

func setupDocumentCreateTestRouter(api *DocumentApi, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "内部服务器错误",
				"details": err.Error(),
			})
		}
	})

	router.POST("/api/v1/writer/projects/:projectId/documents", api.CreateDocument)
	return router
}

func TestDocumentApiCreateDocument_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()
	projectID := "not-a-valid-object-id"
	authorID, err := primitive.ObjectIDFromHex(userID)
	require.NoError(t, err)

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{
			AuthorID: authorID,
		},
		Visibility: writer.VisibilityPrivate,
	}, nil)

	api := NewDocumentApi(document.NewDocumentService(nil, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID)

	reqBody, err := json.Marshal(map[string]any{
		"title": "新文档",
		"type":  "chapter",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(5000), resp["code"])
	assert.Contains(t, resp["details"].(string), "无效的项目ID")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiCreateDocument_ServiceNotInitialized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	projectID := primitive.NewObjectID().Hex()
	api := &DocumentApi{}
	router := setupDocumentCreateTestRouter(api, "")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "服务未初始化")
}

func TestDocumentApiDuplicateDocument_RequiresLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	api := &DocumentApi{}
	router.POST("/api/v1/writer/documents/:id/duplicate", api.DuplicateDocument)

	documentID := primitive.NewObjectID().Hex()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString(`{"position":"inner"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(1002), resp["code"])
}

func TestDocumentApiDuplicateDocument_RejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID().Hex())
		c.Next()
	})

	api := &DocumentApi{}
	router.POST("/api/v1/writer/documents/:id/duplicate", api.DuplicateDocument)

	documentID := primitive.NewObjectID().Hex()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString("{"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDocumentApiMoveDocument_RejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	api := &DocumentApi{}
	router.PUT("/api/v1/writer/documents/:id/move", api.MoveDocument)

	documentID := primitive.NewObjectID().Hex()
	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/documents/"+documentID+"/move", bytes.NewBufferString("{"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDocumentApiCreateDocumentByBody_ServiceNotInitialized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	api := &DocumentApi{}
	router.POST("/api/v1/writer/documents", api.CreateDocumentByBody)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents", bytes.NewBufferString(`{"project_id":"`+primitive.NewObjectID().Hex()+`","title":"新文档"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "服务未初始化")
}

func TestDocumentApiCreateDocumentByBody_RequiresLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	api := &DocumentApi{documentService: document.NewDocumentService(nil, nil, nil, nil)}
	router.POST("/api/v1/writer/documents", api.CreateDocumentByBody)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents", bytes.NewBufferString(`{"project_id":"`+primitive.NewObjectID().Hex()+`","title":"新文档"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(1002), resp["code"])
}

func TestDocumentApiCreateDocumentByBody_RejectsMalformedOrMissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID().Hex())
		c.Next()
	})

	api := &DocumentApi{documentService: document.NewDocumentService(nil, nil, nil, nil)}
	router.POST("/api/v1/writer/documents", api.CreateDocumentByBody)

	cases := []struct {
		name string
		body string
	}{
		{name: "malformed json", body: "{"},
		{name: "missing required fields", body: `{"project_id":"` + primitive.NewObjectID().Hex() + `"}`},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents", bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "参数错误")
		})
	}
}

func TestDocumentApiCreateDocument_ProjectLookupFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID).Return(nil, assert.AnError).Once()

	api := NewDocumentApi(document.NewDocumentService(nil, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID)

	reqBody, err := json.Marshal(map[string]any{
		"title": "新文档",
		"type":  "chapter",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(5000), resp["code"])
	assert.Contains(t, resp["details"].(string), "查询项目失败")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiCreateDocument_ProjectNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID).Return(nil, nil).Once()

	api := NewDocumentApi(document.NewDocumentService(nil, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID)

	reqBody, err := json.Marshal(map[string]any{
		"title": "新文档",
		"type":  "chapter",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(5000), resp["code"])
	assert.Contains(t, resp["details"].(string), "项目不存在")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiUpdateDocument_RejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: &document.DocumentService{}}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/writer/documents/document-1", strings.NewReader("{"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "document-1"}}

	api.UpdateDocument(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "参数错误")
}

func TestDocumentApiDeleteDocument_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: document.NewDocumentService(nil, nil, nil, nil)}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "内部服务器错误",
				"details": err.Error(),
			})
		}
	})
	router.DELETE("/api/v1/writer/documents/:id", api.DeleteDocument)

	req, err := http.NewRequest(http.MethodDelete, "/api/v1/writer/documents/document-1", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
}

func TestDocumentApiGetDocument_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: document.NewDocumentService(nil, nil, nil, nil)}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "内部服务器错误",
				"details": err.Error(),
			})
		}
	})
	router.GET("/api/v1/writer/documents/:id", api.GetDocument)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/writer/documents/document-1", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "用户未登录")
}

func TestDocumentApiGetDocumentTree_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: document.NewDocumentService(nil, nil, nil, nil)}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "内部服务器错误",
				"details": err.Error(),
			})
		}
	})
	router.GET("/api/v1/writer/projects/:projectId/documents/tree", api.GetDocumentTree)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/writer/projects/project-1/documents/tree", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "用户未登录")
}

func TestDocumentApiListDocuments_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewDocumentApi(document.NewDocumentService(&documentRepoStub{getByProjectIDErr: assert.AnError}, nil, nil, nil))
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "内部服务器错误",
				"details": err.Error(),
			})
		}
	})
	router.GET("/api/v1/writer/projects/:projectId/documents", api.ListDocuments)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/writer/projects/project-1/documents", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "查询文档列表失败")
}

func TestDocumentApiReorderDocuments_RejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{}
	router := gin.New()
	router.PUT("/api/v1/writer/projects/:projectId/documents/reorder", api.ReorderDocuments)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/projects/project-1/documents/reorder", bytes.NewBufferString("{"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDocumentApiCreateDocument_ProjectOwnerMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()
	ownerID := primitive.NewObjectID()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{
			AuthorID: ownerID,
		},
		Visibility: writer.VisibilityPrivate,
	}, nil).Once()

	api := NewDocumentApi(document.NewDocumentService(nil, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID)

	reqBody, err := json.Marshal(map[string]any{
		"title": "新文档",
		"type":  "chapter",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(5000), resp["code"])
	assert.Contains(t, resp["details"].(string), "无权限编辑该项目")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiCreateDocument_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	mockProjectRepo := new(mockProjectRepository)
	api := NewDocumentApi(document.NewDocumentService(nil, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID)

	reqBody, err := json.Marshal(map[string]any{
		"title": "新文档",
		"type":  "invalid_type",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(5000), resp["code"])
	assert.Contains(t, resp["details"].(string), "参数验证失败")
	assert.Contains(t, resp["details"].(string), "无效的文档类型")
}
