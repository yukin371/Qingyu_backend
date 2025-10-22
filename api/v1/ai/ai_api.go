package ai

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AIApi AI服务API
type AIApi struct {
	aiService    *aiService.Service
	chatService  *aiService.ChatService
	quotaService *aiService.QuotaService
}

// NewAIApi 创建AI API实例
func NewAIApi(aiService *aiService.Service, chatService *aiService.ChatService, quotaService *aiService.QuotaService) *AIApi {
	return &AIApi{
		aiService:    aiService,
		chatService:  chatService,
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
// @Success 200 {object} response.Response{data=aiService.GenerateContentResponse}
// @Router /api/v1/ai/continue [post]
func (api *AIApi) ContinueWriting(c *gin.Context) {
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
// @Router /api/v1/ai/continue/stream [post]
func (api *AIApi) ContinueWritingStream(c *gin.Context) {
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
// @Success 200 {object} response.Response{data=aiService.GenerateContentResponse}
// @Router /api/v1/ai/rewrite [post]
func (api *AIApi) RewriteText(c *gin.Context) {
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
// @Router /api/v1/ai/rewrite/stream [post]
func (api *AIApi) RewriteTextStream(c *gin.Context) {
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

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID  string              `json:"sessionId"`
	ProjectID  string              `json:"projectId"`
	Message    string              `json:"message" binding:"required"`
	UseContext bool                `json:"useContext"`
	Options    *ai.GenerateOptions `json:"options"`
}

// Chat 聊天
// @Summary AI聊天
// @Description 与AI助手进行对话
// @Tags AI聊天
// @Accept json
// @Produce json
// @Param request body ChatRequest true "聊天请求"
// @Success 200 {object} response.Response{data=aiService.ChatResponse}
// @Router /api/v1/ai/chat [post]
func (api *AIApi) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 转换为Service请求
	serviceReq := &aiService.ChatRequest{
		SessionID:  req.SessionID,
		ProjectID:  req.ProjectID,
		Message:    req.Message,
		UseContext: req.UseContext,
		Options:    req.Options,
	}

	// 调用聊天服务
	result, err := api.chatService.StartChat(c.Request.Context(), serviceReq)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "聊天失败", err.Error())
		return
	}

	// 设置Token使用信息
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)
	c.Set("requestID", uuid.New().String())
	c.Set("aiService", "chat")

	shared.Success(c, http.StatusOK, "聊天成功", result)
}

// ChatStream 聊天（流式）
// @Summary AI聊天（流式）
// @Description 与AI助手进行对话，流式返回结果
// @Tags AI聊天
// @Accept json
// @Produce text/event-stream
// @Param request body ChatRequest true "聊天请求"
// @Success 200 {string} string "SSE流"
// @Router /api/v1/ai/chat/stream [post]
func (api *AIApi) ChatStream(c *gin.Context) {
	var req ChatRequest
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

	// 转换为Service请求
	serviceReq := &aiService.ChatRequest{
		SessionID:  req.SessionID,
		ProjectID:  req.ProjectID,
		Message:    req.Message,
		UseContext: req.UseContext,
		Options:    req.Options,
	}

	// 获取流式响应通道
	streamChan, err := api.chatService.StartChatStream(c.Request.Context(), serviceReq)
	if err != nil {
		c.SSEvent("error", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 流式推送
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false

		case chunk, ok := <-streamChan:
			if !ok {
				return false
			}

			if chunk.IsComplete {
				c.SSEvent("done", gin.H{
					"sessionId":  chunk.SessionID,
					"messageId":  chunk.MessageID,
					"content":    chunk.Content,
					"tokensUsed": chunk.TokensUsed,
					"model":      chunk.Model,
				})

				// 异步消费配额
				go func() {
					userID, _ := c.Get("userId")
					_ = api.quotaService.ConsumeQuota(
						c.Request.Context(),
						userID.(string),
						chunk.TokensUsed,
						"chat",
						chunk.Model,
						chunk.MessageID,
					)
				}()

				return false
			}

			c.SSEvent("message", gin.H{
				"sessionId": chunk.SessionID,
				"messageId": chunk.MessageID,
				"delta":     chunk.Delta,
				"content":   chunk.Content,
				"tokens":    chunk.TokensUsed,
			})

			return true
		}
	})
}

// GetChatSessions 获取聊天会话列表
// @Summary 获取聊天会话列表
// @Description 获取用户的聊天会话列表
// @Tags AI聊天
// @Accept json
// @Produce json
// @Param projectId query string false "项目ID"
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.Response{data=[]dto.ChatSessionDTO}
// @Router /api/v1/ai/sessions [get]
func (api *AIApi) GetChatSessions(c *gin.Context) {
	projectID := c.Query("projectId")

	// 获取分页参数
	limit := 20
	offset := 0
	if l, ok := c.GetQuery("limit"); ok {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	if o, ok := c.GetQuery("offset"); ok {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	sessions, err := api.chatService.ListChatSessions(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取会话列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", sessions)
}

// GetChatHistory 获取聊天历史
// @Summary 获取聊天历史
// @Description 获取指定会话的聊天历史
// @Tags AI聊天
// @Accept json
// @Produce json
// @Param sessionId path string true "会话ID"
// @Success 200 {object} response.Response{data=dto.ChatSessionDTO}
// @Router /api/v1/ai/sessions/:sessionId [get]
func (api *AIApi) GetChatHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "会话ID不能为空")
		return
	}

	session, err := api.chatService.GetChatHistory(c.Request.Context(), sessionID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取聊天历史失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", session)
}

// DeleteChatSession 删除聊天会话
// @Summary 删除聊天会话
// @Description 删除指定的聊天会话
// @Tags AI聊天
// @Accept json
// @Produce json
// @Param sessionId path string true "会话ID"
// @Success 200 {object} response.Response
// @Router /api/v1/ai/sessions/:sessionId [delete]
func (api *AIApi) DeleteChatSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "会话ID不能为空")
		return
	}

	err := api.chatService.DeleteChatSession(c.Request.Context(), sessionID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除会话失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查AI服务状态
// @Tags AI系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/ai/health [get]
func (api *AIApi) HealthCheck(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "ai",
	}

	shared.Success(c, http.StatusOK, "服务正常", status)
}
