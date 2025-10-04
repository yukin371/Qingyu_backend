package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators 注册所有自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
	// 金额验证
	v.RegisterValidation("amount", validateAmount)
	v.RegisterValidation("positive_amount", validatePositiveAmount)
	v.RegisterValidation("amount_range", validateAmountRange)

	// 文件验证
	v.RegisterValidation("file_type", validateFileType)
	v.RegisterValidation("file_size", validateFileSize)

	// 字符串验证
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("strong_password", validateStrongPassword)

	// 业务验证
	v.RegisterValidation("transaction_type", validateTransactionType)
	v.RegisterValidation("withdraw_account", validateWithdrawAccount)
	v.RegisterValidation("content_type", validateContentType)
}

// validateAmount 验证金额格式（最多2位小数）
func validateAmount(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	if amount < 0 {
		return false
	}
	// 检查小数位数（最多2位）
	amountStr := fl.Field().String()
	parts := strings.Split(amountStr, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return false
	}
	return true
}

// validatePositiveAmount 验证正数金额（> 0）
func validatePositiveAmount(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	return amount > 0
}

// validateAmountRange 验证金额范围（0.01 - 1000000）
func validateAmountRange(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	return amount >= 0.01 && amount <= 1000000.00
}

// validateFileType 验证文件类型
func validateFileType(fl validator.FieldLevel) bool {
	fileType := fl.Field().String()
	allowedTypes := []string{
		"image/jpeg", "image/png", "image/gif", "image/webp",
		"application/pdf", "application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/plain", "application/zip",
	}
	for _, t := range allowedTypes {
		if fileType == t {
			return true
		}
	}
	return false
}

// validateFileSize 验证文件大小（最大50MB）
func validateFileSize(fl validator.FieldLevel) bool {
	size := fl.Field().Int()
	maxSize := int64(50 * 1024 * 1024) // 50MB
	return size > 0 && size <= maxSize
}

// validateUsername 验证用户名（3-20个字符，字母数字下划线）
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

// validatePhone 验证手机号（中国大陆）
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// validateStrongPassword 验证强密码（至少8位，包含大小写字母和数字）
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper && hasLower && hasNumber
}

// validateTransactionType 验证交易类型
func validateTransactionType(fl validator.FieldLevel) bool {
	txType := fl.Field().String()
	validTypes := []string{"recharge", "consume", "transfer", "refund", "withdraw"}
	for _, t := range validTypes {
		if txType == t {
			return true
		}
	}
	return false
}

// validateWithdrawAccount 验证提现账号格式
func validateWithdrawAccount(fl validator.FieldLevel) bool {
	account := fl.Field().String()
	// 格式：支付方式:账号 例如 "alipay:user@example.com" 或 "wechat:wxid_xxx"
	parts := strings.Split(account, ":")
	if len(parts) != 2 {
		return false
	}
	method := parts[0]
	accountID := parts[1]

	// 验证支付方式
	validMethods := []string{"alipay", "wechat", "bank"}
	methodValid := false
	for _, m := range validMethods {
		if method == m {
			methodValid = true
			break
		}
	}
	if !methodValid {
		return false
	}

	// 验证账号不为空
	return len(accountID) > 0
}

// validateContentType 验证内容类型
func validateContentType(fl validator.FieldLevel) bool {
	contentType := fl.Field().String()
	validTypes := []string{"book", "chapter", "comment", "review"}
	for _, t := range validTypes {
		if contentType == t {
			return true
		}
	}
	return false
}
