package router

import (
	"Qingyu_backend/config"
	middlewareAuth "Qingyu_backend/internal/middleware/auth"
	securityjwt "Qingyu_backend/internal/security/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// newCanonicalJWTAuthMiddleware 在组合根创建一次 canonical JWT manager，
// 供已完成显式注入迁移的路由组使用。
func newCanonicalJWTAuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	jwtConfig := config.GetJWTConfigEnhanced()
	jwtManager, err := securityjwt.NewManager(jwtConfig.SecretKey, jwtConfig.Expiration, jwtConfig.RefreshDuration)
	if err != nil {
		if logger != nil {
			logger.Warn("创建 canonical JWT manager 失败，回退到兼容鉴权入口", zap.Error(err))
		}
		return middlewareAuth.JWTAuth()
	}

	return middlewareAuth.JWTAuthWithManager(jwtManager)
}
