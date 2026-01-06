package social

import (
	"context"

	"Qingyu_backend/models/social"
)

// MessageRepository 私信仓储接口
type MessageRepository interface {
	// ========== 会话管理 ==========

	// CreateConversation 创建会话
	CreateConversation(ctx context.Context, conversation *social.Conversation) error

	// GetConversationByID 根据ID获取会话
	GetConversationByID(ctx context.Context, conversationID string) (*social.Conversation, error)

	// GetConversationByParticipants 根据参与者获取会话（私聊）
	GetConversationByParticipants(ctx context.Context, participantIDs []string) (*social.Conversation, error)

	// GetUserConversations 获取用户的会话列表
	GetUserConversations(ctx context.Context, userID string, page, size int) ([]*social.Conversation, int64, error)

	// UpdateConversation 更新会话
	UpdateConversation(ctx context.Context, conversationID string, updates map[string]interface{}) error

	// DeleteConversation 删除会话
	DeleteConversation(ctx context.Context, conversationID string) error

	// UpdateLastMessage 更新会话最后一条消息
	UpdateLastMessage(ctx context.Context, conversationID string, message *social.Message) error

	// IncrementUnreadCount 增加未读数
	IncrementUnreadCount(ctx context.Context, conversationID, userID string) error

	// ClearUnreadCount 清空未读数
	ClearUnreadCount(ctx context.Context, conversationID, userID string) error

	// ========== 消息管理 ==========

	// CreateMessage 创建消息
	CreateMessage(ctx context.Context, message *social.Message) error

	// GetMessageByID 根据ID获取消息
	GetMessageByID(ctx context.Context, messageID string) (*social.Message, error)

	// GetMessagesByConversation 获取会话的消息列表
	GetMessagesByConversation(ctx context.Context, conversationID string, page, size int) ([]*social.Message, int64, error)

	// GetMessagesBetweenUsers 获取两个用户之间的消息
	GetMessagesBetweenUsers(ctx context.Context, userID1, userID2 string, page, size int) ([]*social.Message, int64, error)

	// MarkMessageAsRead 标记消息为已读
	MarkMessageAsRead(ctx context.Context, messageID string) error

	// MarkConversationMessagesAsRead 标记会话中所有消息为已读
	MarkConversationMessagesAsRead(ctx context.Context, conversationID, userID string) error

	// DeleteMessage 删除消息（软删除）
	DeleteMessage(ctx context.Context, messageID string) error

	// CountUnreadMessages 统计未读消息数
	CountUnreadMessages(ctx context.Context, userID string) (int, error)

	// ========== @提醒 ==========

	// CreateMention 创建@提醒
	CreateMention(ctx context.Context, mention *social.Mention) error

	// GetMentionByID 根据ID获取@提醒
	GetMentionByID(ctx context.Context, mentionID string) (*social.Mention, error)

	// GetUserMentions 获取用户的@提醒列表
	GetUserMentions(ctx context.Context, userID string, page, size int) ([]*social.Mention, int64, error)

	// GetUnreadMentions 获取未读的@提醒
	GetUnreadMentions(ctx context.Context, userID string, page, size int) ([]*social.Mention, int64, error)

	// MarkMentionAsRead 标记@提醒为已读
	MarkMentionAsRead(ctx context.Context, mentionID string) error

	// MarkAllMentionsAsRead 标记所有@提醒为已读
	MarkAllMentionsAsRead(ctx context.Context, userID string) error

	// CountUnreadMentions 统计未读@提醒数
	CountUnreadMentions(ctx context.Context, userID string) (int, error)

	// DeleteMention 删除@提醒
	DeleteMention(ctx context.Context, mentionID string) error

	// Health 健康检查
	Health(ctx context.Context) error
}
