package recommendation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	recoAPI "Qingyu_backend/api/v1/recommendation"
	"Qingyu_backend/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRecommendationRoutes_TableRoutesReachable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret-key-123456",
			ExpirationHours: 24,
		},
	}
	defer func() {
		config.GlobalConfig = previousConfig
	}()

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

func TestRegisterRecommendationRoutesWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterRecommendationRoutesWithAuth(v1, recoAPI.NewRecommendationAPI(nil), func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/recommendation/admin/tables/manual/65f000000000000000000001", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}
