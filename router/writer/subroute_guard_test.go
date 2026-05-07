package writer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Qingyu_backend/config"
	documentservice "Qingyu_backend/service/writer/document"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newWriterRouteEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func performWriterRequest(engine *gin.Engine, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func TestInitTimelineRoutes_RegisteredRouteReturnsBadRequestInsteadOf404(t *testing.T) {
	engine := newWriterRouteEngine()
	v1Writer := engine.Group("/api/v1/writer")
	InitTimelineRoutes(v1Writer, nil)

	w := performWriterRequest(engine, http.MethodGet, "/api/v1/writer/timelines/timeline-1", "")

	assert.NotEqual(t, http.StatusNotFound, w.Code)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestInitOutlineRoutes_RegisteredRouteReturnsBadRequestInsteadOf404(t *testing.T) {
	engine := newWriterRouteEngine()
	v1Writer := engine.Group("/api/v1/writer")
	InitOutlineRoutes(v1Writer, nil)

	w := performWriterRequest(engine, http.MethodGet, "/api/v1/writer/outlines/outline-1", "")

	assert.NotEqual(t, http.StatusNotFound, w.Code)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestInitStoryHarnessRoutes_ProcessStatusRouteHandlesInvalidJSON(t *testing.T) {
	engine := newWriterRouteEngine()
	v1Writer := engine.Group("/api/v1/writer")
	InitStoryHarnessRoutes(v1Writer, nil, nil, nil)

	w := performWriterRequest(
		engine,
		http.MethodPut,
		"/api/v1/writer/change-requests/cr-1/status",
		"{\"status\":",
	)

	assert.NotEqual(t, http.StatusNotFound, w.Code)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "参数错误")
}

func TestInitStoryHarnessRoutes_TriggerIndexRouteRegistered(t *testing.T) {
	engine := newWriterRouteEngine()
	v1Writer := engine.Group("/api/v1/writer")
	InitStoryHarnessRoutes(v1Writer, nil, nil, nil)

	w := performWriterRequest(
		engine,
		http.MethodPost,
		"/api/v1/writer/projects/project-1/chapters/chapter-1/trigger-index",
		"",
	)

	assert.NotEqual(t, http.StatusNotFound, w.Code)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "索引服务不可用")
}

func TestInitWriterRouter_BatchOperationRouteRegistered(t *testing.T) {
	engine := newWriterRouteEngine()
	v1 := engine.Group("/api/v1")
	originalConfig := config.GlobalConfig
	t.Cleanup(func() {
		config.GlobalConfig = originalConfig
	})
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}

	InitWriterRouter(
		v1,
		nil,
		nil,
		&documentservice.BatchOperationService{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	w := performWriterRequest(engine, http.MethodPost, "/api/v1/writer/batch-operations", "{}")

	assert.NotEqual(t, http.StatusNotFound, w.Code)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
