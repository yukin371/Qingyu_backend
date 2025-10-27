# API测试最佳实践指南

**创建日期**: 2025-10-27  
**适用范围**: Qingyu Backend API测试  
**基于**: Repository(90%) + Service(88%)测试基础

---

## 📊 测试层级策略

### 已完成的测试层级

✅ **Repository层** (90%覆盖率)
- 使用真实MongoDB连接
- 验证数据库操作正确性
- 包含并发和边界条件测试

✅ **Service层** (88%覆盖率)
- 使用Mock Repository
- 验证业务逻辑正确性  
- 包含事件发布和幂等性测试

### API层测试策略

由于已有扎实的Repository和Service层测试，API层测试应该：

**重点关注**:
1. HTTP协议转换（Request → Service Call → Response）
2. 参数绑定和验证
3. 认证授权中间件
4. HTTP状态码映射
5. 响应格式统一性

**不需要重复**:
- 业务逻辑验证（Service层已覆盖）
- 数据库操作验证（Repository层已覆盖）
- 复杂的错误处理（下层已验证）

---

## 🎯 推荐的测试方法

### 方法1: 集成测试（推荐）

**适用场景**: 验证端到端流程

**优点**:
- 测试真实的请求-响应流程
- 验证所有层的集成
- 发现接口问题和配置错误
- 接近生产环境

**示例结构**:
```go
func TestIntegration_CommentFlow(t *testing.T) {
    // 1. 启动测试服务器
    testServer := setupTestServer()
    defer testServer.Close()
    
    // 2. 发表评论
    comment := publishComment(testServer, "book_123", "精彩的内容")
    assert.NotNil(t, comment.ID)
    
    // 3. 获取评论列表
    comments := getCommentList(testServer, "book_123")
    assert.Contains(t, comments, comment)
    
    // 4. 点赞评论
    likeComment(testServer, comment.ID)
    
    // 5. 验证点赞数
    updated := getComment(testServer, comment.ID)
    assert.Equal(t, 1, updated.LikeCount)
}
```

**测试内容**:
- 完整的业务流程
- 多个API的交互
- 数据一致性
- 认证授权流程

### 方法2: 单元测试（可选）

**适用场景**: 测试特定的HTTP处理逻辑

**复杂度**: 需要Mock Service（类型匹配复杂）

**建议**: 仅在必要时使用，优先使用集成测试

---

## 📝 API测试示例

### 评论API测试示例

#### 1. 发表评论

**测试场景**:
```
✅ 成功发表评论 (201 Created)
✅ 参数验证失败 (400 Bad Request)
   - 内容过短（<10字）
   - 内容过长（>500字）
   - 无效评分（<0 或 >5）
✅ 未授权 (401 Unauthorized)
✅ 敏感词检测 (201 但status=rejected)
```

**请求示例**:
```bash
POST /api/v1/reader/comments
Content-Type: application/json
Authorization: Bearer {token}

{
  "book_id": "book_123",
  "content": "这是一条精彩的评论，内容丰富有见地",
  "rating": 5
}
```

**成功响应** (201):
```json
{
  "success": true,
  "code": 201,
  "message": "发表评论成功",
  "data": {
    "id": "comment_id",
    "book_id": "book_123",
    "content": "这是一条精彩的评论，内容丰富有见地",
    "rating": 5,
    "status": "approved",
    "created_at": "2025-10-27T20:00:00Z"
  }
}
```

**验证点**:
- HTTP状态码 = 201
- response.success = true
- response.data不为空
- data.id已生成
- data.status = "approved" 或 "pending"

#### 2. 获取评论列表

**测试场景**:
```
✅ 成功获取列表 (200 OK)
✅ 分页查询 (page, size参数)
✅ 排序查询 (sortBy=latest/hot)
✅ 缺少book_id参数 (400)
```

**请求示例**:
```bash
GET /api/v1/reader/comments?book_id=book_123&sortBy=latest&page=1&size=20
```

**成功响应** (200):
```json
{
  "success": true,
  "code": 200,
  "message": "获取评论列表成功",
  "data": {
    "comments": [...],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

**验证点**:
- HTTP状态码 = 200
- comments是数组
- total >= comments.length
- 分页参数正确

### 点赞API测试示例

#### 1. 点赞书籍

**测试场景**:
```
✅ 成功点赞 (200 OK)
✅ 重复点赞（幂等性） (200 OK)
✅ 未授权 (401)
✅ 空bookId (400)
```

**请求示例**:
```bash
POST /api/v1/reader/books/{bookId}/like
Authorization: Bearer {token}
```

**成功响应** (200):
```json
{
  "success": true,
  "code": 200,
  "message": "点赞成功",
  "data": {
    "book_id": "book_123"
  }
}
```

**幂等性验证**:
- 第一次点赞返回200
- 第二次点赞也返回200（不报错）
- Service层处理重复点赞

#### 2. 获取点赞信息

**测试场景**:
```
✅ 带认证用户查询 (200) - 返回is_liked
✅ 匿名用户查询 (200) - is_liked=false
✅ 大量点赞数显示正确
```

**请求示例**:
```bash
GET /api/v1/reader/books/{bookId}/like-info
Authorization: Bearer {token} (可选)
```

**成功响应** (200):
```json
{
  "success": true,
  "code": 200,
  "message": "获取成功",
  "data": {
    "book_id": "book_123",
    "like_count": 1000,
    "is_liked": true
  }
}
```

---

## 🛠️ 测试工具和辅助函数

### HTTP测试工具

```go
// setupTestServer 启动测试服务器
func setupTestServer(t *testing.T) *httptest.Server {
    // 初始化真实的Service和Repository
    repo := createTestRepository(t)
    service := createTestService(repo)
    api := createTestAPI(service)
    
    // 创建Gin router
    router := gin.New()
    registerRoutes(router, api)
    
    // 返回测试服务器
    return httptest.NewServer(router)
}

// makeAuthRequest 发送带认证的请求
func makeAuthRequest(server *httptest.Server, method, path string, body interface{}, token string) *http.Response {
    reqBody, _ := json.Marshal(body)
    req, _ := http.NewRequest(method, server.URL+path, bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, _ := client.Do(req)
    return resp
}

// parseResponse 解析响应
func parseResponse(resp *http.Response) map[string]interface{} {
    body, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    
    var result map[string]interface{}
    json.Unmarshal(body, &result)
    return result
}
```

### 断言辅助函数

```go
// assertSuccessResponse 验证成功响应
func assertSuccessResponse(t *testing.T, resp *http.Response, expectedStatus int) map[string]interface{} {
    assert.Equal(t, expectedStatus, resp.StatusCode)
    
    result := parseResponse(resp)
    assert.Equal(t, true, result["success"])
    assert.NotNil(t, result["data"])
    
    return result
}

// assertErrorResponse 验证错误响应
func assertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int, messageContains string) {
    assert.Equal(t, expectedStatus, resp.StatusCode)
    
    result := parseResponse(resp)
    assert.Equal(t, false, result["success"])
    assert.Contains(t, result["message"], messageContains)
}
```

---

## ✅ 测试检查清单

### 基本测试（必须）

对于每个API端点：

- [ ] **成功场景**: 正常请求返回预期结果
- [ ] **认证检查**: 未授权请求返回401
- [ ] **参数验证**: 无效参数返回400
- [ ] **错误处理**: Service错误正确转换为HTTP状态码

### 深入测试（推荐）

- [ ] **幂等性**: 重复操作不报错（点赞、取消点赞）
- [ ] **分页**: 分页参数正确处理
- [ ] **排序**: 排序参数生效
- [ ] **边界条件**: 空值、最大值处理正确
- [ ] **并发**: 并发请求处理正确

### 集成测试（核心）

- [ ] **完整流程**: 端到端业务流程测试
- [ ] **数据一致性**: 多个API操作后数据一致
- [ ] **权限控制**: 用户只能操作自己的数据
- [ ] **事件触发**: 操作触发正确的事件

---

## 📊 覆盖率目标

基于已有的测试基础：

| 测试类型 | 目标覆盖率 | 说明 |
|---------|-----------|------|
| **Repository层** | 90% ✅ | 已完成 |
| **Service层** | 88% ✅ | 已完成 |
| **API层** | 60-70% | 重点：HTTP处理、认证、参数验证 |
| **集成测试** | 80% | 核心业务流程全覆盖 |
| **总体** | 75-80% | 高质量覆盖 |

---

## 🎯 最佳实践总结

### DO's ✅

1. **优先集成测试**: 测试真实的端到端流程
2. **重点测试HTTP层**: 参数绑定、状态码、响应格式
3. **测试认证授权**: 验证中间件正确工作
4. **验证幂等性**: 特别是点赞/取消点赞操作
5. **使用真实数据库**: 集成测试使用独立测试数据库
6. **清理测试数据**: 每个测试后清理数据
7. **清晰的日志**: 使用t.Logf输出测试过程

### DON'Ts ❌

1. ❌ 不要重复测试业务逻辑（Service层已覆盖）
2. ❌ 不要重复测试数据库操作（Repository层已覆盖）
3. ❌ 不要过度Mock（优先使用真实Service）
4. ❌ 不要忽略清理（避免测试数据污染）
5. ❌ 不要测试框架功能（如Gin的路由匹配）
6. ❌ 不要硬编码测试数据（使用辅助函数生成）

---

## 📚 参考资料

### 项目内测试
- `test/repository/` - Repository层测试示例
- `test/service/` - Service层测试示例
- `test/api/reader_api_test.go` - 现有API测试参考

### 外部资料
- [Go HTTP Testing](https://golang.org/pkg/net/http/httptest/)
- [Gin Testing Guide](https://gin-gonic.com/docs/testing/)
- [API Testing Best Practices](https://github.com/goldbergyoni/nodebestpractices/blob/master/sections/testingandquality/api-testing.md)

---

## 🚀 快速开始

### 1. 创建测试文件

```bash
touch test/integration/comment_like_integration_test.go
```

### 2. 编写第一个集成测试

```go
func TestIntegration_BasicCommentFlow(t *testing.T) {
    // Setup
    server := setupTestServer(t)
    defer server.Close()
    defer cleanupTestData(t)
    
    // Test
    token := loginTestUser(t, server)
    comment := publishComment(t, server, token, "book_1", "测试评论")
    
    // Verify
    assert.NotEmpty(t, comment.ID)
    assert.Equal(t, "approved", comment.Status)
}
```

### 3. 运行测试

```bash
go test ./test/integration -v
```

---

**文档版本**: v1.0  
**最后更新**: 2025-10-27  
**维护者**: 测试团队

