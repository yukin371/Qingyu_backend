package writer

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 错误码常量
const (
	ErrProjectNotFound      ErrorCode = 40401 // 项目不存在
	ErrDocumentNotFound     ErrorCode = 40402 // 文档不存在
	ErrCommentNotFound      ErrorCode = 40403 // 批注不存在
	ErrCharacterNotFound    ErrorCode = 40404 // 角色不存在
	ErrLocationNotFound     ErrorCode = 40405 // 地点不存在
	ErrTimelineNotFound     ErrorCode = 40406 // 时间线不存在
	ErrVersionNotFound      ErrorCode = 40407 // 版本不存在
	ErrPublicationNotFound  ErrorCode = 40408 // 发布记录不存在
	ErrExportTaskNotFound   ErrorCode = 40409 // 导出任务不存在

	ErrInvalidInput         ErrorCode = 40001 // 输入参数无效
	ErrInvalidProjectID     ErrorCode = 40002 // 无效的项目ID
	ErrInvalidDocumentID    ErrorCode = 40003 // 无效的文档ID
	ErrInvalidContent       ErrorCode = 40004 // 内容无效
	ErrInvalidVersion       ErrorCode = 40005 // 无效的版本号
	ErrInvalidExportFormat  ErrorCode = 40006 // 无效的导出格式
	ErrInvalidRelationType  ErrorCode = 40007 // 无效的关系类型

	ErrProjectAlreadyExists ErrorCode = 40901 // 项目已存在
	ErrDocumentAlreadyExists ErrorCode = 40902 // 文档已存在
	ErrNameAlreadyExists    ErrorCode = 40903 // 名称已存在

	ErrVersionConflict      ErrorCode = 40904 // 版本冲突
	ErrEditConflict         ErrorCode = 40905 // 编辑冲突

	ErrUnauthorized         ErrorCode = 40101 // 未授权
	ErrForbidden            ErrorCode = 40301 // 禁止访问
	ErrNoPermission         ErrorCode = 40302 // 无权限

	ErrPublishFailed        ErrorCode = 50001 // 发布失败
	ErrExportFailed         ErrorCode = 50002 // 导出失败
	ErrInternalError        ErrorCode = 50003 // 内部错误
	ErrStorageError         ErrorCode = 50004 // 存储错误
	ErrExternalServiceError ErrorCode = 50005 // 外部服务错误
)

// HTTPStatus 返回错误码对应的 HTTP 状态码
func (c ErrorCode) HTTPStatus() int {
	switch {
	case c >= 40000 && c < 40100:
		return 400
	case c >= 40100 && c < 40200:
		return 401
	case c >= 40300 && c < 40400:
		return 403
	case c >= 40400 && c < 40500:
		return 404
	case c >= 40900 && c < 41000:
		return 409
	case c >= 50000:
		return 500
	default:
		return http.StatusInternalServerError
	}
}

// String 返回错误码的字符串表示
func (c ErrorCode) String() string {
	switch c {
	case ErrProjectNotFound:
		return "PROJECT_NOT_FOUND"
	case ErrDocumentNotFound:
		return "DOCUMENT_NOT_FOUND"
	case ErrCommentNotFound:
		return "COMMENT_NOT_FOUND"
	case ErrCharacterNotFound:
		return "CHARACTER_NOT_FOUND"
	case ErrLocationNotFound:
		return "LOCATION_NOT_FOUND"
	case ErrTimelineNotFound:
		return "TIMELINE_NOT_FOUND"
	case ErrVersionNotFound:
		return "VERSION_NOT_FOUND"
	case ErrPublicationNotFound:
		return "PUBLICATION_NOT_FOUND"
	case ErrExportTaskNotFound:
		return "EXPORT_TASK_NOT_FOUND"

	case ErrInvalidInput:
		return "INVALID_INPUT"
	case ErrInvalidProjectID:
		return "INVALID_PROJECT_ID"
	case ErrInvalidDocumentID:
		return "INVALID_DOCUMENT_ID"
	case ErrInvalidContent:
		return "INVALID_CONTENT"
	case ErrInvalidVersion:
		return "INVALID_VERSION"
	case ErrInvalidExportFormat:
		return "INVALID_EXPORT_FORMAT"
	case ErrInvalidRelationType:
		return "INVALID_RELATION_TYPE"

	case ErrProjectAlreadyExists:
		return "PROJECT_ALREADY_EXISTS"
	case ErrDocumentAlreadyExists:
		return "DOCUMENT_ALREADY_EXISTS"
	case ErrNameAlreadyExists:
		return "NAME_ALREADY_EXISTS"

	case ErrVersionConflict:
		return "VERSION_CONFLICT"
	case ErrEditConflict:
		return "EDIT_CONFLICT"

	case ErrUnauthorized:
		return "UNAUTHORIZED"
	case ErrForbidden:
		return "FORBIDDEN"
	case ErrNoPermission:
		return "NO_PERMISSION"

	case ErrPublishFailed:
		return "PUBLISH_FAILED"
	case ErrExportFailed:
		return "EXPORT_FAILED"
	case ErrInternalError:
		return "INTERNAL_ERROR"
	case ErrStorageError:
		return "STORAGE_ERROR"
	case ErrExternalServiceError:
		return "EXTERNAL_SERVICE_ERROR"

	default:
		return fmt.Sprintf("UNKNOWN_ERROR(%d)", c)
	}
}

// WriterError Writer 模块结构化错误
type WriterError struct {
	Code     ErrorCode
	Field    string
	Message  string
	Err      error
	Metadata map[string]interface{}
}

// Error 实现 error 接口
func (e *WriterError) Error() string {
	if e.Message != "" {
		if e.Field != "" {
			return fmt.Sprintf("[%s] %s: %s", e.Code.String(), e.Field, e.Message)
		}
		return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
	}
	return e.Code.String()
}

// Unwrap 支持 errors.Unwrap
func (e *WriterError) Unwrap() error {
	return e.Err
}

// WithField 设置字段名
func (e *WriterError) WithField(field string) *WriterError {
	e.Field = field
	return e
}

// WithCause 设置底层错误
func (e *WriterError) WithCause(err error) *WriterError {
	e.Err = err
	return e
}

// WithMetadata 设置元数据
func (e *WriterError) WithMetadata(metadata map[string]interface{}) *WriterError {
	e.Metadata = metadata
	return e
}

// IsRetryable 判断错误是否可重试
func (e *WriterError) IsRetryable() bool {
	// 版本冲突和编辑冲突可以重试
	if e.Code == ErrVersionConflict || e.Code == ErrEditConflict {
		return true
	}
	// 外部服务错误可以重试
	if e.Code == ErrExternalServiceError {
		return true
	}
	return false
}

// NewWriterError 创建新的 Writer 错误
func NewWriterError(code ErrorCode, message string) *WriterError {
	return &WriterError{
		Code:    code,
		Message: message,
	}
}

// IsWriterError 判断是否是 WriterError 类型
func IsWriterError(err error) bool {
	_, ok := err.(*WriterError)
	return ok
}

// IsErrorCode 判断错误是否是指定错误码
func IsErrorCode(err error, code ErrorCode) bool {
	if writerErr, ok := err.(*WriterError); ok {
		return writerErr.Code == code
	}
	return false
}

// ============================================================================
// 专用构造函数
// ============================================================================

// ProjectNotFound 创建项目不存在错误
func ProjectNotFound(identifier string) *WriterError {
	message := "项目不存在"
	if identifier != "" {
		message = fmt.Sprintf("项目 %s 不存在", identifier)
	}
	return &WriterError{Code: ErrProjectNotFound, Message: message}
}

// DocumentNotFound 创建文档不存在错误
func DocumentNotFound(identifier string) *WriterError {
	message := "文档不存在"
	if identifier != "" {
		message = fmt.Sprintf("文档 %s 不存在", identifier)
	}
	return &WriterError{Code: ErrDocumentNotFound, Message: message}
}

// CommentNotFound 创建批注不存在错误
func CommentNotFound() *WriterError {
	return &WriterError{Code: ErrCommentNotFound, Message: "批注不存在"}
}

// CharacterNotFound 创建角色不存在错误
func CharacterNotFound() *WriterError {
	return &WriterError{Code: ErrCharacterNotFound, Message: "角色不存在"}
}

// LocationNotFound 创建地点不存在错误
func LocationNotFound() *WriterError {
	return &WriterError{Code: ErrLocationNotFound, Message: "地点不存在"}
}

// TimelineNotFound 创建时间线不存在错误
func TimelineNotFound() *WriterError {
	return &WriterError{Code: ErrTimelineNotFound, Message: "时间线不存在"}
}

// VersionNotFound 创建版本不存在错误
func VersionNotFound() *WriterError {
	return &WriterError{Code: ErrVersionNotFound, Message: "版本不存在"}
}

// PublicationNotFound 创建发布记录不存在错误
func PublicationNotFound() *WriterError {
	return &WriterError{Code: ErrPublicationNotFound, Message: "发布记录不存在"}
}

// ExportTaskNotFound 创建导出任务不存在错误
func ExportTaskNotFound() *WriterError {
	return &WriterError{Code: ErrExportTaskNotFound, Message: "导出任务不存在"}
}

// InvalidInput 创建输入参数无效错误
func InvalidInput(field string) *WriterError {
	message := "输入参数无效"
	if field != "" {
		message = fmt.Sprintf("%s 输入参数无效", field)
	}
	return &WriterError{Code: ErrInvalidInput, Field: field, Message: message}
}

// InvalidContent 创建内容无效错误
func InvalidContent(message string) *WriterError {
	msg := message
	if msg == "" {
		msg = "内容无效"
	}
	return &WriterError{Code: ErrInvalidContent, Field: "content", Message: msg}
}

// InvalidVersion 创建无效版本号错误
func InvalidVersion() *WriterError {
	return &WriterError{Code: ErrInvalidVersion, Field: "version", Message: "版本号无效"}
}

// InvalidExportFormat 创建无效导出格式错误
func InvalidExportFormat(format string) *WriterError {
	message := "无效的导出格式"
	if format != "" {
		message = fmt.Sprintf("不支持的导出格式: %s", format)
	}
	return &WriterError{Code: ErrInvalidExportFormat, Field: "format", Message: message}
}

// InvalidRelationType 创建无效关系类型错误
func InvalidRelationType(relationType string) *WriterError {
	message := "无效的关系类型"
	if relationType != "" {
		message = fmt.Sprintf("不支持的关系类型: %s", relationType)
	}
	return &WriterError{Code: ErrInvalidRelationType, Field: "relation_type", Message: message}
}

// ProjectAlreadyExists 创建项目已存在错误
func ProjectAlreadyExists(field string) *WriterError {
	var message string
	var fieldName string

	switch field {
	case "name":
		message = "项目名称已存在"
		fieldName = "name"
	default:
		message = "项目已存在"
		fieldName = field
	}

	return &WriterError{Code: ErrProjectAlreadyExists, Field: fieldName, Message: message}
}

// DocumentAlreadyExists 创建文档已存在错误
func DocumentAlreadyExists() *WriterError {
	return &WriterError{Code: ErrDocumentAlreadyExists, Message: "文档已存在"}
}

// NameAlreadyExists 创建名称已存在错误
func NameAlreadyExists(resourceType, name string) *WriterError {
	message := "名称已存在"
	if resourceType != "" && name != "" {
		message = fmt.Sprintf("%s '%s' 已存在", resourceType, name)
	}
	return &WriterError{Code: ErrNameAlreadyExists, Message: message}
}

// VersionConflict 创建版本冲突错误
func VersionConflict(expected, current int) *WriterError {
	message := "版本冲突，请刷新后重试"
	if expected > 0 && current > 0 {
		message = fmt.Sprintf("版本冲突：期望版本 %d，当前版本 %d", expected, current)
	}
	return &WriterError{Code: ErrVersionConflict, Message: message}
}

// EditConflict 创建编辑冲突错误
func EditConflict() *WriterError {
	return &WriterError{Code: ErrEditConflict, Message: "编辑冲突，请刷新后重试"}
}

// Unauthorized 创建未授权错误
func Unauthorized() *WriterError {
	return &WriterError{Code: ErrUnauthorized, Message: "用户未登录"}
}

// Forbidden 创建禁止访问错误
func Forbidden(action string) *WriterError {
	message := "禁止访问"
	if action != "" {
		message = fmt.Sprintf("无权%s", action)
	}
	return &WriterError{Code: ErrForbidden, Message: message}
}

// NoPermission 创建无权限错误
func NoPermission(action string) *WriterError {
	message := "权限不足"
	if action != "" {
		message = fmt.Sprintf("没有%s的权限", action)
	}
	return &WriterError{Code: ErrNoPermission, Message: message}
}

// PublishFailed 创建发布失败错误
func PublishFailed(message string, cause ...error) *WriterError {
	msg := message
	if msg == "" {
		msg = "发布失败"
	}
	err := &WriterError{
		Code:    ErrPublishFailed,
		Message: msg,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}

// ExportFailed 创建导出失败错误
func ExportFailed(message string, cause ...error) *WriterError {
	msg := message
	if msg == "" {
		msg = "导出失败"
	}
	err := &WriterError{
		Code:    ErrExportFailed,
		Message: msg,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}

// InternalError 创建内部错误
func InternalError(message string, cause ...error) *WriterError {
	msg := message
	if msg == "" {
		msg = "内部错误"
	}
	err := &WriterError{
		Code:    ErrInternalError,
		Message: msg,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}

// StorageError 创建存储错误
func StorageError(message string, cause ...error) *WriterError {
	msg := message
	if msg == "" {
		msg = "存储错误"
	}
	err := &WriterError{
		Code:    ErrStorageError,
		Message: msg,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}

// ExternalServiceError 创建外部服务错误
func ExternalServiceError(service, message string, cause ...error) *WriterError {
	msg := message
	if msg == "" {
		msg = "外部服务调用失败"
	}
	if service != "" {
		msg = fmt.Sprintf("[%s] %s", service, msg)
	}
	err := &WriterError{
		Code:    ErrExternalServiceError,
		Message: msg,
	}
	if len(cause) > 0 && cause[0] != nil {
		err.Err = cause[0]
	}
	return err
}
