package reader

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitReaderRouterWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	InitReaderRouterWithAuth(
		v1,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		func(c *gin.Context) {
			c.AbortWithStatus(http.StatusTeapot)
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/books", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
	paths := make(map[string]bool)
	for _, route := range engine.Routes() {
		paths[route.Path] = true
	}
	assert.True(t, paths["/api/v1/reader/books"])
}
