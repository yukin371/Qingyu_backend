package ai

import (
	"net/http"
	"net/http/httptest"
	"testing"

	aiAPI "Qingyu_backend/api/v1/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitCreativeRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	InitCreativeRoutesWithAuth(v1, aiAPI.NewCreativeAPI(nil), func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/creative/outline", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}
