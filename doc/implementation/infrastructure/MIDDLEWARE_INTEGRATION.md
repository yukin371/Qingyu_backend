# P0 中间件集成实施文档

## 概述

本文档记录了青羽平台 P0 (Priority 0) 核心中间件的集成过程。

**实施日期**: 2026-01-04
**状态**: ✅ 已完成

## 什么是 P0 中间件

P0 中间件是指系统中**最高优先级**的基础中间件，它们是系统运行的基础设施，必须首先集成和正确配置。

### P0 中间件列表

1. **RequestID** - 请求 ID 生成和追踪
2. **Recovery** - Panic 恢复和错误处理
3. **Logger** - 结构化日志记录
4. **Metrics** - Prometheus 监控指标
5. **RateLimit** - 请求频率限制
6. **ErrorHandler** - 统一错误处理
7. **CORS** - 跨域资源共享

## 架构设计

### 中间件执行顺序

中间件的执行顺序至关重要，正确的顺序应该是：

```
请求 → RequestID → Recovery → Logger → Metrics → RateLimit → ErrorHandler → CORS → 路由处理
```

### 执行流程

```
┌──────────────────────────────────────────────────────────────┐
│                        Request                               │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 1. RequestID - 生成唯一请求 ID                                │
│    - X-Request-ID header                                     │
│    - 传递给下游处理                                          │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 2. Recovery - 捕获 panic                                      │
│    - 防止程序崩溃                                            │
│    - 记录 panic 堆栈                                         │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 3. Logger - 记录请求日志                                      │
│    - 请求路径、方法                                          │
│    - 请求 ID、客户端 IP                                      │
│    - 执行时长                                                │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 4. Metrics - 收集监控指标                                     │
│    - 请求计数                                                │
│    - 响应时间                                                │
│    - 状态码分布                                              │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 5. RateLimit - 频率限制检查                                   │
│    - 检查请求频率                                            │
│    - 返回 429 Too Many Requests                              │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 6. ErrorHandler - 统一错误处理                                │
│    - 捕获和处理错误                                          │
│    - 返回统一错误响应                                        │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ 7. CORS - 跨域处理                                           │
│    - 处理 OPTIONS 请求                                       │
│    - 添加 CORS headers                                       │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│                    Route Handler                             │
└──────────────────────────────────────────────────────────────┘
```

## 实施步骤

**提交**: `9fe9a01 feat(core): 集成P0中间件到核心服务器`

### 步骤 1: 初始化 zap 日志器

```go
// core/server.go

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func initLogger() (*zap.Logger, error) {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    logger, err := config.Build(
        zap.AddCallerSkip(1),
        zap.AddStacktrace(zapcore.ErrorLevel),
    )
    if err != nil {
        return nil, err
    }

    return logger, nil
}
```

### 步骤 2: 按正确顺序应用中间件

```go
// core/server.go

func SetupServer() *gin.Engine {
    r := gin.New()

    // 1. RequestID 中间件
    r.Use(middleware.RequestIDMiddleware())

    // 2. Recovery 中间件
    r.Use(middleware.RecoveryMiddleware(logger))

    // 3. Logger 中间件
    r.Use(middleware.LoggerMiddleware(logger))

    // 4. Metrics 中间件
    r.Use(middleware.MetricsMiddleware())

    // 5. RateLimit 中间件
    r.Use(middleware.RateLimitMiddleware())

    // 6. ErrorHandler 中间件
    r.Use(middleware.ErrorHandlerMiddleware())

    // 7. CORS 中间件
    r.Use(middleware.CORSMiddleware())

    return r
}
```

### 步骤 3: 注册健康检查端点

```go
// 健康检查
r.GET("/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    })
})

// 存活检查 (Kubernetes)
r.GET("/health/live", func(c *gin.Context) {
    c.Status(http.StatusOK)
})

// 就绪检查 (Kubernetes)
r.GET("/health/ready", func(c *gin.Context) {
    c.Status(http.StatusOK)
})
```

### 步骤 4: 注册 Prometheus 监控端点

```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/segmentio/stats/v4/promstats"
)

// 注册 Prometheus 指标端点
r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

## 中间件实现

### 1. RequestID 中间件

```go
// middleware/request_id.go

func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查是否已有 Request ID
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }

        // 设置到 context 和 header
        c.Set("RequestID", requestID)
        c.Header("X-Request-ID", requestID)

        c.Next()
    }
}

func generateRequestID() string {
    return uuid.New().String() // 或使用其他生成方式
}
```

### 2. Recovery 中间件

```go
// middleware/recovery.go

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录 panic
                logger.Error("panic recovered",
                    zap.Any("error", err),
                    zap.String("request_id", c.GetString("RequestID")),
                    zap.String("path", c.Request.URL.Path),
                    zap.Stack("stack"),
                )

                // 返回 500 错误
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code":    500,
                    "message": "Internal Server Error",
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

### 3. Logger 中间件

```go
// middleware/logger.go

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        // 处理请求
        c.Next()

        // 记录日志
        latency := time.Since(start)
        status := c.Writer.Status()
        method := c.Request.Method
        ip := c.ClientIP()

        logger.Info("request",
            zap.String("request_id", c.GetString("RequestID")),
            zap.String("method", method),
            zap.String("path", path),
            zap.String("query", query),
            zap.Int("status", status),
            zap.Duration("latency", latency),
            zap.String("ip", ip),
            zap.String("user-agent", c.Request.UserAgent()),
        )
    }
}
```

### 4. Metrics 中间件

```go
// middleware/metrics.go

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // 处理请求
        c.Next()

        // 记录指标
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        method := c.Request.Method
        path := c.FullPath()

        httpRequestsTotal.WithLabelValues(method, path, status).Inc()
        httpRequestDuration.WithLabelValues(method, path).Observe(duration)
    }
}
```

### 5. RateLimit 中间件

```go
// middleware/rate_limit.go

import (
    "golang.org/x/time/rate"
)

func RateLimitMiddleware() gin.HandlerFunc {
    // 使用简单的内存限流器
    limiter := rate.NewLimiter(rate.Limit(100), 200) // 100 req/s, burst 200

    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "code":    429,
                "message": "Too Many Requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 6. ErrorHandler 中间件

```go
// middleware/error_handler.go

func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            logger.Error("request error",
                zap.Error(err.Err),
                zap.String("request_id", c.GetString("RequestID")),
            )

            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    500,
                "message": err.Error(),
            })
        }
    }
}
```

### 7. CORS 中间件

```go
// middleware/cors.go

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.Header("Access-Control-Expose-Headers", "Content-Length")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

## 端点列表

### 健康检查端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/health` | GET | 系统健康状态 |
| `/health/live` | GET | 存活检查 (Kubernetes liveness) |
| `/health/ready` | GET | 就绪检查 (Kubernetes readiness) |
| `/ping` | GET | 简单健康检查 |

### 监控端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/metrics` | GET | Prometheus 指标数据 |

## Prometheus 指标

### 收集的指标

| 指标名称 | 类型 | 说明 |
|----------|------|------|
| `http_requests_total` | Counter | HTTP 请求总数 |
| `http_request_duration_seconds` | Histogram | HTTP 请求耗时 |
| `panics_total` | Counter | Panic 总数 |

### 使用 PromQL 查询

```promql
# 请求速率
rate(http_requests_total[5m])

# P95 响应时间
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# 错误率
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])
```

## 配置建议

### 日志级别

- **开发环境**: `debug`
- **测试环境**: `info`
- **生产环境**: `warn` 或 `error`

### 频率限制

- **公开 API**: 100 req/s, burst 200
- **认证 API**: 200 req/s, burst 400
- **管理员 API**: 50 req/s, burst 100

### CORS 策略

- **开发环境**: 允许所有来源
- **生产环境**: 限制特定域名

## 监控和告警

### Grafana 仪表板

建议监控以下指标：

1. **请求速率** - QPS
2. **响应时间** - P50, P95, P99
3. **错误率** - 4xx, 5xx 比例
4. **Panic 次数** - 应该为 0

### 告警规则

```yaml
# prometheus/alerts.yaml
groups:
  - name: api
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        annotations:
          summary: "API 错误率过高"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        annotations:
          summary: "API P95 响应时间超过 1s"
```

## 测试验证

### 测试 Request ID

```bash
curl -i http://localhost:8080/api/v1/health
```

**预期响应头**:
```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

### 测试 Recovery

触发一个 panic，验证系统不会崩溃：

```bash
curl http://localhost:8080/api/v1/debug/panic
```

**预期行为**:
- 服务不会崩溃
- 返回 500 错误
- 日志记录 panic

### 测试 Rate Limit

发送大量请求，触发限流：

```bash
for i in {1..300}; do curl http://localhost:8080/api/v1/health; done
```

**预期行为**:
- 前 200 个请求成功
- 后续请求返回 429

### 测试 Metrics

```bash
curl http://localhost:8080/metrics
```

**预期输出**: Prometheus 指标数据

## 最佳实践

### 1. 中间件顺序

严格遵循中间件执行顺序，错误顺序可能导致功能异常。

### 2. 错误处理

所有中间件都应该正确处理错误，避免泄露敏感信息。

### 3. 性能考虑

- 避免在中间件中执行耗时操作
- 使用异步日志记录
- 合理设置采样率

### 4. 可观测性

- 所有中间件都应该记录日志
- 关键操作应该记录指标
- 使用 Request ID 追踪请求链路

## 相关文档

- [路由修复报告](ROUTER_FIX_REPORT.md)
- [RBAC 权限控制](RBAC_IMPLEMENTATION.md)
- [项目结构总结](../docs/项目结构总结.md)

## 提交历史

```
9fe9a01 - feat(core): 集成P0中间件到核心服务器
18f2e55 - feat(core): 集成P0中间件到核心服务器
```
