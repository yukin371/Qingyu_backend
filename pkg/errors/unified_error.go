package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// ErrorCategory 错误分类
type ErrorCategory string

const (
	CategoryValidation ErrorCategory = "validation"
	CategoryBusiness   ErrorCategory = "business"
	CategorySystem     ErrorCategory = "system"
	CategoryExternal   ErrorCategory = "external"
	CategoryNetwork    ErrorCategory = "network"
	CategoryAuth       ErrorCategory = "auth"
	CategoryDatabase   ErrorCategory = "database"
	CategoryCache      ErrorCategory = "cache"
)

// ErrorLevel 错误级别
type ErrorLevel string

const (
	LevelInfo     ErrorLevel = "info"
	LevelWarning  ErrorLevel = "warning"
	LevelError    ErrorLevel = "error"
	LevelCritical ErrorLevel = "critical"
)

// UnifiedError 统一错误结构
type UnifiedError struct {
	// 基本信息
	ID       string        `json:"id"`
	Code     string        `json:"code"`
	Category ErrorCategory `json:"category"`
	Level    ErrorLevel    `json:"level"`
	Message  string        `json:"message"`
	Details  string        `json:"details,omitempty"`

	// 上下文信息
	Service   string `json:"service,omitempty"`
	Operation string `json:"operation,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`

	// 技术细节
	Cause    error                  `json:"-"`
	Stack    string                 `json:"stack,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// HTTP相关
	Timestamp  time.Time `json:"timestamp"`
	HTTPStatus int       `json:"http_status"`
	Retryable  bool      `json:"retryable"`
}

// Error 实现error接口
func (e *UnifiedError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 支持errors.Unwrap
func (e *UnifiedError) Unwrap() error {
	return e.Cause
}

// GetHTTPStatus 获取HTTP状态码
func (e *UnifiedError) GetHTTPStatus() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	// 根据错误类别和代码映射默认状态码
	switch e.Category {
	case CategoryValidation:
		return http.StatusBadRequest
	case CategoryAuth:
		return http.StatusUnauthorized
	case CategoryBusiness:
		return http.StatusConflict
	case CategoryExternal:
		return http.StatusBadGateway
	case CategoryNetwork:
		return http.StatusServiceUnavailable
	case CategoryDatabase, CategoryCache:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsRetryable 判断是否可重试
func (e *UnifiedError) IsRetryable() bool {
	return e.Retryable
}

// AddMetadata 添加元数据
func (e *UnifiedError) AddMetadata(key string, value interface{}) *UnifiedError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithContext 添加上下文信息
func (e *UnifiedError) WithContext(userID, requestID, traceID string) *UnifiedError {
	e.UserID = userID
	e.RequestID = requestID
	e.TraceID = traceID
	return e
}

// WithService 设置服务信息
func (e *UnifiedError) WithService(service, operation string) *UnifiedError {
	e.Service = service
	e.Operation = operation
	return e
}

// WithOperation 设置操作信息
func (e *UnifiedError) WithOperation(operation string) *UnifiedError {
	e.Operation = operation
	return e
}

// WithMetadata 设置元数据
func (e *UnifiedError) WithMetadata(metadata map[string]interface{}) *UnifiedError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		e.Metadata[k] = v
	}
	return e
}

// ToJSON 转换为JSON格式
func (e *UnifiedError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ErrorBuilder 错误构建器
type ErrorBuilder struct {
	error *UnifiedError
}

// NewErrorBuilder 创建错误构建器
func NewErrorBuilder() *ErrorBuilder {
	return &ErrorBuilder{
		error: &UnifiedError{
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
	}
}

// WithID 设置错误ID
func (b *ErrorBuilder) WithID(id string) *ErrorBuilder {
	b.error.ID = id
	return b
}

// WithCode 设置错误代码
func (b *ErrorBuilder) WithCode(code string) *ErrorBuilder {
	b.error.Code = code
	return b
}

// WithCategory 设置错误分类
func (b *ErrorBuilder) WithCategory(category ErrorCategory) *ErrorBuilder {
	b.error.Category = category
	return b
}

// WithLevel 设置错误级别
func (b *ErrorBuilder) WithLevel(level ErrorLevel) *ErrorBuilder {
	b.error.Level = level
	return b
}

// WithMessage 设置错误消息
func (b *ErrorBuilder) WithMessage(message string) *ErrorBuilder {
	b.error.Message = message
	return b
}

// WithDetails 设置错误详情
func (b *ErrorBuilder) WithDetails(details string) *ErrorBuilder {
	b.error.Details = details
	return b
}

// WithCause 设置原因错误
func (b *ErrorBuilder) WithCause(cause error) *ErrorBuilder {
	b.error.Cause = cause
	return b
}

// WithStack 设置堆栈信息
func (b *ErrorBuilder) WithStack() *ErrorBuilder {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	b.error.Stack = string(buf[:n])
	return b
}

// WithHTTPStatus 设置HTTP状态码
func (b *ErrorBuilder) WithHTTPStatus(status int) *ErrorBuilder {
	b.error.HTTPStatus = status
	return b
}

// WithRetryable 设置是否可重试
func (b *ErrorBuilder) WithRetryable(retryable bool) *ErrorBuilder {
	b.error.Retryable = retryable
	return b
}

// WithService 设置服务信息
func (b *ErrorBuilder) WithService(service, operation string) *ErrorBuilder {
	b.error.Service = service
	b.error.Operation = operation
	return b
}

// WithContext 设置上下文信息
func (b *ErrorBuilder) WithContext(userID, requestID, traceID string) *ErrorBuilder {
	b.error.UserID = userID
	b.error.RequestID = requestID
	b.error.TraceID = traceID
	return b
}

// WithMetadata 设置元数据
func (b *ErrorBuilder) WithMetadata(key string, value interface{}) *ErrorBuilder {
	if b.error.Metadata == nil {
		b.error.Metadata = make(map[string]interface{})
	}
	b.error.Metadata[key] = value
	return b
}

// Build 构建错误
func (b *ErrorBuilder) Build() *UnifiedError {
	return b.error
}
