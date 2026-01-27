package ratelimit

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/internal/middleware/core"
)

// RateLimitMiddleware 限流中间件
//
// 支持多种限流策略：令牌桶、滑动窗口、Redis分布式
type RateLimitMiddleware struct {
	config     *RateLimitConfig
	limiter    Limiter
	keyFunc    KeyFunc
	logger     *zap.Logger
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(config *RateLimitConfig, logger *zap.Logger) (*RateLimitMiddleware, error) {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if logger == nil {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	// 创建限流器
	var limiter Limiter
	var err error

	switch config.Strategy {
	case "token_bucket":
		limiter, err = NewTokenBucketLimiter(config)
	case "sliding_window":
		limiter, err = NewSlidingWindowLimiter(config)
	case "redis":
		limiter, err = NewRedisLimiter(config)
	default:
		return nil, fmt.Errorf("unsupported strategy: %s", config.Strategy)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create limiter: %w", err)
	}

	return &RateLimitMiddleware{
		config:  config,
		limiter: limiter,
		keyFunc: config.GetKeyFunc(),
		logger:  logger,
	}, nil
}

// Name 返回中间件名称
func (m *RateLimitMiddleware) Name() string {
	return "rate_limit"
}

// Priority 返回执行优先级
//
// 返回8，确保限流在认证和授权之后执行
func (m *RateLimitMiddleware) Priority() int {
	return 8
}

// Handler 返回Gin处理函数
func (m *RateLimitMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否启用
		if !m.config.Enabled {
			c.Next()
			return
		}

		// 检查是否跳过该路径
		if m.config.ShouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 构建限流上下文
		rateLimitCtx := &RateLimitContext{
			Context:  c,
			ClientIP: c.ClientIP(),
			Path:     c.Request.URL.Path,
			Method:   c.Request.Method,
			Metadata: make(map[string]interface{}),
		}

		// 尝试从上下文中获取用户ID
		if userID, exists := c.Get("user_id"); exists {
			if userIDStr, ok := userID.(string); ok {
				rateLimitCtx.UserID = userIDStr
			}
		}

		// 生成限流键
		key := m.keyFunc(rateLimitCtx)

		// 检查是否允许请求
		if !m.limiter.Allow(key) {
			m.logger.Warn("Rate limit exceeded",
				zap.String("key", key),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("client_ip", c.ClientIP()),
			)

			c.JSON(m.config.StatusCode, gin.H{
				"code":    42901,
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
func (m *RateLimitMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = DefaultRateLimitConfig()
	}

	// 加载Enabled
	if enabled, ok := config["enabled"].(bool); ok {
		m.config.Enabled = enabled
	}

	// 加载Strategy
	if strategy, ok := config["strategy"].(string); ok {
		m.config.Strategy = strategy
	}

	// 加载Rate
	if rate, ok := config["rate"].(int); ok {
		m.config.Rate = rate
	}

	// 加载Burst
	if burst, ok := config["burst"].(int); ok {
		m.config.Burst = burst
	}

	// 加载WindowSize
	if windowSize, ok := config["window_size"].(int); ok {
		m.config.WindowSize = windowSize
	}

	// 加载KeyFunc
	if keyFunc, ok := config["key_func"].(string); ok {
		m.config.KeyFunc = keyFunc
		m.keyFunc = m.config.GetKeyFunc()
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

	// 加载SkipSuccessful
	if skipSuccessful, ok := config["skip_successful"].(bool); ok {
		m.config.SkipSuccessful = skipSuccessful
	}

	// 加载SkipFailedRequest
	if skipFailed, ok := config["skip_failed_request"].(bool); ok {
		m.config.SkipFailedRequest = skipFailed
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
func (m *RateLimitMiddleware) ValidateConfig() error {
	return m.config.Validate()
}

// Reload 热重载配置
//
// 实现HotReloadMiddleware接口
func (m *RateLimitMiddleware) Reload(config map[string]interface{}) error {
	// 保存旧配置
	oldConfig := *m.config

	// 加载新配置
	if err := m.LoadConfig(config); err != nil {
		// 加载失败，恢复旧配置
		m.config = &oldConfig
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 验证新配置
	if err := m.ValidateConfig(); err != nil {
		// 验证失败，恢复旧配置
		m.config = &oldConfig
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 策略变更时需要重新创建限流器
	if m.config.Strategy != oldConfig.Strategy {
		m.logger.Info("Recreating limiter due to strategy change",
			zap.String("old_strategy", oldConfig.Strategy),
			zap.String("new_strategy", m.config.Strategy),
		)

		// 停止旧限流器
		if stopper, ok := m.limiter.(interface{ Stop() }); ok {
			stopper.Stop()
		}

		// 创建新限流器
		var limiter Limiter
		var err error

		switch m.config.Strategy {
		case "token_bucket":
			limiter, err = NewTokenBucketLimiter(m.config)
		case "sliding_window":
			limiter, err = NewSlidingWindowLimiter(m.config)
		case "redis":
			limiter, err = NewRedisLimiter(m.config)
		default:
			m.config = &oldConfig
			return fmt.Errorf("unsupported strategy: %s", m.config.Strategy)
		}

		if err != nil {
			m.config = &oldConfig
			return fmt.Errorf("failed to create new limiter: %w", err)
		}

		m.limiter = limiter
	}

	m.logger.Info("Rate limit config reloaded",
		zap.String("strategy", m.config.Strategy),
		zap.Int("rate", m.config.Rate),
		zap.Int("burst", m.config.Burst),
	)

	return nil
}

// GetStats 获取统计信息
func (m *RateLimitMiddleware) GetStats(key string) *LimiterStats {
	return m.limiter.GetStats(key)
}

// GetTotalStats 获取总统计信息
func (m *RateLimitMiddleware) GetTotalStats() *LimiterStats {
	if totalLimiter, ok := m.limiter.(interface{ GetTotalStats() *LimiterStats }); ok {
		return totalLimiter.GetTotalStats()
	}
	return nil
}

// 确保实现了核心接口
var _ core.Middleware = (*RateLimitMiddleware)(nil)
var _ core.ConfigurableMiddleware = (*RateLimitMiddleware)(nil)
var _ core.HotReloadMiddleware = (*RateLimitMiddleware)(nil)
