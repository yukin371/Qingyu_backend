# Recommendation API 模块结构说明

## 📁 文件结构

```
api/v1/recommendation/
├── recommendation_api.go    # 推荐API主入口
├── personal.go              # 个性化推荐
├── similar.go               # 相似推荐
└── README.md                # 本文件
```

## 🎯 模块职责

**职责**: 智能推荐系统，为用户提供个性化的书籍推荐

**核心功能**:
- ✅ 个性化推荐
- ✅ 相似书籍推荐
- ✅ 基于协同过滤的推荐
- ✅ 基于内容的推荐
- ✅ 热门推荐

---

## 📋 API端点列表

### 公开端点

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/recommendation/hot | 热门推荐 | RecommendationAPI.GetHotRecommendations |
| GET | /api/v1/recommendation/similar/:bookId | 相似书籍推荐 | SimilarAPI.GetSimilarBooks |
| GET | /api/v1/recommendation/category/:categoryId | 分类推荐 | RecommendationAPI.GetCategoryRecommendations |

### 认证端点（需要JWT Token）

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/recommendation/personal | 个性化推荐 | PersonalAPI.GetPersonalRecommendations |
| GET | /api/v1/recommendation/for-you | 猜你喜欢 | PersonalAPI.GetForYouRecommendations |
| GET | /api/v1/recommendation/based-on/:bookId | 基于某本书的推荐 | PersonalAPI.GetBasedOnRecommendations |

---

## 🔄 推荐算法

### 1. 个性化推荐（Personal Recommendations）

**算法**: 协同过滤 + 内容过滤

**输入**:
- 用户阅读历史
- 用户评分记录
- 用户收藏记录
- 用户偏好设置

**输出**:
- 推荐书籍列表
- 推荐理由
- 推荐评分

**流程**:
```
1. 获取用户画像
2. 分析用户行为
3. 查找相似用户
4. 获取相似用户喜欢的书籍
5. 内容过滤和排序
6. 返回推荐结果
```

---

### 2. 相似书籍推荐（Similar Books）

**算法**: 内容相似度计算

**输入**:
- 目标书籍ID
- 书籍特征（分类、标签、作者等）

**输出**:
- 相似书籍列表
- 相似度评分

**相似度计算维度**:
- 分类相似度（40%）
- 标签相似度（30%）
- 作者相似度（15%）
- 阅读群体相似度（15%）

---

### 3. 热门推荐（Hot Recommendations）

**算法**: 热度计算

**热度因子**:
- 最近阅读量（35%）
- 收藏数（25%）
- 评分（20%）
- 评论数（10%）
- 分享数（10%）

**时间衰减**:
热度值随时间衰减，保证推荐的时效性。

---

## 📊 请求/响应示例

### 个性化推荐
```json
GET /api/v1/recommendation/personal?limit=10
Authorization: Bearer <token>

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "recommendations": [
      {
        "book_id": "book_123",
        "title": "仙逆",
        "author": "耳根",
        "cover": "https://example.com/cover.jpg",
        "score": 9.5,
        "reason": "基于你的阅读历史，你可能喜欢这本书",
        "tags": ["修仙", "玄幻"],
        "similarity": 0.95
      }
    ],
    "total": 50,
    "algorithm": "collaborative_filtering"
  }
}
```

### 相似书籍推荐
```json
GET /api/v1/recommendation/similar/book_456?limit=5

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "source_book": {
      "book_id": "book_456",
      "title": "凡人修仙传"
    },
    "similar_books": [
      {
        "book_id": "book_789",
        "title": "仙逆",
        "author": "耳根",
        "cover": "https://example.com/cover.jpg",
        "similarity": 0.92,
        "common_tags": ["修仙", "凡人流"],
        "reason": "同为经典修仙小说，主角设定相似"
      }
    ]
  }
}
```

### 热门推荐
```json
GET /api/v1/recommendation/hot?category=玄幻&limit=10

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "hot_books": [
      {
        "book_id": "book_101",
        "title": "遮天",
        "author": "辰东",
        "cover": "https://example.com/cover.jpg",
        "hot_score": 98.5,
        "recent_views": 100000,
        "favorites": 50000,
        "rating": 9.3,
        "trend": "rising"
      }
    ],
    "update_time": "2025-10-24T10:00:00Z"
  }
}
```

---

## 🔧 设计原则

### 1. 多样性
避免推荐结果过于单一，保证内容多样性。

### 2. 实时性
定期更新推荐结果，反映最新的用户行为。

### 3. 可解释性
提供推荐理由，让用户理解为什么推荐这本书。

### 4. 隐私保护
推荐算法不泄露用户隐私数据。

### 5. 性能优化
- 推荐结果缓存
- 离线计算
- 增量更新

---

## 📝 开发规范

### 1. 推荐结果缓存

```go
// 个性化推荐缓存1小时
cacheKey := fmt.Sprintf("recommendation:personal:%s", userID)
cache.Set(cacheKey, recommendations, time.Hour)

// 热门推荐缓存30分钟
cacheKey := "recommendation:hot"
cache.Set(cacheKey, hotBooks, 30*time.Minute)
```

### 2. 推荐评分计算

```go
// 综合评分
score := 0.0
score += categoryScore * 0.4
score += tagScore * 0.3
score += authorScore * 0.15
score += audienceScore * 0.15
```

### 3. 冷启动处理

对于新用户或没有足够数据的用户：
- 使用热门推荐
- 基于注册时选择的兴趣标签
- 推荐高质量的经典作品

---

## 🚀 扩展建议

### 未来可添加的功能

1. **深度学习推荐**
   - 基于神经网络的推荐算法
   - 多模态特征融合
   - 实时学习用户偏好

2. **上下文感知推荐**
   - 时间上下文（早晚推荐不同类型）
   - 位置上下文
   - 设备上下文

3. **社交推荐**
   - 好友推荐
   - 社区热门
   - 达人推荐

4. **A/B测试**
   - 多种推荐算法对比
   - 效果评估
   - 动态调整权重

5. **推荐解释**
   - 更详细的推荐理由
   - 可视化展示
   - 用户反馈机制

---

## 📊 推荐效果评估

### 评估指标

1. **准确性指标**
   - 点击率（CTR）
   - 转化率（CVR）
   - 收藏率

2. **多样性指标**
   - 推荐列表多样性
   - 覆盖率

3. **新颖性指标**
   - 推荐新书比例
   - 长尾覆盖率

4. **用户满意度**
   - 用户反馈评分
   - 停留时间
   - 加入书架率

---

## 🔄 与其他模块的关系

### Recommendation vs Bookstore

| 功能 | Recommendation（推荐） | Bookstore（书城） |
|------|----------------------|------------------|
| **定位** | 智能推荐 | 书籍展示 |
| **算法** | 个性化算法 | 基础筛选排序 |
| **数据源** | 用户行为数据 | 书籍元数据 |
| **目标** | 提高发现效率 | 展示所有书籍 |

### Recommendation vs Reader

| 功能 | Recommendation（推荐） | Reader（阅读器） |
|------|----------------------|-----------------|
| **时机** | 选书阶段 | 阅读阶段 |
| **数据** | 消费推荐结果 | 生产行为数据 |
| **关系** | Reader行为驱动推荐 | 推荐引导阅读 |

---

## 📚 相关文档

- [推荐算法设计](../../../doc/design/recommendation/algorithm_design.md)
- [用户画像构建](../../../doc/design/recommendation/user_profile.md)
- [Bookstore API模块](../bookstore/README.md)
- [Reader API模块](../reader/README.md)

---

## 💡 最佳实践

### 1. 推荐策略

```
- 新用户: 热门 + 分类热门
- 活跃用户: 个性化 + 协同过滤
- 流失用户: 新书 + 热门回归
```

### 2. 缓存策略

```
- L1 Cache: 内存缓存（5分钟）
- L2 Cache: Redis缓存（1小时）
- L3 Cache: 数据库（离线计算结果）
```

### 3. 降级策略

```
- 算法服务不可用 → 降级到热门推荐
- 缓存失效 → 降级到分类推荐
- 数据不足 → 降级到编辑推荐
```

---

**版本**: v1.0  
**创建日期**: 2025-10-24  
**维护者**: Recommendation模块开发组

