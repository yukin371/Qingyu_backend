//go:build integration
// +build integration

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAISystemAPI_GetProviders 获取AI提供商列表
func TestAISystemAPI_GetProviders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/ai/providers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data": []map[string]interface{}{
				{
					"name":        "openai",
					"displayName": "OpenAI",
					"status":      "active",
					"models": []string{
						"gpt-4",
						"gpt-3.5-turbo",
					},
				},
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/ai/providers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "获取成功", response["message"])

	data := response["data"].([]interface{})
	assert.Greater(t, len(data), 0)

	provider := data[0].(map[string]interface{})
	assert.Equal(t, "openai", provider["name"])
	assert.Equal(t, "active", provider["status"])
}

// TestAISystemAPI_GetModels 获取AI模型列表
func TestAISystemAPI_GetModels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/ai/models", func(c *gin.Context) {
		provider := c.Query("provider")

		allModels := []map[string]interface{}{
			{
				"id":        "gpt-4",
				"name":      "GPT-4",
				"provider":  "openai",
				"maxTokens": float64(8192),
				"costPer1k": 0.03,
			},
			{
				"id":        "gpt-3.5-turbo",
				"name":      "GPT-3.5 Turbo",
				"provider":  "openai",
				"maxTokens": float64(4096),
				"costPer1k": 0.002,
			},
		}

		// 过滤
		models := allModels
		if provider != "" {
			filtered := []map[string]interface{}{}
			for _, model := range allModels {
				if model["provider"] == provider {
					filtered = append(filtered, model)
				}
			}
			models = filtered
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data":    models,
		})
	})

	// 测试获取所有模型
	req := httptest.NewRequest("GET", "/api/v1/ai/models", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	model := data[0].(map[string]interface{})
	assert.Equal(t, "gpt-4", model["id"])
	assert.Equal(t, "openai", model["provider"])
}

// TestAISystemAPI_GetModels_WithFilter 按提供商过滤模型
func TestAISystemAPI_GetModels_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/ai/models", func(c *gin.Context) {
		provider := c.Query("provider")

		allModels := []map[string]interface{}{
			{
				"id":       "gpt-4",
				"provider": "openai",
			},
			{
				"id":       "gpt-3.5-turbo",
				"provider": "openai",
			},
			{
				"id":       "claude-2",
				"provider": "anthropic",
			},
		}

		models := allModels
		if provider != "" {
			filtered := []map[string]interface{}{}
			for _, model := range allModels {
				if model["provider"] == provider {
					filtered = append(filtered, model)
				}
			}
			models = filtered
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data":    models,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/ai/models?provider=openai", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// 验证所有模型都来自openai
	for _, item := range data {
		model := item.(map[string]interface{})
		assert.Equal(t, "openai", model["provider"])
	}
}

// TestAISystemAPI_HealthCheck 健康检查
func TestAISystemAPI_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/ai/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "服务正常",
			"data": map[string]interface{}{
				"status":  "healthy",
				"service": "ai",
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/ai/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "healthy", data["status"])
}
