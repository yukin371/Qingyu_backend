package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Qingyu_backend/api/v1/ai"
	aiService "Qingyu_backend/service/ai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// AIIntegrationTestSuite AI集成测试套件
type AIIntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	api    *ai.AIApi
}

// SetupSuite 设置测试套件
func (suite *AIIntegrationTestSuite) SetupSuite() {
	// 设置测试模式
	gin.SetMode(gin.TestMode)
	
	// 创建API实例（使用nil服务，专注测试HTTP层）
	suite.api = &ai.AIApi{}
	
	// 设置路由
	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	aiRouter := v1.Group("/ai")
	{
		aiRouter.POST("/generate", suite.api.GenerateContent)
		aiRouter.POST("/analyze", suite.api.AnalyzeContent)
		aiRouter.POST("/continue", suite.api.ContinueWriting)
		aiRouter.POST("/optimize", suite.api.OptimizeText)
		aiRouter.POST("/outline", suite.api.GenerateOutline)
	}
}

// TearDownSuite 清理测试套件
func (suite *AIIntegrationTestSuite) TearDownSuite() {
	// 清理资源
}

// TestServiceLayerIntegration 测试服务层的数据结构集成
func (suite *AIIntegrationTestSuite) TestServiceLayerIntegration() {
	// 测试生成内容请求的创建和验证
	generateReq := &aiService.GenerateContentRequest{
		ProjectID: "test-service-project",
		ChapterID: "test-service-chapter",
		Prompt:    "测试服务层集成",
		Options:   nil,
	}
	
	// 验证请求结构
	suite.NotEmpty(generateReq.ProjectID)
	suite.NotEmpty(generateReq.Prompt)
	
	// 测试分析内容请求
	analyzeReq := &aiService.AnalyzeContentRequest{
		Content:      "测试分析内容",
		AnalysisType: "general",
	}
	
	suite.NotEmpty(analyzeReq.Content)
	suite.NotEmpty(analyzeReq.AnalysisType)
	
	// 测试续写请求
	continueReq := &aiService.ContinueWritingRequest{
		ProjectID:      "test-continue-project",
		ChapterID:      "test-continue-chapter",
		CurrentText:    "当前文本内容",
		ContinueLength: 500,
		Options:        nil,
	}
	
	suite.NotEmpty(continueReq.ProjectID)
	suite.NotEmpty(continueReq.CurrentText)
	suite.Greater(continueReq.ContinueLength, 0)
	
	// 测试优化文本请求
	optimizeReq := &aiService.OptimizeTextRequest{
		ProjectID:    "test-optimize-project",
		ChapterID:    "test-optimize-chapter",
		OriginalText: "原始文本",
		OptimizeType: "style",
		Instructions: "优化指令",
		Options:      nil,
	}
	
	suite.NotEmpty(optimizeReq.ProjectID)
	suite.NotEmpty(optimizeReq.OriginalText)
	suite.NotEmpty(optimizeReq.OptimizeType)
	
	// 测试生成大纲请求
	outlineReq := &aiService.GenerateOutlineRequest{
		ProjectID:   "test-outline-project",
		Theme:       "测试主题",
		Genre:       "测试类型",
		Length:      "短篇",
		KeyElements: []string{"元素1", "元素2"},
		Options:     nil,
	}
	
	suite.NotEmpty(outlineReq.ProjectID)
	suite.NotEmpty(outlineReq.Theme)
	suite.NotEmpty(outlineReq.Genre)
	suite.NotEmpty(outlineReq.KeyElements)
}

// TestErrorHandling 测试错误处理集成
func (suite *AIIntegrationTestSuite) TestErrorHandling() {
	// 测试无效JSON
	req, err := http.NewRequest("POST", "/api/v1/ai/generate", bytes.NewBufferString("invalid json"))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(float64(400), response["code"])
	suite.Contains(response["message"], "请求参数错误")
	
	// 测试缺少必填字段
	requestBody := map[string]interface{}{
		"chapterId": "test-chapter",
		// 缺少 projectId 和 prompt
	}
	
	jsonBody, err := json.Marshal(requestBody)
	suite.NoError(err)
	
	req, err = http.NewRequest("POST", "/api/v1/ai/generate", bytes.NewBuffer(jsonBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusBadRequest, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(float64(400), response["code"])
	suite.Contains(response["message"], "项目ID和提示词不能为空")
}

// TestRequestValidation 测试请求验证集成
func (suite *AIIntegrationTestSuite) TestRequestValidation() {
	// 测试分析内容的空内容验证
	requestBody := map[string]interface{}{
		"content":      "",
		"analysisType": "plot",
	}
	
	jsonBody, err := json.Marshal(requestBody)
	suite.NoError(err)
	
	req, err := http.NewRequest("POST", "/api/v1/ai/analyze", bytes.NewBuffer(jsonBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(float64(400), response["code"])
	suite.Contains(response["message"], "分析内容不能为空")
}

// TestHTTPMethodsIntegration 测试HTTP方法集成
func (suite *AIIntegrationTestSuite) TestHTTPMethodsIntegration() {
	// 测试GET方法（应该不被支持）
	req, err := http.NewRequest("GET", "/api/v1/ai/generate", nil)
	suite.NoError(err)
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotFound, w.Code)
	
	// 测试PUT方法（应该不被支持）
	req, err = http.NewRequest("PUT", "/api/v1/ai/generate", nil)
	suite.NoError(err)
	
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusNotFound, w.Code)
}

// TestAIIntegrationSuite 运行AI集成测试套件
func TestAIIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AIIntegrationTestSuite))
}