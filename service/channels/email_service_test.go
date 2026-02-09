package channels

import (
	"Qingyu_backend/models/messaging"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============ 测试辅助函数 ============

// newTestEmailConfig 创建测试用的EmailConfig
func newTestEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     "localhost",
		SMTPPort:     1025,
		SMTPUsername: "test@example.com",
		SMTPPassword: "testpassword",
		FromAddress:  "test@example.com",
		FromName:     "Test Sender",
		UseTLS:       false,
		Timeout:      5 * time.Second,
	}
}

// newTestEmailRequest 创建测试用的EmailRequest
func newTestEmailRequest() *EmailRequest {
	return &EmailRequest{
		To:      []string{"recipient@example.com"},
		Cc:      []string{},
		Bcc:     []string{},
		Subject: "Test Subject",
		Body:    "Test Body",
		IsHTML:  false,
	}
}

// newTestEmailRequestWithCC 创建带抄送的EmailRequest
func newTestEmailRequestWithCC(cc []string) *EmailRequest {
	req := newTestEmailRequest()
	req.Cc = cc
	return req
}

// newTestEmailRequestWithBCC 创建带密送的EmailRequest
func newTestEmailRequestWithBCC(bcc []string) *EmailRequest {
	req := newTestEmailRequest()
	req.Bcc = bcc
	return req
}

// newTestEmailRequestWithAttachments 创建带附件的EmailRequest
func newTestEmailRequestWithAttachments() *EmailRequest {
	req := newTestEmailRequest()
	req.Attachments = []EmailAttachment{
		{
			Filename: "test.txt",
			Content:  []byte("test content"),
			MimeType: "text/plain",
		},
	}
	return req
}

// newTestContext 创建测试用的context
func newTestContext() context.Context {
	return context.Background()
}

// newActiveTemplate 创建激活的模板
func newActiveTemplate() *messaging.MessageTemplate {
	return &messaging.MessageTemplate{
		ID:        "tpl-001",
		Name:      "welcome",
		Subject:   "Welcome {{Name}}",
		Content:   "Hello {{Name}}, welcome to our service!",
		IsActive:  true,
		Variables: []string{"Name"},
	}
}

// newInactiveTemplate 创建未激活的模板
func newInactiveTemplate() *messaging.MessageTemplate {
	tpl := newActiveTemplate()
	tpl.IsActive = false
	return tpl
}

// newTemplateWithMultipleVariables 创建多变量模板
func newTemplateWithMultipleVariables() *messaging.MessageTemplate {
	return &messaging.MessageTemplate{
		ID:        "tpl-002",
		Name:      "order-confirmation",
		Subject:   "Order Confirmation - Order #{{OrderID}}",
		Content:   "Dear {{Name}}, your order #{{OrderID}} has been confirmed. Total: ${{Amount}}",
		IsActive:  true,
		Variables: []string{"Name", "OrderID", "Amount"},
	}
}

// validVariablesForActiveTemplate 返回ActiveTemplate的有效变量
func validVariablesForActiveTemplate() map[string]string {
	return map[string]string{
		"Name": "John Doe",
	}
}

// validVariablesForMultipleVariables 返回多变量模板的有效变量
func validVariablesForMultipleVariables() map[string]string {
	return map[string]string{
		"Name":    "Alice",
		"OrderID": "12345",
		"Amount":  "99.99",
	}
}

// invalidVariablesForTemplate 返回缺少变量的数据
func invalidVariablesForTemplate() map[string]string {
	return map[string]string{
		"InvalidKey": "InvalidValue",
	}
}

// emptyVariables 返回空变量
func emptyVariables() map[string]string {
	return map[string]string{}
}

// ============ SendEmail 测试 ============

// TestSendEmail_Success 测试成功发送邮件
func TestSendEmail_Success(t *testing.T) {
	// Arrange - 准备测试数据
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequest()
	ctx := newTestContext()

	// Act - 执行被测试的方法
	err := service.SendEmail(ctx, req)

	// Assert - 验证结果
	assert.NoError(t, err, "发送邮件应该成功")
}

// TestSendEmail_InvalidToEmpty 测试收件人为空的情况
func TestSendEmail_InvalidToEmpty(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequest()
	req.To = []string{}
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.Error(t, err, "收件人为空应该返回错误")
	assert.Contains(t, err.Error(), "收件人不能为空", "错误信息应该包含'收件人不能为空'")
}

// TestSendEmail_InvalidSubjectEmpty 测试主题为空的情况
func TestSendEmail_InvalidSubjectEmpty(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequest()
	req.Subject = ""
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.Error(t, err, "主题为空应该返回错误")
	assert.Contains(t, err.Error(), "邮件主题不能为空", "错误信息应该包含'邮件主题不能为空'")
}

// TestSendEmail_InvalidBodyEmpty 测试正文为空的情况
func TestSendEmail_InvalidBodyEmpty(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequest()
	req.Body = ""
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.Error(t, err, "正文为空应该返回错误")
	assert.Contains(t, err.Error(), "邮件内容不能为空", "错误信息应该包含'邮件内容不能为空'")
}

// TestSendEmail_WithCC 测试带抄送的邮件
func TestSendEmail_WithCC(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequestWithCC([]string{"cc@example.com"})
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.NoError(t, err, "带抄送的邮件应该发送成功")
	assert.Len(t, req.Cc, 1, "抄送列表应该有1个收件人")
}

// TestSendEmail_WithBCC 测试带密送的邮件
func TestSendEmail_WithBCC(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequestWithBCC([]string{"bcc@example.com"})
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.NoError(t, err, "带密送的邮件应该发送成功")
	assert.Len(t, req.Bcc, 1, "密送列表应该有1个收件人")
}

// TestSendEmail_WithAttachments 测试带附件的邮件
func TestSendEmail_WithAttachments(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequestWithAttachments()
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.NoError(t, err, "带附件的邮件应该发送成功")
	assert.Len(t, req.Attachments, 1, "附件列表应该有1个附件")
}

// TestSendEmail_HTMLFormat 测试HTML格式邮件
func TestSendEmail_HTMLFormat(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config)
	req := newTestEmailRequest()
	req.IsHTML = true
	req.Body = "<html><body><h1>Test</h1></body></html>"
	ctx := newTestContext()

	// Act
	err := service.SendEmail(ctx, req)

	// Assert
	assert.NoError(t, err, "HTML格式邮件应该发送成功")
	assert.True(t, req.IsHTML, "邮件应该是HTML格式")
}

// ============ SendWithTemplate 测试 ============

// TestSendWithTemplate_Success 测试使用模板成功发送邮件
func TestSendWithTemplate_Success(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	template := newActiveTemplate()
	variables := validVariablesForActiveTemplate()
	to := []string{"recipient@example.com"}
	ctx := newTestContext()

	// Act
	err := service.SendWithTemplate(ctx, to, template, variables)

	// Assert
	assert.NoError(t, err, "使用模板发送邮件应该成功")
}

// TestSendWithTemplate_TemplateNil 测试模板为nil的情况
func TestSendWithTemplate_TemplateNil(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	variables := validVariablesForActiveTemplate()
	to := []string{"recipient@example.com"}
	ctx := newTestContext()

	// Act
	err := service.SendWithTemplate(ctx, to, nil, variables)

	// Assert
	assert.Error(t, err, "模板为nil应该返回错误")
	assert.Contains(t, err.Error(), "模板不能为空", "错误信息应该包含'模板不能为空'")
}

// TestSendWithTemplate_TemplateInactive 测试使用未激活的模板
func TestSendWithTemplate_TemplateInactive(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	template := newInactiveTemplate()
	variables := validVariablesForActiveTemplate()
	to := []string{"recipient@example.com"}
	ctx := newTestContext()

	// Act
	err := service.SendWithTemplate(ctx, to, template, variables)

	// Assert
	assert.Error(t, err, "使用未激活的模板应该返回错误")
	assert.Contains(t, err.Error(), "模板未激活", "错误信息应该包含'模板未激活'")
}

// TestSendWithTemplate_InvalidVariables 测试使用无效变量
func TestSendWithTemplate_InvalidVariables(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	template := newActiveTemplate()
	variables := invalidVariablesForTemplate()
	to := []string{"recipient@example.com"}
	ctx := newTestContext()

	// Act
	err := service.SendWithTemplate(ctx, to, template, variables)

	// Assert
	// 注意：当前实现中，即使变量不匹配也会发送，只是占位符不会被替换
	assert.NoError(t, err, "即使变量无效，当前实现也会发送邮件（占位符不替换）")
}

// TestSendWithTemplate_EmptyVariables 测试使用空变量
func TestSendWithTemplate_EmptyVariables(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	template := newActiveTemplate()
	variables := emptyVariables()
	to := []string{"recipient@example.com"}
	ctx := newTestContext()

	// Act
	err := service.SendWithTemplate(ctx, to, template, variables)

	// Assert
	assert.NoError(t, err, "即使变量为空，当前实现也会发送邮件")
}

// TestSendWithTemplate_MultipleVariables 测试使用多变量模板
func TestSendWithTemplate_MultipleVariables(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	template := newTemplateWithMultipleVariables()
	variables := validVariablesForMultipleVariables()
	to := []string{"recipient@example.com"}
	ctx := newTestContext()

	// Act
	err := service.SendWithTemplate(ctx, to, template, variables)

	// Assert
	assert.NoError(t, err, "使用多变量模板发送邮件应该成功")
}

// ============ SendBatch 测试 ============

// TestSendBatch_AllSuccess 测试批量发送全部成功
func TestSendBatch_AllSuccess(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	recipients := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
	subject := "Test Subject"
	body := "Test Body"
	ctx := newTestContext()

	// Act
	results := service.SendBatch(ctx, recipients, subject, body)

	// Assert
	require.Len(t, results, len(recipients), "结果数量应该与收件人数量相同")
	for _, result := range results {
		assert.True(t, result.Success, "所有邮件都应该发送成功: "+result.Email)
		assert.NoError(t, result.Error, "不应该有错误: "+result.Email)
	}
}

// TestSendBatch_PartialFailure 测试批量发送部分失败
func TestSendBatch_PartialFailure(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	// 使用空的body来触发部分失败
	recipients := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
	subject := "Test Subject"
	body := "" // 空body会导致失败
	ctx := newTestContext()

	// Act
	results := service.SendBatch(ctx, recipients, subject, body)

	// Assert
	require.Len(t, results, len(recipients), "结果数量应该与收件人数量相同")
	for _, result := range results {
		assert.False(t, result.Success, "所有邮件都应该发送失败（因为body为空）")
		assert.Error(t, result.Error, "应该有错误: "+result.Email)
	}
}

// TestSendBatch_AllFailure 测试批量发送全部失败
func TestSendBatch_AllFailure(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	recipients := []string{"user1@example.com", "user2@example.com"}
	subject := "" // 空主题会导致失败
	body := "Test Body"
	ctx := newTestContext()

	// Act
	results := service.SendBatch(ctx, recipients, subject, body)

	// Assert
	require.Len(t, results, len(recipients), "结果数量应该与收件人数量相同")
	for _, result := range results {
		assert.False(t, result.Success, "所有邮件都应该发送失败")
		assert.Error(t, result.Error, "应该有错误: "+result.Email)
	}
}

// TestSendBatch_EmptyRecipients 测试空收件人列表
func TestSendBatch_EmptyRecipients(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	recipients := []string{}
	subject := "Test Subject"
	body := "Test Body"
	ctx := newTestContext()

	// Act
	results := service.SendBatch(ctx, recipients, subject, body)

	// Assert
	assert.Len(t, results, 0, "空收件人列表应该返回空结果")
}

// TestSendBatch_SingleRecipient 测试单个收件人
func TestSendBatch_SingleRecipient(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	recipients := []string{"user1@example.com"}
	subject := "Test Subject"
	body := "Test Body"
	ctx := newTestContext()

	// Act
	results := service.SendBatch(ctx, recipients, subject, body)

	// Assert
	require.Len(t, results, 1, "结果数量应该为1")
	assert.True(t, results[0].Success, "邮件应该发送成功")
	assert.Equal(t, "user1@example.com", results[0].Email, "收件人邮箱应该正确")
}

// ============ ValidateEmail 测试 ============

// TestValidateEmail_Valid 测试有效邮箱验证
func TestValidateEmail_Valid(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.com",
		"user123@test.co.uk",
		"firstname-lastname@example.com",
	}

	// Act & Assert
	for _, email := range validEmails {
		t.Run(email, func(t *testing.T) {
			result := service.ValidateEmail(email)
			assert.True(t, result, "邮箱 '%s' 应该是有效的", email)
		})
	}
}

// TestValidateEmail_Invalid 测试无效邮箱验证
func TestValidateEmail_Invalid(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	invalidEmails := []string{
		"",
		"invalid",
		"@example.com",
		"user@",
		"user@@example.com",
		"user name@example.com",
		"user@example..com",
	}

	// Act & Assert
	for _, email := range invalidEmails {
		t.Run(email, func(t *testing.T) {
			result := service.ValidateEmail(email)
			assert.False(t, result, "邮箱 '%s' 应该是无效的", email)
		})
	}
}

// TestValidateEmail_Empty 测试空字符串
func TestValidateEmail_Empty(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	// Act
	result := service.ValidateEmail("")

	// Assert
	assert.False(t, result, "空字符串应该是无效的邮箱")
}

// TestValidateEmail_NoAtSymbol 测试没有@符号的邮箱
func TestValidateEmail_NoAtSymbol(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	// Act
	result := service.ValidateEmail("invalidemail.com")

	// Assert
	assert.False(t, result, "没有@符号的邮箱应该是无效的")
}

// TestValidateEmail_MultipleAtSymbols 测试多个@符号的邮箱
func TestValidateEmail_MultipleAtSymbols(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	// Act
	result := service.ValidateEmail("user@@example.com")

	// Assert
	assert.False(t, result, "多个@符号的邮箱应该是无效的")
}

// TestValidateEmail_UsernameEmpty 测试用户名为空
func TestValidateEmail_UsernameEmpty(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	// Act
	result := service.ValidateEmail("@example.com")

	// Assert
	assert.False(t, result, "用户名为空的邮箱应该是无效的")
}

// TestValidateEmail_DomainEmpty 测试域名为空
func TestValidateEmail_DomainEmpty(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)

	// Act
	result := service.ValidateEmail("user@")

	// Assert
	assert.False(t, result, "域名为空的邮箱应该是无效的")
}

// ============ Health 测试 ============

// TestHealth_Healthy 测试健康状态
func TestHealth_Healthy(t *testing.T) {
	// Arrange
	config := newTestEmailConfig()
	service := NewEmailService(config).(*EmailServiceImpl)
	ctx := newTestContext()

	// Act
	err := service.Health(ctx)

	// Assert
	assert.NoError(t, err, "配置正确时健康检查应该通过")
}

// TestHealth_Unhealthy 测试不健康状态
func TestHealth_Unhealthy(t *testing.T) {
	// Arrange
	config := &EmailConfig{
		SMTPHost: "", // 空主机名
		SMTPPort: 587,
	}
	service := NewEmailService(config).(*EmailServiceImpl)
	ctx := newTestContext()

	// Act
	err := service.Health(ctx)

	// Assert
	assert.Error(t, err, "SMTP主机未配置时健康检查应该失败")
	assert.Contains(t, err.Error(), "SMTP主机未配置", "错误信息应该包含'SMTP主机未配置'")
}

// ============ NewEmailService 测试 ============

// TestNewEmailService_DefaultValues 测试默认值设置
func TestNewEmailService_DefaultValues(t *testing.T) {
	// Arrange
	config := &EmailConfig{
		SMTPHost:    "localhost",
		SMTPPort:    0, // 使用默认值
		Timeout:     0, // 使用默认值
		FromAddress: "test@example.com",
	}

	// Act
	service := NewEmailService(config).(*EmailServiceImpl)

	// Assert
	assert.Equal(t, 587, service.config.SMTPPort, "默认SMTP端口应该是587")
	assert.Equal(t, 10*time.Second, service.config.Timeout, "默认超时应该是10秒")
}

// TestNewEmailService_CustomValues 测试自定义值
func TestNewEmailService_CustomValues(t *testing.T) {
	// Arrange
	config := &EmailConfig{
		SMTPHost:    "smtp.example.com",
		SMTPPort:    465,
		Timeout:     30 * time.Second,
		FromAddress: "test@example.com",
	}

	// Act
	service := NewEmailService(config).(*EmailServiceImpl)

	// Assert
	assert.Equal(t, "smtp.example.com", service.config.SMTPHost, "SMTP主机应该正确设置")
	assert.Equal(t, 465, service.config.SMTPPort, "SMTP端口应该正确设置")
	assert.Equal(t, 30*time.Second, service.config.Timeout, "超时应该正确设置")
}

// ============ renderTemplate 辅助函数测试 ============

// TestRenderTemplate_SingleVariable 测试单个变量替换
func TestRenderTemplate_SingleVariable(t *testing.T) {
	// Arrange
	template := "Hello {{Name}}!"
	variables := map[string]string{
		"Name": "World",
	}

	// Act
	result := renderTemplate(template, variables)

	// Assert
	assert.Equal(t, "Hello World!", result, "模板渲染结果应该正确")
}

// TestRenderTemplate_MultipleVariables 测试多个变量替换
func TestRenderTemplate_MultipleVariables(t *testing.T) {
	// Arrange
	template := "Dear {{Name}}, your order #{{OrderID}} is confirmed."
	variables := map[string]string{
		"Name":    "Alice",
		"OrderID": "12345",
	}

	// Act
	result := renderTemplate(template, variables)

	// Assert
	assert.Equal(t, "Dear Alice, your order #12345 is confirmed.", result, "模板渲染结果应该正确")
}

// TestRenderTemplate_MissingVariable 测试缺少变量
func TestRenderTemplate_MissingVariable(t *testing.T) {
	// Arrange
	template := "Hello {{Name}}!"
	variables := map[string]string{
		"Other": "Value",
	}

	// Act
	result := renderTemplate(template, variables)

	// Assert
	// 缺少变量时，占位符不会被替换
	assert.Equal(t, "Hello {{Name}}!", result, "缺少变量时占位符应该保持不变")
}

// TestRenderTemplate_EmptyVariables 测试空变量
func TestRenderTemplate_EmptyVariables(t *testing.T) {
	// Arrange
	template := "Hello {{Name}}!"
	variables := map[string]string{}

	// Act
	result := renderTemplate(template, variables)

	// Assert
	assert.Equal(t, "Hello {{Name}}!", result, "空变量时占位符应该保持不变")
}

// TestRenderTemplate_NoPlaceholders 测试无占位符的模板
func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	// Arrange
	template := "Hello World!"
	variables := map[string]string{
		"Name": "Test",
	}

	// Act
	result := renderTemplate(template, variables)

	// Assert
	assert.Equal(t, "Hello World!", result, "无占位符的模板应该保持不变")
}

// TestRenderTemplate_RepeatedVariable 测试重复变量
func TestRenderTemplate_RepeatedVariable(t *testing.T) {
	// Arrange
	template := "{{Name}} {{Name}} {{Name}}"
	variables := map[string]string{
		"Name": "Test",
	}

	// Act
	result := renderTemplate(template, variables)

	// Assert
	assert.Equal(t, "Test Test Test", result, "重复的变量应该都被替换")
}
