# Repository层测试规范

## 概述

Repository层负责与数据库的交互，是数据访问的核心层。本规范定义了Repository层测试的详细要求和最佳实践。

## 核心原则

### ✅ 必须使用真实MongoDB

```go
// ✅ 正确示例
func TestAuthRepository_CreateRole(t *testing.T) {
    // Setup - 使用真实测试数据库
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    role := &authModel.Role{
        Name:        "test_role",
        Description: "Test role",
        Permissions: []string{"test.read"},
        IsSystem:    false,
    }

    // Act - 真实数据库操作
    err := repo.CreateRole(ctx, role)

    // Assert - 验证结果
    require.NoError(t, err)
    assert.NotEmpty(t, role.ID)
    assert.NotZero(t, role.CreatedAt)
}
```

### ❌ 严格禁止

```go
// ❌ 错误：Mock Repository测试Repository层
func TestUserRepository_Create(t *testing.T) {
    mockRepo := new(MockUserRepository)  // 危险！无法发现数据库问题
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
}

// ❌ 错误：使用内存数据库
func TestBookRepository_Query(t *testing.T) {
    db := setupMemoryDB()  // 无法验证MongoDB查询语法
}

// ❌ 错误：使用假数据绕过验证
func TestUserRepository_Create(t *testing.T) {
    user := &User{
        Email: "invalid-email",  // 假数据，无法发现验证问题
    }
}
```

## 测试组织结构

### 文件位置

```
repository/mongodb/{module}/{module}_repository_test.go
```

示例：
```
repository/mongodb/shared/auth_repository_test.go
repository/mongodb/bookstore/bookstore_repository_test.go
repository/mongodb/reader/reading_progress_repository_test.go
```

### 包命名

使用 `package xxx_test` 避免导入循环：

```go
package shared_test

import (
    authModel "Qingyu_backend/models/auth"
    shared "Qingyu_backend/repository/mongodb/shared"
    "Qingyu_backend/test/testutil"
)
```

## 测试数据库Setup

### 使用testutil.SetupTestDB

```go
func TestExample(t *testing.T) {
    // Setup - 获取测试数据库和清理函数
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()  // 测试结束后自动清理

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 测试逻辑...
}
```

### SetupTestDB提供的功能

1. **自动初始化** - MongoDB连接和服务容器
2. **测试数据库** - 独立的qingyu_test数据库
3. **自动清理** - 测试结束后清理所有集合
4. **清理的集合**：
   - users, books, projects, documents
   - roles, permissions, transactions
   - user_behaviors, user_profiles
   - reading_progress, annotations
   - chapters, announcements

## AAA测试模式

所有Repository测试必须遵循AAA模式（Arrange-Act-Assert）：

### Arrange（准备）

```go
// Arrange - 准备测试数据和环境
db, cleanup := testutil.SetupTestDB(t)
defer cleanup()

repo := shared.NewAuthRepository(db)
ctx := context.Background()

role := &authModel.Role{
    Name:        "test_role",
    Permissions: []string{"test.read"},
    IsSystem:    false,
}
```

### Act（执行）

```go
// Act - 执行被测试的操作
err := repo.CreateRole(ctx, role)
```

### Assert（断言）

```go
// Assert - 验证结果
require.NoError(t, err)
assert.NotEmpty(t, role.ID)
assert.NotZero(t, role.CreatedAt)
assert.False(t, role.IsSystem)

// 从数据库查询验证
found, err := repo.GetRole(ctx, role.ID)
require.NoError(t, err)
assert.Equal(t, role.Name, found.Name)
```

## 测试用例设计

### 1. CRUD操作测试

#### Create测试

```go
func TestAuthRepository_CreateRole(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    role := &authModel.Role{
        Name:        "test_role",
        Description: "Test role description",
        Permissions: []string{"test.read", "test.write"},
        IsSystem:    false,
    }

    // Act
    err := repo.CreateRole(ctx, role)

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, role.ID)
    assert.NotZero(t, role.CreatedAt)
    assert.NotZero(t, role.UpdatedAt)
}
```

#### Read测试

```go
func TestAuthRepository_GetRole(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 先创建角色
    role := &authModel.Role{
        Name:        "test_role",
        Description: "Test role",
        Permissions: []string{"test.read"},
        IsSystem:    false,
    }
    err := repo.CreateRole(ctx, role)
    require.NoError(t, err)

    // Act
    result, err := repo.GetRole(ctx, role.ID)

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, role.Name, result.Name)
    assert.Equal(t, role.Description, result.Description)
}
```

#### Update测试

```go
func TestAuthRepository_UpdateRole(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    role := &authModel.Role{
        Name:        "test_role",
        Description: "Original description",
        Permissions: []string{"test.read"},
        IsSystem:    false,
    }
    err := repo.CreateRole(ctx, role)
    require.NoError(t, err)

    // Act - 更新角色
    updates := map[string]interface{}{
        "description": "Updated description",
        "permissions": []string{"new.permission"},
    }
    err = repo.UpdateRole(ctx, role.ID, updates)

    // Assert
    require.NoError(t, err)

    // 验证更新
    result, err := repo.GetRole(ctx, role.ID)
    require.NoError(t, err)
    assert.Equal(t, "Updated description", result.Description)
}
```

#### Delete测试

```go
func TestAuthRepository_DeleteRole(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    role := &authModel.Role{
        Name:        "test_role",
        Description: "Test role",
        Permissions: []string{"test.read"},
        IsSystem:    false,
    }
    err := repo.CreateRole(ctx, role)
    require.NoError(t, err)

    // Act - 删除角色
    err = repo.DeleteRole(ctx, role.ID)

    // Assert
    require.NoError(t, err)

    // 验证已删除
    result, err := repo.GetRole(ctx, role.ID)
    assert.Error(t, err)
    assert.Nil(t, result)
}
```

### 2. 边缘案例测试

```go
func TestAuthRepository_GetRole_NotFound(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 使用不存在的ObjectID
    fakeID := primitive.NewObjectID().Hex()

    // Act
    result, err := repo.GetRole(ctx, fakeID)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "不存在")
}

func TestAuthRepository_GetRole_InvalidID(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // Act
    result, err := repo.GetRole(ctx, "invalid_id")

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "无效")
}
```

### 3. 查询操作测试

```go
func TestAuthRepository_ListRoles(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 创建多个角色
    roles := []*authModel.Role{
        {Name: "role1", Description: "Role 1", IsSystem: false},
        {Name: "role2", Description: "Role 2", IsSystem: false},
        {Name: "role3", Description: "Role 3", IsSystem: false},
    }

    for _, role := range roles {
        err := repo.CreateRole(ctx, role)
        require.NoError(t, err)
    }

    // Act
    result, err := repo.ListRoles(ctx)

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.GreaterOrEqual(t, len(result), 3)
}
```

### 4. 业务规则测试

```go
func TestAuthRepository_DeleteRole_SystemRole(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 创建系统角色
    role := &authModel.Role{
        Name:        "system_role",
        Description: "System role",
        Permissions: []string{"system.*"},
        IsSystem:    true,
    }
    err := repo.CreateRole(ctx, role)
    require.NoError(t, err)

    // Act - 尝试删除系统角色
    err = repo.DeleteRole(ctx, role.ID)

    // Assert
    assert.Error(t, err)  // 系统角色不应被删除
    assert.Contains(t, err.Error(), "不能删除系统角色")
}
```

## Table-Driven测试

对于有多个输入场景的测试，使用表格驱动模式：

```go
func TestAuthRepository_GetRole_TableDriven(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 创建测试角色
    role := &authModel.Role{
        Name:        "admin",
        Description: "Admin role",
        Permissions: []string{"admin.*"},
        IsSystem:    false,
    }
    err := repo.CreateRole(ctx, role)
    require.NoError(t, err)

    tests := []struct {
        name      string
        roleID    string
        setup     func()
        wantErr   bool
        errMsg    string
    }{
        {
            name:   "获取存在的角色",
            roleID: role.ID,
            wantErr: false,
        },
        {
            name:   "获取不存在的角色",
            roleID: primitive.NewObjectID().Hex(),
            wantErr: true,
            errMsg:  "不存在",
        },
        {
            name:   "使用无效ID",
            roleID: "invalid_id",
            wantErr: true,
            errMsg:  "无效",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            result, err := repo.GetRole(ctx, tt.roleID)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, result)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
            }
        })
    }
}
```

## 测试覆盖率目标

### 必须达到的覆盖率

| 方法类型 | 覆盖率目标 | 说明 |
|---------|-----------|------|
| Create | 100% | 数据创建必须完整测试 |
| Read (GetByID, GetByXXX) | ≥90% | 包含成功和失败场景 |
| Update | ≥80% | 包含成功、失败、边缘案例 |
| Delete | ≥80% | 包含成功、失败、业务规则验证 |
| List/Query | ≥70% | 验证查询逻辑和结果 |
| 业务方法 | 100% | 核心业务规则必须完整覆盖 |

### 覆盖率检查命令

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./repository/mongodb/shared

# 查看覆盖率
go tool cover -func=coverage.out | grep "auth_repository.go"

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html
```

## 性能测试

### Benchmark测试

```go
func BenchmarkAuthRepository_CreateRole(b *testing.B) {
    db, cleanup := testutil.SetupTestDB(&testing.T{})
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        role := &authModel.Role{
            Name:        fmt.Sprintf("bench_role_%d", i),
            Description: "Benchmark role",
            Permissions: []string{"bench.read"},
            IsSystem:    false,
        }
        repo.CreateRole(ctx, role)
    }
}

func BenchmarkAuthRepository_GetRole(b *testing.B) {
    db, cleanup := testutil.SetupTestDB(&testing.T{})
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 预创建测试数据
    role := &authModel.Role{
        Name:        "bench_role",
        Description: "Benchmark role",
        Permissions: []string{"bench.read"},
        IsSystem:    false,
    }
    repo.CreateRole(ctx, role)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.GetRole(ctx, role.ID)
    }
}
```

### 运行性能测试

```bash
# 运行Repository层的性能测试
go test -bench=. -benchmem ./repository/mongodb/shared

# 运行特定的benchmark测试
go test -bench=BenchmarkAuthRepository -benchmem ./repository/mongodb/shared
```

## 并发测试

### 测试并发安全性

```go
func TestAuthRepository_CreateRole_Concurrent(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // Act - 并发创建角色
    const concurrency = 10
    var wg sync.WaitGroup
    errors := make(chan error, concurrency)

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            role := &authModel.Role{
                Name:        fmt.Sprintf("concurrent_role_%d", index),
                Description: "Concurrent test role",
                Permissions: []string{"concurrent.read"},
                IsSystem:    false,
            }
            if err := repo.CreateRole(ctx, role); err != nil {
                errors <- err
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Assert - 验证没有错误
    for err := range errors {
        t.Errorf("并发创建角色失败: %v", err)
    }
}
```

### 竞态检测

```bash
# 使用-race标志检测竞态条件
go test -race ./repository/mongodb/shared
```

## 错误处理测试

### 数据库错误处理

```go
func TestAuthRepository_CreateRole_DuplicateKey(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    role := &authModel.Role{
        Name: "unique_role",
        Permissions: []string{"test.read"},
        IsSystem:    false,
    }
    repo.CreateRole(ctx, role)

    // Act - 尝试创建重复名称的角色
    duplicateRole := &authModel.Role{
        Name: "unique_role",  // 重复的名称
        Permissions: []string{"test.write"},
        IsSystem:    false,
    }
    err := repo.CreateRole(ctx, duplicateRole)

    // Assert - 应该返回错误
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "已存在")
}
```

### 网络错误处理

```go
func TestAuthRepository_Health(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // Act - 检查数据库健康状态
    err := repo.Health(ctx)

    // Assert - 应该能ping通
    assert.NoError(t, err)
}
```

## 数据验证测试

### 字段验证

```go
func TestAuthRepository_CreateRole_Validation(t *testing.T) {
    tests := []struct {
        name      string
        role      *authModel.Role
        wantErr   bool
        errMsg    string
    }{
        {
            name: "角色名称为空",
            role: &authModel.Role{
                Name:        "",  // 空名称
                Description: "Test",
                Permissions: []string{"test.read"},
            },
            wantErr: true,
            errMsg:  "名称不能为空",
        },
        {
            name: "权限列表为空",
            role: &authModel.Role{
                Name:        "test_role",
                Description: "Test",
                Permissions: []string{},  // 空权限
            },
            wantErr: true,
            errMsg:  "权限不能为空",
        },
        {
            name: "有效角色",
            role: &authModel.Role{
                Name:        "valid_role",
                Description: "Valid role",
                Permissions: []string{"valid.read"},
                IsSystem:    false,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            db, cleanup := testutil.SetupTestDB(t)
            defer cleanup()

            repo := shared.NewAuthRepository(db)
            ctx := context.Background()

            // Act
            err := repo.CreateRole(ctx, tt.role)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 索引和查询优化测试

### 验证索引使用

```go
func TestAuthRepository_ListRoles_Performance(t *testing.T) {
    // Arrange
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    // 创建大量测试数据
    for i := 0; i < 100; i++ {
        role := &authModel.Role{
            Name:        fmt.Sprintf("role_%d", i),
            Description: fmt.Sprintf("Role %d", i),
            Permissions: []string{"test.read"},
            IsSystem:    i%2 == 0,
        }
        repo.CreateRole(ctx, role)
    }

    // Act - 列出所有角色
    start := time.Now()
    result, err := repo.ListRoles(ctx)
    duration := time.Since(start)

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.GreaterOrEqual(t, len(result), 100)

    // 性能断言 - 100条记录应该在100ms内返回
    assert.Less(t, duration, 100*time.Millisecond, "查询性能应该优化")
}
```

## 常见问题

### Q1: 为什么Repository测试不能用Mock？

**A**: Repository层的职责就是与数据库交互，只有使用真实数据库才能：
- 验证查询语法的正确性
- 发现索引配置问题
- 验证数据类型兼容性
- 测试事务处理逻辑
- 发现性能瓶颈

### Q2: 测试数据如何避免污染？

**A**: 使用`testutil.SetupTestDB`提供的清理机制：
```go
db, cleanup := testutil.SetupTestDB(t)
defer cleanup()  // 测试结束后自动清理所有集合
```

### Q3: 如何处理依赖其他集合的测试？

**A**: 在测试中创建必要的依赖数据：
```go
func TestUserRepository_CreateUser(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    // 先创建必要的角色数据
    role := &authModel.Role{
        Name:        "user",
        Permissions: []string{"user.read"},
        IsSystem:    false,
    }
    roleRepo := shared.NewAuthRepository(db)
    roleRepo.CreateRole(context.Background(), role)

    // 然后创建用户
    user := &users.User{
        Username: "testuser",
        Roles:    []string{role.ID},
    }
    // ...
}
```

## 参考文档

- [testify使用指南](../03_测试工具指南/testify使用指南.md)
- [真实数据库测试规范](../04_真实数据测试规范/真实数据库测试规范.md)
- [禁止hack规范](../04_真实数据测试规范/禁止hack规范.md)
- [Service层测试规范](./service_层测试规范.md)
