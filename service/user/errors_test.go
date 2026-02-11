package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ========== ErrorCode 类型测试 ==========

func TestErrorCode(t *testing.T) {
	t.Run("错误码应该正确定义", func(t *testing.T) {
		// 测试 ErrorCode 类型定义
		assert.Equal(t, ErrorCode(40401), ErrUserNotFound)
		assert.Equal(t, ErrorCode(40001), ErrInvalidEmail)
		assert.Equal(t, ErrorCode(40002), ErrInvalidPassword)
		assert.Equal(t, ErrorCode(40901), ErrUserAlreadyExists)
		assert.Equal(t, ErrorCode(40101), ErrTokenInvalid)
		assert.Equal(t, ErrorCode(40102), ErrTokenExpired)
		assert.Equal(t, ErrorCode(40301), ErrPermissionDenied)
		assert.Equal(t, ErrorCode(50001), ErrInternalError)
	})

	t.Run("用户不存在错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40401), ErrUserNotFound)
	})

	t.Run("邮箱格式无效错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40001), ErrInvalidEmail)
	})

	t.Run("密码格式无效错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40002), ErrInvalidPassword)
	})

	t.Run("用户已存在错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40901), ErrUserAlreadyExists)
	})

	t.Run("令牌无效错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40101), ErrTokenInvalid)
	})

	t.Run("令牌过期错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40102), ErrTokenExpired)
	})

	t.Run("权限不足错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(40301), ErrPermissionDenied)
	})

	t.Run("内部错误错误码", func(t *testing.T) {
		assert.Equal(t, ErrorCode(50001), ErrInternalError)
	})

	t.Run("错误码的 HTTP 状态码映射", func(t *testing.T) {
		// 测试 404 系列错误码
		assert.Equal(t, 404, ErrUserNotFound.HTTPStatus())

		// 测试 400 系列错误码
		assert.Equal(t, 400, ErrInvalidEmail.HTTPStatus())
		assert.Equal(t, 400, ErrInvalidPassword.HTTPStatus())

		// 测试 409 系列错误码
		assert.Equal(t, 409, ErrUserAlreadyExists.HTTPStatus())

		// 测试 401 系列错误码
		assert.Equal(t, 401, ErrTokenInvalid.HTTPStatus())
		assert.Equal(t, 401, ErrTokenExpired.HTTPStatus())

		// 测试 403 系列错误码
		assert.Equal(t, 403, ErrPermissionDenied.HTTPStatus())

		// 测试 500 系列错误码
		assert.Equal(t, 500, ErrInternalError.HTTPStatus())
	})

	t.Run("错误码的字符串表示", func(t *testing.T) {
		tests := []struct {
			code     ErrorCode
			expected string
		}{
			{ErrUserNotFound, "USER_NOT_FOUND"},
			{ErrInvalidEmail, "INVALID_EMAIL"},
			{ErrInvalidPassword, "INVALID_PASSWORD"},
			{ErrUserAlreadyExists, "USER_ALREADY_EXISTS"},
			{ErrTokenInvalid, "TOKEN_INVALID"},
			{ErrTokenExpired, "TOKEN_EXPIRED"},
			{ErrPermissionDenied, "PERMISSION_DENIED"},
			{ErrInternalError, "INTERNAL_ERROR"},
		}

		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.code.String())
			})
		}
	})
}

// ========== UserError 类型测试 ==========

func TestUserError(t *testing.T) {
	t.Run("UserError 结构体定义", func(t *testing.T) {
		// 测试基本字段
		err := &UserError{
			Code:    ErrUserNotFound,
			Field:   "user_id",
			Message: "用户不存在",
			Err:     nil,
		}

		assert.Equal(t, ErrorCode(40401), err.Code)
		assert.Equal(t, "user_id", err.Field)
		assert.Equal(t, "用户不存在", err.Message)
		assert.Nil(t, err.Err)
	})

	t.Run("UserError 实现 error 接口", func(t *testing.T) {
		err := &UserError{
			Code:    ErrUserNotFound,
			Message: "用户不存在",
		}

		expected := "[USER_NOT_FOUND] 用户不存在"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("UserError 带字段的错误消息", func(t *testing.T) {
		err := &UserError{
			Code:    ErrInvalidEmail,
			Field:   "email",
			Message: "邮箱格式不正确",
		}

		expected := "[INVALID_EMAIL] email: 邮箱格式不正确"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("UserError 实现 Unwrap 接口", func(t *testing.T) {
		internalErr := errors.New("database connection failed")
		err := &UserError{
			Code:    ErrInternalError,
			Field:   "database",
			Message: "数据库连接失败",
			Err:     internalErr,
		}

		assert.Equal(t, internalErr, errors.Unwrap(err))
	})

	t.Run("UserError 不包含底层错误时 Unwrap 返回 nil", func(t *testing.T) {
		err := &UserError{
			Code:    ErrUserNotFound,
			Message: "用户不存在",
		}

		assert.Nil(t, errors.Unwrap(err))
	})
}

// ========== 构造函数测试 ==========

func TestNewUserError(t *testing.T) {
	t.Run("创建基本错误", func(t *testing.T) {
		err := NewUserError(ErrUserNotFound, "用户不存在")

		assert.Equal(t, ErrUserNotFound, err.Code)
		assert.Equal(t, "用户不存在", err.Message)
		assert.Empty(t, err.Field)
		assert.Nil(t, err.Err)
	})

	t.Run("创建带字段的错误", func(t *testing.T) {
		err := NewUserError(ErrInvalidEmail, "邮箱格式不正确").
			WithField("email")

		assert.Equal(t, ErrInvalidEmail, err.Code)
		assert.Equal(t, "邮箱格式不正确", err.Message)
		assert.Equal(t, "email", err.Field)
	})

	t.Run("创建带底层错误的错误", func(t *testing.T) {
		internalErr := errors.New("internal error")
		err := NewUserError(ErrInternalError, "内部错误").
			WithCause(internalErr)

		assert.Equal(t, ErrInternalError, err.Code)
		assert.Equal(t, "internal error", errors.Unwrap(err).Error())
	})

	t.Run("链式调用", func(t *testing.T) {
		internalErr := errors.New("db error")
		err := NewUserError(ErrInternalError, "数据库错误").
			WithField("password").
			WithCause(internalErr)

		assert.Equal(t, ErrInternalError, err.Code)
		assert.Equal(t, "数据库错误", err.Message)
		assert.Equal(t, "password", err.Field)
		assert.Equal(t, internalErr, errors.Unwrap(err))
	})
}

func TestNotFound(t *testing.T) {
	t.Run("创建用户不存在错误", func(t *testing.T) {
		err := NotFound("user_123")

		assert.Equal(t, ErrUserNotFound, err.Code)
		assert.Contains(t, err.Message, "user_123")
	})

	t.Run("创建带字段的不存在错误", func(t *testing.T) {
		err := NotFound("email@example.com").WithField("email")

		assert.Equal(t, ErrUserNotFound, err.Code)
		assert.Equal(t, "email", err.Field)
	})
}

func TestInvalidEmail(t *testing.T) {
	t.Run("创建邮箱格式无效错误", func(t *testing.T) {
		err := InvalidEmail("invalid-email")

		assert.Equal(t, ErrInvalidEmail, err.Code)
		assert.Equal(t, "email", err.Field)
		assert.Contains(t, err.Message, "invalid-email")
	})

	t.Run("使用默认消息", func(t *testing.T) {
		err := InvalidEmail("")

		assert.Equal(t, ErrInvalidEmail, err.Code)
		assert.Equal(t, "email", err.Field)
		assert.Contains(t, err.Message, "邮箱格式")
	})
}

func TestInvalidPassword(t *testing.T) {
	t.Run("创建密码格式无效错误", func(t *testing.T) {
		err := InvalidPassword("密码太短")

		assert.Equal(t, ErrInvalidPassword, err.Code)
		assert.Equal(t, "password", err.Field)
		assert.Equal(t, "密码太短", err.Message)
	})

	t.Run("使用默认消息", func(t *testing.T) {
		err := InvalidPassword("")

		assert.Equal(t, ErrInvalidPassword, err.Code)
		assert.Equal(t, "password", err.Field)
		assert.Contains(t, err.Message, "密码")
	})
}

func TestUserAlreadyExists(t *testing.T) {
	t.Run("创建用户已存在错误 - 用户名", func(t *testing.T) {
		err := UserAlreadyExists("username")

		assert.Equal(t, ErrUserAlreadyExists, err.Code)
		assert.Equal(t, "username", err.Field)
		assert.Contains(t, err.Message, "用户名")
	})

	t.Run("创建用户已存在错误 - 邮箱", func(t *testing.T) {
		err := UserAlreadyExists("email")

		assert.Equal(t, ErrUserAlreadyExists, err.Code)
		assert.Equal(t, "email", err.Field)
		assert.Contains(t, err.Message, "邮箱")
	})
}

func TestTokenInvalid(t *testing.T) {
	t.Run("创建令牌无效错误", func(t *testing.T) {
		err := TokenInvalid()

		assert.Equal(t, ErrTokenInvalid, err.Code)
		assert.Contains(t, err.Message, "令牌")
		assert.Contains(t, err.Message, "无效")
	})
}

func TestTokenExpired(t *testing.T) {
	t.Run("创建令牌过期错误", func(t *testing.T) {
		err := TokenExpired()

		assert.Equal(t, ErrTokenExpired, err.Code)
		assert.Contains(t, err.Message, "令牌")
		assert.Contains(t, err.Message, "过期")
	})
}

func TestPermissionDenied(t *testing.T) {
	t.Run("创建权限不足错误", func(t *testing.T) {
		err := PermissionDenied("删除用户")

		assert.Equal(t, ErrPermissionDenied, err.Code)
		assert.Contains(t, err.Message, "删除用户")
		assert.Contains(t, err.Message, "权限")
	})
}

func TestInternalError(t *testing.T) {
	t.Run("创建内部错误 - 不带原因", func(t *testing.T) {
		err := InternalError("服务异常")

		assert.Equal(t, ErrInternalError, err.Code)
		assert.Equal(t, "服务异常", err.Message)
		assert.Nil(t, errors.Unwrap(err))
	})

	t.Run("创建内部错误 - 带原因", func(t *testing.T) {
		cause := errors.New("database connection failed")
		err := InternalError("数据库连接失败", cause)

		assert.Equal(t, ErrInternalError, err.Code)
		assert.Equal(t, "数据库连接失败", err.Message)
		assert.Equal(t, cause, errors.Unwrap(err))
	})
}

// ========== 错误类型判断测试 ==========

func TestIsUserError(t *testing.T) {
	t.Run("识别 UserError 类型", func(t *testing.T) {
		userErr := &UserError{Code: ErrUserNotFound}
		stdErr := errors.New("standard error")

		assert.True(t, IsUserError(userErr))
		assert.False(t, IsUserError(stdErr))
		assert.False(t, IsUserError(nil))
	})

	t.Run("识别具体的错误码", func(t *testing.T) {
		tests := []struct {
			name     string
			err      error
			code     ErrorCode
			expected bool
		}{
			{"用户不存在", &UserError{Code: ErrUserNotFound}, ErrUserNotFound, true},
			{"用户不存在与其他错误", &UserError{Code: ErrUserNotFound}, ErrInvalidEmail, false},
			{"标准错误", errors.New("standard"), ErrUserNotFound, false},
			{"nil 错误", nil, ErrUserNotFound, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, IsErrorCode(tt.err, tt.code))
			})
		}
	})
}

// ========== HTTP 状态码映射测试 ==========

func TestHTTPStatusMapping(t *testing.T) {
	t.Run("所有错误码都有对应的 HTTP 状态码", func(t *testing.T) {
		tests := []struct {
			code           ErrorCode
			expectedStatus int
		}{
			{ErrUserNotFound, 404},
			{ErrInvalidEmail, 400},
			{ErrInvalidPassword, 400},
			{ErrUserAlreadyExists, 409},
			{ErrTokenInvalid, 401},
			{ErrTokenExpired, 401},
			{ErrPermissionDenied, 403},
			{ErrInternalError, 500},
		}

		for _, tt := range tests {
			t.Run(tt.code.String(), func(t *testing.T) {
				assert.Equal(t, tt.expectedStatus, tt.code.HTTPStatus())
			})
		}
	})
}

// ========== 集成测试 ==========

func TestUserErrorIntegration(t *testing.T) {
	t.Run("完整的错误链", func(t *testing.T) {
		// 模拟一个完整的错误场景
		dbErr := errors.New("connection timeout")

		userErr := InternalError("数据库操作失败", dbErr).
			WithField("password").
			WithMetadata(map[string]interface{}{
				"timeout":   30,
				"operation": "update_password",
			})

		// 验证错误链
		assert.Equal(t, ErrInternalError, userErr.Code)
		assert.Equal(t, "password", userErr.Field)
		assert.Equal(t, dbErr, errors.Unwrap(userErr))

		// 验证错误消息包含关键信息
		errorMsg := userErr.Error()
		assert.Contains(t, errorMsg, "INTERNAL_ERROR")
		assert.Contains(t, errorMsg, "数据库操作失败")
	})

	t.Run("错误的可重试性", func(t *testing.T) {
		tests := []struct {
			name     string
			err      *UserError
			expected bool
		}{
			{"用户不可重试", &UserError{Code: ErrUserNotFound}, false},
			{"验证错误不可重试", &UserError{Code: ErrInvalidEmail}, false},
			{"权限错误不可重试", &UserError{Code: ErrPermissionDenied}, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.err.IsRetryable())
			})
		}
	})
}

// ========== 边界情况测试 ==========

func TestUserErrorEdgeCases(t *testing.T) {
	t.Run("空字段名", func(t *testing.T) {
		err := &UserError{
			Code:    ErrUserNotFound,
			Field:   "",
			Message: "用户不存在",
		}

		// 空字段时不应该在错误消息中显示字段
		errorMsg := err.Error()
		assert.NotContains(t, errorMsg, ":")
	})

	t.Run("空消息", func(t *testing.T) {
		err := &UserError{
			Code:    ErrUserNotFound,
			Message: "",
		}

		// 空消息时应该使用错误码的默认消息
		errorMsg := err.Error()
		assert.NotEmpty(t, errorMsg)
	})

	t.Run("nil 原因错误", func(t *testing.T) {
		err := &UserError{
			Code:    ErrUserNotFound,
			Message: "用户不存在",
			Err:     nil,
		}

		assert.Nil(t, errors.Unwrap(err))
	})
}

// ========== 性能测试 ==========

func BenchmarkUserErrorCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NotFound("user_123")
	}
}

func BenchmarkUserErrorWithCause(b *testing.B) {
	internalErr := errors.New("internal error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InternalError("error", internalErr)
	}
}
