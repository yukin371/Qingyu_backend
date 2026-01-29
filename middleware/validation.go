package middleware

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// 验证正则表达式
var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{2,31}$`)
	urlRegex      = regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=]+$`)
	htmlTagRegex  = regexp.MustCompile(`<[^>]*>`)
	sqlKeywords   = []string{
		"DROP", "DELETE", "TRUNCATE", "ALTER", "CREATE", "INSERT",
		"UPDATE", "UNION", "SELECT", "EXEC", "EXECUTE",
	}
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidatePhone 验证手机号格式（中国大陆）
func ValidatePhone(phone string) bool {
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

// ValidatePassword 验证密码强度
// 要求：8-128字符，包含大小写字母、数字和特殊字符
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password too short, must be at least 8 characters")
	}
	if len(password) > 128 {
		return errors.New("password too long, must be at most 128 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateUsername 验证用户名格式
// 要求：3-32字符，以字母开头，只能包含字母、数字、下划线和连字符
func ValidateUsername(username string) bool {
	if username == "" {
		return false
	}
	return usernameRegex.MatchString(username)
}

// ValidateURL 验证URL格式
// 要求：以http://或https://开头
func ValidateURL(url string) bool {
	if url == "" {
		return false
	}
	return urlRegex.MatchString(url)
}

// SanitizeString 清理字符串，防止XSS和SQL注入
func SanitizeString(input string) string {
	// 去除首尾空白
	input = strings.TrimSpace(input)

	// 移除SQL注入关键字（在HTML标签清理之前）
	for _, keyword := range sqlKeywords {
		keywordRegex := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, keyword))
		input = keywordRegex.ReplaceAllString(input, "")
	}

	// 去除HTML标签及其内容
	input = htmlTagRegex.ReplaceAllString(input, "")

	// 替换多个空白字符为单个空格
	spaceRegex := regexp.MustCompile(`\s+`)
	input = spaceRegex.ReplaceAllString(input, " ")

	// 再次清理空白
	input = strings.TrimSpace(input)

	return input
}
