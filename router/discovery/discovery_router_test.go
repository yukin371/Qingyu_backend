package discovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	discoveryAPI "Qingyu_backend/api/v1/discovery"

	"github.com/gin-gonic/gin"
)

func TestRegisterDiscoveryRoutes_EndpointsReachable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	api := discoveryAPI.NewDiscoveryAPI(nil)

	v1 := engine.Group("/api/v1")
	RegisterDiscoveryRoutes(v1, api)

	routes := []string{
		"/api/v1/discovery/recommendations",
		"/api/v1/discovery/personalized",
		"/api/v1/discovery/new-releases",
		"/api/v1/discovery/editors-pick",
		"/api/v1/discovery/trending",
		"/api/v1/discovery/topics",
	}

	for _, route := range routes {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		if w.Code == http.StatusNotFound {
			t.Fatalf("expected route %s to be registered, got 404", route)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/discovery/track", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatal("expected route /api/v1/discovery/track to be registered, got 404")
	}

	req = httptest.NewRequest(http.MethodPut, "/api/v1/discovery/preferences", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatal("expected route /api/v1/discovery/preferences to be registered, got 404")
	}
}
