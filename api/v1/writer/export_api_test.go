//go:build integration
// +build integration

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/models/writer"
	writerMocks "Qingyu_backend/service/writer/mocks"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// ExportAPITestSuite 导出API测试套件
type ExportAPITestSuite struct {
	exportAPI     *writer.ExportApi
	router        *gin.Engine
	mockDocRepo   *writerMocks.MockDocumentRepository
	mockContentRepo *writerMocks.MockDocumentContentRepository
	mockProjectRepo *writerMocks.MockProjectRepository
	mockExportRepo *writerMocks.MockExportTaskRepository
	mockFileStorage *writerMocks.MockFileStorage
}

// setupExportAPITest 设置导出API测试环境
func setupExportAPITest(t *testing.T) *ExportAPITestSuite {
	gin.SetMode(gin.TestMode)

	// 创建Mock Repository
	mockDocRepo := new(writerMocks.MockDocumentRepository)
	mockContentRepo := new(writerMocks.MockDocumentContentRepository)
	mockProjectRepo := new(writerMocks.MockProjectRepository)
	mockExportRepo := new(writerMocks.MockExportTaskRepository)
	mockFileStorage := new(writerMocks.MockFileStorage)

	// 创建ExportService
	exportService := writer.NewExportService(
		mockDocRepo,
		mockContentRepo,
		mockProjectRepo,
		mockExportRepo,
		mockFileStorage,
	)

	// 创建ExportAPI
	api := writer.NewExportApi(exportService)

	// 设置路由
	router := gin.New()
	router.Use(gin.Recovery())

	// 需要认证的路由
	authenticated := router.Group("/api/v1/writer")
	authenticated.Use(mockAuthMiddleware())
	{
		// 文档导出
		authenticated.POST("/documents/:id/export", api.ExportDocument)

		// 项目导出
		authenticated.POST("/projects/:id/export", api.ExportProject)
		authenticated.GET("/projects/:projectId/exports", api.ListExportTasks)

		// 导出任务管理
		authenticated.GET("/exports/:id", api.GetExportTask)
		authenticated.GET("/exports/:id/download", api.DownloadExportFile)
		authenticated.DELETE("/exports/:id", api.DeleteExportTask)
		authenticated.POST("/exports/:id/cancel", api.CancelExportTask)
	}

	return &ExportAPITestSuite{
		exportAPI:       api,
		router:          router,
		mockDocRepo:     mockDocRepo,
		mockContentRepo: mockContentRepo,
		mockProjectRepo: mockProjectRepo,
		mockExportRepo:  mockExportRepo,
		mockFileStorage: mockFileStorage,
	}
}

// ============ ExportDocument API测试 ============

// TestExportDocument_Success 测试导出文档成功
func TestExportDocument_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	documentID := "doc123"
	projectID := "proj123"
	userID := "test-user-id-123"

	document := &writer.Document{
		ID:        documentID,
		ProjectID: projectID,
		Title:     "Test Chapter",
		Type:      "chapter",
		Status:    writer.DocumentStatusCompleted,
		WordCount: 1000,
	}

	// 设置Mock期望
	suite.mockDocRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)
	suite.mockExportRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.ExportTask")).Return(nil)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"format":      "txt",
		"includeMeta": false,
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/export?projectId="+projectID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusAccepted, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(202), response["code"])
	assert.Equal(t, "导出任务已创建", response["message"])
	assert.NotNil(t, response["data"])

	// 验证Mock调用
	suite.mockDocRepo.AssertExpectations(t)
	suite.mockExportRepo.AssertExpectations(t)
}

// TestExportDocument_MissingProjectID 测试缺少项目ID
func TestExportDocument_MissingProjectID(t *testing.T) {
	suite := setupExportAPITest(t)

	documentID := "doc123"

	// 准备请求数据
	reqBody := map[string]interface{}{
		"format": "txt",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求 - 缺少projectId参数
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/export", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "projectId不能为空")
}

// TestExportDocument_DocumentNotFound 测试文档不存在
func TestExportDocument_DocumentNotFound(t *testing.T) {
	suite := setupExportAPITest(t)

	documentID := "nonexistent"
	projectID := "proj123"

	// 设置Mock期望 - 文档不存在
	suite.mockDocRepo.On("FindByID", mock.Anything, documentID).Return(nil, assert.AnError)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"format": "txt",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/export?projectId="+projectID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "创建导出任务失败")

	suite.mockDocRepo.AssertExpectations(t)
}

// TestExportDocument_InvalidFormat 测试无效的格式
func TestExportDocument_InvalidFormat(t *testing.T) {
	suite := setupExportAPITest(t)

	documentID := "doc123"
	projectID := "proj123"

	// 准备请求数据 - 无效的格式
	reqBody := map[string]interface{}{
		"format": "invalid",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/export?projectId="+projectID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回参数错误
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============ ExportProject API测试 ============

// TestExportProject_Success 测试导出项目成功
func TestExportProject_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	projectID := "proj123"
	userID := "test-user-id-123"

	project := &writer.Project{
		ID:        projectID,
		AuthorID:  userID,
		Title:     "Test Project",
		Status:    writer.StatusDraft,
	}

	// 设置Mock期望
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	suite.mockExportRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.ExportTask")).Return(nil)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"includeDocuments": true,
		"documentFormats":  "txt",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/export", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusAccepted, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(202), response["code"])
	assert.Equal(t, "导出任务已创建", response["message"])
	assert.NotNil(t, response["data"])

	suite.mockProjectRepo.AssertExpectations(t)
	suite.mockExportRepo.AssertExpectations(t)
}

// TestExportProject_MissingID 测试缺少项目ID
func TestExportProject_MissingID(t *testing.T) {
	suite := setupExportAPITest(t)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"includeDocuments": true,
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求 - URL中缺少ID
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects//export", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============ GetExportTask API测试 ============

// TestGetExportTask_Success 测试获取导出任务成功
func TestGetExportTask_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"

	task := &serviceInterfaces.ExportTask{
		ID:        taskID,
		Type:      serviceInterfaces.ExportTypeDocument,
		ResourceID: "doc123",
		ResourceTitle: "Test Chapter",
		Format:    serviceInterfaces.ExportFormatTXT,
		Status:    serviceInterfaces.ExportStatusCompleted,
		Progress:  100,
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/exports/"+taskID, nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取成功", response["message"])
	assert.NotNil(t, response["data"])

	suite.mockExportRepo.AssertExpectations(t)
}

// TestGetExportTask_NotFound 测试获取不存在的任务
func TestGetExportTask_NotFound(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "nonexistent"

	// 设置Mock期望 - 任务不存在
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(nil, assert.AnError)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/exports/"+taskID, nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "导出任务不存在")

	suite.mockExportRepo.AssertExpectations(t)
}

// ============ DownloadExportFile API测试 ============

// TestDownloadExportFile_Success 测试下载导出文件成功
func TestDownloadExportFile_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"

	task := &serviceInterfaces.ExportTask{
		ID:            taskID,
		Type:          serviceInterfaces.ExportTypeDocument,
		ResourceTitle: "Test Chapter",
		Format:        serviceInterfaces.ExportFormatTXT,
		Status:        serviceInterfaces.ExportStatusCompleted,
		FileURL:       "/exports/test.txt",
		FileSize:      1024,
		ExpiresAt:     time.Now().Add(24 * time.Hour), // 未来时间
	}

	signedURL := "https://storage.example.com/exports/test.txt?signature=xxx"

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
	suite.mockFileStorage.On("GetSignedURL", mock.Anything, task.FileURL, mock.AnythingOfType("time.Duration")).Return(signedURL, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/exports/"+taskID+"/download", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取成功", response["message"])
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Test Chapter.txt", data["filename"])
	assert.Equal(t, signedURL, data["url"])

	suite.mockExportRepo.AssertExpectations(t)
	suite.mockFileStorage.AssertExpectations(t)
}

// TestDownloadExportFile_TaskNotCompleted 测试任务未完成时下载
func TestDownloadExportFile_TaskNotCompleted(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"

	task := &serviceInterfaces.ExportTask{
		ID:        taskID,
		Status:    serviceInterfaces.ExportStatusProcessing, // 处理中
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/exports/"+taskID+"/download", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "获取下载链接失败")

	suite.mockExportRepo.AssertExpectations(t)
}

// ============ ListExportTasks API测试 ============

// TestListExportTasks_Success 测试列出导出任务成功
func TestListExportTasks_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	projectID := "proj123"

	tasks := []*serviceInterfaces.ExportTask{
		{
			ID:        "task1",
			Type:      serviceInterfaces.ExportTypeDocument,
			ResourceTitle: "Chapter 1",
			Status:    serviceInterfaces.ExportStatusCompleted,
		},
		{
			ID:        "task2",
			Type:      serviceInterfaces.ExportTypeDocument,
			ResourceTitle: "Chapter 2",
			Status:    serviceInterfaces.ExportStatusProcessing,
		},
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByProjectID", mock.Anything, projectID, 1, 20).Return(tasks, int64(2), nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/projects/"+projectID+"/exports?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["pagination"])

	suite.mockExportRepo.AssertExpectations(t)
}

// TestListExportTasks_DefaultPagination 测试默认分页参数
func TestListExportTasks_DefaultPagination(t *testing.T) {
	suite := setupExportAPITest(t)

	projectID := "proj123"

	// 设置Mock期望 - 使用默认分页参数
	suite.mockExportRepo.On("FindByProjectID", mock.Anything, projectID, 1, 20).Return([]*serviceInterfaces.ExportTask{}, int64(0), nil)

	// 创建请求 - 不提供分页参数
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/projects/"+projectID+"/exports", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	suite.mockExportRepo.AssertExpectations(t)
}

// ============ DeleteExportTask API测试 ============

// TestDeleteExportTask_Success 测试删除导出任务成功
func TestDeleteExportTask_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"
	userID := "test-user-id-123"

	task := &serviceInterfaces.ExportTask{
		ID:        taskID,
		CreatedBy: userID,
		FileURL:   "/exports/test.txt",
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
	suite.mockFileStorage.On("Delete", mock.Anything, task.FileURL).Return(nil)
	suite.mockExportRepo.On("Delete", mock.Anything, taskID).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/writer/exports/"+taskID, nil)
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "删除成功", response["message"])

	suite.mockExportRepo.AssertExpectations(t)
	suite.mockFileStorage.AssertExpectations(t)
}

// TestDeleteExportTask_Forbidden 测试无权限删除任务
func TestDeleteExportTask_Forbidden(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"
	userID := "test-user-id-123"
	otherUserID := "other-user-id"

	task := &serviceInterfaces.ExportTask{
		ID:        taskID,
		CreatedBy: otherUserID, // 不是当前用户
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/writer/exports/"+taskID, nil)
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "删除导出任务失败")

	suite.mockExportRepo.AssertExpectations(t)
}

// ============ CancelExportTask API测试 ============

// TestCancelExportTask_Success 测试取消导出任务成功
func TestCancelExportTask_Success(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"
	userID := "test-user-id-123"

	task := &serviceInterfaces.ExportTask{
		ID:        taskID,
		CreatedBy: userID,
		Status:    serviceInterfaces.ExportStatusPending,
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
	suite.mockExportRepo.On("Update", mock.Anything, mock.MatchedBy(func(t *serviceInterfaces.ExportTask) bool {
		return t.Status == serviceInterfaces.ExportStatusCancelled
	})).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/exports/"+taskID+"/cancel", nil)
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "取消成功", response["message"])

	suite.mockExportRepo.AssertExpectations(t)
}

// TestCancelExportTask_InvalidStatus 测试取消不允许取消状态的任务
func TestCancelExportTask_InvalidStatus(t *testing.T) {
	suite := setupExportAPITest(t)

	taskID := "task123"
	userID := "test-user-id-123"

	task := &serviceInterfaces.ExportTask{
		ID:        taskID,
		CreatedBy: userID,
		Status:    serviceInterfaces.ExportStatusCompleted, // 已完成，不能取消
	}

	// 设置Mock期望
	suite.mockExportRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/exports/"+taskID+"/cancel", nil)
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "取消导出任务失败")

	suite.mockExportRepo.AssertExpectations(t)
}
