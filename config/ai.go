package config

import "time"

// AIConfig AI服务配置
type AIConfig struct {
	// 外部API配置
	ExternalAPI *ExternalAPIConfig `json:"externalApi"`
	
	// 上下文配置
	Context *ContextConfig `json:"context"`
	
	// 缓存配置
	Cache *CacheConfig `json:"cache"`
	
	// 限流配置
	RateLimit *RateLimitConfig `json:"rateLimit"`
}

// ExternalAPIConfig 外部API配置
type ExternalAPIConfig struct {
	// API提供商 (openai, claude, gemini, etc.)
	Provider string `json:"provider"`
	
	// API密钥
	APIKey string `json:"apiKey"`
	
	// API基础URL
	BaseURL string `json:"baseUrl"`
	
	// 默认模型
	DefaultModel string `json:"defaultModel"`
	
	// 请求超时时间
	Timeout time.Duration `json:"timeout"`
	
	// 最大重试次数
	MaxRetries int `json:"maxRetries"`
	
	// 重试间隔
	RetryInterval time.Duration `json:"retryInterval"`
	
	// 代理设置
	ProxyURL string `json:"proxyUrl,omitempty"`
}

// ContextConfig 上下文配置
type ContextConfig struct {
	// 最大token数量
	MaxTokens int `json:"maxTokens"`
	
	// 默认历史深度
	DefaultHistoryDepth int `json:"defaultHistoryDepth"`
	
	// 默认大纲深度
	DefaultOutlineDepth int `json:"defaultOutlineDepth"`
	
	// 是否默认包含历史
	DefaultIncludeHistory bool `json:"defaultIncludeHistory"`
	
	// 是否默认包含大纲
	DefaultIncludeOutline bool `json:"defaultIncludeOutline"`
	
	// 上下文缓存时间
	CacheExpiration time.Duration `json:"cacheExpiration"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 是否启用缓存
	Enabled bool `json:"enabled"`
	
	// 缓存类型 (memory, redis)
	Type string `json:"type"`
	
	// Redis配置（当type为redis时）
	RedisURL string `json:"redisUrl,omitempty"`
	
	// 默认过期时间
	DefaultExpiration time.Duration `json:"defaultExpiration"`
	
	// 清理间隔
	CleanupInterval time.Duration `json:"cleanupInterval"`
	
	// 最大缓存大小（内存缓存）
	MaxSize int `json:"maxSize"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 是否启用限流
	Enabled bool `json:"enabled"`
	
	// 每分钟请求数限制
	RequestsPerMinute int `json:"requestsPerMinute"`
	
	// 每小时请求数限制
	RequestsPerHour int `json:"requestsPerHour"`
	
	// 每天请求数限制
	RequestsPerDay int `json:"requestsPerDay"`
	
	// 突发请求数量
	BurstSize int `json:"burstSize"`
}

// LoadAIConfig 加载AI配置
func LoadAIConfig() *AIConfig {
	return &AIConfig{
		ExternalAPI: &ExternalAPIConfig{
			Provider:      getEnv("AI_PROVIDER", "openai"),
			APIKey:        getEnv("AI_API_KEY", ""),
			BaseURL:       getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
			DefaultModel:  getEnv("AI_DEFAULT_MODEL", "gpt-3.5-turbo"),
			Timeout:       time.Duration(getEnvAsInt("AI_TIMEOUT", 30)) * time.Second,
			MaxRetries:    getEnvAsInt("AI_MAX_RETRIES", 3),
			RetryInterval: time.Duration(getEnvAsInt("AI_RETRY_INTERVAL", 1)) * time.Second,
			ProxyURL:      getEnv("AI_PROXY_URL", ""),
		},
		Context: &ContextConfig{
			MaxTokens:             getEnvAsInt("AI_MAX_TOKENS", 4000),
			DefaultHistoryDepth:   getEnvAsInt("AI_DEFAULT_HISTORY_DEPTH", 3),
			DefaultOutlineDepth:   getEnvAsInt("AI_DEFAULT_OUTLINE_DEPTH", 2),
			DefaultIncludeHistory: getEnvAsBool("AI_DEFAULT_INCLUDE_HISTORY", true),
			DefaultIncludeOutline: getEnvAsBool("AI_DEFAULT_INCLUDE_OUTLINE", true),
			CacheExpiration:       time.Duration(getEnvAsInt("AI_CONTEXT_CACHE_EXPIRATION", 300)) * time.Second,
		},
		Cache: &CacheConfig{
			Enabled:           getEnvAsBool("AI_CACHE_ENABLED", true),
			Type:              getEnv("AI_CACHE_TYPE", "memory"),
			RedisURL:          getEnv("AI_CACHE_REDIS_URL", ""),
			DefaultExpiration: time.Duration(getEnvAsInt("AI_CACHE_DEFAULT_EXPIRATION", 3600)) * time.Second,
			CleanupInterval:   time.Duration(getEnvAsInt("AI_CACHE_CLEANUP_INTERVAL", 600)) * time.Second,
			MaxSize:           getEnvAsInt("AI_CACHE_MAX_SIZE", 1000),
		},
		RateLimit: &RateLimitConfig{
			Enabled:           getEnvAsBool("AI_RATE_LIMIT_ENABLED", true),
			RequestsPerMinute: getEnvAsInt("AI_RATE_LIMIT_PER_MINUTE", 60),
			RequestsPerHour:   getEnvAsInt("AI_RATE_LIMIT_PER_HOUR", 1000),
			RequestsPerDay:    getEnvAsInt("AI_RATE_LIMIT_PER_DAY", 10000),
			BurstSize:         getEnvAsInt("AI_RATE_LIMIT_BURST_SIZE", 10),
		},
	}
}