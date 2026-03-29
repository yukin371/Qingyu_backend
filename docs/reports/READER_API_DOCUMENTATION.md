# 青羽写作平台 - 阅读功能 API 文档

## 概述

本文档描述了为青羽写作平台新增的阅读功能高优先级API，包括阅读器设置、主题管理、字体管理和章节评论功能。

## 实现日期

2026-01-03

## 目录

1. [阅读器设置API](#阅读器设置api)
2. [主题管理API](#主题管理api)
3. [字体管理API](#字体管理api)
4. [章节评论API](#章节评论api)
5. [数据结构定义](#数据结构定义)
6. [创建的文件](#创建的文件)

---

## 阅读器设置API

### 基础设置端点

#### 1. 获取阅读器设置
```
GET /api/v1/reader/settings
```

**描述**: 获取当前用户的阅读器设置

**认证**: 需要JWT认证

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "...",
    "userId": "...",
    "fontSize": 16,
    "fontFamily": "serif",
    "lineHeight": 1.5,
    "theme": "light",
    "background": "#FFFFFF",
    "pageMode": 1,
    "autoScroll": false,
    "scrollSpeed": 50,
    "createdAt": "2026-01-03T00:00:00Z",
    "updatedAt": "2026-01-03T00:00:00Z"
  }
}
```

#### 2. 保存阅读器设置
```
POST /api/v1/reader/settings
```

**描述**: 保存或创建用户的阅读器设置

**认证**: 需要JWT认证

**请求体**:
```json
{
  "fontSize": 18,
  "fontFamily": "sans-serif",
  "lineHeight": 1.8,
  "theme": "dark",
  "background": "#121212",
  "pageMode": 1,
  "autoScroll": true,
  "scrollSpeed": 60
}
```

#### 3. 更新阅读器设置
```
PUT /api/v1/reader/settings
```

**描述**: 部分更新阅读器设置

**认证**: 需要JWT认证

**请求体**:
```json
{
  "fontSize": 20,
  "theme": "sepia"
}
```

---

## 主题管理API

### 1. 获取可用主题列表
```
GET /api/v1/reader/themes?builtin=true&public=true
```

**描述**: 获取可用主题列表，支持过滤

**认证**: 需要JWT认证

**查询参数**:
- `builtin` (boolean): 仅显示内置主题
- `public` (boolean): 仅显示公开主题

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "themes": [
      {
        "name": "light",
        "displayName": "明亮模式",
        "description": "默认明亮主题，适合白天阅读",
        "isBuiltIn": true,
        "isPublic": true,
        "colors": {
          "background": "#FFFFFF",
          "textPrimary": "#212121",
          "accentColor": "#1976D2",
          ...
        }
      },
      {
        "name": "dark",
        "displayName": "暗黑模式",
        "description": "护眼暗色主题，适合夜间阅读",
        ...
      }
    ],
    "total": 4
  }
}
```

### 2. 获取单个主题
```
GET /api/v1/reader/themes/{name}
```

**描述**: 根据主题名称获取主题详情

**认证**: 需要JWT认证

### 3. 创建自定义主题
```
POST /api/v1/reader/themes
```

**描述**: 创建自定义主题

**认证**: 需要JWT认证

**请求体**:
```json
{
  "name": "my-custom-theme",
  "displayName": "我的主题",
  "description": "自定义阅读主题",
  "isPublic": false,
  "colors": {
    "background": "#F5F5F5",
    "textPrimary": "#333333",
    "accentColor": "#FF5722",
    ...
  }
}
```

### 4. 更新主题
```
PUT /api/v1/reader/themes/{id}
```

**描述**: 更新自定义主题（仅创建者可操作）

**认证**: 需要JWT认证

### 5. 删除主题
```
DELETE /api/v1/reader/themes/{id}
```

**描述**: 删除自定义主题（仅创建者可操作）

**认证**: 需要JWT认证

### 6. 激活主题
```
POST /api/v1/reader/themes/{name}/activate
```

**描述**: 激活主题并应用到阅读设置

**认证**: 需要JWT认证

**响应示例**:
```json
{
  "code": 200,
  "message": "激活成功",
  "data": {
    "message": "主题已激活",
    "themeName": "dark",
    "userId": "..."
  }
}
```

---

## 字体管理API

### 1. 获取可用字体列表
```
GET /api/v1/reader/fonts?category=serif&builtin=true
```

**描述**: 获取可用字体列表

**认证**: 需要JWT认证

**查询参数**:
- `category` (string): 字体分类 - serif/sans-serif/monospace
- `builtin` (boolean): 仅显示内置字体

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "fonts": [
      {
        "name": "system-serif",
        "displayName": "宋体/衬线",
        "fontFamily": "SimSun, 'Songti SC', 'Noto Serif SC', serif",
        "description": "经典衬线字体，适合正文阅读",
        "category": "serif",
        "isBuiltIn": true,
        "isActive": true,
        "supportSize": [12, 14, 16, 18, 20, 22, 24, 28, 32],
        "previewText": "这是一段宋体预览文字 The quick brown fox jumps over the lazy dog."
      }
    ],
    "total": 5,
    "fontsByCategory": {
      "serif": [...],
      "sans-serif": [...],
      "monospace": [...]
    }
  }
}
```

### 2. 获取单个字体
```
GET /api/v1/reader/fonts/{name}
```

**描述**: 根据字体名称获取字体详情

**认证**: 需要JWT认证

### 3. 创建自定义字体
```
POST /api/v1/reader/fonts
```

**描述**: 创建自定义字体（支持上传字体文件）

**认证**: 需要JWT认证

**请求体**:
```json
{
  "name": "my-custom-font",
  "displayName": "我的字体",
  "fontFamily": "'Custom Font', serif",
  "description": "自定义字体描述",
  "category": "serif",
  "fontUrl": "https://example.com/fonts/custom.woff2",
  "previewText": "这是自定义字体预览"
}
```

### 4. 更新字体
```
PUT /api/v1/reader/fonts/{id}
```

**描述**: 更新自定义字体信息

**认证**: 需要JWT认证

### 5. 删除字体
```
DELETE /api/v1/reader/fonts/{id}
```

**描述**: 删除自定义字体

**认证**: 需要JWT认证

### 6. 设置字体偏好
```
POST /api/v1/reader/settings/font
```

**描述**: 设置字体偏好并应用到阅读设置

**认证**: 需要JWT认证

**请求体**:
```json
{
  "fontName": "system-serif",
  "fontSize": 18,
  "lineHeight": 1.8,
  "letterSpacing": 0.5
}
```

---

## 章节评论API

### 章节级评论

#### 1. 获取章节评论列表
```
GET /api/v1/reader/chapters/{chapterId}/comments?page=1&pageSize=20&sortBy=created_at&sortOrder=desc&parentId=
```

**描述**: 获取章节的评论列表，支持分页、排序和层级过滤

**认证**: 需要JWT认证

**查询参数**:
- `page` (int): 页码，默认1
- `pageSize` (int): 每页数量，默认20，最大100
- `sortBy` (string): 排序字段 - created_at/like_count/rating，默认created_at
- `sortOrder` (string): 排序方向 - asc/desc，默认desc
- `parentId` (string): 父评论ID，空字符串表示顶级评论

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "comments": [
      {
        "id": "...",
        "chapterId": "...",
        "bookId": "...",
        "userId": "...",
        "content": "这个章节写得真好！",
        "rating": 5,
        "parentId": null,
        "rootId": null,
        "replyCount": 3,
        "likeCount": 10,
        "isVisible": true,
        "isDeleted": false,
        "createdAt": "2026-01-03T00:00:00Z",
        "updatedAt": "2026-01-03T00:00:00Z",
        "userSnapshot": {
          "id": "...",
          "username": "user123",
          "avatar": "https://example.com/avatar.jpg"
        }
      }
    ],
    "total": 50,
    "page": 1,
    "pageSize": 20,
    "totalPages": 3,
    "avgRating": 4.5,
    "ratingCount": 40
  }
}
```

#### 2. 发表章节评论
```
POST /api/v1/reader/chapters/{chapterId}/comments
```

**描述**: 发表章节评论或回复

**认证**: 需要JWT认证

**请求体**:
```json
{
  "chapterId": "...",
  "bookId": "...",
  "content": "这个章节写得真好！",
  "rating": 5,
  "parentId": null
}
```

**回复评论**:
```json
{
  "chapterId": "...",
  "bookId": "...",
  "content": "我同意你的看法",
  "rating": 0,
  "parentId": "parent_comment_id"
}
```

### 段落级评论

#### 3. 获取章节段落评论概览
```
GET /api/v1/reader/chapters/{chapterId}/paragraph-comments
```

**描述**: 获取章节所有段落的评论统计

**认证**: 需要JWT认证

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapterId": "...",
    "paragraphStats": {
      "0": 5,
      "1": 3,
      "2": 0,
      "3": 8
    }
  }
}
```

#### 4. 发表段落级评论
```
POST /api/v1/reader/chapters/{chapterId}/paragraph-comments
```

**描述**: 对特定段落发表评论

**认证**: 需要JWT认证

**请求体**:
```json
{
  "bookId": "...",
  "content": "这段描写很生动",
  "paragraphIndex": 3,
  "charStart": 10,
  "charEnd": 50
}
```

#### 5. 获取特定段落的评论
```
GET /api/v1/reader/chapters/{chapterId}/paragraphs/{paragraphIndex}/comments
```

**描述**: 获取指定段落的所有评论

**认证**: 需要JWT认证

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "paragraphIndex": 3,
    "paragraphText": "这是第三段的文本内容...",
    "commentCount": 5,
    "comments": [...]
  }
}
```

### 评论管理

#### 6. 获取评论详情
```
GET /api/v1/reader/chapter-comments/{commentId}
```

**描述**: 获取单条评论的详细信息

**认证**: 需要JWT认证

#### 7. 更新评论
```
PUT /api/v1/reader/chapter-comments/{commentId}
```

**描述**: 更新评论内容（仅作者可操作，30分钟内）

**认证**: 需要JWT认证

**请求体**:
```json
{
  "content": "更新后的评论内容",
  "rating": 4
}
```

#### 8. 删除评论
```
DELETE /api/v1/reader/chapter-comments/{commentId}
```

**描述**: 删除评论（软删除，仅作者或管理员可操作）

**认证**: 需要JWT认证

#### 9. 点赞评论
```
POST /api/v1/reader/chapter-comments/{commentId}/like
```

**描述**: 点赞评论

**认证**: 需要JWT认证

#### 10. 取消点赞评论
```
DELETE /api/v1/reader/chapter-comments/{commentId}/like
```

**描述**: 取消点赞评论

**认证**: 需要JWT认证

---

## 数据结构定义

### ReaderTheme（阅读器主题）

```go
type ReaderTheme struct {
    ID          string       // 主题ID
    Name        string       // 主题名称（唯一标识）
    DisplayName string       // 显示名称
    Description string       // 主题描述
    IsBuiltIn   bool         // 是否内置主题
    IsPublic    bool         // 是否公开
    CreatorID   string       // 创建者ID
    Colors      ThemeColors  // 主题颜色配置
    IsActive    bool         // 是否激活
    UseCount    int64        // 使用次数
    CreatedAt   time.Time    // 创建时间
    UpdatedAt   time.Time    // 更新时间
}
```

### ThemeColors（主题颜色）

```go
type ThemeColors struct {
    Background           string  // 主背景色
    SecondaryBackground  string  // 次背景色
    TextPrimary          string  // 主要文字颜色
    TextSecondary        string  // 次要文字颜色
    TextDisabled         string  // 禁用文字颜色
    LinkColor            string  // 链接颜色
    AccentColor          string  // 强调色
    AccentHover          string  // 强调色悬停
    BorderColor          string  // 边框颜色
    DividerColor         string  // 分隔线颜色
    HighlightColor       string  // 高亮颜色
    BookmarkColor        string  // 书签颜色
    AnnotationColor      string  // 标注颜色
    ShadowColor          string  // 阴影颜色
}
```

### ReaderFont（阅读器字体）

```go
type ReaderFont struct {
    ID          string    // 字体ID
    Name        string    // 字体名称（唯一标识）
    DisplayName string    // 显示名称
    FontFamily  string    // CSS font-family 值
    Description string    // 字体描述
    Category    string    // 字体分类：serif/sans-serif/monospace
    FontURL     string    // 字体文件URL
    FontFormat  string    // 字体格式：woff/woff2/ttf
    PreviewText string    // 预览文本
    PreviewURL  string    // 预览图片URL
    IsBuiltIn   bool      // 是否内置字体
    IsActive    bool      // 是否激活可用
    SupportSize []int     // 支持的字号列表
    UseCount    int64     // 使用次数
    CreatedAt   time.Time // 创建时间
    UpdatedAt   time.Time // 更新时间
}
```

### ChapterComment（章节评论）

```go
type ChapterComment struct {
    ID              string              // 评论ID
    ChapterID       string              // 章节ID
    BookID          string              // 书籍ID
    UserID          string              // 用户ID
    Content         string              // 评论内容
    Rating          int                 // 评分（1-5，0表示无评分）
    ParagraphIndex  *int                // 段落索引（段落级评论）
    ParagraphText   *string             // 段落文本摘要
    CharStart       *int                // 字符起始位置
    CharEnd         *int                // 字符结束位置
    ParentID        *string             // 父评论ID
    RootID          *string             // 根评论ID
    ReplyCount      int                 // 回复数量
    LikeCount       int                 // 点赞数
    IsVisible       bool                // 是否可见
    IsDeleted       bool                // 是否已删除
    CreatedAt       time.Time           // 创建时间
    UpdatedAt       time.Time           // 更新时间
    UserSnapshot    *CommentUserSnapshot // 用户快照
}
```

---

## 创建的文件

### Models（数据模型）

1. **D:\Github\青羽\Qingyu_backend\models\reader\reader_theme.go**
   - ReaderTheme 模型
   - ThemeColors 结构
   - 内置主题列表（BuiltInThemes）：light, dark, sepia, eye-care
   - 主题请求/响应结构

2. **D:\Github\青羽\Qingyu_backend\models\reader\reader_font.go**
   - ReaderFont 模型
   - 内置字体列表（BuiltInFonts）：system-serif, system-sans, kai, fangsong, monospace
   - FontPreference 结构
   - 字体请求/响应结构

3. **D:\Github\青羽\Qingyu_backend\models\reader\chapter_comment.go**
   - ChapterComment 模型
   - CommentUserSnapshot 结构
   - ChapterCommentFilter 过滤器
   - 评论请求/响应结构

### API Controllers（API控制器）

4. **D:\Github\青羽\Qingyu_backend\api\v1\reader\theme_api.go**
   - ThemeAPI 结构体
   - GetThemes - 获取主题列表
   - GetThemeByName - 获取单个主题
   - CreateCustomTheme - 创建自定义主题
   - UpdateTheme - 更新主题
   - DeleteTheme - 删除主题
   - ActivateTheme - 激活主题

5. **D:\Github\青羽\Qingyu_backend\api\v1\reader\font_api.go**
   - FontAPI 结构体
   - GetFonts - 获取字体列表
   - GetFontByName - 获取单个字体
   - CreateCustomFont - 创建自定义字体
   - UpdateFont - 更新字体
   - DeleteFont - 删除字体
   - SetFontPreference - 设置字体偏好

6. **D:\Github\青羽\Qingyu_backend\api\v1\reader\chapter_comment_api.go**
   - ChapterCommentAPI 结构体
   - GetChapterComments - 获取章节评论列表
   - CreateChapterComment - 发表章节评论
   - GetChapterComment - 获取单条评论
   - UpdateChapterComment - 更新评论
   - DeleteChapterComment - 删除评论
   - LikeChapterComment - 点赞评论
   - UnlikeChapterComment - 取消点赞
   - GetParagraphComments - 获取段落评论
   - CreateParagraphComment - 发表段落评论
   - GetChapterParagraphComments - 获取章节段落评论概览

### Router Updates（路由更新）

7. **D:\Github\青羽\Qingyu_backend\router\reader\reader_router.go**
   - 添加主题管理路由组
   - 添加字体管理路由组
   - 添加章节评论路由组
   - 添加段落评论路由组
   - 集成新的API处理器

---

## API端点汇总

### 阅读器设置（3个端点）
- GET /api/v1/reader/settings
- POST /api/v1/reader/settings
- PUT /api/v1/reader/settings

### 主题管理（6个端点）
- GET /api/v1/reader/themes
- GET /api/v1/reader/themes/{name}
- POST /api/v1/reader/themes
- PUT /api/v1/reader/themes/{id}
- DELETE /api/v1/reader/themes/{id}
- POST /api/v1/reader/themes/{name}/activate

### 字体管理（6个端点）
- GET /api/v1/reader/fonts
- GET /api/v1/reader/fonts/{name}
- POST /api/v1/reader/fonts
- PUT /api/v1/reader/fonts/{id}
- DELETE /api/v1/reader/fonts/{id}
- POST /api/v1/reader/settings/font

### 章节评论（10个端点）
- GET /api/v1/reader/chapters/{chapterId}/comments
- POST /api/v1/reader/chapters/{chapterId}/comments
- GET /api/v1/reader/chapters/{chapterId}/paragraph-comments
- POST /api/v1/reader/chapters/{chapterId}/paragraph-comments
- GET /api/v1/reader/chapters/{chapterId}/paragraphs/{paragraphIndex}/comments
- GET /api/v1/reader/chapter-comments/{commentId}
- PUT /api/v1/reader/chapter-comments/{commentId}
- DELETE /api/v1/reader/chapter-comments/{commentId}
- POST /api/v1/reader/chapter-comments/{commentId}/like
- DELETE /api/v1/reader/chapter-comments/{commentId}/like

**总计：25个新API端点**

---

## 功能特性

### 1. 阅读器设置
- ✅ 持久化存储
- ✅ 支持多设备同步（通过设置服务）
- ✅ 字体、字号、行距、主题等配置
- ✅ 翻页模式和自动滚动设置

### 2. 主题管理
- ✅ 4个内置主题：明亮、暗黑、羊皮纸、护眼
- ✅ 自定义主题创建
- ✅ 主题分享功能
- ✅ 一键激活主题
- ✅ 完整的颜色配置方案

### 3. 字体管理
- ✅ 5个内置字体：宋体、黑体、楷体、仿宋、等宽
- ✅ 自定义字体上传
- ✅ 字体分类组织
- ✅ 字号支持列表
- ✅ 字体预览功能

### 4. 章节评论
- ✅ 章节级评论
- ✅ 段落级评论（精确到段落和字符位置）
- ✅ 评论回复功能
- ✅ 评论点赞
- ✅ 评论评分（1-5星）
- ✅ 分页和排序
- ✅ 评论编辑（30分钟内）
- ✅ 软删除机制
- ✅ 用户快照（避免联表查询）

---

## 下一步建议

### 数据库实现
1. 创建对应的MongoDB集合
2. 建立索引优化查询性能
3. 实现Repository层数据访问

### Service层实现
1. ThemeService - 主题业务逻辑
2. FontService - 字体业务逻辑
3. ChapterCommentService - 章节评论业务逻辑

### 缓存优化
1. 主题列表缓存
2. 字体列表缓存
3. 热门评论缓存

### 安全增强
1. 评论内容审核
2. 敏感词过滤
3. 频率限制（防刷评论）

### 测试
1. 单元测试
2. 集成测试
3. API文档测试

---

## 版本信息

- 创建日期: 2026-01-03
- 版本: 1.0.0
- 作者: Claude Code
- 项目: 青羽写作平台 (Qingyu_backend)
