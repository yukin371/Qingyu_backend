# Shared Interfaces - 最小领域契约

## 概述

这个包定义了shared模块的最小领域契约（Ports），采用Port/Adapter架构模式，实现依赖倒置和模块解耦。

## 架构模式

### Port/Adapter模式

```
┌─────────────────────────────────────────────────────┐
│                   Application Layer                  │
│                  (业务逻辑层)                         │
└────────────────────┬────────────────────────────────┘
                     │ 依赖
                     ▼
┌─────────────────────────────────────────────────────┐
│                    Ports                             │
│          (接口定义 - 最小领域契约)                     │
│  • StoragePort  • CachePort  • AuthPort              │
└────────────────────┬────────────────────────────────┘
                     │ 实现
                     ▼
┌─────────────────────────────────────────────────────┐
│                   Adapters                           │
│          (适配器 - 连接具体实现)                       │
│  • StorageAdapter  • CacheAdapter  • AuthAdapter     │
└────────────────────┬────────────────────────────────┘
                     │ 包装
                     ▼
┌─────────────────────────────────────────────────────┐
│              Shared Module Implementations           │
│         (现有实现 - 无需修改)                          │
│  • storage.StorageService  • auth.AuthService        │
└─────────────────────────────────────────────────────┘
```

## 接口定义

### StoragePort - 文件存储接口

最小化的文件存储契约，只包含最基本的操作：

```go
type StoragePort interface {
    Upload(ctx context.Context, filename string, data []byte) (string, error)
    Download(ctx context.Context, fileID string) ([]byte, error)
    Delete(ctx context.Context, fileID string) error
}
```

**设计原则**：
- 只暴露核心功能：上传、下载、删除
- 使用简单的数据类型：[]byte、string
- 避免复杂的请求结构

### CachePort - 缓存接口

最小化的缓存契约：

```go
type CachePort interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

**设计原则**：
- 只包含基础操作：Get、Set、Delete
- 使用简单的string值类型
- 支持TTL过期时间

### AuthPort - 认证接口

最小化的认证契约：

```go
type AuthPort interface {
    ValidateToken(ctx context.Context, token string) (string, error)
    CheckPermission(ctx context.Context, userID, permission string) (bool, error)
}
```

**设计原则**：
- 只包含必要的认证和授权功能
- 返回简化的用户标识（userID）
- 权限检查使用简单的字符串匹配

## 适配器

### StorageAdapter

将现有的`storage.StorageService`适配到`StoragePort`接口：

```go
adapter := NewStorageAdapter(storageService)
fileID, err := adapter.Upload(ctx, "file.txt", data)
```

### CacheAdapter

将现有的缓存服务适配到`CachePort`接口：

```go
adapter := NewCacheAdapter(cacheService)
err := adapter.Set(ctx, "key", "value", time.Minute)
```

### AuthAdapter

将现有的`auth.AuthService`适配到`AuthPort`接口：

```go
adapter := NewAuthAdapter(authService)
userID, err := adapter.ValidateToken(ctx, token)
```

## 使用示例

### 在业务代码中使用Ports

```go
// 业务代码只依赖接口，不依赖具体实现
type BookService struct {
    storage shared.StoragePort
    cache   shared.CachePort
    auth    shared.AuthPort
}

func NewBookService(
    storage shared.StoragePort,
    cache shared.CachePort,
    auth shared.AuthPort,
) *BookService {
    return &BookService{
        storage: storage,
        cache:   cache,
        auth:    auth,
    }
}

func (s *BookService) UploadCover(ctx context.Context, data []byte) (string, error) {
    // 使用StoragePort接口，不关心具体实现
    return s.storage.Upload(ctx, "cover.jpg", data)
}
```

### 在容器中注册Adapters

```go
// 在服务容器中注册适配器
container.RegisterProvider(container.Provider{
    Name: "storage-port",
    Factory: func(c *container.ServiceContainer) (interface{}, error) {
        storageService := c.GetStorageService()
        return shared.NewStorageAdapter(storageService), nil
    },
})
```

## 优势

1. **依赖倒置**：业务代码依赖接口，不依赖具体实现
2. **易于测试**：可以使用Mock实现进行单元测试
3. **灵活替换**：可以轻松替换底层实现（如从Redis迁移到Memcached）
4. **最小化契约**：接口只包含必要的方法，降低耦合
5. **向后兼容**：不影响现有业务代码

## 向后兼容性

- 现有的`service/shared/*`模块无需修改
- 适配器层无缝桥接新旧接口
- 业务代码可以逐步迁移到新接口

## 下一步

- [ ] 添加更多Ports（如MessagingPort、NotificationPort）
- [ ] 为关键模块编写基于Ports的集成测试
- [ ] 逐步迁移业务代码使用Ports而非直接依赖shared模块
- [ ] 添加CI检查，禁止跨层直接依赖shared模块实现

## 文件结构

```
service/interfaces/shared/
├── ports.go           # Port接口定义
├── adapters.go        # Adapter实现
├── ports_test.go      # Port接口测试
├── adapters_test.go   # Adapter测试
└── README.md          # 本文档
```

## 测试覆盖

所有接口和适配器都有完整的测试覆盖：

```bash
go test ./service/interfaces/shared/ -v
```

测试结果：
- ✅ 接口定义测试（3个）
- ✅ 接口行为测试（9个）
- ✅ 适配器测试（8个）
- ✅ 类型实现测试（3个）

**总计：23个测试用例，全部通过**
