package ai

import (
	"context"
	"fmt"
	"io"

	"Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WritingApi AI写作API
type WritingApi struct {
	textGenerator aiService.TextGenerator
	quotaService  *aiService.QuotaService
}

// NewWritingApi 创建AI写作API实例
func NewWritingApi(textGenerator aiService.TextGenerator, quotaService *aiService.QuotaService) *WritingApi {
	return &WritingApi{
		textGenerator: textGenerator,
		quotaService:  quotaService,
	}
}

// ContinueWritingRequest 续写请求
type ContinueWritingRequest struct {
	ProjectID      string              `json:"projectId" binding:"required"`
	ChapterID      string              `json:"chapterId"`
	CurrentText    string              `json:"currentText" binding:"required"`
	ContinueLength int                 `json:"continueLength"`
	Options        *ai.GenerateOptions `json:"options"`
}

// ContinueWritingRequestDoc 续写请求文档模型
type ContinueWritingRequestDoc struct {
	ProjectID      string              `json:"projectId"`
	ChapterID      string              `json:"chapterId,omitempty"`
	CurrentText    string              `json:"currentText"`
	ContinueLength int                 `json:"continueLength"`
	Options        *GenerateOptionsDoc `json:"options,omitempty"`
}

// ContinueWriting 智能续写
// @Summary 智能续写
// @Description 基于当前文本进行智能续写
// @Tags AI写作
// @Accept json
// @Produce json
// @Param request body ContinueWritingRequestDoc true "续写请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/writing/continue [post]
func (api *WritingApi) ContinueWriting(c *gin.Context) {
	var req ContinueWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "continue_writing")

	// 调用服务
	generateReq := api.buildContinueTextRequest(req)
	result, err := api.generateText(c.Request.Context(), generateReq)
	if err != nil {
		c.Error(err)
		return
	}

	// 设置Token使用信息（用于配额消费）
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)

	response.SuccessWithMessage(c, "续写成功", gin.H{
		"content":      result.Text,
		"tokensUsed":   result.TokensUsed,
		"model":        result.Model,
		"finishReason": "stop",
	})
}

// ContinueWritingStream 智能续写（流式）
// @Summary 智能续写（流式）
// @Description 基于当前文本进行智能续写，流式返回结果
// @Tags AI写作
// @Accept json
// @Produce text/event-stream
// @Param request body ContinueWritingRequestDoc true "续写请求"
// @Success 200 {string} string "SSE流"
// @Router /api/v1/ai/writing/continue/stream [post]
func (api *WritingApi) ContinueWritingStream(c *gin.Context) {
	var req ContinueWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Allow-Origin", "*")

	// 生成请求ID
	requestID := uuid.New().String()

	result, err := api.generateText(c.Request.Context(), api.buildContinueTextRequest(req))
	if err != nil {
		c.SSEvent("error", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Stream(func(w io.Writer) bool {
		c.SSEvent("message", gin.H{
			"requestId": requestID,
			"delta":     result.Text,
			"content":   result.Text,
			"tokens":    result.TokensUsed,
		})
		c.SSEvent("done", gin.H{
			"requestId":  requestID,
			"content":    result.Text,
			"tokensUsed": result.TokensUsed,
			"model":      result.Model,
		})
		api.consumeQuotaAsync(c, result.TokensUsed, "continue_writing", result.Model, requestID)
		return false
	})
}

// RewriteTextRequest 改写请求
type RewriteTextRequest struct {
	ProjectID    string              `json:"projectId"`
	ChapterID    string              `json:"chapterId"`
	OriginalText string              `json:"originalText" binding:"required"`
	RewriteMode  string              `json:"rewriteMode" binding:"required,oneof=expand shorten polish"`
	Instructions string              `json:"instructions"`
	Options      *ai.GenerateOptions `json:"options"`
}

// RewriteTextRequestDoc 改写请求文档模型
type RewriteTextRequestDoc struct {
	ProjectID    string              `json:"projectId,omitempty"`
	ChapterID    string              `json:"chapterId,omitempty"`
	OriginalText string              `json:"originalText"`
	RewriteMode  string              `json:"rewriteMode"`
	Instructions string              `json:"instructions,omitempty"`
	Options      *GenerateOptionsDoc `json:"options,omitempty"`
}

// RewriteText 改写文本
// @Summary 改写文本
// @Description 对文本进行扩写、缩写或润色
// @Tags AI写作
// @Accept json
// @Produce json
// @Param request body RewriteTextRequestDoc true "改写请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/writing/rewrite [post]
func (api *WritingApi) RewriteText(c *gin.Context) {
	var req RewriteTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 转换改写模式
	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "rewrite")

	// 调用服务
	result, err := api.generateText(c.Request.Context(), api.buildRewriteTextRequest(req))
	if err != nil {
		c.Error(err)
		return
	}

	// 设置Token使用信息
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)

	response.SuccessWithMessage(c, "改写成功", gin.H{
		"content":    result.Text,
		"tokensUsed": result.TokensUsed,
		"model":      result.Model,
		"changes":    []string{"文本已优化"},
	})
}

// RewriteTextStream 改写文本（流式）
// @Summary 改写文本（流式）
// @Description 对文本进行扩写、缩写或润色，流式返回结果
// @Tags AI写作
// @Accept json
// @Produce text/event-stream
// @Param request body RewriteTextRequestDoc true "改写请求"
// @Success 200 {string} string "SSE流"
// @Router /api/v1/ai/writing/rewrite/stream [post]
func (api *WritingApi) RewriteTextStream(c *gin.Context) {
	var req RewriteTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Allow-Origin", "*")

	// 生成请求ID
	requestID := uuid.New().String()

	result, err := api.generateText(c.Request.Context(), api.buildRewriteTextRequest(req))
	if err != nil {
		c.SSEvent("error", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Stream(func(w io.Writer) bool {
		c.SSEvent("message", gin.H{
			"requestId": requestID,
			"delta":     result.Text,
			"content":   result.Text,
			"tokens":    result.TokensUsed,
		})
		c.SSEvent("done", gin.H{
			"requestId":  requestID,
			"content":    result.Text,
			"tokensUsed": result.TokensUsed,
			"model":      result.Model,
		})
		api.consumeQuotaAsync(c, result.TokensUsed, "rewrite", result.Model, requestID)
		return false
	})
}

func (api *WritingApi) buildContinueTextRequest(req ContinueWritingRequest) *aiService.TextGenerateRequest {
	prompt := fmt.Sprintf("请基于以下内容进行续写，保持风格和情节的连贯性：\n\n%s", req.CurrentText)
	if req.ContinueLength > 0 {
		prompt += fmt.Sprintf("\n\n请续写约%d字的内容。", req.ContinueLength)
	}

	generateReq := &aiService.TextGenerateRequest{
		ProjectID:    req.ProjectID,
		ChapterID:    req.ChapterID,
		Prompt:       prompt,
		MaxTokens:    2000,
		Temperature:  0.7,
		WorkflowType: "continue_writing",
	}
	applyGenerateOptions(generateReq, req.Options)
	return generateReq
}

func (api *WritingApi) buildRewriteTextRequest(req RewriteTextRequest) *aiService.TextGenerateRequest {
	var prompt string
	switch req.RewriteMode {
	case "expand":
		prompt = "请对以下文本进行扩写，增加细节描述和情节内容："
	case "shorten":
		prompt = "请对以下文本进行缩写，保留核心内容："
	case "polish":
		prompt = "请对以下文本进行润色，优化表达方式："
	default:
		prompt = "请优化以下文本："
	}

	if req.Instructions != "" {
		prompt += fmt.Sprintf("\n\n具体要求：%s", req.Instructions)
	}
	prompt += fmt.Sprintf("\n\n原文：\n%s", req.OriginalText)

	generateReq := &aiService.TextGenerateRequest{
		ProjectID:    req.ProjectID,
		ChapterID:    req.ChapterID,
		Prompt:       prompt,
		MaxTokens:    2000,
		Temperature:  0.7,
		WorkflowType: "rewrite",
	}
	applyGenerateOptions(generateReq, req.Options)
	return generateReq
}

func applyGenerateOptions(req *aiService.TextGenerateRequest, options *ai.GenerateOptions) {
	if req == nil || options == nil {
		return
	}
	if options.Model != "" {
		req.Model = options.Model
	}
	if options.MaxTokens > 0 {
		req.MaxTokens = options.MaxTokens
	}
	if options.Temperature > 0 {
		req.Temperature = float64(options.Temperature)
	}
}

func (api *WritingApi) generateText(ctx context.Context, req *aiService.TextGenerateRequest) (*aiService.TextGenerateResponse, error) {
	if api == nil || api.textGenerator == nil {
		return nil, fmt.Errorf("AI text generator is not configured")
	}
	return api.textGenerator.GenerateText(ctx, req)
}

func (api *WritingApi) consumeQuotaAsync(c *gin.Context, tokensUsed int, serviceName, model, requestID string) {
	if api.quotaService == nil || tokensUsed <= 0 {
		return
	}
	userIDValue, ok := c.Get("user_id")
	if !ok {
		return
	}
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return
	}

	ctx := c.Request.Context()
	go func() {
		_ = api.quotaService.ConsumeQuota(ctx, userID, tokensUsed, serviceName, model, requestID)
	}()
}
