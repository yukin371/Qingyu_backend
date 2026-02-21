package writer_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
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

// MockImportExportService Mock导入导出服务
type MockImportExportService struct {
	mock.Mock
}

// ExportDocument 导出文档
func (m *MockImportExportService) ExportDocument(ctx context.Context, documentID, projectID, userID string, req *interfaces.ExportDocumentRequest) (*interfaces.ExportTask, error) {
	args := m.Called(ctx, documentID, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportTask), args.Error(1)
}

// GetExportTask 获取导出任务
func (m *MockImportExportService) GetExportTask(ctx context.Context, taskID string) (*interfaces.ExportTask, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportTask), args.Error(1)
}

// DownloadExportFile 下载导出文件
func (m *MockImportExportService) DownloadExportFile(ctx context.Context, taskID string) (*interfaces.ExportFile, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportFile), args.Error(1)
}

// ListExportTasks 列出导出任务
func (m *MockImportExportService) ListExportTasks(ctx context.Context, projectID string, page, pageSize int) ([]*interfaces.ExportTask, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*interfaces.ExportTask), args.Get(1).(int64), args.Error(2)
}

// DeleteExportTask 删除导出任务
func (m *MockImportExportService) DeleteExportTask(ctx context.Context, taskID, userID string) error {
	args := m.Called(ctx, taskID, userID)
	return args.Error(0)
}

// ExportProject 导出项目（异步任务）
func (m *MockImportExportService) ExportProject(ctx context.Context, projectID, userID string, req *interfaces.ExportProjectRequest) (*interfaces.ExportTask, error) {
	args := m.Called(ctx, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ExportTask), args.Error(1)
}

// ExportProjectAsZip 将项目导出为ZIP字节数据（直接返回）
func (m *MockImportExportService) ExportProjectAsZip(ctx context.Context, projectID, userID string) ([]byte, error) {
	args := m.Called(ctx, projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// ImportProject 从ZIP数据导入项目
func (m *MockImportExportService) ImportProject(ctx context.Context, userID string, zipData []byte) (*interfaces.ImportResult, error) {
	args := m.Called(ctx, userID, zipData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ImportResult), args.Error(1)
}

// CancelExportTask 取消导出任务
func (m *MockImportExportService) CancelExportTask(ctx context.Context, taskID, userID string) error {
	args := m.Called(ctx, taskID, userID)
	return args.Error(0)
}

// setupImportExportTestRouter 设置测试路由
func setupImportExportTestRouter(service *MockImportExportService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	importExportAPI := writer.NewImportExportApi(service)
	r.GET("/api/v1/writer/projects/:id/export", importExportAPI.ExportProject)
	r.POST("/api/v1/writer/projects/import", importExportAPI.ImportProject)

	return r
}

// TestImportExportApi_ExportProject_Success 测试成功导出项目
func TestImportExportApi_ExportProject_Success(t *testing.T) {
	// Given
	mockService := new(MockImportExportService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupImportExportTestRouter(mockService, userID)

	// 创建一个简单的 ZIP 文件作为返回值
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	zipFile, _ := zipWriter.Create("测试项目/第一章.txt")
	zipFile.Write([]byte("测试内容"))
	zipWriter.Close()
	mockZipData := buf.Bytes()

	mockService.On("ExportProjectAsZip", mock.Anything, projectID, userID).Return(mockZipData, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/export", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Then
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/zip", recorder.Header().Get("Content-Type"))
	assert.Contains(t, recorder.Header().Get("Content-Disposition"), "attachment")
	assert.True(t, len(recorder.Body.Bytes()) > 0)

	mockService.AssertExpectations(t)
}

// TestImportExportApi_ExportProject_MissingUserID 测试缺少用户ID
func TestImportExportApi_ExportProject_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockImportExportService)
	projectID := primitive.NewObjectID().Hex()
	router := setupImportExportTestRouter(mockService, "")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/export", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1002), response["code"]) // 1002 = Unauthorized
}

// TestImportExportApi_ExportProject_ServiceError 测试服务错误
func TestImportExportApi_ExportProject_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockImportExportService)
	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	router := setupImportExportTestRouter(mockService, userID)

	mockService.On("ExportProjectAsZip", mock.Anything, projectID, userID).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/writer/projects/"+projectID+"/export", nil)
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

// TestImportExportApi_ImportProject_Success 测试成功导入项目
func TestImportExportApi_ImportProject_Success(t *testing.T) {
	// Given
	mockService := new(MockImportExportService)
	userID := primitive.NewObjectID().Hex()
	router := setupImportExportTestRouter(mockService, userID)

	// 创建一个模拟的 ZIP 文件
	body := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(body)
	part, _ := mpWriter.CreateFormFile("file", "test.zip")
	zipBuf := createTestZip(t)
	part.Write(zipBuf)
	mpWriter.Close()

	expectedResult := &interfaces.ImportResult{
		ProjectID:     primitive.NewObjectID().Hex(),
		Title:         "测试项目",
		DocumentCount: 3,
	}

	mockService.On("ImportProject", mock.Anything, userID, mock.Anything).Return(expectedResult, nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/import", body)
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"]) // 0 = Success

	data := response["data"].(map[string]interface{})
	assert.Equal(t, expectedResult.ProjectID, data["projectId"])
	assert.Equal(t, expectedResult.Title, data["title"])
	assert.Equal(t, float64(expectedResult.DocumentCount), data["documentCount"])

	mockService.AssertExpectations(t)
}

// TestImportExportApi_ImportProject_MissingFile 测试缺少文件
func TestImportExportApi_ImportProject_MissingFile(t *testing.T) {
	// Given
	mockService := new(MockImportExportService)
	userID := primitive.NewObjectID().Hex()
	router := setupImportExportTestRouter(mockService, userID)

	body := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(body)
	mpWriter.Close()

	// When
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/import", body)
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 1001 = InvalidParams
}

// TestImportExportApi_ImportProject_InvalidFileType 测试无效的文件类型
func TestImportExportApi_ImportProject_InvalidFileType(t *testing.T) {
	// Given
	mockService := new(MockImportExportService)
	userID := primitive.NewObjectID().Hex()
	router := setupImportExportTestRouter(mockService, userID)

	body := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(body)
	part, _ := mpWriter.CreateFormFile("file", "test.txt")
	part.Write([]byte("not a zip file"))
	mpWriter.Close()

	// When
	req, _ := http.NewRequest("POST", "/api/v1/writer/projects/import", body)
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1001), response["code"]) // 1001 = InvalidParams
}

// createTestZip 创建测试用的 ZIP 文件
func createTestZip(t *testing.T) []byte {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// 创建项目目录和文件
	w1, err := zipWriter.Create("测试项目/第一章.txt")
	assert.NoError(t, err)
	_, err = w1.Write([]byte("这是第一章的内容"))
	assert.NoError(t, err)

	w2, err := zipWriter.Create("测试项目/第二章.txt")
	assert.NoError(t, err)
	_, err = w2.Write([]byte("这是第二章的内容"))
	assert.NoError(t, err)

	err = zipWriter.Close()
	assert.NoError(t, err)

	return buf.Bytes()
}
