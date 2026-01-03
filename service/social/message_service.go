package social

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
)

// MessageService 私信服务
type MessageService struct {
	messageRepo socialRepo.MessageRepository
	eventBus    base.EventBus
	serviceName string
	version     string
}

// NewMessageService 创建私信服务实例
func NewMessageService(
	messageRepo socialRepo.MessageRepository,
	eventBus base.EventBus,
) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		eventBus:    eventBus,
		serviceName: "MessageService",
		version:     "1.0.0",
	}
}

// =========================
// BaseService 接口实现
// =========================

func (s *MessageService) Initialize(ctx context.Context) error {
	return nil
}

func (s *MessageService) Health(ctx context.Context) error {
	if err := s.messageRepo.Health(ctx); err != nil {
		return fmt.Errorf("私信Repository健康检查失败: %w", err)
	}
	return nil
}

func (s *MessageService) Close(ctx context.Context) error {
	return nil
}

func (s *MessageService) GetServiceName() string {
	return s.serviceName
}

func (s *MessageService) GetVersion() string {
	return s.version
}

// =========================
// 会话管理
// =========================

// GetConversations 获取会话列表
func (s *MessageService) GetConversations(ctx context.Context, userID string, page, size int) ([]*social.Conversation, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	conversations, total, err := s.messageRepo.GetUserConversations(ctx, userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取会话列表失败: %w", err)
	}

	return conversations, total, nil
}

// GetConversationMessages 获取会话消息
func (s *MessageService) GetConversationMessages(ctx context.Context, userID, conversationID string, page, size int) ([]*social.Message, int64, error) {
	if userID == "" || conversationID == "" {
		return nil, 0, fmt.Errorf("用户ID和会话ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 50
	}

	// 获取会话验证权限
	conversation, err := s.messageRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, 0, fmt.Errorf("获取会话失败: %w", err)
	}

	// 检查用户是否是参与者
	isParticipant := false
	for _, participantID := range conversation.Participants {
		if participantID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, 0, fmt.Errorf("无权访问该会话")
	}

	messages, total, err := s.messageRepo.GetMessagesByConversation(ctx, conversationID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取消息失败: %w", err)
	}

	return messages, total, nil
}

// SendMessage 发送私信
func (s *MessageService) SendMessage(ctx context.Context, senderID, receiverID, content, messageType string) (*social.Message, error) {
	if senderID == "" || receiverID == "" {
		return nil, fmt.Errorf("发送者和接收者ID不能为空")
	}
	if senderID == receiverID {
		return nil, fmt.Errorf("不能给自己发送消息")
	}
	if content == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}
	if len(content) > 5000 {
		return nil, fmt.Errorf("消息内容最多5000字")
	}

	// 查找或创建会话
	participantIDs := []string{senderID, receiverID}
	conversation, err := s.messageRepo.GetConversationByParticipants(ctx, participantIDs)

	if err != nil || conversation == nil {
		// 创建新会话
		conversation = &social.Conversation{
			Participants: participantIDs,
			UnreadCount:  make(map[string]int),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.messageRepo.CreateConversation(ctx, conversation); err != nil {
			return nil, fmt.Errorf("创建会话失败: %w", err)
		}
	}

	// 创建消息
	message := &social.Message{
		ConversationID: conversation.ID.Hex(),
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Content:        content,
		MessageType:    messageType,
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	if err := s.messageRepo.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("发送消息失败: %w", err)
	}

	// 更新会话最后一条消息
	if err := s.messageRepo.UpdateLastMessage(ctx, conversation.ID.Hex(), message); err != nil {
		fmt.Printf("Warning: Failed to update last message: %v\n", err)
	}

	// 增加接收者未读数
	if err := s.messageRepo.IncrementUnreadCount(ctx, conversation.ID.Hex(), receiverID); err != nil {
		fmt.Printf("Warning: Failed to increment unread count: %v\n", err)
	}

	// 发布事件
	s.publishMessageEvent(ctx, "message.sent", senderID, receiverID, message.ID.Hex())

	return message, nil
}

// MarkMessageAsRead 标记消息已读
func (s *MessageService) MarkMessageAsRead(ctx context.Context, userID, messageID string) error {
	if userID == "" || messageID == "" {
		return fmt.Errorf("用户ID和消息ID不能为空")
	}

	// 获取消息
	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("获取消息失败: %w", err)
	}

	// 权限检查（只有接收者可以标记已读）
	if message.ReceiverID != userID {
		return fmt.Errorf("无权标记该消息")
	}

	if message.IsRead {
		return nil // 已读则跳过
	}

	// 标记已读
	if err := s.messageRepo.MarkMessageAsRead(ctx, messageID); err != nil {
		return fmt.Errorf("标记消息已读失败: %w", err)
	}

	// 减少未读数
	if err := s.messageRepo.ClearUnreadCount(ctx, message.ConversationID, userID); err != nil {
		fmt.Printf("Warning: Failed to clear unread count: %v\n", err)
	}

	return nil
}

// DeleteMessage 删除消息
func (s *MessageService) DeleteMessage(ctx context.Context, userID, messageID string) error {
	if userID == "" || messageID == "" {
		return fmt.Errorf("用户ID和消息ID不能为空")
	}

	// 获取消息
	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("获取消息失败: %w", err)
	}

	// 权限检查（只有发送者可以删除）
	if message.SenderID != userID {
		return fmt.Errorf("无权删除该消息")
	}

	// 软删除
	if err := s.messageRepo.DeleteMessage(ctx, messageID); err != nil {
		return fmt.Errorf("删除消息失败: %w", err)
	}

	// 发布事件
	s.publishMessageEvent(ctx, "message.deleted", message.SenderID, message.ReceiverID, messageID)

	return nil
}

// =========================
// @提醒
// =========================

// CreateMention 创建@提醒
func (s *MessageService) CreateMention(ctx context.Context, senderID, userID, contentType, contentID, content string) error {
	if senderID == "" || userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	mention := &social.Mention{
		UserID:      userID,
		SenderID:    senderID,
		ContentType: contentType,
		ContentID:   contentID,
		Content:     content,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := s.messageRepo.CreateMention(ctx, mention); err != nil {
		return fmt.Errorf("创建@提醒失败: %w", err)
	}

	// 发布事件
	s.publishMentionEvent(ctx, "mention.created", senderID, userID, mention.ID.Hex())

	return nil
}

// GetMentions 获取@提醒列表
func (s *MessageService) GetMentions(ctx context.Context, userID string, page, size int) ([]*social.Mention, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	mentions, total, err := s.messageRepo.GetUserMentions(ctx, userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取@提醒列表失败: %w", err)
	}

	return mentions, total, nil
}

// MarkMentionAsRead 标记@提醒已读
func (s *MessageService) MarkMentionAsRead(ctx context.Context, userID, mentionID string) error {
	if userID == "" || mentionID == "" {
		return fmt.Errorf("用户ID和提醒ID不能为空")
	}

	// 获取提醒
	mention, err := s.messageRepo.GetMentionByID(ctx, mentionID)
	if err != nil {
		return fmt.Errorf("获取@提醒失败: %w", err)
	}

	// 权限检查
	if mention.UserID != userID {
		return fmt.Errorf("无权标记该提醒")
	}

	// 标记已读
	if err := s.messageRepo.MarkMentionAsRead(ctx, mentionID); err != nil {
		return fmt.Errorf("标记@提醒已读失败: %w", err)
	}

	return nil
}

// =========================
// 私有辅助方法
// =========================

func (s *MessageService) publishMessageEvent(ctx context.Context, eventType, senderID, receiverID, messageID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"message_id":  messageID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

func (s *MessageService) publishMentionEvent(ctx context.Context, eventType, senderID, userID, mentionID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"sender_id":  senderID,
			"user_id":    userID,
			"mention_id": mentionID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}
