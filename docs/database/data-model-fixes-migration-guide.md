# 书城数据模型修复与迁移指南

> **文档版本**: v1.0
> **创建时间**: 2026-01-22
> **适用范围**: Qingyu Backend 书城模块

---

## 📋 概述

本文档记录了书城模块数据模型的重大修复，以及相应的数据迁移方案。修复解决了评分范围不一致、MongoDB 序列化问题、状态枚举冲突、金额精度问题和内容哈希计算等关键问题。

---

## 🎯 修复内容

### 1. 评分范围统一 (High 优先级)

**问题**: 评分字段范围不一致，可能导致展示/统计错误

| 模型 | 原范围 | 新范围 | 文件 |
|------|--------|--------|------|
| Book.Rating | 0-10 | 1-5 | models/bookstore/book.go:32 |
| BookDetail.Rating | 0-5 | 1-5 | models/bookstore/book_detail.go:42 |
| BookRating.Rating | 1-5 | 1-5 | models/bookstore/book_rating.go:14 |
| BookStatistics.RatingDistribution | int键 | string键 | models/bookstore/book_statistics.go:19 |

**修复方案**:
```go
// 修复前
Rating float64 `bson:"rating" json:"rating"` // 评分 (0-10)
RatingDistribution map[int]int64 // 键为整数

// 修复后
Rating float64 `bson:"rating" json:"rating" validate:"min=1,max=5"` // 评分 (1-5星)
RatingDistribution map[string]int64 // 键为 "1"-"5" 字符串
```

### 2. MongoDB 键类型修复 (High 优先级)

**问题**: `map[int]int64` 在 MongoDB 中 BSON 序列化时，键会被自动转换为字符串，导致读写不一致

**修复方案**:
- 将 `RatingDistribution` 从 `map[int]int64` 改为 `map[string]int64`
- 更新 `UpdateRating`、`RemoveRating`、`calculateAverageRating` 方法使用字符串键
- 使用 `strconv.Itoa` 和 `strconv.Atoi` 进行安全转换

```go
// 修复后
func (bs *BookStatistics) UpdateRating(rating int) {
    if bs.RatingDistribution == nil {
        bs.RatingDistribution = make(map[string]int64)
    }
    ratingKey := strconv.Itoa(rating) // 1-5 -> "1"-"5"
    bs.RatingDistribution[ratingKey]++
    // ...
}
```

### 3. BookStatus 状态枚举优化 (Medium 优先级)

**问题**: `published` 和 `ongoing` 状态语义重叠，状态流转不清晰

**修复方案**: 移除 `BookStatusPublished`

```go
// 修复前
const (
    BookStatusDraft     BookStatus = "draft"     // 草稿
    BookStatusPublished BookStatus = "published" // 已发布
    BookStatusOngoing   BookStatus = "ongoing"   // 连载中
    BookStatusCompleted BookStatus = "completed" // 已完结
    BookStatusPaused    BookStatus = "paused"    // 暂停更新
)

// 修复后
const (
    BookStatusDraft     BookStatus = "draft"     // 草稿
    BookStatusOngoing   BookStatus = "ongoing"   // 连载中 (已发布且正在更新)
    BookStatusCompleted BookStatus = "completed" // 已完结
    BookStatusPaused    BookStatus = "paused"    // 暂停更新
)
```

**影响的代码文件**:
- `service/bookstore/bookstore_service.go:318`
- `repository/mongodb/bookstore/ranking_repository_mongo.go` (4处)
- `migration/seeds/import_novels.go:223`

### 4. 金额字段精度修复 (Medium 优先级)

**问题**: 使用 `float64` 存储金额存在精度风险，支付结算、折扣计算可能出现误差

**修复方案**: 改为 `int64` 存储"分"

| 模型/字段 | 原类型 | 新类型 | 文件 |
|-----------|--------|--------|------|
| Book.Price | float64 | int64 | models/bookstore/book.go:37 |
| BookDetail.Price | float64 | int64 | models/bookstore/book_detail.go:33 |
| Chapter.Price | float64 | int64 | models/bookstore/chapter.go:18 |
| ChapterPurchase.Price | float64 | int64 | models/bookstore/chapter_purchase.go:15 |
| ChapterPurchaseBatch.TotalPrice | float64 | int64 | models/bookstore/chapter_purchase.go:33 |
| BookPurchase.TotalPrice | float64 | int64 | models/bookstore/chapter_purchase.go:49 |
| BookPurchase.OriginalPrice | float64 | int64 | models/bookstore/chapter_purchase.go:50 |
| ChapterAccessInfo.Price | float64 | int64 | models/bookstore/chapter_purchase.go:96 |
| ChapterCatalogItem.Price | float64 | int64 | models/bookstore/chapter_purchase.go:111 |

**注意**: `BookPurchase.Discount` 折扣字段保持 `float64` (0-1)

### 5. 内容哈希计算修复 (Low 优先级)

**问题**: `CalculateHash` 未基于内容，仅用 ChapterID+Version，无法做内容校验或去重

**修复方案**: 使用 SHA-256 计算真实内容哈希

```go
// 修复前
func (cc *ChapterContent) CalculateHash() string {
    return fmt.Sprintf("%s:%d", cc.ChapterID.Hex(), cc.Version)
}

// 修复后
func (cc *ChapterContent) CalculateHash() string {
    hash := sha256.Sum256([]byte(cc.Content))
    return hex.EncodeToString(hash[:])
}
```

---

## 📊 数据迁移方案

### 评分数据迁移

对于现有的 Book 和 BookDetail 数据，需要将 0-10 范围的评分转换为 1-5 范围：

```javascript
// MongoDB 迁移脚本
db.books.find({ rating: { $exists: true, $gt: 5 } }).forEach(function(doc) {
    var newRating = doc.rating / 2;
    db.books.updateOne(
        { _id: doc._id },
        { $set: { rating: newRating } }
    );
});

db.book_details.find({ rating: { $exists: true, $gt: 5 } }).forEach(function(doc) {
    var newRating = doc.rating / 2;
    db.book_details.updateOne(
        { _id: doc._id },
        { $set: { rating: newRating } }
    );
});
```

### RatingDistribution 键类型迁移

```javascript
// 将整数键转换为字符串键
db.book_statistics.find({ rating_distribution: { $exists: true } }).forEach(function(doc) {
    var newDist = {};
    for (var key in doc.rating_distribution) {
        newDist[String(key)] = doc.rating_distribution[key];
    }
    db.book_statistics.updateOne(
        { _id: doc._id },
        { $set: { rating_distribution: newDist } }
    );
});
```

### BookStatus 状态迁移

```javascript
// 将 published 状态迁移为 ongoing
db.books.updateMany(
    { status: "published" },
    { $set: { status: "ongoing" } }
);

db.book_details.updateMany(
    { status: "published" },
    { $set: { status: "ongoing" } }
);
```

### 金额数据迁移

```javascript
// 将 float64 价格转换为 int64 (分)
// 注意：此迁移假设原始价格单位为元
db.books.find({ price: { $exists: true, $type: "double" } }).forEach(function(doc) {
    var newPrice = Math.round(doc.price * 100);
    db.books.updateOne(
        { _id: doc._id },
        { $set: { price: newPrice } }
    );
});

db.chapters.find({ price: { $exists: true, $type: "double" } }).forEach(function(doc) {
    var newPrice = Math.round(doc.price * 100);
    db.chapters.updateOne(
        { _id: doc._id },
        { $set: { price: newPrice } }
    );
});

// 类似地处理其他价格字段...
```

---

## 🔍 前端影响

### API 响应变更

前端需要适配以下变更：

1. **价格字段**: 从 `number` (元) 变为 `number` (分)，需要除以 100 显示
2. **评分范围**: 从 0-10 变为 1-5，评分组件需要适配
3. **BookStatus**: `published` 状态不再存在，前端代码需要处理该状态为 `ongoing`

### 前端适配示例

```typescript
// 价格显示适配
function formatPrice(priceInCents: number): string {
  return (priceInCents / 100).toFixed(2);
}

// 评分组件适配
// 评分组件的最大值应从 10 改为 5

// 状态处理
function getStatusLabel(status: string): string {
  const statusMap = {
    'draft': '草稿',
    'ongoing': '连载中',
    'completed': '已完结',
    'paused': '暂停更新'
  };
  return statusMap[status] || status;
}
```

---

## ✅ 验证清单

完成修复和迁移后，请验证以下内容：

- [ ] Book.Rating 验证规则为 `min=1,max=5`
- [ ] BookDetail.Rating 验证规则为 `min=1,max=5`
- [ ] BookStatistics.RatingDistribution 键为字符串类型
- [ ] BookStatus 不再包含 `published` 枚举值
- [ ] 所有价格字段类型为 `int64` (分)
- [ ] ChapterContent.CalculateHash 基于 SHA-256 内容计算
- [ ] 数据库中不存在 `status: "published"` 的记录
- [ ] 前端价格显示正确（除以 100）
- [ ] 前端评分组件最大值为 5
- [ ] 相关测试已更新并通过

---

## 📝 注意事项

1. **数据备份**: 执行迁移前务必备份数据库
2. **停机维护**: 建议在低峰期执行迁移，可能需要短暂停机
3. **分步执行**: 建议分步执行迁移，每步验证后再进行下一步
4. **回滚方案**: 准备好回滚脚本以防迁移失败
5. **测试验证**: 在测试环境完整验证后再在生产环境执行

---

## 🔗 相关文件

修改的文件列表：
- `Qingyu_backend/models/bookstore/book.go`
- `Qingyu_backend/models/bookstore/book_detail.go`
- `Qingyu_backend/models/bookstore/book_rating.go`
- `Qingyu_backend/models/bookstore/book_statistics.go`
- `Qingyu_backend/models/bookstore/chapter.go`
- `Qingyu_backend/models/bookstore/chapter_content.go`
- `Qingyu_backend/models/bookstore/chapter_purchase.go`
- `Qingyu_backend/service/bookstore/bookstore_service.go`
- `Qingyu_backend/repository/mongodb/bookstore/ranking_repository_mongo.go`
- `Qingyu_backend/migration/seeds/import_novels.go`
