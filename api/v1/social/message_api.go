package social

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/social/dto"
	modelsMessaging "Qingyu_backend/models/messaging"
	serviceMessaging "Qingyu_backend/service/messaging"
	"Qingyu_backend/realtime/websocket"
)

// MessageAPIV2 消息API处理器（新版本，基于会话）
type MessageAPIV2 struct {
	messageService       *serviceMessaging.MessageService
	conversationService  *serviceMessaging.ConversationService
	wsHub                *websocket.MessagingWSHub
}

// NewMessageAPIV2 创建消息API实例
func NewMessageAPIV2(
	messageService *serviceMessaging.MessageService,
	conversationService *serviceMessaging.ConversationService,
	wsHub *websocket.MessagingWSHub,
) *MessageAPIV2 {
	return &MessageAPIV2{
		messageService:      messageService,
		conversationService: conversationService,
		wsHub:               wsHub,
	}
}

// GetMessages 获取消息列表
// @Summary 获取消息列表
// @Description 分页获取会话的消息列表，支持向上/向下翻页
// @Tags Social Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversationId path string true "会话ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param before query string false "获取此消息之前的消息"
// @Param after query string false "获取此消息之后的消息"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse "参数错误"
// @Failure 403 {object} shared.APIResponse "无权访问"
// @Failure 404 {object} shared.APIResponse "会话不存在"
// @Router /api/v1/social/messages/conversations/{conversationId}/messages [get]
func (api *MessageAPIV2) GetMessages(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	conversationID := c.Param("conversationId")
	if conversationID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "会话ID不能为空")
		return
	}

	var req dto.GetMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}
	req = *req.GetDefaults() // 应用默认值

	// 验证用户是否是会话参与者
	conv, err := api.conversationService.Get(c.Request.Context(), conversationID)
	if err != nil {
		if err == modelsMessaging.ErrConversationNotFound {
			shared.NotFound(c, "会话不存在")
		} else {
			shared.InternalError(c, "获取失败", err)
		}
		return
	}

	if !conv.HasParticipant(userID) {
		shared.Forbidden(c, "无权访问")
		return
	}

	// 解析before/after参数
	var before, after *string
	if req.Before != "" {
		before = &req.Before
	}
	if req.After != "" {
		after = &req.After
	}

	messages, total, err := api.messageService.GetMessages(
		c.Request.Context(),
		conversationID,
		userID,
		req.Page,
		req.PageSize,
		before,
		after,
	)
	if err != nil {
		shared.InternalError(c, "获取失败", err)
		return
	}

	// 转换为响应格式
	messageItems := make([]dto.MessageItem, len(messages))
	for i, msg := range messages {
		messageItems[i] = dto.MessageItem{
			ID:             msg.ID.Hex(),
			ConversationID: msg.ConversationID,
			SenderID:       msg.SenderID,
			ReceiverID:     msg.ReceiverID,
			Content:        msg.Content,
			Type:           string(msg.Type),
			Attachments:    convertAttachments(msg.Extra),
			ReplyTo:        msg.ParentID,
			Read:           msg.IsRead,
			SentAt:         msg.CreatedAt,
		}
	}

	hasMore := len(messages) == req.PageSize

	shared.Success(c, http.StatusOK, "获取成功", dto.GetMessagesResponse{
		Messages: messageItems,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  hasMore,
	})
}

// SendMessage 发送消息
// @Summary 发送消息
// @Description 向会话发送新消息，支持文本、图片、文件等类型
// @Tags Social Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversationId path string true "会话ID"
// @Param request body object true "发送消息请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse "参数错误"
// @Failure 403 {object} shared.APIResponse "无权访问"
// @Failure 404 {object} shared.APIResponse "会话不存在"
// @Router /api/v1/social/messages/conversations/{conversationId}/messages [post]
func (api *MessageAPIV2) SendMessage(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	conversationID := c.Param("conversationId")
	if conversationID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "会话ID不能为空")
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 验证会话权限
	conv, err := api.conversationService.Get(c.Request.Context(), conversationID)
	if err != nil {
		if err == modelsMessaging.ErrConversationNotFound {
			shared.NotFound(c, "会话不存在")
		} else {
			shared.InternalError(c, "获取失败", err)
		}
		return
	}

	if !conv.HasParticipant(userID) {
		shared.Forbidden(c, "无权访问")
		return
	}

	// 解析ReplyTo
	var replyToID *string
	if req.ReplyTo != nil && *req.ReplyTo != "" {
		replyToID = req.ReplyTo
	}

	// 确定接收者
	var receiverID string
	for _, participantID := range conv.ParticipantIDs {
		if participantID != userID {
			receiverID = participantID
			break
		}
	}

	// 创建消息
	message := &modelsMessaging.DirectMessage{
		ConversationID: conversationID,
		Content:        req.Content,
		Type:           modelsMessaging.MessageType(req.Type),
		Status:         modelsMessaging.MessageStatusNormal,
		Extra:          convertAttachmentsDTO(req.Attachments),
	}
	// 设置嵌入的字段
	message.SenderID = userID
	message.ReceiverID = receiverID
	message.ParentID = replyToID

	// 保存消息
	savedMessage, err := api.messageService.Create(c.Request.Context(), message)
	if err != nil {
		shared.InternalError(c, "发送失败", err)
		return
	}

	// 通过WebSocket推送消息
	if api.wsHub != nil {
		// 转换为WebSocket消息格式
		wsMessage := map[string]interface{}{
			"id":             savedMessage.ID.Hex(),
			"conversationId": savedMessage.ConversationID,
			"senderId":       savedMessage.SenderID,
			"receiverId":     savedMessage.ReceiverID,
			"content":        savedMessage.Content,
			"type":           string(savedMessage.Type),
			"attachments":    savedMessage.Extra,
			"replyTo":        savedMessage.ParentID,
			"read":           savedMessage.IsRead,
			"createdAt":      savedMessage.CreatedAt,
		}
		api.wsHub.SendMessage(conversationID, wsMessage, userID)
	}

	// 获取接收者未读数
	unreadCount, _ := api.messageService.GetUnreadCount(
		c.Request.Context(),
		conversationID,
		receiverID,
	)

	shared.Success(c, http.StatusOK, "发送成功", dto.SendMessageResponse{
		MessageID:   savedMessage.ID.Hex(),
		SentAt:      savedMessage.CreatedAt,
		UnreadCount: unreadCount,
	})
}

// CreateConversation 创建会话
// @Summary 创建会话
// @Description 创建新的私信会话，支持一对一私聊
// @Tags Social Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "创建会话请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse "参数错误"
// @Router /api/v1/social/messages/conversations [post]
func (api *MessageAPIV2) CreateConversation(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 验证参与者
	if !contains(req.ParticipantIDs, userID) {
		shared.Error(c, http.StatusBadRequest, "参数错误", "必须包含当前用户")
		return
	}

	// 创建会话
	conv, err := api.conversationService.Create(
		c.Request.Context(),
		req.ParticipantIDs,
		userID,
	)
	if err != nil {
		shared.InternalError(c, "创建失败", err)
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", dto.CreateConversationResponse{
		ConversationID: conv.ID.Hex(),
		Participants:   conv.ParticipantIDs,
		CreatedAt:      conv.CreatedAt,
	})
}

// MarkConversationRead 标记会话已读
// @Summary 标记会话已读
// @Description 标记会话的所有消息为已读状态
// @Tags Social Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversationId path string true "会话ID"
// @Param request body object true "标记已读请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse "参数错误"
// @Router /api/v1/social/messages/conversations/{conversationId}/read [post]
func (api *MessageAPIV2) MarkConversationRead(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	conversationID := c.Param("conversationId")
	if conversationID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "会话ID不能为空")
		return
	}

	var req dto.MarkConversationReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	readAt := time.Unix(req.ReadAt, 0)

	affected, err := api.messageService.MarkConversationRead(
		c.Request.Context(),
		userID,
		conversationID,
		readAt,
	)
	if err != nil {
		shared.InternalError(c, "标记失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "标记成功", dto.MarkAsReadResponse{
		Success: true,
		Message: fmt.Sprintf("已标记%d条消息为已读", affected),
	})
}

// 辅助函数

// contains 检查字符串切片是否包含某元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// convertAttachments 将models中的Extra转换为dto.MessageAttachmentDTO
func convertAttachments(extra map[string]interface{}) []dto.MessageAttachmentDTO {
	if extra == nil {
		return nil
	}

	attachments, ok := extra["attachments"].([]dto.MessageAttachmentDTO)
	if !ok {
		return nil
	}
	return attachments
}

// convertAttachmentsDTO 将dto.MessageAttachmentDTO转换为models中的Extra
func convertAttachmentsDTO(dtoAttachments []dto.MessageAttachmentDTO) map[string]interface{} {
	if len(dtoAttachments) == 0 {
		return nil
	}

	return map[string]interface{}{
		"attachments": dtoAttachments,
	}
}

// convertObjectIDToString 将primitive.ObjectID转换为string指针
func convertObjectIDToString(oid *primitive.ObjectID) *string {
	if oid == nil {
		return nil
	}
	hex := oid.Hex()
	return &hex
}
