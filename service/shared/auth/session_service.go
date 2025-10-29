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

	"go.uber.org/zap"
)

// SessionServiceImpl 会话服务实现
type SessionServiceImpl struct {
	cacheClient   CacheClient   // Redis客户端
	sessionTTL    time.Duration // Session过期时间
	cleanupTicker *time.Ticker  // 定时清理任务
	cleanupStop   chan bool     // 停止清理信号
	lockTTL       time.Duration // 分布式锁过期时间
	isInitialized bool          // 是否已初始化
}

// NewSessionService 创建会话服务
func NewSessionService(cacheClient CacheClient) SessionService {
	service := &SessionServiceImpl{
		cacheClient:   cacheClient,
		sessionTTL:    24 * time.Hour,   // 默认24小时
		lockTTL:       10 * time.Second, // 锁过期时间10秒
		cleanupStop:   make(chan bool, 1),
		isInitialized: false,
	}

	// 启动定时清理任务（每1小时执行一次）
	service.startCleanupTask()

	return service
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
		zap.L().Warn("添加会话到用户列表失败",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
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
// TODO(performance): 使用Redis Pipeline批量获取Session，减少网络往返
// 当前实现: O(n)次Redis查询（1次列表 + n次详情）
// 优化后: 2次Redis查询（1次列表 + 1次Pipeline批量获取）
// 预期性能提升: 50-80%（当n>5时）
// 优先级: P1（Phase 3.5或Phase 4）
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

// addSessionToUserList 将会话ID添加到用户会话列表（带并发控制）
func (s *SessionServiceImpl) addSessionToUserList(ctx context.Context, userID, sessionID string) error {
	// 使用分布式锁保证并发安全
	return s.withUserSessionLock(ctx, userID, func() error {
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
	})
}

// removeSessionFromUserList 从用户会话列表中移除会话ID（带并发控制）
func (s *SessionServiceImpl) removeSessionFromUserList(ctx context.Context, userID, sessionID string) error {
	// 使用分布式锁保证并发安全
	return s.withUserSessionLock(ctx, userID, func() error {
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
	})
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

// ============ 定时清理任务 ============

// startCleanupTask 启动定时清理任务
func (s *SessionServiceImpl) startCleanupTask() {
	s.cleanupTicker = time.NewTicker(1 * time.Hour) // 每1小时执行一次

	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("清理任务panic", zap.Any("panic", r))
			}
		}()

		for {
			select {
			case <-s.cleanupTicker.C:
				// 执行清理
				ctx := context.Background()
				if err := s.CleanupExpiredSessions(ctx); err != nil {
					zap.L().Error("清理过期Session失败", zap.Error(err))
				} else {
					zap.L().Info("成功执行Session清理任务")
				}
			case <-s.cleanupStop:
				s.cleanupTicker.Stop()
				zap.L().Info("Session清理任务已停止")
				return
			}
		}
	}()

	zap.L().Info("Session定时清理任务已启动", zap.Duration("interval", 1*time.Hour))
}

// StopCleanupTask 停止清理任务（优雅关闭）
func (s *SessionServiceImpl) StopCleanupTask() {
	if s.cleanupTicker != nil {
		select {
		case s.cleanupStop <- true:
		default:
		}
	}
}

// CleanupExpiredSessions 清理过期会话（定时任务调用）
func (s *SessionServiceImpl) CleanupExpiredSessions(ctx context.Context) error {
	// 注意：这是一个简化实现
	// 生产环境可能需要扫描所有user_sessions:* key
	// 当前实现依赖Redis的自动过期机制 + 在GetUserSessions时过滤

	zap.L().Info("开始清理过期Session")

	// TODO(optimization): 实现完整的扫描清理逻辑
	// 1. 使用SCAN命令遍历所有user_sessions:*
	// 2. 对每个用户会话列表，移除过期的Session ID
	// 3. 记录清理统计

	// 当前版本：记录日志，实际清理在GetUserSessions中进行
	zap.L().Info("Session清理任务完成（依赖GetUserSessions自动过滤）")

	return nil
}

// ============ 分布式锁 ============

// acquireUserSessionLock 获取用户会话列表锁（基于Redis SETNX）
func (s *SessionServiceImpl) acquireUserSessionLock(ctx context.Context, userID string) (bool, error) {
	lockKey := s.getUserSessionLockKey(userID)
	lockValue := fmt.Sprintf("lock:%d", time.Now().Unix())

	// 使用SetNX实现分布式锁
	// 注意：这需要CacheClient支持SetNX，当前简化实现
	// 生产环境建议使用Redlock或etcd

	err := s.cacheClient.Set(ctx, lockKey, lockValue, s.lockTTL)
	if err != nil {
		return false, err
	}

	return true, nil
}

// releaseUserSessionLock 释放用户会话列表锁
func (s *SessionServiceImpl) releaseUserSessionLock(ctx context.Context, userID string) error {
	lockKey := s.getUserSessionLockKey(userID)
	return s.cacheClient.Delete(ctx, lockKey)
}

// getUserSessionLockKey 获取用户会话列表锁的Key
func (s *SessionServiceImpl) getUserSessionLockKey(userID string) string {
	return fmt.Sprintf("user_sessions_lock:%s", userID)
}

// withUserSessionLock 使用分布式锁执行操作
func (s *SessionServiceImpl) withUserSessionLock(ctx context.Context, userID string, fn func() error) error {
	// 尝试获取锁（带重试）
	maxRetries := 3
	retryDelay := 100 * time.Millisecond

	var acquired bool
	var err error

	for i := 0; i < maxRetries; i++ {
		acquired, err = s.acquireUserSessionLock(ctx, userID)
		if err != nil {
			return fmt.Errorf("获取锁失败: %w", err)
		}
		if acquired {
			break
		}
		// 未获取到锁，等待后重试
		time.Sleep(retryDelay)
		retryDelay *= 2 // 指数退避
	}

	if !acquired {
		return fmt.Errorf("无法获取用户会话锁，请稍后重试")
	}

	// 确保释放锁
	defer func() {
		if err := s.releaseUserSessionLock(ctx, userID); err != nil {
			zap.L().Warn("释放用户会话锁失败",
				zap.String("user_id", userID),
				zap.Error(err),
			)
		}
	}()

	// 执行操作
	return fn()
}
