package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 存储应用配置
type Config struct {
	Database *DatabaseConfig    `mapstructure:"database"`
	Redis    *RedisConfig       `mapstructure:"redis"`
	Server   *ServerConfig      `mapstructure:"server"`
	JWT      *JWTConfig         `mapstructure:"jwt"`
	AI       *AIConfig          `mapstructure:"ai"`
	External *ExternalAPIConfig `mapstructure:"external"`
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

	// JWT默认配置
	v.SetDefault("jwt.secret", "qingyu_secret_key")
	v.SetDefault("jwt.expiration_hours", 24)

	// AI默认配置
	v.SetDefault("ai.api_key", "default_api_key")
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
