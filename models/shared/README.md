# 共享基础模型

> 可复用的基础数据模型混入（Mixins），供各业务领域模型组合使用

---

## 概述

`models/shared/` 目录包含可复用的基础模型组件，通过 **组合模式** 为各业务领域模型提供通用功能。

### 设计原则

1. **组合优于继承** - 通过 `bson:",inline"` 将基础模型嵌入到领域模型中
2. **单一职责** - 每个混入文件专注于一类通用功能
3. **类型安全** - 提供方法封装，确保数据一致性
4. **向后兼容** - 部分模块通过类型别名保持兼容性

---

## 📋 基础模型列表

### 1. base.go - 核心基础模型

**文件**: `base.go`

**核心模型**:
- `BaseEntity` - 通用实体基类（时间戳）
- `IdentifiedEntity` - ID 字段实体
- `ReadStatus` - 已读状态混入
- `Edited` - 编辑追踪混入

```go
// BaseEntity - 时间戳基础
type BaseEntity struct {
    CreatedAt time.Time  `json:"createdAt" bson:"created_at"`
    UpdatedAt time.Time  `json:"updatedAt" bson:"updated_at"`
    DeletedAt time.Time  `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// 方法:
// - Touch(t ...time.Time) - 更新时间戳
// - TouchForCreate() - 创建时设置时间戳
// - SoftDelete() - 软删除
// - IsDeleted() - 判断是否已删除

// IdentifiedEntity - ID 字段
type IdentifiedEntity struct {
    ID string `bson:"_id,omitempty" json:"id"`
}

// 方法:
// - GetID() string - 获取ID
// - SetID(id string) - 设置ID

// ReadStatus - 已读状态
type ReadStatus struct {
    IsRead bool       `json:"isRead" bson:"is_read"`
    ReadAt *time.Time `json:"readAt,omitempty" bson:"read_at,omitempty"`
}

// 方法:
// - MarkAsRead() - 标记为已读
// - MarkAsUnread() - 标记为未读
// - IsRecentlyRead(minutes int) bool - 检查是否在最近N分钟内已读

// Edited - 编辑追踪
type Edited struct {
    LastSavedAt  time.Time `json:"lastSavedAt" bson:"last_saved_at"`
    LastEditedBy string    `json:"lastEditedBy" bson:"last_edited_by"`
}

// 方法:
// - MarkEdited(editorID string) - 标记为已编辑
// - GetLastSavedAt() time.Time - 获取最后保存时间
// - GetLastEditedBy() string - 获取最后编辑人
```

---

### 2. social.go - 社交功能混入

**文件**: `social.go`

**核心模型**:
- `Likable` - 点赞功能混入
- `ThreadedConversation` - 会话嵌套混入

```go
// Likable - 点赞功能
type Likable struct {
    LikeCount int `bson:"like_count" json:"likeCount"`
}

// 方法:
// - AddLike(count int) - 增加点赞数
// - RemoveLike(count int) - 减少点赞数

// ThreadedConversation - 会话嵌套
type ThreadedConversation struct {
    ReplyToCommentID *string `bson:"reply_to_comment_id,omitempty" json:"replyToCommentId,omitempty"`
    RootID           *string `bson:"root_id,omitempty" json:"rootId,omitempty"`
    ReplyCount       int     `bson:"reply_count" json:"replyCount"`
}

// 字段说明:
// - ReplyToCommentID: 直接回复的评论ID（指针类型，可为nil）
// - RootID: 根评论ID（用于多级回复，指针类型，可为nil）
// - ReplyCount: 回复数量
```

---

### 3. communication.go - 通信功能混入

**文件**: `communication.go`

**核心模型**:
- `CommunicationBase` - 通信基础实体

```go
type CommunicationBase struct {
    SenderID   string     `bson:"sender_id" json:"senderId" validate:"required"`
    ReceiverID string     `bson:"receiver_id" json:"receiverId" validate:"required"`
    IsRead     bool       `bson:"is_read" json:"isRead"`
    ReadAt     *time.Time `bson:"read_at,omitempty" json:"readAt,omitempty"`
}

// 方法:
// - MarkAsRead() - 标记为已读
// - MarkAsUnread() - 标记为未读
```

---

### 4. content.go - 内容相关混入

**文件**: `content.go`

**核心模型**:
- `TitledEntity` - 标题实体
- `NamedEntity` - 名称实体
- `DescriptedEntity` - 描述实体

```go
type TitledEntity struct {
    Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`
}

type NamedEntity struct {
    Name string `bson:"name" json:"name" validate:"required,min=1,max=100"`
}

type DescriptedEntity struct {
    Description string `bson:"description,omitempty" json:"description,omitempty" validate:"max=1000"`
}
```

---

### 5. metadata.go - 元数据混入

**文件**: `metadata.go`

**核心模型**:
- `Pinned` - 置顶状态
- `Expirable` - 有效期
- `TargetEntity` - 关联实体

```go
// Pinned - 置顶状态
type Pinned struct {
    IsPinned bool       `bson:"is_pinned" json:"isPinned"`
    PinnedAt *time.Time `bson:"pinned_at,omitempty" json:"pinnedAt,omitempty"`
    PinnedBy *string    `bson:"pinned_by,omitempty" json:"pinnedBy,omitempty"`
}

// 方法:
// - Pin(operatorID string) - 置顶
// - Unpin() - 取消置顶

// Expirable - 有效期
type Expirable struct {
    ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expiresAt,omitempty"`
}

// 方法:
// - IsExpired() bool - 判断是否已过期
// - SetExpiration(duration time.Duration) - 设置过期时间

// TargetEntity - 关联实体
type TargetEntity struct {
    TargetType string `bson:"target_type" json:"targetType"`
    TargetID   string `bson:"target_id" json:"targetId"`
}
```

---

## 📚 使用示例

### 方式一：直接导入 shared 包（推荐）

```go
import "Qingyu_backend/models/shared"

type Comment struct {
    shared.IdentifiedEntity     `bson:",inline"`
    shared.BaseEntity           `bson:",inline"`
    shared.ThreadedConversation `bson:",inline"`
    shared.Likable              `bson:",inline"`

    AuthorID string `bson:"author_id" json:"authorId"`
    Content  string `bson:"content" json:"content"`
}
```

### 方式二：通过领域模块的 base 包（向后兼容）

部分模块（如 `social`、`messaging`、`writer`）提供了 `base/base.go` 文件，通过类型别名重新导出 shared 类型：

```go
// social/base.go
import (shared "Qingyu_backend/models/shared")

type BaseEntity = shared.BaseEntity
type IdentifiedEntity = shared.IdentifiedEntity
type Likable = shared.Likable
type ThreadedConversation = shared.ThreadedConversation
```

使用方式：

```go
import "Qingyu_backend/models/social"

type Comment struct {
    social.IdentifiedEntity     `bson:",inline"`
    social.BaseEntity           `bson:",inline"`
    social.ThreadedConversation `bson:",inline"`
    social.Likable              `bson:",inline"`

    AuthorID string `bson:"author_id" json:"authorId"`
    Content  string `bson:"content" json:"content"`
}
```

### 完整示例：创建评论

```go
import (
    "time"
    "Qingyu_backend/models/shared"
)

func CreateComment(authorID, content string) *Comment {
    comment := &Comment{
        AuthorID: authorID,
        Content:  content,
    }

    // 使用 shared 基础模型的方法
    comment.ID = primitive.NewObjectID().Hex()
    comment.TouchForCreate()  // 设置创建和更新时间

    return comment
}

func MarkCommentRead(comment *Comment) {
    // 如果嵌入了 shared.ReadStatus
    comment.MarkAsRead()
}

func LikeComment(comment *Comment) {
    // 如果嵌入了 shared.Likable
    comment.AddLike(1)
}
```

### 使用编辑追踪（Edited）

```go
import (
    "Qingyu_backend/models/shared"
)

type DocumentContent struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`
    shared.Edited           `bson:",inline"`

    DocumentID string `bson:"document_id" json:"documentId"`
    Content    string `bson:"content" json:"content"`
}

func SaveDocument(doc *DocumentContent, userID string) {
    // 使用 Edited 混入的方法
    doc.MarkEdited(userID)  // 更新 LastSavedAt 和 LastEditedBy
    doc.Touch()             // 更新 UpdatedAt
}
```

### 处理指针字段（ThreadedConversation）

```go
reply := &Comment{
    // ...
    ThreadedConversation: shared.ThreadedConversation{
        ReplyToCommentID: &parentCommentID,  // 使用 & 取地址
        RootID:           &rootCommentID,
    },
}

// 安全检查
if reply.ReplyToCommentID != nil {
    fmt.Println("回复评论ID:", *reply.ReplyToCommentID)
}
```

---

## 🏗️ 模块关系

### 目录结构

```
models/
├── shared/                    # 基础模型混入（本目录）
│   ├── base.go
│   ├── social.go
│   ├── communication.go
│   ├── content.go
│   └── metadata.go
├── auth/                      # 认证授权模型
├── wallet/                    # 钱包模型
├── social/                    # 社交模型
│   └── base.go               # 类型别名（向后兼容）
├── messaging/                 # 消息模型
│   └── base/                 # 类型别名（向后兼容）
├── writer/                    # 写作模型
│   └── base/                 # 类型别名（向后兼容）
└── ... (其他领域模块)
```

### 模块依赖关系

```
┌─────────────────────────────────────────┐
│          models/shared/                  │
│        (基础模型混入层)                   │
│  base.go, social.go, communication.go   │
└──────────────────┬──────────────────────┘
                   │
                   │ 通过类型别名或直接导入
                   │
    ┌──────────────┴──────────────┐
    │                             │
    ▼                             ▼
┌─────────┐                   ┌─────────┐
│ social/ │                   │messaging│
│ (社交)  │                   │ (消息)  │
└─────────┘                   └─────────┘
    │                             │
    ▼                             ▼
┌─────────────────────────────────────┐
│       service/social/                │
│      (使用 social 模型)               │
└─────────────────────────────────────┘
```

---

## 🔍 使用场景对照表

| 需求 | 推荐使用的混入 | 文件 |
|------|--------------|------|
| 需要ID字段 | `IdentifiedEntity` | base.go |
| 需要时间戳 | `BaseEntity` | base.go |
| 需要软删除 | `BaseEntity` + `SoftDelete()` | base.go |
| 需要已读状态 | `ReadStatus` | base.go |
| 需要编辑追踪 | `Edited` | base.go |
| 需要点赞数 | `Likable` | social.go |
| 需要回复嵌套 | `ThreadedConversation` | social.go |
| 需要收发件人 | `CommunicationBase` | communication.go |
| 需要标题 | `TitledEntity` | content.go |
| 需要名称 | `NamedEntity` | content.go |
| 需要描述 | `DescriptedEntity` | content.go |
| 需要置顶 | `Pinned` | metadata.go |
| 需要有效期 | `Expirable` | metadata.go |
| 需要关联对象 | `TargetEntity` | metadata.go |

---

## ⚠️ 注意事项

### 1. ID 类型

`IdentifiedEntity.ID` 是 `string` 类型，**不是** `primitive.ObjectID`：

```go
type IdentifiedEntity struct {
    ID string `bson:"_id,omitempty" json:"id"`
}

// ✅ 正确
comment.ID = primitive.NewObjectID().Hex()

// ❌ 错误
comment.ID = primitive.NewObjectID()  // 类型不匹配
```

### 2. 指针字段

`ThreadedConversation.ReplyToCommentID` 和 `RootID` 是 `*string` 类型，使用前需要 nil 检查：

```go
// ✅ 正确
if reply.ReplyToCommentID != nil {
    fmt.Println(*reply.ReplyToCommentID)
}

// ❌ 错误（可能 panic）
fmt.Println(*reply.ReplyToCommentID)
```

### 3. BSON inline 标签

嵌入时务必使用 `bson:",inline"` 标签，否则字段不会被合并到 MongoDB 文档中：

```go
// ✅ 正确
type Comment struct {
    shared.IdentifiedEntity `bson:",inline"`
}

// ❌ 错误（会创建嵌套对象）
type Comment struct {
    shared.IdentifiedEntity `bson:"identified"`
}
```

### 4. 方法接收者

所有混入的方法都是值接收者，既可以值调用也可以指针调用：

```go
comment.Touch()        // ✅ 值调用
comment.TouchForCreate() // ✅ 值调用
```

---

## 📖 相关文档

- [测试规范](../../doc/standards/testing/) - 测试层级规范
- [重构说明](./P0_REFACTOR_SUMMARY.md) - ID 类型重构总结
- [领域模型](../README.md) - 各领域模型说明

---

*共享基础模型定义完成 ✅*
