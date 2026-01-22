# 统一类型基建规划

> **目标**：为 Qingyu Backend 建立统一的公共类型基础设施，消除跨模块类型不一致问题。

## 目录结构

```
Qingyu_backend/models/shared/
├── types/
│   ├── id.go         # ID 类型转换与校验
│   ├── money.go      # 金额类型（最小货币单位）
│   ├── rating.go     # 评分类型（0-5）
│   ├── progress.go   # 进度类型（0-1）
│   ├── enums.go      # 枚举类型
│   ├── json_bson.go  # JSON/BSON 命名约束
│   └── converter.go  # 通用转换辅助
├── README.md         # 本文档
└── ...existing files
```

---

## 1. ID 类型：`id.go`

### 路径
`Qingyu_backend/models/shared/types/id.go`

### 用途
- 统一 `primitive.ObjectID` ↔ `string` 转换
- ID 格式校验
- 批量 ID 处理
- BSON/JSON 边界规则定义

### 核心功能
```go
package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// ParseObjectID 将 hex 字符串解析为 ObjectID
func ParseObjectID(s string) (primitive.ObjectID, error)

// ToHex 将 ObjectID 转换为 hex 字符串
func ToHex(id primitive.ObjectID) string

// IsValidObjectID 检查字符串是否为有效的 ObjectID hex 格式
func IsValidObjectID(s string) bool

// ParseObjectIDSlice 批量解析 ID 字符串
func ParseObjectIDSlice(ss []string) ([]primitive.ObjectID, map[int]error)

// ToHexSlice 批量转换 ObjectID 为 hex 字符串
func ToHexSlice(ids []primitive.ObjectID) []string

// GenerateNewObjectID 生成新的 ObjectID 并返回 hex 字符串
func GenerateNewObjectID() string
```

### 错误处理
```go
var (
    ErrInvalidIDFormat = errors.New("invalid ID format: must be 24-character hex")
    ErrEmptyID         = errors.New("ID cannot be empty")
)
```

### 使用示例
```go
// Service 层使用
func (s *service) GetBook(ctx context.Context, bookID string) (*dto.BookDTO, error) {
    // string → ObjectID
    oid, err := types.ParseObjectID(bookID)
    if err != nil {
        return nil, err
    }

    book, err := s.repo.FindByID(ctx, oid)
    if err != nil {
        return nil, err
    }

    // ObjectID → string (DTO)
    return &dto.BookDTO{
        ID: types.ToHex(book.ID),
        // ...
    }, nil
}
```

---

## 2. 金额类型：`money.go`

### 路径
`Qingyu_backend/models/shared/types/money.go`

### 用途
- 统一"最小货币单位 int64（分）"类型
- 提供格式化/展示转换
- 金额计算与校验
- 避免浮点误差

### 核心功能
```go
package types

// Money 金额类型（最小货币单位：分）
type Money int64

const (
    // MoneyZero 零金额
    MoneyZero Money = 0
    
    // CentsPerYuan 每元对应的分数
    CentsPerYuan Money = 100
)

// NewMoneyFromYuan 从元创建金额（浮点）
func NewMoneyFromYuan(yuan float64) Money

// NewMoneyFromCents 从分创建金额
func NewMoneyFromCents(cents int64) Money

// ToYuan 转换为元（浮点，仅用于展示）
func (m Money) ToYuan() float64

// ToCents 转换为分（int64）
func (m Money) ToCents() int64

// String 格式化为货币字符串（如 "¥12.99"）
func (m Money) String() string

// Add 金额相加
func (m Money) Add(other Money) Money

// Sub 金额相减
func (m Money) Sub(other Money) Money

// Mul 金额乘法（乘以系数）
func (m Money) Mul(factor float64) Money

// Div 金额除法（除以系数）
func (m Money) Div(divisor float64) Money

// IsZero 是否为零
func (m Money) IsZero() bool

// IsNegative 是否为负
func (m Money) IsNegative() bool

// Compare 比较（-1: <, 0: =, 1: >）
func (m Money) Compare(other Money) int
```

### 使用示例
```go
// Model 层定义
type Book struct {
    ID       primitive.ObjectID `bson:"_id"`
    Price    types.Money        `bson:"price_cents"`    // 价格（分）
    Discount types.Money        `bson:"discount_cents"` // 折扣（分）
}

// Service 层使用
func (s *service) CalculateDiscountedPrice(book *Book) types.Money {
    return book.Price.Sub(book.Discount)
}

// API 层展示
type BookDTO struct {
    Price    string `json:"price"`    // "¥29.99"
    Discount string `json:"discount"` // "¥5.00"
}

// 填充 DTO
dto.Price = book.Price.String()
```

---

## 3. 评分类型：`rating.go`

### 路径
`Qingyu_backend/models/shared/types/rating.go`

### 用途
- 统一 0-5 评分类型
- 评分范围校验
- 评分分布结构（map[string]int64）
- 平均分计算

### 核心功能
```go
package types

// Rating 评分类型（0.0-5.0）
type Rating float32

const (
    // RatingMin 最小评分
    RatingMin Rating = 0.0
    // RatingMax 最大评分
    RatingMax Rating = 5.0
    // RatingDefault 默认评分
    RatingDefault Rating = 0.0
)

// NewRating 创建评分
func NewRating(value float32) (Rating, error)

// MustRating 创建评分（panic on invalid）
func MustRating(value float32) Rating

// IsValid 检查评分是否有效
func (r Rating) IsValid() bool

// ToFloat 转换为 float32
func (r Rating) ToFloat() float32

// String 格式化为字符串（保留 1 位小数）
func (r Rating) String() string

// RatingDistribution 评分分布
type RatingDistribution map[string]int64 // key: "1", "2", "3", "4", "5"

// NewRatingDistribution 创建空分布
func NewRatingDistribution() RatingDistribution

// Add 添加评分
func (rd RatingDistribution) Add(rating Rating) error

// GetCount 获取某分数的个数
func (rd RatingDistribution) GetCount(star int) int64

// GetTotal 获取总评分数
func (rd RatingDistribution) GetTotal() int64

// GetAverage 计算平均分
func (rd RatingDistribution) GetAverage() Rating

// ToBSON 转换为 BSON 兼容格式
func (rd RatingDistribution) ToBSON() map[string]int64

// FromBSON 从 BSON 创建分布
func FromBSON(data map[string]int64) RatingDistribution
```

### 使用示例
```go
// Model 层定义
type BookStatistics struct {
    BookID             primitive.ObjectID    `bson:"_id"`
    AverageRating      types.Rating          `bson:"average_rating"`
    RatingDistribution types.RatingDistribution `bson:"rating_distribution"`
    TotalRatings       int64                 `bson:"total_ratings"`
}

// 添加评分
func (s *stats) AddRating(bookID primitive.ObjectID, rating types.Rating) error {
    stats, _ := s.repo.FindByID(bookID)
    
    // 更新分布
    stats.RatingDistribution.Add(rating)
    stats.TotalRatings++
    
    // 重新计算平均分
    stats.AverageRating = stats.RatingDistribution.GetAverage()
    
    return s.repo.Update(stats)
}
```

---

## 4. 进度类型：`progress.go`

### 路径
`Qingyu_backend/models/shared/types/progress.go`

### 用途
- 统一进度口径（内部 0-1）
- 提供 0-100 互转
- 进度计算与校验

### 核心功能
```go
package types

// Progress 进度类型（0.0-1.0）
type Progress float32

const (
    // ProgressMin 最小进度
    ProgressMin Progress = 0.0
    // ProgressMax 最大进度
    ProgressMax Progress = 1.0
    // ProgressZero 零进度
    ProgressZero Progress = 0.0
    // ProgressFull 完整进度
    ProgressFull Progress = 1.0
)

// NewProgress 创建进度（0-1）
func NewProgress(value float32) (Progress, error)

// NewProgressFromPercent 从百分比创建进度（0-100）
func NewProgressFromPercent(percent int) (Progress, error)

// MustProgress 创建进度（panic on invalid）
func MustProgress(value float32) Progress

// IsValid 检查进度是否有效
func (p Progress) IsValid() bool

// ToPercent 转换为百分比（0-100）
func (p Progress) ToPercent() int

// ToFloat 转换为 float32
func (p Progress) ToFloat() float32

// String 格式化为百分比字符串（如 "75%"）
func (p Progress) String() string

// IsComplete 是否完成
func (p Progress) IsComplete() bool

// IsStarted 是否已开始
func (p Progress) IsStarted() bool

// Add 累加进度
func (p Progress) Add(other Progress) Progress

// Percentage 计算占比
func (p Progress) Percentage(other Progress) int
```

### 使用示例
```go
// Model 层定义
type ReadingProgress struct {
    ID       primitive.ObjectID `bson:"_id"`
    UserID   primitive.ObjectID `bson:"user_id"`
    BookID   primitive.ObjectID `bson:"book_id"`
    Progress types.Progress     `bson:"progress"`     // 0-1
    // ...
}

// Service 层使用
func (s *service) UpdateProgress(userID, bookID primitive.ObjectID, readPages int, totalPages int) error {
    progress := types.Progress(float32(readPages) / float32(totalPages))
    
    return s.repo.Update(&ReadingProgress{
        UserID:   userID,
        BookID:   bookID,
        Progress: progress,
    })
}

// API 层展示
type ProgressDTO struct {
    Progress int    `json:"progress"` // 0-100
    Label    string `json:"label"`    // "75%"
}

// 填充 DTO
dto.Progress = progress.ToPercent()
dto.Label = progress.String()
```

---

## 5. 枚举类型：`enums.go`

### 路径
`Qingyu_backend/models/shared/types/enums.go`

### 用途
- 集中角色、PageMode、写作状态等枚举
- 提供枚举校验
- 避免枚举值散落各处

### 核心功能
```go
package types

// UserRole 用户角色
type UserRole string

const (
    RoleReader UserRole = "reader"
    RoleAuthor UserRole = "author"
    RoleAdmin  UserRole = "admin"
)

// IsValid 检查角色是否有效
func (r UserRole) IsValid() bool

// String 转换为字符串
func (r UserRole) String() string

// ParseUserRole 从字符串解析角色
func ParseUserRole(s string) (UserRole, error)

// PageMode 阅读翻页模式
type PageMode string

const (
    PageModeScroll  PageMode = "scroll"
    PageModePaginate PageMode = "paginate"
)

// IsValid 检查模式是否有效
func (m PageMode) IsValid() bool

// String 转换为字符串
func (m PageMode) String() string

// ParsePageMode 从字符串解析模式
func ParsePageMode(s string) (PageMode, error)

// DocumentStatus 写作文档状态
type DocumentStatus string

const (
    DocumentStatusDraft      DocumentStatus = "draft"
    DocumentStatusPublished  DocumentStatus = "published"
    DocumentStatusArchived   DocumentStatus = "archived"
    DocumentStatusDeleted    DocumentStatus = "deleted"
)

// IsValid 检查状态是否有效
func (s DocumentStatus) IsValid() bool

// String 转换为字符串
func (s DocumentStatus) String() string

// ParseDocumentStatus 从字符串解析状态
func ParseDocumentStatus(s string) (DocumentStatus, error)

// WithdrawalStatus 提现状态
type WithdrawalStatus string

const (
    WithdrawalStatusPending   WithdrawalStatus = "pending"
    WithdrawalStatusApproved  WithdrawalStatus = "approved"
    WithdrawalStatusRejected  WithdrawalStatus = "rejected"
    WithdrawalStatusCompleted WithdrawalStatus = "completed"
)

// IsValid 检查状态是否有效
func (s WithdrawalStatus) IsValid() bool

// String 转换为字符串
func (s WithdrawalStatus) String() string

// ParseWithdrawalStatus 从字符串解析状态
func ParseWithdrawalStatus(s string) (WithdrawalStatus, error)
```

### 使用示例
```go
// Model 层定义
type User struct {
    ID       primitive.ObjectID `bson:"_id"`
    Username string             `bson:"username"`
    Role     types.UserRole     `bson:"role" validate:"required"`
}

// 校验
func (u *User) Validate() error {
    if !u.Role.IsValid() {
        return fmt.Errorf("invalid role: %s", u.Role)
    }
    return nil
}

// API 层使用
type UpdateRoleRequest struct {
    Role string `json:"role" validate:"required,oneof=reader author admin"`
}

func (h *Handler) UpdateRole(c *gin.Context) {
    var req UpdateRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    role, err := types.ParseUserRole(req.Role)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid role"})
        return
    }

    // 更新用户角色
    // ...
}
```

---

## 6. JSON/BSON 命名：`json_bson.go`

### 路径
`Qingyu_backend/models/shared/types/json_bson.go`

### 用途
- 统一 JSON/BSON 命名约束
- 通用 tag/映射辅助
- DTO 转换边界说明

### 核心功能
```go
package types

// NamingConvention 命名规范
type NamingConvention int

const (
    // CamelCase JSON 命名（API 对外）
    CamelCase NamingConvention = iota
    // snake_case BSON 命名（存储）
    SnakeCase
)

// FieldTags 字段标签配置
type FieldTags struct {
    JSON string // JSON 标签
    BSON string // BSON 标签
}

// StandardFieldTags 生成标准字段标签
// JSON: camelCase, BSON: snake_case
func StandardFieldTags(fieldName string) FieldTags

// ModelTags 生成 Model 层标签（只 BSON）
func ModelTags(fieldName string) string

// DTOTags 生成 DTO 层标签（只 JSON）
func DTOTags(fieldName string) string

// TagSet 标签集合（用于生成 struct tag）
type TagSet map[string]string

// String 生成 struct tag 字符串
func (t TagSet) String() string

// CommonTags 常用标签集合
var CommonTags = TagSet{
    "bson":    "",
    "json":    "",
    "validate": "",
}

// WithBSON 添加 BSON 标签
func (t TagSet) WithBSON(name string) TagSet

// WithJSON 添加 JSON 标签
func (t TagSet) WithJSON(name string) TagSet

// WithValidate 添加 validate 标签
func (t TagSet) WithValidate(rule string) TagSet
```

### 使用示例
```go
// Model 层（只 BSON，snake_case）
type Book struct {
    ID       primitive.ObjectID `bson:"_id"`
    AuthorID primitive.ObjectID `bson:"author_id"`
    Title    string             `bson:"title"`
}

// DTO 层（只 JSON，camelCase）
type BookDTO struct {
    ID       string `json:"id" validate:"required"`
    AuthorID string `json:"authorId" validate:"required"`
    Title    string `json:"title"`
}

// 或使用辅助工具
var tags = types.CommonTags.
    WithBSON("author_id").
    WithJSON("authorId").
    WithValidate("required")

// 生成: `bson:"author_id" json:"authorId" validate:"required"`
```

---

## 7. 通用转换：`converter.go`

### 路径
`Qingyu_backend/models/shared/types/converter.go`

### 用途
- 通用 Model ↔ DTO 转换
- 集中管理转换逻辑

### 核心功能
```go
package types

// Converter 转换器接口
type Converter interface {
    ToDTO() interface{}
    FromDTO(dto interface{}) error
}

// DTOConverter DTO 转换辅助
type DTOConverter struct{}

// ModelIDToDTO Model ID → DTO ID (ObjectID → string)
func (DTOConverter) ModelIDToDTO(id primitive.ObjectID) string

// DTOIDToModel DTO ID → Model ID (string → ObjectID)
func (DTOConverter) DTOIDToModel(id string) (primitive.ObjectID, error)

// MoneyToDTO Money → 金额字符串 (Money → "¥12.99")
func (DTOConverter) MoneyToDTO(money Money) string

// DTOMoneyToModel 金额字符串 → Money ("¥12.99" → Money)
func (DTOConverter) DTOMoneyToModel(s string) (Money, error)

// RatingToDTO Rating → 评分字符串 (Rating → "4.5")
func (DTOConverter) RatingToDTO(rating Rating) string

// DTORatingToModel 评分字符串 → Rating ("4.5" → Rating)
func (DTOConverter) DTORatingToModel(s string) (Rating, error)

// ProgressToDTO Progress → 百分比 (Progress → 75)
func (DTOConverter) ProgressToDTO(progress Progress) int

// DTOProgressToModel 百分比 → Progress (75 → Progress)
func (DTOConverter) DTOProgressToModel(percent int) (Progress, error)
```

---

## 8. 文档：`model-consistency-types.md`

### 路径
`docs/architecture/model-consistency-types.md`

### 用途
- 集中记录公共类型定义
- 字段单位、范围与序列化口径
- 作为类型使用的"唯一真源"

### 内容结构
```markdown
# 公共类型定义参考

## 1. ID 类型
- 存储层：`primitive.ObjectID`
- API/Service 层：`string` (hex)
- 转换工具：`types.ParseObjectID()`, `types.ToHex()`

## 2. 金额类型
- 存储/计算：`types.Money` (int64, 分)
- 展示：`money.ToYuan()` (float64, 元)
- 格式化：`money.String()` ("¥12.99")

## 3. 评分类型
- 范围：0.0-5.0
- 分布：`map[string]int64` (key: "1"-"5")
- 校验：`rating.IsValid()`

## 4. 进度类型
- 内部：0.0-1.0
- 展示：0-100 (percent)
- 转换：`progress.ToPercent()`

## 5. 枚举类型
- 角色：reader/author/admin
- PageMode：scroll/paginate
- 文档状态：draft/published/archived/deleted

## 6. 命名规范
- JSON：camelCase
- BSON：snake_case
```

---

## 实施计划

### Phase 1: 创建基础类型（独立任务）

优先级：**High**

1. 创建 `id.go`
2. 创建 `money.go`
3. 创建 `rating.go`
4. 创建 `progress.go`
5. 创建 `enums.go`
6. 创建 `json_bson.go`
7. 创建 `converter.go`

### Phase 2: 文档更新

1. 创建 `model-consistency-types.md`
2. 更新 `model-consistency-fix-guide.md` 引用新类型

### Phase 3: 逐步迁移

按照 `model-consistency-fix-guide.md` 中的优先级，逐模块迁移到新类型

---

## 验收标准

- [ ] 所有类型文件已创建并包含完整实现
- [ ] 所有类型有单元测试覆盖
- [ ] 文档已更新并与代码同步
- [ ] 至少有一个模块已完成迁移作为示例
- [ ] 代码通过 `go vet` 和 `golangci-lint` 检查

---

## 相关文档

- [模型一致性修复指南](./model-consistency-fix-guide.md)
- [ID 类型统一标准](./id-type-unification-standard.md)
- [模型修复任务提示词](./model-fix-task-prompt.md)
