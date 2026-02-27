package validator

import (
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once         sync.Once
	validate     *validator.Validate
	initErr      error
	registrationStatus RegistrationStatus
)

// RegistrationStatus 验证器注册状态
type RegistrationStatus struct {
	Total      int                    // 总共需要注册的验证器数量
	Success    int                    // 成功注册的数量
	Failed     int                    // 注册失败的数量
	FailedTags []string               // 注册失败的验证器标签
	Errors     map[string]error       // 注册失败的错误信息 (tag -> error)
}

// IsComplete 检查是否所有验证器都注册成功
func (rs *RegistrationStatus) IsComplete() bool {
	return rs.Failed == 0 && rs.Total > 0
}

// GetFailedCount 获取注册失败的数量
func (rs *RegistrationStatus) GetFailedCount() int {
	return rs.Failed
}

// GetFailedTags 获取注册失败的验证器标签
func (rs *RegistrationStatus) GetFailedTags() []string {
	return rs.FailedTags
}

// GetValidator 获取全局验证器实例（单例模式）
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		// 注册自定义验证器并记录状态
		registrationStatus = RegisterCustomValidators(validate)
		if !registrationStatus.IsComplete() {
			log.Printf("[WARNING] Validator initialization completed with %d failures out of %d total validators",
				registrationStatus.Failed, registrationStatus.Total)
		}
	})
	return validate
}

// GetInitError 获取验证器初始化错误
func GetInitError() error {
	return initErr
}

// GetRegistrationStatus 获取验证器注册状态
func GetRegistrationStatus() RegistrationStatus {
	return registrationStatus
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
