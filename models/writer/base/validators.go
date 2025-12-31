package base

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	ErrIDRequired          = errors.New("ID不能为空")
	ErrNameRequired        = errors.New("名称不能为空")
	ErrNameTooLong         = errors.New("名称过长")
	ErrTitleRequired       = errors.New("标题不能为空")
	ErrTitleTooLong        = errors.New("标题过长")
	ErrProjectIDRequired   = errors.New("项目ID不能为空")
	ErrAuthorIDRequired    = errors.New("作者ID不能为空")
	ErrInvalidEmail        = errors.New("邮箱格式无效")
	ErrInvalidURL          = errors.New("URL格式无效")
	ErrInvalidStatus       = errors.New("无效的状态")
	ErrInvalidVisibility   = errors.New("无效的可见性")
	ErrInvalidWritingType  = errors.New("无效的写作类型")
	ErrInvalidDocumentType = errors.New("无效的文档类型")
	ErrInvalidLevel        = errors.New("无效的层级")
	ErrInvalidParentType   = errors.New("无效的父类型")
	ErrInvalidRootLevel    = errors.New("根节点层级必须为0")
	ErrInvalidVersion      = errors.New("无效的版本号")
	ErrInvalidEnum         = errors.New("无效的枚举值")
)

// ValidateName 验证名称
func ValidateName(name string, maxLength int) error {
	if name == "" {
		return ErrNameRequired
	}
	if utf8.RuneCountInString(name) > maxLength {
		return ErrNameTooLong
	}
	return nil
}

// ValidateTitle 验证标题
func ValidateTitle(title string, maxLength int) error {
	if title == "" {
		return ErrTitleRequired
	}
	if utf8.RuneCountInString(title) > maxLength {
		return ErrTitleTooLong
	}
	return nil
}

// ValidateURL 验证URL格式
func ValidateURL(url string) error {
	if url == "" {
		return nil // 空URL允许
	}
	// 简单的URL验证
	urlPattern := regexp.MustCompile(`^https?://[^\s]+$`)
	if !urlPattern.MatchString(url) {
		return ErrInvalidURL
	}
	return nil
}

// ValidateStringLength 验证字符串长度
func ValidateStringLength(s string, min, max int) error {
	length := utf8.RuneCountInString(s)
	if length < min {
		return errors.New("字符串过短")
	}
	if max > 0 && length > max {
		return errors.New("字符串过长")
	}
	return nil
}

// ValidateRequired 验证必填字段
func ValidateRequired(field, fieldName string) error {
	if field == "" {
		return errors.New(fieldName + "不能为空")
	}
	return nil
}

// ValidateEnum 验证枚举值
func ValidateEnum(value string, validValues map[string]bool) error {
	if !validValues[value] {
		return errors.New("无效的枚举值")
	}
	return nil
}
