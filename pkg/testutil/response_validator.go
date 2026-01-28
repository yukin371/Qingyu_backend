package testutil

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"Qingyu_backend/pkg/response"
)

// ResponseValidator 响应验证器
type ResponseValidator struct {
	t *testing.T
}

// NewResponseValidator 创建响应验证器
func NewResponseValidator(t *testing.T) *ResponseValidator {
	return &ResponseValidator{t: t}
}

// ValidateSuccess 验证成功响应
func (v *ResponseValidator) ValidateSuccess(w *httptest.ResponseRecorder, expectedHTTPStatus int) {
	// 验证HTTP状态码
	assert.Equal(v.t, expectedHTTPStatus, w.Code)

	// 解析响应
	var resp response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(v.t, err, "响应应该能正确解析")

	// 验证响应格式
	assert.Equal(v.t, 0, resp.Code, "成功响应code应该为0")
	assert.NotEmpty(v.t, resp.Message, "message不能为空")
	assert.NotEmpty(v.t, resp.RequestID, "request_id不能为空")
	assert.Greater(v.t, resp.Timestamp, int64(1700000000000), "timestamp应该是毫秒级")

	// 验证timestamp是13位（毫秒级）- 转换为秒后验证格式
	timestampStr := time.Unix(resp.Timestamp/1000, 0).Format("20060102030405")
	assert.Len(v.t, timestampStr, 14, "timestamp应该是毫秒级（13位数字）")
}

// ValidateError 验证错误响应
func (v *ResponseValidator) ValidateError(w *httptest.ResponseRecorder, expectedHTTPStatus, expectedErrorCode int) {
	// 验证HTTP状态码
	assert.Equal(v.t, expectedHTTPStatus, w.Code)

	// 解析响应
	var resp response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(v.t, err)

	// 验证错误码
	assert.Equal(v.t, expectedErrorCode, resp.Code, "错误码应该匹配")
	assert.NotEmpty(v.t, resp.Message, "错误消息不能为空")
	assert.NotEmpty(v.t, resp.RequestID, "request_id不能为空")
	assert.Greater(v.t, resp.Timestamp, int64(1700000000000))

	// 验证timestamp是13位（毫秒级）- 转换为秒后验证格式
	timestampStr := time.Unix(resp.Timestamp/1000, 0).Format("20060102030405")
	assert.Len(v.t, timestampStr, 14, "timestamp应该是毫秒级（13位数字）")
}

// ValidatePagination 验证分页响应
func (v *ResponseValidator) ValidatePagination(w *httptest.ResponseRecorder) {
	var resp response.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(v.t, err)

	assert.NotNil(v.t, resp.Pagination, "分页信息不能为空")
	assert.Greater(v.t, resp.Pagination.Total, int64(0), "总数应该大于0")
	assert.Greater(v.t, resp.Pagination.Page, 0, "页码应该大于0")
	assert.Greater(v.t, resp.Pagination.PageSize, 0, "页大小应该大于0")
}
