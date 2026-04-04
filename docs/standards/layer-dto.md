# DTO 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 职责边界与依赖关系

### 1.1 职责定义

DTO（Data Transfer Object）层是**数据传输对象层**，负责：

1. **数据传输**：在 API 层和 Service 层之间传递数据
2. **数据隔离**：隐藏内部 Model 结构，避免直接暴露数据库模型
3. **格式转换**：将 Model 数据转换为适合前端使用的格式
4. **验证规则**：定义输入数据的验证规则
5. **API 契约**：定义 API 请求和响应的数据结构

### 1.2 与上下层级的交互边界

```
┌─────────────────────────────────────────────────────────┐
│                      API 层                             │
│              (HTTP 请求/响应处理)                       │
└─────────────────────────────────────────────────────────┘
                         │
                    使用 DTO 类型
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                      DTO 层                             │
│  ┌─────────────────────────────────────────────────────┤
│  │ 职责：                                              │
│  │ - 定义数据传输结构                                  │
│  │ - 定义验证规则（validate 标签）                     │
│  │ - ID/时间格式转换（string 类型）                    │
│  │ - 输入/输出结构分离                                 │
│  └─────────────────────────────────────────────────────┤
│  禁止：                                                  │
│  │ - 包含业务逻辑                                      │
│  │ - 直接操作数据库                                    │
│  │ - 包含数据库标签（bson tag）                        │
└─────────────────────────────────────────────────────────┘
                         │
                    Converter 转换
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Model 层                             │
│              (数据库实体定义)                           │
└─────────────────────────────────────────────────────────┘
```

### 1.3 DTO 与 Model 的区别

| 特性 | DTO | Model |
|------|-----|-------|
| 位置 | `models/dto/` | `models/bookstore/` |
| 用途 | API 数据传输 | 数据库映射 |
| ID 类型 | `string` | `primitive.ObjectID` |
| 时间类型 | `string` (ISO8601) 或 `time.Time` | `time.Time` |
| JSON 标签 | camelCase | camelCase |
| BSON 标签 | ❌ 无 | ✅ snake_case |
| 验证标签 | ✅ validate | ❌ 无 |

### 1.4 依赖关系

```go
// DTO 层允许的依赖
import (
    "time"
    "Qingyu_backend/models/shared/types"  // 类型转换器
)

// DTO 层禁止的依赖
import (
    "go.mongodb.org/mongo-driver/bson"     // ❌ 禁止 BSON 标签
    "Qingyu_backend/service/xxx"          // ❌ 禁止依赖 Service
    "Qingyu_backend/repository/xxx"       // ❌ 禁止依赖 Repository
)
```

---

## 2. 命名与代码规范

### 2.1 文件命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 领域 DTO | `{领域}.go` | `bookstore.go`, `user.go` |
| 功能 DTO | `{功能}_dto.go` | `auth_dto.go`, `writer_dto.go` |
| 转换器 | `{领域}_converter.go` | `writer_converter.go` |
| 测试文件 | `{文件名}_test.go` | `writer_dto_test.go` |

### 2.2 结构体命名规范

| 类型 | 后缀 | 示例 |
|------|------|------|
| 请求 DTO | `Request` | `CreateBookRequest`, `LoginRequest` |
| 更新请求 | `UpdateRequest` | `UpdateBookRequest`, `UpdateProfileRequest` |
| 响应 DTO | `Response` | `BookResponse`, `UserProfileResponse` |
| 列表请求 | `ListRequest` | `ListBooksRequest` |
| 列表响应 | `ListResponse` | `BookListResponse` |
| 数据传输 | `DTO` | `BookDTO` |

### 2.3 字段命名规范

```go
// DTO 字段命名
type BookDTO struct {
    ID          string   `json:"id"`           // ID 使用 string
    CreatedAt   string   `json:"createdAt"`    // 时间使用 string (ISO8601)
    UpdatedAt   string   `json:"updatedAt"`

    // 基本信息
    Title       string   `json:"title"`        // JSON 使用 camelCase
    Author      string   `json:"author"`
    CategoryIDs []string `json:"categoryIds"`  // 切片类型
    Tags        []string `json:"tags"`
    IsFree      bool     `json:"isFree"`       // 布尔类型

    // 可选字段
    Cover       *string  `json:"cover,omitempty"`      // 指针表示可选
    PublishedAt string   `json:"publishedAt,omitempty"` // 空字符串时省略
}
```

### 2.4 目录组织规范

```
models/dto/
├── bookstore.go           # 书城 DTO
├── user.go                # 用户 DTO
├── reader.go              # 读者 DTO
├── writer_dto.go          # 作家 DTO
├── writer_converter.go    # 作家 DTO 转换器
├── content_dto.go         # 内容 DTO
├── content_dto_test.go    # 内容 DTO 测试
├── audit.go               # 审计 DTO
└── README.md              # 快速参考
```

---

## 3. 设计模式与最佳实践

### 3.1 输入/输出分离模式

```go
// ===========================
// 请求 DTO（输入）
// ===========================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
    Title       string   `json:"title" validate:"required,min=1,max=100"`
    Summary     string   `json:"summary,omitempty" validate:"max=500"`
    CoverURL    string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
    Tags        []string `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
    Category    string   `json:"category,omitempty" validate:"max=50"`
    WritingType string   `json:"writingType,omitempty" validate:"omitempty,oneof=novel article script"`
}

// UpdateProjectRequest 更新项目请求（使用指针表示可选字段）
type UpdateProjectRequest struct {
    Title       *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
    Summary     *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
    CoverURL    *string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
    Tags        *[]string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=50"`
    Category    *string   `json:"category,omitempty" validate:"omitempty,max=50"`
    Status      *string   `json:"status,omitempty" validate:"omitempty,oneof=draft serializing completed suspended archived"`
    WritingType *string   `json:"writingType,omitempty" validate:"omitempty,oneof=novel article script"`
}

// ===========================
// 响应 DTO（输出）
// ===========================

// ProjectResponse 项目响应
type ProjectResponse struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Summary   string    `json:"summary"`
    CoverURL  string    `json:"coverUrl"`
    Tags      []string  `json:"tags"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

### 3.2 DTO 转换模式

```go
// bookstore_converter.go

// ToBookDTO Model → DTO 转换
func ToBookDTO(book *bookstore.Book) *dto.BookDTO {
    if book == nil {
        return nil
    }

    var converter types.DTOConverter

    // 处理可空时间字段
    var publishedAt, lastUpdateAt string
    if book.PublishedAt != nil {
        publishedAt = converter.TimeToISO8601(*book.PublishedAt)
    }
    if book.LastUpdateAt != nil {
        lastUpdateAt = converter.TimeToISO8601(*book.LastUpdateAt)
    }

    return &dto.BookDTO{
        ID:          converter.ModelIDToDTO(book.ID),
        CreatedAt:   converter.TimeToISO8601(book.CreatedAt),
        UpdatedAt:   converter.TimeToISO8601(book.UpdatedAt),
        Title:       book.Title,
        Author:      book.Author,
        // ... 其他字段
        PublishedAt: publishedAt,
        LastUpdateAt: lastUpdateAt,
    }
}

// ToBookDTOsFromPtrSlice 批量转换
func ToBookDTOsFromPtrSlice(books []*bookstore.Book) []*dto.BookDTO {
    result := make([]*dto.BookDTO, len(books))
    for i := range books {
        result[i] = ToBookDTO(books[i])
    }
    return result
}

// ToBookModel DTO → Model 转换（用于更新）
func ToBookModel(dto *dto.BookDTO) (*bookstore.Book, error) {
    if dto == nil {
        return nil, nil
    }

    var converter types.DTOConverter

    id, err := converter.DTOIDToModel(dto.ID)
    if err != nil {
        return nil, err
    }

    // ... 转换其他字段

    return &bookstore.Book{
        IdentifiedEntity: shared.IdentifiedEntity{ID: id},
        // ... 其他字段
    }, nil
}
```

### 3.3 分页 DTO 模式

```go
// 列表请求
type ListProjectsRequest struct {
    Page     int    `form:"page" validate:"min=1"`
    PageSize int    `form:"page_size" validate:"min=1,max=100"`
    Status   string `form:"status" validate:"omitempty,oneof=draft serializing completed"`
    Sort     string `form:"sort" validate:"omitempty,oneof=created_at updated_at title"`
    Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

// 列表响应
type ProjectListResponse struct {
    Items    []ProjectResponse `json:"items"`
    Total    int64             `json:"total"`
    Page     int               `json:"page"`
    PageSize int               `json:"pageSize"`
}
```

### 3.4 枚举 DTO 模式

```go
// DocumentType 文档类型枚举
type DocumentType string

const (
    DocumentTypeVolume  DocumentType = "volume"  // 卷
    DocumentTypeChapter DocumentType = "chapter" // 章
    DocumentTypeSection DocumentType = "section" // 节
    DocumentTypeScene   DocumentType = "scene"   // 场景
)

// DocumentStatus 文档状态枚举
type DocumentStatus string

const (
    DocumentStatusPlanned   DocumentStatus = "planned"   // 计划中
    DocumentStatusWriting   DocumentStatus = "writing"   // 写作中
    DocumentStatusCompleted DocumentStatus = "completed" // 已完成
)

// 在 DTO 中使用
type CreateDocumentRequest struct {
    Type   DocumentType   `json:"type" validate:"required,oneof=volume chapter section scene"`
    Status DocumentStatus `json:"status" validate:"omitempty,oneof=planned writing completed"`
}
```

### 3.5 嵌套 DTO 模式

```go
// 树形结构响应
type DocumentTreeResponse struct {
    ProjectID string              `json:"projectId"`
    Documents []*DocumentTreeItem `json:"documents"`
}

// DocumentTreeItem 文档树节点
type DocumentTreeItem struct {
    ID        string              `json:"id"`
    ParentID  *string             `json:"parentId,omitempty"`
    Title     string              `json:"title"`
    Type      DocumentType        `json:"type"`
    Level     int                 `json:"level"`
    OrderKey  string              `json:"orderKey"`
    WordCount int                 `json:"wordCount"`
    Children  []*DocumentTreeItem `json:"children,omitempty"` // 嵌套子节点
}
```

### 3.6 反模式警示

❌ **禁止事项**：

```go
// ❌ 禁止：在 DTO 中使用 BSON 标签
type BookDTO struct {
    ID   string `json:"id" bson:"_id"`  // DTO 不应该有 BSON 标签
    Title string `json:"title" bson:"title"`
}

// ❌ 禁止：DTO 中包含业务逻辑
type BookDTO struct {
    Price float64 `json:"price"`
}

func (d *BookDTO) GetDiscountedPrice() float64 {
    return d.Price * 0.9  // 业务逻辑应该在 Service 层
}

// ❌ 禁止：DTO 直接使用 ObjectID 类型
type BookDTO struct {
    ID primitive.ObjectID `json:"id"`  // 应该使用 string
}

// ❌ 禁止：响应 DTO 包含验证标签
type BookResponse struct {
    Title string `json:"title" validate:"required"`  // 响应不需要验证
}
```

---

## 4. 验证规则规范

### 4.1 常用验证标签

```go
type CreateBookRequest struct {
    // 必填验证
    Title string `json:"title" validate:"required"`

    // 长度验证
    Title       string `json:"title" validate:"required,min=1,max=200"`
    Summary     string `json:"summary" validate:"max=500"`
    Description string `json:"description" validate:"min=10,max=5000"`

    // 数值验证
    Price    float64 `json:"price" validate:"gte=0"`
    Rating   float64 `json:"rating" validate:"min=0,max=5"`
    Age      int     `json:"age" validate:"gte=0,lte=150"`
    Page     int     `json:"page" validate:"min=1"`
    PageSize int     `json:"pageSize" validate:"min=1,max=100"`

    // 格式验证
    Email   string `json:"email" validate:"required,email"`
    URL     string `json:"url" validate:"omitempty,url"`
    Phone   string `json:"phone" validate:"omitempty,e164"`
    HexID   string `json:"id" validate:"hexadecimal,len=24"`

    // 枚举验证
    Status string `json:"status" validate:"required,oneof=draft ongoing completed paused deleted"`
    Sort   string `json:"sort" validate:"omitempty,oneof=asc desc"`

    // 切片验证
    Tags        []string `json:"tags" validate:"max=10"`
    Tags        []string `json:"tags" validate:"max=10,dive,min=1,max=50"`  // 每个元素也验证
    CategoryIDs []string `json:"categoryIds" validate:"required,min=1"`

    // 条件验证
    CoverURL string `json:"coverUrl" validate:"omitempty,url,max=500"`

    // 可选指针字段
    Category *string `json:"category,omitempty" validate:"omitempty,max=50"`
}
```

### 4.2 自定义验证

```go
// 注册自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
    // 注册 "objectid" 验证器
    v.RegisterValidation("objectid", func(fl validator.FieldLevel) bool {
        value := fl.Field().String()
        if value == "" {
            return true  // 空值由 required 处理
        }
        _, err := primitive.ObjectIDFromHex(value)
        return err == nil
    })
}

// 使用自定义验证器
type GetBookRequest struct {
    ID string `json:"id" validate:"required,objectid"`
}
```

### 4.3 验证错误消息

```go
// 验证错误响应格式
type ValidationErrorResponse struct {
    Code    int               `json:"code"`
    Message string            `json:"message"`
    Errors  map[string]string `json:"errors"`  // 字段名 -> 错误消息
}

// 示例响应
{
    "code": 400,
    "message": "请求参数验证失败",
    "errors": {
        "title": "标题不能为空",
        "price": "价格必须大于等于0",
        "email": "邮箱格式不正确"
    }
}
```

---

## 5. 类型转换规范

### 5.1 ID 转换

```go
// DTO 中：string 类型
type BookDTO struct {
    ID         string   `json:"id"`
    AuthorID   string   `json:"authorId"`
    CategoryID string   `json:"categoryId,omitempty"`
    CategoryIDs []string `json:"categoryIds"`
}

// Model 中：ObjectID 类型
type Book struct {
    ID          primitive.ObjectID   `bson:"_id"`
    AuthorID    string               `bson:"author_id"`  // 或 ObjectID
    CategoryIDs []primitive.ObjectID `bson:"category_ids"`
}

// 转换器
var converter types.DTOConverter

// 单个 ID 转换
id := converter.ModelIDToDTO(book.ID)           // ObjectID → string
oid, err := converter.DTOIDToModel(dto.ID)      // string → ObjectID

// 批量 ID 转换
ids := converter.ModelIDsToDTO(book.CategoryIDs)       // []ObjectID → []string
oids, err := converter.DTOIDsToModel(dto.CategoryIDs)  // []string → []ObjectID
```

### 5.2 时间转换

```go
// DTO 中：string (ISO8601) 或 time.Time
type BookDTO struct {
    CreatedAt   string `json:"createdAt"`    // ISO8601 格式
    UpdatedAt   string `json:"updatedAt"`
    PublishedAt string `json:"publishedAt,omitempty"`  // 可选时间
}

// Model 中：time.Time
type Book struct {
    CreatedAt   time.Time  `bson:"created_at"`
    UpdatedAt   time.Time  `bson:"updated_at"`
    PublishedAt *time.Time `bson:"published_at,omitempty"`
}

// 转换器
createdAt := converter.TimeToISO8601(book.CreatedAt)           // time.Time → string
t, err := converter.ISO8601ToTime(dto.CreatedAt)               // string → time.Time
```

### 5.3 金额转换

```go
// DTO 中：string（避免浮点精度问题）
type BookDTO struct {
    Price string `json:"price" validate:"omitempty,numeric"`  // "99.99"
}

// Model 中：int64（分为单位）或 float64
type Book struct {
    Price float64 `bson:"price"`  // 或 int64 存储分
}

// 转换
money := types.NewMoneyFromCents(int64(book.Price))  // 分 → Money
price := money.String()                               // Money → string "99.99"
```

---

## 6. 测试策略

### 6.1 单元测试编写指南

```go
// content_dto_test.go
package dto

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestToBookDTO(t *testing.T) {
    // 1. 准备测试数据
    now := time.Now()
    book := &bookstore.Book{
        IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
        BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
        Title:            "测试书籍",
        Author:           "测试作者",
        Status:           bookstore.BookStatusOngoing,
    }

    // 2. 执行转换
    dto := ToBookDTO(book)

    // 3. 验证结果
    assert.NotNil(t, dto)
    assert.Equal(t, book.Title, dto.Title)
    assert.Equal(t, book.Author, dto.Author)
    assert.NotEmpty(t, dto.ID)
    assert.NotEmpty(t, dto.CreatedAt)
}

func TestToBookModel(t *testing.T) {
    dto := &dto.BookDTO{
        ID:     "507f1f77bcf86cd799439011",
        Title:  "测试书籍",
        Status: "ongoing",
    }

    model, err := ToBookModel(dto)

    assert.NoError(t, err)
    assert.NotNil(t, model)
    assert.Equal(t, dto.Title, model.Title)
}
```

### 6.2 验证测试

```go
func TestCreateBookRequest_Validation(t *testing.T) {
    validate := validator.New()

    // 测试有效请求
    validReq := CreateBookRequest{
        Title:  "测试书籍",
        Author: "测试作者",
        Price:  99.99,
    }
    err := validate.Struct(validReq)
    assert.NoError(t, err)

    // 测试无效请求（标题为空）
    invalidReq := CreateBookRequest{
        Title:  "",
        Author: "测试作者",
    }
    err = validate.Struct(invalidReq)
    assert.Error(t, err)
}
```

---

## 7. 完整代码示例

### 7.1 完整 DTO 定义示例

```go
// writer_dto.go
package dto

import "time"

// ===========================
// Writer DTO（符合分层架构规范）
// ===========================
//
// 命名和标签规范：
// - DTO 结构体使用驼峰命名（PascalCase）
// - JSON 字段标签使用驼峰命名（camelCase）
// - 对应的 MongoDB 模型（位于 models/writer/）使用蛇形命名（snake_case）的 BSON 标签
//
// 用途：
// - 用于 Service 层和 API 层之间的数据传输
// - ID 和时间字段统一使用字符串类型
// - 避免直接暴露 MongoDB 模型到 API 层

// ===========================
// Project DTOs
// ===========================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
    Title       string   `json:"title" validate:"required,min=1,max=100"`
    Summary     string   `json:"summary,omitempty" validate:"max=500"`
    CoverURL    string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
    Tags        []string `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
    Category    string   `json:"category,omitempty" validate:"max=50"`
    WritingType string   `json:"writingType,omitempty" validate:"omitempty,oneof=novel article script"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
    Title       *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
    Summary     *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
    CoverURL    *string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
    Tags        *[]string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=50"`
    Category    *string   `json:"category,omitempty" validate:"omitempty,max=50"`
    Status      *string   `json:"status,omitempty" validate:"omitempty,oneof=draft serializing completed suspended archived"`
    WritingType *string   `json:"writingType,omitempty" validate:"omitempty,oneof=novel article script"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Summary   string    `json:"summary"`
    CoverURL  string    `json:"coverUrl"`
    Tags      []string  `json:"tags"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// ListProjectsRequest 查询参数用于列出项目
type ListProjectsRequest struct {
    Page     int    `form:"page" validate:"min=1"`
    PageSize int    `form:"page_size" validate:"min=1,max=100"`
    Status   string `form:"status" validate:"omitempty,oneof=draft serializing completed suspended archived"`
    Sort     string `form:"sort" validate:"omitempty,oneof=created_at updated_at title"`
    Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

// ProjectListResponse 分页项目列表响应
type ProjectListResponse struct {
    Items    []ProjectResponse `json:"items"`
    Total    int64             `json:"total"`
    Page     int               `json:"page"`
    PageSize int               `json:"pageSize"`
}

// ===========================
// Document DTOs
// ===========================

// DocumentType 文档类型枚举
type DocumentType string

const (
    DocumentTypeVolume  DocumentType = "volume"  // 卷
    DocumentTypeChapter DocumentType = "chapter" // 章
    DocumentTypeSection DocumentType = "section" // 节
    DocumentTypeScene   DocumentType = "scene"   // 场景
)

// DocumentStatus 文档状态枚举
type DocumentStatus string

const (
    DocumentStatusPlanned   DocumentStatus = "planned"   // 计划中
    DocumentStatusWriting   DocumentStatus = "writing"   // 写作中
    DocumentStatusCompleted DocumentStatus = "completed" // 已完成
)

// CreateDocumentRequest 创建文档请求（树形结构）
type CreateDocumentRequest struct {
    ProjectID    string       `json:"projectId" validate:"required"`
    ParentID     *string      `json:"parentId,omitempty"`
    Title        string       `json:"title" validate:"required,min=1,max=200"`
    Type         DocumentType `json:"type" validate:"required,oneof=volume chapter section scene"`
    Level        int          `json:"level" validate:"min=0,max=2"`
    OrderKey     string       `json:"orderKey,omitempty"`
    CharacterIDs []string     `json:"characterIds,omitempty"`
    LocationIDs  []string     `json:"locationIds,omitempty"`
    TimelineIDs  []string     `json:"timelineIds,omitempty"`
    Tags         []string     `json:"tags,omitempty"`
    Notes        string       `json:"notes,omitempty"`
}

// UpdateDocumentRequest 更新文档元数据请求
type UpdateDocumentRequest struct {
    Title        *string         `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
    Status       *DocumentStatus `json:"status,omitempty" validate:"omitempty,oneof=planned writing completed"`
    CharacterIDs *[]string       `json:"characterIds,omitempty" validate:"omitempty,max=50"`
    LocationIDs  *[]string       `json:"locationIds,omitempty" validate:"omitempty,max=50"`
    TimelineIDs  *[]string       `json:"timelineIds,omitempty" validate:"omitempty,max=50"`
    Tags         *[]string       `json:"tags,omitempty" validate:"omitempty,max=20"`
    Notes        *string         `json:"notes,omitempty" validate:"omitempty,max=1000"`
    OrderKey     *string         `json:"orderKey,omitempty"`
}

// DocumentResponse 文档响应（树形结构）
type DocumentResponse struct {
    ID           string         `json:"id"`
    ProjectID    string         `json:"projectId"`
    ParentID     *string        `json:"parentId,omitempty"`
    Title        string         `json:"title"`
    Type         DocumentType   `json:"type"`
    Level        int            `json:"level"`
    Order        int            `json:"order"`
    OrderKey     string         `json:"orderKey"`
    Status       DocumentStatus `json:"status"`
    WordCount    int            `json:"wordCount"`
    CharacterIDs []string       `json:"characterIds,omitempty"`
    LocationIDs  []string       `json:"locationIds,omitempty"`
    TimelineIDs  []string       `json:"timelineIds,omitempty"`
    Tags         []string       `json:"tags,omitempty"`
    Notes        string         `json:"notes,omitempty"`
    CreatedAt    time.Time      `json:"createdAt"`
    UpdatedAt    time.Time      `json:"updatedAt"`
}

// DocumentTreeResponse 文档树响应（嵌套结构）
type DocumentTreeResponse struct {
    ProjectID string              `json:"projectId"`
    Documents []*DocumentTreeItem `json:"documents"`
}

// DocumentTreeItem 文档树节点
type DocumentTreeItem struct {
    ID        string              `json:"id"`
    ParentID  *string             `json:"parentId,omitempty"`
    Title     string              `json:"title"`
    Type      DocumentType        `json:"type"`
    Level     int                 `json:"level"`
    OrderKey  string              `json:"orderKey"`
    WordCount int                 `json:"wordCount"`
    Children  []*DocumentTreeItem `json:"children,omitempty"`
}
```

---

## 8. 参考资料

- [DTO 层快速参考](../models/dto/README.md)
- [API 层设计说明](./layer-api.md)
- [Model 层设计说明](./layer-models.md)
- [类型转换器文档](../models/shared/types/README.md)

---

*最后更新：2026-03-19*
