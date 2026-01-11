package interfaces

import (
	"context"
	socialModel "Qingyu_backend/models/social"
)

// MessageService 私信服务接口
type MessageService interface {
	// =========================
	// 会话管理
	// =========================
	GetConversations(ctx context.Context, userID string, page, size int) ([]*socialModel.Conversation, int64, error)
	GetConversationMessages(ctx context.Context, userID, conversationID string, page, size int) ([]*socialModel.Message, int64, error)
	SendMessage(ctx context.Context, senderID, receiverID, content, messageType string) (*socialModel.Message, error)

	// =========================
	// 消息管理
	// =========================
	MarkMessageAsRead(ctx context.Context, userID, messageID string) error
	DeleteMessage(ctx context.Context, userID, messageID string) error

	// =========================
	// @提醒
	// =========================
	CreateMention(ctx context.Context, senderID, userID, contentType, contentID, content string) error
	GetMentions(ctx context.Context, userID string, page, size int) ([]*socialModel.Mention, int64, error)
	MarkMentionAsRead(ctx context.Context, userID, mentionID string) error
}
