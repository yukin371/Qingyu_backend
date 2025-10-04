package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	authModel "Qingyu_backend/models/shared/auth"
)

// SessionServiceImpl 会话服务实现
type SessionServiceImpl struct {
	cacheClient CacheClient // Redis客户端
	sessionTTL  time.Duration
}

// NewSessionService 创建会话服务
func NewSessionService(cacheClient CacheClient) SessionService {
	return &SessionServiceImpl{
		cacheClient: cacheClient,
		sessionTTL:  24 * time.Hour, // 默认24小时
	}
}

// ============ 会话管理 ============

// CreateSession 创建会话
func (s *SessionServiceImpl) CreateSession(ctx context.Context, userID string) (*Session, error) {
	// 1. 生成会话ID
	sessionID, err := s.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("生成会话ID失败: %w", err)
	}

	// 2. 创建会话
	now := time.Now()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}

	// 3. 存储到Redis
	key := s.getSessionKey(sessionID)
	value := fmt.Sprintf("%s|%d|%d", userID, session.CreatedAt.Unix(), session.ExpiresAt.Unix())
	if err := s.cacheClient.Set(ctx, key, value, s.sessionTTL); err != nil {
		return nil, fmt.Errorf("存储会话失败: %w", err)
	}

	return session, nil
}

// GetSession 获取会话
func (s *SessionServiceImpl) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	// 1. 从Redis获取
	key := s.getSessionKey(sessionID)
	value, err := s.cacheClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("会话不存在或已过期")
	}

	// 2. 解析会话数据
	var userID string
	var createdAt, expiresAt int64
	_, err = fmt.Sscanf(value, "%s|%d|%d", &userID, &createdAt, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("解析会话数据失败: %w", err)
	}

	// 3. 构建会话对象
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Unix(createdAt, 0),
		ExpiresAt: time.Unix(expiresAt, 0),
	}

	// 4. 检查是否过期
	if time.Now().After(session.ExpiresAt) {
		_ = s.DestroySession(ctx, sessionID)
		return nil, fmt.Errorf("会话已过期")
	}

	return session, nil
}

// DestroySession 销毁会话
func (s *SessionServiceImpl) DestroySession(ctx context.Context, sessionID string) error {
	key := s.getSessionKey(sessionID)
	if err := s.cacheClient.Delete(ctx, key); err != nil {
		return fmt.Errorf("销毁会话失败: %w", err)
	}
	return nil
}

// RefreshSession 刷新会话
func (s *SessionServiceImpl) RefreshSession(ctx context.Context, sessionID string) error {
	// 1. 获取现有会话
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("刷新会话失败: %w", err)
	}

	// 2. 更新过期时间
	session.ExpiresAt = time.Now().Add(s.sessionTTL)

	// 3. 重新存储
	key := s.getSessionKey(sessionID)
	value := fmt.Sprintf("%s|%d|%d", session.UserID, session.CreatedAt.Unix(), session.ExpiresAt.Unix())
	if err := s.cacheClient.Set(ctx, key, value, s.sessionTTL); err != nil {
		return fmt.Errorf("刷新会话失败: %w", err)
	}

	return nil
}

// ValidateSession 验证会话
func (s *SessionServiceImpl) ValidateSession(ctx context.Context, sessionID string) (bool, error) {
	_, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetUserSessions 获取用户的所有会话（可选实现）
func (s *SessionServiceImpl) GetUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	// TODO: 需要维护用户ID到会话ID的映射
	return nil, fmt.Errorf("功能待实现")
}

// DestroyUserSessions 销毁用户的所有会话（可选实现）
func (s *SessionServiceImpl) DestroyUserSessions(ctx context.Context, userID string) error {
	// TODO: 需要维护用户ID到会话ID的映射
	return fmt.Errorf("功能待实现")
}

// UpdateSession 更新会话（实现接口）
func (s *SessionServiceImpl) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	// 1. 获取现有会话
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("更新会话失败: %w", err)
	}

	// 2. 更新数据（这里简化处理，实际可扩展存储更多数据）
	key := s.getSessionKey(sessionID)
	value := fmt.Sprintf("%s|%d|%d", session.UserID, session.CreatedAt.Unix(), session.ExpiresAt.Unix())

	// 3. 重新存储
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("会话已过期")
	}

	if err := s.cacheClient.Set(ctx, key, value, ttl); err != nil {
		return fmt.Errorf("更新会话失败: %w", err)
	}

	return nil
}

// ============ 辅助方法 ============

// generateSessionID 生成会话ID
func (s *SessionServiceImpl) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// getSessionKey 获取会话Redis Key
func (s *SessionServiceImpl) getSessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

// UpdateSessionModel 更新Session模型（与AuthService集成）
func (s *SessionServiceImpl) UpdateSessionModel(session *Session) *authModel.Session {
	return &authModel.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Data:      make(map[string]interface{}),
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}
