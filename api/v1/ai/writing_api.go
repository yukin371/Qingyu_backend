package ai

import (
	"fmt"
	"io"
	"net/http"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WritingApi AI写作API
type WritingApi struct {
	aiService    *aiService.Service
	quotaService *aiService.QuotaService
}

// NewWritingApi 创建AI写作API实例
func NewWritingApi(aiService *aiService.Service, quotaService *aiService.QuotaService) *WritingApi {
	return &WritingApi{
		aiService:    aiService,
		quotaService: quotaService,
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

// ContinueWriting 智能续写
// @Summary 智能续写
// @Description 基于当前文本进行智能续写
// @Tags AI写作
// @Accept json
// @Produce json
// @Param request body ContinueWritingRequest true "续写请求"
// @Success 200 {object} shared.APIResponse{data=aiService.GenerateContentResponse}
// @Router /api/v1/ai/writing/continue [post]
func (api *WritingApi) ContinueWriting(c *gin.Context) {
	var req ContinueWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 转换为Service请求
	serviceReq := &aiService.ContinueWritingRequest{
		ProjectID:      req.ProjectID,
		ChapterID:      req.ChapterID,
		CurrentText:    req.CurrentText,
		ContinueLength: req.ContinueLength,
		Options:        req.Options,
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "continue_writing")

	// 调用服务
	result, err := api.aiService.ContinueWriting(c.Request.Context(), serviceReq)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "续写失败", err.Error())
		return
	}

	// 设置Token使用信息（用于配额消费）
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)

	shared.Success(c, http.StatusOK, "续写成功", result)
}

// ContinueWritingStream 智能续写（流式）
// @Summary 智能续写（流式）
// @Description 基于当前文本进行智能续写，流式返回结果
// @Tags AI写作
// @Accept json
// @Produce text/event-stream
// @Param request body ContinueWritingRequest true "续写请求"
// @Success 200 {string} string "SSE流"
// @Router /api/v1/ai/writing/continue/stream [post]
func (api *WritingApi) ContinueWritingStream(c *gin.Context) {
	var req ContinueWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
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

	// 转换为Service请求
	serviceReq := &aiService.GenerateContentRequest{
		ProjectID: req.ProjectID,
		ChapterID: req.ChapterID,
		Prompt:    fmt.Sprintf("请基于以下内容进行续写，保持风格和情节的连贯性：\n\n%s", req.CurrentText),
		Options:   req.Options,
	}

	if req.ContinueLength > 0 {
		serviceReq.Prompt += fmt.Sprintf("\n\n请续写约%d字的内容。", req.ContinueLength)
	}

	// 获取流式响应通道
	streamChan, err := api.aiService.GenerateContentStream(c.Request.Context(), serviceReq)
	if err != nil {
		c.SSEvent("error", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 流式推送
	var totalTokens int
	var model string
	fullContent := ""

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			// 客户端断开连接
			return false

		case chunk, ok := <-streamChan:
			if !ok {
				// channel关闭，发送完成事件
				c.SSEvent("done", gin.H{
					"requestId":  requestID,
					"content":    fullContent,
					"tokensUsed": totalTokens,
					"model":      model,
				})

				// 异步消费配额
				go func() {
					userID, _ := c.Get("userId")
					_ = api.quotaService.ConsumeQuota(
						c.Request.Context(),
						userID.(string),
						totalTokens,
						"continue_writing",
						model,
						requestID,
					)
				}()

				return false
			}

			// 累计内容和Token
			fullContent += chunk.Content
			totalTokens = chunk.TokensUsed
			model = chunk.Model

			// 发送增量数据
			c.SSEvent("message", gin.H{
				"requestId": requestID,
				"delta":     chunk.Content,
				"content":   fullContent,
				"tokens":    totalTokens,
			})

			return true
		}
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

// RewriteText 改写文本
// @Summary 改写文本
// @Description 对文本进行扩写、缩写或润色
// @Tags AI写作
// @Accept json
// @Produce json
// @Param request body RewriteTextRequest true "改写请求"
// @Success 200 {object} shared.APIResponse{data=aiService.GenerateContentResponse}
// @Router /api/v1/ai/writing/rewrite [post]
func (api *WritingApi) RewriteText(c *gin.Context) {
	var req RewriteTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 转换改写模式
	var optimizeType string
	switch req.RewriteMode {
	case "expand":
		optimizeType = "expand"
	case "shorten":
		optimizeType = "shorten"
	case "polish":
		optimizeType = "style"
	default:
		optimizeType = "style"
	}

	// 转换为Service请求
	serviceReq := &aiService.OptimizeTextRequest{
		ProjectID:    req.ProjectID,
		ChapterID:    req.ChapterID,
		OriginalText: req.OriginalText,
		OptimizeType: optimizeType,
		Instructions: req.Instructions,
		Options:      req.Options,
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "rewrite")

	// 调用服务
	result, err := api.aiService.OptimizeText(c.Request.Context(), serviceReq)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "改写失败", err.Error())
		return
	}

	// 设置Token使用信息
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)

	shared.Success(c, http.StatusOK, "改写成功", result)
}

// RewriteTextStream 改写文本（流式）
// @Summary 改写文本（流式）
// @Description 对文本进行扩写、缩写或润色，流式返回结果
// @Tags AI写作
// @Accept json
// @Produce text/event-stream
// @Param request body RewriteTextRequest true "改写请求"
// @Success 200 {string} string "SSE流"
// @Router /api/v1/ai/writing/rewrite/stream [post]
func (api *WritingApi) RewriteTextStream(c *gin.Context) {
	var req RewriteTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
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

	// 构建Prompt
	var prompt string
	switch req.RewriteMode {
	case "expand":
		prompt = "请对以下文本进行扩写，增加细节描述和情节内容："
	case "shorten":
		prompt = "请对以下文本进行缩写，保留核心内容："
	case "polish":
		prompt = "请对以下文本进行润色，优化表达方式："
	}

	if req.Instructions != "" {
		prompt += fmt.Sprintf("\n\n具体要求：%s", req.Instructions)
	}

	prompt += fmt.Sprintf("\n\n原文：\n%s", req.OriginalText)

	// 转换为Service请求
	serviceReq := &aiService.GenerateContentRequest{
		ProjectID: req.ProjectID,
		ChapterID: req.ChapterID,
		Prompt:    prompt,
		Options:   req.Options,
	}

	// 获取流式响应通道
	streamChan, err := api.aiService.GenerateContentStream(c.Request.Context(), serviceReq)
	if err != nil {
		c.SSEvent("error", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 流式推送
	var totalTokens int
	var model string
	fullContent := ""

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false

		case chunk, ok := <-streamChan:
			if !ok {
				c.SSEvent("done", gin.H{
					"requestId":  requestID,
					"content":    fullContent,
					"tokensUsed": totalTokens,
					"model":      model,
				})

				// 异步消费配额
				go func() {
					userID, _ := c.Get("userId")
					_ = api.quotaService.ConsumeQuota(
						c.Request.Context(),
						userID.(string),
						totalTokens,
						"rewrite",
						model,
						requestID,
					)
				}()

				return false
			}

			fullContent += chunk.Content
			totalTokens = chunk.TokensUsed
			model = chunk.Model

			c.SSEvent("message", gin.H{
				"requestId": requestID,
				"delta":     chunk.Content,
				"content":   fullContent,
				"tokens":    totalTokens,
			})

			return true
		}
	})
}
