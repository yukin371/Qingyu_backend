# API迁移工具套件

> **版本**: v1.0
> **更新日期**: 2026-01-29
> **用途**: 辅助API从shared包迁移到response包

## 📋 目录

1. [工具概述](#工具概述)
2. [安装使用](#安装使用)
3. [命令详解](#命令详解)
4. [使用示例](#使用示例)
5. [开发说明](#开发说明)

---

## 工具概述

本工具套件提供4个核心功能：

| 工具 | 功能 | 用途 |
|------|------|------|
| **analyze** | 分析API文件 | 统计响应调用，评估复杂度 |
| **migrate** | 迁移API文件 | 自动替换shared为response |
| **validate** | 验证迁移质量 | 检查迁移完整性 |
| **testgen** | 生成测试框架 | 辅助生成单元测试 |

### 特性

- ✅ **AST解析**: 精确解析代码，而非简单正则替换
- ✅ **安全机制**: 自动备份，dry-run模式，可回滚
- ✅ **批量处理**: 支持单文件和批量处理
- ✅ **详细报告**: JSON格式输出，便于集成

---

## 安装使用

### 编译工具

```bash
cd scripts/migration-tools
go build -o migration-tools.exe .
```

### 添加到PATH

```bash
# Windows
set PATH=%PATH%;E:\Github\Qingyu\.worktrees\block7-tdd-reader-pilot\scripts\migration-tools

# Linux/Mac
export PATH=$PATH:/path/to/migration-tools
```

---

## 命令详解

### analyze - 分析工具

分析API文件中的响应调用情况。

#### 语法

```bash
migration-tools analyze [options]
```

#### 选项

| 选项 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--path` | string | API目录或文件路径 | 必需 |
| `--output` | string | 输出文件路径（JSON） | 无（打印到控制台） |
| `--verbose` | bool | 详细输出 | false |

#### 示例

```bash
# 分析Writer模块
migration-tools analyze --path api/v1/writer

# 分析并输出到文件
migration-tools analyze --path api/v1/writer --output analysis.json

# 详细输出
migration-tools analyze --path api/v1/writer --verbose
```

#### 输出示例

```json
{
  "module": "writer",
  "total_files": 17,
  "total_calls": 395,
  "breakdown": {
    "shared_success": 65,
    "shared_error": 62,
    "shared_validation_error": 11,
    "response_success": 120,
    "response_bad_request": 45,
    ...
  },
  "complexity": "medium",
  "recommendation": "建议按复杂度从低到高逐步迁移"
}
```

---

### migrate - 迁移工具

自动迁移API文件从shared包到response包。

#### 语法

```bash
migration-tools migrate [options]
```

#### 选项

| 选项 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--file` | string | API文件路径 | 必需 |
| `--dry-run` | bool | 预览模式，不实际修改 | false |
| `--backup` | bool | 创建备份文件 | true |
| `--verbose` | bool | 详细输出 | false |

#### 迁移规则

| 原调用 | 迁移后 | 说明 |
|--------|--------|------|
| `shared.Success(c, http.StatusOK, msg, data)` | `response.Success(c, data)` | 移除状态码和消息 |
| `shared.Success(c, http.StatusCreated, msg, data)` | `response.Created(c, data)` | 改用Created |
| `shared.Error(c, http.StatusBadRequest, msg, details)` | `response.BadRequest(c, msg, details)` | 移除状态码 |
| `shared.Error(c, http.StatusUnauthorized, msg, details)` | `response.Unauthorized(c, msg)` | 移除状态码 |
| `shared.Error(c, http.StatusForbidden, msg, details)` | `response.Forbidden(c, msg)` | 移除状态码 |
| `shared.Error(c, http.StatusNotFound, msg, details)` | `response.NotFound(c, msg)` | 移除状态码 |
| `shared.Error(c, http.StatusConflict, msg, details)` | `response.Conflict(c, msg, details)` | 移除状态码 |
| `shared.Error(c, 5xx, msg, err)` | `response.InternalError(c, err)` | 移除状态码和消息 |
| `shared.ValidationError(c, err)` | `response.BadRequest(c, "参数错误", err.Error())` | 添加消息 |

#### 示例

```bash
# 预览迁移（不实际修改）
migration-tools migrate --file api/v1/writer/audit_api.go --dry-run

# 实际迁移（会创建备份）
migration-tools migrate --file api/v1/writer/audit_api.go

# 迁移但不创建备份
migration-tools migrate --file api/v1/writer/audit_api.go --backup=false

# 详细输出
migration-tools migrate --file api/v1/writer/audit_api.go --verbose
```

#### 输出示例

```
=== 迁移结果: api/v1/writer/audit_api.go ===
Success调用迁移: 8
Error调用迁移: 12
ValidationError迁移: 2
总调用迁移: 22
添加导入: [response]
移除导入: [shared net/http]
备份文件: api/v1/writer/audit_api.go.bak
耗时: 0.15秒
```

---

### validate - 验证工具

验证迁移是否完整。

#### 语法

```bash
migration-tools validate [options]
```

#### 选项

| 选项 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--path` | string | API目录或文件路径 | 必需 |
| `--checks` | string | 检查项（逗号分隔） | all |
| `--verbose` | bool | 详细输出 | false |

#### 检查项

| 检查项 | 说明 |
|--------|------|
| `imports` | 检查导入是否清理（shared, net/http） |
| `no_shared_calls` | 检查是否还有shared调用 |
| `swagger` | 检查Swagger注释是否更新 |
| `all` | 执行所有检查 |

#### 示例

```bash
# 验证所有检查项
migration-tools validate --path api/v1/writer

# 只检查shared调用
migration-tools validate --path api/v1/writer --checks=no_shared_calls

# 详细输出
migration-tools validate --path api/v1/writer --verbose
```

#### 输出示例

```
=== 验证结果: api/v1/writer ===
文件总数: 17
整体状态: ❌ 失败

检查项:
  导入清理: ❌
  无shared调用: ❌
  Swagger更新: ✅

问题列表 (42):
  [1] shared_call_found: 发现shared.Error调用，应该迁移到response包
  [2] shared_import_found: 发现shared包导入，应该移除
  ...
```

---

### testgen - 测试生成助手

生成API文件的测试框架（TODO功能）。

#### 语法

```bash
migration-tools testgen [options]
```

#### 选项

| 选项 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--file` | string | API文件路径 | 必需 |
| `--output` | string | 输出测试文件路径 | 自动生成 |
| `--verbose` | bool | 详细输出 | false |

#### 示例

```bash
# 生成测试框架
migration-tools testgen --file api/v1/writer/audit_api.go

# 指定输出路径
migration-tools testgen --file api/v1/writer/audit_api.go --output api/v1/writer/audit_api_generated_test.go
```

---

## 使用示例

### 完整工作流

```bash
# 1. 分析Writer模块
migration-tools analyze --path api/v1/writer --output writer_analysis.json

# 2. 查看分析结果
cat writer_analysis.json

# 3. 预览迁移
migration-tools migrate --file api/v1/writer/audit_api.go --dry-run

# 4. 实际迁移
migration-tools migrate --file api/v1/writer/audit_api.go

# 5. 验证迁移
migration-tools validate --path api/v1/writer/audit_api.go

# 6. 运行测试
cd api/v1/writer
go test -v -run TestAudit
```

### 批量迁移

```bash
# 批量迁移多个文件（使用shell循环）
for file in api/v1/writer/*_api.go; do
    echo "Migrating $file..."
    migration-tools migrate --file "$file"
done

# 验证所有文件
migration-tools validate --path api/v1/writer
```

---

## 开发说明

### 项目结构

```
scripts/migration-tools/
├── main.go      # CLI入口
├── analyze.go   # 分析工具实现
├── migrate.go   # 迁移工具实现
├── validate.go  # 验证工具实现
├── testgen.go   # 测试生成实现（TODO）
├── go.mod       # Go模块定义
└── README.md    # 本文档
```

### 核心概念

#### AST解析

工具使用Go的AST（抽象语法树）解析器来精确分析代码：

```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)

// 遍历AST
ast.Inspect(node, func(n ast.Node) bool {
    // 处理节点
    return true
})
```

#### 代码重写

通过修改AST节点来重写代码：

```go
// 找到函数调用
call, ok := n.(*ast.CallExpr)
if ok {
    // 修改调用
    sel := call.Fun.(*ast.SelectorExpr)
    sel.X.Name = "response"  // 修改包名
    sel.Sel.Name = "Success" // 修改函数名
}
```

### 扩展开发

#### 添加新的检查项

在`validate.go`中添加新的检查函数：

```go
func checkCustomRule(node *ast.Node, validation *FileValidation) bool {
    // 实现自定义检查逻辑
    return true
}
```

#### 添加新的迁移规则

在`migrate.go`中添加新的迁移处理：

```go
case "CustomFunc":
    r.migrateCustomFuncCall(call, sel)
```

---

## 常见问题

### Q: 工具会修改原始代码吗？

A: 默认会创建备份文件（.bak）。使用`--dry-run`可以预览变更而不实际修改。

### Q: 如何回滚迁移？

A: 使用备份文件：
```bash
cp api/v1/writer/audit_api.go.bak api/v1/writer/audit_api.go
```

### Q: 工具支持哪些Go版本？

A: Go 1.21+（使用了AST解析特性）

### Q: 可以迁移非API文件吗？

A: 工具主要设计用于API文件，但理论上可以处理任何使用shared包的Go文件。

### Q: 工具会处理注释吗？

A: 不会。Swagger注释需要手动更新或使用validate工具检查。

---

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork本仓库
2. 创建feature分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

### 代码规范

- 遵循Go标准代码风格
- 添加必要的注释和文档
- 确保所有测试通过
- 更新相关文档

---

## 许可证

MIT License

---

**版本**: v1.0
**最后更新**: 2026-01-29
**维护者**: Backend Team
