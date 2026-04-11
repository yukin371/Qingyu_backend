# 后端遗留债务清理 Phase 2 指南

日期：2026-04-07
状态：已完成（2026-04-12 收尾验证通过）
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

### 首要候选 ✅ 已完成

- 彻底移除 `TokenBlacklistRepository` 及其 Redis 仓储实现路径

原因：

- phase 1 已将 `repository/interfaces/shared` 收缩到只剩一个生产接口
- 当前代码搜索显示 `repository/redis.NewTokenBlacklistRepository` 已无活跃运行时调用方
- auth 运行时实际已经使用 `service/auth.RedisAdapter` 与 `service/auth.InMemoryTokenBlacklist`
- 继续保留死亡仓储接口只会制造噪音，不会保护当前行为

### 次级候选 ✅ 已完成

- 清退 `repository/interfaces/shared/mocks` 下残余 legacy mock
- `repository/interfaces/shared` 目录已从受管代码中完全消失

## 非目标（仍延期）

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

### 检查点 1：确认影响面 ✅

1. 确认以下符号不存在活跃的非测试调用方：
   - `repository/interfaces/shared.TokenBlacklistRepository`
   - `repository/redis.NewTokenBlacklistRepository`
   - `repository/redis.NewTokenBlacklistRepositoryWithConfig`
2. 确认 auth 运行时仍然使用：
   - `service/auth.NewRedisAdapter`
   - `service/auth.NewInMemoryTokenBlacklist`

### 检查点 2：删除死亡仓储黑名单路径 ✅

1. `repository/interfaces/shared/token_blacklist_repository.go` 已删除
2. `repository/redis/token_blacklist_repository_redis.go` 已删除
3. 只保护死亡路径的自循环仓储测试已删除

### 检查点 3：清掉残余 shared mock ✅

1. `repository/interfaces/shared/mocks/*` 已审核并删除
2. `repository/interfaces/shared` 目录已从受管代码中完全消失
3. grep 确认无残留引用

### 检查点 4：加固 live auth token 路径 ✅

1. `service/auth.JWTServiceImpl` 黑名单行为已有测试覆盖：
   - `TestJWTService_RevokeTokenMarksTokenAsRevoked`
   - `TestJWTService_ValidateTokenRejectsRevokedToken`
   - `TestInMemoryTokenBlacklistBasicFlow`
2. `service/user.UserServiceImpl` 的 token 生命周期调用已接到 live auth 抽象：
   - `LogoutUser` 通过 `TokenLifecycleService.Logout`
   - `ValidateToken` 通过 `TokenLifecycleService.ValidateTokenUserID`
3. 2026-04-12 补齐 nil-lifecycle 降级路径测试
4. 已吊销 token 与 token 校验流程都通过当前活跃 auth 抽象生效

### 检查点 5：同步文档与验收说明 ✅

1. 清理报告 / assessment 中的措辞已从"剩余迁移"更新为"死亡路径清退"
2. repository 架构文档已移除 `interfaces/shared` 尾部包描述
3. 后续长期规划继续以父仓库计划文档为准

## 建议执行顺序

### 当 phase 2 选择清退死亡黑名单路径时 ✅ 已完成

1. GitNexus 风格影响面分析 + grep 确认当前调用方与 import ✅
2. 确认运行时不依赖后，再删除死亡接口与实现 ✅
3. 用 live auth 黑名单测试替代死亡路径测试 ✅
4. 删除 `repository/interfaces/shared/mocks` 残余文件 ✅
5. 更新架构与清理文档 ✅

## 验证基线（2026-04-12 执行结果）

```bash
# 定向测试：全部 PASS（20/20）
go test ./service/auth -run 'TestJWTService|TestInMemoryTokenBlacklist' -count=1     # 3/3 PASS
go test ./service/user -run 'TestUserService_(LogoutUser|ValidateToken)' -count=1    # 14/14 PASS
go test ./api/v1/user/handler -run 'TestAuthHandler_Logout' -count=1                 # 3/3 PASS

# 受影响包编译空跑：全部通过
go test ./service/auth ./service/user ./service/container ./router/user ./repository/interfaces/... -run '^$' -count=1
go test ./api/v1/user ./api/v1/user/handler ./api/v1/writer ./service/writer/document ./service/writer/impl ./router/writer -run '^$' -count=1
```

## 合并标准 ✅ 全部满足

- 死亡黑名单仓储路径已从代码与文档中消失 ✅
- `repository/interfaces/shared` 已完全退场 ✅
- 面向用户的 token 生命周期入口不再返回占位结果 ✅
- 受影响包编译通过，且定向 auth / runtime 测试通过 ✅
- 延期项仍然保持延期，没有混入同一分支 ✅
