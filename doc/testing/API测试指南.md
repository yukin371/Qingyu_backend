# API测试指南

**版本**: 1.0  
**更新日期**: 2025-10-17  
**维护者**: 青羽后端团队

---

## 📖 概述

本文档提供青羽后端API测试的完整指南，包括单元测试、集成测试和端到端测试的最佳实践。

---

## 🎯 API测试层次

```
┌─────────────────────────────────────┐
│   E2E Tests (端到端测试)             │  ← 完整用户流程
├─────────────────────────────────────┤
│   Integration Tests (集成测试)       │  ← API + Service + DB
├─────────────────────────────────────┤
│   Unit Tests (单元测试)              │  ← 纯Handler逻辑
└─────────────────────────────────────┘
```

---

## 1. API单元测试

### 1.1 测试Handler逻辑

**位置**: `api/v1/xxx/xxx_test.go`

**示例**：测试用户注册API
```go
// api/v1/system/sys_user_test.go
package system

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "Qingyu_backend/service"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockUserService 模拟用户服务
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) Register(req *service.RegisterRequest) (*service.RegisterResponse, error) {
    args := m.Called(req)
    return args.Get(0).(*service.RegisterResponse), args.Error(1)
}

func TestUserApi_Register(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    
    mockService := new(MockUserService)
    api := NewUserApi(mockService)
    
    // 准备请求
    reqBody := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    }
    jsonData, _ := json.Marshal(reqBody)
    
    // 设置Mock期望
    mockService.On("Register", mock.MatchedBy(func(req *service.RegisterRequest) bool {
        return req.Username == "testuser" && req.Email == "test@example.com"
    })).Return(&service.RegisterResponse{
        UserID:   "user-123",
        Username: "testuser",
        Email:    "test@example.com",
    }, nil)
    
    // Act
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
    c.Request.Header.Set("Content-Type", "application/json")
    
    api.Register(c)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    assert.Equal(t, float64(201), response["code"])
    assert.Equal(t, "用户注册成功", response["message"])
    assert.NotNil(t, response["data"])
    
    mockService.AssertExpectations(t)
}
```

### 1.2 测试参数验证

```go
func TestUserApi_Register_InvalidEmail(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := new(MockUserService)
    api := NewUserApi(mockService)
    
    // 无效的邮箱格式
    reqBody := map[string]interface{}{
        "username": "testuser",
        "email":    "invalid-email",
        "password": "password123",
    }
    jsonData, _ := json.Marshal(reqBody)
    
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
    c.Request.Header.Set("Content-Type", "application/json")
    
    api.Register(c)
    
    // 验证返回400错误
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Contains(t, response["error"], "email")
}
```

---

## 2. API集成测试

### 2.1 完整流程测试

**位置**: `test/api/xxx_api_test.go`

**示例**：测试书店API
```go
// test/api/bookstore_api_test.go
package api_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "Qingyu_backend/api/v1/bookstore"
    "Qingyu_backend/service"
    "Qingyu_backend/repository/mongodb"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func setupBookstoreAPI(t *testing.T) (*gin.Engine, func()) {
    // 设置测试环境
    gin.SetMode(gin.TestMode)
    r := gin.New()
    
    // 初始化真实的Repository和Service
    factory := mongodb.NewMongoRepositoryFactory(testConfig)
    bookRepo := factory.CreateBookRepository()
    bookService := service.NewBookService(bookRepo)
    bookApi := bookstore.NewBookApi(bookService)
    
    // 注册路由
    api := r.Group("/api/v1")
    api.GET("/books", bookApi.ListBooks)
    api.GET("/books/:id", bookApi.GetBook)
    api.POST("/books", bookApi.CreateBook)
    
    cleanup := func() {
        // 清理测试数据
        factory.Close()
    }
    
    return r, cleanup
}

func TestBookstoreAPI_CreateAndGetBook(t *testing.T) {
    router, cleanup := setupBookstoreAPI(t)
    defer cleanup()
    
    // 1. 创建书籍
    createReq := map[string]interface{}{
        "title":       "测试书籍",
        "author":      "测试作者",
        "description": "测试描述",
        "price":       99.99,
    }
    jsonData, _ := json.Marshal(createReq)
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var createResp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &createResp)
    bookID := createResp["data"].(map[string]interface{})["id"].(string)
    
    // 2. 获取创建的书籍
    w = httptest.NewRecorder()
    req = httptest.NewRequest("GET", "/api/v1/books/"+bookID, nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var getResp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &getResp)
    book := getResp["data"].(map[string]interface{})
    
    assert.Equal(t, "测试书籍", book["title"])
    assert.Equal(t, "测试作者", book["author"])
}
```

### 2.2 认证和权限测试

```go
func TestBookstoreAPI_RequiresAuthentication(t *testing.T) {
    router, cleanup := setupBookstoreAPI(t)
    defer cleanup()
    
    // 没有token的请求
    w := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewBuffer([]byte("{}")))
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBookstoreAPI_RequiresAdminPermission(t *testing.T) {
    router, cleanup := setupBookstoreAPI(t)
    defer cleanup()
    
    // 普通用户token
    token := generateUserToken("user-123", "user")
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest("DELETE", "/api/v1/books/book-123", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusForbidden, w.Code)
}
```

---

## 3. 端到端测试

### 3.1 完整用户流程

**位置**: `test/integration/xxx_e2e_test.go`

**示例**：用户注册到购买的完整流程
```go
// test/integration/e2e_user_lifecycle_test.go
package integration_test

import (
    "testing"
    "net/http"
    "bytes"
    "encoding/json"
    
    "github.com/stretchr/testify/assert"
)

func TestE2E_UserPurchaseFlow(t *testing.T) {
    // 1. 用户注册
    registerResp := apiRequest(t, "POST", "/api/v1/register", map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    })
    
    userID := extractField(registerResp, "data.id")
    assert.NotEmpty(t, userID)
    
    // 2. 用户登录
    loginResp := apiRequest(t, "POST", "/api/v1/login", map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    })
    
    token := extractField(loginResp, "data.token")
    assert.NotEmpty(t, token)
    
    // 3. 浏览书籍列表
    booksResp := apiRequestWithAuth(t, "GET", "/api/v1/books", nil, token)
    books := extractField(booksResp, "data.items")
    assert.NotEmpty(t, books)
    
    // 4. 查看书籍详情
    bookID := books.([]interface{})[0].(map[string]interface{})["id"].(string)
    bookResp := apiRequestWithAuth(t, "GET", "/api/v1/books/"+bookID, nil, token)
    assert.Equal(t, 200, int(bookResp["code"].(float64)))
    
    // 5. 添加到购物车
    cartResp := apiRequestWithAuth(t, "POST", "/api/v1/cart/add", map[string]interface{}{
        "bookId":   bookID,
        "quantity": 1,
    }, token)
    assert.Equal(t, 200, int(cartResp["code"].(float64)))
    
    // 6. 创建订单
    orderResp := apiRequestWithAuth(t, "POST", "/api/v1/orders", map[string]interface{}{
        "items": []map[string]interface{}{
            {"bookId": bookID, "quantity": 1},
        },
    }, token)
    
    orderID := extractField(orderResp, "data.orderId")
    assert.NotEmpty(t, orderID)
    
    // 7. 支付订单
    paymentResp := apiRequestWithAuth(t, "POST", "/api/v1/orders/"+orderID.(string)+"/pay", map[string]interface{}{
        "paymentMethod": "alipay",
    }, token)
    
    assert.Equal(t, 200, int(paymentResp["code"].(float64)))
}
```

---

## 4. 测试工具函数

### 4.1 HTTP辅助函数

```go
// test/testutil/api_helpers.go
package testutil

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
)

// APIRequest 发送API请求
func APIRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
    var reqBody []byte
    if body != nil {
        reqBody, _ = json.Marshal(body)
    }
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    return w
}

// APIRequestWithAuth 发送带认证的API请求
func APIRequestWithAuth(t *testing.T, router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
    var reqBody []byte
    if body != nil {
        reqBody, _ = json.Marshal(body)
    }
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    
    router.ServeHTTP(w, req)
    return w
}

// ParseResponse 解析响应
func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    return response
}

// AssertAPISuccess 断言API成功
func AssertAPISuccess(t *testing.T, w *httptest.ResponseRecorder) {
    assert.Equal(t, http.StatusOK, w.Code)
    
    resp := ParseResponse(w)
    code, ok := resp["code"].(float64)
    assert.True(t, ok)
    assert.True(t, code >= 200 && code < 300)
}

// AssertAPIError 断言API错误
func AssertAPIError(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
    assert.Equal(t, expectedCode, w.Code)
    
    resp := ParseResponse(w)
    assert.NotEmpty(t, resp["error"])
}
```

---

## 5. 表驱动测试

### 5.1 参数验证测试

```go
func TestUserApi_Register_ValidationErrors(t *testing.T) {
    tests := []struct {
        name       string
        request    map[string]interface{}
        expectCode int
        expectErr  string
    }{
        {
            name: "缺少用户名",
            request: map[string]interface{}{
                "email":    "test@example.com",
                "password": "password123",
            },
            expectCode: http.StatusBadRequest,
            expectErr:  "username",
        },
        {
            name: "邮箱格式错误",
            request: map[string]interface{}{
                "username": "testuser",
                "email":    "invalid-email",
                "password": "password123",
            },
            expectCode: http.StatusBadRequest,
            expectErr:  "email",
        },
        {
            name: "密码太短",
            request: map[string]interface{}{
                "username": "testuser",
                "email":    "test@example.com",
                "password": "123",
            },
            expectCode: http.StatusBadRequest,
            expectErr:  "password",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockService := new(MockUserService)
            api := NewUserApi(mockService)
            
            jsonData, _ := json.Marshal(tt.request)
            
            // Act
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
            c.Request.Header.Set("Content-Type", "application/json")
            
            api.Register(c)
            
            // Assert
            assert.Equal(t, tt.expectCode, w.Code)
            assert.Contains(t, w.Body.String(), tt.expectErr)
        })
    }
}
```

---

## 6. Postman集成

### 6.1 导出Postman测试

参考：[Postman测试指南](./Postman测试指南.md)

### 6.2 Newman自动化

```bash
# 安装Newman
npm install -g newman

# 运行Postman集合
newman run postman_collection.json -e postman_environment.json

# 生成HTML报告
newman run postman_collection.json -r html --reporter-html-export report.html
```

---

## 7. 最佳实践

### 7.1 测试隔离

✅ **DO**:
```go
func TestAPI_Feature(t *testing.T) {
    // 每个测试独立创建环境
    router, cleanup := setupTestAPI(t)
    defer cleanup()
    
    // 测试逻辑
}
```

❌ **DON'T**:
```go
var globalRouter *gin.Engine  // 不要使用全局变量

func TestAPI_Feature(t *testing.T) {
    // 依赖全局状态
}
```

### 7.2 测试数据管理

✅ **DO**:
```go
func TestAPI_CreateBook(t *testing.T) {
    // 使用工厂创建测试数据
    book := fixtures.NewBookFactory().Create()
    
    // 测试后清理
    defer cleanup(book.ID)
}
```

### 7.3 响应验证

✅ **DO**:
```go
// 验证完整响应结构
assert.Equal(t, http.StatusOK, w.Code)

var response Response
json.Unmarshal(w.Body.Bytes(), &response)

assert.Equal(t, 200, response.Code)
assert.Equal(t, "success", response.Message)
assert.NotNil(t, response.Data)
```

---

## 8. 错误场景测试

### 8.1 网络错误模拟

```go
func TestAPI_HandleServiceUnavailable(t *testing.T) {
    // Mock服务返回错误
    mockService.On("GetBook", mock.Anything).Return(nil, errors.New("service unavailable"))
    
    w := apiRequest(t, "GET", "/api/v1/books/123")
    
    assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
```

### 8.2 超时测试

```go
func TestAPI_RequestTimeout(t *testing.T) {
    // 设置短超时
    router := setupRouterWithTimeout(1 * time.Millisecond)
    
    // 触发慢查询
    w := apiRequest(t, router, "GET", "/api/v1/slow-endpoint")
    
    assert.Equal(t, http.StatusRequestTimeout, w.Code)
}
```

---

## 9. 性能测试

### 9.1 基准测试

```go
func BenchmarkAPI_ListBooks(b *testing.B) {
    router, cleanup := setupTestAPI(&testing.T{})
    defer cleanup()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        req := httptest.NewRequest("GET", "/api/v1/books", nil)
        router.ServeHTTP(w, req)
    }
}
```

---

## 10. CI/CD集成

### 10.1 GitHub Actions示例

```yaml
name: API Tests

on: [push, pull_request]

jobs:
  api-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Start MongoDB
        uses: supercharge/mongodb-github-action@1.8.0
      
      - name: Run API Tests
        run: |
          go test -v ./test/api/...
          go test -v ./test/integration/...
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v2
```

---

## 📚 相关资源

- [测试组织规范](./测试组织规范.md)
- [Postman测试指南](./Postman测试指南.md)
- [性能测试规范](./性能测试规范.md)
- [共享服务测试文档](./共享服务测试文档.md)

---

**最后更新**: 2025-10-17  
**维护者**: 青羽后端团队

