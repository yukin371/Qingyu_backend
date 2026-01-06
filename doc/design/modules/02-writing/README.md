# 02 - 写作创作模块

> **模块编号**: 02
> **模块名称**: Writing & Creation
> **负责功能**: 项目管理、文档编辑、版本控制、AI辅助写作
> **完成度**: 🟡 70%

## 📋 目录结构

```
写作创作模块/
├── api/v1/
│   ├── writer/                   # 写作API
│   │   ├── project_api.go       # 项目管理
│   │   ├── document_api.go      # 文档管理
│   │   ├── version_api.go       # 版本控制
│   │   └── editor_api.go        # 编辑器功能
│   └── projects/                 # 项目API（兼容路由）
├── service/writer/               # 写作服务层
│   ├── project_service.go       # 项目服务
│   ├── document_service.go      # 文档服务
│   ├── version_service.go       # 版本服务
│   └── editor_service.go        # 编辑器服务
├── repository/interfaces/writer/ # 仓储接口
├── repository/mongodb/writer/    # MongoDB仓储实现
│   ├── project_repository_mongo.go
│   ├── document_repository_mongo.go
│   └── version_repository_mongo.go
└── models/writer/                # 数据模型
    ├── project.go                # 项目实体
    ├── document.go               # 文档实体
    ├── version.go                # 版本实体
    ├── editor.go                 # 编辑器状态
    └── statistics.go             # 统计数据
```

## 🎯 核心功能

### 1. 项目管理

- **创建项目**: 支持多种分类（小说、散文、诗歌等）
- **项目编辑**: 标题、描述、封面、标签
- **协作管理**: 添加协作者、权限分配
- **项目统计**: 字数、章节数、阅读量
- **状态管理**: 草稿、连载、完结

### 2. 文档管理

- **层级结构**: 支持卷/章节的树形结构
- **文档类型**: 卷、章节、前言、后记、番外
- **文档状态**: 草稿、审核中、已发布
- **拖拽排序**: 支持文档重新排序
- **批量操作**: 批量移动、删除

### 3. 编辑器功能

- **自动保存**: 每30秒自动保存
- **文档锁定**: 防止多人同时编辑
- **字数统计**: 实时字数、字符数统计
- **敏感词检测**: 自动检测违规内容
- **快捷键**: 自定义快捷键

### 4. 版本控制

- **版本历史**: 记录每次修改
- **版本对比**: Diff功能对比差异
- **版本回滚**: 恢复到历史版本
- **版本标签**: 标记重要版本

### 5. AI辅助

- **AI续写**: 根据上下文续写内容
- **AI改写**: 改写段落表达
- **AI扩写**: 扩展段落内容
- **AI润色**: 优化文字表达
- **AI大纲**: 生成章节大纲

## 📊 数据模型

### Project (项目实体)

```go
type Project struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    Title           string               `bson:"title" json:"title"`
    Description     string               `bson:"description" json:"description"`
    Cover           string               `bson:"cover" json:"cover"`
    Category        ProjectCategory      `bson:"category" json:"category"`
    Status          ProjectStatus        `bson:"status" json:"status"`
    Tags            []string             `bson:"tags" json:"tags"`
    IsPublic        bool                 `bson:"is_public" json:"isPublic"`
    IsCompleted     bool                 `bson:"is_completed" json:"isCompleted"`

    // 协作信息
    Collaborators   []Collaborator       `bson:"collaborators" json:"collaborators"`

    // 统计信息
    WordCount       int64                `bson:"word_count" json:"wordCount"`
    ChapterCount    int                  `bson:"chapter_count" json:"chapterCount"`
    ReadingCount    int64                `bson:"reading_count" json:"readingCount"`
    CommentCount    int                  `bson:"comment_count" json:"commentCount"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
    PublishedAt     *time.Time           `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
    CompletedAt     *time.Time           `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
}

type ProjectCategory string
const (
    CategoryNovel     ProjectCategory = "novel"
    CategoryShort     ProjectCategory = "short"
    CategoryProse     ProjectCategory = "prose"
    CategoryPoetry    ProjectCategory = "poetry"
    CategoryScript    ProjectCategory = "script"
    CategoryOther     ProjectCategory = "other"
)

type ProjectStatus string
const (
    StatusDraft      ProjectStatus = "draft"
    StatusOngoing    ProjectStatus = "ongoing"
    StatusCompleted  ProjectStatus = "completed"
    StatusPaused     ProjectStatus = "paused"
)
```

### Document (文档实体)

```go
type Document struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    ProjectID       primitive.ObjectID   `bson:"project_id" json:"projectId"`
    ParentID        *primitive.ObjectID  `bson:"parent_id,omitempty" json:"parentId,omitempty"`
    Title           string               `bson:"title" json:"title"`
    Content         string               `bson:"content" json:"content"`
    ContentType     ContentType          `bson:"content_type" json:"contentType"`
    DocumentType    DocumentType         `bson:"document_type" json:"documentType"`
    Order           int                  `bson:"order" json:"order"`
    Depth           int                  `bson:"depth" json:"depth"`
    Status          DocumentStatus       `bson:"status" json:"status"`

    // 权限控制
    IsLocked        bool                 `bson:"is_locked" json:"isLocked"`
    LockedBy        *primitive.ObjectID  `bson:"locked_by,omitempty" json:"lockedBy,omitempty"`
    LockedAt        *time.Time           `bson:"locked_at,omitempty" json:"lockedAt,omitempty"`

    // 版本控制
    CurrentVersion  int                  `bson:"current_version" json:"currentVersion"`

    // 统计信息
    WordCount       int                  `bson:"word_count" json:"wordCount"`
    CharCount       int                  `bson:"char_count" json:"charCount"`
    ReadingTime     int                  `bson:"reading_time" json:"readingTime"`

    // 发布信息
    IsPublished     bool                 `bson:"is_published" json:"isPublished"`
    PublishedAt     *time.Time           `bson:"published_at,omitempty" json:"publishedAt,omitempty"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
    DeletedAt       *time.Time           `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}
```

### Version (版本实体)

```go
type Version struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    DocumentID      primitive.ObjectID   `bson:"document_id" json:"documentId"`
    VersionNumber   int                  `bson:"version_number" json:"versionNumber"`
    Content         string               `bson:"content" json:"content"`
    WordCount       int                  `bson:"word_count" json:"wordCount"`
    CharCount       int                  `bson:"char_count" json:"charCount"`
    ChangeSummary   string               `bson:"change_summary" json:"changeSummary"`
    IsTagged        bool                 `bson:"is_tagged" json:"isTagged"`
    Tag             string               `bson:"tag,omitempty" json:"tag,omitempty"`
    CreatedBy       primitive.ObjectID   `bson:"created_by" json:"createdBy"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}
```

## 🔄 核心流程

### 项目创建流程

```
1. 作者点击"创建项目"
   ↓
2. 填写项目信息（标题、描述、分类等）
   ↓
3. 选择协作设置（公开/私有）
   ↓
4. 后端验证并创建项目
   ↓
5. 初始化项目文档结构（创建默认卷/章节）
   ↓
6. 返回项目ID
```

### 文档编辑流程

```
1. 用户打开文档
   ↓
2. 后端检查文档锁定状态
   ↓
3. 如果未锁定，锁定文档（锁定时间：30分钟）
   ↓
4. 返回文档内容
   ↓
5. 用户编辑
   ↓
6. 自动保存机制触发（每30秒）
   ↓
7. 保存新版本
   ↓
8. 更新字数统计
   ↓
9. 用户关闭文档，释放锁定
```

### 版本控制流程

```
1. 用户修改文档并保存
   ↓
2. 系统创建新版本
   ↓
3. 版本号递增
   ↓
4. 保存版本快照
   ↓
5. 用户可查看版本历史
   ↓
6. 用户可选择版本对比
   ↓
7. 用户可恢复到任意版本
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/writer/projects | 创建项目 | 是 |
| GET | /api/v1/writer/projects | 获取项目列表 | 是 |
| GET | /api/v1/writer/projects/:id | 获取项目详情 | 是 |
| PUT | /api/v1/writer/projects/:id | 更新项目 | 是 |
| DELETE | /api/v1/writer/projects/:id | 删除项目 | 是 |
| POST | /api/v1/writer/project/:projectId/documents | 创建文档 | 是 |
| GET | /api/v1/writer/project/:projectId/documents | 获取文档列表 | 是 |
| GET | /api/v1/writer/documents/:id | 获取文档详情 | 是 |
| PUT | /api/v1/writer/documents/:id | 更新文档 | 是 |
| DELETE | /api/v1/writer/documents/:id | 删除文档 | 是 |
| PUT | /api/v1/writer/documents/:id/move | 移动文档 | 是 |
| POST | /api/v1/writer/documents/:id/autosave | 自动保存 | 是 |
| GET | /api/v1/writer/documents/:id/versions | 获取版本历史 | 是 |
| GET | /api/v1/writer/document/:documentId/versions/:versionId | 获取特定版本 | 是 |
| POST | /api/v1/writer/document/:documentId/versions/compare | 版本对比 | 是 |
| POST | /api/v1/writer/document/:documentId/versions/:versionId/restore | 恢复版本 | 是 |

## 🤝 协作功能

### 协作者角色权限

| 角色 | 读取 | 编辑 | 删除 | 管理 |
|------|------|------|------|------|
| Owner | ✓ | ✓ | ✓ | ✓ |
| Editor | ✓ | ✓ | ✗ | ✗ |
| Viewer | ✓ | ✗ | ✗ | ✗ |

### 实时协作机制

```
1. 文档锁定机制
   - 防止多人同时编辑同一文档
   - 锁定超时自动释放（30分钟）

2. 变更通知
   - 协作者更新通知
   - 评论和批注通知

3. 权限管理
   - 动态添加/移除协作者
   - 角色权限调整
```

## 🤖 AI集成

### AI写作辅助

- **续写**: 根据上下文续写内容
- **改写**: 改变文字表达方式
- **扩写**: 扩展段落内容
- **润色**: 优化文字表达
- **大纲生成**: 生成章节大纲

### AI配额管理

- 每日免费AI调用次数限制
- 不同用户等级不同配额
- 超出配额需购买

## 🔒 内容安全

### 敏感词检测

- 政治、暴力、色情等敏感内容
- 可配置敏感词库
- 检测结果可标记或自动过滤

### 内容审核

- 自动检测违规内容
- 人工审核流程
- 违规内容处理

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 用户身份验证
- **09 - AI模块**: AI辅助功能

### 被依赖的模块
- **03 - 阅读器**: 获取已发布内容
- **06 - 书城**: 获取作品信息

### 外部服务
- **AI服务**: 用于AI辅助写作
- **存储服务**: 封面图片存储

## 📈 扩展点

1. **协作功能增强**
   - 实时协作编辑（WebSocket + OT算法）
   - 评论和批注系统
   - 变更建议和审核流程

2. **导入导出**
   - 支持导入Word、TXT、Markdown
   - 导出为EPUB、PDF、MOBI

3. **模板系统**
   - 预设项目模板
   - 章节模板
   - 风格指南

4. **插件系统**
   - 自定义编辑器插件
   - 第三方工具集成

## 🚀 性能优化

1. **文档缓存**
   - Redis 缓存热点文档
   - CDN 分发文档内容

2. **自动保存优化**
   - 防抖处理，避免频繁保存
   - 批量保存变更

3. **搜索优化**
   - 全文索引（MongoDB Atlas Search）
   - 搜索结果缓存

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/writer/`
**相关设计**: [写作端模块设计文档](../../writing/)
