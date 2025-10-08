# 项目节点文档CRUD设计

## 1. 需求概述

### 1.1 功能描述
项目节点文档CRUD功能是青羽写作系统的核心功能，负责管理项目中的文档节点，包括章节、段落、角色卡等各种类型的文档内容。该功能提供完整的创建、读取、更新、删除操作，支持层次化的文档结构管理。

### 1.2 业务价值
- **结构化管理**：提供层次化的文档组织方式，便于作者管理复杂的创作内容
- **版本控制**：支持文档版本管理，保障内容安全
- **协作支持**：为未来的多人协作功能提供基础
- **AI集成**：为AI辅助创作提供数据基础

### 1.3 用户场景
- 作者创建新的章节或段落
- 编辑现有文档内容
- 查看文档历史版本
- 删除不需要的文档节点
- 调整文档结构和层次关系

### 1.4 功能边界
- 支持文本类型的文档节点
- 支持层次化的父子关系
- 支持软删除和恢复
- 不包含实时协作编辑功能（后续版本支持）

## 2. 架构设计

### 2.1 整体架构
```
┌─────────────────┐
│   Frontend      │ ← 用户界面
├─────────────────┤
│   Router        │ ← 路由层：/api/v1/document/*
├─────────────────┤
│   API Layer     │ ← HTTP接口处理
├─────────────────┤
│   Service Layer │ ← 业务逻辑处理
├─────────────────┤
│   Model Layer   │ ← 数据模型定义
├─────────────────┤
│   MongoDB       │ ← 数据持久化
└─────────────────┘
```

### 2.2 模块划分
- **DocumentAPI**：HTTP接口处理模块
- **DocumentService**：文档业务逻辑模块
- **DocumentModel**：文档数据模型
- **VersionService**：版本控制服务
- **PermissionService**：权限控制服务

### 2.3 数据流向
1. 客户端发送HTTP请求
2. Router层路由到对应API处理函数
3. API层验证参数并调用Service层
4. Service层执行业务逻辑并操作Model层
5. Model层与MongoDB交互
6. 返回结果沿原路径返回

### 2.4 技术选型
- **Web框架**：Gin
- **数据库**：MongoDB
- **ORM**：MongoDB官方驱动
- **认证**：JWT中间件
- **日志**：Zap

## 3. 详细设计

### 3.1 Router层设计
```go
// router/document/document.go
func InitDocumentRouter(router *gin.RouterGroup) {
    documentApi := api.NewDocumentApi()
    
    docGroup := router.Group("/document")
    docGroup.Use(middleware.AuthMiddleware()) // 需要认证
    {
        docGroup.POST("", documentApi.CreateDocument)
        docGroup.GET("/:id", documentApi.GetDocument)
        docGroup.PUT("/:id", documentApi.UpdateDocument)
        docGroup.DELETE("/:id", documentApi.DeleteDocument)
        docGroup.GET("/project/:projectId", documentApi.ListDocuments)
        docGroup.POST("/:id/restore", documentApi.RestoreDocument)
    }
}
```

### 3.2 API层设计
```go
// api/v1/document/document.go
type DocumentApi struct {
    documentService *service.DocumentService
}

// CreateDocument 创建文档节点
func (api *DocumentApi) CreateDocument(c *gin.Context) {
    var req CreateDocumentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.FailWithMessage("参数验证失败", c)
        return
    }
    
    userID := c.GetString("userID")
    document, err := api.documentService.CreateDocument(userID, &req)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    
    response.OkWithData(document, c)
}
```

### 3.3 Service层设计
```go
// service/document/document.go
type DocumentService struct {
    // 依赖注入
}

func (s *DocumentService) CreateDocument(userID string, req *CreateDocumentRequest) (*models.Document, error) {
    // 1. 权限验证
    if err := s.validatePermission(userID, req.ProjectID); err != nil {
        return nil, err
    }
    
    // 2. 业务逻辑处理
    document := &models.Document{
        ProjectID: req.ProjectID,
        ParentID:  req.ParentID,
        Title:     req.Title,
        Content:   req.Content,
        Type:      req.Type,
        CreatedBy: userID,
    }
    
    // 3. 数据持久化
    if err := s.saveDocument(document); err != nil {
        return nil, err
    }
    
    return document, nil
}
```

### 3.4 Model层设计
```go
// models/document/document.go
type Document struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    ProjectID string    `bson:"project_id" json:"projectId"`
    ParentID  string    `bson:"parent_id,omitempty" json:"parentId,omitempty"`
    Title     string    `bson:"title" json:"title"`
    Content   string    `bson:"content" json:"content"`
    Type      DocType   `bson:"type" json:"type"`
    Order     int       `bson:"order" json:"order"`
    CreatedBy string    `bson:"created_by" json:"createdBy"`
    CreatedAt time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
    DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

type DocType string

const (
    DocTypeChapter   DocType = "chapter"
    DocTypeParagraph DocType = "paragraph"
    DocTypeCharacter DocType = "character"
    DocTypeOutline   DocType = "outline"
)
```

## 4. 数据设计

### 4.1 数据模型
参见Model层设计中的Document结构体定义。

### 4.2 数据关系
- **项目关系**：Document.ProjectID → Project.ID
- **层次关系**：Document.ParentID → Document.ID（自引用）
- **用户关系**：Document.CreatedBy → User.ID

### 4.3 索引策略
```javascript
// MongoDB索引设计
db.documents.createIndex({"project_id": 1, "deleted_at": 1})
db.documents.createIndex({"parent_id": 1, "order": 1})
db.documents.createIndex({"created_by": 1})
db.documents.createIndex({"created_at": -1})
```

### 4.4 数据迁移
初始版本无需数据迁移，后续版本升级时需要考虑：
- 新增字段的默认值设置
- 索引的创建和删除
- 数据格式的转换

## 5. 接口设计

### 5.1 API规范
遵循RESTful设计原则，使用标准HTTP方法：
- POST：创建资源
- GET：获取资源
- PUT：更新资源
- DELETE：删除资源

### 5.2 请求示例

#### 创建文档
```http
POST /api/v1/document
Content-Type: application/json
Authorization: Bearer <token>

{
    "projectId": "project123",
    "parentId": "parent456",
    "title": "第一章",
    "content": "章节内容...",
    "type": "chapter"
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "doc789",
        "projectId": "project123",
        "parentId": "parent456",
        "title": "第一章",
        "content": "章节内容...",
        "type": "chapter",
        "order": 1,
        "createdBy": "user123",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
    }
}
```

### 5.3 错误处理
```json
{
    "code": 400,
    "message": "参数验证失败",
    "data": null
}
```

### 5.4 性能要求
- 文档创建响应时间 < 200ms
- 文档查询响应时间 < 100ms
- 支持并发用户数 > 1000

## 6. 安全设计

### 6.1 权限控制
- **项目权限**：只有项目所有者可以操作项目下的文档
- **文档权限**：基于项目权限继承
- **操作权限**：创建、编辑、删除需要相应权限

### 6.2 数据安全
- **输入验证**：严格验证所有输入参数
- **SQL注入防护**：使用参数化查询
- **XSS防护**：对输出内容进行转义
- **CSRF防护**：使用CSRF Token

### 6.3 输入验证
```go
type CreateDocumentRequest struct {
    ProjectID string  `json:"projectId" binding:"required"`
    ParentID  string  `json:"parentId"`
    Title     string  `json:"title" binding:"required,max=200"`
    Content   string  `json:"content" binding:"max=1000000"`
    Type      DocType `json:"type" binding:"required,oneof=chapter paragraph character outline"`
}
```

### 6.4 审计日志
记录所有文档操作的审计日志：
- 操作用户
- 操作时间
- 操作类型
- 操作对象
- 操作结果

## 7. 测试设计

### 7.1 测试策略
- **单元测试**：测试Service层业务逻辑
- **集成测试**：测试API接口功能
- **性能测试**：测试并发性能
- **安全测试**：测试权限控制和输入验证

### 7.2 测试用例

#### 功能测试用例
1. **创建文档**
   - 正常创建文档
   - 无权限创建文档
   - 参数验证失败

2. **查询文档**
   - 查询存在的文档
   - 查询不存在的文档
   - 查询已删除的文档

3. **更新文档**
   - 正常更新文档
   - 更新不存在的文档
   - 无权限更新文档

4. **删除文档**
   - 软删除文档
   - 删除不存在的文档
   - 无权限删除文档

### 7.3 性能测试
- **负载测试**：模拟1000并发用户
- **压力测试**：测试系统极限
- **稳定性测试**：长时间运行测试

### 7.4 安全测试
- **权限测试**：验证权限控制机制
- **注入测试**：测试SQL注入防护
- **XSS测试**：测试跨站脚本防护

## 8. 部署和运维

### 8.1 部署方案
- **容器化部署**：使用Docker容器
- **负载均衡**：使用Nginx进行负载均衡
- **数据库集群**：MongoDB副本集部署
- **缓存层**：Redis缓存热点数据

### 8.2 监控指标
- **响应时间**：API响应时间监控
- **错误率**：接口错误率监控
- **并发数**：当前并发用户数
- **数据库性能**：MongoDB性能指标

### 8.3 日志策略
```go
// 结构化日志记录
logger.Info("创建文档",
    zap.String("userId", userID),
    zap.String("projectId", req.ProjectID),
    zap.String("documentId", document.ID),
    zap.Duration("duration", time.Since(start)),
)
```

### 8.4 故障处理
- **数据库连接失败**：自动重连机制
- **服务异常**：健康检查和自动重启
- **数据备份**：定期数据备份和恢复测试

## 9. 风险评估

### 9.1 技术风险
- **数据库性能**：大量文档可能影响查询性能
  - **应对措施**：合理设计索引，使用分页查询
- **并发冲突**：多用户同时编辑可能产生冲突
  - **应对措施**：使用乐观锁机制

### 9.2 业务风险
- **数据丢失**：误删除或系统故障导致数据丢失
  - **应对措施**：软删除机制，定期备份
- **权限泄露**：权限控制不当导致数据泄露
  - **应对措施**：严格的权限验证机制

### 9.3 性能风险
- **响应延迟**：大文档或高并发导致响应延迟
  - **应对措施**：内容分页，缓存优化
- **存储空间**：大量文档内容占用存储空间
  - **应对措施**：内容压缩，归档机制

### 9.4 安全风险
- **注入攻击**：恶意输入导致安全问题
  - **应对措施**：严格输入验证，参数化查询
- **权限绕过**：权限控制漏洞
  - **应对措施**：多层权限验证，安全审计

## 10. 实施计划

### 10.1 开发阶段
**第一阶段（2周）**：基础CRUD功能
- 实现Document模型
- 实现基础API接口
- 完成单元测试

**第二阶段（1周）**：权限和安全
- 实现权限控制
- 添加输入验证
- 完成安全测试

**第三阶段（1周）**：优化和完善
- 性能优化
- 错误处理完善
- 文档和部署

### 10.2 测试阶段
- **单元测试**：开发过程中同步进行
- **集成测试**：第二阶段完成后进行
- **性能测试**：第三阶段进行
- **用户验收测试**：所有功能完成后进行

### 10.3 上线计划
- **灰度发布**：先在测试环境验证
- **分批上线**：逐步开放给用户
- **监控观察**：密切监控系统状态
- **快速回滚**：准备回滚方案

### 10.4 后续优化
- **性能优化**：根据使用情况优化性能
- **功能扩展**：添加协作编辑功能
- **用户体验**：根据用户反馈改进界面
- **AI集成**：集成AI辅助创作功能