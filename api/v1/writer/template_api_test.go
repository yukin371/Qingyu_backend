package writer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	writerModel "Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/writer/document"
)

type captureTemplateRepository struct {
	filter *writerRepo.TemplateFilter
}

func (r *captureTemplateRepository) Create(ctx context.Context, template *writerModel.Template) error {
	return nil
}

func (r *captureTemplateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*writerModel.Template, error) {
	return &writerModel.Template{}, nil
}

func (r *captureTemplateRepository) Update(ctx context.Context, template *writerModel.Template) error {
	return nil
}

func (r *captureTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (r *captureTemplateRepository) ListByProject(ctx context.Context, projectID primitive.ObjectID, filter *writerRepo.TemplateFilter) ([]*writerModel.Template, error) {
	r.filter = filter
	return []*writerModel.Template{}, nil
}

func (r *captureTemplateRepository) ListByWorkspace(ctx context.Context, workspaceID string, filter *writerRepo.TemplateFilter) ([]*writerModel.Template, error) {
	r.filter = filter
	return []*writerModel.Template{}, nil
}

func (r *captureTemplateRepository) ListGlobal(ctx context.Context, filter *writerRepo.TemplateFilter) ([]*writerModel.Template, error) {
	r.filter = filter
	return []*writerModel.Template{}, nil
}

func (r *captureTemplateRepository) CountByProject(ctx context.Context, projectID primitive.ObjectID) (int64, error) {
	return 0, nil
}

func (r *captureTemplateRepository) ExistsByName(ctx context.Context, projectID *primitive.ObjectID, name string, excludeID *primitive.ObjectID) (bool, error) {
	return false, nil
}

func TestTemplateAPIListTemplatesDefaultsInvalidSortOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &captureTemplateRepository{}
	service := document.NewTemplateService(repo, zap.NewNop())
	api := NewTemplateAPI(service, zap.NewNop())
	router := gin.New()
	router.GET("/templates", api.ListTemplates)

	req, _ := http.NewRequest(http.MethodGet, "/templates?sortOrder=invalid", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if repo.filter == nil {
		t.Fatal("expected repository filter to be captured")
	}
	if repo.filter.SortOrder != -1 {
		t.Fatalf("expected invalid sortOrder to default to -1, got %d", repo.filter.SortOrder)
	}
}

func TestTemplateAPIListTemplatesKeepsValidSortOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &captureTemplateRepository{}
	service := document.NewTemplateService(repo, zap.NewNop())
	api := NewTemplateAPI(service, zap.NewNop())
	router := gin.New()
	router.GET("/templates", api.ListTemplates)

	req, _ := http.NewRequest(http.MethodGet, "/templates?sortOrder=1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if repo.filter == nil {
		t.Fatal("expected repository filter to be captured")
	}
	if repo.filter.SortOrder != 1 {
		t.Fatalf("expected sortOrder 1, got %d", repo.filter.SortOrder)
	}
}
