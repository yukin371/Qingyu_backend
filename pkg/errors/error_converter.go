package errors

import (
	"fmt"
	"time"

	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/base"
)

// ErrorConverter 错误转换器
type ErrorConverter struct{}

// NewErrorConverter 创建错误转换器
func NewErrorConverter() *ErrorConverter {
	return &ErrorConverter{}
}

// ConvertRepositoryError 转换Repository错误
func (c *ErrorConverter) ConvertRepositoryError(err *infrastructure.RepositoryError, service string) *UnifiedError {
	if err == nil {
		return nil
	}

	var category ErrorCategory
	var httpStatus int
	var level ErrorLevel

	switch err.Type {
	case infrastructure.ValidationError:
		category = CategoryValidation
		httpStatus = 400
		level = LevelWarning
	case infrastructure.NotFoundError:
		category = CategoryBusiness
		httpStatus = 404
		level = LevelWarning
	case infrastructure.DuplicateError:
		category = CategoryBusiness
		httpStatus = 409
		level = LevelWarning
	case infrastructure.InternalError:
		category = CategorySystem
		httpStatus = 500
		level = LevelError
	default:
		category = CategorySystem
		httpStatus = 500
		level = LevelError
	}

	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       string(err.Type),
		Category:   category,
		Level:      level,
		Message:    err.Message,
		Details:    err.Detail,
		Service:    service,
		Timestamp:  time.Now(),
		HTTPStatus: httpStatus,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// ConvertServiceError 转换Service错误
func (c *ErrorConverter) ConvertServiceError(err *base.ServiceError, service string) *UnifiedError {
	if err == nil {
		return nil
	}

	var category ErrorCategory
	var httpStatus int
	var level ErrorLevel

	switch err.Type {
	case base.VALIDATION:
		category = CategoryValidation
		httpStatus = 400
		level = LevelWarning
	case base.BUSINESS:
		category = CategoryBusiness
		httpStatus = 409
		level = LevelError
	case base.NOT_FOUND:
		category = CategoryBusiness
		httpStatus = 404
		level = LevelWarning
	case base.INTERNAL:
		category = CategorySystem
		httpStatus = 500
		level = LevelError
	default:
		category = CategorySystem
		httpStatus = 500
		level = LevelError
	}

	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       string(err.Type),
		Category:   category,
		Level:      level,
		Message:    err.Message,
		Details:    err.Details,
		Service:    service,
		Cause:      err.Cause,
		Timestamp:  time.Now(),
		HTTPStatus: httpStatus,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// ConvertGenericError 转换通用错误
func (c *ErrorConverter) ConvertGenericError(err error, service, operation string) *UnifiedError {
	if err == nil {
		return nil
	}

	// 检查是否已经是UnifiedError
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	// 检查是否是RepositoryError
	if repoErr, ok := err.(*infrastructure.RepositoryError); ok {
		return c.ConvertRepositoryError(repoErr, service)
	}

	// 检查是否是ServiceError
	if serviceErr, ok := err.(*base.ServiceError); ok {
		return c.ConvertServiceError(serviceErr, service)
	}

	// 默认转换为内部错误
	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "INTERNAL_ERROR",
		Category:   CategorySystem,
		Level:      LevelError,
		Message:    "Internal server error",
		Details:    err.Error(),
		Service:    service,
		Operation:  operation,
		Cause:      err,
		Timestamp:  time.Now(),
		HTTPStatus: 500,
		Retryable:  false,
		Metadata:   make(map[string]interface{}),
	}
}

// ConvertPanic 转换panic为错误
func (c *ErrorConverter) ConvertPanic(panicValue interface{}, service, operation string) *UnifiedError {
	var message string
	var details string

	switch v := panicValue.(type) {
	case string:
		message = "Panic occurred"
		details = v
	case error:
		message = "Panic occurred"
		details = v.Error()
	default:
		message = "Panic occurred"
		details = fmt.Sprintf("Unknown panic type: %T", v)
	}

	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "PANIC",
		Category:   CategorySystem,
		Level:      LevelCritical,
		Message:    message,
		Details:    details,
		Service:    service,
		Operation:  operation,
		Timestamp:  time.Now(),
		HTTPStatus: 500,
		Retryable:  false,
		Metadata: map[string]interface{}{
			"panic_value": panicValue,
		},
	}
}

// 全局转换器实例
var DefaultConverter = NewErrorConverter()
