# Models 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **最后更新**: 2026-03-19

---

## 目录

1. [职责边界与依赖关系](#1-职责边界与依赖关系)
2. [命名与代码规范](#2-命名与代码规范)
3. [设计模式与最佳实践](#3-设计模式与最佳实践)
4. [接口与契约规范](#4-接口与契约规范)
5. [测试策略](#5-测试策略)
6. [完整代码示例](#6-完整代码示例)

---

## 1. 职责边界与依赖关系

### 1.1 核心职责

Models 层是数据结构的定义层，负责：

- **数据结构定义**：定义所有业务实体的数据结构
- **MongoDB 映射**：通过 `bson` tag 定义 MongoDB 字段映射
- **JSON 序列化**：通过 `json` tag 定义 API 响应格式
- **字段验证**：通过 `validate` tag 定义字段验证规则
- **枚举定义**：定义业务状态枚举及其行为方法
- **过滤器结构**：定义查询过滤条件结构体

### 1.2 依赖关系图

```
┌─────────────────────────────────────────────────────────┐
│                    上层依赖方                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Service   │  │ Repository  │  │     DTO     │     │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘     │
└─────────┼────────────────┼────────────────┼─────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────┐
│                     Models 层                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │
│  │  admin  │  │   ai    │  │  auth   │  │bookstore│   │
│  ├─────────┤  ├─────────┤  ├─────────┤  ├─────────┤   │
│  │ reader  │  │ social  │  │ finance │  │  ...    │   │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘   │
│       │            │            │            │         │
│       └────────────┴────────────┴────────────┘         │
│                          │                              │
│                          ▼                              │
│  ┌─────────────────────────────────────────────────┐   │
│  │              models/shared                       │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐       │   │
│  │  │  base.go │  │  types/  │  │ enums.go │       │   │
│  │  └──────────┘  └──────────┘  └──────────┘       │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    外部依赖                              │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │ go.mongodb.org/ │  │   time / fmt    │              │
│  │    mongo-driver │  │                 │              │
│  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────┘
```

### 1.3 层级边界

| 可以做 | 不应该做 |
|--------|----------|
| 定义数据结构 | 包含业务逻辑 |
| 定义枚举和常量 | 直接操作数据库 |
| 定义验证标签 | 包含数据库连接代码 |
| 定义辅助方法（如 `IsValid()`） | 调用 Service 或 Repository |
| 引用 shared 包 | 引用上层包（service, repository） |

### 1.4 模块划分

```
models/
├── admin/          # 管理员相关（封禁记录、导出历史）
├── ai/             # AI服务（会话、上下文、配额）
├── audit/          # 审计（敏感词、违规记录）
├── auth/           # 认证授权（JWT、会话、角色）
├── bookstore/      # 书城（书籍、章节、分类、评论）
├── dto/            # 数据传输对象
├── finance/        # 财务（订单、支付、钱包）
├── messaging/      # 消息（私信、群聊）
├── notification/   # 通知（系统通知、推送）
├── reader/         # 阅读（进度、书架、笔记）
├── recommendation/ # 推荐
├── search/         # 搜索索引
├── shared/         # 共享基础类型
│   └── types/      # 类型定义（ID、金额、评分等）
├── social/         # 社交（关注、点赞）
├── stats/          # 统计
└── storage/        # 存储
```

---

## 2. 命名与代码规范

### 2.1 文件命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 实体模型 | 小写单数名词 | `book.go`, `user.go`, `chapter.go` |
| 枚举定义 | 与主实体同名 | `book.go` 中的 `BookStatus` |
| 过滤器 | 实体名 + `_filter` | `book.go` 中的 `BookFilter` |
| 测试文件 | 原文件名 + `_test` | `book_test.go` |
| 文档说明 | 大写 README | `README.md` |

### 2.2 结构体命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 主实体 | PascalCase 单数 | `Book`, `User`, `Chapter` |
| 详情实体 | 实体名 + Detail | `BookDetail` |
| 统计实体 | 实体名 + Statistics | `BookStatistics` |
| 过滤器 | 实体名 + Filter | `BookFilter` |
| 基础实体 | Base + 用途 | `BaseEntity`, `IdentifiedEntity` |
| 混入结构 | 用途描述 | `ReadStatus`, `Edited` |

### 2.3 字段命名与标签

```go
type Book struct {
    // 字段名：PascalCase
    // bson tag：snake_case（MongoDB 字段名）
    // json tag：camelCase（API 响应）
    // validate tag：验证规则

    Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`

    // 可选字段使用 omitempty
    DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`

    // 外键字段使用 _id 后缀
    AuthorID string `bson:"author_id,omitempty" json:"authorId,omitempty"`

    // 列表字段使用复数
    CategoryIDs []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
}
```

### 2.4 枚举命名

```go
// 类型名：实体名 + 属性
type BookStatus string

// 常量：类型名前缀 + 状态值
const (
    BookStatusDraft     BookStatus = "draft"
    BookStatusOngoing   BookStatus = "ongoing"
    BookStatusCompleted BookStatus = "completed"
    BookStatusPaused    BookStatus = "paused"
)
```

### 2.5 方法命名

| 方法类型 | 命名规范 | 示例 |
|----------|----------|------|
| 验证方法 | `IsValid()` | `func (s BookStatus) IsValid() bool` |
| 状态检查 | `IsXxx()` | `func (s BookStatus) IsPublic() bool` |
| 转换方法 | `String()`, `ToXxx()` | `func (s BookStatus) String() string` |
| 业务行为 | 动词短语 | `func (b *BaseEntity) SoftDelete()` |
| 工厂函数 | `New` + 类型名 | `func NewBookFilter() *BookFilter` |

---

## 3. 设计模式与最佳实践

### 3.1 基础模型嵌入模式

使用结构体嵌入复用通用字段：

```go
// 基础实体 - 提供时间戳和软删除
type BaseEntity struct {
    CreatedAt time.Time  `json:"createdAt" bson:"created_at"`
    UpdatedAt time.Time  `json:"updatedAt" bson:"updated_at"`
    DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// ID 实体 - 提供 ID 字段
type IdentifiedEntity struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

// 业务实体 - 嵌入基础实体
type Book struct {
    IdentifiedEntity `bson:",inline"`  // 内联嵌入
    BaseEntity       `bson:",inline"`  // 内联嵌入

    Title string `bson:"title" json:"title"`
    // ... 其他业务字段
}
```

### 3.2 混入模式

为特定行为定义可复用的混入结构：

```go
// 已读状态混入
type ReadStatus struct {
    IsRead bool       `json:"isRead" bson:"is_read"`
    ReadAt *time.Time `json:"readAt,omitempty" bson:"read_at,omitempty"`
}

// 使用混入
type Notification struct {
    IdentifiedEntity `bson:",inline"`
    BaseEntity       `bson:",inline"`
    ReadStatus       `bson:",inline"`  // 混入已读状态

    Title   string `bson:"title" json:"title"`
    Content string `bson:"content" json:"content"`
}
```

### 3.3 枚举与行为方法

枚举类型应附带业务行为方法：

```go
type BookStatus string

const (
    BookStatusDraft     BookStatus = "draft"
    BookStatusOngoing   BookStatus = "ongoing"
    BookStatusCompleted BookStatus = "completed"
    BookStatusPaused    BookStatus = "paused"
)

// 验证方法
func (s BookStatus) IsValid() bool {
    switch s {
    case BookStatusDraft, BookStatusOngoing,
         BookStatusCompleted, BookStatusPaused:
        return true
    default:
        return false
    }
}

// 业务判断方法
func (s BookStatus) IsPublic() bool {
    return s == BookStatusOngoing || s == BookStatusCompleted
}

func (s BookStatus) CanEdit() bool {
    return s == BookStatusDraft || s == BookStatusPaused
}

// 字符串转换
func (s BookStatus) String() string {
    return string(s)
}
```

### 3.4 过滤器模式

为复杂查询定义过滤器结构体：

```go
type BookFilter struct {
    // 精确匹配
    CategoryID *string `json:"categoryId,omitempty"`
    AuthorID   *string `json:"authorId,omitempty"`
    Status     *BookStatus `json:"status,omitempty"`

    // 布尔过滤
    IsRecommended *bool `json:"isRecommended,omitempty"`
    IsFree        *bool `json:"isFree,omitempty"`

    // 范围过滤
    MinPrice *float64 `json:"minPrice,omitempty"`
    MaxPrice *float64 `json:"maxPrice,omitempty"`

    // 列表过滤
    Tags []string `json:"tags,omitempty"`

    // 模糊搜索
    Keyword *string `json:"keyword,omitempty"`

    // 分页
    SortBy    string  `json:"sortBy,omitempty"`
    SortOrder string  `json:"sortOrder,omitempty"`
    Limit     int     `json:"limit,omitempty"`
    Offset    int     `json:"offset,omitempty"`
    Cursor    *string `json:"cursor,omitempty"`
}
```

### 3.5 反模式警示

| 反模式 | 问题 | 正确做法 |
|--------|------|----------|
| 在模型中包含数据库操作 | 违反单一职责 | 模型只定义结构，操作放 Repository |
| 循环依赖 | 导致编译错误 | 使用接口或共享包解耦 |
| 过度使用指针 | 增加复杂度 | 仅在需要 nil 时使用指针 |
| 忽略 bson tag | 字段映射错误 | 始终显式定义 bson tag |
| 硬编码状态值 | 难以维护 | 使用枚举类型 |

---

## 4. 接口与契约规范

### 4.1 字段类型选择

| 场景 | 推荐类型 | 示例 |
|------|----------|------|
| 主键 | `primitive.ObjectID` | `ID primitive.ObjectID` |
| 外键（MongoDB） | `primitive.ObjectID` | `AuthorID primitive.ObjectID` |
| 外键（跨服务） | `string` | `ProjectID *string` |
| 金额 | `float64` 或 `int64`（分） | `Price float64` |
| 时间 | `time.Time` | `CreatedAt time.Time` |
| 可空时间 | `*time.Time` | `DeletedAt *time.Time` |
| 状态枚举 | 自定义 string 类型 | `BookStatus string` |
| 列表 | `[]T` | `Tags []string` |

### 4.2 Tag 规范

```go
type Example struct {
    // 必填字段
    Title string `bson:"title" json:"title" validate:"required"`

    // 可选字段（可能为空）
    Subtitle string `bson:"subtitle,omitempty" json:"subtitle,omitempty"`

    // 范围验证
    Age int `bson:"age" json:"age" validate:"min=0,max=150"`

    // 长度验证
    Name string `bson:"name" json:"name" validate:"min=1,max=100"`

    // 格式验证
    Email string `bson:"email" json:"email" validate:"email"`

    // URL 验证
    Cover string `bson:"cover" json:"cover" validate:"url"`

    // 枚举验证
    Status string `bson:"status" json:"status" validate:"oneof=draft published archived"`
}
```

### 4.3 错误处理

模型层不直接返回错误，但枚举方法可以：

```go
// 验证方法返回布尔值
func (s BookStatus) IsValid() bool

// 转换方法可以返回错误
func ParseBookStatus(s string) (BookStatus, error) {
    status := BookStatus(s)
    if !status.IsValid() {
        return "", fmt.Errorf("invalid book status: %s", s)
    }
    return status, nil
}
```

---

## 5. 测试策略

### 5.1 测试文件组织

```
models/
├── bookstore/
│   ├── book.go
│   ├── book_test.go           # 单元测试
│   └── book_status_test.go    # 枚举测试
└── shared/
    └── types/
        ├── enums.go
        └── enums_test.go
```

### 5.2 测试覆盖范围

| 测试类型 | 覆盖内容 | 示例 |
|----------|----------|------|
| 枚举验证 | `IsValid()` 方法 | 测试所有有效和无效值 |
| 状态判断 | `IsXxx()` 方法 | 测试各种状态组合 |
| JSON 序列化 | Marshal/Unmarshal | 验证 tag 正确性 |
| BSON 映射 | MongoDB 序列化 | 验证字段映射正确 |
| 验证标签 | validate tag | 使用 validator 库测试 |

### 5.3 测试示例

```go
// book_status_test.go
package bookstore

import "testing"

func TestBookStatus_IsValid(t *testing.T) {
    tests := []struct {
        status   BookStatus
        expected bool
    }{
        {BookStatusDraft, true},
        {BookStatusOngoing, true},
        {BookStatusCompleted, true},
        {BookStatusPaused, true},
        {BookStatus("invalid"), false},
        {BookStatus(""), false},
    }

    for _, tt := range tests {
        t.Run(string(tt.status), func(t *testing.T) {
            if got := tt.status.IsValid(); got != tt.expected {
                t.Errorf("BookStatus(%q).IsValid() = %v, want %v",
                    tt.status, got, tt.expected)
            }
        })
    }
}

func TestBookStatus_IsPublic(t *testing.T) {
    if !BookStatusOngoing.IsPublic() {
        t.Error("BookStatusOngoing should be public")
    }
    if !BookStatusCompleted.IsPublic() {
        t.Error("BookStatusCompleted should be public")
    }
    if BookStatusDraft.IsPublic() {
        t.Error("BookStatusDraft should not be public")
    }
}
```

### 5.4 测试覆盖率要求

| 类型 | 最低覆盖率 |
|------|-----------|
| 枚举方法 | 100% |
| 基础实体方法 | 90% |
| 验证逻辑 | 80% |

---

## 6. 完整代码示例

### 6.1 标准实体模型

```go
// bookstore/chapter.go
package bookstore

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"

    "Qingyu_backend/models/shared"
    "Qingyu_backend/models/shared/types"
)

// ChapterStatus 章节状态枚举
type ChapterStatus string

const (
    ChapterStatusDraft     ChapterStatus = "draft"
    ChapterStatusPublished ChapterStatus = "published"
    ChapterStatusRemoved   ChapterStatus = "removed"
)

// IsValid 检查状态是否有效
func (s ChapterStatus) IsValid() bool {
    switch s {
    case ChapterStatusDraft, ChapterStatusPublished, ChapterStatusRemoved:
        return true
    default:
        return false
    }
}

// IsPublished 是否已发布
func (s ChapterStatus) IsPublished() bool {
    return s == ChapterStatusPublished
}

// Chapter 章节模型
type Chapter struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联字段
    BookID primitive.ObjectID `bson:"book_id" json:"bookId" validate:"required"`

    // 基本信息字段
    Title       string `bson:"title" json:"title" validate:"required,min=1,max=200"`
    Content     string `bson:"content" json:"content,omitempty"`
    ContentURL  string `bson:"content_url,omitempty" json:"contentUrl,omitempty"`
    WordCount   int    `bson:"word_count" json:"wordCount" validate:"min=0"`
    ChapterNum  int    `bson:"chapter_num" json:"chapterNum" validate:"min=1"`

    // 状态与定价
    Status    ChapterStatus `bson:"status" json:"status" validate:"required"`
    IsVip     bool          `bson:"is_vip" json:"isVip"`
    Price     float64       `bson:"price" json:"price" validate:"min=0"`
    IsFree    bool          `bson:"is_free" json:"isFree"`

    // 统计信息
    ViewCount  int64 `bson:"view_count" json:"viewCount"`
    LikeCount  int64 `bson:"like_count" json:"likeCount"`

    // 时间字段
    PublishedAt *time.Time `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
}

// ChapterFilter 章节查询过滤器
type ChapterFilter struct {
    BookID    *string        `json:"bookId,omitempty"`
    Status    *ChapterStatus `json:"status,omitempty"`
    IsVip     *bool          `json:"isVip,omitempty"`
    MinPrice  *float64       `json:"minPrice,omitempty"`
    MaxPrice  *float64       `json:"maxPrice,omitempty"`
    Keyword   *string        `json:"keyword,omitempty"`
    SortBy    string         `json:"sortBy,omitempty"`
    SortOrder string         `json:"sortOrder,omitempty"`
    Limit     int            `json:"limit,omitempty"`
    Offset    int            `json:"offset,omitempty"`
}

// GetID 实现 Cacheable 接口
func (c Chapter) GetID() string {
    return c.ID.Hex()
}
```

### 6.2 基础实体定义

```go
// shared/base.go
package shared

import (
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseEntity 通用实体基类
type BaseEntity struct {
    CreatedAt time.Time  `json:"createdAt" bson:"created_at"`
    UpdatedAt time.Time  `json:"updatedAt" bson:"updated_at"`
    DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// Touch 更新时间戳
func (b *BaseEntity) Touch(t ...time.Time) {
    if len(t) > 0 {
        b.UpdatedAt = t[0]
    } else {
        b.UpdatedAt = time.Now()
    }
}

// TouchForCreate 创建时设置时间戳
func (b *BaseEntity) TouchForCreate() {
    now := time.Now()
    if b.CreatedAt.IsZero() {
        b.CreatedAt = now
    }
    if b.UpdatedAt.IsZero() {
        b.UpdatedAt = now
    }
}

// SoftDelete 软删除
func (b *BaseEntity) SoftDelete() {
    now := time.Now()
    b.DeletedAt = &now
    b.Touch(now)
}

// IsDeleted 判断是否已删除
func (b *BaseEntity) IsDeleted() bool {
    return b.DeletedAt != nil && !b.DeletedAt.IsZero()
}

// IdentifiedEntity 包含ID字段的基础实体
type IdentifiedEntity struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

// GetID 获取ID
func (i *IdentifiedEntity) GetID() primitive.ObjectID {
    return i.ID
}

// SetID 设置ID
func (i *IdentifiedEntity) SetID(id primitive.ObjectID) {
    i.ID = id
}

// GenerateID 生成新的ObjectID
func (i *IdentifiedEntity) GenerateID() {
    if i.ID.IsZero() {
        i.ID = primitive.NewObjectID()
    }
}
```

### 6.3 共享类型定义

```go
// shared/types/rating.go
package types

import (
    "errors"
    "math"
)

// ErrInvalidRating 无效评分错误
var ErrInvalidRating = errors.New("rating must be between 0 and 5")

// Rating 评分类型（0-5星）
type Rating struct {
    Value float64 `bson:"value" json:"value"`
    Count int64   `bson:"count" json:"count"`
}

// NewRating 创建新评分
func NewRating(value float64, count int64) (Rating, error) {
    if value < 0 || value > 5 {
        return Rating{}, ErrInvalidRating
    }
    return Rating{
        Value: math.Round(value*100) / 100, // 保留两位小数
        Count: count,
    }, nil
}

// AddRating 添加评分并更新平均值
func (r *Rating) AddRating(newRating float64) error {
    if newRating < 0 || newRating > 5 {
        return ErrInvalidRating
    }

    totalValue := r.Value * float64(r.Count)
    r.Count++
    r.Value = (totalValue + newRating) / float64(r.Count)
    r.Value = math.Round(r.Value*100) / 100

    return nil
}
```

---

## 附录

### A. 相关文档

- [Repository 层设计说明](./layer-repository.md)
- [Service 层设计说明](./layer-service.md)
- [DTO 层设计说明](./layer-dto.md)

### B. 参考资源

- [MongoDB Go Driver 文档](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [Go 结构体标签规范](https://pkg.go.dev/reflect#StructTag)
- [go-playground/validator](https://github.com/go-playground/validator)

---

*最后更新：2026-03-19*
