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

// PublishAPITestSuite 发布API测试套件
type PublishAPITestSuite struct {
	publishAPI     *writer.PublishApi
	router         *gin.Engine
	mockProjectRepo  *writerMocks.MockProjectRepository
	mockDocumentRepo *writerMocks.MockDocumentRepository
	mockPublicationRepo *writerMocks.MockPublicationRepository
	mockBookstoreClient *writerMocks.MockBookstoreClient
	mockEventBus     *writerMocks.MockEventBus
}

// setupPublishAPITest 设置发布API测试环境
func setupPublishAPITest(t *testing.T) *PublishAPITestSuite {
	gin.SetMode(gin.TestMode)

	// 创建Mock Repository
	mockProjectRepo := new(writerMocks.MockProjectRepository)
	mockDocumentRepo := new(writerMocks.MockDocumentRepository)
	mockPublicationRepo := new(writerMocks.MockPublicationRepository)
	mockBookstoreClient := new(writerMocks.MockBookstoreClient)
	mockEventBus := new(writerMocks.MockEventBus)

	// 创建PublishService
	publishService := writer.NewPublishService(
		mockProjectRepo,
		mockDocumentRepo,
		mockPublicationRepo,
		mockBookstoreClient,
		mockEventBus,
	)

	// 创建PublishAPI
	api := writer.NewPublishApi(publishService)

	// 设置路由
	router := gin.New()
	router.Use(gin.Recovery())

	// 需要认证的路由
	authenticated := router.Group("/api/v1/writer")
	authenticated.Use(mockAuthMiddleware())
	{
		// 项目发布
		authenticated.POST("/projects/:id/publish", api.PublishProject)
		authenticated.POST("/projects/:id/unpublish", api.UnpublishProject)
		authenticated.GET("/projects/:id/publication-status", api.GetProjectPublicationStatus)
		authenticated.GET("/projects/:projectId/publications", api.GetPublicationRecords)

		// 文档发布
		authenticated.POST("/documents/:id/publish", api.PublishDocument)
		authenticated.PUT("/documents/:id/publish-status", api.UpdateDocumentPublishStatus)

		// 批量发布
		authenticated.POST("/projects/:projectId/documents/batch-publish", api.BatchPublishDocuments)

		// 发布记录
		authenticated.GET("/publications/:id", api.GetPublicationRecord)
	}

	return &PublishAPITestSuite{
		publishAPI:          api,
		router:              router,
		mockProjectRepo:     mockProjectRepo,
		mockDocumentRepo:    mockDocumentRepo,
		mockPublicationRepo: mockPublicationRepo,
		mockBookstoreClient: mockBookstoreClient,
		mockEventBus:        mockEventBus,
	}
}

// ============ PublishProject API测试 ============

// TestPublishProject_Success 测试发布项目成功
func TestPublishProject_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"
	userID := "test-user-id-123"

	project := &writer.Project{
		ID:       projectID,
		AuthorID: userID,
		Title:    "Test Project",
		Status:   writer.StatusDraft,
	}

	// 准备请求数据
	price := 9.99
	reqBody := map[string]interface{}{
		"bookstoreId":    "bookstore123",
		"categoryId":     "category123",
		"tags":           []string{"小说", "玄幻"},
		"description":    "测试项目描述",
		"publishType":    "serial",
		"price":          price,
		"freeChapters":   10,
		"authorNote":     "作者的话",
		"enableComment":  true,
		"enableShare":    true,
	}
	body, _ := json.Marshal(reqBody)

	// 设置Mock期望
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil)
	suite.mockPublicationRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.PublicationRecord")).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBuffer(body))
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
	assert.Equal(t, "发布任务已创建", response["message"])
	assert.NotNil(t, response["data"])

	suite.mockProjectRepo.AssertExpectations(t)
	suite.mockPublicationRepo.AssertExpectations(t)
}

// TestPublishProject_ProjectNotFound 测试项目不存在
func TestPublishProject_ProjectNotFound(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "nonexistent"
	userID := "test-user-id-123"

	// 准备请求数据
	reqBody := map[string]interface{}{
		"bookstoreId": "bookstore123",
		"categoryId":  "category123",
		"publishType": "serial",
	}
	body, _ := json.Marshal(reqBody)

	// 设置Mock期望 - 项目不存在
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(nil, assert.AnError)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "发布项目失败")

	suite.mockProjectRepo.AssertExpectations(t)
}

// TestPublishProject_AlreadyPublished 测试项目已发布
func TestPublishProject_AlreadyPublished(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"
	userID := "test-user-id-123"

	project := &writer.Project{
		ID:       projectID,
		AuthorID: userID,
		Title:    "Test Project",
	}

	existingRecord := &serviceInterfaces.PublicationRecord{
		ID:         "record123",
		ResourceID: projectID,
		Status:     serviceInterfaces.PublicationStatusPublished,
	}

	// 准备请求数据
	reqBody := map[string]interface{}{
		"bookstoreId": "bookstore123",
		"categoryId":  "category123",
		"publishType": "serial",
	}
	body, _ := json.Marshal(reqBody)

	// 设置Mock期望 - 项目已发布
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(existingRecord, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "发布项目失败")

	suite.mockProjectRepo.AssertExpectations(t)
	suite.mockPublicationRepo.AssertExpectations(t)
}

// TestPublishProject_InvalidRequestBody 测试无效的请求体
func TestPublishProject_InvalidRequestBody(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"

	// 准备请求数据 - 缺少必填字段
	reqBody := map[string]interface{}{
		"bookstoreId": "bookstore123",
		// 缺少 categoryId
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回参数错误
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============ UnpublishProject API测试 ============

// TestUnpublishProject_Success 测试取消发布成功
func TestUnpublishProject_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"
	userID := "test-user-id-123"
	bookstoreID := "bookstore123"

	record := &serviceInterfaces.PublicationRecord{
		ID:          "record123",
		ResourceID:  projectID,
		BookstoreID: bookstoreID,
		Status:      serviceInterfaces.PublicationStatusPublished,
		CreatedBy:   userID,
	}

	// 设置Mock期望
	suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(record, nil)
	suite.mockBookstoreClient.On("UnpublishProject", mock.Anything, projectID, bookstoreID).Return(nil)
	suite.mockPublicationRepo.On("Update", mock.Anything, mock.MatchedBy(func(r *serviceInterfaces.PublicationRecord) bool {
		return r.Status == serviceInterfaces.PublicationStatusUnpublished
	})).Return(nil)
	suite.mockEventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/unpublish", nil)
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
	assert.Equal(t, "取消发布成功", response["message"])

	suite.mockPublicationRepo.AssertExpectations(t)
	suite.mockBookstoreClient.AssertExpectations(t)
	suite.mockEventBus.AssertExpectations(t)
}

// TestUnpublishProject_NotPublished 测试项目未发布
func TestUnpublishProject_NotPublished(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"
	userID := "test-user-id-123"

	// 设置Mock期望 - 没有发布记录
	suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, assert.AnError)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/unpublish", nil)
	req.Header.Set("X-Test-User-ID", userID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "取消发布失败")

	suite.mockPublicationRepo.AssertExpectations(t)
}

// ============ GetProjectPublicationStatus API测试 ============

// TestGetProjectPublicationStatus_Success 测试获取项目发布状态成功
func TestGetProjectPublicationStatus_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"

	project := &writer.Project{
		ID:       projectID,
		AuthorID: "user123",
		Title:    "Test Project",
	}

	bookstoreID := "bookstore123"
	publishTime := time.Now()
	record := &serviceInterfaces.PublicationRecord{
		ID:            "record123",
		ResourceID:    projectID,
		BookstoreID:   bookstoreID,
		Status:        serviceInterfaces.PublicationStatusPublished,
		PublishTime:   &publishTime,
	}

	documents := []*writer.Document{
		{
			ID:        "doc1",
			ProjectID: projectID,
			Title:     "Chapter 1",
		},
		{
			ID:        "doc2",
			ProjectID: projectID,
			Title:     "Chapter 2",
		},
	}

	statistics := &serviceInterfaces.PublicationStatistics{
		TotalViews:  1000,
		TotalLikes:  100,
		TotalComments: 50,
		LastSyncTime: time.Now(),
	}

	// 设置Mock期望
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(record, nil)
	suite.mockBookstoreClient.On("GetStatistics", mock.Anything, projectID, bookstoreID).Return(statistics, nil)
	suite.mockDocumentRepo.On("FindByProjectID", mock.Anything, projectID).Return(documents, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/projects/"+projectID+"/publication-status", nil)
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
	assert.True(t, data["isPublished"].(bool))
	assert.Equal(t, "Test Project", data["projectTitle"])
	assert.Equal(t, float64(2), data["totalChapters"])

	suite.mockProjectRepo.AssertExpectations(t)
	suite.mockPublicationRepo.AssertExpectations(t)
	suite.mockBookstoreClient.AssertExpectations(t)
	suite.mockDocumentRepo.AssertExpectations(t)
}

// ============ PublishDocument API测试 ============

// TestPublishDocument_Success 测试发布文档成功
func TestPublishDocument_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	documentID := "doc123"
	projectID := "proj123"
	userID := "test-user-id-123"

	document := &writer.Document{
		ID:        documentID,
		ProjectID: projectID,
		Title:     "Test Chapter",
		Type:      "chapter",
		Status:    writer.DocumentStatusCompleted,
	}

	project := &writer.Project{
		ID:       projectID,
		AuthorID: userID,
		Title:    "Test Project",
	}

	// 准备请求数据
	reqBody := map[string]interface{}{
		"chapterTitle":  "第一章",
		"chapterNumber": 1,
		"isFree":        false,
		"authorNote":    "章节作者的话",
	}
	body, _ := json.Marshal(reqBody)

	// 设置Mock期望
	suite.mockDocumentRepo.On("FindByID", mock.Anything, documentID).Return(document, nil)
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)
	suite.mockPublicationRepo.On("FindByResourceID", mock.Anything, documentID).Return(nil, nil)
	suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil)
	suite.mockPublicationRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.PublicationRecord")).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/publish?projectId="+projectID, bytes.NewBuffer(body))
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
	assert.Equal(t, "发布任务已创建", response["message"])
	assert.NotNil(t, response["data"])

	suite.mockDocumentRepo.AssertExpectations(t)
	suite.mockProjectRepo.AssertExpectations(t)
	suite.mockPublicationRepo.AssertExpectations(t)
}

// TestPublishDocument_MissingProjectID 测试缺少项目ID
func TestPublishDocument_MissingProjectID(t *testing.T) {
	suite := setupPublishAPITest(t)

	documentID := "doc123"

	// 准备请求数据
	reqBody := map[string]interface{}{
		"chapterTitle":  "第一章",
		"chapterNumber": 1,
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求 - 缺少projectId参数
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/documents/"+documentID+"/publish", bytes.NewBuffer(body))
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

// ============ BatchPublishDocuments API测试 ============

// TestBatchPublishDocuments_Success 测试批量发布文档成功
func TestBatchPublishDocuments_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"
	userID := "test-user-id-123"

	project := &writer.Project{
		ID:       projectID,
		AuthorID: userID,
		Title:    "Test Project",
	}

	documentIDs := []string{"doc1", "doc2"}

	doc1 := &writer.Document{
		ID:        "doc1",
		ProjectID: projectID,
		Title:     "Chapter 1",
	}

	doc2 := &writer.Document{
		ID:        "doc2",
		ProjectID: projectID,
		Title:     "Chapter 2",
	}

	// 准备请求数据
	reqBody := map[string]interface{}{
		"documentIds":   documentIDs,
		"autoNumbering": true,
		"startNumber":   1,
		"isFree":        false,
	}
	body, _ := json.Marshal(reqBody)

	// 设置Mock期望
	suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil)

	for _, doc := range []*writer.Document{doc1, doc2} {
		suite.mockDocumentRepo.On("FindByID", mock.Anything, doc.ID).Return(doc, nil).Once()
		suite.mockProjectRepo.On("FindByID", mock.Anything, projectID).Return(project, nil).Once()
		suite.mockPublicationRepo.On("FindByResourceID", mock.Anything, doc.ID).Return(nil, nil).Once()
		suite.mockPublicationRepo.On("FindPublishedByProjectID", mock.Anything, projectID).Return(nil, nil).Once()
		suite.mockPublicationRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.PublicationRecord")).Return(nil).Once()
	}

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents/batch-publish", bytes.NewBuffer(body))
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
	assert.Equal(t, "批量发布任务已创建", response["message"])
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["successCount"])
	assert.Equal(t, float64(0), data["failCount"])

	suite.mockProjectRepo.AssertExpectations(t)
	suite.mockDocumentRepo.AssertExpectations(t)
	suite.mockPublicationRepo.AssertExpectations(t)
}

// TestBatchPublishDocuments_EmptyDocumentIDs 测试空文档ID列表
func TestBatchPublishDocuments_EmptyDocumentIDs(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"

	// 准备请求数据 - 空的文档ID列表
	reqBody := map[string]interface{}{
		"documentIds":   []string{},
		"autoNumbering": true,
		"startNumber":   1,
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writer/projects/"+projectID+"/documents/batch-publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回参数错误
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============ GetPublicationRecords API测试 ============

// TestGetPublicationRecords_Success 测试获取发布记录列表成功
func TestGetPublicationRecords_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	projectID := "proj123"

	records := []*serviceInterfaces.PublicationRecord{
		{
			ID:            "record1",
			Type:          "document",
			ResourceID:    "doc1",
			ResourceTitle: "Chapter 1",
			Status:        serviceInterfaces.PublicationStatusPublished,
		},
		{
			ID:            "record2",
			Type:          "document",
			ResourceID:    "doc2",
			ResourceTitle: "Chapter 2",
			Status:        serviceInterfaces.PublicationStatusPending,
		},
	}

	// 设置Mock期望
	suite.mockPublicationRepo.On("FindByProjectID", mock.Anything, projectID, 1, 20).Return(records, int64(2), nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/projects/"+projectID+"/publications?page=1&pageSize=20", nil)
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

	suite.mockPublicationRepo.AssertExpectations(t)
}

// ============ GetPublicationRecord API测试 ============

// TestGetPublicationRecord_Success 测试获取发布记录详情成功
func TestGetPublicationRecord_Success(t *testing.T) {
	suite := setupPublishAPITest(t)

	recordID := "record123"

	record := &serviceInterfaces.PublicationRecord{
		ID:            recordID,
		Type:          "document",
		ResourceID:    "doc1",
		ResourceTitle: "Chapter 1",
		Status:        serviceInterfaces.PublicationStatusPublished,
		BookstoreID:   "bookstore123",
	}

	// 设置Mock期望
	suite.mockPublicationRepo.On("FindByID", mock.Anything, recordID).Return(record, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/publications/"+recordID, nil)
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
	assert.Equal(t, recordID, data["id"])

	suite.mockPublicationRepo.AssertExpectations(t)
}

// TestGetPublicationRecord_NotFound 测试获取不存在的发布记录
func TestGetPublicationRecord_NotFound(t *testing.T) {
	suite := setupPublishAPITest(t)

	recordID := "nonexistent"

	// 设置Mock期望 - 记录不存在
	suite.mockPublicationRepo.On("FindByID", mock.Anything, recordID).Return(nil, assert.AnError)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writer/publications/"+recordID, nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "发布记录不存在")

	suite.mockPublicationRepo.AssertExpectations(t)
}
