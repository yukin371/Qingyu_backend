package system

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitSystemRoutes_RegistersUnderAPIV1Prefix(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	v1 := engine.Group("/api/v1")

	InitSystemRoutes(v1)

	paths := map[string]bool{}
	for _, route := range engine.Routes() {
		paths[route.Path] = true
	}

	assert.True(t, paths["/api/v1/system/health"])
	assert.True(t, paths["/api/v1/system/health/:service"])
	assert.True(t, paths["/api/v1/system/metrics"])
	assert.True(t, paths["/api/v1/system/metrics/:service"])
	assert.True(t, paths["/api/v1/errors/report"])
	assert.False(t, paths["/system/health"])
	assert.False(t, paths["/system/metrics"])
}
