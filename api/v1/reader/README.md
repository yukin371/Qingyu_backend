# Reader API 模块 - 阅读器

## 📁 模块职责

**Reader（阅读器）**模块负责所有与阅读体验相关的功能，提供完整的沉浸式阅读环境。

## 🎯 核心功能

### 1. 书架管理
- 个人书架
- 添加/移除书籍
- 书架分类
- 最近阅读

### 2. 章节阅读
- 章节完整内容
- 上一章/下一章
- 章节导航
- 阅读位置记忆

### 3. 阅读进度
- 自动保存进度
- 章节进度显示
- 全书进度显示
- 阅读时长统计
- 阅读历史

### 4. 标注管理
- 文本标注（高亮、下划线）
- 笔记记录
- 标注搜索
- 标注导出
- 按书籍/章节筛选

### 5. 阅读设置
- 字体大小
- 行间距
- 背景主题
- 翻页方式
- 屏幕亮度
- 夜间模式

### 6. 阅读统计
- 每日阅读时长
- 每周阅读统计
- 阅读习惯分析
- 阅读成就

## 📦 文件结构

```
api/v1/reader/
├── books_api.go               # 书架管理
├── chapters_api.go            # 章节内容
├── progress.go                # 阅读进度
├── annotations_api.go         # 标注管理
├── annotations_api_optimized.go # 标注优化版本
├── setting_api.go             # 阅读设置
└── README.md                  # 本文档
```

## 🌐 API路由总览

### 所有接口都需要JWT认证

#### 书架管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/reader/books | 获取书架 | BooksAPI.GetBookshelf |
| POST | /api/v1/reader/books/:bookId | 添加到书架 | BooksAPI.AddToBookshelf |
| DELETE | /api/v1/reader/books/:bookId | 从书架移除 | BooksAPI.RemoveFromBookshelf |

#### 章节内容

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/reader/chapters/:id | 获取章节信息 | ChaptersAPI.GetChapter |
| GET | /api/v1/reader/chapters/:id/content | 获取章节内容 | ChaptersAPI.GetContent |
| GET | /api/v1/reader/chapters/:id/next | 获取下一章 | ChaptersAPI.GetNextChapter |
| GET | /api/v1/reader/chapters/:id/prev | 获取上一章 | ChaptersAPI.GetPrevChapter |

#### 阅读进度

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/reader/progress/:bookId | 获取阅读进度 | ProgressAPI.GetProgress |
| POST | /api/v1/reader/progress | 保存阅读进度 | ProgressAPI.SaveProgress |
| POST | /api/v1/reader/progress/time | 更新阅读时长 | ProgressAPI.UpdateReadingTime |
| GET | /api/v1/reader/progress/history | 获取阅读历史 | ProgressAPI.GetHistory |
| GET | /api/v1/reader/progress/statistics | 获取阅读统计 | ProgressAPI.GetStatistics |
| GET | /api/v1/reader/progress/statistics/daily | 获取每日统计 | ProgressAPI.GetDailyStats |
| GET | /api/v1/reader/progress/statistics/weekly | 获取每周统计 | ProgressAPI.GetWeeklyStats |

#### 标注管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/reader/annotations | 获取标注列表 | AnnotationsAPI.GetAnnotations |
| GET | /api/v1/reader/annotations/:id | 获取标注详情 | AnnotationsAPI.GetAnnotation |
| POST | /api/v1/reader/annotations | 创建标注 | AnnotationsAPI.CreateAnnotation |
| PUT | /api/v1/reader/annotations/:id | 更新标注 | AnnotationsAPI.UpdateAnnotation |
| DELETE | /api/v1/reader/annotations/:id | 删除标注 | AnnotationsAPI.DeleteAnnotation |
| DELETE | /api/v1/reader/annotations | 批量删除标注 | AnnotationsAPI.BatchDelete |
| GET | /api/v1/reader/annotations/book/:bookId | 获取书籍标注 | AnnotationsAPI.GetBookAnnotations |
| GET | /api/v1/reader/annotations/chapter/:chapterId | 获取章节标注 | AnnotationsAPI.GetChapterAnnotations |
| GET | /api/v1/reader/annotations/search | 搜索标注 | AnnotationsAPI.SearchAnnotations |
| GET | /api/v1/reader/annotations/export | 导出标注 | AnnotationsAPI.ExportAnnotations |

#### 阅读设置

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/reader/settings | 获取阅读设置 | SettingAPI.GetSettings |
| PUT | /api/v1/reader/settings | 更新阅读设置 | SettingAPI.UpdateSettings |
| POST | /api/v1/reader/settings/reset | 重置设置 | SettingAPI.ResetSettings |

## 🔄 与Bookstore模块的区别

| 功能 | Bookstore（书店） | Reader（阅读器） |
|------|------------------|-----------------|
| **定位** | 发现和浏览 | 阅读和学习 |
| **用户场景** | 找书、选书 | 读书、记笔记 |
| **核心功能** | 搜索、推荐、详情 | 阅读、进度、标注 |
| **章节** | 预览（前几章） | 完整内容 |
| **认证要求** | 多为公开 | 必须认证 |
| **数据存储** | 书籍元数据 | 用户阅读数据 |
| **数据隔离** | 全局共享 | 用户私有 |

## 🎨 使用场景

### 场景1：开始阅读
```
1. 从书架获取书籍 → GET /reader/books
2. 选择一本书
3. 获取上次阅读位置 → GET /reader/progress/:bookId
4. 获取章节内容 → GET /reader/chapters/:id/content
5. 开始阅读，自动保存进度
```

### 场景2：做笔记
```
1. 阅读过程中遇到重要内容
2. 选中文本
3. 创建标注 → POST /reader/annotations
4. 添加笔记内容
5. 稍后可以搜索和导出标注
```

### 场景3：查看阅读统计
```
1. 访问阅读统计 → GET /reader/progress/statistics
2. 查看每日阅读时长
3. 查看每周统计图表
4. 查看阅读成就
```

### 场景4：个性化设置
```
1. 进入阅读设置 → GET /reader/settings
2. 调整字体大小、行间距
3. 选择主题（白天/夜间）
4. 保存设置 → PUT /reader/settings
5. 设置实时生效
```

## 🔧 技术特点

### 1. 离线支持
- 章节内容缓存
- 离线阅读
- 进度同步

### 2. 实时同步
- 阅读进度实时保存
- 多设备同步
- 断点续读

### 3. 性能优化
- 章节预加载
- 图片懒加载
- 虚拟滚动

### 4. 用户体验
- 流畅翻页动画
- 手势操作
- 护眼模式
- 沉浸式阅读

## 📊 数据模型

### ReadingProgress（阅读进度）
```go
type ReadingProgress struct {
    UserID        string
    BookID        string
    ChapterID     string
    Progress      float64  // 0-1
    LastReadAt    time.Time
    ReadingTime   int64    // 秒
}
```

### Annotation（标注）
```go
type Annotation struct {
    ID            string
    UserID        string
    BookID        string
    ChapterID     string
    SelectedText  string
    Note          string
    Type          string   // highlight, underline, note
    Color         string
    CreatedAt     time.Time
}
```

### ReadingSettings（阅读设置）
```go
type ReadingSettings struct {
    UserID        string
    FontSize      int
    LineHeight    float64
    Theme         string
    PageMode      string
    Brightness    int
    AutoSave      bool
}
```

### ReadingStatistics（阅读统计）
```go
type ReadingStatistics struct {
    UserID           string
    TotalReadingTime int64
    BooksRead        int
    DailyAverage     int64
    WeeklyData       []DailyStats
}
```

## 🚀 后续规划

### Phase 1（已完成）
- ✅ 书架管理
- ✅ 章节阅读
- ✅ 阅读进度
- ✅ 标注管理
- ✅ 阅读设置

### Phase 2（进行中）
- 🔄 阅读统计优化
- 🔄 离线阅读
- 🔄 多设备同步

### Phase 3（计划中）
- 📋 朗读功能（TTS）
- 📋 翻译功能
- 📋 AI辅助理解
- 📋 社交分享
- 📋 阅读成就系统

## 💡 最佳实践

### 1. 进度保存策略
- 每15秒自动保存一次
- 切换章节时保存
- 退出阅读时保存
- 后台运行时保存

### 2. 标注管理
- 支持多种标注类型
- 标注颜色分类
- 支持全文搜索
- 支持导出为Markdown

### 3. 性能优化
- 章节内容分页加载
- 预加载下一章
- 标注懒加载
- 图片压缩

### 4. 用户体验
- 记住上次阅读位置
- 平滑滚动和翻页
- 手势控制
- 快捷键支持

## 📚 相关文档

- [Bookstore API 模块](../bookstore/README.md)
- [Reader Service 设计](../../../doc/design/reader/README.md)
- [阅读器UI设计](../../../doc/design/reader/ui_design.md)
- [数据库设计](../../../doc/database/reader_schema.md)

---

**版本**: v2.0  
**更新日期**: 2025-10-22  
**维护者**: Reader模块开发组

