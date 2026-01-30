package integration

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/pkg/testutil"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/pkg/errors"
)

// TestValidationErrorFormat T4.1: 参数错误响应（400, 100001）
func TestValidationErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/v1/reader/annotations", func(c *gin.Context) {
		response.BadRequest(c, "参数验证失败", "type字段不能为空")
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	// 发送缺少type参数的请求
	reqData := map[string]string{
		"bookId":    "book123",
		"chapterId": "chapter123",
	}

	w := th.PerformRequest("POST", "/api/v1/reader/annotations", reqData)

	// 验证错误响应格式
	validator.ValidateError(w, http.StatusBadRequest, int(errors.InvalidParams))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.InvalidParams), resp["code"])
	assert.NotEmpty(t, resp["message"])
	assert.NotEmpty(t, resp["request_id"])
	assert.Greater(t, resp["timestamp"].(float64), float64(1700000000000))
	// 注意：新的response.BadRequest不返回details字段
}

// TestUnauthorizedErrorFormat T4.3: 未授权响应（401, 100601）
func TestUnauthorizedErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/annotations/chapter", func(c *gin.Context) {
		response.Unauthorized(c, "未授权访问")
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("GET", "/api/v1/reader/annotations/chapter", nil)

	// 验证错误响应格式
	validator.ValidateError(w, http.StatusUnauthorized, int(errors.Unauthorized))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.Unauthorized), resp["code"])
	assert.Contains(t, resp["message"], "未授权")
}

// TestBadRequestErrorFormat T4.1: 参数错误-缺少必填字段
func TestBadRequestErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/annotations/chapter", func(c *gin.Context) {
		response.BadRequest(c, "参数错误", "书籍ID和章节ID不能为空")
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	// 不提供查询参数
	w := th.PerformRequest("GET", "/api/v1/reader/annotations/chapter", nil)

	// 验证错误响应格式
	validator.ValidateError(w, http.StatusBadRequest, int(errors.InvalidParams))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.InvalidParams), resp["code"])
	// 注意：新的response.BadRequest使用自定义消息，不包含原message
}

// TestErrorCodeMapping T4.6: 错误码映射验证
func TestErrorCodeMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/v1/reader/annotations", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      100001, // InvalidParams
			"message":   "请求参数无效",
			"timestamp": 1706684800000,
		})
	})

	th := testutil.NewTestHelper(t, router)

	w := th.PerformRequest("POST", "/api/v1/reader/annotations", map[string]string{})

	resp := th.ParseResponse(w)
	code := int(resp["code"].(float64))

	// 验证是6位错误码
	// TODO 需要改为4位错误码
	assert.GreaterOrEqual(t, code, 100000, "错误码应该是6位数字")
	assert.Less(t, code, 1000000, "错误码应该是6位数字")
	assert.Equal(t, 100001, code)
}
