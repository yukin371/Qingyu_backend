package middleware

import (
	"fmt"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// MiddlewareType 中间件类型
type MiddlewareType string

const (
	MiddlewareTypeLogger     MiddlewareType = "logger"
	MiddlewareTypeCORS       MiddlewareType = "cors"
	MiddlewareTypeRateLimit  MiddlewareType = "rate_limit"
	MiddlewareTypeRecovery   MiddlewareType = "recovery"
	MiddlewareTypeAuth       MiddlewareType = "auth"
	MiddlewareTypePermission MiddlewareType = "permission"
	MiddlewareTypeSecurity   MiddlewareType = "security"
	MiddlewareTypeTimeout    MiddlewareType = "timeout"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Type     MiddlewareType         `json:"type" yaml:"type"`
	Enabled  bool                   `json:"enabled" yaml:"enabled"`
	Priority int                    `json:"priority" yaml:"priority"`
	Config   map[string]interface{} `json:"config" yaml:"config"`
}

// MiddlewareFactory 中间件工厂
type MiddlewareFactory struct {
	creators map[MiddlewareType]MiddlewareCreator
}

// MiddlewareCreator 中间件创建器
type MiddlewareCreator func(config map[string]interface{}) (gin.HandlerFunc, error)

// NewMiddlewareFactory 创建中间件工厂
func NewMiddlewareFactory() *MiddlewareFactory {
	factory := &MiddlewareFactory{
		creators: make(map[MiddlewareType]MiddlewareCreator),
	}
	
	// 注册默认中间件创建器
	factory.RegisterCreator(MiddlewareTypeLogger, CreateLoggerMiddleware)
	factory.RegisterCreator(MiddlewareTypeCORS, CreateCORSMiddleware)
	factory.RegisterCreator(MiddlewareTypeRateLimit, CreateRateLimitMiddleware)
	factory.RegisterCreator(MiddlewareTypeRecovery, CreateRecoveryMiddleware)
	factory.RegisterCreator(MiddlewareTypeAuth, CreateAuthMiddleware)
	factory.RegisterCreator(MiddlewareTypePermission, CreatePermissionMiddleware)
	factory.RegisterCreator(MiddlewareTypeSecurity, CreateSecurityMiddleware)
	factory.RegisterCreator(MiddlewareTypeTimeout, CreateTimeoutMiddleware)
	
	return factory
}

// RegisterCreator 注册中间件创建器
func (f *MiddlewareFactory) RegisterCreator(middlewareType MiddlewareType, creator MiddlewareCreator) {
	f.creators[middlewareType] = creator
}

// CreateMiddleware 创建单个中间件
func (f *MiddlewareFactory) CreateMiddleware(config MiddlewareConfig) (gin.HandlerFunc, error) {
	creator, exists := f.creators[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown middleware type: %s", config.Type)
	}
	
	if !config.Enabled {
		return nil, nil
	}
	
	return creator(config.Config)
}

// CreateMiddlewares 创建多个中间件
func (f *MiddlewareFactory) CreateMiddlewares(configs []MiddlewareConfig) ([]gin.HandlerFunc, error) {
	var middlewares []gin.HandlerFunc
	
	// 按优先级排序
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Priority < configs[j].Priority
	})
	
	for _, config := range configs {
		middleware, err := f.CreateMiddleware(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create middleware %s: %w", config.Type, err)
		}
		
		if middleware != nil {
			middlewares = append(middlewares, middleware)
		}
	}
	
	return middlewares, nil
}

// MiddlewareChainBuilder 中间件链构建器
type MiddlewareChainBuilder struct {
	factory     *MiddlewareFactory
	middlewares []MiddlewareConfig
}

// NewMiddlewareChainBuilder 创建中间件链构建器
func NewMiddlewareChainBuilder() *MiddlewareChainBuilder {
	return &MiddlewareChainBuilder{
		factory:     NewMiddlewareFactory(),
		middlewares: make([]MiddlewareConfig, 0),
	}
}

// Add 添加中间件配置
func (b *MiddlewareChainBuilder) Add(config MiddlewareConfig) *MiddlewareChainBuilder {
	b.middlewares = append(b.middlewares, config)
	return b
}

// AddLogger 添加日志中间件
func (b *MiddlewareChainBuilder) AddLogger(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeLogger,
		Enabled:  true,
		Priority: 1,
		Config:   config,
	})
}

// AddCORS 添加CORS中间件
func (b *MiddlewareChainBuilder) AddCORS(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeCORS,
		Enabled:  true,
		Priority: 2,
		Config:   config,
	})
}

// AddRateLimit 添加限流中间件
func (b *MiddlewareChainBuilder) AddRateLimit(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeRateLimit,
		Enabled:  true,
		Priority: 3,
		Config:   config,
	})
}

// AddRecovery 添加恢复中间件
func (b *MiddlewareChainBuilder) AddRecovery(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeRecovery,
		Enabled:  true,
		Priority: 4,
		Config:   config,
	})
}

// AddAuth 添加认证中间件
func (b *MiddlewareChainBuilder) AddAuth(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeAuth,
		Enabled:  true,
		Priority: 5,
		Config:   config,
	})
}

// AddPermission 添加权限中间件
func (b *MiddlewareChainBuilder) AddPermission(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypePermission,
		Enabled:  true,
		Priority: 6,
		Config:   config,
	})
}

// AddSecurity 添加安全中间件
func (b *MiddlewareChainBuilder) AddSecurity(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeSecurity,
		Enabled:  true,
		Priority: 7,
		Config:   config,
	})
}

// AddTimeout 添加超时中间件
func (b *MiddlewareChainBuilder) AddTimeout(config map[string]interface{}) *MiddlewareChainBuilder {
	return b.Add(MiddlewareConfig{
		Type:     MiddlewareTypeTimeout,
		Enabled:  true,
		Priority: 8,
		Config:   config,
	})
}

// Build 构建中间件链
func (b *MiddlewareChainBuilder) Build() ([]gin.HandlerFunc, error) {
	return b.factory.CreateMiddlewares(b.middlewares)
}

// RouteGroupConfig 路由组配置
type RouteGroupConfig struct {
	Path        string             `json:"path" yaml:"path"`
	Middlewares []MiddlewareConfig `json:"middlewares" yaml:"middlewares"`
	Routes      []RouteConfig      `json:"routes" yaml:"routes"`
}

// RouteConfig 路由配置
type RouteConfig struct {
	Method      string             `json:"method" yaml:"method"`
	Path        string             `json:"path" yaml:"path"`
	Handler     string             `json:"handler" yaml:"handler"`
	Middlewares []MiddlewareConfig `json:"middlewares" yaml:"middlewares"`
}

// ApplyRouteConfig 应用路由配置
func ApplyRouteConfig(engine *gin.Engine, groupConfigs []RouteGroupConfig) error {
	factory := NewMiddlewareFactory()
	
	for _, groupConfig := range groupConfigs {
		// 创建路由组中间件
		groupMiddlewares, err := factory.CreateMiddlewares(groupConfig.Middlewares)
		if err != nil {
			return fmt.Errorf("failed to create group middlewares for %s: %w", groupConfig.Path, err)
		}
		
		// 创建路由组
		_ = engine.Group(groupConfig.Path, groupMiddlewares...)
		
		// 应用路由
		for _, routeConfig := range groupConfig.Routes {
			// 创建路由中间件
			_, err := factory.CreateMiddlewares(routeConfig.Middlewares)
			if err != nil {
				return fmt.Errorf("failed to create route middlewares for %s %s: %w", 
					routeConfig.Method, routeConfig.Path, err)
			}
			
			// 注册路由（这里需要根据实际的处理器注册逻辑来实现）
			// 示例：group.Handle(routeConfig.Method, routeConfig.Path, append(routeMiddlewares, handler)...)
		}
	}
	
	return nil
}

// CreateAuthMiddleware 创建认证中间件（占位符实现）
func CreateAuthMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	// 这里应该根据实际的认证逻辑来实现
	// 暂时返回JWT中间件
	return JWTAuth(), nil
}

// CreateLoggerMiddleware 创建日志中间件的包装函数
func CreateLoggerMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	loggerConfig := LoggerConfig{
		EnableColor:       true,
		EnableReqBody:     false,
		EnableRespBody:    false,
		MaxReqBodySize:    1024,
		MaxRespBodySize:   1024,
		SkipPaths:         []string{"/health", "/metrics"},
		SlowThreshold:     time.Second * 2,
		EnablePerformance: true,
	}
	
	// 解析配置
	if enableColor, ok := config["enable_color"].(bool); ok {
		loggerConfig.EnableColor = enableColor
	}
	if enableReqBody, ok := config["enable_req_body"].(bool); ok {
		loggerConfig.EnableReqBody = enableReqBody
	}
	if maxReqBodySize, ok := config["max_req_body_size"].(int); ok {
		loggerConfig.MaxReqBodySize = maxReqBodySize
	}
	if skipPaths, ok := config["skip_paths"].([]string); ok {
		loggerConfig.SkipPaths = skipPaths
	}
	if slowThreshold, ok := config["slow_threshold"].(string); ok {
		if duration, err := time.ParseDuration(slowThreshold); err == nil {
			loggerConfig.SlowThreshold = duration
		}
	}
	if enablePerformance, ok := config["enable_performance"].(bool); ok {
		loggerConfig.EnablePerformance = enablePerformance
	}
	
	return LoggerWithConfig(loggerConfig), nil
}