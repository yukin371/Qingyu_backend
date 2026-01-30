package messages

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/interfaces"
	"Qingyu_backend/pkg/response"
	"fmt"
)

// MessageAPI 私信API处理器
type MessageAPI struct {
	messageService interfaces.MessageService
}

// NewMessageAPI 创建私信API实例
func NewMessageAPI(messageService interfaces.MessageService) *MessageAPI {
	return &MessageAPI{
		messageService: messageService,
	}
}

// =========================
// 会话管理
// =========================

// GetConversations 获取会话列表
// @Summary 获取会话列表
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/messages/conversations [get]
// @Security Bearer
func (api *MessageAPI) GetConversations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	conversations, total, err := api.messageService.GetConversations(
		c.Request.Context(),
		userID.(string),
		params.Page,
		params.Size,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取会话列表成功", gin.H{
		"list":  conversations,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// GetConversationMessages 获取会话消息
// @Summary 获取会话消息
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param conversationId path string true "会话ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(50)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/messages/{conversationId} [get]
// @Security Bearer
func (api *MessageAPI) GetConversationMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		response.BadRequest(c,  "参数错误", "会话ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 50

	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	messages, total, err := api.messageService.GetConversationMessages(
		c.Request.Context(),
		userID.(string),
		conversationID,
		params.Page,
		params.Size,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取消息成功", gin.H{
		"list":  messages,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// =========================
// 消息管理
// =========================

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID  string `json:"receiver_id" binding:"required"`
	Content     string `json:"content" binding:"required,max=5000"`
	MessageType string `json:"message_type" binding:"required,oneof=text image system"`
}

// SendMessage 发送私信
// @Summary 发送私信
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param request body SendMessageRequest true "消息信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/messages [post]
// @Security Bearer
func (api *MessageAPI) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	senderID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	message, err := api.messageService.SendMessage(
		c.Request.Context(),
		senderID.(string),
		req.ReceiverID,
		req.Content,
		req.MessageType,
	)

	if err != nil {
		errMsg := err.Error()
		if errMsg == "不能给自己发送消息" {
			response.BadRequest(c,  "操作失败", errMsg)
		} else {
			response.InternalError(c, fmt.Errorf("发送消息失败: %s", errMsg))
		}
		return
	}

	response.Created(c, message)
}

// MarkMessageAsRead 标记消息已读
// @Summary 标记消息已读
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param id path string true "消息ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/messages/{id}/read [put]
// @Security Bearer
func (api *MessageAPI) MarkMessageAsRead(c *gin.Context) {
	messageID := c.Param("id")
	if messageID == "" {
		response.BadRequest(c,  "参数错误", "消息ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	err := api.messageService.MarkMessageAsRead(
		c.Request.Context(),
		userID.(string),
		messageID,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "标记消息已读成功", nil)
}

// DeleteMessage 删除消息
// @Summary 删除消息
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param id path string true "消息ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/messages/{id} [delete]
// @Security Bearer
func (api *MessageAPI) DeleteMessage(c *gin.Context) {
	messageID := c.Param("id")
	if messageID == "" {
		response.BadRequest(c,  "参数错误", "消息ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	err := api.messageService.DeleteMessage(
		c.Request.Context(),
		userID.(string),
		messageID,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除消息成功", nil)
}

// =========================
// @提醒
// =========================

// CreateMentionRequest 创建@提醒请求
type CreateMentionRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	ContentType string `json:"content_type" binding:"required,oneof=comment review message"`
	ContentID   string `json:"content_id" binding:"required"`
	Content     string `json:"content" binding:"required,max=200"`
}

// CreateMention 创建@提醒
// @Summary 创建@提醒
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param request body CreateMentionRequest true "提醒信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/mentions [post]
// @Security Bearer
func (api *MessageAPI) CreateMention(c *gin.Context) {
	var req CreateMentionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	senderID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	err := api.messageService.CreateMention(
		c.Request.Context(),
		senderID.(string),
		req.UserID,
		req.ContentType,
		req.ContentID,
		req.Content,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, nil)
}

// GetMentions 获取@提醒列表
// @Summary 获取@提醒列表
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/mentions [get]
// @Security Bearer
func (api *MessageAPI) GetMentions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	mentions, total, err := api.messageService.GetMentions(
		c.Request.Context(),
		userID.(string),
		params.Page,
		params.Size,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取@提醒列表成功", gin.H{
		"list":  mentions,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// MarkMentionAsRead 标记@提醒已读
// @Summary 标记@提醒已读
// @Tags 社交-私信
// @Accept json
// @Produce json
// @Param id path string true "提醒ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/mentions/{id}/read [put]
// @Security Bearer
func (api *MessageAPI) MarkMentionAsRead(c *gin.Context) {
	mentionID := c.Param("id")
	if mentionID == "" {
		response.BadRequest(c,  "参数错误", "提醒ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	err := api.messageService.MarkMentionAsRead(
		c.Request.Context(),
		userID.(string),
		mentionID,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "标记@提醒已读成功", nil)
}
