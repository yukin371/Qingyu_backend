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

// TestBookmarkValidationErrorFormat T4.1: 参数错误响应（400, 100001）
func TestBookmarkValidationErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/v1/reader/books/book123/bookmarks", func(c *gin.Context) {
		response.BadRequest(c, "参数验证失败", "position字段不能为空")
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	reqData := map[string]interface{}{
		"bookId":    "book123",
		"chapterId": "chapter123",
	}

	w := th.PerformRequest("POST", "/api/v1/reader/books/book123/bookmarks", reqData)

	validator.ValidateError(w, http.StatusBadRequest, int(errors.InvalidParams))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.InvalidParams), resp["code"])
	assert.NotEmpty(t, resp["message"])
	assert.NotEmpty(t, resp["request_id"])
	assert.Greater(t, resp["timestamp"].(float64), float64(1700000000000))
}

// TestBookmarkUnauthorizedErrorFormat T4.3: 未授权响应（401, 100601）
func TestBookmarkUnauthorizedErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/bookmarks", func(c *gin.Context) {
		response.Unauthorized(c, "未授权访问")
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("GET", "/api/v1/reader/bookmarks", nil)

	validator.ValidateError(w, http.StatusUnauthorized, int(errors.Unauthorized))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.Unauthorized), resp["code"])
	assert.Contains(t, resp["message"], "未授权")
}

// TestBookmarkNotFoundErrorFormat T4.4: 书签不存在响应（404, 100404）
func TestBookmarkNotFoundErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/v1/reader/bookmarks/:id", func(c *gin.Context) {
		response.NotFound(c, "书签不存在")
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	w := th.PerformRequest("GET", "/api/v1/reader/bookmarks/nonexistent", nil)

	validator.ValidateError(w, http.StatusNotFound, int(errors.NotFound))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.NotFound), resp["code"])
	assert.Contains(t, resp["message"], "不存在")
}

// TestBookmarkConflictErrorFormat T4.5: 书签已存在响应（409, 100409）
func TestBookmarkConflictErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/v1/reader/books/book123/bookmarks", func(c *gin.Context) {
		response.Conflict(c, "书签已存在", nil)
	})

	th := testutil.NewTestHelper(t, router)
	validator := testutil.NewResponseValidator(t)

	reqData := map[string]interface{}{
		"bookId":    "book123",
		"chapterId": "chapter123",
		"position":  100,
	}

	w := th.PerformRequest("POST", "/api/v1/reader/books/book123/bookmarks", reqData)

	validator.ValidateError(w, http.StatusConflict, int(errors.Conflict))

	resp := th.ParseResponse(w)
	assert.Equal(t, float64(errors.Conflict), resp["code"])
	assert.Contains(t, resp["message"], "已存在")
}
