---
name: Shared Service
description: 跨模块共享的基础设施服务层，提供配置管理、RBAC权限、Redis缓存、文件存储和指标收集能力
type: module
---

# Shared Service

> 最后更新：2026-05-21

## 职责

为所有业务模块（writer、bookstore、social、finance 等）提供跨领域基础设施服务：动态配置管理（热更新+备份恢复）、完整 RBAC 权限控制、Redis 缓存（含策略化 TTL）、多后端文件存储（本地/MinIO）和运行时指标收集。不包含业务逻辑本身。

## 数据流

```
各业务模块 Service
    ├→ ConfigService          → config.yaml (读写 + 热重载)
    ├→ PermissionService      → auth.PermissionRepository → MongoDB
    ├→ RedisCacheService      → Redis (基础/批量/Hash/Set/ZSet)
    │       └→ CacheStrategy  → 按 key 前缀匹配 TTL 策略
    ├→ StorageServiceImpl     → BackendFactory → LocalBackend / MinIOBackend
    │       └→ MultipartUploadService → 分片上传(5MB chunk, 24h TTL)
    │       └→ ImageProcessor → 缩略图/压缩
    └→ ServiceMetrics         → 内存聚合(调用次数/耗时/健康状态)
```

## 约定 & 陷阱

- **存储后端切换通过 BackendFactory**：业务层只依赖 `StorageService` 接口，不直接构造 `LocalBackend` 或 `MinIOBackend`；新增后端（OSS/COS/S3）只需实现 `StorageBackend` 接口并注册到工厂
- **RepositoryAdapter 桥接两层接口**：storage 子包的 `FileRepository` 接口通过 `RepositoryAdapter` 适配到 `repository/interfaces/storage.StorageRepository`，避免循环依赖
- **缓存策略按 key 前缀匹配**：`StrategyManager` 使用 `strings.HasPrefix` 匹配，注册顺序影响优先级；未匹配的 key 回退到默认 TTL（5min）
- **ConfigService 读写全局配置文件**：通过 `config.GlobalConfig` 单例读取，更新时写文件并 `reloadConfig()`；备份文件存于同目录，恢复依赖自动备份
- **PermissionService 是薄封装**：直接委托 `auth.PermissionRepository`，不含缓存层，高频权限检查场景需调用方自行缓存
- **已迁移模块不可用**：auth 模块已迁移至 `service/auth/`，messaging 已迁移至 `service/channels/`，旧的 `service/shared/auth` 和 `service/shared/messaging` 路径不存在
