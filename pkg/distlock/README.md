# 分布式锁模块 (distlock)

## 概述

分布式锁模块提供基于 Redis 的分布式锁实现，用于在分布式环境中保证共享资源的互斥访问。

## 核心组件

### RedisLockService

基于 Redis 的分布式锁服务实现。

#### 特性

- **原子性获取锁**: 使用 `SETNX` 实现，确保只有一个客户端能成功获取锁
- **原子性释放锁**: 使用 Lua 脚本确保只有锁的持有者才能释放锁，防止误删
- **原子性延期**: 使用 Lua 脚本实现锁的自动延期
- **可重试机制**: 支持设置最大重试次数和重试间隔

#### 锁的获取 (Acquire)

```
SET key value NX PX ttl
```

- `NX`: 仅在键不存在时设置
- `PX ttl`: 设置毫秒级过期时间

#### 锁的释放 (Release)

使用 Lua 脚本保证原子性：

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
```

只有当锁的值等于请求者持有的 `lockID` 时才删除，防止误删其他客户端的锁。

#### 锁的延期 (Extend)

使用 Lua 脚本保证原子性：

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("pexpire", KEYS[1], ARGV[2])
else
    return 0
end
```

只有锁的持有者才能延长锁的过期时间。

## 使用方式

### 1. 通过服务容器创建（推荐）

在 `router/writer/writer_router_init.go` 中：

```go
import "Qingyu_backend/pkg/distlock"

// 获取 Redis 客户端
redisClient := serviceContainer.GetRedisClient()

// 创建分布式锁服务
var distLockSvc *distlock.RedisLockService
if redisClient != nil {
    distLockSvc = distlock.NewRedisLockService(redisClient, "distlock")
}

// 注入到需要使用锁的服务
syncSvc := writerservice.NewOutlineDocumentSyncService(
    outlineRepo, documentRepo, projectRepo,
    outlineSvc.(*writerservice.OutlineService),
    distLockSvc,
)
```

### 2. 直接通过地址创建

```go
distLockSvc, err := distlock.NewRedisLockServiceFromAddr(
    "localhost:6379",    // 地址
    "",                  // 密码（无密码则空字符串）
    0,                   // DB 编号
    "myapp",             // 键前缀
)
if err != nil {
    log.Fatalf("创建分布式锁服务失败: %v", err)
}
defer distLockSvc.Close()
```

## API 接口

### Lock 接口

```go
type Lock interface {
    Acquire(ctx context.Context, ttl time.Duration) (string, error)
    Release(ctx context.Context, lockID string) error
    Extend(ctx context.Context, lockID string, ttl time.Duration) error
}
```

### RedisLockService 方法

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `Acquire` | 尝试获取锁 | `ctx`: 上下文, `lockKey`: 锁键, `ttl`: 过期时间 | `lockID`: 锁标识（用于释放）, `error`: 错误 |
| `AcquireWithRetry` | 获取锁，带重试 | `ctx`: 上下文, `lockKey`: 锁键, `ttl`: 过期时间, `maxRetries`: 最大重试次数, `retryInterval`: 重试间隔 | `lockID`: 锁标识, `error`: 错误 |
| `Release` | 释放锁 | `ctx`: 上下文, `lockKey`: 锁键, `lockID`: 锁标识 | `error`: 错误 |
| `Extend` | 延长锁的过期时间 | `ctx`: 上下文, `lockKey`: 锁键, `lockID`: 锁标识, `ttl`: 新的过期时间 | `error`: 错误 |
| `Close` | 关闭 Redis 客户端 | - | `error`: 错误 |

## 应用场景

### 1. 保护全局总纲的并发创建

在 `OutlineDocumentSyncService.findOrCreateGlobalOutline` 中，使用分布式锁防止并发创建多个全局总纲：

```go
func (s *OutlineDocumentSyncService) findOrCreateGlobalOutline(ctx context.Context, projectID string) (string, error) {
    if s.distLock == nil {
        return s.findOrCreateGlobalOutlineWithUpsert(ctx, projectID)
    }

    // 使用分布式锁保证并发安全
    lockKey := fmt.Sprintf("global_outline:%s", projectID)
    lockID, err := s.distLock.AcquireWithRetry(ctx, lockKey, 5*time.Second, 3, 500*time.Millisecond)
    if err != nil {
        log.Printf("[OutlineDocSync] 获取分布式锁失败: %v", err)
        return s.findOrCreateGlobalOutlineWithUpsert(ctx, projectID)
    }
    defer func() {
        if releaseErr := s.distLock.Release(ctx, lockKey, lockID); releaseErr != nil {
            log.Printf("[OutlineDocSync] 释放分布式锁失败: %v", releaseErr)
        }
    }()

    // 临界区内使用原子性 upsert
    projectOID, _ := primitive.ObjectIDFromHex(projectID)
    existing, err := s.outlineRepo.FindByGlobalOutline(ctx, projectOID)
    // ...
}
```

### 2. 其他需要分布式互斥的场景

- 分布式任务调度
- 共享资源的互斥访问
- 防止重复操作（如重复订单、重复消息处理等）

## 错误处理

| 错误 | 说明 |
|------|------|
| `ErrLockAcquisitionFailed` | 获取锁失败（锁已被其他客户端持有） |
| `ErrLockNotHeld` | 锁未持有（释放或延期时锁已过期或被其他客户端持有） |
| `ErrLockTimeout` | 获取锁超时（超过最大重试次数） |

## 注意事项

1. **锁的过期时间**: 建议设置合理的 TTL，既要防止锁永久持有（客户端崩溃），也要确保业务操作有足够时间完成
2. **锁的键命名**: 使用有意义的键名前缀，如 `distlock:global_outline:projectID`
3. **释放锁的 defer**: 在获取锁后应立即使用 defer 释放锁，确保即使业务逻辑失败也能释放锁
4. **降级处理**: 如果分布式锁服务不可用，应有降级方案（如纯数据库原子操作）
5. **客户端崩溃**: 如果持有锁的客户端崩溃，锁会自动过期，其他客户端可以继续获取锁
