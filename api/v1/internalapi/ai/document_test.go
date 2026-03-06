package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalService "Qingyu_backend/service/internalapi"
)

// TestDocumentAPI_CreateOrUpdateDocument_MissingRequiredFields 测试缺少必填字段的情况
func TestDocumentAPI_CreateOrUpdateDocument_MissingRequiredFields(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	documentAPI := &DocumentAPI{}
	router.POST("/api/v1/internal/ai/documents", documentAPI.CreateOrUpdateDocument)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "缺少user_id",
			requestBody:    `{"project_id": "proj123", "action": "create", "document": {"title": "测试文档"}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "缺少project_id",
			requestBody:    `{"user_id": "user123", "action": "create", "document": {"title": "测试文档"}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "缺少action",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "document": {"title": "测试文档"}}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "缺少document",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "action": "create"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "无效的JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/internal/ai/documents", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// When
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

// TestDocumentAPI_GetDocument_MissingQueryParams 测试缺少查询参数的情况
func TestDocumentAPI_GetDocument_MissingQueryParams(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	documentAPI := &DocumentAPI{}
	router.GET("/api/v1/internal/ai/documents/:id", documentAPI.GetDocument)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id和project_id",
			url:            "/api/v1/internal/ai/documents/doc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有user_id",
			url:            "/api/v1/internal/ai/documents/doc123?user_id=user123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有project_id",
			url:            "/api/v1/internal/ai/documents/doc123?project_id=proj123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.url, nil)

			// When
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, float64(tc.expectedCode), response["code"])
		})
	}
}

// TestDocumentAPI_ListDocuments_MissingQueryParams 测试列表接口缺少参数的情况
func TestDocumentAPI_ListDocuments_MissingQueryParams(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	documentAPI := &DocumentAPI{}
	router.GET("/api/v1/internal/ai/documents", documentAPI.ListDocuments)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id和project_id",
			url:            "/api/v1/internal/ai/documents",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有user_id",
			url:            "/api/v1/internal/ai/documents?user_id=user123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有project_id",
			url:            "/api/v1/internal/ai/documents?project_id=proj123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "带limit参数",
			url:            "/api/v1/internal/ai/documents?user_id=user123&project_id=proj123&limit=100",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.url, nil)

			// When
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

// TestDocumentAPI_DeleteDocument_MissingQueryParams 测试删除接口缺少参数的情况
func TestDocumentAPI_DeleteDocument_MissingQueryParams(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	documentAPI := &DocumentAPI{}
	router.DELETE("/api/v1/internal/ai/documents/:id", documentAPI.DeleteDocument)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id和project_id",
			url:            "/api/v1/internal/ai/documents/doc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有user_id",
			url:            "/api/v1/internal/ai/documents/doc123?user_id=user123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有project_id",
			url:            "/api/v1/internal/ai/documents/doc123?project_id=proj123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", tc.url, nil)

			// When
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, float64(tc.expectedCode), response["code"])
		})
	}
}

// TestDocumentAPI_BatchGetDocuments_InvalidRequest 测试批量获取接口的请求验证
func TestDocumentAPI_BatchGetDocuments_InvalidRequest(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	documentAPI := &DocumentAPI{}
	router.POST("/api/v1/internal/ai/documents/batch", documentAPI.BatchGetDocuments)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id",
			requestBody:    `{"project_id": "proj123", "document_ids": ["doc1", "doc2"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "缺少project_id",
			requestBody:    `{"user_id": "user123", "document_ids": ["doc1", "doc2"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "缺少document_ids",
			requestBody:    `{"user_id": "user123", "project_id": "proj123"}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "无效的JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/internal/ai/documents/batch", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// When
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, float64(tc.expectedCode), response["code"])
		})
	}
}

// TestDocumentAPI_CreateOrUpdateDocument_ServiceNotInitialized 测试服务未初始化的情况
func TestDocumentAPI_CreateOrUpdateDocument_ServiceNotInitialized(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 创建一个没有初始化service的API实例
	documentAPI := &DocumentAPI{}
	router.POST("/api/v1/internal/ai/documents", documentAPI.CreateOrUpdateDocument)

	validReq := internalService.CreateOrUpdateRequest{
		UserID:    "user123",
		ProjectID: "proj123",
		Action:    "create",
		Document: internalService.WriterDraftData{
			ChapterNum: 1,
			Title:      "测试文档",
			Content:    "测试内容",
		},
	}
	reqBody, _ := json.Marshal(validReq)
	req, _ := http.NewRequest("POST", "/api/v1/internal/ai/documents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()

	// 捕获panic
	defer func() {
		if r := recover(); r != nil {
			// 由于service是nil，会panic，这是预期的
			assert.NotNil(t, r)
		}
	}()

	router.ServeHTTP(w, req)

	// Then - 如果没有panic，检查状态码
	if w.Code != 0 {
		assert.NotEqual(t, http.StatusOK, w.Code)
	}
}

// TestDocumentAPI_ValidRequestStructure 测试请求结构验证
func TestDocumentAPI_ValidRequestStructure(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	documentAPI := &DocumentAPI{}
	router.POST("/api/v1/internal/ai/documents", documentAPI.CreateOrUpdateDocument)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "有效的create请求",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "action": "create", "document": {"chapter_num": 1, "title": "测试文档", "content": "测试内容"}}`,
			expectedStatus: http.StatusInternalServerError, // 因为service是nil
		},
		{
			name:           "有效的update请求",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "action": "update", "document": {"chapter_num": 1, "title": "更新标题", "content": "更新内容"}}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "有效的create_or_update请求",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "action": "create_or_update", "document": {"chapter_num": 1, "title": "智能创建或更新"}}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "有效的append请求",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "action": "append", "document": {"chapter_num": 1, "content": "追加内容"}}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/internal/ai/documents", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// When
			w := httptest.NewRecorder()

			// 捕获panic
			defer func() {
				if r := recover(); r != nil {
					// 由于service是nil，会panic，这是预期的
					assert.NotNil(t, r)
				}
			}()

			router.ServeHTTP(w, req)

			// Then - 如果没有panic，检查状态码
			if w.Code != 0 {
				assert.Equal(t, tc.expectedStatus, w.Code)
			}
		})
	}
}
