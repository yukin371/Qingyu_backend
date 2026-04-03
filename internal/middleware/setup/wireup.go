package setup

import (
	"Qingyu_backend/internal/middleware/auth"
	"Qingyu_backend/internal/middleware/builtin"
	"Qingyu_backend/internal/middleware/core"
	"Qingyu_backend/internal/middleware/ratelimit"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterBuiltinMiddlewareFactories 注册所有内置中间件的工厂函数
// 将真实实现连接到 core 包中的占位中间件，打破 core -> builtin/auth/ratelimit 的循环依赖
func RegisterBuiltinMiddlewareFactories(logger *zap.Logger) {
	// RequestID - 无依赖
	core.RegisterMiddlewareFactory("request_id", func() gin.HandlerFunc {
		return builtin.NewRequestIDMiddleware().Handler()
	})

	// Security - 无依赖
	core.RegisterMiddlewareFactory("security", func() gin.HandlerFunc {
		return builtin.NewSecurityMiddleware().Handler()
	})

	// Logger - 依赖 zap.Logger
	core.RegisterMiddlewareFactory("logger", func() gin.HandlerFunc {
		return builtin.NewLoggerMiddleware(logger).Handler()
	})

	// Compression - 无依赖
	core.RegisterMiddlewareFactory("compression", func() gin.HandlerFunc {
		return builtin.NewCompressionMiddleware().Handler()
	})

	// RateLimit - 依赖配置和 logger
	// 注意：RateLimit 需要配置参数，这里注册的是默认配置的工厂
	// 实际使用时应通过 RegisterRateLimitFactory 传入自定义配置
	core.RegisterMiddlewareFactory("rate_limit", func() gin.HandlerFunc {
		rlConfig := ratelimit.DefaultRateLimitConfig()
		rlConfig.Enabled = true
		impl, err := ratelimit.NewRateLimitMiddleware(rlConfig, logger)
		if err != nil {
			if logger != nil {
				logger.Warn("RateLimit中间件创建失败，降级为放行",
					zap.Error(err),
				)
			}
			return func(c *gin.Context) { c.Next() }
		}
		return impl.Handler()
	})

	// Auth - 依赖 JWTManager 和 logger
	// 注意：Auth 需要密钥配置，这里注册的是占位工厂
	// 实际使用时应通过 RegisterAuthFactory 传入真实依赖
	core.RegisterMiddlewareFactory("auth", func() gin.HandlerFunc {
		if logger != nil {
			logger.Warn("Auth中间件使用占位工厂，请通过 RegisterAuthFactory 注册真实实现")
		}
		return func(c *gin.Context) { c.Next() }
	})

	// Permission - 依赖配置和 logger
	// 注意：Permission 需要配置参数，这里注册的是占位工厂
	// 实际使用时应通过 RegisterPermissionFactory 传入真实依赖
	core.RegisterMiddlewareFactory("permission", func() gin.HandlerFunc {
		if logger != nil {
			logger.Warn("Permission中间件使用占位工厂，请通过 RegisterPermissionFactory 注册真实实现")
		}
		return func(c *gin.Context) { c.Next() }
	})
}

// RegisterAuthFactory 注册带真实依赖的 Auth 中间件工厂
func RegisterAuthFactory(jwtMgr auth.JWTManager, blacklist auth.Blacklist, logger *zap.Logger) {
	core.RegisterMiddlewareFactory("auth", func() gin.HandlerFunc {
		impl := auth.NewJWTAuthMiddleware(jwtMgr, blacklist, logger)
		return impl.Handler()
	})
}

// RegisterRateLimitFactory 注册带自定义配置的 RateLimit 中间件工厂
func RegisterRateLimitFactory(rlConfig *ratelimit.RateLimitConfig, logger *zap.Logger) {
	core.RegisterMiddlewareFactory("rate_limit", func() gin.HandlerFunc {
		impl, err := ratelimit.NewRateLimitMiddleware(rlConfig, logger)
		if err != nil {
			if logger != nil {
				logger.Warn("RateLimit中间件创建失败，降级为放行",
					zap.Error(err),
				)
			}
			return func(c *gin.Context) { c.Next() }
		}
		return impl.Handler()
	})
}

// RegisterPermissionFactory 注册带自定义配置的 Permission 中间件工厂
func RegisterPermissionFactory(permConfig *auth.PermissionConfig, logger *zap.Logger) {
	core.RegisterMiddlewareFactory("permission", func() gin.HandlerFunc {
		impl, err := auth.NewPermissionMiddleware(permConfig, logger)
		if err != nil {
			if logger != nil {
				logger.Warn("Permission中间件创建失败，降级为放行",
					zap.Error(err),
				)
			}
			return func(c *gin.Context) { c.Next() }
		}
		return impl.Handler()
	})
}
