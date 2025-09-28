# API设计规范

## 1. 概述

本文档定义了青羽项目API设计的标准和规范，确保API的一致性、可维护性和易用性。

## 2. RESTful API设计原则

### 2.1 URL设计规范

#### 2.1.1 基本规则
- 使用名词而非动词
- 使用复数形式表示资源集合
- 使用小写字母，单词间用连字符分隔
- 避免深层嵌套，最多3层

#### 2.1.2 URL结构
```
/api/v1/{resource}/{id}/{sub-resource}/{sub-id}
```

#### 2.1.3 示例
```
GET    /api/v1/projects           # 获取项目列表
GET    /api/v1/projects/123       # 获取特定项目
POST   /api/v1/projects           # 创建新项目
PUT    /api/v1/projects/123       # 更新项目
DELETE /api/v1/projects/123       # 删除项目

GET    /api/v1/projects/123/documents    # 获取项目下的文档
POST   /api/v1/projects/123/documents    # 在项目下创建文档
```

### 2.2 HTTP方法使用规范

| 方法 | 用途 | 幂等性 | 安全性 |
|------|------|--------|--------|
| GET | 获取资源 | ✅ | ✅ |
| POST | 创建资源 | ❌ | ❌ |
| PUT | 更新/替换资源 | ✅ | ❌ |
| PATCH | 部分更新资源 | ❌ | ❌ |
| DELETE | 删除资源 | ✅ | ❌ |

### 2.3 状态码使用规范

#### 2.3.1 成功响应
- `200 OK` - 请求成功
- `201 Created` - 资源创建成功
- `204 No Content` - 请求成功但无返回内容

#### 2.3.2 客户端错误
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 未认证
- `403 Forbidden` - 无权限
- `404 Not Found` - 资源不存在
- `409 Conflict` - 资源冲突
- `422 Unprocessable Entity` - 请求格式正确但语义错误

#### 2.3.3 服务器错误
- `500 Internal Server Error` - 服务器内部错误
- `502 Bad Gateway` - 网关错误
- `503 Service Unavailable` - 服务不可用

## 3. 请求和响应格式

### 3.1 请求格式

#### 3.1.1 Content-Type
- JSON请求：`application/json`
- 文件上传：`multipart/form-data`
- 表单提交：`application/x-www-form-urlencoded`

#### 3.1.2 请求头规范
```
Content-Type: application/json
Authorization: Bearer {token}
X-Request-ID: {unique-id}
```

### 3.2 响应格式

#### 3.2.1 统一响应结构
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "uuid"
}
```

#### 3.2.2 分页响应结构
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "uuid"
}
```

#### 3.2.3 错误响应结构
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
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "uuid"
}
```

## 4. 参数规范

### 4.1 查询参数

#### 4.1.1 分页参数
- `page` - 页码，从1开始
- `page_size` - 每页数量，默认20，最大100
- `sort` - 排序字段
- `order` - 排序方向：asc/desc

#### 4.1.2 过滤参数
- `filter[field]` - 字段过滤
- `search` - 全文搜索
- `created_at_start` - 创建时间开始
- `created_at_end` - 创建时间结束

#### 4.1.3 示例
```
GET /api/v1/projects?page=1&page_size=20&sort=created_at&order=desc&filter[status]=active&search=小说
```

### 4.2 路径参数
- 使用有意义的标识符
- 支持UUID和数字ID
- 参数名使用下划线命名

### 4.3 请求体参数
- 使用驼峰命名法（camelCase）
- 必填参数明确标注
- 提供参数类型和格式说明

## 5. 认证和授权

### 5.1 认证方式
- 使用JWT Bearer Token
- Token放在Authorization头中
- 支持Token刷新机制

### 5.2 权限控制
- 基于角色的访问控制（RBAC）
- 资源级权限验证
- 操作级权限验证

## 6. 版本控制

### 6.1 版本策略
- URL路径版本控制：`/api/v1/`
- 语义化版本号：major.minor.patch
- 向后兼容性保证

### 6.2 版本升级
- 新版本保持向后兼容
- 废弃功能提前通知
- 提供迁移指南

## 7. 错误处理

### 7.1 错误分类
- 客户端错误（4xx）
- 服务器错误（5xx）
- 业务逻辑错误

### 7.2 错误信息
- 提供清晰的错误描述
- 包含错误代码和类型
- 提供解决建议

## 8. 性能优化

### 8.1 缓存策略
- 使用适当的缓存头
- 支持条件请求
- 实现ETag机制

### 8.2 限流控制
- 基于用户的限流
- 基于IP的限流
- 提供限流信息头

## 9. 文档规范

### 9.1 API文档要求
- 使用OpenAPI 3.0规范
- 提供完整的示例
- 包含错误场景说明

### 9.2 文档维护
- 代码变更同步更新文档
- 定期审查文档准确性
- 提供交互式文档

## 10. 测试规范

### 10.1 测试覆盖
- 单元测试覆盖率 > 80%
- 集成测试覆盖主要流程
- 端到端测试覆盖关键场景

### 10.2 测试数据
- 使用模拟数据
- 避免依赖外部服务
- 保证测试环境一致性

## 11. 监控和日志

### 11.1 日志规范
- 统一日志格式
- 记录关键操作
- 包含请求追踪ID

### 11.2 监控指标
- 响应时间
- 错误率
- 吞吐量
- 可用性

## 12. 安全规范

### 12.1 输入验证
- 严格验证所有输入
- 防止SQL注入
- 防止XSS攻击

### 12.2 数据保护
- 敏感数据加密
- 不在日志中记录敏感信息
- 实现数据脱敏

## 13. 附录

### 13.1 常用HTTP状态码
详见第2.3节状态码使用规范

### 13.2 API设计检查清单
- [ ] URL设计符合RESTful规范
- [ ] 使用正确的HTTP方法
- [ ] 返回适当的状态码
- [ ] 响应格式统一
- [ ] 错误处理完善
- [ ] 文档完整准确
- [ ] 安全措施到位
- [ ] 性能优化合理