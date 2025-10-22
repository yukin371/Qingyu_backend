# 书城系统 API 参考

> **版本**: v1.0  
> **最后更新**: 2025-10-18  
> **基础路径**: `/api/v1/bookstore`

---

## 1. 概述

书城系统提供书籍浏览、搜索、分类管理、榜单展示等功能，是平台的核心阅读内容展示模块。

### 1.1 基础信息

- **认证要求**: 大部分接口无需认证，部分统计接口需要登录
- **响应格式**: 统一 JSON 格式
- **分页支持**: 是
- **搜索支持**: 是

### 1.2 功能特性

- ✅ 首页数据聚合（Banner + 推荐 + 精选）
- ✅ 书籍详情查看
- ✅ 分类树结构
- ✅ 高级搜索和过滤
- ✅ 多维度榜单（实时/周榜/月榜/新人榜）
- ✅ 浏览量和点击量统计

---

## 2. 接口列表

### 2.1 首页相关

#### 2.1.1 获取首页数据

**接口说明**: 一次性获取首页所需的所有数据（Banner、推荐书籍、精选书籍、热门分类）

**请求**
```
GET /api/v1/bookstore/homepage
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取首页数据成功",
  "data": {
    "banners": [
      {
        "id": "banner123",
        "title": "新书推荐",
        "image": "https://example.com/banner.jpg",
        "link": "/books/book123",
        "sort": 1
      }
    ],
    "recommendedBooks": [...],
    "featuredBooks": [...],
    "categories": [...]
  }
}
```

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/homepage"
```

**JavaScript/Axios 示例**
```javascript
const getHomepage = async () => {
  try {
    const response = await axios.get('/api/v1/bookstore/homepage');
    return response.data.data;
  } catch (error) {
    console.error('获取首页数据失败:', error);
    throw error;
  }
};
```

---

### 2.2 书籍相关

#### 2.2.1 获取书籍详情

**接口说明**: 根据书籍ID获取书籍的详细信息

**请求**
```
GET /api/v1/bookstore/books/{id}
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 书籍ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取书籍详情成功",
  "data": {
    "id": "book123",
    "title": "示例书籍",
    "author": "作者名",
    "cover": "https://example.com/cover.jpg",
    "description": "书籍简介",
    "categoryId": "cat123",
    "categoryName": "玄幻",
    "status": "ongoing",
    "wordCount": 1000000,
    "chapterCount": 100,
    "rating": 4.5,
    "viewCount": 10000,
    "likeCount": 500,
    "tags": ["热血", "冒险"],
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-10-18T00:00:00Z"
  }
}
```

**错误响应**
- `400` - 书籍ID不能为空
- `404` - 书籍不存在或不可用
- `500` - 服务器错误

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books/book123"
```

**JavaScript/Axios 示例**
```javascript
const getBookDetail = async (bookId) => {
  try {
    const response = await axios.get(`/api/v1/bookstore/books/${bookId}`);
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 404) {
      console.error('书籍不存在');
    }
    throw error;
  }
};
```

#### 2.2.2 根据分类获取书籍列表

**接口说明**: 获取指定分类下的书籍列表，支持分页

**请求**
```
GET /api/v1/bookstore/categories/{categoryId}/books
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| categoryId | string | 是 | 分类ID |

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| size | int | 否 | 20 | 每页数量（最大100） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取分类书籍成功",
  "data": [
    {
      "id": "book123",
      "title": "示例书籍",
      "author": "作者名",
      "cover": "https://example.com/cover.jpg",
      "rating": 4.5,
      "viewCount": 10000
    }
  ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/categories/cat123/books?page=1&size=20"
```

#### 2.2.3 获取推荐书籍

**接口说明**: 获取平台推荐的书籍列表

**请求**
```
GET /api/v1/bookstore/books/recommended
```

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| size | int | 否 | 20 | 每页数量 |

#### 2.2.4 获取精选书籍

**接口说明**: 获取编辑精选的书籍列表

**请求**
```
GET /api/v1/bookstore/books/featured
```

**Query 参数**: 同推荐书籍

#### 2.2.5 搜索书籍

**接口说明**: 根据关键词和过滤条件搜索书籍

**请求**
```
GET /api/v1/bookstore/books/search
```

**Query 参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词 |
| categoryId | string | 否 | 分类ID |
| author | string | 否 | 作者名 |
| minRating | number | 否 | 最低评分 |
| tags | []string | 否 | 标签列表 |
| sortBy | string | 否 | 排序字段（created_at/updated_at/view_count/like_count/rating） |
| sortOrder | string | 否 | 排序方向（asc/desc），默认 desc |
| page | int | 否 | 页码，默认1 |
| size | int | 否 | 每页数量，默认20 |

**请求注意事项**
- 至少需要提供 `keyword`、`categoryId` 或 `author` 中的一个参数
- `tags` 可以传递多个值

**响应示例**
```json
{
  "code": 200,
  "message": "搜索书籍成功",
  "data": [
    {
      "id": "book123",
      "title": "示例书籍",
      "author": "作者名",
      "cover": "https://example.com/cover.jpg",
      "rating": 4.5
    }
  ],
  "total": 50,
  "page": 1,
  "size": 20
}
```

**cURL 示例**
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books/search?keyword=玄幻&sortBy=view_count&sortOrder=desc&page=1&size=20"
```

**JavaScript/Axios 示例**
```javascript
const searchBooks = async (params) => {
  try {
    const response = await axios.get('/api/v1/bookstore/books/search', {
      params: {
        keyword: params.keyword,
        categoryId: params.categoryId,
        author: params.author,
        tags: params.tags, // 数组会自动转换为 tags[]=tag1&tags[]=tag2
        sortBy: params.sortBy || 'created_at',
        sortOrder: params.sortOrder || 'desc',
        page: params.page || 1,
        size: params.size || 20
      }
    });
    return response.data;
  } catch (error) {
    console.error('搜索书籍失败:', error);
    throw error;
  }
};

// 使用示例
const results = await searchBooks({
  keyword: '玄幻',
  tags: ['热血', '冒险'],
  sortBy: 'view_count',
  page: 1
});
```

#### 2.2.6 增加书籍浏览量

**接口说明**: 记录用户浏览书籍，增加浏览量统计

**请求**
```
POST /api/v1/bookstore/books/{id}/view
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 书籍ID |

**响应示例**
```json
{
  "code": 200,
  "message": "浏览量增加成功"
}
```

**JavaScript/Axios 示例**
```javascript
const incrementBookView = async (bookId) => {
  try {
    await axios.post(`/api/v1/bookstore/books/${bookId}/view`);
  } catch (error) {
    console.error('增加浏览量失败:', error);
  }
};
```

---

### 2.3 分类相关

#### 2.3.1 获取分类树

**接口说明**: 获取完整的分类树结构（包含父子关系）

**请求**
```
GET /api/v1/bookstore/categories/tree
```

**响应示例**
```json
{
  "code": 200,
  "message": "获取分类树成功",
  "data": [
    {
      "id": "cat1",
      "name": "玄幻",
      "description": "玄幻类小说",
      "icon": "https://example.com/icon.png",
      "sort": 1,
      "children": [
        {
          "id": "cat1-1",
          "name": "东方玄幻",
          "parentId": "cat1",
          "sort": 1
        },
        {
          "id": "cat1-2",
          "name": "异世大陆",
          "parentId": "cat1",
          "sort": 2
        }
      ]
    },
    {
      "id": "cat2",
      "name": "都市",
      "description": "都市类小说",
      "sort": 2,
      "children": []
    }
  ]
}
```

**JavaScript/Axios 示例**
```javascript
const getCategoryTree = async () => {
  try {
    const response = await axios.get('/api/v1/bookstore/categories/tree');
    return response.data.data;
  } catch (error) {
    console.error('获取分类树失败:', error);
    throw error;
  }
};
```

#### 2.3.2 获取分类详情

**接口说明**: 根据分类ID获取分类的详细信息

**请求**
```
GET /api/v1/bookstore/categories/{id}
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 分类ID |

**响应示例**
```json
{
  "code": 200,
  "message": "获取分类详情成功",
  "data": {
    "id": "cat1",
    "name": "玄幻",
    "description": "玄幻类小说",
    "icon": "https://example.com/icon.png",
    "parentId": null,
    "sort": 1,
    "bookCount": 1000
  }
}
```

---

### 2.4 Banner 相关

#### 2.4.1 获取激活的 Banner 列表

**接口说明**: 获取当前激活的 Banner 列表

**请求**
```
GET /api/v1/bookstore/banners
```

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| limit | int | 否 | 5 | 数量限制（最大20） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取Banner列表成功",
  "data": [
    {
      "id": "banner1",
      "title": "新书推荐",
      "image": "https://example.com/banner.jpg",
      "link": "/books/book123",
      "sort": 1,
      "startTime": "2024-01-01T00:00:00Z",
      "endTime": "2024-12-31T23:59:59Z"
    }
  ]
}
```

#### 2.4.2 增加 Banner 点击次数

**接口说明**: 记录用户点击 Banner，增加点击次数统计

**请求**
```
POST /api/v1/bookstore/banners/{id}/click
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | Banner ID |

**响应示例**
```json
{
  "code": 200,
  "message": "点击次数增加成功"
}
```

**JavaScript/Axios 示例**
```javascript
const trackBannerClick = async (bannerId) => {
  try {
    await axios.post(`/api/v1/bookstore/banners/${bannerId}/click`);
  } catch (error) {
    console.error('记录Banner点击失败:', error);
  }
};
```

---

### 2.5 榜单相关

#### 2.5.1 获取实时榜

**接口说明**: 获取当日实时榜单数据

**请求**
```
GET /api/v1/bookstore/rankings/realtime
```

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| limit | int | 否 | 20 | 限制数量（最大100） |

**响应示例**
```json
{
  "code": 200,
  "message": "获取实时榜成功",
  "data": [
    {
      "rank": 1,
      "bookId": "book123",
      "title": "示例书籍",
      "author": "作者名",
      "cover": "https://example.com/cover.jpg",
      "score": 9500,
      "change": 2
    }
  ]
}
```

#### 2.5.2 获取周榜

**接口说明**: 获取指定周期的周榜单数据

**请求**
```
GET /api/v1/bookstore/rankings/weekly
```

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| period | string | 否 | 当前周 | 周期（格式: 2024-W01） |
| limit | int | 否 | 20 | 限制数量 |

**响应格式**: 同实时榜

#### 2.5.3 获取月榜

**接口说明**: 获取指定月份的月榜单数据

**请求**
```
GET /api/v1/bookstore/rankings/monthly
```

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| period | string | 否 | 当前月 | 月份（格式: 2024-01） |
| limit | int | 否 | 20 | 限制数量 |

#### 2.5.4 获取新人榜

**接口说明**: 获取指定月份的新人榜单数据

**请求**
```
GET /api/v1/bookstore/rankings/newbie
```

**Query 参数**: 同月榜

#### 2.5.5 根据类型获取榜单（统一接口）

**接口说明**: 根据榜单类型获取指定周期的榜单数据

**请求**
```
GET /api/v1/bookstore/rankings/{type}
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 榜单类型（realtime/weekly/monthly/newbie） |

**Query 参数**
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| period | string | 否 | - | 周期 |
| limit | int | 否 | 20 | 限制数量 |

**JavaScript/Axios 示例**
```javascript
const getRanking = async (type, period = null, limit = 20) => {
  try {
    const response = await axios.get(`/api/v1/bookstore/rankings/${type}`, {
      params: { period, limit }
    });
    return response.data.data;
  } catch (error) {
    console.error('获取榜单失败:', error);
    throw error;
  }
};

// 使用示例
const realtimeRanking = await getRanking('realtime');
const weeklyRanking = await getRanking('weekly', '2024-W42');
const monthlyRanking = await getRanking('monthly', '2024-10');
```

---

## 3. 数据结构

### 3.1 Book（书籍）

```typescript
interface Book {
  id: string;                 // 书籍ID
  title: string;              // 标题
  author: string;             // 作者
  cover: string;              // 封面URL
  description: string;        // 简介
  categoryId: string;         // 分类ID
  categoryName: string;       // 分类名称
  status: string;             // 状态（ongoing/completed/paused）
  wordCount: number;          // 字数
  chapterCount: number;       // 章节数
  rating: number;             // 评分（0-5）
  viewCount: number;          // 浏览量
  likeCount: number;          // 点赞数
  tags: string[];             // 标签
  createdAt: string;          // 创建时间（ISO 8601）
  updatedAt: string;          // 更新时间（ISO 8601）
}
```

### 3.2 Category（分类）

```typescript
interface Category {
  id: string;                 // 分类ID
  name: string;               // 分类名称
  description: string;        // 分类描述
  icon: string;               // 图标URL
  parentId: string | null;    // 父分类ID
  sort: number;               // 排序值
  bookCount: number;          // 书籍数量
  children?: Category[];      // 子分类（仅在树形结构中）
}
```

### 3.3 Banner

```typescript
interface Banner {
  id: string;                 // Banner ID
  title: string;              // 标题
  image: string;              // 图片URL
  link: string;               // 跳转链接
  sort: number;               // 排序值
  startTime: string;          // 开始时间
  endTime: string;            // 结束时间
  clickCount: number;         // 点击次数
}
```

### 3.4 RankingItem（榜单项）

```typescript
interface RankingItem {
  rank: number;               // 排名
  bookId: string;             // 书籍ID
  title: string;              // 书籍标题
  author: string;             // 作者
  cover: string;              // 封面URL
  score: number;              // 榜单分数
  change: number;             // 排名变化（正数上升，负数下降）
}
```

---

## 4. 错误码说明

### 4.1 通用错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 参数错误（ID为空、格式无效等） |
| 404 | 资源不存在（书籍、分类、Banner等） |
| 500 | 服务器内部错误 |

### 4.2 业务错误码

| 错误信息 | 说明 |
|----------|------|
| "书籍ID不能为空" | 未提供书籍ID |
| "分类ID不能为空" | 未提供分类ID |
| "分类ID格式无效" | 分类ID格式不正确 |
| "请提供搜索关键词或过滤条件" | 搜索接口缺少必要参数 |
| "无效的榜单类型" | 榜单类型不在允许范围内 |
| "书籍不存在或不可用" | 书籍已删除或未上架 |

---

## 5. 完整示例

### 5.1 书城首页组件（Vue 3）

```vue
<template>
  <div class="bookstore-home">
    <!-- Banner 轮播 -->
    <el-carousel height="300px" v-if="homepage.banners">
      <el-carousel-item v-for="banner in homepage.banners" :key="banner.id">
        <img :src="banner.image" :alt="banner.title" @click="onBannerClick(banner)" />
      </el-carousel-item>
    </el-carousel>

    <!-- 推荐书籍 -->
    <section class="books-section">
      <h2>编辑推荐</h2>
      <div class="books-grid">
        <book-card
          v-for="book in homepage.recommendedBooks"
          :key="book.id"
          :book="book"
          @click="goToBookDetail(book.id)"
        />
      </div>
    </section>

    <!-- 精选书籍 -->
    <section class="books-section">
      <h2>精选书籍</h2>
      <div class="books-grid">
        <book-card
          v-for="book in homepage.featuredBooks"
          :key="book.id"
          :book="book"
        />
      </div>
    </section>

    <!-- 分类导航 -->
    <section class="categories">
      <h2>热门分类</h2>
      <div class="category-list">
        <div
          v-for="category in homepage.categories"
          :key="category.id"
          class="category-item"
          @click="goToCategory(category.id)"
        >
          {{ category.name }}
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const homepage = ref({});
const loading = ref(false);

// 获取首页数据
const fetchHomepage = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/v1/bookstore/homepage');
    homepage.value = response.data.data;
  } catch (error) {
    console.error('获取首页数据失败:', error);
  } finally {
    loading.value = false;
  }
};

// Banner 点击处理
const onBannerClick = async (banner) => {
  // 记录点击
  await axios.post(`/api/v1/bookstore/banners/${banner.id}/click`).catch(() => {});
  
  // 跳转
  if (banner.link.startsWith('/')) {
    router.push(banner.link);
  } else {
    window.open(banner.link, '_blank');
  }
};

// 跳转到书籍详情
const goToBookDetail = (bookId) => {
  router.push(`/books/${bookId}`);
};

// 跳转到分类页面
const goToCategory = (categoryId) => {
  router.push(`/categories/${categoryId}`);
};

onMounted(() => {
  fetchHomepage();
});
</script>
```

### 5.2 搜索页面组件

```vue
<template>
  <div class="search-page">
    <!-- 搜索表单 -->
    <el-form :model="searchForm" inline>
      <el-form-item label="关键词">
        <el-input v-model="searchForm.keyword" placeholder="输入书名或作者" />
      </el-form-item>
      <el-form-item label="分类">
        <el-select v-model="searchForm.categoryId" placeholder="选择分类">
          <el-option
            v-for="cat in categories"
            :key="cat.id"
            :label="cat.name"
            :value="cat.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="排序">
        <el-select v-model="searchForm.sortBy">
          <el-option label="最新更新" value="updated_at" />
          <el-option label="最高人气" value="view_count" />
          <el-option label="最高评分" value="rating" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="search">搜索</el-button>
      </el-form-item>
    </el-form>

    <!-- 搜索结果 -->
    <div class="search-results">
      <el-empty v-if="!loading && results.length === 0" description="暂无搜索结果" />
      <div v-else class="books-grid">
        <book-card
          v-for="book in results"
          :key="book.id"
          :book="book"
          @click="goToBookDetail(book.id)"
        />
      </div>

      <!-- 分页 -->
      <el-pagination
        v-if="total > 0"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="total"
        layout="total, prev, pager, next, jumper"
        @current-change="search"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const loading = ref(false);
const results = ref([]);
const total = ref(0);
const categories = ref([]);

const searchForm = reactive({
  keyword: '',
  categoryId: '',
  sortBy: 'updated_at'
});

const pagination = reactive({
  page: 1,
  size: 20
});

// 搜索书籍
const search = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/v1/bookstore/books/search', {
      params: {
        keyword: searchForm.keyword,
        categoryId: searchForm.categoryId,
        sortBy: searchForm.sortBy,
        sortOrder: 'desc',
        page: pagination.page,
        size: pagination.size
      }
    });
    
    results.value = response.data.data;
    total.value = response.data.total;
  } catch (error) {
    console.error('搜索失败:', error);
  } finally {
    loading.value = false;
  }
};

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await axios.get('/api/v1/bookstore/categories/tree');
    categories.value = response.data.data.flatMap(cat => [
      cat,
      ...(cat.children || [])
    ]);
  } catch (error) {
    console.error('获取分类失败:', error);
  }
};

const goToBookDetail = (bookId) => {
  router.push(`/books/${bookId}`);
};

onMounted(() => {
  fetchCategories();
});
</script>
```

---

## 6. 最佳实践

### 6.1 首页数据加载优化

```javascript
// 使用 Promise.all 并行加载多个数据（如果不使用 homepage 聚合接口）
const loadHomePageData = async () => {
  try {
    const [banners, recommended, featured, categories] = await Promise.all([
      axios.get('/api/v1/bookstore/banners?limit=5'),
      axios.get('/api/v1/bookstore/books/recommended?size=10'),
      axios.get('/api/v1/bookstore/books/featured?size=10'),
      axios.get('/api/v1/bookstore/categories/tree')
    ]);
    
    return {
      banners: banners.data.data,
      recommendedBooks: recommended.data.data,
      featuredBooks: featured.data.data,
      categories: categories.data.data
    };
  } catch (error) {
    console.error('加载首页数据失败:', error);
    throw error;
  }
};
```

### 6.2 浏览量统计（防抖）

```javascript
import { debounce } from 'lodash-es';

// 防抖处理，避免短时间内重复请求
const trackBookView = debounce(async (bookId) => {
  try {
    await axios.post(`/api/v1/bookstore/books/${bookId}/view`);
  } catch (error) {
    // 静默失败，不影响用户体验
    console.warn('记录浏览量失败:', error);
  }
}, 3000, { leading: true, trailing: false });

// 在书籍详情页挂载时调用
onMounted(() => {
  trackBookView(bookId.value);
});
```

### 6.3 搜索防抖

```javascript
import { ref, watch } from 'vue';
import { debounce } from 'lodash-es';

const keyword = ref('');

// 防抖搜索
const debouncedSearch = debounce(async (value) => {
  if (!value) return;
  
  try {
    const response = await axios.get('/api/v1/bookstore/books/search', {
      params: { keyword: value, page: 1, size: 20 }
    });
    // 处理搜索结果
  } catch (error) {
    console.error('搜索失败:', error);
  }
}, 500);

// 监听关键词变化
watch(keyword, (newValue) => {
  debouncedSearch(newValue);
});
```

---

## 7. 常见问题

### Q1: 如何实现无限滚动加载？

**A**: 使用分页参数配合滚动监听：

```javascript
import { ref, onMounted, onUnmounted } from 'vue';

const books = ref([]);
const page = ref(1);
const hasMore = ref(true);
const loading = ref(false);

const loadMore = async () => {
  if (loading.value || !hasMore.value) return;
  
  loading.value = true;
  try {
    const response = await axios.get('/api/v1/bookstore/books/recommended', {
      params: { page: page.value, size: 20 }
    });
    
    const newBooks = response.data.data;
    books.value.push(...newBooks);
    
    // 判断是否还有更多数据
    hasMore.value = newBooks.length === 20;
    page.value++;
  } catch (error) {
    console.error('加载更多失败:', error);
  } finally {
    loading.value = false;
  }
};

// 滚动监听
const handleScroll = () => {
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
  const windowHeight = window.innerHeight;
  const documentHeight = document.documentElement.scrollHeight;
  
  if (scrollTop + windowHeight >= documentHeight - 100) {
    loadMore();
  }
};

onMounted(() => {
  loadMore();
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
```

### Q2: 如何缓存分类树数据？

**A**: 使用 Pinia/Vuex 或 localStorage：

```javascript
// stores/category.js
import { defineStore } from 'pinia';
import axios from 'axios';

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categories: [],
    lastFetch: null
  }),
  
  actions: {
    async fetchCategories(force = false) {
      // 如果数据存在且不强制刷新，且缓存未过期（1小时）
      if (
        !force &&
        this.categories.length > 0 &&
        this.lastFetch &&
        Date.now() - this.lastFetch < 3600000
      ) {
        return this.categories;
      }
      
      try {
        const response = await axios.get('/api/v1/bookstore/categories/tree');
        this.categories = response.data.data;
        this.lastFetch = Date.now();
        return this.categories;
      } catch (error) {
        console.error('获取分类失败:', error);
        throw error;
      }
    }
  }
});
```

### Q3: 如何优化书籍列表的加载性能？

**A**: 
1. **图片懒加载**: 使用 `v-lazy` 或 Intersection Observer
2. **虚拟滚动**: 对于超长列表使用 `vue-virtual-scroller`
3. **分页加载**: 避免一次加载过多数据
4. **预加载**: 提前加载下一页数据

```javascript
// 预加载下一页
const prefetchNextPage = async () => {
  if (pagination.page < totalPages.value) {
    try {
      await axios.get('/api/v1/bookstore/books/recommended', {
        params: { page: pagination.page + 1, size: 20 }
      });
      // 数据会被浏览器缓存
    } catch (error) {
      // 忽略错误
    }
  }
};

// 在当前页加载完成后，延迟预加载下一页
setTimeout(() => {
  prefetchNextPage();
}, 1000);
```

---

**文档版本**: v1.0  
**最后更新**: 2025-10-18  
**维护者**: 青羽后端团队

