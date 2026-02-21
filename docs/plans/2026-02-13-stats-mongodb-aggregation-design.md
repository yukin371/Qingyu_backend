# Stats MongoDB 聚合管道设计

**日期**: 2026-02-13  
**范围**: `service/shared/stats`  
**目标**: 为后续 `Task 2.4/2.5` 提供可直接落地的 MongoDB 聚合方案，先完成设计，不改业务行为。

## 设计原则

1. 只做读聚合，不在查询链路写回业务数据。  
2. 所有时间窗口统一使用 `[startDate, endDate)` 半开区间。  
3. 高基数维度先聚合后投影，避免在应用层做大规模 map/reduce。  
4. pipeline 输出字段与现有 `stats_service.go` DTO 一一映射。  

## 数据源与索引建议

1. 用户维度
- 集合：`users`
- 关键字段：`_id`, `created_at`, `last_login_at`, `is_vip`, `status`
- 索引建议：`{created_at:1}`, `{last_login_at:1}`, `{status:1,is_vip:1}`

2. 内容维度
- 集合：`books`, `chapters`
- 关键字段：`author_id`, `status`, `category`, `word_count`, `created_at`
- 索引建议：`books:{author_id:1,status:1,created_at:1}`, `chapters:{book_id:1,word_count:1}`

3. 行为与收益（后续接入）
- 集合建议：`reading_events`, `wallet_transactions`
- 索引建议：`{user_id:1,created_at:1}`, `{author_id:1,type:1,created_at:1}`

## 管道 A：平台用户统计（PlatformUserStats）

目标字段：`TotalUsers/NewUsers/ActiveUsers/VIPUsers/RetentionRate/AverageActiveDay`

Pipeline（users）：
```javascript
[
  {
    $facet: {
      total_users: [{ $count: "value" }],
      new_users: [
        { $match: { created_at: { $gte: startDate, $lt: endDate } } },
        { $count: "value" }
      ],
      active_users: [
        { $match: { last_login_at: { $gte: startDate, $lt: endDate } } },
        { $count: "value" }
      ],
      vip_users: [
        { $match: { is_vip: true, status: "active" } },
        { $count: "value" }
      ]
    }
  },
  {
    $project: {
      total_users: { $ifNull: [{ $arrayElemAt: ["$total_users.value", 0] }, 0] },
      new_users: { $ifNull: [{ $arrayElemAt: ["$new_users.value", 0] }, 0] },
      active_users: { $ifNull: [{ $arrayElemAt: ["$active_users.value", 0] }, 0] },
      vip_users: { $ifNull: [{ $arrayElemAt: ["$vip_users.value", 0] }, 0] }
    }
  }
]
```

计算规则（应用层）：
- `RetentionRate = active_users / max(total_users,1)`  
- `AverageActiveDay` 在未接入 `reading_events` 前保持占位 0（与当前行为一致）。

## 管道 B：平台内容统计（PlatformContentStats）

目标字段：`TotalBooks/NewBooks/TotalChapters/TotalWords/TotalViews/AverageRating/PopularCategories`

Pipeline（books + chapters）：
```javascript
[
  {
    $facet: {
      books_all: [{ $count: "value" }],
      books_new: [
        { $match: { created_at: { $gte: startDate, $lt: endDate } } },
        { $count: "value" }
      ],
      categories: [
        { $match: { status: "published" } },
        { $group: { _id: "$category", count: { $sum: 1 } } },
        { $sort: { count: -1 } },
        { $limit: 10 }
      ]
    }
  }
]
```

补充聚合（chapters）：
```javascript
[
  { $group: {
      _id: null,
      total_chapters: { $sum: 1 },
      total_words: { $sum: { $ifNull: ["$word_count", 0] } }
  }}
]
```

说明：
- `TotalViews`、`AverageRating` 依赖阅读/评分数据源，未接入前保持占位值。

## 管道 C：作者内容统计（ContentStats by userID）

目标字段：`TotalBooks/PublishedBooks/TotalChapters/TotalWords/TotalViews/...`

Pipeline（books）：
```javascript
[
  { $match: { author_id: userID } },
  {
    $facet: {
      total_books: [{ $count: "value" }],
      published_books: [
        { $match: { status: "published" } },
        { $count: "value" }
      ],
      book_ids: [{ $project: { _id: 1 } }]
    }
  }
]
```

随后按 `book_ids` 在 `chapters` 侧聚合 `total_chapters/total_words`。

## 管道 D：收益统计（RevenueStats，后续接入）

集合：`wallet_transactions`
```javascript
[
  { $match: {
      author_id: userID,
      created_at: { $gte: startDate, $lt: endDate },
      status: "success"
  }},
  { $group: {
      _id: "$type",
      amount: { $sum: "$amount" }
  }}
]
```

说明：当前服务未接入 wallet 聚合，该管道仅作为后续实现输入。

## 应用层映射约束

1. 所有缺省值使用 0/空数组，保持当前 API 输出兼容。  
2. 聚合错误按“可降级字段”与“硬失败字段”区分：  
- 硬失败：`GetUserStats/GetContentStats` 主查询失败直接返回错误。  
- 可降级：views/rating/revenue 等未接入字段回落 0。  

## 分阶段落地建议

1. 第一步：接入 A/B/C 主聚合（不改 DTO）。  
2. 第二步：接入 `reading_events`，补齐 `TotalViews/ActiveDays`。  
3. 第三步：接入 `wallet_transactions`，补齐 `RevenueStats`。  

