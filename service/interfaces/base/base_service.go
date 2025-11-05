package base

import (
	"context"
	"fmt"
	"time"
)

// BaseService 基础Service接口
type BaseService interface {
	// 服务生命周期
	Initialize(ctx context.Context) error
	Health(ctx context.Context) error
	Close(ctx context.Context) error

	// 服务信息
	GetServiceName() string
	GetVersion() string
}

// ServiceError 服务错误类型
type ServiceError struct {
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Cause     error     `json:"cause,omitempty"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Service, e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Service, e.Type, e.Message)
}

// 错误类型常量
const (
	ErrorTypeValidation   = "VALIDATION"
	ErrorTypeBusiness     = "BUSINESS"
	ErrorTypeNotFound     = "NOT_FOUND"
	ErrorTypeUnauthorized = "UNAUTHORIZED"
	ErrorTypeForbidden    = "FORBIDDEN"
	ErrorTypeInternal     = "INTERNAL"
	ErrorTypeTimeout      = "TIMEOUT"
	ErrorTypeExternal     = "EXTERNAL"
)

// NewServiceError 创建服务错误
func NewServiceError(service, errorType, message string, cause error) *ServiceError {
	return &ServiceError{
		Type:      errorType,
		Message:   message,
		Cause:     cause,
		Service:   service,
		Timestamp: time.Now(),
	}
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	if serviceErr, ok := err.(*ServiceError); ok {
		return serviceErr.Type == ErrorTypeValidation
	}
	return false
}

// IsBusinessError 检查是否为业务错误
func IsBusinessError(err error) bool {
	if serviceErr, ok := err.(*ServiceError); ok {
		return serviceErr.Type == ErrorTypeBusiness
	}
	return false
}

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	if serviceErr, ok := err.(*ServiceError); ok {
		return serviceErr.Type == ErrorTypeNotFound
	}
	return false
}

// Event 事件接口
type Event interface {
	GetEventType() string
	GetEventData() interface{}
	GetTimestamp() time.Time
	GetSource() string
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	GetHandlerName() string
	GetSupportedEventTypes() []string
}

// EventBus 事件总线接口
type EventBus interface {
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handlerName string) error
	Publish(ctx context.Context, event Event) error
	PublishAsync(ctx context.Context, event Event) error
}
