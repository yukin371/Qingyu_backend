package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// 全局密码重置Token管理器实例
var (
	globalPasswordResetTokenManager *PasswordResetTokenManager
	passwordResetTokenManagerOnce  sync.Once
)

// GetGlobalPasswordResetTokenManager 获取全局密码重置Token管理器（单例）
func GetGlobalPasswordResetTokenManager() *PasswordResetTokenManager {
	passwordResetTokenManagerOnce.Do(func() {
		globalPasswordResetTokenManager = &PasswordResetTokenManager{
			tokens: make(map[string]*ResetTokenInfo),
		}
	})
	return globalPasswordResetTokenManager
}

// PasswordResetTokenManager 密码重置Token管理器
type PasswordResetTokenManager struct {
	tokens map[string]*ResetTokenInfo // email -> token info
	mu     sync.RWMutex
}

// ResetTokenInfo 重置Token信息
type ResetTokenInfo struct {
	Token     string
	Email     string
	ExpiresAt time.Time
	Used      bool
}

// NewPasswordResetTokenManager 创建密码重置Token管理器
func NewPasswordResetTokenManager() *PasswordResetTokenManager {
	return &PasswordResetTokenManager{
		tokens: make(map[string]*ResetTokenInfo),
	}
}

// GenerateToken 生成重置Token
func (m *PasswordResetTokenManager) GenerateToken(ctx context.Context, email string) (string, error) {
	// 生成32字节随机Token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("生成Token失败: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 存储Token信息（有效期1小时）
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens[email] = &ResetTokenInfo{
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	return token, nil
}

// ValidateToken 验证Token
func (m *PasswordResetTokenManager) ValidateToken(ctx context.Context, email, token string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tokenInfo, exists := m.tokens[email]
	if !exists {
		return fmt.Errorf("无效的重置Token")
	}

	if tokenInfo.Token != token {
		return fmt.Errorf("无效的重置Token")
	}

	if tokenInfo.Used {
		return fmt.Errorf("重置Token已使用")
	}

	if time.Now().After(tokenInfo.ExpiresAt) {
		return fmt.Errorf("重置Token已过期")
	}

	return nil
}

// MarkTokenAsUsed 标记Token为已使用
func (m *PasswordResetTokenManager) MarkTokenAsUsed(ctx context.Context, email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokenInfo, exists := m.tokens[email]
	if !exists {
		return fmt.Errorf("Token不存在")
	}

	tokenInfo.Used = true
	return nil
}

// CleanExpiredTokens 清理过期的Token
func (m *PasswordResetTokenManager) CleanExpiredTokens(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for email, tokenInfo := range m.tokens {
		if now.After(tokenInfo.ExpiresAt) {
			delete(m.tokens, email)
		}
	}
}
