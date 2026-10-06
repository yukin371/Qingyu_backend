package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRegisterAuthRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1/shared")

	RegisterAuthRoutesWithAuth(v1, nil, nil, zap.NewNop(), func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}

func TestRegisterStorageRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Use a fresh group so this test cannot depend on the auth routes above.
	storageEngine := gin.New()
	storageV1 := storageEngine.Group("/api/v1/shared")
	RegisterStorageRoutesWithAuth(storageV1, nil, nil, nil, func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/shared/storage/files", nil)
	rec := httptest.NewRecorder()
	storageEngine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}
