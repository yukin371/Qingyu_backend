# Writer API 模块结构说明

## 📁 文件结构

```
api/v1/writer/
├── project_api.go      # 项目管理API
├── document_api.go     # 文档管理API
├── editor_api.go       # 编辑器API
├── audit_api.go        # 审核API（作者端）
├── stats_api.go        # 统计API
├── version_api.go      # 版本管理API
├── types.go            # 公共DTO定义
└── README.md           # 本文件
```

## 🎯 模块职责划分

### 1. ProjectAPI (`project_api.go`)

**职责**: 写作项目管理

**核心功能**:
- ✅ 创建项目
- ✅ 获取项目列表
- ✅ 获取项目详情
- ✅ 更新项目
- ✅ 删除项目
- ✅ 项目设置管理

**API端点**:
```
POST   /api/v1/projects              # 创建项目
GET    /api/v1/projects              # 获取项目列表
GET    /api/v1/projects/:id          # 获取项目详情
PUT    /api/v1/projects/:id          # 更新项目
DELETE /api/v1/projects/:id          # 删除项目
GET    /api/v1/projects/:id/settings # 获取项目设置
PUT    /api/v1/projects/:id/settings # 更新项目设置
```

---

### 2. DocumentAPI (`document_api.go`)

**职责**: 文档/章节管理

**核心功能**:
- ✅ 创建文档
- ✅ 获取文档列表
- ✅ 获取文档详情
- ✅ 更新文档内容
- ✅ 删除文档
- ✅ 文档排序

**API端点**:
```
POST   /api/v1/projects/:id/documents        # 创建文档
GET    /api/v1/projects/:id/documents        # 获取文档列表
GET    /api/v1/documents/:id                 # 获取文档详情
PUT    /api/v1/documents/:id                 # 更新文档
DELETE /api/v1/documents/:id                 # 删除文档
PUT    /api/v1/documents/:id/order           # 更新文档顺序
POST   /api/v1/documents/:id/publish         # 发布文档
```

---

### 3. EditorAPI (`editor_api.go`)

**职责**: 编辑器功能

**核心功能**:
- ✅ 自动保存
- ✅ 保存草稿
- ✅ 字数统计
- ✅ 敏感词检测
- ✅ 锁定文档（防止并发编辑）

**API端点**:
```
POST /api/v1/documents/:id/autosave      # 自动保存
POST /api/v1/documents/:id/draft         # 保存草稿
GET  /api/v1/documents/:id/wordcount     # 字数统计
POST /api/v1/documents/:id/check         # 敏感词检测
POST /api/v1/documents/:id/lock          # 锁定文档
POST /api/v1/documents/:id/unlock        # 解锁文档
```

---

### 4. AuditAPI (`audit_api.go`)

**职责**: 审核功能（作者端）

**核心功能**:
- ✅ 提交审核
- ✅ 查看审核结果
- ✅ 查看违规记录
- ✅ 申诉审核结果

**API端点**:
```
POST /api/v1/audit/check                     # 实时检测内容
POST /api/v1/documents/:id/audit             # 全文审核文档
GET  /api/v1/documents/:id/audit             # 获取审核结果
POST /api/v1/audit/:id/appeal                # 申诉审核结果
GET  /api/v1/users/:userId/violations        # 获取用户违规记录
GET  /api/v1/users/:userId/violations/summary # 获取用户违规统计
```

**注意**: 管理员审核功能已迁移到 `admin` 模块。

---

### 5. StatsAPI (`stats_api.go`)

**职责**: 统计数据

**核心功能**:
- ✅ 项目统计
- ✅ 文档统计
- ✅ 阅读数据统计
- ✅ 收入统计

**API端点**:
```
GET /api/v1/projects/:id/stats           # 项目统计
GET /api/v1/documents/:id/stats          # 文档统计
GET /api/v1/writer/stats/overview        # 作者总览统计
GET /api/v1/writer/stats/income          # 收入统计
GET /api/v1/writer/stats/readers         # 读者统计
```

---

### 6. VersionAPI (`version_api.go`)

**职责**: 版本控制

**核心功能**:
- ✅ 创建版本
- ✅ 版本历史
- ✅ 版本对比
- ✅ 恢复版本

**API端点**:
```
POST /api/v1/documents/:id/versions        # 创建版本
GET  /api/v1/documents/:id/versions        # 版本历史
GET  /api/v1/documents/:id/versions/:vid   # 获取版本详情
POST /api/v1/documents/:id/versions/:vid/restore # 恢复版本
GET  /api/v1/documents/:id/versions/compare      # 版本对比
```

---

## 🔄 API调用流程

### 标准流程（需要认证）
```
客户端请求 
  → Router 
  → JWTAuth中间件（验证Token）
  → API Handler 
  → Service层 
  → Repository层 
  → 数据库
```

### 审核流程
```
作者编辑内容 
  → 实时检测（快速反馈）
  → 保存草稿 
  → 提交审核（全文审核）
  → 审核服务分析 
  → 返回审核结果
  → （如有问题）作者修改或申诉
  → 管理员复核
```

---

## 🛡️ 中间件配置

### 1. JWT认证中间件
所有Writer接口都需要JWT认证：
```go
writerGroup.Use(middleware.JWTAuth())
```

### 2. 作者权限中间件
某些接口需要作者角色：
```go
writerGroup.Use(middleware.RequireRole("author", "admin"))
```

---

## 📊 请求/响应示例

### 创建项目
```json
POST /api/v1/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "我的小说",
  "description": "一个精彩的故事",
  "category": "玄幻",
  "tags": ["修仙", "热血"],
  "settings": {
    "autoSave": true,
    "autoSaveInterval": 30
  }
}

Response:
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "project_id": "proj_123",
    "title": "我的小说",
    "created_at": "2025-10-24T10:00:00Z"
  }
}
```

### 创建文档
```json
POST /api/v1/projects/proj_123/documents
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "第一章：开端",
  "content": "故事从这里开始...",
  "order": 1
}

Response:
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "document_id": "doc_456",
    "title": "第一章：开端",
    "word_count": 15
  }
}
```

### 自动保存
```json
POST /api/v1/documents/doc_456/autosave
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "故事从这里开始...（更新的内容）",
  "cursor_position": 120
}

Response:
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "saved_at": "2025-10-24T10:05:30Z",
    "word_count": 25
  }
}
```

### 提交审核
```json
POST /api/v1/documents/doc_456/audit
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "完整的章节内容..."
}

Response:
{
  "code": 200,
  "message": "审核完成",
  "data": {
    "audit_id": "audit_789",
    "status": "approved",
    "risk_level": 0,
    "violations": []
  }
}
```

### 查看统计
```json
GET /api/v1/projects/proj_123/stats
Authorization: Bearer <token>

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "total_documents": 10,
    "total_words": 50000,
    "total_views": 1000,
    "total_favorites": 50,
    "today_views": 100,
    "today_words": 2000
  }
}
```

---

## 🔧 设计原则

### 1. 以作者为中心
所有功能围绕作者创作体验设计。

### 2. 自动保存优先
防止数据丢失，提供稳定的编辑体验。

### 3. 实时反馈
敏感词检测、字数统计等功能实时反馈。

### 4. 版本管理
支持版本历史和回滚，保护创作成果。

### 5. 审核透明
清晰的审核结果和申诉流程。

---

## 📝 开发规范

### 1. 命名规范
- **API结构体**：`<功能>Api`（如 `ProjectApi`、`DocumentApi`）
- **构造函数**：`New<功能>Api`
- **方法名**：动词+名词（如 `CreateProject`、`UpdateDocument`）

### 2. 错误处理
```go
if err != nil {
    shared.Error(c, http.StatusInternalServerError, "操作失败", err.Error())
    return
}
```

### 3. 获取当前用户
```go
userID, ok := shared.GetUserID(c)
if !ok {
    return
}
```

### 4. 权限验证
只能操作自己的项目和文档：
```go
if project.AuthorID != userID {
    shared.Error(c, http.StatusForbidden, "无权限", "只能操作自己的项目")
    return
}
```

---

## 🚀 扩展建议

### 未来可添加的功能

1. **协作功能**
   - 多人协作编辑
   - 评论和批注
   - 权限管理

2. **AI辅助**
   - 智能续写
   - 文本润色
   - 剧情建议

3. **数据分析**
   - 读者画像
   - 阅读热力图
   - 章节质量分析

4. **发布管理**
   - 定时发布
   - 发布预览
   - 多平台同步

5. **素材管理**
   - 角色卡片
   - 世界观设定
   - 大纲管理

---

## 🔄 与其他模块的关系

### Writer vs Reader
| 功能 | Writer（写作端） | Reader（阅读端） |
|------|-----------------|-----------------|
| **定位** | 内容创作 | 内容消费 |
| **用户** | 作者 | 读者 |
| **核心功能** | 编辑、审核、统计 | 阅读、进度、标注 |
| **数据** | 项目、文档草稿 | 已发布内容、阅读记录 |

### Writer vs Bookstore
| 功能 | Writer（写作端） | Bookstore（书城） |
|------|-----------------|------------------|
| **定位** | 创作管理 | 展示发现 |
| **视角** | 作者视角 | 读者视角 |
| **内容** | 草稿和未发布 | 已发布和上架 |

### Writer vs Admin
| 功能 | Writer（作者端） | Admin（管理端） |
|------|-----------------|----------------|
| **审核** | 提交审核、查看结果 | 审核内容、处理申诉 |
| **统计** | 个人统计 | 全站统计 |
| **权限** | 管理自己的内容 | 管理所有内容 |

---

## 📚 相关文档

- [项目服务设计](../../../doc/design/project/README.md)
- [文档管理设计](../../../doc/design/document/README.md)
- [审核服务设计](../../../doc/design/audit/README.md)
- [编辑器API文档](../../../doc/api/编辑器API文档.md)
- [写作端API完整文档](../../../doc/api/写作端API完整文档.md)

---

## 📋 API端点总览

### 项目管理
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/projects | 创建项目 |
| GET | /api/v1/projects | 获取项目列表 |
| GET | /api/v1/projects/:id | 获取项目详情 |
| PUT | /api/v1/projects/:id | 更新项目 |
| DELETE | /api/v1/projects/:id | 删除项目 |

### 文档管理
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/projects/:id/documents | 创建文档 |
| GET | /api/v1/projects/:id/documents | 获取文档列表 |
| GET | /api/v1/documents/:id | 获取文档详情 |
| PUT | /api/v1/documents/:id | 更新文档 |
| DELETE | /api/v1/documents/:id | 删除文档 |

### 编辑器
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/documents/:id/autosave | 自动保存 |
| POST | /api/v1/documents/:id/draft | 保存草稿 |
| GET | /api/v1/documents/:id/wordcount | 字数统计 |
| POST | /api/v1/documents/:id/check | 敏感词检测 |

### 审核（作者端）
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/audit/check | 实时检测 |
| POST | /api/v1/documents/:id/audit | 提交审核 |
| GET | /api/v1/documents/:id/audit | 获取审核结果 |
| POST | /api/v1/audit/:id/appeal | 申诉 |

### 统计
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/projects/:id/stats | 项目统计 |
| GET | /api/v1/documents/:id/stats | 文档统计 |
| GET | /api/v1/writer/stats/overview | 总览统计 |

### 版本管理
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/documents/:id/versions | 创建版本 |
| GET | /api/v1/documents/:id/versions | 版本历史 |
| POST | /api/v1/documents/:id/versions/:vid/restore | 恢复版本 |

---

**版本**: v1.0  
**创建日期**: 2025-10-24  
**维护者**: Writer模块开发组

