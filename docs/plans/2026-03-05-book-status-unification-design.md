# BookStatus 枚举统一设计方案

**设计日期**: 2026-03-05
**设计者**: Kore
**优先级**: 🔴 P0
**问题来源**: Issue #011-A: 前后端数据类型不一致

---

## 问题描述

### 当前状态

后端存在**三套不同的BookStatus定义**，违反了DRY原则，可能导致数据不一致：

| 位置                                     | 定义                            | 状态值                                              |
| -------------------------------------- | ----------------------------- | ------------------------------------------------ |
| `models/bookstore/book.go:15-20`       | `models/bookstore.BookStatus` | draft, ongoing, completed, paused                |
| `models/shared/types/enums.go:210-218` | `types.BookStatus`            | draft, **published**, completed, paused, deleted |
| `internal/domain/book.go:76-91`        | `domain.BookStatus`           | draft, ongoing, completed, paused, deleted       |

### 关键冲突

1. **`published` vs `ongoing` 冲突**
   - `enums.go` 使用 `published`（已发布）
   - `book.go` 和 `domain.go` 使用 `ongoing`（连载中）

2. **`deleted` 状态缺失**
   - `book.go` 缺少 `deleted` 状态

3. **前端使用 `serializing`**
   - 前端手动定义使用 `'serializing'`（与后端 `ongoing` 不匹配）

---

## 统一方案

### 选择标准定义

选择 `internal/domain/book.go` 作为**唯一定义源**，原因：

1. ✅ **状态最完整** - 包含所有5个状态（draft, ongoing, completed, paused, deleted）
2. ✅ **语义更清晰** - `ongoing` 比 `published` 更准确表达"连载中"状态
3. ✅ **位于domain层** - 符合DDD设计原则，业务概念应该在domain层
4. ✅ **已有验证方法** - `IsValid()` 方法已包含完整状态验证

### 统一后的定义

```go
// internal/domain/book.go
package domain

// BookStatus 书籍状态枚举
type BookStatus string

const (
    // BookStatusDraft 草稿
    BookStatusDraft BookStatus = "draft"

    // BookStatusOngoing 连载中 (已发布且正在更新)
    BookStatusOngoing BookStatus = "ongoing"

    // BookStatusCompleted 已完结
    BookStatusCompleted BookStatus = "completed"

    // BookStatusPaused 已暂停更新
    BookStatusPaused BookStatus = "paused"

    // BookStatusDeleted 已删除
    BookStatusDeleted BookStatus = "deleted"
)
```

---

## 迁移方案

### Phase 1: 清理重复定义（删除冲突定义）

#### 1.1 删除 `models/bookstore/book.go` 中的定义

**当前** (第12-20行):
```go
// BookStatus 书籍状态枚举
type BookStatus string

const (
    BookStatusDraft     BookStatus = "draft"
    BookStatusOngoing   BookStatus = "ongoing"
    BookStatusCompleted BookStatus = "completed"
    BookStatusPaused    BookStatus = "paused"
)
```

**修改后**:
```go
// 删除此定义，改为导入domain包
import "Qingyu_backend/internal/domain"

// Book 使用 domain.BookStatus
type Book struct {
    // ...
    Status domain.BookStatus `bson:"status" json:"status" validate:"required"`
    // ...
}
```

#### 1.2 删除 `models/shared/types/enums.go` 中的定义

**当前** (第209-218行):
```go
// BookStatus 书籍状态
type BookStatus string

const (
    BookStatusDraft     BookStatus = "draft"
    BookStatusPublished BookStatus = "published"  // ← 冲突值
    BookStatusCompleted BookStatus = "completed"
    BookStatusPaused    BookStatus = "paused"
    BookStatusDeleted   BookStatus = "deleted"
)
```

**修改后**:
```go
// 删除BookStatus定义，保留其他类型定义
// BookStatus 已迁移至 internal/domain/book.go
```

**注意**: 需要检查并更新所有导入 `types.BookStatus` 的引用

---

### Phase 2: 更新所有引用

#### 2.1 查找所有使用BookStatus的文件

```bash
# 搜索使用 models/bookstore.BookStatus 的文件
grep -r "bookstore\.BookStatus\|bookstore\.Book" --include="*.go" .

# 搜索使用 types.BookStatus 的文件
grep -r "types\.BookStatus\|types\.Book" --include="*.go" .
```

#### 2.2 更新导入语句

**需要更新的包**（根据审查）:
1. `service/bookstore/` - 更新为使用 `domain.BookStatus`
2. `repository/mongodb/bookstore/` - 更新为使用 `domain.BookStatus`
3. `handler/bookstore/` - 更新为使用 `domain.BookStatus`
4. `dto/bookstore/` - 更新为使用 `domain.BookStatus`

**更新模式**:
```go
// 旧的导入
import "Qingyu_backend/models/bookstore"
// 使用: bookstore.BookStatus

// 新的导入
import "Qingyu_backend/internal/domain"
// 使用: domain.BookStatus
```

---

### Phase 3: 前端同步更新

#### 3.1 更新前端BookStatus枚举

**当前** (`Qingyu_fronted/src/types/bookstore.ts`):
```typescript
export type BookStatus = 'serializing' | 'completed' | 'paused'
```

**修改后**:
```typescript
export type BookStatus = 'draft' | 'ongoing' | 'completed' | 'paused' | 'deleted'
```

#### 3.2 更新使用BookStatus的组件

需要搜索并更新所有使用BookStatus的地方：
```bash
# 在前端目录搜索
grep -r "BookStatus\|serializing\|'completed'" --include="*.ts" --include="*.vue" src/
```

---

## 数据迁移

### 现有数据处理

#### 场景1: 数据库中存在 `published` 状态

如果数据库中有书籍使用 `published` 状态：

```sql
-- 检查是否有published状态
db.books.find({status: "published"}).count()

-- 如果存在，需要迁移
db.books.updateMany(
  {status: "published"},
  {$set: {status: "ongoing"}}
)
```

**迁移脚本** (`cmd/migrate/migrate_book_status.go`):
```go
func MigrateBookStatus(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("books")

    // 迁移 published → ongoing
    result, err := collection.UpdateMany(
        ctx,
        bson.M{"status": "published"},
        bson.M{"$set": bson.M{"status": "ongoing"}},
    )

    fmt.Printf("Migrated %d books from 'published' to 'ongoing'\n", result.ModifiedCount)
    return err
}
```

---

## 验证清单

### 后端验证
- [ ] `models/bookstore/book.go` 删除重复定义
- [ ] `models/shared/types/enums.go` 删除重复定义
- [ ] 所有Service更新为使用 `domain.BookStatus`
- [ ] 所有Repository更新为使用 `domain.BookStatus`
- [ ] 所有Handler更新为使用 `domain.BookStatus`
- [ ] 所有DTO更新为使用 `domain.BookStatus`

### 前端验证
- [ ] `src/types/bookstore.ts` 更新BookStatus类型
- [ ] 所有组件更新使用新的枚举值
- [ ] API调用正确传递新状态值

### 数据验证
- [ ] 检查数据库中是否存在 `published` 状态
- [ ] 如需要，执行数据迁移脚本
- [ ] 验证迁移后的数据完整性

### 测试验证
- [ ] 单元测试更新
- [ ] 集成测试验证
- [ ] API测试验证前后端状态同步

---

## 实施计划

### Step 1: 准备（30分钟）
- [ ] 备份当前代码
- [ ] 创建新分支 `feature/book-status-unification`
- [ ] 检查数据库中是否有 `published` 状态数据

### Step 2: 后端迁移（2-3小时）
- [ ] 删除 `models/bookstore/book.go` 中的BookStatus定义
- [ ] 删除 `models/shared/types/enums.go` 中的BookStatus定义
- [ ] 更新所有导入语句
- [ ] 运行测试验证

### Step 3: 前端迁移（1-2小时）
- [ ] 更新 `src/types/bookstore.ts`
- [ ] 搜索并更新所有使用BookStatus的组件
- [ ] 运行前端测试

### Step 4: 数据迁移（如需要，30分钟）
- [ ] 执行数据迁移脚本
- [ ] 验证数据完整性

### Step 5: 集成测试（1小时）
- [ ] 前后端联调测试
- [ ] 验证状态流转正确
- [ ] 提交代码

---

## 回滚方案

如果迁移后出现问题：

1. **代码回滚**: `git revert <commit>`
2. **数据回滚**: 如已执行数据迁移，准备回滚脚本
3. **分支保护**: 使用feature分支，不影响主分支

---

## 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 数据库中存在published状态 | 中 | 高 | 先检查数据，准备迁移脚本 |
| 其他代码依赖旧定义 | 中 | 中 | 全面搜索引用，使用IDE重构工具 |
| 前端组件兼容性 | 低 | 中 | 前端使用TypeScript，编译期检查 |
| API调用失败 | 低 | 高 | 充分的测试验证 |

---

## 相关文档

- [Issue #011: 前后端数据类型不一致](../issues/011-frontend-backend-data-type-inconsistency.md)
- [P0问题审查报告](../reports/2026-03-05-p0-issue-audit-report.md)
- [数据模型V2设计](./2026-03-04-editor-data-model-v2-schema.md)

---

**设计完成时间**: 2026-03-05
**预计实施时间**: 4-6小时
**建议执行者**: 后端 + 前端协同
