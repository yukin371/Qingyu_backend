# 设计规范文档索引

**版本**: v2.0
**更新**: 2026-01-08
**状态**: ✅ 正式实施

---

## 📚 规范体系

```
doc/standards/
├── architecture/    # 架构设计规范
├── api/            # API接口规范
├── testing/        # 测试规范
├── engineering/    # 软件工程规范
└── coding/         # 编码规范
```

---

## 一、架构设计规范

### 1.1 核心架构文档

| 文档 | 说明 | 适用对象 |
|------|------|---------|
| [架构设计规范](architecture/架构设计规范.md) | 系统整体架构、分层设计、各层职责 | 全体开发 |
| [路由层设计规范](architecture/路由层设计规范.md) | 路由定义、中间件配置、请求分发 | 后端开发 |
| [依赖管理规范](architecture/依赖管理规范.md) | 依赖倒置、接口使用、模块解耦 | 架构师、后端开发 |

**核心内容**：
- ✅ 分层架构（Router → API → Service → Repository）
- ✅ 依赖注入和接口设计
- ✅ Repository模式
- ✅ 错误处理机制
- ✅ 数据访问规范

---

## 二、API设计规范

### 2.1 API文档

| 文档 | 说明 | 适用对象 |
|------|------|---------|
| [API设计规范](api/API设计规范.md) | RESTful设计、请求响应、认证授权 | 后端开发、前端开发 |

**核心内容**：
- ✅ RESTful URL设计
- ✅ HTTP方法和状态码使用
- ✅ 统一响应格式
- ✅ 参数规范（查询、路径、请求体）
- ✅ 错误处理和响应
- ✅ 认证授权机制
- ✅ 版本控制策略
- ✅ 性能优化和安全规范

---

## 三、测试规范

### 3.1 测试文档

| 文档 | 说明 | 适用对象 |
|------|------|---------|
| [测试架构设计规范](testing/测试架构设计规范.md) | 测试分层、单元测试、集成测试 | 全体开发 |

**核心内容**：
- ✅ 测试金字塔（60%单元 + 30%集成 + 10%E2E）
- ✅ 单元测试规范（Mock使用、AAA模式）
- ✅ 集成测试规范（TestHelper使用）
- ✅ 测试覆盖率目标（≥80%）
- ✅ 性能测试和基准测试
- ✅ CI/CD集成

---

## 四、软件工程规范

### 4.1 工程文档

| 文档 | 说明 | 适用对象 |
|------|------|---------|
| [软件工程规范](engineering/软件工程规范.md) | 编码规范、版本控制、项目管理 | 全体开发 |

**核心内容**：
- ✅ Go语言编码规范
- ✅ Git分支和Commit规范
- ✅ 代码审查流程
- ✅ 文档编写规范
- ✅ 需求和任务管理
- ✅ 质量保证工具
- ✅ 安全规范
- ✅ 部署发布流程

---

## 五、专项规范

### 5.1 模块设计文档

| 文档 | 说明 | 适用对象 |
|------|------|---------|
| [Repository层设计规范](coding/repository层设计规范.md) | 数据访问层设计、接口定义 | 后端开发 |
| [日志规范](coding/日志规范.md) | 日志级别、格式、采集 | 后端开发、运维 |
| [缓存设计规范](coding/缓存设计规范.md) | 缓存策略、失效、穿透 | 后端开发 |

---

## 六、快速参考

### 6.1 项目结构

```
Qingyu_backend/
├── api/v1/           # API层
├── service/          # Service层
│   ├── interfaces/   # 服务接口
│   └── {module}/     # 服务实现
├── repository/       # Repository层
│   ├── interfaces/   # 仓储接口
│   └── mongodb/      # MongoDB实现
├── models/           # 数据模型
├── router/           # 路由层
├── middleware/       # 中间件
└── config/           # 配置
```

### 6.2 分层调用规则

```
✅ 正确的调用：
Router → API → Service → Repository → Database

❌ 错误的调用：
API → Repository（跳过Service）
Service → Database（跳过Repository）
```

### 6.3 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 接口 | `{模块}Repository` / `{模块}Service` | `BookRepository` |
| 实现 | `{技术}{模块}Repository` | `MongoBookRepository` |
| 方法 | 动词开头 | `GetBookByID` |
| 测试 | `{filename}_test.go` | `book_service_test.go` |

### 6.4 错误处理

```go
// Repository层
return &RepositoryError{
    Code:    "NOT_FOUND",
    Message: "书籍不存在",
    Err:     err,
}

// Service层
return &ServiceError{
    Code:    "BOOK_NOT_FOUND",
    Message: "书籍不存在",
}

// API层
shared.NotFound(c, "书籍不存在")
```

### 6.5 依赖注入

```go
// Service构造函数
func NewBookstoreService(
    bookRepo repository.BookRepository,
    cache cache.Cache,
) BookstoreService {
    return &BookstoreServiceImpl{
        bookRepo: bookRepo,
        cache:    cache,
    }
}
```

---

## 七、检查清单

### 7.1 代码提交前

- [ ] 代码符合编码规范
- [ ] 通过所有单元测试
- [ ] 测试覆盖率达标
- [ ] 添加/更新了文档
- [ ] 通过静态分析检查
- [ ] 自我Review代码

### 7.2 PR提交前

- [ ] 完成功能开发
- [ ] 测试全部通过
- [ ] 更新相关文档
- [ ] 代码审查通过
- [ ] 无安全漏洞
- [ ] 性能测试通过

### 7.3 发布前

- [ ] 所有测试通过
- [ ] 性能指标达标
- [ ] 安全检查通过
- [ ] 文档完整
- [ ] 回滚方案准备
- [ ] 监控告警配置

---

## 八、最佳实践

### 8.1 设计原则

**SOLID原则**：
- **S**ingle Responsibility - 单一职责
- **O**pen/Closed - 开闭原则
- **L**iskov Substitution - 里氏替换
- **I**nterface Segregation - 接口隔离
- **D**ependency Inversion - 依赖倒置

### 8.2 开发原则

**DRY** - Don't Repeat Yourself（避免重复）
**KISS** - Keep It Simple, Stupid（保持简单）
**YAGNI** - You Aren't Gonna Need It（不要过度设计）

### 8.3 代码质量

**可读性** > **聪明才智**
**测试** > **调试**
**文档** > **口头**
**迭代** > **完美**

---

## 九、更新记录

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| v2.0 | 2026-01-08 | 精简重构，统一规范体系 | 架构团队 |
| v1.3 | 2025-10-25 | 添加测试规范 | 测试团队 |
| v1.0 | 2025-10-06 | 初始版本 | 架构团队 |

---

## 十、联系方式

**规范维护**：架构团队
**问题反馈**：提交Issue或PR
**文档更新**：遵循本规范进行更新

---

**重要提示**：
1. 所有新功能开发必须遵循相关规范
2. 规范会持续优化，定期Review
3. 特殊情况可申请豁免，需架构团队审批
4. 违反规范可能导致代码审查不通过
