# 阅读器系统 API 参考

> **版本**: v1.3 ⭐️已更新  
> **最后更新**: 2025-10-25  
> **基础路径**: `/api/v1/reader`  
> **主要更新**: 新增评论点赞功能、统一响应格式、阅读历史优化

---

## 1. 概述

阅读器系统提供章节管理、阅读进度保存、阅读设置等核心阅读功能，是用户实际阅读书籍内容的主要模块。

### 1.1 基础信息

- **认证要求**: ✅ **所有接口均需要登录认证**
- **响应格式**: 统一 JSON 格式
- **分页支持**: 是
- **实时同步**: 支持阅读进度跨设备同步

### 1.2 功能特性

- ✅ 章节信息获取
- ✅ 章节内容加载
- ✅ 阅读进度保存与同步
- ✅ 阅读设置个性化
- ✅ 章节导航（上一章/下一章）
- ✅ 书签管理
- ✅ 阅读历史记录
- ✅ **评论功能（发表、回复、点赞）** ⭐️v1.3新增
- ✅ **评论点赞/取消点赞** ⭐️v1.3新增

---

## 1.3 统一响应格式 ⭐️v1.3更新

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 业务数据
  },
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde"
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "详细错误信息",
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde"
}
```

### 分页响应
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...],
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde",
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false
  }
}
```

**新增字段说明**:
- `timestamp`: Unix时间戳，服务器响应时间
- `request_id`: 请求追踪ID，便于调试和日志追踪（可选字段）
- `pagination`: 分页信息对象（替代原来的平铺字段）

---

## 1.4 TypeScript 类型定义 ⭐️v1.3新增

```typescript
// src/types/api.ts

// 基础响应类型
export interface APIResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  timestamp: number;
  request_id?: string;
}

// 错误响应类型
export interface ErrorResponse {
  code: number;
  message: string;
  error?: string;
  timestamp: number;
  request_id?: string;
}

// 分页信息类型
export interface Pagination {
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
  has_next: boolean;
  has_previous: boolean;
}

// 分页响应类型
export interface PaginatedResponse<T = any> extends APIResponse<T[]> {
  pagination: Pagination;
}

// 章节信息类型
export interface ChapterInfo {
  id: string;
  bookId: string;
  title: string;
  chapterNum: number;
  wordCount: number;
  isFree: boolean;
  price: number;
  publishTime: string;
  prevChapterId: string | null;
  nextChapterId: string | null;
}

// 评论类型
export interface Comment {
  id: string;
  user_id: string;
  username: string;
  avatar: string;
  book_id: string;
  chapter_id?: string;
  content: string;
  rating?: number;
  like_count: number;
  reply_count: number;
  is_liked: boolean;
  created_at: string;
}

// 评论列表响应
export interface CommentListResponse {
  comments: Comment[];
  total: number;
  page: number;
  size: number;
}

// 阅读历史类型
export interface ReadingHistory {
  id: string;
  book_id: string;
  book_title: string;
  chapter_id: string;
  chapter_title: string;
  progress: number;
  duration: number;
  last_read_at: string;
}

// 阅读统计类型
export interface ReadingStats {
  summary: {
    total_books: number;
    total_chapters: number;
    total_duration: number;
    average_daily_duration: number;
  };
  daily_stats: Array<{
    date: string;
    books_read: number;
    chapters_read: number;
    duration: number;
  }>;
}
```

---

## 2. 接口列表

### 2.1 章节相关

#### 2.1.1 获取章节信息

**接口说明**: 根据章节ID获取章节的基础信息（不包含内容）

**请求**
```
GET /api/v1/reader/chapters/{id}
```

**认证**: 🔒 需要 JWT Token

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 章节ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "chapter123",
    "bookId": "book123",
    "title": "第一章 开始",
    "chapterNum": 1,
    "wordCount": 3500,
    "isFree": true,
    "price": 0,
    "publishTime": "2024-01-01T00:00:00Z",
    "prevChapterId": null,
    "nextChapterId": "chapter124"
  }
}
```

**错误响应**
- `404` - 章节不存在
- `401` - 未授权（未登录）

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters/chapter123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**JavaScript/Axios 示例**
```javascript
const getChapterInfo = async (chapterId) => {
  try {
    const response = await axios.get(`/api/v1/reader/chapters/${chapterId}`);
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 404) {
      console.error('章节不存在');
    }
    throw error;
  }
};
```

#### 2.1.2 获取章节内容

**接口说明**: 根据章节ID获取章节的完整内容（需要权限验证）

**请求**
```
GET /api/v1/reader/chapters/{id}/content
```

**认证**: 🔒 需要 JWT Token

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 章节ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "content": "这是章节的完整内容文本......"
  }
}
```

**错误响应**
- `401` - 未授权（未登录）
- `403` - 权限不足（需要购买或解锁）
- `404` - 章节不存在

**权限说明**:
- 免费章节：所有登录用户可读
- VIP 章节：需要 VIP 会员
- 付费章节：需要购买章节或全本

**JavaScript/Axios 示例**
```javascript
const getChapterContent = async (chapterId) => {
  try {
    const response = await axios.get(
      `/api/v1/reader/chapters/${chapterId}/content`
    );
    return response.data.data.content;
  } catch (error) {
    if (error.response?.status === 403) {
      // 处理权限不足
      console.error('需要购买或解锁');
      // 跳转到购买页面
    }
    throw error;
  }
};
```

#### 2.1.3 获取书籍章节列表

**接口说明**: 获取指定书籍的章节列表，支持分页

**请求**
```
GET /api/v1/reader/chapters
```

**认证**: 🔒 需要 JWT Token

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| bookId | string | 是 | - | 书籍ID |
| page | int | 否 | 1 | 页码 |
| size | int | 否 | 20 | 每页数量（最大100） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapters": [
      {
        "id": "chapter123",
        "title": "第一章 开始",
        "chapterNum": 1,
        "wordCount": 3500,
        "isFree": true,
        "publishTime": "2024-01-01T00:00:00Z"
      },
      {
        "id": "chapter124",
        "title": "第二章 修炼",
        "chapterNum": 2,
        "wordCount": 4200,
        "isFree": false,
        "publishTime": "2024-01-02T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters?bookId=book123&page=1&size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**JavaScript/Axios 示例**
```javascript
const getBookChapters = async (bookId, page = 1, size = 20) => {
  try {
    const response = await axios.get('/api/v1/reader/chapters', {
      params: { bookId, page, size }
    });
    return response.data.data;
  } catch (error) {
    console.error('获取章节列表失败:', error);
    throw error;
  }
};
```

#### 2.1.4 获取章节导航

**接口说明**: 获取当前章节的上一章和下一章信息

**请求**
```
GET /api/v1/reader/chapters/navigation
```

**认证**: 🔒 需要 JWT Token

**Query 参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bookId | string | 是 | 书籍ID |
| chapterNum | int | 是 | 当前章节号 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "prevChapter": {
      "id": "chapter122",
      "title": "序章",
      "chapterNum": 0
    },
    "nextChapter": {
      "id": "chapter124",
      "title": "第二章 修炼",
      "chapterNum": 2
    }
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getChapterNavigation = async (bookId, chapterNum) => {
  try {
    const response = await axios.get('/api/v1/reader/chapters/navigation', {
      params: { bookId, chapterNum }
    });
    return response.data.data;
  } catch (error) {
    console.error('获取章节导航失败:', error);
    return { prevChapter: null, nextChapter: null };
  }
};
```

#### 2.1.5 获取第一章

**接口说明**: 快速获取指定书籍的第一章

**请求**
```
GET /api/v1/reader/chapters/first
```

**认证**: 🔒 需要 JWT Token

**Query 参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bookId | string | 是 | 书籍ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "chapter123",
    "title": "第一章 开始",
    "chapterNum": 1
  }
}
```

#### 2.1.6 获取最后一章

**接口说明**: 快速获取指定书籍的最新章节

**请求**
```
GET /api/v1/reader/chapters/last
```

**认证**: 🔒 需要 JWT Token

**Query 参数**: 同"获取第一章"

---

### 2.2 阅读进度相关

#### 2.2.1 保存阅读进度

**接口说明**: 保存用户当前的阅读进度（支持跨设备同步）

**请求**
```
POST /api/v1/reader/progress
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "bookId": "book123",
  "chapterId": "chapter124",
  "chapterNum": 2,
  "progress": 65,
  "scrollPosition": 1200
}
```

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bookId | string | 是 | 书籍ID |
| chapterId | string | 是 | 当前章节ID |
| chapterNum | int | 是 | 章节号 |
| progress | int | 是 | 章节阅读进度（0-100） |
| scrollPosition | int | 否 | 滚动位置（px） |

**响应示例**
```json
{
  "code": 200,
  "message": "保存成功"
}
```

**JavaScript/Axios 示例**
```javascript
import { debounce } from 'lodash-es';

// 防抖保存进度，避免频繁请求
const saveProgressDebounced = debounce(async (progressData) => {
  try {
    await axios.post('/api/v1/reader/progress', progressData);
  } catch (error) {
    console.error('保存阅读进度失败:', error);
  }
}, 2000);

// 在滚动时调用
const handleScroll = () => {
  const progress = calculateProgress(); // 计算当前进度
  const scrollPosition = window.pageYOffset;
  
  saveProgressDebounced({
    bookId: currentBook.value.id,
    chapterId: currentChapter.value.id,
    chapterNum: currentChapter.value.chapterNum,
    progress,
    scrollPosition
  });
};
```

#### 2.2.2 获取阅读进度

**接口说明**: 获取用户在指定书籍的阅读进度

**请求**
```
GET /api/v1/reader/progress
```

**认证**: 🔒 需要 JWT Token

**Query 参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bookId | string | 是 | 书籍ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "bookId": "book123",
    "chapterId": "chapter124",
    "chapterNum": 2,
    "chapterTitle": "第二章 修炼",
    "progress": 65,
    "scrollPosition": 1200,
    "lastReadTime": "2024-10-18T10:30:00Z"
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getReadingProgress = async (bookId) => {
  try {
    const response = await axios.get('/api/v1/reader/progress', {
      params: { bookId }
    });
    return response.data.data;
  } catch (error) {
    console.error('获取阅读进度失败:', error);
    return null;
  }
};
```

---

### 2.3 阅读设置相关

#### 2.3.1 获取阅读设置

**接口说明**: 获取用户的个性化阅读设置

**请求**
```
GET /api/v1/reader/settings
```

**认证**: 🔒 需要 JWT Token

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "fontSize": 18,
    "fontFamily": "SimSun",
    "lineHeight": 1.8,
    "theme": "light",
    "pageWidth": 800,
    "autoSave": true,
    "pageMode": "scroll"
  }
}
```

#### 2.3.2 更新阅读设置

**接口说明**: 更新用户的阅读设置

**请求**
```
POST /api/v1/reader/settings
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "fontSize": 20,
  "fontFamily": "Microsoft YaHei",
  "lineHeight": 2.0,
  "theme": "dark",
  "pageWidth": 900,
  "autoSave": true,
  "pageMode": "page"
}
```

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| fontSize | int | 否 | 字体大小（12-30） |
| fontFamily | string | 否 | 字体（SimSun/Microsoft YaHei/KaiTi） |
| lineHeight | float | 否 | 行高（1.0-3.0） |
| theme | string | 否 | 主题（light/dark/sepia） |
| pageWidth | int | 否 | 页面宽度（600-1200） |
| autoSave | boolean | 否 | 自动保存进度 |
| pageMode | string | 否 | 翻页模式（scroll/page） |

**响应示例**
```json
{
  "code": 200,
  "message": "设置已保存"
}
```

**JavaScript/Axios 示例**
```javascript
const updateReaderSettings = async (settings) => {
  try {
    await axios.post('/api/v1/reader/settings', settings);
  } catch (error) {
    console.error('保存阅读设置失败:', error);
    throw error;
  }
};
```

---

### 2.4 书签相关

#### 2.4.1 添加书签

**接口说明**: 在当前阅读位置添加书签

**请求**
```
POST /api/v1/reader/bookmarks
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "bookId": "book123",
  "chapterId": "chapter124",
  "chapterTitle": "第二章 修炼",
  "position": 1200,
  "note": "重要情节"
}
```

**响应示例**
```json
{
  "code": 200,
  "message": "书签添加成功",
  "data": {
    "id": "bookmark123",
    "createdAt": "2024-10-18T10:30:00Z"
  }
}
```

#### 2.4.2 获取书签列表

**接口说明**: 获取指定书籍的所有书签

**请求**
```
GET /api/v1/reader/bookmarks
```

**认证**: 🔒 需要 JWT Token

**Query 参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bookId | string | 是 | 书籍ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "bookmark123",
      "chapterId": "chapter124",
      "chapterTitle": "第二章 修炼",
      "position": 1200,
      "note": "重要情节",
      "createdAt": "2024-10-18T10:30:00Z"
    }
  ]
}
```

#### 2.4.3 删除书签

**接口说明**: 删除指定的书签

**请求**
```
DELETE /api/v1/reader/bookmarks/{id}
```

**认证**: 🔒 需要 JWT Token

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 书签ID |

**响应示例**
```json
{
  "code": 200,
  "message": "书签删除成功"
}
```

---

## 3. 数据结构

### 3.1 Chapter（章节）

```typescript
interface Chapter {
  id: string;                 // 章节ID
  bookId: string;             // 书籍ID
  title: string;              // 章节标题
  chapterNum: number;         // 章节序号
  wordCount: number;          // 字数
  isFree: boolean;            // 是否免费
  price: number;              // 价格（付费章节）
  publishTime: string;        // 发布时间（ISO 8601）
  prevChapterId: string | null; // 上一章ID
  nextChapterId: string | null; // 下一章ID
}
```

### 3.2 ReadingProgress（阅读进度）

```typescript
interface ReadingProgress {
  bookId: string;             // 书籍ID
  chapterId: string;          // 当前章节ID
  chapterNum: number;         // 章节号
  chapterTitle: string;       // 章节标题
  progress: number;           // 章节阅读进度（0-100）
  scrollPosition: number;     // 滚动位置（px）
  lastReadTime: string;       // 最后阅读时间（ISO 8601）
}
```

### 3.3 ReaderSettings（阅读设置）

```typescript
interface ReaderSettings {
  fontSize: number;           // 字体大小（12-30）
  fontFamily: string;         // 字体
  lineHeight: number;         // 行高（1.0-3.0）
  theme: string;              // 主题（light/dark/sepia）
  pageWidth: number;          // 页面宽度（600-1200）
  autoSave: boolean;          // 自动保存进度
  pageMode: string;           // 翻页模式（scroll/page）
}
```

### 3.4 Bookmark（书签）

```typescript
interface Bookmark {
  id: string;                 // 书签ID
  bookId: string;             // 书籍ID
  chapterId: string;          // 章节ID
  chapterTitle: string;       // 章节标题
  position: number;           // 位置（滚动位置）
  note: string;               // 备注
  createdAt: string;          // 创建时间（ISO 8601）
}
```

---

## 4. 完整示例

### 4.1 阅读器页面组件（Vue 3）

```vue
<template>
  <div class="reader-page" :style="readerStyles">
    <!-- 顶部工具栏 -->
    <div class="reader-header">
      <el-button @click="goBack">返回</el-button>
      <span class="chapter-title">{{ currentChapter?.title }}</span>
      <el-button @click="showSettings = true">设置</el-button>
    </div>

    <!-- 章节内容 -->
    <div 
      ref="contentRef" 
      class="chapter-content"
      @scroll="handleScroll"
    >
      <h1 class="chapter-title">{{ currentChapter?.title }}</h1>
      <div class="chapter-text" v-html="chapterContent"></div>
    </div>

    <!-- 底部导航 -->
    <div class="reader-footer">
      <el-button 
        :disabled="!navigation.prevChapter" 
        @click="goToPrevChapter"
      >
        上一章
      </el-button>
      <el-button @click="showBookmarks = true">书签</el-button>
      <el-button 
        :disabled="!navigation.nextChapter" 
        @click="goToNextChapter"
      >
        下一章
      </el-button>
    </div>

    <!-- 设置抽屉 -->
    <el-drawer v-model="showSettings" title="阅读设置">
      <el-form :model="settings" label-width="80px">
        <el-form-item label="字体大小">
          <el-slider v-model="settings.fontSize" :min="12" :max="30" />
        </el-form-item>
        <el-form-item label="字体">
          <el-select v-model="settings.fontFamily">
            <el-option label="宋体" value="SimSun" />
            <el-option label="微软雅黑" value="Microsoft YaHei" />
            <el-option label="楷体" value="KaiTi" />
          </el-select>
        </el-form-item>
        <el-form-item label="行高">
          <el-slider v-model="settings.lineHeight" :min="1.0" :max="3.0" :step="0.1" />
        </el-form-item>
        <el-form-item label="主题">
          <el-radio-group v-model="settings.theme">
            <el-radio label="light">明亮</el-radio>
            <el-radio label="dark">暗黑</el-radio>
            <el-radio label="sepia">护眼</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveSettings">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- 书签列表 -->
    <el-drawer v-model="showBookmarks" title="我的书签">
      <div v-for="bookmark in bookmarks" :key="bookmark.id" class="bookmark-item">
        <div @click="goToBookmark(bookmark)">
          <h4>{{ bookmark.chapterTitle }}</h4>
          <p>{{ bookmark.note }}</p>
          <small>{{ formatDate(bookmark.createdAt) }}</small>
        </div>
        <el-button type="danger" size="small" @click="deleteBookmark(bookmark.id)">
          删除
        </el-button>
      </div>
      <el-button type="primary" @click="addBookmark">添加书签</el-button>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { debounce } from 'lodash-es';

const route = useRoute();
const router = useRouter();

// 数据
const currentChapter = ref(null);
const chapterContent = ref('');
const navigation = ref({ prevChapter: null, nextChapter: null });
const settings = ref({
  fontSize: 18,
  fontFamily: 'SimSun',
  lineHeight: 1.8,
  theme: 'light',
  pageWidth: 800
});
const bookmarks = ref([]);
const showSettings = ref(false);
const showBookmarks = ref(false);
const contentRef = ref(null);

// 计算样式
const readerStyles = computed(() => ({
  '--font-size': settings.value.fontSize + 'px',
  '--font-family': settings.value.fontFamily,
  '--line-height': settings.value.lineHeight,
  '--page-width': settings.value.pageWidth + 'px',
  '--bg-color': themeColors[settings.value.theme].bg,
  '--text-color': themeColors[settings.value.theme].text
}));

const themeColors = {
  light: { bg: '#fff', text: '#333' },
  dark: { bg: '#1a1a1a', text: '#e0e0e0' },
  sepia: { bg: '#f5f1e8', text: '#4a4a4a' }
};

// 加载章节内容
const loadChapter = async (chapterId) => {
  try {
    // 加载章节信息
    const chapterInfo = await axios.get(`/api/v1/reader/chapters/${chapterId}`);
    currentChapter.value = chapterInfo.data.data;

    // 加载章节内容
    const content = await axios.get(`/api/v1/reader/chapters/${chapterId}/content`);
    chapterContent.value = content.data.data.content;

    // 加载导航信息
    const nav = await axios.get('/api/v1/reader/chapters/navigation', {
      params: {
        bookId: currentChapter.value.bookId,
        chapterNum: currentChapter.value.chapterNum
      }
    });
    navigation.value = nav.data.data;

    // 恢复阅读进度
    await restoreProgress();
  } catch (error) {
    console.error('加载章节失败:', error);
    if (error.response?.status === 403) {
      // 跳转到购买页面
      router.push(`/purchase/${currentChapter.value.bookId}`);
    }
  }
};

// 恢复阅读进度
const restoreProgress = async () => {
  try {
    const progress = await axios.get('/api/v1/reader/progress', {
      params: { bookId: currentChapter.value.bookId }
    });
    
    if (progress.data.data?.scrollPosition) {
      contentRef.value.scrollTop = progress.data.data.scrollPosition;
    }
  } catch (error) {
    console.warn('恢复阅读进度失败:', error);
  }
};

// 保存阅读进度（防抖）
const saveProgress = debounce(async () => {
  if (!currentChapter.value) return;

  const scrollTop = contentRef.value.scrollTop;
  const scrollHeight = contentRef.value.scrollHeight;
  const clientHeight = contentRef.value.clientHeight;
  const progress = Math.round((scrollTop / (scrollHeight - clientHeight)) * 100);

  try {
    await axios.post('/api/v1/reader/progress', {
      bookId: currentChapter.value.bookId,
      chapterId: currentChapter.value.id,
      chapterNum: currentChapter.value.chapterNum,
      progress,
      scrollPosition: scrollTop
    });
  } catch (error) {
    console.error('保存阅读进度失败:', error);
  }
}, 2000);

// 滚动处理
const handleScroll = () => {
  saveProgress();
};

// 章节导航
const goToPrevChapter = () => {
  if (navigation.value.prevChapter) {
    router.push(`/reader/${navigation.value.prevChapter.id}`);
  }
};

const goToNextChapter = () => {
  if (navigation.value.nextChapter) {
    router.push(`/reader/${navigation.value.nextChapter.id}`);
  }
};

// 阅读设置
const loadSettings = async () => {
  try {
    const response = await axios.get('/api/v1/reader/settings');
    settings.value = response.data.data;
  } catch (error) {
    console.error('加载阅读设置失败:', error);
  }
};

const saveSettings = async () => {
  try {
    await axios.post('/api/v1/reader/settings', settings.value);
    showSettings.value = false;
  } catch (error) {
    console.error('保存阅读设置失败:', error);
  }
};

// 书签管理
const loadBookmarks = async () => {
  if (!currentChapter.value) return;
  
  try {
    const response = await axios.get('/api/v1/reader/bookmarks', {
      params: { bookId: currentChapter.value.bookId }
    });
    bookmarks.value = response.data.data;
  } catch (error) {
    console.error('加载书签失败:', error);
  }
};

const addBookmark = async () => {
  try {
    await axios.post('/api/v1/reader/bookmarks', {
      bookId: currentChapter.value.bookId,
      chapterId: currentChapter.value.id,
      chapterTitle: currentChapter.value.title,
      position: contentRef.value.scrollTop,
      note: '书签' + (bookmarks.value.length + 1)
    });
    await loadBookmarks();
  } catch (error) {
    console.error('添加书签失败:', error);
  }
};

const deleteBookmark = async (id) => {
  try {
    await axios.delete(`/api/v1/reader/bookmarks/${id}`);
    await loadBookmarks();
  } catch (error) {
    console.error('删除书签失败:', error);
  }
};

const goToBookmark = (bookmark) => {
  if (bookmark.chapterId !== currentChapter.value.id) {
    router.push(`/reader/${bookmark.chapterId}`);
  }
  contentRef.value.scrollTop = bookmark.position;
  showBookmarks.value = false;
};

// 生命周期
onMounted(async () => {
  await loadSettings();
  await loadChapter(route.params.chapterId);
  await loadBookmarks();
});

// 监听路由变化
watch(() => route.params.chapterId, (newId) => {
  if (newId) {
    loadChapter(newId);
  }
});

// 返回
const goBack = () => {
  router.back();
};

// 工具函数
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('zh-CN');
};
</script>

<style scoped>
.reader-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-color);
  color: var(--text-color);
}

.chapter-content {
  flex: 1;
  overflow-y: auto;
  padding: 40px 20px;
  max-width: var(--page-width);
  margin: 0 auto;
}

.chapter-text {
  font-size: var(--font-size);
  font-family: var(--font-family);
  line-height: var(--line-height);
  text-align: justify;
  white-space: pre-wrap;
}

.reader-header,
.reader-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: rgba(0, 0, 0, 0.05);
}

.bookmark-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}

.bookmark-item:hover {
  background-color: #f5f5f5;
}
</style>
```

---

## 5. 最佳实践

### 5.1 阅读进度自动保存

```javascript
// 使用 requestIdleCallback 在浏览器空闲时保存进度
const saveProgressIdle = (progressData) => {
  if ('requestIdleCallback' in window) {
    requestIdleCallback(() => {
      axios.post('/api/v1/reader/progress', progressData).catch(() => {});
    });
  } else {
    setTimeout(() => {
      axios.post('/api/v1/reader/progress', progressData).catch(() => {});
    }, 100);
  }
};
```

### 5.2 章节内容预加载

```javascript
// 预加载下一章，提升阅读体验
const prefetchNextChapter = async (nextChapterId) => {
  if (!nextChapterId) return;
  
  try {
    // 使用 fetch API 预加载
    fetch(`/api/v1/reader/chapters/${nextChapterId}/content`, {
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
  } catch (error) {
    // 忽略错误
  }
};

// 在用户阅读到 70% 时预加载下一章
watch(() => readProgress.value, (progress) => {
  if (progress >= 70 && navigation.value.nextChapter) {
    prefetchNextChapter(navigation.value.nextChapter.id);
  }
});
```

### 5.3 离线阅读支持

```javascript
// 使用 IndexedDB 缓存章节内容
import { openDB } from 'idb';

const dbPromise = openDB('reader-cache', 1, {
  upgrade(db) {
    db.createObjectStore('chapters');
  }
});

const cacheChapter = async (chapterId, content) => {
  const db = await dbPromise;
  await db.put('chapters', content, chapterId);
};

const getCachedChapter = async (chapterId) => {
  const db = await dbPromise;
  return await db.get('chapters', chapterId);
};

// 加载章节时先尝试从缓存读取
const loadChapterWithCache = async (chapterId) => {
  // 先读缓存
  const cached = await getCachedChapter(chapterId);
  if (cached) {
    chapterContent.value = cached;
  }
  
  // 异步更新
  try {
    const response = await axios.get(`/api/v1/reader/chapters/${chapterId}/content`);
    const content = response.data.data.content;
    chapterContent.value = content;
    await cacheChapter(chapterId, content);
  } catch (error) {
    if (!cached) throw error;
  }
};
```

---

## 6. 评论功能 ⭐️v1.3新增

### 6.1 发表评论

**接口说明**: 对书籍或章节发表评论

**请求**
```
POST /api/v1/reader/comments
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "book_id": "book123",
  "chapter_id": "chapter456",
  "content": "这章写得真不错！",
  "rating": 5
}
```

**响应示例**
```json
{
  "code": 201,
  "message": "发表评论成功",
  "data": {
    "id": "comment789",
    "user_id": "user123",
    "book_id": "book123",
    "chapter_id": "chapter456",
    "content": "这章写得真不错！",
    "rating": 5,
    "like_count": 0,
    "reply_count": 0,
    "is_liked": false,
    "created_at": "2025-10-25T10:00:00Z"
  },
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

### 6.2 获取评论列表

**接口说明**: 获取书籍的评论列表

**请求**
```
GET /api/v1/reader/comments?book_id={bookId}&sortBy=latest&page=1&size=20
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| book_id | string | 是 | 书籍ID |
| sortBy | string | 否 | 排序方式：latest（最新）/hot（最热），默认latest |
| page | int | 否 | 页码，默认1 |
| size | int | 否 | 每页数量，默认20 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "comments": [
      {
        "id": "comment789",
        "user_id": "user123",
        "username": "张三",
        "avatar": "https://avatar.url",
        "book_id": "book123",
        "content": "这章写得真不错！",
        "rating": 5,
        "like_count": 10,
        "reply_count": 3,
        "is_liked": false,
        "created_at": "2025-10-25T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  },
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

### 6.3 点赞评论 ⭐️v1.3新增

**接口说明**: 对评论进行点赞

**请求**
```
POST /api/v1/reader/comments/:id/like
```

**认证**: 🔒 需要 JWT Token

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 评论ID |

**响应示例**
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": null,
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

### 6.4 取消点赞 ⭐️v1.3新增

**接口说明**: 取消对评论的点赞

**请求**
```
DELETE /api/v1/reader/comments/:id/like
```

**认证**: 🔒 需要 JWT Token

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 评论ID |

**响应示例**
```json
{
  "code": 200,
  "message": "取消点赞成功",
  "data": null,
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

### 6.5 回复评论

**接口说明**: 回复一条评论

**请求**
```
POST /api/v1/reader/comments/:id/reply
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "content": "我也觉得很精彩！"
}
```

**响应示例**
```json
{
  "code": 201,
  "message": "回复成功",
  "data": {
    "id": "comment790",
    "parent_comment_id": "comment789",
    "user_id": "user456",
    "content": "我也觉得很精彩！",
    "created_at": "2025-10-25T10:05:00Z"
  },
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

### 6.6 TypeScript API 封装

**评论相关API封装**:
```typescript
// src/api/reader.ts
import type { APIResponse, Comment, CommentListResponse } from '@/types/api';
import request from '@/utils/request';

// 发表评论
export interface CreateCommentParams {
  book_id: string;
  chapter_id?: string;
  content: string;
  rating?: number;
}

export const createComment = (data: CreateCommentParams) => {
  return request.post<APIResponse<Comment>>('/reader/comments', data);
};

// 获取评论列表
export interface GetCommentListParams {
  book_id: string;
  sortBy?: 'latest' | 'hot';
  page?: number;
  size?: number;
}

export const getCommentList = (params: GetCommentListParams) => {
  return request.get<APIResponse<CommentListResponse>>('/reader/comments', { params });
};

// 点赞评论
export const likeComment = (commentId: string) => {
  return request.post<APIResponse<null>>(`/reader/comments/${commentId}/like`);
};

// 取消点赞
export const unlikeComment = (commentId: string) => {
  return request.delete<APIResponse<null>>(`/reader/comments/${commentId}/like`);
};

// 回复评论
export const replyComment = (commentId: string, content: string) => {
  return request.post<APIResponse<Comment>>(`/reader/comments/${commentId}/reply`, { content });
};
```

### 6.7 Vue 3 + TypeScript 组件示例

**评论列表组件**:
```vue
<template>
  <div class="comment-list">
    <!-- 评论列表 -->
    <div v-for="comment in comments" :key="comment.id" class="comment-item">
      <div class="comment-header">
        <img :src="comment.avatar" class="avatar" alt="avatar" />
        <span class="username">{{ comment.username }}</span>
        <span class="time">{{ formatTime(comment.created_at) }}</span>
      </div>
      
      <div class="comment-content">{{ comment.content }}</div>
      
      <div class="comment-actions">
        <!-- 点赞按钮 -->
        <button 
          @click="toggleLike(comment)" 
          :class="{ 'active': comment.is_liked }"
          class="like-btn"
        >
          <i :class="comment.is_liked ? 'icon-liked' : 'icon-like'"></i>
          {{ comment.like_count }}
        </button>
        
        <!-- 回复按钮 -->
        <button @click="showReply(comment)" class="reply-btn">
          <i class="icon-reply"></i>
          回复 ({{ comment.reply_count }})
        </button>
      </div>
    </div>
    
    <!-- 加载更多 -->
    <button v-if="hasMore" @click="loadMore" class="load-more">
      加载更多
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getCommentList, likeComment, unlikeComment } from '@/api/reader';
import type { Comment } from '@/types/api';

interface Props {
  bookId: string;
}

const props = defineProps<Props>();

const comments = ref<Comment[]>([]);
const page = ref(1);
const hasMore = ref(true);

// 加载评论列表
const loadComments = async () => {
  try {
    const response = await getCommentList({
      book_id: props.bookId,
      sortBy: 'latest',
      page: page.value,
      size: 20
    });
    
    const { comments: newComments, total } = response.data!;
    
    if (page.value === 1) {
      comments.value = newComments;
    } else {
      comments.value.push(...newComments);
    }
    
    hasMore.value = comments.value.length < total;
  } catch (error) {
    ElMessage.error('加载评论失败');
  }
};

// 点赞/取消点赞
const toggleLike = async (comment: Comment) => {
  try {
    if (comment.is_liked) {
      await unlikeComment(comment.id);
      comment.is_liked = false;
      comment.like_count--;
      ElMessage.success('已取消点赞');
    } else {
      await likeComment(comment.id);
      comment.is_liked = true;
      comment.like_count++;
      ElMessage.success('点赞成功');
    }
  } catch (error) {
    ElMessage.error('操作失败');
  }
};

// 显示回复框
const showReply = (comment: Comment) => {
  // 实现回复逻辑
  console.log('回复评论:', comment.id);
};

// 加载更多
const loadMore = () => {
  page.value++;
  loadComments();
};

// 格式化时间
const formatTime = (time: string): string => {
  const date = new Date(time);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  
  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`;
  return date.toLocaleDateString();
};

onMounted(() => {
  loadComments();
});
</script>

<style scoped>
.comment-list {
  padding: 16px;
}

.comment-item {
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.username {
  font-weight: bold;
  font-size: 14px;
}

.time {
  color: #999;
  font-size: 12px;
}

.comment-content {
  margin: 12px 0;
  line-height: 1.6;
}

.comment-actions {
  display: flex;
  gap: 16px;
}

.like-btn,
.reply-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
}

.like-btn.active {
  color: #f56c6c;
  border-color: #f56c6c;
  background: #fef0f0;
}

.like-btn:hover,
.reply-btn:hover {
  border-color: #409eff;
  color: #409eff;
}

.load-more {
  width: 100%;
  padding: 12px;
  margin-top: 16px;
  border: 1px dashed #ddd;
  background: #fff;
  color: #666;
  cursor: pointer;
  border-radius: 4px;
}

.load-more:hover {
  border-color: #409eff;
  color: #409eff;
}
</style>
```

---

## 7. 阅读历史 ⭐️v1.3优化

### 7.1 记录阅读历史

**接口说明**: 记录用户的阅读行为

**请求**
```
POST /api/v1/reader/reading-history
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "book_id": "book123",
  "chapter_id": "chapter456",
  "start_time": "2025-10-25T10:00:00Z",
  "end_time": "2025-10-25T10:15:00Z",
  "progress": 75.5,
  "device_type": "web",
  "device_id": "device_123"
}
```

**响应示例**
```json
{
  "code": 201,
  "message": "记录成功",
  "data": {},
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

> **注意**: 响应状态码为 `201 Created`，遵循 RESTful 规范。

### 7.2 获取阅读历史

**接口说明**: 获取用户的阅读历史记录

**请求**
```
GET /api/v1/reader/reading-history?page=1&page_size=20&book_id=book123
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| book_id | string | 否 | 筛选指定书籍 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "histories": [
      {
        "id": "history123",
        "book_id": "book123",
        "book_title": "修仙传",
        "chapter_id": "chapter456",
        "chapter_title": "第一章",
        "progress": 75.5,
        "duration": 900,
        "last_read_at": "2025-10-25T10:15:00Z"
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "page_size": 20,
      "total_pages": 3,
      "has_next": true,
      "has_previous": false
    }
  },
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

### 7.3 获取阅读统计

**接口说明**: 获取用户的阅读统计数据

**请求**
```
GET /api/v1/reader/reading-history/stats?days=30
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "summary": {
      "total_books": 15,
      "total_chapters": 120,
      "total_duration": 36000,
      "average_daily_duration": 1200
    },
    "daily_stats": [
      {
        "date": "2025-10-25",
        "books_read": 3,
        "chapters_read": 8,
        "duration": 1800
      }
    ]
  },
  "timestamp": 1729875123,
  "request_id": "req-xxx"
}
```

---

## 8. 常见问题

### Q1: 如何处理付费章节的购买流程？

**A**: 检测 403 错误并引导用户购买：

```javascript
const loadChapterContent = async (chapterId) => {
  try {
    const response = await axios.get(`/api/v1/reader/chapters/${chapterId}/content`);
    return response.data.data.content;
  } catch (error) {
    if (error.response?.status === 403) {
      // 弹出购买对话框
      showPurchaseDialog.value = true;
      purchaseInfo.value = {
        chapterId,
        price: currentChapter.value.price
      };
      throw new Error('需要购买章节');
    }
    throw error;
  }
};
```

### Q2: 如何实现跨设备阅读进度同步？

**A**: 进度数据已经自动保存在服务器，只需在打开书籍时恢复：

```javascript
const openBook = async (bookId) => {
  try {
    // 获取阅读进度
    const progressRes = await axios.get('/api/v1/reader/progress', {
      params: { bookId }
    });
    
    const progress = progressRes.data.data;
    if (progress && progress.chapterId) {
      // 跳转到上次阅读的章节
      router.push(`/reader/${progress.chapterId}`);
    } else {
      // 从第一章开始
      const firstChapterRes = await axios.get('/api/v1/reader/chapters/first', {
        params: { bookId }
      });
      router.push(`/reader/${firstChapterRes.data.data.id}`);
    }
  } catch (error) {
    console.error('打开书籍失败:', error);
  }
};
```

### Q3: 如何优化长章节的渲染性能？

**A**: 使用虚拟滚动或分段加载：

```javascript
// 分段加载长内容
const loadChapterInChunks = async (chapterId) => {
  const response = await axios.get(`/api/v1/reader/chapters/${chapterId}/content`);
  const fullContent = response.data.data.content;
  
  // 分割成段落
  const paragraphs = fullContent.split('\n\n');
  const chunkSize = 50; // 每次渲染50个段落
  
  // 初始渲染
  renderedParagraphs.value = paragraphs.slice(0, chunkSize);
  
  // 滚动到底部时加载更多
  const loadMore = () => {
    const currentLength = renderedParagraphs.value.length;
    if (currentLength < paragraphs.length) {
      renderedParagraphs.value.push(
        ...paragraphs.slice(currentLength, currentLength + chunkSize)
      );
    }
  };
  
  return loadMore;
};
```

---

**文档版本**: v1.3 ⭐️  
**最后更新**: 2025-10-25  
**维护者**: 青羽后端团队

**v1.3 主要更新**:
- ✅ 新增评论点赞/取消点赞功能
- ✅ 统一响应格式（timestamp, request_id）
- ✅ 优化阅读历史API
- ✅ 新增完整的Vue组件示例

