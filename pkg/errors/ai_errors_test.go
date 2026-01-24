package errors

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAIError_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      *AIError
		expected bool
	}{
		{"AI服务不可用", &AIError{Type: ErrAIUnavailable}, true},
		{"超时", &AIError{Type: ErrAITimeout}, true},
		{"频率限制", &AIError{Type: ErrAIRateLimit}, true},
		{"模型不可用", &AIError{Type: ErrModelNotAvailable}, true},
		{"配额不足", &AIError{Type: ErrQuotaExhausted}, false},
		{"参数无效", &AIError{Type: ErrInvalidRequest}, false},
		{"请求过大", &AIError{Type: ErrRequestTooLarge}, false},
		{"上下文超限", &AIError{Type: ErrContextExceeded}, false},
		{"配额检查失败", &AIError{Type: ErrQuotaCheckFail}, false},
		{"模型不存在", &AIError{Type: ErrModelNotFound}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsRetryable(); got != tt.expected {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAIError_Error(t *testing.T) {
	err := &AIError{
		Type:    ErrAIUnavailable,
		Message: "service down",
	}

	expected := "[1001] service down"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %v, want %v", got, expected)
	}
}

func TestAIError_Unwrap(t *testing.T) {
	internalErr := errors.New("internal error")
	aiErr := &AIError{
		Type:     ErrAIUnavailable,
		Message:  "service down",
		Internal: internalErr,
	}

	if got := errors.Unwrap(aiErr); got != internalErr {
		t.Errorf("Unwrap() = %v, want %v", got, internalErr)
	}
}

func TestAIError_GRPCCode(t *testing.T) {
	tests := []struct {
		name     string
		errType  AIErrorType
		expected codes.Code
	}{
		{"AI服务不可用", ErrAIUnavailable, codes.Unavailable},
		{"超时", ErrAITimeout, codes.DeadlineExceeded},
		{"频率限制", ErrAIRateLimit, codes.ResourceExhausted},
		{"配额不足", ErrQuotaExhausted, codes.ResourceExhausted},
		{"参数无效", ErrInvalidRequest, codes.InvalidArgument},
		{"请求过大", ErrRequestTooLarge, codes.InvalidArgument},
		{"上下文超限", ErrContextExceeded, codes.InvalidArgument},
		{"模型不存在", ErrModelNotFound, codes.NotFound},
		{"模型不可用", ErrModelNotAvailable, codes.NotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &AIError{Type: tt.errType}
			if got := err.GRPCCode(); got != tt.expected {
				t.Errorf("GRPCCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAIError_ToErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		errType  AIErrorType
		expected ErrorCode
	}{
		{"配额不足", ErrQuotaExhausted, InsufficientQuota},
		{"配额检查失败", ErrQuotaCheckFail, InsufficientQuota},
		{"参数无效", ErrInvalidRequest, InvalidParams},
		{"AI服务不可用", ErrAIUnavailable, ExternalAPIError},
		{"超时", ErrAITimeout, ExternalAPIError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &AIError{Type: tt.errType}
			if got := err.ToErrorCode(); got != tt.expected {
				t.Errorf("ToErrorCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAIError_ToStatus(t *testing.T) {
	err := &AIError{
		Type:    ErrAIUnavailable,
		Message: "service down",
	}

	st := err.ToStatus()

	if st.Code() != codes.Unavailable {
		t.Errorf("ToStatus().Code() = %v, want %v", st.Code(), codes.Unavailable)
	}

	if st.Message() != "service down" {
		t.Errorf("ToStatus().Message() = %v, want %v", st.Message(), "service down")
	}
}

func TestNewAIError(t *testing.T) {
	internalErr := errors.New("internal")
	err := NewAIError(ErrAIUnavailable, "test message", internalErr)

	if err.Type != ErrAIUnavailable {
		t.Errorf("Type = %v, want %v", err.Type, ErrAIUnavailable)
	}

	if err.Message != "test message" {
		t.Errorf("Message = %v, want %v", err.Message, "test message")
	}

	if err.Internal != internalErr {
		t.Errorf("Internal = %v, want %v", err.Internal, internalErr)
	}
}

func TestNewAIErrorWithDetails(t *testing.T) {
	details := map[string]interface{}{
		"retry_after": 60,
		"endpoint":    "localhost:50051",
	}

	err := NewAIErrorWithDetails(ErrAIRateLimit, "rate limited", nil, details)

	if err.Details == nil {
		t.Fatal("Details should not be nil")
	}

	if len(err.Details) != 2 {
		t.Errorf("Details length = %v, want %v", len(err.Details), 2)
	}
}

func TestIsAIError(t *testing.T) {
	aiErr := &AIError{Type: ErrAIUnavailable}
	stdErr := errors.New("standard error")

	if !IsAIError(aiErr) {
		t.Error("IsAIError(aiErr) should return true")
	}

	if IsAIError(stdErr) {
		t.Error("IsAIError(stdErr) should return false")
	}
}

func TestConvertGRPCErrorToAIError(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		check    func(*testing.T, *AIError)
	}{
		{
			name:     "nil error",
			inputErr: nil,
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr != nil {
					t.Error("Expected nil for nil input")
				}
			},
		},
		{
			name:     "already AIError",
			inputErr: &AIError{Type: ErrAIUnavailable},
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr.Type != ErrAIUnavailable {
					t.Errorf("Type = %v, want %v", aiErr.Type, ErrAIUnavailable)
				}
			},
		},
		{
			name:     "gRPC unavailable error",
			inputErr: status.Error(codes.Unavailable, "service unavailable"),
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr.Type != ErrAIUnavailable {
					t.Errorf("Type = %v, want %v", aiErr.Type, ErrAIUnavailable)
				}
			},
		},
		{
			name:     "gRPC deadline exceeded",
			inputErr: status.Error(codes.DeadlineExceeded, "timeout"),
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr.Type != ErrAITimeout {
					t.Errorf("Type = %v, want %v", aiErr.Type, ErrAITimeout)
				}
			},
		},
		{
			name:     "gRPC resource exhausted",
			inputErr: status.Error(codes.ResourceExhausted, "quota exceeded"),
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr.Type != ErrQuotaExhausted {
					t.Errorf("Type = %v, want %v", aiErr.Type, ErrQuotaExhausted)
				}
			},
		},
		{
			name:     "gRPC invalid argument",
			inputErr: status.Error(codes.InvalidArgument, "bad request"),
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr.Type != ErrInvalidRequest {
					t.Errorf("Type = %v, want %v", aiErr.Type, ErrInvalidRequest)
				}
			},
		},
		{
			name:     "gRPC not found",
			inputErr: status.Error(codes.NotFound, "model not found"),
			check: func(t *testing.T, aiErr *AIError) {
				if aiErr.Type != ErrModelNotFound {
					t.Errorf("Type = %v, want %v", aiErr.Type, ErrModelNotFound)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertGRPCErrorToAIError(tt.inputErr)
			tt.check(t, result)
		})
	}
}
