package ai

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWritingAPI_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewWritingApi(nil, nil)
	router := gin.New()
	router.POST("/continue", api.ContinueWriting)
	router.POST("/continue/stream", api.ContinueWritingStream)
	router.POST("/rewrite", api.RewriteText)
	router.POST("/rewrite/stream", api.RewriteTextStream)

	t.Run("continue invalid json should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/continue", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("continue missing required fields should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/continue", bytes.NewBufferString(`{"projectId":"p1"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rewrite invalid mode should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/rewrite", bytes.NewBufferString(
			`{"originalText":"test-content","rewriteMode":"invalid-mode"}`,
		))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("continue stream invalid body should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/continue/stream", bytes.NewBufferString(`{"projectId":"p1"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rewrite stream invalid body should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/rewrite/stream", bytes.NewBufferString(
			`{"originalText":"x","rewriteMode":"unknown"}`,
		))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

