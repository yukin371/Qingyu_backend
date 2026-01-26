# Search Module

统一搜索模块，提供书籍、项目、文档、用户和向量搜索能力。

## 目录结构

```
service/search/
├── search.go                  # SearchService 主服务
├── config.go                  # 搜索配置
├── query_optimizer.go         # 查询优化器
├── engine/                    # 搜索引擎实现
│   ├── engine.go             # Engine 接口定义
│   ├── elasticsearch.go      # ES 引擎
│   ├── milvus.go             # Milvus 引擎
│   └── mongodb.go            # MongoDB 兼容引擎
├── provider/                  # 业务搜索提供者
│   ├── provider.go           # Provider 接口定义
│   ├── book_provider.go      # 书籍搜索
│   ├── project_provider.go   # 项目搜索
│   ├── document_provider.go  # 文档搜索
│   ├── user_provider.go      # 用户搜索
│   └── vector_provider.go    # 向量搜索
├── cache/                     # 缓存管理
│   └── search_cache.go       # 搜索缓存
└── sync/                      # 数据同步
    ├── change_stream.go      # Change Stream 监听器
    └── sync_worker.go        # 同步 Worker
```

## 设计模式

### 分层架构

1. **API 层** (`api/v1/search/`)
   - 提供统一的 REST API 接口
   - 参数验证和转换
   - 响应格式化

2. **Service 层** (`service/search/`)
   - SearchService: 统一搜索入口
   - Provider: 业务搜索提供者
   - Engine: 搜索引擎抽象

3. **Repository 层** (`repository/search/`)
   - Elasticsearch Repository
   - Milvus Repository
   - Cache Repository

4. **Models 层** (`models/search/`)
   - 统一的请求/响应模型
   - 错误定义
   - 配置模型

## 核心接口

### Engine 接口

```go
type Engine interface {
    Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error)
    Index(ctx context.Context, index string, documents []Document) error
    Update(ctx context.Context, index string, id string, document Document) error
    Delete(ctx context.Context, index string, id string) error
    CreateIndex(ctx context.Context, index string, mapping interface{}) error
    Health(ctx context.Context) error
}
```

### Provider 接口

```go
type Provider interface {
    Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
    Type() SearchType
    Validate(req *SearchRequest) error
}
```

## 搜索类型

| 类型 | 路径 | 认证 | 说明 |
|------|------|------|------|
| books | `/api/v1/search/books` | 无需认证 | 公开书籍搜索 |
| projects | `/api/v1/search/projects` | 需要认证 | 创作项目搜索 |
| documents | `/api/v1/search/documents` | 需要认证 | 文档搜索 |
| users | `/api/v1/search/users` | 无需认证 | 用户搜索 |
| vector | `/api/v1/search/vector` | 需要认证 | 向量搜索 |

## 权限控制

### 书籍搜索
- 强制过滤: `status = published` 且 `is_private = false`

### 项目搜索
- 强制过滤: `author_id = 当前用户ID`

### 文档搜索
- 强制过滤: `user_id = 当前用户ID`
- 可选过滤: `project_id`

## 数据同步

```
MongoDB (主数据库)
    ↓ Change Stream
Redis Queue
    ↓ Worker
Elasticsearch (搜索引擎)
```

## 配置

```yaml
search:
  cache:
    enabled: true
    default_ttl: 300s
    hot_ttl: 600s
  rate_limit:
    enabled: true
    requests_per_minute: 60
  optimizer:
    max_results: 10000
    max_page_size: 100
    min_query_length: 2
  books:
    allowed_statuses: ["completed", "published", "serializing"]
    allowed_privacy: [false]
```

## TODO

- [ ] 实现 Elasticsearch Engine
- [ ] 实现 Milvus Engine
- [ ] 实现 MongoEngine
- [ ] 实现 BookProvider
- [ ] 实现 ProjectProvider
- [ ] 实现 DocumentProvider
- [ ] 实现 UserProvider
- [ ] 实现 VectorProvider
- [ ] 实现搜索缓存
- [ ] 实现 Change Stream 监听
- [ ] 实现同步 Worker
- [ ] 集成到主路由
- [ ] 单元测试
- [ ] 集成测试

## 参考

- [统一搜索模块设计文档](../../../docs/plans/2026-01-25-unified-search-design.md)
