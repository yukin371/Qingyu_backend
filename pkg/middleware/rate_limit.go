package middleware

import (
	"fmt"
	"sync"
	"time"

	"Qingyu_backend/pkg/errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	// 每秒允许的请求数
	Rate float64
	// 桶容量（突发流量）
	Burst int
	// 限流键生成函数
	KeyFunc func(*gin.Context) string
	// 是否使用Redis分布式限流
	UseRedis bool
	// Redis配置（可选）
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// DefaultRateLimiterConfig 默认限流配置
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Rate:    10,  // 每秒10个请求
		Burst:   100, // 桶容量100
		KeyFunc: DefaultKeyFunc,
	}
}

// DefaultKeyFunc 默认的限流键生成函数（使用IP）
func DefaultKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// UserKeyFunc 基于用户ID的限流键生成函数
func UserKeyFunc(c *gin.Context) string {
	userID := c.GetString("userId")
	if userID == "" {
		return c.ClientIP()
	}
	return fmt.Sprintf("user:%s", userID)
}

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	limiters map[string]*limiter
	mu       sync.RWMutex
	config   *RateLimiterConfig
}

type limiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter(config *RateLimiterConfig) *IPRateLimiter {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	return &IPRateLimiter{
		limiters: make(map[string]*limiter),
		config:   config,
	}
}

// GetLimiter 获取限流器
func (rl *IPRateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 清理过期的限流器（5分钟未使用）
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// 检查是否已存在
	if l, exists := rl.limiters[key]; exists {
		l.lastSeen = time.Now()
		return l.limiter
	}

	// 创建新的限流器
	l := &limiter{
		limiter:  rate.NewLimiter(rate.Limit(rl.config.Rate), rl.config.Burst),
		lastSeen: time.Now(),
	}
	rl.limiters[key] = l

	return l.limiter
}

// CleanupLimiters 清理过期的限流器
func (rl *IPRateLimiter) CleanupLimiters() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, l := range rl.limiters {
		if now.Sub(l.lastSeen) > 5*time.Minute {
			delete(rl.limiters, key)
		}
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config *RateLimiterConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	limiter := NewIPRateLimiter(config)

	// 启动清理协程
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.CleanupLimiters()
		}
	}()

	return func(c *gin.Context) {
		// 生成限流键
		key := config.KeyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		// 获取限流器
		l := limiter.GetLimiter(key)

		// 检查是否允许请求
		if !l.Allow() {
			ErrorResponse(c, errors.NewRateLimit())
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdvancedRateLimiter 高级限流器（支持不同的限流策略）
type AdvancedRateLimiter struct {
	globalLimiter *rate.Limiter
	ipLimiters    map[string]*rate.Limiter
	userLimiters  map[string]*rate.Limiter
	mu            sync.RWMutex
	config        *RateLimiterConfig
}

// NewAdvancedRateLimiter 创建高级限流器
func NewAdvancedRateLimiter(config *RateLimiterConfig) *AdvancedRateLimiter {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	return &AdvancedRateLimiter{
		globalLimiter: rate.NewLimiter(rate.Limit(config.Rate)*10, config.Burst*10), // 全局限制
		ipLimiters:    make(map[string]*rate.Limiter),
		userLimiters:  make(map[string]*rate.Limiter),
		config:        config,
	}
}

// Allow 检查是否允许请求
func (rl *AdvancedRateLimiter) Allow(c *gin.Context) bool {
	// 首先检查全局限制
	if !rl.globalLimiter.Allow() {
		return false
	}

	// 检查IP限制
	ip := c.ClientIP()
	rl.mu.RLock()
	ipLimiter, ipExists := rl.ipLimiters[ip]
	rl.mu.RUnlock()

	if !ipExists {
		rl.mu.Lock()
		ipLimiter = rate.NewLimiter(rate.Limit(rl.config.Rate), rl.config.Burst)
		rl.ipLimiters[ip] = ipLimiter
		rl.mu.Unlock()
	}

	if !ipLimiter.Allow() {
		return false
	}

	// 检查用户限制（如果已登录）
	userID := c.GetString("userId")
	if userID != "" {
		rl.mu.RLock()
		userLimiter, userExists := rl.userLimiters[userID]
		rl.mu.RUnlock()

		if !userExists {
			rl.mu.Lock()
			// 用户限制可以比IP限制更宽松
			userLimiter = rate.NewLimiter(rate.Limit(rl.config.Rate)*2, rl.config.Burst*2)
			rl.userLimiters[userID] = userLimiter
			rl.mu.Unlock()
		}

		if !userLimiter.Allow() {
			return false
		}
	}

	return true
}

// AdvancedRateLimitMiddleware 高级限流中间件
func AdvancedRateLimitMiddleware(config *RateLimiterConfig) gin.HandlerFunc {
	limiter := NewAdvancedRateLimiter(config)

	return func(c *gin.Context) {
		if !limiter.Allow(c) {
			ErrorResponse(c, errors.NewRateLimit())
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitByPath 按路径限流
// 不同路径可以有不同的限流策略
type RateLimitByPath struct {
	limiters map[string]*rate.Limiter
	configs  map[string]*RateLimiterConfig
	mu       sync.RWMutex
}

// NewRateLimitByPath 创建按路径限流的限流器
func NewRateLimitByPath() *RateLimitByPath {
	return &RateLimitByPath{
		limiters: make(map[string]*rate.Limiter),
		configs:  make(map[string]*RateLimiterConfig),
	}
}

// AddPath 添加路径限流规则
func (rl *RateLimitByPath) AddPath(path string, config *RateLimiterConfig) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.configs[path] = config
	rl.limiters[path] = rate.NewLimiter(rate.Limit(config.Rate), config.Burst)
}

// GetLimiter 获取路径对应的限流器
func (rl *RateLimitByPath) GetLimiter(path string) *rate.Limiter {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	// 精确匹配
	if l, exists := rl.limiters[path]; exists {
		return l
	}

	// 返回默认限流器
	return rate.NewLimiter(10, 100)
}

// RateLimitByPathMiddleware 按路径限流中间件
func RateLimitByPathMiddleware(rl *RateLimitByPath) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		limiter := rl.GetLimiter(path)

		if !limiter.Allow() {
			ErrorResponse(c, errors.NewRateLimit())
			c.Abort()
			return
		}

		c.Next()
	}
}

// SlidingWindowLog 滑动窗口日志限流器
type SlidingWindowLog struct {
	windows map[string][]time.Time
	mu      sync.RWMutex
	limit   int
	window  time.Duration
}

// NewSlidingWindowLog 创建滑动窗口限流器
func NewSlidingWindowLog(limit int, window time.Duration) *SlidingWindowLog {
	return &SlidingWindowLog{
		windows: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

// Allow 检查是否允许请求
func (sw *SlidingWindowLog) Allow(key string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-sw.window)

	// 获取窗口内的请求记录
	timestamps, exists := sw.windows[key]
	if !exists {
		timestamps = []time.Time{}
		sw.windows[key] = timestamps
	}

	// 清理窗口外的记录
	validTimestamps := []time.Time{}
	for _, ts := range timestamps {
		if ts.After(windowStart) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// 检查是否超过限制
	if len(validTimestamps) >= sw.limit {
		return false
	}

	// 添加当前请求
	sw.windows[key] = append(validTimestamps, now)
	return true
}

// SlidingWindowMiddleware 滑动窗口限流中间件
func SlidingWindowMiddleware(limit int, window time.Duration, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	if keyFunc == nil {
		keyFunc = DefaultKeyFunc
	}

	sw := NewSlidingWindowLog(limit, window)

	return func(c *gin.Context) {
		key := keyFunc(c)
		if !sw.Allow(key) {
			ErrorResponse(c, errors.NewRateLimit())
			c.Abort()
			return
		}
		c.Next()
	}
}
