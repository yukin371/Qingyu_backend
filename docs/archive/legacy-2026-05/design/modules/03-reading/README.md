# 03 - 阅读器模块

> **模块编号**: 03
> **模块名称**: Reader
> **负责功能**: 阅读进度管理、书架管理、笔记标注、阅读设置
> **完成度**: 🟢 75%

## 📋 目录结构

```
阅读器模块/
├── api/v1/reader/                    # 阅读器API
│   ├── books_api.go                # 书架管理
│   ├── progress_api.go             # 阅读进度
│   ├── annotations_api.go          # 笔记标注
│   ├── settings_api.go             # 阅读设置
│   └── history_api.go              # 阅读历史
├── service/reader/                  # 阅读器服务层
│   ├── bookshelf_service.go        # 书架服务
│   ├── progress_service.go         # 进度服务
│   ├── annotation_service.go       # 标注服务
│   └── settings_service.go         # 设置服务
├── repository/interfaces/reader/    # 仓储接口
├── repository/mongodb/reader/       # MongoDB仓储实现
│   ├── bookshelf_repository_mongo.go
│   ├── progress_repository_mongo.go
│   └── annotation_repository_mongo.go
└── models/reader/                   # 数据模型
    ├── bookshelf.go                # 书架
    ├── progress.go                 # 阅读进度
    ├── annotation.go               # 标注
    ├── bookmark.go                 # 书签
    └── settings.go                 # 阅读设置
```

## 🎯 核心功能

### 1. 书架管理

- **添加书籍**: 添加书籍到个人书架
- **书架分类**: 最近阅读、收藏、未读完、已读完
- **书架排序**: 按时间、进度、评分排序
- **书架筛选**: 按标签、状态筛选
- **移除书籍**: 从书架移除

### 2. 阅读进度

- **位置追踪**: 记录阅读位置（章节、字符位置）
- **进度计算**: 章节进度、全书进度
- **阅读时长**: 记录每次阅读时长
- **阅读历史**: 历史阅读记录
- **阅读统计**: 总阅读量、完读率、连续阅读天数

### 3. 标注系统

- **高亮标记**: 标记重要文本
- **笔记批注**: 添加个人想法
- **书签管理**: 保存阅读位置
- **标注分类**: 按颜色、标签分类
- **标注导出**: 导出标注和笔记

### 4. 阅读设置

- **字体设置**: 字号、字体、行高
- **主题切换**: 日间/夜间/护眼模式
- **翻页模式**: 仿真/滑动/滚动翻页
- **自动滚动**: 自动翻页设置
- **全屏模式**: 沉浸式阅读

## 📊 数据模型

### BookshelfEntry (书架条目)

```go
type BookshelfEntry struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    BookID          primitive.ObjectID   `bson:"book_id" json:"bookId"`
    ChapterID       *primitive.ObjectID  `bson:"chapter_id,omitempty" json:"chapterId,omitempty"`

    // 阅读状态
    ReadingStatus   ReadingStatus       `bson:"reading_status" json:"readingStatus"`
    Progress        float64              `bson:"progress" json:"progress"`
    LastChapter     int                  `bson:"last_chapter" json:"lastChapter"`
    LastPosition    int                  `bson:"last_position" json:"lastPosition"`

    // 个人标记
    IsFavorite      bool                 `bson:"is_favorite" json:"isFavorite"`
    Rating          int                  `bson:"rating" json:"rating"`
    Tags            []string             `bson:"tags" json:"tags"`
    Notes           string               `bson:"notes" json:"notes"`

    // 统计信息
    TotalReadTime   int64                `bson:"total_read_time" json:"totalReadTime"`
    ReadCount       int                  `bson:"read_count" json:"readCount"`

    // 时间戳
    AddedAt         time.Time            `bson:"added_at" json:"addedAt"`
    LastReadAt      *time.Time           `bson:"last_read_at,omitempty" json:"lastReadAt,omitempty"`
    FinishedAt      *time.Time           `bson:"finished_at,omitempty" json:"finishedAt,omitempty"`
}

type ReadingStatus string
const (
    ReadingStatusNotStarted   ReadingStatus = "not_started"
    ReadingStatusReading      ReadingStatus = "reading"
    ReadingStatusPaused       ReadingStatus = "paused"
    ReadingStatusCompleted    ReadingStatus = "completed"
    ReadingStatusAbandoned    ReadingStatus = "abandoned"
)
```

### ReadingProgress (阅读进度)

```go
type ReadingProgress struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    BookID          primitive.ObjectID   `bson:"book_id" json:"bookId"`
    ChapterID       primitive.ObjectID   `bson:"chapter_id" json:"chapterId"`

    // 进度信息
    ChapterPosition  int                  `bson:"chapter_position" json:"chapterPosition"`
    ChapterProgress float64              `bson:"chapter_progress" json:"chapterProgress"`
    BookProgress    float64              `bson:"book_progress" json:"bookProgress"`

    // 阅读统计
    ReadTime        int64                `bson:"read_time" json:"readTime"`
    TotalReadTime   int64                `bson:"total_read_time" json:"totalReadTime"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}
```

### Annotation (标注)

```go
type Annotation struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    BookID          primitive.ObjectID   `bson:"book_id" json:"bookId"`
    ChapterID       primitive.ObjectID   `bson:"chapter_id" json:"chapterId"`

    // 标注内容
    Type            AnnotationType       `bson:"type" json:"type"`
    Content         string               `bson:"content" json:"content"`
    Note            string               `bson:"note" json:"note"`
    Color           string               `bson:"color" json:"color"`

    // 位置信息
    StartPosition   int                  `bson:"start_position" json:"startPosition"`
    EndPosition     int                  `bson:"end_position" json:"endPosition"`
    StartOffset     int                  `bson:"start_offset" json:"startOffset"`
    EndOffset       int                  `bson:"end_offset" json:"endOffset"`
    ChapterNumber   int                  `bson:"chapter_number" json:"chapterNumber"`

    // 元数据
    Tags            []string             `bson:"tags" json:"tags"`
    IsPublic        bool                 `bson:"is_public" json:"isPublic"`
    IsShared        bool                 `bson:"is_shared" json:"isShared"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}

type AnnotationType string
const (
    AnnotationTypeHighlight  AnnotationType = "highlight"
    AnnotationTypeNote       AnnotationType = "note"
    AnnotationTypeBookmark   AnnotationType = "bookmark"
)
```

### ReaderSettings (阅读设置)

```go
type ReaderSettings struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`

    // 显示设置
    FontSize        int                  `bson:"font_size" json:"fontSize"`
    FontFamily      string               `bson:"font_family" json:"fontFamily"`
    LineHeight      float64              `bson:"line_height" json:"lineHeight"`
    ParagraphSpacing int                 `bson:"paragraph_spacing" json:"paragraphSpacing"`

    // 主题设置
    Theme           ReaderTheme          `bson:"theme" json:"theme"`
    BackgroundColor string               `bson:"background_color" json:"backgroundColor"`
    TextColor       string               `bson:"text_color" json:"textColor"`

    // 阅读模式
    Mode            ReaderMode           `bson:"mode" json:"mode"`
    AutoScroll      bool                 `bson:"auto_scroll" json:"autoScroll"`
    ScrollSpeed     int                  `bson:"scroll_speed" json:"scrollSpeed"`

    // 翻页设置
    FlipMode        FlipMode             `bson:"flip_mode" json:"flipMode"`
    TapToTurn       bool                 `bson:"tap_to_turn" json:"tapToTurn"`

    // 其他设置
    ShowComment     bool                 `bson:"show_comment" json:"showComment"`
    ShowNavigation  bool                 `bson:"show_navigation" json:"showNavigation"`
    FullScreenMode  bool                 `bson:"full_screen_mode" json:"fullScreenMode"`

    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}

type ReaderTheme string
const (
    ThemeLight      ReaderTheme = "light"
    ThemeDark       ReaderTheme = "dark"
    ThemeSepia      ReaderTheme = "sepia"
    ThemeCustom     ReaderTheme = "custom"
)

type ReaderMode string
const (
    ModeDay         ReaderMode = "day"
    ModeNight       ReaderMode = "night"
    ModeEyeCare     ReaderMode = "eye_care"
)

type FlipMode string
const (
    FlipSimulation  FlipMode = "simulation"
    FlipSlide       FlipMode = "slide"
    FlipScroll      FlipMode = "scroll"
    FlipNone        FlipMode = "none"
)
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/reader/books | 获取书架 | 是 |
| GET | /api/v1/reader/books/recent | 最近阅读 | 是 |
| GET | /api/v1/reader/books/unfinished | 未读完 | 是 |
| GET | /api/v1/reader/books/finished | 已读完 | 是 |
| POST | /api/v1/reader/books/:bookId | 添加到书架 | 是 |
| DELETE | /api/v1/reader/books/:bookId | 从书架移除 | 是 |
| GET | /api/v1/reader/progress/:bookId | 获取阅读进度 | 是 |
| POST | /api/v1/reader/progress | 保存阅读进度 | 是 |
| POST | /api/v1/reader/progress/time | 更新阅读时长 | 是 |
| GET | /api/v1/reader/progress/stats | 阅读统计 | 是 |
| POST | /api/v1/reader/annotations | 创建标注 | 是 |
| PUT | /api/v1/reader/annotations/:id | 更新标注 | 是 |
| DELETE | /api/v1/reader/annotations/:id | 删除标注 | 是 |
| GET | /api/v1/reader/annotations/notes | 获取笔记 | 是 |
| GET | /api/v1/reader/annotations/bookmarks | 获取书签 | 是 |
| GET | /api/v1/reader/annotations/highlights | 获取高亮 | 是 |
| GET | /api/v1/reader/settings | 获取阅读设置 | 是 |
| POST | /api/v1/reader/settings | 保存阅读设置 | 是 |
| PUT | /api/v1/reader/settings | 更新阅读设置 | 是 |

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 用户身份验证
- **06 - 书城**: 获取书籍和章节信息
- **02 - 写作创作**: 获取未发布内容（作者本人）

### 被依赖的模块
- **04 - 社交互动**: 分享阅读记录和标注
- **09 - AI**: 基于阅读历史推荐内容

## 🚀 性能优化

1. **进度同步优化**
   - 防抖处理，避免频繁保存
   - 批量更新进度

2. **标注查询优化**
   - 按章节索引标注
   - 标注分页加载

3. **缓存策略**
   - Redis 缓存阅读进度
   - 本地存储阅读设置

## 📈 扩展点

1. **朗读功能**
   - TTS文本转语音
   - 听书进度同步

2. **社交阅读**
   - 阅读排行
   - 阅读时长挑战
   - 阅读成就系统

3. **智能推荐**
   - 基于阅读历史推荐书籍
   - 基于标注内容推荐相似段落

4. **跨设备同步**
   - 阅读进度云端同步
   - 标注和设置同步

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/reader/`
**相关设计**: [阅读端模块设计文档](../../reading/), [阅读端](../../阅读端/)
