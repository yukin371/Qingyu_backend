package ai

import (
	"net/http"
	"time"

	"Qingyu_backend/api/v1/shared"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// SystemApi AI系统API
type SystemApi struct {
	aiService *aiService.Service
}

// NewSystemApi 创建AI系统API实例
func NewSystemApi(aiService *aiService.Service) *SystemApi {
	return &SystemApi{
		aiService: aiService,
	}
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查AI服务状态
// @Tags AI系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/ai/health [get]
func (api *SystemApi) HealthCheck(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "ai",
	}

	shared.Success(c, http.StatusOK, "服务正常", status)
}

// GetProviders 获取AI提供商列表
// @Summary 获取AI提供商列表
// @Description 获取系统支持的AI提供商列表及其状态
// @Tags AI系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/ai/providers [get]
func (api *SystemApi) GetProviders(c *gin.Context) {
	// TODO: 实现获取提供商列表的逻辑
	providers := []gin.H{
		{
			"name":        "openai",
			"displayName": "OpenAI",
			"status":      "active",
			"models": []string{
				"gpt-4",
				"gpt-3.5-turbo",
			},
		},
	}

	shared.Success(c, http.StatusOK, "获取成功", providers)
}

// GetModels 获取可用模型列表
// @Summary 获取可用模型列表
// @Description 获取指定提供商的可用模型列表
// @Tags AI系统
// @Accept json
// @Produce json
// @Param provider query string false "提供商名称"
// @Success 200 {object} response.Response
// @Router /api/v1/ai/models [get]
func (api *SystemApi) GetModels(c *gin.Context) {
	provider := c.Query("provider")

	// TODO: 实现获取模型列表的逻辑
	models := []gin.H{
		{
			"id":        "gpt-4",
			"name":      "GPT-4",
			"provider":  "openai",
			"maxTokens": 8192,
			"costPer1k": 0.03,
		},
		{
			"id":        "gpt-3.5-turbo",
			"name":      "GPT-3.5 Turbo",
			"provider":  "openai",
			"maxTokens": 4096,
			"costPer1k": 0.002,
		},
	}

	// 按提供商过滤
	if provider != "" {
		filtered := []gin.H{}
		for _, model := range models {
			if model["provider"] == provider {
				filtered = append(filtered, model)
			}
		}
		models = filtered
	}

	shared.Success(c, http.StatusOK, "获取成功", models)
}
