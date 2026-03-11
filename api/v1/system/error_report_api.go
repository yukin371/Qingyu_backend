package system

import (
	"log"
	"strings"
	"time"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type clientErrorReport struct {
	ErrorCode    int                    `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage" binding:"required"`
	ErrorType    string                 `json:"errorType"`
	URL          string                 `json:"url"`
	UserAgent    string                 `json:"userAgent"`
	UserID       string                 `json:"userId"`
	Timestamp    string                 `json:"timestamp"`
	Details      map[string]interface{} `json:"details"`
}

type clientErrorReportBatchRequest struct {
	Errors []clientErrorReport `json:"errors" binding:"required,min=1,max=50"`
}

// ReportClientErrors 接收前端批量错误上报
// @Summary 接收前端错误上报
// @Description 接收浏览器端批量错误上报，当前仅记录日志并返回成功
// @Tags 系统监控
// @Accept json
// @Produce json
// @Param request body clientErrorReportBatchRequest true "错误上报批次"
// @Success 200 {object} response.APIResponse "接收成功"
// @Failure 400 {object} response.APIResponse "参数错误"
// @Router /api/v1/errors/report [post]
func (api *HealthAPI) ReportClientErrors(c *gin.Context) {
	var req clientErrorReportBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "错误上报参数无效", err.Error())
		return
	}

	for _, item := range req.Errors {
		log.Printf(
			"[client-error] code=%d type=%s url=%s user=%s message=%s timestamp=%s details=%v ua=%s",
			item.ErrorCode,
			item.ErrorType,
			item.URL,
			item.UserID,
			item.ErrorMessage,
			normalizeClientErrorTimestamp(item.Timestamp),
			item.Details,
			truncateClientUserAgent(item.UserAgent),
		)
	}

	response.SuccessWithMessage(c, "错误上报已接收", gin.H{
		"accepted": len(req.Errors),
	})
}

func normalizeClientErrorTimestamp(raw string) string {
	if raw == "" {
		return time.Now().Format(time.RFC3339)
	}
	if _, err := time.Parse(time.RFC3339, raw); err == nil {
		return raw
	}
	return time.Now().Format(time.RFC3339)
}

func truncateClientUserAgent(userAgent string) string {
	const maxLen = 160
	if len(userAgent) <= maxLen {
		return userAgent
	}
	return strings.TrimSpace(userAgent[:maxLen]) + "..."
}
