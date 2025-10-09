package errors

import (
	"fmt"
	"strings"
)

// ErrorChain 错误链，用于管理一系列相关错误
type ErrorChain struct {
	errors []*UnifiedError
}

// NewErrorChain 创建错误链
func NewErrorChain() *ErrorChain {
	return &ErrorChain{
		errors: make([]*UnifiedError, 0),
	}
}

// Add 添加错误到链中
func (ec *ErrorChain) Add(err *UnifiedError) *ErrorChain {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
	return ec
}

// HasErrors 检查是否有错误
func (ec *ErrorChain) HasErrors() bool {
	return len(ec.errors) > 0
}

// GetErrors 获取所有错误
func (ec *ErrorChain) GetErrors() []*UnifiedError {
	return ec.errors
}

// GetFirst 获取第一个错误
func (ec *ErrorChain) GetFirst() *UnifiedError {
	if len(ec.errors) > 0 {
		return ec.errors[0]
	}
	return nil
}

// GetLast 获取最后一个错误
func (ec *ErrorChain) GetLast() *UnifiedError {
	if len(ec.errors) > 0 {
		return ec.errors[len(ec.errors)-1]
	}
	return nil
}

// Count 获取错误数量
func (ec *ErrorChain) Count() int {
	return len(ec.errors)
}

// Clear 清空错误链
func (ec *ErrorChain) Clear() {
	ec.errors = ec.errors[:0]
}

// ErrorAggregator 错误聚合器，用于收集和分组错误
type ErrorAggregator struct {
	errors map[string][]*UnifiedError
}

// NewErrorAggregator 创建错误聚合器
func NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{
		errors: make(map[string][]*UnifiedError),
	}
}

// Add 添加错误到指定键
func (ea *ErrorAggregator) Add(key string, err *UnifiedError) {
	if err != nil {
		ea.errors[key] = append(ea.errors[key], err)
	}
}

// Get 获取指定键的错误列表
func (ea *ErrorAggregator) Get(key string) []*UnifiedError {
	return ea.errors[key]
}

// GetAll 获取所有错误
func (ea *ErrorAggregator) GetAll() map[string][]*UnifiedError {
	return ea.errors
}

// HasErrors 检查是否有错误
func (ea *ErrorAggregator) HasErrors() bool {
	return len(ea.errors) > 0
}

// HasErrorsForKey 检查指定键是否有错误
func (ea *ErrorAggregator) HasErrorsForKey(key string) bool {
	return len(ea.errors[key]) > 0
}

// Clear 清空所有错误
func (ea *ErrorAggregator) Clear() {
	ea.errors = make(map[string][]*UnifiedError)
}

// ClearKey 清空指定键的错误
func (ea *ErrorAggregator) ClearKey(key string) {
	delete(ea.errors, key)
}

// ErrorContext 错误上下文包装器
type ErrorContext struct {
	userID    string
	requestID string
	traceID   string
	service   string
	operation string
}

// NewErrorContext 创建错误上下文
func NewErrorContext(userID, requestID, traceID, service, operation string) *ErrorContext {
	return &ErrorContext{
		userID:    userID,
		requestID: requestID,
		traceID:   traceID,
		service:   service,
		operation: operation,
	}
}

// Wrap 包装错误并添加上下文
func (ec *ErrorContext) Wrap(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	// 如果已经是UnifiedError，添加上下文信息
	if unifiedErr, ok := err.(*UnifiedError); ok {
		unifiedErr.UserID = ec.userID
		unifiedErr.RequestID = ec.requestID
		unifiedErr.TraceID = ec.traceID
		if unifiedErr.Service == "" {
			unifiedErr.Service = ec.service
		}
		if unifiedErr.Operation == "" {
			unifiedErr.Operation = ec.operation
		}
		return unifiedErr
	}

	// 转换为UnifiedError并添加上下文
	return DefaultConverter.ConvertGenericError(err, ec.service, ec.operation).
		WithContext(ec.userID, ec.requestID, ec.traceID)
}

// 工具函数

// IsErrorType 检查错误是否为指定类型
func IsErrorType(err error, category ErrorCategory) bool {
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr.Category == category
	}
	return false
}

// IsRetryableError 检查错误是否可重试
func IsRetryableError(err error) bool {
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr.IsRetryable()
	}
	return false
}

// GetErrorCode 获取错误代码
func GetErrorCode(err error) string {
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr.Code
	}
	return "UNKNOWN"
}

// GetErrorCategory 获取错误分类
func GetErrorCategory(err error) ErrorCategory {
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr.Category
	}
	return CategorySystem
}

// GetHTTPStatus 获取HTTP状态码
func GetHTTPStatus(err error) int {
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr.GetHTTPStatus()
	}
	return 500
}

// FormatErrorChain 格式化错误链为字符串
func FormatErrorChain(chain *ErrorChain) string {
	if !chain.HasErrors() {
		return ""
	}

	var parts []string
	for i, err := range chain.GetErrors() {
		parts = append(parts, fmt.Sprintf("%d. [%s] %s", i+1, err.Code, err.Message))
	}

	return strings.Join(parts, "\n")
}

// FormatErrorAggregator 格式化错误聚合器为字符串
func FormatErrorAggregator(aggregator *ErrorAggregator) string {
	if !aggregator.HasErrors() {
		return ""
	}

	var parts []string
	for key, errors := range aggregator.GetAll() {
		parts = append(parts, fmt.Sprintf("%s:", key))
		for _, err := range errors {
			parts = append(parts, fmt.Sprintf("  - [%s] %s", err.Code, err.Message))
		}
	}

	return strings.Join(parts, "\n")
}

// CombineErrors 合并多个错误为一个
func CombineErrors(errors ...*UnifiedError) *UnifiedError {
	if len(errors) == 0 {
		return nil
	}

	// 过滤nil错误
	validErrors := make([]*UnifiedError, 0, len(errors))
	for _, err := range errors {
		if err != nil {
			validErrors = append(validErrors, err)
		}
	}

	if len(validErrors) == 0 {
		return nil
	}

	if len(validErrors) == 1 {
		return validErrors[0]
	}

	// 创建组合错误
	first := validErrors[0]
	var messages []string
	var details []string

	for _, err := range validErrors {
		messages = append(messages, err.Message)
		if err.Details != "" {
			details = append(details, err.Details)
		}
	}

	return &UnifiedError{
		ID:         generateErrorID(),
		Code:       "MULTIPLE_ERRORS",
		Category:   first.Category,
		Level:      first.Level,
		Message:    strings.Join(messages, "; "),
		Details:    strings.Join(details, "; "),
		Service:    first.Service,
		Operation:  first.Operation,
		UserID:     first.UserID,
		RequestID:  first.RequestID,
		TraceID:    first.TraceID,
		HTTPStatus: first.HTTPStatus,
		Retryable:  false,
		Metadata: map[string]interface{}{
			"error_count": len(validErrors),
			"error_codes": func() []string {
				codes := make([]string, len(validErrors))
				for i, err := range validErrors {
					codes[i] = err.Code
				}
				return codes
			}(),
		},
	}
}
