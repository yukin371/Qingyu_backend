package writer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/service/interfaces"
)

type mockExportAPIService struct {
	mock.Mock
}

func (m *mockExportAPIService) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *interfaces.ExportDocumentRequest) (*interfaces.ExportTask, error) {
	args := m.Called(ctx, documentID, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportTask), args.Error(1)
}

func (m *mockExportAPIService) GetExportTask(ctx context.Context, taskID string) (*interfaces.ExportTask, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportTask), args.Error(1)
}

func (m *mockExportAPIService) DownloadExportFile(ctx context.Context, taskID string) (*interfaces.ExportFile, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportFile), args.Error(1)
}

func (m *mockExportAPIService) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*interfaces.ExportTask, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*interfaces.ExportTask), args.Get(1).(int64), args.Error(2)
}

func (m *mockExportAPIService) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	return m.Called(ctx, taskID, userID).Error(0)
}

func (m *mockExportAPIService) ExportProject(ctx context.Context, projectID, userID string, req *interfaces.ExportProjectRequest) (*interfaces.ExportTask, error) {
	args := m.Called(ctx, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportTask), args.Error(1)
}

func (m *mockExportAPIService) ExportProjectAsZip(ctx context.Context, projectID, userID string) ([]byte, error) {
	args := m.Called(ctx, projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockExportAPIService) ImportProject(ctx context.Context, userID string, zipData []byte) (*interfaces.ImportResult, error) {
	args := m.Called(ctx, userID, zipData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ImportResult), args.Error(1)
}

func (m *mockExportAPIService) CancelExportTask(ctx context.Context, taskID, userID string) error {
	return m.Called(ctx, taskID, userID).Error(0)
}

func TestExportApi_DownloadExportFile_StreamsBinary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := new(mockExportAPIService)
	api := NewExportApi(service)

	router := gin.New()
	router.GET("/api/v1/writer/exports/:id/download", api.DownloadExportFile)

	service.On("DownloadExportFile", mock.Anything, "task-1").Return(&interfaces.ExportFile{
		Filename: "chapter.md",
		Content:  []byte("# hello"),
		MimeType: "text/markdown; charset=utf-8",
		FileSize: int64(len("# hello")),
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/exports/task-1/download", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "text/markdown; charset=utf-8", resp.Header().Get("Content-Type"))
	assert.Contains(t, resp.Header().Get("Content-Disposition"), "attachment; filename=\"chapter.md\"")
	assert.Equal(t, "# hello", resp.Body.String())
	service.AssertExpectations(t)
}

func TestExportApi_DownloadExportFile_FallsBackToJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := new(mockExportAPIService)
	api := NewExportApi(service)

	router := gin.New()
	router.GET("/api/v1/writer/exports/:id/download", api.DownloadExportFile)

	service.On("DownloadExportFile", mock.Anything, "task-2").Return(&interfaces.ExportFile{
		Filename: "chapter.md",
		URL:      "/signed/task-2",
		MimeType: "text/markdown; charset=utf-8",
		FileSize: 128,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/exports/task-2/download", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var payload map[string]any
	err := json.Unmarshal(resp.Body.Bytes(), &payload)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), payload["code"])
	service.AssertExpectations(t)
}
