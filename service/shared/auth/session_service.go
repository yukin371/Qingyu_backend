package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

	// 4. 将会话ID添加到用户会话列表（用于多端登录限制）
	// MVP: 使用JSON存储会话ID列表（简化实现，避免扩展CacheClient接口）
	if err := s.addSessionToUserList(ctx, userID, sessionID); err != nil {
		// 非关键错误，记录但不中断
		fmt.Printf("[Session] 添加会话到用户列表失败: %v\n", err)
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
	// 修复: 使用strings.Split替代fmt.Sscanf，因为%s会读取到空格而不是止于|
	parts := strings.Split(value, "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("会话数据格式无效")
	}

	userID := parts[0]
	createdAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("解析创建时间失败: %w", err)
	}
	expiresAt, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("解析过期时间失败: %w", err)
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
	// 1. 获取会话以确定userID
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		// 会话可能已过期，继续删除
	}

	// 2. 删除会话
	key := s.getSessionKey(sessionID)
	if err := s.cacheClient.Delete(ctx, key); err != nil {
		return fmt.Errorf("销毁会话失败: %w", err)
	}

	// 3. 从用户会话列表中移除
	if session != nil {
		_ = s.removeSessionFromUserList(ctx, session.UserID, sessionID)
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

// GetUserSessions 获取用户的所有会话（MVP实现）
func (s *SessionServiceImpl) GetUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	// 1. 从Redis获取用户的会话ID列表（JSON格式）
	userSessionsKey := s.getUserSessionsKey(userID)
	value, err := s.cacheClient.Get(ctx, userSessionsKey)
	if err != nil {
		// 列表不存在，返回空数组
		return []*Session{}, nil
	}

	// 2. 解析会话ID列表
	var sessionIDs []string
	if err := json.Unmarshal([]byte(value), &sessionIDs); err != nil {
		return nil, fmt.Errorf("解析用户会话列表失败: %w", err)
	}

	// 3. 逐个获取会话详情，过滤过期会话
	sessions := make([]*Session, 0, len(sessionIDs))
	validSessionIDs := make([]string, 0, len(sessionIDs))
	for _, sessionID := range sessionIDs {
		session, err := s.GetSession(ctx, sessionID)
		if err != nil {
			// 会话已过期，跳过
			continue
		}
		sessions = append(sessions, session)
		validSessionIDs = append(validSessionIDs, sessionID)
	}

	// 4. 更新会话列表（移除过期会话）
	if len(validSessionIDs) != len(sessionIDs) {
		_ = s.saveUserSessionList(ctx, userID, validSessionIDs)
	}

	return sessions, nil
}

// DestroyUserSessions 销毁用户的所有会话（MVP实现）
func (s *SessionServiceImpl) DestroyUserSessions(ctx context.Context, userID string) error {
	// 1. 获取用户的所有会话
	sessions, err := s.GetUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户会话失败: %w", err)
	}

	// 2. 逐个销毁会话
	for _, session := range sessions {
		_ = s.DestroySession(ctx, session.ID)
	}

	// 3. 删除用户会话列表
	userSessionsKey := s.getUserSessionsKey(userID)
	_ = s.cacheClient.Delete(ctx, userSessionsKey)

	return nil
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

// getUserSessionsKey 获取用户会话列表Redis Key
func (s *SessionServiceImpl) getUserSessionsKey(userID string) string {
	return fmt.Sprintf("user_sessions:%s", userID)
}

// CheckDeviceLimit 检查设备数量限制（MVP: 多端登录限制）
// maxDevices: 最大允许设备数，默认5
func (s *SessionServiceImpl) CheckDeviceLimit(ctx context.Context, userID string, maxDevices int) error {
	if maxDevices <= 0 {
		maxDevices = 5 // 默认最多5台设备
	}

	// 获取当前用户的所有活跃会话
	sessions, err := s.GetUserSessions(ctx, userID)
	if err != nil {
		// 查询失败时采用宽松策略，允许登录
		return nil
	}

	if len(sessions) >= maxDevices {
		return fmt.Errorf("登录设备数量已达上限（%d台），请退出其他设备后重试", maxDevices)
	}

	return nil
}

// addSessionToUserList 将会话ID添加到用户会话列表
func (s *SessionServiceImpl) addSessionToUserList(ctx context.Context, userID, sessionID string) error {
	userSessionsKey := s.getUserSessionsKey(userID)

	// 获取现有列表
	value, err := s.cacheClient.Get(ctx, userSessionsKey)
	var sessionIDs []string
	if err == nil {
		_ = json.Unmarshal([]byte(value), &sessionIDs)
	}

	// 添加新会话ID（去重）
	found := false
	for _, id := range sessionIDs {
		if id == sessionID {
			found = true
			break
		}
	}
	if !found {
		sessionIDs = append(sessionIDs, sessionID)
	}

	// 保存列表
	return s.saveUserSessionList(ctx, userID, sessionIDs)
}

// removeSessionFromUserList 从用户会话列表中移除会话ID
func (s *SessionServiceImpl) removeSessionFromUserList(ctx context.Context, userID, sessionID string) error {
	userSessionsKey := s.getUserSessionsKey(userID)

	// 获取现有列表
	value, err := s.cacheClient.Get(ctx, userSessionsKey)
	if err != nil {
		return nil // 列表不存在，无需操作
	}

	var sessionIDs []string
	if err := json.Unmarshal([]byte(value), &sessionIDs); err != nil {
		return err
	}

	// 移除会话ID
	newSessionIDs := make([]string, 0, len(sessionIDs))
	for _, id := range sessionIDs {
		if id != sessionID {
			newSessionIDs = append(newSessionIDs, id)
		}
	}

	// 保存列表
	return s.saveUserSessionList(ctx, userID, newSessionIDs)
}

// saveUserSessionList 保存用户会话列表
func (s *SessionServiceImpl) saveUserSessionList(ctx context.Context, userID string, sessionIDs []string) error {
	userSessionsKey := s.getUserSessionsKey(userID)

	if len(sessionIDs) == 0 {
		// 列表为空，删除key
		return s.cacheClient.Delete(ctx, userSessionsKey)
	}

	// 序列化为JSON
	data, err := json.Marshal(sessionIDs)
	if err != nil {
		return err
	}

	// 保存到Redis（过期时间比单个会话略长）
	return s.cacheClient.Set(ctx, userSessionsKey, string(data), s.sessionTTL+24*time.Hour)
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
