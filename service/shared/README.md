# 共享底层服务模块

> 青羽后端 - 模块化单体架构的共享服务层
>
> **创建时间**: 2025-09-30  
> **当前进度**: 阶段1 已完成 ✅

---

## 📋 概述

共享底层服务是青羽后端的基础设施层，为阅读端和写作端提供统一的核心功能。

**架构模式**: 模块化单体（Modular Monolith）

---

## 🏗️ 模块结构

### 已实现模块

| 模块 | 状态 | 说明 | 接口文件 |
|------|------|------|---------|
| **Auth** | 📝 接口定义 | 认证与权限管理 | [interfaces.go](./auth/interfaces.go) |
| **Wallet** | 📝 接口定义 | 钱包与交易系统 | [interfaces.go](./wallet/interfaces.go) |
| **Recommendation** | 📝 接口定义 | 推荐服务 | [interfaces.go](./recommendation/interfaces.go) |
| **Messaging** | 📝 接口定义 | 消息队列 | [interfaces.go](./messaging/interfaces.go) |
| **Storage** | 📝 接口定义 | 文件存储 | [interfaces.go](./storage/interfaces.go) |
| **Admin** | 📝 接口定义 | 管理后台 | [interfaces.go](./admin/interfaces.go) |

---

## 📦 目录结构

```
service/shared/
├── auth/                      # 认证与权限模块
│   └── interfaces.go          # ✅ 已创建
│
├── wallet/                    # 钱包系统模块
│   └── interfaces.go          # ✅ 已创建
│
├── recommendation/            # 推荐服务模块
│   └── interfaces.go          # ✅ 已创建
│
├── messaging/                 # 消息队列模块
│   └── interfaces.go          # ✅ 已创建
│
├── storage/                   # 文件存储模块
│   └── interfaces.go          # ✅ 已创建
│
├── admin/                     # 管理后台模块
│   └── interfaces.go          # ✅ 已创建
│
└── README.md                  # 本文档
```

---

## 🎯 各模块功能

### 1. Auth 模块

**核心功能**:
- ✅ JWT Token 生成与验证
- ✅ 用户注册与登录
- ✅ 基于 RBAC 的权限系统
- ✅ 会话管理（Redis）
- ✅ 角色与权限管理

**接口数量**: 4个服务接口 + 多个请求/响应结构

**文档**: [账号权限系统设计](../../doc/design/shared/账号权限系统设计.md)

---

### 2. Wallet 模块

**核心功能**:
- ✅ 钱包创建与管理
- ✅ 充值、消费、转账
- ✅ 交易记录查询
- ✅ 提现申请与审核
- ✅ 支付集成（预留接口）

**接口数量**: 1个主服务接口 + 多个数据结构

**文档**: [钱包系统设计](../../doc/design/shared/钱包系统设计.md)

---

### 3. Recommendation 模块

**核心功能**:
- ✅ 个性化推荐
- ✅ 相似内容推荐
- ✅ 热门内容推荐
- ✅ 用户行为收集
- ✅ 推荐结果缓存

**接口数量**: 1个主服务接口

**文档**: [推荐服务设计](../../doc/design/shared/推荐服务设计.md)

---

### 4. Messaging 模块

**核心功能**:
- ✅ 消息发布订阅
- ✅ 延迟消息
- ✅ Topic 管理
- ✅ 基于 Redis Streams

**接口数量**: 1个主服务接口 + 预定义Topic常量

**文档**: [消息队列设计](../../doc/design/shared/消息队列设计.md)

---

### 5. Storage 模块

**核心功能**:
- ✅ 文件上传下载
- ✅ 文件权限控制
- ✅ 文件元数据管理
- ✅ 临时访问URL生成
- ✅ 本地存储 + 云存储支持

**接口数量**: 1个主服务接口

**文档**: [文件存储设计](../../doc/design/shared/文件存储设计.md)

---

### 6. Admin 模块

**核心功能**:
- ✅ 内容审核
- ✅ 用户管理（封禁、解封）
- ✅ 提现审核
- ✅ 操作日志记录
- ✅ 统计数据查询

**接口数量**: 1个主服务接口

**文档**: [管理后台设计](../../doc/design/shared/管理后台设计.md)

---

## 📊 实施进度

### 阶段 1: 目录结构创建 ✅ (100%)

**已完成**:
- [x] 创建所有服务模块目录
- [x] 创建所有模型目录
- [x] 创建所有Repository接口目录
- [x] 创建6个服务接口文件
- [x] 创建6个核心模型文件
- [x] 创建Repository接口定义

**完成时间**: 2025-09-30

---

### 阶段 2: Auth 模块基础功能 ⏸️ (0%)

**待实现**:
- [ ] JWT服务实现
- [ ] 角色模型与Repository
- [ ] 角色服务实现
- [ ] 权限服务实现

**预计工作量**: 8小时

---

### 阶段 3-12: 待开始

详见: [共享底层服务实施计划](../../doc/implementation/共享底层服务实施计划.md)

---

## 🔧 使用说明

### 导入服务接口

```go
import (
    authService "Qingyu_backend/service/shared/auth"
    walletService "Qingyu_backend/service/shared/wallet"
    // ... 其他模块
)
```

### 调用示例（待实现后更新）

```go
// 示例：用户注册
req := &authService.RegisterRequest{
    Username: "user123",
    Email: "user@example.com",
    Password: "password123",
}

response, err := authSvc.Register(ctx, req)
if err != nil {
    // 错误处理
}

// 获取Token
token := response.Token
```

---

## 📂 相关文件

### 模型文件

```
models/shared/
├── auth/
│   ├── role.go              # ✅ 角色与权限模型
│   └── session.go           # ✅ 会话模型
│
├── wallet/
│   └── wallet.go            # ✅ 钱包与交易模型
│
├── recommendation/
│   └── recommendation.go    # ✅ 推荐与行为模型
│
├── storage/
│   └── file.go              # ✅ 文件存储模型
│
└── admin/
    └── admin.go             # ✅ 审核与日志模型
```

### Repository 接口

```
repository/interfaces/shared/
└── shared_repository.go     # ✅ 所有Repository接口定义
```

---

## 🔗 依赖关系

### 模块依赖图

```
┌─────────────┐
│   Messaging │ (基础设施)
└──────┬──────┘
       │
┌──────▼──────┐     ┌────────────┐
│    Auth     │◄────┤   Storage  │
└──────┬──────┘     └────────────┘
       │
       ├──────────────────┐
       │                  │
┌──────▼──────┐     ┌─────▼─────────┐
│   Wallet    │     │ Recommendation│
└─────────────┘     └───────────────┘
       │                  │
       └──────────┬───────┘
                  │
           ┌──────▼──────┐
           │    Admin    │
           └─────────────┘
```

### 实施优先级

1. **优先级1（基础设施）**: Auth, Messaging
2. **优先级2（核心业务）**: Wallet, Storage
3. **优先级3（增值服务）**: Recommendation, Admin

---

## 📚 相关文档

### 设计文档
- [共享底层服务设计文档](../../doc/design/shared/README_共享底层服务设计文档.md)
- [账号权限系统设计](../../doc/design/shared/账号权限系统设计.md)
- [钱包系统设计](../../doc/design/shared/钱包系统设计.md)
- [推荐服务设计](../../doc/design/shared/推荐服务设计.md)
- [消息队列设计](../../doc/design/shared/消息队列设计.md)

### 实施文档
- [共享底层服务实施计划](../../doc/implementation/共享底层服务实施计划.md)
- [架构设计规范](../../doc/architecture/架构设计规范.md)
- [项目开发规则](../../doc/architecture/项目开发规则.md)

---

## 🚀 下一步计划

### 本周目标 (Week 1)
1. 完成 Auth 模块基础功能（JWT + 角色权限）
2. 完成 Auth Repository 实现
3. 编写 Auth 模块单元测试

### 本月目标 (Month 1)
1. 完成 Auth 和 Wallet 两大核心模块
2. 完成基础 API 层
3. 完成基础集成测试

---

## ⚠️ 注意事项

1. **模块独立性**: 各模块通过接口交互，禁止跨模块直接调用内部实现
2. **接口稳定性**: 接口一旦定义应保持稳定，变更需谨慎
3. **错误处理**: 统一使用项目错误处理规范
4. **日志记录**: 关键操作必须记录日志
5. **测试覆盖**: 核心功能必须有单元测试

---

## 📞 协作指南

### 如何贡献

1. 选择一个待实现的模块或功能
2. 阅读对应的设计文档
3. 实现服务接口
4. 编写单元测试
5. 更新文档

### Commit 规范

```
feat(shared/auth): 实现JWT服务
fix(shared/wallet): 修复余额计算错误
test(shared/recommendation): 添加推荐算法测试
docs(shared): 更新使用说明
```

---

*持续更新中... 📝*
