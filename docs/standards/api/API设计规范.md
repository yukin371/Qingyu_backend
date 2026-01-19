# API设计规范

**版本**: v3.0
**更新**: 2026-01-09
**状态**: ✅ 正式实施

**变更内容**：
- 新增第十一章：接口化设计规范
- 新增第十二章：API层测试规范

---

## 一、RESTful设计

### 1.1 URL规范

**基本规则**：
- ✅ 使用名词而非动词
- ✅ 使用复数形式表示资源集合
- ✅ 使用小写字母和连字符
- ❌ 避免深层嵌套（最多3层）

**URL结构**：
```
/api/v1/{resource}/{id}/{sub-resource}/{sub-id}
```

**示例**：
```
GET    /api/v1/books                    # 获取书籍列表
GET    /api/v1/books/{id}               # 获取特定书籍
POST   /api/v1/books                    # 创建书籍
PUT    /api/v1/books/{id}               # 更新书籍
DELETE /api/v1/books/{id}               # 删除书籍

GET    /api/v1/books/{id}/chapters      # 获取书籍的章节
POST   /api/v1/books/{id}/chapters      # 为书籍创建章节
```

### 1.2 HTTP方法

| 方法 | 用途 | 幂等性 | 示例 |
|------|------|--------|------|
| GET | 获取资源 | ✅ | 获取书籍列表 |
| POST | 创建资源 | ❌ | 创建新书 |
| PUT | 完整更新 | ✅ | 更新书籍全部信息 |
| PATCH | 部分更新 | ❌ | 更新书籍状态 |
| DELETE | 删除资源 | ✅ | 删除书籍 |

### 1.3 状态码

**成功响应**：
- `200 OK` - 请求成功
- `201 Created` - 资源创建成功
- `204 No Content` - 成功但无返回内容

**客户端错误**：
- `400 Bad Request` - 参数错误
- `401 Unauthorized` - 未认证
- `403 Forbidden` - 无权限
- `404 Not Found` - 资源不存在
- `409 Conflict` - 资源冲突
- `422 Unprocessable Entity` - 语义错误

**服务器错误**：
- `500 Internal Server Error` - 服务器错误
- `502 Bad Gateway` - 网关错误
- `503 Service Unavailable` - 服务不可用

---

## 二、请求规范

### 2.1 请求头

```
Content-Type: application/json
Authorization: Bearer {token}
X-Request-ID: {unique-id}
```

### 2.2 查询参数

**分页参数**：
```
page         页码，从1开始
page_size    每页数量，默认20，最大100
sort         排序字段
order        排序方向：asc/desc
```

**过滤参数**：
```
filter[field]    字段过滤
search           全文搜索
created_at_start 创建时间起始
created_at_end   创建时间结束
```

**示例**：
```
GET /api/v1/books?page=1&page_size=20&sort=created_at&order=desc&filter[status]=published&search=玄幻
```

### 2.3 路径参数

- 使用有意义的标识符
- 支持UUID和数字ID
- 使用下划线命名

```
/api/v1/books/{book_id}
/api/v1/users/{user_id}/books/{book_id}
```

### 2.4 请求体

- 使用驼峰命名法（camelCase）
- 必填参数明确标注
- 提供参数类型和验证规则

**示例**：
```json
{
  "title": "玄幻小说",
  "author": "张三",
  "description": "这是一部玄幻小说",
  "tags": ["玄幻", "修真"],
  "status": "published"
}
```

---

## 三、响应规范

### 3.1 统一响应结构

**成功响应**：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {},
  "timestamp": "2026-01-08T10:00:00Z",
  "request_id": "uuid"
}
```

**分页响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  },
  "timestamp": "2026-01-08T10:00:00Z",
  "request_id": "uuid"
}
```

**错误响应**：
```json
{
  "code": 400,
  "message": "参数错误",
  "error": {
    "type": "validation_error",
    "details": [
      {
        "field": "email",
        "message": "邮箱格式不正确"
      }
    ]
  },
  "timestamp": "2026-01-08T10:00:00Z",
  "request_id": "uuid"
}
```

### 3.2 响应方法

**使用统一响应方法**：
```go
// api/v1/shared/response.go
func Success(c *gin.Context, code int, message string, data interface{})
func Paginated(c *gin.Context, items interface{}, total int64, page, pageSize int, message string)
func BadRequest(c *gin.Context, message string, details string)
func Unauthorized(c *gin.Context, message string)
func Forbidden(c *gin.Context, message string)
func NotFound(c *gin.Context, message string)
func InternalError(c *gin.Context, message string, err error)
```

---

## 四、认证授权

### 4.1 认证方式

**JWT Bearer Token**：
```
Authorization: Bearer {access_token}
```

**Token刷新**：
```
POST /api/v1/auth/refresh
{
  "refresh_token": "..."
}
```

### 4.2 权限控制

#### 4.2.1 角色系统（RBAC）

**角色定义**（支持多角色）：
- `reader` - 读者（所有用户默认角色）
- `author` - 作者（可发布作品的认证作者）
- `admin` - 管理员（系统管理权限）

**角色权限矩阵**：
| 功能 | reader | author | admin |
|------|--------|--------|-------|
| 阅读作品 | ✅ | ✅ | ✅ |
| 发表评论 | ✅ | ✅ | ✅ |
| 点赞/收藏 | ✅ | ✅ | ✅ |
| 关注作者 | ✅ | ✅ | ✅ |
| 创建书单 | ✅ | ✅ | ✅ |
| 发布作品 | ❌ | ✅ | ✅ |
| 管理作品 | ❌ | 仅自己 | ✅ |
| 查看收入 | ❌ | 仅自己 | ✅ |
| 用户管理 | ❌ | ❌ | ✅ |
| 内容审核 | ❌ | ❌ | ✅ |

**多角色示例**：
```json
// 普通读者
{
  "roles": ["reader"]
}

// 认证作者（同时拥有读者权限）
{
  "roles": ["reader", "author"]
}

// 管理员（拥有所有权限）
{
  "roles": ["reader", "author", "admin"]
}
```

#### 4.2.2 VIP会员等级

**VIP等级**（独立于角色，决定内容访问范围）：
- Level 0 - 非VIP（默认）
- Level 1 - 基础VIP
- Level 2 - VIP Plus
- Level 3 - VIP Pro
- Level 4 - VIP Ultra
- Level 5 - 超级VIP

**VIP权益示例**：
- 提前阅读：VIP Level 3可提前7天阅读新章节
- 专属作品：VIP Level 4可访问专属VIP作品区
- 阅读折扣：VIP Level 2+享受8折购书优惠

**重要**：VIP是会员等级（access level），不是功能角色（functional role）

#### 4.2.3 资源级权限

**所有权原则**：
- 用户只能操作自己的资源
- 作者只能编辑自己的作品
- 管理员可操作所有资源

**权限检查流程**：
```
1. 认证检查（JWT验证）
2. 角色检查（是否拥有该角色）
3. 资源所有权检查（Service层验证）
4. VIP等级检查（内容访问限制）
```

#### 4.2.4 公共端点

**无需认证的端点**：
```
GET  /api/v1/books              # 获取书籍列表
GET  /api/v1/books/{id}         # 获取书籍详情
GET  /api/v1/recommendations    # 获取推荐
GET  /api/v1/categories         # 获取分类
```

**需要认证的端点**：
```
POST /api/v1/comments           # 发表评论
POST /api/v1/likes              # 点赞
POST /api/v1/collections        # 收藏
```

### 4.3 认证流程

#### 4.3.1 未认证用户访问流程

```
公共端点请求
  ↓
无需JWT Token
  ↓
返回公共数据（带访问限制）
```

**示例**：
- 浏览书籍列表（无个性化推荐）
- 查看书籍详情（部分章节预览）
- 阅读推荐内容（热门/最新）

#### 4.3.2 已认证用户访问流程

```
请求携带JWT Token
  ↓
JWT中间件验证
  ↓
解析用户信息（userId, roles, vipLevel）
  ↓
角色中间件检查（如需要）
  ↓
权限中间件检查（如需要）
  ↓
Service层资源所有权验证
  ↓
返回完整数据
```

#### 4.3.3 权限检查示例

**示例1：发布作品**（需要author角色）
```go
// Middleware: RequireRole("author")
// Service: 验证userId是否拥有author角色
```

**示例2：更新作品**（需要author角色 + 资源所有权）
```go
// Middleware: RequireRole("author")
// Service: 验证作品.AuthorID == userId
```

**示例3：阅读VIP章节**（需要VIP等级）
```go
// Middleware: RequireVIPLevel(3)
// Service: 验证用户VIP等级 >= 章节要求等级
```

---

## 五、版本控制

### 5.1 版本策略

**URL路径版本控制**：
```
/api/v1/books
/api/v2/books  # 未来版本
```

**版本规则**：
- 破坏性变更：升级主版本
- 新增功能：升级次版本
- Bug修复：升级修订号

### 5.2 兼容性

- 新版本保持向后兼容
- 废弃功能提前通知
- 提供迁移指南

---

## 六、错误处理

### 6.1 错误分类

**客户端错误（4xx）**：
- 参数验证失败
- 认证失败
- 权限不足
- 资源不存在
- 资源冲突

**服务器错误（5xx）**：
- 内部错误
- 数据库错误
- 第三方服务错误
- 系统过载

### 6.2 错误信息

**提供清晰信息**：
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": {
    "type": "VALIDATION_ERROR",
    "field": "email",
    "message": "邮箱格式不正确",
    "hint": "请输入有效的邮箱地址"
  }
}
```

**不泄露敏感信息**：
- ❌ 不返回数据库错误详情
- ❌ 不返回内部堆栈
- ✅ 记录详细日志供排查

---

## 七、性能优化

### 7.1 缓存策略

**HTTP缓存**：
```
Cache-Control: public, max-age=3600
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
```

**数据缓存**：
- 热点数据缓存
- 合理设置过期时间
- 缓存失效策略

### 7.2 分页优化

**默认分页**：
- page_size默认20
- 最大100
- 使用cursor分页处理大数据集

### 7.3 请求优化

**减少请求次数**：
- 批量操作接口
- 字段过滤（投影）
- 数据聚合接口

---

## 八、安全规范

### 8.1 输入验证

**严格验证所有输入**：
- 类型检查
- 长度限制
- 格式验证
- SQL注入防护
- XSS防护

### 8.2 敏感数据

**不记录敏感信息**：
- ❌ 密码
- ❌ Token完整信息
- ❌ 个人隐私信息

**传输加密**：
- ✅ HTTPS
- ✅ 敏感字段加密

---

## 九、文档规范

### 9.1 Swagger注解

**必须包含**：
```go
// @Summary      创建书籍
// @Description  创建新的书籍
// @Tags         书籍管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body BookCreateRequest true "书籍信息"
// @Success      201  {object}  APIResponse{data=Book}
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /books [post]
func (api *BookAPI) CreateBook(c *gin.Context) {
    // ...
}
```

### 9.2 接口文档

**更新流程**：
1. 代码变更时同步更新注解
2. 生成最新Swagger文档
3. 团队Review接口变更

---

## 十、最佳实践

### 10.1 设计原则

✅ **推荐**：
- RESTful风格
- 统一的响应格式
- 清晰的错误信息
- 完善的API文档
- 版本控制

❌ **避免**：
- 动词作为URL
- 深层嵌套
- 不一致的命名
- 过度设计
- 忽略错误处理

### 10.2 性能考虑

- 合理使用缓存
- 避免过度查询
- 使用批量操作
- 限制返回数据量
- 异步处理耗时操作

### 10.3 可维护性

- 清晰的接口命名
- 统一的错误处理
- 完善的日志记录
- 版本兼容性
- 详细的文档

**权限控制最佳实践**：

✅ **推荐**：
- 使用中间件统一处理角色验证
- 资源所有权检查在Service层进行
- VIP等级检查与业务逻辑分离
- 权限拒绝时返回明确的错误信息
- 记录权限检查失败的审计日志

❌ **避免**：
- 在API层硬编码用户ID
- 绕过中间件直接检查角色
- 返回模糊的权限错误（如"权限不足"而不说明原因）
- 在客户端进行权限验证（永远不可信）

**多角色系统设计原则**：
```go
// ✅ 正确：使用中间件检查角色
router.POST("/books", middleware.RequireRole("author"), api.CreateBook)

// ✅ 正确：Service层验证资源所有权
func (s *BookService) UpdateBook(ctx context.Context, userID, bookID string) error {
    book, err := s.repo.GetByID(ctx, bookID)
    if err != nil {
        return err
    }
    if book.AuthorID != userID && !user.HasRole("admin") {
        return ErrNoPermission
    }
    // 更新逻辑...
}

// ❌ 错误：在API层硬编码权限检查
func (api *BookAPI) UpdateBook(c *gin.Context) {
    userID := c.GetString("user_id")
    // 直接检查数据库，绕过中间件
}
```

---

## 十一、接口化设计规范

### 11.1 依赖倒置原则在API层的应用

**核心原则**：
- ✅ API层依赖Service接口，而非具体实现
- ✅ 通过构造函数注入依赖
- ✅ 使用接口实现解耦和可测试性

**代码对比**：
```go
// ❌ 错误：依赖具体实现
type BookstoreAPI struct {
    service *bookstore.BookstoreService
}

// ✅ 正确：依赖接口
type BookstoreAPI struct {
    service interfaces.BookstoreService
}
```

### 11.2 服务接口定义规范

**接口定义位置**：
```
service/interfaces/
├── comment_service_interface.go
├── like_service_interface.go
├── bookstore_service_interface.go
├── reader_service_interface.go
└── ...
```

**命名规范**：
- 接口名：`{模块名}Service`
- 文件名：`{模块名}_service_interface.go`

**接口定义模板**：
```go
package interfaces

import (
    "context"
    "Qingyu_backend/models/{module}"
)

type {Module}Service interface {
    // CRUD操作
    GetByID(ctx context.Context, id string) (*models.{Entity}, error)
    List(ctx context.Context, page, size int) ([]*models.{Entity}, int64, error)

    // 业务操作
    {BusinessAction}(ctx context.Context, ...) error
}
```

### 11.3 API层结构设计模式

**标准API结构**：
```go
type {Module}API struct {
    {service} interfaces.{Module}Service
}

func New{Module}API(
    {service} interfaces.{Module}Service,
) *{Module}API {
    return &{Module}API{
        {service}: {service},
    }
}
```

### 11.4 可测试性设计规范

**测试覆盖率要求**：
- API层整体：≥ 80%
- 核心业务逻辑：100%
- 错误处理：≥ 70%

**Mock对象创建**：
- 使用 `testify/mock`
- 实现所有接口方法
- 遵循 Given-When-Then 模式

---

## 十二、API层测试规范

### 12.1 测试策略

**测试分布**：
- 单元测试：60% - 使用Mock对象测试单个方法
- 集成测试：30% - 测试多个组件协作
- E2E测试：10% - 端到端业务流程测试

### 12.2 测试文件结构

**文件组织**：
```
api/v1/social/
├── comment_api.go
└── comment_api_test.go
```

**测试命名规范**：
```go
// 测试文件命名：{module}_api_test.go
// 测试函数命名：Test{API名}_{方法名}_{场景}

func TestCommentAPI_CreateComment_Success(t *testing.T)
func TestCommentAPI_CreateComment_MissingBookID(t *testing.T)
func TestCommentAPI_CreateComment_Unauthorized(t *testing.T)
```

### 12.3 测试模式：Given-When-Then

**标准测试结构**：
```go
func TestCommentAPI_CreateComment_Success(t *testing.T) {
    // Given - 准备测试数据
    mockService := new(MockCommentService)
    userID := primitive.NewObjectID().Hex()
    router := setupTestRouter(mockService, userID)

    reqBody := map[string]interface{}{
        "book_id": "book123",
        "content": "这是一本非常好的书！",
        "rating":  5,
    }

    expectedComment := &community.Comment{}
    expectedComment.ID = primitive.NewObjectID().Hex()
    expectedComment.AuthorID = userID
    expectedComment.Content = "这是一本非常好的书！"
    expectedComment.Rating = 5

    mockService.On("PublishComment", mock.Anything, userID, "book123",
        "这是一本非常好的书！", 5).Return(expectedComment, nil)

    jsonBody, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/api/v1/reader/comments",
        bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")

    // When - 执行操作
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Then - 验证结果
    assert.Equal(t, http.StatusCreated, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, float64(http.StatusCreated), response["code"])
    assert.NotNil(t, response["data"])

    mockService.AssertExpectations(t)
}
```

### 12.4 Mock对象创建

**Mock服务实现**：
```go
type MockCommentService struct {
    mock.Mock
}

func (m *MockCommentService) PublishComment(
    ctx context.Context,
    userID, bookID, chapterID, content string,
    rating int,
) (*community.Comment, error) {
    args := m.Called(ctx, userID, bookID, chapterID, content, rating)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*community.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentList(
    ctx context.Context,
    bookID string,
    sortBy string,
    page, size int,
) ([]*community.Comment, int64, error) {
    args := m.Called(ctx, bookID, sortBy, page, size)
    if args.Get(0) == nil {
        return nil, args.Get(1).(int64), args.Error(2)
    }
    return args.Get(0).([]*community.Comment), args.Get(1).(int64), args.Error(2)
}
```

### 12.5 测试辅助函数

**设置测试路由**：
```go
func setupTestRouter(commentService interfaces.CommentService, userID string) *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()

    // 添加middleware来设置userId（用于需要认证的端点）
    r.Use(func(c *gin.Context) {
        if userID != "" {
            c.Set("user_id", userID)
        }
        c.Next()
    })

    api := socialAPI.NewCommentAPI(commentService)

    v1 := r.Group("/api/v1/reader")
    {
        v1.POST("/comments", api.CreateComment)
        v1.GET("/comments", api.GetCommentList)
        v1.GET("/comments/:id", api.GetCommentDetail)
        v1.PUT("/comments/:id", api.UpdateComment)
        v1.DELETE("/comments/:id", api.DeleteComment)
    }

    return r
}
```

### 12.6 测试覆盖场景

**必须测试的场景**：
1. ✅ 正常流程（Success）
2. ✅ 参数验证失败（Missing Required Fields）
3. ✅ 权限验证（Unauthorized）
4. ✅ 资源不存在（Not Found）
5. ✅ 业务规则冲突（Conflict）
6. ✅ 服务错误（Service Error）

**示例测试用例**：
```go
// 1. 正常流程
func TestCommentAPI_CreateComment_Success(t *testing.T)

// 2. 参数验证
func TestCommentAPI_CreateComment_MissingBookID(t *testing.T)
func TestCommentAPI_CreateComment_ContentTooShort(t *testing.T)

// 3. 权限验证
func TestCommentAPI_CreateComment_Unauthorized(t *testing.T)

// 4. 资源不存在
func TestCommentAPI_GetCommentDetail_NotFound(t *testing.T)

// 5. 权限检查
func TestCommentAPI_DeleteComment_NoPermission(t *testing.T)

// 6. 批量操作
func TestCommentAPI_BatchDeleteComments_Success(t *testing.T)
```

---

**相关文档**：
- [架构设计规范](../architecture/架构设计规范.md)
- [路由层设计规范](../architecture/路由层设计规范.md)
- [测试架构规范](../testing/测试架构规范.md)
