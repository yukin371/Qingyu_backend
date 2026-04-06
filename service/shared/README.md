# 共享底层服务模块

> 青羽后端 - 模块化单体架构的共享服务层
>
> **创建时间**: 2025-09-30
> **最后更新**: 2026-01-12 (清理重复模块)

---

## 📋 概述

共享底层服务是青羽后端的基础设施层，为阅读端和写作端提供统一的核心功能。

**架构模式**: 模块化单体（Modular Monolith）

**注意**: 以下模块已移至业务目录，不再在此目录维护：
- **Wallet** → `service/finance/wallet/`
- **Recommendation** → `service/recommendation/`
- **Admin** → `service/admin/`

---

## 🏗️ 模块结构

### 当前模块

| 模块 | 状态 | 说明 | 主要文件 |
|------|------|------|---------|
| **Auth** | 📦 已迁移 | 认证与权限管理 | 已迁移至 [service/auth/](../auth/) |
| **Messaging** | 📦 已迁移 | 消息队列与通知 | 已迁移至 [service/channels/](../channels/) |
| **Storage** | ✅ 接口 | 文件存储接口 | [interfaces.go](./storage/interfaces.go) |
| **Cache** | ✅ 实现 | Redis 缓存服务 | [redis_cache_service.go](./cache/redis_cache_service.go) |
| **Search** | ✅ 实现 | 搜索服务 | [search_service.go](./search/search_service.go) |
| **Metrics** | ✅ 实现 | 服务指标收集 | [service_metrics.go](./metrics/service_metrics.go) |
| **Stats** | ✅ 实现 | 平台统计服务 | [stats_service.go](./stats/stats_service.go) |

---

## 📦 目录结构

```
service/shared/
├── storage/                   # 文件存储模块
│   └── interfaces.go          # 存储服务接口
│
├── cache/                     # 缓存模块
│   └── redis_cache_service.go # Redis 缓存实现
│
├── search/                    # 搜索模块
│   └── search_service.go      # 搜索服务实现
│
├── metrics/                   # 指标模块
│   └── service_metrics.go     # 服务指标收集
│
├── stats/                     # 统计模块
│   └── stats_service.go       # 平台统计服务
│
├── config_service.go          # 动态配置管理服务
├── permission_service.go      # 权限服务（完整 RBAC）
├── messaging_compat.go        # 消息模块兼容层
│
└── README.md                  # 本文档
```

**已迁移模块**:
- **auth/** → 已迁移至 [`service/auth/`](../auth/)
- **messaging/** → 已迁移至 [`service/channels/`](../channels/)

---

## 🎯 各模块功能

### 已迁移模块

#### 1. Auth 模块 (已迁移至 `service/auth/`)

**核心功能**:
- JWT Token 生成与验证
- 用户注册与登录
- 基于 RBAC 的权限系统
- 会话管理（Redis）
- 角色与权限管理

**新位置**: [`service/auth/`](../auth/)

#### 2. Messaging 模块 (已迁移至 `service/channels/`)

**核心功能**:
- 消息发布订阅（基于 Redis Streams）
- 邮件发送服务
- 通知发送服务
- 延迟消息支持

**新位置**: [`service/channels/`](../channels/)

**兼容层**: 本目录的 `messaging_compat.go` 提供了向后兼容的别名和类型

---

### 1. Storage 模块

**核心功能**:
- 文件上传下载接口定义
- 文件权限控制
- 文件元数据管理
- 临时访问 URL 生成

**主要文件**:
- `storage/interfaces.go` - 存储服务接口定义

---

### 2. Cache 模块

**核心功能**:
- 通用 Redis 缓存服务
- 基础缓存操作（GET/SET/DELETE）
- 批量操作支持
- 哈希、集合、有序集合操作

**主要文件**:
- `cache/redis_cache_service.go` - Redis 缓存实现
- `cache/redis_cache_service_test.go` - 基础操作测试
- `cache/redis_cache_advanced_test.go` - 批量操作、高级操作、哈希和集合测试
- `cache/redis_cache_sorted_set_test.go` - 有序集合和服务管理测试

**测试覆盖**:
- 测试覆盖率: 82.9%
- 测试用例数: 60+
- 测试框架: miniredis + testify
- 详细测试报告: 参见 `docs/plans/submodules/backend/shared-and-layering/shared-cache-tdd-plan.md`

---

### 3. Search 模块

**核心功能**:
- 搜索服务接口定义

**主要文件**:
- `search/search_service.go` - 搜索服务实现

---

### 4. Metrics 模块

**核心功能**:
- 服务指标收集与上报

**主要文件**:
- `metrics/service_metrics.go` - 服务指标收集

---

### 5. Stats 模块

**核心功能**:
- 平台统计数据服务

**主要文件**:
- `stats/stats_service.go` - 平台统计服务

---

### 6. Config 服务

**核心功能**:
- 动态配置管理
- 配置热更新
- 配置备份与恢复

**主要文件**:
- `config_service.go` - 配置管理服务

---

### 7. Permission 服务

**核心功能**:
- 完整 RBAC 权限实现
- 权限检查接口
- 角色管理接口

**主要文件**:
- `permission_service.go` - 权限服务实现

---

## 🔧 使用说明

### 导入服务接口

```go
import (
    authService "Qingyu_backend/service/auth"           // 认证服务
    cacheService "Qingyu_backend/service/shared/cache"  // 缓存服务
    channelsService "Qingyu_backend/service/channels"   // 消息通知服务
    storageService "Qingyu_backend/service/shared/storage" // 存储服务
)
```

**注意**:
- Auth模块已迁移至 `service/auth`
- Messaging模块已迁移至 `service/channels`
- 旧的import路径 `service/shared/auth` 和 `service/shared/messaging` 不再可用
- 可使用 `service/shared` 的兼容层作为过渡（`messaging_compat.go`）

---

## ⚠️ 注意事项

1. **模块独立性**: 各模块通过接口交互，禁止跨模块直接调用内部实现
2. **接口稳定性**: 接口一旦定义应保持稳定，变更需谨慎
3. **错误处理**: 统一使用项目错误处理规范
4. **日志记录**: 关键操作必须记录日志
5. **测试覆盖**: 核心功能必须有单元测试

---

## 📞 Commit 规范

```
feat(shared/auth): 实现JWT服务
fix(shared/cache): 修复缓存失效问题
test(shared/messaging): 添加消息队列测试
docs(shared): 更新使用说明
refactor(shared): 清理重复模块
```

---

*最后更新: 2026-01-12 - 清理重复模块*
