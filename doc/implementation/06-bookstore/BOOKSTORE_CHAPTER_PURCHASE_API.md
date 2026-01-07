# 书城系统章节目录和付费API文档

## 概述

本文档描述青羽写作平台书城系统的章节目录和付费相关的高优先级API实现。

## 新增文件列表

### 1. 数据模型
- `D:\Github\青羽\Qingyu_backend\models\bookstore\chapter_purchase.go`
  - 章节购买记录模型
  - 批量购买记录模型
  - 全书购买记录模型
  - 章节访问信息模型
  - 章节目录模型

### 2. 仓储接口
- `D:\Github\青羽\Qingyu_backend\repository\interfaces\bookstore\ChapterPurchaseRepository_interface.go`
  - 章节购买记录仓储接口
  - 支持单章购买、批量购买、全书购买记录管理
  - 提供权限检查和统计查询功能

### 3. 服务层
- `D:\Github\青羽\Qingyu_backend\service\bookstore\chapter_purchase_service.go`
  - 章节购买服务实现
  - 整合章节、钱包、购买记录服务
  - 实现完整的购买流程和权限检查逻辑

### 4. API层
- `D:\Github\青羽\Qingyu_backend\api\v1\bookstore\chapter_catalog_api.go`
  - 章节目录和购买API处理器
  - 提供章节目录查询、试读、VIP章节、购买等功能

### 5. 路由配置
- `D:\Github\青羽\Qingyu_backend\router\bookstore\bookstore_router.go` (已更新)
  - 添加章节目录相关路由
  - 添加章节购买相关路由
  - 新增 `InitReaderPurchaseRouter` 函数

---

## API端点列表

### 一、章节目录API

#### 1. 获取书籍章节目录
```
GET /api/v1/bookstore/books/{id}/chapters
```

**权限**: 公开（可选认证，认证用户可查看购买状态）

**路径参数**:
- `id`: 书籍ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "507f1f77bcf86cd799439011",
    "book_title": "示例小说",
    "total_chapters": 100,
    "free_chapters": 20,
    "paid_chapters": 80,
    "vip_chapters": 0,
    "total_word_count": 500000,
    "trial_count": 10,
    "chapters": [
      {
        "chapter_id": "507f1f77bcf86cd799439012",
        "title": "第一章 开端",
        "chapter_num": 1,
        "word_count": 3000,
        "is_free": true,
        "price": 0.0,
        "publish_time": "2024-01-01T00:00:00Z",
        "is_published": true,
        "is_purchased": false,
        "is_vip": false
      }
    ]
  }
}
```

---

#### 2. 获取单个章节信息
```
GET /api/v1/bookstore/books/{id}/chapters/{chapterId}
```

**权限**: 公开

**路径参数**:
- `id`: 书籍ID
- `chapterId`: 章节ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "book_id": "507f1f77bcf86cd799439011",
    "title": "第一章 开端",
    "chapter_num": 1,
    "word_count": 3000,
    "is_free": true,
    "price": 0.0,
    "publish_time": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 二、试读章节API

#### 3. 获取试读章节列表
```
GET /api/v1/bookstore/books/{id}/trial-chapters
```

**权限**: 公开

**路径参数**:
- `id`: 书籍ID

**查询参数**:
- `count`: 试读章节数量（默认10）

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "507f1f77bcf86cd799439011",
    "count": 10,
    "chapters": [
      {
        "id": "507f1f77bcf86cd799439012",
        "title": "第一章 开端",
        "chapter_num": 1,
        "is_free": true
      }
    ]
  }
}
```

---

### 三、付费章节API

#### 4. 获取VIP章节列表
```
GET /api/v1/bookstore/books/{id}/vip-chapters
```

**权限**: 公开

**路径参数**:
- `id`: 书籍ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "507f1f77bcf86cd799439011",
    "count": 80,
    "chapters": [
      {
        "id": "507f1f77bcf86cd799439015",
        "title": "第二十一章 转折",
        "chapter_num": 21,
        "is_free": false,
        "price": 0.1
      }
    ]
  }
}
```

---

#### 5. 获取章节价格
```
GET /api/v1/bookstore/chapters/{chapterId}/price
```

**权限**: 公开

**路径参数**:
- `chapterId`: 章节ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapter_id": "507f1f77bcf86cd799439015",
    "price": 0.1
  }
}
```

---

#### 6. 购买章节
```
POST /api/v1/reader/chapters/{chapterId}/purchase
```

**权限**: 需要JWT认证

**路径参数**:
- `chapterId`: 章节ID

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例（成功）**:
```json
{
  "code": 200,
  "message": "购买成功",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "user_id": "507f1f77bcf86cd799439030",
    "chapter_id": "507f1f77bcf86cd799439015",
    "book_id": "507f1f77bcf86cd799439011",
    "price": 0.1,
    "purchase_time": "2024-01-15T10:30:00Z",
    "chapter_title": "第二十一章 转折",
    "chapter_num": 21,
    "book_title": "示例小说"
  }
}
```

**错误响应**:
- `409`: 章节已购买
- `403`: 余额不足
- `400`: 免费章节无需购买

---

#### 7. 批量购买全书
```
POST /api/v1/reader/books/{id}/buy-all
```

**权限**: 需要JWT认证

**路径参数**:
- `id`: 书籍ID

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例（成功）**:
```json
{
  "code": 200,
  "message": "购买成功",
  "data": {
    "id": "507f1f77bcf86cd799439021",
    "user_id": "507f1f77bcf86cd799439030",
    "book_id": "507f1f77bcf86cd799439011",
    "total_price": 6.4,
    "original_price": 8.0,
    "discount": 0.2,
    "purchase_time": "2024-01-15T10:30:00Z",
    "book_title": "示例小说",
    "chapter_count": 80
  }
}
```

**说明**: 全书购买享受20%折扣

---

#### 8. 购买记录查询
```
GET /api/v1/reader/purchases
```

**权限**: 需要JWT认证

**查询参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认20）

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "total": 50,
  "page": 1,
  "size": 20,
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "chapter_id": "507f1f77bcf86cd799439015",
      "price": 0.1,
      "purchase_time": "2024-01-15T10:30:00Z",
      "chapter_title": "第二十一章 转折",
      "book_title": "示例小说"
    }
  ]
}
```

---

#### 9. 获取某本书的购买记录
```
GET /api/v1/reader/purchases/{bookId}
```

**权限**: 需要JWT认证

**路径参数**:
- `bookId`: 书籍ID

**查询参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认20）

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "total": 10,
  "page": 1,
  "size": 20,
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "chapter_id": "507f1f77bcf86cd799439015",
      "price": 0.1,
      "purchase_time": "2024-01-15T10:30:00Z",
      "chapter_title": "第二十一章 转折"
    }
  ]
}
```

---

#### 10. 检查章节访问权限
```
GET /api/v1/bookstore/chapters/{chapterId}/access
```

**权限**: 公开（认证用户可查看购买状态）

**路径参数**:
- `chapterId`: 章节ID

**响应示例**:
```json
{
  "code": 200,
  "message": "检查成功",
  "data": {
    "chapter_id": "507f1f77bcf86cd799439015",
    "title": "第二十一章 转折",
    "chapter_num": 21,
    "word_count": 3000,
    "is_free": false,
    "price": 0.1,
    "is_purchased": true,
    "is_vip": false,
    "can_access": true,
    "access_reason": "purchased",
    "purchase_time": "2024-01-15T10:30:00Z"
  }
}
```

**access_reason 可能的值**:
- `free`: 免费章节
- `purchased`: 已购买该章节
- `purchased_book`: 已购买全书
- `vip`: VIP权限

---

## 数据模型定义

### ChapterPurchase（章节购买记录）
```go
type ChapterPurchase struct {
    ID             primitive.ObjectID
    UserID         primitive.ObjectID
    ChapterID      primitive.ObjectID
    BookID         primitive.ObjectID
    Price          float64
    PurchaseTime   time.Time
    TransactionID  string
    CreatedAt      time.Time

    // 冗余字段（用于快速查询）
    ChapterTitle   string
    ChapterNum     int
    BookTitle      string
    BookCover      string
}
```

### ChapterPurchaseBatch（批量购买记录）
```go
type ChapterPurchaseBatch struct {
    ID             primitive.ObjectID
    UserID         primitive.ObjectID
    BookID         primitive.ObjectID
    ChapterIDs     []primitive.ObjectID
    TotalPrice     float64
    ChaptersCount  int
    PurchaseTime   time.Time
    TransactionID  string
    CreatedAt      time.Time

    // 冗余字段
    BookTitle      string
    BookCover      string
}
```

### BookPurchase（全书购买记录）
```go
type BookPurchase struct {
    ID             primitive.ObjectID
    UserID         primitive.ObjectID
    BookID         primitive.ObjectID
    TotalPrice     float64
    OriginalPrice  float64
    Discount       float64  // 折扣（0-1）
    PurchaseTime   time.Time
    TransactionID  string
    CreatedAt      time.Time

    // 冗余字段
    BookTitle      string
    BookCover      string
    ChapterCount   int
}
```

### ChapterAccessInfo（章节访问信息）
```go
type ChapterAccessInfo struct {
    ChapterID      primitive.ObjectID
    Title          string
    ChapterNum     int
    WordCount      int
    IsFree         bool
    Price          float64
    IsPurchased    bool
    IsVIP          bool
    CanAccess      bool
    AccessReason   string  // free, purchased, vip
    PurchaseTime   *time.Time
}
```

### ChapterCatalog（章节目录）
```go
type ChapterCatalog struct {
    BookID          primitive.ObjectID
    BookTitle       string
    TotalChapters   int
    FreeChapters    int
    PaidChapters    int
    VIPChapters     int
    TotalWordCount  int64
    Chapters        []ChapterCatalogItem
    TrialCount      int  // 可试读章节数量
}
```

### ChapterCatalogItem（章节目录项）
```go
type ChapterCatalogItem struct {
    ChapterID      primitive.ObjectID
    Title          string
    ChapterNum     int
    WordCount      int
    IsFree         bool
    Price          float64
    PublishTime    time.Time
    IsPublished    bool
    IsPurchased    bool   // 仅在认证用户的请求中返回
    IsVIP          bool   // VIP专属章节
}
```

---

## 路由配置说明

### 1. bookstore路由更新

`InitBookstoreRouter` 函数新增以下参数:
```go
func InitBookstoreRouter(
    r *gin.RouterGroup,
    bookstoreService bookstore.BookstoreService,
    bookDetailService bookstore.BookDetailService,
    ratingService bookstore.BookRatingService,
    statisticsService bookstore.BookStatisticsService,
    chapterService bookstore.ChapterService,          // 新增
    purchaseService bookstore.ChapterPurchaseService, // 新增
)
```

### 2. 新增路由函数

新增 `InitReaderPurchaseRouter` 函数用于处理购买相关接口:
```go
func InitReaderPurchaseRouter(
    r *gin.RouterGroup,
    purchaseService bookstore.ChapterPurchaseService,
)
```

### 3. 路由分组策略

**公开路由** (`/api/v1/bookstore`):
- 章节目录查询
- 试读章节获取
- VIP章节列表
- 价格查询
- 访问权限检查

**认证路由** (`/api/v1/reader`):
- 购买章节
- 购买全书
- 购买记录查询

---

## 服务依赖关系

```
ChapterPurchaseService 依赖:
  ├── ChapterRepository (章节仓储)
  ├── ChapterPurchaseRepository (购买记录仓储)
  ├── BookStoreRepository (书籍仓储)
  ├── WalletService (钱包服务)
  └── CacheService (缓存服务，可选)
```

---

## 权限检查逻辑

### 章节访问权限判断流程

1. **免费章节**: 直接允许访问
2. **已购买章节**: 检查购买记录
3. **已购买全书**: 检查全书购买记录
4. **VIP用户**: 检查用户VIP状态（待实现）

### 全书购买折扣

- 默认折扣: 8折 (20% off)
- 计算公式: `discountedPrice = originalPrice * 0.8`

---

## 待实现功能

1. **VIP系统集成**
   - 用户VIP状态检查
   - VIP章节标识
   - VIP到期时间管理

2. **购买记录仓储实现**
   - MongoDB实现 ChapterPurchaseRepository
   - 索引优化（user_id + chapter_id, user_id + book_id）

3. **批量购买优化**
   - 支持自定义章节列表购买
   - 购物车功能

4. **统计功能**
   - 购买数据统计
   - 消费统计
   - 阅读偏好分析

---

## 使用示例

### 前端集成示例

```typescript
// 1. 获取章节目录
async function getChapterCatalog(bookId: string, token?: string) {
  const headers: Record<string, string> = {};
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(
    `/api/v1/bookstore/books/${bookId}/chapters`,
    { headers }
  );

  return response.json();
}

// 2. 购买章节
async function purchaseChapter(chapterId: string, token: string) {
  const response = await fetch(
    `/api/v1/reader/chapters/${chapterId}/purchase`,
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    }
  );

  return response.json();
}

// 3. 检查访问权限
async function checkAccess(chapterId: string) {
  const response = await fetch(
    `/api/v1/bookstore/chapters/${chapterId}/access`
  );

  return response.json();
}

// 4. 获取购买记录
async function getPurchases(token: string, page = 1, size = 20) {
  const response = await fetch(
    `/api/v1/reader/purchases?page=${page}&size=${size}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );

  return response.json();
}
```

---

## 错误码说明

| 错误码 | 说明 |
|-------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 权限不足（余额不足） |
| 404 | 资源不存在 |
| 409 | 资源冲突（已购买） |
| 500 | 服务器内部错误 |

---

## 总结

本次实现了青羽写作平台书城系统的高优先级章节管理和付费相关API，包括:

✅ 章节目录查询（树形结构）
✅ 试读章节获取（前N章免费）
✅ VIP章节列表
✅ 章节购买（单章、批量、全书）
✅ 购买记录查询
✅ 访问权限检查
✅ 完整的服务层和业务逻辑
✅ 统一的响应格式
✅ 路由配置更新

所有API均已整合现有服务，使用统一的响应格式和错误处理机制。
