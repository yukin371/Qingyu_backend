package main

import (
	"fmt"
	"log"

	"Qingyu_backend/api/v1/system"
	"Qingyu_backend/config"
	"Qingyu_backend/core"
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

// @title           青羽写作平台 API (完整演示版)
// @version         1.0
// @description     青羽写作平台后端服务API文档 - 此版本展示完整的服务初始化流程
// @host            localhost:9090
// @BasePath        /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. 加载配置
	if err := loadConfiguration(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 2. 启用配置热重载
	enableConfigHotReload()

	// 3. 初始化日志系统
	if err := initializeLogger(); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}

	// 4. 初始化服务容器（包含数据库连接）
	if err := initializeServiceContainer(); err != nil {
		log.Fatalf("服务容器初始化失败: %v", err)
	}

	// 5. 设置Gin运行模式
	setupGinMode()

	// 6. 创建Gin引擎实例
	r := gin.New()

	// 7. 应用中间件（按P0优先级顺序）
	applyMiddleware(r)

	// 8. 注册系统路由（健康检查、监控）
	registerSystemRoutes(r)

	// 9. 注册业务路由
	registerBusinessRoutes(r)

	// 10. 注册Swagger文档路由
	registerSwaggerRoute(r)

	// 11. 启动服务器
	if err := startServer(r); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}

// loadConfiguration 加载配置文件
func loadConfiguration() error {
	_, err := config.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	log.Println("[主程序] 配置加载成功")
	return nil
}

// enableConfigHotReload 启用配置热重载
func enableConfigHotReload() {
	// 注册配置重载处理器
	config.RegisterReloadHandler("database", func() {
		log.Println("[配置热重载] 数据库配置已变更，请重启服务以生效")
	})

	config.RegisterReloadHandler("server", func() {
		log.Println("[配置热重载] 服务器配置已重新加载")
	})

	config.EnableHotReload()
	log.Println("[主程序] 配置热重载已启用")
}

// initializeLogger 初始化日志系统
func initializeLogger() error {
	cfg := config.GlobalConfig
	if cfg == nil || cfg.Server == nil {
		return fmt.Errorf("配置未加载")
	}

	// 构建日志配置
	logCfg := cfg.Log
	if logCfg == nil {
		logCfg = &config.LogConfig{
			Level:       "info",
			Format:      "json",
			Output:      "stdout",
			Filename:    "logs/app.log",
			Development: cfg.Server.Mode == "debug",
			Mode:        "normal",
		}
	}

	loggerConfig := &logger.Config{
		Level:       logCfg.Level,
		Format:      logCfg.Format,
		Output:      logCfg.Output,
		Filename:    logCfg.Filename,
		Development: logCfg.Development || cfg.Server.Mode == "debug",
		StrictMode:  logCfg.Mode == "strict",
	}

	if err := logger.Init(loggerConfig); err != nil {
		return fmt.Errorf("日志系统初始化失败: %w", err)
	}

	logger.Info("日志系统初始化完成",
		zap.String("level", logCfg.Level),
		zap.String("format", logCfg.Format),
	)
	return nil
}

// initializeServiceContainer 初始化服务容器
// ServiceContainer会自动创建MongoDB连接、Repository工厂和所有业务服务
func initializeServiceContainer() error {
	if err := core.InitServices(); err != nil {
		return fmt.Errorf("服务容器初始化失败: %w", err)
	}

	logger.Info("服务容器初始化完成",
		zap.String("status", "ready"),
	)
	return nil
}

// setupGinMode 设置Gin运行模式
func setupGinMode() {
	cfg := config.GlobalConfig.Server
	gin.SetMode(cfg.Mode)
	logger.Info("Gin运行模式设置完成", zap.String("mode", cfg.Mode))
}

// applyMiddleware 应用中间件（按P0优先级顺序）
// 中间件顺序非常重要，请勿随意更改
func applyMiddleware(r *gin.Engine) {
	// 1. RequestIDMiddleware - 请求ID（最优先，为整个请求生成唯一ID）
	requestIDMW := builtin.NewRequestIDMiddleware()
	r.Use(requestIDMW.Handler())
	logger.Info("中间件已注册: RequestID")

	// 2. RecoveryMiddleware - 异常恢复（捕获panic，防止服务崩溃）
	recoveryMW := builtin.NewRecoveryMiddleware(logger.Get().Logger)
	r.Use(recoveryMW.Handler())
	logger.Info("中间件已注册: Recovery")

	// 3. LoggerMiddleware - 结构化日志记录
	loggerMW := builtin.NewLoggerMiddleware(logger.Get().Logger)
	if loggerConfig := buildLoggerMiddlewareConfig(); loggerConfig != nil {
		if err := loggerMW.LoadConfig(loggerConfig); err != nil {
			logger.Warn("日志中间件配置加载失败", zap.Error(err))
		}
	}
	r.Use(loggerMW.Handler())
	logger.Info("中间件已注册: Logger")

	// 4. PrometheusMiddleware - 监控指标收集
	r.Use(metrics.Middleware())
	logger.Info("中间件已注册: Prometheus")

	// 5. ErrorHandlerMiddleware - 统一错误处理
	errorHandlerMW := builtin.NewErrorHandlerMiddleware(logger.Get().Logger)
	r.Use(errorHandlerMW.Handler())
	logger.Info("中间件已注册: ErrorHandler")

	// 6. CORSMiddleware - 跨域处理
	corsMW := builtin.NewCORSMiddleware()
	r.Use(corsMW.Handler())
	logger.Info("中间件已注册: CORS")

	// 7. RateLimitMiddleware - API限流（可选，根据配置启用）
	applyRateLimitMiddleware(r)
}

// applyRateLimitMiddleware 应用限流中间件（可选）
func applyRateLimitMiddleware(r *gin.Engine) {
	cfg := config.GlobalConfig.RateLimit
	if cfg == nil || !cfg.Enabled {
		logger.Info("限流中间件: 已禁用")
		return
	}

	rateLimitConfig := buildRateLimitConfig(cfg)
	rateLimitMW, err := ratelimit.NewRateLimitMiddleware(rateLimitConfig, logger.Get().Logger)
	if err != nil {
		logger.Warn("限流中间件创建失败", zap.Error(err))
		return
	}

	r.Use(rateLimitMW.Handler())
	logger.Info("限流中间件已启用",
		zap.Int("rate", rateLimitConfig.Rate),
		zap.Int("burst", rateLimitConfig.Burst),
	)
}

// registerSystemRoutes 注册系统路由（健康检查、监控）
func registerSystemRoutes(r *gin.Engine) {
	healthAPI := system.NewHealthAPI()

	// 系统整体健康检查
	r.GET("/health", healthAPI.SystemHealth)

	// Kubernetes存活检查
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "alive"})
	})

	// Kubernetes就绪检查
	r.GET("/health/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	// Prometheus指标端点
	r.GET("/metrics", metrics.GinMetricsHandler())

	logger.Info("系统路由已注册: /health, /metrics")
}

// registerBusinessRoutes 注册业务路由
func registerBusinessRoutes(r *gin.Engine) {
	router.RegisterRoutes(r)
	logger.Info("业务路由已注册")
}

// registerSwaggerRoute 注册Swagger文档路由
func registerSwaggerRoute(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info("Swagger文档路由已注册: /swagger/*")
}

// startServer 启动HTTP服务器
func startServer(r *gin.Engine) error {
	cfg := config.GlobalConfig.Server
	addr := fmt.Sprintf(":%s", cfg.Port)

	fmt.Printf("\n========================================\n")
	fmt.Printf("  青羽写作平台后端服务启动成功\n")
	fmt.Printf("  地址: http://localhost:%s\n", cfg.Port)
	fmt.Printf("  模式: %s\n", cfg.Mode)
	fmt.Printf("  文档: http://localhost:%s/swagger/index.html\n", cfg.Port)
	fmt.Printf("========================================\n\n")

	return r.Run(addr)
}

// buildLoggerMiddlewareConfig 构建日志中间件配置
func buildLoggerMiddlewareConfig() map[string]interface{} {
	logCfg := config.GlobalConfig.Log
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

// buildRateLimitConfig 构建限流中间件配置
func buildRateLimitConfig(cfg *config.RateLimitConfig) *ratelimit.RateLimitConfig {
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

// stringSliceToInterfaceSlice 字符串切片转换为接口切片
func stringSliceToInterfaceSlice(values []string) []interface{} {
	out := make([]interface{}, 0, len(values))
	for _, v := range values {
		out = append(out, v)
	}
	return out
}
