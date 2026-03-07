# Redis缓存策略优化指南

> **文档版本**: v1.0
> **创建日期**: 2026-01-06
> **适用系统**: 青羽写作平台后端

## 📋 概述

本文档说明如何优化青羽平台的Redis缓存策略，以提升系统性能和降低数据库压力。

## 🎯 缓存设计原则

### 1. 缓存使用原则

| 原则 | 说明 | 示例 |
|------|------|------|
| **热点数据优先** | 频繁访问的数据优先缓存 | 首页推荐、热门书籍 |
| **读多写少** | 读取频繁、修改少的数据适合缓存 | 书籍详情、用户信息 |
| **计算昂贵** | 复杂计算结果应该缓存 | 搜索结果、统计数据 |
| **数据一致性** | 根据一致性要求选择缓存策略 | 实时数据不缓存、允许短暂延迟则缓存 |

### 2. 缓存穿透防护

```go
// 方案1: 布隆过滤器
func (s *BookService) GetBook(id string) (*Book, error) {
    // 1. 检查布隆过滤器
    exists := s.bloomFilter.Exists(id)
    if !exists {
        return nil, ErrBookNotFound
    }

    // 2. 查询缓存
    book, err := s.cache.Get(ctx, "book:"+id)
    if err == nil {
        return book, nil
    }

    // 3. 查询数据库
    book, err = s.repository.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 4. 写入缓存
    s.cache.Set(ctx, "book:"+id, book, 1*time.Hour)
    return book, nil
}

// 方案2: 缓存空值
func (s *BookService) GetBook(id string) (*Book, error) {
    // 1. 查询缓存
    val, err := s.cache.Get(ctx, "book:"+id)
    if err == nil {
        if val == "NULL" {
            return nil, ErrBookNotFound
        }
        return decodeBook(val), nil
    }

    // 2. 查询数据库
    book, err := s.repository.Get(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 缓存空值，防止穿透
            s.cache.Set(ctx, "book:"+id, "NULL", 5*time.Minute)
            return nil, ErrBookNotFound
        }
        return nil, err
    }

    // 3. 写入缓存
    s.cache.Set(ctx, "book:"+id, encodeBook(book), 1*time.Hour)
    return book, nil
}
```

### 3. 缓存击穿防护

```go
// 方案1: 互斥锁 (Mutex)
func (s *BookService) GetBook(id string) (*Book, error) {
    // 1. 查询缓存
    book, err := s.cache.Get(ctx, "book:"+id)
    if err == nil {
        return decodeBook(book), nil
    }

    // 2. 获取分布式锁
    lockKey := "lock:book:" + id
    lock, err := s.locker.Acquire(ctx, lockKey, 10*time.Second)
    if err != nil {
        // 获取锁失败，等待片刻后重试
        time.Sleep(100 * time.Millisecond)
        return s.GetBook(id) // 递归重试
    }
    defer lock.Release()

    // 3. 双重检查：其他线程可能已经加载了缓存
    book, err = s.cache.Get(ctx, "book:"+id)
    if err == nil {
        return decodeBook(book), nil
    }

    // 4. 查询数据库
    book, err = s.repository.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 5. 写入缓存
    s.cache.Set(ctx, "book:"+id, encodeBook(book), 1*time.Hour)
    return book, nil
}

// 方案2: 逻辑过期 (Logical Expiration)
type CacheValue struct {
    Data      []byte `json:"data"`
    ExpireAt int64  `json:"expire_at"` // 逻辑过期时间
}

func (s *BookService) GetBook(id string) (*Book, error) {
    // 1. 查询缓存
    val, err := s.cache.Get(ctx, "book:"+id)
    if err == nil {
        cacheVal := decodeCacheValue(val)
        if time.Now().Unix() < cacheVal.ExpireAt {
            // 未过期，直接返回
            return decodeBook(cacheVal.Data), nil
        }

        // 已过期，异步刷新
        go func() {
            s.refreshBookCache(id)
        }()

        // 返回过期数据（保证可用性）
        return decodeBook(cacheVal.Data), nil
    }

    // 2. 缓存不存在，直接查询数据库
    book, err := s.repository.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    cacheVal := &CacheValue{
        Data:      encodeBook(book),
        ExpireAt: time.Now().Add(1*time.Hour).Unix(),
    }
    s.cache.Set(ctx, "book:"+id, encodeCacheValue(cacheVal), 1*time.Hour+10*time.Minute)
    return book, nil
}
```

### 4. 缓存雪崩防护

```go
// 方案1: 随机过期时间
func (s *BookService) SetBookCache(id string, book *Book) error {
    // 基础过期时间 1 小时
    baseTTL := 1 * time.Hour
    // 随机增加 0-10 分钟
    randomTTL := time.Duration(rand.Intn(600)) * time.Second

    finalTTL := baseTTL + randomTTL
    return s.cache.Set(ctx, "book:"+id, encodeBook(book), finalTTL)
}

// 方案2: 多级缓存
func (s *BookService) GetBook(id string) (*Book, error) {
    // L1: 本地缓存（内存）
    if book, ok := s.localCache.Get("book:" + id); ok {
        return book, nil
    }

    // L2: Redis缓存
    val, err := s.cache.Get(ctx, "book:"+id)
    if err == nil {
        book := decodeBook(val)
        s.localCache.Set("book:"+id, book, 5*time.Minute)
        return book, nil
    }

    // L3: 数据库
    book, err := s.repository.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 回写缓存
    s.cache.Set(ctx, "book:"+id, encodeBook(book), 1*time.Hour)
    s.localCache.Set("book:"+id, book, 5*time.Minute)
    return book, nil
}
```

## 📊 缓存策略分类

### 1. 按数据类型分类

| 数据类型 | 缓存策略 | 过期时间 | 更新策略 |
|---------|---------|---------|---------|
| **热点书籍** | Cache Aside + 预热 | 1小时 | 定时刷新 |
| **用户信息** | Cache Aside | 30分钟 | 写入时更新 |
| **阅读进度** | Write Through | 永久 | 实时更新 |
| **统计数据** | Write Behind | 5分钟 | 批量更新 |
| **搜索结果** | Cache Aside | 10分钟 | 失效后重建 |
| **推荐列表** | Cache Aside + 预热 | 15分钟 | 定时刷新 |
| **配置信息** | Cache Aside | 1小时 | 主动刷新 |
| **限流计数** | Redis INCR | 动态 | 滑动窗口 |

### 2. 按更新频率分类

| 频率 | 策略 | 过期时间 | 示例 |
|------|------|---------|------|
| **实时更新** | 不缓存或短时间缓存 | 1-5分钟 | 在线人数、阅读进度 |
| **频繁更新** | Cache Aside + 短过期 | 10-30分钟 | 点赞数、评论数 |
| **偶尔更新** | Cache Aside + 长过期 | 1-6小时 | 书籍详情、章节内容 |
| **极少更新** | Cache Aside + 超长过期 | 1-7天 | 系统配置、分类列表 |

### 3. 按数据大小分类

| 大小 | 存储方式 | 示例 |
|------|---------|------|
| **小数据 (<1KB)** | String | 用户ID、配置项 |
| **中等数据 (1KB-100KB)** | String/Hash | 书籍详情、用户信息 |
| **大数据 (100KB-1MB)** | Hash/压缩 | 章节内容、长文本 |
| **超大数据 (>1MB)** | 分片存储或不要缓存 | 完整书稿、大文件 |

## 🚀 缓存策略实现

### 1. Cache Aside (旁路缓存)

**适用场景**: 读多写少的数据

```go
// 读取
func (s *BookService) GetBook(id string) (*Book, error) {
    // 1. 查询缓存
    val, err := s.cache.Get(ctx, "book:"+id)
    if err == nil {
        return decodeBook(val), nil
    }

    // 2. 查询数据库
    book, err := s.repository.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    s.cache.Set(ctx, "book:"+id, encodeBook(book), 1*time.Hour)
    return book, nil
}

// 更新
func (s *BookService) UpdateBook(book *Book) error {
    // 1. 更新数据库
    err := s.repository.Update(ctx, book)
    if err != nil {
        return err
    }

    // 2. 删除缓存（而不是更新缓存）
    s.cache.Delete(ctx, "book:"+book.ID)

    return nil
}
```

### 2. Write Through (直写缓存)

**适用场景**: 需要强一致性的数据

```go
func (s *BookService) UpdateBook(book *Book) error {
    // 同时更新缓存和数据库
    err := s.cache.Set(ctx, "book:"+book.ID, encodeBook(book), 1*time.Hour)
    if err != nil {
        return err
    }

    err = s.repository.Update(ctx, book)
    if err != nil {
        // 回滚缓存
        s.cache.Delete(ctx, "book:"+book.ID)
        return err
    }

    return nil
}
```

### 3. Write Behind (异步写回)

**适用场景**: 写入频繁、可接受短暂不一致的数据

```go
func (s *BookService) IncrementReadCount(bookID string) error {
    // 1. 只更新缓存
    key := "book:" + bookID + ":read_count"
    _, err := s.cache.Increment(ctx, key)
    if err != nil {
        return err
    }

    // 2. 异步批量写入数据库
    go func() {
        s.batchUpdateReadCount()
    }()

    return nil
}

func (s *BookService) batchUpdateReadCount() {
    // 每10秒或积累100条更新时批量写入
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.flushReadCount()
        }
    }
}
```

### 4. 缓存预热 (Cache Warm-up)

```go
func (s *BookService) WarmUpCache() error {
    // 1. 获取热门书籍列表
    hotBooks, err := s.repository.GetHotBooks(ctx, 1000)
    if err != nil {
        return err
    }

    // 2. 批量加载到缓存
    pipe := s.cache.Pipeline()
    for _, book := range hotBooks {
        key := "book:" + book.ID
        pipe.Set(ctx, key, encodeBook(book), 1*time.Hour)
    }

    // 3. 执行批量操作
    _, err = pipe.Exec(ctx)
    return err
}
```

## 🔧 缓存Key设计

### 1. Key命名规范

```go
// 格式: {业务模块}:{数据类型}:{唯一标识}[:{子项}]

// 示例
"user:info:123"              // 用户信息
"user:session:123"           // 用户会话
"book:detail:456"            // 书籍详情
"book:chapter:456:1"         // 章节1
"book:catalog:456"           // 目录
"book:hot:novel"             // 热门小说
"search:result:abc123"       // 搜索结果
"stats:daily:2026-01-06"     // 每日统计
"lock:book:456"              // 分布式锁
```

### 2. Key过期策略

| Key模式 | 过期时间 | 说明 |
|---------|---------|------|
| `user:session:*` | 7天 | 用户会话 |
| `user:info:*` | 30分钟 | 用户信息 |
| `book:detail:*` | 1小时 | 书籍详情 |
| `book:content:*` | 6小时 | 章节内容 |
| `book:hot:*` | 15分钟 | 热门列表 |
| `search:result:*` | 10分钟 | 搜索结果 |
| `stats:*` | 5分钟 | 统计数据 |
| `lock:*` | 10秒 | 分布式锁 |

### 3. Hash vs String 选择

```go
// 使用 Hash 的场景：
// 1. 对象的部分字段需要更新
// 2. 对象字段较多且经常单独访问
// 3. 需要获取部分字段

func (s *UserService) UpdateUserField(userID, field, value string) error {
    key := "user:info:" + userID
    return s.cache.HSet(ctx, key, field, value)
}

// 使用 String 的场景：
// 1. 对象整体读写
// 2. 对象较小
// 3. 需要原子性更新

func (s *BookService) GetBook(id string) (*Book, error) {
    key := "book:" + id
    val, err := s.cache.Get(ctx, key)
    if err != nil {
        return nil, err
    }
    return json.Unmarshal(val)
}
```

## 📈 缓存监控

### 1. 缓存命中率监控

```go
// 缓存统计
type CacheStats struct {
    Hits   int64 `json:"hits"`
    Misses int64 `json:"misses"`
}

func (s *CacheService) RecordHit() {
    atomic.AddInt64(&s.stats.Hits, 1)
}

func (s *CacheService) RecordMiss() {
    atomic.AddInt64(&s.stats.Misses, 1)
}

func (s *CacheService) GetHitRate() float64 {
    hits := atomic.LoadInt64(&s.stats.Hits)
    misses := atomic.LoadInt64(&s.stats.Misses)
    total := hits + misses
    if total == 0 {
        return 0
    }
    return float64(hits) / float64(total) * 100
}

// Prometheus 指标
var (
    cacheHits = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache_type"},
    )

    cacheMisses = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_misses_total",
            Help: "Total number of cache misses",
        },
        []string{"cache_type"},
    )

    cacheHitRate = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "cache_hit_rate",
            Help: "Cache hit rate percentage",
        },
        []string{"cache_type"},
    )
)

func (s *CacheService) recordMetrics(cacheType string) {
    hitRate := s.GetHitRate()
    cacheHits.WithLabelValues(cacheType).Inc()
    cacheHitRate.WithLabelValues(cacheType).Set(hitRate)
}
```

### 2. 缓存健康检查

```go
func (s *CacheService) HealthCheck() error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Ping检查
    err := s.client.Ping(ctx).Err()
    if err != nil {
        return fmt.Errorf("cache ping failed: %w", err)
    }

    // 读写测试
    testKey := "health:check"
    testValue := "ok"

    err = s.client.Set(ctx, testKey, testValue, 10*time.Second).Err()
    if err != nil {
        return fmt.Errorf("cache write failed: %w", err)
    }

    val, err := s.client.Get(ctx, testKey).Result()
    if err != nil || val != testValue {
        return fmt.Errorf("cache read failed: %w", err)
    }

    return nil
}
```

## 🎯 最佳实践

### 1. 缓存使用检查清单

- [ ] 所有热点数据都进行了缓存
- [ ] 缓存Key设计合理，包含业务前缀
- [ ] 缓存过期时间根据数据特性设置
- [ ] 实现了缓存穿透防护
- [ ] 实现了缓存击穿防护
- [ ] 实现了缓存雪崩防护
- [ ] 缓存更新策略选择正确
- [ ] 缓存数据大小合理
- [ ] 实现了缓存命中率监控
- [ ] 实现了缓存健康检查

### 2. 性能优化建议

1. **使用Pipeline批量操作**
```go
pipe := s.client.Pipeline()
for _, item := range items {
    pipe.Set(ctx, item.Key, item.Value, item.TTL)
}
_, err := pipe.Exec(ctx)
```

2. **压缩大对象**
```go
import "github.com/klauspost/compress/gzip"

func compress(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    gw := gzip.NewWriter(&buf)
    _, err := gw.Write(data)
    if err != nil {
        return nil, err
    }
    gw.Close()
    return buf.Bytes(), nil
}
```

3. **使用本地缓存减少Redis访问**
```go
import "github.com/patrickmn/go-cache"

localCache := cache.New(5*time.Minute, 10*time.Minute)
```

### 3. 常见错误避免

❌ **错误1**: 将大对象完整缓存
```go
// 不推荐
s.cache.Set(ctx, "book:123", entireBookWithChapters, ttl)
```

✅ **正确**: 只缓存必要字段
```go
// 推荐
s.cache.Set(ctx, "book:123", bookMetadata, ttl)
s.cache.Set(ctx, "book:123:catalog", chapterCatalog, ttl)
```

❌ **错误2**: 缓存频繁变化的数据
```go
// 不推荐
s.cache.Set(ctx, "online_count", count, ttl)  // 每秒变化
```

✅ **正确**: 使用计数器或实时查询
```go
// 推荐
s.client.Incr(ctx, "online_count")  // 原子操作
```

❌ **错误3**: 缓存Key冲突
```go
// 不推荐
s.cache.Set(ctx, id, data, ttl)  // 直接使用ID作为Key
```

✅ **正确**: 使用命名空间
```go
// 推荐
s.cache.Set(ctx, fmt.Sprintf("book:%s", id), data, ttl)
```

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2026-01-06
