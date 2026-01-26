package messaging

import (
	"context"
	"time"

	"Qingyu_backend/models/messaging"
)

// ConversationService 会话服务
type ConversationService struct {
	conversationRepo messaging.ConversationRepository
	messageRepo      messaging.MessageRepository
}

// NewConversationService 创建会话服务
func NewConversationService(
	conversationRepo messaging.ConversationRepository,
	messageRepo messaging.MessageRepository,
) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

// Get 获取会话
func (s *ConversationService) Get(ctx context.Context, conversationID string) (*messaging.Conversation, error) {
	return s.conversationRepo.FindByID(ctx, conversationID)
}

// FindByParticipants 根据参与者查找会话
func (s *ConversationService) FindByParticipants(
	ctx context.Context,
	participantIDs []string,
) (*messaging.Conversation, error) {
	return s.conversationRepo.FindByParticipants(ctx, participantIDs)
}

// Create 创建会话
func (s *ConversationService) Create(
	ctx context.Context,
	participantIDs []string,
	currentUserID string,
) (*messaging.Conversation, error) {
	// 检查是否已存在相同参与者的会话
	existingConv, err := s.conversationRepo.FindByParticipants(ctx, participantIDs)
	if err == nil && existingConv != nil {
		return existingConv, nil
	}

	// 创建新会话
	conv := &messaging.Conversation{
		ParticipantIDs: participantIDs,
		Type:           messaging.ConversationTypeDirect,
		IsActive:       true,
		IsArchived:     false,
		CreatedBy:      currentUserID,
	}

	err = s.conversationRepo.Create(ctx, conv)
	if err != nil {
		return nil, err
	}

	return conv, nil
}

// UpdateLastMessage 更新会话的最后消息
func (s *ConversationService) UpdateLastMessage(
	ctx context.Context,
	conversationID string,
	messageID string,
	content string,
	senderID string,
	sentAt time.Time,
) error {
	conv, err := s.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return err
	}

	// 更新最后消息信息
	conv.UpdateLastMessage(messageID, content, senderID, sentAt)

	return s.conversationRepo.Update(ctx, conv)
}

// IncrementUnreadCount 增加未读计数
func (s *ConversationService) IncrementUnreadCount(
	ctx context.Context,
	conversationID string,
	userID string,
) error {
	return s.conversationRepo.IncrementUnreadCount(ctx, conversationID, userID)
}

// GetConversations 获取用户的会话列表
func (s *ConversationService) GetConversations(
	ctx context.Context,
	userID string,
	page, pageSize int,
) ([]*messaging.Conversation, int64, error) {
	return s.conversationRepo.FindByUserID(ctx, userID, page, pageSize)
}
