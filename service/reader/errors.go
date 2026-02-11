package reader

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 错误码常量
const (
	ErrProgressNotFound   ErrorCode = 40401 // 阅读进度不存在
	ErrAnnotationNotFound ErrorCode = 40402 // 标注不存在
	ErrSettingsNotFound   ErrorCode = 40403 // 阅读设置不存在
	// ErrChapterNotFound 已在 chapter_service.go 中定义，使用 ErrReaderChapterNotFound
	ErrReaderChapterNotFound ErrorCode = 40404 // 章节不存在
	ErrBookNotFound          ErrorCode = 40405 // 书籍不存在

	ErrInvalidProgress      ErrorCode = 40001 // 无效的阅读进度
	ErrInvalidAnnotation    ErrorCode = 40002 // 无效的标注
	ErrInvalidSettings      ErrorCode = 40003 // 无效的阅读设置
	ErrInvalidStatus        ErrorCode = 40004 // 无效的书籍状态
	ErrInvalidChapterNumber ErrorCode = 40005 // 无效的章节号

	// ErrAccessDenied 已在 chapter_service.go 中定义，使用 ErrReaderAccessDenied
	ErrReaderAccessDenied ErrorCode = 40301 // 访问被拒绝
	ErrPermissionDenied   ErrorCode = 40302 // 权限不足

	ErrSyncFailed    ErrorCode = 50001 // 同步失败
	ErrInternalError ErrorCode = 50002 // 内部错误
)

// HTTPStatus 返回错误码对应的 HTTP 状态码
func (c ErrorCode) HTTPStatus() int {
	switch {
	case c >= 40000 && c < 40100:
		return 400
	case c >= 40300 && c < 40400:
		return 403
	case c >= 40400 && c < 40500:
		return 404
	case c >= 50000:
		return 500
	default:
		return http.StatusInternalServerError
	}
}

// String 返回错误码的字符串表示
func (c ErrorCode) String() string {
	switch c {
	case ErrProgressNotFound:
		return "PROGRESS_NOT_FOUND"
	case ErrAnnotationNotFound:
		return "ANNOTATION_NOT_FOUND"
	case ErrSettingsNotFound:
		return "SETTINGS_NOT_FOUND"
	case ErrReaderChapterNotFound:
		return "CHAPTER_NOT_FOUND"
	case ErrBookNotFound:
		return "BOOK_NOT_FOUND"
	case ErrInvalidProgress:
		return "INVALID_PROGRESS"
	case ErrInvalidAnnotation:
		return "INVALID_ANNOTATION"
	case ErrInvalidSettings:
		return "INVALID_SETTINGS"
	case ErrInvalidStatus:
		return "INVALID_STATUS"
	case ErrInvalidChapterNumber:
		return "INVALID_CHAPTER_NUMBER"
	case ErrReaderAccessDenied:
		return "ACCESS_DENIED"
	case ErrPermissionDenied:
		return "PERMISSION_DENIED"
	case ErrSyncFailed:
		return "SYNC_FAILED"
	case ErrInternalError:
		return "INTERNAL_ERROR"
	default:
		return fmt.Sprintf("UNKNOWN_ERROR(%d)", c)
	}
}

// ReaderError Reader模块结构化错误
type ReaderError struct {
	Code     ErrorCode
	Field    string
	Message  string
	Err      error
	Metadata map[string]interface{}
}

// Error 实现 error 接口
func (e *ReaderError) Error() string {
	if e.Message != "" {
		if e.Field != "" {
			return fmt.Sprintf("[%s] %s: %s", e.Code.String(), e.Field, e.Message)
		}
		return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
	}
	return e.Code.String()
}

// Unwrap 支持 errors.Unwrap
func (e *ReaderError) Unwrap() error {
	return e.Err
}

// WithField 设置字段名
func (e *ReaderError) WithField(field string) *ReaderError {
	e.Field = field
	return e
}

// WithCause 设置底层错误
func (e *ReaderError) WithCause(err error) *ReaderError {
	e.Err = err
	return e
}

// WithMetadata 设置元数据
func (e *ReaderError) WithMetadata(metadata map[string]interface{}) *ReaderError {
	e.Metadata = metadata
	return e
}

// IsRetryable 判断错误是否可重试
func (e *ReaderError) IsRetryable() bool {
	// Reader 错误通常不可重试
	return false
}

// NewReaderError 创建新的 Reader 错误
func NewReaderError(code ErrorCode, message string) *ReaderError {
	return &ReaderError{
		Code:    code,
		Message: message,
	}
}

// IsReaderError 判断是否是 ReaderError 类型
func IsReaderError(err error) bool {
	_, ok := err.(*ReaderError)
	return ok
}

// IsErrorCode 判断错误是否是指定错误码
func IsErrorCode(err error, code ErrorCode) bool {
	if readerErr, ok := err.(*ReaderError); ok {
		return readerErr.Code == code
	}
	return false
}

// ============================================================================
// 专用构造函数
// ============================================================================

// ProgressNotFound 创建阅读进度不存在错误
func ProgressNotFound(identifier string) *ReaderError {
	message := "阅读进度不存在"
	if identifier != "" {
		message = fmt.Sprintf("阅读进度 %s 不存在", identifier)
	}
	return &ReaderError{Code: ErrProgressNotFound, Message: message}
}

// AnnotationNotFound 创建标注不存在错误
func AnnotationNotFound(identifier string) *ReaderError {
	message := "标注不存在"
	if identifier != "" {
		message = fmt.Sprintf("标注 %s 不存在", identifier)
	}
	return &ReaderError{Code: ErrAnnotationNotFound, Message: message}
}

// SettingsNotFound 创建阅读设置不存在错误
func SettingsNotFound() *ReaderError {
	return &ReaderError{Code: ErrSettingsNotFound, Message: "阅读设置不存在"}
}

// ChapterNotFound 创建章节不存在错误
func ChapterNotFound(identifier string) *ReaderError {
	message := "章节不存在"
	if identifier != "" {
		message = fmt.Sprintf("章节 %s 不存在", identifier)
	}
	return &ReaderError{Code: ErrReaderChapterNotFound, Message: message}
}

// BookNotFound 创建书籍不存在错误
func BookNotFound(identifier string) *ReaderError {
	message := "书籍不存在"
	if identifier != "" {
		message = fmt.Sprintf("书籍 %s 不存在", identifier)
	}
	return &ReaderError{Code: ErrBookNotFound, Message: message}
}

// InvalidProgress 创建无效阅读进度错误
func InvalidProgress(message string) *ReaderError {
	msg := message
	if msg == "" {
		msg = "阅读进度值必须在0-1之间"
	}
	return &ReaderError{Code: ErrInvalidProgress, Field: "progress", Message: msg}
}

// InvalidAnnotation 创建无效标注错误
func InvalidAnnotation(field string) *ReaderError {
	message := "标注参数无效"
	if field != "" {
		message = fmt.Sprintf("%s 参数无效", field)
	}
	return &ReaderError{Code: ErrInvalidAnnotation, Field: field, Message: message}
}

// InvalidSettings 创建无效阅读设置错误
func InvalidSettings(field string) *ReaderError {
	message := "阅读设置无效"
	if field != "" {
		message = fmt.Sprintf("%s 设置无效", field)
	}
	return &ReaderError{Code: ErrInvalidSettings, Field: field, Message: message}
}

// InvalidStatus 创建无效状态错误
func InvalidStatus(status string) *ReaderError {
	message := "无效的书籍状态"
	if status != "" {
		message = fmt.Sprintf("无效的书籍状态: %s", status)
	}
	return &ReaderError{Code: ErrInvalidStatus, Field: "status", Message: message}
}

// InvalidChapterNumber 创建无效章节号错误
func InvalidChapterNumber(chapterNum int) *ReaderError {
	return &ReaderError{
		Code:    ErrInvalidChapterNumber,
		Field:   "chapter_number",
		Message: fmt.Sprintf("无效的章节号: %d", chapterNum),
	}
}

// AccessDenied 创建访问被拒绝错误
func AccessDenied(resource string) *ReaderError {
	message := "访问被拒绝"
	if resource != "" {
		message = fmt.Sprintf("没有访问 %s 的权限", resource)
	}
	return &ReaderError{Code: ErrReaderAccessDenied, Message: message}
}

// PermissionDenied 创建权限不足错误
func PermissionDenied(action string) *ReaderError {
	message := "权限不足"
	if action != "" {
		message = fmt.Sprintf("没有%s的权限", action)
	}
	return &ReaderError{Code: ErrPermissionDenied, Message: message}
}

// SyncFailed 创建同步失败错误
func SyncFailed(message string, cause ...error) *ReaderError {
	err := &ReaderError{
		Code:    ErrSyncFailed,
		Message: message,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}

// InternalError 创建内部错误
func InternalError(message string, cause ...error) *ReaderError {
	err := &ReaderError{
		Code:    ErrInternalError,
		Message: message,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}
