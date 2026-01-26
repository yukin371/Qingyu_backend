package messaging

import (
	"context"
	"time"

	"Qingyu_backend/models/messaging"
)

// MessageService 消息服务
type MessageService struct {
	messageRepo         messaging.MessageRepository
	conversationService *ConversationService
}

// NewMessageService 创建消息服务
func NewMessageService(
	messageRepo messaging.MessageRepository,
	conversationService *ConversationService,
) *MessageService {
	return &MessageService{
		messageRepo:         messageRepo,
		conversationService: conversationService,
	}
}

// GetMessages 获取消息列表
func (s *MessageService) GetMessages(
	ctx context.Context,
	conversationID string,
	userID string,
	page, pageSize int,
	before, after *string,
) ([]*messaging.DirectMessage, int, error) {
	messages, total, err := s.messageRepo.FindByConversationID(
		ctx,
		conversationID,
		userID,
		page,
		pageSize,
		before,
		after,
	)

	return messages, total, err
}

// Create 创建消息
func (s *MessageService) Create(
	ctx context.Context,
	message *messaging.DirectMessage,
) (*messaging.DirectMessage, error) {
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	err := s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	// 更新会话的最后消息
	err = s.conversationService.UpdateLastMessage(
		ctx,
		message.ConversationID,
		message.ID.Hex(),
		message.Content,
		message.SenderID,
		message.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 增加未读计数
	err = s.conversationService.IncrementUnreadCount(
		ctx,
		message.ConversationID,
		message.ReceiverID,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// MarkConversationRead 标记会话已读
func (s *MessageService) MarkConversationRead(
	ctx context.Context,
	userID, conversationID string,
	readAt time.Time,
) (int, error) {
	// 标记会话的所有未读消息为已读
	return s.messageRepo.MarkConversationRead(
		ctx,
		conversationID,
		userID,
		readAt,
	)
}

// GetUnreadCount 获取未读消息数
func (s *MessageService) GetUnreadCount(
	ctx context.Context,
	conversationID string,
	userID string,
) (int, error) {
	return s.messageRepo.CountUnreadInConversation(ctx, conversationID, userID)
}
