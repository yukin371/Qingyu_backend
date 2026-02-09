# Shared模块重构计划

**文档版本**: v1.0
**创建日期**: 2026-02-09
**任务**: ARCH-006模块边界收口

---

## 现状分析

### 当前shared模块结构

```
service/shared/
├── auth/           # 认证服务（15个文件）
├── cache/          # 缓存服务（1个文件）
├── messaging/      # 消息服务（7个文件）
├── metrics/        # 指标收集（待分析）
├── stats/          # 统计服务（1个文件）
├── storage/        # 存储服务（7个文件）
├── config_service.go
├── permission_service.go
└── README.md
```

### 模块职责分析

| 模块 | 职责 | 复杂度 | 依赖方 |
|------|------|--------|--------|
| **auth** | 用户认证、令牌管理、权限控制、OAuth、会话管理 | 高 | 多个业务服务、API层 |
| **messaging** | 消息队列、邮件通知、站内通知 | 中 | 少数服务 |
| **storage** | 文件存储、MinIO/本地存储适配 | 中 | 多个业务服务、API层 |
| **cache** | Redis缓存封装 | 低 | 几乎所有服务 |
| **stats** | 统计数据收集 | 低 | API层 |
| **metrics** | 指标数据收集 | 待分析 | 待分析 |

---

## 重构决策

### 原则

1. **业务独立模块** → 提升为独立模块
2. **基础设施模块** → 保留在shared
3. **通用工具模块** → 保留在shared

### 模块分类

#### ✅ 应该提升为独立模块

**1. auth → service/auth/**
- **理由**：
  - 核心业务功能，职责复杂
  - 被多个业务服务广泛依赖
  - 有独立的接口定义
  - 可以独立演进和测试

- **影响范围**：
  - service/container
  - service/interfaces/shared
  - api/v1/auth
  - api/v1/shared
  - middleware/auth

- **迁移难度**: 中等（已有Port接口）

**2. messaging → service/notification/**
- **理由**：
  - 业务相关功能（通知、消息）
  - 有独立的使用场景
  - 可以独立演进

- **影响范围**：
  - service/user
  - service/interfaces

- **迁移难度**: 低（依赖方少）

#### ✅ 应该保留在shared

**1. cache → service/shared/cache/**
- **理由**：
  - 基础设施组件
  - 被几乎所有服务依赖
  - 简单稳定的抽象

- **迁移难度**: 不适用（保留）

**2. storage → service/shared/storage/**
- **理由**：
  - 基础设施组件
  - 通用文件存储能力
  - 已有清晰的抽象

- **迁移难度**: 不适用（保留）

**3. stats → service/shared/stats/**
- **理由**：
  - 简单的统计功能
  - 使用范围有限

- **迁移难度**: 不适用（保留）

**4. metrics → service/shared/metrics/**
- **理由**：
  - 基础设施组件（待确认）

- **迁移难度**: 待分析

---

## 重构方案

### 阶段1: 准备工作（1天）

**目标**: 确保迁移基础就绪

- [x] 创建Port接口定义（已完成）
- [x] 创建依赖检查工具（已完成）
- [ ] 分析auth模块的完整依赖关系
- [ ] 分析messaging模块的完整依赖关系

**验收标准**:
- 依赖关系图谱完成
- 迁移风险识别完成

### 阶段2: 迁移auth模块（2-3天）

**目标**: 将auth提升为独立模块，但不破坏现有功能

#### 步骤1: 创建新目录结构

```
service/auth/
├── interfaces.go          # 从shared/auth/移入
├── services/
│   ├── auth_service.go   # 从shared/auth/移入
│   ├── jwt_service.go
│   ├── oauth_service.go
│   ├── session_service.go
│   ├── permission_service.go
│   ├── role_service.go
│   └── password_validator.go
├── adapters/
│   └── redis_adapter.go  # 从shared/auth/移入
├── README.md
└── _migration/           # 迁移兼容层
    └── shared_compat.go   # 重新导出到shared
```

#### 步骤2: 创建迁移兼容层

在`service/auth/_migration/shared_compat.go`中：

```go
// Package auth 提供认证服务
package auth

// 为了向后兼容，重新导出所有公共符号到shared路径
// 这样现有代码可以继续使用Qingyu_backend/service/shared/auth

// 当所有代码迁移完成后，可以删除这个兼容层
```

#### 步骤3: 更新import路径

- [ ] 更新service/container中的import
- [ ] 更新service/interfaces/shared中的import
- [ ] 更新api/v1/auth中的import
- [ ] 更新middleware中的import

#### 步骤4: 更新Port接口

- [ ] 修改service/interfaces/shared/adapters.go
- [ ] 从新路径导入auth服务

#### 步骤5: 测试验证

- [ ] 运行所有单元测试
- [ ] 运行集成测试
- [ ] 手动测试关键功能

**验收标准**:
- 所有测试通过
- API功能正常
- 性能无明显下降

### 阶段3: 迁移messaging模块（1-2天）

**目标**: 将messaging提升为独立的notification模块

#### 步骤1: 创建新目录结构

```
service/notification/
├── interfaces.go
├── services/
│   ├── email_service.go
│   ├── notification_service.go
│   └── messaging_service.go
├── queue/
│   └── redis_queue_client.go
└── README.md
```

#### 步骤2-5: 同auth模块的步骤

---

## 迁移兼容策略

### 短期兼容（迁移期）

在`service/shared/auth/`创建软链接或重新导出文件：

```go
// service/shared/auth/compat.go
package auth

// 为了向后兼容，重新导出新auth包的内容
// import "Qingyu_backend/service/auth" as newauth
// 这样现有代码可以继续使用Qingyu_backend/service/shared/auth

// TODO: 在v1.2中移除此兼容层
```

### 长期目标（清理阶段）

- [ ] 所有代码迁移到新import路径
- [ ] 删除兼容层
- [ ] 更新文档
- [ ] 更新CI规则，禁止使用旧路径

---

## 风险评估

### 高风险

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| import路径变更导致编译失败 | 高 | 使用兼容层，分阶段迁移 |
| 循环依赖引入 | 高 | 严格的依赖检查，Port接口隔离 |
| 运行时错误 | 中 | 完整的测试覆盖 |

### 中风险

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 性能下降 | 低 | 基准测试，性能验证 |
| 遗漏某些依赖 | 中 | 依赖分析工具检查 |

---

## 时间估算

| 阶段 | 任务 | 预计时间 |
|------|------|---------|
| 1 | 准备工作 | 1天 |
| 2 | 迁移auth模块 | 2-3天 |
| 3 | 迁移messaging模块 | 1-2天 |
| 4 | 测试验证 | 1天 |
| 5 | 清理兼容层 | 0.5天 |
| **总计** | | **6-8天** |

---

## 实施顺序

### 推荐顺序

1. **auth模块**（最重要，影响最广）
2. **messaging模块**（影响较小，作为验证）

### 不推荐的顺序

- ❌ 同时迁移多个模块
- ❌ 先迁移storage/cache（应该保留）
- ❌ 没有兼容层直接迁移

---

## 成功标准

### 每个模块迁移完成后

- [ ] 所有测试通过
- [ ] API功能正常
- [ ] 无性能下降
- [ ] 文档已更新

### 整体完成后

- [ ] shared模块只包含基础设施组件
- [ ] 业务模块提升为独立模块
- [ ] 依赖检查无违规
- [ ] 代码审查通过

---

## 回滚计划

如果迁移出现问题：

1. **立即回滚**: 使用git revert恢复到迁移前状态
2. **修复后重试**: 识别问题并修复后重新迁移
3. **分阶段回滚**: 如果是多个模块，可以单独回滚某个模块

---

## 相关文档

- [依赖规则文档](../architecture/dependency-rules.md)
- [接口契约文档](../../service/interfaces/shared/README.md)
- [迁移计划](../../plan/2026-02-09-migration-plan.md)
