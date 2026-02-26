package errors_test

import (
	stderrors "errors"
	"net/http"
	"testing"

	"Qingyu_backend/pkg/errors"
)

// ============================================================================
// 错误创建测试
// ============================================================================

func TestErrorBuilder_Build(t *testing.T) {
	tests := []struct {
		name         string
		builder      *errors.ErrorBuilder
		wantCode     string
		wantCategory errors.ErrorCategory
		wantLevel    errors.ErrorLevel
		wantMessage  string
		wantDetails  string
	}{
		{
			name: "创建验证错误",
			builder: errors.NewErrorBuilder().
				WithCode("1001").
				WithCategory(errors.CategoryValidation).
				WithLevel(errors.LevelWarning).
				WithMessage("参数无效").
				WithDetails("用户名长度必须在3-20个字符之间").
				WithHTTPStatus(http.StatusBadRequest),
			wantCode:     "1001",
			wantCategory: errors.CategoryValidation,
			wantLevel:    errors.LevelWarning,
			wantMessage:  "参数无效",
			wantDetails:  "用户名长度必须在3-20个字符之间",
		},
		{
			name: "创建业务错误",
			builder: errors.NewErrorBuilder().
				WithCode("2003").
				WithCategory(errors.CategoryBusiness).
				WithLevel(errors.LevelError).
				WithMessage("用户名已被使用"),
			wantCode:     "2003",
			wantCategory: errors.CategoryBusiness,
			wantLevel:    errors.LevelError,
			wantMessage:  "用户名已被使用",
		},
		{
			name: "创建系统错误",
			builder: errors.NewErrorBuilder().
				WithCode("5000").
				WithCategory(errors.CategorySystem).
				WithLevel(errors.LevelCritical).
				WithMessage("内部错误").
				WithHTTPStatus(http.StatusInternalServerError),
			wantCode:     "5000",
			wantCategory: errors.CategorySystem,
			wantLevel:    errors.LevelCritical,
			wantMessage:  "内部错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.builder.Build()

			if err == nil {
				t.Fatal("Build() 返回 nil")
			}

			if err.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", err.Code, tt.wantCode)
			}

			if err.Category != tt.wantCategory {
				t.Errorf("Category = %s, want %s", err.Category, tt.wantCategory)
			}

			if err.Level != tt.wantLevel {
				t.Errorf("Level = %s, want %s", err.Level, tt.wantLevel)
			}

			if err.Message != tt.wantMessage {
				t.Errorf("Message = %s, want %s", err.Message, tt.wantMessage)
			}

			if tt.wantDetails != "" && err.Details != tt.wantDetails {
				t.Errorf("Details = %s, want %s", err.Details, tt.wantDetails)
			}
		})
	}
}

// ============================================================================
// 错误工厂测试
// ============================================================================

func TestErrorFactory_ValidationError(t *testing.T) {
	factory := errors.NewErrorFactory("test-service")

	err := factory.ValidationError("1001", "参数无效", "用户名长度必须在3-20个字符之间")

	if err == nil {
		t.Fatal("ValidationError() 返回 nil")
	}

	if err.Code != "1001" {
		t.Errorf("Code = %s, want 1001", err.Code)
	}

	if err.Category != errors.CategoryValidation {
		t.Errorf("Category = %s, want %s", err.Category, errors.CategoryValidation)
	}

	if err.Level != errors.LevelWarning {
		t.Errorf("Level = %s, want %s", err.Level, errors.LevelWarning)
	}

	if err.Message != "参数无效" {
		t.Errorf("Message = %s, want 参数无效", err.Message)
	}

	if err.Service != "test-service" {
		t.Errorf("Service = %s, want test-service", err.Service)
	}
}

func TestErrorFactory_BusinessError(t *testing.T) {
	factory := errors.NewErrorFactory("test-service")

	err := factory.BusinessError("2003", "用户名已被使用")

	if err == nil {
		t.Fatal("BusinessError() 返回 nil")
	}

	if err.Code != "2003" {
		t.Errorf("Code = %s, want 2003", err.Code)
	}

	if err.Category != errors.CategoryBusiness {
		t.Errorf("Category = %s, want %s", err.Category, errors.CategoryBusiness)
	}

	if err.Level != errors.LevelError {
		t.Errorf("Level = %s, want %s", err.Level, errors.LevelError)
	}
}

func TestErrorFactory_NotFoundError(t *testing.T) {
	factory := errors.NewErrorFactory("test-service")

	tests := []struct {
		name     string
		resource string
		id       string
		wantCode string
	}{
		{
			name:     "资源不存在 - 带ID",
			resource: "用户",
			id:       "123",
			wantCode: "NOT_FOUND", // ErrorFactory 使用字符串错误码
		},
		{
			name:     "资源不存在 - 不带ID",
			resource: "用户",
			id:       "",
			wantCode: "NOT_FOUND", // ErrorFactory 使用字符串错误码
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.NotFoundError(tt.resource, tt.id)

			if err == nil {
				t.Fatal("NotFoundError() 返回 nil")
			}

			if err.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", err.Code, tt.wantCode)
			}

			if err.Category != errors.CategoryBusiness {
				t.Errorf("Category = %s, want %s", err.Category, errors.CategoryBusiness)
			}

			if err.Level != errors.LevelWarning {
				t.Errorf("Level = %s, want %s", err.Level, errors.LevelWarning)
			}

			if err.GetHTTPStatus() != http.StatusNotFound {
				t.Errorf("HTTPStatus = %d, want %d", err.GetHTTPStatus(), http.StatusNotFound)
			}
		})
	}
}

func TestErrorFactory_InternalError(t *testing.T) {
	factory := errors.NewErrorFactory("test-service")

	cause := stderrors.New("database connection failed")
	err := factory.InternalError("5001", "查询用户失败", cause)

	if err == nil {
		t.Fatal("InternalError() 返回 nil")
	}

	if err.Code != "5001" {
		t.Errorf("Code = %s, want 5001", err.Code)
	}

	if err.Category != errors.CategorySystem {
		t.Errorf("Category = %s, want %s", err.Category, errors.CategorySystem)
	}

	if err.Level != errors.LevelError {
		t.Errorf("Level = %s, want %s", err.Level, errors.LevelError)
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap() 应该返回 cause")
	}

	// ErrorFactory 不自动添加堆栈信息
	_ = err.Stack
}

// ============================================================================
// HTTP状态码映射测试
// ============================================================================

func TestUnifiedError_GetHTTPStatus(t *testing.T) {
	tests := []struct {
		name         string
		category     errors.ErrorCategory
		explicitCode int
		wantStatus   int
	}{
		{
			name:       "验证错误",
			category:   errors.CategoryValidation,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "认证错误",
			category:   errors.CategoryAuth,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "业务错误",
			category:   errors.CategoryBusiness,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "外部服务错误",
			category:   errors.CategoryExternal,
			wantStatus: http.StatusBadGateway,
		},
		{
			name:       "网络错误",
			category:   errors.CategoryNetwork,
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:         "显式设置HTTP状态码",
			category:     errors.CategoryValidation,
			explicitCode: http.StatusForbidden,
			wantStatus:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := errors.NewErrorBuilder().
				WithCategory(tt.category).
				WithMessage("测试错误")

			if tt.explicitCode != 0 {
				builder = builder.WithHTTPStatus(tt.explicitCode)
			}

			err := builder.Build()

			if err.GetHTTPStatus() != tt.wantStatus {
				t.Errorf("GetHTTPStatus() = %d, want %d", err.GetHTTPStatus(), tt.wantStatus)
			}
		})
	}
}

// ============================================================================
// 错误转换测试
// ============================================================================

func TestToUnifiedError(t *testing.T) {
	tests := []struct {
		name      string
		service   string
		err       error
		wantNil   bool
		wantCode  string
	}{
		{
			name:    "nil错误",
			service: "test-service",
			err:     nil,
			wantNil: true,
		},
		{
			name:      "统一错误",
			service:   "test-service",
			err:       errors.NewErrorBuilder().WithCode("1001").WithMessage("测试").Build(),
			wantNil:   false,
			wantCode:  "1001",
		},
		{
			name:      "普通错误",
			service:   "test-service",
			err:       stderrors.New("普通错误"),
			wantNil:   false,
			wantCode:  "5001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.ToUnifiedError(tt.service, tt.err)

			if tt.wantNil {
				if result != nil {
					t.Errorf("ToUnifiedError() 应该返回 nil")
				}
				return
			}

			if result == nil {
				t.Fatal("ToUnifiedError() 返回 nil")
			}

			if result.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", result.Code, tt.wantCode)
			}
		})
	}
}

func TestWrapNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		id       string
		wantCode string
	}{
		{
			name:     "带ID的资源不存在",
			resource: "用户",
			id:       "123",
			wantCode: "1004",
		},
		{
			name:     "不带ID的资源不存在",
			resource: "用户",
			id:       "",
			wantCode: "1004",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.WrapNotFoundError(tt.resource, tt.id)

			if err == nil {
				t.Fatal("WrapNotFoundError() 返回 nil")
			}

			if err.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", err.Code, tt.wantCode)
			}

			if err.GetHTTPStatus() != http.StatusNotFound {
				t.Errorf("HTTPStatus = %d, want %d", err.GetHTTPStatus(), http.StatusNotFound)
			}
		})
	}
}

// ============================================================================
// HTTP响应构建测试
// ============================================================================

func TestToHTTPResponse(t *testing.T) {
	tests := []struct {
		name         string
		err          *errors.UnifiedError
		requestID    string
		traceID      string
		wantStatus   int
		wantCode     string
		wantMessage  string
	}{
		{
			name: "成功",
			err: &errors.UnifiedError{
				Code:    "0",
				Message: "成功",
			},
			wantStatus:  http.StatusOK,
			wantCode:    "0",
			wantMessage: "成功",
		},
		{
			name: "验证错误",
			err: errors.NewErrorBuilder().
				WithCode("1001").
				WithMessage("参数无效").
				WithHTTPStatus(http.StatusBadRequest).
				Build(),
			requestID:    "req-123",
			traceID:      "trace-456",
			wantStatus:   http.StatusBadRequest,
			wantCode:     "1001",
			wantMessage:  "参数无效",
		},
		{
			name: "资源不存在",
			err: errors.NewErrorBuilder().
				WithCode("1004").
				WithMessage("资源不存在").
				WithHTTPStatus(http.StatusNotFound).
				Build(),
			wantStatus:  http.StatusNotFound,
			wantCode:    "1004",
			wantMessage: "资源不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, response := errors.ToHTTPResponse(tt.err, tt.requestID, tt.traceID)

			if status != tt.wantStatus {
				t.Errorf("Status = %d, want %d", status, tt.wantStatus)
			}

			if response.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", response.Code, tt.wantCode)
			}

			if response.Message != tt.wantMessage {
				t.Errorf("Message = %s, want %s", response.Message, tt.wantMessage)
			}

			if tt.requestID != "" && response.RequestID != tt.requestID {
				t.Errorf("RequestID = %s, want %s", response.RequestID, tt.requestID)
			}

			if tt.traceID != "" && response.TraceID != tt.traceID {
				t.Errorf("TraceID = %s, want %s", response.TraceID, tt.traceID)
			}
		})
	}
}

// ============================================================================
// 错误码测试
// ============================================================================

func TestErrorCode_GetHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		code       errors.ErrorCode
		wantStatus int
	}{
		{
			name:       "成功",
			code:       errors.Success,
			wantStatus: http.StatusOK,
		},
		{
			name:       "参数无效",
			code:       errors.InvalidParams,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "未授权",
			code:       errors.Unauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "禁止访问",
			code:       errors.Forbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "资源不存在",
			code:       errors.NotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "频率限制",
			code:       errors.RateLimitExceeded,
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name:       "内部错误",
			code:       errors.InternalError,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "外部API错误",
			code:       errors.ExternalAPIError,
			wantStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := errors.GetHTTPStatus(tt.code)
			if status != tt.wantStatus {
				t.Errorf("GetHTTPStatus() = %d, want %d", status, tt.wantStatus)
			}
		})
	}
}

func TestErrorCode_GetDefaultMessage(t *testing.T) {
	tests := []struct {
		name          string
		code          errors.ErrorCode
		wantMessage   string
	}{
		{
			name:        "成功",
			code:        errors.Success,
			wantMessage: "成功",
		},
		{
			name:        "用户不存在",
			code:        errors.UserNotFound,
			wantMessage: "用户不存在",
		},
		{
			name:        "Token过期",
			code:        errors.TokenExpired,
			wantMessage: "登录已过期，请重新登录",
		},
		{
			name:        "书籍不存在",
			code:        errors.BookNotFound,
			wantMessage: "书籍不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := errors.GetDefaultMessage(tt.code)
			if message != tt.wantMessage {
				t.Errorf("GetDefaultMessage() = %s, want %s", message, tt.wantMessage)
			}
		})
	}
}
