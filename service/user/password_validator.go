package user

import (
	"regexp"
	"strings"
)

// PasswordValidator 密码验证器
type PasswordValidator struct {
	minLength        int
	requireUppercase bool
	requireLowercase bool
	requireDigit     bool
	requireSpecial   bool
	commonPasswords  map[string]bool
}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:        8,
		requireUppercase: true,
		requireLowercase: true,
		requireDigit:     true,
		requireSpecial:   false, // 特殊字符可选
		commonPasswords:  loadCommonPasswords(),
	}
}

// ValidateStrength 验证密码强度
func (v *PasswordValidator) ValidateStrength(password string) (bool, string) {
	// 1. 检查长度
	if len(password) < v.minLength {
		return false, "密码长度不能少于8位"
	}

	// 2. 检查大写字母
	if v.requireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false, "密码必须包含大写字母"
	}

	// 3. 检查小写字母
	if v.requireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false, "密码必须包含小写字母"
	}

	// 4. 检查数字
	if v.requireDigit && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false, "密码必须包含数字"
	}

	// 5. 检查特殊字符（可选）
	if v.requireSpecial && !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return false, "密码必须包含特殊字符"
	}

	// 6. 检查是否是常见弱密码
	if v.IsCommonPassword(password) {
		return false, "密码过于常见，请使用更复杂的密码"
	}

	// 7. 检查连续字符
	if hasSequentialChars(password) {
		return false, "密码不能包含连续的字符（如123、abc）"
	}

	return true, ""
}

// IsCommonPassword 检查是否是常见密码
func (v *PasswordValidator) IsCommonPassword(password string) bool {
	lowerPassword := strings.ToLower(password)
	return v.commonPasswords[lowerPassword]
}

// GetStrengthScore 获取密码强度评分（0-100）
func (v *PasswordValidator) GetStrengthScore(password string) int {
	score := 0

	// 长度评分（最多30分）
	if len(password) >= 8 {
		score += 10
	}
	if len(password) >= 12 {
		score += 10
	}
	if len(password) >= 16 {
		score += 10
	}

	// 字符类型评分（每种15分）
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score += 15
	}
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score += 15
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		score += 15
	}
	if regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		score += 15
	}

	// 扣分项
	if v.IsCommonPassword(password) {
		score -= 30
	}
	if hasSequentialChars(password) {
		score -= 20
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// GetStrengthLevel 获取密码强度等级
func (v *PasswordValidator) GetStrengthLevel(password string) string {
	score := v.GetStrengthScore(password)

	if score >= 80 {
		return "强"
	} else if score >= 60 {
		return "中等"
	} else if score >= 40 {
		return "一般"
	} else {
		return "弱"
	}
}

// ============ 辅助函数 ============

// hasSequentialChars 检查是否有连续字符
func hasSequentialChars(s string) bool {
	// 检查连续数字（123, 234等）
	for i := 0; i < len(s)-2; i++ {
		if s[i] >= '0' && s[i] <= '7' {
			if s[i]+1 == s[i+1] && s[i]+2 == s[i+2] {
				return true
			}
		}
	}

	// 检查连续字母（abc, bcd等）
	lower := strings.ToLower(s)
	for i := 0; i < len(lower)-2; i++ {
		if lower[i] >= 'a' && lower[i] <= 'x' {
			if lower[i]+1 == lower[i+1] && lower[i]+2 == lower[i+2] {
				return true
			}
		}
	}

	return false
}

// loadCommonPasswords 加载常见弱密码列表
func loadCommonPasswords() map[string]bool {
	// TODO(Phase3): 从文件或数据库加载完整的弱密码字典
	common := []string{
		"password", "password123", "admin123", "user1234",
		"test1234", "123456", "12345678", "123456789",
		"qwerty", "abc123", "111111", "000000",
		"admin", "root", "user", "test",
		"welcome", "login", "passw0rd", "letmein",
	}

	passwordMap := make(map[string]bool)
	for _, pwd := range common {
		passwordMap[strings.ToLower(pwd)] = true
	}

	return passwordMap
}
