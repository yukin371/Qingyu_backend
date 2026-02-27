# Shared 模块 P1 架构重构规划

**文档日期**: 2026-02-12
**版本**: v1.1
**状态**: 执行中
**优先级**: P1
**关联任务**: ARCH-003 模块边界收口

---

## 执行摘要

基于 P0 架构修复的经验，需要对 Qingyu_backend 的 `service/shared/` 模块进行 P1 级别的架构重构。当前 shared 模块存在大量 TODO 注释、未实现功能，以及与 Port/Adapter 架构模式的不完整适配。

本规划旨在通过渐进式重构，明确 shared 模块各子模块的职责边界，完善接口定义，实现缺失功能，最终达成清晰的架构分层和完整的测试覆盖。

---

## 问题分析

### 当前 shared 模块结构（规划基线）

```
service/shared/
├── cache/
│   └── redis_cache_service.go       # Redis 缓存服务实现
├── storage/
│   ├── backend_factory.go           # 存储后端工厂
│   ├── image_processor.go          # 图片处理服务（部分实现）
│   ├── interfaces.go              # 存储服务接口定义
│   ├── local_backend.go           # 本地存储后端
│   ├── minio_backend.go          # MinIO 存储后端
│   ├── multipart_upload_service.go  # 分片上传服务（已实现，持续重构中）
│   ├── repository_adapter.go       # Repository 适配器
│   └── storage_service.go         # 存储服务主实现
├── stats/
│   └── stats_service.go           # 统计服务（大量 TODO）
├── metrics/
│   └── service_metrics.go         # 服务指标收集
├── config_service.go              # 配置服务
├── permission_service.go         # 权限服务
└── messaging_compat.go           # 消息兼容层
```

### 代码现状分析

#### 1. Stats 服务问题

| 文件 | 行数 | TODO 数量 | 主要问题 |
|------|------|----------|---------|
| `stats_service.go` | 477 | 15+ | 大量聚合查询未实现 |
| - | - | - | 阅读行为统计缺失 |
| - | - | - | 活跃度计算缺失 |
| - | - | - | 钱包收益统计缺失 |

**具体 TODO 列表**:
- `TODO(Task3)`: 扩展 BookRepository 支持 string author_id 查询 (第108, 231行)
- `TODO(Task3-聚合查询)`: MongoDB 聚合管道统计 (第161, 273行)
- `TODO(Task3-活跃度统计)`: 需要实现活跃度记录表 (第329行)
- `TODO(Task3-收益统计)`: 需要实现钱包交易记录 (第380行)
- `TODO(Task3)`: 需要阅读行为统计 (第136行)
- `TODO(Task3)`: 需要点赞统计 (第137行)
- `TODO(Task3)`: 需要评论统计 (第138行)
- `TODO(Task3)`: 需要钱包统计 (第139行)
- `TODO(Task3)`: 需要活跃度统计 (第143行)

#### 2. Storage 服务问题（规划基线）

| 功能模块 | 状态 | 描述 |
|---------|------|------|
| 基础上传/下载 | ✅ 已实现 | 文件上传、下载、删除基本功能 |
| 图片处理 | 🟡 部分实现 | 压缩、裁剪已实现，水印未完成 |
| 分片上传 | 🟡 已实现/待收敛 | 服务逻辑和测试已补齐，仍需 API 端点联调 |
| 断点续传 | ❌ 未实现 | 缺少分片管理和续传逻辑 |
| 云存储集成 | 🟡 部分实现 | MinIO 已集成，其他云存储未实现 |

**具体 TODO 列表** (storage_service.go 第15-19行):
```go
// TODO: 完善文件上传功能（分片上传、断点续传）
// TODO: 完善文件下载功能（断点续传、流式下载）
// TODO: 添加图片处理功能（压缩、裁剪、水印）
// TODO: 集成云存储服务（阿里云OSS、腾讯云COS、AWS S3）
// TODO: 实现文件版本管理
```

#### 3. Cache 服务状态

| 功能 | 状态 | 描述 |
|------|------|------|
| 基础 CRUD | ✅ 已实现 | Get/Set/Delete/Exists |
| 批量操作 | ✅ 已实现 | MGet/MSet/MDelete |
| 高级操作 | ✅ 已实现 | Expire/TTL/Increment/Decrement |
| 哈希操作 | ✅ 已实现 | HGet/HSet/HGetAll/HDelete |
| 集合操作 | ✅ 已实现 | SAdd/SMembers/SIsMember/SRemove |
| 有序集合 | ✅ 已实现 | ZAdd/ZRange/ZRemove |
| 缓存策略 | ❌ 未实现 | 缺少缓存失效策略和预热机制 |

#### 4. Port/Adapter 适配现状

| Port | 定义状态 | Adapter 状态 | 依赖方 |
|-------|---------|-------------|--------|
| StoragePort | ✅ 已定义 (ports.go) | ✅ 已实现 (adapters.go) | 部分服务 |
| CachePort | ✅ 已定义 (ports.go) | ✅ 已实现 (adapters.go) | 部分服务 |
| AuthPort | ✅ 已定义 (ports.go) | ✅ 已实现 (adapters.go) | 部分服务 |
| StatsPort | ❌ 未定义 | ❌ 未实现 | 无 |

---

## 重构方案

### 架构原则

遵循项目已建立的 Port/Adapter 模式，参考 P0 架构修复的成功经验：

1. **渐进式重构**: 分阶段实施，每阶段可独立验证和回滚
2. **接口先行**: 先定义清晰的 Port 接口，再实现具体 Adapter
3. **兼容性保证**: 保留旧接口，通过适配层渐进迁移
4. **测试驱动**: 每个功能实现前先定义测试用例

### 分层设计

```
┌─────────────────────────────────────────────────────────┐
│                    API Layer                        │
│           (api/v1/shared/*.go)                    │
└─────────────────────┬───────────────────────────────┘
                      │ 依赖 Port 接口
┌─────────────────────▼───────────────────────────────┐
│              Service Interfaces (Ports)               │
│        service/interfaces/shared/ports.go             │
│  - StoragePort    - CachePort    - StatsPort       │
│  - UploadPort     - ImagePort    - MetricsPort     │
└─────────────────────┬───────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────┐
│            Service Implementations                    │
│         service/shared/*/_service.go               │
│  - StorageService  - CacheService  - StatsService  │
└─────────────────────┬───────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────┐
│               Adapters                             │
│      service/shared/*/_adapter.go                   │
│  - RedisAdapter    - MongoAdapter   - S3Adapter   │
└─────────────────────┴───────────────────────────────┘
```

### 目录结构调整

```
service/shared/
├── storage/
│   ├── ports.go                  # 新增：存储端口定义
│   ├── storage_service.go        # 现有：保持不变
│   ├── adapters/
│   │   ├── local_adapter.go     # 从 local_backend.go 重构
│   │   ├── minio_adapter.go     # 从 minio_backend.go 重构
│   │   ├── s3_adapter.go       # 新增：AWS S3 适配器
│   │   ├── oss_adapter.go       # 新增：阿里云 OSS 适配器
│   │   └── cos_adapter.go       # 新增：腾讯云 COS 适配器
│   ├── image_processor.go       # 现有：增强水印功能
│   ├── upload/
│   │   ├── chunked_upload.go    # 新增：分片上传
│   │   ├── resumable_upload.go  # 新增：断点续传
│   │   └── upload_manager.go    # 新增：上传管理器
│   └── README.md               # 新增：模块说明
├── stats/
│   ├── ports.go                # 新增：统计端口定义
│   ├── stats_service.go       # 现有：完善聚合查询
│   ├── aggregators/
│   │   ├── user_aggregator.go   # 新增：用户统计聚合器
│   │   ├── content_aggregator.go # 新增：内容统计聚合器
│   │   └── revenue_aggregator.go # 新增：收益统计聚合器
│   ├── activity/
│   │   ├── tracker.go          # 新增：活跃度追踪器
│   │   └── calculator.go      # 新增：活跃度计算器
│   └── README.md              # 新增：模块说明
├── cache/
│   ├── cache_service.go       # 现有：保持不变
│   ├── strategies/
│   │   ├── cache_strategy.go    # 新增：缓存策略接口
│   │   ├── lru_strategy.go    # 新增：LRU 淘汰策略
│   │   └── ttl_strategy.go    # 新增：TTL 过期策略
│   └── README.md              # 新增：模块说明
└── metrics/
    ├── metrics_service.go      # 现有：保持不变
    └── README.md              # 新增：模块说明
```

---

## 任务分解

### 最新进展（2026-02-13）

本节记录当前 worktree 的实际落地进度，用于代码评审和后续合并跟踪。

#### 执行策略调整（架构优先）

- 自 2026-02-13 起，功能扩展项（如断点续传细节实现、图片水印、多云存储适配器）暂缓。
- 优先推进“降低维护成本”的架构类改造：依赖倒置、端口收口、接口去重、测试稳定性与可替换性提升。
- 所有新提交默认要求“不改变业务行为”的结构优化优先，功能新增按 `deferred` 状态记录。

#### 代码落地 Checklist

- [x] `service/shared/storage/multipart_upload_service.go`: 收口状态常量、分片大小常量、默认分类常量。
- [x] `service/shared/storage/multipart_upload_service.go`: 增加统一参数校验和 `context` 取消检查。
- [x] `service/shared/storage/multipart_upload_service.go`: 抽取 helper（分片校验、路径构建、状态判断、chunk reader 构建）。
- [x] `service/shared/storage/multipart_upload_service.go`: 修复 MD5 校验读取后 reader 状态问题（改为字节缓冲后校验与保存）。
- [x] `service/shared/storage/multipart_upload_service.go`: `extractCategory` 改为按路径真实解析，不再固定返回 `general`。
- [x] `service/shared/storage/mock/mocks.go`: 新增 `SetMultipartUploadExpiresAt` 以稳定构造上传过期场景。
- [x] `service/shared/storage/multipart_upload_service_test.go`: 过期用例改为真实断言错误。
- [x] `service/shared/storage/multipart_upload_service_test.go`: context 取消用例移除 `Skip`，改为断言 `context.Canceled`。
- [x] `service/shared/storage/storage_service.go`: 统一分页默认值和下载链接 TTL 默认值。
- [x] `service/shared/storage/storage_service.go`: 增加输入参数校验、`context` 取消检查和分类归一化。
- [x] `service/shared/storage/storage_service.go`: 上传流程改为单次读取并复用字节，避免 MD5 与存储读取冲突。
- [x] `service/shared/storage/storage_service.go`: 删除流程改为返回明确错误，不依赖 `fmt.Printf` 输出。
- [x] `service/shared/storage/storage_service_test.go`: 新增行为测试（默认分类、MD5、上下文取消、删除失败语义、默认下载 TTL）。
- [x] `api/v1/shared/storage_api.go`: 修复 multipart init 绑定逻辑（服务端注入 `uploaded_by`），并补充 `upload_id` 空值参数校验。
- [x] `api/v1/shared/storage_api_test.go`: 新增 multipart API 回归测试，覆盖 init/upload/progress/complete/abort 主链路与关键参数错误。
- [x] `service/shared/stats/ports.go`: 新增 `StatsPort/AggregatorPort` 首批接口定义。
- [x] `service/shared/stats/aggregators/*.go`: 新增用户/内容聚合器骨架实现。
- [x] `service/shared/stats/aggregators/aggregators_test.go`: 新增细粒度聚合器端口契约测试（`UserAggregatorPort`/`ContentAggregatorPort`）。
- [x] `service/shared/cache/strategies/*.go`: 新增缓存策略接口、默认 TTL 策略与策略管理器实现。
- [x] `service/shared/cache/strategies/strategies_test.go`: 新增策略匹配与默认 TTL 测试。
- [x] `service/container/service_container.go`: 移除 `Set/GetStorageServiceImpl` 具体实现入口，仅保留接口端口访问，减少容器层实现耦合。
- [x] `router/shared/shared_router.go` + `router/enter.go`: 路由注册全链路统一使用 `storage` 端口接口，不再依赖具体实现类型。
- [x] `router/shared/shared_router.go`: 删除未被运行时代码使用的 `RegisterRoutes` 兼容壳函数与专属依赖导入，降低路由维护面。
- [x] `service/shared/stats` + `api/v1/user/handler` + `router/user` + `router/enter`: 统计链路统一改用 `StatsPort`，移除 `PlatformStatsService` 接口别名，减少重复抽象命名。
- [x] `router/user/user_router.go` + `api/v1/user/handler/public_user_handler.go` + `router/enter.go`: 用户路由链路移除 `interface{}` 参数与运行时类型断言，改为明确端口类型（`UserRepository`/`BookstoreService`）。
- [x] `router/user/user_router.go` + `router/enter.go`: 用户路由显式接收可选 `StorageService` 端口并完成注入，移除“TODO 注释式依赖”。
- [x] `router/user/user_router.go`: 路由注册签名改为暴露本地 `BookstoreService` 端口类型，避免直接泄漏 `api handler` 层类型。
- [x] `api/v1/user/handler/profile_handler.go`: 统一 storage 端口命名（`sharedStorage`），消除导入别名/字段/参数同名歧义，降低阅读维护成本。
- [x] `router/enter.go`: 统一存储相关局部变量命名（`sharedStorageSvc/sharedStorageErr` 等），减少路由入口层命名歧义。
- [x] `service/container/service_container.go`: `SetupDefaultServices` 中存储初始化局部变量统一语义命名（Repository/Backend/Svc），降低容器初始化段阅读负担。
- [x] `service/container/service_container.go`: 存储初始化段注释与日志文案去掉过时“快速通道”措辞并统一术语，降低认知噪音。
- [x] `router/enter.go`: 清理过时 TODO 与大段注释伪代码（章节购买/用户关系/搜索服务说明），改为单行现状说明，降低入口文件维护噪音。
- [x] `router/enter.go`: 清理 `nolint` 占位注释（`ineffassign`/`nilness`）并改为显式未接入说明，降低静态检查噪音耦合。
- [x] `router/enter.go`: 用户路由注册段变量命名统一（`bookstoreSvcForUM`/`userRepoForUM`），去除历史 `...Interface`/`...Instance` 命名噪音。
- [x] `router/user/user_router.go`: `service/user` 双别名收敛为单一 `userService` 别名，消除同包多别名认知负担。
- [x] `router/enter.go`: 路由跳过日志文案统一为“跳过...路由注册”口径，减少运行日志语义不一致。
- [x] `service/container/service_container.go`: 清理容器内历史 `TODO:` 注释为“当前状态/占位说明”，避免误导为短期实现任务。

#### 验证 Checklist

- [x] `go test ./service/shared/storage -count=1`
- [x] `go test ./service/shared/storage ./service/interfaces/shared -count=1`
- [x] `go test ./api/v1/shared ./service/shared/storage ./service/shared/cache/strategies ./service/shared/stats/... -count=1`
- [x] `go test ./service/container -run Test__CompileOnly__ -count=1`
- [x] `go test ./api/v1/shared ./router/shared ./service/shared/storage ./service/shared/stats/... ./service/shared/cache ./service/shared/cache/strategies -count=1`
- [x] `go test ./router/shared ./router -run Test__CompileOnly__ -count=1`
- [x] `go test ./api/v1/user/handler ./router/user ./router ./service/shared/stats ./service/shared/stats/aggregators -count=1`
- [x] `go test ./api/v1/user/handler ./router/user ./router -count=1`
- [x] `go test ./api/v1/shared ./service/container ./service/shared/storage -run Test__CompileOnly__ -count=1`

#### 后续 Checklist（未完成）

- [x] `Task 2.1`: API 端点联调（分片上传相关路由与 handler）。
- [x] `Task 2.1`: 分片上传集成回归测试（API 回归测试已覆盖 init/upload/progress/complete/abort）。
- [x] `Task 1.2`: `stats` 模块 P1 主体任务启动并落地首批接口/聚合器（`ports.go` + `aggregators` 骨架）。
- [x] `Task 1.3`: `cache` 模块 P1 主体任务启动并落地策略接口（`strategies` 包 + 基础测试）。
- [x] 架构解耦：`StorageAPI` 从依赖具体实现改为依赖端口接口（API 层依赖倒置）。
- [x] 架构收口：`PlatformStatsService` 与 `StatsPort` 去重，统一单一端口接口来源。
- [x] 架构收口：`RedisCacheService` 接入 `StrategyManager`（仅结构接入，不新增功能行为）。
- [x] 架构收口：`storage` 模块统一导出 API 依赖端口（`StorageService`/`MultipartUploadManager`/`ImageProcessorService`），移除 API 层重复接口定义。
- [x] 架构收口：`stats` 聚合端口拆分为细粒度接口（`UserAggregatorPort`/`ContentAggregatorPort`），保留兼容组合接口。
- [x] 架构收口：`ServiceContainer`/`router` 存储链路剔除具体实现依赖，仅通过端口接口完成注入与路由注册。
- [x] 架构收口：补齐 `stats` 聚合器端口契约测试，确保后续实现替换不影响接口约束。

### 单页总 Checklist（汇总）

#### Task 1.1 Storage Port/Adapter

- [x] Port 接口定义完成并通过审查（`service/shared/storage/interfaces.go` + API/router/container 依赖倒置完成）。
- [x] LocalAdapter 实现通过测试（`service/shared/storage` 包测试通过）。
- [x] MinIO 适配到新接口（`MinIOBackend` 增加 `StorageBackend` 契约断言 + 工厂配置校验测试）。
- [x] 单元测试覆盖率 > 80%（Task 1.1 范围）。

#### Task 1.2 Stats Port/Adapter

- [x] Port 接口定义完成（`service/shared/stats/ports.go`）。
- [x] Aggregator 基础结构实现（`service/shared/stats/aggregators/*.go`）。
- [x] MongoDB 聚合管道设计文档完成（`docs/plans/2026-02-13-stats-mongodb-aggregation-design.md`）。
- [x] 单元测试框架搭建完成（已落地聚合器端口契约测试 `service/shared/stats/aggregators/aggregators_test.go`）。

#### Task 1.3 Cache 策略管理

- [x] 策略接口定义完成（`service/shared/cache/strategies/cache_strategy.go`）。
- [x] TTL 策略实现（`service/shared/cache/strategies/ttl_strategy.go`）。
- [x] RedisCacheService 集成策略管理（`service/shared/cache/redis_cache_service.go`）。
- [x] 单元测试通过（`service/shared/cache/strategies/strategies_test.go`）。

#### Task 2.1 分片上传与断点续传

- [x] 分片上传功能完成（Service 层）。
- [ ] 断点续传功能完成（deferred：功能实现暂缓）。
- [x] 单元测试覆盖率 > 80%（`service/shared/storage` 包内）。
- [x] 集成测试通过（API 回归链路）。
- [x] API 端点实现。

#### Task 2.2 图片处理

- [ ] 图片水印功能实现（deferred：功能实现暂缓）。
- [ ] 文字水印功能实现（deferred：功能实现暂缓）。
- [ ] 性能基准测试通过（deferred）。
- [ ] 单元测试覆盖率 > 80%（deferred）。

#### Task 2.3 云存储适配器

- [ ] S3 适配器实现（deferred：功能实现暂缓）。
- [ ] OSS 适配器实现（deferred：功能实现暂缓）。
- [ ] COS 适配器实现（deferred：功能实现暂缓）。
- [ ] 工厂模式实现（deferred：功能实现暂缓）。
- [ ] 各云存储单元测试通过（deferred）。

#### Task 2.4 Stats 聚合查询

- [ ] 所有 TODO 转换为实现。
- [ ] 聚合查询性能测试通过。
- [ ] 单元测试覆盖率 > 80%。
- [ ] API 返回真实数据。

#### Task 2.5 活跃度追踪

- [ ] 活跃度追踪实现。
- [ ] 活跃度计算实现。
- [ ] 性能测试通过。
- [ ] 单元测试覆盖率 > 80%。

#### Task 2.6 Cache 优化

- [ ] LRU 策略实现。
- [ ] 缓存预热实现。
- [ ] 失效策略实现。
- [ ] 性能对比测试通过。

#### Task 3.1 测试覆盖

- [ ] 单元测试覆盖率 > 90%。
- [ ] 集成测试覆盖率 > 80%。
- [ ] 性能测试基准建立。
- [ ] 混沌测试（网络中断、存储故障）。

#### Task 3.2 性能优化

- [ ] 性能基准对比报告。
- [ ] P95 延迟降低 > 20%。
- [ ] 吞吐量提升 > 30%。

#### Task 3.3 监控与可观测性

- [ ] 所有指标导出到 Prometheus。
- [ ] Grafana 仪表板配置。
- [ ] 告警规则配置。

#### 全局验收门槛

- [ ] 所有 TODO 清理或转换为实现。
- [ ] Port 接口定义清晰。
- [ ] ServiceContainer 正确管理依赖。
- [ ] 单元测试通过。
- [ ] 编译零警告。
- [ ] 单元测试覆盖率 > 80%（一般标准）。
- [ ] 集成测试覆盖率 > 70%（一般标准）。
- [ ] 性能基准建立（一般标准）。
- [ ] API 文档更新（一般标准）。
- [ ] 单元测试覆盖率 > 90%（理想标准）。
- [ ] 集成测试覆盖率 > 80%（理想标准）。
- [ ] 性能优化报告（理想标准）。
- [ ] 监控指标完整（理想标准）。
- [ ] 架构文档更新（理想标准）。

### P0 任务（核心基础，预计 3-4 天）

#### Task 1.1: 设计 Storage 层 Port/Adapter

**目标**: 定义清晰的存储接口，适配现有实现到新架构

**文件**:
- 创建: `service/shared/storage/ports.go`
- 修改: `service/interfaces/shared/ports.go` (扩展 StoragePort)
- 创建: `service/shared/storage/adapters/local_adapter.go`
- 修改: `service/shared/storage/local_backend.go` (重构为 adapter)

**接口定义**:
```go
package storage

type StoragePort interface {
    // 基础操作
    Upload(ctx context.Context, req *UploadRequest) (*FileInfo, error)
    Download(ctx context.Context, fileID string) (io.ReadCloser, error)
    Delete(ctx context.Context, fileID string) error
    Exists(ctx context.Context, fileID string) (bool, error)

    // URL 生成
    GetURL(ctx context.Context, fileID string, expiresIn time.Duration) (string, error)
}

type ChunkedUploadPort interface {
    // 分片上传
    InitUpload(ctx context.Context, req *InitUploadRequest) (*UploadSession, error)
    UploadChunk(ctx context.Context, sessionID string, chunk *Chunk) error
    CompleteUpload(ctx context.Context, sessionID string) (*FileInfo, error)
    CancelUpload(ctx context.Context, sessionID string) error
}
```

**验收标准**:
- [ ] Port 接口定义完成并通过审查
- [ ] LocalAdapter 实现通过测试
- [ ] MinIO 适配到新接口
- [ ] 单元测试覆盖率 > 80%

---

#### Task 1.2: 设计 Stats 层 Port/Adapter

**目标**: 为统计服务定义接口，支持聚合查询

**文件**:
- 创建: `service/shared/stats/ports.go`
- 修改: `service/interfaces/shared/ports.go` (添加 StatsPort)
- 创建: `service/shared/stats/aggregators/user_aggregator.go`
- 创建: `service/shared/stats/aggregators/content_aggregator.go`

**接口定义**:
```go
package stats

type StatsPort interface {
    // 用户统计
    GetUserStats(ctx context.Context, userID string) (*UserStats, error)
    GetPlatformUserStats(ctx context.Context, filter *StatsFilter) (*PlatformUserStats, error)

    // 内容统计
    GetContentStats(ctx context.Context, userID string) (*ContentStats, error)
    GetPlatformContentStats(ctx context.Context, filter *StatsFilter) (*PlatformContentStats, error)

    // 活跃度统计
    GetUserActivity(ctx context.Context, userID string, days int) (*ActivityStats, error)
}

type AggregatorPort interface {
    AggregateUserStats(ctx context.Context, filter *StatsFilter) (*UserStats, error)
    AggregateContentStats(ctx context.Context, filter *StatsFilter) (*ContentStats, error)
}
```

**验收标准**:
- [ ] Port 接口定义完成
- [ ] Aggregator 基础结构实现
- [ ] MongoDB 聚合管道设计文档
- [ ] 单元测试框架搭建完成

---

#### Task 1.3: 设计 Cache 层策略管理

**目标**: 为缓存服务添加策略管理能力

**文件**:
- 创建: `service/shared/cache/strategies/cache_strategy.go`
- 创建: `service/shared/cache/strategies/ttl_strategy.go`
- 修改: `service/shared/cache/redis_cache_service.go`

**接口定义**:
```go
package strategies

type CacheStrategy interface {
    // ShouldCache 判断是否应该缓存
    ShouldCache(key string, value interface{}) bool

    // GetTTL 获取缓存过期时间
    GetTTL(key string) time.Duration

    // OnMiss 缓存未命中时的处理
    OnMiss(key string) error
}

type CacheManager interface {
    RegisterStrategy(pattern string, strategy CacheStrategy)
    GetStrategy(key string) CacheStrategy
    SetWithStrategy(ctx context.Context, key string, value interface{}) error
}
```

**验收标准**:
- [ ] 策略接口定义完成
- [ ] TTL 策略实现
- [ ] RedisCacheService 集成策略管理
- [ ] 单元测试通过

---

### P1 任务（功能完善，预计 5-7 天）

#### Task 2.1: 实现分片上传

**目标**: 完成文件分片上传和断点续传功能

**文件**:
- 创建: `service/shared/storage/upload/chunked_upload.go`
- 创建: `service/shared/storage/upload/resumable_upload.go`
- 创建: `service/shared/storage/upload/upload_manager.go`
- 修改: `service/shared/storage/multipart_upload_service.go`

**实现内容**:
1. **ChunkedUploadService**:
   - `InitUpload()` - 初始化上传会话
   - `UploadChunk()` - 上传单个分片
   - `CompleteUpload()` - 合并分片
   - `CancelUpload()` - 取消上传

2. **ResumableUpload**:
   - 分片状态管理（MongoDB 持久化）
   - 断点续传逻辑
   - 分片校验（MD5）

3. **UploadManager**:
   - 并发上传控制
   - 进度跟踪
   - 错误重试

**数据模型**:
```go
type UploadSession struct {
    ID           string    `json:"id" bson:"_id"`
    UserID       string    `json:"user_id" bson:"user_id"`
    Filename     string    `json:"filename" bson:"filename"`
    TotalSize    int64     `json:"total_size" bson:"total_size"`
    ChunkSize    int64     `json:"chunk_size" bson:"chunk_size"`
    ChunkCount   int       `json:"chunk_count" bson:"chunk_count"`
    UploadedChunks []int    `json:"uploaded_chunks" bson:"uploaded_chunks"`
    Status       string    `json:"status" bson:"status"` // uploading, completed, cancelled
    CreatedAt    time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
```

**验收标准**:
- [x] 分片上传功能完成（Service 层）
- [ ] 断点续传功能完成
- [x] 单元测试覆盖率 > 80%（`service/shared/storage` 包内）
- [ ] 集成测试通过
- [ ] API 端点实现

---

#### Task 2.2: 完善图片处理

**目标**: 完成图片水印功能，优化处理流程

**文件**:
- 修改: `service/shared/storage/image_processor.go`
- 创建: `service/shared/storage/watermark.go`

**实现内容**:
1. **水印功能**:
   - 图片水印（支持 PNG 透明图）
   - 文字水印（自定义字体、颜色、位置）
   - 水印位置计算
   - 透明度控制

2. **性能优化**:
   - 并行处理多张图片
   - 图片缓存（避免重复处理）
   - 流式处理（减少内存占用）

**接口扩展**:
```go
type WatermarkOptions struct {
    Type        string  // "image" or "text"
    ImagePath   string  // 图片水印路径
    Text        string  // 文字水印内容
    Font       string  // 字体路径
    FontSize   int      // 字体大小
    Color      string  // 文字颜色
    Position   string  // 位置: top-left, center, bottom-right 等
    Opacity   float64 // 透明度 0-1
    Padding    int      // 边距
}

func (p *ImageProcessor) ApplyWatermark(ctx context.Context, sourcePath string, opts *WatermarkOptions) (string, error)
```

**验收标准**:
- [ ] 图片水印功能实现
- [ ] 文字水印功能实现
- [ ] 性能基准测试通过
- [ ] 单元测试覆盖率 > 80%

---

#### Task 2.3: 实现云存储适配器

**目标**: 添加阿里云 OSS、腾讯云 COS、AWS S3 支持

**文件**:
- 创建: `service/shared/storage/adapters/s3_adapter.go`
- 创建: `service/shared/storage/adapters/oss_adapter.go`
- 创建: `service/shared/storage/adapters/cos_adapter.go`
- 修改: `service/shared/storage/backend_factory.go`

**实现内容**:
1. **S3Adapter**: AWS S3 兼容接口
2. **OSSAdapter**: 阿里云 OSS
3. **COSAdapter**: 腾讯云 COS
4. **BackendFactory**: 动态选择存储后端

**配置**:
```go
type StorageConfig struct {
    Backend   string                 // "local", "minio", "s3", "oss", "cos"
    Local     *LocalBackendConfig
    MinIO     *MinIOConfig
    S3        *S3Config
    OSS       *OSSConfig
    COS       *COSConfig
}

func NewStorageBackend(config *StorageConfig) (StorageBackend, error)
```

**验收标准**:
- [ ] S3 适配器实现
- [ ] OSS 适配器实现
- [ ] COS 适配器实现
- [ ] 工厂模式实现
- [ ] 各云存储单元测试通过

---

#### Task 2.4: 实现 Stats 聚合查询

**目标**: 完成 MongoDB 聚合查询实现，替换所有 TODO

**文件**:
- 创建: `service/shared/stats/aggregators/user_aggregator.go`
- 创建: `service/shared/stats/aggregators/content_aggregator.go`
- 创建: `service/shared/stats/aggregators/revenue_aggregator.go`
- 修改: `service/shared/stats/stats_service.go`

**实现内容**:
1. **用户统计聚合**:
   - 总用户数、新增用户
   - 活跃用户、VIP 用户
   - 留存率计算

2. **内容统计聚合**:
   - 总书籍数、新增书籍
   - 总章节数、总字数
   - 热门分类统计

3. **收益统计聚合**:
   - 总收益、期间收益
   - 按书籍分组收益
   - 按收益类型分组

**MongoDB 聚合管道示例**:
```go
pipeline := mongo.Pipeline{
    {{"$match", bson.D{
        {"created_at", bson.D{
            {"$gte", startDate},
            {"$lte", endDate},
        }},
    }}},
    {{"$group", bson.D{
        {"_id", "$category"},
        {"count", bson.D{{"$sum", 1}}},
        {"total_words", bson.D{{"$sum", "$word_count"}}},
    }}},
    {{"$sort", bson.D{{"count", -1}}},
}
```

**验收标准**:
- [ ] 所有 TODO 转换为实现
- [ ] 聚合查询性能测试通过
- [ ] 单元测试覆盖率 > 80%
- [ ] API 返回真实数据

---

#### Task 2.5: 实现活跃度追踪

**目标**: 实现用户活跃度记录和计算

**文件**:
- 创建: `service/shared/stats/activity/tracker.go`
- 创建: `service/shared/stats/activity/calculator.go`
- 创建: `models/user_activity_log.go`

**数据模型**:
```go
type UserActivityLog struct {
    ID          string    `json:"id" bson:"_id"`
    UserID      string    `json:"user_id" bson:"user_id"`
    Action      string    `json:"action" bson:"action"` // read, write, comment, like
    TargetID    string    `json:"target_id" bson:"target_id"`
    TargetType  string    `json:"target_type" bson:"target_type"` // book, chapter, comment
    Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
    Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
}
```

**实现内容**:
1. **ActivityTracker**:
   - 记录用户行为
   - 异步批量写入（使用 EventBus）
   - 数据过期策略

2. **ActivityCalculator**:
   - 计算活跃天数
   - 计算活跃时段
   - 操作类型分布

**验收标准**:
- [ ] 活跃度追踪实现
- [ ] 活跃度计算实现
- [ ] 性能测试通过
- [ ] 单元测试覆盖率 > 80%

---

#### Task 2.6: 优化 Cache 策略

**目标**: 实现缓存预热、失效策略

**文件**:
- 创建: `service/shared/cache/strategies/lru_strategy.go`
- 创建: `service/shared/cache/strategies/prefetch.go`
- 修改: `service/shared/cache/redis_cache_service.go`

**实现内容**:
1. **LRU 淘汰策略**:
   - 基于 Redis ZSet 实现
   - 可配置淘汰阈值
   - 淘汰回调

2. **缓存预热**:
   - 启动时加载热门数据
   - 基于访问模式预测
   - 异步预热

3. **失效策略**:
   - 主动失效（数据更新时）
   - 被动失效（TTL 过期）
   - 定时刷新（关键数据）

**验收标准**:
- [ ] LRU 策略实现
- [ ] 缓存预热实现
- [ ] 失效策略实现
- [ ] 性能对比测试通过

---

### P2 任务（优化增强，预计 3-4 天）

#### Task 3.1: 完善测试覆盖

**目标**: 提升测试覆盖率到理想标准

**文件**:
- 创建: `service/shared/storage/storage_service_test.go`
- 创建: `service/shared/stats/stats_service_test.go`
- 创建: `service/shared/cache/cache_service_test.go`
- 创建: `test/integration/shared_integration_test.go`

**验收标准**:
- [ ] 单元测试覆盖率 > 90%
- [ ] 集成测试覆盖率 > 80%
- [ ] 性能测试基准建立
- [ ] 混沌测试（网络中断、存储故障）

---

#### Task 3.2: 性能优化

**目标**: 优化关键路径性能

**优化点**:
1. **Stats 聚合**:
   - 添加索引
   - 结果缓存
   - 异步计算

2. **Storage 上传**:
   - 并行分片上传
   - 断点续传优化
   - CDN 集成

3. **Cache 缓存**:
   - Pipeline 批量操作
   - 本地多级缓存
   - 缓存压缩

**验收标准**:
- [ ] 性能基准对比报告
- [ ] P95 延迟降低 > 20%
- [ ] 吞吐量提升 > 30%

---

#### Task 3.3: 监控和可观测性

**目标**: 添加完善的监控指标

**文件**:
- 修改: `service/shared/metrics/service_metrics.go`
- 创建: `service/shared/storage/metrics.go`
- 创建: `service/shared/stats/metrics.go`

**指标**:
1. **Storage 指标**:
   - 上传/下载计数
   - 存储空间使用
   - 错误率

2. **Stats 指标**:
   - 聚合查询耗时
   - 缓存命中率
   - 数据新鲜度

3. **Cache 指标**:
   - 命中率
   - 淘汰率
   - 内存使用

**验收标准**:
- [ ] 所有指标导出到 Prometheus
- [ ] Grafana 仪表板配置
- [ ] 告警规则配置

---

## 验收标准

### 最低标准（必须满足）

- [ ] 所有 TODO 清理或转换为实现
- [ ] Port 接口定义清晰
- [ ] ServiceContainer 正确管理依赖
- [ ] 单元测试通过
- [ ] 编译零警告

### 一般标准（推荐满足）

- [ ] 单元测试覆盖率 > 80%
- [ ] 集成测试覆盖率 > 70%
- [ ] 性能基准建立
- [ ] API 文档更新

### 理想标准（尽量满足）

- [ ] 单元测试覆盖率 > 90%
- [ ] 集成测试覆盖率 > 80%
- [ ] 性能优化报告
- [ ] 监控指标完整
- [ ] 架构文档更新

---

## 检查点与里程碑

### 检查点清单

| 检查点 | 检查内容 | 预期时间 | 状态 |
|--------|----------|----------|------|
| CP1 | P0 任务完成，Port 接口定义完成 | Day 4 | ⬜ |
| CP2 | 分片上传功能实现完成 | Day 8 | 🟡 |
| CP3 | Stats 聚合查询实现完成 | Day 12 | ⬜ |
| CP4 | 云存储适配器实现完成 | Day 15 | ⬜ |
| CP5 | 所有 TODO 清理完成 | Day 18 | ⬜ |
| CP6 | 测试覆盖率达到一般标准 | Day 21 | 🟡 |

### 里程碑

| 里程碑 | 描述 | 预期时间 | 状态 |
|--------|------|----------|------|
| M1 | Storage 层 Port/Adapter 架构完成 | Day 4 | ⬜ |
| M2 | Stats 层 Port/Adapter 架构完成 | Day 8 | ⬜ |
| M3 | 分片上传功能可用 | Day 10 | 🟡 |
| M4 | 所有核心功能实现完成 | Day 18 | ⬜ |
| M5 | 重构完成，可合并主分支 | Day 21 | ⬜ |

---

## 风险与应对

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| 聚合查询性能问题 | 高 | 中 | 提前进行性能测试，必要时引入预计算 |
| 云存储 API 兼容性 | 中 | 中 | 统一接口抽象，添加适配层测试 |
| 分片上传状态管理复杂 | 高 | 高 | 使用现有 Event Store 模式持久化状态 |
| 测试数据准备 | 中 | 低 | 使用现有的 seed_data 工具扩展 |
| 性能回归 | 高 | 低 | 建立性能基准，持续监控 |

---

## 回滚计划

### 回滚触发条件

- 核心测试失败超过 30%
- 性能回归超过 50%
- 出现新的 P0/P1 问题

### 回滚策略

1. **单任务回滚**: 单个 Task 失败时，仅回滚该 Task
2. **阶段回滚**: 整个 P0/P1 阶段失败时，回滚到上一稳定版本
3. **完整回滚**: 所有任务失败时，回滚到重构前状态

### 回滚步骤

```bash
# 1. 创建备份分支
git checkout -b backup/pre-refactor

# 2. 回滚到重构前
git revert <refactor-commit-range>

# 3. 验证系统正常
go test ./...
go run cmd/server/main.go
```

---

## 关联文档

- [P0 架构修复计划](./2026-02-12-p0-architecture-fix-implementation.md)
- [架构改进计划](./2026-02-07-architecture-refinement-plan.md)
- [Shared 架构简化 PR 草案](./2026-02-13-shared-architecture-simplification-pr-draft.md)
- [Stats MongoDB 聚合管道设计](./2026-02-13-stats-mongodb-aggregation-design.md)
- [Shared 模块重构计划](../plan/shared-module-refactor-plan.md)
- [Port/Adapter 设计文档](../architecture/port-adapter-pattern.md)

---

## 工作目录

**Worktree**: `Qingyu_backend_shared-p1-refactor`
**分支**: `feature/shared-p1-refactor`
**基准分支**: `main`

---

## 文档维护

| 项目 | 内容 |
|------|------|
| 更新频率 | 每完成一个 Task 更新一次 |
| 责任人 | Kore |
| 版本历史 | v1.0 (2026-02-12): 初始版本；v1.1 (2026-02-13): 补充 storage 模块阶段性落地进度与验证记录 |

---

## 附录：技术栈与依赖

### 现有依赖

| 包名 | 版本 | 用途 |
|------|------|------|
| github.com/redis/go-redis/v9 | - | Redis 客户端 |
| github.com/disintegration/imaging | - | 图片处理 |
| github.com/minio/minio-go/v7 | - | MinIO 客户端 |

### 新增依赖

| 包名 | 用途 |
|------|------|
| github.com/aws/aws-sdk-go-v2 | AWS S3 支持 |
| github.com/aliyun/aliyun-oss-go-sdk | 阿里云 OSS 支持 |
| github.com/tencentyun/cos-go-sdk-v5 | 腾讯云 COS 支持 |
| github.com/go-playground/validator/v10 | 参数验证增强 |

---

**文档结束**
