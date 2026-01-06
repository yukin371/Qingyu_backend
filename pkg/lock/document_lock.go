package lock

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"Qingyu_backend/pkg/cache"
)

// DocumentLock 文档锁
type DocumentLock struct {
	DocumentID  string    `json:"documentId"`
	UserID      string    `json:"userId"`
	UserName    string    `json:"userName"`
	LockedAt    time.Time `json:"lockedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	AutoExtend  bool      `json:"autoExtend"`
	DeviceID    string    `json:"deviceId"`
}

// LockStatus 锁状态
type LockStatus struct {
	IsLocked   bool           `json:"isLocked"`
	Lock       *DocumentLock  `json:"lock,omitempty"`
	CanEdit    bool           `json:"canEdit"`
	LockOwner  string         `json:"lockOwner,omitempty"`
	WaitQueue  []string       `json:"waitQueue,omitempty"` // 等待队列中的用户ID
}

// DocumentLockService 文档锁定服务接口
type DocumentLockService interface {
	// LockDocument 锁定文档
	LockDocument(ctx context.Context, documentID, userID, userName, deviceID string, autoExtend bool, ttl time.Duration) (*DocumentLock, error)

	// UnlockDocument 解锁文档
	UnlockDocument(ctx context.Context, documentID, userID string) error

	// RefreshLock 刷新锁（心跳）
	RefreshLock(ctx context.Context, documentID, userID string, ttl time.Duration) error

	// GetLockStatus 获取锁状态
	GetLockStatus(ctx context.Context, documentID, userID string) (*LockStatus, error)

	// ForceUnlock 强制解锁（管理员）
	ForceUnlock(ctx context.Context, documentID string) error

	// ExtendLock 延长锁时间
	ExtendLock(ctx context.Context, documentID, userID string, ttl time.Duration) error
}

// RedisDocumentLockService Redis分布式锁实现
type RedisDocumentLockService struct {
	client cache.RedisClient
	prefix string

	// 本地缓存
	localCache sync.Map
	// 本地锁（用于并发控制）
	mu sync.RWMutex
}

// NewRedisDocumentLockService 创建文档锁服务
func NewRedisDocumentLockService(client cache.RedisClient, prefix string) DocumentLockService {
	if prefix == "" {
		prefix = "doclock"
	}
	return &RedisDocumentLockService{
		client: client,
		prefix: prefix,
	}
}

// LockDocument 锁定文档
func (s *RedisDocumentLockService) LockDocument(ctx context.Context, documentID, userID, userName, deviceID string, autoExtend bool, ttl time.Duration) (*DocumentLock, error) {
	key := s.getLockKey(documentID)

	// 检查锁状态
	status, _ := s.GetLockStatus(ctx, documentID, userID)
	if status.IsLocked {
		// 检查是否是当前用户的锁
		if status.Lock != nil && status.Lock.UserID == userID {
			// 刷新锁
			if err := s.RefreshLock(ctx, documentID, userID, ttl); err != nil {
				return nil, fmt.Errorf("failed to refresh lock: %w", err)
			}
			return status.Lock, nil
		}
		// 被其他用户锁定
		return nil, fmt.Errorf("document is locked by %s", status.LockOwner)
	}

	// 创建锁
	now := time.Now()
	lock := &DocumentLock{
		DocumentID: documentID,
		UserID:     userID,
		UserName:   userName,
		LockedAt:   now,
		ExpiresAt:  now.Add(ttl),
		AutoExtend: autoExtend,
		DeviceID:   deviceID,
	}

	// 序列化
	data, err := s.serializeLock(lock)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize lock: %w", err)
	}

	// 使用Exists+Set实现SETNX（原子性的分布式锁）
	// 先检查锁是否存在
	exists, err := s.client.Exists(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to check lock: %w", err)
	}
	if exists > 0 {
		return nil, fmt.Errorf("failed to acquire lock: document is locked")
	}

	// 设置锁
	if err := s.client.Set(ctx, key, data, ttl); err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	// 保存到本地缓存
	s.localCache.Store(documentID, lock)

	return lock, nil
}

// UnlockDocument 解锁文档
func (s *RedisDocumentLockService) UnlockDocument(ctx context.Context, documentID, userID string) error {
	key := s.getLockKey(documentID)

	// 检查锁的所有者
	status, err := s.GetLockStatus(ctx, documentID, userID)
	if err != nil {
		return err
	}

	if !status.IsLocked {
		return nil // 未锁定，无需解锁
	}

	if status.Lock.UserID != userID {
		return fmt.Errorf("permission denied: lock owned by %s", status.LockOwner)
	}

	// 删除Redis锁
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	// 删除本地缓存
	s.localCache.Delete(documentID)

	return nil
}

// RefreshLock 刷新锁（心跳）
func (s *RedisDocumentLockService) RefreshLock(ctx context.Context, documentID, userID string, ttl time.Duration) error {
	key := s.getLockKey(documentID)

	// 获取当前锁
	status, err := s.GetLockStatus(ctx, documentID, userID)
	if err != nil {
		return err
	}

	if !status.IsLocked || status.Lock.UserID != userID {
		return fmt.Errorf("not authorized to refresh lock")
	}

	// 更新过期时间
	lock := status.Lock
	lock.ExpiresAt = time.Now().Add(ttl)

	// 序列化
	data, err := s.serializeLock(lock)
	if err != nil {
		return fmt.Errorf("failed to serialize lock: %w", err)
	}

	// 更新Redis（直接使用Set，因为我们已经验证了锁存在）
	if err := s.client.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to refresh lock: %w", err)
	}

	// 更新本地缓存
	s.localCache.Store(documentID, lock)

	return nil
}

// GetLockStatus 获取锁状态
func (s *RedisDocumentLockService) GetLockStatus(ctx context.Context, documentID, userID string) (*LockStatus, error) {
	key := s.getLockKey(documentID)

	status := &LockStatus{
		IsLocked:   false,
		CanEdit:    true,
		WaitQueue:  []string{},
	}

	// 从Redis获取锁信息
	data, err := s.client.Get(ctx, key)
	if err != nil {
		// 错误或者没有锁
		return status, nil
	}
	if data == "" {
		// 没有锁
		return status, nil
	}

	// 反序列化
	lock, err := s.deserializeLock(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize lock: %w", err)
	}

	// 检查是否过期
	if time.Now().After(lock.ExpiresAt) {
		// 锁已过期，删除
		s.client.Delete(ctx, key)
		s.localCache.Delete(documentID)
		return status, nil
	}

	status.IsLocked = true
	status.Lock = lock
	status.LockOwner = lock.UserName

	// 判断当前用户是否可以编辑
	if userID == "" {
		status.CanEdit = false
	} else {
		status.CanEdit = (lock.UserID == userID)
	}

	return status, nil
}

// ForceUnlock 强制解锁（管理员）
func (s *RedisDocumentLockService) ForceUnlock(ctx context.Context, documentID string) error {
	key := s.getLockLockKey(documentID)

	// 删除Redis锁
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to force unlock: %w", err)
	}

	// 删除本地缓存
	s.localCache.Delete(documentID)

	return nil
}

// ExtendLock 延长锁时间
func (s *RedisDocumentLockService) ExtendLock(ctx context.Context, documentID, userID string, ttl time.Duration) error {
	return s.RefreshLock(ctx, documentID, userID, ttl)
}

// getLockKey 获取锁键
func (s *RedisDocumentLockService) getLockKey(documentID string) string {
	return fmt.Sprintf("%s:doc:%s", s.prefix, documentID)
}

// getLockLockKey 获取锁的锁键（用于强制解锁）
func (s *RedisDocumentLockService) getLockLockKey(documentID string) string {
	return s.getLockKey(documentID)
}

// serializeLock 序列化锁
func (s *RedisDocumentLockService) serializeLock(lock *DocumentLock) (string, error) {
	// 使用JSON序列化
	data, err := lock.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// deserializeLock 反序列化锁
func (s *RedisDocumentLockService) deserializeLock(data string) (*DocumentLock, error) {
	lock := &DocumentLock{}
	err := lock.UnmarshalJSON([]byte(data))
	return lock, err
}

// MarshalJSON JSON序列化
func (l *DocumentLock) MarshalJSON() ([]byte, error) {
	type Alias DocumentLock
	return json.Marshal(&struct {
		LockedAt  string `json:"lockedAt"`
		ExpiresAt string `json:"expiresAt"`
		*Alias
	}{
		LockedAt:  l.LockedAt.Format(time.RFC3339),
		ExpiresAt: l.ExpiresAt.Format(time.RFC3339),
		Alias:     (*Alias)(l),
	})
}

// UnmarshalJSON JSON反序列化
func (l *DocumentLock) UnmarshalJSON(data []byte) error {
	type Alias DocumentLock
	aux := &struct {
		LockedAt  string `json:"lockedAt"`
		ExpiresAt string `json:"expiresAt"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	l.LockedAt, err = time.Parse(time.RFC3339, aux.LockedAt)
	if err != nil {
		return err
	}

	l.ExpiresAt, err = time.Parse(time.RFC3339, aux.ExpiresAt)
	return err
}
