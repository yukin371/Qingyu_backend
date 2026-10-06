package community

import (
	"net/http"
	"net/http/httptest"
	"testing"

	communityAPI "Qingyu_backend/api/v1/community"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCommunityRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterCommunityRoutesWithAuth(v1, (*communityAPI.PostAPI)(nil), func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/community/posts", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
	paths := make(map[string]bool)
	for _, route := range engine.Routes() {
		paths[route.Path] = true
	}
	assert.True(t, paths["/api/v1/community/posts"])
}
