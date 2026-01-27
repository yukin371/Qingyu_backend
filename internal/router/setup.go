package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/internal/middleware/auth"
	"Qingyu_backend/internal/middleware/builtin"
)

// SetupMiddleware 配置并初始化所有中间件
//
// 这是中间件初始化的统一入口，负责：
// 1. 创建所有中间件实例
// 2. 配置中间件参数
// 3. 按优先级返回中间件列表
func SetupMiddleware(logger *zap.Logger) ([]gin.HandlerFunc, error) {
	var handlers []gin.HandlerFunc

	// 1. RequestID 中间件（优先级 1）
	requestIDMW := builtin.NewRequestIDMiddleware()
	handlers = append(handlers, requestIDMW.Handler())

	// 2. Recovery 中间件（优先级 2）
	recoveryMW := builtin.NewRecoveryMiddleware(logger)
	handlers = append(handlers, recoveryMW.Handler())

	// 3. ErrorHandler 中间件（优先级 3）
	errorHandlerMW := builtin.NewErrorHandlerMiddleware(logger)
	handlers = append(handlers, errorHandlerMW.Handler())

	// 4. Security 中间件（优先级 4）
	securityMW := builtin.NewSecurityMiddleware()
	handlers = append(handlers, securityMW.Handler())

	// 5. CORS 中间件（优先级 5）
	corsMW := builtin.NewCORSMiddleware()
	handlers = append(handlers, corsMW.Handler())

	// 6. Logger 中间件（优先级 7）
	loggerMW := builtin.NewLoggerMiddleware(logger)
	handlers = append(handlers, loggerMW.Handler())

	// 7. Compression 中间件（优先级 12）
	compressionMW := builtin.NewCompressionMiddleware()
	handlers = append(handlers, compressionMW.Handler())

	return handlers, nil
}

// SetupAuthMiddleware 配置认证中间件
//
// JWT 认证中间件，验证用户身份并提取用户信息到上下文
func SetupAuthMiddleware(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	// TODO: 实现JWT认证中间件
	// 这里先返回一个简单的中间件，后续需要完整实现
	return func(c *gin.Context) {
		// 暂时跳过认证
		c.Next()
	}
}

// SetupPermissionMiddleware 配置权限中间件
//
// 基于RBAC的权限检查中间件
func SetupPermissionMiddleware(configPath string, logger *zap.Logger) (*auth.PermissionMiddleware, error) {
	permConfig := &auth.PermissionConfig{
		Enabled:    true,
		Strategy:   "rbac",
		ConfigPath: configPath,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/api/v1/auth/login",
			"/api/v1/auth/register",
		},
		Message:    "权限不足，无法访问该资源",
		StatusCode: 403,
	}

	permMW, err := auth.NewPermissionMiddleware(permConfig, logger)
	if err != nil {
		return nil, err
	}

	return permMW, nil
}

// SetupRoutes 配置所有路由
//
// 这是路由配置的主入口，负责：
// 1. 注册所有中间件到全局路由
// 2. 配置各个功能模块的路由
// 3. 设置公开路由和受保护路由
func SetupRoutes(router *gin.Engine, logger *zap.Logger) error {
	// ========== 配置全局中间件 ==========
	middlewareHandlers, err := SetupMiddleware(logger)
	if err != nil {
		return err
	}

	for _, handler := range middlewareHandlers {
		router.Use(handler)
	}

	// ========== 配置权限中间件 ==========
	permMW, err := SetupPermissionMiddleware("configs/permissions.yaml", logger)
	if err != nil {
		logger.Warn("Failed to setup permission middleware", zap.Error(err))
		// 权限中间件失败不应阻止服务启动
	}

	// ========== API v1 路由组 ==========
	v1 := router.Group("/api/v1")
	{
		// 公开路由（不需要认证和权限）
		setupPublicRoutes(v1, logger)

		// 需要认证的路由
		authGroup := v1.Group("")
		authGroup.Use(SetupAuthMiddleware("", logger))
		{
			// 需要权限的路由
			if permMW != nil {
				setupProtectedRoutes(authGroup, permMW, logger)
			}
		}
	}

	// ========== 健康检查 ==========
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "qingyu-backend",
		})
	})

	return nil
}

// setupPublicRoutes 配置公开路由
//
// 公开路由不需要认证和权限检查
func setupPublicRoutes(v1 *gin.RouterGroup, logger *zap.Logger) {
	// 认证相关
	authGroup := v1.Group("/auth")
	{
		// 登录、注册等公开接口
		// TODO: 添加实际的handler
		authGroup.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Login endpoint"})
		})
		authGroup.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Register endpoint"})
		})
	}
}

// setupProtectedRoutes 配置受保护路由
//
// 受保护路由需要用户认证和权限检查
func setupProtectedRoutes(v1 *gin.RouterGroup, permMW *auth.PermissionMiddleware, logger *zap.Logger) {
	// 权限检查会自动应用
	// 路由会根据用户的角色和权限进行检查

	// 书店相关路由
	bookstoreGroup := v1.Group("/bookstore")
	{
		// 这些路由会自动进行权限检查
		// 具体权限根据用户角色自动判断
		bookstoreGroup.GET("/books", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "List books"})
		})
		bookstoreGroup.POST("/books", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Create book"})
		})
	}

	// 管理员路由（需要admin角色）
	adminGroup := v1.Group("/admin")
	{
		adminGroup.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "List users"})
		})
		adminGroup.POST("/config", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Update config"})
		})
	}
}
