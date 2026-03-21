# API 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 职责边界与依赖关系

### 1.1 职责定义

API 层是后端系统的**请求处理层**，负责：

1. **HTTP 请求处理**：接收 HTTP 请求，解析参数
2. **参数验证**：使用 `binding` 标签进行请求体验证
3. **调用 Service 层**：将请求转发给 Service 层处理业务逻辑
4. **响应封装**：将 Service 层返回的数据转换为标准响应格式
5. **错误处理**：捕获并转换错误为标准 HTTP 响应

### 1.2 与上下层级的交互边界

```
┌─────────────────────────────────────────────────────────┐
│                     客户端请求                           │
│                  (HTTP/JSON)                            │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                      API 层                             │
│  ┌─────────────────────────────────────────────────────┤
│  │ 职责：                                              │
│  │ - 路由匹配                                          │
│  │ - 参数解析与验证                                    │
│  │ - 调用 Service                                      │
│  │ - DTO 转换                                          │
│  │ - 响应封装                                          │
│  │ - 错误处理                                          │
│  └─────────────────────────────────────────────────────┤
│  禁止：                                                  │
│  │ - 直接访问数据库                                    │
│  │ - 编写业务逻辑                                      │
│  │ - 跨模块调用其他 API                                │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Service 层                           │
│              (业务逻辑编排)                             │
└─────────────────────────────────────────────────────────┘
```

### 1.3 依赖关系

```go
// API 层允许的依赖
import (
    "github.com/gin-gonic/gin"
    "Qingyu_backend/pkg/response"       // 响应封装
    "Qingyu_backend/pkg/logger"         // 日志
    "Qingyu_backend/models/dto"         // DTO 类型
    "Qingyu_backend/service/xxx"        // Service 接口
    "Qingyu_backend/models/xxx"         // 领域模型（仅用于转换）
)

// API 层禁止的依赖
import (
    "Qingyu_backend/repository/xxx"     // ❌ 禁止直接访问 Repository
    "go.mongodb.org/mongo-driver"       // ❌ 禁止直接操作数据库
)
```

---

## 2. 命名与代码规范

### 2.1 文件命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| API 处理器 | `{模块}_{功能}_api.go` | `bookstore_api.go`, `user_auth_api.go` |
| DTO 定义 | `{模块}_dto.go` | `auth_dto.go`, `user_dto.go` |
| 转换器 | `{模块}_converter.go` | `bookstore_converter.go` |
| 类型定义 | `types.go` | `types.go` |
| 路由注册 | `routes.go` | `routes.go` |
| 测试文件 | `{文件名}_test.go` | `bookstore_api_test.go` |

### 2.2 结构体命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| API 处理器 | `{模块}API` 或 `{模块}Handler` | `BookstoreAPI`, `UserHandler` |
| 请求 DTO | `{操作}{资源}Request` | `CreateBookRequest`, `UpdateProfileRequest` |
| 响应 DTO | `{资源}Response` | `BookResponse`, `UserProfileResponse` |
| 列表响应 | `{资源}ListResponse` | `BookListResponse` |
| 分页响应 | 使用通用 `PaginatedResponse` | - |

### 2.3 方法命名规范

| HTTP 方法 | 前缀 | 示例 |
|-----------|------|------|
| GET (单个) | `Get` | `GetBookByID`, `GetUserProfile` |
| GET (列表) | `List` 或 `Get` | `ListBooks`, `GetBooks` |
| POST | `Create` | `CreateBook`, `CreateProject` |
| PUT | `Update` | `UpdateBook`, `UpdateProfile` |
| DELETE | `Delete` | `DeleteBook`, `DeleteComment` |
| 特殊操作 | 动词 | `Login`, `Logout`, `Search` |

### 2.4 目录组织规范

```
api/v1/
├── admin/              # 管理后台 API
│   ├── user_admin_api.go
│   ├── permission_api.go
│   └── audit_admin_api.go
├── auth/               # 认证相关 API
│   └── auth_api.go
├── bookstore/          # 书城 API
│   ├── bookstore_api.go
│   ├── bookstore_converter.go
│   └── chapter_catalog_api.go
├── reader/             # 阅读器 API
│   ├── bookmark_api.go
│   ├── chapter_api.go
│   └── types.go
├── shared/             # 共享工具
│   ├── request_validator.go
│   ├── types.go
│   └── response.go (概念上的)
├── social/             # 社交 API
│   ├── comment_api.go
│   ├── like_api.go
│   └── follow_api.go
└── writer/             # 作家 API
    ├── project_api.go
    ├── document_api.go
    └── types.go
```

---

## 3. 设计模式与最佳实践

### 3.1 API 处理器结构模式

```go
// API 处理器结构
type BookstoreAPI struct {
    service       bookstoreService.BookstoreService  // 业务服务
    searchService *search.SearchService              // 可选：搜索服务
    logger        *logger.Logger                     // 日志
}

// 构造函数：依赖注入
func NewBookstoreAPI(
    service bookstoreService.BookstoreService,
    searchService *search.SearchService,
    logger *logger.Logger,
) *BookstoreAPI {
    return &BookstoreAPI{
        service:       service,
        searchService: searchService,
        logger:        logger,
    }
}
```

### 3.2 请求处理流程

```
┌─────────────────────────────────────────────────────────┐
│                   请求处理标准流程                        │
└─────────────────────────────────────────────────────────┘
                          │
         ┌────────────────┼────────────────┐
         ▼                ▼                ▼
   ┌──────────┐    ┌──────────┐    ┌──────────┐
   │ 1. 参数   │    │ 2. 业务   │    │ 3. 响应   │
   │    解析   │───▶│    处理   │───▶│    封装   │
   └──────────┘    └──────────┘    └──────────┘
         │                │                │
         ▼                ▼                ▼
   - 绑定请求体     - 调用 Service    - 成功响应
   - 验证参数       - 错误处理        - 分页响应
   - 参数预处理     - DTO 转换        - 错误响应
```

### 3.3 响应封装模式

使用 `pkg/response` 包的统一响应方法：

```go
// 成功响应
response.SuccessWithMessage(c, "操作成功", data)
response.Success(c, data)  // 无消息

// 分页响应
response.Paginated(c, items, total, page, size, "获取成功")

// 错误响应
response.BadRequest(c, "参数错误", "详细说明")
response.Unauthorized(c, "未认证")
response.Forbidden(c, "权限不足")
response.NotFound(c, "资源不存在")
response.ErrorJSON(c, http.StatusInternalServerError, "服务器错误")
```

### 3.4 参数验证模式

#### 使用 binding 标签

```go
type CreateBookRequest struct {
    Title       string   `json:"title" binding:"required,min=1,max=200"`
    Author      string   `json:"author" binding:"required,max=100"`
    CategoryID  string   `json:"categoryId" binding:"required"`
    Tags        []string `json:"tags" binding:"max=10,dive,min=1,max=50"`
    Price       float64  `json:"price" binding:"gte=0"`
    IsFree      bool     `json:"isFree"`
}
```

#### 手动验证

```go
func (api *BookstoreAPI) GetBookByID(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.BadRequest(c, "参数错误", "书籍ID不能为空")
        return
    }

    // 验证 ObjectID 格式
    if _, err := primitive.ObjectIDFromHex(id); err != nil {
        response.BadRequest(c, "参数错误", "书籍ID格式无效")
        return
    }

    // 继续处理...
}
```

### 3.5 分页参数规范

```go
// 标准分页参数获取
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

// 参数边界校验
if page < 1 {
    page = 1
}
if size < 1 || size > 100 {
    size = 20  // 默认值
}
```

### 3.6 错误处理模式

```go
func (api *BookstoreAPI) GetBookByID(c *gin.Context) {
    // ... 参数验证 ...

    book, err := api.service.GetBookByID(ctx, id)
    if err != nil {
        // 1. 特定错误处理
        if err.Error() == "book not found" || err.Error() == "book not available" {
            response.SuccessWithMessage(c, "书籍不存在或不可用", nil)
            return
        }

        // 2. 通用错误处理：交给中间件
        c.Error(err)
        return
    }

    // 3. 成功响应
    bookDTO := ToBookDTO(book)
    response.SuccessWithMessage(c, "获取书籍详情成功", bookDTO)
}
```

### 3.7 反模式警示

❌ **禁止事项**：

```go
// ❌ 禁止：直接操作数据库
func (api *BookstoreAPI) GetBook(c *gin.Context) {
    collection := api.db.Collection("books")
    collection.FindOne(...)
}

// ❌ 禁止：在 API 层编写业务逻辑
func (api *BookstoreAPI) CreateBook(c *gin.Context) {
    if book.Price > 100 {
        // 价格计算逻辑应该在 Service 层
    }
}

// ❌ 禁止：直接返回 Model
func (api *BookstoreAPI) GetBook(c *gin.Context) {
    book, _ := api.service.GetBookByID(ctx, id)
    c.JSON(200, book)  // 应该转换为 DTO
}

// ❌ 禁止：跨模块调用其他 API
func (api *BookstoreAPI) GetBook(c *gin.Context) {
    userAPI.GetUser(...)  // 应该通过 Service 层交互
}
```

---

## 4. 接口与契约规范

### 4.1 统一响应格式

#### 成功响应

```json
{
  "code": 200,
  "message": "操作成功",
  "data": { ... }
}
```

#### 分页响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [ ... ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "参数错误"
}
```

#### 验证错误响应

```json
{
  "code": 400,
  "message": "请求参数验证失败",
  "errors": {
    "title": "标题不能为空",
    "price": "价格必须大于等于0"
  }
}
```

### 4.2 HTTP 状态码规范

| 状态码 | 场景 | 使用方法 |
|--------|------|----------|
| 200 | 成功（GET/PUT/DELETE） | `response.SuccessWithMessage(c, ...)` |
| 201 | 创建成功 | `response.Success(c, data)` (配合 HTTP 201) |
| 400 | 请求参数错误 | `response.BadRequest(c, ...)` |
| 401 | 未认证 | `response.Unauthorized(c, ...)` |
| 403 | 禁止访问 | `response.Forbidden(c, ...)` |
| 404 | 资源不存在 | `response.NotFound(c, ...)` |
| 500 | 服务器错误 | `response.ErrorJSON(c, 500, ...)` |

### 4.3 业务错误码规范

定义在 `pkg/response/codes.go`：

```go
const (
    // 成功
    CodeSuccess = 0

    // 通用客户端错误 (1000-1999)
    CodeParamError       = 1001
    CodeUnauthorized     = 1002
    CodeForbidden        = 1003
    CodeNotFound         = 1004

    // 用户相关错误 (2000-2999)
    CodeUserNotFound      = 2001
    CodeInvalidCredential = 2002

    // 业务逻辑错误 (3000-3999)
    CodeBookNotFound      = 3001
    CodeChapterNotFound   = 3002

    // 服务端错误 (5000-5999)
    CodeInternalError     = 5000
)
```

### 4.4 Swagger 注释规范

```go
// GetBookByID 根据ID获取书籍详情
//
//	@Summary		获取书籍详情
//	@Description	根据书籍ID获取书籍的详细信息
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id} [get]
func (api *BookstoreAPI) GetBookByID(c *gin.Context) {
    // ...
}
```

---

## 5. 测试策略

### 5.1 单元测试编写指南

#### 测试文件结构

```go
// bookstore_api_test.go
package bookstore

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock Service
type MockBookstoreService struct {
    mock.Mock
}

func (m *MockBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*bookstore.Book), args.Error(1)
}
```

#### 测试用例示例

```go
func TestBookstoreAPI_GetBookByID_Success(t *testing.T) {
    // 1. 准备
    gin.SetMode(gin.TestMode)
    mockService := new(MockBookstoreService)
    api := NewBookstoreAPI(mockService, nil, logger.NewNopLogger())

    // 设置 Mock 返回值
    mockService.On("GetBookByID", mock.Anything, "507f1f77bcf86cd799439011").
        Return(&bookstore.Book{Title: "测试书籍"}, nil)

    // 2. 执行
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "507f1f77bcf86cd799439011"}}

    api.GetBookByID(c)

    // 3. 验证
    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}

func TestBookstoreAPI_GetBookByID_InvalidID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockService := new(MockBookstoreService)
    api := NewBookstoreAPI(mockService, nil, logger.NewNopLogger())

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: ""}}  // 空 ID

    api.GetBookByID(c)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

### 5.2 测试覆盖率要求

| 测试类型 | 覆盖率要求 |
|----------|------------|
| 单元测试 | ≥ 70% |
| 关键路径 | 100% |
| 错误处理 | ≥ 80% |

### 5.3 测试场景清单

- [ ] 正常流程：成功获取/创建/更新/删除
- [ ] 参数验证：必填字段缺失、格式错误
- [ ] 边界条件：分页参数、字符串长度
- [ ] 错误处理：Service 层返回错误
- [ ] 权限检查：未认证、无权限

---

## 6. 完整代码示例

### 6.1 完整 API 处理器示例

```go
// bookstore_api.go
package bookstore

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.uber.org/zap"

    "Qingyu_backend/models/bookstore"
    "Qingyu_backend/pkg/logger"
    "Qingyu_backend/pkg/response"
    bookstoreService "Qingyu_backend/service/bookstore"
)

// BookstoreAPI 书城 API 处理器
type BookstoreAPI struct {
    service       bookstoreService.BookstoreService
    searchService *search.SearchService
    logger        *logger.Logger
}

// NewBookstoreAPI 创建书城 API 实例
func NewBookstoreAPI(
    service bookstoreService.BookstoreService,
    searchService *search.SearchService,
    logger *logger.Logger,
) *BookstoreAPI {
    return &BookstoreAPI{
        service:       service,
        searchService: searchService,
        logger:        logger,
    }
}

// GetHomepage 获取首页数据
//
//	@Summary		获取书城首页数据
//	@Description	获取书城首页的Banner、推荐书籍、精选书籍、分类等数据
//	@Tags			书城
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/homepage [get]
func (api *BookstoreAPI) GetHomepage(c *gin.Context) {
    data, err := api.service.GetHomepageData(c.Request.Context())
    if err != nil {
        c.Error(err)
        return
    }

    response.SuccessWithMessage(c, "获取首页数据成功", data)
}

// GetBooks 获取书籍列表
//
//	@Summary		获取书籍列表
//	@Description	获取所有书籍列表，支持分页和关键词搜索
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"页码"	default(1)
//	@Param			size	query		int	false	"每页数量"	default(20)
//	@Param			q		query		string	false	"搜索关键词"
//	@Success		200	{object}	response.PaginatedResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books [get]
func (api *BookstoreAPI) GetBooks(c *gin.Context) {
    // 1. 参数解析
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    keyword := c.Query("q")

    // 2. 参数验证
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 20
    }

    // 3. 调用 Service
    var books []*bookstore.Book
    var total int64
    var err error

    if keyword != "" {
        filter := &bookstore.BookFilter{
            Keyword:   &keyword,
            SortBy:    "created_at",
            SortOrder: "desc",
            Limit:     size,
            Offset:    (page - 1) * size,
        }
        books, total, err = api.service.SearchBooksWithFilter(c.Request.Context(), filter)
    } else {
        books, total, err = api.service.GetAllBooks(c.Request.Context(), page, size)
    }

    if err != nil {
        c.Error(err)
        return
    }

    // 4. DTO 转换
    bookDTOs := ToBookDTOsFromPtrSlice(books)

    // 5. 响应封装
    response.Paginated(c, bookDTOs, total, page, size, "获取书籍列表成功")
}

// GetBookByID 根据ID获取书籍详情
//
//	@Summary		获取书籍详情
//	@Description	根据书籍ID获取书籍的详细信息
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id} [get]
func (api *BookstoreAPI) GetBookByID(c *gin.Context) {
    // 1. 参数解析
    id := c.Param("id")
    if id == "" {
        response.BadRequest(c, "参数错误", "书籍ID不能为空")
        return
    }

    // 2. 参数验证
    if _, err := primitive.ObjectIDFromHex(id); err != nil {
        response.BadRequest(c, "参数错误", "书籍ID格式无效")
        return
    }

    // 3. 调用 Service
    book, err := api.service.GetBookByID(c.Request.Context(), id)
    if err != nil {
        if err.Error() == "book not found" || err.Error() == "book not available" {
            response.SuccessWithMessage(c, "书籍不存在或不可用", nil)
            return
        }
        c.Error(err)
        return
    }

    // 4. DTO 转换
    bookDTO := ToBookDTO(book)

    // 5. 响应封装
    response.SuccessWithMessage(c, "获取书籍详情成功", bookDTO)
}
```

### 6.2 路由注册示例

```go
// routes.go
package bookstore

import (
    "github.com/gin-gonic/gin"
)

// RegisterRoutes 注册书城路由
func RegisterRoutes(r *gin.RouterGroup, api *BookstoreAPI) {
    // 书城路由组
    bookstore := r.Group("/bookstore")
    {
        // 首页
        bookstore.GET("/homepage", api.GetHomepage)

        // 书籍
        books := bookstore.Group("/books")
        {
            books.GET("", api.GetBooks)
            books.GET("/search", api.SearchBooks)
            books.GET("/recommended", api.GetRecommendedBooks)
            books.GET("/featured", api.GetFeaturedBooks)
            books.GET("/:id", api.GetBookByID)
            books.GET("/:id/similar", api.GetSimilarBooks)
            books.POST("/:id/view", api.IncrementBookView)
        }

        // 分类
        categories := bookstore.Group("/categories")
        {
            categories.GET("/tree", api.GetCategoryTree)
            categories.GET("/:id", api.GetCategoryByID)
            categories.GET("/:id/books", api.GetBooksByCategory)
        }

        // 榜单
        rankings := bookstore.Group("/rankings")
        {
            rankings.GET("/realtime", api.GetRealtimeRanking)
            rankings.GET("/weekly", api.GetWeeklyRanking)
            rankings.GET("/monthly", api.GetMonthlyRanking)
            rankings.GET("/:type", api.GetRankingByType)
        }
    }
}
```

---

## 7. 参考资料

- [API 层快速参考](../api/v1/README.md)
- [DTO 层设计说明](./layer-dto.md)
- [Service 层设计说明](./layer-service.md)
- [错误码规范](./api-status-code-standard.md)

---

*最后更新：2026-03-19*
