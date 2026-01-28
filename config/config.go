package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	authModel "Qingyu_backend/models/auth"
)

// Config 存储应用配置
type Config struct {
	Database *DatabaseConfig                   `mapstructure:"database"`
	Redis    *RedisConfig                      `mapstructure:"redis"`
	Cache    *CacheConfig                      `mapstructure:"cache"`
	Server   *ServerConfig                     `mapstructure:"server"`
	Log      *LogConfig                        `mapstructure:"log"`
	JWT      *JWTConfig                        `mapstructure:"jwt"`
	AI       *AIConfig                         `mapstructure:"ai"`
	External *ExternalAPIConfig                `mapstructure:"external"`
	AIQuota  *AIQuotaConfig                    `mapstructure:"ai_quota"`
	Email    *EmailConfig                      `mapstructure:"email"`
	Payment  *PaymentConfig                    `mapstructure:"payment"`
	RateLimit *RateLimitConfig                 `mapstructure:"rate_limit"`
	OAuth    map[string]*authModel.OAuthConfig `mapstructure:"oauth"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level       string            `mapstructure:"level"`
	Format      string            `mapstructure:"format"`
	Output      string            `mapstructure:"output"`
	Filename    string            `mapstructure:"filename"`
	Development bool              `mapstructure:"development"`
	Mode        string            `mapstructure:"mode"` // normal|strict
	Request     *LogRequestConfig `mapstructure:"request"`
	RedactKeys  []string          `mapstructure:"redact_keys"`
}

// LogRequestConfig 请求日志配置
type LogRequestConfig struct {
	EnableBody     bool     `mapstructure:"enable_body"`
	MaxBodySize    int      `mapstructure:"max_body_size"`
	SkipPaths      []string `mapstructure:"skip_paths"`
	BodyAllowPaths []string `mapstructure:"body_allow_paths"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

// AIConfig AI配置
type AIConfig struct {
	APIKey        string                 `mapstructure:"api_key"`
	BaseURL       string                 `mapstructure:"base_url"`
	MaxTokens     int                    `mapstructure:"max_tokens"`
	Temperature   int                    `mapstructure:"temperature"`
	PythonService *PythonAIServiceConfig `mapstructure:"python_service"`
	AIService     *AIServiceConfig       `mapstructure:"ai_service"`
}

// AIServiceConfig AI服务配置（gRPC）
type AIServiceConfig struct {
	Endpoint       string `mapstructure:"endpoint"`        // gRPC服务端点
	EnableFallback bool   `mapstructure:"enable_fallback"` // 是否启用适配器回退
	Timeout        int    `mapstructure:"timeout"`         // 请求超时（秒）
}

// PythonAIServiceConfig Python AI服务配置
type PythonAIServiceConfig struct {
	Host           string `mapstructure:"host"`
	GrpcPort       int    `mapstructure:"grpc_port"`
	EmbeddingModel string `mapstructure:"embedding_model"`
	MilvusHost     string `mapstructure:"milvus_host"`
	MilvusPort     int    `mapstructure:"milvus_port"`
	RedisHost      string `mapstructure:"redis_host"`
	RedisPort      int    `mapstructure:"redis_port"`
}

// ExternalAPIConfig 外部API配置（用于适配器管理器）
type ExternalAPIConfig struct {
	DefaultProvider string                     `mapstructure:"default_provider"`
	Providers       map[string]*ProviderConfig `mapstructure:"providers"`
}

// ProviderConfig 提供商配置
type ProviderConfig struct {
	Name            string   `mapstructure:"name"`
	APIKey          string   `mapstructure:"api_key"`
	SecretKey       string   `mapstructure:"secret_key"`
	BaseURL         string   `mapstructure:"base_url"`
	Priority        int      `mapstructure:"priority"`
	Enabled         bool     `mapstructure:"enabled"`
	SupportedModels []string `mapstructure:"supported_models"`
}

// AIQuotaConfig AI配额配置
type AIQuotaConfig struct {
	DefaultQuotas    *DefaultQuotasConfig `mapstructure:"default_quotas"`
	Reset            *QuotaResetConfig    `mapstructure:"reset"`
	WarningThreshold float64              `mapstructure:"warning_threshold"`
	AllowOverdraft   bool                 `mapstructure:"allow_overdraft"`
	OverdraftLimit   int                  `mapstructure:"overdraft_limit"`
}

// DefaultQuotasConfig 默认配额配置
type DefaultQuotasConfig struct {
	Reader map[string]int `mapstructure:"reader"` // normal, vip
	Writer map[string]int `mapstructure:"writer"` // novice, signed, master
	Admin  map[string]int `mapstructure:"admin"`  // normal
}

// QuotaResetConfig 配额重置配置
type QuotaResetConfig struct {
	DailyResetHour  int  `mapstructure:"daily_reset_hour"`
	EnableAutoReset bool `mapstructure:"enable_auto_reset"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	SMTPHost    string `mapstructure:"smtp_host"`
	SMTPPort    int    `mapstructure:"smtp_port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	FromAddress string `mapstructure:"from_address"`
	FromName    string `mapstructure:"from_name"`
	UseTLS      bool   `mapstructure:"use_tls"`
	UseSSL      bool   `mapstructure:"use_ssl"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	Enabled         bool             `mapstructure:"enabled"`
	DefaultProvider string           `mapstructure:"default_provider"` // alipay, wechat
	Alipay          *AlipayConfig    `mapstructure:"alipay"`
	Wechat          *WechatPayConfig `mapstructure:"wechat"`
	NotifyURL       string           `mapstructure:"notify_url"`
	ReturnURL       string           `mapstructure:"return_url"`
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	Enabled        bool     `mapstructure:"enabled" json:"enabled"`
	RequestsPerSec float64  `mapstructure:"requests_per_sec" json:"requests_per_sec"`
	Burst          int      `mapstructure:"burst" json:"burst"`
	SkipPaths      []string `mapstructure:"skip_paths" json:"skip_paths"`
}

// Validate validates the rate limit configuration
func (c *RateLimitConfig) Validate() error {
	// Skip validation if disabled
	if !c.Enabled {
		return nil
	}

	if c.RequestsPerSec <= 0 {
		return fmt.Errorf("rate_limit.requests_per_sec must be positive, got %f", c.RequestsPerSec)
	}
	if c.Burst < 0 {
		return fmt.Errorf("rate_limit.burst cannot be negative, got %d", c.Burst)
	}
	if c.Burst > 10000 {
		return fmt.Errorf("rate_limit.burst too large (%d), maximum is 10000", c.Burst)
	}
	return nil
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	AppID           string `mapstructure:"app_id"`
	PrivateKey      string `mapstructure:"private_key"`
	PublicKey       string `mapstructure:"public_key"`
	AlipayPublicKey string `mapstructure:"alipay_public_key"`
	SignType        string `mapstructure:"sign_type"` // RSA2
	Sandbox         bool   `mapstructure:"sandbox"`
}

// WechatPayConfig 微信支付配置
type WechatPayConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	AppID      string `mapstructure:"app_id"`
	MchID      string `mapstructure:"mch_id"`
	APIKey     string `mapstructure:"api_key"`
	APIv3Key   string `mapstructure:"apiv3_key"`
	SerialNo   string `mapstructure:"serial_no"`
	PrivateKey string `mapstructure:"private_key"`
	Sandbox    bool   `mapstructure:"sandbox"`
}

// GetDefaultQuota 从配置获取默认配额
func (c *AIQuotaConfig) GetDefaultQuota(userRole, membershipLevel string) int {
	if c == nil || c.DefaultQuotas == nil {
		// 如果配置不存在，返回硬编码的默认值
		return 5
	}

	switch userRole {
	case "reader":
		if quota, ok := c.DefaultQuotas.Reader[membershipLevel]; ok {
			return quota
		}
		// 如果没有找到对应等级，尝试normal
		if quota, ok := c.DefaultQuotas.Reader["normal"]; ok {
			return quota
		}
	case "writer":
		if quota, ok := c.DefaultQuotas.Writer[membershipLevel]; ok {
			return quota
		}
		if quota, ok := c.DefaultQuotas.Writer["novice"]; ok {
			return quota
		}
	case "admin":
		if quota, ok := c.DefaultQuotas.Admin["normal"]; ok {
			return quota
		}
	}

	// 最后的默认值
	return 5
}

// DefaultRateLimitConfig returns a default rate limit configuration
// Note: Actual defaults are set in setDefaults() for Viper
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Enabled: true,
		// Other fields will be populated by Viper defaults
	}
}

var (
	// GlobalConfig 全局配置实例
	GlobalConfig *Config
	// v viper实例
	v *viper.Viper
)

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	v = viper.New()

	// 设置默认值
	setDefaults()

	// 检查是否是文件路径
	if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
		// 直接设置配置文件
		v.SetConfigFile(configPath)
	} else {
		// 优先尝试加载 config.test.yaml（测试配置优先）
		testConfigPaths := []string{
			configPath,
			"./config",
			"../../config",
			".",
		}

		testConfigFound := false
		for _, path := range testConfigPaths {
			testConfigFile := fmt.Sprintf("%s/config.test.yaml", path)
			if _, err := os.Stat(testConfigFile); err == nil {
				v.SetConfigFile(testConfigFile)
				testConfigFound = true
				fmt.Printf("[Config] Using test config: %s\n", testConfigFile)
				break
			}
		}

		// 如果没有找到 test 配置，使用默认配置查找
		if !testConfigFound {
			v.SetConfigName("config")       // 配置文件名（不带扩展名）
			v.SetConfigType("yaml")         // 配置文件类型
			v.AddConfigPath(configPath)     // 配置文件路径
			v.AddConfigPath("./config")     // config子目录
			v.AddConfigPath("../../config") // 从cmd/server运行时
			v.AddConfigPath(".")            // 当前目录
			fmt.Println("[Config] Test config not found, using default config search")
		}
	}

	v.AutomaticEnv()                                   // 读取环境变量
	v.SetEnvPrefix("QINGYU")                           // 环境变量前缀
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 环境变量键替换器

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果有默认值，可以继续
			fmt.Printf("Warning: Config file not found: %v\n", err)
		} else {
			// 其他错误
			return nil, fmt.Errorf("fatal error reading config file: %w", err)
		}
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// 设置默认值
	setConfigDefaults(config)

	// 验证配置
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// 设置全局配置实例
	GlobalConfig = config

	return config, nil
}

// LoadDatabaseConfigFromViper 从viper加载数据库配置
func LoadDatabaseConfigFromViper() (*DatabaseConfig, error) {
	if v == nil {
		return nil, fmt.Errorf("viper未初始化，请先调用LoadConfig")
	}

	var dbConfig DatabaseConfig
	if err := v.UnmarshalKey("database", &dbConfig); err != nil {
		return nil, fmt.Errorf("解析数据库配置失败: %w", err)
	}

	// 验证配置
	if err := dbConfig.Validate(); err != nil {
		return nil, fmt.Errorf("数据库配置验证失败: %w", err)
	}

	return &dbConfig, nil
}

// setDefaults 设置Viper默认值
func setDefaults() {
	// 数据库默认配置 - 使用database.go中的高级配置结构
	v.SetDefault("database.primary.type", "mongodb")
	v.SetDefault("database.primary.mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("database.primary.mongodb.database", "qingyu")
	v.SetDefault("database.primary.mongodb.max_pool_size", 100)
	v.SetDefault("database.primary.mongodb.min_pool_size", 5)
	v.SetDefault("database.primary.mongodb.connect_timeout", 10*time.Second)
	v.SetDefault("database.primary.mongodb.server_timeout", 30*time.Second)
	v.SetDefault("database.primary.mongodb.profiling_level", 1)
	v.SetDefault("database.primary.mongodb.slow_ms", 100)
	v.SetDefault("database.primary.mongodb.profiler_size_mb", 100)

	// 索引配置
	v.SetDefault("database.indexing.auto_create", true)
	v.SetDefault("database.indexing.background", true)

	// 验证配置
	v.SetDefault("database.validation.enabled", true)
	v.SetDefault("database.validation.strict_mode", false)

	// 迁移配置
	v.SetDefault("database.migration.enabled", false)
	v.SetDefault("database.migration.batch_size", 1000)
	v.SetDefault("database.migration.timeout", 30*time.Minute)
	v.SetDefault("database.migration.strategy.type", "full")
	v.SetDefault("database.migration.strategy.verify_data", true)
	v.SetDefault("database.migration.strategy.rollback_on_fail", true)

	// 同步配置
	v.SetDefault("database.sync.enabled", false)
	v.SetDefault("database.sync.interval", 5*time.Minute)
	v.SetDefault("database.sync.batch_size", 100)
	v.SetDefault("database.sync.direction", "source_to_target")
	v.SetDefault("database.sync.conflict_resolution", "latest_wins")

	// 服务器默认配置
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")

	// 日志默认配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.filename", "logs/app.log")
	v.SetDefault("log.development", false)
	v.SetDefault("log.mode", "normal")
	v.SetDefault("log.request.enable_body", false)
	v.SetDefault("log.request.max_body_size", 2048)
	v.SetDefault("log.request.skip_paths", []string{"/health", "/metrics", "/swagger"})
	v.SetDefault("log.request.body_allow_paths", []string{})
	v.SetDefault("log.redact_keys", []string{"authorization", "password", "token", "cookie"})

	// JWT默认配置
	v.SetDefault("jwt.secret", "qingyu_secret_key")
	v.SetDefault("jwt.expiration_hours", 24)

	// AI默认配置
	v.SetDefault("ai.api_key", "default_api_key")
	v.SetDefault("ai.base_url", "https://api.openai.com/v1")
	v.SetDefault("ai.max_tokens", 2000)
	v.SetDefault("ai.temperature", 7)

	// Python AI服务默认配置
	v.SetDefault("ai.python_service.host", "localhost")
	v.SetDefault("ai.python_service.grpc_port", 50052)
	v.SetDefault("ai.python_service.embedding_model", "BAAI/bge-large-zh-v1.5")
	v.SetDefault("ai.python_service.milvus_host", "localhost")
	v.SetDefault("ai.python_service.milvus_port", 19530)
	v.SetDefault("ai.python_service.redis_host", "localhost")
	v.SetDefault("ai.python_service.redis_port", 6379)

	// AI服务默认配置（gRPC）
	v.SetDefault("ai.ai_service.endpoint", "localhost:50052")
	v.SetDefault("ai.ai_service.enable_fallback", false)
	v.SetDefault("ai.ai_service.timeout", 30)

	// 邮件默认配置
	v.SetDefault("email.enabled", false)
	v.SetDefault("email.smtp_host", "smtp.example.com")
	v.SetDefault("email.smtp_port", 587)
	v.SetDefault("email.username", "")
	v.SetDefault("email.password", "")
	v.SetDefault("email.from_address", "noreply@qingyu.com")
	v.SetDefault("email.from_name", "青羽阅读")
	v.SetDefault("email.use_tls", true)
	v.SetDefault("email.use_ssl", false)

	// 支付默认配置
	v.SetDefault("payment.enabled", false)
	v.SetDefault("payment.default_provider", "alipay")
	v.SetDefault("payment.notify_url", "")
	v.SetDefault("payment.return_url", "")

	// 支付宝默认配置
	v.SetDefault("payment.alipay.enabled", false)
	v.SetDefault("payment.alipay.sandbox", true)
	v.SetDefault("payment.alipay.sign_type", "RSA2")

	// 微信支付默认配置
	v.SetDefault("payment.wechat.enabled", false)
	v.SetDefault("payment.wechat.sandbox", true)

	// 缓存默认配置
	v.SetDefault("cache.enabled", false)
	v.SetDefault("cache.double_delete_delay", 1*time.Second)
	v.SetDefault("cache.null_cache_ttl", 30*time.Second)
	v.SetDefault("cache.null_cache_prefix", "@@NULL@@")
	v.SetDefault("cache.breaker_max_requests", 3)
	v.SetDefault("cache.breaker_interval", 10*time.Second)
	v.SetDefault("cache.breaker_timeout", 30*time.Second)
	v.SetDefault("cache.breaker_threshold", 0.6)

	// 速率限制默认配置
	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.requests_per_sec", 100)
	v.SetDefault("rate_limit.burst", 200)
	v.SetDefault("rate_limit.skip_paths", []string{"/health", "/metrics"})
}

// WatchConfig 启用配置热重载
func WatchConfig(onChange func()) {
	v.WatchConfig()
	if onChange != nil {
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("Config file changed: %s\n", e.Name)

			// 重新加载配置
			config := &Config{}
			if err := v.Unmarshal(config); err != nil {
				fmt.Printf("Error reloading config: %v\n", err)
				return
			}

			// 设置默认值
			setConfigDefaults(config)

			// 验证配置
			if err := ValidateConfig(config); err != nil {
				fmt.Printf("Error validating reloaded config: %v\n", err)
				return
			}

			// 更新全局配置
			GlobalConfig = config

			// 调用回调函数
			onChange()
		})
	}
}

// GetViper 获取viper实例
func GetViper() *viper.Viper {
	return v
}
