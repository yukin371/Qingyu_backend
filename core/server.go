package core

import (
	"fmt"

	"Qingyu_backend/api/v1/system"
	"Qingyu_backend/config"
	"Qingyu_backend/internal/middleware/builtin"
	"Qingyu_backend/internal/middleware/ratelimit"
	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/pkg/metrics"
	"Qingyu_backend/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// InitServer 初始化服务器
func InitServer() (*gin.Engine, error) {
	cfg := config.GlobalConfig.Server
	if cfg == nil {
		return nil, fmt.Errorf("server configuration is missing")
	}

	// 1. 初始化日志系统（P0中间件）
	logCfg := config.GlobalConfig.Log
	if logCfg == nil {
		logCfg = &config.LogConfig{
			Level:       "info",
			Format:      "json",
			Output:      "stdout",
			Filename:    "logs/app.log",
			Development: cfg.Mode == "debug",
			Mode:        "normal",
		}
	}

	loggerConfig := &logger.Config{
		Level:       logCfg.Level,
		Format:      logCfg.Format,
		Output:      logCfg.Output,
		Filename:    logCfg.Filename,
		Development: logCfg.Development || cfg.Mode == "debug",
		StrictMode:  logCfg.Mode == "strict",
	}
	if err := logger.Init(loggerConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	logger.Info("Logger initialized", zap.String("module", "init"))

	// 初始化服务容器
	if err := InitServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	// 设置gin模式
	gin.SetMode(cfg.Mode)

	// 创建gin实例
	r := gin.New()

	// 2. 应用P0中间件（顺序很重要）
	// 使用新架构的中间件 internal/middleware

	// RequestIDMiddleware - 请求ID（最优先，为整个请求生成唯一ID）
	r.Use(builtin.NewRequestIDMiddleware().Handler())

	// RecoveryMiddleware - 异常恢复（捕获panic）
	recoveryMW := builtin.NewRecoveryMiddleware(logger.Get().Logger)
	r.Use(recoveryMW.Handler())

	// LoggerMiddleware - 结构化日志记录
	// 使用 builtin 版本，支持敏感信息脱敏、严格模式等功能
	loggerMW := builtin.NewLoggerMiddleware(logger.Get().Logger)
	if loggerConfig := buildLoggerMiddlewareConfig(logCfg); loggerConfig != nil {
		if err := loggerMW.LoadConfig(loggerConfig); err != nil {
			logger.Warn("Failed to load logger config", zap.Error(err))
		}
	}
	r.Use(loggerMW.Handler())

	// PrometheusMiddleware - 监控指标收集
	r.Use(metrics.Middleware())

	// ErrorHandler - 统一错误处理
	errorHandlerMW := builtin.NewErrorHandlerMiddleware(logger.Get().Logger)
	r.Use(errorHandlerMW.Handler())

	// CORSMiddleware - 跨域处理
	corsMW := builtin.NewCORSMiddleware()
	r.Use(corsMW.Handler())

	// RateLimitMiddleware - API限流（支持配置化启用/禁用）
	if config.GlobalConfig.RateLimit != nil && config.GlobalConfig.RateLimit.Enabled {
		// 使用新架构的限流中间件
		rateLimitConfig := buildRateLimitMiddlewareConfig(config.GlobalConfig.RateLimit)

		// 创建限流中间件
		rateLimitMW, err := ratelimit.NewRateLimitMiddleware(rateLimitConfig, logger.Get().Logger)
		if err != nil {
			logger.Warn("Failed to create rate limit middleware", zap.Error(err))
		} else {
			r.Use(rateLimitMW.Handler())
			logger.Info("Rate limit middleware enabled",
				zap.Int("rate", rateLimitConfig.Rate),
				zap.Int("burst", rateLimitConfig.Burst),
				zap.Strings("skip_paths", rateLimitConfig.SkipPaths))
		}
	} else {
		logger.Info("Rate limit middleware disabled")
	}

	// 自定义JSON渲染 - 设置不转义HTML字符（包括中文）
	// 注意：我们通过在API层使用response.JSON()函数来实现
	// 在pkg/response/json_renderer.go中提供了JsonWithNoEscape函数

	// 3. 注册健康检查和监控端点
	healthAPI := system.NewHealthAPI()
	r.GET("/health", healthAPI.SystemHealth)     // 系统整体健康检查
	r.GET("/health/live", func(c *gin.Context) { // K8s存活检查
		c.JSON(200, gin.H{"status": "alive"})
	})
	r.GET("/health/ready", func(c *gin.Context) { // K8s就绪检查
		c.JSON(200, gin.H{"status": "ready"})
	})
	r.GET("/metrics", metrics.GinMetricsHandler()) // Prometheus指标端点（Gin处理器）

	// 注册业务路由
	router.RegisterRoutes(r)

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("Server initialized successfully", zap.String("module", "init"))

	return r, nil
}

// RunServer 运行服务器
func RunServer(r *gin.Engine) error {
	cfg := config.GlobalConfig.Server
	if cfg == nil {
		return fmt.Errorf("server configuration is missing")
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server is running on port %s in %s mode\n", cfg.Port, cfg.Mode)
	return r.Run(addr)
}

func buildLoggerMiddlewareConfig(logCfg *config.LogConfig) map[string]interface{} {
	if logCfg == nil || logCfg.Request == nil {
		return nil
	}

	return map[string]interface{}{
		"skip_paths":          stringSliceToInterfaceSlice(logCfg.Request.SkipPaths),
		"body_allow_paths":    stringSliceToInterfaceSlice(logCfg.Request.BodyAllowPaths),
		"enable_request_body": logCfg.Request.EnableBody || logCfg.Mode == "strict",
		"max_body_size":       logCfg.Request.MaxBodySize,
		"redact_keys":         stringSliceToInterfaceSlice(logCfg.RedactKeys),
		"mode":                logCfg.Mode,
	}
}

func buildRateLimitMiddlewareConfig(cfg *config.RateLimitConfig) *ratelimit.RateLimitConfig {
	if cfg == nil {
		return nil
	}

	rate := int(cfg.RequestsPerSec)
	if rate <= 0 {
		rate = 1
	}

	burst := cfg.Burst
	if burst < rate {
		burst = rate
	}

	return &ratelimit.RateLimitConfig{
		Enabled:         cfg.Enabled,
		Strategy:        "token_bucket",
		Rate:            rate,
		Burst:           burst,
		KeyFunc:         "ip",
		SkipPaths:       cfg.SkipPaths,
		Message:         "请求过于频繁，请稍后再试",
		StatusCode:      429,
		CleanupInterval: 300,
	}
}

func stringSliceToInterfaceSlice(values []string) []interface{} {
	out := make([]interface{}, 0, len(values))
	for _, v := range values {
		out = append(out, v)
	}
	return out
}
