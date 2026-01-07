# 实施文档目录

> 青羽后端项目的实施文档、进度跟踪和测试报告

---

## 新模块化文档结构

文档已按照模块化设计重新组织，与 `doc/design/modules` 结构保持一致。

### 模块文档目录

| 模块 | 说明 | 状态 | 文档路径 |
|------|------|------|----------|
| **01-auth** | 认证/权限控制 | ✅ 完成 | [01-auth/](./01-auth/) |
| **02-writing** | 写作工具/导出发布 | ✅ 完成 | [02-writing/](./02-writing/) |
| **03-reading** | 阅读/统计系统 | ✅ 完成 | [03-reading/](./03-reading/) |
| **04-social** | 社交/书单功能 | ✅ 完成 | [04-social/](./04-social/) |
| **05-communication** | 通信/通知系统 | ✅ 完成 | [05-communication/](./05-communication/) |
| **06-bookstore** | 书城系统 | ✅ 完成 | [06-bookstore/](./06-bookstore/) |
| **07-admin** | 管理后台 | - | [07-admin/](./07-admin/) |
| **08-finance** | 财务/支付 | ✅ 完成 | [08-finance/](./08-finance/) |
| **09-ai** | AI服务 | - | [09-ai/](./09-ai/) |

### 跨模块文档

| 分类 | 说明 | 文档路径 |
|------|------|----------|
| **infrastructure** | 基础设施/中间件/MCP | [infrastructure/](./infrastructure/) |
| **testing** | 测试报告 | [testing/](./testing/) |
| **docs** | 项目文档/启动指南 | [docs/](./docs/) |

---

## 快速导航

### 核心实施文档

#### 认证与权限
- [RBAC 权限控制实施](./01-auth/RBAC_IMPLEMENTATION.md)
  - 基于角色的访问控制系统
  - 权限检查中间件实现
  - 角色管理和用户授权

#### 写作功能
- [导出和发布管理 API](./02-writing/WRITING_FEATURES_API.md)
  - 文档导出功能 (TXT/MD/DOCX)
  - 项目导出为 ZIP
  - 发布到书城平台

#### 阅读功能
- [阅读统计模块实施](./03-reading/READING_STATS_IMPLEMENTATION.md)
  - 章节/作品统计
  - 读者行为记录
  - 热力图和趋势分析

#### 社交功能
- [书单系统模块实施](./04-social/BOOKLIST_MODULE_IMPLEMENTATION.md)
  - 书单管理
  - 热门书单推荐
  - 书单收藏功能

#### 通信功能
- [通知系统实施](./05-communication/NOTIFICATION_SYSTEM_IMPLEMENTATION.md)
  - 消息通知系统
  - WebSocket 实时推送
- [通知 API 快速参考](./05-communication/NOTIFICATION_API_QUICK_REFERENCE.md)

#### 书城系统
- [书城 API 实施总结](./06-bookstore/BOOKSTORE_API_IMPLEMENTATION_SUMMARY.md)
  - 书籍管理
  - 章节系统
- [章节购买 API](./06-bookstore/BOOKSTORE_CHAPTER_PURCHASE_API.md)
- [书城测试报告](./06-bookstore/BOOKSTORE_TEST_REPORT.md)

#### 财务系统
- [支付系统实施](./08-finance/FINANCE_PAYMENT_SYSTEM_IMPLEMENTATION.md)
  - 会员系统
  - 作者收入
  - 支付集成

### 基础设施文档

- [P0 中间件集成](./infrastructure/MIDDLEWARE_INTEGRATION.md)
  - RequestID、Recovery、Logger
  - Metrics、RateLimit、ErrorHandler、CORS
- [路由冲突修复](./infrastructure/ROUTER_FIX_REPORT.md)
  - 重复路由问题诊断
  - 路由架构优化
- [MCP 服务器配置](./infrastructure/MCP_SERVERS.md)
- [浏览器 MCP 设置](./infrastructure/BROWSER_MCP_SETUP.md)

### 测试报告

- [前端功能测试](./testing/FRONTEND_FUNCTIONALITY_TEST_REPORT.md)
- [搜索功能验证](./testing/SEARCH_FUNCTIONALITY_VERIFICATION_REPORT.md)
- [前端搜索修复](./testing/FRONTEND_SEARCH_FIX_REPORT.md)
- [搜索编码修复](./testing/SEARCH_ENCODING_FIX_REPORT.md)
- [API 缺失报告](./testing/MISSING_API_REPORT.md)
- [共享 API 修复](./testing/SHARED_API_FIX.md)

### 项目文档

- [启动服务指南](./docs/START_SERVICES.md)
- [项目结构总结](./docs/项目结构总结.md)

---

## 原有进度指导目录

原有的中文目录结构保持不变，包含详细的进度跟踪和历史记录。

### 进度指导

- [00进度指导/](./00进度指导/) - 各阶段进度跟踪报告

### 基础设施

- [01基础设施/](./01基础设施/) - 基础设施实施文档

### 共享服务

- [02共享底层服务/](./02共享底层服务/) - Auth、Wallet、Storage 等

### 用户管理

- [03用户管理模块/](./03用户管理模块/) - 用户管理实施文档

### 写作模块

- [04写作端模块/](./04写作端模块/) - 写作端实施文档

### AI 服务

- [05AI服务模块/](./05AI服务模块/) - AI 服务实施文档

### 阅读模块

- [06阅读端模块/](./06阅读端模块/) - 阅读端实施文档

### 测试实施

- [测试实施/](./测试实施/) - 测试实施报告

### 修复记录

- [修复/](./修复/) - 问题修复记录

---

## Git 提交历史

### 最近提交

```
73e1ba6 - docs: 整理实施文档并创建进度指导文档
dd5f79a - fix(router): 修复重复路由注册冲突
3c738ae - feat(service): 启用 ReadingStatsService 并注册 reading-stats 路由
a5a8762 - feat(service): ReadingStatsService 实现 BaseService 接口
4129a26 - feat(repository): 在 RepositoryFactory 中添加 stats 模块仓储
```

### 关键功能提交

| 功能 | 提交 | 说明 |
|------|------|------|
| Stats 模块 | 4129a26, a5a8762, 3c738ae | 仓储、BaseService、路由注册 |
| BookList 模块 | 28018e7, 4c72d79, bb97a60 | MongoDB 仓储、工厂方法、路由 |
| P0 中间件 | 9fe9a01, 18f2e55 | 核心中间件集成 |
| RBAC 权限 | 11ad486 | 权限控制系统 |

---

## 模块对应关系

### 新模块 ↔ 旧目录映射

| 新模块 | 对应旧目录 | 功能 |
|--------|-----------|------|
| 01-auth | 03用户管理模块 | 认证、权限、用户管理 |
| 02-writing | 04写作端模块 | 写作工具、导出发布 |
| 03-reading | 06阅读端模块 | 阅读、统计、推荐 |
| 04-social | 04写作端模块(社交部分) | 书单、评论、点赞 |
| 05-communication | 02共享底层服务(Messaging) | 通知、消息 |
| 06-bookstore | 06阅读端模块(书城) | 书城、章节、购买 |
| 07-admin | 03用户管理模块(管理) | 后台管理 |
| 08-finance | - | 财务、支付 |
| 09-ai | 05AI服务模块 | AI 服务 |

---

## 文档使用指南

### 查找实施文档

1. **按模块查找**: 使用上方的模块文档目录表
2. **按功能查找**: 使用快速导航中的功能分类
3. **按问题查找**: 查看测试报告和修复记录

### 更新文档

1. 新模块实施文档放在对应的 `XX-功能名/` 目录
2. 测试报告放在 `testing/` 目录
3. 基础设施相关放在 `infrastructure/` 目录

### 文档命名规范

- 实施文档: `FEATURE_MODULE_IMPLEMENTATION.md`
- 测试报告: `FEATURE_TEST_REPORT.md`
- 修复报告: `ISSUE_FIX_REPORT.md`
- 快速参考: `FEATURE_API_QUICK_REFERENCE.md`

---

*文档最后更新: 2026-01-07*
