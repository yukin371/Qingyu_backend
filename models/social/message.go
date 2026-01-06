package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/common"
)

// Conversation 会话
type Conversation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Participants []string           `bson:"participants" json:"participants"` // 参与者ID列表
	LastMessage  *Message           `bson:"last_message,omitempty" json:"last_message"`
	UnreadCount  map[string]int     `bson:"unread_count" json:"unread_count"` // 每个参与者的未读数
	common.BaseEntity                                                        // 嵌入时间戳字段
}

// Message 消息
type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversationID string             `bson:"conversation_id" json:"conversation_id"`
	SenderID       string             `bson:"sender_id" json:"sender_id"`
	ReceiverID     string             `bson:"receiver_id" json:"receiver_id"`
	Content        string             `bson:"content" json:"content"`
	MessageType    string             `bson:"message_type" json:"message_type"` // text, image, system
	common.ReadStatus                                                        // 嵌入已读状态
	IsDeleted      bool               `bson:"is_deleted" json:"is_deleted"`
	DeletedAt      *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	common.BaseEntity                                                        // 嵌入时间戳字段
}

// Mention @提醒
type Mention struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`           // 被@的用户ID
	SenderID    string             `bson:"sender_id" json:"sender_id"`       // @发起者ID
	ContentType string             `bson:"content_type" json:"content_type"` // 提及类型: comment, review, message
	ContentID   string             `bson:"content_id" json:"content_id"`     // 关联内容ID
	Content     string             `bson:"content" json:"content"`           // 提及内容摘要
	common.ReadStatus                                                        // 嵌入已读状态
	common.BaseEntity                                                        // 嵌入时间戳字段
}

// ConversationInfo 会话信息
type ConversationInfo struct {
	ID           primitive.ObjectID `json:"id"`
	Participants []*ParticipantInfo `json:"participants"`
	LastMessage  *MessageInfo       `json:"last_message,omitempty"`
	UnreadCount  int                `json:"unread_count"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// ParticipantInfo 参与者信息
type ParticipantInfo struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserAvatar string `json:"user_avatar,omitempty"`
	Online    bool   `json:"online"`
}

// MessageInfo 消息信息
type MessageInfo struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	SenderName string   `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar,omitempty"`
	Content   string    `json:"content"`
	MessageType string  `json:"message_type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
