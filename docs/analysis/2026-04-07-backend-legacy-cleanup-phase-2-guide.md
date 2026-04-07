# 后端遗留债务清理 Phase 2 指南

日期：2026-04-07
分支：`chore/backend-legacy-cleanup-phase-2`
基线：`dev@ac7a2530`

## 目标

phase 2 应在不重新打开 phase 1 已完成的大范围低风险切片的前提下，继续推进后端遗留债务清理。

基于 2026-04-07 的新一轮引用扫描，当前更准确的剩余目标是：

- 清退最后一条已死亡的仓储层 token 黑名单路径
- 彻底移除空壳 `repository/interfaces/shared`
- 把验证重点从死亡路径测试转到 live auth 黑名单路径

本阶段仍然应坚持一次只处理一个收敛主题。

## 推荐优先级

### 首要候选

- 彻底移除 `TokenBlacklistRepository` 及其 Redis 仓储实现路径

原因：

- phase 1 已将 `repository/interfaces/shared` 收缩到只剩一个生产接口
- 当前代码搜索显示 `repository/redis.NewTokenBlacklistRepository` 已无活跃运行时调用方
- auth 运行时实际已经使用 `service/auth.RedisAdapter` 与 `service/auth.InMemoryTokenBlacklist`
- 继续保留死亡仓储接口只会制造噪音，不会保护当前行为

### 次级候选

只有在首要候选再次被明确延期时，才考虑其中之一：

- 清退 `repository/interfaces/shared/mocks` 下残余 legacy mock
- 在完成活跃模块 DTO 与兼容别名分离后，继续 writer 响应侧 DTO 清理

## 非目标

除非被拆成新的独立切片，否则不要把以下内容打包进 phase 2：

- `service/shared/stats` 重构
- AI 旧版 service wrapper 移除
- 与 `TokenBlacklistRepository` 无关的广义 auth/user 架构改造
- 当前清理主题之外的顺手重构

## 工作规则

1. 编辑待迁移或待删除符号前，先做 GitNexus 风格影响面确认。
2. 遇到 `HIGH` 或 `CRITICAL` 影响面时，应缩小切片，而不是强行推进。
3. 对死亡兼容路径，优先删除，不优先搬迁。
4. 除非本阶段明确包含行为修正，否则尽量保持运行时行为不变。
5. 文档只同步本阶段真正触达的模块路径与契约。

## Phase 2 检查点

### 检查点 1：确认影响面

1. 确认以下符号不存在活跃的非测试调用方：
   - `repository/interfaces/shared.TokenBlacklistRepository`
   - `repository/redis.NewTokenBlacklistRepository`
   - `repository/redis.NewTokenBlacklistRepositoryWithConfig`
2. 确认 auth 运行时仍然使用：
   - `service/auth.NewRedisAdapter`
   - `service/auth.NewInMemoryTokenBlacklist`
3. 一旦发现生产调用方，立即停止删除方案，改回包路径迁移方案

### 检查点 2：删除死亡仓储黑名单路径

1. 删除 `repository/interfaces/shared/token_blacklist_repository.go`
2. 删除 `repository/redis/token_blacklist_repository_redis.go`
3. 删除只保护死亡路径的自循环仓储测试

### 检查点 3：清掉残余 shared mock

1. 审核 `repository/interfaces/shared/mocks/*`
2. 若无代码引用，则删除相关文件
3. 确认 `repository/interfaces/shared` 可以彻底消失

### 检查点 4：加固 live auth token 路径

1. 为以下能力补充或扩展测试：
   - `service/auth.JWTServiceImpl` 黑名单行为
   - 如有缺口，则补 `service/auth.InMemoryTokenBlacklist` 兜底行为
2. 将 `service/user.UserServiceImpl` 的 token 生命周期调用接到 live auth 抽象
3. 确认已吊销 token 与 token 校验流程都通过当前活跃 auth 抽象生效

### 检查点 5：同步文档与验收说明

1. 将清理报告 / assessment 中的措辞从“剩余迁移”更新为“死亡路径清退”
2. 更新 repository 架构文档，移除 `interfaces/shared` 尾部包描述
3. 后续长期规划继续以父仓库计划文档为准

## 建议执行顺序

### 当 phase 2 选择清退死亡黑名单路径时

1. 先用 GitNexus 风格影响面分析加 `rg` 确认当前调用方与 import
2. 确认运行时不依赖后，再删除死亡接口与实现
3. 用 live auth 黑名单测试替代死亡路径测试
4. 若无引用，再删除 `repository/interfaces/shared/mocks` 残余文件
5. 更新架构与清理文档，确保不再把 `shared` 描述为活跃仓储接口包

## 验证基线

每完成一个有意义的子切片后，至少重跑以下校验：

```bash
go test ./service/auth -run 'TestJWTService|TestInMemoryTokenBlacklist' -count=1
go test ./service/auth ./service/user ./service/container ./router/user ./repository/interfaces/... -run '^$' -count=1
go test ./api/v1/writer ./service/writer/document ./service/writer/impl ./router/writer -run '^$' -count=1
```

若该切片暴露出新的隐藏调用方，再补更窄的定向测试。

## 合并标准

phase 2 只有在以下条件都满足时才可合并：

- 死亡黑名单仓储路径已从代码与文档中消失
- `repository/interfaces/shared` 已完全退场
- 面向用户的 token 生命周期入口不再返回占位结果
- 受影响包编译通过，且定向 auth / runtime 测试通过
- 延期项仍然保持延期，没有混入同一分支
