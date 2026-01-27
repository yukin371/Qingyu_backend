package core

import "github.com/gin-gonic/gin"

// Middleware 中间件核心接口
//
// 所有中间件必须实现此接口，提供统一的中间件抽象。
type Middleware interface {
	// Name 返回中间件唯一标识
	//
	// 名称用于日志记录、配置管理、中间件注册等场景。
	// 命名规范：使用小写字母和下划线，如 "rate_limit"、"request_id"。
	Name() string

	// Priority 返回执行优先级
	//
	// 返回值越小，中间件越先执行。
	// 默认优先级参考：
	//   1-5: 基础设施（RequestID、Recovery、Security、CORS）
	//   6-8: 监控和日志（Timeout、Logger、Metrics）
	//   9-10: 认证授权（RateLimit、Auth、Permission）
	//   11-12: 业务层（Validation、Compression）
	//
	// 可通过配置文件覆盖此默认值（参见 Manager.Register）。
	Priority() int

	// Handler 返回Gin处理函数
	//
	// 这是中间件的实际执行逻辑。
	// 应该遵循Gin中间件的最佳实践：
	//   - 使用 c.Next() 调用后续处理
	//   - 使用 c.Abort() 中断请求处理
	//   - 使用 c.Set() 传递数据到后续中间件
	Handler() gin.HandlerFunc
}

// ConfigurableMiddleware 可配置中间件接口（可选）
//
// 如果中间件需要支持配置加载和验证，应该实现此接口。
type ConfigurableMiddleware interface {
	Middleware

	// LoadConfig 从配置加载参数
	//
	// config 是从YAML配置文件解析而来的map[string]interface{}。
	// 实现时需要类型断言和错误处理。
	//
	// 示例：
	//   func (m *MyMiddleware) LoadConfig(config map[string]interface{}) error {
	//       if rate, ok := config["rate"].(int); ok {
	//           m.config.Rate = rate
	//       }
	//       return nil
	//   }
	LoadConfig(config map[string]interface{}) error

	// ValidateConfig 验证配置有效性
	//
	// 在LoadConfig之后调用，确保配置值合法。
	// 如果配置无效，返回错误以阻止服务启动。
	//
	// 示例：
	//   func (m *MyMiddleware) ValidateConfig() error {
	//       if m.config.Rate <= 0 {
	//           return errors.New("rate must be positive")
	//       }
	//       return nil
	//   }
	ValidateConfig() error
}

// HotReloadMiddleware 支持热更新的中间件（v2.0新增）
//
// 仅动态配置中间件（限流、认证、权限）需要实现此接口。
// 热更新允许在运行时重新加载配置，无需重启服务。
type HotReloadMiddleware interface {
	ConfigurableMiddleware

	// Reload 热重载配置
	//
	// 在不重启服务的情况下重新加载配置。
	// 实现时应该：
	//   1. 验证新配置有效性
	//   2. 原子性地切换到新配置
	//   3. 失败时保持旧配置不变
	//
	// 示例：
	//   func (m *RateLimitMiddleware) Reload(config map[string]interface{}) error {
	//       newConfig := m.config
	//       if err := m.LoadConfig(config); err != nil {
	//           return err
	//       }
	//       if err := m.ValidateConfig(); err != nil {
	//           return err
	//       }
	//       m.config = newConfig
	//       return nil
	//   }
	Reload(config map[string]interface{}) error
}
