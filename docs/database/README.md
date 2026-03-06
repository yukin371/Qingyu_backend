# 数据库迁移工具使用文档

> **创建日期**: 2025-10-12  
> **版本**: v1.0  
> **状态**: ✅ 已完成

---

## 📋 概述

本迁移工具用于管理MongoDB数据库的schema变更和数据迁移，提供版本控制、正向迁移、回滚和种子数据功能。

### 主要特性

- ✅ **版本管理**: 追踪所有已应用的迁移
- ✅ **正向迁移**: 升级数据库schema
- ✅ **回滚功能**: 回退到之前的版本
- ✅ **种子数据**: 快速填充测试数据
- ✅ **状态查询**: 查看当前迁移状态
- ✅ **历史记录**: 记录迁移执行时间

---

## 🚀 快速开始

### 1. 构建迁移工具

```bash
cd Qingyu_backend
go build -o migrate cmd/migrate/main.go
```

### 2. 运行迁移

```bash
# 升级到最新版本
./migrate -command=up

# 查看迁移状态
./migrate -command=status

# 回滚一步
./migrate -command=down -steps=1

# 运行种子数据
./migrate -command=seed
```

### 3. 校验索引规范

```bash
go run ./cmd/verify_indexes
```

如需指定规范文件：

```bash
go run ./cmd/verify_indexes -spec docs/database/indexes.yaml
```

### 4. 慢查询监控

应用启动后会自动接入 MongoDB `CommandMonitor`，并按以下配置输出慢查询告警日志：

```yaml
database:
  primary:
    mongodb:
      profiling_level: 1  # 0=关闭, 1=仅慢查询, 2=记录所有命令
      slow_ms: 100        # 慢查询阈值（毫秒）
```

如需手动校准 MongoDB profiler，可使用：

```bash
go run ./cmd/mongodb-profiler
```

---

## 📚 命令详解

### up - 升级迁移

执行所有未应用的迁移，按版本号顺序执行。

```bash
./migrate -command=up
```

**输出示例**:
```
Applying migration 001: Add indexes to users collection
  ✓ Created username unique index
  ✓ Created email unique index
  ✓ Created phone unique index
  ✓ Created created_at index
✓ Migration 001 applied successfully

Applying migration 002: Add view_count and like_count fields to books
  ✓ Updated 5 books with new fields
  ✓ Created view_count and like_count index
✓ Migration 002 applied successfully

✓ Command completed successfully
```

---

### down - 回滚迁移

回滚最近的迁移。

```bash
# 回滚最近的1个迁移
./migrate -command=down -steps=1

# 回滚最近的2个迁移
./migrate -command=down -steps=2

# 回滚所有迁移
./migrate -command=down -steps=0
```

**输出示例**:
```
Rolling back migration 002: Add view_count and like_count fields to books
  ✓ Removed fields from 5 books
  ✓ Dropped view_count and like_count index
✓ Migration 002 rolled back successfully

✓ Command completed successfully
```

---

### status - 查看状态

查看所有迁移的当前状态。

```bash
./migrate -command=status
```

**输出示例**:
```
=== Migration Status ===

VERSION              STATUS     DESCRIPTION                                       
--------------------------------------------------------------------------------
001                  Applied    Add indexes to users collection                   
002                  Pending    Add view_count and like_count fields to books    

Total: 2 migrations, 1 applied, 1 pending

✓ Command completed successfully
```

---

### seed - 运行种子数据

填充测试数据到数据库。

```bash
./migrate -command=seed
```

**输出示例**:
```
=== Running Seeds ===

✓ Seeded 4 users
  Test accounts:
    - admin:admin@qingyu.com (password: password123)
    - author1:author1@qingyu.com (password: password123)
    - reader1:reader1@qingyu.com (password: password123)
    - reader2:reader2@qingyu.com (password: password123)

✓ Seeded 8 categories
  - 玄幻 (ID: ...)
  - 都市 (ID: ...)
  - 仙侠 (ID: ...)
  ...

✓ Seeded 5 books
  - 修真世界 by 方想 (ID: ...)
  - 诡秘之主 by 爱潜水的乌贼 (ID: ...)
  ...

✓ All seeds completed

✓ Command completed successfully
```

---

### reset - 重置所有迁移

⚠️ **危险操作！** 回滚所有迁移并删除迁移记录。

```bash
./migrate -command=reset
```

系统会要求确认：
```
⚠️  WARNING: This will rollback all migrations!
Are you sure? (yes/no): yes

Rolling back migration 002: Add view_count and like_count fields to books
✓ Migration 002 rolled back successfully

Rolling back migration 001: Add indexes to users collection
✓ Migration 001 rolled back successfully

✓ All migrations reset successfully

✓ Command completed successfully
```

---

## 📝 编写迁移脚本

### 迁移接口

每个迁移必须实现以下接口：

```go
type Migration interface {
	Version() string
	Description() string
	Up(ctx context.Context, db *mongo.Database) error
	Down(ctx context.Context, db *mongo.Database) error
}
```

### 迁移示例1：添加索引

**文件**: `migration/examples/001_add_user_indexes.go`

```go
package examples

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddUserIndexes struct{}

func (m *AddUserIndexes) Version() string {
	return "001"
}

func (m *AddUserIndexes) Description() string {
	return "Add indexes to users collection"
}

func (m *AddUserIndexes) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 创建唯一索引
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, usernameIndex)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	fmt.Println("  ✓ Created username unique index")
	return nil
}

func (m *AddUserIndexes) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	_, err := collection.Indexes().DropOne(ctx, "username_1")
	if err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}

	fmt.Println("  ✓ Dropped username index")
	return nil
}
```

### 迁移示例2：添加字段

**文件**: `migration/examples/002_add_book_fields.go`

```go
package examples

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddBookFields struct{}

func (m *AddBookFields) Version() string {
	return "002"
}

func (m *AddBookFields) Description() string {
	return "Add view_count and like_count fields to books"
}

func (m *AddBookFields) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("books")

	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"view_count": 0,
			"like_count": 0,
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add fields: %w", err)
	}

	fmt.Printf("  ✓ Updated %d books\n", result.ModifiedCount)
	return nil
}

func (m *AddBookFields) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("books")

	filter := bson.M{}
	update := bson.M{
		"$unset": bson.M{
			"view_count": "",
			"like_count": "",
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove fields: %w", err)
	}

	fmt.Printf("  ✓ Removed fields from %d books\n", result.ModifiedCount)
	return nil
}
```

### 注册迁移

在 `cmd/migrate/main.go` 中注册新迁移：

```go
func registerMigrations(manager *migration.Manager) {
	manager.RegisterMultiple(
		&examples.AddUserIndexes{},
		&examples.AddBookFields{},
		// 添加你的新迁移
		&examples.YourNewMigration{},
	)
}
```

---

## 🌱 种子数据

### 编写种子数据

**文件**: `migration/seeds/users.go`

```go
package seeds

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func SeedUsers(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 检查是否已有数据
	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("Users already exist, skipping seed\n")
		return nil
	}

	// 插入测试数据
	users := []interface{}{
		// ... 用户数据
	}

	_, err = collection.InsertMany(ctx, users)
	return err
}
```

### 可用种子数据

当前提供的种子数据：

1. **用户数据** (`seeds/users.go`)
   - 4个测试账号（admin、author1、reader1、reader2）
   - 密码统一为：`password123`

2. **分类数据** (`seeds/categories.go`)
   - 8个书籍分类（玄幻、都市、仙侠等）

3. **书籍数据** (`seeds/books.go`)
   - 5本示例书籍（修真世界、诡秘之主等）

---

## 🔧 配置

### 数据库配置

迁移工具使用项目的配置文件：

```yaml
# config/config.yaml
database:
  primary:
    type: mongodb
    mongodb:
      uri: mongodb://localhost:27017
      database: qingyu
      max_pool_size: 100
      min_pool_size: 10
```

### 指定配置文件

```bash
# 使用特定配置文件
./migrate -config=/path/to/config.yaml -command=up

# 使用Docker配置
./migrate -config=config/config.docker.yaml -command=up
```

---

## 📂 目录结构

```
migration/
├── manager.go              # 迁移管理器
├── examples/               # 迁移示例
│   ├── 001_add_user_indexes.go
│   └── 002_add_book_fields.go
├── seeds/                  # 种子数据
│   ├── users.go
│   ├── categories.go
│   └── books.go
└── README.md               # 本文档

cmd/
└── migrate/
    └── main.go             # 命令行工具
```

---

## 💡 最佳实践

### 1. 版本号规则

- 使用3位数字版本号：`001`, `002`, `003`...
- 按时间顺序递增
- 不要修改已应用的迁移

### 2. 迁移命名

文件命名规则：`{version}_{description}.go`

示例：
- `001_add_user_indexes.go`
- `002_add_book_fields.go`
- `003_update_user_schema.go`

### 3. 迁移原则

- **向后兼容**: 确保迁移不会破坏现有功能
- **可回滚**: 每个Up必须有对应的Down
- **幂等性**: 多次执行相同迁移应该安全
- **小步迭代**: 每次迁移只做一件事

### 4. 测试迁移

在生产环境之前，务必测试：

```bash
# 1. 执行迁移
./migrate -command=up

# 2. 验证结果
./migrate -command=status

# 3. 测试回滚
./migrate -command=down -steps=1

# 4. 重新执行
./migrate -command=up
```

### 5. 生产环境

生产环境执行迁移前：

1. **备份数据库**
2. **在staging环境测试**
3. **制定回滚计划**
4. **监控迁移过程**
5. **验证迁移结果**

---

## 🐛 故障排除

### 问题1：迁移失败

**症状**: 迁移执行失败，错误信息不清晰

**解决方案**:
1. 检查数据库连接
2. 查看错误日志
3. 验证迁移脚本语法
4. 检查数据库权限

### 问题2：索引冲突

**症状**: 创建索引时提示已存在

**解决方案**:
```go
// 先检查索引是否存在
indexes, _ := collection.Indexes().List(ctx)
// 如果存在则跳过创建
```

### 问题3：迁移历史不一致

**症状**: 迁移记录与实际状态不符

**解决方案**:
```bash
# 手动检查数据库
use qingyu
db.migrations.find()

# 如必要，手动修复记录
db.migrations.deleteOne({version: "xxx"})
```

---

## 📊 监控与日志

### 迁移记录

所有迁移记录存储在 `migrations` 集合：

```javascript
// 查询迁移历史
db.migrations.find().sort({applied_at: -1})

// 查看特定迁移
db.migrations.findOne({version: "001"})

// 查看回滚记录
db.migrations.find({rolled_back: true})
```

### 日志格式

迁移执行时的日志格式：

```
Applying migration 001: Add indexes to users collection
  ✓ Created username unique index
  ✓ Created email unique index
✓ Migration 001 applied successfully
```

---

## 🚀 示例工作流

### 开发流程

```bash
# 1. 编写新迁移
vim migration/examples/003_my_migration.go

# 2. 注册迁移
vim cmd/migrate/main.go

# 3. 测试迁移
./migrate -command=up

# 4. 验证结果
./migrate -command=status

# 5. 测试回滚
./migrate -command=down -steps=1

# 6. 重新执行
./migrate -command=up
```

### 生产部署

```bash
# 1. 备份数据库
mongodump --db qingyu --out /backup

# 2. 查看待执行迁移
./migrate -command=status

# 3. 执行迁移
./migrate -config=config/config.prod.yaml -command=up

# 4. 验证结果
./migrate -config=config/config.prod.yaml -command=status

# 5. 运行种子数据（如需要）
./migrate -config=config/config.prod.yaml -command=seed
```

---

## 📚 相关文档

- [MongoDB官方文档](https://www.mongodb.com/docs/)
- [Go MongoDB Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- [项目配置文档](../config/README.md)
- [下一步规划](../doc/implementation/01基础设施/下一步规划.md)

---

## ✅ 总结

### 核心功能

- ✅ 版本管理
- ✅ 正向迁移
- ✅ 回滚功能
- ✅ 状态查询
- ✅ 种子数据

### 使用建议

1. 始终先在测试环境验证
2. 迁移前备份数据
3. 保持迁移的小而专注
4. 确保可回滚性
5. 记录迁移历史

---

**文档版本**: v1.0  
**最后更新**: 2025年10月12日  
**维护者**: 青羽开发团队







