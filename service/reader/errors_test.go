package reader

import (
	"errors"
	"testing"
)

// TestErrorCodeHTTPStatus 测试错误码对应的 HTTP 状态码
func TestErrorCodeHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected int
	}{
		{"ProgressNotFound returns 404", ErrProgressNotFound, 404},
		{"AnnotationNotFound returns 404", ErrAnnotationNotFound, 404},
		{"SettingsNotFound returns 404", ErrSettingsNotFound, 404},
		{"ChapterNotFound returns 404", ErrReaderChapterNotFound, 404},
		{"BookNotFound returns 404", ErrBookNotFound, 404},
		{"InvalidProgress returns 400", ErrInvalidProgress, 400},
		{"InvalidAnnotation returns 400", ErrInvalidAnnotation, 400},
		{"InvalidSettings returns 400", ErrInvalidSettings, 400},
		{"InvalidStatus returns 400", ErrInvalidStatus, 400},
		{"InvalidChapterNumber returns 400", ErrInvalidChapterNumber, 400},
		{"AccessDenied returns 403", ErrReaderAccessDenied, 403},
		{"PermissionDenied returns 403", ErrPermissionDenied, 403},
		{"SyncFailed returns 500", ErrSyncFailed, 500},
		{"InternalError returns 500", ErrInternalError, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.HTTPStatus(); got != tt.expected {
				t.Errorf("ErrorCode.HTTPStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestErrorCodeString 测试错误码的字符串表示
func TestErrorCodeString(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected string
	}{
		{"ProgressNotFound string", ErrProgressNotFound, "PROGRESS_NOT_FOUND"},
		{"AnnotationNotFound string", ErrAnnotationNotFound, "ANNOTATION_NOT_FOUND"},
		{"SettingsNotFound string", ErrSettingsNotFound, "SETTINGS_NOT_FOUND"},
		{"ChapterNotFound string", ErrReaderChapterNotFound, "CHAPTER_NOT_FOUND"},
		{"BookNotFound string", ErrBookNotFound, "BOOK_NOT_FOUND"},
		{"InvalidProgress string", ErrInvalidProgress, "INVALID_PROGRESS"},
		{"InvalidAnnotation string", ErrInvalidAnnotation, "INVALID_ANNOTATION"},
		{"InvalidSettings string", ErrInvalidSettings, "INVALID_SETTINGS"},
		{"InvalidStatus string", ErrInvalidStatus, "INVALID_STATUS"},
		{"InvalidChapterNumber string", ErrInvalidChapterNumber, "INVALID_CHAPTER_NUMBER"},
		{"AccessDenied string", ErrReaderAccessDenied, "ACCESS_DENIED"},
		{"PermissionDenied string", ErrPermissionDenied, "PERMISSION_DENIED"},
		{"SyncFailed string", ErrSyncFailed, "SYNC_FAILED"},
		{"InternalError string", ErrInternalError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.expected {
				t.Errorf("ErrorCode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestReaderErrorError 测试 ReaderError 的 Error() 方法
func TestReaderErrorError(t *testing.T) {
	t.Run("Error with message only", func(t *testing.T) {
		err := NewReaderError(ErrInvalidProgress, "进度值无效")
		expected := "[INVALID_PROGRESS] 进度值无效"
		if got := err.Error(); got != expected {
			t.Errorf("ReaderError.Error() = %v, want %v", got, expected)
		}
	})

	t.Run("Error with field and message", func(t *testing.T) {
		err := InvalidProgress("进度值必须在0-1之间")
		expected := "[INVALID_PROGRESS] progress: 进度值必须在0-1之间"
		if got := err.Error(); got != expected {
			t.Errorf("ReaderError.Error() = %v, want %v", got, expected)
		}
	})

	t.Run("Error without message", func(t *testing.T) {
		err := &ReaderError{Code: ErrInternalError}
		expected := "INTERNAL_ERROR"
		if got := err.Error(); got != expected {
			t.Errorf("ReaderError.Error() = %v, want %v", got, expected)
		}
	})
}

// TestReaderErrorUnwrap 测试 ReaderError 的 Unwrap() 方法
func TestReaderErrorUnwrap(t *testing.T) {
	cause := errors.New("底层错误")
	err := InternalError("内部错误", cause)

	if got := errors.Unwrap(err); got != cause {
		t.Errorf("ReaderError.Unwrap() = %v, want %v", got, cause)
	}
}

// TestReaderErrorWithField 测试 WithField 方法
func TestReaderErrorWithField(t *testing.T) {
	err := InvalidProgress("")
	err = err.WithField("progress")

	if err.Field != "progress" {
		t.Errorf("ReaderError.Field = %v, want %v", err.Field, "progress")
	}
}

// TestReaderErrorWithCause 测试 WithCause 方法
func TestReaderErrorWithCause(t *testing.T) {
	cause := errors.New("底层错误")
	err := InternalError("")
	err = err.WithCause(cause)

	if err.Err != cause {
		t.Errorf("ReaderError.Err = %v, want %v", err.Err, cause)
	}
}

// TestReaderErrorWithMetadata 测试 WithMetadata 方法
func TestReaderErrorWithMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"user_id": "user1",
		"book_id": "book1",
	}
	err := InvalidProgress("")
	err = err.WithMetadata(metadata)

	if len(err.Metadata) != 2 {
		t.Errorf("ReaderError.Metadata length = %v, want %v", len(err.Metadata), 2)
	}
}

// TestReaderErrorIsRetryable 测试 IsRetryable 方法
func TestReaderErrorIsRetryable(t *testing.T) {
	err := InvalidProgress("")
	if err.IsRetryable() {
		t.Errorf("ReaderError.IsRetryable() = true, want false")
	}
}

// TestIsReaderError 测试 IsReaderError 函数
func TestIsReaderError(t *testing.T) {
	t.Run("IsReaderError returns true for ReaderError", func(t *testing.T) {
		err := InvalidProgress("")
		if !IsReaderError(err) {
			t.Errorf("IsReaderError() = false, want true")
		}
	})

	t.Run("IsReaderError returns false for other error", func(t *testing.T) {
		err := errors.New("other error")
		if IsReaderError(err) {
			t.Errorf("IsReaderError() = true, want false")
		}
	})
}

// TestIsErrorCode 测试 IsErrorCode 函数
func TestIsErrorCode(t *testing.T) {
	t.Run("IsErrorCode returns true for matching code", func(t *testing.T) {
		err := InvalidProgress("")
		if !IsErrorCode(err, ErrInvalidProgress) {
			t.Errorf("IsErrorCode() = false, want true")
		}
	})

	t.Run("IsErrorCode returns false for different code", func(t *testing.T) {
		err := InvalidProgress("")
		if IsErrorCode(err, ErrInvalidAnnotation) {
			t.Errorf("IsErrorCode() = true, want false")
		}
	})

	t.Run("IsErrorCode returns false for non-ReaderError", func(t *testing.T) {
		err := errors.New("other error")
		if IsErrorCode(err, ErrInvalidProgress) {
			t.Errorf("IsErrorCode() = true, want false")
		}
	})
}

// Test专用构造函数 测试专用构造函数
func Test专用构造函数(t *testing.T) {
	t.Run("ProgressNotFound", func(t *testing.T) {
		err := ProgressNotFound("user1-book1")
		if err.Code != ErrProgressNotFound {
			t.Errorf("ProgressNotFound() code = %v, want %v", err.Code, ErrProgressNotFound)
		}
		if err.Message == "" {
			t.Errorf("ProgressNotFound() message should not be empty")
		}
	})

	t.Run("AnnotationNotFound", func(t *testing.T) {
		err := AnnotationNotFound("annotation1")
		if err.Code != ErrAnnotationNotFound {
			t.Errorf("AnnotationNotFound() code = %v, want %v", err.Code, ErrAnnotationNotFound)
		}
	})

	t.Run("SettingsNotFound", func(t *testing.T) {
		err := SettingsNotFound()
		if err.Code != ErrSettingsNotFound {
			t.Errorf("SettingsNotFound() code = %v, want %v", err.Code, ErrSettingsNotFound)
		}
	})

	t.Run("ChapterNotFound", func(t *testing.T) {
		err := ChapterNotFound("chapter1")
		if err.Code != ErrReaderChapterNotFound {
			t.Errorf("ChapterNotFound() code = %v, want %v", err.Code, ErrReaderChapterNotFound)
		}
	})

	t.Run("BookNotFound", func(t *testing.T) {
		err := BookNotFound("book1")
		if err.Code != ErrBookNotFound {
			t.Errorf("BookNotFound() code = %v, want %v", err.Code, ErrBookNotFound)
		}
	})

	t.Run("InvalidProgress", func(t *testing.T) {
		err := InvalidProgress("进度值必须在0-1之间")
		if err.Code != ErrInvalidProgress {
			t.Errorf("InvalidProgress() code = %v, want %v", err.Code, ErrInvalidProgress)
		}
		if err.Field != "progress" {
			t.Errorf("InvalidProgress() field = %v, want %v", err.Field, "progress")
		}
	})

	t.Run("InvalidAnnotation", func(t *testing.T) {
		err := InvalidAnnotation("chapter_id")
		if err.Code != ErrInvalidAnnotation {
			t.Errorf("InvalidAnnotation() code = %v, want %v", err.Code, ErrInvalidAnnotation)
		}
		if err.Field != "chapter_id" {
			t.Errorf("InvalidAnnotation() field = %v, want %v", err.Field, "chapter_id")
		}
	})

	t.Run("InvalidSettings", func(t *testing.T) {
		err := InvalidSettings("font_size")
		if err.Code != ErrInvalidSettings {
			t.Errorf("InvalidSettings() code = %v, want %v", err.Code, ErrInvalidSettings)
		}
		if err.Field != "font_size" {
			t.Errorf("InvalidSettings() field = %v, want %v", err.Field, "font_size")
		}
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		err := InvalidStatus("unknown")
		if err.Code != ErrInvalidStatus {
			t.Errorf("InvalidStatus() code = %v, want %v", err.Code, ErrInvalidStatus)
		}
		if err.Field != "status" {
			t.Errorf("InvalidStatus() field = %v, want %v", err.Field, "status")
		}
	})

	t.Run("InvalidChapterNumber", func(t *testing.T) {
		err := InvalidChapterNumber(-1)
		if err.Code != ErrInvalidChapterNumber {
			t.Errorf("InvalidChapterNumber() code = %v, want %v", err.Code, ErrInvalidChapterNumber)
		}
		if err.Field != "chapter_number" {
			t.Errorf("InvalidChapterNumber() field = %v, want %v", err.Field, "chapter_number")
		}
	})

	t.Run("AccessDenied", func(t *testing.T) {
		err := AccessDenied("书籍内容")
		if err.Code != ErrReaderAccessDenied {
			t.Errorf("AccessDenied() code = %v, want %v", err.Code, ErrReaderAccessDenied)
		}
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		err := PermissionDenied("阅读VIP章节")
		if err.Code != ErrPermissionDenied {
			t.Errorf("PermissionDenied() code = %v, want %v", err.Code, ErrPermissionDenied)
		}
	})

	t.Run("SyncFailed", func(t *testing.T) {
		cause := errors.New("网络错误")
		err := SyncFailed("同步失败", cause)
		if err.Code != ErrSyncFailed {
			t.Errorf("SyncFailed() code = %v, want %v", err.Code, ErrSyncFailed)
		}
		if err.Err != cause {
			t.Errorf("SyncFailed() cause = %v, want %v", err.Err, cause)
		}
	})

	t.Run("InternalError", func(t *testing.T) {
		cause := errors.New("数据库错误")
		err := InternalError("内部错误", cause)
		if err.Code != ErrInternalError {
			t.Errorf("InternalError() code = %v, want %v", err.Code, ErrInternalError)
		}
		if err.Err != cause {
			t.Errorf("InternalError() cause = %v, want %v", err.Err, cause)
		}
	})
}
