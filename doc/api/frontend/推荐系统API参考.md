# 推荐系统 API 参考

> **版本**: v1.0  
> **最后更新**: 2025-10-18  
> **基础路径**: `/api/v1/recommendation`

---

## 1. 概述

推荐系统通过分析用户行为和内容特征，为用户提供个性化的书籍推荐，提升用户发现好书的效率和阅读体验。

### 1.1 基础信息

- **认证要求**: 部分接口需要登录（个性化推荐、行为记录）
- **响应格式**: 统一 JSON 格式
- **推荐策略**: 协同过滤 + 内容推荐 + 热门推荐

### 1.2 功能特性

- ✅ 个性化推荐（基于用户历史行为）
- ✅ 相似物品推荐（看了这本的还看了...）
- ✅ 用户行为追踪
- ✅ 首页混合推荐
- ✅ 热门推荐
- ✅ 分类推荐

---

## 2. 接口列表

### 2.1 个性化推荐

#### 2.1.1 获取个性化推荐

**接口说明**: 根据用户历史阅读行为，推荐可能感兴趣的书籍

**请求**
```
GET /api/v1/recommendation/personalized
```

**认证**: 🔒 需要 JWT Token

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| limit | int | 否 | 10 | 推荐数量（最大50） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "recommendations": [
      {
        "bookId": "book123",
        "title": "示例书籍",
        "author": "作者名",
        "cover": "https://example.com/cover.jpg",
        "rating": 4.5,
        "reason": "因为您喜欢《类似书籍》",
        "score": 0.92
      }
    ],
    "count": 10
  }
}
```

**推荐理由类型**:
- "因为您喜欢《xxx》"
- "基于您的阅读偏好"
- "与您阅读的书籍相似"

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/recommendation/personalized?limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**JavaScript/Axios 示例**
```javascript
const getPersonalizedRecommendations = async (limit = 10) => {
  try {
    const response = await axios.get('/api/v1/recommendation/personalized', {
      params: { limit }
    });
    return response.data.data.recommendations;
  } catch (error) {
    console.error('获取个性化推荐失败:', error);
    throw error;
  }
};
```

---

### 2.2 相似推荐

#### 2.2.1 获取相似物品推荐

**接口说明**: 根据指定书籍，推荐相似的书籍（"看了这本的还看了..."）

**请求**
```
GET /api/v1/recommendation/similar
```

**认证**: ❌ 无需认证（但登录后推荐质量更高）

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| itemId | string | 是 | - | 书籍ID |
| limit | int | 否 | 10 | 推荐数量 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "similar_items": [
      {
        "bookId": "book456",
        "title": "相似书籍",
        "author": "作者B",
        "cover": "https://example.com/cover2.jpg",
        "rating": 4.3,
        "similarity": 0.88
      }
    ],
    "count": 10
  }
}
```

**相似度说明**:
- `similarity`: 相似度分数（0-1），越高越相似
- 基于内容特征（分类、标签、作者）和协同过滤

**JavaScript/Axios 示例**
```javascript
const getSimilarBooks = async (bookId, limit = 10) => {
  try {
    const response = await axios.get('/api/v1/recommendation/similar', {
      params: { itemId: bookId, limit }
    });
    return response.data.data.similar_items;
  } catch (error) {
    console.error('获取相似推荐失败:', error);
    return [];
  }
};
```

---

### 2.3 行为追踪

#### 2.3.1 记录用户行为

**接口说明**: 记录用户的阅读行为，用于优化推荐算法

**请求**
```
POST /api/v1/recommendation/behavior
```

**认证**: 🔒 需要 JWT Token

**请求体**
```json
{
  "itemId": "book123",
  "chapterId": "chapter456",
  "behaviorType": "read",
  "value": 1.0,
  "metadata": {
    "readTime": 300,
    "progress": 50
  }
}
```

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| itemId | string | 是 | 书籍ID |
| chapterId | string | 否 | 章节ID |
| behaviorType | string | 是 | 行为类型（view/click/collect/read/finish/like/share） |
| value | float | 否 | 行为权重（默认1.0） |
| metadata | object | 否 | 附加元数据 |

**行为类型说明**:
| 类型 | 权重 | 说明 |
|------|------|------|
| view | 0.1 | 浏览书籍详情 |
| click | 0.2 | 点击书籍封面 |
| collect | 0.5 | 收藏书籍 |
| read | 0.7 | 开始阅读 |
| finish | 1.0 | 完成阅读 |
| like | 0.8 | 点赞书籍 |
| share | 0.6 | 分享书籍 |

**响应示例**
```json
{
  "code": 200,
  "message": "记录成功"
}
```

**JavaScript/Axios 示例**
```javascript
// 行为追踪工具类
class BehaviorTracker {
  // 记录浏览行为
  static async trackView(bookId) {
    try {
      await axios.post('/api/v1/recommendation/behavior', {
        itemId: bookId,
        behaviorType: 'view',
        value: 0.1
      });
    } catch (error) {
      console.warn('记录浏览行为失败:', error);
    }
  }

  // 记录阅读行为
  static async trackRead(bookId, chapterId, readTime) {
    try {
      await axios.post('/api/v1/recommendation/behavior', {
        itemId: bookId,
        chapterId: chapterId,
        behaviorType: 'read',
        value: 0.7,
        metadata: { readTime }
      });
    } catch (error) {
      console.warn('记录阅读行为失败:', error);
    }
  }

  // 记录收藏行为
  static async trackCollect(bookId) {
    try {
      await axios.post('/api/v1/recommendation/behavior', {
        itemId: bookId,
        behaviorType: 'collect',
        value: 0.5
      });
    } catch (error) {
      console.warn('记录收藏行为失败:', error);
    }
  }

  // 记录点赞行为
  static async trackLike(bookId) {
    try {
      await axios.post('/api/v1/recommendation/behavior', {
        itemId: bookId,
        behaviorType: 'like',
        value: 0.8
      });
    } catch (error) {
      console.warn('记录点赞行为失败:', error);
    }
  }
}

// 使用示例
// 在书籍详情页
onMounted(() => {
  BehaviorTracker.trackView(bookId.value);
});

// 开始阅读时
const startReading = () => {
  BehaviorTracker.trackRead(bookId.value, chapterId.value, 0);
  router.push(`/reader/${chapterId.value}`);
};

// 收藏书籍时
const collectBook = async () => {
  await collectAPI.addToLibrary(bookId.value);
  BehaviorTracker.trackCollect(bookId.value);
};
```

---

### 2.4 混合推荐

#### 2.4.1 获取首页推荐

**接口说明**: 获取首页的混合推荐（包含个性化 + 热门 + 新书）

**请求**
```
GET /api/v1/recommendation/homepage
```

**认证**: ❌ 无需认证（登录后包含个性化推荐）

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| limit | int | 否 | 20 | 推荐数量 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "recommendations": [
      {
        "bookId": "book123",
        "title": "推荐书籍",
        "author": "作者名",
        "cover": "https://example.com/cover.jpg",
        "rating": 4.5,
        "tag": "个性化推荐",
        "reason": "基于您的阅读偏好"
      },
      {
        "bookId": "book456",
        "title": "热门书籍",
        "author": "作者B",
        "cover": "https://example.com/cover2.jpg",
        "rating": 4.8,
        "tag": "热门推荐",
        "reason": "近期热门书籍"
      }
    ],
    "count": 20
  }
}
```

**推荐标签类型**:
- "个性化推荐" - 基于用户历史
- "热门推荐" - 近期热度高
- "新书推荐" - 最新上架
- "编辑精选" - 编辑推荐

**JavaScript/Axios 示例**
```javascript
const getHomepageRecommendations = async (limit = 20) => {
  try {
    const response = await axios.get('/api/v1/recommendation/homepage', {
      params: { limit }
    });
    return response.data.data.recommendations;
  } catch (error) {
    console.error('获取首页推荐失败:', error);
    return [];
  }
};
```

---

### 2.5 热门推荐

#### 2.5.1 获取热门推荐

**接口说明**: 获取近期热门书籍推荐

**请求**
```
GET /api/v1/recommendation/hot
```

**认证**: ❌ 无需认证

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| limit | int | 否 | 20 | 推荐数量 |
| days | int | 否 | 7 | 统计天数（7/15/30） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "recommendations": [
      {
        "bookId": "book123",
        "title": "热门书籍",
        "author": "作者名",
        "cover": "https://example.com/cover.jpg",
        "rating": 4.7,
        "hotScore": 9500,
        "trend": "up"
      }
    ],
    "count": 20
  }
}
```

**热度分数说明**:
- `hotScore`: 综合热度分数（浏览量 + 阅读量 + 收藏量等权重计算）
- `trend`: 趋势（up/down/stable）

**JavaScript/Axios 示例**
```javascript
const getHotBooks = async (days = 7, limit = 20) => {
  try {
    const response = await axios.get('/api/v1/recommendation/hot', {
      params: { days, limit }
    });
    return response.data.data.recommendations;
  } catch (error) {
    console.error('获取热门推荐失败:', error);
    return [];
  }
};
```

---

### 2.6 分类推荐

#### 2.6.1 获取分类推荐

**接口说明**: 获取指定分类的推荐书籍

**请求**
```
GET /api/v1/recommendation/category
```

**认证**: ❌ 无需认证（登录后推荐质量更高）

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| category | string | 是 | - | 分类名称或ID |
| limit | int | 否 | 20 | 推荐数量 |

**响应示例**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "recommendations": [
      {
        "bookId": "book123",
        "title": "玄幻书籍",
        "author": "作者名",
        "cover": "https://example.com/cover.jpg",
        "rating": 4.6,
        "categoryScore": 0.95
      }
    ],
    "count": 20,
    "category": "玄幻"
  }
}
```

**JavaScript/Axios 示例**
```javascript
const getCategoryRecommendations = async (category, limit = 20) => {
  try {
    const response = await axios.get('/api/v1/recommendation/category', {
      params: { category, limit }
    });
    return response.data.data.recommendations;
  } catch (error) {
    console.error('获取分类推荐失败:', error);
    return [];
  }
};
```

---

## 3. 数据结构

### 3.1 RecommendationItem（推荐项）

```typescript
interface RecommendationItem {
  bookId: string;             // 书籍ID
  title: string;              // 标题
  author: string;             // 作者
  cover: string;              // 封面URL
  rating: number;             // 评分（0-5）
  reason?: string;            // 推荐理由
  score?: number;             // 推荐分数（0-1）
  tag?: string;               // 推荐标签
  similarity?: number;        // 相似度（0-1）
  hotScore?: number;          // 热度分数
  trend?: string;             // 趋势（up/down/stable）
}
```

### 3.2 Behavior（用户行为）

```typescript
interface Behavior {
  userId: string;             // 用户ID（自动从Token获取）
  itemId: string;             // 书籍ID
  chapterId?: string;         // 章节ID（可选）
  behaviorType: string;       // 行为类型
  value: number;              // 行为权重
  metadata?: Record<string, any>; // 附加元数据
  timestamp?: string;         // 时间戳（服务器自动生成）
}
```

---

## 4. 完整示例

### 4.1 推荐页面组件（Vue 3）

```vue
<template>
  <div class="recommendation-page">
    <!-- 个性化推荐 -->
    <section class="recommendation-section" v-if="isLoggedIn">
      <h2>为您推荐</h2>
      <div class="books-grid">
        <book-card
          v-for="book in personalizedBooks"
          :key="book.bookId"
          :book="book"
          :show-reason="true"
          @click="goToBookDetail(book.bookId)"
        />
      </div>
      <el-button @click="loadMore('personalized')">加载更多</el-button>
    </section>

    <!-- 热门推荐 -->
    <section class="recommendation-section">
      <h2>本周热门</h2>
      <el-tabs v-model="hotDays" @tab-change="loadHotBooks">
        <el-tab-pane label="7天" name="7"></el-tab-pane>
        <el-tab-pane label="15天" name="15"></el-tab-pane>
        <el-tab-pane label="30天" name="30"></el-tab-pane>
      </el-tabs>
      <div class="books-grid">
        <book-card
          v-for="book in hotBooks"
          :key="book.bookId"
          :book="book"
          :show-trend="true"
          @click="goToBookDetail(book.bookId)"
        />
      </div>
    </section>

    <!-- 分类推荐 -->
    <section class="recommendation-section">
      <h2>分类推荐</h2>
      <el-select v-model="selectedCategory" @change="loadCategoryBooks">
        <el-option
          v-for="cat in categories"
          :key="cat.id"
          :label="cat.name"
          :value="cat.name"
        />
      </el-select>
      <div class="books-grid">
        <book-card
          v-for="book in categoryBooks"
          :key="book.bookId"
          :book="book"
          @click="goToBookDetail(book.bookId)"
        />
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import axios from 'axios';

const router = useRouter();
const userStore = useUserStore();

// 数据
const personalizedBooks = ref([]);
const hotBooks = ref([]);
const categoryBooks = ref([]);
const categories = ref([]);
const selectedCategory = ref('');
const hotDays = ref('7');

// 计算属性
const isLoggedIn = computed(() => userStore.isLoggedIn);

// 加载个性化推荐
const loadPersonalizedBooks = async (limit = 20) => {
  if (!isLoggedIn.value) return;
  
  try {
    const response = await axios.get('/api/v1/recommendation/personalized', {
      params: { limit }
    });
    personalizedBooks.value = response.data.data.recommendations;
  } catch (error) {
    console.error('加载个性化推荐失败:', error);
  }
};

// 加载热门推荐
const loadHotBooks = async () => {
  try {
    const response = await axios.get('/api/v1/recommendation/hot', {
      params: { days: parseInt(hotDays.value), limit: 20 }
    });
    hotBooks.value = response.data.data.recommendations;
  } catch (error) {
    console.error('加载热门推荐失败:', error);
  }
};

// 加载分类推荐
const loadCategoryBooks = async () => {
  if (!selectedCategory.value) return;
  
  try {
    const response = await axios.get('/api/v1/recommendation/category', {
      params: { category: selectedCategory.value, limit: 20 }
    });
    categoryBooks.value = response.data.data.recommendations;
  } catch (error) {
    console.error('加载分类推荐失败:', error);
  }
};

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await axios.get('/api/v1/bookstore/categories/tree');
    categories.value = response.data.data.flatMap(cat => [cat, ...(cat.children || [])]);
    if (categories.value.length > 0) {
      selectedCategory.value = categories.value[0].name;
      await loadCategoryBooks();
    }
  } catch (error) {
    console.error('加载分类列表失败:', error);
  }
};

// 加载更多
const loadMore = async (type) => {
  if (type === 'personalized') {
    await loadPersonalizedBooks(personalizedBooks.value.length + 20);
  }
};

// 跳转到书籍详情
const goToBookDetail = (bookId) => {
  router.push(`/books/${bookId}`);
};

// 生命周期
onMounted(async () => {
  await Promise.all([
    loadPersonalizedBooks(),
    loadHotBooks(),
    loadCategories()
  ]);
});
</script>

<style scoped>
.recommendation-page {
  padding: 20px;
}

.recommendation-section {
  margin-bottom: 40px;
}

.recommendation-section h2 {
  font-size: 24px;
  margin-bottom: 20px;
}

.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}
</style>
```

### 4.2 书籍详情页行为追踪

```vue
<template>
  <div class="book-detail">
    <div class="book-info">
      <img :src="book.cover" :alt="book.title" />
      <div class="book-meta">
        <h1>{{ book.title }}</h1>
        <p>作者：{{ book.author }}</p>
        <el-button type="primary" @click="startReading">开始阅读</el-button>
        <el-button @click="collectBook">{{ isCollected ? '已收藏' : '收藏' }}</el-button>
        <el-button @click="likeBook">{{ isLiked ? '已点赞' : '点赞' }}</el-button>
      </div>
    </div>

    <!-- 相似推荐 -->
    <section class="similar-section">
      <h2>相似推荐</h2>
      <div class="books-grid">
        <book-card
          v-for="similar in similarBooks"
          :key="similar.bookId"
          :book="similar"
          @click="goToBook(similar.bookId)"
        />
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';

const route = useRoute();
const router = useRouter();

// 数据
const book = ref(null);
const similarBooks = ref([]);
const isCollected = ref(false);
const isLiked = ref(false);
const readingStartTime = ref(null);

// 加载书籍详情
const loadBookDetail = async (bookId) => {
  try {
    const response = await axios.get(`/api/v1/bookstore/books/${bookId}`);
    book.value = response.data.data;
    
    // 记录浏览行为
    trackBehavior('view', bookId);
    
    // 加载相似推荐
    await loadSimilarBooks(bookId);
  } catch (error) {
    console.error('加载书籍详情失败:', error);
  }
};

// 加载相似推荐
const loadSimilarBooks = async (bookId) => {
  try {
    const response = await axios.get('/api/v1/recommendation/similar', {
      params: { itemId: bookId, limit: 10 }
    });
    similarBooks.value = response.data.data.similar_items;
  } catch (error) {
    console.error('加载相似推荐失败:', error);
  }
};

// 记录行为
const trackBehavior = async (type, bookId, metadata = {}) => {
  try {
    await axios.post('/api/v1/recommendation/behavior', {
      itemId: bookId,
      behaviorType: type,
      value: getBehaviorWeight(type),
      metadata
    });
  } catch (error) {
    console.warn('记录行为失败:', error);
  }
};

// 获取行为权重
const getBehaviorWeight = (type) => {
  const weights = {
    view: 0.1,
    click: 0.2,
    collect: 0.5,
    read: 0.7,
    like: 0.8
  };
  return weights[type] || 1.0;
};

// 开始阅读
const startReading = () => {
  readingStartTime.value = Date.now();
  trackBehavior('read', book.value.id);
  router.push(`/reader/${book.value.firstChapterId}`);
};

// 收藏书籍
const collectBook = async () => {
  if (isCollected.value) return;
  
  try {
    await axios.post('/api/v1/library/collect', { bookId: book.value.id });
    isCollected.value = true;
    trackBehavior('collect', book.value.id);
  } catch (error) {
    console.error('收藏失败:', error);
  }
};

// 点赞书籍
const likeBook = async () => {
  if (isLiked.value) return;
  
  try {
    await axios.post('/api/v1/books/like', { bookId: book.value.id });
    isLiked.value = true;
    trackBehavior('like', book.value.id);
  } catch (error) {
    console.error('点赞失败:', error);
  }
};

// 跳转到其他书籍
const goToBook = (bookId) => {
  router.push(`/books/${bookId}`);
};

// 页面卸载时记录阅读时长
const recordReadingTime = () => {
  if (readingStartTime.value) {
    const duration = Math.floor((Date.now() - readingStartTime.value) / 1000);
    trackBehavior('view', book.value.id, { viewDuration: duration });
  }
};

// 生命周期
onMounted(() => {
  loadBookDetail(route.params.bookId);
  
  // 监听页面离开
  window.addEventListener('beforeunload', recordReadingTime);
});

// 监听路由变化
watch(() => route.params.bookId, (newId) => {
  if (newId) {
    recordReadingTime();
    readingStartTime.value = null;
    loadBookDetail(newId);
  }
});
</script>
```

---

## 5. 最佳实践

### 5.1 行为追踪批量上传

```javascript
// 批量缓存行为，定时上传
class BehaviorBatcher {
  constructor() {
    this.queue = [];
    this.flushInterval = 5000; // 5秒上传一次
    this.maxQueueSize = 50;    // 最多缓存50条
    
    this.startTimer();
  }

  add(behavior) {
    this.queue.push({
      ...behavior,
      timestamp: new Date().toISOString()
    });
    
    if (this.queue.length >= this.maxQueueSize) {
      this.flush();
    }
  }

  async flush() {
    if (this.queue.length === 0) return;
    
    const behaviors = [...this.queue];
    this.queue = [];
    
    try {
      await axios.post('/api/v1/recommendation/behaviors/batch', {
        behaviors
      });
    } catch (error) {
      console.error('批量上传行为失败:', error);
      // 失败时放回队列
      this.queue.unshift(...behaviors);
    }
  }

  startTimer() {
    setInterval(() => {
      this.flush();
    }, this.flushInterval);
  }
}

const batcher = new BehaviorBatcher();

// 使用
batcher.add({
  itemId: 'book123',
  behaviorType: 'view',
  value: 0.1
});
```

### 5.2 推荐结果缓存

```javascript
// 使用 Vuex/Pinia 缓存推荐结果
import { defineStore } from 'pinia';

export const useRecommendationStore = defineStore('recommendation', {
  state: () => ({
    personalizedCache: null,
    hotCache: {},
    cacheExpiry: 5 * 60 * 1000 // 5分钟过期
  }),
  
  actions: {
    async getPersonalized(limit = 20, force = false) {
      if (!force && this.personalizedCache && this.isCacheValid(this.personalizedCache.timestamp)) {
        return this.personalizedCache.data;
      }
      
      const response = await axios.get('/api/v1/recommendation/personalized', {
        params: { limit }
      });
      
      this.personalizedCache = {
        data: response.data.data.recommendations,
        timestamp: Date.now()
      };
      
      return this.personalizedCache.data;
    },
    
    isCacheValid(timestamp) {
      return Date.now() - timestamp < this.cacheExpiry;
    }
  }
});
```

### 5.3 推荐多样性优化

```javascript
// 前端对推荐结果进行去重和多样性处理
const diversifyRecommendations = (recommendations, seenBookIds = []) => {
  // 过滤已看过的书籍
  let filtered = recommendations.filter(book => !seenBookIds.includes(book.bookId));
  
  // 按分类分组
  const byCategory = {};
  filtered.forEach(book => {
    const cat = book.categoryName || 'other';
    if (!byCategory[cat]) byCategory[cat] = [];
    byCategory[cat].push(book);
  });
  
  // 每个分类最多取3本，保证多样性
  const diversified = [];
  Object.values(byCategory).forEach(books => {
    diversified.push(...books.slice(0, 3));
  });
  
  return diversified;
};
```

---

## 6. 常见问题

### Q1: 如何提高推荐质量？

**A**: 
1. **多记录用户行为**: 尽可能详细地记录用户的各种行为
2. **提供反馈机制**: 让用户可以标记"不感兴趣"
3. **冷启动处理**: 新用户引导选择偏好分类

```javascript
// 新用户偏好选择
const setupUserPreferences = async (categoryIds) => {
  try {
    await axios.post('/api/v1/user/preferences', {
      categories: categoryIds
    });
    
    // 为每个选择的分类记录一次兴趣行为
    for (const catId of categoryIds) {
      await axios.post('/api/v1/recommendation/behavior', {
        itemId: 'category_' + catId,
        behaviorType: 'interest',
        value: 1.0
      });
    }
  } catch (error) {
    console.error('设置偏好失败:', error);
  }
};
```

### Q2: 推荐结果为空怎么办？

**A**: 提供降级策略：

```javascript
const getRecommendations = async (userId) => {
  try {
    // 先尝试个性化推荐
    let recs = await axios.get('/api/v1/recommendation/personalized');
    if (recs.data.data.count > 0) {
      return recs.data.data.recommendations;
    }
    
    // 降级到热门推荐
    recs = await axios.get('/api/v1/recommendation/hot');
    if (recs.data.data.count > 0) {
      return recs.data.data.recommendations;
    }
    
    // 最后降级到编辑推荐
    recs = await axios.get('/api/v1/bookstore/books/featured');
    return recs.data.data;
  } catch (error) {
    console.error('获取推荐失败:', error);
    return [];
  }
};
```

### Q3: 如何避免过度追踪用户行为？

**A**: 实施隐私保护策略：

```javascript
// 只在用户同意的情况下追踪
const trackBehaviorWithConsent = async (behavior) => {
  // 检查用户隐私设置
  const hasConsent = localStorage.getItem('trackingConsent') === 'true';
  if (!hasConsent) return;
  
  try {
    await axios.post('/api/v1/recommendation/behavior', behavior);
  } catch (error) {
    console.warn('记录行为失败:', error);
  }
};

// 提供隐私设置
const updatePrivacySettings = async (allowTracking) => {
  localStorage.setItem('trackingConsent', allowTracking);
  
  // 如果关闭追踪，删除服务器端的行为数据
  if (!allowTracking) {
    await axios.delete('/api/v1/recommendation/behaviors');
  }
};
```

---

**文档版本**: v1.0  
**最后更新**: 2025-10-18  
**维护者**: 青羽后端团队

