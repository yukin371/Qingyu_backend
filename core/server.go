package core

import (
	"fmt"

	"Qingyu_backend/api/v1/system"
	"Qingyu_backend/config"
	"Qingyu_backend/middleware"
	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/pkg/metrics"
	pkgmiddleware "Qingyu_backend/pkg/middleware"
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
	// RequestIDMiddleware - 请求ID（最优先，为整个请求生成唯一ID）
	r.Use(pkgmiddleware.RequestIDMiddleware())

	// RecoveryMiddleware - 异常恢复（捕获panic）
	r.Use(pkgmiddleware.RecoveryMiddleware())

	// LoggerMiddleware - 结构化日志记录
	accessCfg := pkgmiddleware.DefaultAccessLogConfig()
	accessCfg.Mode = logCfg.Mode
	accessCfg.RedactKeys = logCfg.RedactKeys
	if logCfg.Request != nil {
		accessCfg.SkipPaths = logCfg.Request.SkipPaths
		accessCfg.BodyAllowPaths = logCfg.Request.BodyAllowPaths
		accessCfg.EnableBody = logCfg.Request.EnableBody || logCfg.Mode == "strict"
		accessCfg.MaxBodySize = logCfg.Request.MaxBodySize
	}
	r.Use(pkgmiddleware.LoggerMiddleware(accessCfg))

	// PrometheusMiddleware - 监控指标收集
	r.Use(metrics.Middleware())

	// RateLimitMiddleware - API限流（IP + 用户双重限流）
	rateLimitConfig := pkgmiddleware.DefaultRateLimiterConfig()
	rateLimitConfig.Rate = 100  // 每秒100个请求
	rateLimitConfig.Burst = 200 // 桶容量200
	r.Use(pkgmiddleware.RateLimitMiddleware(rateLimitConfig))

	// ErrorHandler - 统一错误处理（最后执行，处理所有错误）
	r.Use(pkgmiddleware.ErrorHandler())

	// 保留原有的中间件
	r.Use(middleware.CORSMiddleware())

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
