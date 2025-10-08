package errors

import (
	"fmt"
	"time"
)

// ErrorFactory 错误工厂
type ErrorFactory struct {
	service string
}

// NewErrorFactory 创建错误工厂
func NewErrorFactory(service string) *ErrorFactory {
	return &ErrorFactory{service: service}
}

// ValidationError 创建验证错误
func (f *ErrorFactory) ValidationError(code, message string, details ...string) *UnifiedError {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       code,
		Category:   CategoryValidation,
		Level:      LevelWarning,
		Message:    message,
		Details:    detail,
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 400,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// BusinessError 创建业务错误
func (f *ErrorFactory) BusinessError(code, message string, details ...string) *UnifiedError {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       code,
		Category:   CategoryBusiness,
		Level:      LevelError,
		Message:    message,
		Details:    detail,
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 409,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// NotFoundError 创建未找到错误
func (f *ErrorFactory) NotFoundError(resource, id string) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "NOT_FOUND",
		Category:   CategoryBusiness,
		Level:      LevelWarning,
		Message:    fmt.Sprintf("%s not found", resource),
		Details:    fmt.Sprintf("Resource %s with ID %s was not found", resource, id),
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 404,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// AuthError 创建认证错误
func (f *ErrorFactory) AuthError(code, message string) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       code,
		Category:   CategoryAuth,
		Level:      LevelWarning,
		Message:    message,
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 401,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// InternalError 创建内部错误
func (f *ErrorFactory) InternalError(code, message string, cause error) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       code,
		Category:   CategorySystem,
		Level:      LevelError,
		Message:    message,
		Service:    f.service,
		Cause:      cause,
		Timestamp:  time.Now(),
		HTTPStatus: 500,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// ExternalError 创建外部服务错误
func (f *ErrorFactory) ExternalError(code, message string, retryable bool) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       code,
		Category:   CategoryExternal,
		Level:      LevelError,
		Message:    message,
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 502,
		Retryable:  retryable,
		Metadata:   make(map[string]interface{}),
	}
}

// NetworkError 创建网络错误
func (f *ErrorFactory) NetworkError(message string) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "NETWORK_ERROR",
		Category:   CategoryNetwork,
		Level:      LevelError,
		Message:    message,
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 503,
		Retryable:  true,
		Metadata:   make(map[string]interface{}),
	}
}

// TimeoutError 创建超时错误
func (f *ErrorFactory) TimeoutError(operation string) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "TIMEOUT",
		Category:   CategoryNetwork,
		Level:      LevelError,
		Message:    fmt.Sprintf("Operation %s timed out", operation),
		Service:    f.service,
		Operation:  operation,
		Timestamp:  time.Now(),
		HTTPStatus: 408,
		Retryable:  true,
		Metadata:   make(map[string]interface{}),
	}
}

// RateLimitError 创建限流错误
func (f *ErrorFactory) RateLimitError(limit int) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "RATE_LIMIT_EXCEEDED",
		Category:   CategorySystem,
		Level:      LevelWarning,
		Message:    "Rate limit exceeded",
		Details:    fmt.Sprintf("Request rate limit of %d per minute exceeded", limit),
		Service:    f.service,
		Timestamp:  time.Now(),
		HTTPStatus: 429,
		Retryable:  true,
		Metadata:   make(map[string]interface{}),
	}
}

// DatabaseError 创建数据库错误
func (f *ErrorFactory) DatabaseError(operation string, cause error) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "DATABASE_ERROR",
		Category:   CategoryDatabase,
		Level:      LevelError,
		Message:    fmt.Sprintf("Database operation failed: %s", operation),
		Service:    f.service,
		Operation:  operation,
		Cause:      cause,
		Timestamp:  time.Now(),
		HTTPStatus: 500,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// CacheError 创建缓存错误
func (f *ErrorFactory) CacheError(operation string, cause error) *UnifiedError {
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "CACHE_ERROR",
		Category:   CategoryCache,
		Level:      LevelWarning,
		Message:    fmt.Sprintf("Cache operation failed: %s", operation),
		Service:    f.service,
		Operation:  operation,
		Cause:      cause,
		Timestamp:  time.Now(),
		HTTPStatus: 500,
		Retryable:  true,
		Metadata:   make(map[string]interface{}),
	}
}

// 全局错误工厂实例
var (
	AIServiceFactory        = NewErrorFactory("ai-service")
	UserServiceFactory      = NewErrorFactory("user-service")
	DocumentServiceFactory  = NewErrorFactory("document-service")
	ProjectFactory          = NewErrorFactory("project-service")
	BookstoreServiceFactory = NewErrorFactory("bookstore-service")
	ReaderServiceFactory    = NewErrorFactory("reader-service")
	WriterServiceFactory    = NewErrorFactory("writer-service")
)

// generateErrorID 生成错误ID
func generateErrorID() string {
	return fmt.Sprintf("err_%d", time.Now().UnixNano())
}
