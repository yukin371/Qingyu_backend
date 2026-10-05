package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/pkg/response"
)

func setupStatsTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	statsHandler := NewStatsHandler(nil, nil)

	authorized := func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	}

	router.GET("/stats/my/activity", authorized, statsHandler.GetMyActivityStats)
	router.GET("/stats/my/revenue", authorized, statsHandler.GetMyRevenueStats)

	return router
}

func TestStatsHandler_GetMyActivityStats_NotImplemented(t *testing.T) {
	router := setupStatsTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/stats/my/activity?days=30", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotImplemented, w.Code)

	var resp response.APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, response.CodeInternalError, resp.Code)
	assert.Equal(t, "活跃度统计功能开发中", resp.Message)
}

func TestStatsHandler_GetMyRevenueStats_NotImplemented(t *testing.T) {
	router := setupStatsTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/stats/my/revenue?start_date=2026-01-01&end_date=2026-01-31", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotImplemented, w.Code)

	var resp response.APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, response.CodeInternalError, resp.Code)
	assert.Equal(t, "收益统计功能开发中", resp.Message)
}

func TestStatsHandler_RejectsMalformedUserIDContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	statsHandler := NewStatsHandler(nil, nil)
	router.GET("/stats/my", func(c *gin.Context) {
		c.Set("user_id", 123)
		c.Next()
	}, statsHandler.GetMyStats)
	router.GET("/stats/my/content", func(c *gin.Context) {
		c.Set("user_id", 123)
		c.Next()
	}, statsHandler.GetMyContentStats)

	for _, target := range []string{"/stats/my", "/stats/my/content"} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, target, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, target)
	}
}
