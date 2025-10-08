package shared

import (
	"time"
)

// APIResponse 统一API响应格式
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`            // 响应时间戳
	RequestID string      `json:"request_id,omitempty"` // 请求ID（用于追踪）
}

// PaginatedResponse 分页响应格式
type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Timestamp  int64       `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
	Pagination *Pagination `json:"pagination"` // 分页信息
}

// Pagination 分页信息
type Pagination struct {
	Total       int64 `json:"total"`        // 总记录数
	Page        int   `json:"page"`         // 当前页码
	PageSize    int   `json:"page_size"`    // 每页大小
	TotalPages  int   `json:"total_pages"`  // 总页数
	HasNext     bool  `json:"has_next"`     // 是否有下一页
	HasPrevious bool  `json:"has_previous"` // 是否有上一页
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// SuccessResponse 成功响应辅助函数
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// SuccessResponseWithRequestID 带RequestID的成功响应
func SuccessResponseWithRequestID(data interface{}, message, requestID string) APIResponse {
	return APIResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: requestID,
	}
}

// ErrorResponseWithCode 错误响应辅助函数
func ErrorResponseWithCode(code int, message string, err error) ErrorResponse {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return ErrorResponse{
		Code:      code,
		Message:   message,
		Error:     errMsg,
		Timestamp: time.Now().Unix(),
	}
}

// ErrorResponseWithRequestID 带RequestID的错误响应
func ErrorResponseWithRequestID(code int, message string, err error, requestID string) ErrorResponse {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return ErrorResponse{
		Code:      code,
		Message:   message,
		Error:     errMsg,
		Timestamp: time.Now().Unix(),
		RequestID: requestID,
	}
}

// PaginatedResponseHelper 分页响应辅助函数
func PaginatedResponseHelper(data interface{}, total int64, page, pageSize int, message string) PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		Pagination: &Pagination{
			Total:       total,
			Page:        page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}
}

// PaginatedResponseWithRequestID 带RequestID的分页响应
func PaginatedResponseWithRequestID(data interface{}, total int64, page, pageSize int, message, requestID string) PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: requestID,
		Pagination: &Pagination{
			Total:       total,
			Page:        page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}
}

// NewPagination 创建分页信息
func NewPagination(total int64, page, pageSize int) *Pagination {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &Pagination{
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}
