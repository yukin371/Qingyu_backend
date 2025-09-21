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
	// 默认提供商
	DefaultProvider string `json:"defaultProvider"`
	
	// 提供商配置映射
	Providers map[string]*ProviderConfig `json:"providers"`
	
	// 全局请求超时时间
	Timeout time.Duration `json:"timeout"`
	
	// 全局最大重试次数
	MaxRetries int `json:"maxRetries"`
	
	// 全局重试间隔
	RetryInterval time.Duration `json:"retryInterval"`
	
	// 全局代理设置
	ProxyURL string `json:"proxyUrl,omitempty"`
}

// ProviderConfig 单个提供商配置
type ProviderConfig struct {
	// 提供商名称 (openai, claude, gemini, wenxin, qwen)
	Name string `json:"name"`
	
	// API密钥
	APIKey string `json:"apiKey"`
	
	// API密钥2 (用于百度文心一言的SecretKey)
	SecretKey string `json:"secretKey,omitempty"`
	
	// API基础URL
	BaseURL string `json:"baseUrl"`
	
	// 默认模型
	DefaultModel string `json:"defaultModel"`
	
	// 支持的模型列表
	SupportedModels []string `json:"supportedModels"`
	
	// 是否启用
	Enabled bool `json:"enabled"`
	
	// 优先级 (数字越小优先级越高)
	Priority int `json:"priority"`
	
	// 提供商特定的超时时间
	Timeout time.Duration `json:"timeout,omitempty"`
	
	// 提供商特定的最大重试次数
	MaxRetries int `json:"maxRetries,omitempty"`
	
	// 提供商特定的重试间隔
	RetryInterval time.Duration `json:"retryInterval,omitempty"`
	
	// 提供商特定的代理设置
	ProxyURL string `json:"proxyUrl,omitempty"`
	
	// 额外配置参数
	ExtraConfig map[string]interface{} `json:"extraConfig,omitempty"`
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
			DefaultProvider: getEnv("AI_DEFAULT_PROVIDER", "openai"),
			Providers: map[string]*ProviderConfig{
				"openai": {
					Name:            "openai",
					APIKey:          getEnv("OPENAI_API_KEY", ""),
					BaseURL:         getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
					DefaultModel:    getEnv("OPENAI_DEFAULT_MODEL", "gpt-3.5-turbo"),
					SupportedModels: []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo", "gpt-3.5-turbo-16k"},
					Enabled:         getEnvAsBool("OPENAI_ENABLED", true),
					Priority:        getEnvAsInt("OPENAI_PRIORITY", 1),
				},
				"claude": {
					Name:            "claude",
					APIKey:          getEnv("CLAUDE_API_KEY", ""),
					BaseURL:         getEnv("CLAUDE_BASE_URL", "https://api.anthropic.com"),
					DefaultModel:    getEnv("CLAUDE_DEFAULT_MODEL", "claude-3-5-sonnet-20241022"),
					SupportedModels: []string{"claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022", "claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"},
					Enabled:         getEnvAsBool("CLAUDE_ENABLED", false),
					Priority:        getEnvAsInt("CLAUDE_PRIORITY", 2),
				},
				"gemini": {
					Name:            "gemini",
					APIKey:          getEnv("GEMINI_API_KEY", ""),
					BaseURL:         getEnv("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com"),
					DefaultModel:    getEnv("GEMINI_DEFAULT_MODEL", "gemini-1.5-pro"),
					SupportedModels: []string{"gemini-1.5-pro", "gemini-1.5-flash", "gemini-1.0-pro", "gemini-pro-vision"},
					Enabled:         getEnvAsBool("GEMINI_ENABLED", false),
					Priority:        getEnvAsInt("GEMINI_PRIORITY", 3),
				},
				"wenxin": {
					Name:            "wenxin",
					APIKey:          getEnv("WENXIN_API_KEY", ""),
					SecretKey:       getEnv("WENXIN_SECRET_KEY", ""),
					BaseURL:         getEnv("WENXIN_BASE_URL", "https://aip.baidubce.com"),
					DefaultModel:    getEnv("WENXIN_DEFAULT_MODEL", "ernie-3.5-8k"),
					SupportedModels: []string{"ernie-4.0-8k", "ernie-4.0-8k-preview", "ernie-3.5-8k", "ernie-3.5-8k-0205", "ernie-turbo-8k", "ernie-speed-8k", "ernie-lite-8k", "ernie-tiny-8k"},
					Enabled:         getEnvAsBool("WENXIN_ENABLED", false),
					Priority:        getEnvAsInt("WENXIN_PRIORITY", 4),
				},
				"qwen": {
					Name:            "qwen",
					APIKey:          getEnv("QWEN_API_KEY", ""),
					BaseURL:         getEnv("QWEN_BASE_URL", "https://dashscope.aliyuncs.com"),
					DefaultModel:    getEnv("QWEN_DEFAULT_MODEL", "qwen-turbo"),
					SupportedModels: []string{"qwen-turbo", "qwen-plus", "qwen-max", "qwen-max-1201", "qwen-max-longcontext", "qwen-7b-chat", "qwen-14b-chat", "qwen-72b-chat", "qwen1.5-7b-chat", "qwen1.5-14b-chat", "qwen1.5-72b-chat", "qwen2-7b-instruct", "qwen2-72b-instruct"},
					Enabled:         getEnvAsBool("QWEN_ENABLED", false),
					Priority:        getEnvAsInt("QWEN_PRIORITY", 5),
				},
			},
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