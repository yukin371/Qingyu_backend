# Models 层快速参考

## 职责

定义数据结构、MongoDB 映射、JSON 序列化、字段验证和枚举类型。

## 目录结构

```
models/
├── admin/          # 管理员相关
├── ai/             # AI服务
├── auth/           # 认证授权
├── bookstore/      # 书城（书籍、章节、分类）
├── dto/            # 数据传输对象
├── finance/        # 财务
├── messaging/      # 消息
├── notification/   # 通知
├── reader/         # 阅读（进度、书架）
├── shared/         # 共享基础类型
│   └── types/      # ID、金额、评分等
├── social/         # 社交
└── stats/          # 统计
```

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | 小写单数 | `book.go`, `chapter.go` |
| 结构体 | PascalCase | `Book`, `Chapter` |
| 字段 | PascalCase | `Title`, `AuthorID` |
| 枚举 | 实体名 + 属性 | `BookStatus` |
| 枚举常量 | 类型前缀 | `BookStatusDraft` |

## Tag 规范

```go
type Book struct {
    // 必填字段
    Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`

    // 可选字段
    Subtitle string `bson:"subtitle,omitempty" json:"subtitle,omitempty"`

    // 外键
    AuthorID primitive.ObjectID `bson:"author_id" json:"authorId"`
}
```

## 基础实体嵌入

```go
type Book struct {
    shared.IdentifiedEntity `bson:",inline"`  // ID字段
    shared.BaseEntity       `bson:",inline"`  // 时间戳、软删除

    Title string `bson:"title" json:"title"`
    // ... 业务字段
}
```

## 枚举定义

```go
type BookStatus string

const (
    BookStatusDraft     BookStatus = "draft"
    BookStatusOngoing   BookStatus = "ongoing"
    BookStatusCompleted BookStatus = "completed"
)

func (s BookStatus) IsValid() bool {
    switch s {
    case BookStatusDraft, BookStatusOngoing, BookStatusCompleted:
        return true
    default:
        return false
    }
}
```

## 过滤器模式

```go
type BookFilter struct {
    CategoryID    *string     `json:"categoryId,omitempty"`
    Status        *BookStatus `json:"status,omitempty"`
    IsRecommended *bool       `json:"isRecommended,omitempty"`
    MinPrice      *float64    `json:"minPrice,omitempty"`
    MaxPrice      *float64    `json:"maxPrice,omitempty"`
    Tags          []string    `json:"tags,omitempty"`
    Keyword       *string     `json:"keyword,omitempty"`
    SortBy        string      `json:"sortBy,omitempty"`
    Limit         int         `json:"limit,omitempty"`
    Offset        int         `json:"offset,omitempty"`
}
```

## 禁止事项

- ❌ 在模型中包含数据库操作
- ❌ 循环依赖
- ❌ 调用 Service 或 Repository
- ❌ 硬编码状态值（使用枚举）

## 详见

完整设计文档: [docs/standards/layer-models.md](../docs/standards/layer-models.md)
