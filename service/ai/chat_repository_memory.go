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
	var totalTokens int

	for _, session := range r.sessions {
		if session.ProjectID == projectID {
			totalSessions++
			totalMessages += len(session.Messages)
			for _, msg := range session.Messages {
				totalTokens += msg.TokenUsed
			}
		}
	}

	return &ChatStatistics{
		TotalSessions: totalSessions,
		TotalMessages: totalMessages,
	}, nil
}
