package ai

import (
	"context"
	"fmt"
	"sync"

	"Qingyu_backend/models/ai"
)

// InMemoryChatRepository 内存Chat Repository实现（用于临时测试）
type InMemoryChatRepository struct {
	sessions map[string]*ai.ChatSession
	messages map[string][]*ai.ChatMessage
	mu       sync.RWMutex
}

// NewInMemoryChatRepository 创建内存Chat Repository
func NewInMemoryChatRepository() ChatRepositoryInterface {
	return &InMemoryChatRepository{
		sessions: make(map[string]*ai.ChatSession),
		messages: make(map[string][]*ai.ChatMessage),
	}
}

// CreateSession 创建会话
func (r *InMemoryChatRepository) CreateSession(ctx context.Context, session *ai.ChatSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.SessionID]; exists {
		return fmt.Errorf("会话已存在")
	}

	r.sessions[session.SessionID] = session
	return nil
}

// GetSessionByID 获取会话
func (r *InMemoryChatRepository) GetSessionByID(ctx context.Context, sessionID string) (*ai.ChatSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("会话不存在")
	}

	return session, nil
}

// UpdateSession 更新会话
func (r *InMemoryChatRepository) UpdateSession(ctx context.Context, session *ai.ChatSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.SessionID]; !exists {
		return fmt.Errorf("会话不存在")
	}

	r.sessions[session.SessionID] = session
	return nil
}

// DeleteSession 删除会话
func (r *InMemoryChatRepository) DeleteSession(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, sessionID)
	delete(r.messages, sessionID)
	return nil
}

// GetSessionsByProjectID 获取项目会话列表
func (r *InMemoryChatRepository) GetSessionsByProjectID(ctx context.Context, projectID string, limit, offset int) ([]*ai.ChatSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sessions []*ai.ChatSession
	for _, session := range r.sessions {
		if session.ProjectID == projectID {
			sessions = append(sessions, session)
		}
	}

	// 简单分页
	start := offset
	end := offset + limit

	if start >= len(sessions) {
		return []*ai.ChatSession{}, nil
	}

	if end > len(sessions) {
		end = len(sessions)
	}

	return sessions[start:end], nil
}

// CreateMessage 创建消息
func (r *InMemoryChatRepository) CreateMessage(ctx context.Context, message *ai.ChatMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[message.SessionID] = append(r.messages[message.SessionID], message)
	return nil
}

// GetSessionStatistics 获取会话统计信息
func (r *InMemoryChatRepository) GetSessionStatistics(ctx context.Context, projectID string) (*ChatStatistics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var totalSessions int
	var totalMessages int
	var activeSessions int

	for _, session := range r.sessions {
		if session.ProjectID == projectID {
			totalSessions++
			if session.Status == "active" {
				activeSessions++
			}
			// 从独立消息集合统计
			totalMessages += len(r.messages[session.SessionID])
		}
	}

	return &ChatStatistics{
		TotalSessions:  totalSessions,
		ActiveSessions: activeSessions,
		TotalMessages:  totalMessages,
	}, nil
}

// GetMessagesBySessionID 获取会话的所有消息
func (r *InMemoryChatRepository) GetMessagesBySessionID(ctx context.Context, sessionID string) ([]ai.ChatMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, exists := r.messages[sessionID]
	if !exists {
		return []ai.ChatMessage{}, nil
	}

	result := make([]ai.ChatMessage, len(msgs))
	for i, msg := range msgs {
		result[i] = *msg
	}

	return result, nil
}

// GetRecentMessagesBySessionID 获取会话的最近消息
func (r *InMemoryChatRepository) GetRecentMessagesBySessionID(ctx context.Context, sessionID string, limit int) ([]ai.ChatMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, exists := r.messages[sessionID]
	if !exists {
		return []ai.ChatMessage{}, nil
	}

	// 获取最近的 limit 条消息
	start := 0
	if len(msgs) > limit {
		start = len(msgs) - limit
	}

	result := make([]ai.ChatMessage, 0, limit)
	for i := start; i < len(msgs); i++ {
		result = append(result, *msgs[i])
	}

	return result, nil
}

// GetMessagesBySessionIDPaginated 分页获取消息
func (r *InMemoryChatRepository) GetMessagesBySessionIDPaginated(ctx context.Context, sessionID string, limit, offset int) ([]ai.ChatMessage, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, exists := r.messages[sessionID]
	if !exists {
		return []ai.ChatMessage{}, 0, nil
	}

	total := int64(len(msgs))

	// 分页
	start := offset
	if start >= len(msgs) {
		return []ai.ChatMessage{}, total, nil
	}

	end := offset + limit
	if end > len(msgs) {
		end = len(msgs)
	}

	result := make([]ai.ChatMessage, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, *msgs[i])
	}

	return result, total, nil
}

// GetMessageCountBySessionID 获取会话的消息总数
func (r *InMemoryChatRepository) GetMessageCountBySessionID(ctx context.Context, sessionID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, exists := r.messages[sessionID]
	if !exists {
		return 0, nil
	}

	return int64(len(msgs)), nil
}
