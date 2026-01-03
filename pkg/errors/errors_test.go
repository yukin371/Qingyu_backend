package errors

import (
	"errors"
	"net/http"
	"testing"
)

// TestNew 测试创建错误
func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		code      ErrorCode
		message   string
		wantCode  ErrorCode
		wantMsg   string
		wantHTTP  int
	}{
		{
			name:     "无效参数错误",
			code:     InvalidParams,
			message:  "测试参数无效",
			wantCode: InvalidParams,
			wantMsg:  "测试参数无效",
			wantHTTP: http.StatusBadRequest,
		},
		{
			name:     "资源不存在",
			code:     NotFound,
			message:  "用户不存在",
			wantCode: NotFound,
			wantMsg:  "用户不存在",
			wantHTTP: http.StatusNotFound,
		},
		{
			name:     "使用默认消息",
			code:     InternalError,
			message:  "",
			wantCode: InternalError,
			wantMsg:  GetDefaultMessage(InternalError),
			wantHTTP: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message)

			if err.Info.Code != tt.wantCode {
				t.Errorf("New().Info.Code = %v, want %v", err.Info.Code, tt.wantCode)
			}

			if err.Info.Message != tt.wantMsg {
				t.Errorf("New().Info.Message = %v, want %v", err.Info.Message, tt.wantMsg)
			}

			if err.Info.HTTPStatus != tt.wantHTTP {
				t.Errorf("New().Info.HTTPStatus = %v, want %v", err.Info.HTTPStatus, tt.wantHTTP)
			}
		})
	}
}

// TestWrap 测试包装错误
func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, DatabaseError, "数据库操作失败")

	if wrappedErr.Err != originalErr {
		t.Errorf("Wrap() Err = %v, want %v", wrappedErr.Err, originalErr)
	}

	if wrappedErr.Info.Code != DatabaseError {
		t.Errorf("Wrap() Info.Code = %v, want %v", wrappedErr.Info.Code, DatabaseError)
	}
}

// TestWithDetails 测试添加错误详情
func TestWithDetails(t *testing.T) {
	details := map[string]string{
		"field": "email",
		"error": "invalid format",
	}

	err := New(InvalidParams, "").WithDetails(details)

	if err.Info.Details == nil {
		t.Fatal("WithDetails() Details = nil, want non-nil")
	}

	errDetails, ok := err.Info.Details.(map[string]string)
	if !ok {
		t.Fatal("WithDetails() Details type assertion failed")
	}

	if errDetails["field"] != "email" {
		t.Errorf("WithDetails() Details[field] = %v, want email", errDetails["field"])
	}
}

// TestConvenienceFunctions 测试便捷函数
func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func() *AppError
		want ErrorCode
	}{
		{
			name: "InvalidParam",
			fn:   func() *AppError { return InvalidParam("email") },
			want: InvalidParams,
		},
		{
			name: "Unauthorized",
			fn:   func() *AppError { return Unauthorized("未登录") },
			want: Unauthorized,
		},
		{
			name: "NotFound",
			fn:   func() *AppError { return NotFound("用户") },
			want: NotFound,
		},
		{
			name: "InsufficientBalance",
			fn:   func() *AppError { return InsufficientBalance() },
			want: InsufficientBalance,
		},
		{
			name: "TokenExpired",
			fn:   func() *AppError { return TokenExpired() },
			want: TokenExpired,
		},
		{
			name: "RateLimit",
			fn:   func() *AppError { return RateLimit() },
			want: RateLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err.Info.Code != tt.want {
				t.Errorf("%s() Info.Code = %v, want %v", tt.name, err.Info.Code, tt.want)
			}
		})
	}
}

// TestErrorHandler_Handle 测试错误处理器
func TestErrorHandler_Handle(t *testing.T) {
	handler := &ErrorHandler{}

	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantCode    ErrorCode
		wantMessage string
	}{
		{
			name:        "无错误",
			err:         nil,
			wantStatus:  http.StatusOK,
			wantCode:    Success,
			wantMessage: "success",
		},
		{
			name:        "应用错误",
			err:         New(InvalidParams, "参数无效"),
			wantStatus:  http.StatusBadRequest,
			wantCode:    InvalidParams,
			wantMessage: "参数无效",
		},
		{
			name:        "普通错误",
			err:         errors.New("普通错误"),
			wantStatus:  http.StatusInternalServerError,
			wantCode:    InternalError,
			wantMessage: GetDefaultMessage(InternalError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, response := handler.Handle(tt.err)

			if status != tt.wantStatus {
				t.Errorf("Handle() status = %v, want %v", status, tt.wantStatus)
			}

			respMap, ok := response.(map[string]interface{})
			if !ok {
				t.Fatal("Handle() response type assertion failed")
			}

			code := respMap["code"].(ErrorCode)
			if code != tt.wantCode {
				t.Errorf("Handle() code = %v, want %v", code, tt.wantCode)
			}

			message := respMap["message"].(string)
			if message != tt.wantMessage {
				t.Errorf("Handle() message = %v, want %v", message, tt.wantMessage)
			}
		})
	}
}

// BenchmarkNew 性能测试 - 创建错误
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(InvalidParams, "测试参数")
	}
}

// BenchmarkWrap 性能测试 - 包装错误
func BenchmarkWrap(b *testing.B) {
	originalErr := errors.New("original error")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Wrap(originalErr, DatabaseError, "数据库错误")
	}
}
