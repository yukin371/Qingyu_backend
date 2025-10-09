package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ViperConfigManager viper配置管理器
type ViperConfigManager struct {
	viper *viper.Viper
}

// NewViperConfigManager 创建viper配置管理器
func NewViperConfigManager() *ViperConfigManager {
	return &ViperConfigManager{
		viper: viper.New(),
	}
}

// LoadFromFile 从文件加载配置
func (m *ViperConfigManager) LoadFromFile(configPath string) error {
	m.viper.SetConfigName("config")
	m.viper.SetConfigType("yaml")
	m.viper.AddConfigPath(configPath)
	m.viper.AddConfigPath(".")

	// 设置环境变量支持
	m.viper.AutomaticEnv()
	m.viper.SetEnvPrefix("QINGYU")
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	m.setViperDefaults()

	return m.viper.ReadInConfig()
}

// LoadDatabaseConfig 加载数据库配置
func (m *ViperConfigManager) LoadDatabaseConfig() (*DatabaseConfig, error) {
	var config DatabaseConfig
	if err := m.viper.UnmarshalKey("database", &config); err != nil {
		return nil, fmt.Errorf("解析数据库配置失败: %w", err)
	}

	// 应用环境变量覆盖
	if err := m.applyDatabaseEnvOverrides(&config); err != nil {
		return nil, fmt.Errorf("应用环境变量覆盖失败: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("数据库配置验证失败: %w", err)
	}

	return &config, nil
}

// LoadAppConfig 加载应用配置
func (m *ViperConfigManager) LoadAppConfig() (*Config, error) {
	var config Config
	if err := m.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析应用配置失败: %w", err)
	}

	// 加载数据库配置
	dbConfig, err := m.LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("加载数据库配置失败: %w", err)
	}
	config.Database = dbConfig

	// 设置默认值
	setConfigDefaults(&config)

	// 验证配置
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("应用配置验证失败: %w", err)
	}

	return &config, nil
}

// setViperDefaults 设置viper默认值
func (m *ViperConfigManager) setViperDefaults() {
	// 数据库默认配置
	m.viper.SetDefault("database.primary.type", "mongodb")
	m.viper.SetDefault("database.primary.mongodb.uri", "mongodb://localhost:27017")
	m.viper.SetDefault("database.primary.mongodb.database", "qingyu")
	m.viper.SetDefault("database.primary.mongodb.max_pool_size", 100)
	m.viper.SetDefault("database.primary.mongodb.min_pool_size", 10)
	m.viper.SetDefault("database.primary.mongodb.connect_timeout", 10*time.Second)
	m.viper.SetDefault("database.primary.mongodb.server_timeout", 30*time.Second)

	// 索引配置
	m.viper.SetDefault("database.primary.mongodb.indexes.auto_create", true)

	// 验证配置
	m.viper.SetDefault("database.primary.mongodb.validation.enabled", true)

	// 迁移配置
	m.viper.SetDefault("database.migration.enabled", false)
	m.viper.SetDefault("database.migration.batch_size", 1000)
	m.viper.SetDefault("database.migration.timeout", 30*time.Minute)
	m.viper.SetDefault("database.migration.strategy.type", "full")
	m.viper.SetDefault("database.migration.strategy.verify_data", true)
	m.viper.SetDefault("database.migration.strategy.rollback_on_fail", true)

	// 同步配置
	m.viper.SetDefault("database.sync.enabled", false)
	m.viper.SetDefault("database.sync.interval", 5*time.Minute)
	m.viper.SetDefault("database.sync.batch_size", 100)
	m.viper.SetDefault("database.sync.direction", "source_to_target")
	m.viper.SetDefault("database.sync.conflict_resolution", "latest_wins")

	// 服务器配置
	m.viper.SetDefault("server.port", "8080")
	m.viper.SetDefault("server.mode", "debug")

	// JWT配置
	m.viper.SetDefault("jwt.secret", "qingyu_secret_key")
	m.viper.SetDefault("jwt.expiration_hours", 24)

	// AI配置
	m.viper.SetDefault("ai.base_url", "https://api.openai.com/v1")
	m.viper.SetDefault("ai.max_tokens", 2000)
	m.viper.SetDefault("ai.temperature", 7)
}

// applyDatabaseEnvOverrides 应用数据库环境变量覆盖
func (m *ViperConfigManager) applyDatabaseEnvOverrides(config *DatabaseConfig) error {
	// MongoDB环境变量覆盖
	if config.Primary.Type == DatabaseTypeMongoDB && config.Primary.MongoDB != nil {
		if uri := m.viper.GetString("QINGYU_DATABASE_PRIMARY_MONGODB_URI"); uri != "" {
			config.Primary.MongoDB.URI = uri
		}
		if db := m.viper.GetString("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE"); db != "" {
			config.Primary.MongoDB.Database = db
		}
		if poolSize := m.viper.GetUint64("QINGYU_DATABASE_PRIMARY_MONGODB_MAX_POOL_SIZE"); poolSize > 0 {
			config.Primary.MongoDB.MaxPoolSize = poolSize
		}
	}

	// PostgreSQL环境变量覆盖
	if config.Primary.Type == DatabaseTypePostgreSQL && config.Primary.PostgreSQL != nil {
		if host := m.viper.GetString("QINGYU_DATABASE_PRIMARY_POSTGRESQL_HOST"); host != "" {
			config.Primary.PostgreSQL.Host = host
		}
		if port := m.viper.GetInt("QINGYU_DATABASE_PRIMARY_POSTGRESQL_PORT"); port > 0 {
			config.Primary.PostgreSQL.Port = port
		}
		if db := m.viper.GetString("QINGYU_DATABASE_PRIMARY_POSTGRESQL_DATABASE"); db != "" {
			config.Primary.PostgreSQL.Database = db
		}
		if user := m.viper.GetString("QINGYU_DATABASE_PRIMARY_POSTGRESQL_USERNAME"); user != "" {
			config.Primary.PostgreSQL.Username = user
		}
		if password := m.viper.GetString("QINGYU_DATABASE_PRIMARY_POSTGRESQL_PASSWORD"); password != "" {
			config.Primary.PostgreSQL.Password = password
		}
	}

	return nil
}

// WatchConfig 监听配置变化
func (m *ViperConfigManager) WatchConfig(onChange func(*Config)) {
	m.viper.WatchConfig()
	m.viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件变更: %s\n", e.Name)

		// 重新加载配置
		config, err := m.LoadAppConfig()
		if err != nil {
			fmt.Printf("重新加载配置失败: %v\n", err)
			return
		}

		// 调用回调函数
		if onChange != nil {
			onChange(config)
		}
	})
}

// GetViper 获取viper实例
func (m *ViperConfigManager) GetViper() *viper.Viper {
	return m.viper
}

// 全局viper配置管理器实例
var GlobalViperManager *ViperConfigManager

// InitViperManager 初始化全局viper管理器
func InitViperManager(configPath string) error {
	GlobalViperManager = NewViperConfigManager()
	return GlobalViperManager.LoadFromFile(configPath)
}

// GetDatabaseConfigFromViper 从全局viper管理器获取数据库配置
func GetDatabaseConfigFromViper() (*DatabaseConfig, error) {
	if GlobalViperManager == nil {
		return nil, fmt.Errorf("viper管理器未初始化")
	}
	return GlobalViperManager.LoadDatabaseConfig()
}
