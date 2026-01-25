package search

import "fmt"

// SearchError 搜索错误
type SearchError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现 error 接口
func (e *SearchError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 支持错误包装
func (e *SearchError) Unwrap() error {
	return e.Err
}

// 错误代码常量
const (
	ErrCodeInvalidRequest    = "INVALID_REQUEST"
	ErrCodeUnsupportedType   = "UNSUPPORTED_SEARCH_TYPE"
	ErrCodeEngineFailure     = "ENGINE_FAILURE"
	ErrCodeIndexNotFound     = "INDEX_NOT_FOUND"
	ErrCodeDocumentNotFound  = "DOCUMENT_NOT_FOUND"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeQueryParseError   = "QUERY_PARSE_ERROR"
	ErrCodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrCodeCacheFailure      = "CACHE_FAILURE"
	ErrCodeSyncFailure       = "SYNC_FAILURE"
)

// 预定义错误
var (
	ErrInvalidRequest     = &SearchError{Code: ErrCodeInvalidRequest, Message: "Invalid search request"}
	ErrUnauthorized       = &SearchError{Code: ErrCodeUnauthorized, Message: "Authentication required"}
	ErrIndexNotFound      = &SearchError{Code: ErrCodeIndexNotFound, Message: "Search index not found"}
	ErrDocumentNotFound   = &SearchError{Code: ErrCodeDocumentNotFound, Message: "Document not found"}
	ErrUnsupportedType    = &SearchError{Code: ErrCodeUnsupportedType, Message: "Unsupported search type"}
	ErrQueryParseError    = &SearchError{Code: ErrCodeQueryParseError, Message: "Query parse error"}
	ErrRateLimitExceeded  = &SearchError{Code: ErrCodeRateLimitExceeded, Message: "Rate limit exceeded"}
)

// NewSearchError 创建搜索错误
func NewSearchError(code, message string, err error) *SearchError {
	return &SearchError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError 包装错误
func WrapError(err error, code, message string) *SearchError {
	return &SearchError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
