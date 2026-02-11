package auth

import (
	"fmt"
	"regexp"
	"unicode"
)

// PasswordValidator 密码强度验证器
// MVP实现：基础规则验证
type PasswordValidator struct {
	minLength    int
	maxLength    int
	requireNum   bool
	requireAlpha bool
}

// NewPasswordValidator 创建密码验证器
// MVP: 使用默认规则
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:    8,  // 最短8位
		maxLength:    32, // 最长32位
		requireNum:   true,
		requireAlpha: true,
	}
}

// ValidatePassword 验证密码强度（MVP实现）
// 规则：
//   - 长度: 8-32位
//   - 至少包含一个数字
//   - 至少包含一个字母
func (v *PasswordValidator) ValidatePassword(password string) error {
	// 1. 检查长度
	if len(password) < v.minLength {
		return fmt.Errorf("密码长度不能少于%d位", v.minLength)
	}

	if len(password) > v.maxLength {
		return fmt.Errorf("密码长度不能超过%d位", v.maxLength)
	}

	// 2. 检查是否包含数字
	if v.requireNum {
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		if !hasNumber {
			return fmt.Errorf("密码必须包含至少一个数字")
		}
	}

	// 3. 检查是否包含字母
	if v.requireAlpha {
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
		if !hasLetter {
			return fmt.Errorf("密码必须包含至少一个字母")
		}
	}

	return nil
}

// --- Phase 2增强功能占位符 ---

// TODO: Phase 2增强功能
// - [ ] 检查是否包含特殊字符（!@#$%^&*等）
// - [ ] 检查是否包含大小写字母组合
// - [ ] 检查常见弱密码黑名单（123456, password等）
// - [ ] 检查与用户名相似度（不能太相似）
// - [ ] 密码历史记录（不能与最近5次密码相同）
// - [ ] 密码复杂度评分（弱/中/强）
// - [ ] 密码有效期检查
// - [ ] 支持自定义规则配置

// ValidatePasswordStrength 密码强度评分（未来实现）
// 返回：weak, medium, strong
func (v *PasswordValidator) ValidatePasswordStrength(password string) string {
	score := 0

	// 基础长度分
	if len(password) >= 8 {
		score += 1
	}
	if len(password) >= 12 {
		score += 1
	}

	// 字符类型分
	hasNumber := false
	hasLower := false
	hasUpper := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if hasNumber {
		score += 1
	}
	if hasLower && hasUpper {
		score += 2
	} else if hasLower || hasUpper {
		score += 1
	}
	if hasSpecial {
		score += 2
	}

	// 评级
	if score <= 2 {
		return "weak"
	} else if score <= 5 {
		return "medium"
	} else {
		return "strong"
	}
}

// CheckWeakPassword 检查是否为常见弱密码（未来实现）
func (v *PasswordValidator) CheckWeakPassword(password string) bool {
	// 常见弱密码列表
	weakPasswords := []string{
		"123456", "password", "123456789", "12345678", "12345",
		"1234567", "password1", "123123", "1234567890", "1234",
		"qwerty", "abc123", "111111", "admin", "root",
	}

	for _, weak := range weakPasswords {
		if password == weak {
			return true
		}
	}

	return false
}

// GetPasswordRequirements 获取密码要求说明（用于前端提示）
func (v *PasswordValidator) GetPasswordRequirements() string {
	return fmt.Sprintf(
		"密码要求：长度%d-%d位，至少包含一个字母和一个数字",
		v.minLength,
		v.maxLength,
	)
}

// --- 配置方法 ---

// SetMinLength 设置最小长度
func (v *PasswordValidator) SetMinLength(length int) *PasswordValidator {
	v.minLength = length
	return v
}

// SetMaxLength 设置最大长度
func (v *PasswordValidator) SetMaxLength(length int) *PasswordValidator {
	v.maxLength = length
	return v
}

// SetRequireNumber 设置是否要求数字
func (v *PasswordValidator) SetRequireNumber(require bool) *PasswordValidator {
	v.requireNum = require
	return v
}

// SetRequireAlpha 设置是否要求字母
func (v *PasswordValidator) SetRequireAlpha(require bool) *PasswordValidator {
	v.requireAlpha = require
	return v
}
