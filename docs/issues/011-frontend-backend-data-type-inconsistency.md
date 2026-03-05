# Issue #011: 前后端数据类型不一致

**优先级**: 高 (P0)
**类型**: 兼容性问题
**状态**: 待处理
**创建日期**: 2026-03-05
**来源报告**: [前后端数据类型对比报告](../reports/archived/2026-03-04-frontend-backend-data-type-comparison-report.md)、[类型转换兼容性分析](../reports/archived/type-conversion-compatibility-analysis.md)

---

## 问题描述

前后端数据类型定义存在28处不一致，导致数据传输和显示问题。

### 问题规模

| 严重程度 | 问题数量 | 说明 |
|---------|----------|------|
| P0（阻塞） | 5 | 必须立即修复 |
| P1（重要） | 12 | 应该尽快处理 |
| P2（一般） | 11 | 可以后续优化 |

---

## 主要问题分类

### 1. 枚举值不一致 🔴 P0

#### BookStatus 枚举不匹配

**问题**: 前端和后端的 BookStatus 枚举值不一致。

```typescript
// ❌ 前端: src/types/bookstore.ts
export type BookStatus = 'serializing' | 'completed' | 'paused'
```

```go
// ✅ 后端: models/bookstore/book.go
const (
    BookStatusDraft     BookStatus = "draft"      // 前端没有
    BookStatusOngoing   BookStatus = "ongoing"    // 前端期望 'serializing'
    BookStatusCompleted BookStatus = "completed"  // 匹配
    BookStatusPaused    BookStatus = "paused"     // 匹配
    BookStatusDeleted   BookStatus = "deleted"    // 前端没有
)
```

**影响**: 书籍列表、详情页面、筛选功能

**修复方案**: 修改前端枚举定义为 `'draft' | 'ongoing' | 'completed' | 'paused' | 'deleted'`

---

### 2. 布尔字段转换不一致 🔴 P0

#### is_* 字段转换遗漏

**问题**: 后端使用 snake_case (`is_free`, `is_hot`)，前端期望 camelCase，某些字段可能遗漏转换。

```go
// 后端正确配置
type Book struct {
    IsFree        bool `bson:"is_free" json:"isFree"`
    IsRecommended bool `bson:"is_recommended" json:"isRecommended"`
    IsFeatured    bool `bson:"is_featured" json:"isFeatured"`
    IsHot         bool `bson:"is_hot" json:"isHot"`
}
```

**需要检查的字段**:
| 后端BSON | 后端JSON | 前端期望 | 状态 |
|---------|---------|---------|------|
| `is_free` | `isFree` | `isFree` | ✅ |
| `is_vip` | `isVip` | `isVip` | ⚠️ 需验证 |
| `has_next` | `hasNext` | `hasNext` | ⚠️ 需验证 |

**修复方案**: 批量检查所有 `is_*` 字段的 JSON 标签

---

### 3. 数组类型不匹配 🔴 P0

#### CategoryIDs 类型不一致

**问题**: 后端使用 ObjectId 数组，前端使用 string 单值或数组。

```go
// 后端: models/bookstore/book.go
type Book struct {
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
    Categories    []string             `bson:"categories" json:"categories"` // 冗余字段
}
```

```typescript
// ❌ 前端: src/types/bookstore.ts
interface Book {
  categoryId: string        // 单值，不匹配后端数组
  categoryName?: string
  category?: string
}
```

**修复方案**:
1. 前端改为 `categoryIds: string[]`
2. 后端 DTO 转换 `ObjectId[]` → `string[]`

---

### 4. 响应拦截器处理不一致 🔴 P0

**问题**: 前端响应拦截器处理分页响应的方式可能导致字段丢失。

```typescript
// src/core/services/http.service.ts
apiClient.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res && typeof res === 'object' && 'code' in res && 'data' in res) {
      if ('pagination' in res) {
        return res  // 保留pagination
      }
      return res.data  // 只返回data
    }
    return res
  }
)
```

**问题**:
1. 某些接口返回 `{code, data, pagination}`
2. 某些接口返回 `{code, data}`
3. 拦截器处理不一致，可能丢失字段

**修复方案**: 统一响应格式，拦截器始终返回完整响应

---

### 5. 金额单位未转换 🟡 P1

#### Price 字段类型和单位不一致

**问题**:
- **后端**: `Price float64` (分，但使用 float64)
- **前端**: `price?: number` (期望元)

```go
// 后端: models/bookstore/book.go
Price float64 `bson:"price" json:"price" validate:"min=0"` // 分，使用float64
```

**影响**: 价格显示错误，可能是实际价格的100倍

**修复方案**:
1. 前端统一除以100转换
2. 长期迁移到 `types.Money` 类型

```typescript
// 前端转换工具
function formatPrice(cents: number): string {
  return `¥${(cents / 100).toFixed(2)}`
}
```

---

### 6. 时间字段命名不一致 🟡 P1

**问题**: 不同模型使用不同的时间字段名。

| 字段用途 | 后端命名 | 前端期望 | 状态 |
|---------|---------|---------|------|
| 创建时间 | `CreatedAt` | `createdAt` | ✅ |
| 更新时间 | `UpdatedAt` | `updateTime` | ⚠️ 不一致 |
| 发布时间 | `PublishedAt` | `publishTime` | ⚠️ 不一致 |

**修复方案**: 统一命名规范 `createdAt`, `updatedAt`, `publishedAt`

---

### 7. 时间格式处理不一致 🟡 P1

**问题**:
- **后端**: `*time.Time` (Go时间类型)
- **前端期望**: ISO 8601 string

```go
// 后端可能序列化为多种格式
PublishedAt   *time.Time `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
```

**问题**:
1. `*time.Time` 可能被序列化为多种格式
2. `omitempty` 可能导致前端收到 null

**修复方案**:
1. 后端统一使用 RFC3339 格式
2. 前端添加时间解析工具
3. 处理 null/undefined 情况

---

### 8. 其他枚举不一致 🟡 P1

#### UserRole 枚举命名不一致

```typescript
// 前端
type UserRole = 'admin' | 'writer' | 'user'
```

```go
// 后端
const (
    RoleReader Role = "reader"
    RoleAuthor Role = "author"  // 前端使用 'writer'
    RoleAdmin  Role = "admin"
)
```

**修复方案**: 前端改为 `'admin' | 'author' | 'user' | 'reader'`

#### BehaviorType 枚举不一致

```typescript
// 前端独有
'collect' | 'search'

// 后端独有
'finish' | 'share'
```

**修复方案**: 合并前后端定义

---

### 9. DocumentContent V2 支持缺失 🟡 P1

**问题**: 后端已将 Document 内容分离到 DocumentContent 集合，前端可能仍在使用旧的 Document.content 字段。

```go
// V2 (新)
type Document struct {
    // 内容已移除
}

type DocumentContent struct {
    DocumentID primitive.ObjectID `bson:"document_id"`
    Content    string             `bson:"content"`
    Version    int                `bson:"version"`
}
```

**修复方案**:
1. 前端添加 DocumentContent 类型定义
2. 更新 API 调用使用新端点
3. 实现内容版本切换 UI

---

### 10. ID 类型转换边界不清晰 🟡 P1

**问题**: 根据 `id-type-unification-standard.md`，转换应该在 Service 层进行，但实际代码中可能存在混用。

**理想架构**:
```
API Layer (string) → Service Layer (string) → Repository Layer (ObjectID)
```

**实际问题**:
- 某些 API Handler 可能直接使用 `ObjectID`
- 某些 Service 可能返回包含 `ObjectID` 的结构体

**修复方案**:
1. 审查所有 API Handler，确保 DTO 使用 string
2. 审查所有 Service 接口，确保边界使用 string
3. 在 Service 层调用 Repository 前转换 string → ObjectID

---

## 实施计划

### Phase 1: P0 阻塞问题（1-2 天）

| 问题 | 修复方案 | 预计时间 |
|------|---------|----------|
| BookStatus 枚举值不一致 | 前端改为 `'draft'\|'ongoing'\|'completed'\|'paused'\|'deleted'` | 2小时 |
| is_* 字段转换遗漏 | 批量检查所有 Model JSON 标签 | 2小时 |
| CategoryIDs 数组类型 | 前端改为数组，后端 DTO 转换 | 3小时 |
| 响应拦截器处理不一致 | 统一响应格式处理 | 2小时 |
| snake_case → camelCase 遗漏 | 批量检查脚本 | 3小时 |

### Phase 2: P1 重要问题（3-5 天）

| 问题 | 修复方案 | 预计时间 |
|------|---------|----------|
| Price 字段类型和单位 | 前端统一除以100转换 | 2小时 |
| 时间字段命名不一致 | 统一命名规范 | 3小时 |
| 时间格式处理不一致 | 统一使用 ISO 8601 | 3小时 |
| UserRole 枚举命名不一致 | 前端改为 'author' | 1小时 |
| BehaviorType 枚举不一致 | 合并前后端定义 | 2小时 |
| DocumentContent V2 支持 | 前端添加新类型 | 1天 |
| ID 类型转换边界不清晰 | 明确 Service 层边界 | 4小时 |
| 可选字段指针处理 | 明确 null 处理 | 3小时 |

### Phase 3: P2 一般问题（长期优化）

- 统一分页参数命名
- 明确软删除策略
- 统一 ID 类型 (UUID vs ObjectId)
- 优化时间字段命名
- 添加类型注释和文档

---

## 检查清单

### 字段转换验证
- [ ] 所有 `is_*` 字段正确转换为 `isXxx`
- [ ] 所有 snake_case 字段转换为 camelCase
- [ ] BookStatus 枚举值前后端一致
- [ ] 时间字段使用 ISO 8601 格式

### 类型转换验证
- [ ] ID 类型在 Service 层正确转换
- [ ] CategoryIDs 正确转换为数组
- [ ] Price 字段正确转换(分 ↔ 元)
- [ ] 可选字段正确处理 null

### 功能验证
- [ ] 书籍列表正确显示
- [ ] 书籍详情正确加载
- [ ] 分类筛选正常工作
- [ ] 分页功能正常
- [ ] 表单提交成功
- [ ] 时间正确显示

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [前后端数据类型对比报告](../reports/archived/2026-03-04-frontend-backend-data-type-comparison-report.md) | 完整类型对比分析 |
| [类型转换兼容性分析](../reports/archived/type-conversion-compatibility-analysis.md) | 详细问题清单 |
| [后端 Book 模型](../../models/bookstore/book.go) | 后端数据模型 |
| [前端 Bookstore 类型](../../../Qingyu_fronted/src/types/bookstore.ts) | 前端类型定义 |

---

## 相关 Issue

- [#001: 统一模型层 ID 字段类型](./001-unify-id-type-in-models.md)
- [#005: API 标准化问题](./005-api-standardization-issues.md)
- [#006: 数据库索引问题](./006-database-index-issues.md)
