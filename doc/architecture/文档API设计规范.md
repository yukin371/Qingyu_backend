# 文档API设计规范

## 1. 需求概述

### 1.1 功能描述
文档API是青羽写作系统的核心接口，提供对文档资源的创建、读取、更新和删除等基本操作，支持文档内容的管理和版本控制。该API是前端编辑器与后端存储系统之间的桥梁，确保用户创作内容的安全存储和高效访问。

### 1.2 业务价值
- **内容管理**：提供统一的文档内容管理接口，支持多种文档类型
- **版本控制**：支持文档版本管理，保障内容安全
- **数据一致性**：确保前后端数据交互的一致性和可靠性
- **性能优化**：通过合理的API设计提升系统响应速度和用户体验

### 1.3 用户场景
- 用户创建新文档
- 获取文档列表和单个文档详情
- 更新文档内容和元数据
- 删除不需要的文档

### 1.4 功能边界
- 仅处理文档基本CRUD操作
- 不包含复杂的文档内容解析和渲染
- 不包含权限控制（由认证中间件处理）
- 不包含实时协作功能（由专门的协作API处理）

## 2. API设计

### 2.1 API概览

| 方法   | 路径                | 功能描述         | 权限要求 |
|--------|---------------------|-----------------|----------|
| POST   | /api/v1/document/   | 创建新文档       | 用户登录 |
| GET    | /api/v1/document/   | 获取文档列表     | 用户登录 |
| GET    | /api/v1/document/doc/:id | 获取单个文档 | 用户登录 |
| PUT    | /api/v1/document/doc/:id | 更新文档     | 用户登录 |
| DELETE | /api/v1/document/doc/:id | 删除文档     | 用户登录 |

### 2.2 请求与响应格式

#### 2.2.1 创建文档 (POST /api/v1/document/)

**请求体**:
```json
{
  "projectId": "string",  // 可选，项目ID
  "nodeId": "string",     // 可选，节点ID
  "title": "string",      // 必填，文档标题
  "content": "string",    // 必填，文档内容
  "format": "string"      // 必填，文档格式（如markdown）
}
```

**响应体**:
```json
{
  "id": "string",
  "projectId": "string",
  "nodeId": "string",
  "title": "string",
  "content": "string",
  "format": "string",
  "words": 0,
  "version": 1,
  "createdAt": "2023-01-01T00:00:00Z",
  "updatedAt": "2023-01-01T00:00:00Z"
}
```

**状态码**:
- 201: 创建成功
- 400: 请求参数错误
- 401: 未授权
- 500: 服务器错误

#### 2.2.2 获取文档列表 (GET /api/v1/document/)

**查询参数**:
- userId: 用户ID（可选）
- limit: 每页数量（可选，默认50）
- offset: 偏移量（可选，默认0）

**响应体**:
```json
[
  {
    "id": "string",
    "projectId": "string",
    "nodeId": "string",
    "title": "string",
    "content": "string",
    "format": "string",
    "words": 0,
    "version": 1,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
]
```

**状态码**:
- 200: 成功
- 401: 未授权
- 500: 服务器错误

#### 2.2.3 获取单个文档 (GET /api/v1/document/doc/:id)

**路径参数**:
- id: 文档ID

**响应体**:
```json
{
  "id": "string",
  "projectId": "string",
  "nodeId": "string",
  "title": "string",
  "content": "string",
  "format": "string",
  "words": 0,
  "version": 1,
  "createdAt": "2023-01-01T00:00:00Z",
  "updatedAt": "2023-01-01T00:00:00Z"
}
```

**状态码**:
- 200: 成功
- 401: 未授权
- 404: 文档不存在
- 500: 服务器错误

#### 2.2.4 更新文档 (PUT /api/v1/document/doc/:id)

**路径参数**:
- id: 文档ID

**请求体**:
```json
{
  "title": "string",    // 可选，文档标题
  "content": "string",  // 可选，文档内容
  "format": "string"    // 可选，文档格式
}
```

**响应体**:
```json
{
  "id": "string",
  "projectId": "string",
  "nodeId": "string",
  "title": "string",
  "content": "string",
  "format": "string",
  "words": 0,
  "version": 2,
  "createdAt": "2023-01-01T00:00:00Z",
  "updatedAt": "2023-01-01T00:00:00Z"
}
```

**状态码**:
- 200: 更新成功
- 400: 请求参数错误
- 401: 未授权
- 404: 文档不存在
- 500: 服务器错误

#### 2.2.5 删除文档 (DELETE /api/v1/document/doc/:id)

**路径参数**:
- id: 文档ID

**响应体**:
```json
{
  "success": true
}
```

**状态码**:
- 200: 删除成功
- 401: 未授权
- 404: 文档不存在
- 500: 服务器错误

## 3. 实现细节

### 3.1 Router层设计
```go
// 文档路由组
docRouter := r.Group("/document")
{
    // 文档相关路由
    d := api.NewDocumentApi()
    
    // 基础文档操作
    docRouter.POST("/", d.Create)
    docRouter.GET("/", d.List)
    docRouter.GET("/doc/:id", d.Get)
    docRouter.PUT("/doc/:id", d.Update)
    docRouter.DELETE("/doc/:id", d.Delete)
}
```

### 3.2 API层设计
```go
// DocumentApi 文档相关API
type DocumentApi struct {
    service *svc.DocumentService
}

func NewDocumentApi() *DocumentApi {
    return &DocumentApi{service: &svc.DocumentService{}}
}

// Create 创建文档
func (a *DocumentApi) Create(c *gin.Context) {
    var req model.Document
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }
    created, err := a.service.Create(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, created)
}

// 其他API方法实现...
```

### 3.3 Service层设计
```go
// DocumentService 处理文档相关业务逻辑
type DocumentService struct{}

// Create 创建文档
func (s *DocumentService) Create(doc *model.Document) (*model.Document, error) {
    if doc == nil {
        return nil, errors.New("文档不能为空")
    }
    doc.ID = primitive.NewObjectID().Hex()
    doc.TouchForCreate()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := getCollection().InsertOne(ctx, doc)
    if err != nil {
        return nil, err
    }
    return doc, nil
}

// 其他Service方法实现...
```

## 4. 错误处理

### 4.1 错误码定义

| 错误码 | 描述                 | 处理方式                     |
|--------|---------------------|----------------------------|
| 400    | 请求参数错误         | 检查请求参数格式和必填项     |
| 401    | 未授权              | 用户需要登录或重新认证       |
| 403    | 权限不足            | 用户无权访问该资源           |
| 404    | 资源不存在          | 检查资源ID是否正确           |
| 409    | 资源冲突            | 处理版本冲突或重复操作       |
| 500    | 服务器内部错误      | 联系系统管理员               |

### 4.2 错误响应格式
```json
{
  "error": "错误描述信息",
  "code": "错误代码",
  "details": {
    // 可选的错误详情
  }
}
```

## 5. 安全设计

### 5.1 认证与授权
- 使用JWT进行用户认证
- 基于角色的访问控制（RBAC）
- 文档所有者和协作者权限区分

### 5.2 输入验证
- 请求参数类型和格式验证
- 内容安全过滤（防XSS等）
- 请求大小限制

### 5.3 数据安全
- 敏感数据加密存储
- 传输层安全（TLS）
- 访问日志记录

## 6. 性能考虑

### 6.1 优化策略
- 分页查询优化
- 大文档内容压缩
- 缓存常用文档
- 异步处理长操作

### 6.2 监控指标
- API响应时间
- 错误率
- 并发请求数
- 资源使用率

## 7. 测试策略

### 7.1 单元测试
- API处理函数测试
- 参数验证测试
- 错误处理测试

### 7.2 集成测试
- 端到端API流程测试
- 权限控制测试
- 性能负载测试

## 8. 版本管理

### 8.1 API版本控制
- 使用URL路径版本控制（/api/v1/）
- 向后兼容性保证
- 版本升级策略

### 8.2 变更管理
- API变更文档
- 废弃流程
- 迁移指南