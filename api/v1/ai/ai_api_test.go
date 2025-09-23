package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter 设置测试路由，不使用实际服务
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 创建一个AIApi实例，但不初始化service（避免数据库连接）
	api := &AIApi{service: nil}
	
	v1 := router.Group("/api/v1")
	aiRouter := v1.Group("/ai")
	{
		aiRouter.POST("/generate", api.GenerateContent)
		aiRouter.POST("/analyze", api.AnalyzeContent)
		aiRouter.POST("/continue", api.ContinueWriting)
		aiRouter.POST("/optimize", api.OptimizeText)
		aiRouter.POST("/outline", api.GenerateOutline)
	}
	
	return router
}

// TestAIApi_GenerateContent_BadRequest 测试生成内容API的错误请求
func TestAIApi_GenerateContent_BadRequest(t *testing.T) {
	router := setupTestRouter()
	
	// 测试无效的JSON
	req, _ := http.NewRequest("POST", "/api/v1/ai/generate", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "请求参数错误")
}

// TestAIApi_GenerateContent_MissingFields 测试生成内容API缺少必填字段
func TestAIApi_GenerateContent_MissingFields(t *testing.T) {
	router := setupTestRouter()
	
	// 测试缺少必填字段
	reqData := map[string]interface{}{
		"chapterId": "test-chapter",
		// 缺少 projectId 和 prompt
	}
	
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "项目ID和提示词不能为空")
}

// TestAIApi_AnalyzeContent_BadRequest 测试分析内容API的错误请求
func TestAIApi_AnalyzeContent_BadRequest(t *testing.T) {
	router := setupTestRouter()
	
	// 测试空请求体
	req, _ := http.NewRequest("POST", "/api/v1/ai/analyze", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "分析内容不能为空")
}

// TestAIApi_ContinueWriting_MissingFields 测试续写API缺少必填字段
func TestAIApi_ContinueWriting_MissingFields(t *testing.T) {
	router := setupTestRouter()
	
	// 测试缺少必填字段
	reqData := map[string]interface{}{
		"projectId": "test-project",
		// 缺少其他必填字段
	}
	
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/continue", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
}

// TestAIApi_OptimizeText_MissingFields 测试优化文本API缺少必填字段
func TestAIApi_OptimizeText_MissingFields(t *testing.T) {
	router := setupTestRouter()
	
	// 测试缺少必填字段
	reqData := map[string]interface{}{
		"projectId": "test-project",
		// 缺少其他必填字段
	}
	
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/optimize", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
}

// TestAIApi_GenerateOutline_MissingFields 测试生成大纲API缺少必填字段
func TestAIApi_GenerateOutline_MissingFields(t *testing.T) {
	router := setupTestRouter()
	
	// 测试缺少必填字段
	reqData := map[string]interface{}{
		"projectId": "test-project",
		// 缺少其他必填字段
	}
	
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/outline", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
}

// TestAIApi_InvalidJSON 测试无效JSON请求
func TestAIApi_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	// 测试所有端点的无效JSON处理
	endpoints := []string{
		"/api/v1/ai/generate",
		"/api/v1/ai/analyze",
		"/api/v1/ai/continue",
		"/api/v1/ai/optimize",
		"/api/v1/ai/outline",
	}
	
	for _, endpoint := range endpoints {
		req, _ := http.NewRequest("POST", endpoint, bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code, "端点 %s 应该返回400错误", endpoint)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), response["code"])
	}
}

// TestAIApi_MethodNotAllowed 测试不支持的HTTP方法
func TestAIApi_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()
	
	// 测试GET方法（应该只支持POST）
	endpoints := []string{
		"/api/v1/ai/generate",
		"/api/v1/ai/analyze",
		"/api/v1/ai/continue",
		"/api/v1/ai/optimize",
		"/api/v1/ai/outline",
	}
	
	for _, endpoint := range endpoints {
		req, _ := http.NewRequest("GET", endpoint, nil)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code, "端点 %s 应该不支持GET方法", endpoint)
	}
}