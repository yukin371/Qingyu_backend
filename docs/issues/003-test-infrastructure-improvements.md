# Issue #003: 测试基础设施改进

**优先级**: 中 (P1)
**类型**: 测试基础设施
**状态**: 部分修复
**创建日期**: 2026-03-05
**相关报告**: [Writer DTO 重构总结报告](../reports/2026-03-05-dto-refactoring-summary.md#23-测试基础设施问题)

---

## 问题描述

测试基础设施存在多个问题，导致测试不稳定、结果不可靠，影响开发效率。

### 已修复问题
- ✅ GlobalConfig 全局变量竞态条件
- ✅ 测试数据库隔离失效（`loadLocalConfigWithFallback` 未使用唯一数据库名）

### 待改进问题
- ⚠️ 测试 Nil 检查缺失
- ⚠️ 测试集合清理不完整
- ⚠️ 并发测试支持不足

---

## 子问题

### 3.1 测试数据库隔离 ✅ 已修复

**问题**: `loadLocalConfigWithFallback` 直接使用配置文件中的固定数据库名称，导致多个测试共享数据库。

**修复**:
```go
// 修复前
func loadLocalConfigWithFallback() (*config.Config, error) {
    cfg, err := config.LoadConfig(candidate)
    if err == nil {
        return cfg, nil  // ❌ 使用固定数据库名
    }
}

// 修复后
func loadLocalConfigWithFallback() (*config.Config, error) {
    cfg, err := config.LoadConfig(candidate)
    if err == nil {
        // ✅ 强制使用唯一测试数据库名
        cfg.Database.Primary.MongoDB.Database = resolveTestMongoDatabaseName(
            cfg.Database.Primary.MongoDB.Database
        )
        return cfg, nil
    }
}
```

### 3.2 GlobalConfig 竞态条件 ✅ 已修复

**问题**: 多个测试并发修改 `config.GlobalConfig` 全局变量，导致测试失败。

**修复**:
```go
// test/testutil/database.go
func SetupTestDB(t *testing.T) (*mongo.Database, func()) {
    // 保存旧的配置
    oldConfig := config.GlobalConfig  // ✅

    // 初始化全局配置
    config.GlobalConfig = cfg

    cleanup := func() {
        config.GlobalConfig = oldConfig  // ✅ 恢复
        // ... 其他清理
    }
}
```

### 3.3 测试 Nil 检查缺失 ⚠️ 待修复

**问题**: 许多测试未检查返回值是否为 nil 就直接访问字段，导致 panic。

**示例**:
```go
// 危险的模式
found, err := repo.GetByID(ctx, id)
require.NoError(t, err)
assert.Equal(t, expected, found.Field)  // ❌ found 为 nil 时 panic

// 安全的模式
found, err := repo.GetByID(ctx, id)
require.NoError(t, err)
assert.NotNil(t, found)  // ✅ 先检查 nil
assert.Equal(t, expected, found.Field)
```

**需要修复的测试**:
- [ ] `repository/mongodb/writer/*_test.go`
- [ ] `repository/mongodb/social/*_test.go`
- [ ] `repository/mongodb/reader/*_test.go`
- [ ] 其他测试文件

### 3.4 测试集合清理不完整 ⚠️ 待改进

**问题**: `SetupTestDB` 的清理函数只清理部分集合，新添加的集合可能未被清理。

**当前清理的集合**:
```go
// test/testutil/database.go
_ = db.Collection("user_behaviors").Drop(ctx)
_ = db.Collection("user_profiles").Drop(ctx)
// ... 约 30 个集合
```

**改进方案**:
1. 自动发现并清理所有非系统集合
2. 使用独立数据库，直接 Drop 数据库（当前方案）
3. 添加清理验证，确保所有测试集合已清理

### 3.5 并发测试支持不足 ⚠️ 待改进

**问题**: 部分测试不支持并发执行，必须使用 `-parallel=1`。

**原因**:
1. 共享全局状态（GlobalConfig）
2. 共享测试数据库
3. 测试之间有依赖关系

**改进方向**:
1. 消除全局状态依赖
2. 每个测试使用独立数据库
3. 确保测试独立性

---

## 改进计划

### Phase 1: 修复 Nil 检查问题（高优先级）

```bash
# 扫描需要添加 Nil 检查的测试
grep -r "require.NoError(t, err)" repository/*/*_test.go | \
  while read line; do
    file=$(echo $line | cut -d: -f1)
    line_num=$(echo $line | cut -d: -f2)
    # 检查下一行是否直接访问字段
  done
```

修复模式：
1. `GetByID` 后添加 `assert.NotNil(t, found)`
2. `FindByX` 后检查切片是否为 nil（可选）
3. `Create` 后验证 ID 已设置

### Phase 2: 完善测试清理逻辑

```go
// 自动清理所有测试集合
func cleanupAllTestCollections(db *mongo.Database) {
    ctx := context.Background()
    collections, err := db.ListCollectionNames(ctx, bson.M{"name": bson.M{
        "$not": bson.M{"$in": []string{"system.*"}},
    }})
    if err != nil {
        return
    }

    for _, coll := range collections {
        _ = db.Collection(coll).Drop(ctx)
    }
}
```

### Phase 3: 提升并发测试支持

1. **消除全局状态**
   - 使用依赖注入而非全局变量
   - 每个测试创建独立的 ServiceContainer

2. **测试独立性检查**
   - 添加随机执行顺序测试
   - 并发执行测试（`go test -parallel=N`）

3. **测试隔离**
   - 每个测试使用独立数据库（已实现）
   - 每个 goroutine 使用独立 context

---

## 测试最佳实践

### 1. Setup-Cleanup 模式

```go
func TestXxx(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    repo, ctx, cleanup := setupRepo(t)
    defer cleanup()  // ✅ 使用 defer 确保清理

    // 测试逻辑...
}
```

### 2. Nil 检查模式

```go
func TestGetByID(t *testing.T) {
    repo, ctx, cleanup := setupRepo(t)
    defer cleanup()

    // 测试找到记录
    found, err := repo.GetByID(ctx, existingID)
    require.NoError(t, err)
    assert.NotNil(t, found)  // ✅ 总是检查 nil
    assert.Equal(t, "expected", found.Field)

    // 测试找不到记录
    found, err = repo.GetByID(ctx, nonExistentID)
    require.NoError(t, err)
    assert.Nil(t, found)  // ✅ 验证为 nil
}
```

### 3. 独立性原则

```go
// ❌ 不好：测试之间有依赖
func TestA(t *testing.T) {
    // 依赖 TestB 先执行
}

func TestB(t *testing.T) {
    // 必须先执行
}

// ✅ 好：每个测试独立
func TestA(t *testing.T) {
    repo, ctx, cleanup := setupRepo(t)
    defer cleanup()
    // 完全独立的测试
}
```

---

## 验证清单

每个测试修复后：
- [ ] 单独运行测试通过
- [ ] 并发运行测试通过（`-parallel=4`）
- [ ] 随机顺序运行测试通过（`-shuffle=on`）
- [ ] 全量测试通过

---

## 相关文件

| 文件 | 状态 | 说明 |
|------|------|------|
| `test/testutil/database.go` | ✅ 已修复 | 数据库隔离和 GlobalConfig |
| `repository/mongodb/writer/*_test.go` | ⚠️ 待修复 | 需添加 Nil 检查 |
| `repository/mongodb/social/*_test.go` | ⚠️ 待修复 | 需添加 Nil 检查 |

---

## 参考链接

- [Go Testing 最佳实践](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- 相关 Issue: [#002 Create 方法未回设 ID](./002-create-method-id-not-set-bug.md)
