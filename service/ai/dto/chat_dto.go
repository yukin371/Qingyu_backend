package dto

import (
	"time"

	"Qingyu_backend/models/ai"
)

// ChatMessageDTO 聊天消息 DTO（API层使用）
// 用于API响应，提供更清晰的数据结构
type ChatMessageDTO struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"` // user, assistant, system
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ChatSessionDTO 聊天会话 DTO（API层使用）
// 用于API响应，只包含必要的字段
// 注意：不再包含 Messages 字段，消息需要单独查询
type ChatSessionDTO struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId,omitempty"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PaginatedMessagesDTO 分页消息 DTO
type PaginatedMessagesDTO struct {
	Messages []*ChatMessageDTO `json:"messages"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// ChatHistoryDTO 聊天历史 DTO
// 包含会话信息和分页消息
type ChatHistoryDTO struct {
	Session  *ChatSessionDTO      `json:"session"`
	Messages *PaginatedMessagesDTO `json:"messages"`
}

// ToMessageDTO 将模型转换为 DTO
// 转换 models/ai.ChatMessage 为 API 响应格式
func ToMessageDTO(msg *ai.ChatMessage) *ChatMessageDTO {
	if msg == nil {
		return nil
	}

	return &ChatMessageDTO{
		ID:        msg.ID.Hex(),
		Role:      msg.Role,
		Content:   msg.Content,
		Timestamp: msg.Timestamp,
		Metadata:  make(map[string]interface{}),
	}
}

// ToSessionDTO 将模型转换为 DTO
// 转换 models/ai.ChatSession 为 API 响应格式
// 注意：不再处理 Messages 字段
func ToSessionDTO(session *ai.ChatSession) *ChatSessionDTO {
	if session == nil {
		return nil
	}

	return &ChatSessionDTO{
		ID:        session.SessionID,
		ProjectID: session.ProjectID,
		Title:     session.Title,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}

// ToMessagesDTO 将消息列表转换为分页 DTO
func ToMessagesDTO(messages []ai.ChatMessage, total int64, limit, offset int) *PaginatedMessagesDTO {
	messageDTOs := make([]*ChatMessageDTO, len(messages))
	for i, msg := range messages {
		messageDTOs[i] = ToMessageDTO(&msg)
	}

	return &PaginatedMessagesDTO{
		Messages: messageDTOs,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}
}

// ToChatHistoryDTO 创建聊天历史 DTO
func ToChatHistoryDTO(session *ai.ChatSession, messages []ai.ChatMessage, total int64, limit, offset int) *ChatHistoryDTO {
	return &ChatHistoryDTO{
		Session:  ToSessionDTO(session),
		Messages: ToMessagesDTO(messages, total, limit, offset),
	}
}
