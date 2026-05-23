# 项目四层架构CRUD设计

> 最后整理: 2026-05-22  
> 当前状态: `legacy-live`
> 历史架构版本: v2.1  
> 创建日期: 2025-10-21  
> 历史维护者: 青羽架构组

## Page Role

- 这里负责：项目/节点/文档/内容四层结构的历史 CRUD 设计稿，包括模型拆分、协作流程和 API 草图。
- 不负责：当前项目模块真实分层、当前路由与 API 事实、现行版本管理实现结论。

## Recommended Read Path

1. [README.md](./README.md)
2. [项目路径管理设计.md](./项目路径管理设计.md)
3. [../../../architecture/README.md](../../../architecture/README.md)
4. [../../../api/README.md](../../../api/README.md)
5. [../../../implementation/README.md](../../../implementation/README.md)

## Boundary

- 如果你要找“当前项目架构边界”，优先看 [../../../architecture/README.md](../../../architecture/README.md)。
- 如果你要找“当前接口入口”，优先看 [../../../api/README.md](../../../api/README.md)。
- 如果你要找“已落地实施记录”，优先看 [../../../implementation/README.md](../../../implementation/README.md)。
- 本页更适合解释四层 CRUD 当时如何规划，不适合作为当前实现事实的唯一依据。

## Quick Section Map

| 如果你想看 | 直接跳到 |
|------|------|
| 四层结构与整体架构 | [概述](#概述) / [架构设计](#架构设计) |
| 四层数据模型 | [四层数据模型详解](#四层数据模型详解) |
| 典型协作流程 | [四层协作流程](#四层协作流程) |
| 数据库与 API 设计 | [数据库设计](#数据库设计) / [API设计](#api设计) |
| 方案优势与关联资料 | [架构优势总结](#架构优势总结) / [相关文档](#相关文档) |

## Quick Takeaways

- 这篇最有价值的地方，是它把项目结构拆成“Project / Node / Document / DocumentContent”四层，并配了典型流程。
- 如果你只关心项目文档管理主线，优先看“架构设计”“四层数据模型详解”“四层协作流程”“API设计”。
- 文中的四层划分和接口示意属于历史设计口径，不应直接视为当前代码结构事实。

## Skip Guide

- 只想知道“为什么要拆四层”：看 [概述](#概述) 和 [架构设计](#架构设计)。
- 只想知道“项目创建、建章节、编辑内容怎么走”：看 [四层协作流程](#四层协作流程)。
- 如果你当前只关心现行架构或接口事实，请优先回到 `docs/architecture/`、`docs/api/`、`docs/implementation/`。

## 📋 概述

本文档详细描述青羽写作平台的**四层架构CRUD设计**，该架构通过分层设计实现高性能、可扩展的文档管理系统。

### 四层架构

```
Project (项目层)
  ↓
Node (节点层 - 树形结构)
  ↓
Document (文档元数据层)
  ↓
DocumentContent (文档内容层 - 支持GridFS大文本)
```

---

## 🏗️ 架构设计

### 1. 整体架构图

```
┌─────────────────────────────────────────────┐
│  Frontend (用户界面)                         │
├─────────────────────────────────────────────┤
│  Router (路由层)                             │
│  /api/v1/projects/*                         │
│  /api/v1/nodes/*                            │
│  /api/v1/documents/*                        │
│  /api/v1/documents/:id/content              │
├─────────────────────────────────────────────┤
│  API Layer (接口层)                         │
│  - ProjectApi                               │
│  - NodeApi                                  │
│  - DocumentApi                              │
│  - DocumentContentApi                       │
├─────────────────────────────────────────────┤
│  Service Layer (业务逻辑层)                 │
│  - ProjectService                           │
│  - NodeService (树形结构管理)               │
│  - DocumentService (元数据管理)             │
│  - DocumentContentService (内容管理+GridFS) │
├─────────────────────────────────────────────┤
│  Repository Layer (数据访问层)              │
│  - ProjectRepository                        │
│  - NodeRepository                           │
│  - DocumentRepository                       │
│  - DocumentContentRepository                │
│  - GridFSRepository (大文件存储)            │
├─────────────────────────────────────────────┤
│  Database Layer (数据库层)                  │
│  - MongoDB Collections:                     │
│    • projects (轻量)                        │
│    • nodes (树形结构)                       │
│    • documents (元数据 - 轻量)              │
│    • document_contents (实际内容 - 重量)    │
│    • fs.files & fs.chunks (GridFS)         │
└─────────────────────────────────────────────┘
```

---

## 📊 四层数据模型详解

### 第一层：Project（项目层）

**职责**：项目级别的管理和元数据

```go
type Project struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    OwnerID     string    `bson:"owner_id" json:"ownerId"`
    Name        string    `bson:"name" json:"name"`
    Description string    `bson:"description,omitempty" json:"description"`
    Status      string    `bson:"status" json:"status"` // draft | active | archived
    Type        string    `bson:"type" json:"type"`     // novel | essay | script
    RootNodeID  string    `bson:"root_node_id" json:"rootNodeId"` // 根节点ID
    
    // 统计信息（冗余字段，从下层同步）
    TotalWords      int       `bson:"total_words" json:"totalWords"`
    TotalDocuments  int       `bson:"total_documents" json:"totalDocuments"`
    
    CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
    DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}
```

**特点**：
- ✅ 轻量级，只包含项目级别元数据
- ✅ 软删除支持
- ✅ 统计信息冗余（性能优化）

**使用场景**：
- 项目列表查询
- 项目创建/删除
- 项目统计信息

---

### 第二层：Node（节点层 - 树形结构）

**职责**：管理文档的树形结构关系

```go
type Node struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    ProjectID   string    `bson:"project_id" json:"projectId"`
    ParentID    string    `bson:"parent_id,omitempty" json:"parentId"` // null表示根节点
    
    // 树形结构字段
    Level       int       `bson:"level" json:"level"`           // 层级：0(根), 1(卷), 2(章), 3(节)
    Order       int       `bson:"order" json:"order"`           // 同级排序
    Path        string    `bson:"path" json:"path"`             // 路径：/1/2/5 (快速查询祖先)
    
    // 节点类型
    Type        NodeType  `bson:"type" json:"type"`             // root | volume | chapter | section
    
    // 关联文档（可选，叶子节点才有）
    DocumentID  string    `bson:"document_id,omitempty" json:"documentId,omitempty"`
    
    // 元数据
    Title       string    `bson:"title" json:"title"`
    Icon        string    `bson:"icon,omitempty" json:"icon"`
    
    // 统计（冗余）
    ChildCount  int       `bson:"child_count" json:"childCount"`
    
    CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
    DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

type NodeType string

const (
    NodeTypeRoot    NodeType = "root"    // 根节点
    NodeTypeVolume  NodeType = "volume"  // 卷
    NodeTypeChapter NodeType = "chapter" // 章
    NodeTypeSection NodeType = "section" // 节
)
```

**设计要点**：

**1. Path字段加速祖先查询**
```go
// 查询某节点的所有祖先
ancestors := strings.Split(node.Path, "/")
// /1/2/5 → [1, 2, 5]
```

**2. Level字段控制层级**
```go
// 最多3层（不含root）
maxLevel := 3
if node.Level >= maxLevel {
    return errors.New("超过最大层级限制")
}
```

**3. Order字段支持拖拽排序**
```go
// 同级节点按order排序
nodes, _ := repo.GetChildNodes(parentID) // ORDER BY order ASC
```

**特点**：
- ✅ 树形结构管理
- ✅ 快速查询子节点
- ✅ 支持拖拽排序
- ✅ Path字段加速祖先查询

**使用场景**：
- 文档树展示
- 节点拖拽排序
- 章节目录生成

---

### 第三层：Document（文档元数据层）

**职责**：管理文档的元数据和关联关系（不包含实际内容）

```go
type Document struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    ProjectID   string    `bson:"project_id" json:"projectId"`
    NodeID      string    `bson:"node_id" json:"nodeId"` // 关联到Node（一对一）
    
    // 元数据
    Title       string    `bson:"title" json:"title"`
    Status      string    `bson:"status" json:"status"`   // draft | published
    Tags        []string  `bson:"tags,omitempty" json:"tags"`
    Notes       string    `bson:"notes,omitempty" json:"notes"`
    
    // 关联设定（角色、地点、时间线等）
    CharacterIDs []string  `bson:"character_ids,omitempty" json:"characterIds"`
    LocationIDs  []string  `bson:"location_ids,omitempty" json:"locationIds"`
    TimelineIDs  []string  `bson:"timeline_ids,omitempty" json:"timelineIds"`
    
    // 统计信息（从DocumentContent同步）
    WordCount    int       `bson:"word_count" json:"wordCount"`
    CharCount    int       `bson:"char_count" json:"charCount"`
    
    // 版本控制
    Version      int       `bson:"version" json:"version"`     // 乐观锁版本号
    
    CreatedBy    string    `bson:"created_by" json:"createdBy"`
    CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
    DeletedAt    *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}
```

**关键设计**：

**❌ 不包含Content字段**
```go
// 错误示例
type Document struct {
    Content string `bson:"content"` // ❌ 不应该在Document中
}

// ✅ 正确设计：内容在DocumentContent中
```

**✅ WordCount冗余字段**
```go
// 查询文档列表时，可以直接显示字数
documents, _ := repo.GetDocuments(projectID)
// 每个document.WordCount已经同步好，无需查询DocumentContent
```

**特点**：
- ✅ 轻量级，查询快速
- ✅ 不包含实际内容
- ✅ 统计信息冗余（性能优化）
- ✅ 支持设定关联

**使用场景**：
- 文档列表查询
- 文档元数据管理
- 设定关联查询
- 统计信息展示

---

### 第四层：DocumentContent（文档内容层）

**职责**：管理文档的实际内容，支持GridFS大文本存储

```go
type DocumentContent struct {
    ID            string    `bson:"_id,omitempty" json:"id"`
    DocumentID    string    `bson:"document_id" json:"documentId"` // 关联到Document（一对一）
    
    // 实际内容（小文件 <1MB）
    Content       string    `bson:"content,omitempty" json:"content"`
    ContentType   string    `bson:"content_type" json:"contentType"` // markdown | richtext | html
    
    // 大文件支持（>1MB）
    GridFSID      string    `bson:"gridfs_id,omitempty" json:"gridfsId"`
    FileSize      int64     `bson:"file_size" json:"fileSize"` // 字节数
    
    // 统计信息
    WordCount     int       `bson:"word_count" json:"wordCount"`
    CharCount     int       `bson:"char_count" json:"charCount"`
    
    // 版本控制（乐观锁）
    Version       int       `bson:"version" json:"version"`
    
    // 自动保存信息
    LastSavedAt   time.Time `bson:"last_saved_at" json:"lastSavedAt"`
    IsAutoSave    bool      `bson:"is_auto_save" json:"isAutoSave"`
    
    CreatedAt     time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt     time.Time `bson:"updated_at" json:"updatedAt"`
}
```

**大文件存储策略**：

```go
func (s *DocumentContentService) SaveContent(documentID, content string) error {
    contentSize := len([]byte(content))
    
    if contentSize > 1*1024*1024 { // 1MB阈值
        // 大文件：存储到GridFS
        gridfsID, err := s.gridfsRepo.Upload(content)
        if err != nil {
            return err
        }
        
        return s.repo.Save(&DocumentContent{
            DocumentID: documentID,
            Content:    "",              // 内容为空
            GridFSID:   gridfsID,        // GridFS文件ID
            FileSize:   int64(contentSize),
            WordCount:  countWords(content),
        })
    } else {
        // 小文件：直接存储
        return s.repo.Save(&DocumentContent{
            DocumentID: documentID,
            Content:    content,         // 直接存储
            GridFSID:   "",
            FileSize:   int64(contentSize),
            WordCount:  countWords(content),
        })
    }
}
```

**读取内容**：

```go
func (s *DocumentContentService) GetContent(documentID string) (string, error) {
    docContent, err := s.repo.GetByDocumentID(documentID)
    if err != nil {
        return "", err
    }
    
    if docContent.GridFSID != "" {
        // 从GridFS读取
        return s.gridfsRepo.Download(docContent.GridFSID)
    } else {
        // 直接返回
        return docContent.Content, nil
    }
}
```

**特点**：
- ✅ 支持小文件直接存储
- ✅ 支持大文件GridFS存储
- ✅ 自动选择存储方式
- ✅ 乐观锁版本控制
- ✅ 自动保存支持

**使用场景**：
- 编辑器加载内容
- 自动保存（30秒间隔）
- 手动保存
- 大文本小说支持（>100万字）

---

## 🔄 四层协作流程

### 场景1：创建新项目

```
1. ProjectService.CreateProject()
   ↓ 创建Project记录
   ↓
2. NodeService.CreateRootNode()
   ↓ 创建根节点（type=root）
   ↓
3. 返回Project和RootNode信息
```

**事务保证**：
```go
func (s *ProjectService) CreateProject(req *CreateProjectRequest) error {
    return s.mongoClient.UseSession(ctx, func(sc mongo.SessionContext) error {
        sc.StartTransaction()
        
        // 1. 创建Project
        project := &Project{...}
        if err := s.projectRepo.Create(sc, project); err != nil {
            sc.AbortTransaction(sc)
            return err
        }
        
        // 2. 创建根节点
        rootNode := &Node{
            ProjectID: project.ID,
            Type:      NodeTypeRoot,
            Level:     0,
            Title:     project.Name,
        }
        if err := s.nodeRepo.Create(sc, rootNode); err != nil {
            sc.AbortTransaction(sc)
            return err
        }
        
        // 3. 更新Project的RootNodeID
        project.RootNodeID = rootNode.ID
        if err := s.projectRepo.Update(sc, project); err != nil {
            sc.AbortTransaction(sc)
            return err
        }
        
        return sc.CommitTransaction(sc)
    })
}
```

---

### 场景2：创建新章节（含内容）

```
1. NodeService.CreateNode()
   ↓ 创建Node（type=chapter, level=2）
   ↓ 更新父节点的ChildCount
   ↓
2. DocumentService.CreateDocument()
   ↓ 创建Document（关联NodeID）
   ↓
3. DocumentContentService.CreateContent()
   ↓ 创建DocumentContent（关联DocumentID）
   ↓ 根据大小选择：直接存储 or GridFS
   ↓
4. 返回完整的章节信息
```

**代码示例**：
```go
func (api *ChapterApi) CreateChapter(c *gin.Context) {
    var req CreateChapterRequest
    c.ShouldBindJSON(&req)
    
    // 1. 创建Node
    node, err := api.nodeService.CreateNode(&NodeRequest{
        ProjectID: req.ProjectID,
        ParentID:  req.ParentID,
        Type:      NodeTypeChapter,
        Title:     req.Title,
    })
    
    // 2. 创建Document
    document, err := api.documentService.CreateDocument(&DocumentRequest{
        ProjectID: req.ProjectID,
        NodeID:    node.ID,
        Title:     req.Title,
    })
    
    // 3. 创建DocumentContent（如果有初始内容）
    if req.Content != "" {
        err := api.contentService.SaveContent(document.ID, req.Content)
    }
    
    response.Success(c, node)
}
```

---

### 场景3：查询文档树（高频操作）

```
GET /api/v1/projects/:id/tree

只查询两层：
1. 查询所有Node（按level和order排序）
   ↓ db.nodes.find({project_id: "xxx"}).sort({level: 1, order: 1})
   ↓ 返回树形结构（不含content）
   ↓
2. （可选）查询Document元数据
   ↓ 批量查询：db.documents.find({node_id: {$in: nodeIDs}})
   ↓ 返回WordCount等统计信息
```

**性能优化**：
- ✅ 只查询Node和Document（轻量级）
- ✅ 不查询DocumentContent（重量级）
- ✅ 响应时间：<300ms（100个节点）

**返回数据示例**：
```json
{
  "projectId": "proj123",
  "tree": [
    {
      "nodeId": "node1",
      "type": "volume",
      "title": "第一卷",
      "level": 1,
      "order": 1,
      "wordCount": 50000,  // 从Document获取
      "children": [
        {
          "nodeId": "node2",
          "type": "chapter",
          "title": "第一章",
          "level": 2,
          "order": 1,
          "wordCount": 5000
        }
      ]
    }
  ]
}
```

---

### 场景4：编辑文档内容

```
用户点击章节 → 加载编辑器

1. GET /api/v1/documents/:id/content
   ↓
2. DocumentContentService.GetContent(documentID)
   ↓
3. 根据GridFSID判断：
   - 有GridFSID → 从GridFS读取
   - 无GridFSID → 直接返回Content字段
   ↓
4. 返回内容给编辑器
```

**代码示例**：
```go
func (api *DocumentApi) GetContent(c *gin.Context) {
    documentID := c.Param("id")
    
    // 查询DocumentContent
    content, err := api.contentService.GetContent(documentID)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    // 查询Document元数据（可选）
    document, _ := api.documentService.GetByID(documentID)
    
    response.Success(c, gin.H{
        "content":     content,
        "wordCount":   document.WordCount,
        "version":     document.Version,
        "lastSavedAt": document.UpdatedAt,
    })
}
```

---

### 场景5：自动保存（30秒间隔）

```
编辑器自动保存

1. POST /api/v1/documents/:id/autosave
   ↓
2. DocumentContentService.AutoSave()
   ↓ 使用乐观锁更新DocumentContent
   ↓ 更新WordCount
   ↓
3. DocumentService.UpdateWordCount()
   ↓ 同步WordCount到Document（冗余字段）
   ↓
4. （不创建Version记录）
```

**乐观锁示例**：
```go
func (s *DocumentContentService) AutoSave(documentID, content string, version int) error {
    wordCount := countWords(content)
    
    // 乐观锁更新
    result, err := s.repo.UpdateWithVersion(documentID, version, &DocumentContent{
        Content:    content,
        WordCount:  wordCount,
        Version:    version + 1,
        IsAutoSave: true,
        LastSavedAt: time.Now(),
    })
    
    if result.MatchedCount == 0 {
        return errors.New("版本冲突，内容已被其他用户修改")
    }
    
    // 同步WordCount到Document
    s.documentService.UpdateWordCount(documentID, wordCount)
    
    return nil
}
```

---

### 场景6：拖拽排序节点

```
用户拖拽节点调整顺序

1. PUT /api/v1/nodes/:id/reorder
   ↓
2. NodeService.Reorder()
   ↓ 更新同级节点的Order字段
   ↓
3. 返回更新后的节点列表
```

**批量更新Order**：
```go
func (s *NodeService) Reorder(nodeID, newParentID string, newOrder int) error {
    // 1. 获取节点
    node, _ := s.repo.GetByID(nodeID)
    
    // 2. 如果父节点改变，更新Path和Level
    if node.ParentID != newParentID {
        newParent, _ := s.repo.GetByID(newParentID)
        node.ParentID = newParentID
        node.Level = newParent.Level + 1
        node.Path = newParent.Path + "/" + nodeID
    }
    
    // 3. 更新Order（需要调整同级节点）
    siblings, _ := s.repo.GetChildNodes(newParentID)
    for i, sibling := range siblings {
        if sibling.ID == nodeID {
            sibling.Order = newOrder
        } else if sibling.Order >= newOrder {
            sibling.Order += 1
        }
        s.repo.Update(sibling)
    }
    
    return nil
}
```

---

## 📋 数据库设计

### 集合结构

#### projects集合（轻量）
```javascript
{
  "_id": "proj123",
  "owner_id": "user456",
  "name": "我的小说",
  "status": "active",
  "root_node_id": "node_root",
  "total_words": 100000,
  "total_documents": 50,
  "created_at": ISODate("2025-01-01T00:00:00Z")
}
```

#### nodes集合（树形结构）
```javascript
{
  "_id": "node2",
  "project_id": "proj123",
  "parent_id": "node1",       // null表示根节点
  "level": 2,                 // 层级
  "order": 1,                 // 同级排序
  "path": "/node1/node2",     // 路径（快速查询祖先）
  "type": "chapter",          // root | volume | chapter | section
  "document_id": "doc5",      // 关联文档（叶子节点）
  "title": "第一章",
  "child_count": 3,
  "created_at": ISODate
}
```

#### documents集合（元数据 - 轻量）
```javascript
{
  "_id": "doc5",
  "project_id": "proj123",
  "node_id": "node2",
  "title": "第一章",
  "status": "draft",
  "tags": ["玄幻", "修仙"],
  "character_ids": ["char1", "char2"],
  "word_count": 5000,         // 冗余字段
  "version": 3,               // 乐观锁
  "created_at": ISODate
}
```

#### document_contents集合（实际内容 - 重量）
```javascript
{
  "_id": "content10",
  "document_id": "doc5",
  
  // 小文件（<1MB）
  "content": "章节内容...",   // 实际内容
  "gridfs_id": "",
  
  // 大文件（>1MB）时：
  // "content": "",
  // "gridfs_id": "gridfs_abc123",
  
  "content_type": "markdown",
  "word_count": 5000,
  "file_size": 50000,         // 字节数
  "version": 3,
  "is_auto_save": false,
  "last_saved_at": ISODate
}
```

---

### 索引设计

```javascript
// projects集合
db.projects.createIndex({"owner_id": 1, "status": 1})
db.projects.createIndex({"created_at": -1})

// nodes集合（树形查询优化）
db.nodes.createIndex({"project_id": 1, "level": 1, "order": 1})
db.nodes.createIndex({"parent_id": 1, "order": 1})
db.nodes.createIndex({"path": 1})  // 祖先查询

// documents集合
db.documents.createIndex({"project_id": 1})
db.documents.createIndex({"node_id": 1})
db.documents.createIndex({"character_ids": 1})  // 多值索引

// document_contents集合
db.document_contents.createIndex({"document_id": 1}, {"unique": true})
db.document_contents.createIndex({"last_saved_at": -1})
```

---

## 🔌 API设计

### Project API

```http
# 创建项目
POST /api/v1/projects
{
  "name": "我的小说",
  "type": "novel"
}

# 查询项目列表
GET /api/v1/projects?status=active

# 查询单个项目
GET /api/v1/projects/:id

# 更新项目
PUT /api/v1/projects/:id

# 删除项目（软删除）
DELETE /api/v1/projects/:id
```

### Node API

```http
# 查询文档树
GET /api/v1/projects/:id/tree

# 创建节点
POST /api/v1/nodes
{
  "projectId": "proj123",
  "parentId": "node1",
  "type": "chapter",
  "title": "第一章"
}

# 拖拽排序
PUT /api/v1/nodes/:id/reorder
{
  "newParentId": "node2",
  "newOrder": 3
}

# 删除节点（级联删除子节点）
DELETE /api/v1/nodes/:id
```

### Document API

```http
# 创建文档
POST /api/v1/documents
{
  "nodeId": "node2",
  "title": "第一章"
}

# 查询文档元数据
GET /api/v1/documents/:id

# 更新元数据
PUT /api/v1/documents/:id
{
  "tags": ["玄幻"],
  "characterIds": ["char1"]
}
```

### DocumentContent API

```http
# 获取文档内容
GET /api/v1/documents/:id/content

# 自动保存
POST /api/v1/documents/:id/autosave
{
  "content": "更新的内容...",
  "version": 3
}

# 手动保存（创建版本）
PUT /api/v1/documents/:id/content
{
  "content": "更新的内容...",
  "version": 3,
  "comment": "完成第一章修改"
}
```

---

## ✅ 架构优势总结

### 1. 性能优化 🚀

| 操作 | 三层架构 | 四层架构 | 提升 |
|------|---------|---------|------|
| 查询文档树 | 5-10秒 | <300ms | **20-30倍** |
| 加载编辑器 | 3-5秒 | <500ms | **6-10倍** |
| 自动保存 | 1秒 | <200ms | **5倍** |

### 2. 可扩展性 📈

- ✅ 支持超大文本（>10MB，GridFS自动处理）
- ✅ 支持100万+文档
- ✅ 支持1000+并发编辑

### 3. 数据安全 🔒

- ✅ 乐观锁防止并发冲突
- ✅ 软删除支持数据恢复
- ✅ 版本历史完整保留

### 4. 开发友好 👨‍💻

- ✅ 职责分离清晰
- ✅ 按需加载
- ✅ 易于测试

---

## 📚 相关文档

- [核心功能设计入口](../README.md) - Document/Content 分层与核心基础设施入口
- [后端架构总览](../../../architecture/README.md) - 当前实现侧架构入口
- [版本控制](./版本控制.md) - 版本管理详细设计
- [后端架构总览](../../../architecture/README.md) - 数据访问层规范与当前分层入口

---

**文档版本**: v1.0  
**最后更新**: 2025-10-21  
**维护者**: 青羽架构组

