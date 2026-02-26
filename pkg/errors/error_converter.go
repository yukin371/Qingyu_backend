package errors

import (
	"fmt"
	"net/http"
)

// ============================================================================
// 错误转换器
// 用于将旧的错误类型转换为统一的 UnifiedError
// ============================================================================

// LegacyErrorConverter 旧错误转换器
type LegacyErrorConverter struct {
	service string
}

// NewLegacyErrorConverter 创建旧错误转换器
func NewLegacyErrorConverter(service string) *LegacyErrorConverter {
	return &LegacyErrorConverter{service: service}
}

// ConvertFromUserError 转换 UserError 为 UnifiedError
func (c *LegacyErrorConverter) ConvertFromUserError(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	// 如果已经是 UnifiedError，直接返回
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	_ = UserServiceFactory // 预留工厂供后续使用

	// 尝试获取错误消息
	message := err.Error()

	// 根据错误消息判断错误类型
	// 这里简化处理，实际应该根据 UserError.Code 判断
	return NewErrorBuilder().
		WithCode("2001").
		WithCategory(CategoryBusiness).
		WithLevel(LevelError).
		WithMessage(message).
		WithService(c.service, "ConvertFromUserError").
		Build()
}

// ConvertFromReaderError 转换 ReaderError 为 UnifiedError
func (c *LegacyErrorConverter) ConvertFromReaderError(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	_ = ReaderServiceFactory // 预留工厂供后续使用

	return NewErrorBuilder().
		WithCode("3401").
		WithCategory(CategoryBusiness).
		WithLevel(LevelError).
		WithMessage(err.Error()).
		WithService(c.service, "ConvertFromReaderError").
		Build()
}

// ConvertFromWriterError 转换 WriterError 为 UnifiedError
func (c *LegacyErrorConverter) ConvertFromWriterError(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	_ = WriterServiceFactory // 预留工厂供后续使用

	return NewErrorBuilder().
		WithCode("3301").
		WithCategory(CategoryBusiness).
		WithLevel(LevelError).
		WithMessage(err.Error()).
		WithService(c.service, "ConvertFromWriterError").
		Build()
}

// ConvertFromSearchError 转换 SearchError 为 UnifiedError
func (c *LegacyErrorConverter) ConvertFromSearchError(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	return NewErrorBuilder().
		WithCode("5000").
		WithCategory(CategorySystem).
		WithLevel(LevelError).
		WithMessage("搜索服务错误").
		WithDetails(err.Error()).
		WithService(c.service, "ConvertFromSearchError").
		Build()
}

// ConvertFromAIError 转换 AIError 为 UnifiedError
func (c *LegacyErrorConverter) ConvertFromAIError(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	// AI 错误映射到统一错误码
	// 这里简化处理，实际应该根据 AIError.Type 判断
	return NewErrorBuilder().
		WithCode("3501").
		WithCategory(CategoryExternal).
		WithLevel(LevelError).
		WithMessage("AI服务错误").
		WithDetails(err.Error()).
		WithService(c.service, "ConvertFromAIError").
		WithRetryable(true).
		Build()
}

// ConvertFromRepositoryError 转换 RepositoryError 为 UnifiedError
func (c *LegacyErrorConverter) ConvertFromRepositoryError(err error) *UnifiedError {
	if err == nil {
		return nil
	}

	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	return NewErrorBuilder().
		WithCode("5001").
		WithCategory(CategoryDatabase).
		WithLevel(LevelCritical).
		WithMessage("数据库错误").
		WithDetails(err.Error()).
		WithService(c.service, "ConvertFromRepositoryError").
		Build()
}

// ============================================================================
// 全局转换函数
// ============================================================================

// ToUnifiedError 将任意错误转换为 UnifiedError
func ToUnifiedError(service string, err error) *UnifiedError {
	if err == nil {
		return nil
	}

	// 如果已经是 UnifiedError，直接返回
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr
	}

	converter := NewLegacyErrorConverter(service)

	// 尝试识别错误类型并转换
	// 这里简化处理，实际可以使用类型断言或错误包装
	return converter.ConvertFromRepositoryError(err)
}

// WrapError 包装错误为 UnifiedError
func WrapError(code, message string, cause error) *UnifiedError {
	return NewErrorBuilder().
		WithCode(code).
		WithCategory(CategorySystem).
		WithLevel(LevelError).
		WithMessage(message).
		WithCause(cause).
		Build()
}

// WrapValidationError 包装验证错误
func WrapValidationError(field, message string) *UnifiedError {
	return NewErrorBuilder().
		WithCode("1001").
		WithCategory(CategoryValidation).
		WithLevel(LevelWarning).
		WithMessage(message).
		WithDetails(fmt.Sprintf("field: %s", field)).
		Build()
}

// WrapNotFoundError 包装资源不存在错误
func WrapNotFoundError(resource, id string) *UnifiedError {
	message := fmt.Sprintf("%s不存在", resource)
	if id != "" {
		message = fmt.Sprintf("%s '%s' 不存在", resource, id)
	}

	return NewErrorBuilder().
		WithCode("1004").
		WithCategory(CategoryBusiness).
		WithLevel(LevelInfo).
		WithMessage(message).
		WithHTTPStatus(http.StatusNotFound).
		Build()
}

// WrapBusinessError 包装业务逻辑错误
func WrapBusinessError(code, message string) *UnifiedError {
	return NewErrorBuilder().
		WithCode(code).
		WithCategory(CategoryBusiness).
		WithLevel(LevelError).
		WithMessage(message).
		Build()
}

// WrapInternalError 包装内部错误
func WrapInternalError(message string, cause error) *UnifiedError {
	builder := NewErrorBuilder().
		WithCode("5000").
		WithCategory(CategorySystem).
		WithLevel(LevelCritical).
		WithMessage(message)

	if cause != nil {
		builder = builder.WithCause(cause).WithStack()
	}

	return builder.Build()
}

// WrapAuthError 包装认证错误
func WrapAuthError(code, message string) *UnifiedError {
	return NewErrorBuilder().
		WithCode(code).
		WithCategory(CategoryAuth).
		WithLevel(LevelWarning).
		WithMessage(message).
		WithHTTPStatus(http.StatusUnauthorized).
		Build()
}

// WrapForbiddenError 包装授权错误
func WrapForbiddenError(message string) *UnifiedError {
	return NewErrorBuilder().
		WithCode("1003").
		WithCategory(CategoryAuth).
		WithLevel(LevelWarning).
		WithMessage(message).
		WithHTTPStatus(http.StatusForbidden).
		Build()
}

// WrapConflictError 包装冲突错误
func WrapConflictError(message string) *UnifiedError {
	return NewErrorBuilder().
		WithCode("1006").
		WithCategory(CategoryBusiness).
		WithLevel(LevelWarning).
		WithMessage(message).
		WithHTTPStatus(http.StatusConflict).
		Build()
}

// WrapRateLimitError 包装频率限制错误
func WrapRateLimitError(limitType string) *UnifiedError {
	message := "请求过于频繁"
	if limitType != "" {
		message = fmt.Sprintf("%s频率限制超出", limitType)
	}

	return NewErrorBuilder().
		WithCode("4000").
		WithCategory(CategoryBusiness).
		WithLevel(LevelWarning).
		WithMessage(message).
		WithHTTPStatus(http.StatusTooManyRequests).
		WithRetryable(true).
		Build()
}

// ============================================================================
// HTTP 错误响应构建器
// ============================================================================

// HTTPErrorResponse HTTP错误响应
type HTTPErrorResponse struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Timestamp  string      `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
	TraceID    string      `json:"trace_id,omitempty"`
}

// ToHTTPResponse 将 UnifiedError 转换为 HTTP 响应
func ToHTTPResponse(err *UnifiedError, requestID, traceID string) (int, HTTPErrorResponse) {
	if err == nil {
		return http.StatusOK, HTTPErrorResponse{
			Code:      "0",
			Message:   "成功",
			Timestamp: err.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	statusCode := err.GetHTTPStatus()

	return statusCode, HTTPErrorResponse{
		Code:      err.Code,
		Message:   err.Message,
		Details:   err.Details,
		Timestamp: err.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		RequestID: requestID,
		TraceID:   traceID,
	}
}

// ToHTTPResponseWithError 将任意错误转换为 HTTP 响应
func ToHTTPResponseWithError(service string, err error, requestID, traceID string) (int, HTTPErrorResponse) {
	if err == nil {
		return http.StatusOK, HTTPErrorResponse{
			Code:      "0",
			Message:   "成功",
			Timestamp: "",
		}
	}

	unifiedErr := ToUnifiedError(service, err)
	return ToHTTPResponse(unifiedErr, requestID, traceID)
}
