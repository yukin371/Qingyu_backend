# ID 类型统一参考标准

> **目标**：建立 Qingyu Backend 项目中 ID 类型的统一使用标准，消除跨模块 ID 类型混用问题。

## 核心原则

```
┌─────────────────────────────────────────────────────────────────┐
│                        ID 类型分层标准                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐    ┌─────────────────┐    ┌──────────────┐ │
│  │   API Layer     │    │  Service Layer  │    │ Repository   │ │
│  │  (DTO/Request)  │◄──►│  (Business)     │◄──►│ /DAO Layer   │ │
│  │                 │    │                 │    │              │ │
│  │   string        │    │    string       │    │ ObjectID     │ │
│  │   (hex)         │    │    (hex)        │    │ (primitive)  │ │
│  └─────────────────┘    └─────────────────┘    └──────────────┘ │
│         │                       │                      │        │
│         │                       │                      │        │
│         ▼                       ▼                      ▼        │
│  对外统一使用 string     业务逻辑使用 string      存储使用 ObjectID │
│  便于 JSON 序列化       与 API 层一致           原生 Mongo 类型  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

转换边界：
  API ↔ Service: string ↔ string (无需转换，保持一致性)
  Service ↔ Repository: string ↔ ObjectID (边界转换)
```

## 各层职责规范

### 1. 存储层（Models/DAO）

**规则**：
- 统一使用 `primitive.ObjectID`
- 所有 MongoDB 文档的主键必须是 `bson:"_id"`
- 外键关联字段也使用 `ObjectID`

**示例**：

```go
// models/bookstore/book.go
package bookstore

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
    ID       primitive.ObjectID `bson:"_id" json:"id"`
    AuthorID primitive.ObjectID `bson:"author_id" json:"-"` // 对外不直接暴露
    Title    string             `bson:"title" json:"title"`
    // ...
}
```

### 2. Repository/DAO 层

**规则**：
- 只接收和返回 `primitive.ObjectID`
- 不处理 ID 转换逻辑
- 查询条件使用 ObjectID

**示例**：

```go
// repository/bookstore/book_repository.go
package bookstore

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookRepository interface {
    FindByID(ctx context.Context, id primitive.ObjectID) (*models.Book, error)
    FindByAuthor(ctx context.Context, authorID primitive.ObjectID) ([]*models.Book, error)
    Create(ctx context.Context, book *models.Book) error
}

func (r *bookRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Book, error) {
    filter := bson.M{"_id": id}
    // ...
}
```

### 3. Service 层

**规则**：
- 对外接口只接收和返回 `string`
- 在调用 Repository 前转换 string → ObjectID
- 在返回结果前转换 ObjectID → string

**示例**：

```go
// service/bookstore/book_service.go
package bookstore

type BookService interface {
    GetBook(ctx context.Context, bookID string) (*dto.BookDTO, error)
    GetBooksByAuthor(ctx context.Context, authorID string) ([]*dto.BookDTO, error)
}

func (s *bookServiceImpl) GetBook(ctx context.Context, bookID string) (*dto.BookDTO, error) {
    // Service 层入口：string → ObjectID
    oid, err := idutil.ParseObjectID(bookID)
    if err != nil {
        return nil, fmt.Errorf("invalid book ID: %w", err)
    }

    // 调用 Repository（传递 ObjectID）
    book, err := s.repo.FindByID(ctx, oid)
    if err != nil {
        return nil, err
    }

    // 转换为 DTO（ObjectID → string）
    return s.toDTO(book), nil
}
```

### 4. API 层（Handler/Controller）

**规则**：
- 所有请求/响应 DTO 使用 `string`
- 不直接处理 `primitive.ObjectID`
- 传递给 Service 层的已经是 string

**示例**：

```go
// api/v1/bookstore/book_api.go
package bookstore

type BookDTO struct {
    ID       string `json:"id"`
    AuthorID string `json:"authorId"`
    Title    string `json:"title"`
    // ...
}

type GetBookRequest struct {
    ID string `uri:"id" validate:"required"`
}

func (h *BookAPI) GetBook(c *gin.Context) {
    var req GetBookRequest
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 直接传递 string 给 Service
    book, err := h.service.GetBook(c.Request.Context(), req.ID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, book)
}
```

## ID 转换工具包

### 工具包位置

```
Qingyu_backend/pkg/idutil/converter.go
```

或者扩展现有的：
```
Qingyu_backend/models/base/identifiers.go
```

### 转换函数实现

```go
// pkg/idutil/converter.go
package idutil

import (
    "errors"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
    ErrInvalidIDFormat = errors.New("invalid ID format")
    ErrEmptyID         = errors.New("ID cannot be empty")
)

// ParseObjectID 将 hex 字符串解析为 ObjectID
// 输入：24字符的 hex 字符串
// 输出：primitive.ObjectID 或 error
func ParseObjectID(s string) (primitive.ObjectID, error) {
    if s == "" {
        return primitive.NilObjectID, ErrEmptyID
    }
    
    oid, err := primitive.ObjectIDFromHex(s)
    if err != nil {
        return primitive.NilObjectID, fmt.Errorf("%w: %s", ErrInvalidIDFormat, s)
    }
    
    return oid, nil
}

// MustParseObjectID 解析 ObjectID，panic on error
// 仅在测试或确定 ID 有效时使用
func MustParseObjectID(s string) primitive.ObjectID {
    oid, err := ParseObjectID(s)
    if err != nil {
        panic(err)
    }
    return oid
}

// ToHex 将 ObjectID 转换为 hex 字符串
// 输入：primitive.ObjectID
// 输出：24字符的 hex 字符串
func ToHex(id primitive.ObjectID) string {
    return id.Hex()
}

// IsValidObjectID 检查字符串是否为有效的 ObjectID hex 格式
func IsValidObjectID(s string) bool {
    _, err := primitive.ObjectIDFromHex(s)
    return err == nil
}

// ParseObjectIDSlice 批量解析 ID 字符串
// 返回：成功解析的 ObjectID 列表和失败索引的映射
func ParseObjectIDSlice(ss []string) ([]primitive.ObjectID, map[int]error) {
    oids := make([]primitive.ObjectID, 0, len(ss))
    errs := make(map[int]error)
    
    for i, s := range ss {
        oid, err := ParseObjectID(s)
        if err != nil {
            errs[i] = err
            continue
        }
        oids = append(oids, oid)
    }
    
    return oids, errs
}

// ToHexSlice 批量转换 ObjectID 为 hex 字符串
func ToHexSlice(ids []primitive.ObjectID) []string {
    result := make([]string, len(ids))
    for i, id := range ids {
        result[i] = ToHex(id)
    }
    return result
}

// GenerateNewObjectID 生成新的 ObjectID 并返回 hex 字符串
func GenerateNewObjectID() string {
    return primitive.NewObjectID().Hex()
}
```

### DTO 转换辅助函数

```go
// pkg/idutil/dto.go
package idutil

import "go.mongodb.org/mongo-driver/bson/primitive"

// DTOConverter 提供 Model ↔ DTO 的 ID 字段转换
type DTOConverter struct{}

// ModelIDToDTO Model 层的 ObjectID → DTO 层的 string
func (DTOConverter) ModelIDToDTO(id primitive.ObjectID) string {
    return ToHex(id)
}

// ModelIDsToDTO 批量转换
func (DTOConverter) ModelIDsToDTO(ids []primitive.ObjectID) []string {
    return ToHexSlice(ids)
}

// DTOIDToModel DTO 层的 string → Model 层的 ObjectID
func (DTOConverter) DTOIDToModel(id string) (primitive.ObjectID, error) {
    return ParseObjectID(id)
}

// DTOIDsToModel 批量转换
func (DTOConverter) DTOIDsToModel(ids []string) ([]primitive.ObjectID, error) {
    oids, errs := ParseObjectIDSlice(ids)
    if len(errs) > 0 {
        return nil, fmt.Errorf("failed to parse %d IDs", len(errs))
    }
    return oids, nil
}
```

## 现有代码迁移指南

### 迁移步骤

#### Phase 1: 创建转换工具（不影响现有代码）

```bash
# 1. 创建工具包
mkdir -p Qingyu_backend/pkg/idutil
touch Qingyu_backend/pkg/idutil/converter.go
touch Qingyu_backend/pkg/idutil/dto.go
touch Qingyu_backend/pkg/idutil/converter_test.go
```

#### Phase 2: 逐模块迁移

**对于每个模块**：

1. **Models 层**：确保使用 `primitive.ObjectID`
2. **Repository 层**：确保接口使用 `ObjectID`
3. **Service 层**：修改接口为 `string`，内部转换
4. **API 层**：DTO 使用 `string`

**示例迁移**：

```go
// ❌ 迁移前（混用情况）
type BookService interface {
    GetBook(id string) (*Book, error)              // 返回 Model
    GetBookByID(oid primitive.ObjectID) (*Book, error) // 混乱
}

// ✅ 迁移后
type BookService interface {
    GetBook(ctx context.Context, bookID string) (*dto.BookDTO, error)
}

func (s *serviceImpl) GetBook(ctx context.Context, bookID string) (*dto.BookDTO, error) {
    // 1. 转换 ID
    oid, err := idutil.ParseObjectID(bookID)
    if err != nil {
        return nil, err
    }
    
    // 2. 调用 Repository（ObjectID）
    book, err := s.repo.FindByID(ctx, oid)
    if err != nil {
        return nil, err
    }
    
    // 3. 转换为 DTO
    return &dto.BookDTO{
        ID:    idutil.ToHex(book.ID),
        Title: book.Title,
        // ...
    }, nil
}
```

### 迁移检查清单

- [ ] 所有 Model 的 ID 字段是 `primitive.ObjectID`
- [ ] 所有 Repository 接口参数是 `ObjectID`
- [ ] 所有 Service 接口参数/返回值是 `string`
- [ ] 所有 API DTO 字段是 `string`
- [ ] 不存在 `string` 和 `ObjectID` 的直接混用
- [ ] 所有转换都通过 `idutil` 包完成

## 常见错误模式

### ❌ 错误模式 1: Model 层使用 string

```go
// 错误：Model 层不应该使用 string
type Book struct {
    ID       string `bson:"_id"` // ❌
    AuthorID string `bson:"author_id"` // ❌
}
```

### ❌ 错误模式 2: Repository 层接收 string

```go
// 错误：Repository 不应该处理转换
func (r *repo) FindByID(ctx context.Context, id string) (*Book, error) {
    oid, _ := primitive.ObjectIDFromHex(id) // ❌ 转换应该在 Service 层
    // ...
}
```

### ❌ 错误模式 3: Service 层返回 ObjectID

```go
// 错误：Service 对外应该暴露 string
func (s *service) GetBook(id string) (*Book, error) { // ❌ 返回 Model
    // ...
}
```

### ❌ 错误模式 4: 中间层混用

```go
// 错误：不要在中间层来回转换
func SomeFunction(id string) {
    oid, _ := primitive.ObjectIDFromHex(id)
    str := oid.Hex()
    oid2, _ := primitive.ObjectIDFromHex(str) // ❌ 毫无意义的转换链
}
```

## 单元测试示例

```go
// pkg/idutil/converter_test.go
package idutil

import (
    "testing"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseObjectID(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    primitive.ObjectID
        wantErr bool
    }{
        {
            name:    "valid hex",
            input:   "507f1f77bcf86cd799439011",
            want:    primitive.MustObjectIDFromHex("507f1f77bcf86cd799439011"),
            wantErr: false,
        },
        {
            name:    "empty string",
            input:   "",
            want:    primitive.NilObjectID,
            wantErr: true,
        },
        {
            name:    "invalid hex",
            input:   "invalid",
            want:    primitive.NilObjectID,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseObjectID(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseObjectID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ParseObjectID() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestToHex(t *testing.T) {
    oid := primitive.NewObjectID()
    got := ToHex(oid)
    want := oid.Hex()
    
    if got != want {
        t.Errorf("ToHex() = %v, want %v", got, want)
    }
}

func TestRoundTrip(t *testing.T) {
    original := primitive.NewObjectID()
    
    // ObjectID → string → ObjectID
    hex := ToHex(original)
    parsed, err := ParseObjectID(hex)
    if err != nil {
        t.Fatalf("ParseObjectID() error = %v", err)
    }
    
    if parsed != original {
        t.Errorf("RoundTrip failed: got %v, want %v", parsed, original)
    }
}
```

## 相关文档

- [模型一致性修复指南](./model-consistency-fix-guide.md)
- [Qingyu_backend/models/writer/base/identifiers.go](../../Qingyu_backend/models/writer/base/identifiers.go)
