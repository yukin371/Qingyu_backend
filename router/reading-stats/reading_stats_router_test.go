package readingstats

import (
	"net/http"
	"net/http/httptest"
	"testing"

	readingStatsService "Qingyu_backend/service/reader/stats"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterReadingStatsRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")
	statsService := readingStatsService.NewReadingStatsService(nil, nil, nil)

	RegisterReadingStatsRoutesWithAuth(v1, statsService, func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/my/stats", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}
