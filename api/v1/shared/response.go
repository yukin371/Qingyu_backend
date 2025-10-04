package shared

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse 分页响应格式
type PaginatedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Size    int         `json:"size,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse 成功响应辅助函数
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponseWithCode 错误响应辅助函数
func ErrorResponseWithCode(code int, message string, err error) ErrorResponse {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return ErrorResponse{
		Code:    code,
		Message: message,
		Error:   errMsg,
	}
}

// PaginatedResponseHelper 分页响应辅助函数
func PaginatedResponseHelper(data interface{}, total int64, page, size int, message string) PaginatedResponse {
	return PaginatedResponse{
		Code:    200,
		Message: message,
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	}
}
