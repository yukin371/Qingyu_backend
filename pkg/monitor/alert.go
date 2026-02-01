package monitor

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// ConsoleAlerter 控制台告警（用于开发测试）
type ConsoleAlerter struct{}

// NewConsoleAlerter 创建控制台告警器
func NewConsoleAlerter() *ConsoleAlerter {
	return &ConsoleAlerter{}
}

// SendAlert 发送告警到控制台
func (a *ConsoleAlerter) SendAlert(ctx context.Context, message string) error {
	log.Printf("[ALERT] %s", message)
	return nil
}

// EmailAlerter 邮件告警（用于生产环境）
type EmailAlerter struct {
	smtpHost     string
	smtpPort     int
	from         string
	to           []string
	username     string
	password     string
	fromName     string
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	From     string
	To       []string
	Username string
	Password string
	FromName string
}

// NewEmailAlerter 创建邮件告警器
func NewEmailAlerter(config EmailConfig) *EmailAlerter {
	return &EmailAlerter{
		smtpHost: config.SMTPHost,
		smtpPort: config.SMTPPort,
		from:     config.From,
		to:       config.To,
		username: config.Username,
		password: config.Password,
		fromName: config.FromName,
	}
}

// SendAlert 发送邮件告警
func (a *EmailAlerter) SendAlert(ctx context.Context, message string) error {
	// 构建邮件内容
	subject := "【数据质量告警】Qingyu系统"
	body := fmt.Sprintf("告警时间: %s\n告警内容:\n%s\n\n请及时处理！", ctx.Value("timestamp"), message)

	// 设置邮件头
	headers := make(map[string]string)
	headers["From"] = a.formatAddress(a.fromName, a.from)
	headers["To"] = a.formatAddress("", a.to[0])
	if len(a.to) > 1 {
		for _, addr := range a.to[1:] {
			headers["To"] += ", " + addr
		}
	}
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	// 构建邮件消息
	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	// 发送邮件
	auth := smtp.PlainAuth("", a.username, a.password, a.smtpHost)
	addr := fmt.Sprintf("%s:%d", a.smtpHost, a.smtpPort)

	err := smtp.SendMail(addr, auth, a.from, a.to, []byte(msg))
	if err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// formatAddress 格式化邮件地址
func (a *EmailAlerter) formatAddress(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

// WebhookAlerter Webhook告警（支持钉钉、企业微信等）
type WebhookAlerter struct {
	webhookURL string
}

// NewWebhookAlerter 创建Webhook告警器
func NewWebhookAlerter(webhookURL string) *WebhookAlerter {
	return &WebhookAlerter{
		webhookURL: webhookURL,
	}
}

// SendAlert 发送Webhook告警
func (a *WebhookAlerter) SendAlert(ctx context.Context, message string) error {
	// TODO: 实现Webhook调用
	// 可以使用http.Post发送JSON格式的告警
	log.Printf("[WEBHOOK ALERT] URL: %s, Message: %s", a.webhookURL, message)
	return nil
}

// CompositeAlerter 组合告警器（支持多种告警方式）
type CompositeAlerter struct {
	alerters []Alerter
}

// NewCompositeAlerter 创建组合告警器
func NewCompositeAlerter(alerters ...Alerter) *CompositeAlerter {
	return &CompositeAlerter{
		alerters: alerters,
	}
}

// SendAlert 发送告警到所有配置的告警器
func (a *CompositeAlerter) SendAlert(ctx context.Context, message string) error {
	// 尝试发送到所有告警器
	var lastErr error
	for _, alerter := range a.alerters {
		if err := alerter.SendAlert(ctx, message); err != nil {
			log.Printf("告警发送失败: %v", err)
			lastErr = err
		}
	}
	return lastErr
}

// GetAlerterFromEnv 从环境变量获取告警器
func GetAlerterFromEnv() Alerter {
	// 检查是否配置了邮件告警
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	to := os.Getenv("SMTP_TO")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	if smtpHost != "" && from != "" && to != "" {
		port := 587
		if smtpPort != "" {
			fmt.Sscanf(smtpPort, "%d", &port)
		}

		config := EmailConfig{
			SMTPHost: smtpHost,
			SMTPPort: port,
			From:     from,
			To:       []string{to},
			Username: username,
			Password: password,
			FromName: "Qingyu数据质量监控",
		}
		return NewEmailAlerter(config)
	}

	// 检查是否配置了Webhook
	webhookURL := os.Getenv("ALERT_WEBHOOK_URL")
	if webhookURL != "" {
		return NewWebhookAlerter(webhookURL)
	}

	// 默认使用控制台告警
	return NewConsoleAlerter()
}
