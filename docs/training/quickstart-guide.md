# 中间件快速上手指南

## 概述

本指南提供具体的代码示例，帮助开发者快速上手新的中间件架构。

**前置要求**：
- Go 1.21+
- 熟悉 Gin 框架
- 了解基本的中间件概念

**预计阅读时间**: 30 分钟

---

## 目录

1. [示例1: 创建静态配置中间件](#示例1-创建静态配置中间件)
2. [示例2: 创建动态配置中间件](#示例2-创建动态配置中间件)
3. [配置文件示例](#配置文件示例)
4. [单元测试示例](#单元测试示例)
5. [常见问题](#常见问题)

---

## 示例1: 创建静态配置中间件

### 场景

创建一个 `RequestID` 中间件，为每个请求生成唯一 ID。

### 完整代码

#### 1. 定义配置结构体

```go
// internal/middleware/base/request_id.go
package base

import (
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
    "strings"
)

// RequestIDConfig 配置结构
type RequestIDConfig struct {
    Enabled  bool   `json:"enabled" yaml:"enabled"`
    Generator string `json:"generator" yaml:"generator"` // uuid, snowflake
    Header   string `json:"header" yaml:"header"`
}

// DefaultRequestIDConfig 默认配置
func DefaultRequestIDConfig() *RequestIDConfig {
    return &RequestIDConfig{
        Enabled:   true,
        Generator: "uuid",
        Header:    "X-Request-ID",
    }
}
```

#### 2. 实现中间件

```go
// RequestIDMiddleware RequestID 中间件
type RequestIDMiddleware struct {
    config *RequestIDConfig
}

// NewRequestIDMiddleware 创建中间件
func NewRequestIDMiddleware() *RequestIDMiddleware {
    return &RequestIDMiddleware{
        config: DefaultRequestIDConfig(),
    }
}

// Handler 返回 gin.HandlerFunc
func (rm *RequestIDMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 如果未启用，跳过
        if !rm.config.Enabled {
            c.Next()
            return
        }

        // 尝试从请求头获取
        requestID := c.GetHeader(rm.config.Header)

        // 如果不存在，生成新的
        if requestID == "" {
            switch strings.ToLower(rm.config.Generator) {
            case "uuid":
                requestID = uuid.New().String()
            case "snowflake":
                // 使用 Snowflake 算法生成
                requestID = generateSnowflakeID()
            default:
                requestID = uuid.New().String()
            }
        }

        // 存储到上下文
        c.Set("request_id", requestID)

        // 设置响应头
        c.Header(rm.config.Header, requestID)

        c.Next()
    }
}

// LoadConfig 加载配置
func (rm *RequestIDMiddleware) LoadConfig(config map[string]interface{}) error {
    if enabled, ok := config["enabled"].(bool); ok {
        rm.config.Enabled = enabled
    }
    if generator, ok := config["generator"].(string); ok {
        rm.config.Generator = generator
    }
    if header, ok := config["header"].(string); ok {
        rm.config.Header = header
    }
    return nil
}

// ValidateConfig 验证配置
func (rm *RequestIDMiddleware) ValidateConfig() error {
    if rm.config.Header == "" {
        return errors.New("header 不能为空")
    }
    if rm.config.Generator != "uuid" && rm.config.Generator != "snowflake" {
        return errors.New("generator 必须是 uuid 或 snowflake")
    }
    return nil
}
```

#### 3. 注册中间件

```go
// internal/middleware/registry.go
package middleware

import (
    "qingyu/backend/internal/middleware/base"
)

func init() {
    // 创建 RequestID 中间件
    requestID := base.NewRequestIDMiddleware()

    // 注册到注册表
    registry.Register("request_id", MiddlewareConfig{
        Name:     "request_id",
        Priority: 900, // 高优先级，在 Recovery 之后
        Handler:  requestID.Handler(),
        Enabled:  true,
        Static:   true,
        Config:   requestID.config,
        Dependencies: []string{"recovery"},
    })
}
```

#### 4. 配置文件

```yaml
# config/middleware.yaml
static:
  request_id:
    enabled: true
    generator: "uuid"
    header: "X-Request-ID"
```

#### 5. 使用示例

```go
// 在 handler 中使用
func GetUserHandler(c *gin.Context) {
    // 获取 request_id
    requestID := c.GetString("request_id")

    // 记录日志
    log.Printf("Request ID: %s", requestID)

    // 返回响应
    c.JSON(200, gin.H{
        "request_id": requestID,
        "user":       currentUser,
    })
}
```

---

## 示例2: 创建动态配置中间件

### 场景

创建一个 `RateLimit` 限流中间件，支持动态更新限流规则。

### 完整代码

#### 1. 定义配置结构体

```go
// internal/middleware/security/rate_limit.go
package security

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
    "sync"
)

// RateLimitConfig 配置结构
type RateLimitConfig struct {
    Enabled bool   `json:"enabled" yaml:"enabled"`
    Driver  string `json:"driver" yaml:"driver"` // memory, redis

    // 全局限速
    Global struct {
        Rate  float64 `json:"rate" yaml:"rate"`   // 每秒请求数
        Burst int     `json:"burst" yaml:"burst"` // 突发容量
    } `json:"global" yaml:"global"`

    // 按 IP 限速
    PerIP struct {
        Rate  float64 `json:"rate" yaml:"rate"`
        Burst int     `json:"burst" yaml:"burst"`
    } `json:"per_ip" yaml:"per_ip"`

    // 按用户限速
    PerUser struct {
        Rate  float64 `json:"rate" yaml:"rate"`
        Burst int     `json:"burst" yaml:"burst"`
    } `json:"per_user" yaml:"per_user"`
}

// DefaultRateLimitConfig 默认配置
func DefaultRateLimitConfig() *RateLimitConfig {
    config := &RateLimitConfig{
        Enabled: true,
        Driver:  "memory",
    }
    config.Global.Rate = 1000
    config.Global.Burst = 2000
    config.PerIP.Rate = 100
    config.PerIP.Burst = 200
    config.PerUser.Rate = 50
    config.PerUser.Burst = 100
    return config
}
```

#### 2. 实现存储接口

```go
// RateLimitStore 限流存储接口
type RateLimitStore interface {
    Allow(key string, rate float64, burst int) bool
    Reset(key string) error
}

// MemoryRateLimitStore 内存实现
type MemoryRateLimitStore struct {
    limiters sync.Map // map[string]*rate.Limiter
    mu       sync.RWMutex
}

func NewMemoryRateLimitStore() *MemoryRateLimitStore {
    return &MemoryRateLimitStore{}
}

func (m *MemoryRateLimitStore) Allow(key string, rate float64, burst int) bool {
    // 获取或创建 limiter
    limiter, _ := m.limiters.LoadOrStore(key, rate.NewLimiter(rate.Limit(rate), burst))

    // 检查是否允许
    return limiter.(*rate.Limiter).Allow()
}

func (m *MemoryRateLimitStore) Reset(key string) error {
    m.limiters.Delete(key)
    return nil
}

// RedisRateLimitStore Redis 实现（使用 Redis-cell 算法）
type RedisRateLimitStore struct {
    client *redis.Client
}

func NewRedisRateLimitStore(client *redis.Client) *RedisRateLimitStore {
    return &RedisRateLimitStore{client: client}
}

func (r *RedisRateLimitStore) Allow(key string, rate float64, burst int) bool {
    // 使用 Redis CL.THROTTLE 命令
    result, err := r.client.Do(
        "CL.THROTTLE",
        key,
        rate,    // 令牌速率
        burst,   // 容量
        burst,   // 初始容量
        1,       // 获取1个令牌
    ).Result()

    if err != nil {
        return true // 出错时允许通过
    }

    // 解析结果：[状态, 限流重置时间, 获取令牌数, ...]
    status := result.([]interface{})[0]
    return status == int64(0) // 0 表示允许
}

func (r *RedisRateLimitStore) Reset(key string) error {
    return r.client.Del(key).Err()
}
```

#### 3. 实现中间件

```go
// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
    config *RateLimitConfig
    store  RateLimitStore
}

// NewRateLimitMiddleware 创建中间件
func NewRateLimitMiddleware(store RateLimitStore) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        config: DefaultRateLimitConfig(),
        store:  store,
    }
}

// Handler 返回 gin.HandlerFunc
func (rl *RateLimitMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !rl.config.Enabled {
            c.Next()
            return
        }

        // 1. 全局限流检查
        if !rl.store.Allow("global", rl.config.Global.Rate, rl.config.Global.Burst) {
            c.JSON(429, gin.H{
                "error": "全局请求过多，请稍后再试",
                "code":  "RATE_LIMIT_EXCEEDED",
            })
            c.Abort()
            return
        }

        // 2. IP 限流检查
        ip := c.ClientIP()
        key := "ip:" + ip
        if !rl.store.Allow(key, rl.config.PerIP.Rate, rl.config.PerIP.Burst) {
            c.JSON(429, gin.H{
                "error": "您的请求过于频繁，请稍后再试",
                "code":  "IP_RATE_LIMIT_EXCEEDED",
            })
            c.Abort()
            return
        }

        // 3. 用户限流检查（如果已登录）
        if userID, exists := c.Get("user_id"); exists {
            key := "user:" + userID.(string)
            if !rl.store.Allow(key, rl.config.PerUser.Rate, rl.config.PerUser.Burst) {
                c.JSON(429, gin.H{
                    "error": "您的请求过于频繁，请稍后再试",
                    "code":  "USER_RATE_LIMIT_EXCEEDED",
                })
                c.Abort()
                return
            }
        }

        c.Next()
    }
}

// LoadConfig 加载配置
func (rl *RateLimitMiddleware) LoadConfig(config map[string]interface{}) error {
    if enabled, ok := config["enabled"].(bool); ok {
        rl.config.Enabled = enabled
    }
    if driver, ok := config["driver"].(string); ok {
        rl.config.Driver = driver
    }

    // 加载全局配置
    if global, ok := config["global"].(map[string]interface{}); ok {
        if rate, ok := global["rate"].(float64); ok {
            rl.config.Global.Rate = rate
        }
        if burst, ok := global["burst"].(int); ok {
            rl.config.Global.Burst = burst
        }
    }

    // 加载 IP 配置
    if perIP, ok := config["per_ip"].(map[string]interface{}); ok {
        if rate, ok := perIP["rate"].(float64); ok {
            rl.config.PerIP.Rate = rate
        }
        if burst, ok := perIP["burst"].(int); ok {
            rl.config.PerIP.Burst = burst
        }
    }

    // 加载用户配置
    if perUser, ok := config["per_user"].(map[string]interface{}); ok {
        if rate, ok := perUser["rate"].(float64); ok {
            rl.config.PerUser.Rate = rate
        }
        if burst, ok := perUser["burst"].(int); ok {
            rl.config.PerUser.Burst = burst
        }
    }

    return nil
}

// ValidateConfig 验证配置
func (rl *RateLimitMiddleware) ValidateConfig() error {
    if rl.config.Global.Rate <= 0 {
        return errors.New("global rate 必须大于 0")
    }
    if rl.config.Global.Burst <= 0 {
        return errors.New("global burst 必须大于 0")
    }
    return nil
}
```

#### 4. 热更新配置

```go
// WatchConfigChanges 监听配置变化
func (rl *RateLimitMiddleware) WatchConfigChanges(configPath string) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }

    err = watcher.Add(configPath)
    if err != nil {
        return err
    }

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    // 重新加载配置
                    data, _ := ioutil.ReadFile(configPath)
                    config := make(map[string]interface{})
                    yaml.Unmarshal(data, &config)

                    dynamicConfig := config["dynamic"].(map[string]interface{})
                    rateLimitConfig := dynamicConfig["rate_limit"].(map[string]interface{})

                    rl.LoadConfig(rateLimitConfig)
                    log.Info("RateLimit 配置已热更新")
                }
            case err := <-watcher.Errors:
                log.Error("Watcher error:", err)
            }
        }
    }()

    return nil
}
```

#### 5. 配置文件

```yaml
# config/middleware.yaml
dynamic:
  rate_limit:
    enabled: true
    driver: "redis"

    # 全局限速
    global:
      rate: 1000   # 每秒 1000 个请求
      burst: 2000  # 突发容量 2000

    # 按 IP 限速
    per_ip:
      rate: 100    # 每 IP 每秒 100 个请求
      burst: 200

    # 按用户限速
    per_user:
      rate: 50     # 每用户每秒 50 个请求
      burst: 100
```

---

## 配置文件示例

### 完整配置文件

```yaml
# config/middleware.yaml

# ========== 静态配置中间件 ==========
static:
  # Recovery 中间件
  recovery:
    enabled: true
    stack: true
    stack_size: 4096  # 4KB

  # RequestID 中间件
  request_id:
    enabled: true
    generator: "uuid"    # uuid, snowflake
    header: "X-Request-ID"

  # CORS 中间件
  cors:
    enabled: true
    allow_origins:
      - "http://localhost:3000"
      - "https://qingyu.example.com"
    allow_methods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
      - "OPTIONS"
    allow_headers:
      - "Content-Type"
      - "Authorization"
      - "X-Request-ID"
    expose_headers:
      - "X-Request-ID"
    max_age: 86400       # 24小时
    allow_credentials: true

  # Timeout 中间件
  timeout:
    enabled: true
    timeout: 30s         # 30秒超时
    message: "请求超时"

# ========== 动态配置中间件 ==========
dynamic:
  # JWT 认证中间件
  jwt:
    enabled: true
    secret: "${JWT_SECRET}"
    expire: 86400        # 24小时
    issuer: "qingyu"
    # 白名单路径
    whitelist:
      - "/api/auth/login"
      - "/api/auth/register"
      - "/api/public/*"

  # RBAC 权限中间件
  rbac:
    enabled: true
    driver: "redis"      # memory, redis, database
    cache_ttl: 300       # 缓存5分钟
    # 权限规则
    rules:
      - role: "admin"
        resources:
          - pattern: "/api/admin/*"
            actions: ["*"]
      - role: "user"
        resources:
          - pattern: "/api/user/profile"
            actions: ["read", "update"]
          - pattern: "/api/books/*"
            actions: ["read"]

  # 限流中间件
  rate_limit:
    enabled: true
    driver: "redis"
    # 全局限速
    global:
      rate: 1000
      burst: 2000
    # 按 IP 限速
    per_ip:
      rate: 100
      burst: 200
    # 按用户限速
    per_user:
      rate: 50
      burst: 100

  # CSRF 防护中间件
  csrf:
    enabled: false       # 暂时禁用
    secret: "${CSRF_SECRET}"
    expire: 3600         # 1小时
    # 白名单方法
    whitelist_methods:
      - "GET"
      - "HEAD"
      - "OPTIONS"
    # 白名单路径
    whitelist_paths:
      - "/api/public/*"

  # 访问日志中间件
  access_log:
    enabled: true
    # 日志格式
    format: "json"       # json, text
    # 包含字段
    fields:
      - "request_id"
      - "method"
      - "path"
      - "status"
      - "latency"
      - "ip"
      - "user_agent"
    # 跳过路径
    skip_paths:
      - "/health"
      - "/metrics"

  # Prometheus 监控中间件
  prometheus:
    enabled: true
    # 指标路径
    path: "/metrics"
    # 指标标签
    labels:
      - "method"
      - "path"
      - "status"

  # 缓存中间件
  cache:
    enabled: true
    driver: "redis"      # memory, redis
    ttl: 300             # 默认缓存5分钟
    # 缓存规则
    rules:
      - pattern: "/api/books/*"
        ttl: 600         # 书籍列表缓存10分钟
      - pattern: "/api/categories/*"
        ttl: 3600        # 分类缓存1小时

  # 压缩中间件
  compress:
    enabled: true
    # 压缩级别
    level: 5             # 1-9，5是平衡点
    # 压缩类型
    types:
      - "application/json"
      - "text/html"
      - "text/plain"
      - "text/css"
      - "application/javascript"
    # 最小大小
    min_length: 1024     # 大于1KB才压缩
```

### 环境变量替换

配置文件支持环境变量替换：

```yaml
jwt:
  secret: "${JWT_SECRET}"        # 从环境变量读取
  expire: "${JWT_EXPIRE:86400}"  # 默认值86400

csrf:
  secret: "${CSRF_SECRET}"
```

使用方法：

```bash
# 设置环境变量
export JWT_SECRET="your-secret-key"
export CSRF_SECRET="your-csrf-secret"

# 启动应用
./qingyu-backend
```

---

## 单元测试示例

### 测试 RequestID 中间件

```go
// internal/middleware/base/request_id_test.go
package base

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
    // 设置 Gin 为测试模式
    gin.SetMode(gin.TestMode)

    t.Run("启用状态", func(t *testing.T) {
        // 创建中间件
        mw := NewRequestIDMiddleware()
        mw.config.Enabled = true

        // 创建测试路由
        router := gin.New()
        router.Use(mw.Handler())
        router.GET("/test", func(c *gin.Context) {
            requestID := c.GetString("request_id")
            c.JSON(200, gin.H{"request_id": requestID})
        })

        // 执行请求
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/test", nil)
        router.ServeHTTP(w, req)

        // 验证
        assert.Equal(t, 200, w.Code)
        assert.NotEmpty(t, w.Header().Get("X-Request-ID"))

        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.NotEmpty(t, response["request_id"])
    })

    t.Run("禁用状态", func(t *testing.T) {
        // 创建中间件
        mw := NewRequestIDMiddleware()
        mw.config.Enabled = false

        // 创建测试路由
        router := gin.New()
        router.Use(mw.Handler())
        router.GET("/test", func(c *gin.Context) {
            requestID := c.GetString("request_id")
            c.JSON(200, gin.H{"request_id": requestID})
        })

        // 执行请求
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/test", nil)
        router.ServeHTTP(w, req)

        // 验证
        assert.Equal(t, 200, w.Code)
        assert.Empty(t, w.Header().Get("X-Request-ID"))

        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.Empty(t, response["request_id"])
    })

    t.Run("使用请求头中的 RequestID", func(t *testing.T) {
        // 创建中间件
        mw := NewRequestIDMiddleware()
        mw.config.Enabled = true

        // 创建测试路由
        router := gin.New()
        router.Use(mw.Handler())
        router.GET("/test", func(c *gin.Context) {
            requestID := c.GetString("request_id")
            c.JSON(200, gin.H{"request_id": requestID})
        })

        // 执行请求（带 RequestID）
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/test", nil)
        req.Header.Set("X-Request-ID", "test-request-id-123")
        router.ServeHTTP(w, req)

        // 验证
        assert.Equal(t, 200, w.Code)
        assert.Equal(t, "test-request-id-123", w.Header().Get("X-Request-ID"))

        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.Equal(t, "test-request-id-123", response["request_id"])
    })
}
```

### 测试 RateLimit 中间件

```go
// internal/middleware/security/rate_limit_test.go
package security

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("全局限流", func(t *testing.T) {
        // 创建内存存储
        store := NewMemoryRateLimitStore()

        // 创建中间件（设置极低限流）
        mw := NewRateLimitMiddleware(store)
        mw.config.Enabled = true
        mw.config.Global.Rate = 1   // 每秒1个请求
        mw.config.Global.Burst = 1  // 容量1

        // 创建测试路由
        router := gin.New()
        router.Use(mw.Handler())
        router.GET("/test", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "ok"})
        })

        // 第一个请求应该成功
        w1 := httptest.NewRecorder()
        req1, _ := http.NewRequest("GET", "/test", nil)
        router.ServeHTTP(w1, req1)
        assert.Equal(t, 200, w1.Code)

        // 第二个请求应该被限流
        w2 := httptest.NewRecorder()
        req2, _ := http.NewRequest("GET", "/test", nil)
        router.ServeHTTP(w2, req2)
        assert.Equal(t, 429, w2.Code)
    })

    t.Run("IP 限流", func(t *testing.T) {
        store := NewMemoryRateLimitStore()

        mw := NewRateLimitMiddleware(store)
        mw.config.Enabled = true
        mw.config.PerIP.Rate = 1
        mw.config.PerIP.Burst = 1

        router := gin.New()
        router.Use(mw.Handler())
        router.GET("/test", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "ok"})
        })

        // 第一个 IP 的第一个请求成功
        w1 := httptest.NewRecorder()
        req1, _ := http.NewRequest("GET", "/test", nil)
        req1.RemoteAddr = "192.168.1.100:1234"
        router.ServeHTTP(w1, req1)
        assert.Equal(t, 200, w1.Code)

        // 同一个 IP 的第二个请求被限流
        w2 := httptest.NewRecorder()
        req2, _ := http.NewRequest("GET", "/test", nil)
        req2.RemoteAddr = "192.168.1.100:1234"
        router.ServeHTTP(w2, req2)
        assert.Equal(t, 429, w2.Code)

        // 不同 IP 的请求成功
        w3 := httptest.NewRecorder()
        req3, _ := http.NewRequest("GET", "/test", nil)
        req3.RemoteAddr = "192.168.1.101:1234"
        router.ServeHTTP(w3, req3)
        assert.Equal(t, 200, w3.Code)
    })

    t.Run("禁用状态", func(t *testing.T) {
        store := NewMemoryRateLimitStore()

        mw := NewRateLimitMiddleware(store)
        mw.config.Enabled = false  // 禁用

        router := gin.New()
        router.Use(mw.Handler())
        router.GET("/test", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "ok"})
        })

        // 多个请求都应该成功
        for i := 0; i < 10; i++ {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/test", nil)
            router.ServeHTTP(w, req)
            assert.Equal(t, 200, w.Code)
        }
    })
}
```

### 表驱动测试

```go
// internal/middleware/base/request_id_validation_test.go
package base

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRequestIDConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  *RequestIDConfig
        wantErr bool
        errMsg  string
    }{
        {
            name: "有效配置 - UUID",
            config: &RequestIDConfig{
                Enabled:   true,
                Generator: "uuid",
                Header:    "X-Request-ID",
            },
            wantErr: false,
        },
        {
            name: "有效配置 - Snowflake",
            config: &RequestIDConfig{
                Enabled:   true,
                Generator: "snowflake",
                Header:    "X-Request-ID",
            },
            wantErr: false,
        },
        {
            name: "无效配置 - 空 header",
            config: &RequestIDConfig{
                Enabled:   true,
                Generator: "uuid",
                Header:    "",
            },
            wantErr: true,
            errMsg:  "header 不能为空",
        },
        {
            name: "无效配置 - 无效 generator",
            config: &RequestIDConfig{
                Enabled:   true,
                Generator: "invalid",
                Header:    "X-Request-ID",
            },
            wantErr: true,
            errMsg:  "generator 必须是 uuid 或 snowflake",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mw := &RequestIDMiddleware{config: tt.config}
            err := mw.ValidateConfig()

            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## 常见问题

### Q1: 如何禁用某个中间件？

**A**: 在配置文件中设置 `enabled: false`：

```yaml
static:
  cors:
    enabled: false  # 禁用 CORS
```

或运行时禁用：

```go
manager.Disable("cors")
```

### Q2: 如何自定义错误响应格式？

**A**: 实现自定义 `ErrorHandler`：

```go
func CustomErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors[0].Err

            // 自定义格式
            c.JSON(500, gin.H{
                "success": false,
                "error": gin.H{
                    "code":    getErrorCode(err),
                    "message": err.Error(),
                    "request_id": c.GetString("request_id"),
                },
            })
        }
    }
}
```

### Q3: 如何添加中间件到特定路由？

**A**: 使用路由分组：

```go
// 仅对 admin 路由应用 RBAC
adminGroup := router.Group("/api/admin")
adminGroup.Use(rbac.Handler())
{
    adminGroup.GET("/users", GetUsersHandler)
    adminGroup.POST("/users", CreateUserHandler)
}

// 公开路由不需要认证
publicGroup := router.Group("/api/public")
{
    publicGroup.GET("/books", GetBooksHandler)
}
```

### Q4: 如何调试中间件执行顺序？

**A**: 添加日志中间件：

```go
func DebugMiddleware(name string) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("[DEBUG] %s: Before request", name)
        c.Next()
        log.Printf("[DEBUG] %s: After request (status=%d)", name, c.Writer.Status())
    }
}

// 使用
router.Use(
    DebugMiddleware("Recovery"),
    recovery.Handler(),
    DebugMiddleware("RequestID"),
    requestID.Handler(),
)
```

### Q5: 如何测试中间件？

**A**: 使用 `httptest` 包：

```go
func TestMyMiddleware(t *testing.T) {
    // 创建测试路由
    router := gin.New()
    router.Use(myMiddleware.Handler())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // 执行请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    router.ServeHTTP(w, req)

    // 验证
    assert.Equal(t, 200, w.Code)
}
```

### Q6: 如何实现中间件间通信？

**A**: 使用 Gin 的 `Set` 和 `Get`：

```go
// 第一个中间件设置数据
func FirstMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("user_id", "12345")
        c.Next()
    }
}

// 第二个中间件获取数据
func SecondMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        log.Printf("User ID: %s", userID)
        c.Next()
    }
}
```

### Q7: 如何处理中间件中的异步操作？

**A**: 使用 goroutine 和 context：

```go
func AsyncMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 在后台执行异步任务
        go func() {
            // 使用请求的副本
            requestID := c.GetString("request_id")

            // 执行异步操作
            result := doAsyncWork()

            // 记录结果
            log.Printf("Request %s: async result = %v", requestID, result)
        }()

        c.Next()
    }
}
```

### Q8: 如何实现中间件链超时？

**A**: 使用 `context.WithTimeout`：

```go
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()

        c.Request = c.Request.WithContext(ctx)

        finished := make(chan struct{})
        go func() {
            c.Next()
            close(finished)
        }()

        select {
        case <-finished:
            // 正常完成
        case <-ctx.Done():
            // 超时
            c.JSON(408, gin.H{"error": "请求超时"})
            c.Abort()
        }
    }
}
```

---

## 最佳实践

### 1. 中间件顺序

```go
// 推荐顺序
router.Use(
    recovery.Handler(),       // 1. 恢复 panic
    requestID.Handler(),      // 2. 生成请求 ID
    timeout.Handler(),        // 3. 超时控制
    cors.Handler(),           // 4. CORS
    jwt.Handler(),            // 5. 认证
    rbac.Handler(),           // 6. 授权
    rateLimit.Handler(),      // 7. 限流
    accessLog.Handler(),      // 8. 访问日志
    prometheus.Handler(),     // 9. 监控
    cache.Handler(),          // 10. 缓存
    compress.Handler(),       // 11. 压缩
    errorHandler.Handler(),   // 12. 错误处理
)
```

### 2. 错误处理

```go
// 不要在中间件中直接 panic
if err != nil {
    panic(err)  // ❌ 错误
}

// 应该使用 c.Error()
if err != nil {
    c.Error(err)  // ✅ 正确
    return
}
```

### 3. 性能优化

```go
// 避免重复计算
func MyMiddleware() gin.HandlerFunc {
    // 预计算
    config := loadConfig()

    return func(c *gin.Context) {
        // 使用预计算的配置
        if config.Enabled {
            // ...
        }
        c.Next()
    }
}
```

### 4. 测试覆盖

```go
// 测试所有场景
- 中间件启用
- 中间件禁用
- 正常请求
- 错误请求
- 边界条件
- 并发请求
```

---

## 下一步

- 阅读[架构培训文档](./middleware-architecture-training.md)
- 查看[API 文档](../api/middleware-api.md)
- 运行[示例代码](../examples/)

---

**文档版本**: v1.0
**最后更新**: 2026-01-27
**维护者**: 架构师团队
