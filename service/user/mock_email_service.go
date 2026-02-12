package user

import (
	"context"
	"sync"
)

// MockEmailService Mock邮件服务
// 用于集成测试中模拟邮件发送，避免实际发送邮件
type MockEmailService struct {
	mu               sync.Mutex
	lastToken        string
	lastResetToken   string
	sentEmails       []EmailRecord
	emailEnabled     bool
}

// EmailRecord 邮件记录
type EmailRecord struct {
	To      string
	Subject string
	Body    string
	Token   string
}

// NewMockEmailService 创建Mock邮件服务
func NewMockEmailService() *MockEmailService {
	return &MockEmailService{
		sentEmails:   make([]EmailRecord, 0),
		emailEnabled: false, // 默认禁用实际发送
	}
}

// SetEmailEnabled 设置是否启用邮件发送
func (m *MockEmailService) SetEmailEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emailEnabled = enabled
}

// SendVerificationEmail 发送验证邮件
func (m *MockEmailService) SendVerificationEmail(email, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastToken = token
	m.sentEmails = append(m.sentEmails, EmailRecord{
		To:      email,
		Subject: "邮箱验证",
		Body:    "您的验证码是: " + token,
		Token:   token,
	})

	return nil
}

// SendPasswordResetEmail 发送密码重置邮件
func (m *MockEmailService) SendPasswordResetEmail(email, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastResetToken = token
	m.sentEmails = append(m.sentEmails, EmailRecord{
		To:      email,
		Subject: "密码重置",
		Body:    "您的重置Token是: " + token,
		Token:   token,
	})

	return nil
}

// GetLastVerificationToken 获取最后的验证Token
func (m *MockEmailService) GetLastVerificationToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastToken
}

// GetLastResetToken 获取最后的重置Token
func (m *MockEmailService) GetLastResetToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastResetToken
}

// GetSentEmails 获取所有发送的邮件记录
func (m *MockEmailService) GetSentEmails() []EmailRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 返回副本
	emails := make([]EmailRecord, len(m.sentEmails))
	copy(emails, m.sentEmails)
	return emails
}

// GetEmailsTo 获取发送到指定邮箱的所有记录
func (m *MockEmailService) GetEmailsTo(email string) []EmailRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []EmailRecord
	for _, record := range m.sentEmails {
		if record.To == email {
			result = append(result, record)
		}
	}
	return result
}

// HasEmailTo 检查是否发送过邮件到指定邮箱
func (m *MockEmailService) HasEmailTo(email string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, record := range m.sentEmails {
		if record.To == email {
			return true
		}
	}
	return false
}

// Clear 清除记录
func (m *MockEmailService) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastToken = ""
	m.lastResetToken = ""
	m.sentEmails = make([]EmailRecord, 0)
}

// GetEmailCount 获取发送邮件数量
func (m *MockEmailService) GetEmailCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.sentEmails)
}

// Health 健康检查（实现EmailService接口）
func (m *MockEmailService) Health(ctx context.Context) error {
	return nil
}
