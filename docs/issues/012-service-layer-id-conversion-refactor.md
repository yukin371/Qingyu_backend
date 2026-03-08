# Issue #012: Service层ID转换重构

**创建日期**: 2026-03-08
**优先级**: P1 (重要)
**状态**: 待处理
**类型**: 架构重构

---

## 问题描述

当前项目中ID转换(`string` <-> `primitive.ObjectID`)分散在多个层级，违反了分层架构原则：

| 层级 | 当前ID转换数量 |
|------|----------------|
| Service层 | 169处 |
| Repository层 | 208处 |
| API层 | 91处 |

### 核心问题

1. **Service层混乱** - 部分Service直接使用`primitive.ObjectID`，部分使用`string`
2. **接口不一致** - 部分Service接口暴露了MongoDB实现细节
3. **违反分层原则** - Service层不应该知道数据库的技术实现

---

## 推荐方案

**Repository层负责ID转换**

| 层级 | ID类型 | 职责 |
|------|--------|------|
| **API层** | `string` | 从请求中获取string ID |
| **Service层** | `string` | 业务逻辑，传递string ID给Repository |
| **Repository接口** | `string` | 定义为string类型 |
| **Repository实现** | 内部转换为 `ObjectID` | 与MongoDB交互 |

### 理由

1. **符合现有接口设计** - Repository接口已定义为`string`类型
2. **符合分层架构原则** - Repository是唯一知道MongoDB细节的层
3. **便于数据库迁移** - 未来切换数据库只需修改Repository层
4. **已有基础设施** - `BaseMongoRepository`已提供`ParseID()`方法

---

## 受影响的文件

### 高优先级

| 文件 | 问题 | 修改内容 |
|------|------|----------|
| `service/bookstore/book_rating_service.go` | 接口使用`primitive.ObjectID` | 改为`string` |
| `service/bookstore/chapter_purchase_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/bookstore/book_statistics_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/finance/author_revenue_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/finance/membership_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/messaging/announcement_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/audit/content_audit_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/admin/user_admin_service.go` | Service层ID转换 | 移除转换，传递string |
| `service/container/service_container.go` | Service层ID转换 | 移除转换，传递string |

### 中优先级

| 文件 | 问题 |
|------|------|
| `service/bookstore/chapter_purchase_service.go` | 大量内联ID转换 |
| `service/user/transaction_manager.go` | 事务管理器特殊情况 |

---

## 修改规范

### Service层修改示例

**修改前**:
```go
func (s *ChapterPurchaseService) PurchaseChapter(ctx context.Context, userID, bookID, chapterID string) error {
    userOID, _ := primitive.ObjectIDFromHex(userID)
    bookOID, _ := primitive.ObjectIDFromHex(bookID)
    chapterOID, _ := primitive.ObjectIDFromHex(chapterID)
    // 使用ObjectID调用Repository...
}
```

**修改后**:
```go
func (s *ChapterPurchaseService) PurchaseChapter(ctx context.Context, userID, bookID, chapterID string) error {
    // 直接传递string给Repository
    return s.purchaseRepo.Create(ctx, userID, bookID, chapterID)
}
```

### Service接口修改示例

**修改前**:
```go
type BookRatingService interface {
    GetRatingByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookRating, error)
}
```

**修改后**:
```go
type BookRatingService interface {
    GetRatingByID(ctx context.Context, id string) (*bookstore.BookRating, error)
}
```

### Repository层确保规范

Repository实现应遵循:
```go
func (r *MongoBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
    objectID, err := r.ParseID(id)  // 使用BaseMongoRepository的方法
    if err != nil {
        return nil, types.ErrInvalidIDFormat
    }
    // 内部使用ObjectID进行查询...
}
```

---

## 实施计划

### 阶段1: 接口统一 (1-2天)
1. 修改所有Service接口，将`primitive.ObjectID`参数改为`string`
2. 更新相关的Mock文件

### 阶段2: Service层重构 (2-3天)
1. 移除Service层的ID转换代码
2. 确保Service层只传递string给Repository

### 阶段3: 测试验证 (1-2天)
1. 更新单元测试
2. 运行集成测试验证

### 阶段4: 文档更新 (0.5天)
1. 更新架构文档
2. 添加开发规范

---

## 验收标准

- [ ] 所有Service接口使用`string`类型ID参数
- [ ] Service层不再直接使用`primitive.ObjectIDFromHex`
- [ ] Repository接口保持`string`类型
- [ ] 所有测试通过
- [ ] CodeQL无新增警报

---

## 关联Issue

- #010 Repository层业务逻辑渗透
- #011 前后端数据类型不一致

---

*文档由 Kore 创建于 2026-03-08*
