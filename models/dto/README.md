# DTO 层快速参考

## 职责

定义数据传输对象，用于 API 层和 Service 层之间的数据传输，隔离内部 Model 结构。

## 目录结构

```
models/dto/
├── bookstore.go           # 书城 DTO
├── user.go                # 用户 DTO
├── reader.go              # 读者 DTO
├── writer_dto.go          # 作家 DTO
├── writer_converter.go    # 作家 DTO 转换器
├── content_dto.go         # 内容 DTO
├── audit.go               # 审计 DTO
└── README.md              # 本文档
```

## DTO 与 Model 的区别

| 特性 | DTO | Model |
|------|-----|-------|
| 位置 | `models/dto/` | `models/bookstore/` |
| 用途 | API 数据传输 | 数据库映射 |
| ID 类型 | `string` | `primitive.ObjectID` |
| 时间类型 | `string` (ISO8601) 或 `time.Time` | `time.Time` |
| JSON 标签 | camelCase | camelCase |
| BSON 标签 | ❌ 无 | ✅ snake_case |
| 验证标签 | ✅ validate | ❌ 无 |

## 命名规范

| 类型 | 后缀 | 示例 |
|------|------|------|
| 请求 DTO | `Request` | `CreateBookRequest`, `LoginRequest` |
| 更新请求 | `UpdateRequest` | `UpdateBookRequest` |
| 响应 DTO | `Response` | `BookResponse`, `UserProfileResponse` |
| 列表请求 | `ListRequest` | `ListBooksRequest` |
| 列表响应 | `ListResponse` | `BookListResponse` |
| 数据传输 | `DTO` | `BookDTO` |

## 常用验证标签

```go
type CreateBookRequest struct {
    // 必填验证
    Title string `json:"title" validate:"required"`

    // 长度验证
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Summary string `json:"summary" validate:"max=500"`

    // 数值验证
    Price    float64 `json:"price" validate:"gte=0"`
    Rating   float64 `json:"rating" validate:"min=0,max=5"`
    PageSize int     `json:"pageSize" validate:"min=1,max=100"`

    // 格式验证
    Email string `json:"email" validate:"required,email"`
    URL   string `json:"url" validate:"omitempty,url"`

    // 枚举验证
    Status string `json:"status" validate:"oneof=draft ongoing completed"`

    // 切片验证
    Tags []string `json:"tags" validate:"max=10,dive,min=1,max=50"`
}
```

## 输入/输出分离

```go
// 请求 DTO（创建 - 全部必填）
type CreateBookRequest struct {
    Title  string  `json:"title" validate:"required"`
    Price  float64 `json:"price" validate:"gte=0"`
}

// 请求 DTO（更新 - 使用指针表示可选）
type UpdateBookRequest struct {
    Title  *string  `json:"title,omitempty"`
    Price  *float64 `json:"price,omitempty"`
}

// 响应 DTO
type BookResponse struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    CreatedAt string `json:"createdAt"`
}
```

## 转换器模式

```go
// Model → DTO
func ToBookDTO(book *bookstore.Book) *dto.BookDTO {
    if book == nil {
        return nil
    }

    var converter types.DTOConverter
    return &dto.BookDTO{
        ID:        converter.ModelIDToDTO(book.ID),
        CreatedAt: converter.TimeToISO8601(book.CreatedAt),
        Title:     book.Title,
    }
}

// 批量转换
func ToBookDTOsFromPtrSlice(books []*bookstore.Book) []*dto.BookDTO {
    result := make([]*dto.BookDTO, len(books))
    for i := range books {
        result[i] = ToBookDTO(books[i])
    }
    return result
}

// DTO → Model
func ToBookModel(dto *dto.BookDTO) (*bookstore.Book, error) {
    var converter types.DTOConverter
    id, err := converter.DTOIDToModel(dto.ID)
    if err != nil {
        return nil, err
    }
    return &bookstore.Book{
        IdentifiedEntity: shared.IdentifiedEntity{ID: id},
        Title:            dto.Title,
    }, nil
}
```

## 禁止事项

- ❌ 在 DTO 中使用 BSON 标签
- ❌ DTO 中包含业务逻辑
- ❌ DTO 直接使用 ObjectID 类型
- ❌ 响应 DTO 包含验证标签

## 详见

完整设计文档: [docs/standards/layer-dto.md](../../docs/standards/layer-dto.md)
