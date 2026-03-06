# Issue 关联关系分析

**更新日期**: 2026-03-06
**活跃Issue总数**: 6
**已归档Issue总数**: 7

## 已归档

- [#013: 测试用户种子数据ID未设置问题](./archived/013-test-user-seed-id-not-set.md) - ✅ 已修复并归档（2026-03-06）
- [#001: 统一模型层 ID 字段类型](./archived/001-unify-id-type-in-models.md) - ✅ 模型层 `_id` 统一已完成并归档（2026-03-06）
- [#002: Repository Create 方法未回设 ID](./archived/002-create-method-id-not-set-bug.md) - ✅ 核心与低频批量入口已修复并归档（2026-03-06）
- [#006: 数据库索引问题](./archived/006-database-index-issues.md) - ✅ 核心问题已解决并归档（2026-03-06）
- [#007: Service 层事务管理缺失](./archived/007-transaction-management.md) - ✅ 高风险业务事务问题已解决并归档（2026-03-06）
- [#008: 中间件架构问题](./archived/008-middleware-architecture-issues.md) - ✅ 核心阻塞已解决并归档（2026-03-06）
- [#012: 401认证错误和权限配置问题](./archived/012-auth-401-and-permission-issues.md) - ✅ 核心认证/权限阻塞已解决并归档（2026-03-06）

---

## 关联矩阵

| Issue | 关联Issue | 关联类型 | 说明 |
|-------|----------|---------|------|
| #003 | #009 | 🔗 相关 | 测试基础设施和覆盖率 |
| #004 | - | ⚪ 独立 | 代码质量问题 |
| #005 | #011 | 🔗 相关 | API响应格式与数据类型 |
| #009 | #003 | 🔗 相关 | 测试覆盖率与基础设施 |
| #010 | #007(已归档) | 🔄 依赖 | Repository层问题继续参考已归档事务治理 |
| #011 | #005 | 🔗 相关 | 数据类型与API |

---

## 需要联合处理的Issue组合

### 组合1: 类型与接口问题组 🔄

**关联Issue**: #011, #005

**核心问题**: 前后端数据类型与 API 格式仍需继续统一

**处理策略**:
1. 已归档参考: #001 模型层 `_id` 已统一
2. 已归档参考: #002 Repository Create 回设已收口
3. #011: 继续处理前后端类型转换与枚举不一致
4. #005: 继续处理 API 标准化

---

### 组合2: 测试问题组 🔗

**关联Issue**: #003, #009

**核心问题**: 测试基础设施不足导致覆盖率低

**处理策略**: 并行处理
- #003: 测试基础设施改进
- #009: 测试覆盖率提升
- 已归档参考: #013 种子数据修复

**协同重点**:
- 种子数据问题（#013）已修复，当前重点转向测试覆盖率（#009）
- 基础设施改进（#003）支持覆盖率提升（#009）

---

### 组合3: API标准化组 🔗

**关联Issue**: #005, #011

**核心问题**: API响应格式和数据类型不一致

**处理策略**: 联合处理
- #005: 统一响应码、URL前缀、分页格式
- #011: 统一前后端数据类型（包括响应格式）

**协同重点**:
- 响应拦截器处理（#011）与响应码统一（#005）
- 分页格式（#005）与数据类型（#011）

---

### 组合4: 架构重构组 🔄

**关联Issue**: #010

**核心问题**: Repository层业务逻辑渗透仍需继续外移

**处理策略**:
1. 已归档参考: #007 已完成高风险事务收敛
2. #010: 继续将Repository层业务逻辑移到Service层

---

## 需要拆分的大型Issue

### Issue #010: Repository 层业务逻辑渗透

**当前规模**: 58个方法，23个文件

**建议拆分为**:

#### #010-A: Bookstore域Repository重构（P0）
- ranking_repository_mongo.go: CalculateRealtimeRanking, CalculateWeeklyRanking, CalculateMonthlyRanking, CalculateNewbieRanking, UpdateRankings
- book_statistics_repository_mongo.go: UpdateRating, RemoveRating, BatchRecalculateStatistics

#### #010-B: Reader域Repository重构（P1）
- reading_progress_repository_mongo.go: SaveProgress, UpdateReadingTime, GetUnfinishedBooks, GetFinishedBooks
- collection_repository_mongo.go: validateCollectionTag

#### #010-C: Finance域Repository重构（P0）
- wallet_repository_mongo.go: UpdateBalance

#### #010-D: Writer域Repository重构（P1）
- project_repository_mongo.go: Create（默认值设置）
- document_repository_mongo.go: Create
- batch_operation_repository_mongo.go: UpdateItemStatus
- outline_repository_mongo.go: normalizeAndValidateOutlineQueryID

#### #010-E: Stats域Repository重构（P2）
- reader_behavior_repository_mongo.go: 所有Calculate*方法
- book_stats_repository_mongo.go: 所有Calculate*方法

#### #010-F: Social域Repository重构（P2）
- follow_repository_mongo.go: sanitizeFollowType, UpdateMutualStatus

#### #010-G: User域Repository重构（P2）
- user_repository_mongo.go: ValidateUser, GetActiveUsers

---

### Issue #011: 前后端数据类型不一致

**当前规模**: 28个问题（P0:5, P1:12, P2:11）

**建议拆分为**:

#### #011-A: 枚举值不一致问题（P0）
- BookStatus枚举值不匹配
- UserRole枚举命名不一致
- BehaviorType枚举不一致
- DocumentStatus枚举差异

#### #011-B: 字段类型和转换问题（P0）
- is_*布尔字段转换
- CategoryIDs数组类型
- 响应拦截器处理不一致

#### #011-C: 时间和金额处理问题（P1）
- Price字段类型和单位
- 时间字段命名不一致
- 时间格式处理不一致

#### #011-D: V2架构兼容性问题（P1）
- DocumentContent V2支持
- stableRef/orderKey字段

#### #011-E: 通用类型转换问题（P2）
- 分页参数命名
- Tags数组一致性
- 可选字段指针处理

---

## 优先级建议

### 第一优先级（立即处理）: P0问题组合
1. **架构重构组**: #007 → #010-A, #010-C
2. **枚举值问题**: #011-A
3. **API标准化组**: #005, #011

### 第二优先级（本周处理）: P1问题
1. **测试问题组**: #003, #009, #010-B, #010-D
2. **API标准化组**: #005, #011-B, #011-C
3. **认证权限组**: 已归档，转入后续治理项

### 第三优先级（持续优化）: P2问题
1. **代码质量**: #004
2. **其他Repository重构**: #010-E, #010-F, #010-G
3. **通用类型转换**: #011-E

---

## 更新计划

1. **拆分Issue**: 将 #010 和 #011 拆分为更小的Issue
2. **添加关联链接**: 在每个Issue中添加"相关Issue"和"依赖Issue"字段
3. **创建Issue看板**: 按优先级和状态组织Issue
4. **设置依赖关系**: 在项目管理工具中设置Issue依赖
