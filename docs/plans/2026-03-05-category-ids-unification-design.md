# CategoryIDs 类型统一设计方案

**设计日期**: 2026-03-05
**设计者**: Kore
**优先级**: 🔴 P0
**问题来源**: Issue #011-B: 前后端数据类型不一致

---

## 问题描述

### 当前问题

**后端存在两套不同的 CategoryIDs 类型定义**：

| 模型 | 文件 | 类型 |
|------|------|------|
| `Book` | `models/bookstore/book.go:32` | `[]primitive.ObjectID` |
| `BookDetail` | `models/bookstore/book_detail.go:30` | `[]string` |
| `BookDetailService` | `service/bookstore/book_detail_service.go:19` | `[]string` |

### 关键冲突

1. **后端内部类型不一致**
   - `Book` 模型使用 `[]primitive.ObjectID`
   - `BookDetail` 模型使用 `[]string`
   - 违反了单一数据源原则

2. **前后端数据结构不匹配**
   - 后端：数组 `categoryIds: string[]`
   - 前端：单值 `categoryId: string`
   - 导致多分类书籍数据丢失

3. **API参数命名不一致**
   - 查询参数：`categoryId`（单数）
   - 数据字段：`categoryIds`（复数）
   - 容易引起混淆

### 问题证据

```go
// models/bookstore/book.go:32
type Book struct {
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
}

// models/bookstore/book_detail.go:30
type BookDetail struct {
    CategoryIDs  []string   `bson:"category_ids" json:"categoryIds"`
}

// 前端 src/types/bookstore.ts
interface Book {
  categoryId: string        // 单值！
  categoryName?: string
}
```

**影响**：
- 多分类书籍只能显示第一个分类
- 数据库查询条件不一致
- 类型转换代码冗余

---

## 统一方案

### 设计原则

根据 `docs/standards/id-type-unification-standard.md` 中的分层转换原则：

```
┌─────────────────────────────────────────────┐
│         API / Handler 层                     │
│         使用 string 类型                      │
└─────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────┐
│         Service 层                           │
│         使用 string 类型                      │
└─────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────┐
│         Repository 层                        │
│         使用 primitive.ObjectID              │
└─────────────────────────────────────────────┘
```

### 选择标准定义

**选择 `[]primitive.ObjectID` 作为数据库存储类型**，原因：

1. ✅ **符合ID类型统一标准** - Issue #001 的目标是将所有ID统一为 `primitive.ObjectID`
2. ✅ **性能优势** - MongoDB 索引查询更高效
3. ✅ **与现有Book模型一致** - 减少迁移成本
4. ✅ **支持多分类** - 数组类型天然支持一对多关系

### 统一后的类型定义

```go
// models/bookstore/book.go - 保持不变
type Book struct {
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
}

// models/bookstore/book_detail.go - 需要修改
type BookDetail struct {
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`  // 改为 ObjectID
}
```

---

## 迁移方案

### Phase 1: 后端类型统一（1天）

#### 1.1 修改 BookDetail 模型

**当前** (第30行):
```go
type BookDetail struct {
    CategoryIDs  []string   `bson:"category_ids" json:"categoryIds"`
}
```

**修改后**:
```go
type BookDetail struct {
    CategoryIDs  []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
}
```

#### 1.2 更新 DTO 转换逻辑

**当前 DTO** (`dto/bookstore/book_dto.go`):
```go
type BookDTO struct {
    CategoryIDs []string `json:"categoryIds"`
}
```

**保持不变** - DTO 层应该继续使用 `[]string`，在 API 层进行转换。

**转换函数** (`api/v1/bookstore/bookstore_converter.go`):
```go
// Model → DTO: ObjectID → string
func ModelIDsToDTO(ids []primitive.ObjectID) []string {
    if ids == nil {
        return nil
    }
    result := make([]string, len(ids))
    for i, id := range ids {
        result[i] = id.Hex()
    }
    return result
}

// DTO → Model: string → ObjectID
func DTOIDsToModel(ids []string) ([]primitive.ObjectID, error) {
    if ids == nil {
        return nil, nil
    }
    result := make([]primitive.ObjectID, len(ids))
    for i, idStr := range ids {
        id, err := primitive.ObjectIDFromHex(idStr)
        if err != nil {
            return nil, fmt.Errorf("invalid category ID: %s", idStr)
        }
        result[i] = id
    }
    return result, nil
}
```

#### 1.3 更新 BookDetailService

**当前** (`service/bookstore/book_detail_service.go:19`):
```go
type BookDetailService struct {
    CategoryIDs []string `json:"category_ids,omitempty"`
}
```

**修改后**:
```go
type BookDetailService struct {
    CategoryIDs []primitive.ObjectID `json:"categoryIds"`
}
```

#### 1.4 更新 Repository 查询

需要检查所有使用 CategoryIDs 的 Repository 方法：

```go
// repository/mongodb/bookstore/book_repository_mongo.go

// 按分类查询
func (r *BookRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]*Book, error) {
    oid, err := primitive.ObjectIDFromHex(categoryID)
    if err != nil {
        return nil, err
    }

    filter := bson.M{
        "category_ids": oid,  // 查询数组中包含该ID
    }

    return r.Find(ctx, filter)
}
```

---

### Phase 2: 前端同步更新（1天）

#### 2.1 更新前端类型定义

**当前** (`Qingyu_fronted/src/types/bookstore.ts`):
```typescript
export interface Book {
  categoryId: string        // 单值
  categoryName?: string
  category?: string
}
```

**修改后**:
```typescript
export interface Book {
  categoryIds: string[]     // 改为数组
  categoryNames?: string[]  // 对应的名称数组
}

// 辅助函数获取主分类（兼容旧代码）
export function getPrimaryCategory(book: Book): string | undefined {
  return book.categoryIds?.[0]
}
```

#### 2.2 更新组件使用

需要搜索并更新所有使用 `categoryId` 的组件：

```bash
# 在前端目录搜索
grep -r "categoryId" --include="*.ts" --include="*.vue" src/
```

**需要修改的典型代码**:

```vue
<!-- 旧代码 -->
<el-select v-model="book.categoryId" placeholder="选择分类">
  <el-option
    v-for="cat in categories"
    :key="cat.id"
    :label="cat.name"
    :value="cat.id"
  />
</el-select>

<!-- 新代码 -->
<el-select
  v-model="book.categoryIds"
  multiple
  placeholder="选择分类（可多选）"
>
  <el-option
    v-for="cat in categories"
    :key="cat.id"
    :label="cat.name"
    :value="cat.id"
  />
</el-select>
```

#### 2.3 更新 API 调用

**当前**:
```typescript
// 获取某分类下的书籍
async getBooksByCategory(categoryId: string) {
  return this.http.get(`/api/v1/bookstore/books?categoryId=${categoryId}`)
}
```

**修改后**:
```typescript
// 支持多分类查询
async getBooksByCategories(categoryIds: string[]) {
  const params = new URLSearchParams()
  categoryIds.forEach(id => params.append('categoryIds', id))
  return this.http.get(`/api/v1/bookstore/books?${params}`)
}
```

---

### Phase 3: API 参数统一（半天）

#### 3.1 统一 API 参数命名

**当前 API 定义** (`api/v1/bookstore/bookstore_api.go`):
```go
//	@Param		categoryId	query		string	false	"分类ID"
categoryID := c.Query("categoryId")
```

**修改后**:
```go
//	@Param		categoryIds	query		[]string	false	"分类ID列表（可多选）"
categoryIds := c.QueryArray("categoryIds")
```

#### 3.2 更新查询逻辑

```go
// api/v1/bookstore/bookstore_api.go

func (api *BookstoreAPI) GetBooks(c *gin.Context) {
    // 获取多分类参数
    categoryIds := c.QueryArray("categoryIds")

    filter := bson.M{}

    if len(categoryIds) > 0 {
        oids, err := converter.DTOIDsToModel(categoryIds)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid category IDs"})
            return
        }
        filter["category_ids"] = bson.M{"$in": oids}
    }

    // 继续查询...
}
```

---

## 数据迁移

### 现有数据处理

#### 场景1: BookDetail 使用 string ID

如果 `book_details` 集合中存储的是 string 类型的 ID：

```javascript
// 检查数据类型
db.book_details.findOne({}, {category_ids: 1})

// 如果是 string 数组，需要迁移
```

**迁移脚本** (`cmd/migrate/migrate_category_ids.go`):
```go
func MigrateCategoryIDs(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("book_details")

    // 查找所有文档
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)

    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return err
    }

    for _, doc := range results {
        // 获取 category_ids
        categoryIDsRaw := doc["category_ids"]

        var newCategoryIDs []primitive.ObjectID

        switch v := categoryIDsRaw.(type) {
        case []string:
            // string 数组 → ObjectID 数组
            for _, idStr := range v {
                oid, err := primitive.ObjectIDFromHex(idStr)
                if err == nil {
                    newCategoryIDs = append(newCategoryIDs, oid)
                }
            }
        case []interface{}:
            // interface{} 数组（可能是 ObjectID 或 string）
            for _, item := range v {
                switch idVal := item.(type) {
                case primitive.ObjectID:
                    newCategoryIDs = append(newCategoryIDs, idVal)
                case string:
                    oid, err := primitive.ObjectIDFromHex(idVal)
                    if err == nil {
                        newCategoryIDs = append(newCategoryIDs, oid)
                    }
                }
            }
        case []primitive.ObjectID:
            // 已经是正确类型，跳过
            continue
        }

        if len(newCategoryIDs) > 0 {
            _, err := collection.UpdateOne(
                ctx,
                bson.M{"_id": doc["_id"]},
                bson.M{"$set": bson.M{"category_ids": newCategoryIDs}},
            )
            if err != nil {
                log.Printf("Failed to update book detail %s: %v", doc["_id"], err)
            }
        }
    }

    return nil
}
```

---

## 验证清单

### 后端验证
- [ ] `BookDetail` 模型改为 `[]primitive.ObjectID`
- [ ] `BookDetailService` 类型更新
- [ ] DTO 转换函数正常工作
- [ ] Repository 查询正确处理 ObjectID 数组
- [ ] API 参数从 `categoryId` 改为 `categoryIds`

### 前端验证
- [ ] `Book` 类型改为 `categoryIds: string[]`
- [ ] 表单组件支持多选
- [ ] 列表显示支持多分类
- [ ] API 调用传递数组参数

### 数据验证
- [ ] 检查 `book_details` 集合中 category_ids 的实际类型
- [ ] 如需要，执行数据迁移脚本
- [ ] 验证迁移后的数据完整性

### 测试验证
- [ ] 单元测试更新
- [ ] 集成测试验证多分类查询
- [ ] API 测试验证参数正确性

---

## 实施计划

### Step 1: 准备（30分钟）
- [ ] 备份当前代码
- [ ] 创建新分支 `feature/category-ids-unification`
- [ ] 检查数据库中 category_ids 的实际类型

### Step 2: 后端迁移（2-3小时）
- [ ] 修改 `BookDetail` 模型类型
- [ ] 更新 `BookDetailService`
- [ ] 验证 DTO 转换函数
- [ ] 更新 API 参数和查询逻辑
- [ ] 运行测试验证

### Step 3: 前端迁移（2-3小时）
- [ ] 更新 `src/types/bookstore.ts`
- [ ] 搜索并更新所有使用 categoryId 的组件
- [ ] 更新表单组件支持多选
- [ ] 更新列表显示逻辑
- [ ] 运行前端测试

### Step 4: 数据迁移（如需要，1小时）
- [ ] 执行数据迁移脚本
- [ ] 验证数据完整性
- [ ] 回滚准备

### Step 5: 集成测试（1小时）
- [ ] 前后端联调测试
- [ ] 验证多分类创建/编辑
- [ ] 验证多分类查询
- [ ] 提交代码

---

## 回滚方案

如果迁移后出现问题：

1. **代码回滚**: `git revert <commit>`
2. **数据回滚**: 如已执行数据迁移，准备回滚脚本
3. **分支保护**: 使用 feature 分支，不影响主分支

---

## 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 数据库中存储string ID | 中 | 高 | 先检查数据，准备迁移脚本 |
| 前端组件兼容性 | 低 | 中 | 提供辅助函数平滑过渡 |
| API参数变更影响调用方 | 中 | 高 | 保留旧参数兼容性，标记废弃 |
| 多选UI改造复杂度 | 低 | 低 | 使用现有 Element Plus 组件 |

---

## 相关文档

- [Issue #011: 前后端数据类型不一致](../issues/011-frontend-backend-data-type-inconsistency.md)
- [ID类型统一标准](../standards/id-type-unification-standard.md)
- [P0问题审查报告](../reports/2026-03-05-p0-issue-audit-report.md)
- [BookStatus 枚举统一设计](./2026-03-05-book-status-unification-design.md)

---

**设计完成时间**: 2026-03-05
**预计实施时间**: 6-8小时
**建议执行者**: 后端 + 前端协同
