package integration

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/pkg/testutil"
	"Qingyu_backend/pkg/response"
)

// TestCreateBookmarkResponseFormat T3.1: 创建书签成功响应格式验证
func TestCreateBookmarkResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/v1/reader/books/book123/bookmarks", func(c *gin.Context) {
		response.Created(c, gin.H{"id": "bookmark123"})
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	reqData := map[string]interface{}{
		"bookId":    "book123",
		"chapterId": "chapter123",
		"position":  100,
	}

	w := th.PerformRequest("POST", "/api/v1/reader/books/book123/bookmarks", reqData)

	assert.Equal(t, http.StatusCreated, w.Code)
	validator.ValidateSuccess(w, http.StatusCreated)

	resp := th.ParseResponse(w)
	assert.Equal(t, 0.0, resp["code"])
	assert.NotEmpty(t, resp["message"])
	assert.NotEmpty(t, resp["data"])
	assert.NotEmpty(t, resp["request_id"])
	assert.Greater(t, resp["timestamp"].(float64), float64(1700000000000))
}

// TestGetBookmarksResponseFormat T3.2: 获取书签列表成功响应格式验证
func TestGetBookmarksResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/bookmarks", func(c *gin.Context) {
		response.Success(c, []interface{}{})
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("GET", "/api/v1/reader/bookmarks?page=1&size=20", nil)

	validator.ValidateSuccess(w, http.StatusOK)

	resp := th.ParseResponse(w)
	assert.Equal(t, 0.0, resp["code"])
	assert.NotNil(t, resp["data"])
}

// TestDeleteBookmarkResponseFormat T3.4: 删除书签成功响应验证
func TestDeleteBookmarkResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/api/v1/reader/bookmarks/:id", func(c *gin.Context) {
		response.Success(c, nil)
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("DELETE", "/api/v1/reader/bookmarks/bookmark123", nil)

	validator.ValidateSuccess(w, http.StatusOK)

	resp := th.ParseResponse(w)
	assert.Equal(t, "操作成功", resp["message"])
}

// TestBookmarkRequestIDExists T3.5: request_id存在性验证
func TestBookmarkRequestIDExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/bookmarks/test", func(c *gin.Context) {
		response.Success(c, nil)
	})

	th := testutil.NewTestHelper(t, router)
	w := th.PerformRequest("GET", "/api/v1/reader/bookmarks/test", nil)

	resp := th.ParseResponse(w)
	assert.NotEmpty(t, resp["request_id"], "request_id不能为空")
}

// TestBookmarkTimestampPrecision T3.6: timestamp精度验证（毫秒级）
func TestBookmarkTimestampPrecision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/bookmarks/test", func(c *gin.Context) {
		response.Success(c, nil)
	})

	th := testutil.NewTestHelper(t, router)
	w := th.PerformRequest("GET", "/api/v1/reader/bookmarks/test", nil)

	resp := th.ParseResponse(w)
	timestamp := int64(resp["timestamp"].(float64))

	// 验证是毫秒级时间戳（13位数字）
	assert.GreaterOrEqual(t, timestamp, int64(1000000000000))
	assert.Less(t, timestamp, int64(9999999999999))
}
