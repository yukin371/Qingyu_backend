package admin

import (
	"testing"

	"Qingyu_backend/config"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func collectRoutePaths(engine *gin.Engine) map[string]bool {
	paths := make(map[string]bool, len(engine.Routes()))
	for _, route := range engine.Routes() {
		paths[route.Path] = true
	}
	return paths
}

func TestRegisterAdminRoutes_RegistersQuotaRoutesWhenServicesAvailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}

	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterAdminRoutes(
		v1,
		nil,
		&aiService.QuotaService{},
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
		&aiService.QuotaAdminService{},
		&aiService.QuotaDashboardService{},
		&aiService.QuotaPolicyService{},
		&aiService.QuotaAlertService{},
	)

	paths := collectRoutePaths(engine)

	assert.True(t, paths["/api/v1/admin/quota/dashboard"])
	assert.True(t, paths["/api/v1/admin/quota/statistics/global"])
	assert.True(t, paths["/api/v1/admin/quota/statistics/trend"])
	assert.True(t, paths["/api/v1/admin/quota/dashboard/refresh"])
	assert.True(t, paths["/api/v1/admin/quota/users"])
	assert.True(t, paths["/api/v1/admin/quota/users/:userId"])
	assert.True(t, paths["/api/v1/admin/quota/policies"])
	assert.True(t, paths["/api/v1/admin/quota/alerts"])
}

func TestRegisterAdminRoutes_DoesNotRegisterQuotaRoutesWhenServicesMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}

	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterAdminRoutes(
		v1,
		nil,
		&aiService.QuotaService{},
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
		&aiService.QuotaDashboardService{},
		&aiService.QuotaPolicyService{},
		&aiService.QuotaAlertService{},
	)

	paths := collectRoutePaths(engine)

	assert.False(t, paths["/api/v1/admin/quota/dashboard"])
	assert.False(t, paths["/api/v1/admin/quota/users"])
	assert.False(t, paths["/api/v1/admin/quota/policies"])
	assert.False(t, paths["/api/v1/admin/quota/alerts"])
}
