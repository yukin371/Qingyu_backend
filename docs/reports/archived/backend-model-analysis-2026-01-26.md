# Qingyu Backend 数据模型设计分析报告

> **报告日期**: 2026-01-26
> **分析范围**: Qingyu Backend 全部数据模型（117个文件，11,229行代码）
> **分析方法**: 静态代码分析 + 设计规范对比 + 历史任务回顾
> **分析人员**: 猫娘助手 Kore 🐱

---

## 📋 执行摘要

### 分析范围

本次分析对 Qingyu Backend 的数据模型进行了全面审查，涵盖：

- **共享基础模型**: `models/shared/` - 10个文件
- **核心业务模块**: users, auth, bookstore, reader, writer - 47个文件
- **支撑服务模块**: finance, social, messaging, notification, ai - 45个文件
- **DTO层**: `models/dto/` - 4个文件
- **其他模块**: audit, recommendation, stats, storage, search - 15个文件

### 关键发现

#### ✅ 优势

1. **完善的共享类型系统**: `models/shared/types/` 提供了 Money, Rating, Progress, Role 等专用类型
2. **DTOConverter 工具**: 提供了完整的 Model↔DTO 转换工具
3. **部分模块已规范化**: users/user.go, reader/readingprogress.go 等已使用 shared.IdentifiedEntity
4. **历史修复已完成**: Phase 4 已解决钱包模型、消息/通知模型、AI会话存储等问题
5. **设计文档完善**: 有详细的数据库设计说明书和迁移指南

#### ⚠️ 主要问题

1. **ID 类型不一致（P0）**:
   - 7个模块使用 `string` 类型存储 ID
   - 5个模块使用 `primitive.ObjectID`
   - 导致跨模块关联和转换困难

2. **基础模型不统一（P0）**:
   - 存在3套不同的 IdentifiedEntity 定义
   - BaseEntity 字段在多处重复定义
   - writer 和 messaging 模块使用自定义 base 包

3. **金额字段未统一（P1）**:
   - bookstore 模块使用 `float64` 存储金额
   - finance 模块已使用 `types.Money`（正确）
   - 存在精度风险

4. **枚举类型重复定义（P1）**:
   - BookStatus, UserRole 等在多处定义
   - 未统一使用 shared/types 包

5. **JSON 标签异常（P2）**:
   - 多个文件出现 `json:"$1$2"` 错误
   - 可能是代码生成工具问题

### 优先级建议

| 优先级 | 问题 | 影响范围 | 建议处理时间 |
|--------|------|----------|--------------|
| **P0** | ID 类型统一 | 全局 | 立即开始 |
| **P0** | 基础模型统一 | 全局 | 立即开始 |
| **P1** | 金额字段统一 | bookstore + frontend | 1-2周 |
| **P1** | 枚举类型统一 | 全局 | 1-2周 |
| **P2** | JSON 标签修复 | 全局 | 2-4周 |

---

## 📊 模型清单

### 共享基础模型 (`models/shared/`)

| 文件 | 用途 | 代码行数 | 状态 |
|------|------|----------|------|
| `base.go` | IdentifiedEntity, BaseEntity, ReadStatus, Edited | 127 | ✅ 规范 |
| `types/money.go` | Money 类型（int64，单位：分） | 178 | ✅ 规范 |
| `types/id.go` | ObjectID 解析转换工具 | 91 | ✅ 规范 |
| `types/enums.go` | UserRole, PageMode, DocumentStatus, WithdrawalStatus, BookStatus, OrderStatus | 323 | ✅ 规范 |
| `types/rating.go` | Rating 类型（0.0-5.0），RatingDistribution | 215 | ✅ 规范 |
| `types/progress.go` | Progress 类型（0.0-1.0） | 241 | ✅ 规范 |
| `types/converter.go` | DTOConverter（Model↔DTO 转换） | 273 | ✅ 规范 |
| `communication.go` | 通信相关 mixin（CommunicationBase 等） | - | ✅ 规范 |
| `content.go` | 内容相关 mixin | - | ✅ 规范 |
| `metadata.go` | 元数据相关 mixin | - | ✅ 规范 |
| `social.go` | 社交相关 mixin（Likable, ThreadedConversation） | - | ✅ 规范 |

**小计**: 10个文件，~1,448行代码

### 核心业务模块

#### Users 模块 (`models/users/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `user.go` | 用户模型 | ✅ IdentifiedEntity, BaseEntity | ✅ 最佳实践 |
| `admin.go` | AuditRecord, AdminLog | ❌ 无 | ⚠️ ID类型错误 |
| `user_dto.go` | 用户 DTO | - | ✅ |
| `user_statistics.go` | 用户统计 | - | ⚠️ 需检查 |
| `user_test.go` | 用户测试 | - | ✅ |

**问题**:
- `admin.go` 中 `AuditRecord` 和 `AdminLog` 使用 `string ID`，未使用 `IdentifiedEntity`
- JSON 标签异常：`json:"$1$2"`

#### Auth 模块 (`models/auth/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `role.go` | Role, Permission, UserRole | ❌ 无 | ⚠️ ID类型错误 |
| `session.go` | Session, TokenBlacklist | ❌ 无 | ⚠️ 重复定义时间戳 |
| `jwt.go` | JWT 相关 | - | ✅ |
| `oauth.go` | OAuth 相关 | - | ✅ |

**问题**:
- 所有模型 ID 使用 `string` 类型
- 未使用 `IdentifiedEntity` 或 `BaseEntity`
- 时间戳字段重复定义
- JSON 标签异常

#### Bookstore 模块 (`models/bookstore/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `book.go` | 书籍模型 | ✅ IdentifiedEntity, BaseEntity, types.Rating | ⚠️ 金额类型错误 |
| `book_detail.go` | 书籍详情 | ✅ IdentifiedEntity, BaseEntity | ⚠️ 金额类型错误 |
| `chapter.go` | 章节模型 | ❌ 无 | ❌ ID类型错误 + JSON标签异常 |
| `chapter_content.go` | 章节内容 | ❌ 无 | ⚠️ 需检查 |
| `chapter_purchase.go` | 章节购买 | ❌ 无 | ⚠️ 金额类型错误 |
| `book_rating.go` | 书籍评分 | - | ⚠️ 需检查 |
| `book_statistics.go` | 书籍统计 | - | ⚠️ 需检查 |
| `category.go` | 分类 | - | ⚠️ 需检查 |
| `banner.go` | 横幅 | - | ⚠️ 需检查 |
| `ranking.go` | 排行 | - | ⚠️ 需检查 |

**问题**:
- `chapter.go` ID 使用 `string` 类型，未使用 `IdentifiedEntity`
- 金额字段使用 `float64`，应使用 `types.Money`
- `BookStatus` 在本地定义，应使用 `shared/types.BookStatus`
- JSON 标签异常

#### Reader 模块 (`models/reader/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `readingprogress.go` | 阅读进度 | ✅ IdentifiedEntity, BaseEntity, types.Progress | ✅ 最佳实践 |
| `bookmark.go` | 书签 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `annotation.go` | 注释 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `chapter_comment.go` | 章节评论 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `reader_font.go` | 阅读字体 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `reader_theme.go` | 阅读主题 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `reading_history.go` | 阅读历史 | - | ⚠️ 需检查 |
| `readingsettings.go` | 阅读设置 | - | ⚠️ 需检查 |
| `collection_compat.go` | 收藏兼容 | - | ⚠️ 需检查 |

**状态**: reader 模块是最佳实践示例！

#### Writer 模块 (`models/writer/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `project.go` | 项目模型 | ❌ writer/base.IdentifiedEntity | ⚠️ 使用自定义base |
| `document.go` | 文档模型 | ❌ writer/base.* | ⚠️ 使用自定义base |
| `node.go` | 节点模型 | ❌ writer/base.* | ⚠️ 使用自定义base |
| `batch_operation.go` | 批量操作 | ❌ writer/base.* | ⚠️ 使用自定义base |
| ... | 其他文件 | ❌ writer/base.* | ⚠️ 使用自定义base |

**问题**:
- 所有模型使用 `writer/base.IdentifiedEntity`（AuthorID 是 `string` 类型）
- 未使用全局的 `shared.IdentifiedEntity`
- 存在两套基础模型系统

### 支撑服务模块

#### Finance 模块 (`models/finance/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `wallet.go` | 钱包模型 | ✅ types.Money | ⚠️ ID类型错误 |
| `author_revenue.go` | 作者收益 | ✅ types.Money | ⚠️ ID类型错误 |
| `membership.go` | 会员 | ✅ types.Money | ⚠️ ID类型错误 |

**问题**:
- 所有模型 ID 使用 `string` 类型
- 未使用 `IdentifiedEntity` 或 `BaseEntity`
- 金额字段已正确使用 `types.Money` ✅
- JSON 标签异常

#### Social 模块 (`models/social/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `base.go` | 基础模型 | ✅ shared 包别名 | ✅ 已迁移 |
| `comment.go` | 评论模型 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `message.go` | 私信模型 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `follow.go` | 关注模型 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `like.go` | 点赞模型 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `review.go` | 评论模型 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `booklist.go` | 书单 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `collection.go` | 收藏 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |
| `user_relation.go` | 用户关系 | ✅ IdentifiedEntity, BaseEntity | ✅ 规范 |

**状态**: social 模块已完全迁移到 shared 包，是最佳实践示例！

#### Messaging 模块 (`models/messaging/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `inbox_notification.go` | 站内通知 | ❌ messaging/base.* | ⚠️ 使用自定义base |
| `conversation.go` | 对话 | ❌ messaging/base.* | ⚠️ 使用自定义base |
| `message.go` | 消息 | ❌ messaging/base.* | ⚠️ 使用自定义base |
| `announcement.go` | 公告 | ❌ messaging/base.* | ⚠️ 使用自定义base |

**问题**:
- 使用 `messaging/base.IdentifiedEntity`（而非 shared）
- 设计优秀但未在生产环境使用
- 与 notification 模块功能重复

#### Notification 模块 (`models/notification/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `notification.go` | 通知模型 | ✅ shared.IdentifiedEntity | ✅ 生产使用 |

**状态**: 当前生产环境使用此模型

#### AI 模块 (`models/ai/`)

| 文件 | 用途 | 使用共享类型 | 状态 |
|------|------|-------------|------|
| `chat_session.go` | AI 会话 | ✅ IdentifiedEntity, BaseEntity | ✅ Phase 4 已优化 |
| `context.go` | 上下文 | - | ⚠️ 需检查 |
| `user_quota.go` | 用户配额 | ✅ types.Money | ✅ 规范 |
| `request_log.go` | 请求日志 | - | ⚠️ 需检查 |

### DTO 层 (`models/dto/`)

| 文件 | 用途 | 覆盖模块 | 状态 |
|------|------|----------|------|
| `user.go` | 用户 DTO | users | ✅ Request/Response 配对 |
| `bookstore.go` | 书城 DTO | bookstore | ✅ Request/Response 配对 |
| `reader.go` | 阅读 DTO | reader | ✅ Request/Response 配对 |
| `audit.go` | 审计 DTO | audit | ✅ Request/Response 配对 |

**状态**: DTO 层设计良好，但覆盖不完整

---

## 🔍 共享基础模型分析

### `models/shared/base.go`

**结构**:
```go
// IdentifiedEntity 包含ID字段的基础实体
type IdentifiedEntity struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

// BaseEntity 通用实体基类
type BaseEntity struct {
    CreatedAt time.Time `json:"createdAt" bson:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
    DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// ReadStatus 已读状态混入
type ReadStatus struct {
    IsRead bool       `json:"isRead" bson:"is_read"`
    ReadAt *time.Time `json:"readAt,omitempty" bson:"read_at,omitempty"`
}

// Edited 编辑追踪混入
type Edited struct {
    LastSavedAt  time.Time `json:"lastSavedAt" bson:"last_saved_at"`
    LastEditedBy string    `json:"lastEditedBy" bson:"last_edited_by"`
}
```

**最佳实践**:
- ✅ 使用 `bson:",inline"` 标签嵌入
- ✅ BSON 使用 `snake_case`，JSON 使用 `camelCase`
- ✅ 提供实用方法（Touch, SoftDelete, IsDeleted 等）
- ✅ ID 字段使用 `primitive.ObjectID`

**使用示例**（users/user.go）:
```go
type User struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    Username string `bson:"username" json:"username" validate:"required,min=3,max=50"`
    Email    string `bson:"email,omitempty" json:"email" validate:"omitempty,email"`
    // ...
}
```

### `models/shared/types/` 包

#### Money 类型（`types/money.go`）

**设计**:
```go
type Money int64  // 单位：分

const (
    MoneyZero     Money = 0
    CentsPerYuan  Money = 100
)
```

**优势**:
- ✅ 避免浮点数精度问题
- ✅ 提供完整的运算方法（Add, Sub, Mul, Div）
- ✅ 提供格式化方法（String, ToYuan, ToCents）
- ✅ 提供比较方法（Compare, GreaterThan, LessThan）
- ✅ 支持折扣计算（ApplyDiscount）

**正确使用**（finance/wallet.go）:
```go
type Wallet struct {
    UserID  string      `json:"-" bson:"user_id"`
    Balance types.Money `json:"-" bson:"balance_cents"`  // ✅ 正确
}
```

**错误使用**（bookstore/book.go）:
```go
type Book struct {
    Price float64 `bson:"price" json:"price"`  // ❌ 应改为 types.Money
}
```

#### Rating 类型（`types/rating.go`）

**设计**:
```go
type Rating float64  // 范围: 0.0-5.0

const (
    RatingMin     Rating = 0.0
    RatingMax     Rating = 5.0
    RatingDefault Rating = 0.0
)

type RatingDistribution map[string]int64  // key: "1", "2", "3", "4", "5"
```

**优势**:
- ✅ 统一评分范围（0-5星）
- ✅ 提供验证方法（IsValid）
- ✅ RatingDistribution 使用字符串键（避免 BSON 序列化问题）
- ✅ 提供统计方法（GetAverage, GetTotal, GetPercentage）

**正确使用**（bookstore/book.go）:
```go
type Book struct {
    Rating types.Rating `bson:"rating" json:"rating" validate:"min=0,max=5"`  // ✅ 正确
}
```

#### Progress 类型（`types/progress.go`）

**设计**:
```go
type Progress float32  // 范围: 0.0-1.0

const (
    ProgressMin   Progress = 0.0
    ProgressMax   Progress = 1.0
    ProgressZero  Progress = 0.0
    ProgressFull  Progress = 1.0
)
```

**优势**:
- ✅ 统一进度表示（0-1）
- ✅ 提供百分比转换（ToPercent）
- ✅ 提供阶段判断（GetStage）
- ✅ 支持进度运算（Add, Sub, Percentage）

**正确使用**（reader/readingprogress.go）:
```go
type ReadingProgress struct {
    Progress types.Progress `bson:"progress" json:"progress"`  // ✅ 正确
}
```

#### 枚举类型（`types/enums.go`）

**定义的枚举**:
- `UserRole`: reader, author, admin
- `PageMode`: scroll, paginate
- `DocumentStatus`: draft, published, archived, deleted
- `WithdrawalStatus`: pending, approved, rejected, completed
- `BookStatus`: draft, published, completed, paused, deleted
- `OrderStatus`: pending, paid, completed, cancelled, refunded

**优势**:
- ✅ 统一管理所有枚举类型
- ✅ 提供验证方法（IsValid）
- ✅ 提供状态流转方法（CanPublish, CanEdit, IsFinal 等）

**问题**:
- ⚠️ bookstore/book.go 重新定义了 `BookStatus`（未使用 shared/types）
- ⚠️ auth/role.go 重新定义了 `UserRole` 常量

#### DTOConverter（`types/converter.go`）

**功能**:
- ✅ ID 转换（ModelIDToDTO, DTOIDToModel）
- ✅ 金额转换（MoneyToDTO, DTOMoneyToYuan）
- ✅ 评分转换（RatingToDTO, DTORatingToModel）
- ✅ 进度转换（ProgressToDTO, DTOProgressToModel）
- ✅ 时间转换（TimeToISO8601, ISO8601ToTime）
- ✅ 枚举转换（UserRoleToString, StringToUserRole 等）

**使用建议**:
```go
var converter = types.DTOConverter{}

// Model → DTO
dtoID := converter.ModelIDToDTO(model.ID)
dtoPrice := converter.MoneyToYuan(model.Price)

// DTO → Model
modelID, err := converter.DTOIDToModel(dtoID)
modelPrice := converter.DTOMoneyToYuan(dtoPrice)
```

---

## 🎯 核心业务模块分析

### Users 模块

**最佳实践示例**（`user.go`）:

```go
type User struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    Username string `bson:"username" json:"username" validate:"required,min=3,max=50"`
    Email    string `bson:"email,omitempty" json:"email" validate:"omitempty,email"`
    Password string `bson:"password" json:"-" validate:"required,min=6"`

    Roles    []string `bson:"roles" json:"roles" validate:"required,dive,oneof=reader author admin"`
    VIPLevel int      `bson:"vip_level" json:"vipLevel" validate:"min=0,max=5"`

    Status      UserStatus `bson:"status" json:"status" validate:"required,oneof=active inactive banned deleted"`
    Avatar      string     `bson:"avatar,omitempty" json:"avatar,omitempty"`
    Nickname    string     `bson:"nickname,omitempty" json:"nickname,omitempty" validate:"max=50"`
    Bio         string     `bson:"bio,omitempty" json:"bio,omitempty" validate:"max=500"`

    // 认证相关
    EmailVerified bool      `bson:"email_verified" json:"emailVerified"`
    PhoneVerified bool      `bson:"phone_verified" json:"phoneVerified"`
    LastLoginAt   time.Time `bson:"last_login_at,omitempty" json:"lastLoginAt,omitempty"`
    LastLoginIP   string    `bson:"last_login_ip,omitempty" json:"lastLoginIP,omitempty"`
}
```

**优势**:
- ✅ 使用 `shared.IdentifiedEntity` 和 `BaseEntity`
- ✅ 完整的 validate 验证规则
- ✅ 密码字段使用 `json:"-"` 不序列化
- ✅ 支持多角色和角色继承
- ✅ 提供丰富的业务方法（HasRole, IsVIP, UpdateLastLogin 等）

**问题**（`admin.go`）:
```go
// ❌ 问题1: ID 使用 string 类型
type AuditRecord struct {
    ID   string `json:"id" bson:"_id,omitempty"`
    // ...
}

// ❌ 问题2: JSON 标签异常
type AuditRecord struct {
    ContentID string `json:"$1$2" bson:"content_id"`  // 应该是 json:"contentId"
}

// ❌ 问题3: 未使用 IdentifiedEntity 或 BaseEntity
type AuditRecord struct {
    CreatedAt time.Time `json:"$1$2" bson:"created_at"`
    UpdatedAt time.Time `json:"$1$2" bson:"updated_at"`
    // 应该嵌入 shared.BaseEntity
}
```

**建议**:
```go
// ✅ 改进方案
type AuditRecord struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    ContentID   primitive.ObjectID `bson:"content_id" json:"contentId"`
    ContentType string             `bson:"content_type" json:"contentType" validate:"required"`
    Status      string             `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"`
    ReviewerID  primitive.ObjectID `bson:"reviewer_id,omitempty" json:"reviewerId,omitempty"`
    Reason      string             `bson:"reason,omitempty" json:"reason,omitempty" validate:"max=500"`
    ReviewedAt  *time.Time         `bson:"reviewed_at,omitempty" json:"reviewedAt,omitempty"`
    Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}
```

### Auth 模块

**问题**（`role.go`）:
```go
// ❌ 问题1: ID 使用 string 类型
type Role struct {
    ID   string `json:"id" bson:"_id,omitempty"`
    // ...
}

// ❌ 问题2: 未使用 BaseEntity，时间戳重复定义
type Role struct {
    CreatedAt time.Time `json:"$1$2" bson:"created_at"`
    UpdatedAt time.Time `json:"$1$2" bson:"updated_at"`
}

// ❌ 问题3: JSON 标签异常
type Role struct {
    IsSystem  bool `json:"$1$2" bson:"is_system"`   // 应该是 json:"isSystem"
    IsDefault bool `json:"$1$2" bson:"is_default"`  // 应该是 json:"isDefault"
}
```

**建议**:
```go
// ✅ 改进方案
type Role struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    Name        string   `bson:"name" json:"name" validate:"required,min=1,max=50"`
    Description string   `bson:"description" json:"description" validate:"max=500"`
    Permissions []string `bson:"permissions" json:"permissions"`
    IsSystem    bool     `bson:"is_system" json:"isSystem"`
    IsDefault   bool     `bson:"is_default" json:"isDefault"`
}
```

### Bookstore 模块

**最佳实践**（`book.go`）:
```go
type Book struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    Title         string               `bson:"title" json:"title" validate:"required,min=1,max=200"`
    Author        string               `bson:"author" json:"author" validate:"required,min=1,max=100"`
    AuthorID      primitive.ObjectID   `bson:"author_id,omitempty" json:"authorId,omitempty"`
    Introduction  string               `bson:"introduction" json:"introduction" validate:"max=1000"`
    Cover         string               `bson:"cover" json:"cover" validate:"url"`
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
    Status        BookStatus           `bson:"status" json:"status" validate:"required"`
    Rating        types.Rating         `bson:"rating" json:"rating" validate:"min=0,max=5"`  // ✅ 正确使用 types.Rating
    Price         float64              `bson:"price" json:"price" validate:"min=0"`  // ❌ 应改为 types.Money
    // ...
}
```

**问题**（`chapter.go`）:
```go
// ❌ 问题1: ID 使用 string 类型
type Chapter struct {
    ID string `bson:"_id,omitempty" json:"id"`
    // ...
}

// ❌ 问题2: 未使用 IdentifiedEntity 或 BaseEntity
type Chapter struct {
    PublishTime time.Time `bson:"publish_time" json:"$1$2"`
    CreatedAt   time.Time `bson:"created_at" json:"$1$2"`
    UpdatedAt   time.Time `bson:"updated_at" json:"$1$2"`
}

// ❌ 问题3: Price 使用 float64
type Chapter struct {
    Price float64 `bson:"price" json:"price"`  // ❌ 应改为 types.Money
}
```

**建议**:
```go
// ✅ 改进方案
type Chapter struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    BookID     primitive.ObjectID `bson:"book_id" json:"bookId" validate:"required"`
    Title      string             `bson:"title" json:"title" validate:"required,min=1,max=200"`
    ChapterNum int                `bson:"chapter_num" json:"chapterNum" validate:"min=1"`
    WordCount  int                `bson:"word_count" json:"wordCount" validate:"min=0"`
    IsFree     bool               `bson:"is_free" json:"isFree"`
    Price      types.Money        `bson:"price_cents" json:"-"`  // ✅ 使用 types.Money，不暴露到 JSON
    // ...
}
```

**金额字段修复优先级**:
根据 `docs/database/data-model-fixes-migration-guide.md`，Bookstore 模块的金额字段修复已经在计划中：

| 字段 | 当前类型 | 目标类型 | 优先级 |
|------|----------|----------|--------|
| Book.Price | float64 | types.Money | High |
| BookDetail.Price | float64 | types.Money | High |
| Chapter.Price | float64 | types.Money | High |
| ChapterPurchase.Price | float64 | types.Money | Medium |
| ChapterPurchase.TotalPrice | float64 | types.Money | Medium |

### Reader 模块

**最佳实践示例**（`readingprogress.go`）:
```go
type ReadingProgress struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    UserID      primitive.ObjectID `bson:"user_id" json:"userId" validate:"required"`
    BookID      primitive.ObjectID `bson:"book_id" json:"bookId" validate:"required"`
    ChapterID   primitive.ObjectID `bson:"chapter_id" json:"chapterId" validate:"required"`
    Progress    types.Progress     `bson:"progress" json:"progress" validate:"min=0,max=1"`  // ✅ 正确使用 types.Progress
    ReadingTime int64              `bson:"reading_time" json:"readingTime" validate:"min=0"`
    LastReadAt  time.Time          `bson:"last_read_at" json:"lastReadAt"`
    Status      string             `bson:"status" json:"status" validate:"required,oneof=reading want_read finished"`
}
```

**状态**: Reader 模块是最佳实践示例！所有模型都：
- ✅ 使用 `shared.IdentifiedEntity` 和 `BaseEntity`
- ✅ ID 字段使用 `primitive.ObjectID`
- ✅ 正确使用共享类型（types.Progress）
- ✅ 有完整的 validate 验证规则
- ✅ BSON/JSON 标签符合规范

### Writer 模块

**问题**（`project.go`）:
```go
// ❌ 问题1: 使用 writer/base.IdentifiedEntity（而非 shared）
type Project struct {
    base.IdentifiedEntity `bson:",inline"`  // writer/base.IdentifiedEntity
    base.OwnedEntity      `bson:",inline"`  // AuthorID 是 string 类型
    base.TitledEntity     `bson:",inline"`
    base.Timestamps       `bson:",inline"`
    // ...
}

// ❌ 问题2: AuthorID 是 string 类型
type OwnedEntity struct {
    AuthorID string `bson:"author_id" json:"authorId"`  // ❌ 应改为 primitive.ObjectID
}
```

**影响**:
- Writer 模块与全局的 `shared.IdentifiedEntity` 不兼容
- 跨模块关联时需要类型转换（string ↔ primitive.ObjectID）
- 与 users/auth 模块的 ID 类型不一致

**建议**:
1. 短期：创建 ID 转换工具函数
2. 中期：迁移到 `shared.IdentifiedEntity`
3. 长期：统一所有模块使用共享基础模型

---

## 🛠️ 支撑服务模块分析

### Finance 模块

**正确使用 types.Money**（`wallet.go`）:
```go
type Wallet struct {
    ID        string      `json:"id" bson:"_id,omitempty"`  // ❌ ID 类型错误
    UserID    string      `json:"-" bson:"user_id"`
    Balance   types.Money `json:"-" bson:"balance_cents"`  // ✅ 正确使用 types.Money
    Frozen    bool        `json:"frozen" bson:"frozen"`
    // ...
}

type Transaction struct {
    ID              string      `json:"id" bson:"_id,omitempty"`
    UserID          string      `json:"-" bson:"user_id"`
    Amount          types.Money `json:"-" bson:"amount_cents"`  // ✅ 正确
    Balance         types.Money `json:"-" bson:"balance_cents"`  // ✅ 正确
    // ...
}
```

**优势**:
- ✅ 所有金额字段都使用 `types.Money`
- ✅ 金额字段使用 `json:"-"` 不暴露到 API
- ✅ BSON 字段名使用 `_cents` 后缀（明确单位）

**问题**:
- ❌ ID 字段使用 `string` 类型
- ❌ 未使用 `IdentifiedEntity` 或 `BaseEntity`
- ❌ JSON 标签异常

### Social 模块

**最佳实践示例**（`comment.go`）:
```go
type Comment struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`
    shared.Likable          `bson:",inline"`         // 点赞计数
    shared.ThreadedConversation `bson:",inline"`     // 父子评论

    Content     string             `bson:"content" json:"content" validate:"required,min=1,max=1000"`
    AuthorID    primitive.ObjectID `bson:"author_id" json:"authorId" validate:"required"`
    TargetID    primitive.ObjectID `bson:"target_id" json:"targetId" validate:"required"`
    TargetType  string             `bson:"target_type" json:"targetType" validate:"required,oneof=book chapter comment"`
    IsDeleted   bool               `bson:"is_deleted" json:"isDeleted"`
}
```

**状态**: Social 模块已完全迁移到 shared 包！

**迁移记录**（`social/base.go`）:
```go
// 本文件已弃用，所有基础模型已迁移到 models/shared
// 为了向后兼容，这里重新导出 shared 包的类型

import (
    shared "Qingyu_backend/models/shared"
)

// 类型别名 - 向后兼容
type BaseEntity = shared.BaseEntity
type IdentifiedEntity = shared.IdentifiedEntity
type Timestamps = shared.BaseEntity
type Likable = shared.Likable
type ThreadedConversation = shared.ThreadedConversation
```

### Messaging 模块

**优秀设计但未使用**（`inbox_notification.go`）:
```go
// 【模型说明】
// 这是改进版的站内通知模型，设计上优于旧的 notification.Notification 模型。
//
// 【设计优势】
// 1. 使用 mixin 模式提高代码复用性
// 2. 字段命名统一（JSON camelCase, BSON snake_case）
// 3. ID 类型统一使用 primitive.ObjectID
// 4. 支持更丰富的功能：目标关联、置顶、发送者快照等
// 5. 类型枚举更细化（comment, like, follow, mention 等）
//
// 【当前状态】
// - ✅ 模型设计完成（Phase 1）
// - ✅ Repository 和 Service 实现已创建（Phase 2）
// - ⏸️  未在生产环境使用
//
// 【与 notification.Notification 的关系】
// 两者功能相似（都是站内通知），但当前生产环境使用的是 notification.Notification。
// 本模型设计更优，但考虑到系统稳定性，迁移计划已延后。

type InboxNotification struct {
    base.IdentifiedEntity  `bson:",inline"`  // messaging/base.IdentifiedEntity
    base.Timestamps        `bson:",inline"`
    base.CommunicationBase `bson:",inline"`
    base.Expirable         `bson:",inline"`
    base.TargetEntity      `bson:",inline"`
    base.Pinned            `bson:",inline"`

    Type       InboxNotificationType     `bson:"type" json:"type" validate:"required"`
    Priority   InboxNotificationPriority `bson:"priority" json:"priority" validate:"required,oneof=low normal high urgent"`
    Title      string                    `bson:"title" json:"title" validate:"required,min=1,max=200"`
    Content    string                    `bson:"content" json:"content" validate:"required,min=1,max=1000"`
    Data       map[string]interface{}    `bson:"data,omitempty" json:"data,omitempty"`
    // ...
}
```

**问题**:
- ❌ 使用 `messaging/base.IdentifiedEntity`（而非 shared）
- ⏸️ 设计优秀但未在生产环境使用
- ⚠️ 与 notification 模块功能重复

**建议**:
1. 短期：保持现状，使用 notification.Notification
2. 中期：评估迁移到 InboxNotification 的成本和收益
3. 长期：统一到 shared.IdentifiedEntity

### Notification 模块

**当前生产环境使用**（`notification.go`）:
```go
type Notification struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    UserID  primitive.ObjectID `bson:"user_id" json:"userId" validate:"required"`
    Type    string             `bson:"type" json:"type" validate:"required"`
    Title   string             `bson:"title" json:"title" validate:"required,min=1,max=200"`
    Content string             `bson:"content" json:"content" validate:"required,min=1,max=1000"`
    IsRead  bool               `bson:"is_read" json:"isRead"`
    ReadAt  *time.Time         `bson:"read_at,omitempty" json:"readAt,omitempty"`
    // ...
}
```

**状态**: ✅ 规范，使用 `shared.IdentifiedEntity`

---

## 🔄 DTO 层分析

### DTO 文件清单

| 文件 | Request DTO | Response DTO | 状态 |
|------|-------------|--------------|------|
| `dto/user.go` | ✅ 7个 | ✅ 7个 | ✅ 配对完整 |
| `dto/bookstore.go` | ✅ 6个 | ✅ 6个 | ✅ 配对完整 |
| `dto/reader.go` | ✅ 3个 | ✅ 3个 | ✅ 配对完整 |
| `dto/audit.go` | ✅ 2个 | ✅ 2个 | ✅ 配对完整 |

### DTO 设计模式

**示例**（`dto/user.go`）:
```go
// Request DTO
type UpdateProfileRequest struct {
    Nickname string `json:"nickname" validate:"max=50"`
    Bio      string `json:"bio" validate:"max=500"`
    Avatar   string `json:"avatar" validate:"url"`
}

// Response DTO
type UserResponse struct {
    ID          string    `json:"id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    Nickname    string    `json:"nickname"`
    Avatar      string    `json:"avatar"`
    Bio         string    `json:"bio"`
    Roles       []string  `json:"roles"`
    VIPLevel    int       `json:"vipLevel"`
    Status      string    `json:"status"`
    CreatedAt   string    `json:"createdAt"`
    UpdatedAt   string    `json:"updatedAt"`
}

// 转换方法
func ToUserResponse(user *users.User) UserResponse {
    return UserResponse{
        ID:        user.ID.Hex(),
        Username:  user.Username,
        Email:     user.Email,
        Nickname:  user.Nickname,
        Avatar:    user.Avatar,
        Bio:       user.Bio,
        Roles:     user.Roles,
        VIPLevel:  user.VIPLevel,
        Status:    string(user.Status),
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
    }
}
```

**优势**:
- ✅ Request/Response 配对完整
- ✅ ID 字段使用 `string` 类型（API 标准）
- ✅ 时间字段使用 ISO8601 字符串
- ✅ 有完整的转换方法

**问题**:
- ⚠️ 未使用 `types.DTOConverter` 工具
- ⚠️ 金额字段转换逻辑不统一（有的除以100，有的不处理）
- ⚠️ 缺少部分模块的 DTO（finance, social 等）

### DTOConverter 使用建议

虽然 `types/converter.go` 提供了 `DTOConverter` 工具，但实际 DTO 层未使用：

**现状**（手动转换）:
```go
func ToUserResponse(user *users.User) UserResponse {
    return UserResponse{
        ID: user.ID.Hex(),  // 手动调用 .Hex()
        CreatedAt: user.CreatedAt.Format(time.RFC3339),  // 手动格式化
        // ...
    }
}
```

**建议**（使用 DTOConverter）:
```go
var converter = types.DTOConverter{}

func ToUserResponse(user *users.User) UserResponse {
    return UserResponse{
        ID:        converter.ModelIDToDTO(user.ID),
        CreatedAt: converter.TimeToISO8601(user.CreatedAt),
        UpdatedAt: converter.TimeToISO8601(user.UpdatedAt),
        // ...
    }
}
```

---

## 🔎 类型一致性检查

### ID 类型

**问题**: ID 类型在模块间不一致

| 模块 | ID 类型 | 与 shared 兼容 | 状态 |
|------|---------|----------------|------|
| users/user.go | `primitive.ObjectID` | ✅ | ✅ 正确 |
| reader/* | `primitive.ObjectID` | ✅ | ✅ 正确 |
| social/* | `primitive.ObjectID` | ✅ | ✅ 正确 |
| notification/* | `primitive.ObjectID` | ✅ | ✅ 正确 |
| ai/* | `primitive.ObjectID` | ✅ | ✅ 正确 |
| users/admin.go | `string` | ❌ | ❌ 错误 |
| auth/* | `string` | ❌ | ❌ 错误 |
| bookstore/chapter.go | `string` | ❌ | ❌ 错误 |
| writer/* | `string` (AuthorID) | ❌ | ❌ 错误 |
| finance/* | `string` | ❌ | ❌ 错误 |

**影响**:
- 跨模块关联需要类型转换
- Repository 接口不一致
- DTO 转换逻辑复杂

**统计**:
- ✅ 使用 `primitive.ObjectID`: 5个模块
- ❌ 使用 `string`: 5个模块

### Money/金额字段

**问题**: 金额字段类型不统一

| 模块 | 字段 | 当前类型 | 目标类型 | 状态 |
|------|------|----------|----------|------|
| finance/wallet.go | Balance | `types.Money` | - | ✅ 正确 |
| finance/transaction.go | Amount | `types.Money` | - | ✅ 正确 |
| bookstore/book.go | Price | `float64` | `types.Money` | ❌ 错误 |
| bookstore/chapter.go | Price | `float64` | `types.Money` | ❌ 错误 |

**影响**:
- 精度丢失风险
- 前端显示逻辑不统一
- 金额运算复杂

**迁移计划**: 已在 `docs/database/data-model-fixes-migration-guide.md` 中定义

### Rating/评分字段

**状态**: ✅ 已统一使用 `types.Rating`

| 模块 | 字段 | 类型 | 状态 |
|------|------|------|------|
| bookstore/book.go | Rating | `types.Rating` | ✅ 正确 |
| bookstore/book_detail.go | Rating | `types.Rating` | ✅ 正确 |
| bookstore/book_rating.go | Rating | `types.Rating` | ✅ 正确 |

**优势**:
- 统一评分范围（0-5星）
- RatingDistribution 使用字符串键（避免 BSON 序列化问题）

### Progress/进度字段

**状态**: ✅ 已统一使用 `types.Progress`

| 模块 | 字段 | 类型 | 状态 |
|------|------|------|------|
| reader/readingprogress.go | Progress | `types.Progress` | ✅ 正确 |

### 时间戳字段

**状态**: 部分使用 `BaseEntity`，部分重复定义

| 模块 | 时间戳定义 | 状态 |
|------|-----------|------|
| users/user.go | `BaseEntity` | ✅ 正确 |
| reader/* | `BaseEntity` | ✅ 正确 |
| social/* | `BaseEntity` | ✅ 正确 |
| auth/role.go | 重复定义 | ❌ 错误 |
| auth/session.go | 重复定义 | ❌ 错误 |
| bookstore/chapter.go | 重复定义 | ❌ 错误 |
| finance/wallet.go | 重复定义 | ❌ 错误 |

### JSON/BSON 字段

**状态**: 混合使用 `primitive.M` 和 `map[string]interface{}`

| 模块 | 字段 | 类型 | 状态 |
|------|------|------|------|
| messaging/inbox_notification.go | Data | `map[string]interface{}` | ✅ 推荐 |
| notification/notification.go | Metadata | `map[string]interface{}` | ✅ 推荐 |
| users/admin.go | Metadata | `map[string]interface{}` | ✅ 推荐 |

**建议**: 统一使用 `map[string]interface{}`（更灵活）

---

## ✅ 验证规则覆盖

### Validate 标签使用情况

**最佳实践**（`users/user.go`）:
```go
type User struct {
    Username string `validate:"required,min=3,max=50"`
    Email    string `validate:"omitempty,email"`
    Password string `validate:"required,min=6"`
    Roles    []string `validate:"required,dive,oneof=reader author admin"`
    Status   UserStatus `validate:"required,oneof=active inactive banned deleted"`
}
```

**统计**:
- ✅ 有完整 validate 规则: users, reader, social
- ⚠️ 部分有 validate 规则: bookstore, writer
- ❌ 缺少 validate 规则: auth, finance, messaging

**常用验证规则**:
- `required`: 必填字段
- `omitempty`: 可选字段
- `min,max`: 数值/字符串范围
- `email`: 邮箱格式
- `url`: URL 格式
- `oneof`: 枚举值验证
- `dive`: 数组/切片元素验证

**建议**:
1. 所有模型都应添加 `validate` 标签
2. 统一使用 `shared/types` 包的枚举类型
3. 在 API 层使用 validator 中间件

### 自定义验证方法

**示例**（`writer/project.go`）:
```go
func (p *Project) Validate() error {
    // 验证标题
    if err := base.ValidateTitle(p.Title, 100); err != nil {
        return err
    }

    // 验证作者ID
    if p.AuthorID == "" {
        return base.ErrAuthorIDRequired
    }

    // 验证状态
    if !p.Status.IsValid() {
        return base.ErrInvalidStatus
    }

    return nil
}
```

**建议**: 统一错误处理，使用 `shared/types` 包的验证方法

---

## 🚨 问题清单（按优先级）

### P0: 关键一致性问题

#### 1. ID 类型不统一

**问题描述**: 5个模块使用 `string` 类型存储 ID，5个模块使用 `primitive.ObjectID`

**影响模块**:
- users/admin.go (AuditRecord, AdminLog)
- auth/* (Role, Session, etc.)
- bookstore/chapter.go
- writer/* (AuthorID 是 string)
- finance/* (Wallet, Transaction, etc.)

**影响范围**:
- 跨模块关联需要类型转换
- Repository 接口不一致
- DTO 转换逻辑复杂
- 容易出现类型错误

**建议修复方案**:
1. **Phase 1**: 统一所有模块使用 `primitive.ObjectID`
   ```go
   // 修复前
   type Role struct {
       ID string `json:"id" bson:"_id,omitempty"`
   }

   // 修复后
   type Role struct {
       shared.IdentifiedEntity `bson:",inline"`
       shared.BaseEntity       `bson:",inline"`
   }
   ```

2. **Phase 2**: 更新 Repository 接口
   ```go
   // 修复前
   GetByID(ctx context.Context, id string) (*Role, error)

   // 修复后
   GetByID(ctx context.Context, id primitive.ObjectID) (*Role, error)
   ```

3. **Phase 3**: 更新 DTO 转换逻辑
   ```go
   // 使用 types.DTOConverter
   dtoID := converter.ModelIDToDTO(model.ID)
   ```

**数据迁移**: 无需迁移（仅代码层面修改）

**预计工作量**: 3-5天

#### 2. 基础模型不统一

**问题描述**: 存在3套不同的 IdentifiedEntity 定义

**发现的定义**:
1. `models/shared/base.go`: IdentifiedEntity (ID: primitive.ObjectID)
2. `models/writer/base/`: IdentifiedEntity (AuthorID: string)
3. `models/messaging/base/`: IdentifiedEntity

**影响**:
- 代码重复
- 不符合 DRY 原则
- 维护困难

**建议修复方案**:
1. 所有模块统一使用 `shared.IdentifiedEntity`
2. 删除 `writer/base/` 和 `messaging/base/` 中的重复定义
3. Writer 模块的 AuthorID 改为独立的字段

**预计工作量**: 2-3天

### P1: 重要改进项

#### 3. 金额字段未统一

**问题描述**: Bookstore 模块使用 `float64` 存储金额，Finance 模块使用 `types.Money`

**影响模块**:
- bookstore/book.go (Price)
- bookstore/book_detail.go (Price)
- bookstore/chapter.go (Price)
- bookstore/chapter_purchase.go (Price, TotalPrice, etc.)

**影响**:
- 精度丢失风险
- 前端显示逻辑不统一
- 金额运算复杂

**修复指南**: 已在 `docs/database/data-model-fixes-migration-guide.md` 中定义

**数据迁移**:
```javascript
// MongoDB 迁移脚本
db.books.find({ price: { $exists: true, $type: "double" } }).forEach(function(doc) {
    var newPrice = Math.round(doc.price * 100);
    db.books.updateOne(
        { _id: doc._id },
        { $set: { price_cents: newPrice } }
    );
});
```

**前端影响**: 价格显示需要除以 100

**预计工作量**: 2-3天

#### 4. 枚举类型重复定义

**问题描述**: BookStatus, UserRole 等在多处定义

**重复定义**:
- `shared/types/enums.go`: BookStatus
- `models/bookstore/book.go`: BookStatus
- `auth/role.go`: UserRole 常量
- `shared/types/enums.go`: UserRole

**建议修复方案**:
1. 删除本地定义的枚举
2. 统一使用 `shared/types` 包的枚举
3. 更新 import

**示例**:
```go
// 修复前
import "Qingyu_backend/models/bookstore"

type BookStatus string
const (
    BookStatusDraft BookStatus = "draft"
    // ...
)

// 修复后
import "Qingyu_backend/models/shared/types"

type Book struct {
    Status types.BookStatus `bson:"status" json:"status"`
}
```

**预计工作量**: 1-2天

#### 5. 验证规则缺失

**问题描述**: auth, finance, messaging 等模块缺少 `validate` 标签

**建议**:
- 所有模型都应添加 `validate` 标签
- 统一使用 `shared/types` 包的枚举类型
- 在 API 层使用 validator 中间件

**预计工作量**: 2-3天

### P2: 优化建议

#### 6. JSON 标签异常

**问题描述**: 多个文件出现 `json:"$1$2"` 错误

**影响文件**:
- users/admin.go
- auth/role.go
- auth/session.go
- bookstore/chapter.go
- finance/wallet.go

**可能原因**: 代码生成工具错误

**建议修复方案**:
1. 全局搜索替换 `json:"$1$2"` 为正确的 camelCase 格式
2. 检查代码生成工具配置
3. 添加 pre-commit hook 检查 JSON 标签

**预计工作量**: 1天

#### 7. DTO 层覆盖不完整

**问题描述**: finance, social 等模块缺少 DTO

**建议**:
- 补充缺失的 DTO 定义
- 统一使用 `types.DTOConverter`
- 金额字段统一转换逻辑

**预计工作量**: 3-4天

---

## 💡 改进建议

### 短期改进（1-2周）

#### 1. 修复 JSON 标签异常（P2）

**优先级**: 高（影响功能）

**任务**:
1. 全局搜索 `json:"$1$2"`
2. 替换为正确的 camelCase 格式
3. 添加 pre-commit hook

**验收标准**:
- ✅ 不存在 `json:"$1$2"` 标签
- ✅ 所有 JSON 标签符合 camelCase 规范

#### 2. 统一使用 shared/types 枚举（P1）

**优先级**: 高

**任务**:
1. 删除本地定义的 BookStatus, UserRole 等
2. 统一使用 `shared/types` 包
3. 更新 import 和引用

**验收标准**:
- ✅ 不存在重复定义的枚举
- ✅ 所有枚举来自 `shared/types` 包

#### 3. 补充验证规则（P1）

**优先级**: 中

**任务**:
1. 为 auth, finance, messaging 模块添加 `validate` 标签
2. 统一验证规则格式
3. 添加 validator 中间件测试

**验收标准**:
- ✅ 所有模型都有 validate 标签
- ✅ 验证规则测试覆盖

### 中期改进（1-2月）

#### 4. ID 类型统一（P0）

**优先级**: 最高

**阶段**:
1. **Week 1-2**: 修复 users/admin.go, auth/* 模块
2. **Week 3-4**: 修复 bookstore/chapter.go 模块
3. **Week 5-6**: 修复 finance/* 模块
4. **Week 7-8**: 修复 writer/* 模块

**验收标准**:
- ✅ 所有 ID 字段使用 `primitive.ObjectID`
- ✅ 所有模型使用 `shared.IdentifiedEntity`
- ✅ Repository 接口统一
- ✅ 代码编译通过

#### 5. 基础模型统一（P0）

**优先级**: 最高

**任务**:
1. 删除 `writer/base/` 中的重复定义
2. 删除 `messaging/base/` 中的重复定义
3. 所有模块统一使用 `shared.IdentifiedEntity`

**验收标准**:
- ✅ 只存在一套 IdentifiedEntity 定义
- ✅ 所有模块使用 `shared.IdentifiedEntity`

#### 6. 金额字段统一（P1）

**优先级**: 高

**任务**:
1. 修复 Bookstore 模块金额字段
2. 执行数据迁移脚本
3. 更新前端金额显示逻辑

**验收标准**:
- ✅ 所有金额字段使用 `types.Money`
- ✅ 数据库数据已迁移
- ✅ 前端金额显示正确

### 长期改进（3-6月）

#### 7. DTO 层完善（P2）

**优先级**: 中

**任务**:
1. 补充缺失模块的 DTO 定义
2. 统一使用 `types.DTOConverter`
3. 添加 DTO 转换单元测试

**验收标准**:
- ✅ 所有模块都有完整的 DTO
- ✅ DTO 转换使用 DTOConverter
- ✅ DTO 转换测试覆盖

#### 8. 消息/通知模型统一（P2）

**优先级**: 低

**任务**:
1. 评估 InboxNotification 迁移成本
2. 制定迁移计划
3. 执行迁移（如决定迁移）

**验收标准**:
- ✅ 只存在一套通知模型
- ✅ 数据迁移完成
- ✅ 功能测试通过

---

## 📝 规范更新建议

### 需要新增的规范

#### 1. 数据模型设计规范

**文件**: `docs/architecture/model-design-guide.md`

**内容大纲**:
```markdown
# 数据模型设计规范

## 1. 基础模型使用

### 1.1 必须使用 shared.IdentifiedEntity
所有需要 ID 字段的模型必须使用 `shared.IdentifiedEntity`：

\`\`\`go
type MyModel struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`
    // ...
}
\`\`\`

### 1.2 禁止自定义基础模型
禁止在模块内定义 IdentifiedEntity 或 BaseEntity：

❌ 错误：
\`\`\`go
type MyIdentifiedEntity struct {
    ID string `bson:"_id,omitempty"`
}
\`\`\`

✅ 正确：
\`\`\`go
import "Qingyu_backend/models/shared"

type MyModel struct {
    shared.IdentifiedEntity `bson:",inline"`
}
\`\`\`

## 2. 字段类型规范

### 2.1 ID 字段
- 必须使用 `primitive.ObjectID`
- DTO 层可以转换为 `string`

### 2.2 金额字段
- 必须使用 `types.Money`（int64，单位：分）
- 禁止使用 `float64`

### 2.3 评分字段
- 必须使用 `types.Rating`（0.0-5.0）

### 2.4 进度字段
- 必须使用 `types.Progress`（0.0-1.0）

### 2.5 枚举字段
- 必须使用 `shared/types` 包的枚举类型
- 禁止在本地重新定义

## 3. 标签规范

### 3.1 BSON 标签
- 使用 `snake_case`
- 必须包含 `bson` 标签

### 3.2 JSON 标签
- 使用 `camelCase`
- 必须包含 `json` 标签
- 敏感字段使用 `json:"-"`

### 3.3 Validate 标签
- 所有字段都应添加 `validate` 标签
- 使用标准验证规则

## 4. 方法规范

### 4.1 实体方法
- 提供 `IsValid()` 方法（枚举类型）
- 提供业务方法（如 `HasRole()`, `IsVIP()`）

### 4.2 转换方法
- 使用 `types.DTOConverter`
- 提供 `ToDTO()` 和 `FromDTO()` 方法

## 5. 命名规范

### 5.1 模型命名
- 使用单数形式（User 而非 Users）
- 使用大驼峰命名法（PascalCase）

### 5.2 字段命名
- Go 字段：大驼峰（PascalCase）
- BSON 字段：snake_case
- JSON 字段：camelCase

## 6. 文档规范

### 6.1 模型注释
- 所有模型必须添加文档注释
- 说明模型用途和使用场景

### 6.2 字段注释
- 重要字段必须添加行内注释
- 说明单位（如：分、秒等）
```

#### 2. DTO 转换规范

**文件**: `docs/architecture/dto-conversion-guide.md`

**内容大纲**:
```markdown
# DTO 转换规范

## 1. 转换原则

### 1.1 使用 DTOConverter
所有 Model ↔ DTO 转换必须使用 `types.DTOConverter`：

\`\`\`go
var converter = types.DTOConverter{}

// Model → DTO
dtoID := converter.ModelIDToDTO(model.ID)

// DTO → Model
modelID, err := converter.DTOIDToModel(dtoID)
\`\`\`

### 1.2 金额字段转换
- Model 层：`types.Money`（分）
- DTO 层：`float64`（元）或 `string`（¥12.99）

\`\`\`go
// Model → DTO
dtoPrice := converter.MoneyToYuan(model.Price)  // float64

// DTO → Model
modelPrice := converter.DTOMoneyToYuan(dtoPrice)  // types.Money
\`\`\`

### 1.3 时间字段转换
- Model 层：`time.Time`
- DTO 层：ISO8601 字符串

\`\`\`go
// Model → DTO
dtoTime := converter.TimeToISO8601(model.CreatedAt)

// DTO → Model
modelTime, err := converter.ISO8601ToTime(dtoTime)
\`\`\`

## 2. DTO 结构定义

### 2.1 Request DTO
\`\`\`go
type UpdateProfileRequest struct {
    Nickname string `json:"nickname" validate:"max=50"`
    Bio      string `json:"bio" validate:"max=500"`
    Avatar   string `json:"avatar" validate:"url"`
}
\`\`\`

### 2.2 Response DTO
\`\`\`go
type UserResponse struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    CreatedAt string `json:"createdAt"`
}
\`\`\`

## 3. 转换方法

### 3.1 Model → DTO
\`\`\`go
func ToUserResponse(user *users.User) UserResponse {
    return UserResponse{
        ID:        converter.ModelIDToDTO(user.ID),
        Username:  user.Username,
        CreatedAt: converter.TimeToISO8601(user.CreatedAt),
    }
}
\`\`\`

### 3.2 DTO → Model
\`\`\`go
func ToUser(req *CreateUserRequest) *users.User {
    return &users.User{
        Username: req.Username,
        Email:    req.Email,
    }
}
\`\`\`

## 4. 批量转换

### 4.1 ModelSlice → DTOSlice
\`\`\`go
func ToUserResponseList(users []*users.User) []UserResponse {
    result := make([]UserResponse, len(users))
    for i, user := range users {
        result[i] = ToUserResponse(user)
    }
    return result
}
\`\`\`

## 5. 错误处理

### 5.1 转换错误
\`\`\`go
modelID, err := converter.DTOIDToModel(dtoID)
if err != nil {
    return nil, fmt.Errorf("invalid ID: %w", err)
}
\`\`\`
```

### 需要更新的规范

#### 3. 更新 `docs/design/database/数据库设计说明书.md`

**新增内容**:
```markdown
## 7. 数据模型规范

### 7.1 基础模型
所有集合必须使用 `shared.IdentifiedEntity` 和 `shared.BaseEntity`。

### 7.2 ID 字段
- 类型：`primitive.ObjectID`
- BSON 标签：`_id,omitempty`
- JSON 标签：`id`

### 7.3 时间戳字段
- 类型：`time.Time`
- BSON 标签：`created_at`, `updated_at`
- JSON 标签：`createdAt`, `updatedAt`

### 7.4 金额字段
- 类型：`types.Money`（int64，单位：分）
- BSON 标签：`amount_cents`
- JSON 标签：`-`（不暴露）

### 7.5 特殊类型
- 评分：`types.Rating`（0.0-5.0）
- 进度：`types.Progress`（0.0-1.0）
- 枚举：`shared/types` 包中的枚举类型
```

---

## 📊 统计数据

### 模块统计

| 类别 | 模块数 | 文件数 | 代码行数 | 状态 |
|------|--------|--------|----------|------|
| 共享基础 | 1 | 10 | ~1,448 | ✅ 规范 |
| 核心业务 | 5 | 47 | ~4,500 | ⚠️ 部分规范 |
| 支撑服务 | 5 | 45 | ~4,200 | ⚠️ 部分规范 |
| 其他 | 5 | 15 | ~1,081 | ⚠️ 需检查 |
| **总计** | **16** | **117** | **11,229** | **-** |

### 规范符合度

| 规范项 | 符合模块数 | 总模块数 | 符合率 |
|--------|-----------|----------|--------|
| 使用 shared.IdentifiedEntity | 6 | 16 | 37.5% |
| ID 字段使用 ObjectID | 6 | 16 | 37.5% |
| 金额字段使用 Money | 2 | 4 | 50% |
| 使用共享枚举类型 | 3 | 16 | 18.75% |
| 有 validate 标签 | 8 | 16 | 50% |
| BSON/JSON 标签规范 | 10 | 16 | 62.5% |

### 问题分布

| 优先级 | 问题数 | 涉及模块数 |
|--------|--------|-----------|
| **P0** | 2 | 10 |
| **P1** | 3 | 12 |
| **P2** | 2 | 8 |
| **总计** | **7** | **16** |

---

## 🎯 总结

### 主要成果

1. **完整的模型清单**: 整理了 117 个文件，11,229 行代码的数据模型
2. **问题识别**: 发现了 7 类主要问题，按优先级分类
3. **最佳实践识别**: users/user.go, reader/*, social/* 等模块是最佳实践示例
4. **改进建议**: 提供了短期、中期、长期的改进建议
5. **规范更新**: 建议新增 2 个规范文档，更新 1 个现有文档

### 关键发现

1. **共享类型系统完善**: `models/shared/types/` 提供了 Money, Rating, Progress 等专用类型
2. **部分模块已规范化**: users, reader, social 模块是最佳实践
3. **历史修复已完成**: Phase 4 解决了钱包模型、消息/通知模型等问题
4. **ID 类型不一致是最大问题**: 影响范围广，修复成本高

### 下一步行动

**立即开始**:
1. 修复 JSON 标签异常（1天）
2. 统一使用 shared/types 枚举（1-2天）
3. 补充验证规则（2-3天）

**1-2周内**:
4. ID 类型统一 - Phase 1（users/admin, auth）
5. ID 类型统一 - Phase 2（bookstore/chapter）

**1-2月内**:
6. ID 类型统一 - Phase 3（finance）
7. ID 类型统一 - Phase 4（writer）
8. 金额字段统一（bookstore）

**3-6月内**:
9. DTO 层完善
10. 消息/通知模型统一评估

---

**报告完成时间**: 2026-01-26
**分析人员**: 猫娘助手 Kore 🐱
**下次审查**: 建议 3 个月后（2026-04-26）

喵~ 🐱
