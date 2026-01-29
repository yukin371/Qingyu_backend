package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccess(t *testing.T) {
	c, w := setupTestContext()

	data := map[string]string{"key": "value"}
	Success(c, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"code":0`)
	assert.Contains(t, w.Body.String(), `"message":"操作成功"`)
	assert.Contains(t, w.Body.String(), `"data":`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
}

func TestCreated(t *testing.T) {
	c, w := setupTestContext()

	data := map[string]string{"id": "123"}
	Created(c, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"code":0`)
	assert.Contains(t, w.Body.String(), `"message":"创建成功"`)
	assert.Contains(t, w.Body.String(), `"data":`)
}

func TestNoContent(t *testing.T) {
	c, _ := setupTestContext()

	// 初始化request（gin测试框架需要）
	c.Request = httptest.NewRequest("GET", "/test", nil)

	NoContent(c)

	// 注意：gin测试框架中，c.Status()不会自动写入response
	// 这里我们检查context是否设置了正确的状态码
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
}

func TestBadRequest(t *testing.T) {
	c, w := setupTestContext()

	BadRequest(c, "参数错误", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"code":1001`)
	assert.Contains(t, w.Body.String(), `"message":"参数错误"`)
}

func TestBadRequestWithDetails(t *testing.T) {
	c, w := setupTestContext()

	details := map[string]string{"field": "error"}
	BadRequest(c, "参数错误", details)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"details":`)
}

func TestUnauthorized(t *testing.T) {
	c, w := setupTestContext()

	Unauthorized(c, "未授权")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"code":1002`)
	assert.Contains(t, w.Body.String(), `"message":"未授权"`)
}

func TestForbidden(t *testing.T) {
	c, w := setupTestContext()

	Forbidden(c, "禁止访问")

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), `"code":1003`)
	assert.Contains(t, w.Body.String(), `"message":"禁止访问"`)
}

func TestNotFound(t *testing.T) {
	c, w := setupTestContext()

	NotFound(c, "资源不存在")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"code":1004`)
	assert.Contains(t, w.Body.String(), `"message":"资源不存在"`)
}

func TestConflict(t *testing.T) {
	c, w := setupTestContext()

	Conflict(c, "资源冲突", nil)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), `"code":1006`)
	assert.Contains(t, w.Body.String(), `"message":"资源冲突"`)
}

func TestConflictWithDetails(t *testing.T) {
	c, w := setupTestContext()

	details := map[string]string{"field": "conflict"}
	Conflict(c, "资源冲突", details)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), `"details":`)
}

func TestInternalError(t *testing.T) {
	c, w := setupTestContext()

	err := errors.New("database error")
	InternalError(c, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"code":5000`)
	assert.Contains(t, w.Body.String(), `"message":"服务器内部错误"`)
	assert.Contains(t, w.Body.String(), `"error":"database error"`)
}

func TestInternalErrorNil(t *testing.T) {
	c, w := setupTestContext()

	InternalError(c, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"code":5000`)
	assert.Contains(t, w.Body.String(), `"message":"服务器内部错误"`)
}

func TestPaginated(t *testing.T) {
	c, w := setupTestContext()

	data := []string{"item1", "item2", "item3"}
	Paginated(c, data, 100, 2, 10, "获取成功")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"code":0`)
	assert.Contains(t, w.Body.String(), `"message":"获取成功"`)
	assert.Contains(t, w.Body.String(), `"pagination":`)
	assert.Contains(t, w.Body.String(), `"total":100`)
	assert.Contains(t, w.Body.String(), `"page":2`)
	assert.Contains(t, w.Body.String(), `"page_size":10`)
	assert.Contains(t, w.Body.String(), `"total_pages":10`)
	assert.Contains(t, w.Body.String(), `"has_next":true`)
	assert.Contains(t, w.Body.String(), `"has_previous":true`)
}

func TestPaginatedFirstPage(t *testing.T) {
	c, w := setupTestContext()

	data := []string{"item1"}
	Paginated(c, data, 100, 1, 10, "获取成功")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"has_next":true`)
	assert.Contains(t, w.Body.String(), `"has_previous":false`)
}

func TestPaginatedLastPage(t *testing.T) {
	c, w := setupTestContext()

	data := []string{"item1"}
	Paginated(c, data, 100, 10, 10, "获取成功")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"has_next":false`)
	assert.Contains(t, w.Body.String(), `"has_previous":true`)
}

func TestSuccessWithMessage(t *testing.T) {
	c, w := setupTestContext()

	data := map[string]string{"key": "value"}
	SuccessWithMessage(c, "自定义消息", data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"code":0`)
	assert.Contains(t, w.Body.String(), `"message":"自定义消息"`)
	assert.Contains(t, w.Body.String(), `"data":`)
}

func TestNewPagination(t *testing.T) {
	pagination := NewPagination(100, 2, 10)

	assert.Equal(t, int64(100), pagination.Total)
	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 10, pagination.PageSize)
	assert.Equal(t, 10, pagination.TotalPages)
	assert.True(t, pagination.HasNext)
	assert.True(t, pagination.HasPrevious)
}

func TestNewPaginationFirstPage(t *testing.T) {
	pagination := NewPagination(100, 1, 10)

	assert.Equal(t, int64(100), pagination.Total)
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.PageSize)
	assert.Equal(t, 10, pagination.TotalPages)
	assert.True(t, pagination.HasNext)
	assert.False(t, pagination.HasPrevious)
}

func TestNewPaginationLastPage(t *testing.T) {
	pagination := NewPagination(100, 10, 10)

	assert.Equal(t, int64(100), pagination.Total)
	assert.Equal(t, 10, pagination.Page)
	assert.Equal(t, 10, pagination.PageSize)
	assert.Equal(t, 10, pagination.TotalPages)
	assert.False(t, pagination.HasNext)
	assert.True(t, pagination.HasPrevious)
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试从header获取request_id
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("X-Request-ID", "test-request-id")

	requestID := getRequestID(c)
	assert.Equal(t, "test-request-id", requestID)
}

func TestGetRequestIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试从context获取request_id
	c.Set("requestId", "context-request-id")

	requestID := getRequestID(c)
	assert.Equal(t, "context-request-id", requestID)
}
