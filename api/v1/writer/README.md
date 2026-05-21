# Writer API 模块结构说明

## 文件结构

```
api/v1/writer/
├── project_api.go                   # 项目管理
├── document_api.go                  # 文档管理
├── editor_api.go                    # 编辑器功能（自动保存、内容编辑、字数统计、快捷键）
├── version_api.go                   # 版本控制
├── publish_api.go                   # 发布管理
├── comment_api.go                   # 批注系统
├── template_api.go                  # 模板管理
├── character_api.go                 # 角色管理
├── location_api.go                  # 地点管理
├── timeline_api.go                  # 时间线管理
├── outline_api.go                   # 大纲管理
├── entity_api.go                    # 统一实体管理（角色/地点/物品/组织/伏笔）
├── keyword_api.go                   # 关键词检索（角色/地点前缀补全）
├── encyclopedia_api.go              # 设定百科（概念词条）
├── lock_api.go                      # 文档锁定（并发编辑控制）
├── batch_operation_api.go           # 批量操作（删除/移动/复制）
├── search_api.go                    # 文档搜索
├── export_api.go                    # 导出管理（TXT/MD/DOCX）
├── import_export_api.go             # 导入导出（ZIP 项目包）
├── story_harness_api.go             # Story Harness（章节上下文/索引）
├── change_request_api.go            # 变更建议（章节校对建议）
├── dashboard_api.go                 # 作者仪表板统计
├── stats_api.go                     # 阅读/书店统计
├── writer_stats_aggregate_api.go    # 作者聚合统计
├── audit_api.go                     # 内容审核
├── types.go                         # 公共响应类型定义
├── *_test.go                        # 单元测试
└── README.md                        # 本文件
```

## 模块职责划分

### 1. ProjectApi (`project_api.go`)

**职责**: 写作项目的全生命周期管理

**核心功能**:
- ✅ 创建项目
- ✅ 获取项目详情
- ✅ 获取项目列表（支持分页与筛选）
- ✅ 更新项目
- ✅ 删除项目
- ✅ 更新项目统计信息

**API端点**:
```
POST   /api/v1/writer/projects                    # 创建项目
GET    /api/v1/writer/projects                    # 获取项目列表
GET    /api/v1/writer/projects/:id                # 获取项目详情
PUT    /api/v1/writer/projects/:id                # 更新项目
DELETE /api/v1/writer/projects/:id                # 删除项目
PUT    /api/v1/writer/projects/:id/statistics     # 更新项目统计信息
```

**依赖服务**: `project.ProjectService`

---

### 2. DocumentApi (`document_api.go`)

**职责**: 项目内文档（章节）的增删改查与组织管理

**核心功能**:
- ✅ 创建文档（支持路径参数和请求体两种方式）
- ✅ 获取文档详情
- ✅ 获取文档树（层级结构）
- ✅ 更新文档
- ✅ 删除文档
- ✅ 移动文档
- ✅ 重新排序文档
- ✅ 复制文档

**API端点**:
```
POST   /api/v1/writer/project/:projectId/documents              # 创建文档
GET    /api/v1/writer/project/:projectId/documents              # 获取文档列表
GET    /api/v1/writer/project/:projectId/documents/tree         # 获取文档树
PUT    /api/v1/writer/project/:projectId/documents/reorder      # 重新排序
POST   /api/v1/writer/documents                                 # 创建文档（请求体方式）
GET    /api/v1/writer/documents/:id                             # 获取文档详情
PUT    /api/v1/writer/documents/:id                             # 更新文档
DELETE /api/v1/writer/documents/:id                             # 删除文档
PUT    /api/v1/writer/documents/:id/move                        # 移动文档
POST   /api/v1/writer/documents/:id/duplicate                   # 复制文档
```

**依赖服务**: `document.DocumentService`

---

### 3. EditorApi (`editor_api.go`)

**职责**: 编辑器核心功能，包括自动保存、内容读写、字数统计和快捷键配置

**核心功能**:
- ✅ 自动保存（支持版本冲突检测）
- ✅ 获取保存状态
- ✅ 获取/更新文档内容（单篇与批量）
- ✅ 重建内容索引
- ✅ 计算字数
- ✅ 用户快捷键配置（获取/更新/重置/帮助）

**API端点**:
```
POST   /api/v1/writer/documents/:id/autosave          # 自动保存
GET    /api/v1/writer/documents/:id/save-status       # 获取保存状态
GET    /api/v1/writer/documents/:id/content           # 获取文档内容
PUT    /api/v1/writer/documents/:id/content           # 更新文档内容
GET    /api/v1/writer/documents/:id/contents          # 批量获取内容
PUT    /api/v1/writer/documents/:id/contents          # 批量更新内容
POST   /api/v1/writer/documents/:id/contents/reindex  # 重建内容索引
POST   /api/v1/writer/documents/:id/word-count        # 计算字数
GET    /api/v1/writer/user/shortcuts                  # 获取快捷键配置
PUT    /api/v1/writer/user/shortcuts                  # 更新快捷键配置
POST   /api/v1/writer/user/shortcuts/reset            # 重置快捷键
GET    /api/v1/writer/user/shortcuts/help             # 获取快捷键帮助
```

**依赖服务**: `document.DocumentService`, `document.WordCountService`, `document.ShortcutService`

---

### 4. VersionApi (`version_api.go`)

**职责**: 文档版本历史管理与版本回溯

**核心功能**:
- ✅ 获取版本历史列表
- ✅ 获取特定版本内容
- ✅ 比较两个版本差异
- ✅ 恢复到指定版本

**API端点**:
```
GET    /api/v1/writer/document/:documentId/versions                         # 获取版本历史
GET    /api/v1/writer/document/:documentId/versions/:versionId              # 获取特定版本
GET    /api/v1/writer/document/:documentId/versions/compare                 # 比较版本
POST   /api/v1/writer/document/:documentId/versions/:versionId/restore      # 恢复版本
```

**依赖服务**: `project.VersionService`

---

### 5. PublishApi (`publish_api.go`)

**职责**: 作品发布到书城的全流程管理，包括项目发布、章节发布、发布记录查询

**核心功能**:
- ✅ 发布项目到书城
- ✅ 取消发布项目
- ✅ 获取项目发布状态
- ✅ 发布单个文档（章节）
- ✅ 更新文档发布状态
- ✅ 批量发布文档
- ✅ 获取发布记录列表
- ✅ 获取发布记录详情

**API端点**:
```
POST   /api/v1/writer/projects/:id/publish                      # 发布项目
POST   /api/v1/writer/projects/:id/unpublish                    # 取消发布
GET    /api/v1/writer/projects/:id/publication-status            # 获取发布状态
POST   /api/v1/writer/documents/:id/publish                     # 发布文档
PUT    /api/v1/writer/documents/:id/publish-status               # 更新发布状态
POST   /api/v1/writer/projects/:id/documents/batch-publish      # 批量发布
GET    /api/v1/writer/projects/:id/publications                  # 发布记录列表
GET    /api/v1/writer/publications/:id                           # 发布记录详情
```

**依赖服务**: `interfaces.PublishService`

---

### 6. CommentApi (`comment_api.go`)

**职责**: 文档批注系统，支持创建、回复、解决、搜索批注

**核心功能**:
- ✅ 创建批注
- ✅ 获取批注列表（支持分页/筛选）
- ✅ 获取批注详情
- ✅ 更新批注
- ✅ 删除批注
- ✅ 解决/取消解决批注
- ✅ 回复批注
- ✅ 获取批注线程
- ✅ 获取批注统计
- ✅ 搜索批注
- ✅ 批量删除批注

**API端点**:
```
POST   /api/v1/writer/documents/:id/comments              # 创建批注
GET    /api/v1/writer/documents/:id/comments              # 获取批注列表
GET    /api/v1/writer/documents/:id/comments/stats        # 批注统计
GET    /api/v1/writer/documents/:id/comments/search       # 搜索批注
GET    /api/v1/writer/comments/:id                        # 获取批注详情
PUT    /api/v1/writer/comments/:id                        # 更新批注
DELETE /api/v1/writer/comments/:id                        # 删除批注
POST   /api/v1/writer/comments/:id/resolve                # 解决批注
POST   /api/v1/writer/comments/:id/unresolve              # 取消解决
POST   /api/v1/writer/comments/:id/reply                  # 回复批注
GET    /api/v1/writer/comments/threads/:threadId           # 获取批注线程
POST   /api/v1/writer/comments/batch-delete                # 批量删除
```

**依赖服务**: `writerservice.CommentService`

---

### 7. TemplateApi (`template_api.go`)

**职责**: 文档模板的创建、管理和应用

**核心功能**:
- ✅ 创建模板
- ✅ 获取模板列表（支持分页与筛选）
- ✅ 获取模板详情
- ✅ 更新模板
- ✅ 删除模板
- ✅ 应用模板到文档

**API端点**:
```
POST   /api/v1/writer/templates              # 创建模板
GET    /api/v1/writer/templates              # 获取模板列表
GET    /api/v1/writer/templates/:id          # 获取模板详情
PUT    /api/v1/writer/templates/:id          # 更新模板
DELETE /api/v1/writer/templates/:id          # 删除模板
POST   /api/v1/writer/templates/:id/apply    # 应用模板
```

**依赖服务**: `document.TemplateService`

---

### 8. CharacterApi (`character_api.go`)

**职责**: 角色管理，包括角色 CRUD、角色关系、关系时序变化事件

**核心功能**:
- ✅ 创建/获取/更新/删除角色
- ✅ 获取项目角色列表
- ✅ 创建/删除角色关系
- ✅ 获取角色关系列表
- ✅ 获取角色关系图
- ✅ 创建/获取/更新/删除关系时序事件

**API端点**:
```
POST   /api/v1/writer/projects/:id/characters                               # 创建角色
GET    /api/v1/writer/projects/:id/characters                               # 获取角色列表
GET    /api/v1/writer/projects/:id/characters/relations                     # 获取关系列表
GET    /api/v1/writer/projects/:id/characters/graph                         # 获取关系图
GET    /api/v1/writer/characters/:characterId                               # 获取角色详情
PUT    /api/v1/writer/characters/:characterId                               # 更新角色
DELETE /api/v1/writer/characters/:characterId                               # 删除角色
POST   /api/v1/writer/characters/relations                                  # 创建关系
DELETE /api/v1/writer/characters/relations/:relationId                      # 删除关系
POST   /api/v1/writer/characters/relations/:relationId/timeline             # 创建时序事件
GET    /api/v1/writer/characters/relations/:relationId/timeline             # 获取时序历史
PUT    /api/v1/writer/characters/relations/timeline-events/:eventId         # 更新时序事件
DELETE /api/v1/writer/characters/relations/timeline-events/:eventId         # 删除时序事件
```

**依赖服务**: `interfaces.CharacterService`

---

### 9. LocationApi (`location_api.go`)

**职责**: 故事地点管理，包括地点 CRUD、地点层级树、地点关系

**核心功能**:
- ✅ 创建/获取/更新/删除地点
- ✅ 获取项目地点列表
- ✅ 获取地点层级树
- ✅ 创建/删除地点关系
- ✅ 获取地点关系列表

**API端点**:
```
POST   /api/v1/writer/projects/:id/locations                    # 创建地点
GET    /api/v1/writer/projects/:id/locations                    # 获取地点列表
GET    /api/v1/writer/projects/:id/locations/tree               # 获取地点树
GET    /api/v1/writer/projects/:id/locations/relations          # 获取关系列表
GET    /api/v1/writer/locations/:locationId                     # 获取地点详情
PUT    /api/v1/writer/locations/:locationId                     # 更新地点
DELETE /api/v1/writer/locations/:locationId                     # 删除地点
POST   /api/v1/writer/locations/relations                       # 创建关系
DELETE /api/v1/writer/locations/relations/:relationId           # 删除关系
```

**依赖服务**: `interfaces.LocationService`

---

### 10. TimelineApi (`timeline_api.go`)

**职责**: 故事时间线管理，包括时间线 CRUD、事件管理和可视化

**核心功能**:
- ✅ 创建/获取/删除时间线
- ✅ 获取项目时间线列表
- ✅ 创建/获取/更新/删除时间线事件
- ✅ 获取时间线事件列表
- ✅ 获取时间线可视化数据

**API端点**:
```
POST   /api/v1/writer/projects/:id/timelines                          # 创建时间线
GET    /api/v1/writer/projects/:id/timelines                          # 获取时间线列表
GET    /api/v1/writer/timelines/:timelineId                           # 获取时间线详情
DELETE /api/v1/writer/timelines/:timelineId                           # 删除时间线
POST   /api/v1/writer/timelines/:timelineId/events                    # 创建事件
GET    /api/v1/writer/timelines/:timelineId/events                    # 获取事件列表
GET    /api/v1/writer/timelines/:timelineId/visualization             # 获取可视化数据
GET    /api/v1/writer/timeline-events/:eventId                       # 获取事件详情
PUT    /api/v1/writer/timeline-events/:eventId                       # 更新事件
DELETE /api/v1/writer/timeline-events/:eventId                       # 删除事件
```

**依赖服务**: `interfaces.TimelineService`

---

### 11. OutlineApi (`outline_api.go`)

**职责**: 大纲管理，支持树形大纲结构、大纲-文档双向同步

**核心功能**:
- ✅ 创建大纲节点
- ✅ 获取大纲详情
- ✅ 获取项目大纲列表
- ✅ 获取大纲树
- ✅ 更新/删除大纲
- ✅ 获取子节点列表

**API端点**:
```
POST   /api/v1/writer/projects/:id/outlines               # 创建大纲节点
GET    /api/v1/writer/projects/:id/outlines               # 获取大纲列表
GET    /api/v1/writer/projects/:id/outlines/tree          # 获取大纲树
GET    /api/v1/writer/projects/:id/outlines/children      # 获取子节点
GET    /api/v1/writer/outlines/:outlineId                 # 获取大纲详情
PUT    /api/v1/writer/outlines/:outlineId                 # 更新大纲
DELETE /api/v1/writer/outlines/:outlineId                 # 删除大纲
```

**依赖服务**: `interfaces.OutlineService`

---

### 12. EntityApi (`entity_api.go`)

**职责**: 统一实体管理，聚合角色、地点、物品、组织、伏笔等实体类型，提供实体图谱

**核心功能**:
- ✅ 获取项目实体列表（支持按类型筛选）
- ✅ 获取项目实体图谱（实体关系图）
- ✅ 更新实体状态字段

**API端点**:
```
GET    /api/v1/writer/projects/:id/entities                   # 获取实体列表
GET    /api/v1/writer/projects/:id/entities/graph             # 获取实体图谱
PUT    /api/v1/writer/entities/:entityId/state-fields         # 更新实体状态字段
```

**依赖服务**: `interfaces.EntityService`

---

### 13. KeywordApi (`keyword_api.go`)

**职责**: 关键词检索，支持角色和地点的前缀补全搜索

**核心功能**:
- ✅ 按项目搜索角色/地点关键词（前缀补全）

**API端点**:
```
GET    /api/v1/writer/projects/:id/keywords/search            # 搜索关键词
```

**依赖服务**: `interfaces.CharacterService`, `interfaces.LocationService`

---

### 14. EncyclopediaApi (`encyclopedia_api.go`)

**职责**: 设定百科，管理项目内的概念词条（世界观设定等）

**核心功能**:
- ✅ 获取项目概念列表
- ✅ 搜索概念
- ✅ 获取概念详情

**API端点**:
```
GET    /api/v1/writer/projects/:id/concepts                      # 获取概念列表
GET    /api/v1/writer/projects/:id/concepts/search               # 搜索概念
GET    /api/v1/writer/projects/:id/concepts/:conceptId           # 获取概念详情
```

**依赖服务**: `writerInterface.ConceptRepository`（当前为 nil，返回空列表）

---

### 15. LockApi (`lock_api.go`)

**职责**: 文档并发编辑控制，基于 Redis 实现分布式锁

**核心功能**:
- ✅ 锁定文档
- ✅ 解锁文档
- ✅ 刷新锁
- ✅ 获取锁状态
- ✅ 强制解锁
- ✅ 延长锁定时间

**API端点**:
```
POST   /api/v1/writer/documents/:id/lock              # 锁定文档
DELETE /api/v1/writer/documents/:id/lock               # 解锁文档
PUT    /api/v1/writer/documents/:id/lock/refresh       # 刷新锁
GET    /api/v1/writer/documents/:id/lock/status        # 获取锁状态
POST   /api/v1/writer/documents/:id/lock/force         # 强制解锁
POST   /api/v1/writer/documents/:id/lock/extend        # 延长锁定
```

**依赖服务**: `lock.DocumentLockService`（Redis 实现）

---

### 16. BatchOperationApi (`batch_operation_api.go`)

**职责**: 文档批量操作，支持批量删除、移动、复制，含 Preflight 预检查

**核心功能**:
- ✅ 提交批量操作（含预检查）
- ✅ 获取批量操作状态
- ✅ 取消批量操作
- ✅ 撤销批量操作

**API端点**:
```
POST   /api/v1/writer/batch-operations                   # 提交批量操作
GET    /api/v1/writer/batch-operations/:id               # 获取操作状态
POST   /api/v1/writer/batch-operations/:id/cancel        # 取消操作
POST   /api/v1/writer/batch-operations/:id/undo          # 撤销操作
```

**依赖服务**: `document.BatchOperationService`

---

### 17. SearchApi (`search_api.go`)

**职责**: 写作端文档全文搜索

**核心功能**:
- ✅ 按关键词搜索用户项目中的文档

**API端点**:
```
GET    /api/v1/writer/search/documents                  # 搜索文档
```

**依赖服务**: `searchservice.SearchService`

---

### 18. ExportApi (`export_api.go`)

**职责**: 文档/项目导出管理，支持 TXT、MD、DOCX 格式，异步任务模式

**核心功能**:
- ✅ 导出文档（异步任务）
- ✅ 导出项目（异步任务）
- ✅ 获取导出任务状态
- ✅ 下载导出文件
- ✅ 获取导出任务列表
- ✅ 删除导出任务
- ✅ 取消导出任务

**API端点**:
```
POST   /api/v1/writer/documents/:id/export              # 导出文档
POST   /api/v1/writer/projects/:id/export               # 导出项目（任务模式）
GET    /api/v1/writer/projects/:id/exports              # 获取导出任务列表
GET    /api/v1/writer/exports/:id                       # 获取任务状态
GET    /api/v1/writer/exports/:id/download              # 下载文件
DELETE /api/v1/writer/exports/:id                       # 删除任务
POST   /api/v1/writer/exports/:id/cancel                # 取消任务
```

**依赖服务**: `interfaces.ExportService`

---

### 19. ImportExportApi (`import_export_api.go`)

**职责**: 项目级别 ZIP 导入导出（直接下载/上传）

**核心功能**:
- ✅ 导出项目为 ZIP（直接下载）
- ✅ 从 ZIP 导入项目

**API端点**:
```
GET    /api/v1/writer/projects/:id/export               # 导出项目为 ZIP
POST   /api/v1/writer/projects/import                    # 从 ZIP 导入项目
```

**依赖服务**: `interfaces.ExportService`

---

### 20. StoryHarnessApi (`story_harness_api.go`)

**职责**: Story Harness 核心功能，提供章节上下文（Context Lens）和章节索引

**核心功能**:
- ✅ 获取章节上下文（角色快照、关系、作用域）
- ✅ 手动触发章节索引
- ✅ 重建章节投影

**API端点**:
```
GET    /api/v1/writer/projects/:id/chapters/:chapterId/context              # 获取章节上下文
POST   /api/v1/writer/projects/:id/chapters/:chapterId/trigger-index        # 触发章节索引
POST   /api/v1/writer/projects/:id/chapters/:chapterId/rebuild-projection   # 重建投影
```

**依赖服务**: `storyharness.ContextService`, `storyharness.IndexerService`, `storyharness.ChangeRequestService`

---

### 21. ChangeRequestApi (`change_request_api.go`)

**职责**: 章节变更建议管理，配合 Story Harness 的索引结果

**核心功能**:
- ✅ 获取章节建议列表（支持状态筛选）
- ✅ 获取建议详情
- ✅ 处理建议（接受/忽略/延后）

**API端点**:
```
GET    /api/v1/writer/projects/:id/chapters/:chapterId/change-requests    # 获取建议列表
GET    /api/v1/writer/change-requests/:requestId                          # 获取建议详情
PUT    /api/v1/writer/change-requests/:requestId/status                   # 处理建议
```

**依赖服务**: `storyharness.ChangeRequestService`

---

### 22. DashboardApi (`dashboard_api.go`)

**职责**: 作者工作台仪表板统计

**核心功能**:
- ✅ 获取仪表板统计数据（总字数、项目数、今日字数、待审核数、连续写作天数）

**API端点**:
```
GET    /api/v1/writer/dashboard/stats                    # 获取仪表板统计
```

**依赖服务**: `writer.DashboardService`

---

### 23. StatsApi (`stats_api.go`)

**职责**: 阅读/书店统计数据，包括作品统计、章节统计、收入、热力图等

**核心功能**:
- ✅ 获取作者作品列表
- ✅ 获取作品统计数据
- ✅ 获取章节统计数据
- ✅ 获取阅读热力图
- ✅ 获取收入统计
- ✅ 获取热门章节
- ✅ 获取每日统计
- ✅ 获取跳出点分析
- ✅ 获取留存率
- ✅ 获取订阅趋势
- ✅ 获取读者活跃度
- ✅ 作品对比
- ✅ 记录读者行为

**API端点**:
```
GET    /api/v1/writer/books                                    # 获取作品列表
GET    /api/v1/writer/books/:book_id/stats                     # 获取作品统计
GET    /api/v1/writer/books/:book_id/stats/subscribers         # 订阅趋势
GET    /api/v1/writer/books/:book_id/stats/chapters            # 章节统计
GET    /api/v1/writer/books/:book_id/stats/reader-activity     # 读者活跃度
GET    /api/v1/writer/books/:book_id/heatmap                   # 阅读热力图
GET    /api/v1/writer/books/:book_id/revenue                   # 收入统计
GET    /api/v1/writer/books/:book_id/top-chapters              # 热门章节
GET    /api/v1/writer/books/:book_id/daily-stats               # 每日统计
GET    /api/v1/writer/books/:book_id/drop-off-points           # 跳出点分析
GET    /api/v1/writer/books/:book_id/retention                 # 留存率
POST   /api/v1/writer/stats/compare                            # 作品对比
GET    /api/v1/writer/chapters/:chapter_id/stats               # 章节统计
POST   /api/v1/writer/reader/behavior                          # 记录读者行为
```

**依赖服务**: `readingStats.ReadingStatsService`, `bookstoreRepo.BookRepository`

---

### 24. WriterStatsAggregateApi (`writer_stats_aggregate_api.go`)

**职责**: 面向作者工作台的聚合统计接口，提供概览、趋势、今日数据

**核心功能**:
- ✅ 获取作者统计概览（总阅读量、总收入、总订阅等）
- ✅ 获取阅读量趋势
- ✅ 获取订阅趋势
- ✅ 获取章节统计
- ✅ 获取今日统计

**API端点**:
```
GET    /api/v1/writer/stats/overview                          # 统计概览
GET    /api/v1/writer/stats/views                             # 阅读量趋势
GET    /api/v1/writer/stats/subscribers                       # 订阅趋势
GET    /api/v1/writer/stats/chapters                          # 章节统计
GET    /api/v1/writer/stats/today                             # 今日统计
```

**依赖服务**: `readingStats.ReadingStatsService`, `bookstoreRepo.BookRepository`

---

### 25. AuditApi (`audit_api.go`)

**职责**: 内容审核，包括实时检测、全文审核、申诉、复核和违规记录管理

**核心功能**:
- ✅ 实时检测内容（不创建审核记录）
- ✅ 全文审核文档
- ✅ 获取审核结果
- ✅ 提交申诉
- ✅ 管理员待复核列表
- ✅ 复核审核结果
- ✅ 复核申诉
- ✅ 获取用户违规记录
- ✅ 获取用户违规统计
- ✅ 获取高风险审核记录

**API端点**:
```
POST   /api/v1/audit/check                                    # 实时检测
POST   /api/v1/documents/:id/audit                            # 全文审核
GET    /api/v1/documents/:id/audit-result                     # 获取审核结果
POST   /api/v1/audit/:id/appeal                               # 提交申诉
GET    /api/v1/users/:userId/violations                       # 获取违规记录
GET    /api/v1/users/:userId/violation-summary                # 获取违规统计
GET    /api/v1/admin/audit/pending                            # 待复核列表（管理员）
GET    /api/v1/admin/audit/high-risk                          # 高风险记录（管理员）
POST   /api/v1/admin/audit/:id/review                         # 复核审核（管理员）
POST   /api/v1/admin/audit/:id/appeal/review                  # 复核申诉（管理员）
```

**依赖服务**: `audit.ContentAuditService`

---

## 依赖关系说明

1. **路由注册入口**: 所有路由通过 `router/writer/` 目录下的路由文件统一注册到 `/api/v1/writer` 路径组下，并强制使用 JWT 认证中间件。

2. **服务依赖链**: API Handler -> Service 层 -> Repository 层 -> MongoDB/Redis。服务实例由 `RegisterWriterRoutes` 中统一创建并注入。

3. **关键基础设施依赖**:
   - **MongoDB**: 核心数据存储，用于项目、文档、版本、角色、地点等全部业务数据
   - **Redis**: 文档锁定服务（分布式锁）、分布式锁服务（大纲并发保护）
   - **EventBus**: 跨模块事件通信（文档创建/标题变更触发大纲同步、内容保存触发章节索引）

4. **双向同步机制**: 大纲（Outline）与文档（Document）之间存在双向同步：创建文档时自动创建大纲节点，修改文档标题时同步更新大纲标题。该机制由 `OutlineDocumentSyncService` 实现。

5. **Story Harness 集成**: 文档内容保存时自动触发章节索引（`IndexerService`），生成变更建议（`ChangeRequest`），支持章节上下文查询和投影重建。

6. **审核模块独立性**: 审核路由（`AuditApi`）不挂载在 `/writer` 路径组下，而是独立挂载于 `/api/v1/audit` 和 `/api/v1/admin/audit`，管理接口额外要求管理员角色权限。

7. **统计数据源**: `StatsApi` 和 `WriterStatsAggregateApi` 依赖阅读统计服务（`ReadingStatsService`）和书城仓库（`BookRepository`），数据来源于书城模块而非 writer 模块自身。

---

**版本**: v2.0
**更新日期**: 2026-05-21
**维护者**: Writer模块开发组
