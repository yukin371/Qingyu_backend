package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aiService "Qingyu_backend/service/ai"
	"Qingyu_backend/service/ai/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// setupTestRouter 设置测试路由
func setupTestRouter(api *WritingAssistantApi) *gin.Engine {
	router := gin.New()

	// 注册路由
	router.POST("/api/v1/ai/writing/summarize", api.SummarizeContent)
	router.POST("/api/v1/ai/writing/summarize-chapter", api.SummarizeChapter)
	router.POST("/api/v1/ai/writing/proofread", api.ProofreadContent)
	router.GET("/api/v1/ai/writing/suggestions/:id", api.GetProofreadSuggestion)
	router.POST("/api/v1/ai/audit/sensitive-words", api.CheckSensitiveWords)
	router.GET("/api/v1/ai/audit/sensitive-words/:id", api.GetSensitiveWordsDetail)

	return router
}

// TestWritingAssistantApi_SummarizeContent_Success 测试成功总结内容
func TestWritingAssistantApi_SummarizeContent_Success(t *testing.T) {
	// 创建模拟服务
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	// 准备请求
	reqBody := dto.SummarizeRequest{
		Content:     "这是需要总结的测试内容",
		SummaryType: "brief",
		MaxLength:   500,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 验证响应结构
	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// TestWritingAssistantApi_SummarizeContent_InvalidJSON 测试无效JSON
func TestWritingAssistantApi_SummarizeContent_InvalidJSON(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	// 无效JSON
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestWritingAssistantApi_SummarizeContent_EmptyContent 测试空内容
func TestWritingAssistantApi_SummarizeContent_EmptyContent(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.SummarizeRequest{
		Content: "",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 应该返回错误（因为内容验证在服务层）
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestWritingAssistantApi_SummarizeChapter_Success 测试章节总结
func TestWritingAssistantApi_SummarizeChapter_Success(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.ChapterSummaryRequest{
		ChapterID:    "chapter-123",
		ProjectID:    "project-456",
		OutlineLevel: 3,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize-chapter", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// TestWritingAssistantApi_SummarizeChapter_MissingChapterID 测试缺少章节ID
func TestWritingAssistantApi_SummarizeChapter_MissingChapterID(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.ChapterSummaryRequest{
		ChapterID: "",
		ProjectID: "project-456",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize-chapter", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestWritingAssistantApi_ProofreadContent_Success 测试成功校对
func TestWritingAssistantApi_ProofreadContent_Success(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.ProofreadRequest{
		Content: "这是需要校对的测试内容",
		CheckTypes: []string{"grammar", "spelling"},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/proofread", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// TestWritingAssistantApi_ProofreadContent_EmptyContent 测试空内容校对
func TestWritingAssistantApi_ProofreadContent_EmptyContent(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.ProofreadRequest{
		Content: "",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/proofread", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestWritingAssistantApi_GetProofreadSuggestion_Success 测试获取校对建议
func TestWritingAssistantApi_GetProofreadSuggestion_Success(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	suggestionID := "suggestion-123"
	req, _ := http.NewRequest("GET", "/api/v1/ai/writing/suggestions/"+suggestionID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// TestWritingAssistantApi_GetProofreadSuggestion_EmptyID 测试空建议ID
func TestWritingAssistantApi_GetProofreadSuggestion_EmptyID(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	req, _ := http.NewRequest("GET", "/api/v1/ai/writing/suggestions/", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 应该返回404（路由不匹配）
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestWritingAssistantApi_CheckSensitiveWords_Success 测试成功检测敏感词
func TestWritingAssistantApi_CheckSensitiveWords_Success(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.SensitiveWordsCheckRequest{
		Content:     "这是需要检测敏感词的测试内容",
		CustomWords: []string{"测试词"},
		Category:    "all",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/audit/sensitive-words", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")

	// 验证数据结构
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "checkId")
	assert.Contains(t, data, "isSafe")
	assert.Contains(t, data, "totalMatches")
}

// TestWritingAssistantApi_CheckSensitiveWords_EmptyContent 测试空内容
func TestWritingAssistantApi_CheckSensitiveWords_EmptyContent(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.SensitiveWordsCheckRequest{
		Content: "",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/audit/sensitive-words", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestWritingAssistantApi_GetSensitiveWordsDetail_Success 测试获取检测详情
func TestWritingAssistantApi_GetSensitiveWordsDetail_Success(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	checkID := "check-123"
	req, _ := http.NewRequest("GET", "/api/v1/ai/audit/sensitive-words/"+checkID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// TestWritingAssistantApi_GetSensitiveWordsDetail_EmptyID 测试空检测ID
func TestWritingAssistantApi_GetSensitiveWordsDetail_EmptyID(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	req, _ := http.NewRequest("GET", "/api/v1/ai/audit/sensitive-words/", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestWritingAssistantApi_ResponseHeaders 测试响应头
func TestWritingAssistantApi_ResponseHeaders(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.SummarizeRequest{
		Content:     "测试内容",
		SummaryType: "brief",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应头
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// TestWritingAssistantApi_HTTPMethods 测试HTTP方法
func TestWritingAssistantApi_HTTPMethods(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	t.Run("GET请求应该返回404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/ai/writing/summarize", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST请求应该正常处理", func(t *testing.T) {
		reqBody := dto.SummarizeRequest{
			Content:     "测试",
			SummaryType: "brief",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestWritingAssistantApi_MissingContentType 测试缺少Content-Type
func TestWritingAssistantApi_MissingContentType(t *testing.T) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.SummarizeRequest{
		Content:     "测试",
		SummaryType: "brief",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
	// 不设置Content-Type

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 应该仍然能处理（Gin会自动检测JSON）
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestWritingAssistantApi_Integration 测试API集成场景
func TestWritingAssistantApi_Integration(t *testing.T) {
	t.Run("完整的API请求流程", func(t *testing.T) {
		summarizeService := &aiService.SummarizeService{}
		proofreadService := &aiService.ProofreadService{}
		sensitiveWordsService := &aiService.SensitiveWordsService{}

		api := NewWritingAssistantApi(
			summarizeService,
			proofreadService,
			sensitiveWordsService,
		)

		router := setupTestRouter(api)

		// 测试1: 总结内容
		t.Run("总结内容API", func(t *testing.T) {
			reqBody := dto.SummarizeRequest{
				Content:       `这是一篇较长的文章内容。包含了很多文字和段落。
				需要AI来总结关键信息和要点。`,
				SummaryType:   "detailed",
				IncludeQuotes: true,
				MaxLength:     1000,
			}

			body, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, float64(200), response["code"])
			assert.Equal(t, "总结成功", response["message"])
		})

		// 测试2: 校对内容
		t.Run("校对内容API", func(t *testing.T) {
			reqBody := dto.ProofreadRequest{
				Content: "这是需要校对的内容。包含一些可能的语法错误。",
				CheckTypes: []string{"grammar", "spelling", "punctuation"},
			}

			body, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/writing/proofread", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, float64(200), response["code"])
			assert.Equal(t, "校对完成", response["message"])
		})

		// 测试3: 敏感词检测
		t.Run("敏感词检测API", func(t *testing.T) {
			reqBody := dto.SensitiveWordsCheckRequest{
				Content:     "这是需要检测敏感词的内容",
				CustomWords: []string{"测试词", "违规词"},
				Category:    "all",
			}

			body, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/audit/sensitive-words", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, float64(200), response["code"])
			assert.Equal(t, "检测完成", response["message"])
		})
	})
}

// BenchmarkWritingAssistantApi_SummarizeContent API性能测试
func BenchmarkWritingAssistantApi_SummarizeContent(b *testing.B) {
	summarizeService := &aiService.SummarizeService{}
	proofreadService := &aiService.ProofreadService{}
	sensitiveWordsService := &aiService.SensitiveWordsService{}

	api := NewWritingAssistantApi(
		summarizeService,
		proofreadService,
		sensitiveWordsService,
	)

	router := setupTestRouter(api)

	reqBody := dto.SummarizeRequest{
		Content:     "测试内容",
		SummaryType: "brief",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/ai/writing/summarize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
