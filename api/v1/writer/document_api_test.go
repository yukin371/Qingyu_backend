package writer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	"Qingyu_backend/repository/interfaces/infrastructure"
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
