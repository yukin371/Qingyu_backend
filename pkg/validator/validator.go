package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// GetValidator 获取全局验证器实例（单例模式）
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		// 注册自定义验证器
		RegisterCustomValidators(validate)
	})
	return validate
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

// ValidateStructWithErrors 验证结构体并返回友好错误
func ValidateStructWithErrors(s interface{}) ValidationErrors {
	err := ValidateStruct(s)
	if err == nil {
		return nil
	}
	return TranslateError(err)
}
