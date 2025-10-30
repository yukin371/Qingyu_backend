package api

import (
	documentModel "Qingyu_backend/models/writer"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/document"
)

// setupSimpleEditorRouter 设置简化的测试路由（用于不需要DocumentService的测试）
func setupSimpleEditorRouter(userID string) (*gin.Engine, *writer.EditorApi) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 模拟认证中间件
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userID", userID)
			ctx := context.WithValue(c.Request.Context(), "userID", userID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	})

	// 创建DocumentService（用于EditorApi初始化）
	mockDocRepo := new(MockDocumentRepository)
	mockDocContentRepo := new(MockDocumentContentRepository)
	mockProjRepo := new(MockProjectRepository)
	eventBus := &MockEventBus{}
	docService := document.NewDocumentService(mockDocRepo, mockDocContentRepo, mockProjRepo, eventBus)

	editorApi := writer.NewEditorApi(docService)

	// 注册路由
	documents := r.Group("/api/v1/documents")
	{
		documents.POST("/:id/word-count", editorApi.CalculateWordCount)
	}

	user := r.Group("/api/v1/user")
	{
		user.GET("/shortcuts", editorApi.GetUserShortcuts)
		user.PUT("/shortcuts", editorApi.UpdateUserShortcuts)
		user.POST("/shortcuts/reset", editorApi.ResetUserShortcuts)
		user.GET("/shortcuts/help", editorApi.GetShortcutHelp)
	}

	return r, editorApi
}

// TestEditorApi_CalculateWordCount 测试计算字数
func TestEditorApi_CalculateWordCount(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		requestBody    writer.WordCountRequest
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:       "成功计算字数（不过滤Markdown）",
			documentID: "doc123",
			requestBody: writer.WordCountRequest{
				Content:        "这是一段测试文本，包含一些**粗体**和*斜体*内容。",
				FilterMarkdown: false,
			},
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "计算成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				// 注意：JSON字段是camelCase，如totalCount而不是total_words
				assert.NotZero(t, data["totalCount"])
				assert.NotZero(t, data["chineseCount"])
				assert.Equal(t, float64(1), data["readingTime"])
			},
		},
		{
			name:       "成功计算字数（过滤Markdown）",
			documentID: "doc123",
			requestBody: writer.WordCountRequest{
				Content:        "这是一段测试文本，包含一些**粗体**和*斜体*内容。\n\n这是第二段。",
				FilterMarkdown: true,
			},
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "计算成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.NotZero(t, data["totalCount"])
				assert.NotZero(t, data["chineseCount"])
				assert.NotZero(t, data["paragraphCount"])
			},
		},
		{
			name:       "空内容返回零值",
			documentID: "doc123",
			requestBody: writer.WordCountRequest{
				Content: "",
			},
			userID:         "user123",
			expectedStatus: http.StatusOK, // API允许空内容，返回零值
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "计算成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(0), data["totalCount"])
			},
		},
		{
			name:       "计算空格和特殊字符",
			documentID: "doc123",
			requestBody: writer.WordCountRequest{
				Content:        "   测试   空格   处理   ",
				FilterMarkdown: false,
			},
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				data := resp["data"].(map[string]interface{})
				assert.NotZero(t, data["totalCount"])
				assert.Equal(t, float64(6), data["chineseCount"]) // "测试空格处理"共6个汉字
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := setupSimpleEditorRouter(tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/documents/"+tt.documentID+"/word-count", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// TestEditorApi_GetUserShortcuts 测试获取用户快捷键配置
func TestEditorApi_GetUserShortcuts(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "成功获取快捷键配置",
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.NotNil(t, data)
				// 验证默认快捷键存在
				assert.NotNil(t, data["shortcuts"])
			},
		},
		{
			name:           "未认证用户",
			userID:         "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未授权", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := setupSimpleEditorRouter(tt.userID)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/shortcuts", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// TestEditorApi_UpdateUserShortcuts 测试更新用户快捷键配置
func TestEditorApi_UpdateUserShortcuts(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    writer.UpdateShortcutsRequest
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "成功更新快捷键配置",
			requestBody: writer.UpdateShortcutsRequest{
				Shortcuts: map[string]documentModel.Shortcut{
					"save": {
						Key:         "Ctrl+S",
						Description: "保存文档",
					},
					"undo": {
						Key:         "Ctrl+Z",
						Description: "撤销",
					},
				},
			},
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "更新成功", resp["message"])
			},
		},
		{
			name: "未认证用户",
			requestBody: writer.UpdateShortcutsRequest{
				Shortcuts: map[string]documentModel.Shortcut{},
			},
			userID:         "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未授权", resp["message"])
			},
		},
		{
			name: "空的快捷键配置",
			requestBody: writer.UpdateShortcutsRequest{
				Shortcuts: map[string]documentModel.Shortcut{},
			},
			userID:         "user123",
			expectedStatus: http.StatusInternalServerError, // 空配置导致ShortcutService错误
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
				assert.Equal(t, "更新快捷键配置失败", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := setupSimpleEditorRouter(tt.userID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/user/shortcuts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// TestEditorApi_ResetUserShortcuts 测试重置用户快捷键配置
func TestEditorApi_ResetUserShortcuts(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "成功重置快捷键配置",
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "重置成功", resp["message"])
			},
		},
		{
			name:           "未认证用户",
			userID:         "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未授权", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := setupSimpleEditorRouter(tt.userID)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/shortcuts/reset", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// TestEditorApi_GetShortcutHelp 测试获取快捷键帮助
func TestEditorApi_GetShortcutHelp(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "成功获取快捷键帮助",
			userID:         "user123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取成功", resp["message"])
				data := resp["data"]
				assert.NotNil(t, data)
				// 验证帮助信息包含分类
				helpData, ok := data.([]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, helpData)
			},
		},
		{
			name:           "未认证用户",
			userID:         "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未授权", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := setupSimpleEditorRouter(tt.userID)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/shortcuts/help", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// Note: MockDocumentRepository, MockProjectRepository和MockEventBus定义在test_helpers.go和document_api_test.go中
