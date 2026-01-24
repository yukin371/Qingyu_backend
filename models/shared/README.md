# 共享基础模型

> 可复用的基础数据模型和类型系统，为项目提供统一的类型定义和转换工具

---

## 概述

`models/shared/` 目录包含可复用的基础模型组件和类型系统：

1. **基础模型混入（Mixins）** - 通过组合模式为各业务领域模型提供通用功能
2. **统一类型系统（types/）** - 提供跨模块的统一类型定义和转换工具

### 设计原则

1. **组合优于继承** - 通过 `bson:",inline"` 将基础模型嵌入到领域模型中
2. **单一职责** - 每个混入文件专注于一类通用功能
3. **类型安全** - 提供方法封装，确保数据一致性
4. **分层架构** - Model 层使用 ObjectID，API/DTO 层使用 string

---

## 📁 目录结构

```
models/shared/
├── base.go           # 核心基础模型混入
├── social.go         # 社交功能混入
├── communication.go  # 通信功能混入
├── content.go        # 内容相关混入
├── metadata.go       # 元数据混入
├── types/            # 统一类型系统
│   ├── id.go         # ID 类型转换
│   ├── money.go      # 金额类型
│   ├── rating.go     # 评分类型
│   ├── progress.go   # 进度类型
│   ├── enums.go      # 枚举类型
│   ├── converter.go  # DTO 转换辅助
│   └── README.md     # 类型系统文档
└── README.md         # 本文档
```

---

## 🏗️ 分层架构（方案B）

```
┌─────────────────────────────────────────────┐
│  API 层 (DTO)         → string id           │  ← 对外接口，JSON友好
├─────────────────────────────────────────────┤
│  Service 层          → 转换逻辑              │  ← Model↔DTO转换
├─────────────────────────────────────────────┤
│  Model 层            → ObjectID             │  ← 数据库存储，高效
└─────────────────────────────────────────────┘
```

### ID 类型说明

| 层级 | ID 类型 | JSON 标签 | BSON 标签 |
|------|---------|-----------|-----------|
| Model 层 | `primitive.ObjectID` | `-` (不暴露) | `_id` |
| DTO 层 | `string` | `id` | N/A |
| Service 层 | `string` | N/A | N/A |

---

## 📋 基础模型列表

### 1. base.go - 核心基础模型

**文件**: `base.go`

**核心模型**:
- `BaseEntity` - 通用实体基类（时间戳）
- `IdentifiedEntity` - ID 字段实体（Model 层）
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

// IdentifiedEntity - ID 字段（Model 层使用）
type IdentifiedEntity struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`
}

// 方法:
// - GetID() primitive.ObjectID - 获取ID
// - SetID(id primitive.ObjectID) - 设置ID

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

## 🧬 types/ 统一类型系统

**文件**: `types/README.md`

**包含**:
- `id.go` - ID 转换工具（ObjectID ↔ string）
- `money.go` - 金额类型（int64，最小单位：分）
- `rating.go` - 评分类型（0-5）
- `progress.go` - 进度类型（0-1）
- `enums.go` - 枚举类型（角色、状态等）
- `converter.go` - DTO 转换辅助

**使用示例**:
```go
import "Qingyu_backend/models/shared/types"

// Model → DTO 转换
var converter types.DTOConverter

dto.ID = converter.ModelIDToDTO(model.ID)              // ObjectID → string
dto.CreatedAt = converter.TimeToISO8601(model.CreatedAt) // time.Time → string

// DTO → Model 转换
id, err := converter.DTOIDToModel(dto.ID)               // string → ObjectID
createdAt, err := converter.ISO8601ToTime(dto.CreatedAt) // string → time.Time
```

详见 [types/README.md](./types/README.md)

---

## 📚 使用示例

### Model 层使用 shared.IdentifiedEntity

```go
import "Qingyu_backend/models/shared"

type User struct {
    shared.IdentifiedEntity `bson:",inline"`  // ID = primitive.ObjectID
    shared.BaseEntity       `bson:",inline"`  // 时间戳

    Username string `bson:"username" json:"username"`
    Email    string `bson:"email" json:"email"`
}
```

### DTO 层使用 string ID

```go
// api/v1/shared/user_types.go
type UserDTO struct {
    ID        string `json:"id"`                     // string ID
    CreatedAt string `json:"createdAt"`              // ISO8601 时间字符串
    UpdatedAt string `json:"updatedAt"`
    Username  string `json:"username"`
    Email     string `json:"email"`
}
```

### Service 层转换

```go
import (
    "Qingyu_backend/models/shared"
    "Qingyu_backend/models/shared/types"
)

// Model → DTO
func ToUserDTO(user *User) *UserDTO {
    var converter types.DTOConverter
    return &UserDTO{
        ID:        converter.ModelIDToDTO(user.ID),
        CreatedAt: converter.TimeToISO8601(user.CreatedAt),
        UpdatedAt: converter.TimeToISO8601(user.UpdatedAt),
        Username:  user.Username,
        Email:     user.Email,
    }
}

// DTO → Model
func ToUser(dto *UserDTO) (*User, error) {
    var converter types.DTOConverter
    id, createdAt, updatedAt, err := converter.ParseBaseFields(
        dto.ID, dto.CreatedAt, dto.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &User{
        IdentifiedEntity: shared.IdentifiedEntity{ID: id},
        BaseEntity:       shared.BaseEntity{CreatedAt: createdAt, UpdatedAt: updatedAt},
        Username:          dto.Username,
        Email:            dto.Email,
    }, nil
}
```

---

## ⚠️ 重要注意事项

### 1. ID 类型区分

**Model 层**:
```go
type User struct {
    shared.IdentifiedEntity `bson:",inline"`  // ID 是 primitive.ObjectID
}

// ✅ 正确
user.ID = primitive.NewObjectID()

// ❌ 错误
user.ID = "abc123"  // 类型不匹配
```

**DTO 层**:
```go
type UserDTO struct {
    ID string `json:"id"`  // ID 是 string
}
```

### 2. JSON 序列化

`IdentifiedEntity` 的 JSON 标签是 `json:"-"`，不会序列化到 JSON：

```go
user := &User{
    IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
    Username: "test",
}

// JSON 序列化
// {"username": "test"}  ← ID 不会出现
```

如果需要返回 ID 给前端，使用 DTO：
```go
dto := &UserDTO{
    ID: user.ID.Hex(),  // 手动转换为 string
    Username: user.Username,
}

// JSON 序列化
// {"id": "507f1f77bcf86cd799439011", "username": "test"}
```

### 3. BSON inline 标签

嵌入时务必使用 `bson:",inline"` 标签：

```go
// ✅ 正确
type User struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`
}

// ❌ 错误（会创建嵌套对象）
type User struct {
    shared.IdentifiedEntity `bson:"identified"`
}
```

---

## 🔍 使用场景对照表

| 需求 | 推荐使用的混入 | 文件 |
|------|--------------|------|
| 需要ID字段（Model层） | `IdentifiedEntity` | base.go |
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

## 📖 相关文档

- [类型系统文档](./types/README.md) - types/ 包详细说明
- [模型一致性修复指南](../../docs/architecture/model-consistency-fix-guide.md) - 模型重构指南
- [分层架构重构计划](../../.serena/memories/backend-layered-architecture-plan.md) - 方案B实施计划

---

## 📝 更新历史

- **2026-01-23**: 清理冗余文件（删除 json_bson.go），完善 types/converter.go，更新分层架构说明
- **2026-01-22**: 初始版本，包含基础模型混入

---

*共享基础模型定义完成 ✅*
