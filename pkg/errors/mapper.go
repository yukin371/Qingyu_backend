package errors

import (
	"net/http"
)

// HTTPStatuser 错误状态码接口
// 任何实现了此接口的错误类型都可以自动映射到HTTP状态码
type HTTPStatuser interface {
	HTTPStatus() int
}

// MessageProvider 错误消息接口
// 任何实现了此接口的错误类型都可以提供友好的错误消息
type MessageProvider interface {
	ErrorMessage() string
}

// MapToHTTPStatus 将错误映射到HTTP状态码
//
// 支持的错误类型优先级：
// 1. UnifiedError - 使用内置的GetHTTPStatus()方法
// 2. 实现HTTPStatuser接口的错误 - 使用其HTTPStatus()方法
// 3. 其他错误 - 默认500
func MapToHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// 1. 检查UnifiedError
	if unifiedErr, ok := err.(*UnifiedError); ok {
		return unifiedErr.GetHTTPStatus()
	}

	// 2. 检查HTTPStatuser接口（ReaderError, WriterError, UserError等都实现了此接口）
	if statusErr, ok := err.(HTTPStatuser); ok {
		return statusErr.HTTPStatus()
	}

	// 默认返回500
	return http.StatusInternalServerError
}

// GetErrorMessage 从错误中提取用户友好的错误消息
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	// 1. UnifiedError
	if unifiedErr, ok := err.(*UnifiedError); ok {
		if unifiedErr.Message != "" {
			return unifiedErr.Message
		}
		return "内部错误"
	}

	// 2. MessageProvider接口
	if msgErr, ok := err.(MessageProvider); ok {
		return msgErr.ErrorMessage()
	}

	// 3. 标准错误 - 返回错误消息
	return err.Error()
}

// IsNotFoundError 判断是否为404类型错误
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return MapToHTTPStatus(err) == http.StatusNotFound
}

// IsUnauthorizedError 判断是否为401类型错误
func IsUnauthorizedError(err error) bool {
	if err == nil {
		return false
	}
	return MapToHTTPStatus(err) == http.StatusUnauthorized
}

// IsForbiddenError 判断是否为403类型错误
func IsForbiddenError(err error) bool {
	if err == nil {
		return false
	}
	return MapToHTTPStatus(err) == http.StatusForbidden
}

// IsConflictError 判断是否为409类型错误
func IsConflictError(err error) bool {
	if err == nil {
		return false
	}
	return MapToHTTPStatus(err) == http.StatusConflict
}

// IsClientError 判断是否为4xx客户端错误
func IsClientError(err error) bool {
	if err == nil {
		return false
	}
	status := MapToHTTPStatus(err)
	return status >= 400 && status < 500
}

// IsServerError 判断是否为5xx服务端错误
func IsServerError(err error) bool {
	if err == nil {
		return false
	}
	status := MapToHTTPStatus(err)
	return status >= 500 && status < 600
}
