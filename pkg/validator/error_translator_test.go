package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestTranslateError(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	type TestStruct struct {
		Username string  `validate:"required,username"`
		Email    string  `validate:"required,email"`
		Amount   float64 `validate:"required,positive_amount"`
		Phone    string  `validate:"omitempty,phone"`
	}

	tests := []struct {
		name          string
		data          TestStruct
		expectedCount int
		checkMessages []string
	}{
		{
			name:          "所有字段有效",
			data:          TestStruct{Username: "user123", Email: "test@example.com", Amount: 100.00, Phone: "13812345678"},
			expectedCount: 0,
		},
		{
			name:          "用户名无效",
			data:          TestStruct{Username: "ab", Email: "test@example.com", Amount: 100.00},
			expectedCount: 1,
			checkMessages: []string{"必须是3-20个字符"},
		},
		{
			name:          "邮箱无效",
			data:          TestStruct{Username: "user123", Email: "invalid", Amount: 100.00},
			expectedCount: 1,
			checkMessages: []string{"必须是有效的邮箱"},
		},
		{
			name:          "金额无效",
			data:          TestStruct{Username: "user123", Email: "test@example.com", Amount: -10.00},
			expectedCount: 1,
			checkMessages: []string{"必须是正数"},
		},
		{
			name:          "多个字段无效",
			data:          TestStruct{Username: "", Email: "", Amount: 0.00},
			expectedCount: 3,
			checkMessages: []string{"是必填字段", "是必填字段"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.data)
			errors := TranslateError(err)

			assert.Equal(t, tt.expectedCount, len(errors))

			if tt.expectedCount > 0 {
				for _, checkMsg := range tt.checkMessages {
					found := false
					for _, e := range errors {
						if containsString(e.Message, checkMsg) {
							found = true
							break
						}
					}
					assert.True(t, found, "期望找到包含 '%s' 的错误消息", checkMsg)
				}
			}
		})
	}
}

func TestValidationErrors_FormatErrors(t *testing.T) {
	errors := ValidationErrors{
		{Field: "username", Message: "用户名格式错误", Tag: "username"},
		{Field: "email", Message: "邮箱格式错误", Tag: "email"},
	}

	formatted := errors.FormatErrors()
	assert.Contains(t, formatted, "用户名格式错误")
	assert.Contains(t, formatted, "邮箱格式错误")
	assert.Contains(t, formatted, ";")
}

func TestValidationErrors_GetFieldErrors(t *testing.T) {
	errors := ValidationErrors{
		{Field: "username", Message: "用户名格式错误", Tag: "username"},
		{Field: "email", Message: "邮箱格式错误", Tag: "email"},
	}

	fieldErrors := errors.GetFieldErrors()
	assert.Equal(t, 2, len(fieldErrors))
	assert.Equal(t, "用户名格式错误", fieldErrors["username"])
	assert.Equal(t, "邮箱格式错误", fieldErrors["email"])
}

func TestGetErrorMessage_CustomValidators(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	tests := []struct {
		name        string
		validate    string
		value       interface{}
		expectedMsg string
	}{
		{
			name:        "positive_amount",
			validate:    "positive_amount",
			value:       -10.0,
			expectedMsg: "必须是正数",
		},
		{
			name:        "amount_range",
			validate:    "amount_range",
			value:       10000000.0,
			expectedMsg: "必须在 0.01 到 1,000,000.00 之间",
		},
		{
			name:        "username",
			validate:    "username",
			value:       "ab",
			expectedMsg: "必须是3-20个字符",
		},
		{
			name:        "phone",
			validate:    "phone",
			value:       "123",
			expectedMsg: "必须是有效的手机号",
		},
		{
			name:        "strong_password",
			validate:    "strong_password",
			value:       "weak",
			expectedMsg: "必须至少8位",
		},
		{
			name:        "withdraw_account",
			validate:    "withdraw_account",
			value:       "invalid",
			expectedMsg: "格式不正确",
		},
		{
			name:        "content_type",
			validate:    "content_type",
			value:       "invalid",
			expectedMsg: "必须是有效的内容类型",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Field interface{} `validate:"required"`
			}

			// 使用反射动态设置验证标签
			v := validator.New()
			RegisterCustomValidators(v)

			// 手动构造错误并验证消息
			err := v.Var(tt.value, tt.validate)
			if err != nil {
				errors := TranslateError(err)
				assert.NotEmpty(t, errors)
				if len(errors) > 0 {
					assert.Contains(t, errors[0].Message, tt.expectedMsg)
				}
			}
		})
	}
}

// 辅助函数
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
