# 阅读器Repository层实施文档

> **模块**: 阅读器系统 - Repository层  
> **阶段**: 阶段二  
> **完成时间**: 2025-10-08  
> **实施人员**: 青羽后端团队

---

## 📋 目录

1. [实施概述](#1-实施概述)
2. [ChapterRepository实现](#2-chapterrepository实现)
3. [ReadingProgressRepository实现](#3-readingprogressrepository实现)
4. [AnnotationRepository实现](#4-annotationrepository实现)
5. [数据库索引设计](#5-数据库索引设计)
6. [实施检查清单](#6-实施检查清单)

---

## 1. 实施概述

### 1.1 目标

完成阅读器系统Repository层的实现，包括：
- ✅ 章节管理Repository
- ✅ 阅读进度Repository
- ✅ 标注（笔记/书签）Repository

### 1.2 架构设计

遵循项目统一的Repository模式：

```
repository/
├── interfaces/reading/              # Repository接口定义
│   ├── chapter_repository.go        # 章节Repository接口
│   ├── reading_progress_repository.go  # 阅读进度Repository接口
│   └── annotation_repository.go     # 标注Repository接口
└── mongodb/reading/                 # MongoDB实现
    ├── chapter_repository_mongo.go
    ├── reading_progress_repository_mongo.go
    └── annotation_repository_mongo.go
```

### 1.3 技术栈

- **语言**: Go 1.21+
- **数据库**: MongoDB
- **驱动**: mongo-driver (go.mongodb.org/mongo-driver)
- **架构模式**: Repository Pattern

---

## 2. ChapterRepository实现

### 2.1 接口设计

**文件路径**: `repository/interfaces/reading/chapter_repository.go`

#### 核心功能

| 功能分类 | 方法名称 | 说明 |
|---------|---------|------|
| 基础CRUD | Create, GetByID, Update, Delete | 基本的增删改查 |
| 章节查询 | GetByBookID, GetByChapterNum | 按书籍/章节号查询 |
| 章节导航 | GetPrevChapter, GetNextChapter | 上一章/下一章 |
| 状态查询 | GetPublishedChapters, GetVIPChapters | 按状态筛选 |
| 统计功能 | CountByBookID, CountVIPChapters | 章节统计 |
| 批量操作 | BatchCreate, BatchUpdateStatus | 批量处理 |
| VIP管理 | CheckVIPAccess, GetChapterPrice | VIP权限检查 |
| 内容管理 | GetChapterContent, UpdateChapterContent | 章节内容 |

#### 接口定义示例

```go
type ChapterRepository interface {
    // 基础CRUD操作
    Create(ctx context.Context, chapter *reader.Chapter) error
    GetByID(ctx context.Context, id string) (*reader.Chapter, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
    Delete(ctx context.Context, id string) error

    // 章节查询
    GetByBookID(ctx context.Context, bookID string) ([]*reader.Chapter, error)
    GetByChapterNum(ctx context.Context, bookID string, chapterNum int) (*reader.Chapter, error)
    
    // 章节导航
    GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error)
    GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error)
    
    // VIP权限
    CheckVIPAccess(ctx context.Context, chapterID string) (bool, error)
    GetChapterPrice(ctx context.Context, chapterID string) (int64, error)
    
    // 健康检查
    Health(ctx context.Context) error
}
```

### 2.2 MongoDB实现

**文件路径**: `repository/mongodb/reading/chapter_repository_mongo.go`

#### 核心实现要点

**1. 章节导航实现**

```go
// GetNextChapter 获取下一章
func (r *MongoChapterRepository) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
    var chapter reader.Chapter
    filter := bson.M{
        "book_id":     bookID,
        "chapter_num": bson.M{"$gt": currentChapterNum},
        "status":      1,
    }
    opts := options.FindOne().SetSort(bson.D{{Key: "chapter_num", Value: 1}})
    
    err := r.collection.FindOne(ctx, filter, opts).Decode(&chapter)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // 已经是最后一章
        }
        return nil, fmt.Errorf("查询下一章失败: %w", err)
    }
    
    return &chapter, nil
}
```

**2. 分页查询**

```go
func (r *MongoChapterRepository) GetByBookIDWithPagination(ctx context.Context, bookID string, limit, offset int64) ([]*reader.Chapter, error) {
    filter := bson.M{
        "book_id": bookID,
        "status":  1,
    }
    
    opts := options.Find().
        SetSkip(offset).
        SetLimit(limit).
        SetSort(bson.D{{Key: "chapter_num", Value: 1}})
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, fmt.Errorf("查询章节列表失败: %w", err)
    }
    defer cursor.Close(ctx)
    
    var chapters []*reader.Chapter
    if err = cursor.All(ctx, &chapters); err != nil {
        return nil, fmt.Errorf("解析章节数据失败: %w", err)
    }
    
    return chapters, nil
}
```

**3. 批量操作**

```go
func (r *MongoChapterRepository) BatchCreate(ctx context.Context, chapters []*reader.Chapter) error {
    if len(chapters) == 0 {
        return nil
    }
    
    docs := make([]interface{}, len(chapters))
    now := time.Now()
    for i, chapter := range chapters {
        if chapter.ID == "" {
            chapter.ID = generateID()
        }
        chapter.CreatedAt = now
        chapter.UpdatedAt = now
        docs[i] = chapter
    }
    
    _, err := r.collection.InsertMany(ctx, docs)
    if err != nil {
        return fmt.Errorf("批量创建章节失败: %w", err)
    }
    
    return nil
}
```

#### 性能优化

1. **查询优化**
   - 章节导航使用索引排序
   - 内容查询使用Projection限制字段
   - 批量操作使用InsertMany/UpdateMany

2. **索引设计**（见第5章）

---

## 3. ReadingProgressRepository实现

### 3.1 接口设计

**文件路径**: `repository/interfaces/reading/reading_progress_repository.go`

#### 核心功能

| 功能分类 | 方法名称 | 说明 |
|---------|---------|------|
| 基础CRUD | Create, GetByID, Update, Delete | 基本增删改查 |
| 进度查询 | GetByUserAndBook, GetRecentReadingByUser | 查询阅读进度 |
| 进度保存 | SaveProgress, UpdateReadingTime | 保存和更新 |
| 统计查询 | GetTotalReadingTime, GetReadingTimeByPeriod | 时长统计 |
| 阅读记录 | GetReadingHistory, GetUnfinishedBooks | 阅读历史 |
| 数据同步 | SyncProgress, GetProgressesByUser | 进度同步 |

#### 接口定义示例

```go
type ReadingProgressRepository interface {
    // 基础CRUD
    Create(ctx context.Context, progress *reader.ReadingProgress) error
    GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error)
    
    // 进度查询
    GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error)
    GetRecentReadingByUser(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error)
    
    // 进度保存
    SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error
    UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error
    
    // 统计查询
    GetTotalReadingTime(ctx context.Context, userID string) (int64, error)
    GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error)
    
    // 数据同步
    SyncProgress(ctx context.Context, userID string, progresses []*reader.ReadingProgress) error
}
```

### 3.2 MongoDB实现

**文件路径**: `repository/mongodb/reading/reading_progress_repository_mongo.go`

#### 核心实现要点

**1. Upsert操作实现进度保存**

```go
func (r *MongoReadingProgressRepository) SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
    filter := bson.M{
        "user_id": userID,
        "book_id": bookID,
    }
    
    update := bson.M{
        "$set": bson.M{
            "chapter_id":   chapterID,
            "progress":     progress,
            "last_read_at": time.Now(),
            "updated_at":   time.Now(),
        },
        "$setOnInsert": bson.M{
            "_id":          generateProgressID(),
            "reading_time": int64(0),
            "created_at":   time.Now(),
        },
    }
    
    opts := options.Update().SetUpsert(true)
    _, err := r.collection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return fmt.Errorf("保存阅读进度失败: %w", err)
    }
    
    return nil
}
```

**2. 增量更新阅读时长**

```go
func (r *MongoReadingProgressRepository) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
    filter := bson.M{
        "user_id": userID,
        "book_id": bookID,
    }
    
    update := bson.M{
        "$inc": bson.M{
            "reading_time": duration, // 增量更新
        },
        "$set": bson.M{
            "last_read_at": time.Now(),
            "updated_at":   time.Now(),
        },
    }
    
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("更新阅读时长失败: %w", err)
    }
    
    // 如果没有记录，创建新记录
    if result.MatchedCount == 0 {
        progress := &reader.ReadingProgress{
            ID:          generateProgressID(),
            UserID:      userID,
            BookID:      bookID,
            ReadingTime: duration,
            Progress:    0,
            LastReadAt:  time.Now(),
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        return r.Create(ctx, progress)
    }
    
    return nil
}
```

**3. 聚合查询统计阅读时长**

```go
func (r *MongoReadingProgressRepository) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
    pipeline := mongo.Pipeline{
        {{Key: "$match", Value: bson.M{"user_id": userID}}},
        {{Key: "$group", Value: bson.M{
            "_id":   nil,
            "total": bson.M{"$sum": "$reading_time"},
        }}},
    }
    
    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return 0, fmt.Errorf("统计总阅读时长失败: %w", err)
    }
    defer cursor.Close(ctx)
    
    var result []struct {
        Total int64 `bson:"total"`
    }
    if err = cursor.All(ctx, &result); err != nil {
        return 0, fmt.Errorf("解析统计结果失败: %w", err)
    }
    
    if len(result) == 0 {
        return 0, nil
    }
    
    return result[0].Total, nil
}
```

**4. 批量同步**

```go
func (r *MongoReadingProgressRepository) BatchUpdateProgress(ctx context.Context, progresses []*reader.ReadingProgress) error {
    if len(progresses) == 0 {
        return nil
    }
    
    models := make([]mongo.WriteModel, len(progresses))
    for i, progress := range progresses {
        filter := bson.M{
            "user_id": progress.UserID,
            "book_id": progress.BookID,
        }
        
        update := bson.M{
            "$set": bson.M{
                "chapter_id":   progress.ChapterID,
                "progress":     progress.Progress,
                "reading_time": progress.ReadingTime,
                "last_read_at": progress.LastReadAt,
                "updated_at":   time.Now(),
            },
            "$setOnInsert": bson.M{
                "_id":        progress.ID,
                "created_at": time.Now(),
            },
        }
        
        models[i] = mongo.NewUpdateOneModel().
            SetFilter(filter).
            SetUpdate(update).
            SetUpsert(true)
    }
    
    opts := options.BulkWrite().SetOrdered(false)
    _, err := r.collection.BulkWrite(ctx, models, opts)
    if err != nil {
        return fmt.Errorf("批量更新阅读进度失败: %w", err)
    }
    
    return nil
}
```

#### 业务特点

1. **频繁更新**: 使用Upsert避免重复插入
2. **增量统计**: 使用`$inc`操作符累加阅读时长
3. **数据同步**: 支持批量Upsert
4. **时间筛选**: 支持按时间段统计

---

## 4. AnnotationRepository实现

### 4.1 接口设计

**文件路径**: `repository/interfaces/reading/annotation_repository.go`

#### 核心功能

| 功能分类 | 方法名称 | 说明 |
|---------|---------|------|
| 基础CRUD | Create, GetByID, Update, Delete | 基本增删改查 |
| 笔记操作 | GetNotes, GetNotesByChapter, SearchNotes | 笔记管理 |
| 书签操作 | GetBookmarks, GetBookmarkByPosition | 书签管理 |
| 高亮操作 | GetHighlights, GetHighlightsByChapter | 高亮管理 |
| 统计功能 | CountByUser, CountByBook, CountByType | 标注统计 |
| 批量操作 | BatchCreate, BatchDelete | 批量处理 |
| 数据同步 | SyncAnnotations, GetRecentAnnotations | 数据同步 |
| 分享功能 | GetPublicAnnotations, GetSharedAnnotations | 公开分享 |

#### 接口定义示例

```go
type AnnotationRepository interface {
    // 基础CRUD
    Create(ctx context.Context, annotation *reader.Annotation) error
    GetByID(ctx context.Context, id string) (*reader.Annotation, error)
    
    // 笔记操作
    GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
    SearchNotes(ctx context.Context, userID string, keyword string) ([]*reader.Annotation, error)
    
    // 书签操作
    GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
    GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, startOffset int) (*reader.Annotation, error)
    
    // 高亮操作
    GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
    
    // 数据同步
    SyncAnnotations(ctx context.Context, userID string, annotations []*reader.Annotation) error
    GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error)
    
    // 分享功能
    GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error)
}
```

### 4.2 MongoDB实现

**文件路径**: `repository/mongodb/reading/annotation_repository_mongo.go`

#### 核心实现要点

**1. 按类型查询标注**

```go
func (r *MongoAnnotationRepository) GetByType(ctx context.Context, userID, bookID string, annotationType int) ([]*reader.Annotation, error) {
    filter := bson.M{
        "user_id": userID,
        "book_id": bookID,
        "type":    annotationType,
    }
    opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, fmt.Errorf("查询标注失败: %w", err)
    }
    defer cursor.Close(ctx)
    
    var annotations []*reader.Annotation
    if err = cursor.All(ctx, &annotations); err != nil {
        return nil, fmt.Errorf("解析标注数据失败: %w", err)
    }
    
    return annotations, nil
}

// 笔记: type = 1
func (r *MongoAnnotationRepository) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
    return r.GetByType(ctx, userID, bookID, 1)
}

// 书签: type = 2
func (r *MongoAnnotationRepository) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
    return r.GetByType(ctx, userID, bookID, 2)
}

// 高亮: type = 3
func (r *MongoAnnotationRepository) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
    return r.GetByType(ctx, userID, bookID, 3)
}
```

**2. 全文搜索笔记**

```go
func (r *MongoAnnotationRepository) SearchNotes(ctx context.Context, userID string, keyword string) ([]*reader.Annotation, error) {
    filter := bson.M{
        "user_id": userID,
        "type":    1,
        "$or": []bson.M{
            {"content": bson.M{"$regex": keyword, "$options": "i"}},
            {"note": bson.M{"$regex": keyword, "$options": "i"}},
        },
    }
    opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, fmt.Errorf("搜索笔记失败: %w", err)
    }
    defer cursor.Close(ctx)
    
    var annotations []*reader.Annotation
    if err = cursor.All(ctx, &annotations); err != nil {
        return nil, fmt.Errorf("解析笔记数据失败: %w", err)
    }
    
    return annotations, nil
}
```

**3. 按位置查询书签**

```go
func (r *MongoAnnotationRepository) GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, startOffset int) (*reader.Annotation, error) {
    var annotation reader.Annotation
    filter := bson.M{
        "user_id":      userID,
        "book_id":      bookID,
        "chapter_id":   chapterID,
        "type":         2,
        "start_offset": startOffset,
    }
    
    err := r.collection.FindOne(ctx, filter).Decode(&annotation)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // 书签不存在
        }
        return nil, fmt.Errorf("查询书签失败: %w", err)
    }
    
    return &annotation, nil
}
```

**4. 批量同步**

```go
func (r *MongoAnnotationRepository) SyncAnnotations(ctx context.Context, userID string, annotations []*reader.Annotation) error {
    if len(annotations) == 0 {
        return nil
    }
    
    models := make([]mongo.WriteModel, len(annotations))
    for i, annotation := range annotations {
        if annotation.ID == "" {
            annotation.ID = generateAnnotationID()
        }
        
        filter := bson.M{"_id": annotation.ID}
        update := bson.M{
            "$set": bson.M{
                "user_id":      annotation.UserID,
                "book_id":      annotation.BookID,
                "chapter_id":   annotation.ChapterID,
                "type":         annotation.Type,
                "content":      annotation.Content,
                "note":         annotation.Note,
                "color":        annotation.Color,
                "start_offset": annotation.StartOffset,
                "end_offset":   annotation.EndOffset,
                "is_public":    annotation.IsPublic,
                "updated_at":   time.Now(),
            },
            "$setOnInsert": bson.M{
                "created_at": annotation.CreatedAt,
            },
        }
        
        models[i] = mongo.NewUpdateOneModel().
            SetFilter(filter).
            SetUpdate(update).
            SetUpsert(true)
    }
    
    opts := options.BulkWrite().SetOrdered(false)
    _, err := r.collection.BulkWrite(ctx, models, opts)
    if err != nil {
        return fmt.Errorf("同步标注失败: %w", err)
    }
    
    return nil
}
```

**5. 公开标注查询**

```go
func (r *MongoAnnotationRepository) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
    filter := bson.M{
        "book_id":    bookID,
        "chapter_id": chapterID,
        "is_public":  true,
    }
    opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, fmt.Errorf("查询公开标注失败: %w", err)
    }
    defer cursor.Close(ctx)
    
    var annotations []*reader.Annotation
    if err = cursor.All(ctx, &annotations); err != nil {
        return nil, fmt.Errorf("解析标注数据失败: %w", err)
    }
    
    return annotations, nil
}
```

#### 业务特点

1. **多类型管理**: 笔记、书签、高亮统一存储，按type区分
2. **全文搜索**: 支持笔记内容搜索
3. **位置定位**: 书签按start_offset精确定位
4. **公开分享**: 支持标注公开和查询
5. **批量同步**: 支持多端数据同步

---

## 5. 数据库索引设计

### 5.1 chapters 集合

```javascript
// 章节集合索引
db.chapters.createIndex({ "book_id": 1, "chapter_num": 1 }, { unique: true })
db.chapters.createIndex({ "book_id": 1, "status": 1, "chapter_num": 1 })
db.chapters.createIndex({ "book_id": 1, "is_vip": 1 })
db.chapters.createIndex({ "publish_time": 1 })
```

**索引说明**:
- `book_id + chapter_num`: 唯一索引，保证同一本书章节号不重复
- `book_id + status + chapter_num`: 查询已发布章节的复合索引
- `book_id + is_vip`: VIP章节筛选
- `publish_time`: 按发布时间排序

### 5.2 reading_progress 集合

```javascript
// 阅读进度集合索引
db.reading_progress.createIndex({ "user_id": 1, "book_id": 1 }, { unique: true })
db.reading_progress.createIndex({ "user_id": 1, "last_read_at": -1 })
db.reading_progress.createIndex({ "user_id": 1, "progress": 1 })
db.reading_progress.createIndex({ "last_read_at": 1 }) // 用于清理旧数据
```

**索引说明**:
- `user_id + book_id`: 唯一索引，一个用户一本书只有一条进度记录
- `user_id + last_read_at`: 查询最近阅读记录
- `user_id + progress`: 查询未读完/已读完书籍
- `last_read_at`: 用于定期清理旧数据

### 5.3 annotations 集合

```javascript
// 标注集合索引
db.annotations.createIndex({ "user_id": 1, "book_id": 1, "chapter_id": 1 })
db.annotations.createIndex({ "user_id": 1, "book_id": 1, "type": 1 })
db.annotations.createIndex({ "user_id": 1, "created_at": -1 })
db.annotations.createIndex({ "book_id": 1, "chapter_id": 1, "is_public": 1 })
db.annotations.createIndex({ "user_id": 1, "book_id": 1, "chapter_id": 1, "type": 1, "start_offset": 1 })

// 全文索引（用于笔记搜索）
db.annotations.createIndex({ "content": "text", "note": "text" })
```

**索引说明**:
- `user_id + book_id + chapter_id`: 章节标注查询
- `user_id + book_id + type`: 按类型筛选标注
- `user_id + created_at`: 最近标注查询
- `book_id + chapter_id + is_public`: 公开标注查询
- `user_id + book_id + chapter_id + type + start_offset`: 书签位置定位
- `content + note`: 全文搜索索引

---

## 6. 实施检查清单

### 6.1 代码实现检查

- [x] ChapterRepository接口定义完整
- [x] ChapterRepository MongoDB实现
- [x] ReadingProgressRepository接口定义完整
- [x] ReadingProgressRepository MongoDB实现
- [x] AnnotationRepository接口定义完整
- [x] AnnotationRepository MongoDB实现
- [x] 所有方法都有错误处理
- [x] 所有方法都有Context支持
- [x] 时间戳自动更新(CreatedAt/UpdatedAt)
- [x] 支持健康检查(Health方法)

### 6.2 性能优化检查

- [x] 查询使用了合适的索引
- [x] 批量操作使用InsertMany/BulkWrite
- [x] 大字段查询使用Projection限制
- [x] 排序字段建立索引
- [x] 唯一约束使用unique索引

### 6.3 业务逻辑检查

- [x] 章节导航（上/下一章）正确实现
- [x] VIP权限检查逻辑完整
- [x] 阅读进度支持Upsert
- [x] 阅读时长增量更新
- [x] 标注按类型区分
- [x] 支持公开标注查询
- [x] 支持批量数据同步

### 6.4 数据一致性检查

- [x] 用户-书籍进度唯一约束
- [x] 书籍-章节号唯一约束
- [x] 时间戳字段自动维护
- [x] 删除操作级联考虑

### 6.5 测试检查

- [ ] 单元测试编写（下一阶段）
- [ ] 集成测试编写（下一阶段）
- [ ] 性能测试（下一阶段）
- [ ] 并发测试（下一阶段）

---

## 7. 下一步计划

### 7.1 Service层实现（阶段三）

实现阅读器业务逻辑层：
1. ChapterService - 章节获取、内容管理
2. ReadingProgressService - 进度保存、统计
3. AnnotationService - 标注管理、搜索

### 7.2 API层实现（阶段四）

实现阅读器HTTP接口：
1. ChapterAPI - 章节相关接口
2. ProgressAPI - 进度相关接口
3. AnnotationAPI - 标注相关接口
4. 路由配置

### 7.3 测试完善

1. 编写Repository层单元测试
2. 编写集成测试
3. 性能测试和优化

---

**文档维护**: 青羽后端团队  
**完成时间**: 2025-10-08  
**下一步**: 实现阅读器Service层

