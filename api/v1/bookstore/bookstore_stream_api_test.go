package bookstore

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStreamSearchBooksReturnsUnavailableWhenServiceIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := &BookstoreAPI{}
	r.GET("/books/stream", api.StreamSearchBooks)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/books/stream?keyword=qingyu", nil)
	r.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), "bookstore stream service is not configured")
}
