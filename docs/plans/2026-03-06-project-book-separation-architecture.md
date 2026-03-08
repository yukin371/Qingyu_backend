# Project-Book 分离架构设计

**设计日期**: 2026-03-06
**设计者**: Kore
**优先级**: 🔴 P0
**问题来源**: BookStatus枚举统一讨论中发现架构缺陷

---

## 问题描述

### 当前架构问题

当前系统采用 **Project + Book 分离架构**，但**两者之间缺少关联关系**：

```
┌────────────────────┐                    ┌────────────────────┐
│     Project        │                    │        Book        │
│   (作者创作管理)     │                    │     (读者阅读)      │
├────────────────────┤                    ├────────────────────┤
│ ID: project_001    │                    │ ID: book_001       │
│ Title: "星河骑士"   │    ❌ 无关联         │ Title: "星河骑士"   │
│ AuthorID: user_123 │                    │ AuthorID: user_123 │
│ Status:            │                    │ Status: ongoing    │
│   serializing      │                    │                    │
│ Chapters: [1,2,3]  │                    │ Chapters: [1,2,3] │
└────────────────────┘                    └────────────────────┘
         │                                         ▲
         │  发布时？                                │
         └─────────────────────────────────────────┘
                    如何关联？如何同步？
```

**导致的问题**：
1. ❌ 发布时无法知道 Book 对应哪个 Project
2. ❌ Project 更新时无法同步到 Book
3. ❌ 同作者同名 Book 无法区分来源
4. ❌ 无法实现"从 Project 创建 Book"功能

---

## 架构设计

### 1. 添加关联字段

#### Book 模型修改

```go
// models/bookstore/book.go

type Book struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 基本信息
    Title         string `bson:"title" json:"title"`
    Author        string `bson:"author" json:"author"`
    AuthorID      string `bson:"author_id,omitempty" json:"authorId,omitempty"`

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  新增：关联字段
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ProjectID     *string `bson:"project_id,omitempty" json:"projectId,omitempty"` // ← 关联的Project ID
    SourceType    SourceType `bson:"source_type,omitempty" json:"sourceType,omitempty"` // ← 来源类型
    SyncMode      SyncMode `bson:"sync_mode,omitempty" json:"syncMode,omitempty"` // ← 同步模式
    LastSyncedAt  *time.Time `bson:"last_synced_at,omitempty" json:"lastSyncedAt,omitempty"` // ← 最后同步时间
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    // 阅读相关字段
    Introduction  string `bson:"introduction" json:"introduction"`
    Cover         string `bson:"cover" json:"cover"`
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
    Tags          []string `bson:"tags" json:"tags"`
    Status        BookStatus `bson:"status" json:"status"`
    Rating        types.Rating `bson:"rating" json:"rating"`
    ViewCount     int64 `bson:"view_count" json:"viewCount"`
    WordCount     int64 `bson:"word_count" json:"wordCount"`
    ChapterCount  int `bson:"chapter_count" json:"chapterCount"`
    Price         float64 `bson:"price" json:"price"`
    IsFree        bool `bson:"is_free" json:"isFree"`
    PublishedAt   *time.Time `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
    LastUpdateAt  *time.Time `bson:"last_update_at,omitempty" json:"lastUpdateAt,omitempty"`
}

// SourceType 书籍来源类型
type SourceType string

const (
    SourceTypeProject SourceType = "project" // 来自Writer项目
    SourceTypeManual  SourceType = "manual"  // 手动创建
    SourceTypeImport  SourceType = "import"  // 外部导入
)

// SyncMode 同步模式
type SyncMode string

const (
    SyncModeManual  SyncMode = "manual"  // 手动同步（默认）
    SyncModeAuto    SyncMode = "auto"    // 自动同步
    SyncModeOneTime SyncMode = "onetime" // 一次性同步（快照模式）
)
```

#### Chapter 模型修改

```go
// models/bookstore/chapter.go

type Chapter struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    BookID       primitive.ObjectID `bson:"book_id" json:"bookId"`

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  新增：关联字段
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ProjectID    *string `bson:"project_id,omitempty" json:"projectId,omitempty"` // ← 来源Project ID
    ProjectChapterID *string `bson:"project_chapter_id,omitempty" json:"projectChapterId,omitempty"` // ← 来源章节ID
    ContentHash  string `bson:"content_hash,omitempty" json:"contentHash,omitempty"` // ← 内容哈希（用于变更检测）
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    ChapterNumber int    `bson:"chapter_number" json:"chapterNumber"`
    Title         string `bson:"title" json:"title"`
    Content       string `bson:"content" json:"content"`
    WordCount     int    `bson:"word_count" json:"wordCount"`
    PublishedAt   *time.Time `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
    Price         float64 `bson:"price" json:"price"`
    IsFree        bool   `bson:"is_free" json:"isFree"`
    ViewCount     int64  `bson:"view_count" json:"viewCount"`
}
```

---

### 2. 发布流程设计

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         发布流程                                        │
└─────────────────────────────────────────────────────────────────────────┘

Author                        Service Layer                 Repository
   │                              │                              │
   │  1. 发布Project                │                              │
   │  POST /projects/:id/publish    │                              │
   │─────────────────────────────>│                              │
   │                              │                              │
   │                              │  2. 检查Project是否存在Book  │
   │                              │──────────────────────────>│
   │                              │                              │
   │                              │  3. 返回关联的Book（如有）  │
   │                              │<──────────────────────────│
   │                              │                              │
   │                              │  4. 根据SyncMode决定：      │
   │                              │     - onetime: 创建新Book   │
   │                              │     - auto: 更新现有Book    │
   │                              │     - manual: 不自动同步    │
   │                              │                              │
   │  5. 返回Book ID               │                              │
   │<─────────────────────────────│                              │
   │                              │                              │
```

#### 发布服务实现

```go
// service/writer/publish_service.go

type PublishService struct {
    projectRepo   interfaces.ProjectRepository
    bookRepo      interfaces.BookRepository
    chapterRepo   interfaces.ChapterRepository
    syncService   *BookSyncService
}

// PublishProject 发布项目为书籍
func (s *PublishService) PublishProject(
    ctx context.Context,
    projectID string,
    req *PublishRequest,
) (*Book, error) {
    // 1. 获取Project
    project, err := s.projectRepo.GetByID(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("项目不存在: %w", err)
    }

    // 2. 检查是否已存在关联的Book
    existingBook, _ := s.bookRepo.GetByProjectID(ctx, projectID)

    var book *Book

    if existingBook != nil && existingBook.SyncMode == SyncModeAuto {
        // 自动同步模式：更新现有Book
        book, err = s.syncService.SyncFromProject(ctx, project, existingBook)
    } else {
        // 手动模式或首次发布：创建新Book
        book, err = s.createBookFromProject(ctx, project, req)
    }

    if err != nil {
        return nil, err
    }

    // 3. 更新Project状态
    project.Status = writer.StatusSerializing
    project.PublishedAt = &time.Time{}
    *project.PublishedAt = time.Now()

    if err := s.projectRepo.Update(ctx, project); err != nil {
        return nil, err
    }

    return book, nil
}

// createBookFromProject 从Project创建Book
func (s *PublishService) createBookFromProject(
    ctx context.Context,
    project *writer.Project,
    req *PublishRequest,
) (*Book, error) {
    // 创建Book
    book := &Book{
        Title:        project.Title,
        Author:       project.AuthorID, // 需要从用户服务获取
        AuthorID:     project.AuthorID,
        Introduction: project.Summary,
        Cover:        project.CoverURL,
        Status:       BookStatusOngoing,
        ProjectID:    &project.ID.Hex(), // ← 关联Project
        SourceType:   SourceTypeProject,
        SyncMode:     req.SyncMode,       // ← 同步模式
        PublishedAt:  project.PublishedAt,
    }

    if err := s.bookRepo.Create(ctx, book); err != nil {
        return nil, err
    }

    // 同步章节
    for _, chapter := range project.Chapters {
        bookChapter := &Chapter{
            BookID:            book.ID,
            ProjectID:         &project.ID.Hex(),
            ProjectChapterID:  &chapter.ID.Hex(),
            ChapterNumber:     chapter.ChapterNumber,
            Title:             chapter.Title,
            Content:           chapter.Content,
            WordCount:         chapter.WordCount,
            ContentHash:       calculateHash(chapter.Content),
        }

        if err := s.chapterRepo.Create(ctx, bookChapter); err != nil {
            return nil, err
        }
    }

    return book, nil
}
```

---

### 3. 同步机制设计

#### 3.1 手动同步（默认）

```go
// service/writer/book_sync_service.go

type BookSyncService struct {
    projectRepo   interfaces.ProjectRepository
    bookRepo      interfaces.BookRepository
    chapterRepo   interfaces.ChapterRepository
}

// SyncFromProject 手动触发同步
func (s *BookSyncService) SyncFromProject(
    ctx context.Context,
    project *writer.Project,
    book *Book,
) (*Book, error) {
    // 1. 更新Book基本信息
    book.Title = project.Title
    book.Introduction = project.Summary
    book.Cover = project.CoverURL
    book.LastSyncedAt = &time.Time{}
    *book.LastSyncedAt = time.Now()

    if err := s.bookRepo.Update(ctx, book); err != nil {
        return nil, err
    }

    // 2. 同步章节（根据ContentHash检测变更）
    projectChapters, err := s.getProjectChapters(ctx, project.ID.Hex())
    if err != nil {
        return nil, err
    }

    for _, pc := range projectChapters {
        existingChapter, _ := s.chapterRepo.GetByProjectChapterID(
            ctx, pc.ID.Hex(),
        )

        newHash := calculateHash(pc.Content)

        if existingChapter == nil {
            // 新章节：创建
            s.createChapter(ctx, book.ID, pc)
        } else if existingChapter.ContentHash != newHash {
            // 内容变更：更新
            s.updateChapterContent(ctx, existingChapter, pc)
        }
    }

    return book, nil
}
```

#### 3.2 自动同步（可选）

```go
// 内部/messaging/sync_event_handler.go

type SyncEventHandler struct {
    syncService *BookSyncService
}

// HandleProjectUpdate 处理Project更新事件
func (h *SyncEventHandler) HandleProjectUpdate(
    ctx context.Context,
    event *ProjectUpdatedEvent,
) error {
    // 查找关联的Book
    book, err := h.bookRepo.GetByProjectID(ctx, event.ProjectID)
    if err != nil {
        return nil // 没有关联Book，跳过
    }

    // 检查同步模式
    if book.SyncMode != SyncModeAuto {
        return nil // 非自动同步模式，跳过
    }

    // 执行同步
    _, err = h.syncService.SyncFromProject(ctx, event.Project, book)
    return err
}
```

---

### 4. 数据库索引

```go
// migration/mongodb/book_indexes.go

func CreateBookIndexes(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("books")

    indexes := []mongo.IndexModel{
        // Project ID关联索引
        {
            Keys: bson.D{{Key: "project_id", Value: 1}},
            Options: options.Index().
                SetSparse(true). // 允许为空（手动创建的Book）
                SetName("idx_project_id"),
        },
        // 组合索引：作者 + Project ID
        {
            Keys: bson.D{
                {Key: "author_id", Value: 1},
                {Key: "project_id", Value: 1},
            },
            Options: options.Index().SetName("idx_author_project"),
        },
        // 同步模式索引
        {
            Keys: bson.D{{Key: "sync_mode", Value: 1}},
            Options: options.Index().SetName("idx_sync_mode"),
        },
    }

    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

---

### 5. API 设计

```go
// api/v1/writer/publish_api.go

// PublishProject 发布项目
// @Summary 发布项目为书籍
// @Description 将Writer项目发布到书城
// @Tags Writer
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body PublishRequest true "发布参数"
// @Success 200 {object} Book
// @Router /writer/projects/{id}/publish [post]
func (api *PublishAPI) PublishProject(c *gin.Context) {
    projectID := c.Param("id")

    var req PublishRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    book, err := api.publishService.PublishProject(c.Request.Context(), projectID, &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, book)
}

// SyncBook 同步书籍
// @Summary 手动同步书籍内容
// @Description 从Project同步最新内容到Book
// @Tags Writer
// @Accept json
// @Produce json
// @Param id path string true "书籍ID"
// @Success 200 {object} Book
// @Router /writer/books/{id}/sync [post]
func (api *PublishAPI) SyncBook(c *gin.Context) {
    bookID := c.Param("id")

    book, err := api.bookRepo.GetByID(c.Request.Context(), bookID)
    if err != nil {
        c.JSON(404, gin.H{"error": "书籍不存在"})
        return
    }

    if book.ProjectID == nil {
        c.JSON(400, gin.H{"error": "该书籍未关联Project，无法同步"})
        return
    }

    project, err := api.projectRepo.GetByID(c.Request.Context(), *book.ProjectID)
    if err != nil {
        c.JSON(404, gin.H{"error": "关联的Project不存在"})
        return
    }

    updatedBook, err := api.syncService.SyncFromProject(
        c.Request.Context(),
        project,
        book,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, updatedBook)
}

// UpdateSyncMode 更新同步模式
// @Summary 更新书籍同步模式
// @Tags Writer
// @Accept json
// @Produce json
// @Param id path string true "书籍ID"
// @Param request body UpdateSyncModeRequest true "同步模式"
// @Success 200 {object} Book
// @Router /writer/books/{id}/sync-mode [put]
func (api *PublishAPI) UpdateSyncMode(c *gin.Context) {
    bookID := c.Param("id")

    var req UpdateSyncModeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    book, err := api.bookRepo.GetByID(c.Request.Context(), bookID)
    if err != nil {
        c.JSON(404, gin.H{"error": "书籍不存在"})
        return
    }

    book.SyncMode = req.SyncMode
    if err := api.bookRepo.Update(c.Request.Context(), book); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, book)
}

type PublishRequest struct {
    SyncMode    SyncMode `json:"syncMode" validate:"required"`    // 同步模式
    BookTitle   *string   `json:"bookTitle,omitempty"`           // 自定义书名（可选）
    Price       *float64  `json:"price,omitempty"`               // 价格（可选）
}

type UpdateSyncModeRequest struct {
    SyncMode SyncMode `json:"syncMode" validate:"required"`
}
```

---

## 前端影响

### 需要更新的API调用

```typescript
// src/modules/writer/api/publish.ts

/**
 * 发布项目为书籍
 */
export async function publishProject(
  projectId: string,
  params: {
    syncMode: 'manual' | 'auto' | 'onetime'
    bookTitle?: string
    price?: number
  }
): Promise<Book> {
  return http.post(`/api/v1/writer/projects/${projectId}/publish`, params)
}

/**
 * 手动同步书籍
 */
export async function syncBook(bookId: string): Promise<Book> {
  return http.post(`/api/v1/writer/books/${bookId}/sync`)
}

/**
 * 更新同步模式
 */
export async function updateSyncMode(
  bookId: string,
  syncMode: 'manual' | 'auto' | 'onetime'
): Promise<Book> {
  return http.put(`/api/v1/writer/books/${bookId}/sync-mode`, { syncMode })
}
```

### UI更新建议

```
┌─────────────────────────────────────────────────────────────┐
│  项目详情页 - 发布按钮                                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  [发布到书城]                                                 │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  同步模式:                                           │     │
│  │  ○ 手动同步 - 发布后不再自动更新                      │     │
│  │  ○ 自动同步 - Project更新时自动同步到Book            │     │
│  │  ○ 一次性 - 创建快照后不再关联                        │     │
│  │                                                      │     │
│  │  书名: [星河骑士] (可覆盖)                            │     │
│  │  价格: [100] 分                                       │     │
│  │                                                      │     │
│  │           [取消]              [发布]                  │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 实施计划

### Phase 1: 数据模型（1天）

- [ ] Book 模型添加关联字段
- [ ] Chapter 模型添加关联字段
- [ ] 创建数据库迁移脚本

### Phase 2: 服务层（2天）

- [ ] 实现 PublishService
- [ ] 实现 BookSyncService
- [ ] 实现 Repository 方法（GetByProjectID等）

### Phase 3: API层（1天）

- [ ] 发布API
- [ ] 同步API
- [ ] 更新同步模式API

### Phase 4: 前端（1天）

- [ ] 发布对话框UI
- [ ] 同步按钮
- [ ] 同步模式切换

### Phase 5: 测试（1天）

- [ ] 单元测试
- [ ] 集成测试
- [ ] E2E测试

---

## 注意事项

### 数据一致性

1. **同步失败处理**：记录同步日志，允许手动重试
2. **并发控制**：使用分布式锁防止同步冲突
3. **内容哈希**：使用ContentHash避免不必要的更新

### 性能优化

1. **增量同步**：只同步变更的章节
2. **异步处理**：大项目同步使用异步任务
3. **缓存策略**：Book数据缓存，SyncMode=Manual时长期缓存

### 安全考虑

1. **权限验证**：只有Project所有者可以发布
2. **价格控制**：防止恶意设置高价
3. **内容审核**：发布前进行内容审核

---

## 相关文档

- [BookStatus 枚举统一设计](./2026-03-05-book-status-unification-design.md)
- [事务管理器设计](./2026-03-05-transaction-manager-design.md)
- [CategoryIDs 类型统一设计](./2026-03-05-category-ids-unification-design.md)

---

**设计完成时间**: 2026-03-06
**预计实施时间**: 6天
**建议执行者**: 后端团队
