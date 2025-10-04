package errors

import (
	"fmt"
	"time"
)

// RepositoryErrorType 仓储错误类型
type RepositoryErrorType string

const (
	RepositoryErrorValidation RepositoryErrorType = "validation" // 验证错误
	RepositoryErrorNotFound   RepositoryErrorType = "not_found"  // 未找到错误
	RepositoryErrorConflict   RepositoryErrorType = "conflict"   // 冲突错误
	RepositoryErrorDuplicate  RepositoryErrorType = "duplicate"  // 重复错误
	RepositoryErrorInternal   RepositoryErrorType = "internal"   // 内部错误
)

// RepositoryError 仓储错误
type RepositoryError struct {
	Type    RepositoryErrorType
	Message string
	Detail  error
}

// Error 实现error接口
func (e *RepositoryError) Error() string {
	return e.Message
}

// NewRepositoryError 创建仓储错误
func NewRepositoryError(t RepositoryErrorType, msg string, detail error) *RepositoryError {
	return &RepositoryError{
		Type:    t,
		Message: msg,
		Detail:  detail,
	}
}

// ServiceErrorType 服务错误类型
type ServiceErrorType string

const (
	ServiceErrorValidation   ServiceErrorType = "VALIDATION"
	ServiceErrorBusiness     ServiceErrorType = "BUSINESS"
	ServiceErrorNotFound     ServiceErrorType = "NOT_FOUND"
	ServiceErrorUnauthorized ServiceErrorType = "UNAUTHORIZED"
	ServiceErrorForbidden    ServiceErrorType = "FORBIDDEN"
	ServiceErrorInternal     ServiceErrorType = "INTERNAL"
	ServiceErrorTimeout      ServiceErrorType = "TIMEOUT"
	ServiceErrorExternal     ServiceErrorType = "EXTERNAL"
)

// ServiceError 服务错误
type ServiceError struct {
	Type      ServiceErrorType
	Message   string
	Details   string
	Cause     error
	Service   string
	Timestamp time.Time
}

// Error 实现error接口
func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Service, e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Service, e.Type, e.Message)
}

// NewServiceError 创建服务错误
func NewServiceError(service string, errorType ServiceErrorType, message string, details string, cause error) *ServiceError {
	return &ServiceError{
		Type:      errorType,
		Message:   message,
		Details:   details,
		Cause:     cause,
		Service:   service,
		Timestamp: time.Now(),
	}
}
