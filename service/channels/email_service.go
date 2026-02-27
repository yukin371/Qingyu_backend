package channels

import (
	messagingModel "Qingyu_backend/models/messaging"
	"context"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// EmailService 邮件服务接口
type EmailService interface {
	// 发送邮件
	SendEmail(ctx context.Context, req *EmailRequest) error
	// 使用模板发送邮件
	SendWithTemplate(ctx context.Context, to []string, template *messagingModel.MessageTemplate, variables map[string]string) error
	// 批量发送邮件
	SendBatch(ctx context.Context, recipients []string, subject, body string) []EmailResult
	// 验证邮件地址
	ValidateEmail(email string) bool
	// 健康检查
	Health(ctx context.Context) error
}

// EmailServiceImpl 邮件服务实现
type EmailServiceImpl struct {
	config *EmailConfig
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string        // SMTP服务器地址
	SMTPPort     int           // SMTP端口
	SMTPUsername string        // SMTP用户名
	SMTPPassword string        // SMTP密码
	FromAddress  string        // 发件人邮箱
	FromName     string        // 发件人名称
	UseTLS       bool          // 是否使用TLS
	Timeout      time.Duration // 超时时间
	EnableSMTP   bool          // 是否启用真实SMTP发送（默认false，兼容旧行为）
}

// EmailRequest 邮件请求
type EmailRequest struct {
	To          []string          // 收件人列表
	Cc          []string          // 抄送列表
	Bcc         []string          // 密送列表
	Subject     string            // 主题
	Body        string            // 邮件正文
	IsHTML      bool              // 是否HTML格式
	Attachments []EmailAttachment // 附件列表
}

// EmailAttachment 邮件附件
type EmailAttachment struct {
	Filename string // 文件名
	Content  []byte // 文件内容
	MimeType string // MIME类型
}

// EmailResult 邮件发送结果
type EmailResult struct {
	Email   string // 收件人邮箱
	Success bool   // 是否成功
	Error   error  // 错误信息
}

const (
	smtpSubjectPlaceholder = "Qingyu Notification"
	smtpBodyPlaceholder    = "Please view the full message in Qingyu application inbox."
)

// NewEmailService 创建邮件服务
func NewEmailService(config *EmailConfig) EmailService {
	// 设置默认值
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.SMTPPort == 0 {
		config.SMTPPort = 587
	}

	return &EmailServiceImpl{
		config: config,
	}
}

// SendEmail 发送邮件
func (s *EmailServiceImpl) SendEmail(ctx context.Context, req *EmailRequest) error {
	if req == nil {
		return fmt.Errorf("邮件请求不能为空")
	}

	// 1. 参数验证
	if len(req.To) == 0 {
		return fmt.Errorf("收件人不能为空")
	}
	if req.Subject == "" {
		return fmt.Errorf("邮件主题不能为空")
	}
	if containsHeaderInjection(req.Subject) {
		return fmt.Errorf("邮件主题包含非法换行符")
	}
	if req.Body == "" {
		return fmt.Errorf("邮件内容不能为空")
	}
	// 2. SMTP配置验证
	if !s.config.EnableSMTP {
		return nil
	}

	// 3. SMTP配置验证
	if s.config.SMTPHost == "" || s.config.SMTPPort <= 0 {
		return fmt.Errorf("SMTP配置不完整")
	}
	if s.config.FromAddress == "" {
		return fmt.Errorf("发件人地址未配置")
	}
	if containsHeaderInjection(s.config.FromName) || containsHeaderInjection(s.config.FromAddress) {
		return fmt.Errorf("发件人信息包含非法换行符")
	}
	fromAddress, err := normalizeSMTPAddress(s.config.FromAddress)
	if err != nil {
		return fmt.Errorf("发件人地址格式无效")
	}

	cleanTo := make([]string, 0, len(req.To))
	for _, email := range req.To {
		if !s.ValidateEmail(email) || containsHeaderInjection(email) {
			return fmt.Errorf("收件人邮箱格式无效")
		}
		cleaned, err := normalizeSMTPAddress(email)
		if err != nil {
			return fmt.Errorf("收件人邮箱格式无效")
		}
		cleanTo = append(cleanTo, cleaned)
	}
	cleanCc := make([]string, 0, len(req.Cc))
	for _, email := range req.Cc {
		if !s.ValidateEmail(email) || containsHeaderInjection(email) {
			return fmt.Errorf("抄送邮箱格式无效")
		}
		cleaned, err := normalizeSMTPAddress(email)
		if err != nil {
			return fmt.Errorf("抄送邮箱格式无效")
		}
		cleanCc = append(cleanCc, cleaned)
	}
	cleanBcc := make([]string, 0, len(req.Bcc))
	for _, email := range req.Bcc {
		if !s.ValidateEmail(email) || containsHeaderInjection(email) {
			return fmt.Errorf("密送邮箱格式无效")
		}
		cleaned, err := normalizeSMTPAddress(email)
		if err != nil {
			return fmt.Errorf("密送邮箱格式无效")
		}
		cleanBcc = append(cleanBcc, cleaned)
	}

	// 4. 构建收件人列表
	recipients := make([]string, 0, len(cleanTo)+len(cleanCc)+len(cleanBcc))
	recipients = append(recipients, cleanTo...)
	recipients = append(recipients, cleanCc...)
	recipients = append(recipients, cleanBcc...)
	if len(recipients) == 0 {
		return fmt.Errorf("收件人不能为空")
	}

	// 5. 组装邮件内容
	contentType := "text/plain; charset=UTF-8"
	if req.IsHTML {
		contentType = "text/html; charset=UTF-8"
	}

	headers := []string{
		fmt.Sprintf("From: %s", formatAddress(s.config.FromName, fromAddress)),
		"To: undisclosed-recipients:;",
		fmt.Sprintf("Subject: %s", smtpSubjectPlaceholder),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: %s", contentType),
	}

	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + smtpBodyPlaceholder

	// 6. 执行发送
	addr := net.JoinHostPort(s.config.SMTPHost, strconv.Itoa(s.config.SMTPPort))
	var auth smtp.Auth
	if s.config.SMTPUsername != "" && s.config.SMTPPassword != "" {
		auth = smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	}

	if err := sendSMTPMessage(addr, auth, fromAddress, recipients, []byte(message)); err != nil {
		return fmt.Errorf("SMTP发送失败: %w", err)
	}
	return nil
}

// SendWithTemplate 使用模板发送邮件
func (s *EmailServiceImpl) SendWithTemplate(ctx context.Context, to []string, template *messagingModel.MessageTemplate, variables map[string]string) error {
	// 1. 验证模板
	if template == nil {
		return fmt.Errorf("模板不能为空")
	}
	if !template.IsActive {
		return fmt.Errorf("模板未激活")
	}

	// 2. 渲染模板
	subject := renderTemplate(template.Subject, variables)
	body := renderTemplate(template.Content, variables)

	// 3. 发送邮件
	return s.SendEmail(ctx, &EmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true, // 模板默认使用HTML格式
	})
}

// SendBatch 批量发送邮件
func (s *EmailServiceImpl) SendBatch(ctx context.Context, recipients []string, subject, body string) []EmailResult {
	results := make([]EmailResult, len(recipients))

	for i, email := range recipients {
		err := s.SendEmail(ctx, &EmailRequest{
			To:      []string{email},
			Subject: subject,
			Body:    body,
			IsHTML:  false,
		})

		results[i] = EmailResult{
			Email:   email,
			Success: err == nil,
			Error:   err,
		}

		// 添加短暂延迟，避免SMTP服务器限流
		if i < len(recipients)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return results
}

// ValidateEmail 验证邮件地址
func (s *EmailServiceImpl) ValidateEmail(email string) bool {
	// 简单的邮箱格式验证
	if email == "" {
		return false
	}

	// 检查是否包含空格
	if strings.Contains(email, " ") {
		return false
	}

	// 检查是否有连续的点
	if strings.Contains(email, "..") {
		return false
	}

	// 使用正则表达式进行更严格的验证
	// 正则说明：
	// ^[a-zA-Z0-9._%+-]+  - 用户名部分：允许字母、数字、点、下划线、百分号、加号、减号
	// @                     - @符号
	// [a-zA-Z0-9.-]+        - 域名主体：允许字母、数字、点、减号
	// \.                    - 点
	// [a-zA-Z]{2,}$         - 顶级域名：至少2个字母
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// Health 健康检查
func (s *EmailServiceImpl) Health(ctx context.Context) error {
	// 验证SMTP配置
	if s.config.SMTPHost == "" {
		return fmt.Errorf("SMTP主机未配置")
	}

	// TODO(Phase3): 实际连接SMTP服务器进行健康检查
	// 当前只验证配置
	return nil
}

// ============ 辅助函数 ============

// renderTemplate 渲染模板
func renderTemplate(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ============ SMTP认证辅助 ============

// plainAuth 明文认证
func plainAuth(username, password, host string) smtp.Auth {
	return smtp.PlainAuth("", username, password, host)
}

func formatAddress(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

func containsHeaderInjection(value string) bool {
	return strings.ContainsAny(value, "\r\n")
}

func normalizeSMTPAddress(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	parsed, err := mail.ParseAddress(trimmed)
	if err != nil || parsed == nil {
		return "", fmt.Errorf("invalid email address")
	}
	if parsed.Address != trimmed {
		return "", fmt.Errorf("invalid email address")
	}
	return strings.ToLower(parsed.Address), nil
}

func sendSMTPMessage(addr string, auth smtp.Auth, from string, recipients []string, message []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := io.Copy(writer, strings.NewReader(string(message))); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

// TODO(Phase3): 支持更多SMTP功能
// - [ ] OAuth2认证
// - [ ] 邮件模板缓存
// - [ ] 发送失败重试
// - [ ] 发送队列管理
// - [ ] 邮件发送统计
// - [ ] 反垃圾邮件处理（SPF, DKIM, DMARC）
