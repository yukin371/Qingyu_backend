package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors 验证错误列表
type ValidationErrors []ValidationError

// TranslateError 将validator错误转换为友好的中文错误消息
func TranslateError(err error) ValidationErrors {
	var errors ValidationErrors

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   getFieldName(e),
				Message: getErrorMessage(e),
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
			})
		}
	}

	return errors
}

// getFieldName 获取字段名（转换为小写下划线格式）
func getFieldName(fe validator.FieldError) string {
	field := fe.Field()
	// 转换驼峰命名为下划线命名
	var result strings.Builder
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getErrorMessage 根据验证标签返回友好的错误消息
func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s 是必填字段", field)
	case "email":
		return fmt.Sprintf("%s 必须是有效的邮箱地址", field)
	case "min":
		return fmt.Sprintf("%s 长度不能小于 %s", field, param)
	case "max":
		return fmt.Sprintf("%s 长度不能大于 %s", field, param)
	case "len":
		return fmt.Sprintf("%s 长度必须等于 %s", field, param)
	case "gt":
		return fmt.Sprintf("%s 必须大于 %s", field, param)
	case "gte":
		return fmt.Sprintf("%s 必须大于等于 %s", field, param)
	case "lt":
		return fmt.Sprintf("%s 必须小于 %s", field, param)
	case "lte":
		return fmt.Sprintf("%s 必须小于等于 %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s 必须是以下值之一: %s", field, param)
	case "url":
		return fmt.Sprintf("%s 必须是有效的URL", field)
	case "uri":
		return fmt.Sprintf("%s 必须是有效的URI", field)
	case "alpha":
		return fmt.Sprintf("%s 只能包含字母", field)
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", field)
	case "numeric":
		return fmt.Sprintf("%s 必须是数字", field)

	// 自定义验证器
	case "amount":
		return fmt.Sprintf("%s 必须是有效的金额格式（最多2位小数）", field)
	case "positive_amount":
		return fmt.Sprintf("%s 必须是正数", field)
	case "amount_range":
		return fmt.Sprintf("%s 必须在 0.01 到 1,000,000.00 之间", field)
	case "file_type":
		return fmt.Sprintf("%s 文件类型不支持", field)
	case "file_size":
		return fmt.Sprintf("%s 文件大小不能超过 50MB", field)
	case "username":
		return fmt.Sprintf("%s 必须是3-20个字符的字母、数字或下划线", field)
	case "phone":
		return fmt.Sprintf("%s 必须是有效的手机号", field)
	case "strong_password":
		return fmt.Sprintf("%s 必须至少8位，包含大小写字母和数字", field)
	case "transaction_type":
		return fmt.Sprintf("%s 必须是有效的交易类型（recharge/consume/transfer/refund/withdraw）", field)
	case "withdraw_account":
		return fmt.Sprintf("%s 格式不正确，应为 '支付方式:账号'（如 alipay:user@example.com）", field)
	case "content_type":
		return fmt.Sprintf("%s 必须是有效的内容类型（book/chapter/comment/review）", field)

	default:
		return fmt.Sprintf("%s 验证失败", field)
	}
}

// FormatErrors 格式化错误消息为字符串
func (ve ValidationErrors) FormatErrors() string {
	if len(ve) == 0 {
		return ""
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// GetFieldErrors 获取字段级错误映射
func (ve ValidationErrors) GetFieldErrors() map[string]string {
	fieldErrors := make(map[string]string)
	for _, err := range ve {
		fieldErrors[err.Field] = err.Message
	}
	return fieldErrors
}
