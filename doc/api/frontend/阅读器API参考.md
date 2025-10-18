# 阅读器系统 API 参考

> **版本**: v1.0  
> **最后更新**: 2025-10-18  
> **基础路径**: `/api/v1/reader`

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

## 6. 常见问题

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

**文档版本**: v1.0  
**最后更新**: 2025-10-18  
**维护者**: 青羽后端团队

