# 核心功能设计文档

> 最后整理: 2026-05-22  
> 当前状态: `legacy-live / summary-draft`

本文档是 `docs/design/core/` 的早期核心能力汇总稿，用来快速回看认证、架构、配置、事务、事件等主题曾如何被整体描述。它不再是当前目录的标准入口。

## Page Role

- 这里负责：核心基础设施历史设计主题的汇总导航。
- 不负责：当前架构 owner、当前实现状态、当前目录标准入口。

## Recommended Read Path

1. [README.md](./README.md)
2. [../../architecture/README.md](../../architecture/README.md)
3. [../../standards/README.md](../../standards/README.md)
4. [../../database/README.md](../../database/README.md)

## Boundary

- 当前 `core/` 目录标准入口是 [README.md](./README.md)，不是本文件。
- 如果你要找“当前后端架构与模块边界”，优先看 [../../architecture/README.md](../../architecture/README.md)。
- 如果你要找“当前规则与分层约束”，优先看 [../../standards/README.md](../../standards/README.md)。
- 如果你要找数据库迁移与操作专题，请回到 [../../database/README.md](../../database/README.md)。

> **架构版本**: v2.1  
> **最后更新**: 2025-10-21  
> **模块定位**: 青羽平台核心基础设施

## 📋 概述

核心功能模块包含青羽平台的基础设施设计，涵盖认证、架构、配置、事务、事件驱动等核心机制。

## 📚 文档索引

### 身份认证与授权

#### [JWT身份认证设计](./auth/JWT身份认证设计.md)
**功能范围**：
- JWT token生成和验证
- token刷新机制
- 权限验证

**状态**：✅ 已完成

### 架构设计

#### [Repository层与Service层架构重新设计](../shared/README.md)
**功能范围**：
- Repository模式设计
- Service层设计规范
- 依赖注入
- 接口优先原则

**状态**：✅ 已完成

#### [后端架构与分层入口](../../architecture/README.md)
**文档类型**：架构重构总结  
**内容**：Repository和Service层重构的总结报告  
**状态**：✅ 已完成

### v2.1新增架构设计 ⭐

#### [Saga事务模式设计](./saga/Saga事务模式设计.md)
**功能范围**：
- Saga事务编排
- 补偿机制设计
- 关键业务场景Saga（用户购买章节、作者提现）
- Saga Log设计
- 最终一致性保证

**v2.1架构特性**：
- ✅ 分布式事务一致性保证
- ✅ 自动补偿机制
- ✅ 事务状态管理

**应用场景**：
```
场景：用户购买章节
1. 扣除代币（Wallet Service）
2. 解锁章节（Bookstore Service）
3. 增加作者收益（Monetization Service）

如果步骤2或3失败，自动执行补偿：
- 补偿操作：返还代币
```

**状态**：✅ v2.1设计完成

#### [事件驱动架构设计](./event/事件驱动架构设计.md)
**功能范围**：
- EventBus增强设计
- 强制事件方法定义
- 核心业务事件
- EventHandler注册机制

**v2.1架构特性**：
- ✅ Service接口强制事件方法
- ✅ 异步事件发布
- ✅ 事件订阅和处理

**核心事件**：
- `task.completed` - 任务完成
- `user.level_up` - 用户升级
- `document.version_created` - 文档版本创建
- `chapter.purchased` - 章节购买

**状态**：✅ v2.1设计完成

### 项目管理

#### [项目四层架构CRUD设计](./project/项目四层架构CRUD设计.md) ⭐
**功能范围**：
- **四层架构**：Project → Node → Document → Content
- GridFS大文本支持（>1MB）
- 树形结构管理（Path字段加速查询）
- Document/Content分离（性能优化20-30倍）
- 乐观锁版本控制

**v2.1架构特性**：
- ✅ 支持超大文本（>10MB）
- ✅ 查询文档树 <300ms（100个节点）
- ✅ 按需加载内容
- ✅ 自动选择存储方式（小文件直接存储，大文件GridFS）

**状态**：✅ v2.1版本已完成

#### [项目路径管理设计](./project/项目路径管理设计.md)
**功能范围**：
- 项目路径管理
- 路径组织规范

**状态**：✅ 已完成

#### 项目CRUD.md - ⚠️ 已过时
> 本文档已被更详细的设计取代。  
> **替代文档**：[项目四层架构CRUD设计](./project/项目四层架构CRUD设计.md)

#### 项目\_节点\_文档\_CRUD\_Design.md - ⚠️ 已过时
> 本文档基于三层架构，已被四层架构取代。  
> **替代文档**：[项目四层架构CRUD设计](./project/项目四层架构CRUD设计.md)
> **过时原因**：缺少Content层、不支持GridFS大文本存储

### 版本控制

#### [版本控制](./project/版本控制.md)
**功能范围**：
- 文档版本管理
- 历史记录
- 版本对比和回滚

**状态**：✅ 已完成

### 系统配置

#### [配置系统设计](./config/配置系统设计.md)
**功能范围**：
- 配置管理
- 环境变量
- 动态配置

**状态**：✅ 已完成

### 待办与搁置功能

#### [核心功能设计待办清单](./核心功能设计待办清单.md)
**文档类型**：待办清单  
**状态**：📝 维护中

#### [协作编辑系统（搁置）](./核心功能设计待办清单.md)
**功能范围**：实时协作编辑  
**状态**：⏸️ 已搁置

#### [插件扩展系统(搁置)](./_archived/)
**功能范围**：插件系统设计  
**状态**：⏸️ 已搁置

## 🎯 v2.1架构要点

### 事件驱动架构

所有核心业务流程采用事件驱动：

**强制事件方法示例**：
```go
type UserService interface {
    // 业务方法
    CreateUser(...)
    
    // 强制事件方法
    OnUserCreated(ctx context.Context, event *UserCreatedEvent) error
    OnTaskCompleted(ctx context.Context, event *TaskCompletedEvent) error
}
```

### Saga事务模式

跨服务事务采用Saga模式保证最终一致性：

**Saga流程**：
1. 执行业务步骤
2. 记录Saga Log
3. 如果失败，执行补偿操作
4. 保证最终一致性

### 设计原则

1. **单一职责原则**：每个模块只负责一个明确的功能领域
2. **依赖倒置原则**：高层模块不依赖低层模块，通过接口解耦
3. **开闭原则**：对扩展开放，对修改关闭
4. **接口优先原则**：所有跨层依赖必须通过接口
5. **事件驱动原则**：核心业务必须发布事件

## 🔗 相关文档

### 架构总览
- [后端架构总览](../../architecture/README.md)

### 重构规划
- [事件驱动改造清单](../重构规划/事件驱动改造清单.md) - 需要添加事件的Service方法清单
- [架构重构总结v3.0](../重构规划/架构重构总结v3.0.md) - 架构升级总结

### 相关模块
- [数据库设计](../database/README.md) - Repository层实现
- [中间件设计](../middleware/README.md) - 认证、日志、错误处理
- [安全设计](../security/README.md) - 安全策略

## 📞 联系方式

**文档维护者**：青羽架构组  
**最后更新**：2025-10-21  
**架构版本**：v2.1
