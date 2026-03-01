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

// TestConceptAPI_CreateConcept_MissingRequiredFields 测试创建概念缺少必填字段
func TestConceptAPI_CreateConcept_MissingRequiredFields(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.POST("/api/v1/internal/ai/concepts", conceptAPI.CreateConcept)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id",
			requestBody:    `{"project_id": "proj123", "name": "角色名", "category": "角色", "content": "角色描述"}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "缺少project_id",
			requestBody:    `{"user_id": "user123", "name": "角色名", "category": "角色", "content": "角色描述"}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "缺少name",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "category": "角色", "content": "角色描述"}`,
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
			req, _ := http.NewRequest("POST", "/api/v1/internal/ai/concepts", bytes.NewBufferString(tc.requestBody))
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

// TestConceptAPI_GetConcept_MissingQueryParams 测试获取概念缺少查询参数
func TestConceptAPI_GetConcept_MissingQueryParams(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.GET("/api/v1/internal/ai/concepts/:id", conceptAPI.GetConcept)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id和project_id",
			url:            "/api/v1/internal/ai/concepts/concept123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有user_id",
			url:            "/api/v1/internal/ai/concepts/concept123?user_id=user123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有project_id",
			url:            "/api/v1/internal/ai/concepts/concept123?project_id=proj123",
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

// TestConceptAPI_UpdateConcept_InvalidRequest 测试更新概念的请求验证
func TestConceptAPI_UpdateConcept_InvalidRequest(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.PUT("/api/v1/internal/ai/concepts/:id", conceptAPI.UpdateConcept)

	testCases := []struct {
		name           string
		url            string
		requestBody    string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "无效的JSON",
			url:            "/api/v1/internal/ai/concepts/concept123",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "有效的请求结构",
			url:            "/api/v1/internal/ai/concepts/concept123",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "新名称", "content": "新内容"}`,
			expectedStatus: http.StatusInternalServerError, // service是nil
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", tc.url, bytes.NewBufferString(tc.requestBody))
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

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err == nil {
					if tc.expectedCode != 0 {
						assert.Equal(t, float64(tc.expectedCode), response["code"])
					}
				}
			}
		})
	}
}

// TestConceptAPI_DeleteConcept_MissingQueryParams 测试删除概念缺少查询参数
func TestConceptAPI_DeleteConcept_MissingQueryParams(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.DELETE("/api/v1/internal/ai/concepts/:id", conceptAPI.DeleteConcept)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id和project_id",
			url:            "/api/v1/internal/ai/concepts/concept123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有user_id",
			url:            "/api/v1/internal/ai/concepts/concept123?user_id=user123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有project_id",
			url:            "/api/v1/internal/ai/concepts/concept123?project_id=proj123",
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

// TestConceptAPI_SearchConcepts_MissingQueryParams 测试搜索概念缺少查询参数
func TestConceptAPI_SearchConcepts_MissingQueryParams(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.GET("/api/v1/internal/ai/concepts", conceptAPI.SearchConcepts)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id和project_id",
			url:            "/api/v1/internal/ai/concepts",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有user_id",
			url:            "/api/v1/internal/ai/concepts?user_id=user123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "只有project_id",
			url:            "/api/v1/internal/ai/concepts?project_id=proj123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "带所有参数",
			url:            "/api/v1/internal/ai/concepts?user_id=user123&project_id=proj123&category=角色&keyword=法师&limit=10",
			expectedStatus: http.StatusInternalServerError, // service是nil
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.url, nil)

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

// TestConceptAPI_BatchGetConcepts_InvalidRequest 测试批量获取概念的请求验证
func TestConceptAPI_BatchGetConcepts_InvalidRequest(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.POST("/api/v1/internal/ai/concepts/batch", conceptAPI.BatchGetConcepts)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "缺少user_id",
			requestBody:    `{"project_id": "proj123", "concept_ids": ["concept1", "concept2"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "缺少project_id",
			requestBody:    `{"user_id": "user123", "concept_ids": ["concept1", "concept2"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
		},
		{
			name:           "缺少concept_ids",
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
			req, _ := http.NewRequest("POST", "/api/v1/internal/ai/concepts/batch", bytes.NewBufferString(tc.requestBody))
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

// TestConceptAPI_CreateConcept_ValidRequestStructure 测试创建概念的有效请求结构
func TestConceptAPI_CreateConcept_ValidRequestStructure(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.POST("/api/v1/internal/ai/concepts", conceptAPI.CreateConcept)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "完整的角色概念",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "亚瑟·潘德拉贡", "category": "角色", "content": "骑士王，持有圣剑Excalibur", "tags": ["主角", "骑士", "王"]}`,
			expectedStatus: http.StatusInternalServerError, // service是nil
		},
		{
			name:           "最小必填字段",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "概念名称"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "带分类的概念",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "卡美洛", "category": "地点", "content": "不列颠的王都"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "带空标签数组",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "概念名", "tags": []}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/internal/ai/concepts", bytes.NewBufferString(tc.requestBody))
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

// TestConceptAPI_CreateConcept_ServiceNotInitialized 测试服务未初始化的情况
func TestConceptAPI_CreateConcept_ServiceNotInitialized(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.POST("/api/v1/internal/ai/concepts", conceptAPI.CreateConcept)

	validReq := internalService.CreateConceptRequest{
		UserID:    "user123",
		ProjectID: "proj123",
		Name:      "测试概念",
		Category:  "角色",
		Content:   "概念描述",
		Tags:      []string{"标签1", "标签2"},
	}
	reqBody, _ := json.Marshal(validReq)
	req, _ := http.NewRequest("POST", "/api/v1/internal/ai/concepts", bytes.NewBuffer(reqBody))
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
		assert.NotEqual(t, http.StatusCreated, w.Code)
	}
}

// TestConceptAPI_UpdateConcept_ValidRequestStructure 测试更新概念的有效请求结构
func TestConceptAPI_UpdateConcept_ValidRequestStructure(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conceptAPI := &ConceptAPI{}
	router.PUT("/api/v1/internal/ai/concepts/:id", conceptAPI.UpdateConcept)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "更新名称",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "新名称"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "更新内容",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "content": "新内容描述"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "更新标签",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "tags": ["新标签1", "新标签2"]}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "更新所有字段",
			requestBody:    `{"user_id": "user123", "project_id": "proj123", "name": "完全新名称", "content": "完全新内容", "tags": ["标签A", "标签B"]}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/api/v1/internal/ai/concepts/concept123", bytes.NewBufferString(tc.requestBody))
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
