package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/internal/middleware/core"
)

// PermissionMiddleware 权限中间件
//
// 使用Checker检查用户权限
type PermissionMiddleware struct {
	checker Checker
	config  *PermissionConfig
	logger  *zap.Logger
}

// PermissionConfig 权限中间件配置
type PermissionConfig struct {
	// Enabled 是否启用权限检查
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Strategy 权限检查策略（rbac, casbin）
	Strategy string `json:"strategy" yaml:"strategy"`

	// ConfigPath 权限配置文件路径
	ConfigPath string `json:"config_path" yaml:"config_path"`

	// SkipPaths 跳过权限检查的路径
	SkipPaths []string `json:"skip_paths" yaml:"skip_paths"`

	// Message 权限不足时的提示信息
	Message string `json:"message" yaml:"message"`

	// StatusCode 权限不足时返回的状态码
	StatusCode int `json:"status_code" yaml:"status_code"`
}

// DefaultPermissionConfig 返回默认配置
func DefaultPermissionConfig() *PermissionConfig {
	return &PermissionConfig{
		Enabled:    true,
		Strategy:   "rbac",
		ConfigPath: "configs/permissions.yaml",
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/api/v1/auth/login",
			"/api/v1/auth/register",
		},
		Message:    "权限不足，无法访问该资源",
		StatusCode: 403,
	}
}

// NewPermissionMiddleware 创建权限中间件
func NewPermissionMiddleware(config *PermissionConfig, logger *zap.Logger) (*PermissionMiddleware, error) {
	if config == nil {
		config = DefaultPermissionConfig()
	}

	if logger == nil {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	// 创建权限检查器
	checker, err := CreateChecker(&CheckerConfig{
		Strategy:   config.Strategy,
		ConfigPath: config.ConfigPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create checker: %w", err)
	}

	return &PermissionMiddleware{
		checker: checker,
		config:  config,
		logger:  logger,
	}, nil
}

// Name 返回中间件名称
func (m *PermissionMiddleware) Name() string {
	return "permission"
}

// Priority 返回执行优先级
//
// 返回10，确保权限检查在认证之后执行
func (m *PermissionMiddleware) Priority() int {
	return 10
}

// Handler 返回Gin处理函数
func (m *PermissionMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否启用
		if !m.config.Enabled {
			c.Next()
			return
		}

		// 检查是否跳过该路径
		if m.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 从上下文获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			// 未认证用户
			m.logger.Warn("User not authenticated",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			c.JSON(401, gin.H{
				"code":    40101,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		subject, ok := userID.(string)
		if !ok {
			m.logger.Error("Invalid user_id type in context",
				zap.String("user_id", fmt.Sprintf("%v", userID)),
			)
			c.JSON(500, gin.H{
				"code":    50001,
				"message": "用户信息格式错误",
			})
			c.Abort()
			return
		}

		// 构建权限
		perm := Permission{
			Resource: getResourceFromPath(c.Request.URL.Path),
			Action:   getActionFromMethod(c.Request.Method),
		}

		// 检查权限
		allowed, err := m.checker.Check(c.Request.Context(), subject, perm)
		if err != nil {
			m.logger.Error("Permission check failed",
				zap.String("user_id", subject),
				zap.String("resource", perm.Resource),
				zap.String("action", perm.Action),
				zap.Error(err),
			)
			c.JSON(500, gin.H{
				"code":    50002,
				"message": "权限检查失败",
			})
			c.Abort()
			return
		}

		if !allowed {
			m.logger.Warn("Permission denied",
				zap.String("user_id", subject),
				zap.String("resource", perm.Resource),
				zap.String("action", perm.Action),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			c.JSON(m.config.StatusCode, gin.H{
				"code":    40301,
				"message": m.config.Message,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *PermissionMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = DefaultPermissionConfig()
	}

	// 加载Enabled
	if enabled, ok := config["enabled"].(bool); ok {
		m.config.Enabled = enabled
	}

	// 加载Strategy
	if strategy, ok := config["strategy"].(string); ok {
		m.config.Strategy = strategy
	}

	// 加载ConfigPath
	if configPath, ok := config["config_path"].(string); ok {
		m.config.ConfigPath = configPath
	}

	// 加载SkipPaths
	if skipPaths, ok := config["skip_paths"].([]interface{}); ok {
		m.config.SkipPaths = make([]string, len(skipPaths))
		for i, v := range skipPaths {
			if str, ok := v.(string); ok {
				m.config.SkipPaths[i] = str
			}
		}
	}

	// 加载Message
	if message, ok := config["message"].(string); ok {
		m.config.Message = message
	}

	// 加载StatusCode
	if statusCode, ok := config["status_code"].(int); ok {
		m.config.StatusCode = statusCode
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *PermissionMiddleware) ValidateConfig() error {
	if m.config == nil {
		return fmt.Errorf("config is nil")
	}

	if m.config.Strategy == "" {
		return fmt.Errorf("strategy cannot be empty")
	}

	if m.config.StatusCode < 100 || m.config.StatusCode > 599 {
		return fmt.Errorf("invalid status code: %d", m.config.StatusCode)
	}

	return nil
}

// Reload 热重载配置
//
// 实现HotReloadMiddleware接口
func (m *PermissionMiddleware) Reload(config map[string]interface{}) error {
	// 保存旧配置
	oldConfig := *m.config

	// 加载新配置
	if err := m.LoadConfig(config); err != nil {
		m.config = &oldConfig
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 验证新配置
	if err := m.ValidateConfig(); err != nil {
		m.config = &oldConfig
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 如果策略改变，需要重新创建检查器
	if m.config.Strategy != oldConfig.Strategy {
		m.logger.Info("Recreating checker due to strategy change",
			zap.String("old_strategy", oldConfig.Strategy),
			zap.String("new_strategy", m.config.Strategy),
		)

		// 关闭旧检查器
		if err := m.checker.Close(); err != nil {
			m.logger.Error("Failed to close old checker", zap.Error(err))
		}

		// 创建新检查器
		checker, err := CreateChecker(&CheckerConfig{
			Strategy:   m.config.Strategy,
			ConfigPath: m.config.ConfigPath,
		})
		if err != nil {
			m.config = &oldConfig
			return fmt.Errorf("failed to create new checker: %w", err)
		}

		m.checker = checker
	}

	m.logger.Info("Permission config reloaded",
		zap.String("strategy", m.config.Strategy),
		zap.Bool("enabled", m.config.Enabled),
	)

	return nil
}

// GetChecker 获取权限检查器
//
// 用于外部直接调用检查器
func (m *PermissionMiddleware) GetChecker() Checker {
	return m.checker
}

// shouldSkipPath 检查是否应该跳过权限检查
func (m *PermissionMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range m.config.SkipPaths {
		// 支持前缀匹配
		if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
			return true
		}
	}
	return false
}

// ========== 辅助函数 ==========

// getResourceFromPath 从路径提取资源类型
//
// 例如:
//   /api/v1/projects/123 -> project
//   /api/v1/users/456 -> user
//   /api/v1/documents -> document
func getResourceFromPath(path string) string {
	// 简单实现：提取路径中的第一个资源段
	// 实际项目中可能需要更复杂的路由解析

	// 跳过API版本前缀
	if len(path) > 8 && path[:8] == "/api/v1/" {
		path = path[8:]
	}

	// 提取资源名
	resource := path
	for i, ch := range path {
		if ch == '/' {
			if i > 0 {
				resource = path[:i]
			} else {
				// 跳过开头的 /
				resource = path[1:]
			}
			break
		}
	}

	// 如果还有更多路径，继续查找第一个资源段
	for i, ch := range resource {
		if ch == '/' || ch == '?' {
			resource = resource[:i]
			break
		}
	}

	// 去掉复数形式（例如 projects -> project）
	if len(resource) > 1 && resource[len(resource)-1] == 's' {
		resource = resource[:len(resource)-1]
	}

	return resource
}

// getActionFromMethod 从HTTP方法提取操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

// ========== 确保实现了核心接口 ==========

var _ core.Middleware = (*PermissionMiddleware)(nil)
var _ core.ConfigurableMiddleware = (*PermissionMiddleware)(nil)
var _ core.HotReloadMiddleware = (*PermissionMiddleware)(nil)
