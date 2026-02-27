package ai

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCError gRPC调用错误
type GRPCError struct {
	Code     string
	Message  string
	Original error
}

// Error 实现error接口
func (e *GRPCError) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Original)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现错误包装接口，支持errors.Unwrap
func (e *GRPCError) Unwrap() error {
	return e.Original
}

// 错误码常量
const (
	// 服务不可用错误
	ErrCodeUnavailable   = "AI_SERVICE_UNAVAILABLE"
	ErrCodeTimeout       = "AI_SERVICE_TIMEOUT"
	ErrCodeNetworkError  = "AI_SERVICE_NETWORK_ERROR"

	// 请求错误
	ErrCodeInvalidReq    = "INVALID_REQUEST"
	ErrCodeMissingParam  = "MISSING_PARAMETER"
	ErrCodeInvalidParam  = "INVALID_PARAMETER"

	// 业务逻辑错误
	ErrCodeExecutionFailed = "EXECUTION_FAILED"
	ErrCodeWorkflowFailed  = "WORKFLOW_FAILED"
	ErrCodeGenerationFailed = "GENERATION_FAILED"

	// 资源错误
	ErrCodeResourceExhausted = "RESOURCE_EXHAUSTED"
	ErrCodeQuotaExceeded     = "QUOTA_EXCEEDED"

	// 权限错误
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"

	// 内部错误
	ErrCodeInternal   = "INTERNAL_ERROR"
	ErrCodeUnknown    = "UNKNOWN_ERROR"
)

// gRPC状态码到错误码的映射
var grpcCodeMapping = map[codes.Code]string{
	codes.Unavailable:        ErrCodeUnavailable,
	codes.DeadlineExceeded:    ErrCodeTimeout,
	codes.InvalidArgument:     ErrCodeInvalidReq,
	codes.NotFound:            ErrCodeInvalidParam,
	codes.AlreadyExists:       ErrCodeInvalidParam,
	codes.PermissionDenied:    ErrCodeForbidden,
	codes.Unauthenticated:     ErrCodeUnauthorized,
	codes.ResourceExhausted:   ErrCodeResourceExhausted,
	codes.FailedPrecondition:  ErrCodeExecutionFailed,
	codes.OutOfRange:          ErrCodeInvalidParam,
	codes.Unimplemented:       ErrCodeInternal,
	codes.Internal:            ErrCodeInternal,
	codes.DataLoss:            ErrCodeInternal,
	codes.Unknown:             ErrCodeUnknown,
}

// 错误码到用户友好消息的映射
var errorMessages = map[string]string{
	ErrCodeUnavailable:        "AI服务暂时不可用，请稍后重试",
	ErrCodeTimeout:            "请求超时，请稍后重试",
	ErrCodeNetworkError:       "网络连接失败，请检查网络设置",
	ErrCodeInvalidReq:         "请求参数无效",
	ErrCodeMissingParam:       "缺少必需参数",
	ErrCodeInvalidParam:       "参数值无效",
	ErrCodeExecutionFailed:    "AI执行失败",
	ErrCodeWorkflowFailed:     "工作流执行失败",
	ErrCodeGenerationFailed:   "内容生成失败",
	ErrCodeResourceExhausted:  "服务资源不足，请稍后重试",
	ErrCodeQuotaExceeded:      "配额已用完，请升级套餐",
	ErrCodeUnauthorized:       "未授权访问",
	ErrCodeForbidden:          "无权访问此资源",
	ErrCodeInternal:           "服务内部错误",
	ErrCodeUnknown:            "未知错误",
}

// WrapGRPCError 包装gRPC错误为统一的GRPCError
// 将原始gRPC错误转换为带有标准错误码的GRPCError
func WrapGRPCError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	// 如果已经是GRPCError，直接返回
	if grpcErr, ok := err.(*GRPCError); ok {
		return grpcErr
	}

	// 从gRPC状态中提取错误码和消息
	st, ok := status.FromError(err)
	if !ok {
		// 不是gRPC状态错误，返回通用错误
		return &GRPCError{
			Code:     ErrCodeUnknown,
			Message:  "未知错误",
			Original: err,
		}
	}

	// 映射gRPC状态码到自定义错误码
	errorCode, exists := grpcCodeMapping[st.Code()]
	if !exists {
		errorCode = ErrCodeUnknown
	}

	// 获取用户友好的错误消息
	errorMessage := getErrorMessage(errorCode, st.Message())

	return &GRPCError{
		Code:     errorCode,
		Message:  errorMessage,
		Original: err,
	}
}

// NewGRPCError 创建新的GRPCError
func NewGRPCError(code, message string, original error) *GRPCError {
	return &GRPCError{
		Code:     code,
		Message:  message,
		Original: original,
	}
}

// IsRetryableError 判断错误是否可重试
// 根据错误码判断是否应该重试请求
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	grpcErr := WrapGRPCError(err)
	if grpcErr == nil {
		return false
	}

	switch grpcErr.Code {
	case ErrCodeUnavailable,
	     ErrCodeTimeout,
	     ErrCodeNetworkError,
	     ErrCodeResourceExhausted:
		return true
	default:
		return false
	}
}

// IsTemporaryError 判断错误是否是临时性的
// 临时性错误通常可以通过重试解决
func IsTemporaryError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable,
	     codes.DeadlineExceeded,
	     codes.ResourceExhausted,
	     codes.Aborted:
		return true
	default:
		return false
	}
}

// GetErrorCode 从错误中提取错误码
func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}

	grpcErr := WrapGRPCError(err)
	if grpcErr != nil {
		return grpcErr.Code
	}

	return ErrCodeUnknown
}

// GetErrorMessage 从错误中提取用户友好的错误消息
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	grpcErr := WrapGRPCError(err)
	if grpcErr != nil {
		return grpcErr.Message
	}

	return "未知错误"
}

// getErrorMessage 获取错误消息
// 优先使用预定义的消息，如果没有则使用原始消息
func getErrorMessage(code string, originalMessage string) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}

	// 如果没有预定义消息，使用原始消息但进行清理
	message := strings.TrimSpace(originalMessage)
	if message == "" {
		message = "操作失败，请稍后重试"
	}

	return message
}

// ParseWorkflowStatus 解析工作流状态
// 将工作流状态字符串转换为标准错误码
func ParseWorkflowStatus(status string) string {
	switch status {
	case "completed":
		return "" // 成功状态返回空
	case "failed", "error":
		return ErrCodeExecutionFailed
	case "timeout":
		return ErrCodeTimeout
	case "cancelled":
		return ErrCodeInvalidReq
	default:
		return ErrCodeWorkflowFailed
	}
}

// ErrorWithCode 创建带有错误码的错误
func ErrorWithCode(code, message string) error {
	return &GRPCError{
		Code:    code,
		Message: message,
	}
}

// ErrorWithCodeAndOriginal 创建带有错误码和原始错误的错误
func ErrorWithCodeAndOriginal(code, message string, original error) error {
	return &GRPCError{
		Code:     code,
		Message:  message,
		Original: original,
	}
}

// IsUnavailable 判断错误是否为服务不可用
func IsUnavailable(err error) bool {
	return GetErrorCode(err) == ErrCodeUnavailable
}

// IsTimeout 判断错误是否为超时
func IsTimeout(err error) bool {
	return GetErrorCode(err) == ErrCodeTimeout
}

// IsInvalidRequest 判断错误是否为无效请求
func IsInvalidRequest(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeInvalidReq || code == ErrCodeMissingParam || code == ErrCodeInvalidParam
}
