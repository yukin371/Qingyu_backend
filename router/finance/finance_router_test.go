package finance

import (
	"testing"

	financeAPI "Qingyu_backend/api/v1/finance"
	"Qingyu_backend/config"

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

func TestRegisterFinanceRoutes_RegistersWalletCompatibilityAliases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.GlobalConfig = &config.Config{
		JWT: &config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}

	engine := gin.New()
	v1 := engine.Group("/api/v1")

	RegisterFinanceRoutes(v1, financeAPI.NewWalletAPI(nil), nil, nil)

	paths := collectRoutePaths(engine)

	assert.True(t, paths["/api/v1/finance/wallet"])
	assert.True(t, paths["/api/v1/finance/wallet/balance"])
	assert.True(t, paths["/api/v1/finance/wallet/detail"])
	assert.True(t, paths["/api/v1/finance/wallet/withdraws"])
	assert.True(t, paths["/api/v1/finance/wallet/withdrawals"])
}
