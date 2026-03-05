# Issue #004: 代码质量改进

**优先级**: 中 (P1)
**类型**: 代码质量
**状态**: 待处理
**创建日期**: 2026-03-05
**相关报告**: [Writer DTO 重构总结报告](../reports/2026-03-05-dto-refactoring-summary.md#24-代码质量问题)

---

## 问题描述

项目中存在多处代码质量问题，影响代码的可维护性、一致性和健壮性。

---

## 问题清单

### 4.1 错误消息语言不统一

**优先级**: 中

**问题**: 代码中错误消息混用中文和英文，不利于国际化支持。

```go
// 有些地方用中文
return errors.New("模板不存在")
return fmt.Errorf("无效的书单ID: %w", err)

// 有些地方用英文
return errors.New("template not found")
return fmt.Errorf("invalid book list ID: %w", err)
```

**影响**:
- 日志分析困难
- 国际化支持复杂
- 错误处理不一致

**解决方案**:

1. **建立错误码体系**

```go
// pkg/errors/codes.go
package errors

const (
    // 通用错误码 (1000-1999)
    ErrCodeNotFound       = 1001
    ErrCodeInvalidInput   = 1002
    ErrCodeInternalError  = 1003

    // Auth 错误码 (2000-2999)
    ErrCodeUserNotFound      = 2001
    ErrCodeInvalidPassword   = 2002

    // Writer 错误码 (3000-3999)
    ErrCodeProjectNotFound   = 3001
    ErrCodeDocumentNotFound  = 3002
)

// RepositoryError 结构
type RepositoryError struct {
    Code    int
    Message string // 英文消息（默认）
    CN      string // 中文消息
    Err     error
}
```

2. **统一使用英文错误消息**

```go
// 修改为
return errors.NewRepositoryError(
    errors.ErrCodeNotFound,
    "template not found",
    err,
)

// 错误处理层负责转换为用户语言
func HandleError(err error) string {
    if repoErr, ok := err.(*errors.RepositoryError); ok {
        if userLang == "zh" {
            return repoErr.CN
        }
        return repoErr.Message
    }
    return err.Error()
}
```

### 4.2 MongoDB 排序参数类型错误

**优先级**: 中

**问题**: 使用无序的 `bson.M` 作为排序参数，可能导致排序不稳定。

```go
// ❌ 错误：使用无序的 bson.M
opts := options.Find().SetSort(bson.M{"category": 1, "created_at": -1})
// MongoDB 可能返回错误的排序顺序
```

**解决方案**:

```go
// ✅ 正确：使用有序的 bson.D
opts := options.Find().SetSort(bson.D{
    {Key: "category", Value: 1},
    {Key: "created_at", Value: -1},
})
```

**需要修复的位置**:

```bash
# 搜索所有使用 bson.M 排序的地方
grep -r "SetSort.*bson.M" repository/
```

已知的需要修复：
- [ ] `repository/mongodb/auth/permission_template_repository_mongo.go:191`
- [ ] 其他使用 `SetSort(bson.M{...})` 的地方

### 4.3 测试中的 Nil 检查缺失

**优先级**: 高

**问题**: 测试未检查返回值是否为 nil 就直接访问字段。

```go
// ❌ 危险
found, err := repo.GetByID(ctx, id)
require.NoError(t, err)
assert.Equal(t, 0, found.LikeCount)  // panic if found is nil
```

**解决方案**:

```go
// ✅ 安全
found, err := repo.GetByID(ctx, id)
require.NoError(t, err)
assert.NotNil(t, found)  // 先检查
assert.Equal(t, 0, found.LikeCount)
```

详见 [Issue #003: 测试基础设施改进](./003-test-infrastructure-improvements.md)

### 4.4 Magic Number 和 Magic String

**优先级**: 低

**问题**: 代码中存在魔法数字和魔法字符串。

```go
// ❌ Magic Number
limit := 10  // 这里的 10 是什么含义？
if len(title) > 200 {  // 200 的含义？
    return errors.New("title too long")
}

// ❌ Magic String
if user.Role == "admin" {  // 字符串硬编码
    // ...
}
```

**解决方案**:

```go
// ✅ 使用常量
const (
    MaxTitleLength      = 200
    DefaultPageSize     = 10
    MaxPageSize         = 100
    RoleAdmin          = "admin"
    RoleAuthor         = "author"
    RoleReader         = "reader"
)

limit := DefaultPageSize
if len(title) > MaxTitleLength {
    return errors.New("title too long")
}

if user.Role == RoleAdmin {
    // ...
}
```

### 4.5 过时的注释和 TODO

**优先级**: 低

**问题**: 代码中存在过时的注释和未完成的 TODO。

```bash
# 搜索所有 TODO
grep -r "TODO" --include="*.go" . | wc -l

# 搜索 FIXME
grep -r "FIXME" --include="*.go" . | wc -l
```

**建议**:
1. 定期审查 TODO 和 FIXME
2. 为 TODO 添加负责人和预期时间
3. 创建 Issue 跟踪重要的 TODO

```go
// ❌ 不好的 TODO
// TODO: 优化这个查询

// ✅ 好的 TODO
// TODO(@yukin, 2026-03-31): 优化这个查询
// 当前复杂度 O(n²)，预期优化到 O(n log n)
// 参考: https://github.com/xxx/issues/123
```

---

## 改进计划

### Phase 1: 紧急修复（1-2 周）

1. **修复 MongoDB 排序参数**
   - 扫描所有 `SetSort(bson.M{...})`
   - 替换为 `SetSort(bson.D{...})`
   - 添加测试验证排序正确性

2. **添加关键测试的 Nil 检查**
   - 优先修复 GetByID 相关测试
   - 修复 Create/Update 相关测试

### Phase 2: 中期改进（1-2 月）

1. **建立错误码体系**
   - 设计错误码规范
   - 创建 `pkg/errors` 包
   - 逐步迁移现有错误处理

2. **消除 Magic Number**
   - 定义常用常量
   - 替换硬编码值

### Phase 3: 长期维护（持续）

1. **代码审查清单**
   - 添加错误消息语言检查
   - 添加排序参数类型检查
   - 添加 Nil 检查要求

2. **定期清理**
   - 每月审查 TODO/FIXME
   - 更新过时的注释
   - 清理未使用的代码

---

## 代码审查清单

在提交代码前，请检查：

### 错误处理
- [ ] 错误消息使用英文
- [ ] 错误码正确设置
- [ ] 错误被正确包装（`fmt.Errorf` with `%w`）

### MongoDB 查询
- [ ] 排序使用 `bson.D` 而非 `bson.M`
- [ ] ID 转换正确处理
- [ ] 查询参数经过验证

### 测试代码
- [ ] 可能返回 nil 的值有 Nil 检查
- [ ] 测试独立运行通过
- [ ] 测试并发运行通过

### 代码风格
- [ ] 无 Magic Number（使用常量）
- [ ] 无 Magic String（使用常量）
- [ ] 注释准确且不过时

---

## 自动化检查

可以使用以下工具自动检查部分问题：

```bash
# 使用 golangci-lint
golangci-lint run

# 搜索潜在问题
grep -r "SetSort.*bson.M" --include="*.go" .
grep -r "require.NoError" --include="*_test.go" -A1 | grep "assert\|found\." | grep -v "assert.NotNil"
```

---

## 相关 Issue

- [#001: 统一模型层 ID 字段类型](./001-unify-id-type-in-models.md)
- [#002: Repository Create 方法未回设 ID](./002-create-method-id-not-set-bug.md)
- [#003: 测试基础设施改进](./003-test-infrastructure-improvements.md)
