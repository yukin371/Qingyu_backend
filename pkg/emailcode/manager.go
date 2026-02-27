package emailcode

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"sync"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/service/channels"
)

const (
	defaultCodeLength = 6
	defaultCodeTTL    = 10 * time.Minute
	defaultCooldown   = 60 * time.Second
	maxVerifyAttempts = 5
)

type codeRecord struct {
	CodeHash  [32]byte
	ExpiresAt time.Time
	SentAt    time.Time
	Attempts  int
}

// Manager 管理邮箱验证码发送与校验
type Manager struct {
	enabled bool

	emailService channels.EmailService
	fromName     string

	codeTTL  time.Duration
	cooldown time.Duration

	mu    sync.RWMutex
	codes map[string]*codeRecord
}

// NewManager 创建验证码管理器
func NewManager() *Manager {
	m := &Manager{
		codeTTL:  defaultCodeTTL,
		cooldown: defaultCooldown,
		codes:    make(map[string]*codeRecord),
	}

	if config.GlobalConfig == nil || config.GlobalConfig.Email == nil {
		return m
	}

	emailCfg := config.GlobalConfig.Email
	if !emailCfg.Enabled {
		return m
	}

	if emailCfg.SMTPHost == "" || emailCfg.SMTPPort == 0 || emailCfg.Username == "" || emailCfg.Password == "" || emailCfg.FromAddress == "" {
		return m
	}

	m.emailService = channels.NewEmailService(&channels.EmailConfig{
		SMTPHost:     emailCfg.SMTPHost,
		SMTPPort:     emailCfg.SMTPPort,
		SMTPUsername: emailCfg.Username,
		SMTPPassword: emailCfg.Password,
		FromAddress:  emailCfg.FromAddress,
		FromName:     emailCfg.FromName,
		UseTLS:       emailCfg.UseTLS,
		Timeout:      10 * time.Second,
		EnableSMTP:   true,
	})
	m.fromName = emailCfg.FromName
	m.enabled = true
	return m
}

// Enabled 返回邮箱验证码是否启用
func (m *Manager) Enabled() bool {
	return m.enabled
}

// SendRegisterCode 发送注册验证码
func (m *Manager) SendRegisterCode(ctx context.Context, email string) error {
	if !m.enabled {
		return fmt.Errorf("邮箱验证码功能未启用，请先配置 QINGYU_EMAIL_ENABLED 和 SMTP 参数")
	}

	normalized, err := normalizeEmail(email)
	if err != nil {
		return err
	}

	now := time.Now()

	m.mu.Lock()
	existing, ok := m.codes[normalized]
	if ok && now.Sub(existing.SentAt) < m.cooldown {
		wait := m.cooldown - now.Sub(existing.SentAt)
		m.mu.Unlock()
		return fmt.Errorf("验证码发送过于频繁，请在 %d 秒后重试", int(wait.Seconds())+1)
	}
	m.mu.Unlock()

	code, err := generateNumericCode(defaultCodeLength)
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	subject := "青羽阅读邮箱验证码"
	body := fmt.Sprintf("您的注册验证码是：%s\n\n验证码 %d 分钟内有效，请勿泄露给他人。", code, int(m.codeTTL.Minutes()))
	if m.fromName != "" {
		subject = fmt.Sprintf("%s - 邮箱验证码", m.fromName)
	}

	if err := m.emailService.SendEmail(ctx, &channels.EmailRequest{
		To:      []string{normalized},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}); err != nil {
		return fmt.Errorf("发送验证码失败: %w", err)
	}

	record := &codeRecord{
		CodeHash:  sha256.Sum256([]byte(code)),
		ExpiresAt: now.Add(m.codeTTL),
		SentAt:    now,
		Attempts:  0,
	}

	m.mu.Lock()
	m.codes[normalized] = record
	m.mu.Unlock()
	return nil
}

// VerifyRegisterCode 校验注册验证码
func (m *Manager) VerifyRegisterCode(email, code string) error {
	if !m.enabled {
		return nil
	}

	normalized, err := normalizeEmail(email)
	if err != nil {
		return err
	}
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("验证码不能为空")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	record, ok := m.codes[normalized]
	if !ok {
		return fmt.Errorf("验证码不存在或已过期")
	}

	now := time.Now()
	if now.After(record.ExpiresAt) {
		delete(m.codes, normalized)
		return fmt.Errorf("验证码已过期")
	}

	if record.Attempts >= maxVerifyAttempts {
		delete(m.codes, normalized)
		return fmt.Errorf("验证码错误次数过多，请重新获取")
	}

	inputHash := sha256.Sum256([]byte(strings.TrimSpace(code)))
	if subtle.ConstantTimeCompare(inputHash[:], record.CodeHash[:]) != 1 {
		record.Attempts++
		return fmt.Errorf("验证码错误")
	}

	delete(m.codes, normalized)
	return nil
}

func normalizeEmail(email string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(email))
	if trimmed == "" {
		return "", fmt.Errorf("邮箱不能为空")
	}
	if _, err := mail.ParseAddress(trimmed); err != nil {
		return "", fmt.Errorf("邮箱格式不正确")
	}
	return trimmed, nil
}

func generateNumericCode(length int) (string, error) {
	var b strings.Builder
	b.Grow(length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b.WriteByte(byte('0') + byte(n.Int64()))
	}
	return b.String(), nil
}
