# 依赖规则文档

**文档版本**: v1.0
**创建日期**: 2026-02-09
**适用范围**: Qingyu_backend架构重构阶段2

---

## 概述

本文档定义了项目模块间的依赖规则，目的是实现清晰的模块边界，防止跨层依赖和循环依赖。

## 依赖方向原则

### 允许的依赖方向

```
┌─────────────────────────────────────────────────────┐
│                   API Layer                         │
│                  (api/v1/*)                          │
└────────────────────┬────────────────────────────────┘
                     │ 可以依赖
                     ▼
┌─────────────────────────────────────────────────────┐
│                 Service Layer                       │
│                 (service/*)                          │
└────────────────────┬────────────────────────────────┘
                     │ 可以依赖
                     ▼
┌─────────────────────────────────────────────────────┐
│              Repository Layer                        │
│              (repository/*)                          │
└────────────────────┬────────────────────────────────┘
                     │ 可以依赖
                     ▼
┌─────────────────────────────────────────────────────┐
│                 Models Layer                        │
│                 (models/*)                           │
└─────────────────────────────────────────────────────┘
```

### 禁止的依赖方向

1. **下层不能依赖上层**
   - ❌ models → repository
   - ❌ repository → service
   - ❌ service → api

2. **禁止跨层直接依赖**
   - ❌ api → repository（必须通过service）
   - ❌ api → models（某些场景允许，但要谨慎）

3. **禁止横向依赖**
   - ❌ service/user → service/writer（应该通过接口）
   - ❌ 两个业务模块直接相互依赖

---

## Shared模块依赖规则

### 允许的依赖

以下模块**可以**直接依赖`service/shared/*`：

```
✅ service/container/    (容器初始化)
✅ service/interfaces/shared/  (适配器层)
✅ router/shared/       (共享路由)
```

### 禁止的依赖

以下模块**不应该**直接依赖`service/shared/*`的实现细节：

```
❌ service/user/        (业务服务)
❌ service/writer/       (业务服务)
❌ service/reader/       (业务服务)
❌ api/v1/user/         (API层)
❌ api/v1/writer/        (API层)
```

### 正确做法：通过Port接口

业务模块应该依赖`service/interfaces/shared`中定义的Port接口，而不是直接依赖shared实现：

```go
// ❌ 错误：直接依赖shared实现
import "Qingyu_backend/service/shared/auth"

type UserService struct {
    authService *auth.AuthService  // 直接依赖具体实现
}

// ✅ 正确：依赖Port接口
import "Qingyu_backend/service/interfaces/shared"

type UserService struct {
    authPort shared.AuthPort  // 依赖接口
}
```

---

## 允许的Import清单

### API层可以导入

```
api/v1/* 可以导入:
  - github.com/gin-gonic/gin
  - Qingyu_backend/service/*
  - Qingyu_backend/service/interfaces/*
  - Qingyu_backend/models/*
  - Qingyu_backend/pkg/*
  - Qingyu_backend/middleware/*
```

### Service层可以导入

```
service/* 可以导入:
  - context
  - Qingyu_backend/models/*
  - Qingyu_backend/repository/*
  - Qingyu_backend/service/interfaces/*
  - Qingyu_backend/pkg/errors
  - Qingyu_backend/pkg/quota
```

### Service层不应该导入

```
service/{user,writer,reader,ai} 不应该导入:
  ❌ Qingyu_backend/service/shared/* (直接导入)
  ❌ Qingyu_backend/api/*
```

**例外**：以下情况允许导入shared
- `service/container/` (容器初始化)
- 适配器实现

---

## 接口层依赖规则

### service/interfaces/ 职责

`service/interfaces/`定义的是服务契约，应该：

```
✅ 定义接口（Port）
✅ 定义请求/响应结构
✅ 被其他模块依赖
```

```
❌ 包含具体实现
❌ 依赖service/shared/*
❌ 依赖业务service/*
```

---

## 特定模块规则

### Auth模块

**允许依赖auth的模块**：
- `service/container/` (初始化)
- `service/interfaces/shared/` (适配器)
- `api/v1/auth/` (认证API)
- `api/v1/shared/auth_api.go` (共享认证API)
- `middleware/auth/` (认证中间件)

**禁止直接依赖auth的模块**：
- `service/user/` (应该通过AuthPort)
- `service/writer/` (应该通过AuthPort)

### Storage模块

**允许依赖storage的模块**：
- `service/container/`
- `service/interfaces/shared/`
- `api/v1/shared/storage_api.go`

**禁止直接依赖storage的模块**：
- 业务服务 (应该通过StoragePort)

### Cache模块

**允许依赖cache的模块**：
- `service/container/`
- `service/interfaces/shared/`

**禁止直接依赖cache的模块**：
- 业务服务 (应该通过CachePort)

---

## 循环依赖检测

### 禁止的循环依赖模式

```
❌ service/A → service/B → service/A
❌ service/A → service/interfaces/A → service/A
```

### 检测方法

使用`go mod graph`或专用工具检测：

```bash
# 查看依赖图
go mod graph

# 使用循环依赖检测工具
go run ./scripts/check-dependencies.sh
```

---

## 违规处理

### CI检查失败

如果CI检查发现违规依赖：

1. **查看报告**: 检查`check-dependencies.log`
2. **修复依赖**:
   - 移除违规的import
   - 通过接口注入替代直接依赖
3. **验证修复**: 本地运行检查脚本
4. **提交代码**

### 临时豁免

如果确实需要临时豁免某些依赖：

1. 在代码中添加注释说明原因
2. 创建GitHub Issue跟踪
3. 设置明确的解决期限
4. 定期审查豁免列表

```go
// TODO(ARCH-001): 临时依赖，计划在v1.2移除
// 原因: 需要等待AuthPort完整实现
// 跟踪: https://github.com/yukin371/Qingyu_backend/issues/xxx
import "Qingyu_backend/service/shared/auth"
```

---

## 迁移路径

### 阶段1: 建立规则（当前）

- ✅ 定义依赖规则
- ✅ 创建检查脚本
- ⏳ 集成到CI

### 阶段2: 修复违规

- [ ] 识别所有违规依赖
- [ ] 创建Port接口
- [ ] 实现适配器
- [ ] 逐步迁移业务代码

### 阶段3: 强制执行

- [ ] CI检查失败阻止合并
- [ ] 移除临时豁免
- [ ] 定期审计依赖关系

---

## 附录：快速参考

### 常见违规模式

| 违规类型 | 示例 | 正确做法 |
|---------|------|---------|
| API直接依赖shared | `api/v1/user` import `service/shared/auth` | 使用Port接口 |
| 服务直接依赖shared | `service/user` import `service/shared/storage` | 使用StoragePort |
| 业务服务相互依赖 | `service/user` import `service/writer` | 通过接口或事件 |
| 接口层依赖实现 | `service/interfaces` import `service/shared` | 只定义接口 |

### 检查命令

```bash
# 本地检查
./scripts/check-dependencies.sh

# CI检查
go test ./scripts/... -v
```

---

**维护**: 本文档随架构演进持续更新
**反馈**: 如有问题请在GitHub Issue中提出
