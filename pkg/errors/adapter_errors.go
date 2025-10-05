package errors

import "fmt"

// AdapterErrorType 适配器错误类型
type AdapterErrorType string

const (
	AdapterErrorInit        AdapterErrorType = "INIT_ERROR"
	AdapterErrorConfig      AdapterErrorType = "CONFIG_ERROR"
	AdapterErrorValidation  AdapterErrorType = "VALIDATION_ERROR"
	AdapterErrorNotFound    AdapterErrorType = "NOT_FOUND"
	AdapterErrorUnavailable AdapterErrorType = "UNAVAILABLE"
	AdapterErrorInternal    AdapterErrorType = "INTERNAL_ERROR"
)

// AdapterError 适配器错误
type AdapterError struct {
	Type      AdapterErrorType
	Message   string
	AdapterID string
	Provider  string
	Cause     error
}

func (e *AdapterError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[Adapter:%s/%s] %s: %s (cause: %v)",
			e.Provider, e.AdapterID, e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[Adapter:%s/%s] %s: %s",
		e.Provider, e.AdapterID, e.Type, e.Message)
}

// NewAdapterError 创建适配器错误
func NewAdapterError(adapterID, provider string, errorType AdapterErrorType, message string, cause error) *AdapterError {
	return &AdapterError{
		Type:      errorType,
		Message:   message,
		AdapterID: adapterID,
		Provider:  provider,
		Cause:     cause,
	}
}
