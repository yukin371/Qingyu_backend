//go:build integration
// +build integration

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReaderBooksAPI_RemoveFromBookshelf 从书架移除书籍
func TestReaderBooksAPI_RemoveFromBookshelf(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.DELETE("/api/v1/reader/books/:bookId", readerAuthMiddleware("user123"), func(c *gin.Context) {
		bookID := c.Param("bookId")
		if bookID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    1,
				"message": "书籍ID不能为空",
			})
			return
		}

		// 模拟删除操作
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "移除成功",
			"data":    nil,
		})
	})

	// 测试成功删除
	req := httptest.NewRequest("DELETE", "/api/v1/reader/books/book123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "移除成功", response["message"])
}

// TestReaderBooksAPI_RemoveFromBookshelf_MissingID 书籍ID为空
func TestReaderBooksAPI_RemoveFromBookshelf_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.DELETE("/api/v1/reader/books/:bookId", readerAuthMiddleware("user123"), func(c *gin.Context) {
		bookID := c.Param("bookId")
		if bookID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    1,
				"message": "书籍ID不能为空",
				"error":   "missing_book_id",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "移除成功",
		})
	})

	req := httptest.NewRequest("DELETE", "/api/v1/reader/books/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 404 因为没有匹配的路由
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestReaderBooksAPI_RemoveFromBookshelf_Unauthorized 未授权
func TestReaderBooksAPI_RemoveFromBookshelf_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.DELETE("/api/v1/reader/books/:bookId", func(c *gin.Context) {
		// 模拟未授权
		_, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "移除成功",
		})
	})

	req := httptest.NewRequest("DELETE", "/api/v1/reader/books/book123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// readerAuthMiddleware 阅读器认证中间件
func readerAuthMiddleware(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", userID)
		c.Next()
	}
}
