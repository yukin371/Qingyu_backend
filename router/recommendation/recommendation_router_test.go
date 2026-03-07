package recommendation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	recoAPI "Qingyu_backend/api/v1/recommendation"

	"github.com/gin-gonic/gin"
)

func TestRegisterRecommendationRoutes_TableRoutesReachable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := recoAPI.NewRecommendationAPI(nil)

	v1 := engine.Group("/api/v1")
	RegisterRecommendationRoutes(v1, api)

	// Public table route should be registered and handled (table service nil => 500, not 404)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/recommendation/tables", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("expected /tables route to be registered, got 404")
	}

	// Admin auto route should be registered and hit auth middleware (401, not 404)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/recommendation/admin/tables/auto/weekly/2026-W10", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("expected admin auto route to be registered, got 404")
	}

	// Admin manual route should be registered and hit auth middleware (401, not 404)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/recommendation/admin/tables/manual/65f000000000000000000001", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("expected admin manual route to be registered, got 404")
	}
}
