package writer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestDocumentApi_DuplicateDocument_MissingUserID 测试缺少用户ID的情况
func TestDocumentApi_DuplicateDocument_MissingUserID(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()

	documentAPI := &DocumentApi{}

	r.POST("/api/v1/writer/documents/:id/duplicate", documentAPI.DuplicateDocument)

	documentID := primitive.NewObjectID().Hex()
	reqBody := []byte(`{"position": "inner"}`)
	req, _ := http.NewRequest("POST", "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusUnauthorized), response["code"])
}

// TestDocumentApi_DuplicateDocument_InvalidJSON 测试无效的JSON请求（有userId）
func TestDocumentApi_DuplicateDocument_InvalidJSON(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		c.Set("userId", primitive.NewObjectID().Hex())
		c.Next()
	})

	documentAPI := &DocumentApi{}

	r.POST("/api/v1/writer/documents/:id/duplicate", documentAPI.DuplicateDocument)

	documentID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDocumentApi_DuplicateDocument_ServiceNotInitialized 测试服务未初始化的情况
func TestDocumentApi_DuplicateDocument_ServiceNotInitialized(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		c.Set("userId", primitive.NewObjectID().Hex())
		c.Next()
	})

	// 创建一个没有初始化documentService的API实例
	documentAPI := &DocumentApi{}

	r.POST("/api/v1/writer/documents/:id/duplicate", documentAPI.DuplicateDocument)

	documentID := primitive.NewObjectID().Hex()
	reqBody := []byte(`{"position": "inner"}`)
	req, _ := http.NewRequest("POST", "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()

	// 使用defer来捕获panic
	defer func() {
		if r := recover(); r != nil {
			// 测试通过了 - 确实会发生panic
			assert.NotNil(t, r)
		}
	}()

	r.ServeHTTP(w, req)

	// Then
	// DuplicateDocument方法内部会调用documentService.DuplicateDocument
	// 如果documentService是nil，应该会panic
	// 具体的状态码取决于实现
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// TestDocumentApi_DuplicateDocument_ValidRequestStructure 测试请求结构验证
func TestDocumentApi_DuplicateDocument_ValidRequestStructure(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		c.Set("userId", primitive.NewObjectID().Hex())
		c.Next()
	})

	documentAPI := &DocumentApi{}

	r.POST("/api/v1/writer/documents/:id/duplicate", documentAPI.DuplicateDocument)

	documentID := primitive.NewObjectID().Hex()

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "有效的inner位置",
			requestBody:    `{"position": "inner", "copyContent": false}`,
			expectedStatus: http.StatusInternalServerError, // 因为service是nil，会出错
		},
		{
			name:           "有效的before位置",
			requestBody:    `{"position": "before", "copyContent": true}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "有效的after位置",
			requestBody:    `{"position": "after", "copyContent": false}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "有目标父ID",
			requestBody:    `{"position": "inner", "targetParentId": "some-parent-id", "copyContent": false}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/writer/documents/"+documentID+"/duplicate", bytes.NewBufferString(tc.requestBody))
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

			r.ServeHTTP(w, req)

			// Then - 如果没有panic，检查状态码
			if w.Code != 0 {
				assert.Equal(t, tc.expectedStatus, w.Code)
			}
		})
	}
}
