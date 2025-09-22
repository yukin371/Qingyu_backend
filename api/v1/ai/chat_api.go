package ai

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// ChatApi AI聊天API控制器
type ChatApi struct {
	chatService ChatServiceInterface
}

// ChatServiceInterface 聊天服务接口
type ChatServiceInterface interface {
	StartChat(ctx context.Context, req *aiService.ChatRequest) (*aiService.ChatResponse, error)
	GetChatHistory(ctx context.Context, sessionID string) (*aiService.ChatSession, error)
	ListChatSessions(ctx context.Context, projectID string, limit, offset int) ([]*aiService.ChatSession, error)
	DeleteChatSession(ctx context.Context, sessionID string) error
	GetStatistics(projectID string) (*aiService.ChatStatistics, error)
}

// NewChatApi 创建聊天API实例
func NewChatApi(chatService ChatServiceInterface) *ChatApi {
	return &ChatApi{
		chatService: chatService,
	}
}

// StartChat 开始聊天
// POST /api/v1/ai/chat
func (a *ChatApi) StartChat(c *gin.Context) {
	var req aiService.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "消息内容不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用服务开始聊天
	response, err := a.chatService.StartChat(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "开始聊天失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// GetChatHistory 获取聊天历史
// GET /api/v1/ai/chat/{sessionId}/history
func (a *ChatApi) GetChatHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "会话ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	session, err := a.chatService.GetChatHistory(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      404,
			"message":   "会话不存在: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      session,
		"timestamp": getTimestamp(),
	})
}

// ListChatSessions 列出聊天会话
// GET /api/v1/ai/chat/sessions?projectId={projectId}&limit={limit}&offset={offset}
func (a *ChatApi) ListChatSessions(c *gin.Context) {
	projectID := c.Query("projectId")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "limit参数格式错误",
			"timestamp": getTimestamp(),
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "offset参数格式错误",
			"timestamp": getTimestamp(),
		})
		return
	}

	sessions, err := a.chatService.ListChatSessions(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "获取会话列表失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"sessions": sessions,
			"total":    len(sessions),
			"limit":    limit,
			"offset":   offset,
		},
		"timestamp": getTimestamp(),
	})
}

// DeleteChatSession 删除聊天会话
// DELETE /api/v1/ai/chat/{sessionId}
func (a *ChatApi) DeleteChatSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "会话ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	err := a.chatService.DeleteChatSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "删除会话失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"timestamp": getTimestamp(),
	})
}

// ContinueChat 继续聊天（流式响应）
// POST /api/v1/ai/chat/stream
func (a *ChatApi) ContinueChat(c *gin.Context) {
	var req aiService.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "消息内容不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 设置流式响应
	if req.Options == nil {
		req.Options = &ai.GenerateOptions{}
	}
	req.Options.Stream = true

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 这里需要实现流式响应逻辑
	// 暂时使用普通响应
	response, err := a.chatService.StartChat(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "聊天处理失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// UpdateChatSession 更新聊天会话信息
// PUT /api/v1/ai/chat/{sessionId}
func (a *ChatApi) UpdateChatSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "会话ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 这里需要实现更新会话逻辑
	// 暂时返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"sessionId": sessionID,
			"title":     req.Title,
			"updatedAt": time.Now(),
		},
		"timestamp": getTimestamp(),
	})
}

// GetChatStatistics 获取聊天统计信息
// GET /api/v1/ai/chat/statistics?projectId={projectId}
func (a *ChatApi) GetChatStatistics(c *gin.Context) {
	projectID := c.Query("projectId")

	// 调用服务获取统计信息
	statistics, err := a.chatService.GetStatistics(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "获取统计信息失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      statistics,
		"timestamp": getTimestamp(),
	})
}

// ExportChatHistory 导出聊天历史
// GET /api/v1/ai/chat/{sessionId}/export?format={format}
func (a *ChatApi) ExportChatHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")
	format := c.DefaultQuery("format", "json")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "会话ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	session, err := a.chatService.GetChatHistory(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      404,
			"message":   "会话不存在: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	switch format {
	case "json":
		c.Header("Content-Disposition", "attachment; filename=chat_history.json")
		c.JSON(http.StatusOK, session)
	case "txt":
		c.Header("Content-Type", "text/plain")
		c.Header("Content-Disposition", "attachment; filename=chat_history.txt")

		// 转换为文本格式
		text := formatChatAsText(session)
		c.String(http.StatusOK, text)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "不支持的导出格式",
			"timestamp": getTimestamp(),
		})
	}
}

func formatChatAsText(session *aiService.ChatSession) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("聊天会话: %s\n", session.Title))
	result.WriteString(fmt.Sprintf("创建时间: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("更新时间: %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05")))
	result.WriteString("\n--- 消息记录 ---\n\n")
	
	for _, msg := range session.Messages {
		result.WriteString(fmt.Sprintf("[%s] %s: %s\n", 
			msg.Timestamp.Format("15:04:05"), 
			msg.Role, 
			msg.Content))
		result.WriteString("\n")
	}
	
	return result.String()
}