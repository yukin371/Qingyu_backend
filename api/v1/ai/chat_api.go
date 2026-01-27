package ai

import (
	"io"
	"net/http"
	"strconv"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
	aiDto "Qingyu_backend/service/ai/dto" // Imported for Swagger annotations
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChatApi AI聊天API
type ChatApi struct {
	chatService  *aiService.ChatService
	quotaService *aiService.QuotaService
}

// NewChatApi 创建AI聊天API实例
func NewChatApi(chatService *aiService.ChatService, quotaService *aiService.QuotaService) *ChatApi {
	return &ChatApi{
		chatService:  chatService,
		quotaService: quotaService,
	}
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
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/ai/chat [post]
func (api *ChatApi) Chat(c *gin.Context) {
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
func (api *ChatApi) ChatStream(c *gin.Context) {
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
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/ai/chat/sessions [get]
func (api *ChatApi) GetChatSessions(c *gin.Context) {
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
// @Description 获取指定会话的聊天历史，支持分页
// @Tags AI聊天
// @Accept json
// @Produce json
// @Param sessionId path string true "会话ID"
// @Param limit query int false "每页数量" default(50) minimum(1) maximum(100)
// @Param offset query int false "偏移量" default(0) minimum(0)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/ai/chat/sessions/:sessionId [get]
func (api *ChatApi) GetChatHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "会话ID不能为空")
		return
	}

	// 获取分页参数，设置默认值和最大值
	limit := 50  // 默认50条
	offset := 0  // 默认从头开始

	if l, ok := c.GetQuery("limit"); ok {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			if n > 100 {
				limit = 100  // 最大100条
			} else {
				limit = n
			}
		}
	}

	if o, ok := c.GetQuery("offset"); ok {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	session, err := api.chatService.GetChatHistory(c.Request.Context(), sessionID, limit, offset)
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
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/ai/chat/sessions/:sessionId [delete]
func (api *ChatApi) DeleteChatSession(c *gin.Context) {
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

var _ = aiDto.ChatSessionDTO{}
