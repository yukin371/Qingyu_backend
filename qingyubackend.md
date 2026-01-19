---
description: Qingyu Backend 项目开发规则和分层架构指导 v2.1
globs:
alwaysApply: true
---

# Qingyu Backend 项目开发规则 v2.1

> **快速参考指南**：本文档提供核心开发规则和架构原则概览。详细设计和实现请参考 `docs/` 目录下的文档。

## 项目介绍

**项目名称**：青羽写作后端服务（Qingyu Backend）

青羽写作后端服务是为 AI 辅助写作应用提供核心支持的高性能 API 项目。采用现代化分层架构和设计模式，支持用户管理、文档存储、AI 文本生成、阅读社区等核心功能。

### 技术栈

**核心技术**
- **编程语言**：Go 1.21+
- **Web框架**：Gin（RESTful API）
- **数据库**：MongoDB + 支持多数据库扩展
- **认证授权**：JWT + RBAC
- **配置管理**：Viper（多环境配置）
- **日志系统**：Zap（结构化日志）
- **缓存**：Redis
- **容器化**：Docker + Docker Compose

**架构模式**
- **分层架构**：Router → API → Service → Repository → Model
- **依赖注入**：接口驱动、服务容器管理
- **事件驱动**：EventBus
- **工厂模式**：Repository 工厂
- **策略模式**：多提供商适配器

## 核心架构原则

### 1. 分层架构原则

项目采用严格的五层架构，各层职责清晰，单向依赖：

```
Router → API → Service → Repository → Model
```

**核心规则**：
- ✅ 上层依赖下层（通过接口）
- ❌ 下层不能依赖上层
- ✅ 同层通过接口交互
- ✅ 依赖接口而非具体实现

### 2. 依赖注入原则

- 所有服务通过构造函数注入依赖
- 依赖接口而非具体实现
- 使用 ServiceContainer 管理服务生命周期
- 通过 RepositoryFactory 创建数据访问层

### 3. 接口优先原则

- 所有跨层依赖必须定义接口
- Repository 层提供数据库无关的抽象
- Service 层定义业务能力接口
- 便于 Mock 测试和多实现扩展

> **详细说明**：参见 `docs/architecture/架构设计规范.md`

---

## 各层级职责说明

### 1. Model 层 (`models/`)

**核心职责**：
- ✅ 定义数据结构和字段标签（bson, json, validate）
- ✅ 提供基础数据方法（格式化、转换）
- ❌ 不包含业务逻辑、数据库操作、HTTP 处理

**组织结构**：按业务域分包（users, document, ai, reading, shared 等）

> **详细说明**：参见实际代码 `models/` 目录

---

### 2. Repository 层 (`repository/`)

**核心职责**：
- ✅ 数据持久化和查询封装
- ✅ 事务管理和缓存策略
- ✅ 提供数据库无关的接口抽象
- ❌ 不包含业务逻辑、HTTP 处理

**组织结构**：
- `interfaces/` - 接口定义（按业务域分包）
- `mongodb/` - MongoDB 具体实现
- `Mock/` - Mock 实现（测试用）

**核心模式**：
- **Factory 模式**：通过 RepositoryFactory 创建实例
- **接口抽象**：基础 CRUD + 业务特定查询
- **QueryBuilder**：复杂查询条件构建器

**最佳实践**：
- ✅ 所有数据库操作必须在 Repository 层
- ✅ 使用接口定义，便于多实现和测试
- ✅ 统一错误处理，支持健康检查
- ❌ 不处理业务逻辑

> **详细说明**：参见 `docs/architecture/repository层设计规范.md`

---

### 3. Service 层 (`service/`)

**核心职责**：
- ✅ 业务逻辑处理和规则验证
- ✅ 流程控制和事务协调
- ✅ 事件发布（EventBus）
- ✅ 调用 Repository 进行数据操作
- ❌ 不直接操作数据库、不处理 HTTP 请求、不操作 gin.Context

**组织结构**：
- `base/` - 基础服务接口
- `interfaces/` - Service 接口定义
- `container/` - 服务容器（依赖注入）
- 业务服务包：user, ai, project, bookstore, reading, shared 等

**核心模式**：
- **BaseService 接口**：Initialize, Health, Close, GetServiceName, GetVersion
- **依赖注入**：通过构造函数注入 Repository 接口和 EventBus
- **ServiceContainer**：管理服务注册、获取和初始化
- **事件驱动**：业务操作后发布事件

**最佳实践**：
- ✅ 所有业务逻辑必须在 Service 层
- ✅ 通过 Repository 接口访问数据
- ✅ 使用 EventBus 发布业务事件
- ✅ 统一的参数验证和错误处理
- ❌ 不直接操作数据库和 HTTP

> **详细说明**：参见 `docs/architecture/架构设计规范.md`

---

### 4. API 层 (`api/v1/`)

**核心职责**：
- ✅ HTTP 请求处理：参数绑定和验证
- ✅ 调用 Service 层处理业务
- ✅ 构建 HTTP 响应和错误转换
- ❌ 不包含业务逻辑、不直接操作数据库、不直接调用 Repository

**组织结构**：按业务域分包（system, document, ai, reading, reader, recommendation, shared, writer 等）

**核心模式**：
- **统一响应格式**：Success, Error, ValidationError
- **错误转换**：Service 层错误 → HTTP 状态码
- **上下文获取**：从 gin.Context 获取用户信息（由中间件注入）

**最佳实践**：
- ✅ 只处理 HTTP 协议相关逻辑
- ✅ 参数绑定后调用 Service 层
- ✅ 统一响应格式和错误处理
- ❌ 不包含业务逻辑

> **详细说明**：参见 `docs/api/API设计规范.md`

---

### 5. Router 层 (`router/`)

**核心职责**：
- ✅ 路由定义和注册
- ✅ 中间件配置（认证、权限、限流等）
- ✅ API 版本管理和路由分组
- ❌ 不处理业务逻辑、不处理请求参数

**组织结构**：
- `enter.go` - 路由入口和全局中间件
- 业务路由包：users, project, ai, bookstore, reading, shared, writer 等

**核心模式**：
- **路由分组**：public（公开）、authenticated（需认证）、admin（管理员）
- **中间件链**：Logger → CORS → Recovery → JWTAuth → Permission
- **RESTful 风格**：GET/POST/PUT/DELETE + 资源路径

**最佳实践**：
- ✅ 清晰的路由分组和版本管理
- ✅ RESTful 风格的路径设计
- ✅ 合理的中间件顺序
- ✅ 基于角色的访问控制
- ❌ 不处理业务逻辑

> **详细说明**：参见 `docs/architecture/路由层设计规范.md`

---

## 统一错误处理

**错误分类**：
- `CategoryValidation` - 参数验证错误（400）
- `CategoryBusiness` - 业务逻辑错误（409/404）
- `CategoryAuth` - 认证授权错误（401/403）
- `CategorySystem` - 系统错误（500）
- `CategoryDatabase` - 数据库错误
- `CategoryExternal` - 外部服务错误

**UnifiedError 结构**：包含 Code, Category, Message, Details, Service, Operation, HTTPStatus 等字段

**创建函数**：
- `NewValidationError()` - 参数验证错误
- `NewBusinessError()` - 业务错误
- `NewNotFoundError()` - 资源未找到
- `NewInternalError()` - 内部错误
- `NewAuthError()` - 认证错误

**使用规范**：
- **Repository 层**：返回基础错误或特定错误类型
- **Service 层**：使用 UnifiedError 包装，添加业务上下文
- **API 层**：转换为 HTTP 状态码和响应

> **详细说明**：参见 `pkg/errors/` 包实现

---

## 事件驱动架构

**核心接口**：
- `Event` - 事件接口（GetEventType, GetEventData, GetTimestamp, GetSource）
- `EventHandler` - 事件处理器接口（Handle, GetHandlerName, GetSupportedEventTypes）
- `EventBus` - 事件总线接口（Subscribe, Unsubscribe, Publish, PublishAsync）

**使用流程**：
1. **定义事件**：定义事件数据结构
2. **发布事件**：Service 层在业务操作后通过 EventBus 发布
3. **订阅事件**：实现 EventHandler 接口
4. **注册处理器**：在应用启动时注册到 EventBus

**典型场景**：
- 用户注册后发送欢迎邮件
- 订单创建后更新库存
- 文档更新后同步搜索索引

> **详细说明**：参见 `service/base/` 包实现

---

## 开发规范

### 1. 命名规范

**文件命名**：
- 小写字母 + 下划线：`user_service.go`、`user_repository_mongo.go`
- 测试文件：`user_service_test.go`
- 接口文件：`UserRepository_interface.go`

**包名/结构体/方法命名**：
- **包名**：小写单数（user, project, ai）
- **结构体**：大驼峰（UserService, MongoUserRepository, UserApi）
- **接口**：名词 + 功能（UserRepository, UserService），不用 Interface 后缀
- **方法**：动词开头（CreateUser, GetUserByID），私有方法小写开头（validateUser）
- **变量**：小驼峰（userRepo, eventBus），缩写一致（userID 不是 userId）

> **详细规范**：参见 `docs/engineering/软件工程规范_v2.0.md`

### 2. 注释规范

- **包注释**：说明包的功能和用途
- **结构体注释**：说明结构体职责
- **方法注释**：说明功能、参数、返回值

### 3. 错误处理规范

- **Repository 层**：使用基础错误或 infrastructure 错误类型
- **Service 层**：使用 UnifiedError 包装，添加上下文
- **API 层**：转换为 HTTP 响应

### 4. 测试规范

- **单元测试**：使用 Mock Repository 和 EventBus
- **集成测试**：使用真实 MongoDB 连接
- **测试组织**：按层级组织（`test/repository/`, `test/service/`, `test/api/`）

> **详细说明**：参见 `test/` 目录和 `docs/testing/`

### 5. 日志规范

使用结构化日志（logrus），包含 service, operation, request_id 等字段

> **详细说明**：参见 `docs/engineering/软件工程规范_v2.0.md`

---

## 项目初始化流程

1. **加载配置**：使用 Viper 加载多环境配置
2. **初始化数据库**：创建 MongoDB 连接
3. **创建 Repository 工厂**：初始化 RepositoryFactory
4. **创建服务容器**：初始化 ServiceContainer
5. **注册服务**：通过工厂创建 Repository，注入到 Service，注册到容器
6. **初始化服务**：调用 container.Initialize()
7. **初始化 Gin 引擎**：配置中间件
8. **注册路由**：调用 router.InitRoutes()
9. **启动服务器**：运行 HTTP 服务器

> **详细说明**：参见 `cmd/server/main.go`

---

## 架构检查清单

### Repository 层
- [ ] 实现统一 Repository 接口
- [ ] 使用 Repository 工厂模式
- [ ] 不包含业务逻辑
- [ ] 统一错误处理和健康检查
- [ ] 查询使用 QueryBuilder 封装

### Service 层
- [ ] 实现 BaseService 接口
- [ ] 使用依赖注入，通过 Repository 接口访问数据
- [ ] 不直接操作数据库
- [ ] 统一参数验证和错误处理
- [ ] 发布业务事件，支持单元测试

### API 层
- [ ] 只处理 HTTP 请求响应
- [ ] 不包含业务逻辑
- [ ] 参数绑定后调用 Service 层
- [ ] 统一响应格式和错误转换
- [ ] 不直接调用 Repository

### Router 层
- [ ] 清晰的路由分组和 RESTful 风格
- [ ] 合理的中间件顺序
- [ ] 基于角色的访问控制
- [ ] 不处理业务逻辑

### 代码质量
- [ ] 遵循命名规范和注释规范
- [ ] 统一错误处理
- [ ] 单元测试和集成测试覆盖
- [ ] 代码审查通过

---

## 常见问题解答

**Q1: 为什么使用 Repository 模式？**
- 业务逻辑与数据访问分离
- 便于单元测试（Mock Repository）
- 支持多数据库实现
- 提高可维护性和可扩展性

**Q2: Service 层应包含什么？**
- ✅ 业务逻辑、参数验证、业务规则检查、事务协调、事件发布
- ❌ 不包含数据库操作和 HTTP 处理

**Q3: 如何处理跨 Service 事务？**
- Repository 层提供事务支持
- 使用事件驱动架构
- 使用 Saga 模式
- 考虑最终一致性

**Q4: 如何进行依赖注入？**
- 使用 ServiceContainer 管理服务依赖
- 通过构造函数注入 Repository 接口
- 通过接口注入，便于测试

---

## 参考文档

**架构设计**：
- `docs/architecture/架构设计规范.md` - 完整架构设计说明
- `docs/architecture/repository层设计规范.md` - Repository 层详细设计
- `docs/architecture/路由层设计规范.md` - Router 层设计规范
- `docs/architecture/项目开发规则.md` - 项目开发规则

**工程规范**：
- `docs/engineering/软件工程规范_v2.0.md` - 编码规范和最佳实践
- `docs/engineering/需求分析文档_v2.0.md` - 需求分析文档

**API 文档**：
- `docs/api/API设计规范.md` - API 设计规范
- `docs/api/` - 各业务模块 API 文档

**测试文档**：
- `docs/testing/` - 测试规范和测试报告
- `test/README.md` - 测试运行指南

**部署运维**：
- `docs/ops/` - 部署和运维文档
- `docker/README.md` - Docker 环境说明

---

## 版本历史

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| 2.1 | 2025-10-22 | 简化文档，详细内容指向 docs/ 目录 |
| 2.0 | 2025-10-06 | 重构版本，Repository 模式、依赖注入、事件驱动 |
| 1.0 | 2025-04-25 | 初始版本，基础分层架构 |

---

**最后更新**：2025-10-22
**维护者**：青羽后端架构团队
