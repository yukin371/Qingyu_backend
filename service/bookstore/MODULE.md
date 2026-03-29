# Bookstore Service

> 最后更新：2026-03-29

## 职责

书城核心业务层，管理书籍展示、分类体系、搜索、排行榜、Banner 和首页数据聚合。不处理创作、发布流程（由 Writer/Publish 负责）。

## 数据流

```
API Handler → BookstoreServiceImpl → Repository → MongoDB
                ↓
         SearchService（搜索委托，Milvus 全文检索）
                ↓
         EventBus（搜索事件、浏览事件）
```

## 约定 & 陷阱

- **搜索委托模式**：BookstoreService 不直接实现搜索，通过 `SetSearchService` 注入外部 SearchService，调用前必须检查是否已注入
- **排行榜策略**：`rankingStrategy` 函数根据排行榜类型选择不同计算公式，新增排行类型必须同步添加策略
- **首页数据聚合**：`GetHomepageData` 一次调用聚合 Banner + 推荐 + 热门 + 新书，注意并发控制和超时
- **书籍状态过滤**：所有列表查询默认只返回 `published` 状态的书籍，未发布书不会出现在书城
