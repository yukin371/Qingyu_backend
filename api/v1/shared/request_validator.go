package shared

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/validator"
)

// ValidateRequest 验证请求并返回友好错误
func ValidateRequest(c *gin.Context, req interface{}) bool {
	// 读取原始请求体用于调试
	bodyBytes, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 绑定请求
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数格式错误",
			Error:   err.Error(),
			Debug:   string(bodyBytes), // 添加调试信息
		})
		return false
	}

	// 验证请求
	validationErrors := validator.ValidateStructWithErrors(req)
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
	validationErrors := validator.ValidateStructWithErrors(params)
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

// ValidationErrorResponse 验证错误响应（字段级错误）
type ValidationErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"` // 字段名 -> 错误消息
}
