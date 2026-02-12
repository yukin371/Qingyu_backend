package writer

import (
	"errors"
	"testing"
)

// TestErrorCodeString 测试错误码字符串表示
func TestErrorCodeString(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected string
	}{
		{"ProjectNotFound", ErrProjectNotFound, "PROJECT_NOT_FOUND"},
		{"DocumentNotFound", ErrDocumentNotFound, "DOCUMENT_NOT_FOUND"},
		{"InvalidInput", ErrInvalidInput, "INVALID_INPUT"},
		{"VersionConflict", ErrVersionConflict, "VERSION_CONFLICT"},
		{"Unauthorized", ErrUnauthorized, "UNAUTHORIZED"},
		{"Forbidden", ErrForbidden, "FORBIDDEN"},
		{"PublishFailed", ErrPublishFailed, "PUBLISH_FAILED"},
		{"ExportFailed", ErrExportFailed, "EXPORT_FAILED"},
		{"InternalError", ErrInternalError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.expected {
				t.Errorf("ErrorCode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestErrorCodeHTTPStatus 测试错误码对应的 HTTP 状态码
func TestErrorCodeHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected int
	}{
		{"InvalidInput", ErrInvalidInput, 400},
		{"Unauthorized", ErrUnauthorized, 401},
		{"Forbidden", ErrForbidden, 403},
		{"ProjectNotFound", ErrProjectNotFound, 404},
		{"VersionConflict", ErrVersionConflict, 409},
		{"InternalError", ErrInternalError, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.HTTPStatus(); got != tt.expected {
				t.Errorf("ErrorCode.HTTPStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestWriterError_Error 测试 WriterError 错误消息
func TestWriterError_Error(t *testing.T) {
	err := NewWriterError(ErrProjectNotFound, "项目不存在")
	if got := err.Error(); got != "[PROJECT_NOT_FOUND] 项目不存在" {
		t.Errorf("WriterError.Error() = %v, want %v", got, "[PROJECT_NOT_FOUND] 项目不存在")
	}

	// 带 Field 的错误
	err = InvalidInput("title")
	if got := err.Error(); got != "[INVALID_INPUT] title: title 输入参数无效" {
		t.Errorf("WriterError.Error() with field = %v", got)
	}
}

// TestWriterError_WithField 测试 WithField 方法
func TestWriterError_WithField(t *testing.T) {
	err := NewWriterError(ErrInvalidInput, "输入无效")
	err = err.WithField("title")

	if err.Field != "title" {
		t.Errorf("WithField() failed, Field = %v, want %v", err.Field, "title")
	}
}

// TestWriterError_WithCause 测试 WithCause 方法
func TestWriterError_WithCause(t *testing.T) {
	cause := errors.New("底层错误")
	err := NewWriterError(ErrInternalError, "内部错误")
	err = err.WithCause(cause)

	if err.Err != cause {
		t.Errorf("WithCause() failed, Err = %v, want %v", err.Err, cause)
	}
}

// TestWriterError_Unwrap 测试 Unwrap 方法
func TestWriterError_Unwrap(t *testing.T) {
	cause := errors.New("底层错误")
	err := NewWriterError(ErrInternalError, "内部错误")
	err = err.WithCause(cause)

	if unwrapped := errors.Unwrap(err); unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

// TestWriterError_WithMetadata 测试 WithMetadata 方法
func TestWriterError_WithMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"project_id": "123",
		"user_id":    "456",
	}

	err := NewWriterError(ErrForbidden, "禁止访问")
	err = err.WithMetadata(metadata)

	if len(err.Metadata) != 2 {
		t.Errorf("WithMetadata() failed, metadata length = %v, want 2", len(err.Metadata))
	}
}

// TestWriterError_IsRetryable 测试 IsRetryable 方法
func TestWriterError_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      *WriterError
		expected bool
	}{
		{"VersionConflict", VersionConflict(1, 2), true},
		{"EditConflict", EditConflict(), true},
		{"ExternalServiceError", ExternalServiceError("bookstore", "调用失败"), true},
		{"ProjectNotFound", ProjectNotFound("123"), false},
		{"InvalidInput", InvalidInput("title"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsRetryable(); got != tt.expected {
				t.Errorf("WriterError.IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestIsWriterError 测试 IsWriterError 函数
func TestIsWriterError(t *testing.T) {
	writerErr := NewWriterError(ErrProjectNotFound, "项目不存在")
	standardErr := errors.New("标准错误")

	if !IsWriterError(writerErr) {
		t.Error("IsWriterError(writerErr) returned false, want true")
	}

	if IsWriterError(standardErr) {
		t.Error("IsWriterError(standardErr) returned true, want false")
	}
}

// TestIsErrorCode 测试 IsErrorCode 函数
func TestIsErrorCode(t *testing.T) {
	err := ProjectNotFound("123")

	if !IsErrorCode(err, ErrProjectNotFound) {
		t.Error("IsErrorCode() returned false, want true")
	}

	if IsErrorCode(err, ErrDocumentNotFound) {
		t.Error("IsErrorCode() returned true, want false")
	}
}

// TestProjectNotFound 测试 ProjectNotFound 构造函数
func TestProjectNotFound(t *testing.T) {
	// 无参数
	err := ProjectNotFound("")
	if err.Code != ErrProjectNotFound {
		t.Errorf("ProjectNotFound() wrong code, got %v", err.Code)
	}

	// 带参数
	err = ProjectNotFound("project-123")
	if err.Message != "项目 project-123 不存在" {
		t.Errorf("ProjectNotFound() wrong message, got %v", err.Message)
	}
}

// TestInvalidInput 测试 InvalidInput 构造函数
func TestInvalidInput(t *testing.T) {
	err := InvalidInput("title")

	if err.Code != ErrInvalidInput {
		t.Errorf("InvalidInput() wrong code, got %v", err.Code)
	}

	if err.Field != "title" {
		t.Errorf("InvalidInput() wrong field, got %v", err.Field)
	}
}

// TestVersionConflict 测试 VersionConflict 构造函数
func TestVersionConflict(t *testing.T) {
	err := VersionConflict(1, 2)

	if err.Code != ErrVersionConflict {
		t.Errorf("VersionConflict() wrong code, got %v", err.Code)
	}

	// 检查消息包含版本信息
	if err.Message == "" {
		t.Error("VersionConflict() empty message")
	}
}

// TestForbidden 测试 Forbidden 构造函数
func TestForbidden(t *testing.T) {
	// 无参数
	err := Forbidden("")
	if err.Code != ErrForbidden {
		t.Errorf("Forbidden() wrong code, got %v", err.Code)
	}

	// 带参数
	err = Forbidden("编辑此项目")
	if err.Message != "无权编辑此项目" {
		t.Errorf("Forbidden() wrong message, got %v", err.Message)
	}
}

// TestPublishFailed 测试 PublishFailed 构造函数
func TestPublishFailed(t *testing.T) {
	cause := errors.New("书城API错误")

	// 无原因
	err := PublishFailed("发布失败")
	if err.Code != ErrPublishFailed {
		t.Errorf("PublishFailed() wrong code, got %v", err.Code)
	}

	// 带原因
	err = PublishFailed("发布失败", cause)
	if err.Err != cause {
		t.Errorf("PublishFailed() with cause failed, Err = %v, want %v", err.Err, cause)
	}
}

// TestExportFailed 测试 ExportFailed 构造函数
func TestExportFailed(t *testing.T) {
	err := ExportFailed("导出失败")

	if err.Code != ErrExportFailed {
		t.Errorf("ExportFailed() wrong code, got %v", err.Code)
	}
}

// TestInternalError 测试 InternalError 构造函数
func TestInternalError(t *testing.T) {
	cause := errors.New("数据库错误")

	err := InternalError("内部错误", cause)

	if err.Code != ErrInternalError {
		t.Errorf("InternalError() wrong code, got %v", err.Code)
	}

	if err.Err != cause {
		t.Errorf("InternalError() with cause failed, Err = %v, want %v", err.Err, cause)
	}
}

// TestExternalServiceError 测试 ExternalServiceError 构造函数
func TestExternalServiceError(t *testing.T) {
	// 不指定服务
	err := ExternalServiceError("", "调用失败")
	if err.Code != ErrExternalServiceError {
		t.Errorf("ExternalServiceError() wrong code, got %v", err.Code)
	}

	// 指定服务
	err = ExternalServiceError("bookstore", "书城API调用失败")
	if err.Message != "[bookstore] 书城API调用失败" {
		t.Errorf("ExternalServiceError() wrong message, got %v", err.Message)
	}
}

// TestStorageError 测试 StorageError 构造函数
func TestStorageError(t *testing.T) {
	cause := errors.New("MinIO错误")

	err := StorageError("文件上传失败", cause)

	if err.Code != ErrStorageError {
		t.Errorf("StorageError() wrong code, got %v", err.Code)
	}

	if err.Err != cause {
		t.Errorf("StorageError() with cause failed, Err = %v, want %v", err.Err, cause)
	}
}
