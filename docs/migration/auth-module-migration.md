# Auth模块迁移指南

## 概述

auth模块已从`service/shared/auth`迁移到`service/auth`，作为架构重构的一部分。本文档提供完整的迁移指南。

## 迁移背景

### 为什么要迁移？

1. **架构清晰**：将auth模块提升为独立服务，符合DDD的领域划分
2. **依赖方向**：避免业务服务直接依赖shared实现
3. **可维护性**：独立的auth模块更容易维护和测试
4. **扩展性**：为未来的微服务拆分做准备

### 迁移时间线

- **2026-02-09**: auth模块迁移完成
- **2026-02-xx**: compat 兼容层删除（`service/shared/auth` 停用）

## Import路径变更

### 旧路径（已废弃）

```go
import "Qingyu_backend/service/shared/auth"
```

### 新路径（推荐）

```go
import "Qingyu_backend/service/auth"
```

## 迁移步骤

### 1. 更新Import路径

```go
// ❌ 旧的导入
import "Qingyu_backend/service/shared/auth"

// ✅ 新的导入
import "Qingyu_backend/service/auth"
```

### 2. 更新类型引用

```go
// ❌ 旧代码
type UserService struct {
    authService *shared.AuthService
}

// ✅ 新代码
type UserService struct {
    authService *auth.Service
}
```

### 3. 更新容器注册

```go
// ❌ 旧的注册方式
container.ProvideSharedAuthService()

// ✅ 新的注册方式
container.ProvideAuthService()
```

### 4. 更新Port接口

如果使用Port接口模式：

```go
// ✅ 通过Port接口使用（推荐）
import "Qingyu_backend/service/interfaces/shared"

type UserService struct {
    authPort shared.AuthPort
}
```

## 兼容性说明

`service/shared/auth` 已删除，不再提供运行时兼容。请直接使用 `service/auth`。
密码验证器归属 `service/user`，中间件入口统一为 `internal/middleware/auth`。

### CI检查规则

- **生产代码**: 使用旧路径会触发❌错误
- **测试代码**: 使用旧路径会触发❌错误
- **新路径**: 不产生任何警告

## 验证迁移

### 1. 运行依赖检查

```bash
go run ./scripts/check-dependencies
```

预期输出：
```
🔍 检查代码依赖关系...
✅ 未发现依赖违规
```

如果发现违规：
```
❌ [1] service/user/user_service.go:15
   导入: Qingyu_backend/service/shared/auth
   规则: auth模块已迁移到service/auth，请使用新的导入路径
```

### 2. 运行测试

```bash
# 单元测试
go test ./service/auth/...

# 集成测试
go test ./test/integration/...

# 依赖检查
go test ./scripts/check-dependencies/...
```

### 3. 构建验证

```bash
go build ./...
```

## 常见问题

### Q1: 我必须立即迁移吗？

**A**: 是。`service/shared/auth` 已删除，必须迁移到新路径：
- 认证服务使用 `service/auth`
- 密码验证器使用 `service/user`
- 中间件使用 `internal/middleware/auth`

### Q2: 迁移会影响现有功能吗？

**A**: 会影响仍在使用旧导入路径的代码。迁移后行为不变，但旧路径不再可用。

### Q3: 如何确认迁移成功？

**A**: 运行以下检查：
1. 依赖检查无警告：`go run ./scripts/check-dependencies`
2. 所有测试通过：`go test ./...`
3. 构建成功：`go build ./...`

### Q4: 测试文件需要迁移吗？

**A**: 需要。测试文件中的旧路径同样会被依赖检查拦截。

## 代码示例

### 示例1: 简单导入更新

```go
// package service/user

// ❌ 迁移前（已不可用）
package user

import "Qingyu_backend/service/shared/auth"

type UserService struct {
    auth *auth.AuthService
}

// ✅ 迁移后
package user

import "Qingyu_backend/service/auth"

type UserService struct {
    auth *auth.Service
}
```

### 示例2: Port接口模式（推荐）

```go
// package service/user

// ✅ 使用Port接口
package user

import "Qingyu_backend/service/interfaces/shared"

type UserService struct {
    authPort shared.AuthPort
}

func (s *UserService) Login(username, password string) (*auth.User, error) {
    return s.authPort.Login(username, password)
}
```

### 示例3: 容器注册

```go
// package container

// ❌ 迁移前
func ProvideSharedAuthService() *auth.Service {
    return auth.NewService()
}

// ✅ 迁移后
func ProvideAuthService() *auth.Service {
    return auth.NewService()
}
```

## CI/CD集成

### GitHub Actions

依赖检查已集成到CI中：

```yaml
# .github/workflows/architecture-ci.yml
dependency-check:
  name: 依赖关系检查
  runs-on: ubuntu-latest
  steps:
    - name: 检查代码依赖关系
      run: |
        go run ./scripts/check-dependencies
```

### 本地Pre-commit Hook

```bash
# .git/hooks/pre-commit
#!/bin/bash
echo "检查依赖关系..."
go run ./scripts/check-dependencies
if [ $? -ne 0 ]; then
    echo "❌ 依赖检查失败，请修复后再提交"
    exit 1
fi
```

## 相关文档

- [依赖规则文档](../architecture/dependency-rules.md)
- [架构设计](../plan/2026-02-09-new-architecture-design.md)
- [迁移计划](../plan/2026-02-09-migration-plan.md)
- [依赖检查工具](../../scripts/check-dependencies/README.md)

## 支持

如有问题或需要帮助，请：
1. 查看上述相关文档
2. 运行依赖检查工具诊断问题
3. 联系架构团队

---

**最后更新**: 2026-02-09
**维护者**: 架构团队
**状态**: 活跃维护中
