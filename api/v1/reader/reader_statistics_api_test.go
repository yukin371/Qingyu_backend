package reader

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestReaderStatisticsAPIRejectsMalformedUserIDContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/reader/statistics/heatmap", func(c *gin.Context) {
		c.Set("user_id", 123)
		(&ReaderStatisticsAPI{}).GetHeatmap(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/reader/statistics/heatmap", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestReaderStatisticsAPIRejectsEmptyUserIDContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/reader/statistics/heatmap", func(c *gin.Context) {
		c.Set("user_id", "")
		(&ReaderStatisticsAPI{}).GetHeatmap(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/reader/statistics/heatmap", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}
