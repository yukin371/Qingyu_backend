# 依赖检查工具

## 概述

这个工具用于检查代码中的依赖关系是否符合项目的架构规则，防止跨层依赖和违规的直接依赖。

## 功能

- ✅ 检查业务服务是否直接依赖shared模块实现
- ✅ 识别应该使用Port接口的违规导入
- ✅ 生成详细的违规报告
- ✅ 提供修复建议

## 使用方法

### 方式1: 直接运行Go程序

```bash
cd scripts/check-dependencies
go run main.go
```

### 方式2: 编译后运行

```bash
# 编译
cd scripts/check-dependencies
go build -o check-deps

# 运行
./check-deps
```

### 方式3: 在项目根目录运行

```bash
# 从项目根目录运行
go run ./scripts/check-dependencies
```

## 输出示例

### 无违规时

```
🔍 检查代码依赖关系...

✅ 未发现依赖违规
```

### 有违规时

```
🔍 检查代码依赖关系...

❌ 发现 3 个错误, 1 个警告

❌ [1] service/user/user_service.go:15
   导入: Qingyu_backend/service/shared/auth
   规则: 不应该直接导入shared模块，请使用service/interfaces/shared中的Port接口

❌ [2] service/writer/book_service.go:20
   导入: Qingyu_backend/service/shared/storage
   规则: 不应该直接导入shared模块，请使用service/interfaces/shared中的Port接口

💡 修复建议:
   1. 使用service/interfaces/shared中定义的Port接口
   2. 通过依赖注入而非直接导入
   3. 参考文档: docs/architecture/dependency-rules.md
```

## 依赖规则

### 允许直接导入shared的模块

以下模块**可以**直接导入`service/shared/*`：

- `service/container/` - 容器初始化
- `service/interfaces/shared/` - 适配器层
- `router/shared/` - 共享路由
- `api/v1/auth/` - 认证API
- `api/v1/shared/` - 共享API
- `middleware/` - 中间件

### 禁止直接导入shared的模块

以下模块**不应该**直接导入`service/shared/*`：

- `service/user/` - 用户服务
- `service/writer/` - 写作服务
- `service/reader/` - 阅读服务
- `service/ai/` - AI服务

### 正确做法

```go
// ❌ 错误：直接依赖shared实现
import "Qingyu_backend/service/shared/auth"

type UserService struct {
    authService *auth.AuthService
}

// ✅ 正确：依赖Port接口
import "Qingyu_backend/service/interfaces/shared"

type UserService struct {
    authPort shared.AuthPort
}
```

## CI集成

### GitHub Actions

在`.github/workflows/ci.yml`中添加：

```yaml
name: CI

on: [push, pull_request]

jobs:
  check-dependencies:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Check dependencies
        run: |
          go run ./scripts/check-dependencies
```

### 本地Git Hook

在`.git/hooks/pre-commit`中添加：

```bash
#!/bin/bash
echo "检查依赖关系..."
go run ./scripts/check-dependencies
if [ $? -ne 0 ]; then
    echo "❌ 依赖检查失败，请修复后再提交"
    exit 1
fi
```

## 扩展规则

要添加新的依赖规则，修改`main.go`中的配置：

```go
var forbiddenImports = map[string]string{
    `service/user`:     `不应该直接导入shared模块`,
    `service/new-module`: `添加新模块的规则`,
}

var allowedSharedImporters = map[string]bool{
    `service/container`:  true,
    `service/new-allowed`: true,
}
```

## 测试

运行测试：

```bash
cd scripts/check-dependencies
go test -v
```

## 性能

工具会遍历项目中所有的`.go`文件（排除测试文件和vendor目录）。

在大型项目中的性能：
- 100个文件: ~100ms
- 1000个文件: ~500ms
- 10000个文件: ~2s

## 故障排除

### 问题：误报

如果发现误报，可以：

1. 检查文件路径是否正确
2. 添加到允许列表
3. 使用`//nolint:depcheck`注释临时忽略

### 问题：检查太慢

可以：

1. 只检查特定目录：修改`filepath.Walk`的过滤条件
2. 并行处理：使用goroutine并发检查
3. 缓存结果：只检查变更的文件

## 贡献

欢迎贡献改进：

1. 添加更多检查规则
2. 改进错误报告
3. 优化性能
4. 添加更多测试

## 相关文档

- [依赖规则文档](../../docs/architecture/dependency-rules.md)
- [架构设计](../../docs/plan/2026-02-09-new-architecture-design.md)
- [迁移计划](../../docs/plan/2026-02-09-migration-plan.md)
