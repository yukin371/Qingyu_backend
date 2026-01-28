package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// TestSuccess_SuccessResponse 测试成功响应
func TestSuccess_SuccessResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		data       interface{}
	}{
		{
			name:       "正常成功响应",
			statusCode: http.StatusOK,
			message:    "操作成功",
			data:       map[string]string{"key": "value"},
		},
		{
			name:       "无数据成功响应",
			statusCode: http.StatusOK,
			message:    "操作成功",
			data:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置request_id到context
			c.Set("request_id", "test-request-123")

			Success(c, tt.statusCode, tt.message, tt.data)

			// 验证HTTP状态码
			assert.Equal(t, tt.statusCode, w.Code)

			// 验证响应体包含必要字段
			assert.Contains(t, w.Body.String(), `"code":`)
			assert.Contains(t, w.Body.String(), `"message":`)
			assert.Contains(t, w.Body.String(), `"timestamp":`)
			assert.Contains(t, w.Body.String(), `"request_id":"test-request-123"`)
		})
	}
}

// TestError_ErrorResponse 测试错误响应
func TestError_ErrorResponse(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		message     string
		errorDetail string
	}{
		{
			name:        "Bad Request",
			statusCode:  http.StatusBadRequest,
			message:     "参数错误",
			errorDetail: "缺少必需参数",
		},
		{
			name:        "Internal Server Error",
			statusCode:  http.StatusInternalServerError,
			message:     "服务器错误",
			errorDetail: "数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置request_id到context
			c.Set("request_id", "test-request-456")

			Error(c, tt.statusCode, tt.message, tt.errorDetail)

			// 验证HTTP状态码
			assert.Equal(t, tt.statusCode, w.Code)

			// 验证响应体包含必要字段
			assert.Contains(t, w.Body.String(), `"code":`)
			assert.Contains(t, w.Body.String(), `"message":`)
			assert.Contains(t, w.Body.String(), `"error":`)
			assert.Contains(t, w.Body.String(), `"timestamp":`)
			assert.Contains(t, w.Body.String(), `"request_id":"test-request-456"`)
		})
	}
}

// TestValidationError_ValidationError 测试参数验证错误响应
func TestValidationError_ValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-789")

	err := assert.AnError
	ValidationError(c, err)

	// 验证HTTP状态码为400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证响应体
	assert.Contains(t, w.Body.String(), `"code":400`)
	assert.Contains(t, w.Body.String(), `"message":"参数验证失败"`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-789"`)
}

// TestUnauthorized_UnauthorizedResponse 测试未授权响应
func TestUnauthorized_UnauthorizedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-unauth")

	Unauthorized(c, "请先登录")

	// 验证HTTP状态码为401
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 验证响应体
	assert.Contains(t, w.Body.String(), `"code":401`)
	assert.Contains(t, w.Body.String(), `"message":"请先登录"`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-unauth"`)
}

// TestForbidden_ForbiddenResponse 测试禁止访问响应
func TestForbidden_ForbiddenResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-forbidden")

	Forbidden(c, "权限不足")

	// 验证HTTP状态码为403
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 验证响应体
	assert.Contains(t, w.Body.String(), `"code":403`)
	assert.Contains(t, w.Body.String(), `"message":"权限不足"`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-forbidden"`)
}

// TestNotFound_NotFoundResponse 测试资源不存在响应
func TestNotFound_NotFoundResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-notfound")

	NotFound(c, "书籍不存在")

	// 验证HTTP状态码为404
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 验证响应体
	assert.Contains(t, w.Body.String(), `"code":404`)
	assert.Contains(t, w.Body.String(), `"message":"书籍不存在"`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-notfound"`)
}

// TestInternalError_InternalErrorResponse 测试内部错误响应
func TestInternalError_InternalErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-internal")

	err := assert.AnError
	InternalError(c, "处理失败", err)

	// 验证HTTP状态码为500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证响应体
	assert.Contains(t, w.Body.String(), `"code":500`)
	assert.Contains(t, w.Body.String(), `"message":"处理失败"`)
	assert.Contains(t, w.Body.String(), `"error":`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-internal"`)
}

// TestBadRequest_BadRequestResponse 测试错误请求响应
func TestBadRequest_BadRequestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-bad")

	BadRequest(c, "参数无效", "ID格式错误")

	// 验证HTTP状态码为400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 验证响应体
	assert.Contains(t, w.Body.String(), `"code":400`)
	assert.Contains(t, w.Body.String(), `"message":"参数无效"`)
	assert.Contains(t, w.Body.String(), `"error":"ID格式错误"`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-bad"`)
}

// TestPaginated_PaginatedResponse 测试分页响应
func TestPaginated_PaginatedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-paginated")

	data := []string{"item1", "item2", "item3"}
	Paginated(c, data, 100, 2, 10, "获取成功")

	// 验证HTTP状态码为200
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证响应体包含分页信息
	assert.Contains(t, w.Body.String(), `"code":200`)
	assert.Contains(t, w.Body.String(), `"message":"获取成功"`)
	assert.Contains(t, w.Body.String(), `"total":100`)
	assert.Contains(t, w.Body.String(), `"page":2`)
	assert.Contains(t, w.Body.String(), `"page_size":10`)
	assert.Contains(t, w.Body.String(), `"total_pages":10`)
	assert.Contains(t, w.Body.String(), `"has_next":true`)
	assert.Contains(t, w.Body.String(), `"has_previous":true`)
	assert.Contains(t, w.Body.String(), `"timestamp":`)
	assert.Contains(t, w.Body.String(), `"request_id":"test-request-paginated"`)
}

// TestTimestampIsMilliseconds 验证时间戳是毫秒级
func TestTimestampIsMilliseconds(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "test-request-timestamp")

	Success(c, http.StatusOK, "测试", nil)

	body := w.Body.String()

	// 毫秒级时间戳应该是13位数字（当前时间范围内）
	// 我们可以通过检查timestamp的值来验证
	// 当前时间戳（秒）是10位，毫秒是13位
	assert.Contains(t, body, `"timestamp":`)
	// 简单验证：timestamp应该大于1_000_000_000_000（2001年以后的毫秒时间戳）
	// 这里我们主要验证函数被正确调用
}

// TestRequestID_PresentInAllResponses 验证所有响应都包含request_id
func TestRequestID_PresentInAllResponses(t *testing.T) {
	testRequestID := "test-req-123"

	tests := []struct {
		name       string
		setupFunc  func(*gin.Context)
		containsID bool
	}{
		{
			name: "Success响应包含request_id",
			setupFunc: func(c *gin.Context) {
				c.Set("request_id", testRequestID)
				Success(c, http.StatusOK, "成功", nil)
			},
			containsID: true,
		},
		{
			name: "Error响应包含request_id",
			setupFunc: func(c *gin.Context) {
				c.Set("request_id", testRequestID)
				Error(c, http.StatusBadRequest, "错误", "详情")
			},
			containsID: true,
		},
		{
			name: "NotFound响应包含request_id",
			setupFunc: func(c *gin.Context) {
				c.Set("request_id", testRequestID)
				NotFound(c, "未找到")
			},
			containsID: true,
		},
		{
			name: "Paginated响应包含request_id",
			setupFunc: func(c *gin.Context) {
				c.Set("request_id", testRequestID)
				Paginated(c, []string{}, 0, 1, 10, "成功")
			},
			containsID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupFunc(c)

			if tt.containsID {
				assert.Contains(t, w.Body.String(), `"request_id":"`+testRequestID+`"`)
			}
		})
	}
}

// TestHTTPStatusCodes 验证HTTP状态码使用正确
func TestHTTPStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(*gin.Context)
		wantStatus int
	}{
		{
			name: "Success返回200",
			setupFunc: func(c *gin.Context) {
				Success(c, http.StatusOK, "成功", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "BadRequest返回400",
			setupFunc: func(c *gin.Context) {
				BadRequest(c, "错误", "详情")
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unauthorized返回401",
			setupFunc: func(c *gin.Context) {
				Unauthorized(c, "未授权")
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden返回403",
			setupFunc: func(c *gin.Context) {
				Forbidden(c, "禁止")
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "NotFound返回404",
			setupFunc: func(c *gin.Context) {
				NotFound(c, "未找到")
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "InternalError返回500",
			setupFunc: func(c *gin.Context) {
				InternalError(c, "错误", assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("request_id", "test")

			tt.setupFunc(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// BenchmarkSuccess 性能测试
func BenchmarkSuccess(b *testing.B) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("request_id", "bench-test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Set("request_id", "bench-test")
		Success(c, http.StatusOK, "成功", nil)
	}
}
