package dto

import "time"

// GetMessagesRequest 获取消息列表请求
type GetMessagesRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Before   string `form:"before,omitempty"` // 获取此消息之前的消息（分页）
	After    string `form:"after,omitempty"`  // 获取此消息之后的消息（分页）
}

// GetDefaults 返回默认值
func (r *GetMessagesRequest) GetDefaults() *GetMessagesRequest {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
	return r
}

// GetMessagesResponse 获取消息列表响应
type GetMessagesResponse struct {
	Messages []MessageItem `json:"messages"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	HasMore  bool          `json:"has_more"`
}

// MessageItem 消息项
type MessageItem struct {
	ID             string                `json:"id"`
	ConversationID string                `json:"conversation_id"`
	SenderID       string                `json:"sender_id"`
	ReceiverID     string                `json:"receiver_id"`
	Content        string                `json:"content"`
	Type           string                `json:"type"`
	Attachments    []MessageAttachmentDTO `json:"attachments,omitempty"`
	ReplyTo        *string               `json:"reply_to,omitempty"`
	Read           bool                  `json:"read"`
	SentAt         time.Time             `json:"sent_at"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	Content     string                `json:"content" binding:"required,min=1,max=5000"`
	Type        string                `json:"type" binding:"required,oneof=text image file"`
	Attachments []MessageAttachmentDTO `json:"attachments,omitempty"`
	ReplyTo     *string               `json:"reply_to,omitempty"`
}

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
	MessageID    string    `json:"message_id"`
	SentAt       time.Time `json:"sent_at"`
	UnreadCount  int       `json:"unread_count"`
}

// CreateConversationRequest 创建会话请求
type CreateConversationRequest struct {
	ParticipantIDs []string `json:"participant_ids" binding:"required,min=2"`
}

// CreateConversationResponse 创建会话响应
type CreateConversationResponse struct {
	ConversationID string    `json:"conversation_id"`
	Participants   []string  `json:"participants"`
	CreatedAt      time.Time `json:"created_at"`
}

// MarkConversationReadRequest 标记会话已读请求
type MarkConversationReadRequest struct {
	ReadAt int64 `json:"read_at" binding:"required"`
}

// MarkAsReadResponse 标记已读响应
type MarkAsReadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// MessageAttachmentDTO 消息附件DTO（用于API层，区别于models层）
type MessageAttachmentDTO struct {
	Type      string `json:"type" binding:"required,oneof=image file"`
	URL       string `json:"url" binding:"required"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime_type"`
}
