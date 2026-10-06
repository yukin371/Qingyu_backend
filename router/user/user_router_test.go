package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterUserRoutesWithAuth(v1, nil, nil, nil, nil, nil, nil, nil, func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}
