// Package models 提供消息相关数据模型
package models

import "time"

// Message 私信消息
type Message struct {
	ID             string    `bson:"_id"`
	ConversationID string    `bson:"conversation_id"`
	SenderID       string    `bson:"sender_id"`
	ReceiverID     string    `bson:"receiver_id"`
	Content        string    `bson:"content"`
	MessageType    string    `bson:"message_type"`
	IsRead         bool      `bson:"is_read"`
	ReadAt         time.Time `bson:"read_at"`
	CreatedAt      time.Time `bson:"created_at"`
}

// Conversation 对话
type Conversation struct {
	ID            string    `bson:"_id"`
	Participants  []string  `bson:"participants"`
	Type          string    `bson:"type"`
	Title         string    `bson:"title"`
	LastMessageAt time.Time `bson:"last_message_at"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

// Announcement 公告
type Announcement struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Content     string    `bson:"content"`
	Type        string    `bson:"type"`
	Priority    int       `bson:"priority"`
	IsPublished bool      `bson:"is_published"`
	PublishedAt time.Time `bson:"published_at"`
	ExpireAt    time.Time `bson:"expire_at"`
	CreatedAt   time.Time `bson:"created_at"`
}
