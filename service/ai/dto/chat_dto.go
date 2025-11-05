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
type ChatSessionDTO struct {
	ID        string            `json:"id"`
	ProjectID string            `json:"projectId,omitempty"`
	Title     string            `json:"title"`
	Messages  []*ChatMessageDTO `json:"messages"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
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
func ToSessionDTO(session *ai.ChatSession) *ChatSessionDTO {
	if session == nil {
		return nil
	}

	messages := make([]*ChatMessageDTO, len(session.Messages))
	for i, msg := range session.Messages {
		messages[i] = ToMessageDTO(&msg)
	}

	return &ChatSessionDTO{
		ID:        session.SessionID,
		ProjectID: session.ProjectID,
		Title:     session.Title,
		Messages:  messages,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}
