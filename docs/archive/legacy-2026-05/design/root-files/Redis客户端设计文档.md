# Redis客户端设计文档

**版本**: v1.0  
**创建日期**: 2025-10-24  
**最后更新**: 2025-10-27  
**状态**: ✅ 已完成实施

---

## 1. 设计目标

### 1.1 核心目标
- 提供统一的Redis访问接口
- 支持连接池管理
- 支持自动重连
- 支持健康检查
- 易于测试和Mock

### 1.2 非功能性目标
- 高性能：响应时间<10ms
- 高可用：自动重连机制
- 易用性：简洁的API设计
- 可测试：便于单元测试

---

## 2. 配置设计

### 2.1 配置结构

```go
// RedisConfig Redis配置
type RedisConfig struct {
    // 连接配置
    Host     string `mapstructure:"host" json:"host"`         // Redis主机地址
    Port     int    `mapstructure:"port" json:"port"`         // Redis端口
    Password string `mapstructure:"password" json:"password"` // 密码
    DB       int    `mapstructure:"db" json:"db"`            // 数据库索引
    
    // 连接池配置
    PoolSize     int `mapstructure:"pool_size" json:"pool_size"`         // 连接池大小
    MinIdleConns int `mapstructure:"min_idle_conns" json:"min_idle_conns"` // 最小空闲连接数
    MaxIdleConns int `mapstructure:"max_idle_conns" json:"max_idle_conns"` // 最大空闲连接数
    
    // 超时配置
    DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout"`   // 连接超时
    ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`   // 读超时
    WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"` // 写超时
    PoolTimeout  time.Duration `mapstructure:"pool_timeout" json:"pool_timeout"`   // 连接池超时
    
    // 重连配置
    MaxRetries      int           `mapstructure:"max_retries" json:"max_retries"`           // 最大重试次数
    MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff" json:"min_retry_backoff"` // 最小重试间隔
    MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff" json:"max_retry_backoff"` // 最大重试间隔
}
```

### 2.2 默认配置

```go
func DefaultRedisConfig() *RedisConfig {
    return &RedisConfig{
        Host:     "localhost",
        Port:     6379,
        Password: "",
        DB:       0,
        
        PoolSize:     10,
        MinIdleConns: 5,
        MaxIdleConns: 10,
        
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolTimeout:  4 * time.Second,
        
        MaxRetries:      3,
        MinRetryBackoff: 8 * time.Millisecond,
        MaxRetryBackoff: 512 * time.Millisecond,
    }
}
```

### 2.3 配置文件示例

```yaml
# config.yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  max_idle_conns: 10
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
  pool_timeout: 4s
  max_retries: 3
  min_retry_backoff: 8ms
  max_retry_backoff: 512ms
```

---

## 3. 客户端接口设计

### 3.1 核心接口

```go
// RedisClient Redis客户端接口
type RedisClient interface {
    // 基础操作
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, keys ...string) (int64, error)
    
    // 批量操作
    MGet(ctx context.Context, keys ...string) ([]interface{}, error)
    MSet(ctx context.Context, pairs ...interface{}) error
    
    // 过期时间管理
    Expire(ctx context.Context, key string, expiration time.Duration) error
    TTL(ctx context.Context, key string) (time.Duration, error)
    
    // Hash操作
    HGet(ctx context.Context, key, field string) (string, error)
    HSet(ctx context.Context, key string, values ...interface{}) error
    HGetAll(ctx context.Context, key string) (map[string]string, error)
    HDel(ctx context.Context, key string, fields ...string) error
    
    // Set操作
    SAdd(ctx context.Context, key string, members ...interface{}) error
    SMembers(ctx context.Context, key string) ([]string, error)
    SRem(ctx context.Context, key string, members ...interface{}) error
    
    // 原子操作
    Incr(ctx context.Context, key string) (int64, error)
    Decr(ctx context.Context, key string) (int64, error)
    IncrBy(ctx context.Context, key string, value int64) (int64, error)
    
    // 健康检查
    Ping(ctx context.Context) error
    
    // 生命周期
    Close() error
    
    // 获取原始客户端（用于高级操作）
    GetClient() interface{}
}
```

---

## 4. 实现设计

### 4.1 实现结构

```go
// redisClientImpl Redis客户端实现
type redisClientImpl struct {
    client *redis.Client
    config *RedisConfig
    mu     sync.RWMutex
}
```

### 4.2 初始化流程

```
1. 加载配置
   ↓
2. 创建redis.Client
   ↓
3. 设置连接池参数
   ↓
4. 设置超时参数
   ↓
5. 设置重连参数
   ↓
6. 健康检查（Ping）
   ↓
7. 返回客户端实例
```

### 4.3 错误处理

```go
// 错误类型
var (
    ErrRedisNil        = errors.New("redis: nil returned")
    ErrKeyNotFound     = errors.New("redis: key not found")
    ErrConnectionFailed = errors.New("redis: connection failed")
    ErrTimeout         = errors.New("redis: operation timeout")
)

// 错误包装
func wrapRedisError(err error) error {
    if err == redis.Nil {
        return ErrRedisNil
    }
    if strings.Contains(err.Error(), "timeout") {
        return ErrTimeout
    }
    return err
}
```

---

## 5. 使用场景设计

### 5.1 缓存场景

```go
// 场景1：缓存用户信息
func CacheUserInfo(ctx context.Context, userID string, user *User) error {
    data, _ := json.Marshal(user)
    return redisClient.Set(ctx, "user:"+userID, data, 5*time.Minute)
}

func GetUserInfoFromCache(ctx context.Context, userID string) (*User, error) {
    data, err := redisClient.Get(ctx, "user:"+userID)
    if err != nil {
        return nil, err
    }
    var user User
    json.Unmarshal([]byte(data), &user)
    return &user, nil
}
```

### 5.2 会话管理场景

```go
// 场景2：会话管理
func StoreSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
    values := make([]interface{}, 0)
    for k, v := range data {
        values = append(values, k, v)
    }
    return redisClient.HSet(ctx, "session:"+sessionID, values...)
}

func GetSession(ctx context.Context, sessionID string) (map[string]string, error) {
    return redisClient.HGetAll(ctx, "session:"+sessionID)
}
```

### 5.3 Token黑名单场景

```go
// 场景3：Token黑名单
func AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
    return redisClient.Set(ctx, "blacklist:"+token, "1", expiration)
}

func IsInBlacklist(ctx context.Context, token string) (bool, error) {
    count, err := redisClient.Exists(ctx, "blacklist:"+token)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
```

### 5.4 配额管理场景

```go
// 场景4：配额管理
func DeductQuota(ctx context.Context, userID string, amount int64) (int64, error) {
    remaining, err := redisClient.IncrBy(ctx, "quota:"+userID, -amount)
    if err != nil {
        return 0, err
    }
    return remaining, nil
}

func GetQuota(ctx context.Context, userID string) (int64, error) {
    data, err := redisClient.Get(ctx, "quota:"+userID)
    if err != nil {
        return 0, err
    }
    return strconv.ParseInt(data, 10, 64)
}
```

---

## 6. 集成设计

### 6.1 服务容器集成

```go
// ServiceContainer中添加
type ServiceContainer struct {
    // ...
    redisClient RedisClient
}

// 初始化
func (c *ServiceContainer) initRedis() error {
    cfg := config.GetRedisConfig()
    client, err := NewRedisClient(cfg)
    if err != nil {
        return fmt.Errorf("初始化Redis客户端失败: %w", err)
    }
    c.redisClient = client
    return nil
}

// Getter
func (c *ServiceContainer) GetRedisClient() RedisClient {
    return c.redisClient
}
```

### 6.2 依赖注入

```go
// 示例：JWTService使用Redis
type JWTServiceImpl struct {
    redisClient RedisClient
    // ...
}

func NewJWTService(redisClient RedisClient, ...) JWTService {
    return &JWTServiceImpl{
        redisClient: redisClient,
        // ...
    }
}
```

---

## 7. 测试设计

### 7.1 单元测试

```go
// 使用Mock进行单元测试
type MockRedisClient struct {
    mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
    args := m.Called(ctx, key)
    return args.String(0), args.Error(1)
}

// 测试示例
func TestCacheUserInfo(t *testing.T) {
    mockRedis := new(MockRedisClient)
    mockRedis.On("Set", mock.Anything, "user:123", mock.Anything, mock.Anything).
        Return(nil)
    
    // 测试逻辑...
}
```

### 7.2 集成测试

```go
// 使用testcontainers启动真实Redis进行集成测试
func TestRedisClient_Integration(t *testing.T) {
    // 启动Redis容器
    ctx := context.Background()
    redisContainer, err := testcontainers.GenericContainer(ctx, ...)
    
    // 创建客户端
    client, err := NewRedisClient(cfg)
    
    // 测试操作
    err = client.Set(ctx, "test", "value", time.Minute)
    assert.NoError(t, err)
    
    val, err := client.Get(ctx, "test")
    assert.NoError(t, err)
    assert.Equal(t, "value", val)
}
```

---

## 8. 性能考虑

### 8.1 性能目标

| 操作 | 目标响应时间 | 备注 |
|------|-------------|------|
| Get | <5ms | 本地Redis |
| Set | <5ms | 本地Redis |
| MGet | <10ms | 批量操作 |
| Hash操作 | <8ms | - |

### 8.2 优化策略

1. **连接池优化**
   - 合理设置连接池大小
   - 预热连接池

2. **批量操作**
   - 使用Pipeline减少RTT
   - 使用MGet/MSet批量处理

3. **序列化优化**
   - 使用MessagePack替代JSON（可选）
   - 压缩大数据

---

## 9. 安全考虑

### 9.1 安全措施

1. **密码保护**
   - 配置文件中的密码加密
   - 环境变量注入

2. **连接安全**
   - TLS加密连接（生产环境）
   - IP白名单

3. **数据安全**
   - 敏感数据加密后存储
   - 设置合理的过期时间

---

## 10. 监控设计

### 10.1 监控指标

```go
type RedisMetrics struct {
    TotalCommands    int64
    SuccessCommands  int64
    FailedCommands   int64
    TotalLatency     time.Duration
    ConnectionErrors int64
    TimeoutErrors    int64
}
```

### 10.2 指标收集

- 在每次操作后记录指标
- 定期上报到Prometheus
- 异常情况告警

---

## 11. 实施计划

### 阶段1：核心实现（Day 1）
- [ ] 配置结构实现
- [ ] 客户端基础实现
- [ ] 基础操作（Get、Set、Delete）

### 阶段2：扩展功能（Day 1-2）
- [ ] 批量操作
- [ ] Hash/Set操作
- [ ] 健康检查

### 阶段3：测试（Day 2）
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能测试

### 阶段4：集成（Day 2）
- [ ] 服务容器集成
- [ ] 文档完善

---

## 12. 附录

### 12.1 参考资料
- [go-redis文档](https://redis.uptrace.dev/)
- [Redis最佳实践](https://redis.io/docs/manual/patterns/)

### 12.2 待优化项
- [ ] Pipeline支持
- [ ] Pub/Sub支持
- [ ] Lua脚本支持
- [ ] 集群支持

---

**设计状态**: ✅ 完成  
**审核状态**: 待审核  
**最后更新**: 2025-10-24

