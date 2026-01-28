package shared

import (
	"net/http"
	"time"

	apperrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/internal/middleware/builtin"

	"github.com/gin-gonic/gin"
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
	Debug     string `json:"debug,omitempty"` // 调试信息，仅开发环境使用
	Timestamp int64  `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// SuccessResponse 成功响应辅助函数
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

// SuccessResponseWithRequestID 带RequestID的成功响应
func SuccessResponseWithRequestID(data interface{}, message, requestID string) APIResponse {
	return APIResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
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
		Timestamp: time.Now().UnixMilli(),
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
		Timestamp: time.Now().UnixMilli(),
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
		Timestamp: time.Now().UnixMilli(),
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
		Timestamp: time.Now().UnixMilli(),
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

// HandleServiceError 处理Service层错误
// 将ServiceError转换为HTTP响应
func HandleServiceError(c *gin.Context, err error) {
	// 检查是否为UnifiedError类型
	if unifiedErr, ok := err.(*apperrors.UnifiedError); ok {
		HandleUnifiedError(c, unifiedErr)
		return
	}

	// 默认错误处理
	Error(c, http.StatusInternalServerError, "操作失败", err.Error())
}

// HandleUnifiedError 处理UnifiedError错误
// 将UnifiedError转换为标准HTTP响应
func HandleUnifiedError(c *gin.Context, err *apperrors.UnifiedError) {
	statusCode := err.GetHTTPStatus()
	response := ErrorResponse{
		Code:      statusCode,
		Message:   err.Message,
		Error:     err.Code,
		Timestamp: time.Now().UnixMilli(),
		RequestID: err.RequestID,
	}

	// 在开发环境中添加调试信息
	// TODO: 从配置中读取环境变量
	if err.Details != "" {
		response.Debug = err.Details
	}

	c.JSON(statusCode, response)
}

// HandleError 处理任意类型的错误
// 自动识别错误类型并返回适当的HTTP响应
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 检查是否为UnifiedError
	if unifiedErr, ok := err.(*apperrors.UnifiedError); ok {
		HandleUnifiedError(c, unifiedErr)
		return
	}

	// 默认错误处理
	Error(c, http.StatusInternalServerError, "操作失败", err.Error())
}

// WrapServiceError 包装服务层错误为UnifiedError
// 用于在API层将服务层错误转换为统一的错误格式
func WrapServiceError(serviceName, operation string, err error) *apperrors.UnifiedError {
	if err == nil {
		return nil
	}

	// 如果已经是UnifiedError，直接返回
	if unifiedErr, ok := err.(*apperrors.UnifiedError); ok {
		return unifiedErr.WithService(serviceName, operation)
	}

	// 否则创建新的UnifiedError
	factory := apperrors.NewErrorFactory(serviceName)
	return factory.InternalError("INTERNAL_ERROR", "操作失败", err).
		WithService(serviceName, operation)
}

// =========================
// Gin便捷响应函数
// =========================

// Success 返回成功响应
// 统一的成功响应格式，用于所有API成功情况
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Code:      statusCode,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// SuccessData 返回成功响应（简化版，使用默认状态码和消息）
func SuccessData(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, "操作成功", data)
}

// Error 返回错误响应
// 统一的错误响应格式，用于所有API错误情况
func Error(c *gin.Context, statusCode int, message string, errorDetail string) {
	c.JSON(statusCode, ErrorResponse{
		Code:      statusCode,
		Message:   message,
		Error:     errorDetail,
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// ValidationError 返回参数验证错误响应
// 专门用于处理参数验证失败的情况
func ValidationError(c *gin.Context, err error) {
	errorDetail := ""
	if err != nil {
		errorDetail = err.Error()
	}
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:      http.StatusBadRequest,
		Message:   "参数验证失败",
		Error:     errorDetail,
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// Paginated 返回分页响应
// 统一的分页响应格式
func Paginated(c *gin.Context, data interface{}, total int64, page, pageSize int, message string) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
		Pagination: &Pagination{
			Total:       total,
			Page:        page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	})
}

// SuccessWithRequestID 返回带RequestID的成功响应
func SuccessWithRequestID(c *gin.Context, statusCode int, message string, data interface{}, requestID string) {
	c.JSON(statusCode, APIResponse{
		Code:      statusCode,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: requestID,
	})
}

// ErrorWithRequestID 返回带RequestID的错误响应
func ErrorWithRequestID(c *gin.Context, statusCode int, message string, errorDetail string, requestID string) {
	c.JSON(statusCode, ErrorResponse{
		Code:      statusCode,
		Message:   message,
		Error:     errorDetail,
		Timestamp: time.Now().UnixMilli(),
		RequestID: requestID,
	})
}

// Unauthorized 返回未授权错误
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:      http.StatusUnauthorized,
		Message:   message,
		Error:     "请先登录或提供有效的访问凭证",
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// Forbidden 返回禁止访问错误
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Code:      http.StatusForbidden,
		Message:   message,
		Error:     "您没有权限访问此资源",
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// NotFound 返回资源不存在错误
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:      http.StatusNotFound,
		Message:   message,
		Error:     "请求的资源不存在",
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// InternalError 返回内部服务器错误
func InternalError(c *gin.Context, message string, err error) {
	errorDetail := ""
	if err != nil {
		errorDetail = err.Error()
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:      http.StatusInternalServerError,
		Message:   message,
		Error:     errorDetail,
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}

// BadRequest 返回错误请求
func BadRequest(c *gin.Context, message string, errorDetail string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:      http.StatusBadRequest,
		Message:   message,
		Error:     errorDetail,
		Timestamp: time.Now().UnixMilli(),
		RequestID: builtin.GetRequestID(c),
	})
}
