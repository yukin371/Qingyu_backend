package infrastructure

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation" // 验证错误
	ErrorTypeNotFound   ErrorType = "not_found"  // 未找到错误
	ErrorTypeConflict   ErrorType = "conflict"   // 冲突错误
	ErrorTypeDuplicate  ErrorType = "duplicate"  // 重复错误
	ErrorTypeInternal   ErrorType = "internal"   // 内部错误
)

// RepositoryError 仓储错误
type RepositoryError struct {
	Type    ErrorType
	Message string
	Detail  error
}

// Error 实现error接口
func (e *RepositoryError) Error() string {
	return e.Message
}

// NewRepositoryError 创建仓储错误
func NewRepositoryError(t ErrorType, msg string, detail error) *RepositoryError {
	return &RepositoryError{
		Type:    t,
		Message: msg,
		Detail:  detail,
	}
}
