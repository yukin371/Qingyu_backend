package ai

import (
	"context"
	"fmt"

	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"Qingyu_backend/pkg/response"
)

// RAGApi RAG检索增强API
type RAGApi struct {
	ragService *aiService.RAGService
}

// NewRAGApi 创建RAG API实例
func NewRAGApi(ragService *aiService.RAGService) *RAGApi {
	return &RAGApi{
		ragService: ragService,
	}
}

// RetrieveAndGenerateRequest 检索增强生成请求
type RetrieveAndGenerateRequest struct {
	Query           string `json:"query" binding:"required"`
	TopK            int    `json:"topK"`
	ProjectID       string `json:"projectId"`
	IncludeOutline  bool   `json:"includeOutline"`
}

// RetrieveAndGenerateResponse 检索增强生成响应
type RetrieveAndGenerateResponse struct {
	Answer       string   `json:"answer"`
	Sources      []string `json:"sources"`
	References   []string `json:"references"`
	TokensUsed   int      `json:"tokensUsed"`
	ExecutionTime int64   `json:"executionTime"`
}

// RetrieveAndGenerate RAG检索增强生成
// @Summary RAG检索增强生成
// @Description 通过向量检索获取相关文档片段，结合LLM生成准确答案
// @Tags RAG检索
// @Accept json
// @Produce json
// @Param request body RetrieveAndGenerateRequest true "检索生成请求"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/ai/rag/retrieve [post]
func (api *RAGApi) RetrieveAndGenerate(c *gin.Context) {
	var req RetrieveAndGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "无法获取用户信息")
		return
	}

	ctx := c.Request.Context()

	// 调用RAG服务执行检索增强生成
	result, err := api.ragService.RetrieveAndGenerate(
		ctx,
		req.Query,
		userID.(string),
		req.ProjectID,
		req.TopK,
		req.IncludeOutline,
	)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 构造响应
	resp := &RetrieveAndGenerateResponse{
		Answer:       result.Answer,
		Sources:      result.Sources,
		References:   result.References,
		TokensUsed:   result.TokensUsed,
		ExecutionTime: result.ExecutionTime,
	}

	response.SuccessWithMessage(c, "检索增强生成成功", resp)
}

// SearchSimilar 搜索相似内容
// @Summary 搜索相似内容
// @Description 通过向量检索查找与查询内容相似的文档片段
// @Tags RAG检索
// @Accept json
// @Produce json
// @Param query query string true "查询内容"
// @Param topK topK int false "返回数量" default(5)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/ai/rag/search [get]
func (api *RAGApi) SearchSimilar(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		response.BadRequest(c, "查询内容不能为空", nil)
		return
	}

	topK := 5
	if topKStr := c.Query("topK"); topKStr != "" {
		if k, err := parseTopK(topKStr); err == nil {
			topK = k
		}
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "无法获取用户信息")
		return
	}

	ctx := c.Request.Context()

	// 调用RAG服务执行相似内容搜索
	sources, err := api.ragService.SearchSimilar(
		ctx,
		query,
		userID.(string),
		topK,
	)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, sources)
}

// parseTopK 解析topK参数
func parseTopK(s string) (int, error) {
	var k int
	_, err := fmt.Sscanf(s, "%d", &k)
	if err != nil {
		return 0, err
	}
	if k < 1 || k > 100 {
		return 0, fmt.Errorf("topK must be between 1 and 100")
	}
	return k, nil
}
