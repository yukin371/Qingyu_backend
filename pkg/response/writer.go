package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应格式
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// PaginatedResponse 分页响应格式
type PaginatedResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	Pagination *Pagination `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// getRequestID 从context中获取request_id
func getRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	// 优先从context中获取（由middleware设置）
	if requestID := c.GetString("requestId"); requestID != "" {
		return requestID
	}
	// 其次从header中获取
	if c.Request != nil {
		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
			return requestID
		}
	}
	// 最后生成新的request_id（使用UUID或随机字符串）
	return generateRequestID()
}

// generateRequestID 生成新的request_id
func generateRequestID() string {
	return "req_" + time.Now().Format("20060102150405")
}

// Success 返回成功响应（200 OK）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:      0, // 成功响应code为0
		Message:   "操作成功",
		Data:      data,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
}

// Created 返回创建成功响应（201 Created）
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Code:      0, // 成功响应code为0
		Message:   "创建成功",
		Data:      data,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
}

// NoContent 返回无内容响应（204 No Content）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest 返回错误请求响应（400 Bad Request）
func BadRequest(c *gin.Context, message string, details interface{}) {
	response := APIResponse{
		Code:      CodeParamError, // 参数错误
		Message:   message,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	}
	if details != nil {
		response.Data = map[string]interface{}{
			"details": details,
		}
	}
	c.JSON(http.StatusBadRequest, response)
}

// Unauthorized 返回未授权响应（401 Unauthorized）
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Code:      CodeUnauthorized, // 未授权
		Message:   message,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
}

// Forbidden 返回禁止访问响应（403 Forbidden）
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Code:      CodeForbidden, // 禁止访问
		Message:   message,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
}

// NotFound 返回资源不存在响应（404 Not Found）
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Code:      CodeNotFound, // 资源不存在
		Message:   message,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
}

// Conflict 返回冲突响应（409 Conflict）
func Conflict(c *gin.Context, message string, details interface{}) {
	response := APIResponse{
		Code:      CodeConflict, // 资源冲突
		Message:   message,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	}
	if details != nil {
		response.Data = map[string]interface{}{
			"details": details,
		}
	}
	c.JSON(http.StatusConflict, response)
}

// InternalError 返回内部服务器错误响应（500 Internal Server Error）
func InternalError(c *gin.Context, err error) {
	message := "服务器内部错误"
	errorDetail := ""
	if err != nil {
		errorDetail = err.Error()
	}
	c.JSON(http.StatusInternalServerError, APIResponse{
		Code:      CodeInternalError, // 内部错误
		Message:   message,
		Data: map[string]interface{}{
			"error": errorDetail,
		},
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
}

// Paginated 返回分页响应（200 OK）
func Paginated(c *gin.Context, data interface{}, total int64, page, pageSize int, message string) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:      0, // 成功响应code为0
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
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

// SuccessWithMessage 返回带自定义消息的成功响应（200 OK）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:      0, // 成功响应code为0
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(), // 毫秒级时间戳
		RequestID: getRequestID(c),
	})
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
