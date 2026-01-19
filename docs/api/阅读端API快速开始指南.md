# 阅读端API快速开始指南

> **版本**: v1.0  
> **最后更新**: 2025-10-16  
> **适用对象**: 前端开发者、移动端开发者

---

## 📚 简介

本指南提供了阅读端API的快速入门教程和常见使用场景的代码示例，帮助开发者快速集成青羽阅读端功能。

---

## 🚀 快速开始

### 环境准备

**基础URL**:
- 开发环境: `http://localhost:8080`
- 生产环境: `https://api.qingyu.com`

**必要条件**:
- 有效的API访问权限
- （可选）注册账号并获取JWT Token

---

## 📖 常见使用场景

### 场景1：首页展示

#### 需求
展示书城首页，包括Banner、推荐书籍、分类等。

#### 实现步骤

**1. 获取首页数据**

```javascript
async function getHomepageData() {
  const response = await fetch('http://localhost:8080/api/v1/bookstore/homepage');
  const data = await response.json();
  
  if (data.code === 200) {
    return data.data;
  }
  
  throw new Error(data.message);
}

// 使用示例
const homepageData = await getHomepageData();
console.log('Banners:', homepageData.banners);
console.log('推荐书籍:', homepageData.recommendedBooks);
```

**2. 渲染示例（React）**

```jsx
function Homepage() {
  const [data, setData] = useState(null);
  
  useEffect(() => {
    getHomepageData().then(setData);
  }, []);
  
  if (!data) return <Loading />;
  
  return (
    <div>
      {/* Banner轮播 */}
      <BannerCarousel banners={data.banners} />
      
      {/* 推荐书籍 */}
      <BookList 
        title="编辑推荐" 
        books={data.recommendedBooks} 
      />
      
      {/* 分类导航 */}
      <CategoryNav categories={data.categories} />
    </div>
  );
}
```

---

### 场景2：书籍列表和搜索

#### 需求
显示分类书籍列表，支持分页和搜索。

#### 实现步骤

**1. 获取分类书籍列表**

```javascript
async function getBooksByCategory(category, page = 1, size = 20) {
  const url = `http://localhost:8080/api/v1/bookstore/books/category/${encodeURIComponent(category)}`;
  const params = new URLSearchParams({ page, size });
  
  const response = await fetch(`${url}?${params}`);
  const data = await response.json();
  
  return data;
}

// 使用示例
const result = await getBooksByCategory('玄幻', 1, 20);
console.log('书籍列表:', result.data);
console.log('总数:', result.total);
```

**2. 书籍搜索**

```javascript
async function searchBooks(keyword, page = 1) {
  const params = new URLSearchParams({
    q: keyword,
    type: 'all',
    page: page,
    size: 20
  });
  
  const response = await fetch(
    `http://localhost:8080/api/v1/books/search?${params}`
  );
  
  return await response.json();
}

// 使用示例
const searchResult = await searchBooks('斗破');
```

**3. 完整的列表组件（React）**

```jsx
function BookListPage({ category }) {
  const [books, setBooks] = useState([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  
  // 加载书籍
  const loadBooks = async () => {
    setLoading(true);
    try {
      const result = await getBooksByCategory(category, page, 20);
      setBooks(result.data);
      setTotal(result.total);
    } catch (error) {
      console.error('加载失败', error);
    } finally {
      setLoading(false);
    }
  };
  
  useEffect(() => {
    loadBooks();
  }, [category, page]);
  
  return (
    <div>
      <h1>{category}分类</h1>
      
      {loading ? (
        <Loading />
      ) : (
        <>
          <div className="book-grid">
            {books.map(book => (
              <BookCard key={book.id} book={book} />
            ))}
          </div>
          
          <Pagination 
            current={page}
            total={total}
            pageSize={20}
            onChange={setPage}
          />
        </>
      )}
    </div>
  );
}
```

---

### 场景3：书籍详情页

#### 需求
展示书籍详细信息、章节列表、评分等。

#### 实现步骤

**1. 获取书籍详情**

```javascript
async function getBookDetail(bookId) {
  const response = await fetch(
    `http://localhost:8080/api/v1/bookstore/books/${bookId}`
  );
  return await response.json();
}
```

**2. 获取章节列表**

```javascript
async function getChapterList(bookId, page = 1, size = 50) {
  const params = new URLSearchParams({ page, size });
  const response = await fetch(
    `http://localhost:8080/api/v1/books/${bookId}/chapters?${params}`
  );
  return await response.json();
}
```

**3. 获取评分信息**

```javascript
async function getBookRatings(bookId, page = 1) {
  const params = new URLSearchParams({ page, limit: 10 });
  const response = await fetch(
    `http://localhost:8080/api/v1/reading/books/${bookId}/ratings?${params}`
  );
  return await response.json();
}

async function getBookRatingStats(bookId) {
  const response = await fetch(
    `http://localhost:8080/api/v1/reading/books/${bookId}/ratings/stats`
  );
  return await response.json();
}
```

**4. 完整的详情页组件**

```jsx
function BookDetailPage({ bookId }) {
  const [book, setBook] = useState(null);
  const [chapters, setChapters] = useState([]);
  const [ratingStats, setRatingStats] = useState(null);
  
  useEffect(() => {
    // 并行加载所有数据
    Promise.all([
      getBookDetail(bookId),
      getChapterList(bookId, 1, 50),
      getBookRatingStats(bookId)
    ]).then(([bookData, chaptersData, statsData]) => {
      setBook(bookData.data);
      setChapters(chaptersData.data);
      setRatingStats(statsData.data);
    });
  }, [bookId]);
  
  if (!book) return <Loading />;
  
  return (
    <div>
      {/* 书籍基本信息 */}
      <BookInfo book={book} ratingStats={ratingStats} />
      
      {/* 操作按钮 */}
      <ActionButtons bookId={bookId} />
      
      {/* 章节列表 */}
      <ChapterList chapters={chapters} />
      
      {/* 评分和评论 */}
      <RatingSection bookId={bookId} stats={ratingStats} />
    </div>
  );
}
```

---

### 场景4：阅读器实现

#### 需求
实现完整的阅读器功能，包括章节内容、进度保存、翻页等。

#### 实现步骤

**1. 获取章节内容（需要登录）**

```javascript
async function getChapterContent(chapterId, token) {
  const response = await fetch(
    `http://localhost:8080/api/v1/reader/chapters/${chapterId}/content`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  return await response.json();
}
```

**2. 保存阅读进度**

```javascript
async function saveReadingProgress(bookId, chapterId, progress, token) {
  const response = await fetch(
    'http://localhost:8080/api/v1/reader/progress',
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        bookId,
        chapterId,
        progress
      })
    }
  );
  return await response.json();
}
```

**3. 获取章节导航**

```javascript
async function getChapterNavigation(bookId, chapterNum) {
  const params = new URLSearchParams({ bookId, chapterNum });
  const response = await fetch(
    `http://localhost:8080/api/v1/reader/chapters/navigation?${params}`
  );
  return await response.json();
}
```

**4. 完整的阅读器组件**

```jsx
function ReaderPage({ bookId, chapterId }) {
  const [content, setContent] = useState('');
  const [chapter, setChapter] = useState(null);
  const [navigation, setNavigation] = useState(null);
  const [progress, setProgress] = useState(0);
  const token = useAuth(); // 获取用户Token
  
  // 加载章节
  useEffect(() => {
    loadChapter();
  }, [chapterId]);
  
  const loadChapter = async () => {
    try {
      // 获取章节信息
      const chapterInfo = await fetch(
        `http://localhost:8080/api/v1/reader/chapters/${chapterId}`
      ).then(r => r.json());
      
      setChapter(chapterInfo.data);
      
      // 获取章节内容
      const contentData = await getChapterContent(chapterId, token);
      setContent(contentData.data.content);
      
      // 获取导航
      const navData = await getChapterNavigation(
        bookId, 
        chapterInfo.data.chapterNumber
      );
      setNavigation(navData.data);
    } catch (error) {
      console.error('加载章节失败', error);
    }
  };
  
  // 滚动监听 - 自动保存进度
  useEffect(() => {
    const handleScroll = debounce(() => {
      const scrollHeight = document.documentElement.scrollHeight;
      const scrollTop = window.scrollY;
      const clientHeight = window.innerHeight;
      
      const currentProgress = (scrollTop + clientHeight) / scrollHeight;
      setProgress(currentProgress);
      
      // 自动保存进度
      saveReadingProgress(bookId, chapterId, currentProgress, token);
    }, 1000);
    
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [bookId, chapterId, token]);
  
  return (
    <div className="reader">
      {/* 顶部导航栏 */}
      <ReaderHeader 
        title={chapter?.title}
        progress={progress}
      />
      
      {/* 章节内容 */}
      <div className="reader-content">
        <h2>{chapter?.title}</h2>
        <div 
          className="content"
          dangerouslySetInnerHTML={{ __html: content }}
        />
      </div>
      
      {/* 底部翻页 */}
      <ReaderFooter
        prevChapter={navigation?.prevChapter}
        nextChapter={navigation?.nextChapter}
        onNavigate={loadChapter}
      />
    </div>
  );
}
```

---

### 场景5：注记功能

#### 需求
实现高亮标注、书签、笔记等注记功能。

#### 实现步骤

**1. 创建注记**

```javascript
async function createAnnotation(annotation, token) {
  const response = await fetch(
    'http://localhost:8080/api/v1/reader/annotations',
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(annotation)
    }
  );
  return await response.json();
}

// 使用示例 - 创建高亮
const highlight = await createAnnotation({
  bookId: 'book123',
  chapterId: 'chapter50',
  type: 'highlight',
  text: '选中的文本',
  range: '100-150'
}, token);

// 创建书签
const bookmark = await createAnnotation({
  bookId: 'book123',
  chapterId: 'chapter50',
  type: 'bookmark',
  text: '第50章重要位置'
}, token);

// 创建笔记
const note = await createAnnotation({
  bookId: 'book123',
  chapterId: 'chapter50',
  type: 'note',
  text: '选中的文本',
  note: '这里是我的笔记内容',
  range: '200-250'
}, token);
```

**2. 获取注记列表**

```javascript
async function getBookAnnotations(bookId, type, token) {
  const params = new URLSearchParams({ type });
  const response = await fetch(
    `http://localhost:8080/api/v1/reader/annotations/book/${bookId}?${params}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  return await response.json();
}

// 获取所有高亮
const highlights = await getBookAnnotations('book123', 'highlight', token);

// 获取所有书签
const bookmarks = await getBookAnnotations('book123', 'bookmark', token);
```

**3. 文本选择和标注组件**

```jsx
function TextSelection({ bookId, chapterId }) {
  const [selectedText, setSelectedText] = useState('');
  const [range, setRange] = useState('');
  const [showMenu, setShowMenu] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ x: 0, y: 0 });
  const token = useAuth();
  
  // 处理文本选择
  const handleTextSelect = () => {
    const selection = window.getSelection();
    const text = selection.toString().trim();
    
    if (text) {
      setSelectedText(text);
      
      // 计算选择范围（简化版）
      const range = {
        start: selection.anchorOffset,
        end: selection.focusOffset
      };
      setRange(`${range.start}-${range.end}`);
      
      // 显示操作菜单
      const rect = selection.getRangeAt(0).getBoundingClientRect();
      setMenuPosition({
        x: rect.left + rect.width / 2,
        y: rect.top - 40
      });
      setShowMenu(true);
    } else {
      setShowMenu(false);
    }
  };
  
  // 创建高亮
  const handleHighlight = async () => {
    await createAnnotation({
      bookId,
      chapterId,
      type: 'highlight',
      text: selectedText,
      range
    }, token);
    
    setShowMenu(false);
    // 重新渲染内容以显示高亮
  };
  
  // 创建笔记
  const handleNote = async () => {
    const note = prompt('请输入笔记内容：');
    if (note) {
      await createAnnotation({
        bookId,
        chapterId,
        type: 'note',
        text: selectedText,
        note,
        range
      }, token);
      
      setShowMenu(false);
    }
  };
  
  return (
    <div onMouseUp={handleTextSelect}>
      {/* 内容区域 */}
      <div className="content">
        {/* 渲染章节内容 */}
      </div>
      
      {/* 选择菜单 */}
      {showMenu && (
        <div 
          className="selection-menu"
          style={{
            position: 'fixed',
            left: menuPosition.x,
            top: menuPosition.y
          }}
        >
          <button onClick={handleHighlight}>高亮</button>
          <button onClick={handleNote}>笔记</button>
        </div>
      )}
    </div>
  );
}
```

---

### 场景6：阅读设置

#### 需求
实现阅读器个性化设置，如字体、主题等。

#### 实现步骤

**1. 获取和保存设置**

```javascript
async function getReadingSettings(token) {
  const response = await fetch(
    'http://localhost:8080/api/v1/reader/settings',
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  return await response.json();
}

async function updateReadingSettings(settings, token) {
  const response = await fetch(
    'http://localhost:8080/api/v1/reader/settings',
    {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(settings)
    }
  );
  return await response.json();
}
```

**2. 设置面板组件**

```jsx
function ReaderSettings() {
  const [settings, setSettings] = useState(null);
  const token = useAuth();
  
  useEffect(() => {
    getReadingSettings(token).then(data => {
      setSettings(data.data);
    });
  }, [token]);
  
  const updateSetting = async (key, value) => {
    const newSettings = { ...settings, [key]: value };
    setSettings(newSettings);
    
    // 保存到服务器
    await updateReadingSettings({ [key]: value }, token);
    
    // 应用设置
    applySettings(newSettings);
  };
  
  const applySettings = (settings) => {
    const contentEl = document.querySelector('.reader-content');
    if (contentEl) {
      contentEl.style.fontSize = `${settings.fontSize}px`;
      contentEl.style.fontFamily = settings.fontFamily;
      contentEl.style.lineHeight = settings.lineHeight;
      contentEl.style.backgroundColor = settings.backgroundColor;
      contentEl.style.color = settings.textColor;
    }
  };
  
  if (!settings) return null;
  
  return (
    <div className="settings-panel">
      {/* 字体大小 */}
      <div className="setting-item">
        <label>字体大小</label>
        <input
          type="range"
          min="12"
          max="24"
          value={settings.fontSize}
          onChange={(e) => updateSetting('fontSize', parseInt(e.target.value))}
        />
        <span>{settings.fontSize}px</span>
      </div>
      
      {/* 字体类型 */}
      <div className="setting-item">
        <label>字体</label>
        <select
          value={settings.fontFamily}
          onChange={(e) => updateSetting('fontFamily', e.target.value)}
        >
          <option value="宋体">宋体</option>
          <option value="黑体">黑体</option>
          <option value="楷体">楷体</option>
          <option value="微软雅黑">微软雅黑</option>
        </select>
      </div>
      
      {/* 主题 */}
      <div className="setting-item">
        <label>主题</label>
        <div className="theme-selector">
          <button 
            className={settings.theme === 'default' ? 'active' : ''}
            onClick={() => updateSetting('theme', 'default')}
          >
            默认
          </button>
          <button 
            className={settings.theme === 'night' ? 'active' : ''}
            onClick={() => updateSetting('theme', 'night')}
          >
            夜间
          </button>
          <button 
            className={settings.theme === 'sepia' ? 'active' : ''}
            onClick={() => updateSetting('theme', 'sepia')}
          >
            护眼
          </button>
        </div>
      </div>
      
      {/* 翻页模式 */}
      <div className="setting-item">
        <label>翻页模式</label>
        <select
          value={settings.pageMode}
          onChange={(e) => updateSetting('pageMode', e.target.value)}
        >
          <option value="scroll">滚动翻页</option>
          <option value="paginate">仿真翻页</option>
        </select>
      </div>
    </div>
  );
}
```

---

### 场景7：用户评分和评论

#### 需求
用户可以对书籍进行评分和评论。

#### 实现步骤

**1. 提交评分**

```javascript
async function submitRating(bookId, rating, review, token) {
  const response = await fetch(
    'http://localhost:8080/api/v1/reading/ratings',
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        bookId,
        rating,
        review
      })
    }
  );
  return await response.json();
}
```

**2. 评分组件**

```jsx
function RatingForm({ bookId }) {
  const [rating, setRating] = useState(0);
  const [review, setReview] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const token = useAuth();
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (rating === 0) {
      alert('请选择评分');
      return;
    }
    
    setSubmitting(true);
    try {
      await submitRating(bookId, rating, review, token);
      alert('评分成功！');
      // 刷新页面或评分列表
    } catch (error) {
      alert('评分失败：' + error.message);
    } finally {
      setSubmitting(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit} className="rating-form">
      <h3>给这本书评分</h3>
      
      {/* 星级评分 */}
      <div className="star-rating">
        {[1, 2, 3, 4, 5].map(star => (
          <span
            key={star}
            className={star <= rating ? 'star active' : 'star'}
            onClick={() => setRating(star)}
          >
            ★
          </span>
        ))}
      </div>
      
      {/* 评论文本 */}
      <textarea
        value={review}
        onChange={(e) => setReview(e.target.value)}
        placeholder="写下你的评价（可选）"
        rows={4}
      />
      
      <button type="submit" disabled={submitting}>
        {submitting ? '提交中...' : '提交评分'}
      </button>
    </form>
  );
}
```

---

## 🛠️ 工具函数库

### API请求封装

```javascript
// api.js
const BASE_URL = 'http://localhost:8080';

class APIClient {
  constructor(baseURL = BASE_URL) {
    this.baseURL = baseURL;
    this.token = null;
  }
  
  setToken(token) {
    this.token = token;
  }
  
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers
    };
    
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }
    
    try {
      const response = await fetch(url, {
        ...options,
        headers
      });
      
      const data = await response.json();
      
      if (data.code !== 200 && data.code !== 201) {
        throw new Error(data.message || '请求失败');
      }
      
      return data;
    } catch (error) {
      console.error('API请求错误:', error);
      throw error;
    }
  }
  
  get(endpoint, params = {}) {
    const query = new URLSearchParams(params).toString();
    const url = query ? `${endpoint}?${query}` : endpoint;
    return this.request(url);
  }
  
  post(endpoint, data) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }
  
  put(endpoint, data) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }
  
  delete(endpoint) {
    return this.request(endpoint, {
      method: 'DELETE'
    });
  }
}

// 创建全局实例
export const api = new APIClient();

// 使用示例
api.setToken(userToken);
const bookData = await api.get('/api/v1/bookstore/books/123');
```

### 防抖函数

```javascript
// utils.js
export function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}
```

### 本地存储

```javascript
// storage.js
export const storage = {
  set(key, value) {
    localStorage.setItem(key, JSON.stringify(value));
  },
  
  get(key) {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : null;
  },
  
  remove(key) {
    localStorage.removeItem(key);
  },
  
  clear() {
    localStorage.clear();
  }
};

// 使用示例
storage.set('userToken', token);
const token = storage.get('userToken');
```

---

## 📱 移动端适配

### React Native示例

```javascript
import AsyncStorage from '@react-native-async-storage/async-storage';

// API请求
async function fetchBookDetail(bookId) {
  try {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/bookstore/books/${bookId}`
    );
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('请求失败', error);
    throw error;
  }
}

// Token管理
async function saveToken(token) {
  await AsyncStorage.setItem('userToken', token);
}

async function getToken() {
  return await AsyncStorage.getItem('userToken');
}

// 阅读进度本地缓存
async function saveProgressLocal(bookId, progress) {
  const key = `progress_${bookId}`;
  await AsyncStorage.setItem(key, JSON.stringify(progress));
}
```

---

## ⚠️ 常见问题

### Q1: Token过期怎么处理？

```javascript
async function requestWithTokenRefresh(endpoint, options) {
  let token = storage.get('userToken');
  
  try {
    const response = await api.request(endpoint, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${token}`
      }
    });
    return response;
  } catch (error) {
    if (error.message.includes('token') || error.message.includes('401')) {
      // Token过期，刷新Token
      token = await refreshToken();
      storage.set('userToken', token);
      
      // 重试请求
      return await api.request(endpoint, {
        ...options,
        headers: {
          ...options.headers,
          'Authorization': `Bearer ${token}`
        }
      });
    }
    throw error;
  }
}
```

### Q2: 如何实现离线阅读？

```javascript
// 使用IndexedDB或LocalStorage缓存章节内容
async function cacheChapter(chapterId, content) {
  await localforage.setItem(`chapter_${chapterId}`, content);
}

async function getChapterFromCache(chapterId) {
  return await localforage.getItem(`chapter_${chapterId}`);
}

// 读取章节时先查缓存
async function readChapter(chapterId, token) {
  // 先从缓存读取
  let content = await getChapterFromCache(chapterId);
  
  if (content) {
    return content;
  }
  
  // 缓存未命中，从服务器获取
  const data = await getChapterContent(chapterId, token);
  content = data.data.content;
  
  // 缓存章节
  await cacheChapter(chapterId, content);
  
  return content;
}
```

### Q3: 如何优化大量章节列表的加载？

```javascript
// 使用虚拟滚动 + 分页加载
function ChapterListVirtual({ bookId }) {
  const [chapters, setChapters] = useState([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  
  const loadMore = async () => {
    if (!hasMore) return;
    
    const result = await getChapterList(bookId, page, 50);
    setChapters([...chapters, ...result.data]);
    setPage(page + 1);
    setHasMore(chapters.length < result.total);
  };
  
  return (
    <VirtualList
      data={chapters}
      itemHeight={60}
      onScrollEnd={loadMore}
      renderItem={(chapter) => (
        <ChapterItem chapter={chapter} />
      )}
    />
  );
}
```

---

## 🎯 性能优化建议

### 1. 请求合并

```javascript
// 合并多个请求
async function loadPageData(bookId) {
  const [book, chapters, ratings] = await Promise.all([
    getBookDetail(bookId),
    getChapterList(bookId),
    getBookRatingStats(bookId)
  ]);
  
  return { book, chapters, ratings };
}
```

### 2. 结果缓存

```javascript
// 简单的内存缓存
const cache = new Map();

async function getCachedData(key, fetcher, ttl = 5 * 60 * 1000) {
  const cached = cache.get(key);
  
  if (cached && Date.now() - cached.time < ttl) {
    return cached.data;
  }
  
  const data = await fetcher();
  cache.set(key, { data, time: Date.now() });
  
  return data;
}

// 使用
const book = await getCachedData(
  `book_${bookId}`,
  () => getBookDetail(bookId)
);
```

### 3. 图片懒加载

```jsx
function BookCover({ src, alt }) {
  const [loaded, setLoaded] = useState(false);
  const imgRef = useRef(null);
  
  useEffect(() => {
    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting) {
        setLoaded(true);
        observer.disconnect();
      }
    });
    
    if (imgRef.current) {
      observer.observe(imgRef.current);
    }
    
    return () => observer.disconnect();
  }, []);
  
  return (
    <div ref={imgRef} className="book-cover">
      {loaded && <img src={src} alt={alt} />}
    </div>
  );
}
```

---

## 📞 技术支持

遇到问题？

1. 查看完整API文档: [阅读端API使用文档.md](./阅读端API使用文档.md)
2. 提交Issue到项目仓库
3. 联系技术支持: dev@qingyu.com

---

**文档版本**: v1.0  
**最后更新**: 2025-10-16  
**维护者**: 青羽后端团队

