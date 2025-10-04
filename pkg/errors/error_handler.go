package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse HTTP错误响应结构
type ErrorResponse struct {
	Success   bool                   `json:"success"`
	Error     *ErrorInfo             `json:"error"`
	RequestID string                 `json:"request_id,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Details  string `json:"details,omitempty"`
	Category string `json:"category"`
	Level    string `json:"level"`
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	enableLogging    bool
	enableStackTrace bool
	converter        *ErrorConverter
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(enableLogging, enableStackTrace bool) *ErrorHandler {
	return &ErrorHandler{
		enableLogging:    enableLogging,
		enableStackTrace: enableStackTrace,
		converter:        NewErrorConverter(),
	}
}

// HandleError 处理错误并返回HTTP响应
func (h *ErrorHandler) HandleError(c *gin.Context, err error, service, operation string) {
	// 转换为统一错误
	unifiedErr := h.converter.ConvertGenericError(err, service, operation)
	
	// 添加请求上下文
	if requestID := c.GetString("request_id"); requestID != "" {
		unifiedErr.RequestID = requestID
	}
	if traceID := c.GetString("trace_id"); traceID != "" {
		unifiedErr.TraceID = traceID
	}
	if userID := c.GetString("user_id"); userID != "" {
		unifiedErr.UserID = userID
	}

	// 记录日志
	if h.enableLogging {
		h.logError(unifiedErr)
	}

	// 构建响应
	response := h.buildErrorResponse(unifiedErr)
	
	// 返回HTTP响应
	c.JSON(unifiedErr.GetHTTPStatus(), response)
}

// HandlePanic 处理panic
func (h *ErrorHandler) HandlePanic(c *gin.Context, panicValue interface{}, service, operation string) {
	// 转换panic为错误
	unifiedErr := h.converter.ConvertPanic(panicValue, service, operation)
	
	// 添加堆栈信息
	if h.enableStackTrace {
		unifiedErr.Stack = string(debug.Stack())
	}

	// 添加请求上下文
	if requestID := c.GetString("request_id"); requestID != "" {
		unifiedErr.RequestID = requestID
	}
	if traceID := c.GetString("trace_id"); traceID != "" {
		unifiedErr.TraceID = traceID
	}

	// 记录日志
	if h.enableLogging {
		h.logError(unifiedErr)
	}

	// 构建响应
	response := h.buildErrorResponse(unifiedErr)
	
	// 返回HTTP响应
	c.JSON(500, response)
}

// buildErrorResponse 构建错误响应
func (h *ErrorHandler) buildErrorResponse(err *UnifiedError) *ErrorResponse {
	response := &ErrorResponse{
		Success:   false,
		RequestID: err.RequestID,
		TraceID:   err.TraceID,
		Timestamp: err.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		Error: &ErrorInfo{
			Code:     err.Code,
			Message:  err.Message,
			Details:  err.Details,
			Category: string(err.Category),
			Level:    string(err.Level),
		},
	}

	// 添加元数据（排除敏感信息）
	if err.Metadata != nil && len(err.Metadata) > 0 {
		response.Metadata = make(map[string]interface{})
		for k, v := range err.Metadata {
			// 过滤敏感信息
			if !h.isSensitiveKey(k) {
				response.Metadata[k] = v
			}
		}
	}

	return response
}

// logError 记录错误日志
func (h *ErrorHandler) logError(err *UnifiedError) {
	logData := map[string]interface{}{
		"error_id":   err.ID,
		"code":       err.Code,
		"category":   err.Category,
		"level":      err.Level,
		"message":    err.Message,
		"service":    err.Service,
		"operation":  err.Operation,
		"user_id":    err.UserID,
		"request_id": err.RequestID,
		"trace_id":   err.TraceID,
		"timestamp":  err.Timestamp,
	}

	if err.Cause != nil {
		logData["cause"] = err.Cause.Error()
	}

	if err.Stack != "" && h.enableStackTrace {
		logData["stack"] = err.Stack
	}

	logJSON, _ := json.Marshal(logData)
	log.Printf("ERROR: %s", string(logJSON))
}

// isSensitiveKey 检查是否为敏感键
func (h *ErrorHandler) isSensitiveKey(key string) bool {
	sensitiveKeys := []string{
		"password", "token", "secret", "key", "auth",
		"credential", "private", "confidential",
	}
	
	for _, sensitive := range sensitiveKeys {
		if key == sensitive {
			return true
		}
	}
	return false
}

// ErrorMiddleware Gin错误处理中间件
func ErrorMiddleware(service string) gin.HandlerFunc {
	handler := NewErrorHandler(true, false)
	
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				handler.HandlePanic(c, r, service, c.Request.URL.Path)
				c.Abort()
			}
		}()
		
		c.Next()
		
		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handler.HandleError(c, err, service, c.Request.URL.Path)
			c.Abort()
		}
	}
}

// BusinessErrorMiddleware 业务错误处理中间件
func BusinessErrorMiddleware(service string) gin.HandlerFunc {
	handler := NewErrorHandler(true, false)
	
	return func(c *gin.Context) {
		c.Next()
		
		// 检查响应状态码
		if c.Writer.Status() >= 400 {
			// 如果已经设置了错误状态码但没有响应体，创建默认错误
			if c.Writer.Size() <= 0 {
				err := &UnifiedError{
					ID:         generateErrorID(),
					Code:       http.StatusText(c.Writer.Status()),
					Category:   CategorySystem,
					Level:      LevelError,
					Message:    "Request failed",
					Service:    service,
					HTTPStatus: c.Writer.Status(),
				}
				
				response := handler.buildErrorResponse(err)
				c.JSON(c.Writer.Status(), response)
			}
		}
	}
}

// 全局错误处理器实例
var DefaultHandler = NewErrorHandler(true, false)