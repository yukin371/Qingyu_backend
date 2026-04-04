# Config 层快速参考

## 职责

配置管理层，负责配置加载、验证、管理、热重载、默认值管理。

## 目录结构

```
config/
├── config.go               # 主配置结构和加载函数
├── database.go             # 数据库配置
├── redis.go                # Redis配置
├── cache.go                # 缓存配置
├── jwt.go                  # JWT配置
├── oauth_config.go         # OAuth配置
├── feature_flags.go        # 功能开关配置
├── validation.go           # 配置验证
├── reload.go               # 热重载
├── env.go                  # 环境变量
└── viper_integration.go    # Viper集成
```

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | `{功能}.go` | `database.go`, `jwt.go` |
| 主配置 | `Config` | `Config` |
| 子配置 | `{功能}Config` | `DatabaseConfig`, `JWTConfig` |

## 配置结构

```go
type Config struct {
    Database      *DatabaseConfig      `mapstructure:"database"`
    Redis         *RedisConfig         `mapstructure:"redis"`
    Server        *ServerConfig        `mapstructure:"server"`
    Log           *LogConfig           `mapstructure:"log"`
    JWT           *JWTConfig           `mapstructure:"jwt"`
    RateLimit     *RateLimitConfig     `mapstructure:"rate_limit"`
    // ...
}
```

## 快速示例

```go
// 加载配置
config, err := config.LoadConfig("./configs")
if err != nil {
    log.Fatal(err)
}

// 访问配置
port := config.GlobalConfig.Server.Port
secret := config.GlobalConfig.JWT.Secret

// 通过 Viper 访问
secret := viper.GetString("jwt.secret")
```

## 配置优先级

```
1. 命令行参数（最高）
2. 环境变量
3. 配置文件
4. 默认值（最低）
```

## 环境变量命名

```
前缀: QINGYU_
分隔符: _ (点号替换为下划线)

示例:
  database.primary.mongodb.uri → QINGYU_DATABASE_PRIMARY_MONGODB_URI
  server.port → QINGYU_SERVER_PORT
  jwt.secret → QINGYU_JWT_SECRET
```

## 配置文件示例

```yaml
database:
  primary:
    type: mongodb
    mongodb:
      uri: mongodb://localhost:27017
      database: qingyu

server:
  port: "9090"
  mode: debug

jwt:
  secret: your-secret-key
  expiration_hours: 24

rate_limit:
  enabled: true
  requests_per_sec: 100
```

## 热重载

```go
config.WatchConfig(func() {
    fmt.Println("配置已更新")
    // 重新初始化依赖配置的组件
})
```

## 验证

```go
func (c *RateLimitConfig) Validate() error {
    if c.RequestsPerSec <= 0 {
        return fmt.Errorf("requests_per_sec must be positive")
    }
    return nil
}
```

## 禁止事项

- ❌ 硬编码配置值
- ❌ 配置结构体包含业务逻辑
- ❌ 配置验证时执行副作用
- ❌ 直接修改全局配置

## 详见

完整设计文档: [docs/standards/layer-config.md](../docs/standards/layer-config.md)
