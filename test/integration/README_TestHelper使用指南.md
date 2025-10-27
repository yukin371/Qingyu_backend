# TestHelper 使用指南

**版本**: 1.0  
**创建时间**: 2025-10-25

---

## 一、快速开始

### 1.1 基础用法

```go
package integration

import (
    "testing"
    "Qingyu_backend/test/integration"
)

func TestExample(t *testing.T) {
    // 1. 创建TestHelper
    helper := integration.NewTestHelper(t, router)
    
    // 2. 登录获取token
    token := helper.LoginTestUser()  // 使用默认测试用户
    // 或
    token := helper.LoginUser("custom_user", "password")
    
    // 3. 发送请求
    w := helper.DoAuthRequest("POST", integration.ReaderCollectionsPath, reqBody, token)
    
    // 4. 断言结果
    helper.AssertSuccess(w, 201, "添加收藏失败")
}
```

---

## 二、核心功能

### 2.1 认证功能

#### LoginUser - 登录指定用户
```go
token := helper.LoginUser("username", "password")
```
- 自动处理登录请求
- 返回JWT token
- 登录失败返回空字符串并记录日志

#### LoginTestUser - 登录默认测试用户
```go
token := helper.LoginTestUser()
```
- 使用默认测试用户: `test_user01` / `Test@123456`
- 最常用的登录方式

---

### 2.2 HTTP请求功能

#### DoRequest - 通用请求
```go
w := helper.DoRequest(
    "POST",                           // HTTP方法
    "/api/v1/reader/collections",    // 路径
    map[string]interface{}{          // 请求body（可为nil）
        "book_id": "xxx",
        "note": "我的笔记",
    },
    "",                               // token（可为空）
)
```

#### DoAuthRequest - 认证请求
```go
w := helper.DoAuthRequest("POST", path, body, token)
```
- 自动添加 Authorization header
- Token不能为空，否则会导致测试失败

**最佳实践**:
```go
// ✅ 推荐：使用路径常量
w := helper.DoAuthRequest("POST", integration.ReaderCollectionsPath, body, token)

// ❌ 不推荐：硬编码路径
w := helper.DoAuthRequest("POST", "/api/v1/reader/collections", body, token)
```

---

### 2.3 响应断言功能

#### AssertSuccess - 断言成功响应
```go
response := helper.AssertSuccess(w, 201, "添加收藏失败")
```

**自动验证**:
- ✅ 状态码匹配
- ✅ 响应是有效的JSON
- ✅ 返回解析后的响应数据

**失败时输出**:
```
添加收藏失败
期望状态码: 201
实际状态码: 400
响应内容: {"code":40001,"message":"该书籍已经收藏","data":null}
```

#### AssertError - 断言错误响应
```go
helper.AssertError(w, 400, "已经收藏", "重复收藏检测失败")
```

**验证内容**:
- ✅ 状态码匹配
- ✅ 错误信息包含指定文本
- ✅ 响应格式正确

---

### 2.4 数据库功能

#### GetTestBook - 获取单个测试书籍
```go
bookID := helper.GetTestBook()
if bookID == "" {
    t.Skip("没有测试书籍，跳过测试")
}
```

#### GetTestBooks - 获取多个测试书籍
```go
bookIDs := helper.GetTestBooks(5)  // 获取5本书
if len(bookIDs) < 3 {
    t.Skip("测试书籍不足，跳过测试")
}
```

#### VerifyBookExists - 验证书籍存在
```go
if !helper.VerifyBookExists(bookID) {
    t.Fatalf("书籍 %s 不存在", bookID)
}
```

#### CleanupTestData - 清理测试数据
```go
// 测试结束后清理
defer helper.CleanupTestData("collections", "reading_progress")
```

---

### 2.5 日志功能

#### 日志级别
```go
helper.LogSuccess("操作成功")   // ✓ 成功
helper.LogInfo("信息提示")      // ℹ 信息
helper.LogWarning("警告信息")   // ⚠ 警告
helper.LogError("错误信息")     // ❌ 错误
```

**输出示例**:
```
✓ 登录成功: test_user01 (token: eyJhbGciOiJIUzI1NiIs...)
ℹ 获取5本测试书籍
⚠ 数据库中没有测试书籍
❌ 添加收藏失败
```

---

## 三、路径常量

### 3.1 完整列表

```go
// 基础路径
APIBasePath = "/api/v1"

// 认证
LoginPath    = "/api/v1/login"
RegisterPath = "/api/v1/register"

// 用户
UserProfilePath  = "/api/v1/users/profile"
UserPasswordPath = "/api/v1/users/password"

// 阅读器
ReaderBooksPath       = "/api/v1/reader/books"
ReaderChaptersPath    = "/api/v1/reader/chapters"
ReaderProgressPath    = "/api/v1/reader/progress"
ReaderAnnotationsPath = "/api/v1/reader/annotations"
ReaderCommentsPath    = "/api/v1/reader/comments"
ReaderCollectionsPath = "/api/v1/reader/collections"
ReaderLikesPath       = "/api/v1/reader/likes"

// 书城
BookstoreHomePath    = "/api/v1/bookstore/homepage"
BookstoreBooksPath   = "/api/v1/bookstore/books"
BookstoreRankingPath = "/api/v1/bookstore/rankings"
```

### 3.2 使用示例

```go
// ✅ 推荐
w := helper.DoRequest("GET", integration.BookstoreHomePath, nil, "")

// ❌ 不推荐
w := helper.DoRequest("GET", "/api/v1/bookstore/homepage", nil, "")
```

---

## 四、完整示例

### 4.1 收藏功能测试

```go
func TestCollectionFlow(t *testing.T) {
    // 初始化
    helper := integration.NewTestHelper(t, router)
    token := helper.LoginTestUser()
    bookID := helper.GetTestBook()
    
    t.Run("添加收藏", func(t *testing.T) {
        // 准备请求
        reqBody := map[string]interface{}{
            "book_id": bookID,
            "note":    "测试笔记",
            "tags":    []string{"测试", "收藏"},
        }
        
        // 发送请求
        w := helper.DoAuthRequest("POST", integration.ReaderCollectionsPath, reqBody, token)
        
        // 断言成功
        response := helper.AssertSuccess(w, 201, "添加收藏失败")
        
        // 验证返回数据
        data := response["data"].(map[string]interface{})
        assert.Equal(t, bookID, data["book_id"])
        assert.Equal(t, "测试笔记", data["note"])
        
        helper.LogSuccess("添加收藏成功")
    })
    
    t.Run("重复收藏", func(t *testing.T) {
        // 重复添加相同书籍
        reqBody := map[string]interface{}{
            "book_id": bookID,
            "note":    "重复收藏",
        }
        
        w := helper.DoAuthRequest("POST", integration.ReaderCollectionsPath, reqBody, token)
        
        // 断言失败且错误信息正确
        helper.AssertError(w, 400, "已经收藏", "应该拒绝重复收藏")
    })
    
    t.Run("获取收藏列表", func(t *testing.T) {
        w := helper.DoAuthRequest("GET", integration.ReaderCollectionsPath+"?page=1&size=10", nil, token)
        
        response := helper.AssertSuccess(w, 200, "获取收藏列表失败")
        
        // 验证列表不为空
        data := response["data"].(map[string]interface{})
        list := data["list"].([]interface{})
        assert.Greater(t, len(list), 0, "收藏列表应该有数据")
        
        helper.LogSuccess("获取收藏列表成功，共%d条", len(list))
    })
    
    // 清理测试数据
    defer helper.CleanupTestData("collections")
}
```

### 4.2 认证流程测试

```go
func TestAuthFlow(t *testing.T) {
    helper := integration.NewTestHelper(t, router)
    
    t.Run("成功登录", func(t *testing.T) {
        token := helper.LoginUser("test_user01", "Test@123456")
        assert.NotEmpty(t, token, "Token不应该为空")
        helper.LogSuccess("登录成功")
    })
    
    t.Run("密码错误", func(t *testing.T) {
        reqBody := map[string]interface{}{
            "username": "test_user01",
            "password": "wrong_password",
        }
        
        w := helper.DoRequest("POST", integration.LoginPath, reqBody, "")
        helper.AssertError(w, 401, "密码错误", "应该拒绝错误密码")
    })
    
    t.Run("访问需要认证的接口", func(t *testing.T) {
        // 不带token访问
        w := helper.DoRequest("GET", integration.ReaderBooksPath, nil, "")
        helper.AssertError(w, 401, "未授权", "应该要求登录")
        
        // 带token访问
        token := helper.LoginTestUser()
        w = helper.DoAuthRequest("GET", integration.ReaderBooksPath, nil, token)
        helper.AssertSuccess(w, 200, "认证后应该可以访问")
    })
}
```

### 4.3 数据准备测试

```go
func TestWithTestData(t *testing.T) {
    helper := integration.NewTestHelper(t, router)
    
    // 检查测试数据
    bookIDs := helper.GetTestBooks(3)
    if len(bookIDs) < 3 {
        t.Skip("测试书籍不足（需要至少3本），跳过测试")
    }
    
    token := helper.LoginTestUser()
    
    // 使用测试数据
    for i, bookID := range bookIDs {
        t.Run(fmt.Sprintf("添加书籍%d到书架", i+1), func(t *testing.T) {
            reqBody := map[string]interface{}{
                "book_id": bookID,
            }
            
            w := helper.DoAuthRequest("POST", integration.ReaderBooksPath+"/"+bookID, reqBody, token)
            helper.AssertSuccess(w, 201, "添加书籍%d到书架失败", i+1)
        })
    }
    
    // 清理
    defer helper.CleanupTestData("reading_progress", "bookshelf")
}
```

---

## 五、最佳实践

### 5.1 测试组织

```go
func TestFeature(t *testing.T) {
    // 1. 初始化（只做一次）
    helper := integration.NewTestHelper(t, router)
    token := helper.LoginTestUser()
    
    // 2. 数据准备（如果需要）
    bookID := helper.GetTestBook()
    if bookID == "" {
        t.Skip("没有测试数据，跳过测试")
    }
    
    // 3. 子测试
    t.Run("场景1", func(t *testing.T) {
        // 测试实现
    })
    
    t.Run("场景2", func(t *testing.T) {
        // 测试实现
    })
    
    // 4. 清理（defer确保执行）
    defer helper.CleanupTestData("test_collection")
}
```

### 5.2 错误处理

```go
// ✅ 推荐：详细的错误信息
helper.AssertSuccess(w, 201, 
    "添加收藏失败 - 书籍ID: %s, 用户: %s", 
    bookID, username)

// ❌ 不推荐：模糊的错误信息
helper.AssertSuccess(w, 201, "失败")
```

### 5.3 测试独立性

```go
// ✅ 推荐：每个测试独立准备数据
t.Run("测试1", func(t *testing.T) {
    data := prepareTestData()
    defer cleanupTestData(data)
    // 测试逻辑
})

// ❌ 不推荐：测试间共享可变状态
sharedData := prepareTestData()
t.Run("测试1", func(t *testing.T) {
    modifyData(sharedData)  // 影响后续测试
})
```

### 5.4 断言精确性

```go
// ✅ 推荐：精确断言
response := helper.AssertSuccess(w, 201, "添加失败")
data := response["data"].(map[string]interface{})
assert.Equal(t, expectedID, data["id"])
assert.Equal(t, expectedName, data["name"])

// ❌ 不推荐：只检查状态码
helper.AssertSuccess(w, 201, "添加失败")
// 没有验证返回的数据
```

---

## 六、常见问题

### Q1: TestHelper 和原来的写法比有什么优势？

**A**: 主要优势：
1. **代码量减少70%** - 无需重复编写请求构造代码
2. **错误信息详细** - 自动显示完整的请求/响应上下文
3. **统一性** - 所有测试使用相同的方式
4. **可维护性** - API变更时只需修改常量定义

### Q2: 如何调试测试失败？

**A**: 使用详细的日志：
```go
// 1. 使用LogInfo记录关键步骤
helper.LogInfo("准备添加收藏: bookID=%s", bookID)

// 2. AssertSuccess会自动输出详细错误
helper.AssertSuccess(w, 201, "添加收藏失败")

// 3. 手动打印响应（如果需要）
t.Logf("完整响应: %s", w.Body.String())
```

### Q3: 如何处理需要特殊header的请求？

**A**: 使用DoRequest手动构造：
```go
var bodyReader io.Reader
if body != nil {
    bodyBytes, _ := json.Marshal(body)
    bodyReader = bytes.NewReader(bodyBytes)
}

req := httptest.NewRequest("POST", path, bodyReader)
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer "+token)
req.Header.Set("X-Custom-Header", "value")

w := httptest.NewRecorder()
helper.router.ServeHTTP(w, req)
```

### Q4: 测试数据如何准备？

**A**: 三种方式：
```go
// 1. 使用现有数据库数据
bookID := helper.GetTestBook()

// 2. 在测试中创建临时数据
bookID := createTestBook(t)
defer deleteTestBook(t, bookID)

// 3. 使用固定的测试数据
bookID := "507f1f77bcf86cd799439011"  // 已知存在的测试数据
```

### Q5: 如何测试并发场景？

**A**: 示例：
```go
func TestConcurrent(t *testing.T) {
    helper := integration.NewTestHelper(t, router)
    token := helper.LoginTestUser()
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            
            reqBody := map[string]interface{}{
                "book_id": fmt.Sprintf("book_%d", index),
            }
            
            w := helper.DoAuthRequest("POST", integration.ReaderCollectionsPath, reqBody, token)
            // 注意：并发时的断言可能需要额外处理
            if w.Code != 201 {
                t.Logf("并发请求%d失败: %d", index, w.Code)
            }
        }(i)
    }
    wg.Wait()
}
```

---

## 七、迁移指南

### 7.1 迁移步骤

**步骤1: 引入TestHelper**
```go
// 在测试函数开头添加
helper := integration.NewTestHelper(t, router)
```

**步骤2: 替换登录代码**
```go
// 替换前（~20行）
loginData := map[string]interface{}{...}
body, _ := json.Marshal(loginData)
req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
// ... 更多代码
token := ...

// 替换后（1行）
token := helper.LoginTestUser()
```

**步骤3: 替换请求代码**
```go
// 替换前（~10行）
reqBody := map[string]interface{}{...}
body, _ := json.Marshal(reqBody)
req := httptest.NewRequest("POST", path, bytes.NewReader(body))
req.Header.Set("Authorization", "Bearer "+token)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

// 替换后（1行）
w := helper.DoAuthRequest("POST", integration.ReaderCollectionsPath, reqBody, token)
```

**步骤4: 替换断言代码**
```go
// 替换前
assert.Equal(t, 201, w.Code)
var response map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &response)

// 替换后
response := helper.AssertSuccess(w, 201, "添加收藏失败")
```

### 7.2 迁移检查清单

- [ ] 使用 `NewTestHelper` 创建helper
- [ ] 使用 `LoginTestUser()` 或 `LoginUser()` 登录
- [ ] 使用路径常量而不是硬编码路径
- [ ] 使用 `DoAuthRequest()` 发送认证请求
- [ ] 使用 `AssertSuccess()` 或 `AssertError()` 断言
- [ ] 使用 `LogSuccess()` 等记录关键步骤
- [ ] 使用 `CleanupTestData()` 清理数据

---

## 八、更新日志

| 版本 | 日期 | 变更内容 |
|-----|------|---------|
| 1.0 | 2025-10-25 | 初始版本发布 |

---

**维护者**: 后端开发团队  
**问题反馈**: 请在项目仓库提Issue

