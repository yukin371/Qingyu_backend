package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// DatabaseType 数据库类型
type DatabaseType string

const (
	DatabaseTypeMongoDB    DatabaseType = "mongodb"
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeMySQL      DatabaseType = "mysql"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type       string                `yaml:"type" json:"type" mapstructure:"type"`
	Primary    DatabaseConnection    `yaml:"primary" json:"primary" mapstructure:"primary"`
	Replicas   []DatabaseConnection  `yaml:"replicas" json:"replicas" mapstructure:"replicas"`
	Indexing   IndexingConfig        `yaml:"indexing" json:"indexing" mapstructure:"indexing"`
	Validation ValidationConfig      `yaml:"validation" json:"validation" mapstructure:"validation"`
	Sync       SynchronizationConfig `yaml:"sync" json:"sync" mapstructure:"sync"`
}

// DatabaseConnection 数据库连接配置
type DatabaseConnection struct {
	Type       DatabaseType      `yaml:"type" json:"type" mapstructure:"type"`
	MongoDB    *MongoDBConfig    `yaml:"mongodb,omitempty" json:"mongodb,omitempty" mapstructure:"mongodb,omitempty"`
	PostgreSQL *PostgreSQLConfig `yaml:"postgresql,omitempty" json:"postgresql,omitempty" mapstructure:"postgresql,omitempty"`
	MySQL      *MySQLConfig      `yaml:"mysql,omitempty" json:"mysql,omitempty" mapstructure:"mysql,omitempty"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI            string        `yaml:"uri" json:"uri" mapstructure:"uri"`
	Database       string        `yaml:"database" json:"database" mapstructure:"database"`
	MaxPoolSize    uint64        `yaml:"max_pool_size" json:"max_pool_size" mapstructure:"max_pool_size"`
	MinPoolSize    uint64        `yaml:"min_pool_size" json:"min_pool_size" mapstructure:"min_pool_size"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" json:"connect_timeout" mapstructure:"connect_timeout"`
	ServerTimeout  time.Duration `yaml:"server_timeout" json:"server_timeout" mapstructure:"server_timeout"`

	// Profiling配置
	ProfilingLevel int   `yaml:"profiling_level" json:"profiling_level" mapstructure:"profiling_level"`
	SlowMS         int64 `yaml:"slow_ms" json:"slow_ms" mapstructure:"slow_ms"`
	ProfilerSizeMB int64 `yaml:"profiler_size_mb" json:"profiler_size_mb" mapstructure:"profiler_size_mb"`
}

// PostgreSQLConfig PostgreSQL配置
type PostgreSQLConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	Username     string        `yaml:"username" json:"username"`
	Password     string        `yaml:"password" json:"password"`
	SSLMode      string        `yaml:"ssl_mode" json:"ssl_mode"`
	MaxOpenConns int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnTimeout  time.Duration `yaml:"conn_timeout" json:"conn_timeout"`

	// 迁移配置
	MigrationsPath string `yaml:"migrations_path" json:"migrations_path"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	Username     string        `yaml:"username" json:"username"`
	Password     string        `yaml:"password" json:"password"`
	Charset      string        `yaml:"charset" json:"charset"`
	MaxOpenConns int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnTimeout  time.Duration `yaml:"conn_timeout" json:"conn_timeout"`
}

// IndexingConfig 索引配置
type IndexingConfig struct {
	AutoCreate bool `yaml:"auto_create" json:"auto_create" mapstructure:"auto_create"`
	Background bool `yaml:"background" json:"background" mapstructure:"background"`
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	Enabled    bool `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	StrictMode bool `yaml:"strict_mode" json:"strict_mode" mapstructure:"strict_mode"`
}

// SynchronizationConfig 同步配置
type SynchronizationConfig struct {
	Enabled  bool          `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	Interval time.Duration `yaml:"interval" json:"interval" mapstructure:"interval"`
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig(configPath string) (*DatabaseConfig, error) {
	var config DatabaseConfig

	// 如果没有指定配置文件路径，使用默认配置
	if configPath == "" {
		config = *getDefaultDatabaseConfig()
	} else {
		// 读取配置文件
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 应用环境变量覆盖
	if err := applyEnvironmentOverrides(&config); err != nil {
		return nil, fmt.Errorf("应用环境变量失败: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// getDefaultDatabaseConfig 获取默认数据库配置
func getDefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type: "mongodb",
		Primary: DatabaseConnection{
			Type: DatabaseTypeMongoDB,
			MongoDB: &MongoDBConfig{
				URI:            getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017"),
				Database:       getEnvOrDefault("MONGODB_DATABASE", "qingyu"),
				MaxPoolSize:    100,
				MinPoolSize:    5,
				ConnectTimeout: 10 * time.Second,
				ServerTimeout:  30 * time.Second,
				ProfilingLevel: 1,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
		},
		Replicas: []DatabaseConnection{},
		Indexing: IndexingConfig{
			AutoCreate: true,
			Background: true,
		},
		Validation: ValidationConfig{
			Enabled:    true,
			StrictMode: false,
		},
		Sync: SynchronizationConfig{
			Enabled:  false,
			Interval: 5 * time.Minute,
		},
	}
}

// applyEnvironmentOverrides 应用环境变量覆盖
func applyEnvironmentOverrides(config *DatabaseConfig) error {
	// MongoDB环境变量覆盖
	if config.Primary.Type == DatabaseTypeMongoDB && config.Primary.MongoDB != nil {
		if uri := os.Getenv("MONGODB_URI"); uri != "" {
			config.Primary.MongoDB.URI = uri
		}
		if db := os.Getenv("MONGODB_DATABASE"); db != "" {
			config.Primary.MongoDB.Database = db
		}
		if poolSize := os.Getenv("MONGODB_MAX_POOL_SIZE"); poolSize != "" {
			if size, err := strconv.ParseUint(poolSize, 10, 64); err == nil {
				config.Primary.MongoDB.MaxPoolSize = size
			}
		}

		// MongoDB Profiling环境变量覆盖
		if level := os.Getenv("MONGODB_PROFILING_LEVEL"); level != "" {
			if l, err := strconv.Atoi(level); err == nil && l >= 0 && l <= 2 {
				config.Primary.MongoDB.ProfilingLevel = l
			}
		}
		if slowMS := os.Getenv("MONGODB_SLOW_MS"); slowMS != "" {
			if ms, err := strconv.ParseInt(slowMS, 10, 64); err == nil && ms >= 0 {
				config.Primary.MongoDB.SlowMS = ms
			}
		}
		if sizeMB := os.Getenv("MONGODB_PROFILER_SIZE_MB"); sizeMB != "" {
			if mb, err := strconv.ParseInt(sizeMB, 10, 64); err == nil && mb >= 1 {
				config.Primary.MongoDB.ProfilerSizeMB = mb
			}
		}
	}

	// PostgreSQL环境变量覆盖
	if config.Primary.Type == DatabaseTypePostgreSQL && config.Primary.PostgreSQL != nil {
		if host := os.Getenv("POSTGRES_HOST"); host != "" {
			config.Primary.PostgreSQL.Host = host
		}
		if port := os.Getenv("POSTGRES_PORT"); port != "" {
			if p, err := strconv.Atoi(port); err == nil {
				config.Primary.PostgreSQL.Port = p
			}
		}
		if db := os.Getenv("POSTGRES_DB"); db != "" {
			config.Primary.PostgreSQL.Database = db
		}
		if user := os.Getenv("POSTGRES_USER"); user != "" {
			config.Primary.PostgreSQL.Username = user
		}
		if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
			config.Primary.PostgreSQL.Password = password
		}
	}

	return nil
}

// Validate 验证配置
func (c *DatabaseConfig) Validate() error {
	// 验证主数据库配置
	if err := c.Primary.Validate(); err != nil {
		return fmt.Errorf("主数据库配置无效: %w", err)
	}

	// 验证副本数据库配置
	for i, replica := range c.Replicas {
		if err := replica.Validate(); err != nil {
			return fmt.Errorf("副本数据库 %d 配置无效: %w", i, err)
		}
	}

	return nil
}

// Validate 验证数据库实例配置
func (c *DatabaseConnection) Validate() error {
	switch c.Type {
	case DatabaseTypeMongoDB:
		if c.MongoDB == nil {
			return fmt.Errorf("MongoDB配置不能为空")
		}
		return c.MongoDB.Validate()
	case DatabaseTypePostgreSQL:
		if c.PostgreSQL == nil {
			return fmt.Errorf("PostgreSQL配置不能为空")
		}
		return c.PostgreSQL.Validate()
	case DatabaseTypeMySQL:
		if c.MySQL == nil {
			return fmt.Errorf("MySQL配置不能为空")
		}
		return c.MySQL.Validate()
	default:
		return fmt.Errorf("不支持的数据库类型: %s", c.Type)
	}
}

// Validate 验证MongoDB配置
func (c *MongoDBConfig) Validate() error {
	if c.URI == "" {
		return fmt.Errorf("MongoDB URI不能为空")
	}
	if c.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if c.ProfilingLevel < 0 || c.ProfilingLevel > 2 {
		return fmt.Errorf("ProfilingLevel必须在0-2之间")
	}
	if c.SlowMS < 0 {
		return fmt.Errorf("SlowMS必须非负")
	}
	if c.ProfilerSizeMB < 1 {
		return fmt.Errorf("ProfilerSizeMB必须至少为1MB")
	}
	if c.MaxPoolSize == 0 {
		c.MaxPoolSize = 100
	}
	if c.MinPoolSize == 0 {
		c.MinPoolSize = 5
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10 * time.Second
	}
	if c.ServerTimeout == 0 {
		c.ServerTimeout = 30 * time.Second
	}
	return nil
}

// Validate 验证PostgreSQL配置
func (c *PostgreSQLConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("PostgreSQL主机地址不能为空")
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if c.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 100
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnTimeout == 0 {
		c.ConnTimeout = 10 * time.Second
	}
	return nil
}

// Validate 验证MySQL配置
func (c *MySQLConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("MySQL主机地址不能为空")
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if c.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if c.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 100
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnTimeout == 0 {
		c.ConnTimeout = 10 * time.Second
	}
	return nil
}

// ToRepositoryConfig 转换为仓储配置
// ToRepositoryConfig 转换为仓储配置 - 返回通用配置接口
func (c *MongoDBConfig) ToRepositoryConfig() map[string]interface{} {
	return map[string]interface{}{
		"type":             "mongodb",
		"uri":              c.URI,
		"database":         c.Database,
		"max_pool_size":    c.MaxPoolSize,
		"min_pool_size":    c.MinPoolSize,
		"connect_timeout":  c.ConnectTimeout,
		"server_timeout":   c.ServerTimeout,
		"profiling_level":  c.ProfilingLevel,
		"slow_ms":          c.SlowMS,
		"profiler_size_mb": c.ProfilerSizeMB,
	}
}

// ToRepositoryConfig 转换为仓储配置 - 返回通用配置接口
func (c *PostgreSQLConfig) ToRepositoryConfig() map[string]interface{} {
	return map[string]interface{}{
		"type":           "postgresql",
		"host":           c.Host,
		"port":           c.Port,
		"database":       c.Database,
		"username":       c.Username,
		"password":       c.Password,
		"ssl_mode":       c.SSLMode,
		"max_open_conns": c.MaxOpenConns,
		"max_idle_conns": c.MaxIdleConns,
		"conn_timeout":   c.ConnTimeout,
	}
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SaveDatabaseConfig 保存数据库配置到文件
func SaveDatabaseConfig(config *DatabaseConfig, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// DatabaseConfigManager 数据库配置管理器
type DatabaseConfigManager struct {
	config     *DatabaseConfig
	configPath string
}

// NewDatabaseConfigManager 创建数据库配置管理器
func NewDatabaseConfigManager(configPath string) (*DatabaseConfigManager, error) {
	config, err := LoadDatabaseConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &DatabaseConfigManager{
		config:     config,
		configPath: configPath,
	}, nil
}

// GetConfig 获取配置
func (m *DatabaseConfigManager) GetConfig() *DatabaseConfig {
	return m.config
}

// UpdateConfig 更新配置
func (m *DatabaseConfigManager) UpdateConfig(config *DatabaseConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	m.config = config

	// 保存到文件
	if m.configPath != "" {
		return SaveDatabaseConfig(config, m.configPath)
	}

	return nil
}

// EnableMigration 启用迁移
func (m *DatabaseConfigManager) EnableMigration(secondary DatabaseConnection) error {
	// 添加副本数据库
	m.config.Replicas = append(m.config.Replicas, secondary)
	return m.UpdateConfig(m.config)
}

// DisableMigration 禁用迁移
func (m *DatabaseConfigManager) DisableMigration() error {
	// 清空副本数据库
	m.config.Replicas = []DatabaseConnection{}
	return m.UpdateConfig(m.config)
}

// EnableSync 启用同步
func (m *DatabaseConfigManager) EnableSync(syncConfig SynchronizationConfig) error {
	m.config.Sync = syncConfig
	m.config.Sync.Enabled = true

	return m.UpdateConfig(m.config)
}

// DisableSync 禁用同步
func (m *DatabaseConfigManager) DisableSync() error {
	m.config.Sync = SynchronizationConfig{Enabled: false}
	return m.UpdateConfig(m.config)
}
