package social

import (
	"net/http"
	"net/http/httptest"
	"testing"

	socialAPI "Qingyu_backend/api/v1/social"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSocialRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterSocialRoutesWithAuth(
		v1,
		nil,
		socialAPI.NewCommentAPI(nil),
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/social/comments?book_id=book-1", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
	paths := make(map[string]bool)
	for _, route := range engine.Routes() {
		paths[route.Path] = true
	}
	assert.True(t, paths["/api/v1/social/comments"])
}
