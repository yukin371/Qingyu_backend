# API设计规范

**版本**: v2.0
**更新**: 2026-01-08
**状态**: ✅ 正式实施

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

**基于角色的访问控制（RBAC）**：
- 普通用户
- VIP用户
- 作者
- 管理员

**资源级权限**：
- 用户只能操作自己的资源
- 作者只能编辑自己的作品

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

---

**相关文档**：
- [架构设计规范](../architecture/架构设计规范.md)
- [路由层设计规范](../architecture/路由层设计规范.md)
