package shared

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	appValidator "Qingyu_backend/pkg/validator"
)

// GetValidator 获取全局验证器实例
func GetValidator() *validator.Validate {
	return appValidator.GetValidator()
}

// ValidateRequest 验证请求并返回友好错误
// 注意：此函数不会绑定请求体，假设已经通过BindJSON等函数完成绑定
func ValidateRequest(c *gin.Context, req interface{}) bool {
	// 向后兼容：自动尝试绑定 JSON 请求体。
	// 旧代码大量直接调用 ValidateRequest 而未显式绑定，导致字段始终为空。
	method := c.Request.Method
	contentType := strings.ToLower(c.GetHeader("Content-Type"))
	shouldBindJSON := (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) &&
		strings.Contains(contentType, "application/json")
	if shouldBindJSON {
		bodyBytes, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    400,
				Message: "请求体读取失败",
				Error:   err.Error(),
			})
			return false
		}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, req); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    400,
					Message: "请求体格式错误",
					Error:   err.Error(),
				})
				return false
			}
			// 还原 Body，避免后续中间件/处理器读取为空。
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	// 验证请求
	validationErrors := appValidator.ValidateStructWithErrors(req)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Code:    400,
			Message: "请求参数验证失败",
			Errors:  validationErrors.GetFieldErrors(),
		})
		return false
	}

	return true
}

// ValidateQueryParams 验证查询参数
func ValidateQueryParams(c *gin.Context, params interface{}) bool {
	// 绑定查询参数
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "查询参数格式错误",
			Error:   err.Error(),
		})
		return false
	}

	// 验证参数
	validationErrors := appValidator.ValidateStructWithErrors(params)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Code:    400,
			Message: "查询参数验证失败",
			Errors:  validationErrors.GetFieldErrors(),
		})
		return false
	}

	return true
}

// HandleValidationError 处理验证错误
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Code:    400,
		Message: "请求参数验证失败",
		Errors: map[string]string{
			"validation": err.Error(),
		},
	})
}

// ValidationErrorResponse 验证错误响应（字段级错误）
type ValidationErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"` // 字段名 -> 错误消息
}
