# Issue 关联关系分析

**更新日期**: 2026-03-05
**Issue总数**: 13

---

## 关联矩阵

| Issue | 关联Issue | 关联类型 | 说明 |
|-------|----------|---------|------|
| #001 | #002, #011, #013 | 🔄 依赖 | ID类型问题需要联合处理 |
| #002 | #001, #013 | 🔄 依赖 | Create方法ID回设与ID类型相关 |
| #003 | #009, #013 | 🔗 相关 | 测试基础设施和覆盖率 |
| #004 | - | ⚪ 独立 | 代码质量问题 |
| #005 | #011 | 🔗 相关 | API响应格式与数据类型 |
| #006 | - | ⚪ 独立 | 数据库索引问题 |
| #007 | #010 | 🔄 依赖 | 事务管理与Repository业务逻辑 |
| #008 | #012 | 🔗 相关 | 中间件与权限认证 |
| #009 | #003, #013 | 🔗 相关 | 测试覆盖率与基础设施 |
| #010 | #001, #002, #007 | 🔄 依赖 | Repository层问题关联多个 |
| #011 | #001, #005 | 🔗 相关 | 数据类型与ID类型/API |
| #012 | #008 | 🔗 相关 | 权限配置与中间件 |
| #013 | #001, #002, #003, #009 | 🔄 依赖 | 种子数据关联ID和测试 |

---

## 需要联合处理的Issue组合

### 组合1: ID类型问题组 🔄

**关联Issue**: #001, #002, #011, #013

**核心问题**: ID字段类型不一致导致的查询和转换问题

**处理策略**: 按顺序处理
1. #001: 统一模型层 ID 字段类型（基础设施）
2. #002: 修复 Create 方法 ID 回设
3. #011: 修复前后端 ID 类型转换
4. #013: 修复测试用户种子数据 ID

**依赖关系**: #001 → #002 → #011 → #013

---

### 组合2: 测试问题组 🔗

**关联Issue**: #003, #009, #013

**核心问题**: 测试基础设施不足导致覆盖率低

**处理策略**: 并行处理
- #003: 测试基础设施改进
- #009: 测试覆盖率提升
- #013: 种子数据修复

**协同重点**:
- 修复测试数据问题（#013）后才能提升覆盖率（#009）
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

**关联Issue**: #007, #010

**核心问题**: Service层和Repository层职责不清

**处理策略**: 先#007后#010
1. #007: 实现事务管理器（基础设施）
2. #010: 将Repository层业务逻辑移到Service层

**依赖关系**: #007 → #010

---

### 组合5: 认证权限组 🔗

**关联Issue**: #008, #012

**核心问题**: 中间件架构和权限配置问题

**处理策略**: 联合处理
- #008: 修正中间件目录结构、CORS位置、限流统一
- #012: 修复权限配置、实现权限检查中间件

**协同重点**:
- 权限中间件实现（#012）需要中间件架构支持（#008）
- CORS位置修复（#008）影响认证流程（#012）

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
1. **ID类型问题组**: #001 → #002 → #013
2. **架构重构组**: #007 → #010-A, #010-C
3. **枚举值问题**: #011-A

### 第二优先级（本周处理）: P1问题
1. **测试问题组**: #003, #009, #010-B, #010-D
2. **API标准化组**: #005, #011-B, #011-C
3. **认证权限组**: #008, #012

### 第三优先级（持续优化）: P2问题
1. **代码质量**: #004
2. **数据库索引**: #006
3. **其他Repository重构**: #010-E, #010-F, #010-G
4. **通用类型转换**: #011-E

---

## 更新计划

1. **拆分Issue**: 将 #010 和 #011 拆分为更小的Issue
2. **添加关联链接**: 在每个Issue中添加"相关Issue"和"依赖Issue"字段
3. **创建Issue看板**: 按优先级和状态组织Issue
4. **设置依赖关系**: 在项目管理工具中设置Issue依赖
