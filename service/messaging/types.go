package messaging

import (
	"errors"
	"time"
)

var (
	// ErrConversationNotFound 会话不存在
	ErrConversationNotFound = errors.New("会话不存在")

	// ErrNotParticipant 不是会话参与者
	ErrNotParticipant = errors.New("不是会话参与者")

	// ErrConversationExists 会话已存在
	ErrConversationExists = errors.New("会话已存在")

	// ErrInvalidParticipantIDs 无效的参与者ID列表
	ErrInvalidParticipantIDs = errors.New("参与者ID列表无效")

	// ErrMessageNotFound 消息不存在
	ErrMessageNotFound = errors.New("消息不存在")

	// ErrInvalidMessageType 无效的消息类型
	ErrInvalidMessageType = errors.New("无效的消息类型")

	// ErrInvalidReplyToID 无效的回复消息ID
	ErrInvalidReplyToID = errors.New("无效的回复消息ID")
)

// Attachment 附件
type Attachment struct {
	Type     string `bson:"type" json:"type"`         // image, file
	URL      string `bson:"url" json:"url"`           // 文件URL
	Size     int64  `bson:"size" json:"size"`         // 文件大小（字节）
	Name     string `bson:"name" json:"name"`         // 文件名
	MimeType string `bson:"mime_type" json:"mime_type"` // MIME类型
}

// LastMessage 最后消息信息
type LastMessage struct {
	Content string    `bson:"content" json:"content"`
	SentAt  time.Time `bson:"sent_at" json:"sent_at"`
	SentBy  string    `bson:"sent_by" json:"sent_by"`
}

// Message 消息模型（服务层使用）
type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	ReceiverID     string
	Content        string
	Type           string
	Read           bool
	ReadAt         *time.Time
	SentAt         time.Time
	CreatedAt      time.Time
	Attachments    []Attachment
	ReplyTo        *string
}

// Conversation 会话模型（服务层使用）
type Conversation struct {
	ID            string
	Participants  []string
	Type          string
	LastMessage   *LastMessage
	UnreadCount   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Muted         []string
}
