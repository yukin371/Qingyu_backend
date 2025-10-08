# 项目CRUD设计文档

## 1. 概述

本文档详细描述了青羽后端系统中项目管理模块的CRUD（创建、读取、更新、删除）功能设计。项目管理是整个文档管理系统的核心模块，负责管理用户的小说工程项目，包括项目的基本信息、状态管理、权限控制等功能。

### 1.1 设计目标

- **数据完整性**：确保项目数据的一致性和完整性
- **权限安全**：实现基于所有者的权限控制机制
- **性能优化**：通过合理的索引设计提升查询性能
- **扩展性**：支持未来功能扩展，如协作、分享等
- **软删除**：支持项目的软删除和恢复机制

### 1.2 技术架构

项目CRUD功能采用分层架构设计：

```
API层 (HTTP接口) → Service层 (业务逻辑) → Model层 (数据模型) → Database层 (MongoDB)
```

## 2. 数据模型设计

### 2.1 Project 模型

项目模型是整个项目管理系统的核心数据结构：

```go
// Project 表示一本小说工程
type Project struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    OwnerID     string    `bson:"owner_id" json:"ownerId"`
    Name        string    `bson:"name" json:"name"`
    Status      string    `bson:"status" json:"status"` // public | private
    Description string    `bson:"description,omitempty" json:"description,omitempty"`
    CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}
```

#### 字段说明

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| ID | string | 是 | 项目唯一标识符，使用ObjectID |
| OwnerID | string | 是 | 项目所有者用户ID |
| Name | string | 是 | 项目名称 |
| Status | string | 是 | 项目状态：public(公开) / private(私有) |
| Description | string | 否 | 项目描述信息 |
| CreatedAt | time.Time | 是 | 创建时间 |
| UpdatedAt | time.Time | 是 | 最后更新时间 |

### 2.2 项目状态枚举

```go
// ProjectStatus 工程状态
type ProjectStatus string

const (
    ProjectStatusPublic  ProjectStatus = "public"
    ProjectStatusPrivate ProjectStatus = "private"
)

// IsValidProjectStatus 验证工程状态是否合法
func IsValidProjectStatus(status string) bool {
    switch ProjectStatus(status) {
    case ProjectStatusPublic, ProjectStatusPrivate:
        return true
    default:
        return false
    }
}
```

## 3. 数据库设计

### 3.1 集合结构

- **集合名称**：`projects`
- **存储引擎**：MongoDB WiredTiger
- **分片策略**：按 `owner_id` 进行分片（可选）

### 3.2 索引设计

为了优化查询性能，设计了以下索引：

```go
// EnsureIndexes 创建项目相关的 MongoDB 索引（幂等）
func (s *ProjectService) EnsureIndexes(ctx context.Context) error {
    indexModels := []mongo.IndexModel{
        {
            // 项目ID唯一索引
            Keys:    bson.M{"_id": 1},
            Options: options.Index().SetUnique(true),
        },
        {
            // 所有者ID索引，用于按所有者查询项目
            Keys: bson.M{"owner_id": 1},
        },
        {
            // 状态索引，用于按状态过滤项目
            Keys: bson.M{"status": 1},
        },
        {
            // 复合索引：所有者ID和状态，用于同时按所有者和状态过滤
            Keys: bson.M{"owner_id": 1, "status": 1},
        },
        {
            // 创建时间索引，用于排序
            Keys: bson.M{"created_at": -1},
        },
        {
            // 删除时间索引，用于软删除查询
            Keys:    bson.M{"deleted_at": 1},
            Options: options.Index().SetSparse(true), // 稀疏索引
        },
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexModels)
    return err
}
```

#### 索引说明

| 索引名称 | 字段 | 类型 | 用途 |
|----------|------|------|------|
| 主键索引 | _id | 唯一 | 主键查询 |
| 所有者索引 | owner_id | 普通 | 按用户查询项目 |
| 状态索引 | status | 普通 | 按状态过滤 |
| 复合索引 | owner_id + status | 复合 | 用户项目状态查询 |
| 时间索引 | created_at | 降序 | 时间排序 |
| 软删除索引 | deleted_at | 稀疏 | 软删除查询 |

## 4. Service层设计

### 4.1 ProjectService 结构

```go
type ProjectService struct{}

func projectCol() *mongo.Collection { 
    return global.DB.Collection("projects") 
}
```

### 4.2 CRUD 方法详解

#### 4.2.1 Create - 创建项目

```go
func (s *ProjectService) CreateProject(ctx context.Context, p *model.Project) error
```

**功能描述**：创建新的项目

**参数说明**：
- `ctx`：上下文对象，用于超时控制和取消操作
- `p`：项目对象指针

**业务逻辑**：
1. 验证必填字段（Name、ID不能为空）
2. 设置创建时间和更新时间
3. 生成新的ObjectID作为项目ID
4. 插入数据库

**实现代码**：
```go
func (s *ProjectService) CreateProject(ctx context.Context, p *model.Project) error {
    if p.Name == "" || p.ID == "" {
        return errors.New("invalid arguments")
    }
    p.CreatedAt = time.Now()
    p.UpdatedAt = p.CreatedAt
    p.ID = primitive.NewObjectID().Hex()
    if _, err := projectCol().InsertOne(ctx, p); err != nil {
        return err
    }
    return nil
}
```

**错误处理**：
- 参数验证失败：返回 "invalid arguments" 错误
- 数据库插入失败：返回具体的数据库错误

#### 4.2.2 Read - 读取项目

##### 单个项目查询

```go
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*model.Project, error)
```

**功能描述**：根据项目ID获取项目详情

**参数说明**：
- `ctx`：上下文对象
- `projectID`：项目ID

**业务逻辑**：
1. 验证项目ID不为空
2. 根据ID查询项目
3. 返回项目对象

**实现代码**：
```go
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*model.Project, error) {
    if projectID == "" {
        return nil, errors.New("未提供项目id")
    }
    var p model.Project
    if err := projectCol().FindOne(ctx, bson.M{"_id": projectID}).Decode(&p); err != nil {
        return nil, err
    }
    return &p, nil
}
```

**错误处理**：
- 项目ID为空：返回 "未提供项目id" 错误
- 项目不存在：返回 MongoDB 的 ErrNoDocuments 错误

##### 项目列表查询

```go
func (s *ProjectService) List(ctx context.Context, ownerID string, status string) ([]*model.Project, error)
```

**功能描述**：获取项目列表，支持按所有者和状态过滤

**参数说明**：
- `ctx`：上下文对象
- `ownerID`：所有者ID（可选）
- `status`：项目状态（可选）

**业务逻辑**：
1. 构建查询过滤条件
2. 按创建时间降序排序
3. 执行查询并返回结果列表

**实现代码**：
```go
func (s *ProjectService) List(ctx context.Context, ownerID string, status string) ([]*model.Project, error) {
    filter := bson.M{}
    if ownerID != "" {
        filter["owner_id"] = ownerID
    }
    if status != "" {
        filter["status"] = status
    }
    opt := options.Find().SetSort(bson.M{"created_at": -1})
    cur, err := projectCol().Find(ctx, filter, opt)
    if err != nil {
        return nil, err
    }
    list := make([]*model.Project, 0)
    if err := cur.All(ctx, &list); err != nil {
        return nil, err
    }
    return list, nil
}
```

**查询优化**：
- 使用复合索引 `owner_id + status` 优化查询性能
- 支持分页查询（可扩展）

#### 4.2.3 Update - 更新项目

```go
func (s *ProjectService) Update(ctx context.Context, projectID, ownerID string, upd model.Project) error
```

**功能描述**：更新项目信息（仅允许修改 name、description、status）

**参数说明**：
- `ctx`：上下文对象
- `projectID`：项目ID
- `ownerID`：所有者ID（权限验证）
- `upd`：更新的项目信息

**业务逻辑**：
1. 构建更新字段集合
2. 强制验证所有者权限（只能修改自己的项目）
3. 更新 updated_at 时间戳
4. 执行更新操作

**实现代码**：
```go
func (s *ProjectService) Update(ctx context.Context, projectID, ownerID string, upd model.Project) error {
    set := bson.M{"updated_at": time.Now()}
    if upd.Name != "" {
        set["name"] = upd.Name
    }
    if upd.Description != "" {
        set["description"] = upd.Description
    }
    if upd.Status != "" {
        set["status"] = upd.Status
    }
    res, err := projectCol().UpdateOne(ctx,
        bson.M{"_id": projectID, "owner_id": ownerID}, // 强制只能改自己的
        bson.M{"$set": set})
    if res.MatchedCount == 0 {
        return errors.New("project not found")
    }
    if mongo.IsDuplicateKeyError(err) {
        return errors.New("project name duplicate")
    }
    return err
}
```

**权限控制**：
- 通过 `owner_id` 字段确保用户只能修改自己的项目
- 如果匹配的文档数为0，返回 "project not found" 错误

**错误处理**：
- 项目名重复：检测 MongoDB 重复键错误，返回 "project name duplicate"
- 权限不足：返回 "project not found" 错误

#### 4.2.4 Delete - 删除项目

##### 软删除

```go
func (s *ProjectService) Delete(ctx context.Context, projectID, ownerID string) error
```

**功能描述**：软删除项目（标记删除，不实际删除数据）

**参数说明**：
- `ctx`：上下文对象
- `projectID`：项目ID
- `ownerID`：所有者ID（权限验证）

**业务逻辑**：
1. 验证所有者权限
2. 设置 `deleted_at` 时间戳
3. 将状态更改为 "deleted"
4. 保留原始数据以支持恢复

**实现代码**：
```go
func (s *ProjectService) Delete(ctx context.Context, projectID, ownerID string) error {
    res, err := projectCol().UpdateOne(ctx,
        bson.M{"_id": projectID, "owner_id": ownerID},
        bson.M{"$set": bson.M{"deleted_at": time.Now(), "status": "deleted"}})
    if res.MatchedCount == 0 {
        return errors.New("project not found")
    }
    return err
}
```

**软删除优势**：
- 数据安全：避免误删除造成的数据丢失
- 审计追踪：保留删除记录
- 恢复功能：支持后续恢复操作

##### 硬删除

```go
func (s *ProjectService) DeleteHard(ctx context.Context, projectID string) error
```

**功能描述**：硬删除项目（物理删除，管理后台使用）

**实现代码**：
```go
func (s *ProjectService) DeleteHard(ctx context.Context, projectID string) error {
    _, err := projectCol().DeleteOne(ctx, bson.M{"_id": projectID})
    return err
}
```

**使用场景**：
- 管理员清理无效数据
- 系统维护和数据清理
- 合规要求的数据销毁

**注意事项**：
- 此操作不可逆
- 建议仅在管理后台使用
- 需要额外的权限验证

### 4.3 辅助方法

#### 权限验证

```go
func (s *ProjectService) IsOwner(ctx context.Context, projectID, userID string) bool
```

**功能描述**：判断用户是否为项目所有者

**实现代码**：
```go
func (s *ProjectService) IsOwner(ctx context.Context, projectID, userID string) bool {
    return projectCol().FindOne(ctx, bson.M{"_id": projectID, "owner_id": userID}).Err() == nil
}
```

**用途**：
- 权限中间件验证
- API接口权限检查
- 业务逻辑权限控制

#### 事务操作

```go
func (s *ProjectService) CreateWithRootNode(ctx context.Context, p *model.Project, rootNode *model.Node) error
```

**功能描述**：事务创建项目并初始化根节点

**实现代码**：
```go
func (s *ProjectService) CreateWithRootNode(ctx context.Context, p *model.Project, rootNode *model.Node) error {
    return global.MongoClient.UseSession(ctx, func(sc mongo.SessionContext) error {
        if err := sc.StartTransaction(); err != nil {
            return err
        }
        defer sc.EndSession(ctx)

        if err := s.CreateProject(ctx, p); err != nil {
            sc.AbortTransaction(sc)
            return err
        }
        rootNode.ProjectID = p.ID
        rootNode.TouchForCreate()
        if _, err := global.DB.Collection("nodes").InsertOne(sc, rootNode); err != nil {
            sc.AbortTransaction(sc)
            return err
        }
        return sc.CommitTransaction(sc)
    })
}
```

**事务流程**：
1. 开启 MongoDB 事务
2. 创建项目记录
3. 创建根节点记录
4. 提交事务或回滚

**事务优势**：
- 数据一致性：确保项目和根节点同时创建成功
- 原子性：要么全部成功，要么全部失败
- 隔离性：避免并发操作的数据冲突

## 5. API层设计

### 5.1 API接口规范

基于 RESTful 设计原则，项目API接口遵循以下规范：

- **基础路径**：`/api/v1/projects`
- **HTTP方法**：GET、POST、PUT、DELETE
- **响应格式**：JSON
- **错误处理**：统一错误码和错误信息

### 5.2 接口定义

#### 5.2.1 创建项目

```http
POST /api/v1/projects
Content-Type: application/json

{
    "name": "我的小说项目",
    "status": "private",
    "description": "这是一个科幻小说项目"
}
```

**响应示例**：
```json
{
    "id": "507f1f77bcf86cd799439011",
    "ownerId": "507f1f77bcf86cd799439012",
    "name": "我的小说项目",
    "status": "private",
    "description": "这是一个科幻小说项目",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### 5.2.2 获取项目列表

```http
GET /api/v1/projects?ownerId=507f1f77bcf86cd799439012&status=private
```

**查询参数**：
- `ownerId`：所有者ID（可选）
- `status`：项目状态（可选）
- `page`：页码（可选，默认1）
- `limit`：每页数量（可选，默认20）

**响应示例**：
```json
{
    "projects": [
        {
            "id": "507f1f77bcf86cd799439011",
            "ownerId": "507f1f77bcf86cd799439012",
            "name": "我的小说项目",
            "status": "private",
            "description": "这是一个科幻小说项目",
            "createdAt": "2024-01-15T10:30:00Z",
            "updatedAt": "2024-01-15T10:30:00Z"
        }
    ],
    "total": 1,
    "page": 1,
    "limit": 20
}
```

#### 5.2.3 获取单个项目

```http
GET /api/v1/projects/507f1f77bcf86cd799439011
```

**响应示例**：
```json
{
    "id": "507f1f77bcf86cd799439011",
    "ownerId": "507f1f77bcf86cd799439012",
    "name": "我的小说项目",
    "status": "private",
    "description": "这是一个科幻小说项目",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### 5.2.4 更新项目

```http
PUT /api/v1/projects/507f1f77bcf86cd799439011
Content-Type: application/json

{
    "name": "更新后的项目名称",
    "status": "public",
    "description": "更新后的项目描述"
}
```

#### 5.2.5 删除项目

```http
DELETE /api/v1/projects/507f1f77bcf86cd799439011
```

**响应示例**：
```json
{
    "deleted": true,
    "message": "项目已成功删除"
}
```

### 5.3 错误处理

#### 错误响应格式

```json
{
    "error": "错误信息",
    "code": "错误码",
    "timestamp": "2024-01-15T10:30:00Z"
}
```

#### 常见错误码

| HTTP状态码 | 错误码 | 错误信息 | 说明 |
|------------|--------|----------|------|
| 400 | INVALID_REQUEST | 请求参数无效 | 请求体格式错误或必填字段缺失 |
| 401 | UNAUTHORIZED | 未授权访问 | 用户未登录或token无效 |
| 403 | FORBIDDEN | 权限不足 | 用户无权限访问该资源 |
| 404 | NOT_FOUND | 项目不存在 | 指定的项目ID不存在 |
| 409 | CONFLICT | 项目名称重复 | 项目名称已存在 |
| 500 | INTERNAL_ERROR | 服务器内部错误 | 数据库连接失败等系统错误 |

## 6. 总结

本文档详细描述了青羽后端系统项目CRUD功能的完整设计方案，涵盖了从数据模型到API接口、从业务逻辑到性能优化的各个方面。

### 6.1 核心特性

- **完整的CRUD操作**：支持项目的创建、读取、更新、删除
- **权限控制**：基于所有者的权限管理机制
- **软删除支持**：安全的删除和恢复机制
- **性能优化**：合理的索引设计和查询优化
- **扩展性设计**：支持未来功能扩展

### 6.2 技术亮点

- **分层架构**：清晰的代码组织结构
- **事务支持**：确保数据一致性
- **错误处理**：完善的错误处理机制
- **索引优化**：针对查询场景的索引设计
- **权限验证**：安全的权限控制机制

### 6.3 函数功能说明

#### ProjectService 核心方法

1. **EnsureIndexes(ctx context.Context) error**
   - 功能：创建项目相关的 MongoDB 索引
   - 特点：幂等操作，可重复执行
   - 用途：系统初始化时创建必要的索引

2. **CreateProject(ctx context.Context, p *model.Project) error**
   - 功能：创建新项目
   - 验证：检查必填字段
   - 处理：自动生成ID和时间戳

3. **GetProjectByID(ctx context.Context, projectID string) (*model.Project, error)**
   - 功能：根据ID获取单个项目
   - 验证：检查项目ID有效性
   - 返回：项目详细信息

4. **List(ctx context.Context, ownerID string, status string) ([]*model.Project, error)**
   - 功能：获取项目列表
   - 过滤：支持按所有者和状态过滤
   - 排序：按创建时间降序排列

5. **Update(ctx context.Context, projectID, ownerID string, upd model.Project) error**
   - 功能：更新项目信息
   - 权限：仅允许所有者更新
   - 字段：支持更新名称、描述、状态

6. **Delete(ctx context.Context, projectID, ownerID string) error**
   - 功能：软删除项目
   - 安全：标记删除而非物理删除
   - 权限：仅允许所有者删除

7. **DeleteHard(ctx context.Context, projectID string) error**
   - 功能：硬删除项目
   - 用途：管理后台数据清理
   - 风险：不可逆操作

8. **IsOwner(ctx context.Context, projectID, userID string) bool**
   - 功能：验证用户是否为项目所有者
   - 用途：权限验证和中间件
   - 返回：布尔值表示权限状态

9. **CreateWithRootNode(ctx context.Context, p *model.Project, rootNode *model.Node) error**
   - 功能：事务创建项目和根节点
   - 特点：保证数据一致性
   - 机制：MongoDB事务支持

通过本设计文档的实施，可以构建一个功能完善、性能优良、易于维护的项目管理系统，为青羽平台的用户提供优质的项目管理体验。

## 关联文档
- 软件需求规格说明书(SRS) ../软件需求规格说明书(SRS).md
- 架构设计说明书 ../架构设计说明书.md
- API 接口总览 ../API接口总览.md
- 数据库设计说明书 ../数据库设计说明书.md
- 测试计划与用例 ../测试计划与用例.md
- 部署与运维指南 ../部署与运维指南.md
- 安全设计与威胁建模 ../安全设计与威胁建模.md
- 日志与监控 ../日志与监控.md
- 需求追踪矩阵 ../需求追踪矩阵.md