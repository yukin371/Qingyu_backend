# API层测试规范

## 概述

API层负责处理HTTP请求和响应，是系统与外部交互的入口。本规范定义了API层测试的详细要求和最佳实践。

## 核心原则

### ✅ 必须使用集成测试

```go
// ✅ 正确示例：完整集成测试
func TestAuthAPI_Login(t *testing.T) {
    // Setup - 真实数据库 + HTTP服务器
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 准备测试数据
    user := helper.CreateTestUser(&users.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123", // 会被hash
    })

    // Act - 真实HTTP请求
    reqBody := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    }
    w := helper.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)

    // Assert - 验证HTTP响应
    helper.AssertSuccess(w, 200, "登录失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NotEmpty(t, resp["token"])
}
```

### ❌ 严格禁止

```go
// ❌ 错误：只测试单个handler函数
func TestAuthAPI_LoginHandler_UnitTest(t *testing.T) {
    // 问题：绕过了路由、中间件、真实请求解析
    handler := NewAuthHandler(mockService)

    req := httptest.NewRequest("POST", "/login", strings.NewReader(`{"email":"test@example.com"}`))
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    // 测试通过了，但真实请求可能失败：
    // - 路由配置错误
    // - 中间件缺失
    // - 请求解析问题
}
```

## 测试组织结构

### 文件位置

```
test/api/v1/{module}_test.go
```

示例：
```
test/api/v1/auth_test.go
test/api/v1/bookstore_test.go
test/api/v1/reader_test.go
```

### 包命名

```go
package v1_test

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "Qingyu_backend/test/integration"
)
```

## TestHelper使用

### TestHelper封装

```go
// test/integration/test_helper.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type TestHelper struct {
    T      *testing.T
    Router *gin.Engine
}

func NewTestHelper(t *testing.T, router *gin.Engine) *TestHelper {
    return &TestHelper{
        T:      t,
        Router: router,
    }
}

// DoRequest 执行HTTP请求（无认证）
func (h *TestHelper) DoRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
    h.T.Helper()

    var bodyReader *bytes.Reader
    if body != nil {
        bodyBytes, _ := json.Marshal(body)
        bodyReader = bytes.NewReader(bodyBytes)
    } else {
        bodyReader = bytes.NewReader([]byte{})
    }

    req := httptest.NewRequest(method, path, bodyReader)
    req.Header.Set("Content-Type", "application/json")

    for k, v := range headers {
        req.Header.Set(k, v)
    }

    w := httptest.NewRecorder()
    h.Router.ServeHTTP(w, req)

    return w
}

// DoAuthRequest 执行需要认证的HTTP请求
func (h *TestHelper) DoAuthRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
    h.T.Helper()

    headers := map[string]string{
        "Authorization": "Bearer " + token,
    }
    return h.DoRequest(method, path, body, headers)
}

// AssertSuccess 断言成功响应
func (h *TestHelper) AssertSuccess(w *httptest.ResponseRecorder, expectedStatus int, msg string) {
    h.T.Helper()

    assert.Equal(h.T, expectedStatus, w.Code, "%s, 状态码: %d, 响应: %s", msg, w.Code, w.Body.String())
}

// AssertError 断言错误响应
func (h *TestHelper) AssertError(w *httptest.ResponseRecorder, expectedStatus int, expectedMsg string) {
    h.T.Helper()

    assert.Equal(h.T, expectedStatus, w.Code, "期望状态码 %d", expectedStatus)

    var resp map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    require.NoError(h.T, err)

    if expectedMsg != "" {
        assert.Contains(h.T, resp["message"], expectedMsg)
    }
}

// AssertResponseStructure 断言响应结构
func (h *TestHelper) AssertResponseStructure(w *httptest.ResponseRecorder, fields []string) {
    h.T.Helper()

    var resp map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    require.NoError(h.T, err)

    for _, field := range fields {
        _, exists := resp[field]
        assert.True(h.T, exists, "响应缺少字段: %s", field)
    }
}

// CreateTestUser 创建测试用户
func (h *TestHelper) CreateTestUser(user *users.User) *users.User {
    h.T.Helper()

    if user == nil {
        user = &users.User{
            Username: "testuser_" + primitive.NewObjectID().Hex(),
            Email:    "test_" + primitive.NewObjectID().Hex() + "@example.com",
            Password: "password123",
        }
    }

    // 直接操作数据库创建用户
    db := GetTestDB(h.T)
    userRepo := repository.NewUserRepository(db)

    // Hash密码
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPassword)

    err := userRepo.Create(context.Background(), user)
    require.NoError(h.T, err)

    return user
}

// LoginTestUser 登录测试用户并返回token
func (h *TestHelper) LoginTestUser(email, password string) string {
    h.T.Helper()

    reqBody := map[string]interface{}{
        "email":    email,
        "password": password,
    }
    w := h.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)
    h.AssertSuccess(w, 200, "登录失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    token, ok := resp["data"].(map[string]interface{})["token"].(string)
    require.True(h.T, ok, "响应中缺少token")

    return token
}
```

## AAA测试模式

所有API测试必须遵循AAA模式（Arrange-Act-Assert）：

### Arrange（准备）

```go
// Arrange - 准备测试环境和数据
router, cleanup := integration.SetupTestEnvironment(t)
defer cleanup()

helper := integration.NewTestHelper(t, router)

// 创建测试用户
user := helper.CreateTestUser(&users.User{
    Username: "testuser",
    Email:    "test@example.com",
    Password: "password123",
})

// 登录获取token
token := helper.LoginTestUser(user.Email, "password123")
```

### Act（执行）

```go
// Act - 执行HTTP请求
reqBody := map[string]interface{}{
    "title":       "新书籍",
    "description": "这是一本测试书籍",
    "price":       100,
}
w := helper.DoAuthRequest("POST", "/api/v1/books", reqBody, token)
```

### Assert（断言）

```go
// Assert - 验证HTTP响应
helper.AssertSuccess(w, 201, "创建书籍失败")

// 验证响应结构
helper.AssertResponseStructure(w, []string{"id", "title", "price"})

// 验证返回数据
var resp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &resp)

data := resp["data"].(map[string]interface{})
assert.Equal(t, "新书籍", data["title"])
assert.Equal(t, float64(100), data["price"])
```

## 测试用例设计

### 1. 认证相关测试

#### 登录成功测试

```go
func TestAuthAPI_Login_Success(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建测试用户
    user := helper.CreateTestUser(&users.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    })

    // Act
    reqBody := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    }
    w := helper.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)

    // Assert
    helper.AssertSuccess(w, 200, "登录失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    assert.Equal(t, float64(200), resp["code"])
    assert.NotEmpty(t, resp["data"].(map[string]interface{})["token"])

    // 验证token有效性
    token := resp["data"].(map[string]interface{})["token"].(string)
    claims, err := jwtService.VerifyToken(context.Background(), token)
    require.NoError(t, err)
    assert.Equal(t, user.ID, claims.UserID)
}
```

#### 登录失败测试

```go
func TestAuthAPI_Login_WrongPassword(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    helper.CreateTestUser(&users.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    })

    // Act - 使用错误密码
    reqBody := map[string]interface{}{
        "email":    "test@example.com",
        "password": "wrongpassword",
    }
    w := helper.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)

    // Assert
    helper.AssertError(w, 401, "密码错误")
}

func TestAuthAPI_Login_UserNotFound(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // Act - 用户不存在
    reqBody := map[string]interface{}{
        "email":    "notfound@example.com",
        "password": "password123",
    }
    w := helper.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)

    // Assert
    helper.AssertError(w, 404, "用户不存在")
}
```

### 2. CRUD操作测试

#### 创建资源

```go
func TestBookAPI_CreateBook_Success(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建并登录作者
    author := helper.CreateTestUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    token := helper.LoginTestUser(author.Email, "password123")

    // Act
    reqBody := map[string]interface{}{
        "title":       "测试书籍",
        "description": "这是一本测试书籍",
        "category":    "小说",
        "price":       0, // 免费
    }
    w := helper.DoAuthRequest("POST", "/api/v1/books", reqBody, token)

    // Assert
    helper.AssertSuccess(w, 201, "创建书籍失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    data := resp["data"].(map[string]interface{})
    assert.NotEmpty(t, data["id"])
    assert.Equal(t, "测试书籍", data["title"])
    assert.Equal(t, "author", data["author_name"])

    // 验证数据库中的记录
    db := integration.GetTestDB(t)
    book, err := repository.NewBookRepository(db).GetByID(context.Background(), data["id"].(string))
    require.NoError(t, err)
    assert.Equal(t, "测试书籍", book.Title)
}
```

#### 查询资源

```go
func TestBookAPI_GetBook_Success(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建测试书籍
    db := integration.GetTestDB(t)
    bookRepo := repository.NewBookRepository(db)
    book := &bookstore.Book{
        Title:       "测试书籍",
        Description: "测试描述",
        AuthorID:    "author123",
        Price:       100,
    }
    err := bookRepo.Create(context.Background(), book)
    require.NoError(t, err)

    // Act
    w := helper.DoRequest("GET", "/api/v1/books/"+book.ID, nil, nil)

    // Assert
    helper.AssertSuccess(w, 200, "获取书籍失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    data := resp["data"].(map[string]interface{})
    assert.Equal(t, book.ID, data["id"])
    assert.Equal(t, "测试书籍", data["title"])
}
```

#### 更新资源

```go
func TestBookAPI_UpdateBook_Success(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建作者和书籍
    author := helper.CreateTestUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    token := helper.LoginTestUser(author.Email, "password123")

    db := integration.GetTestDB(t)
    bookRepo := repository.NewBookRepository(db)
    book := &bookstore.Book{
        Title:       "原标题",
        Description: "原描述",
        AuthorID:    author.ID,
        Price:       100,
    }
    err := bookRepo.Create(context.Background(), book)
    require.NoError(t, err)

    // Act - 更新书籍
    reqBody := map[string]interface{}{
        "title":       "新标题",
        "description": "新描述",
        "price":       200,
    }
    w := helper.DoAuthRequest("PUT", "/api/v1/books/"+book.ID, reqBody, token)

    // Assert
    helper.AssertSuccess(w, 200, "更新书籍失败")

    // 验证数据库
    updatedBook, err := bookRepo.GetByID(context.Background(), book.ID)
    require.NoError(t, err)
    assert.Equal(t, "新标题", updatedBook.Title)
    assert.Equal(t, "新描述", updatedBook.Description)
    assert.Equal(t, 200, updatedBook.Price)
}
```

#### 删除资源

```go
func TestBookAPI_DeleteBook_Success(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建作者和书籍
    author := helper.CreateTestUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    token := helper.LoginTestUser(author.Email, "password123")

    db := integration.GetTestDB(t)
    bookRepo := repository.NewBookRepository(db)
    book := &bookstore.Book{
        Title:    "要删除的书籍",
        AuthorID: author.ID,
    }
    err := bookRepo.Create(context.Background(), book)
    require.NoError(t, err)

    // Act - 删除书籍
    w := helper.DoAuthRequest("DELETE", "/api/v1/books/"+book.ID, nil, token)

    // Assert
    helper.AssertSuccess(w, 200, "删除书籍失败")

    // 验证数据库中已删除
    _, err = bookRepo.GetByID(context.Background(), book.ID)
    assert.Error(t, err)
}
```

### 3. 权限控制测试

#### 未授权访问测试

```go
func TestBookAPI_UpdateBook_Unauthorized(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建作者A和书籍
    authorA := helper.CreateTestUser(&users.User{
        Username: "author_a",
        Email:    "author_a@example.com",
        Password: "password123",
    })

    db := integration.GetTestDB(t)
    bookRepo := repository.NewBookRepository(db)
    book := &bookstore.Book{
        Title:    "作者A的书籍",
        AuthorID: authorA.ID,
    }
    bookRepo.Create(context.Background(), book)

    // 创建作者B并登录
    authorB := helper.CreateTestUser(&users.User{
        Username: "author_b",
        Email:    "author_b@example.com",
        Password: "password123",
    })
    tokenB := helper.LoginTestUser(authorB.Email, "password123")

    // Act - 作者B尝试修改作者A的书籍
    reqBody := map[string]interface{}{
        "title": "被篡改的标题",
    }
    w := helper.DoAuthRequest("PUT", "/api/v1/books/"+book.ID, reqBody, tokenB)

    // Assert
    helper.AssertError(w, 403, "没有权限")

    // 验证数据未被修改
    originalBook, _ := bookRepo.GetByID(context.Background(), book.ID)
    assert.Equal(t, "作者A的书籍", originalBook.Title)
}
```

#### 需要认证的接口测试

```go
func TestBookAPI_CreateBook_RequiresAuth(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // Act - 没有token的请求
    reqBody := map[string]interface{}{
        "title": "测试书籍",
    }
    w := helper.DoRequest("POST", "/api/v1/books", reqBody, nil)

    // Assert
    helper.AssertError(w, 401, "需要认证")
}
```

### 4. 输入验证测试

```go
func TestBookAPI_CreateBook_Validation(t *testing.T) {
    tests := []struct {
        name         string
        reqBody      map[string]interface{}
        expectedCode int
        expectedMsg  string
    }{
        {
            name: "标题为空",
            reqBody: map[string]interface{}{
                "title": "",
            },
            expectedCode: 400,
            expectedMsg:  "标题不能为空",
        },
        {
            name: "价格负数",
            reqBody: map[string]interface{}{
                "title": "测试书籍",
                "price": -100,
            },
            expectedCode: 400,
            expectedMsg:  "价格不能为负数",
        },
        {
            name: "分类无效",
            reqBody: map[string]interface{}{
                "title":    "测试书籍",
                "category": "不存在的分类",
            },
            expectedCode: 400,
            expectedMsg:  "无效的分类",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            router, cleanup := integration.SetupTestEnvironment(t)
            defer cleanup()

            helper := integration.NewTestHelper(t, router)

            author := helper.CreateTestUser(&users.User{
                Username: "author",
                Email:    "author@example.com",
                Password: "password123",
            })
            token := helper.LoginTestUser(author.Email, "password123")

            // Act
            w := helper.DoAuthRequest("POST", "/api/v1/books", tt.reqBody, token)

            // Assert
            helper.AssertError(w, tt.expectedCode, tt.expectedMsg)
        })
    }
}
```

### 5. 分页查询测试

```go
func TestBookAPI_ListBooks_Pagination(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建多本测试书籍
    db := integration.GetTestDB(t)
    bookRepo := repository.NewBookRepository(db)
    for i := 1; i <= 25; i++ {
        book := &bookstore.Book{
            Title:    fmt.Sprintf("书籍 %d", i),
            AuthorID: "author123",
            Status:   "published",
        }
        bookRepo.Create(context.Background(), book)
    }

    // Act - 第一页
    w := helper.DoRequest("GET", "/api/v1/books?page=1&page_size=10", nil, nil)

    // Assert
    helper.AssertSuccess(w, 200, "获取书籍列表失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    data := resp["data"].(map[string]interface{})
    items := data["items"].([]interface{})
    assert.Len(t, items, 10, "第一页应该有10条记录")

    total := data["total"].(float64)
    assert.Equal(t, float64(25), total, "总共应该有25条记录")
}
```

### 6. 文件上传测试

```go
func TestBookAPI_UploadCover_Success(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    author := helper.CreateTestUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    token := helper.LoginTestUser(author.Email, "password123")

    // 创建测试图片
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, _ := writer.CreateFormFile("cover", "test_cover.jpg")
    part.Write([]byte("fake image content"))
    writer.Close()

    req := httptest.NewRequest("POST", "/api/v1/books/book123/cover", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", "Bearer "+token)

    // Act
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    helper.AssertSuccess(w, 200, "上传封面失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    data := resp["data"].(map[string]interface{})
    assert.NotEmpty(t, data["cover_url"])
}
```

## Table-Driven测试

```go
func TestAuthAPI_Register_TableDriven(t *testing.T) {
    tests := []struct {
        name         string
        reqBody      map[string]interface{}
        expectedCode int
        expectedMsg  string
        setupFunc    func(*integration.TestHelper) string
    }{
        {
            name: "成功注册",
            reqBody: map[string]interface{}{
                "username": "newuser",
                "email":    "newuser@example.com",
                "password": "password123",
            },
            expectedCode: 201,
            expectedMsg:  "注册成功",
            setupFunc:    nil,
        },
        {
            name: "用户名太短",
            reqBody: map[string]interface{}{
                "username": "ab",
                "email":    "test@example.com",
                "password": "password123",
            },
            expectedCode: 400,
            expectedMsg:  "用户名长度至少3位",
            setupFunc:    nil,
        },
        {
            name: "邮箱格式错误",
            reqBody: map[string]interface{}{
                "username": "testuser",
                "email":    "invalid-email",
                "password": "password123",
            },
            expectedCode: 400,
            expectedMsg:  "邮箱格式无效",
            setupFunc:    nil,
        },
        {
            name: "密码太短",
            reqBody: map[string]interface{}{
                "username": "testuser",
                "email":    "test@example.com",
                "password": "12345",
            },
            expectedCode: 400,
            expectedMsg:  "密码长度至少8位",
            setupFunc:    nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            router, cleanup := integration.SetupTestEnvironment(t)
            defer cleanup()

            helper := integration.NewTestHelper(t, router)

            // Act
            w := helper.DoRequest("POST", "/api/v1/auth/register", tt.reqBody, nil)

            // Assert
            if tt.expectedCode >= 200 && tt.expectedCode < 300 {
                helper.AssertSuccess(w, tt.expectedCode, tt.expectedMsg)
            } else {
                helper.AssertError(w, tt.expectedCode, tt.expectedMsg)
            }
        })
    }
}
```

## 错误处理测试

### 404错误测试

```go
func TestBookAPI_GetBook_NotFound(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // Act - 访问不存在的书籍
    fakeID := primitive.NewObjectID().Hex()
    w := helper.DoRequest("GET", "/api/v1/books/"+fakeID, nil, nil)

    // Assert
    helper.AssertError(w, 404, "书籍不存在")
}
```

### 服务器错误测试

```go
func TestBookAPI_CreateBook_DatabaseError(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 创建用户
    author := helper.CreateTestUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    token := helper.LoginTestUser(author.Email, "password123")

    // Act - 发送超长标题导致数据库错误
    longTitle := strings.Repeat("a", 10000)
    reqBody := map[string]interface{}{
        "title": longTitle,
    }
    w := helper.DoAuthRequest("POST", "/api/v1/books", reqBody, token)

    // Assert
    // 应该返回500而不是panic
    assert.Equal(t, 500, w.Code)
}
```

## 测试环境Setup

### SetupTestEnvironment实现

```go
// test/integration/setup.go
package integration

import (
    "context"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/require"

    "Qingyu_backend/config"
    "Qingyu_backend/container"
    "Qingyu_backend/models"
    "Qingyu_backend/repository/mongodb/shared"
    "Qingyu_backend/router"
)

func SetupTestEnvironment(t *testing.T) (*gin.Engine, func()) {
    t.Helper()

    // 加载测试配置
    cfg, err := config.LoadConfig("config/config_test.yaml")
    require.NoError(t, err, "加载配置失败")

    // 初始化服务容器
    c := container.NewServiceContainer(cfg)
    err = c.Initialize(context.Background())
    require.NoError(t, err, "初始化容器失败")

    // 设置路由
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.Use(gin.Recovery())

    // 注册所有路由
    routerGroup := router.Group("/api/v1")
    router.RegisterRoutes(routerGroup, c)

    // 清理函数
    cleanup := func() {
        // 清理测试数据库
        db := c.GetMongoDB()
        collections := []string{
            "users", "books", "projects", "documents",
            "roles", "permissions", "transactions",
            "user_behaviors", "user_profiles",
            "reading_progress", "annotations",
            "chapters", "announcements",
        }
        ctx := context.Background()
        for _, coll := range collections {
            _ = db.Collection(coll).Drop(ctx)
        }
        _ = c.Close(ctx)
    }

    return router, cleanup
}

func GetTestDB(t *testing.T) *mongo.Database {
    t.Helper()

    cfg, err := config.LoadConfig("config/config_test.yaml")
    require.NoError(t, err)

    c := container.NewServiceContainer(cfg)
    err = c.Initialize(context.Background())
    require.NoError(t, err)

    return c.GetMongoDB()
}
```

## 测试覆盖率目标

### 必须达到的覆盖率

| 接口类型 | 覆盖率目标 | 说明 |
|---------|-----------|------|
| 认证接口 | 100% | 登录、注册、token刷新 |
| CRUD接口 | ≥80% | 包含成功和失败场景 |
| 权限控制 | 100% | 所有权限验证必须测试 |
| 输入验证 | ≥90% | 各种无效输入场景 |
| 错误处理 | ≥80% | 各种错误响应 |

### 覆盖率检查命令

```bash
# 生成API测试覆盖率
go test -coverprofile=coverage.out ./test/api/v1

# 查看覆盖率
go tool cover -func=coverage.out
```

## 并发测试

```go
func TestBookAPI_CreateBook_Concurrent(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    author := helper.CreateTestUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    token := helper.LoginTestUser(author.Email, "password123")

    // Act - 并发创建书籍
    const concurrency = 10
    var wg sync.WaitGroup
    errors := make(chan error, concurrency)

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()

            reqBody := map[string]interface{}{
                "title": fmt.Sprintf("并发书籍 %d", index),
            }
            w := helper.DoAuthRequest("POST", "/api/v1/books", reqBody, token)

            if w.Code != 201 {
                errors <- fmt.Errorf("创建失败, 状态码: %d", w.Code)
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Assert - 验证没有错误
    for err := range errors {
        t.Errorf("并发创建书籍失败: %v", err)
    }
}
```

## 常见问题

### Q1: 为什么API测试必须用集成测试而不是单元测试？

**A**: API层测试需要验证：
- **路由配置**：URL是否正确映射到handler
- **中间件**：认证、日志、CORS等是否正常工作
- **请求解析**：JSON、表单、文件上传等
- **响应格式**：状态码、响应结构、错误信息
- **数据库交互**：完整的请求-响应-数据库流程

只有集成测试能覆盖所有这些方面。

### Q2: 如何处理测试数据库数据隔离？

**A**: 使用`SetupTestEnvironment`的清理机制：
```go
func TestExample(t *testing.T) {
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()  // 测试结束后自动清理所有集合

    // 测试逻辑...
}
```

每个测试都有独立的数据库，测试结束后自动清理，不会互相影响。

### Q3: 如何测试需要异步处理的接口？

**A**: 使用轮询或eventually模式：
```go
func TestAsyncAPI(t *testing.T) {
    // Arrange
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // Act - 触发异步操作
    w := helper.DoRequest("POST", "/api/v1/async-job", nil, token)
    jobID := resp["data"].(map[string]interface{})["job_id"].(string)

    // Assert - 轮询检查结果
    assert.Eventually(t, func() bool {
        w := helper.DoRequest("GET", "/api/v1/async-job/"+jobID, nil, token)
        return w.Code == 200 && resp["data"].(map[string]interface{})["status"] == "completed"
    }, 10*time.Second, 100*time.Millisecond, "异步任务未完成")
}
```

## 参考文档

- [testify使用指南](../03_测试工具指南/testify使用指南.md)
- [集成测试详细规范](../02_测试类型规范/集成测试详细规范.md)
- [Service层测试规范](./service_层测试规范.md)
- [E2E测试规范](./e2e_测试规范.md)
