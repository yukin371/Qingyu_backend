# Config 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 职责边界与依赖关系

### 1.1 职责定义

Config 层是后端系统的**配置管理层**，负责：

1. **配置加载**：从文件、环境变量加载配置
2. **配置验证**：验证配置值的有效性
3. **配置管理**：统一管理应用配置
4. **配置热重载**：支持运行时重新加载配置
5. **默认值管理**：提供合理的默认配置值

### 1.2 与上下层级的交互边界

```
┌─────────────────────────────────────────────────────────┐
│                   应用启动                              │
│              (main.go)                                 │
└─────────────────────────────────────────────────────────┘
                         │
                    加载配置
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Config 层                            │
│  ┌─────────────────────────────────────────────────────┤
│  │ 职责：                                              │
│  │ - 读取配置文件 (YAML/JSON)                         │
│  │ - 读取环境变量                                      │
│  │ - 设置默认值                                        │
│  │ - 验证配置有效性                                    │
│  │ - 热重载支持                                        │
│  │ - 配置访问接口                                      │
│  └─────────────────────────────────────────────────────┤
│  输出：                                                  │
│  │ - GlobalConfig 全局配置实例                        │
│  │ - 各子配置结构体                                    │
└─────────────────────────────────────────────────────────┘
                         │
                    提供配置
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│              其他层级（Service、Repository等）          │
│              (消费配置)                                │
└─────────────────────────────────────────────────────────┘
```

### 1.3 依赖关系

```go
// Config 层允许的依赖
import (
    "github.com/spf13/viper"        // 配置管理
    "github.com/fsnotify/fsnotify"  // 文件监控
)

// Config 层禁止的依赖
import (
    "Qingyu_backend/service/xxx"    // ❌ 禁止依赖 Service
    "Qingyu_backend/api/xxx"        // ❌ 禁止依赖 API
    "go.mongodb.org/mongo-driver"   // ❌ 禁止直接操作数据库
)
```

---

## 2. 命名与代码规范

### 2.1 文件命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 主配置 | `config.go` | `config.go` |
| 子配置 | `{功能}.go` | `database.go`, `jwt.go`, `cache.go` |
| 验证 | `validation.go` | `validation.go` |
| 热重载 | `reload.go` | `reload.go` |
| 测试 | `{文件名}_test.go` | `database_test.go` |

### 2.2 结构体命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 主配置 | `Config` | `Config` |
| 子配置 | `{功能}Config` | `DatabaseConfig`, `JWTConfig`, `RedisConfig` |
| 配置选项 | `{功能}Options` | `ConnectionOptions`, `PoolOptions` |

### 2.3 目录组织规范

```
config/
├── config.go               # 主配置结构和加载函数
├── database.go             # 数据库配置
├── database_test.go        # 数据库配置测试
├── redis.go                # Redis配置
├── cache.go                # 缓存配置
├── jwt.go                  # JWT配置
├── oauth_config.go         # OAuth配置
├── rate_limit_test.go      # 限流配置测试
├── feature_flags.go        # 功能开关配置
├── validation.go           # 配置验证
├── reload.go               # 热重载
├── env.go                  # 环境变量
└── viper_integration.go    # Viper集成
```

---

## 3. 设计模式与最佳实践

### 3.1 配置结构模式

```go
// config.go
package config

// Config 存储应用配置
type Config struct {
    Database      *DatabaseConfig                   `mapstructure:"database"`
    Redis         *RedisConfig                      `mapstructure:"redis"`
    Cache         *CacheConfig                      `mapstructure:"cache"`
    Server        *ServerConfig                     `mapstructure:"server"`
    Log           *LogConfig                        `mapstructure:"log"`
    JWT           *JWTConfig                        `mapstructure:"jwt"`
    AI            *AIConfig                         `mapstructure:"ai"`
    Elasticsearch *ElasticsearchConfig              `mapstructure:"elasticsearch"`
    External      *ExternalAPIConfig                `mapstructure:"external"`
    AIQuota       *AIQuotaConfig                    `mapstructure:"ai_quota"`
    Email         *EmailConfig                      `mapstructure:"email"`
    Payment       *PaymentConfig                    `mapstructure:"payment"`
    RateLimit     *RateLimitConfig                  `mapstructure:"rate_limit"`
    OAuth         map[string]*authModel.OAuthConfig `mapstructure:"oauth"`
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
```

### 3.2 配置加载模式

```go
// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
    v = viper.New()

    // 1. 设置默认值
    setDefaults()

    // 2. 检查配置文件路径
    if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
        v.SetConfigFile(configPath)
    } else {
        v.SetConfigName("config")
        v.SetConfigType("yaml")
        v.AddConfigPath(configPath)
        v.AddConfigPath("./configs")
        v.AddConfigPath(".")
    }

    // 3. 启用环境变量
    v.AutomaticEnv()
    v.SetEnvPrefix("QINGYU")
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 4. 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            fmt.Printf("Warning: Config file not found: %v\n", err)
        } else {
            return nil, fmt.Errorf("fatal error reading config file: %w", err)
        }
    }

    // 5. 解析配置
    config := &Config{}
    if err := v.Unmarshal(config); err != nil {
        return nil, fmt.Errorf("unable to decode config into struct: %w", err)
    }

    // 6. 设置默认值（处理嵌套结构）
    setConfigDefaults(config)

    // 7. 验证配置
    if err := ValidateConfig(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    // 8. 设置全局配置实例
    GlobalConfig = config

    return config, nil
}
```

### 3.3 默认值模式

```go
// setDefaults 设置Viper默认值
func setDefaults() {
    // 数据库默认配置
    v.SetDefault("database.primary.type", "mongodb")
    v.SetDefault("database.primary.mongodb.uri", "mongodb://localhost:27017")
    v.SetDefault("database.primary.mongodb.database", "qingyu")
    v.SetDefault("database.primary.mongodb.max_pool_size", 100)
    v.SetDefault("database.primary.mongodb.min_pool_size", 5)

    // 服务器默认配置
    v.SetDefault("server.port", "9090")
    v.SetDefault("server.mode", "debug")

    // 日志默认配置
    v.SetDefault("log.level", "info")
    v.SetDefault("log.format", "json")
    v.SetDefault("log.output", "stdout")

    // JWT默认配置
    v.SetDefault("jwt.secret", "qingyu_secret_key")
    v.SetDefault("jwt.expiration_hours", 24)

    // 速率限制默认配置
    v.SetDefault("rate_limit.enabled", true)
    v.SetDefault("rate_limit.requests_per_sec", 100)
    v.SetDefault("rate_limit.burst", 200)
}
```

### 3.4 配置验证模式

```go
// validation.go

// ValidateConfig 验证配置
func ValidateConfig(config *Config) error {
    if config == nil {
        return errors.New("config is nil")
    }

    // 验证数据库配置
    if err := config.Database.Validate(); err != nil {
        return fmt.Errorf("database config: %w", err)
    }

    // 验证JWT配置
    if err := config.JWT.Validate(); err != nil {
        return fmt.Errorf("jwt config: %w", err)
    }

    // 验证限流配置
    if err := config.RateLimit.Validate(); err != nil {
        return fmt.Errorf("rate_limit config: %w", err)
    }

    return nil
}

// Validate 验证速率限制配置
func (c *RateLimitConfig) Validate() error {
    if !c.Enabled {
        return nil  // 禁用时不验证
    }

    if c.RequestsPerSec <= 0 {
        return fmt.Errorf("requests_per_sec must be positive, got %f", c.RequestsPerSec)
    }
    if c.Burst < 0 {
        return fmt.Errorf("burst cannot be negative, got %d", c.Burst)
    }
    if c.Burst > 10000 {
        return fmt.Errorf("burst too large (%d), maximum is 10000", c.Burst)
    }
    return nil
}
```

### 3.5 热重载模式

```go
// reload.go

// WatchConfig 启用配置热重载
func WatchConfig(onChange func()) {
    v.WatchConfig()
    if onChange != nil {
        v.OnConfigChange(func(e fsnotify.Event) {
            fmt.Printf("Config file changed: %s\n", e.Name)

            // 1. 重新加载配置
            config := &Config{}
            if err := v.Unmarshal(config); err != nil {
                fmt.Printf("Error reloading config: %v\n", err)
                return
            }

            // 2. 设置默认值
            setConfigDefaults(config)

            // 3. 验证配置
            if err := ValidateConfig(config); err != nil {
                fmt.Printf("Error validating reloaded config: %v\n", err)
                return
            }

            // 4. 原子更新全局配置
            GlobalConfig = config

            // 5. 调用回调函数
            onChange()
        })
    }
}
```

### 3.6 环境变量覆盖模式

```go
// env.go

// 环境变量命名规范：
// 前缀: QINGYU_
// 分隔符: _ (点号替换为下划线)
// 示例:
//   database.primary.mongodb.uri -> QINGYU_DATABASE_PRIMARY_MONGODB_URI
//   server.port -> QINGYU_SERVER_PORT
//   jwt.secret -> QINGYU_JWT_SECRET

// GetEnvOrDefault 获取环境变量或默认值
func GetEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// GetEnvIntOrDefault 获取整数环境变量或默认值
func GetEnvIntOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

// GetEnvBoolOrDefault 获取布尔环境变量或默认值
func GetEnvBoolOrDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        switch strings.ToLower(value) {
        case "true", "1", "yes", "on":
            return true
        case "false", "0", "no", "off":
            return false
        }
    }
    return defaultValue
}
```

### 3.7 配置文件结构

```yaml
# config.yaml 示例
database:
  primary:
    type: mongodb
    mongodb:
      uri: mongodb://localhost:27017
      database: qingyu
      max_pool_size: 100
      min_pool_size: 5

redis:
  enabled: true
  host: localhost
  port: 6379
  password: ""
  db: 0

server:
  port: "9090"
  mode: debug

log:
  level: info
  format: json
  output: stdout

jwt:
  secret: your-secret-key
  expiration_hours: 24

rate_limit:
  enabled: true
  requests_per_sec: 100
  burst: 200
  skip_paths:
    - /health
    - /metrics
```

### 3.8 反模式警示

❌ **禁止事项**：

```go
// ❌ 禁止：硬编码配置值
func NewService() *Service {
    return &Service{
        timeout: 30 * time.Second,  // 应该从配置读取
    }
}

// ❌ 禁止：在配置结构体中包含业务逻辑
type DatabaseConfig struct {
    URI string
}
func (c *DatabaseConfig) Connect() {  // 配置只存储数据，不执行操作
}

// ❌ 禁止：配置验证时执行副作用
func (c *Config) Validate() error {
    os.Create(c.Log.Filename)  // 禁止在验证中创建文件
    return nil
}

// ❌ 禁止：全局可变配置
GlobalConfig.Database.URI = "new-uri"  // 配置应该是只读的
```

---

## 4. 配置访问规范

### 4.1 全局配置访问

```go
// 通过全局变量访问
if config.GlobalConfig != nil {
    secret := config.GlobalConfig.JWT.Secret
    port := config.GlobalConfig.Server.Port
}

// 通过 Viper 访问
secret := viper.GetString("jwt.secret")
port := viper.GetString("server.port")
```

### 4.2 子配置加载

```go
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
```

### 4.3 配置优先级

```
1. 命令行参数（最高优先级）
2. 环境变量
3. 配置文件
4. 默认值（最低优先级）
```

---

## 5. 测试策略

### 5.1 配置测试编写指南

```go
// database_test.go
package config

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDatabaseConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  DatabaseConfig
        wantErr bool
    }{
        {
            name: "valid config",
            config: DatabaseConfig{
                Primary: &PrimaryDBConfig{
                    Type: "mongodb",
                    MongoDB: &MongoDBConfig{
                        URI:      "mongodb://localhost:27017",
                        Database: "qingyu",
                    },
                },
            },
            wantErr: false,
        },
        {
            name: "missing URI",
            config: DatabaseConfig{
                Primary: &PrimaryDBConfig{
                    Type: "mongodb",
                    MongoDB: &MongoDBConfig{
                        Database: "qingyu",
                    },
                },
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestLoadConfig(t *testing.T) {
    config, err := LoadConfig("../../configs")

    assert.NoError(t, err)
    assert.NotNil(t, config)
    assert.NotEmpty(t, config.Server.Port)
}
```

### 5.2 测试覆盖率要求

| 测试类型 | 覆盖率要求 |
|----------|------------|
| 配置验证 | 100% |
| 默认值 | 100% |
| 加载逻辑 | ≥ 80% |

---

## 6. 完整代码示例

### 6.1 完整配置模块示例

```go
// database.go
package config

import (
    "errors"
    "fmt"
    "time"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    Primary   *PrimaryDBConfig   `mapstructure:"primary"`
    Secondary *SecondaryDBConfig `mapstructure:"secondary,omitempty"`
    Indexing  *IndexingConfig    `mapstructure:"indexing"`
    Migration *MigrationConfig   `mapstructure:"migration"`
    Sync      *SyncConfig        `mapstructure:"sync"`
}

// PrimaryDBConfig 主数据库配置
type PrimaryDBConfig struct {
    Type     string        `mapstructure:"type"`
    MongoDB  *MongoDBConfig `mapstructure:"mongodb"`
    MySQL    *MySQLConfig  `mapstructure:"mysql,omitempty"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
    URI            string        `mapstructure:"uri"`
    Database       string        `mapstructure:"database"`
    MaxPoolSize    uint64        `mapstructure:"max_pool_size"`
    MinPoolSize    uint64        `mapstructure:"min_pool_size"`
    ConnectTimeout time.Duration `mapstructure:"connect_timeout"`
    ServerTimeout  time.Duration `mapstructure:"server_timeout"`
}

// Validate 验证数据库配置
func (c *DatabaseConfig) Validate() error {
    if c.Primary == nil {
        return errors.New("primary database config is required")
    }

    switch c.Primary.Type {
    case "mongodb":
        if c.Primary.MongoDB == nil {
            return errors.New("mongodb config is required when type is mongodb")
        }
        if c.Primary.MongoDB.URI == "" {
            return errors.New("mongodb uri is required")
        }
        if c.Primary.MongoDB.Database == "" {
            return errors.New("mongodb database name is required")
        }
    case "mysql":
        if c.Primary.MySQL == nil {
            return errors.New("mysql config is required when type is mysql")
        }
    default:
        return fmt.Errorf("unsupported database type: %s", c.Primary.Type)
    }

    return nil
}

// GetConnectionString 获取连接字符串
func (c *DatabaseConfig) GetConnectionString() string {
    if c.Primary == nil {
        return ""
    }
    switch c.Primary.Type {
    case "mongodb":
        return c.Primary.MongoDB.URI
    case "mysql":
        return c.Primary.MySQL.DSN
    }
    return ""
}

// GetDatabaseName 获取数据库名称
func (c *DatabaseConfig) GetDatabaseName() string {
    if c.Primary == nil {
        return ""
    }
    switch c.Primary.Type {
    case "mongodb":
        return c.Primary.MongoDB.Database
    case "mysql":
        return c.Primary.MySQL.Database
    }
    return ""
}
```

---

## 7. 参考资料

- [Config 层快速参考](../config/README.md)
- [Viper 文档](https://github.com/spf13/viper)
- [Service 层设计说明](./layer-service.md)

---

*最后更新：2026-03-19*
