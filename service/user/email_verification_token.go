package user

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// EmailVerificationTokenManager 邮箱验证Token管理器
type EmailVerificationTokenManager struct {
	tokens map[string]*VerificationTokenInfo // email -> token info
	mu     sync.RWMutex
}

// VerificationTokenInfo 验证Token信息
type VerificationTokenInfo struct {
	Code      string
	Email     string
	UserID    string
	ExpiresAt time.Time
	Used      bool
}

// NewEmailVerificationTokenManager 创建邮箱验证Token管理器
func NewEmailVerificationTokenManager() *EmailVerificationTokenManager {
	return &EmailVerificationTokenManager{
		tokens: make(map[string]*VerificationTokenInfo),
	}
}

// GenerateCode 生成6位数字验证码
func (m *EmailVerificationTokenManager) GenerateCode(ctx context.Context, userID, email string) (string, error) {
	// 生成6位随机数字验证码
	codeBytes := make([]byte, 3)
	if _, err := rand.Read(codeBytes); err != nil {
		return "", fmt.Errorf("生成验证码失败: %w", err)
	}
	// 转换为6位数字字符串
	code := fmt.Sprintf("%06d", int(codeBytes[0])<<16|int(codeBytes[1])<<8|int(codeBytes[2]))
	// 取后6位确保是6位数
	if len(code) > 6 {
		code = code[len(code)-6:]
	}

	// 存储验证码信息（有效期30分钟）
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens[email] = &VerificationTokenInfo{
		Code:      code,
		Email:     email,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Used:      false,
	}

	return code, nil
}

// ValidateCode 验证验证码
func (m *EmailVerificationTokenManager) ValidateCode(ctx context.Context, userID, email, code string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tokenInfo, exists := m.tokens[email]
	if !exists {
		return fmt.Errorf("验证码不存在或已过期")
	}

	// 验证用户ID是否匹配
	if tokenInfo.UserID != userID {
		return fmt.Errorf("验证码不匹配")
	}

	if tokenInfo.Code != code {
		return fmt.Errorf("验证码错误")
	}

	if tokenInfo.Used {
		return fmt.Errorf("验证码已使用")
	}

	if time.Now().After(tokenInfo.ExpiresAt) {
		return fmt.Errorf("验证码已过期")
	}

	return nil
}

// MarkCodeAsUsed 标记验证码为已使用
func (m *EmailVerificationTokenManager) MarkCodeAsUsed(ctx context.Context, email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokenInfo, exists := m.tokens[email]
	if !exists {
		return fmt.Errorf("验证码不存在")
	}

	tokenInfo.Used = true
	return nil
}

// CleanExpiredCodes 清理过期的验证码
func (m *EmailVerificationTokenManager) CleanExpiredCodes(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for email, tokenInfo := range m.tokens {
		if now.After(tokenInfo.ExpiresAt) {
			delete(m.tokens, email)
		}
	}
}
