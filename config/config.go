package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 存储应用配置
type Config struct {
	Database *DatabaseConfig `mapstructure:"database"`
	Server   *ServerConfig   `mapstructure:"server"`
	JWT      *JWTConfig     `mapstructure:"jwt"`
	AI       *AIConfig      `mapstructure:"ai"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MongoURI        string        `mapstructure:"uri"`
	DBName          string        `mapstructure:"name"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"`
	MaxPoolSize     uint64        `mapstructure:"max_pool_size"`
	MinPoolSize     uint64        `mapstructure:"min_pool_size"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
	RetryWrites     bool          `mapstructure:"retry_writes"`
	RetryReads      bool          `mapstructure:"retry_reads"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

// AIConfig AI配置
type AIConfig struct {
	APIKey      string `mapstructure:"api_key"`
	BaseURL     string `mapstructure:"base_url"`
	MaxTokens   int    `mapstructure:"max_tokens"`
	Temperature int    `mapstructure:"temperature"`
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

	// 配置Viper
	v.SetConfigName("config")    // 配置文件名（不带扩展名）
	v.SetConfigType("yaml")      // 配置文件类型
	v.AddConfigPath(configPath)  // 配置文件路径
	v.AddConfigPath(".")         // 当前目录
	v.AutomaticEnv()            // 读取环境变量
	v.SetEnvPrefix("QINGYU")    // 环境变量前缀

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

// setDefaults 设置Viper默认值
func setDefaults() {
	// 数据库默认配置
	v.SetDefault("database.uri", "mongodb://localhost:27017")
	v.SetDefault("database.name", "Qingyu_writer")
	v.SetDefault("database.connect_timeout", 10*time.Second)
	v.SetDefault("database.max_pool_size", 100)
	v.SetDefault("database.min_pool_size", 10)
	v.SetDefault("database.max_conn_idle_time", 60*time.Second)
	v.SetDefault("database.retry_writes", true)
	v.SetDefault("database.retry_reads", true)

	// 服务器默认配置
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")

	// JWT默认配置
	v.SetDefault("jwt.secret", "qingyu_secret_key")
	v.SetDefault("jwt.expiration_hours", 24)

	// AI默认配置
	v.SetDefault("ai.base_url", "https://api.openai.com/v1")
	v.SetDefault("ai.max_tokens", 2000)
	v.SetDefault("ai.temperature", 7)
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
