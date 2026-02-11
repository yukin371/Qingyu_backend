package user

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 错误码常量
const (
	ErrUserNotFound      ErrorCode = 40401 // 用户不存在
	ErrInvalidEmail      ErrorCode = 40001 // 邮箱格式无效
	ErrInvalidPassword   ErrorCode = 40002 // 密码格式无效
	ErrUserAlreadyExists ErrorCode = 40901 // 用户已存在
	ErrTokenInvalid      ErrorCode = 40101 // 令牌无效
	ErrTokenExpired      ErrorCode = 40102 // 令牌过期
	ErrPermissionDenied  ErrorCode = 40301 // 权限不足
	ErrInternalError     ErrorCode = 50001 // 内部错误
)

// HTTPStatus 返回错误码对应的 HTTP 状态码
func (c ErrorCode) HTTPStatus() int {
	switch {
	case c >= 40000 && c < 40100:
		return 400
	case c >= 40100 && c < 40200:
		return 401
	case c >= 40300 && c < 40400:
		return 403
	case c >= 40400 && c < 40500:
		return 404
	case c >= 40900 && c < 41000:
		return 409
	case c >= 50000:
		return 500
	default:
		return http.StatusInternalServerError
	}
}

// String 返回错误码的字符串表示
func (c ErrorCode) String() string {
	switch c {
	case ErrUserNotFound:
		return "USER_NOT_FOUND"
	case ErrInvalidEmail:
		return "INVALID_EMAIL"
	case ErrInvalidPassword:
		return "INVALID_PASSWORD"
	case ErrUserAlreadyExists:
		return "USER_ALREADY_EXISTS"
	case ErrTokenInvalid:
		return "TOKEN_INVALID"
	case ErrTokenExpired:
		return "TOKEN_EXPIRED"
	case ErrPermissionDenied:
		return "PERMISSION_DENIED"
	case ErrInternalError:
		return "INTERNAL_ERROR"
	default:
		return fmt.Sprintf("UNKNOWN_ERROR(%d)", c)
	}
}

// UserError 用户模块结构化错误
type UserError struct {
	Code     ErrorCode
	Field    string
	Message  string
	Err      error
	Metadata map[string]interface{}
}

// Error 实现 error 接口
func (e *UserError) Error() string {
	if e.Message != "" {
		if e.Field != "" {
			return fmt.Sprintf("[%s] %s: %s", e.Code.String(), e.Field, e.Message)
		}
		return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
	}
	return e.Code.String()
}

// Unwrap 支持 errors.Unwrap
func (e *UserError) Unwrap() error {
	return e.Err
}

// WithField 设置字段名
func (e *UserError) WithField(field string) *UserError {
	e.Field = field
	return e
}

// WithCause 设置底层错误
func (e *UserError) WithCause(err error) *UserError {
	e.Err = err
	return e
}

// WithMetadata 设置元数据
func (e *UserError) WithMetadata(metadata map[string]interface{}) *UserError {
	e.Metadata = metadata
	return e
}

// IsRetryable 判断错误是否可重试
func (e *UserError) IsRetryable() bool {
	// 用户错误通常不可重试
	return false
}

// NewUserError 创建新的用户错误
func NewUserError(code ErrorCode, message string) *UserError {
	return &UserError{
		Code:    code,
		Message: message,
	}
}

// IsUserError 判断是否是 UserError 类型
func IsUserError(err error) bool {
	_, ok := err.(*UserError)
	return ok
}

// IsErrorCode 判断错误是否是指定错误码
func IsErrorCode(err error, code ErrorCode) bool {
	if userErr, ok := err.(*UserError); ok {
		return userErr.Code == code
	}
	return false
}

// 专用构造函数

// NotFound 创建用户不存在错误
func NotFound(identifier string) *UserError {
	message := "用户不存在"
	if identifier != "" {
		message = fmt.Sprintf("用户 %s 不存在", identifier)
	}
	return &UserError{Code: ErrUserNotFound, Message: message}
}

// InvalidEmail 创建邮箱格式无效错误
func InvalidEmail(field string) *UserError {
	message := "邮箱格式无效"
	if field != "" {
		message = fmt.Sprintf("%s 邮箱格式无效", field)
	}
	return &UserError{Code: ErrInvalidEmail, Field: "email", Message: message}
}

// InvalidPassword 创建密码格式无效错误
func InvalidPassword(message string) *UserError {
	msg := message
	if msg == "" {
		msg = "密码格式无效"
	}
	return &UserError{Code: ErrInvalidPassword, Field: "password", Message: msg}
}

// UserAlreadyExists 创建用户已存在错误
func UserAlreadyExists(field string) *UserError {
	var message string
	var fieldName string

	switch field {
	case "username":
		message = "用户名已存在"
		fieldName = "username"
	case "email":
		message = "邮箱已被使用"
		fieldName = "email"
	default:
		message = "用户已存在"
		fieldName = field
	}

	return &UserError{Code: ErrUserAlreadyExists, Field: fieldName, Message: message}
}

// TokenInvalid 创建令牌无效错误
func TokenInvalid() *UserError {
	return &UserError{Code: ErrTokenInvalid, Message: "令牌无效"}
}

// TokenExpired 创建令牌过期错误
func TokenExpired() *UserError {
	return &UserError{Code: ErrTokenExpired, Message: "令牌已过期"}
}

// PermissionDenied 创建权限不足错误
func PermissionDenied(action string) *UserError {
	message := "权限不足"
	if action != "" {
		message = fmt.Sprintf("没有%s的权限", action)
	}
	return &UserError{Code: ErrPermissionDenied, Message: message}
}

// InternalError 创建内部错误
func InternalError(message string, cause ...error) *UserError {
	err := &UserError{
		Code:    ErrInternalError,
		Message: message,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}
