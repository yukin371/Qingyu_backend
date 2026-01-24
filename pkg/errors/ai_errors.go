package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AIErrorType AI服务错误类型
type AIErrorType int

const (
	// AI服务不可用错误 1001-1099
	ErrAIUnavailable AIErrorType = 1001 // AI 服务不可用
	ErrAITimeout     AIErrorType = 1002 // AI 服务超时
	ErrAIRateLimit   AIErrorType = 1003 // AI 服务频率限制

	// 配额相关错误 2001-2099
	ErrQuotaExhausted AIErrorType = 2001 // 配额不足
	ErrQuotaCheckFail AIErrorType = 2002 // 配额检查失败

	// 请求相关错误 3001-3099
	ErrInvalidRequest  AIErrorType = 3001 // 请求参数无效
	ErrRequestTooLarge AIErrorType = 3002 // 请求过大
	ErrContextExceeded AIErrorType = 3003 // 上下文超限

	// 模型相关错误 4001-4099
	ErrModelNotFound   AIErrorType = 4001 // 模型不存在
	ErrModelNotAvailable AIErrorType = 4002 // 模型不可用
)

// AIError AI服务错误
type AIError struct {
	Type     AIErrorType
	Message  string
	Details  map[string]interface{}
	Internal error // 内部错误（用于日志）
}

// Error 实现 error 接口
func (e *AIError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Type, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AIError) Unwrap() error {
	return e.Internal
}

// IsRetryable 判断错误是否可重试
func (e *AIError) IsRetryable() bool {
	switch e.Type {
	case ErrAIUnavailable, ErrAITimeout, ErrAIRateLimit, ErrModelNotAvailable:
		return true
	default:
		return false
	}
}

// GRPCCode 转换为 gRPC 状态码
func (e *AIError) GRPCCode() codes.Code {
	switch e.Type {
	case ErrAIUnavailable:
		return codes.Unavailable
	case ErrAITimeout:
		return codes.DeadlineExceeded
	case ErrAIRateLimit:
		return codes.ResourceExhausted
	case ErrQuotaExhausted, ErrQuotaCheckFail:
		return codes.ResourceExhausted
	case ErrInvalidRequest:
		return codes.InvalidArgument
	case ErrRequestTooLarge, ErrContextExceeded:
		return codes.InvalidArgument
	case ErrModelNotFound, ErrModelNotAvailable:
		return codes.NotFound
	default:
		return codes.Internal
	}
}

// ToStatus 转换为 gRPC status
func (e *AIError) ToStatus() *status.Status {
	return status.New(e.GRPCCode(), e.Message)
}

// ToErrorCode 转换为项目通用 ErrorCode
func (e *AIError) ToErrorCode() ErrorCode {
	switch e.Type {
	case ErrQuotaExhausted, ErrQuotaCheckFail:
		return InsufficientQuota
	case ErrInvalidRequest:
		return InvalidParams
	case ErrAIUnavailable, ErrAITimeout:
		return ExternalAPIError
	default:
		return InternalError
	}
}

// NewAIError 创建 AI 错误
func NewAIError(errorType AIErrorType, message string, internal error) *AIError {
	return &AIError{
		Type:     errorType,
		Message:  message,
		Internal: internal,
	}
}

// NewAIErrorWithDetails 创建带详情的 AI 错误
func NewAIErrorWithDetails(errorType AIErrorType, message string, internal error, details map[string]interface{}) *AIError {
	return &AIError{
		Type:     errorType,
		Message:  message,
		Details:  details,
		Internal: internal,
	}
}

// IsAIError 判断是否为 AI 错误
func IsAIError(err error) bool {
	_, ok := err.(*AIError)
	return ok
}

// ConvertGRPCErrorToAIError 将 gRPC 错误转换为 AI 错误
func ConvertGRPCErrorToAIError(err error) *AIError {
	if err == nil {
		return nil
	}

	// 如果已经是 AIError，直接返回
	if aiErr, ok := err.(*AIError); ok {
		return aiErr
	}

	// 解析 gRPC status
	st, ok := status.FromError(err)
	if !ok {
		return NewAIError(ErrAIUnavailable, "AI service error", err)
	}

	// 根据 gRPC 状态码映射错误类型
	switch st.Code() {
	case codes.Unavailable:
		return NewAIError(ErrAIUnavailable, st.Message(), err)
	case codes.DeadlineExceeded:
		return NewAIError(ErrAITimeout, st.Message(), err)
	case codes.ResourceExhausted:
		return NewAIError(ErrQuotaExhausted, st.Message(), err)
	case codes.InvalidArgument:
		return NewAIError(ErrInvalidRequest, st.Message(), err)
	case codes.NotFound:
		return NewAIError(ErrModelNotFound, st.Message(), err)
	default:
		return NewAIError(ErrAIUnavailable, st.Message(), err)
	}
}
