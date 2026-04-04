package core

import (
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"Qingyu_backend/internal/middleware"
)

// InitializerImpl 中间件初始化器实现
type InitializerImpl struct {
	config      *middleware.Config
	middlewares map[string]Middleware
	mu          sync.RWMutex
	logger      *zap.Logger
}

// NewInitializer 创建初始化器
func NewInitializer(logger *zap.Logger) *InitializerImpl {
	return &InitializerImpl{
		middlewares: make(map[string]Middleware),
		logger:      logger,
	}
}

// LoadFromConfig 从配置文件加载中间件配置
func (i *InitializerImpl) LoadFromConfig(configPath string) (*middleware.Config, error) {
	i.logger.Info("Loading middleware config", zap.String("path", configPath))

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML
	var config middleware.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	i.mu.Lock()
	i.config = &config
	i.mu.Unlock()

	i.logger.Info("Config loaded successfully",
		zap.Int("static_configs", i.countStaticConfigs(&config)),
		zap.Int("dynamic_configs", i.countDynamicConfigs(&config)))

	return &config, nil
}

// Initialize 初始化所有中间件
func (i *InitializerImpl) Initialize() ([]Middleware, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.config == nil {
		return nil, fmt.Errorf("config not loaded, call LoadFromConfig first")
	}

	i.logger.Info("Initializing middlewares")

	// 清空已有中间件
	i.middlewares = make(map[string]Middleware)

	var middlewareList []Middleware
	var errs []error

	// ========== 初始化静态配置中间件 ==========

	// RequestID
	if i.config.Middleware.RequestID != nil {
		mw, err := i.createRequestIDMiddleware(i.config.Middleware.RequestID)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create request_id middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// Recovery
	if i.config.Middleware.Recovery != nil {
		mw, err := i.createRecoveryMiddleware(i.config.Middleware.Recovery)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create recovery middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// Security
	if i.config.Middleware.Security != nil {
		mw, err := i.createSecurityMiddleware(i.config.Middleware.Security)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create security middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// Logger
	if i.config.Middleware.Logger != nil {
		mw, err := i.createLoggerMiddleware(i.config.Middleware.Logger)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create logger middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// Compression
	if i.config.Middleware.Compression != nil && i.config.Middleware.Compression.Enabled {
		mw, err := i.createCompressionMiddleware(i.config.Middleware.Compression)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create compression middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// ========== 初始化动态配置中间件 ==========

	// RateLimit
	if i.config.Middleware.RateLimit != nil && i.config.Middleware.RateLimit.Enabled {
		mw, err := i.createRateLimitMiddleware(i.config.Middleware.RateLimit)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create rate_limit middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// Auth
	if i.config.Middleware.Auth != nil && i.config.Middleware.Auth.Enabled {
		mw, err := i.createAuthMiddleware(i.config.Middleware.Auth)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create auth middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// Permission
	if i.config.Middleware.Permission != nil && i.config.Middleware.Permission.Enabled {
		mw, err := i.createPermissionMiddleware(i.config.Middleware.Permission)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create permission middleware: %w", err))
		} else {
			i.middlewares[mw.Name()] = mw
			middlewareList = append(middlewareList, mw)
		}
	}

	// 处理错误
	if len(errs) > 0 {
		i.logger.Error("Some middlewares failed to initialize",
			zap.Int("failed", len(errs)),
			zap.Int("succeeded", len(middlewareList)))
		return middlewareList, fmt.Errorf("%d middlewares failed to initialize: %v", len(errs), errs)
	}

	i.logger.Info("All middlewares initialized successfully",
		zap.Int("total", len(middlewareList)))

	return middlewareList, nil
}

// GetMiddleware 获取指定中间件实例
func (i *InitializerImpl) GetMiddleware(name string) (Middleware, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	mw, exists := i.middlewares[name]
	if !exists {
		return nil, fmt.Errorf("middleware %s not found", name)
	}

	return mw, nil
}

// ListMiddlewares 列出所有已初始化的中间件
func (i *InitializerImpl) ListMiddlewares() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	names := make([]string, 0, len(i.middlewares))
	for name := range i.middlewares {
		names = append(names, name)
	}

	return names
}

// ========== 中间件创建方法 ==========

// createRequestIDMiddleware 创建RequestID中间件
func (i *InitializerImpl) createRequestIDMiddleware(config *middleware.RequestIDConfig) (Middleware, error) {
	mw := &RequestIDMiddleware{}
	if factory, ok := registeredFactories["request_id"]; ok {
		mw.handler = factory()
	}
	return mw, nil
}

// createRecoveryMiddleware 创建Recovery中间件
func (i *InitializerImpl) createRecoveryMiddleware(config *middleware.RecoveryConfig) (Middleware, error) {
	return &RecoveryMiddleware{
		stackSize:    config.StackSize,
		disablePrint: config.DisablePrint,
		logger:       i.logger,
	}, nil
}

// createSecurityMiddleware 创建Security中间件
func (i *InitializerImpl) createSecurityMiddleware(config *middleware.SecurityConfig) (Middleware, error) {
	mw := &SecurityMiddleware{}
	if factory, ok := registeredFactories["security"]; ok {
		mw.handler = factory()
	}
	return mw, nil
}

// createLoggerMiddleware 创建Logger中间件
func (i *InitializerImpl) createLoggerMiddleware(config *middleware.LoggerConfig) (Middleware, error) {
	mw := &LoggerMiddleware{}
	if factory, ok := registeredFactories["logger"]; ok {
		mw.handler = factory()
	}
	return mw, nil
}

// createCompressionMiddleware 创建Compression中间件
func (i *InitializerImpl) createCompressionMiddleware(config *middleware.CompressionConfig) (Middleware, error) {
	// 验证压缩级别
	if config.Level < 1 || config.Level > 9 {
		return nil, fmt.Errorf("invalid compression level: %d (must be 1-9)", config.Level)
	}

	mw := &CompressionMiddleware{}
	if factory, ok := registeredFactories["compression"]; ok {
		mw.handler = factory()
	}
	return mw, nil
}

// createRateLimitMiddleware 创建RateLimit中间件
func (i *InitializerImpl) createRateLimitMiddleware(config *middleware.RateLimitConfig) (Middleware, error) {
	// 验证策略
	if config.Strategy == "" {
		return nil, fmt.Errorf("rate limit strategy cannot be empty")
	}

	mw := &RateLimitMiddleware{}
	if factory, ok := registeredFactories["rate_limit"]; ok {
		mw.handler = factory()
	} else {
		i.logger.Warn("RateLimit中间件工厂未注册，使用降级模式（放行所有请求）")
	}
	return mw, nil
}

// createAuthMiddleware 创建Auth中间件
func (i *InitializerImpl) createAuthMiddleware(config *middleware.AuthConfig) (Middleware, error) {
	if config.Secret == "" {
		return nil, fmt.Errorf("auth secret cannot be empty")
	}

	mw := &AuthMiddleware{}
	if factory, ok := registeredFactories["auth"]; ok {
		mw.handler = factory()
	} else {
		i.logger.Warn("Auth中间件工厂未注册，使用降级模式（放行所有请求）")
	}
	return mw, nil
}

// createPermissionMiddleware 创建Permission中间件
func (i *InitializerImpl) createPermissionMiddleware(config *middleware.PermissionConfig) (Middleware, error) {
	if config.Strategy == "" {
		return nil, fmt.Errorf("permission strategy cannot be empty")
	}
	if config.ConfigPath == "" {
		return nil, fmt.Errorf("permission config path cannot be empty")
	}

	mw := &PermissionMiddleware{}
	if factory, ok := registeredFactories["permission"]; ok {
		mw.handler = factory()
	} else {
		i.logger.Warn("Permission中间件工厂未注册，使用降级模式（放行所有请求）")
	}
	return mw, nil
}

// ========== 辅助方法 ==========

// countStaticConfigs 统计静态配置数量
func (i *InitializerImpl) countStaticConfigs(config *middleware.Config) int {
	count := 0
	if config.Middleware.RequestID != nil {
		count++
	}
	if config.Middleware.Recovery != nil {
		count++
	}
	if config.Middleware.Security != nil {
		count++
	}
	if config.Middleware.Logger != nil {
		count++
	}
	if config.Middleware.Compression != nil {
		count++
	}
	return count
}

// countDynamicConfigs 统计动态配置数量
func (i *InitializerImpl) countDynamicConfigs(config *middleware.Config) int {
	count := 0
	if config.Middleware.RateLimit != nil {
		count++
	}
	if config.Middleware.Auth != nil {
		count++
	}
	if config.Middleware.Permission != nil {
		count++
	}
	return count
}

// ========== 中间件实现（委托到真实实现）==========

// MiddlewareFunc 是创建中间件 Handler 的工厂函数类型
type MiddlewareFunc func() gin.HandlerFunc

// registeredFactories 存储已注册的中间件工厂函数
var registeredFactories = map[string]MiddlewareFunc{}

// RegisterMiddlewareFactory 注册中间件工厂函数
// 由 init 阶段的 wireup 代码调用，将真实实现注入到 core 包
func RegisterMiddlewareFactory(name string, factory MiddlewareFunc) {
	registeredFactories[name] = factory
}

// RequestIDMiddleware 请求ID中间件
type RequestIDMiddleware struct {
	handler gin.HandlerFunc
}

func (m *RequestIDMiddleware) Name() string  { return "request_id" }
func (m *RequestIDMiddleware) Priority() int { return 1 }
func (m *RequestIDMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}

// RecoveryMiddleware 异常恢复中间件
type RecoveryMiddleware struct {
	stackSize    int
	disablePrint bool
	logger       *zap.Logger
}

func (m *RecoveryMiddleware) Name() string        { return "recovery" }
func (m *RecoveryMiddleware) Priority() int       { return 2 }
func (m *RecoveryMiddleware) Handler() gin.HandlerFunc {
	return gin.Recovery()
}

// SecurityMiddleware 安全头中间件
type SecurityMiddleware struct {
	handler gin.HandlerFunc
}

func (m *SecurityMiddleware) Name() string  { return "security" }
func (m *SecurityMiddleware) Priority() int { return 3 }
func (m *SecurityMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}

// LoggerMiddleware 日志中间件
type LoggerMiddleware struct {
	handler gin.HandlerFunc
}

func (m *LoggerMiddleware) Name() string  { return "logger" }
func (m *LoggerMiddleware) Priority() int { return 6 }
func (m *LoggerMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}

// CompressionMiddleware 压缩中间件
type CompressionMiddleware struct {
	handler gin.HandlerFunc
}

func (m *CompressionMiddleware) Name() string  { return "compression" }
func (m *CompressionMiddleware) Priority() int { return 12 }
func (m *CompressionMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	handler gin.HandlerFunc
}

func (m *RateLimitMiddleware) Name() string  { return "rate_limit" }
func (m *RateLimitMiddleware) Priority() int { return 8 }
func (m *RateLimitMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	handler gin.HandlerFunc
}

func (m *AuthMiddleware) Name() string  { return "auth" }
func (m *AuthMiddleware) Priority() int { return 9 }
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}

// PermissionMiddleware 权限中间件
type PermissionMiddleware struct {
	handler gin.HandlerFunc
}

func (m *PermissionMiddleware) Name() string  { return "permission" }
func (m *PermissionMiddleware) Priority() int { return 10 }
func (m *PermissionMiddleware) Handler() gin.HandlerFunc {
	if m.handler != nil {
		return m.handler
	}
	return func(c *gin.Context) { c.Next() }
}
