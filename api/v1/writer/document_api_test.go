package writer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Qingyu_backend/models/dto"
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
	getByIDDoc        *writer.Document
	getByIDErr        error
	createErr         error
	createdDocs       []*writer.Document
	updateErr         error
	updateByProjectErr error
	lastUpdateByProjectDocumentID string
	lastUpdateByProjectProjectID  string
	lastUpdateByProjectUpdates    map[string]interface{}
	getByProjectIDCallCount int
	firstProjectID          string
	firstLimit              int64
	firstOffset             int64
}

type documentContentRepoStub struct {
	content        *writer.DocumentContent
	getErr         error
	createErr      error
	createdContents []*writer.DocumentContent
}

func (s *documentRepoStub) Create(ctx context.Context, document *writer.Document) error {
	s.createdDocs = append(s.createdDocs, document)
	return s.createErr
}
func (s *documentRepoStub) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	return s.getByIDDoc, nil
}
func (s *documentRepoStub) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.updateErr
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
	s.getByProjectIDCallCount++
	if s.getByProjectIDCallCount == 1 {
		s.firstProjectID = projectID
		s.firstLimit = limit
		s.firstOffset = offset
	}
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
	s.lastUpdateByProjectDocumentID = documentID
	s.lastUpdateByProjectProjectID = projectID
	s.lastUpdateByProjectUpdates = updates
	return s.updateByProjectErr
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

func (s *documentContentRepoStub) Create(ctx context.Context, content *writer.DocumentContent) error {
	s.createdContents = append(s.createdContents, content)
	return s.createErr
}
func (s *documentContentRepoStub) GetByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	return nil, nil
}
func (s *documentContentRepoStub) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}
func (s *documentContentRepoStub) Delete(ctx context.Context, id string) error { return nil }
func (s *documentContentRepoStub) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.DocumentContent, error) {
	return nil, nil
}
func (s *documentContentRepoStub) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}
func (s *documentContentRepoStub) Exists(ctx context.Context, id string) (bool, error) { return false, nil }
func (s *documentContentRepoStub) Health(ctx context.Context) error { return nil }
func (s *documentContentRepoStub) GetByDocumentID(ctx context.Context, documentID string) (*writer.DocumentContent, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.content, nil
}
func (s *documentContentRepoStub) UpdateWithVersion(ctx context.Context, documentID string, updates map[string]interface{}, expectedVersion int) error {
	return nil
}
func (s *documentContentRepoStub) BatchUpdateContent(ctx context.Context, updates map[string]string) error { return nil }
func (s *documentContentRepoStub) GetContentStats(ctx context.Context, documentID string) (wordCount, charCount int, err error) {
	return 0, 0, nil
}
func (s *documentContentRepoStub) StoreToGridFS(ctx context.Context, documentID string, content []byte) (gridFSID string, err error) {
	return "", nil
}
func (s *documentContentRepoStub) LoadFromGridFS(ctx context.Context, gridFSID string) (content []byte, err error) {
	return nil, nil
}
func (s *documentContentRepoStub) CreateWithTransaction(ctx context.Context, content *writer.DocumentContent, callback func(ctx context.Context) error) error {
	return nil
}

var _ writerRepo.DocumentContentRepository = (*documentContentRepoStub)(nil)

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

func TestDocumentApiMoveDocument_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()

	api := NewDocumentApi(document.NewDocumentService(&documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Level:            2,
		},
		updateErr: assert.AnError,
	}, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.PUT("/api/v1/writer/documents/:id/move", api.MoveDocument)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/documents/"+documentID+"/move", bytes.NewBufferString(`{"orderKey":"chapter-2"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "移动文档失败")

	mockProjectRepo.AssertExpectations(t)
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

func TestDocumentApiUpdateDocument_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()

	api := NewDocumentApi(document.NewDocumentService(&documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Level:            2,
		},
		updateByProjectErr: assert.AnError,
	}, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.PUT("/api/v1/writer/documents/:id", api.UpdateDocument)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/documents/"+documentID, bytes.NewBufferString(`{"title":"新标题"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "更新文档失败")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiUpdateDocument_ForwardsSelectedFieldsToRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()
	title := "新标题"
	notes := "新备注"

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()

	repo := &documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Level:            2,
		},
	}
	api := NewDocumentApi(document.NewDocumentService(repo, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.PUT("/api/v1/writer/documents/:id", api.UpdateDocument)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/documents/"+documentID, bytes.NewBufferString(`{"title":"新标题","notes":"新备注"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, documentID, repo.lastUpdateByProjectDocumentID)
	assert.Equal(t, projectID.Hex(), repo.lastUpdateByProjectProjectID)
	require.Len(t, repo.lastUpdateByProjectUpdates, 2)
	assert.Equal(t, title, repo.lastUpdateByProjectUpdates["title"])
	assert.Equal(t, notes, repo.lastUpdateByProjectUpdates["notes"])

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiUpdateDocument_ForwardsCompositeFieldsToRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()

	repo := &documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Level:            2,
		},
	}
	api := NewDocumentApi(document.NewDocumentService(repo, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.PUT("/api/v1/writer/documents/:id", api.UpdateDocument)

	reqBody := `{"status":"completed","tags":["tag-a","tag-b"],"characterIds":["char-1","char-2"]}`
	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/documents/"+documentID, bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, documentID, repo.lastUpdateByProjectDocumentID)
	assert.Equal(t, projectID.Hex(), repo.lastUpdateByProjectProjectID)
	require.Len(t, repo.lastUpdateByProjectUpdates, 3)
	assert.Equal(t, "completed", string(repo.lastUpdateByProjectUpdates["status"].(dto.DocumentStatus)))

	tagsPtr, ok := repo.lastUpdateByProjectUpdates["tags"].(*[]string)
	require.True(t, ok)
	assert.Equal(t, []string{"tag-a", "tag-b"}, *tagsPtr)

	characterIDsPtr, ok := repo.lastUpdateByProjectUpdates["character_ids"].(*[]string)
	require.True(t, ok)
	assert.Equal(t, []string{"char-1", "char-2"}, *characterIDsPtr)

	mockProjectRepo.AssertExpectations(t)
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

func TestDocumentApiListDocuments_UsesDefaultPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &documentRepoStub{}
	api := NewDocumentApi(document.NewDocumentService(repo, nil, nil, nil))
	router := gin.New()
	router.GET("/api/v1/writer/projects/:projectId/documents", api.ListDocuments)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/writer/projects/project-1/documents", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Data struct {
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
			Total    int `json:"total"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Data.Page)
	assert.Equal(t, 20, resp.Data.PageSize)
	assert.Equal(t, 0, resp.Data.Total)
	assert.Equal(t, "project-1", repo.firstProjectID)
	assert.EqualValues(t, 20, repo.firstLimit)
	assert.EqualValues(t, 0, repo.firstOffset)
	assert.Equal(t, 2, repo.getByProjectIDCallCount)
}

func TestDocumentApiListDocuments_NormalizesPaginationAndForwardsToService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &documentRepoStub{}
	api := NewDocumentApi(document.NewDocumentService(repo, nil, nil, nil))
	router := gin.New()
	router.GET("/api/v1/writer/projects/:projectId/documents", api.ListDocuments)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/writer/projects/project-1/documents?page=0&size=999", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Data struct {
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Data.Page)
	assert.Equal(t, 100, resp.Data.PageSize)
	assert.Equal(t, "project-1", repo.firstProjectID)
	assert.EqualValues(t, 100, repo.firstLimit)
	assert.EqualValues(t, 0, repo.firstOffset)
	assert.Equal(t, 2, repo.getByProjectIDCallCount)
}

func TestDocumentApiListDocuments_ForwardsNonDefaultOffset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &documentRepoStub{}
	api := NewDocumentApi(document.NewDocumentService(repo, nil, nil, nil))
	router := gin.New()
	router.GET("/api/v1/writer/projects/:projectId/documents", api.ListDocuments)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/writer/projects/project-1/documents?page=3&size=15", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Data struct {
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 3, resp.Data.Page)
	assert.Equal(t, 15, resp.Data.PageSize)
	assert.Equal(t, "project-1", repo.firstProjectID)
	assert.EqualValues(t, 15, repo.firstLimit)
	assert.EqualValues(t, 30, repo.firstOffset)
	assert.Equal(t, 2, repo.getByProjectIDCallCount)
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

func TestDocumentApiReorderDocuments_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID().Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()

	api := NewDocumentApi(document.NewDocumentService(&documentRepoStub{updateErr: assert.AnError}, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.PUT("/api/v1/writer/projects/:projectId/documents/reorder", api.ReorderDocuments)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/writer/projects/"+projectID+"/documents/reorder", bytes.NewBufferString(`{"items":[{"documentId":"doc-1","orderKey":"001"}]}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "更新文档 doc-1 排序失败")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiDuplicateDocument_SurfacesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()

	api := NewDocumentApi(document.NewDocumentService(&documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Title:            "原文档",
			Type:             "chapter",
			StableRef:        "doc-1",
			OrderKey:         "aaaa",
		},
		createErr: assert.AnError,
	}, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.POST("/api/v1/writer/documents/:id/duplicate", api.DuplicateDocument)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString(`{"copyContent":false}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "内部服务器错误")
	assert.Contains(t, w.Body.String(), "创建新文档失败")

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiDuplicateDocument_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()
	mockProjectRepo.On("Update", mock.Anything, projectID.Hex(), mock.Anything).Return(nil).Maybe()

	repo := &documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Title:            "原文档",
			Type:             "chapter",
			StableRef:        "doc-1",
			OrderKey:         "aaaa",
		},
	}
	api := NewDocumentApi(document.NewDocumentService(repo, nil, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.POST("/api/v1/writer/documents/:id/duplicate", api.DuplicateDocument)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString(`{"copyContent":false}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Len(t, repo.createdDocs, 1)
	assert.Equal(t, "Copy - 原文档", repo.createdDocs[0].Title)
	assert.Equal(t, "doc-1-copy", repo.createdDocs[0].StableRef)

	var resp struct {
		Data struct {
			NewDocumentID string `json:"newDocumentId"`
			Title         string `json:"title"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Data.NewDocumentID)
	assert.Equal(t, "Copy - 原文档", resp.Data.Title)

	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentApiDuplicateDocument_CopyContentSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentObjectID := primitive.NewObjectID()
	documentID := documentObjectID.Hex()

	mockProjectRepo := new(mockProjectRepository)
	mockProjectRepo.On("GetByID", mock.Anything, projectID.Hex()).Return(&writer.Project{
		OwnedEntity: base.OwnedEntity{AuthorID: userID},
	}, nil).Once()
	mockProjectRepo.On("Update", mock.Anything, projectID.Hex(), mock.Anything).Return(nil).Maybe()

	docRepo := &documentRepoStub{
		getByIDDoc: &writer.Document{
			IdentifiedEntity: base.IdentifiedEntity{ID: documentObjectID},
			ProjectID:        projectID,
			Title:            "原文档",
			Type:             "chapter",
			StableRef:        "doc-1",
			OrderKey:         "aaaa",
		},
	}
	contentRepo := &documentContentRepoStub{
		content: &writer.DocumentContent{
			DocumentID:  documentObjectID,
			Content:     "正文内容",
			ContentType: "markdown",
			WordCount:   4,
			CharCount:   4,
			Version:     1,
		},
	}
	api := NewDocumentApi(document.NewDocumentService(docRepo, contentRepo, mockProjectRepo, nil))
	router := setupDocumentCreateTestRouter(api, userID.Hex())
	router.POST("/api/v1/writer/documents/:id/duplicate", api.DuplicateDocument)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString(`{"copyContent":true}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Len(t, docRepo.createdDocs, 1)
	require.Len(t, contentRepo.createdContents, 1)
	assert.Equal(t, "正文内容", contentRepo.createdContents[0].Content)
	assert.Equal(t, "markdown", contentRepo.createdContents[0].ContentType)
	assert.Equal(t, docRepo.createdDocs[0].ID, contentRepo.createdContents[0].DocumentID)

	mockProjectRepo.AssertExpectations(t)
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
