package integration

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/pkg/testutil"
	"Qingyu_backend/pkg/response"
)

// TestCreateAnnotationResponseFormat T3.1: 创建成功响应格式验证
func TestCreateAnnotationResponseFormat(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 创建一个简单的Handler来模拟CreateAnnotation
	router.POST("/api/v1/reader/annotations", func(c *gin.Context) {
		response.Created(c, gin.H{"id": "annotation123"})
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	// 准备请求数据
	reqData := map[string]string{
		"bookId":    "book123",
		"chapterId": "chapter123",
		"type":      "bookmark",
	}

	w := th.PerformRequest("POST", "/api/v1/reader/annotations", reqData)

	// 验证HTTP状态码
	assert.Equal(t, http.StatusCreated, w.Code)

	// 验证响应格式
	validator.ValidateSuccess(w, http.StatusCreated)

	resp := th.ParseResponse(w)
	assert.Equal(t, 0.0, resp["code"])
	assert.NotEmpty(t, resp["message"])
	assert.NotEmpty(t, resp["data"])
	assert.NotEmpty(t, resp["request_id"])
	assert.Greater(t, resp["timestamp"].(float64), float64(1700000000000))
}

// TestGetAnnotationsByChapterResponseFormat T3.2: 查询成功响应格式验证
func TestGetAnnotationsByChapterResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 模拟GetAnnotationsByChapter
	router.GET("/api/v1/reader/annotations/chapter", func(c *gin.Context) {
		response.Success(c, []interface{}{})
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("GET", "/api/v1/reader/annotations/chapter?bookId=book123&chapterId=chapter123", nil)

	// 验证响应格式
	validator.ValidateSuccess(w, http.StatusOK)

	resp := th.ParseResponse(w)
	assert.Equal(t, 0.0, resp["code"])
	assert.NotNil(t, resp["data"])
}

// TestDeleteAnnotationNoContent T3.4: 删除成功响应验证（200状态码）
func TestDeleteAnnotationResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 模拟DeleteAnnotation
	router.DELETE("/api/v1/reader/annotations/:id", func(c *gin.Context) {
		response.Success(c, nil)
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("DELETE", "/api/v1/reader/annotations/annotation123", nil)

	// 验证响应格式
	validator.ValidateSuccess(w, http.StatusOK)

	resp := th.ParseResponse(w)
	assert.Equal(t, "操作成功", resp["message"])
}

// TestRequestIDExists T3.5: request_id存在性验证
func TestRequestIDExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/annotations/test", func(c *gin.Context) {
		response.Success(c, nil)
	})

	th := testutil.NewTestHelper(t, router)
	w := th.PerformRequest("GET", "/api/v1/reader/annotations/test", nil)

	resp := th.ParseResponse(w)
	assert.NotEmpty(t, resp["request_id"], "request_id不能为空")
}

// TestTimestampPrecision T3.6: timestamp精度验证（毫秒级）
func TestTimestampPrecision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/annotations/test", func(c *gin.Context) {
		response.Success(c, nil)
	})

	th := testutil.NewTestHelper(t, router)
	w := th.PerformRequest("GET", "/api/v1/reader/annotations/test", nil)

	resp := th.ParseResponse(w)
	timestamp := int64(resp["timestamp"].(float64))

	// 验证是毫秒级时间戳（13位数字）
	assert.GreaterOrEqual(t, timestamp, int64(1000000000000))
	assert.Less(t, timestamp, int64(9999999999999))
}
